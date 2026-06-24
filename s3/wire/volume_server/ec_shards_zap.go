// Code generated from volume_server.proto; DO NOT EDIT.

package volume_serverwire

import (
	zap "github.com/zap-proto/go"
)

// --- VolumeEcShardsGenerateRequest ---

const (
	volumeEcShardsGenerateRequestVolumeIDOff   = 0
	volumeEcShardsGenerateRequestCollectionOff = 8
	volumeEcShardsGenerateRequestSize          = 16
)

// VolumeEcShardsGenerateRequest is a zero-copy view into a ZAP-encoded
// VolumeEcShardsGenerateRequest message.
type VolumeEcShardsGenerateRequest struct{ o zap.Object }

// WrapVolumeEcShardsGenerateRequest parses b and returns a typed view. Returns
// an error if the wire-level checks (magic, version, size) fail.
func WrapVolumeEcShardsGenerateRequest(b []byte) (VolumeEcShardsGenerateRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeEcShardsGenerateRequest{}, err
	}
	return VolumeEcShardsGenerateRequest{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t VolumeEcShardsGenerateRequest) VolumeID() uint32 {
	return t.o.Uint32(volumeEcShardsGenerateRequestVolumeIDOff)
}

// Collection reads the collection field (proto field 2, string).
func (t VolumeEcShardsGenerateRequest) Collection() string {
	return t.o.Text(volumeEcShardsGenerateRequestCollectionOff)
}

// VolumeEcShardsGenerateRequestInput collects the field values for
// NewVolumeEcShardsGenerateRequest.
type VolumeEcShardsGenerateRequestInput struct {
	VolumeID   uint32
	Collection string
}

// NewVolumeEcShardsGenerateRequest builds a ZAP-encoded
// VolumeEcShardsGenerateRequest message from in and returns the bytes.
func NewVolumeEcShardsGenerateRequest(in VolumeEcShardsGenerateRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeEcShardsGenerateRequestSize)
	ob.SetUint32(volumeEcShardsGenerateRequestVolumeIDOff, in.VolumeID)
	ob.SetText(volumeEcShardsGenerateRequestCollectionOff, in.Collection)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeEcShardsGenerateResponse ---

const volumeEcShardsGenerateResponseSize = 0

// VolumeEcShardsGenerateResponse is a zero-copy view into a ZAP-encoded
// VolumeEcShardsGenerateResponse message. VolumeEcShardsGenerateResponse carries
// no fields.
type VolumeEcShardsGenerateResponse struct{ o zap.Object }

// WrapVolumeEcShardsGenerateResponse parses b and returns a typed view. Returns
// an error if the wire-level checks (magic, version, size) fail.
func WrapVolumeEcShardsGenerateResponse(b []byte) (VolumeEcShardsGenerateResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeEcShardsGenerateResponse{}, err
	}
	return VolumeEcShardsGenerateResponse{o: m.Root()}, nil
}

// VolumeEcShardsGenerateResponseInput collects the field values for
// NewVolumeEcShardsGenerateResponse. VolumeEcShardsGenerateResponse has no
// fields.
type VolumeEcShardsGenerateResponseInput struct{}

// NewVolumeEcShardsGenerateResponse builds a ZAP-encoded
// VolumeEcShardsGenerateResponse message and returns the bytes.
func NewVolumeEcShardsGenerateResponse(in VolumeEcShardsGenerateResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeEcShardsGenerateResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeEcShardsRebuildRequest ---

const (
	volumeEcShardsRebuildRequestVolumeIDOff            = 0
	volumeEcShardsRebuildRequestCollectionOff          = 8
	volumeEcShardsRebuildRequestUnsafeIgnoreSidecarOff = 16
	volumeEcShardsRebuildRequestSize                   = 17
)

// VolumeEcShardsRebuildRequest is a zero-copy view into a ZAP-encoded
// VolumeEcShardsRebuildRequest message.
type VolumeEcShardsRebuildRequest struct{ o zap.Object }

// WrapVolumeEcShardsRebuildRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVolumeEcShardsRebuildRequest(b []byte) (VolumeEcShardsRebuildRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeEcShardsRebuildRequest{}, err
	}
	return VolumeEcShardsRebuildRequest{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t VolumeEcShardsRebuildRequest) VolumeID() uint32 {
	return t.o.Uint32(volumeEcShardsRebuildRequestVolumeIDOff)
}

// Collection reads the collection field (proto field 2, string).
func (t VolumeEcShardsRebuildRequest) Collection() string {
	return t.o.Text(volumeEcShardsRebuildRequestCollectionOff)
}

// UnsafeIgnoreSidecar reads the unsafe_ignore_sidecar field (proto field 3,
// bool).
func (t VolumeEcShardsRebuildRequest) UnsafeIgnoreSidecar() bool {
	return t.o.Bool(volumeEcShardsRebuildRequestUnsafeIgnoreSidecarOff)
}

// VolumeEcShardsRebuildRequestInput collects the field values for
// NewVolumeEcShardsRebuildRequest.
type VolumeEcShardsRebuildRequestInput struct {
	VolumeID            uint32
	Collection          string
	UnsafeIgnoreSidecar bool
}

// NewVolumeEcShardsRebuildRequest builds a ZAP-encoded
// VolumeEcShardsRebuildRequest message from in and returns the bytes.
func NewVolumeEcShardsRebuildRequest(in VolumeEcShardsRebuildRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeEcShardsRebuildRequestSize)
	ob.SetUint32(volumeEcShardsRebuildRequestVolumeIDOff, in.VolumeID)
	ob.SetText(volumeEcShardsRebuildRequestCollectionOff, in.Collection)
	ob.SetBool(volumeEcShardsRebuildRequestUnsafeIgnoreSidecarOff, in.UnsafeIgnoreSidecar)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeEcShardsRebuildResponse ---

const (
	volumeEcShardsRebuildResponseRebuiltShardIDsOff = 0
	volumeEcShardsRebuildResponseSize               = 8
)

// VolumeEcShardsRebuildResponse is a zero-copy view into a ZAP-encoded
// VolumeEcShardsRebuildResponse message.
type VolumeEcShardsRebuildResponse struct{ o zap.Object }

// WrapVolumeEcShardsRebuildResponse parses b and returns a typed view. Returns
// an error if the wire-level checks (magic, version, size) fail.
func WrapVolumeEcShardsRebuildResponse(b []byte) (VolumeEcShardsRebuildResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeEcShardsRebuildResponse{}, err
	}
	return VolumeEcShardsRebuildResponse{o: m.Root()}, nil
}

// RebuiltShardIDsLen reports the number of rebuilt_shard_ids elements (proto
// field 1, repeated uint32).
func (t VolumeEcShardsRebuildResponse) RebuiltShardIDsLen() int {
	return t.o.List(volumeEcShardsRebuildResponseRebuiltShardIDsOff).Len()
}

// RebuiltShardIDAt returns the i-th rebuilt_shard_ids element (uint32).
func (t VolumeEcShardsRebuildResponse) RebuiltShardIDAt(i int) uint32 {
	return t.o.List(volumeEcShardsRebuildResponseRebuiltShardIDsOff).Uint32(i)
}

// VolumeEcShardsRebuildResponseInput collects the field values for
// NewVolumeEcShardsRebuildResponse.
type VolumeEcShardsRebuildResponseInput struct {
	RebuiltShardIDs []uint32
}

// NewVolumeEcShardsRebuildResponse builds a ZAP-encoded
// VolumeEcShardsRebuildResponse message from in and returns the bytes.
func NewVolumeEcShardsRebuildResponse(in VolumeEcShardsRebuildResponseInput) []byte {
	b := zap.NewBuilder(256)
	listOff, listLen := addUint32List(b, in.RebuiltShardIDs)
	ob := b.StartObject(volumeEcShardsRebuildResponseSize)
	ob.SetList(volumeEcShardsRebuildResponseRebuiltShardIDsOff, listOff, listLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeEcShardsCopyRequest ---

const (
	volumeEcShardsCopyRequestVolumeIDOff       = 0
	volumeEcShardsCopyRequestCollectionOff     = 8
	volumeEcShardsCopyRequestShardIDsOff       = 16
	volumeEcShardsCopyRequestCopyEcxFileOff    = 24
	volumeEcShardsCopyRequestSourceDataNodeOff = 32
	volumeEcShardsCopyRequestCopyEcjFileOff    = 40
	volumeEcShardsCopyRequestCopyVifFileOff    = 41
	volumeEcShardsCopyRequestDiskIDOff         = 44
	volumeEcShardsCopyRequestCopyEcsumFileOff  = 48
	volumeEcShardsCopyRequestSize              = 49
)

// VolumeEcShardsCopyRequest is a zero-copy view into a ZAP-encoded
// VolumeEcShardsCopyRequest message.
type VolumeEcShardsCopyRequest struct{ o zap.Object }

// WrapVolumeEcShardsCopyRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVolumeEcShardsCopyRequest(b []byte) (VolumeEcShardsCopyRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeEcShardsCopyRequest{}, err
	}
	return VolumeEcShardsCopyRequest{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t VolumeEcShardsCopyRequest) VolumeID() uint32 {
	return t.o.Uint32(volumeEcShardsCopyRequestVolumeIDOff)
}

// Collection reads the collection field (proto field 2, string).
func (t VolumeEcShardsCopyRequest) Collection() string {
	return t.o.Text(volumeEcShardsCopyRequestCollectionOff)
}

// ShardIDsLen reports the number of shard_ids elements (proto field 3, repeated
// uint32).
func (t VolumeEcShardsCopyRequest) ShardIDsLen() int {
	return t.o.List(volumeEcShardsCopyRequestShardIDsOff).Len()
}

// ShardIDAt returns the i-th shard_ids element (uint32).
func (t VolumeEcShardsCopyRequest) ShardIDAt(i int) uint32 {
	return t.o.List(volumeEcShardsCopyRequestShardIDsOff).Uint32(i)
}

// CopyEcxFile reads the copy_ecx_file field (proto field 4, bool).
func (t VolumeEcShardsCopyRequest) CopyEcxFile() bool {
	return t.o.Bool(volumeEcShardsCopyRequestCopyEcxFileOff)
}

// SourceDataNode reads the source_data_node field (proto field 5, string).
func (t VolumeEcShardsCopyRequest) SourceDataNode() string {
	return t.o.Text(volumeEcShardsCopyRequestSourceDataNodeOff)
}

// CopyEcjFile reads the copy_ecj_file field (proto field 6, bool).
func (t VolumeEcShardsCopyRequest) CopyEcjFile() bool {
	return t.o.Bool(volumeEcShardsCopyRequestCopyEcjFileOff)
}

// CopyVifFile reads the copy_vif_file field (proto field 7, bool).
func (t VolumeEcShardsCopyRequest) CopyVifFile() bool {
	return t.o.Bool(volumeEcShardsCopyRequestCopyVifFileOff)
}

// DiskID reads the disk_id field (proto field 8, uint32).
func (t VolumeEcShardsCopyRequest) DiskID() uint32 {
	return t.o.Uint32(volumeEcShardsCopyRequestDiskIDOff)
}

// CopyEcsumFile reads the copy_ecsum_file field (proto field 9, bool).
func (t VolumeEcShardsCopyRequest) CopyEcsumFile() bool {
	return t.o.Bool(volumeEcShardsCopyRequestCopyEcsumFileOff)
}

// VolumeEcShardsCopyRequestInput collects the field values for
// NewVolumeEcShardsCopyRequest.
type VolumeEcShardsCopyRequestInput struct {
	VolumeID       uint32
	Collection     string
	ShardIDs       []uint32
	CopyEcxFile    bool
	SourceDataNode string
	CopyEcjFile    bool
	CopyVifFile    bool
	DiskID         uint32
	CopyEcsumFile  bool
}

// NewVolumeEcShardsCopyRequest builds a ZAP-encoded VolumeEcShardsCopyRequest
// message from in and returns the bytes.
func NewVolumeEcShardsCopyRequest(in VolumeEcShardsCopyRequestInput) []byte {
	b := zap.NewBuilder(256)
	listOff, listLen := addUint32List(b, in.ShardIDs)
	ob := b.StartObject(volumeEcShardsCopyRequestSize)
	ob.SetUint32(volumeEcShardsCopyRequestVolumeIDOff, in.VolumeID)
	ob.SetText(volumeEcShardsCopyRequestCollectionOff, in.Collection)
	ob.SetList(volumeEcShardsCopyRequestShardIDsOff, listOff, listLen)
	ob.SetBool(volumeEcShardsCopyRequestCopyEcxFileOff, in.CopyEcxFile)
	ob.SetText(volumeEcShardsCopyRequestSourceDataNodeOff, in.SourceDataNode)
	ob.SetBool(volumeEcShardsCopyRequestCopyEcjFileOff, in.CopyEcjFile)
	ob.SetBool(volumeEcShardsCopyRequestCopyVifFileOff, in.CopyVifFile)
	ob.SetUint32(volumeEcShardsCopyRequestDiskIDOff, in.DiskID)
	ob.SetBool(volumeEcShardsCopyRequestCopyEcsumFileOff, in.CopyEcsumFile)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeEcShardsCopyResponse ---

const volumeEcShardsCopyResponseSize = 0

// VolumeEcShardsCopyResponse is a zero-copy view into a ZAP-encoded
// VolumeEcShardsCopyResponse message. VolumeEcShardsCopyResponse carries no
// fields.
type VolumeEcShardsCopyResponse struct{ o zap.Object }

// WrapVolumeEcShardsCopyResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVolumeEcShardsCopyResponse(b []byte) (VolumeEcShardsCopyResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeEcShardsCopyResponse{}, err
	}
	return VolumeEcShardsCopyResponse{o: m.Root()}, nil
}

// VolumeEcShardsCopyResponseInput collects the field values for
// NewVolumeEcShardsCopyResponse. VolumeEcShardsCopyResponse has no fields.
type VolumeEcShardsCopyResponseInput struct{}

// NewVolumeEcShardsCopyResponse builds a ZAP-encoded VolumeEcShardsCopyResponse
// message and returns the bytes.
func NewVolumeEcShardsCopyResponse(in VolumeEcShardsCopyResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeEcShardsCopyResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeEcShardsDeleteRequest ---

const (
	volumeEcShardsDeleteRequestVolumeIDOff     = 0
	volumeEcShardsDeleteRequestCollectionOff   = 8
	volumeEcShardsDeleteRequestShardIDsOff     = 16
	volumeEcShardsDeleteRequestFullTeardownOff = 24
	volumeEcShardsDeleteRequestEncodeTsNsOff   = 32
	volumeEcShardsDeleteRequestSize            = 40
)

// VolumeEcShardsDeleteRequest is a zero-copy view into a ZAP-encoded
// VolumeEcShardsDeleteRequest message.
type VolumeEcShardsDeleteRequest struct{ o zap.Object }

// WrapVolumeEcShardsDeleteRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVolumeEcShardsDeleteRequest(b []byte) (VolumeEcShardsDeleteRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeEcShardsDeleteRequest{}, err
	}
	return VolumeEcShardsDeleteRequest{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t VolumeEcShardsDeleteRequest) VolumeID() uint32 {
	return t.o.Uint32(volumeEcShardsDeleteRequestVolumeIDOff)
}

// Collection reads the collection field (proto field 2, string).
func (t VolumeEcShardsDeleteRequest) Collection() string {
	return t.o.Text(volumeEcShardsDeleteRequestCollectionOff)
}

// ShardIDsLen reports the number of shard_ids elements (proto field 3, repeated
// uint32).
func (t VolumeEcShardsDeleteRequest) ShardIDsLen() int {
	return t.o.List(volumeEcShardsDeleteRequestShardIDsOff).Len()
}

// ShardIDAt returns the i-th shard_ids element (uint32).
func (t VolumeEcShardsDeleteRequest) ShardIDAt(i int) uint32 {
	return t.o.List(volumeEcShardsDeleteRequestShardIDsOff).Uint32(i)
}

// FullTeardown reads the full_teardown field (proto field 4, bool).
func (t VolumeEcShardsDeleteRequest) FullTeardown() bool {
	return t.o.Bool(volumeEcShardsDeleteRequestFullTeardownOff)
}

// EncodeTsNs reads the encode_ts_ns field (proto field 5, int64).
func (t VolumeEcShardsDeleteRequest) EncodeTsNs() int64 {
	return t.o.Int64(volumeEcShardsDeleteRequestEncodeTsNsOff)
}

// VolumeEcShardsDeleteRequestInput collects the field values for
// NewVolumeEcShardsDeleteRequest.
type VolumeEcShardsDeleteRequestInput struct {
	VolumeID     uint32
	Collection   string
	ShardIDs     []uint32
	FullTeardown bool
	EncodeTsNs   int64
}

// NewVolumeEcShardsDeleteRequest builds a ZAP-encoded VolumeEcShardsDeleteRequest
// message from in and returns the bytes.
func NewVolumeEcShardsDeleteRequest(in VolumeEcShardsDeleteRequestInput) []byte {
	b := zap.NewBuilder(256)
	listOff, listLen := addUint32List(b, in.ShardIDs)
	ob := b.StartObject(volumeEcShardsDeleteRequestSize)
	ob.SetUint32(volumeEcShardsDeleteRequestVolumeIDOff, in.VolumeID)
	ob.SetText(volumeEcShardsDeleteRequestCollectionOff, in.Collection)
	ob.SetList(volumeEcShardsDeleteRequestShardIDsOff, listOff, listLen)
	ob.SetBool(volumeEcShardsDeleteRequestFullTeardownOff, in.FullTeardown)
	ob.SetInt64(volumeEcShardsDeleteRequestEncodeTsNsOff, in.EncodeTsNs)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeEcShardsDeleteResponse ---

const (
	volumeEcShardsDeleteResponseFullTeardownDoneOff = 0
	volumeEcShardsDeleteResponseSize                = 1
)

// VolumeEcShardsDeleteResponse is a zero-copy view into a ZAP-encoded
// VolumeEcShardsDeleteResponse message.
type VolumeEcShardsDeleteResponse struct{ o zap.Object }

// WrapVolumeEcShardsDeleteResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVolumeEcShardsDeleteResponse(b []byte) (VolumeEcShardsDeleteResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeEcShardsDeleteResponse{}, err
	}
	return VolumeEcShardsDeleteResponse{o: m.Root()}, nil
}

// FullTeardownDone reads the full_teardown_done field (proto field 1, bool).
func (t VolumeEcShardsDeleteResponse) FullTeardownDone() bool {
	return t.o.Bool(volumeEcShardsDeleteResponseFullTeardownDoneOff)
}

// VolumeEcShardsDeleteResponseInput collects the field values for
// NewVolumeEcShardsDeleteResponse.
type VolumeEcShardsDeleteResponseInput struct {
	FullTeardownDone bool
}

// NewVolumeEcShardsDeleteResponse builds a ZAP-encoded
// VolumeEcShardsDeleteResponse message from in and returns the bytes.
func NewVolumeEcShardsDeleteResponse(in VolumeEcShardsDeleteResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeEcShardsDeleteResponseSize)
	ob.SetBool(volumeEcShardsDeleteResponseFullTeardownDoneOff, in.FullTeardownDone)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeEcShardsMountRequest ---

const (
	volumeEcShardsMountRequestVolumeIDOff       = 0
	volumeEcShardsMountRequestCollectionOff     = 8
	volumeEcShardsMountRequestShardIDsOff       = 16
	volumeEcShardsMountRequestSourceDiskTypeOff = 24
	volumeEcShardsMountRequestSize              = 32
)

// VolumeEcShardsMountRequest is a zero-copy view into a ZAP-encoded
// VolumeEcShardsMountRequest message.
type VolumeEcShardsMountRequest struct{ o zap.Object }

// WrapVolumeEcShardsMountRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVolumeEcShardsMountRequest(b []byte) (VolumeEcShardsMountRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeEcShardsMountRequest{}, err
	}
	return VolumeEcShardsMountRequest{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t VolumeEcShardsMountRequest) VolumeID() uint32 {
	return t.o.Uint32(volumeEcShardsMountRequestVolumeIDOff)
}

// Collection reads the collection field (proto field 2, string).
func (t VolumeEcShardsMountRequest) Collection() string {
	return t.o.Text(volumeEcShardsMountRequestCollectionOff)
}

// ShardIDsLen reports the number of shard_ids elements (proto field 3, repeated
// uint32).
func (t VolumeEcShardsMountRequest) ShardIDsLen() int {
	return t.o.List(volumeEcShardsMountRequestShardIDsOff).Len()
}

// ShardIDAt returns the i-th shard_ids element (uint32).
func (t VolumeEcShardsMountRequest) ShardIDAt(i int) uint32 {
	return t.o.List(volumeEcShardsMountRequestShardIDsOff).Uint32(i)
}

// SourceDiskType reads the source_disk_type field (proto field 4, string).
func (t VolumeEcShardsMountRequest) SourceDiskType() string {
	return t.o.Text(volumeEcShardsMountRequestSourceDiskTypeOff)
}

// VolumeEcShardsMountRequestInput collects the field values for
// NewVolumeEcShardsMountRequest.
type VolumeEcShardsMountRequestInput struct {
	VolumeID       uint32
	Collection     string
	ShardIDs       []uint32
	SourceDiskType string
}

// NewVolumeEcShardsMountRequest builds a ZAP-encoded VolumeEcShardsMountRequest
// message from in and returns the bytes.
func NewVolumeEcShardsMountRequest(in VolumeEcShardsMountRequestInput) []byte {
	b := zap.NewBuilder(256)
	listOff, listLen := addUint32List(b, in.ShardIDs)
	ob := b.StartObject(volumeEcShardsMountRequestSize)
	ob.SetUint32(volumeEcShardsMountRequestVolumeIDOff, in.VolumeID)
	ob.SetText(volumeEcShardsMountRequestCollectionOff, in.Collection)
	ob.SetList(volumeEcShardsMountRequestShardIDsOff, listOff, listLen)
	ob.SetText(volumeEcShardsMountRequestSourceDiskTypeOff, in.SourceDiskType)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeEcShardsMountResponse ---

const volumeEcShardsMountResponseSize = 0

// VolumeEcShardsMountResponse is a zero-copy view into a ZAP-encoded
// VolumeEcShardsMountResponse message. VolumeEcShardsMountResponse carries no
// fields.
type VolumeEcShardsMountResponse struct{ o zap.Object }

// WrapVolumeEcShardsMountResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVolumeEcShardsMountResponse(b []byte) (VolumeEcShardsMountResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeEcShardsMountResponse{}, err
	}
	return VolumeEcShardsMountResponse{o: m.Root()}, nil
}

// VolumeEcShardsMountResponseInput collects the field values for
// NewVolumeEcShardsMountResponse. VolumeEcShardsMountResponse has no fields.
type VolumeEcShardsMountResponseInput struct{}

// NewVolumeEcShardsMountResponse builds a ZAP-encoded VolumeEcShardsMountResponse
// message and returns the bytes.
func NewVolumeEcShardsMountResponse(in VolumeEcShardsMountResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeEcShardsMountResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeEcShardsUnmountRequest ---
//
// Proto field 2 is absent (skipped in the .proto); present fields take slots in
// declaration order.

const (
	volumeEcShardsUnmountRequestVolumeIDOff   = 0
	volumeEcShardsUnmountRequestShardIDsOff   = 8
	volumeEcShardsUnmountRequestEncodeTsNsOff = 16
	volumeEcShardsUnmountRequestSize          = 24
)

// VolumeEcShardsUnmountRequest is a zero-copy view into a ZAP-encoded
// VolumeEcShardsUnmountRequest message.
type VolumeEcShardsUnmountRequest struct{ o zap.Object }

// WrapVolumeEcShardsUnmountRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVolumeEcShardsUnmountRequest(b []byte) (VolumeEcShardsUnmountRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeEcShardsUnmountRequest{}, err
	}
	return VolumeEcShardsUnmountRequest{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t VolumeEcShardsUnmountRequest) VolumeID() uint32 {
	return t.o.Uint32(volumeEcShardsUnmountRequestVolumeIDOff)
}

// ShardIDsLen reports the number of shard_ids elements (proto field 3, repeated
// uint32).
func (t VolumeEcShardsUnmountRequest) ShardIDsLen() int {
	return t.o.List(volumeEcShardsUnmountRequestShardIDsOff).Len()
}

// ShardIDAt returns the i-th shard_ids element (uint32).
func (t VolumeEcShardsUnmountRequest) ShardIDAt(i int) uint32 {
	return t.o.List(volumeEcShardsUnmountRequestShardIDsOff).Uint32(i)
}

// EncodeTsNs reads the encode_ts_ns field (proto field 4, int64).
func (t VolumeEcShardsUnmountRequest) EncodeTsNs() int64 {
	return t.o.Int64(volumeEcShardsUnmountRequestEncodeTsNsOff)
}

// VolumeEcShardsUnmountRequestInput collects the field values for
// NewVolumeEcShardsUnmountRequest.
type VolumeEcShardsUnmountRequestInput struct {
	VolumeID   uint32
	ShardIDs   []uint32
	EncodeTsNs int64
}

// NewVolumeEcShardsUnmountRequest builds a ZAP-encoded
// VolumeEcShardsUnmountRequest message from in and returns the bytes.
func NewVolumeEcShardsUnmountRequest(in VolumeEcShardsUnmountRequestInput) []byte {
	b := zap.NewBuilder(256)
	listOff, listLen := addUint32List(b, in.ShardIDs)
	ob := b.StartObject(volumeEcShardsUnmountRequestSize)
	ob.SetUint32(volumeEcShardsUnmountRequestVolumeIDOff, in.VolumeID)
	ob.SetList(volumeEcShardsUnmountRequestShardIDsOff, listOff, listLen)
	ob.SetInt64(volumeEcShardsUnmountRequestEncodeTsNsOff, in.EncodeTsNs)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeEcShardsUnmountResponse ---

const volumeEcShardsUnmountResponseSize = 0

// VolumeEcShardsUnmountResponse is a zero-copy view into a ZAP-encoded
// VolumeEcShardsUnmountResponse message. VolumeEcShardsUnmountResponse carries
// no fields.
type VolumeEcShardsUnmountResponse struct{ o zap.Object }

// WrapVolumeEcShardsUnmountResponse parses b and returns a typed view. Returns
// an error if the wire-level checks (magic, version, size) fail.
func WrapVolumeEcShardsUnmountResponse(b []byte) (VolumeEcShardsUnmountResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeEcShardsUnmountResponse{}, err
	}
	return VolumeEcShardsUnmountResponse{o: m.Root()}, nil
}

// VolumeEcShardsUnmountResponseInput collects the field values for
// NewVolumeEcShardsUnmountResponse. VolumeEcShardsUnmountResponse has no fields.
type VolumeEcShardsUnmountResponseInput struct{}

// NewVolumeEcShardsUnmountResponse builds a ZAP-encoded
// VolumeEcShardsUnmountResponse message and returns the bytes.
func NewVolumeEcShardsUnmountResponse(in VolumeEcShardsUnmountResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeEcShardsUnmountResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeEcShardReadRequest ---
//
// Proto field 6 is reserved; present fields take slots in declaration order.

const (
	volumeEcShardReadRequestVolumeIDOff   = 0
	volumeEcShardReadRequestShardIDOff    = 4
	volumeEcShardReadRequestOffsetOff     = 8
	volumeEcShardReadRequestSizeFOff      = 16
	volumeEcShardReadRequestFileKeyOff    = 24
	volumeEcShardReadRequestEncodeTsNsOff = 32
	volumeEcShardReadRequestSize          = 40
)

// VolumeEcShardReadRequest is a zero-copy view into a ZAP-encoded
// VolumeEcShardReadRequest message.
type VolumeEcShardReadRequest struct{ o zap.Object }

// WrapVolumeEcShardReadRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVolumeEcShardReadRequest(b []byte) (VolumeEcShardReadRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeEcShardReadRequest{}, err
	}
	return VolumeEcShardReadRequest{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t VolumeEcShardReadRequest) VolumeID() uint32 {
	return t.o.Uint32(volumeEcShardReadRequestVolumeIDOff)
}

// ShardID reads the shard_id field (proto field 2, uint32).
func (t VolumeEcShardReadRequest) ShardID() uint32 {
	return t.o.Uint32(volumeEcShardReadRequestShardIDOff)
}

// Offset reads the offset field (proto field 3, int64).
func (t VolumeEcShardReadRequest) Offset() int64 {
	return t.o.Int64(volumeEcShardReadRequestOffsetOff)
}

// Size reads the size field (proto field 4, int64).
func (t VolumeEcShardReadRequest) Size() int64 { return t.o.Int64(volumeEcShardReadRequestSizeFOff) }

// FileKey reads the file_key field (proto field 5, uint64).
func (t VolumeEcShardReadRequest) FileKey() uint64 {
	return t.o.Uint64(volumeEcShardReadRequestFileKeyOff)
}

// EncodeTsNs reads the encode_ts_ns field (proto field 7, int64).
func (t VolumeEcShardReadRequest) EncodeTsNs() int64 {
	return t.o.Int64(volumeEcShardReadRequestEncodeTsNsOff)
}

// VolumeEcShardReadRequestInput collects the field values for
// NewVolumeEcShardReadRequest.
type VolumeEcShardReadRequestInput struct {
	VolumeID   uint32
	ShardID    uint32
	Offset     int64
	Size       int64
	FileKey    uint64
	EncodeTsNs int64
}

// NewVolumeEcShardReadRequest builds a ZAP-encoded VolumeEcShardReadRequest
// message from in and returns the bytes.
func NewVolumeEcShardReadRequest(in VolumeEcShardReadRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeEcShardReadRequestSize)
	ob.SetUint32(volumeEcShardReadRequestVolumeIDOff, in.VolumeID)
	ob.SetUint32(volumeEcShardReadRequestShardIDOff, in.ShardID)
	ob.SetInt64(volumeEcShardReadRequestOffsetOff, in.Offset)
	ob.SetInt64(volumeEcShardReadRequestSizeFOff, in.Size)
	ob.SetUint64(volumeEcShardReadRequestFileKeyOff, in.FileKey)
	ob.SetInt64(volumeEcShardReadRequestEncodeTsNsOff, in.EncodeTsNs)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeEcShardReadResponse ---

const (
	volumeEcShardReadResponseDataOff       = 0
	volumeEcShardReadResponseIsDeletedOff  = 8
	volumeEcShardReadResponseEncodeTsNsOff = 16
	volumeEcShardReadResponseSize          = 24
)

// VolumeEcShardReadResponse is a zero-copy view into a ZAP-encoded
// VolumeEcShardReadResponse message.
type VolumeEcShardReadResponse struct{ o zap.Object }

// WrapVolumeEcShardReadResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVolumeEcShardReadResponse(b []byte) (VolumeEcShardReadResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeEcShardReadResponse{}, err
	}
	return VolumeEcShardReadResponse{o: m.Root()}, nil
}

// Data reads the data field (proto field 1, bytes).
func (t VolumeEcShardReadResponse) Data() []byte { return t.o.Bytes(volumeEcShardReadResponseDataOff) }

// IsDeleted reads the is_deleted field (proto field 2, bool).
func (t VolumeEcShardReadResponse) IsDeleted() bool {
	return t.o.Bool(volumeEcShardReadResponseIsDeletedOff)
}

// EncodeTsNs reads the encode_ts_ns field (proto field 3, int64).
func (t VolumeEcShardReadResponse) EncodeTsNs() int64 {
	return t.o.Int64(volumeEcShardReadResponseEncodeTsNsOff)
}

// VolumeEcShardReadResponseInput collects the field values for
// NewVolumeEcShardReadResponse.
type VolumeEcShardReadResponseInput struct {
	Data       []byte
	IsDeleted  bool
	EncodeTsNs int64
}

// NewVolumeEcShardReadResponse builds a ZAP-encoded VolumeEcShardReadResponse
// message from in and returns the bytes.
func NewVolumeEcShardReadResponse(in VolumeEcShardReadResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeEcShardReadResponseSize)
	ob.SetBytes(volumeEcShardReadResponseDataOff, in.Data)
	ob.SetBool(volumeEcShardReadResponseIsDeletedOff, in.IsDeleted)
	ob.SetInt64(volumeEcShardReadResponseEncodeTsNsOff, in.EncodeTsNs)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeEcBlobDeleteRequest ---

const (
	volumeEcBlobDeleteRequestVolumeIDOff   = 0
	volumeEcBlobDeleteRequestCollectionOff = 8
	volumeEcBlobDeleteRequestFileKeyOff    = 16
	volumeEcBlobDeleteRequestVersionOff    = 24
	volumeEcBlobDeleteRequestSize          = 28
)

// VolumeEcBlobDeleteRequest is a zero-copy view into a ZAP-encoded
// VolumeEcBlobDeleteRequest message.
type VolumeEcBlobDeleteRequest struct{ o zap.Object }

// WrapVolumeEcBlobDeleteRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVolumeEcBlobDeleteRequest(b []byte) (VolumeEcBlobDeleteRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeEcBlobDeleteRequest{}, err
	}
	return VolumeEcBlobDeleteRequest{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t VolumeEcBlobDeleteRequest) VolumeID() uint32 {
	return t.o.Uint32(volumeEcBlobDeleteRequestVolumeIDOff)
}

// Collection reads the collection field (proto field 2, string).
func (t VolumeEcBlobDeleteRequest) Collection() string {
	return t.o.Text(volumeEcBlobDeleteRequestCollectionOff)
}

// FileKey reads the file_key field (proto field 3, uint64).
func (t VolumeEcBlobDeleteRequest) FileKey() uint64 {
	return t.o.Uint64(volumeEcBlobDeleteRequestFileKeyOff)
}

// Version reads the version field (proto field 4, uint32).
func (t VolumeEcBlobDeleteRequest) Version() uint32 {
	return t.o.Uint32(volumeEcBlobDeleteRequestVersionOff)
}

// VolumeEcBlobDeleteRequestInput collects the field values for
// NewVolumeEcBlobDeleteRequest.
type VolumeEcBlobDeleteRequestInput struct {
	VolumeID   uint32
	Collection string
	FileKey    uint64
	Version    uint32
}

// NewVolumeEcBlobDeleteRequest builds a ZAP-encoded VolumeEcBlobDeleteRequest
// message from in and returns the bytes.
func NewVolumeEcBlobDeleteRequest(in VolumeEcBlobDeleteRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeEcBlobDeleteRequestSize)
	ob.SetUint32(volumeEcBlobDeleteRequestVolumeIDOff, in.VolumeID)
	ob.SetText(volumeEcBlobDeleteRequestCollectionOff, in.Collection)
	ob.SetUint64(volumeEcBlobDeleteRequestFileKeyOff, in.FileKey)
	ob.SetUint32(volumeEcBlobDeleteRequestVersionOff, in.Version)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeEcBlobDeleteResponse ---

const volumeEcBlobDeleteResponseSize = 0

// VolumeEcBlobDeleteResponse is a zero-copy view into a ZAP-encoded
// VolumeEcBlobDeleteResponse message. VolumeEcBlobDeleteResponse carries no
// fields.
type VolumeEcBlobDeleteResponse struct{ o zap.Object }

// WrapVolumeEcBlobDeleteResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVolumeEcBlobDeleteResponse(b []byte) (VolumeEcBlobDeleteResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeEcBlobDeleteResponse{}, err
	}
	return VolumeEcBlobDeleteResponse{o: m.Root()}, nil
}

// VolumeEcBlobDeleteResponseInput collects the field values for
// NewVolumeEcBlobDeleteResponse. VolumeEcBlobDeleteResponse has no fields.
type VolumeEcBlobDeleteResponseInput struct{}

// NewVolumeEcBlobDeleteResponse builds a ZAP-encoded VolumeEcBlobDeleteResponse
// message and returns the bytes.
func NewVolumeEcBlobDeleteResponse(in VolumeEcBlobDeleteResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeEcBlobDeleteResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeEcShardsToVolumeRequest ---

const (
	volumeEcShardsToVolumeRequestVolumeIDOff   = 0
	volumeEcShardsToVolumeRequestCollectionOff = 8
	volumeEcShardsToVolumeRequestSize          = 16
)

// VolumeEcShardsToVolumeRequest is a zero-copy view into a ZAP-encoded
// VolumeEcShardsToVolumeRequest message.
type VolumeEcShardsToVolumeRequest struct{ o zap.Object }

// WrapVolumeEcShardsToVolumeRequest parses b and returns a typed view. Returns
// an error if the wire-level checks (magic, version, size) fail.
func WrapVolumeEcShardsToVolumeRequest(b []byte) (VolumeEcShardsToVolumeRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeEcShardsToVolumeRequest{}, err
	}
	return VolumeEcShardsToVolumeRequest{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t VolumeEcShardsToVolumeRequest) VolumeID() uint32 {
	return t.o.Uint32(volumeEcShardsToVolumeRequestVolumeIDOff)
}

// Collection reads the collection field (proto field 2, string).
func (t VolumeEcShardsToVolumeRequest) Collection() string {
	return t.o.Text(volumeEcShardsToVolumeRequestCollectionOff)
}

// VolumeEcShardsToVolumeRequestInput collects the field values for
// NewVolumeEcShardsToVolumeRequest.
type VolumeEcShardsToVolumeRequestInput struct {
	VolumeID   uint32
	Collection string
}

// NewVolumeEcShardsToVolumeRequest builds a ZAP-encoded
// VolumeEcShardsToVolumeRequest message from in and returns the bytes.
func NewVolumeEcShardsToVolumeRequest(in VolumeEcShardsToVolumeRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeEcShardsToVolumeRequestSize)
	ob.SetUint32(volumeEcShardsToVolumeRequestVolumeIDOff, in.VolumeID)
	ob.SetText(volumeEcShardsToVolumeRequestCollectionOff, in.Collection)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeEcShardsToVolumeResponse ---

const volumeEcShardsToVolumeResponseSize = 0

// VolumeEcShardsToVolumeResponse is a zero-copy view into a ZAP-encoded
// VolumeEcShardsToVolumeResponse message. VolumeEcShardsToVolumeResponse carries
// no fields.
type VolumeEcShardsToVolumeResponse struct{ o zap.Object }

// WrapVolumeEcShardsToVolumeResponse parses b and returns a typed view. Returns
// an error if the wire-level checks (magic, version, size) fail.
func WrapVolumeEcShardsToVolumeResponse(b []byte) (VolumeEcShardsToVolumeResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeEcShardsToVolumeResponse{}, err
	}
	return VolumeEcShardsToVolumeResponse{o: m.Root()}, nil
}

// VolumeEcShardsToVolumeResponseInput collects the field values for
// NewVolumeEcShardsToVolumeResponse. VolumeEcShardsToVolumeResponse has no
// fields.
type VolumeEcShardsToVolumeResponseInput struct{}

// NewVolumeEcShardsToVolumeResponse builds a ZAP-encoded
// VolumeEcShardsToVolumeResponse message and returns the bytes.
func NewVolumeEcShardsToVolumeResponse(in VolumeEcShardsToVolumeResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeEcShardsToVolumeResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeEcShardsInfoRequest ---

const (
	volumeEcShardsInfoRequestVolumeIDOff = 0
	volumeEcShardsInfoRequestSize        = 4
)

// VolumeEcShardsInfoRequest is a zero-copy view into a ZAP-encoded
// VolumeEcShardsInfoRequest message.
type VolumeEcShardsInfoRequest struct{ o zap.Object }

// WrapVolumeEcShardsInfoRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVolumeEcShardsInfoRequest(b []byte) (VolumeEcShardsInfoRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeEcShardsInfoRequest{}, err
	}
	return VolumeEcShardsInfoRequest{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t VolumeEcShardsInfoRequest) VolumeID() uint32 {
	return t.o.Uint32(volumeEcShardsInfoRequestVolumeIDOff)
}

// VolumeEcShardsInfoRequestInput collects the field values for
// NewVolumeEcShardsInfoRequest.
type VolumeEcShardsInfoRequestInput struct {
	VolumeID uint32
}

// NewVolumeEcShardsInfoRequest builds a ZAP-encoded VolumeEcShardsInfoRequest
// message from in and returns the bytes.
func NewVolumeEcShardsInfoRequest(in VolumeEcShardsInfoRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeEcShardsInfoRequestSize)
	ob.SetUint32(volumeEcShardsInfoRequestVolumeIDOff, in.VolumeID)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeEcShardsInfoResponse ---

const (
	volumeEcShardsInfoResponseEcShardInfosOff     = 0
	volumeEcShardsInfoResponseVolumeSizeOff       = 8
	volumeEcShardsInfoResponseFileCountOff        = 16
	volumeEcShardsInfoResponseFileDeletedCountOff = 24
	volumeEcShardsInfoResponseSize                = 32
)

// VolumeEcShardsInfoResponse is a zero-copy view into a ZAP-encoded
// VolumeEcShardsInfoResponse message.
type VolumeEcShardsInfoResponse struct{ o zap.Object }

// WrapVolumeEcShardsInfoResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVolumeEcShardsInfoResponse(b []byte) (VolumeEcShardsInfoResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeEcShardsInfoResponse{}, err
	}
	return VolumeEcShardsInfoResponse{o: m.Root()}, nil
}

// EcShardInfosLen reports the number of ec_shard_infos elements (proto field 1,
// repeated EcShardInfo).
func (t VolumeEcShardsInfoResponse) EcShardInfosLen() int {
	return t.o.List(volumeEcShardsInfoResponseEcShardInfosOff).Len()
}

// EcShardInfoAt returns the i-th ec_shard_infos element. The bool is false when
// i is out of range or the element fails to parse.
func (t VolumeEcShardsInfoResponse) EcShardInfoAt(i int) (EcShardInfo, bool) {
	o := t.o.List(volumeEcShardsInfoResponseEcShardInfosOff).ObjectAt(i)
	if o.IsNull() {
		return EcShardInfo{}, false
	}
	return EcShardInfo{o: o}, true
}

// VolumeSize reads the volume_size field (proto field 2, uint64).
func (t VolumeEcShardsInfoResponse) VolumeSize() uint64 {
	return t.o.Uint64(volumeEcShardsInfoResponseVolumeSizeOff)
}

// FileCount reads the file_count field (proto field 3, uint64).
func (t VolumeEcShardsInfoResponse) FileCount() uint64 {
	return t.o.Uint64(volumeEcShardsInfoResponseFileCountOff)
}

// FileDeletedCount reads the file_deleted_count field (proto field 4, uint64).
func (t VolumeEcShardsInfoResponse) FileDeletedCount() uint64 {
	return t.o.Uint64(volumeEcShardsInfoResponseFileDeletedCountOff)
}

// VolumeEcShardsInfoResponseInput collects the field values for
// NewVolumeEcShardsInfoResponse. Each EcShardInfos entry is an EcShardInfo
// message's own ZAP buffer (from NewEcShardInfo).
type VolumeEcShardsInfoResponseInput struct {
	EcShardInfos     [][]byte
	VolumeSize       uint64
	FileCount        uint64
	FileDeletedCount uint64
}

// NewVolumeEcShardsInfoResponse builds a ZAP-encoded VolumeEcShardsInfoResponse
// message from in and returns the bytes.
func NewVolumeEcShardsInfoResponse(in VolumeEcShardsInfoResponseInput) []byte {
	b := zap.NewBuilder(256)
	listOff, listLen := addObjectList(b, in.EcShardInfos)
	ob := b.StartObject(volumeEcShardsInfoResponseSize)
	ob.SetList(volumeEcShardsInfoResponseEcShardInfosOff, listOff, listLen)
	ob.SetUint64(volumeEcShardsInfoResponseVolumeSizeOff, in.VolumeSize)
	ob.SetUint64(volumeEcShardsInfoResponseFileCountOff, in.FileCount)
	ob.SetUint64(volumeEcShardsInfoResponseFileDeletedCountOff, in.FileDeletedCount)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- EcShardInfo ---

const (
	ecShardInfoShardIDOff    = 0
	ecShardInfoSizeFOff      = 8
	ecShardInfoCollectionOff = 16
	ecShardInfoVolumeIDOff   = 24
	ecShardInfoSize          = 28
)

// EcShardInfo is a zero-copy view into a ZAP-encoded EcShardInfo message.
type EcShardInfo struct{ o zap.Object }

// WrapEcShardInfo parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapEcShardInfo(b []byte) (EcShardInfo, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return EcShardInfo{}, err
	}
	return EcShardInfo{o: m.Root()}, nil
}

// ShardID reads the shard_id field (proto field 1, uint32).
func (t EcShardInfo) ShardID() uint32 { return t.o.Uint32(ecShardInfoShardIDOff) }

// Size reads the size field (proto field 2, int64).
func (t EcShardInfo) Size() int64 { return t.o.Int64(ecShardInfoSizeFOff) }

// Collection reads the collection field (proto field 3, string).
func (t EcShardInfo) Collection() string { return t.o.Text(ecShardInfoCollectionOff) }

// VolumeID reads the volume_id field (proto field 4, uint32).
func (t EcShardInfo) VolumeID() uint32 { return t.o.Uint32(ecShardInfoVolumeIDOff) }

// EcShardInfoInput collects the field values for NewEcShardInfo.
type EcShardInfoInput struct {
	ShardID    uint32
	Size       int64
	Collection string
	VolumeID   uint32
}

// NewEcShardInfo builds a ZAP-encoded EcShardInfo message from in and returns
// the bytes.
func NewEcShardInfo(in EcShardInfoInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(ecShardInfoSize)
	ob.SetUint32(ecShardInfoShardIDOff, in.ShardID)
	ob.SetInt64(ecShardInfoSizeFOff, in.Size)
	ob.SetText(ecShardInfoCollectionOff, in.Collection)
	ob.SetUint32(ecShardInfoVolumeIDOff, in.VolumeID)
	ob.FinishAsRoot()
	return b.Finish()
}
