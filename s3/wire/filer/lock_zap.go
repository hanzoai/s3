// Code generated from filer.proto; DO NOT EDIT.

package filerwire

import (
	zap "github.com/zap-proto/go"
)

// --- LockRequest ---

const (
	lockRequestNameOff          = 0  // name (field 1, string)
	lockRequestSecondsToLockOff = 8  // seconds_to_lock (field 2, int64)
	lockRequestRenewTokenOff    = 16 // renew_token (field 3, string)
	lockRequestIsMovedOff       = 24 // is_moved (field 4, bool)
	lockRequestOwnerOff         = 32 // owner (field 5, string)
	lockRequestSize             = 40
)

// LockRequest is a zero-copy view into a ZAP-encoded LockRequest message.
type LockRequest struct{ o zap.Object }

// WrapLockRequest parses b and returns a typed view.
func WrapLockRequest(b []byte) (LockRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return LockRequest{}, err
	}
	return LockRequest{o: m.Root()}, nil
}

func (t LockRequest) Name() string         { return t.o.Text(lockRequestNameOff) }
func (t LockRequest) SecondsToLock() int64 { return t.o.Int64(lockRequestSecondsToLockOff) }
func (t LockRequest) RenewToken() string   { return t.o.Text(lockRequestRenewTokenOff) }
func (t LockRequest) IsMoved() bool        { return t.o.Bool(lockRequestIsMovedOff) }
func (t LockRequest) Owner() string        { return t.o.Text(lockRequestOwnerOff) }

// LockRequestInput collects the field values for NewLockRequest.
type LockRequestInput struct {
	Name          string
	SecondsToLock int64
	RenewToken    string
	IsMoved       bool
	Owner         string
}

// NewLockRequest builds a ZAP-encoded LockRequest message from in.
func NewLockRequest(in LockRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(lockRequestSize)
	ob.SetText(lockRequestNameOff, in.Name)
	ob.SetInt64(lockRequestSecondsToLockOff, in.SecondsToLock)
	ob.SetText(lockRequestRenewTokenOff, in.RenewToken)
	ob.SetBool(lockRequestIsMovedOff, in.IsMoved)
	ob.SetText(lockRequestOwnerOff, in.Owner)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- LockResponse ---

const (
	lockResponseRenewTokenOff      = 0  // renew_token (field 1, string)
	lockResponseLockOwnerOff       = 8  // lock_owner (field 2, string)
	lockResponseLockHostMovedToOff = 16 // lock_host_moved_to (field 3, string)
	lockResponseErrorOff           = 24 // error (field 4, string)
	lockResponseGenerationOff      = 32 // generation (field 5, int64)
	lockResponseSize               = 40
)

// LockResponse is a zero-copy view into a ZAP-encoded LockResponse message.
type LockResponse struct{ o zap.Object }

// WrapLockResponse parses b and returns a typed view.
func WrapLockResponse(b []byte) (LockResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return LockResponse{}, err
	}
	return LockResponse{o: m.Root()}, nil
}

func (t LockResponse) RenewToken() string      { return t.o.Text(lockResponseRenewTokenOff) }
func (t LockResponse) LockOwner() string       { return t.o.Text(lockResponseLockOwnerOff) }
func (t LockResponse) LockHostMovedTo() string { return t.o.Text(lockResponseLockHostMovedToOff) }
func (t LockResponse) Error() string           { return t.o.Text(lockResponseErrorOff) }
func (t LockResponse) Generation() int64       { return t.o.Int64(lockResponseGenerationOff) }

// LockResponseInput collects the field values for NewLockResponse.
type LockResponseInput struct {
	RenewToken      string
	LockOwner       string
	LockHostMovedTo string
	Error           string
	Generation      int64
}

// NewLockResponse builds a ZAP-encoded LockResponse message from in.
func NewLockResponse(in LockResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(lockResponseSize)
	ob.SetText(lockResponseRenewTokenOff, in.RenewToken)
	ob.SetText(lockResponseLockOwnerOff, in.LockOwner)
	ob.SetText(lockResponseLockHostMovedToOff, in.LockHostMovedTo)
	ob.SetText(lockResponseErrorOff, in.Error)
	ob.SetInt64(lockResponseGenerationOff, in.Generation)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- UnlockRequest ---

const (
	unlockRequestNameOff       = 0  // name (field 1, string)
	unlockRequestRenewTokenOff = 8  // renew_token (field 2, string)
	unlockRequestIsMovedOff    = 16 // is_moved (field 3, bool)
	unlockRequestSize          = 24
)

// UnlockRequest is a zero-copy view into a ZAP-encoded UnlockRequest message.
type UnlockRequest struct{ o zap.Object }

// WrapUnlockRequest parses b and returns a typed view.
func WrapUnlockRequest(b []byte) (UnlockRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return UnlockRequest{}, err
	}
	return UnlockRequest{o: m.Root()}, nil
}

func (t UnlockRequest) Name() string       { return t.o.Text(unlockRequestNameOff) }
func (t UnlockRequest) RenewToken() string { return t.o.Text(unlockRequestRenewTokenOff) }
func (t UnlockRequest) IsMoved() bool      { return t.o.Bool(unlockRequestIsMovedOff) }

// UnlockRequestInput collects the field values for NewUnlockRequest.
type UnlockRequestInput struct {
	Name       string
	RenewToken string
	IsMoved    bool
}

// NewUnlockRequest builds a ZAP-encoded UnlockRequest message from in.
func NewUnlockRequest(in UnlockRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(unlockRequestSize)
	ob.SetText(unlockRequestNameOff, in.Name)
	ob.SetText(unlockRequestRenewTokenOff, in.RenewToken)
	ob.SetBool(unlockRequestIsMovedOff, in.IsMoved)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- UnlockResponse ---

const (
	unlockResponseErrorOff   = 0 // error (field 1, string)
	unlockResponseMovedToOff = 8 // moved_to (field 2, string)
	unlockResponseSize       = 16
)

// UnlockResponse is a zero-copy view into a ZAP-encoded UnlockResponse message.
type UnlockResponse struct{ o zap.Object }

// WrapUnlockResponse parses b and returns a typed view.
func WrapUnlockResponse(b []byte) (UnlockResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return UnlockResponse{}, err
	}
	return UnlockResponse{o: m.Root()}, nil
}

func (t UnlockResponse) Error() string   { return t.o.Text(unlockResponseErrorOff) }
func (t UnlockResponse) MovedTo() string { return t.o.Text(unlockResponseMovedToOff) }

// UnlockResponseInput collects the field values for NewUnlockResponse.
type UnlockResponseInput struct {
	Error   string
	MovedTo string
}

// NewUnlockResponse builds a ZAP-encoded UnlockResponse message from in.
func NewUnlockResponse(in UnlockResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(unlockResponseSize)
	ob.SetText(unlockResponseErrorOff, in.Error)
	ob.SetText(unlockResponseMovedToOff, in.MovedTo)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- FindLockOwnerRequest ---

const (
	findLockOwnerRequestNameOff    = 0 // name (field 1, string)
	findLockOwnerRequestIsMovedOff = 8 // is_moved (field 2, bool)
	findLockOwnerRequestSize       = 16
)

// FindLockOwnerRequest is a zero-copy view into a ZAP-encoded FindLockOwnerRequest message.
type FindLockOwnerRequest struct{ o zap.Object }

// WrapFindLockOwnerRequest parses b and returns a typed view.
func WrapFindLockOwnerRequest(b []byte) (FindLockOwnerRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return FindLockOwnerRequest{}, err
	}
	return FindLockOwnerRequest{o: m.Root()}, nil
}

func (t FindLockOwnerRequest) Name() string  { return t.o.Text(findLockOwnerRequestNameOff) }
func (t FindLockOwnerRequest) IsMoved() bool { return t.o.Bool(findLockOwnerRequestIsMovedOff) }

// FindLockOwnerRequestInput collects the field values for NewFindLockOwnerRequest.
type FindLockOwnerRequestInput struct {
	Name    string
	IsMoved bool
}

// NewFindLockOwnerRequest builds a ZAP-encoded FindLockOwnerRequest message from in.
func NewFindLockOwnerRequest(in FindLockOwnerRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(findLockOwnerRequestSize)
	ob.SetText(findLockOwnerRequestNameOff, in.Name)
	ob.SetBool(findLockOwnerRequestIsMovedOff, in.IsMoved)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- FindLockOwnerResponse ---

const (
	findLockOwnerResponseOwnerOff = 0 // owner (field 1, string)
	findLockOwnerResponseSize     = 8
)

// FindLockOwnerResponse is a zero-copy view into a ZAP-encoded FindLockOwnerResponse message.
type FindLockOwnerResponse struct{ o zap.Object }

// WrapFindLockOwnerResponse parses b and returns a typed view.
func WrapFindLockOwnerResponse(b []byte) (FindLockOwnerResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return FindLockOwnerResponse{}, err
	}
	return FindLockOwnerResponse{o: m.Root()}, nil
}

func (t FindLockOwnerResponse) Owner() string { return t.o.Text(findLockOwnerResponseOwnerOff) }

// FindLockOwnerResponseInput collects the field values for NewFindLockOwnerResponse.
type FindLockOwnerResponseInput struct {
	Owner string
}

// NewFindLockOwnerResponse builds a ZAP-encoded FindLockOwnerResponse message from in.
func NewFindLockOwnerResponse(in FindLockOwnerResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(findLockOwnerResponseSize)
	ob.SetText(findLockOwnerResponseOwnerOff, in.Owner)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- Lock ---

const (
	lockNameOff        = 0  // name (field 1, string)
	lockRenewTokenOff  = 8  // renew_token (field 2, string)
	lockExpiredAtNsOff = 16 // expired_at_ns (field 3, int64)
	lockOwnerOff       = 24 // owner (field 4, string)
	lockGenerationOff  = 32 // generation (field 5, int64)
	lockIsBackupOff    = 40 // is_backup (field 6, bool)
	lockSeqOff         = 48 // seq (field 7, int64)
	lockSize           = 56
)

// Lock is a zero-copy view into a ZAP-encoded Lock message.
type Lock struct{ o zap.Object }

// WrapLock parses b and returns a typed view.
func WrapLock(b []byte) (Lock, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return Lock{}, err
	}
	return Lock{o: m.Root()}, nil
}

func (t Lock) Name() string       { return t.o.Text(lockNameOff) }
func (t Lock) RenewToken() string { return t.o.Text(lockRenewTokenOff) }
func (t Lock) ExpiredAtNs() int64 { return t.o.Int64(lockExpiredAtNsOff) }
func (t Lock) Owner() string      { return t.o.Text(lockOwnerOff) }
func (t Lock) Generation() int64  { return t.o.Int64(lockGenerationOff) }
func (t Lock) IsBackup() bool     { return t.o.Bool(lockIsBackupOff) }
func (t Lock) Seq() int64         { return t.o.Int64(lockSeqOff) }

// LockInput collects the field values for NewLock.
type LockInput struct {
	Name        string
	RenewToken  string
	ExpiredAtNs int64
	Owner       string
	Generation  int64
	IsBackup    bool
	Seq         int64
}

// NewLock builds a ZAP-encoded Lock message from in.
func NewLock(in LockInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(lockSize)
	ob.SetText(lockNameOff, in.Name)
	ob.SetText(lockRenewTokenOff, in.RenewToken)
	ob.SetInt64(lockExpiredAtNsOff, in.ExpiredAtNs)
	ob.SetText(lockOwnerOff, in.Owner)
	ob.SetInt64(lockGenerationOff, in.Generation)
	ob.SetBool(lockIsBackupOff, in.IsBackup)
	ob.SetInt64(lockSeqOff, in.Seq)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- TransferLocksRequest ---

const (
	transferLocksRequestLocksOff = 0 // locks (field 1, repeated Lock)
	transferLocksRequestSize     = 8
)

// TransferLocksRequest is a zero-copy view into a ZAP-encoded TransferLocksRequest message.
type TransferLocksRequest struct{ o zap.Object }

// WrapTransferLocksRequest parses b and returns a typed view.
func WrapTransferLocksRequest(b []byte) (TransferLocksRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return TransferLocksRequest{}, err
	}
	return TransferLocksRequest{o: m.Root()}, nil
}

// LocksLen returns the number of locks elements (proto field 1, repeated Lock).
func (t TransferLocksRequest) LocksLen() int { return t.o.List(transferLocksRequestLocksOff).Len() }

// Locks returns the i-th locks element as a typed Lock view.
func (t TransferLocksRequest) Locks(i int) Lock {
	return Lock{o: t.o.List(transferLocksRequestLocksOff).ObjectAt(i)}
}

// TransferLocksRequestInput collects the field values for NewTransferLocksRequest.
// Each Locks element is a standalone Lock buffer (from NewLock).
type TransferLocksRequestInput struct {
	Locks [][]byte
}

// NewTransferLocksRequest builds a ZAP-encoded TransferLocksRequest message from in.
func NewTransferLocksRequest(in TransferLocksRequestInput) []byte {
	b := zap.NewBuilder(256)
	locksOff, locksLen := addObjectList(b, in.Locks)
	ob := b.StartObject(transferLocksRequestSize)
	ob.SetList(transferLocksRequestLocksOff, locksOff, locksLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- TransferLocksResponse ---

const transferLocksResponseSize = 0

// TransferLocksResponse is a zero-copy view into a ZAP-encoded TransferLocksResponse
// message. TransferLocksResponse carries no fields.
type TransferLocksResponse struct{ o zap.Object }

// WrapTransferLocksResponse parses b and returns a typed view.
func WrapTransferLocksResponse(b []byte) (TransferLocksResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return TransferLocksResponse{}, err
	}
	return TransferLocksResponse{o: m.Root()}, nil
}

// TransferLocksResponseInput collects the field values for NewTransferLocksResponse.
// TransferLocksResponse has no fields.
type TransferLocksResponseInput struct{}

// NewTransferLocksResponse builds a ZAP-encoded TransferLocksResponse message.
func NewTransferLocksResponse(in TransferLocksResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(transferLocksResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ReplicateLockRequest ---

const (
	replicateLockRequestNameOff        = 0  // name (field 1, string)
	replicateLockRequestRenewTokenOff  = 8  // renew_token (field 2, string)
	replicateLockRequestExpiredAtNsOff = 16 // expired_at_ns (field 3, int64)
	replicateLockRequestOwnerOff       = 24 // owner (field 4, string)
	replicateLockRequestGenerationOff  = 32 // generation (field 5, int64)
	replicateLockRequestIsUnlockOff    = 40 // is_unlock (field 6, bool)
	replicateLockRequestSeqOff         = 48 // seq (field 7, int64)
	replicateLockRequestSize           = 56
)

// ReplicateLockRequest is a zero-copy view into a ZAP-encoded ReplicateLockRequest message.
type ReplicateLockRequest struct{ o zap.Object }

// WrapReplicateLockRequest parses b and returns a typed view.
func WrapReplicateLockRequest(b []byte) (ReplicateLockRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ReplicateLockRequest{}, err
	}
	return ReplicateLockRequest{o: m.Root()}, nil
}

func (t ReplicateLockRequest) Name() string       { return t.o.Text(replicateLockRequestNameOff) }
func (t ReplicateLockRequest) RenewToken() string { return t.o.Text(replicateLockRequestRenewTokenOff) }
func (t ReplicateLockRequest) ExpiredAtNs() int64 {
	return t.o.Int64(replicateLockRequestExpiredAtNsOff)
}
func (t ReplicateLockRequest) Owner() string     { return t.o.Text(replicateLockRequestOwnerOff) }
func (t ReplicateLockRequest) Generation() int64 { return t.o.Int64(replicateLockRequestGenerationOff) }
func (t ReplicateLockRequest) IsUnlock() bool    { return t.o.Bool(replicateLockRequestIsUnlockOff) }
func (t ReplicateLockRequest) Seq() int64        { return t.o.Int64(replicateLockRequestSeqOff) }

// ReplicateLockRequestInput collects the field values for NewReplicateLockRequest.
type ReplicateLockRequestInput struct {
	Name        string
	RenewToken  string
	ExpiredAtNs int64
	Owner       string
	Generation  int64
	IsUnlock    bool
	Seq         int64
}

// NewReplicateLockRequest builds a ZAP-encoded ReplicateLockRequest message from in.
func NewReplicateLockRequest(in ReplicateLockRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(replicateLockRequestSize)
	ob.SetText(replicateLockRequestNameOff, in.Name)
	ob.SetText(replicateLockRequestRenewTokenOff, in.RenewToken)
	ob.SetInt64(replicateLockRequestExpiredAtNsOff, in.ExpiredAtNs)
	ob.SetText(replicateLockRequestOwnerOff, in.Owner)
	ob.SetInt64(replicateLockRequestGenerationOff, in.Generation)
	ob.SetBool(replicateLockRequestIsUnlockOff, in.IsUnlock)
	ob.SetInt64(replicateLockRequestSeqOff, in.Seq)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ReplicateLockResponse ---

const replicateLockResponseSize = 0

// ReplicateLockResponse is a zero-copy view into a ZAP-encoded ReplicateLockResponse
// message. ReplicateLockResponse carries no fields.
type ReplicateLockResponse struct{ o zap.Object }

// WrapReplicateLockResponse parses b and returns a typed view.
func WrapReplicateLockResponse(b []byte) (ReplicateLockResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ReplicateLockResponse{}, err
	}
	return ReplicateLockResponse{o: m.Root()}, nil
}

// ReplicateLockResponseInput collects the field values for NewReplicateLockResponse.
// ReplicateLockResponse has no fields.
type ReplicateLockResponseInput struct{}

// NewReplicateLockResponse builds a ZAP-encoded ReplicateLockResponse message.
func NewReplicateLockResponse(in ReplicateLockResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(replicateLockResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}
