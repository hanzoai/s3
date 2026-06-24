// Code generated from filer.proto; DO NOT EDIT.

package filerwire

import (
	zap "github.com/zap-proto/go"
)

// --- CreateEntryRequest ---

const (
	createEntryRequestDirectoryOff                = 0  // directory (field 1, string)
	createEntryRequestEntryOff                    = 8  // entry (field 2, Entry object)
	createEntryRequestOExclOff                    = 12 // o_excl (field 3, bool)
	createEntryRequestIsFromOtherClusterOff       = 13 // is_from_other_cluster (field 4, bool)
	createEntryRequestSignaturesOff               = 16 // signatures (field 5, repeated int32)
	createEntryRequestSkipCheckParentDirectoryOff = 24 // skip_check_parent_directory (field 6, bool)
	createEntryRequestConditionOff                = 28 // condition (field 7, WriteCondition object)
	createEntryRequestSize                        = 32
)

// CreateEntryRequest is a zero-copy view into a ZAP-encoded CreateEntryRequest message.
type CreateEntryRequest struct{ o zap.Object }

// WrapCreateEntryRequest parses b and returns a typed view.
func WrapCreateEntryRequest(b []byte) (CreateEntryRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return CreateEntryRequest{}, err
	}
	return CreateEntryRequest{o: m.Root()}, nil
}

func (t CreateEntryRequest) Directory() string { return t.o.Text(createEntryRequestDirectoryOff) }
func (t CreateEntryRequest) OExcl() bool       { return t.o.Bool(createEntryRequestOExclOff) }
func (t CreateEntryRequest) IsFromOtherCluster() bool {
	return t.o.Bool(createEntryRequestIsFromOtherClusterOff)
}
func (t CreateEntryRequest) SkipCheckParentDirectory() bool {
	return t.o.Bool(createEntryRequestSkipCheckParentDirectoryOff)
}

// Entry returns the nested Entry (null-object if unset).
func (t CreateEntryRequest) Entry() Entry {
	return Entry{o: t.o.Object(createEntryRequestEntryOff)}
}

// HasEntry reports whether the entry field is present.
func (t CreateEntryRequest) HasEntry() bool {
	return !t.o.Object(createEntryRequestEntryOff).IsNull()
}

// Condition returns the nested WriteCondition (null-object if unset).
func (t CreateEntryRequest) Condition() WriteCondition {
	return WriteCondition{o: t.o.Object(createEntryRequestConditionOff)}
}

// HasCondition reports whether the condition field is present.
func (t CreateEntryRequest) HasCondition() bool {
	return !t.o.Object(createEntryRequestConditionOff).IsNull()
}

// SignaturesLen returns the number of signatures elements (proto field 5,
// repeated int32).
func (t CreateEntryRequest) SignaturesLen() int {
	return t.o.List(createEntryRequestSignaturesOff).Len()
}

// Signature returns the i-th signatures element.
func (t CreateEntryRequest) Signature(i int) int32 {
	return int32(t.o.List(createEntryRequestSignaturesOff).Uint32(i))
}

// CreateEntryRequestInput collects the field values for NewCreateEntryRequest.
// Entry and Condition, when non-nil, are standalone buffers.
type CreateEntryRequestInput struct {
	Directory                string
	Entry                    []byte
	OExcl                    bool
	IsFromOtherCluster       bool
	Signatures               []int32
	SkipCheckParentDirectory bool
	Condition                []byte
}

// NewCreateEntryRequest builds a ZAP-encoded CreateEntryRequest message from in.
func NewCreateEntryRequest(in CreateEntryRequestInput) []byte {
	b := zap.NewBuilder(512)
	sigOff, sigLen := addInt32List(b, in.Signatures)
	ob := b.StartObject(createEntryRequestSize)
	ob.SetText(createEntryRequestDirectoryOff, in.Directory)
	setNestedObject(b, ob, createEntryRequestEntryOff, in.Entry)
	ob.SetBool(createEntryRequestOExclOff, in.OExcl)
	ob.SetBool(createEntryRequestIsFromOtherClusterOff, in.IsFromOtherCluster)
	ob.SetList(createEntryRequestSignaturesOff, sigOff, sigLen)
	ob.SetBool(createEntryRequestSkipCheckParentDirectoryOff, in.SkipCheckParentDirectory)
	setNestedObject(b, ob, createEntryRequestConditionOff, in.Condition)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- CreateEntryResponse ---

const (
	createEntryResponseErrorOff         = 0  // error (field 1, string)
	createEntryResponseMetadataEventOff = 8  // metadata_event (field 2, SubscribeMetadataResponse object)
	createEntryResponseErrorCodeOff     = 12 // error_code (field 3, enum FilerError)
	createEntryResponseSize             = 16
)

// CreateEntryResponse is a zero-copy view into a ZAP-encoded CreateEntryResponse message.
type CreateEntryResponse struct{ o zap.Object }

// WrapCreateEntryResponse parses b and returns a typed view.
func WrapCreateEntryResponse(b []byte) (CreateEntryResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return CreateEntryResponse{}, err
	}
	return CreateEntryResponse{o: m.Root()}, nil
}

func (t CreateEntryResponse) Error() string { return t.o.Text(createEntryResponseErrorOff) }

// ErrorCode reads the error_code enum (see FilerError* constants).
func (t CreateEntryResponse) ErrorCode() uint32 {
	return t.o.Uint32(createEntryResponseErrorCodeOff)
}

// MetadataEvent returns the nested SubscribeMetadataResponse (null-object if unset).
func (t CreateEntryResponse) MetadataEvent() SubscribeMetadataResponse {
	return SubscribeMetadataResponse{o: t.o.Object(createEntryResponseMetadataEventOff)}
}

// HasMetadataEvent reports whether the metadata_event field is present.
func (t CreateEntryResponse) HasMetadataEvent() bool {
	return !t.o.Object(createEntryResponseMetadataEventOff).IsNull()
}

// CreateEntryResponseInput collects the field values for NewCreateEntryResponse.
// MetadataEvent, when non-nil, is a standalone SubscribeMetadataResponse buffer.
type CreateEntryResponseInput struct {
	Error         string
	MetadataEvent []byte
	ErrorCode     uint32
}

// NewCreateEntryResponse builds a ZAP-encoded CreateEntryResponse message from in.
func NewCreateEntryResponse(in CreateEntryResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(createEntryResponseSize)
	ob.SetText(createEntryResponseErrorOff, in.Error)
	setNestedObject(b, ob, createEntryResponseMetadataEventOff, in.MetadataEvent)
	ob.SetUint32(createEntryResponseErrorCodeOff, in.ErrorCode)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- UpdateEntryRequestExpectedExtendedEntry (map<string,bytes>) ---

const (
	updateEntryRequestExpectedExtendedEntryKeyOff   = 0
	updateEntryRequestExpectedExtendedEntryValueOff = 8
	updateEntryRequestExpectedExtendedEntrySize     = 16
)

// UpdateEntryRequestExpectedExtendedEntry is a zero-copy view of one
// map<string,bytes> pair of UpdateEntryRequest.expected_extended (proto field 5).
type UpdateEntryRequestExpectedExtendedEntry struct{ o zap.Object }

// WrapUpdateEntryRequestExpectedExtendedEntry parses b and returns a typed view.
func WrapUpdateEntryRequestExpectedExtendedEntry(b []byte) (UpdateEntryRequestExpectedExtendedEntry, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return UpdateEntryRequestExpectedExtendedEntry{}, err
	}
	return UpdateEntryRequestExpectedExtendedEntry{o: m.Root()}, nil
}

func (t UpdateEntryRequestExpectedExtendedEntry) Key() string {
	return t.o.Text(updateEntryRequestExpectedExtendedEntryKeyOff)
}
func (t UpdateEntryRequestExpectedExtendedEntry) Value() []byte {
	return t.o.Bytes(updateEntryRequestExpectedExtendedEntryValueOff)
}

// UpdateEntryRequestExpectedExtendedEntryInput collects one map<string,bytes> pair.
type UpdateEntryRequestExpectedExtendedEntryInput struct {
	Key   string
	Value []byte
}

// NewUpdateEntryRequestExpectedExtendedEntry builds one encoded expected_extended
// map entry.
func NewUpdateEntryRequestExpectedExtendedEntry(in UpdateEntryRequestExpectedExtendedEntryInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(updateEntryRequestExpectedExtendedEntrySize)
	ob.SetText(updateEntryRequestExpectedExtendedEntryKeyOff, in.Key)
	ob.SetBytes(updateEntryRequestExpectedExtendedEntryValueOff, in.Value)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- UpdateEntryRequest ---

const (
	updateEntryRequestDirectoryOff          = 0  // directory (field 1, string)
	updateEntryRequestEntryOff              = 8  // entry (field 2, Entry object)
	updateEntryRequestIsFromOtherClusterOff = 12 // is_from_other_cluster (field 3, bool)
	updateEntryRequestSignaturesOff         = 16 // signatures (field 4, repeated int32)
	updateEntryRequestExpectedExtendedOff   = 24 // expected_extended (field 5, map<string,bytes>)
	updateEntryRequestSize                  = 32
)

// UpdateEntryRequest is a zero-copy view into a ZAP-encoded UpdateEntryRequest message.
type UpdateEntryRequest struct{ o zap.Object }

// WrapUpdateEntryRequest parses b and returns a typed view.
func WrapUpdateEntryRequest(b []byte) (UpdateEntryRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return UpdateEntryRequest{}, err
	}
	return UpdateEntryRequest{o: m.Root()}, nil
}

func (t UpdateEntryRequest) Directory() string { return t.o.Text(updateEntryRequestDirectoryOff) }
func (t UpdateEntryRequest) IsFromOtherCluster() bool {
	return t.o.Bool(updateEntryRequestIsFromOtherClusterOff)
}

// Entry returns the nested Entry (null-object if unset).
func (t UpdateEntryRequest) Entry() Entry {
	return Entry{o: t.o.Object(updateEntryRequestEntryOff)}
}

// HasEntry reports whether the entry field is present.
func (t UpdateEntryRequest) HasEntry() bool {
	return !t.o.Object(updateEntryRequestEntryOff).IsNull()
}

// SignaturesLen returns the number of signatures elements (proto field 4,
// repeated int32).
func (t UpdateEntryRequest) SignaturesLen() int {
	return t.o.List(updateEntryRequestSignaturesOff).Len()
}

// Signature returns the i-th signatures element.
func (t UpdateEntryRequest) Signature(i int) int32 {
	return int32(t.o.List(updateEntryRequestSignaturesOff).Uint32(i))
}

// ExpectedExtendedLen returns the number of expected_extended map entries (proto
// field 5, map<string,bytes>).
func (t UpdateEntryRequest) ExpectedExtendedLen() int {
	return t.o.List(updateEntryRequestExpectedExtendedOff).Len()
}

// ExpectedExtended returns the i-th expected_extended map entry.
func (t UpdateEntryRequest) ExpectedExtended(i int) UpdateEntryRequestExpectedExtendedEntry {
	return UpdateEntryRequestExpectedExtendedEntry{o: t.o.List(updateEntryRequestExpectedExtendedOff).ObjectAt(i)}
}

// UpdateEntryRequestInput collects the field values for NewUpdateEntryRequest.
// Entry, when non-nil, is a standalone Entry buffer; each ExpectedExtended
// element is a pre-built UpdateEntryRequestExpectedExtendedEntry buffer.
type UpdateEntryRequestInput struct {
	Directory          string
	Entry              []byte
	IsFromOtherCluster bool
	Signatures         []int32
	ExpectedExtended   [][]byte
}

// NewUpdateEntryRequest builds a ZAP-encoded UpdateEntryRequest message from in.
func NewUpdateEntryRequest(in UpdateEntryRequestInput) []byte {
	b := zap.NewBuilder(512)
	sigOff, sigLen := addInt32List(b, in.Signatures)
	expOff, expLen := addObjectList(b, in.ExpectedExtended)
	ob := b.StartObject(updateEntryRequestSize)
	ob.SetText(updateEntryRequestDirectoryOff, in.Directory)
	setNestedObject(b, ob, updateEntryRequestEntryOff, in.Entry)
	ob.SetBool(updateEntryRequestIsFromOtherClusterOff, in.IsFromOtherCluster)
	ob.SetList(updateEntryRequestSignaturesOff, sigOff, sigLen)
	ob.SetList(updateEntryRequestExpectedExtendedOff, expOff, expLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- UpdateEntryResponse ---

const (
	updateEntryResponseMetadataEventOff = 0 // metadata_event (field 1, SubscribeMetadataResponse object)
	updateEntryResponseSize             = 8
)

// UpdateEntryResponse is a zero-copy view into a ZAP-encoded UpdateEntryResponse message.
type UpdateEntryResponse struct{ o zap.Object }

// WrapUpdateEntryResponse parses b and returns a typed view.
func WrapUpdateEntryResponse(b []byte) (UpdateEntryResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return UpdateEntryResponse{}, err
	}
	return UpdateEntryResponse{o: m.Root()}, nil
}

// MetadataEvent returns the nested SubscribeMetadataResponse (null-object if unset).
func (t UpdateEntryResponse) MetadataEvent() SubscribeMetadataResponse {
	return SubscribeMetadataResponse{o: t.o.Object(updateEntryResponseMetadataEventOff)}
}

// HasMetadataEvent reports whether the metadata_event field is present.
func (t UpdateEntryResponse) HasMetadataEvent() bool {
	return !t.o.Object(updateEntryResponseMetadataEventOff).IsNull()
}

// UpdateEntryResponseInput collects the field values for NewUpdateEntryResponse.
// MetadataEvent, when non-nil, is a standalone SubscribeMetadataResponse buffer.
type UpdateEntryResponseInput struct {
	MetadataEvent []byte
}

// NewUpdateEntryResponse builds a ZAP-encoded UpdateEntryResponse message from in.
func NewUpdateEntryResponse(in UpdateEntryResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(updateEntryResponseSize)
	setNestedObject(b, ob, updateEntryResponseMetadataEventOff, in.MetadataEvent)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- TouchAccessTimeRequest ---

const (
	touchAccessTimeRequestDirectoryOff     = 0  // directory (field 1, string)
	touchAccessTimeRequestNameOff          = 8  // name (field 2, string)
	touchAccessTimeRequestClientAtimeNsOff = 16 // client_atime_ns (field 3, int64)
	touchAccessTimeRequestSize             = 24
)

// TouchAccessTimeRequest is a zero-copy view into a ZAP-encoded TouchAccessTimeRequest message.
type TouchAccessTimeRequest struct{ o zap.Object }

// WrapTouchAccessTimeRequest parses b and returns a typed view.
func WrapTouchAccessTimeRequest(b []byte) (TouchAccessTimeRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return TouchAccessTimeRequest{}, err
	}
	return TouchAccessTimeRequest{o: m.Root()}, nil
}

func (t TouchAccessTimeRequest) Directory() string {
	return t.o.Text(touchAccessTimeRequestDirectoryOff)
}
func (t TouchAccessTimeRequest) Name() string { return t.o.Text(touchAccessTimeRequestNameOff) }
func (t TouchAccessTimeRequest) ClientAtimeNs() int64 {
	return t.o.Int64(touchAccessTimeRequestClientAtimeNsOff)
}

// TouchAccessTimeRequestInput collects the field values for NewTouchAccessTimeRequest.
type TouchAccessTimeRequestInput struct {
	Directory     string
	Name          string
	ClientAtimeNs int64
}

// NewTouchAccessTimeRequest builds a ZAP-encoded TouchAccessTimeRequest message from in.
func NewTouchAccessTimeRequest(in TouchAccessTimeRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(touchAccessTimeRequestSize)
	ob.SetText(touchAccessTimeRequestDirectoryOff, in.Directory)
	ob.SetText(touchAccessTimeRequestNameOff, in.Name)
	ob.SetInt64(touchAccessTimeRequestClientAtimeNsOff, in.ClientAtimeNs)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- TouchAccessTimeResponse ---

const (
	touchAccessTimeResponsePersistedAtimeNsOff = 0 // persisted_atime_ns (field 1, int64)
	touchAccessTimeResponseUpdatedOff          = 8 // updated (field 2, bool)
	touchAccessTimeResponseSize                = 16
)

// TouchAccessTimeResponse is a zero-copy view into a ZAP-encoded TouchAccessTimeResponse message.
type TouchAccessTimeResponse struct{ o zap.Object }

// WrapTouchAccessTimeResponse parses b and returns a typed view.
func WrapTouchAccessTimeResponse(b []byte) (TouchAccessTimeResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return TouchAccessTimeResponse{}, err
	}
	return TouchAccessTimeResponse{o: m.Root()}, nil
}

func (t TouchAccessTimeResponse) PersistedAtimeNs() int64 {
	return t.o.Int64(touchAccessTimeResponsePersistedAtimeNsOff)
}
func (t TouchAccessTimeResponse) Updated() bool {
	return t.o.Bool(touchAccessTimeResponseUpdatedOff)
}

// TouchAccessTimeResponseInput collects the field values for NewTouchAccessTimeResponse.
type TouchAccessTimeResponseInput struct {
	PersistedAtimeNs int64
	Updated          bool
}

// NewTouchAccessTimeResponse builds a ZAP-encoded TouchAccessTimeResponse message from in.
func NewTouchAccessTimeResponse(in TouchAccessTimeResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(touchAccessTimeResponseSize)
	ob.SetInt64(touchAccessTimeResponsePersistedAtimeNsOff, in.PersistedAtimeNs)
	ob.SetBool(touchAccessTimeResponseUpdatedOff, in.Updated)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- AppendToEntryRequest ---

const (
	appendToEntryRequestDirectoryOff = 0  // directory (field 1, string)
	appendToEntryRequestEntryNameOff = 8  // entry_name (field 2, string)
	appendToEntryRequestChunksOff    = 16 // chunks (field 3, repeated FileChunk)
	appendToEntryRequestSize         = 24
)

// AppendToEntryRequest is a zero-copy view into a ZAP-encoded AppendToEntryRequest message.
type AppendToEntryRequest struct{ o zap.Object }

// WrapAppendToEntryRequest parses b and returns a typed view.
func WrapAppendToEntryRequest(b []byte) (AppendToEntryRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return AppendToEntryRequest{}, err
	}
	return AppendToEntryRequest{o: m.Root()}, nil
}

func (t AppendToEntryRequest) Directory() string {
	return t.o.Text(appendToEntryRequestDirectoryOff)
}
func (t AppendToEntryRequest) EntryName() string {
	return t.o.Text(appendToEntryRequestEntryNameOff)
}

// ChunksLen returns the number of chunks elements (proto field 3, repeated FileChunk).
func (t AppendToEntryRequest) ChunksLen() int { return t.o.List(appendToEntryRequestChunksOff).Len() }

// Chunks returns the i-th chunks element as a typed FileChunk view.
func (t AppendToEntryRequest) Chunks(i int) FileChunk {
	return FileChunk{o: t.o.List(appendToEntryRequestChunksOff).ObjectAt(i)}
}

// AppendToEntryRequestInput collects the field values for NewAppendToEntryRequest.
// Each Chunks element is a standalone FileChunk buffer (from NewFileChunk).
type AppendToEntryRequestInput struct {
	Directory string
	EntryName string
	Chunks    [][]byte
}

// NewAppendToEntryRequest builds a ZAP-encoded AppendToEntryRequest message from in.
func NewAppendToEntryRequest(in AppendToEntryRequestInput) []byte {
	b := zap.NewBuilder(512)
	chunksOff, chunksLen := addObjectList(b, in.Chunks)
	ob := b.StartObject(appendToEntryRequestSize)
	ob.SetText(appendToEntryRequestDirectoryOff, in.Directory)
	ob.SetText(appendToEntryRequestEntryNameOff, in.EntryName)
	ob.SetList(appendToEntryRequestChunksOff, chunksOff, chunksLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- AppendToEntryResponse ---

const appendToEntryResponseSize = 0

// AppendToEntryResponse is a zero-copy view into a ZAP-encoded AppendToEntryResponse
// message. AppendToEntryResponse carries no fields.
type AppendToEntryResponse struct{ o zap.Object }

// WrapAppendToEntryResponse parses b and returns a typed view.
func WrapAppendToEntryResponse(b []byte) (AppendToEntryResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return AppendToEntryResponse{}, err
	}
	return AppendToEntryResponse{o: m.Root()}, nil
}

// AppendToEntryResponseInput collects the field values for NewAppendToEntryResponse.
// AppendToEntryResponse has no fields.
type AppendToEntryResponseInput struct{}

// NewAppendToEntryResponse builds a ZAP-encoded AppendToEntryResponse message.
func NewAppendToEntryResponse(in AppendToEntryResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(appendToEntryResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- DeleteEntryRequest ---

const (
	deleteEntryRequestDirectoryOff            = 0  // directory (field 1, string)
	deleteEntryRequestNameOff                 = 8  // name (field 2, string)
	deleteEntryRequestIsDeleteDataOff         = 16 // is_delete_data (field 4, bool)
	deleteEntryRequestIsRecursiveOff          = 17 // is_recursive (field 5, bool)
	deleteEntryRequestIgnoreRecursiveErrorOff = 18 // ignore_recursive_error (field 6, bool)
	deleteEntryRequestIsFromOtherClusterOff   = 19 // is_from_other_cluster (field 7, bool)
	deleteEntryRequestSignaturesOff           = 24 // signatures (field 8, repeated int32)
	deleteEntryRequestIfNotModifiedAfterOff   = 32 // if_not_modified_after (field 9, int64)
	deleteEntryRequestSize                    = 40
)

// DeleteEntryRequest is a zero-copy view into a ZAP-encoded DeleteEntryRequest message.
type DeleteEntryRequest struct{ o zap.Object }

// WrapDeleteEntryRequest parses b and returns a typed view.
func WrapDeleteEntryRequest(b []byte) (DeleteEntryRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return DeleteEntryRequest{}, err
	}
	return DeleteEntryRequest{o: m.Root()}, nil
}

func (t DeleteEntryRequest) Directory() string { return t.o.Text(deleteEntryRequestDirectoryOff) }
func (t DeleteEntryRequest) Name() string      { return t.o.Text(deleteEntryRequestNameOff) }
func (t DeleteEntryRequest) IsDeleteData() bool {
	return t.o.Bool(deleteEntryRequestIsDeleteDataOff)
}
func (t DeleteEntryRequest) IsRecursive() bool {
	return t.o.Bool(deleteEntryRequestIsRecursiveOff)
}
func (t DeleteEntryRequest) IgnoreRecursiveError() bool {
	return t.o.Bool(deleteEntryRequestIgnoreRecursiveErrorOff)
}
func (t DeleteEntryRequest) IsFromOtherCluster() bool {
	return t.o.Bool(deleteEntryRequestIsFromOtherClusterOff)
}
func (t DeleteEntryRequest) IfNotModifiedAfter() int64 {
	return t.o.Int64(deleteEntryRequestIfNotModifiedAfterOff)
}

// SignaturesLen returns the number of signatures elements (proto field 8,
// repeated int32).
func (t DeleteEntryRequest) SignaturesLen() int {
	return t.o.List(deleteEntryRequestSignaturesOff).Len()
}

// Signature returns the i-th signatures element.
func (t DeleteEntryRequest) Signature(i int) int32 {
	return int32(t.o.List(deleteEntryRequestSignaturesOff).Uint32(i))
}

// DeleteEntryRequestInput collects the field values for NewDeleteEntryRequest.
type DeleteEntryRequestInput struct {
	Directory            string
	Name                 string
	IsDeleteData         bool
	IsRecursive          bool
	IgnoreRecursiveError bool
	IsFromOtherCluster   bool
	Signatures           []int32
	IfNotModifiedAfter   int64
}

// NewDeleteEntryRequest builds a ZAP-encoded DeleteEntryRequest message from in.
func NewDeleteEntryRequest(in DeleteEntryRequestInput) []byte {
	b := zap.NewBuilder(256)
	sigOff, sigLen := addInt32List(b, in.Signatures)
	ob := b.StartObject(deleteEntryRequestSize)
	ob.SetText(deleteEntryRequestDirectoryOff, in.Directory)
	ob.SetText(deleteEntryRequestNameOff, in.Name)
	ob.SetBool(deleteEntryRequestIsDeleteDataOff, in.IsDeleteData)
	ob.SetBool(deleteEntryRequestIsRecursiveOff, in.IsRecursive)
	ob.SetBool(deleteEntryRequestIgnoreRecursiveErrorOff, in.IgnoreRecursiveError)
	ob.SetBool(deleteEntryRequestIsFromOtherClusterOff, in.IsFromOtherCluster)
	ob.SetList(deleteEntryRequestSignaturesOff, sigOff, sigLen)
	ob.SetInt64(deleteEntryRequestIfNotModifiedAfterOff, in.IfNotModifiedAfter)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- DeleteEntryResponse ---

const (
	deleteEntryResponseErrorOff         = 0 // error (field 1, string)
	deleteEntryResponseMetadataEventOff = 8 // metadata_event (field 2, SubscribeMetadataResponse object)
	deleteEntryResponseSize             = 16
)

// DeleteEntryResponse is a zero-copy view into a ZAP-encoded DeleteEntryResponse message.
type DeleteEntryResponse struct{ o zap.Object }

// WrapDeleteEntryResponse parses b and returns a typed view.
func WrapDeleteEntryResponse(b []byte) (DeleteEntryResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return DeleteEntryResponse{}, err
	}
	return DeleteEntryResponse{o: m.Root()}, nil
}

func (t DeleteEntryResponse) Error() string { return t.o.Text(deleteEntryResponseErrorOff) }

// MetadataEvent returns the nested SubscribeMetadataResponse (null-object if unset).
func (t DeleteEntryResponse) MetadataEvent() SubscribeMetadataResponse {
	return SubscribeMetadataResponse{o: t.o.Object(deleteEntryResponseMetadataEventOff)}
}

// HasMetadataEvent reports whether the metadata_event field is present.
func (t DeleteEntryResponse) HasMetadataEvent() bool {
	return !t.o.Object(deleteEntryResponseMetadataEventOff).IsNull()
}

// DeleteEntryResponseInput collects the field values for NewDeleteEntryResponse.
// MetadataEvent, when non-nil, is a standalone SubscribeMetadataResponse buffer.
type DeleteEntryResponseInput struct {
	Error         string
	MetadataEvent []byte
}

// NewDeleteEntryResponse builds a ZAP-encoded DeleteEntryResponse message from in.
func NewDeleteEntryResponse(in DeleteEntryResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(deleteEntryResponseSize)
	ob.SetText(deleteEntryResponseErrorOff, in.Error)
	setNestedObject(b, ob, deleteEntryResponseMetadataEventOff, in.MetadataEvent)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- AtomicRenameEntryRequest ---

const (
	atomicRenameEntryRequestOldDirectoryOff = 0  // old_directory (field 1, string)
	atomicRenameEntryRequestOldNameOff      = 8  // old_name (field 2, string)
	atomicRenameEntryRequestNewDirectoryOff = 16 // new_directory (field 3, string)
	atomicRenameEntryRequestNewNameOff      = 24 // new_name (field 4, string)
	atomicRenameEntryRequestSignaturesOff   = 32 // signatures (field 5, repeated int32)
	atomicRenameEntryRequestSize            = 40
)

// AtomicRenameEntryRequest is a zero-copy view into a ZAP-encoded AtomicRenameEntryRequest message.
type AtomicRenameEntryRequest struct{ o zap.Object }

// WrapAtomicRenameEntryRequest parses b and returns a typed view.
func WrapAtomicRenameEntryRequest(b []byte) (AtomicRenameEntryRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return AtomicRenameEntryRequest{}, err
	}
	return AtomicRenameEntryRequest{o: m.Root()}, nil
}

func (t AtomicRenameEntryRequest) OldDirectory() string {
	return t.o.Text(atomicRenameEntryRequestOldDirectoryOff)
}
func (t AtomicRenameEntryRequest) OldName() string {
	return t.o.Text(atomicRenameEntryRequestOldNameOff)
}
func (t AtomicRenameEntryRequest) NewDirectory() string {
	return t.o.Text(atomicRenameEntryRequestNewDirectoryOff)
}
func (t AtomicRenameEntryRequest) NewName() string {
	return t.o.Text(atomicRenameEntryRequestNewNameOff)
}

// SignaturesLen returns the number of signatures elements (proto field 5,
// repeated int32).
func (t AtomicRenameEntryRequest) SignaturesLen() int {
	return t.o.List(atomicRenameEntryRequestSignaturesOff).Len()
}

// Signature returns the i-th signatures element.
func (t AtomicRenameEntryRequest) Signature(i int) int32 {
	return int32(t.o.List(atomicRenameEntryRequestSignaturesOff).Uint32(i))
}

// AtomicRenameEntryRequestInput collects the field values for NewAtomicRenameEntryRequest.
type AtomicRenameEntryRequestInput struct {
	OldDirectory string
	OldName      string
	NewDirectory string
	NewName      string
	Signatures   []int32
}

// NewAtomicRenameEntryRequest builds a ZAP-encoded AtomicRenameEntryRequest message from in.
func NewAtomicRenameEntryRequest(in AtomicRenameEntryRequestInput) []byte {
	b := zap.NewBuilder(256)
	sigOff, sigLen := addInt32List(b, in.Signatures)
	ob := b.StartObject(atomicRenameEntryRequestSize)
	ob.SetText(atomicRenameEntryRequestOldDirectoryOff, in.OldDirectory)
	ob.SetText(atomicRenameEntryRequestOldNameOff, in.OldName)
	ob.SetText(atomicRenameEntryRequestNewDirectoryOff, in.NewDirectory)
	ob.SetText(atomicRenameEntryRequestNewNameOff, in.NewName)
	ob.SetList(atomicRenameEntryRequestSignaturesOff, sigOff, sigLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- AtomicRenameEntryResponse ---

const atomicRenameEntryResponseSize = 0

// AtomicRenameEntryResponse is a zero-copy view into a ZAP-encoded
// AtomicRenameEntryResponse message. AtomicRenameEntryResponse carries no fields.
type AtomicRenameEntryResponse struct{ o zap.Object }

// WrapAtomicRenameEntryResponse parses b and returns a typed view.
func WrapAtomicRenameEntryResponse(b []byte) (AtomicRenameEntryResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return AtomicRenameEntryResponse{}, err
	}
	return AtomicRenameEntryResponse{o: m.Root()}, nil
}

// AtomicRenameEntryResponseInput collects the field values for
// NewAtomicRenameEntryResponse. AtomicRenameEntryResponse has no fields.
type AtomicRenameEntryResponseInput struct{}

// NewAtomicRenameEntryResponse builds a ZAP-encoded AtomicRenameEntryResponse message.
func NewAtomicRenameEntryResponse(in AtomicRenameEntryResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(atomicRenameEntryResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- StreamRenameEntryRequest ---

const (
	streamRenameEntryRequestOldDirectoryOff = 0  // old_directory (field 1, string)
	streamRenameEntryRequestOldNameOff      = 8  // old_name (field 2, string)
	streamRenameEntryRequestNewDirectoryOff = 16 // new_directory (field 3, string)
	streamRenameEntryRequestNewNameOff      = 24 // new_name (field 4, string)
	streamRenameEntryRequestSignaturesOff   = 32 // signatures (field 5, repeated int32)
	streamRenameEntryRequestSize            = 40
)

// StreamRenameEntryRequest is a zero-copy view into a ZAP-encoded StreamRenameEntryRequest message.
type StreamRenameEntryRequest struct{ o zap.Object }

// WrapStreamRenameEntryRequest parses b and returns a typed view.
func WrapStreamRenameEntryRequest(b []byte) (StreamRenameEntryRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return StreamRenameEntryRequest{}, err
	}
	return StreamRenameEntryRequest{o: m.Root()}, nil
}

func (t StreamRenameEntryRequest) OldDirectory() string {
	return t.o.Text(streamRenameEntryRequestOldDirectoryOff)
}
func (t StreamRenameEntryRequest) OldName() string {
	return t.o.Text(streamRenameEntryRequestOldNameOff)
}
func (t StreamRenameEntryRequest) NewDirectory() string {
	return t.o.Text(streamRenameEntryRequestNewDirectoryOff)
}
func (t StreamRenameEntryRequest) NewName() string {
	return t.o.Text(streamRenameEntryRequestNewNameOff)
}

// SignaturesLen returns the number of signatures elements (proto field 5,
// repeated int32).
func (t StreamRenameEntryRequest) SignaturesLen() int {
	return t.o.List(streamRenameEntryRequestSignaturesOff).Len()
}

// Signature returns the i-th signatures element.
func (t StreamRenameEntryRequest) Signature(i int) int32 {
	return int32(t.o.List(streamRenameEntryRequestSignaturesOff).Uint32(i))
}

// StreamRenameEntryRequestInput collects the field values for NewStreamRenameEntryRequest.
type StreamRenameEntryRequestInput struct {
	OldDirectory string
	OldName      string
	NewDirectory string
	NewName      string
	Signatures   []int32
}

// NewStreamRenameEntryRequest builds a ZAP-encoded StreamRenameEntryRequest message from in.
func NewStreamRenameEntryRequest(in StreamRenameEntryRequestInput) []byte {
	b := zap.NewBuilder(256)
	sigOff, sigLen := addInt32List(b, in.Signatures)
	ob := b.StartObject(streamRenameEntryRequestSize)
	ob.SetText(streamRenameEntryRequestOldDirectoryOff, in.OldDirectory)
	ob.SetText(streamRenameEntryRequestOldNameOff, in.OldName)
	ob.SetText(streamRenameEntryRequestNewDirectoryOff, in.NewDirectory)
	ob.SetText(streamRenameEntryRequestNewNameOff, in.NewName)
	ob.SetList(streamRenameEntryRequestSignaturesOff, sigOff, sigLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- StreamRenameEntryResponse ---

const (
	streamRenameEntryResponseDirectoryOff         = 0  // directory (field 1, string)
	streamRenameEntryResponseEventNotificationOff = 8  // event_notification (field 2, EventNotification object)
	streamRenameEntryResponseTsNsOff              = 16 // ts_ns (field 3, int64)
	streamRenameEntryResponseSize                 = 24
)

// StreamRenameEntryResponse is a zero-copy view into a ZAP-encoded StreamRenameEntryResponse message.
type StreamRenameEntryResponse struct{ o zap.Object }

// WrapStreamRenameEntryResponse parses b and returns a typed view.
func WrapStreamRenameEntryResponse(b []byte) (StreamRenameEntryResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return StreamRenameEntryResponse{}, err
	}
	return StreamRenameEntryResponse{o: m.Root()}, nil
}

func (t StreamRenameEntryResponse) Directory() string {
	return t.o.Text(streamRenameEntryResponseDirectoryOff)
}
func (t StreamRenameEntryResponse) TsNs() int64 {
	return t.o.Int64(streamRenameEntryResponseTsNsOff)
}

// EventNotification returns the nested EventNotification (null-object if unset).
func (t StreamRenameEntryResponse) EventNotification() EventNotification {
	return EventNotification{o: t.o.Object(streamRenameEntryResponseEventNotificationOff)}
}

// HasEventNotification reports whether the event_notification field is present.
func (t StreamRenameEntryResponse) HasEventNotification() bool {
	return !t.o.Object(streamRenameEntryResponseEventNotificationOff).IsNull()
}

// StreamRenameEntryResponseInput collects the field values for
// NewStreamRenameEntryResponse. EventNotification, when non-nil, is a standalone
// EventNotification buffer.
type StreamRenameEntryResponseInput struct {
	Directory         string
	EventNotification []byte
	TsNs              int64
}

// NewStreamRenameEntryResponse builds a ZAP-encoded StreamRenameEntryResponse message from in.
func NewStreamRenameEntryResponse(in StreamRenameEntryResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(streamRenameEntryResponseSize)
	ob.SetText(streamRenameEntryResponseDirectoryOff, in.Directory)
	setNestedObject(b, ob, streamRenameEntryResponseEventNotificationOff, in.EventNotification)
	ob.SetInt64(streamRenameEntryResponseTsNsOff, in.TsNs)
	ob.FinishAsRoot()
	return b.Finish()
}
