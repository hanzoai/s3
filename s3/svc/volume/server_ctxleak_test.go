// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// server_ctxleak_test.go — RED adversarial probe of vector #5: the streaming
// server-stream Send adapters in stream_server.go hand the engine
// context.Background() (serverBackend.ctx), NOT a per-stream context cancelled
// on conn-drop. The filer template (s3/wire/filer/filerstream) wires
// Context() -> transport.Stream.Context() on every server-stream type; the
// volume wire (s3/wire/volume_server) does not, so a context-respecting
// VolumeTailSender handler that gates its idle wait on stream.Context().Done()
// LEAKS its goroutine when the client disconnects idle.
//
// This is exactly the live-follower pattern: a volume tail-sender that, while
// caught up (isLastOne), waits for either new data OR cancellation. The real
// s3/server/volume_grpc_tail.go loop relies on Send() erroring; but ANY handler
// (this one, or a future ctx-aware rewrite, or the auth checkGrpcAdminAuth path
// that already calls stream.Context()) that blocks on Context().Done() will
// never wake under ZAP because the context is Background().

package volume

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/zap-proto/go/transport"

	volume_server_pb "github.com/hanzoai/s3/s3/pb/volume_server_pb"
	vsw "github.com/hanzoai/s3/s3/wire/volume_server"
)

// ctxGatedVolumeServer's VolumeTailSender blocks until the stream context is
// cancelled — the canonical "idle follower waits for cancel" shape. It closes
// released when the handler returns, so the test can observe whether the
// handler goroutine is freed when the client conn drops.
type ctxGatedVolumeServer struct {
	volume_server_pb.UnimplementedVolumeServerServer
	entered  chan struct{}
	released chan struct{}
	once     sync.Once
}

func (f *ctxGatedVolumeServer) VolumeTailSender(req *volume_server_pb.VolumeTailSenderRequest, stream volume_server_pb.VolumeServer_VolumeTailSenderServer) error {
	f.once.Do(func() { close(f.entered) })
	// A live follower that is caught up waits for new data or cancellation.
	// Here we model the cancellation arm only: when the client disconnects,
	// the stream context MUST fire. Under ZAP-volume it is context.Background().
	<-stream.Context().Done()
	close(f.released)
	return stream.Context().Err()
}

// TestVolumeTailSender_ContextCancelOnConnDrop opens a tail stream, waits for
// the handler to block on Context().Done(), drops the client connection, and
// asserts the handler releases. It will FAIL on the current volume wire because
// the Send adapter's Context() returns context.Background() (never cancelled).
//
// When wire/volume_server.serverStream exposes Context() (mirroring
// filerstream's `func (x *XStream) Context() context.Context { return x.s.Context() }`)
// and the *ServerStream types + the volume Send adapters thread it through,
// this test passes.
func TestVolumeTailSender_ContextCancelOnConnDrop(t *testing.T) {
	fake := &ctxGatedVolumeServer{
		entered:  make(chan struct{}),
		released: make(chan struct{}),
	}
	store := NewServerBackend(fake)
	srv, err := transport.ListenStream("tcp", "127.0.0.1:0", vsw.Dispatch(store), vsw.StreamHandler(store))
	if err != nil {
		t.Fatalf("ListenStream: %v", err)
	}
	defer srv.Close()

	conn, err := transport.Dial("tcp", srv.Addr().String())
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	cli := New(conn, nil)

	// Open the server-stream. The request rides as the stream init frame.
	stream, err := cli.VolumeTailSender(context.Background(), &volume_server_pb.VolumeTailSenderRequest{
		VolumeId: 1, SinceNs: 0, IdleTimeoutSeconds: 0,
	})
	if err != nil {
		t.Fatalf("VolumeTailSender open: %v", err)
	}
	_ = stream

	// Wait for the handler to be blocked on Context().Done().
	select {
	case <-fake.entered:
	case <-time.After(2 * time.Second):
		t.Fatal("handler never entered")
	}

	// Client disconnects idle. The server-side stream context MUST fire so the
	// handler goroutine is released.
	conn.Close()

	select {
	case <-fake.released:
		// handler released — no leak
	case <-time.After(3 * time.Second):
		t.Fatal("LEAK: VolumeTailSender handler goroutine did not release after client conn drop; " +
			"stream Context() is context.Background() (stream_server.go a.ctx=b.ctx=Background), " +
			"not transport.Stream.Context(). Wire/volume_server.serverStream must expose Context().")
	}
}
