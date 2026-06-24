// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// Native ZAP service adapter for HanzoMessagingAgent — the hand-written glue
// that binds the generated mq_agent wire messages to a real message-queue
// backend over the canonical github.com/zap-proto/go transport. This is the
// mq_agent analogue of s3/zapsvc/service.go: it kills gRPC for the agent's
// four RPCs and replaces them with the SAME wire hanzo and lux share.
//
// Shape:
//   - Unary RPCs (StartPublishSession, ClosePublishSession) mirror the object
//     service's Get/Put: Wrap the request view, call the backend, build the
//     reply buffer. Dispatched by the generated ordinal table over a plain
//     Call/Response envelope.
//   - Streaming RPCs (PublishRecord, SubscribeRecord) are bidirectional in the
//     proto (`stream … returns stream …`). zapgen cannot express streams, so —
//     exactly like s3/admin/plugin/worker_stream_zap.go — they are hand-wired
//     over transport.OpenStream / transport.StreamHandler / *transport.Stream.
//     Each stream frame is a zero-copy New*/Wrap* buffer (same doctrine as
//     unary); the opening request rides as the stream's init payload and is
//     replayed as the server's first Recv. NOTHING is faked or buffered into a
//     single round-trip.
//
// The backend is a Go interface (MessagingAgentStore) defined in agent terms,
// NOT an import of the real filer/broker/publisher engine — integration wires
// the concrete *agent.MessageQueueAgent to it later (see the package doc and
// the integration notes returned with this change).

package mq_agentwire

import (
	"github.com/zap-proto/go/rpc"
	"github.com/zap-proto/go/transport"
)

// MessagingAgentStore is the backend the ZAP agent service delegates to. It is
// expressed purely in raw ZAP message buffers (the New*/Wrap* doctrine): the
// real engine (s3/mq/agent.MessageQueueAgent) implements it; tests use an
// in-memory fake.
//
// Unary methods take/return one message buffer. The streaming methods take a
// typed stream the handler drives until the peer half-closes — they NEVER
// return a single buffer (that would collapse the stream).
type MessagingAgentStore interface {
	// StartPublishSession opens a publisher session for a topic. req is a
	// StartPublishSessionRequest buffer; the return is a
	// StartPublishSessionResponse buffer.
	StartPublishSession(req []byte) (resp []byte, err error)
	// ClosePublishSession closes a publisher session. req is a
	// ClosePublishSessionRequest buffer; the return is a
	// ClosePublishSessionResponse buffer.
	ClosePublishSession(req []byte) (resp []byte, err error)
	// PublishRecord drives one bidirectional PublishRecord stream: the client
	// streams PublishRecordRequest frames (the first carries the session id),
	// the server streams PublishRecordResponse acks. Returning ends the stream.
	PublishRecord(stream PublishRecordServerStream) error
	// SubscribeRecord drives one bidirectional SubscribeRecord stream: the
	// client sends an init frame then ack frames, the server streams
	// SubscribeRecordResponse data frames. Returning ends the stream.
	SubscribeRecord(stream SubscribeRecordServerStream) error
}

// --- unary dispatch (StartPublishSession, ClosePublishSession) ---

// unaryHandler adapts a MessagingAgentStore to the two unary ordinals. The
// streaming ordinals (3, 4) are NOT served here — they arrive as stream-open
// frames routed to StreamHandler, never as Call envelopes — so an unexpected
// unary call on them yields StatusNotFound.
type unaryHandler struct{ store MessagingAgentStore }

func (h unaryHandler) dispatch(envelope []byte) ([]byte, error) {
	call, err := rpc.ParseRequest(envelope)
	if err != nil {
		return nil, err
	}
	switch call.Method {
	case HanzoMessagingAgentStartPublishSessionOrdinal:
		body, err := h.store.StartPublishSession(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	case HanzoMessagingAgentClosePublishSessionOrdinal:
		body, err := h.store.ClosePublishSession(call.Payload)
		if err != nil {
			// Carry the handler error message so the caller can reconstruct sentinels.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(err.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	default:
		// PublishRecord/SubscribeRecord are streaming-only; anything else is unknown.
		return rpc.BuildResponse(rpc.StatusNotFound, call.PromiseID, nil), nil
	}
}

// Dispatch is the unary HanzoMessagingAgent server dispatch bound to store —
// pass it as the unary leg of transport.ListenStream. It serves only
// StartPublishSession and ClosePublishSession; the streaming methods are served
// by StreamHandler.
func Dispatch(store MessagingAgentStore) transport.Dispatch {
	h := unaryHandler{store: store}
	return h.dispatch
}

// --- streaming server side (PublishRecord, SubscribeRecord) ---

// PublishRecordServerStream is the server view of one PublishRecord stream:
// receive PublishRecordRequest frames, send PublishRecordResponse acks. The
// first request frame arrives as the stream's init payload and is replayed as
// the first Recv (mirroring s3/admin/plugin/worker_stream_zap.go).
type PublishRecordServerStream interface {
	// Recv returns the next PublishRecordRequest buffer, or io.EOF once the
	// client half-closes. Re-wrap with WrapPublishRecordRequest.
	Recv() ([]byte, error)
	// Send streams one PublishRecordResponse buffer (from NewPublishRecordResponse).
	Send(resp []byte) error
}

// SubscribeRecordServerStream is the server view of one SubscribeRecord stream:
// receive SubscribeRecordRequest frames (first = init, then acks), send
// SubscribeRecordResponse data frames.
type SubscribeRecordServerStream interface {
	// Recv returns the next SubscribeRecordRequest buffer, or io.EOF once the
	// client half-closes. Re-wrap with WrapSubscribeRecordRequest.
	Recv() ([]byte, error)
	// Send streams one SubscribeRecordResponse buffer (from NewSubscribeRecordResponse).
	Send(resp []byte) error
}

// serverStream adapts a *transport.Stream to the *ServerStream interfaces.
// Both streaming methods share the same frame plumbing (raw ZAP buffers in/out);
// only the typed wrappers around the buffers differ, so one adapter serves both.
// The opening request frame rides as init and is replayed as the first Recv.
type serverStream struct {
	s     *transport.Stream
	init  []byte
	hasIn bool
}

func newServerStream(init []byte, s *transport.Stream) *serverStream {
	return &serverStream{s: s, init: init, hasIn: len(init) > 0}
}

func (z *serverStream) Recv() ([]byte, error) {
	if z.hasIn {
		z.hasIn = false
		return z.init, nil
	}
	return z.s.Recv()
}

func (z *serverStream) Send(body []byte) error { return z.s.Send(body) }

// StreamHandler is the transport.StreamHandler for HanzoMessagingAgent's two
// streaming ordinals: it adapts each accepted ZAP stream to the matching
// server-stream interface and runs the backend's handler against it. Pass it as
// the stream leg of transport.ListenStream alongside Dispatch(store).
//
// Returning from the backend handler half-closes the send side (the client's
// Recv then sees io.EOF), per transport.StreamHandler semantics.
func StreamHandler(store MessagingAgentStore) transport.StreamHandler {
	return func(method uint32, init []byte, s *transport.Stream) {
		switch method {
		case HanzoMessagingAgentPublishRecordOrdinal:
			_ = store.PublishRecord(newServerStream(init, s))
		case HanzoMessagingAgentSubscribeRecordOrdinal:
			_ = store.SubscribeRecord(newServerStream(init, s))
		default:
			// Unknown stream ordinal: returning half-closes; client sees io.EOF.
		}
	}
}

// Serve starts the native ZAP HanzoMessagingAgent service on network/addr (e.g.
// "tcp", ":16777"), backed by store. Unary calls (StartPublishSession,
// ClosePublishSession) route through Dispatch; the two bidi streams route
// through StreamHandler — both over ONE listener, exactly like the worker ZAP
// server (s3/admin/dash/worker_grpc_server.go). Returns the running server;
// Close stops it.
func Serve(network, addr string, store MessagingAgentStore) (*transport.Server, error) {
	return transport.ListenStream(network, addr, Dispatch(store), StreamHandler(store))
}

// --- client side ---

// Client is the typed HanzoMessagingAgent ZAP client internal services hold. It
// owns the transport connection to one agent endpoint and serves both the unary
// calls and the streaming opens over it.
type Client struct {
	conn *transport.Conn
	rpc  *HanzoMessagingAgentClient
}

// Dial connects to the agent ZAP service at addr over plain TCP (e.g.
// "mq-agent.hanzo.svc:16777"). For the PQ-secured mesh, build the
// *transport.Conn via the TLS/QUIC path and use NewClient.
func Dial(network, addr string) (*Client, error) {
	conn, err := transport.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn, rpc: NewHanzoMessagingAgentClient(conn, nil)}, nil
}

// NewClient wraps an already-established transport.Conn (TCP, Unix, or PQ-TLS).
func NewClient(conn *transport.Conn) *Client {
	return &Client{conn: conn, rpc: NewHanzoMessagingAgentClient(conn, nil)}
}

// Close releases the underlying connection.
func (c *Client) Close() error { return c.conn.Close() }

// StartPublishSession opens a publisher session over ZAP. in/out are the typed
// message inputs/views; the bytes never leave this package as structs.
func (c *Client) StartPublishSession(in StartPublishSessionRequestInput) (StartPublishSessionResponse, error) {
	_, body, err := c.rpc.StartPublishSession(NewStartPublishSessionRequest(in))
	if err != nil {
		return StartPublishSessionResponse{}, err
	}
	return WrapStartPublishSessionResponse(body)
}

// ClosePublishSession closes a publisher session over ZAP.
func (c *Client) ClosePublishSession(in ClosePublishSessionRequestInput) (ClosePublishSessionResponse, error) {
	_, body, err := c.rpc.ClosePublishSession(NewClosePublishSessionRequest(in))
	if err != nil {
		return ClosePublishSessionResponse{}, err
	}
	return WrapClosePublishSessionResponse(body)
}

// PublishRecordClientStream is the client view of one PublishRecord stream: send
// PublishRecordRequest frames, receive PublishRecordResponse acks, half-close
// when done.
type PublishRecordClientStream struct{ s *transport.Stream }

// Send streams one PublishRecordRequest frame.
func (p *PublishRecordClientStream) Send(in PublishRecordRequestInput) error {
	return p.s.Send(NewPublishRecordRequest(in))
}

// Recv returns the next PublishRecordResponse ack, or io.EOF once the server
// half-closes.
func (p *PublishRecordClientStream) Recv() (PublishRecordResponse, error) {
	body, err := p.s.Recv()
	if err != nil {
		return PublishRecordResponse{}, err
	}
	return WrapPublishRecordResponse(body)
}

// CloseSend half-closes the outbound half; the server's Recv then sees io.EOF.
func (p *PublishRecordClientStream) CloseSend() error { return p.s.CloseSend() }

// PublishRecord opens a PublishRecord bidi stream. The first record (carrying
// the session id) rides as the stream's init payload; subsequent records go via
// Send on the returned stream.
func (c *Client) PublishRecord(first PublishRecordRequestInput) (*PublishRecordClientStream, error) {
	s, err := c.conn.OpenStream(HanzoMessagingAgentPublishRecordOrdinal, NewPublishRecordRequest(first))
	if err != nil {
		return nil, err
	}
	return &PublishRecordClientStream{s: s}, nil
}

// SubscribeRecordClientStream is the client view of one SubscribeRecord stream:
// send SubscribeRecordRequest ack frames, receive SubscribeRecordResponse data
// frames.
type SubscribeRecordClientStream struct{ s *transport.Stream }

// Send streams one SubscribeRecordRequest frame (an ack after the init).
func (p *SubscribeRecordClientStream) Send(in SubscribeRecordRequestInput) error {
	return p.s.Send(NewSubscribeRecordRequest(in))
}

// Recv returns the next SubscribeRecordResponse data frame, or io.EOF once the
// server half-closes.
func (p *SubscribeRecordClientStream) Recv() (SubscribeRecordResponse, error) {
	body, err := p.s.Recv()
	if err != nil {
		return SubscribeRecordResponse{}, err
	}
	return WrapSubscribeRecordResponse(body)
}

// CloseSend half-closes the outbound half; the server's Recv then sees io.EOF.
func (p *SubscribeRecordClientStream) CloseSend() error { return p.s.CloseSend() }

// SubscribeRecord opens a SubscribeRecord bidi stream. The init request (the
// InitSubscribeRecordRequest, wrapped in a SubscribeRecordRequest) rides as the
// stream's init payload; subsequent ack frames go via Send.
func (c *Client) SubscribeRecord(init SubscribeRecordRequestInput) (*SubscribeRecordClientStream, error) {
	s, err := c.conn.OpenStream(HanzoMessagingAgentSubscribeRecordOrdinal, NewSubscribeRecordRequest(init))
	if err != nil {
		return nil, err
	}
	return &SubscribeRecordClientStream{s: s}, nil
}
