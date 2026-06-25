// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// stream_server.go is the server-side half of the broker's 7 streaming RPCs: a
// mq_brokerwire.StreamBroker that wraps the existing mq_pb.HanzoMessagingServer
// so its streaming methods answer over transport.ListenStream. ServeStream
// switches on the rpc ordinal, wraps the inbound transport.Stream into the gRPC
// streaming-server interface each engine method expects, and calls the method:
//
//   - PublisherToPubBalancer / SubscriberToSubCoordinator / PublishMessage /
//     SubscribeMessage / PublishFollowMe are bidi: the adapter's Recv decodes
//     each inbound wire frame to *mq_pb.<Rpc>Request, Send encodes each
//     *mq_pb.<Rpc>Response to a wire frame.
//   - SubscribeFollowMe is client-stream: Recv until io.EOF, then SendAndClose
//     ships the single summary response.
//   - GetUnflushedMessages is server-stream: its sole request rode the opener
//     (init), and Send ships each response item.
//
// Each adapter serves Context() from the per-stream s.Context() (zap-proto/go
// v1.6.2+): cancelled when the stream ends OR the peer drops the connection, so a
// long-lived idle subscription that then disconnects is observed and the engine
// handler returns — no goroutine leak. The adapters satisfy the grpc-free
// HanzoMessaging_*Server seams (Send/Recv/Context, plus SendAndClose for the
// client-stream RPC) directly — no grpc.ServerStream plumbing.

package mq

import (
	"context"

	"github.com/zap-proto/go/transport"

	mq_pb "github.com/hanzoai/s3/s3/pb/mq_pb"
	mq_brokerwire "github.com/hanzoai/s3/s3/wire/mq_broker"
)

// streamServerBackend adapts a mq_pb.HanzoMessagingServer to mq_brokerwire.StreamBroker.
type streamServerBackend struct {
	srv mq_pb.HanzoMessagingServer
}

// NewStreamServer returns a mq_brokerwire.StreamBroker that serves srv's 7
// streaming RPCs over ZAP. Pass it to mq_brokerwire.Serve (as the stream arg) or
// use its ServeStream as the transport.StreamHandler.
func NewStreamServer(srv mq_pb.HanzoMessagingServer) mq_brokerwire.StreamBroker {
	return streamServerBackend{srv: srv}
}

var _ mq_brokerwire.StreamBroker = streamServerBackend{}

// ServeStream dispatches one inbound stream to the matching engine method. It
// runs on its own goroutine (per transport.StreamHandler contract); returning
// half-closes the send side. Engine errors end the stream — the peer's Recv sees
// the connection close — matching the gRPC behaviour where a non-nil return
// terminates the RPC.
func (b streamServerBackend) ServeStream(method uint32, init []byte, s transport.Stream) {
	switch method {
	case mq_brokerwire.HanzoMessagingPublisherToPubBalancerOrdinal:
		_ = b.srv.PublisherToPubBalancer(&pubBalancerServerStream{ctx: s.Context(), s: s})
	case mq_brokerwire.HanzoMessagingSubscriberToSubCoordinatorOrdinal:
		_ = b.srv.SubscriberToSubCoordinator(&subCoordinatorServerStream{ctx: s.Context(), s: s})
	case mq_brokerwire.HanzoMessagingPublishMessageOrdinal:
		_ = b.srv.PublishMessage(&publishServerStream{ctx: s.Context(), s: s})
	case mq_brokerwire.HanzoMessagingSubscribeMessageOrdinal:
		_ = b.srv.SubscribeMessage(&subscribeServerStream{ctx: s.Context(), s: s})
	case mq_brokerwire.HanzoMessagingPublishFollowMeOrdinal:
		_ = b.srv.PublishFollowMe(&publishFollowMeServerStream{ctx: s.Context(), s: s})
	case mq_brokerwire.HanzoMessagingSubscribeFollowMeOrdinal:
		_ = b.srv.SubscribeFollowMe(&subscribeFollowMeServerStream{ctx: s.Context(), s: s})
	case mq_brokerwire.HanzoMessagingGetUnflushedMessagesOrdinal:
		req, err := GetUnflushedMessagesReqFromWire(init)
		if err != nil {
			return
		}
		_ = b.srv.GetUnflushedMessages(req, &getUnflushedServerStream{ctx: s.Context(), s: s})
	default:
		// Unknown streaming ordinal: end immediately.
	}
}

// GetUnflushedMessagesReqFromWire decodes the server-stream's opener request.
func GetUnflushedMessagesReqFromWire(b []byte) (*mq_pb.GetUnflushedMessagesRequest, error) {
	v, err := mq_brokerwire.WrapGetUnflushedMessagesRequest(b)
	if err != nil {
		return nil, err
	}
	return &mq_pb.GetUnflushedMessagesRequest{
		Topic:             topicFromWire(v.Topic()),
		Partition:         partitionFromWire(v.Partition()),
		StartBufferOffset: v.StartBufferOffset(),
	}, nil
}

// --- PublisherToPubBalancer (bidi) ---

type pubBalancerServerStream struct {
	ctx context.Context
	s   transport.Stream
}

func (x *pubBalancerServerStream) Context() context.Context { return x.ctx }
func (x *pubBalancerServerStream) Send(resp *mq_pb.PublisherToPubBalancerResponse) error {
	return x.s.Send(mq_brokerwire.NewPublisherToPubBalancerResponse(mq_brokerwire.PublisherToPubBalancerResponseInput{}))
}
func (x *pubBalancerServerStream) Recv() (*mq_pb.PublisherToPubBalancerRequest, error) {
	b, err := x.s.Recv()
	if err != nil {
		return nil, err
	}
	return PublisherToPubBalancerReqFromWire(b)
}

// --- SubscriberToSubCoordinator (bidi) ---

type subCoordinatorServerStream struct {
	ctx context.Context
	s   transport.Stream
}

func (x *subCoordinatorServerStream) Context() context.Context { return x.ctx }
func (x *subCoordinatorServerStream) Send(resp *mq_pb.SubscriberToSubCoordinatorResponse) error {
	return x.s.Send(SubscriberToSubCoordinatorRespToWire(resp))
}
func (x *subCoordinatorServerStream) Recv() (*mq_pb.SubscriberToSubCoordinatorRequest, error) {
	b, err := x.s.Recv()
	if err != nil {
		return nil, err
	}
	return SubscriberToSubCoordinatorReqFromWire(b)
}

// --- PublishMessage (bidi) ---

type publishServerStream struct {
	ctx context.Context
	s   transport.Stream
}

func (x *publishServerStream) Context() context.Context { return x.ctx }
func (x *publishServerStream) Send(resp *mq_pb.PublishMessageResponse) error {
	return x.s.Send(PublishMessageRespToWire(resp))
}
func (x *publishServerStream) Recv() (*mq_pb.PublishMessageRequest, error) {
	b, err := x.s.Recv()
	if err != nil {
		return nil, err
	}
	return PublishMessageReqFromWire(b)
}

// --- SubscribeMessage (bidi) ---

type subscribeServerStream struct {
	ctx context.Context
	s   transport.Stream
}

func (x *subscribeServerStream) Context() context.Context { return x.ctx }
func (x *subscribeServerStream) Send(resp *mq_pb.SubscribeMessageResponse) error {
	return x.s.Send(SubscribeMessageRespToWire(resp))
}
func (x *subscribeServerStream) Recv() (*mq_pb.SubscribeMessageRequest, error) {
	b, err := x.s.Recv()
	if err != nil {
		return nil, err
	}
	return SubscribeMessageReqFromWire(b)
}

// --- PublishFollowMe (bidi) ---

type publishFollowMeServerStream struct {
	ctx context.Context
	s   transport.Stream
}

func (x *publishFollowMeServerStream) Context() context.Context { return x.ctx }
func (x *publishFollowMeServerStream) Send(resp *mq_pb.PublishFollowMeResponse) error {
	return x.s.Send(PublishFollowMeRespToWire(resp))
}
func (x *publishFollowMeServerStream) Recv() (*mq_pb.PublishFollowMeRequest, error) {
	b, err := x.s.Recv()
	if err != nil {
		return nil, err
	}
	return PublishFollowMeReqFromWire(b)
}

// --- SubscribeFollowMe (client-stream) ---
//
// The engine Recv's until io.EOF then calls SendAndClose with the single summary.
type subscribeFollowMeServerStream struct {
	ctx context.Context
	s   transport.Stream
}

func (x *subscribeFollowMeServerStream) Context() context.Context { return x.ctx }
func (x *subscribeFollowMeServerStream) Recv() (*mq_pb.SubscribeFollowMeRequest, error) {
	b, err := x.s.Recv()
	if err != nil {
		return nil, err
	}
	return SubscribeFollowMeReqFromWire(b)
}
func (x *subscribeFollowMeServerStream) SendAndClose(resp *mq_pb.SubscribeFollowMeResponse) error {
	if err := x.s.Send(SubscribeFollowMeRespToWire(resp)); err != nil {
		return err
	}
	return x.s.CloseSend()
}

// --- GetUnflushedMessages (server-stream) ---
//
// The request rode the opener (decoded in ServeStream); the engine only Send's.
type getUnflushedServerStream struct {
	ctx context.Context
	s   transport.Stream
}

func (x *getUnflushedServerStream) Context() context.Context { return x.ctx }
func (x *getUnflushedServerStream) Send(resp *mq_pb.GetUnflushedMessagesResponse) error {
	return x.s.Send(GetUnflushedMessagesRespToWire(resp))
}
