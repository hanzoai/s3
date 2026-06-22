package sub_coordinator

import (
	"github.com/hanzoai/s3/s3/mq/topic"
	"github.com/hanzoai/s3/s3/pb/mq_pb"
)

type ConsumerGroupInstanceId string

type ConsumerGroupInstance struct {
	InstanceId         ConsumerGroupInstanceId
	AssignedPartitions []topic.Partition
	ResponseChan       chan *mq_pb.SubscriberToSubCoordinatorResponse
	MaxPartitionCount  int32
}

func NewConsumerGroupInstance(instanceId string, maxPartitionCount int32) *ConsumerGroupInstance {
	return &ConsumerGroupInstance{
		InstanceId:        ConsumerGroupInstanceId(instanceId),
		ResponseChan:      make(chan *mq_pb.SubscriberToSubCoordinatorResponse, 1),
		MaxPartitionCount: maxPartitionCount,
	}
}
