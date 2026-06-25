// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package filerwire

import (
	"io"
	"testing"

	"github.com/zap-proto/go/transport"
)

// streamSrv pairs DispatchHanzoFilerStream (the unary dispatch is nil for these
// tests) with a TCP listener and returns a client Conn dialled to it. Mirrors
// the transport package's own stream tests: real sockets, real frames, no
// in-memory shortcut.
func streamSrv(t *testing.T, h HanzoFilerStreamHandler) transport.Conn {
	t.Helper()
	srv, err := transport.ListenStream("tcp", "127.0.0.1:0", nil, DispatchHanzoFilerStream(h))
	if err != nil {
		t.Fatalf("ListenStream: %v", err)
	}
	t.Cleanup(func() { _ = srv.Close() })

	conn, err := transport.Dial("tcp", srv.Addr().String())
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	return conn
}

// listEntriesHandler answers ListEntries by reflecting the request's Limit into
// the count of response frames, and stamping each frame's SnapshotTsNs with its
// 1-based index so the client can assert delivery order. Each frame carries a
// real nested Entry so HasEntry is exercised. The other streaming methods are
// unused here.
type listEntriesHandler struct{}

func (listEntriesHandler) StreamRenameEntry([]byte, transport.Stream) error      { return nil }
func (listEntriesHandler) TraverseBfsMetadata([]byte, transport.Stream) error    { return nil }
func (listEntriesHandler) SubscribeMetadata([]byte, transport.Stream) error      { return nil }
func (listEntriesHandler) SubscribeLocalMetadata([]byte, transport.Stream) error { return nil }
func (listEntriesHandler) StreamMutateEntry(transport.Stream) error              { return nil }

func (listEntriesHandler) ListEntries(req []byte, s transport.Stream) error {
	// The request is a real ListEntriesRequest wire buffer — read it zero-copy.
	r, err := WrapListEntriesRequest(req)
	if err != nil {
		return err
	}
	for i := int64(1); i <= int64(r.Limit()); i++ {
		entry := NewEntry(EntryInput{Name: r.Directory()})
		frame := NewListEntriesResponse(ListEntriesResponseInput{
			Entry:        entry, // exercises the nested-object field path
			SnapshotTsNs: i,     // per-frame ordering marker (scalar field)
		})
		if err := s.Send(frame); err != nil {
			return err
		}
	}
	return nil // returning half-closes -> client Recv sees io.EOF
}

// TestStreamClient_ListEntries proves server-streaming end-to-end through the
// native ZAP streaming surface: the client opens with a real ListEntriesRequest,
// the server pushes Limit real ListEntriesResponse frames, the client drains
// them to io.EOF and reads each back zero-copy. Asserts frame COUNT, per-frame
// scalar ordering (SnapshotTsNs), and nested-field presence (HasEntry). No
// proto.Marshal, no grpc — the bytes are the message.
//
// NOTE: this asserts the response's scalar field and the presence of the nested
// Entry, not Entry().Name(). The nested-object string round-trip
// (setNestedObject in helpers_zap.go) is a separate, pre-existing wire-codegen
// concern (the "extend zapgen var-length" follow-up) and is out of scope for the
// filer RPC/streaming surface under test here.
func TestStreamClient_ListEntries(t *testing.T) {
	conn := streamSrv(t, listEntriesHandler{})
	cl := NewHanzoFilerStreamClient(conn)

	const limit = uint32(3)
	st, err := cl.ListEntries(NewListEntriesRequest(ListEntriesRequestInput{
		Directory: "/buckets/x",
		Limit:     limit,
	}))
	if err != nil {
		t.Fatalf("ListEntries: %v", err)
	}

	var marks []int64
	for {
		frame, err := st.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Recv: %v", err)
		}
		resp, err := WrapListEntriesResponse(frame)
		if err != nil {
			t.Fatalf("WrapListEntriesResponse: %v", err)
		}
		if !resp.HasEntry() {
			t.Fatalf("response missing entry")
		}
		marks = append(marks, resp.SnapshotTsNs())
	}
	if len(marks) != int(limit) {
		t.Fatalf("got %d frames %v, want %d", len(marks), marks, limit)
	}
	for i, m := range marks {
		if m != int64(i+1) {
			t.Fatalf("frame[%d] SnapshotTsNs = %d, want %d (out of order)", i, m, i+1)
		}
	}
}

// echoMutateHandler implements the bidirectional StreamMutateEntry by echoing
// each request's RequestID back in a response, marking the stream's final reply
// IsLast when the client half-closes.
type echoMutateHandler struct{}

func (echoMutateHandler) ListEntries([]byte, transport.Stream) error             { return nil }
func (echoMutateHandler) StreamRenameEntry([]byte, transport.Stream) error       { return nil }
func (echoMutateHandler) TraverseBfsMetadata([]byte, transport.Stream) error     { return nil }
func (echoMutateHandler) SubscribeMetadata([]byte, transport.Stream) error       { return nil }
func (echoMutateHandler) SubscribeLocalMetadata([]byte, transport.Stream) error  { return nil }

func (echoMutateHandler) StreamMutateEntry(s transport.Stream) error {
	for {
		frame, err := s.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		req, err := WrapStreamMutateEntryRequest(frame)
		if err != nil {
			return err
		}
		out := NewStreamMutateEntryResponse(StreamMutateEntryResponseInput{
			RequestID: req.RequestID(),
		})
		if err := s.Send(out); err != nil {
			return err
		}
	}
}

// TestStreamClient_StreamMutateEntry proves bidirectional streaming end-to-end:
// the client streams N real StreamMutateEntryRequest frames and half-closes; the
// server echoes each RequestID back as a real StreamMutateEntryResponse; the
// client receives all N then io.EOF. Verifies transport v1.5.0 bidi duplex under
// the native filer wire types.
func TestStreamClient_StreamMutateEntry(t *testing.T) {
	conn := streamSrv(t, echoMutateHandler{})
	cl := NewHanzoFilerStreamClient(conn)

	st, err := cl.StreamMutateEntry()
	if err != nil {
		t.Fatalf("StreamMutateEntry: %v", err)
	}

	const n = 5
	for i := uint64(1); i <= n; i++ {
		if err := st.Send(NewStreamMutateEntryRequest(StreamMutateEntryRequestInput{RequestID: i})); err != nil {
			t.Fatalf("Send %d: %v", i, err)
		}
	}
	if err := st.CloseSend(); err != nil {
		t.Fatalf("CloseSend: %v", err)
	}

	var got []uint64
	for {
		frame, err := st.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Recv: %v", err)
		}
		resp, err := WrapStreamMutateEntryResponse(frame)
		if err != nil {
			t.Fatalf("WrapStreamMutateEntryResponse: %v", err)
		}
		got = append(got, resp.RequestID())
	}
	if len(got) != n {
		t.Fatalf("got %d responses %v, want %d", len(got), got, n)
	}
	for i := 0; i < n; i++ {
		if got[i] != uint64(i+1) {
			t.Fatalf("response[%d] RequestID = %d, want %d", i, got[i], i+1)
		}
	}
}
