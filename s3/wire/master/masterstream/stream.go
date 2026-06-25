// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// Package masterstream is the native ZAP streaming adapter for the 3 streaming
// Hanzo (master) RPCs. It is to the generated s3/wire/master (masterwire)
// package what s3/wire/filer/filerstream is to filerwire: the ONLY hand-written
// glue binding the zero-copy wire messages to a real transport stream. No gRPC,
// no protobuf framing, no struct marshaling — every frame that crosses the
// socket is a masterwire New*/Wrap* buffer framed as a stream message.
//
// zapgen does NOT emit streaming (its method grammar has no `stream` keyword),
// so the generated masterwire service exposes only single-frame unary stubs for
// the streaming ordinals, and masterwire.Dispatch short-circuits them with
// ErrStreamingNotWired rather than degrading a stream to one request->one
// response. This adapter wires the REAL duplex primitive shipped by
// github.com/zap-proto/go/transport (Conn.OpenStream on the client,
// transport.ListenStream + a StreamHandler on the server, Stream.Send/Recv on
// both ends). The reserved masterwire ordinals are reused as each stream's
// open-method id, so the wire numbering stays stable and the unary stub and this
// streaming path agree on every method number.
//
// All 3 streaming RPCs are bidirectional (stream Request <-> stream Response):
//
//	SendHeartbeat  bidirectional  (Heartbeat*            <-> HeartbeatResponse*)
//	KeepConnected  bidirectional  (KeepConnectedRequest* <-> KeepConnectedResponse*)
//	StreamAssign   bidirectional  (AssignRequest*        <-> AssignResponse*)
//
// The backend (Server below) is defined in master terms and does NOT import the
// real master engine — the integrator implements it against the live topology;
// tests use an in-memory stub (see stream_roundtrip_test.go).

package masterstream

import (
	"context"
	"io"

	masterwire "github.com/hanzoai/s3/s3/wire/master"
	"github.com/zap-proto/go/transport"
)

// Ordinals each streaming RPC opens with — the SAME ids the generated unary
// stubs reserve, so the streaming open frame and any unary probe of the method
// agree on the wire number.
const (
	SendHeartbeatOrdinal = masterwire.HanzoSendHeartbeatOrdinal
	KeepConnectedOrdinal = masterwire.HanzoKeepConnectedOrdinal
	StreamAssignOrdinal  = masterwire.HanzoStreamAssignOrdinal
)

// IsEOF reports whether err is the stream's half-close sentinel, so callers can
// end a Recv loop without importing io.
func IsEOF(err error) bool { return err == io.EOF }

// ---------------------------------------------------------------------------
// Server side
// ---------------------------------------------------------------------------

// Server is the server-side backend the streaming adapter delegates to: the
// real master engine implements it; tests use an in-memory stub. Each method
// handles the full lifetime of one bidirectional stream of its RPC: init is the
// first request frame (already wrapped), s.Recv yields the rest, s.Send ships
// response frames; returning ends the stream (the peer's Recv then sees io.EOF).
//
// A returned error just ends the stream (the transport half-closes the send
// side); per-frame errors are carried in the response messages' own error
// fields, not at the transport.
type Server interface {
	// SendHeartbeat serves one volume server's heartbeat stream.
	SendHeartbeat(init masterwire.Heartbeat, s *HeartbeatStream) error
	// KeepConnected serves one client's keep-connected subscription stream.
	KeepConnected(init masterwire.KeepConnectedRequest, s *KeepConnectedStream) error
	// StreamAssign serves one client's batched assign stream.
	StreamAssign(init masterwire.AssignRequest, s *AssignStream) error
}

// HeartbeatStream is the server's half of the bidirectional SendHeartbeat
// stream: Recv yields successive Heartbeat frames, Send ships HeartbeatResponse
// frames.
type HeartbeatStream struct{ s transport.Stream }

// Context is cancelled when this stream ends or the connection drops.
func (x *HeartbeatStream) Context() context.Context { return x.s.Context() }

// Recv returns the next request frame as a zero-copy view, or io.EOF once the
// client half-closes its send side.
func (x *HeartbeatStream) Recv() (masterwire.Heartbeat, error) {
	b, err := x.s.Recv()
	if err != nil {
		return masterwire.Heartbeat{}, err
	}
	return masterwire.WrapHeartbeat(b)
}

// RecvBytes returns the next request frame as the raw ZAP buffer.
func (x *HeartbeatStream) RecvBytes() ([]byte, error) { return x.s.Recv() }

// Send ships one HeartbeatResponse frame (the bytes ARE the message).
func (x *HeartbeatStream) Send(frame []byte) error { return x.s.Send(frame) }

// CloseSend half-closes the response direction.
func (x *HeartbeatStream) CloseSend() error { return x.s.CloseSend() }

// KeepConnectedStream is the server's half of the bidirectional KeepConnected
// stream: Recv yields KeepConnectedRequest frames, Send ships
// KeepConnectedResponse frames.
type KeepConnectedStream struct{ s transport.Stream }

// Context is cancelled when this stream ends or the connection drops.
func (x *KeepConnectedStream) Context() context.Context { return x.s.Context() }

// Recv returns the next request frame as a zero-copy view, or io.EOF once the
// client half-closes its send side.
func (x *KeepConnectedStream) Recv() (masterwire.KeepConnectedRequest, error) {
	b, err := x.s.Recv()
	if err != nil {
		return masterwire.KeepConnectedRequest{}, err
	}
	return masterwire.WrapKeepConnectedRequest(b)
}

// RecvBytes returns the next request frame as the raw ZAP buffer.
func (x *KeepConnectedStream) RecvBytes() ([]byte, error) { return x.s.Recv() }

// Send ships one KeepConnectedResponse frame.
func (x *KeepConnectedStream) Send(frame []byte) error { return x.s.Send(frame) }

// CloseSend half-closes the response direction.
func (x *KeepConnectedStream) CloseSend() error { return x.s.CloseSend() }

// AssignStream is the server's half of the bidirectional StreamAssign stream:
// Recv yields AssignRequest frames, Send ships AssignResponse frames.
type AssignStream struct{ s transport.Stream }

// Context is cancelled when this stream ends or the connection drops.
func (x *AssignStream) Context() context.Context { return x.s.Context() }

// Recv returns the next request frame as a zero-copy view, or io.EOF once the
// client half-closes its send side.
func (x *AssignStream) Recv() (masterwire.AssignRequest, error) {
	b, err := x.s.Recv()
	if err != nil {
		return masterwire.AssignRequest{}, err
	}
	return masterwire.WrapAssignRequest(b)
}

// RecvBytes returns the next request frame as the raw ZAP buffer.
func (x *AssignStream) RecvBytes() ([]byte, error) { return x.s.Recv() }

// Send ships one AssignResponse frame.
func (x *AssignStream) Send(frame []byte) error { return x.s.Send(frame) }

// CloseSend half-closes the response direction.
func (x *AssignStream) CloseSend() error { return x.s.CloseSend() }

// Handler builds the transport StreamHandler that routes an inbound stream-open
// to the matching Server method by ordinal. Pass it to transport.ListenStream
// (or use Serve, which does that). An open for an ordinal that is not one of the
// 3 streaming RPCs is ignored (the stream half-closes immediately). A malformed
// init buffer ends the stream without invoking the backend.
func Handler(srv Server) transport.StreamHandler {
	return func(method uint32, init []byte, s transport.Stream) {
		switch method {
		case SendHeartbeatOrdinal:
			v, err := masterwire.WrapHeartbeat(init)
			if err != nil {
				return
			}
			_ = srv.SendHeartbeat(v, &HeartbeatStream{s: s})
		case KeepConnectedOrdinal:
			v, err := masterwire.WrapKeepConnectedRequest(init)
			if err != nil {
				return
			}
			_ = srv.KeepConnected(v, &KeepConnectedStream{s: s})
		case StreamAssignOrdinal:
			v, err := masterwire.WrapAssignRequest(init)
			if err != nil {
				return
			}
			_ = srv.StreamAssign(v, &AssignStream{s: s})
		default:
			return // unknown stream method: returning half-closes + ends it
		}
	}
}

// StreamServer is the running native ZAP Hanzo streaming endpoint. Close stops
// it and tears down live streams.
type StreamServer struct{ srv *transport.Server }

// Serve starts the native ZAP Hanzo streaming endpoint on network/addr (e.g.
// "tcp", ":19333"), backed by srv. It registers NO unary dispatch (the unary
// RPCs are served by masterwire.Serve); inbound stream-opens route to srv via
// Handler.
//
// To serve BOTH the unary and streaming RPCs on ONE listener, call
// transport.ListenStream(network, addr, masterwire.Dispatch(backend), Handler(srv))
// directly instead of using two separate listeners. For the PQ-secured mesh,
// establish the listener over the transport TLS/QUIC path and call
// transport.ListenStream on it.
func Serve(network, addr string, srv Server) (*StreamServer, error) {
	s, err := transport.ListenStream(network, addr, nil, Handler(srv))
	if err != nil {
		return nil, err
	}
	return &StreamServer{srv: s}, nil
}

// Addr is the listener address (useful with ":0").
func (s *StreamServer) Addr() string { return s.srv.Addr().String() }

// Close stops the server.
func (s *StreamServer) Close() error { return s.srv.Close() }

// ---------------------------------------------------------------------------
// Client side
// ---------------------------------------------------------------------------

// Client is the caller's typed connection to one master streaming endpoint. It
// owns the transport connection; it may open multiple concurrent streams (one
// per outstanding streaming RPC) over the single connection.
type Client struct{ conn transport.Conn }

// Dial connects to the master streaming endpoint at addr over plain TCP. For the
// PQ-secured mesh use transport.DialTLS and NewClient.
func Dial(network, addr string) (*Client, error) {
	conn, err := transport.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn}, nil
}

// NewClient wraps an already-established transport.Conn (TCP, Unix, or PQ-TLS).
// The conn's read loop routes stream frames, so a plain call-only Dial conn is
// sufficient. This also lets a caller open streams over the SAME connection a
// masterwire.Client uses for unary calls (pass masterwire.Client's conn).
func NewClient(conn transport.Conn) *Client { return &Client{conn: conn} }

// Close releases the underlying connection.
func (c *Client) Close() error { return c.conn.Close() }

// SendHeartbeat opens the bidirectional SendHeartbeat stream, sending first (a
// masterwire NewHeartbeat buffer) as the opening frame. The returned stream
// Sends further request frames and Recvs response frames.
func (c *Client) SendHeartbeat(first []byte) (*ClientHeartbeatStream, error) {
	s, err := c.conn.OpenStream(SendHeartbeatOrdinal, first)
	if err != nil {
		return nil, err
	}
	return &ClientHeartbeatStream{s: s}, nil
}

// KeepConnected opens the bidirectional KeepConnected stream, sending first (a
// masterwire NewKeepConnectedRequest buffer) as the opening frame.
func (c *Client) KeepConnected(first []byte) (*ClientKeepConnectedStream, error) {
	s, err := c.conn.OpenStream(KeepConnectedOrdinal, first)
	if err != nil {
		return nil, err
	}
	return &ClientKeepConnectedStream{s: s}, nil
}

// StreamAssign opens the bidirectional StreamAssign stream, sending first (a
// masterwire NewAssignRequest buffer) as the opening frame.
func (c *Client) StreamAssign(first []byte) (*ClientAssignStream, error) {
	s, err := c.conn.OpenStream(StreamAssignOrdinal, first)
	if err != nil {
		return nil, err
	}
	return &ClientAssignStream{s: s}, nil
}

// ClientHeartbeatStream is the caller's half of the bidirectional SendHeartbeat
// stream: Send ships Heartbeat frames, Recv yields HeartbeatResponse frames.
type ClientHeartbeatStream struct{ s transport.Stream }

// Send ships one Heartbeat frame (a masterwire NewHeartbeat buffer).
func (s *ClientHeartbeatStream) Send(frame []byte) error { return s.s.Send(frame) }

// Recv returns the next HeartbeatResponse as a zero-copy view, or io.EOF once
// the server ends the stream.
func (s *ClientHeartbeatStream) Recv() (masterwire.HeartbeatResponse, error) {
	b, err := s.s.Recv()
	if err != nil {
		return masterwire.HeartbeatResponse{}, err
	}
	return masterwire.WrapHeartbeatResponse(b)
}

// RecvBytes returns the next response frame as the raw ZAP buffer.
func (s *ClientHeartbeatStream) RecvBytes() ([]byte, error) { return s.s.Recv() }

// CloseSend half-closes the request direction (no more heartbeats); the server's
// Recv then returns io.EOF. Idempotent.
func (s *ClientHeartbeatStream) CloseSend() error { return s.s.CloseSend() }

// ClientKeepConnectedStream is the caller's half of the bidirectional
// KeepConnected stream: Send ships requests, Recv yields responses.
type ClientKeepConnectedStream struct{ s transport.Stream }

// Send ships one KeepConnectedRequest frame.
func (s *ClientKeepConnectedStream) Send(frame []byte) error { return s.s.Send(frame) }

// Recv returns the next KeepConnectedResponse as a zero-copy view.
func (s *ClientKeepConnectedStream) Recv() (masterwire.KeepConnectedResponse, error) {
	b, err := s.s.Recv()
	if err != nil {
		return masterwire.KeepConnectedResponse{}, err
	}
	return masterwire.WrapKeepConnectedResponse(b)
}

// RecvBytes returns the next response frame as the raw ZAP buffer.
func (s *ClientKeepConnectedStream) RecvBytes() ([]byte, error) { return s.s.Recv() }

// CloseSend half-closes the request direction. Idempotent.
func (s *ClientKeepConnectedStream) CloseSend() error { return s.s.CloseSend() }

// ClientAssignStream is the caller's half of the bidirectional StreamAssign
// stream: Send ships requests, Recv yields responses.
type ClientAssignStream struct{ s transport.Stream }

// Send ships one AssignRequest frame.
func (s *ClientAssignStream) Send(frame []byte) error { return s.s.Send(frame) }

// Recv returns the next AssignResponse as a zero-copy view.
func (s *ClientAssignStream) Recv() (masterwire.AssignResponse, error) {
	b, err := s.s.Recv()
	if err != nil {
		return masterwire.AssignResponse{}, err
	}
	return masterwire.WrapAssignResponse(b)
}

// RecvBytes returns the next response frame as the raw ZAP buffer.
func (s *ClientAssignStream) RecvBytes() ([]byte, error) { return s.s.Recv() }

// CloseSend half-closes the request direction. Idempotent.
func (s *ClientAssignStream) CloseSend() error { return s.s.CloseSend() }
