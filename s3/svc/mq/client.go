// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// client.go is the canonical strangler seam from the mq_pb.HanzoMessagingClient
// contract onto the native ZAP transport. It lives in package mqzap — the package
// that owns the per-RPC <Rpc>ReqToWire/<Rpc>RespFromWire converters — so the
// contract is bridged in exactly ONE place. It implements
// mq_pb.HanzoMessagingClient but routes every call over a
// github.com/zap-proto/go transport.Conn instead of gRPC:
//
//   - the 14 unary RPCs go through mq_brokerwire.HanzoMessagingClient (over the
//     conn's Call channel), translating *mq_pb.<Rpc>Request <-> ZAP buffer <->
//     *mq_pb.<Rpc>Response with the converters in rpc.go;
//   - the 7 streaming RPCs open a transport.Stream via conn.OpenStream(ordinal,
//     init) and return the rpc.* stream seam (rpc.BidiStream / rpc.ClientStream /
//     rpc.ServerStream) whose Send/Recv encode/decode each frame with the
//     converters in stream.go.
//
// Because the adapter satisfies mq_pb.HanzoMessagingClient, no call site changes:
// callers keep building *mq_pb requests and reading *mq_pb responses. The bytes
// are ZAP at every hop; protobuf framing and gRPC are gone. New is the ONE place
// the contract is bridged.

package mq

import (
	"context"

	"github.com/zap-proto/go/transport"

	mq_pb "github.com/hanzoai/s3/s3/pb/mq_pb"
	"github.com/hanzoai/s3/s3/pb/rpc"
	mq_brokerwire "github.com/hanzoai/s3/s3/wire/mq_broker"
)

// client binds a transport.Conn to the mq_pb.HanzoMessagingClient contract. unary
// issues calls over the mq_brokerwire client; streaming opens transport streams on
// the SAME connection.
type client struct {
	conn  transport.Conn
	unary *mq_brokerwire.HanzoMessagingClient
}

// New wraps an established transport.Conn as a mq_pb.HanzoMessagingClient that
// routes over ZAP. capability is the optional capability token attached to every
// unary request (nil for none). The caller owns conn's lifecycle (Close when
// done); the adapter pools nothing.
func New(conn transport.Conn, capability []byte) mq_pb.HanzoMessagingClient {
	return &client{conn: conn, unary: mq_brokerwire.NewHanzoMessagingClient(conn, capability)}
}

var _ mq_pb.HanzoMessagingClient = (*client)(nil)

// --- unary RPCs ------------------------------------------------------------

func (a *client) FindBrokerLeader(_ context.Context, in *mq_pb.FindBrokerLeaderRequest) (*mq_pb.FindBrokerLeaderResponse, error) {
	_, body, err := a.unary.FindBrokerLeader(FindBrokerLeaderReqToWire(in))
	if err != nil {
		return nil, err
	}
	return FindBrokerLeaderRespFromWire(body)
}

func (a *client) BalanceTopics(_ context.Context, in *mq_pb.BalanceTopicsRequest) (*mq_pb.BalanceTopicsResponse, error) {
	_, body, err := a.unary.BalanceTopics(BalanceTopicsReqToWire(in))
	if err != nil {
		return nil, err
	}
	return BalanceTopicsRespFromWire(body)
}

func (a *client) ListTopics(_ context.Context, in *mq_pb.ListTopicsRequest) (*mq_pb.ListTopicsResponse, error) {
	_, body, err := a.unary.ListTopics(ListTopicsReqToWire(in))
	if err != nil {
		return nil, err
	}
	return ListTopicsRespFromWire(body)
}

func (a *client) TopicExists(_ context.Context, in *mq_pb.TopicExistsRequest) (*mq_pb.TopicExistsResponse, error) {
	_, body, err := a.unary.TopicExists(TopicExistsReqToWire(in))
	if err != nil {
		return nil, err
	}
	return TopicExistsRespFromWire(body)
}

func (a *client) ConfigureTopic(_ context.Context, in *mq_pb.ConfigureTopicRequest) (*mq_pb.ConfigureTopicResponse, error) {
	_, body, err := a.unary.ConfigureTopic(ConfigureTopicReqToWire(in))
	if err != nil {
		return nil, err
	}
	return ConfigureTopicRespFromWire(body)
}

func (a *client) LookupTopicBrokers(_ context.Context, in *mq_pb.LookupTopicBrokersRequest) (*mq_pb.LookupTopicBrokersResponse, error) {
	_, body, err := a.unary.LookupTopicBrokers(LookupTopicBrokersReqToWire(in))
	if err != nil {
		return nil, err
	}
	return LookupTopicBrokersRespFromWire(body)
}

func (a *client) GetTopicConfiguration(_ context.Context, in *mq_pb.GetTopicConfigurationRequest) (*mq_pb.GetTopicConfigurationResponse, error) {
	_, body, err := a.unary.GetTopicConfiguration(GetTopicConfigurationReqToWire(in))
	if err != nil {
		return nil, err
	}
	return GetTopicConfigurationRespFromWire(body)
}

func (a *client) GetTopicPublishers(_ context.Context, in *mq_pb.GetTopicPublishersRequest) (*mq_pb.GetTopicPublishersResponse, error) {
	_, body, err := a.unary.GetTopicPublishers(GetTopicPublishersReqToWire(in))
	if err != nil {
		return nil, err
	}
	return GetTopicPublishersRespFromWire(body)
}

func (a *client) GetTopicSubscribers(_ context.Context, in *mq_pb.GetTopicSubscribersRequest) (*mq_pb.GetTopicSubscribersResponse, error) {
	_, body, err := a.unary.GetTopicSubscribers(GetTopicSubscribersReqToWire(in))
	if err != nil {
		return nil, err
	}
	return GetTopicSubscribersRespFromWire(body)
}

func (a *client) AssignTopicPartitions(_ context.Context, in *mq_pb.AssignTopicPartitionsRequest) (*mq_pb.AssignTopicPartitionsResponse, error) {
	_, body, err := a.unary.AssignTopicPartitions(AssignTopicPartitionsReqToWire(in))
	if err != nil {
		return nil, err
	}
	return AssignTopicPartitionsRespFromWire(body)
}

func (a *client) ClosePublishers(_ context.Context, in *mq_pb.ClosePublishersRequest) (*mq_pb.ClosePublishersResponse, error) {
	_, body, err := a.unary.ClosePublishers(ClosePublishersReqToWire(in))
	if err != nil {
		return nil, err
	}
	return ClosePublishersRespFromWire(body)
}

func (a *client) CloseSubscribers(_ context.Context, in *mq_pb.CloseSubscribersRequest) (*mq_pb.CloseSubscribersResponse, error) {
	_, body, err := a.unary.CloseSubscribers(CloseSubscribersReqToWire(in))
	if err != nil {
		return nil, err
	}
	return CloseSubscribersRespFromWire(body)
}

func (a *client) FetchMessage(_ context.Context, in *mq_pb.FetchMessageRequest) (*mq_pb.FetchMessageResponse, error) {
	_, body, err := a.unary.FetchMessage(FetchMessageReqToWire(in))
	if err != nil {
		return nil, err
	}
	return FetchMessageRespFromWire(body)
}

func (a *client) GetPartitionRangeInfo(_ context.Context, in *mq_pb.GetPartitionRangeInfoRequest) (*mq_pb.GetPartitionRangeInfoResponse, error) {
	_, body, err := a.unary.GetPartitionRangeInfo(GetPartitionRangeInfoReqToWire(in))
	if err != nil {
		return nil, err
	}
	return GetPartitionRangeInfoRespFromWire(body)
}

// --- streaming RPCs --------------------------------------------------------
//
// Each bidi/client streaming RPC opens with an empty opener frame; the caller
// ships every message (init included) via Send, mirroring the gRPC contract, and
// the server's StreamBroker.ServeStream Recv-loops over those frames.
// GetUnflushedMessages is server-streaming: its sole request rides the opener and
// the client only Recv's.

func (a *client) PublisherToPubBalancer(_ context.Context) (rpc.BidiStream[mq_pb.PublisherToPubBalancerRequest, mq_pb.PublisherToPubBalancerResponse], error) {
	s, err := a.conn.OpenStream(mq_brokerwire.HanzoMessagingPublisherToPubBalancerOrdinal, nil)
	if err != nil {
		return nil, err
	}
	return &pubBalancerStream{s: s}, nil
}

func (a *client) SubscriberToSubCoordinator(_ context.Context) (rpc.BidiStream[mq_pb.SubscriberToSubCoordinatorRequest, mq_pb.SubscriberToSubCoordinatorResponse], error) {
	s, err := a.conn.OpenStream(mq_brokerwire.HanzoMessagingSubscriberToSubCoordinatorOrdinal, nil)
	if err != nil {
		return nil, err
	}
	return &subCoordinatorStream{s: s}, nil
}

func (a *client) PublishMessage(_ context.Context) (rpc.BidiStream[mq_pb.PublishMessageRequest, mq_pb.PublishMessageResponse], error) {
	s, err := a.conn.OpenStream(mq_brokerwire.HanzoMessagingPublishMessageOrdinal, nil)
	if err != nil {
		return nil, err
	}
	return &publishStream{s: s}, nil
}

func (a *client) SubscribeMessage(_ context.Context) (rpc.BidiStream[mq_pb.SubscribeMessageRequest, mq_pb.SubscribeMessageResponse], error) {
	s, err := a.conn.OpenStream(mq_brokerwire.HanzoMessagingSubscribeMessageOrdinal, nil)
	if err != nil {
		return nil, err
	}
	return &subscribeStream{s: s}, nil
}

func (a *client) PublishFollowMe(_ context.Context) (rpc.BidiStream[mq_pb.PublishFollowMeRequest, mq_pb.PublishFollowMeResponse], error) {
	s, err := a.conn.OpenStream(mq_brokerwire.HanzoMessagingPublishFollowMeOrdinal, nil)
	if err != nil {
		return nil, err
	}
	return &publishFollowMeStream{s: s}, nil
}

func (a *client) SubscribeFollowMe(_ context.Context) (rpc.ClientStream[mq_pb.SubscribeFollowMeRequest, mq_pb.SubscribeFollowMeResponse], error) {
	s, err := a.conn.OpenStream(mq_brokerwire.HanzoMessagingSubscribeFollowMeOrdinal, nil)
	if err != nil {
		return nil, err
	}
	return &subscribeFollowMeStream{s: s}, nil
}

func (a *client) GetUnflushedMessages(_ context.Context, in *mq_pb.GetUnflushedMessagesRequest) (rpc.ServerStream[mq_pb.GetUnflushedMessagesResponse], error) {
	s, err := a.conn.OpenStream(mq_brokerwire.HanzoMessagingGetUnflushedMessagesOrdinal, GetUnflushedMessagesReqToWire(in))
	if err != nil {
		return nil, err
	}
	if err := s.CloseSend(); err != nil {
		return nil, err
	}
	return &getUnflushedStream{s: s}, nil
}

// The stream wrappers below implement exactly the rpc.* seam each RPC returns:
// rpc.BidiStream needs Send/Recv/CloseSend; rpc.ClientStream needs
// Send/CloseAndRecv; rpc.ServerStream needs Recv. Each frame is decoded/encoded
// with the same converters as the unary path.

// pubBalancerStream adapts a transport.Stream to the PublisherToPubBalancer bidi.
type pubBalancerStream struct{ s transport.Stream }

func (x *pubBalancerStream) Send(in *mq_pb.PublisherToPubBalancerRequest) error {
	return x.s.Send(PublisherToPubBalancerReqToWire(in))
}
func (x *pubBalancerStream) Recv() (*mq_pb.PublisherToPubBalancerResponse, error) {
	b, err := x.s.Recv()
	if err != nil {
		return nil, err
	}
	return PublisherToPubBalancerRespFromWire(b)
}
func (x *pubBalancerStream) CloseSend() error { return x.s.CloseSend() }

// subCoordinatorStream adapts a transport.Stream to the SubscriberToSubCoordinator bidi.
type subCoordinatorStream struct{ s transport.Stream }

func (x *subCoordinatorStream) Send(in *mq_pb.SubscriberToSubCoordinatorRequest) error {
	return x.s.Send(SubscriberToSubCoordinatorReqToWire(in))
}
func (x *subCoordinatorStream) Recv() (*mq_pb.SubscriberToSubCoordinatorResponse, error) {
	b, err := x.s.Recv()
	if err != nil {
		return nil, err
	}
	return SubscriberToSubCoordinatorRespFromWire(b)
}
func (x *subCoordinatorStream) CloseSend() error { return x.s.CloseSend() }

// publishStream adapts a transport.Stream to the PublishMessage bidi.
type publishStream struct{ s transport.Stream }

func (x *publishStream) Send(in *mq_pb.PublishMessageRequest) error {
	return x.s.Send(PublishMessageReqToWire(in))
}
func (x *publishStream) Recv() (*mq_pb.PublishMessageResponse, error) {
	b, err := x.s.Recv()
	if err != nil {
		return nil, err
	}
	return PublishMessageRespFromWire(b)
}
func (x *publishStream) CloseSend() error { return x.s.CloseSend() }

// subscribeStream adapts a transport.Stream to the SubscribeMessage bidi.
type subscribeStream struct{ s transport.Stream }

func (x *subscribeStream) Send(in *mq_pb.SubscribeMessageRequest) error {
	return x.s.Send(SubscribeMessageReqToWire(in))
}
func (x *subscribeStream) Recv() (*mq_pb.SubscribeMessageResponse, error) {
	b, err := x.s.Recv()
	if err != nil {
		return nil, err
	}
	return SubscribeMessageRespFromWire(b)
}
func (x *subscribeStream) CloseSend() error { return x.s.CloseSend() }

// publishFollowMeStream adapts a transport.Stream to the PublishFollowMe bidi.
type publishFollowMeStream struct{ s transport.Stream }

func (x *publishFollowMeStream) Send(in *mq_pb.PublishFollowMeRequest) error {
	return x.s.Send(PublishFollowMeReqToWire(in))
}
func (x *publishFollowMeStream) Recv() (*mq_pb.PublishFollowMeResponse, error) {
	b, err := x.s.Recv()
	if err != nil {
		return nil, err
	}
	return PublishFollowMeRespFromWire(b)
}
func (x *publishFollowMeStream) CloseSend() error { return x.s.CloseSend() }

// subscribeFollowMeStream adapts a transport.Stream to the SubscribeFollowMe
// client-stream: the client Sends until it CloseAndRecv's the single response.
type subscribeFollowMeStream struct{ s transport.Stream }

func (x *subscribeFollowMeStream) Send(in *mq_pb.SubscribeFollowMeRequest) error {
	return x.s.Send(SubscribeFollowMeReqToWire(in))
}
func (x *subscribeFollowMeStream) CloseAndRecv() (*mq_pb.SubscribeFollowMeResponse, error) {
	if err := x.s.CloseSend(); err != nil {
		return nil, err
	}
	b, err := x.s.Recv()
	if err != nil {
		return nil, err
	}
	return SubscribeFollowMeRespFromWire(b)
}

// getUnflushedStream adapts a transport.Stream to the GetUnflushedMessages
// server-stream: the request rode the opener, so the client only Recv's.
type getUnflushedStream struct{ s transport.Stream }

func (x *getUnflushedStream) Recv() (*mq_pb.GetUnflushedMessagesResponse, error) {
	b, err := x.s.Recv()
	if err != nil {
		return nil, err
	}
	return GetUnflushedMessagesRespFromWire(b)
}
