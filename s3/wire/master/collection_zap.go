// Code generated from master.proto; DO NOT EDIT.

package masterwire

import (
	zap "github.com/zap-proto/go"
)

// --- Collection ---

const (
	collectionNameOff = 0
	collectionSize    = 8
)

// Collection is a zero-copy view into a ZAP-encoded Collection message.
type Collection struct{ o zap.Object }

// WrapCollection parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapCollection(b []byte) (Collection, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return Collection{}, err
	}
	return Collection{o: m.Root()}, nil
}

// Name reads the name field (proto field 1, string).
func (t Collection) Name() string { return t.o.Text(collectionNameOff) }

// CollectionInput collects the field values for NewCollection.
type CollectionInput struct {
	Name string
}

// NewCollection builds a ZAP-encoded Collection message from in and returns the
// bytes.
func NewCollection(in CollectionInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(collectionSize)
	ob.SetText(collectionNameOff, in.Name)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- CollectionListRequest ---

const (
	collectionListRequestIncludeNormalVolumesOff = 0
	collectionListRequestIncludeEcVolumesOff     = 1
	collectionListRequestSize                    = 2
)

// CollectionListRequest is a zero-copy view into a ZAP-encoded
// CollectionListRequest message.
type CollectionListRequest struct{ o zap.Object }

// WrapCollectionListRequest parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapCollectionListRequest(b []byte) (CollectionListRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return CollectionListRequest{}, err
	}
	return CollectionListRequest{o: m.Root()}, nil
}

// IncludeNormalVolumes reads the include_normal_volumes field (proto field 1,
// bool).
func (t CollectionListRequest) IncludeNormalVolumes() bool {
	return t.o.Bool(collectionListRequestIncludeNormalVolumesOff)
}

// IncludeEcVolumes reads the include_ec_volumes field (proto field 2, bool).
func (t CollectionListRequest) IncludeEcVolumes() bool {
	return t.o.Bool(collectionListRequestIncludeEcVolumesOff)
}

// CollectionListRequestInput collects the field values for
// NewCollectionListRequest.
type CollectionListRequestInput struct {
	IncludeNormalVolumes bool
	IncludeEcVolumes     bool
}

// NewCollectionListRequest builds a ZAP-encoded CollectionListRequest message
// from in and returns the bytes.
func NewCollectionListRequest(in CollectionListRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(collectionListRequestSize)
	ob.SetBool(collectionListRequestIncludeNormalVolumesOff, in.IncludeNormalVolumes)
	ob.SetBool(collectionListRequestIncludeEcVolumesOff, in.IncludeEcVolumes)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- CollectionListResponse ---

const (
	collectionListResponseCollectionsOff = 0
	collectionListResponseSize           = 8
)

// CollectionListResponse is a zero-copy view into a ZAP-encoded
// CollectionListResponse message.
type CollectionListResponse struct{ o zap.Object }

// WrapCollectionListResponse parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapCollectionListResponse(b []byte) (CollectionListResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return CollectionListResponse{}, err
	}
	return CollectionListResponse{o: m.Root()}, nil
}

// CollectionsLen reports the number of collections elements (proto field 1,
// repeated Collection).
func (t CollectionListResponse) CollectionsLen() int {
	return t.o.List(collectionListResponseCollectionsOff).Len()
}

// CollectionAt returns the i-th collections element.
func (t CollectionListResponse) CollectionAt(i int) (Collection, bool) {
	o := t.o.List(collectionListResponseCollectionsOff).ObjectAt(i)
	if o.IsNull() {
		return Collection{}, false
	}
	return Collection{o: o}, true
}

// CollectionListResponseInput collects the field values for
// NewCollectionListResponse. Collections takes each element as its own ZAP
// sub-buffer (from NewCollection).
type CollectionListResponseInput struct {
	Collections [][]byte
}

// NewCollectionListResponse builds a ZAP-encoded CollectionListResponse message
// from in and returns the bytes.
func NewCollectionListResponse(in CollectionListResponseInput) []byte {
	b := zap.NewBuilder(256)
	collectionsOff, collectionsLen := setMsgList(b, in.Collections)
	ob := b.StartObject(collectionListResponseSize)
	ob.SetList(collectionListResponseCollectionsOff, collectionsOff, collectionsLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- CollectionDeleteRequest ---

const (
	collectionDeleteRequestNameOff = 0
	collectionDeleteRequestSize    = 8
)

// CollectionDeleteRequest is a zero-copy view into a ZAP-encoded
// CollectionDeleteRequest message.
type CollectionDeleteRequest struct{ o zap.Object }

// WrapCollectionDeleteRequest parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapCollectionDeleteRequest(b []byte) (CollectionDeleteRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return CollectionDeleteRequest{}, err
	}
	return CollectionDeleteRequest{o: m.Root()}, nil
}

// Name reads the name field (proto field 1, string).
func (t CollectionDeleteRequest) Name() string { return t.o.Text(collectionDeleteRequestNameOff) }

// CollectionDeleteRequestInput collects the field values for
// NewCollectionDeleteRequest.
type CollectionDeleteRequestInput struct {
	Name string
}

// NewCollectionDeleteRequest builds a ZAP-encoded CollectionDeleteRequest
// message from in and returns the bytes.
func NewCollectionDeleteRequest(in CollectionDeleteRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(collectionDeleteRequestSize)
	ob.SetText(collectionDeleteRequestNameOff, in.Name)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- CollectionDeleteResponse ---

const collectionDeleteResponseSize = 0

// CollectionDeleteResponse is a zero-copy view into a ZAP-encoded
// CollectionDeleteResponse message. CollectionDeleteResponse carries no fields.
type CollectionDeleteResponse struct{ o zap.Object }

// WrapCollectionDeleteResponse parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapCollectionDeleteResponse(b []byte) (CollectionDeleteResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return CollectionDeleteResponse{}, err
	}
	return CollectionDeleteResponse{o: m.Root()}, nil
}

// CollectionDeleteResponseInput collects the field values for
// NewCollectionDeleteResponse. CollectionDeleteResponse has no fields.
type CollectionDeleteResponseInput struct{}

// NewCollectionDeleteResponse builds a ZAP-encoded CollectionDeleteResponse
// message and returns the bytes.
func NewCollectionDeleteResponse(in CollectionDeleteResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(collectionDeleteResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}
