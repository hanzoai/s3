// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Package filerclient is the replication subsystem's strangler seam onto the
// native ZAP transport for the HanzoFiler service. It implements the
// filer_pb.HanzoFilerClient interface — the contract the replication source,
// filer sink and signature reader already speak — but routes every call over a
// github.com/zap-proto/go transport.Conn instead of gRPC:
//
//   - the 28 unary RPCs go through filerwire.HanzoFilerClient (over the conn's
//     Call channel), translating *filer_pb.<Rpc>Request <-> ZAP buffer <->
//     *filer_pb.<Rpc>Response with the filerzap converters;
//   - the 6 streaming RPCs open a transport stream via the filerstream client
//     and return an rpc.ServerStream / rpc.BidiStream whose Recv()/Send()
//     decode/encode each frame with the same filerzap converters.
//
// Because the adapter satisfies filer_pb.HanzoFilerClient, no replication call
// site changes: callers keep building *filer_pb requests and reading *filer_pb
// responses; the bytes are ZAP at every hop. This is the ONE hand-written client
// glue for replication — DRY across source.FilerSource.WithFilerClient,
// filersink.FilerSink.WithFilerClient and replicator.ReadFilerSignature.
//
// It mirrors the FUSE mount's filer_client_adapter.go (the reference cut) and
// reuses the same filerzap converters, so the two adapters stay in lockstep on
// the wire numbering. Replication only invokes unary RPCs today; the streaming
// methods are implemented for full interface fidelity.
package filerclient

import (
	"context"
	"io"

	"github.com/hanzoai/s3/s3/svc/filer"
	filer_pb "github.com/hanzoai/s3/s3/pb/filer_pb"
	"github.com/hanzoai/s3/s3/pb/rpc"
	filerwire "github.com/hanzoai/s3/s3/wire/filer"
	filerstream "github.com/hanzoai/s3/s3/wire/filer/filerstream"

	"github.com/zap-proto/go/transport"
)

// adapter binds a transport.Conn to the filer_pb.HanzoFilerClient contract.
// unary issues calls over the filerwire client; stream opens streams over the
// filerstream client on the SAME connection.
type adapter struct {
	conn   transport.Conn
	unary  *filerwire.HanzoFilerClient
	stream *filerstream.Client
}

// New wraps an established transport.Conn as a filer_pb.HanzoFilerClient that
// routes over the native ZAP transport. The conn's lifetime is owned by the
// caller (close it when the call scope returns).
func New(conn transport.Conn) filer_pb.HanzoFilerClient {
	return &adapter{
		conn:   conn,
		unary:  filerwire.NewHanzoFilerClient(conn, nil),
		stream: filerstream.NewClient(conn),
	}
}

var _ filer_pb.HanzoFilerClient = (*adapter)(nil)

// --- unary RPCs ------------------------------------------------------------

func (a *adapter) LookupDirectoryEntry(ctx context.Context, in *filer_pb.LookupDirectoryEntryRequest) (*filer_pb.LookupDirectoryEntryResponse, error) {
	_, body, err := a.unary.LookupDirectoryEntry(ctx, filer.LookupDirectoryEntryReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filer.LookupDirectoryEntryRespFromWire(body)
}

func (a *adapter) CreateEntry(ctx context.Context, in *filer_pb.CreateEntryRequest) (*filer_pb.CreateEntryResponse, error) {
	_, body, err := a.unary.CreateEntry(ctx, filer.CreateEntryReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filer.CreateEntryRespFromWire(body)
}

func (a *adapter) UpdateEntry(ctx context.Context, in *filer_pb.UpdateEntryRequest) (*filer_pb.UpdateEntryResponse, error) {
	_, body, err := a.unary.UpdateEntry(ctx, filer.UpdateEntryReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filer.UpdateEntryRespFromWire(body)
}

func (a *adapter) TouchAccessTime(ctx context.Context, in *filer_pb.TouchAccessTimeRequest) (*filer_pb.TouchAccessTimeResponse, error) {
	_, body, err := a.unary.TouchAccessTime(ctx, filer.TouchAccessTimeReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filer.TouchAccessTimeRespFromWire(body)
}

func (a *adapter) AppendToEntry(ctx context.Context, in *filer_pb.AppendToEntryRequest) (*filer_pb.AppendToEntryResponse, error) {
	_, body, err := a.unary.AppendToEntry(ctx, filer.AppendToEntryReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filer.AppendToEntryRespFromWire(body)
}

func (a *adapter) DeleteEntry(ctx context.Context, in *filer_pb.DeleteEntryRequest) (*filer_pb.DeleteEntryResponse, error) {
	_, body, err := a.unary.DeleteEntry(ctx, filer.DeleteEntryReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filer.DeleteEntryRespFromWire(body)
}

func (a *adapter) ObjectTransaction(ctx context.Context, in *filer_pb.ObjectTransactionRequest) (*filer_pb.ObjectTransactionResponse, error) {
	_, body, err := a.unary.ObjectTransaction(ctx, filer.ObjectTransactionReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filer.ObjectTransactionRespFromWire(body)
}

func (a *adapter) ObjectTransactionBatch(ctx context.Context, in *filer_pb.ObjectTransactionBatchRequest) (*filer_pb.ObjectTransactionBatchResponse, error) {
	_, body, err := a.unary.ObjectTransactionBatch(ctx, filer.ObjectTransactionBatchReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filer.ObjectTransactionBatchRespFromWire(body)
}

func (a *adapter) PosixLock(ctx context.Context, in *filer_pb.PosixLockRequest) (*filer_pb.PosixLockResponse, error) {
	_, body, err := a.unary.PosixLock(ctx, filer.PosixLockReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filer.PosixLockRespFromWire(body)
}

func (a *adapter) AtomicRenameEntry(ctx context.Context, in *filer_pb.AtomicRenameEntryRequest) (*filer_pb.AtomicRenameEntryResponse, error) {
	_, body, err := a.unary.AtomicRenameEntry(ctx, filer.AtomicRenameEntryReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filer.AtomicRenameEntryRespFromWire(body)
}

func (a *adapter) AssignVolume(ctx context.Context, in *filer_pb.AssignVolumeRequest) (*filer_pb.AssignVolumeResponse, error) {
	_, body, err := a.unary.AssignVolume(ctx, filer.AssignVolumeReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filer.AssignVolumeRespFromWire(body)
}

func (a *adapter) LookupVolume(ctx context.Context, in *filer_pb.LookupVolumeRequest) (*filer_pb.LookupVolumeResponse, error) {
	_, body, err := a.unary.LookupVolume(ctx, filer.LookupVolumeReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filer.LookupVolumeRespFromWire(body)
}

func (a *adapter) CollectionList(ctx context.Context, in *filer_pb.CollectionListRequest) (*filer_pb.CollectionListResponse, error) {
	_, body, err := a.unary.CollectionList(ctx, filer.CollectionListReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filer.CollectionListRespFromWire(body)
}

func (a *adapter) DeleteCollection(ctx context.Context, in *filer_pb.DeleteCollectionRequest) (*filer_pb.DeleteCollectionResponse, error) {
	_, body, err := a.unary.DeleteCollection(ctx, filer.DeleteCollectionReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filer.DeleteCollectionRespFromWire(body)
}

func (a *adapter) Statistics(ctx context.Context, in *filer_pb.StatisticsRequest) (*filer_pb.StatisticsResponse, error) {
	_, body, err := a.unary.Statistics(ctx, filer.StatisticsReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filer.StatisticsRespFromWire(body)
}

func (a *adapter) Ping(ctx context.Context, in *filer_pb.PingRequest) (*filer_pb.PingResponse, error) {
	_, body, err := a.unary.Ping(ctx, filer.PingReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filer.PingRespFromWire(body)
}

func (a *adapter) GetFilerConfiguration(ctx context.Context, in *filer_pb.GetFilerConfigurationRequest) (*filer_pb.GetFilerConfigurationResponse, error) {
	_, body, err := a.unary.GetFilerConfiguration(ctx, filer.GetFilerConfigurationReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filer.GetFilerConfigurationRespFromWire(body)
}

func (a *adapter) ListMetadataSubscribers(ctx context.Context, in *filer_pb.ListMetadataSubscribersRequest) (*filer_pb.ListMetadataSubscribersResponse, error) {
	_, body, err := a.unary.ListMetadataSubscribers(ctx, filer.ListMetadataSubscribersReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filer.ListMetadataSubscribersRespFromWire(body)
}

func (a *adapter) KvGet(ctx context.Context, in *filer_pb.KvGetRequest) (*filer_pb.KvGetResponse, error) {
	_, body, err := a.unary.KvGet(ctx, filer.KvGetReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filer.KvGetRespFromWire(body)
}

func (a *adapter) KvPut(ctx context.Context, in *filer_pb.KvPutRequest) (*filer_pb.KvPutResponse, error) {
	_, body, err := a.unary.KvPut(ctx, filer.KvPutReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filer.KvPutRespFromWire(body)
}

func (a *adapter) CacheRemoteObjectToLocalCluster(ctx context.Context, in *filer_pb.CacheRemoteObjectToLocalClusterRequest) (*filer_pb.CacheRemoteObjectToLocalClusterResponse, error) {
	_, body, err := a.unary.CacheRemoteObjectToLocalCluster(ctx, filer.CacheRemoteObjectToLocalClusterReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filer.CacheRemoteObjectToLocalClusterRespFromWire(body)
}

func (a *adapter) DistributedLock(ctx context.Context, in *filer_pb.LockRequest) (*filer_pb.LockResponse, error) {
	_, body, err := a.unary.DistributedLock(ctx, filer.DistributedLockReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filer.DistributedLockRespFromWire(body)
}

func (a *adapter) DistributedUnlock(ctx context.Context, in *filer_pb.UnlockRequest) (*filer_pb.UnlockResponse, error) {
	_, body, err := a.unary.DistributedUnlock(ctx, filer.DistributedUnlockReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filer.DistributedUnlockRespFromWire(body)
}

func (a *adapter) FindLockOwner(ctx context.Context, in *filer_pb.FindLockOwnerRequest) (*filer_pb.FindLockOwnerResponse, error) {
	_, body, err := a.unary.FindLockOwner(ctx, filer.FindLockOwnerReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filer.FindLockOwnerRespFromWire(body)
}

func (a *adapter) TransferLocks(ctx context.Context, in *filer_pb.TransferLocksRequest) (*filer_pb.TransferLocksResponse, error) {
	_, body, err := a.unary.TransferLocks(ctx, filer.TransferLocksReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filer.TransferLocksRespFromWire(body)
}

func (a *adapter) ReplicateLock(ctx context.Context, in *filer_pb.ReplicateLockRequest) (*filer_pb.ReplicateLockResponse, error) {
	_, body, err := a.unary.ReplicateLock(ctx, filer.ReplicateLockReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filer.ReplicateLockRespFromWire(body)
}

func (a *adapter) MountRegister(ctx context.Context, in *filer_pb.MountRegisterRequest) (*filer_pb.MountRegisterResponse, error) {
	_, body, err := a.unary.MountRegister(ctx, filer.MountRegisterReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filer.MountRegisterRespFromWire(body)
}

func (a *adapter) MountList(ctx context.Context, in *filer_pb.MountListRequest) (*filer_pb.MountListResponse, error) {
	_, body, err := a.unary.MountList(ctx, filer.MountListReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filer.MountListRespFromWire(body)
}

// --- streaming RPCs --------------------------------------------------------

func (a *adapter) ListEntries(ctx context.Context, in *filer_pb.ListEntriesRequest) (rpc.ServerStream[filer_pb.ListEntriesResponse], error) {
	s, err := a.stream.ListEntries(filer.ListEntriesReqToWire(in))
	if err != nil {
		return nil, err
	}
	return &listEntriesStream{s: s}, nil
}

func (a *adapter) StreamRenameEntry(ctx context.Context, in *filer_pb.StreamRenameEntryRequest) (rpc.ServerStream[filer_pb.StreamRenameEntryResponse], error) {
	s, err := a.stream.StreamRenameEntry(filer.StreamRenameEntryReqToWire(in))
	if err != nil {
		return nil, err
	}
	return &streamRenameStream{s: s}, nil
}

func (a *adapter) TraverseBfsMetadata(ctx context.Context, in *filer_pb.TraverseBfsMetadataRequest) (rpc.ServerStream[filer_pb.TraverseBfsMetadataResponse], error) {
	s, err := a.stream.TraverseBfsMetadata(filer.TraverseBfsMetadataReqToWire(in))
	if err != nil {
		return nil, err
	}
	return &traverseBfsStream{s: s}, nil
}

func (a *adapter) SubscribeMetadata(ctx context.Context, in *filer_pb.SubscribeMetadataRequest) (rpc.ServerStream[filer_pb.SubscribeMetadataResponse], error) {
	s, err := a.stream.SubscribeMetadata(filer.SubscribeMetadataReqToWire(in))
	if err != nil {
		return nil, err
	}
	return &subscribeMetadataStream{s: s}, nil
}

func (a *adapter) SubscribeLocalMetadata(ctx context.Context, in *filer_pb.SubscribeMetadataRequest) (rpc.ServerStream[filer_pb.SubscribeMetadataResponse], error) {
	s, err := a.stream.SubscribeLocalMetadata(filer.SubscribeMetadataReqToWire(in))
	if err != nil {
		return nil, err
	}
	return &subscribeMetadataStream{s: s}, nil
}

func (a *adapter) StreamMutateEntry(_ context.Context) (rpc.BidiStream[filer_pb.StreamMutateEntryRequest, filer_pb.StreamMutateEntryResponse], error) {
	// Open the bidi stream with an empty (Which=0, "no member") request frame so
	// the server has a parseable opener that yields no response; requests then
	// flow via Send.
	s, err := a.stream.StreamMutateEntry(filerwire.NewStreamMutateEntryRequest(filerwire.StreamMutateEntryRequestInput{}))
	if err != nil {
		return nil, err
	}
	return &streamMutateStream{s: s}, nil
}

// The stream wrappers below implement exactly the rpc.* seam each RPC returns:
// rpc.ServerStream needs Recv (CloseSend is offered to release the stream early);
// the bidi wrapper adds Send.

type listEntriesStream struct {
	s *filerstream.ClientListEntriesStream
}

func (x *listEntriesStream) Recv() (*filer_pb.ListEntriesResponse, error) {
	b, err := x.s.RecvBytes()
	if err != nil {
		return nil, err
	}
	return filer.ListEntriesRespFromWire(b)
}
func (x *listEntriesStream) CloseSend() error { return x.s.Close() }

type streamRenameStream struct {
	s *filerstream.ClientStreamRenameEntryStream
}

func (x *streamRenameStream) Recv() (*filer_pb.StreamRenameEntryResponse, error) {
	b, err := x.s.RecvBytes()
	if err != nil {
		return nil, err
	}
	return filer.StreamRenameEntryRespFromWire(b)
}
func (x *streamRenameStream) CloseSend() error { return x.s.Close() }

type traverseBfsStream struct {
	s *filerstream.ClientTraverseBfsMetadataStream
}

func (x *traverseBfsStream) Recv() (*filer_pb.TraverseBfsMetadataResponse, error) {
	// replication does not consume TraverseBfsMetadata frames; drain to EOF.
	if _, err := x.s.RecvBytes(); err != nil {
		return nil, err
	}
	return nil, io.EOF
}
func (x *traverseBfsStream) CloseSend() error { return x.s.Close() }

type subscribeMetadataStream struct {
	s *filerstream.ClientSubscribeMetadataStream
}

func (x *subscribeMetadataStream) Recv() (*filer_pb.SubscribeMetadataResponse, error) {
	b, err := x.s.RecvBytes()
	if err != nil {
		return nil, err
	}
	return filer.SubscribeMetadataResponseFromWire(b)
}
func (x *subscribeMetadataStream) CloseSend() error { return x.s.Close() }

type streamMutateStream struct {
	s *filerstream.ClientStreamMutateEntryStream
}

func (x *streamMutateStream) Send(in *filer_pb.StreamMutateEntryRequest) error {
	return x.s.Send(filer.StreamMutateEntryReqToWire(in))
}
func (x *streamMutateStream) Recv() (*filer_pb.StreamMutateEntryResponse, error) {
	b, err := x.s.RecvBytes()
	if err != nil {
		return nil, err
	}
	return filer.StreamMutateEntryRespFromWire(b)
}
func (x *streamMutateStream) CloseSend() error { return x.s.CloseSend() }
