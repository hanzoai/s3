// Code generated from filer.proto; DO NOT EDIT.

package filerwire

import (
	zap "github.com/zap-proto/go"
)

// --- PosixLockRange ---

const (
	posixLockRangeStartOff   = 0  // start (field 1, uint64)
	posixLockRangeEndOff     = 8  // end (field 2, uint64)
	posixLockRangeTypeOff    = 16 // type (field 3, uint32)
	posixLockRangeSidOff     = 24 // sid (field 4, uint64)
	posixLockRangeOwnerOff   = 32 // owner (field 5, uint64)
	posixLockRangePidOff     = 40 // pid (field 6, uint32)
	posixLockRangeIsFlockOff = 44 // is_flock (field 7, bool)
	posixLockRangeSize       = 48
)

// PosixLockRange is a zero-copy view into a ZAP-encoded PosixLockRange message.
type PosixLockRange struct{ o zap.Object }

// WrapPosixLockRange parses b and returns a typed view.
func WrapPosixLockRange(b []byte) (PosixLockRange, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PosixLockRange{}, err
	}
	return PosixLockRange{o: m.Root()}, nil
}

func (t PosixLockRange) Start() uint64 { return t.o.Uint64(posixLockRangeStartOff) }
func (t PosixLockRange) End() uint64   { return t.o.Uint64(posixLockRangeEndOff) }
func (t PosixLockRange) Type() uint32  { return t.o.Uint32(posixLockRangeTypeOff) }
func (t PosixLockRange) Sid() uint64   { return t.o.Uint64(posixLockRangeSidOff) }
func (t PosixLockRange) Owner() uint64 { return t.o.Uint64(posixLockRangeOwnerOff) }
func (t PosixLockRange) Pid() uint32   { return t.o.Uint32(posixLockRangePidOff) }
func (t PosixLockRange) IsFlock() bool { return t.o.Bool(posixLockRangeIsFlockOff) }

// PosixLockRangeInput collects the field values for NewPosixLockRange.
type PosixLockRangeInput struct {
	Start   uint64
	End     uint64
	Type    uint32
	Sid     uint64
	Owner   uint64
	Pid     uint32
	IsFlock bool
}

// NewPosixLockRange builds a ZAP-encoded PosixLockRange message from in.
func NewPosixLockRange(in PosixLockRangeInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(posixLockRangeSize)
	ob.SetUint64(posixLockRangeStartOff, in.Start)
	ob.SetUint64(posixLockRangeEndOff, in.End)
	ob.SetUint32(posixLockRangeTypeOff, in.Type)
	ob.SetUint64(posixLockRangeSidOff, in.Sid)
	ob.SetUint64(posixLockRangeOwnerOff, in.Owner)
	ob.SetUint32(posixLockRangePidOff, in.Pid)
	ob.SetBool(posixLockRangeIsFlockOff, in.IsFlock)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- PosixLockRequest ---

const (
	posixLockRequestKeyOff          = 0  // key (field 1, string)
	posixLockRequestIsMovedOff      = 8  // is_moved (field 2, bool)
	posixLockRequestOpOff           = 12 // op (field 3, enum PosixLockOp)
	posixLockRequestLockOff         = 16 // lock (field 4, PosixLockRange object)
	posixLockRequestLocksOff        = 24 // locks (field 5, repeated PosixLockRange)
	posixLockRequestCoolingProbeOff = 32 // cooling_probe (field 6, bool)
	posixLockRequestSize            = 40
)

// PosixLockRequest is a zero-copy view into a ZAP-encoded PosixLockRequest message.
type PosixLockRequest struct{ o zap.Object }

// WrapPosixLockRequest parses b and returns a typed view.
func WrapPosixLockRequest(b []byte) (PosixLockRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PosixLockRequest{}, err
	}
	return PosixLockRequest{o: m.Root()}, nil
}

func (t PosixLockRequest) Key() string   { return t.o.Text(posixLockRequestKeyOff) }
func (t PosixLockRequest) IsMoved() bool { return t.o.Bool(posixLockRequestIsMovedOff) }

// Op reads the op enum (see PosixLockOp* constants).
func (t PosixLockRequest) Op() uint32         { return t.o.Uint32(posixLockRequestOpOff) }
func (t PosixLockRequest) CoolingProbe() bool { return t.o.Bool(posixLockRequestCoolingProbeOff) }

// Lock returns the nested PosixLockRange (null-object if unset).
func (t PosixLockRequest) Lock() PosixLockRange {
	return PosixLockRange{o: t.o.Object(posixLockRequestLockOff)}
}

// HasLock reports whether the lock field is present.
func (t PosixLockRequest) HasLock() bool { return !t.o.Object(posixLockRequestLockOff).IsNull() }

// LocksLen returns the number of locks elements (proto field 5, repeated
// PosixLockRange).
func (t PosixLockRequest) LocksLen() int { return t.o.List(posixLockRequestLocksOff).Len() }

// Locks returns the i-th locks element as a typed PosixLockRange view.
func (t PosixLockRequest) Locks(i int) PosixLockRange {
	return PosixLockRange{o: t.o.List(posixLockRequestLocksOff).ObjectAt(i)}
}

// PosixLockRequestInput collects the field values for NewPosixLockRequest. Lock,
// when non-nil, is a standalone PosixLockRange buffer; each Locks element is a
// standalone PosixLockRange buffer.
type PosixLockRequestInput struct {
	Key          string
	IsMoved      bool
	Op           uint32
	Lock         []byte
	Locks        [][]byte
	CoolingProbe bool
}

// NewPosixLockRequest builds a ZAP-encoded PosixLockRequest message from in.
func NewPosixLockRequest(in PosixLockRequestInput) []byte {
	b := zap.NewBuilder(256)
	locksOff, locksLen := addObjectList(b, in.Locks)
	ob := b.StartObject(posixLockRequestSize)
	ob.SetText(posixLockRequestKeyOff, in.Key)
	ob.SetBool(posixLockRequestIsMovedOff, in.IsMoved)
	ob.SetUint32(posixLockRequestOpOff, in.Op)
	setNestedObject(b, ob, posixLockRequestLockOff, in.Lock)
	ob.SetList(posixLockRequestLocksOff, locksOff, locksLen)
	ob.SetBool(posixLockRequestCoolingProbeOff, in.CoolingProbe)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- PosixLockResponse ---

const (
	posixLockResponseGrantedOff     = 0 // granted (field 1, bool)
	posixLockResponseHasConflictOff = 1 // has_conflict (field 2, bool)
	posixLockResponseConflictOff    = 4 // conflict (field 3, PosixLockRange object)
	posixLockResponseSize           = 8
)

// PosixLockResponse is a zero-copy view into a ZAP-encoded PosixLockResponse message.
type PosixLockResponse struct{ o zap.Object }

// WrapPosixLockResponse parses b and returns a typed view.
func WrapPosixLockResponse(b []byte) (PosixLockResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PosixLockResponse{}, err
	}
	return PosixLockResponse{o: m.Root()}, nil
}

func (t PosixLockResponse) Granted() bool     { return t.o.Bool(posixLockResponseGrantedOff) }
func (t PosixLockResponse) HasConflict() bool { return t.o.Bool(posixLockResponseHasConflictOff) }

// Conflict returns the nested PosixLockRange (null-object if unset).
func (t PosixLockResponse) Conflict() PosixLockRange {
	return PosixLockRange{o: t.o.Object(posixLockResponseConflictOff)}
}

// HasConflictRange reports whether the conflict field object is present.
func (t PosixLockResponse) HasConflictRange() bool {
	return !t.o.Object(posixLockResponseConflictOff).IsNull()
}

// PosixLockResponseInput collects the field values for NewPosixLockResponse.
// Conflict, when non-nil, is a standalone PosixLockRange buffer.
type PosixLockResponseInput struct {
	Granted     bool
	HasConflict bool
	Conflict    []byte
}

// NewPosixLockResponse builds a ZAP-encoded PosixLockResponse message from in.
func NewPosixLockResponse(in PosixLockResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(posixLockResponseSize)
	ob.SetBool(posixLockResponseGrantedOff, in.Granted)
	ob.SetBool(posixLockResponseHasConflictOff, in.HasConflict)
	setNestedObject(b, ob, posixLockResponseConflictOff, in.Conflict)
	ob.FinishAsRoot()
	return b.Finish()
}
