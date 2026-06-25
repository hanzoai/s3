// Code generated from filer.proto; DO NOT EDIT.

package filerwire

import (
	"github.com/zap-proto/go/transport"
)

// Streaming surface for the HanzoFiler service.
//
// The unary methods on HanzoFilerClient ride a one-shot Call channel
// (HanzoFilerChannel). The six streaming methods cannot: a server-streaming
// RPC answers one request with many response frames, and StreamMutateEntry is
// fully bidirectional. Those ride transport.Conn's stream layer
// (Conn.OpenStream on the client, transport.StreamHandler on the server), one
// stream per RPC, correlated by the open frame's PromiseID.
//
// Every frame on a stream is an opaque filer wire buffer of this package:
// build it with the matching New* and read it with the matching Wrap*. The
// stream itself moves bytes; the bytes ARE the message (zero-copy, no marshal).
//
// Shapes (filer.proto):
//   - ListEntries, StreamRenameEntry, TraverseBfsMetadata, SubscribeMetadata,
//     SubscribeLocalMetadata: server-streaming — client sends one request as
//     the open payload, then reads response frames until io.EOF.
//   - StreamMutateEntry: bidirectional — client and server both Send and Recv
//     until each half-closes.

// HanzoFilerStreamClient issues the HanzoFiler streaming RPCs over a
// transport.Conn. Unlike the unary HanzoFilerClient (which needs only a Call
// channel), streaming needs the full duplex Conn to open per-RPC streams.
type HanzoFilerStreamClient struct {
	conn transport.Conn
}

// NewHanzoFilerStreamClient returns a streaming client bound to conn. The same
// conn may also back a unary HanzoFilerClient (it satisfies HanzoFilerChannel),
// so one connection serves both call styles.
func NewHanzoFilerStreamClient(conn transport.Conn) *HanzoFilerStreamClient {
	return &HanzoFilerStreamClient{conn: conn}
}

// openServerStream opens a server-streaming RPC: it ships req as the open
// payload and half-closes the send side immediately (the client sends nothing
// further), leaving a stream the caller drains with Recv until io.EOF.
func (c *HanzoFilerStreamClient) openServerStream(ordinal uint32, req []byte) (transport.Stream, error) {
	s, err := c.conn.OpenStream(ordinal, req)
	if err != nil {
		return nil, err
	}
	if err := s.CloseSend(); err != nil {
		return nil, err
	}
	return s, nil
}

// ListEntries opens the server-streaming ListEntries RPC. Drain the returned
// stream with Recv: each frame is a ListEntriesResponse buffer (WrapListEntriesResponse),
// terminated by io.EOF.
func (c *HanzoFilerStreamClient) ListEntries(req []byte) (transport.Stream, error) {
	return c.openServerStream(HanzoFilerListEntriesOrdinal, req)
}

// StreamRenameEntry opens the server-streaming StreamRenameEntry RPC. Each
// response frame is a StreamRenameEntryResponse buffer.
func (c *HanzoFilerStreamClient) StreamRenameEntry(req []byte) (transport.Stream, error) {
	return c.openServerStream(HanzoFilerStreamRenameEntryOrdinal, req)
}

// TraverseBfsMetadata opens the server-streaming TraverseBfsMetadata RPC. Each
// response frame is a TraverseBfsMetadataResponse buffer.
func (c *HanzoFilerStreamClient) TraverseBfsMetadata(req []byte) (transport.Stream, error) {
	return c.openServerStream(HanzoFilerTraverseBfsMetadataOrdinal, req)
}

// SubscribeMetadata opens the server-streaming SubscribeMetadata RPC. Each
// response frame is a SubscribeMetadataResponse buffer.
func (c *HanzoFilerStreamClient) SubscribeMetadata(req []byte) (transport.Stream, error) {
	return c.openServerStream(HanzoFilerSubscribeMetadataOrdinal, req)
}

// SubscribeLocalMetadata opens the server-streaming SubscribeLocalMetadata RPC
// (it shares SubscribeMetadata's request/response message types). Each response
// frame is a SubscribeMetadataResponse buffer.
func (c *HanzoFilerStreamClient) SubscribeLocalMetadata(req []byte) (transport.Stream, error) {
	return c.openServerStream(HanzoFilerSubscribeLocalMetadataOrdinal, req)
}

// StreamMutateEntry opens the bidirectional StreamMutateEntry RPC. The caller
// Sends StreamMutateEntryRequest buffers and Recvs StreamMutateEntryResponse
// buffers concurrently, half-closing the send side with CloseSend when done.
// The open carries no initial payload; the first request is the first Send.
func (c *HanzoFilerStreamClient) StreamMutateEntry() (transport.Stream, error) {
	return c.conn.OpenStream(HanzoFilerStreamMutateEntryOrdinal, nil)
}

// HanzoFilerStreamHandler is the server contract for the HanzoFiler streaming
// RPCs. Implement each method and route inbound streams to it with
// DispatchHanzoFilerStream.
//
// For the five server-streaming methods, init is the opening request buffer and
// the handler pushes response frames with s.Send; returning closes the stream
// (the client's Recv sees io.EOF). For the bidirectional StreamMutateEntry, the
// handler Recvs request frames and Sends response frames on s until the client
// half-closes (Recv returns io.EOF) and the handler returns. Each buffer is an
// opaque filer wire message of this package (build with New*, read with Wrap*).
type HanzoFilerStreamHandler interface {
	ListEntries(req []byte, s transport.Stream) error              // STREAMING (server)
	StreamRenameEntry(req []byte, s transport.Stream) error        // STREAMING (server)
	TraverseBfsMetadata(req []byte, s transport.Stream) error      // STREAMING (server)
	SubscribeMetadata(req []byte, s transport.Stream) error        // STREAMING (server)
	SubscribeLocalMetadata(req []byte, s transport.Stream) error   // STREAMING (server)
	StreamMutateEntry(s transport.Stream) error                    // STREAMING (bidirectional)
}

// DispatchHanzoFilerStream returns a transport.StreamHandler that routes one
// inbound stream by its open method ordinal to h. It is the streaming sibling
// of DispatchHanzoFiler: pass it as the stream argument of
// transport.ListenStream alongside func(env []byte) ([]byte, error) {
// return filerwire.DispatchHanzoFiler(unaryHandler, env) }.
//
// An unknown ordinal returns without serving the stream, so the opener's Recv
// observes the half-close (io.EOF) immediately. A handler error half-closes the
// send side the same way (the streaming RPCs report terminal status in-band via
// their response messages, mirroring the gRPC server methods they back).
func DispatchHanzoFilerStream(h HanzoFilerStreamHandler) transport.StreamHandler {
	return func(method uint32, init []byte, s transport.Stream) {
		switch method {
		case HanzoFilerListEntriesOrdinal:
			_ = h.ListEntries(init, s)
		case HanzoFilerStreamRenameEntryOrdinal:
			_ = h.StreamRenameEntry(init, s)
		case HanzoFilerTraverseBfsMetadataOrdinal:
			_ = h.TraverseBfsMetadata(init, s)
		case HanzoFilerSubscribeMetadataOrdinal:
			_ = h.SubscribeMetadata(init, s)
		case HanzoFilerSubscribeLocalMetadataOrdinal:
			_ = h.SubscribeLocalMetadata(init, s)
		case HanzoFilerStreamMutateEntryOrdinal:
			_ = h.StreamMutateEntry(s)
		}
	}
}
