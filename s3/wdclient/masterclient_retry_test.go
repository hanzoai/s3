package wdclient

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hanzoai/s3/s3/pb"
	"github.com/hanzoai/s3/s3/pb/master_pb"
	"github.com/hanzoai/s3/s3/svc/master"
	masterwire "github.com/hanzoai/s3/s3/wire/master"
	masterstream "github.com/hanzoai/s3/s3/wire/master/masterstream"
	"github.com/zap-proto/go/transport"
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
		masterwire.Dispatch(master.NewServerBackend(srv)),
		masterstream.Handler(master.NewStreamServer(srv)))
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
	if calls := srv.callCount.Load(); calls != 4 {
		t.Errorf("expected 4 calls (3 unavailable + 1 success), got %d", calls)
	}
}

func TestLookupVolumeIdsStopsOnContextCancel(t *testing.T) {
	srv := &fakeLookupServer{unavailableCount: 1000}
	addr := startFakeMasterServer(t, srv)

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
	if calls := srv.callCount.Load(); calls <= 1 {
		t.Errorf("expected multiple retry attempts, got %d calls", calls)
	}
	if elapsed > 5*time.Second {
		t.Errorf("took %v, expected to stop near context deadline of 2s", elapsed)
	}
}

// warmingMasterServer knows no volumes for its first coldCalls lookups, which is
// what a master does before volume-server heartbeats land. Then it resolves.
type warmingMasterServer struct {
	master_pb.UnimplementedHanzoServer
	coldCalls int32
	callCount atomic.Int32
}

func (s *warmingMasterServer) LookupVolume(_ context.Context, req *master_pb.LookupVolumeRequest) (*master_pb.LookupVolumeResponse, error) {
	n := s.callCount.Add(1)
	resp := &master_pb.LookupVolumeResponse{}
	for _, vid := range req.VolumeOrFileIds {
		loc := &master_pb.LookupVolumeResponse_VolumeIdLocation{VolumeOrFileId: vid}
		if n <= s.coldCalls {
			// A successful RPC carrying a per-volume error — not an RPC failure,
			// which is why the retry never used to see it.
			loc.Error = "volume id " + vid + " not found"
		} else {
			loc.Locations = []*master_pb.Location{{Url: "127.0.0.1:8080", PublicUrl: "127.0.0.1:8080"}}
		}
		resp.VolumeIdLocations = append(resp.VolumeIdLocations, loc)
	}
	return resp, nil
}

func newTestMasterClient(t *testing.T, addr pb.ServerAddress) *MasterClient {
	t.Helper()
	mc := NewMasterClient(pb.DialOption{}, "", "test", "", "", "", pb.ServerDiscovery{})
	mc.setCurrentMaster(addr)
	mc.grpcTimeout = 5 * time.Second
	return mc
}

// A master with no topology is retried, not believed. This is the 2026-08-15
// outage: hanzo.ai's index.html lookup hit an empty master, got "volume id 576
// not found", and the edge served 404 while the object was intact.
func TestAColdMasterIsWaitedForRatherThanBelieved(t *testing.T) {
	srv := &warmingMasterServer{coldCalls: 2}
	mc := newTestMasterClient(t, startFakeMasterServer(t, srv))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := (&masterVolumeProvider{masterClient: mc}).LookupVolumeIds(ctx, []string{"576"})
	if err != nil {
		t.Fatalf("expected the warm master's answer to win, got: %v", err)
	}
	if _, ok := result["576"]; !ok {
		t.Fatalf("volume 576 missing; a cold 'not found' was taken as authoritative (result=%v)", result)
	}
	if n := srv.callCount.Load(); n < 3 {
		t.Errorf("expected the lookup to be retried past the cold answers, got %d calls", n)
	}
}

// A volume that is really gone reports the master's own error, not the sentinel.
func TestADeletedVolumeStillReportsNotFound(t *testing.T) {
	// Never warms — indistinguishable from a deleted volume. Bounded wait, then
	// the real error.
	srv := &warmingMasterServer{coldCalls: 1 << 30}
	mc := newTestMasterClient(t, startFakeMasterServer(t, srv))

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	result, err := (&masterVolumeProvider{masterClient: mc}).LookupVolumeIds(ctx, []string{"999"})
	if err == nil {
		t.Fatal("expected an error for a volume no master knows")
	}
	if len(result) != 0 {
		t.Errorf("expected no locations, got %v", result)
	}
	if errors.Is(err, errTopologyNotLoaded) {
		t.Errorf("reported the WAIT as the reason; caller needs the master's own answer: %v", err)
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected the master's per-volume error to survive, got: %v", err)
	}
}

// Resolving some means the topology is loaded, so a missing one costs no wait.
func TestAPartialAnswerIsNotTreatedAsColdness(t *testing.T) {
	srv := &partialMasterServer{}
	mc := newTestMasterClient(t, startFakeMasterServer(t, srv))

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	start := time.Now()
	result, err := (&masterVolumeProvider{masterClient: mc}).LookupVolumeIds(ctx, []string{"1", "999"})
	elapsed := time.Since(start)

	if _, ok := result["1"]; !ok {
		t.Errorf("expected the known volume to resolve, got %v", result)
	}
	if err == nil {
		t.Error("expected the unknown volume to be reported")
	}
	if elapsed > masterWarmup {
		t.Errorf("waited %v on an authoritative answer; a partial result means the topology IS loaded", elapsed)
	}
	if n := srv.callCount.Load(); n != 1 {
		t.Errorf("expected exactly 1 call for an authoritative answer, got %d", n)
	}
}

// partialMasterServer knows volume 1 only: a loaded topology.
type partialMasterServer struct {
	master_pb.UnimplementedHanzoServer
	callCount atomic.Int32
}

func (s *partialMasterServer) LookupVolume(_ context.Context, req *master_pb.LookupVolumeRequest) (*master_pb.LookupVolumeResponse, error) {
	s.callCount.Add(1)
	resp := &master_pb.LookupVolumeResponse{}
	for _, vid := range req.VolumeOrFileIds {
		loc := &master_pb.LookupVolumeResponse_VolumeIdLocation{VolumeOrFileId: vid}
		if vid == "1" {
			loc.Locations = []*master_pb.Location{{Url: "127.0.0.1:8080", PublicUrl: "127.0.0.1:8080"}}
		} else {
			loc.Error = "volume id " + vid + " not found"
		}
		resp.VolumeIdLocations = append(resp.VolumeIdLocations, loc)
	}
	return resp, nil
}

// masterWarmup must actually apply. Spelling errTopologyNotLoaded with an
// "Unavailable" prefix made isUnavailableErr match first, handing it the 30s
// budget — visible only as a lookup that took too long.
func TestTheWarmupBoundIsShorterThanTheTransportBudget(t *testing.T) {
	if isUnavailableErr(errTopologyNotLoaded) {
		t.Fatal("errTopologyNotLoaded is classified as a transient transport error, " +
			"so it inherits the 30s budget and masterWarmup never applies")
	}
	srv := &warmingMasterServer{coldCalls: 1 << 30} // never warms
	mc := newTestMasterClient(t, startFakeMasterServer(t, srv))

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	start := time.Now()
	_, err := (&masterVolumeProvider{masterClient: mc}).LookupVolumeIds(ctx, []string{"999"})
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected an error for a volume no master ever knows")
	}
	// A backoff step may straddle the bound; just keep it off the 30s budget.
	if elapsed > masterWarmup+10*time.Second {
		t.Errorf("waited %v; masterWarmup is %v, so this fell through to the transport budget", elapsed, masterWarmup)
	}
}
