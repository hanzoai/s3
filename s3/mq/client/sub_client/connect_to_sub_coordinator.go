package sub_client

import (
	"io"
	"time"

	"github.com/hanzoai/s3/s3/glog"
	"github.com/hanzoai/s3/s3/mq/broker/brokerpb"
	"github.com/hanzoai/s3/s3/pb/mq_pb"
	mq_brokerwire "github.com/hanzoai/s3/s3/wire/mq_broker"
	"github.com/zap-proto/go/transport"
)

func (sub *TopicSubscriber) doKeepConnectedToSubCoordinator() {
	waitTime := 1 * time.Second
	for {
		for _, broker := range sub.bootstrapBrokers {

			select {
			case <-sub.ctx.Done():
				return
			default:
			}

			// lookup topic brokers
			var brokerLeader string
			err := func() error {
				conn, err := transport.Dial("tcp", broker)
				if err != nil {
					return err
				}
				defer conn.Close()
				client := mq_brokerwire.NewClient(conn)
				resp, err := client.FindBrokerLeader(mq_brokerwire.FindBrokerLeaderRequestInput{})
				if err != nil {
					return err
				}
				brokerLeader = resp.Broker()
				return nil
			}()
			if err != nil {
				glog.V(0).Infof("broker coordinator on %s: %v", broker, err)
				continue
			}
			glog.V(0).Infof("found broker coordinator: %v", brokerLeader)

			// connect to the balancer
			func() error {
				conn, err := transport.Dial("tcp", brokerLeader)
				if err != nil {
					glog.V(0).Infof("subscriber %s dial leader: %v", sub.ContentConfig.Topic, err)
					return err
				}
				defer conn.Close()

				initBuf := brokerpb.SubscriberToSubCoordinatorInitToWire(&mq_pb.SubscriberToSubCoordinatorRequest_InitMessage{
					ConsumerGroup:           sub.SubscriberConfig.ConsumerGroup,
					ConsumerGroupInstanceId: sub.SubscriberConfig.ConsumerGroupInstanceId,
					Topic:                   sub.ContentConfig.Topic.ToPbTopic(),
					MaxPartitionCount:       sub.SubscriberConfig.MaxPartitionCount,
				})
				stream, err := conn.OpenStream(mq_brokerwire.HanzoMessagingSubscriberToSubCoordinatorOrdinal, initBuf)
				if err != nil {
					glog.V(0).Infof("subscriber %s: %v", sub.ContentConfig.Topic, err)
					return err
				}
				waitTime = 1 * time.Second

				// Maybe later: subscribe to multiple topics instead of just one

				go func() {
					for reply := range sub.brokerPartitionAssignmentAckChan {

						select {
						case <-sub.ctx.Done():
							return
						default:
						}

						glog.V(0).Infof("subscriber instance %s ack %+v", sub.SubscriberConfig.ConsumerGroupInstanceId, reply)
						if err := stream.Send(brokerpb.SubscriberToSubCoordinatorRequestToWire(reply)); err != nil {
							glog.V(0).Infof("subscriber %s reply: %v", sub.ContentConfig.Topic, err)
							return
						}
					}
				}()

				// keep receiving messages from the sub coordinator
				for {
					body, err := stream.Recv()
					if err != nil {
						if err == io.EOF {
							return nil
						}
						glog.V(0).Infof("subscriber %s receive: %v", sub.ContentConfig.Topic, err)
						return err
					}

					select {
					case <-sub.ctx.Done():
						return nil
					default:
					}

					rv, err := mq_brokerwire.WrapSubscriberToSubCoordinatorResponse(body)
					if err != nil {
						glog.V(0).Infof("subscriber %s decode assignment: %v", sub.ContentConfig.Topic, err)
						return err
					}
					resp := brokerpb.SubscriberToSubCoordinatorResponseFromWire(rv)
					sub.brokerPartitionAssignmentChan <- resp
					glog.V(0).Infof("Received assignment: %+v", resp)
				}
			}()
		}
		glog.V(0).Infof("subscriber %s/%s waiting for more assignments", sub.ContentConfig.Topic, sub.SubscriberConfig.ConsumerGroup)
		if waitTime < 10*time.Second {
			waitTime += 1 * time.Second
		}
		time.Sleep(waitTime)
	}
}
