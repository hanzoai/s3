package mount

// filer_client_adapter.go is the FUSE mount's strangler seam onto the native
// ZAP transport. It implements the filer_pb.HanzoFilerClient interface — the
// contract the entire FUSE client path already speaks — but routes every call
// over a github.com/zap-proto/go transport.Conn instead of gRPC:
//
//   - the 28 unary RPCs go through filerwire.HanzoFilerClient (over the conn's
//     Call channel), translating *filer_pb.<Rpc>Request <-> ZAP buffer <->
//     *filer_pb.<Rpc>Response with the filerzap converters;
//   - the 6 streaming RPCs open a transport stream via the filerstream client
//     and return a grpc.ServerStreamingClient / grpc.BidiStreamingClient whose
//     Recv()/Send() decode/encode each frame with the same filerzap converters.
//
// Because the adapter satisfies filer_pb.HanzoFilerClient, none of the ~25 FUSE
// call sites change: they keep building *filer_pb requests and reading *filer_pb
// responses. The bytes are ZAP at every hop; protobuf framing and gRPC are gone.

import (
	"context"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/hanzoai/s3/s3/filerzap"
	filer_pb "github.com/hanzoai/s3/s3/pb/filer_pb"
	filerstream "github.com/hanzoai/s3/s3/wire/filer/filerstream"
	filerwire "github.com/hanzoai/s3/s3/wire/filer"

	"github.com/zap-proto/go/transport"
)

// filerClientAdapter binds a transport.Conn to the filer_pb.HanzoFilerClient
// contract. unary issues calls over the filerwire client; stream opens streams
// over the filerstream client on the SAME connection.
type filerClientAdapter struct {
	conn   *transport.Conn
	unary  *filerwire.HanzoFilerClient
	stream *filerstream.Client
}

// newFilerClientAdapter wraps an established transport.Conn.
func newFilerClientAdapter(conn *transport.Conn) *filerClientAdapter {
	return &filerClientAdapter{
		conn:   conn,
		unary:  filerwire.NewHanzoFilerClient(conn, nil),
		stream: filerstream.NewClient(conn),
	}
}

var _ filer_pb.HanzoFilerClient = (*filerClientAdapter)(nil)

// --- unary RPCs ------------------------------------------------------------

func (a *filerClientAdapter) LookupDirectoryEntry(_ context.Context, in *filer_pb.LookupDirectoryEntryRequest, _ ...grpc.CallOption) (*filer_pb.LookupDirectoryEntryResponse, error) {
	_, body, err := a.unary.LookupDirectoryEntry(filerzap.LookupDirectoryEntryReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.LookupDirectoryEntryRespFromWire(body)
}

func (a *filerClientAdapter) CreateEntry(_ context.Context, in *filer_pb.CreateEntryRequest, _ ...grpc.CallOption) (*filer_pb.CreateEntryResponse, error) {
	_, body, err := a.unary.CreateEntry(filerzap.CreateEntryReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.CreateEntryRespFromWire(body)
}

func (a *filerClientAdapter) UpdateEntry(_ context.Context, in *filer_pb.UpdateEntryRequest, _ ...grpc.CallOption) (*filer_pb.UpdateEntryResponse, error) {
	_, body, err := a.unary.UpdateEntry(filerzap.UpdateEntryReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.UpdateEntryRespFromWire(body)
}

func (a *filerClientAdapter) TouchAccessTime(_ context.Context, in *filer_pb.TouchAccessTimeRequest, _ ...grpc.CallOption) (*filer_pb.TouchAccessTimeResponse, error) {
	_, body, err := a.unary.TouchAccessTime(filerzap.TouchAccessTimeReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.TouchAccessTimeRespFromWire(body)
}

func (a *filerClientAdapter) AppendToEntry(_ context.Context, in *filer_pb.AppendToEntryRequest, _ ...grpc.CallOption) (*filer_pb.AppendToEntryResponse, error) {
	_, body, err := a.unary.AppendToEntry(filerzap.AppendToEntryReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.AppendToEntryRespFromWire(body)
}

func (a *filerClientAdapter) DeleteEntry(_ context.Context, in *filer_pb.DeleteEntryRequest, _ ...grpc.CallOption) (*filer_pb.DeleteEntryResponse, error) {
	_, body, err := a.unary.DeleteEntry(filerzap.DeleteEntryReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.DeleteEntryRespFromWire(body)
}

func (a *filerClientAdapter) ObjectTransaction(_ context.Context, in *filer_pb.ObjectTransactionRequest, _ ...grpc.CallOption) (*filer_pb.ObjectTransactionResponse, error) {
	_, body, err := a.unary.ObjectTransaction(filerzap.ObjectTransactionReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.ObjectTransactionRespFromWire(body)
}

func (a *filerClientAdapter) ObjectTransactionBatch(_ context.Context, in *filer_pb.ObjectTransactionBatchRequest, _ ...grpc.CallOption) (*filer_pb.ObjectTransactionBatchResponse, error) {
	_, body, err := a.unary.ObjectTransactionBatch(filerzap.ObjectTransactionBatchReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.ObjectTransactionBatchRespFromWire(body)
}

func (a *filerClientAdapter) PosixLock(_ context.Context, in *filer_pb.PosixLockRequest, _ ...grpc.CallOption) (*filer_pb.PosixLockResponse, error) {
	_, body, err := a.unary.PosixLock(filerzap.PosixLockReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.PosixLockRespFromWire(body)
}

func (a *filerClientAdapter) AtomicRenameEntry(_ context.Context, in *filer_pb.AtomicRenameEntryRequest, _ ...grpc.CallOption) (*filer_pb.AtomicRenameEntryResponse, error) {
	_, body, err := a.unary.AtomicRenameEntry(filerzap.AtomicRenameEntryReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.AtomicRenameEntryRespFromWire(body)
}

func (a *filerClientAdapter) AssignVolume(_ context.Context, in *filer_pb.AssignVolumeRequest, _ ...grpc.CallOption) (*filer_pb.AssignVolumeResponse, error) {
	_, body, err := a.unary.AssignVolume(filerzap.AssignVolumeReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.AssignVolumeRespFromWire(body)
}

func (a *filerClientAdapter) LookupVolume(_ context.Context, in *filer_pb.LookupVolumeRequest, _ ...grpc.CallOption) (*filer_pb.LookupVolumeResponse, error) {
	_, body, err := a.unary.LookupVolume(filerzap.LookupVolumeReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.LookupVolumeRespFromWire(body)
}

func (a *filerClientAdapter) CollectionList(_ context.Context, in *filer_pb.CollectionListRequest, _ ...grpc.CallOption) (*filer_pb.CollectionListResponse, error) {
	_, body, err := a.unary.CollectionList(filerzap.CollectionListReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.CollectionListRespFromWire(body)
}

func (a *filerClientAdapter) DeleteCollection(_ context.Context, in *filer_pb.DeleteCollectionRequest, _ ...grpc.CallOption) (*filer_pb.DeleteCollectionResponse, error) {
	_, body, err := a.unary.DeleteCollection(filerzap.DeleteCollectionReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.DeleteCollectionRespFromWire(body)
}

func (a *filerClientAdapter) Statistics(_ context.Context, in *filer_pb.StatisticsRequest, _ ...grpc.CallOption) (*filer_pb.StatisticsResponse, error) {
	_, body, err := a.unary.Statistics(filerzap.StatisticsReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.StatisticsRespFromWire(body)
}

func (a *filerClientAdapter) Ping(_ context.Context, in *filer_pb.PingRequest, _ ...grpc.CallOption) (*filer_pb.PingResponse, error) {
	_, body, err := a.unary.Ping(filerzap.PingReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.PingRespFromWire(body)
}

func (a *filerClientAdapter) GetFilerConfiguration(_ context.Context, in *filer_pb.GetFilerConfigurationRequest, _ ...grpc.CallOption) (*filer_pb.GetFilerConfigurationResponse, error) {
	_, body, err := a.unary.GetFilerConfiguration(filerzap.GetFilerConfigurationReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.GetFilerConfigurationRespFromWire(body)
}

func (a *filerClientAdapter) ListMetadataSubscribers(_ context.Context, in *filer_pb.ListMetadataSubscribersRequest, _ ...grpc.CallOption) (*filer_pb.ListMetadataSubscribersResponse, error) {
	_, body, err := a.unary.ListMetadataSubscribers(filerzap.ListMetadataSubscribersReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.ListMetadataSubscribersRespFromWire(body)
}

func (a *filerClientAdapter) KvGet(_ context.Context, in *filer_pb.KvGetRequest, _ ...grpc.CallOption) (*filer_pb.KvGetResponse, error) {
	_, body, err := a.unary.KvGet(filerzap.KvGetReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.KvGetRespFromWire(body)
}

func (a *filerClientAdapter) KvPut(_ context.Context, in *filer_pb.KvPutRequest, _ ...grpc.CallOption) (*filer_pb.KvPutResponse, error) {
	_, body, err := a.unary.KvPut(filerzap.KvPutReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.KvPutRespFromWire(body)
}

func (a *filerClientAdapter) CacheRemoteObjectToLocalCluster(_ context.Context, in *filer_pb.CacheRemoteObjectToLocalClusterRequest, _ ...grpc.CallOption) (*filer_pb.CacheRemoteObjectToLocalClusterResponse, error) {
	_, body, err := a.unary.CacheRemoteObjectToLocalCluster(filerzap.CacheRemoteObjectToLocalClusterReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.CacheRemoteObjectToLocalClusterRespFromWire(body)
}

func (a *filerClientAdapter) DistributedLock(_ context.Context, in *filer_pb.LockRequest, _ ...grpc.CallOption) (*filer_pb.LockResponse, error) {
	_, body, err := a.unary.DistributedLock(filerzap.DistributedLockReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.DistributedLockRespFromWire(body)
}

func (a *filerClientAdapter) DistributedUnlock(_ context.Context, in *filer_pb.UnlockRequest, _ ...grpc.CallOption) (*filer_pb.UnlockResponse, error) {
	_, body, err := a.unary.DistributedUnlock(filerzap.DistributedUnlockReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.DistributedUnlockRespFromWire(body)
}

func (a *filerClientAdapter) FindLockOwner(_ context.Context, in *filer_pb.FindLockOwnerRequest, _ ...grpc.CallOption) (*filer_pb.FindLockOwnerResponse, error) {
	_, body, err := a.unary.FindLockOwner(filerzap.FindLockOwnerReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.FindLockOwnerRespFromWire(body)
}

func (a *filerClientAdapter) TransferLocks(_ context.Context, in *filer_pb.TransferLocksRequest, _ ...grpc.CallOption) (*filer_pb.TransferLocksResponse, error) {
	_, body, err := a.unary.TransferLocks(filerzap.TransferLocksReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.TransferLocksRespFromWire(body)
}

func (a *filerClientAdapter) ReplicateLock(_ context.Context, in *filer_pb.ReplicateLockRequest, _ ...grpc.CallOption) (*filer_pb.ReplicateLockResponse, error) {
	_, body, err := a.unary.ReplicateLock(filerzap.ReplicateLockReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.ReplicateLockRespFromWire(body)
}

func (a *filerClientAdapter) MountRegister(_ context.Context, in *filer_pb.MountRegisterRequest, _ ...grpc.CallOption) (*filer_pb.MountRegisterResponse, error) {
	_, body, err := a.unary.MountRegister(filerzap.MountRegisterReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.MountRegisterRespFromWire(body)
}

func (a *filerClientAdapter) MountList(_ context.Context, in *filer_pb.MountListRequest, _ ...grpc.CallOption) (*filer_pb.MountListResponse, error) {
	_, body, err := a.unary.MountList(filerzap.MountListReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.MountListRespFromWire(body)
}

// --- streaming RPCs --------------------------------------------------------

func (a *filerClientAdapter) ListEntries(_ context.Context, in *filer_pb.ListEntriesRequest, _ ...grpc.CallOption) (grpc.ServerStreamingClient[filer_pb.ListEntriesResponse], error) {
	s, err := a.stream.ListEntries(filerzap.ListEntriesReqToWire(in))
	if err != nil {
		return nil, err
	}
	return &listEntriesClientStream{s: s}, nil
}

func (a *filerClientAdapter) StreamRenameEntry(_ context.Context, in *filer_pb.StreamRenameEntryRequest, _ ...grpc.CallOption) (grpc.ServerStreamingClient[filer_pb.StreamRenameEntryResponse], error) {
	s, err := a.stream.StreamRenameEntry(filerzap.StreamRenameEntryReqToWire(in))
	if err != nil {
		return nil, err
	}
	return &streamRenameClientStream{s: s}, nil
}

func (a *filerClientAdapter) TraverseBfsMetadata(_ context.Context, in *filer_pb.TraverseBfsMetadataRequest, _ ...grpc.CallOption) (grpc.ServerStreamingClient[filer_pb.TraverseBfsMetadataResponse], error) {
	s, err := a.stream.TraverseBfsMetadata(filerzap.TraverseBfsMetadataReqToWire(in))
	if err != nil {
		return nil, err
	}
	return &traverseBfsClientStream{s: s}, nil
}

func (a *filerClientAdapter) SubscribeMetadata(_ context.Context, in *filer_pb.SubscribeMetadataRequest, _ ...grpc.CallOption) (grpc.ServerStreamingClient[filer_pb.SubscribeMetadataResponse], error) {
	s, err := a.stream.SubscribeMetadata(filerzap.SubscribeMetadataReqToWire(in))
	if err != nil {
		return nil, err
	}
	return &subscribeMetadataClientStream{s: s}, nil
}

func (a *filerClientAdapter) SubscribeLocalMetadata(_ context.Context, in *filer_pb.SubscribeMetadataRequest, _ ...grpc.CallOption) (grpc.ServerStreamingClient[filer_pb.SubscribeMetadataResponse], error) {
	s, err := a.stream.SubscribeLocalMetadata(filerzap.SubscribeMetadataReqToWire(in))
	if err != nil {
		return nil, err
	}
	return &subscribeMetadataClientStream{s: s}, nil
}

func (a *filerClientAdapter) StreamMutateEntry(_ context.Context, _ ...grpc.CallOption) (grpc.BidiStreamingClient[filer_pb.StreamMutateEntryRequest, filer_pb.StreamMutateEntryResponse], error) {
	// The bidi stream carries requests via Send; open it with an empty (Which=0,
	// "no member") request frame so the server has a parseable opener that yields
	// no response.
	s, err := a.stream.StreamMutateEntry(filerwire.NewStreamMutateEntryRequest(filerwire.StreamMutateEntryRequestInput{}))
	if err != nil {
		return nil, err
	}
	return &streamMutateClientStream{s: s}, nil
}

// --- grpc.ClientStream stub shared by every adapter stream -----------------

// clientStreamStub satisfies the non-Recv/Send half of grpc.ClientStream. The
// FUSE consumers only ever call Recv()/Send()/CloseSend(); the remaining
// methods exist solely to satisfy the interface.
type clientStreamStub struct{}

func (clientStreamStub) Header() (metadata.MD, error) { return nil, nil }
func (clientStreamStub) Trailer() metadata.MD         { return nil }
func (clientStreamStub) Context() context.Context     { return context.Background() }
func (clientStreamStub) SendMsg(any) error            { return nil }
func (clientStreamStub) RecvMsg(any) error            { return nil }

// listEntriesClientStream adapts a filerstream ListEntries stream to
// grpc.ServerStreamingClient[ListEntriesResponse].
type listEntriesClientStream struct {
	clientStreamStub
	s *filerstream.ClientListEntriesStream
}

func (x *listEntriesClientStream) Recv() (*filer_pb.ListEntriesResponse, error) {
	b, err := x.s.RecvBytes()
	if err != nil {
		return nil, err
	}
	return filerzap.ListEntriesRespFromWire(b)
}
func (x *listEntriesClientStream) CloseSend() error { return x.s.Close() }

// streamRenameClientStream adapts a filerstream StreamRenameEntry stream.
type streamRenameClientStream struct {
	clientStreamStub
	s *filerstream.ClientStreamRenameEntryStream
}

func (x *streamRenameClientStream) Recv() (*filer_pb.StreamRenameEntryResponse, error) {
	b, err := x.s.RecvBytes()
	if err != nil {
		return nil, err
	}
	return filerzap.StreamRenameEntryRespFromWire(b)
}
func (x *streamRenameClientStream) CloseSend() error { return x.s.Close() }

// traverseBfsClientStream adapts a filerstream TraverseBfsMetadata stream. mount
// does not currently consume this RPC, but the adapter implements it so the full
// filer_pb.HanzoFilerClient contract is honored.
type traverseBfsClientStream struct {
	clientStreamStub
	s *filerstream.ClientTraverseBfsMetadataStream
}

func (x *traverseBfsClientStream) Recv() (*filer_pb.TraverseBfsMetadataResponse, error) {
	// mount never reads TraverseBfsMetadata frames; drain to EOF.
	if _, err := x.s.RecvBytes(); err != nil {
		return nil, err
	}
	return nil, io.EOF
}
func (x *traverseBfsClientStream) CloseSend() error { return x.s.Close() }

// subscribeMetadataClientStream adapts a filerstream SubscribeMetadata stream
// (used for both cluster-wide and local subscriptions).
type subscribeMetadataClientStream struct {
	clientStreamStub
	s *filerstream.ClientSubscribeMetadataStream
}

func (x *subscribeMetadataClientStream) Recv() (*filer_pb.SubscribeMetadataResponse, error) {
	b, err := x.s.RecvBytes()
	if err != nil {
		return nil, err
	}
	return filerzap.SubscribeMetadataResponseFromWire(b)
}
func (x *subscribeMetadataClientStream) CloseSend() error { return x.s.Close() }

// streamMutateClientStream adapts the bidirectional StreamMutateEntry stream.
type streamMutateClientStream struct {
	clientStreamStub
	s *filerstream.ClientStreamMutateEntryStream
}

func (x *streamMutateClientStream) Send(in *filer_pb.StreamMutateEntryRequest) error {
	return x.s.Send(filerzap.StreamMutateEntryReqToWire(in))
}
func (x *streamMutateClientStream) Recv() (*filer_pb.StreamMutateEntryResponse, error) {
	b, err := x.s.RecvBytes()
	if err != nil {
		return nil, err
	}
	return filerzap.StreamMutateEntryRespFromWire(b)
}
func (x *streamMutateClientStream) CloseSend() error { return x.s.CloseSend() }
