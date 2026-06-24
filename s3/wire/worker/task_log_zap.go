// Code generated from worker.proto; DO NOT EDIT.

package workerwire

import (
	zap "github.com/zap-proto/go"
)

// --- TaskLogRequest ---
//
// Fixed payload: task_id(8 ptr) worker_id(8 ptr) include_metadata(1)
// max_entries(4) log_level(8 ptr) start_time(8) end_time(8).

const (
	taskLogRequestTaskIDOff          = 0
	taskLogRequestWorkerIDOff        = 8
	taskLogRequestIncludeMetadataOff = 16
	taskLogRequestMaxEntriesOff      = 20
	taskLogRequestLogLevelOff        = 24
	taskLogRequestStartTimeOff       = 32
	taskLogRequestEndTimeOff         = 40
	taskLogRequestSize               = 48
)

// TaskLogRequest is a zero-copy view into a ZAP-encoded TaskLogRequest message.
type TaskLogRequest struct{ o zap.Object }

// WrapTaskLogRequest parses b and returns a typed view.
func WrapTaskLogRequest(b []byte) (TaskLogRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return TaskLogRequest{}, err
	}
	return TaskLogRequest{o: m.Root()}, nil
}

func (t TaskLogRequest) TaskID() string   { return t.o.Text(taskLogRequestTaskIDOff) }
func (t TaskLogRequest) WorkerID() string { return t.o.Text(taskLogRequestWorkerIDOff) }
func (t TaskLogRequest) IncludeMetadata() bool {
	return t.o.Bool(taskLogRequestIncludeMetadataOff)
}
func (t TaskLogRequest) MaxEntries() int32 { return t.o.Int32(taskLogRequestMaxEntriesOff) }
func (t TaskLogRequest) LogLevel() string  { return t.o.Text(taskLogRequestLogLevelOff) }
func (t TaskLogRequest) StartTime() int64  { return t.o.Int64(taskLogRequestStartTimeOff) }
func (t TaskLogRequest) EndTime() int64    { return t.o.Int64(taskLogRequestEndTimeOff) }

// TaskLogRequestInput collects the field values for NewTaskLogRequest.
type TaskLogRequestInput struct {
	TaskID          string
	WorkerID        string
	IncludeMetadata bool
	MaxEntries      int32
	LogLevel        string
	StartTime       int64
	EndTime         int64
}

// NewTaskLogRequest builds a ZAP-encoded TaskLogRequest message from in.
func NewTaskLogRequest(in TaskLogRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(taskLogRequestSize)
	ob.SetText(taskLogRequestTaskIDOff, in.TaskID)
	ob.SetText(taskLogRequestWorkerIDOff, in.WorkerID)
	ob.SetBool(taskLogRequestIncludeMetadataOff, in.IncludeMetadata)
	ob.SetInt32(taskLogRequestMaxEntriesOff, in.MaxEntries)
	ob.SetText(taskLogRequestLogLevelOff, in.LogLevel)
	ob.SetInt64(taskLogRequestStartTimeOff, in.StartTime)
	ob.SetInt64(taskLogRequestEndTimeOff, in.EndTime)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- TaskLogEntry ---
//
// Fixed payload: timestamp(8) level(8 ptr) message(8 ptr) fields(8 list ptr)
// progress(4 float) status(8 ptr). fields is map<string,string>.

const (
	taskLogEntryTimestampOff = 0
	taskLogEntryLevelOff     = 8
	taskLogEntryMessageOff   = 16
	taskLogEntryFieldsOff    = 24
	taskLogEntryProgressOff  = 32
	taskLogEntryStatusOff    = 36
	taskLogEntrySize         = 44
)

// TaskLogEntry is a zero-copy view into a ZAP-encoded TaskLogEntry message.
type TaskLogEntry struct{ o zap.Object }

// WrapTaskLogEntry parses b and returns a typed view.
func WrapTaskLogEntry(b []byte) (TaskLogEntry, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return TaskLogEntry{}, err
	}
	return TaskLogEntry{o: m.Root()}, nil
}

func (t TaskLogEntry) Timestamp() int64  { return t.o.Int64(taskLogEntryTimestampOff) }
func (t TaskLogEntry) Level() string     { return t.o.Text(taskLogEntryLevelOff) }
func (t TaskLogEntry) Message() string   { return t.o.Text(taskLogEntryMessageOff) }
func (t TaskLogEntry) Progress() float32 { return t.o.Float32(taskLogEntryProgressOff) }
func (t TaskLogEntry) Status() string    { return t.o.Text(taskLogEntryStatusOff) }

// FieldsLen returns the number of fields entries.
func (t TaskLogEntry) FieldsLen() int { return t.o.List(taskLogEntryFieldsOff).Len() }

// Field returns the i-th fields entry (key/value).
func (t TaskLogEntry) Field(i int) MapEntry {
	return mapEntryAt(t.o.List(taskLogEntryFieldsOff), i)
}

// TaskLogEntryInput collects the field values for NewTaskLogEntry.
type TaskLogEntryInput struct {
	Timestamp int64
	Level     string
	Message   string
	Fields    []MapEntry
	Progress  float32
	Status    string
}

// NewTaskLogEntry builds a ZAP-encoded TaskLogEntry message from in.
func NewTaskLogEntry(in TaskLogEntryInput) []byte {
	b := zap.NewBuilder(256)
	fieldsOff, fieldsLen := 0, 0
	if len(in.Fields) > 0 {
		el := startElemList(b)
		for _, e := range in.Fields {
			el.addStringEntry(e.Key, e.Value)
		}
		fieldsOff, fieldsLen = el.finish()
	}
	ob := b.StartObject(taskLogEntrySize)
	ob.SetInt64(taskLogEntryTimestampOff, in.Timestamp)
	ob.SetText(taskLogEntryLevelOff, in.Level)
	ob.SetText(taskLogEntryMessageOff, in.Message)
	ob.SetList(taskLogEntryFieldsOff, fieldsOff, fieldsLen)
	ob.SetFloat32(taskLogEntryProgressOff, in.Progress)
	ob.SetText(taskLogEntryStatusOff, in.Status)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- TaskLogMetadata ---
//
// Fixed payload: task_id(8 ptr) task_type(8 ptr) worker_id(8 ptr) start_time(8)
// end_time(8) duration_ms(8) status(8 ptr) progress(4 float) volume_id(4)
// server(8 ptr) collection(8 ptr) log_file_path(8 ptr) created_at(8)
// custom_data(8 list ptr). custom_data is map<string,string>.

const (
	taskLogMetadataTaskIDOff      = 0
	taskLogMetadataTaskTypeOff    = 8
	taskLogMetadataWorkerIDOff    = 16
	taskLogMetadataStartTimeOff   = 24
	taskLogMetadataEndTimeOff     = 32
	taskLogMetadataDurationMsOff  = 40
	taskLogMetadataStatusOff      = 48
	taskLogMetadataProgressOff    = 56
	taskLogMetadataVolumeIDOff    = 60
	taskLogMetadataServerOff      = 64
	taskLogMetadataCollectionOff  = 72
	taskLogMetadataLogFilePathOff = 80
	taskLogMetadataCreatedAtOff   = 88
	taskLogMetadataCustomDataOff  = 96
	taskLogMetadataSize           = 104
)

// TaskLogMetadata is a zero-copy view into a ZAP-encoded TaskLogMetadata message.
type TaskLogMetadata struct{ o zap.Object }

// WrapTaskLogMetadata parses b and returns a typed view.
func WrapTaskLogMetadata(b []byte) (TaskLogMetadata, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return TaskLogMetadata{}, err
	}
	return TaskLogMetadata{o: m.Root()}, nil
}

func (t TaskLogMetadata) TaskID() string     { return t.o.Text(taskLogMetadataTaskIDOff) }
func (t TaskLogMetadata) TaskType() string   { return t.o.Text(taskLogMetadataTaskTypeOff) }
func (t TaskLogMetadata) WorkerID() string   { return t.o.Text(taskLogMetadataWorkerIDOff) }
func (t TaskLogMetadata) StartTime() int64   { return t.o.Int64(taskLogMetadataStartTimeOff) }
func (t TaskLogMetadata) EndTime() int64     { return t.o.Int64(taskLogMetadataEndTimeOff) }
func (t TaskLogMetadata) DurationMs() int64  { return t.o.Int64(taskLogMetadataDurationMsOff) }
func (t TaskLogMetadata) Status() string     { return t.o.Text(taskLogMetadataStatusOff) }
func (t TaskLogMetadata) Progress() float32  { return t.o.Float32(taskLogMetadataProgressOff) }
func (t TaskLogMetadata) VolumeID() uint32   { return t.o.Uint32(taskLogMetadataVolumeIDOff) }
func (t TaskLogMetadata) Server() string     { return t.o.Text(taskLogMetadataServerOff) }
func (t TaskLogMetadata) Collection() string { return t.o.Text(taskLogMetadataCollectionOff) }
func (t TaskLogMetadata) LogFilePath() string {
	return t.o.Text(taskLogMetadataLogFilePathOff)
}
func (t TaskLogMetadata) CreatedAt() int64 { return t.o.Int64(taskLogMetadataCreatedAtOff) }

// CustomDataLen returns the number of custom_data entries.
func (t TaskLogMetadata) CustomDataLen() int {
	return t.o.List(taskLogMetadataCustomDataOff).Len()
}

// CustomData returns the i-th custom_data entry (key/value).
func (t TaskLogMetadata) CustomData(i int) MapEntry {
	return mapEntryAt(t.o.List(taskLogMetadataCustomDataOff), i)
}

// TaskLogMetadataInput collects the field values for NewTaskLogMetadata.
type TaskLogMetadataInput struct {
	TaskID      string
	TaskType    string
	WorkerID    string
	StartTime   int64
	EndTime     int64
	DurationMs  int64
	Status      string
	Progress    float32
	VolumeID    uint32
	Server      string
	Collection  string
	LogFilePath string
	CreatedAt   int64
	CustomData  []MapEntry
}

// NewTaskLogMetadata builds a ZAP-encoded TaskLogMetadata message from in.
func NewTaskLogMetadata(in TaskLogMetadataInput) []byte {
	b := zap.NewBuilder(512)
	customOff, customLen := 0, 0
	if len(in.CustomData) > 0 {
		el := startElemList(b)
		for _, e := range in.CustomData {
			el.addStringEntry(e.Key, e.Value)
		}
		customOff, customLen = el.finish()
	}
	ob := b.StartObject(taskLogMetadataSize)
	ob.SetText(taskLogMetadataTaskIDOff, in.TaskID)
	ob.SetText(taskLogMetadataTaskTypeOff, in.TaskType)
	ob.SetText(taskLogMetadataWorkerIDOff, in.WorkerID)
	ob.SetInt64(taskLogMetadataStartTimeOff, in.StartTime)
	ob.SetInt64(taskLogMetadataEndTimeOff, in.EndTime)
	ob.SetInt64(taskLogMetadataDurationMsOff, in.DurationMs)
	ob.SetText(taskLogMetadataStatusOff, in.Status)
	ob.SetFloat32(taskLogMetadataProgressOff, in.Progress)
	ob.SetUint32(taskLogMetadataVolumeIDOff, in.VolumeID)
	ob.SetText(taskLogMetadataServerOff, in.Server)
	ob.SetText(taskLogMetadataCollectionOff, in.Collection)
	ob.SetText(taskLogMetadataLogFilePathOff, in.LogFilePath)
	ob.SetInt64(taskLogMetadataCreatedAtOff, in.CreatedAt)
	ob.SetList(taskLogMetadataCustomDataOff, customOff, customLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- TaskLogResponse ---
//
// Fixed payload: task_id(8 ptr) worker_id(8 ptr) success(1) error_message(8 ptr)
// metadata(4 obj ptr) log_entries(8 list ptr). metadata is a singular
// TaskLogMetadata nested object; log_entries is repeated TaskLogEntry.

const (
	taskLogResponseTaskIDOff       = 0
	taskLogResponseWorkerIDOff     = 8
	taskLogResponseSuccessOff      = 16
	taskLogResponseErrorMessageOff = 24
	taskLogResponseMetadataOff     = 32
	taskLogResponseLogEntriesOff   = 40
	taskLogResponseSize            = 48
)

// TaskLogResponse is a zero-copy view into a ZAP-encoded TaskLogResponse message.
type TaskLogResponse struct{ o zap.Object }

// WrapTaskLogResponse parses b and returns a typed view.
func WrapTaskLogResponse(b []byte) (TaskLogResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return TaskLogResponse{}, err
	}
	return TaskLogResponse{o: m.Root()}, nil
}

func (t TaskLogResponse) TaskID() string       { return t.o.Text(taskLogResponseTaskIDOff) }
func (t TaskLogResponse) WorkerID() string     { return t.o.Text(taskLogResponseWorkerIDOff) }
func (t TaskLogResponse) Success() bool        { return t.o.Bool(taskLogResponseSuccessOff) }
func (t TaskLogResponse) ErrorMessage() string { return t.o.Text(taskLogResponseErrorMessageOff) }

// Metadata returns the nested TaskLogMetadata (null-object if unset).
func (t TaskLogResponse) Metadata() TaskLogMetadata {
	return TaskLogMetadata{o: t.o.Object(taskLogResponseMetadataOff)}
}

// HasMetadata reports whether the metadata field is present.
func (t TaskLogResponse) HasMetadata() bool {
	return !t.o.Object(taskLogResponseMetadataOff).IsNull()
}

// LogEntriesLen returns the number of log_entries elements.
func (t TaskLogResponse) LogEntriesLen() int {
	return t.o.List(taskLogResponseLogEntriesOff).Len()
}

// LogEntry returns the i-th log_entries element as a typed TaskLogEntry view.
func (t TaskLogResponse) LogEntry(i int) TaskLogEntry {
	return TaskLogEntry{o: t.o.List(taskLogResponseLogEntriesOff).ObjectAt(i)}
}

// TaskLogResponseInput collects the field values for NewTaskLogResponse. Metadata,
// when non-nil, is a standalone TaskLogMetadata buffer; LogEntries elements are
// standalone TaskLogEntry buffers.
type TaskLogResponseInput struct {
	TaskID       string
	WorkerID     string
	Success      bool
	ErrorMessage string
	Metadata     []byte
	LogEntries   [][]byte
}

// NewTaskLogResponse builds a ZAP-encoded TaskLogResponse message from in.
func NewTaskLogResponse(in TaskLogResponseInput) []byte {
	b := zap.NewBuilder(512)
	metaOff := 0
	if len(in.Metadata) > 0 {
		metaOff = embedMessage(b, in.Metadata)
	}
	entriesOff, entriesLen := 0, 0
	if len(in.LogEntries) > 0 {
		el := startElemList(b)
		for _, e := range in.LogEntries {
			el.addElem(e)
		}
		entriesOff, entriesLen = el.finish()
	}
	ob := b.StartObject(taskLogResponseSize)
	ob.SetText(taskLogResponseTaskIDOff, in.TaskID)
	ob.SetText(taskLogResponseWorkerIDOff, in.WorkerID)
	ob.SetBool(taskLogResponseSuccessOff, in.Success)
	ob.SetText(taskLogResponseErrorMessageOff, in.ErrorMessage)
	ob.SetObject(taskLogResponseMetadataOff, metaOff)
	ob.SetList(taskLogResponseLogEntriesOff, entriesOff, entriesLen)
	ob.FinishAsRoot()
	return b.Finish()
}
