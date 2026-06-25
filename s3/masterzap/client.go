// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// client.go is the canonical strangler seam from the master_pb.HanzoClient
// contract onto the native ZAP transport. It lives in package masterzap — the
// package that owns the per-RPC <Rpc>ReqToWire/<Rpc>RespFromWire converters — so
// the contract is bridged in exactly ONE place. It implements
// master_pb.HanzoClient but routes every call over a
// github.com/zap-proto/go transport.Conn instead of gRPC:
//
//   - the 21 unary RPCs go through masterwire.HanzoClient (over the conn's Call
//     channel), translating *master_pb.<Rpc>Request <-> ZAP buffer <->
//     *master_pb.<Rpc>Response with the converters in rpc.go;
//   - the 3 streaming RPCs (SendHeartbeat, KeepConnected, StreamAssign) are
//     bidirectional; each opens a transport stream via the masterstream client
//     and returns a grpc.BidiStreamingClient whose Send()/Recv() encode/decode
//     each frame with the converters in stream.go.
//
// Because the adapter satisfies master_pb.HanzoClient, no call site changes:
// callers keep building *master_pb requests and reading *master_pb responses.
// The bytes are ZAP at every hop; protobuf framing and gRPC are gone. This is
// the master analogue of filerzap.NewZapFilerClient — same structure, same
// streaming-adapter pattern.

package masterzap

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/zap-proto/go/transport"

	master_pb "github.com/hanzoai/s3/s3/pb/master_pb"
	masterwire "github.com/hanzoai/s3/s3/wire/master"
	masterstream "github.com/hanzoai/s3/s3/wire/master/masterstream"
)

// zapMasterClient binds a transport.Conn to the master_pb.HanzoClient contract.
// unary issues calls over the masterwire client; stream opens streams over the
// masterstream client on the SAME connection.
type zapMasterClient struct {
	conn   transport.Conn
	unary  *masterwire.HanzoClient
	stream *masterstream.Client
}

// New wraps an established transport.Conn as a master_pb.HanzoClient that routes
// over ZAP. capability is the optional capability token attached to every unary
// request (nil for none). The caller owns conn's lifecycle (Close when done);
// the adapter pools nothing.
func New(conn transport.Conn, capability []byte) master_pb.HanzoClient {
	return &zapMasterClient{
		conn:   conn,
		unary:  masterwire.NewHanzoClient(conn, capability),
		stream: masterstream.NewClient(conn),
	}
}

var _ master_pb.HanzoClient = (*zapMasterClient)(nil)

// --- unary RPCs ------------------------------------------------------------

func (a *zapMasterClient) LookupVolume(_ context.Context, in *master_pb.LookupVolumeRequest, _ ...grpc.CallOption) (*master_pb.LookupVolumeResponse, error) {
	_, body, err := a.unary.LookupVolume(LookupVolumeReqToWire(in))
	if err != nil {
		return nil, err
	}
	return LookupVolumeRespFromWire(body)
}

func (a *zapMasterClient) Assign(_ context.Context, in *master_pb.AssignRequest, _ ...grpc.CallOption) (*master_pb.AssignResponse, error) {
	_, body, err := a.unary.Assign(AssignReqToWire(in))
	if err != nil {
		return nil, err
	}
	return AssignRespFromWire(body)
}

func (a *zapMasterClient) Statistics(_ context.Context, in *master_pb.StatisticsRequest, _ ...grpc.CallOption) (*master_pb.StatisticsResponse, error) {
	_, body, err := a.unary.Statistics(StatisticsReqToWire(in))
	if err != nil {
		return nil, err
	}
	return StatisticsRespFromWire(body)
}

func (a *zapMasterClient) CollectionList(_ context.Context, in *master_pb.CollectionListRequest, _ ...grpc.CallOption) (*master_pb.CollectionListResponse, error) {
	_, body, err := a.unary.CollectionList(CollectionListReqToWire(in))
	if err != nil {
		return nil, err
	}
	return CollectionListRespFromWire(body)
}

func (a *zapMasterClient) CollectionDelete(_ context.Context, in *master_pb.CollectionDeleteRequest, _ ...grpc.CallOption) (*master_pb.CollectionDeleteResponse, error) {
	_, body, err := a.unary.CollectionDelete(CollectionDeleteReqToWire(in))
	if err != nil {
		return nil, err
	}
	return CollectionDeleteRespFromWire(body)
}

func (a *zapMasterClient) VolumeList(_ context.Context, in *master_pb.VolumeListRequest, _ ...grpc.CallOption) (*master_pb.VolumeListResponse, error) {
	_, body, err := a.unary.VolumeList(VolumeListReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeListRespFromWire(body)
}

func (a *zapMasterClient) LookupEcVolume(_ context.Context, in *master_pb.LookupEcVolumeRequest, _ ...grpc.CallOption) (*master_pb.LookupEcVolumeResponse, error) {
	_, body, err := a.unary.LookupEcVolume(LookupEcVolumeReqToWire(in))
	if err != nil {
		return nil, err
	}
	return LookupEcVolumeRespFromWire(body)
}

func (a *zapMasterClient) VacuumVolume(_ context.Context, in *master_pb.VacuumVolumeRequest, _ ...grpc.CallOption) (*master_pb.VacuumVolumeResponse, error) {
	_, body, err := a.unary.VacuumVolume(VacuumVolumeReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VacuumVolumeRespFromWire(body)
}

func (a *zapMasterClient) DisableVacuum(_ context.Context, in *master_pb.DisableVacuumRequest, _ ...grpc.CallOption) (*master_pb.DisableVacuumResponse, error) {
	_, body, err := a.unary.DisableVacuum(DisableVacuumReqToWire(in))
	if err != nil {
		return nil, err
	}
	return DisableVacuumRespFromWire(body)
}

func (a *zapMasterClient) EnableVacuum(_ context.Context, in *master_pb.EnableVacuumRequest, _ ...grpc.CallOption) (*master_pb.EnableVacuumResponse, error) {
	_, body, err := a.unary.EnableVacuum(EnableVacuumReqToWire(in))
	if err != nil {
		return nil, err
	}
	return EnableVacuumRespFromWire(body)
}

func (a *zapMasterClient) VolumeMarkReadonly(_ context.Context, in *master_pb.VolumeMarkReadonlyRequest, _ ...grpc.CallOption) (*master_pb.VolumeMarkReadonlyResponse, error) {
	_, body, err := a.unary.VolumeMarkReadonly(VolumeMarkReadonlyReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeMarkReadonlyRespFromWire(body)
}

func (a *zapMasterClient) GetMasterConfiguration(_ context.Context, in *master_pb.GetMasterConfigurationRequest, _ ...grpc.CallOption) (*master_pb.GetMasterConfigurationResponse, error) {
	_, body, err := a.unary.GetMasterConfiguration(GetMasterConfigurationReqToWire(in))
	if err != nil {
		return nil, err
	}
	return GetMasterConfigurationRespFromWire(body)
}

func (a *zapMasterClient) ListClusterNodes(_ context.Context, in *master_pb.ListClusterNodesRequest, _ ...grpc.CallOption) (*master_pb.ListClusterNodesResponse, error) {
	_, body, err := a.unary.ListClusterNodes(ListClusterNodesReqToWire(in))
	if err != nil {
		return nil, err
	}
	return ListClusterNodesRespFromWire(body)
}

func (a *zapMasterClient) LeaseAdminToken(_ context.Context, in *master_pb.LeaseAdminTokenRequest, _ ...grpc.CallOption) (*master_pb.LeaseAdminTokenResponse, error) {
	_, body, err := a.unary.LeaseAdminToken(LeaseAdminTokenReqToWire(in))
	if err != nil {
		return nil, err
	}
	return LeaseAdminTokenRespFromWire(body)
}

func (a *zapMasterClient) ReleaseAdminToken(_ context.Context, in *master_pb.ReleaseAdminTokenRequest, _ ...grpc.CallOption) (*master_pb.ReleaseAdminTokenResponse, error) {
	_, body, err := a.unary.ReleaseAdminToken(ReleaseAdminTokenReqToWire(in))
	if err != nil {
		return nil, err
	}
	return ReleaseAdminTokenRespFromWire(body)
}

func (a *zapMasterClient) Ping(_ context.Context, in *master_pb.PingRequest, _ ...grpc.CallOption) (*master_pb.PingResponse, error) {
	_, body, err := a.unary.Ping(PingReqToWire(in))
	if err != nil {
		return nil, err
	}
	return PingRespFromWire(body)
}

func (a *zapMasterClient) RaftListClusterServers(_ context.Context, in *master_pb.RaftListClusterServersRequest, _ ...grpc.CallOption) (*master_pb.RaftListClusterServersResponse, error) {
	_, body, err := a.unary.RaftListClusterServers(RaftListClusterServersReqToWire(in))
	if err != nil {
		return nil, err
	}
	return RaftListClusterServersRespFromWire(body)
}

func (a *zapMasterClient) RaftAddServer(_ context.Context, in *master_pb.RaftAddServerRequest, _ ...grpc.CallOption) (*master_pb.RaftAddServerResponse, error) {
	_, body, err := a.unary.RaftAddServer(RaftAddServerReqToWire(in))
	if err != nil {
		return nil, err
	}
	return RaftAddServerRespFromWire(body)
}

func (a *zapMasterClient) RaftRemoveServer(_ context.Context, in *master_pb.RaftRemoveServerRequest, _ ...grpc.CallOption) (*master_pb.RaftRemoveServerResponse, error) {
	_, body, err := a.unary.RaftRemoveServer(RaftRemoveServerReqToWire(in))
	if err != nil {
		return nil, err
	}
	return RaftRemoveServerRespFromWire(body)
}

func (a *zapMasterClient) RaftLeadershipTransfer(_ context.Context, in *master_pb.RaftLeadershipTransferRequest, _ ...grpc.CallOption) (*master_pb.RaftLeadershipTransferResponse, error) {
	_, body, err := a.unary.RaftLeadershipTransfer(RaftLeadershipTransferReqToWire(in))
	if err != nil {
		return nil, err
	}
	return RaftLeadershipTransferRespFromWire(body)
}

func (a *zapMasterClient) VolumeGrow(_ context.Context, in *master_pb.VolumeGrowRequest, _ ...grpc.CallOption) (*master_pb.VolumeGrowResponse, error) {
	_, body, err := a.unary.VolumeGrow(VolumeGrowReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeGrowRespFromWire(body)
}

// --- streaming RPCs --------------------------------------------------------

func (a *zapMasterClient) SendHeartbeat(_ context.Context, _ ...grpc.CallOption) (grpc.BidiStreamingClient[master_pb.Heartbeat, master_pb.HeartbeatResponse], error) {
	// The bidi stream carries Heartbeat frames via Send; open it with an empty
	// Heartbeat opener frame so the server has a parseable opener.
	s, err := a.stream.SendHeartbeat(masterwire.NewHeartbeat(masterwire.HeartbeatInput{}))
	if err != nil {
		return nil, err
	}
	return &zapHeartbeatClientStream{s: s}, nil
}

func (a *zapMasterClient) KeepConnected(_ context.Context, _ ...grpc.CallOption) (grpc.BidiStreamingClient[master_pb.KeepConnectedRequest, master_pb.KeepConnectedResponse], error) {
	s, err := a.stream.KeepConnected(masterwire.NewKeepConnectedRequest(masterwire.KeepConnectedRequestInput{}))
	if err != nil {
		return nil, err
	}
	return &zapKeepConnectedClientStream{s: s}, nil
}

func (a *zapMasterClient) StreamAssign(_ context.Context, _ ...grpc.CallOption) (grpc.BidiStreamingClient[master_pb.AssignRequest, master_pb.AssignResponse], error) {
	s, err := a.stream.StreamAssign(masterwire.NewAssignRequest(masterwire.AssignRequestInput{}))
	if err != nil {
		return nil, err
	}
	return &zapAssignClientStream{s: s}, nil
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

// zapHeartbeatClientStream adapts a masterstream SendHeartbeat stream to
// grpc.BidiStreamingClient[Heartbeat, HeartbeatResponse]. The stream is opened
// with an empty Heartbeat frame; each Send ships a real Heartbeat as a data
// frame and each Recv decodes a HeartbeatResponse.
type zapHeartbeatClientStream struct {
	zapClientStreamStub
	s *masterstream.ClientHeartbeatStream
}

func (x *zapHeartbeatClientStream) Send(in *master_pb.Heartbeat) error {
	return x.s.Send(HeartbeatReqToWire(in))
}
func (x *zapHeartbeatClientStream) Recv() (*master_pb.HeartbeatResponse, error) {
	b, err := x.s.RecvBytes()
	if err != nil {
		return nil, err
	}
	return HeartbeatRespFromWire(b)
}
func (x *zapHeartbeatClientStream) CloseSend() error { return x.s.CloseSend() }

// zapKeepConnectedClientStream adapts a masterstream KeepConnected stream.
type zapKeepConnectedClientStream struct {
	zapClientStreamStub
	s *masterstream.ClientKeepConnectedStream
}

func (x *zapKeepConnectedClientStream) Send(in *master_pb.KeepConnectedRequest) error {
	return x.s.Send(KeepConnectedReqToWire(in))
}
func (x *zapKeepConnectedClientStream) Recv() (*master_pb.KeepConnectedResponse, error) {
	b, err := x.s.RecvBytes()
	if err != nil {
		return nil, err
	}
	return KeepConnectedRespFromWire(b)
}
func (x *zapKeepConnectedClientStream) CloseSend() error { return x.s.CloseSend() }

// zapAssignClientStream adapts a masterstream StreamAssign stream.
type zapAssignClientStream struct {
	zapClientStreamStub
	s *masterstream.ClientAssignStream
}

func (x *zapAssignClientStream) Send(in *master_pb.AssignRequest) error {
	return x.s.Send(AssignReqToWire(in))
}
func (x *zapAssignClientStream) Recv() (*master_pb.AssignResponse, error) {
	b, err := x.s.RecvBytes()
	if err != nil {
		return nil, err
	}
	return AssignRespFromWire(b)
}
func (x *zapAssignClientStream) CloseSend() error { return x.s.CloseSend() }
