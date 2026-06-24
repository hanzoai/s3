package pub_balancer

import (
	"fmt"

	mq_brokerwire "github.com/hanzoai/s3/s3/wire/mq_broker"
	mq_schemawire "github.com/hanzoai/s3/s3/wire/mq_schema"
	"github.com/zap-proto/go/transport"
	"google.golang.org/grpc"
)

// PubBalancer <= PublisherToPubBalancer() <= Broker <=> Publish()
// ExecuteBalanceActionMove from PubBalancer => AssignTopicPartitions() => Broker => Publish()

func (balancer *PubBalancer) ExecuteBalanceActionMove(move *BalanceActionMove, _ grpc.DialOption) error {
	if _, found := balancer.Brokers.Get(move.SourceBroker); !found {
		return fmt.Errorf("source broker %s not found", move.SourceBroker)
	}
	if _, found := balancer.Brokers.Get(move.TargetBroker); !found {
		return fmt.Errorf("target broker %s not found", move.TargetBroker)
	}

	topicBuf := mq_schemawire.NewTopic(mq_schemawire.TopicInput{
		Namespace: move.TopicPartition.Namespace,
		Name:      move.TopicPartition.Name,
	})
	partitionBuf := mq_schemawire.NewPartition(mq_schemawire.PartitionInput{
		RingSize:   move.TopicPartition.RingSize,
		RangeStart: move.TopicPartition.RangeStart,
		RangeStop:  move.TopicPartition.RangeStop,
		UnixTimeNs: move.TopicPartition.UnixTimeNs,
	})

	assign := func(broker string, draining bool) error {
		conn, err := transport.Dial("tcp", broker)
		if err != nil {
			return err
		}
		defer conn.Close()
		client := mq_brokerwire.NewClient(conn)
		_, err = client.AssignTopicPartitions(mq_brokerwire.AssignTopicPartitionsRequestInput{
			Topic: topicBuf,
			BrokerPartitionAssignments: [][]byte{
				mq_brokerwire.NewBrokerPartitionAssignment(mq_brokerwire.BrokerPartitionAssignmentInput{
					Partition: partitionBuf,
				}),
			},
			IsLeader:   true,
			IsDraining: draining,
		})
		return err
	}

	if err := assign(move.TargetBroker, false); err != nil {
		return fmt.Errorf("assign topic partition %v to %s: %v", move.TopicPartition, move.TargetBroker, err)
	}
	if err := assign(move.SourceBroker, true); err != nil {
		return fmt.Errorf("assign topic partition %v to %s: %v", move.TopicPartition, move.SourceBroker, err)
	}

	return nil
}
