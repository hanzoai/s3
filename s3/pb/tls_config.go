package pb

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strings"

	"github.com/hanzoai/s3/s3/glog"
	"github.com/hanzoai/s3/s3/util"
)

// DialOption is the grpc-free client dial option threaded through the s3 RPC
// stack. The real dialing is the ZAP transport (transport.Dial/DialTLS with the
// PQ X-Wing curve), so the only thing a dial needs from configuration is the
// optional client *tls.Config — TLS is nil for plaintext (loopback / dev).
// LoadClientTLS builds it; the volume/broker grpc fallback and the ZAP dials
// consume TLS.
type DialOption struct {
	TLS *tls.Config
}

// --- raw *tls.Config for the native ZAP transport (PQ-TLS) ---
// ServerTLSConfig/ClientTLSConfig build a *tls.Config from the
// <component>.cert/.key + grpc.ca + <component>.allowed_commonNames config so
// the ZAP filer/master enforce the same mTLS the legacy gRPC clients did. Wrap
// the result with transport.PQTLSConfig to pin the X25519MLKEM768 PQ X-Wing
// curve. Both return nil when no cert/key is configured (caller serves/dials
// plaintext).

func certPoolFromFile(caFile string) (*x509.CertPool, error) {
	data, err := os.ReadFile(caFile)
	if err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(data) {
		return nil, fmt.Errorf("no certificates parsed from %s", caFile)
	}
	return pool, nil
}

// applyCommonNameVerification gates the peer leaf cert's CommonName against
// <component>.allowed_commonNames / grpc.allowed_wildcard_domain.
func applyCommonNameVerification(cfg *tls.Config, config *util.ViperProxy, component string) {
	allowedCN := config.GetString(component + ".allowed_commonNames")
	allowedWildcard := config.GetString("grpc.allowed_wildcard_domain")
	if allowedCN == "" && allowedWildcard == "" {
		return
	}
	allowed := make(map[string]bool)
	for _, s := range strings.Split(allowedCN, ",") {
		if s != "" {
			allowed[s] = true
		}
	}
	cfg.VerifyPeerCertificate = func(rawCerts [][]byte, _ [][]*x509.Certificate) error {
		if len(rawCerts) == 0 {
			return fmt.Errorf("no peer certificate presented")
		}
		leaf, err := x509.ParseCertificate(rawCerts[0])
		if err != nil {
			return err
		}
		cn := leaf.Subject.CommonName
		if allowedWildcard != "" && strings.HasSuffix(cn, allowedWildcard) {
			return nil
		}
		if allowed[cn] {
			return nil
		}
		return fmt.Errorf("invalid peer common name: %s", cn)
	}
}

// ServerTLSConfig returns a server *tls.Config (mTLS: RequireAndVerifyClientCert
// + CN gate) for the ZAP transport, or nil if <component>.cert/.key is unset.
func ServerTLSConfig(config *util.ViperProxy, component string) *tls.Config {
	if config == nil {
		return nil
	}
	certFile, keyFile := config.GetString(component+".cert"), config.GetString(component+".key")
	if certFile == "" || keyFile == "" {
		return nil
	}
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		glog.Warningf("ServerTLSConfig(%s): load keypair: %v", component, err)
		return nil
	}
	cfg := &tls.Config{Certificates: []tls.Certificate{cert}, MinVersion: tls.VersionTLS13}
	caFile := config.GetString("grpc.ca")
	if caFile == "" {
		// Fail closed: a mesh component with a server cert but no client CA would
		// present a cert and accept ANY client (one-way TLS, no mutual auth) — a
		// silent mTLS downgrade. Internal ZAP services are mTLS-only; refuse to
		// serve rather than silently drop client authentication.
		glog.Fatalf("ServerTLSConfig(%s): %s.cert/.key set but grpc.ca is empty — refusing to serve one-way TLS on the internal mesh; set grpc.ca for mutual PQ-TLS", component, component)
	}
	pool, err := certPoolFromFile(caFile)
	if err != nil {
		glog.Fatalf("ServerTLSConfig(%s): grpc.ca %q unreadable: %v — refusing to serve without client-cert verification", component, caFile, err)
	}
	cfg.ClientCAs = pool
	cfg.ClientAuth = tls.RequireAndVerifyClientCert
	applyCommonNameVerification(cfg, config, component)
	return cfg
}

// ClientTLSConfig returns a client *tls.Config (presents <component> cert, trusts
// grpc.ca) for the ZAP transport, or nil if <component>.cert/.key is unset.
func ClientTLSConfig(config *util.ViperProxy, component string) *tls.Config {
	if config == nil {
		return nil
	}
	certFile, keyFile := config.GetString(component+".cert"), config.GetString(component+".key")
	if certFile == "" || keyFile == "" {
		return nil
	}
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		glog.Warningf("ClientTLSConfig(%s): load keypair: %v", component, err)
		return nil
	}
	cfg := &tls.Config{Certificates: []tls.Certificate{cert}, MinVersion: tls.VersionTLS13}
	if caFile := config.GetString("grpc.ca"); caFile != "" {
		pool, err := certPoolFromFile(caFile)
		if err != nil {
			glog.Warningf("ClientTLSConfig(%s): ca: %v", component, err)
			return nil
		}
		cfg.RootCAs = pool
	}
	return cfg
}
