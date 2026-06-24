// Code generated from master.proto; DO NOT EDIT.

package masterwire

import (
	zap "github.com/zap-proto/go"
)

// --- KeepConnectedRequest ---

const (
	keepConnectedRequestClientTypeOff    = 0
	keepConnectedRequestClientAddressOff = 8
	keepConnectedRequestVersionOff       = 16
	keepConnectedRequestFilerGroupOff    = 24
	keepConnectedRequestDataCenterOff    = 32
	keepConnectedRequestRackOff          = 40
	keepConnectedRequestSize             = 48
)

// KeepConnectedRequest is a zero-copy view into a ZAP-encoded
// KeepConnectedRequest message.
type KeepConnectedRequest struct{ o zap.Object }

// WrapKeepConnectedRequest parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapKeepConnectedRequest(b []byte) (KeepConnectedRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return KeepConnectedRequest{}, err
	}
	return KeepConnectedRequest{o: m.Root()}, nil
}

// ClientType reads the client_type field (proto field 1, string).
func (t KeepConnectedRequest) ClientType() string {
	return t.o.Text(keepConnectedRequestClientTypeOff)
}

// ClientAddress reads the client_address field (proto field 3, string).
func (t KeepConnectedRequest) ClientAddress() string {
	return t.o.Text(keepConnectedRequestClientAddressOff)
}

// Version reads the version field (proto field 4, string).
func (t KeepConnectedRequest) Version() string { return t.o.Text(keepConnectedRequestVersionOff) }

// FilerGroup reads the filer_group field (proto field 5, string).
func (t KeepConnectedRequest) FilerGroup() string {
	return t.o.Text(keepConnectedRequestFilerGroupOff)
}

// DataCenter reads the data_center field (proto field 6, string).
func (t KeepConnectedRequest) DataCenter() string {
	return t.o.Text(keepConnectedRequestDataCenterOff)
}

// Rack reads the rack field (proto field 7, string).
func (t KeepConnectedRequest) Rack() string { return t.o.Text(keepConnectedRequestRackOff) }

// KeepConnectedRequestInput collects the field values for
// NewKeepConnectedRequest.
type KeepConnectedRequestInput struct {
	ClientType    string
	ClientAddress string
	Version       string
	FilerGroup    string
	DataCenter    string
	Rack          string
}

// NewKeepConnectedRequest builds a ZAP-encoded KeepConnectedRequest message from
// in and returns the bytes.
func NewKeepConnectedRequest(in KeepConnectedRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(keepConnectedRequestSize)
	ob.SetText(keepConnectedRequestClientTypeOff, in.ClientType)
	ob.SetText(keepConnectedRequestClientAddressOff, in.ClientAddress)
	ob.SetText(keepConnectedRequestVersionOff, in.Version)
	ob.SetText(keepConnectedRequestFilerGroupOff, in.FilerGroup)
	ob.SetText(keepConnectedRequestDataCenterOff, in.DataCenter)
	ob.SetText(keepConnectedRequestRackOff, in.Rack)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeLocation ---

const (
	volumeLocationUrlOff           = 0
	volumeLocationPublicUrlOff     = 8
	volumeLocationNewVidsOff       = 16
	volumeLocationDeletedVidsOff   = 24
	volumeLocationLeaderOff        = 32
	volumeLocationDataCenterOff    = 40
	volumeLocationGrpcPortOff      = 48
	volumeLocationNewEcVidsOff     = 52
	volumeLocationDeletedEcVidsOff = 60
	volumeLocationSize             = 68
)

// VolumeLocation is a zero-copy view into a ZAP-encoded VolumeLocation message.
type VolumeLocation struct{ o zap.Object }

// WrapVolumeLocation parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapVolumeLocation(b []byte) (VolumeLocation, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeLocation{}, err
	}
	return VolumeLocation{o: m.Root()}, nil
}

// Url reads the url field (proto field 1, string).
func (t VolumeLocation) Url() string { return t.o.Text(volumeLocationUrlOff) }

// PublicUrl reads the public_url field (proto field 2, string).
func (t VolumeLocation) PublicUrl() string { return t.o.Text(volumeLocationPublicUrlOff) }

// NewVidsLen reports the number of new_vids elements (proto field 3, repeated
// uint32).
func (t VolumeLocation) NewVidsLen() int { return t.o.List(volumeLocationNewVidsOff).Len() }

// NewVidAt returns the i-th new_vids element.
func (t VolumeLocation) NewVidAt(i int) uint32 { return t.o.List(volumeLocationNewVidsOff).Uint32(i) }

// DeletedVidsLen reports the number of deleted_vids elements (proto field 4,
// repeated uint32).
func (t VolumeLocation) DeletedVidsLen() int { return t.o.List(volumeLocationDeletedVidsOff).Len() }

// DeletedVidAt returns the i-th deleted_vids element.
func (t VolumeLocation) DeletedVidAt(i int) uint32 {
	return t.o.List(volumeLocationDeletedVidsOff).Uint32(i)
}

// Leader reads the leader field (proto field 5, string).
func (t VolumeLocation) Leader() string { return t.o.Text(volumeLocationLeaderOff) }

// DataCenter reads the data_center field (proto field 6, string).
func (t VolumeLocation) DataCenter() string { return t.o.Text(volumeLocationDataCenterOff) }

// GrpcPort reads the grpc_port field (proto field 7, uint32).
func (t VolumeLocation) GrpcPort() uint32 { return t.o.Uint32(volumeLocationGrpcPortOff) }

// NewEcVidsLen reports the number of new_ec_vids elements (proto field 8,
// repeated uint32).
func (t VolumeLocation) NewEcVidsLen() int { return t.o.List(volumeLocationNewEcVidsOff).Len() }

// NewEcVidAt returns the i-th new_ec_vids element.
func (t VolumeLocation) NewEcVidAt(i int) uint32 {
	return t.o.List(volumeLocationNewEcVidsOff).Uint32(i)
}

// DeletedEcVidsLen reports the number of deleted_ec_vids elements (proto field
// 9, repeated uint32).
func (t VolumeLocation) DeletedEcVidsLen() int {
	return t.o.List(volumeLocationDeletedEcVidsOff).Len()
}

// DeletedEcVidAt returns the i-th deleted_ec_vids element.
func (t VolumeLocation) DeletedEcVidAt(i int) uint32 {
	return t.o.List(volumeLocationDeletedEcVidsOff).Uint32(i)
}

// VolumeLocationInput collects the field values for NewVolumeLocation.
type VolumeLocationInput struct {
	Url           string
	PublicUrl     string
	NewVids       []uint32
	DeletedVids   []uint32
	Leader        string
	DataCenter    string
	GrpcPort      uint32
	NewEcVids     []uint32
	DeletedEcVids []uint32
}

// NewVolumeLocation builds a ZAP-encoded VolumeLocation message from in and
// returns the bytes.
func NewVolumeLocation(in VolumeLocationInput) []byte {
	b := zap.NewBuilder(256)
	newVidsOff, newVidsLen := setUint32List(b, in.NewVids)
	deletedVidsOff, deletedVidsLen := setUint32List(b, in.DeletedVids)
	newEcVidsOff, newEcVidsLen := setUint32List(b, in.NewEcVids)
	deletedEcVidsOff, deletedEcVidsLen := setUint32List(b, in.DeletedEcVids)
	ob := b.StartObject(volumeLocationSize)
	ob.SetText(volumeLocationUrlOff, in.Url)
	ob.SetText(volumeLocationPublicUrlOff, in.PublicUrl)
	ob.SetList(volumeLocationNewVidsOff, newVidsOff, newVidsLen)
	ob.SetList(volumeLocationDeletedVidsOff, deletedVidsOff, deletedVidsLen)
	ob.SetText(volumeLocationLeaderOff, in.Leader)
	ob.SetText(volumeLocationDataCenterOff, in.DataCenter)
	ob.SetUint32(volumeLocationGrpcPortOff, in.GrpcPort)
	ob.SetList(volumeLocationNewEcVidsOff, newEcVidsOff, newEcVidsLen)
	ob.SetList(volumeLocationDeletedEcVidsOff, deletedEcVidsOff, deletedEcVidsLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ClusterNodeUpdate ---

const (
	clusterNodeUpdateNodeTypeOff    = 0
	clusterNodeUpdateAddressOff     = 8
	clusterNodeUpdateIsAddOff       = 16
	clusterNodeUpdateFilerGroupOff  = 17
	clusterNodeUpdateCreatedAtNsOff = 25
	clusterNodeUpdateSize           = 33
)

// ClusterNodeUpdate is a zero-copy view into a ZAP-encoded ClusterNodeUpdate
// message.
type ClusterNodeUpdate struct{ o zap.Object }

// WrapClusterNodeUpdate parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapClusterNodeUpdate(b []byte) (ClusterNodeUpdate, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ClusterNodeUpdate{}, err
	}
	return ClusterNodeUpdate{o: m.Root()}, nil
}

// NodeType reads the node_type field (proto field 1, string).
func (t ClusterNodeUpdate) NodeType() string { return t.o.Text(clusterNodeUpdateNodeTypeOff) }

// Address reads the address field (proto field 2, string).
func (t ClusterNodeUpdate) Address() string { return t.o.Text(clusterNodeUpdateAddressOff) }

// IsAdd reads the is_add field (proto field 4, bool).
func (t ClusterNodeUpdate) IsAdd() bool { return t.o.Bool(clusterNodeUpdateIsAddOff) }

// FilerGroup reads the filer_group field (proto field 5, string).
func (t ClusterNodeUpdate) FilerGroup() string { return t.o.Text(clusterNodeUpdateFilerGroupOff) }

// CreatedAtNs reads the created_at_ns field (proto field 6, int64).
func (t ClusterNodeUpdate) CreatedAtNs() int64 { return t.o.Int64(clusterNodeUpdateCreatedAtNsOff) }

// ClusterNodeUpdateInput collects the field values for NewClusterNodeUpdate.
type ClusterNodeUpdateInput struct {
	NodeType    string
	Address     string
	IsAdd       bool
	FilerGroup  string
	CreatedAtNs int64
}

// NewClusterNodeUpdate builds a ZAP-encoded ClusterNodeUpdate message from in
// and returns the bytes.
func NewClusterNodeUpdate(in ClusterNodeUpdateInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(clusterNodeUpdateSize)
	ob.SetText(clusterNodeUpdateNodeTypeOff, in.NodeType)
	ob.SetText(clusterNodeUpdateAddressOff, in.Address)
	ob.SetBool(clusterNodeUpdateIsAddOff, in.IsAdd)
	ob.SetText(clusterNodeUpdateFilerGroupOff, in.FilerGroup)
	ob.SetInt64(clusterNodeUpdateCreatedAtNsOff, in.CreatedAtNs)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- LockRingUpdate ---

const (
	lockRingUpdateFilerGroupOff = 0
	lockRingUpdateServersOff    = 8
	lockRingUpdateVersionOff    = 16
	lockRingUpdateSize          = 24
)

// LockRingUpdate is a zero-copy view into a ZAP-encoded LockRingUpdate message.
type LockRingUpdate struct{ o zap.Object }

// WrapLockRingUpdate parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapLockRingUpdate(b []byte) (LockRingUpdate, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return LockRingUpdate{}, err
	}
	return LockRingUpdate{o: m.Root()}, nil
}

// FilerGroup reads the filer_group field (proto field 1, string).
func (t LockRingUpdate) FilerGroup() string { return t.o.Text(lockRingUpdateFilerGroupOff) }

// ServersLen reports the number of servers elements (proto field 2, repeated
// string).
func (t LockRingUpdate) ServersLen() int { return t.o.List(lockRingUpdateServersOff).Len() }

// ServerAt returns the i-th servers element.
func (t LockRingUpdate) ServerAt(i int) string {
	return stringAt(t.o.List(lockRingUpdateServersOff), i)
}

// Version reads the version field (proto field 3, int64).
func (t LockRingUpdate) Version() int64 { return t.o.Int64(lockRingUpdateVersionOff) }

// LockRingUpdateInput collects the field values for NewLockRingUpdate.
type LockRingUpdateInput struct {
	FilerGroup string
	Servers    []string
	Version    int64
}

// NewLockRingUpdate builds a ZAP-encoded LockRingUpdate message from in and
// returns the bytes.
func NewLockRingUpdate(in LockRingUpdateInput) []byte {
	b := zap.NewBuilder(256)
	serversOff, serversLen := setStringList(b, in.Servers)
	ob := b.StartObject(lockRingUpdateSize)
	ob.SetText(lockRingUpdateFilerGroupOff, in.FilerGroup)
	ob.SetList(lockRingUpdateServersOff, serversOff, serversLen)
	ob.SetInt64(lockRingUpdateVersionOff, in.Version)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- KeepConnectedResponse ---

const (
	keepConnectedResponseVolumeLocationOff    = 0
	keepConnectedResponseClusterNodeUpdateOff = 8
	keepConnectedResponseLockRingUpdateOff    = 16
	keepConnectedResponseSize                 = 24
)

// KeepConnectedResponse is a zero-copy view into a ZAP-encoded
// KeepConnectedResponse message.
type KeepConnectedResponse struct{ o zap.Object }

// WrapKeepConnectedResponse parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapKeepConnectedResponse(b []byte) (KeepConnectedResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return KeepConnectedResponse{}, err
	}
	return KeepConnectedResponse{o: m.Root()}, nil
}

// VolumeLocation reads the volume_location field (proto field 1, message
// VolumeLocation). The bool is false when the field is absent.
func (t KeepConnectedResponse) VolumeLocation() (VolumeLocation, bool) {
	b := t.o.Bytes(keepConnectedResponseVolumeLocationOff)
	if len(b) == 0 {
		return VolumeLocation{}, false
	}
	v, err := WrapVolumeLocation(b)
	if err != nil {
		return VolumeLocation{}, false
	}
	return v, true
}

// ClusterNodeUpdate reads the cluster_node_update field (proto field 2, message
// ClusterNodeUpdate). The bool is false when the field is absent.
func (t KeepConnectedResponse) ClusterNodeUpdate() (ClusterNodeUpdate, bool) {
	b := t.o.Bytes(keepConnectedResponseClusterNodeUpdateOff)
	if len(b) == 0 {
		return ClusterNodeUpdate{}, false
	}
	c, err := WrapClusterNodeUpdate(b)
	if err != nil {
		return ClusterNodeUpdate{}, false
	}
	return c, true
}

// LockRingUpdate reads the lock_ring_update field (proto field 3, message
// LockRingUpdate). The bool is false when the field is absent.
func (t KeepConnectedResponse) LockRingUpdate() (LockRingUpdate, bool) {
	b := t.o.Bytes(keepConnectedResponseLockRingUpdateOff)
	if len(b) == 0 {
		return LockRingUpdate{}, false
	}
	l, err := WrapLockRingUpdate(b)
	if err != nil {
		return LockRingUpdate{}, false
	}
	return l, true
}

// KeepConnectedResponseInput collects the field values for
// NewKeepConnectedResponse. Each message field is the child's own ZAP buffer;
// nil encodes an absent field.
type KeepConnectedResponseInput struct {
	VolumeLocation    []byte
	ClusterNodeUpdate []byte
	LockRingUpdate    []byte
}

// NewKeepConnectedResponse builds a ZAP-encoded KeepConnectedResponse message
// from in and returns the bytes.
func NewKeepConnectedResponse(in KeepConnectedResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(keepConnectedResponseSize)
	ob.SetBytes(keepConnectedResponseVolumeLocationOff, in.VolumeLocation)
	ob.SetBytes(keepConnectedResponseClusterNodeUpdateOff, in.ClusterNodeUpdate)
	ob.SetBytes(keepConnectedResponseLockRingUpdateOff, in.LockRingUpdate)
	ob.FinishAsRoot()
	return b.Finish()
}
