package agent

import (
	"fmt"

	"github.com/hanzoai/s3/s3/mq/agent/agentconv"
	mq_agentwire "github.com/hanzoai/s3/s3/wire/mq_agent"
	"github.com/zap-proto/go/transport"
)

// PublishRecord serves the HanzoMessagingAgent PublishRecord bidirectional
// stream over the ZAP transport: the first frame carries the session id, each
// subsequent frame a key/value record to publish. It mirrors the prior gRPC
// stream loop, reading PublishRecordRequest envelopes off s until the peer
// half-closes (io.EOF).
func (a *MessageQueueAgent) PublishRecord(init []byte, s *transport.Stream) error {
	m, err := mq_agentwire.WrapPublishRecordRequest(init)
	if err != nil {
		return err
	}
	sessionId := SessionId(m.SessionId())
	a.publishersLock.RLock()
	publisherEntry, found := a.publishers[sessionId]
	a.publishersLock.RUnlock()
	if !found {
		return fmt.Errorf("publish session id %d not found", sessionId)
	}
	defer func() {
		a.publishersLock.Lock()
		delete(a.publishers, sessionId)
		a.publishersLock.Unlock()
	}()

	if value := m.Value(); value != nil {
		if err := publisherEntry.entry.PublishRecord(m.Key(), agentconv.RecordValueFromWire(value)); err != nil {
			return err
		}
	}

	for {
		frame, err := s.Recv()
		if err != nil {
			return err
		}
		m, err := mq_agentwire.WrapPublishRecordRequest(frame)
		if err != nil {
			return err
		}
		value := m.Value()
		if value == nil {
			continue
		}
		if err := publisherEntry.entry.PublishRecord(m.Key(), agentconv.RecordValueFromWire(value)); err != nil {
			return err
		}
	}
}
