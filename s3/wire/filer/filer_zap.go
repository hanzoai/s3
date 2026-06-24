// Code generated from filer.proto; DO NOT EDIT.

package filerwire

import (
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
type HanzoFilerChannel interface {
	Call(envelope []byte) (rpc.Response, error)
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
// response body, mapping a non-OK status to an error tagged with method.
func (c *HanzoFilerClient) unaryCall(method, target uint32, payload []byte, name string) (rpc.Promise, []byte, error) {
	p := c.sess.Next()
	resp, err := c.ch.Call(rpc.BuildRequest(rpc.Call{
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
		return p, nil, fmt.Errorf("HanzoFiler.%s: status %d", name, resp.Status)
	}
	return p, resp.Body, nil
}

func (c *HanzoFilerClient) LookupDirectoryEntry(req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerLookupDirectoryEntryOrdinal, rpc.NoTarget, req, "LookupDirectoryEntry")
}

// LookupDirectoryEntryOn issues LookupDirectoryEntry as a dependent call
// pipelined on the answer of on.
func (c *HanzoFilerClient) LookupDirectoryEntryOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerLookupDirectoryEntryOrdinal, on.ID, nil, "LookupDirectoryEntry")
}

func (c *HanzoFilerClient) CreateEntry(req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerCreateEntryOrdinal, rpc.NoTarget, req, "CreateEntry")
}

// CreateEntryOn issues CreateEntry as a dependent call pipelined on on.
func (c *HanzoFilerClient) CreateEntryOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerCreateEntryOrdinal, on.ID, nil, "CreateEntry")
}

func (c *HanzoFilerClient) UpdateEntry(req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerUpdateEntryOrdinal, rpc.NoTarget, req, "UpdateEntry")
}

// UpdateEntryOn issues UpdateEntry as a dependent call pipelined on on.
func (c *HanzoFilerClient) UpdateEntryOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerUpdateEntryOrdinal, on.ID, nil, "UpdateEntry")
}

func (c *HanzoFilerClient) TouchAccessTime(req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerTouchAccessTimeOrdinal, rpc.NoTarget, req, "TouchAccessTime")
}

// TouchAccessTimeOn issues TouchAccessTime as a dependent call pipelined on on.
func (c *HanzoFilerClient) TouchAccessTimeOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerTouchAccessTimeOrdinal, on.ID, nil, "TouchAccessTime")
}

func (c *HanzoFilerClient) AppendToEntry(req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerAppendToEntryOrdinal, rpc.NoTarget, req, "AppendToEntry")
}

// AppendToEntryOn issues AppendToEntry as a dependent call pipelined on on.
func (c *HanzoFilerClient) AppendToEntryOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerAppendToEntryOrdinal, on.ID, nil, "AppendToEntry")
}

func (c *HanzoFilerClient) DeleteEntry(req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerDeleteEntryOrdinal, rpc.NoTarget, req, "DeleteEntry")
}

// DeleteEntryOn issues DeleteEntry as a dependent call pipelined on on.
func (c *HanzoFilerClient) DeleteEntryOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerDeleteEntryOrdinal, on.ID, nil, "DeleteEntry")
}

func (c *HanzoFilerClient) ObjectTransaction(req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerObjectTransactionOrdinal, rpc.NoTarget, req, "ObjectTransaction")
}

// ObjectTransactionOn issues ObjectTransaction as a dependent call pipelined on on.
func (c *HanzoFilerClient) ObjectTransactionOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerObjectTransactionOrdinal, on.ID, nil, "ObjectTransaction")
}

func (c *HanzoFilerClient) ObjectTransactionBatch(req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerObjectTransactionBatchOrdinal, rpc.NoTarget, req, "ObjectTransactionBatch")
}

// ObjectTransactionBatchOn issues ObjectTransactionBatch as a dependent call
// pipelined on on.
func (c *HanzoFilerClient) ObjectTransactionBatchOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerObjectTransactionBatchOrdinal, on.ID, nil, "ObjectTransactionBatch")
}

func (c *HanzoFilerClient) PosixLock(req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerPosixLockOrdinal, rpc.NoTarget, req, "PosixLock")
}

// PosixLockOn issues PosixLock as a dependent call pipelined on on.
func (c *HanzoFilerClient) PosixLockOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerPosixLockOrdinal, on.ID, nil, "PosixLock")
}

func (c *HanzoFilerClient) AtomicRenameEntry(req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerAtomicRenameEntryOrdinal, rpc.NoTarget, req, "AtomicRenameEntry")
}

// AtomicRenameEntryOn issues AtomicRenameEntry as a dependent call pipelined on on.
func (c *HanzoFilerClient) AtomicRenameEntryOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerAtomicRenameEntryOrdinal, on.ID, nil, "AtomicRenameEntry")
}

// StreamRenameEntry is a server-streaming RPC.
// STREAMING: the streaming body lands when the transport streaming primitive
// ships. Until then this unary stub issues a single request and returns the
// first response frame, so the ordinal and wire types are exercised end-to-end.
func (c *HanzoFilerClient) StreamRenameEntry(req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerStreamRenameEntryOrdinal, rpc.NoTarget, req, "StreamRenameEntry")
}

// StreamMutateEntry is a bidirectional-streaming RPC.
// STREAMING: the streaming body lands when the transport streaming primitive
// ships. Until then this unary stub issues a single request and returns the
// first response frame, so the ordinal and wire types are exercised end-to-end.
func (c *HanzoFilerClient) StreamMutateEntry(req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerStreamMutateEntryOrdinal, rpc.NoTarget, req, "StreamMutateEntry")
}

func (c *HanzoFilerClient) AssignVolume(req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerAssignVolumeOrdinal, rpc.NoTarget, req, "AssignVolume")
}

// AssignVolumeOn issues AssignVolume as a dependent call pipelined on on.
func (c *HanzoFilerClient) AssignVolumeOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerAssignVolumeOrdinal, on.ID, nil, "AssignVolume")
}

func (c *HanzoFilerClient) LookupVolume(req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerLookupVolumeOrdinal, rpc.NoTarget, req, "LookupVolume")
}

// LookupVolumeOn issues LookupVolume as a dependent call pipelined on on.
func (c *HanzoFilerClient) LookupVolumeOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerLookupVolumeOrdinal, on.ID, nil, "LookupVolume")
}

func (c *HanzoFilerClient) CollectionList(req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerCollectionListOrdinal, rpc.NoTarget, req, "CollectionList")
}

// CollectionListOn issues CollectionList as a dependent call pipelined on on.
func (c *HanzoFilerClient) CollectionListOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerCollectionListOrdinal, on.ID, nil, "CollectionList")
}

func (c *HanzoFilerClient) DeleteCollection(req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerDeleteCollectionOrdinal, rpc.NoTarget, req, "DeleteCollection")
}

// DeleteCollectionOn issues DeleteCollection as a dependent call pipelined on on.
func (c *HanzoFilerClient) DeleteCollectionOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerDeleteCollectionOrdinal, on.ID, nil, "DeleteCollection")
}

func (c *HanzoFilerClient) Statistics(req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerStatisticsOrdinal, rpc.NoTarget, req, "Statistics")
}

// StatisticsOn issues Statistics as a dependent call pipelined on on.
func (c *HanzoFilerClient) StatisticsOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerStatisticsOrdinal, on.ID, nil, "Statistics")
}

func (c *HanzoFilerClient) Ping(req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerPingOrdinal, rpc.NoTarget, req, "Ping")
}

// PingOn issues Ping as a dependent call pipelined on on.
func (c *HanzoFilerClient) PingOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerPingOrdinal, on.ID, nil, "Ping")
}

func (c *HanzoFilerClient) GetFilerConfiguration(req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerGetFilerConfigurationOrdinal, rpc.NoTarget, req, "GetFilerConfiguration")
}

// GetFilerConfigurationOn issues GetFilerConfiguration as a dependent call
// pipelined on on.
func (c *HanzoFilerClient) GetFilerConfigurationOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerGetFilerConfigurationOrdinal, on.ID, nil, "GetFilerConfiguration")
}

// TraverseBfsMetadata is a server-streaming RPC.
// STREAMING: the streaming body lands when the transport streaming primitive
// ships. Until then this unary stub issues a single request and returns the
// first response frame, so the ordinal and wire types are exercised end-to-end.
func (c *HanzoFilerClient) TraverseBfsMetadata(req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerTraverseBfsMetadataOrdinal, rpc.NoTarget, req, "TraverseBfsMetadata")
}

// SubscribeMetadata is a server-streaming RPC.
// STREAMING: the streaming body lands when the transport streaming primitive
// ships. Until then this unary stub issues a single request and returns the
// first response frame, so the ordinal and wire types are exercised end-to-end.
func (c *HanzoFilerClient) SubscribeMetadata(req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerSubscribeMetadataOrdinal, rpc.NoTarget, req, "SubscribeMetadata")
}

// SubscribeLocalMetadata is a server-streaming RPC (shares the SubscribeMetadata
// request/response message types).
// STREAMING: the streaming body lands when the transport streaming primitive
// ships. Until then this unary stub issues a single request and returns the
// first response frame, so the ordinal and wire types are exercised end-to-end.
func (c *HanzoFilerClient) SubscribeLocalMetadata(req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerSubscribeLocalMetadataOrdinal, rpc.NoTarget, req, "SubscribeLocalMetadata")
}

func (c *HanzoFilerClient) ListMetadataSubscribers(req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerListMetadataSubscribersOrdinal, rpc.NoTarget, req, "ListMetadataSubscribers")
}

// ListMetadataSubscribersOn issues ListMetadataSubscribers as a dependent call
// pipelined on on.
func (c *HanzoFilerClient) ListMetadataSubscribersOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerListMetadataSubscribersOrdinal, on.ID, nil, "ListMetadataSubscribers")
}

func (c *HanzoFilerClient) KvGet(req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerKvGetOrdinal, rpc.NoTarget, req, "KvGet")
}

// KvGetOn issues KvGet as a dependent call pipelined on on.
func (c *HanzoFilerClient) KvGetOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerKvGetOrdinal, on.ID, nil, "KvGet")
}

func (c *HanzoFilerClient) KvPut(req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerKvPutOrdinal, rpc.NoTarget, req, "KvPut")
}

// KvPutOn issues KvPut as a dependent call pipelined on on.
func (c *HanzoFilerClient) KvPutOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerKvPutOrdinal, on.ID, nil, "KvPut")
}

func (c *HanzoFilerClient) CacheRemoteObjectToLocalCluster(req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerCacheRemoteObjectToLocalClusterOrdinal, rpc.NoTarget, req, "CacheRemoteObjectToLocalCluster")
}

// CacheRemoteObjectToLocalClusterOn issues CacheRemoteObjectToLocalCluster as a
// dependent call pipelined on on.
func (c *HanzoFilerClient) CacheRemoteObjectToLocalClusterOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerCacheRemoteObjectToLocalClusterOrdinal, on.ID, nil, "CacheRemoteObjectToLocalCluster")
}

func (c *HanzoFilerClient) DistributedLock(req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerDistributedLockOrdinal, rpc.NoTarget, req, "DistributedLock")
}

// DistributedLockOn issues DistributedLock as a dependent call pipelined on on.
func (c *HanzoFilerClient) DistributedLockOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerDistributedLockOrdinal, on.ID, nil, "DistributedLock")
}

func (c *HanzoFilerClient) DistributedUnlock(req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerDistributedUnlockOrdinal, rpc.NoTarget, req, "DistributedUnlock")
}

// DistributedUnlockOn issues DistributedUnlock as a dependent call pipelined on on.
func (c *HanzoFilerClient) DistributedUnlockOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerDistributedUnlockOrdinal, on.ID, nil, "DistributedUnlock")
}

func (c *HanzoFilerClient) FindLockOwner(req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerFindLockOwnerOrdinal, rpc.NoTarget, req, "FindLockOwner")
}

// FindLockOwnerOn issues FindLockOwner as a dependent call pipelined on on.
func (c *HanzoFilerClient) FindLockOwnerOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerFindLockOwnerOrdinal, on.ID, nil, "FindLockOwner")
}

func (c *HanzoFilerClient) TransferLocks(req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerTransferLocksOrdinal, rpc.NoTarget, req, "TransferLocks")
}

// TransferLocksOn issues TransferLocks as a dependent call pipelined on on.
func (c *HanzoFilerClient) TransferLocksOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerTransferLocksOrdinal, on.ID, nil, "TransferLocks")
}

func (c *HanzoFilerClient) ReplicateLock(req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerReplicateLockOrdinal, rpc.NoTarget, req, "ReplicateLock")
}

// ReplicateLockOn issues ReplicateLock as a dependent call pipelined on on.
func (c *HanzoFilerClient) ReplicateLockOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerReplicateLockOrdinal, on.ID, nil, "ReplicateLock")
}

func (c *HanzoFilerClient) MountRegister(req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerMountRegisterOrdinal, rpc.NoTarget, req, "MountRegister")
}

// MountRegisterOn issues MountRegister as a dependent call pipelined on on.
func (c *HanzoFilerClient) MountRegisterOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerMountRegisterOrdinal, on.ID, nil, "MountRegister")
}

func (c *HanzoFilerClient) MountList(req []byte) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerMountListOrdinal, rpc.NoTarget, req, "MountList")
}

// MountListOn issues MountList as a dependent call pipelined on on.
func (c *HanzoFilerClient) MountListOn(on rpc.Promise) (rpc.Promise, []byte, error) {
	return c.unaryCall(HanzoFilerMountListOrdinal, on.ID, nil, "MountList")
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
			return rpc.BuildResponse(rpc.StatusInternal, call.PromiseID, nil), nil
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
