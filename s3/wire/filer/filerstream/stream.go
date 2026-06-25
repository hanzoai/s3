// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// Package filerstream is the native ZAP streaming adapter for the 6 streaming
// HanzoFiler RPCs. It is to the generated s3/wire/filer (filerwire) package what
// s3/wire/plugin/zapstream is to s3/wire/plugin, and what s3/zapsvc is to
// s3/wire/object: the ONLY hand-written glue binding the zero-copy wire messages
// to a real transport. No gRPC, no protobuf framing, no struct marshaling —
// every frame that crosses the socket is a filerwire New*/Wrap* buffer framed as
// a stream message.
//
// zapgen does NOT emit streaming (its method grammar has no `stream` keyword),
// so the generated filerwire service exposes only single-frame unary stubs for
// the streaming ordinals. This adapter wires the REAL duplex primitive shipped
// by github.com/zap-proto/go/transport v1.5.0+ (Conn.OpenStream on the client,
// transport.ListenStream + a StreamHandler on the server, Stream.Send/Recv on
// both ends). The reserved filerwire ordinals are reused as each stream's
// open-method id, so the wire numbering stays stable and the unary stub and this
// streaming path agree on every method number.
//
// The 6 streaming RPCs and their shapes:
//
//	ListEntries            server-stream  (ListEntriesRequest        -> ListEntriesResponse*)
//	StreamRenameEntry      server-stream  (StreamRenameEntryRequest  -> StreamRenameEntryResponse*)
//	TraverseBfsMetadata    server-stream  (TraverseBfsMetadataRequest-> TraverseBfsMetadataResponse*)
//	SubscribeMetadata      server-stream  (SubscribeMetadataRequest  -> SubscribeMetadataResponse*)
//	SubscribeLocalMetadata server-stream  (SubscribeMetadataRequest  -> SubscribeMetadataResponse*)
//	StreamMutateEntry      bidirectional  (StreamMutateEntryRequest* <-> StreamMutateEntryResponse*)
//
// SubscribeMetadata and SubscribeLocalMetadata share the request/response
// message types but are distinct RPCs (distinct ordinals); they get distinct
// backend methods so the server can scope "local" vs cluster-wide.
//
// The backend (Server below) is defined in filer terms and does NOT import the
// real filer/master/volume engine — the integrator implements it against the
// live engine; tests use an in-memory stub (see stream_roundtrip_test.go).

package filerstream

import (
	"context"
	"io"

	filerwire "github.com/hanzoai/s3/s3/wire/filer"
	"github.com/zap-proto/go/transport"
)

// Ordinals each streaming RPC opens with — the SAME ids the generated unary
// stubs reserve, so the streaming open frame and any unary probe of the method
// agree on the wire number.
const (
	ListEntriesOrdinal            = filerwire.HanzoFilerListEntriesOrdinal
	StreamRenameEntryOrdinal      = filerwire.HanzoFilerStreamRenameEntryOrdinal
	StreamMutateEntryOrdinal      = filerwire.HanzoFilerStreamMutateEntryOrdinal
	TraverseBfsMetadataOrdinal    = filerwire.HanzoFilerTraverseBfsMetadataOrdinal
	SubscribeMetadataOrdinal      = filerwire.HanzoFilerSubscribeMetadataOrdinal
	SubscribeLocalMetadataOrdinal = filerwire.HanzoFilerSubscribeLocalMetadataOrdinal
)

// IsEOF reports whether err is the stream's half-close sentinel, so callers can
// end a Recv loop without importing io.
func IsEOF(err error) bool { return err == io.EOF }

// ---------------------------------------------------------------------------
// Server side
// ---------------------------------------------------------------------------

// Server is the server-side backend the streaming adapter delegates to: the
// real filer engine implements it; tests use an in-memory stub. Each method
// handles the full lifetime of one stream of its RPC.
//
// For the 5 server-streaming RPCs, init is the opening request (already wrapped)
// and the stream's Send ships successive response frames; returning ends the
// stream (the client's Recv then sees io.EOF). For the one bidirectional RPC
// (StreamMutateEntry) both halves interleave: Recv yields request frames, Send
// ships response frames.
//
// A returned error just ends the stream (the transport half-closes the send
// side); per-frame errors are carried in the response messages' own error
// fields (e.g. StreamMutateEntryResponse.Error / .Errno), not at the transport.
type Server interface {
	// ListEntries streams the directory listing for init.
	ListEntries(init filerwire.ListEntriesRequest, s *ListEntriesStream) error
	// StreamRenameEntry streams the rename event notifications for init.
	StreamRenameEntry(init filerwire.StreamRenameEntryRequest, s *StreamRenameEntryStream) error
	// TraverseBfsMetadata streams a breadth-first metadata walk from init.
	TraverseBfsMetadata(init filerwire.TraverseBfsMetadataRequest, s *TraverseBfsMetadataStream) error
	// SubscribeMetadata streams cluster-wide metadata change events from init.since_ns.
	SubscribeMetadata(init filerwire.SubscribeMetadataRequest, s *SubscribeMetadataStream) error
	// SubscribeLocalMetadata streams this filer's local metadata change events.
	SubscribeLocalMetadata(init filerwire.SubscribeMetadataRequest, s *SubscribeMetadataStream) error
	// StreamMutateEntry serves one mount's ordered bidirectional mutation stream.
	// init is the first request frame; s.Recv yields the rest, s.Send replies.
	StreamMutateEntry(init filerwire.StreamMutateEntryRequest, s *StreamMutateEntryStream) error
}

// ListEntriesStream is the server's half of a ListEntries server-stream: Send
// ships one ListEntriesResponse (a filerwire NewListEntriesResponse buffer).
type ListEntriesStream struct{ s transport.Stream }

// Context is cancelled when this stream ends or the connection drops.
func (x *ListEntriesStream) Context() context.Context { return x.s.Context() }

// Send ships one ListEntriesResponse frame (the bytes ARE the message).
func (s *ListEntriesStream) Send(frame []byte) error { return s.s.Send(frame) }

// CloseSend half-closes the stream (idempotent; the handler returning does this).
func (s *ListEntriesStream) CloseSend() error { return s.s.CloseSend() }

// StreamRenameEntryStream is the server's half of a StreamRenameEntry stream.
type StreamRenameEntryStream struct{ s transport.Stream }

// Context is cancelled when this stream ends or the connection drops.
func (x *StreamRenameEntryStream) Context() context.Context { return x.s.Context() }

// Send ships one StreamRenameEntryResponse frame.
func (s *StreamRenameEntryStream) Send(frame []byte) error { return s.s.Send(frame) }

// CloseSend half-closes the stream.
func (s *StreamRenameEntryStream) CloseSend() error { return s.s.CloseSend() }

// TraverseBfsMetadataStream is the server's half of a TraverseBfsMetadata stream.
type TraverseBfsMetadataStream struct{ s transport.Stream }

// Context is cancelled when this stream ends or the connection drops.
func (x *TraverseBfsMetadataStream) Context() context.Context { return x.s.Context() }

// Send ships one TraverseBfsMetadataResponse frame.
func (s *TraverseBfsMetadataStream) Send(frame []byte) error { return s.s.Send(frame) }

// CloseSend half-closes the stream.
func (s *TraverseBfsMetadataStream) CloseSend() error { return s.s.CloseSend() }

// SubscribeMetadataStream is the server's half of a (local or cluster-wide)
// SubscribeMetadata stream.
type SubscribeMetadataStream struct{ s transport.Stream }

// Context is cancelled when this stream ends or the connection drops.
func (x *SubscribeMetadataStream) Context() context.Context { return x.s.Context() }

// Send ships one SubscribeMetadataResponse frame.
func (s *SubscribeMetadataStream) Send(frame []byte) error { return s.s.Send(frame) }

// CloseSend half-closes the stream.
func (s *SubscribeMetadataStream) CloseSend() error { return s.s.CloseSend() }

// StreamMutateEntryStream is the server's half of the bidirectional
// StreamMutateEntry stream: Recv yields successive StreamMutateEntryRequest
// frames (request_id-correlated), Send ships StreamMutateEntryResponse frames.
type StreamMutateEntryStream struct{ s transport.Stream }

// Context is cancelled when this stream ends or the connection drops.
func (x *StreamMutateEntryStream) Context() context.Context { return x.s.Context() }

// Recv returns the next request frame as a zero-copy view, or io.EOF once the
// client half-closes its send side.
func (s *StreamMutateEntryStream) Recv() (filerwire.StreamMutateEntryRequest, error) {
	b, err := s.s.Recv()
	if err != nil {
		return filerwire.StreamMutateEntryRequest{}, err
	}
	return filerwire.WrapStreamMutateEntryRequest(b)
}

// RecvBytes returns the next request frame as the raw ZAP buffer.
func (s *StreamMutateEntryStream) RecvBytes() ([]byte, error) { return s.s.Recv() }

// Send ships one StreamMutateEntryResponse frame.
func (s *StreamMutateEntryStream) Send(frame []byte) error { return s.s.Send(frame) }

// CloseSend half-closes the response direction.
func (s *StreamMutateEntryStream) CloseSend() error { return s.s.CloseSend() }

// Handler builds the transport StreamHandler that routes an inbound stream-open
// to the matching Server method by ordinal. Pass it to transport.ListenStream
// (or use Serve, which does that). An open for an ordinal that is not one of the
// 6 streaming RPCs is ignored (the stream half-closes immediately). A malformed
// init buffer ends the stream without invoking the backend.
func Handler(srv Server) transport.StreamHandler {
	return func(method uint32, init []byte, s transport.Stream) {
		switch method {
		case ListEntriesOrdinal:
			v, err := filerwire.WrapListEntriesRequest(init)
			if err != nil {
				return
			}
			_ = srv.ListEntries(v, &ListEntriesStream{s: s})
		case StreamRenameEntryOrdinal:
			v, err := filerwire.WrapStreamRenameEntryRequest(init)
			if err != nil {
				return
			}
			_ = srv.StreamRenameEntry(v, &StreamRenameEntryStream{s: s})
		case TraverseBfsMetadataOrdinal:
			v, err := filerwire.WrapTraverseBfsMetadataRequest(init)
			if err != nil {
				return
			}
			_ = srv.TraverseBfsMetadata(v, &TraverseBfsMetadataStream{s: s})
		case SubscribeMetadataOrdinal:
			v, err := filerwire.WrapSubscribeMetadataRequest(init)
			if err != nil {
				return
			}
			_ = srv.SubscribeMetadata(v, &SubscribeMetadataStream{s: s})
		case SubscribeLocalMetadataOrdinal:
			v, err := filerwire.WrapSubscribeMetadataRequest(init)
			if err != nil {
				return
			}
			_ = srv.SubscribeLocalMetadata(v, &SubscribeMetadataStream{s: s})
		case StreamMutateEntryOrdinal:
			v, err := filerwire.WrapStreamMutateEntryRequest(init)
			if err != nil {
				return
			}
			_ = srv.StreamMutateEntry(v, &StreamMutateEntryStream{s: s})
		default:
			return // unknown stream method: returning half-closes + ends it
		}
	}
}

// StreamServer is the running native ZAP HanzoFiler streaming endpoint. Close
// stops it and tears down live streams.
type StreamServer struct{ srv *transport.Server }

// Serve starts the native ZAP HanzoFiler streaming endpoint on network/addr
// (e.g. "tcp", ":18889"), backed by srv. It registers NO unary dispatch (the
// unary RPCs are served by filerwire.Serve); inbound stream-opens route to srv
// via Handler.
//
// To serve BOTH the unary and streaming RPCs on ONE listener, call
// transport.ListenStream(network, addr, filerwire.Dispatch(backend), Handler(srv))
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

// Client is the caller's typed connection to one filer streaming endpoint. It
// owns the transport connection; it may open multiple concurrent streams (one
// per outstanding streaming RPC) over the single connection.
type Client struct{ conn transport.Conn }

// Dial connects to the filer streaming endpoint at addr (e.g.
// "filer.hanzo.svc:18889") over plain TCP. For the PQ-secured mesh use
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
// sufficient. This also lets a caller open streams over the SAME connection a
// filerwire.Client uses for unary calls (pass filerwire.Client.Conn()).
func NewClient(conn transport.Conn) *Client { return &Client{conn: conn} }

// Close releases the underlying connection.
func (c *Client) Close() error { return c.conn.Close() }

// ListEntries opens the ListEntries server-stream with req (a filerwire
// NewListEntriesRequest buffer) as the opening frame. Recv the returned
// ClientListEntriesStream until io.EOF.
func (c *Client) ListEntries(req []byte) (*ClientListEntriesStream, error) {
	s, err := c.conn.OpenStream(ListEntriesOrdinal, req)
	if err != nil {
		return nil, err
	}
	return &ClientListEntriesStream{s: s}, nil
}

// StreamRenameEntry opens the StreamRenameEntry server-stream.
func (c *Client) StreamRenameEntry(req []byte) (*ClientStreamRenameEntryStream, error) {
	s, err := c.conn.OpenStream(StreamRenameEntryOrdinal, req)
	if err != nil {
		return nil, err
	}
	return &ClientStreamRenameEntryStream{s: s}, nil
}

// TraverseBfsMetadata opens the TraverseBfsMetadata server-stream.
func (c *Client) TraverseBfsMetadata(req []byte) (*ClientTraverseBfsMetadataStream, error) {
	s, err := c.conn.OpenStream(TraverseBfsMetadataOrdinal, req)
	if err != nil {
		return nil, err
	}
	return &ClientTraverseBfsMetadataStream{s: s}, nil
}

// SubscribeMetadata opens the cluster-wide SubscribeMetadata server-stream.
func (c *Client) SubscribeMetadata(req []byte) (*ClientSubscribeMetadataStream, error) {
	s, err := c.conn.OpenStream(SubscribeMetadataOrdinal, req)
	if err != nil {
		return nil, err
	}
	return &ClientSubscribeMetadataStream{s: s}, nil
}

// SubscribeLocalMetadata opens the local SubscribeMetadata server-stream.
func (c *Client) SubscribeLocalMetadata(req []byte) (*ClientSubscribeMetadataStream, error) {
	s, err := c.conn.OpenStream(SubscribeLocalMetadataOrdinal, req)
	if err != nil {
		return nil, err
	}
	return &ClientSubscribeMetadataStream{s: s}, nil
}

// StreamMutateEntry opens the bidirectional StreamMutateEntry stream, sending
// first (a filerwire NewStreamMutateEntryRequest buffer) as the opening frame.
// The returned stream Sends further request frames and Recvs response frames.
func (c *Client) StreamMutateEntry(first []byte) (*ClientStreamMutateEntryStream, error) {
	s, err := c.conn.OpenStream(StreamMutateEntryOrdinal, first)
	if err != nil {
		return nil, err
	}
	return &ClientStreamMutateEntryStream{s: s}, nil
}

// ClientListEntriesStream is the caller's half of a ListEntries server-stream.
type ClientListEntriesStream struct{ s transport.Stream }

// Recv returns the next ListEntriesResponse as a zero-copy view, or io.EOF once
// the server ends the stream.
func (s *ClientListEntriesStream) Recv() (filerwire.ListEntriesResponse, error) {
	b, err := s.s.Recv()
	if err != nil {
		return filerwire.ListEntriesResponse{}, err
	}
	return filerwire.WrapListEntriesResponse(b)
}

// RecvBytes returns the next response frame as the raw ZAP buffer.
func (s *ClientListEntriesStream) RecvBytes() ([]byte, error) { return s.s.Recv() }

// Close releases the caller's send half (server-stream: the caller never sends
// after the open, so this just half-closes).
func (s *ClientListEntriesStream) Close() error { return s.s.CloseSend() }

// ClientStreamRenameEntryStream is the caller's half of a StreamRenameEntry stream.
type ClientStreamRenameEntryStream struct{ s transport.Stream }

// Recv returns the next StreamRenameEntryResponse as a zero-copy view.
func (s *ClientStreamRenameEntryStream) Recv() (filerwire.StreamRenameEntryResponse, error) {
	b, err := s.s.Recv()
	if err != nil {
		return filerwire.StreamRenameEntryResponse{}, err
	}
	return filerwire.WrapStreamRenameEntryResponse(b)
}

// RecvBytes returns the next response frame as the raw ZAP buffer.
func (s *ClientStreamRenameEntryStream) RecvBytes() ([]byte, error) { return s.s.Recv() }

// Close half-closes the caller's send half.
func (s *ClientStreamRenameEntryStream) Close() error { return s.s.CloseSend() }

// ClientTraverseBfsMetadataStream is the caller's half of a TraverseBfsMetadata stream.
type ClientTraverseBfsMetadataStream struct{ s transport.Stream }

// Recv returns the next TraverseBfsMetadataResponse as a zero-copy view.
func (s *ClientTraverseBfsMetadataStream) Recv() (filerwire.TraverseBfsMetadataResponse, error) {
	b, err := s.s.Recv()
	if err != nil {
		return filerwire.TraverseBfsMetadataResponse{}, err
	}
	return filerwire.WrapTraverseBfsMetadataResponse(b)
}

// RecvBytes returns the next response frame as the raw ZAP buffer.
func (s *ClientTraverseBfsMetadataStream) RecvBytes() ([]byte, error) { return s.s.Recv() }

// Close half-closes the caller's send half.
func (s *ClientTraverseBfsMetadataStream) Close() error { return s.s.CloseSend() }

// ClientSubscribeMetadataStream is the caller's half of a (local or cluster-wide)
// SubscribeMetadata server-stream.
type ClientSubscribeMetadataStream struct{ s transport.Stream }

// Recv returns the next SubscribeMetadataResponse as a zero-copy view.
func (s *ClientSubscribeMetadataStream) Recv() (filerwire.SubscribeMetadataResponse, error) {
	b, err := s.s.Recv()
	if err != nil {
		return filerwire.SubscribeMetadataResponse{}, err
	}
	return filerwire.WrapSubscribeMetadataResponse(b)
}

// RecvBytes returns the next response frame as the raw ZAP buffer.
func (s *ClientSubscribeMetadataStream) RecvBytes() ([]byte, error) { return s.s.Recv() }

// Close half-closes the caller's send half (ending the subscription).
func (s *ClientSubscribeMetadataStream) Close() error { return s.s.CloseSend() }

// ClientStreamMutateEntryStream is the caller's half of the bidirectional
// StreamMutateEntry stream: Send ships request frames, Recv yields responses.
type ClientStreamMutateEntryStream struct{ s transport.Stream }

// Send ships one StreamMutateEntryRequest frame (a filerwire
// NewStreamMutateEntryRequest buffer).
func (s *ClientStreamMutateEntryStream) Send(frame []byte) error { return s.s.Send(frame) }

// Recv returns the next StreamMutateEntryResponse as a zero-copy view, or io.EOF
// once the server ends the stream.
func (s *ClientStreamMutateEntryStream) Recv() (filerwire.StreamMutateEntryResponse, error) {
	b, err := s.s.Recv()
	if err != nil {
		return filerwire.StreamMutateEntryResponse{}, err
	}
	return filerwire.WrapStreamMutateEntryResponse(b)
}

// RecvBytes returns the next response frame as the raw ZAP buffer.
func (s *ClientStreamMutateEntryStream) RecvBytes() ([]byte, error) { return s.s.Recv() }

// CloseSend half-closes the request direction (no more requests); the server's
// Recv then returns io.EOF. Idempotent.
func (s *ClientStreamMutateEntryStream) CloseSend() error { return s.s.CloseSend() }
