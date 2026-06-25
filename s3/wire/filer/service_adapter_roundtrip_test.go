// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package filerwire

import (
	"bytes"
	"context"
	"testing"
)

// memFiler is an in-memory Backend for the unary round-trip test. It does NOT
// touch the real filer engine — it stores entries in a map so the test proves
// the bytes survive a real socket as zero-copy ZAP buffers. Only CreateEntry and
// LookupDirectoryEntry carry logic; the other 26 unary methods return a trivial
// built response so the Backend is fully satisfied.
type memFiler struct {
	entries map[string][]byte // "dir/name" -> a NewEntry buffer
}

func newMemFiler() *memFiler { return &memFiler{entries: map[string][]byte{}} }

func (m *memFiler) CreateEntry(req CreateEntryRequest) ([]byte, error) {
	e := req.Entry() // zero-copy view of the nested Entry
	key := req.Directory() + "/" + e.Name()
	// Re-encode the entry as a standalone buffer for storage (the bytes ARE the
	// message; a single-field copy is enough for the round-trip assertion).
	m.entries[key] = NewEntry(EntryInput{
		Name:        e.Name(),
		IsDirectory: e.IsDirectory(),
		Content:     append([]byte(nil), e.Content()...),
	})
	return NewCreateEntryResponse(CreateEntryResponseInput{ErrorCode: FilerErrorOK}), nil
}

func (m *memFiler) LookupDirectoryEntry(req LookupDirectoryEntryRequest) ([]byte, error) {
	key := req.Directory() + "/" + req.Name()
	entryBuf := m.entries[key] // nil -> response with no entry
	return NewLookupDirectoryEntryResponse(LookupDirectoryEntryResponseInput{Entry: entryBuf}), nil
}

// The remaining 26 unary methods: trivial built responses (the test does not
// exercise them, but the Backend contract must be total).
func (m *memFiler) UpdateEntry(UpdateEntryRequest) ([]byte, error) {
	return NewUpdateEntryResponse(UpdateEntryResponseInput{}), nil
}
func (m *memFiler) TouchAccessTime(TouchAccessTimeRequest) ([]byte, error) {
	return NewTouchAccessTimeResponse(TouchAccessTimeResponseInput{}), nil
}
func (m *memFiler) AppendToEntry(AppendToEntryRequest) ([]byte, error) {
	return NewAppendToEntryResponse(AppendToEntryResponseInput{}), nil
}
func (m *memFiler) DeleteEntry(DeleteEntryRequest) ([]byte, error) {
	return NewDeleteEntryResponse(DeleteEntryResponseInput{}), nil
}
func (m *memFiler) ObjectTransaction(ObjectTransactionRequest) ([]byte, error) {
	return NewObjectTransactionResponse(ObjectTransactionResponseInput{ErrorCode: FilerErrorOK}), nil
}
func (m *memFiler) ObjectTransactionBatch(ObjectTransactionBatchRequest) ([]byte, error) {
	return NewObjectTransactionBatchResponse(ObjectTransactionBatchResponseInput{}), nil
}
func (m *memFiler) PosixLock(PosixLockRequest) ([]byte, error) {
	return NewPosixLockResponse(PosixLockResponseInput{}), nil
}
func (m *memFiler) AtomicRenameEntry(AtomicRenameEntryRequest) ([]byte, error) {
	return NewAtomicRenameEntryResponse(AtomicRenameEntryResponseInput{}), nil
}
func (m *memFiler) AssignVolume(AssignVolumeRequest) ([]byte, error) {
	return NewAssignVolumeResponse(AssignVolumeResponseInput{}), nil
}
func (m *memFiler) LookupVolume(LookupVolumeRequest) ([]byte, error) {
	return NewLookupVolumeResponse(LookupVolumeResponseInput{}), nil
}
func (m *memFiler) CollectionList(CollectionListRequest) ([]byte, error) {
	return NewCollectionListResponse(CollectionListResponseInput{}), nil
}
func (m *memFiler) DeleteCollection(DeleteCollectionRequest) ([]byte, error) {
	return NewDeleteCollectionResponse(DeleteCollectionResponseInput{}), nil
}
func (m *memFiler) Statistics(StatisticsRequest) ([]byte, error) {
	return NewStatisticsResponse(StatisticsResponseInput{}), nil
}
func (m *memFiler) Ping(PingRequest) ([]byte, error) {
	return NewPingResponse(PingResponseInput{}), nil
}
func (m *memFiler) GetFilerConfiguration(GetFilerConfigurationRequest) ([]byte, error) {
	return NewGetFilerConfigurationResponse(GetFilerConfigurationResponseInput{}), nil
}
func (m *memFiler) ListMetadataSubscribers(ListMetadataSubscribersRequest) ([]byte, error) {
	return NewListMetadataSubscribersResponse(ListMetadataSubscribersResponseInput{}), nil
}
func (m *memFiler) KvGet(KvGetRequest) ([]byte, error) {
	return NewKvGetResponse(KvGetResponseInput{}), nil
}
func (m *memFiler) KvPut(KvPutRequest) ([]byte, error) {
	return NewKvPutResponse(KvPutResponseInput{}), nil
}
func (m *memFiler) CacheRemoteObjectToLocalCluster(CacheRemoteObjectToLocalClusterRequest) ([]byte, error) {
	return NewCacheRemoteObjectToLocalClusterResponse(CacheRemoteObjectToLocalClusterResponseInput{}), nil
}
func (m *memFiler) DistributedLock(LockRequest) ([]byte, error) {
	return NewLockResponse(LockResponseInput{}), nil
}
func (m *memFiler) DistributedUnlock(UnlockRequest) ([]byte, error) {
	return NewUnlockResponse(UnlockResponseInput{}), nil
}
func (m *memFiler) FindLockOwner(FindLockOwnerRequest) ([]byte, error) {
	return NewFindLockOwnerResponse(FindLockOwnerResponseInput{}), nil
}
func (m *memFiler) TransferLocks(TransferLocksRequest) ([]byte, error) {
	return NewTransferLocksResponse(TransferLocksResponseInput{}), nil
}
func (m *memFiler) ReplicateLock(ReplicateLockRequest) ([]byte, error) {
	return NewReplicateLockResponse(ReplicateLockResponseInput{}), nil
}
func (m *memFiler) MountRegister(MountRegisterRequest) ([]byte, error) {
	return NewMountRegisterResponse(MountRegisterResponseInput{}), nil
}
func (m *memFiler) MountList(MountListRequest) ([]byte, error) {
	return NewMountListResponse(MountListResponseInput{}), nil
}

// TestUnaryRoundTrip proves CreateEntry then LookupDirectoryEntry over the
// canonical github.com/zap-proto/go transport across a real TCP socket: the
// bytes cross the wire as ZAP RPC envelopes carrying zero-copy filerwire
// payloads — no HTTP, no protobuf, no struct marshaling.
func TestUnaryRoundTrip(t *testing.T) {
	store := newMemFiler()

	srv, err := Serve("tcp", "127.0.0.1:0", store)
	if err != nil {
		t.Fatalf("Serve: %v", err)
	}
	defer srv.Close()

	cli, err := Dial("tcp", srv.Addr().String())
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer cli.Close()

	wantContent := []byte("the bytes are the message")
	entry := NewEntry(EntryInput{Name: "x.bin", IsDirectory: false, Content: wantContent})

	createResp, err := cli.CreateEntry(context.Background(), NewCreateEntryRequest(CreateEntryRequestInput{
		Directory:                "/buckets/blobs/team-go",
		Entry:                    entry,
		SkipCheckParentDirectory: true,
	}))
	if err != nil {
		t.Fatalf("CreateEntry: %v", err)
	}
	if createResp.ErrorCode() != FilerErrorOK {
		t.Fatalf("CreateEntry error_code = %d, want OK", createResp.ErrorCode())
	}

	lookupResp, err := cli.LookupDirectoryEntry(context.Background(), NewLookupDirectoryEntryRequest(LookupDirectoryEntryRequestInput{
		Directory: "/buckets/blobs/team-go",
		Name:      "x.bin",
	}))
	if err != nil {
		t.Fatalf("LookupDirectoryEntry: %v", err)
	}
	if !lookupResp.HasEntry() {
		t.Fatalf("LookupDirectoryEntry: entry missing")
	}
	got := lookupResp.Entry()
	if got.Name() != "x.bin" {
		t.Fatalf("entry name = %q, want x.bin", got.Name())
	}
	if !bytes.Equal(got.Content(), wantContent) {
		t.Fatalf("entry content = %q, want %q", got.Content(), wantContent)
	}

	// A streaming RPC invoked on the unary path must be refused, not faked: the
	// handler returns ErrStreamingEndpoint -> StatusInternal -> a client error.
	if _, body, err := cli.rpc.SubscribeMetadata(context.Background(), NewSubscribeMetadataRequest(SubscribeMetadataRequestInput{ClientName: "probe"})); err == nil {
		t.Fatalf("SubscribeMetadata unary call should be refused, got body %d bytes", len(body))
	}
}
