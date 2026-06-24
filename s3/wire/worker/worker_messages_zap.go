// Code generated from worker.proto; DO NOT EDIT.

package workerwire

import (
	zap "github.com/zap-proto/go"
)

// --- WorkerRegistration ---
//
// Fixed payload: worker_id(8 ptr) address(8 ptr) capabilities(8 list ptr)
// max_concurrent(4) metadata(8 list ptr). capabilities is repeated string;
// metadata is map<string,string>.

const (
	workerRegistrationWorkerIDOff      = 0
	workerRegistrationAddressOff       = 8
	workerRegistrationCapabilitiesOff  = 16
	workerRegistrationMaxConcurrentOff = 24
	workerRegistrationMetadataOff      = 28
	workerRegistrationSize             = 36
)

// WorkerRegistration is a zero-copy view into a ZAP-encoded WorkerRegistration message.
type WorkerRegistration struct{ o zap.Object }

// WrapWorkerRegistration parses b and returns a typed view.
func WrapWorkerRegistration(b []byte) (WorkerRegistration, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return WorkerRegistration{}, err
	}
	return WorkerRegistration{o: m.Root()}, nil
}

func (t WorkerRegistration) WorkerID() string { return t.o.Text(workerRegistrationWorkerIDOff) }
func (t WorkerRegistration) Address() string  { return t.o.Text(workerRegistrationAddressOff) }
func (t WorkerRegistration) MaxConcurrent() int32 {
	return t.o.Int32(workerRegistrationMaxConcurrentOff)
}

// CapabilitiesLen returns the number of capabilities elements.
func (t WorkerRegistration) CapabilitiesLen() int {
	return t.o.List(workerRegistrationCapabilitiesOff).Len()
}

// Capability returns the i-th capabilities element.
func (t WorkerRegistration) Capability(i int) string {
	return elemText(t.o.List(workerRegistrationCapabilitiesOff), i)
}

// MetadataLen returns the number of metadata entries.
func (t WorkerRegistration) MetadataLen() int {
	return t.o.List(workerRegistrationMetadataOff).Len()
}

// Metadata returns the i-th metadata entry (key/value).
func (t WorkerRegistration) Metadata(i int) MapEntry {
	return mapEntryAt(t.o.List(workerRegistrationMetadataOff), i)
}

// WorkerRegistrationInput collects the field values for NewWorkerRegistration.
type WorkerRegistrationInput struct {
	WorkerID      string
	Address       string
	Capabilities  []string
	MaxConcurrent int32
	Metadata      []MapEntry
}

// NewWorkerRegistration builds a ZAP-encoded WorkerRegistration message from in.
func NewWorkerRegistration(in WorkerRegistrationInput) []byte {
	b := zap.NewBuilder(256)
	capsOff, capsLen := 0, 0
	if len(in.Capabilities) > 0 {
		el := startElemList(b)
		for _, s := range in.Capabilities {
			el.addText(s)
		}
		capsOff, capsLen = el.finish()
	}
	metaOff, metaLen := 0, 0
	if len(in.Metadata) > 0 {
		el := startElemList(b)
		for _, e := range in.Metadata {
			el.addStringEntry(e.Key, e.Value)
		}
		metaOff, metaLen = el.finish()
	}
	ob := b.StartObject(workerRegistrationSize)
	ob.SetText(workerRegistrationWorkerIDOff, in.WorkerID)
	ob.SetText(workerRegistrationAddressOff, in.Address)
	ob.SetList(workerRegistrationCapabilitiesOff, capsOff, capsLen)
	ob.SetInt32(workerRegistrationMaxConcurrentOff, in.MaxConcurrent)
	ob.SetList(workerRegistrationMetadataOff, metaOff, metaLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- RegistrationResponse ---
//
// Fixed payload: success(1) message(8 ptr) assigned_worker_id(8 ptr).

const (
	registrationResponseSuccessOff          = 0
	registrationResponseMessageOff          = 8
	registrationResponseAssignedWorkerIDOff = 16
	registrationResponseSize                = 24
)

// RegistrationResponse is a zero-copy view into a ZAP-encoded RegistrationResponse message.
type RegistrationResponse struct{ o zap.Object }

// WrapRegistrationResponse parses b and returns a typed view.
func WrapRegistrationResponse(b []byte) (RegistrationResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return RegistrationResponse{}, err
	}
	return RegistrationResponse{o: m.Root()}, nil
}

func (t RegistrationResponse) Success() bool   { return t.o.Bool(registrationResponseSuccessOff) }
func (t RegistrationResponse) Message() string { return t.o.Text(registrationResponseMessageOff) }
func (t RegistrationResponse) AssignedWorkerID() string {
	return t.o.Text(registrationResponseAssignedWorkerIDOff)
}

// RegistrationResponseInput collects the field values for NewRegistrationResponse.
type RegistrationResponseInput struct {
	Success          bool
	Message          string
	AssignedWorkerID string
}

// NewRegistrationResponse builds a ZAP-encoded RegistrationResponse message from in.
func NewRegistrationResponse(in RegistrationResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(registrationResponseSize)
	ob.SetBool(registrationResponseSuccessOff, in.Success)
	ob.SetText(registrationResponseMessageOff, in.Message)
	ob.SetText(registrationResponseAssignedWorkerIDOff, in.AssignedWorkerID)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- WorkerHeartbeat ---
//
// Fixed payload: worker_id(8 ptr) status(8 ptr) current_load(4)
// max_concurrent(4) current_task_ids(8 list ptr) tasks_completed(4)
// tasks_failed(4) uptime_seconds(8). current_task_ids is repeated string.

const (
	workerHeartbeatWorkerIDOff       = 0
	workerHeartbeatStatusOff         = 8
	workerHeartbeatCurrentLoadOff    = 16
	workerHeartbeatMaxConcurrentOff  = 20
	workerHeartbeatCurrentTaskIDsOff = 24
	workerHeartbeatTasksCompletedOff = 32
	workerHeartbeatTasksFailedOff    = 36
	workerHeartbeatUptimeSecondsOff  = 40
	workerHeartbeatSize              = 48
)

// WorkerHeartbeat is a zero-copy view into a ZAP-encoded WorkerHeartbeat message.
type WorkerHeartbeat struct{ o zap.Object }

// WrapWorkerHeartbeat parses b and returns a typed view.
func WrapWorkerHeartbeat(b []byte) (WorkerHeartbeat, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return WorkerHeartbeat{}, err
	}
	return WorkerHeartbeat{o: m.Root()}, nil
}

func (t WorkerHeartbeat) WorkerID() string   { return t.o.Text(workerHeartbeatWorkerIDOff) }
func (t WorkerHeartbeat) Status() string     { return t.o.Text(workerHeartbeatStatusOff) }
func (t WorkerHeartbeat) CurrentLoad() int32 { return t.o.Int32(workerHeartbeatCurrentLoadOff) }
func (t WorkerHeartbeat) MaxConcurrent() int32 {
	return t.o.Int32(workerHeartbeatMaxConcurrentOff)
}
func (t WorkerHeartbeat) TasksCompleted() int32 {
	return t.o.Int32(workerHeartbeatTasksCompletedOff)
}
func (t WorkerHeartbeat) TasksFailed() int32 { return t.o.Int32(workerHeartbeatTasksFailedOff) }
func (t WorkerHeartbeat) UptimeSeconds() int64 {
	return t.o.Int64(workerHeartbeatUptimeSecondsOff)
}

// CurrentTaskIDsLen returns the number of current_task_ids elements.
func (t WorkerHeartbeat) CurrentTaskIDsLen() int {
	return t.o.List(workerHeartbeatCurrentTaskIDsOff).Len()
}

// CurrentTaskID returns the i-th current_task_ids element.
func (t WorkerHeartbeat) CurrentTaskID(i int) string {
	return elemText(t.o.List(workerHeartbeatCurrentTaskIDsOff), i)
}

// WorkerHeartbeatInput collects the field values for NewWorkerHeartbeat.
type WorkerHeartbeatInput struct {
	WorkerID       string
	Status         string
	CurrentLoad    int32
	MaxConcurrent  int32
	CurrentTaskIDs []string
	TasksCompleted int32
	TasksFailed    int32
	UptimeSeconds  int64
}

// NewWorkerHeartbeat builds a ZAP-encoded WorkerHeartbeat message from in.
func NewWorkerHeartbeat(in WorkerHeartbeatInput) []byte {
	b := zap.NewBuilder(256)
	tasksOff, tasksLen := 0, 0
	if len(in.CurrentTaskIDs) > 0 {
		el := startElemList(b)
		for _, s := range in.CurrentTaskIDs {
			el.addText(s)
		}
		tasksOff, tasksLen = el.finish()
	}
	ob := b.StartObject(workerHeartbeatSize)
	ob.SetText(workerHeartbeatWorkerIDOff, in.WorkerID)
	ob.SetText(workerHeartbeatStatusOff, in.Status)
	ob.SetInt32(workerHeartbeatCurrentLoadOff, in.CurrentLoad)
	ob.SetInt32(workerHeartbeatMaxConcurrentOff, in.MaxConcurrent)
	ob.SetList(workerHeartbeatCurrentTaskIDsOff, tasksOff, tasksLen)
	ob.SetInt32(workerHeartbeatTasksCompletedOff, in.TasksCompleted)
	ob.SetInt32(workerHeartbeatTasksFailedOff, in.TasksFailed)
	ob.SetInt64(workerHeartbeatUptimeSecondsOff, in.UptimeSeconds)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- HeartbeatResponse ---
//
// Fixed payload: success(1) message(8 ptr).

const (
	heartbeatResponseSuccessOff = 0
	heartbeatResponseMessageOff = 8
	heartbeatResponseSize       = 16
)

// HeartbeatResponse is a zero-copy view into a ZAP-encoded HeartbeatResponse message.
type HeartbeatResponse struct{ o zap.Object }

// WrapHeartbeatResponse parses b and returns a typed view.
func WrapHeartbeatResponse(b []byte) (HeartbeatResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return HeartbeatResponse{}, err
	}
	return HeartbeatResponse{o: m.Root()}, nil
}

func (t HeartbeatResponse) Success() bool   { return t.o.Bool(heartbeatResponseSuccessOff) }
func (t HeartbeatResponse) Message() string { return t.o.Text(heartbeatResponseMessageOff) }

// HeartbeatResponseInput collects the field values for NewHeartbeatResponse.
type HeartbeatResponseInput struct {
	Success bool
	Message string
}

// NewHeartbeatResponse builds a ZAP-encoded HeartbeatResponse message from in.
func NewHeartbeatResponse(in HeartbeatResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(heartbeatResponseSize)
	ob.SetBool(heartbeatResponseSuccessOff, in.Success)
	ob.SetText(heartbeatResponseMessageOff, in.Message)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- TaskRequest ---
//
// Fixed payload: worker_id(8 ptr) capabilities(8 list ptr) available_slots(4).
// capabilities is repeated string.

const (
	taskRequestWorkerIDOff       = 0
	taskRequestCapabilitiesOff   = 8
	taskRequestAvailableSlotsOff = 16
	taskRequestSize              = 20
)

// TaskRequest is a zero-copy view into a ZAP-encoded TaskRequest message.
type TaskRequest struct{ o zap.Object }

// WrapTaskRequest parses b and returns a typed view.
func WrapTaskRequest(b []byte) (TaskRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return TaskRequest{}, err
	}
	return TaskRequest{o: m.Root()}, nil
}

func (t TaskRequest) WorkerID() string      { return t.o.Text(taskRequestWorkerIDOff) }
func (t TaskRequest) AvailableSlots() int32 { return t.o.Int32(taskRequestAvailableSlotsOff) }

// CapabilitiesLen returns the number of capabilities elements.
func (t TaskRequest) CapabilitiesLen() int { return t.o.List(taskRequestCapabilitiesOff).Len() }

// Capability returns the i-th capabilities element.
func (t TaskRequest) Capability(i int) string {
	return elemText(t.o.List(taskRequestCapabilitiesOff), i)
}

// TaskRequestInput collects the field values for NewTaskRequest.
type TaskRequestInput struct {
	WorkerID       string
	Capabilities   []string
	AvailableSlots int32
}

// NewTaskRequest builds a ZAP-encoded TaskRequest message from in.
func NewTaskRequest(in TaskRequestInput) []byte {
	b := zap.NewBuilder(256)
	capsOff, capsLen := 0, 0
	if len(in.Capabilities) > 0 {
		el := startElemList(b)
		for _, s := range in.Capabilities {
			el.addText(s)
		}
		capsOff, capsLen = el.finish()
	}
	ob := b.StartObject(taskRequestSize)
	ob.SetText(taskRequestWorkerIDOff, in.WorkerID)
	ob.SetList(taskRequestCapabilitiesOff, capsOff, capsLen)
	ob.SetInt32(taskRequestAvailableSlotsOff, in.AvailableSlots)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- TaskUpdate ---
//
// Fixed payload: task_id(8 ptr) worker_id(8 ptr) status(8 ptr) progress(4 float)
// message(8 ptr) metadata(8 list ptr). metadata is map<string,string>.

const (
	taskUpdateTaskIDOff   = 0
	taskUpdateWorkerIDOff = 8
	taskUpdateStatusOff   = 16
	taskUpdateProgressOff = 24
	taskUpdateMessageOff  = 28
	taskUpdateMetadataOff = 36
	taskUpdateSize        = 44
)

// TaskUpdate is a zero-copy view into a ZAP-encoded TaskUpdate message.
type TaskUpdate struct{ o zap.Object }

// WrapTaskUpdate parses b and returns a typed view.
func WrapTaskUpdate(b []byte) (TaskUpdate, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return TaskUpdate{}, err
	}
	return TaskUpdate{o: m.Root()}, nil
}

func (t TaskUpdate) TaskID() string    { return t.o.Text(taskUpdateTaskIDOff) }
func (t TaskUpdate) WorkerID() string  { return t.o.Text(taskUpdateWorkerIDOff) }
func (t TaskUpdate) Status() string    { return t.o.Text(taskUpdateStatusOff) }
func (t TaskUpdate) Progress() float32 { return t.o.Float32(taskUpdateProgressOff) }
func (t TaskUpdate) Message() string   { return t.o.Text(taskUpdateMessageOff) }

// MetadataLen returns the number of metadata entries.
func (t TaskUpdate) MetadataLen() int { return t.o.List(taskUpdateMetadataOff).Len() }

// Metadata returns the i-th metadata entry (key/value).
func (t TaskUpdate) Metadata(i int) MapEntry {
	return mapEntryAt(t.o.List(taskUpdateMetadataOff), i)
}

// TaskUpdateInput collects the field values for NewTaskUpdate.
type TaskUpdateInput struct {
	TaskID   string
	WorkerID string
	Status   string
	Progress float32
	Message  string
	Metadata []MapEntry
}

// NewTaskUpdate builds a ZAP-encoded TaskUpdate message from in.
func NewTaskUpdate(in TaskUpdateInput) []byte {
	b := zap.NewBuilder(256)
	metaOff, metaLen := 0, 0
	if len(in.Metadata) > 0 {
		el := startElemList(b)
		for _, e := range in.Metadata {
			el.addStringEntry(e.Key, e.Value)
		}
		metaOff, metaLen = el.finish()
	}
	ob := b.StartObject(taskUpdateSize)
	ob.SetText(taskUpdateTaskIDOff, in.TaskID)
	ob.SetText(taskUpdateWorkerIDOff, in.WorkerID)
	ob.SetText(taskUpdateStatusOff, in.Status)
	ob.SetFloat32(taskUpdateProgressOff, in.Progress)
	ob.SetText(taskUpdateMessageOff, in.Message)
	ob.SetList(taskUpdateMetadataOff, metaOff, metaLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- TaskComplete ---
//
// Fixed payload: task_id(8 ptr) worker_id(8 ptr) success(1) error_message(8 ptr)
// completion_time(8) result_metadata(8 list ptr). result_metadata is
// map<string,string>.

const (
	taskCompleteTaskIDOff         = 0
	taskCompleteWorkerIDOff       = 8
	taskCompleteSuccessOff        = 16
	taskCompleteErrorMessageOff   = 24
	taskCompleteCompletionTimeOff = 32
	taskCompleteResultMetadataOff = 40
	taskCompleteSize              = 48
)

// TaskComplete is a zero-copy view into a ZAP-encoded TaskComplete message.
type TaskComplete struct{ o zap.Object }

// WrapTaskComplete parses b and returns a typed view.
func WrapTaskComplete(b []byte) (TaskComplete, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return TaskComplete{}, err
	}
	return TaskComplete{o: m.Root()}, nil
}

func (t TaskComplete) TaskID() string       { return t.o.Text(taskCompleteTaskIDOff) }
func (t TaskComplete) WorkerID() string     { return t.o.Text(taskCompleteWorkerIDOff) }
func (t TaskComplete) Success() bool        { return t.o.Bool(taskCompleteSuccessOff) }
func (t TaskComplete) ErrorMessage() string { return t.o.Text(taskCompleteErrorMessageOff) }
func (t TaskComplete) CompletionTime() int64 {
	return t.o.Int64(taskCompleteCompletionTimeOff)
}

// ResultMetadataLen returns the number of result_metadata entries.
func (t TaskComplete) ResultMetadataLen() int {
	return t.o.List(taskCompleteResultMetadataOff).Len()
}

// ResultMetadata returns the i-th result_metadata entry (key/value).
func (t TaskComplete) ResultMetadata(i int) MapEntry {
	return mapEntryAt(t.o.List(taskCompleteResultMetadataOff), i)
}

// TaskCompleteInput collects the field values for NewTaskComplete.
type TaskCompleteInput struct {
	TaskID         string
	WorkerID       string
	Success        bool
	ErrorMessage   string
	CompletionTime int64
	ResultMetadata []MapEntry
}

// NewTaskComplete builds a ZAP-encoded TaskComplete message from in.
func NewTaskComplete(in TaskCompleteInput) []byte {
	b := zap.NewBuilder(256)
	metaOff, metaLen := 0, 0
	if len(in.ResultMetadata) > 0 {
		el := startElemList(b)
		for _, e := range in.ResultMetadata {
			el.addStringEntry(e.Key, e.Value)
		}
		metaOff, metaLen = el.finish()
	}
	ob := b.StartObject(taskCompleteSize)
	ob.SetText(taskCompleteTaskIDOff, in.TaskID)
	ob.SetText(taskCompleteWorkerIDOff, in.WorkerID)
	ob.SetBool(taskCompleteSuccessOff, in.Success)
	ob.SetText(taskCompleteErrorMessageOff, in.ErrorMessage)
	ob.SetInt64(taskCompleteCompletionTimeOff, in.CompletionTime)
	ob.SetList(taskCompleteResultMetadataOff, metaOff, metaLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- TaskCancellation ---
//
// Fixed payload: task_id(8 ptr) reason(8 ptr) force(1).

const (
	taskCancellationTaskIDOff = 0
	taskCancellationReasonOff = 8
	taskCancellationForceOff  = 16
	taskCancellationSize      = 17
)

// TaskCancellation is a zero-copy view into a ZAP-encoded TaskCancellation message.
type TaskCancellation struct{ o zap.Object }

// WrapTaskCancellation parses b and returns a typed view.
func WrapTaskCancellation(b []byte) (TaskCancellation, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return TaskCancellation{}, err
	}
	return TaskCancellation{o: m.Root()}, nil
}

func (t TaskCancellation) TaskID() string { return t.o.Text(taskCancellationTaskIDOff) }
func (t TaskCancellation) Reason() string { return t.o.Text(taskCancellationReasonOff) }
func (t TaskCancellation) Force() bool    { return t.o.Bool(taskCancellationForceOff) }

// TaskCancellationInput collects the field values for NewTaskCancellation.
type TaskCancellationInput struct {
	TaskID string
	Reason string
	Force  bool
}

// NewTaskCancellation builds a ZAP-encoded TaskCancellation message from in.
func NewTaskCancellation(in TaskCancellationInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(taskCancellationSize)
	ob.SetText(taskCancellationTaskIDOff, in.TaskID)
	ob.SetText(taskCancellationReasonOff, in.Reason)
	ob.SetBool(taskCancellationForceOff, in.Force)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- WorkerShutdown ---
//
// Fixed payload: worker_id(8 ptr) reason(8 ptr) pending_task_ids(8 list ptr).
// pending_task_ids is repeated string.

const (
	workerShutdownWorkerIDOff       = 0
	workerShutdownReasonOff         = 8
	workerShutdownPendingTaskIDsOff = 16
	workerShutdownSize              = 24
)

// WorkerShutdown is a zero-copy view into a ZAP-encoded WorkerShutdown message.
type WorkerShutdown struct{ o zap.Object }

// WrapWorkerShutdown parses b and returns a typed view.
func WrapWorkerShutdown(b []byte) (WorkerShutdown, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return WorkerShutdown{}, err
	}
	return WorkerShutdown{o: m.Root()}, nil
}

func (t WorkerShutdown) WorkerID() string { return t.o.Text(workerShutdownWorkerIDOff) }
func (t WorkerShutdown) Reason() string   { return t.o.Text(workerShutdownReasonOff) }

// PendingTaskIDsLen returns the number of pending_task_ids elements.
func (t WorkerShutdown) PendingTaskIDsLen() int {
	return t.o.List(workerShutdownPendingTaskIDsOff).Len()
}

// PendingTaskID returns the i-th pending_task_ids element.
func (t WorkerShutdown) PendingTaskID(i int) string {
	return elemText(t.o.List(workerShutdownPendingTaskIDsOff), i)
}

// WorkerShutdownInput collects the field values for NewWorkerShutdown.
type WorkerShutdownInput struct {
	WorkerID       string
	Reason         string
	PendingTaskIDs []string
}

// NewWorkerShutdown builds a ZAP-encoded WorkerShutdown message from in.
func NewWorkerShutdown(in WorkerShutdownInput) []byte {
	b := zap.NewBuilder(256)
	pendOff, pendLen := 0, 0
	if len(in.PendingTaskIDs) > 0 {
		el := startElemList(b)
		for _, s := range in.PendingTaskIDs {
			el.addText(s)
		}
		pendOff, pendLen = el.finish()
	}
	ob := b.StartObject(workerShutdownSize)
	ob.SetText(workerShutdownWorkerIDOff, in.WorkerID)
	ob.SetText(workerShutdownReasonOff, in.Reason)
	ob.SetList(workerShutdownPendingTaskIDsOff, pendOff, pendLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- AdminShutdown ---
//
// Fixed payload: reason(8 ptr) graceful_shutdown_seconds(4).

const (
	adminShutdownReasonOff                  = 0
	adminShutdownGracefulShutdownSecondsOff = 8
	adminShutdownSize                       = 12
)

// AdminShutdown is a zero-copy view into a ZAP-encoded AdminShutdown message.
type AdminShutdown struct{ o zap.Object }

// WrapAdminShutdown parses b and returns a typed view.
func WrapAdminShutdown(b []byte) (AdminShutdown, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return AdminShutdown{}, err
	}
	return AdminShutdown{o: m.Root()}, nil
}

func (t AdminShutdown) Reason() string { return t.o.Text(adminShutdownReasonOff) }
func (t AdminShutdown) GracefulShutdownSeconds() int32 {
	return t.o.Int32(adminShutdownGracefulShutdownSecondsOff)
}

// AdminShutdownInput collects the field values for NewAdminShutdown.
type AdminShutdownInput struct {
	Reason                  string
	GracefulShutdownSeconds int32
}

// NewAdminShutdown builds a ZAP-encoded AdminShutdown message from in.
func NewAdminShutdown(in AdminShutdownInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(adminShutdownSize)
	ob.SetText(adminShutdownReasonOff, in.Reason)
	ob.SetInt32(adminShutdownGracefulShutdownSecondsOff, in.GracefulShutdownSeconds)
	ob.FinishAsRoot()
	return b.Finish()
}
