// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// stream_server.go is the server-side half of the volume server's 11 streaming
// RPCs: the vsw.VolumeServerStore streaming methods that wrap the existing
// volume_server_pb.VolumeServerServer so its streaming methods answer over
// transport.ListenStream.
//
//   - 10 server-streaming (VacuumVolumeCompact, VolumeIncrementalCopy,
//     VolumeCopy, CopyFile, ReadAllNeedles, VolumeTailSender, VolumeEcShardRead,
//     VolumeTierMoveDatToRemote, VolumeTierMoveDatFromRemote, Query): read the
//     opening request view (Init), build the pb request, and drive the engine
//     method through a Send-adapter that ships each pb response as a zero-copy
//     wire frame.
//   - 1 client-streaming (ReceiveFile): the engine method drives stream.Recv()
//     until io.EOF then stream.SendAndClose(resp); the adapter maps each inbound
//     wire frame to the pb oneof request and ships the single terminal reply.
//
// The grpc.ServerStream embedded in each adapter supplies the streaming
// interface's remaining methods (Header/Trailer/SetHeader/...), which the volume
// engine does not call on these paths — exactly the filerzap doctrine.

package volumezap

import (
	"context"

	"google.golang.org/grpc"

	volume_server_pb "github.com/hanzoai/s3/s3/pb/volume_server_pb"
	vsw "github.com/hanzoai/s3/s3/wire/volume_server"
)

// --- server-streaming (10): Init -> pb request -> engine method(req, sendAdapter) ---

func (b serverBackend) VacuumVolumeCompact(s vsw.VacuumVolumeCompactServerStream) error {
	v, err := s.Init()
	if err != nil {
		return err
	}
	return b.vs.VacuumVolumeCompact(VacuumVolumeCompactReqFromView(v), &vacuumVolumeCompactSend{ctx: b.ctx, out: s})
}

func (b serverBackend) VolumeIncrementalCopy(s vsw.VolumeIncrementalCopyServerStream) error {
	v, err := s.Init()
	if err != nil {
		return err
	}
	return b.vs.VolumeIncrementalCopy(VolumeIncrementalCopyReqFromView(v), &volumeIncrementalCopySend{ctx: b.ctx, out: s})
}

func (b serverBackend) VolumeCopy(s vsw.VolumeCopyServerStream) error {
	v, err := s.Init()
	if err != nil {
		return err
	}
	return b.vs.VolumeCopy(VolumeCopyReqFromView(v), &volumeCopySend{ctx: b.ctx, out: s})
}

func (b serverBackend) CopyFile(s vsw.CopyFileServerStream) error {
	v, err := s.Init()
	if err != nil {
		return err
	}
	return b.vs.CopyFile(CopyFileReqFromView(v), &copyFileSend{ctx: b.ctx, out: s})
}

func (b serverBackend) ReadAllNeedles(s vsw.ReadAllNeedlesServerStream) error {
	v, err := s.Init()
	if err != nil {
		return err
	}
	return b.vs.ReadAllNeedles(ReadAllNeedlesReqFromView(v), &readAllNeedlesSend{ctx: b.ctx, out: s})
}

func (b serverBackend) VolumeTailSender(s vsw.VolumeTailSenderServerStream) error {
	v, err := s.Init()
	if err != nil {
		return err
	}
	return b.vs.VolumeTailSender(VolumeTailSenderReqFromView(v), &volumeTailSenderSend{ctx: b.ctx, out: s})
}

func (b serverBackend) VolumeEcShardRead(s vsw.VolumeEcShardReadServerStream) error {
	v, err := s.Init()
	if err != nil {
		return err
	}
	return b.vs.VolumeEcShardRead(VolumeEcShardReadReqFromView(v), &volumeEcShardReadSend{ctx: b.ctx, out: s})
}

func (b serverBackend) VolumeTierMoveDatToRemote(s vsw.VolumeTierMoveDatToRemoteServerStream) error {
	v, err := s.Init()
	if err != nil {
		return err
	}
	return b.vs.VolumeTierMoveDatToRemote(VolumeTierMoveDatToRemoteReqFromView(v), &volumeTierToRemoteSend{ctx: b.ctx, out: s})
}

func (b serverBackend) VolumeTierMoveDatFromRemote(s vsw.VolumeTierMoveDatFromRemoteServerStream) error {
	v, err := s.Init()
	if err != nil {
		return err
	}
	return b.vs.VolumeTierMoveDatFromRemote(VolumeTierMoveDatFromRemoteReqFromView(v), &volumeTierFromRemoteSend{ctx: b.ctx, out: s})
}

func (b serverBackend) Query(s vsw.QueryServerStream) error {
	v, err := s.Init()
	if err != nil {
		return err
	}
	return b.vs.Query(QueryReqFromView(v), &querySend{ctx: b.ctx, out: s})
}

// --- client-streaming (1): engine drives Recv()/SendAndClose() ---

func (b serverBackend) ReceiveFile(s vsw.ReceiveFileServerStream) error {
	return b.vs.ReceiveFile(&receiveFileRecv{ctx: b.ctx, in: s})
}

// --- server-stream Send adapters (pb Send -> wire frame on the ZAP stream) ---

type vacuumVolumeCompactSend struct {
	grpc.ServerStream
	ctx context.Context
	out vsw.VacuumVolumeCompactServerStream
}

func (a *vacuumVolumeCompactSend) Send(resp *volume_server_pb.VacuumVolumeCompactResponse) error {
	return a.out.Send(VacuumVolumeCompactRespToInput(resp))
}
func (a *vacuumVolumeCompactSend) Context() context.Context { return a.ctx }

type volumeIncrementalCopySend struct {
	grpc.ServerStream
	ctx context.Context
	out vsw.VolumeIncrementalCopyServerStream
}

func (a *volumeIncrementalCopySend) Send(resp *volume_server_pb.VolumeIncrementalCopyResponse) error {
	return a.out.Send(VolumeIncrementalCopyRespToInput(resp))
}
func (a *volumeIncrementalCopySend) Context() context.Context { return a.ctx }

type volumeCopySend struct {
	grpc.ServerStream
	ctx context.Context
	out vsw.VolumeCopyServerStream
}

func (a *volumeCopySend) Send(resp *volume_server_pb.VolumeCopyResponse) error {
	return a.out.Send(VolumeCopyRespToInput(resp))
}
func (a *volumeCopySend) Context() context.Context { return a.ctx }

type copyFileSend struct {
	grpc.ServerStream
	ctx context.Context
	out vsw.CopyFileServerStream
}

func (a *copyFileSend) Send(resp *volume_server_pb.CopyFileResponse) error {
	return a.out.Send(CopyFileRespToInput(resp))
}
func (a *copyFileSend) Context() context.Context { return a.ctx }

type readAllNeedlesSend struct {
	grpc.ServerStream
	ctx context.Context
	out vsw.ReadAllNeedlesServerStream
}

func (a *readAllNeedlesSend) Send(resp *volume_server_pb.ReadAllNeedlesResponse) error {
	return a.out.Send(ReadAllNeedlesRespToInput(resp))
}
func (a *readAllNeedlesSend) Context() context.Context { return a.ctx }

type volumeTailSenderSend struct {
	grpc.ServerStream
	ctx context.Context
	out vsw.VolumeTailSenderServerStream
}

func (a *volumeTailSenderSend) Send(resp *volume_server_pb.VolumeTailSenderResponse) error {
	return a.out.Send(VolumeTailSenderRespToInput(resp))
}
func (a *volumeTailSenderSend) Context() context.Context { return a.ctx }

type volumeEcShardReadSend struct {
	grpc.ServerStream
	ctx context.Context
	out vsw.VolumeEcShardReadServerStream
}

func (a *volumeEcShardReadSend) Send(resp *volume_server_pb.VolumeEcShardReadResponse) error {
	return a.out.Send(VolumeEcShardReadRespToInput(resp))
}
func (a *volumeEcShardReadSend) Context() context.Context { return a.ctx }

type volumeTierToRemoteSend struct {
	grpc.ServerStream
	ctx context.Context
	out vsw.VolumeTierMoveDatToRemoteServerStream
}

func (a *volumeTierToRemoteSend) Send(resp *volume_server_pb.VolumeTierMoveDatToRemoteResponse) error {
	return a.out.Send(VolumeTierMoveDatToRemoteRespToInput(resp))
}
func (a *volumeTierToRemoteSend) Context() context.Context { return a.ctx }

type volumeTierFromRemoteSend struct {
	grpc.ServerStream
	ctx context.Context
	out vsw.VolumeTierMoveDatFromRemoteServerStream
}

func (a *volumeTierFromRemoteSend) Send(resp *volume_server_pb.VolumeTierMoveDatFromRemoteResponse) error {
	return a.out.Send(VolumeTierMoveDatFromRemoteRespToInput(resp))
}
func (a *volumeTierFromRemoteSend) Context() context.Context { return a.ctx }

type querySend struct {
	grpc.ServerStream
	ctx context.Context
	out vsw.QueryServerStream
}

func (a *querySend) Send(resp *volume_server_pb.QueriedStripe) error {
	return a.out.Send(QueriedStripeToInput(resp))
}
func (a *querySend) Context() context.Context { return a.ctx }

// --- client-stream Recv adapter (wire frame -> pb oneof request; reply on close) ---

// receiveFileRecv implements grpc.ClientStreamingServer[ReceiveFileRequest,
// ReceiveFileResponse] (= volume_server_pb.VolumeServer_ReceiveFileServer): the
// engine loops Recv() until io.EOF (the wire stream's Recv returns it when the
// client half-closes) then calls SendAndClose once with the terminal reply.
type receiveFileRecv struct {
	grpc.ServerStream
	ctx context.Context
	in  vsw.ReceiveFileServerStream
}

func (a *receiveFileRecv) Recv() (*volume_server_pb.ReceiveFileRequest, error) {
	v, err := a.in.Recv()
	if err != nil {
		return nil, err
	}
	return ReceiveFileReqFromView(v), nil
}
func (a *receiveFileRecv) SendAndClose(resp *volume_server_pb.ReceiveFileResponse) error {
	return a.in.Reply(ReceiveFileRespToInput(resp))
}
func (a *receiveFileRecv) Context() context.Context { return a.ctx }
