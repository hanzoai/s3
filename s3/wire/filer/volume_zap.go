// Code generated from filer.proto; DO NOT EDIT.

package filerwire

import (
	zap "github.com/zap-proto/go"
)

// --- AssignVolumeRequest ---

const (
	assignVolumeRequestCountOff            = 0  // count (field 1, int32)
	assignVolumeRequestCollectionOff       = 8  // collection (field 2, string)
	assignVolumeRequestReplicationOff      = 16 // replication (field 3, string)
	assignVolumeRequestTtlSecOff           = 24 // ttl_sec (field 4, int32)
	assignVolumeRequestDataCenterOff       = 32 // data_center (field 5, string)
	assignVolumeRequestPathOff             = 40 // path (field 6, string)
	assignVolumeRequestRackOff             = 48 // rack (field 7, string)
	assignVolumeRequestDataNodeOff         = 56 // data_node (field 9, string)
	assignVolumeRequestDiskTypeOff         = 64 // disk_type (field 8, string)
	assignVolumeRequestExpectedDataSizeOff = 72 // expected_data_size (field 10, uint64)
	assignVolumeRequestSize                = 80
)

// AssignVolumeRequest is a zero-copy view into a ZAP-encoded AssignVolumeRequest message.
type AssignVolumeRequest struct{ o zap.Object }

// WrapAssignVolumeRequest parses b and returns a typed view.
func WrapAssignVolumeRequest(b []byte) (AssignVolumeRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return AssignVolumeRequest{}, err
	}
	return AssignVolumeRequest{o: m.Root()}, nil
}

func (t AssignVolumeRequest) Count() int32        { return t.o.Int32(assignVolumeRequestCountOff) }
func (t AssignVolumeRequest) Collection() string  { return t.o.Text(assignVolumeRequestCollectionOff) }
func (t AssignVolumeRequest) Replication() string { return t.o.Text(assignVolumeRequestReplicationOff) }
func (t AssignVolumeRequest) TtlSec() int32       { return t.o.Int32(assignVolumeRequestTtlSecOff) }
func (t AssignVolumeRequest) DataCenter() string  { return t.o.Text(assignVolumeRequestDataCenterOff) }
func (t AssignVolumeRequest) Path() string        { return t.o.Text(assignVolumeRequestPathOff) }
func (t AssignVolumeRequest) Rack() string        { return t.o.Text(assignVolumeRequestRackOff) }
func (t AssignVolumeRequest) DataNode() string    { return t.o.Text(assignVolumeRequestDataNodeOff) }
func (t AssignVolumeRequest) DiskType() string    { return t.o.Text(assignVolumeRequestDiskTypeOff) }
func (t AssignVolumeRequest) ExpectedDataSize() uint64 {
	return t.o.Uint64(assignVolumeRequestExpectedDataSizeOff)
}

// AssignVolumeRequestInput collects the field values for NewAssignVolumeRequest.
type AssignVolumeRequestInput struct {
	Count            int32
	Collection       string
	Replication      string
	TtlSec           int32
	DataCenter       string
	Path             string
	Rack             string
	DataNode         string
	DiskType         string
	ExpectedDataSize uint64
}

// NewAssignVolumeRequest builds a ZAP-encoded AssignVolumeRequest message from in.
func NewAssignVolumeRequest(in AssignVolumeRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(assignVolumeRequestSize)
	ob.SetInt32(assignVolumeRequestCountOff, in.Count)
	ob.SetText(assignVolumeRequestCollectionOff, in.Collection)
	ob.SetText(assignVolumeRequestReplicationOff, in.Replication)
	ob.SetInt32(assignVolumeRequestTtlSecOff, in.TtlSec)
	ob.SetText(assignVolumeRequestDataCenterOff, in.DataCenter)
	ob.SetText(assignVolumeRequestPathOff, in.Path)
	ob.SetText(assignVolumeRequestRackOff, in.Rack)
	ob.SetText(assignVolumeRequestDataNodeOff, in.DataNode)
	ob.SetText(assignVolumeRequestDiskTypeOff, in.DiskType)
	ob.SetUint64(assignVolumeRequestExpectedDataSizeOff, in.ExpectedDataSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- Location ---

const (
	locationURLOff        = 0  // url (field 1, string)
	locationPublicURLOff  = 8  // public_url (field 2, string)
	locationGrpcPortOff   = 16 // grpc_port (field 3, uint32)
	locationDataCenterOff = 24 // data_center (field 4, string)
	locationSize          = 32
)

// Location is a zero-copy view into a ZAP-encoded Location message.
type Location struct{ o zap.Object }

// WrapLocation parses b and returns a typed view.
func WrapLocation(b []byte) (Location, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return Location{}, err
	}
	return Location{o: m.Root()}, nil
}

func (t Location) URL() string        { return t.o.Text(locationURLOff) }
func (t Location) PublicURL() string  { return t.o.Text(locationPublicURLOff) }
func (t Location) GrpcPort() uint32   { return t.o.Uint32(locationGrpcPortOff) }
func (t Location) DataCenter() string { return t.o.Text(locationDataCenterOff) }

// LocationInput collects the field values for NewLocation.
type LocationInput struct {
	URL        string
	PublicURL  string
	GrpcPort   uint32
	DataCenter string
}

// NewLocation builds a ZAP-encoded Location message from in.
func NewLocation(in LocationInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(locationSize)
	ob.SetText(locationURLOff, in.URL)
	ob.SetText(locationPublicURLOff, in.PublicURL)
	ob.SetUint32(locationGrpcPortOff, in.GrpcPort)
	ob.SetText(locationDataCenterOff, in.DataCenter)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- Locations ---

const (
	locationsLocationsOff = 0 // locations (field 1, repeated Location)
	locationsSize         = 8
)

// Locations is a zero-copy view into a ZAP-encoded Locations message.
type Locations struct{ o zap.Object }

// WrapLocations parses b and returns a typed view.
func WrapLocations(b []byte) (Locations, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return Locations{}, err
	}
	return Locations{o: m.Root()}, nil
}

// LocationsLen returns the number of locations elements (proto field 1, repeated Location).
func (t Locations) LocationsLen() int { return t.o.List(locationsLocationsOff).Len() }

// Locations returns the i-th locations element as a typed Location view.
func (t Locations) Locations(i int) Location {
	return Location{o: t.o.List(locationsLocationsOff).ObjectAt(i)}
}

// LocationsInput collects the field values for NewLocations. Each Locations
// element is a standalone Location buffer (from NewLocation).
type LocationsInput struct {
	Locations [][]byte
}

// NewLocations builds a ZAP-encoded Locations message from in.
func NewLocations(in LocationsInput) []byte {
	b := zap.NewBuilder(256)
	locOff, locLen := addObjectList(b, in.Locations)
	ob := b.StartObject(locationsSize)
	ob.SetList(locationsLocationsOff, locOff, locLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- AssignVolumeResponse ---

const (
	assignVolumeResponseFileIDOff      = 0  // file_id (field 1, string)
	assignVolumeResponseCountOff       = 8  // count (field 4, int32)
	assignVolumeResponseAuthOff        = 16 // auth (field 5, string)
	assignVolumeResponseCollectionOff  = 24 // collection (field 6, string)
	assignVolumeResponseReplicationOff = 32 // replication (field 7, string)
	assignVolumeResponseErrorOff       = 40 // error (field 8, string)
	assignVolumeResponseLocationOff    = 48 // location (field 9, Location object)
	assignVolumeResponseSize           = 56
)

// AssignVolumeResponse is a zero-copy view into a ZAP-encoded AssignVolumeResponse message.
type AssignVolumeResponse struct{ o zap.Object }

// WrapAssignVolumeResponse parses b and returns a typed view.
func WrapAssignVolumeResponse(b []byte) (AssignVolumeResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return AssignVolumeResponse{}, err
	}
	return AssignVolumeResponse{o: m.Root()}, nil
}

func (t AssignVolumeResponse) FileID() string     { return t.o.Text(assignVolumeResponseFileIDOff) }
func (t AssignVolumeResponse) Count() int32       { return t.o.Int32(assignVolumeResponseCountOff) }
func (t AssignVolumeResponse) Auth() string       { return t.o.Text(assignVolumeResponseAuthOff) }
func (t AssignVolumeResponse) Collection() string { return t.o.Text(assignVolumeResponseCollectionOff) }
func (t AssignVolumeResponse) Replication() string {
	return t.o.Text(assignVolumeResponseReplicationOff)
}
func (t AssignVolumeResponse) Error() string { return t.o.Text(assignVolumeResponseErrorOff) }

// Location returns the nested Location (null-object if unset).
func (t AssignVolumeResponse) Location() Location {
	return Location{o: t.o.Object(assignVolumeResponseLocationOff)}
}

// HasLocation reports whether the location field is present.
func (t AssignVolumeResponse) HasLocation() bool {
	return !t.o.Object(assignVolumeResponseLocationOff).IsNull()
}

// AssignVolumeResponseInput collects the field values for NewAssignVolumeResponse.
// Location, when non-nil, is a standalone Location buffer (from NewLocation).
type AssignVolumeResponseInput struct {
	FileID      string
	Count       int32
	Auth        string
	Collection  string
	Replication string
	Error       string
	Location    []byte
}

// NewAssignVolumeResponse builds a ZAP-encoded AssignVolumeResponse message from in.
func NewAssignVolumeResponse(in AssignVolumeResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(assignVolumeResponseSize)
	ob.SetText(assignVolumeResponseFileIDOff, in.FileID)
	ob.SetInt32(assignVolumeResponseCountOff, in.Count)
	ob.SetText(assignVolumeResponseAuthOff, in.Auth)
	ob.SetText(assignVolumeResponseCollectionOff, in.Collection)
	ob.SetText(assignVolumeResponseReplicationOff, in.Replication)
	ob.SetText(assignVolumeResponseErrorOff, in.Error)
	setNestedObject(b, ob, assignVolumeResponseLocationOff, in.Location)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- LookupVolumeRequest ---

const (
	lookupVolumeRequestVolumeIDsOff = 0 // volume_ids (field 1, repeated string)
	lookupVolumeRequestSize         = 8
)

// LookupVolumeRequest is a zero-copy view into a ZAP-encoded LookupVolumeRequest message.
type LookupVolumeRequest struct{ o zap.Object }

// WrapLookupVolumeRequest parses b and returns a typed view.
func WrapLookupVolumeRequest(b []byte) (LookupVolumeRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return LookupVolumeRequest{}, err
	}
	return LookupVolumeRequest{o: m.Root()}, nil
}

// VolumeIDsLen returns the number of volume_ids elements (proto field 1, repeated string).
func (t LookupVolumeRequest) VolumeIDsLen() int {
	return t.o.List(lookupVolumeRequestVolumeIDsOff).Len()
}

// VolumeID returns the i-th volume_ids element.
func (t LookupVolumeRequest) VolumeID(i int) string {
	return elemText(t.o.List(lookupVolumeRequestVolumeIDsOff), i)
}

// LookupVolumeRequestInput collects the field values for NewLookupVolumeRequest.
type LookupVolumeRequestInput struct {
	VolumeIDs []string
}

// NewLookupVolumeRequest builds a ZAP-encoded LookupVolumeRequest message from in.
func NewLookupVolumeRequest(in LookupVolumeRequestInput) []byte {
	b := zap.NewBuilder(256)
	idsOff, idsLen := addStringList(b, in.VolumeIDs)
	ob := b.StartObject(lookupVolumeRequestSize)
	ob.SetList(lookupVolumeRequestVolumeIDsOff, idsOff, idsLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- LookupVolumeResponseLocationsMapEntry (map<string, Locations>) ---

const (
	lookupVolumeResponseLocationsMapEntryKeyOff   = 0
	lookupVolumeResponseLocationsMapEntryValueOff = 8
	lookupVolumeResponseLocationsMapEntrySize     = 12
)

// LookupVolumeResponseLocationsMapEntry is a zero-copy view of one
// map<string,Locations> pair of LookupVolumeResponse.locations_map (proto field 1).
type LookupVolumeResponseLocationsMapEntry struct{ o zap.Object }

// WrapLookupVolumeResponseLocationsMapEntry parses b and returns a typed view.
func WrapLookupVolumeResponseLocationsMapEntry(b []byte) (LookupVolumeResponseLocationsMapEntry, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return LookupVolumeResponseLocationsMapEntry{}, err
	}
	return LookupVolumeResponseLocationsMapEntry{o: m.Root()}, nil
}

func (t LookupVolumeResponseLocationsMapEntry) Key() string {
	return t.o.Text(lookupVolumeResponseLocationsMapEntryKeyOff)
}

// Value returns the nested Locations value (null-object if unset).
func (t LookupVolumeResponseLocationsMapEntry) Value() Locations {
	return Locations{o: t.o.Object(lookupVolumeResponseLocationsMapEntryValueOff)}
}

// LookupVolumeResponseLocationsMapEntryInput collects one map<string,Locations>
// pair. Value is a pre-built Locations sub-buffer (from NewLocations).
type LookupVolumeResponseLocationsMapEntryInput struct {
	Key   string
	Value []byte
}

// NewLookupVolumeResponseLocationsMapEntry builds one encoded locations_map entry.
func NewLookupVolumeResponseLocationsMapEntry(in LookupVolumeResponseLocationsMapEntryInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(lookupVolumeResponseLocationsMapEntrySize)
	ob.SetText(lookupVolumeResponseLocationsMapEntryKeyOff, in.Key)
	setNestedObject(b, ob, lookupVolumeResponseLocationsMapEntryValueOff, in.Value)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- LookupVolumeResponse ---

const (
	lookupVolumeResponseLocationsMapOff = 0 // locations_map (field 1, map<string,Locations>)
	lookupVolumeResponseSize            = 8
)

// LookupVolumeResponse is a zero-copy view into a ZAP-encoded LookupVolumeResponse message.
type LookupVolumeResponse struct{ o zap.Object }

// WrapLookupVolumeResponse parses b and returns a typed view.
func WrapLookupVolumeResponse(b []byte) (LookupVolumeResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return LookupVolumeResponse{}, err
	}
	return LookupVolumeResponse{o: m.Root()}, nil
}

// LocationsMapLen returns the number of locations_map entries (proto field 1,
// map<string,Locations>).
func (t LookupVolumeResponse) LocationsMapLen() int {
	return t.o.List(lookupVolumeResponseLocationsMapOff).Len()
}

// LocationsMap returns the i-th locations_map entry.
func (t LookupVolumeResponse) LocationsMap(i int) LookupVolumeResponseLocationsMapEntry {
	return LookupVolumeResponseLocationsMapEntry{o: t.o.List(lookupVolumeResponseLocationsMapOff).ObjectAt(i)}
}

// LookupVolumeResponseInput collects the field values for NewLookupVolumeResponse.
// Each LocationsMap element is a standalone LookupVolumeResponseLocationsMapEntry
// buffer.
type LookupVolumeResponseInput struct {
	LocationsMap [][]byte
}

// NewLookupVolumeResponse builds a ZAP-encoded LookupVolumeResponse message from in.
func NewLookupVolumeResponse(in LookupVolumeResponseInput) []byte {
	b := zap.NewBuilder(256)
	mapOff, mapLen := addObjectList(b, in.LocationsMap)
	ob := b.StartObject(lookupVolumeResponseSize)
	ob.SetList(lookupVolumeResponseLocationsMapOff, mapOff, mapLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- Collection ---

const (
	collectionNameOff = 0 // name (field 1, string)
	collectionSize    = 8
)

// Collection is a zero-copy view into a ZAP-encoded Collection message.
type Collection struct{ o zap.Object }

// WrapCollection parses b and returns a typed view.
func WrapCollection(b []byte) (Collection, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return Collection{}, err
	}
	return Collection{o: m.Root()}, nil
}

func (t Collection) Name() string { return t.o.Text(collectionNameOff) }

// CollectionInput collects the field values for NewCollection.
type CollectionInput struct {
	Name string
}

// NewCollection builds a ZAP-encoded Collection message from in.
func NewCollection(in CollectionInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(collectionSize)
	ob.SetText(collectionNameOff, in.Name)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- CollectionListRequest ---

const (
	collectionListRequestIncludeNormalVolumesOff = 0 // include_normal_volumes (field 1, bool)
	collectionListRequestIncludeEcVolumesOff     = 1 // include_ec_volumes (field 2, bool)
	collectionListRequestSize                    = 2
)

// CollectionListRequest is a zero-copy view into a ZAP-encoded CollectionListRequest message.
type CollectionListRequest struct{ o zap.Object }

// WrapCollectionListRequest parses b and returns a typed view.
func WrapCollectionListRequest(b []byte) (CollectionListRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return CollectionListRequest{}, err
	}
	return CollectionListRequest{o: m.Root()}, nil
}

func (t CollectionListRequest) IncludeNormalVolumes() bool {
	return t.o.Bool(collectionListRequestIncludeNormalVolumesOff)
}
func (t CollectionListRequest) IncludeEcVolumes() bool {
	return t.o.Bool(collectionListRequestIncludeEcVolumesOff)
}

// CollectionListRequestInput collects the field values for NewCollectionListRequest.
type CollectionListRequestInput struct {
	IncludeNormalVolumes bool
	IncludeEcVolumes     bool
}

// NewCollectionListRequest builds a ZAP-encoded CollectionListRequest message from in.
func NewCollectionListRequest(in CollectionListRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(collectionListRequestSize)
	ob.SetBool(collectionListRequestIncludeNormalVolumesOff, in.IncludeNormalVolumes)
	ob.SetBool(collectionListRequestIncludeEcVolumesOff, in.IncludeEcVolumes)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- CollectionListResponse ---

const (
	collectionListResponseCollectionsOff = 0 // collections (field 1, repeated Collection)
	collectionListResponseSize           = 8
)

// CollectionListResponse is a zero-copy view into a ZAP-encoded CollectionListResponse message.
type CollectionListResponse struct{ o zap.Object }

// WrapCollectionListResponse parses b and returns a typed view.
func WrapCollectionListResponse(b []byte) (CollectionListResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return CollectionListResponse{}, err
	}
	return CollectionListResponse{o: m.Root()}, nil
}

// CollectionsLen returns the number of collections elements (proto field 1,
// repeated Collection).
func (t CollectionListResponse) CollectionsLen() int {
	return t.o.List(collectionListResponseCollectionsOff).Len()
}

// Collections returns the i-th collections element as a typed Collection view.
func (t CollectionListResponse) Collections(i int) Collection {
	return Collection{o: t.o.List(collectionListResponseCollectionsOff).ObjectAt(i)}
}

// CollectionListResponseInput collects the field values for NewCollectionListResponse.
// Each Collections element is a standalone Collection buffer (from NewCollection).
type CollectionListResponseInput struct {
	Collections [][]byte
}

// NewCollectionListResponse builds a ZAP-encoded CollectionListResponse message from in.
func NewCollectionListResponse(in CollectionListResponseInput) []byte {
	b := zap.NewBuilder(256)
	colOff, colLen := addObjectList(b, in.Collections)
	ob := b.StartObject(collectionListResponseSize)
	ob.SetList(collectionListResponseCollectionsOff, colOff, colLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- DeleteCollectionRequest ---

const (
	deleteCollectionRequestCollectionOff = 0 // collection (field 1, string)
	deleteCollectionRequestSize          = 8
)

// DeleteCollectionRequest is a zero-copy view into a ZAP-encoded DeleteCollectionRequest message.
type DeleteCollectionRequest struct{ o zap.Object }

// WrapDeleteCollectionRequest parses b and returns a typed view.
func WrapDeleteCollectionRequest(b []byte) (DeleteCollectionRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return DeleteCollectionRequest{}, err
	}
	return DeleteCollectionRequest{o: m.Root()}, nil
}

func (t DeleteCollectionRequest) Collection() string {
	return t.o.Text(deleteCollectionRequestCollectionOff)
}

// DeleteCollectionRequestInput collects the field values for NewDeleteCollectionRequest.
type DeleteCollectionRequestInput struct {
	Collection string
}

// NewDeleteCollectionRequest builds a ZAP-encoded DeleteCollectionRequest message from in.
func NewDeleteCollectionRequest(in DeleteCollectionRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(deleteCollectionRequestSize)
	ob.SetText(deleteCollectionRequestCollectionOff, in.Collection)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- DeleteCollectionResponse ---

const deleteCollectionResponseSize = 0

// DeleteCollectionResponse is a zero-copy view into a ZAP-encoded
// DeleteCollectionResponse message. DeleteCollectionResponse carries no fields.
type DeleteCollectionResponse struct{ o zap.Object }

// WrapDeleteCollectionResponse parses b and returns a typed view.
func WrapDeleteCollectionResponse(b []byte) (DeleteCollectionResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return DeleteCollectionResponse{}, err
	}
	return DeleteCollectionResponse{o: m.Root()}, nil
}

// DeleteCollectionResponseInput collects the field values for
// NewDeleteCollectionResponse. DeleteCollectionResponse has no fields.
type DeleteCollectionResponseInput struct{}

// NewDeleteCollectionResponse builds a ZAP-encoded DeleteCollectionResponse message.
func NewDeleteCollectionResponse(in DeleteCollectionResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(deleteCollectionResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- StatisticsRequest ---

const (
	statisticsRequestReplicationOff = 0  // replication (field 1, string)
	statisticsRequestCollectionOff  = 8  // collection (field 2, string)
	statisticsRequestTtlOff         = 16 // ttl (field 3, string)
	statisticsRequestDiskTypeOff    = 24 // disk_type (field 4, string)
	statisticsRequestSize           = 32
)

// StatisticsRequest is a zero-copy view into a ZAP-encoded StatisticsRequest message.
type StatisticsRequest struct{ o zap.Object }

// WrapStatisticsRequest parses b and returns a typed view.
func WrapStatisticsRequest(b []byte) (StatisticsRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return StatisticsRequest{}, err
	}
	return StatisticsRequest{o: m.Root()}, nil
}

func (t StatisticsRequest) Replication() string { return t.o.Text(statisticsRequestReplicationOff) }
func (t StatisticsRequest) Collection() string  { return t.o.Text(statisticsRequestCollectionOff) }
func (t StatisticsRequest) Ttl() string         { return t.o.Text(statisticsRequestTtlOff) }
func (t StatisticsRequest) DiskType() string    { return t.o.Text(statisticsRequestDiskTypeOff) }

// StatisticsRequestInput collects the field values for NewStatisticsRequest.
type StatisticsRequestInput struct {
	Replication string
	Collection  string
	Ttl         string
	DiskType    string
}

// NewStatisticsRequest builds a ZAP-encoded StatisticsRequest message from in.
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
	statisticsResponseTotalSizeOff = 0  // total_size (field 4, uint64)
	statisticsResponseUsedSizeOff  = 8  // used_size (field 5, uint64)
	statisticsResponseFileCountOff = 16 // file_count (field 6, uint64)
	statisticsResponseSize         = 24
)

// StatisticsResponse is a zero-copy view into a ZAP-encoded StatisticsResponse message.
type StatisticsResponse struct{ o zap.Object }

// WrapStatisticsResponse parses b and returns a typed view.
func WrapStatisticsResponse(b []byte) (StatisticsResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return StatisticsResponse{}, err
	}
	return StatisticsResponse{o: m.Root()}, nil
}

func (t StatisticsResponse) TotalSize() uint64 { return t.o.Uint64(statisticsResponseTotalSizeOff) }
func (t StatisticsResponse) UsedSize() uint64  { return t.o.Uint64(statisticsResponseUsedSizeOff) }
func (t StatisticsResponse) FileCount() uint64 { return t.o.Uint64(statisticsResponseFileCountOff) }

// StatisticsResponseInput collects the field values for NewStatisticsResponse.
type StatisticsResponseInput struct {
	TotalSize uint64
	UsedSize  uint64
	FileCount uint64
}

// NewStatisticsResponse builds a ZAP-encoded StatisticsResponse message from in.
func NewStatisticsResponse(in StatisticsResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(statisticsResponseSize)
	ob.SetUint64(statisticsResponseTotalSizeOff, in.TotalSize)
	ob.SetUint64(statisticsResponseUsedSizeOff, in.UsedSize)
	ob.SetUint64(statisticsResponseFileCountOff, in.FileCount)
	ob.FinishAsRoot()
	return b.Finish()
}
