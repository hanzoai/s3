// Package rpc holds the grpc-free streaming contracts shared by the generated
// *_pb packages. It is a leaf: it imports nothing, so every *_pb package can
// depend on it without a cycle.
//
// These mirror the client-side halves of grpc's generic stream interfaces with
// the gRPC plumbing cut out: no embedded grpc.ClientStream, no
// opts ...grpc.CallOption. No caller passes a CallOption and none reaches for
// Context/Header/Trailer on a client stream, so the sets are just the methods
// callers actually use.
package rpc

// ServerStream is the client side of a server-streaming RPC: one request, many
// responses read with Recv until io.EOF.
type ServerStream[T any] interface {
	Recv() (*T, error)
}

// BidiStream is the client side of a bidirectional-streaming RPC: many requests
// via Send, many responses via Recv, CloseSend to half-close the send side.
type BidiStream[Req, Resp any] interface {
	Send(*Req) error
	Recv() (*Resp, error)
	CloseSend() error
}

// ClientStream is the client side of a client-streaming RPC: many requests via
// Send, then CloseAndRecv to half-close and read the single response.
type ClientStream[Req, Resp any] interface {
	Send(*Req) error
	CloseAndRecv() (*Resp, error)
}
