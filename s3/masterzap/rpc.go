// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// rpc.go holds the per-RPC request builders and response decoders for the 21
// unary Hanzo (master) RPCs, built on the leaf converters in convert.go. Naming:
// <Rpc>ReqToWire(*pb.<Rpc>Request) []byte and <Rpc>RespFromWire([]byte) (*pb.<Rpc>Response, error).

package masterzap

import (
	master_pb "github.com/hanzoai/s3/s3/pb/master_pb"
	masterwire "github.com/hanzoai/s3/s3/wire/master"
)

// LookupVolume
func LookupVolumeReqToWire(r *master_pb.LookupVolumeRequest) []byte {
	return masterwire.NewLookupVolumeRequest(masterwire.LookupVolumeRequestInput{
		VolumeOrFileIds: r.VolumeOrFileIds, Collection: r.Collection,
	})
}
func LookupVolumeRespFromWire(b []byte) (*master_pb.LookupVolumeResponse, error) {
	v, err := masterwire.WrapLookupVolumeResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &master_pb.LookupVolumeResponse{}
	for i := 0; i < v.VolumeIdLocationsLen(); i++ {
		l, ok := v.VolumeIdLocationAt(i)
		if !ok {
			continue
		}
		loc := &master_pb.LookupVolumeResponse_VolumeIdLocation{
			VolumeOrFileId: l.VolumeOrFileId(), Error: l.Error(), Auth: l.Auth(),
		}
		for j := 0; j < l.LocationsLen(); j++ {
			if lv, ok := l.LocationAt(j); ok {
				loc.Locations = append(loc.Locations, locationFromView(lv))
			}
		}
		resp.VolumeIdLocations = append(resp.VolumeIdLocations, loc)
	}
	return resp, nil
}

// Assign
func AssignReqToWire(r *master_pb.AssignRequest) []byte {
	return masterwire.NewAssignRequest(masterwire.AssignRequestInput{
		Count:               r.Count,
		Replication:         r.Replication,
		Collection:          r.Collection,
		Ttl:                 r.Ttl,
		DataCenter:          r.DataCenter,
		Rack:                r.Rack,
		DataNode:            r.DataNode,
		MemoryMapMaxSizeMb:  r.MemoryMapMaxSizeMb,
		WritableVolumeCount: r.WritableVolumeCount,
		DiskType:            r.DiskType,
		ExpectedDataSize:    r.ExpectedDataSize,
	})
}
func AssignRespFromWire(b []byte) (*master_pb.AssignResponse, error) {
	v, err := masterwire.WrapAssignResponse(b)
	if err != nil {
		return nil, err
	}
	return assignRespFromView(v), nil
}
func assignRespFromView(v masterwire.AssignResponse) *master_pb.AssignResponse {
	resp := &master_pb.AssignResponse{
		Fid: v.Fid(), Count: v.Count(), Error: v.Error(), Auth: v.Auth(),
	}
	for i := 0; i < v.ReplicasLen(); i++ {
		if r, ok := v.ReplicaAt(i); ok {
			resp.Replicas = append(resp.Replicas, locationFromView(r))
		}
	}
	if loc, ok := v.Location(); ok {
		resp.Location = locationFromView(loc)
	}
	return resp
}

// AssignReqFromWire / AssignRespToWire are the server-direction halves used by
// the StreamAssign ZAP server (the client direction is ReqToWire/RespFromWire).
func AssignReqFromWire(b []byte) (*master_pb.AssignRequest, error) {
	v, err := masterwire.WrapAssignRequest(b)
	if err != nil {
		return nil, err
	}
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
	}, nil
}

func AssignRespToWire(r *master_pb.AssignResponse) []byte {
	var replicas [][]byte
	for _, rep := range r.Replicas {
		replicas = append(replicas, locationToWire(rep))
	}
	return masterwire.NewAssignResponse(masterwire.AssignResponseInput{
		Fid:      r.Fid,
		Count:    r.Count,
		Error:    r.Error,
		Auth:     r.Auth,
		Replicas: replicas,
		Location: locationToWire(r.Location),
	})
}

// Statistics
func StatisticsReqToWire(r *master_pb.StatisticsRequest) []byte {
	return masterwire.NewStatisticsRequest(masterwire.StatisticsRequestInput{
		Replication: r.Replication, Collection: r.Collection, Ttl: r.Ttl, DiskType: r.DiskType,
	})
}
func StatisticsRespFromWire(b []byte) (*master_pb.StatisticsResponse, error) {
	v, err := masterwire.WrapStatisticsResponse(b)
	if err != nil {
		return nil, err
	}
	return &master_pb.StatisticsResponse{TotalSize: v.TotalSize(), UsedSize: v.UsedSize(), FileCount: v.FileCount()}, nil
}

// CollectionList
func CollectionListReqToWire(r *master_pb.CollectionListRequest) []byte {
	return masterwire.NewCollectionListRequest(masterwire.CollectionListRequestInput{
		IncludeNormalVolumes: r.IncludeNormalVolumes, IncludeEcVolumes: r.IncludeEcVolumes,
	})
}
func CollectionListRespFromWire(b []byte) (*master_pb.CollectionListResponse, error) {
	v, err := masterwire.WrapCollectionListResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &master_pb.CollectionListResponse{}
	for i := 0; i < v.CollectionsLen(); i++ {
		if c, ok := v.CollectionAt(i); ok {
			resp.Collections = append(resp.Collections, &master_pb.Collection{Name: c.Name()})
		}
	}
	return resp, nil
}

// CollectionDelete
func CollectionDeleteReqToWire(r *master_pb.CollectionDeleteRequest) []byte {
	return masterwire.NewCollectionDeleteRequest(masterwire.CollectionDeleteRequestInput{Name: r.Name})
}
func CollectionDeleteRespFromWire(b []byte) (*master_pb.CollectionDeleteResponse, error) {
	if _, err := masterwire.WrapCollectionDeleteResponse(b); err != nil {
		return nil, err
	}
	return &master_pb.CollectionDeleteResponse{}, nil
}

// VolumeList
func VolumeListReqToWire(r *master_pb.VolumeListRequest) []byte {
	return masterwire.NewVolumeListRequest(masterwire.VolumeListRequestInput{})
}
func VolumeListRespFromWire(b []byte) (*master_pb.VolumeListResponse, error) {
	v, err := masterwire.WrapVolumeListResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &master_pb.VolumeListResponse{VolumeSizeLimitMb: v.VolumeSizeLimitMb()}
	if ti, ok := v.TopologyInfo(); ok {
		resp.TopologyInfo = topologyInfoFromView(ti)
	}
	return resp, nil
}

// LookupEcVolume
func LookupEcVolumeReqToWire(r *master_pb.LookupEcVolumeRequest) []byte {
	return masterwire.NewLookupEcVolumeRequest(masterwire.LookupEcVolumeRequestInput{VolumeId: r.VolumeId})
}
func LookupEcVolumeRespFromWire(b []byte) (*master_pb.LookupEcVolumeResponse, error) {
	v, err := masterwire.WrapLookupEcVolumeResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &master_pb.LookupEcVolumeResponse{VolumeId: v.VolumeId()}
	for i := 0; i < v.ShardIdLocationsLen(); i++ {
		s, ok := v.ShardIdLocationAt(i)
		if !ok {
			continue
		}
		sl := &master_pb.LookupEcVolumeResponse_EcShardIdLocation{ShardId: s.ShardId()}
		for j := 0; j < s.LocationsLen(); j++ {
			if lv, ok := s.LocationAt(j); ok {
				sl.Locations = append(sl.Locations, locationFromView(lv))
			}
		}
		resp.ShardIdLocations = append(resp.ShardIdLocations, sl)
	}
	return resp, nil
}

// VacuumVolume
func VacuumVolumeReqToWire(r *master_pb.VacuumVolumeRequest) []byte {
	return masterwire.NewVacuumVolumeRequest(masterwire.VacuumVolumeRequestInput{
		GarbageThreshold: r.GarbageThreshold, VolumeId: r.VolumeId, Collection: r.Collection,
	})
}
func VacuumVolumeRespFromWire(b []byte) (*master_pb.VacuumVolumeResponse, error) {
	if _, err := masterwire.WrapVacuumVolumeResponse(b); err != nil {
		return nil, err
	}
	return &master_pb.VacuumVolumeResponse{}, nil
}

// DisableVacuum
func DisableVacuumReqToWire(r *master_pb.DisableVacuumRequest) []byte {
	return masterwire.NewDisableVacuumRequest(masterwire.DisableVacuumRequestInput{ByPlugin: r.ByPlugin})
}
func DisableVacuumRespFromWire(b []byte) (*master_pb.DisableVacuumResponse, error) {
	if _, err := masterwire.WrapDisableVacuumResponse(b); err != nil {
		return nil, err
	}
	return &master_pb.DisableVacuumResponse{}, nil
}

// EnableVacuum
func EnableVacuumReqToWire(r *master_pb.EnableVacuumRequest) []byte {
	return masterwire.NewEnableVacuumRequest(masterwire.EnableVacuumRequestInput{ByPlugin: r.ByPlugin})
}
func EnableVacuumRespFromWire(b []byte) (*master_pb.EnableVacuumResponse, error) {
	if _, err := masterwire.WrapEnableVacuumResponse(b); err != nil {
		return nil, err
	}
	return &master_pb.EnableVacuumResponse{}, nil
}

// VolumeMarkReadonly
func VolumeMarkReadonlyReqToWire(r *master_pb.VolumeMarkReadonlyRequest) []byte {
	return masterwire.NewVolumeMarkReadonlyRequest(masterwire.VolumeMarkReadonlyRequestInput{
		Ip:               r.Ip,
		Port:             r.Port,
		VolumeId:         r.VolumeId,
		Collection:       r.Collection,
		ReplicaPlacement: r.ReplicaPlacement,
		Version:          r.Version,
		Ttl:              r.Ttl,
		DiskType:         r.DiskType,
		IsReadonly:       r.IsReadonly,
	})
}
func VolumeMarkReadonlyRespFromWire(b []byte) (*master_pb.VolumeMarkReadonlyResponse, error) {
	if _, err := masterwire.WrapVolumeMarkReadonlyResponse(b); err != nil {
		return nil, err
	}
	return &master_pb.VolumeMarkReadonlyResponse{}, nil
}

// GetMasterConfiguration
func GetMasterConfigurationReqToWire(r *master_pb.GetMasterConfigurationRequest) []byte {
	return masterwire.NewGetMasterConfigurationRequest(masterwire.GetMasterConfigurationRequestInput{})
}
func GetMasterConfigurationRespFromWire(b []byte) (*master_pb.GetMasterConfigurationResponse, error) {
	v, err := masterwire.WrapGetMasterConfigurationResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &master_pb.GetMasterConfigurationResponse{
		MetricsAddress:          v.MetricsAddress(),
		MetricsIntervalSeconds:  v.MetricsIntervalSeconds(),
		DefaultReplication:      v.DefaultReplication(),
		Leader:                  v.Leader(),
		VolumeSizeLimitMB:       v.VolumeSizeLimitMB(),
		VolumePreallocate:       v.VolumePreallocate(),
		MaintenanceScripts:      v.MaintenanceScripts(),
		MaintenanceSleepMinutes: v.MaintenanceSleepMinutes(),
	}
	for i := 0; i < v.StorageBackendsLen(); i++ {
		if sb, ok := v.StorageBackendAt(i); ok {
			resp.StorageBackends = append(resp.StorageBackends, storageBackendFromView(sb))
		}
	}
	return resp, nil
}

// ListClusterNodes
func ListClusterNodesReqToWire(r *master_pb.ListClusterNodesRequest) []byte {
	return masterwire.NewListClusterNodesRequest(masterwire.ListClusterNodesRequestInput{
		ClientType: r.ClientType, FilerGroup: r.FilerGroup, Limit: r.Limit,
	})
}
func ListClusterNodesRespFromWire(b []byte) (*master_pb.ListClusterNodesResponse, error) {
	v, err := masterwire.WrapListClusterNodesResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &master_pb.ListClusterNodesResponse{}
	for i := 0; i < v.ClusterNodesLen(); i++ {
		if n, ok := v.ClusterNodeAt(i); ok {
			resp.ClusterNodes = append(resp.ClusterNodes, &master_pb.ListClusterNodesResponse_ClusterNode{
				Address: n.Address(), Version: n.Version(), CreatedAtNs: n.CreatedAtNs(),
				DataCenter: n.DataCenter(), Rack: n.Rack(),
			})
		}
	}
	return resp, nil
}

// LeaseAdminToken
func LeaseAdminTokenReqToWire(r *master_pb.LeaseAdminTokenRequest) []byte {
	return masterwire.NewLeaseAdminTokenRequest(masterwire.LeaseAdminTokenRequestInput{
		PreviousToken: r.PreviousToken, PreviousLockTime: r.PreviousLockTime,
		LockName: r.LockName, ClientName: r.ClientName, Message: r.Message,
	})
}
func LeaseAdminTokenRespFromWire(b []byte) (*master_pb.LeaseAdminTokenResponse, error) {
	v, err := masterwire.WrapLeaseAdminTokenResponse(b)
	if err != nil {
		return nil, err
	}
	return &master_pb.LeaseAdminTokenResponse{Token: v.Token(), LockTsNs: v.LockTsNs()}, nil
}

// ReleaseAdminToken
func ReleaseAdminTokenReqToWire(r *master_pb.ReleaseAdminTokenRequest) []byte {
	return masterwire.NewReleaseAdminTokenRequest(masterwire.ReleaseAdminTokenRequestInput{
		PreviousToken: r.PreviousToken, PreviousLockTime: r.PreviousLockTime, LockName: r.LockName,
	})
}
func ReleaseAdminTokenRespFromWire(b []byte) (*master_pb.ReleaseAdminTokenResponse, error) {
	if _, err := masterwire.WrapReleaseAdminTokenResponse(b); err != nil {
		return nil, err
	}
	return &master_pb.ReleaseAdminTokenResponse{}, nil
}

// Ping
func PingReqToWire(r *master_pb.PingRequest) []byte {
	return masterwire.NewPingRequest(masterwire.PingRequestInput{Target: r.Target, TargetType: r.TargetType})
}
func PingRespFromWire(b []byte) (*master_pb.PingResponse, error) {
	v, err := masterwire.WrapPingResponse(b)
	if err != nil {
		return nil, err
	}
	return &master_pb.PingResponse{StartTimeNs: v.StartTimeNs(), RemoteTimeNs: v.RemoteTimeNs(), StopTimeNs: v.StopTimeNs()}, nil
}

// RaftListClusterServers
func RaftListClusterServersReqToWire(r *master_pb.RaftListClusterServersRequest) []byte {
	return masterwire.NewRaftListClusterServersRequest(masterwire.RaftListClusterServersRequestInput{})
}
func RaftListClusterServersRespFromWire(b []byte) (*master_pb.RaftListClusterServersResponse, error) {
	v, err := masterwire.WrapRaftListClusterServersResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &master_pb.RaftListClusterServersResponse{}
	for i := 0; i < v.ClusterServersLen(); i++ {
		if s, ok := v.ClusterServerAt(i); ok {
			resp.ClusterServers = append(resp.ClusterServers, &master_pb.RaftListClusterServersResponse_ClusterServers{
				Id: s.Id(), Address: s.Address(), Suffrage: s.Suffrage(), IsLeader: s.IsLeader(),
			})
		}
	}
	return resp, nil
}

// RaftAddServer
func RaftAddServerReqToWire(r *master_pb.RaftAddServerRequest) []byte {
	return masterwire.NewRaftAddServerRequest(masterwire.RaftAddServerRequestInput{
		Id: r.Id, Address: r.Address, Voter: r.Voter,
	})
}
func RaftAddServerRespFromWire(b []byte) (*master_pb.RaftAddServerResponse, error) {
	if _, err := masterwire.WrapRaftAddServerResponse(b); err != nil {
		return nil, err
	}
	return &master_pb.RaftAddServerResponse{}, nil
}

// RaftRemoveServer
func RaftRemoveServerReqToWire(r *master_pb.RaftRemoveServerRequest) []byte {
	return masterwire.NewRaftRemoveServerRequest(masterwire.RaftRemoveServerRequestInput{Id: r.Id, Force: r.Force})
}
func RaftRemoveServerRespFromWire(b []byte) (*master_pb.RaftRemoveServerResponse, error) {
	if _, err := masterwire.WrapRaftRemoveServerResponse(b); err != nil {
		return nil, err
	}
	return &master_pb.RaftRemoveServerResponse{}, nil
}

// RaftLeadershipTransfer
func RaftLeadershipTransferReqToWire(r *master_pb.RaftLeadershipTransferRequest) []byte {
	return masterwire.NewRaftLeadershipTransferRequest(masterwire.RaftLeadershipTransferRequestInput{
		TargetId: r.TargetId, TargetAddress: r.TargetAddress,
	})
}
func RaftLeadershipTransferRespFromWire(b []byte) (*master_pb.RaftLeadershipTransferResponse, error) {
	v, err := masterwire.WrapRaftLeadershipTransferResponse(b)
	if err != nil {
		return nil, err
	}
	return &master_pb.RaftLeadershipTransferResponse{PreviousLeader: v.PreviousLeader(), NewLeader: v.NewLeader()}, nil
}

// VolumeGrow
func VolumeGrowReqToWire(r *master_pb.VolumeGrowRequest) []byte {
	return masterwire.NewVolumeGrowRequest(masterwire.VolumeGrowRequestInput{
		WritableVolumeCount: r.WritableVolumeCount,
		Replication:         r.Replication,
		Collection:          r.Collection,
		Ttl:                 r.Ttl,
		DataCenter:          r.DataCenter,
		Rack:                r.Rack,
		DataNode:            r.DataNode,
		MemoryMapMaxSizeMb:  r.MemoryMapMaxSizeMb,
		DiskType:            r.DiskType,
	})
}
func VolumeGrowRespFromWire(b []byte) (*master_pb.VolumeGrowResponse, error) {
	if _, err := masterwire.WrapVolumeGrowResponse(b); err != nil {
		return nil, err
	}
	return &master_pb.VolumeGrowResponse{}, nil
}

// --- server-direction unary converters --------------------------------------
//
// The bridge (masterZapBridge) decodes each request from its ZAP buffer with
// <Rpc>ReqFromWire, runs the unchanged (ctx,*pb.Req)->(*pb.Resp,err) handler,
// then re-encodes the reply with <Rpc>RespToWire. These are the inverses of the
// client-direction ReqToWire/RespFromWire above; the nested-tree heavy lifting
// (topology, locations, storage backends) reuses the shared helpers in
// convert.go, which are already bidirectional.

// LookupVolume
func LookupVolumeReqFromWire(b []byte) (*master_pb.LookupVolumeRequest, error) {
	v, err := masterwire.WrapLookupVolumeRequest(b)
	if err != nil {
		return nil, err
	}
	req := &master_pb.LookupVolumeRequest{Collection: v.Collection()}
	for i := 0; i < v.VolumeOrFileIdsLen(); i++ {
		req.VolumeOrFileIds = append(req.VolumeOrFileIds, v.VolumeOrFileIdAt(i))
	}
	return req, nil
}
func LookupVolumeRespToWire(r *master_pb.LookupVolumeResponse) []byte {
	var locs [][]byte
	for _, l := range r.VolumeIdLocations {
		var inner [][]byte
		for _, loc := range l.Locations {
			inner = append(inner, locationToWire(loc))
		}
		locs = append(locs, masterwire.NewLookupVolumeResponseVolumeIdLocation(masterwire.LookupVolumeResponseVolumeIdLocationInput{
			VolumeOrFileId: l.VolumeOrFileId, Locations: inner, Error: l.Error, Auth: l.Auth,
		}))
	}
	return masterwire.NewLookupVolumeResponse(masterwire.LookupVolumeResponseInput{VolumeIdLocations: locs})
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
	var sls [][]byte
	for _, s := range r.ShardIdLocations {
		var inner [][]byte
		for _, loc := range s.Locations {
			inner = append(inner, locationToWire(loc))
		}
		sls = append(sls, masterwire.NewLookupEcVolumeResponseEcShardIdLocation(masterwire.LookupEcVolumeResponseEcShardIdLocationInput{
			ShardId: s.ShardId, Locations: inner,
		}))
	}
	return masterwire.NewLookupEcVolumeResponse(masterwire.LookupEcVolumeResponseInput{
		VolumeId: r.VolumeId, ShardIdLocations: sls,
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
	var cols [][]byte
	for _, c := range r.Collections {
		cols = append(cols, masterwire.NewCollection(masterwire.CollectionInput{Name: c.Name}))
	}
	return masterwire.NewCollectionListResponse(masterwire.CollectionListResponseInput{Collections: cols})
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
	var nodes [][]byte
	for _, n := range r.ClusterNodes {
		nodes = append(nodes, masterwire.NewListClusterNodesResponseClusterNode(masterwire.ListClusterNodesResponseClusterNodeInput{
			Address: n.Address, Version: n.Version, CreatedAtNs: n.CreatedAtNs,
			DataCenter: n.DataCenter, Rack: n.Rack,
		}))
	}
	return masterwire.NewListClusterNodesResponse(masterwire.ListClusterNodesResponseInput{ClusterNodes: nodes})
}

// --- server-direction raft converters ---------------------------------------
//
// The 4 raft membership RPCs answer over the SAME ZAP unary Backend as the rest
// of the master service (the gRPC raft listener carries only the seaweedfs/raft
// transport, not these admin RPCs). Inverses of the client-direction readers
// above.

func RaftListClusterServersReqFromWire(b []byte) (*master_pb.RaftListClusterServersRequest, error) {
	if _, err := masterwire.WrapRaftListClusterServersRequest(b); err != nil {
		return nil, err
	}
	return &master_pb.RaftListClusterServersRequest{}, nil
}
func RaftListClusterServersRespToWire(r *master_pb.RaftListClusterServersResponse) []byte {
	var servers [][]byte
	for _, s := range r.ClusterServers {
		servers = append(servers, masterwire.NewRaftListClusterServersResponseClusterServers(masterwire.RaftListClusterServersResponseClusterServersInput{
			Id: s.Id, Address: s.Address, Suffrage: s.Suffrage, IsLeader: s.IsLeader,
		}))
	}
	return masterwire.NewRaftListClusterServersResponse(masterwire.RaftListClusterServersResponseInput{ClusterServers: servers})
}

func RaftAddServerReqFromWire(b []byte) (*master_pb.RaftAddServerRequest, error) {
	v, err := masterwire.WrapRaftAddServerRequest(b)
	if err != nil {
		return nil, err
	}
	return &master_pb.RaftAddServerRequest{Id: v.Id(), Address: v.Address(), Voter: v.Voter()}, nil
}
func RaftAddServerRespToWire(r *master_pb.RaftAddServerResponse) []byte {
	return masterwire.NewRaftAddServerResponse(masterwire.RaftAddServerResponseInput{})
}

func RaftRemoveServerReqFromWire(b []byte) (*master_pb.RaftRemoveServerRequest, error) {
	v, err := masterwire.WrapRaftRemoveServerRequest(b)
	if err != nil {
		return nil, err
	}
	return &master_pb.RaftRemoveServerRequest{Id: v.Id(), Force: v.Force()}, nil
}
func RaftRemoveServerRespToWire(r *master_pb.RaftRemoveServerResponse) []byte {
	return masterwire.NewRaftRemoveServerResponse(masterwire.RaftRemoveServerResponseInput{})
}

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
