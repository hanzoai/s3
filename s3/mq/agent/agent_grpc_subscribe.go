package agent

import (
	"context"

	"github.com/zap-proto/go/transport"
	"google.golang.org/protobuf/proto"

	"github.com/hanzoai/s3/s3/glog"
	"github.com/hanzoai/s3/s3/mq/agent/agentconv"
	"github.com/hanzoai/s3/s3/mq/client/sub_client"
	"github.com/hanzoai/s3/s3/mq/topic"
	"github.com/hanzoai/s3/s3/pb"
	"github.com/hanzoai/s3/s3/pb/mq_pb"
	"github.com/hanzoai/s3/s3/pb/schema_pb"
	"github.com/hanzoai/s3/s3/util"
	mq_agentwire "github.com/hanzoai/s3/s3/wire/mq_agent"
	mq_schemawire "github.com/hanzoai/s3/s3/wire/mq_schema"
)

// SubscribeRecord serves the HanzoMessagingAgent SubscribeRecord bidirectional
// stream over the ZAP transport: the first frame carries the init request, each
// subsequent inbound frame an ack; each delivered record is sent back as a
// SubscribeRecordResponse frame.
func (a *MessageQueueAgent) SubscribeRecord(init []byte, s transport.Stream) error {
	// the first message is the subscribe request; it carries the init request
	initMessage, err := mq_agentwire.WrapSubscribeRecordRequest(init)
	if err != nil {
		return err
	}

	subscriber, err := a.handleInitSubscribeRecordRequest(context.Background(), initMessage.Init())
	if err != nil {
		return err
	}

	var lastErr error
	executors := util.NewLimitedConcurrentExecutor(int(subscriber.SubscriberConfig.SlidingWindowSize))
	subscriber.SetOnDataMessageFn(func(m *mq_pb.SubscribeMessageResponse_Data) {
		executors.Execute(func() {
			record := &schema_pb.RecordValue{}
			err := proto.Unmarshal(m.Data.Value, record)
			if err != nil {
				glog.V(0).Infof("unmarshal record value: %v", err)
				if lastErr == nil {
					lastErr = err
				}
				return
			}
			resp := mq_agentwire.NewSubscribeRecordResponse(mq_agentwire.SubscribeRecordResponseInput{
				Key:   m.Data.Key,
				Value: agentconv.RecordValueToWire(record),
				TsNs:  m.Data.TsNs,
			})
			if sendErr := s.Send(resp); sendErr != nil {
				glog.V(0).Infof("send record: %v", sendErr)
				if lastErr == nil {
					lastErr = sendErr
				}
			}
		})
	})

	go func() {
		subErr := subscriber.Subscribe()
		if subErr != nil {
			glog.V(0).Infof("subscriber %s subscribe: %v", subscriber.SubscriberConfig.String(), subErr)
			if lastErr == nil {
				lastErr = subErr
			}
		}
	}()

	for {
		frame, err := s.Recv()
		if err != nil {
			glog.V(0).Infof("subscriber %s receive: %v", subscriber.SubscriberConfig.String(), err)
			return err
		}
		m, err := mq_agentwire.WrapSubscribeRecordRequest(frame)
		if err != nil {
			glog.V(0).Infof("subscriber %s parse ack: %v", subscriber.SubscriberConfig.String(), err)
			return err
		}
		subscriber.PartitionOffsetChan <- sub_client.KeyedTimestamp{
			Key:  m.AckKey(),
			TsNs: m.AckSequence(), // Note: AckSequence should be renamed to AckTsNs in agent protocol
		}
	}
}

func (a *MessageQueueAgent) handleInitSubscribeRecordRequest(ctx context.Context, initBytes []byte) (*sub_client.TopicSubscriber, error) {
	req, err := mq_agentwire.WrapInitSubscribeRecordRequest(initBytes)
	if err != nil {
		return nil, err
	}

	subscriberConfig := &sub_client.SubscriberConfiguration{
		ConsumerGroup:           req.ConsumerGroup(),
		ConsumerGroupInstanceId: req.ConsumerGroupInstanceId(),
		GrpcDialOption:          pb.DialOption{},
		MaxPartitionCount:       req.MaxSubscribedPartitions(),
		SlidingWindowSize:       req.SlidingWindowSize(),
	}

	partitionOffsets := make([]*schema_pb.PartitionOffset, 0, req.PartitionOffsetsLen())
	for i := 0; i < req.PartitionOffsetsLen(); i++ {
		po, err := mq_schemawire.WrapPartitionOffset(req.PartitionOffsetAt(i))
		if err != nil {
			return nil, err
		}
		partitionOffsets = append(partitionOffsets, agentconv.PartitionOffsetsFromWire(po, i))
	}

	contentConfig := &sub_client.ContentConfiguration{
		Topic:            topic.FromPbTopic(agentconv.TopicFromWire(req.Topic())),
		Filter:           req.Filter(),
		PartitionOffsets: partitionOffsets,
		OffsetType:       schema_pb.OffsetType(req.OffsetType()),
		OffsetTsNs:       req.OffsetTsNs(),
	}

	topicSubscriber := sub_client.NewTopicSubscriber(
		ctx,
		a.brokersList(),
		subscriberConfig,
		contentConfig,
		make(chan sub_client.KeyedTimestamp, 1024),
	)

	return topicSubscriber, nil
}
