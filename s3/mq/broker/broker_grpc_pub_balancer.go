package broker

import (
	"fmt"

	"github.com/hanzoai/s3/s3/mq/pub_balancer"
	"github.com/hanzoai/s3/s3/pb/mq_pb"
)

// PublisherToPubBalancer receives connections from brokers and collects stats
func (b *MessageQueueBroker) PublisherToPubBalancer(stream mq_pb.HanzoMessaging_PublisherToPubBalancerServer) error {
	if !b.isLockOwner() {
		return fmt.Errorf("Unavailable: not current broker balancer")
	}
	req, err := stream.Recv()
	if err != nil {
		return fmt.Errorf("receive init message: %w", err)
	}

	// process init message
	initMessage := req.GetInit()
	var brokerStats *pub_balancer.BrokerStats
	if initMessage != nil {
		brokerStats = b.PubBalancer.AddBroker(initMessage.Broker)
	} else {
		return fmt.Errorf("InvalidArgument: balancer init message is empty")
	}
	defer func() {
		b.PubBalancer.RemoveBroker(initMessage.Broker, brokerStats)
	}()

	// process stats message
	for {
		req, err := stream.Recv()
		if err != nil {
			return fmt.Errorf("receive stats message from %s: %v", initMessage.Broker, err)
		}
		if !b.isLockOwner() {
			return fmt.Errorf("Unavailable: not current broker balancer")
		}
		if receivedStats := req.GetStats(); receivedStats != nil {
			b.PubBalancer.OnBrokerStatsUpdated(initMessage.Broker, brokerStats, receivedStats)
			// glog.V(4).Infof("received from %v: %+v", initMessage.Broker, receivedStats)
		}
	}
}
