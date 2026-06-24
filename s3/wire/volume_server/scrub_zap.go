// Code generated from volume_server.proto; DO NOT EDIT.

package volume_serverwire

import (
	zap "github.com/zap-proto/go"
)

// --- ScrubVolumeRequest ---

const (
	scrubVolumeRequestModeOff                      = 0
	scrubVolumeRequestVolumeIDsOff                 = 8
	scrubVolumeRequestMarkBrokenVolumesReadonlyOff = 16
	scrubVolumeRequestSize                         = 17
)

// ScrubVolumeRequest is a zero-copy view into a ZAP-encoded ScrubVolumeRequest
// message.
type ScrubVolumeRequest struct{ o zap.Object }

// WrapScrubVolumeRequest parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapScrubVolumeRequest(b []byte) (ScrubVolumeRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ScrubVolumeRequest{}, err
	}
	return ScrubVolumeRequest{o: m.Root()}, nil
}

// Mode reads the mode field (proto field 1, enum VolumeScrubMode).
func (t ScrubVolumeRequest) Mode() uint32 { return t.o.Uint32(scrubVolumeRequestModeOff) }

// VolumeIDsLen reports the number of volume_ids elements (proto field 2,
// repeated uint32).
func (t ScrubVolumeRequest) VolumeIDsLen() int {
	return t.o.List(scrubVolumeRequestVolumeIDsOff).Len()
}

// VolumeIDAt returns the i-th volume_ids element (uint32).
func (t ScrubVolumeRequest) VolumeIDAt(i int) uint32 {
	return t.o.List(scrubVolumeRequestVolumeIDsOff).Uint32(i)
}

// MarkBrokenVolumesReadonly reads the mark_broken_volumes_readonly field (proto
// field 3, bool).
func (t ScrubVolumeRequest) MarkBrokenVolumesReadonly() bool {
	return t.o.Bool(scrubVolumeRequestMarkBrokenVolumesReadonlyOff)
}

// ScrubVolumeRequestInput collects the field values for NewScrubVolumeRequest.
type ScrubVolumeRequestInput struct {
	Mode                      uint32
	VolumeIDs                 []uint32
	MarkBrokenVolumesReadonly bool
}

// NewScrubVolumeRequest builds a ZAP-encoded ScrubVolumeRequest message from in
// and returns the bytes.
func NewScrubVolumeRequest(in ScrubVolumeRequestInput) []byte {
	b := zap.NewBuilder(256)
	listOff, listLen := addUint32List(b, in.VolumeIDs)
	ob := b.StartObject(scrubVolumeRequestSize)
	ob.SetUint32(scrubVolumeRequestModeOff, in.Mode)
	ob.SetList(scrubVolumeRequestVolumeIDsOff, listOff, listLen)
	ob.SetBool(scrubVolumeRequestMarkBrokenVolumesReadonlyOff, in.MarkBrokenVolumesReadonly)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ScrubVolumeResponse ---

const (
	scrubVolumeResponseTotalVolumesOff    = 0
	scrubVolumeResponseTotalFilesOff      = 8
	scrubVolumeResponseBrokenVolumeIDsOff = 16
	scrubVolumeResponseDetailsOff         = 24
	scrubVolumeResponseSize               = 32
)

// ScrubVolumeResponse is a zero-copy view into a ZAP-encoded ScrubVolumeResponse
// message.
type ScrubVolumeResponse struct{ o zap.Object }

// WrapScrubVolumeResponse parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapScrubVolumeResponse(b []byte) (ScrubVolumeResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ScrubVolumeResponse{}, err
	}
	return ScrubVolumeResponse{o: m.Root()}, nil
}

// TotalVolumes reads the total_volumes field (proto field 1, uint64).
func (t ScrubVolumeResponse) TotalVolumes() uint64 {
	return t.o.Uint64(scrubVolumeResponseTotalVolumesOff)
}

// TotalFiles reads the total_files field (proto field 2, uint64).
func (t ScrubVolumeResponse) TotalFiles() uint64 {
	return t.o.Uint64(scrubVolumeResponseTotalFilesOff)
}

// BrokenVolumeIDsLen reports the number of broken_volume_ids elements (proto
// field 3, repeated uint32).
func (t ScrubVolumeResponse) BrokenVolumeIDsLen() int {
	return t.o.List(scrubVolumeResponseBrokenVolumeIDsOff).Len()
}

// BrokenVolumeIDAt returns the i-th broken_volume_ids element (uint32).
func (t ScrubVolumeResponse) BrokenVolumeIDAt(i int) uint32 {
	return t.o.List(scrubVolumeResponseBrokenVolumeIDsOff).Uint32(i)
}

// DetailsLen reports the number of details elements (proto field 4, repeated
// string).
func (t ScrubVolumeResponse) DetailsLen() int {
	return t.o.List(scrubVolumeResponseDetailsOff).Len()
}

// DetailAt returns the i-th details element (string).
func (t ScrubVolumeResponse) DetailAt(i int) string {
	return elemText(t.o.List(scrubVolumeResponseDetailsOff), i)
}

// ScrubVolumeResponseInput collects the field values for NewScrubVolumeResponse.
type ScrubVolumeResponseInput struct {
	TotalVolumes    uint64
	TotalFiles      uint64
	BrokenVolumeIDs []uint32
	Details         []string
}

// NewScrubVolumeResponse builds a ZAP-encoded ScrubVolumeResponse message from
// in and returns the bytes.
func NewScrubVolumeResponse(in ScrubVolumeResponseInput) []byte {
	b := zap.NewBuilder(256)
	brokenOff, brokenLen := addUint32List(b, in.BrokenVolumeIDs)
	detailsOff, detailsLen := addStringList(b, in.Details)
	ob := b.StartObject(scrubVolumeResponseSize)
	ob.SetUint64(scrubVolumeResponseTotalVolumesOff, in.TotalVolumes)
	ob.SetUint64(scrubVolumeResponseTotalFilesOff, in.TotalFiles)
	ob.SetList(scrubVolumeResponseBrokenVolumeIDsOff, brokenOff, brokenLen)
	ob.SetList(scrubVolumeResponseDetailsOff, detailsOff, detailsLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ScrubEcVolumeRequest ---

const (
	scrubEcVolumeRequestModeOff      = 0
	scrubEcVolumeRequestVolumeIDsOff = 8
	scrubEcVolumeRequestSize         = 16
)

// ScrubEcVolumeRequest is a zero-copy view into a ZAP-encoded
// ScrubEcVolumeRequest message.
type ScrubEcVolumeRequest struct{ o zap.Object }

// WrapScrubEcVolumeRequest parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapScrubEcVolumeRequest(b []byte) (ScrubEcVolumeRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ScrubEcVolumeRequest{}, err
	}
	return ScrubEcVolumeRequest{o: m.Root()}, nil
}

// Mode reads the mode field (proto field 1, enum VolumeScrubMode).
func (t ScrubEcVolumeRequest) Mode() uint32 { return t.o.Uint32(scrubEcVolumeRequestModeOff) }

// VolumeIDsLen reports the number of volume_ids elements (proto field 2,
// repeated uint32).
func (t ScrubEcVolumeRequest) VolumeIDsLen() int {
	return t.o.List(scrubEcVolumeRequestVolumeIDsOff).Len()
}

// VolumeIDAt returns the i-th volume_ids element (uint32).
func (t ScrubEcVolumeRequest) VolumeIDAt(i int) uint32 {
	return t.o.List(scrubEcVolumeRequestVolumeIDsOff).Uint32(i)
}

// ScrubEcVolumeRequestInput collects the field values for
// NewScrubEcVolumeRequest.
type ScrubEcVolumeRequestInput struct {
	Mode      uint32
	VolumeIDs []uint32
}

// NewScrubEcVolumeRequest builds a ZAP-encoded ScrubEcVolumeRequest message from
// in and returns the bytes.
func NewScrubEcVolumeRequest(in ScrubEcVolumeRequestInput) []byte {
	b := zap.NewBuilder(256)
	listOff, listLen := addUint32List(b, in.VolumeIDs)
	ob := b.StartObject(scrubEcVolumeRequestSize)
	ob.SetUint32(scrubEcVolumeRequestModeOff, in.Mode)
	ob.SetList(scrubEcVolumeRequestVolumeIDsOff, listOff, listLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ScrubEcVolumeResponse ---

const (
	scrubEcVolumeResponseTotalVolumesOff     = 0
	scrubEcVolumeResponseTotalFilesOff       = 8
	scrubEcVolumeResponseBrokenVolumeIDsOff  = 16
	scrubEcVolumeResponseBrokenShardInfosOff = 24
	scrubEcVolumeResponseDetailsOff          = 32
	scrubEcVolumeResponseSize                = 40
)

// ScrubEcVolumeResponse is a zero-copy view into a ZAP-encoded
// ScrubEcVolumeResponse message.
type ScrubEcVolumeResponse struct{ o zap.Object }

// WrapScrubEcVolumeResponse parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapScrubEcVolumeResponse(b []byte) (ScrubEcVolumeResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ScrubEcVolumeResponse{}, err
	}
	return ScrubEcVolumeResponse{o: m.Root()}, nil
}

// TotalVolumes reads the total_volumes field (proto field 1, uint64).
func (t ScrubEcVolumeResponse) TotalVolumes() uint64 {
	return t.o.Uint64(scrubEcVolumeResponseTotalVolumesOff)
}

// TotalFiles reads the total_files field (proto field 2, uint64).
func (t ScrubEcVolumeResponse) TotalFiles() uint64 {
	return t.o.Uint64(scrubEcVolumeResponseTotalFilesOff)
}

// BrokenVolumeIDsLen reports the number of broken_volume_ids elements (proto
// field 3, repeated uint32).
func (t ScrubEcVolumeResponse) BrokenVolumeIDsLen() int {
	return t.o.List(scrubEcVolumeResponseBrokenVolumeIDsOff).Len()
}

// BrokenVolumeIDAt returns the i-th broken_volume_ids element (uint32).
func (t ScrubEcVolumeResponse) BrokenVolumeIDAt(i int) uint32 {
	return t.o.List(scrubEcVolumeResponseBrokenVolumeIDsOff).Uint32(i)
}

// BrokenShardInfosLen reports the number of broken_shard_infos elements (proto
// field 4, repeated EcShardInfo).
func (t ScrubEcVolumeResponse) BrokenShardInfosLen() int {
	return t.o.List(scrubEcVolumeResponseBrokenShardInfosOff).Len()
}

// BrokenShardInfoAt returns the i-th broken_shard_infos element. The bool is
// false when i is out of range or the element fails to parse.
func (t ScrubEcVolumeResponse) BrokenShardInfoAt(i int) (EcShardInfo, bool) {
	o := t.o.List(scrubEcVolumeResponseBrokenShardInfosOff).ObjectAt(i)
	if o.IsNull() {
		return EcShardInfo{}, false
	}
	return EcShardInfo{o: o}, true
}

// DetailsLen reports the number of details elements (proto field 5, repeated
// string).
func (t ScrubEcVolumeResponse) DetailsLen() int {
	return t.o.List(scrubEcVolumeResponseDetailsOff).Len()
}

// DetailAt returns the i-th details element (string).
func (t ScrubEcVolumeResponse) DetailAt(i int) string {
	return elemText(t.o.List(scrubEcVolumeResponseDetailsOff), i)
}

// ScrubEcVolumeResponseInput collects the field values for
// NewScrubEcVolumeResponse. Each BrokenShardInfos entry is an EcShardInfo
// message's own ZAP buffer (from NewEcShardInfo).
type ScrubEcVolumeResponseInput struct {
	TotalVolumes     uint64
	TotalFiles       uint64
	BrokenVolumeIDs  []uint32
	BrokenShardInfos [][]byte
	Details          []string
}

// NewScrubEcVolumeResponse builds a ZAP-encoded ScrubEcVolumeResponse message
// from in and returns the bytes.
func NewScrubEcVolumeResponse(in ScrubEcVolumeResponseInput) []byte {
	b := zap.NewBuilder(256)
	brokenOff, brokenLen := addUint32List(b, in.BrokenVolumeIDs)
	shardOff, shardLen := addObjectList(b, in.BrokenShardInfos)
	detailsOff, detailsLen := addStringList(b, in.Details)
	ob := b.StartObject(scrubEcVolumeResponseSize)
	ob.SetUint64(scrubEcVolumeResponseTotalVolumesOff, in.TotalVolumes)
	ob.SetUint64(scrubEcVolumeResponseTotalFilesOff, in.TotalFiles)
	ob.SetList(scrubEcVolumeResponseBrokenVolumeIDsOff, brokenOff, brokenLen)
	ob.SetList(scrubEcVolumeResponseBrokenShardInfosOff, shardOff, shardLen)
	ob.SetList(scrubEcVolumeResponseDetailsOff, detailsOff, detailsLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeNeedleStatusRequest ---

const (
	volumeNeedleStatusRequestVolumeIDOff = 0
	volumeNeedleStatusRequestNeedleIDOff = 8
	volumeNeedleStatusRequestSize        = 16
)

// VolumeNeedleStatusRequest is a zero-copy view into a ZAP-encoded
// VolumeNeedleStatusRequest message.
type VolumeNeedleStatusRequest struct{ o zap.Object }

// WrapVolumeNeedleStatusRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVolumeNeedleStatusRequest(b []byte) (VolumeNeedleStatusRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeNeedleStatusRequest{}, err
	}
	return VolumeNeedleStatusRequest{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t VolumeNeedleStatusRequest) VolumeID() uint32 {
	return t.o.Uint32(volumeNeedleStatusRequestVolumeIDOff)
}

// NeedleID reads the needle_id field (proto field 2, uint64).
func (t VolumeNeedleStatusRequest) NeedleID() uint64 {
	return t.o.Uint64(volumeNeedleStatusRequestNeedleIDOff)
}

// VolumeNeedleStatusRequestInput collects the field values for
// NewVolumeNeedleStatusRequest.
type VolumeNeedleStatusRequestInput struct {
	VolumeID uint32
	NeedleID uint64
}

// NewVolumeNeedleStatusRequest builds a ZAP-encoded VolumeNeedleStatusRequest
// message from in and returns the bytes.
func NewVolumeNeedleStatusRequest(in VolumeNeedleStatusRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeNeedleStatusRequestSize)
	ob.SetUint32(volumeNeedleStatusRequestVolumeIDOff, in.VolumeID)
	ob.SetUint64(volumeNeedleStatusRequestNeedleIDOff, in.NeedleID)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeNeedleStatusResponse ---

const (
	volumeNeedleStatusResponseNeedleIDOff     = 0
	volumeNeedleStatusResponseCookieOff       = 8
	volumeNeedleStatusResponseSizeFOff        = 12
	volumeNeedleStatusResponseLastModifiedOff = 16
	volumeNeedleStatusResponseCrcOff          = 24
	volumeNeedleStatusResponseTTLOff          = 28
	volumeNeedleStatusResponseSize            = 36
)

// VolumeNeedleStatusResponse is a zero-copy view into a ZAP-encoded
// VolumeNeedleStatusResponse message.
type VolumeNeedleStatusResponse struct{ o zap.Object }

// WrapVolumeNeedleStatusResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVolumeNeedleStatusResponse(b []byte) (VolumeNeedleStatusResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeNeedleStatusResponse{}, err
	}
	return VolumeNeedleStatusResponse{o: m.Root()}, nil
}

// NeedleID reads the needle_id field (proto field 1, uint64).
func (t VolumeNeedleStatusResponse) NeedleID() uint64 {
	return t.o.Uint64(volumeNeedleStatusResponseNeedleIDOff)
}

// Cookie reads the cookie field (proto field 2, uint32).
func (t VolumeNeedleStatusResponse) Cookie() uint32 {
	return t.o.Uint32(volumeNeedleStatusResponseCookieOff)
}

// Size reads the size field (proto field 3, uint32).
func (t VolumeNeedleStatusResponse) Size() uint32 {
	return t.o.Uint32(volumeNeedleStatusResponseSizeFOff)
}

// LastModified reads the last_modified field (proto field 4, uint64).
func (t VolumeNeedleStatusResponse) LastModified() uint64 {
	return t.o.Uint64(volumeNeedleStatusResponseLastModifiedOff)
}

// Crc reads the crc field (proto field 5, uint32).
func (t VolumeNeedleStatusResponse) Crc() uint32 {
	return t.o.Uint32(volumeNeedleStatusResponseCrcOff)
}

// TTL reads the ttl field (proto field 6, string).
func (t VolumeNeedleStatusResponse) TTL() string {
	return t.o.Text(volumeNeedleStatusResponseTTLOff)
}

// VolumeNeedleStatusResponseInput collects the field values for
// NewVolumeNeedleStatusResponse.
type VolumeNeedleStatusResponseInput struct {
	NeedleID     uint64
	Cookie       uint32
	Size         uint32
	LastModified uint64
	Crc          uint32
	TTL          string
}

// NewVolumeNeedleStatusResponse builds a ZAP-encoded VolumeNeedleStatusResponse
// message from in and returns the bytes.
func NewVolumeNeedleStatusResponse(in VolumeNeedleStatusResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeNeedleStatusResponseSize)
	ob.SetUint64(volumeNeedleStatusResponseNeedleIDOff, in.NeedleID)
	ob.SetUint32(volumeNeedleStatusResponseCookieOff, in.Cookie)
	ob.SetUint32(volumeNeedleStatusResponseSizeFOff, in.Size)
	ob.SetUint64(volumeNeedleStatusResponseLastModifiedOff, in.LastModified)
	ob.SetUint32(volumeNeedleStatusResponseCrcOff, in.Crc)
	ob.SetText(volumeNeedleStatusResponseTTLOff, in.TTL)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- PingRequest ---

const (
	pingRequestTargetOff     = 0
	pingRequestTargetTypeOff = 8
	pingRequestSize          = 16
)

// PingRequest is a zero-copy view into a ZAP-encoded PingRequest message.
type PingRequest struct{ o zap.Object }

// WrapPingRequest parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapPingRequest(b []byte) (PingRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PingRequest{}, err
	}
	return PingRequest{o: m.Root()}, nil
}

// Target reads the target field (proto field 1, string).
func (t PingRequest) Target() string { return t.o.Text(pingRequestTargetOff) }

// TargetType reads the target_type field (proto field 2, string).
func (t PingRequest) TargetType() string { return t.o.Text(pingRequestTargetTypeOff) }

// PingRequestInput collects the field values for NewPingRequest.
type PingRequestInput struct {
	Target     string
	TargetType string
}

// NewPingRequest builds a ZAP-encoded PingRequest message from in and returns
// the bytes.
func NewPingRequest(in PingRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(pingRequestSize)
	ob.SetText(pingRequestTargetOff, in.Target)
	ob.SetText(pingRequestTargetTypeOff, in.TargetType)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- PingResponse ---

const (
	pingResponseStartTimeNsOff  = 0
	pingResponseRemoteTimeNsOff = 8
	pingResponseStopTimeNsOff   = 16
	pingResponseSize            = 24
)

// PingResponse is a zero-copy view into a ZAP-encoded PingResponse message.
type PingResponse struct{ o zap.Object }

// WrapPingResponse parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapPingResponse(b []byte) (PingResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PingResponse{}, err
	}
	return PingResponse{o: m.Root()}, nil
}

// StartTimeNs reads the start_time_ns field (proto field 1, int64).
func (t PingResponse) StartTimeNs() int64 { return t.o.Int64(pingResponseStartTimeNsOff) }

// RemoteTimeNs reads the remote_time_ns field (proto field 2, int64).
func (t PingResponse) RemoteTimeNs() int64 { return t.o.Int64(pingResponseRemoteTimeNsOff) }

// StopTimeNs reads the stop_time_ns field (proto field 3, int64).
func (t PingResponse) StopTimeNs() int64 { return t.o.Int64(pingResponseStopTimeNsOff) }

// PingResponseInput collects the field values for NewPingResponse.
type PingResponseInput struct {
	StartTimeNs  int64
	RemoteTimeNs int64
	StopTimeNs   int64
}

// NewPingResponse builds a ZAP-encoded PingResponse message from in and returns
// the bytes.
func NewPingResponse(in PingResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(pingResponseSize)
	ob.SetInt64(pingResponseStartTimeNsOff, in.StartTimeNs)
	ob.SetInt64(pingResponseRemoteTimeNsOff, in.RemoteTimeNs)
	ob.SetInt64(pingResponseStopTimeNsOff, in.StopTimeNs)
	ob.FinishAsRoot()
	return b.Finish()
}
