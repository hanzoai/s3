package agent_client

import (
	"fmt"

	"github.com/hanzoai/s3/s3/mq/agent/agentconv"
	"github.com/hanzoai/s3/s3/mq/schema"
	"github.com/hanzoai/s3/s3/pb/schema_pb"
	mq_agentwire "github.com/hanzoai/s3/s3/wire/mq_agent"
	mq_schemawire "github.com/hanzoai/s3/s3/wire/mq_schema"
	"github.com/zap-proto/go/transport"
)

type PublishSession struct {
	schema         *schema.Schema
	partitionCount int
	publisherName  string
	conn           transport.Conn
	stream         transport.Stream
}

func NewPublishSession(agentAddress string, topicSchema *schema.Schema, partitionCount int, publisherName string) (*PublishSession, error) {

	// call local agent ZAP server to create a new session
	conn, err := transport.Dial("tcp", agentAddress)
	if err != nil {
		return nil, fmt.Errorf("dial agent server %s: %v", agentAddress, err)
	}
	agentClient := mq_agentwire.NewHanzoMessagingAgentClient(conn, nil)

	startReq := mq_agentwire.NewStartPublishSessionRequest(mq_agentwire.StartPublishSessionRequestInput{
		Topic: mq_schemawire.NewTopic(mq_schemawire.TopicInput{
			Namespace: topicSchema.Namespace,
			Name:      topicSchema.Name,
		}),
		PartitionCount: int32(partitionCount),
		RecordType:     agentconv.RecordTypeToWire(topicSchema.RecordType),
		PublisherName:  publisherName,
	})
	_, body, err := agentClient.StartPublishSession(startReq)
	if err != nil {
		conn.Close()
		return nil, err
	}
	resp, err := mq_agentwire.WrapStartPublishSessionResponse(body)
	if err != nil {
		conn.Close()
		return nil, err
	}
	if resp.Error() != "" {
		conn.Close()
		return nil, fmt.Errorf("start publish session: %v", resp.Error())
	}

	// open the PublishRecord bidirectional stream; the first frame carries the
	// session id.
	stream, err := conn.OpenStream(mq_agentwire.HanzoMessagingAgentPublishRecordOrdinal,
		mq_agentwire.NewPublishRecordRequest(mq_agentwire.PublishRecordRequestInput{
			SessionId: resp.SessionId(),
		}))
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("publish record: %w", err)
	}

	return &PublishSession{
		schema:         topicSchema,
		partitionCount: partitionCount,
		publisherName:  publisherName,
		conn:           conn,
		stream:         stream,
	}, nil
}

func (a *PublishSession) CloseSession() error {
	if a.schema == nil {
		return nil
	}
	err := a.stream.CloseSend()
	if err != nil {
		return fmt.Errorf("close send: %w", err)
	}
	a.schema = nil
	if a.conn != nil {
		a.conn.Close()
	}
	return err
}

func (a *PublishSession) PublishMessageRecord(key []byte, record *schema_pb.RecordValue) error {
	return a.stream.Send(mq_agentwire.NewPublishRecordRequest(mq_agentwire.PublishRecordRequestInput{
		Key:   key,
		Value: agentconv.RecordValueToWire(record),
	}))
}
