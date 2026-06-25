// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// server_roundtrip_test.go proves the SERVER adapter (NewServerBackend +
// NewStreamServer) that wraps a real master_pb.HanzoServer and serves it over
// transport.ListenStream — the exact path command/master.go and
// command/master_follower.go now use in place of the killed gRPC HanzoServer
// registration. A fake master_pb.HanzoServer (the engine) is wrapped, served,
// and hit by the real master.New client, so the full loop is exercised end to
// end: pb request -> ToWire (client) -> dispatch -> ReqFromWire (server) ->
// engine -> RespToWire (server) -> RespFromWire (client) -> pb response, for the
// representative unary RPCs and ALL 3 bidirectional streams. This is the
// server-side mirror of rpc_roundtrip_test.go (which drives a hand-written wire
// backend); together they close both directions over a real socket.

package master

import (
	"context"
	"io"
	"testing"

	master_pb "github.com/hanzoai/s3/s3/pb/master_pb"
	masterwire "github.com/hanzoai/s3/s3/wire/master"
	masterstream "github.com/hanzoai/s3/s3/wire/master/masterstream"
	"github.com/zap-proto/go/transport"
)

// fakeMaster is a master_pb.HanzoServer that records what it received and returns
// deterministic responses, so the test asserts on both directions of the loop.
// Embedding UnimplementedHanzoServer keeps it a total HanzoServer; only the
// methods the test drives are overridden.
type fakeMaster struct {
	master_pb.UnimplementedHanzoServer

	gotPingTarget    string
	gotLookupIDs     []string
	gotAssignColl    string
	gotVolumeGrowTtl string
}

func (f *fakeMaster) Ping(_ context.Context, in *master_pb.PingRequest) (*master_pb.PingResponse, error) {
	f.gotPingTarget = in.Target
	return &master_pb.PingResponse{StartTimeNs: 11, RemoteTimeNs: 22, StopTimeNs: 33}, nil
}

func (f *fakeMaster) LookupVolume(_ context.Context, in *master_pb.LookupVolumeRequest) (*master_pb.LookupVolumeResponse, error) {
	f.gotLookupIDs = in.VolumeOrFileIds
	resp := &master_pb.LookupVolumeResponse{}
	for _, id := range in.VolumeOrFileIds {
		resp.VolumeIdLocations = append(resp.VolumeIdLocations, &master_pb.LookupVolumeResponse_VolumeIdLocation{
			VolumeOrFileId: id,
			Auth:           "tok-" + id,
			Locations:      []*master_pb.Location{{Url: "vol-" + id, GrpcPort: 18080, DataCenter: "dc1"}},
		})
	}
	return resp, nil
}

func (f *fakeMaster) VolumeList(_ context.Context, _ *master_pb.VolumeListRequest) (*master_pb.VolumeListResponse, error) {
	// A deep topology graph: Topology -> DC -> Rack -> DataNode -> DiskInfo, plus a
	// map<string,DiskInfo> at the topology root carrying a volume. Proves the
	// server-side topologyInfoToWire over a real engine response.
	disk := &master_pb.DiskInfo{
		Type: "hdd", MaxVolumeCount: 100,
		VolumeInfos: []*master_pb.VolumeInformationMessage{{Id: 7, Size: 1 << 20, Collection: "c"}},
		Tags:        []string{"t1"},
	}
	node := &master_pb.DataNodeInfo{
		Id: "n1", Address: "10.0.0.1:18080", GrpcPort: 18080,
		DiskInfos: map[string]*master_pb.DiskInfo{"hdd": disk},
	}
	rack := &master_pb.RackInfo{Id: "r1", DataNodeInfos: []*master_pb.DataNodeInfo{node}}
	dc := &master_pb.DataCenterInfo{Id: "dc1", RackInfos: []*master_pb.RackInfo{rack}}
	return &master_pb.VolumeListResponse{
		VolumeSizeLimitMb: 30000,
		TopologyInfo: &master_pb.TopologyInfo{
			Id: "topo", DataCenterInfos: []*master_pb.DataCenterInfo{dc},
			DiskInfos: map[string]*master_pb.DiskInfo{"hdd": disk},
		},
	}, nil
}

func (f *fakeMaster) VolumeGrow(_ context.Context, in *master_pb.VolumeGrowRequest) (*master_pb.VolumeGrowResponse, error) {
	f.gotVolumeGrowTtl = in.Ttl
	return &master_pb.VolumeGrowResponse{}, nil
}

// StreamAssign echoes each AssignRequest as an AssignResponse keyed on the
// request's collection, until the client half-closes (io.EOF). Records the first
// non-empty collection seen.
func (f *fakeMaster) StreamAssign(stream master_pb.Hanzo_StreamAssignServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if req.Collection != "" && f.gotAssignColl == "" {
			f.gotAssignColl = req.Collection
		}
		if err := stream.Send(&master_pb.AssignResponse{
			Fid:      "fid-" + req.Collection,
			Count:    req.Count,
			Location: &master_pb.Location{Url: "assign", GrpcPort: 1},
		}); err != nil {
			return err
		}
	}
}

// SendHeartbeat answers each Heartbeat with a HeartbeatResponse carrying the
// leader, reflecting the heartbeat id into a storage-backend id so the request
// fields are observably round-tripped.
func (f *fakeMaster) SendHeartbeat(stream master_pb.Hanzo_SendHeartbeatServer) error {
	for {
		hb, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if err := stream.Send(&master_pb.HeartbeatResponse{
			Leader: "leader:9333", VolumeSizeLimit: 30000, MetricsIntervalSeconds: 15,
			StorageBackends: []*master_pb.StorageBackend{{Type: "echo", Id: hb.Id}},
		}); err != nil {
			return err
		}
	}
}

// KeepConnected answers the open subscription with one VolumeLocation update,
// then ends (so the client sees io.EOF).
func (f *fakeMaster) KeepConnected(stream master_pb.Hanzo_KeepConnectedServer) error {
	req, err := stream.Recv()
	if err != nil {
		return err
	}
	return stream.Send(&master_pb.KeepConnectedResponse{
		VolumeLocation: &master_pb.VolumeLocation{
			Url: "vs-1:18080", DataCenter: req.DataCenter, NewVids: []uint32{9},
		},
	})
}

// serveEngine stands up ONE transport listener serving the fake master_pb engine
// over the master SERVER adapter (NewServerBackend unary + NewStreamServer
// streaming), exactly as command/master.go now does, and returns the real
// master client plus a cleanup func and the fake (for received-side asserts).
func serveEngine(t *testing.T) (master_pb.HanzoClient, *fakeMaster, func()) {
	t.Helper()
	ms := &fakeMaster{}
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
	cli := New(conn, nil)
	return cli, ms, func() { conn.Close(); srv.Close() }
}

func TestServerUnaryRoundTrip(t *testing.T) {
	cli, ms, done := serveEngine(t)
	defer done()
	ctx := context.Background()

	t.Run("Ping", func(t *testing.T) {
		resp, err := cli.Ping(ctx, &master_pb.PingRequest{Target: "m1", TargetType: "master"})
		if err != nil {
			t.Fatal(err)
		}
		if ms.gotPingTarget != "m1" {
			t.Fatalf("server saw target %q, want m1", ms.gotPingTarget)
		}
		if resp.StartTimeNs != 11 || resp.RemoteTimeNs != 22 || resp.StopTimeNs != 33 {
			t.Fatalf("ping scalars lost: %+v", resp)
		}
	})

	t.Run("LookupVolume", func(t *testing.T) {
		resp, err := cli.LookupVolume(ctx, &master_pb.LookupVolumeRequest{
			VolumeOrFileIds: []string{"3", "5"}, Collection: "c",
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(ms.gotLookupIDs) != 2 || ms.gotLookupIDs[0] != "3" || ms.gotLookupIDs[1] != "5" {
			t.Fatalf("server saw ids %+v, want [3 5]", ms.gotLookupIDs)
		}
		if len(resp.VolumeIdLocations) != 2 {
			t.Fatalf("want 2 locations, got %d", len(resp.VolumeIdLocations))
		}
		got := resp.VolumeIdLocations[0]
		if got.VolumeOrFileId != "3" || got.Auth != "tok-3" {
			t.Fatalf("location scalars lost: %+v", got)
		}
		if len(got.Locations) != 1 || got.Locations[0].Url != "vol-3" ||
			got.Locations[0].GrpcPort != 18080 || got.Locations[0].DataCenter != "dc1" {
			t.Fatalf("nested location lost: %+v", got.Locations)
		}
	})

	t.Run("VolumeListTopologyGraph", func(t *testing.T) {
		resp, err := cli.VolumeList(ctx, &master_pb.VolumeListRequest{})
		if err != nil {
			t.Fatal(err)
		}
		if resp.VolumeSizeLimitMb != 30000 {
			t.Fatalf("scalar lost: %d", resp.VolumeSizeLimitMb)
		}
		ti := resp.TopologyInfo
		if ti == nil || ti.Id != "topo" {
			t.Fatalf("topology root lost: %+v", ti)
		}
		if len(ti.DataCenterInfos) != 1 || ti.DataCenterInfos[0].Id != "dc1" {
			t.Fatalf("dc lost: %+v", ti.DataCenterInfos)
		}
		node := ti.DataCenterInfos[0].RackInfos[0].DataNodeInfos
		if len(node) != 1 || node[0].Address != "10.0.0.1:18080" {
			t.Fatalf("node lost: %+v", node)
		}
		disk, ok := node[0].DiskInfos["hdd"]
		if !ok || disk.Type != "hdd" || len(disk.VolumeInfos) != 1 || disk.VolumeInfos[0].Id != 7 {
			t.Fatalf("disk map / volume lost: %+v", node[0].DiskInfos)
		}
		if len(disk.Tags) != 1 || disk.Tags[0] != "t1" {
			t.Fatalf("disk tags lost: %+v", disk.Tags)
		}
		if _, ok := ti.DiskInfos["hdd"]; !ok {
			t.Fatalf("topology-root disk map lost: %+v", ti.DiskInfos)
		}
	})

	t.Run("VolumeGrowEmptyResponse", func(t *testing.T) {
		if _, err := cli.VolumeGrow(ctx, &master_pb.VolumeGrowRequest{
			WritableVolumeCount: 3, Collection: "c", Ttl: "7d",
		}); err != nil {
			t.Fatal(err)
		}
		if ms.gotVolumeGrowTtl != "7d" {
			t.Fatalf("server saw ttl %q, want 7d", ms.gotVolumeGrowTtl)
		}
	})
}

func TestServerStreamAssignBidi(t *testing.T) {
	cli, ms, done := serveEngine(t)
	defer done()

	stream, err := cli.StreamAssign(context.Background())
	if err != nil {
		t.Fatalf("StreamAssign: %v", err)
	}
	// The empty opener is discarded server-side; the engine sees only the real
	// frames the client Sends. So the first Recv below pairs with the first Send,
	// with no opener echo.
	for _, coll := range []string{"alpha", "beta"} {
		if err := stream.Send(&master_pb.AssignRequest{Count: 2, Collection: coll}); err != nil {
			t.Fatalf("Send %s: %v", coll, err)
		}
		resp, err := stream.Recv()
		if err != nil {
			t.Fatalf("Recv %s: %v", coll, err)
		}
		if resp.Fid != "fid-"+coll || resp.Count != 2 {
			t.Fatalf("assign response lost: %+v", resp)
		}
		if resp.Location == nil || resp.Location.Url != "assign" {
			t.Fatalf("assign location lost: %+v", resp.Location)
		}
	}
	if ms.gotAssignColl != "alpha" {
		t.Fatalf("server saw first collection %q, want alpha", ms.gotAssignColl)
	}
	if err := stream.CloseSend(); err != nil {
		t.Fatalf("CloseSend: %v", err)
	}
	if _, err := stream.Recv(); err != nil && err != io.EOF {
		t.Fatalf("final Recv: want io.EOF, got %v", err)
	}
}

func TestServerSendHeartbeatBidi(t *testing.T) {
	cli, _, done := serveEngine(t)
	defer done()

	stream, err := cli.SendHeartbeat(context.Background())
	if err != nil {
		t.Fatalf("SendHeartbeat: %v", err)
	}
	// The empty opener is discarded server-side; the first real Heartbeat the
	// client Sends is the engine's first Recv.
	if err := stream.Send(&master_pb.Heartbeat{Ip: "10.0.0.2", Port: 18080, Id: "vs-1"}); err != nil {
		t.Fatalf("Send: %v", err)
	}
	resp, err := stream.Recv()
	if err != nil {
		t.Fatalf("Recv: %v", err)
	}
	if resp.Leader != "leader:9333" || resp.VolumeSizeLimit != 30000 || resp.MetricsIntervalSeconds != 15 {
		t.Fatalf("heartbeat response lost: %+v", resp)
	}
	// The heartbeat id reflected back through the storage-backend id proves the
	// request fields round-tripped server-side.
	if len(resp.StorageBackends) != 1 || resp.StorageBackends[0].Id != "vs-1" {
		t.Fatalf("heartbeat request id not reflected: %+v", resp.StorageBackends)
	}
	if err := stream.CloseSend(); err != nil {
		t.Fatalf("CloseSend: %v", err)
	}
	if _, err := stream.Recv(); err != nil && err != io.EOF {
		t.Fatalf("final Recv: want io.EOF, got %v", err)
	}
}

func TestServerKeepConnectedBidi(t *testing.T) {
	cli, _, done := serveEngine(t)
	defer done()

	stream, err := cli.KeepConnected(context.Background())
	if err != nil {
		t.Fatalf("KeepConnected: %v", err)
	}
	// The opener frame carries the subscription; the engine answers once with a
	// VolumeLocation update, reflecting the request's data center.
	if err := stream.Send(&master_pb.KeepConnectedRequest{
		ClientType: "filer", ClientAddress: "10.0.0.3:8888", DataCenter: "dc1",
	}); err != nil {
		t.Fatalf("Send: %v", err)
	}
	resp, err := stream.Recv()
	if err != nil {
		t.Fatalf("Recv: %v", err)
	}
	if resp.VolumeLocation == nil || resp.VolumeLocation.Url != "vs-1:18080" {
		t.Fatalf("keep-connected volume location lost: %+v", resp.VolumeLocation)
	}
	if resp.VolumeLocation.DataCenter != "dc1" {
		t.Fatalf("keep-connected request dc not reflected: %+v", resp.VolumeLocation)
	}
	if len(resp.VolumeLocation.NewVids) != 1 || resp.VolumeLocation.NewVids[0] != 9 {
		t.Fatalf("keep-connected new vids lost: %+v", resp.VolumeLocation.NewVids)
	}
}
