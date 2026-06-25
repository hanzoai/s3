package topic

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zap-proto/go/transport"

	"github.com/hanzoai/s3/s3/glog"
	"github.com/hanzoai/s3/s3/mq/agent/agentconv"
	"github.com/hanzoai/s3/s3/mq/broker/brokerpb"
	"github.com/hanzoai/s3/s3/pb"
	"github.com/hanzoai/s3/s3/pb/filer_pb"
	"github.com/hanzoai/s3/s3/pb/mq_pb"
	"github.com/hanzoai/s3/s3/util/log_buffer"
	mq_brokerwire "github.com/hanzoai/s3/s3/wire/mq_broker"
)

type LocalPartition struct {
	ListenersWaits int64
	AckTsNs        int64

	// notifying clients
	ListenersLock sync.Mutex
	ListenersCond *sync.Cond

	Partition
	LogBuffer   *log_buffer.LogBuffer
	Publishers  *LocalPartitionPublishers
	Subscribers *LocalPartitionSubscribers

	publishFolloweMeStream transport.Stream
	followerConn           transport.Conn
	Follower               string

	// Track last activity for idle cleanup
	lastActivityTime atomic.Int64 // Unix nano timestamp
}

var TIME_FORMAT = "2006-01-02-15-04-05"
var PartitionGenerationFormat = "v2006-01-02-15-04-05"

func NewLocalPartition(partition Partition, logFlushInterval int, logFlushFn log_buffer.LogFlushFuncType, readFromDiskFn log_buffer.LogReadFromDiskFuncType) *LocalPartition {
	lp := &LocalPartition{
		Partition:   partition,
		Publishers:  NewLocalPartitionPublishers(),
		Subscribers: NewLocalPartitionSubscribers(),
	}
	lp.ListenersCond = sync.NewCond(&lp.ListenersLock)
	lp.lastActivityTime.Store(time.Now().UnixNano()) // Initialize with current time

	// Ensure a minimum flush interval to prevent busy-loop when set to 0
	// A flush interval of 0 would cause time.Sleep(0) creating a CPU-consuming busy loop
	flushInterval := time.Duration(logFlushInterval) * time.Second
	if flushInterval == 0 {
		flushInterval = 1 * time.Second // Minimum 1 second to avoid busy-loop, allow near-immediate flushing
	}

	lp.LogBuffer = log_buffer.NewLogBuffer(fmt.Sprintf("%d/%04d-%04d", partition.UnixTimeNs, partition.RangeStart, partition.RangeStop),
		flushInterval, logFlushFn, readFromDiskFn, func() {
			if atomic.LoadInt64(&lp.ListenersWaits) > 0 {
				lp.ListenersCond.Broadcast()
			}
		})
	return lp
}

func (p *LocalPartition) Publish(message *mq_pb.DataMessage) error {
	if err := p.LogBuffer.AddToBuffer(message); err != nil {
		return fmt.Errorf("failed to add message to log buffer: %w", err)
	}
	p.UpdateActivity() // Track publish activity for idle cleanup

	// maybe send to the follower
	if p.publishFolloweMeStream != nil {
		// println("recv", string(message.Key), message.TsNs)
		if followErr := p.publishFolloweMeStream.Send(mq_brokerwire.NewPublishFollowMeRequest(mq_brokerwire.PublishFollowMeRequestInput{
			MessageWhich: mq_brokerwire.PublishFollowMeRequestMessageData,
			MessageValue: brokerpb.DataMessageToWire(message),
		})); followErr != nil {
			return fmt.Errorf("send to follower %s: %v", p.Follower, followErr)
		}
	} else {
		atomic.StoreInt64(&p.AckTsNs, message.TsNs)
	}

	return nil
}

func (p *LocalPartition) Subscribe(clientName string, startPosition log_buffer.MessagePosition,
	onNoMessageFn func() bool, eachMessageFn log_buffer.EachLogEntryFuncType) error {
	var processedPosition log_buffer.MessagePosition
	var readPersistedLogErr error
	var readInMemoryLogErr error
	var isDone bool

	p.UpdateActivity() // Track subscribe activity for idle cleanup

	// CRITICAL FIX: Use offset-based functions if startPosition is offset-based
	// This allows reading historical data by offset, not just by timestamp
	if startPosition.IsOffsetBased {
		// Wrap eachMessageFn to match the signature expected by LoopProcessLogDataWithOffset
		// Also update activity when messages are processed
		eachMessageWithOffsetFn := func(logEntry *filer_pb.LogEntry, offset int64) (bool, error) {
			p.UpdateActivity() // Track message read activity
			return eachMessageFn(logEntry)
		}

		// Wrap eachMessageFn for disk reads to also update activity
		eachMessageWithActivityFn := func(logEntry *filer_pb.LogEntry) (bool, error) {
			p.UpdateActivity() // Track disk read activity for idle cleanup
			return eachMessageFn(logEntry)
		}

		// Always attempt initial disk read for historical data
		// This is fast if no data on disk, and ensures we don't miss old data
		// The memory read loop below handles new data with instant notifications
		glog.V(2).Infof("%s reading historical data from disk starting at offset %d", clientName, startPosition.Offset)
		processedPosition, isDone, readPersistedLogErr = p.LogBuffer.ReadFromDiskFn(startPosition, 0, eachMessageWithActivityFn)
		if readPersistedLogErr != nil {
			glog.V(2).Infof("%s read %v persisted log: %v", clientName, p.Partition, readPersistedLogErr)
			return readPersistedLogErr
		}
		if isDone {
			return nil
		}

		// Update position after reading from disk
		if processedPosition.Time.UnixNano() != 0 || processedPosition.IsOffsetBased {
			startPosition = processedPosition
		}

		// Step 2: Enter the main loop - read from in-memory buffer, occasionally checking disk
		for {
			// Read from in-memory buffer (this is the hot path - handles streaming data)
			glog.V(4).Infof("SUBSCRIBE: Reading from in-memory buffer for %s at offset %d", clientName, startPosition.Offset)
			processedPosition, isDone, readInMemoryLogErr = p.LogBuffer.LoopProcessLogDataWithOffset(clientName, startPosition, 0, onNoMessageFn, eachMessageWithOffsetFn)

			if isDone {
				return nil
			}

			// Update position
			// CRITICAL FIX: For offset-based reads, Time is zero, so check Offset instead
			if processedPosition.Time.UnixNano() != 0 || processedPosition.IsOffsetBased {
				startPosition = processedPosition
			}

			// If we get ResumeFromDiskError, it means data was flushed to disk
			// Read from disk ONCE to catch up, then continue with in-memory buffer
			if readInMemoryLogErr == log_buffer.ResumeFromDiskError {
				glog.V(4).Infof("SUBSCRIBE: ResumeFromDiskError - reading flushed data from disk for %s at offset %d", clientName, startPosition.Offset)
				processedPosition, isDone, readPersistedLogErr = p.LogBuffer.ReadFromDiskFn(startPosition, 0, eachMessageWithActivityFn)
				if readPersistedLogErr != nil {
					glog.V(2).Infof("%s read %v persisted log after flush: %v", clientName, p.Partition, readPersistedLogErr)
					return readPersistedLogErr
				}
				if isDone {
					return nil
				}

				// Update position and continue the loop (back to in-memory buffer)
				// CRITICAL FIX: For offset-based reads, Time is zero, so check Offset instead
				if processedPosition.Time.UnixNano() != 0 || processedPosition.IsOffsetBased {
					startPosition = processedPosition
				}
				// Loop continues - back to reading from in-memory buffer
				continue
			}

			// Any other error is a real error
			if readInMemoryLogErr != nil {
				glog.V(2).Infof("%s read %v in memory log: %v", clientName, p.Partition, readInMemoryLogErr)
				return readInMemoryLogErr
			}

			// If we get here with no error and not done, something is wrong
			glog.V(1).Infof("SUBSCRIBE: Unexpected state for %s - no error but not done, continuing", clientName)
		}
	}

	// Original timestamp-based subscription logic
	// Wrap eachMessageFn for disk reads to also update activity
	eachMessageWithActivityFn := func(logEntry *filer_pb.LogEntry) (bool, error) {
		p.UpdateActivity() // Track disk read activity for idle cleanup
		return eachMessageFn(logEntry)
	}

	for {
		processedPosition, isDone, readPersistedLogErr = p.LogBuffer.ReadFromDiskFn(startPosition, 0, eachMessageWithActivityFn)
		if readPersistedLogErr != nil {
			glog.V(0).Infof("%s read %v persisted log: %v", clientName, p.Partition, readPersistedLogErr)
			return readPersistedLogErr
		}
		if isDone {
			return nil
		}

		// CRITICAL FIX: For offset-based reads, Time is zero, so check Offset instead
		if processedPosition.Time.UnixNano() != 0 || processedPosition.IsOffsetBased {
			startPosition = processedPosition
		}
		processedPosition, isDone, readInMemoryLogErr = p.LogBuffer.LoopProcessLogData(clientName, startPosition, 0, onNoMessageFn, eachMessageFn)
		if isDone {
			return nil
		}
		// CRITICAL FIX: For offset-based reads, Time is zero, so check Offset instead
		if processedPosition.Time.UnixNano() != 0 || processedPosition.IsOffsetBased {
			startPosition = processedPosition
		}

		if readInMemoryLogErr == log_buffer.ResumeFromDiskError {
			continue
		}
		if readInMemoryLogErr != nil {
			glog.V(0).Infof("%s read %v in memory log: %v", clientName, p.Partition, readInMemoryLogErr)
			return readInMemoryLogErr
		}
	}
}

func (p *LocalPartition) GetEarliestMessageTimeInMemory() time.Time {
	return p.LogBuffer.GetEarliestTime()
}

func (p *LocalPartition) HasData() bool {
	return !p.LogBuffer.GetEarliestTime().IsZero()
}

func (p *LocalPartition) GetEarliestInMemoryMessagePosition() log_buffer.MessagePosition {
	return p.LogBuffer.GetEarliestPosition()
}

func (p *LocalPartition) closePublishers() {
	p.Publishers.SignalShutdown()
}
func (p *LocalPartition) closeSubscribers() {
	p.Subscribers.SignalShutdown()
}

func (p *LocalPartition) WaitUntilNoPublishers() {
	for {
		if p.Publishers.Size() == 0 {
			return
		}
		time.Sleep(113 * time.Millisecond)
	}
}

func (p *LocalPartition) MaybeConnectToFollowers(initMessage *mq_pb.PublishMessageRequest_InitMessage, _ pb.DialOption) (err error) {
	if p.publishFolloweMeStream != nil {
		return nil
	}
	if initMessage.FollowerBroker == "" {
		return nil
	}

	p.Follower = initMessage.FollowerBroker
	p.followerConn, err = transport.Dial("tcp", p.Follower)
	if err != nil {
		return fmt.Errorf("fail to dial %s: %v", p.Follower, err)
	}
	initBuf := mq_brokerwire.NewPublishFollowMeRequest(mq_brokerwire.PublishFollowMeRequestInput{
		MessageWhich: mq_brokerwire.PublishFollowMeRequestMessageInit,
		MessageValue: mq_brokerwire.NewPublishFollowMeRequestInitMessage(mq_brokerwire.PublishFollowMeRequestInitMessageInput{
			Topic:     agentconv.TopicToWire(initMessage.Topic),
			Partition: agentconv.PartitionToWire(initMessage.Partition),
		}),
	})
	p.publishFolloweMeStream, err = p.followerConn.OpenStream(mq_brokerwire.HanzoMessagingPublishFollowMeOrdinal, initBuf)
	if err != nil {
		p.followerConn.Close()
		p.followerConn = nil
		return fmt.Errorf("fail to create publish client: %w", err)
	}

	// start receiving ack from follower
	go func() {
		for {
			ackBody, err := p.publishFolloweMeStream.Recv()
			if err != nil {
				if err == io.EOF {
					glog.V(0).Infof("local partition %v follower %v stopped", p.Partition, p.Follower)
					return
				}
				glog.Errorf("Receiving local partition %v  follower %s ack: %v", p.Partition, p.Follower, err)
				return
			}
			ack, err := mq_brokerwire.WrapPublishFollowMeResponse(ackBody)
			if err != nil {
				glog.Errorf("decode local partition %v follower %s ack: %v", p.Partition, p.Follower, err)
				return
			}
			atomic.StoreInt64(&p.AckTsNs, ack.AckTsNs())
		}
	}()
	return nil
}

func (p *LocalPartition) MaybeShutdownLocalPartition() (hasShutdown bool) {

	if p.Publishers.Size() == 0 && p.Subscribers.Size() == 0 {
		p.LogBuffer.ShutdownLogBuffer()
		for !p.LogBuffer.IsAllFlushed() {
			time.Sleep(113 * time.Millisecond)
		}
		if p.publishFolloweMeStream != nil {
			// send close to the follower
			if followErr := p.publishFolloweMeStream.Send(mq_brokerwire.NewPublishFollowMeRequest(mq_brokerwire.PublishFollowMeRequestInput{
				MessageWhich: mq_brokerwire.PublishFollowMeRequestMessageClose,
				MessageValue: mq_brokerwire.NewPublishFollowMeRequestCloseMessage(mq_brokerwire.PublishFollowMeRequestCloseMessageInput{}),
			})); followErr != nil {
				glog.Errorf("Error closing follower stream: %v", followErr)
			}
			glog.V(4).Infof("closing connection to follower")
			p.followerConn.Close()
			p.publishFolloweMeStream = nil
			p.followerConn = nil
			p.Follower = ""
		}

		hasShutdown = true
	}

	glog.V(0).Infof("local partition %v Publisher:%d Subscriber:%d follower:%s shutdown %v", p.Partition, p.Publishers.Size(), p.Subscribers.Size(), p.Follower, hasShutdown)
	return
}

// MaybeShutdownLocalPartitionForTopic is a topic-aware version that considers system topic retention
func (p *LocalPartition) MaybeShutdownLocalPartitionForTopic(topicName string) (hasShutdown bool) {
	// For system topics like _schemas, be more conservative about shutdown
	if isSystemTopic(topicName) {
		glog.V(0).Infof("System topic %s - skipping aggressive shutdown for partition %v (Publishers:%d Subscribers:%d)",
			topicName, p.Partition, p.Publishers.Size(), p.Subscribers.Size())
		return false
	}

	// For regular topics, use the standard shutdown logic
	return p.MaybeShutdownLocalPartition()
}

// isSystemTopic checks if a topic should have special retention behavior
func isSystemTopic(topicName string) bool {
	systemTopics := []string{
		"_schemas",            // Schema Registry topic
		"__consumer_offsets",  // Kafka consumer offsets topic
		"__transaction_state", // Kafka transaction state topic
	}

	for _, systemTopic := range systemTopics {
		if topicName == systemTopic {
			return true
		}
	}

	// Also check for topics with system prefixes
	return strings.HasPrefix(topicName, "_") || strings.HasPrefix(topicName, "__")
}

func (p *LocalPartition) Shutdown() {
	p.closePublishers()
	p.closeSubscribers()
	p.LogBuffer.ShutdownLogBuffer()
	glog.V(0).Infof("local partition %v shutting down", p.Partition)
}

func (p *LocalPartition) NotifyLogFlushed(flushTsNs int64) {
	if p.publishFolloweMeStream != nil {
		if followErr := p.publishFolloweMeStream.Send(mq_brokerwire.NewPublishFollowMeRequest(mq_brokerwire.PublishFollowMeRequestInput{
			MessageWhich: mq_brokerwire.PublishFollowMeRequestMessageFlush,
			MessageValue: mq_brokerwire.NewPublishFollowMeRequestFlushMessage(mq_brokerwire.PublishFollowMeRequestFlushMessageInput{
				TsNs: flushTsNs,
			}),
		})); followErr != nil {
			glog.Errorf("send follower %s flush message: %v", p.Follower, followErr)
		}
		// println("notifying", p.Follower, "flushed at", flushTsNs)
	}
}

// UpdateActivity updates the last activity timestamp for this partition
// Should be called whenever a publisher publishes or a subscriber reads
func (p *LocalPartition) UpdateActivity() {
	p.lastActivityTime.Store(time.Now().UnixNano())
}

// IsIdle returns true if the partition has no publishers and no subscribers
func (p *LocalPartition) IsIdle() bool {
	return p.Publishers.Size() == 0 && p.Subscribers.Size() == 0
}

// GetIdleDuration returns how long the partition has been idle
func (p *LocalPartition) GetIdleDuration() time.Duration {
	lastActivity := p.lastActivityTime.Load()
	return time.Since(time.Unix(0, lastActivity))
}

// ShouldCleanup returns true if the partition should be cleaned up
// A partition should be cleaned up if:
// 1. It has no publishers and no subscribers
// 2. It has been idle for longer than the idle timeout
func (p *LocalPartition) ShouldCleanup(idleTimeout time.Duration) bool {
	if !p.IsIdle() {
		return false
	}
	return p.GetIdleDuration() > idleTimeout
}
