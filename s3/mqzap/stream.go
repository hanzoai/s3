// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// stream.go holds the per-frame request builders and response decoders for the
// HanzoMessaging STREAMING RPCs (PublisherToPubBalancer, SubscriberToSubCoordinator,
// PublishMessage, SubscribeMessage, PublishFollowMe, SubscribeFollowMe,
// GetUnflushedMessages). Each request/response is a oneof; the converter encodes
// the set variant into the matching mq_brokerwire <Msg>Input (or decodes the view
// back). Naming mirrors rpc.go: <Rpc>ReqToWire(*mq_pb.<Rpc>Request) []byte and
// <Rpc>RespFromWire([]byte) (*mq_pb.<Rpc>Response, error).

package mqzap

import (
	mq_pb "github.com/hanzoai/s3/s3/pb/mq_pb"
	mq_brokerwire "github.com/hanzoai/s3/s3/wire/mq_broker"
)

// PublisherToPubBalancer (bidi)

func PublisherToPubBalancerReqToWire(r *mq_pb.PublisherToPubBalancerRequest) []byte {
	in := mq_brokerwire.PublisherToPubBalancerRequestInput{MessageWhich: mq_brokerwire.PublisherToPubBalancerRequestMessageNone}
	switch m := r.Message.(type) {
	case *mq_pb.PublisherToPubBalancerRequest_Init:
		in.MessageWhich = mq_brokerwire.PublisherToPubBalancerRequestMessageInit
		in.MessageValue = mq_brokerwire.NewPublisherToPubBalancerRequestInitMessage(mq_brokerwire.PublisherToPubBalancerRequestInitMessageInput{Broker: m.Init.Broker})
	case *mq_pb.PublisherToPubBalancerRequest_Stats:
		in.MessageWhich = mq_brokerwire.PublisherToPubBalancerRequestMessageStats
		in.MessageValue = brokerStatsToWire(m.Stats)
	}
	return mq_brokerwire.NewPublisherToPubBalancerRequest(in)
}

func PublisherToPubBalancerRespFromWire(b []byte) (*mq_pb.PublisherToPubBalancerResponse, error) {
	if _, err := mq_brokerwire.WrapPublisherToPubBalancerResponse(b); err != nil {
		return nil, err
	}
	return &mq_pb.PublisherToPubBalancerResponse{}, nil
}

func brokerStatsToWire(s *mq_pb.BrokerStats) []byte {
	if s == nil {
		return nil
	}
	stats := make([]mq_brokerwire.StringMsgEntry, 0, len(s.Stats))
	for k, v := range s.Stats {
		stats = append(stats, mq_brokerwire.StringMsgEntry{Key: k, Value: topicPartitionStatsToWire(v)})
	}
	return mq_brokerwire.NewBrokerStats(mq_brokerwire.BrokerStatsInput{CpuUsagePercent: s.CpuUsagePercent, Stats: stats})
}

func topicPartitionStatsToWire(s *mq_pb.TopicPartitionStats) []byte {
	if s == nil {
		return nil
	}
	return mq_brokerwire.NewTopicPartitionStats(mq_brokerwire.TopicPartitionStatsInput{
		Topic:           topicToWire(s.Topic),
		Partition:       partitionToWire(s.Partition),
		PublisherCount:  s.PublisherCount,
		SubscriberCount: s.SubscriberCount,
		Follower:        s.Follower,
	})
}

// SubscriberToSubCoordinator (bidi)

func SubscriberToSubCoordinatorReqToWire(r *mq_pb.SubscriberToSubCoordinatorRequest) []byte {
	in := mq_brokerwire.SubscriberToSubCoordinatorRequestInput{MessageWhich: mq_brokerwire.SubscriberToSubCoordinatorRequestMessageNone}
	switch m := r.Message.(type) {
	case *mq_pb.SubscriberToSubCoordinatorRequest_Init:
		in.MessageWhich = mq_brokerwire.SubscriberToSubCoordinatorRequestMessageInit
		in.MessageValue = mq_brokerwire.NewSubscriberToSubCoordinatorRequestInitMessage(mq_brokerwire.SubscriberToSubCoordinatorRequestInitMessageInput{
			ConsumerGroup:           m.Init.ConsumerGroup,
			ConsumerGroupInstanceId: m.Init.ConsumerGroupInstanceId,
			Topic:                   topicToWire(m.Init.Topic),
			MaxPartitionCount:       m.Init.MaxPartitionCount,
			RebalanceSeconds:        m.Init.RebalanceSeconds,
		})
	case *mq_pb.SubscriberToSubCoordinatorRequest_AckAssignment:
		in.MessageWhich = mq_brokerwire.SubscriberToSubCoordinatorRequestMessageAckAssignment
		in.MessageValue = mq_brokerwire.NewSubscriberToSubCoordinatorRequestAckAssignmentMessage(mq_brokerwire.SubscriberToSubCoordinatorRequestAckAssignmentMessageInput{
			Partition: partitionToWire(m.AckAssignment.Partition),
		})
	case *mq_pb.SubscriberToSubCoordinatorRequest_AckUnAssignment:
		in.MessageWhich = mq_brokerwire.SubscriberToSubCoordinatorRequestMessageAckUnAssignment
		in.MessageValue = mq_brokerwire.NewSubscriberToSubCoordinatorRequestAckUnAssignmentMessage(mq_brokerwire.SubscriberToSubCoordinatorRequestAckUnAssignmentMessageInput{
			Partition: partitionToWire(m.AckUnAssignment.Partition),
		})
	}
	return mq_brokerwire.NewSubscriberToSubCoordinatorRequest(in)
}

func SubscriberToSubCoordinatorRespFromWire(b []byte) (*mq_pb.SubscriberToSubCoordinatorResponse, error) {
	v, err := mq_brokerwire.WrapSubscriberToSubCoordinatorResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &mq_pb.SubscriberToSubCoordinatorResponse{}
	switch v.WhichMessage() {
	case mq_brokerwire.SubscriberToSubCoordinatorResponseMessageAssignment:
		a := &mq_pb.SubscriberToSubCoordinatorResponse_Assignment{}
		if pa, ok := v.Assignment().PartitionAssignment(); ok {
			a.PartitionAssignment = assignmentFromView(pa)
		}
		resp.Message = &mq_pb.SubscriberToSubCoordinatorResponse_Assignment_{Assignment: a}
	case mq_brokerwire.SubscriberToSubCoordinatorResponseMessageUnAssignment:
		resp.Message = &mq_pb.SubscriberToSubCoordinatorResponse_UnAssignment_{
			UnAssignment: &mq_pb.SubscriberToSubCoordinatorResponse_UnAssignment{
				Partition: partitionFromWire(v.UnAssignment().Partition()),
			},
		}
	}
	return resp, nil
}

// PublishMessage (bidi)

func PublishMessageReqToWire(r *mq_pb.PublishMessageRequest) []byte {
	in := mq_brokerwire.PublishMessageRequestInput{MessageWhich: mq_brokerwire.PublishMessageRequestMessageNone}
	switch m := r.Message.(type) {
	case *mq_pb.PublishMessageRequest_Init:
		in.MessageWhich = mq_brokerwire.PublishMessageRequestMessageInit
		in.MessageValue = mq_brokerwire.NewPublishMessageRequestInitMessage(mq_brokerwire.PublishMessageRequestInitMessageInput{
			Topic:          topicToWire(m.Init.Topic),
			Partition:      partitionToWire(m.Init.Partition),
			AckInterval:    m.Init.AckInterval,
			FollowerBroker: m.Init.FollowerBroker,
			PublisherName:  m.Init.PublisherName,
		})
	case *mq_pb.PublishMessageRequest_Data:
		in.MessageWhich = mq_brokerwire.PublishMessageRequestMessageData
		in.MessageValue = dataMessageToWire(m.Data)
	}
	return mq_brokerwire.NewPublishMessageRequest(in)
}

func PublishMessageRespFromWire(b []byte) (*mq_pb.PublishMessageResponse, error) {
	v, err := mq_brokerwire.WrapPublishMessageResponse(b)
	if err != nil {
		return nil, err
	}
	return &mq_pb.PublishMessageResponse{
		AckTsNs:        v.AckTsNs(),
		Error:          v.Error(),
		ShouldClose:    v.ShouldClose(),
		ErrorCode:      v.ErrorCode(),
		AssignedOffset: v.AssignedOffset(),
	}, nil
}

// SubscribeMessage (bidi)

func SubscribeMessageReqToWire(r *mq_pb.SubscribeMessageRequest) []byte {
	in := mq_brokerwire.SubscribeMessageRequestInput{MessageWhich: mq_brokerwire.SubscribeMessageRequestMessageNone}
	switch m := r.Message.(type) {
	case *mq_pb.SubscribeMessageRequest_Init:
		in.MessageWhich = mq_brokerwire.SubscribeMessageRequestMessageInit
		in.MessageValue = mq_brokerwire.NewSubscribeMessageRequestInitMessage(mq_brokerwire.SubscribeMessageRequestInitMessageInput{
			ConsumerGroup:     m.Init.ConsumerGroup,
			ConsumerId:        m.Init.ConsumerId,
			ClientId:          m.Init.ClientId,
			Topic:             topicToWire(m.Init.Topic),
			PartitionOffset:   partitionOffsetToWire(m.Init.PartitionOffset),
			OffsetType:        uint32(m.Init.OffsetType),
			Filter:            m.Init.Filter,
			FollowerBroker:    m.Init.FollowerBroker,
			SlidingWindowSize: m.Init.SlidingWindowSize,
		})
	case *mq_pb.SubscribeMessageRequest_Ack:
		in.MessageWhich = mq_brokerwire.SubscribeMessageRequestMessageAck
		in.MessageValue = mq_brokerwire.NewSubscribeMessageRequestAckMessage(mq_brokerwire.SubscribeMessageRequestAckMessageInput{
			TsNs: m.Ack.TsNs, Key: m.Ack.Key,
		})
	case *mq_pb.SubscribeMessageRequest_Seek:
		in.MessageWhich = mq_brokerwire.SubscribeMessageRequestMessageSeek
		in.MessageValue = mq_brokerwire.NewSubscribeMessageRequestSeekMessage(mq_brokerwire.SubscribeMessageRequestSeekMessageInput{
			Offset: m.Seek.Offset, OffsetType: uint32(m.Seek.OffsetType),
		})
	}
	return mq_brokerwire.NewSubscribeMessageRequest(in)
}

func SubscribeMessageRespFromWire(b []byte) (*mq_pb.SubscribeMessageResponse, error) {
	v, err := mq_brokerwire.WrapSubscribeMessageResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &mq_pb.SubscribeMessageResponse{}
	switch v.WhichMessage() {
	case mq_brokerwire.SubscribeMessageResponseMessageCtrl:
		c := v.Ctrl()
		resp.Message = &mq_pb.SubscribeMessageResponse_Ctrl{Ctrl: &mq_pb.SubscribeMessageResponse_SubscribeCtrlMessage{
			Error: c.Error(), IsEndOfStream: c.IsEndOfStream(), IsEndOfTopic: c.IsEndOfTopic(),
		}}
	case mq_brokerwire.SubscribeMessageResponseMessageData:
		resp.Message = &mq_pb.SubscribeMessageResponse_Data{Data: dataMessageFromView(v.Data())}
	}
	return resp, nil
}

// PublishFollowMe (bidi)

func PublishFollowMeReqToWire(r *mq_pb.PublishFollowMeRequest) []byte {
	in := mq_brokerwire.PublishFollowMeRequestInput{MessageWhich: mq_brokerwire.PublishFollowMeRequestMessageNone}
	switch m := r.Message.(type) {
	case *mq_pb.PublishFollowMeRequest_Init:
		in.MessageWhich = mq_brokerwire.PublishFollowMeRequestMessageInit
		in.MessageValue = mq_brokerwire.NewPublishFollowMeRequestInitMessage(mq_brokerwire.PublishFollowMeRequestInitMessageInput{
			Topic: topicToWire(m.Init.Topic), Partition: partitionToWire(m.Init.Partition),
		})
	case *mq_pb.PublishFollowMeRequest_Data:
		in.MessageWhich = mq_brokerwire.PublishFollowMeRequestMessageData
		in.MessageValue = dataMessageToWire(m.Data)
	case *mq_pb.PublishFollowMeRequest_Flush:
		in.MessageWhich = mq_brokerwire.PublishFollowMeRequestMessageFlush
		in.MessageValue = mq_brokerwire.NewPublishFollowMeRequestFlushMessage(mq_brokerwire.PublishFollowMeRequestFlushMessageInput{TsNs: m.Flush.TsNs})
	case *mq_pb.PublishFollowMeRequest_Close:
		in.MessageWhich = mq_brokerwire.PublishFollowMeRequestMessageClose
		in.MessageValue = mq_brokerwire.NewPublishFollowMeRequestCloseMessage(mq_brokerwire.PublishFollowMeRequestCloseMessageInput{})
	}
	return mq_brokerwire.NewPublishFollowMeRequest(in)
}

func PublishFollowMeRespFromWire(b []byte) (*mq_pb.PublishFollowMeResponse, error) {
	v, err := mq_brokerwire.WrapPublishFollowMeResponse(b)
	if err != nil {
		return nil, err
	}
	return &mq_pb.PublishFollowMeResponse{AckTsNs: v.AckTsNs()}, nil
}

// SubscribeFollowMe (client-stream)

func SubscribeFollowMeReqToWire(r *mq_pb.SubscribeFollowMeRequest) []byte {
	in := mq_brokerwire.SubscribeFollowMeRequestInput{MessageWhich: mq_brokerwire.SubscribeFollowMeRequestMessageNone}
	switch m := r.Message.(type) {
	case *mq_pb.SubscribeFollowMeRequest_Init:
		in.MessageWhich = mq_brokerwire.SubscribeFollowMeRequestMessageInit
		in.MessageValue = mq_brokerwire.NewSubscribeFollowMeRequestInitMessage(mq_brokerwire.SubscribeFollowMeRequestInitMessageInput{
			Topic: topicToWire(m.Init.Topic), Partition: partitionToWire(m.Init.Partition), ConsumerGroup: m.Init.ConsumerGroup,
		})
	case *mq_pb.SubscribeFollowMeRequest_Ack:
		in.MessageWhich = mq_brokerwire.SubscribeFollowMeRequestMessageAck
		in.MessageValue = mq_brokerwire.NewSubscribeFollowMeRequestAckMessage(mq_brokerwire.SubscribeFollowMeRequestAckMessageInput{TsNs: m.Ack.TsNs})
	case *mq_pb.SubscribeFollowMeRequest_Close:
		in.MessageWhich = mq_brokerwire.SubscribeFollowMeRequestMessageClose
		in.MessageValue = mq_brokerwire.NewSubscribeFollowMeRequestCloseMessage(mq_brokerwire.SubscribeFollowMeRequestCloseMessageInput{})
	}
	return mq_brokerwire.NewSubscribeFollowMeRequest(in)
}

func SubscribeFollowMeRespFromWire(b []byte) (*mq_pb.SubscribeFollowMeResponse, error) {
	v, err := mq_brokerwire.WrapSubscribeFollowMeResponse(b)
	if err != nil {
		return nil, err
	}
	return &mq_pb.SubscribeFollowMeResponse{AckTsNs: v.AckTsNs()}, nil
}

// GetUnflushedMessages (server-stream)

func GetUnflushedMessagesReqToWire(r *mq_pb.GetUnflushedMessagesRequest) []byte {
	return mq_brokerwire.NewGetUnflushedMessagesRequest(mq_brokerwire.GetUnflushedMessagesRequestInput{
		Topic: topicToWire(r.Topic), Partition: partitionToWire(r.Partition), StartBufferOffset: r.StartBufferOffset,
	})
}

func GetUnflushedMessagesRespFromWire(b []byte) (*mq_pb.GetUnflushedMessagesResponse, error) {
	v, err := mq_brokerwire.WrapGetUnflushedMessagesResponse(b)
	if err != nil {
		return nil, err
	}
	return &mq_pb.GetUnflushedMessagesResponse{
		Message:     logEntryFromWire(v.Message()),
		Error:       v.Error(),
		EndOfStream: v.EndOfStream(),
	}, nil
}
