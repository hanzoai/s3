// Code generated from volume_server.proto; DO NOT EDIT.

package volume_serverwire

import (
	zap "github.com/zap-proto/go"
)

// --- VolumeCopyRequest ---

const (
	volumeCopyRequestVolumeIDOff        = 0
	volumeCopyRequestCollectionOff      = 8
	volumeCopyRequestReplicationOff     = 16
	volumeCopyRequestTTLOff             = 24
	volumeCopyRequestSourceDataNodeOff  = 32
	volumeCopyRequestDiskTypeOff        = 40
	volumeCopyRequestIoBytePerSecondOff = 48
	volumeCopyRequestSize               = 56
)

// VolumeCopyRequest is a zero-copy view into a ZAP-encoded VolumeCopyRequest
// message.
type VolumeCopyRequest struct{ o zap.Object }

// WrapVolumeCopyRequest parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapVolumeCopyRequest(b []byte) (VolumeCopyRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeCopyRequest{}, err
	}
	return VolumeCopyRequest{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t VolumeCopyRequest) VolumeID() uint32 { return t.o.Uint32(volumeCopyRequestVolumeIDOff) }

// Collection reads the collection field (proto field 2, string).
func (t VolumeCopyRequest) Collection() string { return t.o.Text(volumeCopyRequestCollectionOff) }

// Replication reads the replication field (proto field 3, string).
func (t VolumeCopyRequest) Replication() string { return t.o.Text(volumeCopyRequestReplicationOff) }

// TTL reads the ttl field (proto field 4, string).
func (t VolumeCopyRequest) TTL() string { return t.o.Text(volumeCopyRequestTTLOff) }

// SourceDataNode reads the source_data_node field (proto field 5, string).
func (t VolumeCopyRequest) SourceDataNode() string {
	return t.o.Text(volumeCopyRequestSourceDataNodeOff)
}

// DiskType reads the disk_type field (proto field 6, string).
func (t VolumeCopyRequest) DiskType() string { return t.o.Text(volumeCopyRequestDiskTypeOff) }

// IoBytePerSecond reads the io_byte_per_second field (proto field 7, int64).
func (t VolumeCopyRequest) IoBytePerSecond() int64 {
	return t.o.Int64(volumeCopyRequestIoBytePerSecondOff)
}

// VolumeCopyRequestInput collects the field values for NewVolumeCopyRequest.
type VolumeCopyRequestInput struct {
	VolumeID        uint32
	Collection      string
	Replication     string
	TTL             string
	SourceDataNode  string
	DiskType        string
	IoBytePerSecond int64
}

// NewVolumeCopyRequest builds a ZAP-encoded VolumeCopyRequest message from in and
// returns the bytes.
func NewVolumeCopyRequest(in VolumeCopyRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeCopyRequestSize)
	ob.SetUint32(volumeCopyRequestVolumeIDOff, in.VolumeID)
	ob.SetText(volumeCopyRequestCollectionOff, in.Collection)
	ob.SetText(volumeCopyRequestReplicationOff, in.Replication)
	ob.SetText(volumeCopyRequestTTLOff, in.TTL)
	ob.SetText(volumeCopyRequestSourceDataNodeOff, in.SourceDataNode)
	ob.SetText(volumeCopyRequestDiskTypeOff, in.DiskType)
	ob.SetInt64(volumeCopyRequestIoBytePerSecondOff, in.IoBytePerSecond)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeCopyResponse ---

const (
	volumeCopyResponseLastAppendAtNsOff = 0
	volumeCopyResponseProcessedBytesOff = 8
	volumeCopyResponseSize              = 16
)

// VolumeCopyResponse is a zero-copy view into a ZAP-encoded VolumeCopyResponse
// message.
type VolumeCopyResponse struct{ o zap.Object }

// WrapVolumeCopyResponse parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapVolumeCopyResponse(b []byte) (VolumeCopyResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeCopyResponse{}, err
	}
	return VolumeCopyResponse{o: m.Root()}, nil
}

// LastAppendAtNs reads the last_append_at_ns field (proto field 1, uint64).
func (t VolumeCopyResponse) LastAppendAtNs() uint64 {
	return t.o.Uint64(volumeCopyResponseLastAppendAtNsOff)
}

// ProcessedBytes reads the processed_bytes field (proto field 2, int64).
func (t VolumeCopyResponse) ProcessedBytes() int64 {
	return t.o.Int64(volumeCopyResponseProcessedBytesOff)
}

// VolumeCopyResponseInput collects the field values for NewVolumeCopyResponse.
type VolumeCopyResponseInput struct {
	LastAppendAtNs uint64
	ProcessedBytes int64
}

// NewVolumeCopyResponse builds a ZAP-encoded VolumeCopyResponse message from in
// and returns the bytes.
func NewVolumeCopyResponse(in VolumeCopyResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeCopyResponseSize)
	ob.SetUint64(volumeCopyResponseLastAppendAtNsOff, in.LastAppendAtNs)
	ob.SetInt64(volumeCopyResponseProcessedBytesOff, in.ProcessedBytes)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- CopyFileRequest ---

const (
	copyFileRequestVolumeIDOff                 = 0
	copyFileRequestExtOff                      = 8
	copyFileRequestCompactionRevisionOff       = 16
	copyFileRequestStopOffsetOff               = 24
	copyFileRequestCollectionOff               = 32
	copyFileRequestIsEcVolumeOff               = 40
	copyFileRequestIgnoreSourceFileNotFoundOff = 41
	copyFileRequestSize                        = 42
)

// CopyFileRequest is a zero-copy view into a ZAP-encoded CopyFileRequest
// message.
type CopyFileRequest struct{ o zap.Object }

// WrapCopyFileRequest parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapCopyFileRequest(b []byte) (CopyFileRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return CopyFileRequest{}, err
	}
	return CopyFileRequest{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t CopyFileRequest) VolumeID() uint32 { return t.o.Uint32(copyFileRequestVolumeIDOff) }

// Ext reads the ext field (proto field 2, string).
func (t CopyFileRequest) Ext() string { return t.o.Text(copyFileRequestExtOff) }

// CompactionRevision reads the compaction_revision field (proto field 3,
// uint32).
func (t CopyFileRequest) CompactionRevision() uint32 {
	return t.o.Uint32(copyFileRequestCompactionRevisionOff)
}

// StopOffset reads the stop_offset field (proto field 4, uint64).
func (t CopyFileRequest) StopOffset() uint64 { return t.o.Uint64(copyFileRequestStopOffsetOff) }

// Collection reads the collection field (proto field 5, string).
func (t CopyFileRequest) Collection() string { return t.o.Text(copyFileRequestCollectionOff) }

// IsEcVolume reads the is_ec_volume field (proto field 6, bool).
func (t CopyFileRequest) IsEcVolume() bool { return t.o.Bool(copyFileRequestIsEcVolumeOff) }

// IgnoreSourceFileNotFound reads the ignore_source_file_not_found field (proto
// field 7, bool).
func (t CopyFileRequest) IgnoreSourceFileNotFound() bool {
	return t.o.Bool(copyFileRequestIgnoreSourceFileNotFoundOff)
}

// CopyFileRequestInput collects the field values for NewCopyFileRequest.
type CopyFileRequestInput struct {
	VolumeID                 uint32
	Ext                      string
	CompactionRevision       uint32
	StopOffset               uint64
	Collection               string
	IsEcVolume               bool
	IgnoreSourceFileNotFound bool
}

// NewCopyFileRequest builds a ZAP-encoded CopyFileRequest message from in and
// returns the bytes.
func NewCopyFileRequest(in CopyFileRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(copyFileRequestSize)
	ob.SetUint32(copyFileRequestVolumeIDOff, in.VolumeID)
	ob.SetText(copyFileRequestExtOff, in.Ext)
	ob.SetUint32(copyFileRequestCompactionRevisionOff, in.CompactionRevision)
	ob.SetUint64(copyFileRequestStopOffsetOff, in.StopOffset)
	ob.SetText(copyFileRequestCollectionOff, in.Collection)
	ob.SetBool(copyFileRequestIsEcVolumeOff, in.IsEcVolume)
	ob.SetBool(copyFileRequestIgnoreSourceFileNotFoundOff, in.IgnoreSourceFileNotFound)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- CopyFileResponse ---

const (
	copyFileResponseFileContentOff  = 0
	copyFileResponseModifiedTsNsOff = 8
	copyFileResponseSize            = 16
)

// CopyFileResponse is a zero-copy view into a ZAP-encoded CopyFileResponse
// message.
type CopyFileResponse struct{ o zap.Object }

// WrapCopyFileResponse parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapCopyFileResponse(b []byte) (CopyFileResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return CopyFileResponse{}, err
	}
	return CopyFileResponse{o: m.Root()}, nil
}

// FileContent reads the file_content field (proto field 1, bytes).
func (t CopyFileResponse) FileContent() []byte { return t.o.Bytes(copyFileResponseFileContentOff) }

// ModifiedTsNs reads the modified_ts_ns field (proto field 2, int64).
func (t CopyFileResponse) ModifiedTsNs() int64 {
	return t.o.Int64(copyFileResponseModifiedTsNsOff)
}

// CopyFileResponseInput collects the field values for NewCopyFileResponse.
type CopyFileResponseInput struct {
	FileContent  []byte
	ModifiedTsNs int64
}

// NewCopyFileResponse builds a ZAP-encoded CopyFileResponse message from in and
// returns the bytes.
func NewCopyFileResponse(in CopyFileResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(copyFileResponseSize)
	ob.SetBytes(copyFileResponseFileContentOff, in.FileContent)
	ob.SetInt64(copyFileResponseModifiedTsNsOff, in.ModifiedTsNs)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ReceiveFileInfo ---

const (
	receiveFileInfoVolumeIDOff   = 0
	receiveFileInfoExtOff        = 8
	receiveFileInfoCollectionOff = 16
	receiveFileInfoIsEcVolumeOff = 24
	receiveFileInfoShardIDOff    = 28
	receiveFileInfoFileSizeOff   = 32
	receiveFileInfoDiskIDOff     = 40
	receiveFileInfoSize          = 44
)

// ReceiveFileInfo is a zero-copy view into a ZAP-encoded ReceiveFileInfo
// message.
type ReceiveFileInfo struct{ o zap.Object }

// WrapReceiveFileInfo parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapReceiveFileInfo(b []byte) (ReceiveFileInfo, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ReceiveFileInfo{}, err
	}
	return ReceiveFileInfo{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t ReceiveFileInfo) VolumeID() uint32 { return t.o.Uint32(receiveFileInfoVolumeIDOff) }

// Ext reads the ext field (proto field 2, string).
func (t ReceiveFileInfo) Ext() string { return t.o.Text(receiveFileInfoExtOff) }

// Collection reads the collection field (proto field 3, string).
func (t ReceiveFileInfo) Collection() string { return t.o.Text(receiveFileInfoCollectionOff) }

// IsEcVolume reads the is_ec_volume field (proto field 4, bool).
func (t ReceiveFileInfo) IsEcVolume() bool { return t.o.Bool(receiveFileInfoIsEcVolumeOff) }

// ShardID reads the shard_id field (proto field 5, uint32).
func (t ReceiveFileInfo) ShardID() uint32 { return t.o.Uint32(receiveFileInfoShardIDOff) }

// FileSize reads the file_size field (proto field 6, uint64).
func (t ReceiveFileInfo) FileSize() uint64 { return t.o.Uint64(receiveFileInfoFileSizeOff) }

// DiskID reads the disk_id field (proto field 7, uint32).
func (t ReceiveFileInfo) DiskID() uint32 { return t.o.Uint32(receiveFileInfoDiskIDOff) }

// ReceiveFileInfoInput collects the field values for NewReceiveFileInfo.
type ReceiveFileInfoInput struct {
	VolumeID   uint32
	Ext        string
	Collection string
	IsEcVolume bool
	ShardID    uint32
	FileSize   uint64
	DiskID     uint32
}

// NewReceiveFileInfo builds a ZAP-encoded ReceiveFileInfo message from in and
// returns the bytes.
func NewReceiveFileInfo(in ReceiveFileInfoInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(receiveFileInfoSize)
	ob.SetUint32(receiveFileInfoVolumeIDOff, in.VolumeID)
	ob.SetText(receiveFileInfoExtOff, in.Ext)
	ob.SetText(receiveFileInfoCollectionOff, in.Collection)
	ob.SetBool(receiveFileInfoIsEcVolumeOff, in.IsEcVolume)
	ob.SetUint32(receiveFileInfoShardIDOff, in.ShardID)
	ob.SetUint64(receiveFileInfoFileSizeOff, in.FileSize)
	ob.SetUint32(receiveFileInfoDiskIDOff, in.DiskID)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ReceiveFileRequest (oneof data) ---
//
// STREAMING: ReceiveFile is a client-streaming RPC; ReceiveFileRequest is the
// streamed element. The data oneof is encoded as a uint32 discriminator (Data)
// plus a single variant payload: the info variant stores the ReceiveFileInfo
// child message's ZAP buffer, and file_content stores its raw bytes.

const (
	receiveFileRequestDataOff    = 0 // uint32 discriminator
	receiveFileRequestPayloadOff = 8 // bytes: selected variant payload
	receiveFileRequestSize       = 16
)

// ReceiveFileRequest is a zero-copy view into a ZAP-encoded ReceiveFileRequest
// message.
type ReceiveFileRequest struct{ o zap.Object }

// WrapReceiveFileRequest parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapReceiveFileRequest(b []byte) (ReceiveFileRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ReceiveFileRequest{}, err
	}
	return ReceiveFileRequest{o: m.Root()}, nil
}

// Data reports which oneof variant is set (one of the ReceiveFileData* values).
func (t ReceiveFileRequest) Data() uint32 { return t.o.Uint32(receiveFileRequestDataOff) }

func (t ReceiveFileRequest) payloadIf(which uint32) ([]byte, bool) {
	if t.Data() != which {
		return nil, false
	}
	return t.o.Bytes(receiveFileRequestPayloadOff), true
}

// Info reads the info variant (proto field 1, message ReceiveFileInfo). The bool
// is false when info is not the selected variant.
func (t ReceiveFileRequest) Info() (ReceiveFileInfo, bool) {
	bts, ok := t.payloadIf(ReceiveFileDataInfo)
	if !ok || len(bts) == 0 {
		return ReceiveFileInfo{}, false
	}
	v, err := WrapReceiveFileInfo(bts)
	if err != nil {
		return ReceiveFileInfo{}, false
	}
	return v, true
}

// FileContent reads the file_content variant (proto field 2, bytes). The bool is
// false when file_content is not the selected variant.
func (t ReceiveFileRequest) FileContent() ([]byte, bool) {
	return t.payloadIf(ReceiveFileDataFileContent)
}

// ReceiveFileRequestInput collects the field values for NewReceiveFileRequest.
// Set Data to the chosen variant and supply exactly that variant's input: Info
// is the ReceiveFileInfo message's own ZAP buffer (from NewReceiveFileInfo),
// FileContent is the raw bytes payload.
type ReceiveFileRequestInput struct {
	Data        uint32
	Info        []byte
	FileContent []byte
}

// NewReceiveFileRequest builds a ZAP-encoded ReceiveFileRequest message from in
// and returns the bytes.
func NewReceiveFileRequest(in ReceiveFileRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(receiveFileRequestSize)
	ob.SetUint32(receiveFileRequestDataOff, in.Data)
	switch in.Data {
	case ReceiveFileDataInfo:
		ob.SetBytes(receiveFileRequestPayloadOff, in.Info)
	case ReceiveFileDataFileContent:
		ob.SetBytes(receiveFileRequestPayloadOff, in.FileContent)
	}
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ReceiveFileResponse ---

const (
	receiveFileResponseBytesWrittenOff = 0
	receiveFileResponseErrorOff        = 8
	receiveFileResponseSize            = 16
)

// ReceiveFileResponse is a zero-copy view into a ZAP-encoded ReceiveFileResponse
// message.
type ReceiveFileResponse struct{ o zap.Object }

// WrapReceiveFileResponse parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapReceiveFileResponse(b []byte) (ReceiveFileResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ReceiveFileResponse{}, err
	}
	return ReceiveFileResponse{o: m.Root()}, nil
}

// BytesWritten reads the bytes_written field (proto field 1, uint64).
func (t ReceiveFileResponse) BytesWritten() uint64 {
	return t.o.Uint64(receiveFileResponseBytesWrittenOff)
}

// Error reads the error field (proto field 2, string).
func (t ReceiveFileResponse) Error() string { return t.o.Text(receiveFileResponseErrorOff) }

// ReceiveFileResponseInput collects the field values for
// NewReceiveFileResponse.
type ReceiveFileResponseInput struct {
	BytesWritten uint64
	Error        string
}

// NewReceiveFileResponse builds a ZAP-encoded ReceiveFileResponse message from
// in and returns the bytes.
func NewReceiveFileResponse(in ReceiveFileResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(receiveFileResponseSize)
	ob.SetUint64(receiveFileResponseBytesWrittenOff, in.BytesWritten)
	ob.SetText(receiveFileResponseErrorOff, in.Error)
	ob.FinishAsRoot()
	return b.Finish()
}
