package wdclient

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

// fakeLookupBackend is a masterwire.Backend whose LookupVolume returns
// Unavailable for the first unavailableCount calls, then succeeds — the warmup
// behaviour the production master signals as a plain "Unavailable: ..." error
// over ZAP (no gRPC status). The retry loop classifies it by
// strings.Contains(err, "Unavailable"). Other methods are unimplemented via the
// embedded interface.
type fakeLookupBackend struct {
	masterwire.Backend
	unavailableCount int32
	callCount        atomic.Int32
}

func (s *fakeLookupBackend) LookupVolume(reqBytes []byte) ([]byte, error) {
	req, err := masterzap.LookupVolumeReqFromWire(reqBytes)
	if err != nil {
		return nil, err
	}
	n := s.callCount.Add(1)
	if n <= s.unavailableCount {
		return nil, fmt.Errorf("Unavailable: master is warming up")
	}
	resp := &master_pb.LookupVolumeResponse{}
	for _, vid := range req.VolumeOrFileIds {
		resp.VolumeIdLocations = append(resp.VolumeIdLocations, &master_pb.LookupVolumeResponse_VolumeIdLocation{
			VolumeOrFileId: vid,
			Locations: []*master_pb.Location{
				{Url: "127.0.0.1:8080", PublicUrl: "127.0.0.1:8080"},
			},
		})
	}
	return masterzap.LookupVolumeRespToWire(resp), nil
}

// startFakeMasterServer serves backend over the native ZAP transport and returns
// an HTTP-style master address whose ServerAddress.ToMasterZapAddress() (grpcPort
// +10000, i.e. http+20000) resolves back to the listener, so the master client
// reaches it via pb.WithMasterClient.
func startFakeMasterServer(t *testing.T, backend masterwire.Backend) pb.ServerAddress {
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

func TestLookupVolumeIdsRetriesOnUnavailable(t *testing.T) {
	backend := &fakeLookupBackend{unavailableCount: 3}
	addr := startFakeMasterServer(t, backend)

	mc := NewMasterClient(
		pb.DialOption{},
		"", "test", "", "", "",
		pb.ServerDiscovery{},
	)
	mc.setCurrentMaster(addr)
	mc.grpcTimeout = 5 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	provider := &masterVolumeProvider{masterClient: mc}
	result, err := provider.LookupVolumeIds(ctx, []string{"1"})

	if err != nil {
		t.Fatalf("expected success after retries, got error: %v", err)
	}
	if _, ok := result["1"]; !ok {
		t.Error("expected volume 1 in result")
	}
	if calls := backend.callCount.Load(); calls != 4 {
		t.Errorf("expected 4 calls (3 unavailable + 1 success), got %d", calls)
	}
}

func TestLookupVolumeIdsStopsOnContextCancel(t *testing.T) {
	backend := &fakeLookupBackend{unavailableCount: 1000}
	addr := startFakeMasterServer(t, backend)

	mc := NewMasterClient(
		pb.DialOption{},
		"", "test", "", "", "",
		pb.ServerDiscovery{},
	)
	mc.setCurrentMaster(addr)
	mc.grpcTimeout = 5 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	provider := &masterVolumeProvider{masterClient: mc}
	start := time.Now()
	_, err := provider.LookupVolumeIds(ctx, []string{"1"})
	elapsed := time.Since(start)

	// Verify the error is from context deadline
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context.DeadlineExceeded, got: %v", err)
	}
	// Verify the loop actually retried (not just an immediate failure)
	if calls := backend.callCount.Load(); calls <= 1 {
		t.Errorf("expected multiple retry attempts, got %d calls", calls)
	}
	if elapsed > 5*time.Second {
		t.Errorf("took %v, expected to stop near context deadline of 2s", elapsed)
	}
}
