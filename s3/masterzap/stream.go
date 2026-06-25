// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// stream.go holds the per-frame request builders and response decoders for the
// Hanzo (master) STREAMING RPCs (SendHeartbeat, KeepConnected, StreamAssign).
// They compose the leaf converters in convert.go so the streaming seam carries
// the same *master_pb message graph the topology/volume engine already speaks
// while every frame crosses the socket as a masterwire buffer. Naming mirrors
// rpc.go: <Rpc>ReqToWire(*pb.<Rpc>Request) []byte and
// <Rpc>RespFromWire([]byte) (*pb.<Rpc>Response, error).
//
// StreamAssign carries AssignRequest/AssignResponse — its converters live in
// rpc.go (AssignReqToWire / AssignRespFromWire) since the unary Assign RPC uses
// the same messages; the bidi adapter in client.go reuses them directly.

package masterzap

import (
	master_pb "github.com/hanzoai/s3/s3/pb/master_pb"
	masterwire "github.com/hanzoai/s3/s3/wire/master"
)

// SendHeartbeat (bidi): Heartbeat -> HeartbeatResponse.

func HeartbeatReqToWire(r *master_pb.Heartbeat) []byte {
	in := masterwire.HeartbeatInput{
		Ip:              r.Ip,
		Port:            r.Port,
		PublicUrl:       r.PublicUrl,
		MaxFileKey:      r.MaxFileKey,
		DataCenter:      r.DataCenter,
		Rack:            r.Rack,
		AdminPort:       r.AdminPort,
		HasNoVolumes:    r.HasNoVolumes,
		HasNoEcShards:   r.HasNoEcShards,
		GrpcPort:        r.GrpcPort,
		LocationUuids:   r.LocationUuids,
		Id:              r.Id,
		State:           stateToWire(r.State),
		Volumes:         volumeInfoListToWire(r.Volumes),
		NewVolumes:      volumeShortInfoListToWire(r.NewVolumes),
		DeletedVolumes:  volumeShortInfoListToWire(r.DeletedVolumes),
		EcShards:        ecShardInfoListToWire(r.EcShards),
		NewEcShards:     ecShardInfoListToWire(r.NewEcShards),
		DeletedEcShards: ecShardInfoListToWire(r.DeletedEcShards),
		DiskTags:        diskTagListToWire(r.DiskTags),
	}
	for k, v := range r.MaxVolumeCounts {
		in.MaxVolumeCounts = append(in.MaxVolumeCounts, masterwire.StringUint32Entry{Key: k, Value: v})
	}
	return masterwire.NewHeartbeat(in)
}

func HeartbeatReqFromWire(b []byte) (*master_pb.Heartbeat, error) {
	v, err := masterwire.WrapHeartbeat(b)
	if err != nil {
		return nil, err
	}
	r := &master_pb.Heartbeat{
		Ip:            v.Ip(),
		Port:          v.Port(),
		PublicUrl:     v.PublicUrl(),
		MaxFileKey:    v.MaxFileKey(),
		DataCenter:    v.DataCenter(),
		Rack:          v.Rack(),
		AdminPort:     v.AdminPort(),
		HasNoVolumes:  v.HasNoVolumes(),
		HasNoEcShards: v.HasNoEcShards(),
		GrpcPort:      v.GrpcPort(),
		Id:            v.Id(),
	}
	if n := v.VolumesLen(); n > 0 {
		r.Volumes = make([]*master_pb.VolumeInformationMessage, 0, n)
		for i := 0; i < n; i++ {
			if m, ok := v.VolumeAt(i); ok {
				r.Volumes = append(r.Volumes, volumeInfoFromView(m))
			}
		}
	}
	if n := v.EcShardsLen(); n > 0 {
		r.EcShards = make([]*master_pb.VolumeEcShardInformationMessage, 0, n)
		for i := 0; i < n; i++ {
			if m, ok := v.EcShardAt(i); ok {
				r.EcShards = append(r.EcShards, ecShardInfoFromView(m))
			}
		}
	}
	if n := v.LocationUuidsLen(); n > 0 {
		r.LocationUuids = make([]string, n)
		for i := 0; i < n; i++ {
			r.LocationUuids[i] = v.LocationUuidAt(i)
		}
	}
	if n := v.MaxVolumeCountsLen(); n > 0 {
		r.MaxVolumeCounts = make(map[string]uint32, n)
		for i := 0; i < n; i++ {
			e := v.MaxVolumeCountAt(i)
			r.MaxVolumeCounts[e.Key] = e.Value
		}
	}
	if s, ok := v.State(); ok {
		r.State = stateFromView(s)
	}
	return r, nil
}

func HeartbeatRespToWire(r *master_pb.HeartbeatResponse) []byte {
	return masterwire.NewHeartbeatResponse(masterwire.HeartbeatResponseInput{
		VolumeSizeLimit:        r.VolumeSizeLimit,
		Leader:                 r.Leader,
		MetricsAddress:         r.MetricsAddress,
		MetricsIntervalSeconds: r.MetricsIntervalSeconds,
		StorageBackends:        storageBackendsToWire(r.StorageBackends),
		DuplicatedUuids:        r.DuplicatedUuids,
		Preallocate:            r.Preallocate,
	})
}

func HeartbeatRespFromWire(b []byte) (*master_pb.HeartbeatResponse, error) {
	v, err := masterwire.WrapHeartbeatResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &master_pb.HeartbeatResponse{
		VolumeSizeLimit:        v.VolumeSizeLimit(),
		Leader:                 v.Leader(),
		MetricsAddress:         v.MetricsAddress(),
		MetricsIntervalSeconds: v.MetricsIntervalSeconds(),
		Preallocate:            v.Preallocate(),
	}
	for i := 0; i < v.StorageBackendsLen(); i++ {
		if sb, ok := v.StorageBackendAt(i); ok {
			resp.StorageBackends = append(resp.StorageBackends, storageBackendFromView(sb))
		}
	}
	if n := v.DuplicatedUuidsLen(); n > 0 {
		resp.DuplicatedUuids = make([]string, n)
		for i := 0; i < n; i++ {
			resp.DuplicatedUuids[i] = v.DuplicatedUuidAt(i)
		}
	}
	return resp, nil
}

// KeepConnected (bidi): KeepConnectedRequest -> KeepConnectedResponse.

func KeepConnectedReqToWire(r *master_pb.KeepConnectedRequest) []byte {
	return masterwire.NewKeepConnectedRequest(masterwire.KeepConnectedRequestInput{
		ClientType: r.ClientType, ClientAddress: r.ClientAddress, Version: r.Version,
		FilerGroup: r.FilerGroup, DataCenter: r.DataCenter, Rack: r.Rack,
	})
}

func KeepConnectedReqFromWire(b []byte) (*master_pb.KeepConnectedRequest, error) {
	v, err := masterwire.WrapKeepConnectedRequest(b)
	if err != nil {
		return nil, err
	}
	return &master_pb.KeepConnectedRequest{
		ClientType: v.ClientType(), ClientAddress: v.ClientAddress(), Version: v.Version(),
		FilerGroup: v.FilerGroup(), DataCenter: v.DataCenter(), Rack: v.Rack(),
	}, nil
}

func KeepConnectedRespToWire(r *master_pb.KeepConnectedResponse) []byte {
	return masterwire.NewKeepConnectedResponse(masterwire.KeepConnectedResponseInput{
		VolumeLocation:    volumeLocationToWire(r.VolumeLocation),
		ClusterNodeUpdate: clusterNodeUpdateToWire(r.ClusterNodeUpdate),
		LockRingUpdate:    lockRingUpdateToWire(r.LockRingUpdate),
	})
}

func KeepConnectedRespFromWire(b []byte) (*master_pb.KeepConnectedResponse, error) {
	v, err := masterwire.WrapKeepConnectedResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &master_pb.KeepConnectedResponse{}
	if l, ok := v.VolumeLocation(); ok {
		resp.VolumeLocation = volumeLocationFromView(l)
	}
	if c, ok := v.ClusterNodeUpdate(); ok {
		resp.ClusterNodeUpdate = clusterNodeUpdateFromView(c)
	}
	if l, ok := v.LockRingUpdate(); ok {
		resp.LockRingUpdate = lockRingUpdateFromView(l)
	}
	return resp, nil
}

// --- repeated-message list helpers for the Heartbeat builder ---

func volumeInfoListToWire(ms []*master_pb.VolumeInformationMessage) [][]byte {
	if len(ms) == 0 {
		return nil
	}
	out := make([][]byte, len(ms))
	for i, m := range ms {
		out[i] = volumeInfoToWire(m)
	}
	return out
}

func volumeShortInfoListToWire(ms []*master_pb.VolumeShortInformationMessage) [][]byte {
	if len(ms) == 0 {
		return nil
	}
	out := make([][]byte, len(ms))
	for i, m := range ms {
		out[i] = volumeShortInfoToWire(m)
	}
	return out
}

func ecShardInfoListToWire(ms []*master_pb.VolumeEcShardInformationMessage) [][]byte {
	if len(ms) == 0 {
		return nil
	}
	out := make([][]byte, len(ms))
	for i, m := range ms {
		out[i] = ecShardInfoToWire(m)
	}
	return out
}

func diskTagListToWire(ds []*master_pb.DiskTag) [][]byte {
	if len(ds) == 0 {
		return nil
	}
	out := make([][]byte, len(ds))
	for i, d := range ds {
		out[i] = diskTagToWire(d)
	}
	return out
}
