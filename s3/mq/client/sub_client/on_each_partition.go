package sub_client

import (
	"errors"
	"fmt"
	"io"

	"github.com/hanzoai/s3/s3/glog"
	"github.com/hanzoai/s3/s3/mq/broker/brokerpb"
	"github.com/hanzoai/s3/s3/pb/mq_pb"
	"github.com/hanzoai/s3/s3/pb/schema_pb"
	mq_brokerwire "github.com/hanzoai/s3/s3/wire/mq_broker"
	"github.com/zap-proto/go/transport"
)

type KeyedTimestamp struct {
	Key  []byte
	TsNs int64 // Timestamp in nanoseconds for acknowledgment
}

func (sub *TopicSubscriber) onEachPartition(assigned *mq_pb.BrokerPartitionAssignment, stopCh chan struct{}, onDataMessageFn OnDataMessageFn) error {
	// connect to the partition broker
	conn, err := transport.Dial("tcp", assigned.LeaderBroker)
	if err != nil {
		return fmt.Errorf("dial partition broker %s: %w", assigned.LeaderBroker, err)
	}
	defer conn.Close()

	slidingWindowSize := sub.SubscriberConfig.SlidingWindowSize
	if slidingWindowSize <= 0 {
		slidingWindowSize = 1
	}

	po := findPartitionOffset(sub.ContentConfig.PartitionOffsets, assigned.Partition)
	if po == nil {
		po = &schema_pb.PartitionOffset{
			Partition: assigned.Partition,
			StartTsNs: sub.ContentConfig.OffsetTsNs,
		}
	}

	initBuf := brokerpb.SubscribeMessageInitToWire(&mq_pb.SubscribeMessageRequest_InitMessage{
		ConsumerGroup:     sub.SubscriberConfig.ConsumerGroup,
		ConsumerId:        sub.SubscriberConfig.ConsumerGroupInstanceId,
		Topic:             sub.ContentConfig.Topic.ToPbTopic(),
		PartitionOffset:   po,
		OffsetType:        sub.ContentConfig.OffsetType,
		Filter:            sub.ContentConfig.Filter,
		FollowerBroker:    assigned.FollowerBroker,
		SlidingWindowSize: slidingWindowSize,
	})
	subscribeClient, err := conn.OpenStream(mq_brokerwire.HanzoMessagingSubscribeMessageOrdinal, initBuf)
	if err != nil {
		return fmt.Errorf("create subscribe client: %w", err)
	}

	glog.V(0).Infof("subscriber %s connected to partition %+v at %v", sub.ContentConfig.Topic, assigned.Partition, assigned.LeaderBroker)

	if sub.OnCompletionFunc != nil {
		defer sub.OnCompletionFunc()
	}

	go func() {
		for {
			select {
			case <-sub.ctx.Done():
				subscribeClient.CloseSend()
				return
			case <-stopCh:
				subscribeClient.CloseSend()
				return
			case ack, ok := <-sub.PartitionOffsetChan:
				if !ok {
					subscribeClient.CloseSend()
					return
				}
				subscribeClient.Send(brokerpb.SubscribeMessageAckToWire(ack.Key, ack.TsNs))
			}
		}
	}()

	for {
		// glog.V(0).Infof("subscriber %s/%s waiting for message", sub.ContentConfig.Topic, sub.SubscriberConfig.ConsumerGroup)
		body, err := subscribeClient.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return fmt.Errorf("subscribe recv: %w", err)
		}
		rv, err := mq_brokerwire.WrapSubscribeMessageResponse(body)
		if err != nil {
			return fmt.Errorf("subscribe decode: %w", err)
		}
		resp := brokerpb.SubscribeMessageResponseFromWire(rv)
		if resp.Message == nil {
			glog.V(0).Infof("subscriber %s/%s received nil message", sub.ContentConfig.Topic, sub.SubscriberConfig.ConsumerGroup)
			continue
		}

		select {
		case <-sub.ctx.Done():
			return nil
		case <-stopCh:
			return nil
		default:
		}

		switch m := resp.Message.(type) {
		case *mq_pb.SubscribeMessageResponse_Data:
			if m.Data.Ctrl != nil {
				glog.V(2).Infof("subscriber %s received control from producer:%s isClose:%v", sub.SubscriberConfig.ConsumerGroup, m.Data.Ctrl.PublisherName, m.Data.Ctrl.IsClose)
				continue
			}
			if len(m.Data.Key) == 0 {
				// fmt.Printf("empty key %+v, type %v\n", m, reflect.TypeOf(m))
				continue
			}
			onDataMessageFn(m)
		case *mq_pb.SubscribeMessageResponse_Ctrl:
			// glog.V(0).Infof("subscriber %s/%s/%s received control %+v", sub.ContentConfig.Namespace, sub.ContentConfig.Topic, sub.SubscriberConfig.ConsumerGroup, m.Ctrl)
			if m.Ctrl.IsEndOfStream || m.Ctrl.IsEndOfTopic {
				return io.EOF
			}
		}
	}
}

func findPartitionOffset(partitionOffsets []*schema_pb.PartitionOffset, partition *schema_pb.Partition) *schema_pb.PartitionOffset {
	for _, po := range partitionOffsets {
		if po.Partition == partition {
			return po
		}
	}
	return nil
}
