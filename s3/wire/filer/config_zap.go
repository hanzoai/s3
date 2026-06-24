// Code generated from filer.proto; DO NOT EDIT.

package filerwire

import (
	zap "github.com/zap-proto/go"
)

// --- PingRequest ---

const (
	pingRequestTargetOff     = 0 // target (field 1, string)
	pingRequestTargetTypeOff = 8 // target_type (field 2, string)
	pingRequestSize          = 16
)

// PingRequest is a zero-copy view into a ZAP-encoded PingRequest message.
type PingRequest struct{ o zap.Object }

// WrapPingRequest parses b and returns a typed view.
func WrapPingRequest(b []byte) (PingRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PingRequest{}, err
	}
	return PingRequest{o: m.Root()}, nil
}

func (t PingRequest) Target() string     { return t.o.Text(pingRequestTargetOff) }
func (t PingRequest) TargetType() string { return t.o.Text(pingRequestTargetTypeOff) }

// PingRequestInput collects the field values for NewPingRequest.
type PingRequestInput struct {
	Target     string
	TargetType string
}

// NewPingRequest builds a ZAP-encoded PingRequest message from in.
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
	pingResponseStartTimeNsOff  = 0  // start_time_ns (field 1, int64)
	pingResponseRemoteTimeNsOff = 8  // remote_time_ns (field 2, int64)
	pingResponseStopTimeNsOff   = 16 // stop_time_ns (field 3, int64)
	pingResponseSize            = 24
)

// PingResponse is a zero-copy view into a ZAP-encoded PingResponse message.
type PingResponse struct{ o zap.Object }

// WrapPingResponse parses b and returns a typed view.
func WrapPingResponse(b []byte) (PingResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PingResponse{}, err
	}
	return PingResponse{o: m.Root()}, nil
}

func (t PingResponse) StartTimeNs() int64  { return t.o.Int64(pingResponseStartTimeNsOff) }
func (t PingResponse) RemoteTimeNs() int64 { return t.o.Int64(pingResponseRemoteTimeNsOff) }
func (t PingResponse) StopTimeNs() int64   { return t.o.Int64(pingResponseStopTimeNsOff) }

// PingResponseInput collects the field values for NewPingResponse.
type PingResponseInput struct {
	StartTimeNs  int64
	RemoteTimeNs int64
	StopTimeNs   int64
}

// NewPingResponse builds a ZAP-encoded PingResponse message from in.
func NewPingResponse(in PingResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(pingResponseSize)
	ob.SetInt64(pingResponseStartTimeNsOff, in.StartTimeNs)
	ob.SetInt64(pingResponseRemoteTimeNsOff, in.RemoteTimeNs)
	ob.SetInt64(pingResponseStopTimeNsOff, in.StopTimeNs)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- GetFilerConfigurationRequest ---

const getFilerConfigurationRequestSize = 0

// GetFilerConfigurationRequest is a zero-copy view into a ZAP-encoded
// GetFilerConfigurationRequest message. It carries no fields.
type GetFilerConfigurationRequest struct{ o zap.Object }

// WrapGetFilerConfigurationRequest parses b and returns a typed view.
func WrapGetFilerConfigurationRequest(b []byte) (GetFilerConfigurationRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return GetFilerConfigurationRequest{}, err
	}
	return GetFilerConfigurationRequest{o: m.Root()}, nil
}

// GetFilerConfigurationRequestInput collects the field values for
// NewGetFilerConfigurationRequest. It has no fields.
type GetFilerConfigurationRequestInput struct{}

// NewGetFilerConfigurationRequest builds a ZAP-encoded GetFilerConfigurationRequest message.
func NewGetFilerConfigurationRequest(in GetFilerConfigurationRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(getFilerConfigurationRequestSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- GetFilerConfigurationResponse ---

const (
	getFilerConfigurationResponseMastersOff            = 0  // masters (field 1, repeated string)
	getFilerConfigurationResponseReplicationOff        = 8  // replication (field 2, string)
	getFilerConfigurationResponseCollectionOff         = 16 // collection (field 3, string)
	getFilerConfigurationResponseMaxMbOff              = 24 // max_mb (field 4, uint32)
	getFilerConfigurationResponseDirBucketsOff         = 32 // dir_buckets (field 5, string)
	getFilerConfigurationResponseCipherOff             = 40 // cipher (field 7, bool)
	getFilerConfigurationResponseSignatureOff          = 44 // signature (field 8, int32)
	getFilerConfigurationResponseMetricsAddressOff     = 48 // metrics_address (field 9, string)
	getFilerConfigurationResponseMetricsIntervalSecOff = 56 // metrics_interval_sec (field 10, int32)
	getFilerConfigurationResponseVersionOff            = 64 // version (field 11, string)
	getFilerConfigurationResponseClusterIDOff          = 72 // cluster_id (field 12, string)
	getFilerConfigurationResponseFilerGroupOff         = 80 // filer_group (field 13, string)
	getFilerConfigurationResponseMajorVersionOff       = 88 // major_version (field 14, int32)
	getFilerConfigurationResponseMinorVersionOff       = 92 // minor_version (field 15, int32)
	getFilerConfigurationResponseSize                  = 96
)

// GetFilerConfigurationResponse is a zero-copy view into a ZAP-encoded
// GetFilerConfigurationResponse message.
type GetFilerConfigurationResponse struct{ o zap.Object }

// WrapGetFilerConfigurationResponse parses b and returns a typed view.
func WrapGetFilerConfigurationResponse(b []byte) (GetFilerConfigurationResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return GetFilerConfigurationResponse{}, err
	}
	return GetFilerConfigurationResponse{o: m.Root()}, nil
}

func (t GetFilerConfigurationResponse) Replication() string {
	return t.o.Text(getFilerConfigurationResponseReplicationOff)
}
func (t GetFilerConfigurationResponse) Collection() string {
	return t.o.Text(getFilerConfigurationResponseCollectionOff)
}
func (t GetFilerConfigurationResponse) MaxMb() uint32 {
	return t.o.Uint32(getFilerConfigurationResponseMaxMbOff)
}
func (t GetFilerConfigurationResponse) DirBuckets() string {
	return t.o.Text(getFilerConfigurationResponseDirBucketsOff)
}
func (t GetFilerConfigurationResponse) Cipher() bool {
	return t.o.Bool(getFilerConfigurationResponseCipherOff)
}
func (t GetFilerConfigurationResponse) Signature() int32 {
	return t.o.Int32(getFilerConfigurationResponseSignatureOff)
}
func (t GetFilerConfigurationResponse) MetricsAddress() string {
	return t.o.Text(getFilerConfigurationResponseMetricsAddressOff)
}
func (t GetFilerConfigurationResponse) MetricsIntervalSec() int32 {
	return t.o.Int32(getFilerConfigurationResponseMetricsIntervalSecOff)
}
func (t GetFilerConfigurationResponse) Version() string {
	return t.o.Text(getFilerConfigurationResponseVersionOff)
}
func (t GetFilerConfigurationResponse) ClusterID() string {
	return t.o.Text(getFilerConfigurationResponseClusterIDOff)
}
func (t GetFilerConfigurationResponse) FilerGroup() string {
	return t.o.Text(getFilerConfigurationResponseFilerGroupOff)
}
func (t GetFilerConfigurationResponse) MajorVersion() int32 {
	return t.o.Int32(getFilerConfigurationResponseMajorVersionOff)
}
func (t GetFilerConfigurationResponse) MinorVersion() int32 {
	return t.o.Int32(getFilerConfigurationResponseMinorVersionOff)
}

// MastersLen returns the number of masters elements (proto field 1, repeated string).
func (t GetFilerConfigurationResponse) MastersLen() int {
	return t.o.List(getFilerConfigurationResponseMastersOff).Len()
}

// Master returns the i-th masters element.
func (t GetFilerConfigurationResponse) Master(i int) string {
	return elemText(t.o.List(getFilerConfigurationResponseMastersOff), i)
}

// GetFilerConfigurationResponseInput collects the field values for
// NewGetFilerConfigurationResponse.
type GetFilerConfigurationResponseInput struct {
	Masters            []string
	Replication        string
	Collection         string
	MaxMb              uint32
	DirBuckets         string
	Cipher             bool
	Signature          int32
	MetricsAddress     string
	MetricsIntervalSec int32
	Version            string
	ClusterID          string
	FilerGroup         string
	MajorVersion       int32
	MinorVersion       int32
}

// NewGetFilerConfigurationResponse builds a ZAP-encoded
// GetFilerConfigurationResponse message from in.
func NewGetFilerConfigurationResponse(in GetFilerConfigurationResponseInput) []byte {
	b := zap.NewBuilder(512)
	mastersOff, mastersLen := addStringList(b, in.Masters)
	ob := b.StartObject(getFilerConfigurationResponseSize)
	ob.SetList(getFilerConfigurationResponseMastersOff, mastersOff, mastersLen)
	ob.SetText(getFilerConfigurationResponseReplicationOff, in.Replication)
	ob.SetText(getFilerConfigurationResponseCollectionOff, in.Collection)
	ob.SetUint32(getFilerConfigurationResponseMaxMbOff, in.MaxMb)
	ob.SetText(getFilerConfigurationResponseDirBucketsOff, in.DirBuckets)
	ob.SetBool(getFilerConfigurationResponseCipherOff, in.Cipher)
	ob.SetInt32(getFilerConfigurationResponseSignatureOff, in.Signature)
	ob.SetText(getFilerConfigurationResponseMetricsAddressOff, in.MetricsAddress)
	ob.SetInt32(getFilerConfigurationResponseMetricsIntervalSecOff, in.MetricsIntervalSec)
	ob.SetText(getFilerConfigurationResponseVersionOff, in.Version)
	ob.SetText(getFilerConfigurationResponseClusterIDOff, in.ClusterID)
	ob.SetText(getFilerConfigurationResponseFilerGroupOff, in.FilerGroup)
	ob.SetInt32(getFilerConfigurationResponseMajorVersionOff, in.MajorVersion)
	ob.SetInt32(getFilerConfigurationResponseMinorVersionOff, in.MinorVersion)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- FilerConfPathConf (FilerConf.PathConf nested message) ---

const (
	filerConfPathConfLocationPrefixOff           = 0  // location_prefix (field 1, string)
	filerConfPathConfCollectionOff               = 8  // collection (field 2, string)
	filerConfPathConfReplicationOff              = 16 // replication (field 3, string)
	filerConfPathConfTtlOff                      = 24 // ttl (field 4, string)
	filerConfPathConfDiskTypeOff                 = 32 // disk_type (field 5, string)
	filerConfPathConfFsyncOff                    = 40 // fsync (field 6, bool)
	filerConfPathConfVolumeGrowthCountOff        = 44 // volume_growth_count (field 7, uint32)
	filerConfPathConfReadOnlyOff                 = 48 // read_only (field 8, bool)
	filerConfPathConfDataCenterOff               = 56 // data_center (field 9, string)
	filerConfPathConfRackOff                     = 64 // rack (field 10, string)
	filerConfPathConfDataNodeOff                 = 72 // data_node (field 11, string)
	filerConfPathConfMaxFileNameLengthOff        = 80 // max_file_name_length (field 12, uint32)
	filerConfPathConfDisableChunkDeletionOff     = 84 // disable_chunk_deletion (field 13, bool)
	filerConfPathConfWormOff                     = 85 // worm (field 14, bool)
	filerConfPathConfWormGracePeriodSecondsOff   = 88 // worm_grace_period_seconds (field 15, uint64)
	filerConfPathConfWormRetentionTimeSecondsOff = 96 // worm_retention_time_seconds (field 16, uint64)
	filerConfPathConfSize                        = 104
)

// FilerConfPathConf is a zero-copy view into one encoded FilerConf.PathConf.
type FilerConfPathConf struct{ o zap.Object }

// WrapFilerConfPathConf parses b and returns a typed view.
func WrapFilerConfPathConf(b []byte) (FilerConfPathConf, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return FilerConfPathConf{}, err
	}
	return FilerConfPathConf{o: m.Root()}, nil
}

func (t FilerConfPathConf) LocationPrefix() string {
	return t.o.Text(filerConfPathConfLocationPrefixOff)
}
func (t FilerConfPathConf) Collection() string  { return t.o.Text(filerConfPathConfCollectionOff) }
func (t FilerConfPathConf) Replication() string { return t.o.Text(filerConfPathConfReplicationOff) }
func (t FilerConfPathConf) Ttl() string         { return t.o.Text(filerConfPathConfTtlOff) }
func (t FilerConfPathConf) DiskType() string    { return t.o.Text(filerConfPathConfDiskTypeOff) }
func (t FilerConfPathConf) Fsync() bool         { return t.o.Bool(filerConfPathConfFsyncOff) }
func (t FilerConfPathConf) VolumeGrowthCount() uint32 {
	return t.o.Uint32(filerConfPathConfVolumeGrowthCountOff)
}
func (t FilerConfPathConf) ReadOnly() bool     { return t.o.Bool(filerConfPathConfReadOnlyOff) }
func (t FilerConfPathConf) DataCenter() string { return t.o.Text(filerConfPathConfDataCenterOff) }
func (t FilerConfPathConf) Rack() string       { return t.o.Text(filerConfPathConfRackOff) }
func (t FilerConfPathConf) DataNode() string   { return t.o.Text(filerConfPathConfDataNodeOff) }
func (t FilerConfPathConf) MaxFileNameLength() uint32 {
	return t.o.Uint32(filerConfPathConfMaxFileNameLengthOff)
}
func (t FilerConfPathConf) DisableChunkDeletion() bool {
	return t.o.Bool(filerConfPathConfDisableChunkDeletionOff)
}
func (t FilerConfPathConf) Worm() bool { return t.o.Bool(filerConfPathConfWormOff) }
func (t FilerConfPathConf) WormGracePeriodSeconds() uint64 {
	return t.o.Uint64(filerConfPathConfWormGracePeriodSecondsOff)
}
func (t FilerConfPathConf) WormRetentionTimeSeconds() uint64 {
	return t.o.Uint64(filerConfPathConfWormRetentionTimeSecondsOff)
}

// FilerConfPathConfInput collects the field values for NewFilerConfPathConf.
type FilerConfPathConfInput struct {
	LocationPrefix           string
	Collection               string
	Replication              string
	Ttl                      string
	DiskType                 string
	Fsync                    bool
	VolumeGrowthCount        uint32
	ReadOnly                 bool
	DataCenter               string
	Rack                     string
	DataNode                 string
	MaxFileNameLength        uint32
	DisableChunkDeletion     bool
	Worm                     bool
	WormGracePeriodSeconds   uint64
	WormRetentionTimeSeconds uint64
}

// NewFilerConfPathConf builds a ZAP-encoded FilerConf.PathConf from in.
func NewFilerConfPathConf(in FilerConfPathConfInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(filerConfPathConfSize)
	ob.SetText(filerConfPathConfLocationPrefixOff, in.LocationPrefix)
	ob.SetText(filerConfPathConfCollectionOff, in.Collection)
	ob.SetText(filerConfPathConfReplicationOff, in.Replication)
	ob.SetText(filerConfPathConfTtlOff, in.Ttl)
	ob.SetText(filerConfPathConfDiskTypeOff, in.DiskType)
	ob.SetBool(filerConfPathConfFsyncOff, in.Fsync)
	ob.SetUint32(filerConfPathConfVolumeGrowthCountOff, in.VolumeGrowthCount)
	ob.SetBool(filerConfPathConfReadOnlyOff, in.ReadOnly)
	ob.SetText(filerConfPathConfDataCenterOff, in.DataCenter)
	ob.SetText(filerConfPathConfRackOff, in.Rack)
	ob.SetText(filerConfPathConfDataNodeOff, in.DataNode)
	ob.SetUint32(filerConfPathConfMaxFileNameLengthOff, in.MaxFileNameLength)
	ob.SetBool(filerConfPathConfDisableChunkDeletionOff, in.DisableChunkDeletion)
	ob.SetBool(filerConfPathConfWormOff, in.Worm)
	ob.SetUint64(filerConfPathConfWormGracePeriodSecondsOff, in.WormGracePeriodSeconds)
	ob.SetUint64(filerConfPathConfWormRetentionTimeSecondsOff, in.WormRetentionTimeSeconds)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- FilerConf ---

const (
	filerConfVersionOff   = 0 // version (field 1, int32)
	filerConfLocationsOff = 8 // locations (field 2, repeated PathConf)
	filerConfSize         = 16
)

// FilerConf is a zero-copy view into a ZAP-encoded FilerConf message.
type FilerConf struct{ o zap.Object }

// WrapFilerConf parses b and returns a typed view.
func WrapFilerConf(b []byte) (FilerConf, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return FilerConf{}, err
	}
	return FilerConf{o: m.Root()}, nil
}

func (t FilerConf) Version() int32 { return t.o.Int32(filerConfVersionOff) }

// LocationsLen returns the number of locations elements (proto field 2, repeated PathConf).
func (t FilerConf) LocationsLen() int { return t.o.List(filerConfLocationsOff).Len() }

// Locations returns the i-th locations element as a typed FilerConfPathConf view.
func (t FilerConf) Locations(i int) FilerConfPathConf {
	return FilerConfPathConf{o: t.o.List(filerConfLocationsOff).ObjectAt(i)}
}

// FilerConfInput collects the field values for NewFilerConf. Each Locations
// element is a standalone FilerConfPathConf buffer (from NewFilerConfPathConf).
type FilerConfInput struct {
	Version   int32
	Locations [][]byte
}

// NewFilerConf builds a ZAP-encoded FilerConf message from in.
func NewFilerConf(in FilerConfInput) []byte {
	b := zap.NewBuilder(512)
	locOff, locLen := addObjectList(b, in.Locations)
	ob := b.StartObject(filerConfSize)
	ob.SetInt32(filerConfVersionOff, in.Version)
	ob.SetList(filerConfLocationsOff, locOff, locLen)
	ob.FinishAsRoot()
	return b.Finish()
}
