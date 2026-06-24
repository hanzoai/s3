package mount

import (
	"context"
	"io"
	"testing"

	"github.com/hanzoai/s3/s3/util/chunk_cache"
	mount_peerwire "github.com/hanzoai/s3/s3/wire/mount_peer"
)

// newTestPeerGrpc starts a real PeerGrpcServer on a loopback port and returns a
// ZAP client connected to it plus a cleanup. The client speaks the wire
// protocol over the transport, exactly as production does.
func newTestPeerGrpc(t *testing.T, cache chunk_cache.ChunkCache, ownerFor func(fid string) string, selfAddr string) (MountPeerClient, func()) {
	t.Helper()
	dir := NewPeerDirectory()
	srv := NewPeerGrpcServer(cache, dir, ownerFor, selfAddr)
	if err := srv.Start("127.0.0.1:0"); err != nil {
		t.Fatalf("start: %v", err)
	}

	dial := DefaultMountPeerDialer()
	client, closeFn, err := dial(context.Background(), srv.Addr())
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	cleanup := func() {
		closeFn()
		srv.Stop()
	}
	return client, cleanup
}

func TestPeerGrpcServer_AnnounceAndLookup(t *testing.T) {
	self := "self:18080"
	client, cleanup := newTestPeerGrpc(t, nil, func(fid string) string {
		return self
	}, self)
	defer cleanup()

	_, err := client.ChunkAnnounce(mount_peerwire.ChunkAnnounceRequestInput{
		FileIds:    [][]byte{[]byte("3,a"), []byte("3,b")},
		PeerAddr:   "holder:18080",
		TtlSeconds: 60,
	})
	if err != nil {
		t.Fatalf("announce: %v", err)
	}

	resp, err := client.ChunkLookup(mount_peerwire.ChunkLookupRequestInput{
		FileIds: [][]byte{[]byte("3,a"), []byte("3,missing")},
	})
	if err != nil {
		t.Fatalf("lookup: %v", err)
	}

	peersFor := func(fid string) (mount_peerwire.PeerSet, bool) {
		entries := resp.PeersByFid()
		for i := 0; i < entries.Length(); i++ {
			entry, werr := mount_peerwire.WrapChunkLookupResponsePeersByFidEntry(entries.BytesAt(i))
			if werr == nil && entry.Key() == fid {
				return entry.Value(), true
			}
		}
		return mount_peerwire.PeerSet{}, false
	}

	if set, ok := peersFor("3,a"); !ok || set.Peers().Length() != 1 {
		t.Errorf("3,a: expected 1 holder, got ok=%v len=%d", ok, set.Peers().Length())
	}
	if _, ok := peersFor("3,missing"); !ok {
		t.Errorf("3,missing should have a (nil/empty) peer set entry")
	}
}

func TestPeerGrpcServer_OwnerMismatch(t *testing.T) {
	client, cleanup := newTestPeerGrpc(t, nil, func(fid string) string {
		return "some-other:18080"
	}, "self:18080")
	defer cleanup()

	ann, err := client.ChunkAnnounce(mount_peerwire.ChunkAnnounceRequestInput{
		FileIds:    [][]byte{[]byte("3,x")},
		PeerAddr:   "holder:18080",
		TtlSeconds: 60,
	})
	if err != nil {
		t.Fatalf("announce: %v", err)
	}
	rejected := ann.RejectedFileIds()
	if rejected.Length() != 1 || string(rejected.BytesAt(0)) != "3,x" {
		t.Errorf("expected 3,x rejected, got len=%d", rejected.Length())
	}
}

// TestPeerGrpcServer_FetchChunk_StreamsHit exercises the byte-transfer path: a
// cached chunk returned as multiple ZAP stream frames, concatenated by the
// caller into the original payload.
func TestPeerGrpcServer_FetchChunk_StreamsHit(t *testing.T) {
	cache := newFakeChunkCache()
	// A payload deliberately larger than fetchChunkStreamSize to force
	// multiple Send() calls.
	payload := make([]byte, fetchChunkStreamSize*2+123)
	for i := range payload {
		payload[i] = byte(i & 0xff)
	}
	cache.Put("3,stream", payload)

	client, cleanup := newTestPeerGrpc(t, cache, nil, "self:18080")
	defer cleanup()

	stream, err := client.FetchChunk(mount_peerwire.FetchChunkRequestInput{FileId: "3,stream"})
	if err != nil {
		t.Fatalf("FetchChunk: %v", err)
	}
	var got []byte
	for {
		frame, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Recv: %v", err)
		}
		resp, werr := mount_peerwire.WrapFetchChunkResponse(frame)
		if werr != nil {
			t.Fatalf("wrap frame: %v", werr)
		}
		got = append(got, resp.Data()...)
	}
	if len(got) != len(payload) {
		t.Fatalf("got %d bytes, want %d", len(got), len(payload))
	}
	for i, b := range got {
		if b != payload[i] {
			t.Fatalf("byte mismatch at offset %d: got %02x want %02x", i, b, payload[i])
		}
	}
}

// TestPeerGrpcServer_FetchChunk_NotFound verifies that a miss yields a
// zero-frame stream (the server half-closes without sending any
// FetchChunkResponse), which the client reads as io.EOF with no bytes.
func TestPeerGrpcServer_FetchChunk_NotFound(t *testing.T) {
	cache := newFakeChunkCache()
	client, cleanup := newTestPeerGrpc(t, cache, nil, "self:18080")
	defer cleanup()

	stream, err := client.FetchChunk(mount_peerwire.FetchChunkRequestInput{FileId: "3,missing"})
	if err != nil {
		t.Fatalf("FetchChunk open: %v", err)
	}
	frame, err := stream.Recv()
	if err != io.EOF {
		t.Fatalf("expected io.EOF on missing fid, got frame=%d err=%v", len(frame), err)
	}
}
