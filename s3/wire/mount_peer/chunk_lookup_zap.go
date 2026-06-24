// Code generated from mount_peer.proto; DO NOT EDIT.

package mount_peerwire

import (
	zap "github.com/zap-proto/go"
)

// --- ChunkLookupRequest ---

const (
	chunkLookupRequestFileIdsOff = 0
	chunkLookupRequestSize       = 8
)

// ChunkLookupRequest is a zero-copy view into a ZAP-encoded ChunkLookupRequest
// message.
type ChunkLookupRequest struct{ o zap.Object }

// WrapChunkLookupRequest parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapChunkLookupRequest(b []byte) (ChunkLookupRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ChunkLookupRequest{}, err
	}
	return ChunkLookupRequest{o: m.Root()}, nil
}

// FileIds reads the file_ids field (proto field 1, repeated string). Each
// element is a raw-bytes entry; read it with List.BytesAt(i).
func (t ChunkLookupRequest) FileIds() zap.List { return t.o.List(chunkLookupRequestFileIdsOff) }

// ChunkLookupRequestInput collects the field values for NewChunkLookupRequest.
// FileIds is a list of pre-encoded raw-bytes entries (one per file id).
type ChunkLookupRequestInput struct {
	FileIds [][]byte
}

// NewChunkLookupRequest builds a ZAP-encoded ChunkLookupRequest message from in
// and returns the bytes.
func NewChunkLookupRequest(in ChunkLookupRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(chunkLookupRequestSize)
	fileIdsLB := b.StartList(0)
	for _, elem := range in.FileIds {
		fileIdsLB.AddObjectBytes(elem)
	}
	ob.SetList(chunkLookupRequestFileIdsOff, fileIdsLB.FinishOffset(), len(in.FileIds))
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ChunkLookupResponse ---

const (
	chunkLookupResponsePeersByFidOff      = 0
	chunkLookupResponseNotOwnerFileIdsOff = 8
	chunkLookupResponseSize               = 16
)

// ChunkLookupResponse is a zero-copy view into a ZAP-encoded ChunkLookupResponse
// message.
type ChunkLookupResponse struct{ o zap.Object }

// WrapChunkLookupResponse parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapChunkLookupResponse(b []byte) (ChunkLookupResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ChunkLookupResponse{}, err
	}
	return ChunkLookupResponse{o: m.Root()}, nil
}

// PeersByFid reads the peers_by_fid field (proto field 1,
// map<string, PeerSet>): a list of {key, value} entry objects. Read element i
// with List.ObjectAt(i) and wrap it as a ChunkLookupResponsePeersByFidEntry.
func (t ChunkLookupResponse) PeersByFid() zap.List {
	return t.o.List(chunkLookupResponsePeersByFidOff)
}

// NotOwnerFileIds reads the not_owner_file_ids field (proto field 2, repeated
// string): fids the receiver does not own. Each element is a raw-bytes entry;
// read it with List.BytesAt(i).
func (t ChunkLookupResponse) NotOwnerFileIds() zap.List {
	return t.o.List(chunkLookupResponseNotOwnerFileIdsOff)
}

// ChunkLookupResponseInput collects the field values for
// NewChunkLookupResponse. PeersByFid is a list of pre-encoded
// ChunkLookupResponsePeersByFidEntry buffers (build each with
// NewChunkLookupResponsePeersByFidEntry); NotOwnerFileIds is a list of
// pre-encoded raw-bytes entries.
type ChunkLookupResponseInput struct {
	PeersByFid      [][]byte
	NotOwnerFileIds [][]byte
}

// NewChunkLookupResponse builds a ZAP-encoded ChunkLookupResponse message from
// in and returns the bytes.
func NewChunkLookupResponse(in ChunkLookupResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(chunkLookupResponseSize)
	peersByFidLB := b.StartList(0)
	for _, elem := range in.PeersByFid {
		peersByFidLB.AddObjectBytes(elem)
	}
	ob.SetList(chunkLookupResponsePeersByFidOff, peersByFidLB.FinishOffset(), len(in.PeersByFid))
	notOwnerFileIdsLB := b.StartList(0)
	for _, elem := range in.NotOwnerFileIds {
		notOwnerFileIdsLB.AddObjectBytes(elem)
	}
	ob.SetList(chunkLookupResponseNotOwnerFileIdsOff, notOwnerFileIdsLB.FinishOffset(), len(in.NotOwnerFileIds))
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ChunkLookupResponsePeersByFidEntry (map<string, PeerSet> entry) ---

const (
	chunkLookupResponsePeersByFidEntryKeyOff   = 0
	chunkLookupResponsePeersByFidEntryValueOff = 8
	chunkLookupResponsePeersByFidEntrySize     = 12
)

// ChunkLookupResponsePeersByFidEntry is a zero-copy view into one entry of the
// peers_by_fid map: a string key and a nested PeerSet value.
type ChunkLookupResponsePeersByFidEntry struct{ o zap.Object }

// WrapChunkLookupResponsePeersByFidEntry parses b and returns a typed view.
// Returns an error if the wire-level checks (magic, version, size) fail.
func WrapChunkLookupResponsePeersByFidEntry(b []byte) (ChunkLookupResponsePeersByFidEntry, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ChunkLookupResponsePeersByFidEntry{}, err
	}
	return ChunkLookupResponsePeersByFidEntry{o: m.Root()}, nil
}

// Key reads the map entry key (string).
func (t ChunkLookupResponsePeersByFidEntry) Key() string {
	return t.o.Text(chunkLookupResponsePeersByFidEntryKeyOff)
}

// Value reads the map entry value (nested PeerSet message).
func (t ChunkLookupResponsePeersByFidEntry) Value() PeerSet {
	return PeerSet{o: t.o.Object(chunkLookupResponsePeersByFidEntryValueOff)}
}

// ChunkLookupResponsePeersByFidEntryInput collects the field values for
// NewChunkLookupResponsePeersByFidEntry. Value is a pre-built PeerSet sub-buffer
// (build it with NewPeerSet).
type ChunkLookupResponsePeersByFidEntryInput struct {
	Key   string
	Value []byte
}

// NewChunkLookupResponsePeersByFidEntry builds a ZAP-encoded peers_by_fid map
// entry from in and returns the bytes.
func NewChunkLookupResponsePeersByFidEntry(in ChunkLookupResponsePeersByFidEntryInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(chunkLookupResponsePeersByFidEntrySize)
	ob.SetText(chunkLookupResponsePeersByFidEntryKeyOff, in.Key)
	embedObject(b, ob, chunkLookupResponsePeersByFidEntryValueOff, in.Value)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- PeerSet ---

const (
	peerSetPeersOff = 0
	peerSetSize     = 8
)

// PeerSet is a zero-copy view into a ZAP-encoded PeerSet message.
type PeerSet struct{ o zap.Object }

// WrapPeerSet parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapPeerSet(b []byte) (PeerSet, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PeerSet{}, err
	}
	return PeerSet{o: m.Root()}, nil
}

// Peers reads the peers field (proto field 1, repeated PeerInfo): a list of
// out-of-line PeerInfo objects. Read element i with List.ObjectAt(i) and wrap
// it as a PeerInfo.
func (t PeerSet) Peers() zap.List { return t.o.List(peerSetPeersOff) }

// PeerSetInput collects the field values for NewPeerSet. Peers is a list of
// pre-built PeerInfo sub-buffers (build each with NewPeerInfo).
type PeerSetInput struct {
	Peers [][]byte
}

// NewPeerSet builds a ZAP-encoded PeerSet message from in and returns the bytes.
func NewPeerSet(in PeerSetInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(peerSetSize)
	peersLB := b.StartList(0)
	for _, elem := range in.Peers {
		peersLB.AddObjectBytes(elem)
	}
	ob.SetList(peerSetPeersOff, peersLB.FinishOffset(), len(in.Peers))
	ob.FinishAsRoot()
	return b.Finish()
}

// --- PeerInfo ---

const (
	peerInfoPeerAddrOff   = 0
	peerInfoRackOff       = 8
	peerInfoDataCenterOff = 16
	peerInfoSize          = 24
)

// PeerInfo is a zero-copy view into a ZAP-encoded PeerInfo message.
type PeerInfo struct{ o zap.Object }

// WrapPeerInfo parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapPeerInfo(b []byte) (PeerInfo, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PeerInfo{}, err
	}
	return PeerInfo{o: m.Root()}, nil
}

// PeerAddr reads the peer_addr field (proto field 1, string).
func (t PeerInfo) PeerAddr() string { return t.o.Text(peerInfoPeerAddrOff) }

// Rack reads the rack field (proto field 2, string).
func (t PeerInfo) Rack() string { return t.o.Text(peerInfoRackOff) }

// DataCenter reads the data_center field (proto field 3, string).
func (t PeerInfo) DataCenter() string { return t.o.Text(peerInfoDataCenterOff) }

// PeerInfoInput collects the field values for NewPeerInfo.
type PeerInfoInput struct {
	PeerAddr   string
	Rack       string
	DataCenter string
}

// NewPeerInfo builds a ZAP-encoded PeerInfo message from in and returns the
// bytes.
func NewPeerInfo(in PeerInfoInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(peerInfoSize)
	ob.SetText(peerInfoPeerAddrOff, in.PeerAddr)
	ob.SetText(peerInfoRackOff, in.Rack)
	ob.SetText(peerInfoDataCenterOff, in.DataCenter)
	ob.FinishAsRoot()
	return b.Finish()
}
