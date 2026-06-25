package pluginworkers

import (
	"context"
	"errors"

	"github.com/hanzoai/s3/s3/pb/volume_server_pb"
)

// unimplementedVolumeServer fills the volume_server_pb.VolumeServerServer methods
// the test fake does not exercise, returning errFakeVolumeUnsupported. It is the
// test-local stand-in for the retired grpc UnimplementedVolumeServerServer, so
// the fake satisfies the full grpc-free server contract volume dispatches.
type unimplementedVolumeServer struct{}

var errFakeVolumeUnsupported = errors.New("fake volume server: RPC not implemented")

func (unimplementedVolumeServer) BatchDelete(context.Context, *volume_server_pb.BatchDeleteRequest) (*volume_server_pb.BatchDeleteResponse, error) {
	return nil, errFakeVolumeUnsupported
}

func (unimplementedVolumeServer) DeleteCollection(context.Context, *volume_server_pb.DeleteCollectionRequest) (*volume_server_pb.DeleteCollectionResponse, error) {
	return nil, errFakeVolumeUnsupported
}

func (unimplementedVolumeServer) AllocateVolume(context.Context, *volume_server_pb.AllocateVolumeRequest) (*volume_server_pb.AllocateVolumeResponse, error) {
	return nil, errFakeVolumeUnsupported
}

func (unimplementedVolumeServer) VolumeSyncStatus(context.Context, *volume_server_pb.VolumeSyncStatusRequest) (*volume_server_pb.VolumeSyncStatusResponse, error) {
	return nil, errFakeVolumeUnsupported
}

func (unimplementedVolumeServer) VolumeIncrementalCopy(*volume_server_pb.VolumeIncrementalCopyRequest, volume_server_pb.VolumeServer_VolumeIncrementalCopyServer) error {
	return errFakeVolumeUnsupported
}

func (unimplementedVolumeServer) VolumeUnmount(context.Context, *volume_server_pb.VolumeUnmountRequest) (*volume_server_pb.VolumeUnmountResponse, error) {
	return nil, errFakeVolumeUnsupported
}

func (unimplementedVolumeServer) VolumeConfigure(context.Context, *volume_server_pb.VolumeConfigureRequest) (*volume_server_pb.VolumeConfigureResponse, error) {
	return nil, errFakeVolumeUnsupported
}

func (unimplementedVolumeServer) VolumeStatus(context.Context, *volume_server_pb.VolumeStatusRequest) (*volume_server_pb.VolumeStatusResponse, error) {
	return nil, errFakeVolumeUnsupported
}

func (unimplementedVolumeServer) GetState(context.Context, *volume_server_pb.GetStateRequest) (*volume_server_pb.GetStateResponse, error) {
	return nil, errFakeVolumeUnsupported
}

func (unimplementedVolumeServer) SetState(context.Context, *volume_server_pb.SetStateRequest) (*volume_server_pb.SetStateResponse, error) {
	return nil, errFakeVolumeUnsupported
}

func (unimplementedVolumeServer) ReadNeedleBlob(context.Context, *volume_server_pb.ReadNeedleBlobRequest) (*volume_server_pb.ReadNeedleBlobResponse, error) {
	return nil, errFakeVolumeUnsupported
}

func (unimplementedVolumeServer) ReadNeedleMeta(context.Context, *volume_server_pb.ReadNeedleMetaRequest) (*volume_server_pb.ReadNeedleMetaResponse, error) {
	return nil, errFakeVolumeUnsupported
}

func (unimplementedVolumeServer) WriteNeedleBlob(context.Context, *volume_server_pb.WriteNeedleBlobRequest) (*volume_server_pb.WriteNeedleBlobResponse, error) {
	return nil, errFakeVolumeUnsupported
}

func (unimplementedVolumeServer) ReadAllNeedles(*volume_server_pb.ReadAllNeedlesRequest, volume_server_pb.VolumeServer_ReadAllNeedlesServer) error {
	return errFakeVolumeUnsupported
}

func (unimplementedVolumeServer) VolumeTailSender(*volume_server_pb.VolumeTailSenderRequest, volume_server_pb.VolumeServer_VolumeTailSenderServer) error {
	return errFakeVolumeUnsupported
}

func (unimplementedVolumeServer) VolumeEcShardsGenerate(context.Context, *volume_server_pb.VolumeEcShardsGenerateRequest) (*volume_server_pb.VolumeEcShardsGenerateResponse, error) {
	return nil, errFakeVolumeUnsupported
}

func (unimplementedVolumeServer) VolumeEcShardsRebuild(context.Context, *volume_server_pb.VolumeEcShardsRebuildRequest) (*volume_server_pb.VolumeEcShardsRebuildResponse, error) {
	return nil, errFakeVolumeUnsupported
}

func (unimplementedVolumeServer) VolumeEcShardsCopy(context.Context, *volume_server_pb.VolumeEcShardsCopyRequest) (*volume_server_pb.VolumeEcShardsCopyResponse, error) {
	return nil, errFakeVolumeUnsupported
}

func (unimplementedVolumeServer) VolumeEcShardRead(*volume_server_pb.VolumeEcShardReadRequest, volume_server_pb.VolumeServer_VolumeEcShardReadServer) error {
	return errFakeVolumeUnsupported
}

func (unimplementedVolumeServer) VolumeEcBlobDelete(context.Context, *volume_server_pb.VolumeEcBlobDeleteRequest) (*volume_server_pb.VolumeEcBlobDeleteResponse, error) {
	return nil, errFakeVolumeUnsupported
}

func (unimplementedVolumeServer) VolumeEcShardsToVolume(context.Context, *volume_server_pb.VolumeEcShardsToVolumeRequest) (*volume_server_pb.VolumeEcShardsToVolumeResponse, error) {
	return nil, errFakeVolumeUnsupported
}

func (unimplementedVolumeServer) VolumeTierMoveDatToRemote(*volume_server_pb.VolumeTierMoveDatToRemoteRequest, volume_server_pb.VolumeServer_VolumeTierMoveDatToRemoteServer) error {
	return errFakeVolumeUnsupported
}

func (unimplementedVolumeServer) VolumeTierMoveDatFromRemote(*volume_server_pb.VolumeTierMoveDatFromRemoteRequest, volume_server_pb.VolumeServer_VolumeTierMoveDatFromRemoteServer) error {
	return errFakeVolumeUnsupported
}

func (unimplementedVolumeServer) VolumeServerStatus(context.Context, *volume_server_pb.VolumeServerStatusRequest) (*volume_server_pb.VolumeServerStatusResponse, error) {
	return nil, errFakeVolumeUnsupported
}

func (unimplementedVolumeServer) VolumeServerLeave(context.Context, *volume_server_pb.VolumeServerLeaveRequest) (*volume_server_pb.VolumeServerLeaveResponse, error) {
	return nil, errFakeVolumeUnsupported
}

func (unimplementedVolumeServer) FetchAndWriteNeedle(context.Context, *volume_server_pb.FetchAndWriteNeedleRequest) (*volume_server_pb.FetchAndWriteNeedleResponse, error) {
	return nil, errFakeVolumeUnsupported
}

func (unimplementedVolumeServer) ScrubVolume(context.Context, *volume_server_pb.ScrubVolumeRequest) (*volume_server_pb.ScrubVolumeResponse, error) {
	return nil, errFakeVolumeUnsupported
}

func (unimplementedVolumeServer) ScrubEcVolume(context.Context, *volume_server_pb.ScrubEcVolumeRequest) (*volume_server_pb.ScrubEcVolumeResponse, error) {
	return nil, errFakeVolumeUnsupported
}

func (unimplementedVolumeServer) Query(*volume_server_pb.QueryRequest, volume_server_pb.VolumeServer_QueryServer) error {
	return errFakeVolumeUnsupported
}

func (unimplementedVolumeServer) VolumeNeedleStatus(context.Context, *volume_server_pb.VolumeNeedleStatusRequest) (*volume_server_pb.VolumeNeedleStatusResponse, error) {
	return nil, errFakeVolumeUnsupported
}

func (unimplementedVolumeServer) Ping(context.Context, *volume_server_pb.PingRequest) (*volume_server_pb.PingResponse, error) {
	return nil, errFakeVolumeUnsupported
}
