// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// server_rpc.go is the server-side mirror of rpc.go: the per-RPC request
// DECODERS and response BUILDERS for the 21 unary Hanzo (master) RPCs. Where
// rpc.go provides the client direction — <Rpc>ReqToWire (build the request) and
// <Rpc>RespFromWire (decode the reply) — this file provides the inverse pair the
// ZAP server backend needs: <Rpc>ReqFromWire (decode the request wire view into
// *master_pb.<Rpc>Request) and <Rpc>RespToWire (encode the *master_pb.<Rpc>Response
// the engine returned as the reply wire buffer). Together the two directions
// close the loop so a request built by the client is decoded here, run through
// the unchanged master engine, and its response re-encoded for the client to
// decode — all over ZAP, no protobuf framing. Built on the leaf converters in
// convert.go; this is the ONLY place the server-direction master unary mapping
// lives (server.go composes them).

package master

import (
	master_pb "github.com/hanzoai/s3/s3/pb/master_pb"
	masterwire "github.com/hanzoai/s3/s3/wire/master"
)

// LookupVolume
func LookupVolumeReqFromWire(b []byte) (*master_pb.LookupVolumeRequest, error) {
	v, err := masterwire.WrapLookupVolumeRequest(b)
	if err != nil {
		return nil, err
	}
	req := &master_pb.LookupVolumeRequest{Collection: v.Collection()}
	if n := v.VolumeOrFileIdsLen(); n > 0 {
		req.VolumeOrFileIds = make([]string, n)
		for i := 0; i < n; i++ {
			req.VolumeOrFileIds[i] = v.VolumeOrFileIdAt(i)
		}
	}
	return req, nil
}
func LookupVolumeRespToWire(r *master_pb.LookupVolumeResponse) []byte {
	locs := make([][]byte, len(r.VolumeIdLocations))
	for i, l := range r.VolumeIdLocations {
		locs[i] = masterwire.NewLookupVolumeResponseVolumeIdLocation(masterwire.LookupVolumeResponseVolumeIdLocationInput{
			VolumeOrFileId: l.VolumeOrFileId,
			Locations:      locationsToWire(l.Locations),
			Error:          l.Error,
			Auth:           l.Auth,
		})
	}
	return masterwire.NewLookupVolumeResponse(masterwire.LookupVolumeResponseInput{VolumeIdLocations: locs})
}

// Assign
func AssignReqFromWire(b []byte) (*master_pb.AssignRequest, error) {
	v, err := masterwire.WrapAssignRequest(b)
	if err != nil {
		return nil, err
	}
	return assignReqFromView(v), nil
}
func assignReqFromView(v masterwire.AssignRequest) *master_pb.AssignRequest {
	return &master_pb.AssignRequest{
		Count:               v.Count(),
		Replication:         v.Replication(),
		Collection:          v.Collection(),
		Ttl:                 v.Ttl(),
		DataCenter:          v.DataCenter(),
		Rack:                v.Rack(),
		DataNode:            v.DataNode(),
		MemoryMapMaxSizeMb:  v.MemoryMapMaxSizeMb(),
		WritableVolumeCount: v.WritableVolumeCount(),
		DiskType:            v.DiskType(),
		ExpectedDataSize:    v.ExpectedDataSize(),
	}
}
func AssignRespToWire(r *master_pb.AssignResponse) []byte {
	return masterwire.NewAssignResponse(masterwire.AssignResponseInput{
		Fid:      r.Fid,
		Location: locationToWire(r.Location),
		Count:    r.Count,
		Error:    r.Error,
		Auth:     r.Auth,
		Replicas: locationsToWire(r.Replicas),
	})
}

// Statistics
func StatisticsReqFromWire(b []byte) (*master_pb.StatisticsRequest, error) {
	v, err := masterwire.WrapStatisticsRequest(b)
	if err != nil {
		return nil, err
	}
	return &master_pb.StatisticsRequest{
		Replication: v.Replication(), Collection: v.Collection(), Ttl: v.Ttl(), DiskType: v.DiskType(),
	}, nil
}
func StatisticsRespToWire(r *master_pb.StatisticsResponse) []byte {
	return masterwire.NewStatisticsResponse(masterwire.StatisticsResponseInput{
		TotalSize: r.TotalSize, UsedSize: r.UsedSize, FileCount: r.FileCount,
	})
}

// CollectionList
func CollectionListReqFromWire(b []byte) (*master_pb.CollectionListRequest, error) {
	v, err := masterwire.WrapCollectionListRequest(b)
	if err != nil {
		return nil, err
	}
	return &master_pb.CollectionListRequest{
		IncludeNormalVolumes: v.IncludeNormalVolumes(), IncludeEcVolumes: v.IncludeEcVolumes(),
	}, nil
}
func CollectionListRespToWire(r *master_pb.CollectionListResponse) []byte {
	cols := make([][]byte, len(r.Collections))
	for i, c := range r.Collections {
		cols[i] = masterwire.NewCollection(masterwire.CollectionInput{Name: c.Name})
	}
	return masterwire.NewCollectionListResponse(masterwire.CollectionListResponseInput{Collections: cols})
}

// CollectionDelete
func CollectionDeleteReqFromWire(b []byte) (*master_pb.CollectionDeleteRequest, error) {
	v, err := masterwire.WrapCollectionDeleteRequest(b)
	if err != nil {
		return nil, err
	}
	return &master_pb.CollectionDeleteRequest{Name: v.Name()}, nil
}
func CollectionDeleteRespToWire(*master_pb.CollectionDeleteResponse) []byte {
	return masterwire.NewCollectionDeleteResponse(masterwire.CollectionDeleteResponseInput{})
}

// VolumeList
func VolumeListReqFromWire(b []byte) (*master_pb.VolumeListRequest, error) {
	if _, err := masterwire.WrapVolumeListRequest(b); err != nil {
		return nil, err
	}
	return &master_pb.VolumeListRequest{}, nil
}
func VolumeListRespToWire(r *master_pb.VolumeListResponse) []byte {
	return masterwire.NewVolumeListResponse(masterwire.VolumeListResponseInput{
		TopologyInfo:      topologyInfoToWire(r.TopologyInfo),
		VolumeSizeLimitMb: r.VolumeSizeLimitMb,
	})
}

// LookupEcVolume
func LookupEcVolumeReqFromWire(b []byte) (*master_pb.LookupEcVolumeRequest, error) {
	v, err := masterwire.WrapLookupEcVolumeRequest(b)
	if err != nil {
		return nil, err
	}
	return &master_pb.LookupEcVolumeRequest{VolumeId: v.VolumeId()}, nil
}
func LookupEcVolumeRespToWire(r *master_pb.LookupEcVolumeResponse) []byte {
	shards := make([][]byte, len(r.ShardIdLocations))
	for i, s := range r.ShardIdLocations {
		shards[i] = masterwire.NewLookupEcVolumeResponseEcShardIdLocation(masterwire.LookupEcVolumeResponseEcShardIdLocationInput{
			ShardId:   s.ShardId,
			Locations: locationsToWire(s.Locations),
		})
	}
	return masterwire.NewLookupEcVolumeResponse(masterwire.LookupEcVolumeResponseInput{
		VolumeId: r.VolumeId, ShardIdLocations: shards,
	})
}

// VacuumVolume
func VacuumVolumeReqFromWire(b []byte) (*master_pb.VacuumVolumeRequest, error) {
	v, err := masterwire.WrapVacuumVolumeRequest(b)
	if err != nil {
		return nil, err
	}
	return &master_pb.VacuumVolumeRequest{
		GarbageThreshold: v.GarbageThreshold(), VolumeId: v.VolumeId(), Collection: v.Collection(),
	}, nil
}
func VacuumVolumeRespToWire(*master_pb.VacuumVolumeResponse) []byte {
	return masterwire.NewVacuumVolumeResponse(masterwire.VacuumVolumeResponseInput{})
}

// DisableVacuum
func DisableVacuumReqFromWire(b []byte) (*master_pb.DisableVacuumRequest, error) {
	v, err := masterwire.WrapDisableVacuumRequest(b)
	if err != nil {
		return nil, err
	}
	return &master_pb.DisableVacuumRequest{ByPlugin: v.ByPlugin()}, nil
}
func DisableVacuumRespToWire(*master_pb.DisableVacuumResponse) []byte {
	return masterwire.NewDisableVacuumResponse(masterwire.DisableVacuumResponseInput{})
}

// EnableVacuum
func EnableVacuumReqFromWire(b []byte) (*master_pb.EnableVacuumRequest, error) {
	v, err := masterwire.WrapEnableVacuumRequest(b)
	if err != nil {
		return nil, err
	}
	return &master_pb.EnableVacuumRequest{ByPlugin: v.ByPlugin()}, nil
}
func EnableVacuumRespToWire(*master_pb.EnableVacuumResponse) []byte {
	return masterwire.NewEnableVacuumResponse(masterwire.EnableVacuumResponseInput{})
}

// VolumeMarkReadonly
func VolumeMarkReadonlyReqFromWire(b []byte) (*master_pb.VolumeMarkReadonlyRequest, error) {
	v, err := masterwire.WrapVolumeMarkReadonlyRequest(b)
	if err != nil {
		return nil, err
	}
	return &master_pb.VolumeMarkReadonlyRequest{
		Ip:               v.Ip(),
		Port:             v.Port(),
		VolumeId:         v.VolumeId(),
		Collection:       v.Collection(),
		ReplicaPlacement: v.ReplicaPlacement(),
		Version:          v.Version(),
		Ttl:              v.Ttl(),
		DiskType:         v.DiskType(),
		IsReadonly:       v.IsReadonly(),
	}, nil
}
func VolumeMarkReadonlyRespToWire(*master_pb.VolumeMarkReadonlyResponse) []byte {
	return masterwire.NewVolumeMarkReadonlyResponse(masterwire.VolumeMarkReadonlyResponseInput{})
}

// GetMasterConfiguration
func GetMasterConfigurationReqFromWire(b []byte) (*master_pb.GetMasterConfigurationRequest, error) {
	if _, err := masterwire.WrapGetMasterConfigurationRequest(b); err != nil {
		return nil, err
	}
	return &master_pb.GetMasterConfigurationRequest{}, nil
}
func GetMasterConfigurationRespToWire(r *master_pb.GetMasterConfigurationResponse) []byte {
	return masterwire.NewGetMasterConfigurationResponse(masterwire.GetMasterConfigurationResponseInput{
		MetricsAddress:          r.MetricsAddress,
		MetricsIntervalSeconds:  r.MetricsIntervalSeconds,
		StorageBackends:         storageBackendsToWire(r.StorageBackends),
		DefaultReplication:      r.DefaultReplication,
		Leader:                  r.Leader,
		VolumeSizeLimitMB:       r.VolumeSizeLimitMB,
		VolumePreallocate:       r.VolumePreallocate,
		MaintenanceScripts:      r.MaintenanceScripts,
		MaintenanceSleepMinutes: r.MaintenanceSleepMinutes,
	})
}

// ListClusterNodes
func ListClusterNodesReqFromWire(b []byte) (*master_pb.ListClusterNodesRequest, error) {
	v, err := masterwire.WrapListClusterNodesRequest(b)
	if err != nil {
		return nil, err
	}
	return &master_pb.ListClusterNodesRequest{
		ClientType: v.ClientType(), FilerGroup: v.FilerGroup(), Limit: v.Limit(),
	}, nil
}
func ListClusterNodesRespToWire(r *master_pb.ListClusterNodesResponse) []byte {
	nodes := make([][]byte, len(r.ClusterNodes))
	for i, n := range r.ClusterNodes {
		nodes[i] = masterwire.NewListClusterNodesResponseClusterNode(masterwire.ListClusterNodesResponseClusterNodeInput{
			Address: n.Address, Version: n.Version, CreatedAtNs: n.CreatedAtNs,
			DataCenter: n.DataCenter, Rack: n.Rack,
		})
	}
	return masterwire.NewListClusterNodesResponse(masterwire.ListClusterNodesResponseInput{ClusterNodes: nodes})
}

// LeaseAdminToken
func LeaseAdminTokenReqFromWire(b []byte) (*master_pb.LeaseAdminTokenRequest, error) {
	v, err := masterwire.WrapLeaseAdminTokenRequest(b)
	if err != nil {
		return nil, err
	}
	return &master_pb.LeaseAdminTokenRequest{
		PreviousToken: v.PreviousToken(), PreviousLockTime: v.PreviousLockTime(),
		LockName: v.LockName(), ClientName: v.ClientName(), Message: v.Message(),
	}, nil
}
func LeaseAdminTokenRespToWire(r *master_pb.LeaseAdminTokenResponse) []byte {
	return masterwire.NewLeaseAdminTokenResponse(masterwire.LeaseAdminTokenResponseInput{
		Token: r.Token, LockTsNs: r.LockTsNs,
	})
}

// ReleaseAdminToken
func ReleaseAdminTokenReqFromWire(b []byte) (*master_pb.ReleaseAdminTokenRequest, error) {
	v, err := masterwire.WrapReleaseAdminTokenRequest(b)
	if err != nil {
		return nil, err
	}
	return &master_pb.ReleaseAdminTokenRequest{
		PreviousToken: v.PreviousToken(), PreviousLockTime: v.PreviousLockTime(), LockName: v.LockName(),
	}, nil
}
func ReleaseAdminTokenRespToWire(*master_pb.ReleaseAdminTokenResponse) []byte {
	return masterwire.NewReleaseAdminTokenResponse(masterwire.ReleaseAdminTokenResponseInput{})
}

// Ping
func PingReqFromWire(b []byte) (*master_pb.PingRequest, error) {
	v, err := masterwire.WrapPingRequest(b)
	if err != nil {
		return nil, err
	}
	return &master_pb.PingRequest{Target: v.Target(), TargetType: v.TargetType()}, nil
}
func PingRespToWire(r *master_pb.PingResponse) []byte {
	return masterwire.NewPingResponse(masterwire.PingResponseInput{
		StartTimeNs: r.StartTimeNs, RemoteTimeNs: r.RemoteTimeNs, StopTimeNs: r.StopTimeNs,
	})
}

// RaftListClusterServers
func RaftListClusterServersReqFromWire(b []byte) (*master_pb.RaftListClusterServersRequest, error) {
	if _, err := masterwire.WrapRaftListClusterServersRequest(b); err != nil {
		return nil, err
	}
	return &master_pb.RaftListClusterServersRequest{}, nil
}
func RaftListClusterServersRespToWire(r *master_pb.RaftListClusterServersResponse) []byte {
	servers := make([][]byte, len(r.ClusterServers))
	for i, s := range r.ClusterServers {
		servers[i] = masterwire.NewRaftListClusterServersResponseClusterServers(masterwire.RaftListClusterServersResponseClusterServersInput{
			Id: s.Id, Address: s.Address, Suffrage: s.Suffrage, IsWriter: s.IsWriter,
		})
	}
	return masterwire.NewRaftListClusterServersResponse(masterwire.RaftListClusterServersResponseInput{ClusterServers: servers})
}

// RaftAddServer
func RaftAddServerReqFromWire(b []byte) (*master_pb.RaftAddServerRequest, error) {
	v, err := masterwire.WrapRaftAddServerRequest(b)
	if err != nil {
		return nil, err
	}
	return &master_pb.RaftAddServerRequest{Id: v.Id(), Address: v.Address(), Voter: v.Voter()}, nil
}
func RaftAddServerRespToWire(*master_pb.RaftAddServerResponse) []byte {
	return masterwire.NewRaftAddServerResponse(masterwire.RaftAddServerResponseInput{})
}

// RaftRemoveServer
func RaftRemoveServerReqFromWire(b []byte) (*master_pb.RaftRemoveServerRequest, error) {
	v, err := masterwire.WrapRaftRemoveServerRequest(b)
	if err != nil {
		return nil, err
	}
	return &master_pb.RaftRemoveServerRequest{Id: v.Id(), Force: v.Force()}, nil
}
func RaftRemoveServerRespToWire(*master_pb.RaftRemoveServerResponse) []byte {
	return masterwire.NewRaftRemoveServerResponse(masterwire.RaftRemoveServerResponseInput{})
}

// RaftLeadershipTransfer
func RaftLeadershipTransferReqFromWire(b []byte) (*master_pb.RaftLeadershipTransferRequest, error) {
	v, err := masterwire.WrapRaftLeadershipTransferRequest(b)
	if err != nil {
		return nil, err
	}
	return &master_pb.RaftLeadershipTransferRequest{TargetId: v.TargetId(), TargetAddress: v.TargetAddress()}, nil
}
func RaftLeadershipTransferRespToWire(r *master_pb.RaftLeadershipTransferResponse) []byte {
	return masterwire.NewRaftLeadershipTransferResponse(masterwire.RaftLeadershipTransferResponseInput{
		PreviousLeader: r.PreviousLeader, NewLeader: r.NewLeader,
	})
}

// VolumeGrow
func VolumeGrowReqFromWire(b []byte) (*master_pb.VolumeGrowRequest, error) {
	v, err := masterwire.WrapVolumeGrowRequest(b)
	if err != nil {
		return nil, err
	}
	return &master_pb.VolumeGrowRequest{
		WritableVolumeCount: v.WritableVolumeCount(),
		Replication:         v.Replication(),
		Collection:          v.Collection(),
		Ttl:                 v.Ttl(),
		DataCenter:          v.DataCenter(),
		Rack:                v.Rack(),
		DataNode:            v.DataNode(),
		MemoryMapMaxSizeMb:  v.MemoryMapMaxSizeMb(),
		DiskType:            v.DiskType(),
	}, nil
}
func VolumeGrowRespToWire(*master_pb.VolumeGrowResponse) []byte {
	return masterwire.NewVolumeGrowResponse(masterwire.VolumeGrowResponseInput{})
}
