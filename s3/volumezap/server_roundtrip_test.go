// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// server_roundtrip_test.go proves the volume server answers over the native ZAP
// transport end to end: a proto-shaped fake volume_server_pb.VolumeServerServer
// is wrapped by NewServerBackend, served over transport.ListenStream (unary
// Dispatch + StreamHandler on one TCP listener — the exact wiring
// command/volume.go now uses), and reached by the volumezap.New client (the exact
// client pb.WithVolumeServerClient now uses). Every assertion is on the *proto*
// values that survive request decode and response encode through the
// server_rpc.go / server_convert.go converters — no HTTP, no protobuf framing on
// the wire, no struct marshaling. This is the volume analogue of the filer
// cutover smoke test.

package volumezap

import (
	"context"
	"io"
	"sync"
	"testing"

	"github.com/zap-proto/go/transport"

	volume_server_pb "github.com/hanzoai/s3/s3/pb/volume_server_pb"
	vsw "github.com/hanzoai/s3/s3/wire/volume_server"
)

// fakeVolumeServer is a proto-shaped VolumeServerServer exercising a unary RPC
// (VolumeStatus), a unary RPC with a nested/repeated response (VolumeServerStatus),
// a server-streaming RPC (VolumeTailSender), and the client-streaming RPC
// (ReceiveFile). It embeds UnimplementedVolumeServerServer so the untested RPCs
// have safe defaults. Recorded request fields are mutex-guarded (written on the
// server goroutine, read on the test goroutine).
type fakeVolumeServer struct {
	volume_server_pb.UnimplementedVolumeServerServer

	mu sync.Mutex

	gotStatusVolID uint32

	gotTailVolID uint32
	gotTailSince uint64
	gotTailIdle  uint32

	gotRecvVolID  uint32
	gotRecvExt    string
	gotRecvChunks [][]byte
}

func (f *fakeVolumeServer) VolumeStatus(_ context.Context, req *volume_server_pb.VolumeStatusRequest) (*volume_server_pb.VolumeStatusResponse, error) {
	f.mu.Lock()
	f.gotStatusVolID = req.VolumeId
	f.mu.Unlock()
	return &volume_server_pb.VolumeStatusResponse{
		IsReadOnly: true, VolumeSize: 1 << 40, FileCount: 123456, FileDeletedCount: 789,
	}, nil
}

func (f *fakeVolumeServer) VolumeServerStatus(_ context.Context, _ *volume_server_pb.VolumeServerStatusRequest) (*volume_server_pb.VolumeServerStatusResponse, error) {
	return &volume_server_pb.VolumeServerStatusResponse{
		Version: "ZAP", DataCenter: "dc1", Rack: "r1",
		DiskStatuses: []*volume_server_pb.DiskStatus{
			{Dir: "/data", All: 1000, Used: 400, Free: 600, PercentFree: 60, PercentUsed: 40, DiskType: "ssd"},
		},
		MemoryStatus: &volume_server_pb.MemStatus{Goroutines: 12, All: 2048, Used: 512, Heap: 256},
		State:        &volume_server_pb.VolumeServerState{Maintenance: true, Version: 7},
	}, nil
}

func (f *fakeVolumeServer) VolumeTailSender(req *volume_server_pb.VolumeTailSenderRequest, stream volume_server_pb.VolumeServer_VolumeTailSenderServer) error {
	f.mu.Lock()
	f.gotTailVolID = req.VolumeId
	f.gotTailSince = req.SinceNs
	f.gotTailIdle = req.IdleTimeoutSeconds
	f.mu.Unlock()
	for i := 0; i < 3; i++ {
		last := i == 2
		if err := stream.Send(&volume_server_pb.VolumeTailSenderResponse{
			NeedleHeader: []byte{byte(i), 0xAA},
			NeedleBody:   []byte{byte(i), byte(i), byte(i), byte(i)},
			IsLastChunk:  last,
			Version:      3,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (f *fakeVolumeServer) ReceiveFile(stream volume_server_pb.VolumeServer_ReceiveFileServer) error {
	var total uint64
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		f.mu.Lock()
		switch d := req.Data.(type) {
		case *volume_server_pb.ReceiveFileRequest_Info:
			if d.Info != nil {
				f.gotRecvVolID = d.Info.VolumeId
				f.gotRecvExt = d.Info.Ext
			}
		case *volume_server_pb.ReceiveFileRequest_FileContent:
			chunk := append([]byte(nil), d.FileContent...)
			f.gotRecvChunks = append(f.gotRecvChunks, chunk)
			total += uint64(len(chunk))
		}
		f.mu.Unlock()
	}
	return stream.SendAndClose(&volume_server_pb.ReceiveFileResponse{BytesWritten: total})
}

// serve stands up the volumezap server backend over a real TCP ZAP listener
// (unary Dispatch + StreamHandler, the exact command/volume.go wiring) and
// returns a pb client over the same transport (the exact pb.WithVolumeServerClient
// wiring) plus a cleanup.
func serve(t *testing.T, vs volume_server_pb.VolumeServerServer) (volume_server_pb.VolumeServerClient, func()) {
	t.Helper()
	store := NewServerBackend(vs)
	srv, err := transport.ListenStream("tcp", "127.0.0.1:0", vsw.Dispatch(store), vsw.StreamHandler(store))
	if err != nil {
		t.Fatalf("ListenStream: %v", err)
	}
	conn, err := transport.Dial("tcp", srv.Addr().String())
	if err != nil {
		srv.Close()
		t.Fatalf("Dial: %v", err)
	}
	cli := New(conn, nil)
	return cli, func() {
		conn.Close()
		srv.Close()
	}
}

// TestVolumeStatusOverZAP proves a scalar unary RPC: VolumeStatusRequest crosses
// as a ZAP envelope, the adapter decodes it to proto, the engine runs, and the
// VolumeStatusResponse comes back through the response encoder.
func TestVolumeStatusOverZAP(t *testing.T) {
	fake := &fakeVolumeServer{}
	cli, stop := serve(t, fake)
	defer stop()

	resp, err := cli.VolumeStatus(context.Background(), &volume_server_pb.VolumeStatusRequest{VolumeId: 4242})
	if err != nil {
		t.Fatalf("VolumeStatus: %v", err)
	}
	fake.mu.Lock()
	gotVolID := fake.gotStatusVolID
	fake.mu.Unlock()
	if gotVolID != 4242 {
		t.Fatalf("server saw VolumeId=%d, want 4242", gotVolID)
	}
	if !resp.IsReadOnly || resp.VolumeSize != 1<<40 || resp.FileCount != 123456 || resp.FileDeletedCount != 789 {
		t.Fatalf("VolumeStatus response lost data: %+v", resp)
	}
}

// TestVolumeServerStatusOverZAP proves a unary RPC with a nested/repeated/oneof
// response graph (disk statuses, memory status, server state) round-trips.
func TestVolumeServerStatusOverZAP(t *testing.T) {
	fake := &fakeVolumeServer{}
	cli, stop := serve(t, fake)
	defer stop()

	resp, err := cli.VolumeServerStatus(context.Background(), &volume_server_pb.VolumeServerStatusRequest{})
	if err != nil {
		t.Fatalf("VolumeServerStatus: %v", err)
	}
	if resp.Version != "ZAP" || resp.DataCenter != "dc1" || resp.Rack != "r1" {
		t.Fatalf("scalars lost: %+v", resp)
	}
	if len(resp.DiskStatuses) != 1 || resp.DiskStatuses[0].Dir != "/data" ||
		resp.DiskStatuses[0].All != 1000 || resp.DiskStatuses[0].PercentUsed != 40 {
		t.Fatalf("disk status lost: %+v", resp.DiskStatuses)
	}
	if resp.MemoryStatus == nil || resp.MemoryStatus.Goroutines != 12 || resp.MemoryStatus.Heap != 256 {
		t.Fatalf("memory status lost: %+v", resp.MemoryStatus)
	}
	if resp.State == nil || !resp.State.Maintenance || resp.State.Version != 7 {
		t.Fatalf("state lost: %+v", resp.State)
	}
}

// TestVolumeTailSenderStreamOverZAP proves a server-streaming RPC: the request
// crosses as the stream init frame, the adapter decodes it, and each engine
// stream.Send ships a response frame the client decodes — three frames, last
// flagged, then io.EOF.
func TestVolumeTailSenderStreamOverZAP(t *testing.T) {
	fake := &fakeVolumeServer{}
	cli, stop := serve(t, fake)
	defer stop()

	stream, err := cli.VolumeTailSender(context.Background(), &volume_server_pb.VolumeTailSenderRequest{
		VolumeId: 9, SinceNs: 100, IdleTimeoutSeconds: 5,
	})
	if err != nil {
		t.Fatalf("VolumeTailSender: %v", err)
	}
	var frames []*volume_server_pb.VolumeTailSenderResponse
	for {
		f, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Recv: %v", err)
		}
		frames = append(frames, f)
	}
	if len(frames) != 3 {
		t.Fatalf("got %d frames, want 3", len(frames))
	}
	if frames[0].Version != 3 || len(frames[0].NeedleHeader) != 2 || frames[0].NeedleHeader[1] != 0xAA {
		t.Fatalf("frame 0 lost data: %+v", frames[0])
	}
	if frames[2].IsLastChunk != true || frames[0].IsLastChunk != false {
		t.Fatalf("last-chunk flag wrong: %+v", frames)
	}
	fake.mu.Lock()
	defer fake.mu.Unlock()
	if fake.gotTailVolID != 9 || fake.gotTailSince != 100 || fake.gotTailIdle != 5 {
		t.Fatalf("server saw wrong tail request: vol=%d since=%d idle=%d", fake.gotTailVolID, fake.gotTailSince, fake.gotTailIdle)
	}
}

// TestReceiveFileStreamOverZAP proves the client-streaming RPC: the client sends
// an info frame then content chunks (each a ZAP frame carrying the oneof variant),
// the adapter decodes each to the proto oneof, and the single terminal reply
// (total bytes) comes back on half-close.
func TestReceiveFileStreamOverZAP(t *testing.T) {
	fake := &fakeVolumeServer{}
	cli, stop := serve(t, fake)
	defer stop()

	stream, err := cli.ReceiveFile(context.Background())
	if err != nil {
		t.Fatalf("ReceiveFile: %v", err)
	}
	if err := stream.Send(&volume_server_pb.ReceiveFileRequest{
		Data: &volume_server_pb.ReceiveFileRequest_Info{Info: &volume_server_pb.ReceiveFileInfo{VolumeId: 17, Ext: ".dat"}},
	}); err != nil {
		t.Fatalf("Send info: %v", err)
	}
	chunks := [][]byte{[]byte("hello"), []byte(" "), []byte("world")}
	for _, c := range chunks {
		if err := stream.Send(&volume_server_pb.ReceiveFileRequest{
			Data: &volume_server_pb.ReceiveFileRequest_FileContent{FileContent: c},
		}); err != nil {
			t.Fatalf("Send content: %v", err)
		}
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		t.Fatalf("CloseAndRecv: %v", err)
	}
	if resp.BytesWritten != 11 {
		t.Fatalf("BytesWritten=%d, want 11", resp.BytesWritten)
	}
	fake.mu.Lock()
	defer fake.mu.Unlock()
	if fake.gotRecvVolID != 17 || fake.gotRecvExt != ".dat" {
		t.Fatalf("server saw wrong info: vol=%d ext=%q", fake.gotRecvVolID, fake.gotRecvExt)
	}
	if len(fake.gotRecvChunks) != 3 || string(fake.gotRecvChunks[0]) != "hello" || string(fake.gotRecvChunks[2]) != "world" {
		t.Fatalf("server saw wrong chunks: %q", fake.gotRecvChunks)
	}
}
