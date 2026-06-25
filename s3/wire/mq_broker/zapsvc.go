// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// Native ZAP service adapter for HanzoMessaging — the hand-written glue that
// binds the generated Dispatch/Client (mq_broker_zap.go) to a real message-queue
// broker backend. This is the mq_broker peer of s3/zapsvc/service.go (object) and
// s3/wire/s3_lifecycle/zapsvc.go (lifecycle): requests/responses are the
// zero-copy mq_brokerwire message buffers, carried over the canonical
// github.com/zap-proto/go transport. The unary RPCs are dispatched by the
// generated DispatchHanzoMessaging; the seven streaming RPCs are wired by hand
// over transport.ListenStream / Conn.OpenStream — the v1.5.0 streaming primitive
// the generator does NOT emit. No struct marshaling, no protobuf, no gRPC, no
// HTTP — the bytes ARE the message.
//
// The backend is a Go interface (Broker) expressed in pure wire terms, so this
// file imports neither the mq_pb domain model nor the real broker/filer/master
// engine. The MQ broker server wires its live implementation to this interface
// at integration time; tests use an in-memory fake (see zapsvc_roundtrip_test.go).
//
// UNARY (14): FindBrokerLeader, BalanceTopics, ListTopics, TopicExists,
// ConfigureTopic, LookupTopicBrokers, GetTopicConfiguration, GetTopicPublishers,
// GetTopicSubscribers, AssignTopicPartitions, ClosePublishers, CloseSubscribers,
// FetchMessage, GetPartitionRangeInfo — each is a request buffer -> response
// buffer call, mirroring object GetObject/PutObject exactly.
//
// STREAMING (7): PublisherToPubBalancer (bidi), SubscriberToSubCoordinator
// (bidi), PublishMessage (bidi), SubscribeMessage (bidi), PublishFollowMe (bidi),
// SubscribeFollowMe (client-stream), GetUnflushedMessages (server-stream). Each
// rides a transport.Stream: the backend's StreamBroker.ServeStream is invoked
// with the rpc ordinal, the opening payload, and the live stream to Send/Recv on.
// Frames are the same zero-copy New<Msg>/Wrap<Msg> buffers as the unary path.

package mq_brokerwire

import (
	"github.com/zap-proto/go/transport"
)

// Broker is the unary backend the ZAP service delegates to. The real MQ broker
// engine implements this; tests use an in-memory fake. Each method receives the
// zero-copy view over its request buffer (read fields with the accessor methods)
// and returns the response as the matching Input struct, which the handler builds
// into a wire buffer. The view aliases the request buffer, so copy out any bytes
// the implementation needs to retain beyond the call.
type Broker interface {
	FindBrokerLeader(req FindBrokerLeaderRequest) (FindBrokerLeaderResponseInput, error)
	BalanceTopics(req BalanceTopicsRequest) (BalanceTopicsResponseInput, error)
	ListTopics(req ListTopicsRequest) (ListTopicsResponseInput, error)
	TopicExists(req TopicExistsRequest) (TopicExistsResponseInput, error)
	ConfigureTopic(req ConfigureTopicRequest) (ConfigureTopicResponseInput, error)
	LookupTopicBrokers(req LookupTopicBrokersRequest) (LookupTopicBrokersResponseInput, error)
	GetTopicConfiguration(req GetTopicConfigurationRequest) (GetTopicConfigurationResponseInput, error)
	GetTopicPublishers(req GetTopicPublishersRequest) (GetTopicPublishersResponseInput, error)
	GetTopicSubscribers(req GetTopicSubscribersRequest) (GetTopicSubscribersResponseInput, error)
	AssignTopicPartitions(req AssignTopicPartitionsRequest) (AssignTopicPartitionsResponseInput, error)
	ClosePublishers(req ClosePublishersRequest) (ClosePublishersResponseInput, error)
	CloseSubscribers(req CloseSubscribersRequest) (CloseSubscribersResponseInput, error)
	FetchMessage(req FetchMessageRequest) (FetchMessageResponseInput, error)
	GetPartitionRangeInfo(req GetPartitionRangeInfoRequest) (GetPartitionRangeInfoResponseInput, error)
}

// StreamBroker is the streaming backend. ServeStream serves ONE inbound stream:
// method is the rpc ordinal from the open envelope (switch on the
// HanzoMessaging*Ordinal constants for the seven streaming RPCs), init is the
// opening payload (a New<Req> buffer — conventionally the InitMessage-bearing
// request), and s is the bidirectional stream to Recv from / Send on. It runs on
// its own goroutine; returning half-closes the send side (the peer's Recv sees
// io.EOF). This one method covers all three streaming shapes — server-, client-,
// and bidirectional — it is just which side calls Send vs Recv:
//
//   - PublisherToPubBalancer / SubscriberToSubCoordinator / PublishMessage /
//     SubscribeMessage / PublishFollowMe are bidi: interleave Send/Recv.
//   - SubscribeFollowMe is client-stream: Recv until io.EOF, then Send one summary.
//   - GetUnflushedMessages is server-stream: Send per item, return to end.
//
// Every Send/Recv body is a zero-copy mq_brokerwire buffer (New<Msg> / Wrap<Msg>),
// same doctrine as the unary path, just framed as a stream message.
type StreamBroker interface {
	ServeStream(method uint32, init []byte, s transport.Stream)
}

// handler adapts a Broker to the generated HanzoMessagingHandler: it Wraps each
// request buffer into its zero-copy view, calls the backend, and builds the
// response buffer with the matching New<Resp>. The streaming methods on the
// generated interface are never reached over a unary call — they are served via
// ServeStream — so they return a (harmless) empty response buffer if a caller
// ever issues them single-shot.
type handler struct{ b Broker }

func (h handler) FindBrokerLeader(req []byte) ([]byte, error) {
	v, err := WrapFindBrokerLeaderRequest(req)
	if err != nil {
		return nil, err
	}
	out, err := h.b.FindBrokerLeader(v)
	if err != nil {
		return nil, err
	}
	return NewFindBrokerLeaderResponse(out), nil
}

func (h handler) BalanceTopics(req []byte) ([]byte, error) {
	v, err := WrapBalanceTopicsRequest(req)
	if err != nil {
		return nil, err
	}
	out, err := h.b.BalanceTopics(v)
	if err != nil {
		return nil, err
	}
	return NewBalanceTopicsResponse(out), nil
}

func (h handler) ListTopics(req []byte) ([]byte, error) {
	v, err := WrapListTopicsRequest(req)
	if err != nil {
		return nil, err
	}
	out, err := h.b.ListTopics(v)
	if err != nil {
		return nil, err
	}
	return NewListTopicsResponse(out), nil
}

func (h handler) TopicExists(req []byte) ([]byte, error) {
	v, err := WrapTopicExistsRequest(req)
	if err != nil {
		return nil, err
	}
	out, err := h.b.TopicExists(v)
	if err != nil {
		return nil, err
	}
	return NewTopicExistsResponse(out), nil
}

func (h handler) ConfigureTopic(req []byte) ([]byte, error) {
	v, err := WrapConfigureTopicRequest(req)
	if err != nil {
		return nil, err
	}
	out, err := h.b.ConfigureTopic(v)
	if err != nil {
		return nil, err
	}
	return NewConfigureTopicResponse(out), nil
}

func (h handler) LookupTopicBrokers(req []byte) ([]byte, error) {
	v, err := WrapLookupTopicBrokersRequest(req)
	if err != nil {
		return nil, err
	}
	out, err := h.b.LookupTopicBrokers(v)
	if err != nil {
		return nil, err
	}
	return NewLookupTopicBrokersResponse(out), nil
}

func (h handler) GetTopicConfiguration(req []byte) ([]byte, error) {
	v, err := WrapGetTopicConfigurationRequest(req)
	if err != nil {
		return nil, err
	}
	out, err := h.b.GetTopicConfiguration(v)
	if err != nil {
		return nil, err
	}
	return NewGetTopicConfigurationResponse(out), nil
}

func (h handler) GetTopicPublishers(req []byte) ([]byte, error) {
	v, err := WrapGetTopicPublishersRequest(req)
	if err != nil {
		return nil, err
	}
	out, err := h.b.GetTopicPublishers(v)
	if err != nil {
		return nil, err
	}
	return NewGetTopicPublishersResponse(out), nil
}

func (h handler) GetTopicSubscribers(req []byte) ([]byte, error) {
	v, err := WrapGetTopicSubscribersRequest(req)
	if err != nil {
		return nil, err
	}
	out, err := h.b.GetTopicSubscribers(v)
	if err != nil {
		return nil, err
	}
	return NewGetTopicSubscribersResponse(out), nil
}

func (h handler) AssignTopicPartitions(req []byte) ([]byte, error) {
	v, err := WrapAssignTopicPartitionsRequest(req)
	if err != nil {
		return nil, err
	}
	out, err := h.b.AssignTopicPartitions(v)
	if err != nil {
		return nil, err
	}
	return NewAssignTopicPartitionsResponse(out), nil
}

func (h handler) ClosePublishers(req []byte) ([]byte, error) {
	v, err := WrapClosePublishersRequest(req)
	if err != nil {
		return nil, err
	}
	out, err := h.b.ClosePublishers(v)
	if err != nil {
		return nil, err
	}
	return NewClosePublishersResponse(out), nil
}

func (h handler) CloseSubscribers(req []byte) ([]byte, error) {
	v, err := WrapCloseSubscribersRequest(req)
	if err != nil {
		return nil, err
	}
	out, err := h.b.CloseSubscribers(v)
	if err != nil {
		return nil, err
	}
	return NewCloseSubscribersResponse(out), nil
}

func (h handler) FetchMessage(req []byte) ([]byte, error) {
	v, err := WrapFetchMessageRequest(req)
	if err != nil {
		return nil, err
	}
	out, err := h.b.FetchMessage(v)
	if err != nil {
		return nil, err
	}
	return NewFetchMessageResponse(out), nil
}

func (h handler) GetPartitionRangeInfo(req []byte) ([]byte, error) {
	v, err := WrapGetPartitionRangeInfoRequest(req)
	if err != nil {
		return nil, err
	}
	out, err := h.b.GetPartitionRangeInfo(v)
	if err != nil {
		return nil, err
	}
	return NewGetPartitionRangeInfoResponse(out), nil
}

// --- streaming methods of the generated interface ---
//
// These are served by ServeStream over a transport.Stream, never by a unary
// Call. The generated DispatchHanzoMessaging still has an arm for each (it cannot
// know a method is streaming), so we satisfy the interface with a no-op that
// returns an empty response buffer. A single-shot caller therefore gets a valid
// empty envelope rather than a decode error; real traffic always uses the stream.

func (h handler) PublisherToPubBalancer(req []byte) ([]byte, error) {
	return NewPublisherToPubBalancerResponse(PublisherToPubBalancerResponseInput{}), nil
}

func (h handler) SubscriberToSubCoordinator(req []byte) ([]byte, error) {
	return NewSubscriberToSubCoordinatorResponse(SubscriberToSubCoordinatorResponseInput{}), nil
}

func (h handler) PublishMessage(req []byte) ([]byte, error) {
	return NewPublishMessageResponse(PublishMessageResponseInput{}), nil
}

func (h handler) SubscribeMessage(req []byte) ([]byte, error) {
	return NewSubscribeMessageResponse(SubscribeMessageResponseInput{}), nil
}

func (h handler) PublishFollowMe(req []byte) ([]byte, error) {
	return NewPublishFollowMeResponse(PublishFollowMeResponseInput{}), nil
}

func (h handler) SubscribeFollowMe(req []byte) ([]byte, error) {
	return NewSubscribeFollowMeResponse(SubscribeFollowMeResponseInput{}), nil
}

func (h handler) GetUnflushedMessages(req []byte) ([]byte, error) {
	return NewGetUnflushedMessagesResponse(GetUnflushedMessagesResponseInput{}), nil
}

// Dispatch is the HanzoMessaging unary server dispatch bound to b — pass it to
// transport.Listen or as the unary arm of transport.ListenStream. Exposed so
// callers can compose it (e.g. behind PQ-TLS).
func Dispatch(b Broker) transport.Dispatch {
	h := handler{b: b}
	return func(envelope []byte) ([]byte, error) {
		return DispatchHanzoMessaging(h, envelope)
	}
}

// Serve starts the native ZAP HanzoMessaging service on network/addr (e.g.
// "tcp", ":17777"), backed by b for the 14 unary RPCs and stream for the 7
// streaming RPCs. stream may be nil if the deployment serves unary-only. Returns
// the running server; Close stops it. For the PQ-secured mesh, establish the
// listener via the TLS/QUIC path and use transport.ListenStream on it (or
// compose Dispatch behind a transport.PQTLSConfig listener).
func Serve(network, addr string, b Broker, stream StreamBroker) (*transport.Server, error) {
	var sh transport.StreamHandler
	if stream != nil {
		sh = stream.ServeStream
	}
	return transport.ListenStream(network, addr, Dispatch(b), sh)
}

// Client is the typed HanzoMessaging ZAP service client internal callers (the
// publisher / subscriber / balancer / Kafka gateway) hold. It owns the transport
// connection to one broker endpoint. Unary RPCs return decoded views; streaming
// RPCs return a live transport.Stream to drive.
type Client struct {
	conn transport.Conn
	rpc  *HanzoMessagingClient
}

// Dial connects to the HanzoMessaging ZAP service at addr (e.g.
// "broker.hanzo.svc:17777") over plain TCP. For the PQ-secured mesh, establish a
// transport.Conn via transport.DialTLS with a transport.PQTLSConfig and use
// NewClient.
func Dial(network, addr string) (*Client, error) {
	conn, err := transport.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn, rpc: NewHanzoMessagingClient(conn, nil)}, nil
}

// NewClient wraps an already-established transport.Conn (TCP, Unix, or PQ-TLS).
func NewClient(conn transport.Conn) *Client {
	return &Client{conn: conn, rpc: NewHanzoMessagingClient(conn, nil)}
}

// Conn exposes the underlying transport connection so callers can open streams
// for the seven streaming RPCs (see the Open* helpers below) or share it.
func (c *Client) Conn() transport.Conn { return c.conn }

// Close releases the underlying connection.
func (c *Client) Close() error { return c.conn.Close() }

// --- unary client calls (one per unary RPC) ---

// FindBrokerLeader issues the FindBrokerLeader RPC and returns the decoded view.
func (c *Client) FindBrokerLeader(in FindBrokerLeaderRequestInput) (FindBrokerLeaderResponse, error) {
	_, body, err := c.rpc.FindBrokerLeader(NewFindBrokerLeaderRequest(in))
	if err != nil {
		return FindBrokerLeaderResponse{}, err
	}
	return WrapFindBrokerLeaderResponse(body)
}

// BalanceTopics issues the BalanceTopics RPC and returns the decoded view.
func (c *Client) BalanceTopics(in BalanceTopicsRequestInput) (BalanceTopicsResponse, error) {
	_, body, err := c.rpc.BalanceTopics(NewBalanceTopicsRequest(in))
	if err != nil {
		return BalanceTopicsResponse{}, err
	}
	return WrapBalanceTopicsResponse(body)
}

// ListTopics issues the ListTopics RPC and returns the decoded view.
func (c *Client) ListTopics(in ListTopicsRequestInput) (ListTopicsResponse, error) {
	_, body, err := c.rpc.ListTopics(NewListTopicsRequest(in))
	if err != nil {
		return ListTopicsResponse{}, err
	}
	return WrapListTopicsResponse(body)
}

// TopicExists issues the TopicExists RPC and returns the decoded view.
func (c *Client) TopicExists(in TopicExistsRequestInput) (TopicExistsResponse, error) {
	_, body, err := c.rpc.TopicExists(NewTopicExistsRequest(in))
	if err != nil {
		return TopicExistsResponse{}, err
	}
	return WrapTopicExistsResponse(body)
}

// ConfigureTopic issues the ConfigureTopic RPC and returns the decoded view.
func (c *Client) ConfigureTopic(in ConfigureTopicRequestInput) (ConfigureTopicResponse, error) {
	_, body, err := c.rpc.ConfigureTopic(NewConfigureTopicRequest(in))
	if err != nil {
		return ConfigureTopicResponse{}, err
	}
	return WrapConfigureTopicResponse(body)
}

// LookupTopicBrokers issues the LookupTopicBrokers RPC and returns the view.
func (c *Client) LookupTopicBrokers(in LookupTopicBrokersRequestInput) (LookupTopicBrokersResponse, error) {
	_, body, err := c.rpc.LookupTopicBrokers(NewLookupTopicBrokersRequest(in))
	if err != nil {
		return LookupTopicBrokersResponse{}, err
	}
	return WrapLookupTopicBrokersResponse(body)
}

// GetTopicConfiguration issues the GetTopicConfiguration RPC and returns the view.
func (c *Client) GetTopicConfiguration(in GetTopicConfigurationRequestInput) (GetTopicConfigurationResponse, error) {
	_, body, err := c.rpc.GetTopicConfiguration(NewGetTopicConfigurationRequest(in))
	if err != nil {
		return GetTopicConfigurationResponse{}, err
	}
	return WrapGetTopicConfigurationResponse(body)
}

// GetTopicPublishers issues the GetTopicPublishers RPC and returns the view.
func (c *Client) GetTopicPublishers(in GetTopicPublishersRequestInput) (GetTopicPublishersResponse, error) {
	_, body, err := c.rpc.GetTopicPublishers(NewGetTopicPublishersRequest(in))
	if err != nil {
		return GetTopicPublishersResponse{}, err
	}
	return WrapGetTopicPublishersResponse(body)
}

// GetTopicSubscribers issues the GetTopicSubscribers RPC and returns the view.
func (c *Client) GetTopicSubscribers(in GetTopicSubscribersRequestInput) (GetTopicSubscribersResponse, error) {
	_, body, err := c.rpc.GetTopicSubscribers(NewGetTopicSubscribersRequest(in))
	if err != nil {
		return GetTopicSubscribersResponse{}, err
	}
	return WrapGetTopicSubscribersResponse(body)
}

// AssignTopicPartitions issues the AssignTopicPartitions RPC and returns the view.
func (c *Client) AssignTopicPartitions(in AssignTopicPartitionsRequestInput) (AssignTopicPartitionsResponse, error) {
	_, body, err := c.rpc.AssignTopicPartitions(NewAssignTopicPartitionsRequest(in))
	if err != nil {
		return AssignTopicPartitionsResponse{}, err
	}
	return WrapAssignTopicPartitionsResponse(body)
}

// ClosePublishers issues the ClosePublishers RPC and returns the decoded view.
func (c *Client) ClosePublishers(in ClosePublishersRequestInput) (ClosePublishersResponse, error) {
	_, body, err := c.rpc.ClosePublishers(NewClosePublishersRequest(in))
	if err != nil {
		return ClosePublishersResponse{}, err
	}
	return WrapClosePublishersResponse(body)
}

// CloseSubscribers issues the CloseSubscribers RPC and returns the decoded view.
func (c *Client) CloseSubscribers(in CloseSubscribersRequestInput) (CloseSubscribersResponse, error) {
	_, body, err := c.rpc.CloseSubscribers(NewCloseSubscribersRequest(in))
	if err != nil {
		return CloseSubscribersResponse{}, err
	}
	return WrapCloseSubscribersResponse(body)
}

// FetchMessage issues the stateless FetchMessage RPC and returns the decoded view.
func (c *Client) FetchMessage(in FetchMessageRequestInput) (FetchMessageResponse, error) {
	_, body, err := c.rpc.FetchMessage(NewFetchMessageRequest(in))
	if err != nil {
		return FetchMessageResponse{}, err
	}
	return WrapFetchMessageResponse(body)
}

// GetPartitionRangeInfo issues the GetPartitionRangeInfo RPC and returns the view.
func (c *Client) GetPartitionRangeInfo(in GetPartitionRangeInfoRequestInput) (GetPartitionRangeInfoResponse, error) {
	_, body, err := c.rpc.GetPartitionRangeInfo(NewGetPartitionRangeInfoRequest(in))
	if err != nil {
		return GetPartitionRangeInfoResponse{}, err
	}
	return WrapGetPartitionRangeInfoResponse(body)
}

// --- streaming client opens (one per streaming RPC) ---
//
// Each returns a live transport.Stream the caller drives with Send/Recv/
// CloseSend. init is the opening payload (a New<Req> buffer carrying the
// InitMessage variant); subsequent frames are New<Msg> buffers, read back with
// Wrap<Msg>. The server's StreamBroker.ServeStream is invoked with the matching
// ordinal.

// OpenPublisherToPubBalancer opens the PublisherToPubBalancer bidi stream.
func (c *Client) OpenPublisherToPubBalancer(init []byte) (transport.Stream, error) {
	return c.conn.OpenStream(HanzoMessagingPublisherToPubBalancerOrdinal, init)
}

// OpenSubscriberToSubCoordinator opens the SubscriberToSubCoordinator bidi stream.
func (c *Client) OpenSubscriberToSubCoordinator(init []byte) (transport.Stream, error) {
	return c.conn.OpenStream(HanzoMessagingSubscriberToSubCoordinatorOrdinal, init)
}

// OpenPublishMessage opens the PublishMessage bidi stream (data plane).
func (c *Client) OpenPublishMessage(init []byte) (transport.Stream, error) {
	return c.conn.OpenStream(HanzoMessagingPublishMessageOrdinal, init)
}

// OpenSubscribeMessage opens the SubscribeMessage bidi stream (data plane).
func (c *Client) OpenSubscribeMessage(init []byte) (transport.Stream, error) {
	return c.conn.OpenStream(HanzoMessagingSubscribeMessageOrdinal, init)
}

// OpenPublishFollowMe opens the PublishFollowMe bidi stream (lead -> follower).
func (c *Client) OpenPublishFollowMe(init []byte) (transport.Stream, error) {
	return c.conn.OpenStream(HanzoMessagingPublishFollowMeOrdinal, init)
}

// OpenSubscribeFollowMe opens the SubscribeFollowMe client-stream.
func (c *Client) OpenSubscribeFollowMe(init []byte) (transport.Stream, error) {
	return c.conn.OpenStream(HanzoMessagingSubscribeFollowMeOrdinal, init)
}

// OpenGetUnflushedMessages opens the GetUnflushedMessages server-stream (SQL).
func (c *Client) OpenGetUnflushedMessages(init []byte) (transport.Stream, error) {
	return c.conn.OpenStream(HanzoMessagingGetUnflushedMessagesOrdinal, init)
}
