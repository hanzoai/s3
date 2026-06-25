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
// Each adapter's Context() delegates to the wire stream's per-stream context
// (vsw.<Rpc>ServerStream.Context() -> transport.Stream.Context()), which is
// cancelled when the stream ends OR the peer drops the connection (zap-proto/go
// v1.6.2+). A backend handler that gates an idle wait on stream.Context().Done()
// — or passes it to request-scoped auth (checkGrpcAdminAuth) — is therefore
// released on client disconnect, no goroutine leak. This mirrors the filerzap
// stream adapters exactly. Each adapter satisfies the grpc-free seam in
// volume_server_pb (Send/Context for server-streams; Recv/SendAndClose/Context
// for ReceiveFile) — no grpc plumbing.

package volumezap

import (
	"context"

	volume_server_pb "github.com/hanzoai/s3/s3/pb/volume_server_pb"
	vsw "github.com/hanzoai/s3/s3/wire/volume_server"
)

// --- server-streaming (10): Init -> pb request -> engine method(req, sendAdapter) ---

func (b serverBackend) VacuumVolumeCompact(s vsw.VacuumVolumeCompactServerStream) error {
	v, err := s.Init()
	if err != nil {
		return err
	}
	return b.vs.VacuumVolumeCompact(VacuumVolumeCompactReqFromView(v), &vacuumVolumeCompactSend{out: s})
}

func (b serverBackend) VolumeIncrementalCopy(s vsw.VolumeIncrementalCopyServerStream) error {
	v, err := s.Init()
	if err != nil {
		return err
	}
	return b.vs.VolumeIncrementalCopy(VolumeIncrementalCopyReqFromView(v), &volumeIncrementalCopySend{out: s})
}

func (b serverBackend) VolumeCopy(s vsw.VolumeCopyServerStream) error {
	v, err := s.Init()
	if err != nil {
		return err
	}
	return b.vs.VolumeCopy(VolumeCopyReqFromView(v), &volumeCopySend{out: s})
}

func (b serverBackend) CopyFile(s vsw.CopyFileServerStream) error {
	v, err := s.Init()
	if err != nil {
		return err
	}
	return b.vs.CopyFile(CopyFileReqFromView(v), &copyFileSend{out: s})
}

func (b serverBackend) ReadAllNeedles(s vsw.ReadAllNeedlesServerStream) error {
	v, err := s.Init()
	if err != nil {
		return err
	}
	return b.vs.ReadAllNeedles(ReadAllNeedlesReqFromView(v), &readAllNeedlesSend{out: s})
}

func (b serverBackend) VolumeTailSender(s vsw.VolumeTailSenderServerStream) error {
	v, err := s.Init()
	if err != nil {
		return err
	}
	return b.vs.VolumeTailSender(VolumeTailSenderReqFromView(v), &volumeTailSenderSend{out: s})
}

func (b serverBackend) VolumeEcShardRead(s vsw.VolumeEcShardReadServerStream) error {
	v, err := s.Init()
	if err != nil {
		return err
	}
	return b.vs.VolumeEcShardRead(VolumeEcShardReadReqFromView(v), &volumeEcShardReadSend{out: s})
}

func (b serverBackend) VolumeTierMoveDatToRemote(s vsw.VolumeTierMoveDatToRemoteServerStream) error {
	v, err := s.Init()
	if err != nil {
		return err
	}
	return b.vs.VolumeTierMoveDatToRemote(VolumeTierMoveDatToRemoteReqFromView(v), &volumeTierToRemoteSend{out: s})
}

func (b serverBackend) VolumeTierMoveDatFromRemote(s vsw.VolumeTierMoveDatFromRemoteServerStream) error {
	v, err := s.Init()
	if err != nil {
		return err
	}
	return b.vs.VolumeTierMoveDatFromRemote(VolumeTierMoveDatFromRemoteReqFromView(v), &volumeTierFromRemoteSend{out: s})
}

func (b serverBackend) Query(s vsw.QueryServerStream) error {
	v, err := s.Init()
	if err != nil {
		return err
	}
	return b.vs.Query(QueryReqFromView(v), &querySend{out: s})
}

// --- client-streaming (1): engine drives Recv()/SendAndClose() ---

func (b serverBackend) ReceiveFile(s vsw.ReceiveFileServerStream) error {
	return b.vs.ReceiveFile(&receiveFileRecv{in: s})
}

// --- server-stream Send adapters (pb Send -> wire frame on the ZAP stream) ---
//
// Context() delegates to the wire stream's per-stream context, so engine
// handlers see the real cancel-on-disconnect context, not context.Background().

type vacuumVolumeCompactSend struct {
	out vsw.VacuumVolumeCompactServerStream
}

func (a *vacuumVolumeCompactSend) Send(resp *volume_server_pb.VacuumVolumeCompactResponse) error {
	return a.out.Send(VacuumVolumeCompactRespToInput(resp))
}
func (a *vacuumVolumeCompactSend) Context() context.Context { return a.out.Context() }

type volumeIncrementalCopySend struct {
	out vsw.VolumeIncrementalCopyServerStream
}

func (a *volumeIncrementalCopySend) Send(resp *volume_server_pb.VolumeIncrementalCopyResponse) error {
	return a.out.Send(VolumeIncrementalCopyRespToInput(resp))
}
func (a *volumeIncrementalCopySend) Context() context.Context { return a.out.Context() }

type volumeCopySend struct {
	out vsw.VolumeCopyServerStream
}

func (a *volumeCopySend) Send(resp *volume_server_pb.VolumeCopyResponse) error {
	return a.out.Send(VolumeCopyRespToInput(resp))
}
func (a *volumeCopySend) Context() context.Context { return a.out.Context() }

type copyFileSend struct {
	out vsw.CopyFileServerStream
}

func (a *copyFileSend) Send(resp *volume_server_pb.CopyFileResponse) error {
	return a.out.Send(CopyFileRespToInput(resp))
}
func (a *copyFileSend) Context() context.Context { return a.out.Context() }

type readAllNeedlesSend struct {
	out vsw.ReadAllNeedlesServerStream
}

func (a *readAllNeedlesSend) Send(resp *volume_server_pb.ReadAllNeedlesResponse) error {
	return a.out.Send(ReadAllNeedlesRespToInput(resp))
}
func (a *readAllNeedlesSend) Context() context.Context { return a.out.Context() }

type volumeTailSenderSend struct {
	out vsw.VolumeTailSenderServerStream
}

func (a *volumeTailSenderSend) Send(resp *volume_server_pb.VolumeTailSenderResponse) error {
	return a.out.Send(VolumeTailSenderRespToInput(resp))
}
func (a *volumeTailSenderSend) Context() context.Context { return a.out.Context() }

type volumeEcShardReadSend struct {
	out vsw.VolumeEcShardReadServerStream
}

func (a *volumeEcShardReadSend) Send(resp *volume_server_pb.VolumeEcShardReadResponse) error {
	return a.out.Send(VolumeEcShardReadRespToInput(resp))
}
func (a *volumeEcShardReadSend) Context() context.Context { return a.out.Context() }

type volumeTierToRemoteSend struct {
	out vsw.VolumeTierMoveDatToRemoteServerStream
}

func (a *volumeTierToRemoteSend) Send(resp *volume_server_pb.VolumeTierMoveDatToRemoteResponse) error {
	return a.out.Send(VolumeTierMoveDatToRemoteRespToInput(resp))
}
func (a *volumeTierToRemoteSend) Context() context.Context { return a.out.Context() }

type volumeTierFromRemoteSend struct {
	out vsw.VolumeTierMoveDatFromRemoteServerStream
}

func (a *volumeTierFromRemoteSend) Send(resp *volume_server_pb.VolumeTierMoveDatFromRemoteResponse) error {
	return a.out.Send(VolumeTierMoveDatFromRemoteRespToInput(resp))
}
func (a *volumeTierFromRemoteSend) Context() context.Context { return a.out.Context() }

type querySend struct {
	out vsw.QueryServerStream
}

func (a *querySend) Send(resp *volume_server_pb.QueriedStripe) error {
	return a.out.Send(QueriedStripeToInput(resp))
}
func (a *querySend) Context() context.Context { return a.out.Context() }

// --- client-stream Recv adapter (wire frame -> pb oneof request; reply on close) ---

// receiveFileRecv implements volume_server_pb.VolumeServer_ReceiveFileServer: the
// engine loops Recv() until io.EOF (the wire stream's Recv returns it when the
// client half-closes) then calls SendAndClose once with the terminal reply.
type receiveFileRecv struct {
	in vsw.ReceiveFileServerStream
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
func (a *receiveFileRecv) Context() context.Context { return a.in.Context() }
