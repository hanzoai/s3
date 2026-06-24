// Code generated from master.proto; DO NOT EDIT.

package masterwire

import (
	zap "github.com/zap-proto/go"
)

// --- DiskInfo ---

const (
	diskInfoTypeOff              = 0
	diskInfoVolumeCountOff       = 8
	diskInfoMaxVolumeCountOff    = 16
	diskInfoFreeVolumeCountOff   = 24
	diskInfoActiveVolumeCountOff = 32
	diskInfoVolumeInfosOff       = 40
	diskInfoEcShardInfosOff      = 48
	diskInfoRemoteVolumeCountOff = 56
	diskInfoDiskIdOff            = 64
	diskInfoTagsOff              = 68
	diskInfoSize                 = 76
)

// DiskInfo is a zero-copy view into a ZAP-encoded DiskInfo message.
type DiskInfo struct{ o zap.Object }

// WrapDiskInfo parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapDiskInfo(b []byte) (DiskInfo, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return DiskInfo{}, err
	}
	return DiskInfo{o: m.Root()}, nil
}

// Type reads the type field (proto field 1, string).
func (t DiskInfo) Type() string { return t.o.Text(diskInfoTypeOff) }

// VolumeCount reads the volume_count field (proto field 2, int64).
func (t DiskInfo) VolumeCount() int64 { return t.o.Int64(diskInfoVolumeCountOff) }

// MaxVolumeCount reads the max_volume_count field (proto field 3, int64).
func (t DiskInfo) MaxVolumeCount() int64 { return t.o.Int64(diskInfoMaxVolumeCountOff) }

// FreeVolumeCount reads the free_volume_count field (proto field 4, int64).
func (t DiskInfo) FreeVolumeCount() int64 { return t.o.Int64(diskInfoFreeVolumeCountOff) }

// ActiveVolumeCount reads the active_volume_count field (proto field 5, int64).
func (t DiskInfo) ActiveVolumeCount() int64 { return t.o.Int64(diskInfoActiveVolumeCountOff) }

// VolumeInfosLen reports the number of volume_infos elements (proto field 6,
// repeated VolumeInformationMessage).
func (t DiskInfo) VolumeInfosLen() int { return t.o.List(diskInfoVolumeInfosOff).Len() }

// VolumeInfoAt returns the i-th volume_infos element.
func (t DiskInfo) VolumeInfoAt(i int) (VolumeInformationMessage, bool) {
	o := t.o.List(diskInfoVolumeInfosOff).ObjectAt(i)
	if o.IsNull() {
		return VolumeInformationMessage{}, false
	}
	return VolumeInformationMessage{o: o}, true
}

// EcShardInfosLen reports the number of ec_shard_infos elements (proto field 7,
// repeated VolumeEcShardInformationMessage).
func (t DiskInfo) EcShardInfosLen() int { return t.o.List(diskInfoEcShardInfosOff).Len() }

// EcShardInfoAt returns the i-th ec_shard_infos element.
func (t DiskInfo) EcShardInfoAt(i int) (VolumeEcShardInformationMessage, bool) {
	o := t.o.List(diskInfoEcShardInfosOff).ObjectAt(i)
	if o.IsNull() {
		return VolumeEcShardInformationMessage{}, false
	}
	return VolumeEcShardInformationMessage{o: o}, true
}

// RemoteVolumeCount reads the remote_volume_count field (proto field 8, int64).
func (t DiskInfo) RemoteVolumeCount() int64 { return t.o.Int64(diskInfoRemoteVolumeCountOff) }

// DiskId reads the disk_id field (proto field 9, uint32).
func (t DiskInfo) DiskId() uint32 { return t.o.Uint32(diskInfoDiskIdOff) }

// TagsLen reports the number of tags elements (proto field 10, repeated string).
func (t DiskInfo) TagsLen() int { return t.o.List(diskInfoTagsOff).Len() }

// TagAt returns the i-th tags element.
func (t DiskInfo) TagAt(i int) string { return stringAt(t.o.List(diskInfoTagsOff), i) }

// DiskInfoInput collects the field values for NewDiskInfo. VolumeInfos and
// EcShardInfos take each element as its own ZAP sub-buffer.
type DiskInfoInput struct {
	Type              string
	VolumeCount       int64
	MaxVolumeCount    int64
	FreeVolumeCount   int64
	ActiveVolumeCount int64
	VolumeInfos       [][]byte
	EcShardInfos      [][]byte
	RemoteVolumeCount int64
	DiskId            uint32
	Tags              []string
}

// NewDiskInfo builds a ZAP-encoded DiskInfo message from in and returns the
// bytes.
func NewDiskInfo(in DiskInfoInput) []byte {
	b := zap.NewBuilder(256)
	volumeInfosOff, volumeInfosLen := setMsgList(b, in.VolumeInfos)
	ecShardInfosOff, ecShardInfosLen := setMsgList(b, in.EcShardInfos)
	tagsOff, tagsLen := setStringList(b, in.Tags)
	ob := b.StartObject(diskInfoSize)
	ob.SetText(diskInfoTypeOff, in.Type)
	ob.SetInt64(diskInfoVolumeCountOff, in.VolumeCount)
	ob.SetInt64(diskInfoMaxVolumeCountOff, in.MaxVolumeCount)
	ob.SetInt64(diskInfoFreeVolumeCountOff, in.FreeVolumeCount)
	ob.SetInt64(diskInfoActiveVolumeCountOff, in.ActiveVolumeCount)
	ob.SetList(diskInfoVolumeInfosOff, volumeInfosOff, volumeInfosLen)
	ob.SetList(diskInfoEcShardInfosOff, ecShardInfosOff, ecShardInfosLen)
	ob.SetInt64(diskInfoRemoteVolumeCountOff, in.RemoteVolumeCount)
	ob.SetUint32(diskInfoDiskIdOff, in.DiskId)
	ob.SetList(diskInfoTagsOff, tagsOff, tagsLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- DataNodeInfo ---

const (
	dataNodeInfoIdOff        = 0
	dataNodeInfoDiskInfosOff = 8
	dataNodeInfoGrpcPortOff  = 16
	dataNodeInfoAddressOff   = 20
	dataNodeInfoSize         = 28
)

// DataNodeInfo is a zero-copy view into a ZAP-encoded DataNodeInfo message.
type DataNodeInfo struct{ o zap.Object }

// WrapDataNodeInfo parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapDataNodeInfo(b []byte) (DataNodeInfo, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return DataNodeInfo{}, err
	}
	return DataNodeInfo{o: m.Root()}, nil
}

// Id reads the id field (proto field 1, string).
func (t DataNodeInfo) Id() string { return t.o.Text(dataNodeInfoIdOff) }

// DiskInfosLen reports the number of diskInfos entries (proto field 2,
// map<string,DiskInfo>).
func (t DataNodeInfo) DiskInfosLen() int { return t.o.List(dataNodeInfoDiskInfosOff).Len() }

// DiskInfoAt returns the i-th diskInfos entry. The entry's Value is the raw
// DiskInfo sub-message bytes; decode it with WrapDiskInfo.
func (t DataNodeInfo) DiskInfoAt(i int) StringMsgEntry {
	return stringMsgEntryAt(t.o.List(dataNodeInfoDiskInfosOff), i)
}

// GrpcPort reads the grpc_port field (proto field 3, uint32).
func (t DataNodeInfo) GrpcPort() uint32 { return t.o.Uint32(dataNodeInfoGrpcPortOff) }

// Address reads the address field (proto field 4, string).
func (t DataNodeInfo) Address() string { return t.o.Text(dataNodeInfoAddressOff) }

// DataNodeInfoInput collects the field values for NewDataNodeInfo. Each
// DiskInfos entry's Value is a DiskInfo message's own ZAP buffer (from
// NewDiskInfo).
type DataNodeInfoInput struct {
	Id        string
	DiskInfos []StringMsgEntry
	GrpcPort  uint32
	Address   string
}

// NewDataNodeInfo builds a ZAP-encoded DataNodeInfo message from in and returns
// the bytes.
func NewDataNodeInfo(in DataNodeInfoInput) []byte {
	b := zap.NewBuilder(256)
	diskInfosOff, diskInfosLen := setStringMsgMap(b, in.DiskInfos)
	ob := b.StartObject(dataNodeInfoSize)
	ob.SetText(dataNodeInfoIdOff, in.Id)
	ob.SetList(dataNodeInfoDiskInfosOff, diskInfosOff, diskInfosLen)
	ob.SetUint32(dataNodeInfoGrpcPortOff, in.GrpcPort)
	ob.SetText(dataNodeInfoAddressOff, in.Address)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- RackInfo ---

const (
	rackInfoIdOff            = 0
	rackInfoDataNodeInfosOff = 8
	rackInfoDiskInfosOff     = 16
	rackInfoSize             = 24
)

// RackInfo is a zero-copy view into a ZAP-encoded RackInfo message.
type RackInfo struct{ o zap.Object }

// WrapRackInfo parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapRackInfo(b []byte) (RackInfo, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return RackInfo{}, err
	}
	return RackInfo{o: m.Root()}, nil
}

// Id reads the id field (proto field 1, string).
func (t RackInfo) Id() string { return t.o.Text(rackInfoIdOff) }

// DataNodeInfosLen reports the number of data_node_infos elements (proto field
// 2, repeated DataNodeInfo).
func (t RackInfo) DataNodeInfosLen() int { return t.o.List(rackInfoDataNodeInfosOff).Len() }

// DataNodeInfoAt returns the i-th data_node_infos element.
func (t RackInfo) DataNodeInfoAt(i int) (DataNodeInfo, bool) {
	o := t.o.List(rackInfoDataNodeInfosOff).ObjectAt(i)
	if o.IsNull() {
		return DataNodeInfo{}, false
	}
	return DataNodeInfo{o: o}, true
}

// DiskInfosLen reports the number of diskInfos entries (proto field 3,
// map<string,DiskInfo>).
func (t RackInfo) DiskInfosLen() int { return t.o.List(rackInfoDiskInfosOff).Len() }

// DiskInfoAt returns the i-th diskInfos entry. The entry's Value is the raw
// DiskInfo sub-message bytes; decode it with WrapDiskInfo.
func (t RackInfo) DiskInfoAt(i int) StringMsgEntry {
	return stringMsgEntryAt(t.o.List(rackInfoDiskInfosOff), i)
}

// RackInfoInput collects the field values for NewRackInfo. DataNodeInfos takes
// each element as its own ZAP sub-buffer; each DiskInfos entry's Value is a
// DiskInfo message's own ZAP buffer.
type RackInfoInput struct {
	Id            string
	DataNodeInfos [][]byte
	DiskInfos     []StringMsgEntry
}

// NewRackInfo builds a ZAP-encoded RackInfo message from in and returns the
// bytes.
func NewRackInfo(in RackInfoInput) []byte {
	b := zap.NewBuilder(256)
	dataNodeInfosOff, dataNodeInfosLen := setMsgList(b, in.DataNodeInfos)
	diskInfosOff, diskInfosLen := setStringMsgMap(b, in.DiskInfos)
	ob := b.StartObject(rackInfoSize)
	ob.SetText(rackInfoIdOff, in.Id)
	ob.SetList(rackInfoDataNodeInfosOff, dataNodeInfosOff, dataNodeInfosLen)
	ob.SetList(rackInfoDiskInfosOff, diskInfosOff, diskInfosLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- DataCenterInfo ---

const (
	dataCenterInfoIdOff        = 0
	dataCenterInfoRackInfosOff = 8
	dataCenterInfoDiskInfosOff = 16
	dataCenterInfoSize         = 24
)

// DataCenterInfo is a zero-copy view into a ZAP-encoded DataCenterInfo message.
type DataCenterInfo struct{ o zap.Object }

// WrapDataCenterInfo parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapDataCenterInfo(b []byte) (DataCenterInfo, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return DataCenterInfo{}, err
	}
	return DataCenterInfo{o: m.Root()}, nil
}

// Id reads the id field (proto field 1, string).
func (t DataCenterInfo) Id() string { return t.o.Text(dataCenterInfoIdOff) }

// RackInfosLen reports the number of rack_infos elements (proto field 2,
// repeated RackInfo).
func (t DataCenterInfo) RackInfosLen() int { return t.o.List(dataCenterInfoRackInfosOff).Len() }

// RackInfoAt returns the i-th rack_infos element.
func (t DataCenterInfo) RackInfoAt(i int) (RackInfo, bool) {
	o := t.o.List(dataCenterInfoRackInfosOff).ObjectAt(i)
	if o.IsNull() {
		return RackInfo{}, false
	}
	return RackInfo{o: o}, true
}

// DiskInfosLen reports the number of diskInfos entries (proto field 3,
// map<string,DiskInfo>).
func (t DataCenterInfo) DiskInfosLen() int { return t.o.List(dataCenterInfoDiskInfosOff).Len() }

// DiskInfoAt returns the i-th diskInfos entry. The entry's Value is the raw
// DiskInfo sub-message bytes; decode it with WrapDiskInfo.
func (t DataCenterInfo) DiskInfoAt(i int) StringMsgEntry {
	return stringMsgEntryAt(t.o.List(dataCenterInfoDiskInfosOff), i)
}

// DataCenterInfoInput collects the field values for NewDataCenterInfo. RackInfos
// takes each element as its own ZAP sub-buffer; each DiskInfos entry's Value is
// a DiskInfo message's own ZAP buffer.
type DataCenterInfoInput struct {
	Id        string
	RackInfos [][]byte
	DiskInfos []StringMsgEntry
}

// NewDataCenterInfo builds a ZAP-encoded DataCenterInfo message from in and
// returns the bytes.
func NewDataCenterInfo(in DataCenterInfoInput) []byte {
	b := zap.NewBuilder(256)
	rackInfosOff, rackInfosLen := setMsgList(b, in.RackInfos)
	diskInfosOff, diskInfosLen := setStringMsgMap(b, in.DiskInfos)
	ob := b.StartObject(dataCenterInfoSize)
	ob.SetText(dataCenterInfoIdOff, in.Id)
	ob.SetList(dataCenterInfoRackInfosOff, rackInfosOff, rackInfosLen)
	ob.SetList(dataCenterInfoDiskInfosOff, diskInfosOff, diskInfosLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- TopologyInfo ---

const (
	topologyInfoIdOff              = 0
	topologyInfoDataCenterInfosOff = 8
	topologyInfoDiskInfosOff       = 16
	topologyInfoSize               = 24
)

// TopologyInfo is a zero-copy view into a ZAP-encoded TopologyInfo message.
type TopologyInfo struct{ o zap.Object }

// WrapTopologyInfo parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapTopologyInfo(b []byte) (TopologyInfo, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return TopologyInfo{}, err
	}
	return TopologyInfo{o: m.Root()}, nil
}

// Id reads the id field (proto field 1, string).
func (t TopologyInfo) Id() string { return t.o.Text(topologyInfoIdOff) }

// DataCenterInfosLen reports the number of data_center_infos elements (proto
// field 2, repeated DataCenterInfo).
func (t TopologyInfo) DataCenterInfosLen() int {
	return t.o.List(topologyInfoDataCenterInfosOff).Len()
}

// DataCenterInfoAt returns the i-th data_center_infos element.
func (t TopologyInfo) DataCenterInfoAt(i int) (DataCenterInfo, bool) {
	o := t.o.List(topologyInfoDataCenterInfosOff).ObjectAt(i)
	if o.IsNull() {
		return DataCenterInfo{}, false
	}
	return DataCenterInfo{o: o}, true
}

// DiskInfosLen reports the number of diskInfos entries (proto field 3,
// map<string,DiskInfo>).
func (t TopologyInfo) DiskInfosLen() int { return t.o.List(topologyInfoDiskInfosOff).Len() }

// DiskInfoAt returns the i-th diskInfos entry. The entry's Value is the raw
// DiskInfo sub-message bytes; decode it with WrapDiskInfo.
func (t TopologyInfo) DiskInfoAt(i int) StringMsgEntry {
	return stringMsgEntryAt(t.o.List(topologyInfoDiskInfosOff), i)
}

// TopologyInfoInput collects the field values for NewTopologyInfo.
// DataCenterInfos takes each element as its own ZAP sub-buffer; each DiskInfos
// entry's Value is a DiskInfo message's own ZAP buffer.
type TopologyInfoInput struct {
	Id              string
	DataCenterInfos [][]byte
	DiskInfos       []StringMsgEntry
}

// NewTopologyInfo builds a ZAP-encoded TopologyInfo message from in and returns
// the bytes.
func NewTopologyInfo(in TopologyInfoInput) []byte {
	b := zap.NewBuilder(256)
	dataCenterInfosOff, dataCenterInfosLen := setMsgList(b, in.DataCenterInfos)
	diskInfosOff, diskInfosLen := setStringMsgMap(b, in.DiskInfos)
	ob := b.StartObject(topologyInfoSize)
	ob.SetText(topologyInfoIdOff, in.Id)
	ob.SetList(topologyInfoDataCenterInfosOff, dataCenterInfosOff, dataCenterInfosLen)
	ob.SetList(topologyInfoDiskInfosOff, diskInfosOff, diskInfosLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeListRequest ---

const volumeListRequestSize = 0

// VolumeListRequest is a zero-copy view into a ZAP-encoded VolumeListRequest
// message. VolumeListRequest carries no fields.
type VolumeListRequest struct{ o zap.Object }

// WrapVolumeListRequest parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapVolumeListRequest(b []byte) (VolumeListRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeListRequest{}, err
	}
	return VolumeListRequest{o: m.Root()}, nil
}

// VolumeListRequestInput collects the field values for NewVolumeListRequest.
// VolumeListRequest has no fields.
type VolumeListRequestInput struct{}

// NewVolumeListRequest builds a ZAP-encoded VolumeListRequest message and
// returns the bytes.
func NewVolumeListRequest(in VolumeListRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeListRequestSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeListResponse ---

const (
	volumeListResponseTopologyInfoOff      = 0
	volumeListResponseVolumeSizeLimitMbOff = 8
	volumeListResponseSize                 = 16
)

// VolumeListResponse is a zero-copy view into a ZAP-encoded VolumeListResponse
// message.
type VolumeListResponse struct{ o zap.Object }

// WrapVolumeListResponse parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapVolumeListResponse(b []byte) (VolumeListResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeListResponse{}, err
	}
	return VolumeListResponse{o: m.Root()}, nil
}

// TopologyInfo reads the topology_info field (proto field 1, message
// TopologyInfo). The bool is false when the field is absent.
func (t VolumeListResponse) TopologyInfo() (TopologyInfo, bool) {
	b := t.o.Bytes(volumeListResponseTopologyInfoOff)
	if len(b) == 0 {
		return TopologyInfo{}, false
	}
	ti, err := WrapTopologyInfo(b)
	if err != nil {
		return TopologyInfo{}, false
	}
	return ti, true
}

// VolumeSizeLimitMb reads the volume_size_limit_mb field (proto field 2,
// uint64).
func (t VolumeListResponse) VolumeSizeLimitMb() uint64 {
	return t.o.Uint64(volumeListResponseVolumeSizeLimitMbOff)
}

// VolumeListResponseInput collects the field values for NewVolumeListResponse.
// TopologyInfo is the child message's own ZAP buffer (from NewTopologyInfo);
// nil encodes an absent topology_info.
type VolumeListResponseInput struct {
	TopologyInfo      []byte
	VolumeSizeLimitMb uint64
}

// NewVolumeListResponse builds a ZAP-encoded VolumeListResponse message from in
// and returns the bytes.
func NewVolumeListResponse(in VolumeListResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeListResponseSize)
	ob.SetBytes(volumeListResponseTopologyInfoOff, in.TopologyInfo)
	ob.SetUint64(volumeListResponseVolumeSizeLimitMbOff, in.VolumeSizeLimitMb)
	ob.FinishAsRoot()
	return b.Finish()
}
