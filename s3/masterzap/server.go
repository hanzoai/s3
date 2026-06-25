// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// server.go is the server-side mirror of client.go: a masterwire.Backend that
// wraps the existing gRPC-shaped master (master_pb.HanzoServer — the live
// *MasterServer engine) so it answers over the canonical
// github.com/zap-proto/go transport. Each unary RPC is one hop each way: wire
// request buffer -> <Rpc>ReqFromWire -> the existing master_pb.HanzoServer
// method -> <Rpc>RespToWire -> wire response bytes (see server_rpc.go). This is
// the master analogue of filerzap.NewServerBackend: with NewServerBackend wired
// into masterwire.Dispatch (the analogue of the retired gRPC master registration), the
// ZAP client in client.go reaches the real master engine over ZAP — no protobuf
// framing on the wire, all topology/raft/vacuum logic reused unchanged.
//
// The three streaming RPCs (SendHeartbeat, KeepConnected, StreamAssign) are NOT
// served on this unary path; they ride the duplex transport stream via
// NewStreamServer (stream_server.go) and masterstream.Handler. They are present
// here only to satisfy masterwire.Backend, and they fail loudly
// (ErrStreamingNotWired) — masterwire.Dispatch short-circuits them before they
// are reached, so a misrouted streaming call never silently degrades to a
// one-shot exchange.

package masterzap

import (
	"context"

	master_pb "github.com/hanzoai/s3/s3/pb/master_pb"
	masterwire "github.com/hanzoai/s3/s3/wire/master"
)

// serverBackend adapts a master_pb.HanzoServer to masterwire.Backend.
type serverBackend struct {
	ms  master_pb.HanzoServer
	ctx context.Context
}

// NewServerBackend returns a masterwire.Backend that serves ms's unary RPCs over
// ZAP. Pass it to masterwire.Dispatch / transport.ListenStream. The streaming
// RPCs are served separately by NewStreamServer.
func NewServerBackend(ms master_pb.HanzoServer) masterwire.Backend {
	return serverBackend{ms: ms, ctx: context.Background()}
}

// --- streaming (served by NewStreamServer; Dispatch never reaches these) ---

func (b serverBackend) SendHeartbeat([]byte) ([]byte, error) {
	return nil, masterwire.ErrStreamingNotWired
}
func (b serverBackend) KeepConnected([]byte) ([]byte, error) {
	return nil, masterwire.ErrStreamingNotWired
}
func (b serverBackend) StreamAssign([]byte) ([]byte, error) {
	return nil, masterwire.ErrStreamingNotWired
}

// --- unary RPCs (decode request -> engine -> encode response) ---

func (b serverBackend) LookupVolume(req []byte) ([]byte, error) {
	in, err := LookupVolumeReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.ms.LookupVolume(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return LookupVolumeRespToWire(resp), nil
}

func (b serverBackend) Assign(req []byte) ([]byte, error) {
	in, err := AssignReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.ms.Assign(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return AssignRespToWire(resp), nil
}

func (b serverBackend) Statistics(req []byte) ([]byte, error) {
	in, err := StatisticsReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.ms.Statistics(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return StatisticsRespToWire(resp), nil
}

func (b serverBackend) CollectionList(req []byte) ([]byte, error) {
	in, err := CollectionListReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.ms.CollectionList(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return CollectionListRespToWire(resp), nil
}

func (b serverBackend) CollectionDelete(req []byte) ([]byte, error) {
	in, err := CollectionDeleteReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.ms.CollectionDelete(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return CollectionDeleteRespToWire(resp), nil
}

func (b serverBackend) VolumeList(req []byte) ([]byte, error) {
	in, err := VolumeListReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.ms.VolumeList(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return VolumeListRespToWire(resp), nil
}

func (b serverBackend) LookupEcVolume(req []byte) ([]byte, error) {
	in, err := LookupEcVolumeReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.ms.LookupEcVolume(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return LookupEcVolumeRespToWire(resp), nil
}

func (b serverBackend) VacuumVolume(req []byte) ([]byte, error) {
	in, err := VacuumVolumeReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.ms.VacuumVolume(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return VacuumVolumeRespToWire(resp), nil
}

func (b serverBackend) DisableVacuum(req []byte) ([]byte, error) {
	in, err := DisableVacuumReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.ms.DisableVacuum(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return DisableVacuumRespToWire(resp), nil
}

func (b serverBackend) EnableVacuum(req []byte) ([]byte, error) {
	in, err := EnableVacuumReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.ms.EnableVacuum(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return EnableVacuumRespToWire(resp), nil
}

func (b serverBackend) VolumeMarkReadonly(req []byte) ([]byte, error) {
	in, err := VolumeMarkReadonlyReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.ms.VolumeMarkReadonly(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return VolumeMarkReadonlyRespToWire(resp), nil
}

func (b serverBackend) GetMasterConfiguration(req []byte) ([]byte, error) {
	in, err := GetMasterConfigurationReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.ms.GetMasterConfiguration(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return GetMasterConfigurationRespToWire(resp), nil
}

func (b serverBackend) ListClusterNodes(req []byte) ([]byte, error) {
	in, err := ListClusterNodesReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.ms.ListClusterNodes(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return ListClusterNodesRespToWire(resp), nil
}

func (b serverBackend) LeaseAdminToken(req []byte) ([]byte, error) {
	in, err := LeaseAdminTokenReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.ms.LeaseAdminToken(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return LeaseAdminTokenRespToWire(resp), nil
}

func (b serverBackend) ReleaseAdminToken(req []byte) ([]byte, error) {
	in, err := ReleaseAdminTokenReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.ms.ReleaseAdminToken(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return ReleaseAdminTokenRespToWire(resp), nil
}

func (b serverBackend) Ping(req []byte) ([]byte, error) {
	in, err := PingReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.ms.Ping(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return PingRespToWire(resp), nil
}

func (b serverBackend) RaftListClusterServers(req []byte) ([]byte, error) {
	in, err := RaftListClusterServersReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.ms.RaftListClusterServers(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return RaftListClusterServersRespToWire(resp), nil
}

func (b serverBackend) RaftAddServer(req []byte) ([]byte, error) {
	in, err := RaftAddServerReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.ms.RaftAddServer(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return RaftAddServerRespToWire(resp), nil
}

func (b serverBackend) RaftRemoveServer(req []byte) ([]byte, error) {
	in, err := RaftRemoveServerReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.ms.RaftRemoveServer(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return RaftRemoveServerRespToWire(resp), nil
}

func (b serverBackend) RaftLeadershipTransfer(req []byte) ([]byte, error) {
	in, err := RaftLeadershipTransferReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.ms.RaftLeadershipTransfer(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return RaftLeadershipTransferRespToWire(resp), nil
}

func (b serverBackend) VolumeGrow(req []byte) ([]byte, error) {
	in, err := VolumeGrowReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.ms.VolumeGrow(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return VolumeGrowRespToWire(resp), nil
}
