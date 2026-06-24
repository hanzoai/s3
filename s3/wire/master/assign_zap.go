// Code generated from master.proto; DO NOT EDIT.

package masterwire

import (
	zap "github.com/zap-proto/go"
)

// --- Location ---

const (
	locationUrlOff        = 0
	locationPublicUrlOff  = 8
	locationGrpcPortOff   = 16
	locationDataCenterOff = 20
	locationSize          = 28
)

// Location is a zero-copy view into a ZAP-encoded Location message.
type Location struct{ o zap.Object }

// WrapLocation parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapLocation(b []byte) (Location, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return Location{}, err
	}
	return Location{o: m.Root()}, nil
}

// Url reads the url field (proto field 1, string).
func (t Location) Url() string { return t.o.Text(locationUrlOff) }

// PublicUrl reads the public_url field (proto field 2, string).
func (t Location) PublicUrl() string { return t.o.Text(locationPublicUrlOff) }

// GrpcPort reads the grpc_port field (proto field 3, uint32).
func (t Location) GrpcPort() uint32 { return t.o.Uint32(locationGrpcPortOff) }

// DataCenter reads the data_center field (proto field 4, string).
func (t Location) DataCenter() string { return t.o.Text(locationDataCenterOff) }

// LocationInput collects the field values for NewLocation.
type LocationInput struct {
	Url        string
	PublicUrl  string
	GrpcPort   uint32
	DataCenter string
}

// NewLocation builds a ZAP-encoded Location message from in and returns the
// bytes.
func NewLocation(in LocationInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(locationSize)
	ob.SetText(locationUrlOff, in.Url)
	ob.SetText(locationPublicUrlOff, in.PublicUrl)
	ob.SetUint32(locationGrpcPortOff, in.GrpcPort)
	ob.SetText(locationDataCenterOff, in.DataCenter)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- LookupVolumeRequest ---

const (
	lookupVolumeRequestVolumeOrFileIdsOff = 0
	lookupVolumeRequestCollectionOff      = 8
	lookupVolumeRequestSize               = 16
)

// LookupVolumeRequest is a zero-copy view into a ZAP-encoded LookupVolumeRequest
// message.
type LookupVolumeRequest struct{ o zap.Object }

// WrapLookupVolumeRequest parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapLookupVolumeRequest(b []byte) (LookupVolumeRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return LookupVolumeRequest{}, err
	}
	return LookupVolumeRequest{o: m.Root()}, nil
}

// VolumeOrFileIdsLen reports the number of volume_or_file_ids elements (proto
// field 1, repeated string).
func (t LookupVolumeRequest) VolumeOrFileIdsLen() int {
	return t.o.List(lookupVolumeRequestVolumeOrFileIdsOff).Len()
}

// VolumeOrFileIdAt returns the i-th volume_or_file_ids element.
func (t LookupVolumeRequest) VolumeOrFileIdAt(i int) string {
	return stringAt(t.o.List(lookupVolumeRequestVolumeOrFileIdsOff), i)
}

// Collection reads the collection field (proto field 2, string).
func (t LookupVolumeRequest) Collection() string {
	return t.o.Text(lookupVolumeRequestCollectionOff)
}

// LookupVolumeRequestInput collects the field values for
// NewLookupVolumeRequest.
type LookupVolumeRequestInput struct {
	VolumeOrFileIds []string
	Collection      string
}

// NewLookupVolumeRequest builds a ZAP-encoded LookupVolumeRequest message from
// in and returns the bytes.
func NewLookupVolumeRequest(in LookupVolumeRequestInput) []byte {
	b := zap.NewBuilder(256)
	idsOff, idsLen := setStringList(b, in.VolumeOrFileIds)
	ob := b.StartObject(lookupVolumeRequestSize)
	ob.SetList(lookupVolumeRequestVolumeOrFileIdsOff, idsOff, idsLen)
	ob.SetText(lookupVolumeRequestCollectionOff, in.Collection)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- LookupVolumeResponse.VolumeIdLocation ---

const (
	lookupVolumeResponseVolumeIdLocationVolumeOrFileIdOff = 0
	lookupVolumeResponseVolumeIdLocationLocationsOff      = 8
	lookupVolumeResponseVolumeIdLocationErrorOff          = 16
	lookupVolumeResponseVolumeIdLocationAuthOff           = 24
	lookupVolumeResponseVolumeIdLocationSize              = 32
)

// LookupVolumeResponseVolumeIdLocation is a zero-copy view into a ZAP-encoded
// LookupVolumeResponse.VolumeIdLocation message.
type LookupVolumeResponseVolumeIdLocation struct{ o zap.Object }

// WrapLookupVolumeResponseVolumeIdLocation parses b and returns a typed view.
// Returns an error if the wire-level checks (magic, version, size) fail.
func WrapLookupVolumeResponseVolumeIdLocation(b []byte) (LookupVolumeResponseVolumeIdLocation, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return LookupVolumeResponseVolumeIdLocation{}, err
	}
	return LookupVolumeResponseVolumeIdLocation{o: m.Root()}, nil
}

// VolumeOrFileId reads the volume_or_file_id field (proto field 1, string).
func (t LookupVolumeResponseVolumeIdLocation) VolumeOrFileId() string {
	return t.o.Text(lookupVolumeResponseVolumeIdLocationVolumeOrFileIdOff)
}

// LocationsLen reports the number of locations elements (proto field 2, repeated
// Location).
func (t LookupVolumeResponseVolumeIdLocation) LocationsLen() int {
	return t.o.List(lookupVolumeResponseVolumeIdLocationLocationsOff).Len()
}

// LocationAt returns the i-th locations element.
func (t LookupVolumeResponseVolumeIdLocation) LocationAt(i int) (Location, bool) {
	o := t.o.List(lookupVolumeResponseVolumeIdLocationLocationsOff).ObjectAt(i)
	if o.IsNull() {
		return Location{}, false
	}
	return Location{o: o}, true
}

// Error reads the error field (proto field 3, string).
func (t LookupVolumeResponseVolumeIdLocation) Error() string {
	return t.o.Text(lookupVolumeResponseVolumeIdLocationErrorOff)
}

// Auth reads the auth field (proto field 4, string).
func (t LookupVolumeResponseVolumeIdLocation) Auth() string {
	return t.o.Text(lookupVolumeResponseVolumeIdLocationAuthOff)
}

// LookupVolumeResponseVolumeIdLocationInput collects the field values for
// NewLookupVolumeResponseVolumeIdLocation. Locations takes each Location element
// as its own ZAP sub-buffer.
type LookupVolumeResponseVolumeIdLocationInput struct {
	VolumeOrFileId string
	Locations      [][]byte
	Error          string
	Auth           string
}

// NewLookupVolumeResponseVolumeIdLocation builds a ZAP-encoded
// LookupVolumeResponse.VolumeIdLocation message from in and returns the bytes.
func NewLookupVolumeResponseVolumeIdLocation(in LookupVolumeResponseVolumeIdLocationInput) []byte {
	b := zap.NewBuilder(256)
	locationsOff, locationsLen := setMsgList(b, in.Locations)
	ob := b.StartObject(lookupVolumeResponseVolumeIdLocationSize)
	ob.SetText(lookupVolumeResponseVolumeIdLocationVolumeOrFileIdOff, in.VolumeOrFileId)
	ob.SetList(lookupVolumeResponseVolumeIdLocationLocationsOff, locationsOff, locationsLen)
	ob.SetText(lookupVolumeResponseVolumeIdLocationErrorOff, in.Error)
	ob.SetText(lookupVolumeResponseVolumeIdLocationAuthOff, in.Auth)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- LookupVolumeResponse ---

const (
	lookupVolumeResponseVolumeIdLocationsOff = 0
	lookupVolumeResponseSize                 = 8
)

// LookupVolumeResponse is a zero-copy view into a ZAP-encoded
// LookupVolumeResponse message.
type LookupVolumeResponse struct{ o zap.Object }

// WrapLookupVolumeResponse parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapLookupVolumeResponse(b []byte) (LookupVolumeResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return LookupVolumeResponse{}, err
	}
	return LookupVolumeResponse{o: m.Root()}, nil
}

// VolumeIdLocationsLen reports the number of volume_id_locations elements (proto
// field 1, repeated LookupVolumeResponse.VolumeIdLocation).
func (t LookupVolumeResponse) VolumeIdLocationsLen() int {
	return t.o.List(lookupVolumeResponseVolumeIdLocationsOff).Len()
}

// VolumeIdLocationAt returns the i-th volume_id_locations element.
func (t LookupVolumeResponse) VolumeIdLocationAt(i int) (LookupVolumeResponseVolumeIdLocation, bool) {
	o := t.o.List(lookupVolumeResponseVolumeIdLocationsOff).ObjectAt(i)
	if o.IsNull() {
		return LookupVolumeResponseVolumeIdLocation{}, false
	}
	return LookupVolumeResponseVolumeIdLocation{o: o}, true
}

// LookupVolumeResponseInput collects the field values for
// NewLookupVolumeResponse. VolumeIdLocations takes each element as its own ZAP
// sub-buffer (from NewLookupVolumeResponseVolumeIdLocation).
type LookupVolumeResponseInput struct {
	VolumeIdLocations [][]byte
}

// NewLookupVolumeResponse builds a ZAP-encoded LookupVolumeResponse message from
// in and returns the bytes.
func NewLookupVolumeResponse(in LookupVolumeResponseInput) []byte {
	b := zap.NewBuilder(256)
	locsOff, locsLen := setMsgList(b, in.VolumeIdLocations)
	ob := b.StartObject(lookupVolumeResponseSize)
	ob.SetList(lookupVolumeResponseVolumeIdLocationsOff, locsOff, locsLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- AssignRequest ---

const (
	assignRequestCountOff               = 0
	assignRequestReplicationOff         = 8
	assignRequestCollectionOff          = 16
	assignRequestTtlOff                 = 24
	assignRequestDataCenterOff          = 32
	assignRequestRackOff                = 40
	assignRequestDataNodeOff            = 48
	assignRequestMemoryMapMaxSizeMbOff  = 56
	assignRequestWritableVolumeCountOff = 60
	assignRequestDiskTypeOff            = 64
	assignRequestExpectedDataSizeOff    = 72
	assignRequestSize                   = 80
)

// AssignRequest is a zero-copy view into a ZAP-encoded AssignRequest message.
type AssignRequest struct{ o zap.Object }

// WrapAssignRequest parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapAssignRequest(b []byte) (AssignRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return AssignRequest{}, err
	}
	return AssignRequest{o: m.Root()}, nil
}

// Count reads the count field (proto field 1, uint64).
func (t AssignRequest) Count() uint64 { return t.o.Uint64(assignRequestCountOff) }

// Replication reads the replication field (proto field 2, string).
func (t AssignRequest) Replication() string { return t.o.Text(assignRequestReplicationOff) }

// Collection reads the collection field (proto field 3, string).
func (t AssignRequest) Collection() string { return t.o.Text(assignRequestCollectionOff) }

// Ttl reads the ttl field (proto field 4, string).
func (t AssignRequest) Ttl() string { return t.o.Text(assignRequestTtlOff) }

// DataCenter reads the data_center field (proto field 5, string).
func (t AssignRequest) DataCenter() string { return t.o.Text(assignRequestDataCenterOff) }

// Rack reads the rack field (proto field 6, string).
func (t AssignRequest) Rack() string { return t.o.Text(assignRequestRackOff) }

// DataNode reads the data_node field (proto field 7, string).
func (t AssignRequest) DataNode() string { return t.o.Text(assignRequestDataNodeOff) }

// MemoryMapMaxSizeMb reads the memory_map_max_size_mb field (proto field 8,
// uint32).
func (t AssignRequest) MemoryMapMaxSizeMb() uint32 {
	return t.o.Uint32(assignRequestMemoryMapMaxSizeMbOff)
}

// WritableVolumeCount reads the writable_volume_count field (proto field 9,
// uint32).
func (t AssignRequest) WritableVolumeCount() uint32 {
	return t.o.Uint32(assignRequestWritableVolumeCountOff)
}

// DiskType reads the disk_type field (proto field 10, string).
func (t AssignRequest) DiskType() string { return t.o.Text(assignRequestDiskTypeOff) }

// ExpectedDataSize reads the expected_data_size field (proto field 11, uint64).
func (t AssignRequest) ExpectedDataSize() uint64 {
	return t.o.Uint64(assignRequestExpectedDataSizeOff)
}

// AssignRequestInput collects the field values for NewAssignRequest.
type AssignRequestInput struct {
	Count               uint64
	Replication         string
	Collection          string
	Ttl                 string
	DataCenter          string
	Rack                string
	DataNode            string
	MemoryMapMaxSizeMb  uint32
	WritableVolumeCount uint32
	DiskType            string
	ExpectedDataSize    uint64
}

// NewAssignRequest builds a ZAP-encoded AssignRequest message from in and
// returns the bytes.
func NewAssignRequest(in AssignRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(assignRequestSize)
	ob.SetUint64(assignRequestCountOff, in.Count)
	ob.SetText(assignRequestReplicationOff, in.Replication)
	ob.SetText(assignRequestCollectionOff, in.Collection)
	ob.SetText(assignRequestTtlOff, in.Ttl)
	ob.SetText(assignRequestDataCenterOff, in.DataCenter)
	ob.SetText(assignRequestRackOff, in.Rack)
	ob.SetText(assignRequestDataNodeOff, in.DataNode)
	ob.SetUint32(assignRequestMemoryMapMaxSizeMbOff, in.MemoryMapMaxSizeMb)
	ob.SetUint32(assignRequestWritableVolumeCountOff, in.WritableVolumeCount)
	ob.SetText(assignRequestDiskTypeOff, in.DiskType)
	ob.SetUint64(assignRequestExpectedDataSizeOff, in.ExpectedDataSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- VolumeGrowRequest ---

const (
	volumeGrowRequestWritableVolumeCountOff = 0
	volumeGrowRequestReplicationOff         = 4
	volumeGrowRequestCollectionOff          = 12
	volumeGrowRequestTtlOff                 = 20
	volumeGrowRequestDataCenterOff          = 28
	volumeGrowRequestRackOff                = 36
	volumeGrowRequestDataNodeOff            = 44
	volumeGrowRequestMemoryMapMaxSizeMbOff  = 52
	volumeGrowRequestDiskTypeOff            = 56
	volumeGrowRequestSize                   = 64
)

// VolumeGrowRequest is a zero-copy view into a ZAP-encoded VolumeGrowRequest
// message.
type VolumeGrowRequest struct{ o zap.Object }

// WrapVolumeGrowRequest parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapVolumeGrowRequest(b []byte) (VolumeGrowRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeGrowRequest{}, err
	}
	return VolumeGrowRequest{o: m.Root()}, nil
}

// WritableVolumeCount reads the writable_volume_count field (proto field 1,
// uint32).
func (t VolumeGrowRequest) WritableVolumeCount() uint32 {
	return t.o.Uint32(volumeGrowRequestWritableVolumeCountOff)
}

// Replication reads the replication field (proto field 2, string).
func (t VolumeGrowRequest) Replication() string { return t.o.Text(volumeGrowRequestReplicationOff) }

// Collection reads the collection field (proto field 3, string).
func (t VolumeGrowRequest) Collection() string { return t.o.Text(volumeGrowRequestCollectionOff) }

// Ttl reads the ttl field (proto field 4, string).
func (t VolumeGrowRequest) Ttl() string { return t.o.Text(volumeGrowRequestTtlOff) }

// DataCenter reads the data_center field (proto field 5, string).
func (t VolumeGrowRequest) DataCenter() string { return t.o.Text(volumeGrowRequestDataCenterOff) }

// Rack reads the rack field (proto field 6, string).
func (t VolumeGrowRequest) Rack() string { return t.o.Text(volumeGrowRequestRackOff) }

// DataNode reads the data_node field (proto field 7, string).
func (t VolumeGrowRequest) DataNode() string { return t.o.Text(volumeGrowRequestDataNodeOff) }

// MemoryMapMaxSizeMb reads the memory_map_max_size_mb field (proto field 8,
// uint32).
func (t VolumeGrowRequest) MemoryMapMaxSizeMb() uint32 {
	return t.o.Uint32(volumeGrowRequestMemoryMapMaxSizeMbOff)
}

// DiskType reads the disk_type field (proto field 9, string).
func (t VolumeGrowRequest) DiskType() string { return t.o.Text(volumeGrowRequestDiskTypeOff) }

// VolumeGrowRequestInput collects the field values for NewVolumeGrowRequest.
type VolumeGrowRequestInput struct {
	WritableVolumeCount uint32
	Replication         string
	Collection          string
	Ttl                 string
	DataCenter          string
	Rack                string
	DataNode            string
	MemoryMapMaxSizeMb  uint32
	DiskType            string
}

// NewVolumeGrowRequest builds a ZAP-encoded VolumeGrowRequest message from in
// and returns the bytes.
func NewVolumeGrowRequest(in VolumeGrowRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeGrowRequestSize)
	ob.SetUint32(volumeGrowRequestWritableVolumeCountOff, in.WritableVolumeCount)
	ob.SetText(volumeGrowRequestReplicationOff, in.Replication)
	ob.SetText(volumeGrowRequestCollectionOff, in.Collection)
	ob.SetText(volumeGrowRequestTtlOff, in.Ttl)
	ob.SetText(volumeGrowRequestDataCenterOff, in.DataCenter)
	ob.SetText(volumeGrowRequestRackOff, in.Rack)
	ob.SetText(volumeGrowRequestDataNodeOff, in.DataNode)
	ob.SetUint32(volumeGrowRequestMemoryMapMaxSizeMbOff, in.MemoryMapMaxSizeMb)
	ob.SetText(volumeGrowRequestDiskTypeOff, in.DiskType)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- AssignResponse ---

const (
	assignResponseFidOff      = 0
	assignResponseCountOff    = 8
	assignResponseErrorOff    = 16
	assignResponseAuthOff     = 24
	assignResponseReplicasOff = 32
	assignResponseLocationOff = 40
	assignResponseSize        = 48
)

// AssignResponse is a zero-copy view into a ZAP-encoded AssignResponse message.
type AssignResponse struct{ o zap.Object }

// WrapAssignResponse parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapAssignResponse(b []byte) (AssignResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return AssignResponse{}, err
	}
	return AssignResponse{o: m.Root()}, nil
}

// Fid reads the fid field (proto field 1, string).
func (t AssignResponse) Fid() string { return t.o.Text(assignResponseFidOff) }

// Count reads the count field (proto field 4, uint64).
func (t AssignResponse) Count() uint64 { return t.o.Uint64(assignResponseCountOff) }

// Error reads the error field (proto field 5, string).
func (t AssignResponse) Error() string { return t.o.Text(assignResponseErrorOff) }

// Auth reads the auth field (proto field 6, string).
func (t AssignResponse) Auth() string { return t.o.Text(assignResponseAuthOff) }

// ReplicasLen reports the number of replicas elements (proto field 7, repeated
// Location).
func (t AssignResponse) ReplicasLen() int { return t.o.List(assignResponseReplicasOff).Len() }

// ReplicaAt returns the i-th replicas element.
func (t AssignResponse) ReplicaAt(i int) (Location, bool) {
	o := t.o.List(assignResponseReplicasOff).ObjectAt(i)
	if o.IsNull() {
		return Location{}, false
	}
	return Location{o: o}, true
}

// Location reads the location field (proto field 8, message Location). The bool
// is false when the field is absent.
func (t AssignResponse) Location() (Location, bool) {
	b := t.o.Bytes(assignResponseLocationOff)
	if len(b) == 0 {
		return Location{}, false
	}
	l, err := WrapLocation(b)
	if err != nil {
		return Location{}, false
	}
	return l, true
}

// AssignResponseInput collects the field values for NewAssignResponse. Replicas
// takes each Location element as its own ZAP sub-buffer; Location is a single
// Location message's own ZAP buffer (from NewLocation), nil for absent.
type AssignResponseInput struct {
	Fid      string
	Count    uint64
	Error    string
	Auth     string
	Replicas [][]byte
	Location []byte
}

// NewAssignResponse builds a ZAP-encoded AssignResponse message from in and
// returns the bytes.
func NewAssignResponse(in AssignResponseInput) []byte {
	b := zap.NewBuilder(256)
	replicasOff, replicasLen := setMsgList(b, in.Replicas)
	ob := b.StartObject(assignResponseSize)
	ob.SetText(assignResponseFidOff, in.Fid)
	ob.SetUint64(assignResponseCountOff, in.Count)
	ob.SetText(assignResponseErrorOff, in.Error)
	ob.SetText(assignResponseAuthOff, in.Auth)
	ob.SetList(assignResponseReplicasOff, replicasOff, replicasLen)
	ob.SetBytes(assignResponseLocationOff, in.Location)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- StatisticsRequest ---

const (
	statisticsRequestReplicationOff = 0
	statisticsRequestCollectionOff  = 8
	statisticsRequestTtlOff         = 16
	statisticsRequestDiskTypeOff    = 24
	statisticsRequestSize           = 32
)

// StatisticsRequest is a zero-copy view into a ZAP-encoded StatisticsRequest
// message.
type StatisticsRequest struct{ o zap.Object }

// WrapStatisticsRequest parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapStatisticsRequest(b []byte) (StatisticsRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return StatisticsRequest{}, err
	}
	return StatisticsRequest{o: m.Root()}, nil
}

// Replication reads the replication field (proto field 1, string).
func (t StatisticsRequest) Replication() string { return t.o.Text(statisticsRequestReplicationOff) }

// Collection reads the collection field (proto field 2, string).
func (t StatisticsRequest) Collection() string { return t.o.Text(statisticsRequestCollectionOff) }

// Ttl reads the ttl field (proto field 3, string).
func (t StatisticsRequest) Ttl() string { return t.o.Text(statisticsRequestTtlOff) }

// DiskType reads the disk_type field (proto field 4, string).
func (t StatisticsRequest) DiskType() string { return t.o.Text(statisticsRequestDiskTypeOff) }

// StatisticsRequestInput collects the field values for NewStatisticsRequest.
type StatisticsRequestInput struct {
	Replication string
	Collection  string
	Ttl         string
	DiskType    string
}

// NewStatisticsRequest builds a ZAP-encoded StatisticsRequest message from in
// and returns the bytes.
func NewStatisticsRequest(in StatisticsRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(statisticsRequestSize)
	ob.SetText(statisticsRequestReplicationOff, in.Replication)
	ob.SetText(statisticsRequestCollectionOff, in.Collection)
	ob.SetText(statisticsRequestTtlOff, in.Ttl)
	ob.SetText(statisticsRequestDiskTypeOff, in.DiskType)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- StatisticsResponse ---

const (
	statisticsResponseTotalSizeOff = 0
	statisticsResponseUsedSizeOff  = 8
	statisticsResponseFileCountOff = 16
	statisticsResponseSize         = 24
)

// StatisticsResponse is a zero-copy view into a ZAP-encoded StatisticsResponse
// message.
type StatisticsResponse struct{ o zap.Object }

// WrapStatisticsResponse parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapStatisticsResponse(b []byte) (StatisticsResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return StatisticsResponse{}, err
	}
	return StatisticsResponse{o: m.Root()}, nil
}

// TotalSize reads the total_size field (proto field 4, uint64).
func (t StatisticsResponse) TotalSize() uint64 { return t.o.Uint64(statisticsResponseTotalSizeOff) }

// UsedSize reads the used_size field (proto field 5, uint64).
func (t StatisticsResponse) UsedSize() uint64 { return t.o.Uint64(statisticsResponseUsedSizeOff) }

// FileCount reads the file_count field (proto field 6, uint64).
func (t StatisticsResponse) FileCount() uint64 { return t.o.Uint64(statisticsResponseFileCountOff) }

// StatisticsResponseInput collects the field values for NewStatisticsResponse.
type StatisticsResponseInput struct {
	TotalSize uint64
	UsedSize  uint64
	FileCount uint64
}

// NewStatisticsResponse builds a ZAP-encoded StatisticsResponse message from in
// and returns the bytes.
func NewStatisticsResponse(in StatisticsResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(statisticsResponseSize)
	ob.SetUint64(statisticsResponseTotalSizeOff, in.TotalSize)
	ob.SetUint64(statisticsResponseUsedSizeOff, in.UsedSize)
	ob.SetUint64(statisticsResponseFileCountOff, in.FileCount)
	ob.FinishAsRoot()
	return b.Finish()
}
