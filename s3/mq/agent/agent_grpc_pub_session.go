package agent

import (
	"log/slog"
	"math/rand/v2"

	"github.com/hanzoai/s3/s3/mq/agent/agentconv"
	"github.com/hanzoai/s3/s3/mq/client/pub_client"
	"github.com/hanzoai/s3/s3/mq/topic"
	mq_agentwire "github.com/hanzoai/s3/s3/wire/mq_agent"
)

// StartPublishSession implements the mq_agentwire.HanzoMessagingAgent unary
// method over the ZAP transport: it decodes the request envelope, opens a
// publisher session, and returns the response envelope.
func (a *MessageQueueAgent) StartPublishSession(req []byte) ([]byte, error) {
	request, err := mq_agentwire.WrapStartPublishSessionRequest(req)
	if err != nil {
		return nil, err
	}
	pbTopic := agentconv.TopicFromWire(request.Topic())

	sessionId := rand.Int64()

	topicPublisher, err := pub_client.NewTopicPublisher(
		&pub_client.PublisherConfiguration{
			Topic:          topic.NewTopic(pbTopic.Namespace, pbTopic.Name),
			PartitionCount: request.PartitionCount(),
			Brokers:        a.brokersList(),
			PublisherName:  request.PublisherName(),
			RecordType:     agentconv.RecordTypeFromWire(request.RecordType()),
		})
	if err != nil {
		return nil, err
	}

	a.publishersLock.Lock()
	a.publishers[SessionId(sessionId)] = &SessionEntry[*pub_client.TopicPublisher]{
		entry: topicPublisher,
	}
	a.publishersLock.Unlock()

	return mq_agentwire.NewStartPublishSessionResponse(mq_agentwire.StartPublishSessionResponseInput{
		SessionId: sessionId,
	}), nil
}

// ClosePublishSession implements the mq_agentwire.HanzoMessagingAgent unary
// method over the ZAP transport.
func (a *MessageQueueAgent) ClosePublishSession(req []byte) ([]byte, error) {
	request, err := mq_agentwire.WrapClosePublishSessionRequest(req)
	if err != nil {
		return nil, err
	}
	sessionId := SessionId(request.SessionId())

	var finishErr string
	a.publishersLock.Lock()
	publisherEntry, found := a.publishers[sessionId]
	if found {
		if err := publisherEntry.entry.FinishPublish(); err != nil {
			finishErr = err.Error()
			slog.Warn("failed to finish publish", "error", err)
		}
		delete(a.publishers, sessionId)
	}
	a.publishersLock.Unlock()

	return mq_agentwire.NewClosePublishSessionResponse(mq_agentwire.ClosePublishSessionResponseInput{
		Error: finishErr,
	}), nil
}
