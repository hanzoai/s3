// Package certreload provides a native, gRPC-free reloading X.509 keypair so
// both TLS servers (s3/security) and TLS clients (s3/util/http/client) can share
// one hot-reload implementation without an import cycle between them.
//
// It replaces the former dependency on grpc's
// credentials/tls/certprovider + pemfile (which S3 used purely as a file
// watcher, never for gRPC transport): a *Provider stats the cert/key files on a
// refresh tick and re-parses only when their mtime/size changes, so the hot path
// (the GetCertificate / GetClientCertificate callback fired on each TLS
// handshake) stays a cheap cached read. Rotated certs (e.g. from Kubernetes
// cert-manager / Vault) are picked up without a restart.
//
// Lives in its own subpackage because s3/security already imports
// s3/util/http/client (for LoadHTTPClientFromFile).
package certreload

import (
	"crypto/tls"
	"fmt"
	"os"
	"sync"
	"time"
)

// RefreshIntervalEnv names an environment variable that overrides the
// refresh cadence. Accepts any time.ParseDuration value (e.g. "30m",
// "500ms"). Primarily a hook for integration tests that need rotation
// to complete in seconds, but also useful in production when paired
// with short-lived certs (e.g. Vault-issued).
const RefreshIntervalEnv = "WEED_TLS_CERT_REFRESH_INTERVAL"

// DefaultRefreshInterval is the cadence at which the provider stats the
// cert/key files on disk. It re-parses only when mtime/size change, so the hot
// path (current() on each TLS handshake) stays cheap.
//
// 5 hours matches the prior constant used for gRPC mTLS. Resolved once
// at process start from RefreshIntervalEnv if set.
var DefaultRefreshInterval = resolveRefreshInterval(5 * time.Hour)

func resolveRefreshInterval(fallback time.Duration) time.Duration {
	if s := os.Getenv(RefreshIntervalEnv); s != "" {
		if d, err := time.ParseDuration(s); err == nil && d > 0 {
			return d
		}
	}
	return fallback
}

// Provider holds a parsed X.509 keypair and reloads it from disk on each refresh
// tick when the underlying files change. It is the native, gRPC-free replacement
// for grpc's certprovider.Provider; consumers hold it only to Close() the
// background refresh at shutdown.
type Provider struct {
	certFile, keyFile string

	mu       sync.RWMutex
	cert     *tls.Certificate
	certMod  time.Time
	certSize int64
	keyMod   time.Time
	keySize  int64

	stop     chan struct{}
	stopOnce sync.Once
}

func newProvider(certFile, keyFile string, refresh time.Duration) (*Provider, error) {
	if certFile == "" || keyFile == "" {
		return nil, fmt.Errorf("both certFile and keyFile are required")
	}
	p := &Provider{
		certFile: certFile,
		keyFile:  keyFile,
		stop:     make(chan struct{}),
	}
	// Load once eagerly so a bad cert/key fails fast at construction (parity with
	// pemfile.NewProvider, which errored on an unreadable initial pair).
	if err := p.reload(); err != nil {
		return nil, err
	}
	go p.refreshLoop(refresh)
	return p, nil
}

// reload re-reads and re-parses the keypair if either file changed since the
// last load. Safe to call concurrently with current().
func (p *Provider) reload() error {
	certInfo, err := os.Stat(p.certFile)
	if err != nil {
		return fmt.Errorf("stat cert %s: %w", p.certFile, err)
	}
	keyInfo, err := os.Stat(p.keyFile)
	if err != nil {
		return fmt.Errorf("stat key %s: %w", p.keyFile, err)
	}

	p.mu.RLock()
	unchanged := p.cert != nil &&
		certInfo.ModTime().Equal(p.certMod) && certInfo.Size() == p.certSize &&
		keyInfo.ModTime().Equal(p.keyMod) && keyInfo.Size() == p.keySize
	p.mu.RUnlock()
	if unchanged {
		return nil
	}

	cert, err := tls.LoadX509KeyPair(p.certFile, p.keyFile)
	if err != nil {
		return fmt.Errorf("load keypair (%s, %s): %w", p.certFile, p.keyFile, err)
	}

	p.mu.Lock()
	p.cert = &cert
	p.certMod, p.certSize = certInfo.ModTime(), certInfo.Size()
	p.keyMod, p.keySize = keyInfo.ModTime(), keyInfo.Size()
	p.mu.Unlock()
	return nil
}

func (p *Provider) refreshLoop(refresh time.Duration) {
	if refresh <= 0 {
		refresh = DefaultRefreshInterval
	}
	t := time.NewTicker(refresh)
	defer t.Stop()
	for {
		select {
		case <-p.stop:
			return
		case <-t.C:
			// Best-effort: a transient read error (e.g. mid-rotation) keeps the
			// last good cert until the next tick rather than dropping it.
			_ = p.reload()
		}
	}
}

// current returns the cached parsed certificate.
func (p *Provider) current() (*tls.Certificate, error) {
	p.mu.RLock()
	cert := p.cert
	p.mu.RUnlock()
	if cert == nil {
		return nil, fmt.Errorf("no TLS key material available for %s", p.certFile)
	}
	return cert, nil
}

// Close stops the background refresh goroutine. Idempotent and nil-safe.
func (p *Provider) Close() {
	if p == nil {
		return
	}
	p.stopOnce.Do(func() { close(p.stop) })
}

// NewServerGetCertificate returns a callback suitable for
// tls.Config.GetCertificate, backed by a reloading keypair so rotated certs
// (e.g. from k8s cert-manager) are picked up without a restart. Caller should
// Close() the returned Provider at shutdown.
func NewServerGetCertificate(certFile, keyFile string) (func(*tls.ClientHelloInfo) (*tls.Certificate, error), *Provider, error) {
	return newServerGetCertificate(certFile, keyFile, DefaultRefreshInterval)
}

func newServerGetCertificate(certFile, keyFile string, refresh time.Duration) (func(*tls.ClientHelloInfo) (*tls.Certificate, error), *Provider, error) {
	provider, err := newProvider(certFile, keyFile, refresh)
	if err != nil {
		return nil, nil, err
	}
	get := func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
		return provider.current()
	}
	return get, provider, nil
}

// NewClientGetCertificate returns a callback suitable for
// tls.Config.GetClientCertificate. Fires per TLS handshake, so long-lived
// HTTPS clients (FUSE mount, backup, filer→volume, etc.) pick up rotated
// client mTLS certs as pooled connections recycle.
func NewClientGetCertificate(certFile, keyFile string) (func(*tls.CertificateRequestInfo) (*tls.Certificate, error), *Provider, error) {
	return newClientGetCertificate(certFile, keyFile, DefaultRefreshInterval)
}

func newClientGetCertificate(certFile, keyFile string, refresh time.Duration) (func(*tls.CertificateRequestInfo) (*tls.Certificate, error), *Provider, error) {
	provider, err := newProvider(certFile, keyFile, refresh)
	if err != nil {
		return nil, nil, err
	}
	get := func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
		return provider.current()
	}
	return get, provider, nil
}
