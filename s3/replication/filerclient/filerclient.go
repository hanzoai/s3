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
//     and return a grpc.ServerStreamingClient / grpc.BidiStreamingClient whose
//     Recv()/Send() decode/encode each frame with the same filerzap converters.
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

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/hanzoai/s3/s3/filerzap"
	filer_pb "github.com/hanzoai/s3/s3/pb/filer_pb"
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

func (a *adapter) LookupDirectoryEntry(_ context.Context, in *filer_pb.LookupDirectoryEntryRequest, _ ...grpc.CallOption) (*filer_pb.LookupDirectoryEntryResponse, error) {
	_, body, err := a.unary.LookupDirectoryEntry(filerzap.LookupDirectoryEntryReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.LookupDirectoryEntryRespFromWire(body)
}

func (a *adapter) CreateEntry(_ context.Context, in *filer_pb.CreateEntryRequest, _ ...grpc.CallOption) (*filer_pb.CreateEntryResponse, error) {
	_, body, err := a.unary.CreateEntry(filerzap.CreateEntryReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.CreateEntryRespFromWire(body)
}

func (a *adapter) UpdateEntry(_ context.Context, in *filer_pb.UpdateEntryRequest, _ ...grpc.CallOption) (*filer_pb.UpdateEntryResponse, error) {
	_, body, err := a.unary.UpdateEntry(filerzap.UpdateEntryReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.UpdateEntryRespFromWire(body)
}

func (a *adapter) TouchAccessTime(_ context.Context, in *filer_pb.TouchAccessTimeRequest, _ ...grpc.CallOption) (*filer_pb.TouchAccessTimeResponse, error) {
	_, body, err := a.unary.TouchAccessTime(filerzap.TouchAccessTimeReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.TouchAccessTimeRespFromWire(body)
}

func (a *adapter) AppendToEntry(_ context.Context, in *filer_pb.AppendToEntryRequest, _ ...grpc.CallOption) (*filer_pb.AppendToEntryResponse, error) {
	_, body, err := a.unary.AppendToEntry(filerzap.AppendToEntryReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.AppendToEntryRespFromWire(body)
}

func (a *adapter) DeleteEntry(_ context.Context, in *filer_pb.DeleteEntryRequest, _ ...grpc.CallOption) (*filer_pb.DeleteEntryResponse, error) {
	_, body, err := a.unary.DeleteEntry(filerzap.DeleteEntryReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.DeleteEntryRespFromWire(body)
}

func (a *adapter) ObjectTransaction(_ context.Context, in *filer_pb.ObjectTransactionRequest, _ ...grpc.CallOption) (*filer_pb.ObjectTransactionResponse, error) {
	_, body, err := a.unary.ObjectTransaction(filerzap.ObjectTransactionReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.ObjectTransactionRespFromWire(body)
}

func (a *adapter) ObjectTransactionBatch(_ context.Context, in *filer_pb.ObjectTransactionBatchRequest, _ ...grpc.CallOption) (*filer_pb.ObjectTransactionBatchResponse, error) {
	_, body, err := a.unary.ObjectTransactionBatch(filerzap.ObjectTransactionBatchReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.ObjectTransactionBatchRespFromWire(body)
}

func (a *adapter) PosixLock(_ context.Context, in *filer_pb.PosixLockRequest, _ ...grpc.CallOption) (*filer_pb.PosixLockResponse, error) {
	_, body, err := a.unary.PosixLock(filerzap.PosixLockReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.PosixLockRespFromWire(body)
}

func (a *adapter) AtomicRenameEntry(_ context.Context, in *filer_pb.AtomicRenameEntryRequest, _ ...grpc.CallOption) (*filer_pb.AtomicRenameEntryResponse, error) {
	_, body, err := a.unary.AtomicRenameEntry(filerzap.AtomicRenameEntryReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.AtomicRenameEntryRespFromWire(body)
}

func (a *adapter) AssignVolume(_ context.Context, in *filer_pb.AssignVolumeRequest, _ ...grpc.CallOption) (*filer_pb.AssignVolumeResponse, error) {
	_, body, err := a.unary.AssignVolume(filerzap.AssignVolumeReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.AssignVolumeRespFromWire(body)
}

func (a *adapter) LookupVolume(_ context.Context, in *filer_pb.LookupVolumeRequest, _ ...grpc.CallOption) (*filer_pb.LookupVolumeResponse, error) {
	_, body, err := a.unary.LookupVolume(filerzap.LookupVolumeReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.LookupVolumeRespFromWire(body)
}

func (a *adapter) CollectionList(_ context.Context, in *filer_pb.CollectionListRequest, _ ...grpc.CallOption) (*filer_pb.CollectionListResponse, error) {
	_, body, err := a.unary.CollectionList(filerzap.CollectionListReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.CollectionListRespFromWire(body)
}

func (a *adapter) DeleteCollection(_ context.Context, in *filer_pb.DeleteCollectionRequest, _ ...grpc.CallOption) (*filer_pb.DeleteCollectionResponse, error) {
	_, body, err := a.unary.DeleteCollection(filerzap.DeleteCollectionReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.DeleteCollectionRespFromWire(body)
}

func (a *adapter) Statistics(_ context.Context, in *filer_pb.StatisticsRequest, _ ...grpc.CallOption) (*filer_pb.StatisticsResponse, error) {
	_, body, err := a.unary.Statistics(filerzap.StatisticsReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.StatisticsRespFromWire(body)
}

func (a *adapter) Ping(_ context.Context, in *filer_pb.PingRequest, _ ...grpc.CallOption) (*filer_pb.PingResponse, error) {
	_, body, err := a.unary.Ping(filerzap.PingReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.PingRespFromWire(body)
}

func (a *adapter) GetFilerConfiguration(_ context.Context, in *filer_pb.GetFilerConfigurationRequest, _ ...grpc.CallOption) (*filer_pb.GetFilerConfigurationResponse, error) {
	_, body, err := a.unary.GetFilerConfiguration(filerzap.GetFilerConfigurationReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.GetFilerConfigurationRespFromWire(body)
}

func (a *adapter) ListMetadataSubscribers(_ context.Context, in *filer_pb.ListMetadataSubscribersRequest, _ ...grpc.CallOption) (*filer_pb.ListMetadataSubscribersResponse, error) {
	_, body, err := a.unary.ListMetadataSubscribers(filerzap.ListMetadataSubscribersReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.ListMetadataSubscribersRespFromWire(body)
}

func (a *adapter) KvGet(_ context.Context, in *filer_pb.KvGetRequest, _ ...grpc.CallOption) (*filer_pb.KvGetResponse, error) {
	_, body, err := a.unary.KvGet(filerzap.KvGetReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.KvGetRespFromWire(body)
}

func (a *adapter) KvPut(_ context.Context, in *filer_pb.KvPutRequest, _ ...grpc.CallOption) (*filer_pb.KvPutResponse, error) {
	_, body, err := a.unary.KvPut(filerzap.KvPutReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.KvPutRespFromWire(body)
}

func (a *adapter) CacheRemoteObjectToLocalCluster(_ context.Context, in *filer_pb.CacheRemoteObjectToLocalClusterRequest, _ ...grpc.CallOption) (*filer_pb.CacheRemoteObjectToLocalClusterResponse, error) {
	_, body, err := a.unary.CacheRemoteObjectToLocalCluster(filerzap.CacheRemoteObjectToLocalClusterReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.CacheRemoteObjectToLocalClusterRespFromWire(body)
}

func (a *adapter) DistributedLock(_ context.Context, in *filer_pb.LockRequest, _ ...grpc.CallOption) (*filer_pb.LockResponse, error) {
	_, body, err := a.unary.DistributedLock(filerzap.DistributedLockReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.DistributedLockRespFromWire(body)
}

func (a *adapter) DistributedUnlock(_ context.Context, in *filer_pb.UnlockRequest, _ ...grpc.CallOption) (*filer_pb.UnlockResponse, error) {
	_, body, err := a.unary.DistributedUnlock(filerzap.DistributedUnlockReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.DistributedUnlockRespFromWire(body)
}

func (a *adapter) FindLockOwner(_ context.Context, in *filer_pb.FindLockOwnerRequest, _ ...grpc.CallOption) (*filer_pb.FindLockOwnerResponse, error) {
	_, body, err := a.unary.FindLockOwner(filerzap.FindLockOwnerReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.FindLockOwnerRespFromWire(body)
}

func (a *adapter) TransferLocks(_ context.Context, in *filer_pb.TransferLocksRequest, _ ...grpc.CallOption) (*filer_pb.TransferLocksResponse, error) {
	_, body, err := a.unary.TransferLocks(filerzap.TransferLocksReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.TransferLocksRespFromWire(body)
}

func (a *adapter) ReplicateLock(_ context.Context, in *filer_pb.ReplicateLockRequest, _ ...grpc.CallOption) (*filer_pb.ReplicateLockResponse, error) {
	_, body, err := a.unary.ReplicateLock(filerzap.ReplicateLockReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.ReplicateLockRespFromWire(body)
}

func (a *adapter) MountRegister(_ context.Context, in *filer_pb.MountRegisterRequest, _ ...grpc.CallOption) (*filer_pb.MountRegisterResponse, error) {
	_, body, err := a.unary.MountRegister(filerzap.MountRegisterReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.MountRegisterRespFromWire(body)
}

func (a *adapter) MountList(_ context.Context, in *filer_pb.MountListRequest, _ ...grpc.CallOption) (*filer_pb.MountListResponse, error) {
	_, body, err := a.unary.MountList(filerzap.MountListReqToWire(in))
	if err != nil {
		return nil, err
	}
	return filerzap.MountListRespFromWire(body)
}

// --- streaming RPCs --------------------------------------------------------

func (a *adapter) ListEntries(_ context.Context, in *filer_pb.ListEntriesRequest, _ ...grpc.CallOption) (grpc.ServerStreamingClient[filer_pb.ListEntriesResponse], error) {
	s, err := a.stream.ListEntries(filerzap.ListEntriesReqToWire(in))
	if err != nil {
		return nil, err
	}
	return &listEntriesStream{s: s}, nil
}

func (a *adapter) StreamRenameEntry(_ context.Context, in *filer_pb.StreamRenameEntryRequest, _ ...grpc.CallOption) (grpc.ServerStreamingClient[filer_pb.StreamRenameEntryResponse], error) {
	s, err := a.stream.StreamRenameEntry(filerzap.StreamRenameEntryReqToWire(in))
	if err != nil {
		return nil, err
	}
	return &streamRenameStream{s: s}, nil
}

func (a *adapter) TraverseBfsMetadata(_ context.Context, in *filer_pb.TraverseBfsMetadataRequest, _ ...grpc.CallOption) (grpc.ServerStreamingClient[filer_pb.TraverseBfsMetadataResponse], error) {
	s, err := a.stream.TraverseBfsMetadata(filerzap.TraverseBfsMetadataReqToWire(in))
	if err != nil {
		return nil, err
	}
	return &traverseBfsStream{s: s}, nil
}

func (a *adapter) SubscribeMetadata(_ context.Context, in *filer_pb.SubscribeMetadataRequest, _ ...grpc.CallOption) (grpc.ServerStreamingClient[filer_pb.SubscribeMetadataResponse], error) {
	s, err := a.stream.SubscribeMetadata(filerzap.SubscribeMetadataReqToWire(in))
	if err != nil {
		return nil, err
	}
	return &subscribeMetadataStream{s: s}, nil
}

func (a *adapter) SubscribeLocalMetadata(_ context.Context, in *filer_pb.SubscribeMetadataRequest, _ ...grpc.CallOption) (grpc.ServerStreamingClient[filer_pb.SubscribeMetadataResponse], error) {
	s, err := a.stream.SubscribeLocalMetadata(filerzap.SubscribeMetadataReqToWire(in))
	if err != nil {
		return nil, err
	}
	return &subscribeMetadataStream{s: s}, nil
}

func (a *adapter) StreamMutateEntry(_ context.Context, _ ...grpc.CallOption) (grpc.BidiStreamingClient[filer_pb.StreamMutateEntryRequest, filer_pb.StreamMutateEntryResponse], error) {
	// Open the bidi stream with an empty (Which=0, "no member") request frame so
	// the server has a parseable opener that yields no response; requests then
	// flow via Send.
	s, err := a.stream.StreamMutateEntry(filerwire.NewStreamMutateEntryRequest(filerwire.StreamMutateEntryRequestInput{}))
	if err != nil {
		return nil, err
	}
	return &streamMutateStream{s: s}, nil
}

// --- grpc.ClientStream stub shared by every adapter stream -----------------

// clientStreamStub satisfies the non-Recv/Send half of grpc.ClientStream. The
// consumers only ever call Recv()/Send()/CloseSend(); the remaining methods
// exist solely to satisfy the interface.
type clientStreamStub struct{}

func (clientStreamStub) Header() (metadata.MD, error) { return nil, nil }
func (clientStreamStub) Trailer() metadata.MD         { return nil }
func (clientStreamStub) Context() context.Context     { return context.Background() }
func (clientStreamStub) SendMsg(any) error            { return nil }
func (clientStreamStub) RecvMsg(any) error            { return nil }

type listEntriesStream struct {
	clientStreamStub
	s *filerstream.ClientListEntriesStream
}

func (x *listEntriesStream) Recv() (*filer_pb.ListEntriesResponse, error) {
	b, err := x.s.RecvBytes()
	if err != nil {
		return nil, err
	}
	return filerzap.ListEntriesRespFromWire(b)
}
func (x *listEntriesStream) CloseSend() error { return x.s.Close() }

type streamRenameStream struct {
	clientStreamStub
	s *filerstream.ClientStreamRenameEntryStream
}

func (x *streamRenameStream) Recv() (*filer_pb.StreamRenameEntryResponse, error) {
	b, err := x.s.RecvBytes()
	if err != nil {
		return nil, err
	}
	return filerzap.StreamRenameEntryRespFromWire(b)
}
func (x *streamRenameStream) CloseSend() error { return x.s.Close() }

type traverseBfsStream struct {
	clientStreamStub
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
	clientStreamStub
	s *filerstream.ClientSubscribeMetadataStream
}

func (x *subscribeMetadataStream) Recv() (*filer_pb.SubscribeMetadataResponse, error) {
	b, err := x.s.RecvBytes()
	if err != nil {
		return nil, err
	}
	return filerzap.SubscribeMetadataResponseFromWire(b)
}
func (x *subscribeMetadataStream) CloseSend() error { return x.s.Close() }

type streamMutateStream struct {
	clientStreamStub
	s *filerstream.ClientStreamMutateEntryStream
}

func (x *streamMutateStream) Send(in *filer_pb.StreamMutateEntryRequest) error {
	return x.s.Send(filerzap.StreamMutateEntryReqToWire(in))
}
func (x *streamMutateStream) Recv() (*filer_pb.StreamMutateEntryResponse, error) {
	b, err := x.s.RecvBytes()
	if err != nil {
		return nil, err
	}
	return filerzap.StreamMutateEntryRespFromWire(b)
}
func (x *streamMutateStream) CloseSend() error { return x.s.CloseSend() }
