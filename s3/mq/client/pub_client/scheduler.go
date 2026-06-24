package pub_client

import (
	"fmt"
	"io"
	"log"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hanzoai/s3/s3/glog"
	"github.com/hanzoai/s3/s3/mq/agent/agentconv"
	"github.com/hanzoai/s3/s3/mq/broker/brokerpb"
	"github.com/hanzoai/s3/s3/mq/topic"
	"github.com/hanzoai/s3/s3/pb"
	"github.com/hanzoai/s3/s3/pb/mq_pb"
	"github.com/hanzoai/s3/s3/pb/schema_pb"
	"github.com/hanzoai/s3/s3/util/buffered_queue"
	mq_brokerwire "github.com/hanzoai/s3/s3/wire/mq_broker"
	mq_schemawire "github.com/hanzoai/s3/s3/wire/mq_schema"
	"github.com/zap-proto/go/transport"
)

type EachPartitionError struct {
	*mq_pb.BrokerPartitionAssignment
	Err        error
	generation int
}

type EachPartitionPublishJob struct {
	*mq_pb.BrokerPartitionAssignment
	stopChan   chan bool
	wg         sync.WaitGroup
	generation int
	inputQueue *buffered_queue.BufferedQueue[*mq_pb.DataMessage]
}

func (p *TopicPublisher) startSchedulerThread(wg *sync.WaitGroup) error {

	if err := p.doConfigureTopic(); err != nil {
		wg.Done()
		return fmt.Errorf("configure topic %s: %v", p.config.Topic, err)
	}

	log.Printf("start scheduler thread for topic %s", p.config.Topic)

	generation := 0
	var errChan chan EachPartitionError
	for {
		glog.V(0).Infof("lookup partitions gen %d topic %s", generation+1, p.config.Topic)
		if assignments, err := p.doLookupTopicPartitions(); err == nil {
			generation++
			glog.V(0).Infof("start generation %d with %d assignments", generation, len(assignments))
			if errChan == nil {
				errChan = make(chan EachPartitionError, len(assignments))
			}
			p.onEachAssignments(generation, assignments, errChan)
		} else {
			glog.Errorf("lookup topic %s: %v", p.config.Topic, err)
			time.Sleep(5 * time.Second)
			continue
		}

		if generation == 1 {
			wg.Done()
		}

		// wait for any error to happen. If so, consume all remaining errors, and retry
		for {
			select {
			case eachErr := <-errChan:
				glog.Errorf("gen %d publish to topic %s partition %v: %v", eachErr.generation, p.config.Topic, eachErr.Partition, eachErr.Err)
				if eachErr.generation < generation {
					continue
				}
				break
			}
		}
	}
}

func (p *TopicPublisher) onEachAssignments(generation int, assignments []*mq_pb.BrokerPartitionAssignment, errChan chan EachPartitionError) {
	// TODO assuming this is not re-configured so the partitions are fixed.
	sort.Slice(assignments, func(i, j int) bool {
		return assignments[i].Partition.RangeStart < assignments[j].Partition.RangeStart
	})
	var jobs []*EachPartitionPublishJob
	hasExistingJob := len(p.jobs) == len(assignments)
	for i, assignment := range assignments {
		if assignment.LeaderBroker == "" {
			continue
		}
		if hasExistingJob {
			var existingJob *EachPartitionPublishJob
			existingJob = p.jobs[i]
			if pb.ServerAddress(existingJob.BrokerPartitionAssignment.LeaderBroker).Equals(pb.ServerAddress(assignment.LeaderBroker)) {
				existingJob.generation = generation
				jobs = append(jobs, existingJob)
				continue
			} else {
				if existingJob.LeaderBroker != "" {
					close(existingJob.stopChan)
					existingJob.LeaderBroker = ""
					existingJob.wg.Wait()
				}
			}
		}

		// start a go routine to publish to this partition
		job := &EachPartitionPublishJob{
			BrokerPartitionAssignment: assignment,
			stopChan:                  make(chan bool, 1),
			generation:                generation,
			inputQueue:                buffered_queue.NewBufferedQueue[*mq_pb.DataMessage](1024),
		}
		job.wg.Add(1)
		go func(job *EachPartitionPublishJob) {
			defer job.wg.Done()
			if err := p.doPublishToPartition(job); err != nil {
				log.Printf("publish to %s partition %v: %v", p.config.Topic, job.Partition, err)
				errChan <- EachPartitionError{assignment, err, generation}
			}
		}(job)
		jobs = append(jobs, job)
		// TODO assuming this is not re-configured so the partitions are fixed.
		// better just re-use the existing job
		p.partition2Buffer.Insert(assignment.Partition.RangeStart, assignment.Partition.RangeStop, job.inputQueue)
	}
	p.jobs = jobs
}

func (p *TopicPublisher) doPublishToPartition(job *EachPartitionPublishJob) error {

	log.Printf("connecting to %v for topic partition %+v", job.LeaderBroker, job.Partition)

	initBuf := mq_brokerwire.NewPublishMessageRequest(mq_brokerwire.PublishMessageRequestInput{
		MessageWhich: mq_brokerwire.PublishMessageRequestMessageInit,
		MessageValue: mq_brokerwire.NewPublishMessageRequestInitMessage(mq_brokerwire.PublishMessageRequestInitMessageInput{
			Topic:          wireTopic(p.config.Topic),
			Partition:      wirePartition(job.Partition),
			AckInterval:    128,
			FollowerBroker: job.FollowerBroker,
			PublisherName:  p.config.PublisherName,
		}),
	})

	conn, err := transport.Dial("tcp", job.LeaderBroker)
	if err != nil {
		return fmt.Errorf("dial broker %s: %v", job.LeaderBroker, err)
	}
	defer conn.Close()
	stream, err := conn.OpenStream(mq_brokerwire.HanzoMessagingPublishMessageOrdinal, initBuf)
	if err != nil {
		return fmt.Errorf("create publish client: %w", err)
	}
	publishClient := &PublishClient{
		stream: stream,
		Broker: job.LeaderBroker,
	}
	// process the hello message
	helloBody, err := stream.Recv()
	if err != nil {
		return fmt.Errorf("recv init response: %w", err)
	}
	hello, err := mq_brokerwire.WrapPublishMessageResponse(helloBody)
	if err != nil {
		return fmt.Errorf("decode init response: %w", err)
	}
	if hello.Error() != "" {
		return fmt.Errorf("init response error: %v", hello.Error())
	}

	var publishedTsNs int64
	hasMoreData := int32(1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			ackBody, err := publishClient.Recv()
			if err != nil {
				if err == io.EOF {
					log.Printf("publish to %s EOF", publishClient.Broker)
					return
				}
				publishClient.Err = err
				log.Printf("publish1 to %s error: %v\n", publishClient.Broker, err)
				return
			}
			ackResp, err := mq_brokerwire.WrapPublishMessageResponse(ackBody)
			if err != nil {
				publishClient.Err = err
				log.Printf("publish1 to %s decode error: %v\n", publishClient.Broker, err)
				return
			}
			if ackResp.Error() != "" {
				publishClient.Err = fmt.Errorf("ack error: %v", ackResp.Error())
				log.Printf("publish2 to %s error: %v\n", publishClient.Broker, ackResp.Error())
				return
			}
			ackTsNs := ackResp.AckTsNs()
			if ackTsNs > 0 {
				log.Printf("ack %d published %d hasMoreData:%d", ackTsNs, atomic.LoadInt64(&publishedTsNs), atomic.LoadInt32(&hasMoreData))
			}
			if atomic.LoadInt64(&publishedTsNs) <= ackTsNs && atomic.LoadInt32(&hasMoreData) == 0 {
				return
			}
		}
	}()

	publishCounter := 0
	for data, hasData := job.inputQueue.Dequeue(); hasData; data, hasData = job.inputQueue.Dequeue() {
		if data.Ctrl != nil && data.Ctrl.IsClose {
			// need to set this before sending to brokers, to avoid timing issue
			atomic.StoreInt32(&hasMoreData, 0)
		}
		if err := publishClient.Send(mq_brokerwire.NewPublishMessageRequest(mq_brokerwire.PublishMessageRequestInput{
			MessageWhich: mq_brokerwire.PublishMessageRequestMessageData,
			MessageValue: brokerpb.DataMessageToWire(data),
		})); err != nil {
			return fmt.Errorf("send publish data: %w", err)
		}
		publishCounter++
		atomic.StoreInt64(&publishedTsNs, data.TsNs)
	}
	if publishCounter > 0 {
		wg.Wait()
	} else {
		// CloseSend would cancel the context on the server side
		if err := publishClient.CloseSend(); err != nil {
			return fmt.Errorf("close send: %w", err)
		}
	}

	log.Printf("published %d messages to %v for topic partition %+v", publishCounter, job.LeaderBroker, job.Partition)

	return nil
}

func (p *TopicPublisher) doConfigureTopic() (err error) {
	if len(p.config.Brokers) == 0 {
		return fmt.Errorf("topic configuring found no bootstrap brokers")
	}
	var lastErr error
	for _, brokerAddress := range p.config.Brokers {
		err = func() error {
			conn, err := transport.Dial("tcp", brokerAddress)
			if err != nil {
				return err
			}
			defer conn.Close()
			client := mq_brokerwire.NewClient(conn)
			_, err = client.ConfigureTopic(mq_brokerwire.ConfigureTopicRequestInput{
				Topic:             wireTopic(p.config.Topic),
				PartitionCount:    p.config.PartitionCount,
				MessageRecordType: agentconv.RecordTypeToWire(p.config.RecordType), // Flat schema
			})
			return err
		}()
		if err == nil {
			lastErr = nil
			return nil
		} else {
			lastErr = fmt.Errorf("%s: %v", brokerAddress, err)
		}
	}

	if lastErr != nil {
		return fmt.Errorf("doConfigureTopic %s: %v", p.config.Topic, err)
	}
	return nil
}

func (p *TopicPublisher) doLookupTopicPartitions() (assignments []*mq_pb.BrokerPartitionAssignment, err error) {
	if len(p.config.Brokers) == 0 {
		return nil, fmt.Errorf("lookup found no bootstrap brokers")
	}
	var lastErr error
	for _, brokerAddress := range p.config.Brokers {
		err := func() error {
			conn, err := transport.Dial("tcp", brokerAddress)
			if err != nil {
				return err
			}
			defer conn.Close()
			client := mq_brokerwire.NewClient(conn)
			lookupResp, err := client.LookupTopicBrokers(mq_brokerwire.LookupTopicBrokersRequestInput{
				Topic: wireTopic(p.config.Topic),
			})
			if err != nil {
				return err
			}
			glog.V(0).Infof("lookup topic %s: %d assignments", p.config.Topic, lookupResp.BrokerPartitionAssignmentsLen())

			n := lookupResp.BrokerPartitionAssignmentsLen()
			if n == 0 {
				return fmt.Errorf("no broker partition assignments")
			}

			out := make([]*mq_pb.BrokerPartitionAssignment, 0, n)
			for i := 0; i < n; i++ {
				bpa, ok := lookupResp.BrokerPartitionAssignmentAt(i)
				if !ok {
					continue
				}
				out = append(out, brokerpb.BrokerPartitionAssignmentFromWire(bpa))
			}
			assignments = out
			return nil
		}()
		if err == nil {
			return assignments, nil
		} else {
			lastErr = err
		}
	}

	return nil, fmt.Errorf("lookup topic %s: %v", p.config.Topic, lastErr)

}

// --- schema_pb<->wire conversions for the ZAP transport boundary ---

func wireTopic(t topic.Topic) []byte {
	return mq_schemawire.NewTopic(mq_schemawire.TopicInput{Namespace: t.Namespace, Name: t.Name})
}

func wirePartition(p *schema_pb.Partition) []byte {
	return agentconv.PartitionToWire(p)
}
