// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// Package pluginstream is the native ZAP streaming adapter for the
// PluginControlService.WorkerStream bidirectional RPC — the duplex control
// plane between an admin (server) and an external worker (client). It is to the
// generated s3/wire/plugin (pluginwire) package what s3/object is to
// s3/wire/object: the ONLY hand-written glue binding the zero-copy wire
// messages to a real transport. No gRPC, no protobuf framing, no struct
// marshaling — every frame that crosses the socket is a pluginwire
// New*/Wrap* buffer.
//
// WorkerStream is `rpc WorkerStream(stream WorkerToAdminMessage) returns
// (stream AdminToWorkerMessage)` — true bidirectional streaming. zapgen does
// NOT emit streaming (its method grammar has no `stream` keyword), and the
// generated pluginwire service exposes only a single-frame unary stub for the
// ordinal. This adapter wires the REAL duplex primitive shipped by
// github.com/zap-proto/go/transport v1.5.0+ (Conn.OpenStream on the client,
// transport.ListenStream + a StreamHandler on the server, Stream.Send/Recv on
// both ends). The reserved ordinal pluginwire.PluginControlServiceWorkerStreamOrdinal
// is reused as the stream's open-method id so the wire stays stable and the
// generated unary stub and this streaming path agree on the method number.
//
// The worker (client) opens the stream and is the first to speak: its opening
// WorkerToAdminMessage (conventionally a WorkerHello) is the OpenStream init
// payload. Thereafter both halves interleave — the worker Sends
// WorkerToAdminMessage frames and Recvs AdminToWorkerMessage frames; the admin
// does the mirror image.
//
// The backend is the Admin interface below — defined in worker/admin terms,
// NOT importing the real admin control plane (filer/master/maintenance). The
// integrator implements Admin against the real control plane and starts Serve;
// tests use an in-memory stub (see roundtrip_test.go).
package pluginstream

import (
	"io"

	pluginwire "github.com/hanzoai/s3/s3/wire/plugin"
	"github.com/zap-proto/go/transport"
)

// WorkerStreamOrdinal is the rpc method id used to open a WorkerStream. It is
// the SAME ordinal the generated unary stub reserves, so the streaming open
// frame and any future unary probe of this method agree on the wire number.
const WorkerStreamOrdinal = pluginwire.PluginControlServiceWorkerStreamOrdinal

// ---------------------------------------------------------------------------
// Server side
// ---------------------------------------------------------------------------

// Admin is the server-side backend the streaming adapter delegates to: the
// real admin control plane implements it; tests use an in-memory stub. One
// call to ServeWorker handles the full lifetime of a single worker's duplex
// stream. The init buffer is the worker's opening WorkerToAdminMessage (its
// WorkerHello); s is the live stream to Recv further worker frames from and
// Send admin frames on. Returning from ServeWorker half-closes the admin→worker
// direction (the worker's Recv then sees io.EOF) and ends the stream.
type Admin interface {
	ServeWorker(init pluginwire.WorkerToAdminMessage, s *ServerStream) error
}

// ServerStream is the admin's half of one worker's duplex stream. It is a
// thin zero-copy wrapper over the transport stream: Recv yields the next
// WorkerToAdminMessage (worker→admin), Send ships one AdminToWorkerMessage
// (admin→worker). Safe for one concurrent Recv and one concurrent Send — the
// standard streaming-RPC discipline.
type ServerStream struct{ s transport.Stream }

// Recv returns the next worker→admin frame as a zero-copy view, or io.EOF once
// the worker half-closes its send side.
func (s *ServerStream) Recv() (pluginwire.WorkerToAdminMessage, error) {
	b, err := s.s.Recv()
	if err != nil {
		return pluginwire.WorkerToAdminMessage{}, err
	}
	return pluginwire.WrapWorkerToAdminMessage(b)
}

// RecvBytes returns the next worker→admin frame as the raw ZAP buffer (for
// callers that forward bytes without inspecting fields).
func (s *ServerStream) RecvBytes() ([]byte, error) { return s.s.Recv() }

// Send ships one admin→worker frame. frame must be a pluginwire
// NewAdminToWorkerMessage buffer (the bytes ARE the message).
func (s *ServerStream) Send(frame []byte) error { return s.s.Send(frame) }

// CloseSend half-closes the admin→worker direction; the worker's Recv then
// returns io.EOF. Idempotent. (The handler returning does this automatically.)
func (s *ServerStream) CloseSend() error { return s.s.CloseSend() }

// Handler builds the transport StreamHandler that routes an inbound stream-open
// to admin.ServeWorker. Pass it to transport.ListenStream (or use Serve, which
// does that). An open for any ordinal other than WorkerStreamOrdinal is ignored
// (the stream half-closes immediately).
func Handler(admin Admin) transport.StreamHandler {
	return func(method uint32, init []byte, s transport.Stream) {
		if method != WorkerStreamOrdinal {
			return // unknown stream method: returning half-closes + ends it
		}
		hello, err := pluginwire.WrapWorkerToAdminMessage(init)
		if err != nil {
			return
		}
		_ = admin.ServeWorker(hello, &ServerStream{s: s})
	}
}

// Server is the running native ZAP PluginControlService endpoint. Close stops
// it and tears down live worker streams.
type Server struct{ srv *transport.Server }

// Serve starts the native ZAP PluginControlService on network/addr (e.g.
// "tcp", ":18920"), backed by admin. It registers NO unary dispatch (the only
// RPC is the WorkerStream stream); inbound stream-opens route to admin via
// Handler. For the PQ-secured mesh, establish the listener over the transport
// TLS/QUIC path and call transport.ListenStream on it directly (see the
// package doc), then wrap with NewServer-equivalent composition.
func Serve(network, addr string, admin Admin) (*Server, error) {
	srv, err := transport.ListenStream(network, addr, nil, Handler(admin))
	if err != nil {
		return nil, err
	}
	return &Server{srv: srv}, nil
}

// Addr is the listener address (useful with ":0").
func (s *Server) Addr() string { return s.srv.Addr().String() }

// Close stops the server.
func (s *Server) Close() error { return s.srv.Close() }

// ---------------------------------------------------------------------------
// Client side (the worker)
// ---------------------------------------------------------------------------

// Client is the worker's typed connection to one admin endpoint. It owns the
// transport connection; one Client may open one WorkerStream at a time over it
// (the control plane is a single long-lived stream per worker).
type Client struct{ conn transport.Conn }

// Dial connects to the admin PluginControlService at addr (e.g.
// "admin.hanzo.svc:18920") over plain TCP. For the PQ-secured mesh use
// transport.DialTLS and NewClient.
func Dial(network, addr string) (*Client, error) {
	conn, err := transport.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn}, nil
}

// NewClient wraps an already-established transport.Conn (TCP, Unix, or PQ-TLS).
// The conn's read loop routes stream frames, so a plain call-only Dial conn is
// sufficient for the worker side.
func NewClient(conn transport.Conn) *Client { return &Client{conn: conn} }

// Close releases the underlying connection.
func (c *Client) Close() error { return c.conn.Close() }

// OpenWorkerStream opens the duplex control stream, sending hello (a pluginwire
// NewWorkerToAdminMessage buffer, conventionally a WorkerHello) as the opening
// frame. The returned ClientStream Sends further worker→admin frames and Recvs
// admin→worker frames.
func (c *Client) OpenWorkerStream(hello []byte) (*ClientStream, error) {
	s, err := c.conn.OpenStream(WorkerStreamOrdinal, hello)
	if err != nil {
		return nil, err
	}
	return &ClientStream{s: s}, nil
}

// ClientStream is the worker's half of the duplex control stream. Send ships
// one WorkerToAdminMessage (worker→admin); Recv yields the next
// AdminToWorkerMessage (admin→worker). Safe for one concurrent Send and one
// concurrent Recv.
type ClientStream struct{ s transport.Stream }

// Send ships one worker→admin frame. frame must be a pluginwire
// NewWorkerToAdminMessage buffer.
func (s *ClientStream) Send(frame []byte) error { return s.s.Send(frame) }

// Recv returns the next admin→worker frame as a zero-copy view, or io.EOF once
// the admin half-closes (e.g. its ServeWorker returned).
func (s *ClientStream) Recv() (pluginwire.AdminToWorkerMessage, error) {
	b, err := s.s.Recv()
	if err != nil {
		return pluginwire.AdminToWorkerMessage{}, err
	}
	return pluginwire.WrapAdminToWorkerMessage(b)
}

// RecvBytes returns the next admin→worker frame as the raw ZAP buffer.
func (s *ClientStream) RecvBytes() ([]byte, error) { return s.s.Recv() }

// CloseSend half-closes the worker→admin direction (no more worker frames);
// the admin's Recv then returns io.EOF. Idempotent.
func (s *ClientStream) CloseSend() error { return s.s.CloseSend() }

// IsEOF reports whether err is the stream's half-close sentinel, so callers can
// end their Recv loop without importing io.
func IsEOF(err error) bool { return err == io.EOF }
