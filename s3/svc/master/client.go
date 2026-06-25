// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// client.go is the canonical strangler seam from the master_pb.HanzoClient
// contract onto the native ZAP transport. It lives in package master — the
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
//     and returns an rpc.BidiStream whose Send()/Recv() encode/decode each frame
//     with the converters in stream.go.
//
// Because the adapter satisfies master_pb.HanzoClient, no call site changes:
// callers keep building *master_pb requests and reading *master_pb responses.
// The bytes are ZAP at every hop; protobuf framing and gRPC are gone. This is
// the master analogue of filer.NewZapFilerClient — same structure, same
// streaming-adapter pattern.

package master

import (
	"context"

	"github.com/zap-proto/go/transport"

	master_pb "github.com/hanzoai/s3/s3/pb/master_pb"
	"github.com/hanzoai/s3/s3/pb/rpc"
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

func (a *zapMasterClient) LookupVolume(ctx context.Context, in *master_pb.LookupVolumeRequest) (*master_pb.LookupVolumeResponse, error) {
	_, body, err := a.unary.LookupVolume(ctx, LookupVolumeReqToWire(in))
	if err != nil {
		return nil, err
	}
	return LookupVolumeRespFromWire(body)
}

func (a *zapMasterClient) Assign(ctx context.Context, in *master_pb.AssignRequest) (*master_pb.AssignResponse, error) {
	_, body, err := a.unary.Assign(ctx, AssignReqToWire(in))
	if err != nil {
		return nil, err
	}
	return AssignRespFromWire(body)
}

func (a *zapMasterClient) Statistics(ctx context.Context, in *master_pb.StatisticsRequest) (*master_pb.StatisticsResponse, error) {
	_, body, err := a.unary.Statistics(ctx, StatisticsReqToWire(in))
	if err != nil {
		return nil, err
	}
	return StatisticsRespFromWire(body)
}

func (a *zapMasterClient) CollectionList(ctx context.Context, in *master_pb.CollectionListRequest) (*master_pb.CollectionListResponse, error) {
	_, body, err := a.unary.CollectionList(ctx, CollectionListReqToWire(in))
	if err != nil {
		return nil, err
	}
	return CollectionListRespFromWire(body)
}

func (a *zapMasterClient) CollectionDelete(ctx context.Context, in *master_pb.CollectionDeleteRequest) (*master_pb.CollectionDeleteResponse, error) {
	_, body, err := a.unary.CollectionDelete(ctx, CollectionDeleteReqToWire(in))
	if err != nil {
		return nil, err
	}
	return CollectionDeleteRespFromWire(body)
}

func (a *zapMasterClient) VolumeList(ctx context.Context, in *master_pb.VolumeListRequest) (*master_pb.VolumeListResponse, error) {
	_, body, err := a.unary.VolumeList(ctx, VolumeListReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeListRespFromWire(body)
}

func (a *zapMasterClient) LookupEcVolume(ctx context.Context, in *master_pb.LookupEcVolumeRequest) (*master_pb.LookupEcVolumeResponse, error) {
	_, body, err := a.unary.LookupEcVolume(ctx, LookupEcVolumeReqToWire(in))
	if err != nil {
		return nil, err
	}
	return LookupEcVolumeRespFromWire(body)
}

func (a *zapMasterClient) VacuumVolume(ctx context.Context, in *master_pb.VacuumVolumeRequest) (*master_pb.VacuumVolumeResponse, error) {
	_, body, err := a.unary.VacuumVolume(ctx, VacuumVolumeReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VacuumVolumeRespFromWire(body)
}

func (a *zapMasterClient) DisableVacuum(ctx context.Context, in *master_pb.DisableVacuumRequest) (*master_pb.DisableVacuumResponse, error) {
	_, body, err := a.unary.DisableVacuum(ctx, DisableVacuumReqToWire(in))
	if err != nil {
		return nil, err
	}
	return DisableVacuumRespFromWire(body)
}

func (a *zapMasterClient) EnableVacuum(ctx context.Context, in *master_pb.EnableVacuumRequest) (*master_pb.EnableVacuumResponse, error) {
	_, body, err := a.unary.EnableVacuum(ctx, EnableVacuumReqToWire(in))
	if err != nil {
		return nil, err
	}
	return EnableVacuumRespFromWire(body)
}

func (a *zapMasterClient) VolumeMarkReadonly(ctx context.Context, in *master_pb.VolumeMarkReadonlyRequest) (*master_pb.VolumeMarkReadonlyResponse, error) {
	_, body, err := a.unary.VolumeMarkReadonly(ctx, VolumeMarkReadonlyReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeMarkReadonlyRespFromWire(body)
}

func (a *zapMasterClient) GetMasterConfiguration(ctx context.Context, in *master_pb.GetMasterConfigurationRequest) (*master_pb.GetMasterConfigurationResponse, error) {
	_, body, err := a.unary.GetMasterConfiguration(ctx, GetMasterConfigurationReqToWire(in))
	if err != nil {
		return nil, err
	}
	return GetMasterConfigurationRespFromWire(body)
}

func (a *zapMasterClient) ListClusterNodes(ctx context.Context, in *master_pb.ListClusterNodesRequest) (*master_pb.ListClusterNodesResponse, error) {
	_, body, err := a.unary.ListClusterNodes(ctx, ListClusterNodesReqToWire(in))
	if err != nil {
		return nil, err
	}
	return ListClusterNodesRespFromWire(body)
}

func (a *zapMasterClient) LeaseAdminToken(ctx context.Context, in *master_pb.LeaseAdminTokenRequest) (*master_pb.LeaseAdminTokenResponse, error) {
	_, body, err := a.unary.LeaseAdminToken(ctx, LeaseAdminTokenReqToWire(in))
	if err != nil {
		return nil, err
	}
	return LeaseAdminTokenRespFromWire(body)
}

func (a *zapMasterClient) ReleaseAdminToken(ctx context.Context, in *master_pb.ReleaseAdminTokenRequest) (*master_pb.ReleaseAdminTokenResponse, error) {
	_, body, err := a.unary.ReleaseAdminToken(ctx, ReleaseAdminTokenReqToWire(in))
	if err != nil {
		return nil, err
	}
	return ReleaseAdminTokenRespFromWire(body)
}

func (a *zapMasterClient) Ping(ctx context.Context, in *master_pb.PingRequest) (*master_pb.PingResponse, error) {
	_, body, err := a.unary.Ping(ctx, PingReqToWire(in))
	if err != nil {
		return nil, err
	}
	return PingRespFromWire(body)
}

func (a *zapMasterClient) RaftListClusterServers(ctx context.Context, in *master_pb.RaftListClusterServersRequest) (*master_pb.RaftListClusterServersResponse, error) {
	_, body, err := a.unary.RaftListClusterServers(ctx, RaftListClusterServersReqToWire(in))
	if err != nil {
		return nil, err
	}
	return RaftListClusterServersRespFromWire(body)
}

func (a *zapMasterClient) RaftAddServer(ctx context.Context, in *master_pb.RaftAddServerRequest) (*master_pb.RaftAddServerResponse, error) {
	_, body, err := a.unary.RaftAddServer(ctx, RaftAddServerReqToWire(in))
	if err != nil {
		return nil, err
	}
	return RaftAddServerRespFromWire(body)
}

func (a *zapMasterClient) RaftRemoveServer(ctx context.Context, in *master_pb.RaftRemoveServerRequest) (*master_pb.RaftRemoveServerResponse, error) {
	_, body, err := a.unary.RaftRemoveServer(ctx, RaftRemoveServerReqToWire(in))
	if err != nil {
		return nil, err
	}
	return RaftRemoveServerRespFromWire(body)
}

func (a *zapMasterClient) RaftLeadershipTransfer(ctx context.Context, in *master_pb.RaftLeadershipTransferRequest) (*master_pb.RaftLeadershipTransferResponse, error) {
	_, body, err := a.unary.RaftLeadershipTransfer(ctx, RaftLeadershipTransferReqToWire(in))
	if err != nil {
		return nil, err
	}
	return RaftLeadershipTransferRespFromWire(body)
}

func (a *zapMasterClient) VolumeGrow(ctx context.Context, in *master_pb.VolumeGrowRequest) (*master_pb.VolumeGrowResponse, error) {
	_, body, err := a.unary.VolumeGrow(ctx, VolumeGrowReqToWire(in))
	if err != nil {
		return nil, err
	}
	return VolumeGrowRespFromWire(body)
}

// --- streaming RPCs --------------------------------------------------------

func (a *zapMasterClient) SendHeartbeat(_ context.Context) (rpc.BidiStream[master_pb.Heartbeat, master_pb.HeartbeatResponse], error) {
	// The bidi stream carries Heartbeat frames via Send; open it with an empty
	// Heartbeat opener frame so the server has a parseable opener.
	s, err := a.stream.SendHeartbeat(masterwire.NewHeartbeat(masterwire.HeartbeatInput{}))
	if err != nil {
		return nil, err
	}
	return &zapHeartbeatClientStream{s: s}, nil
}

func (a *zapMasterClient) KeepConnected(_ context.Context) (rpc.BidiStream[master_pb.KeepConnectedRequest, master_pb.KeepConnectedResponse], error) {
	s, err := a.stream.KeepConnected(masterwire.NewKeepConnectedRequest(masterwire.KeepConnectedRequestInput{}))
	if err != nil {
		return nil, err
	}
	return &zapKeepConnectedClientStream{s: s}, nil
}

func (a *zapMasterClient) StreamAssign(_ context.Context) (rpc.BidiStream[master_pb.AssignRequest, master_pb.AssignResponse], error) {
	s, err := a.stream.StreamAssign(masterwire.NewAssignRequest(masterwire.AssignRequestInput{}))
	if err != nil {
		return nil, err
	}
	return &zapAssignClientStream{s: s}, nil
}

// The three adapter streams below satisfy rpc.BidiStream[Req, Resp]
// (Send/Recv/CloseSend) — exactly the methods the master callers use. No grpc
// ClientStream plumbing (Header/Trailer/Context/SendMsg/RecvMsg) is needed.

// zapHeartbeatClientStream adapts a masterstream SendHeartbeat stream to
// rpc.BidiStream[Heartbeat, HeartbeatResponse]. The stream is opened with an
// empty Heartbeat frame; each Send ships a real Heartbeat as a data frame and
// each Recv decodes a HeartbeatResponse.
type zapHeartbeatClientStream struct {
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
