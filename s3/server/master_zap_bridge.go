// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// masterZapBridge adapts the existing *MasterServer handlers to the native ZAP
// masterwire.Backend contract, so the master serves its unary RPCs over the
// canonical zap-proto transport. Each unary RPC is decoded from its ZAP request
// buffer (masterzap.<Rpc>ReqFromWire), threaded through the unchanged
// (ctx,*pb.Req)->(*pb.Resp,err) handler, then re-encoded as the ZAP response
// (masterzap.<Rpc>RespToWire). The three bidi streams are served separately by
// masterZapStream (master_zap_stream.go) over the same listener.
//
// NOT bridged here:
//   - raft membership: RaftListClusterServers/AddServer/RemoveServer/
//     LeadershipTransfer (entangled with the raftServer; consensus migration is
//     separate — see consensus_server.go). These ride the gRPC raft listener.
//   - streaming dispatch: SendHeartbeat/KeepConnected/StreamAssign are served by
//     masterstream.Handler, NOT this unary Backend; the methods below exist only
//     to satisfy masterwire.Backend and must never be reached over the unary
//     path.
package s3server

import (
	"context"
	"errors"

	"github.com/hanzoai/s3/s3/masterzap"
	"github.com/hanzoai/s3/s3/pb/master_pb"
	masterwire "github.com/hanzoai/s3/s3/wire/master"
)

// errZapUnaryStream guards the streaming ordinals on the unary Backend: those
// RPCs are served by masterstream.Handler, so reaching them here is a wiring
// bug, not a runtime path.
var errZapUnaryStream = errors.New("masterzap: streaming RPC must be served via masterstream.Handler, not the unary Backend")

// masterZapBridge implements masterwire.Backend by delegating to the live
// *MasterServer handlers.
type masterZapBridge struct{ ms *MasterServer }

// NewMasterZapBackend returns the masterwire.Backend backed by ms.
func NewMasterZapBackend(ms *MasterServer) masterwire.Backend { return masterZapBridge{ms: ms} }

// --- streaming (served by masterstream.Handler, never the unary path) ---

func (z masterZapBridge) SendHeartbeat(req []byte) ([]byte, error) { return nil, errZapUnaryStream }
func (z masterZapBridge) KeepConnected(req []byte) ([]byte, error) { return nil, errZapUnaryStream }
func (z masterZapBridge) StreamAssign(req []byte) ([]byte, error)  { return nil, errZapUnaryStream }

// --- nested-message unary ---

func (z masterZapBridge) LookupVolume(req []byte) ([]byte, error) {
	in, err := masterzap.LookupVolumeReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := z.ms.LookupVolume(context.Background(), in)
	if err != nil {
		return nil, err
	}
	return masterzap.LookupVolumeRespToWire(resp), nil
}

func (z masterZapBridge) Assign(req []byte) ([]byte, error) {
	in, err := masterzap.AssignReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := z.ms.Assign(context.Background(), in)
	if err != nil {
		return nil, err
	}
	return masterzap.AssignRespToWire(resp), nil
}

func (z masterZapBridge) VolumeList(req []byte) ([]byte, error) {
	in, err := masterzap.VolumeListReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := z.ms.VolumeList(context.Background(), in)
	if err != nil {
		return nil, err
	}
	return masterzap.VolumeListRespToWire(resp), nil
}

func (z masterZapBridge) LookupEcVolume(req []byte) ([]byte, error) {
	in, err := masterzap.LookupEcVolumeReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := z.ms.LookupEcVolume(context.Background(), in)
	if err != nil {
		return nil, err
	}
	return masterzap.LookupEcVolumeRespToWire(resp), nil
}

func (z masterZapBridge) CollectionList(req []byte) ([]byte, error) {
	in, err := masterzap.CollectionListReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := z.ms.CollectionList(context.Background(), in)
	if err != nil {
		return nil, err
	}
	return masterzap.CollectionListRespToWire(resp), nil
}

func (z masterZapBridge) GetMasterConfiguration(req []byte) ([]byte, error) {
	in, err := masterzap.GetMasterConfigurationReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := z.ms.GetMasterConfiguration(context.Background(), in)
	if err != nil {
		return nil, err
	}
	return masterzap.GetMasterConfigurationRespToWire(resp), nil
}

func (z masterZapBridge) ListClusterNodes(req []byte) ([]byte, error) {
	in, err := masterzap.ListClusterNodesReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := z.ms.ListClusterNodes(context.Background(), in)
	if err != nil {
		return nil, err
	}
	return masterzap.ListClusterNodesRespToWire(resp), nil
}

// --- raft membership admin RPCs (the gRPC listener carries only the
//     seaweedfs/raft transport; these admin calls answer over ZAP) ---

func (z masterZapBridge) RaftListClusterServers(req []byte) ([]byte, error) {
	in, err := masterzap.RaftListClusterServersReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := z.ms.RaftListClusterServers(context.Background(), in)
	if err != nil {
		return nil, err
	}
	return masterzap.RaftListClusterServersRespToWire(resp), nil
}

func (z masterZapBridge) RaftAddServer(req []byte) ([]byte, error) {
	in, err := masterzap.RaftAddServerReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := z.ms.RaftAddServer(context.Background(), in)
	if err != nil {
		return nil, err
	}
	return masterzap.RaftAddServerRespToWire(resp), nil
}

func (z masterZapBridge) RaftRemoveServer(req []byte) ([]byte, error) {
	in, err := masterzap.RaftRemoveServerReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := z.ms.RaftRemoveServer(context.Background(), in)
	if err != nil {
		return nil, err
	}
	return masterzap.RaftRemoveServerRespToWire(resp), nil
}

func (z masterZapBridge) RaftLeadershipTransfer(req []byte) ([]byte, error) {
	in, err := masterzap.RaftLeadershipTransferReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := z.ms.RaftLeadershipTransfer(context.Background(), in)
	if err != nil {
		return nil, err
	}
	return masterzap.RaftLeadershipTransferRespToWire(resp), nil
}

// --- scalar-shaped unary (bridged) ---

func (z masterZapBridge) Statistics(req []byte) ([]byte, error) {
	in, err := masterwire.WrapStatisticsRequest(req)
	if err != nil {
		return nil, err
	}
	resp, err := z.ms.Statistics(context.Background(), &master_pb.StatisticsRequest{
		Replication: in.Replication(),
		Collection:  in.Collection(),
		Ttl:         in.Ttl(),
		DiskType:    in.DiskType(),
	})
	if err != nil {
		return nil, err
	}
	return masterwire.NewStatisticsResponse(masterwire.StatisticsResponseInput{
		TotalSize: resp.TotalSize,
		UsedSize:  resp.UsedSize,
		FileCount: resp.FileCount,
	}), nil
}

func (z masterZapBridge) VacuumVolume(req []byte) ([]byte, error) {
	in, err := masterwire.WrapVacuumVolumeRequest(req)
	if err != nil {
		return nil, err
	}
	if _, err := z.ms.VacuumVolume(context.Background(), &master_pb.VacuumVolumeRequest{
		GarbageThreshold: in.GarbageThreshold(),
		VolumeId:         in.VolumeId(),
		Collection:       in.Collection(),
	}); err != nil {
		return nil, err
	}
	return masterwire.NewVacuumVolumeResponse(masterwire.VacuumVolumeResponseInput{}), nil
}

func (z masterZapBridge) DisableVacuum(req []byte) ([]byte, error) {
	in, err := masterwire.WrapDisableVacuumRequest(req)
	if err != nil {
		return nil, err
	}
	if _, err := z.ms.DisableVacuum(context.Background(), &master_pb.DisableVacuumRequest{
		ByPlugin: in.ByPlugin(),
	}); err != nil {
		return nil, err
	}
	return masterwire.NewDisableVacuumResponse(masterwire.DisableVacuumResponseInput{}), nil
}

func (z masterZapBridge) EnableVacuum(req []byte) ([]byte, error) {
	in, err := masterwire.WrapEnableVacuumRequest(req)
	if err != nil {
		return nil, err
	}
	if _, err := z.ms.EnableVacuum(context.Background(), &master_pb.EnableVacuumRequest{
		ByPlugin: in.ByPlugin(),
	}); err != nil {
		return nil, err
	}
	return masterwire.NewEnableVacuumResponse(masterwire.EnableVacuumResponseInput{}), nil
}

func (z masterZapBridge) VolumeMarkReadonly(req []byte) ([]byte, error) {
	in, err := masterwire.WrapVolumeMarkReadonlyRequest(req)
	if err != nil {
		return nil, err
	}
	if _, err := z.ms.VolumeMarkReadonly(context.Background(), &master_pb.VolumeMarkReadonlyRequest{
		Ip:               in.Ip(),
		Port:             in.Port(),
		VolumeId:         in.VolumeId(),
		Collection:       in.Collection(),
		ReplicaPlacement: in.ReplicaPlacement(),
		Version:          in.Version(),
		Ttl:              in.Ttl(),
		DiskType:         in.DiskType(),
		IsReadonly:       in.IsReadonly(),
	}); err != nil {
		return nil, err
	}
	return masterwire.NewVolumeMarkReadonlyResponse(masterwire.VolumeMarkReadonlyResponseInput{}), nil
}

func (z masterZapBridge) VolumeGrow(req []byte) ([]byte, error) {
	in, err := masterwire.WrapVolumeGrowRequest(req)
	if err != nil {
		return nil, err
	}
	if _, err := z.ms.VolumeGrow(context.Background(), &master_pb.VolumeGrowRequest{
		WritableVolumeCount: in.WritableVolumeCount(),
		Replication:         in.Replication(),
		Collection:          in.Collection(),
		Ttl:                 in.Ttl(),
		DataCenter:          in.DataCenter(),
		Rack:                in.Rack(),
		DataNode:            in.DataNode(),
		MemoryMapMaxSizeMb:  in.MemoryMapMaxSizeMb(),
		DiskType:            in.DiskType(),
	}); err != nil {
		return nil, err
	}
	return masterwire.NewVolumeGrowResponse(masterwire.VolumeGrowResponseInput{}), nil
}

func (z masterZapBridge) CollectionDelete(req []byte) ([]byte, error) {
	in, err := masterwire.WrapCollectionDeleteRequest(req)
	if err != nil {
		return nil, err
	}
	if _, err := z.ms.CollectionDelete(context.Background(), &master_pb.CollectionDeleteRequest{
		Name: in.Name(),
	}); err != nil {
		return nil, err
	}
	return masterwire.NewCollectionDeleteResponse(masterwire.CollectionDeleteResponseInput{}), nil
}

func (z masterZapBridge) LeaseAdminToken(req []byte) ([]byte, error) {
	in, err := masterwire.WrapLeaseAdminTokenRequest(req)
	if err != nil {
		return nil, err
	}
	resp, err := z.ms.LeaseAdminToken(context.Background(), &master_pb.LeaseAdminTokenRequest{
		PreviousToken:    in.PreviousToken(),
		PreviousLockTime: in.PreviousLockTime(),
		LockName:         in.LockName(),
		ClientName:       in.ClientName(),
		Message:          in.Message(),
	})
	if err != nil {
		return nil, err
	}
	return masterwire.NewLeaseAdminTokenResponse(masterwire.LeaseAdminTokenResponseInput{
		Token:    resp.Token,
		LockTsNs: resp.LockTsNs,
	}), nil
}

func (z masterZapBridge) ReleaseAdminToken(req []byte) ([]byte, error) {
	in, err := masterwire.WrapReleaseAdminTokenRequest(req)
	if err != nil {
		return nil, err
	}
	if _, err := z.ms.ReleaseAdminToken(context.Background(), &master_pb.ReleaseAdminTokenRequest{
		PreviousToken:    in.PreviousToken(),
		PreviousLockTime: in.PreviousLockTime(),
		LockName:         in.LockName(),
	}); err != nil {
		return nil, err
	}
	return masterwire.NewReleaseAdminTokenResponse(masterwire.ReleaseAdminTokenResponseInput{}), nil
}

func (z masterZapBridge) Ping(req []byte) ([]byte, error) {
	in, err := masterwire.WrapPingRequest(req)
	if err != nil {
		return nil, err
	}
	resp, err := z.ms.Ping(context.Background(), &master_pb.PingRequest{
		Target:     in.Target(),
		TargetType: in.TargetType(),
	})
	if err != nil {
		return nil, err
	}
	return masterwire.NewPingResponse(masterwire.PingResponseInput{
		StartTimeNs:  resp.StartTimeNs,
		RemoteTimeNs: resp.RemoteTimeNs,
		StopTimeNs:   resp.StopTimeNs,
	}), nil
}
