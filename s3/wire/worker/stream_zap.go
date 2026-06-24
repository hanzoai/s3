// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// Native ZAP bidirectional streaming for the WorkerService.WorkerStream RPC.
//
// worker.proto declares the only WorkerService method as a full-duplex stream:
//
//	rpc WorkerStream(stream WorkerMessage) returns (stream AdminMessage)
//
// The unary Channel/Client/Handler/Dispatch in worker_zap.go is a single-shot
// fallback (one WorkerMessage -> one AdminMessage); it cannot carry the duplex
// loop. This file IS the duplex loop, wired over the canonical
// github.com/zap-proto/go/transport stream primitive (Stream + OpenStream +
// ListenStream), which the runtime ships as of v1.5.0. zapgen does not emit
// streaming, so — exactly as the migration recipe prescribes — the stream is
// hand-wired here alongside the generated unary surface, reusing the generated
// ordinal so the wire id stays stable.
//
// Payloads stay zero-copy: every frame is a NewWorkerMessage / NewAdminMessage
// buffer, read with WrapWorkerMessage / WrapAdminMessage. The bytes ARE the
// message — no struct, no marshal, no protobuf. The worker (client) Sends
// WorkerMessage and Recvs AdminMessage; the admin (server) Recvs WorkerMessage
// and Sends AdminMessage.
//
// The server backend is a Go interface (WorkerStreamServer) — the real worker
// subsystem (s3/admin/dash) implements it; this package never imports it. That
// is the same decoupling ObjectStore gives the object service.

package workerwire

import (
	"github.com/zap-proto/go/transport"
)

// streamChannel is the subset of *transport.Conn this stream needs to open a
// client stream. *transport.Conn satisfies it; tests may substitute a fake.
type streamChannel interface {
	OpenStream(method uint32, init []byte) (*transport.Stream, error)
}

// ---------------------------------------------------------------------------
// Client side — the worker dials admin and drives the WorkerMessage half.
// ---------------------------------------------------------------------------

// WorkerStream is the worker-side view of an open WorkerService.WorkerStream:
// it Sends WorkerMessage frames upstream and Recvs AdminMessage frames back.
// Safe for one concurrent Send and one concurrent Recv (standard streaming-RPC
// usage). Payloads are opaque ZAP buffers — build a WorkerMessage to Send with
// NewWorkerMessage, and Wrap a received frame with WrapAdminMessage.
type WorkerStream struct {
	s *transport.Stream
}

// OpenWorkerStream opens the bidi WorkerStream over ch (a *transport.Conn). The
// optional init frame is the first WorkerMessage (build it with
// NewWorkerMessage); pass nil to open with no initial message. The peer's
// stream handler is invoked with the WorkerStream ordinal and init.
func OpenWorkerStream(ch streamChannel, init []byte) (*WorkerStream, error) {
	s, err := ch.OpenStream(WorkerServiceWorkerStreamOrdinal, init)
	if err != nil {
		return nil, err
	}
	return &WorkerStream{s: s}, nil
}

// Send ships one WorkerMessage frame (a NewWorkerMessage buffer) upstream.
func (w *WorkerStream) Send(workerMessage []byte) error { return w.s.Send(workerMessage) }

// Recv returns the next AdminMessage frame as raw bytes (Wrap it with
// WrapAdminMessage), io.EOF once admin half-closes, or transport.ErrClosed if
// the connection drops.
func (w *WorkerStream) Recv() ([]byte, error) { return w.s.Recv() }

// CloseSend half-closes the worker's outbound half; admin's Recv then returns
// io.EOF. Idempotent.
func (w *WorkerStream) CloseSend() error { return w.s.CloseSend() }

// ---------------------------------------------------------------------------
// Server side — admin accepts the stream and drives the AdminMessage half.
// ---------------------------------------------------------------------------

// AdminStream is the admin-side view of an accepted WorkerService.WorkerStream:
// it Recvs WorkerMessage frames from the worker and Sends AdminMessage frames
// back. Init holds the opening WorkerMessage payload (may be empty). Payloads
// are opaque ZAP buffers — Wrap a received frame with WrapWorkerMessage, build
// a reply with NewAdminMessage.
type AdminStream struct {
	init []byte
	s    *transport.Stream
}

// Init returns the opening WorkerMessage frame the worker sent with OpenStream
// (Wrap it with WrapWorkerMessage). It is nil/empty when the stream was opened
// without an initial message.
func (a *AdminStream) Init() []byte { return a.init }

// Recv returns the next WorkerMessage frame as raw bytes (Wrap it with
// WrapWorkerMessage), io.EOF once the worker half-closes, or
// transport.ErrClosed if the connection drops.
func (a *AdminStream) Recv() ([]byte, error) { return a.s.Recv() }

// Send ships one AdminMessage frame (a NewAdminMessage buffer) back to the
// worker.
func (a *AdminStream) Send(adminMessage []byte) error { return a.s.Send(adminMessage) }

// CloseSend half-closes admin's outbound half early; the worker's Recv then
// returns io.EOF. Optional — returning from ServeWorkerStream half-closes it
// automatically. Idempotent.
func (a *AdminStream) CloseSend() error { return a.s.CloseSend() }

// WorkerStreamServer is the admin-side backend contract for the bidi
// WorkerStream RPC. The real worker dashboard (s3/admin/dash) implements it;
// this package never imports that engine, exactly as ObjectStore decouples the
// object service. ServeWorkerStream runs on its own goroutine per accepted
// stream; returning half-closes the admin->worker direction.
type WorkerStreamServer interface {
	ServeWorkerStream(stream *AdminStream)
}

// WorkerStreamHandler adapts a WorkerStreamServer to a transport.StreamHandler,
// to hand to transport.ListenStream as the stream dispatcher. It routes the
// WorkerStream ordinal to srv.ServeWorkerStream and ignores any other ordinal
// (the service declares no other streaming method).
func WorkerStreamHandler(srv WorkerStreamServer) transport.StreamHandler {
	return func(method uint32, init []byte, s *transport.Stream) {
		switch method {
		case WorkerServiceWorkerStreamOrdinal:
			srv.ServeWorkerStream(&AdminStream{init: init, s: s})
		default:
			// Unknown stream ordinal: returning half-closes and ends it.
		}
	}
}

// ListenWorkerStream binds network/addr and serves the bidi WorkerStream on
// every accepted connection, routing each stream to srv. unaryDispatch serves
// the (single-shot) unary fallback and may be nil; pass
// DispatchWorkerService-bound dispatch to keep both surfaces live on one port.
// Returns the running server; Close stops it.
func ListenWorkerStream(network, addr string, srv WorkerStreamServer, unaryDispatch transport.Dispatch) (*transport.Server, error) {
	return transport.ListenStream(network, addr, unaryDispatch, WorkerStreamHandler(srv))
}
