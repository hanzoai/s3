// Code generated from master.proto; DO NOT EDIT.

package masterwire

import (
	zap "github.com/zap-proto/go"

	volume_serverwire "github.com/hanzoai/s3/s3/wire/volume_server"
)

// --- DiskTag ---

const (
	diskTagDiskIdOff = 0
	diskTagTagsOff   = 4
	diskTagSize      = 12
)

// DiskTag is a zero-copy view into a ZAP-encoded DiskTag message.
type DiskTag struct{ o zap.Object }

// WrapDiskTag parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapDiskTag(b []byte) (DiskTag, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return DiskTag{}, err
	}
	return DiskTag{o: m.Root()}, nil
}

// DiskId reads the disk_id field (proto field 1, uint32).
func (t DiskTag) DiskId() uint32 { return t.o.Uint32(diskTagDiskIdOff) }

// TagsLen reports the number of tags elements (proto field 2, repeated string).
func (t DiskTag) TagsLen() int { return t.o.List(diskTagTagsOff).Len() }

// TagAt returns the i-th tags element.
func (t DiskTag) TagAt(i int) string { return stringAt(t.o.List(diskTagTagsOff), i) }

// DiskTagInput collects the field values for NewDiskTag.
type DiskTagInput struct {
	DiskId uint32
	Tags   []string
}

// NewDiskTag builds a ZAP-encoded DiskTag message from in and returns the bytes.
func NewDiskTag(in DiskTagInput) []byte {
	b := zap.NewBuilder(256)
	tagsOff, tagsLen := setStringList(b, in.Tags)
	ob := b.StartObject(diskTagSize)
	ob.SetUint32(diskTagDiskIdOff, in.DiskId)
	ob.SetList(diskTagTagsOff, tagsOff, tagsLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- Heartbeat ---

const (
	heartbeatIpOff              = 0
	heartbeatPortOff            = 8
	heartbeatPublicUrlOff       = 12
	heartbeatMaxFileKeyOff      = 20
	heartbeatDataCenterOff      = 28
	heartbeatRackOff            = 36
	heartbeatAdminPortOff       = 44
	heartbeatVolumesOff         = 48
	heartbeatNewVolumesOff      = 56
	heartbeatDeletedVolumesOff  = 64
	heartbeatHasNoVolumesOff    = 72
	heartbeatEcShardsOff        = 73
	heartbeatNewEcShardsOff     = 81
	heartbeatDeletedEcShardsOff = 89
	heartbeatHasNoEcShardsOff   = 97
	heartbeatMaxVolumeCountsOff = 98
	heartbeatGrpcPortOff        = 106
	heartbeatLocationUuidsOff   = 110
	heartbeatIdOff              = 118
	heartbeatStateOff           = 126
	heartbeatDiskTagsOff        = 134
	heartbeatSize               = 142
)

// Heartbeat is a zero-copy view into a ZAP-encoded Heartbeat message.
type Heartbeat struct{ o zap.Object }

// WrapHeartbeat parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapHeartbeat(b []byte) (Heartbeat, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return Heartbeat{}, err
	}
	return Heartbeat{o: m.Root()}, nil
}

// Ip reads the ip field (proto field 1, string).
func (t Heartbeat) Ip() string { return t.o.Text(heartbeatIpOff) }

// Port reads the port field (proto field 2, uint32).
func (t Heartbeat) Port() uint32 { return t.o.Uint32(heartbeatPortOff) }

// PublicUrl reads the public_url field (proto field 3, string).
func (t Heartbeat) PublicUrl() string { return t.o.Text(heartbeatPublicUrlOff) }

// MaxFileKey reads the max_file_key field (proto field 5, uint64).
func (t Heartbeat) MaxFileKey() uint64 { return t.o.Uint64(heartbeatMaxFileKeyOff) }

// DataCenter reads the data_center field (proto field 6, string).
func (t Heartbeat) DataCenter() string { return t.o.Text(heartbeatDataCenterOff) }

// Rack reads the rack field (proto field 7, string).
func (t Heartbeat) Rack() string { return t.o.Text(heartbeatRackOff) }

// AdminPort reads the admin_port field (proto field 8, uint32).
func (t Heartbeat) AdminPort() uint32 { return t.o.Uint32(heartbeatAdminPortOff) }

// VolumesLen reports the number of volumes elements (proto field 9, repeated
// VolumeInformationMessage).
func (t Heartbeat) VolumesLen() int { return t.o.List(heartbeatVolumesOff).Len() }

// VolumeAt returns the i-th volumes element. The bool is false when i is out of
// range or the element fails to parse.
func (t Heartbeat) VolumeAt(i int) (VolumeInformationMessage, bool) {
	o := t.o.List(heartbeatVolumesOff).ObjectAt(i)
	if o.IsNull() {
		return VolumeInformationMessage{}, false
	}
	return VolumeInformationMessage{o: o}, true
}

// NewVolumesLen reports the number of new_volumes elements (proto field 10,
// repeated VolumeShortInformationMessage).
func (t Heartbeat) NewVolumesLen() int { return t.o.List(heartbeatNewVolumesOff).Len() }

// NewVolumeAt returns the i-th new_volumes element.
func (t Heartbeat) NewVolumeAt(i int) (VolumeShortInformationMessage, bool) {
	o := t.o.List(heartbeatNewVolumesOff).ObjectAt(i)
	if o.IsNull() {
		return VolumeShortInformationMessage{}, false
	}
	return VolumeShortInformationMessage{o: o}, true
}

// DeletedVolumesLen reports the number of deleted_volumes elements (proto field
// 11, repeated VolumeShortInformationMessage).
func (t Heartbeat) DeletedVolumesLen() int { return t.o.List(heartbeatDeletedVolumesOff).Len() }

// DeletedVolumeAt returns the i-th deleted_volumes element.
func (t Heartbeat) DeletedVolumeAt(i int) (VolumeShortInformationMessage, bool) {
	o := t.o.List(heartbeatDeletedVolumesOff).ObjectAt(i)
	if o.IsNull() {
		return VolumeShortInformationMessage{}, false
	}
	return VolumeShortInformationMessage{o: o}, true
}

// HasNoVolumes reads the has_no_volumes field (proto field 12, bool).
func (t Heartbeat) HasNoVolumes() bool { return t.o.Bool(heartbeatHasNoVolumesOff) }

// EcShardsLen reports the number of ec_shards elements (proto field 16, repeated
// VolumeEcShardInformationMessage).
func (t Heartbeat) EcShardsLen() int { return t.o.List(heartbeatEcShardsOff).Len() }

// EcShardAt returns the i-th ec_shards element.
func (t Heartbeat) EcShardAt(i int) (VolumeEcShardInformationMessage, bool) {
	o := t.o.List(heartbeatEcShardsOff).ObjectAt(i)
	if o.IsNull() {
		return VolumeEcShardInformationMessage{}, false
	}
	return VolumeEcShardInformationMessage{o: o}, true
}

// NewEcShardsLen reports the number of new_ec_shards elements (proto field 17,
// repeated VolumeEcShardInformationMessage).
func (t Heartbeat) NewEcShardsLen() int { return t.o.List(heartbeatNewEcShardsOff).Len() }

// NewEcShardAt returns the i-th new_ec_shards element.
func (t Heartbeat) NewEcShardAt(i int) (VolumeEcShardInformationMessage, bool) {
	o := t.o.List(heartbeatNewEcShardsOff).ObjectAt(i)
	if o.IsNull() {
		return VolumeEcShardInformationMessage{}, false
	}
	return VolumeEcShardInformationMessage{o: o}, true
}

// DeletedEcShardsLen reports the number of deleted_ec_shards elements (proto
// field 18, repeated VolumeEcShardInformationMessage).
func (t Heartbeat) DeletedEcShardsLen() int { return t.o.List(heartbeatDeletedEcShardsOff).Len() }

// DeletedEcShardAt returns the i-th deleted_ec_shards element.
func (t Heartbeat) DeletedEcShardAt(i int) (VolumeEcShardInformationMessage, bool) {
	o := t.o.List(heartbeatDeletedEcShardsOff).ObjectAt(i)
	if o.IsNull() {
		return VolumeEcShardInformationMessage{}, false
	}
	return VolumeEcShardInformationMessage{o: o}, true
}

// HasNoEcShards reads the has_no_ec_shards field (proto field 19, bool).
func (t Heartbeat) HasNoEcShards() bool { return t.o.Bool(heartbeatHasNoEcShardsOff) }

// MaxVolumeCountsLen reports the number of max_volume_counts entries (proto
// field 4, map<string,uint32>).
func (t Heartbeat) MaxVolumeCountsLen() int { return t.o.List(heartbeatMaxVolumeCountsOff).Len() }

// MaxVolumeCountAt returns the i-th max_volume_counts entry.
func (t Heartbeat) MaxVolumeCountAt(i int) StringUint32Entry {
	return stringUint32EntryAt(t.o.List(heartbeatMaxVolumeCountsOff), i)
}

// GrpcPort reads the grpc_port field (proto field 20, uint32).
func (t Heartbeat) GrpcPort() uint32 { return t.o.Uint32(heartbeatGrpcPortOff) }

// LocationUuidsLen reports the number of location_uuids elements (proto field
// 21, repeated string).
func (t Heartbeat) LocationUuidsLen() int { return t.o.List(heartbeatLocationUuidsOff).Len() }

// LocationUuidAt returns the i-th location_uuids element.
func (t Heartbeat) LocationUuidAt(i int) string {
	return stringAt(t.o.List(heartbeatLocationUuidsOff), i)
}

// Id reads the id field (proto field 22, string).
func (t Heartbeat) Id() string { return t.o.Text(heartbeatIdOff) }

// State reads the state field (proto field 23, message
// volume_server_pb.VolumeServerState). The bool is false when the field is
// absent. The value is decoded with volume_serverwire.WrapVolumeServerState.
func (t Heartbeat) State() (volume_serverwire.VolumeServerState, bool) {
	b := t.o.Bytes(heartbeatStateOff)
	if len(b) == 0 {
		return volume_serverwire.VolumeServerState{}, false
	}
	s, err := volume_serverwire.WrapVolumeServerState(b)
	if err != nil {
		return volume_serverwire.VolumeServerState{}, false
	}
	return s, true
}

// DiskTagsLen reports the number of disk_tags elements (proto field 24, repeated
// DiskTag).
func (t Heartbeat) DiskTagsLen() int { return t.o.List(heartbeatDiskTagsOff).Len() }

// DiskTagAt returns the i-th disk_tags element.
func (t Heartbeat) DiskTagAt(i int) (DiskTag, bool) {
	o := t.o.List(heartbeatDiskTagsOff).ObjectAt(i)
	if o.IsNull() {
		return DiskTag{}, false
	}
	return DiskTag{o: o}, true
}

// HeartbeatInput collects the field values for NewHeartbeat. Repeated-message
// and map fields take each element as its own ZAP sub-buffer; State is the
// VolumeServerState message's own ZAP buffer (from
// volume_serverwire.NewVolumeServerState), nil for absent.
type HeartbeatInput struct {
	Ip              string
	Port            uint32
	PublicUrl       string
	MaxFileKey      uint64
	DataCenter      string
	Rack            string
	AdminPort       uint32
	Volumes         [][]byte
	NewVolumes      [][]byte
	DeletedVolumes  [][]byte
	HasNoVolumes    bool
	EcShards        [][]byte
	NewEcShards     [][]byte
	DeletedEcShards [][]byte
	HasNoEcShards   bool
	MaxVolumeCounts []StringUint32Entry
	GrpcPort        uint32
	LocationUuids   []string
	Id              string
	State           []byte
	DiskTags        [][]byte
}

// NewHeartbeat builds a ZAP-encoded Heartbeat message from in and returns the
// bytes.
func NewHeartbeat(in HeartbeatInput) []byte {
	b := zap.NewBuilder(512)
	volumesOff, volumesLen := setMsgList(b, in.Volumes)
	newVolumesOff, newVolumesLen := setMsgList(b, in.NewVolumes)
	deletedVolumesOff, deletedVolumesLen := setMsgList(b, in.DeletedVolumes)
	ecShardsOff, ecShardsLen := setMsgList(b, in.EcShards)
	newEcShardsOff, newEcShardsLen := setMsgList(b, in.NewEcShards)
	deletedEcShardsOff, deletedEcShardsLen := setMsgList(b, in.DeletedEcShards)
	maxVolumeCountsOff, maxVolumeCountsLen := setStringUint32Map(b, in.MaxVolumeCounts)
	locationUuidsOff, locationUuidsLen := setStringList(b, in.LocationUuids)
	diskTagsOff, diskTagsLen := setMsgList(b, in.DiskTags)
	ob := b.StartObject(heartbeatSize)
	ob.SetText(heartbeatIpOff, in.Ip)
	ob.SetUint32(heartbeatPortOff, in.Port)
	ob.SetText(heartbeatPublicUrlOff, in.PublicUrl)
	ob.SetUint64(heartbeatMaxFileKeyOff, in.MaxFileKey)
	ob.SetText(heartbeatDataCenterOff, in.DataCenter)
	ob.SetText(heartbeatRackOff, in.Rack)
	ob.SetUint32(heartbeatAdminPortOff, in.AdminPort)
	ob.SetList(heartbeatVolumesOff, volumesOff, volumesLen)
	ob.SetList(heartbeatNewVolumesOff, newVolumesOff, newVolumesLen)
	ob.SetList(heartbeatDeletedVolumesOff, deletedVolumesOff, deletedVolumesLen)
	ob.SetBool(heartbeatHasNoVolumesOff, in.HasNoVolumes)
	ob.SetList(heartbeatEcShardsOff, ecShardsOff, ecShardsLen)
	ob.SetList(heartbeatNewEcShardsOff, newEcShardsOff, newEcShardsLen)
	ob.SetList(heartbeatDeletedEcShardsOff, deletedEcShardsOff, deletedEcShardsLen)
	ob.SetBool(heartbeatHasNoEcShardsOff, in.HasNoEcShards)
	ob.SetList(heartbeatMaxVolumeCountsOff, maxVolumeCountsOff, maxVolumeCountsLen)
	ob.SetUint32(heartbeatGrpcPortOff, in.GrpcPort)
	ob.SetList(heartbeatLocationUuidsOff, locationUuidsOff, locationUuidsLen)
	ob.SetText(heartbeatIdOff, in.Id)
	ob.SetBytes(heartbeatStateOff, in.State)
	ob.SetList(heartbeatDiskTagsOff, diskTagsOff, diskTagsLen)
	ob.FinishAsRoot()
	return b.Finish()
}
