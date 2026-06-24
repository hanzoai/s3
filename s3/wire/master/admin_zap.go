// Code generated from master.proto; DO NOT EDIT.

package masterwire

import (
	zap "github.com/zap-proto/go"
)

// --- GetMasterConfigurationRequest ---

const getMasterConfigurationRequestSize = 0

// GetMasterConfigurationRequest is a zero-copy view into a ZAP-encoded
// GetMasterConfigurationRequest message. It carries no fields.
type GetMasterConfigurationRequest struct{ o zap.Object }

// WrapGetMasterConfigurationRequest parses b and returns a typed view. Returns
// an error if the wire-level checks (magic, version, size) fail.
func WrapGetMasterConfigurationRequest(b []byte) (GetMasterConfigurationRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return GetMasterConfigurationRequest{}, err
	}
	return GetMasterConfigurationRequest{o: m.Root()}, nil
}

// GetMasterConfigurationRequestInput collects the field values for
// NewGetMasterConfigurationRequest. It has no fields.
type GetMasterConfigurationRequestInput struct{}

// NewGetMasterConfigurationRequest builds a ZAP-encoded
// GetMasterConfigurationRequest message and returns the bytes.
func NewGetMasterConfigurationRequest(in GetMasterConfigurationRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(getMasterConfigurationRequestSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- GetMasterConfigurationResponse ---

const (
	getMasterConfigurationResponseMetricsAddressOff          = 0
	getMasterConfigurationResponseMetricsIntervalSecondsOff  = 8
	getMasterConfigurationResponseStorageBackendsOff         = 12
	getMasterConfigurationResponseDefaultReplicationOff      = 20
	getMasterConfigurationResponseLeaderOff                  = 28
	getMasterConfigurationResponseVolumeSizeLimitMBOff       = 36
	getMasterConfigurationResponseVolumePreallocateOff       = 40
	getMasterConfigurationResponseMaintenanceScriptsOff      = 41
	getMasterConfigurationResponseMaintenanceSleepMinutesOff = 49
	getMasterConfigurationResponseSize                       = 53
)

// GetMasterConfigurationResponse is a zero-copy view into a ZAP-encoded
// GetMasterConfigurationResponse message.
type GetMasterConfigurationResponse struct{ o zap.Object }

// WrapGetMasterConfigurationResponse parses b and returns a typed view. Returns
// an error if the wire-level checks (magic, version, size) fail.
func WrapGetMasterConfigurationResponse(b []byte) (GetMasterConfigurationResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return GetMasterConfigurationResponse{}, err
	}
	return GetMasterConfigurationResponse{o: m.Root()}, nil
}

// MetricsAddress reads the metrics_address field (proto field 1, string).
func (t GetMasterConfigurationResponse) MetricsAddress() string {
	return t.o.Text(getMasterConfigurationResponseMetricsAddressOff)
}

// MetricsIntervalSeconds reads the metrics_interval_seconds field (proto field
// 2, uint32).
func (t GetMasterConfigurationResponse) MetricsIntervalSeconds() uint32 {
	return t.o.Uint32(getMasterConfigurationResponseMetricsIntervalSecondsOff)
}

// StorageBackendsLen reports the number of storage_backends elements (proto
// field 3, repeated StorageBackend).
func (t GetMasterConfigurationResponse) StorageBackendsLen() int {
	return t.o.List(getMasterConfigurationResponseStorageBackendsOff).Len()
}

// StorageBackendAt returns the i-th storage_backends element.
func (t GetMasterConfigurationResponse) StorageBackendAt(i int) (StorageBackend, bool) {
	o := t.o.List(getMasterConfigurationResponseStorageBackendsOff).ObjectAt(i)
	if o.IsNull() {
		return StorageBackend{}, false
	}
	return StorageBackend{o: o}, true
}

// DefaultReplication reads the default_replication field (proto field 4,
// string).
func (t GetMasterConfigurationResponse) DefaultReplication() string {
	return t.o.Text(getMasterConfigurationResponseDefaultReplicationOff)
}

// Leader reads the leader field (proto field 5, string).
func (t GetMasterConfigurationResponse) Leader() string {
	return t.o.Text(getMasterConfigurationResponseLeaderOff)
}

// VolumeSizeLimitMB reads the volume_size_limit_m_b field (proto field 6,
// uint32).
func (t GetMasterConfigurationResponse) VolumeSizeLimitMB() uint32 {
	return t.o.Uint32(getMasterConfigurationResponseVolumeSizeLimitMBOff)
}

// VolumePreallocate reads the volume_preallocate field (proto field 7, bool).
func (t GetMasterConfigurationResponse) VolumePreallocate() bool {
	return t.o.Bool(getMasterConfigurationResponseVolumePreallocateOff)
}

// MaintenanceScripts reads the maintenance_scripts field (proto field 8,
// string).
func (t GetMasterConfigurationResponse) MaintenanceScripts() string {
	return t.o.Text(getMasterConfigurationResponseMaintenanceScriptsOff)
}

// MaintenanceSleepMinutes reads the maintenance_sleep_minutes field (proto field
// 9, uint32).
func (t GetMasterConfigurationResponse) MaintenanceSleepMinutes() uint32 {
	return t.o.Uint32(getMasterConfigurationResponseMaintenanceSleepMinutesOff)
}

// GetMasterConfigurationResponseInput collects the field values for
// NewGetMasterConfigurationResponse. StorageBackends takes each element as its
// own ZAP sub-buffer.
type GetMasterConfigurationResponseInput struct {
	MetricsAddress          string
	MetricsIntervalSeconds  uint32
	StorageBackends         [][]byte
	DefaultReplication      string
	Leader                  string
	VolumeSizeLimitMB       uint32
	VolumePreallocate       bool
	MaintenanceScripts      string
	MaintenanceSleepMinutes uint32
}

// NewGetMasterConfigurationResponse builds a ZAP-encoded
// GetMasterConfigurationResponse message from in and returns the bytes.
func NewGetMasterConfigurationResponse(in GetMasterConfigurationResponseInput) []byte {
	b := zap.NewBuilder(256)
	storageBackendsOff, storageBackendsLen := setMsgList(b, in.StorageBackends)
	ob := b.StartObject(getMasterConfigurationResponseSize)
	ob.SetText(getMasterConfigurationResponseMetricsAddressOff, in.MetricsAddress)
	ob.SetUint32(getMasterConfigurationResponseMetricsIntervalSecondsOff, in.MetricsIntervalSeconds)
	ob.SetList(getMasterConfigurationResponseStorageBackendsOff, storageBackendsOff, storageBackendsLen)
	ob.SetText(getMasterConfigurationResponseDefaultReplicationOff, in.DefaultReplication)
	ob.SetText(getMasterConfigurationResponseLeaderOff, in.Leader)
	ob.SetUint32(getMasterConfigurationResponseVolumeSizeLimitMBOff, in.VolumeSizeLimitMB)
	ob.SetBool(getMasterConfigurationResponseVolumePreallocateOff, in.VolumePreallocate)
	ob.SetText(getMasterConfigurationResponseMaintenanceScriptsOff, in.MaintenanceScripts)
	ob.SetUint32(getMasterConfigurationResponseMaintenanceSleepMinutesOff, in.MaintenanceSleepMinutes)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ListClusterNodesRequest ---

const (
	listClusterNodesRequestClientTypeOff = 0
	listClusterNodesRequestFilerGroupOff = 8
	listClusterNodesRequestLimitOff      = 16
	listClusterNodesRequestSize          = 20
)

// ListClusterNodesRequest is a zero-copy view into a ZAP-encoded
// ListClusterNodesRequest message.
type ListClusterNodesRequest struct{ o zap.Object }

// WrapListClusterNodesRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapListClusterNodesRequest(b []byte) (ListClusterNodesRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ListClusterNodesRequest{}, err
	}
	return ListClusterNodesRequest{o: m.Root()}, nil
}

// ClientType reads the client_type field (proto field 1, string).
func (t ListClusterNodesRequest) ClientType() string {
	return t.o.Text(listClusterNodesRequestClientTypeOff)
}

// FilerGroup reads the filer_group field (proto field 2, string).
func (t ListClusterNodesRequest) FilerGroup() string {
	return t.o.Text(listClusterNodesRequestFilerGroupOff)
}

// Limit reads the limit field (proto field 4, int32).
func (t ListClusterNodesRequest) Limit() int32 {
	return t.o.Int32(listClusterNodesRequestLimitOff)
}

// ListClusterNodesRequestInput collects the field values for
// NewListClusterNodesRequest.
type ListClusterNodesRequestInput struct {
	ClientType string
	FilerGroup string
	Limit      int32
}

// NewListClusterNodesRequest builds a ZAP-encoded ListClusterNodesRequest
// message from in and returns the bytes.
func NewListClusterNodesRequest(in ListClusterNodesRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(listClusterNodesRequestSize)
	ob.SetText(listClusterNodesRequestClientTypeOff, in.ClientType)
	ob.SetText(listClusterNodesRequestFilerGroupOff, in.FilerGroup)
	ob.SetInt32(listClusterNodesRequestLimitOff, in.Limit)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ListClusterNodesResponse.ClusterNode ---

const (
	listClusterNodesResponseClusterNodeAddressOff     = 0
	listClusterNodesResponseClusterNodeVersionOff     = 8
	listClusterNodesResponseClusterNodeCreatedAtNsOff = 16
	listClusterNodesResponseClusterNodeDataCenterOff  = 24
	listClusterNodesResponseClusterNodeRackOff        = 32
	listClusterNodesResponseClusterNodeSize           = 40
)

// ListClusterNodesResponseClusterNode is a zero-copy view into a ZAP-encoded
// ListClusterNodesResponse.ClusterNode message.
type ListClusterNodesResponseClusterNode struct{ o zap.Object }

// WrapListClusterNodesResponseClusterNode parses b and returns a typed view.
// Returns an error if the wire-level checks (magic, version, size) fail.
func WrapListClusterNodesResponseClusterNode(b []byte) (ListClusterNodesResponseClusterNode, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ListClusterNodesResponseClusterNode{}, err
	}
	return ListClusterNodesResponseClusterNode{o: m.Root()}, nil
}

// Address reads the address field (proto field 1, string).
func (t ListClusterNodesResponseClusterNode) Address() string {
	return t.o.Text(listClusterNodesResponseClusterNodeAddressOff)
}

// Version reads the version field (proto field 2, string).
func (t ListClusterNodesResponseClusterNode) Version() string {
	return t.o.Text(listClusterNodesResponseClusterNodeVersionOff)
}

// CreatedAtNs reads the created_at_ns field (proto field 4, int64).
func (t ListClusterNodesResponseClusterNode) CreatedAtNs() int64 {
	return t.o.Int64(listClusterNodesResponseClusterNodeCreatedAtNsOff)
}

// DataCenter reads the data_center field (proto field 5, string).
func (t ListClusterNodesResponseClusterNode) DataCenter() string {
	return t.o.Text(listClusterNodesResponseClusterNodeDataCenterOff)
}

// Rack reads the rack field (proto field 6, string).
func (t ListClusterNodesResponseClusterNode) Rack() string {
	return t.o.Text(listClusterNodesResponseClusterNodeRackOff)
}

// ListClusterNodesResponseClusterNodeInput collects the field values for
// NewListClusterNodesResponseClusterNode.
type ListClusterNodesResponseClusterNodeInput struct {
	Address     string
	Version     string
	CreatedAtNs int64
	DataCenter  string
	Rack        string
}

// NewListClusterNodesResponseClusterNode builds a ZAP-encoded
// ListClusterNodesResponse.ClusterNode message from in and returns the bytes.
func NewListClusterNodesResponseClusterNode(in ListClusterNodesResponseClusterNodeInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(listClusterNodesResponseClusterNodeSize)
	ob.SetText(listClusterNodesResponseClusterNodeAddressOff, in.Address)
	ob.SetText(listClusterNodesResponseClusterNodeVersionOff, in.Version)
	ob.SetInt64(listClusterNodesResponseClusterNodeCreatedAtNsOff, in.CreatedAtNs)
	ob.SetText(listClusterNodesResponseClusterNodeDataCenterOff, in.DataCenter)
	ob.SetText(listClusterNodesResponseClusterNodeRackOff, in.Rack)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ListClusterNodesResponse ---

const (
	listClusterNodesResponseClusterNodesOff = 0
	listClusterNodesResponseSize            = 8
)

// ListClusterNodesResponse is a zero-copy view into a ZAP-encoded
// ListClusterNodesResponse message.
type ListClusterNodesResponse struct{ o zap.Object }

// WrapListClusterNodesResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapListClusterNodesResponse(b []byte) (ListClusterNodesResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ListClusterNodesResponse{}, err
	}
	return ListClusterNodesResponse{o: m.Root()}, nil
}

// ClusterNodesLen reports the number of cluster_nodes elements (proto field 1,
// repeated ListClusterNodesResponse.ClusterNode).
func (t ListClusterNodesResponse) ClusterNodesLen() int {
	return t.o.List(listClusterNodesResponseClusterNodesOff).Len()
}

// ClusterNodeAt returns the i-th cluster_nodes element.
func (t ListClusterNodesResponse) ClusterNodeAt(i int) (ListClusterNodesResponseClusterNode, bool) {
	o := t.o.List(listClusterNodesResponseClusterNodesOff).ObjectAt(i)
	if o.IsNull() {
		return ListClusterNodesResponseClusterNode{}, false
	}
	return ListClusterNodesResponseClusterNode{o: o}, true
}

// ListClusterNodesResponseInput collects the field values for
// NewListClusterNodesResponse. ClusterNodes takes each element as its own ZAP
// sub-buffer (from NewListClusterNodesResponseClusterNode).
type ListClusterNodesResponseInput struct {
	ClusterNodes [][]byte
}

// NewListClusterNodesResponse builds a ZAP-encoded ListClusterNodesResponse
// message from in and returns the bytes.
func NewListClusterNodesResponse(in ListClusterNodesResponseInput) []byte {
	b := zap.NewBuilder(256)
	clusterNodesOff, clusterNodesLen := setMsgList(b, in.ClusterNodes)
	ob := b.StartObject(listClusterNodesResponseSize)
	ob.SetList(listClusterNodesResponseClusterNodesOff, clusterNodesOff, clusterNodesLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- LeaseAdminTokenRequest ---

const (
	leaseAdminTokenRequestPreviousTokenOff    = 0
	leaseAdminTokenRequestPreviousLockTimeOff = 8
	leaseAdminTokenRequestLockNameOff         = 16
	leaseAdminTokenRequestClientNameOff       = 24
	leaseAdminTokenRequestMessageOff          = 32
	leaseAdminTokenRequestSize                = 40
)

// LeaseAdminTokenRequest is a zero-copy view into a ZAP-encoded
// LeaseAdminTokenRequest message.
type LeaseAdminTokenRequest struct{ o zap.Object }

// WrapLeaseAdminTokenRequest parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapLeaseAdminTokenRequest(b []byte) (LeaseAdminTokenRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return LeaseAdminTokenRequest{}, err
	}
	return LeaseAdminTokenRequest{o: m.Root()}, nil
}

// PreviousToken reads the previous_token field (proto field 1, int64).
func (t LeaseAdminTokenRequest) PreviousToken() int64 {
	return t.o.Int64(leaseAdminTokenRequestPreviousTokenOff)
}

// PreviousLockTime reads the previous_lock_time field (proto field 2, int64).
func (t LeaseAdminTokenRequest) PreviousLockTime() int64 {
	return t.o.Int64(leaseAdminTokenRequestPreviousLockTimeOff)
}

// LockName reads the lock_name field (proto field 3, string).
func (t LeaseAdminTokenRequest) LockName() string {
	return t.o.Text(leaseAdminTokenRequestLockNameOff)
}

// ClientName reads the client_name field (proto field 4, string).
func (t LeaseAdminTokenRequest) ClientName() string {
	return t.o.Text(leaseAdminTokenRequestClientNameOff)
}

// Message reads the message field (proto field 5, string).
func (t LeaseAdminTokenRequest) Message() string {
	return t.o.Text(leaseAdminTokenRequestMessageOff)
}

// LeaseAdminTokenRequestInput collects the field values for
// NewLeaseAdminTokenRequest.
type LeaseAdminTokenRequestInput struct {
	PreviousToken    int64
	PreviousLockTime int64
	LockName         string
	ClientName       string
	Message          string
}

// NewLeaseAdminTokenRequest builds a ZAP-encoded LeaseAdminTokenRequest message
// from in and returns the bytes.
func NewLeaseAdminTokenRequest(in LeaseAdminTokenRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(leaseAdminTokenRequestSize)
	ob.SetInt64(leaseAdminTokenRequestPreviousTokenOff, in.PreviousToken)
	ob.SetInt64(leaseAdminTokenRequestPreviousLockTimeOff, in.PreviousLockTime)
	ob.SetText(leaseAdminTokenRequestLockNameOff, in.LockName)
	ob.SetText(leaseAdminTokenRequestClientNameOff, in.ClientName)
	ob.SetText(leaseAdminTokenRequestMessageOff, in.Message)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- LeaseAdminTokenResponse ---

const (
	leaseAdminTokenResponseTokenOff    = 0
	leaseAdminTokenResponseLockTsNsOff = 8
	leaseAdminTokenResponseSize        = 16
)

// LeaseAdminTokenResponse is a zero-copy view into a ZAP-encoded
// LeaseAdminTokenResponse message.
type LeaseAdminTokenResponse struct{ o zap.Object }

// WrapLeaseAdminTokenResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapLeaseAdminTokenResponse(b []byte) (LeaseAdminTokenResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return LeaseAdminTokenResponse{}, err
	}
	return LeaseAdminTokenResponse{o: m.Root()}, nil
}

// Token reads the token field (proto field 1, int64).
func (t LeaseAdminTokenResponse) Token() int64 {
	return t.o.Int64(leaseAdminTokenResponseTokenOff)
}

// LockTsNs reads the lock_ts_ns field (proto field 2, int64).
func (t LeaseAdminTokenResponse) LockTsNs() int64 {
	return t.o.Int64(leaseAdminTokenResponseLockTsNsOff)
}

// LeaseAdminTokenResponseInput collects the field values for
// NewLeaseAdminTokenResponse.
type LeaseAdminTokenResponseInput struct {
	Token    int64
	LockTsNs int64
}

// NewLeaseAdminTokenResponse builds a ZAP-encoded LeaseAdminTokenResponse
// message from in and returns the bytes.
func NewLeaseAdminTokenResponse(in LeaseAdminTokenResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(leaseAdminTokenResponseSize)
	ob.SetInt64(leaseAdminTokenResponseTokenOff, in.Token)
	ob.SetInt64(leaseAdminTokenResponseLockTsNsOff, in.LockTsNs)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ReleaseAdminTokenRequest ---

const (
	releaseAdminTokenRequestPreviousTokenOff    = 0
	releaseAdminTokenRequestPreviousLockTimeOff = 8
	releaseAdminTokenRequestLockNameOff         = 16
	releaseAdminTokenRequestSize                = 24
)

// ReleaseAdminTokenRequest is a zero-copy view into a ZAP-encoded
// ReleaseAdminTokenRequest message.
type ReleaseAdminTokenRequest struct{ o zap.Object }

// WrapReleaseAdminTokenRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapReleaseAdminTokenRequest(b []byte) (ReleaseAdminTokenRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ReleaseAdminTokenRequest{}, err
	}
	return ReleaseAdminTokenRequest{o: m.Root()}, nil
}

// PreviousToken reads the previous_token field (proto field 1, int64).
func (t ReleaseAdminTokenRequest) PreviousToken() int64 {
	return t.o.Int64(releaseAdminTokenRequestPreviousTokenOff)
}

// PreviousLockTime reads the previous_lock_time field (proto field 2, int64).
func (t ReleaseAdminTokenRequest) PreviousLockTime() int64 {
	return t.o.Int64(releaseAdminTokenRequestPreviousLockTimeOff)
}

// LockName reads the lock_name field (proto field 3, string).
func (t ReleaseAdminTokenRequest) LockName() string {
	return t.o.Text(releaseAdminTokenRequestLockNameOff)
}

// ReleaseAdminTokenRequestInput collects the field values for
// NewReleaseAdminTokenRequest.
type ReleaseAdminTokenRequestInput struct {
	PreviousToken    int64
	PreviousLockTime int64
	LockName         string
}

// NewReleaseAdminTokenRequest builds a ZAP-encoded ReleaseAdminTokenRequest
// message from in and returns the bytes.
func NewReleaseAdminTokenRequest(in ReleaseAdminTokenRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(releaseAdminTokenRequestSize)
	ob.SetInt64(releaseAdminTokenRequestPreviousTokenOff, in.PreviousToken)
	ob.SetInt64(releaseAdminTokenRequestPreviousLockTimeOff, in.PreviousLockTime)
	ob.SetText(releaseAdminTokenRequestLockNameOff, in.LockName)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ReleaseAdminTokenResponse ---

const releaseAdminTokenResponseSize = 0

// ReleaseAdminTokenResponse is a zero-copy view into a ZAP-encoded
// ReleaseAdminTokenResponse message. It carries no fields.
type ReleaseAdminTokenResponse struct{ o zap.Object }

// WrapReleaseAdminTokenResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapReleaseAdminTokenResponse(b []byte) (ReleaseAdminTokenResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ReleaseAdminTokenResponse{}, err
	}
	return ReleaseAdminTokenResponse{o: m.Root()}, nil
}

// ReleaseAdminTokenResponseInput collects the field values for
// NewReleaseAdminTokenResponse. It has no fields.
type ReleaseAdminTokenResponseInput struct{}

// NewReleaseAdminTokenResponse builds a ZAP-encoded ReleaseAdminTokenResponse
// message and returns the bytes.
func NewReleaseAdminTokenResponse(in ReleaseAdminTokenResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(releaseAdminTokenResponseSize)
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

// --- RaftAddServerRequest ---

const (
	raftAddServerRequestIdOff      = 0
	raftAddServerRequestAddressOff = 8
	raftAddServerRequestVoterOff   = 16
	raftAddServerRequestSize       = 17
)

// RaftAddServerRequest is a zero-copy view into a ZAP-encoded
// RaftAddServerRequest message.
type RaftAddServerRequest struct{ o zap.Object }

// WrapRaftAddServerRequest parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapRaftAddServerRequest(b []byte) (RaftAddServerRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return RaftAddServerRequest{}, err
	}
	return RaftAddServerRequest{o: m.Root()}, nil
}

// Id reads the id field (proto field 1, string).
func (t RaftAddServerRequest) Id() string { return t.o.Text(raftAddServerRequestIdOff) }

// Address reads the address field (proto field 2, string).
func (t RaftAddServerRequest) Address() string { return t.o.Text(raftAddServerRequestAddressOff) }

// Voter reads the voter field (proto field 3, bool).
func (t RaftAddServerRequest) Voter() bool { return t.o.Bool(raftAddServerRequestVoterOff) }

// RaftAddServerRequestInput collects the field values for
// NewRaftAddServerRequest.
type RaftAddServerRequestInput struct {
	Id      string
	Address string
	Voter   bool
}

// NewRaftAddServerRequest builds a ZAP-encoded RaftAddServerRequest message from
// in and returns the bytes.
func NewRaftAddServerRequest(in RaftAddServerRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(raftAddServerRequestSize)
	ob.SetText(raftAddServerRequestIdOff, in.Id)
	ob.SetText(raftAddServerRequestAddressOff, in.Address)
	ob.SetBool(raftAddServerRequestVoterOff, in.Voter)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- RaftAddServerResponse ---

const raftAddServerResponseSize = 0

// RaftAddServerResponse is a zero-copy view into a ZAP-encoded
// RaftAddServerResponse message. It carries no fields.
type RaftAddServerResponse struct{ o zap.Object }

// WrapRaftAddServerResponse parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapRaftAddServerResponse(b []byte) (RaftAddServerResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return RaftAddServerResponse{}, err
	}
	return RaftAddServerResponse{o: m.Root()}, nil
}

// RaftAddServerResponseInput collects the field values for
// NewRaftAddServerResponse. It has no fields.
type RaftAddServerResponseInput struct{}

// NewRaftAddServerResponse builds a ZAP-encoded RaftAddServerResponse message
// and returns the bytes.
func NewRaftAddServerResponse(in RaftAddServerResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(raftAddServerResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- RaftRemoveServerRequest ---

const (
	raftRemoveServerRequestIdOff    = 0
	raftRemoveServerRequestForceOff = 8
	raftRemoveServerRequestSize     = 9
)

// RaftRemoveServerRequest is a zero-copy view into a ZAP-encoded
// RaftRemoveServerRequest message.
type RaftRemoveServerRequest struct{ o zap.Object }

// WrapRaftRemoveServerRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapRaftRemoveServerRequest(b []byte) (RaftRemoveServerRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return RaftRemoveServerRequest{}, err
	}
	return RaftRemoveServerRequest{o: m.Root()}, nil
}

// Id reads the id field (proto field 1, string).
func (t RaftRemoveServerRequest) Id() string { return t.o.Text(raftRemoveServerRequestIdOff) }

// Force reads the force field (proto field 2, bool).
func (t RaftRemoveServerRequest) Force() bool { return t.o.Bool(raftRemoveServerRequestForceOff) }

// RaftRemoveServerRequestInput collects the field values for
// NewRaftRemoveServerRequest.
type RaftRemoveServerRequestInput struct {
	Id    string
	Force bool
}

// NewRaftRemoveServerRequest builds a ZAP-encoded RaftRemoveServerRequest
// message from in and returns the bytes.
func NewRaftRemoveServerRequest(in RaftRemoveServerRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(raftRemoveServerRequestSize)
	ob.SetText(raftRemoveServerRequestIdOff, in.Id)
	ob.SetBool(raftRemoveServerRequestForceOff, in.Force)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- RaftRemoveServerResponse ---

const raftRemoveServerResponseSize = 0

// RaftRemoveServerResponse is a zero-copy view into a ZAP-encoded
// RaftRemoveServerResponse message. It carries no fields.
type RaftRemoveServerResponse struct{ o zap.Object }

// WrapRaftRemoveServerResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapRaftRemoveServerResponse(b []byte) (RaftRemoveServerResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return RaftRemoveServerResponse{}, err
	}
	return RaftRemoveServerResponse{o: m.Root()}, nil
}

// RaftRemoveServerResponseInput collects the field values for
// NewRaftRemoveServerResponse. It has no fields.
type RaftRemoveServerResponseInput struct{}

// NewRaftRemoveServerResponse builds a ZAP-encoded RaftRemoveServerResponse
// message and returns the bytes.
func NewRaftRemoveServerResponse(in RaftRemoveServerResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(raftRemoveServerResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- RaftListClusterServersRequest ---

const raftListClusterServersRequestSize = 0

// RaftListClusterServersRequest is a zero-copy view into a ZAP-encoded
// RaftListClusterServersRequest message. It carries no fields.
type RaftListClusterServersRequest struct{ o zap.Object }

// WrapRaftListClusterServersRequest parses b and returns a typed view. Returns
// an error if the wire-level checks (magic, version, size) fail.
func WrapRaftListClusterServersRequest(b []byte) (RaftListClusterServersRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return RaftListClusterServersRequest{}, err
	}
	return RaftListClusterServersRequest{o: m.Root()}, nil
}

// RaftListClusterServersRequestInput collects the field values for
// NewRaftListClusterServersRequest. It has no fields.
type RaftListClusterServersRequestInput struct{}

// NewRaftListClusterServersRequest builds a ZAP-encoded
// RaftListClusterServersRequest message and returns the bytes.
func NewRaftListClusterServersRequest(in RaftListClusterServersRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(raftListClusterServersRequestSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- RaftListClusterServersResponse.ClusterServers ---

const (
	raftListClusterServersResponseClusterServersIdOff       = 0
	raftListClusterServersResponseClusterServersAddressOff  = 8
	raftListClusterServersResponseClusterServersSuffrageOff = 16
	raftListClusterServersResponseClusterServersIsLeaderOff = 24
	raftListClusterServersResponseClusterServersSize        = 25
)

// RaftListClusterServersResponseClusterServers is a zero-copy view into a
// ZAP-encoded RaftListClusterServersResponse.ClusterServers message.
type RaftListClusterServersResponseClusterServers struct{ o zap.Object }

// WrapRaftListClusterServersResponseClusterServers parses b and returns a typed
// view. Returns an error if the wire-level checks (magic, version, size) fail.
func WrapRaftListClusterServersResponseClusterServers(b []byte) (RaftListClusterServersResponseClusterServers, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return RaftListClusterServersResponseClusterServers{}, err
	}
	return RaftListClusterServersResponseClusterServers{o: m.Root()}, nil
}

// Id reads the id field (proto field 1, string).
func (t RaftListClusterServersResponseClusterServers) Id() string {
	return t.o.Text(raftListClusterServersResponseClusterServersIdOff)
}

// Address reads the address field (proto field 2, string).
func (t RaftListClusterServersResponseClusterServers) Address() string {
	return t.o.Text(raftListClusterServersResponseClusterServersAddressOff)
}

// Suffrage reads the suffrage field (proto field 3, string).
func (t RaftListClusterServersResponseClusterServers) Suffrage() string {
	return t.o.Text(raftListClusterServersResponseClusterServersSuffrageOff)
}

// IsLeader reads the isLeader field (proto field 4, bool).
func (t RaftListClusterServersResponseClusterServers) IsLeader() bool {
	return t.o.Bool(raftListClusterServersResponseClusterServersIsLeaderOff)
}

// RaftListClusterServersResponseClusterServersInput collects the field values
// for NewRaftListClusterServersResponseClusterServers.
type RaftListClusterServersResponseClusterServersInput struct {
	Id       string
	Address  string
	Suffrage string
	IsLeader bool
}

// NewRaftListClusterServersResponseClusterServers builds a ZAP-encoded
// RaftListClusterServersResponse.ClusterServers message from in and returns the
// bytes.
func NewRaftListClusterServersResponseClusterServers(in RaftListClusterServersResponseClusterServersInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(raftListClusterServersResponseClusterServersSize)
	ob.SetText(raftListClusterServersResponseClusterServersIdOff, in.Id)
	ob.SetText(raftListClusterServersResponseClusterServersAddressOff, in.Address)
	ob.SetText(raftListClusterServersResponseClusterServersSuffrageOff, in.Suffrage)
	ob.SetBool(raftListClusterServersResponseClusterServersIsLeaderOff, in.IsLeader)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- RaftListClusterServersResponse ---

const (
	raftListClusterServersResponseClusterServersListOff = 0
	raftListClusterServersResponseSize                  = 8
)

// RaftListClusterServersResponse is a zero-copy view into a ZAP-encoded
// RaftListClusterServersResponse message.
type RaftListClusterServersResponse struct{ o zap.Object }

// WrapRaftListClusterServersResponse parses b and returns a typed view. Returns
// an error if the wire-level checks (magic, version, size) fail.
func WrapRaftListClusterServersResponse(b []byte) (RaftListClusterServersResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return RaftListClusterServersResponse{}, err
	}
	return RaftListClusterServersResponse{o: m.Root()}, nil
}

// ClusterServersLen reports the number of cluster_servers elements (proto field
// 1, repeated RaftListClusterServersResponse.ClusterServers).
func (t RaftListClusterServersResponse) ClusterServersLen() int {
	return t.o.List(raftListClusterServersResponseClusterServersListOff).Len()
}

// ClusterServerAt returns the i-th cluster_servers element.
func (t RaftListClusterServersResponse) ClusterServerAt(i int) (RaftListClusterServersResponseClusterServers, bool) {
	o := t.o.List(raftListClusterServersResponseClusterServersListOff).ObjectAt(i)
	if o.IsNull() {
		return RaftListClusterServersResponseClusterServers{}, false
	}
	return RaftListClusterServersResponseClusterServers{o: o}, true
}

// RaftListClusterServersResponseInput collects the field values for
// NewRaftListClusterServersResponse. ClusterServers takes each element as its
// own ZAP sub-buffer (from NewRaftListClusterServersResponseClusterServers).
type RaftListClusterServersResponseInput struct {
	ClusterServers [][]byte
}

// NewRaftListClusterServersResponse builds a ZAP-encoded
// RaftListClusterServersResponse message from in and returns the bytes.
func NewRaftListClusterServersResponse(in RaftListClusterServersResponseInput) []byte {
	b := zap.NewBuilder(256)
	clusterServersOff, clusterServersLen := setMsgList(b, in.ClusterServers)
	ob := b.StartObject(raftListClusterServersResponseSize)
	ob.SetList(raftListClusterServersResponseClusterServersListOff, clusterServersOff, clusterServersLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- RaftLeadershipTransferRequest ---

const (
	raftLeadershipTransferRequestTargetIdOff      = 0
	raftLeadershipTransferRequestTargetAddressOff = 8
	raftLeadershipTransferRequestSize             = 16
)

// RaftLeadershipTransferRequest is a zero-copy view into a ZAP-encoded
// RaftLeadershipTransferRequest message.
type RaftLeadershipTransferRequest struct{ o zap.Object }

// WrapRaftLeadershipTransferRequest parses b and returns a typed view. Returns
// an error if the wire-level checks (magic, version, size) fail.
func WrapRaftLeadershipTransferRequest(b []byte) (RaftLeadershipTransferRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return RaftLeadershipTransferRequest{}, err
	}
	return RaftLeadershipTransferRequest{o: m.Root()}, nil
}

// TargetId reads the target_id field (proto field 1, string).
func (t RaftLeadershipTransferRequest) TargetId() string {
	return t.o.Text(raftLeadershipTransferRequestTargetIdOff)
}

// TargetAddress reads the target_address field (proto field 2, string).
func (t RaftLeadershipTransferRequest) TargetAddress() string {
	return t.o.Text(raftLeadershipTransferRequestTargetAddressOff)
}

// RaftLeadershipTransferRequestInput collects the field values for
// NewRaftLeadershipTransferRequest.
type RaftLeadershipTransferRequestInput struct {
	TargetId      string
	TargetAddress string
}

// NewRaftLeadershipTransferRequest builds a ZAP-encoded
// RaftLeadershipTransferRequest message from in and returns the bytes.
func NewRaftLeadershipTransferRequest(in RaftLeadershipTransferRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(raftLeadershipTransferRequestSize)
	ob.SetText(raftLeadershipTransferRequestTargetIdOff, in.TargetId)
	ob.SetText(raftLeadershipTransferRequestTargetAddressOff, in.TargetAddress)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- RaftLeadershipTransferResponse ---

const (
	raftLeadershipTransferResponsePreviousLeaderOff = 0
	raftLeadershipTransferResponseNewLeaderOff      = 8
	raftLeadershipTransferResponseSize              = 16
)

// RaftLeadershipTransferResponse is a zero-copy view into a ZAP-encoded
// RaftLeadershipTransferResponse message.
type RaftLeadershipTransferResponse struct{ o zap.Object }

// WrapRaftLeadershipTransferResponse parses b and returns a typed view. Returns
// an error if the wire-level checks (magic, version, size) fail.
func WrapRaftLeadershipTransferResponse(b []byte) (RaftLeadershipTransferResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return RaftLeadershipTransferResponse{}, err
	}
	return RaftLeadershipTransferResponse{o: m.Root()}, nil
}

// PreviousLeader reads the previous_leader field (proto field 1, string).
func (t RaftLeadershipTransferResponse) PreviousLeader() string {
	return t.o.Text(raftLeadershipTransferResponsePreviousLeaderOff)
}

// NewLeader reads the new_leader field (proto field 2, string).
func (t RaftLeadershipTransferResponse) NewLeader() string {
	return t.o.Text(raftLeadershipTransferResponseNewLeaderOff)
}

// RaftLeadershipTransferResponseInput collects the field values for
// NewRaftLeadershipTransferResponse.
type RaftLeadershipTransferResponseInput struct {
	PreviousLeader string
	NewLeader      string
}

// NewRaftLeadershipTransferResponse builds a ZAP-encoded
// RaftLeadershipTransferResponse message from in and returns the bytes.
func NewRaftLeadershipTransferResponse(in RaftLeadershipTransferResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(raftLeadershipTransferResponseSize)
	ob.SetText(raftLeadershipTransferResponsePreviousLeaderOff, in.PreviousLeader)
	ob.SetText(raftLeadershipTransferResponseNewLeaderOff, in.NewLeader)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeGrowResponse ---

const volumeGrowResponseSize = 0

// VolumeGrowResponse is a zero-copy view into a ZAP-encoded VolumeGrowResponse
// message. It carries no fields.
type VolumeGrowResponse struct{ o zap.Object }

// WrapVolumeGrowResponse parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapVolumeGrowResponse(b []byte) (VolumeGrowResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeGrowResponse{}, err
	}
	return VolumeGrowResponse{o: m.Root()}, nil
}

// VolumeGrowResponseInput collects the field values for NewVolumeGrowResponse.
// It has no fields.
type VolumeGrowResponseInput struct{}

// NewVolumeGrowResponse builds a ZAP-encoded VolumeGrowResponse message and
// returns the bytes.
func NewVolumeGrowResponse(in VolumeGrowResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeGrowResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}
