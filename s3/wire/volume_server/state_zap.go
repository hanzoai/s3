// Code generated from volume_server.proto; DO NOT EDIT.

package volume_serverwire

import (
	zap "github.com/zap-proto/go"
)

// --- VolumeServerState ---

const (
	volumeServerStateMaintenanceOff = 0
	volumeServerStateVersionOff     = 4
	volumeServerStateSize           = 8
)

// VolumeServerState is a zero-copy view into a ZAP-encoded VolumeServerState
// message.
type VolumeServerState struct{ o zap.Object }

// WrapVolumeServerState parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapVolumeServerState(b []byte) (VolumeServerState, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeServerState{}, err
	}
	return VolumeServerState{o: m.Root()}, nil
}

// Maintenance reads the maintenance field (proto field 1, bool).
func (t VolumeServerState) Maintenance() bool { return t.o.Bool(volumeServerStateMaintenanceOff) }

// Version reads the version field (proto field 2, uint32).
func (t VolumeServerState) Version() uint32 { return t.o.Uint32(volumeServerStateVersionOff) }

// VolumeServerStateInput collects the field values for NewVolumeServerState.
type VolumeServerStateInput struct {
	Maintenance bool
	Version     uint32
}

// NewVolumeServerState builds a ZAP-encoded VolumeServerState message from in
// and returns the bytes.
func NewVolumeServerState(in VolumeServerStateInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeServerStateSize)
	ob.SetBool(volumeServerStateMaintenanceOff, in.Maintenance)
	ob.SetUint32(volumeServerStateVersionOff, in.Version)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- BatchDeleteRequest ---

const (
	batchDeleteRequestFileIDsOff         = 0
	batchDeleteRequestSkipCookieCheckOff = 8
	batchDeleteRequestSize               = 9
)

// BatchDeleteRequest is a zero-copy view into a ZAP-encoded BatchDeleteRequest
// message.
type BatchDeleteRequest struct{ o zap.Object }

// WrapBatchDeleteRequest parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapBatchDeleteRequest(b []byte) (BatchDeleteRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return BatchDeleteRequest{}, err
	}
	return BatchDeleteRequest{o: m.Root()}, nil
}

// FileIDsLen reports the number of file_ids elements (proto field 1, repeated
// string).
func (t BatchDeleteRequest) FileIDsLen() int { return t.o.List(batchDeleteRequestFileIDsOff).Len() }

// FileIDAt returns the i-th file_ids element (string).
func (t BatchDeleteRequest) FileIDAt(i int) string {
	return elemText(t.o.List(batchDeleteRequestFileIDsOff), i)
}

// SkipCookieCheck reads the skip_cookie_check field (proto field 2, bool).
func (t BatchDeleteRequest) SkipCookieCheck() bool {
	return t.o.Bool(batchDeleteRequestSkipCookieCheckOff)
}

// BatchDeleteRequestInput collects the field values for NewBatchDeleteRequest.
type BatchDeleteRequestInput struct {
	FileIDs         []string
	SkipCookieCheck bool
}

// NewBatchDeleteRequest builds a ZAP-encoded BatchDeleteRequest message from in
// and returns the bytes.
func NewBatchDeleteRequest(in BatchDeleteRequestInput) []byte {
	b := zap.NewBuilder(256)
	listOff, listLen := addStringList(b, in.FileIDs)
	ob := b.StartObject(batchDeleteRequestSize)
	ob.SetList(batchDeleteRequestFileIDsOff, listOff, listLen)
	ob.SetBool(batchDeleteRequestSkipCookieCheckOff, in.SkipCookieCheck)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- BatchDeleteResponse ---

const (
	batchDeleteResponseResultsOff = 0
	batchDeleteResponseSize       = 8
)

// BatchDeleteResponse is a zero-copy view into a ZAP-encoded BatchDeleteResponse
// message.
type BatchDeleteResponse struct{ o zap.Object }

// WrapBatchDeleteResponse parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapBatchDeleteResponse(b []byte) (BatchDeleteResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return BatchDeleteResponse{}, err
	}
	return BatchDeleteResponse{o: m.Root()}, nil
}

// ResultsLen reports the number of results elements (proto field 1, repeated
// DeleteResult).
func (t BatchDeleteResponse) ResultsLen() int {
	return t.o.List(batchDeleteResponseResultsOff).Len()
}

// ResultAt returns the i-th results element. The bool is false when i is out of
// range or the element fails to parse.
func (t BatchDeleteResponse) ResultAt(i int) (DeleteResult, bool) {
	o := t.o.List(batchDeleteResponseResultsOff).ObjectAt(i)
	if o.IsNull() {
		return DeleteResult{}, false
	}
	return DeleteResult{o: o}, true
}

// BatchDeleteResponseInput collects the field values for
// NewBatchDeleteResponse. Each Results entry is a DeleteResult message's own ZAP
// buffer (from NewDeleteResult).
type BatchDeleteResponseInput struct {
	Results [][]byte
}

// NewBatchDeleteResponse builds a ZAP-encoded BatchDeleteResponse message from
// in and returns the bytes.
func NewBatchDeleteResponse(in BatchDeleteResponseInput) []byte {
	b := zap.NewBuilder(256)
	listOff, listLen := addObjectList(b, in.Results)
	ob := b.StartObject(batchDeleteResponseSize)
	ob.SetList(batchDeleteResponseResultsOff, listOff, listLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- DeleteResult ---

const (
	deleteResultFileIDOff  = 0
	deleteResultStatusOff  = 8
	deleteResultErrorOff   = 12
	deleteResultSizeOff    = 20
	deleteResultVersionOff = 24
	deleteResultSize       = 28
)

// DeleteResult is a zero-copy view into a ZAP-encoded DeleteResult message.
type DeleteResult struct{ o zap.Object }

// WrapDeleteResult parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapDeleteResult(b []byte) (DeleteResult, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return DeleteResult{}, err
	}
	return DeleteResult{o: m.Root()}, nil
}

// FileID reads the file_id field (proto field 1, string).
func (t DeleteResult) FileID() string { return t.o.Text(deleteResultFileIDOff) }

// Status reads the status field (proto field 2, int32).
func (t DeleteResult) Status() int32 { return t.o.Int32(deleteResultStatusOff) }

// Error reads the error field (proto field 3, string).
func (t DeleteResult) Error() string { return t.o.Text(deleteResultErrorOff) }

// Size reads the size field (proto field 4, uint32).
func (t DeleteResult) Size() uint32 { return t.o.Uint32(deleteResultSizeOff) }

// Version reads the version field (proto field 5, uint32).
func (t DeleteResult) Version() uint32 { return t.o.Uint32(deleteResultVersionOff) }

// DeleteResultInput collects the field values for NewDeleteResult.
type DeleteResultInput struct {
	FileID  string
	Status  int32
	Error   string
	Size    uint32
	Version uint32
}

// NewDeleteResult builds a ZAP-encoded DeleteResult message from in and returns
// the bytes.
func NewDeleteResult(in DeleteResultInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(deleteResultSize)
	ob.SetText(deleteResultFileIDOff, in.FileID)
	ob.SetInt32(deleteResultStatusOff, in.Status)
	ob.SetText(deleteResultErrorOff, in.Error)
	ob.SetUint32(deleteResultSizeOff, in.Size)
	ob.SetUint32(deleteResultVersionOff, in.Version)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- Empty ---

const emptySize = 0

// Empty is a zero-copy view into a ZAP-encoded Empty message. Empty carries no
// fields.
type Empty struct{ o zap.Object }

// WrapEmpty parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapEmpty(b []byte) (Empty, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return Empty{}, err
	}
	return Empty{o: m.Root()}, nil
}

// EmptyInput collects the field values for NewEmpty. Empty has no fields.
type EmptyInput struct{}

// NewEmpty builds a ZAP-encoded Empty message and returns the bytes.
func NewEmpty(in EmptyInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(emptySize)
	ob.FinishAsRoot()
	return b.Finish()
}
