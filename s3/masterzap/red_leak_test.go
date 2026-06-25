// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// RED adversarial probes — NOT Blue's. These document attacks on the
// master-over-ZAP cut. Keep in the worktree for integration triage.

package masterzap

import (
	"context"
	"io"
	"sync/atomic"
	"testing"
	"time"

	master_pb "github.com/hanzoai/s3/s3/pb/master_pb"
	masterwire "github.com/hanzoai/s3/s3/wire/master"
	masterstream "github.com/hanzoai/s3/s3/wire/master/masterstream"
	"github.com/zap-proto/go/transport"
)

// leakMaster counts entries/exits of each streaming handler so a test can prove
// whether a handler returns (no leak) or stays blocked (leak) after the client
// abandons its half of the stream.
type leakMaster struct {
	master_pb.UnimplementedHanzoServer
	kcLive atomic.Int32
	hbLive atomic.Int32
}

func (m *leakMaster) KeepConnected(stream master_pb.Hanzo_KeepConnectedServer) error {
	m.kcLive.Add(1)
	defer m.kcLive.Add(-1)
	// Mirror the real engine: read the subscription, answer once, then keep
	// reading until the client half-closes (io.EOF) or the conn drops.
	req, err := stream.Recv()
	if err != nil {
		return err
	}
	if err := stream.Send(&master_pb.KeepConnectedResponse{
		VolumeLocation: &master_pb.VolumeLocation{Url: "vs-1:18080", DataCenter: req.DataCenter},
	}); err != nil {
		return err
	}
	for {
		if _, err := stream.Recv(); err != nil {
			return err
		}
	}
}

func (m *leakMaster) SendHeartbeat(stream master_pb.Hanzo_SendHeartbeatServer) error {
	m.hbLive.Add(1)
	defer m.hbLive.Add(-1)
	for {
		if _, err := stream.Recv(); err != nil {
			return err
		}
	}
}

func serveLeak(t *testing.T) (master_pb.HanzoClient, *leakMaster, transport.Conn, func()) {
	t.Helper()
	ms := &leakMaster{}
	srv, err := transport.ListenStream("tcp", "127.0.0.1:0",
		masterwire.Dispatch(NewServerBackend(ms)),
		masterstream.Handler(NewStreamServer(ms)))
	if err != nil {
		t.Fatalf("ListenStream: %v", err)
	}
	conn, err := transport.Dial("tcp", srv.Addr().String())
	if err != nil {
		srv.Close()
		t.Fatalf("Dial: %v", err)
	}
	return New(conn, nil), ms, conn, func() { conn.Close(); srv.Close() }
}

func waitFor(d time.Duration, cond func() bool) bool {
	deadline := time.Now().Add(d)
	for time.Now().Before(deadline) {
		if cond() {
			return true
		}
		time.Sleep(2 * time.Millisecond)
	}
	return cond()
}

// TestRedKeepConnectedAbandonLeaksServerHandler models wdclient.tryConnectToMaster
// on the leader-redirect path: it opens KeepConnected, Sends the subscription,
// Recvs the first response, then RETURNS from fn WITHOUT CloseSend (exactly what
// masterclient.go does at `return nil` on redirect) while the pooled conn stays
// alive (the Pool does not Close it). The server handler must observe the abandon
// and return. If it stays blocked, every leader redirect leaks a server goroutine
// (and a client streams-map entry) for the life of the shared pooled conn.
func TestRedKeepConnectedAbandonLeaksServerHandler(t *testing.T) {
	cli, ms, _, done := serveLeak(t)
	defer done()

	stream, err := cli.KeepConnected(context.Background())
	if err != nil {
		t.Fatalf("KeepConnected: %v", err)
	}
	if err := stream.Send(&master_pb.KeepConnectedRequest{ClientType: "filer", DataCenter: "dc1"}); err != nil {
		t.Fatalf("Send: %v", err)
	}
	if _, err := stream.Recv(); err != nil {
		t.Fatalf("Recv: %v", err)
	}
	if !waitFor(time.Second, func() bool { return ms.kcLive.Load() == 1 }) {
		t.Fatalf("handler never started; live=%d", ms.kcLive.Load())
	}

	// fn returns here on the real redirect path. The conn is a SHARED pooled conn
	// that stays alive, so the caller MUST half-close — the fix adds
	// `defer stream.CloseSend()` to tryConnectToMaster (masterclient.go) and
	// runHeartbeatStream (volume_grpc_client_to_master.go). Mimic that defer and
	// assert the server KeepConnected handler releases (no goroutine leak per
	// leader-redirect). Without the CloseSend this stays live → the original leak.
	if err := stream.CloseSend(); err != nil {
		t.Fatalf("CloseSend: %v", err)
	}

	released := waitFor(2*time.Second, func() bool { return ms.kcLive.Load() == 0 })
	if !released {
		t.Fatalf("LEAK NOT FIXED: server KeepConnected handler still live (%d) 2s after the caller's defer stream.CloseSend(); the abandon-path fix is ineffective.", ms.kcLive.Load())
	}
}

// TestRedKeepConnectedCloseSendReleases is the control: WITH CloseSend the server
// handler returns. Proves the leak above is caused solely by the missing
// CloseSend on the caller's redirect/stop path, not by the adapter itself.
func TestRedKeepConnectedCloseSendReleases(t *testing.T) {
	cli, ms, _, done := serveLeak(t)
	defer done()

	stream, err := cli.KeepConnected(context.Background())
	if err != nil {
		t.Fatalf("KeepConnected: %v", err)
	}
	if err := stream.Send(&master_pb.KeepConnectedRequest{ClientType: "filer", DataCenter: "dc1"}); err != nil {
		t.Fatalf("Send: %v", err)
	}
	if _, err := stream.Recv(); err != nil {
		t.Fatalf("Recv: %v", err)
	}
	if !waitFor(time.Second, func() bool { return ms.kcLive.Load() == 1 }) {
		t.Fatalf("handler never started")
	}
	if err := stream.CloseSend(); err != nil {
		t.Fatalf("CloseSend: %v", err)
	}
	if !waitFor(2*time.Second, func() bool { return ms.kcLive.Load() == 0 }) {
		t.Fatalf("handler did not release even after CloseSend; live=%d", ms.kcLive.Load())
	}
	_ = io.EOF
}

// TestRedConnDropReleasesHandler proves the OTHER release path: when the whole
// conn drops, the per-stream Context cancels and the handler returns. This is the
// path that DOES work (conn-level drop) — contrasted with the abandon-on-alive-conn
// leak above.
func TestRedConnDropReleasesHandler(t *testing.T) {
	cli, ms, conn, done := serveLeak(t)
	defer done()

	stream, err := cli.SendHeartbeat(context.Background())
	if err != nil {
		t.Fatalf("SendHeartbeat: %v", err)
	}
	if err := stream.Send(&master_pb.Heartbeat{Id: "vs-1"}); err != nil {
		t.Fatalf("Send: %v", err)
	}
	if !waitFor(time.Second, func() bool { return ms.hbLive.Load() == 1 }) {
		t.Fatalf("handler never started")
	}
	conn.Close() // hard conn drop
	if !waitFor(2*time.Second, func() bool { return ms.hbLive.Load() == 0 }) {
		t.Fatalf("handler did not release on conn drop; live=%d", ms.hbLive.Load())
	}
}
