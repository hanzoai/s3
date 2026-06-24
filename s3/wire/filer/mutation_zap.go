// Code generated from filer.proto; DO NOT EDIT.

package filerwire

import (
	zap "github.com/zap-proto/go"
)

// --- ObjectMutationSetExtendedEntry (ObjectMutation.set_extended map<string,bytes>) ---

const (
	objectMutationSetExtendedEntryKeyOff   = 0
	objectMutationSetExtendedEntryValueOff = 8
	objectMutationSetExtendedEntrySize     = 16
)

// ObjectMutationSetExtendedEntry is a zero-copy view of one map<string,bytes>
// pair of ObjectMutation.set_extended (proto field 5).
type ObjectMutationSetExtendedEntry struct{ o zap.Object }

// WrapObjectMutationSetExtendedEntry parses b and returns a typed view.
func WrapObjectMutationSetExtendedEntry(b []byte) (ObjectMutationSetExtendedEntry, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ObjectMutationSetExtendedEntry{}, err
	}
	return ObjectMutationSetExtendedEntry{o: m.Root()}, nil
}

func (t ObjectMutationSetExtendedEntry) Key() string {
	return t.o.Text(objectMutationSetExtendedEntryKeyOff)
}
func (t ObjectMutationSetExtendedEntry) Value() []byte {
	return t.o.Bytes(objectMutationSetExtendedEntryValueOff)
}

// ObjectMutationSetExtendedEntryInput collects one map<string,bytes> pair.
type ObjectMutationSetExtendedEntryInput struct {
	Key   string
	Value []byte
}

// NewObjectMutationSetExtendedEntry builds one encoded set_extended map entry.
func NewObjectMutationSetExtendedEntry(in ObjectMutationSetExtendedEntryInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(objectMutationSetExtendedEntrySize)
	ob.SetText(objectMutationSetExtendedEntryKeyOff, in.Key)
	ob.SetBytes(objectMutationSetExtendedEntryValueOff, in.Value)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- RecomputeCopyExtendedEntry (Recompute.copy_extended map<string,string>) ---

const (
	recomputeCopyExtendedEntryKeyOff   = 0
	recomputeCopyExtendedEntryValueOff = 8
	recomputeCopyExtendedEntrySize     = 16
)

// RecomputeCopyExtendedEntry is a zero-copy view of one map<string,string> pair
// of Recompute.copy_extended (proto field 3).
type RecomputeCopyExtendedEntry struct{ o zap.Object }

// WrapRecomputeCopyExtendedEntry parses b and returns a typed view.
func WrapRecomputeCopyExtendedEntry(b []byte) (RecomputeCopyExtendedEntry, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return RecomputeCopyExtendedEntry{}, err
	}
	return RecomputeCopyExtendedEntry{o: m.Root()}, nil
}

func (t RecomputeCopyExtendedEntry) Key() string {
	return t.o.Text(recomputeCopyExtendedEntryKeyOff)
}
func (t RecomputeCopyExtendedEntry) Value() string {
	return t.o.Text(recomputeCopyExtendedEntryValueOff)
}

// RecomputeCopyExtendedEntryInput collects one map<string,string> pair.
type RecomputeCopyExtendedEntryInput struct {
	Key   string
	Value string
}

// NewRecomputeCopyExtendedEntry builds one encoded copy_extended map entry.
func NewRecomputeCopyExtendedEntry(in RecomputeCopyExtendedEntryInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(recomputeCopyExtendedEntrySize)
	ob.SetText(recomputeCopyExtendedEntryKeyOff, in.Key)
	ob.SetText(recomputeCopyExtendedEntryValueOff, in.Value)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- Recompute ---

const (
	recomputeScanDirOff      = 0  // scan_dir (field 1, string)
	recomputeDescendingOff   = 8  // descending (field 2, bool)
	recomputeCopyExtendedOff = 16 // copy_extended (field 3, map<string,string>)
	recomputeNameToKeyOff    = 24 // name_to_key (field 4, string)
	recomputeSizeToKeyOff    = 32 // size_to_key (field 5, string)
	recomputeMtimeToKeyOff   = 40 // mtime_to_key (field 6, string)
	recomputeDemoteKeyOff    = 48 // demote_key (field 7, string)
	recomputeDemoteValueOff  = 56 // demote_value (field 8, bytes)
	recomputeExcludeNameOff  = 64 // exclude_name (field 9, string)
	recomputeSize            = 72
)

// Recompute is a zero-copy view into a ZAP-encoded Recompute message.
type Recompute struct{ o zap.Object }

// WrapRecompute parses b and returns a typed view.
func WrapRecompute(b []byte) (Recompute, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return Recompute{}, err
	}
	return Recompute{o: m.Root()}, nil
}

func (t Recompute) ScanDir() string     { return t.o.Text(recomputeScanDirOff) }
func (t Recompute) Descending() bool    { return t.o.Bool(recomputeDescendingOff) }
func (t Recompute) NameToKey() string   { return t.o.Text(recomputeNameToKeyOff) }
func (t Recompute) SizeToKey() string   { return t.o.Text(recomputeSizeToKeyOff) }
func (t Recompute) MtimeToKey() string  { return t.o.Text(recomputeMtimeToKeyOff) }
func (t Recompute) DemoteKey() string   { return t.o.Text(recomputeDemoteKeyOff) }
func (t Recompute) DemoteValue() []byte { return t.o.Bytes(recomputeDemoteValueOff) }
func (t Recompute) ExcludeName() string { return t.o.Text(recomputeExcludeNameOff) }

// CopyExtendedLen returns the number of copy_extended map entries (proto field 3,
// map<string,string>).
func (t Recompute) CopyExtendedLen() int { return t.o.List(recomputeCopyExtendedOff).Len() }

// CopyExtended returns the i-th copy_extended map entry.
func (t Recompute) CopyExtended(i int) RecomputeCopyExtendedEntry {
	return RecomputeCopyExtendedEntry{o: t.o.List(recomputeCopyExtendedOff).ObjectAt(i)}
}

// RecomputeInput collects the field values for NewRecompute. Each CopyExtended
// element is a standalone RecomputeCopyExtendedEntry buffer.
type RecomputeInput struct {
	ScanDir      string
	Descending   bool
	CopyExtended [][]byte
	NameToKey    string
	SizeToKey    string
	MtimeToKey   string
	DemoteKey    string
	DemoteValue  []byte
	ExcludeName  string
}

// NewRecompute builds a ZAP-encoded Recompute message from in.
func NewRecompute(in RecomputeInput) []byte {
	b := zap.NewBuilder(256)
	copyOff, copyLen := addObjectList(b, in.CopyExtended)
	ob := b.StartObject(recomputeSize)
	ob.SetText(recomputeScanDirOff, in.ScanDir)
	ob.SetBool(recomputeDescendingOff, in.Descending)
	ob.SetList(recomputeCopyExtendedOff, copyOff, copyLen)
	ob.SetText(recomputeNameToKeyOff, in.NameToKey)
	ob.SetText(recomputeSizeToKeyOff, in.SizeToKey)
	ob.SetText(recomputeMtimeToKeyOff, in.MtimeToKey)
	ob.SetText(recomputeDemoteKeyOff, in.DemoteKey)
	ob.SetBytes(recomputeDemoteValueOff, in.DemoteValue)
	ob.SetText(recomputeExcludeNameOff, in.ExcludeName)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ObjectMutation ---

const (
	objectMutationTypeOff           = 0  // type (field 1, enum ObjectMutation.Type)
	objectMutationDirectoryOff      = 8  // directory (field 2, string)
	objectMutationNameOff           = 16 // name (field 3, string)
	objectMutationEntryOff          = 24 // entry (field 4, Entry object)
	objectMutationSetExtendedOff    = 32 // set_extended (field 5, map<string,bytes>)
	objectMutationDeleteExtendedOff = 40 // delete_extended (field 6, repeated string)
	objectMutationIsDeleteDataOff   = 48 // is_delete_data (field 7, bool)
	objectMutationIsRecursiveOff    = 49 // is_recursive (field 8, bool)
	objectMutationRecomputeOff      = 52 // recompute (field 9, Recompute object)
	objectMutationSetContentOff     = 56 // set_content (field 10, bool)
	objectMutationContentOff        = 64 // content (field 11, bytes)
	objectMutationTouchMtimeOff     = 72 // touch_mtime (field 12, bool)
	objectMutationSize              = 80
)

// ObjectMutation is a zero-copy view into a ZAP-encoded ObjectMutation message.
type ObjectMutation struct{ o zap.Object }

// WrapObjectMutation parses b and returns a typed view.
func WrapObjectMutation(b []byte) (ObjectMutation, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ObjectMutation{}, err
	}
	return ObjectMutation{o: m.Root()}, nil
}

// Type reads the type enum (see ObjectMutationType* constants).
func (t ObjectMutation) Type() uint32       { return t.o.Uint32(objectMutationTypeOff) }
func (t ObjectMutation) Directory() string  { return t.o.Text(objectMutationDirectoryOff) }
func (t ObjectMutation) Name() string       { return t.o.Text(objectMutationNameOff) }
func (t ObjectMutation) IsDeleteData() bool { return t.o.Bool(objectMutationIsDeleteDataOff) }
func (t ObjectMutation) IsRecursive() bool  { return t.o.Bool(objectMutationIsRecursiveOff) }
func (t ObjectMutation) SetContent() bool   { return t.o.Bool(objectMutationSetContentOff) }
func (t ObjectMutation) Content() []byte    { return t.o.Bytes(objectMutationContentOff) }
func (t ObjectMutation) TouchMtime() bool   { return t.o.Bool(objectMutationTouchMtimeOff) }

// Entry returns the nested Entry (null-object if unset).
func (t ObjectMutation) Entry() Entry { return Entry{o: t.o.Object(objectMutationEntryOff)} }

// HasEntry reports whether the entry field is present.
func (t ObjectMutation) HasEntry() bool { return !t.o.Object(objectMutationEntryOff).IsNull() }

// Recompute returns the nested Recompute (null-object if unset).
func (t ObjectMutation) Recompute() Recompute {
	return Recompute{o: t.o.Object(objectMutationRecomputeOff)}
}

// HasRecompute reports whether the recompute field is present.
func (t ObjectMutation) HasRecompute() bool {
	return !t.o.Object(objectMutationRecomputeOff).IsNull()
}

// SetExtendedLen returns the number of set_extended map entries (proto field 5,
// map<string,bytes>).
func (t ObjectMutation) SetExtendedLen() int {
	return t.o.List(objectMutationSetExtendedOff).Len()
}

// SetExtended returns the i-th set_extended map entry.
func (t ObjectMutation) SetExtended(i int) ObjectMutationSetExtendedEntry {
	return ObjectMutationSetExtendedEntry{o: t.o.List(objectMutationSetExtendedOff).ObjectAt(i)}
}

// DeleteExtendedLen returns the number of delete_extended elements (proto field
// 6, repeated string).
func (t ObjectMutation) DeleteExtendedLen() int {
	return t.o.List(objectMutationDeleteExtendedOff).Len()
}

// DeleteExtended returns the i-th delete_extended element.
func (t ObjectMutation) DeleteExtended(i int) string {
	return elemText(t.o.List(objectMutationDeleteExtendedOff), i)
}

// ObjectMutationInput collects the field values for NewObjectMutation. Entry and
// Recompute, when non-nil, are standalone buffers; each SetExtended element is a
// pre-built ObjectMutationSetExtendedEntry buffer.
type ObjectMutationInput struct {
	Type           uint32
	Directory      string
	Name           string
	Entry          []byte
	SetExtended    [][]byte
	DeleteExtended []string
	IsDeleteData   bool
	IsRecursive    bool
	Recompute      []byte
	SetContent     bool
	Content        []byte
	TouchMtime     bool
}

// NewObjectMutation builds a ZAP-encoded ObjectMutation message from in.
func NewObjectMutation(in ObjectMutationInput) []byte {
	b := zap.NewBuilder(512)
	setExtOff, setExtLen := addObjectList(b, in.SetExtended)
	delExtOff, delExtLen := addStringList(b, in.DeleteExtended)
	ob := b.StartObject(objectMutationSize)
	ob.SetUint32(objectMutationTypeOff, in.Type)
	ob.SetText(objectMutationDirectoryOff, in.Directory)
	ob.SetText(objectMutationNameOff, in.Name)
	setNestedObject(b, ob, objectMutationEntryOff, in.Entry)
	ob.SetList(objectMutationSetExtendedOff, setExtOff, setExtLen)
	ob.SetList(objectMutationDeleteExtendedOff, delExtOff, delExtLen)
	ob.SetBool(objectMutationIsDeleteDataOff, in.IsDeleteData)
	ob.SetBool(objectMutationIsRecursiveOff, in.IsRecursive)
	setNestedObject(b, ob, objectMutationRecomputeOff, in.Recompute)
	ob.SetBool(objectMutationSetContentOff, in.SetContent)
	ob.SetBytes(objectMutationContentOff, in.Content)
	ob.SetBool(objectMutationTouchMtimeOff, in.TouchMtime)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ObjectTransactionRequest ---

const (
	objectTransactionRequestLockKeyOff            = 0  // lock_key (field 1, string)
	objectTransactionRequestConditionOff          = 8  // condition (field 2, WriteCondition object)
	objectTransactionRequestMutationsOff          = 16 // mutations (field 3, repeated ObjectMutation)
	objectTransactionRequestIsFromOtherClusterOff = 24 // is_from_other_cluster (field 4, bool)
	objectTransactionRequestSignaturesOff         = 32 // signatures (field 5, repeated int32)
	objectTransactionRequestConditionKeyOff       = 40 // condition_key (field 6, string)
	objectTransactionRequestRouteKeyOff           = 48 // route_key (field 7, string)
	objectTransactionRequestIsMovedOff            = 56 // is_moved (field 8, bool)
	objectTransactionRequestSize                  = 64
)

// ObjectTransactionRequest is a zero-copy view into a ZAP-encoded
// ObjectTransactionRequest message.
type ObjectTransactionRequest struct{ o zap.Object }

// WrapObjectTransactionRequest parses b and returns a typed view.
func WrapObjectTransactionRequest(b []byte) (ObjectTransactionRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ObjectTransactionRequest{}, err
	}
	return ObjectTransactionRequest{o: m.Root()}, nil
}

func (t ObjectTransactionRequest) LockKey() string {
	return t.o.Text(objectTransactionRequestLockKeyOff)
}
func (t ObjectTransactionRequest) IsFromOtherCluster() bool {
	return t.o.Bool(objectTransactionRequestIsFromOtherClusterOff)
}
func (t ObjectTransactionRequest) ConditionKey() string {
	return t.o.Text(objectTransactionRequestConditionKeyOff)
}
func (t ObjectTransactionRequest) RouteKey() string {
	return t.o.Text(objectTransactionRequestRouteKeyOff)
}
func (t ObjectTransactionRequest) IsMoved() bool {
	return t.o.Bool(objectTransactionRequestIsMovedOff)
}

// Condition returns the nested WriteCondition (null-object if unset).
func (t ObjectTransactionRequest) Condition() WriteCondition {
	return WriteCondition{o: t.o.Object(objectTransactionRequestConditionOff)}
}

// HasCondition reports whether the condition field is present.
func (t ObjectTransactionRequest) HasCondition() bool {
	return !t.o.Object(objectTransactionRequestConditionOff).IsNull()
}

// MutationsLen returns the number of mutations elements (proto field 3, repeated
// ObjectMutation).
func (t ObjectTransactionRequest) MutationsLen() int {
	return t.o.List(objectTransactionRequestMutationsOff).Len()
}

// Mutations returns the i-th mutations element as a typed ObjectMutation view.
func (t ObjectTransactionRequest) Mutations(i int) ObjectMutation {
	return ObjectMutation{o: t.o.List(objectTransactionRequestMutationsOff).ObjectAt(i)}
}

// SignaturesLen returns the number of signatures elements (proto field 5,
// repeated int32).
func (t ObjectTransactionRequest) SignaturesLen() int {
	return t.o.List(objectTransactionRequestSignaturesOff).Len()
}

// Signature returns the i-th signatures element.
func (t ObjectTransactionRequest) Signature(i int) int32 {
	return int32(t.o.List(objectTransactionRequestSignaturesOff).Uint32(i))
}

// ObjectTransactionRequestInput collects the field values for
// NewObjectTransactionRequest. Condition, when non-nil, is a standalone
// WriteCondition buffer; each Mutations element is a standalone ObjectMutation
// buffer.
type ObjectTransactionRequestInput struct {
	LockKey            string
	Condition          []byte
	Mutations          [][]byte
	IsFromOtherCluster bool
	Signatures         []int32
	ConditionKey       string
	RouteKey           string
	IsMoved            bool
}

// NewObjectTransactionRequest builds a ZAP-encoded ObjectTransactionRequest
// message from in.
func NewObjectTransactionRequest(in ObjectTransactionRequestInput) []byte {
	b := zap.NewBuilder(512)
	mutsOff, mutsLen := addObjectList(b, in.Mutations)
	sigOff, sigLen := addInt32List(b, in.Signatures)
	ob := b.StartObject(objectTransactionRequestSize)
	ob.SetText(objectTransactionRequestLockKeyOff, in.LockKey)
	setNestedObject(b, ob, objectTransactionRequestConditionOff, in.Condition)
	ob.SetList(objectTransactionRequestMutationsOff, mutsOff, mutsLen)
	ob.SetBool(objectTransactionRequestIsFromOtherClusterOff, in.IsFromOtherCluster)
	ob.SetList(objectTransactionRequestSignaturesOff, sigOff, sigLen)
	ob.SetText(objectTransactionRequestConditionKeyOff, in.ConditionKey)
	ob.SetText(objectTransactionRequestRouteKeyOff, in.RouteKey)
	ob.SetBool(objectTransactionRequestIsMovedOff, in.IsMoved)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ObjectTransactionResponse ---

const (
	objectTransactionResponseErrorOff     = 0 // error (field 1, string)
	objectTransactionResponseErrorCodeOff = 8 // error_code (field 2, enum FilerError)
	objectTransactionResponseSize         = 12
)

// ObjectTransactionResponse is a zero-copy view into a ZAP-encoded
// ObjectTransactionResponse message.
type ObjectTransactionResponse struct{ o zap.Object }

// WrapObjectTransactionResponse parses b and returns a typed view.
func WrapObjectTransactionResponse(b []byte) (ObjectTransactionResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ObjectTransactionResponse{}, err
	}
	return ObjectTransactionResponse{o: m.Root()}, nil
}

func (t ObjectTransactionResponse) Error() string {
	return t.o.Text(objectTransactionResponseErrorOff)
}

// ErrorCode reads the error_code enum (see FilerError* constants).
func (t ObjectTransactionResponse) ErrorCode() uint32 {
	return t.o.Uint32(objectTransactionResponseErrorCodeOff)
}

// ObjectTransactionResponseInput collects the field values for
// NewObjectTransactionResponse.
type ObjectTransactionResponseInput struct {
	Error     string
	ErrorCode uint32
}

// NewObjectTransactionResponse builds a ZAP-encoded ObjectTransactionResponse
// message from in.
func NewObjectTransactionResponse(in ObjectTransactionResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(objectTransactionResponseSize)
	ob.SetText(objectTransactionResponseErrorOff, in.Error)
	ob.SetUint32(objectTransactionResponseErrorCodeOff, in.ErrorCode)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ObjectTransactionBatchRequest ---

const (
	objectTransactionBatchRequestTransactionsOff = 0
	objectTransactionBatchRequestSize            = 8
)

// ObjectTransactionBatchRequest is a zero-copy view into a ZAP-encoded
// ObjectTransactionBatchRequest message.
type ObjectTransactionBatchRequest struct{ o zap.Object }

// WrapObjectTransactionBatchRequest parses b and returns a typed view.
func WrapObjectTransactionBatchRequest(b []byte) (ObjectTransactionBatchRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ObjectTransactionBatchRequest{}, err
	}
	return ObjectTransactionBatchRequest{o: m.Root()}, nil
}

// TransactionsLen returns the number of transactions elements (proto field 1,
// repeated ObjectTransactionRequest).
func (t ObjectTransactionBatchRequest) TransactionsLen() int {
	return t.o.List(objectTransactionBatchRequestTransactionsOff).Len()
}

// Transactions returns the i-th transactions element as a typed
// ObjectTransactionRequest view.
func (t ObjectTransactionBatchRequest) Transactions(i int) ObjectTransactionRequest {
	return ObjectTransactionRequest{o: t.o.List(objectTransactionBatchRequestTransactionsOff).ObjectAt(i)}
}

// ObjectTransactionBatchRequestInput collects the field values for
// NewObjectTransactionBatchRequest. Each Transactions element is a standalone
// ObjectTransactionRequest buffer.
type ObjectTransactionBatchRequestInput struct {
	Transactions [][]byte
}

// NewObjectTransactionBatchRequest builds a ZAP-encoded
// ObjectTransactionBatchRequest message from in.
func NewObjectTransactionBatchRequest(in ObjectTransactionBatchRequestInput) []byte {
	b := zap.NewBuilder(512)
	txOff, txLen := addObjectList(b, in.Transactions)
	ob := b.StartObject(objectTransactionBatchRequestSize)
	ob.SetList(objectTransactionBatchRequestTransactionsOff, txOff, txLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ObjectTransactionBatchResponse ---

const (
	objectTransactionBatchResponseResponsesOff = 0
	objectTransactionBatchResponseSize         = 8
)

// ObjectTransactionBatchResponse is a zero-copy view into a ZAP-encoded
// ObjectTransactionBatchResponse message.
type ObjectTransactionBatchResponse struct{ o zap.Object }

// WrapObjectTransactionBatchResponse parses b and returns a typed view.
func WrapObjectTransactionBatchResponse(b []byte) (ObjectTransactionBatchResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ObjectTransactionBatchResponse{}, err
	}
	return ObjectTransactionBatchResponse{o: m.Root()}, nil
}

// ResponsesLen returns the number of responses elements (proto field 1, repeated
// ObjectTransactionResponse).
func (t ObjectTransactionBatchResponse) ResponsesLen() int {
	return t.o.List(objectTransactionBatchResponseResponsesOff).Len()
}

// Responses returns the i-th responses element as a typed
// ObjectTransactionResponse view.
func (t ObjectTransactionBatchResponse) Responses(i int) ObjectTransactionResponse {
	return ObjectTransactionResponse{o: t.o.List(objectTransactionBatchResponseResponsesOff).ObjectAt(i)}
}

// ObjectTransactionBatchResponseInput collects the field values for
// NewObjectTransactionBatchResponse. Each Responses element is a standalone
// ObjectTransactionResponse buffer.
type ObjectTransactionBatchResponseInput struct {
	Responses [][]byte
}

// NewObjectTransactionBatchResponse builds a ZAP-encoded
// ObjectTransactionBatchResponse message from in.
func NewObjectTransactionBatchResponse(in ObjectTransactionBatchResponseInput) []byte {
	b := zap.NewBuilder(512)
	respOff, respLen := addObjectList(b, in.Responses)
	ob := b.StartObject(objectTransactionBatchResponseSize)
	ob.SetList(objectTransactionBatchResponseResponsesOff, respOff, respLen)
	ob.FinishAsRoot()
	return b.Finish()
}
