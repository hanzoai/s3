// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// server_convert.go is the server-side mirror of rpc.go/convert.go: where the
// client converters map *pb.<Req> -> wire bytes (request) and wire bytes ->
// *pb.<Resp> (response), the server converters do the inverse — wire request
// bytes -> *pb.<Req> (so the engine method can run) and *pb.<Resp> -> wire
// response bytes (so the reply ships zero-copy). The per-nested-message leaf
// converters here are the exact inverse of convert.go's leaves: <Msg>ToWire for
// response leaves, <msg>FromView for request leaves. Same wire Input fields,
// same view accessors as the client direction — so the mapping stays
// byte-faithful both ways.

package volumezap

import (
	remote_pb "github.com/hanzoai/s3/s3/pb/remote_pb"
	volume_server_pb "github.com/hanzoai/s3/s3/pb/volume_server_pb"
	remotewire "github.com/hanzoai/s3/s3/wire/remote"
	vsw "github.com/hanzoai/s3/s3/wire/volume_server"
)

// --- response leaf encoders (inverse of convert.go's *FromView) ---

func deleteResultToWire(r *volume_server_pb.DeleteResult) []byte {
	if r == nil {
		return nil
	}
	return vsw.NewDeleteResult(vsw.DeleteResultInput{
		FileID:  r.FileId,
		Status:  r.Status,
		Error:   r.Error,
		Size:    r.Size,
		Version: r.Version,
	})
}

func ecShardInfoToWire(r *volume_server_pb.EcShardInfo) []byte {
	if r == nil {
		return nil
	}
	return vsw.NewEcShardInfo(vsw.EcShardInfoInput{
		ShardID:    r.ShardId,
		Size:       r.Size,
		Collection: r.Collection,
		VolumeID:   r.VolumeId,
	})
}

func diskStatusToWire(r *volume_server_pb.DiskStatus) []byte {
	if r == nil {
		return nil
	}
	return vsw.NewDiskStatus(vsw.DiskStatusInput{
		Dir:         r.Dir,
		All:         r.All,
		Used:        r.Used,
		Free:        r.Free,
		PercentFree: r.PercentFree,
		PercentUsed: r.PercentUsed,
		DiskType:    r.DiskType,
		Error:       r.Error,
	})
}

func memStatusToWire(r *volume_server_pb.MemStatus) []byte {
	if r == nil {
		return nil
	}
	return vsw.NewMemStatus(vsw.MemStatusInput{
		Goroutines: r.Goroutines,
		All:        r.All,
		Used:       r.Used,
		Free:       r.Free,
		Self:       r.Self,
		Heap:       r.Heap,
		Stack:      r.Stack,
	})
}

func remoteFileToWire(r *volume_server_pb.RemoteFile) []byte {
	if r == nil {
		return nil
	}
	return vsw.NewRemoteFile(vsw.RemoteFileInput{
		BackendType:  r.BackendType,
		BackendID:    r.BackendId,
		Key:          r.Key,
		Offset:       r.Offset,
		FileSize:     r.FileSize,
		ModifiedTime: r.ModifiedTime,
		Extension:    r.Extension,
	})
}

func ecShardConfigToWire(r *volume_server_pb.EcShardConfig) []byte {
	if r == nil {
		return nil
	}
	return vsw.NewEcShardConfig(vsw.EcShardConfigInput{
		DataShards:   r.DataShards,
		ParityShards: r.ParityShards,
		EncodeTsNs:   r.EncodeTsNs,
	})
}

func volumeInfoToWire(r *volume_server_pb.VolumeInfo) []byte {
	if r == nil {
		return nil
	}
	files := make([][]byte, len(r.Files))
	for i, f := range r.Files {
		files[i] = remoteFileToWire(f)
	}
	return vsw.NewVolumeInfo(vsw.VolumeInfoInput{
		Files:         files,
		Version:       r.Version,
		Replication:   r.Replication,
		BytesOffset:   r.BytesOffset,
		DatFileSize:   r.DatFileSize,
		ExpireAtSec:   r.ExpireAtSec,
		ReadOnly:      r.ReadOnly,
		EcShardConfig: ecShardConfigToWire(r.EcShardConfig),
	})
}

// --- request leaf decoders (inverse of convert.go's *ToWire) ---

func fetchReplicaFromView(v vsw.FetchAndWriteNeedleReplica) *volume_server_pb.FetchAndWriteNeedleRequest_Replica {
	return &volume_server_pb.FetchAndWriteNeedleRequest_Replica{
		Url:       v.URL(),
		PublicUrl: v.PublicURL(),
		GrpcPort:  v.GrpcPort(),
	}
}

func remoteConfFromView(v remotewire.RemoteConf) *remote_pb.RemoteConf {
	return &remote_pb.RemoteConf{
		Type:                            v.Type(),
		Name:                            v.Name(),
		S3AccessKey:                     v.S3AccessKey(),
		S3SecretKey:                     v.S3SecretKey(),
		S3Region:                        v.S3Region(),
		S3Endpoint:                      v.S3Endpoint(),
		S3StorageClass:                  v.S3StorageClass(),
		S3ForcePathStyle:                v.S3ForcePathStyle(),
		S3SupportTagging:                v.S3SupportTagging(),
		S3V4Signature:                   v.S3V4Signature(),
		GcsGoogleApplicationCredentials: v.GcsGoogleApplicationCredentials(),
		GcsProjectId:                    v.GcsProjectID(),
		AzureAccountName:                v.AzureAccountName(),
		AzureAccountKey:                 v.AzureAccountKey(),
		BackblazeKeyId:                  v.BackblazeKeyID(),
		BackblazeApplicationKey:         v.BackblazeApplicationKey(),
		BackblazeEndpoint:               v.BackblazeEndpoint(),
		BackblazeRegion:                 v.BackblazeRegion(),
		AliyunAccessKey:                 v.AliyunAccessKey(),
		AliyunSecretKey:                 v.AliyunSecretKey(),
		AliyunEndpoint:                  v.AliyunEndpoint(),
		AliyunRegion:                    v.AliyunRegion(),
		TencentSecretId:                 v.TencentSecretID(),
		TencentSecretKey:                v.TencentSecretKey(),
		TencentEndpoint:                 v.TencentEndpoint(),
		BaiduAccessKey:                  v.BaiduAccessKey(),
		BaiduSecretKey:                  v.BaiduSecretKey(),
		BaiduEndpoint:                   v.BaiduEndpoint(),
		BaiduRegion:                     v.BaiduRegion(),
		WasabiAccessKey:                 v.WasabiAccessKey(),
		WasabiSecretKey:                 v.WasabiSecretKey(),
		WasabiEndpoint:                  v.WasabiEndpoint(),
		WasabiRegion:                    v.WasabiRegion(),
		FilebaseAccessKey:               v.FilebaseAccessKey(),
		FilebaseSecretKey:               v.FilebaseSecretKey(),
		FilebaseEndpoint:                v.FilebaseEndpoint(),
		StorjAccessKey:                  v.StorjAccessKey(),
		StorjSecretKey:                  v.StorjSecretKey(),
		StorjEndpoint:                   v.StorjEndpoint(),
		ContaboAccessKey:                v.ContaboAccessKey(),
		ContaboSecretKey:                v.ContaboSecretKey(),
		ContaboEndpoint:                 v.ContaboEndpoint(),
		ContaboRegion:                   v.ContaboRegion(),
	}
}

func remoteLocationFromView(v remotewire.RemoteStorageLocation) *remote_pb.RemoteStorageLocation {
	return &remote_pb.RemoteStorageLocation{
		Name:                   v.Name(),
		Bucket:                 v.Bucket(),
		Path:                   v.Path(),
		ListingCacheTtlSeconds: v.ListingCacheTTLSeconds(),
	}
}
