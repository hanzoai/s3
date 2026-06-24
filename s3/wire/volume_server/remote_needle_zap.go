// Code generated from volume_server.proto; DO NOT EDIT.

package volume_serverwire

import (
	zap "github.com/zap-proto/go"

	remotewire "github.com/hanzoai/s3/s3/wire/remote"
)

// --- FetchAndWriteNeedleReplica (FetchAndWriteNeedleRequest.Replica) ---
//
// FetchAndWriteNeedleReplica is the flattened form of the message type declared
// inline as FetchAndWriteNeedleRequest.Replica in the .proto.

const (
	fetchAndWriteNeedleReplicaURLOff       = 0
	fetchAndWriteNeedleReplicaPublicURLOff = 8
	fetchAndWriteNeedleReplicaGrpcPortOff  = 16
	fetchAndWriteNeedleReplicaSize         = 20
)

// FetchAndWriteNeedleReplica is a zero-copy view into a ZAP-encoded
// FetchAndWriteNeedleRequest.Replica message.
type FetchAndWriteNeedleReplica struct{ o zap.Object }

// WrapFetchAndWriteNeedleReplica parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapFetchAndWriteNeedleReplica(b []byte) (FetchAndWriteNeedleReplica, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return FetchAndWriteNeedleReplica{}, err
	}
	return FetchAndWriteNeedleReplica{o: m.Root()}, nil
}

// URL reads the url field (proto field 1, string).
func (t FetchAndWriteNeedleReplica) URL() string {
	return t.o.Text(fetchAndWriteNeedleReplicaURLOff)
}

// PublicURL reads the public_url field (proto field 2, string).
func (t FetchAndWriteNeedleReplica) PublicURL() string {
	return t.o.Text(fetchAndWriteNeedleReplicaPublicURLOff)
}

// GrpcPort reads the grpc_port field (proto field 3, int32).
func (t FetchAndWriteNeedleReplica) GrpcPort() int32 {
	return t.o.Int32(fetchAndWriteNeedleReplicaGrpcPortOff)
}

// FetchAndWriteNeedleReplicaInput collects the field values for
// NewFetchAndWriteNeedleReplica.
type FetchAndWriteNeedleReplicaInput struct {
	URL       string
	PublicURL string
	GrpcPort  int32
}

// NewFetchAndWriteNeedleReplica builds a ZAP-encoded
// FetchAndWriteNeedleRequest.Replica message from in and returns the bytes.
func NewFetchAndWriteNeedleReplica(in FetchAndWriteNeedleReplicaInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(fetchAndWriteNeedleReplicaSize)
	ob.SetText(fetchAndWriteNeedleReplicaURLOff, in.URL)
	ob.SetText(fetchAndWriteNeedleReplicaPublicURLOff, in.PublicURL)
	ob.SetInt32(fetchAndWriteNeedleReplicaGrpcPortOff, in.GrpcPort)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- FetchAndWriteNeedleRequest ---
//
// remote_conf (field 15) and remote_location (field 16) are the cross-service
// remote_pb.RemoteConf / remote_pb.RemoteStorageLocation types; they are
// embedded as nested objects whose buffers are built by the remotewire package
// (NewRemoteConf / NewRemoteStorageLocation).

const (
	fetchAndWriteNeedleRequestVolumeIDOff            = 0
	fetchAndWriteNeedleRequestNeedleIDOff            = 8
	fetchAndWriteNeedleRequestCookieOff              = 16
	fetchAndWriteNeedleRequestOffsetOff              = 24
	fetchAndWriteNeedleRequestSizeFOff               = 32
	fetchAndWriteNeedleRequestReplicasOff            = 40
	fetchAndWriteNeedleRequestAuthOff                = 48
	fetchAndWriteNeedleRequestDownloadConcurrencyOff = 56
	fetchAndWriteNeedleRequestRemoteConfOff          = 60
	fetchAndWriteNeedleRequestRemoteLocationOff      = 68
	fetchAndWriteNeedleRequestSize                   = 76
)

// FetchAndWriteNeedleRequest is a zero-copy view into a ZAP-encoded
// FetchAndWriteNeedleRequest message.
type FetchAndWriteNeedleRequest struct{ o zap.Object }

// WrapFetchAndWriteNeedleRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapFetchAndWriteNeedleRequest(b []byte) (FetchAndWriteNeedleRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return FetchAndWriteNeedleRequest{}, err
	}
	return FetchAndWriteNeedleRequest{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t FetchAndWriteNeedleRequest) VolumeID() uint32 {
	return t.o.Uint32(fetchAndWriteNeedleRequestVolumeIDOff)
}

// NeedleID reads the needle_id field (proto field 2, uint64).
func (t FetchAndWriteNeedleRequest) NeedleID() uint64 {
	return t.o.Uint64(fetchAndWriteNeedleRequestNeedleIDOff)
}

// Cookie reads the cookie field (proto field 3, uint32).
func (t FetchAndWriteNeedleRequest) Cookie() uint32 {
	return t.o.Uint32(fetchAndWriteNeedleRequestCookieOff)
}

// Offset reads the offset field (proto field 4, int64).
func (t FetchAndWriteNeedleRequest) Offset() int64 {
	return t.o.Int64(fetchAndWriteNeedleRequestOffsetOff)
}

// Size reads the size field (proto field 5, int64).
func (t FetchAndWriteNeedleRequest) Size() int64 {
	return t.o.Int64(fetchAndWriteNeedleRequestSizeFOff)
}

// ReplicasLen reports the number of replicas elements (proto field 6, repeated
// Replica).
func (t FetchAndWriteNeedleRequest) ReplicasLen() int {
	return t.o.List(fetchAndWriteNeedleRequestReplicasOff).Len()
}

// ReplicaAt returns the i-th replicas element. The bool is false when i is out
// of range or the element fails to parse.
func (t FetchAndWriteNeedleRequest) ReplicaAt(i int) (FetchAndWriteNeedleReplica, bool) {
	o := t.o.List(fetchAndWriteNeedleRequestReplicasOff).ObjectAt(i)
	if o.IsNull() {
		return FetchAndWriteNeedleReplica{}, false
	}
	return FetchAndWriteNeedleReplica{o: o}, true
}

// Auth reads the auth field (proto field 7, string).
func (t FetchAndWriteNeedleRequest) Auth() string {
	return t.o.Text(fetchAndWriteNeedleRequestAuthOff)
}

// DownloadConcurrency reads the download_concurrency field (proto field 8,
// int32).
func (t FetchAndWriteNeedleRequest) DownloadConcurrency() int32 {
	return t.o.Int32(fetchAndWriteNeedleRequestDownloadConcurrencyOff)
}

// RemoteConf reads the remote_conf field (proto field 15, message
// remote_pb.RemoteConf). The bool is false when the field is absent.
func (t FetchAndWriteNeedleRequest) RemoteConf() (remotewire.RemoteConf, bool) {
	bts := t.o.Bytes(fetchAndWriteNeedleRequestRemoteConfOff)
	if len(bts) == 0 {
		return remotewire.RemoteConf{}, false
	}
	v, err := remotewire.WrapRemoteConf(bts)
	if err != nil {
		return remotewire.RemoteConf{}, false
	}
	return v, true
}

// RemoteLocation reads the remote_location field (proto field 16, message
// remote_pb.RemoteStorageLocation). The bool is false when the field is absent.
func (t FetchAndWriteNeedleRequest) RemoteLocation() (remotewire.RemoteStorageLocation, bool) {
	bts := t.o.Bytes(fetchAndWriteNeedleRequestRemoteLocationOff)
	if len(bts) == 0 {
		return remotewire.RemoteStorageLocation{}, false
	}
	v, err := remotewire.WrapRemoteStorageLocation(bts)
	if err != nil {
		return remotewire.RemoteStorageLocation{}, false
	}
	return v, true
}

// FetchAndWriteNeedleRequestInput collects the field values for
// NewFetchAndWriteNeedleRequest. Each Replicas entry is a
// FetchAndWriteNeedleReplica buffer; RemoteConf and RemoteLocation are the
// cross-service remote messages' own ZAP buffers (from remotewire.NewRemoteConf
// / remotewire.NewRemoteStorageLocation).
type FetchAndWriteNeedleRequestInput struct {
	VolumeID            uint32
	NeedleID            uint64
	Cookie              uint32
	Offset              int64
	Size                int64
	Replicas            [][]byte
	Auth                string
	DownloadConcurrency int32
	RemoteConf          []byte
	RemoteLocation      []byte
}

// NewFetchAndWriteNeedleRequest builds a ZAP-encoded FetchAndWriteNeedleRequest
// message from in and returns the bytes.
func NewFetchAndWriteNeedleRequest(in FetchAndWriteNeedleRequestInput) []byte {
	b := zap.NewBuilder(512)
	listOff, listLen := addObjectList(b, in.Replicas)
	ob := b.StartObject(fetchAndWriteNeedleRequestSize)
	ob.SetUint32(fetchAndWriteNeedleRequestVolumeIDOff, in.VolumeID)
	ob.SetUint64(fetchAndWriteNeedleRequestNeedleIDOff, in.NeedleID)
	ob.SetUint32(fetchAndWriteNeedleRequestCookieOff, in.Cookie)
	ob.SetInt64(fetchAndWriteNeedleRequestOffsetOff, in.Offset)
	ob.SetInt64(fetchAndWriteNeedleRequestSizeFOff, in.Size)
	ob.SetList(fetchAndWriteNeedleRequestReplicasOff, listOff, listLen)
	ob.SetText(fetchAndWriteNeedleRequestAuthOff, in.Auth)
	ob.SetInt32(fetchAndWriteNeedleRequestDownloadConcurrencyOff, in.DownloadConcurrency)
	ob.SetBytes(fetchAndWriteNeedleRequestRemoteConfOff, in.RemoteConf)
	ob.SetBytes(fetchAndWriteNeedleRequestRemoteLocationOff, in.RemoteLocation)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- FetchAndWriteNeedleResponse ---

const (
	fetchAndWriteNeedleResponseETagOff = 0
	fetchAndWriteNeedleResponseSize    = 8
)

// FetchAndWriteNeedleResponse is a zero-copy view into a ZAP-encoded
// FetchAndWriteNeedleResponse message.
type FetchAndWriteNeedleResponse struct{ o zap.Object }

// WrapFetchAndWriteNeedleResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapFetchAndWriteNeedleResponse(b []byte) (FetchAndWriteNeedleResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return FetchAndWriteNeedleResponse{}, err
	}
	return FetchAndWriteNeedleResponse{o: m.Root()}, nil
}

// ETag reads the e_tag field (proto field 1, string).
func (t FetchAndWriteNeedleResponse) ETag() string {
	return t.o.Text(fetchAndWriteNeedleResponseETagOff)
}

// FetchAndWriteNeedleResponseInput collects the field values for
// NewFetchAndWriteNeedleResponse.
type FetchAndWriteNeedleResponseInput struct {
	ETag string
}

// NewFetchAndWriteNeedleResponse builds a ZAP-encoded FetchAndWriteNeedleResponse
// message from in and returns the bytes.
func NewFetchAndWriteNeedleResponse(in FetchAndWriteNeedleResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(fetchAndWriteNeedleResponseSize)
	ob.SetText(fetchAndWriteNeedleResponseETagOff, in.ETag)
	ob.FinishAsRoot()
	return b.Finish()
}
