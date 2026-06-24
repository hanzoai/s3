// Code generated from filer.proto; DO NOT EDIT.

package filerwire

import (
	zap "github.com/zap-proto/go"
)

// --- SubscribeMetadataRequest ---

const (
	subscribeMetadataRequestClientNameOff                   = 0  // client_name (field 1, string)
	subscribeMetadataRequestPathPrefixOff                   = 8  // path_prefix (field 2, string)
	subscribeMetadataRequestSinceNsOff                      = 16 // since_ns (field 3, int64)
	subscribeMetadataRequestSignatureOff                    = 24 // signature (field 4, int32)
	subscribeMetadataRequestPathPrefixesOff                 = 28 // path_prefixes (field 6, repeated string)
	subscribeMetadataRequestClientIDOff                     = 36 // client_id (field 7, int32)
	subscribeMetadataRequestUntilNsOff                      = 40 // until_ns (field 8, int64)
	subscribeMetadataRequestClientEpochOff                  = 48 // client_epoch (field 9, int32)
	subscribeMetadataRequestDirectoriesOff                  = 56 // directories (field 10, repeated string)
	subscribeMetadataRequestClientSupportsBatchingOff       = 64 // client_supports_batching (field 11, bool)
	subscribeMetadataRequestClientSupportsMetadataChunksOff = 65 // client_supports_metadata_chunks (field 12, bool)
	subscribeMetadataRequestClientSupportsIdleHeartbeatOff  = 66 // client_supports_idle_heartbeat (field 13, bool)
	subscribeMetadataRequestSize                            = 72
)

// SubscribeMetadataRequest is a zero-copy view into a ZAP-encoded SubscribeMetadataRequest message.
type SubscribeMetadataRequest struct{ o zap.Object }

// WrapSubscribeMetadataRequest parses b and returns a typed view.
func WrapSubscribeMetadataRequest(b []byte) (SubscribeMetadataRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return SubscribeMetadataRequest{}, err
	}
	return SubscribeMetadataRequest{o: m.Root()}, nil
}

func (t SubscribeMetadataRequest) ClientName() string {
	return t.o.Text(subscribeMetadataRequestClientNameOff)
}
func (t SubscribeMetadataRequest) PathPrefix() string {
	return t.o.Text(subscribeMetadataRequestPathPrefixOff)
}
func (t SubscribeMetadataRequest) SinceNs() int64 {
	return t.o.Int64(subscribeMetadataRequestSinceNsOff)
}
func (t SubscribeMetadataRequest) Signature() int32 {
	return t.o.Int32(subscribeMetadataRequestSignatureOff)
}
func (t SubscribeMetadataRequest) ClientID() int32 {
	return t.o.Int32(subscribeMetadataRequestClientIDOff)
}
func (t SubscribeMetadataRequest) UntilNs() int64 {
	return t.o.Int64(subscribeMetadataRequestUntilNsOff)
}
func (t SubscribeMetadataRequest) ClientEpoch() int32 {
	return t.o.Int32(subscribeMetadataRequestClientEpochOff)
}
func (t SubscribeMetadataRequest) ClientSupportsBatching() bool {
	return t.o.Bool(subscribeMetadataRequestClientSupportsBatchingOff)
}
func (t SubscribeMetadataRequest) ClientSupportsMetadataChunks() bool {
	return t.o.Bool(subscribeMetadataRequestClientSupportsMetadataChunksOff)
}
func (t SubscribeMetadataRequest) ClientSupportsIdleHeartbeat() bool {
	return t.o.Bool(subscribeMetadataRequestClientSupportsIdleHeartbeatOff)
}

// PathPrefixesLen returns the number of path_prefixes elements (proto field 6,
// repeated string).
func (t SubscribeMetadataRequest) PathPrefixesLen() int {
	return t.o.List(subscribeMetadataRequestPathPrefixesOff).Len()
}

// PathPrefixe returns the i-th path_prefixes element.
func (t SubscribeMetadataRequest) PathPrefixe(i int) string {
	return elemText(t.o.List(subscribeMetadataRequestPathPrefixesOff), i)
}

// DirectoriesLen returns the number of directories elements (proto field 10,
// repeated string).
func (t SubscribeMetadataRequest) DirectoriesLen() int {
	return t.o.List(subscribeMetadataRequestDirectoriesOff).Len()
}

// Directory returns the i-th directories element.
func (t SubscribeMetadataRequest) Directory(i int) string {
	return elemText(t.o.List(subscribeMetadataRequestDirectoriesOff), i)
}

// SubscribeMetadataRequestInput collects the field values for NewSubscribeMetadataRequest.
type SubscribeMetadataRequestInput struct {
	ClientName                   string
	PathPrefix                   string
	SinceNs                      int64
	Signature                    int32
	PathPrefixes                 []string
	ClientID                     int32
	UntilNs                      int64
	ClientEpoch                  int32
	Directories                  []string
	ClientSupportsBatching       bool
	ClientSupportsMetadataChunks bool
	ClientSupportsIdleHeartbeat  bool
}

// NewSubscribeMetadataRequest builds a ZAP-encoded SubscribeMetadataRequest message from in.
func NewSubscribeMetadataRequest(in SubscribeMetadataRequestInput) []byte {
	b := zap.NewBuilder(512)
	prefixesOff, prefixesLen := addStringList(b, in.PathPrefixes)
	dirsOff, dirsLen := addStringList(b, in.Directories)
	ob := b.StartObject(subscribeMetadataRequestSize)
	ob.SetText(subscribeMetadataRequestClientNameOff, in.ClientName)
	ob.SetText(subscribeMetadataRequestPathPrefixOff, in.PathPrefix)
	ob.SetInt64(subscribeMetadataRequestSinceNsOff, in.SinceNs)
	ob.SetInt32(subscribeMetadataRequestSignatureOff, in.Signature)
	ob.SetList(subscribeMetadataRequestPathPrefixesOff, prefixesOff, prefixesLen)
	ob.SetInt32(subscribeMetadataRequestClientIDOff, in.ClientID)
	ob.SetInt64(subscribeMetadataRequestUntilNsOff, in.UntilNs)
	ob.SetInt32(subscribeMetadataRequestClientEpochOff, in.ClientEpoch)
	ob.SetList(subscribeMetadataRequestDirectoriesOff, dirsOff, dirsLen)
	ob.SetBool(subscribeMetadataRequestClientSupportsBatchingOff, in.ClientSupportsBatching)
	ob.SetBool(subscribeMetadataRequestClientSupportsMetadataChunksOff, in.ClientSupportsMetadataChunks)
	ob.SetBool(subscribeMetadataRequestClientSupportsIdleHeartbeatOff, in.ClientSupportsIdleHeartbeat)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- LogFileChunkRef ---

const (
	logFileChunkRefChunksOff   = 0  // chunks (field 1, repeated FileChunk)
	logFileChunkRefFileTsNsOff = 8  // file_ts_ns (field 2, int64)
	logFileChunkRefFilerIDOff  = 16 // filer_id (field 3, string)
	logFileChunkRefSize        = 24
)

// LogFileChunkRef is a zero-copy view into a ZAP-encoded LogFileChunkRef message.
type LogFileChunkRef struct{ o zap.Object }

// WrapLogFileChunkRef parses b and returns a typed view.
func WrapLogFileChunkRef(b []byte) (LogFileChunkRef, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return LogFileChunkRef{}, err
	}
	return LogFileChunkRef{o: m.Root()}, nil
}

func (t LogFileChunkRef) FileTsNs() int64 { return t.o.Int64(logFileChunkRefFileTsNsOff) }
func (t LogFileChunkRef) FilerID() string { return t.o.Text(logFileChunkRefFilerIDOff) }

// ChunksLen returns the number of chunks elements (proto field 1, repeated FileChunk).
func (t LogFileChunkRef) ChunksLen() int { return t.o.List(logFileChunkRefChunksOff).Len() }

// Chunks returns the i-th chunks element as a typed FileChunk view.
func (t LogFileChunkRef) Chunks(i int) FileChunk {
	return FileChunk{o: t.o.List(logFileChunkRefChunksOff).ObjectAt(i)}
}

// LogFileChunkRefInput collects the field values for NewLogFileChunkRef. Each
// Chunks element is a standalone FileChunk buffer (from NewFileChunk).
type LogFileChunkRefInput struct {
	Chunks   [][]byte
	FileTsNs int64
	FilerID  string
}

// NewLogFileChunkRef builds a ZAP-encoded LogFileChunkRef message from in.
func NewLogFileChunkRef(in LogFileChunkRefInput) []byte {
	b := zap.NewBuilder(256)
	chunksOff, chunksLen := addObjectList(b, in.Chunks)
	ob := b.StartObject(logFileChunkRefSize)
	ob.SetList(logFileChunkRefChunksOff, chunksOff, chunksLen)
	ob.SetInt64(logFileChunkRefFileTsNsOff, in.FileTsNs)
	ob.SetText(logFileChunkRefFilerIDOff, in.FilerID)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- SubscribeMetadataResponse ---

const (
	subscribeMetadataResponseDirectoryOff         = 0  // directory (field 1, string)
	subscribeMetadataResponseEventNotificationOff = 8  // event_notification (field 2, EventNotification object)
	subscribeMetadataResponseTsNsOff              = 16 // ts_ns (field 3, int64)
	subscribeMetadataResponseEventsOff            = 24 // events (field 4, repeated SubscribeMetadataResponse)
	subscribeMetadataResponseLogFileRefsOff       = 32 // log_file_refs (field 5, repeated LogFileChunkRef)
	subscribeMetadataResponseSize                 = 40
)

// SubscribeMetadataResponse is a zero-copy view into a ZAP-encoded
// SubscribeMetadataResponse message. The events field is a recursive list of the
// same message type (backlog catch-up batch).
type SubscribeMetadataResponse struct{ o zap.Object }

// WrapSubscribeMetadataResponse parses b and returns a typed view.
func WrapSubscribeMetadataResponse(b []byte) (SubscribeMetadataResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return SubscribeMetadataResponse{}, err
	}
	return SubscribeMetadataResponse{o: m.Root()}, nil
}

func (t SubscribeMetadataResponse) Directory() string {
	return t.o.Text(subscribeMetadataResponseDirectoryOff)
}
func (t SubscribeMetadataResponse) TsNs() int64 {
	return t.o.Int64(subscribeMetadataResponseTsNsOff)
}

// EventNotification returns the nested EventNotification (null-object if unset).
func (t SubscribeMetadataResponse) EventNotification() EventNotification {
	return EventNotification{o: t.o.Object(subscribeMetadataResponseEventNotificationOff)}
}

// HasEventNotification reports whether the event_notification field is present.
func (t SubscribeMetadataResponse) HasEventNotification() bool {
	return !t.o.Object(subscribeMetadataResponseEventNotificationOff).IsNull()
}

// EventsLen returns the number of events elements (proto field 4, repeated
// SubscribeMetadataResponse).
func (t SubscribeMetadataResponse) EventsLen() int {
	return t.o.List(subscribeMetadataResponseEventsOff).Len()
}

// Events returns the i-th events element as a typed SubscribeMetadataResponse view.
func (t SubscribeMetadataResponse) Events(i int) SubscribeMetadataResponse {
	return SubscribeMetadataResponse{o: t.o.List(subscribeMetadataResponseEventsOff).ObjectAt(i)}
}

// LogFileRefsLen returns the number of log_file_refs elements (proto field 5,
// repeated LogFileChunkRef).
func (t SubscribeMetadataResponse) LogFileRefsLen() int {
	return t.o.List(subscribeMetadataResponseLogFileRefsOff).Len()
}

// LogFileRef returns the i-th log_file_refs element as a typed LogFileChunkRef view.
func (t SubscribeMetadataResponse) LogFileRef(i int) LogFileChunkRef {
	return LogFileChunkRef{o: t.o.List(subscribeMetadataResponseLogFileRefsOff).ObjectAt(i)}
}

// SubscribeMetadataResponseInput collects the field values for
// NewSubscribeMetadataResponse. EventNotification, when non-nil, is a standalone
// EventNotification buffer; each Events element is a standalone
// SubscribeMetadataResponse buffer; each LogFileRefs element is a standalone
// LogFileChunkRef buffer.
type SubscribeMetadataResponseInput struct {
	Directory         string
	EventNotification []byte
	TsNs              int64
	Events            [][]byte
	LogFileRefs       [][]byte
}

// NewSubscribeMetadataResponse builds a ZAP-encoded SubscribeMetadataResponse message from in.
func NewSubscribeMetadataResponse(in SubscribeMetadataResponseInput) []byte {
	b := zap.NewBuilder(512)
	eventsOff, eventsLen := addObjectList(b, in.Events)
	refsOff, refsLen := addObjectList(b, in.LogFileRefs)
	ob := b.StartObject(subscribeMetadataResponseSize)
	ob.SetText(subscribeMetadataResponseDirectoryOff, in.Directory)
	setNestedObject(b, ob, subscribeMetadataResponseEventNotificationOff, in.EventNotification)
	ob.SetInt64(subscribeMetadataResponseTsNsOff, in.TsNs)
	ob.SetList(subscribeMetadataResponseEventsOff, eventsOff, eventsLen)
	ob.SetList(subscribeMetadataResponseLogFileRefsOff, refsOff, refsLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ListMetadataSubscribersRequest ---

const (
	listMetadataSubscribersRequestClientTypesOff = 0 // client_types (field 1, repeated string)
	listMetadataSubscribersRequestSize           = 8
)

// ListMetadataSubscribersRequest is a zero-copy view into a ZAP-encoded
// ListMetadataSubscribersRequest message.
type ListMetadataSubscribersRequest struct{ o zap.Object }

// WrapListMetadataSubscribersRequest parses b and returns a typed view.
func WrapListMetadataSubscribersRequest(b []byte) (ListMetadataSubscribersRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ListMetadataSubscribersRequest{}, err
	}
	return ListMetadataSubscribersRequest{o: m.Root()}, nil
}

// ClientTypesLen returns the number of client_types elements (proto field 1,
// repeated string).
func (t ListMetadataSubscribersRequest) ClientTypesLen() int {
	return t.o.List(listMetadataSubscribersRequestClientTypesOff).Len()
}

// ClientType returns the i-th client_types element.
func (t ListMetadataSubscribersRequest) ClientType(i int) string {
	return elemText(t.o.List(listMetadataSubscribersRequestClientTypesOff), i)
}

// ListMetadataSubscribersRequestInput collects the field values for
// NewListMetadataSubscribersRequest.
type ListMetadataSubscribersRequestInput struct {
	ClientTypes []string
}

// NewListMetadataSubscribersRequest builds a ZAP-encoded
// ListMetadataSubscribersRequest message from in.
func NewListMetadataSubscribersRequest(in ListMetadataSubscribersRequestInput) []byte {
	b := zap.NewBuilder(256)
	typesOff, typesLen := addStringList(b, in.ClientTypes)
	ob := b.StartObject(listMetadataSubscribersRequestSize)
	ob.SetList(listMetadataSubscribersRequestClientTypesOff, typesOff, typesLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- MetadataSubscriber ---

const (
	metadataSubscriberClientNameOff    = 0  // client_name (field 1, string)
	metadataSubscriberClientTypeOff    = 8  // client_type (field 2, string)
	metadataSubscriberAddressOff       = 16 // address (field 3, string)
	metadataSubscriberPathPrefixOff    = 24 // path_prefix (field 4, string)
	metadataSubscriberClientIDOff      = 32 // client_id (field 5, int32)
	metadataSubscriberClientEpochOff   = 36 // client_epoch (field 6, int32)
	metadataSubscriberConnectedAtNsOff = 40 // connected_at_ns (field 7, int64)
	metadataSubscriberFilerAddressOff  = 48 // filer_address (field 8, string)
	metadataSubscriberSize             = 56
)

// MetadataSubscriber is a zero-copy view into a ZAP-encoded MetadataSubscriber message.
type MetadataSubscriber struct{ o zap.Object }

// WrapMetadataSubscriber parses b and returns a typed view.
func WrapMetadataSubscriber(b []byte) (MetadataSubscriber, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return MetadataSubscriber{}, err
	}
	return MetadataSubscriber{o: m.Root()}, nil
}

func (t MetadataSubscriber) ClientName() string { return t.o.Text(metadataSubscriberClientNameOff) }
func (t MetadataSubscriber) ClientType() string { return t.o.Text(metadataSubscriberClientTypeOff) }
func (t MetadataSubscriber) Address() string    { return t.o.Text(metadataSubscriberAddressOff) }
func (t MetadataSubscriber) PathPrefix() string { return t.o.Text(metadataSubscriberPathPrefixOff) }
func (t MetadataSubscriber) ClientID() int32    { return t.o.Int32(metadataSubscriberClientIDOff) }
func (t MetadataSubscriber) ClientEpoch() int32 { return t.o.Int32(metadataSubscriberClientEpochOff) }
func (t MetadataSubscriber) ConnectedAtNs() int64 {
	return t.o.Int64(metadataSubscriberConnectedAtNsOff)
}
func (t MetadataSubscriber) FilerAddress() string {
	return t.o.Text(metadataSubscriberFilerAddressOff)
}

// MetadataSubscriberInput collects the field values for NewMetadataSubscriber.
type MetadataSubscriberInput struct {
	ClientName    string
	ClientType    string
	Address       string
	PathPrefix    string
	ClientID      int32
	ClientEpoch   int32
	ConnectedAtNs int64
	FilerAddress  string
}

// NewMetadataSubscriber builds a ZAP-encoded MetadataSubscriber message from in.
func NewMetadataSubscriber(in MetadataSubscriberInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(metadataSubscriberSize)
	ob.SetText(metadataSubscriberClientNameOff, in.ClientName)
	ob.SetText(metadataSubscriberClientTypeOff, in.ClientType)
	ob.SetText(metadataSubscriberAddressOff, in.Address)
	ob.SetText(metadataSubscriberPathPrefixOff, in.PathPrefix)
	ob.SetInt32(metadataSubscriberClientIDOff, in.ClientID)
	ob.SetInt32(metadataSubscriberClientEpochOff, in.ClientEpoch)
	ob.SetInt64(metadataSubscriberConnectedAtNsOff, in.ConnectedAtNs)
	ob.SetText(metadataSubscriberFilerAddressOff, in.FilerAddress)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ListMetadataSubscribersResponse ---

const (
	listMetadataSubscribersResponseSubscribersOff = 0 // subscribers (field 1, repeated MetadataSubscriber)
	listMetadataSubscribersResponseSize           = 8
)

// ListMetadataSubscribersResponse is a zero-copy view into a ZAP-encoded
// ListMetadataSubscribersResponse message.
type ListMetadataSubscribersResponse struct{ o zap.Object }

// WrapListMetadataSubscribersResponse parses b and returns a typed view.
func WrapListMetadataSubscribersResponse(b []byte) (ListMetadataSubscribersResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ListMetadataSubscribersResponse{}, err
	}
	return ListMetadataSubscribersResponse{o: m.Root()}, nil
}

// SubscribersLen returns the number of subscribers elements (proto field 1,
// repeated MetadataSubscriber).
func (t ListMetadataSubscribersResponse) SubscribersLen() int {
	return t.o.List(listMetadataSubscribersResponseSubscribersOff).Len()
}

// Subscribers returns the i-th subscribers element as a typed MetadataSubscriber view.
func (t ListMetadataSubscribersResponse) Subscribers(i int) MetadataSubscriber {
	return MetadataSubscriber{o: t.o.List(listMetadataSubscribersResponseSubscribersOff).ObjectAt(i)}
}

// ListMetadataSubscribersResponseInput collects the field values for
// NewListMetadataSubscribersResponse. Each Subscribers element is a standalone
// MetadataSubscriber buffer (from NewMetadataSubscriber).
type ListMetadataSubscribersResponseInput struct {
	Subscribers [][]byte
}

// NewListMetadataSubscribersResponse builds a ZAP-encoded
// ListMetadataSubscribersResponse message from in.
func NewListMetadataSubscribersResponse(in ListMetadataSubscribersResponseInput) []byte {
	b := zap.NewBuilder(256)
	subsOff, subsLen := addObjectList(b, in.Subscribers)
	ob := b.StartObject(listMetadataSubscribersResponseSize)
	ob.SetList(listMetadataSubscribersResponseSubscribersOff, subsOff, subsLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- TraverseBfsMetadataRequest ---

const (
	traverseBfsMetadataRequestDirectoryOff        = 0 // directory (field 1, string)
	traverseBfsMetadataRequestExcludedPrefixesOff = 8 // excluded_prefixes (field 2, repeated string)
	traverseBfsMetadataRequestSize                = 16
)

// TraverseBfsMetadataRequest is a zero-copy view into a ZAP-encoded TraverseBfsMetadataRequest message.
type TraverseBfsMetadataRequest struct{ o zap.Object }

// WrapTraverseBfsMetadataRequest parses b and returns a typed view.
func WrapTraverseBfsMetadataRequest(b []byte) (TraverseBfsMetadataRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return TraverseBfsMetadataRequest{}, err
	}
	return TraverseBfsMetadataRequest{o: m.Root()}, nil
}

func (t TraverseBfsMetadataRequest) Directory() string {
	return t.o.Text(traverseBfsMetadataRequestDirectoryOff)
}

// ExcludedPrefixesLen returns the number of excluded_prefixes elements (proto
// field 2, repeated string).
func (t TraverseBfsMetadataRequest) ExcludedPrefixesLen() int {
	return t.o.List(traverseBfsMetadataRequestExcludedPrefixesOff).Len()
}

// ExcludedPrefixe returns the i-th excluded_prefixes element.
func (t TraverseBfsMetadataRequest) ExcludedPrefixe(i int) string {
	return elemText(t.o.List(traverseBfsMetadataRequestExcludedPrefixesOff), i)
}

// TraverseBfsMetadataRequestInput collects the field values for NewTraverseBfsMetadataRequest.
type TraverseBfsMetadataRequestInput struct {
	Directory        string
	ExcludedPrefixes []string
}

// NewTraverseBfsMetadataRequest builds a ZAP-encoded TraverseBfsMetadataRequest message from in.
func NewTraverseBfsMetadataRequest(in TraverseBfsMetadataRequestInput) []byte {
	b := zap.NewBuilder(256)
	exOff, exLen := addStringList(b, in.ExcludedPrefixes)
	ob := b.StartObject(traverseBfsMetadataRequestSize)
	ob.SetText(traverseBfsMetadataRequestDirectoryOff, in.Directory)
	ob.SetList(traverseBfsMetadataRequestExcludedPrefixesOff, exOff, exLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- TraverseBfsMetadataResponse ---

const (
	traverseBfsMetadataResponseDirectoryOff = 0 // directory (field 1, string)
	traverseBfsMetadataResponseEntryOff     = 8 // entry (field 2, Entry object)
	traverseBfsMetadataResponseSize         = 16
)

// TraverseBfsMetadataResponse is a zero-copy view into a ZAP-encoded TraverseBfsMetadataResponse message.
type TraverseBfsMetadataResponse struct{ o zap.Object }

// WrapTraverseBfsMetadataResponse parses b and returns a typed view.
func WrapTraverseBfsMetadataResponse(b []byte) (TraverseBfsMetadataResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return TraverseBfsMetadataResponse{}, err
	}
	return TraverseBfsMetadataResponse{o: m.Root()}, nil
}

func (t TraverseBfsMetadataResponse) Directory() string {
	return t.o.Text(traverseBfsMetadataResponseDirectoryOff)
}

// Entry returns the nested Entry (null-object if unset).
func (t TraverseBfsMetadataResponse) Entry() Entry {
	return Entry{o: t.o.Object(traverseBfsMetadataResponseEntryOff)}
}

// HasEntry reports whether the entry field is present.
func (t TraverseBfsMetadataResponse) HasEntry() bool {
	return !t.o.Object(traverseBfsMetadataResponseEntryOff).IsNull()
}

// TraverseBfsMetadataResponseInput collects the field values for
// NewTraverseBfsMetadataResponse. Entry, when non-nil, is a standalone Entry buffer.
type TraverseBfsMetadataResponseInput struct {
	Directory string
	Entry     []byte
}

// NewTraverseBfsMetadataResponse builds a ZAP-encoded TraverseBfsMetadataResponse message from in.
func NewTraverseBfsMetadataResponse(in TraverseBfsMetadataResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(traverseBfsMetadataResponseSize)
	ob.SetText(traverseBfsMetadataResponseDirectoryOff, in.Directory)
	setNestedObject(b, ob, traverseBfsMetadataResponseEntryOff, in.Entry)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- LogEntry ---

const (
	logEntryTsNsOff             = 0  // ts_ns (field 1, int64)
	logEntryPartitionKeyHashOff = 8  // partition_key_hash (field 2, int32)
	logEntryDataOff             = 16 // data (field 3, bytes)
	logEntryKeyOff              = 24 // key (field 4, bytes)
	logEntryOffsetOff           = 32 // offset (field 5, int64)
	logEntrySize                = 40
)

// LogEntry is a zero-copy view into a ZAP-encoded LogEntry message.
type LogEntry struct{ o zap.Object }

// WrapLogEntry parses b and returns a typed view.
func WrapLogEntry(b []byte) (LogEntry, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return LogEntry{}, err
	}
	return LogEntry{o: m.Root()}, nil
}

func (t LogEntry) TsNs() int64             { return t.o.Int64(logEntryTsNsOff) }
func (t LogEntry) PartitionKeyHash() int32 { return t.o.Int32(logEntryPartitionKeyHashOff) }
func (t LogEntry) Data() []byte            { return t.o.Bytes(logEntryDataOff) }
func (t LogEntry) Key() []byte             { return t.o.Bytes(logEntryKeyOff) }
func (t LogEntry) Offset() int64           { return t.o.Int64(logEntryOffsetOff) }

// LogEntryInput collects the field values for NewLogEntry.
type LogEntryInput struct {
	TsNs             int64
	PartitionKeyHash int32
	Data             []byte
	Key              []byte
	Offset           int64
}

// NewLogEntry builds a ZAP-encoded LogEntry message from in.
func NewLogEntry(in LogEntryInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(logEntrySize)
	ob.SetInt64(logEntryTsNsOff, in.TsNs)
	ob.SetInt32(logEntryPartitionKeyHashOff, in.PartitionKeyHash)
	ob.SetBytes(logEntryDataOff, in.Data)
	ob.SetBytes(logEntryKeyOff, in.Key)
	ob.SetInt64(logEntryOffsetOff, in.Offset)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- KeepConnectedRequest ---

const (
	keepConnectedRequestNameOff      = 0  // name (field 1, string)
	keepConnectedRequestGrpcPortOff  = 8  // grpc_port (field 2, uint32)
	keepConnectedRequestResourcesOff = 16 // resources (field 3, repeated string)
	keepConnectedRequestSize         = 24
)

// KeepConnectedRequest is a zero-copy view into a ZAP-encoded KeepConnectedRequest message.
type KeepConnectedRequest struct{ o zap.Object }

// WrapKeepConnectedRequest parses b and returns a typed view.
func WrapKeepConnectedRequest(b []byte) (KeepConnectedRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return KeepConnectedRequest{}, err
	}
	return KeepConnectedRequest{o: m.Root()}, nil
}

func (t KeepConnectedRequest) Name() string     { return t.o.Text(keepConnectedRequestNameOff) }
func (t KeepConnectedRequest) GrpcPort() uint32 { return t.o.Uint32(keepConnectedRequestGrpcPortOff) }

// ResourcesLen returns the number of resources elements (proto field 3, repeated string).
func (t KeepConnectedRequest) ResourcesLen() int {
	return t.o.List(keepConnectedRequestResourcesOff).Len()
}

// Resource returns the i-th resources element.
func (t KeepConnectedRequest) Resource(i int) string {
	return elemText(t.o.List(keepConnectedRequestResourcesOff), i)
}

// KeepConnectedRequestInput collects the field values for NewKeepConnectedRequest.
type KeepConnectedRequestInput struct {
	Name      string
	GrpcPort  uint32
	Resources []string
}

// NewKeepConnectedRequest builds a ZAP-encoded KeepConnectedRequest message from in.
func NewKeepConnectedRequest(in KeepConnectedRequestInput) []byte {
	b := zap.NewBuilder(256)
	resOff, resLen := addStringList(b, in.Resources)
	ob := b.StartObject(keepConnectedRequestSize)
	ob.SetText(keepConnectedRequestNameOff, in.Name)
	ob.SetUint32(keepConnectedRequestGrpcPortOff, in.GrpcPort)
	ob.SetList(keepConnectedRequestResourcesOff, resOff, resLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- KeepConnectedResponse ---

const keepConnectedResponseSize = 0

// KeepConnectedResponse is a zero-copy view into a ZAP-encoded KeepConnectedResponse
// message. KeepConnectedResponse carries no fields.
type KeepConnectedResponse struct{ o zap.Object }

// WrapKeepConnectedResponse parses b and returns a typed view.
func WrapKeepConnectedResponse(b []byte) (KeepConnectedResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return KeepConnectedResponse{}, err
	}
	return KeepConnectedResponse{o: m.Root()}, nil
}

// KeepConnectedResponseInput collects the field values for NewKeepConnectedResponse.
// KeepConnectedResponse has no fields.
type KeepConnectedResponseInput struct{}

// NewKeepConnectedResponse builds a ZAP-encoded KeepConnectedResponse message.
func NewKeepConnectedResponse(in KeepConnectedResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(keepConnectedResponseSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- LocateBrokerRequest ---

const (
	locateBrokerRequestResourceOff = 0 // resource (field 1, string)
	locateBrokerRequestSize        = 8
)

// LocateBrokerRequest is a zero-copy view into a ZAP-encoded LocateBrokerRequest message.
type LocateBrokerRequest struct{ o zap.Object }

// WrapLocateBrokerRequest parses b and returns a typed view.
func WrapLocateBrokerRequest(b []byte) (LocateBrokerRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return LocateBrokerRequest{}, err
	}
	return LocateBrokerRequest{o: m.Root()}, nil
}

func (t LocateBrokerRequest) Resource() string { return t.o.Text(locateBrokerRequestResourceOff) }

// LocateBrokerRequestInput collects the field values for NewLocateBrokerRequest.
type LocateBrokerRequestInput struct {
	Resource string
}

// NewLocateBrokerRequest builds a ZAP-encoded LocateBrokerRequest message from in.
func NewLocateBrokerRequest(in LocateBrokerRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(locateBrokerRequestSize)
	ob.SetText(locateBrokerRequestResourceOff, in.Resource)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- LocateBrokerResponseResource (LocateBrokerResponse.Resource nested message) ---

const (
	locateBrokerResponseResourceGrpcAddressesOff = 0 // grpc_addresses (field 1, string)
	locateBrokerResponseResourceResourceCountOff = 8 // resource_count (field 2, int32)
	locateBrokerResponseResourceSize             = 16
)

// LocateBrokerResponseResource is a zero-copy view into one encoded
// LocateBrokerResponse.Resource.
type LocateBrokerResponseResource struct{ o zap.Object }

// WrapLocateBrokerResponseResource parses b and returns a typed view.
func WrapLocateBrokerResponseResource(b []byte) (LocateBrokerResponseResource, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return LocateBrokerResponseResource{}, err
	}
	return LocateBrokerResponseResource{o: m.Root()}, nil
}

func (t LocateBrokerResponseResource) GrpcAddresses() string {
	return t.o.Text(locateBrokerResponseResourceGrpcAddressesOff)
}
func (t LocateBrokerResponseResource) ResourceCount() int32 {
	return t.o.Int32(locateBrokerResponseResourceResourceCountOff)
}

// LocateBrokerResponseResourceInput collects the field values for
// NewLocateBrokerResponseResource.
type LocateBrokerResponseResourceInput struct {
	GrpcAddresses string
	ResourceCount int32
}

// NewLocateBrokerResponseResource builds a ZAP-encoded LocateBrokerResponse.Resource from in.
func NewLocateBrokerResponseResource(in LocateBrokerResponseResourceInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(locateBrokerResponseResourceSize)
	ob.SetText(locateBrokerResponseResourceGrpcAddressesOff, in.GrpcAddresses)
	ob.SetInt32(locateBrokerResponseResourceResourceCountOff, in.ResourceCount)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- LocateBrokerResponse ---

const (
	locateBrokerResponseFoundOff     = 0 // found (field 1, bool)
	locateBrokerResponseResourcesOff = 8 // resources (field 2, repeated Resource)
	locateBrokerResponseSize         = 16
)

// LocateBrokerResponse is a zero-copy view into a ZAP-encoded LocateBrokerResponse message.
type LocateBrokerResponse struct{ o zap.Object }

// WrapLocateBrokerResponse parses b and returns a typed view.
func WrapLocateBrokerResponse(b []byte) (LocateBrokerResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return LocateBrokerResponse{}, err
	}
	return LocateBrokerResponse{o: m.Root()}, nil
}

func (t LocateBrokerResponse) Found() bool { return t.o.Bool(locateBrokerResponseFoundOff) }

// ResourcesLen returns the number of resources elements (proto field 2, repeated Resource).
func (t LocateBrokerResponse) ResourcesLen() int {
	return t.o.List(locateBrokerResponseResourcesOff).Len()
}

// Resources returns the i-th resources element as a typed LocateBrokerResponseResource view.
func (t LocateBrokerResponse) Resources(i int) LocateBrokerResponseResource {
	return LocateBrokerResponseResource{o: t.o.List(locateBrokerResponseResourcesOff).ObjectAt(i)}
}

// LocateBrokerResponseInput collects the field values for NewLocateBrokerResponse.
// Each Resources element is a standalone LocateBrokerResponseResource buffer.
type LocateBrokerResponseInput struct {
	Found     bool
	Resources [][]byte
}

// NewLocateBrokerResponse builds a ZAP-encoded LocateBrokerResponse message from in.
func NewLocateBrokerResponse(in LocateBrokerResponseInput) []byte {
	b := zap.NewBuilder(256)
	resOff, resLen := addObjectList(b, in.Resources)
	ob := b.StartObject(locateBrokerResponseSize)
	ob.SetBool(locateBrokerResponseFoundOff, in.Found)
	ob.SetList(locateBrokerResponseResourcesOff, resOff, resLen)
	ob.FinishAsRoot()
	return b.Finish()
}
