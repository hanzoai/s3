package mount

import (
	"context"
	"sync"
	"time"

	"github.com/zap-proto/go/transport"
)

// PeerConnPool caches one long-lived ZAP transport connection per peer
// address. Both the announcer flush loop and the read-path fetcher hit
// the same handful of directory owners repeatedly; without a cache each
// call would pay TCP handshake cost (and on TLS, an additional
// handshake). The cache makes steady-state owner RPCs effectively free
// after the first call.
//
// Connections whose read loop has torn down (peer EOF / local close) are
// transparently replaced on next access. Sizing: entries are ~1 KB + the
// conn itself; bounded at maxPeerConnPoolEntries to contain runaway
// growth.
type PeerConnPool struct {
	mu    sync.Mutex
	conns map[string]*transport.Conn
}

// maxPeerConnPoolEntries caps live peer conns per mount. A 10k-mount
// fleet with HRW-sharded directory reaches only ~200 distinct owner
// addresses per mount in the worst case, so this is far above any real
// footprint while still bounding pathological growth.
const maxPeerConnPoolEntries = 4096

// NewPeerConnPool returns an empty pool. Peer RPCs travel over the native
// ZAP transport (plaintext TCP); the pool only manages connection reuse.
func NewPeerConnPool() *PeerConnPool {
	return &PeerConnPool{
		conns: map[string]*transport.Conn{},
	}
}

// Dialer returns a MountPeerDialer bound to this pool. The returned
// closeFn is a no-op — the pool owns the connection lifecycle. Tests
// that want per-call dials can keep using the non-pooled variant.
func (p *PeerConnPool) Dialer() MountPeerDialer {
	return func(ctx context.Context, peerAddr string) (*transport.Conn, func(), error) {
		conn, err := p.get(peerAddr)
		if err != nil {
			return nil, func() {}, err
		}
		return conn, func() {}, nil
	}
}

func (p *PeerConnPool) get(peerAddr string) (*transport.Conn, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if conn, ok := p.conns[peerAddr]; ok {
		if conn.IsClosed() {
			// Pooled conn's read loop has torn down — drop it and redial.
			_ = conn.Close()
			delete(p.conns, peerAddr)
		} else {
			return conn, nil
		}
	}
	if len(p.conns) >= maxPeerConnPoolEntries {
		// Evict one arbitrary entry. Simple over LRU: the pool is small
		// in practice, and the victim will be re-dialed if needed.
		for k, c := range p.conns {
			_ = c.Close()
			delete(p.conns, k)
			break
		}
	}
	conn, err := transport.Dial("tcp", peerAddr)
	if err != nil {
		return nil, err
	}
	p.conns[peerAddr] = conn
	return conn, nil
}

// Close tears down every cached connection. Safe to call multiple times.
func (p *PeerConnPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for addr, c := range p.conns {
		_ = c.Close()
		delete(p.conns, addr)
	}
}

// Size returns the current number of cached connections — useful for
// tests and metrics exports.
func (p *PeerConnPool) Size() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.conns)
}

// DefaultMountPeerDialer returns a per-call dialer (no pooling). Kept for
// tests and for any caller that genuinely wants a fresh connection per
// invocation. Production code should prefer PeerConnPool.Dialer().
func DefaultMountPeerDialer() MountPeerDialer {
	return func(ctx context.Context, peerAddr string) (*transport.Conn, func(), error) {
		conn, err := transport.Dial("tcp", peerAddr)
		if err != nil {
			return nil, func() {}, err
		}
		return conn, func() { _ = conn.Close() }, nil
	}
}

// peerConnMaxAge lets external callers (or future metrics) decide when a
// pool entry is "stale" for monitoring purposes. The pool itself does not
// expire on age — the transport tears a conn down on peer EOF and the
// pool redials on next access.
var peerConnMaxAge = 10 * time.Minute

var _ = peerConnMaxAge // suppress unused-lint until a metric consumes it
