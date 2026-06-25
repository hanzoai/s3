// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// convert.go is the proto<->wire leaf converter set shared by the per-RPC
// request builders and response decoders in rpc.go. Each nested message gets a
// <Msg>ToWire(*pb.<Msg>) []byte builder and a <msg>FromView(wire.<Msg>) *pb.<Msg>
// decoder; the response decoders pull views via the volume_server wire
// convention of (View, bool) presence accessors. This is the ONLY place
// proto<->wire mapping for VolumeServer's nested message graph lives.

package volumezap

import (
	remote_pb "github.com/hanzoai/s3/s3/pb/remote_pb"
	volume_server_pb "github.com/hanzoai/s3/s3/pb/volume_server_pb"
	remotewire "github.com/hanzoai/s3/s3/wire/remote"
	vsw "github.com/hanzoai/s3/s3/wire/volume_server"
)

// --- VolumeServerState ---

func volumeServerStateToWire(s *volume_server_pb.VolumeServerState) []byte {
	if s == nil {
		return nil
	}
	return vsw.NewVolumeServerState(vsw.VolumeServerStateInput{
		Maintenance: s.Maintenance,
		Version:     s.Version,
	})
}

func volumeServerStateFromView(v vsw.VolumeServerState) *volume_server_pb.VolumeServerState {
	return &volume_server_pb.VolumeServerState{Maintenance: v.Maintenance(), Version: v.Version()}
}

// --- DeleteResult ---

func deleteResultFromView(v vsw.DeleteResult) *volume_server_pb.DeleteResult {
	return &volume_server_pb.DeleteResult{
		FileId:  v.FileID(),
		Status:  v.Status(),
		Error:   v.Error(),
		Size:    v.Size(),
		Version: v.Version(),
	}
}

// --- EcShardInfo ---

func ecShardInfoFromView(v vsw.EcShardInfo) *volume_server_pb.EcShardInfo {
	return &volume_server_pb.EcShardInfo{
		ShardId:    v.ShardID(),
		Size:       v.Size(),
		Collection: v.Collection(),
		VolumeId:   v.VolumeID(),
	}
}

// --- DiskStatus / MemStatus ---

func diskStatusFromView(v vsw.DiskStatus) *volume_server_pb.DiskStatus {
	return &volume_server_pb.DiskStatus{
		Dir:         v.Dir(),
		All:         v.All(),
		Used:        v.Used(),
		Free:        v.Free(),
		PercentFree: v.PercentFree(),
		PercentUsed: v.PercentUsed(),
		DiskType:    v.DiskType(),
		Error:       v.Error(),
	}
}

func memStatusFromView(v vsw.MemStatus) *volume_server_pb.MemStatus {
	return &volume_server_pb.MemStatus{
		Goroutines: v.Goroutines(),
		All:        v.All(),
		Used:       v.Used(),
		Free:       v.Free(),
		Self:       v.Self(),
		Heap:       v.Heap(),
		Stack:      v.Stack(),
	}
}

// --- RemoteFile / EcShardConfig / VolumeInfo ---

func remoteFileFromView(v vsw.RemoteFile) *volume_server_pb.RemoteFile {
	return &volume_server_pb.RemoteFile{
		BackendType:  v.BackendType(),
		BackendId:    v.BackendID(),
		Key:          v.Key(),
		Offset:       v.Offset(),
		FileSize:     v.FileSize(),
		ModifiedTime: v.ModifiedTime(),
		Extension:    v.Extension(),
	}
}

func ecShardConfigFromView(v vsw.EcShardConfig) *volume_server_pb.EcShardConfig {
	return &volume_server_pb.EcShardConfig{
		DataShards:   v.DataShards(),
		ParityShards: v.ParityShards(),
		EncodeTsNs:   v.EncodeTsNs(),
	}
}

func volumeInfoFromView(v vsw.VolumeInfo) *volume_server_pb.VolumeInfo {
	info := &volume_server_pb.VolumeInfo{
		Version:     v.Version(),
		Replication: v.Replication(),
		BytesOffset: v.BytesOffset(),
		DatFileSize: v.DatFileSize(),
		ExpireAtSec: v.ExpireAtSec(),
		ReadOnly:    v.ReadOnly(),
	}
	for i := 0; i < v.FilesLen(); i++ {
		if f, ok := v.FileAt(i); ok {
			info.Files = append(info.Files, remoteFileFromView(f))
		}
	}
	if cfg, ok := v.EcShardConfig(); ok {
		info.EcShardConfig = ecShardConfigFromView(cfg)
	}
	return info
}

// --- FetchAndWriteNeedle replica + remote storage messages ---

func fetchReplicaToWire(r *volume_server_pb.FetchAndWriteNeedleRequest_Replica) []byte {
	if r == nil {
		return nil
	}
	return vsw.NewFetchAndWriteNeedleReplica(vsw.FetchAndWriteNeedleReplicaInput{
		URL: r.Url, PublicURL: r.PublicUrl, GrpcPort: r.GrpcPort,
	})
}

func remoteConfToWire(c *remote_pb.RemoteConf) []byte {
	if c == nil {
		return nil
	}
	return remotewire.NewRemoteConf(remotewire.RemoteConfInput{
		Type:                            c.Type,
		Name:                            c.Name,
		S3AccessKey:                     c.S3AccessKey,
		S3SecretKey:                     c.S3SecretKey,
		S3Region:                        c.S3Region,
		S3Endpoint:                      c.S3Endpoint,
		S3StorageClass:                  c.S3StorageClass,
		S3ForcePathStyle:                c.S3ForcePathStyle,
		S3SupportTagging:                c.S3SupportTagging,
		S3V4Signature:                   c.S3V4Signature,
		GcsGoogleApplicationCredentials: c.GcsGoogleApplicationCredentials,
		GcsProjectID:                    c.GcsProjectId,
		AzureAccountName:                c.AzureAccountName,
		AzureAccountKey:                 c.AzureAccountKey,
		BackblazeKeyID:                  c.BackblazeKeyId,
		BackblazeApplicationKey:         c.BackblazeApplicationKey,
		BackblazeEndpoint:               c.BackblazeEndpoint,
		BackblazeRegion:                 c.BackblazeRegion,
		AliyunAccessKey:                 c.AliyunAccessKey,
		AliyunSecretKey:                 c.AliyunSecretKey,
		AliyunEndpoint:                  c.AliyunEndpoint,
		AliyunRegion:                    c.AliyunRegion,
		TencentSecretID:                 c.TencentSecretId,
		TencentSecretKey:                c.TencentSecretKey,
		TencentEndpoint:                 c.TencentEndpoint,
		BaiduAccessKey:                  c.BaiduAccessKey,
		BaiduSecretKey:                  c.BaiduSecretKey,
		BaiduEndpoint:                   c.BaiduEndpoint,
		BaiduRegion:                     c.BaiduRegion,
		WasabiAccessKey:                 c.WasabiAccessKey,
		WasabiSecretKey:                 c.WasabiSecretKey,
		WasabiEndpoint:                  c.WasabiEndpoint,
		WasabiRegion:                    c.WasabiRegion,
		FilebaseAccessKey:               c.FilebaseAccessKey,
		FilebaseSecretKey:               c.FilebaseSecretKey,
		FilebaseEndpoint:                c.FilebaseEndpoint,
		StorjAccessKey:                  c.StorjAccessKey,
		StorjSecretKey:                  c.StorjSecretKey,
		StorjEndpoint:                   c.StorjEndpoint,
		ContaboAccessKey:                c.ContaboAccessKey,
		ContaboSecretKey:                c.ContaboSecretKey,
		ContaboEndpoint:                 c.ContaboEndpoint,
		ContaboRegion:                   c.ContaboRegion,
	})
}

func remoteLocationToWire(l *remote_pb.RemoteStorageLocation) []byte {
	if l == nil {
		return nil
	}
	return remotewire.NewRemoteStorageLocation(remotewire.RemoteStorageLocationInput{
		Name: l.Name, Bucket: l.Bucket, Path: l.Path, ListingCacheTTLSeconds: l.ListingCacheTtlSeconds,
	})
}
