// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// server_rpc.go completes the filer ZAP server backend: the remaining unary
// HanzoFiler RPCs beyond entry CRUD (server.go). Each is the server-side mirror
// of the client converter in rpc.go — read the request wire view into proto,
// call the existing filer_pb.HanzoFilerServer, build the response wire from
// proto. Together with server.go this implements every unary method of the
// filerwire.Backend, so filerwire.Serve(NewServerBackend(fs)) serves the whole
// filer over ZAP with no gRPC.

package filerzap

import (
	filer_pb "github.com/hanzoai/s3/s3/pb/filer_pb"
	filerwire "github.com/hanzoai/s3/s3/wire/filer"
)

func (b serverBackend) TouchAccessTime(v filerwire.TouchAccessTimeRequest) ([]byte, error) {
	resp, err := b.fs.TouchAccessTime(b.ctx, &filer_pb.TouchAccessTimeRequest{
		Directory: v.Directory(), Name: v.Name(), ClientAtimeNs: v.ClientAtimeNs(),
	})
	if err != nil {
		return nil, err
	}
	return filerwire.NewTouchAccessTimeResponse(filerwire.TouchAccessTimeResponseInput{
		PersistedAtimeNs: resp.PersistedAtimeNs, Updated: resp.Updated,
	}), nil
}

func (b serverBackend) AppendToEntry(v filerwire.AppendToEntryRequest) ([]byte, error) {
	_, err := b.fs.AppendToEntry(b.ctx, &filer_pb.AppendToEntryRequest{
		Directory: v.Directory(), EntryName: v.EntryName(),
		Chunks: chunksFromView(v.ChunksLen(), v.Chunks),
	})
	if err != nil {
		return nil, err
	}
	return filerwire.NewAppendToEntryResponse(filerwire.AppendToEntryResponseInput{}), nil
}

func (b serverBackend) ObjectTransaction(v filerwire.ObjectTransactionRequest) ([]byte, error) {
	resp, err := b.fs.ObjectTransaction(b.ctx, objectTransactionReqFromView(v))
	if err != nil {
		return nil, err
	}
	return filerwire.NewObjectTransactionResponse(filerwire.ObjectTransactionResponseInput{
		Error: resp.Error, ErrorCode: uint32(resp.ErrorCode),
	}), nil
}

func (b serverBackend) ObjectTransactionBatch(v filerwire.ObjectTransactionBatchRequest) ([]byte, error) {
	req := &filer_pb.ObjectTransactionBatchRequest{}
	for i := 0; i < v.TransactionsLen(); i++ {
		req.Transactions = append(req.Transactions, objectTransactionReqFromView(v.Transactions(i)))
	}
	resp, err := b.fs.ObjectTransactionBatch(b.ctx, req)
	if err != nil {
		return nil, err
	}
	outs := make([][]byte, len(resp.Responses))
	for i, r := range resp.Responses {
		outs[i] = filerwire.NewObjectTransactionResponse(filerwire.ObjectTransactionResponseInput{
			Error: r.Error, ErrorCode: uint32(r.ErrorCode),
		})
	}
	return filerwire.NewObjectTransactionBatchResponse(filerwire.ObjectTransactionBatchResponseInput{Responses: outs}), nil
}

func (b serverBackend) PosixLock(v filerwire.PosixLockRequest) ([]byte, error) {
	req := &filer_pb.PosixLockRequest{
		Key: v.Key(), IsMoved: v.IsMoved(), Op: filer_pb.PosixLockOp(v.Op()), CoolingProbe: v.CoolingProbe(),
	}
	if v.HasLock() {
		req.Lock = posixLockRangeFromView(v.Lock())
	}
	for i := 0; i < v.LocksLen(); i++ {
		req.Locks = append(req.Locks, posixLockRangeFromView(v.Locks(i)))
	}
	resp, err := b.fs.PosixLock(b.ctx, req)
	if err != nil {
		return nil, err
	}
	in := filerwire.PosixLockResponseInput{Granted: resp.Granted, HasConflict: resp.HasConflict}
	if resp.Conflict != nil {
		in.Conflict = posixLockRangeToWire(resp.Conflict)
	}
	return filerwire.NewPosixLockResponse(in), nil
}

func (b serverBackend) AtomicRenameEntry(v filerwire.AtomicRenameEntryRequest) ([]byte, error) {
	_, err := b.fs.AtomicRenameEntry(b.ctx, &filer_pb.AtomicRenameEntryRequest{
		OldDirectory: v.OldDirectory(), OldName: v.OldName(),
		NewDirectory: v.NewDirectory(), NewName: v.NewName(),
		Signatures: signaturesFromView(v.SignaturesLen(), v.Signature),
	})
	if err != nil {
		return nil, err
	}
	return filerwire.NewAtomicRenameEntryResponse(filerwire.AtomicRenameEntryResponseInput{}), nil
}

func (b serverBackend) AssignVolume(v filerwire.AssignVolumeRequest) ([]byte, error) {
	resp, err := b.fs.AssignVolume(b.ctx, &filer_pb.AssignVolumeRequest{
		Count: v.Count(), Collection: v.Collection(), Replication: v.Replication(),
		TtlSec: v.TtlSec(), DataCenter: v.DataCenter(), Path: v.Path(), Rack: v.Rack(),
		DataNode: v.DataNode(), DiskType: v.DiskType(), ExpectedDataSize: v.ExpectedDataSize(),
	})
	if err != nil {
		return nil, err
	}
	in := filerwire.AssignVolumeResponseInput{
		FileID: resp.FileId, Count: resp.Count, Auth: resp.Auth,
		Collection: resp.Collection, Replication: resp.Replication, Error: resp.Error,
	}
	if resp.Location != nil {
		in.Location = locationToWire(resp.Location)
	}
	return filerwire.NewAssignVolumeResponse(in), nil
}

func (b serverBackend) LookupVolume(v filerwire.LookupVolumeRequest) ([]byte, error) {
	ids := make([]string, v.VolumeIDsLen())
	for i := range ids {
		ids[i] = v.VolumeID(i)
	}
	resp, err := b.fs.LookupVolume(b.ctx, &filer_pb.LookupVolumeRequest{VolumeIds: ids})
	if err != nil {
		return nil, err
	}
	entries := make([][]byte, 0, len(resp.LocationsMap))
	for k, locs := range resp.LocationsMap {
		locBufs := make([][]byte, len(locs.Locations))
		for i, l := range locs.Locations {
			locBufs[i] = locationToWire(l)
		}
		entries = append(entries, filerwire.NewLookupVolumeResponseLocationsMapEntry(filerwire.LookupVolumeResponseLocationsMapEntryInput{
			Key: k, Value: filerwire.NewLocations(filerwire.LocationsInput{Locations: locBufs}),
		}))
	}
	return filerwire.NewLookupVolumeResponse(filerwire.LookupVolumeResponseInput{LocationsMap: entries}), nil
}

func (b serverBackend) CollectionList(v filerwire.CollectionListRequest) ([]byte, error) {
	resp, err := b.fs.CollectionList(b.ctx, &filer_pb.CollectionListRequest{
		IncludeNormalVolumes: v.IncludeNormalVolumes(), IncludeEcVolumes: v.IncludeEcVolumes(),
	})
	if err != nil {
		return nil, err
	}
	cols := make([][]byte, len(resp.Collections))
	for i, c := range resp.Collections {
		cols[i] = filerwire.NewCollection(filerwire.CollectionInput{Name: c.Name})
	}
	return filerwire.NewCollectionListResponse(filerwire.CollectionListResponseInput{Collections: cols}), nil
}

func (b serverBackend) DeleteCollection(v filerwire.DeleteCollectionRequest) ([]byte, error) {
	_, err := b.fs.DeleteCollection(b.ctx, &filer_pb.DeleteCollectionRequest{Collection: v.Collection()})
	if err != nil {
		return nil, err
	}
	return filerwire.NewDeleteCollectionResponse(filerwire.DeleteCollectionResponseInput{}), nil
}

func (b serverBackend) Statistics(v filerwire.StatisticsRequest) ([]byte, error) {
	resp, err := b.fs.Statistics(b.ctx, &filer_pb.StatisticsRequest{
		Replication: v.Replication(), Collection: v.Collection(), Ttl: v.Ttl(), DiskType: v.DiskType(),
	})
	if err != nil {
		return nil, err
	}
	return filerwire.NewStatisticsResponse(filerwire.StatisticsResponseInput{
		TotalSize: resp.TotalSize, UsedSize: resp.UsedSize, FileCount: resp.FileCount,
	}), nil
}

func (b serverBackend) Ping(v filerwire.PingRequest) ([]byte, error) {
	resp, err := b.fs.Ping(b.ctx, &filer_pb.PingRequest{Target: v.Target(), TargetType: v.TargetType()})
	if err != nil {
		return nil, err
	}
	return filerwire.NewPingResponse(filerwire.PingResponseInput{
		StartTimeNs: resp.StartTimeNs, RemoteTimeNs: resp.RemoteTimeNs, StopTimeNs: resp.StopTimeNs,
	}), nil
}

func (b serverBackend) GetFilerConfiguration(v filerwire.GetFilerConfigurationRequest) ([]byte, error) {
	resp, err := b.fs.GetFilerConfiguration(b.ctx, &filer_pb.GetFilerConfigurationRequest{})
	if err != nil {
		return nil, err
	}
	return filerwire.NewGetFilerConfigurationResponse(filerwire.GetFilerConfigurationResponseInput{
		Masters: resp.Masters, Replication: resp.Replication, Collection: resp.Collection,
		MaxMb: resp.MaxMb, DirBuckets: resp.DirBuckets, Cipher: resp.Cipher, Signature: resp.Signature,
		MetricsAddress: resp.MetricsAddress, MetricsIntervalSec: resp.MetricsIntervalSec,
		Version: resp.Version, ClusterID: resp.ClusterId, FilerGroup: resp.FilerGroup,
		MajorVersion: resp.MajorVersion, MinorVersion: resp.MinorVersion,
	}), nil
}

func (b serverBackend) ListMetadataSubscribers(v filerwire.ListMetadataSubscribersRequest) ([]byte, error) {
	types := make([]string, v.ClientTypesLen())
	for i := range types {
		types[i] = v.ClientType(i)
	}
	resp, err := b.fs.ListMetadataSubscribers(b.ctx, &filer_pb.ListMetadataSubscribersRequest{ClientTypes: types})
	if err != nil {
		return nil, err
	}
	subs := make([][]byte, len(resp.Subscribers))
	for i, s := range resp.Subscribers {
		subs[i] = filerwire.NewMetadataSubscriber(filerwire.MetadataSubscriberInput{
			ClientName: s.ClientName, ClientType: s.ClientType, Address: s.Address,
			PathPrefix: s.PathPrefix, ClientID: s.ClientId, ClientEpoch: s.ClientEpoch,
			ConnectedAtNs: s.ConnectedAtNs, FilerAddress: s.FilerAddress,
		})
	}
	return filerwire.NewListMetadataSubscribersResponse(filerwire.ListMetadataSubscribersResponseInput{Subscribers: subs}), nil
}

func (b serverBackend) KvGet(v filerwire.KvGetRequest) ([]byte, error) {
	resp, err := b.fs.KvGet(b.ctx, &filer_pb.KvGetRequest{Key: v.Key()})
	if err != nil {
		return nil, err
	}
	return filerwire.NewKvGetResponse(filerwire.KvGetResponseInput{Value: resp.Value, Error: resp.Error}), nil
}

func (b serverBackend) KvPut(v filerwire.KvPutRequest) ([]byte, error) {
	resp, err := b.fs.KvPut(b.ctx, &filer_pb.KvPutRequest{Key: v.Key(), Value: v.Value()})
	if err != nil {
		return nil, err
	}
	return filerwire.NewKvPutResponse(filerwire.KvPutResponseInput{Error: resp.Error}), nil
}

func (b serverBackend) DistributedLock(v filerwire.LockRequest) ([]byte, error) {
	resp, err := b.fs.DistributedLock(b.ctx, &filer_pb.LockRequest{
		Name: v.Name(), SecondsToLock: v.SecondsToLock(), RenewToken: v.RenewToken(),
		IsMoved: v.IsMoved(), Owner: v.Owner(),
	})
	if err != nil {
		return nil, err
	}
	return filerwire.NewLockResponse(filerwire.LockResponseInput{
		RenewToken: resp.RenewToken, LockOwner: resp.LockOwner, LockHostMovedTo: resp.LockHostMovedTo,
		Error: resp.Error, Generation: resp.Generation,
	}), nil
}

func (b serverBackend) DistributedUnlock(v filerwire.UnlockRequest) ([]byte, error) {
	resp, err := b.fs.DistributedUnlock(b.ctx, &filer_pb.UnlockRequest{
		Name: v.Name(), RenewToken: v.RenewToken(), IsMoved: v.IsMoved(),
	})
	if err != nil {
		return nil, err
	}
	return filerwire.NewUnlockResponse(filerwire.UnlockResponseInput{Error: resp.Error, MovedTo: resp.MovedTo}), nil
}

func (b serverBackend) FindLockOwner(v filerwire.FindLockOwnerRequest) ([]byte, error) {
	resp, err := b.fs.FindLockOwner(b.ctx, &filer_pb.FindLockOwnerRequest{Name: v.Name(), IsMoved: v.IsMoved()})
	if err != nil {
		return nil, err
	}
	return filerwire.NewFindLockOwnerResponse(filerwire.FindLockOwnerResponseInput{Owner: resp.Owner}), nil
}

func (b serverBackend) TransferLocks(v filerwire.TransferLocksRequest) ([]byte, error) {
	req := &filer_pb.TransferLocksRequest{}
	for i := 0; i < v.LocksLen(); i++ {
		l := v.Locks(i)
		req.Locks = append(req.Locks, &filer_pb.Lock{
			Name: l.Name(), RenewToken: l.RenewToken(), ExpiredAtNs: l.ExpiredAtNs(),
			Owner: l.Owner(), Generation: l.Generation(), IsBackup: l.IsBackup(), Seq: l.Seq(),
		})
	}
	if _, err := b.fs.TransferLocks(b.ctx, req); err != nil {
		return nil, err
	}
	return filerwire.NewTransferLocksResponse(filerwire.TransferLocksResponseInput{}), nil
}

func (b serverBackend) ReplicateLock(v filerwire.ReplicateLockRequest) ([]byte, error) {
	_, err := b.fs.ReplicateLock(b.ctx, &filer_pb.ReplicateLockRequest{
		Name: v.Name(), RenewToken: v.RenewToken(), ExpiredAtNs: v.ExpiredAtNs(),
		Owner: v.Owner(), Generation: v.Generation(), IsUnlock: v.IsUnlock(), Seq: v.Seq(),
	})
	if err != nil {
		return nil, err
	}
	return filerwire.NewReplicateLockResponse(filerwire.ReplicateLockResponseInput{}), nil
}

func (b serverBackend) MountRegister(v filerwire.MountRegisterRequest) ([]byte, error) {
	_, err := b.fs.MountRegister(b.ctx, &filer_pb.MountRegisterRequest{
		PeerAddr: v.PeerAddr(), Rack: v.Rack(), TtlSeconds: v.TtlSeconds(), DataCenter: v.DataCenter(),
	})
	if err != nil {
		return nil, err
	}
	return filerwire.NewMountRegisterResponse(filerwire.MountRegisterResponseInput{}), nil
}

func (b serverBackend) MountList(v filerwire.MountListRequest) ([]byte, error) {
	resp, err := b.fs.MountList(b.ctx, &filer_pb.MountListRequest{})
	if err != nil {
		return nil, err
	}
	mounts := make([][]byte, len(resp.Mounts))
	for i, m := range resp.Mounts {
		mounts[i] = filerwire.NewMountInfo(filerwire.MountInfoInput{
			PeerAddr: m.PeerAddr, Rack: m.Rack, LastSeenNs: m.LastSeenNs, DataCenter: m.DataCenter,
		})
	}
	return filerwire.NewMountListResponse(filerwire.MountListResponseInput{Mounts: mounts}), nil
}

// --- request view -> proto converters (shared by server.go and StreamMutateEntry) ---

func createEntryReqFromView(v filerwire.CreateEntryRequest) *filer_pb.CreateEntryRequest {
	req := &filer_pb.CreateEntryRequest{
		Directory:                v.Directory(),
		OExcl:                    v.OExcl(),
		IsFromOtherCluster:       v.IsFromOtherCluster(),
		SkipCheckParentDirectory: v.SkipCheckParentDirectory(),
		Signatures:               signaturesFromView(v.SignaturesLen(), v.Signature),
		Condition:                writeConditionFromView(v.HasCondition(), v.Condition()),
	}
	if v.HasEntry() {
		req.Entry = entryFromView(v.Entry())
	}
	return req
}

func updateEntryReqFromView(v filerwire.UpdateEntryRequest) *filer_pb.UpdateEntryRequest {
	req := &filer_pb.UpdateEntryRequest{
		Directory:          v.Directory(),
		IsFromOtherCluster: v.IsFromOtherCluster(),
		Signatures:         signaturesFromView(v.SignaturesLen(), v.Signature),
		ExpectedExtended:   extendedFromView(v.ExpectedExtendedLen(), v.ExpectedExtended),
	}
	if v.HasEntry() {
		req.Entry = entryFromView(v.Entry())
	}
	return req
}

func deleteEntryReqFromView(v filerwire.DeleteEntryRequest) *filer_pb.DeleteEntryRequest {
	return &filer_pb.DeleteEntryRequest{
		Directory:            v.Directory(),
		Name:                 v.Name(),
		IsDeleteData:         v.IsDeleteData(),
		IsRecursive:          v.IsRecursive(),
		IgnoreRecursiveError: v.IgnoreRecursiveError(),
		IsFromOtherCluster:   v.IsFromOtherCluster(),
		IfNotModifiedAfter:   v.IfNotModifiedAfter(),
		Signatures:           signaturesFromView(v.SignaturesLen(), v.Signature),
	}
}

func streamRenameEntryReqFromView(v filerwire.StreamRenameEntryRequest) *filer_pb.StreamRenameEntryRequest {
	return &filer_pb.StreamRenameEntryRequest{
		OldDirectory: v.OldDirectory(), OldName: v.OldName(),
		NewDirectory: v.NewDirectory(), NewName: v.NewName(),
		Signatures: signaturesFromView(v.SignaturesLen(), v.Signature),
	}
}

// --- server-side leaf converters ---

func locationToWire(l *filer_pb.Location) []byte {
	return filerwire.NewLocation(filerwire.LocationInput{
		URL: l.Url, PublicURL: l.PublicUrl, GrpcPort: l.GrpcPort, DataCenter: l.DataCenter,
	})
}

func objectTransactionReqFromView(v filerwire.ObjectTransactionRequest) *filer_pb.ObjectTransactionRequest {
	req := &filer_pb.ObjectTransactionRequest{
		LockKey: v.LockKey(), IsFromOtherCluster: v.IsFromOtherCluster(),
		ConditionKey: v.ConditionKey(), RouteKey: v.RouteKey(), IsMoved: v.IsMoved(),
		Signatures: signaturesFromView(v.SignaturesLen(), v.Signature),
	}
	if v.HasCondition() {
		req.Condition = writeConditionFromView(true, v.Condition())
	}
	for i := 0; i < v.MutationsLen(); i++ {
		req.Mutations = append(req.Mutations, objectMutationFromView(v.Mutations(i)))
	}
	return req
}

func objectMutationFromView(m filerwire.ObjectMutation) *filer_pb.ObjectMutation {
	out := &filer_pb.ObjectMutation{
		Type: filer_pb.ObjectMutation_Type(m.Type()), Directory: m.Directory(), Name: m.Name(),
		IsDeleteData: m.IsDeleteData(), IsRecursive: m.IsRecursive(),
		SetContent: m.SetContent(), Content: m.Content(), TouchMtime: m.TouchMtime(),
	}
	for i := 0; i < m.DeleteExtendedLen(); i++ {
		out.DeleteExtended = append(out.DeleteExtended, m.DeleteExtended(i))
	}
	if m.HasEntry() {
		out.Entry = entryFromView(m.Entry())
	}
	if m.HasRecompute() {
		out.Recompute = recomputeFromView(m.Recompute())
	}
	if n := m.SetExtendedLen(); n > 0 {
		out.SetExtended = make(map[string][]byte, n)
		for i := 0; i < n; i++ {
			kv := m.SetExtended(i)
			out.SetExtended[kv.Key()] = kv.Value()
		}
	}
	return out
}

func recomputeFromView(v filerwire.Recompute) *filer_pb.Recompute {
	out := &filer_pb.Recompute{
		ScanDir: v.ScanDir(), Descending: v.Descending(),
		NameToKey: v.NameToKey(), SizeToKey: v.SizeToKey(), MtimeToKey: v.MtimeToKey(),
		DemoteKey: v.DemoteKey(), DemoteValue: v.DemoteValue(), ExcludeName: v.ExcludeName(),
	}
	if n := v.CopyExtendedLen(); n > 0 {
		out.CopyExtended = make(map[string]string, n)
		for i := 0; i < n; i++ {
			kv := v.CopyExtended(i)
			out.CopyExtended[kv.Key()] = string(kv.Value())
		}
	}
	return out
}
