package agent_client

import (
	"fmt"

	"github.com/hanzoai/s3/s3/mq/agent/agentconv"
	"github.com/hanzoai/s3/s3/mq/topic"
	"github.com/hanzoai/s3/s3/pb/schema_pb"
	mq_agentwire "github.com/hanzoai/s3/s3/wire/mq_agent"
	mq_schemawire "github.com/hanzoai/s3/s3/wire/mq_schema"
	"github.com/zap-proto/go/transport"
)

type SubscribeOption struct {
	ConsumerGroup           string
	ConsumerGroupInstanceId string
	Topic                   topic.Topic
	OffsetType              schema_pb.OffsetType
	OffsetTsNs              int64
	Filter                  string
	MaxSubscribedPartitions int32
	SlidingWindowSize       int32
}

type SubscribeSession struct {
	Option *SubscribeOption
	conn   transport.Conn
	stream transport.Stream
}

func NewSubscribeSession(agentAddress string, option *SubscribeOption) (*SubscribeSession, error) {
	// call local agent ZAP server to create a new session
	conn, err := transport.Dial("tcp", agentAddress)
	if err != nil {
		return nil, fmt.Errorf("dial agent server %s: %v", agentAddress, err)
	}

	initRequest := mq_agentwire.NewInitSubscribeRecordRequest(mq_agentwire.InitSubscribeRecordRequestInput{
		ConsumerGroup:           option.ConsumerGroup,
		ConsumerGroupInstanceId: option.ConsumerGroupInstanceId,
		Topic: mq_schemawire.NewTopic(mq_schemawire.TopicInput{
			Namespace: option.Topic.Namespace,
			Name:      option.Topic.Name,
		}),
		OffsetType:              uint32(option.OffsetType),
		OffsetTsNs:              option.OffsetTsNs,
		MaxSubscribedPartitions: option.MaxSubscribedPartitions,
		Filter:                  option.Filter,
		SlidingWindowSize:       option.SlidingWindowSize,
	})

	// open the SubscribeRecord bidirectional stream; the first frame carries
	// the init request.
	stream, err := conn.OpenStream(mq_agentwire.HanzoMessagingAgentSubscribeRecordOrdinal,
		mq_agentwire.NewSubscribeRecordRequest(mq_agentwire.SubscribeRecordRequestInput{
			Init: initRequest,
		}))
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("subscribe record: %w", err)
	}

	return &SubscribeSession{
		Option: option,
		conn:   conn,
		stream: stream,
	}, nil
}

func (s *SubscribeSession) CloseSession() error {
	err := s.stream.CloseSend()
	if s.conn != nil {
		s.conn.Close()
	}
	return err
}

func (a *SubscribeSession) SubscribeMessageRecord(
	onEachMessageFn func(key []byte, record *schema_pb.RecordValue),
	onCompletionFn func()) error {
	for {
		frame, err := a.stream.Recv()
		if err != nil {
			if onCompletionFn != nil {
				onCompletionFn()
			}
			return err
		}
		resp, err := mq_agentwire.WrapSubscribeRecordResponse(frame)
		if err != nil {
			if onCompletionFn != nil {
				onCompletionFn()
			}
			return err
		}
		onEachMessageFn(resp.Key(), agentconv.RecordValueFromWire(resp.Value()))
	}
}
