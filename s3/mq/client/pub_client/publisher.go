package pub_client

import (
	"log"
	"sync"

	"github.com/rdleal/intervalst/interval"
	"github.com/zap-proto/go/transport"

	"github.com/hanzoai/s3/s3/mq/pub_balancer"
	"github.com/hanzoai/s3/s3/mq/topic"
	"github.com/hanzoai/s3/s3/pb"
	"github.com/hanzoai/s3/s3/pb/mq_pb"
	"github.com/hanzoai/s3/s3/pb/schema_pb"
	"github.com/hanzoai/s3/s3/util/buffered_queue"
)

type PublisherConfiguration struct {
	Topic          topic.Topic
	PartitionCount int32
	Brokers        []string
	PublisherName  string // for debugging
	RecordType     *schema_pb.RecordType
}

type PublishClient struct {
	stream transport.Stream
	Broker string
	Err    error
}

// Send transmits one PublishMessageRequest frame over the ZAP stream.
func (c *PublishClient) Send(req []byte) error { return c.stream.Send(req) }

// Recv reads the next PublishMessageResponse frame from the ZAP stream.
func (c *PublishClient) Recv() ([]byte, error) { return c.stream.Recv() }

// CloseSend half-closes the send side of the ZAP stream.
func (c *PublishClient) CloseSend() error { return c.stream.CloseSend() }

type TopicPublisher struct {
	partition2Buffer *interval.SearchTree[*buffered_queue.BufferedQueue[*mq_pb.DataMessage], int32]
	grpcDialOption   pb.DialOption
	sync.Mutex       // protects grpc
	config           *PublisherConfiguration
	jobs             []*EachPartitionPublishJob
}

func NewTopicPublisher(config *PublisherConfiguration) (tp *TopicPublisher, err error) {
	tp = &TopicPublisher{
		partition2Buffer: interval.NewSearchTree[*buffered_queue.BufferedQueue[*mq_pb.DataMessage]](func(a, b int32) int {
			return int(a - b)
		}),
		grpcDialOption: pb.DialOption{},
		config:         config,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		if err = tp.startSchedulerThread(&wg); err != nil {
			log.Println(err)
			return
		}
	}()

	wg.Wait()

	return
}

func (p *TopicPublisher) Shutdown() error {

	if inputBuffers, found := p.partition2Buffer.AllIntersections(0, pub_balancer.MaxPartitionCount); found {
		for _, inputBuffer := range inputBuffers {
			inputBuffer.CloseInput()
		}
	}

	for _, job := range p.jobs {
		job.wg.Wait()
	}

	return nil
}
