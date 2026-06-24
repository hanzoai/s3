// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package filerzap

import (
	"context"
	"io"
	"testing"

	filer_pb "github.com/hanzoai/s3/s3/pb/filer_pb"
	filerwire "github.com/hanzoai/s3/s3/wire/filer"
	"github.com/hanzoai/s3/s3/wire/filer/filerstream"
	"github.com/zap-proto/go/transport"
	"google.golang.org/grpc"
)

// subFakeFiler streams a fixed set of metadata events for SubscribeMetadata.
type subFakeFiler struct {
	filer_pb.UnimplementedHanzoFilerServer
	dirs []string
}

func (f *subFakeFiler) SubscribeMetadata(req *filer_pb.SubscribeMetadataRequest, stream grpc.ServerStreamingServer[filer_pb.SubscribeMetadataResponse]) error {
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
