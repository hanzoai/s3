// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// receivefile_oneof_test.go localizes the ReceiveFile client-streaming oneof
// drop independently of the slow erasure-coding e2e. It drives the EXACT path the
// EC shard-distribution caller uses — cli.ReceiveFile(ctx) to open, Send(Info)
// for the first frame, Send(FileContent) for chunks, CloseAndRecv to half-close
// and read the reply — and asserts the server saw the Info oneof variant and the
// content chunks. Before the fix the adapter opened the wire stream with an empty
// ReceiveFileRequestInput{} (Data=0) that the server replayed as its first Recv,
// so the production handler's oneof switch hit default ("unknown message type")
// and the real Info frame never landed.

package volume

import (
	"bytes"
	"context"
	"io"
	"sync"
	"testing"

	volume_server_pb "github.com/hanzoai/s3/s3/pb/volume_server_pb"
)

// strictReceiveFileServer mirrors the PRODUCTION ReceiveFile handler
// (s3/server/volume_grpc_copy.go): its oneof switch has a default that rejects
// any frame whose Data variant is unset — exactly the "unknown message type"
// error the EC e2e hits. The fakeVolumeServer's ReceiveFile lacks that default,
// so it silently swallows the empty opener; this strict server is what proves
// the drop at the adapter boundary.
type strictReceiveFileServer struct {
	volume_server_pb.UnimplementedVolumeServerServer
	mu        sync.Mutex
	gotVolID  uint32
	gotExt    string
	gotChunks [][]byte
}

func (s *strictReceiveFileServer) ReceiveFile(stream volume_server_pb.VolumeServer_ReceiveFileServer) error {
	var total uint64
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		switch d := req.Data.(type) {
		case *volume_server_pb.ReceiveFileRequest_Info:
			s.mu.Lock()
			if d.Info != nil {
				s.gotVolID = d.Info.VolumeId
				s.gotExt = d.Info.Ext
			}
			s.mu.Unlock()
		case *volume_server_pb.ReceiveFileRequest_FileContent:
			s.mu.Lock()
			s.gotChunks = append(s.gotChunks, append([]byte(nil), d.FileContent...))
			s.mu.Unlock()
			total += uint64(len(d.FileContent))
		default:
			return stream.SendAndClose(&volume_server_pb.ReceiveFileResponse{Error: "unknown message type"})
		}
	}
	return stream.SendAndClose(&volume_server_pb.ReceiveFileResponse{BytesWritten: total})
}

func TestReceiveFileOneofOverZAP(t *testing.T) {
	fake := &strictReceiveFileServer{}
	cli, stop := serve(t, fake)
	defer stop()

	stream, err := cli.ReceiveFile(context.Background())
	if err != nil {
		t.Fatalf("ReceiveFile open: %v", err)
	}

	// First frame: the Info oneof variant (proto field 1). This is what the EC
	// shard distributor sends before any content.
	if err := stream.Send(&volume_server_pb.ReceiveFileRequest{
		Data: &volume_server_pb.ReceiveFileRequest_Info{
			Info: &volume_server_pb.ReceiveFileInfo{
				VolumeId:   77,
				Ext:        ".ec00",
				Collection: "col",
				IsEcVolume: true,
				ShardId:    0,
				FileSize:   13,
			},
		},
	}); err != nil {
		t.Fatalf("Send(Info): %v", err)
	}

	chunks := [][]byte{[]byte("hello "), []byte("ec-shard")}
	var want uint64
	for _, ch := range chunks {
		if err := stream.Send(&volume_server_pb.ReceiveFileRequest{
			Data: &volume_server_pb.ReceiveFileRequest_FileContent{FileContent: ch},
		}); err != nil {
			t.Fatalf("Send(FileContent): %v", err)
		}
		want += uint64(len(ch))
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		t.Fatalf("CloseAndRecv: %v", err)
	}
	if resp.Error != "" {
		t.Fatalf("server returned error: %q", resp.Error)
	}
	if resp.BytesWritten != want {
		t.Fatalf("BytesWritten = %d, want %d", resp.BytesWritten, want)
	}

	fake.mu.Lock()
	defer fake.mu.Unlock()
	if fake.gotVolID != 77 || fake.gotExt != ".ec00" {
		t.Fatalf("server lost Info oneof: vol=%d ext=%q (want 77/.ec00)", fake.gotVolID, fake.gotExt)
	}
	if len(fake.gotChunks) != 2 ||
		!bytes.Equal(fake.gotChunks[0], chunks[0]) ||
		!bytes.Equal(fake.gotChunks[1], chunks[1]) {
		t.Fatalf("server saw %d content chunks, want 2 matching", len(fake.gotChunks))
	}
}
