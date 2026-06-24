// Code generated from filer.proto; DO NOT EDIT.

package filerwire

import (
	zap "github.com/zap-proto/go"
)

// --- MountRegisterRequest ---

const (
	mountRegisterRequestPeerAddrOff   = 0  // peer_addr (field 1, string)
	mountRegisterRequestRackOff       = 8  // rack (field 2, string)
	mountRegisterRequestTtlSecondsOff = 16 // ttl_seconds (field 3, int32)
	mountRegisterRequestDataCenterOff = 24 // data_center (field 4, string)
	mountRegisterRequestSize          = 32
)

// MountRegisterRequest is a zero-copy view into a ZAP-encoded MountRegisterRequest message.
type MountRegisterRequest struct{ o zap.Object }

// WrapMountRegisterRequest parses b and returns a typed view.
func WrapMountRegisterRequest(b []byte) (MountRegisterRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return MountRegisterRequest{}, err
	}
	return MountRegisterRequest{o: m.Root()}, nil
}

func (t MountRegisterRequest) PeerAddr() string   { return t.o.Text(mountRegisterRequestPeerAddrOff) }
func (t MountRegisterRequest) Rack() string       { return t.o.Text(mountRegisterRequestRackOff) }
func (t MountRegisterRequest) TtlSeconds() int32  { return t.o.Int32(mountRegisterRequestTtlSecondsOff) }
func (t MountRegisterRequest) DataCenter() string { return t.o.Text(mountRegisterRequestDataCenterOff) }

// MountRegisterRequestInput collects the field values for NewMountRegisterRequest.
type MountRegisterRequestInput struct {
	PeerAddr   string
	Rack       string
	TtlSeconds int32
	DataCenter string
}

// NewMountRegisterRequest builds a ZAP-encoded MountRegisterRequest message from in.
func NewMountRegisterRequest(in MountRegisterRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(mountRegisterRequestSize)
	ob.SetText(mountRegisterRequestPeerAddrOff, in.PeerAddr)
	ob.SetText(mountRegisterRequestRackOff, in.Rack)
	ob.SetInt32(mountRegisterRequestTtlSecondsOff, in.TtlSeconds)
	ob.SetText(mountRegisterRequestDataCenterOff, in.DataCenter)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- MountRegisterResponse ---

const mountRegisterResponseSize = 0

// MountRegisterResponse is a zero-copy view into a ZAP-encoded MountRegisterResponse
// message. MountRegisterResponse carries no fields.
type MountRegisterResponse struct{ o zap.Object }

// WrapMountRegisterResponse parses b and returns a typed view.
func WrapMountRegisterResponse(b []byte) (MountRegisterResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return MountRegisterResponse{}, err
	}
	return MountRegisterResponse{o: m.Root()}, nil
}

// MountRegisterResponseInput collects the field values for NewMountRegisterResponse.
// MountRegisterResponse has no fields.
type MountRegisterResponseInput struct{}

// NewMountRegisterResponse builds a ZAP-encoded MountRegisterResponse message.
func NewMountRegisterResponse(in MountRegisterResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(mountRegisterResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- MountListRequest ---

const mountListRequestSize = 0

// MountListRequest is a zero-copy view into a ZAP-encoded MountListRequest message.
// MountListRequest carries no fields.
type MountListRequest struct{ o zap.Object }

// WrapMountListRequest parses b and returns a typed view.
func WrapMountListRequest(b []byte) (MountListRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return MountListRequest{}, err
	}
	return MountListRequest{o: m.Root()}, nil
}

// MountListRequestInput collects the field values for NewMountListRequest.
// MountListRequest has no fields.
type MountListRequestInput struct{}

// NewMountListRequest builds a ZAP-encoded MountListRequest message.
func NewMountListRequest(in MountListRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(mountListRequestSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- MountInfo ---

const (
	mountInfoPeerAddrOff   = 0  // peer_addr (field 1, string)
	mountInfoRackOff       = 8  // rack (field 2, string)
	mountInfoLastSeenNsOff = 16 // last_seen_ns (field 3, int64)
	mountInfoDataCenterOff = 24 // data_center (field 4, string)
	mountInfoSize          = 32
)

// MountInfo is a zero-copy view into a ZAP-encoded MountInfo message.
type MountInfo struct{ o zap.Object }

// WrapMountInfo parses b and returns a typed view.
func WrapMountInfo(b []byte) (MountInfo, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return MountInfo{}, err
	}
	return MountInfo{o: m.Root()}, nil
}

func (t MountInfo) PeerAddr() string   { return t.o.Text(mountInfoPeerAddrOff) }
func (t MountInfo) Rack() string       { return t.o.Text(mountInfoRackOff) }
func (t MountInfo) LastSeenNs() int64  { return t.o.Int64(mountInfoLastSeenNsOff) }
func (t MountInfo) DataCenter() string { return t.o.Text(mountInfoDataCenterOff) }

// MountInfoInput collects the field values for NewMountInfo.
type MountInfoInput struct {
	PeerAddr   string
	Rack       string
	LastSeenNs int64
	DataCenter string
}

// NewMountInfo builds a ZAP-encoded MountInfo message from in.
func NewMountInfo(in MountInfoInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(mountInfoSize)
	ob.SetText(mountInfoPeerAddrOff, in.PeerAddr)
	ob.SetText(mountInfoRackOff, in.Rack)
	ob.SetInt64(mountInfoLastSeenNsOff, in.LastSeenNs)
	ob.SetText(mountInfoDataCenterOff, in.DataCenter)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- MountListResponse ---

const (
	mountListResponseMountsOff = 0 // mounts (field 1, repeated MountInfo)
	mountListResponseSize      = 8
)

// MountListResponse is a zero-copy view into a ZAP-encoded MountListResponse message.
type MountListResponse struct{ o zap.Object }

// WrapMountListResponse parses b and returns a typed view.
func WrapMountListResponse(b []byte) (MountListResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return MountListResponse{}, err
	}
	return MountListResponse{o: m.Root()}, nil
}

// MountsLen returns the number of mounts elements (proto field 1, repeated MountInfo).
func (t MountListResponse) MountsLen() int { return t.o.List(mountListResponseMountsOff).Len() }

// Mounts returns the i-th mounts element as a typed MountInfo view.
func (t MountListResponse) Mounts(i int) MountInfo {
	return MountInfo{o: t.o.List(mountListResponseMountsOff).ObjectAt(i)}
}

// MountListResponseInput collects the field values for NewMountListResponse. Each
// Mounts element is a standalone MountInfo buffer (from NewMountInfo).
type MountListResponseInput struct {
	Mounts [][]byte
}

// NewMountListResponse builds a ZAP-encoded MountListResponse message from in.
func NewMountListResponse(in MountListResponseInput) []byte {
	b := zap.NewBuilder(256)
	mountsOff, mountsLen := addObjectList(b, in.Mounts)
	ob := b.StartObject(mountListResponseSize)
	ob.SetList(mountListResponseMountsOff, mountsOff, mountsLen)
	ob.FinishAsRoot()
	return b.Finish()
}
