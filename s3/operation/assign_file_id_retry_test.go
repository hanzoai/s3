package operation

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hanzoai/s3/s3/masterzap"
	"github.com/hanzoai/s3/s3/pb"
	"github.com/hanzoai/s3/s3/pb/master_pb"
	masterwire "github.com/hanzoai/s3/s3/wire/master"

	"github.com/zap-proto/go/transport"
)

// fakeAssignBackend is a masterwire.Backend whose Assign returns Unavailable for
// the first unavailableCount calls, then succeeds — the same warmup behaviour the
// production master signals as a plain "Unavailable: ..." error over ZAP (no gRPC
// status). The retry loop classifies it by strings.Contains(err, "Unavailable").
// Every other method is unimplemented via the embedded interface.
type fakeAssignBackend struct {
	masterwire.Backend
	unavailableCount int32
	callCount        atomic.Int32
}

func (s *fakeAssignBackend) Assign(req []byte) ([]byte, error) {
	if _, err := masterzap.AssignReqFromWire(req); err != nil {
		return nil, err
	}
	n := s.callCount.Add(1)
	if n <= s.unavailableCount {
		return nil, fmt.Errorf("Unavailable: master is warming up")
	}
	return masterzap.AssignRespToWire(&master_pb.AssignResponse{
		Fid:   "1,abc",
		Count: 1,
		Location: &master_pb.Location{
			Url:       "127.0.0.1:8080",
			PublicUrl: "127.0.0.1:8080",
		},
	}), nil
}

// startFakeMaster serves backend over the native ZAP transport and returns an
// HTTP-style master address whose ServerAddress.ToMasterZapAddress() (grpcPort+
// 10000, i.e. http+20000) resolves back to the listener, so Assign reaches it via
// pb.WithMasterClient.
func startFakeMaster(t *testing.T, backend masterwire.Backend) pb.ServerAddress {
	t.Helper()
	srv, err := transport.ListenStream("tcp", "127.0.0.1:0", masterwire.Dispatch(backend), nil)
	if err != nil {
		t.Fatalf("listen master ZAP: %v", err)
	}
	t.Cleanup(func() { _ = srv.Close() })

	host, portStr, err := net.SplitHostPort(srv.Addr().String())
	if err != nil {
		t.Fatalf("split master addr: %v", err)
	}
	zapPort, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatalf("parse master port: %v", err)
	}
	return pb.ServerAddress(net.JoinHostPort(host, strconv.Itoa(zapPort-20000)))
}

func TestAssignRetriesOnUnavailable(t *testing.T) {
	backend := &fakeAssignBackend{unavailableCount: 3}
	addr := startFakeMaster(t, backend)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ret, err := Assign(ctx, func(_ context.Context) pb.ServerAddress {
		return addr
	}, pb.DialOption{}, &VolumeAssignRequest{Count: 1})

	if err != nil {
		t.Fatalf("expected success after retries, got error: %v", err)
	}
	if ret.Fid != "1,abc" {
		t.Errorf("expected fid '1,abc', got '%s'", ret.Fid)
	}
	if calls := backend.callCount.Load(); calls != 4 {
		t.Errorf("expected 4 calls (3 unavailable + 1 success), got %d", calls)
	}
}

func TestAssignStopsOnContextCancel(t *testing.T) {
	backend := &fakeAssignBackend{unavailableCount: 1000} // never succeeds
	addr := startFakeMaster(t, backend)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	start := time.Now()
	_, err := Assign(ctx, func(_ context.Context) pb.ServerAddress {
		return addr
	}, pb.DialOption{}, &VolumeAssignRequest{Count: 1})

	elapsed := time.Since(start)
	if err == nil {
		t.Fatal("expected error from context cancellation")
	}
	// Should stop within a reasonable time after context deadline
	if elapsed > 5*time.Second {
		t.Errorf("took %v, expected to stop near context deadline of 2s", elapsed)
	}
	// Verify the loop actually retried (not just an immediate failure)
	if calls := backend.callCount.Load(); calls <= 1 {
		t.Errorf("expected multiple retry attempts, got %d calls", calls)
	}
	// Verify the error is from context deadline
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected context.DeadlineExceeded, got: %v", err)
	}
}
