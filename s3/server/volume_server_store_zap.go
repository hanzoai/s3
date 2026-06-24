package s3server

import (
	"context"

	"github.com/hanzoai/s3/s3/pb/remote_pb"
	"github.com/hanzoai/s3/s3/pb/volume_server_pb"
	remotewire "github.com/hanzoai/s3/s3/wire/remote"
	volume_serverwire "github.com/hanzoai/s3/s3/wire/volume_server"
)

// volume_server_store_zap.go adapts *VolumeServer (whose methods carry the gRPC
// (ctx, *pb.Req)->(*pb.Resp, error) shape) to the native ZAP
// volume_serverwire.VolumeServerStore contract: every unary method takes one
// request buffer (Wrap<Req>) and returns one response buffer (New<Resp>). The
// existing handler bodies stay in their lane on the pb domain model — the store
// adapter is a thin conv boundary (wire view -> pb request, pb response -> wire
// builder). This is the volume_server analogue of agentconv: it lets the ZAP
// transport drive the real volume server without gRPC, while the storage layer
// keeps its on-disk protobuf domain types untouched.
//
// Additive strangler: this binds the wire contract to the live handlers without
// removing the gRPC registration, so the tree stays green and the running server
// keeps serving until the registration is flipped to transport.ListenStream and
// the client callbacks move to the wire Client in lockstep.
//
// This file binds the 37 UNARY RPCs. The 11 streaming RPCs (VacuumVolumeCompact,
// VolumeIncrementalCopy, VolumeCopy, CopyFile, ReadAllNeedles, VolumeTailSender,
// VolumeEcShardRead, VolumeTierMoveDatToRemote, VolumeTierMoveDatFromRemote,
// Query, ReceiveFile) are NOT part of this store: they are hand-wired over
// transport.OpenStream / transport.ListenStream when the registration flips,
// exactly like the mq_agent and plugin/worker stream adapters.

// --- nested message converters ---
// Each builds a standalone wire sub-buffer (for embedMessage/addObjectList on
// the response side) or reads a wire view into a pb domain value (on the request
// side). They stay symmetric with the generated *_zap.go builders/views.

// stateProtoToWire builds a wire VolumeServerState buffer from the pb state.
func stateProtoToWire(st *volume_server_pb.VolumeServerState) []byte {
	if st == nil {
		return nil
	}
	return volume_serverwire.NewVolumeServerState(volume_serverwire.VolumeServerStateInput{
		Maintenance: st.GetMaintenance(),
		Version:     st.GetVersion(),
	})
}

// stateWireToProto reads a wire VolumeServerState view into a pb state.
func stateWireToProto(v volume_serverwire.VolumeServerState) *volume_server_pb.VolumeServerState {
	return &volume_server_pb.VolumeServerState{
		Maintenance: v.Maintenance(),
		Version:     v.Version(),
	}
}

// diskStatusProtoToWire builds a wire DiskStatus buffer from a pb DiskStatus.
func diskStatusProtoToWire(d *volume_server_pb.DiskStatus) []byte {
	if d == nil {
		return nil
	}
	return volume_serverwire.NewDiskStatus(volume_serverwire.DiskStatusInput{
		Dir:         d.GetDir(),
		All:         d.GetAll(),
		Used:        d.GetUsed(),
		Free:        d.GetFree(),
		PercentFree: d.GetPercentFree(),
		PercentUsed: d.GetPercentUsed(),
		DiskType:    d.GetDiskType(),
		Error:       d.GetError(),
	})
}

// memStatusProtoToWire builds a wire MemStatus buffer from a pb MemStatus.
func memStatusProtoToWire(m *volume_server_pb.MemStatus) []byte {
	if m == nil {
		return nil
	}
	return volume_serverwire.NewMemStatus(volume_serverwire.MemStatusInput{
		Goroutines: m.GetGoroutines(),
		All:        m.GetAll(),
		Used:       m.GetUsed(),
		Free:       m.GetFree(),
		Self:       m.GetSelf(),
		Heap:       m.GetHeap(),
		Stack:      m.GetStack(),
	})
}

// remoteFileProtoToWire builds a wire RemoteFile buffer from a pb RemoteFile.
func remoteFileProtoToWire(f *volume_server_pb.RemoteFile) []byte {
	if f == nil {
		return nil
	}
	return volume_serverwire.NewRemoteFile(volume_serverwire.RemoteFileInput{
		BackendType:  f.GetBackendType(),
		BackendID:    f.GetBackendId(),
		Key:          f.GetKey(),
		Offset:       f.GetOffset(),
		FileSize:     f.GetFileSize(),
		ModifiedTime: f.GetModifiedTime(),
		Extension:    f.GetExtension(),
	})
}

// volumeInfoProtoToWire builds a wire VolumeInfo buffer from a pb VolumeInfo.
func volumeInfoProtoToWire(vi *volume_server_pb.VolumeInfo) []byte {
	if vi == nil {
		return nil
	}
	var files [][]byte
	for _, f := range vi.GetFiles() {
		files = append(files, remoteFileProtoToWire(f))
	}
	var ecShardConfig []byte
	if c := vi.GetEcShardConfig(); c != nil {
		ecShardConfig = volume_serverwire.NewEcShardConfig(volume_serverwire.EcShardConfigInput{
			DataShards:   c.GetDataShards(),
			ParityShards: c.GetParityShards(),
			EncodeTsNs:   c.GetEncodeTsNs(),
		})
	}
	return volume_serverwire.NewVolumeInfo(volume_serverwire.VolumeInfoInput{
		Files:         files,
		Version:       vi.GetVersion(),
		Replication:   vi.GetReplication(),
		BytesOffset:   vi.GetBytesOffset(),
		DatFileSize:   vi.GetDatFileSize(),
		ExpireAtSec:   vi.GetExpireAtSec(),
		ReadOnly:      vi.GetReadOnly(),
		EcShardConfig: ecShardConfig,
	})
}

// ecShardInfoProtoToWire builds a wire EcShardInfo buffer from a pb EcShardInfo.
func ecShardInfoProtoToWire(e *volume_server_pb.EcShardInfo) []byte {
	if e == nil {
		return nil
	}
	return volume_serverwire.NewEcShardInfo(volume_serverwire.EcShardInfoInput{
		ShardID:    e.GetShardId(),
		Size:       e.GetSize(),
		Collection: e.GetCollection(),
		VolumeID:   e.GetVolumeId(),
	})
}

// deleteResultProtoToWire builds a wire DeleteResult buffer from a pb result.
func deleteResultProtoToWire(r *volume_server_pb.DeleteResult) []byte {
	if r == nil {
		return nil
	}
	return volume_serverwire.NewDeleteResult(volume_serverwire.DeleteResultInput{
		FileID:  r.GetFileId(),
		Status:  r.GetStatus(),
		Error:   r.GetError(),
		Size:    r.GetSize(),
		Version: r.GetVersion(),
	})
}

// remoteConfWireToProto reads a wire RemoteConf view into a pb remote_pb.RemoteConf.
func remoteConfWireToProto(c remotewire.RemoteConf) *remote_pb.RemoteConf {
	return &remote_pb.RemoteConf{
		Type:                            c.Type(),
		Name:                            c.Name(),
		S3AccessKey:                     c.S3AccessKey(),
		S3SecretKey:                     c.S3SecretKey(),
		S3Region:                        c.S3Region(),
		S3Endpoint:                      c.S3Endpoint(),
		S3StorageClass:                  c.S3StorageClass(),
		GcsGoogleApplicationCredentials: c.GcsGoogleApplicationCredentials(),
		GcsProjectId:                    c.GcsProjectID(),
		AzureAccountName:                c.AzureAccountName(),
		AzureAccountKey:                 c.AzureAccountKey(),
		BackblazeKeyId:                  c.BackblazeKeyID(),
		BackblazeApplicationKey:         c.BackblazeApplicationKey(),
		BackblazeEndpoint:               c.BackblazeEndpoint(),
		BackblazeRegion:                 c.BackblazeRegion(),
		AliyunAccessKey:                 c.AliyunAccessKey(),
		AliyunSecretKey:                 c.AliyunSecretKey(),
		AliyunEndpoint:                  c.AliyunEndpoint(),
		AliyunRegion:                    c.AliyunRegion(),
		TencentSecretId:                 c.TencentSecretID(),
		TencentSecretKey:                c.TencentSecretKey(),
		TencentEndpoint:                 c.TencentEndpoint(),
		BaiduAccessKey:                  c.BaiduAccessKey(),
		BaiduSecretKey:                  c.BaiduSecretKey(),
		BaiduEndpoint:                   c.BaiduEndpoint(),
		BaiduRegion:                     c.BaiduRegion(),
		WasabiAccessKey:                 c.WasabiAccessKey(),
		WasabiSecretKey:                 c.WasabiSecretKey(),
		WasabiEndpoint:                  c.WasabiEndpoint(),
		WasabiRegion:                    c.WasabiRegion(),
		FilebaseAccessKey:               c.FilebaseAccessKey(),
		FilebaseSecretKey:               c.FilebaseSecretKey(),
		FilebaseEndpoint:                c.FilebaseEndpoint(),
		StorjAccessKey:                  c.StorjAccessKey(),
		StorjSecretKey:                  c.StorjSecretKey(),
		StorjEndpoint:                   c.StorjEndpoint(),
		ContaboAccessKey:                c.ContaboAccessKey(),
		ContaboSecretKey:                c.ContaboSecretKey(),
		ContaboEndpoint:                 c.ContaboEndpoint(),
		ContaboRegion:                   c.ContaboRegion(),
	}
}

// remoteLocationWireToProto reads a wire RemoteStorageLocation view into pb.
func remoteLocationWireToProto(l remotewire.RemoteStorageLocation) *remote_pb.RemoteStorageLocation {
	return &remote_pb.RemoteStorageLocation{
		Name:                   l.Name(),
		Bucket:                 l.Bucket(),
		Path:                   l.Path(),
		ListingCacheTtlSeconds: l.ListingCacheTTLSeconds(),
	}
}

// volumeServerStore wraps *VolumeServer and satisfies
// volume_serverwire.VolumeServerStore. It is a distinct type (not *VolumeServer
// itself) because the live handlers keep the gRPC (ctx, *pb.Req)->(*pb.Resp)
// signatures, which collide name-for-name with the store's
// (req []byte)->([]byte) contract — exactly the split mq_agentwire's
// unaryHandler uses. Each method Wraps the request view, delegates to the real
// gRPC-shaped handler on the embedded *VolumeServer, then Builds the response
// buffer.
type volumeServerStore struct{ vs *VolumeServer }

// BatchDelete implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) BatchDelete(req []byte) ([]byte, error) {
	view, err := volume_serverwire.WrapBatchDeleteRequest(req)
	if err != nil {
		return nil, err
	}
	fileIDs := make([]string, view.FileIDsLen())
	for i := range fileIDs {
		fileIDs[i] = view.FileIDAt(i)
	}
	resp, err := s.vs.BatchDelete(context.Background(), &volume_server_pb.BatchDeleteRequest{
		FileIds:         fileIDs,
		SkipCookieCheck: view.SkipCookieCheck(),
	})
	if err != nil {
		return nil, err
	}
	var results [][]byte
	for _, r := range resp.GetResults() {
		results = append(results, deleteResultProtoToWire(r))
	}
	return volume_serverwire.NewBatchDeleteResponse(volume_serverwire.BatchDeleteResponseInput{
		Results: results,
	}), nil
}

// VacuumVolumeCheck implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) VacuumVolumeCheck(req []byte) ([]byte, error) {
	view, err := volume_serverwire.WrapVacuumVolumeCheckRequest(req)
	if err != nil {
		return nil, err
	}
	resp, err := s.vs.VacuumVolumeCheck(context.Background(), &volume_server_pb.VacuumVolumeCheckRequest{
		VolumeId: view.VolumeID(),
	})
	if err != nil {
		return nil, err
	}
	return volume_serverwire.NewVacuumVolumeCheckResponse(volume_serverwire.VacuumVolumeCheckResponseInput{
		GarbageRatio: resp.GetGarbageRatio(),
	}), nil
}

// VacuumVolumeCommit implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) VacuumVolumeCommit(req []byte) ([]byte, error) {
	view, err := volume_serverwire.WrapVacuumVolumeCommitRequest(req)
	if err != nil {
		return nil, err
	}
	resp, err := s.vs.VacuumVolumeCommit(context.Background(), &volume_server_pb.VacuumVolumeCommitRequest{
		VolumeId: view.VolumeID(),
	})
	if err != nil {
		return nil, err
	}
	return volume_serverwire.NewVacuumVolumeCommitResponse(volume_serverwire.VacuumVolumeCommitResponseInput{
		IsReadOnly: resp.GetIsReadOnly(),
		VolumeSize: resp.GetVolumeSize(),
	}), nil
}

// VacuumVolumeCleanup implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) VacuumVolumeCleanup(req []byte) ([]byte, error) {
	view, err := volume_serverwire.WrapVacuumVolumeCleanupRequest(req)
	if err != nil {
		return nil, err
	}
	_, err = s.vs.VacuumVolumeCleanup(context.Background(), &volume_server_pb.VacuumVolumeCleanupRequest{
		VolumeId: view.VolumeID(),
	})
	if err != nil {
		return nil, err
	}
	return volume_serverwire.NewVacuumVolumeCleanupResponse(volume_serverwire.VacuumVolumeCleanupResponseInput{}), nil
}

// DeleteCollection implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) DeleteCollection(req []byte) ([]byte, error) {
	view, err := volume_serverwire.WrapDeleteCollectionRequest(req)
	if err != nil {
		return nil, err
	}
	_, err = s.vs.DeleteCollection(context.Background(), &volume_server_pb.DeleteCollectionRequest{
		Collection: view.Collection(),
	})
	if err != nil {
		return nil, err
	}
	return volume_serverwire.NewDeleteCollectionResponse(volume_serverwire.DeleteCollectionResponseInput{}), nil
}

// AllocateVolume implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) AllocateVolume(req []byte) ([]byte, error) {
	view, err := volume_serverwire.WrapAllocateVolumeRequest(req)
	if err != nil {
		return nil, err
	}
	_, err = s.vs.AllocateVolume(context.Background(), &volume_server_pb.AllocateVolumeRequest{
		VolumeId:           view.VolumeID(),
		Collection:         view.Collection(),
		Preallocate:        view.Preallocate(),
		Replication:        view.Replication(),
		Ttl:                view.TTL(),
		MemoryMapMaxSizeMb: view.MemoryMapMaxSizeMb(),
		DiskType:           view.DiskType(),
		Version:            view.Version(),
	})
	if err != nil {
		return nil, err
	}
	return volume_serverwire.NewAllocateVolumeResponse(volume_serverwire.AllocateVolumeResponseInput{}), nil
}

// VolumeSyncStatus implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) VolumeSyncStatus(req []byte) ([]byte, error) {
	view, err := volume_serverwire.WrapVolumeSyncStatusRequest(req)
	if err != nil {
		return nil, err
	}
	resp, err := s.vs.VolumeSyncStatus(context.Background(), &volume_server_pb.VolumeSyncStatusRequest{
		VolumeId: view.VolumeID(),
	})
	if err != nil {
		return nil, err
	}
	return volume_serverwire.NewVolumeSyncStatusResponse(volume_serverwire.VolumeSyncStatusResponseInput{
		VolumeID:        resp.GetVolumeId(),
		Collection:      resp.GetCollection(),
		Replication:     resp.GetReplication(),
		TTL:             resp.GetTtl(),
		TailOffset:      resp.GetTailOffset(),
		CompactRevision: resp.GetCompactRevision(),
		IdxFileSize:     resp.GetIdxFileSize(),
		Version:         resp.GetVersion(),
	}), nil
}

// VolumeMount implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) VolumeMount(req []byte) ([]byte, error) {
	view, err := volume_serverwire.WrapVolumeMountRequest(req)
	if err != nil {
		return nil, err
	}
	_, err = s.vs.VolumeMount(context.Background(), &volume_server_pb.VolumeMountRequest{
		VolumeId: view.VolumeID(),
	})
	if err != nil {
		return nil, err
	}
	return volume_serverwire.NewVolumeMountResponse(volume_serverwire.VolumeMountResponseInput{}), nil
}

// VolumeUnmount implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) VolumeUnmount(req []byte) ([]byte, error) {
	view, err := volume_serverwire.WrapVolumeUnmountRequest(req)
	if err != nil {
		return nil, err
	}
	_, err = s.vs.VolumeUnmount(context.Background(), &volume_server_pb.VolumeUnmountRequest{
		VolumeId: view.VolumeID(),
	})
	if err != nil {
		return nil, err
	}
	return volume_serverwire.NewVolumeUnmountResponse(volume_serverwire.VolumeUnmountResponseInput{}), nil
}

// VolumeDelete implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) VolumeDelete(req []byte) ([]byte, error) {
	view, err := volume_serverwire.WrapVolumeDeleteRequest(req)
	if err != nil {
		return nil, err
	}
	_, err = s.vs.VolumeDelete(context.Background(), &volume_server_pb.VolumeDeleteRequest{
		VolumeId:       view.VolumeID(),
		OnlyEmpty:      view.OnlyEmpty(),
		KeepRemoteData: view.KeepRemoteData(),
	})
	if err != nil {
		return nil, err
	}
	return volume_serverwire.NewVolumeDeleteResponse(volume_serverwire.VolumeDeleteResponseInput{}), nil
}

// VolumeMarkReadonly implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) VolumeMarkReadonly(req []byte) ([]byte, error) {
	view, err := volume_serverwire.WrapVolumeMarkReadonlyRequest(req)
	if err != nil {
		return nil, err
	}
	_, err = s.vs.VolumeMarkReadonly(context.Background(), &volume_server_pb.VolumeMarkReadonlyRequest{
		VolumeId: view.VolumeID(),
		Persist:  view.Persist(),
	})
	if err != nil {
		return nil, err
	}
	return volume_serverwire.NewVolumeMarkReadonlyResponse(volume_serverwire.VolumeMarkReadonlyResponseInput{}), nil
}

// VolumeMarkWritable implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) VolumeMarkWritable(req []byte) ([]byte, error) {
	view, err := volume_serverwire.WrapVolumeMarkWritableRequest(req)
	if err != nil {
		return nil, err
	}
	_, err = s.vs.VolumeMarkWritable(context.Background(), &volume_server_pb.VolumeMarkWritableRequest{
		VolumeId: view.VolumeID(),
	})
	if err != nil {
		return nil, err
	}
	return volume_serverwire.NewVolumeMarkWritableResponse(volume_serverwire.VolumeMarkWritableResponseInput{}), nil
}

// VolumeConfigure implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) VolumeConfigure(req []byte) ([]byte, error) {
	view, err := volume_serverwire.WrapVolumeConfigureRequest(req)
	if err != nil {
		return nil, err
	}
	resp, err := s.vs.VolumeConfigure(context.Background(), &volume_server_pb.VolumeConfigureRequest{
		VolumeId:    view.VolumeID(),
		Replication: view.Replication(),
	})
	if err != nil {
		return nil, err
	}
	return volume_serverwire.NewVolumeConfigureResponse(volume_serverwire.VolumeConfigureResponseInput{
		Error: resp.GetError(),
	}), nil
}

// VolumeStatus implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) VolumeStatus(req []byte) ([]byte, error) {
	view, err := volume_serverwire.WrapVolumeStatusRequest(req)
	if err != nil {
		return nil, err
	}
	resp, err := s.vs.VolumeStatus(context.Background(), &volume_server_pb.VolumeStatusRequest{
		VolumeId: view.VolumeID(),
	})
	if err != nil {
		return nil, err
	}
	return volume_serverwire.NewVolumeStatusResponse(volume_serverwire.VolumeStatusResponseInput{
		IsReadOnly:       resp.GetIsReadOnly(),
		VolumeSize:       resp.GetVolumeSize(),
		FileCount:        resp.GetFileCount(),
		FileDeletedCount: resp.GetFileDeletedCount(),
	}), nil
}

// GetState implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) GetState(req []byte) ([]byte, error) {
	resp, err := s.vs.GetState(context.Background(), &volume_server_pb.GetStateRequest{})
	if err != nil {
		return nil, err
	}
	return volume_serverwire.NewGetStateResponse(volume_serverwire.GetStateResponseInput{
		State: stateProtoToWire(resp.GetState()),
	}), nil
}

// SetState implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) SetState(req []byte) ([]byte, error) {
	view, err := volume_serverwire.WrapSetStateRequest(req)
	if err != nil {
		return nil, err
	}
	var pbState *volume_server_pb.VolumeServerState
	if sv, ok := view.State(); ok {
		pbState = stateWireToProto(sv)
	}
	resp, err := s.vs.SetState(context.Background(), &volume_server_pb.SetStateRequest{State: pbState})
	if err != nil {
		return nil, err
	}
	return volume_serverwire.NewSetStateResponse(volume_serverwire.SetStateResponseInput{
		State: stateProtoToWire(resp.GetState()),
	}), nil
}

// ReadVolumeFileStatus implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) ReadVolumeFileStatus(req []byte) ([]byte, error) {
	view, err := volume_serverwire.WrapReadVolumeFileStatusRequest(req)
	if err != nil {
		return nil, err
	}
	resp, err := s.vs.ReadVolumeFileStatus(context.Background(), &volume_server_pb.ReadVolumeFileStatusRequest{
		VolumeId: view.VolumeID(),
	})
	if err != nil {
		return nil, err
	}
	return volume_serverwire.NewReadVolumeFileStatusResponse(volume_serverwire.ReadVolumeFileStatusResponseInput{
		VolumeID:                resp.GetVolumeId(),
		IdxFileTimestampSeconds: resp.GetIdxFileTimestampSeconds(),
		IdxFileSize:             resp.GetIdxFileSize(),
		DatFileTimestampSeconds: resp.GetDatFileTimestampSeconds(),
		DatFileSize:             resp.GetDatFileSize(),
		FileCount:               resp.GetFileCount(),
		CompactionRevision:      resp.GetCompactionRevision(),
		Collection:              resp.GetCollection(),
		DiskType:                resp.GetDiskType(),
		VolumeInfo:              volumeInfoProtoToWire(resp.GetVolumeInfo()),
		Version:                 resp.GetVersion(),
	}), nil
}

// ReadNeedleBlob implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) ReadNeedleBlob(req []byte) ([]byte, error) {
	view, err := volume_serverwire.WrapReadNeedleBlobRequest(req)
	if err != nil {
		return nil, err
	}
	resp, err := s.vs.ReadNeedleBlob(context.Background(), &volume_server_pb.ReadNeedleBlobRequest{
		VolumeId: view.VolumeID(),
		Offset:   view.Offset(),
		Size:     view.Size(),
	})
	if err != nil {
		return nil, err
	}
	return volume_serverwire.NewReadNeedleBlobResponse(volume_serverwire.ReadNeedleBlobResponseInput{
		NeedleBlob: resp.GetNeedleBlob(),
	}), nil
}

// ReadNeedleMeta implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) ReadNeedleMeta(req []byte) ([]byte, error) {
	view, err := volume_serverwire.WrapReadNeedleMetaRequest(req)
	if err != nil {
		return nil, err
	}
	resp, err := s.vs.ReadNeedleMeta(context.Background(), &volume_server_pb.ReadNeedleMetaRequest{
		VolumeId: view.VolumeID(),
		NeedleId: view.NeedleID(),
		Offset:   view.Offset(),
		Size:     view.Size(),
	})
	if err != nil {
		return nil, err
	}
	return volume_serverwire.NewReadNeedleMetaResponse(volume_serverwire.ReadNeedleMetaResponseInput{
		Cookie:       resp.GetCookie(),
		LastModified: resp.GetLastModified(),
		Crc:          resp.GetCrc(),
		TTL:          resp.GetTtl(),
		AppendAtNs:   resp.GetAppendAtNs(),
	}), nil
}

// WriteNeedleBlob implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) WriteNeedleBlob(req []byte) ([]byte, error) {
	view, err := volume_serverwire.WrapWriteNeedleBlobRequest(req)
	if err != nil {
		return nil, err
	}
	_, err = s.vs.WriteNeedleBlob(context.Background(), &volume_server_pb.WriteNeedleBlobRequest{
		VolumeId:   view.VolumeID(),
		NeedleId:   view.NeedleID(),
		Size:       view.Size(),
		NeedleBlob: view.NeedleBlob(),
	})
	if err != nil {
		return nil, err
	}
	return volume_serverwire.NewWriteNeedleBlobResponse(volume_serverwire.WriteNeedleBlobResponseInput{}), nil
}

// VolumeTailReceiver implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) VolumeTailReceiver(req []byte) ([]byte, error) {
	view, err := volume_serverwire.WrapVolumeTailReceiverRequest(req)
	if err != nil {
		return nil, err
	}
	_, err = s.vs.VolumeTailReceiver(context.Background(), &volume_server_pb.VolumeTailReceiverRequest{
		VolumeId:           view.VolumeID(),
		SinceNs:            view.SinceNs(),
		IdleTimeoutSeconds: view.IdleTimeoutSeconds(),
		SourceVolumeServer: view.SourceVolumeServer(),
	})
	if err != nil {
		return nil, err
	}
	return volume_serverwire.NewVolumeTailReceiverResponse(volume_serverwire.VolumeTailReceiverResponseInput{}), nil
}

// VolumeEcShardsGenerate implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) VolumeEcShardsGenerate(req []byte) ([]byte, error) {
	view, err := volume_serverwire.WrapVolumeEcShardsGenerateRequest(req)
	if err != nil {
		return nil, err
	}
	_, err = s.vs.VolumeEcShardsGenerate(context.Background(), &volume_server_pb.VolumeEcShardsGenerateRequest{
		VolumeId:   view.VolumeID(),
		Collection: view.Collection(),
	})
	if err != nil {
		return nil, err
	}
	return volume_serverwire.NewVolumeEcShardsGenerateResponse(volume_serverwire.VolumeEcShardsGenerateResponseInput{}), nil
}

// VolumeEcShardsRebuild implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) VolumeEcShardsRebuild(req []byte) ([]byte, error) {
	view, err := volume_serverwire.WrapVolumeEcShardsRebuildRequest(req)
	if err != nil {
		return nil, err
	}
	resp, err := s.vs.VolumeEcShardsRebuild(context.Background(), &volume_server_pb.VolumeEcShardsRebuildRequest{
		VolumeId:            view.VolumeID(),
		Collection:          view.Collection(),
		UnsafeIgnoreSidecar: view.UnsafeIgnoreSidecar(),
	})
	if err != nil {
		return nil, err
	}
	return volume_serverwire.NewVolumeEcShardsRebuildResponse(volume_serverwire.VolumeEcShardsRebuildResponseInput{
		RebuiltShardIDs: resp.GetRebuiltShardIds(),
	}), nil
}

// VolumeEcShardsCopy implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) VolumeEcShardsCopy(req []byte) ([]byte, error) {
	view, err := volume_serverwire.WrapVolumeEcShardsCopyRequest(req)
	if err != nil {
		return nil, err
	}
	shardIDs := make([]uint32, view.ShardIDsLen())
	for i := range shardIDs {
		shardIDs[i] = view.ShardIDAt(i)
	}
	_, err = s.vs.VolumeEcShardsCopy(context.Background(), &volume_server_pb.VolumeEcShardsCopyRequest{
		VolumeId:       view.VolumeID(),
		Collection:     view.Collection(),
		ShardIds:       shardIDs,
		CopyEcxFile:    view.CopyEcxFile(),
		SourceDataNode: view.SourceDataNode(),
		CopyEcjFile:    view.CopyEcjFile(),
		CopyVifFile:    view.CopyVifFile(),
		DiskId:         view.DiskID(),
		CopyEcsumFile:  view.CopyEcsumFile(),
	})
	if err != nil {
		return nil, err
	}
	return volume_serverwire.NewVolumeEcShardsCopyResponse(volume_serverwire.VolumeEcShardsCopyResponseInput{}), nil
}

// VolumeEcShardsDelete implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) VolumeEcShardsDelete(req []byte) ([]byte, error) {
	view, err := volume_serverwire.WrapVolumeEcShardsDeleteRequest(req)
	if err != nil {
		return nil, err
	}
	shardIDs := make([]uint32, view.ShardIDsLen())
	for i := range shardIDs {
		shardIDs[i] = view.ShardIDAt(i)
	}
	resp, err := s.vs.VolumeEcShardsDelete(context.Background(), &volume_server_pb.VolumeEcShardsDeleteRequest{
		VolumeId:     view.VolumeID(),
		Collection:   view.Collection(),
		ShardIds:     shardIDs,
		FullTeardown: view.FullTeardown(),
		EncodeTsNs:   view.EncodeTsNs(),
	})
	if err != nil {
		return nil, err
	}
	return volume_serverwire.NewVolumeEcShardsDeleteResponse(volume_serverwire.VolumeEcShardsDeleteResponseInput{
		FullTeardownDone: resp.GetFullTeardownDone(),
	}), nil
}

// VolumeEcShardsMount implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) VolumeEcShardsMount(req []byte) ([]byte, error) {
	view, err := volume_serverwire.WrapVolumeEcShardsMountRequest(req)
	if err != nil {
		return nil, err
	}
	shardIDs := make([]uint32, view.ShardIDsLen())
	for i := range shardIDs {
		shardIDs[i] = view.ShardIDAt(i)
	}
	_, err = s.vs.VolumeEcShardsMount(context.Background(), &volume_server_pb.VolumeEcShardsMountRequest{
		VolumeId:       view.VolumeID(),
		Collection:     view.Collection(),
		ShardIds:       shardIDs,
		SourceDiskType: view.SourceDiskType(),
	})
	if err != nil {
		return nil, err
	}
	return volume_serverwire.NewVolumeEcShardsMountResponse(volume_serverwire.VolumeEcShardsMountResponseInput{}), nil
}

// VolumeEcShardsUnmount implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) VolumeEcShardsUnmount(req []byte) ([]byte, error) {
	view, err := volume_serverwire.WrapVolumeEcShardsUnmountRequest(req)
	if err != nil {
		return nil, err
	}
	shardIDs := make([]uint32, view.ShardIDsLen())
	for i := range shardIDs {
		shardIDs[i] = view.ShardIDAt(i)
	}
	_, err = s.vs.VolumeEcShardsUnmount(context.Background(), &volume_server_pb.VolumeEcShardsUnmountRequest{
		VolumeId:   view.VolumeID(),
		ShardIds:   shardIDs,
		EncodeTsNs: view.EncodeTsNs(),
	})
	if err != nil {
		return nil, err
	}
	return volume_serverwire.NewVolumeEcShardsUnmountResponse(volume_serverwire.VolumeEcShardsUnmountResponseInput{}), nil
}

// VolumeEcBlobDelete implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) VolumeEcBlobDelete(req []byte) ([]byte, error) {
	view, err := volume_serverwire.WrapVolumeEcBlobDeleteRequest(req)
	if err != nil {
		return nil, err
	}
	_, err = s.vs.VolumeEcBlobDelete(context.Background(), &volume_server_pb.VolumeEcBlobDeleteRequest{
		VolumeId:   view.VolumeID(),
		Collection: view.Collection(),
		FileKey:    view.FileKey(),
		Version:    view.Version(),
	})
	if err != nil {
		return nil, err
	}
	return volume_serverwire.NewVolumeEcBlobDeleteResponse(volume_serverwire.VolumeEcBlobDeleteResponseInput{}), nil
}

// VolumeEcShardsToVolume implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) VolumeEcShardsToVolume(req []byte) ([]byte, error) {
	view, err := volume_serverwire.WrapVolumeEcShardsToVolumeRequest(req)
	if err != nil {
		return nil, err
	}
	_, err = s.vs.VolumeEcShardsToVolume(context.Background(), &volume_server_pb.VolumeEcShardsToVolumeRequest{
		VolumeId:   view.VolumeID(),
		Collection: view.Collection(),
	})
	if err != nil {
		return nil, err
	}
	return volume_serverwire.NewVolumeEcShardsToVolumeResponse(volume_serverwire.VolumeEcShardsToVolumeResponseInput{}), nil
}

// VolumeEcShardsInfo implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) VolumeEcShardsInfo(req []byte) ([]byte, error) {
	view, err := volume_serverwire.WrapVolumeEcShardsInfoRequest(req)
	if err != nil {
		return nil, err
	}
	resp, err := s.vs.VolumeEcShardsInfo(context.Background(), &volume_server_pb.VolumeEcShardsInfoRequest{
		VolumeId: view.VolumeID(),
	})
	if err != nil {
		return nil, err
	}
	var ecShardInfos [][]byte
	for _, e := range resp.GetEcShardInfos() {
		ecShardInfos = append(ecShardInfos, ecShardInfoProtoToWire(e))
	}
	return volume_serverwire.NewVolumeEcShardsInfoResponse(volume_serverwire.VolumeEcShardsInfoResponseInput{
		EcShardInfos:     ecShardInfos,
		VolumeSize:       resp.GetVolumeSize(),
		FileCount:        resp.GetFileCount(),
		FileDeletedCount: resp.GetFileDeletedCount(),
	}), nil
}

// VolumeServerStatus implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) VolumeServerStatus(req []byte) ([]byte, error) {
	resp, err := s.vs.VolumeServerStatus(context.Background(), &volume_server_pb.VolumeServerStatusRequest{})
	if err != nil {
		return nil, err
	}
	var diskStatuses [][]byte
	for _, d := range resp.GetDiskStatuses() {
		diskStatuses = append(diskStatuses, diskStatusProtoToWire(d))
	}
	return volume_serverwire.NewVolumeServerStatusResponse(volume_serverwire.VolumeServerStatusResponseInput{
		DiskStatuses: diskStatuses,
		MemoryStatus: memStatusProtoToWire(resp.GetMemoryStatus()),
		Version:      resp.GetVersion(),
		DataCenter:   resp.GetDataCenter(),
		Rack:         resp.GetRack(),
		State:        stateProtoToWire(resp.GetState()),
	}), nil
}

// VolumeServerLeave implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) VolumeServerLeave(req []byte) ([]byte, error) {
	_, err := s.vs.VolumeServerLeave(context.Background(), &volume_server_pb.VolumeServerLeaveRequest{})
	if err != nil {
		return nil, err
	}
	return volume_serverwire.NewVolumeServerLeaveResponse(volume_serverwire.VolumeServerLeaveResponseInput{}), nil
}

// FetchAndWriteNeedle implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) FetchAndWriteNeedle(req []byte) ([]byte, error) {
	view, err := volume_serverwire.WrapFetchAndWriteNeedleRequest(req)
	if err != nil {
		return nil, err
	}
	replicas := make([]*volume_server_pb.FetchAndWriteNeedleRequest_Replica, view.ReplicasLen())
	for i := range replicas {
		r, ok := view.ReplicaAt(i)
		if !ok {
			continue
		}
		replicas[i] = &volume_server_pb.FetchAndWriteNeedleRequest_Replica{
			Url:       r.URL(),
			PublicUrl: r.PublicURL(),
			GrpcPort:  r.GrpcPort(),
		}
	}
	pbReq := &volume_server_pb.FetchAndWriteNeedleRequest{
		VolumeId:            view.VolumeID(),
		NeedleId:            view.NeedleID(),
		Cookie:              view.Cookie(),
		Offset:              view.Offset(),
		Size:                view.Size(),
		Replicas:            replicas,
		Auth:                view.Auth(),
		DownloadConcurrency: view.DownloadConcurrency(),
	}
	if c, ok := view.RemoteConf(); ok {
		pbReq.RemoteConf = remoteConfWireToProto(c)
	}
	if l, ok := view.RemoteLocation(); ok {
		pbReq.RemoteLocation = remoteLocationWireToProto(l)
	}
	resp, err := s.vs.FetchAndWriteNeedle(context.Background(), pbReq)
	if err != nil {
		return nil, err
	}
	return volume_serverwire.NewFetchAndWriteNeedleResponse(volume_serverwire.FetchAndWriteNeedleResponseInput{
		ETag: resp.GetETag(),
	}), nil
}

// ScrubVolume implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) ScrubVolume(req []byte) ([]byte, error) {
	view, err := volume_serverwire.WrapScrubVolumeRequest(req)
	if err != nil {
		return nil, err
	}
	volumeIDs := make([]uint32, view.VolumeIDsLen())
	for i := range volumeIDs {
		volumeIDs[i] = view.VolumeIDAt(i)
	}
	resp, err := s.vs.ScrubVolume(context.Background(), &volume_server_pb.ScrubVolumeRequest{
		Mode:                      volume_server_pb.VolumeScrubMode(view.Mode()),
		VolumeIds:                 volumeIDs,
		MarkBrokenVolumesReadonly: view.MarkBrokenVolumesReadonly(),
	})
	if err != nil {
		return nil, err
	}
	return volume_serverwire.NewScrubVolumeResponse(volume_serverwire.ScrubVolumeResponseInput{
		TotalVolumes:    resp.GetTotalVolumes(),
		TotalFiles:      resp.GetTotalFiles(),
		BrokenVolumeIDs: resp.GetBrokenVolumeIds(),
		Details:         resp.GetDetails(),
	}), nil
}

// ScrubEcVolume implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) ScrubEcVolume(req []byte) ([]byte, error) {
	view, err := volume_serverwire.WrapScrubEcVolumeRequest(req)
	if err != nil {
		return nil, err
	}
	volumeIDs := make([]uint32, view.VolumeIDsLen())
	for i := range volumeIDs {
		volumeIDs[i] = view.VolumeIDAt(i)
	}
	resp, err := s.vs.ScrubEcVolume(context.Background(), &volume_server_pb.ScrubEcVolumeRequest{
		Mode:      volume_server_pb.VolumeScrubMode(view.Mode()),
		VolumeIds: volumeIDs,
	})
	if err != nil {
		return nil, err
	}
	var brokenShardInfos [][]byte
	for _, e := range resp.GetBrokenShardInfos() {
		brokenShardInfos = append(brokenShardInfos, ecShardInfoProtoToWire(e))
	}
	return volume_serverwire.NewScrubEcVolumeResponse(volume_serverwire.ScrubEcVolumeResponseInput{
		TotalVolumes:     resp.GetTotalVolumes(),
		TotalFiles:       resp.GetTotalFiles(),
		BrokenVolumeIDs:  resp.GetBrokenVolumeIds(),
		BrokenShardInfos: brokenShardInfos,
		Details:          resp.GetDetails(),
	}), nil
}

// VolumeNeedleStatus implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) VolumeNeedleStatus(req []byte) ([]byte, error) {
	view, err := volume_serverwire.WrapVolumeNeedleStatusRequest(req)
	if err != nil {
		return nil, err
	}
	resp, err := s.vs.VolumeNeedleStatus(context.Background(), &volume_server_pb.VolumeNeedleStatusRequest{
		VolumeId: view.VolumeID(),
		NeedleId: view.NeedleID(),
	})
	if err != nil {
		return nil, err
	}
	return volume_serverwire.NewVolumeNeedleStatusResponse(volume_serverwire.VolumeNeedleStatusResponseInput{
		NeedleID:     resp.GetNeedleId(),
		Cookie:       resp.GetCookie(),
		Size:         resp.GetSize(),
		LastModified: resp.GetLastModified(),
		Crc:          resp.GetCrc(),
		TTL:          resp.GetTtl(),
	}), nil
}

// Ping implements volume_serverwire.VolumeServerStore.
func (s volumeServerStore) Ping(req []byte) ([]byte, error) {
	view, err := volume_serverwire.WrapPingRequest(req)
	if err != nil {
		return nil, err
	}
	resp, err := s.vs.Ping(context.Background(), &volume_server_pb.PingRequest{
		Target:     view.Target(),
		TargetType: view.TargetType(),
	})
	if err != nil {
		return nil, err
	}
	return volume_serverwire.NewPingResponse(volume_serverwire.PingResponseInput{
		StartTimeNs:  resp.GetStartTimeNs(),
		RemoteTimeNs: resp.GetRemoteTimeNs(),
		StopTimeNs:   resp.GetStopTimeNs(),
	}), nil
}

// volumeServerStore must satisfy the full UNARY surface of the wire contract.
// The compile-time assertion is the strangler's guard rail: it pins every one of
// the 37 unary VolumeServerStore methods to *volumeServerStore so any drift in
// the generated wire contract (a renamed/added unary RPC) breaks THIS build, not
// the running server. The 11 streaming methods are deliberately absent — they are
// bound separately at registration time (transport.ListenStream), so the full
// VolumeServerStore is composed of this store plus the stream handlers; asserting
// the unary subset here keeps the two halves independently checkable.
var _ unaryVolumeServerStore = volumeServerStore{}

// unaryVolumeServerStore is the unary-only projection of
// volume_serverwire.VolumeServerStore (the 11 streaming methods omitted).
type unaryVolumeServerStore interface {
	BatchDelete(req []byte) ([]byte, error)
	VacuumVolumeCheck(req []byte) ([]byte, error)
	VacuumVolumeCommit(req []byte) ([]byte, error)
	VacuumVolumeCleanup(req []byte) ([]byte, error)
	DeleteCollection(req []byte) ([]byte, error)
	AllocateVolume(req []byte) ([]byte, error)
	VolumeSyncStatus(req []byte) ([]byte, error)
	VolumeMount(req []byte) ([]byte, error)
	VolumeUnmount(req []byte) ([]byte, error)
	VolumeDelete(req []byte) ([]byte, error)
	VolumeMarkReadonly(req []byte) ([]byte, error)
	VolumeMarkWritable(req []byte) ([]byte, error)
	VolumeConfigure(req []byte) ([]byte, error)
	VolumeStatus(req []byte) ([]byte, error)
	GetState(req []byte) ([]byte, error)
	SetState(req []byte) ([]byte, error)
	ReadVolumeFileStatus(req []byte) ([]byte, error)
	ReadNeedleBlob(req []byte) ([]byte, error)
	ReadNeedleMeta(req []byte) ([]byte, error)
	WriteNeedleBlob(req []byte) ([]byte, error)
	VolumeTailReceiver(req []byte) ([]byte, error)
	VolumeEcShardsGenerate(req []byte) ([]byte, error)
	VolumeEcShardsRebuild(req []byte) ([]byte, error)
	VolumeEcShardsCopy(req []byte) ([]byte, error)
	VolumeEcShardsDelete(req []byte) ([]byte, error)
	VolumeEcShardsMount(req []byte) ([]byte, error)
	VolumeEcShardsUnmount(req []byte) ([]byte, error)
	VolumeEcBlobDelete(req []byte) ([]byte, error)
	VolumeEcShardsToVolume(req []byte) ([]byte, error)
	VolumeEcShardsInfo(req []byte) ([]byte, error)
	VolumeServerStatus(req []byte) ([]byte, error)
	VolumeServerLeave(req []byte) ([]byte, error)
	FetchAndWriteNeedle(req []byte) ([]byte, error)
	ScrubVolume(req []byte) ([]byte, error)
	ScrubEcVolume(req []byte) ([]byte, error)
	VolumeNeedleStatus(req []byte) ([]byte, error)
	Ping(req []byte) ([]byte, error)
}
