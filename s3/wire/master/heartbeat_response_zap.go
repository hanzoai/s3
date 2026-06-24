// Code generated from master.proto; DO NOT EDIT.

package masterwire

import (
	zap "github.com/zap-proto/go"
)

// --- StorageBackend ---

const (
	storageBackendTypeOff       = 0
	storageBackendIdOff         = 8
	storageBackendPropertiesOff = 16
	storageBackendSize          = 24
)

// StorageBackend is a zero-copy view into a ZAP-encoded StorageBackend message.
type StorageBackend struct{ o zap.Object }

// WrapStorageBackend parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapStorageBackend(b []byte) (StorageBackend, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return StorageBackend{}, err
	}
	return StorageBackend{o: m.Root()}, nil
}

// Type reads the type field (proto field 1, string).
func (t StorageBackend) Type() string { return t.o.Text(storageBackendTypeOff) }

// Id reads the id field (proto field 2, string).
func (t StorageBackend) Id() string { return t.o.Text(storageBackendIdOff) }

// PropertiesLen reports the number of properties entries (proto field 3,
// map<string,string>).
func (t StorageBackend) PropertiesLen() int { return t.o.List(storageBackendPropertiesOff).Len() }

// PropertyAt returns the i-th properties entry.
func (t StorageBackend) PropertyAt(i int) StringStringEntry {
	return stringStringEntryAt(t.o.List(storageBackendPropertiesOff), i)
}

// StorageBackendInput collects the field values for NewStorageBackend.
type StorageBackendInput struct {
	Type       string
	Id         string
	Properties []StringStringEntry
}

// NewStorageBackend builds a ZAP-encoded StorageBackend message from in and
// returns the bytes.
func NewStorageBackend(in StorageBackendInput) []byte {
	b := zap.NewBuilder(256)
	propertiesOff, propertiesLen := setStringStringMap(b, in.Properties)
	ob := b.StartObject(storageBackendSize)
	ob.SetText(storageBackendTypeOff, in.Type)
	ob.SetText(storageBackendIdOff, in.Id)
	ob.SetList(storageBackendPropertiesOff, propertiesOff, propertiesLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- HeartbeatResponse ---

const (
	heartbeatResponseVolumeSizeLimitOff        = 0
	heartbeatResponseLeaderOff                 = 8
	heartbeatResponseMetricsAddressOff         = 16
	heartbeatResponseMetricsIntervalSecondsOff = 24
	heartbeatResponseStorageBackendsOff        = 28
	heartbeatResponseDuplicatedUuidsOff        = 36
	heartbeatResponsePreallocateOff            = 44
	heartbeatResponseSize                      = 45
)

// HeartbeatResponse is a zero-copy view into a ZAP-encoded HeartbeatResponse
// message.
type HeartbeatResponse struct{ o zap.Object }

// WrapHeartbeatResponse parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapHeartbeatResponse(b []byte) (HeartbeatResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return HeartbeatResponse{}, err
	}
	return HeartbeatResponse{o: m.Root()}, nil
}

// VolumeSizeLimit reads the volume_size_limit field (proto field 1, uint64).
func (t HeartbeatResponse) VolumeSizeLimit() uint64 {
	return t.o.Uint64(heartbeatResponseVolumeSizeLimitOff)
}

// Leader reads the leader field (proto field 2, string).
func (t HeartbeatResponse) Leader() string { return t.o.Text(heartbeatResponseLeaderOff) }

// MetricsAddress reads the metrics_address field (proto field 3, string).
func (t HeartbeatResponse) MetricsAddress() string {
	return t.o.Text(heartbeatResponseMetricsAddressOff)
}

// MetricsIntervalSeconds reads the metrics_interval_seconds field (proto field
// 4, uint32).
func (t HeartbeatResponse) MetricsIntervalSeconds() uint32 {
	return t.o.Uint32(heartbeatResponseMetricsIntervalSecondsOff)
}

// StorageBackendsLen reports the number of storage_backends elements (proto
// field 5, repeated StorageBackend).
func (t HeartbeatResponse) StorageBackendsLen() int {
	return t.o.List(heartbeatResponseStorageBackendsOff).Len()
}

// StorageBackendAt returns the i-th storage_backends element.
func (t HeartbeatResponse) StorageBackendAt(i int) (StorageBackend, bool) {
	o := t.o.List(heartbeatResponseStorageBackendsOff).ObjectAt(i)
	if o.IsNull() {
		return StorageBackend{}, false
	}
	return StorageBackend{o: o}, true
}

// DuplicatedUuidsLen reports the number of duplicated_uuids elements (proto
// field 6, repeated string).
func (t HeartbeatResponse) DuplicatedUuidsLen() int {
	return t.o.List(heartbeatResponseDuplicatedUuidsOff).Len()
}

// DuplicatedUuidAt returns the i-th duplicated_uuids element.
func (t HeartbeatResponse) DuplicatedUuidAt(i int) string {
	return stringAt(t.o.List(heartbeatResponseDuplicatedUuidsOff), i)
}

// Preallocate reads the preallocate field (proto field 7, bool).
func (t HeartbeatResponse) Preallocate() bool { return t.o.Bool(heartbeatResponsePreallocateOff) }

// HeartbeatResponseInput collects the field values for NewHeartbeatResponse.
type HeartbeatResponseInput struct {
	VolumeSizeLimit        uint64
	Leader                 string
	MetricsAddress         string
	MetricsIntervalSeconds uint32
	StorageBackends        [][]byte
	DuplicatedUuids        []string
	Preallocate            bool
}

// NewHeartbeatResponse builds a ZAP-encoded HeartbeatResponse message from in
// and returns the bytes.
func NewHeartbeatResponse(in HeartbeatResponseInput) []byte {
	b := zap.NewBuilder(256)
	storageBackendsOff, storageBackendsLen := setMsgList(b, in.StorageBackends)
	duplicatedUuidsOff, duplicatedUuidsLen := setStringList(b, in.DuplicatedUuids)
	ob := b.StartObject(heartbeatResponseSize)
	ob.SetUint64(heartbeatResponseVolumeSizeLimitOff, in.VolumeSizeLimit)
	ob.SetText(heartbeatResponseLeaderOff, in.Leader)
	ob.SetText(heartbeatResponseMetricsAddressOff, in.MetricsAddress)
	ob.SetUint32(heartbeatResponseMetricsIntervalSecondsOff, in.MetricsIntervalSeconds)
	ob.SetList(heartbeatResponseStorageBackendsOff, storageBackendsOff, storageBackendsLen)
	ob.SetList(heartbeatResponseDuplicatedUuidsOff, duplicatedUuidsOff, duplicatedUuidsLen)
	ob.SetBool(heartbeatResponsePreallocateOff, in.Preallocate)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- Empty ---

const emptySize = 0

// Empty is a zero-copy view into a ZAP-encoded Empty message. Empty carries no
// fields.
type Empty struct{ o zap.Object }

// WrapEmpty parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapEmpty(b []byte) (Empty, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return Empty{}, err
	}
	return Empty{o: m.Root()}, nil
}

// EmptyInput collects the field values for NewEmpty. Empty has no fields.
type EmptyInput struct{}

// NewEmpty builds a ZAP-encoded Empty message and returns the bytes.
func NewEmpty(in EmptyInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(emptySize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- SuperBlockExtra.ErasureCoding ---

const (
	superBlockExtraErasureCodingDataOff      = 0
	superBlockExtraErasureCodingParityOff    = 4
	superBlockExtraErasureCodingVolumeIdsOff = 8
	superBlockExtraErasureCodingSize         = 16
)

// SuperBlockExtraErasureCoding is a zero-copy view into a ZAP-encoded
// SuperBlockExtra.ErasureCoding message.
type SuperBlockExtraErasureCoding struct{ o zap.Object }

// WrapSuperBlockExtraErasureCoding parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapSuperBlockExtraErasureCoding(b []byte) (SuperBlockExtraErasureCoding, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return SuperBlockExtraErasureCoding{}, err
	}
	return SuperBlockExtraErasureCoding{o: m.Root()}, nil
}

// Data reads the data field (proto field 1, uint32).
func (t SuperBlockExtraErasureCoding) Data() uint32 {
	return t.o.Uint32(superBlockExtraErasureCodingDataOff)
}

// Parity reads the parity field (proto field 2, uint32).
func (t SuperBlockExtraErasureCoding) Parity() uint32 {
	return t.o.Uint32(superBlockExtraErasureCodingParityOff)
}

// VolumeIdsLen reports the number of volume_ids elements (proto field 3,
// repeated uint32).
func (t SuperBlockExtraErasureCoding) VolumeIdsLen() int {
	return t.o.List(superBlockExtraErasureCodingVolumeIdsOff).Len()
}

// VolumeIdAt returns the i-th volume_ids element.
func (t SuperBlockExtraErasureCoding) VolumeIdAt(i int) uint32 {
	return t.o.List(superBlockExtraErasureCodingVolumeIdsOff).Uint32(i)
}

// SuperBlockExtraErasureCodingInput collects the field values for
// NewSuperBlockExtraErasureCoding.
type SuperBlockExtraErasureCodingInput struct {
	Data      uint32
	Parity    uint32
	VolumeIds []uint32
}

// NewSuperBlockExtraErasureCoding builds a ZAP-encoded
// SuperBlockExtra.ErasureCoding message from in and returns the bytes.
func NewSuperBlockExtraErasureCoding(in SuperBlockExtraErasureCodingInput) []byte {
	b := zap.NewBuilder(256)
	volumeIdsOff, volumeIdsLen := setUint32List(b, in.VolumeIds)
	ob := b.StartObject(superBlockExtraErasureCodingSize)
	ob.SetUint32(superBlockExtraErasureCodingDataOff, in.Data)
	ob.SetUint32(superBlockExtraErasureCodingParityOff, in.Parity)
	ob.SetList(superBlockExtraErasureCodingVolumeIdsOff, volumeIdsOff, volumeIdsLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- SuperBlockExtra ---

const (
	superBlockExtraErasureCodingFieldOff = 0
	superBlockExtraSize                  = 8
)

// SuperBlockExtra is a zero-copy view into a ZAP-encoded SuperBlockExtra message.
type SuperBlockExtra struct{ o zap.Object }

// WrapSuperBlockExtra parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapSuperBlockExtra(b []byte) (SuperBlockExtra, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return SuperBlockExtra{}, err
	}
	return SuperBlockExtra{o: m.Root()}, nil
}

// ErasureCoding reads the erasure_coding field (proto field 1, message
// SuperBlockExtra.ErasureCoding). The bool is false when the field is absent.
func (t SuperBlockExtra) ErasureCoding() (SuperBlockExtraErasureCoding, bool) {
	b := t.o.Bytes(superBlockExtraErasureCodingFieldOff)
	if len(b) == 0 {
		return SuperBlockExtraErasureCoding{}, false
	}
	ec, err := WrapSuperBlockExtraErasureCoding(b)
	if err != nil {
		return SuperBlockExtraErasureCoding{}, false
	}
	return ec, true
}

// SuperBlockExtraInput collects the field values for NewSuperBlockExtra.
// ErasureCoding is the child message's own ZAP buffer (from
// NewSuperBlockExtraErasureCoding); nil encodes an absent erasure_coding.
type SuperBlockExtraInput struct {
	ErasureCoding []byte
}

// NewSuperBlockExtra builds a ZAP-encoded SuperBlockExtra message from in and
// returns the bytes.
func NewSuperBlockExtra(in SuperBlockExtraInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(superBlockExtraSize)
	ob.SetBytes(superBlockExtraErasureCodingFieldOff, in.ErasureCoding)
	ob.FinishAsRoot()
	return b.Finish()
}
