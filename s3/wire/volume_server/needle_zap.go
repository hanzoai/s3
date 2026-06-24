// Code generated from volume_server.proto; DO NOT EDIT.

package volume_serverwire

import (
	zap "github.com/zap-proto/go"
)

// --- ReadNeedleBlobRequest ---
//
// Proto field 2 is absent (skipped in the .proto); present fields take slots in
// declaration order.

const (
	readNeedleBlobRequestVolumeIDOff = 0
	readNeedleBlobRequestOffsetOff   = 8
	readNeedleBlobRequestSizeFOff    = 16
	readNeedleBlobRequestSize        = 20
)

// ReadNeedleBlobRequest is a zero-copy view into a ZAP-encoded
// ReadNeedleBlobRequest message.
type ReadNeedleBlobRequest struct{ o zap.Object }

// WrapReadNeedleBlobRequest parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapReadNeedleBlobRequest(b []byte) (ReadNeedleBlobRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ReadNeedleBlobRequest{}, err
	}
	return ReadNeedleBlobRequest{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t ReadNeedleBlobRequest) VolumeID() uint32 {
	return t.o.Uint32(readNeedleBlobRequestVolumeIDOff)
}

// Offset reads the offset field (proto field 3, int64).
func (t ReadNeedleBlobRequest) Offset() int64 { return t.o.Int64(readNeedleBlobRequestOffsetOff) }

// Size reads the size field (proto field 4, int32).
func (t ReadNeedleBlobRequest) Size() int32 { return t.o.Int32(readNeedleBlobRequestSizeFOff) }

// ReadNeedleBlobRequestInput collects the field values for
// NewReadNeedleBlobRequest.
type ReadNeedleBlobRequestInput struct {
	VolumeID uint32
	Offset   int64
	Size     int32
}

// NewReadNeedleBlobRequest builds a ZAP-encoded ReadNeedleBlobRequest message
// from in and returns the bytes.
func NewReadNeedleBlobRequest(in ReadNeedleBlobRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(readNeedleBlobRequestSize)
	ob.SetUint32(readNeedleBlobRequestVolumeIDOff, in.VolumeID)
	ob.SetInt64(readNeedleBlobRequestOffsetOff, in.Offset)
	ob.SetInt32(readNeedleBlobRequestSizeFOff, in.Size)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ReadNeedleBlobResponse ---

const (
	readNeedleBlobResponseNeedleBlobOff = 0
	readNeedleBlobResponseSize          = 8
)

// ReadNeedleBlobResponse is a zero-copy view into a ZAP-encoded
// ReadNeedleBlobResponse message.
type ReadNeedleBlobResponse struct{ o zap.Object }

// WrapReadNeedleBlobResponse parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapReadNeedleBlobResponse(b []byte) (ReadNeedleBlobResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ReadNeedleBlobResponse{}, err
	}
	return ReadNeedleBlobResponse{o: m.Root()}, nil
}

// NeedleBlob reads the needle_blob field (proto field 1, bytes).
func (t ReadNeedleBlobResponse) NeedleBlob() []byte {
	return t.o.Bytes(readNeedleBlobResponseNeedleBlobOff)
}

// ReadNeedleBlobResponseInput collects the field values for
// NewReadNeedleBlobResponse.
type ReadNeedleBlobResponseInput struct {
	NeedleBlob []byte
}

// NewReadNeedleBlobResponse builds a ZAP-encoded ReadNeedleBlobResponse message
// from in and returns the bytes.
func NewReadNeedleBlobResponse(in ReadNeedleBlobResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(readNeedleBlobResponseSize)
	ob.SetBytes(readNeedleBlobResponseNeedleBlobOff, in.NeedleBlob)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ReadNeedleMetaRequest ---

const (
	readNeedleMetaRequestVolumeIDOff = 0
	readNeedleMetaRequestNeedleIDOff = 8
	readNeedleMetaRequestOffsetOff   = 16
	readNeedleMetaRequestSizeFOff    = 24
	readNeedleMetaRequestSize        = 28
)

// ReadNeedleMetaRequest is a zero-copy view into a ZAP-encoded
// ReadNeedleMetaRequest message.
type ReadNeedleMetaRequest struct{ o zap.Object }

// WrapReadNeedleMetaRequest parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapReadNeedleMetaRequest(b []byte) (ReadNeedleMetaRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ReadNeedleMetaRequest{}, err
	}
	return ReadNeedleMetaRequest{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t ReadNeedleMetaRequest) VolumeID() uint32 {
	return t.o.Uint32(readNeedleMetaRequestVolumeIDOff)
}

// NeedleID reads the needle_id field (proto field 2, uint64).
func (t ReadNeedleMetaRequest) NeedleID() uint64 {
	return t.o.Uint64(readNeedleMetaRequestNeedleIDOff)
}

// Offset reads the offset field (proto field 3, int64).
func (t ReadNeedleMetaRequest) Offset() int64 { return t.o.Int64(readNeedleMetaRequestOffsetOff) }

// Size reads the size field (proto field 4, int32).
func (t ReadNeedleMetaRequest) Size() int32 { return t.o.Int32(readNeedleMetaRequestSizeFOff) }

// ReadNeedleMetaRequestInput collects the field values for
// NewReadNeedleMetaRequest.
type ReadNeedleMetaRequestInput struct {
	VolumeID uint32
	NeedleID uint64
	Offset   int64
	Size     int32
}

// NewReadNeedleMetaRequest builds a ZAP-encoded ReadNeedleMetaRequest message
// from in and returns the bytes.
func NewReadNeedleMetaRequest(in ReadNeedleMetaRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(readNeedleMetaRequestSize)
	ob.SetUint32(readNeedleMetaRequestVolumeIDOff, in.VolumeID)
	ob.SetUint64(readNeedleMetaRequestNeedleIDOff, in.NeedleID)
	ob.SetInt64(readNeedleMetaRequestOffsetOff, in.Offset)
	ob.SetInt32(readNeedleMetaRequestSizeFOff, in.Size)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ReadNeedleMetaResponse ---

const (
	readNeedleMetaResponseCookieOff       = 0
	readNeedleMetaResponseLastModifiedOff = 8
	readNeedleMetaResponseCrcOff          = 16
	readNeedleMetaResponseTTLOff          = 24
	readNeedleMetaResponseAppendAtNsOff   = 32
	readNeedleMetaResponseSize            = 40
)

// ReadNeedleMetaResponse is a zero-copy view into a ZAP-encoded
// ReadNeedleMetaResponse message.
type ReadNeedleMetaResponse struct{ o zap.Object }

// WrapReadNeedleMetaResponse parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapReadNeedleMetaResponse(b []byte) (ReadNeedleMetaResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ReadNeedleMetaResponse{}, err
	}
	return ReadNeedleMetaResponse{o: m.Root()}, nil
}

// Cookie reads the cookie field (proto field 1, uint32).
func (t ReadNeedleMetaResponse) Cookie() uint32 { return t.o.Uint32(readNeedleMetaResponseCookieOff) }

// LastModified reads the last_modified field (proto field 2, uint64).
func (t ReadNeedleMetaResponse) LastModified() uint64 {
	return t.o.Uint64(readNeedleMetaResponseLastModifiedOff)
}

// Crc reads the crc field (proto field 3, uint32).
func (t ReadNeedleMetaResponse) Crc() uint32 { return t.o.Uint32(readNeedleMetaResponseCrcOff) }

// TTL reads the ttl field (proto field 4, string).
func (t ReadNeedleMetaResponse) TTL() string { return t.o.Text(readNeedleMetaResponseTTLOff) }

// AppendAtNs reads the append_at_ns field (proto field 5, uint64).
func (t ReadNeedleMetaResponse) AppendAtNs() uint64 {
	return t.o.Uint64(readNeedleMetaResponseAppendAtNsOff)
}

// ReadNeedleMetaResponseInput collects the field values for
// NewReadNeedleMetaResponse.
type ReadNeedleMetaResponseInput struct {
	Cookie       uint32
	LastModified uint64
	Crc          uint32
	TTL          string
	AppendAtNs   uint64
}

// NewReadNeedleMetaResponse builds a ZAP-encoded ReadNeedleMetaResponse message
// from in and returns the bytes.
func NewReadNeedleMetaResponse(in ReadNeedleMetaResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(readNeedleMetaResponseSize)
	ob.SetUint32(readNeedleMetaResponseCookieOff, in.Cookie)
	ob.SetUint64(readNeedleMetaResponseLastModifiedOff, in.LastModified)
	ob.SetUint32(readNeedleMetaResponseCrcOff, in.Crc)
	ob.SetText(readNeedleMetaResponseTTLOff, in.TTL)
	ob.SetUint64(readNeedleMetaResponseAppendAtNsOff, in.AppendAtNs)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- WriteNeedleBlobRequest ---

const (
	writeNeedleBlobRequestVolumeIDOff   = 0
	writeNeedleBlobRequestNeedleIDOff   = 8
	writeNeedleBlobRequestSizeFOff      = 16
	writeNeedleBlobRequestNeedleBlobOff = 20
	writeNeedleBlobRequestSize          = 28
)

// WriteNeedleBlobRequest is a zero-copy view into a ZAP-encoded
// WriteNeedleBlobRequest message.
type WriteNeedleBlobRequest struct{ o zap.Object }

// WrapWriteNeedleBlobRequest parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapWriteNeedleBlobRequest(b []byte) (WriteNeedleBlobRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return WriteNeedleBlobRequest{}, err
	}
	return WriteNeedleBlobRequest{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t WriteNeedleBlobRequest) VolumeID() uint32 {
	return t.o.Uint32(writeNeedleBlobRequestVolumeIDOff)
}

// NeedleID reads the needle_id field (proto field 2, uint64).
func (t WriteNeedleBlobRequest) NeedleID() uint64 {
	return t.o.Uint64(writeNeedleBlobRequestNeedleIDOff)
}

// Size reads the size field (proto field 3, int32).
func (t WriteNeedleBlobRequest) Size() int32 { return t.o.Int32(writeNeedleBlobRequestSizeFOff) }

// NeedleBlob reads the needle_blob field (proto field 4, bytes).
func (t WriteNeedleBlobRequest) NeedleBlob() []byte {
	return t.o.Bytes(writeNeedleBlobRequestNeedleBlobOff)
}

// WriteNeedleBlobRequestInput collects the field values for
// NewWriteNeedleBlobRequest.
type WriteNeedleBlobRequestInput struct {
	VolumeID   uint32
	NeedleID   uint64
	Size       int32
	NeedleBlob []byte
}

// NewWriteNeedleBlobRequest builds a ZAP-encoded WriteNeedleBlobRequest message
// from in and returns the bytes.
func NewWriteNeedleBlobRequest(in WriteNeedleBlobRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(writeNeedleBlobRequestSize)
	ob.SetUint32(writeNeedleBlobRequestVolumeIDOff, in.VolumeID)
	ob.SetUint64(writeNeedleBlobRequestNeedleIDOff, in.NeedleID)
	ob.SetInt32(writeNeedleBlobRequestSizeFOff, in.Size)
	ob.SetBytes(writeNeedleBlobRequestNeedleBlobOff, in.NeedleBlob)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- WriteNeedleBlobResponse ---

const writeNeedleBlobResponseSize = 0

// WriteNeedleBlobResponse is a zero-copy view into a ZAP-encoded
// WriteNeedleBlobResponse message. WriteNeedleBlobResponse carries no fields.
type WriteNeedleBlobResponse struct{ o zap.Object }

// WrapWriteNeedleBlobResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapWriteNeedleBlobResponse(b []byte) (WriteNeedleBlobResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return WriteNeedleBlobResponse{}, err
	}
	return WriteNeedleBlobResponse{o: m.Root()}, nil
}

// WriteNeedleBlobResponseInput collects the field values for
// NewWriteNeedleBlobResponse. WriteNeedleBlobResponse has no fields.
type WriteNeedleBlobResponseInput struct{}

// NewWriteNeedleBlobResponse builds a ZAP-encoded WriteNeedleBlobResponse
// message and returns the bytes.
func NewWriteNeedleBlobResponse(in WriteNeedleBlobResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(writeNeedleBlobResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ReadAllNeedlesRequest ---

const (
	readAllNeedlesRequestVolumeIDsOff = 0
	readAllNeedlesRequestSize         = 8
)

// ReadAllNeedlesRequest is a zero-copy view into a ZAP-encoded
// ReadAllNeedlesRequest message.
type ReadAllNeedlesRequest struct{ o zap.Object }

// WrapReadAllNeedlesRequest parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapReadAllNeedlesRequest(b []byte) (ReadAllNeedlesRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ReadAllNeedlesRequest{}, err
	}
	return ReadAllNeedlesRequest{o: m.Root()}, nil
}

// VolumeIDsLen reports the number of volume_ids elements (proto field 1,
// repeated uint32).
func (t ReadAllNeedlesRequest) VolumeIDsLen() int {
	return t.o.List(readAllNeedlesRequestVolumeIDsOff).Len()
}

// VolumeIDAt returns the i-th volume_ids element (uint32).
func (t ReadAllNeedlesRequest) VolumeIDAt(i int) uint32 {
	return t.o.List(readAllNeedlesRequestVolumeIDsOff).Uint32(i)
}

// ReadAllNeedlesRequestInput collects the field values for
// NewReadAllNeedlesRequest.
type ReadAllNeedlesRequestInput struct {
	VolumeIDs []uint32
}

// NewReadAllNeedlesRequest builds a ZAP-encoded ReadAllNeedlesRequest message
// from in and returns the bytes.
func NewReadAllNeedlesRequest(in ReadAllNeedlesRequestInput) []byte {
	b := zap.NewBuilder(256)
	listOff, listLen := addUint32List(b, in.VolumeIDs)
	ob := b.StartObject(readAllNeedlesRequestSize)
	ob.SetList(readAllNeedlesRequestVolumeIDsOff, listOff, listLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ReadAllNeedlesResponse ---
//
// Proto field 4 is absent (skipped in the .proto); present fields take slots in
// declaration order.

const (
	readAllNeedlesResponseVolumeIDOff             = 0
	readAllNeedlesResponseNeedleIDOff             = 8
	readAllNeedlesResponseCookieOff               = 16
	readAllNeedlesResponseNeedleBlobOff           = 24
	readAllNeedlesResponseNeedleBlobCompressedOff = 32
	readAllNeedlesResponseLastModifiedOff         = 40
	readAllNeedlesResponseCrcOff                  = 48
	readAllNeedlesResponseNameOff                 = 56
	readAllNeedlesResponseMimeOff                 = 64
	readAllNeedlesResponseSize                    = 72
)

// ReadAllNeedlesResponse is a zero-copy view into a ZAP-encoded
// ReadAllNeedlesResponse message.
type ReadAllNeedlesResponse struct{ o zap.Object }

// WrapReadAllNeedlesResponse parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapReadAllNeedlesResponse(b []byte) (ReadAllNeedlesResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ReadAllNeedlesResponse{}, err
	}
	return ReadAllNeedlesResponse{o: m.Root()}, nil
}

// VolumeID reads the volume_id field (proto field 1, uint32).
func (t ReadAllNeedlesResponse) VolumeID() uint32 {
	return t.o.Uint32(readAllNeedlesResponseVolumeIDOff)
}

// NeedleID reads the needle_id field (proto field 2, uint64).
func (t ReadAllNeedlesResponse) NeedleID() uint64 {
	return t.o.Uint64(readAllNeedlesResponseNeedleIDOff)
}

// Cookie reads the cookie field (proto field 3, uint32).
func (t ReadAllNeedlesResponse) Cookie() uint32 {
	return t.o.Uint32(readAllNeedlesResponseCookieOff)
}

// NeedleBlob reads the needle_blob field (proto field 5, bytes).
func (t ReadAllNeedlesResponse) NeedleBlob() []byte {
	return t.o.Bytes(readAllNeedlesResponseNeedleBlobOff)
}

// NeedleBlobCompressed reads the needle_blob_compressed field (proto field 6,
// bool).
func (t ReadAllNeedlesResponse) NeedleBlobCompressed() bool {
	return t.o.Bool(readAllNeedlesResponseNeedleBlobCompressedOff)
}

// LastModified reads the last_modified field (proto field 7, uint64).
func (t ReadAllNeedlesResponse) LastModified() uint64 {
	return t.o.Uint64(readAllNeedlesResponseLastModifiedOff)
}

// Crc reads the crc field (proto field 8, uint32).
func (t ReadAllNeedlesResponse) Crc() uint32 { return t.o.Uint32(readAllNeedlesResponseCrcOff) }

// Name reads the name field (proto field 9, bytes).
func (t ReadAllNeedlesResponse) Name() []byte { return t.o.Bytes(readAllNeedlesResponseNameOff) }

// Mime reads the mime field (proto field 10, bytes).
func (t ReadAllNeedlesResponse) Mime() []byte { return t.o.Bytes(readAllNeedlesResponseMimeOff) }

// ReadAllNeedlesResponseInput collects the field values for
// NewReadAllNeedlesResponse.
type ReadAllNeedlesResponseInput struct {
	VolumeID             uint32
	NeedleID             uint64
	Cookie               uint32
	NeedleBlob           []byte
	NeedleBlobCompressed bool
	LastModified         uint64
	Crc                  uint32
	Name                 []byte
	Mime                 []byte
}

// NewReadAllNeedlesResponse builds a ZAP-encoded ReadAllNeedlesResponse message
// from in and returns the bytes.
func NewReadAllNeedlesResponse(in ReadAllNeedlesResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(readAllNeedlesResponseSize)
	ob.SetUint32(readAllNeedlesResponseVolumeIDOff, in.VolumeID)
	ob.SetUint64(readAllNeedlesResponseNeedleIDOff, in.NeedleID)
	ob.SetUint32(readAllNeedlesResponseCookieOff, in.Cookie)
	ob.SetBytes(readAllNeedlesResponseNeedleBlobOff, in.NeedleBlob)
	ob.SetBool(readAllNeedlesResponseNeedleBlobCompressedOff, in.NeedleBlobCompressed)
	ob.SetUint64(readAllNeedlesResponseLastModifiedOff, in.LastModified)
	ob.SetUint32(readAllNeedlesResponseCrcOff, in.Crc)
	ob.SetBytes(readAllNeedlesResponseNameOff, in.Name)
	ob.SetBytes(readAllNeedlesResponseMimeOff, in.Mime)
	ob.FinishAsRoot()
	return b.Finish()
}
