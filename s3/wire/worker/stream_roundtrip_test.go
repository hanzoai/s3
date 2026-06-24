// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package workerwire

import (
	"io"
	"testing"

	"github.com/zap-proto/go/transport"
)

// echoAdmin is an in-package WorkerStreamServer used by the round-trip test: it
// Recvs WorkerMessage frames and, for each, Sends one AdminMessage reply whose
// variant mirrors the request (registration -> registration_response, heartbeat
// -> heartbeat_response). It exercises the duplex loop, the oneof unions, and
// the zero-copy accessors end to end. The real admin dashboard implements the
// same WorkerStreamServer interface; this package never imports it.
type echoAdmin struct {
	adminID string
}

func (e echoAdmin) ServeWorkerStream(stream *AdminStream) {
	// Drain the opening frame first if present, then the streamed frames.
	if init := stream.Init(); len(init) > 0 {
		e.reply(stream, init)
	}
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return // worker half-closed; returning half-closes our side.
		}
		if err != nil {
			return
		}
		e.reply(stream, req)
	}
}

func (e echoAdmin) reply(stream *AdminStream, reqEnv []byte) {
	w, err := WrapWorkerMessage(reqEnv)
	if err != nil {
		return
	}
	out := AdminMessageInput{AdminID: e.adminID, Timestamp: w.Timestamp() + 1}
	switch w.WhichMessage() {
	case WorkerMessageMessageRegistration:
		reg := w.Registration()
		out.MessageWhich = AdminMessageMessageRegistrationResponse
		out.MessageValue = NewRegistrationResponse(RegistrationResponseInput{
			Success:          true,
			Message:          "registered " + reg.WorkerID(),
			AssignedWorkerID: reg.WorkerID(),
		})
	case WorkerMessageMessageHeartbeat:
		hb := w.Heartbeat()
		out.MessageWhich = AdminMessageMessageHeartbeatResponse
		out.MessageValue = NewHeartbeatResponse(HeartbeatResponseInput{
			Success: true,
			Message: "ack " + hb.Status(),
		})
	default:
		out.MessageWhich = AdminMessageMessageAdminShutdown
		out.MessageValue = NewAdminShutdown(AdminShutdownInput{Reason: "unknown variant"})
	}
	_ = stream.Send(NewAdminMessage(out))
}

// TestWorkerStreamRoundTrip proves the WorkerService.WorkerStream bidi stream
// over the canonical github.com/zap-proto/go transport across a real TCP
// socket: WorkerMessage frames cross up, AdminMessage frames cross back, all as
// zero-copy ZAP buffers — no HTTP, no protobuf, no marshal.
func TestWorkerStreamRoundTrip(t *testing.T) {
	srv, err := ListenWorkerStream("tcp", "127.0.0.1:0", echoAdmin{adminID: "admin-1"}, nil)
	if err != nil {
		t.Fatalf("ListenWorkerStream: %v", err)
	}
	defer srv.Close()

	conn, err := transport.Dial("tcp", srv.Addr().String())
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer conn.Close()

	// Open the stream with an initial registration frame.
	reg := NewWorkerMessage(WorkerMessageInput{
		WorkerID:     "worker-7",
		Timestamp:    100,
		MessageWhich: WorkerMessageMessageRegistration,
		MessageValue: NewWorkerRegistration(WorkerRegistrationInput{
			WorkerID:      "worker-7",
			Address:       "10.0.0.7:9000",
			Capabilities:  []string{"vacuum", "ec"},
			MaxConcurrent: 4,
		}),
	})
	ws, err := OpenWorkerStream(conn, reg)
	if err != nil {
		t.Fatalf("OpenWorkerStream: %v", err)
	}

	// Reply to the registration frame.
	regRespEnv, err := ws.Recv()
	if err != nil {
		t.Fatalf("Recv registration response: %v", err)
	}
	regResp, err := WrapAdminMessage(regRespEnv)
	if err != nil {
		t.Fatalf("WrapAdminMessage: %v", err)
	}
	if regResp.AdminID() != "admin-1" {
		t.Fatalf("AdminID = %q, want admin-1", regResp.AdminID())
	}
	if regResp.Timestamp() != 101 {
		t.Fatalf("Timestamp = %d, want 101", regResp.Timestamp())
	}
	if regResp.WhichMessage() != AdminMessageMessageRegistrationResponse {
		t.Fatalf("WhichMessage = %d, want RegistrationResponse", regResp.WhichMessage())
	}
	if rr := regResp.RegistrationResponse(); !rr.Success() || rr.AssignedWorkerID() != "worker-7" {
		t.Fatalf("RegistrationResponse = {success:%v, id:%q}", rr.Success(), rr.AssignedWorkerID())
	}

	// Now stream a heartbeat frame and read its ack — proves the loop, not
	// just the opening frame.
	hb := NewWorkerMessage(WorkerMessageInput{
		WorkerID:     "worker-7",
		Timestamp:    200,
		MessageWhich: WorkerMessageMessageHeartbeat,
		MessageValue: NewWorkerHeartbeat(WorkerHeartbeatInput{
			WorkerID:    "worker-7",
			Status:      "busy",
			CurrentLoad: 2,
		}),
	})
	if err := ws.Send(hb); err != nil {
		t.Fatalf("Send heartbeat: %v", err)
	}
	hbRespEnv, err := ws.Recv()
	if err != nil {
		t.Fatalf("Recv heartbeat response: %v", err)
	}
	hbResp, err := WrapAdminMessage(hbRespEnv)
	if err != nil {
		t.Fatalf("WrapAdminMessage heartbeat: %v", err)
	}
	if hbResp.WhichMessage() != AdminMessageMessageHeartbeatResponse {
		t.Fatalf("heartbeat WhichMessage = %d, want HeartbeatResponse", hbResp.WhichMessage())
	}
	if hr := hbResp.HeartbeatResponse(); !hr.Success() || hr.Message() != "ack busy" {
		t.Fatalf("HeartbeatResponse = {success:%v, msg:%q}", hr.Success(), hr.Message())
	}

	// Half-close the worker side; admin's loop ends and half-closes back, so
	// our next Recv sees io.EOF.
	if err := ws.CloseSend(); err != nil {
		t.Fatalf("CloseSend: %v", err)
	}
	if _, err := ws.Recv(); err != io.EOF {
		t.Fatalf("post-CloseSend Recv = %v, want io.EOF", err)
	}
}
