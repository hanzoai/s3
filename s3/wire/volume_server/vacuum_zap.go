// Code generated from volume_server.proto; DO NOT EDIT.

package volume_serverwire

import (
	zap "github.com/zap-proto/go"
)

// --- VacuumVolumeCheckRequest ---

const (
	vacuumVolumeCheckRequestVolumeIDOff = 0
	vacuumVolumeCheckRequestSize        = 4
)

// VacuumVolumeCheckRequest is a zero-copy view into a ZAP-encoded
// VacuumVolumeCheckRequest message.
type VacuumVolumeCheckRequest struct{ o zap.Object }

// WrapVacuumVolumeCheckRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVacuumVolumeCheckRequest(b []byte) (VacuumVolumeCheckRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VacuumVolumeCheckRequest{}, err
	}
	return VacuumVolumeCheckRequest{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t VacuumVolumeCheckRequest) VolumeID() uint32 {
	return t.o.Uint32(vacuumVolumeCheckRequestVolumeIDOff)
}

// VacuumVolumeCheckRequestInput collects the field values for
// NewVacuumVolumeCheckRequest.
type VacuumVolumeCheckRequestInput struct {
	VolumeID uint32
}

// NewVacuumVolumeCheckRequest builds a ZAP-encoded VacuumVolumeCheckRequest
// message from in and returns the bytes.
func NewVacuumVolumeCheckRequest(in VacuumVolumeCheckRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(vacuumVolumeCheckRequestSize)
	ob.SetUint32(vacuumVolumeCheckRequestVolumeIDOff, in.VolumeID)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VacuumVolumeCheckResponse ---

const (
	vacuumVolumeCheckResponseGarbageRatioOff = 0
	vacuumVolumeCheckResponseSize            = 8
)

// VacuumVolumeCheckResponse is a zero-copy view into a ZAP-encoded
// VacuumVolumeCheckResponse message.
type VacuumVolumeCheckResponse struct{ o zap.Object }

// WrapVacuumVolumeCheckResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVacuumVolumeCheckResponse(b []byte) (VacuumVolumeCheckResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VacuumVolumeCheckResponse{}, err
	}
	return VacuumVolumeCheckResponse{o: m.Root()}, nil
}

// GarbageRatio reads the garbage_ratio field (proto field 1, double).
func (t VacuumVolumeCheckResponse) GarbageRatio() float64 {
	return t.o.Float64(vacuumVolumeCheckResponseGarbageRatioOff)
}

// VacuumVolumeCheckResponseInput collects the field values for
// NewVacuumVolumeCheckResponse.
type VacuumVolumeCheckResponseInput struct {
	GarbageRatio float64
}

// NewVacuumVolumeCheckResponse builds a ZAP-encoded VacuumVolumeCheckResponse
// message from in and returns the bytes.
func NewVacuumVolumeCheckResponse(in VacuumVolumeCheckResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(vacuumVolumeCheckResponseSize)
	ob.SetFloat64(vacuumVolumeCheckResponseGarbageRatioOff, in.GarbageRatio)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VacuumVolumeCompactRequest ---

const (
	vacuumVolumeCompactRequestVolumeIDOff    = 0
	vacuumVolumeCompactRequestPreallocateOff = 8
	vacuumVolumeCompactRequestSize           = 16
)

// VacuumVolumeCompactRequest is a zero-copy view into a ZAP-encoded
// VacuumVolumeCompactRequest message.
type VacuumVolumeCompactRequest struct{ o zap.Object }

// WrapVacuumVolumeCompactRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVacuumVolumeCompactRequest(b []byte) (VacuumVolumeCompactRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VacuumVolumeCompactRequest{}, err
	}
	return VacuumVolumeCompactRequest{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t VacuumVolumeCompactRequest) VolumeID() uint32 {
	return t.o.Uint32(vacuumVolumeCompactRequestVolumeIDOff)
}

// Preallocate reads the preallocate field (proto field 2, int64).
func (t VacuumVolumeCompactRequest) Preallocate() int64 {
	return t.o.Int64(vacuumVolumeCompactRequestPreallocateOff)
}

// VacuumVolumeCompactRequestInput collects the field values for
// NewVacuumVolumeCompactRequest.
type VacuumVolumeCompactRequestInput struct {
	VolumeID    uint32
	Preallocate int64
}

// NewVacuumVolumeCompactRequest builds a ZAP-encoded VacuumVolumeCompactRequest
// message from in and returns the bytes.
func NewVacuumVolumeCompactRequest(in VacuumVolumeCompactRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(vacuumVolumeCompactRequestSize)
	ob.SetUint32(vacuumVolumeCompactRequestVolumeIDOff, in.VolumeID)
	ob.SetInt64(vacuumVolumeCompactRequestPreallocateOff, in.Preallocate)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VacuumVolumeCompactResponse ---

const (
	vacuumVolumeCompactResponseProcessedBytesOff = 0
	vacuumVolumeCompactResponseLoadAvg1mOff      = 8
	vacuumVolumeCompactResponseSize              = 12
)

// VacuumVolumeCompactResponse is a zero-copy view into a ZAP-encoded
// VacuumVolumeCompactResponse message.
type VacuumVolumeCompactResponse struct{ o zap.Object }

// WrapVacuumVolumeCompactResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVacuumVolumeCompactResponse(b []byte) (VacuumVolumeCompactResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VacuumVolumeCompactResponse{}, err
	}
	return VacuumVolumeCompactResponse{o: m.Root()}, nil
}

// ProcessedBytes reads the processed_bytes field (proto field 1, int64).
func (t VacuumVolumeCompactResponse) ProcessedBytes() int64 {
	return t.o.Int64(vacuumVolumeCompactResponseProcessedBytesOff)
}

// LoadAvg1m reads the load_avg_1m field (proto field 2, float).
func (t VacuumVolumeCompactResponse) LoadAvg1m() float32 {
	return t.o.Float32(vacuumVolumeCompactResponseLoadAvg1mOff)
}

// VacuumVolumeCompactResponseInput collects the field values for
// NewVacuumVolumeCompactResponse.
type VacuumVolumeCompactResponseInput struct {
	ProcessedBytes int64
	LoadAvg1m      float32
}

// NewVacuumVolumeCompactResponse builds a ZAP-encoded
// VacuumVolumeCompactResponse message from in and returns the bytes.
func NewVacuumVolumeCompactResponse(in VacuumVolumeCompactResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(vacuumVolumeCompactResponseSize)
	ob.SetInt64(vacuumVolumeCompactResponseProcessedBytesOff, in.ProcessedBytes)
	ob.SetFloat32(vacuumVolumeCompactResponseLoadAvg1mOff, in.LoadAvg1m)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VacuumVolumeCommitRequest ---

const (
	vacuumVolumeCommitRequestVolumeIDOff = 0
	vacuumVolumeCommitRequestSize        = 4
)

// VacuumVolumeCommitRequest is a zero-copy view into a ZAP-encoded
// VacuumVolumeCommitRequest message.
type VacuumVolumeCommitRequest struct{ o zap.Object }

// WrapVacuumVolumeCommitRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVacuumVolumeCommitRequest(b []byte) (VacuumVolumeCommitRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VacuumVolumeCommitRequest{}, err
	}
	return VacuumVolumeCommitRequest{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t VacuumVolumeCommitRequest) VolumeID() uint32 {
	return t.o.Uint32(vacuumVolumeCommitRequestVolumeIDOff)
}

// VacuumVolumeCommitRequestInput collects the field values for
// NewVacuumVolumeCommitRequest.
type VacuumVolumeCommitRequestInput struct {
	VolumeID uint32
}

// NewVacuumVolumeCommitRequest builds a ZAP-encoded VacuumVolumeCommitRequest
// message from in and returns the bytes.
func NewVacuumVolumeCommitRequest(in VacuumVolumeCommitRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(vacuumVolumeCommitRequestSize)
	ob.SetUint32(vacuumVolumeCommitRequestVolumeIDOff, in.VolumeID)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VacuumVolumeCommitResponse ---

const (
	vacuumVolumeCommitResponseIsReadOnlyOff = 0
	vacuumVolumeCommitResponseVolumeSizeOff = 8
	vacuumVolumeCommitResponseSize          = 16
)

// VacuumVolumeCommitResponse is a zero-copy view into a ZAP-encoded
// VacuumVolumeCommitResponse message.
type VacuumVolumeCommitResponse struct{ o zap.Object }

// WrapVacuumVolumeCommitResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVacuumVolumeCommitResponse(b []byte) (VacuumVolumeCommitResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VacuumVolumeCommitResponse{}, err
	}
	return VacuumVolumeCommitResponse{o: m.Root()}, nil
}

// IsReadOnly reads the is_read_only field (proto field 1, bool).
func (t VacuumVolumeCommitResponse) IsReadOnly() bool {
	return t.o.Bool(vacuumVolumeCommitResponseIsReadOnlyOff)
}

// VolumeSize reads the volume_size field (proto field 2, uint64).
func (t VacuumVolumeCommitResponse) VolumeSize() uint64 {
	return t.o.Uint64(vacuumVolumeCommitResponseVolumeSizeOff)
}

// VacuumVolumeCommitResponseInput collects the field values for
// NewVacuumVolumeCommitResponse.
type VacuumVolumeCommitResponseInput struct {
	IsReadOnly bool
	VolumeSize uint64
}

// NewVacuumVolumeCommitResponse builds a ZAP-encoded VacuumVolumeCommitResponse
// message from in and returns the bytes.
func NewVacuumVolumeCommitResponse(in VacuumVolumeCommitResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(vacuumVolumeCommitResponseSize)
	ob.SetBool(vacuumVolumeCommitResponseIsReadOnlyOff, in.IsReadOnly)
	ob.SetUint64(vacuumVolumeCommitResponseVolumeSizeOff, in.VolumeSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VacuumVolumeCleanupRequest ---

const (
	vacuumVolumeCleanupRequestVolumeIDOff = 0
	vacuumVolumeCleanupRequestSize        = 4
)

// VacuumVolumeCleanupRequest is a zero-copy view into a ZAP-encoded
// VacuumVolumeCleanupRequest message.
type VacuumVolumeCleanupRequest struct{ o zap.Object }

// WrapVacuumVolumeCleanupRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVacuumVolumeCleanupRequest(b []byte) (VacuumVolumeCleanupRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VacuumVolumeCleanupRequest{}, err
	}
	return VacuumVolumeCleanupRequest{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t VacuumVolumeCleanupRequest) VolumeID() uint32 {
	return t.o.Uint32(vacuumVolumeCleanupRequestVolumeIDOff)
}

// VacuumVolumeCleanupRequestInput collects the field values for
// NewVacuumVolumeCleanupRequest.
type VacuumVolumeCleanupRequestInput struct {
	VolumeID uint32
}

// NewVacuumVolumeCleanupRequest builds a ZAP-encoded VacuumVolumeCleanupRequest
// message from in and returns the bytes.
func NewVacuumVolumeCleanupRequest(in VacuumVolumeCleanupRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(vacuumVolumeCleanupRequestSize)
	ob.SetUint32(vacuumVolumeCleanupRequestVolumeIDOff, in.VolumeID)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VacuumVolumeCleanupResponse ---

const vacuumVolumeCleanupResponseSize = 0

// VacuumVolumeCleanupResponse is a zero-copy view into a ZAP-encoded
// VacuumVolumeCleanupResponse message. VacuumVolumeCleanupResponse carries no
// fields.
type VacuumVolumeCleanupResponse struct{ o zap.Object }

// WrapVacuumVolumeCleanupResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVacuumVolumeCleanupResponse(b []byte) (VacuumVolumeCleanupResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VacuumVolumeCleanupResponse{}, err
	}
	return VacuumVolumeCleanupResponse{o: m.Root()}, nil
}

// VacuumVolumeCleanupResponseInput collects the field values for
// NewVacuumVolumeCleanupResponse. VacuumVolumeCleanupResponse has no fields.
type VacuumVolumeCleanupResponseInput struct{}

// NewVacuumVolumeCleanupResponse builds a ZAP-encoded
// VacuumVolumeCleanupResponse message and returns the bytes.
func NewVacuumVolumeCleanupResponse(in VacuumVolumeCleanupResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(vacuumVolumeCleanupResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}
