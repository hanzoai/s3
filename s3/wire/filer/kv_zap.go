// Code generated from filer.proto; DO NOT EDIT.

package filerwire

import (
	zap "github.com/zap-proto/go"
)

// --- KvGetRequest ---

const (
	kvGetRequestKeyOff = 0 // key (field 1, bytes)
	kvGetRequestSize   = 8
)

// KvGetRequest is a zero-copy view into a ZAP-encoded KvGetRequest message.
type KvGetRequest struct{ o zap.Object }

// WrapKvGetRequest parses b and returns a typed view.
func WrapKvGetRequest(b []byte) (KvGetRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return KvGetRequest{}, err
	}
	return KvGetRequest{o: m.Root()}, nil
}

func (t KvGetRequest) Key() []byte { return t.o.Bytes(kvGetRequestKeyOff) }

// KvGetRequestInput collects the field values for NewKvGetRequest.
type KvGetRequestInput struct {
	Key []byte
}

// NewKvGetRequest builds a ZAP-encoded KvGetRequest message from in.
func NewKvGetRequest(in KvGetRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(kvGetRequestSize)
	ob.SetBytes(kvGetRequestKeyOff, in.Key)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- KvGetResponse ---

const (
	kvGetResponseValueOff = 0 // value (field 1, bytes)
	kvGetResponseErrorOff = 8 // error (field 2, string)
	kvGetResponseSize     = 16
)

// KvGetResponse is a zero-copy view into a ZAP-encoded KvGetResponse message.
type KvGetResponse struct{ o zap.Object }

// WrapKvGetResponse parses b and returns a typed view.
func WrapKvGetResponse(b []byte) (KvGetResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return KvGetResponse{}, err
	}
	return KvGetResponse{o: m.Root()}, nil
}

func (t KvGetResponse) Value() []byte { return t.o.Bytes(kvGetResponseValueOff) }
func (t KvGetResponse) Error() string { return t.o.Text(kvGetResponseErrorOff) }

// KvGetResponseInput collects the field values for NewKvGetResponse.
type KvGetResponseInput struct {
	Value []byte
	Error string
}

// NewKvGetResponse builds a ZAP-encoded KvGetResponse message from in.
func NewKvGetResponse(in KvGetResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(kvGetResponseSize)
	ob.SetBytes(kvGetResponseValueOff, in.Value)
	ob.SetText(kvGetResponseErrorOff, in.Error)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- KvPutRequest ---

const (
	kvPutRequestKeyOff   = 0 // key (field 1, bytes)
	kvPutRequestValueOff = 8 // value (field 2, bytes)
	kvPutRequestSize     = 16
)

// KvPutRequest is a zero-copy view into a ZAP-encoded KvPutRequest message.
type KvPutRequest struct{ o zap.Object }

// WrapKvPutRequest parses b and returns a typed view.
func WrapKvPutRequest(b []byte) (KvPutRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return KvPutRequest{}, err
	}
	return KvPutRequest{o: m.Root()}, nil
}

func (t KvPutRequest) Key() []byte   { return t.o.Bytes(kvPutRequestKeyOff) }
func (t KvPutRequest) Value() []byte { return t.o.Bytes(kvPutRequestValueOff) }

// KvPutRequestInput collects the field values for NewKvPutRequest.
type KvPutRequestInput struct {
	Key   []byte
	Value []byte
}

// NewKvPutRequest builds a ZAP-encoded KvPutRequest message from in.
func NewKvPutRequest(in KvPutRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(kvPutRequestSize)
	ob.SetBytes(kvPutRequestKeyOff, in.Key)
	ob.SetBytes(kvPutRequestValueOff, in.Value)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- KvPutResponse ---

const (
	kvPutResponseErrorOff = 0 // error (field 1, string)
	kvPutResponseSize     = 8
)

// KvPutResponse is a zero-copy view into a ZAP-encoded KvPutResponse message.
type KvPutResponse struct{ o zap.Object }

// WrapKvPutResponse parses b and returns a typed view.
func WrapKvPutResponse(b []byte) (KvPutResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return KvPutResponse{}, err
	}
	return KvPutResponse{o: m.Root()}, nil
}

func (t KvPutResponse) Error() string { return t.o.Text(kvPutResponseErrorOff) }

// KvPutResponseInput collects the field values for NewKvPutResponse.
type KvPutResponseInput struct {
	Error string
}

// NewKvPutResponse builds a ZAP-encoded KvPutResponse message from in.
func NewKvPutResponse(in KvPutResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(kvPutResponseSize)
	ob.SetText(kvPutResponseErrorOff, in.Error)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- CacheRemoteObjectToLocalClusterRequest ---

const (
	cacheRemoteObjectToLocalClusterRequestDirectoryOff           = 0  // directory (field 1, string)
	cacheRemoteObjectToLocalClusterRequestNameOff                = 8  // name (field 2, string)
	cacheRemoteObjectToLocalClusterRequestChunkConcurrencyOff    = 16 // chunk_concurrency (field 3, int32)
	cacheRemoteObjectToLocalClusterRequestDownloadConcurrencyOff = 20 // download_concurrency (field 4, int32)
	cacheRemoteObjectToLocalClusterRequestSize                   = 24
)

// CacheRemoteObjectToLocalClusterRequest is a zero-copy view into a ZAP-encoded
// CacheRemoteObjectToLocalClusterRequest message.
type CacheRemoteObjectToLocalClusterRequest struct{ o zap.Object }

// WrapCacheRemoteObjectToLocalClusterRequest parses b and returns a typed view.
func WrapCacheRemoteObjectToLocalClusterRequest(b []byte) (CacheRemoteObjectToLocalClusterRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return CacheRemoteObjectToLocalClusterRequest{}, err
	}
	return CacheRemoteObjectToLocalClusterRequest{o: m.Root()}, nil
}

func (t CacheRemoteObjectToLocalClusterRequest) Directory() string {
	return t.o.Text(cacheRemoteObjectToLocalClusterRequestDirectoryOff)
}
func (t CacheRemoteObjectToLocalClusterRequest) Name() string {
	return t.o.Text(cacheRemoteObjectToLocalClusterRequestNameOff)
}
func (t CacheRemoteObjectToLocalClusterRequest) ChunkConcurrency() int32 {
	return t.o.Int32(cacheRemoteObjectToLocalClusterRequestChunkConcurrencyOff)
}
func (t CacheRemoteObjectToLocalClusterRequest) DownloadConcurrency() int32 {
	return t.o.Int32(cacheRemoteObjectToLocalClusterRequestDownloadConcurrencyOff)
}

// CacheRemoteObjectToLocalClusterRequestInput collects the field values for
// NewCacheRemoteObjectToLocalClusterRequest.
type CacheRemoteObjectToLocalClusterRequestInput struct {
	Directory           string
	Name                string
	ChunkConcurrency    int32
	DownloadConcurrency int32
}

// NewCacheRemoteObjectToLocalClusterRequest builds a ZAP-encoded
// CacheRemoteObjectToLocalClusterRequest message from in.
func NewCacheRemoteObjectToLocalClusterRequest(in CacheRemoteObjectToLocalClusterRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(cacheRemoteObjectToLocalClusterRequestSize)
	ob.SetText(cacheRemoteObjectToLocalClusterRequestDirectoryOff, in.Directory)
	ob.SetText(cacheRemoteObjectToLocalClusterRequestNameOff, in.Name)
	ob.SetInt32(cacheRemoteObjectToLocalClusterRequestChunkConcurrencyOff, in.ChunkConcurrency)
	ob.SetInt32(cacheRemoteObjectToLocalClusterRequestDownloadConcurrencyOff, in.DownloadConcurrency)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- CacheRemoteObjectToLocalClusterResponse ---

const (
	cacheRemoteObjectToLocalClusterResponseEntryOff         = 0 // entry (field 1, Entry object)
	cacheRemoteObjectToLocalClusterResponseMetadataEventOff = 4 // metadata_event (field 2, SubscribeMetadataResponse object)
	cacheRemoteObjectToLocalClusterResponseSize             = 8
)

// CacheRemoteObjectToLocalClusterResponse is a zero-copy view into a ZAP-encoded
// CacheRemoteObjectToLocalClusterResponse message.
type CacheRemoteObjectToLocalClusterResponse struct{ o zap.Object }

// WrapCacheRemoteObjectToLocalClusterResponse parses b and returns a typed view.
func WrapCacheRemoteObjectToLocalClusterResponse(b []byte) (CacheRemoteObjectToLocalClusterResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return CacheRemoteObjectToLocalClusterResponse{}, err
	}
	return CacheRemoteObjectToLocalClusterResponse{o: m.Root()}, nil
}

// Entry returns the nested Entry (null-object if unset).
func (t CacheRemoteObjectToLocalClusterResponse) Entry() Entry {
	return Entry{o: t.o.Object(cacheRemoteObjectToLocalClusterResponseEntryOff)}
}

// HasEntry reports whether the entry field is present.
func (t CacheRemoteObjectToLocalClusterResponse) HasEntry() bool {
	return !t.o.Object(cacheRemoteObjectToLocalClusterResponseEntryOff).IsNull()
}

// MetadataEvent returns the nested SubscribeMetadataResponse (null-object if unset).
func (t CacheRemoteObjectToLocalClusterResponse) MetadataEvent() SubscribeMetadataResponse {
	return SubscribeMetadataResponse{o: t.o.Object(cacheRemoteObjectToLocalClusterResponseMetadataEventOff)}
}

// HasMetadataEvent reports whether the metadata_event field is present.
func (t CacheRemoteObjectToLocalClusterResponse) HasMetadataEvent() bool {
	return !t.o.Object(cacheRemoteObjectToLocalClusterResponseMetadataEventOff).IsNull()
}

// CacheRemoteObjectToLocalClusterResponseInput collects the field values for
// NewCacheRemoteObjectToLocalClusterResponse. Entry and MetadataEvent, when
// non-nil, are standalone buffers.
type CacheRemoteObjectToLocalClusterResponseInput struct {
	Entry         []byte
	MetadataEvent []byte
}

// NewCacheRemoteObjectToLocalClusterResponse builds a ZAP-encoded
// CacheRemoteObjectToLocalClusterResponse message from in.
func NewCacheRemoteObjectToLocalClusterResponse(in CacheRemoteObjectToLocalClusterResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(cacheRemoteObjectToLocalClusterResponseSize)
	setNestedObject(b, ob, cacheRemoteObjectToLocalClusterResponseEntryOff, in.Entry)
	setNestedObject(b, ob, cacheRemoteObjectToLocalClusterResponseMetadataEventOff, in.MetadataEvent)
	ob.FinishAsRoot()
	return b.Finish()
}
