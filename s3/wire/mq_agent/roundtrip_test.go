// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package mq_agentwire

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

// fakeAgent is an in-memory MessagingAgentStore for the round-trip test. It
// stands in for s3/mq/agent.MessageQueueAgent so the wire package proves out
// without importing the real broker/publisher engine.
type fakeAgent struct {
	nextSession int64
	closed      map[int64]bool
}

func newFakeAgent() *fakeAgent { return &fakeAgent{closed: map[int64]bool{}} }

func (a *fakeAgent) StartPublishSession(req []byte) ([]byte, error) {
	v, err := WrapStartPublishSessionRequest(req)
	if err != nil {
		return nil, err
	}
	if v.PublisherName() == "" {
		return NewStartPublishSessionResponse(StartPublishSessionResponseInput{
			Error: "publisher_name required",
		}), nil
	}
	a.nextSession++
	return NewStartPublishSessionResponse(StartPublishSessionResponseInput{
		SessionId: a.nextSession,
	}), nil
}

func (a *fakeAgent) ClosePublishSession(req []byte) ([]byte, error) {
	v, err := WrapClosePublishSessionRequest(req)
	if err != nil {
		return nil, err
	}
	a.closed[v.SessionId()] = true
	return NewClosePublishSessionResponse(ClosePublishSessionResponseInput{}), nil
}

// PublishRecord: read every request frame (first carries session id), ack each
// with an incrementing offset, then end (half-close) when the client is done.
func (a *fakeAgent) PublishRecord(stream PublishRecordServerStream) error {
	var offset int64
	for {
		buf, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		req, err := WrapPublishRecordRequest(buf)
		if err != nil {
			return err
		}
		// Echo the record's key length back through the offsets so the test can
		// prove the streamed payload arrived intact.
		base := offset
		offset += int64(len(req.Key()))
		if err := stream.Send(NewPublishRecordResponse(PublishRecordResponseInput{
			AckSequence: req.SessionId(),
			BaseOffset:  base,
			LastOffset:  offset,
		})); err != nil {
			return err
		}
	}
}

// SubscribeRecord: read the init frame, stream back a fixed batch of data
// frames keyed by the init's consumer group, then end with an end-of-stream
// marker and return (half-close).
func (a *fakeAgent) SubscribeRecord(stream SubscribeRecordServerStream) error {
	buf, err := stream.Recv()
	if err != nil {
		return err
	}
	req, err := WrapSubscribeRecordRequest(buf)
	if err != nil {
		return err
	}
	initBuf := req.Init()
	init, err := WrapInitSubscribeRecordRequest(initBuf)
	if err != nil {
		return err
	}
	group := init.ConsumerGroup()

	for i := int64(0); i < 3; i++ {
		if err := stream.Send(NewSubscribeRecordResponse(SubscribeRecordResponseInput{
			Key:    []byte(group),
			TsNs:   1000 + i,
			Offset: i,
		})); err != nil {
			return err
		}
	}
	return stream.Send(NewSubscribeRecordResponse(SubscribeRecordResponseInput{
		IsEndOfStream: true,
		Offset:        3,
	}))
}

// TestUnaryRoundTrip proves StartPublishSession then ClosePublishSession over
// the canonical github.com/zap-proto/go transport across a real TCP socket: the
// bytes cross the wire as ZAP RPC envelopes carrying zero-copy mq_agent
// payloads — no HTTP, no protobuf.
func TestUnaryRoundTrip(t *testing.T) {
	srv, err := Serve("tcp", "127.0.0.1:0", newFakeAgent())
	if err != nil {
		t.Fatalf("Serve: %v", err)
	}
	defer srv.Close()

	cli, err := Dial("tcp", srv.Addr().String())
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer cli.Close()

	start, err := cli.StartPublishSession(StartPublishSessionRequestInput{
		PartitionCount: 4,
		PublisherName:  "team-go",
	})
	if err != nil {
		t.Fatalf("StartPublishSession: %v", err)
	}
	if start.Error() != "" {
		t.Fatalf("StartPublishSession error = %q", start.Error())
	}
	if start.SessionId() != 1 {
		t.Fatalf("StartPublishSession session_id = %d, want 1", start.SessionId())
	}

	closed, err := cli.ClosePublishSession(ClosePublishSessionRequestInput{
		SessionId: start.SessionId(),
	})
	if err != nil {
		t.Fatalf("ClosePublishSession: %v", err)
	}
	if closed.Error() != "" {
		t.Fatalf("ClosePublishSession error = %q", closed.Error())
	}
}

// TestPublishRecordStream proves the PublishRecord bidirectional stream over a
// real socket: the client streams several PublishRecordRequest frames (first
// carries the session id) and reads back one PublishRecordResponse ack per
// frame, with offsets the server derived from the streamed keys.
func TestPublishRecordStream(t *testing.T) {
	srv, err := Serve("tcp", "127.0.0.1:0", newFakeAgent())
	if err != nil {
		t.Fatalf("Serve: %v", err)
	}
	defer srv.Close()

	cli, err := Dial("tcp", srv.Addr().String())
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer cli.Close()

	const sessionID = int64(7)
	keys := [][]byte{[]byte("k0"), []byte("key-1"), []byte("k2")}

	// First record rides as the stream init payload.
	stream, err := cli.PublishRecord(PublishRecordRequestInput{
		SessionId: sessionID,
		Key:       keys[0],
		Value:     []byte("v0"),
	})
	if err != nil {
		t.Fatalf("PublishRecord open: %v", err)
	}

	// Subsequent records via Send.
	for _, k := range keys[1:] {
		if err := stream.Send(PublishRecordRequestInput{
			SessionId: sessionID,
			Key:       k,
			Value:     append([]byte("v-"), k...),
		}); err != nil {
			t.Fatalf("PublishRecord send: %v", err)
		}
	}
	if err := stream.CloseSend(); err != nil {
		t.Fatalf("PublishRecord CloseSend: %v", err)
	}

	// One ack per record; offsets accumulate the streamed key lengths.
	var wantBase int64
	for i, k := range keys {
		ack, err := stream.Recv()
		if err != nil {
			t.Fatalf("PublishRecord recv[%d]: %v", i, err)
		}
		if ack.AckSequence() != sessionID {
			t.Fatalf("ack[%d] sequence = %d, want %d", i, ack.AckSequence(), sessionID)
		}
		if ack.BaseOffset() != wantBase {
			t.Fatalf("ack[%d] base_offset = %d, want %d", i, ack.BaseOffset(), wantBase)
		}
		wantBase += int64(len(k))
		if ack.LastOffset() != wantBase {
			t.Fatalf("ack[%d] last_offset = %d, want %d", i, ack.LastOffset(), wantBase)
		}
	}

	// Server half-closes after the last ack.
	if _, err := stream.Recv(); !errors.Is(err, io.EOF) {
		t.Fatalf("PublishRecord final recv = %v, want io.EOF", err)
	}
}

// TestSubscribeRecordStream proves the SubscribeRecord bidirectional stream over
// a real socket: the client sends an init frame and reads back a streamed batch
// of SubscribeRecordResponse data frames terminated by an end-of-stream marker.
func TestSubscribeRecordStream(t *testing.T) {
	srv, err := Serve("tcp", "127.0.0.1:0", newFakeAgent())
	if err != nil {
		t.Fatalf("Serve: %v", err)
	}
	defer srv.Close()

	cli, err := Dial("tcp", srv.Addr().String())
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer cli.Close()

	const group = "analytics"
	initReq := NewInitSubscribeRecordRequest(InitSubscribeRecordRequestInput{
		ConsumerGroup:           group,
		MaxSubscribedPartitions: 2,
		SlidingWindowSize:       8,
	})

	stream, err := cli.SubscribeRecord(SubscribeRecordRequestInput{Init: initReq})
	if err != nil {
		t.Fatalf("SubscribeRecord open: %v", err)
	}

	for i := int64(0); i < 3; i++ {
		data, err := stream.Recv()
		if err != nil {
			t.Fatalf("SubscribeRecord recv[%d]: %v", i, err)
		}
		if !bytes.Equal(data.Key(), []byte(group)) {
			t.Fatalf("data[%d] key = %q, want %q", i, data.Key(), group)
		}
		if data.Offset() != i {
			t.Fatalf("data[%d] offset = %d, want %d", i, data.Offset(), i)
		}
		if data.TsNs() != 1000+i {
			t.Fatalf("data[%d] ts_ns = %d, want %d", i, data.TsNs(), 1000+i)
		}
		if data.IsEndOfStream() {
			t.Fatalf("data[%d] unexpectedly end-of-stream", i)
		}
	}

	end, err := stream.Recv()
	if err != nil {
		t.Fatalf("SubscribeRecord recv end: %v", err)
	}
	if !end.IsEndOfStream() {
		t.Fatalf("final frame is_end_of_stream = false, want true")
	}
	if end.Offset() != 3 {
		t.Fatalf("final frame offset = %d, want 3", end.Offset())
	}

	if err := stream.CloseSend(); err != nil {
		t.Fatalf("SubscribeRecord CloseSend: %v", err)
	}

	// Server returned after the end marker → half-closed.
	if _, err := stream.Recv(); !errors.Is(err, io.EOF) {
		t.Fatalf("SubscribeRecord post-end recv = %v, want io.EOF", err)
	}
}
