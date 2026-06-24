// Code generated from volume_server.proto; DO NOT EDIT.

package volume_serverwire

import (
	zap "github.com/zap-proto/go"
)

// --- ReadVolumeFileStatusRequest ---

const (
	readVolumeFileStatusRequestVolumeIDOff = 0
	readVolumeFileStatusRequestSize        = 4
)

// ReadVolumeFileStatusRequest is a zero-copy view into a ZAP-encoded
// ReadVolumeFileStatusRequest message.
type ReadVolumeFileStatusRequest struct{ o zap.Object }

// WrapReadVolumeFileStatusRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapReadVolumeFileStatusRequest(b []byte) (ReadVolumeFileStatusRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ReadVolumeFileStatusRequest{}, err
	}
	return ReadVolumeFileStatusRequest{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t ReadVolumeFileStatusRequest) VolumeID() uint32 {
	return t.o.Uint32(readVolumeFileStatusRequestVolumeIDOff)
}

// ReadVolumeFileStatusRequestInput collects the field values for
// NewReadVolumeFileStatusRequest.
type ReadVolumeFileStatusRequestInput struct {
	VolumeID uint32
}

// NewReadVolumeFileStatusRequest builds a ZAP-encoded ReadVolumeFileStatusRequest
// message from in and returns the bytes.
func NewReadVolumeFileStatusRequest(in ReadVolumeFileStatusRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(readVolumeFileStatusRequestSize)
	ob.SetUint32(readVolumeFileStatusRequestVolumeIDOff, in.VolumeID)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ReadVolumeFileStatusResponse ---

const (
	readVolumeFileStatusResponseVolumeIDOff                = 0
	readVolumeFileStatusResponseIdxFileTimestampSecondsOff = 8
	readVolumeFileStatusResponseIdxFileSizeOff             = 16
	readVolumeFileStatusResponseDatFileTimestampSecondsOff = 24
	readVolumeFileStatusResponseDatFileSizeOff             = 32
	readVolumeFileStatusResponseFileCountOff               = 40
	readVolumeFileStatusResponseCompactionRevisionOff      = 48
	readVolumeFileStatusResponseCollectionOff              = 52
	readVolumeFileStatusResponseDiskTypeOff                = 60
	readVolumeFileStatusResponseVolumeInfoOff              = 68
	readVolumeFileStatusResponseVersionOff                 = 72
	readVolumeFileStatusResponseSize                       = 76
)

// ReadVolumeFileStatusResponse is a zero-copy view into a ZAP-encoded
// ReadVolumeFileStatusResponse message.
type ReadVolumeFileStatusResponse struct{ o zap.Object }

// WrapReadVolumeFileStatusResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapReadVolumeFileStatusResponse(b []byte) (ReadVolumeFileStatusResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ReadVolumeFileStatusResponse{}, err
	}
	return ReadVolumeFileStatusResponse{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t ReadVolumeFileStatusResponse) VolumeID() uint32 {
	return t.o.Uint32(readVolumeFileStatusResponseVolumeIDOff)
}

// IdxFileTimestampSeconds reads the idx_file_timestamp_seconds field (proto
// field 2, uint64).
func (t ReadVolumeFileStatusResponse) IdxFileTimestampSeconds() uint64 {
	return t.o.Uint64(readVolumeFileStatusResponseIdxFileTimestampSecondsOff)
}

// IdxFileSize reads the idx_file_size field (proto field 3, uint64).
func (t ReadVolumeFileStatusResponse) IdxFileSize() uint64 {
	return t.o.Uint64(readVolumeFileStatusResponseIdxFileSizeOff)
}

// DatFileTimestampSeconds reads the dat_file_timestamp_seconds field (proto
// field 4, uint64).
func (t ReadVolumeFileStatusResponse) DatFileTimestampSeconds() uint64 {
	return t.o.Uint64(readVolumeFileStatusResponseDatFileTimestampSecondsOff)
}

// DatFileSize reads the dat_file_size field (proto field 5, uint64).
func (t ReadVolumeFileStatusResponse) DatFileSize() uint64 {
	return t.o.Uint64(readVolumeFileStatusResponseDatFileSizeOff)
}

// FileCount reads the file_count field (proto field 6, uint64).
func (t ReadVolumeFileStatusResponse) FileCount() uint64 {
	return t.o.Uint64(readVolumeFileStatusResponseFileCountOff)
}

// CompactionRevision reads the compaction_revision field (proto field 7,
// uint32).
func (t ReadVolumeFileStatusResponse) CompactionRevision() uint32 {
	return t.o.Uint32(readVolumeFileStatusResponseCompactionRevisionOff)
}

// Collection reads the collection field (proto field 8, string).
func (t ReadVolumeFileStatusResponse) Collection() string {
	return t.o.Text(readVolumeFileStatusResponseCollectionOff)
}

// DiskType reads the disk_type field (proto field 9, string).
func (t ReadVolumeFileStatusResponse) DiskType() string {
	return t.o.Text(readVolumeFileStatusResponseDiskTypeOff)
}

// VolumeInfo reads the volume_info field (proto field 10, message VolumeInfo).
// The bool is false when the field is absent.
func (t ReadVolumeFileStatusResponse) VolumeInfo() (VolumeInfo, bool) {
	o := t.o.Object(readVolumeFileStatusResponseVolumeInfoOff)
	if o.IsNull() {
		return VolumeInfo{}, false
	}
	return VolumeInfo{o: o}, true
}

// Version reads the version field (proto field 11, uint32).
func (t ReadVolumeFileStatusResponse) Version() uint32 {
	return t.o.Uint32(readVolumeFileStatusResponseVersionOff)
}

// ReadVolumeFileStatusResponseInput collects the field values for
// NewReadVolumeFileStatusResponse. VolumeInfo is the VolumeInfo message's own
// ZAP buffer (from NewVolumeInfo).
type ReadVolumeFileStatusResponseInput struct {
	VolumeID                uint32
	IdxFileTimestampSeconds uint64
	IdxFileSize             uint64
	DatFileTimestampSeconds uint64
	DatFileSize             uint64
	FileCount               uint64
	CompactionRevision      uint32
	Collection              string
	DiskType                string
	VolumeInfo              []byte
	Version                 uint32
}

// NewReadVolumeFileStatusResponse builds a ZAP-encoded
// ReadVolumeFileStatusResponse message from in and returns the bytes.
func NewReadVolumeFileStatusResponse(in ReadVolumeFileStatusResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(readVolumeFileStatusResponseSize)
	ob.SetUint32(readVolumeFileStatusResponseVolumeIDOff, in.VolumeID)
	ob.SetUint64(readVolumeFileStatusResponseIdxFileTimestampSecondsOff, in.IdxFileTimestampSeconds)
	ob.SetUint64(readVolumeFileStatusResponseIdxFileSizeOff, in.IdxFileSize)
	ob.SetUint64(readVolumeFileStatusResponseDatFileTimestampSecondsOff, in.DatFileTimestampSeconds)
	ob.SetUint64(readVolumeFileStatusResponseDatFileSizeOff, in.DatFileSize)
	ob.SetUint64(readVolumeFileStatusResponseFileCountOff, in.FileCount)
	ob.SetUint32(readVolumeFileStatusResponseCompactionRevisionOff, in.CompactionRevision)
	ob.SetText(readVolumeFileStatusResponseCollectionOff, in.Collection)
	ob.SetText(readVolumeFileStatusResponseDiskTypeOff, in.DiskType)
	ob.SetObject(readVolumeFileStatusResponseVolumeInfoOff, embedMessage(b, in.VolumeInfo))
	ob.SetUint32(readVolumeFileStatusResponseVersionOff, in.Version)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- DiskStatus ---

const (
	diskStatusDirOff         = 0
	diskStatusAllOff         = 8
	diskStatusUsedOff        = 16
	diskStatusFreeOff        = 24
	diskStatusPercentFreeOff = 32
	diskStatusPercentUsedOff = 36
	diskStatusDiskTypeOff    = 40
	diskStatusErrorOff       = 48
	diskStatusSize           = 56
)

// DiskStatus is a zero-copy view into a ZAP-encoded DiskStatus message.
type DiskStatus struct{ o zap.Object }

// WrapDiskStatus parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapDiskStatus(b []byte) (DiskStatus, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return DiskStatus{}, err
	}
	return DiskStatus{o: m.Root()}, nil
}

// Dir reads the dir field (proto field 1, string).
func (t DiskStatus) Dir() string { return t.o.Text(diskStatusDirOff) }

// All reads the all field (proto field 2, uint64).
func (t DiskStatus) All() uint64 { return t.o.Uint64(diskStatusAllOff) }

// Used reads the used field (proto field 3, uint64).
func (t DiskStatus) Used() uint64 { return t.o.Uint64(diskStatusUsedOff) }

// Free reads the free field (proto field 4, uint64).
func (t DiskStatus) Free() uint64 { return t.o.Uint64(diskStatusFreeOff) }

// PercentFree reads the percent_free field (proto field 5, float).
func (t DiskStatus) PercentFree() float32 { return t.o.Float32(diskStatusPercentFreeOff) }

// PercentUsed reads the percent_used field (proto field 6, float).
func (t DiskStatus) PercentUsed() float32 { return t.o.Float32(diskStatusPercentUsedOff) }

// DiskType reads the disk_type field (proto field 7, string).
func (t DiskStatus) DiskType() string { return t.o.Text(diskStatusDiskTypeOff) }

// Error reads the error field (proto field 8, string).
func (t DiskStatus) Error() string { return t.o.Text(diskStatusErrorOff) }

// DiskStatusInput collects the field values for NewDiskStatus.
type DiskStatusInput struct {
	Dir         string
	All         uint64
	Used        uint64
	Free        uint64
	PercentFree float32
	PercentUsed float32
	DiskType    string
	Error       string
}

// NewDiskStatus builds a ZAP-encoded DiskStatus message from in and returns the
// bytes.
func NewDiskStatus(in DiskStatusInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(diskStatusSize)
	ob.SetText(diskStatusDirOff, in.Dir)
	ob.SetUint64(diskStatusAllOff, in.All)
	ob.SetUint64(diskStatusUsedOff, in.Used)
	ob.SetUint64(diskStatusFreeOff, in.Free)
	ob.SetFloat32(diskStatusPercentFreeOff, in.PercentFree)
	ob.SetFloat32(diskStatusPercentUsedOff, in.PercentUsed)
	ob.SetText(diskStatusDiskTypeOff, in.DiskType)
	ob.SetText(diskStatusErrorOff, in.Error)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- MemStatus ---

const (
	memStatusGoroutinesOff = 0
	memStatusAllOff        = 8
	memStatusUsedOff       = 16
	memStatusFreeOff       = 24
	memStatusSelfOff       = 32
	memStatusHeapOff       = 40
	memStatusStackOff      = 48
	memStatusSize          = 56
)

// MemStatus is a zero-copy view into a ZAP-encoded MemStatus message.
type MemStatus struct{ o zap.Object }

// WrapMemStatus parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapMemStatus(b []byte) (MemStatus, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return MemStatus{}, err
	}
	return MemStatus{o: m.Root()}, nil
}

// Goroutines reads the goroutines field (proto field 1, int32).
func (t MemStatus) Goroutines() int32 { return t.o.Int32(memStatusGoroutinesOff) }

// All reads the all field (proto field 2, uint64).
func (t MemStatus) All() uint64 { return t.o.Uint64(memStatusAllOff) }

// Used reads the used field (proto field 3, uint64).
func (t MemStatus) Used() uint64 { return t.o.Uint64(memStatusUsedOff) }

// Free reads the free field (proto field 4, uint64).
func (t MemStatus) Free() uint64 { return t.o.Uint64(memStatusFreeOff) }

// Self reads the self field (proto field 5, uint64).
func (t MemStatus) Self() uint64 { return t.o.Uint64(memStatusSelfOff) }

// Heap reads the heap field (proto field 6, uint64).
func (t MemStatus) Heap() uint64 { return t.o.Uint64(memStatusHeapOff) }

// Stack reads the stack field (proto field 7, uint64).
func (t MemStatus) Stack() uint64 { return t.o.Uint64(memStatusStackOff) }

// MemStatusInput collects the field values for NewMemStatus.
type MemStatusInput struct {
	Goroutines int32
	All        uint64
	Used       uint64
	Free       uint64
	Self       uint64
	Heap       uint64
	Stack      uint64
}

// NewMemStatus builds a ZAP-encoded MemStatus message from in and returns the
// bytes.
func NewMemStatus(in MemStatusInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(memStatusSize)
	ob.SetInt32(memStatusGoroutinesOff, in.Goroutines)
	ob.SetUint64(memStatusAllOff, in.All)
	ob.SetUint64(memStatusUsedOff, in.Used)
	ob.SetUint64(memStatusFreeOff, in.Free)
	ob.SetUint64(memStatusSelfOff, in.Self)
	ob.SetUint64(memStatusHeapOff, in.Heap)
	ob.SetUint64(memStatusStackOff, in.Stack)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- RemoteFile ---

const (
	remoteFileBackendTypeOff  = 0
	remoteFileBackendIDOff    = 8
	remoteFileKeyOff          = 16
	remoteFileOffsetOff       = 24
	remoteFileFileSizeOff     = 32
	remoteFileModifiedTimeOff = 40
	remoteFileExtensionOff    = 48
	remoteFileSize            = 56
)

// RemoteFile is a zero-copy view into a ZAP-encoded RemoteFile message.
type RemoteFile struct{ o zap.Object }

// WrapRemoteFile parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapRemoteFile(b []byte) (RemoteFile, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return RemoteFile{}, err
	}
	return RemoteFile{o: m.Root()}, nil
}

// BackendType reads the backend_type field (proto field 1, string).
func (t RemoteFile) BackendType() string { return t.o.Text(remoteFileBackendTypeOff) }

// BackendID reads the backend_id field (proto field 2, string).
func (t RemoteFile) BackendID() string { return t.o.Text(remoteFileBackendIDOff) }

// Key reads the key field (proto field 3, string).
func (t RemoteFile) Key() string { return t.o.Text(remoteFileKeyOff) }

// Offset reads the offset field (proto field 4, uint64).
func (t RemoteFile) Offset() uint64 { return t.o.Uint64(remoteFileOffsetOff) }

// FileSize reads the file_size field (proto field 5, uint64).
func (t RemoteFile) FileSize() uint64 { return t.o.Uint64(remoteFileFileSizeOff) }

// ModifiedTime reads the modified_time field (proto field 6, uint64).
func (t RemoteFile) ModifiedTime() uint64 { return t.o.Uint64(remoteFileModifiedTimeOff) }

// Extension reads the extension field (proto field 7, string).
func (t RemoteFile) Extension() string { return t.o.Text(remoteFileExtensionOff) }

// RemoteFileInput collects the field values for NewRemoteFile.
type RemoteFileInput struct {
	BackendType  string
	BackendID    string
	Key          string
	Offset       uint64
	FileSize     uint64
	ModifiedTime uint64
	Extension    string
}

// NewRemoteFile builds a ZAP-encoded RemoteFile message from in and returns the
// bytes.
func NewRemoteFile(in RemoteFileInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(remoteFileSize)
	ob.SetText(remoteFileBackendTypeOff, in.BackendType)
	ob.SetText(remoteFileBackendIDOff, in.BackendID)
	ob.SetText(remoteFileKeyOff, in.Key)
	ob.SetUint64(remoteFileOffsetOff, in.Offset)
	ob.SetUint64(remoteFileFileSizeOff, in.FileSize)
	ob.SetUint64(remoteFileModifiedTimeOff, in.ModifiedTime)
	ob.SetText(remoteFileExtensionOff, in.Extension)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeInfo ---

const (
	volumeInfoFilesOff         = 0
	volumeInfoVersionOff       = 8
	volumeInfoReplicationOff   = 12
	volumeInfoBytesOffsetOff   = 20
	volumeInfoDatFileSizeOff   = 24
	volumeInfoExpireAtSecOff   = 32
	volumeInfoReadOnlyOff      = 40
	volumeInfoEcShardConfigOff = 44
	volumeInfoSize             = 48
)

// VolumeInfo is a zero-copy view into a ZAP-encoded VolumeInfo message.
type VolumeInfo struct{ o zap.Object }

// WrapVolumeInfo parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapVolumeInfo(b []byte) (VolumeInfo, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeInfo{}, err
	}
	return VolumeInfo{o: m.Root()}, nil
}

// FilesLen reports the number of files elements (proto field 1, repeated
// RemoteFile).
func (t VolumeInfo) FilesLen() int { return t.o.List(volumeInfoFilesOff).Len() }

// FileAt returns the i-th files element. The bool is false when i is out of
// range or the element fails to parse.
func (t VolumeInfo) FileAt(i int) (RemoteFile, bool) {
	o := t.o.List(volumeInfoFilesOff).ObjectAt(i)
	if o.IsNull() {
		return RemoteFile{}, false
	}
	return RemoteFile{o: o}, true
}

// Version reads the version field (proto field 2, uint32).
func (t VolumeInfo) Version() uint32 { return t.o.Uint32(volumeInfoVersionOff) }

// Replication reads the replication field (proto field 3, string).
func (t VolumeInfo) Replication() string { return t.o.Text(volumeInfoReplicationOff) }

// BytesOffset reads the bytes_offset field (proto field 4, uint32).
func (t VolumeInfo) BytesOffset() uint32 { return t.o.Uint32(volumeInfoBytesOffsetOff) }

// DatFileSize reads the dat_file_size field (proto field 5, int64).
func (t VolumeInfo) DatFileSize() int64 { return t.o.Int64(volumeInfoDatFileSizeOff) }

// ExpireAtSec reads the expire_at_sec field (proto field 6, uint64).
func (t VolumeInfo) ExpireAtSec() uint64 { return t.o.Uint64(volumeInfoExpireAtSecOff) }

// ReadOnly reads the read_only field (proto field 7, bool).
func (t VolumeInfo) ReadOnly() bool { return t.o.Bool(volumeInfoReadOnlyOff) }

// EcShardConfig reads the ec_shard_config field (proto field 8, message
// EcShardConfig). The bool is false when the field is absent.
func (t VolumeInfo) EcShardConfig() (EcShardConfig, bool) {
	o := t.o.Object(volumeInfoEcShardConfigOff)
	if o.IsNull() {
		return EcShardConfig{}, false
	}
	return EcShardConfig{o: o}, true
}

// VolumeInfoInput collects the field values for NewVolumeInfo. Each Files entry
// is a RemoteFile message's own ZAP buffer (from NewRemoteFile); EcShardConfig
// is the EcShardConfig message's own ZAP buffer (from NewEcShardConfig).
type VolumeInfoInput struct {
	Files         [][]byte
	Version       uint32
	Replication   string
	BytesOffset   uint32
	DatFileSize   int64
	ExpireAtSec   uint64
	ReadOnly      bool
	EcShardConfig []byte
}

// NewVolumeInfo builds a ZAP-encoded VolumeInfo message from in and returns the
// bytes.
func NewVolumeInfo(in VolumeInfoInput) []byte {
	b := zap.NewBuilder(256)
	listOff, listLen := addObjectList(b, in.Files)
	ob := b.StartObject(volumeInfoSize)
	ob.SetList(volumeInfoFilesOff, listOff, listLen)
	ob.SetUint32(volumeInfoVersionOff, in.Version)
	ob.SetText(volumeInfoReplicationOff, in.Replication)
	ob.SetUint32(volumeInfoBytesOffsetOff, in.BytesOffset)
	ob.SetInt64(volumeInfoDatFileSizeOff, in.DatFileSize)
	ob.SetUint64(volumeInfoExpireAtSecOff, in.ExpireAtSec)
	ob.SetBool(volumeInfoReadOnlyOff, in.ReadOnly)
	ob.SetObject(volumeInfoEcShardConfigOff, embedMessage(b, in.EcShardConfig))
	ob.FinishAsRoot()
	return b.Finish()
}

// --- EcShardConfig ---

const (
	ecShardConfigDataShardsOff   = 0
	ecShardConfigParityShardsOff = 4
	ecShardConfigEncodeTsNsOff   = 8
	ecShardConfigSize            = 16
)

// EcShardConfig is a zero-copy view into a ZAP-encoded EcShardConfig message.
type EcShardConfig struct{ o zap.Object }

// WrapEcShardConfig parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapEcShardConfig(b []byte) (EcShardConfig, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return EcShardConfig{}, err
	}
	return EcShardConfig{o: m.Root()}, nil
}

// DataShards reads the data_shards field (proto field 1, uint32).
func (t EcShardConfig) DataShards() uint32 { return t.o.Uint32(ecShardConfigDataShardsOff) }

// ParityShards reads the parity_shards field (proto field 2, uint32).
func (t EcShardConfig) ParityShards() uint32 { return t.o.Uint32(ecShardConfigParityShardsOff) }

// EncodeTsNs reads the encode_ts_ns field (proto field 3, int64).
func (t EcShardConfig) EncodeTsNs() int64 { return t.o.Int64(ecShardConfigEncodeTsNsOff) }

// EcShardConfigInput collects the field values for NewEcShardConfig.
type EcShardConfigInput struct {
	DataShards   uint32
	ParityShards uint32
	EncodeTsNs   int64
}

// NewEcShardConfig builds a ZAP-encoded EcShardConfig message from in and
// returns the bytes.
func NewEcShardConfig(in EcShardConfigInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(ecShardConfigSize)
	ob.SetUint32(ecShardConfigDataShardsOff, in.DataShards)
	ob.SetUint32(ecShardConfigParityShardsOff, in.ParityShards)
	ob.SetInt64(ecShardConfigEncodeTsNsOff, in.EncodeTsNs)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- EcBitrotProtection ---

const (
	ecBitrotProtectionAlgorithmOff     = 0
	ecBitrotProtectionBlockSizeOff     = 4
	ecBitrotProtectionGenerationOff    = 8
	ecBitrotProtectionEcShardConfigOff = 12
	ecBitrotProtectionShardsOff        = 16
	ecBitrotProtectionEncodeUUIDOff    = 24
	ecBitrotProtectionSize             = 32
)

// EcBitrotProtection is a zero-copy view into a ZAP-encoded EcBitrotProtection
// message.
type EcBitrotProtection struct{ o zap.Object }

// WrapEcBitrotProtection parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapEcBitrotProtection(b []byte) (EcBitrotProtection, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return EcBitrotProtection{}, err
	}
	return EcBitrotProtection{o: m.Root()}, nil
}

// Algorithm reads the algorithm field (proto field 1, enum ChecksumAlgorithm).
func (t EcBitrotProtection) Algorithm() uint32 {
	return t.o.Uint32(ecBitrotProtectionAlgorithmOff)
}

// BlockSize reads the block_size field (proto field 2, uint32).
func (t EcBitrotProtection) BlockSize() uint32 {
	return t.o.Uint32(ecBitrotProtectionBlockSizeOff)
}

// Generation reads the generation field (proto field 3, uint32).
func (t EcBitrotProtection) Generation() uint32 {
	return t.o.Uint32(ecBitrotProtectionGenerationOff)
}

// EcShardConfig reads the ec_shard_config field (proto field 4, message
// EcShardConfig). The bool is false when the field is absent.
func (t EcBitrotProtection) EcShardConfig() (EcShardConfig, bool) {
	o := t.o.Object(ecBitrotProtectionEcShardConfigOff)
	if o.IsNull() {
		return EcShardConfig{}, false
	}
	return EcShardConfig{o: o}, true
}

// ShardsLen reports the number of shards elements (proto field 5, repeated
// EcShardChecksums).
func (t EcBitrotProtection) ShardsLen() int { return t.o.List(ecBitrotProtectionShardsOff).Len() }

// ShardAt returns the i-th shards element. The bool is false when i is out of
// range or the element fails to parse.
func (t EcBitrotProtection) ShardAt(i int) (EcShardChecksums, bool) {
	o := t.o.List(ecBitrotProtectionShardsOff).ObjectAt(i)
	if o.IsNull() {
		return EcShardChecksums{}, false
	}
	return EcShardChecksums{o: o}, true
}

// EncodeUUID reads the encode_uuid field (proto field 6, bytes).
func (t EcBitrotProtection) EncodeUUID() []byte {
	return t.o.Bytes(ecBitrotProtectionEncodeUUIDOff)
}

// EcBitrotProtectionInput collects the field values for NewEcBitrotProtection.
// EcShardConfig is the EcShardConfig message's own ZAP buffer; each Shards entry
// is an EcShardChecksums message's own ZAP buffer.
type EcBitrotProtectionInput struct {
	Algorithm     uint32
	BlockSize     uint32
	Generation    uint32
	EcShardConfig []byte
	Shards        [][]byte
	EncodeUUID    []byte
}

// NewEcBitrotProtection builds a ZAP-encoded EcBitrotProtection message from in
// and returns the bytes.
func NewEcBitrotProtection(in EcBitrotProtectionInput) []byte {
	b := zap.NewBuilder(256)
	listOff, listLen := addObjectList(b, in.Shards)
	ob := b.StartObject(ecBitrotProtectionSize)
	ob.SetUint32(ecBitrotProtectionAlgorithmOff, in.Algorithm)
	ob.SetUint32(ecBitrotProtectionBlockSizeOff, in.BlockSize)
	ob.SetUint32(ecBitrotProtectionGenerationOff, in.Generation)
	ob.SetObject(ecBitrotProtectionEcShardConfigOff, embedMessage(b, in.EcShardConfig))
	ob.SetList(ecBitrotProtectionShardsOff, listOff, listLen)
	ob.SetBytes(ecBitrotProtectionEncodeUUIDOff, in.EncodeUUID)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- EcShardChecksums ---

const (
	ecShardChecksumsShardIDOff     = 0
	ecShardChecksumsCoveredSizeOff = 8
	ecShardChecksumsBlockCrc32cOff = 16
	ecShardChecksumsSize           = 24
)

// EcShardChecksums is a zero-copy view into a ZAP-encoded EcShardChecksums
// message.
type EcShardChecksums struct{ o zap.Object }

// WrapEcShardChecksums parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapEcShardChecksums(b []byte) (EcShardChecksums, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return EcShardChecksums{}, err
	}
	return EcShardChecksums{o: m.Root()}, nil
}

// ShardID reads the shard_id field (proto field 1, uint32).
func (t EcShardChecksums) ShardID() uint32 { return t.o.Uint32(ecShardChecksumsShardIDOff) }

// CoveredSize reads the covered_size field (proto field 2, int64).
func (t EcShardChecksums) CoveredSize() int64 { return t.o.Int64(ecShardChecksumsCoveredSizeOff) }

// BlockCrc32c reads the block_crc32c field (proto field 3, bytes).
func (t EcShardChecksums) BlockCrc32c() []byte {
	return t.o.Bytes(ecShardChecksumsBlockCrc32cOff)
}

// EcShardChecksumsInput collects the field values for NewEcShardChecksums.
type EcShardChecksumsInput struct {
	ShardID     uint32
	CoveredSize int64
	BlockCrc32c []byte
}

// NewEcShardChecksums builds a ZAP-encoded EcShardChecksums message from in and
// returns the bytes.
func NewEcShardChecksums(in EcShardChecksumsInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(ecShardChecksumsSize)
	ob.SetUint32(ecShardChecksumsShardIDOff, in.ShardID)
	ob.SetInt64(ecShardChecksumsCoveredSizeOff, in.CoveredSize)
	ob.SetBytes(ecShardChecksumsBlockCrc32cOff, in.BlockCrc32c)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- OldVersionVolumeInfo ---

const (
	oldVersionVolumeInfoFilesOff       = 0
	oldVersionVolumeInfoVersionOff     = 8
	oldVersionVolumeInfoReplicationOff = 12
	oldVersionVolumeInfoBytesOffsetOff = 20
	oldVersionVolumeInfoDatFileSizeOff = 24
	oldVersionVolumeInfoDestroyTimeOff = 32
	oldVersionVolumeInfoReadOnlyOff    = 40
	oldVersionVolumeInfoSize           = 41
)

// OldVersionVolumeInfo is a zero-copy view into a ZAP-encoded
// OldVersionVolumeInfo message.
type OldVersionVolumeInfo struct{ o zap.Object }

// WrapOldVersionVolumeInfo parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapOldVersionVolumeInfo(b []byte) (OldVersionVolumeInfo, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return OldVersionVolumeInfo{}, err
	}
	return OldVersionVolumeInfo{o: m.Root()}, nil
}

// FilesLen reports the number of files elements (proto field 1, repeated
// RemoteFile).
func (t OldVersionVolumeInfo) FilesLen() int {
	return t.o.List(oldVersionVolumeInfoFilesOff).Len()
}

// FileAt returns the i-th files element. The bool is false when i is out of
// range or the element fails to parse.
func (t OldVersionVolumeInfo) FileAt(i int) (RemoteFile, bool) {
	o := t.o.List(oldVersionVolumeInfoFilesOff).ObjectAt(i)
	if o.IsNull() {
		return RemoteFile{}, false
	}
	return RemoteFile{o: o}, true
}

// Version reads the version field (proto field 2, uint32).
func (t OldVersionVolumeInfo) Version() uint32 {
	return t.o.Uint32(oldVersionVolumeInfoVersionOff)
}

// Replication reads the replication field (proto field 3, string).
func (t OldVersionVolumeInfo) Replication() string {
	return t.o.Text(oldVersionVolumeInfoReplicationOff)
}

// BytesOffset reads the BytesOffset field (proto field 4, uint32).
func (t OldVersionVolumeInfo) BytesOffset() uint32 {
	return t.o.Uint32(oldVersionVolumeInfoBytesOffsetOff)
}

// DatFileSize reads the dat_file_size field (proto field 5, int64).
func (t OldVersionVolumeInfo) DatFileSize() int64 {
	return t.o.Int64(oldVersionVolumeInfoDatFileSizeOff)
}

// DestroyTime reads the DestroyTime field (proto field 6, uint64).
func (t OldVersionVolumeInfo) DestroyTime() uint64 {
	return t.o.Uint64(oldVersionVolumeInfoDestroyTimeOff)
}

// ReadOnly reads the read_only field (proto field 7, bool).
func (t OldVersionVolumeInfo) ReadOnly() bool {
	return t.o.Bool(oldVersionVolumeInfoReadOnlyOff)
}

// OldVersionVolumeInfoInput collects the field values for
// NewOldVersionVolumeInfo. Each Files entry is a RemoteFile message's own ZAP
// buffer (from NewRemoteFile).
type OldVersionVolumeInfoInput struct {
	Files       [][]byte
	Version     uint32
	Replication string
	BytesOffset uint32
	DatFileSize int64
	DestroyTime uint64
	ReadOnly    bool
}

// NewOldVersionVolumeInfo builds a ZAP-encoded OldVersionVolumeInfo message from
// in and returns the bytes.
func NewOldVersionVolumeInfo(in OldVersionVolumeInfoInput) []byte {
	b := zap.NewBuilder(256)
	listOff, listLen := addObjectList(b, in.Files)
	ob := b.StartObject(oldVersionVolumeInfoSize)
	ob.SetList(oldVersionVolumeInfoFilesOff, listOff, listLen)
	ob.SetUint32(oldVersionVolumeInfoVersionOff, in.Version)
	ob.SetText(oldVersionVolumeInfoReplicationOff, in.Replication)
	ob.SetUint32(oldVersionVolumeInfoBytesOffsetOff, in.BytesOffset)
	ob.SetInt64(oldVersionVolumeInfoDatFileSizeOff, in.DatFileSize)
	ob.SetUint64(oldVersionVolumeInfoDestroyTimeOff, in.DestroyTime)
	ob.SetBool(oldVersionVolumeInfoReadOnlyOff, in.ReadOnly)
	ob.FinishAsRoot()
	return b.Finish()
}
