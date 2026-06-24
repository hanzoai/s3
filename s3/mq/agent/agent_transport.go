package agent

import (
	"fmt"

	"github.com/hanzoai/s3/s3/glog"
	mq_agentwire "github.com/hanzoai/s3/s3/wire/mq_agent"
	"github.com/zap-proto/go/transport"
)

// unaryHandler adapts MessageQueueAgent to the wire
// HanzoMessagingAgentHandler contract for the UNARY methods only. The two
// streaming methods are served by MessageQueueAgent.StreamHandler over
// transport.ListenStream and never reach this dispatcher, so their per-frame
// forms here are unreachable guards.
type unaryHandler struct{ a *MessageQueueAgent }

func (h unaryHandler) StartPublishSession(req []byte) ([]byte, error) {
	return h.a.StartPublishSession(req)
}

func (h unaryHandler) ClosePublishSession(req []byte) ([]byte, error) {
	return h.a.ClosePublishSession(req)
}

func (h unaryHandler) PublishRecord(req []byte) ([]byte, error) {
	return nil, fmt.Errorf("PublishRecord is a streaming method; route it through the stream transport")
}

func (h unaryHandler) SubscribeRecord(req []byte) ([]byte, error) {
	return nil, fmt.Errorf("SubscribeRecord is a streaming method; route it through the stream transport")
}

// DispatchUnary routes the agent's unary methods (StartPublishSession,
// ClosePublishSession) off a ZAP request envelope. The streaming methods are
// served by StreamHandler over transport.ListenStream.
func (a *MessageQueueAgent) DispatchUnary(envelope []byte) ([]byte, error) {
	return mq_agentwire.DispatchHanzoMessagingAgent(unaryHandler{a}, envelope)
}

// StreamHandler routes the agent's two bidirectional streaming methods by
// method ordinal to their real stream loops. It satisfies
// transport.StreamHandler.
func (a *MessageQueueAgent) StreamHandler(method uint32, init []byte, s *transport.Stream) {
	switch method {
	case mq_agentwire.HanzoMessagingAgentPublishRecordOrdinal:
		if err := a.PublishRecord(init, s); err != nil {
			glog.V(0).Infof("publish record stream: %v", err)
		}
	case mq_agentwire.HanzoMessagingAgentSubscribeRecordOrdinal:
		if err := a.SubscribeRecord(init, s); err != nil {
			glog.V(0).Infof("subscribe record stream: %v", err)
		}
	}
}
