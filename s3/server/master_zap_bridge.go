// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// masterZapBridge adapts the existing *MasterServer gRPC handlers to the
// native ZAP masterwire.Backend contract, so the master serves its unary RPCs
// over the canonical zap-proto transport (grpcPort+10000) alongside the legacy
// gRPC listener. This is the strangler bridge: scalar-shaped unary RPCs are
// decoded from their ZAP request buffer, threaded through the unchanged
// (ctx,*pb.Req)->(*pb.Resp,err) handler, then re-encoded as the ZAP response.
//
// NOT bridged here (still served only over gRPC until their wire converters /
// stream transport land — masterwire returns an explicit error so a caller that
// dials ZAP for them fails loudly instead of silently degrading):
//   - nested-topology unary: LookupVolume, Assign, VolumeList, LookupEcVolume,
//     ListClusterNodes, CollectionList, GetMasterConfiguration (carry pb message
//     trees with no masterwire New* counterpart yet).
//   - raft membership: RaftListClusterServers/AddServer/RemoveServer/
//     LeadershipTransfer (entangled with the raftServer; consensus migration is
//     separate — see consensus_server.go).
//   - streaming: SendHeartbeat, KeepConnected, StreamAssign (transport stream
//     primitive not wired for master).
package s3server

import (
	"context"
	"errors"

	"github.com/hanzoai/s3/s3/pb/master_pb"
	masterwire "github.com/hanzoai/s3/s3/wire/master"
)

// errZapNotBridged marks a master RPC that is not yet served over the ZAP
// transport. Clients for these RPCs stay on the gRPC path during the strangler.
var errZapNotBridged = errors.New("masterzap: RPC not bridged to ZAP yet (served over gRPC)")

// masterZapBridge implements masterwire.Backend by delegating to the live
// *MasterServer handlers.
type masterZapBridge struct{ ms *MasterServer }

// NewMasterZapBackend returns the masterwire.Backend backed by ms.
func NewMasterZapBackend(ms *MasterServer) masterwire.Backend { return masterZapBridge{ms: ms} }

// --- streaming (not wired) ---

func (z masterZapBridge) SendHeartbeat(req []byte) ([]byte, error) { return nil, errZapNotBridged }
func (z masterZapBridge) KeepConnected(req []byte) ([]byte, error) { return nil, errZapNotBridged }
func (z masterZapBridge) StreamAssign(req []byte) ([]byte, error)  { return nil, errZapNotBridged }

// --- nested-topology unary (not bridged yet) ---

func (z masterZapBridge) LookupVolume(req []byte) ([]byte, error)   { return nil, errZapNotBridged }
func (z masterZapBridge) Assign(req []byte) ([]byte, error)         { return nil, errZapNotBridged }
func (z masterZapBridge) VolumeList(req []byte) ([]byte, error)     { return nil, errZapNotBridged }
func (z masterZapBridge) LookupEcVolume(req []byte) ([]byte, error) { return nil, errZapNotBridged }
func (z masterZapBridge) CollectionList(req []byte) ([]byte, error) { return nil, errZapNotBridged }
func (z masterZapBridge) GetMasterConfiguration(req []byte) ([]byte, error) {
	return nil, errZapNotBridged
}
func (z masterZapBridge) ListClusterNodes(req []byte) ([]byte, error) { return nil, errZapNotBridged }

// --- raft membership (consensus migration is separate) ---

func (z masterZapBridge) RaftListClusterServers(req []byte) ([]byte, error) {
	return nil, errZapNotBridged
}
func (z masterZapBridge) RaftAddServer(req []byte) ([]byte, error)    { return nil, errZapNotBridged }
func (z masterZapBridge) RaftRemoveServer(req []byte) ([]byte, error) { return nil, errZapNotBridged }
func (z masterZapBridge) RaftLeadershipTransfer(req []byte) ([]byte, error) {
	return nil, errZapNotBridged
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
