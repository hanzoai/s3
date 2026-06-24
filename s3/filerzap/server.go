// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// server.go is the server-side mirror of client.go: a filerwire.Backend that
// wraps the existing gRPC-shaped filer (filer_pb.HanzoFilerServer) so it answers
// over the canonical github.com/zap-proto/go transport. Each unary RPC is one
// hop each way: wire request view -> <Rpc>ReqFromView -> the existing
// filer_pb.HanzoFilerServer method -> <Rpc>RespToWire -> wire response bytes.
// This is the second half of the filer strangler: with NewServerBackend wired
// into filerwire.Serve (the analogue of the gRPC RegisterHanzoFilerServer), the
// ZAP client in client.go reaches the real filer engine over ZAP — no protobuf
// framing on the wire, all engine logic reused unchanged.
//
// Core entry CRUD (Lookup/Create/Update/Delete) is migrated. The remaining unary
// RPCs return errFilerRPCNotMigrated until their converters land; the streaming
// RPCs live in the filerstream package. The bytes returned ARE the message.

package filerzap

import (
	"context"
	"errors"
	"strings"

	filer_pb "github.com/hanzoai/s3/s3/pb/filer_pb"
	filerwire "github.com/hanzoai/s3/s3/wire/filer"
)

// errFilerRPCNotMigrated marks a unary filer RPC whose ZAP server-side converter
// is not yet written. It surfaces as a StatusInternal to the ZAP caller.
var errFilerRPCNotMigrated = errors.New("filerzap: filer RPC not yet migrated to the ZAP server backend")

// serverBackend adapts a filer_pb.HanzoFilerServer to filerwire.Backend.
type serverBackend struct {
	fs  filer_pb.HanzoFilerServer
	ctx context.Context
}

// NewServerBackend returns a filerwire.Backend that serves fs over ZAP. Pass it
// to filerwire.Serve / filerwire.Dispatch.
func NewServerBackend(fs filer_pb.HanzoFilerServer) filerwire.Backend {
	return serverBackend{fs: fs, ctx: context.Background()}
}

// --- migrated: entry CRUD ---

func (b serverBackend) LookupDirectoryEntry(v filerwire.LookupDirectoryEntryRequest) ([]byte, error) {
	resp, err := b.fs.LookupDirectoryEntry(b.ctx, &filer_pb.LookupDirectoryEntryRequest{
		Directory: v.Directory(), Name: v.Name(),
	})
	if err != nil {
		// Not-found is conveyed as an empty (nil Entry) OK response, not a
		// transport error: the rpc envelope carries only a status code, so a
		// propagated error would reach the caller as a generic "status 500" and
		// lose the ErrNotFound semantics. filer_pb.LookupEntry maps a nil Entry
		// back to ErrNotFound, which is exactly the gRPC contract callers expect.
		if isFilerNotFound(err) {
			return filerwire.NewLookupDirectoryEntryResponse(filerwire.LookupDirectoryEntryResponseInput{}), nil
		}
		return nil, err
	}
	in := filerwire.LookupDirectoryEntryResponseInput{}
	if resp.Entry != nil {
		in.Entry = EntryToWire(resp.Entry)
	}
	return filerwire.NewLookupDirectoryEntryResponse(in), nil
}

// isFilerNotFound reports whether err is the filer's not-found sentinel, however
// it arrives (the engine's ErrNotFound, a wrapped copy, or a gRPC status whose
// message carries it).
func isFilerNotFound(err error) bool {
	return errors.Is(err, filer_pb.ErrNotFound) ||
		strings.Contains(err.Error(), filer_pb.ErrNotFound.Error())
}

func (b serverBackend) CreateEntry(v filerwire.CreateEntryRequest) ([]byte, error) {
	req := &filer_pb.CreateEntryRequest{
		Directory:                v.Directory(),
		OExcl:                    v.OExcl(),
		IsFromOtherCluster:       v.IsFromOtherCluster(),
		SkipCheckParentDirectory: v.SkipCheckParentDirectory(),
		Signatures:               signaturesFromView(v.SignaturesLen(), v.Signature),
		Condition:                writeConditionFromView(v.HasCondition(), v.Condition()),
	}
	if v.HasEntry() {
		req.Entry = entryFromView(v.Entry())
	}
	resp, err := b.fs.CreateEntry(b.ctx, req)
	if err != nil {
		return nil, err
	}
	return filerwire.NewCreateEntryResponse(filerwire.CreateEntryResponseInput{
		Error:         resp.Error,
		ErrorCode:     uint32(resp.ErrorCode),
		MetadataEvent: SubscribeMetadataResponseToWire(resp.MetadataEvent),
	}), nil
}

func (b serverBackend) UpdateEntry(v filerwire.UpdateEntryRequest) ([]byte, error) {
	req := &filer_pb.UpdateEntryRequest{
		Directory:          v.Directory(),
		IsFromOtherCluster: v.IsFromOtherCluster(),
		Signatures:         signaturesFromView(v.SignaturesLen(), v.Signature),
		ExpectedExtended:   extendedFromView(v.ExpectedExtendedLen(), v.ExpectedExtended),
	}
	if v.HasEntry() {
		req.Entry = entryFromView(v.Entry())
	}
	resp, err := b.fs.UpdateEntry(b.ctx, req)
	if err != nil {
		return nil, err
	}
	return filerwire.NewUpdateEntryResponse(filerwire.UpdateEntryResponseInput{
		MetadataEvent: SubscribeMetadataResponseToWire(resp.MetadataEvent),
	}), nil
}

func (b serverBackend) DeleteEntry(v filerwire.DeleteEntryRequest) ([]byte, error) {
	resp, err := b.fs.DeleteEntry(b.ctx, &filer_pb.DeleteEntryRequest{
		Directory:            v.Directory(),
		Name:                 v.Name(),
		IsDeleteData:         v.IsDeleteData(),
		IsRecursive:          v.IsRecursive(),
		IgnoreRecursiveError: v.IgnoreRecursiveError(),
		IsFromOtherCluster:   v.IsFromOtherCluster(),
		IfNotModifiedAfter:   v.IfNotModifiedAfter(),
		Signatures:           signaturesFromView(v.SignaturesLen(), v.Signature),
	})
	if err != nil {
		return nil, err
	}
	return filerwire.NewDeleteEntryResponse(filerwire.DeleteEntryResponseInput{
		Error:         resp.Error,
		MetadataEvent: SubscribeMetadataResponseToWire(resp.MetadataEvent),
	}), nil
}

// --- not yet migrated (converters pending) ---

func (b serverBackend) TouchAccessTime(filerwire.TouchAccessTimeRequest) ([]byte, error) {
	return nil, errFilerRPCNotMigrated
}
func (b serverBackend) AppendToEntry(filerwire.AppendToEntryRequest) ([]byte, error) {
	return nil, errFilerRPCNotMigrated
}
func (b serverBackend) ObjectTransaction(filerwire.ObjectTransactionRequest) ([]byte, error) {
	return nil, errFilerRPCNotMigrated
}
func (b serverBackend) ObjectTransactionBatch(filerwire.ObjectTransactionBatchRequest) ([]byte, error) {
	return nil, errFilerRPCNotMigrated
}
func (b serverBackend) PosixLock(filerwire.PosixLockRequest) ([]byte, error) {
	return nil, errFilerRPCNotMigrated
}
func (b serverBackend) AtomicRenameEntry(filerwire.AtomicRenameEntryRequest) ([]byte, error) {
	return nil, errFilerRPCNotMigrated
}
func (b serverBackend) AssignVolume(filerwire.AssignVolumeRequest) ([]byte, error) {
	return nil, errFilerRPCNotMigrated
}
func (b serverBackend) LookupVolume(filerwire.LookupVolumeRequest) ([]byte, error) {
	return nil, errFilerRPCNotMigrated
}
func (b serverBackend) CollectionList(filerwire.CollectionListRequest) ([]byte, error) {
	return nil, errFilerRPCNotMigrated
}
func (b serverBackend) DeleteCollection(filerwire.DeleteCollectionRequest) ([]byte, error) {
	return nil, errFilerRPCNotMigrated
}
func (b serverBackend) Statistics(filerwire.StatisticsRequest) ([]byte, error) {
	return nil, errFilerRPCNotMigrated
}
func (b serverBackend) Ping(filerwire.PingRequest) ([]byte, error) {
	return nil, errFilerRPCNotMigrated
}
func (b serverBackend) GetFilerConfiguration(filerwire.GetFilerConfigurationRequest) ([]byte, error) {
	return nil, errFilerRPCNotMigrated
}
func (b serverBackend) ListMetadataSubscribers(filerwire.ListMetadataSubscribersRequest) ([]byte, error) {
	return nil, errFilerRPCNotMigrated
}
func (b serverBackend) KvGet(filerwire.KvGetRequest) ([]byte, error) {
	return nil, errFilerRPCNotMigrated
}
func (b serverBackend) KvPut(filerwire.KvPutRequest) ([]byte, error) {
	return nil, errFilerRPCNotMigrated
}
func (b serverBackend) CacheRemoteObjectToLocalCluster(filerwire.CacheRemoteObjectToLocalClusterRequest) ([]byte, error) {
	return nil, errFilerRPCNotMigrated
}
func (b serverBackend) DistributedLock(filerwire.LockRequest) ([]byte, error) {
	return nil, errFilerRPCNotMigrated
}
func (b serverBackend) DistributedUnlock(filerwire.UnlockRequest) ([]byte, error) {
	return nil, errFilerRPCNotMigrated
}
func (b serverBackend) FindLockOwner(filerwire.FindLockOwnerRequest) ([]byte, error) {
	return nil, errFilerRPCNotMigrated
}
func (b serverBackend) TransferLocks(filerwire.TransferLocksRequest) ([]byte, error) {
	return nil, errFilerRPCNotMigrated
}
func (b serverBackend) ReplicateLock(filerwire.ReplicateLockRequest) ([]byte, error) {
	return nil, errFilerRPCNotMigrated
}
func (b serverBackend) MountRegister(filerwire.MountRegisterRequest) ([]byte, error) {
	return nil, errFilerRPCNotMigrated
}
func (b serverBackend) MountList(filerwire.MountListRequest) ([]byte, error) {
	return nil, errFilerRPCNotMigrated
}

// --- shared server-side leaf converters ---

func signaturesFromView(n int, at func(int) int32) []int32 {
	if n == 0 {
		return nil
	}
	out := make([]int32, n)
	for i := 0; i < n; i++ {
		out[i] = at(i)
	}
	return out
}

func extendedFromView(n int, at func(int) filerwire.UpdateEntryRequestExpectedExtendedEntry) map[string][]byte {
	if n == 0 {
		return nil
	}
	out := make(map[string][]byte, n)
	for i := 0; i < n; i++ {
		kv := at(i)
		out[kv.Key()] = kv.Value()
	}
	return out
}

func writeConditionFromView(has bool, v filerwire.WriteCondition) *filer_pb.WriteCondition {
	if !has {
		return nil
	}
	n := v.ClausesLen()
	clauses := make([]*filer_pb.WriteCondition_Clause, n)
	for i := 0; i < n; i++ {
		cl := v.Clauses(i)
		clauses[i] = &filer_pb.WriteCondition_Clause{
			Kind:      filer_pb.WriteCondition_Kind(cl.Kind()),
			Etags:     etagsFromView(cl.EtagsLen(), cl.Etag),
			UnixTime:  cl.UnixTime(),
			AllowWeak: cl.AllowWeak(),
			ExtKey:    cl.ExtKey(),
			ExtValue:  cl.ExtValue(),
			GateKey:   cl.GateKey(),
			GateValue: cl.GateValue(),
		}
	}
	return &filer_pb.WriteCondition{Clauses: clauses}
}

func etagsFromView(n int, at func(int) string) []string {
	if n == 0 {
		return nil
	}
	out := make([]string, n)
	for i := 0; i < n; i++ {
		out[i] = at(i)
	}
	return out
}
