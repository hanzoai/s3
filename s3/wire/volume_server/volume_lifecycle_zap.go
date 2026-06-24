// Code generated from volume_server.proto; DO NOT EDIT.

package volume_serverwire

import (
	zap "github.com/zap-proto/go"
)

// --- DeleteCollectionRequest ---

const (
	deleteCollectionRequestCollectionOff = 0
	deleteCollectionRequestSize          = 8
)

// DeleteCollectionRequest is a zero-copy view into a ZAP-encoded
// DeleteCollectionRequest message.
type DeleteCollectionRequest struct{ o zap.Object }

// WrapDeleteCollectionRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapDeleteCollectionRequest(b []byte) (DeleteCollectionRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return DeleteCollectionRequest{}, err
	}
	return DeleteCollectionRequest{o: m.Root()}, nil
}

// Collection reads the collection field (proto field 1, string).
func (t DeleteCollectionRequest) Collection() string {
	return t.o.Text(deleteCollectionRequestCollectionOff)
}

// DeleteCollectionRequestInput collects the field values for
// NewDeleteCollectionRequest.
type DeleteCollectionRequestInput struct {
	Collection string
}

// NewDeleteCollectionRequest builds a ZAP-encoded DeleteCollectionRequest
// message from in and returns the bytes.
func NewDeleteCollectionRequest(in DeleteCollectionRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(deleteCollectionRequestSize)
	ob.SetText(deleteCollectionRequestCollectionOff, in.Collection)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- DeleteCollectionResponse ---

const deleteCollectionResponseSize = 0

// DeleteCollectionResponse is a zero-copy view into a ZAP-encoded
// DeleteCollectionResponse message. DeleteCollectionResponse carries no fields.
type DeleteCollectionResponse struct{ o zap.Object }

// WrapDeleteCollectionResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapDeleteCollectionResponse(b []byte) (DeleteCollectionResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return DeleteCollectionResponse{}, err
	}
	return DeleteCollectionResponse{o: m.Root()}, nil
}

// DeleteCollectionResponseInput collects the field values for
// NewDeleteCollectionResponse. DeleteCollectionResponse has no fields.
type DeleteCollectionResponseInput struct{}

// NewDeleteCollectionResponse builds a ZAP-encoded DeleteCollectionResponse
// message and returns the bytes.
func NewDeleteCollectionResponse(in DeleteCollectionResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(deleteCollectionResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- AllocateVolumeRequest ---

const (
	allocateVolumeRequestVolumeIDOff           = 0
	allocateVolumeRequestCollectionOff         = 8
	allocateVolumeRequestPreallocateOff        = 16
	allocateVolumeRequestReplicationOff        = 24
	allocateVolumeRequestTTLOff                = 32
	allocateVolumeRequestMemoryMapMaxSizeMbOff = 40
	allocateVolumeRequestDiskTypeOff           = 44
	allocateVolumeRequestVersionOff            = 52
	allocateVolumeRequestSize                  = 56
)

// AllocateVolumeRequest is a zero-copy view into a ZAP-encoded
// AllocateVolumeRequest message.
type AllocateVolumeRequest struct{ o zap.Object }

// WrapAllocateVolumeRequest parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapAllocateVolumeRequest(b []byte) (AllocateVolumeRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return AllocateVolumeRequest{}, err
	}
	return AllocateVolumeRequest{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t AllocateVolumeRequest) VolumeID() uint32 {
	return t.o.Uint32(allocateVolumeRequestVolumeIDOff)
}

// Collection reads the collection field (proto field 2, string).
func (t AllocateVolumeRequest) Collection() string {
	return t.o.Text(allocateVolumeRequestCollectionOff)
}

// Preallocate reads the preallocate field (proto field 3, int64).
func (t AllocateVolumeRequest) Preallocate() int64 {
	return t.o.Int64(allocateVolumeRequestPreallocateOff)
}

// Replication reads the replication field (proto field 4, string).
func (t AllocateVolumeRequest) Replication() string {
	return t.o.Text(allocateVolumeRequestReplicationOff)
}

// TTL reads the ttl field (proto field 5, string).
func (t AllocateVolumeRequest) TTL() string { return t.o.Text(allocateVolumeRequestTTLOff) }

// MemoryMapMaxSizeMb reads the memory_map_max_size_mb field (proto field 6,
// uint32).
func (t AllocateVolumeRequest) MemoryMapMaxSizeMb() uint32 {
	return t.o.Uint32(allocateVolumeRequestMemoryMapMaxSizeMbOff)
}

// DiskType reads the disk_type field (proto field 7, string).
func (t AllocateVolumeRequest) DiskType() string {
	return t.o.Text(allocateVolumeRequestDiskTypeOff)
}

// Version reads the version field (proto field 8, uint32).
func (t AllocateVolumeRequest) Version() uint32 {
	return t.o.Uint32(allocateVolumeRequestVersionOff)
}

// AllocateVolumeRequestInput collects the field values for
// NewAllocateVolumeRequest.
type AllocateVolumeRequestInput struct {
	VolumeID           uint32
	Collection         string
	Preallocate        int64
	Replication        string
	TTL                string
	MemoryMapMaxSizeMb uint32
	DiskType           string
	Version            uint32
}

// NewAllocateVolumeRequest builds a ZAP-encoded AllocateVolumeRequest message
// from in and returns the bytes.
func NewAllocateVolumeRequest(in AllocateVolumeRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(allocateVolumeRequestSize)
	ob.SetUint32(allocateVolumeRequestVolumeIDOff, in.VolumeID)
	ob.SetText(allocateVolumeRequestCollectionOff, in.Collection)
	ob.SetInt64(allocateVolumeRequestPreallocateOff, in.Preallocate)
	ob.SetText(allocateVolumeRequestReplicationOff, in.Replication)
	ob.SetText(allocateVolumeRequestTTLOff, in.TTL)
	ob.SetUint32(allocateVolumeRequestMemoryMapMaxSizeMbOff, in.MemoryMapMaxSizeMb)
	ob.SetText(allocateVolumeRequestDiskTypeOff, in.DiskType)
	ob.SetUint32(allocateVolumeRequestVersionOff, in.Version)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- AllocateVolumeResponse ---

const allocateVolumeResponseSize = 0

// AllocateVolumeResponse is a zero-copy view into a ZAP-encoded
// AllocateVolumeResponse message. AllocateVolumeResponse carries no fields.
type AllocateVolumeResponse struct{ o zap.Object }

// WrapAllocateVolumeResponse parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapAllocateVolumeResponse(b []byte) (AllocateVolumeResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return AllocateVolumeResponse{}, err
	}
	return AllocateVolumeResponse{o: m.Root()}, nil
}

// AllocateVolumeResponseInput collects the field values for
// NewAllocateVolumeResponse. AllocateVolumeResponse has no fields.
type AllocateVolumeResponseInput struct{}

// NewAllocateVolumeResponse builds a ZAP-encoded AllocateVolumeResponse message
// and returns the bytes.
func NewAllocateVolumeResponse(in AllocateVolumeResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(allocateVolumeResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeSyncStatusRequest ---

const (
	volumeSyncStatusRequestVolumeIDOff = 0
	volumeSyncStatusRequestSize        = 4
)

// VolumeSyncStatusRequest is a zero-copy view into a ZAP-encoded
// VolumeSyncStatusRequest message.
type VolumeSyncStatusRequest struct{ o zap.Object }

// WrapVolumeSyncStatusRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVolumeSyncStatusRequest(b []byte) (VolumeSyncStatusRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeSyncStatusRequest{}, err
	}
	return VolumeSyncStatusRequest{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t VolumeSyncStatusRequest) VolumeID() uint32 {
	return t.o.Uint32(volumeSyncStatusRequestVolumeIDOff)
}

// VolumeSyncStatusRequestInput collects the field values for
// NewVolumeSyncStatusRequest.
type VolumeSyncStatusRequestInput struct {
	VolumeID uint32
}

// NewVolumeSyncStatusRequest builds a ZAP-encoded VolumeSyncStatusRequest
// message from in and returns the bytes.
func NewVolumeSyncStatusRequest(in VolumeSyncStatusRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeSyncStatusRequestSize)
	ob.SetUint32(volumeSyncStatusRequestVolumeIDOff, in.VolumeID)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeSyncStatusResponse ---
//
// Proto field 3 is absent (skipped in the .proto); field offsets follow proto
// field-number order with the gap preserved as unused space is not allocated —
// each present field gets its own slot in declaration order.

const (
	volumeSyncStatusResponseVolumeIDOff        = 0
	volumeSyncStatusResponseCollectionOff      = 8
	volumeSyncStatusResponseReplicationOff     = 16
	volumeSyncStatusResponseTTLOff             = 24
	volumeSyncStatusResponseTailOffsetOff      = 32
	volumeSyncStatusResponseCompactRevisionOff = 40
	volumeSyncStatusResponseIdxFileSizeOff     = 48
	volumeSyncStatusResponseVersionOff         = 56
	volumeSyncStatusResponseSize               = 60
)

// VolumeSyncStatusResponse is a zero-copy view into a ZAP-encoded
// VolumeSyncStatusResponse message.
type VolumeSyncStatusResponse struct{ o zap.Object }

// WrapVolumeSyncStatusResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVolumeSyncStatusResponse(b []byte) (VolumeSyncStatusResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeSyncStatusResponse{}, err
	}
	return VolumeSyncStatusResponse{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t VolumeSyncStatusResponse) VolumeID() uint32 {
	return t.o.Uint32(volumeSyncStatusResponseVolumeIDOff)
}

// Collection reads the collection field (proto field 2, string).
func (t VolumeSyncStatusResponse) Collection() string {
	return t.o.Text(volumeSyncStatusResponseCollectionOff)
}

// Replication reads the replication field (proto field 4, string).
func (t VolumeSyncStatusResponse) Replication() string {
	return t.o.Text(volumeSyncStatusResponseReplicationOff)
}

// TTL reads the ttl field (proto field 5, string).
func (t VolumeSyncStatusResponse) TTL() string { return t.o.Text(volumeSyncStatusResponseTTLOff) }

// TailOffset reads the tail_offset field (proto field 6, uint64).
func (t VolumeSyncStatusResponse) TailOffset() uint64 {
	return t.o.Uint64(volumeSyncStatusResponseTailOffsetOff)
}

// CompactRevision reads the compact_revision field (proto field 7, uint32).
func (t VolumeSyncStatusResponse) CompactRevision() uint32 {
	return t.o.Uint32(volumeSyncStatusResponseCompactRevisionOff)
}

// IdxFileSize reads the idx_file_size field (proto field 8, uint64).
func (t VolumeSyncStatusResponse) IdxFileSize() uint64 {
	return t.o.Uint64(volumeSyncStatusResponseIdxFileSizeOff)
}

// Version reads the version field (proto field 9, uint32).
func (t VolumeSyncStatusResponse) Version() uint32 {
	return t.o.Uint32(volumeSyncStatusResponseVersionOff)
}

// VolumeSyncStatusResponseInput collects the field values for
// NewVolumeSyncStatusResponse.
type VolumeSyncStatusResponseInput struct {
	VolumeID        uint32
	Collection      string
	Replication     string
	TTL             string
	TailOffset      uint64
	CompactRevision uint32
	IdxFileSize     uint64
	Version         uint32
}

// NewVolumeSyncStatusResponse builds a ZAP-encoded VolumeSyncStatusResponse
// message from in and returns the bytes.
func NewVolumeSyncStatusResponse(in VolumeSyncStatusResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeSyncStatusResponseSize)
	ob.SetUint32(volumeSyncStatusResponseVolumeIDOff, in.VolumeID)
	ob.SetText(volumeSyncStatusResponseCollectionOff, in.Collection)
	ob.SetText(volumeSyncStatusResponseReplicationOff, in.Replication)
	ob.SetText(volumeSyncStatusResponseTTLOff, in.TTL)
	ob.SetUint64(volumeSyncStatusResponseTailOffsetOff, in.TailOffset)
	ob.SetUint32(volumeSyncStatusResponseCompactRevisionOff, in.CompactRevision)
	ob.SetUint64(volumeSyncStatusResponseIdxFileSizeOff, in.IdxFileSize)
	ob.SetUint32(volumeSyncStatusResponseVersionOff, in.Version)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeIncrementalCopyRequest ---

const (
	volumeIncrementalCopyRequestVolumeIDOff = 0
	volumeIncrementalCopyRequestSinceNsOff  = 8
	volumeIncrementalCopyRequestSize        = 16
)

// VolumeIncrementalCopyRequest is a zero-copy view into a ZAP-encoded
// VolumeIncrementalCopyRequest message.
type VolumeIncrementalCopyRequest struct{ o zap.Object }

// WrapVolumeIncrementalCopyRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVolumeIncrementalCopyRequest(b []byte) (VolumeIncrementalCopyRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeIncrementalCopyRequest{}, err
	}
	return VolumeIncrementalCopyRequest{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t VolumeIncrementalCopyRequest) VolumeID() uint32 {
	return t.o.Uint32(volumeIncrementalCopyRequestVolumeIDOff)
}

// SinceNs reads the since_ns field (proto field 2, uint64).
func (t VolumeIncrementalCopyRequest) SinceNs() uint64 {
	return t.o.Uint64(volumeIncrementalCopyRequestSinceNsOff)
}

// VolumeIncrementalCopyRequestInput collects the field values for
// NewVolumeIncrementalCopyRequest.
type VolumeIncrementalCopyRequestInput struct {
	VolumeID uint32
	SinceNs  uint64
}

// NewVolumeIncrementalCopyRequest builds a ZAP-encoded
// VolumeIncrementalCopyRequest message from in and returns the bytes.
func NewVolumeIncrementalCopyRequest(in VolumeIncrementalCopyRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeIncrementalCopyRequestSize)
	ob.SetUint32(volumeIncrementalCopyRequestVolumeIDOff, in.VolumeID)
	ob.SetUint64(volumeIncrementalCopyRequestSinceNsOff, in.SinceNs)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeIncrementalCopyResponse ---

const (
	volumeIncrementalCopyResponseFileContentOff = 0
	volumeIncrementalCopyResponseSize           = 8
)

// VolumeIncrementalCopyResponse is a zero-copy view into a ZAP-encoded
// VolumeIncrementalCopyResponse message.
type VolumeIncrementalCopyResponse struct{ o zap.Object }

// WrapVolumeIncrementalCopyResponse parses b and returns a typed view. Returns
// an error if the wire-level checks (magic, version, size) fail.
func WrapVolumeIncrementalCopyResponse(b []byte) (VolumeIncrementalCopyResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeIncrementalCopyResponse{}, err
	}
	return VolumeIncrementalCopyResponse{o: m.Root()}, nil
}

// FileContent reads the file_content field (proto field 1, bytes).
func (t VolumeIncrementalCopyResponse) FileContent() []byte {
	return t.o.Bytes(volumeIncrementalCopyResponseFileContentOff)
}

// VolumeIncrementalCopyResponseInput collects the field values for
// NewVolumeIncrementalCopyResponse.
type VolumeIncrementalCopyResponseInput struct {
	FileContent []byte
}

// NewVolumeIncrementalCopyResponse builds a ZAP-encoded
// VolumeIncrementalCopyResponse message from in and returns the bytes.
func NewVolumeIncrementalCopyResponse(in VolumeIncrementalCopyResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeIncrementalCopyResponseSize)
	ob.SetBytes(volumeIncrementalCopyResponseFileContentOff, in.FileContent)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeMountRequest ---

const (
	volumeMountRequestVolumeIDOff = 0
	volumeMountRequestSize        = 4
)

// VolumeMountRequest is a zero-copy view into a ZAP-encoded VolumeMountRequest
// message.
type VolumeMountRequest struct{ o zap.Object }

// WrapVolumeMountRequest parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapVolumeMountRequest(b []byte) (VolumeMountRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeMountRequest{}, err
	}
	return VolumeMountRequest{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t VolumeMountRequest) VolumeID() uint32 { return t.o.Uint32(volumeMountRequestVolumeIDOff) }

// VolumeMountRequestInput collects the field values for NewVolumeMountRequest.
type VolumeMountRequestInput struct {
	VolumeID uint32
}

// NewVolumeMountRequest builds a ZAP-encoded VolumeMountRequest message from in
// and returns the bytes.
func NewVolumeMountRequest(in VolumeMountRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeMountRequestSize)
	ob.SetUint32(volumeMountRequestVolumeIDOff, in.VolumeID)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeMountResponse ---

const volumeMountResponseSize = 0

// VolumeMountResponse is a zero-copy view into a ZAP-encoded VolumeMountResponse
// message. VolumeMountResponse carries no fields.
type VolumeMountResponse struct{ o zap.Object }

// WrapVolumeMountResponse parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapVolumeMountResponse(b []byte) (VolumeMountResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeMountResponse{}, err
	}
	return VolumeMountResponse{o: m.Root()}, nil
}

// VolumeMountResponseInput collects the field values for NewVolumeMountResponse.
// VolumeMountResponse has no fields.
type VolumeMountResponseInput struct{}

// NewVolumeMountResponse builds a ZAP-encoded VolumeMountResponse message and
// returns the bytes.
func NewVolumeMountResponse(in VolumeMountResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeMountResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeUnmountRequest ---

const (
	volumeUnmountRequestVolumeIDOff = 0
	volumeUnmountRequestSize        = 4
)

// VolumeUnmountRequest is a zero-copy view into a ZAP-encoded
// VolumeUnmountRequest message.
type VolumeUnmountRequest struct{ o zap.Object }

// WrapVolumeUnmountRequest parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapVolumeUnmountRequest(b []byte) (VolumeUnmountRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeUnmountRequest{}, err
	}
	return VolumeUnmountRequest{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t VolumeUnmountRequest) VolumeID() uint32 { return t.o.Uint32(volumeUnmountRequestVolumeIDOff) }

// VolumeUnmountRequestInput collects the field values for
// NewVolumeUnmountRequest.
type VolumeUnmountRequestInput struct {
	VolumeID uint32
}

// NewVolumeUnmountRequest builds a ZAP-encoded VolumeUnmountRequest message from
// in and returns the bytes.
func NewVolumeUnmountRequest(in VolumeUnmountRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeUnmountRequestSize)
	ob.SetUint32(volumeUnmountRequestVolumeIDOff, in.VolumeID)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeUnmountResponse ---

const volumeUnmountResponseSize = 0

// VolumeUnmountResponse is a zero-copy view into a ZAP-encoded
// VolumeUnmountResponse message. VolumeUnmountResponse carries no fields.
type VolumeUnmountResponse struct{ o zap.Object }

// WrapVolumeUnmountResponse parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapVolumeUnmountResponse(b []byte) (VolumeUnmountResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeUnmountResponse{}, err
	}
	return VolumeUnmountResponse{o: m.Root()}, nil
}

// VolumeUnmountResponseInput collects the field values for
// NewVolumeUnmountResponse. VolumeUnmountResponse has no fields.
type VolumeUnmountResponseInput struct{}

// NewVolumeUnmountResponse builds a ZAP-encoded VolumeUnmountResponse message
// and returns the bytes.
func NewVolumeUnmountResponse(in VolumeUnmountResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeUnmountResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeDeleteRequest ---

const (
	volumeDeleteRequestVolumeIDOff       = 0
	volumeDeleteRequestOnlyEmptyOff      = 4
	volumeDeleteRequestKeepRemoteDataOff = 5
	volumeDeleteRequestSize              = 6
)

// VolumeDeleteRequest is a zero-copy view into a ZAP-encoded VolumeDeleteRequest
// message.
type VolumeDeleteRequest struct{ o zap.Object }

// WrapVolumeDeleteRequest parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapVolumeDeleteRequest(b []byte) (VolumeDeleteRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeDeleteRequest{}, err
	}
	return VolumeDeleteRequest{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t VolumeDeleteRequest) VolumeID() uint32 { return t.o.Uint32(volumeDeleteRequestVolumeIDOff) }

// OnlyEmpty reads the only_empty field (proto field 2, bool).
func (t VolumeDeleteRequest) OnlyEmpty() bool { return t.o.Bool(volumeDeleteRequestOnlyEmptyOff) }

// KeepRemoteData reads the keep_remote_data field (proto field 3, bool).
func (t VolumeDeleteRequest) KeepRemoteData() bool {
	return t.o.Bool(volumeDeleteRequestKeepRemoteDataOff)
}

// VolumeDeleteRequestInput collects the field values for NewVolumeDeleteRequest.
type VolumeDeleteRequestInput struct {
	VolumeID       uint32
	OnlyEmpty      bool
	KeepRemoteData bool
}

// NewVolumeDeleteRequest builds a ZAP-encoded VolumeDeleteRequest message from
// in and returns the bytes.
func NewVolumeDeleteRequest(in VolumeDeleteRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeDeleteRequestSize)
	ob.SetUint32(volumeDeleteRequestVolumeIDOff, in.VolumeID)
	ob.SetBool(volumeDeleteRequestOnlyEmptyOff, in.OnlyEmpty)
	ob.SetBool(volumeDeleteRequestKeepRemoteDataOff, in.KeepRemoteData)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeDeleteResponse ---

const volumeDeleteResponseSize = 0

// VolumeDeleteResponse is a zero-copy view into a ZAP-encoded
// VolumeDeleteResponse message. VolumeDeleteResponse carries no fields.
type VolumeDeleteResponse struct{ o zap.Object }

// WrapVolumeDeleteResponse parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapVolumeDeleteResponse(b []byte) (VolumeDeleteResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeDeleteResponse{}, err
	}
	return VolumeDeleteResponse{o: m.Root()}, nil
}

// VolumeDeleteResponseInput collects the field values for
// NewVolumeDeleteResponse. VolumeDeleteResponse has no fields.
type VolumeDeleteResponseInput struct{}

// NewVolumeDeleteResponse builds a ZAP-encoded VolumeDeleteResponse message and
// returns the bytes.
func NewVolumeDeleteResponse(in VolumeDeleteResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeDeleteResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeMarkReadonlyRequest ---

const (
	volumeMarkReadonlyRequestVolumeIDOff = 0
	volumeMarkReadonlyRequestPersistOff  = 4
	volumeMarkReadonlyRequestSize        = 5
)

// VolumeMarkReadonlyRequest is a zero-copy view into a ZAP-encoded
// VolumeMarkReadonlyRequest message.
type VolumeMarkReadonlyRequest struct{ o zap.Object }

// WrapVolumeMarkReadonlyRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVolumeMarkReadonlyRequest(b []byte) (VolumeMarkReadonlyRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeMarkReadonlyRequest{}, err
	}
	return VolumeMarkReadonlyRequest{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t VolumeMarkReadonlyRequest) VolumeID() uint32 {
	return t.o.Uint32(volumeMarkReadonlyRequestVolumeIDOff)
}

// Persist reads the persist field (proto field 2, bool).
func (t VolumeMarkReadonlyRequest) Persist() bool {
	return t.o.Bool(volumeMarkReadonlyRequestPersistOff)
}

// VolumeMarkReadonlyRequestInput collects the field values for
// NewVolumeMarkReadonlyRequest.
type VolumeMarkReadonlyRequestInput struct {
	VolumeID uint32
	Persist  bool
}

// NewVolumeMarkReadonlyRequest builds a ZAP-encoded VolumeMarkReadonlyRequest
// message from in and returns the bytes.
func NewVolumeMarkReadonlyRequest(in VolumeMarkReadonlyRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeMarkReadonlyRequestSize)
	ob.SetUint32(volumeMarkReadonlyRequestVolumeIDOff, in.VolumeID)
	ob.SetBool(volumeMarkReadonlyRequestPersistOff, in.Persist)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeMarkReadonlyResponse ---

const volumeMarkReadonlyResponseSize = 0

// VolumeMarkReadonlyResponse is a zero-copy view into a ZAP-encoded
// VolumeMarkReadonlyResponse message. VolumeMarkReadonlyResponse carries no
// fields.
type VolumeMarkReadonlyResponse struct{ o zap.Object }

// WrapVolumeMarkReadonlyResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVolumeMarkReadonlyResponse(b []byte) (VolumeMarkReadonlyResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeMarkReadonlyResponse{}, err
	}
	return VolumeMarkReadonlyResponse{o: m.Root()}, nil
}

// VolumeMarkReadonlyResponseInput collects the field values for
// NewVolumeMarkReadonlyResponse. VolumeMarkReadonlyResponse has no fields.
type VolumeMarkReadonlyResponseInput struct{}

// NewVolumeMarkReadonlyResponse builds a ZAP-encoded VolumeMarkReadonlyResponse
// message and returns the bytes.
func NewVolumeMarkReadonlyResponse(in VolumeMarkReadonlyResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeMarkReadonlyResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeMarkWritableRequest ---

const (
	volumeMarkWritableRequestVolumeIDOff = 0
	volumeMarkWritableRequestSize        = 4
)

// VolumeMarkWritableRequest is a zero-copy view into a ZAP-encoded
// VolumeMarkWritableRequest message.
type VolumeMarkWritableRequest struct{ o zap.Object }

// WrapVolumeMarkWritableRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVolumeMarkWritableRequest(b []byte) (VolumeMarkWritableRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeMarkWritableRequest{}, err
	}
	return VolumeMarkWritableRequest{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t VolumeMarkWritableRequest) VolumeID() uint32 {
	return t.o.Uint32(volumeMarkWritableRequestVolumeIDOff)
}

// VolumeMarkWritableRequestInput collects the field values for
// NewVolumeMarkWritableRequest.
type VolumeMarkWritableRequestInput struct {
	VolumeID uint32
}

// NewVolumeMarkWritableRequest builds a ZAP-encoded VolumeMarkWritableRequest
// message from in and returns the bytes.
func NewVolumeMarkWritableRequest(in VolumeMarkWritableRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeMarkWritableRequestSize)
	ob.SetUint32(volumeMarkWritableRequestVolumeIDOff, in.VolumeID)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeMarkWritableResponse ---

const volumeMarkWritableResponseSize = 0

// VolumeMarkWritableResponse is a zero-copy view into a ZAP-encoded
// VolumeMarkWritableResponse message. VolumeMarkWritableResponse carries no
// fields.
type VolumeMarkWritableResponse struct{ o zap.Object }

// WrapVolumeMarkWritableResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVolumeMarkWritableResponse(b []byte) (VolumeMarkWritableResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeMarkWritableResponse{}, err
	}
	return VolumeMarkWritableResponse{o: m.Root()}, nil
}

// VolumeMarkWritableResponseInput collects the field values for
// NewVolumeMarkWritableResponse. VolumeMarkWritableResponse has no fields.
type VolumeMarkWritableResponseInput struct{}

// NewVolumeMarkWritableResponse builds a ZAP-encoded VolumeMarkWritableResponse
// message and returns the bytes.
func NewVolumeMarkWritableResponse(in VolumeMarkWritableResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeMarkWritableResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeConfigureRequest ---

const (
	volumeConfigureRequestVolumeIDOff    = 0
	volumeConfigureRequestReplicationOff = 8
	volumeConfigureRequestSize           = 16
)

// VolumeConfigureRequest is a zero-copy view into a ZAP-encoded
// VolumeConfigureRequest message.
type VolumeConfigureRequest struct{ o zap.Object }

// WrapVolumeConfigureRequest parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapVolumeConfigureRequest(b []byte) (VolumeConfigureRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeConfigureRequest{}, err
	}
	return VolumeConfigureRequest{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t VolumeConfigureRequest) VolumeID() uint32 {
	return t.o.Uint32(volumeConfigureRequestVolumeIDOff)
}

// Replication reads the replication field (proto field 2, string).
func (t VolumeConfigureRequest) Replication() string {
	return t.o.Text(volumeConfigureRequestReplicationOff)
}

// VolumeConfigureRequestInput collects the field values for
// NewVolumeConfigureRequest.
type VolumeConfigureRequestInput struct {
	VolumeID    uint32
	Replication string
}

// NewVolumeConfigureRequest builds a ZAP-encoded VolumeConfigureRequest message
// from in and returns the bytes.
func NewVolumeConfigureRequest(in VolumeConfigureRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeConfigureRequestSize)
	ob.SetUint32(volumeConfigureRequestVolumeIDOff, in.VolumeID)
	ob.SetText(volumeConfigureRequestReplicationOff, in.Replication)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeConfigureResponse ---

const (
	volumeConfigureResponseErrorOff = 0
	volumeConfigureResponseSize     = 8
)

// VolumeConfigureResponse is a zero-copy view into a ZAP-encoded
// VolumeConfigureResponse message.
type VolumeConfigureResponse struct{ o zap.Object }

// WrapVolumeConfigureResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVolumeConfigureResponse(b []byte) (VolumeConfigureResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeConfigureResponse{}, err
	}
	return VolumeConfigureResponse{o: m.Root()}, nil
}

// Error reads the error field (proto field 1, string).
func (t VolumeConfigureResponse) Error() string { return t.o.Text(volumeConfigureResponseErrorOff) }

// VolumeConfigureResponseInput collects the field values for
// NewVolumeConfigureResponse.
type VolumeConfigureResponseInput struct {
	Error string
}

// NewVolumeConfigureResponse builds a ZAP-encoded VolumeConfigureResponse
// message from in and returns the bytes.
func NewVolumeConfigureResponse(in VolumeConfigureResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeConfigureResponseSize)
	ob.SetText(volumeConfigureResponseErrorOff, in.Error)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeStatusRequest ---

const (
	volumeStatusRequestVolumeIDOff = 0
	volumeStatusRequestSize        = 4
)

// VolumeStatusRequest is a zero-copy view into a ZAP-encoded VolumeStatusRequest
// message.
type VolumeStatusRequest struct{ o zap.Object }

// WrapVolumeStatusRequest parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapVolumeStatusRequest(b []byte) (VolumeStatusRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeStatusRequest{}, err
	}
	return VolumeStatusRequest{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t VolumeStatusRequest) VolumeID() uint32 { return t.o.Uint32(volumeStatusRequestVolumeIDOff) }

// VolumeStatusRequestInput collects the field values for NewVolumeStatusRequest.
type VolumeStatusRequestInput struct {
	VolumeID uint32
}

// NewVolumeStatusRequest builds a ZAP-encoded VolumeStatusRequest message from
// in and returns the bytes.
func NewVolumeStatusRequest(in VolumeStatusRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeStatusRequestSize)
	ob.SetUint32(volumeStatusRequestVolumeIDOff, in.VolumeID)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeStatusResponse ---

const (
	volumeStatusResponseIsReadOnlyOff       = 0
	volumeStatusResponseVolumeSizeOff       = 8
	volumeStatusResponseFileCountOff        = 16
	volumeStatusResponseFileDeletedCountOff = 24
	volumeStatusResponseSize                = 32
)

// VolumeStatusResponse is a zero-copy view into a ZAP-encoded
// VolumeStatusResponse message.
type VolumeStatusResponse struct{ o zap.Object }

// WrapVolumeStatusResponse parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapVolumeStatusResponse(b []byte) (VolumeStatusResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeStatusResponse{}, err
	}
	return VolumeStatusResponse{o: m.Root()}, nil
}

// IsReadOnly reads the is_read_only field (proto field 1, bool).
func (t VolumeStatusResponse) IsReadOnly() bool {
	return t.o.Bool(volumeStatusResponseIsReadOnlyOff)
}

// VolumeSize reads the volume_size field (proto field 2, uint64).
func (t VolumeStatusResponse) VolumeSize() uint64 {
	return t.o.Uint64(volumeStatusResponseVolumeSizeOff)
}

// FileCount reads the file_count field (proto field 3, uint64).
func (t VolumeStatusResponse) FileCount() uint64 {
	return t.o.Uint64(volumeStatusResponseFileCountOff)
}

// FileDeletedCount reads the file_deleted_count field (proto field 4, uint64).
func (t VolumeStatusResponse) FileDeletedCount() uint64 {
	return t.o.Uint64(volumeStatusResponseFileDeletedCountOff)
}

// VolumeStatusResponseInput collects the field values for
// NewVolumeStatusResponse.
type VolumeStatusResponseInput struct {
	IsReadOnly       bool
	VolumeSize       uint64
	FileCount        uint64
	FileDeletedCount uint64
}

// NewVolumeStatusResponse builds a ZAP-encoded VolumeStatusResponse message from
// in and returns the bytes.
func NewVolumeStatusResponse(in VolumeStatusResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeStatusResponseSize)
	ob.SetBool(volumeStatusResponseIsReadOnlyOff, in.IsReadOnly)
	ob.SetUint64(volumeStatusResponseVolumeSizeOff, in.VolumeSize)
	ob.SetUint64(volumeStatusResponseFileCountOff, in.FileCount)
	ob.SetUint64(volumeStatusResponseFileDeletedCountOff, in.FileDeletedCount)
	ob.FinishAsRoot()
	return b.Finish()
}
