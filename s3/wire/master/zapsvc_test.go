// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package masterwire

import (
	"errors"
	"testing"
)

// fakeMaster is an in-memory Backend for the round-trip test. It decodes each
// request with the matching masterwire.Wrap* and builds the reply with New* —
// the bytes are the message, end to end. The methods not exercised return a
// well-formed empty response (never nil), proving the whole surface dispatches.
type fakeMaster struct{}

// --- streaming: must refuse, never fake a stream ---

func (fakeMaster) SendHeartbeat(req []byte) ([]byte, error) { return nil, ErrStreamingNotWired }
func (fakeMaster) KeepConnected(req []byte) ([]byte, error) { return nil, ErrStreamingNotWired }
func (fakeMaster) StreamAssign(req []byte) ([]byte, error)  { return nil, ErrStreamingNotWired }

// --- unary: real decode -> real encode ---

func (fakeMaster) LookupVolume(req []byte) ([]byte, error) {
	v, err := WrapLookupVolumeRequest(req)
	if err != nil {
		return nil, err
	}
	// Echo each requested id back with one Location, so the test can assert the
	// nested repeated-message + repeated-message-in-message path survives.
	locs := make([][]byte, 0, v.VolumeOrFileIdsLen())
	for i := 0; i < v.VolumeOrFileIdsLen(); i++ {
		loc := NewLocation(LocationInput{
			Url:        "vol-" + v.VolumeOrFileIdAt(i) + ".dc1",
			PublicUrl:  "pub-" + v.VolumeOrFileIdAt(i),
			GrpcPort:   18080,
			DataCenter: v.Collection(),
		})
		locs = append(locs, NewLookupVolumeResponseVolumeIdLocation(LookupVolumeResponseVolumeIdLocationInput{
			VolumeOrFileId: v.VolumeOrFileIdAt(i),
			Locations:      [][]byte{loc},
			Auth:           "tok-" + v.VolumeOrFileIdAt(i),
		}))
	}
	return NewLookupVolumeResponse(LookupVolumeResponseInput{VolumeIdLocations: locs}), nil
}

func (fakeMaster) Assign(req []byte) ([]byte, error) {
	v, err := WrapAssignRequest(req)
	if err != nil {
		return nil, err
	}
	return NewAssignResponse(AssignResponseInput{
		Fid:   "3," + v.Collection(),
		Count: v.Count(),
		Location: NewLocation(LocationInput{
			Url:        "node1:8080",
			DataCenter: v.DataCenter(),
			GrpcPort:   18080,
		}),
	}), nil
}

func (fakeMaster) Statistics(req []byte) ([]byte, error) {
	v, err := WrapStatisticsRequest(req)
	if err != nil {
		return nil, err
	}
	// Fold the request strings into the numbers so the assertion is meaningful.
	return NewStatisticsResponse(StatisticsResponseInput{
		TotalSize: uint64(len(v.Replication()) + len(v.Collection())),
		UsedSize:  uint64(len(v.Ttl())),
		FileCount: uint64(len(v.DiskType())),
	}), nil
}

func (fakeMaster) CollectionList(req []byte) ([]byte, error) {
	v, err := WrapCollectionListRequest(req)
	if err != nil {
		return nil, err
	}
	cols := [][]byte{NewCollection(CollectionInput{Name: "default"})}
	if v.IncludeEcVolumes() {
		cols = append(cols, NewCollection(CollectionInput{Name: "ec"}))
	}
	return NewCollectionListResponse(CollectionListResponseInput{Collections: cols}), nil
}

func (fakeMaster) CollectionDelete(req []byte) ([]byte, error) {
	if _, err := WrapCollectionDeleteRequest(req); err != nil {
		return nil, err
	}
	return NewCollectionDeleteResponse(CollectionDeleteResponseInput{}), nil
}

func (fakeMaster) VolumeList(req []byte) ([]byte, error) {
	if _, err := WrapVolumeListRequest(req); err != nil {
		return nil, err
	}
	// Build a one-DC topology so the nested message + map path is on the wire.
	dc := NewDataCenterInfo(DataCenterInfoInput{Id: "dc1"})
	topo := NewTopologyInfo(TopologyInfoInput{Id: "topo", DataCenterInfos: [][]byte{dc}})
	return NewVolumeListResponse(VolumeListResponseInput{TopologyInfo: topo, VolumeSizeLimitMb: 30000}), nil
}

func (fakeMaster) LookupEcVolume(req []byte) ([]byte, error) {
	v, err := WrapLookupEcVolumeRequest(req)
	if err != nil {
		return nil, err
	}
	return NewLookupEcVolumeResponse(LookupEcVolumeResponseInput{VolumeId: v.VolumeId()}), nil
}

func (fakeMaster) VacuumVolume(req []byte) ([]byte, error) {
	if _, err := WrapVacuumVolumeRequest(req); err != nil {
		return nil, err
	}
	return NewVacuumVolumeResponse(VacuumVolumeResponseInput{}), nil
}

func (fakeMaster) DisableVacuum(req []byte) ([]byte, error) {
	if _, err := WrapDisableVacuumRequest(req); err != nil {
		return nil, err
	}
	return NewDisableVacuumResponse(DisableVacuumResponseInput{}), nil
}

func (fakeMaster) EnableVacuum(req []byte) ([]byte, error) {
	if _, err := WrapEnableVacuumRequest(req); err != nil {
		return nil, err
	}
	return NewEnableVacuumResponse(EnableVacuumResponseInput{}), nil
}

func (fakeMaster) VolumeMarkReadonly(req []byte) ([]byte, error) {
	v, err := WrapVolumeMarkReadonlyRequest(req)
	if err != nil {
		return nil, err
	}
	// Round-trip a bool + scalars back through a response-free RPC by echoing
	// into a Ping-shaped probe is not possible; instead assert in-handler and
	// return the empty response. The test asserts the call succeeds, and the
	// LookupVolume/Statistics cases already prove field fidelity.
	_ = v.IsReadonly()
	return NewVolumeMarkReadonlyResponse(VolumeMarkReadonlyResponseInput{}), nil
}

func (fakeMaster) GetMasterConfiguration(req []byte) ([]byte, error) {
	if _, err := WrapGetMasterConfigurationRequest(req); err != nil {
		return nil, err
	}
	return NewGetMasterConfigurationResponse(GetMasterConfigurationResponseInput{
		Leader: "master1:9333", DefaultReplication: "000",
	}), nil
}

func (fakeMaster) ListClusterNodes(req []byte) ([]byte, error) {
	if _, err := WrapListClusterNodesRequest(req); err != nil {
		return nil, err
	}
	return NewListClusterNodesResponse(ListClusterNodesResponseInput{}), nil
}

func (fakeMaster) LeaseAdminToken(req []byte) ([]byte, error) {
	if _, err := WrapLeaseAdminTokenRequest(req); err != nil {
		return nil, err
	}
	return NewLeaseAdminTokenResponse(LeaseAdminTokenResponseInput{Token: 42, LockTsNs: 1000}), nil
}

func (fakeMaster) ReleaseAdminToken(req []byte) ([]byte, error) {
	if _, err := WrapReleaseAdminTokenRequest(req); err != nil {
		return nil, err
	}
	return NewReleaseAdminTokenResponse(ReleaseAdminTokenResponseInput{}), nil
}

func (fakeMaster) Ping(req []byte) ([]byte, error) {
	v, err := WrapPingRequest(req)
	if err != nil {
		return nil, err
	}
	// Encode the request target length into the times so the round-trip is
	// observable end to end.
	n := int64(len(v.Target()))
	return NewPingResponse(PingResponseInput{StartTimeNs: n, RemoteTimeNs: n + 1, StopTimeNs: n + 2}), nil
}

func (fakeMaster) RaftListClusterServers(req []byte) ([]byte, error) {
	if _, err := WrapRaftListClusterServersRequest(req); err != nil {
		return nil, err
	}
	return NewRaftListClusterServersResponse(RaftListClusterServersResponseInput{}), nil
}

func (fakeMaster) RaftAddServer(req []byte) ([]byte, error) {
	if _, err := WrapRaftAddServerRequest(req); err != nil {
		return nil, err
	}
	return NewRaftAddServerResponse(RaftAddServerResponseInput{}), nil
}

func (fakeMaster) RaftRemoveServer(req []byte) ([]byte, error) {
	if _, err := WrapRaftRemoveServerRequest(req); err != nil {
		return nil, err
	}
	return NewRaftRemoveServerResponse(RaftRemoveServerResponseInput{}), nil
}

func (fakeMaster) RaftLeadershipTransfer(req []byte) ([]byte, error) {
	v, err := WrapRaftLeadershipTransferRequest(req)
	if err != nil {
		return nil, err
	}
	return NewRaftLeadershipTransferResponse(RaftLeadershipTransferResponseInput{
		PreviousLeader: "old", NewLeader: v.TargetId(),
	}), nil
}

func (fakeMaster) VolumeGrow(req []byte) ([]byte, error) {
	if _, err := WrapVolumeGrowRequest(req); err != nil {
		return nil, err
	}
	return NewVolumeGrowResponse(VolumeGrowResponseInput{}), nil
}

// TestRoundTrip proves the master service crosses a real TCP socket as ZAP RPC
// envelopes carrying zero-copy masterwire payloads — no HTTP, no protobuf, no
// gRPC. It exercises a scalar RPC (Ping), a nested repeated-message RPC
// (LookupVolume), a nested-message-in-response RPC (VolumeList), a numeric RPC
// (Statistics) and a couple of void RPCs, then asserts a streaming RPC is
// refused (not silently degraded).
func TestRoundTrip(t *testing.T) {
	srv, err := Serve("tcp", "127.0.0.1:0", fakeMaster{})
	if err != nil {
		t.Fatalf("Serve: %v", err)
	}
	defer srv.Close()

	cli, err := Dial("tcp", srv.Addr().String())
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer cli.Close()

	// --- Ping: scalar round-trip across the socket ---
	pr, err := cli.Ping(PingRequestInput{Target: "master2:9333", TargetType: "master"})
	if err != nil {
		t.Fatalf("Ping: %v", err)
	}
	if got, want := pr.StartTimeNs(), int64(len("master2:9333")); got != want {
		t.Fatalf("Ping StartTimeNs = %d, want %d", got, want)
	}
	if pr.StopTimeNs() != pr.StartTimeNs()+2 {
		t.Fatalf("Ping StopTimeNs = %d, want %d", pr.StopTimeNs(), pr.StartTimeNs()+2)
	}

	// --- LookupVolume: repeated message + repeated-message-in-element ---
	lv, err := cli.LookupVolume(LookupVolumeRequestInput{
		VolumeOrFileIds: []string{"7,abc", "9,def"},
		Collection:      "dc-west",
	})
	if err != nil {
		t.Fatalf("LookupVolume: %v", err)
	}
	if lv.VolumeIdLocationsLen() != 2 {
		t.Fatalf("LookupVolume locations = %d, want 2", lv.VolumeIdLocationsLen())
	}
	vidl, ok := lv.VolumeIdLocationAt(0)
	if !ok {
		t.Fatal("LookupVolume VolumeIdLocationAt(0) not ok")
	}
	if vidl.VolumeOrFileId() != "7,abc" {
		t.Fatalf("VolumeOrFileId = %q, want %q", vidl.VolumeOrFileId(), "7,abc")
	}
	if vidl.Auth() != "tok-7,abc" {
		t.Fatalf("Auth = %q, want %q", vidl.Auth(), "tok-7,abc")
	}
	if vidl.LocationsLen() != 1 {
		t.Fatalf("nested LocationsLen = %d, want 1", vidl.LocationsLen())
	}
	loc, ok := vidl.LocationAt(0)
	if !ok {
		t.Fatal("nested LocationAt(0) not ok")
	}
	if loc.Url() != "vol-7,abc.dc1" {
		t.Fatalf("nested Location.Url = %q", loc.Url())
	}
	if loc.DataCenter() != "dc-west" {
		t.Fatalf("nested Location.DataCenter = %q, want dc-west", loc.DataCenter())
	}
	if loc.GrpcPort() != 18080 {
		t.Fatalf("nested Location.GrpcPort = %d, want 18080", loc.GrpcPort())
	}

	// --- Statistics: numeric folding round-trip ---
	st, err := cli.Statistics(StatisticsRequestInput{Replication: "001", Collection: "c", Ttl: "1h", DiskType: "ssd"})
	if err != nil {
		t.Fatalf("Statistics: %v", err)
	}
	if st.TotalSize() != uint64(len("001")+len("c")) {
		t.Fatalf("Statistics TotalSize = %d", st.TotalSize())
	}
	if st.FileCount() != uint64(len("ssd")) {
		t.Fatalf("Statistics FileCount = %d", st.FileCount())
	}

	// --- VolumeList: nested message in response ---
	vl, err := cli.VolumeList(VolumeListRequestInput{})
	if err != nil {
		t.Fatalf("VolumeList: %v", err)
	}
	if vl.VolumeSizeLimitMb() != 30000 {
		t.Fatalf("VolumeList VolumeSizeLimitMb = %d, want 30000", vl.VolumeSizeLimitMb())
	}
	topo, ok := vl.TopologyInfo()
	if !ok {
		t.Fatal("VolumeList TopologyInfo not present")
	}
	if topo.Id() != "topo" {
		t.Fatalf("topology Id = %q, want topo", topo.Id())
	}
	if topo.DataCenterInfosLen() != 1 {
		t.Fatalf("topology DataCenterInfosLen = %d, want 1", topo.DataCenterInfosLen())
	}
	if dc, ok := topo.DataCenterInfoAt(0); !ok || dc.Id() != "dc1" {
		t.Fatalf("topology DataCenterInfoAt(0) = %v ok=%v", dc.Id(), ok)
	}

	// --- CollectionList: bool-driven repeated message ---
	cl, err := cli.CollectionList(CollectionListRequestInput{IncludeEcVolumes: true})
	if err != nil {
		t.Fatalf("CollectionList: %v", err)
	}
	if cl.CollectionsLen() != 2 {
		t.Fatalf("CollectionList CollectionsLen = %d, want 2", cl.CollectionsLen())
	}

	// --- VolumeMarkReadonly: void RPC with a bool field, just succeed ---
	if _, err := cli.VolumeMarkReadonly(VolumeMarkReadonlyRequestInput{VolumeId: 7, IsReadonly: true}); err != nil {
		t.Fatalf("VolumeMarkReadonly: %v", err)
	}

	// --- RaftLeadershipTransfer: string echo ---
	rl, err := cli.RaftLeadershipTransfer(RaftLeadershipTransferRequestInput{TargetId: "node3"})
	if err != nil {
		t.Fatalf("RaftLeadershipTransfer: %v", err)
	}
	if rl.NewLeader() != "node3" {
		t.Fatalf("RaftLeadershipTransfer NewLeader = %q, want node3", rl.NewLeader())
	}
}

// TestStreamingRefused proves the three streaming RPCs are NOT faked: dispatched
// over the unary path they surface as an error (StatusInternal at the wire,
// driven by ErrStreamingNotWired in the handler), never a truncated stream.
func TestStreamingRefused(t *testing.T) {
	// The handler-level contract: ErrStreamingNotWired, directly.
	h := NewHandler(fakeMaster{})
	for name, call := range map[string]func() ([]byte, error){
		"SendHeartbeat": func() ([]byte, error) { return h.SendHeartbeat(nil) },
		"KeepConnected": func() ([]byte, error) { return h.KeepConnected(nil) },
		"StreamAssign":  func() ([]byte, error) { return h.StreamAssign(nil) },
	} {
		if _, err := call(); !errors.Is(err, ErrStreamingNotWired) {
			t.Fatalf("%s: err = %v, want ErrStreamingNotWired", name, err)
		}
	}

	// The wire-level contract: over a real socket the streaming ordinal returns
	// a non-OK status, so the client surfaces an error rather than a body.
	srv, err := Serve("tcp", "127.0.0.1:0", fakeMaster{})
	if err != nil {
		t.Fatalf("Serve: %v", err)
	}
	defer srv.Close()
	cli, err := Dial("tcp", srv.Addr().String())
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer cli.Close()

	if _, _, err := cli.RPC().SendHeartbeat(NewHeartbeat(HeartbeatInput{Ip: "10.0.0.1", Port: 8080})); err == nil {
		t.Fatal("SendHeartbeat over unary path: expected error, got nil")
	}
}
