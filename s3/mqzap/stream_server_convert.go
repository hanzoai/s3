// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// stream_server_convert.go is the server-side mirror of stream.go: per-frame
// decoders that wrap an inbound mq_brokerwire streaming-request buffer into the
// *mq_pb.<Rpc>Request the engine expects, and encoders that turn each
// *mq_pb.<Rpc>Response the engine emits into a zero-copy mq_brokerwire frame.
// stream.go owns the client direction (*mq_pb -> wire, wire -> *mq_pb); this file
// owns the server direction. Both reuse the SAME leaf converters in convert.go.
// Naming: <Rpc>ReqFromWire([]byte) (*mq_pb.<Rpc>Request, error) and
// <Rpc>RespToWire(*mq_pb.<Rpc>Response) []byte.

package mqzap

import (
	mq_pb "github.com/hanzoai/s3/s3/pb/mq_pb"
	schema_pb "github.com/hanzoai/s3/s3/pb/schema_pb"
	mq_brokerwire "github.com/hanzoai/s3/s3/wire/mq_broker"
)

// --- PublisherToPubBalancer (bidi) ---

func PublisherToPubBalancerReqFromWire(b []byte) (*mq_pb.PublisherToPubBalancerRequest, error) {
	v, err := mq_brokerwire.WrapPublisherToPubBalancerRequest(b)
	if err != nil {
		return nil, err
	}
	req := &mq_pb.PublisherToPubBalancerRequest{}
	switch v.WhichMessage() {
	case mq_brokerwire.PublisherToPubBalancerRequestMessageInit:
		req.Message = &mq_pb.PublisherToPubBalancerRequest_Init{
			Init: &mq_pb.PublisherToPubBalancerRequest_InitMessage{Broker: v.Init().Broker()},
		}
	case mq_brokerwire.PublisherToPubBalancerRequestMessageStats:
		req.Message = &mq_pb.PublisherToPubBalancerRequest_Stats{Stats: brokerStatsFromView(v.Stats())}
	}
	return req, nil
}

// PublisherToPubBalancerResponse is empty; the server replies with the empty
// frame (built inline in stream_server.go's Send).

// --- SubscriberToSubCoordinator (bidi) ---

func SubscriberToSubCoordinatorReqFromWire(b []byte) (*mq_pb.SubscriberToSubCoordinatorRequest, error) {
	v, err := mq_brokerwire.WrapSubscriberToSubCoordinatorRequest(b)
	if err != nil {
		return nil, err
	}
	req := &mq_pb.SubscriberToSubCoordinatorRequest{}
	switch v.WhichMessage() {
	case mq_brokerwire.SubscriberToSubCoordinatorRequestMessageInit:
		i := v.Init()
		req.Message = &mq_pb.SubscriberToSubCoordinatorRequest_Init{
			Init: &mq_pb.SubscriberToSubCoordinatorRequest_InitMessage{
				ConsumerGroup:           i.ConsumerGroup(),
				ConsumerGroupInstanceId: i.ConsumerGroupInstanceId(),
				Topic:                   topicFromWire(i.Topic()),
				MaxPartitionCount:       i.MaxPartitionCount(),
				RebalanceSeconds:        i.RebalanceSeconds(),
			},
		}
	case mq_brokerwire.SubscriberToSubCoordinatorRequestMessageAckAssignment:
		req.Message = &mq_pb.SubscriberToSubCoordinatorRequest_AckAssignment{
			AckAssignment: &mq_pb.SubscriberToSubCoordinatorRequest_AckAssignmentMessage{
				Partition: partitionFromWire(v.AckAssignment().Partition()),
			},
		}
	case mq_brokerwire.SubscriberToSubCoordinatorRequestMessageAckUnAssignment:
		req.Message = &mq_pb.SubscriberToSubCoordinatorRequest_AckUnAssignment{
			AckUnAssignment: &mq_pb.SubscriberToSubCoordinatorRequest_AckUnAssignmentMessage{
				Partition: partitionFromWire(v.AckUnAssignment().Partition()),
			},
		}
	}
	return req, nil
}

func SubscriberToSubCoordinatorRespToWire(r *mq_pb.SubscriberToSubCoordinatorResponse) []byte {
	in := mq_brokerwire.SubscriberToSubCoordinatorResponseInput{MessageWhich: mq_brokerwire.SubscriberToSubCoordinatorResponseMessageNone}
	switch m := r.Message.(type) {
	case *mq_pb.SubscriberToSubCoordinatorResponse_Assignment_:
		in.MessageWhich = mq_brokerwire.SubscriberToSubCoordinatorResponseMessageAssignment
		var pa []byte
		if m.Assignment != nil {
			pa = assignmentToWire(m.Assignment.PartitionAssignment)
		}
		in.MessageValue = mq_brokerwire.NewSubscriberToSubCoordinatorResponseAssignment(mq_brokerwire.SubscriberToSubCoordinatorResponseAssignmentInput{
			PartitionAssignment: pa,
		})
	case *mq_pb.SubscriberToSubCoordinatorResponse_UnAssignment_:
		in.MessageWhich = mq_brokerwire.SubscriberToSubCoordinatorResponseMessageUnAssignment
		var part []byte
		if m.UnAssignment != nil {
			part = partitionToWire(m.UnAssignment.Partition)
		}
		in.MessageValue = mq_brokerwire.NewSubscriberToSubCoordinatorResponseUnAssignment(mq_brokerwire.SubscriberToSubCoordinatorResponseUnAssignmentInput{
			Partition: part,
		})
	}
	return mq_brokerwire.NewSubscriberToSubCoordinatorResponse(in)
}

// --- PublishMessage (bidi) ---

func PublishMessageReqFromWire(b []byte) (*mq_pb.PublishMessageRequest, error) {
	v, err := mq_brokerwire.WrapPublishMessageRequest(b)
	if err != nil {
		return nil, err
	}
	req := &mq_pb.PublishMessageRequest{}
	switch v.WhichMessage() {
	case mq_brokerwire.PublishMessageRequestMessageInit:
		i := v.Init()
		req.Message = &mq_pb.PublishMessageRequest_Init{
			Init: &mq_pb.PublishMessageRequest_InitMessage{
				Topic:          topicFromWire(i.Topic()),
				Partition:      partitionFromWire(i.Partition()),
				AckInterval:    i.AckInterval(),
				FollowerBroker: i.FollowerBroker(),
				PublisherName:  i.PublisherName(),
			},
		}
	case mq_brokerwire.PublishMessageRequestMessageData:
		req.Message = &mq_pb.PublishMessageRequest_Data{Data: dataMessageFromView(v.Data())}
	}
	return req, nil
}

func PublishMessageRespToWire(r *mq_pb.PublishMessageResponse) []byte {
	return mq_brokerwire.NewPublishMessageResponse(mq_brokerwire.PublishMessageResponseInput{
		AckTsNs:        r.AckTsNs,
		Error:          r.Error,
		ShouldClose:    r.ShouldClose,
		ErrorCode:      r.ErrorCode,
		AssignedOffset: r.AssignedOffset,
	})
}

// --- SubscribeMessage (bidi) ---

func SubscribeMessageReqFromWire(b []byte) (*mq_pb.SubscribeMessageRequest, error) {
	v, err := mq_brokerwire.WrapSubscribeMessageRequest(b)
	if err != nil {
		return nil, err
	}
	req := &mq_pb.SubscribeMessageRequest{}
	switch v.WhichMessage() {
	case mq_brokerwire.SubscribeMessageRequestMessageInit:
		i := v.Init()
		req.Message = &mq_pb.SubscribeMessageRequest_Init{
			Init: &mq_pb.SubscribeMessageRequest_InitMessage{
				ConsumerGroup:     i.ConsumerGroup(),
				ConsumerId:        i.ConsumerId(),
				ClientId:          i.ClientId(),
				Topic:             topicFromWire(i.Topic()),
				PartitionOffset:   partitionOffsetFromWire(i.PartitionOffset()),
				OffsetType:        schema_pb.OffsetType(i.OffsetType()),
				Filter:            i.Filter(),
				FollowerBroker:    i.FollowerBroker(),
				SlidingWindowSize: i.SlidingWindowSize(),
			},
		}
	case mq_brokerwire.SubscribeMessageRequestMessageAck:
		a := v.Ack()
		req.Message = &mq_pb.SubscribeMessageRequest_Ack{
			Ack: &mq_pb.SubscribeMessageRequest_AckMessage{TsNs: a.TsNs(), Key: a.Key()},
		}
	case mq_brokerwire.SubscribeMessageRequestMessageSeek:
		s := v.Seek()
		req.Message = &mq_pb.SubscribeMessageRequest_Seek{
			Seek: &mq_pb.SubscribeMessageRequest_SeekMessage{
				Offset: s.Offset(), OffsetType: schema_pb.OffsetType(s.OffsetType()),
			},
		}
	}
	return req, nil
}

func SubscribeMessageRespToWire(r *mq_pb.SubscribeMessageResponse) []byte {
	in := mq_brokerwire.SubscribeMessageResponseInput{MessageWhich: mq_brokerwire.SubscribeMessageResponseMessageNone}
	switch m := r.Message.(type) {
	case *mq_pb.SubscribeMessageResponse_Ctrl:
		in.MessageWhich = mq_brokerwire.SubscribeMessageResponseMessageCtrl
		in.MessageValue = mq_brokerwire.NewSubscribeMessageResponseSubscribeCtrlMessage(mq_brokerwire.SubscribeMessageResponseSubscribeCtrlMessageInput{
			Error:         m.Ctrl.Error,
			IsEndOfStream: m.Ctrl.IsEndOfStream,
			IsEndOfTopic:  m.Ctrl.IsEndOfTopic,
		})
	case *mq_pb.SubscribeMessageResponse_Data:
		in.MessageWhich = mq_brokerwire.SubscribeMessageResponseMessageData
		in.MessageValue = dataMessageToWire(m.Data)
	}
	return mq_brokerwire.NewSubscribeMessageResponse(in)
}

// --- PublishFollowMe (bidi) ---

func PublishFollowMeReqFromWire(b []byte) (*mq_pb.PublishFollowMeRequest, error) {
	v, err := mq_brokerwire.WrapPublishFollowMeRequest(b)
	if err != nil {
		return nil, err
	}
	req := &mq_pb.PublishFollowMeRequest{}
	switch v.WhichMessage() {
	case mq_brokerwire.PublishFollowMeRequestMessageInit:
		i := v.Init()
		req.Message = &mq_pb.PublishFollowMeRequest_Init{
			Init: &mq_pb.PublishFollowMeRequest_InitMessage{
				Topic: topicFromWire(i.Topic()), Partition: partitionFromWire(i.Partition()),
			},
		}
	case mq_brokerwire.PublishFollowMeRequestMessageData:
		req.Message = &mq_pb.PublishFollowMeRequest_Data{Data: dataMessageFromView(v.Data())}
	case mq_brokerwire.PublishFollowMeRequestMessageFlush:
		req.Message = &mq_pb.PublishFollowMeRequest_Flush{
			Flush: &mq_pb.PublishFollowMeRequest_FlushMessage{TsNs: v.Flush().TsNs()},
		}
	case mq_brokerwire.PublishFollowMeRequestMessageClose:
		req.Message = &mq_pb.PublishFollowMeRequest_Close{
			Close: &mq_pb.PublishFollowMeRequest_CloseMessage{},
		}
	}
	return req, nil
}

func PublishFollowMeRespToWire(r *mq_pb.PublishFollowMeResponse) []byte {
	return mq_brokerwire.NewPublishFollowMeResponse(mq_brokerwire.PublishFollowMeResponseInput{AckTsNs: r.AckTsNs})
}

// --- SubscribeFollowMe (client-stream) ---

func SubscribeFollowMeReqFromWire(b []byte) (*mq_pb.SubscribeFollowMeRequest, error) {
	v, err := mq_brokerwire.WrapSubscribeFollowMeRequest(b)
	if err != nil {
		return nil, err
	}
	req := &mq_pb.SubscribeFollowMeRequest{}
	switch v.WhichMessage() {
	case mq_brokerwire.SubscribeFollowMeRequestMessageInit:
		i := v.Init()
		req.Message = &mq_pb.SubscribeFollowMeRequest_Init{
			Init: &mq_pb.SubscribeFollowMeRequest_InitMessage{
				Topic: topicFromWire(i.Topic()), Partition: partitionFromWire(i.Partition()), ConsumerGroup: i.ConsumerGroup(),
			},
		}
	case mq_brokerwire.SubscribeFollowMeRequestMessageAck:
		req.Message = &mq_pb.SubscribeFollowMeRequest_Ack{
			Ack: &mq_pb.SubscribeFollowMeRequest_AckMessage{TsNs: v.Ack().TsNs()},
		}
	case mq_brokerwire.SubscribeFollowMeRequestMessageClose:
		req.Message = &mq_pb.SubscribeFollowMeRequest_Close{
			Close: &mq_pb.SubscribeFollowMeRequest_CloseMessage{},
		}
	}
	return req, nil
}

func SubscribeFollowMeRespToWire(r *mq_pb.SubscribeFollowMeResponse) []byte {
	return mq_brokerwire.NewSubscribeFollowMeResponse(mq_brokerwire.SubscribeFollowMeResponseInput{AckTsNs: r.AckTsNs})
}

// --- GetUnflushedMessages (server-stream) ---

func GetUnflushedMessagesRespToWire(r *mq_pb.GetUnflushedMessagesResponse) []byte {
	return mq_brokerwire.NewGetUnflushedMessagesResponse(mq_brokerwire.GetUnflushedMessagesResponseInput{
		Message:     logEntryToWire(r.Message),
		Error:       r.Error,
		EndOfStream: r.EndOfStream,
	})
}
