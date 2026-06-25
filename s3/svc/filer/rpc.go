// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// rpc.go holds the per-RPC request builders and response decoders for the 27 unary
// HanzoFiler RPCs, built on the leaf converters in convert.go and entry.go. Naming:
// <Rpc>ReqToWire(*pb.<Rpc>Request) []byte and <Rpc>RespFromWire([]byte) (*pb.<Rpc>Response, error).

package filer

import (
	filer_pb "github.com/hanzoai/s3/s3/pb/filer_pb"
	filerwire "github.com/hanzoai/s3/s3/wire/filer"
)

// LookupDirectoryEntry
func LookupDirectoryEntryReqToWire(r *filer_pb.LookupDirectoryEntryRequest) []byte {
	return filerwire.NewLookupDirectoryEntryRequest(filerwire.LookupDirectoryEntryRequestInput{
		Directory: r.Directory, Name: r.Name,
	})
}
func LookupDirectoryEntryRespFromWire(b []byte) (*filer_pb.LookupDirectoryEntryResponse, error) {
	v, err := filerwire.WrapLookupDirectoryEntryResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &filer_pb.LookupDirectoryEntryResponse{}
	if v.HasEntry() {
		resp.Entry = entryFromView(v.Entry())
	}
	return resp, nil
}

// CreateEntry
func CreateEntryReqToWire(r *filer_pb.CreateEntryRequest) []byte {
	return filerwire.NewCreateEntryRequest(filerwire.CreateEntryRequestInput{
		Directory:                r.Directory,
		Entry:                    EntryToWire(r.Entry),
		OExcl:                    r.OExcl,
		IsFromOtherCluster:       r.IsFromOtherCluster,
		Signatures:               r.Signatures,
		SkipCheckParentDirectory: r.SkipCheckParentDirectory,
		Condition:                writeConditionToWire(r.Condition),
	})
}
func CreateEntryRespFromWire(b []byte) (*filer_pb.CreateEntryResponse, error) {
	v, err := filerwire.WrapCreateEntryResponse(b)
	if err != nil {
		return nil, err
	}
	return &filer_pb.CreateEntryResponse{
		Error:         v.Error(),
		ErrorCode:     filer_pb.FilerError(v.ErrorCode()),
		MetadataEvent: metadataEventFromView(v.HasMetadataEvent(), v.MetadataEvent()),
	}, nil
}

// UpdateEntry
func UpdateEntryReqToWire(r *filer_pb.UpdateEntryRequest) []byte {
	return filerwire.NewUpdateEntryRequest(filerwire.UpdateEntryRequestInput{
		Directory:          r.Directory,
		Entry:              EntryToWire(r.Entry),
		IsFromOtherCluster: r.IsFromOtherCluster,
		Signatures:         r.Signatures,
		ExpectedExtended:   extendedToWire(r.ExpectedExtended),
	})
}
func UpdateEntryRespFromWire(b []byte) (*filer_pb.UpdateEntryResponse, error) {
	v, err := filerwire.WrapUpdateEntryResponse(b)
	if err != nil {
		return nil, err
	}
	return &filer_pb.UpdateEntryResponse{
		MetadataEvent: metadataEventFromView(v.HasMetadataEvent(), v.MetadataEvent()),
	}, nil
}

// TouchAccessTime
func TouchAccessTimeReqToWire(r *filer_pb.TouchAccessTimeRequest) []byte {
	return filerwire.NewTouchAccessTimeRequest(filerwire.TouchAccessTimeRequestInput{
		Directory: r.Directory, Name: r.Name, ClientAtimeNs: r.ClientAtimeNs,
	})
}
func TouchAccessTimeRespFromWire(b []byte) (*filer_pb.TouchAccessTimeResponse, error) {
	v, err := filerwire.WrapTouchAccessTimeResponse(b)
	if err != nil {
		return nil, err
	}
	return &filer_pb.TouchAccessTimeResponse{PersistedAtimeNs: v.PersistedAtimeNs(), Updated: v.Updated()}, nil
}

// AppendToEntry
func AppendToEntryReqToWire(r *filer_pb.AppendToEntryRequest) []byte {
	return filerwire.NewAppendToEntryRequest(filerwire.AppendToEntryRequestInput{
		Directory: r.Directory, EntryName: r.EntryName, Chunks: chunksToWire(r.Chunks),
	})
}
func AppendToEntryRespFromWire(b []byte) (*filer_pb.AppendToEntryResponse, error) {
	if _, err := filerwire.WrapAppendToEntryResponse(b); err != nil {
		return nil, err
	}
	return &filer_pb.AppendToEntryResponse{}, nil
}

// DeleteEntry
func DeleteEntryReqToWire(r *filer_pb.DeleteEntryRequest) []byte {
	return filerwire.NewDeleteEntryRequest(filerwire.DeleteEntryRequestInput{
		Directory:            r.Directory,
		Name:                 r.Name,
		IsDeleteData:         r.IsDeleteData,
		IsRecursive:          r.IsRecursive,
		IgnoreRecursiveError: r.IgnoreRecursiveError,
		IsFromOtherCluster:   r.IsFromOtherCluster,
		Signatures:           r.Signatures,
		IfNotModifiedAfter:   r.IfNotModifiedAfter,
	})
}
func DeleteEntryRespFromWire(b []byte) (*filer_pb.DeleteEntryResponse, error) {
	v, err := filerwire.WrapDeleteEntryResponse(b)
	if err != nil {
		return nil, err
	}
	return &filer_pb.DeleteEntryResponse{
		Error:         v.Error(),
		MetadataEvent: metadataEventFromView(v.HasMetadataEvent(), v.MetadataEvent()),
	}, nil
}

// ObjectTransaction
func objectTransactionReqInput(r *filer_pb.ObjectTransactionRequest) filerwire.ObjectTransactionRequestInput {
	muts := make([][]byte, len(r.Mutations))
	for i, m := range r.Mutations {
		muts[i] = objectMutationToWire(m)
	}
	return filerwire.ObjectTransactionRequestInput{
		LockKey:            r.LockKey,
		Condition:          writeConditionToWire(r.Condition),
		Mutations:          muts,
		IsFromOtherCluster: r.IsFromOtherCluster,
		Signatures:         r.Signatures,
		ConditionKey:       r.ConditionKey,
		RouteKey:           r.RouteKey,
		IsMoved:            r.IsMoved,
	}
}
func ObjectTransactionReqToWire(r *filer_pb.ObjectTransactionRequest) []byte {
	return filerwire.NewObjectTransactionRequest(objectTransactionReqInput(r))
}
func objectTransactionRespFromView(v filerwire.ObjectTransactionResponse) *filer_pb.ObjectTransactionResponse {
	return &filer_pb.ObjectTransactionResponse{Error: v.Error(), ErrorCode: filer_pb.FilerError(v.ErrorCode())}
}
func ObjectTransactionRespFromWire(b []byte) (*filer_pb.ObjectTransactionResponse, error) {
	v, err := filerwire.WrapObjectTransactionResponse(b)
	if err != nil {
		return nil, err
	}
	return objectTransactionRespFromView(v), nil
}

// ObjectTransactionBatch
func ObjectTransactionBatchReqToWire(r *filer_pb.ObjectTransactionBatchRequest) []byte {
	txns := make([][]byte, len(r.Transactions))
	for i, t := range r.Transactions {
		txns[i] = ObjectTransactionReqToWire(t)
	}
	return filerwire.NewObjectTransactionBatchRequest(filerwire.ObjectTransactionBatchRequestInput{Transactions: txns})
}
func ObjectTransactionBatchRespFromWire(b []byte) (*filer_pb.ObjectTransactionBatchResponse, error) {
	v, err := filerwire.WrapObjectTransactionBatchResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &filer_pb.ObjectTransactionBatchResponse{}
	if n := v.ResponsesLen(); n > 0 {
		resp.Responses = make([]*filer_pb.ObjectTransactionResponse, n)
		for i := 0; i < n; i++ {
			resp.Responses[i] = objectTransactionRespFromView(v.Responses(i))
		}
	}
	return resp, nil
}

// PosixLock
func PosixLockReqToWire(r *filer_pb.PosixLockRequest) []byte {
	locks := make([][]byte, len(r.Locks))
	for i, l := range r.Locks {
		locks[i] = posixLockRangeToWire(l)
	}
	return filerwire.NewPosixLockRequest(filerwire.PosixLockRequestInput{
		Key:          r.Key,
		IsMoved:      r.IsMoved,
		Op:           uint32(r.Op),
		Lock:         posixLockRangeToWire(r.Lock),
		Locks:        locks,
		CoolingProbe: r.CoolingProbe,
	})
}
func PosixLockRespFromWire(b []byte) (*filer_pb.PosixLockResponse, error) {
	v, err := filerwire.WrapPosixLockResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &filer_pb.PosixLockResponse{Granted: v.Granted(), HasConflict: v.HasConflict()}
	if v.HasConflictRange() {
		resp.Conflict = posixLockRangeFromView(v.Conflict())
	}
	return resp, nil
}

// AtomicRenameEntry
func AtomicRenameEntryReqToWire(r *filer_pb.AtomicRenameEntryRequest) []byte {
	return filerwire.NewAtomicRenameEntryRequest(filerwire.AtomicRenameEntryRequestInput{
		OldDirectory: r.OldDirectory, OldName: r.OldName,
		NewDirectory: r.NewDirectory, NewName: r.NewName, Signatures: r.Signatures,
	})
}
func AtomicRenameEntryRespFromWire(b []byte) (*filer_pb.AtomicRenameEntryResponse, error) {
	if _, err := filerwire.WrapAtomicRenameEntryResponse(b); err != nil {
		return nil, err
	}
	return &filer_pb.AtomicRenameEntryResponse{}, nil
}

// AssignVolume
func AssignVolumeReqToWire(r *filer_pb.AssignVolumeRequest) []byte {
	return filerwire.NewAssignVolumeRequest(filerwire.AssignVolumeRequestInput{
		Count: r.Count, Collection: r.Collection, Replication: r.Replication,
		TtlSec: r.TtlSec, DataCenter: r.DataCenter, Path: r.Path, Rack: r.Rack,
		DataNode: r.DataNode, DiskType: r.DiskType, ExpectedDataSize: r.ExpectedDataSize,
	})
}
func AssignVolumeRespFromWire(b []byte) (*filer_pb.AssignVolumeResponse, error) {
	v, err := filerwire.WrapAssignVolumeResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &filer_pb.AssignVolumeResponse{
		FileId: v.FileID(), Count: v.Count(), Auth: v.Auth(),
		Collection: v.Collection(), Replication: v.Replication(), Error: v.Error(),
	}
	if v.HasLocation() {
		resp.Location = locationFromView(v.Location())
	}
	return resp, nil
}

func locationFromView(v filerwire.Location) *filer_pb.Location {
	return &filer_pb.Location{Url: v.URL(), PublicUrl: v.PublicURL(), GrpcPort: v.GrpcPort(), DataCenter: v.DataCenter()}
}

// LookupVolume
func LookupVolumeReqToWire(r *filer_pb.LookupVolumeRequest) []byte {
	return filerwire.NewLookupVolumeRequest(filerwire.LookupVolumeRequestInput{VolumeIDs: r.VolumeIds})
}
func LookupVolumeRespFromWire(b []byte) (*filer_pb.LookupVolumeResponse, error) {
	v, err := filerwire.WrapLookupVolumeResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &filer_pb.LookupVolumeResponse{}
	if n := v.LocationsMapLen(); n > 0 {
		resp.LocationsMap = make(map[string]*filer_pb.Locations, n)
		for i := 0; i < n; i++ {
			kv := v.LocationsMap(i)
			locs := kv.Value()
			out := &filer_pb.Locations{}
			for j := 0; j < locs.LocationsLen(); j++ {
				out.Locations = append(out.Locations, locationFromView(locs.Locations(j)))
			}
			resp.LocationsMap[kv.Key()] = out
		}
	}
	return resp, nil
}

// CollectionList
func CollectionListReqToWire(r *filer_pb.CollectionListRequest) []byte {
	return filerwire.NewCollectionListRequest(filerwire.CollectionListRequestInput{
		IncludeNormalVolumes: r.IncludeNormalVolumes, IncludeEcVolumes: r.IncludeEcVolumes,
	})
}
func CollectionListRespFromWire(b []byte) (*filer_pb.CollectionListResponse, error) {
	v, err := filerwire.WrapCollectionListResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &filer_pb.CollectionListResponse{}
	for i := 0; i < v.CollectionsLen(); i++ {
		resp.Collections = append(resp.Collections, &filer_pb.Collection{Name: v.Collections(i).Name()})
	}
	return resp, nil
}

// DeleteCollection
func DeleteCollectionReqToWire(r *filer_pb.DeleteCollectionRequest) []byte {
	return filerwire.NewDeleteCollectionRequest(filerwire.DeleteCollectionRequestInput{Collection: r.Collection})
}
func DeleteCollectionRespFromWire(b []byte) (*filer_pb.DeleteCollectionResponse, error) {
	if _, err := filerwire.WrapDeleteCollectionResponse(b); err != nil {
		return nil, err
	}
	return &filer_pb.DeleteCollectionResponse{}, nil
}

// Statistics
func StatisticsReqToWire(r *filer_pb.StatisticsRequest) []byte {
	return filerwire.NewStatisticsRequest(filerwire.StatisticsRequestInput{
		Replication: r.Replication, Collection: r.Collection, Ttl: r.Ttl, DiskType: r.DiskType,
	})
}
func StatisticsRespFromWire(b []byte) (*filer_pb.StatisticsResponse, error) {
	v, err := filerwire.WrapStatisticsResponse(b)
	if err != nil {
		return nil, err
	}
	return &filer_pb.StatisticsResponse{TotalSize: v.TotalSize(), UsedSize: v.UsedSize(), FileCount: v.FileCount()}, nil
}

// Ping
func PingReqToWire(r *filer_pb.PingRequest) []byte {
	return filerwire.NewPingRequest(filerwire.PingRequestInput{Target: r.Target, TargetType: r.TargetType})
}
func PingRespFromWire(b []byte) (*filer_pb.PingResponse, error) {
	v, err := filerwire.WrapPingResponse(b)
	if err != nil {
		return nil, err
	}
	return &filer_pb.PingResponse{StartTimeNs: v.StartTimeNs(), RemoteTimeNs: v.RemoteTimeNs(), StopTimeNs: v.StopTimeNs()}, nil
}

// GetFilerConfiguration
func GetFilerConfigurationReqToWire(r *filer_pb.GetFilerConfigurationRequest) []byte {
	return filerwire.NewGetFilerConfigurationRequest(filerwire.GetFilerConfigurationRequestInput{})
}
func GetFilerConfigurationRespFromWire(b []byte) (*filer_pb.GetFilerConfigurationResponse, error) {
	v, err := filerwire.WrapGetFilerConfigurationResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &filer_pb.GetFilerConfigurationResponse{
		Replication:        v.Replication(),
		Collection:         v.Collection(),
		MaxMb:              v.MaxMb(),
		DirBuckets:         v.DirBuckets(),
		Cipher:             v.Cipher(),
		Signature:          v.Signature(),
		MetricsAddress:     v.MetricsAddress(),
		MetricsIntervalSec: v.MetricsIntervalSec(),
		Version:            v.Version(),
		ClusterId:          v.ClusterID(),
		FilerGroup:         v.FilerGroup(),
		MajorVersion:       v.MajorVersion(),
		MinorVersion:       v.MinorVersion(),
	}
	for i := 0; i < v.MastersLen(); i++ {
		resp.Masters = append(resp.Masters, v.Master(i))
	}
	return resp, nil
}

// ListMetadataSubscribers
func ListMetadataSubscribersReqToWire(r *filer_pb.ListMetadataSubscribersRequest) []byte {
	return filerwire.NewListMetadataSubscribersRequest(filerwire.ListMetadataSubscribersRequestInput{ClientTypes: r.ClientTypes})
}
func ListMetadataSubscribersRespFromWire(b []byte) (*filer_pb.ListMetadataSubscribersResponse, error) {
	v, err := filerwire.WrapListMetadataSubscribersResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &filer_pb.ListMetadataSubscribersResponse{}
	for i := 0; i < v.SubscribersLen(); i++ {
		s := v.Subscribers(i)
		resp.Subscribers = append(resp.Subscribers, &filer_pb.MetadataSubscriber{
			ClientName: s.ClientName(), ClientType: s.ClientType(), Address: s.Address(),
			PathPrefix: s.PathPrefix(), ClientId: s.ClientID(), ClientEpoch: s.ClientEpoch(),
			ConnectedAtNs: s.ConnectedAtNs(), FilerAddress: s.FilerAddress(),
		})
	}
	return resp, nil
}

// KvGet / KvPut
func KvGetReqToWire(r *filer_pb.KvGetRequest) []byte {
	return filerwire.NewKvGetRequest(filerwire.KvGetRequestInput{Key: r.Key})
}
func KvGetRespFromWire(b []byte) (*filer_pb.KvGetResponse, error) {
	v, err := filerwire.WrapKvGetResponse(b)
	if err != nil {
		return nil, err
	}
	return &filer_pb.KvGetResponse{Value: v.Value(), Error: v.Error()}, nil
}
func KvPutReqToWire(r *filer_pb.KvPutRequest) []byte {
	return filerwire.NewKvPutRequest(filerwire.KvPutRequestInput{Key: r.Key, Value: r.Value})
}
func KvPutRespFromWire(b []byte) (*filer_pb.KvPutResponse, error) {
	v, err := filerwire.WrapKvPutResponse(b)
	if err != nil {
		return nil, err
	}
	return &filer_pb.KvPutResponse{Error: v.Error()}, nil
}

// CacheRemoteObjectToLocalCluster
func CacheRemoteObjectToLocalClusterReqToWire(r *filer_pb.CacheRemoteObjectToLocalClusterRequest) []byte {
	return filerwire.NewCacheRemoteObjectToLocalClusterRequest(filerwire.CacheRemoteObjectToLocalClusterRequestInput{
		Directory: r.Directory, Name: r.Name,
		ChunkConcurrency: r.ChunkConcurrency, DownloadConcurrency: r.DownloadConcurrency,
	})
}
func CacheRemoteObjectToLocalClusterRespFromWire(b []byte) (*filer_pb.CacheRemoteObjectToLocalClusterResponse, error) {
	v, err := filerwire.WrapCacheRemoteObjectToLocalClusterResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &filer_pb.CacheRemoteObjectToLocalClusterResponse{
		MetadataEvent: metadataEventFromView(v.HasMetadataEvent(), v.MetadataEvent()),
	}
	if v.HasEntry() {
		resp.Entry = entryFromView(v.Entry())
	}
	return resp, nil
}

// DistributedLock
func DistributedLockReqToWire(r *filer_pb.LockRequest) []byte {
	return filerwire.NewLockRequest(filerwire.LockRequestInput{
		Name: r.Name, SecondsToLock: r.SecondsToLock, RenewToken: r.RenewToken, IsMoved: r.IsMoved, Owner: r.Owner,
	})
}
func DistributedLockRespFromWire(b []byte) (*filer_pb.LockResponse, error) {
	v, err := filerwire.WrapLockResponse(b)
	if err != nil {
		return nil, err
	}
	return &filer_pb.LockResponse{
		RenewToken: v.RenewToken(), LockOwner: v.LockOwner(), LockHostMovedTo: v.LockHostMovedTo(),
		Error: v.Error(), Generation: v.Generation(),
	}, nil
}

// DistributedUnlock
func DistributedUnlockReqToWire(r *filer_pb.UnlockRequest) []byte {
	return filerwire.NewUnlockRequest(filerwire.UnlockRequestInput{Name: r.Name, RenewToken: r.RenewToken, IsMoved: r.IsMoved})
}
func DistributedUnlockRespFromWire(b []byte) (*filer_pb.UnlockResponse, error) {
	v, err := filerwire.WrapUnlockResponse(b)
	if err != nil {
		return nil, err
	}
	return &filer_pb.UnlockResponse{Error: v.Error(), MovedTo: v.MovedTo()}, nil
}

// FindLockOwner
func FindLockOwnerReqToWire(r *filer_pb.FindLockOwnerRequest) []byte {
	return filerwire.NewFindLockOwnerRequest(filerwire.FindLockOwnerRequestInput{Name: r.Name, IsMoved: r.IsMoved})
}
func FindLockOwnerRespFromWire(b []byte) (*filer_pb.FindLockOwnerResponse, error) {
	v, err := filerwire.WrapFindLockOwnerResponse(b)
	if err != nil {
		return nil, err
	}
	return &filer_pb.FindLockOwnerResponse{Owner: v.Owner()}, nil
}

// TransferLocks
func TransferLocksReqToWire(r *filer_pb.TransferLocksRequest) []byte {
	locks := make([][]byte, len(r.Locks))
	for i, l := range r.Locks {
		locks[i] = filerwire.NewLock(filerwire.LockInput{
			Name: l.Name, RenewToken: l.RenewToken, ExpiredAtNs: l.ExpiredAtNs,
			Owner: l.Owner, Generation: l.Generation, IsBackup: l.IsBackup, Seq: l.Seq,
		})
	}
	return filerwire.NewTransferLocksRequest(filerwire.TransferLocksRequestInput{Locks: locks})
}
func TransferLocksRespFromWire(b []byte) (*filer_pb.TransferLocksResponse, error) {
	if _, err := filerwire.WrapTransferLocksResponse(b); err != nil {
		return nil, err
	}
	return &filer_pb.TransferLocksResponse{}, nil
}

// ReplicateLock
func ReplicateLockReqToWire(r *filer_pb.ReplicateLockRequest) []byte {
	return filerwire.NewReplicateLockRequest(filerwire.ReplicateLockRequestInput{
		Name: r.Name, RenewToken: r.RenewToken, ExpiredAtNs: r.ExpiredAtNs,
		Owner: r.Owner, Generation: r.Generation, IsUnlock: r.IsUnlock, Seq: r.Seq,
	})
}
func ReplicateLockRespFromWire(b []byte) (*filer_pb.ReplicateLockResponse, error) {
	if _, err := filerwire.WrapReplicateLockResponse(b); err != nil {
		return nil, err
	}
	return &filer_pb.ReplicateLockResponse{}, nil
}

// MountRegister
func MountRegisterReqToWire(r *filer_pb.MountRegisterRequest) []byte {
	return filerwire.NewMountRegisterRequest(filerwire.MountRegisterRequestInput{
		PeerAddr: r.PeerAddr, Rack: r.Rack, TtlSeconds: r.TtlSeconds, DataCenter: r.DataCenter,
	})
}
func MountRegisterRespFromWire(b []byte) (*filer_pb.MountRegisterResponse, error) {
	if _, err := filerwire.WrapMountRegisterResponse(b); err != nil {
		return nil, err
	}
	return &filer_pb.MountRegisterResponse{}, nil
}

// MountList
func MountListReqToWire(r *filer_pb.MountListRequest) []byte {
	return filerwire.NewMountListRequest(filerwire.MountListRequestInput{})
}
func MountListRespFromWire(b []byte) (*filer_pb.MountListResponse, error) {
	v, err := filerwire.WrapMountListResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &filer_pb.MountListResponse{}
	for i := 0; i < v.MountsLen(); i++ {
		m := v.Mounts(i)
		resp.Mounts = append(resp.Mounts, &filer_pb.MountInfo{
			PeerAddr: m.PeerAddr(), Rack: m.Rack(), LastSeenNs: m.LastSeenNs(), DataCenter: m.DataCenter(),
		})
	}
	return resp, nil
}
