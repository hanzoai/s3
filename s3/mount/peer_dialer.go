package mount

import (
	"context"
	"sync"
	"time"

	mount_peerwire "github.com/hanzoai/s3/s3/wire/mount_peer"
	"github.com/zap-proto/go/transport"
)

// mountPeerConn is the per-peer ZAP client the announcer and fetcher use. It
// pairs the typed wire client (for the unary directory RPCs) with the raw
// transport.Conn (for the FetchChunk server-stream, opened via the transport's
// stream primitive on the MountPeerFetchChunkOrdinal).
type mountPeerConn struct {
	conn   transport.Conn
	client *mount_peerwire.MountPeerClient
}

func newMountPeerConn(conn transport.Conn) *mountPeerConn {
	return &mountPeerConn{conn: conn, client: mount_peerwire.NewMountPeerClient(conn, nil)}
}

// ChunkAnnounce issues the unary ChunkAnnounce RPC and decodes the response.
func (c *mountPeerConn) ChunkAnnounce(in mount_peerwire.ChunkAnnounceRequestInput) (mount_peerwire.ChunkAnnounceResponse, error) {
	_, body, err := c.client.ChunkAnnounce(mount_peerwire.NewChunkAnnounceRequest(in))
	if err != nil {
		return mount_peerwire.ChunkAnnounceResponse{}, err
	}
	return mount_peerwire.WrapChunkAnnounceResponse(body)
}

// ChunkLookup issues the unary ChunkLookup RPC and decodes the response.
func (c *mountPeerConn) ChunkLookup(in mount_peerwire.ChunkLookupRequestInput) (mount_peerwire.ChunkLookupResponse, error) {
	_, body, err := c.client.ChunkLookup(mount_peerwire.NewChunkLookupRequest(in))
	if err != nil {
		return mount_peerwire.ChunkLookupResponse{}, err
	}
	return mount_peerwire.WrapChunkLookupResponse(body)
}

// FetchChunk opens the FetchChunk server-stream on the connection and returns
// the stream to iterate (Recv yields FetchChunkResponse frames until io.EOF).
func (c *mountPeerConn) FetchChunk(in mount_peerwire.FetchChunkRequestInput) (transport.Stream, error) {
	return c.conn.OpenStream(mount_peerwire.MountPeerFetchChunkOrdinal, mount_peerwire.NewFetchChunkRequest(in))
}

// MountPeerClient is the per-peer client contract the announcer and fetcher
// consume. Tests inject a fake; production uses *mountPeerConn over a real ZAP
// connection. ChunkAnnounce / ChunkLookup are unary; FetchChunk is a server
// stream of FetchChunkResponse frames.
type MountPeerClient interface {
	ChunkAnnounce(in mount_peerwire.ChunkAnnounceRequestInput) (mount_peerwire.ChunkAnnounceResponse, error)
	ChunkLookup(in mount_peerwire.ChunkLookupRequestInput) (mount_peerwire.ChunkLookupResponse, error)
	FetchChunk(in mount_peerwire.FetchChunkRequestInput) (transport.Stream, error)
}

// MountPeerDialer opens a MountPeerClient to a given peer over the ZAP
// transport. Tests inject a fake; production uses a real dial backed by a short
// connection cache. The returned closeFn releases the client (a no-op for the
// pooled dialer, which owns the connection lifecycle).
type MountPeerDialer func(ctx context.Context, peerAddr string) (MountPeerClient, func(), error)

// PeerConnPool caches one long-lived ZAP connection per peer address. Both the
// announcer flush loop and the read-path fetcher hit the same handful of
// directory owners repeatedly; without a cache each call would pay TCP
// handshake cost. The cache makes steady-state owner RPCs effectively free
// after the first call.
//
// Sizing: entries are ~1 KB + the conn itself; bounded at
// maxPeerConnPoolEntries to contain runaway growth.
type PeerConnPool struct {
	mu    sync.Mutex
	conns map[string]*mountPeerConn
}

// maxPeerConnPoolEntries caps live peer conns per mount. A 10k-mount
// fleet with HRW-sharded directory reaches only ~200 distinct owner
// addresses per mount in the worst case, so this is far above any real
// footprint while still bounding pathological growth.
const maxPeerConnPoolEntries = 4096

// NewPeerConnPool returns an empty pool.
func NewPeerConnPool() *PeerConnPool {
	return &PeerConnPool{
		conns: map[string]*mountPeerConn{},
	}
}

// Dialer returns a MountPeerDialer bound to this pool. The returned
// closeFn is a no-op — the pool owns the connection lifecycle. Tests
// that want per-call dials can keep using the non-pooled variant.
func (p *PeerConnPool) Dialer() MountPeerDialer {
	return func(ctx context.Context, peerAddr string) (MountPeerClient, func(), error) {
		c, err := p.get(peerAddr)
		if err != nil {
			return nil, func() {}, err
		}
		return c, func() {}, nil
	}
}

func (p *PeerConnPool) get(peerAddr string) (*mountPeerConn, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if c, ok := p.conns[peerAddr]; ok {
		return c, nil
	}
	if len(p.conns) >= maxPeerConnPoolEntries {
		// Evict one arbitrary entry. Simple over LRU: the pool is small
		// in practice, and the victim will be re-dialed if needed.
		for k, c := range p.conns {
			_ = c.conn.Close()
			delete(p.conns, k)
			break
		}
	}
	conn, err := transport.Dial("tcp", peerAddr)
	if err != nil {
		return nil, err
	}
	c := newMountPeerConn(conn)
	p.conns[peerAddr] = c
	return c, nil
}

// Close tears down every cached connection. Safe to call multiple times.
func (p *PeerConnPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for addr, c := range p.conns {
		_ = c.conn.Close()
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
	return func(ctx context.Context, peerAddr string) (MountPeerClient, func(), error) {
		conn, err := transport.Dial("tcp", peerAddr)
		if err != nil {
			return nil, func() {}, err
		}
		return newMountPeerConn(conn), func() { _ = conn.Close() }, nil
	}
}

// peerConnMaxAge lets external callers (or future metrics) decide when a
// pool entry is "stale" for monitoring purposes. The pool itself does not
// expire on age — the transport handles reconnect on next dial.
var peerConnMaxAge = 10 * time.Minute

var _ = peerConnMaxAge // suppress unused-lint until a metric consumes it
