// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package filerzap

import (
	"context"
	"io"
	"testing"
	"time"

	filer_pb "github.com/hanzoai/s3/s3/pb/filer_pb"
	"github.com/hanzoai/s3/s3/pb/filerstub"
	filerwire "github.com/hanzoai/s3/s3/wire/filer"
	"github.com/hanzoai/s3/s3/wire/filer/filerstream"
	"github.com/zap-proto/go/transport"
)

// subFakeFiler streams a fixed set of metadata events for SubscribeMetadata.
type subFakeFiler struct {
	filerstub.FilerServer
	dirs []string
}

func (f *subFakeFiler) SubscribeMetadata(req *filer_pb.SubscribeMetadataRequest, stream filer_pb.HanzoFiler_SubscribeMetadataServer) error {
	for i, d := range f.dirs {
		if err := stream.Send(&filer_pb.SubscribeMetadataResponse{
			Directory: d,
			TsNs:      int64(i + 1),
			EventNotification: &filer_pb.EventNotification{
				NewEntry: &filer_pb.Entry{Name: "evt", Content: []byte(req.ClientName)},
			},
		}); err != nil {
			return err
		}
	}
	return nil
}

// TestStreamServerSubscribeMetadata proves the server-streaming path works over
// ZAP end-to-end: a real fs.SubscribeMetadata streams events through the ZAP
// stream server and the ZAP client receives each one with field fidelity.
func TestStreamServerSubscribeMetadata(t *testing.T) {
	fs := &subFakeFiler{dirs: []string{"/a", "/b", "/c"}}

	srv, err := transport.ListenStream("tcp", "127.0.0.1:0",
		filerwire.Dispatch(NewServerBackend(fs)),
		filerstream.Handler(NewStreamServer(fs)))
	if err != nil {
		t.Fatalf("ListenStream: %v", err)
	}
	defer srv.Close()

	conn, err := transport.Dial("tcp", srv.Addr().String())
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer conn.Close()
	client := NewZapFilerClient(conn, nil)

	stream, err := client.SubscribeMetadata(context.Background(), &filer_pb.SubscribeMetadataRequest{ClientName: "mount-1"})
	if err != nil {
		t.Fatalf("SubscribeMetadata: %v", err)
	}

	var got []string
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Recv: %v", err)
		}
		if resp.EventNotification == nil || resp.EventNotification.NewEntry == nil ||
			string(resp.EventNotification.NewEntry.Content) != "mount-1" {
			t.Fatalf("event fidelity lost: %+v", resp)
		}
		got = append(got, resp.Directory)
	}
	if len(got) != 3 || got[0] != "/a" || got[2] != "/c" {
		t.Fatalf("streamed dirs = %v, want [/a /b /c]", got)
	}
}

// subLeakFiler blocks SubscribeMetadata on the stream context — modelling the
// real engine, which gates its idle wait on stream.Context().Done().
type subLeakFiler struct {
	filerstub.FilerServer
	released chan struct{}
}

func (f *subLeakFiler) SubscribeMetadata(_ *filer_pb.SubscribeMetadataRequest, stream filer_pb.HanzoFiler_SubscribeMetadataServer) error {
	<-stream.Context().Done() // idle: only cancellation frees us
	close(f.released)
	return stream.Context().Err()
}

// TestStreamServerSubscribeReleaseOnDisconnect proves the leak fix end to end:
// an idle SubscribeMetadata handler is released when the ZAP client drops the
// connection, because the per-stream Context (zap-proto/go v1.6.2) cancels.
func TestStreamServerSubscribeReleaseOnDisconnect(t *testing.T) {
	fs := &subLeakFiler{released: make(chan struct{})}
	srv, err := transport.ListenStream("tcp", "127.0.0.1:0",
		filerwire.Dispatch(NewServerBackend(fs)), filerstream.Handler(NewStreamServer(fs)))
	if err != nil {
		t.Fatalf("ListenStream: %v", err)
	}
	defer srv.Close()

	conn, err := transport.Dial("tcp", srv.Addr().String())
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	if _, err := NewZapFilerClient(conn, nil).SubscribeMetadata(
		context.Background(), &filer_pb.SubscribeMetadataRequest{ClientName: "leaky", ClientId: 1, ClientEpoch: 1}); err != nil {
		t.Fatalf("SubscribeMetadata: %v", err)
	}
	_ = conn.Close() // drop the (idle) subscription

	select {
	case <-fs.released:
	case <-time.After(3 * time.Second):
		t.Fatal("idle SubscribeMetadata handler not released after client disconnect — leak")
	}
}
