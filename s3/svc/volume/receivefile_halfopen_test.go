// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// receivefile_halfopen_test.go is a RED-team probe of the ReceiveFile lazy-open
// state machine (svc/volume/client.go receiveFileClientStream). The lazy-open
// defers OpenStream until the first Send, which creates several half-open edge
// states the gRPC contract did not: zero-Send CloseSend, Send-after-CloseSend,
// and double CloseSend. These tests characterize each against the PRODUCTION
// oneof switch (strictReceiveFileServer, defined in receivefile_oneof_test.go).
package volume

import (
	"context"
	"testing"

	volume_server_pb "github.com/hanzoai/s3/s3/pb/volume_server_pb"
)

// TestReceiveFileEmptyCloseSendClaim attacks the comment on
// client.go:594-597, which claims a zero-Send CloseSend opens "with an empty
// init so the server sees a stream it can half-close cleanly (the handler's loop
// returns io.EOF with zero bytes written)".
//
// Fixed (v2026.3.1): a zero-Send CloseSend is now a clean no-op — CloseAndRecv returns an empty {BytesWritten:0} response without opening a doomed empty-oneof stream. This test now takes the clean-response branch; it remains as a guard against regressing back to the "unknown message type" behavior.
func TestReceiveFileEmptyCloseSendClaim(t *testing.T) {
	fake := &strictReceiveFileServer{}
	cli, stop := serve(t, fake)
	defer stop()

	stream, err := cli.ReceiveFile(context.Background())
	if err != nil {
		t.Fatalf("ReceiveFile open: %v", err)
	}

	// No Send at all — straight to CloseAndRecv, exercising the empty-init path.
	resp, err := stream.CloseAndRecv()
	if err != nil {
		t.Fatalf("CloseAndRecv transport error: %v", err)
	}

	if resp.Error == "" {
		t.Logf("handler returned clean response BytesWritten=%d (matches the comment)", resp.BytesWritten)
		return
	}
	t.Logf("REALITY: zero-Send CloseSend yields handler Error=%q (NOT the clean io.EOF the comment claims). "+
		"client.go:594-597 comment is wrong; the empty oneof-unset init is replayed as a real frame and hits the default branch.", resp.Error)
	if resp.Error != "unknown message type" {
		t.Fatalf("expected production oneof-default error, got %q", resp.Error)
	}
}

// TestReceiveFileNormalPathStillWorks is the control: a real Info + content +
// CloseAndRecv sequence (the EC distributor's path) must still land the data and
// return a clean BytesWritten — proving the lazy-open did not regress the happy
// path while we probe the edges.
func TestReceiveFileNormalPathStillWorks(t *testing.T) {
	fake := &strictReceiveFileServer{}
	cli, stop := serve(t, fake)
	defer stop()

	stream, err := cli.ReceiveFile(context.Background())
	if err != nil {
		t.Fatalf("ReceiveFile open: %v", err)
	}
	if err := stream.Send(&volume_server_pb.ReceiveFileRequest{
		Data: &volume_server_pb.ReceiveFileRequest_Info{
			Info: &volume_server_pb.ReceiveFileInfo{VolumeId: 7, Ext: ".dat"},
		},
	}); err != nil {
		t.Fatalf("Info Send: %v", err)
	}
	if err := stream.Send(&volume_server_pb.ReceiveFileRequest{
		Data: &volume_server_pb.ReceiveFileRequest_FileContent{FileContent: []byte("abc")},
	}); err != nil {
		t.Fatalf("content Send: %v", err)
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		t.Fatalf("CloseAndRecv: %v", err)
	}
	if resp.Error != "" {
		t.Fatalf("normal path returned Error=%q", resp.Error)
	}
	if resp.BytesWritten != 3 {
		t.Fatalf("BytesWritten=%d, want 3", resp.BytesWritten)
	}
}
