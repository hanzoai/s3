// Code generated from mount_peer.proto; DO NOT EDIT.

package mount_peerwire

import (
	zap "github.com/zap-proto/go"
)

// --- FetchChunkRequest ---

const (
	fetchChunkRequestFileIdOff       = 0
	fetchChunkRequestExpectedEtagOff = 8
	fetchChunkRequestExpectedSizeOff = 16
	fetchChunkRequestOffsetOff       = 24
	fetchChunkRequestLengthOff       = 32
	fetchChunkRequestSize            = 40
)

// FetchChunkRequest is a zero-copy view into a ZAP-encoded FetchChunkRequest
// message.
type FetchChunkRequest struct{ o zap.Object }

// WrapFetchChunkRequest parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapFetchChunkRequest(b []byte) (FetchChunkRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return FetchChunkRequest{}, err
	}
	return FetchChunkRequest{o: m.Root()}, nil
}

// FileId reads the file_id field (proto field 1, string).
func (t FetchChunkRequest) FileId() string { return t.o.Text(fetchChunkRequestFileIdOff) }

// ExpectedEtag reads the expected_etag field (proto field 2, string): the
// caller's expected MD5 over the full chunk.
func (t FetchChunkRequest) ExpectedEtag() string {
	return t.o.Text(fetchChunkRequestExpectedEtagOff)
}

// ExpectedSize reads the expected_size field (proto field 3, uint64): the
// filer-reported chunk byte count.
func (t FetchChunkRequest) ExpectedSize() uint64 {
	return t.o.Uint64(fetchChunkRequestExpectedSizeOff)
}

// Offset reads the offset field (proto field 4, uint64): byte offset within the
// chunk to start reading from.
func (t FetchChunkRequest) Offset() uint64 { return t.o.Uint64(fetchChunkRequestOffsetOff) }

// Length reads the length field (proto field 5, uint64): number of bytes to
// return from offset (0 means until end of chunk).
func (t FetchChunkRequest) Length() uint64 { return t.o.Uint64(fetchChunkRequestLengthOff) }

// FetchChunkRequestInput collects the field values for NewFetchChunkRequest.
type FetchChunkRequestInput struct {
	FileId       string
	ExpectedEtag string
	ExpectedSize uint64
	Offset       uint64
	Length       uint64
}

// NewFetchChunkRequest builds a ZAP-encoded FetchChunkRequest message from in
// and returns the bytes.
func NewFetchChunkRequest(in FetchChunkRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(fetchChunkRequestSize)
	ob.SetText(fetchChunkRequestFileIdOff, in.FileId)
	ob.SetText(fetchChunkRequestExpectedEtagOff, in.ExpectedEtag)
	ob.SetUint64(fetchChunkRequestExpectedSizeOff, in.ExpectedSize)
	ob.SetUint64(fetchChunkRequestOffsetOff, in.Offset)
	ob.SetUint64(fetchChunkRequestLengthOff, in.Length)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- FetchChunkResponse ---
//
// STREAMING: FetchChunk is a server-streaming RPC; the server emits a sequence
// of FetchChunkResponse frames whose data fields concatenate into the full
// chunk. This schema describes one frame. The streaming method body lands when
// the transport streaming primitive ships.

const (
	fetchChunkResponseDataOff = 0
	fetchChunkResponseSize    = 8
)

// FetchChunkResponse is a zero-copy view into a ZAP-encoded FetchChunkResponse
// message (one stream frame).
type FetchChunkResponse struct{ o zap.Object }

// WrapFetchChunkResponse parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapFetchChunkResponse(b []byte) (FetchChunkResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return FetchChunkResponse{}, err
	}
	return FetchChunkResponse{o: m.Root()}, nil
}

// Data reads the data field (proto field 1, bytes): the next frame of chunk
// bytes; concatenate across the stream.
func (t FetchChunkResponse) Data() []byte { return t.o.Bytes(fetchChunkResponseDataOff) }

// FetchChunkResponseInput collects the field values for NewFetchChunkResponse.
type FetchChunkResponseInput struct {
	Data []byte
}

// NewFetchChunkResponse builds a ZAP-encoded FetchChunkResponse message from in
// and returns the bytes.
func NewFetchChunkResponse(in FetchChunkResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(fetchChunkResponseSize)
	ob.SetBytes(fetchChunkResponseDataOff, in.Data)
	ob.FinishAsRoot()
	return b.Finish()
}
