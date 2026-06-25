// Package filerstub provides embeddable, do-nothing implementations of the
// generated server interfaces (filer_pb.HanzoFilerServer, master_pb.HanzoServer)
// for tests. They replace the gRPC-generated Unimplemented*Server bases the gRPC
// rip deleted: a test double embeds the stub to inherit the full method set and
// overrides only the methods it exercises. Every stub method returns
// errUnimplemented (unary: nil result + error; streaming: error), so an
// unexpected call surfaces loudly instead of compiling against a missing method.
package filerstub

import (
	context "context"
	"errors"

	"github.com/hanzoai/s3/s3/pb/filer_pb"
	"github.com/hanzoai/s3/s3/pb/master_pb"
)

var errUnimplemented = errors.New("filerstub: method not implemented")

// FilerServer is an embeddable filer_pb.HanzoFilerServer whose every method is
// unimplemented. Embed it in a test double and override what the test needs.
type FilerServer struct{}

var _ filer_pb.HanzoFilerServer = FilerServer{}

func (s FilerServer) LookupDirectoryEntry(context.Context, *filer_pb.LookupDirectoryEntryRequest) (*filer_pb.LookupDirectoryEntryResponse, error) {
	return nil, errUnimplemented
}
func (s FilerServer) ListEntries(*filer_pb.ListEntriesRequest, filer_pb.HanzoFiler_ListEntriesServer) error {
	return errUnimplemented
}
func (s FilerServer) CreateEntry(context.Context, *filer_pb.CreateEntryRequest) (*filer_pb.CreateEntryResponse, error) {
	return nil, errUnimplemented
}
func (s FilerServer) UpdateEntry(context.Context, *filer_pb.UpdateEntryRequest) (*filer_pb.UpdateEntryResponse, error) {
	return nil, errUnimplemented
}
func (s FilerServer) TouchAccessTime(context.Context, *filer_pb.TouchAccessTimeRequest) (*filer_pb.TouchAccessTimeResponse, error) {
	return nil, errUnimplemented
}
func (s FilerServer) AppendToEntry(context.Context, *filer_pb.AppendToEntryRequest) (*filer_pb.AppendToEntryResponse, error) {
	return nil, errUnimplemented
}
func (s FilerServer) DeleteEntry(context.Context, *filer_pb.DeleteEntryRequest) (*filer_pb.DeleteEntryResponse, error) {
	return nil, errUnimplemented
}
func (s FilerServer) ObjectTransaction(context.Context, *filer_pb.ObjectTransactionRequest) (*filer_pb.ObjectTransactionResponse, error) {
	return nil, errUnimplemented
}
func (s FilerServer) ObjectTransactionBatch(context.Context, *filer_pb.ObjectTransactionBatchRequest) (*filer_pb.ObjectTransactionBatchResponse, error) {
	return nil, errUnimplemented
}
func (s FilerServer) PosixLock(context.Context, *filer_pb.PosixLockRequest) (*filer_pb.PosixLockResponse, error) {
	return nil, errUnimplemented
}
func (s FilerServer) AtomicRenameEntry(context.Context, *filer_pb.AtomicRenameEntryRequest) (*filer_pb.AtomicRenameEntryResponse, error) {
	return nil, errUnimplemented
}
func (s FilerServer) StreamRenameEntry(*filer_pb.StreamRenameEntryRequest, filer_pb.HanzoFiler_StreamRenameEntryServer) error {
	return errUnimplemented
}
func (s FilerServer) StreamMutateEntry(filer_pb.HanzoFiler_StreamMutateEntryServer) error {
	return errUnimplemented
}
func (s FilerServer) AssignVolume(context.Context, *filer_pb.AssignVolumeRequest) (*filer_pb.AssignVolumeResponse, error) {
	return nil, errUnimplemented
}
func (s FilerServer) LookupVolume(context.Context, *filer_pb.LookupVolumeRequest) (*filer_pb.LookupVolumeResponse, error) {
	return nil, errUnimplemented
}
func (s FilerServer) CollectionList(context.Context, *filer_pb.CollectionListRequest) (*filer_pb.CollectionListResponse, error) {
	return nil, errUnimplemented
}
func (s FilerServer) DeleteCollection(context.Context, *filer_pb.DeleteCollectionRequest) (*filer_pb.DeleteCollectionResponse, error) {
	return nil, errUnimplemented
}
func (s FilerServer) Statistics(context.Context, *filer_pb.StatisticsRequest) (*filer_pb.StatisticsResponse, error) {
	return nil, errUnimplemented
}
func (s FilerServer) Ping(context.Context, *filer_pb.PingRequest) (*filer_pb.PingResponse, error) {
	return nil, errUnimplemented
}
func (s FilerServer) GetFilerConfiguration(context.Context, *filer_pb.GetFilerConfigurationRequest) (*filer_pb.GetFilerConfigurationResponse, error) {
	return nil, errUnimplemented
}
func (s FilerServer) TraverseBfsMetadata(*filer_pb.TraverseBfsMetadataRequest, filer_pb.HanzoFiler_TraverseBfsMetadataServer) error {
	return errUnimplemented
}
func (s FilerServer) SubscribeMetadata(*filer_pb.SubscribeMetadataRequest, filer_pb.HanzoFiler_SubscribeMetadataServer) error {
	return errUnimplemented
}
func (s FilerServer) SubscribeLocalMetadata(*filer_pb.SubscribeMetadataRequest, filer_pb.HanzoFiler_SubscribeLocalMetadataServer) error {
	return errUnimplemented
}
func (s FilerServer) ListMetadataSubscribers(context.Context, *filer_pb.ListMetadataSubscribersRequest) (*filer_pb.ListMetadataSubscribersResponse, error) {
	return nil, errUnimplemented
}
func (s FilerServer) KvGet(context.Context, *filer_pb.KvGetRequest) (*filer_pb.KvGetResponse, error) {
	return nil, errUnimplemented
}
func (s FilerServer) KvPut(context.Context, *filer_pb.KvPutRequest) (*filer_pb.KvPutResponse, error) {
	return nil, errUnimplemented
}
func (s FilerServer) CacheRemoteObjectToLocalCluster(context.Context, *filer_pb.CacheRemoteObjectToLocalClusterRequest) (*filer_pb.CacheRemoteObjectToLocalClusterResponse, error) {
	return nil, errUnimplemented
}
func (s FilerServer) DistributedLock(context.Context, *filer_pb.LockRequest) (*filer_pb.LockResponse, error) {
	return nil, errUnimplemented
}
func (s FilerServer) DistributedUnlock(context.Context, *filer_pb.UnlockRequest) (*filer_pb.UnlockResponse, error) {
	return nil, errUnimplemented
}
func (s FilerServer) FindLockOwner(context.Context, *filer_pb.FindLockOwnerRequest) (*filer_pb.FindLockOwnerResponse, error) {
	return nil, errUnimplemented
}
func (s FilerServer) TransferLocks(context.Context, *filer_pb.TransferLocksRequest) (*filer_pb.TransferLocksResponse, error) {
	return nil, errUnimplemented
}
func (s FilerServer) ReplicateLock(context.Context, *filer_pb.ReplicateLockRequest) (*filer_pb.ReplicateLockResponse, error) {
	return nil, errUnimplemented
}
func (s FilerServer) MountRegister(context.Context, *filer_pb.MountRegisterRequest) (*filer_pb.MountRegisterResponse, error) {
	return nil, errUnimplemented
}
func (s FilerServer) MountList(context.Context, *filer_pb.MountListRequest) (*filer_pb.MountListResponse, error) {
	return nil, errUnimplemented
}

// MasterServer is an embeddable master_pb.HanzoServer whose every method is
// unimplemented. Embed it in a test double and override what the test needs.
type MasterServer struct{}

var _ master_pb.HanzoServer = MasterServer{}

func (s MasterServer) SendHeartbeat(master_pb.Hanzo_SendHeartbeatServer) error {
	return errUnimplemented
}
func (s MasterServer) KeepConnected(master_pb.Hanzo_KeepConnectedServer) error {
	return errUnimplemented
}
func (s MasterServer) LookupVolume(context.Context, *master_pb.LookupVolumeRequest) (*master_pb.LookupVolumeResponse, error) {
	return nil, errUnimplemented
}
func (s MasterServer) Assign(context.Context, *master_pb.AssignRequest) (*master_pb.AssignResponse, error) {
	return nil, errUnimplemented
}
func (s MasterServer) StreamAssign(master_pb.Hanzo_StreamAssignServer) error { return errUnimplemented }
func (s MasterServer) Statistics(context.Context, *master_pb.StatisticsRequest) (*master_pb.StatisticsResponse, error) {
	return nil, errUnimplemented
}
func (s MasterServer) CollectionList(context.Context, *master_pb.CollectionListRequest) (*master_pb.CollectionListResponse, error) {
	return nil, errUnimplemented
}
func (s MasterServer) CollectionDelete(context.Context, *master_pb.CollectionDeleteRequest) (*master_pb.CollectionDeleteResponse, error) {
	return nil, errUnimplemented
}
func (s MasterServer) VolumeList(context.Context, *master_pb.VolumeListRequest) (*master_pb.VolumeListResponse, error) {
	return nil, errUnimplemented
}
func (s MasterServer) LookupEcVolume(context.Context, *master_pb.LookupEcVolumeRequest) (*master_pb.LookupEcVolumeResponse, error) {
	return nil, errUnimplemented
}
func (s MasterServer) VacuumVolume(context.Context, *master_pb.VacuumVolumeRequest) (*master_pb.VacuumVolumeResponse, error) {
	return nil, errUnimplemented
}
func (s MasterServer) DisableVacuum(context.Context, *master_pb.DisableVacuumRequest) (*master_pb.DisableVacuumResponse, error) {
	return nil, errUnimplemented
}
func (s MasterServer) EnableVacuum(context.Context, *master_pb.EnableVacuumRequest) (*master_pb.EnableVacuumResponse, error) {
	return nil, errUnimplemented
}
func (s MasterServer) VolumeMarkReadonly(context.Context, *master_pb.VolumeMarkReadonlyRequest) (*master_pb.VolumeMarkReadonlyResponse, error) {
	return nil, errUnimplemented
}
func (s MasterServer) GetMasterConfiguration(context.Context, *master_pb.GetMasterConfigurationRequest) (*master_pb.GetMasterConfigurationResponse, error) {
	return nil, errUnimplemented
}
func (s MasterServer) ListClusterNodes(context.Context, *master_pb.ListClusterNodesRequest) (*master_pb.ListClusterNodesResponse, error) {
	return nil, errUnimplemented
}
func (s MasterServer) LeaseAdminToken(context.Context, *master_pb.LeaseAdminTokenRequest) (*master_pb.LeaseAdminTokenResponse, error) {
	return nil, errUnimplemented
}
func (s MasterServer) ReleaseAdminToken(context.Context, *master_pb.ReleaseAdminTokenRequest) (*master_pb.ReleaseAdminTokenResponse, error) {
	return nil, errUnimplemented
}
func (s MasterServer) Ping(context.Context, *master_pb.PingRequest) (*master_pb.PingResponse, error) {
	return nil, errUnimplemented
}
func (s MasterServer) RaftListClusterServers(context.Context, *master_pb.RaftListClusterServersRequest) (*master_pb.RaftListClusterServersResponse, error) {
	return nil, errUnimplemented
}
func (s MasterServer) RaftAddServer(context.Context, *master_pb.RaftAddServerRequest) (*master_pb.RaftAddServerResponse, error) {
	return nil, errUnimplemented
}
func (s MasterServer) RaftRemoveServer(context.Context, *master_pb.RaftRemoveServerRequest) (*master_pb.RaftRemoveServerResponse, error) {
	return nil, errUnimplemented
}
func (s MasterServer) RaftLeadershipTransfer(context.Context, *master_pb.RaftLeadershipTransferRequest) (*master_pb.RaftLeadershipTransferResponse, error) {
	return nil, errUnimplemented
}
func (s MasterServer) VolumeGrow(context.Context, *master_pb.VolumeGrowRequest) (*master_pb.VolumeGrowResponse, error) {
	return nil, errUnimplemented
}
