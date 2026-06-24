// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package filerstream

import (
	"fmt"
	"io"
	"testing"

	filerwire "github.com/hanzoai/s3/s3/wire/filer"
)

// stubFiler is an in-memory Server for the streaming round-trip test. It does
// NOT touch the real filer engine — it scripts deterministic streams so the test
// proves the bytes survive a real duplex socket as zero-copy filerwire buffers.
type stubFiler struct {
	listNames []string // entries ListEntries streams back
}

// ListEntries (server-stream): send one ListEntriesResponse per scripted name,
// then return (which half-closes -> client Recv sees io.EOF).
func (f *stubFiler) ListEntries(init filerwire.ListEntriesRequest, s *ListEntriesStream) error {
	_ = init.Directory()
	for i, name := range f.listNames {
		entry := filerwire.NewEntry(filerwire.EntryInput{Name: name})
		frame := filerwire.NewListEntriesResponse(filerwire.ListEntriesResponseInput{
			Entry:        entry,
			SnapshotTsNs: int64(i + 1),
		})
		if err := s.Send(frame); err != nil {
			return err
		}
	}
	return nil
}

// StreamMutateEntry (bidirectional): for each request frame, reply with a
// response carrying the same request_id, until the client half-closes.
func (f *stubFiler) StreamMutateEntry(init filerwire.StreamMutateEntryRequest, s *StreamMutateEntryStream) error {
	handle := func(req filerwire.StreamMutateEntryRequest) error {
		resp := filerwire.NewStreamMutateEntryResponse(filerwire.StreamMutateEntryResponseInput{
			RequestID: req.RequestID(),
			IsLast:    true,
			Which:     filerwire.StreamMutateResponseCreate,
			CreateResponse: filerwire.NewCreateEntryResponse(filerwire.CreateEntryResponseInput{
				ErrorCode: filerwire.FilerErrorOK,
			}),
		})
		return s.Send(resp)
	}
	// First frame is the stream-open init.
	if err := handle(init); err != nil {
		return err
	}
	for {
		req, err := s.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if err := handle(req); err != nil {
			return err
		}
	}
}

// The remaining 4 streaming methods are not exercised here; they immediately
// return (ending the stream) so the Server contract is total.
func (f *stubFiler) StreamRenameEntry(filerwire.StreamRenameEntryRequest, *StreamRenameEntryStream) error {
	return nil
}
func (f *stubFiler) TraverseBfsMetadata(filerwire.TraverseBfsMetadataRequest, *TraverseBfsMetadataStream) error {
	return nil
}
func (f *stubFiler) SubscribeMetadata(filerwire.SubscribeMetadataRequest, *SubscribeMetadataStream) error {
	return nil
}
func (f *stubFiler) SubscribeLocalMetadata(filerwire.SubscribeMetadataRequest, *SubscribeMetadataStream) error {
	return nil
}

// TestListEntriesServerStream proves the ListEntries server-stream over the
// canonical github.com/zap-proto/go transport across a real TCP socket: each
// entry crosses as a zero-copy filerwire buffer framed as a stream message — no
// HTTP, no protobuf, no struct marshaling.
func TestListEntriesServerStream(t *testing.T) {
	want := []string{"a.txt", "b.txt", "c.txt"}
	srv, err := Serve("tcp", "127.0.0.1:0", &stubFiler{listNames: want})
	if err != nil {
		t.Fatalf("Serve: %v", err)
	}
	defer srv.Close()

	cli, err := Dial("tcp", srv.Addr())
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer cli.Close()

	stream, err := cli.ListEntries(filerwire.NewListEntriesRequest(filerwire.ListEntriesRequestInput{
		Directory: "/buckets/blobs",
		Limit:     100,
	}))
	if err != nil {
		t.Fatalf("ListEntries: %v", err)
	}

	var got []string
	for {
		resp, err := stream.Recv()
		if err != nil {
			if IsEOF(err) {
				break
			}
			t.Fatalf("Recv: %v", err)
		}
		if !resp.HasEntry() {
			t.Fatalf("response missing entry")
		}
		got = append(got, resp.Entry().Name())
	}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("ListEntries names = %v, want %v", got, want)
	}
}

// TestStreamMutateEntryBidi proves the bidirectional StreamMutateEntry stream
// over a real TCP socket: the client Sends mutation requests and Recvs the
// matching responses (request_id-correlated), all as zero-copy filerwire buffers.
func TestStreamMutateEntryBidi(t *testing.T) {
	srv, err := Serve("tcp", "127.0.0.1:0", &stubFiler{})
	if err != nil {
		t.Fatalf("Serve: %v", err)
	}
	defer srv.Close()

	cli, err := Dial("tcp", srv.Addr())
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer cli.Close()

	mkReq := func(id uint64, name string) []byte {
		create := filerwire.NewCreateEntryRequest(filerwire.CreateEntryRequestInput{
			Directory: "/buckets/blobs",
			Entry:     filerwire.NewEntry(filerwire.EntryInput{Name: name}),
		})
		return filerwire.NewStreamMutateEntryRequest(filerwire.StreamMutateEntryRequestInput{
			RequestID:     id,
			Which:         filerwire.StreamMutateRequestCreate,
			CreateRequest: create,
		})
	}

	// Open the stream with request_id 1 as the init frame.
	stream, err := cli.StreamMutateEntry(mkReq(1, "one.bin"))
	if err != nil {
		t.Fatalf("StreamMutateEntry: %v", err)
	}

	// First response corresponds to the init frame (request_id 1).
	resp, err := stream.Recv()
	if err != nil {
		t.Fatalf("Recv #1: %v", err)
	}
	if resp.RequestID() != 1 {
		t.Fatalf("response #1 request_id = %d, want 1", resp.RequestID())
	}
	if resp.WhichResponse() != filerwire.StreamMutateResponseCreate {
		t.Fatalf("response #1 which = %d, want create", resp.WhichResponse())
	}
	if resp.CreateResponse().ErrorCode() != filerwire.FilerErrorOK {
		t.Fatalf("response #1 create error_code = %d, want OK", resp.CreateResponse().ErrorCode())
	}

	// Send two more requests over the same stream and read their responses.
	for id := uint64(2); id <= 3; id++ {
		if err := stream.Send(mkReq(id, fmt.Sprintf("e-%d.bin", id))); err != nil {
			t.Fatalf("Send #%d: %v", id, err)
		}
		r, err := stream.Recv()
		if err != nil {
			t.Fatalf("Recv #%d: %v", id, err)
		}
		if r.RequestID() != id {
			t.Fatalf("response #%d request_id = %d, want %d", id, r.RequestID(), id)
		}
	}

	// Half-close the request side; the server's Recv loop ends and it returns.
	if err := stream.CloseSend(); err != nil {
		t.Fatalf("CloseSend: %v", err)
	}
	// After the server returns, our Recv eventually sees io.EOF.
	if _, err := stream.Recv(); err != nil && !IsEOF(err) {
		t.Fatalf("final Recv: expected io.EOF, got %v", err)
	}
}
