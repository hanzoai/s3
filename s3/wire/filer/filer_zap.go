// Code generated from filer.proto; DO NOT EDIT.

package filerwire

import (
	"context"
	"fmt"

	"github.com/zap-proto/go/rpc"
)

// Method ordinals for the HanzoFiler service (stable 1-based wire ids, in
// filer.proto declaration order). Streaming methods keep an ordinal slot so the
// wire numbering is stable once the transport streaming primitive ships.
const (
	HanzoFilerLookupDirectoryEntryOrdinal            uint32 = 1
	HanzoFilerListEntriesOrdinal                     uint32 = 2 // STREAMING (server)
	HanzoFilerCreateEntryOrdinal                     uint32 = 3
	HanzoFilerUpdateEntryOrdinal                     uint32 = 4
	HanzoFilerTouchAccessTimeOrdinal                 uint32 = 5
	HanzoFilerAppendToEntryOrdinal                   uint32 = 6
	HanzoFilerDeleteEntryOrdinal                     uint32 = 7
	HanzoFilerObjectTransactionOrdinal               uint32 = 8
	HanzoFilerObjectTransactionBatchOrdinal          uint32 = 9
	HanzoFilerPosixLockOrdinal                       uint32 = 10
	HanzoFilerAtomicRenameEntryOrdinal               uint32 = 11
	HanzoFilerStreamRenameEntryOrdinal               uint32 = 12 // STREAMING (server)
	HanzoFilerStreamMutateEntryOrdinal               uint32 = 13 // STREAMING (bidirectional)
	HanzoFilerAssignVolumeOrdinal                    uint32 = 14
	HanzoFilerLookupVolumeOrdinal                    uint32 = 15
	HanzoFilerCollectionListOrdinal                  uint32 = 16
	HanzoFilerDeleteCollectionOrdinal                uint32 = 17
	HanzoFilerStatisticsOrdinal                      uint32 = 18
	HanzoFilerPingOrdinal                            uint32 = 19
	HanzoFilerGetFilerConfigurationOrdinal           uint32 = 20
	HanzoFilerTraverseBfsMetadataOrdinal             uint32 = 21 // STREAMING (server)
	HanzoFilerSubscribeMetadataOrdinal               uint32 = 22 // STREAMING (server)
	HanzoFilerSubscribeLocalMetadataOrdinal          uint32 = 23 // STREAMING (server)
	HanzoFilerListMetadataSubscribersOrdinal         uint32 = 24
	HanzoFilerKvGetOrdinal                           uint32 = 25
	HanzoFilerKvPutOrdinal                           uint32 = 26
	HanzoFilerCacheRemoteObjectToLocalClusterOrdinal uint32 = 27
	HanzoFilerDistributedLockOrdinal                 uint32 = 28
	HanzoFilerDistributedUnlockOrdinal               uint32 = 29
	HanzoFilerFindLockOwnerOrdinal                   uint32 = 30
	HanzoFilerTransferLocksOrdinal                   uint32 = 31
	HanzoFilerReplicateLockOrdinal                   uint32 = 32
	HanzoFilerMountRegisterOrdinal                   uint32 = 33
	HanzoFilerMountListOrdinal                       uint32 = 34
)

// HanzoFilerChannel ships one Call envelope and awaits its correlated Response.
// CallContext is Call that also aborts when ctx is done (transport.Conn
// satisfies both).
type HanzoFilerChannel interface {
	Call(envelope []byte) (rpc.Response, error)
	CallContext(ctx context.Context, envelope []byte) (rpc.Response, error)
}

// HanzoFilerClient is a typed RPC client for the HanzoFiler service over a ZAP
// call channel. Each call takes a fresh PromiseID from sess; the pipelined "On"
// form of a method sets Target to a prior call's Promise so the server chains
// them. Request/response payloads cross the transport as opaque ZAP buffers;
// build and read them with the message constructors in this package.
type HanzoFilerClient struct {
	ch   HanzoFilerChannel
	cap  []byte
	sess *rpc.Session
}

// NewHanzoFilerClient returns a client that issues calls over ch, attaching cap
// (which may be nil) to every request.
func NewHanzoFilerClient(ch HanzoFilerChannel, capability []byte) *HanzoFilerClient {
	return &HanzoFilerClient{ch: ch, cap: capability, sess: rpc.NewSession()}
}

// unaryCall issues one request envelope under a fresh promise and returns the
// response body, mapping a non-OK status to an error tagged with method. The
// call aborts when ctx is done.
func (c *HanzoFilerClient) unaryCall(ctx context.Context, method, target uint32, payload []byte, name string) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.CallContext(ctx, rpc.BuildRequest(rpc.Call{
		Method:    method,
		PromiseID: p.ID,
		Target:    target,
		Cap:       c.cap,
		Payload:   payload,
	}))
	if err != nil {
		return p, nil, err
	}
	if resp.Status != rpc.StatusOK {
		if len(resp.Body) > 0 {
			// The server carries the handler error message in the body; surface it
			// so callers can detect sentinels (e.g. filer_pb.ErrNotFound).
			return p, nil, fmt.Errorf("HanzoFiler.%s: %s", name, resp.Body)
		}
		return p, nil, fmt.Errorf("HanzoFiler.%s: status %d", name, resp.Status)
	}
	return p, resp.Body, nil
}

func (c *HanzoFilerClient) LookupDirectoryEntry(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerLookupDirectoryEntryOrdinal, rpc.NoTarget, req, "LookupDirectoryEntry")
}

// LookupDirectoryEntryOn issues LookupDirectoryEntry as a dependent call
// pipelined on the answer of on.
func (c *HanzoFilerClient) LookupDirectoryEntryOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerLookupDirectoryEntryOrdinal, on.ID, nil, "LookupDirectoryEntry")
}

func (c *HanzoFilerClient) CreateEntry(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerCreateEntryOrdinal, rpc.NoTarget, req, "CreateEntry")
}

// CreateEntryOn issues CreateEntry as a dependent call pipelined on on.
func (c *HanzoFilerClient) CreateEntryOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerCreateEntryOrdinal, on.ID, nil, "CreateEntry")
}

func (c *HanzoFilerClient) UpdateEntry(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerUpdateEntryOrdinal, rpc.NoTarget, req, "UpdateEntry")
}

// UpdateEntryOn issues UpdateEntry as a dependent call pipelined on on.
func (c *HanzoFilerClient) UpdateEntryOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerUpdateEntryOrdinal, on.ID, nil, "UpdateEntry")
}

func (c *HanzoFilerClient) TouchAccessTime(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerTouchAccessTimeOrdinal, rpc.NoTarget, req, "TouchAccessTime")
}

// TouchAccessTimeOn issues TouchAccessTime as a dependent call pipelined on on.
func (c *HanzoFilerClient) TouchAccessTimeOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerTouchAccessTimeOrdinal, on.ID, nil, "TouchAccessTime")
}

func (c *HanzoFilerClient) AppendToEntry(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerAppendToEntryOrdinal, rpc.NoTarget, req, "AppendToEntry")
}

// AppendToEntryOn issues AppendToEntry as a dependent call pipelined on on.
func (c *HanzoFilerClient) AppendToEntryOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerAppendToEntryOrdinal, on.ID, nil, "AppendToEntry")
}

func (c *HanzoFilerClient) DeleteEntry(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerDeleteEntryOrdinal, rpc.NoTarget, req, "DeleteEntry")
}

// DeleteEntryOn issues DeleteEntry as a dependent call pipelined on on.
func (c *HanzoFilerClient) DeleteEntryOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerDeleteEntryOrdinal, on.ID, nil, "DeleteEntry")
}

func (c *HanzoFilerClient) ObjectTransaction(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerObjectTransactionOrdinal, rpc.NoTarget, req, "ObjectTransaction")
}

// ObjectTransactionOn issues ObjectTransaction as a dependent call pipelined on on.
func (c *HanzoFilerClient) ObjectTransactionOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerObjectTransactionOrdinal, on.ID, nil, "ObjectTransaction")
}

func (c *HanzoFilerClient) ObjectTransactionBatch(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerObjectTransactionBatchOrdinal, rpc.NoTarget, req, "ObjectTransactionBatch")
}

// ObjectTransactionBatchOn issues ObjectTransactionBatch as a dependent call
// pipelined on on.
func (c *HanzoFilerClient) ObjectTransactionBatchOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerObjectTransactionBatchOrdinal, on.ID, nil, "ObjectTransactionBatch")
}

func (c *HanzoFilerClient) PosixLock(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerPosixLockOrdinal, rpc.NoTarget, req, "PosixLock")
}

// PosixLockOn issues PosixLock as a dependent call pipelined on on.
func (c *HanzoFilerClient) PosixLockOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerPosixLockOrdinal, on.ID, nil, "PosixLock")
}

func (c *HanzoFilerClient) AtomicRenameEntry(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerAtomicRenameEntryOrdinal, rpc.NoTarget, req, "AtomicRenameEntry")
}

// AtomicRenameEntryOn issues AtomicRenameEntry as a dependent call pipelined on on.
func (c *HanzoFilerClient) AtomicRenameEntryOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerAtomicRenameEntryOrdinal, on.ID, nil, "AtomicRenameEntry")
}

// StreamRenameEntry is a server-streaming RPC.
// STREAMING: the streaming body lands when the transport streaming primitive
// ships. Until then this unary stub issues a single request and returns the
// first response frame, so the ordinal and wire types are exercised end-to-end.
func (c *HanzoFilerClient) StreamRenameEntry(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerStreamRenameEntryOrdinal, rpc.NoTarget, req, "StreamRenameEntry")
}

// StreamMutateEntry is a bidirectional-streaming RPC.
// STREAMING: the streaming body lands when the transport streaming primitive
// ships. Until then this unary stub issues a single request and returns the
// first response frame, so the ordinal and wire types are exercised end-to-end.
func (c *HanzoFilerClient) StreamMutateEntry(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerStreamMutateEntryOrdinal, rpc.NoTarget, req, "StreamMutateEntry")
}

func (c *HanzoFilerClient) AssignVolume(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerAssignVolumeOrdinal, rpc.NoTarget, req, "AssignVolume")
}

// AssignVolumeOn issues AssignVolume as a dependent call pipelined on on.
func (c *HanzoFilerClient) AssignVolumeOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerAssignVolumeOrdinal, on.ID, nil, "AssignVolume")
}

func (c *HanzoFilerClient) LookupVolume(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerLookupVolumeOrdinal, rpc.NoTarget, req, "LookupVolume")
}

// LookupVolumeOn issues LookupVolume as a dependent call pipelined on on.
func (c *HanzoFilerClient) LookupVolumeOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerLookupVolumeOrdinal, on.ID, nil, "LookupVolume")
}

func (c *HanzoFilerClient) CollectionList(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerCollectionListOrdinal, rpc.NoTarget, req, "CollectionList")
}

// CollectionListOn issues CollectionList as a dependent call pipelined on on.
func (c *HanzoFilerClient) CollectionListOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerCollectionListOrdinal, on.ID, nil, "CollectionList")
}

func (c *HanzoFilerClient) DeleteCollection(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerDeleteCollectionOrdinal, rpc.NoTarget, req, "DeleteCollection")
}

// DeleteCollectionOn issues DeleteCollection as a dependent call pipelined on on.
func (c *HanzoFilerClient) DeleteCollectionOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerDeleteCollectionOrdinal, on.ID, nil, "DeleteCollection")
}

func (c *HanzoFilerClient) Statistics(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerStatisticsOrdinal, rpc.NoTarget, req, "Statistics")
}

// StatisticsOn issues Statistics as a dependent call pipelined on on.
func (c *HanzoFilerClient) StatisticsOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerStatisticsOrdinal, on.ID, nil, "Statistics")
}

func (c *HanzoFilerClient) Ping(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerPingOrdinal, rpc.NoTarget, req, "Ping")
}

// PingOn issues Ping as a dependent call pipelined on on.
func (c *HanzoFilerClient) PingOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerPingOrdinal, on.ID, nil, "Ping")
}

func (c *HanzoFilerClient) GetFilerConfiguration(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerGetFilerConfigurationOrdinal, rpc.NoTarget, req, "GetFilerConfiguration")
}

// GetFilerConfigurationOn issues GetFilerConfiguration as a dependent call
// pipelined on on.
func (c *HanzoFilerClient) GetFilerConfigurationOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerGetFilerConfigurationOrdinal, on.ID, nil, "GetFilerConfiguration")
}

// TraverseBfsMetadata is a server-streaming RPC.
// STREAMING: the streaming body lands when the transport streaming primitive
// ships. Until then this unary stub issues a single request and returns the
// first response frame, so the ordinal and wire types are exercised end-to-end.
func (c *HanzoFilerClient) TraverseBfsMetadata(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerTraverseBfsMetadataOrdinal, rpc.NoTarget, req, "TraverseBfsMetadata")
}

// SubscribeMetadata is a server-streaming RPC.
// STREAMING: the streaming body lands when the transport streaming primitive
// ships. Until then this unary stub issues a single request and returns the
// first response frame, so the ordinal and wire types are exercised end-to-end.
func (c *HanzoFilerClient) SubscribeMetadata(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerSubscribeMetadataOrdinal, rpc.NoTarget, req, "SubscribeMetadata")
}

// SubscribeLocalMetadata is a server-streaming RPC (shares the SubscribeMetadata
// request/response message types).
// STREAMING: the streaming body lands when the transport streaming primitive
// ships. Until then this unary stub issues a single request and returns the
// first response frame, so the ordinal and wire types are exercised end-to-end.
func (c *HanzoFilerClient) SubscribeLocalMetadata(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerSubscribeLocalMetadataOrdinal, rpc.NoTarget, req, "SubscribeLocalMetadata")
}

func (c *HanzoFilerClient) ListMetadataSubscribers(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerListMetadataSubscribersOrdinal, rpc.NoTarget, req, "ListMetadataSubscribers")
}

// ListMetadataSubscribersOn issues ListMetadataSubscribers as a dependent call
// pipelined on on.
func (c *HanzoFilerClient) ListMetadataSubscribersOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerListMetadataSubscribersOrdinal, on.ID, nil, "ListMetadataSubscribers")
}

func (c *HanzoFilerClient) KvGet(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerKvGetOrdinal, rpc.NoTarget, req, "KvGet")
}

// KvGetOn issues KvGet as a dependent call pipelined on on.
func (c *HanzoFilerClient) KvGetOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerKvGetOrdinal, on.ID, nil, "KvGet")
}

func (c *HanzoFilerClient) KvPut(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerKvPutOrdinal, rpc.NoTarget, req, "KvPut")
}

// KvPutOn issues KvPut as a dependent call pipelined on on.
func (c *HanzoFilerClient) KvPutOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerKvPutOrdinal, on.ID, nil, "KvPut")
}

func (c *HanzoFilerClient) CacheRemoteObjectToLocalCluster(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerCacheRemoteObjectToLocalClusterOrdinal, rpc.NoTarget, req, "CacheRemoteObjectToLocalCluster")
}

// CacheRemoteObjectToLocalClusterOn issues CacheRemoteObjectToLocalCluster as a
// dependent call pipelined on on.
func (c *HanzoFilerClient) CacheRemoteObjectToLocalClusterOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerCacheRemoteObjectToLocalClusterOrdinal, on.ID, nil, "CacheRemoteObjectToLocalCluster")
}

func (c *HanzoFilerClient) DistributedLock(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerDistributedLockOrdinal, rpc.NoTarget, req, "DistributedLock")
}

// DistributedLockOn issues DistributedLock as a dependent call pipelined on on.
func (c *HanzoFilerClient) DistributedLockOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerDistributedLockOrdinal, on.ID, nil, "DistributedLock")
}

func (c *HanzoFilerClient) DistributedUnlock(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerDistributedUnlockOrdinal, rpc.NoTarget, req, "DistributedUnlock")
}

// DistributedUnlockOn issues DistributedUnlock as a dependent call pipelined on on.
func (c *HanzoFilerClient) DistributedUnlockOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerDistributedUnlockOrdinal, on.ID, nil, "DistributedUnlock")
}

func (c *HanzoFilerClient) FindLockOwner(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerFindLockOwnerOrdinal, rpc.NoTarget, req, "FindLockOwner")
}

// FindLockOwnerOn issues FindLockOwner as a dependent call pipelined on on.
func (c *HanzoFilerClient) FindLockOwnerOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerFindLockOwnerOrdinal, on.ID, nil, "FindLockOwner")
}

func (c *HanzoFilerClient) TransferLocks(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerTransferLocksOrdinal, rpc.NoTarget, req, "TransferLocks")
}

// TransferLocksOn issues TransferLocks as a dependent call pipelined on on.
func (c *HanzoFilerClient) TransferLocksOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerTransferLocksOrdinal, on.ID, nil, "TransferLocks")
}

func (c *HanzoFilerClient) ReplicateLock(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerReplicateLockOrdinal, rpc.NoTarget, req, "ReplicateLock")
}

// ReplicateLockOn issues ReplicateLock as a dependent call pipelined on on.
func (c *HanzoFilerClient) ReplicateLockOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerReplicateLockOrdinal, on.ID, nil, "ReplicateLock")
}

func (c *HanzoFilerClient) MountRegister(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerMountRegisterOrdinal, rpc.NoTarget, req, "MountRegister")
}

// MountRegisterOn issues MountRegister as a dependent call pipelined on on.
func (c *HanzoFilerClient) MountRegisterOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerMountRegisterOrdinal, on.ID, nil, "MountRegister")
}

func (c *HanzoFilerClient) MountList(ctx context.Context, req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerMountListOrdinal, rpc.NoTarget, req, "MountList")
}

// MountListOn issues MountList as a dependent call pipelined on on.
func (c *HanzoFilerClient) MountListOn(ctx context.Context, on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(ctx, HanzoFilerMountListOrdinal, on.ID, nil, "MountList")
}

// HanzoFilerHandler is the server contract for the HanzoFiler service. Implement
// each method, then route requests to it with DispatchHanzoFiler. Each req/return
// value is an opaque filer wire buffer (this package); decode it with the
// matching Wrap* and build the reply with New*.
//
// STREAMING methods (ListEntries, StreamRenameEntry, StreamMutateEntry,
// TraverseBfsMetadata, SubscribeMetadata, SubscribeLocalMetadata) take a single
// request frame and return a single response frame in this skeleton; the full
// per-frame streaming contract lands when the transport streaming primitive ships.
type HanzoFilerHandler interface {
	LookupDirectoryEntry(req []byte) ([]byte, error)
	ListEntries(req []byte) ([]byte, error) // STREAMING (server)
	CreateEntry(req []byte) ([]byte, error)
	UpdateEntry(req []byte) ([]byte, error)
	TouchAccessTime(req []byte) ([]byte, error)
	AppendToEntry(req []byte) ([]byte, error)
	DeleteEntry(req []byte) ([]byte, error)
	ObjectTransaction(req []byte) ([]byte, error)
	ObjectTransactionBatch(req []byte) ([]byte, error)
	PosixLock(req []byte) ([]byte, error)
	AtomicRenameEntry(req []byte) ([]byte, error)
	StreamRenameEntry(req []byte) ([]byte, error) // STREAMING (server)
	StreamMutateEntry(req []byte) ([]byte, error) // STREAMING (bidirectional)
	AssignVolume(req []byte) ([]byte, error)
	LookupVolume(req []byte) ([]byte, error)
	CollectionList(req []byte) ([]byte, error)
	DeleteCollection(req []byte) ([]byte, error)
	Statistics(req []byte) ([]byte, error)
	Ping(req []byte) ([]byte, error)
	GetFilerConfiguration(req []byte) ([]byte, error)
	TraverseBfsMetadata(req []byte) ([]byte, error)    // STREAMING (server)
	SubscribeMetadata(req []byte) ([]byte, error)      // STREAMING (server)
	SubscribeLocalMetadata(req []byte) ([]byte, error) // STREAMING (server)
	ListMetadataSubscribers(req []byte) ([]byte, error)
	KvGet(req []byte) ([]byte, error)
	KvPut(req []byte) ([]byte, error)
	CacheRemoteObjectToLocalCluster(req []byte) ([]byte, error)
	DistributedLock(req []byte) ([]byte, error)
	DistributedUnlock(req []byte) ([]byte, error)
	FindLockOwner(req []byte) ([]byte, error)
	TransferLocks(req []byte) ([]byte, error)
	ReplicateLock(req []byte) ([]byte, error)
	MountRegister(req []byte) ([]byte, error)
	MountList(req []byte) ([]byte, error)
}

// DispatchHanzoFiler decodes a Call envelope, routes it by method ordinal to h,
// and returns the response envelope. An unknown ordinal yields a StatusNotFound
// response; a handler error yields StatusInternal.
func DispatchHanzoFiler(h HanzoFilerHandler, envelope []byte) ([]byte, error) {
	call, err := rpc.ParseRequest(envelope)
	if err != nil {
		return nil, err
	}
	reply := func(body []byte, herr error) ([]byte, error) {
		if herr != nil {
			// Carry the handler error message in the response body so the caller
			// can reconstruct it (e.g. filer_pb.ErrNotFound detection by string).
			// The rpc envelope itself carries only a status code.
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, []byte(herr.Error())), nil
		}
		return rpc.BuildResponse(rpc.StatusOK, call.PromiseID, body), nil
	}
	switch call.Method {
	case HanzoFilerLookupDirectoryEntryOrdinal:
		return reply(h.LookupDirectoryEntry(call.Payload))
	case HanzoFilerListEntriesOrdinal: // STREAMING (server)
		return reply(h.ListEntries(call.Payload))
	case HanzoFilerCreateEntryOrdinal:
		return reply(h.CreateEntry(call.Payload))
	case HanzoFilerUpdateEntryOrdinal:
		return reply(h.UpdateEntry(call.Payload))
	case HanzoFilerTouchAccessTimeOrdinal:
		return reply(h.TouchAccessTime(call.Payload))
	case HanzoFilerAppendToEntryOrdinal:
		return reply(h.AppendToEntry(call.Payload))
	case HanzoFilerDeleteEntryOrdinal:
		return reply(h.DeleteEntry(call.Payload))
	case HanzoFilerObjectTransactionOrdinal:
		return reply(h.ObjectTransaction(call.Payload))
	case HanzoFilerObjectTransactionBatchOrdinal:
		return reply(h.ObjectTransactionBatch(call.Payload))
	case HanzoFilerPosixLockOrdinal:
		return reply(h.PosixLock(call.Payload))
	case HanzoFilerAtomicRenameEntryOrdinal:
		return reply(h.AtomicRenameEntry(call.Payload))
	case HanzoFilerStreamRenameEntryOrdinal: // STREAMING (server)
		return reply(h.StreamRenameEntry(call.Payload))
	case HanzoFilerStreamMutateEntryOrdinal: // STREAMING (bidirectional)
		return reply(h.StreamMutateEntry(call.Payload))
	case HanzoFilerAssignVolumeOrdinal:
		return reply(h.AssignVolume(call.Payload))
	case HanzoFilerLookupVolumeOrdinal:
		return reply(h.LookupVolume(call.Payload))
	case HanzoFilerCollectionListOrdinal:
		return reply(h.CollectionList(call.Payload))
	case HanzoFilerDeleteCollectionOrdinal:
		return reply(h.DeleteCollection(call.Payload))
	case HanzoFilerStatisticsOrdinal:
		return reply(h.Statistics(call.Payload))
	case HanzoFilerPingOrdinal:
		return reply(h.Ping(call.Payload))
	case HanzoFilerGetFilerConfigurationOrdinal:
		return reply(h.GetFilerConfiguration(call.Payload))
	case HanzoFilerTraverseBfsMetadataOrdinal: // STREAMING (server)
		return reply(h.TraverseBfsMetadata(call.Payload))
	case HanzoFilerSubscribeMetadataOrdinal: // STREAMING (server)
		return reply(h.SubscribeMetadata(call.Payload))
	case HanzoFilerSubscribeLocalMetadataOrdinal: // STREAMING (server)
		return reply(h.SubscribeLocalMetadata(call.Payload))
	case HanzoFilerListMetadataSubscribersOrdinal:
		return reply(h.ListMetadataSubscribers(call.Payload))
	case HanzoFilerKvGetOrdinal:
		return reply(h.KvGet(call.Payload))
	case HanzoFilerKvPutOrdinal:
		return reply(h.KvPut(call.Payload))
	case HanzoFilerCacheRemoteObjectToLocalClusterOrdinal:
		return reply(h.CacheRemoteObjectToLocalCluster(call.Payload))
	case HanzoFilerDistributedLockOrdinal:
		return reply(h.DistributedLock(call.Payload))
	case HanzoFilerDistributedUnlockOrdinal:
		return reply(h.DistributedUnlock(call.Payload))
	case HanzoFilerFindLockOwnerOrdinal:
		return reply(h.FindLockOwner(call.Payload))
	case HanzoFilerTransferLocksOrdinal:
		return reply(h.TransferLocks(call.Payload))
	case HanzoFilerReplicateLockOrdinal:
		return reply(h.ReplicateLock(call.Payload))
	case HanzoFilerMountRegisterOrdinal:
		return reply(h.MountRegister(call.Payload))
	case HanzoFilerMountListOrdinal:
		return reply(h.MountList(call.Payload))
	default:
		return rpc.BuildResponse(rpc.StatusNotFound, call.PromiseID, nil), nil
	}
}
