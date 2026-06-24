// Package brokerpb transcodes the mq_pb broker-domain messages that ride inside
// the HanzoMessaging RPCs (DataMessage, BrokerPartitionAssignment, the
// SubscribeMessage and SubscriberToSubCoordinator stream frames) between their
// protobuf form — still the in-memory model of the broker/topic/publisher/
// subscriber layers — and their ZAP wire form (mq_brokerwire). These are the
// one-and-only conversions for that boundary, shared by every client and server
// call site so they never re-implement the framing. schema_pb children (Topic,
// Partition, PartitionOffset, RecordType) are delegated to agentconv, the
// canonical schema_pb<->mq_schemawire bridge.
package brokerpb

import (
	"github.com/hanzoai/s3/s3/mq/agent/agentconv"
	"github.com/hanzoai/s3/s3/pb/mq_pb"
	mq_brokerwire "github.com/hanzoai/s3/s3/wire/mq_broker"
	mq_schemawire "github.com/hanzoai/s3/s3/wire/mq_schema"
)

// --- DataMessage ---

// DataMessageToWire encodes an in-memory mq_pb.DataMessage as its ZAP buffer.
func DataMessageToWire(d *mq_pb.DataMessage) []byte {
	if d == nil {
		return nil
	}
	in := mq_brokerwire.DataMessageInput{Key: d.Key, Value: d.Value, TsNs: d.TsNs}
	if d.Ctrl != nil {
		in.Ctrl = mq_brokerwire.NewControlMessage(mq_brokerwire.ControlMessageInput{
			IsClose:       d.Ctrl.IsClose,
			PublisherName: d.Ctrl.PublisherName,
		})
	}
	return mq_brokerwire.NewDataMessage(in)
}

// DataMessageFromWire decodes a DataMessage view into an mq_pb.DataMessage.
func DataMessageFromWire(d mq_brokerwire.DataMessage) *mq_pb.DataMessage {
	out := &mq_pb.DataMessage{Key: d.Key(), Value: d.Value(), TsNs: d.TsNs()}
	if ctrl, ok := d.Ctrl(); ok {
		out.Ctrl = &mq_pb.ControlMessage{IsClose: ctrl.IsClose(), PublisherName: ctrl.PublisherName()}
	}
	return out
}

// --- BrokerPartitionAssignment ---

// BrokerPartitionAssignmentFromWire decodes a wire view into the mq_pb model.
func BrokerPartitionAssignmentFromWire(bpa mq_brokerwire.BrokerPartitionAssignment) *mq_pb.BrokerPartitionAssignment {
	out := &mq_pb.BrokerPartitionAssignment{
		LeaderBroker:   bpa.LeaderBroker(),
		FollowerBroker: bpa.FollowerBroker(),
	}
	if pbuf := bpa.Partition(); len(pbuf) > 0 {
		if pv, err := mq_schemawire.WrapPartition(pbuf); err == nil {
			out.Partition = agentconv.PartitionFromWire(pv)
		}
	}
	return out
}

// --- SubscribeMessageRequest (Init / Ack) ---

// SubscribeMessageInitToWire builds the SubscribeMessageRequest Init frame.
func SubscribeMessageInitToWire(in *mq_pb.SubscribeMessageRequest_InitMessage) []byte {
	return mq_brokerwire.NewSubscribeMessageRequest(mq_brokerwire.SubscribeMessageRequestInput{
		MessageWhich: mq_brokerwire.SubscribeMessageRequestMessageInit,
		MessageValue: mq_brokerwire.NewSubscribeMessageRequestInitMessage(mq_brokerwire.SubscribeMessageRequestInitMessageInput{
			ConsumerGroup:     in.ConsumerGroup,
			ConsumerId:        in.ConsumerId,
			ClientId:          in.ClientId,
			Topic:             agentconv.TopicToWire(in.Topic),
			PartitionOffset:   agentconv.PartitionOffsetToWire(in.PartitionOffset),
			OffsetType:        uint32(in.OffsetType),
			Filter:            in.Filter,
			FollowerBroker:    in.FollowerBroker,
			SlidingWindowSize: in.SlidingWindowSize,
		}),
	})
}

// SubscribeMessageAckToWire builds the SubscribeMessageRequest Ack frame.
func SubscribeMessageAckToWire(key []byte, tsNs int64) []byte {
	return mq_brokerwire.NewSubscribeMessageRequest(mq_brokerwire.SubscribeMessageRequestInput{
		MessageWhich: mq_brokerwire.SubscribeMessageRequestMessageAck,
		MessageValue: mq_brokerwire.NewSubscribeMessageRequestAckMessage(mq_brokerwire.SubscribeMessageRequestAckMessageInput{
			Key:  key,
			TsNs: tsNs,
		}),
	})
}

// SubscribeMessageResponseFromWire decodes one SubscribeMessageResponse frame
// into the mq_pb model, preserving the Data / Ctrl oneof variant.
func SubscribeMessageResponseFromWire(r mq_brokerwire.SubscribeMessageResponse) *mq_pb.SubscribeMessageResponse {
	switch r.WhichMessage() {
	case mq_brokerwire.SubscribeMessageResponseMessageData:
		return &mq_pb.SubscribeMessageResponse{
			Message: &mq_pb.SubscribeMessageResponse_Data{Data: DataMessageFromWire(r.Data())},
		}
	case mq_brokerwire.SubscribeMessageResponseMessageCtrl:
		ctrl := r.Ctrl()
		return &mq_pb.SubscribeMessageResponse{
			Message: &mq_pb.SubscribeMessageResponse_Ctrl{
				Ctrl: &mq_pb.SubscribeMessageResponse_SubscribeCtrlMessage{
					IsEndOfStream: ctrl.IsEndOfStream(),
					IsEndOfTopic:  ctrl.IsEndOfTopic(),
				},
			},
		}
	}
	return &mq_pb.SubscribeMessageResponse{}
}

// --- SubscriberToSubCoordinatorRequest (Init / AckAssignment / AckUnAssignment) ---

// SubscriberToSubCoordinatorInitToWire builds the coordinator Init frame.
func SubscriberToSubCoordinatorInitToWire(in *mq_pb.SubscriberToSubCoordinatorRequest_InitMessage) []byte {
	return mq_brokerwire.NewSubscriberToSubCoordinatorRequest(mq_brokerwire.SubscriberToSubCoordinatorRequestInput{
		MessageWhich: mq_brokerwire.SubscriberToSubCoordinatorRequestMessageInit,
		MessageValue: mq_brokerwire.NewSubscriberToSubCoordinatorRequestInitMessage(mq_brokerwire.SubscriberToSubCoordinatorRequestInitMessageInput{
			ConsumerGroup:           in.ConsumerGroup,
			ConsumerGroupInstanceId: in.ConsumerGroupInstanceId,
			Topic:                   agentconv.TopicToWire(in.Topic),
			MaxPartitionCount:       in.MaxPartitionCount,
			RebalanceSeconds:        in.RebalanceSeconds,
		}),
	})
}

// SubscriberToSubCoordinatorRequestToWire encodes any coordinator request
// (Init / AckAssignment / AckUnAssignment) into its ZAP frame.
func SubscriberToSubCoordinatorRequestToWire(req *mq_pb.SubscriberToSubCoordinatorRequest) []byte {
	switch m := req.Message.(type) {
	case *mq_pb.SubscriberToSubCoordinatorRequest_Init:
		return SubscriberToSubCoordinatorInitToWire(m.Init)
	case *mq_pb.SubscriberToSubCoordinatorRequest_AckAssignment:
		return mq_brokerwire.NewSubscriberToSubCoordinatorRequest(mq_brokerwire.SubscriberToSubCoordinatorRequestInput{
			MessageWhich: mq_brokerwire.SubscriberToSubCoordinatorRequestMessageAckAssignment,
			MessageValue: mq_brokerwire.NewSubscriberToSubCoordinatorRequestAckAssignmentMessage(mq_brokerwire.SubscriberToSubCoordinatorRequestAckAssignmentMessageInput{
				Partition: agentconv.PartitionToWire(m.AckAssignment.Partition),
			}),
		})
	case *mq_pb.SubscriberToSubCoordinatorRequest_AckUnAssignment:
		return mq_brokerwire.NewSubscriberToSubCoordinatorRequest(mq_brokerwire.SubscriberToSubCoordinatorRequestInput{
			MessageWhich: mq_brokerwire.SubscriberToSubCoordinatorRequestMessageAckUnAssignment,
			MessageValue: mq_brokerwire.NewSubscriberToSubCoordinatorRequestAckUnAssignmentMessage(mq_brokerwire.SubscriberToSubCoordinatorRequestAckUnAssignmentMessageInput{
				Partition: agentconv.PartitionToWire(m.AckUnAssignment.Partition),
			}),
		})
	}
	return mq_brokerwire.NewSubscriberToSubCoordinatorRequest(mq_brokerwire.SubscriberToSubCoordinatorRequestInput{})
}

// SubscriberToSubCoordinatorResponseFromWire decodes one coordinator response
// frame into the mq_pb model, preserving the Assignment / UnAssignment variant.
func SubscriberToSubCoordinatorResponseFromWire(r mq_brokerwire.SubscriberToSubCoordinatorResponse) *mq_pb.SubscriberToSubCoordinatorResponse {
	switch r.WhichMessage() {
	case mq_brokerwire.SubscriberToSubCoordinatorResponseMessageAssignment:
		a := r.Assignment()
		out := &mq_pb.SubscriberToSubCoordinatorResponse_Assignment{}
		if bpa, ok := a.PartitionAssignment(); ok {
			out.PartitionAssignment = BrokerPartitionAssignmentFromWire(bpa)
		}
		return &mq_pb.SubscriberToSubCoordinatorResponse{
			Message: &mq_pb.SubscriberToSubCoordinatorResponse_Assignment_{Assignment: out},
		}
	case mq_brokerwire.SubscriberToSubCoordinatorResponseMessageUnAssignment:
		u := r.UnAssignment()
		out := &mq_pb.SubscriberToSubCoordinatorResponse_UnAssignment{}
		if pbuf := u.Partition(); len(pbuf) > 0 {
			if pv, err := mq_schemawire.WrapPartition(pbuf); err == nil {
				out.Partition = agentconv.PartitionFromWire(pv)
			}
		}
		return &mq_pb.SubscriberToSubCoordinatorResponse{
			Message: &mq_pb.SubscriberToSubCoordinatorResponse_UnAssignment_{UnAssignment: out},
		}
	}
	return &mq_pb.SubscriberToSubCoordinatorResponse{}
}
