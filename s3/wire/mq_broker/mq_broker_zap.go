// Code generated from mq_broker.proto; DO NOT EDIT.

package mq_brokerwire

import (
	"fmt"

	"github.com/zap-proto/go/rpc"
)

// Method ordinals for the HanzoMessaging service (stable 1-based wire ids,
// in mq_broker.proto declaration order).
const (
	HanzoMessagingFindBrokerLeaderOrdinal           uint32 = 1
	HanzoMessagingPublisherToPubBalancerOrdinal     uint32 = 2
	HanzoMessagingBalanceTopicsOrdinal              uint32 = 3
	HanzoMessagingListTopicsOrdinal                 uint32 = 4
	HanzoMessagingTopicExistsOrdinal                uint32 = 5
	HanzoMessagingConfigureTopicOrdinal             uint32 = 6
	HanzoMessagingLookupTopicBrokersOrdinal         uint32 = 7
	HanzoMessagingGetTopicConfigurationOrdinal      uint32 = 8
	HanzoMessagingGetTopicPublishersOrdinal         uint32 = 9
	HanzoMessagingGetTopicSubscribersOrdinal        uint32 = 10
	HanzoMessagingAssignTopicPartitionsOrdinal      uint32 = 11
	HanzoMessagingClosePublishersOrdinal            uint32 = 12
	HanzoMessagingCloseSubscribersOrdinal           uint32 = 13
	HanzoMessagingSubscriberToSubCoordinatorOrdinal uint32 = 14
	HanzoMessagingPublishMessageOrdinal             uint32 = 15
	HanzoMessagingSubscribeMessageOrdinal           uint32 = 16
	HanzoMessagingPublishFollowMeOrdinal            uint32 = 17
	HanzoMessagingSubscribeFollowMeOrdinal          uint32 = 18
	HanzoMessagingFetchMessageOrdinal               uint32 = 19
	HanzoMessagingGetUnflushedMessagesOrdinal       uint32 = 20
	HanzoMessagingGetPartitionRangeInfoOrdinal      uint32 = 21
)

// HanzoMessagingChannel ships one Call envelope and awaits its correlated
// Response.
type HanzoMessagingChannel interface {
	Call(envelope []byte) (rpc.Response, error)
	// NextPromiseID allocates a call's PromiseID from the CONNECTION rather than
	// from a per-client counter. PromiseID is the connection's demultiplexing
	// key: the transport keys its in-flight table by it, and OpenStream already
	// allocates stream IDs the same way. A client that numbered its own calls
	// would restart at 1 on every construction, so two clients sharing a pooled
	// conn collide — the transport overwrites the first waiter's slot and its
	// response is dropped, blocking that caller for as long as its context runs.
	NextPromiseID() uint32
}

// HanzoMessagingClient is a typed RPC client for the HanzoMessaging service over
// a ZAP call channel. Each call takes a fresh PromiseID from sess; the pipelined
// "On" form of a method sets Target to a prior call's Promise so the server
// chains them.
type HanzoMessagingClient struct {
	ch   HanzoMessagingChannel
	cap  []byte
}

// NewHanzoMessagingClient returns a client that issues calls over ch, attaching
// cap (which may be nil) to every request.
func NewHanzoMessagingClient(ch HanzoMessagingChannel, capability []byte) *HanzoMessagingClient {
	return &HanzoMessagingClient{ch: ch, cap: capability}
}

// invoke is the shared call path: it takes a fresh promise, ships the request,
// and returns the response body (or a status error). Every per-method wrapper
// funnels through here so the call shape stays in one place.
func (c *HanzoMessagingClient) invoke(method, target uint32, payload []byte, name string) (rpc.Promise, []byte, error) {
	p := rpc.Promise{ID: c.ch.NextPromiseID()}
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    method,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels.
			return p, nil, fmt.Errorf("HanzoMessaging.%s: %s", name, resp.Body)
		}
		return p, nil, fmt.Errorf("HanzoMessaging.%s: status %d", name, resp.Status)
	}
	return p, resp.Body, nil
}

// FindBrokerLeader issues the FindBrokerLeader RPC (unary).
func (c *HanzoMessagingClient) FindBrokerLeader(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingFindBrokerLeaderOrdinal, rpc.NoTarget, req, "FindBrokerLeader")
}

// FindBrokerLeaderOn issues FindBrokerLeader pipelined on the answer of on.
func (c *HanzoMessagingClient) FindBrokerLeaderOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingFindBrokerLeaderOrdinal, on.ID, nil, "FindBrokerLeader")
}

// PublisherToPubBalancer issues the PublisherToPubBalancer RPC.
//
// STREAMING: bidi-streaming (rpc PublisherToPubBalancer(stream
// PublisherToPubBalancerRequest) returns (stream PublisherToPubBalancerResponse)).
// The single-shot form ships one request envelope and returns one response
// envelope; the full duplex stream body lands when the transport streaming
// primitive ships.
func (c *HanzoMessagingClient) PublisherToPubBalancer(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingPublisherToPubBalancerOrdinal, rpc.NoTarget, req, "PublisherToPubBalancer")
}

// PublisherToPubBalancerOn issues PublisherToPubBalancer pipelined on the answer
// of on.
//
// STREAMING: see PublisherToPubBalancer.
func (c *HanzoMessagingClient) PublisherToPubBalancerOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingPublisherToPubBalancerOrdinal, on.ID, nil, "PublisherToPubBalancer")
}

// BalanceTopics issues the BalanceTopics RPC (unary).
func (c *HanzoMessagingClient) BalanceTopics(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingBalanceTopicsOrdinal, rpc.NoTarget, req, "BalanceTopics")
}

// BalanceTopicsOn issues BalanceTopics pipelined on the answer of on.
func (c *HanzoMessagingClient) BalanceTopicsOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingBalanceTopicsOrdinal, on.ID, nil, "BalanceTopics")
}

// ListTopics issues the ListTopics RPC (unary).
func (c *HanzoMessagingClient) ListTopics(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingListTopicsOrdinal, rpc.NoTarget, req, "ListTopics")
}

// ListTopicsOn issues ListTopics pipelined on the answer of on.
func (c *HanzoMessagingClient) ListTopicsOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingListTopicsOrdinal, on.ID, nil, "ListTopics")
}

// TopicExists issues the TopicExists RPC (unary).
func (c *HanzoMessagingClient) TopicExists(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingTopicExistsOrdinal, rpc.NoTarget, req, "TopicExists")
}

// TopicExistsOn issues TopicExists pipelined on the answer of on.
func (c *HanzoMessagingClient) TopicExistsOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingTopicExistsOrdinal, on.ID, nil, "TopicExists")
}

// ConfigureTopic issues the ConfigureTopic RPC (unary).
func (c *HanzoMessagingClient) ConfigureTopic(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingConfigureTopicOrdinal, rpc.NoTarget, req, "ConfigureTopic")
}

// ConfigureTopicOn issues ConfigureTopic pipelined on the answer of on.
func (c *HanzoMessagingClient) ConfigureTopicOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingConfigureTopicOrdinal, on.ID, nil, "ConfigureTopic")
}

// LookupTopicBrokers issues the LookupTopicBrokers RPC (unary).
func (c *HanzoMessagingClient) LookupTopicBrokers(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingLookupTopicBrokersOrdinal, rpc.NoTarget, req, "LookupTopicBrokers")
}

// LookupTopicBrokersOn issues LookupTopicBrokers pipelined on the answer of on.
func (c *HanzoMessagingClient) LookupTopicBrokersOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingLookupTopicBrokersOrdinal, on.ID, nil, "LookupTopicBrokers")
}

// GetTopicConfiguration issues the GetTopicConfiguration RPC (unary).
func (c *HanzoMessagingClient) GetTopicConfiguration(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingGetTopicConfigurationOrdinal, rpc.NoTarget, req, "GetTopicConfiguration")
}

// GetTopicConfigurationOn issues GetTopicConfiguration pipelined on the answer
// of on.
func (c *HanzoMessagingClient) GetTopicConfigurationOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingGetTopicConfigurationOrdinal, on.ID, nil, "GetTopicConfiguration")
}

// GetTopicPublishers issues the GetTopicPublishers RPC (unary).
func (c *HanzoMessagingClient) GetTopicPublishers(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingGetTopicPublishersOrdinal, rpc.NoTarget, req, "GetTopicPublishers")
}

// GetTopicPublishersOn issues GetTopicPublishers pipelined on the answer of on.
func (c *HanzoMessagingClient) GetTopicPublishersOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingGetTopicPublishersOrdinal, on.ID, nil, "GetTopicPublishers")
}

// GetTopicSubscribers issues the GetTopicSubscribers RPC (unary).
func (c *HanzoMessagingClient) GetTopicSubscribers(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingGetTopicSubscribersOrdinal, rpc.NoTarget, req, "GetTopicSubscribers")
}

// GetTopicSubscribersOn issues GetTopicSubscribers pipelined on the answer of
// on.
func (c *HanzoMessagingClient) GetTopicSubscribersOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingGetTopicSubscribersOrdinal, on.ID, nil, "GetTopicSubscribers")
}

// AssignTopicPartitions issues the AssignTopicPartitions RPC (unary).
func (c *HanzoMessagingClient) AssignTopicPartitions(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingAssignTopicPartitionsOrdinal, rpc.NoTarget, req, "AssignTopicPartitions")
}

// AssignTopicPartitionsOn issues AssignTopicPartitions pipelined on the answer
// of on.
func (c *HanzoMessagingClient) AssignTopicPartitionsOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingAssignTopicPartitionsOrdinal, on.ID, nil, "AssignTopicPartitions")
}

// ClosePublishers issues the ClosePublishers RPC (unary).
func (c *HanzoMessagingClient) ClosePublishers(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingClosePublishersOrdinal, rpc.NoTarget, req, "ClosePublishers")
}

// ClosePublishersOn issues ClosePublishers pipelined on the answer of on.
func (c *HanzoMessagingClient) ClosePublishersOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingClosePublishersOrdinal, on.ID, nil, "ClosePublishers")
}

// CloseSubscribers issues the CloseSubscribers RPC (unary).
func (c *HanzoMessagingClient) CloseSubscribers(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingCloseSubscribersOrdinal, rpc.NoTarget, req, "CloseSubscribers")
}

// CloseSubscribersOn issues CloseSubscribers pipelined on the answer of on.
func (c *HanzoMessagingClient) CloseSubscribersOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingCloseSubscribersOrdinal, on.ID, nil, "CloseSubscribers")
}

// SubscriberToSubCoordinator issues the SubscriberToSubCoordinator RPC.
//
// STREAMING: bidi-streaming (rpc SubscriberToSubCoordinator(stream
// SubscriberToSubCoordinatorRequest) returns (stream
// SubscriberToSubCoordinatorResponse)). The single-shot form ships one request
// envelope and returns one response envelope; the full duplex stream body lands
// when the transport streaming primitive ships.
func (c *HanzoMessagingClient) SubscriberToSubCoordinator(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingSubscriberToSubCoordinatorOrdinal, rpc.NoTarget, req, "SubscriberToSubCoordinator")
}

// SubscriberToSubCoordinatorOn issues SubscriberToSubCoordinator pipelined on
// the answer of on.
//
// STREAMING: see SubscriberToSubCoordinator.
func (c *HanzoMessagingClient) SubscriberToSubCoordinatorOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingSubscriberToSubCoordinatorOrdinal, on.ID, nil, "SubscriberToSubCoordinator")
}

// PublishMessage issues the PublishMessage RPC.
//
// STREAMING: bidi-streaming (rpc PublishMessage(stream PublishMessageRequest)
// returns (stream PublishMessageResponse)). The single-shot form ships one
// request envelope and returns one response envelope; the full duplex stream
// body lands when the transport streaming primitive ships.
func (c *HanzoMessagingClient) PublishMessage(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingPublishMessageOrdinal, rpc.NoTarget, req, "PublishMessage")
}

// PublishMessageOn issues PublishMessage pipelined on the answer of on.
//
// STREAMING: see PublishMessage.
func (c *HanzoMessagingClient) PublishMessageOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingPublishMessageOrdinal, on.ID, nil, "PublishMessage")
}

// SubscribeMessage issues the SubscribeMessage RPC.
//
// STREAMING: bidi-streaming (rpc SubscribeMessage(stream SubscribeMessageRequest)
// returns (stream SubscribeMessageResponse)). The single-shot form ships one
// request envelope and returns one response envelope; the full duplex stream
// body lands when the transport streaming primitive ships.
func (c *HanzoMessagingClient) SubscribeMessage(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingSubscribeMessageOrdinal, rpc.NoTarget, req, "SubscribeMessage")
}

// SubscribeMessageOn issues SubscribeMessage pipelined on the answer of on.
//
// STREAMING: see SubscribeMessage.
func (c *HanzoMessagingClient) SubscribeMessageOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingSubscribeMessageOrdinal, on.ID, nil, "SubscribeMessage")
}

// PublishFollowMe issues the PublishFollowMe RPC.
//
// STREAMING: bidi-streaming (rpc PublishFollowMe(stream PublishFollowMeRequest)
// returns (stream PublishFollowMeResponse)). The single-shot form ships one
// request envelope and returns one response envelope; the full duplex stream
// body lands when the transport streaming primitive ships.
func (c *HanzoMessagingClient) PublishFollowMe(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingPublishFollowMeOrdinal, rpc.NoTarget, req, "PublishFollowMe")
}

// PublishFollowMeOn issues PublishFollowMe pipelined on the answer of on.
//
// STREAMING: see PublishFollowMe.
func (c *HanzoMessagingClient) PublishFollowMeOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingPublishFollowMeOrdinal, on.ID, nil, "PublishFollowMe")
}

// SubscribeFollowMe issues the SubscribeFollowMe RPC.
//
// STREAMING: client-streaming (rpc SubscribeFollowMe(stream
// SubscribeFollowMeRequest) returns (SubscribeFollowMeResponse)). The single-shot
// form ships one request envelope and returns the single response envelope; the
// full client-stream body lands when the transport streaming primitive ships.
func (c *HanzoMessagingClient) SubscribeFollowMe(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingSubscribeFollowMeOrdinal, rpc.NoTarget, req, "SubscribeFollowMe")
}

// SubscribeFollowMeOn issues SubscribeFollowMe pipelined on the answer of on.
//
// STREAMING: see SubscribeFollowMe.
func (c *HanzoMessagingClient) SubscribeFollowMeOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingSubscribeFollowMeOrdinal, on.ID, nil, "SubscribeFollowMe")
}

// FetchMessage issues the FetchMessage RPC (unary).
func (c *HanzoMessagingClient) FetchMessage(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingFetchMessageOrdinal, rpc.NoTarget, req, "FetchMessage")
}

// FetchMessageOn issues FetchMessage pipelined on the answer of on.
func (c *HanzoMessagingClient) FetchMessageOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingFetchMessageOrdinal, on.ID, nil, "FetchMessage")
}

// GetUnflushedMessages issues the GetUnflushedMessages RPC.
//
// STREAMING: server-streaming (rpc GetUnflushedMessages(GetUnflushedMessagesRequest)
// returns (stream GetUnflushedMessagesResponse)). The single-shot form ships the
// request envelope and returns one response envelope; the full server-stream body
// lands when the transport streaming primitive ships.
func (c *HanzoMessagingClient) GetUnflushedMessages(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingGetUnflushedMessagesOrdinal, rpc.NoTarget, req, "GetUnflushedMessages")
}

// GetUnflushedMessagesOn issues GetUnflushedMessages pipelined on the answer of
// on.
//
// STREAMING: see GetUnflushedMessages.
func (c *HanzoMessagingClient) GetUnflushedMessagesOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingGetUnflushedMessagesOrdinal, on.ID, nil, "GetUnflushedMessages")
}

// GetPartitionRangeInfo issues the GetPartitionRangeInfo RPC (unary).
func (c *HanzoMessagingClient) GetPartitionRangeInfo(req []byte) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingGetPartitionRangeInfoOrdinal, rpc.NoTarget, req, "GetPartitionRangeInfo")
}

// GetPartitionRangeInfoOn issues GetPartitionRangeInfo pipelined on the answer
// of on.
func (c *HanzoMessagingClient) GetPartitionRangeInfoOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invoke(HanzoMessagingGetPartitionRangeInfoOrdinal, on.ID, nil, "GetPartitionRangeInfo")
}

// HanzoMessagingHandler is the server contract for the HanzoMessaging service.
// Implement each method, then route requests to it with DispatchHanzoMessaging.
//
// STREAMING: PublisherToPubBalancer, SubscriberToSubCoordinator, PublishMessage,
// SubscribeMessage and PublishFollowMe are bidi-streaming; SubscribeFollowMe is
// client-streaming; GetUnflushedMessages is server-streaming. Each handler here
// receives one request envelope and returns one response envelope. The full
// stream contracts land when the transport streaming primitive ships.
type HanzoMessagingHandler interface {
	FindBrokerLeader(req []byte) ([]byte, error)
	PublisherToPubBalancer(req []byte) ([]byte, error) // STREAMING (bidi)
	BalanceTopics(req []byte) ([]byte, error)
	ListTopics(req []byte) ([]byte, error)
	TopicExists(req []byte) ([]byte, error)
	ConfigureTopic(req []byte) ([]byte, error)
	LookupTopicBrokers(req []byte) ([]byte, error)
	GetTopicConfiguration(req []byte) ([]byte, error)
	GetTopicPublishers(req []byte) ([]byte, error)
	GetTopicSubscribers(req []byte) ([]byte, error)
	AssignTopicPartitions(req []byte) ([]byte, error)
	ClosePublishers(req []byte) ([]byte, error)
	CloseSubscribers(req []byte) ([]byte, error)
	SubscriberToSubCoordinator(req []byte) ([]byte, error) // STREAMING (bidi)
	PublishMessage(req []byte) ([]byte, error)             // STREAMING (bidi)
	SubscribeMessage(req []byte) ([]byte, error)           // STREAMING (bidi)
	PublishFollowMe(req []byte) ([]byte, error)            // STREAMING (bidi)
	SubscribeFollowMe(req []byte) ([]byte, error)          // STREAMING (client-stream)
	FetchMessage(req []byte) ([]byte, error)
	GetUnflushedMessages(req []byte) ([]byte, error) // STREAMING (server-stream)
	GetPartitionRangeInfo(req []byte) ([]byte, error)
}

// DispatchHanzoMessaging decodes a Call envelope, routes it by method ordinal to
// h, and returns the response envelope. An unknown ordinal yields a
// StatusNotFound response; a handler error yields StatusInternal.
func DispatchHanzoMessaging(h HanzoMessagingHandler, envelope []byte) ([]byte, error) {
	call, err := rpc.ParseRequest(envelope)
	if err != nil {
		return nil, err
	}
	switch call.Method {
	case HanzoMessagingFindBrokerLeaderOrdinal:
		return dispatchOne(call.PromiseID, h.FindBrokerLeader, call.Payload)
	case HanzoMessagingPublisherToPubBalancerOrdinal:
		// STREAMING (bidi): single-shot dispatch of one request -> one response.
		return dispatchOne(call.PromiseID, h.PublisherToPubBalancer, call.Payload)
	case HanzoMessagingBalanceTopicsOrdinal:
		return dispatchOne(call.PromiseID, h.BalanceTopics, call.Payload)
	case HanzoMessagingListTopicsOrdinal:
		return dispatchOne(call.PromiseID, h.ListTopics, call.Payload)
	case HanzoMessagingTopicExistsOrdinal:
		return dispatchOne(call.PromiseID, h.TopicExists, call.Payload)
	case HanzoMessagingConfigureTopicOrdinal:
		return dispatchOne(call.PromiseID, h.ConfigureTopic, call.Payload)
	case HanzoMessagingLookupTopicBrokersOrdinal:
		return dispatchOne(call.PromiseID, h.LookupTopicBrokers, call.Payload)
	case HanzoMessagingGetTopicConfigurationOrdinal:
		return dispatchOne(call.PromiseID, h.GetTopicConfiguration, call.Payload)
	case HanzoMessagingGetTopicPublishersOrdinal:
		return dispatchOne(call.PromiseID, h.GetTopicPublishers, call.Payload)
	case HanzoMessagingGetTopicSubscribersOrdinal:
		return dispatchOne(call.PromiseID, h.GetTopicSubscribers, call.Payload)
	case HanzoMessagingAssignTopicPartitionsOrdinal:
		return dispatchOne(call.PromiseID, h.AssignTopicPartitions, call.Payload)
	case HanzoMessagingClosePublishersOrdinal:
		return dispatchOne(call.PromiseID, h.ClosePublishers, call.Payload)
	case HanzoMessagingCloseSubscribersOrdinal:
		return dispatchOne(call.PromiseID, h.CloseSubscribers, call.Payload)
	case HanzoMessagingSubscriberToSubCoordinatorOrdinal:
		// STREAMING (bidi): single-shot dispatch of one request -> one response.
		return dispatchOne(call.PromiseID, h.SubscriberToSubCoordinator, call.Payload)
	case HanzoMessagingPublishMessageOrdinal:
		// STREAMING (bidi): single-shot dispatch of one request -> one response.
		return dispatchOne(call.PromiseID, h.PublishMessage, call.Payload)
	case HanzoMessagingSubscribeMessageOrdinal:
		// STREAMING (bidi): single-shot dispatch of one request -> one response.
		return dispatchOne(call.PromiseID, h.SubscribeMessage, call.Payload)
	case HanzoMessagingPublishFollowMeOrdinal:
		// STREAMING (bidi): single-shot dispatch of one request -> one response.
		return dispatchOne(call.PromiseID, h.PublishFollowMe, call.Payload)
	case HanzoMessagingSubscribeFollowMeOrdinal:
		// STREAMING (client-stream): single-shot dispatch of one request -> one response.
		return dispatchOne(call.PromiseID, h.SubscribeFollowMe, call.Payload)
	case HanzoMessagingFetchMessageOrdinal:
		return dispatchOne(call.PromiseID, h.FetchMessage, call.Payload)
	case HanzoMessagingGetUnflushedMessagesOrdinal:
		// STREAMING (server-stream): single-shot dispatch of one request -> one response.
		return dispatchOne(call.PromiseID, h.GetUnflushedMessages, call.Payload)
	case HanzoMessagingGetPartitionRangeInfoOrdinal:
		return dispatchOne(call.PromiseID, h.GetPartitionRangeInfo, call.Payload)
	default:
		return rpc.BuildResponse(rpc.StatusNotFound, call.PromiseID, nil), nil
	}
}

// dispatchOne runs a single request/response handler and maps its result onto a
// response envelope: a handler error becomes StatusInternal, success becomes
// StatusOK with the handler's body. Every dispatch arm funnels through here so
// the status mapping stays in one place.
func dispatchOne(promiseID uint32, fn func([]byte) ([]byte, error), payload []byte) ([]byte, error) {
	body, err := fn(payload)
	if err != nil {
		// Carry the handler error message so the caller can reconstruct sentinels.
		return rpc.BuildResponse(rpc.StatusInternal, promiseID, []byte(err.Error())), nil
	}
	return rpc.BuildResponse(rpc.StatusOK, promiseID, body), nil
}
