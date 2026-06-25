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
//     and return a grpc.ServerStreamingClient / grpc.BidiStreamingClient whose
//     Recv()/Send() decode/encode each frame with the same converters.
//
// Because the adapter satisfies filer_pb.HanzoFilerClient, no call site changes:
// callers keep building *filer_pb requests and reading *filer_pb responses. The
// bytes are ZAP at every hop; protobuf framing and gRPC are gone. The package-pb
// WithFilerClient helpers (pb.NewZapFilerClient), the mount FUSE adapter, and the
// replication source/sink all dial the transport and hand fn a client from here,
// so this is the ONE place the contract is bridged.

package filerzap

import (
	"context"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/zap-proto/go/transport"

	filer_pb "github.com/hanzoai/s3/s3/pb/filer_pb"
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

func (a *zapFilerClient) LookupDirectoryEntry(_ context.Context, in *filer_pb.LookupDirectoryEntryRequest, _ ...grpc.CallOption) (*filer_pb.LookupDirectoryEntryResponse, error) {
	_, body, err := a.unary.LookupDirectoryEntry(LookupDirectoryEntryReqToWire(in))
	if err != nil {
		return nil, err
	}
	return LookupDirectoryEntryRespFromWire(body)
}

func (a *zapFilerClient) CreateEntry(_ context.Context, in *filer_pb.CreateEntryRequest, _ ...grpc.CallOption) (*filer_pb.CreateEntryResponse, error) {
	_, body, err := a.unary.CreateEntry(CreateEntryReqToWire(in))
	if err != nil {
		return nil, err
	}
	return CreateEntryRespFromWire(body)
}

func (a *zapFilerClient) UpdateEntry(_ context.Context, in *filer_pb.UpdateEntryRequest, _ ...grpc.CallOption) (*filer_pb.UpdateEntryResponse, error) {
	_, body, err := a.unary.UpdateEntry(UpdateEntryReqToWire(in))
	if err != nil {
		return nil, err
	}
	return UpdateEntryRespFromWire(body)
}

func (a *zapFilerClient) TouchAccessTime(_ context.Context, in *filer_pb.TouchAccessTimeRequest, _ ...grpc.CallOption) (*filer_pb.TouchAccessTimeResponse, error) {
	_, body, err := a.unary.TouchAccessTime(TouchAccessTimeReqToWire(in))
	if err != nil {
		return nil, err
	}
	return TouchAccessTimeRespFromWire(body)
}

func (a *zapFilerClient) AppendToEntry(_ context.Context, in *filer_pb.AppendToEntryRequest, _ ...grpc.CallOption) (*filer_pb.AppendToEntryResponse, error) {
	_, body, err := a.unary.AppendToEntry(AppendToEntryReqToWire(in))
	if err != nil {
		return nil, err
	}
	return AppendToEntryRespFromWire(body)
}

func (a *zapFilerClient) DeleteEntry(_ context.Context, in *filer_pb.DeleteEntryRequest, _ ...grpc.CallOption) (*filer_pb.DeleteEntryResponse, error) {
	_, body, err := a.unary.DeleteEntry(DeleteEntryReqToWire(in))
	if err != nil {
		return nil, err
	}
	return DeleteEntryRespFromWire(body)
}

func (a *zapFilerClient) ObjectTransaction(_ context.Context, in *filer_pb.ObjectTransactionRequest, _ ...grpc.CallOption) (*filer_pb.ObjectTransactionResponse, error) {
	_, body, err := a.unary.ObjectTransaction(ObjectTransactionReqToWire(in))
	if err != nil {
		return nil, err
	}
	return ObjectTransactionRespFromWire(body)
}

func (a *zapFilerClient) ObjectTransactionBatch(_ context.Context, in *filer_pb.ObjectTransactionBatchRequest, _ ...grpc.CallOption) (*filer_pb.ObjectTransactionBatchResponse, error) {
	_, body, err := a.unary.ObjectTransactionBatch(ObjectTransactionBatchReqToWire(in))
	if err != nil {
		return nil, err
	}
	return ObjectTransactionBatchRespFromWire(body)
}

func (a *zapFilerClient) PosixLock(_ context.Context, in *filer_pb.PosixLockRequest, _ ...grpc.CallOption) (*filer_pb.PosixLockResponse, error) {
	_, body, err := a.unary.PosixLock(PosixLockReqToWire(in))
	if err != nil {
		return nil, err
	}
	return PosixLockRespFromWire(body)
}

func (a *zapFilerClient) AtomicRenameEntry(_ context.Context, in *filer_pb.AtomicRenameEntryRequest, _ ...grpc.CallOption) (*filer_pb.AtomicRenameEntryResponse, error) {
	_, body, err := a.unary.AtomicRenameEntry(AtomicRenameEntryReqToWire(in))
	if err != nil {
		return nil, err
	}
	return AtomicRenameEntryRespFromWire(body)
}

func (a *zapFilerClient) AssignVolume(_ context.Context, in *filer_pb.AssignVolumeRequest, _ ...grpc.CallOption) (*filer_pb.AssignVolumeResponse, error) {
	_, body, err := a.unary.AssignVolume(AssignVolumeReqToWire(in))
	if err != nil {
		return nil, err
	}
	return AssignVolumeRespFromWire(body)
}

func (a *zapFilerClient) LookupVolume(_ context.Context, in *filer_pb.LookupVolumeRequest, _ ...grpc.CallOption) (*filer_pb.LookupVolumeResponse, error) {
	_, body, err := a.unary.LookupVolume(LookupVolumeReqToWire(in))
	if err != nil {
		return nil, err
	}
	return LookupVolumeRespFromWire(body)
}

func (a *zapFilerClient) CollectionList(_ context.Context, in *filer_pb.CollectionListRequest, _ ...grpc.CallOption) (*filer_pb.CollectionListResponse, error) {
	_, body, err := a.unary.CollectionList(CollectionListReqToWire(in))
	if err != nil {
		return nil, err
	}
	return CollectionListRespFromWire(body)
}

func (a *zapFilerClient) DeleteCollection(_ context.Context, in *filer_pb.DeleteCollectionRequest, _ ...grpc.CallOption) (*filer_pb.DeleteCollectionResponse, error) {
	_, body, err := a.unary.DeleteCollection(DeleteCollectionReqToWire(in))
	if err != nil {
		return nil, err
	}
	return DeleteCollectionRespFromWire(body)
}

func (a *zapFilerClient) Statistics(_ context.Context, in *filer_pb.StatisticsRequest, _ ...grpc.CallOption) (*filer_pb.StatisticsResponse, error) {
	_, body, err := a.unary.Statistics(StatisticsReqToWire(in))
	if err != nil {
		return nil, err
	}
	return StatisticsRespFromWire(body)
}

func (a *zapFilerClient) Ping(_ context.Context, in *filer_pb.PingRequest, _ ...grpc.CallOption) (*filer_pb.PingResponse, error) {
	_, body, err := a.unary.Ping(PingReqToWire(in))
	if err != nil {
		return nil, err
	}
	return PingRespFromWire(body)
}

func (a *zapFilerClient) GetFilerConfiguration(_ context.Context, in *filer_pb.GetFilerConfigurationRequest, _ ...grpc.CallOption) (*filer_pb.GetFilerConfigurationResponse, error) {
	_, body, err := a.unary.GetFilerConfiguration(GetFilerConfigurationReqToWire(in))
	if err != nil {
		return nil, err
	}
	return GetFilerConfigurationRespFromWire(body)
}

func (a *zapFilerClient) ListMetadataSubscribers(_ context.Context, in *filer_pb.ListMetadataSubscribersRequest, _ ...grpc.CallOption) (*filer_pb.ListMetadataSubscribersResponse, error) {
	_, body, err := a.unary.ListMetadataSubscribers(ListMetadataSubscribersReqToWire(in))
	if err != nil {
		return nil, err
	}
	return ListMetadataSubscribersRespFromWire(body)
}

func (a *zapFilerClient) KvGet(_ context.Context, in *filer_pb.KvGetRequest, _ ...grpc.CallOption) (*filer_pb.KvGetResponse, error) {
	_, body, err := a.unary.KvGet(KvGetReqToWire(in))
	if err != nil {
		return nil, err
	}
	return KvGetRespFromWire(body)
}

func (a *zapFilerClient) KvPut(_ context.Context, in *filer_pb.KvPutRequest, _ ...grpc.CallOption) (*filer_pb.KvPutResponse, error) {
	_, body, err := a.unary.KvPut(KvPutReqToWire(in))
	if err != nil {
		return nil, err
	}
	return KvPutRespFromWire(body)
}

func (a *zapFilerClient) CacheRemoteObjectToLocalCluster(_ context.Context, in *filer_pb.CacheRemoteObjectToLocalClusterRequest, _ ...grpc.CallOption) (*filer_pb.CacheRemoteObjectToLocalClusterResponse, error) {
	_, body, err := a.unary.CacheRemoteObjectToLocalCluster(CacheRemoteObjectToLocalClusterReqToWire(in))
	if err != nil {
		return nil, err
	}
	return CacheRemoteObjectToLocalClusterRespFromWire(body)
}

func (a *zapFilerClient) DistributedLock(_ context.Context, in *filer_pb.LockRequest, _ ...grpc.CallOption) (*filer_pb.LockResponse, error) {
	_, body, err := a.unary.DistributedLock(DistributedLockReqToWire(in))
	if err != nil {
		return nil, err
	}
	return DistributedLockRespFromWire(body)
}

func (a *zapFilerClient) DistributedUnlock(_ context.Context, in *filer_pb.UnlockRequest, _ ...grpc.CallOption) (*filer_pb.UnlockResponse, error) {
	_, body, err := a.unary.DistributedUnlock(DistributedUnlockReqToWire(in))
	if err != nil {
		return nil, err
	}
	return DistributedUnlockRespFromWire(body)
}

func (a *zapFilerClient) FindLockOwner(_ context.Context, in *filer_pb.FindLockOwnerRequest, _ ...grpc.CallOption) (*filer_pb.FindLockOwnerResponse, error) {
	_, body, err := a.unary.FindLockOwner(FindLockOwnerReqToWire(in))
	if err != nil {
		return nil, err
	}
	return FindLockOwnerRespFromWire(body)
}

func (a *zapFilerClient) TransferLocks(_ context.Context, in *filer_pb.TransferLocksRequest, _ ...grpc.CallOption) (*filer_pb.TransferLocksResponse, error) {
	_, body, err := a.unary.TransferLocks(TransferLocksReqToWire(in))
	if err != nil {
		return nil, err
	}
	return TransferLocksRespFromWire(body)
}

func (a *zapFilerClient) ReplicateLock(_ context.Context, in *filer_pb.ReplicateLockRequest, _ ...grpc.CallOption) (*filer_pb.ReplicateLockResponse, error) {
	_, body, err := a.unary.ReplicateLock(ReplicateLockReqToWire(in))
	if err != nil {
		return nil, err
	}
	return ReplicateLockRespFromWire(body)
}

func (a *zapFilerClient) MountRegister(_ context.Context, in *filer_pb.MountRegisterRequest, _ ...grpc.CallOption) (*filer_pb.MountRegisterResponse, error) {
	_, body, err := a.unary.MountRegister(MountRegisterReqToWire(in))
	if err != nil {
		return nil, err
	}
	return MountRegisterRespFromWire(body)
}

func (a *zapFilerClient) MountList(_ context.Context, in *filer_pb.MountListRequest, _ ...grpc.CallOption) (*filer_pb.MountListResponse, error) {
	_, body, err := a.unary.MountList(MountListReqToWire(in))
	if err != nil {
		return nil, err
	}
	return MountListRespFromWire(body)
}

// --- streaming RPCs --------------------------------------------------------

func (a *zapFilerClient) ListEntries(_ context.Context, in *filer_pb.ListEntriesRequest, _ ...grpc.CallOption) (grpc.ServerStreamingClient[filer_pb.ListEntriesResponse], error) {
	s, err := a.stream.ListEntries(ListEntriesReqToWire(in))
	if err != nil {
		return nil, err
	}
	return &zapListEntriesClientStream{s: s}, nil
}

func (a *zapFilerClient) StreamRenameEntry(_ context.Context, in *filer_pb.StreamRenameEntryRequest, _ ...grpc.CallOption) (grpc.ServerStreamingClient[filer_pb.StreamRenameEntryResponse], error) {
	s, err := a.stream.StreamRenameEntry(StreamRenameEntryReqToWire(in))
	if err != nil {
		return nil, err
	}
	return &zapStreamRenameClientStream{s: s}, nil
}

func (a *zapFilerClient) TraverseBfsMetadata(_ context.Context, in *filer_pb.TraverseBfsMetadataRequest, _ ...grpc.CallOption) (grpc.ServerStreamingClient[filer_pb.TraverseBfsMetadataResponse], error) {
	s, err := a.stream.TraverseBfsMetadata(TraverseBfsMetadataReqToWire(in))
	if err != nil {
		return nil, err
	}
	return &zapTraverseBfsClientStream{s: s}, nil
}

func (a *zapFilerClient) SubscribeMetadata(_ context.Context, in *filer_pb.SubscribeMetadataRequest, _ ...grpc.CallOption) (grpc.ServerStreamingClient[filer_pb.SubscribeMetadataResponse], error) {
	s, err := a.stream.SubscribeMetadata(SubscribeMetadataReqToWire(in))
	if err != nil {
		return nil, err
	}
	return &zapSubscribeMetadataClientStream{s: s}, nil
}

func (a *zapFilerClient) SubscribeLocalMetadata(_ context.Context, in *filer_pb.SubscribeMetadataRequest, _ ...grpc.CallOption) (grpc.ServerStreamingClient[filer_pb.SubscribeMetadataResponse], error) {
	s, err := a.stream.SubscribeLocalMetadata(SubscribeMetadataReqToWire(in))
	if err != nil {
		return nil, err
	}
	return &zapSubscribeMetadataClientStream{s: s}, nil
}

func (a *zapFilerClient) StreamMutateEntry(_ context.Context, _ ...grpc.CallOption) (grpc.BidiStreamingClient[filer_pb.StreamMutateEntryRequest, filer_pb.StreamMutateEntryResponse], error) {
	// The bidi stream carries requests via Send; open it with an empty (Which=0,
	// "no member") request frame so the server has a parseable opener that yields
	// no response.
	s, err := a.stream.StreamMutateEntry(filerwire.NewStreamMutateEntryRequest(filerwire.StreamMutateEntryRequestInput{}))
	if err != nil {
		return nil, err
	}
	return &zapStreamMutateClientStream{s: s}, nil
}

// --- grpc.ClientStream stub shared by every adapter stream -----------------

// zapClientStreamStub satisfies the non-Recv/Send half of grpc.ClientStream. The
// consumers only ever call Recv()/Send()/CloseSend(); the remaining methods
// exist solely to satisfy the interface.
type zapClientStreamStub struct{}

func (zapClientStreamStub) Header() (metadata.MD, error) { return nil, nil }
func (zapClientStreamStub) Trailer() metadata.MD         { return nil }
func (zapClientStreamStub) Context() context.Context     { return context.Background() }
func (zapClientStreamStub) SendMsg(any) error            { return nil }
func (zapClientStreamStub) RecvMsg(any) error            { return nil }

// zapListEntriesClientStream adapts a filerstream ListEntries stream to
// grpc.ServerStreamingClient[ListEntriesResponse].
type zapListEntriesClientStream struct {
	zapClientStreamStub
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
	zapClientStreamStub
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
	zapClientStreamStub
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
	zapClientStreamStub
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
	zapClientStreamStub
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
