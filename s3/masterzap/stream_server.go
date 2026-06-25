// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// stream_server.go is the server-side half of the master's streaming RPCs: a
// masterstream.Server that wraps the existing master_pb.HanzoServer so its
// bidirectional streaming methods answer over transport.ListenStream. All 3 are
// implemented — SendHeartbeat, KeepConnected and StreamAssign. Each is bidi
// (stream Request <-> stream Response): the master engine method takes a
// grpc.BidiStreamingServer; we hand it an adapter whose Recv/Send pump the real
// request/response frames over the ZAP stream as zero-copy wire buffers (the
// converters in stream.go). This is the master analogue of
// filerzap.NewStreamServer.
//
// Opener handling — the ONE difference from filerzap.streamMutateAdapter: the
// masterzap client (client.go) opens each bidi stream with an EMPTY opener frame
// (NewHeartbeat{} / NewKeepConnectedRequest{} / NewAssignRequest{}) purely to
// satisfy transport.OpenStream, then ships the REAL first frame via the caller's
// first Send (this is how the gRPC "create stream, then Send" shape maps onto
// zap's "OpenStream(firstFrame)"). The real master engine expects its first
// Recv() to be the real first request (e.g. KeepConnected reads the subscription
// once). So this adapter DISCARDS the empty opener (init) and reads only the real
// frames the client Sends — replaying the empty opener as a frame would feed the
// engine a bogus zero-valued first request. (filerzap replays init because the
// filer client opens with the REAL first frame; master does not.)
//
// Each handler runs against the per-stream s.Context() (zap-proto/go v1.6.2+):
// it is cancelled when the stream ends OR the peer drops the connection, so a
// long-lived heartbeat/keep-connected subscription that goes idle and then
// disconnects is observed and the engine handler returns — no goroutine leak.

package masterzap

import (
	"context"

	master_pb "github.com/hanzoai/s3/s3/pb/master_pb"
	masterwire "github.com/hanzoai/s3/s3/wire/master"
	masterstream "github.com/hanzoai/s3/s3/wire/master/masterstream"
	"google.golang.org/grpc"
)

// streamServerBackend adapts a master_pb.HanzoServer to masterstream.Server.
type streamServerBackend struct {
	ms master_pb.HanzoServer
}

// NewStreamServer returns a masterstream.Server that serves ms's streaming RPCs
// over ZAP. Pass it to masterstream.Handler / transport.ListenStream. The empty
// opener frame (init) the client opens each stream with is discarded — the real
// frames arrive via the stream's Recv (see the package comment).
func NewStreamServer(ms master_pb.HanzoServer) masterstream.Server {
	return streamServerBackend{ms: ms}
}

func (b streamServerBackend) SendHeartbeat(_ masterwire.Heartbeat, s *masterstream.HeartbeatStream) error {
	return b.ms.SendHeartbeat(&heartbeatStreamAdapter{ctx: s.Context(), s: s})
}

func (b streamServerBackend) KeepConnected(_ masterwire.KeepConnectedRequest, s *masterstream.KeepConnectedStream) error {
	return b.ms.KeepConnected(&keepConnectedStreamAdapter{ctx: s.Context(), s: s})
}

func (b streamServerBackend) StreamAssign(_ masterwire.AssignRequest, s *masterstream.AssignStream) error {
	return b.ms.StreamAssign(&assignStreamAdapter{ctx: s.Context(), s: s})
}

// heartbeatStreamAdapter implements grpc.BidiStreamingServer[Heartbeat,
// HeartbeatResponse] (= master_pb.Hanzo_SendHeartbeatServer): Recv decodes the
// next Heartbeat frame off the ZAP stream, Send ships a HeartbeatResponse frame.
// The embedded grpc.ServerStream supplies the interface's remaining methods
// (SetHeader/SendHeader/SetTrailer/SendMsg/RecvMsg), which the master engine does
// not call on this path.
type heartbeatStreamAdapter struct {
	grpc.ServerStream
	ctx context.Context
	s   *masterstream.HeartbeatStream
}

func (a *heartbeatStreamAdapter) Recv() (*master_pb.Heartbeat, error) {
	v, err := a.s.Recv()
	if err != nil {
		return nil, err
	}
	return heartbeatReqFromView(v), nil
}
func (a *heartbeatStreamAdapter) Send(resp *master_pb.HeartbeatResponse) error {
	return a.s.Send(HeartbeatRespToWire(resp))
}
func (a *heartbeatStreamAdapter) Context() context.Context { return a.ctx }

// keepConnectedStreamAdapter implements grpc.BidiStreamingServer[
// KeepConnectedRequest, KeepConnectedResponse].
type keepConnectedStreamAdapter struct {
	grpc.ServerStream
	ctx context.Context
	s   *masterstream.KeepConnectedStream
}

func (a *keepConnectedStreamAdapter) Recv() (*master_pb.KeepConnectedRequest, error) {
	v, err := a.s.Recv()
	if err != nil {
		return nil, err
	}
	return keepConnectedReqFromView(v), nil
}
func (a *keepConnectedStreamAdapter) Send(resp *master_pb.KeepConnectedResponse) error {
	return a.s.Send(KeepConnectedRespToWire(resp))
}
func (a *keepConnectedStreamAdapter) Context() context.Context { return a.ctx }

// assignStreamAdapter implements grpc.BidiStreamingServer[AssignRequest,
// AssignResponse]. StreamAssign carries the same messages as the unary Assign,
// so it reuses assignReqFromView (request) and AssignRespToWire (response).
type assignStreamAdapter struct {
	grpc.ServerStream
	ctx context.Context
	s   *masterstream.AssignStream
}

func (a *assignStreamAdapter) Recv() (*master_pb.AssignRequest, error) {
	v, err := a.s.Recv()
	if err != nil {
		return nil, err
	}
	return assignReqFromView(v), nil
}
func (a *assignStreamAdapter) Send(resp *master_pb.AssignResponse) error {
	return a.s.Send(AssignRespToWire(resp))
}
func (a *assignStreamAdapter) Context() context.Context { return a.ctx }
