// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// Native ZAP service adapter for the HanzoFiler service — the hand-written glue
// that binds the generated zero-copy wire (filer_zap.go + the *_zap.go message
// views/builders in this package) to a real github.com/zap-proto/go transport.
// It is to s3/wire/filer what s3/zapsvc is to s3/wire/object: the ONLY
// hand-written seam. No gRPC, no protobuf framing, no struct marshaling — every
// frame that crosses the socket is a filerwire New*/Wrap* buffer carried as a
// ZAP RPC envelope and correlated by PromiseID.
//
// SCOPE: this file covers the 28 UNARY HanzoFiler RPCs. The 6 STREAMING RPCs
// (ListEntries, StreamRenameEntry, StreamMutateEntry, TraverseBfsMetadata,
// SubscribeMetadata, SubscribeLocalMetadata) are NOT served here — the unary
// path refuses them with ErrStreamingEndpoint rather than faking a single-frame
// reply. Real per-frame streaming for those lives in the sibling
// s3/wire/filer/filerstream package, wired over transport.ListenStream /
// Conn.OpenStream (the duplex primitive shipped by the transport v1.5.0+). The
// streaming ordinals stay reserved in filer_zap.go so the wire numbering is
// stable across both paths.
//
// The Backend below is defined in filer-engine terms (zero-copy request views
// in, built response buffers out) and does NOT import the real
// filer/master/volume backend — the integrator implements it against the live
// engine; tests use an in-memory stub (see service_adapter_roundtrip_test.go).
// Modeling the seam as (view -> []byte) keeps it orthogonal: the adapter owns
// Wrap/Dispatch/transport, the Backend owns engine logic and builds replies
// with the New* constructors. The bytes ARE the message at every layer.

package filerwire

import (
	"errors"

	"github.com/zap-proto/go/transport"
)

// ErrStreamingEndpoint is returned by the unary handler for the 6 streaming
// HanzoFiler RPCs. They are served over the duplex stream primitive by the
// s3/wire/filer/filerstream package, NOT over the unary Call/Response path, so
// invoking them unary is a programming error, not a silently-buffered stream.
var ErrStreamingEndpoint = errors.New("filerwire: streaming RPC must use the s3/wire/filer/filerstream endpoint, not the unary path")

// Backend is the server-side contract for the 28 unary HanzoFiler RPCs. The
// real S3 filer engine implements it; tests use an in-memory stub. Each method
// receives the request as its zero-copy wire view and returns a response buffer
// built with the matching New*Response constructor (the bytes ARE the message).
// A returned error maps to an rpc.StatusInternal response on the wire.
//
// Streaming RPCs are intentionally absent — implement filerstream.MetadataSource
// / filerstream.Mutator / filerstream.Renamer / filerstream.Lister for those.
type Backend interface {
	LookupDirectoryEntry(req LookupDirectoryEntryRequest) ([]byte, error)
	CreateEntry(req CreateEntryRequest) ([]byte, error)
	UpdateEntry(req UpdateEntryRequest) ([]byte, error)
	TouchAccessTime(req TouchAccessTimeRequest) ([]byte, error)
	AppendToEntry(req AppendToEntryRequest) ([]byte, error)
	DeleteEntry(req DeleteEntryRequest) ([]byte, error)
	ObjectTransaction(req ObjectTransactionRequest) ([]byte, error)
	ObjectTransactionBatch(req ObjectTransactionBatchRequest) ([]byte, error)
	PosixLock(req PosixLockRequest) ([]byte, error)
	AtomicRenameEntry(req AtomicRenameEntryRequest) ([]byte, error)
	AssignVolume(req AssignVolumeRequest) ([]byte, error)
	LookupVolume(req LookupVolumeRequest) ([]byte, error)
	CollectionList(req CollectionListRequest) ([]byte, error)
	DeleteCollection(req DeleteCollectionRequest) ([]byte, error)
	Statistics(req StatisticsRequest) ([]byte, error)
	Ping(req PingRequest) ([]byte, error)
	GetFilerConfiguration(req GetFilerConfigurationRequest) ([]byte, error)
	ListMetadataSubscribers(req ListMetadataSubscribersRequest) ([]byte, error)
	KvGet(req KvGetRequest) ([]byte, error)
	KvPut(req KvPutRequest) ([]byte, error)
	CacheRemoteObjectToLocalCluster(req CacheRemoteObjectToLocalClusterRequest) ([]byte, error)
	DistributedLock(req LockRequest) ([]byte, error)
	DistributedUnlock(req UnlockRequest) ([]byte, error)
	FindLockOwner(req FindLockOwnerRequest) ([]byte, error)
	TransferLocks(req TransferLocksRequest) ([]byte, error)
	ReplicateLock(req ReplicateLockRequest) ([]byte, error)
	MountRegister(req MountRegisterRequest) ([]byte, error)
	MountList(req MountListRequest) ([]byte, error)
}

// handler adapts a Backend to the generated HanzoFilerHandler: it Wraps each
// request buffer into its zero-copy view and delegates to the Backend, which
// returns the already-built response buffer. The 6 streaming methods refuse the
// unary path with ErrStreamingEndpoint (DispatchHanzoFiler turns the error into
// a StatusInternal response; the real stream is served by filerstream).
type handler struct{ b Backend }

func (h handler) LookupDirectoryEntry(req []byte) ([]byte, error) {
	v, err := WrapLookupDirectoryEntryRequest(req)
	if err != nil {
		return nil, err
	}
	return h.b.LookupDirectoryEntry(v)
}

func (h handler) CreateEntry(req []byte) ([]byte, error) {
	v, err := WrapCreateEntryRequest(req)
	if err != nil {
		return nil, err
	}
	return h.b.CreateEntry(v)
}

func (h handler) UpdateEntry(req []byte) ([]byte, error) {
	v, err := WrapUpdateEntryRequest(req)
	if err != nil {
		return nil, err
	}
	return h.b.UpdateEntry(v)
}

func (h handler) TouchAccessTime(req []byte) ([]byte, error) {
	v, err := WrapTouchAccessTimeRequest(req)
	if err != nil {
		return nil, err
	}
	return h.b.TouchAccessTime(v)
}

func (h handler) AppendToEntry(req []byte) ([]byte, error) {
	v, err := WrapAppendToEntryRequest(req)
	if err != nil {
		return nil, err
	}
	return h.b.AppendToEntry(v)
}

func (h handler) DeleteEntry(req []byte) ([]byte, error) {
	v, err := WrapDeleteEntryRequest(req)
	if err != nil {
		return nil, err
	}
	return h.b.DeleteEntry(v)
}

func (h handler) ObjectTransaction(req []byte) ([]byte, error) {
	v, err := WrapObjectTransactionRequest(req)
	if err != nil {
		return nil, err
	}
	return h.b.ObjectTransaction(v)
}

func (h handler) ObjectTransactionBatch(req []byte) ([]byte, error) {
	v, err := WrapObjectTransactionBatchRequest(req)
	if err != nil {
		return nil, err
	}
	return h.b.ObjectTransactionBatch(v)
}

func (h handler) PosixLock(req []byte) ([]byte, error) {
	v, err := WrapPosixLockRequest(req)
	if err != nil {
		return nil, err
	}
	return h.b.PosixLock(v)
}

func (h handler) AtomicRenameEntry(req []byte) ([]byte, error) {
	v, err := WrapAtomicRenameEntryRequest(req)
	if err != nil {
		return nil, err
	}
	return h.b.AtomicRenameEntry(v)
}

func (h handler) AssignVolume(req []byte) ([]byte, error) {
	v, err := WrapAssignVolumeRequest(req)
	if err != nil {
		return nil, err
	}
	return h.b.AssignVolume(v)
}

func (h handler) LookupVolume(req []byte) ([]byte, error) {
	v, err := WrapLookupVolumeRequest(req)
	if err != nil {
		return nil, err
	}
	return h.b.LookupVolume(v)
}

func (h handler) CollectionList(req []byte) ([]byte, error) {
	v, err := WrapCollectionListRequest(req)
	if err != nil {
		return nil, err
	}
	return h.b.CollectionList(v)
}

func (h handler) DeleteCollection(req []byte) ([]byte, error) {
	v, err := WrapDeleteCollectionRequest(req)
	if err != nil {
		return nil, err
	}
	return h.b.DeleteCollection(v)
}

func (h handler) Statistics(req []byte) ([]byte, error) {
	v, err := WrapStatisticsRequest(req)
	if err != nil {
		return nil, err
	}
	return h.b.Statistics(v)
}

func (h handler) Ping(req []byte) ([]byte, error) {
	v, err := WrapPingRequest(req)
	if err != nil {
		return nil, err
	}
	return h.b.Ping(v)
}

func (h handler) GetFilerConfiguration(req []byte) ([]byte, error) {
	v, err := WrapGetFilerConfigurationRequest(req)
	if err != nil {
		return nil, err
	}
	return h.b.GetFilerConfiguration(v)
}

func (h handler) ListMetadataSubscribers(req []byte) ([]byte, error) {
	v, err := WrapListMetadataSubscribersRequest(req)
	if err != nil {
		return nil, err
	}
	return h.b.ListMetadataSubscribers(v)
}

func (h handler) KvGet(req []byte) ([]byte, error) {
	v, err := WrapKvGetRequest(req)
	if err != nil {
		return nil, err
	}
	return h.b.KvGet(v)
}

func (h handler) KvPut(req []byte) ([]byte, error) {
	v, err := WrapKvPutRequest(req)
	if err != nil {
		return nil, err
	}
	return h.b.KvPut(v)
}

func (h handler) CacheRemoteObjectToLocalCluster(req []byte) ([]byte, error) {
	v, err := WrapCacheRemoteObjectToLocalClusterRequest(req)
	if err != nil {
		return nil, err
	}
	return h.b.CacheRemoteObjectToLocalCluster(v)
}

func (h handler) DistributedLock(req []byte) ([]byte, error) {
	v, err := WrapLockRequest(req)
	if err != nil {
		return nil, err
	}
	return h.b.DistributedLock(v)
}

func (h handler) DistributedUnlock(req []byte) ([]byte, error) {
	v, err := WrapUnlockRequest(req)
	if err != nil {
		return nil, err
	}
	return h.b.DistributedUnlock(v)
}

func (h handler) FindLockOwner(req []byte) ([]byte, error) {
	v, err := WrapFindLockOwnerRequest(req)
	if err != nil {
		return nil, err
	}
	return h.b.FindLockOwner(v)
}

func (h handler) TransferLocks(req []byte) ([]byte, error) {
	v, err := WrapTransferLocksRequest(req)
	if err != nil {
		return nil, err
	}
	return h.b.TransferLocks(v)
}

func (h handler) ReplicateLock(req []byte) ([]byte, error) {
	v, err := WrapReplicateLockRequest(req)
	if err != nil {
		return nil, err
	}
	return h.b.ReplicateLock(v)
}

func (h handler) MountRegister(req []byte) ([]byte, error) {
	v, err := WrapMountRegisterRequest(req)
	if err != nil {
		return nil, err
	}
	return h.b.MountRegister(v)
}

func (h handler) MountList(req []byte) ([]byte, error) {
	v, err := WrapMountListRequest(req)
	if err != nil {
		return nil, err
	}
	return h.b.MountList(v)
}

// --- streaming RPCs: refused on the unary path (served by filerstream) ---

func (h handler) ListEntries(req []byte) ([]byte, error)             { return nil, ErrStreamingEndpoint }
func (h handler) StreamRenameEntry(req []byte) ([]byte, error)       { return nil, ErrStreamingEndpoint }
func (h handler) StreamMutateEntry(req []byte) ([]byte, error)       { return nil, ErrStreamingEndpoint }
func (h handler) TraverseBfsMetadata(req []byte) ([]byte, error)     { return nil, ErrStreamingEndpoint }
func (h handler) SubscribeMetadata(req []byte) ([]byte, error)       { return nil, ErrStreamingEndpoint }
func (h handler) SubscribeLocalMetadata(req []byte) ([]byte, error)  { return nil, ErrStreamingEndpoint }

// compile-time assertion that handler satisfies the generated server contract.
var _ HanzoFilerHandler = handler{}

// Dispatch is the HanzoFiler server dispatch bound to b — pass it to
// transport.Listen (or use Serve). Exposed so callers can compose it (e.g.
// behind PQ-TLS, or alongside the filerstream StreamHandler via
// transport.ListenStream).
func Dispatch(b Backend) transport.Dispatch {
	h := handler{b: b}
	return func(envelope []byte) ([]byte, error) {
		return DispatchHanzoFiler(h, envelope)
	}
}

// Serve starts the native ZAP HanzoFiler service on network/addr (e.g. "tcp",
// ":18888" or "unix", sock), backed by b. Returns the running server; Close
// stops it. Unary-only: to also serve the streaming RPCs on the same listener,
// use transport.ListenStream(network, addr, Dispatch(b), filerstream.Handler(...)).
func Serve(network, addr string, b Backend) (*transport.Server, error) {
	return transport.Listen(network, addr, Dispatch(b))
}

// Client is the typed HanzoFiler ZAP service client internal services hold. It
// owns the transport connection to one filer endpoint and issues the unary
// RPCs; the response is returned as its zero-copy wire view. For streaming use
// the filerstream client over the same (or a sibling) connection.
type Client struct {
	conn transport.Conn
	rpc  *HanzoFilerClient
}

// Dial connects to the HanzoFiler ZAP service at addr (e.g. "filer.hanzo.svc:18888")
// over plain TCP. For the PQ-secured mesh establish the transport.Conn via the
// TLS/QUIC path and use NewClient.
func Dial(network, addr string) (*Client, error) {
	conn, err := transport.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn, rpc: NewHanzoFilerClient(conn, nil)}, nil
}

// NewClient wraps an already-established transport.Conn (TCP, Unix, or PQ-TLS).
func NewClient(conn transport.Conn) *Client {
	return &Client{conn: conn, rpc: NewHanzoFilerClient(conn, nil)}
}

// Conn exposes the underlying transport connection so a filerstream client can
// open streams over the same socket as the unary client.
func (c *Client) Conn() transport.Conn { return c.conn }

// Close releases the underlying connection.
func (c *Client) Close() error { return c.conn.Close() }

// The typed unary calls below build the request with New*, ship it over the
// generated HanzoFilerClient, and return the response as its zero-copy view.

func (c *Client) LookupDirectoryEntry(req []byte) (LookupDirectoryEntryResponse, error) {
	_, body, err := c.rpc.LookupDirectoryEntry(req)
	if err != nil {
		return LookupDirectoryEntryResponse{}, err
	}
	return WrapLookupDirectoryEntryResponse(body)
}

func (c *Client) CreateEntry(req []byte) (CreateEntryResponse, error) {
	_, body, err := c.rpc.CreateEntry(req)
	if err != nil {
		return CreateEntryResponse{}, err
	}
	return WrapCreateEntryResponse(body)
}

func (c *Client) UpdateEntry(req []byte) (UpdateEntryResponse, error) {
	_, body, err := c.rpc.UpdateEntry(req)
	if err != nil {
		return UpdateEntryResponse{}, err
	}
	return WrapUpdateEntryResponse(body)
}

func (c *Client) TouchAccessTime(req []byte) (TouchAccessTimeResponse, error) {
	_, body, err := c.rpc.TouchAccessTime(req)
	if err != nil {
		return TouchAccessTimeResponse{}, err
	}
	return WrapTouchAccessTimeResponse(body)
}

func (c *Client) AppendToEntry(req []byte) (AppendToEntryResponse, error) {
	_, body, err := c.rpc.AppendToEntry(req)
	if err != nil {
		return AppendToEntryResponse{}, err
	}
	return WrapAppendToEntryResponse(body)
}

func (c *Client) DeleteEntry(req []byte) (DeleteEntryResponse, error) {
	_, body, err := c.rpc.DeleteEntry(req)
	if err != nil {
		return DeleteEntryResponse{}, err
	}
	return WrapDeleteEntryResponse(body)
}

func (c *Client) ObjectTransaction(req []byte) (ObjectTransactionResponse, error) {
	_, body, err := c.rpc.ObjectTransaction(req)
	if err != nil {
		return ObjectTransactionResponse{}, err
	}
	return WrapObjectTransactionResponse(body)
}

func (c *Client) ObjectTransactionBatch(req []byte) (ObjectTransactionBatchResponse, error) {
	_, body, err := c.rpc.ObjectTransactionBatch(req)
	if err != nil {
		return ObjectTransactionBatchResponse{}, err
	}
	return WrapObjectTransactionBatchResponse(body)
}

func (c *Client) PosixLock(req []byte) (PosixLockResponse, error) {
	_, body, err := c.rpc.PosixLock(req)
	if err != nil {
		return PosixLockResponse{}, err
	}
	return WrapPosixLockResponse(body)
}

func (c *Client) AtomicRenameEntry(req []byte) (AtomicRenameEntryResponse, error) {
	_, body, err := c.rpc.AtomicRenameEntry(req)
	if err != nil {
		return AtomicRenameEntryResponse{}, err
	}
	return WrapAtomicRenameEntryResponse(body)
}

func (c *Client) AssignVolume(req []byte) (AssignVolumeResponse, error) {
	_, body, err := c.rpc.AssignVolume(req)
	if err != nil {
		return AssignVolumeResponse{}, err
	}
	return WrapAssignVolumeResponse(body)
}

func (c *Client) LookupVolume(req []byte) (LookupVolumeResponse, error) {
	_, body, err := c.rpc.LookupVolume(req)
	if err != nil {
		return LookupVolumeResponse{}, err
	}
	return WrapLookupVolumeResponse(body)
}

func (c *Client) CollectionList(req []byte) (CollectionListResponse, error) {
	_, body, err := c.rpc.CollectionList(req)
	if err != nil {
		return CollectionListResponse{}, err
	}
	return WrapCollectionListResponse(body)
}

func (c *Client) DeleteCollection(req []byte) (DeleteCollectionResponse, error) {
	_, body, err := c.rpc.DeleteCollection(req)
	if err != nil {
		return DeleteCollectionResponse{}, err
	}
	return WrapDeleteCollectionResponse(body)
}

func (c *Client) Statistics(req []byte) (StatisticsResponse, error) {
	_, body, err := c.rpc.Statistics(req)
	if err != nil {
		return StatisticsResponse{}, err
	}
	return WrapStatisticsResponse(body)
}

func (c *Client) Ping(req []byte) (PingResponse, error) {
	_, body, err := c.rpc.Ping(req)
	if err != nil {
		return PingResponse{}, err
	}
	return WrapPingResponse(body)
}

func (c *Client) GetFilerConfiguration(req []byte) (GetFilerConfigurationResponse, error) {
	_, body, err := c.rpc.GetFilerConfiguration(req)
	if err != nil {
		return GetFilerConfigurationResponse{}, err
	}
	return WrapGetFilerConfigurationResponse(body)
}

func (c *Client) ListMetadataSubscribers(req []byte) (ListMetadataSubscribersResponse, error) {
	_, body, err := c.rpc.ListMetadataSubscribers(req)
	if err != nil {
		return ListMetadataSubscribersResponse{}, err
	}
	return WrapListMetadataSubscribersResponse(body)
}

func (c *Client) KvGet(req []byte) (KvGetResponse, error) {
	_, body, err := c.rpc.KvGet(req)
	if err != nil {
		return KvGetResponse{}, err
	}
	return WrapKvGetResponse(body)
}

func (c *Client) KvPut(req []byte) (KvPutResponse, error) {
	_, body, err := c.rpc.KvPut(req)
	if err != nil {
		return KvPutResponse{}, err
	}
	return WrapKvPutResponse(body)
}

func (c *Client) CacheRemoteObjectToLocalCluster(req []byte) (CacheRemoteObjectToLocalClusterResponse, error) {
	_, body, err := c.rpc.CacheRemoteObjectToLocalCluster(req)
	if err != nil {
		return CacheRemoteObjectToLocalClusterResponse{}, err
	}
	return WrapCacheRemoteObjectToLocalClusterResponse(body)
}

func (c *Client) DistributedLock(req []byte) (LockResponse, error) {
	_, body, err := c.rpc.DistributedLock(req)
	if err != nil {
		return LockResponse{}, err
	}
	return WrapLockResponse(body)
}

func (c *Client) DistributedUnlock(req []byte) (UnlockResponse, error) {
	_, body, err := c.rpc.DistributedUnlock(req)
	if err != nil {
		return UnlockResponse{}, err
	}
	return WrapUnlockResponse(body)
}

func (c *Client) FindLockOwner(req []byte) (FindLockOwnerResponse, error) {
	_, body, err := c.rpc.FindLockOwner(req)
	if err != nil {
		return FindLockOwnerResponse{}, err
	}
	return WrapFindLockOwnerResponse(body)
}

func (c *Client) TransferLocks(req []byte) (TransferLocksResponse, error) {
	_, body, err := c.rpc.TransferLocks(req)
	if err != nil {
		return TransferLocksResponse{}, err
	}
	return WrapTransferLocksResponse(body)
}

func (c *Client) ReplicateLock(req []byte) (ReplicateLockResponse, error) {
	_, body, err := c.rpc.ReplicateLock(req)
	if err != nil {
		return ReplicateLockResponse{}, err
	}
	return WrapReplicateLockResponse(body)
}

func (c *Client) MountRegister(req []byte) (MountRegisterResponse, error) {
	_, body, err := c.rpc.MountRegister(req)
	if err != nil {
		return MountRegisterResponse{}, err
	}
	return WrapMountRegisterResponse(body)
}

func (c *Client) MountList(req []byte) (MountListResponse, error) {
	_, body, err := c.rpc.MountList(req)
	if err != nil {
		return MountListResponse{}, err
	}
	return WrapMountListResponse(body)
}
