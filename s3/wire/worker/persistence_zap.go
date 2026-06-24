// Code generated from worker.proto; DO NOT EDIT.

package workerwire

import (
	zap "github.com/zap-proto/go"
)

// --- VolumeHealthMetrics ---
//
// Fixed payload: total_size(8) used_size(8) garbage_size(8) garbage_ratio(8 double)
// file_count(4) deleted_file_count(4) last_modified(8) replica_count(4)
// is_ec_volume(1) collection(8 ptr).

const (
	volumeHealthMetricsTotalSizeOff        = 0
	volumeHealthMetricsUsedSizeOff         = 8
	volumeHealthMetricsGarbageSizeOff      = 16
	volumeHealthMetricsGarbageRatioOff     = 24
	volumeHealthMetricsFileCountOff        = 32
	volumeHealthMetricsDeletedFileCountOff = 36
	volumeHealthMetricsLastModifiedOff     = 40
	volumeHealthMetricsReplicaCountOff     = 48
	volumeHealthMetricsIsEcVolumeOff       = 52
	volumeHealthMetricsCollectionOff       = 56
	volumeHealthMetricsSize                = 64
)

// VolumeHealthMetrics is a zero-copy view into a ZAP-encoded VolumeHealthMetrics message.
type VolumeHealthMetrics struct{ o zap.Object }

// WrapVolumeHealthMetrics parses b and returns a typed view.
func WrapVolumeHealthMetrics(b []byte) (VolumeHealthMetrics, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return VolumeHealthMetrics{}, err
	}
	return VolumeHealthMetrics{o: m.Root()}, nil
}

func (t VolumeHealthMetrics) TotalSize() uint64 { return t.o.Uint64(volumeHealthMetricsTotalSizeOff) }
func (t VolumeHealthMetrics) UsedSize() uint64  { return t.o.Uint64(volumeHealthMetricsUsedSizeOff) }
func (t VolumeHealthMetrics) GarbageSize() uint64 {
	return t.o.Uint64(volumeHealthMetricsGarbageSizeOff)
}
func (t VolumeHealthMetrics) GarbageRatio() float64 {
	return t.o.Float64(volumeHealthMetricsGarbageRatioOff)
}
func (t VolumeHealthMetrics) FileCount() int32 {
	return t.o.Int32(volumeHealthMetricsFileCountOff)
}
func (t VolumeHealthMetrics) DeletedFileCount() int32 {
	return t.o.Int32(volumeHealthMetricsDeletedFileCountOff)
}
func (t VolumeHealthMetrics) LastModified() int64 {
	return t.o.Int64(volumeHealthMetricsLastModifiedOff)
}
func (t VolumeHealthMetrics) ReplicaCount() int32 {
	return t.o.Int32(volumeHealthMetricsReplicaCountOff)
}
func (t VolumeHealthMetrics) IsEcVolume() bool {
	return t.o.Bool(volumeHealthMetricsIsEcVolumeOff)
}
func (t VolumeHealthMetrics) Collection() string {
	return t.o.Text(volumeHealthMetricsCollectionOff)
}

// VolumeHealthMetricsInput collects the field values for NewVolumeHealthMetrics.
type VolumeHealthMetricsInput struct {
	TotalSize        uint64
	UsedSize         uint64
	GarbageSize      uint64
	GarbageRatio     float64
	FileCount        int32
	DeletedFileCount int32
	LastModified     int64
	ReplicaCount     int32
	IsEcVolume       bool
	Collection       string
}

// NewVolumeHealthMetrics builds a ZAP-encoded VolumeHealthMetrics message from in.
func NewVolumeHealthMetrics(in VolumeHealthMetricsInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(volumeHealthMetricsSize)
	ob.SetUint64(volumeHealthMetricsTotalSizeOff, in.TotalSize)
	ob.SetUint64(volumeHealthMetricsUsedSizeOff, in.UsedSize)
	ob.SetUint64(volumeHealthMetricsGarbageSizeOff, in.GarbageSize)
	ob.SetFloat64(volumeHealthMetricsGarbageRatioOff, in.GarbageRatio)
	ob.SetInt32(volumeHealthMetricsFileCountOff, in.FileCount)
	ob.SetInt32(volumeHealthMetricsDeletedFileCountOff, in.DeletedFileCount)
	ob.SetInt64(volumeHealthMetricsLastModifiedOff, in.LastModified)
	ob.SetInt32(volumeHealthMetricsReplicaCountOff, in.ReplicaCount)
	ob.SetBool(volumeHealthMetricsIsEcVolumeOff, in.IsEcVolume)
	ob.SetText(volumeHealthMetricsCollectionOff, in.Collection)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- TaskAssignmentRecord ---
//
// Fixed payload: worker_id(8 ptr) worker_address(8 ptr) assigned_at(8)
// unassigned_at(8) reason(8 ptr).

const (
	taskAssignmentRecordWorkerIDOff      = 0
	taskAssignmentRecordWorkerAddressOff = 8
	taskAssignmentRecordAssignedAtOff    = 16
	taskAssignmentRecordUnassignedAtOff  = 24
	taskAssignmentRecordReasonOff        = 32
	taskAssignmentRecordSize             = 40
)

// TaskAssignmentRecord is a zero-copy view into a ZAP-encoded TaskAssignmentRecord message.
type TaskAssignmentRecord struct{ o zap.Object }

// WrapTaskAssignmentRecord parses b and returns a typed view.
func WrapTaskAssignmentRecord(b []byte) (TaskAssignmentRecord, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return TaskAssignmentRecord{}, err
	}
	return TaskAssignmentRecord{o: m.Root()}, nil
}

func (t TaskAssignmentRecord) WorkerID() string {
	return t.o.Text(taskAssignmentRecordWorkerIDOff)
}
func (t TaskAssignmentRecord) WorkerAddress() string {
	return t.o.Text(taskAssignmentRecordWorkerAddressOff)
}
func (t TaskAssignmentRecord) AssignedAt() int64 {
	return t.o.Int64(taskAssignmentRecordAssignedAtOff)
}
func (t TaskAssignmentRecord) UnassignedAt() int64 {
	return t.o.Int64(taskAssignmentRecordUnassignedAtOff)
}
func (t TaskAssignmentRecord) Reason() string { return t.o.Text(taskAssignmentRecordReasonOff) }

// TaskAssignmentRecordInput collects the field values for NewTaskAssignmentRecord.
type TaskAssignmentRecordInput struct {
	WorkerID      string
	WorkerAddress string
	AssignedAt    int64
	UnassignedAt  int64
	Reason        string
}

// NewTaskAssignmentRecord builds a ZAP-encoded TaskAssignmentRecord message from in.
func NewTaskAssignmentRecord(in TaskAssignmentRecordInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(taskAssignmentRecordSize)
	ob.SetText(taskAssignmentRecordWorkerIDOff, in.WorkerID)
	ob.SetText(taskAssignmentRecordWorkerAddressOff, in.WorkerAddress)
	ob.SetInt64(taskAssignmentRecordAssignedAtOff, in.AssignedAt)
	ob.SetInt64(taskAssignmentRecordUnassignedAtOff, in.UnassignedAt)
	ob.SetText(taskAssignmentRecordReasonOff, in.Reason)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- TaskCreationMetrics ---
//
// Fixed payload: trigger_metric(8 ptr) metric_value(8 double) threshold(8 double)
// volume_metrics(4 obj ptr) additional_data(8 list ptr). volume_metrics is a
// singular VolumeHealthMetrics nested object; additional_data is map<string,string>.

const (
	taskCreationMetricsTriggerMetricOff  = 0
	taskCreationMetricsMetricValueOff    = 8
	taskCreationMetricsThresholdOff      = 16
	taskCreationMetricsVolumeMetricsOff  = 24
	taskCreationMetricsAdditionalDataOff = 28
	taskCreationMetricsSize              = 36
)

// TaskCreationMetrics is a zero-copy view into a ZAP-encoded TaskCreationMetrics message.
type TaskCreationMetrics struct{ o zap.Object }

// WrapTaskCreationMetrics parses b and returns a typed view.
func WrapTaskCreationMetrics(b []byte) (TaskCreationMetrics, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return TaskCreationMetrics{}, err
	}
	return TaskCreationMetrics{o: m.Root()}, nil
}

func (t TaskCreationMetrics) TriggerMetric() string {
	return t.o.Text(taskCreationMetricsTriggerMetricOff)
}
func (t TaskCreationMetrics) MetricValue() float64 {
	return t.o.Float64(taskCreationMetricsMetricValueOff)
}
func (t TaskCreationMetrics) Threshold() float64 {
	return t.o.Float64(taskCreationMetricsThresholdOff)
}

// VolumeMetrics returns the nested VolumeHealthMetrics (null-object if unset).
func (t TaskCreationMetrics) VolumeMetrics() VolumeHealthMetrics {
	return VolumeHealthMetrics{o: t.o.Object(taskCreationMetricsVolumeMetricsOff)}
}

// HasVolumeMetrics reports whether the volume_metrics field is present.
func (t TaskCreationMetrics) HasVolumeMetrics() bool {
	return !t.o.Object(taskCreationMetricsVolumeMetricsOff).IsNull()
}

// AdditionalDataLen returns the number of additional_data entries.
func (t TaskCreationMetrics) AdditionalDataLen() int {
	return t.o.List(taskCreationMetricsAdditionalDataOff).Len()
}

// AdditionalData returns the i-th additional_data entry (key/value).
func (t TaskCreationMetrics) AdditionalData(i int) MapEntry {
	return mapEntryAt(t.o.List(taskCreationMetricsAdditionalDataOff), i)
}

// TaskCreationMetricsInput collects the field values for NewTaskCreationMetrics.
// VolumeMetrics, when non-nil, is a standalone VolumeHealthMetrics buffer.
type TaskCreationMetricsInput struct {
	TriggerMetric  string
	MetricValue    float64
	Threshold      float64
	VolumeMetrics  []byte
	AdditionalData []MapEntry
}

// NewTaskCreationMetrics builds a ZAP-encoded TaskCreationMetrics message from in.
func NewTaskCreationMetrics(in TaskCreationMetricsInput) []byte {
	b := zap.NewBuilder(512)
	vmOff := 0
	if len(in.VolumeMetrics) > 0 {
		vmOff = embedMessage(b, in.VolumeMetrics)
	}
	addOff, addLen := 0, 0
	if len(in.AdditionalData) > 0 {
		el := startElemList(b)
		for _, e := range in.AdditionalData {
			el.addStringEntry(e.Key, e.Value)
		}
		addOff, addLen = el.finish()
	}
	ob := b.StartObject(taskCreationMetricsSize)
	ob.SetText(taskCreationMetricsTriggerMetricOff, in.TriggerMetric)
	ob.SetFloat64(taskCreationMetricsMetricValueOff, in.MetricValue)
	ob.SetFloat64(taskCreationMetricsThresholdOff, in.Threshold)
	ob.SetObject(taskCreationMetricsVolumeMetricsOff, vmOff)
	ob.SetList(taskCreationMetricsAdditionalDataOff, addOff, addLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- MaintenanceTaskData ---
//
// Fixed payload: id(8 ptr) type(8 ptr) priority(8 ptr) status(8 ptr)
// volume_id(4) server(8 ptr) collection(8 ptr) typed_params(4 obj ptr)
// reason(8 ptr) created_at(8) scheduled_at(8) started_at(8) completed_at(8)
// worker_id(8 ptr) error(8 ptr) progress(8 double) retry_count(4) max_retries(4)
// created_by(8 ptr) creation_context(8 ptr) assignment_history(8 list ptr)
// detailed_reason(8 ptr) tags(8 list ptr) creation_metrics(4 obj ptr).
// typed_params is a singular TaskParams nested object; assignment_history is
// repeated TaskAssignmentRecord; tags is map<string,string>; creation_metrics
// is a singular TaskCreationMetrics nested object.

const (
	maintenanceTaskDataIDOff                = 0
	maintenanceTaskDataTypeOff              = 8
	maintenanceTaskDataPriorityOff          = 16
	maintenanceTaskDataStatusOff            = 24
	maintenanceTaskDataVolumeIDOff          = 32
	maintenanceTaskDataServerOff            = 36
	maintenanceTaskDataCollectionOff        = 44
	maintenanceTaskDataTypedParamsOff       = 52
	maintenanceTaskDataReasonOff            = 56
	maintenanceTaskDataCreatedAtOff         = 64
	maintenanceTaskDataScheduledAtOff       = 72
	maintenanceTaskDataStartedAtOff         = 80
	maintenanceTaskDataCompletedAtOff       = 88
	maintenanceTaskDataWorkerIDOff          = 96
	maintenanceTaskDataErrorOff             = 104
	maintenanceTaskDataProgressOff          = 112
	maintenanceTaskDataRetryCountOff        = 120
	maintenanceTaskDataMaxRetriesOff        = 124
	maintenanceTaskDataCreatedByOff         = 128
	maintenanceTaskDataCreationContextOff   = 136
	maintenanceTaskDataAssignmentHistoryOff = 144
	maintenanceTaskDataDetailedReasonOff    = 152
	maintenanceTaskDataTagsOff              = 160
	maintenanceTaskDataCreationMetricsOff   = 168
	maintenanceTaskDataSize                 = 172
)

// MaintenanceTaskData is a zero-copy view into a ZAP-encoded MaintenanceTaskData message.
type MaintenanceTaskData struct{ o zap.Object }

// WrapMaintenanceTaskData parses b and returns a typed view.
func WrapMaintenanceTaskData(b []byte) (MaintenanceTaskData, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return MaintenanceTaskData{}, err
	}
	return MaintenanceTaskData{o: m.Root()}, nil
}

func (t MaintenanceTaskData) ID() string         { return t.o.Text(maintenanceTaskDataIDOff) }
func (t MaintenanceTaskData) Type() string       { return t.o.Text(maintenanceTaskDataTypeOff) }
func (t MaintenanceTaskData) Priority() string   { return t.o.Text(maintenanceTaskDataPriorityOff) }
func (t MaintenanceTaskData) Status() string     { return t.o.Text(maintenanceTaskDataStatusOff) }
func (t MaintenanceTaskData) VolumeID() uint32   { return t.o.Uint32(maintenanceTaskDataVolumeIDOff) }
func (t MaintenanceTaskData) Server() string     { return t.o.Text(maintenanceTaskDataServerOff) }
func (t MaintenanceTaskData) Collection() string { return t.o.Text(maintenanceTaskDataCollectionOff) }
func (t MaintenanceTaskData) Reason() string     { return t.o.Text(maintenanceTaskDataReasonOff) }
func (t MaintenanceTaskData) CreatedAt() int64   { return t.o.Int64(maintenanceTaskDataCreatedAtOff) }
func (t MaintenanceTaskData) ScheduledAt() int64 {
	return t.o.Int64(maintenanceTaskDataScheduledAtOff)
}
func (t MaintenanceTaskData) StartedAt() int64 { return t.o.Int64(maintenanceTaskDataStartedAtOff) }
func (t MaintenanceTaskData) CompletedAt() int64 {
	return t.o.Int64(maintenanceTaskDataCompletedAtOff)
}
func (t MaintenanceTaskData) WorkerID() string { return t.o.Text(maintenanceTaskDataWorkerIDOff) }
func (t MaintenanceTaskData) Error() string    { return t.o.Text(maintenanceTaskDataErrorOff) }
func (t MaintenanceTaskData) Progress() float64 {
	return t.o.Float64(maintenanceTaskDataProgressOff)
}
func (t MaintenanceTaskData) RetryCount() int32 {
	return t.o.Int32(maintenanceTaskDataRetryCountOff)
}
func (t MaintenanceTaskData) MaxRetries() int32 {
	return t.o.Int32(maintenanceTaskDataMaxRetriesOff)
}
func (t MaintenanceTaskData) CreatedBy() string {
	return t.o.Text(maintenanceTaskDataCreatedByOff)
}
func (t MaintenanceTaskData) CreationContext() string {
	return t.o.Text(maintenanceTaskDataCreationContextOff)
}
func (t MaintenanceTaskData) DetailedReason() string {
	return t.o.Text(maintenanceTaskDataDetailedReasonOff)
}

// TypedParams returns the nested TaskParams (null-object if unset).
func (t MaintenanceTaskData) TypedParams() TaskParams {
	return TaskParams{o: t.o.Object(maintenanceTaskDataTypedParamsOff)}
}

// HasTypedParams reports whether the typed_params field is present.
func (t MaintenanceTaskData) HasTypedParams() bool {
	return !t.o.Object(maintenanceTaskDataTypedParamsOff).IsNull()
}

// AssignmentHistoryLen returns the number of assignment_history elements.
func (t MaintenanceTaskData) AssignmentHistoryLen() int {
	return t.o.List(maintenanceTaskDataAssignmentHistoryOff).Len()
}

// AssignmentHistory returns the i-th assignment_history element.
func (t MaintenanceTaskData) AssignmentHistory(i int) TaskAssignmentRecord {
	return TaskAssignmentRecord{o: t.o.List(maintenanceTaskDataAssignmentHistoryOff).ObjectAt(i)}
}

// TagsLen returns the number of tags entries.
func (t MaintenanceTaskData) TagsLen() int { return t.o.List(maintenanceTaskDataTagsOff).Len() }

// Tag returns the i-th tags entry (key/value).
func (t MaintenanceTaskData) Tag(i int) MapEntry {
	return mapEntryAt(t.o.List(maintenanceTaskDataTagsOff), i)
}

// CreationMetrics returns the nested TaskCreationMetrics (null-object if unset).
func (t MaintenanceTaskData) CreationMetrics() TaskCreationMetrics {
	return TaskCreationMetrics{o: t.o.Object(maintenanceTaskDataCreationMetricsOff)}
}

// HasCreationMetrics reports whether the creation_metrics field is present.
func (t MaintenanceTaskData) HasCreationMetrics() bool {
	return !t.o.Object(maintenanceTaskDataCreationMetricsOff).IsNull()
}

// MaintenanceTaskDataInput collects the field values for NewMaintenanceTaskData.
// TypedParams / CreationMetrics, when non-nil, are standalone buffers;
// AssignmentHistory elements are standalone TaskAssignmentRecord buffers.
type MaintenanceTaskDataInput struct {
	ID                string
	Type              string
	Priority          string
	Status            string
	VolumeID          uint32
	Server            string
	Collection        string
	TypedParams       []byte
	Reason            string
	CreatedAt         int64
	ScheduledAt       int64
	StartedAt         int64
	CompletedAt       int64
	WorkerID          string
	Error             string
	Progress          float64
	RetryCount        int32
	MaxRetries        int32
	CreatedBy         string
	CreationContext   string
	AssignmentHistory [][]byte
	DetailedReason    string
	Tags              []MapEntry
	CreationMetrics   []byte
}

// NewMaintenanceTaskData builds a ZAP-encoded MaintenanceTaskData message from in.
func NewMaintenanceTaskData(in MaintenanceTaskDataInput) []byte {
	b := zap.NewBuilder(1024)
	typedOff := 0
	if len(in.TypedParams) > 0 {
		typedOff = embedMessage(b, in.TypedParams)
	}
	metricsOff := 0
	if len(in.CreationMetrics) > 0 {
		metricsOff = embedMessage(b, in.CreationMetrics)
	}
	histOff, histLen := 0, 0
	if len(in.AssignmentHistory) > 0 {
		el := startElemList(b)
		for _, e := range in.AssignmentHistory {
			el.addElem(e)
		}
		histOff, histLen = el.finish()
	}
	tagsOff, tagsLen := 0, 0
	if len(in.Tags) > 0 {
		el := startElemList(b)
		for _, e := range in.Tags {
			el.addStringEntry(e.Key, e.Value)
		}
		tagsOff, tagsLen = el.finish()
	}
	ob := b.StartObject(maintenanceTaskDataSize)
	ob.SetText(maintenanceTaskDataIDOff, in.ID)
	ob.SetText(maintenanceTaskDataTypeOff, in.Type)
	ob.SetText(maintenanceTaskDataPriorityOff, in.Priority)
	ob.SetText(maintenanceTaskDataStatusOff, in.Status)
	ob.SetUint32(maintenanceTaskDataVolumeIDOff, in.VolumeID)
	ob.SetText(maintenanceTaskDataServerOff, in.Server)
	ob.SetText(maintenanceTaskDataCollectionOff, in.Collection)
	ob.SetObject(maintenanceTaskDataTypedParamsOff, typedOff)
	ob.SetText(maintenanceTaskDataReasonOff, in.Reason)
	ob.SetInt64(maintenanceTaskDataCreatedAtOff, in.CreatedAt)
	ob.SetInt64(maintenanceTaskDataScheduledAtOff, in.ScheduledAt)
	ob.SetInt64(maintenanceTaskDataStartedAtOff, in.StartedAt)
	ob.SetInt64(maintenanceTaskDataCompletedAtOff, in.CompletedAt)
	ob.SetText(maintenanceTaskDataWorkerIDOff, in.WorkerID)
	ob.SetText(maintenanceTaskDataErrorOff, in.Error)
	ob.SetFloat64(maintenanceTaskDataProgressOff, in.Progress)
	ob.SetInt32(maintenanceTaskDataRetryCountOff, in.RetryCount)
	ob.SetInt32(maintenanceTaskDataMaxRetriesOff, in.MaxRetries)
	ob.SetText(maintenanceTaskDataCreatedByOff, in.CreatedBy)
	ob.SetText(maintenanceTaskDataCreationContextOff, in.CreationContext)
	ob.SetList(maintenanceTaskDataAssignmentHistoryOff, histOff, histLen)
	ob.SetText(maintenanceTaskDataDetailedReasonOff, in.DetailedReason)
	ob.SetList(maintenanceTaskDataTagsOff, tagsOff, tagsLen)
	ob.SetObject(maintenanceTaskDataCreationMetricsOff, metricsOff)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- TaskStateFile ---
//
// Fixed payload: task(4 obj ptr) last_updated(8) admin_version(8 ptr).
// task is a singular MaintenanceTaskData nested object.

const (
	taskStateFileTaskOff         = 0
	taskStateFileLastUpdatedOff  = 8
	taskStateFileAdminVersionOff = 16
	taskStateFileSize            = 24
)

// TaskStateFile is a zero-copy view into a ZAP-encoded TaskStateFile message.
type TaskStateFile struct{ o zap.Object }

// WrapTaskStateFile parses b and returns a typed view.
func WrapTaskStateFile(b []byte) (TaskStateFile, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return TaskStateFile{}, err
	}
	return TaskStateFile{o: m.Root()}, nil
}

func (t TaskStateFile) LastUpdated() int64   { return t.o.Int64(taskStateFileLastUpdatedOff) }
func (t TaskStateFile) AdminVersion() string { return t.o.Text(taskStateFileAdminVersionOff) }

// Task returns the nested MaintenanceTaskData (null-object if unset).
func (t TaskStateFile) Task() MaintenanceTaskData {
	return MaintenanceTaskData{o: t.o.Object(taskStateFileTaskOff)}
}

// HasTask reports whether the task field is present.
func (t TaskStateFile) HasTask() bool {
	return !t.o.Object(taskStateFileTaskOff).IsNull()
}

// TaskStateFileInput collects the field values for NewTaskStateFile. Task, when
// non-nil, is a standalone MaintenanceTaskData buffer.
type TaskStateFileInput struct {
	Task         []byte
	LastUpdated  int64
	AdminVersion string
}

// NewTaskStateFile builds a ZAP-encoded TaskStateFile message from in.
func NewTaskStateFile(in TaskStateFileInput) []byte {
	b := zap.NewBuilder(1024)
	taskOff := 0
	if len(in.Task) > 0 {
		taskOff = embedMessage(b, in.Task)
	}
	ob := b.StartObject(taskStateFileSize)
	ob.SetObject(taskStateFileTaskOff, taskOff)
	ob.SetInt64(taskStateFileLastUpdatedOff, in.LastUpdated)
	ob.SetText(taskStateFileAdminVersionOff, in.AdminVersion)
	ob.FinishAsRoot()
	return b.Finish()
}
