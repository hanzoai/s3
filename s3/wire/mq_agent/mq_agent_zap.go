// Code generated from mq_agent.proto; DO NOT EDIT.

package mq_agentwire

import (
	"fmt"

	"github.com/zap-proto/go/rpc"
)

// Method ordinals for the HanzoMessagingAgent service (stable 1-based wire ids,
// in mq_agent.proto declaration order).
const (
	HanzoMessagingAgentStartPublishSessionOrdinal uint32 = 1
	HanzoMessagingAgentClosePublishSessionOrdinal uint32 = 2
	HanzoMessagingAgentPublishRecordOrdinal       uint32 = 3 // STREAMING (bidi)
	HanzoMessagingAgentSubscribeRecordOrdinal     uint32 = 4 // STREAMING (bidi)
)

// HanzoMessagingAgentChannel ships one Call envelope and awaits its correlated
// Response.
type HanzoMessagingAgentChannel interface {
	Call(envelope []byte) (rpc.Response, error)
}

// HanzoMessagingAgentClient is a typed RPC client for the HanzoMessagingAgent
// service over a ZAP call channel. Each call takes a fresh PromiseID from sess;
// the pipelined "On" form of a method sets Target to a prior call's Promise so
// the server chains them.
type HanzoMessagingAgentClient struct {
	ch   HanzoMessagingAgentChannel
	cap  []byte
	sess *rpc.Session
}

// NewHanzoMessagingAgentClient returns a client that issues calls over ch,
// attaching cap (which may be nil) to every request.
func NewHanzoMessagingAgentClient(ch HanzoMessagingAgentChannel, capability []byte) *HanzoMessagingAgentClient {
	return &HanzoMessagingAgentClient{ch: ch, cap: capability, sess: rpc.NewSession()}
}

func (c *HanzoMessagingAgentClient) StartPublishSession(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeStartPublishSession(rpc.NoTarget, req)
}

// StartPublishSessionOn issues StartPublishSession as a dependent call pipelined
// on the answer of on: the server substitutes on's resolved result for this
// call's payload before dispatch, so it ships without waiting for on to
// round-trip.
func (c *HanzoMessagingAgentClient) StartPublishSessionOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeStartPublishSession(on.ID, nil)
}

func (c *HanzoMessagingAgentClient) invokeStartPublishSession(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    HanzoMessagingAgentStartPublishSessionOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("HanzoMessagingAgent.StartPublishSession: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

func (c *HanzoMessagingAgentClient) ClosePublishSession(req []byte) (rpc.Promise, []byte, error) {
	return c.invokeClosePublishSession(rpc.NoTarget, req)
}

// ClosePublishSessionOn issues ClosePublishSession as a dependent call pipelined
// on the answer of on: the server substitutes on's resolved result for this
// call's payload before dispatch, so it ships without waiting for on to
// round-trip.
func (c *HanzoMessagingAgentClient) ClosePublishSessionOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.invokeClosePublishSession(on.ID, nil)
}

func (c *HanzoMessagingAgentClient) invokeClosePublishSession(target uint32, payload []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    HanzoMessagingAgentClosePublishSessionOrdinal,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("HanzoMessagingAgent.ClosePublishSession: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

// PublishRecord sends one PublishRecordRequest frame and returns the correlated
// PublishRecordResponse body.
//
// STREAMING: PublishRecord is a bidirectional stream (stream
// PublishRecordRequest -> stream PublishRecordResponse). This per-frame send
// is the unit the full duplex stream is built from; the streaming transport
// primitive (multi-frame send/receive over one PromiseID) lands when it ships.
func (c *HanzoMessagingAgentClient) PublishRecord(req []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    HanzoMessagingAgentPublishRecordOrdinal,
		PromiseID: p.ID,
		Target:    rpc.NoTarget,
		Cap:       c.cap,
		Payload:   req,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("HanzoMessagingAgent.PublishRecord: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

// SubscribeRecord sends one SubscribeRecordRequest frame and returns the
// correlated SubscribeRecordResponse body.
//
// STREAMING: SubscribeRecord is a bidirectional stream (stream
// SubscribeRecordRequest -> stream SubscribeRecordResponse). This per-frame send
// is the unit the full duplex stream is built from; the streaming transport
// primitive (multi-frame send/receive over one PromiseID) lands when it ships.
func (c *HanzoMessagingAgentClient) SubscribeRecord(req []byte) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
		Method:    HanzoMessagingAgentSubscribeRecordOrdinal,
		PromiseID: p.ID,
		Target:    rpc.NoTarget,
		Cap:       c.cap,
		Payload:   req,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		return p, nil, fmt.Errorf("HanzoMessagingAgent.SubscribeRecord: status %d", resp.Status)
	}
	return p, resp.Body, nil
}

// HanzoMessagingAgentHandler is the server contract for the HanzoMessagingAgent
// service. Implement each method, then route requests to it with
// DispatchHanzoMessagingAgent. The two streaming methods are dispatched
// per-frame until the streaming transport primitive ships.
type HanzoMessagingAgentHandler interface {
	StartPublishSession(req []byte) ([]byte, error)
	ClosePublishSession(req []byte) ([]byte, error)
	// PublishRecord handles one PublishRecordRequest frame and returns one
	// PublishRecordResponse frame. STREAMING (bidi).
	PublishRecord(req []byte) ([]byte, error)
	// SubscribeRecord handles one SubscribeRecordRequest frame and returns one
	// SubscribeRecordResponse frame. STREAMING (bidi).
	SubscribeRecord(req []byte) ([]byte, error)
}

// DispatchHanzoMessagingAgent decodes a Call envelope, routes it by method
// ordinal to h, and returns the response envelope. An unknown ordinal yields a
// StatusNotFound response; a handler error yields StatusInternal. The two
// streaming methods route per-frame (one request frame -> one response frame)
// until the streaming transport primitive ships.
func DispatchHanzoMessagingAgent(h HanzoMessagingAgentHandler, envelope []byte) ([]byte, error) {
	call, err := rpc.ParseRequest(envelope)
	if err != nil {
		return nil, err
	}
	switch call.Method {
	case HanzoMessagingAgentStartPublishSessionOrdinal:
		body, err := h.StartPublishSession(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case HanzoMessagingAgentClosePublishSessionOrdinal:
		body, err := h.ClosePublishSession(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case HanzoMessagingAgentPublishRecordOrdinal:
		// STREAMING (bidi): per-frame dispatch until the streaming transport ships.
		body, err := h.PublishRecord(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case HanzoMessagingAgentSubscribeRecordOrdinal:
		// STREAMING (bidi): per-frame dispatch until the streaming transport ships.
		body, err := h.SubscribeRecord(call.Payload)
		if err != nil {
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	default:
		return rpc.BuildResponse(rpc.StatusNotFound, call.PromiseID, nil), nil
	}
}
