// Code generated from mount_peer.proto; DO NOT EDIT.

package mount_peerwire

import (
	zap "github.com/zap-proto/go"
)

// --- ChunkAnnounceRequest ---

const (
	chunkAnnounceRequestFileIdsOff    = 0
	chunkAnnounceRequestPeerAddrOff   = 8
	chunkAnnounceRequestRackOff       = 16
	chunkAnnounceRequestTtlSecondsOff = 24
	chunkAnnounceRequestDataCenterOff = 28
	chunkAnnounceRequestSize          = 36
)

// ChunkAnnounceRequest is a zero-copy view into a ZAP-encoded
// ChunkAnnounceRequest message.
type ChunkAnnounceRequest struct{ o zap.Object }

// WrapChunkAnnounceRequest parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapChunkAnnounceRequest(b []byte) (ChunkAnnounceRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ChunkAnnounceRequest{}, err
	}
	return ChunkAnnounceRequest{o: m.Root()}, nil
}

// FileIds reads the file_ids field (proto field 1, repeated string). Each
// element is a raw-bytes entry; read it with List.BytesAt(i).
func (t ChunkAnnounceRequest) FileIds() zap.List { return t.o.List(chunkAnnounceRequestFileIdsOff) }

// PeerAddr reads the peer_addr field (proto field 2, string).
func (t ChunkAnnounceRequest) PeerAddr() string { return t.o.Text(chunkAnnounceRequestPeerAddrOff) }

// Rack reads the rack field (proto field 3, string).
func (t ChunkAnnounceRequest) Rack() string { return t.o.Text(chunkAnnounceRequestRackOff) }

// TtlSeconds reads the ttl_seconds field (proto field 4, int32).
func (t ChunkAnnounceRequest) TtlSeconds() int32 {
	return t.o.Int32(chunkAnnounceRequestTtlSecondsOff)
}

// DataCenter reads the data_center field (proto field 5, string).
func (t ChunkAnnounceRequest) DataCenter() string {
	return t.o.Text(chunkAnnounceRequestDataCenterOff)
}

// ChunkAnnounceRequestInput collects the field values for
// NewChunkAnnounceRequest. FileIds is a list of pre-encoded raw-bytes entries
// (one per file id).
type ChunkAnnounceRequestInput struct {
	FileIds    [][]byte
	PeerAddr   string
	Rack       string
	TtlSeconds int32
	DataCenter string
}

// NewChunkAnnounceRequest builds a ZAP-encoded ChunkAnnounceRequest message
// from in and returns the bytes.
func NewChunkAnnounceRequest(in ChunkAnnounceRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(chunkAnnounceRequestSize)
	ob.SetText(chunkAnnounceRequestPeerAddrOff, in.PeerAddr)
	ob.SetText(chunkAnnounceRequestRackOff, in.Rack)
	ob.SetInt32(chunkAnnounceRequestTtlSecondsOff, in.TtlSeconds)
	ob.SetText(chunkAnnounceRequestDataCenterOff, in.DataCenter)
	fileIdsLB := b.StartList(0)
	for _, elem := range in.FileIds {
		fileIdsLB.AddObjectBytes(elem)
	}
	ob.SetList(chunkAnnounceRequestFileIdsOff, fileIdsLB.FinishOffset(), len(in.FileIds))
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ChunkAnnounceResponse ---

const (
	chunkAnnounceResponseRejectedFileIdsOff = 0
	chunkAnnounceResponseSize               = 8
)

// ChunkAnnounceResponse is a zero-copy view into a ZAP-encoded
// ChunkAnnounceResponse message.
type ChunkAnnounceResponse struct{ o zap.Object }

// WrapChunkAnnounceResponse parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapChunkAnnounceResponse(b []byte) (ChunkAnnounceResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ChunkAnnounceResponse{}, err
	}
	return ChunkAnnounceResponse{o: m.Root()}, nil
}

// RejectedFileIds reads the rejected_file_ids field (proto field 1, repeated
// string): fids the receiver is not the owner of. Each element is a raw-bytes
// entry; read it with List.BytesAt(i).
func (t ChunkAnnounceResponse) RejectedFileIds() zap.List {
	return t.o.List(chunkAnnounceResponseRejectedFileIdsOff)
}

// ChunkAnnounceResponseInput collects the field values for
// NewChunkAnnounceResponse. RejectedFileIds is a list of pre-encoded raw-bytes
// entries (one per file id).
type ChunkAnnounceResponseInput struct {
	RejectedFileIds [][]byte
}

// NewChunkAnnounceResponse builds a ZAP-encoded ChunkAnnounceResponse message
// from in and returns the bytes.
func NewChunkAnnounceResponse(in ChunkAnnounceResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(chunkAnnounceResponseSize)
	rejectedFileIdsLB := b.StartList(0)
	for _, elem := range in.RejectedFileIds {
		rejectedFileIdsLB.AddObjectBytes(elem)
	}
	ob.SetList(chunkAnnounceResponseRejectedFileIdsOff, rejectedFileIdsLB.FinishOffset(), len(in.RejectedFileIds))
	ob.FinishAsRoot()
	return b.Finish()
}
