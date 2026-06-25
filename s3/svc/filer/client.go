// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// client.go is the canonical strangler seam from the filer_pb.HanzoFilerClient
// contract onto the native ZAP transport. It lives in package filerzap — the
// package that owns the per-RPC <Rpc>ReqToWire/<Rpc>RespFromWire converters — so
// the contract is bridged in exactly ONE place. It implements
// filer_pb.HanzoFilerClient but routes every call over a
// github.com/zap-proto/go transport.Conn instead of gRPC:
//
//   - the 28 unary RPCs go through filerwire.HanzoFilerClient (over the conn's
//     Call channel), translating *filer_pb.<Rpc>Request <-> ZAP buffer <->
//     *filer_pb.<Rpc>Response with the converters in rpc.go;
//   - the 6 streaming RPCs open a transport stream via the filerstream client
//     and return an rpc.ServerStream / rpc.BidiStream whose Recv()/Send()
//     decode/encode each frame with the same converters.
//
// Because the adapter satisfies filer_pb.HanzoFilerClient, no call site changes:
// callers keep building *filer_pb requests and reading *filer_pb responses. The
// bytes are ZAP at every hop; protobuf framing and gRPC are gone. The package-pb
// WithFilerClient helpers (pb.NewZapFilerClient), the mount FUSE adapter, and the
// replication source/sink all dial the transport and hand fn a client from here,
// so this is the ONE place the contract is bridged.

package filer

import (
	"context"
	"io"

	"github.com/zap-proto/go/transport"

	filer_pb "github.com/hanzoai/s3/s3/pb/filer_pb"
	"github.com/hanzoai/s3/s3/pb/rpc"
	filerwire "github.com/hanzoai/s3/s3/wire/filer"
	filerstream "github.com/hanzoai/s3/s3/wire/filer/filerstream"
)

// zapFilerClient binds a transport.Conn to the filer_pb.HanzoFilerClient
// contract. unary issues calls over the filerwire client; stream opens streams
// over the filerstream client on the SAME connection.
type zapFilerClient struct {
	conn   transport.Conn
	unary  *filerwire.HanzoFilerClient
	stream *filerstream.Client
}

// NewZapFilerClient wraps an established transport.Conn as a
// filer_pb.HanzoFilerClient that routes over ZAP. cap is the optional capability
// token attached to every unary request (nil for none). The caller owns conn's
// lifecycle (Close when done); the adapter pools nothing.
func NewZapFilerClient(conn transport.Conn, capability []byte) filer_pb.HanzoFilerClient {
	return &zapFilerClient{
		conn:   conn,
		unary:  filerwire.NewHanzoFilerClient(conn, capability),
		stream: filerstream.NewClient(conn),
	}
}

var _ filer_pb.HanzoFilerClient = (*zapFilerClient)(nil)

// --- unary RPCs ------------------------------------------------------------

func (a *zapFilerClient) LookupDirectoryEntry(ctx context.Context, in *filer_pb.LookupDirectoryEntryRequest) (*filer_pb.LookupDirectoryEntryResponse, error) {
	_, body, err := a.unary.LookupDirectoryEntry(ctx, LookupDirectoryEntryReqToWire(in))
	if err != nil {
		return nil, err
	}
	return LookupDirectoryEntryRespFromWire(body)
}

func (a *zapFilerClient) CreateEntry(ctx context.Context, in *filer_pb.CreateEntryRequest) (*filer_pb.CreateEntryResponse, error) {
	_, body, err := a.unary.CreateEntry(ctx, CreateEntryReqToWire(in))
	if err != nil {
		return nil, err
	}
	return CreateEntryRespFromWire(body)
}

func (a *zapFilerClient) UpdateEntry(ctx context.Context, in *filer_pb.UpdateEntryRequest) (*filer_pb.UpdateEntryResponse, error) {
	_, body, err := a.unary.UpdateEntry(ctx, UpdateEntryReqToWire(in))
	if err != nil {
		return nil, err
	}
	return UpdateEntryRespFromWire(body)
}

func (a *zapFilerClient) TouchAccessTime(ctx context.Context, in *filer_pb.TouchAccessTimeRequest) (*filer_pb.TouchAccessTimeResponse, error) {
	_, body, err := a.unary.TouchAccessTime(ctx, TouchAccessTimeReqToWire(in))
	if err != nil {
		return nil, err
	}
	return TouchAccessTimeRespFromWire(body)
}

func (a *zapFilerClient) AppendToEntry(ctx context.Context, in *filer_pb.AppendToEntryRequest) (*filer_pb.AppendToEntryResponse, error) {
	_, body, err := a.unary.AppendToEntry(ctx, AppendToEntryReqToWire(in))
	if err != nil {
		return nil, err
	}
	return AppendToEntryRespFromWire(body)
}

func (a *zapFilerClient) DeleteEntry(ctx context.Context, in *filer_pb.DeleteEntryRequest) (*filer_pb.DeleteEntryResponse, error) {
	_, body, err := a.unary.DeleteEntry(ctx, DeleteEntryReqToWire(in))
	if err != nil {
		return nil, err
	}
	return DeleteEntryRespFromWire(body)
}

func (a *zapFilerClient) ObjectTransaction(ctx context.Context, in *filer_pb.ObjectTransactionRequest) (*filer_pb.ObjectTransactionResponse, error) {
	_, body, err := a.unary.ObjectTransaction(ctx, ObjectTransactionReqToWire(in))
	if err != nil {
		return nil, err
	}
	return ObjectTransactionRespFromWire(body)
}

func (a *zapFilerClient) ObjectTransactionBatch(ctx context.Context, in *filer_pb.ObjectTransactionBatchRequest) (*filer_pb.ObjectTransactionBatchResponse, error) {
	_, body, err := a.unary.ObjectTransactionBatch(ctx, ObjectTransactionBatchReqToWire(in))
	if err != nil {
		return nil, err
	}
	return ObjectTransactionBatchRespFromWire(body)
}

func (a *zapFilerClient) PosixLock(ctx context.Context, in *filer_pb.PosixLockRequest) (*filer_pb.PosixLockResponse, error) {
	_, body, err := a.unary.PosixLock(ctx, PosixLockReqToWire(in))
	if err != nil {
		return nil, err
	}
	return PosixLockRespFromWire(body)
}

func (a *zapFilerClient) AtomicRenameEntry(ctx context.Context, in *filer_pb.AtomicRenameEntryRequest) (*filer_pb.AtomicRenameEntryResponse, error) {
	_, body, err := a.unary.AtomicRenameEntry(ctx, AtomicRenameEntryReqToWire(in))
	if err != nil {
		return nil, err
	}
	return AtomicRenameEntryRespFromWire(body)
}

func (a *zapFilerClient) AssignVolume(ctx context.Context, in *filer_pb.AssignVolumeRequest) (*filer_pb.AssignVolumeResponse, error) {
	_, body, err := a.unary.AssignVolume(ctx, AssignVolumeReqToWire(in))
	if err != nil {
		return nil, err
	}
	return AssignVolumeRespFromWire(body)
}

func (a *zapFilerClient) LookupVolume(ctx context.Context, in *filer_pb.LookupVolumeRequest) (*filer_pb.LookupVolumeResponse, error) {
	_, body, err := a.unary.LookupVolume(ctx, LookupVolumeReqToWire(in))
	if err != nil {
		return nil, err
	}
	return LookupVolumeRespFromWire(body)
}

func (a *zapFilerClient) CollectionList(ctx context.Context, in *filer_pb.CollectionListRequest) (*filer_pb.CollectionListResponse, error) {
	_, body, err := a.unary.CollectionList(ctx, CollectionListReqToWire(in))
	if err != nil {
		return nil, err
	}
	return CollectionListRespFromWire(body)
}

func (a *zapFilerClient) DeleteCollection(ctx context.Context, in *filer_pb.DeleteCollectionRequest) (*filer_pb.DeleteCollectionResponse, error) {
	_, body, err := a.unary.DeleteCollection(ctx, DeleteCollectionReqToWire(in))
	if err != nil {
		return nil, err
	}
	return DeleteCollectionRespFromWire(body)
}

func (a *zapFilerClient) Statistics(ctx context.Context, in *filer_pb.StatisticsRequest) (*filer_pb.StatisticsResponse, error) {
	_, body, err := a.unary.Statistics(ctx, StatisticsReqToWire(in))
	if err != nil {
		return nil, err
	}
	return StatisticsRespFromWire(body)
}

func (a *zapFilerClient) Ping(ctx context.Context, in *filer_pb.PingRequest) (*filer_pb.PingResponse, error) {
	_, body, err := a.unary.Ping(ctx, PingReqToWire(in))
	if err != nil {
		return nil, err
	}
	return PingRespFromWire(body)
}

func (a *zapFilerClient) GetFilerConfiguration(ctx context.Context, in *filer_pb.GetFilerConfigurationRequest) (*filer_pb.GetFilerConfigurationResponse, error) {
	_, body, err := a.unary.GetFilerConfiguration(ctx, GetFilerConfigurationReqToWire(in))
	if err != nil {
		return nil, err
	}
	return GetFilerConfigurationRespFromWire(body)
}

func (a *zapFilerClient) ListMetadataSubscribers(ctx context.Context, in *filer_pb.ListMetadataSubscribersRequest) (*filer_pb.ListMetadataSubscribersResponse, error) {
	_, body, err := a.unary.ListMetadataSubscribers(ctx, ListMetadataSubscribersReqToWire(in))
	if err != nil {
		return nil, err
	}
	return ListMetadataSubscribersRespFromWire(body)
}

func (a *zapFilerClient) KvGet(ctx context.Context, in *filer_pb.KvGetRequest) (*filer_pb.KvGetResponse, error) {
	_, body, err := a.unary.KvGet(ctx, KvGetReqToWire(in))
	if err != nil {
		return nil, err
	}
	return KvGetRespFromWire(body)
}

func (a *zapFilerClient) KvPut(ctx context.Context, in *filer_pb.KvPutRequest) (*filer_pb.KvPutResponse, error) {
	_, body, err := a.unary.KvPut(ctx, KvPutReqToWire(in))
	if err != nil {
		return nil, err
	}
	return KvPutRespFromWire(body)
}

func (a *zapFilerClient) CacheRemoteObjectToLocalCluster(ctx context.Context, in *filer_pb.CacheRemoteObjectToLocalClusterRequest) (*filer_pb.CacheRemoteObjectToLocalClusterResponse, error) {
	_, body, err := a.unary.CacheRemoteObjectToLocalCluster(ctx, CacheRemoteObjectToLocalClusterReqToWire(in))
	if err != nil {
		return nil, err
	}
	return CacheRemoteObjectToLocalClusterRespFromWire(body)
}

func (a *zapFilerClient) DistributedLock(ctx context.Context, in *filer_pb.LockRequest) (*filer_pb.LockResponse, error) {
	_, body, err := a.unary.DistributedLock(ctx, DistributedLockReqToWire(in))
	if err != nil {
		return nil, err
	}
	return DistributedLockRespFromWire(body)
}

func (a *zapFilerClient) DistributedUnlock(ctx context.Context, in *filer_pb.UnlockRequest) (*filer_pb.UnlockResponse, error) {
	_, body, err := a.unary.DistributedUnlock(ctx, DistributedUnlockReqToWire(in))
	if err != nil {
		return nil, err
	}
	return DistributedUnlockRespFromWire(body)
}

func (a *zapFilerClient) FindLockOwner(ctx context.Context, in *filer_pb.FindLockOwnerRequest) (*filer_pb.FindLockOwnerResponse, error) {
	_, body, err := a.unary.FindLockOwner(ctx, FindLockOwnerReqToWire(in))
	if err != nil {
		return nil, err
	}
	return FindLockOwnerRespFromWire(body)
}

func (a *zapFilerClient) TransferLocks(ctx context.Context, in *filer_pb.TransferLocksRequest) (*filer_pb.TransferLocksResponse, error) {
	_, body, err := a.unary.TransferLocks(ctx, TransferLocksReqToWire(in))
	if err != nil {
		return nil, err
	}
	return TransferLocksRespFromWire(body)
}

func (a *zapFilerClient) ReplicateLock(ctx context.Context, in *filer_pb.ReplicateLockRequest) (*filer_pb.ReplicateLockResponse, error) {
	_, body, err := a.unary.ReplicateLock(ctx, ReplicateLockReqToWire(in))
	if err != nil {
		return nil, err
	}
	return ReplicateLockRespFromWire(body)
}

func (a *zapFilerClient) MountRegister(ctx context.Context, in *filer_pb.MountRegisterRequest) (*filer_pb.MountRegisterResponse, error) {
	_, body, err := a.unary.MountRegister(ctx, MountRegisterReqToWire(in))
	if err != nil {
		return nil, err
	}
	return MountRegisterRespFromWire(body)
}

func (a *zapFilerClient) MountList(ctx context.Context, in *filer_pb.MountListRequest) (*filer_pb.MountListResponse, error) {
	_, body, err := a.unary.MountList(ctx, MountListReqToWire(in))
	if err != nil {
		return nil, err
	}
	return MountListRespFromWire(body)
}

// --- streaming RPCs --------------------------------------------------------

func (a *zapFilerClient) ListEntries(ctx context.Context, in *filer_pb.ListEntriesRequest) (rpc.ServerStream[filer_pb.ListEntriesResponse], error) {
	s, err := a.stream.ListEntries(ListEntriesReqToWire(in))
	if err != nil {
		return nil, err
	}
	return &zapListEntriesClientStream{s: s}, nil
}

func (a *zapFilerClient) StreamRenameEntry(ctx context.Context, in *filer_pb.StreamRenameEntryRequest) (rpc.ServerStream[filer_pb.StreamRenameEntryResponse], error) {
	s, err := a.stream.StreamRenameEntry(StreamRenameEntryReqToWire(in))
	if err != nil {
		return nil, err
	}
	return &zapStreamRenameClientStream{s: s}, nil
}

func (a *zapFilerClient) TraverseBfsMetadata(ctx context.Context, in *filer_pb.TraverseBfsMetadataRequest) (rpc.ServerStream[filer_pb.TraverseBfsMetadataResponse], error) {
	s, err := a.stream.TraverseBfsMetadata(TraverseBfsMetadataReqToWire(in))
	if err != nil {
		return nil, err
	}
	return &zapTraverseBfsClientStream{s: s}, nil
}

func (a *zapFilerClient) SubscribeMetadata(ctx context.Context, in *filer_pb.SubscribeMetadataRequest) (rpc.ServerStream[filer_pb.SubscribeMetadataResponse], error) {
	s, err := a.stream.SubscribeMetadata(SubscribeMetadataReqToWire(in))
	if err != nil {
		return nil, err
	}
	return &zapSubscribeMetadataClientStream{s: s}, nil
}

func (a *zapFilerClient) SubscribeLocalMetadata(ctx context.Context, in *filer_pb.SubscribeMetadataRequest) (rpc.ServerStream[filer_pb.SubscribeMetadataResponse], error) {
	s, err := a.stream.SubscribeLocalMetadata(SubscribeMetadataReqToWire(in))
	if err != nil {
		return nil, err
	}
	return &zapSubscribeMetadataClientStream{s: s}, nil
}

func (a *zapFilerClient) StreamMutateEntry(_ context.Context) (rpc.BidiStream[filer_pb.StreamMutateEntryRequest, filer_pb.StreamMutateEntryResponse], error) {
	// The bidi stream carries requests via Send; open it with an empty (Which=0,
	// "no member") request frame so the server has a parseable opener that yields
	// no response.
	s, err := a.stream.StreamMutateEntry(filerwire.NewStreamMutateEntryRequest(filerwire.StreamMutateEntryRequestInput{}))
	if err != nil {
		return nil, err
	}
	return &zapStreamMutateClientStream{s: s}, nil
}

// The stream wrappers below implement exactly the rpc.* seam each RPC returns:
// rpc.ServerStream needs Recv (CloseSend is offered too so callers can release
// the stream early); the bidi wrapper adds Send. Each frame is decoded/encoded
// with the same converters as the unary path.

// zapListEntriesClientStream adapts a filerstream ListEntries stream to
// rpc.ServerStream[ListEntriesResponse].
type zapListEntriesClientStream struct {
	s *filerstream.ClientListEntriesStream
}

func (x *zapListEntriesClientStream) Recv() (*filer_pb.ListEntriesResponse, error) {
	b, err := x.s.RecvBytes()
	if err != nil {
		return nil, err
	}
	return ListEntriesRespFromWire(b)
}
func (x *zapListEntriesClientStream) CloseSend() error { return x.s.Close() }

// zapStreamRenameClientStream adapts a filerstream StreamRenameEntry stream.
type zapStreamRenameClientStream struct {
	s *filerstream.ClientStreamRenameEntryStream
}

func (x *zapStreamRenameClientStream) Recv() (*filer_pb.StreamRenameEntryResponse, error) {
	b, err := x.s.RecvBytes()
	if err != nil {
		return nil, err
	}
	return StreamRenameEntryRespFromWire(b)
}
func (x *zapStreamRenameClientStream) CloseSend() error { return x.s.Close() }

// zapTraverseBfsClientStream adapts a filerstream TraverseBfsMetadata stream. No
// package-level caller consumes TraverseBfsMetadata responses, and there is no
// TraverseBfsMetadataResponse decoder, so this drains to EOF — the adapter still
// implements the RPC to honor the full filer_pb.HanzoFilerClient contract.
type zapTraverseBfsClientStream struct {
	s *filerstream.ClientTraverseBfsMetadataStream
}

func (x *zapTraverseBfsClientStream) Recv() (*filer_pb.TraverseBfsMetadataResponse, error) {
	if _, err := x.s.RecvBytes(); err != nil {
		return nil, err
	}
	return nil, io.EOF
}
func (x *zapTraverseBfsClientStream) CloseSend() error { return x.s.Close() }

// zapSubscribeMetadataClientStream adapts a filerstream SubscribeMetadata stream
// (used for both cluster-wide and local subscriptions).
type zapSubscribeMetadataClientStream struct {
	s *filerstream.ClientSubscribeMetadataStream
}

func (x *zapSubscribeMetadataClientStream) Recv() (*filer_pb.SubscribeMetadataResponse, error) {
	b, err := x.s.RecvBytes()
	if err != nil {
		return nil, err
	}
	return SubscribeMetadataResponseFromWire(b)
}
func (x *zapSubscribeMetadataClientStream) CloseSend() error { return x.s.Close() }

// zapStreamMutateClientStream adapts the bidirectional StreamMutateEntry stream.
type zapStreamMutateClientStream struct {
	s *filerstream.ClientStreamMutateEntryStream
}

func (x *zapStreamMutateClientStream) Send(in *filer_pb.StreamMutateEntryRequest) error {
	return x.s.Send(StreamMutateEntryReqToWire(in))
}
func (x *zapStreamMutateClientStream) Recv() (*filer_pb.StreamMutateEntryResponse, error) {
	b, err := x.s.RecvBytes()
	if err != nil {
		return nil, err
	}
	return StreamMutateEntryRespFromWire(b)
}
func (x *zapStreamMutateClientStream) CloseSend() error { return x.s.CloseSend() }
