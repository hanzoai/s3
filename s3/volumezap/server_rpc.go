// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// server_rpc.go holds the per-RPC server-side converters: the inverse of rpc.go.
// For each unary RPC:
//
//	<Rpc>ReqFromWire([]byte) (*pb.<Rpc>Request, error)   // Wrap<Req> + accessors -> pb
//	<Rpc>RespToWire(*pb.<Rpc>Response) []byte            // pb fields -> New<Resp>
//
// For each streaming RPC the request rides as a typed view (the stream's Init),
// so the request converter takes the view directly:
//
//	<Rpc>ReqFromView(vsw.<Rpc>Request) *pb.<Rpc>Request   // streaming request decode
//	<Rpc>RespToInput(*pb.<Rpc>Response) vsw.<Rpc>ResponseInput // streaming item encode
//
// Every field map is the exact inverse of the matching builder/decoder in
// rpc.go, so the wire shape stays identical in both directions.

package volumezap

import (
	volume_server_pb "github.com/hanzoai/s3/s3/pb/volume_server_pb"
	vsw "github.com/hanzoai/s3/s3/wire/volume_server"
)

// --- unary request decoders + response encoders (37) ---

// BatchDelete
func BatchDeleteReqFromWire(b []byte) (*volume_server_pb.BatchDeleteRequest, error) {
	v, err := vsw.WrapBatchDeleteRequest(b)
	if err != nil {
		return nil, err
	}
	req := &volume_server_pb.BatchDeleteRequest{SkipCookieCheck: v.SkipCookieCheck()}
	for i := 0; i < v.FileIDsLen(); i++ {
		req.FileIds = append(req.FileIds, v.FileIDAt(i))
	}
	return req, nil
}
func BatchDeleteRespToWire(r *volume_server_pb.BatchDeleteResponse) []byte {
	results := make([][]byte, len(r.Results))
	for i, res := range r.Results {
		results[i] = deleteResultToWire(res)
	}
	return vsw.NewBatchDeleteResponse(vsw.BatchDeleteResponseInput{Results: results})
}

// VacuumVolumeCheck
func VacuumVolumeCheckReqFromWire(b []byte) (*volume_server_pb.VacuumVolumeCheckRequest, error) {
	v, err := vsw.WrapVacuumVolumeCheckRequest(b)
	if err != nil {
		return nil, err
	}
	return &volume_server_pb.VacuumVolumeCheckRequest{VolumeId: v.VolumeID()}, nil
}
func VacuumVolumeCheckRespToWire(r *volume_server_pb.VacuumVolumeCheckResponse) []byte {
	return vsw.NewVacuumVolumeCheckResponse(vsw.VacuumVolumeCheckResponseInput{GarbageRatio: r.GarbageRatio})
}

// VacuumVolumeCommit
func VacuumVolumeCommitReqFromWire(b []byte) (*volume_server_pb.VacuumVolumeCommitRequest, error) {
	v, err := vsw.WrapVacuumVolumeCommitRequest(b)
	if err != nil {
		return nil, err
	}
	return &volume_server_pb.VacuumVolumeCommitRequest{VolumeId: v.VolumeID()}, nil
}
func VacuumVolumeCommitRespToWire(r *volume_server_pb.VacuumVolumeCommitResponse) []byte {
	return vsw.NewVacuumVolumeCommitResponse(vsw.VacuumVolumeCommitResponseInput{
		IsReadOnly: r.IsReadOnly, VolumeSize: r.VolumeSize,
	})
}

// VacuumVolumeCleanup
func VacuumVolumeCleanupReqFromWire(b []byte) (*volume_server_pb.VacuumVolumeCleanupRequest, error) {
	v, err := vsw.WrapVacuumVolumeCleanupRequest(b)
	if err != nil {
		return nil, err
	}
	return &volume_server_pb.VacuumVolumeCleanupRequest{VolumeId: v.VolumeID()}, nil
}
func VacuumVolumeCleanupRespToWire(_ *volume_server_pb.VacuumVolumeCleanupResponse) []byte {
	return vsw.NewVacuumVolumeCleanupResponse(vsw.VacuumVolumeCleanupResponseInput{})
}

// DeleteCollection
func DeleteCollectionReqFromWire(b []byte) (*volume_server_pb.DeleteCollectionRequest, error) {
	v, err := vsw.WrapDeleteCollectionRequest(b)
	if err != nil {
		return nil, err
	}
	return &volume_server_pb.DeleteCollectionRequest{Collection: v.Collection()}, nil
}
func DeleteCollectionRespToWire(_ *volume_server_pb.DeleteCollectionResponse) []byte {
	return vsw.NewDeleteCollectionResponse(vsw.DeleteCollectionResponseInput{})
}

// AllocateVolume
func AllocateVolumeReqFromWire(b []byte) (*volume_server_pb.AllocateVolumeRequest, error) {
	v, err := vsw.WrapAllocateVolumeRequest(b)
	if err != nil {
		return nil, err
	}
	return &volume_server_pb.AllocateVolumeRequest{
		VolumeId: v.VolumeID(), Collection: v.Collection(), Preallocate: v.Preallocate(),
		Replication: v.Replication(), Ttl: v.TTL(), MemoryMapMaxSizeMb: v.MemoryMapMaxSizeMb(),
		DiskType: v.DiskType(), Version: v.Version(),
	}, nil
}
func AllocateVolumeRespToWire(_ *volume_server_pb.AllocateVolumeResponse) []byte {
	return vsw.NewAllocateVolumeResponse(vsw.AllocateVolumeResponseInput{})
}

// VolumeSyncStatus
func VolumeSyncStatusReqFromWire(b []byte) (*volume_server_pb.VolumeSyncStatusRequest, error) {
	v, err := vsw.WrapVolumeSyncStatusRequest(b)
	if err != nil {
		return nil, err
	}
	return &volume_server_pb.VolumeSyncStatusRequest{VolumeId: v.VolumeID()}, nil
}
func VolumeSyncStatusRespToWire(r *volume_server_pb.VolumeSyncStatusResponse) []byte {
	return vsw.NewVolumeSyncStatusResponse(vsw.VolumeSyncStatusResponseInput{
		VolumeID: r.VolumeId, Collection: r.Collection, Replication: r.Replication,
		TTL: r.Ttl, TailOffset: r.TailOffset, CompactRevision: r.CompactRevision,
		IdxFileSize: r.IdxFileSize, Version: r.Version,
	})
}

// VolumeMount
func VolumeMountReqFromWire(b []byte) (*volume_server_pb.VolumeMountRequest, error) {
	v, err := vsw.WrapVolumeMountRequest(b)
	if err != nil {
		return nil, err
	}
	return &volume_server_pb.VolumeMountRequest{VolumeId: v.VolumeID()}, nil
}
func VolumeMountRespToWire(_ *volume_server_pb.VolumeMountResponse) []byte {
	return vsw.NewVolumeMountResponse(vsw.VolumeMountResponseInput{})
}

// VolumeUnmount
func VolumeUnmountReqFromWire(b []byte) (*volume_server_pb.VolumeUnmountRequest, error) {
	v, err := vsw.WrapVolumeUnmountRequest(b)
	if err != nil {
		return nil, err
	}
	return &volume_server_pb.VolumeUnmountRequest{VolumeId: v.VolumeID()}, nil
}
func VolumeUnmountRespToWire(_ *volume_server_pb.VolumeUnmountResponse) []byte {
	return vsw.NewVolumeUnmountResponse(vsw.VolumeUnmountResponseInput{})
}

// VolumeDelete
func VolumeDeleteReqFromWire(b []byte) (*volume_server_pb.VolumeDeleteRequest, error) {
	v, err := vsw.WrapVolumeDeleteRequest(b)
	if err != nil {
		return nil, err
	}
	return &volume_server_pb.VolumeDeleteRequest{
		VolumeId: v.VolumeID(), OnlyEmpty: v.OnlyEmpty(), KeepRemoteData: v.KeepRemoteData(),
	}, nil
}
func VolumeDeleteRespToWire(_ *volume_server_pb.VolumeDeleteResponse) []byte {
	return vsw.NewVolumeDeleteResponse(vsw.VolumeDeleteResponseInput{})
}

// VolumeMarkReadonly
func VolumeMarkReadonlyReqFromWire(b []byte) (*volume_server_pb.VolumeMarkReadonlyRequest, error) {
	v, err := vsw.WrapVolumeMarkReadonlyRequest(b)
	if err != nil {
		return nil, err
	}
	return &volume_server_pb.VolumeMarkReadonlyRequest{VolumeId: v.VolumeID(), Persist: v.Persist()}, nil
}
func VolumeMarkReadonlyRespToWire(_ *volume_server_pb.VolumeMarkReadonlyResponse) []byte {
	return vsw.NewVolumeMarkReadonlyResponse(vsw.VolumeMarkReadonlyResponseInput{})
}

// VolumeMarkWritable
func VolumeMarkWritableReqFromWire(b []byte) (*volume_server_pb.VolumeMarkWritableRequest, error) {
	v, err := vsw.WrapVolumeMarkWritableRequest(b)
	if err != nil {
		return nil, err
	}
	return &volume_server_pb.VolumeMarkWritableRequest{VolumeId: v.VolumeID()}, nil
}
func VolumeMarkWritableRespToWire(_ *volume_server_pb.VolumeMarkWritableResponse) []byte {
	return vsw.NewVolumeMarkWritableResponse(vsw.VolumeMarkWritableResponseInput{})
}

// VolumeConfigure
func VolumeConfigureReqFromWire(b []byte) (*volume_server_pb.VolumeConfigureRequest, error) {
	v, err := vsw.WrapVolumeConfigureRequest(b)
	if err != nil {
		return nil, err
	}
	return &volume_server_pb.VolumeConfigureRequest{VolumeId: v.VolumeID(), Replication: v.Replication()}, nil
}
func VolumeConfigureRespToWire(r *volume_server_pb.VolumeConfigureResponse) []byte {
	return vsw.NewVolumeConfigureResponse(vsw.VolumeConfigureResponseInput{Error: r.Error})
}

// VolumeStatus
func VolumeStatusReqFromWire(b []byte) (*volume_server_pb.VolumeStatusRequest, error) {
	v, err := vsw.WrapVolumeStatusRequest(b)
	if err != nil {
		return nil, err
	}
	return &volume_server_pb.VolumeStatusRequest{VolumeId: v.VolumeID()}, nil
}
func VolumeStatusRespToWire(r *volume_server_pb.VolumeStatusResponse) []byte {
	return vsw.NewVolumeStatusResponse(vsw.VolumeStatusResponseInput{
		IsReadOnly: r.IsReadOnly, VolumeSize: r.VolumeSize,
		FileCount: r.FileCount, FileDeletedCount: r.FileDeletedCount,
	})
}

// GetState
func GetStateReqFromWire(b []byte) (*volume_server_pb.GetStateRequest, error) {
	if _, err := vsw.WrapGetStateRequest(b); err != nil {
		return nil, err
	}
	return &volume_server_pb.GetStateRequest{}, nil
}
func GetStateRespToWire(r *volume_server_pb.GetStateResponse) []byte {
	return vsw.NewGetStateResponse(vsw.GetStateResponseInput{State: volumeServerStateToWire(r.State)})
}

// SetState
func SetStateReqFromWire(b []byte) (*volume_server_pb.SetStateRequest, error) {
	v, err := vsw.WrapSetStateRequest(b)
	if err != nil {
		return nil, err
	}
	req := &volume_server_pb.SetStateRequest{}
	if s, ok := v.State(); ok {
		req.State = volumeServerStateFromView(s)
	}
	return req, nil
}
func SetStateRespToWire(r *volume_server_pb.SetStateResponse) []byte {
	return vsw.NewSetStateResponse(vsw.SetStateResponseInput{State: volumeServerStateToWire(r.State)})
}

// ReadVolumeFileStatus
func ReadVolumeFileStatusReqFromWire(b []byte) (*volume_server_pb.ReadVolumeFileStatusRequest, error) {
	v, err := vsw.WrapReadVolumeFileStatusRequest(b)
	if err != nil {
		return nil, err
	}
	return &volume_server_pb.ReadVolumeFileStatusRequest{VolumeId: v.VolumeID()}, nil
}
func ReadVolumeFileStatusRespToWire(r *volume_server_pb.ReadVolumeFileStatusResponse) []byte {
	return vsw.NewReadVolumeFileStatusResponse(vsw.ReadVolumeFileStatusResponseInput{
		VolumeID:                r.VolumeId,
		IdxFileTimestampSeconds: r.IdxFileTimestampSeconds,
		IdxFileSize:             r.IdxFileSize,
		DatFileTimestampSeconds: r.DatFileTimestampSeconds,
		DatFileSize:             r.DatFileSize,
		FileCount:               r.FileCount,
		CompactionRevision:      r.CompactionRevision,
		Collection:              r.Collection,
		DiskType:                r.DiskType,
		VolumeInfo:              volumeInfoToWire(r.VolumeInfo),
		Version:                 r.Version,
	})
}

// ReadNeedleBlob
func ReadNeedleBlobReqFromWire(b []byte) (*volume_server_pb.ReadNeedleBlobRequest, error) {
	v, err := vsw.WrapReadNeedleBlobRequest(b)
	if err != nil {
		return nil, err
	}
	return &volume_server_pb.ReadNeedleBlobRequest{VolumeId: v.VolumeID(), Offset: v.Offset(), Size: v.Size()}, nil
}
func ReadNeedleBlobRespToWire(r *volume_server_pb.ReadNeedleBlobResponse) []byte {
	return vsw.NewReadNeedleBlobResponse(vsw.ReadNeedleBlobResponseInput{NeedleBlob: r.NeedleBlob})
}

// ReadNeedleMeta
func ReadNeedleMetaReqFromWire(b []byte) (*volume_server_pb.ReadNeedleMetaRequest, error) {
	v, err := vsw.WrapReadNeedleMetaRequest(b)
	if err != nil {
		return nil, err
	}
	return &volume_server_pb.ReadNeedleMetaRequest{
		VolumeId: v.VolumeID(), NeedleId: v.NeedleID(), Offset: v.Offset(), Size: v.Size(),
	}, nil
}
func ReadNeedleMetaRespToWire(r *volume_server_pb.ReadNeedleMetaResponse) []byte {
	return vsw.NewReadNeedleMetaResponse(vsw.ReadNeedleMetaResponseInput{
		Cookie: r.Cookie, LastModified: r.LastModified, Crc: r.Crc, TTL: r.Ttl, AppendAtNs: r.AppendAtNs,
	})
}

// WriteNeedleBlob
func WriteNeedleBlobReqFromWire(b []byte) (*volume_server_pb.WriteNeedleBlobRequest, error) {
	v, err := vsw.WrapWriteNeedleBlobRequest(b)
	if err != nil {
		return nil, err
	}
	return &volume_server_pb.WriteNeedleBlobRequest{
		VolumeId: v.VolumeID(), NeedleId: v.NeedleID(), Size: v.Size(), NeedleBlob: v.NeedleBlob(),
	}, nil
}
func WriteNeedleBlobRespToWire(_ *volume_server_pb.WriteNeedleBlobResponse) []byte {
	return vsw.NewWriteNeedleBlobResponse(vsw.WriteNeedleBlobResponseInput{})
}

// VolumeTailReceiver
func VolumeTailReceiverReqFromWire(b []byte) (*volume_server_pb.VolumeTailReceiverRequest, error) {
	v, err := vsw.WrapVolumeTailReceiverRequest(b)
	if err != nil {
		return nil, err
	}
	return &volume_server_pb.VolumeTailReceiverRequest{
		VolumeId: v.VolumeID(), SinceNs: v.SinceNs(), IdleTimeoutSeconds: v.IdleTimeoutSeconds(), SourceVolumeServer: v.SourceVolumeServer(),
	}, nil
}
func VolumeTailReceiverRespToWire(_ *volume_server_pb.VolumeTailReceiverResponse) []byte {
	return vsw.NewVolumeTailReceiverResponse(vsw.VolumeTailReceiverResponseInput{})
}

// VolumeEcShardsGenerate
func VolumeEcShardsGenerateReqFromWire(b []byte) (*volume_server_pb.VolumeEcShardsGenerateRequest, error) {
	v, err := vsw.WrapVolumeEcShardsGenerateRequest(b)
	if err != nil {
		return nil, err
	}
	return &volume_server_pb.VolumeEcShardsGenerateRequest{VolumeId: v.VolumeID(), Collection: v.Collection()}, nil
}
func VolumeEcShardsGenerateRespToWire(_ *volume_server_pb.VolumeEcShardsGenerateResponse) []byte {
	return vsw.NewVolumeEcShardsGenerateResponse(vsw.VolumeEcShardsGenerateResponseInput{})
}

// VolumeEcShardsRebuild
func VolumeEcShardsRebuildReqFromWire(b []byte) (*volume_server_pb.VolumeEcShardsRebuildRequest, error) {
	v, err := vsw.WrapVolumeEcShardsRebuildRequest(b)
	if err != nil {
		return nil, err
	}
	return &volume_server_pb.VolumeEcShardsRebuildRequest{
		VolumeId: v.VolumeID(), Collection: v.Collection(), UnsafeIgnoreSidecar: v.UnsafeIgnoreSidecar(),
	}, nil
}
func VolumeEcShardsRebuildRespToWire(r *volume_server_pb.VolumeEcShardsRebuildResponse) []byte {
	return vsw.NewVolumeEcShardsRebuildResponse(vsw.VolumeEcShardsRebuildResponseInput{RebuiltShardIDs: r.RebuiltShardIds})
}

// VolumeEcShardsCopy
func VolumeEcShardsCopyReqFromWire(b []byte) (*volume_server_pb.VolumeEcShardsCopyRequest, error) {
	v, err := vsw.WrapVolumeEcShardsCopyRequest(b)
	if err != nil {
		return nil, err
	}
	req := &volume_server_pb.VolumeEcShardsCopyRequest{
		VolumeId: v.VolumeID(), Collection: v.Collection(), CopyEcxFile: v.CopyEcxFile(),
		SourceDataNode: v.SourceDataNode(), CopyEcjFile: v.CopyEcjFile(), CopyVifFile: v.CopyVifFile(),
		DiskId: v.DiskID(), CopyEcsumFile: v.CopyEcsumFile(),
	}
	for i := 0; i < v.ShardIDsLen(); i++ {
		req.ShardIds = append(req.ShardIds, v.ShardIDAt(i))
	}
	return req, nil
}
func VolumeEcShardsCopyRespToWire(_ *volume_server_pb.VolumeEcShardsCopyResponse) []byte {
	return vsw.NewVolumeEcShardsCopyResponse(vsw.VolumeEcShardsCopyResponseInput{})
}

// VolumeEcShardsDelete
func VolumeEcShardsDeleteReqFromWire(b []byte) (*volume_server_pb.VolumeEcShardsDeleteRequest, error) {
	v, err := vsw.WrapVolumeEcShardsDeleteRequest(b)
	if err != nil {
		return nil, err
	}
	req := &volume_server_pb.VolumeEcShardsDeleteRequest{
		VolumeId: v.VolumeID(), Collection: v.Collection(),
		FullTeardown: v.FullTeardown(), EncodeTsNs: v.EncodeTsNs(),
	}
	for i := 0; i < v.ShardIDsLen(); i++ {
		req.ShardIds = append(req.ShardIds, v.ShardIDAt(i))
	}
	return req, nil
}
func VolumeEcShardsDeleteRespToWire(r *volume_server_pb.VolumeEcShardsDeleteResponse) []byte {
	return vsw.NewVolumeEcShardsDeleteResponse(vsw.VolumeEcShardsDeleteResponseInput{FullTeardownDone: r.FullTeardownDone})
}

// VolumeEcShardsMount
func VolumeEcShardsMountReqFromWire(b []byte) (*volume_server_pb.VolumeEcShardsMountRequest, error) {
	v, err := vsw.WrapVolumeEcShardsMountRequest(b)
	if err != nil {
		return nil, err
	}
	req := &volume_server_pb.VolumeEcShardsMountRequest{
		VolumeId: v.VolumeID(), Collection: v.Collection(), SourceDiskType: v.SourceDiskType(),
	}
	for i := 0; i < v.ShardIDsLen(); i++ {
		req.ShardIds = append(req.ShardIds, v.ShardIDAt(i))
	}
	return req, nil
}
func VolumeEcShardsMountRespToWire(_ *volume_server_pb.VolumeEcShardsMountResponse) []byte {
	return vsw.NewVolumeEcShardsMountResponse(vsw.VolumeEcShardsMountResponseInput{})
}

// VolumeEcShardsUnmount
func VolumeEcShardsUnmountReqFromWire(b []byte) (*volume_server_pb.VolumeEcShardsUnmountRequest, error) {
	v, err := vsw.WrapVolumeEcShardsUnmountRequest(b)
	if err != nil {
		return nil, err
	}
	req := &volume_server_pb.VolumeEcShardsUnmountRequest{VolumeId: v.VolumeID(), EncodeTsNs: v.EncodeTsNs()}
	for i := 0; i < v.ShardIDsLen(); i++ {
		req.ShardIds = append(req.ShardIds, v.ShardIDAt(i))
	}
	return req, nil
}
func VolumeEcShardsUnmountRespToWire(_ *volume_server_pb.VolumeEcShardsUnmountResponse) []byte {
	return vsw.NewVolumeEcShardsUnmountResponse(vsw.VolumeEcShardsUnmountResponseInput{})
}

// VolumeEcBlobDelete
func VolumeEcBlobDeleteReqFromWire(b []byte) (*volume_server_pb.VolumeEcBlobDeleteRequest, error) {
	v, err := vsw.WrapVolumeEcBlobDeleteRequest(b)
	if err != nil {
		return nil, err
	}
	return &volume_server_pb.VolumeEcBlobDeleteRequest{
		VolumeId: v.VolumeID(), Collection: v.Collection(), FileKey: v.FileKey(), Version: v.Version(),
	}, nil
}
func VolumeEcBlobDeleteRespToWire(_ *volume_server_pb.VolumeEcBlobDeleteResponse) []byte {
	return vsw.NewVolumeEcBlobDeleteResponse(vsw.VolumeEcBlobDeleteResponseInput{})
}

// VolumeEcShardsToVolume
func VolumeEcShardsToVolumeReqFromWire(b []byte) (*volume_server_pb.VolumeEcShardsToVolumeRequest, error) {
	v, err := vsw.WrapVolumeEcShardsToVolumeRequest(b)
	if err != nil {
		return nil, err
	}
	return &volume_server_pb.VolumeEcShardsToVolumeRequest{VolumeId: v.VolumeID(), Collection: v.Collection()}, nil
}
func VolumeEcShardsToVolumeRespToWire(_ *volume_server_pb.VolumeEcShardsToVolumeResponse) []byte {
	return vsw.NewVolumeEcShardsToVolumeResponse(vsw.VolumeEcShardsToVolumeResponseInput{})
}

// VolumeEcShardsInfo
func VolumeEcShardsInfoReqFromWire(b []byte) (*volume_server_pb.VolumeEcShardsInfoRequest, error) {
	v, err := vsw.WrapVolumeEcShardsInfoRequest(b)
	if err != nil {
		return nil, err
	}
	return &volume_server_pb.VolumeEcShardsInfoRequest{VolumeId: v.VolumeID()}, nil
}
func VolumeEcShardsInfoRespToWire(r *volume_server_pb.VolumeEcShardsInfoResponse) []byte {
	infos := make([][]byte, len(r.EcShardInfos))
	for i, info := range r.EcShardInfos {
		infos[i] = ecShardInfoToWire(info)
	}
	return vsw.NewVolumeEcShardsInfoResponse(vsw.VolumeEcShardsInfoResponseInput{
		EcShardInfos: infos, VolumeSize: r.VolumeSize, FileCount: r.FileCount, FileDeletedCount: r.FileDeletedCount,
	})
}

// VolumeServerStatus
func VolumeServerStatusReqFromWire(b []byte) (*volume_server_pb.VolumeServerStatusRequest, error) {
	if _, err := vsw.WrapVolumeServerStatusRequest(b); err != nil {
		return nil, err
	}
	return &volume_server_pb.VolumeServerStatusRequest{}, nil
}
func VolumeServerStatusRespToWire(r *volume_server_pb.VolumeServerStatusResponse) []byte {
	disks := make([][]byte, len(r.DiskStatuses))
	for i, d := range r.DiskStatuses {
		disks[i] = diskStatusToWire(d)
	}
	return vsw.NewVolumeServerStatusResponse(vsw.VolumeServerStatusResponseInput{
		DiskStatuses: disks,
		MemoryStatus: memStatusToWire(r.MemoryStatus),
		Version:      r.Version,
		DataCenter:   r.DataCenter,
		Rack:         r.Rack,
		State:        volumeServerStateToWire(r.State),
	})
}

// VolumeServerLeave
func VolumeServerLeaveReqFromWire(b []byte) (*volume_server_pb.VolumeServerLeaveRequest, error) {
	if _, err := vsw.WrapVolumeServerLeaveRequest(b); err != nil {
		return nil, err
	}
	return &volume_server_pb.VolumeServerLeaveRequest{}, nil
}
func VolumeServerLeaveRespToWire(_ *volume_server_pb.VolumeServerLeaveResponse) []byte {
	return vsw.NewVolumeServerLeaveResponse(vsw.VolumeServerLeaveResponseInput{})
}

// FetchAndWriteNeedle
func FetchAndWriteNeedleReqFromWire(b []byte) (*volume_server_pb.FetchAndWriteNeedleRequest, error) {
	v, err := vsw.WrapFetchAndWriteNeedleRequest(b)
	if err != nil {
		return nil, err
	}
	req := &volume_server_pb.FetchAndWriteNeedleRequest{
		VolumeId: v.VolumeID(), NeedleId: v.NeedleID(), Cookie: v.Cookie(), Offset: v.Offset(), Size: v.Size(),
		Auth: v.Auth(), DownloadConcurrency: v.DownloadConcurrency(),
	}
	for i := 0; i < v.ReplicasLen(); i++ {
		if rep, ok := v.ReplicaAt(i); ok {
			req.Replicas = append(req.Replicas, fetchReplicaFromView(rep))
		}
	}
	if conf, ok := v.RemoteConf(); ok {
		req.RemoteConf = remoteConfFromView(conf)
	}
	if loc, ok := v.RemoteLocation(); ok {
		req.RemoteLocation = remoteLocationFromView(loc)
	}
	return req, nil
}
func FetchAndWriteNeedleRespToWire(r *volume_server_pb.FetchAndWriteNeedleResponse) []byte {
	return vsw.NewFetchAndWriteNeedleResponse(vsw.FetchAndWriteNeedleResponseInput{ETag: r.ETag})
}

// ScrubVolume
func ScrubVolumeReqFromWire(b []byte) (*volume_server_pb.ScrubVolumeRequest, error) {
	v, err := vsw.WrapScrubVolumeRequest(b)
	if err != nil {
		return nil, err
	}
	req := &volume_server_pb.ScrubVolumeRequest{
		Mode: volume_server_pb.VolumeScrubMode(v.Mode()), MarkBrokenVolumesReadonly: v.MarkBrokenVolumesReadonly(),
	}
	for i := 0; i < v.VolumeIDsLen(); i++ {
		req.VolumeIds = append(req.VolumeIds, v.VolumeIDAt(i))
	}
	return req, nil
}
func ScrubVolumeRespToWire(r *volume_server_pb.ScrubVolumeResponse) []byte {
	return vsw.NewScrubVolumeResponse(vsw.ScrubVolumeResponseInput{
		TotalVolumes: r.TotalVolumes, TotalFiles: r.TotalFiles,
		BrokenVolumeIDs: r.BrokenVolumeIds, Details: r.Details,
	})
}

// ScrubEcVolume
func ScrubEcVolumeReqFromWire(b []byte) (*volume_server_pb.ScrubEcVolumeRequest, error) {
	v, err := vsw.WrapScrubEcVolumeRequest(b)
	if err != nil {
		return nil, err
	}
	req := &volume_server_pb.ScrubEcVolumeRequest{Mode: volume_server_pb.VolumeScrubMode(v.Mode())}
	for i := 0; i < v.VolumeIDsLen(); i++ {
		req.VolumeIds = append(req.VolumeIds, v.VolumeIDAt(i))
	}
	return req, nil
}
func ScrubEcVolumeRespToWire(r *volume_server_pb.ScrubEcVolumeResponse) []byte {
	infos := make([][]byte, len(r.BrokenShardInfos))
	for i, info := range r.BrokenShardInfos {
		infos[i] = ecShardInfoToWire(info)
	}
	return vsw.NewScrubEcVolumeResponse(vsw.ScrubEcVolumeResponseInput{
		TotalVolumes: r.TotalVolumes, TotalFiles: r.TotalFiles,
		BrokenVolumeIDs: r.BrokenVolumeIds, BrokenShardInfos: infos, Details: r.Details,
	})
}

// VolumeNeedleStatus
func VolumeNeedleStatusReqFromWire(b []byte) (*volume_server_pb.VolumeNeedleStatusRequest, error) {
	v, err := vsw.WrapVolumeNeedleStatusRequest(b)
	if err != nil {
		return nil, err
	}
	return &volume_server_pb.VolumeNeedleStatusRequest{VolumeId: v.VolumeID(), NeedleId: v.NeedleID()}, nil
}
func VolumeNeedleStatusRespToWire(r *volume_server_pb.VolumeNeedleStatusResponse) []byte {
	return vsw.NewVolumeNeedleStatusResponse(vsw.VolumeNeedleStatusResponseInput{
		NeedleID: r.NeedleId, Cookie: r.Cookie, Size: r.Size,
		LastModified: r.LastModified, Crc: r.Crc, TTL: r.Ttl,
	})
}

// Ping
func PingReqFromWire(b []byte) (*volume_server_pb.PingRequest, error) {
	v, err := vsw.WrapPingRequest(b)
	if err != nil {
		return nil, err
	}
	return &volume_server_pb.PingRequest{Target: v.Target(), TargetType: v.TargetType()}, nil
}
func PingRespToWire(r *volume_server_pb.PingResponse) []byte {
	return vsw.NewPingResponse(vsw.PingResponseInput{
		StartTimeNs: r.StartTimeNs, RemoteTimeNs: r.RemoteTimeNs, StopTimeNs: r.StopTimeNs,
	})
}

// --- streaming request decoders (view -> pb) + response item encoders (pb -> Input) ---

// VacuumVolumeCompact
func VacuumVolumeCompactReqFromView(v vsw.VacuumVolumeCompactRequest) *volume_server_pb.VacuumVolumeCompactRequest {
	return &volume_server_pb.VacuumVolumeCompactRequest{VolumeId: v.VolumeID(), Preallocate: v.Preallocate()}
}
func VacuumVolumeCompactRespToInput(r *volume_server_pb.VacuumVolumeCompactResponse) vsw.VacuumVolumeCompactResponseInput {
	return vsw.VacuumVolumeCompactResponseInput{ProcessedBytes: r.ProcessedBytes, LoadAvg1m: r.LoadAvg_1M}
}

// VolumeIncrementalCopy
func VolumeIncrementalCopyReqFromView(v vsw.VolumeIncrementalCopyRequest) *volume_server_pb.VolumeIncrementalCopyRequest {
	return &volume_server_pb.VolumeIncrementalCopyRequest{VolumeId: v.VolumeID(), SinceNs: v.SinceNs()}
}
func VolumeIncrementalCopyRespToInput(r *volume_server_pb.VolumeIncrementalCopyResponse) vsw.VolumeIncrementalCopyResponseInput {
	return vsw.VolumeIncrementalCopyResponseInput{FileContent: r.FileContent}
}

// VolumeCopy
func VolumeCopyReqFromView(v vsw.VolumeCopyRequest) *volume_server_pb.VolumeCopyRequest {
	return &volume_server_pb.VolumeCopyRequest{
		VolumeId: v.VolumeID(), Collection: v.Collection(), Replication: v.Replication(), Ttl: v.TTL(),
		SourceDataNode: v.SourceDataNode(), DiskType: v.DiskType(), IoBytePerSecond: v.IoBytePerSecond(),
	}
}
func VolumeCopyRespToInput(r *volume_server_pb.VolumeCopyResponse) vsw.VolumeCopyResponseInput {
	return vsw.VolumeCopyResponseInput{LastAppendAtNs: r.LastAppendAtNs, ProcessedBytes: r.ProcessedBytes}
}

// CopyFile
func CopyFileReqFromView(v vsw.CopyFileRequest) *volume_server_pb.CopyFileRequest {
	return &volume_server_pb.CopyFileRequest{
		VolumeId: v.VolumeID(), Ext: v.Ext(), CompactionRevision: v.CompactionRevision(), StopOffset: v.StopOffset(),
		Collection: v.Collection(), IsEcVolume: v.IsEcVolume(), IgnoreSourceFileNotFound: v.IgnoreSourceFileNotFound(),
	}
}
func CopyFileRespToInput(r *volume_server_pb.CopyFileResponse) vsw.CopyFileResponseInput {
	return vsw.CopyFileResponseInput{FileContent: r.FileContent, ModifiedTsNs: r.ModifiedTsNs}
}

// ReadAllNeedles
func ReadAllNeedlesReqFromView(v vsw.ReadAllNeedlesRequest) *volume_server_pb.ReadAllNeedlesRequest {
	req := &volume_server_pb.ReadAllNeedlesRequest{}
	for i := 0; i < v.VolumeIDsLen(); i++ {
		req.VolumeIds = append(req.VolumeIds, v.VolumeIDAt(i))
	}
	return req
}
func ReadAllNeedlesRespToInput(r *volume_server_pb.ReadAllNeedlesResponse) vsw.ReadAllNeedlesResponseInput {
	return vsw.ReadAllNeedlesResponseInput{
		VolumeID: r.VolumeId, NeedleID: r.NeedleId, Cookie: r.Cookie, NeedleBlob: r.NeedleBlob,
		NeedleBlobCompressed: r.NeedleBlobCompressed, LastModified: r.LastModified, Crc: r.Crc,
		Name: r.Name, Mime: r.Mime,
	}
}

// VolumeTailSender
func VolumeTailSenderReqFromView(v vsw.VolumeTailSenderRequest) *volume_server_pb.VolumeTailSenderRequest {
	return &volume_server_pb.VolumeTailSenderRequest{VolumeId: v.VolumeID(), SinceNs: v.SinceNs(), IdleTimeoutSeconds: v.IdleTimeoutSeconds()}
}
func VolumeTailSenderRespToInput(r *volume_server_pb.VolumeTailSenderResponse) vsw.VolumeTailSenderResponseInput {
	return vsw.VolumeTailSenderResponseInput{
		NeedleHeader: r.NeedleHeader, NeedleBody: r.NeedleBody, IsLastChunk: r.IsLastChunk, Version: r.Version,
	}
}

// VolumeEcShardRead
func VolumeEcShardReadReqFromView(v vsw.VolumeEcShardReadRequest) *volume_server_pb.VolumeEcShardReadRequest {
	return &volume_server_pb.VolumeEcShardReadRequest{
		VolumeId: v.VolumeID(), ShardId: v.ShardID(), Offset: v.Offset(), Size: v.Size(), FileKey: v.FileKey(), EncodeTsNs: v.EncodeTsNs(),
	}
}
func VolumeEcShardReadRespToInput(r *volume_server_pb.VolumeEcShardReadResponse) vsw.VolumeEcShardReadResponseInput {
	return vsw.VolumeEcShardReadResponseInput{Data: r.Data, IsDeleted: r.IsDeleted, EncodeTsNs: r.EncodeTsNs}
}

// VolumeTierMoveDatToRemote
func VolumeTierMoveDatToRemoteReqFromView(v vsw.VolumeTierMoveDatToRemoteRequest) *volume_server_pb.VolumeTierMoveDatToRemoteRequest {
	return &volume_server_pb.VolumeTierMoveDatToRemoteRequest{
		VolumeId: v.VolumeID(), Collection: v.Collection(), DestinationBackendName: v.DestinationBackendName(), KeepLocalDatFile: v.KeepLocalDatFile(),
	}
}
func VolumeTierMoveDatToRemoteRespToInput(r *volume_server_pb.VolumeTierMoveDatToRemoteResponse) vsw.VolumeTierMoveDatToRemoteResponseInput {
	return vsw.VolumeTierMoveDatToRemoteResponseInput{Processed: r.Processed, ProcessedPercentage: r.ProcessedPercentage}
}

// VolumeTierMoveDatFromRemote
func VolumeTierMoveDatFromRemoteReqFromView(v vsw.VolumeTierMoveDatFromRemoteRequest) *volume_server_pb.VolumeTierMoveDatFromRemoteRequest {
	return &volume_server_pb.VolumeTierMoveDatFromRemoteRequest{
		VolumeId: v.VolumeID(), Collection: v.Collection(), KeepRemoteDatFile: v.KeepRemoteDatFile(),
	}
}
func VolumeTierMoveDatFromRemoteRespToInput(r *volume_server_pb.VolumeTierMoveDatFromRemoteResponse) vsw.VolumeTierMoveDatFromRemoteResponseInput {
	return vsw.VolumeTierMoveDatFromRemoteResponseInput{Processed: r.Processed, ProcessedPercentage: r.ProcessedPercentage}
}

// Query
func QueryReqFromView(v vsw.QueryRequest) *volume_server_pb.QueryRequest {
	req := &volume_server_pb.QueryRequest{}
	for i := 0; i < v.SelectionsLen(); i++ {
		req.Selections = append(req.Selections, v.SelectionAt(i))
	}
	for i := 0; i < v.FromFileIDsLen(); i++ {
		req.FromFileIds = append(req.FromFileIds, v.FromFileIDAt(i))
	}
	if f, ok := v.Filter(); ok {
		req.Filter = queryFilterFromView(f)
	}
	if s, ok := v.InputSerialization(); ok {
		req.InputSerialization = queryInputSerializationFromView(s)
	}
	if s, ok := v.OutputSerialization(); ok {
		req.OutputSerialization = queryOutputSerializationFromView(s)
	}
	return req
}
func QueriedStripeToInput(r *volume_server_pb.QueriedStripe) vsw.QueriedStripeInput {
	return vsw.QueriedStripeInput{Records: r.Records}
}

// ReceiveFile (client-streaming): each inbound frame -> pb request (oneof)
func ReceiveFileReqFromView(v vsw.ReceiveFileRequest) *volume_server_pb.ReceiveFileRequest {
	switch v.Data() {
	case vsw.ReceiveFileDataInfo:
		req := &volume_server_pb.ReceiveFileRequest{}
		if info, ok := v.Info(); ok {
			req.Data = &volume_server_pb.ReceiveFileRequest_Info{Info: &volume_server_pb.ReceiveFileInfo{
				VolumeId: info.VolumeID(), Ext: info.Ext(), Collection: info.Collection(),
				IsEcVolume: info.IsEcVolume(), ShardId: info.ShardID(), FileSize: info.FileSize(), DiskId: info.DiskID(),
			}}
		}
		return req
	case vsw.ReceiveFileDataFileContent:
		fc, _ := v.FileContent()
		return &volume_server_pb.ReceiveFileRequest{Data: &volume_server_pb.ReceiveFileRequest_FileContent{FileContent: fc}}
	default:
		return &volume_server_pb.ReceiveFileRequest{}
	}
}
func ReceiveFileRespToInput(r *volume_server_pb.ReceiveFileResponse) vsw.ReceiveFileResponseInput {
	return vsw.ReceiveFileResponseInput{BytesWritten: r.BytesWritten, Error: r.Error}
}

// --- Query nested request decoders (inverse of rpc.go's *ToWire) ---

func queryFilterFromView(v vsw.QueryFilter) *volume_server_pb.QueryRequest_Filter {
	return &volume_server_pb.QueryRequest_Filter{Field: v.Field(), Operand: v.Operand(), Value: v.Value()}
}

func queryInputSerializationFromView(v vsw.QueryInputSerialization) *volume_server_pb.QueryRequest_InputSerialization {
	s := &volume_server_pb.QueryRequest_InputSerialization{CompressionType: v.CompressionType()}
	if c, ok := v.CSVInput(); ok {
		s.CsvInput = &volume_server_pb.QueryRequest_InputSerialization_CSVInput{
			FileHeaderInfo: c.FileHeaderInfo(), RecordDelimiter: c.RecordDelimiter(), FieldDelimiter: c.FieldDelimiter(),
			QuoteCharacter: c.QuoteCharacter(), QuoteEscapeCharacter: c.QuoteEscapeCharacter(), Comments: c.Comments(),
			AllowQuotedRecordDelimiter: c.AllowQuotedRecordDelimiter(),
		}
	}
	if j, ok := v.JSONInput(); ok {
		s.JsonInput = &volume_server_pb.QueryRequest_InputSerialization_JSONInput{Type: j.Type()}
	}
	if _, ok := v.ParquetInput(); ok {
		s.ParquetInput = &volume_server_pb.QueryRequest_InputSerialization_ParquetInput{}
	}
	return s
}

func queryOutputSerializationFromView(v vsw.QueryOutputSerialization) *volume_server_pb.QueryRequest_OutputSerialization {
	s := &volume_server_pb.QueryRequest_OutputSerialization{}
	if c, ok := v.CSVOutput(); ok {
		s.CsvOutput = &volume_server_pb.QueryRequest_OutputSerialization_CSVOutput{
			QuoteFields: c.QuoteFields(), RecordDelimiter: c.RecordDelimiter(), FieldDelimiter: c.FieldDelimiter(),
			QuoteCharacter: c.QuoteCharacter(), QuoteEscapeCharacter: c.QuoteEscapeCharacter(),
		}
	}
	if j, ok := v.JSONOutput(); ok {
		s.JsonOutput = &volume_server_pb.QueryRequest_OutputSerialization_JSONOutput{RecordDelimiter: j.RecordDelimiter()}
	}
	return s
}
