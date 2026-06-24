package mount

import (
	"fmt"
	"time"

	"github.com/hanzoai/s3/s3/glog"
	"github.com/hanzoai/s3/s3/util/chunk_cache"
	"github.com/hanzoai/s3/s3/util/mem"
	mount_peerwire "github.com/hanzoai/s3/s3/wire/mount_peer"
	zap "github.com/zap-proto/go"
	"github.com/zap-proto/go/transport"
)

// fetchChunkStreamSize is the frame size used when server-streaming a chunk's
// bytes back to a peer. 1 MiB keeps each ZAP stream frame small so every Recv
// on the client returns quickly and the chunk is assembled with ~16 Recv calls
// for typical 16 MiB chunks.
const fetchChunkStreamSize = 1 * 1024 * 1024

// maxFetchChunkBytes caps the size of a single FetchChunk read buffer.
// The caller's expected_size is untrusted — a misbehaving peer could
// claim 1 TiB and OOM the server. 64 MiB is well above any realistic
// chunk (Hanzo defaults to 16 MiB, filer manifests cap at 32 MiB);
// anything larger is treated as invalid input.
const maxFetchChunkBytes = 64 * 1024 * 1024

// PeerGrpcServer is the single mount-to-mount ZAP endpoint. It serves:
//   - ChunkAnnounce / ChunkLookup — the tier-2 directory RPCs, populated
//     by inbound announces and queried by inbound lookups. Each handler
//     is HRW-gated on the caller-side seed view.
//   - FetchChunk (server stream) — serves bytes from the local
//     chunk_cache to peers. Replaces the earlier HTTP-only peer-serve
//     endpoint: one port, one authentication path, one connection pool.
//
// The unary methods are dispatched via mount_peerwire.DispatchMountPeer over
// the ZAP transport; FetchChunk's byte stream is served through the
// transport's stream handler keyed on the MountPeerFetchChunkOrdinal.
type PeerGrpcServer struct {
	dir      *PeerDirectory
	cache    chunk_cache.ChunkCache
	ownerFor func(fid string) string // HRW owner predicate on current seeds
	selfAddr string
	srv      *transport.Server
	addr     string
	stopped  bool
}

// NewPeerGrpcServer constructs the server. cache is the local chunk_cache
// (used to serve FetchChunk); dir is the local directory shard (used to
// answer ChunkAnnounce / ChunkLookup); ownerFor returns the HRW owner of a
// fid on the current seed view.
func NewPeerGrpcServer(cache chunk_cache.ChunkCache, dir *PeerDirectory, ownerFor func(fid string) string, selfAddr string) *PeerGrpcServer {
	return &PeerGrpcServer{
		cache:    cache,
		dir:      dir,
		ownerFor: ownerFor,
		selfAddr: selfAddr,
	}
}

// Start binds a TCP listener at addr and serves the MountPeer service over the
// ZAP transport (unary dispatch + FetchChunk stream handler).
func (s *PeerGrpcServer) Start(addr string) error {
	srv, err := transport.ListenStream("tcp", addr,
		func(env []byte) ([]byte, error) {
			return mount_peerwire.DispatchMountPeer(s, env)
		},
		s.streamHandler,
	)
	if err != nil {
		return fmt.Errorf("peer listen %s: %w", addr, err)
	}
	s.srv = srv
	s.addr = srv.Addr().String()
	glog.V(0).Infof("peer-zap listening on %s", s.addr)
	return nil
}

// streamHandler serves the FetchChunk server-stream: it decodes the init
// FetchChunkRequest, streams the cached chunk's bytes back as
// FetchChunkResponse frames, then returns (which half-closes the stream so the
// client's Recv sees io.EOF). A miss or invalid request returns without
// sending any frame; the client treats a frameless stream as a miss.
func (s *PeerGrpcServer) streamHandler(method uint32, init []byte, stream *transport.Stream) {
	if method != mount_peerwire.MountPeerFetchChunkOrdinal {
		return
	}
	if err := s.serveFetchChunk(init, stream); err != nil {
		glog.V(2).Infof("peer-fetch stream: %v", err)
	}
}

// Stop halts the ZAP server without waiting for in-flight streams.
func (s *PeerGrpcServer) Stop() {
	if s.stopped {
		return
	}
	s.stopped = true
	if s.srv != nil {
		_ = s.srv.Close()
	}
}

// Addr returns the bound address (useful when the caller used ":0").
func (s *PeerGrpcServer) Addr() string {
	return s.addr
}

// SelfAddr returns the advertise address this server was constructed with —
// the identity the fetcher compares against to avoid dialing itself.
func (s *PeerGrpcServer) SelfAddr() string { return s.selfAddr }

// ChunkAnnounce accepts holder entries for fids this mount owns; rejects
// others so the caller can retry against the correct owner. It implements the
// mount_peerwire.MountPeerHandler contract over the ZAP transport.
func (s *PeerGrpcServer) ChunkAnnounce(reqBytes []byte) ([]byte, error) {
	req, err := mount_peerwire.WrapChunkAnnounceRequest(reqBytes)
	if err != nil {
		return nil, err
	}
	fileIds := stringsFromList(req.FileIds())
	ttl := time.Duration(req.TtlSeconds()) * time.Second
	res := s.dir.Announce(req.PeerAddr(), req.DataCenter(), req.Rack(), fileIds, ttl, s.ownerPredicate)
	return mount_peerwire.NewChunkAnnounceResponse(mount_peerwire.ChunkAnnounceResponseInput{
		RejectedFileIds: bytesList(res.Rejected),
	}), nil
}

// ChunkLookup returns known holders for each requested fid in LRU order. It
// implements the mount_peerwire.MountPeerHandler contract.
func (s *PeerGrpcServer) ChunkLookup(reqBytes []byte) ([]byte, error) {
	req, err := mount_peerwire.WrapChunkLookupRequest(reqBytes)
	if err != nil {
		return nil, err
	}
	res := s.dir.Lookup(stringsFromList(req.FileIds()), s.ownerPredicate)

	entries := make([][]byte, 0, len(res.PeersByFid))
	for fid, holders := range res.PeersByFid {
		peers := make([][]byte, 0, len(holders))
		for _, h := range holders {
			peers = append(peers, mount_peerwire.NewPeerInfo(mount_peerwire.PeerInfoInput{
				PeerAddr:   h.PeerAddr,
				DataCenter: h.DataCenter,
				Rack:       h.Rack,
			}))
		}
		set := mount_peerwire.NewPeerSet(mount_peerwire.PeerSetInput{Peers: peers})
		entries = append(entries, mount_peerwire.NewChunkLookupResponsePeersByFidEntry(
			mount_peerwire.ChunkLookupResponsePeersByFidEntryInput{Key: fid, Value: set}))
	}
	return mount_peerwire.NewChunkLookupResponse(mount_peerwire.ChunkLookupResponseInput{
		PeersByFid:      entries,
		NotOwnerFileIds: bytesList(res.NotOwnerFids),
	}), nil
}

// FetchChunk satisfies the MountPeerHandler interface for a unary call. The
// production path serves FetchChunk as a server stream (see streamHandler /
// serveFetchChunk); a unary caller gets the whole chunk in one response body.
func (s *PeerGrpcServer) FetchChunk(reqBytes []byte) ([]byte, error) {
	req, err := mount_peerwire.WrapFetchChunkRequest(reqBytes)
	if err != nil {
		return nil, err
	}
	buf, lo, hi, err := s.readChunkRange(req)
	if err != nil {
		return nil, err
	}
	defer mem.Free(buf)
	return mount_peerwire.NewFetchChunkResponse(mount_peerwire.FetchChunkResponseInput{
		Data: buf[lo:hi],
	}), nil
}

// serveFetchChunk streams bytes of a cached chunk back to the caller as
// FetchChunkResponse frames of at most fetchChunkStreamSize bytes. A miss or
// invalid request returns an error and no frame; the client reads zero frames
// as a miss.
func (s *PeerGrpcServer) serveFetchChunk(initBytes []byte, stream *transport.Stream) error {
	req, err := mount_peerwire.WrapFetchChunkRequest(initBytes)
	if err != nil {
		return err
	}
	buf, lo, hi, err := s.readChunkRange(req)
	if err != nil {
		return err
	}
	defer mem.Free(buf)

	for off := lo; off < hi; off += fetchChunkStreamSize {
		end := off + fetchChunkStreamSize
		if end > hi {
			end = hi
		}
		frame := mount_peerwire.NewFetchChunkResponse(mount_peerwire.FetchChunkResponseInput{
			Data: buf[off:end],
		})
		if sendErr := stream.Send(frame); sendErr != nil {
			return sendErr
		}
	}
	return nil
}

// readChunkRange validates a FetchChunkRequest, reads the cached chunk into a
// pooled buffer, and returns it with the [lo, hi) byte range to transfer. The
// caller owns the buffer and must mem.Free it.
func (s *PeerGrpcServer) readChunkRange(req mount_peerwire.FetchChunkRequest) (buf []byte, lo, hi int, err error) {
	if s.cache == nil {
		return nil, 0, 0, fmt.Errorf("chunk cache not configured")
	}
	fid := req.FileId()
	if fid == "" {
		return nil, 0, 0, fmt.Errorf("missing file_id")
	}

	// expected_size / length are untrusted — cap to prevent a misbehaving
	// peer from requesting an OOM-sized allocation.
	if req.ExpectedSize() > maxFetchChunkBytes {
		return nil, 0, 0, fmt.Errorf("expected_size %d exceeds max %d", req.ExpectedSize(), maxFetchChunkBytes)
	}
	if req.Length() > maxFetchChunkBytes {
		return nil, 0, 0, fmt.Errorf("length %d exceeds max %d", req.Length(), maxFetchChunkBytes)
	}

	// Size the read buffer to the caller-reported chunk length. The
	// TieredChunkCache.ReadChunkAt wrapper only returns success when n
	// equals len(data), so a buffer larger than the actual stored chunk
	// makes every read look like a miss. Fall back to the max-part-size
	// when the caller left it zero.
	readSize := int(req.ExpectedSize())
	if readSize <= 0 {
		max := s.cache.GetMaxFilePartSizeInCache()
		if max == 0 {
			max = 8 * 1024 * 1024
		}
		readSize = int(max)
	}
	// mem.Allocate rounds up to the nearest power-of-2 slot backed by a
	// shared sync.Pool; avoids an allocation per FetchChunk call.
	buf = mem.Allocate(readSize)

	n, rerr := s.cache.ReadChunkAt(buf, fid, 0)
	if rerr != nil || n <= 0 {
		mem.Free(buf)
		return nil, 0, 0, fmt.Errorf("fid %s not cached", fid)
	}

	// Apply optional offset / length range to the whole-chunk buffer.
	// Default (both 0) is a full-chunk transfer starting at byte 0.
	lo = int(req.Offset())
	if lo < 0 || lo > n {
		mem.Free(buf)
		return nil, 0, 0, fmt.Errorf("offset %d outside chunk length %d", lo, n)
	}
	hi = n
	if req.Length() > 0 {
		hi = lo + int(req.Length())
		if hi > n {
			hi = n
		}
	}
	return buf, lo, hi, nil
}

// stringsFromList materialises a ZAP repeated-string field into a []string.
func stringsFromList(l zap.List) []string {
	out := make([]string, l.Length())
	for i := range out {
		out[i] = string(l.BytesAt(i))
	}
	return out
}

// bytesList encodes each string as a raw-bytes list entry for a ZAP repeated
// string field.
func bytesList(ss []string) [][]byte {
	out := make([][]byte, len(ss))
	for i, s := range ss {
		out[i] = []byte(s)
	}
	return out
}

func (s *PeerGrpcServer) ownerPredicate(fid string) bool {
	if s.ownerFor == nil {
		return true // no HRW configured → accept all (single-mount mode)
	}
	return s.ownerFor(fid) == s.selfAddr
}

// peerDirectorySweepInterval is how often the mount evicts expired
// directory entries. Lookup no longer deletes inline (it takes only an
// RLock), so this sweeper is the sole memory reclamation path.
const peerDirectorySweepInterval = 60 * time.Second

// runPeerDirectorySweeper runs until stopCh closes. Stopping it on
// unmount prevents the goroutine from outliving the WFS, which matters
// for tests and for any embedded use of NewHanzoFileSystem where the
// filesystem object is recreated without a process exit.
func (wfs *WFS) runPeerDirectorySweeper(stopCh <-chan struct{}) {
	ticker := time.NewTicker(peerDirectorySweepInterval)
	defer ticker.Stop()
	for {
		select {
		case <-stopCh:
			return
		case <-ticker.C:
			dir := wfs.peerDirectory
			if dir == nil {
				return
			}
			if evicted := dir.Sweep(); evicted > 0 {
				glog.V(2).Infof("peer directory: evicted %d expired entries", evicted)
			}
		}
	}
}
