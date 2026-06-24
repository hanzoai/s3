// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package pluginstream

import (
	"io"
	"sync"
	"testing"

	pluginwire "github.com/hanzoai/s3/s3/wire/plugin"
)

// stubAdmin is an in-memory Admin for the round-trip test. It does NOT touch
// the real admin control plane — it just echoes a scripted exchange so the test
// proves the bytes survive a real duplex socket.
//
// Protocol exercised:
//   worker --(open: WorkerHello)-->        admin
//   admin  --(AdminHello accepted)-->      worker
//   worker --(WorkerHeartbeat)-->          admin
//   admin  --(ExecuteJobRequest)-->        worker
//   worker --(CloseSend)-->                admin   (Recv -> io.EOF)
//   admin  --(return -> CloseSend)-->       worker  (Recv -> io.EOF)
type stubAdmin struct {
	mu       sync.Mutex
	gotHello string // worker_id from the opening WorkerHello
	gotHB    int32  // detection_slots_total from the WorkerHeartbeat
}

func (a *stubAdmin) ServeWorker(init pluginwire.WorkerToAdminMessage, s *ServerStream) error {
	// init is the worker's opening WorkerToAdminMessage carrying a WorkerHello.
	a.mu.Lock()
	a.gotHello = init.Hello().WorkerID()
	a.mu.Unlock()

	// Reply with an AdminHello(accepted) wrapped in an AdminToWorkerMessage.
	helloBody := pluginwire.NewAdminHello(pluginwire.AdminHelloInput{
		Accepted:                 true,
		Message:                  "welcome",
		HeartbeatIntervalSeconds: 30,
	})
	if err := s.Send(pluginwire.NewAdminToWorkerMessage(pluginwire.AdminToWorkerMessageInput{
		RequestID: "r1",
		Which:     pluginwire.AdminToWorkerBodyHello,
		Hello:     helloBody,
	})); err != nil {
		return err
	}

	// Receive the worker's heartbeat.
	hb, err := s.Recv()
	if err != nil {
		return err
	}
	a.mu.Lock()
	a.gotHB = hb.Heartbeat().DetectionSlotsTotal()
	a.mu.Unlock()

	// Dispatch one job to the worker.
	jobBody := pluginwire.NewExecuteJobRequest(pluginwire.ExecuteJobRequestInput{
		RequestID: "job-req-1",
		Attempt:   1,
	})
	if err := s.Send(pluginwire.NewAdminToWorkerMessage(pluginwire.AdminToWorkerMessageInput{
		RequestID:         "r2",
		Which:             pluginwire.AdminToWorkerBodyExecuteJobRequest,
		ExecuteJobRequest: jobBody,
	})); err != nil {
		return err
	}

	// Drain until the worker half-closes, then return (auto half-close back).
	for {
		if _, err := s.Recv(); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

// TestWorkerStreamRoundTrip proves the bidirectional PluginControlService
// stream over the canonical github.com/zap-proto/go transport across a real TCP
// socket: every frame crosses as a zero-copy pluginwire buffer framed as a
// stream message — no HTTP, no protobuf, no struct marshaling.
func TestWorkerStreamRoundTrip(t *testing.T) {
	admin := &stubAdmin{}

	srv, err := Serve("tcp", "127.0.0.1:0", admin)
	if err != nil {
		t.Fatalf("Serve: %v", err)
	}
	defer srv.Close()

	cli, err := Dial("tcp", srv.Addr())
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer cli.Close()

	// Worker opens the stream with a WorkerHello.
	helloBody := pluginwire.NewWorkerHello(pluginwire.WorkerHelloInput{
		WorkerID:        "w-42",
		Address:         "10.0.0.7:9000",
		WorkerVersion:   "v1.2.3",
		ProtocolVersion: "1",
	})
	stream, err := cli.OpenWorkerStream(pluginwire.NewWorkerToAdminMessage(pluginwire.WorkerToAdminMessageInput{
		WorkerID: "w-42",
		Which:    pluginwire.WorkerToAdminBodyHello,
		Hello:    helloBody,
	}))
	if err != nil {
		t.Fatalf("OpenWorkerStream: %v", err)
	}

	// Admin must answer with AdminHello(accepted).
	first, err := stream.Recv()
	if err != nil {
		t.Fatalf("Recv AdminHello: %v", err)
	}
	if w := first.WhichBody(); w != pluginwire.AdminToWorkerBodyHello {
		t.Fatalf("first frame WhichBody = %d, want Hello(%d)", w, pluginwire.AdminToWorkerBodyHello)
	}
	if !first.Hello().Accepted() {
		t.Fatalf("AdminHello.Accepted = false, want true")
	}
	if got := first.Hello().Message(); got != "welcome" {
		t.Fatalf("AdminHello.Message = %q, want %q", got, "welcome")
	}

	// Worker sends a heartbeat.
	hbBody := pluginwire.NewWorkerHeartbeat(pluginwire.WorkerHeartbeatInput{
		WorkerID:            "w-42",
		DetectionSlotsUsed:  1,
		DetectionSlotsTotal: 4,
	})
	if err := stream.Send(pluginwire.NewWorkerToAdminMessage(pluginwire.WorkerToAdminMessageInput{
		WorkerID:  "w-42",
		Which:     pluginwire.WorkerToAdminBodyHeartbeat,
		Heartbeat: hbBody,
	})); err != nil {
		t.Fatalf("Send heartbeat: %v", err)
	}

	// Admin must dispatch a job.
	second, err := stream.Recv()
	if err != nil {
		t.Fatalf("Recv ExecuteJobRequest: %v", err)
	}
	if w := second.WhichBody(); w != pluginwire.AdminToWorkerBodyExecuteJobRequest {
		t.Fatalf("second frame WhichBody = %d, want ExecuteJobRequest(%d)", w, pluginwire.AdminToWorkerBodyExecuteJobRequest)
	}
	if got := second.ExecuteJobRequest().RequestID(); got != "job-req-1" {
		t.Fatalf("ExecuteJobRequest.RequestID = %q, want %q", got, "job-req-1")
	}

	// Worker is done sending: half-close. Admin's Recv loop then returns and
	// half-closes back, so our next Recv is io.EOF.
	if err := stream.CloseSend(); err != nil {
		t.Fatalf("CloseSend: %v", err)
	}
	if _, err := stream.Recv(); err != io.EOF {
		t.Fatalf("final Recv err = %v, want io.EOF", err)
	}

	// Verify the admin saw the worker-originated fields decoded zero-copy.
	admin.mu.Lock()
	gotHello, gotHB := admin.gotHello, admin.gotHB
	admin.mu.Unlock()
	if gotHello != "w-42" {
		t.Fatalf("admin saw WorkerHello.WorkerID = %q, want %q", gotHello, "w-42")
	}
	if gotHB != 4 {
		t.Fatalf("admin saw WorkerHeartbeat.DetectionSlotsTotal = %d, want 4", gotHB)
	}
}
