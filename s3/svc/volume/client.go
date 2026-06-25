// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// client.go is the strangler seam from the volume_server_pb.VolumeServerClient
// contract onto the native ZAP transport. It lives in package volume — the
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
//     return an rpc.ServerStream / rpc.ClientStream whose Recv()/Send()
//     decode/encode each frame with the same converters.
//
// Because the adapter satisfies volume_server_pb.VolumeServerClient, no call site
// changes: callers keep building *volume_server_pb requests and reading
// *volume_server_pb responses. The bytes are ZAP at every hop; protobuf framing
// and gRPC are gone.

package volume

import (
	"context"

	"github.com/zap-proto/go/transport"

	"github.com/hanzoai/s3/s3/pb/rpc"
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

func (a *zapVolumeClient) BatchDelete(ctx context.Context, in *volume_server_pb.BatchDeleteRequest) (*volume_server_pb.BatchDeleteResponse, error) {
	_, body, err := a.unary.BatchDelete(ctx, BatchDeleteReqToWire(in))
	if err != nil {
		return nil, err
	}
	return BatchDeleteRespFromWire(body)
}

func (a *zapVolumeClient) VacuumVolumeCheck(ctx context.Context, in *volume_server_pb.VacuumVolumeCheckRequest) (*volume_server_pb.VacuumVolumeCheckResponse, error) {
	_, body, err := a.unary.VacuumVolumeCheck(ctx, VacuumVolumeCheckReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VacuumVolumeCheckRespFromWire(body)
}

func (a *zapVolumeClient) VacuumVolumeCommit(ctx context.Context, in *volume_server_pb.VacuumVolumeCommitRequest) (*volume_server_pb.VacuumVolumeCommitResponse, error) {
	_, body, err := a.unary.VacuumVolumeCommit(ctx, VacuumVolumeCommitReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VacuumVolumeCommitRespFromWire(body)
}

func (a *zapVolumeClient) VacuumVolumeCleanup(ctx context.Context, in *volume_server_pb.VacuumVolumeCleanupRequest) (*volume_server_pb.VacuumVolumeCleanupResponse, error) {
	_, body, err := a.unary.VacuumVolumeCleanup(ctx, VacuumVolumeCleanupReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VacuumVolumeCleanupRespFromWire(body)
}

func (a *zapVolumeClient) DeleteCollection(ctx context.Context, in *volume_server_pb.DeleteCollectionRequest) (*volume_server_pb.DeleteCollectionResponse, error) {
	_, body, err := a.unary.DeleteCollection(ctx, DeleteCollectionReqToWire(in))
	if err != nil {
		return nil, err
	}
	return DeleteCollectionRespFromWire(body)
}

func (a *zapVolumeClient) AllocateVolume(ctx context.Context, in *volume_server_pb.AllocateVolumeRequest) (*volume_server_pb.AllocateVolumeResponse, error) {
	_, body, err := a.unary.AllocateVolume(ctx, AllocateVolumeReqToWire(in))
	if err != nil {
		return nil, err
	}
	return AllocateVolumeRespFromWire(body)
}

func (a *zapVolumeClient) VolumeSyncStatus(ctx context.Context, in *volume_server_pb.VolumeSyncStatusRequest) (*volume_server_pb.VolumeSyncStatusResponse, error) {
	_, body, err := a.unary.VolumeSyncStatus(ctx, VolumeSyncStatusReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeSyncStatusRespFromWire(body)
}

func (a *zapVolumeClient) VolumeMount(ctx context.Context, in *volume_server_pb.VolumeMountRequest) (*volume_server_pb.VolumeMountResponse, error) {
	_, body, err := a.unary.VolumeMount(ctx, VolumeMountReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeMountRespFromWire(body)
}

func (a *zapVolumeClient) VolumeUnmount(ctx context.Context, in *volume_server_pb.VolumeUnmountRequest) (*volume_server_pb.VolumeUnmountResponse, error) {
	_, body, err := a.unary.VolumeUnmount(ctx, VolumeUnmountReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeUnmountRespFromWire(body)
}

func (a *zapVolumeClient) VolumeDelete(ctx context.Context, in *volume_server_pb.VolumeDeleteRequest) (*volume_server_pb.VolumeDeleteResponse, error) {
	_, body, err := a.unary.VolumeDelete(ctx, VolumeDeleteReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeDeleteRespFromWire(body)
}

func (a *zapVolumeClient) VolumeMarkReadonly(ctx context.Context, in *volume_server_pb.VolumeMarkReadonlyRequest) (*volume_server_pb.VolumeMarkReadonlyResponse, error) {
	_, body, err := a.unary.VolumeMarkReadonly(ctx, VolumeMarkReadonlyReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeMarkReadonlyRespFromWire(body)
}

func (a *zapVolumeClient) VolumeMarkWritable(ctx context.Context, in *volume_server_pb.VolumeMarkWritableRequest) (*volume_server_pb.VolumeMarkWritableResponse, error) {
	_, body, err := a.unary.VolumeMarkWritable(ctx, VolumeMarkWritableReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeMarkWritableRespFromWire(body)
}

func (a *zapVolumeClient) VolumeConfigure(ctx context.Context, in *volume_server_pb.VolumeConfigureRequest) (*volume_server_pb.VolumeConfigureResponse, error) {
	_, body, err := a.unary.VolumeConfigure(ctx, VolumeConfigureReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeConfigureRespFromWire(body)
}

func (a *zapVolumeClient) VolumeStatus(ctx context.Context, in *volume_server_pb.VolumeStatusRequest) (*volume_server_pb.VolumeStatusResponse, error) {
	_, body, err := a.unary.VolumeStatus(ctx, VolumeStatusReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeStatusRespFromWire(body)
}

func (a *zapVolumeClient) GetState(ctx context.Context, in *volume_server_pb.GetStateRequest) (*volume_server_pb.GetStateResponse, error) {
	_, body, err := a.unary.GetState(ctx, GetStateReqToWire(in))
	if err != nil {
		return nil, err
	}
	return GetStateRespFromWire(body)
}

func (a *zapVolumeClient) SetState(ctx context.Context, in *volume_server_pb.SetStateRequest) (*volume_server_pb.SetStateResponse, error) {
	_, body, err := a.unary.SetState(ctx, SetStateReqToWire(in))
	if err != nil {
		return nil, err
	}
	return SetStateRespFromWire(body)
}

func (a *zapVolumeClient) ReadVolumeFileStatus(ctx context.Context, in *volume_server_pb.ReadVolumeFileStatusRequest) (*volume_server_pb.ReadVolumeFileStatusResponse, error) {
	_, body, err := a.unary.ReadVolumeFileStatus(ctx, ReadVolumeFileStatusReqToWire(in))
	if err != nil {
		return nil, err
	}
	return ReadVolumeFileStatusRespFromWire(body)
}

func (a *zapVolumeClient) ReadNeedleBlob(ctx context.Context, in *volume_server_pb.ReadNeedleBlobRequest) (*volume_server_pb.ReadNeedleBlobResponse, error) {
	_, body, err := a.unary.ReadNeedleBlob(ctx, ReadNeedleBlobReqToWire(in))
	if err != nil {
		return nil, err
	}
	return ReadNeedleBlobRespFromWire(body)
}

func (a *zapVolumeClient) ReadNeedleMeta(ctx context.Context, in *volume_server_pb.ReadNeedleMetaRequest) (*volume_server_pb.ReadNeedleMetaResponse, error) {
	_, body, err := a.unary.ReadNeedleMeta(ctx, ReadNeedleMetaReqToWire(in))
	if err != nil {
		return nil, err
	}
	return ReadNeedleMetaRespFromWire(body)
}

func (a *zapVolumeClient) WriteNeedleBlob(ctx context.Context, in *volume_server_pb.WriteNeedleBlobRequest) (*volume_server_pb.WriteNeedleBlobResponse, error) {
	_, body, err := a.unary.WriteNeedleBlob(ctx, WriteNeedleBlobReqToWire(in))
	if err != nil {
		return nil, err
	}
	return WriteNeedleBlobRespFromWire(body)
}

func (a *zapVolumeClient) VolumeTailReceiver(ctx context.Context, in *volume_server_pb.VolumeTailReceiverRequest) (*volume_server_pb.VolumeTailReceiverResponse, error) {
	_, body, err := a.unary.VolumeTailReceiver(ctx, VolumeTailReceiverReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeTailReceiverRespFromWire(body)
}

func (a *zapVolumeClient) VolumeEcShardsGenerate(ctx context.Context, in *volume_server_pb.VolumeEcShardsGenerateRequest) (*volume_server_pb.VolumeEcShardsGenerateResponse, error) {
	_, body, err := a.unary.VolumeEcShardsGenerate(ctx, VolumeEcShardsGenerateReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeEcShardsGenerateRespFromWire(body)
}

func (a *zapVolumeClient) VolumeEcShardsRebuild(ctx context.Context, in *volume_server_pb.VolumeEcShardsRebuildRequest) (*volume_server_pb.VolumeEcShardsRebuildResponse, error) {
	_, body, err := a.unary.VolumeEcShardsRebuild(ctx, VolumeEcShardsRebuildReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeEcShardsRebuildRespFromWire(body)
}

func (a *zapVolumeClient) VolumeEcShardsCopy(ctx context.Context, in *volume_server_pb.VolumeEcShardsCopyRequest) (*volume_server_pb.VolumeEcShardsCopyResponse, error) {
	_, body, err := a.unary.VolumeEcShardsCopy(ctx, VolumeEcShardsCopyReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeEcShardsCopyRespFromWire(body)
}

func (a *zapVolumeClient) VolumeEcShardsDelete(ctx context.Context, in *volume_server_pb.VolumeEcShardsDeleteRequest) (*volume_server_pb.VolumeEcShardsDeleteResponse, error) {
	_, body, err := a.unary.VolumeEcShardsDelete(ctx, VolumeEcShardsDeleteReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeEcShardsDeleteRespFromWire(body)
}

func (a *zapVolumeClient) VolumeEcShardsMount(ctx context.Context, in *volume_server_pb.VolumeEcShardsMountRequest) (*volume_server_pb.VolumeEcShardsMountResponse, error) {
	_, body, err := a.unary.VolumeEcShardsMount(ctx, VolumeEcShardsMountReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeEcShardsMountRespFromWire(body)
}

func (a *zapVolumeClient) VolumeEcShardsUnmount(ctx context.Context, in *volume_server_pb.VolumeEcShardsUnmountRequest) (*volume_server_pb.VolumeEcShardsUnmountResponse, error) {
	_, body, err := a.unary.VolumeEcShardsUnmount(ctx, VolumeEcShardsUnmountReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeEcShardsUnmountRespFromWire(body)
}

func (a *zapVolumeClient) VolumeEcBlobDelete(ctx context.Context, in *volume_server_pb.VolumeEcBlobDeleteRequest) (*volume_server_pb.VolumeEcBlobDeleteResponse, error) {
	_, body, err := a.unary.VolumeEcBlobDelete(ctx, VolumeEcBlobDeleteReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeEcBlobDeleteRespFromWire(body)
}

func (a *zapVolumeClient) VolumeEcShardsToVolume(ctx context.Context, in *volume_server_pb.VolumeEcShardsToVolumeRequest) (*volume_server_pb.VolumeEcShardsToVolumeResponse, error) {
	_, body, err := a.unary.VolumeEcShardsToVolume(ctx, VolumeEcShardsToVolumeReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeEcShardsToVolumeRespFromWire(body)
}

func (a *zapVolumeClient) VolumeEcShardsInfo(ctx context.Context, in *volume_server_pb.VolumeEcShardsInfoRequest) (*volume_server_pb.VolumeEcShardsInfoResponse, error) {
	_, body, err := a.unary.VolumeEcShardsInfo(ctx, VolumeEcShardsInfoReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeEcShardsInfoRespFromWire(body)
}

func (a *zapVolumeClient) VolumeServerStatus(ctx context.Context, in *volume_server_pb.VolumeServerStatusRequest) (*volume_server_pb.VolumeServerStatusResponse, error) {
	_, body, err := a.unary.VolumeServerStatus(ctx, VolumeServerStatusReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeServerStatusRespFromWire(body)
}

func (a *zapVolumeClient) VolumeServerLeave(ctx context.Context, in *volume_server_pb.VolumeServerLeaveRequest) (*volume_server_pb.VolumeServerLeaveResponse, error) {
	_, body, err := a.unary.VolumeServerLeave(ctx, VolumeServerLeaveReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeServerLeaveRespFromWire(body)
}

func (a *zapVolumeClient) FetchAndWriteNeedle(ctx context.Context, in *volume_server_pb.FetchAndWriteNeedleRequest) (*volume_server_pb.FetchAndWriteNeedleResponse, error) {
	_, body, err := a.unary.FetchAndWriteNeedle(ctx, FetchAndWriteNeedleReqToWire(in))
	if err != nil {
		return nil, err
	}
	return FetchAndWriteNeedleRespFromWire(body)
}

func (a *zapVolumeClient) ScrubVolume(ctx context.Context, in *volume_server_pb.ScrubVolumeRequest) (*volume_server_pb.ScrubVolumeResponse, error) {
	_, body, err := a.unary.ScrubVolume(ctx, ScrubVolumeReqToWire(in))
	if err != nil {
		return nil, err
	}
	return ScrubVolumeRespFromWire(body)
}

func (a *zapVolumeClient) ScrubEcVolume(ctx context.Context, in *volume_server_pb.ScrubEcVolumeRequest) (*volume_server_pb.ScrubEcVolumeResponse, error) {
	_, body, err := a.unary.ScrubEcVolume(ctx, ScrubEcVolumeReqToWire(in))
	if err != nil {
		return nil, err
	}
	return ScrubEcVolumeRespFromWire(body)
}

func (a *zapVolumeClient) VolumeNeedleStatus(ctx context.Context, in *volume_server_pb.VolumeNeedleStatusRequest) (*volume_server_pb.VolumeNeedleStatusResponse, error) {
	_, body, err := a.unary.VolumeNeedleStatus(ctx, VolumeNeedleStatusReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeNeedleStatusRespFromWire(body)
}

func (a *zapVolumeClient) Ping(ctx context.Context, in *volume_server_pb.PingRequest) (*volume_server_pb.PingResponse, error) {
	_, body, err := a.unary.Ping(ctx, PingReqToWire(in))
	if err != nil {
		return nil, err
	}
	return PingRespFromWire(body)
}

// --- server-streaming RPCs -------------------------------------------------

func (a *zapVolumeClient) VacuumVolumeCompact(_ context.Context, in *volume_server_pb.VacuumVolumeCompactRequest) (rpc.ServerStream[volume_server_pb.VacuumVolumeCompactResponse], error) {
	s, err := a.stream.VacuumVolumeCompact(VacuumVolumeCompactReqInput(in))
	if err != nil {
		return nil, err
	}
	return &vacuumVolumeCompactClientStream{s: s}, nil
}

func (a *zapVolumeClient) VolumeIncrementalCopy(_ context.Context, in *volume_server_pb.VolumeIncrementalCopyRequest) (rpc.ServerStream[volume_server_pb.VolumeIncrementalCopyResponse], error) {
	s, err := a.stream.VolumeIncrementalCopy(VolumeIncrementalCopyReqInput(in))
	if err != nil {
		return nil, err
	}
	return &volumeIncrementalCopyClientStream{s: s}, nil
}

func (a *zapVolumeClient) VolumeCopy(_ context.Context, in *volume_server_pb.VolumeCopyRequest) (rpc.ServerStream[volume_server_pb.VolumeCopyResponse], error) {
	s, err := a.stream.VolumeCopy(VolumeCopyReqInput(in))
	if err != nil {
		return nil, err
	}
	return &volumeCopyClientStream{s: s}, nil
}

func (a *zapVolumeClient) CopyFile(_ context.Context, in *volume_server_pb.CopyFileRequest) (rpc.ServerStream[volume_server_pb.CopyFileResponse], error) {
	s, err := a.stream.CopyFile(CopyFileReqInput(in))
	if err != nil {
		return nil, err
	}
	return &copyFileClientStream{s: s}, nil
}

func (a *zapVolumeClient) ReadAllNeedles(_ context.Context, in *volume_server_pb.ReadAllNeedlesRequest) (rpc.ServerStream[volume_server_pb.ReadAllNeedlesResponse], error) {
	s, err := a.stream.ReadAllNeedles(ReadAllNeedlesReqInput(in))
	if err != nil {
		return nil, err
	}
	return &readAllNeedlesClientStream{s: s}, nil
}

func (a *zapVolumeClient) VolumeTailSender(_ context.Context, in *volume_server_pb.VolumeTailSenderRequest) (rpc.ServerStream[volume_server_pb.VolumeTailSenderResponse], error) {
	s, err := a.stream.VolumeTailSender(VolumeTailSenderReqInput(in))
	if err != nil {
		return nil, err
	}
	return &volumeTailSenderClientStream{s: s}, nil
}

func (a *zapVolumeClient) VolumeEcShardRead(_ context.Context, in *volume_server_pb.VolumeEcShardReadRequest) (rpc.ServerStream[volume_server_pb.VolumeEcShardReadResponse], error) {
	s, err := a.stream.VolumeEcShardRead(VolumeEcShardReadReqInput(in))
	if err != nil {
		return nil, err
	}
	return &volumeEcShardReadClientStream{s: s}, nil
}

func (a *zapVolumeClient) VolumeTierMoveDatToRemote(_ context.Context, in *volume_server_pb.VolumeTierMoveDatToRemoteRequest) (rpc.ServerStream[volume_server_pb.VolumeTierMoveDatToRemoteResponse], error) {
	s, err := a.stream.VolumeTierMoveDatToRemote(VolumeTierMoveDatToRemoteReqInput(in))
	if err != nil {
		return nil, err
	}
	return &volumeTierMoveDatToRemoteClientStream{s: s}, nil
}

func (a *zapVolumeClient) VolumeTierMoveDatFromRemote(_ context.Context, in *volume_server_pb.VolumeTierMoveDatFromRemoteRequest) (rpc.ServerStream[volume_server_pb.VolumeTierMoveDatFromRemoteResponse], error) {
	s, err := a.stream.VolumeTierMoveDatFromRemote(VolumeTierMoveDatFromRemoteReqInput(in))
	if err != nil {
		return nil, err
	}
	return &volumeTierMoveDatFromRemoteClientStream{s: s}, nil
}

func (a *zapVolumeClient) Query(_ context.Context, in *volume_server_pb.QueryRequest) (rpc.ServerStream[volume_server_pb.QueriedStripe], error) {
	s, err := a.stream.Query(QueryReqInput(in))
	if err != nil {
		return nil, err
	}
	return &queryClientStream{s: s}, nil
}

// --- client-streaming RPC --------------------------------------------------

func (a *zapVolumeClient) ReceiveFile(_ context.Context) (rpc.ClientStream[volume_server_pb.ReceiveFileRequest, volume_server_pb.ReceiveFileResponse], error) {
	// The gRPC contract opens with no message, but the wire stream's OpenStream
	// payload IS the first frame and the server replays it as its first Recv
	// (serverStream.Recv). So the FIRST Send becomes the open's init frame — the
	// real info frame rides as init exactly like every server-streaming RPC's
	// request does. There is no empty opener: an opener with Data=0 (oneof-unset)
	// would surface to the handler's oneof switch as "unknown message type".
	return &receiveFileClientStream{open: a.stream.ReceiveFile}, nil
}

type vacuumVolumeCompactClientStream struct {
	s *vsw.VacuumVolumeCompactClientStream
}

func (x *vacuumVolumeCompactClientStream) Recv() (*volume_server_pb.VacuumVolumeCompactResponse, error) {
	v, err := x.s.Recv()
	if err != nil {
		return nil, err
	}
	return VacuumVolumeCompactRespFromView(v), nil
}

type volumeIncrementalCopyClientStream struct {
	s *vsw.VolumeIncrementalCopyClientStream
}

func (x *volumeIncrementalCopyClientStream) Recv() (*volume_server_pb.VolumeIncrementalCopyResponse, error) {
	v, err := x.s.Recv()
	if err != nil {
		return nil, err
	}
	return VolumeIncrementalCopyRespFromView(v), nil
}

type volumeCopyClientStream struct {
	s *vsw.VolumeCopyClientStream
}

func (x *volumeCopyClientStream) Recv() (*volume_server_pb.VolumeCopyResponse, error) {
	v, err := x.s.Recv()
	if err != nil {
		return nil, err
	}
	return VolumeCopyRespFromView(v), nil
}

type copyFileClientStream struct {
	s *vsw.CopyFileClientStream
}

func (x *copyFileClientStream) Recv() (*volume_server_pb.CopyFileResponse, error) {
	v, err := x.s.Recv()
	if err != nil {
		return nil, err
	}
	return CopyFileRespFromView(v), nil
}

type readAllNeedlesClientStream struct {
	s *vsw.ReadAllNeedlesClientStream
}

func (x *readAllNeedlesClientStream) Recv() (*volume_server_pb.ReadAllNeedlesResponse, error) {
	v, err := x.s.Recv()
	if err != nil {
		return nil, err
	}
	return ReadAllNeedlesRespFromView(v), nil
}

type volumeTailSenderClientStream struct {
	s *vsw.VolumeTailSenderClientStream
}

func (x *volumeTailSenderClientStream) Recv() (*volume_server_pb.VolumeTailSenderResponse, error) {
	v, err := x.s.Recv()
	if err != nil {
		return nil, err
	}
	return VolumeTailSenderRespFromView(v), nil
}

type volumeEcShardReadClientStream struct {
	s *vsw.VolumeEcShardReadClientStream
}

func (x *volumeEcShardReadClientStream) Recv() (*volume_server_pb.VolumeEcShardReadResponse, error) {
	v, err := x.s.Recv()
	if err != nil {
		return nil, err
	}
	return VolumeEcShardReadRespFromView(v), nil
}

type volumeTierMoveDatToRemoteClientStream struct {
	s *vsw.VolumeTierMoveDatToRemoteClientStream
}

func (x *volumeTierMoveDatToRemoteClientStream) Recv() (*volume_server_pb.VolumeTierMoveDatToRemoteResponse, error) {
	v, err := x.s.Recv()
	if err != nil {
		return nil, err
	}
	return VolumeTierMoveDatToRemoteRespFromView(v), nil
}

type volumeTierMoveDatFromRemoteClientStream struct {
	s *vsw.VolumeTierMoveDatFromRemoteClientStream
}

func (x *volumeTierMoveDatFromRemoteClientStream) Recv() (*volume_server_pb.VolumeTierMoveDatFromRemoteResponse, error) {
	v, err := x.s.Recv()
	if err != nil {
		return nil, err
	}
	return VolumeTierMoveDatFromRemoteRespFromView(v), nil
}

type queryClientStream struct {
	s *vsw.QueryClientStream
}

func (x *queryClientStream) Recv() (*volume_server_pb.QueriedStripe, error) {
	v, err := x.s.Recv()
	if err != nil {
		return nil, err
	}
	return QueriedStripeFromView(v), nil
}

// receiveFileClientStream adapts the one client-streaming RPC. The wire stream
// opens lazily on the FIRST Send: that frame becomes the open's init payload (the
// server replays init as its first Recv), so the real info frame — never an empty
// opener — leads the stream. Subsequent Sends ship content frames; CloseAndRecv
// half-closes and reads the terminal ReceiveFileResponse.
type receiveFileClientStream struct {
	open func(vsw.ReceiveFileRequestInput) (*vsw.ReceiveFileClientStream, error)
	s    *vsw.ReceiveFileClientStream
}

func (x *receiveFileClientStream) Send(in *volume_server_pb.ReceiveFileRequest) error {
	if x.s == nil {
		s, err := x.open(receiveFileReqInput(in))
		if err != nil {
			return err
		}
		x.s = s
		return nil
	}
	return x.s.Send(receiveFileReqInput(in))
}
func (x *receiveFileClientStream) CloseSend() error {
	if x.s == nil {
		// No frame was ever sent; open with an empty init so the server sees a
		// stream it can half-close cleanly (the handler's loop returns io.EOF
		// with zero bytes written). This path is unused by the EC distributor,
		// which always sends an info frame first, but keeps CloseSend total.
		s, err := x.open(vsw.ReceiveFileRequestInput{})
		if err != nil {
			return err
		}
		x.s = s
	}
	return x.s.CloseSend()
}
func (x *receiveFileClientStream) CloseAndRecv() (*volume_server_pb.ReceiveFileResponse, error) {
	if err := x.CloseSend(); err != nil {
		return nil, err
	}
	v, err := x.s.Reply()
	if err != nil {
		return nil, err
	}
	return ReceiveFileRespFromView(v), nil
}
