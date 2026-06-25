// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// server.go is the server-side mirror of client.go: a mq_brokerwire.Broker that
// wraps the existing gRPC-shaped broker (mq_pb.HanzoMessagingServer — the real
// *broker.MessageQueueBroker) so it answers over the canonical
// github.com/zap-proto/go transport. Each unary RPC is one hop each way: the
// zero-copy request view -> <Rpc>ReqFromView -> the existing
// mq_pb.HanzoMessagingServer method -> <Rpc>RespToInput -> the matching
// mq_brokerwire.<Rpc>ResponseInput, which the wire layer's handler encodes back
// to bytes. With NewServerBackend wired into mq_brokerwire.Dispatch (the ZAP
// analogue of the retired gRPC HanzoMessaging registration) and NewStreamServer
// into the stream handler, the whole MQ broker (14 unary + 7 streaming RPCs)
// answers over ZAP — no protobuf framing on the wire, all engine logic reused.
//
// Unary RPCs live here; the 7 streaming RPCs are in stream_server.go.

package mq

import (
	"context"

	mq_pb "github.com/hanzoai/s3/s3/pb/mq_pb"
	mq_brokerwire "github.com/hanzoai/s3/s3/wire/mq_broker"
)

// serverBackend adapts a mq_pb.HanzoMessagingServer to mq_brokerwire.Broker.
type serverBackend struct {
	srv mq_pb.HanzoMessagingServer
	ctx context.Context
}

// NewServerBackend returns a mq_brokerwire.Broker that serves srv's 14 unary RPCs
// over ZAP. Pass it to mq_brokerwire.Dispatch / mq_brokerwire.Serve. For the
// streaming RPCs pair it with NewStreamServer.
func NewServerBackend(srv mq_pb.HanzoMessagingServer) mq_brokerwire.Broker {
	return serverBackend{srv: srv, ctx: context.Background()}
}

var _ mq_brokerwire.Broker = serverBackend{}

func (b serverBackend) FindBrokerLeader(v mq_brokerwire.FindBrokerLeaderRequest) (mq_brokerwire.FindBrokerLeaderResponseInput, error) {
	resp, err := b.srv.FindBrokerLeader(b.ctx, FindBrokerLeaderReqFromView(v))
	if err != nil {
		return mq_brokerwire.FindBrokerLeaderResponseInput{}, err
	}
	return FindBrokerLeaderRespToInput(resp), nil
}

func (b serverBackend) BalanceTopics(v mq_brokerwire.BalanceTopicsRequest) (mq_brokerwire.BalanceTopicsResponseInput, error) {
	resp, err := b.srv.BalanceTopics(b.ctx, BalanceTopicsReqFromView(v))
	if err != nil {
		return mq_brokerwire.BalanceTopicsResponseInput{}, err
	}
	return BalanceTopicsRespToInput(resp), nil
}

func (b serverBackend) ListTopics(v mq_brokerwire.ListTopicsRequest) (mq_brokerwire.ListTopicsResponseInput, error) {
	resp, err := b.srv.ListTopics(b.ctx, ListTopicsReqFromView(v))
	if err != nil {
		return mq_brokerwire.ListTopicsResponseInput{}, err
	}
	return ListTopicsRespToInput(resp), nil
}

func (b serverBackend) TopicExists(v mq_brokerwire.TopicExistsRequest) (mq_brokerwire.TopicExistsResponseInput, error) {
	resp, err := b.srv.TopicExists(b.ctx, TopicExistsReqFromView(v))
	if err != nil {
		return mq_brokerwire.TopicExistsResponseInput{}, err
	}
	return TopicExistsRespToInput(resp), nil
}

func (b serverBackend) ConfigureTopic(v mq_brokerwire.ConfigureTopicRequest) (mq_brokerwire.ConfigureTopicResponseInput, error) {
	resp, err := b.srv.ConfigureTopic(b.ctx, ConfigureTopicReqFromView(v))
	if err != nil {
		return mq_brokerwire.ConfigureTopicResponseInput{}, err
	}
	return ConfigureTopicRespToInput(resp), nil
}

func (b serverBackend) LookupTopicBrokers(v mq_brokerwire.LookupTopicBrokersRequest) (mq_brokerwire.LookupTopicBrokersResponseInput, error) {
	resp, err := b.srv.LookupTopicBrokers(b.ctx, LookupTopicBrokersReqFromView(v))
	if err != nil {
		return mq_brokerwire.LookupTopicBrokersResponseInput{}, err
	}
	return LookupTopicBrokersRespToInput(resp), nil
}

func (b serverBackend) GetTopicConfiguration(v mq_brokerwire.GetTopicConfigurationRequest) (mq_brokerwire.GetTopicConfigurationResponseInput, error) {
	resp, err := b.srv.GetTopicConfiguration(b.ctx, GetTopicConfigurationReqFromView(v))
	if err != nil {
		return mq_brokerwire.GetTopicConfigurationResponseInput{}, err
	}
	return GetTopicConfigurationRespToInput(resp), nil
}

func (b serverBackend) GetTopicPublishers(v mq_brokerwire.GetTopicPublishersRequest) (mq_brokerwire.GetTopicPublishersResponseInput, error) {
	resp, err := b.srv.GetTopicPublishers(b.ctx, GetTopicPublishersReqFromView(v))
	if err != nil {
		return mq_brokerwire.GetTopicPublishersResponseInput{}, err
	}
	return GetTopicPublishersRespToInput(resp), nil
}

func (b serverBackend) GetTopicSubscribers(v mq_brokerwire.GetTopicSubscribersRequest) (mq_brokerwire.GetTopicSubscribersResponseInput, error) {
	resp, err := b.srv.GetTopicSubscribers(b.ctx, GetTopicSubscribersReqFromView(v))
	if err != nil {
		return mq_brokerwire.GetTopicSubscribersResponseInput{}, err
	}
	return GetTopicSubscribersRespToInput(resp), nil
}

func (b serverBackend) AssignTopicPartitions(v mq_brokerwire.AssignTopicPartitionsRequest) (mq_brokerwire.AssignTopicPartitionsResponseInput, error) {
	resp, err := b.srv.AssignTopicPartitions(b.ctx, AssignTopicPartitionsReqFromView(v))
	if err != nil {
		return mq_brokerwire.AssignTopicPartitionsResponseInput{}, err
	}
	return AssignTopicPartitionsRespToInput(resp), nil
}

func (b serverBackend) ClosePublishers(v mq_brokerwire.ClosePublishersRequest) (mq_brokerwire.ClosePublishersResponseInput, error) {
	resp, err := b.srv.ClosePublishers(b.ctx, ClosePublishersReqFromView(v))
	if err != nil {
		return mq_brokerwire.ClosePublishersResponseInput{}, err
	}
	return ClosePublishersRespToInput(resp), nil
}

func (b serverBackend) CloseSubscribers(v mq_brokerwire.CloseSubscribersRequest) (mq_brokerwire.CloseSubscribersResponseInput, error) {
	resp, err := b.srv.CloseSubscribers(b.ctx, CloseSubscribersReqFromView(v))
	if err != nil {
		return mq_brokerwire.CloseSubscribersResponseInput{}, err
	}
	return CloseSubscribersRespToInput(resp), nil
}

func (b serverBackend) FetchMessage(v mq_brokerwire.FetchMessageRequest) (mq_brokerwire.FetchMessageResponseInput, error) {
	resp, err := b.srv.FetchMessage(b.ctx, FetchMessageReqFromView(v))
	if err != nil {
		return mq_brokerwire.FetchMessageResponseInput{}, err
	}
	return FetchMessageRespToInput(resp), nil
}

func (b serverBackend) GetPartitionRangeInfo(v mq_brokerwire.GetPartitionRangeInfoRequest) (mq_brokerwire.GetPartitionRangeInfoResponseInput, error) {
	resp, err := b.srv.GetPartitionRangeInfo(b.ctx, GetPartitionRangeInfoReqFromView(v))
	if err != nil {
		return mq_brokerwire.GetPartitionRangeInfoResponseInput{}, err
	}
	return GetPartitionRangeInfoRespToInput(resp), nil
}
