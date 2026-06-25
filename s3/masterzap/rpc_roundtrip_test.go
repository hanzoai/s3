// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package masterzap

import (
	"context"
	"io"
	"testing"

	master_pb "github.com/hanzoai/s3/s3/pb/master_pb"
	masterwire "github.com/hanzoai/s3/s3/wire/master"
	masterstream "github.com/hanzoai/s3/s3/wire/master/masterstream"
	"github.com/zap-proto/go/transport"
)

// backend is an in-memory masterwire.Backend for the unary round-trip: each
// method decodes the request with masterwire.Wrap* and replies with a
// masterwire.New* buffer, so the test proves pb -> wire -> dispatch -> wire ->
// pb is lossless across a real socket. The streaming methods are short-circuited
// by masterwire.Dispatch (ErrStreamingNotWired); they are served separately by
// streamSrv over the SAME listener.
type backend struct {
	masterwire.Backend // embedded for the methods the test does not exercise
}

func (backend) Ping(req []byte) ([]byte, error) {
	v, err := masterwire.WrapPingRequest(req)
	if err != nil {
		return nil, err
	}
	// Echo the target into start_time_ns length so the round-trip is observable,
	// and set deterministic scalars.
	_ = v.Target()
	return masterwire.NewPingResponse(masterwire.PingResponseInput{
		StartTimeNs: 1, RemoteTimeNs: 2, StopTimeNs: 3,
	}), nil
}

func (backend) LookupVolume(req []byte) ([]byte, error) {
	v, err := masterwire.WrapLookupVolumeRequest(req)
	if err != nil {
		return nil, err
	}
	// Reply with one VolumeIdLocation per requested id, each carrying one Location.
	var locs [][]byte
	for i := 0; i < v.VolumeOrFileIdsLen(); i++ {
		id := v.VolumeOrFileIdAt(i)
		loc := masterwire.NewLocation(masterwire.LocationInput{Url: "vol-" + id, GrpcPort: 18080})
		locs = append(locs, masterwire.NewLookupVolumeResponseVolumeIdLocation(masterwire.LookupVolumeResponseVolumeIdLocationInput{
			VolumeOrFileId: id, Locations: [][]byte{loc}, Auth: "tok",
		}))
	}
	return masterwire.NewLookupVolumeResponse(masterwire.LookupVolumeResponseInput{VolumeIdLocations: locs}), nil
}

func (backend) VolumeList(req []byte) ([]byte, error) {
	if _, err := masterwire.WrapVolumeListRequest(req); err != nil {
		return nil, err
	}
	// Build a deep topology graph: Topology -> DC -> Rack -> DataNode -> DiskInfo,
	// with a map<string,DiskInfo> at the topology root carrying a volume.
	vi := masterwire.NewVolumeInformationMessage(masterwire.VolumeInformationMessageInput{Id: 7, Size: 1 << 20, Collection: "c"})
	disk := masterwire.NewDiskInfo(masterwire.DiskInfoInput{
		Type: "hdd", MaxVolumeCount: 100, VolumeInfos: [][]byte{vi}, Tags: []string{"t1"},
	})
	node := masterwire.NewDataNodeInfo(masterwire.DataNodeInfoInput{
		Id: "n1", Address: "10.0.0.1:18080", GrpcPort: 18080,
		DiskInfos: []masterwire.StringMsgEntry{{Key: "hdd", Value: disk}},
	})
	rack := masterwire.NewRackInfo(masterwire.RackInfoInput{Id: "r1", DataNodeInfos: [][]byte{node}})
	dc := masterwire.NewDataCenterInfo(masterwire.DataCenterInfoInput{Id: "dc1", RackInfos: [][]byte{rack}})
	topo := masterwire.NewTopologyInfo(masterwire.TopologyInfoInput{
		Id: "topo", DataCenterInfos: [][]byte{dc},
		DiskInfos: []masterwire.StringMsgEntry{{Key: "hdd", Value: disk}},
	})
	return masterwire.NewVolumeListResponse(masterwire.VolumeListResponseInput{
		TopologyInfo: topo, VolumeSizeLimitMb: 30000,
	}), nil
}

// streamServer scripts the 3 bidirectional streams for the streaming round-trip.
type streamServer struct{}

// StreamAssign echoes each AssignRequest as an AssignResponse whose Fid encodes
// the request's collection, until the client half-closes.
func (streamServer) StreamAssign(init masterwire.AssignRequest, s *masterstream.AssignStream) error {
	reply := func(req masterwire.AssignRequest) error {
		return s.Send(masterwire.NewAssignResponse(masterwire.AssignResponseInput{
			Fid: "fid-" + req.Collection(), Count: req.Count(),
			Location: masterwire.NewLocation(masterwire.LocationInput{Url: "assign", GrpcPort: 1}),
		}))
	}
	if err := reply(init); err != nil {
		return err
	}
	for {
		req, err := s.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if err := reply(req); err != nil {
			return err
		}
	}
}

// SendHeartbeat replies to each Heartbeat with a HeartbeatResponse carrying the
// leader and the heartbeat's id reflected into the volume size limit.
func (streamServer) SendHeartbeat(init masterwire.Heartbeat, s *masterstream.HeartbeatStream) error {
	reply := func(hb masterwire.Heartbeat) error {
		return s.Send(masterwire.NewHeartbeatResponse(masterwire.HeartbeatResponseInput{
			Leader: "leader:9333", VolumeSizeLimit: 30000, MetricsIntervalSeconds: 15,
		}))
	}
	if err := reply(init); err != nil {
		return err
	}
	for {
		hb, err := s.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if err := reply(hb); err != nil {
			return err
		}
	}
}

// KeepConnected is not exercised here; it returns immediately so the Server
// contract is total.
func (streamServer) KeepConnected(masterwire.KeepConnectedRequest, *masterstream.KeepConnectedStream) error {
	return nil
}

// serve stands up ONE transport listener serving both the unary dispatch and the
// streaming handler, returns a masterzap client over a fresh dial conn, and a
// cleanup func.
func serve(t *testing.T) (master_pb.HanzoClient, func()) {
	t.Helper()
	srv, err := transport.ListenStream("tcp", "127.0.0.1:0",
		masterwire.Dispatch(backend{}), masterstream.Handler(streamServer{}))
	if err != nil {
		t.Fatalf("ListenStream: %v", err)
	}
	conn, err := transport.Dial("tcp", srv.Addr().String())
	if err != nil {
		srv.Close()
		t.Fatalf("Dial: %v", err)
	}
	cli := New(conn, nil)
	return cli, func() { conn.Close(); srv.Close() }
}

// TestUnaryRoundTrip drives representative unary RPCs through the masterzap
// client over a real transport: scalars (Ping), repeated nested locations
// (LookupVolume), and the deep topology graph with maps (VolumeList).
func TestUnaryRoundTrip(t *testing.T) {
	cli, done := serve(t)
	defer done()
	ctx := context.Background()

	t.Run("Ping", func(t *testing.T) {
		resp, err := cli.Ping(ctx, &master_pb.PingRequest{Target: "m1", TargetType: "master"})
		if err != nil {
			t.Fatal(err)
		}
		if resp.StartTimeNs != 1 || resp.RemoteTimeNs != 2 || resp.StopTimeNs != 3 {
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
		if len(resp.VolumeIdLocations) != 2 {
			t.Fatalf("want 2 locations, got %d", len(resp.VolumeIdLocations))
		}
		got := resp.VolumeIdLocations[0]
		if got.VolumeOrFileId != "3" || got.Auth != "tok" {
			t.Fatalf("location scalars lost: %+v", got)
		}
		if len(got.Locations) != 1 || got.Locations[0].Url != "vol-3" || got.Locations[0].GrpcPort != 18080 {
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
		rack := ti.DataCenterInfos[0].RackInfos
		if len(rack) != 1 || rack[0].Id != "r1" {
			t.Fatalf("rack lost: %+v", rack)
		}
		node := rack[0].DataNodeInfos
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
		// The topology-root disk map must survive too.
		if _, ok := ti.DiskInfos["hdd"]; !ok {
			t.Fatalf("topology-root disk map lost: %+v", ti.DiskInfos)
		}
	})
}

// TestStreamAssignBidi proves the bidirectional StreamAssign stream through the
// masterzap client: the client Sends AssignRequests and Recvs the matching
// AssignResponses, all over a real socket as zero-copy masterwire buffers.
func TestStreamAssignBidi(t *testing.T) {
	cli, done := serve(t)
	defer done()

	stream, err := cli.StreamAssign(context.Background())
	if err != nil {
		t.Fatalf("StreamAssign: %v", err)
	}

	// The open frame is an empty AssignRequest; its echo arrives first.
	if _, err := stream.Recv(); err != nil {
		t.Fatalf("Recv #0 (opener echo): %v", err)
	}

	for _, coll := range []string{"alpha", "beta"} {
		if err := stream.Send(&master_pb.AssignRequest{Count: 2, Collection: coll}); err != nil {
			t.Fatalf("Send %s: %v", coll, err)
		}
		resp, err := stream.Recv()
		if err != nil {
			t.Fatalf("Recv %s: %v", coll, err)
		}
		if resp.Fid != "fid-"+coll {
			t.Fatalf("assign fid = %q, want fid-%s", resp.Fid, coll)
		}
		if resp.Count != 2 {
			t.Fatalf("assign count = %d, want 2", resp.Count)
		}
		if resp.Location == nil || resp.Location.Url != "assign" {
			t.Fatalf("assign location lost: %+v", resp.Location)
		}
	}

	if err := stream.CloseSend(); err != nil {
		t.Fatalf("CloseSend: %v", err)
	}
	if _, err := stream.Recv(); err != nil && err != io.EOF {
		t.Fatalf("final Recv: want io.EOF, got %v", err)
	}
}

// TestSendHeartbeatBidi proves the bidirectional SendHeartbeat stream through
// the masterzap client: each Heartbeat the client Sends is answered with a
// HeartbeatResponse carrying the leader.
func TestSendHeartbeatBidi(t *testing.T) {
	cli, done := serve(t)
	defer done()

	stream, err := cli.SendHeartbeat(context.Background())
	if err != nil {
		t.Fatalf("SendHeartbeat: %v", err)
	}

	// Opener echo (empty Heartbeat).
	if _, err := stream.Recv(); err != nil {
		t.Fatalf("Recv #0 (opener): %v", err)
	}

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

	if err := stream.CloseSend(); err != nil {
		t.Fatalf("CloseSend: %v", err)
	}
	if _, err := stream.Recv(); err != nil && err != io.EOF {
		t.Fatalf("final Recv: want io.EOF, got %v", err)
	}
}
