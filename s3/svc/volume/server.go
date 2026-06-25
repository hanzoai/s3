// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// server.go is the server-side mirror of client.go: a vsw.VolumeServerStore that
// wraps the existing gRPC-shaped volume server (volume_server_pb.VolumeServerServer)
// so it answers over the canonical github.com/zap-proto/go transport. Each unary
// RPC is one hop each way: wire request buffer -> <Rpc>ReqFromWire -> the existing
// VolumeServerServer method -> <Rpc>RespToWire -> wire response buffer. The 11
// streaming RPCs are served in stream_server.go.
//
// This is the volume analogue of filer.NewServerBackend: with NewServerBackend
// wired into vsw.Dispatch + vsw.StreamHandler (the analogue of the gRPC
// RegisterVolumeServerServer), the ZAP client in client.go reaches the real
// volume engine over ZAP — no protobuf framing on the wire, all engine logic
// reused unchanged. The bytes returned ARE the message.

package volume

import (
	"context"

	volume_server_pb "github.com/hanzoai/s3/s3/pb/volume_server_pb"
	vsw "github.com/hanzoai/s3/s3/wire/volume_server"
)

// serverBackend adapts a volume_server_pb.VolumeServerServer to vsw.VolumeServerStore.
type serverBackend struct {
	vs  volume_server_pb.VolumeServerServer
	ctx context.Context
}

// NewServerBackend returns a vsw.VolumeServerStore that serves vs over ZAP. Pass
// it to vsw.Dispatch (unary) and vsw.StreamHandler (streaming), or to vsw.Serve /
// transport.ListenStream.
func NewServerBackend(vs volume_server_pb.VolumeServerServer) vsw.VolumeServerStore {
	return serverBackend{vs: vs, ctx: context.Background()}
}

var _ vsw.VolumeServerStore = serverBackend{}

// --- unary RPCs (37): decode req buffer -> engine method -> encode resp buffer ---

func (b serverBackend) BatchDelete(req []byte) ([]byte, error) {
	in, err := BatchDeleteReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.vs.BatchDelete(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return BatchDeleteRespToWire(resp), nil
}

func (b serverBackend) VacuumVolumeCheck(req []byte) ([]byte, error) {
	in, err := VacuumVolumeCheckReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.vs.VacuumVolumeCheck(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return VacuumVolumeCheckRespToWire(resp), nil
}

func (b serverBackend) VacuumVolumeCommit(req []byte) ([]byte, error) {
	in, err := VacuumVolumeCommitReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.vs.VacuumVolumeCommit(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return VacuumVolumeCommitRespToWire(resp), nil
}

func (b serverBackend) VacuumVolumeCleanup(req []byte) ([]byte, error) {
	in, err := VacuumVolumeCleanupReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.vs.VacuumVolumeCleanup(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return VacuumVolumeCleanupRespToWire(resp), nil
}

func (b serverBackend) DeleteCollection(req []byte) ([]byte, error) {
	in, err := DeleteCollectionReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.vs.DeleteCollection(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return DeleteCollectionRespToWire(resp), nil
}

func (b serverBackend) AllocateVolume(req []byte) ([]byte, error) {
	in, err := AllocateVolumeReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.vs.AllocateVolume(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return AllocateVolumeRespToWire(resp), nil
}

func (b serverBackend) VolumeSyncStatus(req []byte) ([]byte, error) {
	in, err := VolumeSyncStatusReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.vs.VolumeSyncStatus(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return VolumeSyncStatusRespToWire(resp), nil
}

func (b serverBackend) VolumeMount(req []byte) ([]byte, error) {
	in, err := VolumeMountReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.vs.VolumeMount(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return VolumeMountRespToWire(resp), nil
}

func (b serverBackend) VolumeUnmount(req []byte) ([]byte, error) {
	in, err := VolumeUnmountReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.vs.VolumeUnmount(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return VolumeUnmountRespToWire(resp), nil
}

func (b serverBackend) VolumeDelete(req []byte) ([]byte, error) {
	in, err := VolumeDeleteReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.vs.VolumeDelete(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return VolumeDeleteRespToWire(resp), nil
}

func (b serverBackend) VolumeMarkReadonly(req []byte) ([]byte, error) {
	in, err := VolumeMarkReadonlyReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.vs.VolumeMarkReadonly(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return VolumeMarkReadonlyRespToWire(resp), nil
}

func (b serverBackend) VolumeMarkWritable(req []byte) ([]byte, error) {
	in, err := VolumeMarkWritableReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.vs.VolumeMarkWritable(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return VolumeMarkWritableRespToWire(resp), nil
}

func (b serverBackend) VolumeConfigure(req []byte) ([]byte, error) {
	in, err := VolumeConfigureReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.vs.VolumeConfigure(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return VolumeConfigureRespToWire(resp), nil
}

func (b serverBackend) VolumeStatus(req []byte) ([]byte, error) {
	in, err := VolumeStatusReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.vs.VolumeStatus(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return VolumeStatusRespToWire(resp), nil
}

func (b serverBackend) GetState(req []byte) ([]byte, error) {
	in, err := GetStateReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.vs.GetState(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return GetStateRespToWire(resp), nil
}

func (b serverBackend) SetState(req []byte) ([]byte, error) {
	in, err := SetStateReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.vs.SetState(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return SetStateRespToWire(resp), nil
}

func (b serverBackend) ReadVolumeFileStatus(req []byte) ([]byte, error) {
	in, err := ReadVolumeFileStatusReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.vs.ReadVolumeFileStatus(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return ReadVolumeFileStatusRespToWire(resp), nil
}

func (b serverBackend) ReadNeedleBlob(req []byte) ([]byte, error) {
	in, err := ReadNeedleBlobReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.vs.ReadNeedleBlob(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return ReadNeedleBlobRespToWire(resp), nil
}

func (b serverBackend) ReadNeedleMeta(req []byte) ([]byte, error) {
	in, err := ReadNeedleMetaReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.vs.ReadNeedleMeta(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return ReadNeedleMetaRespToWire(resp), nil
}

func (b serverBackend) WriteNeedleBlob(req []byte) ([]byte, error) {
	in, err := WriteNeedleBlobReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.vs.WriteNeedleBlob(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return WriteNeedleBlobRespToWire(resp), nil
}

func (b serverBackend) VolumeTailReceiver(req []byte) ([]byte, error) {
	in, err := VolumeTailReceiverReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.vs.VolumeTailReceiver(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return VolumeTailReceiverRespToWire(resp), nil
}

func (b serverBackend) VolumeEcShardsGenerate(req []byte) ([]byte, error) {
	in, err := VolumeEcShardsGenerateReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.vs.VolumeEcShardsGenerate(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return VolumeEcShardsGenerateRespToWire(resp), nil
}

func (b serverBackend) VolumeEcShardsRebuild(req []byte) ([]byte, error) {
	in, err := VolumeEcShardsRebuildReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.vs.VolumeEcShardsRebuild(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return VolumeEcShardsRebuildRespToWire(resp), nil
}

func (b serverBackend) VolumeEcShardsCopy(req []byte) ([]byte, error) {
	in, err := VolumeEcShardsCopyReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.vs.VolumeEcShardsCopy(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return VolumeEcShardsCopyRespToWire(resp), nil
}

func (b serverBackend) VolumeEcShardsDelete(req []byte) ([]byte, error) {
	in, err := VolumeEcShardsDeleteReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.vs.VolumeEcShardsDelete(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return VolumeEcShardsDeleteRespToWire(resp), nil
}

func (b serverBackend) VolumeEcShardsMount(req []byte) ([]byte, error) {
	in, err := VolumeEcShardsMountReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.vs.VolumeEcShardsMount(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return VolumeEcShardsMountRespToWire(resp), nil
}

func (b serverBackend) VolumeEcShardsUnmount(req []byte) ([]byte, error) {
	in, err := VolumeEcShardsUnmountReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.vs.VolumeEcShardsUnmount(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return VolumeEcShardsUnmountRespToWire(resp), nil
}

func (b serverBackend) VolumeEcBlobDelete(req []byte) ([]byte, error) {
	in, err := VolumeEcBlobDeleteReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.vs.VolumeEcBlobDelete(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return VolumeEcBlobDeleteRespToWire(resp), nil
}

func (b serverBackend) VolumeEcShardsToVolume(req []byte) ([]byte, error) {
	in, err := VolumeEcShardsToVolumeReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.vs.VolumeEcShardsToVolume(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return VolumeEcShardsToVolumeRespToWire(resp), nil
}

func (b serverBackend) VolumeEcShardsInfo(req []byte) ([]byte, error) {
	in, err := VolumeEcShardsInfoReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.vs.VolumeEcShardsInfo(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return VolumeEcShardsInfoRespToWire(resp), nil
}

func (b serverBackend) VolumeServerStatus(req []byte) ([]byte, error) {
	in, err := VolumeServerStatusReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.vs.VolumeServerStatus(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return VolumeServerStatusRespToWire(resp), nil
}

func (b serverBackend) VolumeServerLeave(req []byte) ([]byte, error) {
	in, err := VolumeServerLeaveReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.vs.VolumeServerLeave(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return VolumeServerLeaveRespToWire(resp), nil
}

func (b serverBackend) FetchAndWriteNeedle(req []byte) ([]byte, error) {
	in, err := FetchAndWriteNeedleReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.vs.FetchAndWriteNeedle(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return FetchAndWriteNeedleRespToWire(resp), nil
}

func (b serverBackend) ScrubVolume(req []byte) ([]byte, error) {
	in, err := ScrubVolumeReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.vs.ScrubVolume(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return ScrubVolumeRespToWire(resp), nil
}

func (b serverBackend) ScrubEcVolume(req []byte) ([]byte, error) {
	in, err := ScrubEcVolumeReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.vs.ScrubEcVolume(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return ScrubEcVolumeRespToWire(resp), nil
}

func (b serverBackend) VolumeNeedleStatus(req []byte) ([]byte, error) {
	in, err := VolumeNeedleStatusReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.vs.VolumeNeedleStatus(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return VolumeNeedleStatusRespToWire(resp), nil
}

func (b serverBackend) Ping(req []byte) ([]byte, error) {
	in, err := PingReqFromWire(req)
	if err != nil {
		return nil, err
	}
	resp, err := b.vs.Ping(b.ctx, in)
	if err != nil {
		return nil, err
	}
	return PingRespToWire(resp), nil
}
