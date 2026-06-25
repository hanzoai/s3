// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// client.go is the strangler seam from the volume_server_pb.VolumeServerClient
// contract onto the native ZAP transport. It lives in package volumezap — the
// package that owns the per-RPC <Rpc>ReqToWire/<Rpc>RespFromWire converters — so
// the contract is bridged in exactly ONE place. It implements
// volume_server_pb.VolumeServerClient but routes every call over a
// github.com/zap-proto/go transport.Conn instead of gRPC:
//
//   - the 26 unary RPCs go through the generated vsw.VolumeServerClient (over the
//     conn's Call channel), translating *volume_server_pb.<Rpc>Request <-> ZAP
//     buffer <-> *volume_server_pb.<Rpc>Response with the converters in rpc.go;
//   - the 10 server-streaming RPCs and the 1 client-streaming RPC open a transport
//     stream via the vsw.Client streaming opens (which reuse the SAME conn) and
//     return a grpc.ServerStreamingClient / grpc.ClientStreamingClient whose
//     Recv()/Send() decode/encode each frame with the same converters.
//
// Because the adapter satisfies volume_server_pb.VolumeServerClient, no call site
// changes: callers keep building *volume_server_pb requests and reading
// *volume_server_pb responses. The bytes are ZAP at every hop; protobuf framing
// and gRPC are gone.

package volumezap

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/zap-proto/go/transport"

	volume_server_pb "github.com/hanzoai/s3/s3/pb/volume_server_pb"
	vsw "github.com/hanzoai/s3/s3/wire/volume_server"
)

// zapVolumeClient binds a transport.Conn to the VolumeServerClient contract.
// unary issues calls over the generated wire client; stream opens streams over
// the typed wire Client on the SAME connection.
type zapVolumeClient struct {
	conn   transport.Conn
	unary  *vsw.VolumeServerClient
	stream *vsw.Client
}

// New wraps an established transport.Conn as a volume_server_pb.VolumeServerClient
// that routes over ZAP. capability is the optional capability token attached to
// every unary request (nil for none). The caller owns conn's lifecycle (Close
// when done); the adapter pools nothing.
func New(conn transport.Conn, capability []byte) volume_server_pb.VolumeServerClient {
	return &zapVolumeClient{
		conn:   conn,
		unary:  vsw.NewVolumeServerClient(conn, capability),
		stream: vsw.NewClient(conn),
	}
}

var _ volume_server_pb.VolumeServerClient = (*zapVolumeClient)(nil)

// --- unary RPCs ------------------------------------------------------------

func (a *zapVolumeClient) BatchDelete(_ context.Context, in *volume_server_pb.BatchDeleteRequest, _ ...grpc.CallOption) (*volume_server_pb.BatchDeleteResponse, error) {
	_, body, err := a.unary.BatchDelete(BatchDeleteReqToWire(in))
	if err != nil {
		return nil, err
	}
	return BatchDeleteRespFromWire(body)
}

func (a *zapVolumeClient) VacuumVolumeCheck(_ context.Context, in *volume_server_pb.VacuumVolumeCheckRequest, _ ...grpc.CallOption) (*volume_server_pb.VacuumVolumeCheckResponse, error) {
	_, body, err := a.unary.VacuumVolumeCheck(VacuumVolumeCheckReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VacuumVolumeCheckRespFromWire(body)
}

func (a *zapVolumeClient) VacuumVolumeCommit(_ context.Context, in *volume_server_pb.VacuumVolumeCommitRequest, _ ...grpc.CallOption) (*volume_server_pb.VacuumVolumeCommitResponse, error) {
	_, body, err := a.unary.VacuumVolumeCommit(VacuumVolumeCommitReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VacuumVolumeCommitRespFromWire(body)
}

func (a *zapVolumeClient) VacuumVolumeCleanup(_ context.Context, in *volume_server_pb.VacuumVolumeCleanupRequest, _ ...grpc.CallOption) (*volume_server_pb.VacuumVolumeCleanupResponse, error) {
	_, body, err := a.unary.VacuumVolumeCleanup(VacuumVolumeCleanupReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VacuumVolumeCleanupRespFromWire(body)
}

func (a *zapVolumeClient) DeleteCollection(_ context.Context, in *volume_server_pb.DeleteCollectionRequest, _ ...grpc.CallOption) (*volume_server_pb.DeleteCollectionResponse, error) {
	_, body, err := a.unary.DeleteCollection(DeleteCollectionReqToWire(in))
	if err != nil {
		return nil, err
	}
	return DeleteCollectionRespFromWire(body)
}

func (a *zapVolumeClient) AllocateVolume(_ context.Context, in *volume_server_pb.AllocateVolumeRequest, _ ...grpc.CallOption) (*volume_server_pb.AllocateVolumeResponse, error) {
	_, body, err := a.unary.AllocateVolume(AllocateVolumeReqToWire(in))
	if err != nil {
		return nil, err
	}
	return AllocateVolumeRespFromWire(body)
}

func (a *zapVolumeClient) VolumeSyncStatus(_ context.Context, in *volume_server_pb.VolumeSyncStatusRequest, _ ...grpc.CallOption) (*volume_server_pb.VolumeSyncStatusResponse, error) {
	_, body, err := a.unary.VolumeSyncStatus(VolumeSyncStatusReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeSyncStatusRespFromWire(body)
}

func (a *zapVolumeClient) VolumeMount(_ context.Context, in *volume_server_pb.VolumeMountRequest, _ ...grpc.CallOption) (*volume_server_pb.VolumeMountResponse, error) {
	_, body, err := a.unary.VolumeMount(VolumeMountReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeMountRespFromWire(body)
}

func (a *zapVolumeClient) VolumeUnmount(_ context.Context, in *volume_server_pb.VolumeUnmountRequest, _ ...grpc.CallOption) (*volume_server_pb.VolumeUnmountResponse, error) {
	_, body, err := a.unary.VolumeUnmount(VolumeUnmountReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeUnmountRespFromWire(body)
}

func (a *zapVolumeClient) VolumeDelete(_ context.Context, in *volume_server_pb.VolumeDeleteRequest, _ ...grpc.CallOption) (*volume_server_pb.VolumeDeleteResponse, error) {
	_, body, err := a.unary.VolumeDelete(VolumeDeleteReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeDeleteRespFromWire(body)
}

func (a *zapVolumeClient) VolumeMarkReadonly(_ context.Context, in *volume_server_pb.VolumeMarkReadonlyRequest, _ ...grpc.CallOption) (*volume_server_pb.VolumeMarkReadonlyResponse, error) {
	_, body, err := a.unary.VolumeMarkReadonly(VolumeMarkReadonlyReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeMarkReadonlyRespFromWire(body)
}

func (a *zapVolumeClient) VolumeMarkWritable(_ context.Context, in *volume_server_pb.VolumeMarkWritableRequest, _ ...grpc.CallOption) (*volume_server_pb.VolumeMarkWritableResponse, error) {
	_, body, err := a.unary.VolumeMarkWritable(VolumeMarkWritableReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeMarkWritableRespFromWire(body)
}

func (a *zapVolumeClient) VolumeConfigure(_ context.Context, in *volume_server_pb.VolumeConfigureRequest, _ ...grpc.CallOption) (*volume_server_pb.VolumeConfigureResponse, error) {
	_, body, err := a.unary.VolumeConfigure(VolumeConfigureReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeConfigureRespFromWire(body)
}

func (a *zapVolumeClient) VolumeStatus(_ context.Context, in *volume_server_pb.VolumeStatusRequest, _ ...grpc.CallOption) (*volume_server_pb.VolumeStatusResponse, error) {
	_, body, err := a.unary.VolumeStatus(VolumeStatusReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeStatusRespFromWire(body)
}

func (a *zapVolumeClient) GetState(_ context.Context, in *volume_server_pb.GetStateRequest, _ ...grpc.CallOption) (*volume_server_pb.GetStateResponse, error) {
	_, body, err := a.unary.GetState(GetStateReqToWire(in))
	if err != nil {
		return nil, err
	}
	return GetStateRespFromWire(body)
}

func (a *zapVolumeClient) SetState(_ context.Context, in *volume_server_pb.SetStateRequest, _ ...grpc.CallOption) (*volume_server_pb.SetStateResponse, error) {
	_, body, err := a.unary.SetState(SetStateReqToWire(in))
	if err != nil {
		return nil, err
	}
	return SetStateRespFromWire(body)
}

func (a *zapVolumeClient) ReadVolumeFileStatus(_ context.Context, in *volume_server_pb.ReadVolumeFileStatusRequest, _ ...grpc.CallOption) (*volume_server_pb.ReadVolumeFileStatusResponse, error) {
	_, body, err := a.unary.ReadVolumeFileStatus(ReadVolumeFileStatusReqToWire(in))
	if err != nil {
		return nil, err
	}
	return ReadVolumeFileStatusRespFromWire(body)
}

func (a *zapVolumeClient) ReadNeedleBlob(_ context.Context, in *volume_server_pb.ReadNeedleBlobRequest, _ ...grpc.CallOption) (*volume_server_pb.ReadNeedleBlobResponse, error) {
	_, body, err := a.unary.ReadNeedleBlob(ReadNeedleBlobReqToWire(in))
	if err != nil {
		return nil, err
	}
	return ReadNeedleBlobRespFromWire(body)
}

func (a *zapVolumeClient) ReadNeedleMeta(_ context.Context, in *volume_server_pb.ReadNeedleMetaRequest, _ ...grpc.CallOption) (*volume_server_pb.ReadNeedleMetaResponse, error) {
	_, body, err := a.unary.ReadNeedleMeta(ReadNeedleMetaReqToWire(in))
	if err != nil {
		return nil, err
	}
	return ReadNeedleMetaRespFromWire(body)
}

func (a *zapVolumeClient) WriteNeedleBlob(_ context.Context, in *volume_server_pb.WriteNeedleBlobRequest, _ ...grpc.CallOption) (*volume_server_pb.WriteNeedleBlobResponse, error) {
	_, body, err := a.unary.WriteNeedleBlob(WriteNeedleBlobReqToWire(in))
	if err != nil {
		return nil, err
	}
	return WriteNeedleBlobRespFromWire(body)
}

func (a *zapVolumeClient) VolumeTailReceiver(_ context.Context, in *volume_server_pb.VolumeTailReceiverRequest, _ ...grpc.CallOption) (*volume_server_pb.VolumeTailReceiverResponse, error) {
	_, body, err := a.unary.VolumeTailReceiver(VolumeTailReceiverReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeTailReceiverRespFromWire(body)
}

func (a *zapVolumeClient) VolumeEcShardsGenerate(_ context.Context, in *volume_server_pb.VolumeEcShardsGenerateRequest, _ ...grpc.CallOption) (*volume_server_pb.VolumeEcShardsGenerateResponse, error) {
	_, body, err := a.unary.VolumeEcShardsGenerate(VolumeEcShardsGenerateReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeEcShardsGenerateRespFromWire(body)
}

func (a *zapVolumeClient) VolumeEcShardsRebuild(_ context.Context, in *volume_server_pb.VolumeEcShardsRebuildRequest, _ ...grpc.CallOption) (*volume_server_pb.VolumeEcShardsRebuildResponse, error) {
	_, body, err := a.unary.VolumeEcShardsRebuild(VolumeEcShardsRebuildReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeEcShardsRebuildRespFromWire(body)
}

func (a *zapVolumeClient) VolumeEcShardsCopy(_ context.Context, in *volume_server_pb.VolumeEcShardsCopyRequest, _ ...grpc.CallOption) (*volume_server_pb.VolumeEcShardsCopyResponse, error) {
	_, body, err := a.unary.VolumeEcShardsCopy(VolumeEcShardsCopyReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeEcShardsCopyRespFromWire(body)
}

func (a *zapVolumeClient) VolumeEcShardsDelete(_ context.Context, in *volume_server_pb.VolumeEcShardsDeleteRequest, _ ...grpc.CallOption) (*volume_server_pb.VolumeEcShardsDeleteResponse, error) {
	_, body, err := a.unary.VolumeEcShardsDelete(VolumeEcShardsDeleteReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeEcShardsDeleteRespFromWire(body)
}

func (a *zapVolumeClient) VolumeEcShardsMount(_ context.Context, in *volume_server_pb.VolumeEcShardsMountRequest, _ ...grpc.CallOption) (*volume_server_pb.VolumeEcShardsMountResponse, error) {
	_, body, err := a.unary.VolumeEcShardsMount(VolumeEcShardsMountReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeEcShardsMountRespFromWire(body)
}

func (a *zapVolumeClient) VolumeEcShardsUnmount(_ context.Context, in *volume_server_pb.VolumeEcShardsUnmountRequest, _ ...grpc.CallOption) (*volume_server_pb.VolumeEcShardsUnmountResponse, error) {
	_, body, err := a.unary.VolumeEcShardsUnmount(VolumeEcShardsUnmountReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeEcShardsUnmountRespFromWire(body)
}

func (a *zapVolumeClient) VolumeEcBlobDelete(_ context.Context, in *volume_server_pb.VolumeEcBlobDeleteRequest, _ ...grpc.CallOption) (*volume_server_pb.VolumeEcBlobDeleteResponse, error) {
	_, body, err := a.unary.VolumeEcBlobDelete(VolumeEcBlobDeleteReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeEcBlobDeleteRespFromWire(body)
}

func (a *zapVolumeClient) VolumeEcShardsToVolume(_ context.Context, in *volume_server_pb.VolumeEcShardsToVolumeRequest, _ ...grpc.CallOption) (*volume_server_pb.VolumeEcShardsToVolumeResponse, error) {
	_, body, err := a.unary.VolumeEcShardsToVolume(VolumeEcShardsToVolumeReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeEcShardsToVolumeRespFromWire(body)
}

func (a *zapVolumeClient) VolumeEcShardsInfo(_ context.Context, in *volume_server_pb.VolumeEcShardsInfoRequest, _ ...grpc.CallOption) (*volume_server_pb.VolumeEcShardsInfoResponse, error) {
	_, body, err := a.unary.VolumeEcShardsInfo(VolumeEcShardsInfoReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeEcShardsInfoRespFromWire(body)
}

func (a *zapVolumeClient) VolumeServerStatus(_ context.Context, in *volume_server_pb.VolumeServerStatusRequest, _ ...grpc.CallOption) (*volume_server_pb.VolumeServerStatusResponse, error) {
	_, body, err := a.unary.VolumeServerStatus(VolumeServerStatusReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeServerStatusRespFromWire(body)
}

func (a *zapVolumeClient) VolumeServerLeave(_ context.Context, in *volume_server_pb.VolumeServerLeaveRequest, _ ...grpc.CallOption) (*volume_server_pb.VolumeServerLeaveResponse, error) {
	_, body, err := a.unary.VolumeServerLeave(VolumeServerLeaveReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeServerLeaveRespFromWire(body)
}

func (a *zapVolumeClient) FetchAndWriteNeedle(_ context.Context, in *volume_server_pb.FetchAndWriteNeedleRequest, _ ...grpc.CallOption) (*volume_server_pb.FetchAndWriteNeedleResponse, error) {
	_, body, err := a.unary.FetchAndWriteNeedle(FetchAndWriteNeedleReqToWire(in))
	if err != nil {
		return nil, err
	}
	return FetchAndWriteNeedleRespFromWire(body)
}

func (a *zapVolumeClient) ScrubVolume(_ context.Context, in *volume_server_pb.ScrubVolumeRequest, _ ...grpc.CallOption) (*volume_server_pb.ScrubVolumeResponse, error) {
	_, body, err := a.unary.ScrubVolume(ScrubVolumeReqToWire(in))
	if err != nil {
		return nil, err
	}
	return ScrubVolumeRespFromWire(body)
}

func (a *zapVolumeClient) ScrubEcVolume(_ context.Context, in *volume_server_pb.ScrubEcVolumeRequest, _ ...grpc.CallOption) (*volume_server_pb.ScrubEcVolumeResponse, error) {
	_, body, err := a.unary.ScrubEcVolume(ScrubEcVolumeReqToWire(in))
	if err != nil {
		return nil, err
	}
	return ScrubEcVolumeRespFromWire(body)
}

func (a *zapVolumeClient) VolumeNeedleStatus(_ context.Context, in *volume_server_pb.VolumeNeedleStatusRequest, _ ...grpc.CallOption) (*volume_server_pb.VolumeNeedleStatusResponse, error) {
	_, body, err := a.unary.VolumeNeedleStatus(VolumeNeedleStatusReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeNeedleStatusRespFromWire(body)
}

func (a *zapVolumeClient) Ping(_ context.Context, in *volume_server_pb.PingRequest, _ ...grpc.CallOption) (*volume_server_pb.PingResponse, error) {
	_, body, err := a.unary.Ping(PingReqToWire(in))
	if err != nil {
		return nil, err
	}
	return PingRespFromWire(body)
}

// --- server-streaming RPCs -------------------------------------------------

func (a *zapVolumeClient) VacuumVolumeCompact(_ context.Context, in *volume_server_pb.VacuumVolumeCompactRequest, _ ...grpc.CallOption) (grpc.ServerStreamingClient[volume_server_pb.VacuumVolumeCompactResponse], error) {
	s, err := a.stream.VacuumVolumeCompact(VacuumVolumeCompactReqInput(in))
	if err != nil {
		return nil, err
	}
	return &vacuumVolumeCompactClientStream{s: s}, nil
}

func (a *zapVolumeClient) VolumeIncrementalCopy(_ context.Context, in *volume_server_pb.VolumeIncrementalCopyRequest, _ ...grpc.CallOption) (grpc.ServerStreamingClient[volume_server_pb.VolumeIncrementalCopyResponse], error) {
	s, err := a.stream.VolumeIncrementalCopy(VolumeIncrementalCopyReqInput(in))
	if err != nil {
		return nil, err
	}
	return &volumeIncrementalCopyClientStream{s: s}, nil
}

func (a *zapVolumeClient) VolumeCopy(_ context.Context, in *volume_server_pb.VolumeCopyRequest, _ ...grpc.CallOption) (grpc.ServerStreamingClient[volume_server_pb.VolumeCopyResponse], error) {
	s, err := a.stream.VolumeCopy(VolumeCopyReqInput(in))
	if err != nil {
		return nil, err
	}
	return &volumeCopyClientStream{s: s}, nil
}

func (a *zapVolumeClient) CopyFile(_ context.Context, in *volume_server_pb.CopyFileRequest, _ ...grpc.CallOption) (grpc.ServerStreamingClient[volume_server_pb.CopyFileResponse], error) {
	s, err := a.stream.CopyFile(CopyFileReqInput(in))
	if err != nil {
		return nil, err
	}
	return &copyFileClientStream{s: s}, nil
}

func (a *zapVolumeClient) ReadAllNeedles(_ context.Context, in *volume_server_pb.ReadAllNeedlesRequest, _ ...grpc.CallOption) (grpc.ServerStreamingClient[volume_server_pb.ReadAllNeedlesResponse], error) {
	s, err := a.stream.ReadAllNeedles(ReadAllNeedlesReqInput(in))
	if err != nil {
		return nil, err
	}
	return &readAllNeedlesClientStream{s: s}, nil
}

func (a *zapVolumeClient) VolumeTailSender(_ context.Context, in *volume_server_pb.VolumeTailSenderRequest, _ ...grpc.CallOption) (grpc.ServerStreamingClient[volume_server_pb.VolumeTailSenderResponse], error) {
	s, err := a.stream.VolumeTailSender(VolumeTailSenderReqInput(in))
	if err != nil {
		return nil, err
	}
	return &volumeTailSenderClientStream{s: s}, nil
}

func (a *zapVolumeClient) VolumeEcShardRead(_ context.Context, in *volume_server_pb.VolumeEcShardReadRequest, _ ...grpc.CallOption) (grpc.ServerStreamingClient[volume_server_pb.VolumeEcShardReadResponse], error) {
	s, err := a.stream.VolumeEcShardRead(VolumeEcShardReadReqInput(in))
	if err != nil {
		return nil, err
	}
	return &volumeEcShardReadClientStream{s: s}, nil
}

func (a *zapVolumeClient) VolumeTierMoveDatToRemote(_ context.Context, in *volume_server_pb.VolumeTierMoveDatToRemoteRequest, _ ...grpc.CallOption) (grpc.ServerStreamingClient[volume_server_pb.VolumeTierMoveDatToRemoteResponse], error) {
	s, err := a.stream.VolumeTierMoveDatToRemote(VolumeTierMoveDatToRemoteReqInput(in))
	if err != nil {
		return nil, err
	}
	return &volumeTierMoveDatToRemoteClientStream{s: s}, nil
}

func (a *zapVolumeClient) VolumeTierMoveDatFromRemote(_ context.Context, in *volume_server_pb.VolumeTierMoveDatFromRemoteRequest, _ ...grpc.CallOption) (grpc.ServerStreamingClient[volume_server_pb.VolumeTierMoveDatFromRemoteResponse], error) {
	s, err := a.stream.VolumeTierMoveDatFromRemote(VolumeTierMoveDatFromRemoteReqInput(in))
	if err != nil {
		return nil, err
	}
	return &volumeTierMoveDatFromRemoteClientStream{s: s}, nil
}

func (a *zapVolumeClient) Query(_ context.Context, in *volume_server_pb.QueryRequest, _ ...grpc.CallOption) (grpc.ServerStreamingClient[volume_server_pb.QueriedStripe], error) {
	s, err := a.stream.Query(QueryReqInput(in))
	if err != nil {
		return nil, err
	}
	return &queryClientStream{s: s}, nil
}

// --- client-streaming RPC --------------------------------------------------

func (a *zapVolumeClient) ReceiveFile(_ context.Context, _ ...grpc.CallOption) (grpc.ClientStreamingClient[volume_server_pb.ReceiveFileRequest, volume_server_pb.ReceiveFileResponse], error) {
	// gRPC's ReceiveFile opens with no message; the wire open needs an init
	// frame, so open with an empty (Data=0, oneof-unset) opener the server's
	// Recv loop drains without acting on. Every grpc Send then ships a real
	// info/content frame.
	s, err := a.stream.ReceiveFile(vsw.ReceiveFileRequestInput{})
	if err != nil {
		return nil, err
	}
	return &receiveFileClientStream{s: s}, nil
}

// --- grpc.ClientStream stub shared by every adapter stream -----------------

// clientStreamStub satisfies the non-Recv/Send half of grpc.ClientStream. The
// consumers only ever call Recv()/Send()/CloseSend()/CloseAndRecv(); the rest
// exist solely to satisfy the interface.
type clientStreamStub struct{}

func (clientStreamStub) Header() (metadata.MD, error) { return nil, nil }
func (clientStreamStub) Trailer() metadata.MD         { return nil }
func (clientStreamStub) Context() context.Context     { return context.Background() }
func (clientStreamStub) SendMsg(any) error            { return nil }
func (clientStreamStub) RecvMsg(any) error            { return nil }

type vacuumVolumeCompactClientStream struct {
	clientStreamStub
	s *vsw.VacuumVolumeCompactClientStream
}

func (x *vacuumVolumeCompactClientStream) Recv() (*volume_server_pb.VacuumVolumeCompactResponse, error) {
	v, err := x.s.Recv()
	if err != nil {
		return nil, err
	}
	return VacuumVolumeCompactRespFromView(v), nil
}
func (x *vacuumVolumeCompactClientStream) CloseSend() error { return nil }

type volumeIncrementalCopyClientStream struct {
	clientStreamStub
	s *vsw.VolumeIncrementalCopyClientStream
}

func (x *volumeIncrementalCopyClientStream) Recv() (*volume_server_pb.VolumeIncrementalCopyResponse, error) {
	v, err := x.s.Recv()
	if err != nil {
		return nil, err
	}
	return VolumeIncrementalCopyRespFromView(v), nil
}
func (x *volumeIncrementalCopyClientStream) CloseSend() error { return nil }

type volumeCopyClientStream struct {
	clientStreamStub
	s *vsw.VolumeCopyClientStream
}

func (x *volumeCopyClientStream) Recv() (*volume_server_pb.VolumeCopyResponse, error) {
	v, err := x.s.Recv()
	if err != nil {
		return nil, err
	}
	return VolumeCopyRespFromView(v), nil
}
func (x *volumeCopyClientStream) CloseSend() error { return nil }

type copyFileClientStream struct {
	clientStreamStub
	s *vsw.CopyFileClientStream
}

func (x *copyFileClientStream) Recv() (*volume_server_pb.CopyFileResponse, error) {
	v, err := x.s.Recv()
	if err != nil {
		return nil, err
	}
	return CopyFileRespFromView(v), nil
}
func (x *copyFileClientStream) CloseSend() error { return nil }

type readAllNeedlesClientStream struct {
	clientStreamStub
	s *vsw.ReadAllNeedlesClientStream
}

func (x *readAllNeedlesClientStream) Recv() (*volume_server_pb.ReadAllNeedlesResponse, error) {
	v, err := x.s.Recv()
	if err != nil {
		return nil, err
	}
	return ReadAllNeedlesRespFromView(v), nil
}
func (x *readAllNeedlesClientStream) CloseSend() error { return nil }

type volumeTailSenderClientStream struct {
	clientStreamStub
	s *vsw.VolumeTailSenderClientStream
}

func (x *volumeTailSenderClientStream) Recv() (*volume_server_pb.VolumeTailSenderResponse, error) {
	v, err := x.s.Recv()
	if err != nil {
		return nil, err
	}
	return VolumeTailSenderRespFromView(v), nil
}
func (x *volumeTailSenderClientStream) CloseSend() error { return nil }

type volumeEcShardReadClientStream struct {
	clientStreamStub
	s *vsw.VolumeEcShardReadClientStream
}

func (x *volumeEcShardReadClientStream) Recv() (*volume_server_pb.VolumeEcShardReadResponse, error) {
	v, err := x.s.Recv()
	if err != nil {
		return nil, err
	}
	return VolumeEcShardReadRespFromView(v), nil
}
func (x *volumeEcShardReadClientStream) CloseSend() error { return nil }

type volumeTierMoveDatToRemoteClientStream struct {
	clientStreamStub
	s *vsw.VolumeTierMoveDatToRemoteClientStream
}

func (x *volumeTierMoveDatToRemoteClientStream) Recv() (*volume_server_pb.VolumeTierMoveDatToRemoteResponse, error) {
	v, err := x.s.Recv()
	if err != nil {
		return nil, err
	}
	return VolumeTierMoveDatToRemoteRespFromView(v), nil
}
func (x *volumeTierMoveDatToRemoteClientStream) CloseSend() error { return nil }

type volumeTierMoveDatFromRemoteClientStream struct {
	clientStreamStub
	s *vsw.VolumeTierMoveDatFromRemoteClientStream
}

func (x *volumeTierMoveDatFromRemoteClientStream) Recv() (*volume_server_pb.VolumeTierMoveDatFromRemoteResponse, error) {
	v, err := x.s.Recv()
	if err != nil {
		return nil, err
	}
	return VolumeTierMoveDatFromRemoteRespFromView(v), nil
}
func (x *volumeTierMoveDatFromRemoteClientStream) CloseSend() error { return nil }

type queryClientStream struct {
	clientStreamStub
	s *vsw.QueryClientStream
}

func (x *queryClientStream) Recv() (*volume_server_pb.QueriedStripe, error) {
	v, err := x.s.Recv()
	if err != nil {
		return nil, err
	}
	return QueriedStripeFromView(v), nil
}
func (x *queryClientStream) CloseSend() error { return nil }

// receiveFileClientStream adapts the one client-streaming RPC: each Send ships a
// ReceiveFileRequest frame; CloseAndRecv half-closes and reads the terminal
// ReceiveFileResponse.
type receiveFileClientStream struct {
	clientStreamStub
	s *vsw.ReceiveFileClientStream
}

func (x *receiveFileClientStream) Send(in *volume_server_pb.ReceiveFileRequest) error {
	return x.s.Send(receiveFileReqInput(in))
}
func (x *receiveFileClientStream) CloseSend() error { return x.s.CloseSend() }
func (x *receiveFileClientStream) CloseAndRecv() (*volume_server_pb.ReceiveFileResponse, error) {
	if err := x.s.CloseSend(); err != nil {
		return nil, err
	}
	v, err := x.s.Reply()
	if err != nil {
		return nil, err
	}
	return ReceiveFileRespFromView(v), nil
}
