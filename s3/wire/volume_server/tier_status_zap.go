// Code generated from volume_server.proto; DO NOT EDIT.

package volume_serverwire

import (
	zap "github.com/zap-proto/go"
)

// --- VolumeTierMoveDatToRemoteRequest ---

const (
	volumeTierMoveDatToRemoteRequestVolumeIDOff               = 0
	volumeTierMoveDatToRemoteRequestCollectionOff             = 8
	volumeTierMoveDatToRemoteRequestDestinationBackendNameOff = 16
	volumeTierMoveDatToRemoteRequestKeepLocalDatFileOff       = 24
	volumeTierMoveDatToRemoteRequestSize                      = 25
)

// VolumeTierMoveDatToRemoteRequest is a zero-copy view into a ZAP-encoded
// VolumeTierMoveDatToRemoteRequest message.
type VolumeTierMoveDatToRemoteRequest struct{ o zap.Object }

// WrapVolumeTierMoveDatToRemoteRequest parses b and returns a typed view.
// Returns an error if the wire-level checks (magic, version, size) fail.
func WrapVolumeTierMoveDatToRemoteRequest(b []byte) (VolumeTierMoveDatToRemoteRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeTierMoveDatToRemoteRequest{}, err
	}
	return VolumeTierMoveDatToRemoteRequest{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t VolumeTierMoveDatToRemoteRequest) VolumeID() uint32 {
	return t.o.Uint32(volumeTierMoveDatToRemoteRequestVolumeIDOff)
}

// Collection reads the collection field (proto field 2, string).
func (t VolumeTierMoveDatToRemoteRequest) Collection() string {
	return t.o.Text(volumeTierMoveDatToRemoteRequestCollectionOff)
}

// DestinationBackendName reads the destination_backend_name field (proto field
// 3, string).
func (t VolumeTierMoveDatToRemoteRequest) DestinationBackendName() string {
	return t.o.Text(volumeTierMoveDatToRemoteRequestDestinationBackendNameOff)
}

// KeepLocalDatFile reads the keep_local_dat_file field (proto field 4, bool).
func (t VolumeTierMoveDatToRemoteRequest) KeepLocalDatFile() bool {
	return t.o.Bool(volumeTierMoveDatToRemoteRequestKeepLocalDatFileOff)
}

// VolumeTierMoveDatToRemoteRequestInput collects the field values for
// NewVolumeTierMoveDatToRemoteRequest.
type VolumeTierMoveDatToRemoteRequestInput struct {
	VolumeID               uint32
	Collection             string
	DestinationBackendName string
	KeepLocalDatFile       bool
}

// NewVolumeTierMoveDatToRemoteRequest builds a ZAP-encoded
// VolumeTierMoveDatToRemoteRequest message from in and returns the bytes.
func NewVolumeTierMoveDatToRemoteRequest(in VolumeTierMoveDatToRemoteRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeTierMoveDatToRemoteRequestSize)
	ob.SetUint32(volumeTierMoveDatToRemoteRequestVolumeIDOff, in.VolumeID)
	ob.SetText(volumeTierMoveDatToRemoteRequestCollectionOff, in.Collection)
	ob.SetText(volumeTierMoveDatToRemoteRequestDestinationBackendNameOff, in.DestinationBackendName)
	ob.SetBool(volumeTierMoveDatToRemoteRequestKeepLocalDatFileOff, in.KeepLocalDatFile)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeTierMoveDatToRemoteResponse ---

const (
	volumeTierMoveDatToRemoteResponseProcessedOff           = 0
	volumeTierMoveDatToRemoteResponseProcessedPercentageOff = 8
	volumeTierMoveDatToRemoteResponseSize                   = 12
)

// VolumeTierMoveDatToRemoteResponse is a zero-copy view into a ZAP-encoded
// VolumeTierMoveDatToRemoteResponse message.
type VolumeTierMoveDatToRemoteResponse struct{ o zap.Object }

// WrapVolumeTierMoveDatToRemoteResponse parses b and returns a typed view.
// Returns an error if the wire-level checks (magic, version, size) fail.
func WrapVolumeTierMoveDatToRemoteResponse(b []byte) (VolumeTierMoveDatToRemoteResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeTierMoveDatToRemoteResponse{}, err
	}
	return VolumeTierMoveDatToRemoteResponse{o: m.Root()}, nil
}

// Processed reads the processed field (proto field 1, int64).
func (t VolumeTierMoveDatToRemoteResponse) Processed() int64 {
	return t.o.Int64(volumeTierMoveDatToRemoteResponseProcessedOff)
}

// ProcessedPercentage reads the processedPercentage field (proto field 2,
// float).
func (t VolumeTierMoveDatToRemoteResponse) ProcessedPercentage() float32 {
	return t.o.Float32(volumeTierMoveDatToRemoteResponseProcessedPercentageOff)
}

// VolumeTierMoveDatToRemoteResponseInput collects the field values for
// NewVolumeTierMoveDatToRemoteResponse.
type VolumeTierMoveDatToRemoteResponseInput struct {
	Processed           int64
	ProcessedPercentage float32
}

// NewVolumeTierMoveDatToRemoteResponse builds a ZAP-encoded
// VolumeTierMoveDatToRemoteResponse message from in and returns the bytes.
func NewVolumeTierMoveDatToRemoteResponse(in VolumeTierMoveDatToRemoteResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeTierMoveDatToRemoteResponseSize)
	ob.SetInt64(volumeTierMoveDatToRemoteResponseProcessedOff, in.Processed)
	ob.SetFloat32(volumeTierMoveDatToRemoteResponseProcessedPercentageOff, in.ProcessedPercentage)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeTierMoveDatFromRemoteRequest ---

const (
	volumeTierMoveDatFromRemoteRequestVolumeIDOff          = 0
	volumeTierMoveDatFromRemoteRequestCollectionOff        = 8
	volumeTierMoveDatFromRemoteRequestKeepRemoteDatFileOff = 16
	volumeTierMoveDatFromRemoteRequestSize                 = 17
)

// VolumeTierMoveDatFromRemoteRequest is a zero-copy view into a ZAP-encoded
// VolumeTierMoveDatFromRemoteRequest message.
type VolumeTierMoveDatFromRemoteRequest struct{ o zap.Object }

// WrapVolumeTierMoveDatFromRemoteRequest parses b and returns a typed view.
// Returns an error if the wire-level checks (magic, version, size) fail.
func WrapVolumeTierMoveDatFromRemoteRequest(b []byte) (VolumeTierMoveDatFromRemoteRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeTierMoveDatFromRemoteRequest{}, err
	}
	return VolumeTierMoveDatFromRemoteRequest{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t VolumeTierMoveDatFromRemoteRequest) VolumeID() uint32 {
	return t.o.Uint32(volumeTierMoveDatFromRemoteRequestVolumeIDOff)
}

// Collection reads the collection field (proto field 2, string).
func (t VolumeTierMoveDatFromRemoteRequest) Collection() string {
	return t.o.Text(volumeTierMoveDatFromRemoteRequestCollectionOff)
}

// KeepRemoteDatFile reads the keep_remote_dat_file field (proto field 3, bool).
func (t VolumeTierMoveDatFromRemoteRequest) KeepRemoteDatFile() bool {
	return t.o.Bool(volumeTierMoveDatFromRemoteRequestKeepRemoteDatFileOff)
}

// VolumeTierMoveDatFromRemoteRequestInput collects the field values for
// NewVolumeTierMoveDatFromRemoteRequest.
type VolumeTierMoveDatFromRemoteRequestInput struct {
	VolumeID          uint32
	Collection        string
	KeepRemoteDatFile bool
}

// NewVolumeTierMoveDatFromRemoteRequest builds a ZAP-encoded
// VolumeTierMoveDatFromRemoteRequest message from in and returns the bytes.
func NewVolumeTierMoveDatFromRemoteRequest(in VolumeTierMoveDatFromRemoteRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeTierMoveDatFromRemoteRequestSize)
	ob.SetUint32(volumeTierMoveDatFromRemoteRequestVolumeIDOff, in.VolumeID)
	ob.SetText(volumeTierMoveDatFromRemoteRequestCollectionOff, in.Collection)
	ob.SetBool(volumeTierMoveDatFromRemoteRequestKeepRemoteDatFileOff, in.KeepRemoteDatFile)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeTierMoveDatFromRemoteResponse ---

const (
	volumeTierMoveDatFromRemoteResponseProcessedOff           = 0
	volumeTierMoveDatFromRemoteResponseProcessedPercentageOff = 8
	volumeTierMoveDatFromRemoteResponseSize                   = 12
)

// VolumeTierMoveDatFromRemoteResponse is a zero-copy view into a ZAP-encoded
// VolumeTierMoveDatFromRemoteResponse message.
type VolumeTierMoveDatFromRemoteResponse struct{ o zap.Object }

// WrapVolumeTierMoveDatFromRemoteResponse parses b and returns a typed view.
// Returns an error if the wire-level checks (magic, version, size) fail.
func WrapVolumeTierMoveDatFromRemoteResponse(b []byte) (VolumeTierMoveDatFromRemoteResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeTierMoveDatFromRemoteResponse{}, err
	}
	return VolumeTierMoveDatFromRemoteResponse{o: m.Root()}, nil
}

// Processed reads the processed field (proto field 1, int64).
func (t VolumeTierMoveDatFromRemoteResponse) Processed() int64 {
	return t.o.Int64(volumeTierMoveDatFromRemoteResponseProcessedOff)
}

// ProcessedPercentage reads the processedPercentage field (proto field 2,
// float).
func (t VolumeTierMoveDatFromRemoteResponse) ProcessedPercentage() float32 {
	return t.o.Float32(volumeTierMoveDatFromRemoteResponseProcessedPercentageOff)
}

// VolumeTierMoveDatFromRemoteResponseInput collects the field values for
// NewVolumeTierMoveDatFromRemoteResponse.
type VolumeTierMoveDatFromRemoteResponseInput struct {
	Processed           int64
	ProcessedPercentage float32
}

// NewVolumeTierMoveDatFromRemoteResponse builds a ZAP-encoded
// VolumeTierMoveDatFromRemoteResponse message from in and returns the bytes.
func NewVolumeTierMoveDatFromRemoteResponse(in VolumeTierMoveDatFromRemoteResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeTierMoveDatFromRemoteResponseSize)
	ob.SetInt64(volumeTierMoveDatFromRemoteResponseProcessedOff, in.Processed)
	ob.SetFloat32(volumeTierMoveDatFromRemoteResponseProcessedPercentageOff, in.ProcessedPercentage)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeServerStatusRequest ---

const volumeServerStatusRequestSize = 0

// VolumeServerStatusRequest is a zero-copy view into a ZAP-encoded
// VolumeServerStatusRequest message. VolumeServerStatusRequest carries no
// fields.
type VolumeServerStatusRequest struct{ o zap.Object }

// WrapVolumeServerStatusRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVolumeServerStatusRequest(b []byte) (VolumeServerStatusRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeServerStatusRequest{}, err
	}
	return VolumeServerStatusRequest{o: m.Root()}, nil
}

// VolumeServerStatusRequestInput collects the field values for
// NewVolumeServerStatusRequest. VolumeServerStatusRequest has no fields.
type VolumeServerStatusRequestInput struct{}

// NewVolumeServerStatusRequest builds a ZAP-encoded VolumeServerStatusRequest
// message and returns the bytes.
func NewVolumeServerStatusRequest(in VolumeServerStatusRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeServerStatusRequestSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeServerStatusResponse ---

const (
	volumeServerStatusResponseDiskStatusesOff = 0
	volumeServerStatusResponseMemoryStatusOff = 8
	volumeServerStatusResponseVersionOff      = 12
	volumeServerStatusResponseDataCenterOff   = 20
	volumeServerStatusResponseRackOff         = 28
	volumeServerStatusResponseStateOff        = 36
	volumeServerStatusResponseSize            = 40
)

// VolumeServerStatusResponse is a zero-copy view into a ZAP-encoded
// VolumeServerStatusResponse message.
type VolumeServerStatusResponse struct{ o zap.Object }

// WrapVolumeServerStatusResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVolumeServerStatusResponse(b []byte) (VolumeServerStatusResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeServerStatusResponse{}, err
	}
	return VolumeServerStatusResponse{o: m.Root()}, nil
}

// DiskStatusesLen reports the number of disk_statuses elements (proto field 1,
// repeated DiskStatus).
func (t VolumeServerStatusResponse) DiskStatusesLen() int {
	return t.o.List(volumeServerStatusResponseDiskStatusesOff).Len()
}

// DiskStatusAt returns the i-th disk_statuses element. The bool is false when i
// is out of range or the element fails to parse.
func (t VolumeServerStatusResponse) DiskStatusAt(i int) (DiskStatus, bool) {
	o := t.o.List(volumeServerStatusResponseDiskStatusesOff).ObjectAt(i)
	if o.IsNull() {
		return DiskStatus{}, false
	}
	return DiskStatus{o: o}, true
}

// MemoryStatus reads the memory_status field (proto field 2, message
// MemStatus). The bool is false when the field is absent.
func (t VolumeServerStatusResponse) MemoryStatus() (MemStatus, bool) {
	o := t.o.Object(volumeServerStatusResponseMemoryStatusOff)
	if o.IsNull() {
		return MemStatus{}, false
	}
	return MemStatus{o: o}, true
}

// Version reads the version field (proto field 3, string).
func (t VolumeServerStatusResponse) Version() string {
	return t.o.Text(volumeServerStatusResponseVersionOff)
}

// DataCenter reads the data_center field (proto field 4, string).
func (t VolumeServerStatusResponse) DataCenter() string {
	return t.o.Text(volumeServerStatusResponseDataCenterOff)
}

// Rack reads the rack field (proto field 5, string).
func (t VolumeServerStatusResponse) Rack() string {
	return t.o.Text(volumeServerStatusResponseRackOff)
}

// State reads the state field (proto field 6, message VolumeServerState). The
// bool is false when the field is absent.
func (t VolumeServerStatusResponse) State() (VolumeServerState, bool) {
	o := t.o.Object(volumeServerStatusResponseStateOff)
	if o.IsNull() {
		return VolumeServerState{}, false
	}
	return VolumeServerState{o: o}, true
}

// VolumeServerStatusResponseInput collects the field values for
// NewVolumeServerStatusResponse. Each DiskStatuses entry is a DiskStatus
// message's own ZAP buffer; MemoryStatus and State are their messages' own ZAP
// buffers.
type VolumeServerStatusResponseInput struct {
	DiskStatuses [][]byte
	MemoryStatus []byte
	Version      string
	DataCenter   string
	Rack         string
	State        []byte
}

// NewVolumeServerStatusResponse builds a ZAP-encoded VolumeServerStatusResponse
// message from in and returns the bytes.
func NewVolumeServerStatusResponse(in VolumeServerStatusResponseInput) []byte {
	b := zap.NewBuilder(256)
	listOff, listLen := addObjectList(b, in.DiskStatuses)
	ob := b.StartObject(volumeServerStatusResponseSize)
	ob.SetList(volumeServerStatusResponseDiskStatusesOff, listOff, listLen)
	ob.SetObject(volumeServerStatusResponseMemoryStatusOff, embedMessage(b, in.MemoryStatus))
	ob.SetText(volumeServerStatusResponseVersionOff, in.Version)
	ob.SetText(volumeServerStatusResponseDataCenterOff, in.DataCenter)
	ob.SetText(volumeServerStatusResponseRackOff, in.Rack)
	ob.SetObject(volumeServerStatusResponseStateOff, embedMessage(b, in.State))
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeServerLeaveRequest ---

const volumeServerLeaveRequestSize = 0

// VolumeServerLeaveRequest is a zero-copy view into a ZAP-encoded
// VolumeServerLeaveRequest message. VolumeServerLeaveRequest carries no fields.
type VolumeServerLeaveRequest struct{ o zap.Object }

// WrapVolumeServerLeaveRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVolumeServerLeaveRequest(b []byte) (VolumeServerLeaveRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeServerLeaveRequest{}, err
	}
	return VolumeServerLeaveRequest{o: m.Root()}, nil
}

// VolumeServerLeaveRequestInput collects the field values for
// NewVolumeServerLeaveRequest. VolumeServerLeaveRequest has no fields.
type VolumeServerLeaveRequestInput struct{}

// NewVolumeServerLeaveRequest builds a ZAP-encoded VolumeServerLeaveRequest
// message and returns the bytes.
func NewVolumeServerLeaveRequest(in VolumeServerLeaveRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeServerLeaveRequestSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeServerLeaveResponse ---

const volumeServerLeaveResponseSize = 0

// VolumeServerLeaveResponse is a zero-copy view into a ZAP-encoded
// VolumeServerLeaveResponse message. VolumeServerLeaveResponse carries no
// fields.
type VolumeServerLeaveResponse struct{ o zap.Object }

// WrapVolumeServerLeaveResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapVolumeServerLeaveResponse(b []byte) (VolumeServerLeaveResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeServerLeaveResponse{}, err
	}
	return VolumeServerLeaveResponse{o: m.Root()}, nil
}

// VolumeServerLeaveResponseInput collects the field values for
// NewVolumeServerLeaveResponse. VolumeServerLeaveResponse has no fields.
type VolumeServerLeaveResponseInput struct{}

// NewVolumeServerLeaveResponse builds a ZAP-encoded VolumeServerLeaveResponse
// message and returns the bytes.
func NewVolumeServerLeaveResponse(in VolumeServerLeaveResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeServerLeaveResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}
