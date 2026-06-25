// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// rpc.go holds the per-RPC request builders and response decoders for the 37
// VolumeServer RPCs, built on the leaf converters in convert.go. Naming:
// <Rpc>ReqToWire(*pb.<Rpc>Request) []byte and, for unary RPCs,
// <Rpc>RespFromWire([]byte) (*pb.<Rpc>Response, error). Streaming response items
// arrive as typed wire views from the wire client's stream Recv, so each
// streaming RPC also gets <Rpc>RespFromView(wire view) *pb.<Rpc>Response; the
// byte decoder wraps then delegates to it, so the field mapping lives once.

package volumezap

import (
	volume_server_pb "github.com/hanzoai/s3/s3/pb/volume_server_pb"
	vsw "github.com/hanzoai/s3/s3/wire/volume_server"
)

// BatchDelete
func BatchDeleteReqToWire(r *volume_server_pb.BatchDeleteRequest) []byte {
	return vsw.NewBatchDeleteRequest(vsw.BatchDeleteRequestInput{
		FileIDs: r.FileIds, SkipCookieCheck: r.SkipCookieCheck,
	})
}
func BatchDeleteRespFromWire(b []byte) (*volume_server_pb.BatchDeleteResponse, error) {
	v, err := vsw.WrapBatchDeleteResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &volume_server_pb.BatchDeleteResponse{}
	for i := 0; i < v.ResultsLen(); i++ {
		if res, ok := v.ResultAt(i); ok {
			resp.Results = append(resp.Results, deleteResultFromView(res))
		}
	}
	return resp, nil
}

// VacuumVolumeCheck
func VacuumVolumeCheckReqToWire(r *volume_server_pb.VacuumVolumeCheckRequest) []byte {
	return vsw.NewVacuumVolumeCheckRequest(vsw.VacuumVolumeCheckRequestInput{VolumeID: r.VolumeId})
}
func VacuumVolumeCheckRespFromWire(b []byte) (*volume_server_pb.VacuumVolumeCheckResponse, error) {
	v, err := vsw.WrapVacuumVolumeCheckResponse(b)
	if err != nil {
		return nil, err
	}
	return &volume_server_pb.VacuumVolumeCheckResponse{GarbageRatio: v.GarbageRatio()}, nil
}

// VacuumVolumeCommit
func VacuumVolumeCommitReqToWire(r *volume_server_pb.VacuumVolumeCommitRequest) []byte {
	return vsw.NewVacuumVolumeCommitRequest(vsw.VacuumVolumeCommitRequestInput{VolumeID: r.VolumeId})
}
func VacuumVolumeCommitRespFromWire(b []byte) (*volume_server_pb.VacuumVolumeCommitResponse, error) {
	v, err := vsw.WrapVacuumVolumeCommitResponse(b)
	if err != nil {
		return nil, err
	}
	return &volume_server_pb.VacuumVolumeCommitResponse{IsReadOnly: v.IsReadOnly(), VolumeSize: v.VolumeSize()}, nil
}

// VacuumVolumeCleanup
func VacuumVolumeCleanupReqToWire(r *volume_server_pb.VacuumVolumeCleanupRequest) []byte {
	return vsw.NewVacuumVolumeCleanupRequest(vsw.VacuumVolumeCleanupRequestInput{VolumeID: r.VolumeId})
}
func VacuumVolumeCleanupRespFromWire(b []byte) (*volume_server_pb.VacuumVolumeCleanupResponse, error) {
	if _, err := vsw.WrapVacuumVolumeCleanupResponse(b); err != nil {
		return nil, err
	}
	return &volume_server_pb.VacuumVolumeCleanupResponse{}, nil
}

// DeleteCollection
func DeleteCollectionReqToWire(r *volume_server_pb.DeleteCollectionRequest) []byte {
	return vsw.NewDeleteCollectionRequest(vsw.DeleteCollectionRequestInput{Collection: r.Collection})
}
func DeleteCollectionRespFromWire(b []byte) (*volume_server_pb.DeleteCollectionResponse, error) {
	if _, err := vsw.WrapDeleteCollectionResponse(b); err != nil {
		return nil, err
	}
	return &volume_server_pb.DeleteCollectionResponse{}, nil
}

// AllocateVolume
func AllocateVolumeReqToWire(r *volume_server_pb.AllocateVolumeRequest) []byte {
	return vsw.NewAllocateVolumeRequest(vsw.AllocateVolumeRequestInput{
		VolumeID: r.VolumeId, Collection: r.Collection, Preallocate: r.Preallocate,
		Replication: r.Replication, TTL: r.Ttl, MemoryMapMaxSizeMb: r.MemoryMapMaxSizeMb,
		DiskType: r.DiskType, Version: r.Version,
	})
}
func AllocateVolumeRespFromWire(b []byte) (*volume_server_pb.AllocateVolumeResponse, error) {
	if _, err := vsw.WrapAllocateVolumeResponse(b); err != nil {
		return nil, err
	}
	return &volume_server_pb.AllocateVolumeResponse{}, nil
}

// VolumeSyncStatus
func VolumeSyncStatusReqToWire(r *volume_server_pb.VolumeSyncStatusRequest) []byte {
	return vsw.NewVolumeSyncStatusRequest(vsw.VolumeSyncStatusRequestInput{VolumeID: r.VolumeId})
}
func VolumeSyncStatusRespFromWire(b []byte) (*volume_server_pb.VolumeSyncStatusResponse, error) {
	v, err := vsw.WrapVolumeSyncStatusResponse(b)
	if err != nil {
		return nil, err
	}
	return &volume_server_pb.VolumeSyncStatusResponse{
		VolumeId: v.VolumeID(), Collection: v.Collection(), Replication: v.Replication(),
		Ttl: v.TTL(), TailOffset: v.TailOffset(), CompactRevision: v.CompactRevision(),
		IdxFileSize: v.IdxFileSize(), Version: v.Version(),
	}, nil
}

// VolumeMount
func VolumeMountReqToWire(r *volume_server_pb.VolumeMountRequest) []byte {
	return vsw.NewVolumeMountRequest(vsw.VolumeMountRequestInput{VolumeID: r.VolumeId})
}
func VolumeMountRespFromWire(b []byte) (*volume_server_pb.VolumeMountResponse, error) {
	if _, err := vsw.WrapVolumeMountResponse(b); err != nil {
		return nil, err
	}
	return &volume_server_pb.VolumeMountResponse{}, nil
}

// VolumeUnmount
func VolumeUnmountReqToWire(r *volume_server_pb.VolumeUnmountRequest) []byte {
	return vsw.NewVolumeUnmountRequest(vsw.VolumeUnmountRequestInput{VolumeID: r.VolumeId})
}
func VolumeUnmountRespFromWire(b []byte) (*volume_server_pb.VolumeUnmountResponse, error) {
	if _, err := vsw.WrapVolumeUnmountResponse(b); err != nil {
		return nil, err
	}
	return &volume_server_pb.VolumeUnmountResponse{}, nil
}

// VolumeDelete
func VolumeDeleteReqToWire(r *volume_server_pb.VolumeDeleteRequest) []byte {
	return vsw.NewVolumeDeleteRequest(vsw.VolumeDeleteRequestInput{
		VolumeID: r.VolumeId, OnlyEmpty: r.OnlyEmpty, KeepRemoteData: r.KeepRemoteData,
	})
}
func VolumeDeleteRespFromWire(b []byte) (*volume_server_pb.VolumeDeleteResponse, error) {
	if _, err := vsw.WrapVolumeDeleteResponse(b); err != nil {
		return nil, err
	}
	return &volume_server_pb.VolumeDeleteResponse{}, nil
}

// VolumeMarkReadonly
func VolumeMarkReadonlyReqToWire(r *volume_server_pb.VolumeMarkReadonlyRequest) []byte {
	return vsw.NewVolumeMarkReadonlyRequest(vsw.VolumeMarkReadonlyRequestInput{VolumeID: r.VolumeId, Persist: r.Persist})
}
func VolumeMarkReadonlyRespFromWire(b []byte) (*volume_server_pb.VolumeMarkReadonlyResponse, error) {
	if _, err := vsw.WrapVolumeMarkReadonlyResponse(b); err != nil {
		return nil, err
	}
	return &volume_server_pb.VolumeMarkReadonlyResponse{}, nil
}

// VolumeMarkWritable
func VolumeMarkWritableReqToWire(r *volume_server_pb.VolumeMarkWritableRequest) []byte {
	return vsw.NewVolumeMarkWritableRequest(vsw.VolumeMarkWritableRequestInput{VolumeID: r.VolumeId})
}
func VolumeMarkWritableRespFromWire(b []byte) (*volume_server_pb.VolumeMarkWritableResponse, error) {
	if _, err := vsw.WrapVolumeMarkWritableResponse(b); err != nil {
		return nil, err
	}
	return &volume_server_pb.VolumeMarkWritableResponse{}, nil
}

// VolumeConfigure
func VolumeConfigureReqToWire(r *volume_server_pb.VolumeConfigureRequest) []byte {
	return vsw.NewVolumeConfigureRequest(vsw.VolumeConfigureRequestInput{VolumeID: r.VolumeId, Replication: r.Replication})
}
func VolumeConfigureRespFromWire(b []byte) (*volume_server_pb.VolumeConfigureResponse, error) {
	v, err := vsw.WrapVolumeConfigureResponse(b)
	if err != nil {
		return nil, err
	}
	return &volume_server_pb.VolumeConfigureResponse{Error: v.Error()}, nil
}

// VolumeStatus
func VolumeStatusReqToWire(r *volume_server_pb.VolumeStatusRequest) []byte {
	return vsw.NewVolumeStatusRequest(vsw.VolumeStatusRequestInput{VolumeID: r.VolumeId})
}
func VolumeStatusRespFromWire(b []byte) (*volume_server_pb.VolumeStatusResponse, error) {
	v, err := vsw.WrapVolumeStatusResponse(b)
	if err != nil {
		return nil, err
	}
	return &volume_server_pb.VolumeStatusResponse{
		IsReadOnly: v.IsReadOnly(), VolumeSize: v.VolumeSize(),
		FileCount: v.FileCount(), FileDeletedCount: v.FileDeletedCount(),
	}, nil
}

// GetState
func GetStateReqToWire(r *volume_server_pb.GetStateRequest) []byte {
	return vsw.NewGetStateRequest(vsw.GetStateRequestInput{})
}
func GetStateRespFromWire(b []byte) (*volume_server_pb.GetStateResponse, error) {
	v, err := vsw.WrapGetStateResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &volume_server_pb.GetStateResponse{}
	if s, ok := v.State(); ok {
		resp.State = volumeServerStateFromView(s)
	}
	return resp, nil
}

// SetState
func SetStateReqToWire(r *volume_server_pb.SetStateRequest) []byte {
	return vsw.NewSetStateRequest(vsw.SetStateRequestInput{State: volumeServerStateToWire(r.State)})
}
func SetStateRespFromWire(b []byte) (*volume_server_pb.SetStateResponse, error) {
	v, err := vsw.WrapSetStateResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &volume_server_pb.SetStateResponse{}
	if s, ok := v.State(); ok {
		resp.State = volumeServerStateFromView(s)
	}
	return resp, nil
}

// ReadVolumeFileStatus
func ReadVolumeFileStatusReqToWire(r *volume_server_pb.ReadVolumeFileStatusRequest) []byte {
	return vsw.NewReadVolumeFileStatusRequest(vsw.ReadVolumeFileStatusRequestInput{VolumeID: r.VolumeId})
}
func ReadVolumeFileStatusRespFromWire(b []byte) (*volume_server_pb.ReadVolumeFileStatusResponse, error) {
	v, err := vsw.WrapReadVolumeFileStatusResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &volume_server_pb.ReadVolumeFileStatusResponse{
		VolumeId:                v.VolumeID(),
		IdxFileTimestampSeconds: v.IdxFileTimestampSeconds(),
		IdxFileSize:             v.IdxFileSize(),
		DatFileTimestampSeconds: v.DatFileTimestampSeconds(),
		DatFileSize:             v.DatFileSize(),
		FileCount:               v.FileCount(),
		CompactionRevision:      v.CompactionRevision(),
		Collection:              v.Collection(),
		DiskType:                v.DiskType(),
		Version:                 v.Version(),
	}
	if info, ok := v.VolumeInfo(); ok {
		resp.VolumeInfo = volumeInfoFromView(info)
	}
	return resp, nil
}

// ReadNeedleBlob
func ReadNeedleBlobReqToWire(r *volume_server_pb.ReadNeedleBlobRequest) []byte {
	return vsw.NewReadNeedleBlobRequest(vsw.ReadNeedleBlobRequestInput{VolumeID: r.VolumeId, Offset: r.Offset, Size: r.Size})
}
func ReadNeedleBlobRespFromWire(b []byte) (*volume_server_pb.ReadNeedleBlobResponse, error) {
	v, err := vsw.WrapReadNeedleBlobResponse(b)
	if err != nil {
		return nil, err
	}
	return &volume_server_pb.ReadNeedleBlobResponse{NeedleBlob: v.NeedleBlob()}, nil
}

// ReadNeedleMeta
func ReadNeedleMetaReqToWire(r *volume_server_pb.ReadNeedleMetaRequest) []byte {
	return vsw.NewReadNeedleMetaRequest(vsw.ReadNeedleMetaRequestInput{
		VolumeID: r.VolumeId, NeedleID: r.NeedleId, Offset: r.Offset, Size: r.Size,
	})
}
func ReadNeedleMetaRespFromWire(b []byte) (*volume_server_pb.ReadNeedleMetaResponse, error) {
	v, err := vsw.WrapReadNeedleMetaResponse(b)
	if err != nil {
		return nil, err
	}
	return &volume_server_pb.ReadNeedleMetaResponse{
		Cookie: v.Cookie(), LastModified: v.LastModified(), Crc: v.Crc(), Ttl: v.TTL(), AppendAtNs: v.AppendAtNs(),
	}, nil
}

// WriteNeedleBlob
func WriteNeedleBlobReqToWire(r *volume_server_pb.WriteNeedleBlobRequest) []byte {
	return vsw.NewWriteNeedleBlobRequest(vsw.WriteNeedleBlobRequestInput{
		VolumeID: r.VolumeId, NeedleID: r.NeedleId, Size: r.Size, NeedleBlob: r.NeedleBlob,
	})
}
func WriteNeedleBlobRespFromWire(b []byte) (*volume_server_pb.WriteNeedleBlobResponse, error) {
	if _, err := vsw.WrapWriteNeedleBlobResponse(b); err != nil {
		return nil, err
	}
	return &volume_server_pb.WriteNeedleBlobResponse{}, nil
}

// VolumeTailReceiver
func VolumeTailReceiverReqToWire(r *volume_server_pb.VolumeTailReceiverRequest) []byte {
	return vsw.NewVolumeTailReceiverRequest(vsw.VolumeTailReceiverRequestInput{
		VolumeID: r.VolumeId, SinceNs: r.SinceNs, IdleTimeoutSeconds: r.IdleTimeoutSeconds, SourceVolumeServer: r.SourceVolumeServer,
	})
}
func VolumeTailReceiverRespFromWire(b []byte) (*volume_server_pb.VolumeTailReceiverResponse, error) {
	if _, err := vsw.WrapVolumeTailReceiverResponse(b); err != nil {
		return nil, err
	}
	return &volume_server_pb.VolumeTailReceiverResponse{}, nil
}

// VolumeEcShardsGenerate
func VolumeEcShardsGenerateReqToWire(r *volume_server_pb.VolumeEcShardsGenerateRequest) []byte {
	return vsw.NewVolumeEcShardsGenerateRequest(vsw.VolumeEcShardsGenerateRequestInput{VolumeID: r.VolumeId, Collection: r.Collection})
}
func VolumeEcShardsGenerateRespFromWire(b []byte) (*volume_server_pb.VolumeEcShardsGenerateResponse, error) {
	if _, err := vsw.WrapVolumeEcShardsGenerateResponse(b); err != nil {
		return nil, err
	}
	return &volume_server_pb.VolumeEcShardsGenerateResponse{}, nil
}

// VolumeEcShardsRebuild
func VolumeEcShardsRebuildReqToWire(r *volume_server_pb.VolumeEcShardsRebuildRequest) []byte {
	return vsw.NewVolumeEcShardsRebuildRequest(vsw.VolumeEcShardsRebuildRequestInput{
		VolumeID: r.VolumeId, Collection: r.Collection, UnsafeIgnoreSidecar: r.UnsafeIgnoreSidecar,
	})
}
func VolumeEcShardsRebuildRespFromWire(b []byte) (*volume_server_pb.VolumeEcShardsRebuildResponse, error) {
	v, err := vsw.WrapVolumeEcShardsRebuildResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &volume_server_pb.VolumeEcShardsRebuildResponse{}
	for i := 0; i < v.RebuiltShardIDsLen(); i++ {
		resp.RebuiltShardIds = append(resp.RebuiltShardIds, v.RebuiltShardIDAt(i))
	}
	return resp, nil
}

// VolumeEcShardsCopy
func VolumeEcShardsCopyReqToWire(r *volume_server_pb.VolumeEcShardsCopyRequest) []byte {
	return vsw.NewVolumeEcShardsCopyRequest(vsw.VolumeEcShardsCopyRequestInput{
		VolumeID: r.VolumeId, Collection: r.Collection, ShardIDs: r.ShardIds, CopyEcxFile: r.CopyEcxFile,
		SourceDataNode: r.SourceDataNode, CopyEcjFile: r.CopyEcjFile, CopyVifFile: r.CopyVifFile,
		DiskID: r.DiskId, CopyEcsumFile: r.CopyEcsumFile,
	})
}
func VolumeEcShardsCopyRespFromWire(b []byte) (*volume_server_pb.VolumeEcShardsCopyResponse, error) {
	if _, err := vsw.WrapVolumeEcShardsCopyResponse(b); err != nil {
		return nil, err
	}
	return &volume_server_pb.VolumeEcShardsCopyResponse{}, nil
}

// VolumeEcShardsDelete
func VolumeEcShardsDeleteReqToWire(r *volume_server_pb.VolumeEcShardsDeleteRequest) []byte {
	return vsw.NewVolumeEcShardsDeleteRequest(vsw.VolumeEcShardsDeleteRequestInput{
		VolumeID: r.VolumeId, Collection: r.Collection, ShardIDs: r.ShardIds,
		FullTeardown: r.FullTeardown, EncodeTsNs: r.EncodeTsNs,
	})
}
func VolumeEcShardsDeleteRespFromWire(b []byte) (*volume_server_pb.VolumeEcShardsDeleteResponse, error) {
	v, err := vsw.WrapVolumeEcShardsDeleteResponse(b)
	if err != nil {
		return nil, err
	}
	return &volume_server_pb.VolumeEcShardsDeleteResponse{FullTeardownDone: v.FullTeardownDone()}, nil
}

// VolumeEcShardsMount
func VolumeEcShardsMountReqToWire(r *volume_server_pb.VolumeEcShardsMountRequest) []byte {
	return vsw.NewVolumeEcShardsMountRequest(vsw.VolumeEcShardsMountRequestInput{
		VolumeID: r.VolumeId, Collection: r.Collection, ShardIDs: r.ShardIds, SourceDiskType: r.SourceDiskType,
	})
}
func VolumeEcShardsMountRespFromWire(b []byte) (*volume_server_pb.VolumeEcShardsMountResponse, error) {
	if _, err := vsw.WrapVolumeEcShardsMountResponse(b); err != nil {
		return nil, err
	}
	return &volume_server_pb.VolumeEcShardsMountResponse{}, nil
}

// VolumeEcShardsUnmount
func VolumeEcShardsUnmountReqToWire(r *volume_server_pb.VolumeEcShardsUnmountRequest) []byte {
	return vsw.NewVolumeEcShardsUnmountRequest(vsw.VolumeEcShardsUnmountRequestInput{
		VolumeID: r.VolumeId, ShardIDs: r.ShardIds, EncodeTsNs: r.EncodeTsNs,
	})
}
func VolumeEcShardsUnmountRespFromWire(b []byte) (*volume_server_pb.VolumeEcShardsUnmountResponse, error) {
	if _, err := vsw.WrapVolumeEcShardsUnmountResponse(b); err != nil {
		return nil, err
	}
	return &volume_server_pb.VolumeEcShardsUnmountResponse{}, nil
}

// VolumeEcBlobDelete
func VolumeEcBlobDeleteReqToWire(r *volume_server_pb.VolumeEcBlobDeleteRequest) []byte {
	return vsw.NewVolumeEcBlobDeleteRequest(vsw.VolumeEcBlobDeleteRequestInput{
		VolumeID: r.VolumeId, Collection: r.Collection, FileKey: r.FileKey, Version: r.Version,
	})
}
func VolumeEcBlobDeleteRespFromWire(b []byte) (*volume_server_pb.VolumeEcBlobDeleteResponse, error) {
	if _, err := vsw.WrapVolumeEcBlobDeleteResponse(b); err != nil {
		return nil, err
	}
	return &volume_server_pb.VolumeEcBlobDeleteResponse{}, nil
}

// VolumeEcShardsToVolume
func VolumeEcShardsToVolumeReqToWire(r *volume_server_pb.VolumeEcShardsToVolumeRequest) []byte {
	return vsw.NewVolumeEcShardsToVolumeRequest(vsw.VolumeEcShardsToVolumeRequestInput{VolumeID: r.VolumeId, Collection: r.Collection})
}
func VolumeEcShardsToVolumeRespFromWire(b []byte) (*volume_server_pb.VolumeEcShardsToVolumeResponse, error) {
	if _, err := vsw.WrapVolumeEcShardsToVolumeResponse(b); err != nil {
		return nil, err
	}
	return &volume_server_pb.VolumeEcShardsToVolumeResponse{}, nil
}

// VolumeEcShardsInfo
func VolumeEcShardsInfoReqToWire(r *volume_server_pb.VolumeEcShardsInfoRequest) []byte {
	return vsw.NewVolumeEcShardsInfoRequest(vsw.VolumeEcShardsInfoRequestInput{VolumeID: r.VolumeId})
}
func VolumeEcShardsInfoRespFromWire(b []byte) (*volume_server_pb.VolumeEcShardsInfoResponse, error) {
	v, err := vsw.WrapVolumeEcShardsInfoResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &volume_server_pb.VolumeEcShardsInfoResponse{
		VolumeSize: v.VolumeSize(), FileCount: v.FileCount(), FileDeletedCount: v.FileDeletedCount(),
	}
	for i := 0; i < v.EcShardInfosLen(); i++ {
		if info, ok := v.EcShardInfoAt(i); ok {
			resp.EcShardInfos = append(resp.EcShardInfos, ecShardInfoFromView(info))
		}
	}
	return resp, nil
}

// VolumeServerStatus
func VolumeServerStatusReqToWire(r *volume_server_pb.VolumeServerStatusRequest) []byte {
	return vsw.NewVolumeServerStatusRequest(vsw.VolumeServerStatusRequestInput{})
}
func VolumeServerStatusRespFromWire(b []byte) (*volume_server_pb.VolumeServerStatusResponse, error) {
	v, err := vsw.WrapVolumeServerStatusResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &volume_server_pb.VolumeServerStatusResponse{
		Version: v.Version(), DataCenter: v.DataCenter(), Rack: v.Rack(),
	}
	for i := 0; i < v.DiskStatusesLen(); i++ {
		if d, ok := v.DiskStatusAt(i); ok {
			resp.DiskStatuses = append(resp.DiskStatuses, diskStatusFromView(d))
		}
	}
	if mem, ok := v.MemoryStatus(); ok {
		resp.MemoryStatus = memStatusFromView(mem)
	}
	if s, ok := v.State(); ok {
		resp.State = volumeServerStateFromView(s)
	}
	return resp, nil
}

// VolumeServerLeave
func VolumeServerLeaveReqToWire(r *volume_server_pb.VolumeServerLeaveRequest) []byte {
	return vsw.NewVolumeServerLeaveRequest(vsw.VolumeServerLeaveRequestInput{})
}
func VolumeServerLeaveRespFromWire(b []byte) (*volume_server_pb.VolumeServerLeaveResponse, error) {
	if _, err := vsw.WrapVolumeServerLeaveResponse(b); err != nil {
		return nil, err
	}
	return &volume_server_pb.VolumeServerLeaveResponse{}, nil
}

// FetchAndWriteNeedle
func FetchAndWriteNeedleReqToWire(r *volume_server_pb.FetchAndWriteNeedleRequest) []byte {
	replicas := make([][]byte, len(r.Replicas))
	for i, rep := range r.Replicas {
		replicas[i] = fetchReplicaToWire(rep)
	}
	return vsw.NewFetchAndWriteNeedleRequest(vsw.FetchAndWriteNeedleRequestInput{
		VolumeID: r.VolumeId, NeedleID: r.NeedleId, Cookie: r.Cookie, Offset: r.Offset, Size: r.Size,
		Replicas: replicas, Auth: r.Auth, DownloadConcurrency: r.DownloadConcurrency,
		RemoteConf: remoteConfToWire(r.RemoteConf), RemoteLocation: remoteLocationToWire(r.RemoteLocation),
	})
}
func FetchAndWriteNeedleRespFromWire(b []byte) (*volume_server_pb.FetchAndWriteNeedleResponse, error) {
	v, err := vsw.WrapFetchAndWriteNeedleResponse(b)
	if err != nil {
		return nil, err
	}
	return &volume_server_pb.FetchAndWriteNeedleResponse{ETag: v.ETag()}, nil
}

// ScrubVolume
func ScrubVolumeReqToWire(r *volume_server_pb.ScrubVolumeRequest) []byte {
	return vsw.NewScrubVolumeRequest(vsw.ScrubVolumeRequestInput{
		Mode: uint32(r.Mode), VolumeIDs: r.VolumeIds, MarkBrokenVolumesReadonly: r.MarkBrokenVolumesReadonly,
	})
}
func ScrubVolumeRespFromWire(b []byte) (*volume_server_pb.ScrubVolumeResponse, error) {
	v, err := vsw.WrapScrubVolumeResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &volume_server_pb.ScrubVolumeResponse{TotalVolumes: v.TotalVolumes(), TotalFiles: v.TotalFiles()}
	for i := 0; i < v.BrokenVolumeIDsLen(); i++ {
		resp.BrokenVolumeIds = append(resp.BrokenVolumeIds, v.BrokenVolumeIDAt(i))
	}
	for i := 0; i < v.DetailsLen(); i++ {
		resp.Details = append(resp.Details, v.DetailAt(i))
	}
	return resp, nil
}

// ScrubEcVolume
func ScrubEcVolumeReqToWire(r *volume_server_pb.ScrubEcVolumeRequest) []byte {
	return vsw.NewScrubEcVolumeRequest(vsw.ScrubEcVolumeRequestInput{Mode: uint32(r.Mode), VolumeIDs: r.VolumeIds})
}
func ScrubEcVolumeRespFromWire(b []byte) (*volume_server_pb.ScrubEcVolumeResponse, error) {
	v, err := vsw.WrapScrubEcVolumeResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &volume_server_pb.ScrubEcVolumeResponse{TotalVolumes: v.TotalVolumes(), TotalFiles: v.TotalFiles()}
	for i := 0; i < v.BrokenVolumeIDsLen(); i++ {
		resp.BrokenVolumeIds = append(resp.BrokenVolumeIds, v.BrokenVolumeIDAt(i))
	}
	for i := 0; i < v.BrokenShardInfosLen(); i++ {
		if info, ok := v.BrokenShardInfoAt(i); ok {
			resp.BrokenShardInfos = append(resp.BrokenShardInfos, ecShardInfoFromView(info))
		}
	}
	for i := 0; i < v.DetailsLen(); i++ {
		resp.Details = append(resp.Details, v.DetailAt(i))
	}
	return resp, nil
}

// VolumeNeedleStatus
func VolumeNeedleStatusReqToWire(r *volume_server_pb.VolumeNeedleStatusRequest) []byte {
	return vsw.NewVolumeNeedleStatusRequest(vsw.VolumeNeedleStatusRequestInput{VolumeID: r.VolumeId, NeedleID: r.NeedleId})
}
func VolumeNeedleStatusRespFromWire(b []byte) (*volume_server_pb.VolumeNeedleStatusResponse, error) {
	v, err := vsw.WrapVolumeNeedleStatusResponse(b)
	if err != nil {
		return nil, err
	}
	return &volume_server_pb.VolumeNeedleStatusResponse{
		NeedleId: v.NeedleID(), Cookie: v.Cookie(), Size: v.Size(),
		LastModified: v.LastModified(), Crc: v.Crc(), Ttl: v.TTL(),
	}, nil
}

// Ping
func PingReqToWire(r *volume_server_pb.PingRequest) []byte {
	return vsw.NewPingRequest(vsw.PingRequestInput{Target: r.Target, TargetType: r.TargetType})
}
func PingRespFromWire(b []byte) (*volume_server_pb.PingResponse, error) {
	v, err := vsw.WrapPingResponse(b)
	if err != nil {
		return nil, err
	}
	return &volume_server_pb.PingResponse{StartTimeNs: v.StartTimeNs(), RemoteTimeNs: v.RemoteTimeNs(), StopTimeNs: v.StopTimeNs()}, nil
}

// --- streaming request builders + response-view decoders ---

// VacuumVolumeCompact
func VacuumVolumeCompactReqInput(r *volume_server_pb.VacuumVolumeCompactRequest) vsw.VacuumVolumeCompactRequestInput {
	return vsw.VacuumVolumeCompactRequestInput{VolumeID: r.VolumeId, Preallocate: r.Preallocate}
}
func VacuumVolumeCompactRespFromView(v vsw.VacuumVolumeCompactResponse) *volume_server_pb.VacuumVolumeCompactResponse {
	return &volume_server_pb.VacuumVolumeCompactResponse{ProcessedBytes: v.ProcessedBytes(), LoadAvg_1M: v.LoadAvg1m()}
}

// VolumeIncrementalCopy
func VolumeIncrementalCopyReqInput(r *volume_server_pb.VolumeIncrementalCopyRequest) vsw.VolumeIncrementalCopyRequestInput {
	return vsw.VolumeIncrementalCopyRequestInput{VolumeID: r.VolumeId, SinceNs: r.SinceNs}
}
func VolumeIncrementalCopyRespFromView(v vsw.VolumeIncrementalCopyResponse) *volume_server_pb.VolumeIncrementalCopyResponse {
	return &volume_server_pb.VolumeIncrementalCopyResponse{FileContent: v.FileContent()}
}

// VolumeCopy
func VolumeCopyReqInput(r *volume_server_pb.VolumeCopyRequest) vsw.VolumeCopyRequestInput {
	return vsw.VolumeCopyRequestInput{
		VolumeID: r.VolumeId, Collection: r.Collection, Replication: r.Replication, TTL: r.Ttl,
		SourceDataNode: r.SourceDataNode, DiskType: r.DiskType, IoBytePerSecond: r.IoBytePerSecond,
	}
}
func VolumeCopyRespFromView(v vsw.VolumeCopyResponse) *volume_server_pb.VolumeCopyResponse {
	return &volume_server_pb.VolumeCopyResponse{LastAppendAtNs: v.LastAppendAtNs(), ProcessedBytes: v.ProcessedBytes()}
}

// CopyFile
func CopyFileReqInput(r *volume_server_pb.CopyFileRequest) vsw.CopyFileRequestInput {
	return vsw.CopyFileRequestInput{
		VolumeID: r.VolumeId, Ext: r.Ext, CompactionRevision: r.CompactionRevision, StopOffset: r.StopOffset,
		Collection: r.Collection, IsEcVolume: r.IsEcVolume, IgnoreSourceFileNotFound: r.IgnoreSourceFileNotFound,
	}
}
func CopyFileRespFromView(v vsw.CopyFileResponse) *volume_server_pb.CopyFileResponse {
	return &volume_server_pb.CopyFileResponse{FileContent: v.FileContent(), ModifiedTsNs: v.ModifiedTsNs()}
}

// ReadAllNeedles
func ReadAllNeedlesReqInput(r *volume_server_pb.ReadAllNeedlesRequest) vsw.ReadAllNeedlesRequestInput {
	return vsw.ReadAllNeedlesRequestInput{VolumeIDs: r.VolumeIds}
}
func ReadAllNeedlesRespFromView(v vsw.ReadAllNeedlesResponse) *volume_server_pb.ReadAllNeedlesResponse {
	return &volume_server_pb.ReadAllNeedlesResponse{
		VolumeId: v.VolumeID(), NeedleId: v.NeedleID(), Cookie: v.Cookie(), NeedleBlob: v.NeedleBlob(),
		NeedleBlobCompressed: v.NeedleBlobCompressed(), LastModified: v.LastModified(), Crc: v.Crc(),
		Name: v.Name(), Mime: v.Mime(),
	}
}

// VolumeTailSender
func VolumeTailSenderReqInput(r *volume_server_pb.VolumeTailSenderRequest) vsw.VolumeTailSenderRequestInput {
	return vsw.VolumeTailSenderRequestInput{VolumeID: r.VolumeId, SinceNs: r.SinceNs, IdleTimeoutSeconds: r.IdleTimeoutSeconds}
}
func VolumeTailSenderRespFromView(v vsw.VolumeTailSenderResponse) *volume_server_pb.VolumeTailSenderResponse {
	return &volume_server_pb.VolumeTailSenderResponse{
		NeedleHeader: v.NeedleHeader(), NeedleBody: v.NeedleBody(), IsLastChunk: v.IsLastChunk(), Version: v.Version(),
	}
}

// VolumeEcShardRead
func VolumeEcShardReadReqInput(r *volume_server_pb.VolumeEcShardReadRequest) vsw.VolumeEcShardReadRequestInput {
	return vsw.VolumeEcShardReadRequestInput{
		VolumeID: r.VolumeId, ShardID: r.ShardId, Offset: r.Offset, Size: r.Size, FileKey: r.FileKey, EncodeTsNs: r.EncodeTsNs,
	}
}
func VolumeEcShardReadRespFromView(v vsw.VolumeEcShardReadResponse) *volume_server_pb.VolumeEcShardReadResponse {
	return &volume_server_pb.VolumeEcShardReadResponse{Data: v.Data(), IsDeleted: v.IsDeleted(), EncodeTsNs: v.EncodeTsNs()}
}

// VolumeTierMoveDatToRemote
func VolumeTierMoveDatToRemoteReqInput(r *volume_server_pb.VolumeTierMoveDatToRemoteRequest) vsw.VolumeTierMoveDatToRemoteRequestInput {
	return vsw.VolumeTierMoveDatToRemoteRequestInput{
		VolumeID: r.VolumeId, Collection: r.Collection, DestinationBackendName: r.DestinationBackendName, KeepLocalDatFile: r.KeepLocalDatFile,
	}
}
func VolumeTierMoveDatToRemoteRespFromView(v vsw.VolumeTierMoveDatToRemoteResponse) *volume_server_pb.VolumeTierMoveDatToRemoteResponse {
	return &volume_server_pb.VolumeTierMoveDatToRemoteResponse{Processed: v.Processed(), ProcessedPercentage: v.ProcessedPercentage()}
}

// VolumeTierMoveDatFromRemote
func VolumeTierMoveDatFromRemoteReqInput(r *volume_server_pb.VolumeTierMoveDatFromRemoteRequest) vsw.VolumeTierMoveDatFromRemoteRequestInput {
	return vsw.VolumeTierMoveDatFromRemoteRequestInput{
		VolumeID: r.VolumeId, Collection: r.Collection, KeepRemoteDatFile: r.KeepRemoteDatFile,
	}
}
func VolumeTierMoveDatFromRemoteRespFromView(v vsw.VolumeTierMoveDatFromRemoteResponse) *volume_server_pb.VolumeTierMoveDatFromRemoteResponse {
	return &volume_server_pb.VolumeTierMoveDatFromRemoteResponse{Processed: v.Processed(), ProcessedPercentage: v.ProcessedPercentage()}
}

// Query
func QueryReqInput(r *volume_server_pb.QueryRequest) vsw.QueryRequestInput {
	return vsw.QueryRequestInput{
		Selections:          r.Selections,
		FromFileIDs:         r.FromFileIds,
		Filter:              queryFilterToWire(r.Filter),
		InputSerialization:  queryInputSerializationToWire(r.InputSerialization),
		OutputSerialization: queryOutputSerializationToWire(r.OutputSerialization),
	}
}
func QueriedStripeFromView(v vsw.QueriedStripe) *volume_server_pb.QueriedStripe {
	return &volume_server_pb.QueriedStripe{Records: v.Records()}
}

// ReceiveFile (client-streaming)
func ReceiveFileInfoReqInput(r *volume_server_pb.ReceiveFileRequest) vsw.ReceiveFileRequestInput {
	return receiveFileReqInput(r)
}
func ReceiveFileRespFromView(v vsw.ReceiveFileResponse) *volume_server_pb.ReceiveFileResponse {
	return &volume_server_pb.ReceiveFileResponse{BytesWritten: v.BytesWritten(), Error: v.Error()}
}

// receiveFileReqInput maps the ReceiveFileRequest oneof (info vs file_content)
// onto the wire's Data discriminator + variant sub-buffers.
func receiveFileReqInput(r *volume_server_pb.ReceiveFileRequest) vsw.ReceiveFileRequestInput {
	switch d := r.Data.(type) {
	case *volume_server_pb.ReceiveFileRequest_Info:
		var info []byte
		if d.Info != nil {
			info = vsw.NewReceiveFileInfo(vsw.ReceiveFileInfoInput{
				VolumeID: d.Info.VolumeId, Ext: d.Info.Ext, Collection: d.Info.Collection,
				IsEcVolume: d.Info.IsEcVolume, ShardID: d.Info.ShardId, FileSize: d.Info.FileSize, DiskID: d.Info.DiskId,
			})
		}
		return vsw.ReceiveFileRequestInput{Data: vsw.ReceiveFileDataInfo, Info: info}
	case *volume_server_pb.ReceiveFileRequest_FileContent:
		return vsw.ReceiveFileRequestInput{Data: vsw.ReceiveFileDataFileContent, FileContent: d.FileContent}
	default:
		return vsw.ReceiveFileRequestInput{}
	}
}

// --- Query nested message converters ---

func queryFilterToWire(f *volume_server_pb.QueryRequest_Filter) []byte {
	if f == nil {
		return nil
	}
	return vsw.NewQueryFilter(vsw.QueryFilterInput{Field: f.Field, Operand: f.Operand, Value: f.Value})
}

func queryInputSerializationToWire(s *volume_server_pb.QueryRequest_InputSerialization) []byte {
	if s == nil {
		return nil
	}
	in := vsw.QueryInputSerializationInput{CompressionType: s.CompressionType}
	if c := s.CsvInput; c != nil {
		in.CSVInput = vsw.NewQueryInputSerializationCSVInput(vsw.QueryInputSerializationCSVInputInput{
			FileHeaderInfo: c.FileHeaderInfo, RecordDelimiter: c.RecordDelimiter, FieldDelimiter: c.FieldDelimiter,
			QuoteCharacter: c.QuoteCharacter, QuoteEscapeCharacter: c.QuoteEscapeCharacter, Comments: c.Comments,
			AllowQuotedRecordDelimiter: c.AllowQuotedRecordDelimiter,
		})
	}
	if j := s.JsonInput; j != nil {
		in.JSONInput = vsw.NewQueryInputSerializationJSONInput(vsw.QueryInputSerializationJSONInputInput{Type: j.Type})
	}
	if p := s.ParquetInput; p != nil {
		in.ParquetInput = vsw.NewQueryInputSerializationParquetInput(vsw.QueryInputSerializationParquetInputInput{})
	}
	return vsw.NewQueryInputSerialization(in)
}

func queryOutputSerializationToWire(s *volume_server_pb.QueryRequest_OutputSerialization) []byte {
	if s == nil {
		return nil
	}
	out := vsw.QueryOutputSerializationInput{}
	if c := s.CsvOutput; c != nil {
		out.CSVOutput = vsw.NewQueryOutputSerializationCSVOutput(vsw.QueryOutputSerializationCSVOutputInput{
			QuoteFields: c.QuoteFields, RecordDelimiter: c.RecordDelimiter, FieldDelimiter: c.FieldDelimiter,
			QuoteCharacter: c.QuoteCharacter, QuoteEscapeCharacter: c.QuoteEscapeCharacter,
		})
	}
	if j := s.JsonOutput; j != nil {
		out.JSONOutput = vsw.NewQueryOutputSerializationJSONOutput(vsw.QueryOutputSerializationJSONOutputInput{RecordDelimiter: j.RecordDelimiter})
	}
	return vsw.NewQueryOutputSerialization(out)
}
