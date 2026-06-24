// Code generated from filer.proto; DO NOT EDIT.

package filerwire

import (
	zap "github.com/zap-proto/go"
)

// --- LookupDirectoryEntryRequest ---

const (
	lookupDirectoryEntryRequestDirectoryOff = 0
	lookupDirectoryEntryRequestNameOff      = 8
	lookupDirectoryEntryRequestSize         = 16
)

// LookupDirectoryEntryRequest is a zero-copy view into a ZAP-encoded
// LookupDirectoryEntryRequest message.
type LookupDirectoryEntryRequest struct{ o zap.Object }

// WrapLookupDirectoryEntryRequest parses b and returns a typed view.
func WrapLookupDirectoryEntryRequest(b []byte) (LookupDirectoryEntryRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return LookupDirectoryEntryRequest{}, err
	}
	return LookupDirectoryEntryRequest{o: m.Root()}, nil
}

func (t LookupDirectoryEntryRequest) Directory() string {
	return t.o.Text(lookupDirectoryEntryRequestDirectoryOff)
}
func (t LookupDirectoryEntryRequest) Name() string {
	return t.o.Text(lookupDirectoryEntryRequestNameOff)
}

// LookupDirectoryEntryRequestInput collects the field values for
// NewLookupDirectoryEntryRequest.
type LookupDirectoryEntryRequestInput struct {
	Directory string
	Name      string
}

// NewLookupDirectoryEntryRequest builds a ZAP-encoded LookupDirectoryEntryRequest
// message from in.
func NewLookupDirectoryEntryRequest(in LookupDirectoryEntryRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(lookupDirectoryEntryRequestSize)
	ob.SetText(lookupDirectoryEntryRequestDirectoryOff, in.Directory)
	ob.SetText(lookupDirectoryEntryRequestNameOff, in.Name)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- LookupDirectoryEntryResponse ---

const (
	lookupDirectoryEntryResponseEntryOff = 0
	lookupDirectoryEntryResponseSize     = 8
)

// LookupDirectoryEntryResponse is a zero-copy view into a ZAP-encoded
// LookupDirectoryEntryResponse message.
type LookupDirectoryEntryResponse struct{ o zap.Object }

// WrapLookupDirectoryEntryResponse parses b and returns a typed view.
func WrapLookupDirectoryEntryResponse(b []byte) (LookupDirectoryEntryResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return LookupDirectoryEntryResponse{}, err
	}
	return LookupDirectoryEntryResponse{o: m.Root()}, nil
}

// Entry returns the nested Entry (null-object if unset).
func (t LookupDirectoryEntryResponse) Entry() Entry {
	return Entry{o: t.o.Object(lookupDirectoryEntryResponseEntryOff)}
}

// HasEntry reports whether the entry field is present.
func (t LookupDirectoryEntryResponse) HasEntry() bool {
	return !t.o.Object(lookupDirectoryEntryResponseEntryOff).IsNull()
}

// LookupDirectoryEntryResponseInput collects the field values for
// NewLookupDirectoryEntryResponse. Entry, when non-nil, is a standalone Entry
// buffer (from NewEntry).
type LookupDirectoryEntryResponseInput struct {
	Entry []byte
}

// NewLookupDirectoryEntryResponse builds a ZAP-encoded
// LookupDirectoryEntryResponse message from in.
func NewLookupDirectoryEntryResponse(in LookupDirectoryEntryResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(lookupDirectoryEntryResponseSize)
	setNestedObject(b, ob, lookupDirectoryEntryResponseEntryOff, in.Entry)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ListEntriesRequest ---

const (
	listEntriesRequestDirectoryOff          = 0  // directory (field 1, string)
	listEntriesRequestPrefixOff             = 8  // prefix (field 2, string)
	listEntriesRequestStartFromFileNameOff  = 16 // startFromFileName (field 3, string)
	listEntriesRequestInclusiveStartFromOff = 24 // inclusiveStartFrom (field 4, bool)
	listEntriesRequestLimitOff              = 28 // limit (field 5, uint32)
	listEntriesRequestSnapshotTsNsOff       = 32 // snapshot_ts_ns (field 6, int64)
	listEntriesRequestSize                  = 40
)

// ListEntriesRequest is a zero-copy view into a ZAP-encoded ListEntriesRequest message.
type ListEntriesRequest struct{ o zap.Object }

// WrapListEntriesRequest parses b and returns a typed view.
func WrapListEntriesRequest(b []byte) (ListEntriesRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ListEntriesRequest{}, err
	}
	return ListEntriesRequest{o: m.Root()}, nil
}

func (t ListEntriesRequest) Directory() string {
	return t.o.Text(listEntriesRequestDirectoryOff)
}
func (t ListEntriesRequest) Prefix() string { return t.o.Text(listEntriesRequestPrefixOff) }
func (t ListEntriesRequest) StartFromFileName() string {
	return t.o.Text(listEntriesRequestStartFromFileNameOff)
}
func (t ListEntriesRequest) InclusiveStartFrom() bool {
	return t.o.Bool(listEntriesRequestInclusiveStartFromOff)
}
func (t ListEntriesRequest) Limit() uint32 { return t.o.Uint32(listEntriesRequestLimitOff) }
func (t ListEntriesRequest) SnapshotTsNs() int64 {
	return t.o.Int64(listEntriesRequestSnapshotTsNsOff)
}

// ListEntriesRequestInput collects the field values for NewListEntriesRequest.
type ListEntriesRequestInput struct {
	Directory          string
	Prefix             string
	StartFromFileName  string
	InclusiveStartFrom bool
	Limit              uint32
	SnapshotTsNs       int64
}

// NewListEntriesRequest builds a ZAP-encoded ListEntriesRequest message from in.
func NewListEntriesRequest(in ListEntriesRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(listEntriesRequestSize)
	ob.SetText(listEntriesRequestDirectoryOff, in.Directory)
	ob.SetText(listEntriesRequestPrefixOff, in.Prefix)
	ob.SetText(listEntriesRequestStartFromFileNameOff, in.StartFromFileName)
	ob.SetBool(listEntriesRequestInclusiveStartFromOff, in.InclusiveStartFrom)
	ob.SetUint32(listEntriesRequestLimitOff, in.Limit)
	ob.SetInt64(listEntriesRequestSnapshotTsNsOff, in.SnapshotTsNs)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ListEntriesResponse ---

const (
	listEntriesResponseEntryOff        = 0 // entry (field 1, Entry object)
	listEntriesResponseSnapshotTsNsOff = 8 // snapshot_ts_ns (field 2, int64)
	listEntriesResponseSize            = 16
)

// ListEntriesResponse is a zero-copy view into a ZAP-encoded ListEntriesResponse message.
type ListEntriesResponse struct{ o zap.Object }

// WrapListEntriesResponse parses b and returns a typed view.
func WrapListEntriesResponse(b []byte) (ListEntriesResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ListEntriesResponse{}, err
	}
	return ListEntriesResponse{o: m.Root()}, nil
}

// Entry returns the nested Entry (null-object if unset).
func (t ListEntriesResponse) Entry() Entry {
	return Entry{o: t.o.Object(listEntriesResponseEntryOff)}
}

// HasEntry reports whether the entry field is present.
func (t ListEntriesResponse) HasEntry() bool {
	return !t.o.Object(listEntriesResponseEntryOff).IsNull()
}

func (t ListEntriesResponse) SnapshotTsNs() int64 {
	return t.o.Int64(listEntriesResponseSnapshotTsNsOff)
}

// ListEntriesResponseInput collects the field values for NewListEntriesResponse.
// Entry, when non-nil, is a standalone Entry buffer (from NewEntry).
type ListEntriesResponseInput struct {
	Entry        []byte
	SnapshotTsNs int64
}

// NewListEntriesResponse builds a ZAP-encoded ListEntriesResponse message from in.
func NewListEntriesResponse(in ListEntriesResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(listEntriesResponseSize)
	setNestedObject(b, ob, listEntriesResponseEntryOff, in.Entry)
	ob.SetInt64(listEntriesResponseSnapshotTsNsOff, in.SnapshotTsNs)
	ob.FinishAsRoot()
	return b.Finish()
}
