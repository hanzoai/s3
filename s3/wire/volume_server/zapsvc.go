// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// Native ZAP service adapter for VolumeServer — the hand-written glue that binds
// the generated VolumeServer wire (volume_server_zap.go: ordinals + Client +
// Handler + Dispatch) and the zero-copy message views/builders (the *_zap.go
// message files) to a real volume-server backend over the canonical
// github.com/zap-proto/go transport. This is the volume_server analogue of
// s3/zapsvc/service.go (object) and s3/wire/s3_lifecycle/zapsvc.go (lifecycle):
// it kills gRPC for the volume server's RPCs and replaces them with the SAME
// wire hanzo and lux share. No struct marshaling, no protobuf, no HTTP — the
// bytes ARE the message.
//
// Shape (mirrors s3/wire/mq_agent/service.go, the multi-RPC + streaming
// reference):
//   - Unary RPCs (37 of them — BatchDelete, VolumeMount, VolumeStatus, GetState,
//     the EC mount/copy/delete control calls, Ping, …) mirror the object
//     service's Get/Put: the dispatch Wraps nothing itself — it hands the raw
//     request buffer to the backend, which Wraps it, acts, and returns a freshly
//     built response buffer. Dispatched by ordinal over a plain Call/Response
//     envelope.
//   - Streaming RPCs are hand-wired over transport.OpenStream /
//     transport.StreamHandler / transport.Stream, exactly like the mq_agent and
//     plugin/worker stream adapters, because zapgen cannot express streams. Each
//     stream frame is a zero-copy New*/Wrap* buffer (same doctrine as unary); the
//     opening request rides as the stream's init payload and is replayed as the
//     server's first Recv. NOTHING is faked or buffered into a single round-trip.
//     The proto has 11 streams: 10 server-streaming (VacuumVolumeCompact,
//     VolumeIncrementalCopy, VolumeCopy, CopyFile, ReadAllNeedles,
//     VolumeTailSender, VolumeEcShardRead, VolumeTierMoveDatToRemote,
//     VolumeTierMoveDatFromRemote, Query) and 1 client-streaming (ReceiveFile).
//
// The backend is a Go interface (VolumeServerStore) expressed purely in raw ZAP
// message buffers and typed stream handles — it imports neither the
// volume_server_pb domain model nor the real filer/volume/master engine. The S3
// volume server wires its concrete implementation (the methods today registered
// on the gRPC VolumeServer) to this interface at integration time; tests use an
// in-memory fake (see zapsvc_roundtrip_test.go).

package volume_serverwire

import (
	"github.com/zap-proto/go/rpc"
	"github.com/zap-proto/go/transport"
)

// VolumeServerStore is the backend the ZAP volume-server service delegates to.
// It is expressed in the raw New*/Wrap* doctrine: every unary method takes one
// request buffer and returns one response buffer (decode with Wrap<Req>, build
// with New<Resp>); every streaming method takes a typed stream the handler
// drives until the peer half-closes — a streaming method NEVER returns a single
// buffer (that would collapse the stream).
//
// The real volume server (the type that today implements the gRPC
// VolumeServerServer) satisfies this; tests use an in-memory fake. Request views
// alias the request buffer, so copy out any bytes an implementation retains past
// the call.
type VolumeServerStore interface {
	// --- unary (37) ---
	BatchDelete(req []byte) (resp []byte, err error)
	VacuumVolumeCheck(req []byte) (resp []byte, err error)
	VacuumVolumeCommit(req []byte) (resp []byte, err error)
	VacuumVolumeCleanup(req []byte) (resp []byte, err error)
	DeleteCollection(req []byte) (resp []byte, err error)
	AllocateVolume(req []byte) (resp []byte, err error)
	VolumeSyncStatus(req []byte) (resp []byte, err error)
	VolumeMount(req []byte) (resp []byte, err error)
	VolumeUnmount(req []byte) (resp []byte, err error)
	VolumeDelete(req []byte) (resp []byte, err error)
	VolumeMarkReadonly(req []byte) (resp []byte, err error)
	VolumeMarkWritable(req []byte) (resp []byte, err error)
	VolumeConfigure(req []byte) (resp []byte, err error)
	VolumeStatus(req []byte) (resp []byte, err error)
	GetState(req []byte) (resp []byte, err error)
	SetState(req []byte) (resp []byte, err error)
	ReadVolumeFileStatus(req []byte) (resp []byte, err error)
	ReadNeedleBlob(req []byte) (resp []byte, err error)
	ReadNeedleMeta(req []byte) (resp []byte, err error)
	WriteNeedleBlob(req []byte) (resp []byte, err error)
	VolumeTailReceiver(req []byte) (resp []byte, err error)
	VolumeEcShardsGenerate(req []byte) (resp []byte, err error)
	VolumeEcShardsRebuild(req []byte) (resp []byte, err error)
	VolumeEcShardsCopy(req []byte) (resp []byte, err error)
	VolumeEcShardsDelete(req []byte) (resp []byte, err error)
	VolumeEcShardsMount(req []byte) (resp []byte, err error)
	VolumeEcShardsUnmount(req []byte) (resp []byte, err error)
	VolumeEcBlobDelete(req []byte) (resp []byte, err error)
	VolumeEcShardsToVolume(req []byte) (resp []byte, err error)
	VolumeEcShardsInfo(req []byte) (resp []byte, err error)
	VolumeServerStatus(req []byte) (resp []byte, err error)
	VolumeServerLeave(req []byte) (resp []byte, err error)
	FetchAndWriteNeedle(req []byte) (resp []byte, err error)
	ScrubVolume(req []byte) (resp []byte, err error)
	ScrubEcVolume(req []byte) (resp []byte, err error)
	VolumeNeedleStatus(req []byte) (resp []byte, err error)
	Ping(req []byte) (resp []byte, err error)

	// --- server-streaming (10): init request, then the server streams items ---
	VacuumVolumeCompact(stream VacuumVolumeCompactServerStream) error
	VolumeIncrementalCopy(stream VolumeIncrementalCopyServerStream) error
	VolumeCopy(stream VolumeCopyServerStream) error
	CopyFile(stream CopyFileServerStream) error
	ReadAllNeedles(stream ReadAllNeedlesServerStream) error
	VolumeTailSender(stream VolumeTailSenderServerStream) error
	VolumeEcShardRead(stream VolumeEcShardReadServerStream) error
	VolumeTierMoveDatToRemote(stream VolumeTierMoveDatToRemoteServerStream) error
	VolumeTierMoveDatFromRemote(stream VolumeTierMoveDatFromRemoteServerStream) error
	Query(stream QueryServerStream) error

	// --- client-streaming (1): the client streams request frames, then one reply ---
	ReceiveFile(stream ReceiveFileServerStream) error
}

// --- unary dispatch ---

// unaryHandler routes the 37 unary ordinals to the backend store. The 11
// streaming ordinals are NOT served here — they arrive as stream-open frames
// routed to StreamHandler, never as Call envelopes — so an unexpected unary call
// on a streaming ordinal (or any unknown one) yields StatusNotFound. This is the
// same split mq_agentwire uses between its unary Dispatch and its StreamHandler.
type unaryHandler struct{ store VolumeServerStore }

func (h unaryHandler) dispatch(envelope []byte) ([]byte, error) {
	call, err := rpc.ParseRequest(envelope)
	if err != nil {
		return nil, err
	}
	var (
		body []byte
		herr error
	)
	switch call.Method {
	case VolumeServerBatchDeleteOrdinal:
		body, herr = h.store.BatchDelete(call.Payload)
	case VolumeServerVacuumVolumeCheckOrdinal:
		body, herr = h.store.VacuumVolumeCheck(call.Payload)
	case VolumeServerVacuumVolumeCommitOrdinal:
		body, herr = h.store.VacuumVolumeCommit(call.Payload)
	case VolumeServerVacuumVolumeCleanupOrdinal:
		body, herr = h.store.VacuumVolumeCleanup(call.Payload)
	case VolumeServerDeleteCollectionOrdinal:
		body, herr = h.store.DeleteCollection(call.Payload)
	case VolumeServerAllocateVolumeOrdinal:
		body, herr = h.store.AllocateVolume(call.Payload)
	case VolumeServerVolumeSyncStatusOrdinal:
		body, herr = h.store.VolumeSyncStatus(call.Payload)
	case VolumeServerVolumeMountOrdinal:
		body, herr = h.store.VolumeMount(call.Payload)
	case VolumeServerVolumeUnmountOrdinal:
		body, herr = h.store.VolumeUnmount(call.Payload)
	case VolumeServerVolumeDeleteOrdinal:
		body, herr = h.store.VolumeDelete(call.Payload)
	case VolumeServerVolumeMarkReadonlyOrdinal:
		body, herr = h.store.VolumeMarkReadonly(call.Payload)
	case VolumeServerVolumeMarkWritableOrdinal:
		body, herr = h.store.VolumeMarkWritable(call.Payload)
	case VolumeServerVolumeConfigureOrdinal:
		body, herr = h.store.VolumeConfigure(call.Payload)
	case VolumeServerVolumeStatusOrdinal:
		body, herr = h.store.VolumeStatus(call.Payload)
	case VolumeServerGetStateOrdinal:
		body, herr = h.store.GetState(call.Payload)
	case VolumeServerSetStateOrdinal:
		body, herr = h.store.SetState(call.Payload)
	case VolumeServerReadVolumeFileStatusOrdinal:
		body, herr = h.store.ReadVolumeFileStatus(call.Payload)
	case VolumeServerReadNeedleBlobOrdinal:
		body, herr = h.store.ReadNeedleBlob(call.Payload)
	case VolumeServerReadNeedleMetaOrdinal:
		body, herr = h.store.ReadNeedleMeta(call.Payload)
	case VolumeServerWriteNeedleBlobOrdinal:
		body, herr = h.store.WriteNeedleBlob(call.Payload)
	case VolumeServerVolumeTailReceiverOrdinal:
		body, herr = h.store.VolumeTailReceiver(call.Payload)
	case VolumeServerVolumeEcShardsGenerateOrdinal:
		body, herr = h.store.VolumeEcShardsGenerate(call.Payload)
	case VolumeServerVolumeEcShardsRebuildOrdinal:
		body, herr = h.store.VolumeEcShardsRebuild(call.Payload)
	case VolumeServerVolumeEcShardsCopyOrdinal:
		body, herr = h.store.VolumeEcShardsCopy(call.Payload)
	case VolumeServerVolumeEcShardsDeleteOrdinal:
		body, herr = h.store.VolumeEcShardsDelete(call.Payload)
	case VolumeServerVolumeEcShardsMountOrdinal:
		body, herr = h.store.VolumeEcShardsMount(call.Payload)
	case VolumeServerVolumeEcShardsUnmountOrdinal:
		body, herr = h.store.VolumeEcShardsUnmount(call.Payload)
	case VolumeServerVolumeEcBlobDeleteOrdinal:
		body, herr = h.store.VolumeEcBlobDelete(call.Payload)
	case VolumeServerVolumeEcShardsToVolumeOrdinal:
		body, herr = h.store.VolumeEcShardsToVolume(call.Payload)
	case VolumeServerVolumeEcShardsInfoOrdinal:
		body, herr = h.store.VolumeEcShardsInfo(call.Payload)
	case VolumeServerVolumeServerStatusOrdinal:
		body, herr = h.store.VolumeServerStatus(call.Payload)
	case VolumeServerVolumeServerLeaveOrdinal:
		body, herr = h.store.VolumeServerLeave(call.Payload)
	case VolumeServerFetchAndWriteNeedleOrdinal:
		body, herr = h.store.FetchAndWriteNeedle(call.Payload)
	case VolumeServerScrubVolumeOrdinal:
		body, herr = h.store.ScrubVolume(call.Payload)
	case VolumeServerScrubEcVolumeOrdinal:
		body, herr = h.store.ScrubEcVolume(call.Payload)
	case VolumeServerVolumeNeedleStatusOrdinal:
		body, herr = h.store.VolumeNeedleStatus(call.Payload)
	case VolumeServerPingOrdinal:
		body, herr = h.store.Ping(call.Payload)
	default:
		// Streaming ordinals (VacuumVolumeCompact, VolumeIncrementalCopy,
		// VolumeCopy, CopyFile, ReceiveFile, ReadAllNeedles, VolumeTailSender,
		// VolumeEcShardRead, VolumeTierMoveDat*, Query) and any unknown ordinal:
		// not a unary call here.
		return rpc.BuildResponse(rpc.StatusNotFound, call.PromiseID, nil), nil
	}
	if herr != nil {
		// Carry the handler error message so the caller can reconstruct sentinels.
		return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(herr.Error())), nil
	}
	return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
}

// Dispatch is the unary VolumeServer server dispatch bound to store — pass it as
// the unary leg of transport.ListenStream (or, if a deployment never serves the
// streaming RPCs, directly to transport.Listen). It serves only the 37 unary
// methods; the streaming methods are served by StreamHandler. Exposed so callers
// can compose it (e.g. behind PQ-TLS).
func Dispatch(store VolumeServerStore) transport.Dispatch {
	h := unaryHandler{store: store}
	return h.dispatch
}

// --- streaming server side ---

// serverStream adapts a transport.Stream to the per-RPC *ServerStream
// interfaces below. Every streaming method shares the same frame plumbing (raw
// ZAP buffers in/out); only the typed wrappers around the buffers differ, so one
// adapter serves all of them. The opening request frame rides as init and is
// replayed as the first Recv (mirroring mq_agentwire.serverStream and
// s3/admin/plugin/worker_stream_zap.go).
type serverStream struct {
	s     transport.Stream
	init  []byte
	hasIn bool
}

func newServerStream(init []byte, s transport.Stream) *serverStream {
	return &serverStream{s: s, init: init, hasIn: len(init) > 0}
}

// Recv returns the next inbound frame as a raw buffer, replaying the init frame
// first; io.EOF once the peer half-closes.
func (z *serverStream) Recv() ([]byte, error) {
	if z.hasIn {
		z.hasIn = false
		return z.init, nil
	}
	return z.s.Recv()
}

// Send streams one outbound frame.
func (z *serverStream) Send(body []byte) error { return z.s.Send(body) }

// Each server-streaming RPC's server view: Init() yields the (single) opening
// request; Send streams one response item. The client half-closes immediately
// after the open, so the server only ever Sends.
//
//	VacuumVolumeCompact: VacuumVolumeCompactRequest -> stream VacuumVolumeCompactResponse
type VacuumVolumeCompactServerStream struct{ z *serverStream }

// Init returns the opening VacuumVolumeCompactRequest view.
func (s VacuumVolumeCompactServerStream) Init() (VacuumVolumeCompactRequest, error) {
	b, err := s.z.Recv()
	if err != nil {
		return VacuumVolumeCompactRequest{}, err
	}
	return WrapVacuumVolumeCompactRequest(b)
}

// Send streams one VacuumVolumeCompactResponse item.
func (s VacuumVolumeCompactServerStream) Send(in VacuumVolumeCompactResponseInput) error {
	return s.z.Send(NewVacuumVolumeCompactResponse(in))
}

// VolumeIncrementalCopy: VolumeIncrementalCopyRequest -> stream VolumeIncrementalCopyResponse
type VolumeIncrementalCopyServerStream struct{ z *serverStream }

func (s VolumeIncrementalCopyServerStream) Init() (VolumeIncrementalCopyRequest, error) {
	b, err := s.z.Recv()
	if err != nil {
		return VolumeIncrementalCopyRequest{}, err
	}
	return WrapVolumeIncrementalCopyRequest(b)
}

func (s VolumeIncrementalCopyServerStream) Send(in VolumeIncrementalCopyResponseInput) error {
	return s.z.Send(NewVolumeIncrementalCopyResponse(in))
}

// VolumeCopy: VolumeCopyRequest -> stream VolumeCopyResponse
type VolumeCopyServerStream struct{ z *serverStream }

func (s VolumeCopyServerStream) Init() (VolumeCopyRequest, error) {
	b, err := s.z.Recv()
	if err != nil {
		return VolumeCopyRequest{}, err
	}
	return WrapVolumeCopyRequest(b)
}

func (s VolumeCopyServerStream) Send(in VolumeCopyResponseInput) error {
	return s.z.Send(NewVolumeCopyResponse(in))
}

// CopyFile: CopyFileRequest -> stream CopyFileResponse
type CopyFileServerStream struct{ z *serverStream }

func (s CopyFileServerStream) Init() (CopyFileRequest, error) {
	b, err := s.z.Recv()
	if err != nil {
		return CopyFileRequest{}, err
	}
	return WrapCopyFileRequest(b)
}

func (s CopyFileServerStream) Send(in CopyFileResponseInput) error {
	return s.z.Send(NewCopyFileResponse(in))
}

// ReadAllNeedles: ReadAllNeedlesRequest -> stream ReadAllNeedlesResponse
type ReadAllNeedlesServerStream struct{ z *serverStream }

func (s ReadAllNeedlesServerStream) Init() (ReadAllNeedlesRequest, error) {
	b, err := s.z.Recv()
	if err != nil {
		return ReadAllNeedlesRequest{}, err
	}
	return WrapReadAllNeedlesRequest(b)
}

func (s ReadAllNeedlesServerStream) Send(in ReadAllNeedlesResponseInput) error {
	return s.z.Send(NewReadAllNeedlesResponse(in))
}

// VolumeTailSender: VolumeTailSenderRequest -> stream VolumeTailSenderResponse
type VolumeTailSenderServerStream struct{ z *serverStream }

func (s VolumeTailSenderServerStream) Init() (VolumeTailSenderRequest, error) {
	b, err := s.z.Recv()
	if err != nil {
		return VolumeTailSenderRequest{}, err
	}
	return WrapVolumeTailSenderRequest(b)
}

func (s VolumeTailSenderServerStream) Send(in VolumeTailSenderResponseInput) error {
	return s.z.Send(NewVolumeTailSenderResponse(in))
}

// VolumeEcShardRead: VolumeEcShardReadRequest -> stream VolumeEcShardReadResponse
type VolumeEcShardReadServerStream struct{ z *serverStream }

func (s VolumeEcShardReadServerStream) Init() (VolumeEcShardReadRequest, error) {
	b, err := s.z.Recv()
	if err != nil {
		return VolumeEcShardReadRequest{}, err
	}
	return WrapVolumeEcShardReadRequest(b)
}

func (s VolumeEcShardReadServerStream) Send(in VolumeEcShardReadResponseInput) error {
	return s.z.Send(NewVolumeEcShardReadResponse(in))
}

// VolumeTierMoveDatToRemote: VolumeTierMoveDatToRemoteRequest -> stream VolumeTierMoveDatToRemoteResponse
type VolumeTierMoveDatToRemoteServerStream struct{ z *serverStream }

func (s VolumeTierMoveDatToRemoteServerStream) Init() (VolumeTierMoveDatToRemoteRequest, error) {
	b, err := s.z.Recv()
	if err != nil {
		return VolumeTierMoveDatToRemoteRequest{}, err
	}
	return WrapVolumeTierMoveDatToRemoteRequest(b)
}

func (s VolumeTierMoveDatToRemoteServerStream) Send(in VolumeTierMoveDatToRemoteResponseInput) error {
	return s.z.Send(NewVolumeTierMoveDatToRemoteResponse(in))
}

// VolumeTierMoveDatFromRemote: VolumeTierMoveDatFromRemoteRequest -> stream VolumeTierMoveDatFromRemoteResponse
type VolumeTierMoveDatFromRemoteServerStream struct{ z *serverStream }

func (s VolumeTierMoveDatFromRemoteServerStream) Init() (VolumeTierMoveDatFromRemoteRequest, error) {
	b, err := s.z.Recv()
	if err != nil {
		return VolumeTierMoveDatFromRemoteRequest{}, err
	}
	return WrapVolumeTierMoveDatFromRemoteRequest(b)
}

func (s VolumeTierMoveDatFromRemoteServerStream) Send(in VolumeTierMoveDatFromRemoteResponseInput) error {
	return s.z.Send(NewVolumeTierMoveDatFromRemoteResponse(in))
}

// Query: QueryRequest -> stream QueriedStripe
type QueryServerStream struct{ z *serverStream }

func (s QueryServerStream) Init() (QueryRequest, error) {
	b, err := s.z.Recv()
	if err != nil {
		return QueryRequest{}, err
	}
	return WrapQueryRequest(b)
}

func (s QueryServerStream) Send(in QueriedStripeInput) error {
	return s.z.Send(NewQueriedStripe(in))
}

// ReceiveFileServerStream is the server view of the one client-streaming RPC:
// the client streams ReceiveFileRequest frames (the first carries ReceiveFileInfo,
// the rest carry file_content chunks); the server replies with exactly one
// ReceiveFileResponse via Reply when the client half-closes.
//
//	ReceiveFile: stream ReceiveFileRequest -> ReceiveFileResponse
type ReceiveFileServerStream struct{ z *serverStream }

// Recv returns the next ReceiveFileRequest frame, or io.EOF once the client
// half-closes.
func (s ReceiveFileServerStream) Recv() (ReceiveFileRequest, error) {
	b, err := s.z.Recv()
	if err != nil {
		return ReceiveFileRequest{}, err
	}
	return WrapReceiveFileRequest(b)
}

// Reply sends the single terminal ReceiveFileResponse.
func (s ReceiveFileServerStream) Reply(in ReceiveFileResponseInput) error {
	return s.z.Send(NewReceiveFileResponse(in))
}

// StreamHandler is the transport.StreamHandler for VolumeServer's 11 streaming
// ordinals: it adapts each accepted ZAP stream to the matching server-stream
// view and runs the backend's handler against it. Pass it as the stream leg of
// transport.ListenStream alongside Dispatch(store). Returning from a backend
// handler half-closes the send side (the client's Recv then sees io.EOF), per
// transport.StreamHandler semantics.
func StreamHandler(store VolumeServerStore) transport.StreamHandler {
	return func(method uint32, init []byte, s transport.Stream) {
		z := newServerStream(init, s)
		switch method {
		case VolumeServerVacuumVolumeCompactOrdinal:
			_ = store.VacuumVolumeCompact(VacuumVolumeCompactServerStream{z: z})
		case VolumeServerVolumeIncrementalCopyOrdinal:
			_ = store.VolumeIncrementalCopy(VolumeIncrementalCopyServerStream{z: z})
		case VolumeServerVolumeCopyOrdinal:
			_ = store.VolumeCopy(VolumeCopyServerStream{z: z})
		case VolumeServerCopyFileOrdinal:
			_ = store.CopyFile(CopyFileServerStream{z: z})
		case VolumeServerReadAllNeedlesOrdinal:
			_ = store.ReadAllNeedles(ReadAllNeedlesServerStream{z: z})
		case VolumeServerVolumeTailSenderOrdinal:
			_ = store.VolumeTailSender(VolumeTailSenderServerStream{z: z})
		case VolumeServerVolumeEcShardReadOrdinal:
			_ = store.VolumeEcShardRead(VolumeEcShardReadServerStream{z: z})
		case VolumeServerVolumeTierMoveDatToRemoteOrdinal:
			_ = store.VolumeTierMoveDatToRemote(VolumeTierMoveDatToRemoteServerStream{z: z})
		case VolumeServerVolumeTierMoveDatFromRemoteOrdinal:
			_ = store.VolumeTierMoveDatFromRemote(VolumeTierMoveDatFromRemoteServerStream{z: z})
		case VolumeServerQueryOrdinal:
			_ = store.Query(QueryServerStream{z: z})
		case VolumeServerReceiveFileOrdinal:
			_ = store.ReceiveFile(ReceiveFileServerStream{z: z})
		default:
			// Unknown / unary ordinal opened as a stream: returning half-closes;
			// the client's Recv sees io.EOF.
		}
	}
}

// Serve starts the native ZAP VolumeServer service on network/addr (e.g. "tcp",
// ":18906"), backed by store. Unary calls route through Dispatch; the 11 streams
// route through StreamHandler — both over ONE listener (transport.ListenStream),
// exactly like the mq_agent ZAP server. Returns the running server; Close stops
// it. For the PQ-secured mesh, build a TLS/QUIC transport.Conn and serve the
// streams via conn.OpenStream against a ListenStream on the TLS listener.
func Serve(network, addr string, store VolumeServerStore) (*transport.Server, error) {
	return transport.ListenStream(network, addr, Dispatch(store), StreamHandler(store))
}

// --- client side ---

// Client is the typed VolumeServer ZAP client internal callers (master, filer,
// shell/admin, the volume-server-to-volume-server replication paths) hold. It
// owns the transport connection to one volume-server endpoint and serves both
// the unary calls and the streaming opens over it.
type Client struct {
	conn transport.Conn
	rpc  *VolumeServerClient
}

// Dial connects to the volume-server ZAP service at addr over plain TCP (e.g.
// "volume.hanzo.svc:18906"). For the PQ-secured mesh, build the transport.Conn
// via transport.DialTLS with a *tls.Config and use NewClient.
func Dial(network, addr string) (*Client, error) {
	conn, err := transport.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn, rpc: NewVolumeServerClient(conn, nil)}, nil
}

// NewClient wraps an already-established transport.Conn (TCP, Unix, or PQ-TLS).
func NewClient(conn transport.Conn) *Client {
	return &Client{conn: conn, rpc: NewVolumeServerClient(conn, nil)}
}

// Close releases the underlying connection.
func (c *Client) Close() error { return c.conn.Close() }

// Conn exposes the underlying transport connection, for callers that open
// streams directly (the streaming client opens below use it).
func (c *Client) Conn() transport.Conn { return c.conn }

// The typed unary client methods take a New<Req> input and return the Wrap<Resp>
// view. They mirror s3/zapsvc Client.GetObject/PutObject and the mq_agent typed
// client. Each one decodes the reply into its zero-copy view — the bytes never
// leave this package as structs.

func (c *Client) BatchDelete(in BatchDeleteRequestInput) (BatchDeleteResponse, error) {
	_, body, err := c.rpc.BatchDelete(NewBatchDeleteRequest(in))
	if err != nil {
		return BatchDeleteResponse{}, err
	}
	return WrapBatchDeleteResponse(body)
}

func (c *Client) VacuumVolumeCheck(in VacuumVolumeCheckRequestInput) (VacuumVolumeCheckResponse, error) {
	_, body, err := c.rpc.VacuumVolumeCheck(NewVacuumVolumeCheckRequest(in))
	if err != nil {
		return VacuumVolumeCheckResponse{}, err
	}
	return WrapVacuumVolumeCheckResponse(body)
}

func (c *Client) VacuumVolumeCommit(in VacuumVolumeCommitRequestInput) (VacuumVolumeCommitResponse, error) {
	_, body, err := c.rpc.VacuumVolumeCommit(NewVacuumVolumeCommitRequest(in))
	if err != nil {
		return VacuumVolumeCommitResponse{}, err
	}
	return WrapVacuumVolumeCommitResponse(body)
}

func (c *Client) VacuumVolumeCleanup(in VacuumVolumeCleanupRequestInput) (VacuumVolumeCleanupResponse, error) {
	_, body, err := c.rpc.VacuumVolumeCleanup(NewVacuumVolumeCleanupRequest(in))
	if err != nil {
		return VacuumVolumeCleanupResponse{}, err
	}
	return WrapVacuumVolumeCleanupResponse(body)
}

func (c *Client) DeleteCollection(in DeleteCollectionRequestInput) (DeleteCollectionResponse, error) {
	_, body, err := c.rpc.DeleteCollection(NewDeleteCollectionRequest(in))
	if err != nil {
		return DeleteCollectionResponse{}, err
	}
	return WrapDeleteCollectionResponse(body)
}

func (c *Client) AllocateVolume(in AllocateVolumeRequestInput) (AllocateVolumeResponse, error) {
	_, body, err := c.rpc.AllocateVolume(NewAllocateVolumeRequest(in))
	if err != nil {
		return AllocateVolumeResponse{}, err
	}
	return WrapAllocateVolumeResponse(body)
}

func (c *Client) VolumeSyncStatus(in VolumeSyncStatusRequestInput) (VolumeSyncStatusResponse, error) {
	_, body, err := c.rpc.VolumeSyncStatus(NewVolumeSyncStatusRequest(in))
	if err != nil {
		return VolumeSyncStatusResponse{}, err
	}
	return WrapVolumeSyncStatusResponse(body)
}

func (c *Client) VolumeMount(in VolumeMountRequestInput) (VolumeMountResponse, error) {
	_, body, err := c.rpc.VolumeMount(NewVolumeMountRequest(in))
	if err != nil {
		return VolumeMountResponse{}, err
	}
	return WrapVolumeMountResponse(body)
}

func (c *Client) VolumeUnmount(in VolumeUnmountRequestInput) (VolumeUnmountResponse, error) {
	_, body, err := c.rpc.VolumeUnmount(NewVolumeUnmountRequest(in))
	if err != nil {
		return VolumeUnmountResponse{}, err
	}
	return WrapVolumeUnmountResponse(body)
}

func (c *Client) VolumeDelete(in VolumeDeleteRequestInput) (VolumeDeleteResponse, error) {
	_, body, err := c.rpc.VolumeDelete(NewVolumeDeleteRequest(in))
	if err != nil {
		return VolumeDeleteResponse{}, err
	}
	return WrapVolumeDeleteResponse(body)
}

func (c *Client) VolumeMarkReadonly(in VolumeMarkReadonlyRequestInput) (VolumeMarkReadonlyResponse, error) {
	_, body, err := c.rpc.VolumeMarkReadonly(NewVolumeMarkReadonlyRequest(in))
	if err != nil {
		return VolumeMarkReadonlyResponse{}, err
	}
	return WrapVolumeMarkReadonlyResponse(body)
}

func (c *Client) VolumeMarkWritable(in VolumeMarkWritableRequestInput) (VolumeMarkWritableResponse, error) {
	_, body, err := c.rpc.VolumeMarkWritable(NewVolumeMarkWritableRequest(in))
	if err != nil {
		return VolumeMarkWritableResponse{}, err
	}
	return WrapVolumeMarkWritableResponse(body)
}

func (c *Client) VolumeConfigure(in VolumeConfigureRequestInput) (VolumeConfigureResponse, error) {
	_, body, err := c.rpc.VolumeConfigure(NewVolumeConfigureRequest(in))
	if err != nil {
		return VolumeConfigureResponse{}, err
	}
	return WrapVolumeConfigureResponse(body)
}

func (c *Client) VolumeStatus(in VolumeStatusRequestInput) (VolumeStatusResponse, error) {
	_, body, err := c.rpc.VolumeStatus(NewVolumeStatusRequest(in))
	if err != nil {
		return VolumeStatusResponse{}, err
	}
	return WrapVolumeStatusResponse(body)
}

func (c *Client) GetState(in GetStateRequestInput) (GetStateResponse, error) {
	_, body, err := c.rpc.GetState(NewGetStateRequest(in))
	if err != nil {
		return GetStateResponse{}, err
	}
	return WrapGetStateResponse(body)
}

func (c *Client) SetState(in SetStateRequestInput) (SetStateResponse, error) {
	_, body, err := c.rpc.SetState(NewSetStateRequest(in))
	if err != nil {
		return SetStateResponse{}, err
	}
	return WrapSetStateResponse(body)
}

func (c *Client) ReadVolumeFileStatus(in ReadVolumeFileStatusRequestInput) (ReadVolumeFileStatusResponse, error) {
	_, body, err := c.rpc.ReadVolumeFileStatus(NewReadVolumeFileStatusRequest(in))
	if err != nil {
		return ReadVolumeFileStatusResponse{}, err
	}
	return WrapReadVolumeFileStatusResponse(body)
}

func (c *Client) ReadNeedleBlob(in ReadNeedleBlobRequestInput) (ReadNeedleBlobResponse, error) {
	_, body, err := c.rpc.ReadNeedleBlob(NewReadNeedleBlobRequest(in))
	if err != nil {
		return ReadNeedleBlobResponse{}, err
	}
	return WrapReadNeedleBlobResponse(body)
}

func (c *Client) ReadNeedleMeta(in ReadNeedleMetaRequestInput) (ReadNeedleMetaResponse, error) {
	_, body, err := c.rpc.ReadNeedleMeta(NewReadNeedleMetaRequest(in))
	if err != nil {
		return ReadNeedleMetaResponse{}, err
	}
	return WrapReadNeedleMetaResponse(body)
}

func (c *Client) WriteNeedleBlob(in WriteNeedleBlobRequestInput) (WriteNeedleBlobResponse, error) {
	_, body, err := c.rpc.WriteNeedleBlob(NewWriteNeedleBlobRequest(in))
	if err != nil {
		return WriteNeedleBlobResponse{}, err
	}
	return WrapWriteNeedleBlobResponse(body)
}

func (c *Client) VolumeTailReceiver(in VolumeTailReceiverRequestInput) (VolumeTailReceiverResponse, error) {
	_, body, err := c.rpc.VolumeTailReceiver(NewVolumeTailReceiverRequest(in))
	if err != nil {
		return VolumeTailReceiverResponse{}, err
	}
	return WrapVolumeTailReceiverResponse(body)
}

func (c *Client) VolumeEcShardsGenerate(in VolumeEcShardsGenerateRequestInput) (VolumeEcShardsGenerateResponse, error) {
	_, body, err := c.rpc.VolumeEcShardsGenerate(NewVolumeEcShardsGenerateRequest(in))
	if err != nil {
		return VolumeEcShardsGenerateResponse{}, err
	}
	return WrapVolumeEcShardsGenerateResponse(body)
}

func (c *Client) VolumeEcShardsRebuild(in VolumeEcShardsRebuildRequestInput) (VolumeEcShardsRebuildResponse, error) {
	_, body, err := c.rpc.VolumeEcShardsRebuild(NewVolumeEcShardsRebuildRequest(in))
	if err != nil {
		return VolumeEcShardsRebuildResponse{}, err
	}
	return WrapVolumeEcShardsRebuildResponse(body)
}

func (c *Client) VolumeEcShardsCopy(in VolumeEcShardsCopyRequestInput) (VolumeEcShardsCopyResponse, error) {
	_, body, err := c.rpc.VolumeEcShardsCopy(NewVolumeEcShardsCopyRequest(in))
	if err != nil {
		return VolumeEcShardsCopyResponse{}, err
	}
	return WrapVolumeEcShardsCopyResponse(body)
}

func (c *Client) VolumeEcShardsDelete(in VolumeEcShardsDeleteRequestInput) (VolumeEcShardsDeleteResponse, error) {
	_, body, err := c.rpc.VolumeEcShardsDelete(NewVolumeEcShardsDeleteRequest(in))
	if err != nil {
		return VolumeEcShardsDeleteResponse{}, err
	}
	return WrapVolumeEcShardsDeleteResponse(body)
}

func (c *Client) VolumeEcShardsMount(in VolumeEcShardsMountRequestInput) (VolumeEcShardsMountResponse, error) {
	_, body, err := c.rpc.VolumeEcShardsMount(NewVolumeEcShardsMountRequest(in))
	if err != nil {
		return VolumeEcShardsMountResponse{}, err
	}
	return WrapVolumeEcShardsMountResponse(body)
}

func (c *Client) VolumeEcShardsUnmount(in VolumeEcShardsUnmountRequestInput) (VolumeEcShardsUnmountResponse, error) {
	_, body, err := c.rpc.VolumeEcShardsUnmount(NewVolumeEcShardsUnmountRequest(in))
	if err != nil {
		return VolumeEcShardsUnmountResponse{}, err
	}
	return WrapVolumeEcShardsUnmountResponse(body)
}

func (c *Client) VolumeEcBlobDelete(in VolumeEcBlobDeleteRequestInput) (VolumeEcBlobDeleteResponse, error) {
	_, body, err := c.rpc.VolumeEcBlobDelete(NewVolumeEcBlobDeleteRequest(in))
	if err != nil {
		return VolumeEcBlobDeleteResponse{}, err
	}
	return WrapVolumeEcBlobDeleteResponse(body)
}

func (c *Client) VolumeEcShardsToVolume(in VolumeEcShardsToVolumeRequestInput) (VolumeEcShardsToVolumeResponse, error) {
	_, body, err := c.rpc.VolumeEcShardsToVolume(NewVolumeEcShardsToVolumeRequest(in))
	if err != nil {
		return VolumeEcShardsToVolumeResponse{}, err
	}
	return WrapVolumeEcShardsToVolumeResponse(body)
}

func (c *Client) VolumeEcShardsInfo(in VolumeEcShardsInfoRequestInput) (VolumeEcShardsInfoResponse, error) {
	_, body, err := c.rpc.VolumeEcShardsInfo(NewVolumeEcShardsInfoRequest(in))
	if err != nil {
		return VolumeEcShardsInfoResponse{}, err
	}
	return WrapVolumeEcShardsInfoResponse(body)
}

func (c *Client) VolumeServerStatus(in VolumeServerStatusRequestInput) (VolumeServerStatusResponse, error) {
	_, body, err := c.rpc.VolumeServerStatus(NewVolumeServerStatusRequest(in))
	if err != nil {
		return VolumeServerStatusResponse{}, err
	}
	return WrapVolumeServerStatusResponse(body)
}

func (c *Client) VolumeServerLeave(in VolumeServerLeaveRequestInput) (VolumeServerLeaveResponse, error) {
	_, body, err := c.rpc.VolumeServerLeave(NewVolumeServerLeaveRequest(in))
	if err != nil {
		return VolumeServerLeaveResponse{}, err
	}
	return WrapVolumeServerLeaveResponse(body)
}

func (c *Client) FetchAndWriteNeedle(in FetchAndWriteNeedleRequestInput) (FetchAndWriteNeedleResponse, error) {
	_, body, err := c.rpc.FetchAndWriteNeedle(NewFetchAndWriteNeedleRequest(in))
	if err != nil {
		return FetchAndWriteNeedleResponse{}, err
	}
	return WrapFetchAndWriteNeedleResponse(body)
}

func (c *Client) ScrubVolume(in ScrubVolumeRequestInput) (ScrubVolumeResponse, error) {
	_, body, err := c.rpc.ScrubVolume(NewScrubVolumeRequest(in))
	if err != nil {
		return ScrubVolumeResponse{}, err
	}
	return WrapScrubVolumeResponse(body)
}

func (c *Client) ScrubEcVolume(in ScrubEcVolumeRequestInput) (ScrubEcVolumeResponse, error) {
	_, body, err := c.rpc.ScrubEcVolume(NewScrubEcVolumeRequest(in))
	if err != nil {
		return ScrubEcVolumeResponse{}, err
	}
	return WrapScrubEcVolumeResponse(body)
}

func (c *Client) VolumeNeedleStatus(in VolumeNeedleStatusRequestInput) (VolumeNeedleStatusResponse, error) {
	_, body, err := c.rpc.VolumeNeedleStatus(NewVolumeNeedleStatusRequest(in))
	if err != nil {
		return VolumeNeedleStatusResponse{}, err
	}
	return WrapVolumeNeedleStatusResponse(body)
}

func (c *Client) Ping(in PingRequestInput) (PingResponse, error) {
	_, body, err := c.rpc.Ping(NewPingRequest(in))
	if err != nil {
		return PingResponse{}, err
	}
	return WrapPingResponse(body)
}

// --- streaming client opens ---

// Each server-streaming client open issues the open frame carrying the request
// as init, then returns a typed reader. The client immediately half-closes the
// send side (these RPCs send nothing after the open) and Recv()s items until
// io.EOF.

// VacuumVolumeCompactClientStream reads VacuumVolumeCompactResponse items.
type VacuumVolumeCompactClientStream struct{ s transport.Stream }

// Recv returns the next VacuumVolumeCompactResponse, or io.EOF at end of stream.
func (r *VacuumVolumeCompactClientStream) Recv() (VacuumVolumeCompactResponse, error) {
	b, err := r.s.Recv()
	if err != nil {
		return VacuumVolumeCompactResponse{}, err
	}
	return WrapVacuumVolumeCompactResponse(b)
}

// VacuumVolumeCompact opens the server-stream; the request rides as init.
func (c *Client) VacuumVolumeCompact(in VacuumVolumeCompactRequestInput) (*VacuumVolumeCompactClientStream, error) {
	s, err := c.conn.OpenStream(VolumeServerVacuumVolumeCompactOrdinal, NewVacuumVolumeCompactRequest(in))
	if err != nil {
		return nil, err
	}
	_ = s.CloseSend()
	return &VacuumVolumeCompactClientStream{s: s}, nil
}

// VolumeIncrementalCopyClientStream reads VolumeIncrementalCopyResponse items.
type VolumeIncrementalCopyClientStream struct{ s transport.Stream }

func (r *VolumeIncrementalCopyClientStream) Recv() (VolumeIncrementalCopyResponse, error) {
	b, err := r.s.Recv()
	if err != nil {
		return VolumeIncrementalCopyResponse{}, err
	}
	return WrapVolumeIncrementalCopyResponse(b)
}

func (c *Client) VolumeIncrementalCopy(in VolumeIncrementalCopyRequestInput) (*VolumeIncrementalCopyClientStream, error) {
	s, err := c.conn.OpenStream(VolumeServerVolumeIncrementalCopyOrdinal, NewVolumeIncrementalCopyRequest(in))
	if err != nil {
		return nil, err
	}
	_ = s.CloseSend()
	return &VolumeIncrementalCopyClientStream{s: s}, nil
}

// VolumeCopyClientStream reads VolumeCopyResponse items.
type VolumeCopyClientStream struct{ s transport.Stream }

func (r *VolumeCopyClientStream) Recv() (VolumeCopyResponse, error) {
	b, err := r.s.Recv()
	if err != nil {
		return VolumeCopyResponse{}, err
	}
	return WrapVolumeCopyResponse(b)
}

func (c *Client) VolumeCopy(in VolumeCopyRequestInput) (*VolumeCopyClientStream, error) {
	s, err := c.conn.OpenStream(VolumeServerVolumeCopyOrdinal, NewVolumeCopyRequest(in))
	if err != nil {
		return nil, err
	}
	_ = s.CloseSend()
	return &VolumeCopyClientStream{s: s}, nil
}

// CopyFileClientStream reads CopyFileResponse items.
type CopyFileClientStream struct{ s transport.Stream }

func (r *CopyFileClientStream) Recv() (CopyFileResponse, error) {
	b, err := r.s.Recv()
	if err != nil {
		return CopyFileResponse{}, err
	}
	return WrapCopyFileResponse(b)
}

func (c *Client) CopyFile(in CopyFileRequestInput) (*CopyFileClientStream, error) {
	s, err := c.conn.OpenStream(VolumeServerCopyFileOrdinal, NewCopyFileRequest(in))
	if err != nil {
		return nil, err
	}
	_ = s.CloseSend()
	return &CopyFileClientStream{s: s}, nil
}

// ReadAllNeedlesClientStream reads ReadAllNeedlesResponse items.
type ReadAllNeedlesClientStream struct{ s transport.Stream }

func (r *ReadAllNeedlesClientStream) Recv() (ReadAllNeedlesResponse, error) {
	b, err := r.s.Recv()
	if err != nil {
		return ReadAllNeedlesResponse{}, err
	}
	return WrapReadAllNeedlesResponse(b)
}

func (c *Client) ReadAllNeedles(in ReadAllNeedlesRequestInput) (*ReadAllNeedlesClientStream, error) {
	s, err := c.conn.OpenStream(VolumeServerReadAllNeedlesOrdinal, NewReadAllNeedlesRequest(in))
	if err != nil {
		return nil, err
	}
	_ = s.CloseSend()
	return &ReadAllNeedlesClientStream{s: s}, nil
}

// VolumeTailSenderClientStream reads VolumeTailSenderResponse items.
type VolumeTailSenderClientStream struct{ s transport.Stream }

func (r *VolumeTailSenderClientStream) Recv() (VolumeTailSenderResponse, error) {
	b, err := r.s.Recv()
	if err != nil {
		return VolumeTailSenderResponse{}, err
	}
	return WrapVolumeTailSenderResponse(b)
}

func (c *Client) VolumeTailSender(in VolumeTailSenderRequestInput) (*VolumeTailSenderClientStream, error) {
	s, err := c.conn.OpenStream(VolumeServerVolumeTailSenderOrdinal, NewVolumeTailSenderRequest(in))
	if err != nil {
		return nil, err
	}
	_ = s.CloseSend()
	return &VolumeTailSenderClientStream{s: s}, nil
}

// VolumeEcShardReadClientStream reads VolumeEcShardReadResponse items.
type VolumeEcShardReadClientStream struct{ s transport.Stream }

func (r *VolumeEcShardReadClientStream) Recv() (VolumeEcShardReadResponse, error) {
	b, err := r.s.Recv()
	if err != nil {
		return VolumeEcShardReadResponse{}, err
	}
	return WrapVolumeEcShardReadResponse(b)
}

func (c *Client) VolumeEcShardRead(in VolumeEcShardReadRequestInput) (*VolumeEcShardReadClientStream, error) {
	s, err := c.conn.OpenStream(VolumeServerVolumeEcShardReadOrdinal, NewVolumeEcShardReadRequest(in))
	if err != nil {
		return nil, err
	}
	_ = s.CloseSend()
	return &VolumeEcShardReadClientStream{s: s}, nil
}

// VolumeTierMoveDatToRemoteClientStream reads VolumeTierMoveDatToRemoteResponse items.
type VolumeTierMoveDatToRemoteClientStream struct{ s transport.Stream }

func (r *VolumeTierMoveDatToRemoteClientStream) Recv() (VolumeTierMoveDatToRemoteResponse, error) {
	b, err := r.s.Recv()
	if err != nil {
		return VolumeTierMoveDatToRemoteResponse{}, err
	}
	return WrapVolumeTierMoveDatToRemoteResponse(b)
}

func (c *Client) VolumeTierMoveDatToRemote(in VolumeTierMoveDatToRemoteRequestInput) (*VolumeTierMoveDatToRemoteClientStream, error) {
	s, err := c.conn.OpenStream(VolumeServerVolumeTierMoveDatToRemoteOrdinal, NewVolumeTierMoveDatToRemoteRequest(in))
	if err != nil {
		return nil, err
	}
	_ = s.CloseSend()
	return &VolumeTierMoveDatToRemoteClientStream{s: s}, nil
}

// VolumeTierMoveDatFromRemoteClientStream reads VolumeTierMoveDatFromRemoteResponse items.
type VolumeTierMoveDatFromRemoteClientStream struct{ s transport.Stream }

func (r *VolumeTierMoveDatFromRemoteClientStream) Recv() (VolumeTierMoveDatFromRemoteResponse, error) {
	b, err := r.s.Recv()
	if err != nil {
		return VolumeTierMoveDatFromRemoteResponse{}, err
	}
	return WrapVolumeTierMoveDatFromRemoteResponse(b)
}

func (c *Client) VolumeTierMoveDatFromRemote(in VolumeTierMoveDatFromRemoteRequestInput) (*VolumeTierMoveDatFromRemoteClientStream, error) {
	s, err := c.conn.OpenStream(VolumeServerVolumeTierMoveDatFromRemoteOrdinal, NewVolumeTierMoveDatFromRemoteRequest(in))
	if err != nil {
		return nil, err
	}
	_ = s.CloseSend()
	return &VolumeTierMoveDatFromRemoteClientStream{s: s}, nil
}

// QueryClientStream reads QueriedStripe items.
type QueryClientStream struct{ s transport.Stream }

func (r *QueryClientStream) Recv() (QueriedStripe, error) {
	b, err := r.s.Recv()
	if err != nil {
		return QueriedStripe{}, err
	}
	return WrapQueriedStripe(b)
}

func (c *Client) Query(in QueryRequestInput) (*QueryClientStream, error) {
	s, err := c.conn.OpenStream(VolumeServerQueryOrdinal, NewQueryRequest(in))
	if err != nil {
		return nil, err
	}
	_ = s.CloseSend()
	return &QueryClientStream{s: s}, nil
}

// ReceiveFileClientStream is the client view of the one client-streaming RPC:
// Send ReceiveFileRequest frames (first = info, then content chunks), CloseSend,
// then Reply() reads the single terminal ReceiveFileResponse.
type ReceiveFileClientStream struct{ s transport.Stream }

// Send streams one ReceiveFileRequest frame.
func (p *ReceiveFileClientStream) Send(in ReceiveFileRequestInput) error {
	return p.s.Send(NewReceiveFileRequest(in))
}

// CloseSend half-closes the outbound half; the server's Recv then sees io.EOF.
func (p *ReceiveFileClientStream) CloseSend() error { return p.s.CloseSend() }

// Reply reads the single terminal ReceiveFileResponse the server sends after it
// drains the request frames.
func (p *ReceiveFileClientStream) Reply() (ReceiveFileResponse, error) {
	b, err := p.s.Recv()
	if err != nil {
		return ReceiveFileResponse{}, err
	}
	return WrapReceiveFileResponse(b)
}

// ReceiveFile opens the client-stream. The first frame (carrying ReceiveFileInfo)
// rides as init; subsequent content frames go via Send.
func (c *Client) ReceiveFile(first ReceiveFileRequestInput) (*ReceiveFileClientStream, error) {
	s, err := c.conn.OpenStream(VolumeServerReceiveFileOrdinal, NewReceiveFileRequest(first))
	if err != nil {
		return nil, err
	}
	return &ReceiveFileClientStream{s: s}, nil
}
