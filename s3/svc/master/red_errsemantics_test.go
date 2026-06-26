// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// RED adversarial probe: gRPC status Code does not cross the ZAP wire, so any
// master caller that classifies retryability via status.FromError(...).Code()
// (rather than a string match) loses its retry/failover. This reproduces the
// exact predicate s3/operation/assign_file_id.go uses against a master that
// returns codes.Unavailable ("warming up"), proving it no longer retries.

package master

import (
	"context"
	"fmt"
	"strings"
	"testing"

	master_pb "github.com/hanzoai/s3/s3/pb/master_pb"
	masterwire "github.com/hanzoai/s3/s3/wire/master"
	masterstream "github.com/hanzoai/s3/s3/wire/master/masterstream"
	"github.com/zap-proto/go/transport"
)

// warmingMaster returns the exact warmup error the live master emits on Assign
// while topology loads (server/master_grpc_server_assign.go:72).
type warmingMaster struct {
	master_pb.UnimplementedHanzoServer
}

func (warmingMaster) Assign(_ context.Context, _ *master_pb.AssignRequest) (*master_pb.AssignResponse, error) {
	return nil, fmt.Errorf("Unavailable: master is warming up, topology is still loading")
}

// assignRetryPredicate is copied verbatim from s3/operation/assign_file_id.go
// (the retryable classifier). If it returns false for a master-warmup error,
// the assign hot-path fails fast instead of retrying for its 30s budget.
func assignRetryPredicate(err error, ctxLive bool) bool {
	msg := err.Error()
	if strings.Contains(msg, "Unavailable") {
		return true
	}
	if strings.Contains(msg, "Canceled") ||
		strings.Contains(msg, "DeadlineExceeded") {
		return ctxLive
	}
	return false
}

func TestRedAssignUnavailableNotRetryableOverZap(t *testing.T) {
	ms := warmingMaster{}
	srv, err := transport.ListenStream("tcp", "127.0.0.1:0",
		masterwire.Dispatch(NewServerBackend(ms)),
		masterstream.Handler(NewStreamServer(ms)))
	if err != nil {
		t.Fatalf("ListenStream: %v", err)
	}
	defer srv.Close()
	conn, err := transport.Dial("tcp", srv.Addr().String())
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer conn.Close()

	cli := New(conn, nil)
	_, gotErr := cli.Assign(context.Background(), &master_pb.AssignRequest{Count: 1})
	if gotErr == nil {
		t.Fatal("expected an error from warming master")
	}
	t.Logf("ZAP-transported master error: %q", gotErr.Error())

	// The gRPC status CODE does NOT cross the ZAP wire — only the desc text
	// does, and the FIXED classifier matches on it. Guard that a master-warmup
	// error transported over ZAP is still classified retryable so the assign hot
	// path rides out its 30s budget during master restart/warmup.
	if !strings.Contains(gotErr.Error(), "Unavailable") {
		t.Fatalf("expected %q in the wire error, got %q", "Unavailable", gotErr.Error())
	}
	if !assignRetryPredicate(gotErr, true) {
		t.Fatalf("FIX REGRESSED: ZAP-transported master-warmup error not classified retryable — assign hot path would fail fast during warmup")
	}
}
