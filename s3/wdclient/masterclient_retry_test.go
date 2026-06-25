package wdclient

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hanzoai/s3/s3/masterzap"
	"github.com/hanzoai/s3/s3/pb"
	"github.com/hanzoai/s3/s3/pb/master_pb"
	masterwire "github.com/hanzoai/s3/s3/wire/master"
	masterstream "github.com/hanzoai/s3/s3/wire/master/masterstream"
	"github.com/zap-proto/go/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// fakeLookupServer returns Unavailable for the first N calls, then succeeds.
type fakeLookupServer struct {
	master_pb.UnimplementedHanzoServer
	unavailableCount int32
	callCount        atomic.Int32
}

func (s *fakeLookupServer) LookupVolume(_ context.Context, req *master_pb.LookupVolumeRequest) (*master_pb.LookupVolumeResponse, error) {
	n := s.callCount.Add(1)
	if n <= s.unavailableCount {
		return nil, status.Errorf(codes.Unavailable, "master is warming up")
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
	return resp, nil
}

// startFakeMasterServer stands up the fake master over the native ZAP transport
// (the master service is ZAP-only now — see command/master.go), then returns a
// master ServerAddress whose ToMasterZapAddress resolves to the live listener.
// The master client derives the ZAP port as grpcPort+10000, so the address
// encodes grpcPort = (live ZAP port - 10000); ToMasterZapAddress adds 10000 back
// to reach the listener.
func startFakeMasterServer(t *testing.T, srv master_pb.HanzoServer) pb.ServerAddress {
	t.Helper()
	zapSrv, err := transport.ListenStream("tcp", "127.0.0.1:0",
		masterwire.Dispatch(masterzap.NewServerBackend(srv)),
		masterstream.Handler(masterzap.NewStreamServer(srv)))
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { _ = zapSrv.Close() })

	_, zapPortStr, _ := net.SplitHostPort(zapSrv.Addr().String())
	var zapPort int
	fmt.Sscanf(zapPortStr, "%d", &zapPort)
	grpcPort := zapPort - 10000
	return pb.ServerAddress(fmt.Sprintf("127.0.0.1:0.%d", grpcPort))
}

func TestLookupVolumeIdsRetriesOnUnavailable(t *testing.T) {
	srv := &fakeLookupServer{unavailableCount: 3}
	addr := startFakeMasterServer(t, srv)

	mc := NewMasterClient(
		grpc.WithTransportCredentials(insecure.NewCredentials()),
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
	if calls := srv.callCount.Load(); calls != 4 {
		t.Errorf("expected 4 calls (3 unavailable + 1 success), got %d", calls)
	}
}

func TestLookupVolumeIdsStopsOnContextCancel(t *testing.T) {
	srv := &fakeLookupServer{unavailableCount: 1000}
	addr := startFakeMasterServer(t, srv)

	mc := NewMasterClient(
		grpc.WithTransportCredentials(insecure.NewCredentials()),
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
	if calls := srv.callCount.Load(); calls <= 1 {
		t.Errorf("expected multiple retry attempts, got %d calls", calls)
	}
	if elapsed > 5*time.Second {
		t.Errorf("took %v, expected to stop near context deadline of 2s", elapsed)
	}
}
