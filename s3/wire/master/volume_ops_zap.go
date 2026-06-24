// Code generated from master.proto; DO NOT EDIT.

package masterwire

import (
	zap "github.com/zap-proto/go"
)

// --- LookupEcVolumeRequest ---

const (
	lookupEcVolumeRequestVolumeIdOff = 0
	lookupEcVolumeRequestSize        = 4
)

// LookupEcVolumeRequest is a zero-copy view into a ZAP-encoded
// LookupEcVolumeRequest message.
type LookupEcVolumeRequest struct{ o zap.Object }

// WrapLookupEcVolumeRequest parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapLookupEcVolumeRequest(b []byte) (LookupEcVolumeRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return LookupEcVolumeRequest{}, err
	}
	return LookupEcVolumeRequest{o: m.Root()}, nil
}

// VolumeId reads the volume_id field (proto field 1, uint32).
func (t LookupEcVolumeRequest) VolumeId() uint32 {
	return t.o.Uint32(lookupEcVolumeRequestVolumeIdOff)
}

// LookupEcVolumeRequestInput collects the field values for
// NewLookupEcVolumeRequest.
type LookupEcVolumeRequestInput struct {
	VolumeId uint32
}

// NewLookupEcVolumeRequest builds a ZAP-encoded LookupEcVolumeRequest message
// from in and returns the bytes.
func NewLookupEcVolumeRequest(in LookupEcVolumeRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(lookupEcVolumeRequestSize)
	ob.SetUint32(lookupEcVolumeRequestVolumeIdOff, in.VolumeId)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- LookupEcVolumeResponse.EcShardIdLocation ---

const (
	lookupEcVolumeResponseEcShardIdLocationShardIdOff   = 0
	lookupEcVolumeResponseEcShardIdLocationLocationsOff = 4
	lookupEcVolumeResponseEcShardIdLocationSize         = 12
)

// LookupEcVolumeResponseEcShardIdLocation is a zero-copy view into a ZAP-encoded
// LookupEcVolumeResponse.EcShardIdLocation message.
type LookupEcVolumeResponseEcShardIdLocation struct{ o zap.Object }

// WrapLookupEcVolumeResponseEcShardIdLocation parses b and returns a typed view.
// Returns an error if the wire-level checks (magic, version, size) fail.
func WrapLookupEcVolumeResponseEcShardIdLocation(b []byte) (LookupEcVolumeResponseEcShardIdLocation, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return LookupEcVolumeResponseEcShardIdLocation{}, err
	}
	return LookupEcVolumeResponseEcShardIdLocation{o: m.Root()}, nil
}

// ShardId reads the shard_id field (proto field 1, uint32).
func (t LookupEcVolumeResponseEcShardIdLocation) ShardId() uint32 {
	return t.o.Uint32(lookupEcVolumeResponseEcShardIdLocationShardIdOff)
}

// LocationsLen reports the number of locations elements (proto field 2, repeated
// Location).
func (t LookupEcVolumeResponseEcShardIdLocation) LocationsLen() int {
	return t.o.List(lookupEcVolumeResponseEcShardIdLocationLocationsOff).Len()
}

// LocationAt returns the i-th locations element.
func (t LookupEcVolumeResponseEcShardIdLocation) LocationAt(i int) (Location, bool) {
	o := t.o.List(lookupEcVolumeResponseEcShardIdLocationLocationsOff).ObjectAt(i)
	if o.IsNull() {
		return Location{}, false
	}
	return Location{o: o}, true
}

// LookupEcVolumeResponseEcShardIdLocationInput collects the field values for
// NewLookupEcVolumeResponseEcShardIdLocation. Locations takes each Location
// element as its own ZAP sub-buffer.
type LookupEcVolumeResponseEcShardIdLocationInput struct {
	ShardId   uint32
	Locations [][]byte
}

// NewLookupEcVolumeResponseEcShardIdLocation builds a ZAP-encoded
// LookupEcVolumeResponse.EcShardIdLocation message from in and returns the
// bytes.
func NewLookupEcVolumeResponseEcShardIdLocation(in LookupEcVolumeResponseEcShardIdLocationInput) []byte {
	b := zap.NewBuilder(256)
	locationsOff, locationsLen := setMsgList(b, in.Locations)
	ob := b.StartObject(lookupEcVolumeResponseEcShardIdLocationSize)
	ob.SetUint32(lookupEcVolumeResponseEcShardIdLocationShardIdOff, in.ShardId)
	ob.SetList(lookupEcVolumeResponseEcShardIdLocationLocationsOff, locationsOff, locationsLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- LookupEcVolumeResponse ---

const (
	lookupEcVolumeResponseVolumeIdOff         = 0
	lookupEcVolumeResponseShardIdLocationsOff = 4
	lookupEcVolumeResponseSize                = 12
)

// LookupEcVolumeResponse is a zero-copy view into a ZAP-encoded
// LookupEcVolumeResponse message.
type LookupEcVolumeResponse struct{ o zap.Object }

// WrapLookupEcVolumeResponse parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapLookupEcVolumeResponse(b []byte) (LookupEcVolumeResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return LookupEcVolumeResponse{}, err
	}
	return LookupEcVolumeResponse{o: m.Root()}, nil
}

// VolumeId reads the volume_id field (proto field 1, uint32).
func (t LookupEcVolumeResponse) VolumeId() uint32 {
	return t.o.Uint32(lookupEcVolumeResponseVolumeIdOff)
}

// ShardIdLocationsLen reports the number of shard_id_locations elements (proto
// field 2, repeated LookupEcVolumeResponse.EcShardIdLocation).
func (t LookupEcVolumeResponse) ShardIdLocationsLen() int {
	return t.o.List(lookupEcVolumeResponseShardIdLocationsOff).Len()
}

// ShardIdLocationAt returns the i-th shard_id_locations element.
func (t LookupEcVolumeResponse) ShardIdLocationAt(i int) (LookupEcVolumeResponseEcShardIdLocation, bool) {
	o := t.o.List(lookupEcVolumeResponseShardIdLocationsOff).ObjectAt(i)
	if o.IsNull() {
		return LookupEcVolumeResponseEcShardIdLocation{}, false
	}
	return LookupEcVolumeResponseEcShardIdLocation{o: o}, true
}

// LookupEcVolumeResponseInput collects the field values for
// NewLookupEcVolumeResponse. ShardIdLocations takes each element as its own ZAP
// sub-buffer (from NewLookupEcVolumeResponseEcShardIdLocation).
type LookupEcVolumeResponseInput struct {
	VolumeId         uint32
	ShardIdLocations [][]byte
}

// NewLookupEcVolumeResponse builds a ZAP-encoded LookupEcVolumeResponse message
// from in and returns the bytes.
func NewLookupEcVolumeResponse(in LookupEcVolumeResponseInput) []byte {
	b := zap.NewBuilder(256)
	shardIdLocationsOff, shardIdLocationsLen := setMsgList(b, in.ShardIdLocations)
	ob := b.StartObject(lookupEcVolumeResponseSize)
	ob.SetUint32(lookupEcVolumeResponseVolumeIdOff, in.VolumeId)
	ob.SetList(lookupEcVolumeResponseShardIdLocationsOff, shardIdLocationsOff, shardIdLocationsLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VacuumVolumeRequest ---

const (
	vacuumVolumeRequestGarbageThresholdOff = 0
	vacuumVolumeRequestVolumeIdOff         = 4
	vacuumVolumeRequestCollectionOff       = 8
	vacuumVolumeRequestSize                = 16
)

// VacuumVolumeRequest is a zero-copy view into a ZAP-encoded VacuumVolumeRequest
// message.
type VacuumVolumeRequest struct{ o zap.Object }

// WrapVacuumVolumeRequest parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapVacuumVolumeRequest(b []byte) (VacuumVolumeRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VacuumVolumeRequest{}, err
	}
	return VacuumVolumeRequest{o: m.Root()}, nil
}

// GarbageThreshold reads the garbage_threshold field (proto field 1, float).
func (t VacuumVolumeRequest) GarbageThreshold() float32 {
	return t.o.Float32(vacuumVolumeRequestGarbageThresholdOff)
}

// VolumeId reads the volume_id field (proto field 2, uint32).
func (t VacuumVolumeRequest) VolumeId() uint32 { return t.o.Uint32(vacuumVolumeRequestVolumeIdOff) }

// Collection reads the collection field (proto field 3, string).
func (t VacuumVolumeRequest) Collection() string {
	return t.o.Text(vacuumVolumeRequestCollectionOff)
}

// VacuumVolumeRequestInput collects the field values for NewVacuumVolumeRequest.
type VacuumVolumeRequestInput struct {
	GarbageThreshold float32
	VolumeId         uint32
	Collection       string
}

// NewVacuumVolumeRequest builds a ZAP-encoded VacuumVolumeRequest message from
// in and returns the bytes.
func NewVacuumVolumeRequest(in VacuumVolumeRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(vacuumVolumeRequestSize)
	ob.SetFloat32(vacuumVolumeRequestGarbageThresholdOff, in.GarbageThreshold)
	ob.SetUint32(vacuumVolumeRequestVolumeIdOff, in.VolumeId)
	ob.SetText(vacuumVolumeRequestCollectionOff, in.Collection)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VacuumVolumeResponse ---

const vacuumVolumeResponseSize = 0

// VacuumVolumeResponse is a zero-copy view into a ZAP-encoded
// VacuumVolumeResponse message. VacuumVolumeResponse carries no fields.
type VacuumVolumeResponse struct{ o zap.Object }

// WrapVacuumVolumeResponse parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapVacuumVolumeResponse(b []byte) (VacuumVolumeResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VacuumVolumeResponse{}, err
	}
	return VacuumVolumeResponse{o: m.Root()}, nil
}

// VacuumVolumeResponseInput collects the field values for
// NewVacuumVolumeResponse. VacuumVolumeResponse has no fields.
type VacuumVolumeResponseInput struct{}

// NewVacuumVolumeResponse builds a ZAP-encoded VacuumVolumeResponse message and
// returns the bytes.
func NewVacuumVolumeResponse(in VacuumVolumeResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(vacuumVolumeResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- DisableVacuumRequest ---

const (
	disableVacuumRequestByPluginOff = 0
	disableVacuumRequestSize        = 1
)

// DisableVacuumRequest is a zero-copy view into a ZAP-encoded
// DisableVacuumRequest message.
type DisableVacuumRequest struct{ o zap.Object }

// WrapDisableVacuumRequest parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapDisableVacuumRequest(b []byte) (DisableVacuumRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return DisableVacuumRequest{}, err
	}
	return DisableVacuumRequest{o: m.Root()}, nil
}

// ByPlugin reads the by_plugin field (proto field 1, bool).
func (t DisableVacuumRequest) ByPlugin() bool { return t.o.Bool(disableVacuumRequestByPluginOff) }

// DisableVacuumRequestInput collects the field values for
// NewDisableVacuumRequest.
type DisableVacuumRequestInput struct {
	ByPlugin bool
}

// NewDisableVacuumRequest builds a ZAP-encoded DisableVacuumRequest message from
// in and returns the bytes.
func NewDisableVacuumRequest(in DisableVacuumRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(disableVacuumRequestSize)
	ob.SetBool(disableVacuumRequestByPluginOff, in.ByPlugin)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- DisableVacuumResponse ---

const disableVacuumResponseSize = 0

// DisableVacuumResponse is a zero-copy view into a ZAP-encoded
// DisableVacuumResponse message. DisableVacuumResponse carries no fields.
type DisableVacuumResponse struct{ o zap.Object }

// WrapDisableVacuumResponse parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapDisableVacuumResponse(b []byte) (DisableVacuumResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return DisableVacuumResponse{}, err
	}
	return DisableVacuumResponse{o: m.Root()}, nil
}

// DisableVacuumResponseInput collects the field values for
// NewDisableVacuumResponse. DisableVacuumResponse has no fields.
type DisableVacuumResponseInput struct{}

// NewDisableVacuumResponse builds a ZAP-encoded DisableVacuumResponse message
// and returns the bytes.
func NewDisableVacuumResponse(in DisableVacuumResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(disableVacuumResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- EnableVacuumRequest ---

const (
	enableVacuumRequestByPluginOff = 0
	enableVacuumRequestSize        = 1
)

// EnableVacuumRequest is a zero-copy view into a ZAP-encoded EnableVacuumRequest
// message.
type EnableVacuumRequest struct{ o zap.Object }

// WrapEnableVacuumRequest parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapEnableVacuumRequest(b []byte) (EnableVacuumRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return EnableVacuumRequest{}, err
	}
	return EnableVacuumRequest{o: m.Root()}, nil
}

// ByPlugin reads the by_plugin field (proto field 1, bool).
func (t EnableVacuumRequest) ByPlugin() bool { return t.o.Bool(enableVacuumRequestByPluginOff) }

// EnableVacuumRequestInput collects the field values for NewEnableVacuumRequest.
type EnableVacuumRequestInput struct {
	ByPlugin bool
}

// NewEnableVacuumRequest builds a ZAP-encoded EnableVacuumRequest message from
// in and returns the bytes.
func NewEnableVacuumRequest(in EnableVacuumRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(enableVacuumRequestSize)
	ob.SetBool(enableVacuumRequestByPluginOff, in.ByPlugin)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- EnableVacuumResponse ---

const enableVacuumResponseSize = 0

// EnableVacuumResponse is a zero-copy view into a ZAP-encoded
// EnableVacuumResponse message. EnableVacuumResponse carries no fields.
type EnableVacuumResponse struct{ o zap.Object }

// WrapEnableVacuumResponse parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapEnableVacuumResponse(b []byte) (EnableVacuumResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return EnableVacuumResponse{}, err
	}
	return EnableVacuumResponse{o: m.Root()}, nil
}

// EnableVacuumResponseInput collects the field values for
// NewEnableVacuumResponse. EnableVacuumResponse has no fields.
type EnableVacuumResponseInput struct{}

// NewEnableVacuumResponse builds a ZAP-encoded EnableVacuumResponse message and
// returns the bytes.
func NewEnableVacuumResponse(in EnableVacuumResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(enableVacuumResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeMarkReadonlyRequest ---

const (
	volumeMarkReadonlyRequestIpOff               = 0
	volumeMarkReadonlyRequestPortOff             = 8
	volumeMarkReadonlyRequestVolumeIdOff         = 12
	volumeMarkReadonlyRequestCollectionOff       = 16
	volumeMarkReadonlyRequestReplicaPlacementOff = 24
	volumeMarkReadonlyRequestVersionOff          = 28
	volumeMarkReadonlyRequestTtlOff              = 32
	volumeMarkReadonlyRequestDiskTypeOff         = 36
	volumeMarkReadonlyRequestIsReadonlyOff       = 44
	volumeMarkReadonlyRequestSize                = 45
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

// Ip reads the ip field (proto field 1, string).
func (t VolumeMarkReadonlyRequest) Ip() string { return t.o.Text(volumeMarkReadonlyRequestIpOff) }

// Port reads the port field (proto field 2, uint32).
func (t VolumeMarkReadonlyRequest) Port() uint32 {
	return t.o.Uint32(volumeMarkReadonlyRequestPortOff)
}

// VolumeId reads the volume_id field (proto field 4, uint32).
func (t VolumeMarkReadonlyRequest) VolumeId() uint32 {
	return t.o.Uint32(volumeMarkReadonlyRequestVolumeIdOff)
}

// Collection reads the collection field (proto field 5, string).
func (t VolumeMarkReadonlyRequest) Collection() string {
	return t.o.Text(volumeMarkReadonlyRequestCollectionOff)
}

// ReplicaPlacement reads the replica_placement field (proto field 6, uint32).
func (t VolumeMarkReadonlyRequest) ReplicaPlacement() uint32 {
	return t.o.Uint32(volumeMarkReadonlyRequestReplicaPlacementOff)
}

// Version reads the version field (proto field 7, uint32).
func (t VolumeMarkReadonlyRequest) Version() uint32 {
	return t.o.Uint32(volumeMarkReadonlyRequestVersionOff)
}

// Ttl reads the ttl field (proto field 8, uint32).
func (t VolumeMarkReadonlyRequest) Ttl() uint32 {
	return t.o.Uint32(volumeMarkReadonlyRequestTtlOff)
}

// DiskType reads the disk_type field (proto field 9, string).
func (t VolumeMarkReadonlyRequest) DiskType() string {
	return t.o.Text(volumeMarkReadonlyRequestDiskTypeOff)
}

// IsReadonly reads the is_readonly field (proto field 10, bool).
func (t VolumeMarkReadonlyRequest) IsReadonly() bool {
	return t.o.Bool(volumeMarkReadonlyRequestIsReadonlyOff)
}

// VolumeMarkReadonlyRequestInput collects the field values for
// NewVolumeMarkReadonlyRequest.
type VolumeMarkReadonlyRequestInput struct {
	Ip               string
	Port             uint32
	VolumeId         uint32
	Collection       string
	ReplicaPlacement uint32
	Version          uint32
	Ttl              uint32
	DiskType         string
	IsReadonly       bool
}

// NewVolumeMarkReadonlyRequest builds a ZAP-encoded VolumeMarkReadonlyRequest
// message from in and returns the bytes.
func NewVolumeMarkReadonlyRequest(in VolumeMarkReadonlyRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeMarkReadonlyRequestSize)
	ob.SetText(volumeMarkReadonlyRequestIpOff, in.Ip)
	ob.SetUint32(volumeMarkReadonlyRequestPortOff, in.Port)
	ob.SetUint32(volumeMarkReadonlyRequestVolumeIdOff, in.VolumeId)
	ob.SetText(volumeMarkReadonlyRequestCollectionOff, in.Collection)
	ob.SetUint32(volumeMarkReadonlyRequestReplicaPlacementOff, in.ReplicaPlacement)
	ob.SetUint32(volumeMarkReadonlyRequestVersionOff, in.Version)
	ob.SetUint32(volumeMarkReadonlyRequestTtlOff, in.Ttl)
	ob.SetText(volumeMarkReadonlyRequestDiskTypeOff, in.DiskType)
	ob.SetBool(volumeMarkReadonlyRequestIsReadonlyOff, in.IsReadonly)
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
