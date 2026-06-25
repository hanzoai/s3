// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// master_zap_stream.go serves the master's three bidirectional streaming RPCs —
// SendHeartbeat, KeepConnected, StreamAssign — over the native ZAP transport
// (the same listener as the unary masterwire.Backend, on grpcPort+10000). It is
// the streaming analogue of master_zap_bridge.go: it satisfies
// masterstream.Server by wrapping each transport stream half in the grpc-free
// master_pb.Hanzo_*Server seam (Recv/Send/Context over *master_pb messages, with
// the masterzap converters doing the ZAP<->pb translation) and delegating to the
// unchanged *MasterServer streaming handlers. gRPC is gone end to end.
package s3server

import (
	"context"

	"github.com/hanzoai/s3/s3/masterzap"
	"github.com/hanzoai/s3/s3/pb/master_pb"
	masterwire "github.com/hanzoai/s3/s3/wire/master"
	masterstream "github.com/hanzoai/s3/s3/wire/master/masterstream"
)

// masterZapStream implements masterstream.Server by delegating to ms's existing
// streaming handlers via the grpc-free seam adapters below.
type masterZapStream struct{ ms *MasterServer }

// NewMasterZapStreamServer returns the masterstream.Server backed by ms. Pass it
// to masterstream.Handler / transport.ListenStream alongside the unary Dispatch.
func NewMasterZapStreamServer(ms *MasterServer) masterstream.Server { return masterZapStream{ms: ms} }

func (z masterZapStream) SendHeartbeat(_ masterwire.Heartbeat, s *masterstream.HeartbeatStream) error {
	return z.ms.SendHeartbeat(zapHeartbeatServer{s: s})
}

func (z masterZapStream) KeepConnected(_ masterwire.KeepConnectedRequest, s *masterstream.KeepConnectedStream) error {
	return z.ms.KeepConnected(zapKeepConnectedServer{s: s})
}

func (z masterZapStream) StreamAssign(_ masterwire.AssignRequest, s *masterstream.AssignStream) error {
	return z.ms.StreamAssign(zapAssignServer{s: s})
}

// --- seam adapters: transport stream half -> master_pb.Hanzo_*Server ---
//
// The opener frame (the empty *Input the client uses to open the stream) is
// dropped: the masterstream server delivers it as init and s.Recv yields the
// real frames that follow, which is exactly what the handlers read.

type zapHeartbeatServer struct{ s *masterstream.HeartbeatStream }

func (x zapHeartbeatServer) Context() context.Context { return x.s.Context() }
func (x zapHeartbeatServer) Recv() (*master_pb.Heartbeat, error) {
	b, err := x.s.RecvBytes()
	if err != nil {
		return nil, err
	}
	return masterzap.HeartbeatReqFromWire(b)
}
func (x zapHeartbeatServer) Send(r *master_pb.HeartbeatResponse) error {
	return x.s.Send(masterzap.HeartbeatRespToWire(r))
}

type zapKeepConnectedServer struct {
	s *masterstream.KeepConnectedStream
}

func (x zapKeepConnectedServer) Context() context.Context { return x.s.Context() }
func (x zapKeepConnectedServer) Recv() (*master_pb.KeepConnectedRequest, error) {
	b, err := x.s.RecvBytes()
	if err != nil {
		return nil, err
	}
	return masterzap.KeepConnectedReqFromWire(b)
}
func (x zapKeepConnectedServer) Send(r *master_pb.KeepConnectedResponse) error {
	return x.s.Send(masterzap.KeepConnectedRespToWire(r))
}

type zapAssignServer struct{ s *masterstream.AssignStream }

func (x zapAssignServer) Context() context.Context { return x.s.Context() }
func (x zapAssignServer) Recv() (*master_pb.AssignRequest, error) {
	b, err := x.s.RecvBytes()
	if err != nil {
		return nil, err
	}
	return masterzap.AssignReqFromWire(b)
}
func (x zapAssignServer) Send(r *master_pb.AssignResponse) error {
	return x.s.Send(masterzap.AssignRespToWire(r))
}
