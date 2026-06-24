// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package volume_serverwire

import (
	"bytes"
	"errors"
	"io"
	"sync"
	"testing"
)

// errUnimplemented is returned by every VolumeServerStore method the test does
// not exercise. A unary call to such a method comes back StatusInternal; a stream
// open to one returns immediately (half-closing the stream).
var errUnimplemented = errors.New("volume_serverwire: not implemented in test fake")

// unimplementedStore supplies a safe default for all 48 VolumeServerStore
// methods so a fake can embed it and override only the few RPCs under test. This
// is the standard Go "unimplemented base" pattern (the same shape protoc's
// UnimplementedXServer used to provide), kept here so the round-trip fakes stay
// tiny.
type unimplementedStore struct{}

func (unimplementedStore) BatchDelete([]byte) ([]byte, error)          { return nil, errUnimplemented }
func (unimplementedStore) VacuumVolumeCheck([]byte) ([]byte, error)    { return nil, errUnimplemented }
func (unimplementedStore) VacuumVolumeCommit([]byte) ([]byte, error)   { return nil, errUnimplemented }
func (unimplementedStore) VacuumVolumeCleanup([]byte) ([]byte, error)  { return nil, errUnimplemented }
func (unimplementedStore) DeleteCollection([]byte) ([]byte, error)     { return nil, errUnimplemented }
func (unimplementedStore) AllocateVolume([]byte) ([]byte, error)       { return nil, errUnimplemented }
func (unimplementedStore) VolumeSyncStatus([]byte) ([]byte, error)     { return nil, errUnimplemented }
func (unimplementedStore) VolumeMount([]byte) ([]byte, error)          { return nil, errUnimplemented }
func (unimplementedStore) VolumeUnmount([]byte) ([]byte, error)        { return nil, errUnimplemented }
func (unimplementedStore) VolumeDelete([]byte) ([]byte, error)         { return nil, errUnimplemented }
func (unimplementedStore) VolumeMarkReadonly([]byte) ([]byte, error)   { return nil, errUnimplemented }
func (unimplementedStore) VolumeMarkWritable([]byte) ([]byte, error)   { return nil, errUnimplemented }
func (unimplementedStore) VolumeConfigure([]byte) ([]byte, error)      { return nil, errUnimplemented }
func (unimplementedStore) VolumeStatus([]byte) ([]byte, error)         { return nil, errUnimplemented }
func (unimplementedStore) GetState([]byte) ([]byte, error)             { return nil, errUnimplemented }
func (unimplementedStore) SetState([]byte) ([]byte, error)             { return nil, errUnimplemented }
func (unimplementedStore) ReadVolumeFileStatus([]byte) ([]byte, error) { return nil, errUnimplemented }
func (unimplementedStore) ReadNeedleBlob([]byte) ([]byte, error)       { return nil, errUnimplemented }
func (unimplementedStore) ReadNeedleMeta([]byte) ([]byte, error)       { return nil, errUnimplemented }
func (unimplementedStore) WriteNeedleBlob([]byte) ([]byte, error)      { return nil, errUnimplemented }
func (unimplementedStore) VolumeTailReceiver([]byte) ([]byte, error)   { return nil, errUnimplemented }
func (unimplementedStore) VolumeEcShardsGenerate([]byte) ([]byte, error) {
	return nil, errUnimplemented
}
func (unimplementedStore) VolumeEcShardsRebuild([]byte) ([]byte, error) { return nil, errUnimplemented }
func (unimplementedStore) VolumeEcShardsCopy([]byte) ([]byte, error)    { return nil, errUnimplemented }
func (unimplementedStore) VolumeEcShardsDelete([]byte) ([]byte, error)  { return nil, errUnimplemented }
func (unimplementedStore) VolumeEcShardsMount([]byte) ([]byte, error)   { return nil, errUnimplemented }
func (unimplementedStore) VolumeEcShardsUnmount([]byte) ([]byte, error) { return nil, errUnimplemented }
func (unimplementedStore) VolumeEcBlobDelete([]byte) ([]byte, error)    { return nil, errUnimplemented }
func (unimplementedStore) VolumeEcShardsToVolume([]byte) ([]byte, error) {
	return nil, errUnimplemented
}
func (unimplementedStore) VolumeEcShardsInfo([]byte) ([]byte, error)  { return nil, errUnimplemented }
func (unimplementedStore) VolumeServerStatus([]byte) ([]byte, error)  { return nil, errUnimplemented }
func (unimplementedStore) VolumeServerLeave([]byte) ([]byte, error)   { return nil, errUnimplemented }
func (unimplementedStore) FetchAndWriteNeedle([]byte) ([]byte, error) { return nil, errUnimplemented }
func (unimplementedStore) ScrubVolume([]byte) ([]byte, error)         { return nil, errUnimplemented }
func (unimplementedStore) ScrubEcVolume([]byte) ([]byte, error)       { return nil, errUnimplemented }
func (unimplementedStore) VolumeNeedleStatus([]byte) ([]byte, error)  { return nil, errUnimplemented }
func (unimplementedStore) Ping([]byte) ([]byte, error)                { return nil, errUnimplemented }

func (unimplementedStore) VacuumVolumeCompact(VacuumVolumeCompactServerStream) error {
	return errUnimplemented
}
func (unimplementedStore) VolumeIncrementalCopy(VolumeIncrementalCopyServerStream) error {
	return errUnimplemented
}
func (unimplementedStore) VolumeCopy(VolumeCopyServerStream) error { return errUnimplemented }
func (unimplementedStore) CopyFile(CopyFileServerStream) error     { return errUnimplemented }
func (unimplementedStore) ReadAllNeedles(ReadAllNeedlesServerStream) error {
	return errUnimplemented
}
func (unimplementedStore) VolumeTailSender(VolumeTailSenderServerStream) error {
	return errUnimplemented
}
func (unimplementedStore) VolumeEcShardRead(VolumeEcShardReadServerStream) error {
	return errUnimplemented
}
func (unimplementedStore) VolumeTierMoveDatToRemote(VolumeTierMoveDatToRemoteServerStream) error {
	return errUnimplemented
}
func (unimplementedStore) VolumeTierMoveDatFromRemote(VolumeTierMoveDatFromRemoteServerStream) error {
	return errUnimplemented
}
func (unimplementedStore) Query(QueryServerStream) error             { return errUnimplemented }
func (unimplementedStore) ReceiveFile(ReceiveFileServerStream) error { return errUnimplemented }

// fakeVolumeServer is an in-memory VolumeServerStore exercising one unary RPC
// (VolumeStatus), one server-streaming RPC (VolumeTailSender), and the one
// client-streaming RPC (ReceiveFile). Fields the server records off the
// zero-copy request views are guarded by a mutex (written on the server
// goroutine, read on the test goroutine).
type fakeVolumeServer struct {
	unimplementedStore

	mu sync.Mutex

	// VolumeStatus
	gotStatusVolID uint32

	// VolumeTailSender
	gotTailVolID uint32
	gotTailSince uint64
	gotTailIdle  uint32

	// ReceiveFile
	gotRecvVolID  uint32
	gotRecvExt    string
	gotRecvChunks [][]byte
}

func (f *fakeVolumeServer) VolumeStatus(req []byte) ([]byte, error) {
	v, err := WrapVolumeStatusRequest(req)
	if err != nil {
		return nil, err
	}
	f.mu.Lock()
	f.gotStatusVolID = v.VolumeID()
	f.mu.Unlock()
	return NewVolumeStatusResponse(VolumeStatusResponseInput{
		IsReadOnly:       true,
		VolumeSize:       1 << 40,
		FileCount:        123456,
		FileDeletedCount: 789,
	}), nil
}

// VolumeTailSender streams three needle frames then ends (returning half-closes,
// so the client's Recv sees io.EOF after the third).
func (f *fakeVolumeServer) VolumeTailSender(s VolumeTailSenderServerStream) error {
	req, err := s.Init()
	if err != nil {
		return err
	}
	f.mu.Lock()
	f.gotTailVolID = req.VolumeID()
	f.gotTailSince = req.SinceNs()
	f.gotTailIdle = req.IdleTimeoutSeconds()
	f.mu.Unlock()

	for i := 0; i < 3; i++ {
		last := i == 2
		if err := s.Send(VolumeTailSenderResponseInput{
			NeedleHeader: []byte{byte(i), 0xAA},
			NeedleBody:   bytes.Repeat([]byte{byte(i)}, 4),
			IsLastChunk:  last,
			Version:      3,
		}); err != nil {
			return err
		}
	}
	return nil
}

// ReceiveFile drains the request frames (first = info, then content chunks) and
// replies with the total bytes written.
func (f *fakeVolumeServer) ReceiveFile(s ReceiveFileServerStream) error {
	var total uint64
	for {
		req, err := s.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		f.mu.Lock()
		switch req.Data() {
		case ReceiveFileDataInfo:
			if info, ok := req.Info(); ok {
				f.gotRecvVolID = info.VolumeID()
				f.gotRecvExt = info.Ext()
			}
		case ReceiveFileDataFileContent:
			if content, ok := req.FileContent(); ok {
				chunk := append([]byte(nil), content...)
				f.gotRecvChunks = append(f.gotRecvChunks, chunk)
				total += uint64(len(chunk))
			}
		}
		f.mu.Unlock()
	}
	return s.Reply(ReceiveFileResponseInput{BytesWritten: total})
}

// TestVolumeServerUnaryRoundTrip proves a unary RPC (VolumeStatus) over the
// canonical github.com/zap-proto/go transport across a real TCP socket: the
// request crosses as a ZAP RPC envelope carrying the zero-copy
// VolumeStatusRequest payload, the server dispatches it, and the
// VolumeStatusResponse comes back the same way — no HTTP, no protobuf, no struct
// marshaling.
func TestVolumeServerUnaryRoundTrip(t *testing.T) {
	store := &fakeVolumeServer{}

	srv, err := Serve("tcp", "127.0.0.1:0", store)
	if err != nil {
		t.Fatalf("Serve: %v", err)
	}
	defer srv.Close()

	cli, err := Dial("tcp", srv.Addr().String())
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer cli.Close()

	resp, err := cli.VolumeStatus(VolumeStatusRequestInput{VolumeID: 4242})
	if err != nil {
		t.Fatalf("VolumeStatus: %v", err)
	}

	store.mu.Lock()
	gotVolID := store.gotStatusVolID
	store.mu.Unlock()
	if gotVolID != 4242 {
		t.Fatalf("server saw VolumeID = %d, want 4242", gotVolID)
	}

	if !resp.IsReadOnly() {
		t.Fatalf("IsReadOnly = false, want true")
	}
	if resp.VolumeSize() != 1<<40 {
		t.Fatalf("VolumeSize = %d, want %d", resp.VolumeSize(), uint64(1)<<40)
	}
	if resp.FileCount() != 123456 || resp.FileDeletedCount() != 789 {
		t.Fatalf("counts mismatch: fileCount=%d deleted=%d", resp.FileCount(), resp.FileDeletedCount())
	}
}

// TestVolumeServerServerStreamRoundTrip proves a server-streaming RPC
// (VolumeTailSender) over the transport: the request rides as the stream's init
// payload, the server streams three VolumeTailSenderResponse frames as zero-copy
// buffers, and the client reads them until io.EOF.
func TestVolumeServerServerStreamRoundTrip(t *testing.T) {
	store := &fakeVolumeServer{}

	srv, err := Serve("tcp", "127.0.0.1:0", store)
	if err != nil {
		t.Fatalf("Serve: %v", err)
	}
	defer srv.Close()

	cli, err := Dial("tcp", srv.Addr().String())
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer cli.Close()

	stream, err := cli.VolumeTailSender(VolumeTailSenderRequestInput{
		VolumeID:           7,
		SinceNs:            1_700_000_000_000_000_000,
		IdleTimeoutSeconds: 30,
	})
	if err != nil {
		t.Fatalf("VolumeTailSender open: %v", err)
	}

	var got int
	for {
		item, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Recv #%d: %v", got, err)
		}
		wantHdr := []byte{byte(got), 0xAA}
		if !bytes.Equal(item.NeedleHeader(), wantHdr) {
			t.Fatalf("item #%d header = %x, want %x", got, item.NeedleHeader(), wantHdr)
		}
		if !bytes.Equal(item.NeedleBody(), bytes.Repeat([]byte{byte(got)}, 4)) {
			t.Fatalf("item #%d body mismatch", got)
		}
		if item.Version() != 3 {
			t.Fatalf("item #%d version = %d, want 3", got, item.Version())
		}
		if (got == 2) != item.IsLastChunk() {
			t.Fatalf("item #%d IsLastChunk = %v", got, item.IsLastChunk())
		}
		got++
	}
	if got != 3 {
		t.Fatalf("received %d items, want 3", got)
	}

	store.mu.Lock()
	defer store.mu.Unlock()
	if store.gotTailVolID != 7 || store.gotTailSince != 1_700_000_000_000_000_000 || store.gotTailIdle != 30 {
		t.Fatalf("server saw req vol=%d since=%d idle=%d", store.gotTailVolID, store.gotTailSince, store.gotTailIdle)
	}
}

// TestVolumeServerClientStreamRoundTrip proves the one client-streaming RPC
// (ReceiveFile) over the transport: the client streams an info frame (oneof
// member 1) then two content frames (oneof member 2) as zero-copy buffers,
// half-closes, and the server replies with one ReceiveFileResponse.
func TestVolumeServerClientStreamRoundTrip(t *testing.T) {
	store := &fakeVolumeServer{}

	srv, err := Serve("tcp", "127.0.0.1:0", store)
	if err != nil {
		t.Fatalf("Serve: %v", err)
	}
	defer srv.Close()

	cli, err := Dial("tcp", srv.Addr().String())
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer cli.Close()

	// First frame carries the ReceiveFileInfo (oneof member 1). The info sub-
	// message is its own ZAP buffer (NewReceiveFileInfo), embedded as the
	// variant payload.
	stream, err := cli.ReceiveFile(ReceiveFileRequestInput{
		Data: ReceiveFileDataInfo,
		Info: NewReceiveFileInfo(ReceiveFileInfoInput{
			VolumeID: 11,
			Ext:      ".dat",
		}),
	})
	if err != nil {
		t.Fatalf("ReceiveFile open: %v", err)
	}

	chunks := [][]byte{[]byte("hello "), []byte("volume-server")}
	var want uint64
	for _, ch := range chunks {
		if err := stream.Send(ReceiveFileRequestInput{
			Data:        ReceiveFileDataFileContent,
			FileContent: ch,
		}); err != nil {
			t.Fatalf("Send: %v", err)
		}
		want += uint64(len(ch))
	}
	if err := stream.CloseSend(); err != nil {
		t.Fatalf("CloseSend: %v", err)
	}

	reply, err := stream.Reply()
	if err != nil {
		t.Fatalf("Reply: %v", err)
	}
	if reply.BytesWritten() != want {
		t.Fatalf("BytesWritten = %d, want %d", reply.BytesWritten(), want)
	}

	store.mu.Lock()
	defer store.mu.Unlock()
	if store.gotRecvVolID != 11 || store.gotRecvExt != ".dat" {
		t.Fatalf("server saw info vol=%d ext=%q", store.gotRecvVolID, store.gotRecvExt)
	}
	if len(store.gotRecvChunks) != 2 ||
		!bytes.Equal(store.gotRecvChunks[0], chunks[0]) ||
		!bytes.Equal(store.gotRecvChunks[1], chunks[1]) {
		t.Fatalf("server saw %d content chunks, want 2 matching", len(store.gotRecvChunks))
	}
}
