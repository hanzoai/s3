// Code generated from plugin.proto; DO NOT EDIT.

package pluginwire

import (
	zap "github.com/zap-proto/go"
)

// This file carries the worker lifecycle cluster: WorkerHello, AdminHello,
// WorkerHeartbeat, WorkerAcknowledge, RunningWork, JobTypeCapability.

// --- WorkerHello ---

const (
	workerHelloWorkerIDOff         = 0
	workerHelloWorkerInstanceIDOff = 8
	workerHelloAddressOff          = 16
	workerHelloWorkerVersionOff    = 24
	workerHelloProtocolVersionOff  = 32
	workerHelloCapabilitiesOff     = 40
	workerHelloMetadataOff         = 48
	workerHelloSize                = 56
)

// WorkerHello is a zero-copy view into a ZAP-encoded WorkerHello message.
type WorkerHello struct{ o zap.Object }

// WrapWorkerHello parses b and returns a typed view.
func WrapWorkerHello(b []byte) (WorkerHello, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return WorkerHello{}, err
	}
	return WorkerHello{o: m.Root()}, nil
}

func (t WorkerHello) WorkerID() string         { return t.o.Text(workerHelloWorkerIDOff) }
func (t WorkerHello) WorkerInstanceID() string { return t.o.Text(workerHelloWorkerInstanceIDOff) }
func (t WorkerHello) Address() string          { return t.o.Text(workerHelloAddressOff) }
func (t WorkerHello) WorkerVersion() string    { return t.o.Text(workerHelloWorkerVersionOff) }
func (t WorkerHello) ProtocolVersion() string  { return t.o.Text(workerHelloProtocolVersionOff) }

// Capabilities reads capabilities (proto field 6, repeated JobTypeCapability);
// each element is a JobTypeCapability object.
func (t WorkerHello) Capabilities() zap.List { return t.o.List(workerHelloCapabilitiesOff) }

// Metadata reads metadata (proto field 7, map<string,string>); each element is
// a StringMapEntry object.
func (t WorkerHello) Metadata() zap.List { return t.o.List(workerHelloMetadataOff) }

// WorkerHelloInput collects the field values for NewWorkerHello.
type WorkerHelloInput struct {
	WorkerID         string
	WorkerInstanceID string
	Address          string
	WorkerVersion    string
	ProtocolVersion  string
	Capabilities     [][]byte
	Metadata         [][]byte
}

// NewWorkerHello builds a ZAP-encoded WorkerHello message.
func NewWorkerHello(in WorkerHelloInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(workerHelloSize)
	ob.SetText(workerHelloWorkerIDOff, in.WorkerID)
	ob.SetText(workerHelloWorkerInstanceIDOff, in.WorkerInstanceID)
	ob.SetText(workerHelloAddressOff, in.Address)
	ob.SetText(workerHelloWorkerVersionOff, in.WorkerVersion)
	ob.SetText(workerHelloProtocolVersionOff, in.ProtocolVersion)
	cb := b.StartList(0)
	for _, elem := range in.Capabilities {
		cb.AddObjectBytes(elem)
	}
	ob.SetList(workerHelloCapabilitiesOff, cb.FinishOffset(), len(in.Capabilities))
	mb := b.StartList(0)
	for _, elem := range in.Metadata {
		mb.AddObjectBytes(elem)
	}
	ob.SetList(workerHelloMetadataOff, mb.FinishOffset(), len(in.Metadata))
	ob.FinishAsRoot()
	return b.Finish()
}

// --- AdminHello ---

const (
	adminHelloAcceptedOff                 = 0
	adminHelloMessageOff                  = 8
	adminHelloHeartbeatIntervalSecondsOff = 16
	adminHelloReconnectDelaySecondsOff    = 20
	adminHelloSize                        = 24
)

// AdminHello is a zero-copy view into a ZAP-encoded AdminHello message.
type AdminHello struct{ o zap.Object }

// WrapAdminHello parses b and returns a typed view.
func WrapAdminHello(b []byte) (AdminHello, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return AdminHello{}, err
	}
	return AdminHello{o: m.Root()}, nil
}

func (t AdminHello) Accepted() bool  { return t.o.Bool(adminHelloAcceptedOff) }
func (t AdminHello) Message() string { return t.o.Text(adminHelloMessageOff) }
func (t AdminHello) HeartbeatIntervalSeconds() int32 {
	return t.o.Int32(adminHelloHeartbeatIntervalSecondsOff)
}
func (t AdminHello) ReconnectDelaySeconds() int32 {
	return t.o.Int32(adminHelloReconnectDelaySecondsOff)
}

// AdminHelloInput collects the field values for NewAdminHello.
type AdminHelloInput struct {
	Accepted                 bool
	Message                  string
	HeartbeatIntervalSeconds int32
	ReconnectDelaySeconds    int32
}

// NewAdminHello builds a ZAP-encoded AdminHello message.
func NewAdminHello(in AdminHelloInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(adminHelloSize)
	ob.SetBool(adminHelloAcceptedOff, in.Accepted)
	ob.SetText(adminHelloMessageOff, in.Message)
	ob.SetInt32(adminHelloHeartbeatIntervalSecondsOff, in.HeartbeatIntervalSeconds)
	ob.SetInt32(adminHelloReconnectDelaySecondsOff, in.ReconnectDelaySeconds)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- WorkerHeartbeat ---

const (
	workerHeartbeatWorkerIDOff            = 0
	workerHeartbeatRunningWorkOff         = 8
	workerHeartbeatDetectionSlotsUsedOff  = 16
	workerHeartbeatDetectionSlotsTotalOff = 20
	workerHeartbeatExecutionSlotsUsedOff  = 24
	workerHeartbeatExecutionSlotsTotalOff = 28
	workerHeartbeatQueuedJobsByTypeOff    = 32
	workerHeartbeatMetadataOff            = 40
	workerHeartbeatSize                   = 48
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

func (t WorkerHeartbeat) WorkerID() string { return t.o.Text(workerHeartbeatWorkerIDOff) }

// RunningWork reads running_work (proto field 2, repeated RunningWork); each
// element is a RunningWork object.
func (t WorkerHeartbeat) RunningWork() zap.List { return t.o.List(workerHeartbeatRunningWorkOff) }

func (t WorkerHeartbeat) DetectionSlotsUsed() int32 {
	return t.o.Int32(workerHeartbeatDetectionSlotsUsedOff)
}
func (t WorkerHeartbeat) DetectionSlotsTotal() int32 {
	return t.o.Int32(workerHeartbeatDetectionSlotsTotalOff)
}
func (t WorkerHeartbeat) ExecutionSlotsUsed() int32 {
	return t.o.Int32(workerHeartbeatExecutionSlotsUsedOff)
}
func (t WorkerHeartbeat) ExecutionSlotsTotal() int32 {
	return t.o.Int32(workerHeartbeatExecutionSlotsTotalOff)
}

// QueuedJobsByType reads queued_jobs_by_type (proto field 7, map<string,int32>);
// each element is an Int32MapEntry object.
func (t WorkerHeartbeat) QueuedJobsByType() zap.List {
	return t.o.List(workerHeartbeatQueuedJobsByTypeOff)
}

// Metadata reads metadata (proto field 8, map<string,string>); each element is
// a StringMapEntry object.
func (t WorkerHeartbeat) Metadata() zap.List { return t.o.List(workerHeartbeatMetadataOff) }

// WorkerHeartbeatInput collects the field values for NewWorkerHeartbeat.
type WorkerHeartbeatInput struct {
	WorkerID            string
	RunningWork         [][]byte
	DetectionSlotsUsed  int32
	DetectionSlotsTotal int32
	ExecutionSlotsUsed  int32
	ExecutionSlotsTotal int32
	QueuedJobsByType    [][]byte
	Metadata            [][]byte
}

// NewWorkerHeartbeat builds a ZAP-encoded WorkerHeartbeat message.
func NewWorkerHeartbeat(in WorkerHeartbeatInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(workerHeartbeatSize)
	ob.SetText(workerHeartbeatWorkerIDOff, in.WorkerID)
	rb := b.StartList(0)
	for _, elem := range in.RunningWork {
		rb.AddObjectBytes(elem)
	}
	ob.SetList(workerHeartbeatRunningWorkOff, rb.FinishOffset(), len(in.RunningWork))
	ob.SetInt32(workerHeartbeatDetectionSlotsUsedOff, in.DetectionSlotsUsed)
	ob.SetInt32(workerHeartbeatDetectionSlotsTotalOff, in.DetectionSlotsTotal)
	ob.SetInt32(workerHeartbeatExecutionSlotsUsedOff, in.ExecutionSlotsUsed)
	ob.SetInt32(workerHeartbeatExecutionSlotsTotalOff, in.ExecutionSlotsTotal)
	qb := b.StartList(0)
	for _, elem := range in.QueuedJobsByType {
		qb.AddObjectBytes(elem)
	}
	ob.SetList(workerHeartbeatQueuedJobsByTypeOff, qb.FinishOffset(), len(in.QueuedJobsByType))
	mb := b.StartList(0)
	for _, elem := range in.Metadata {
		mb.AddObjectBytes(elem)
	}
	ob.SetList(workerHeartbeatMetadataOff, mb.FinishOffset(), len(in.Metadata))
	ob.FinishAsRoot()
	return b.Finish()
}

// --- WorkerAcknowledge ---

const (
	workerAcknowledgeRequestIDOff = 0
	workerAcknowledgeAcceptedOff  = 8
	workerAcknowledgeMessageOff   = 16
	workerAcknowledgeSize         = 24
)

// WorkerAcknowledge is a zero-copy view into a ZAP-encoded WorkerAcknowledge.
type WorkerAcknowledge struct{ o zap.Object }

// WrapWorkerAcknowledge parses b and returns a typed view.
func WrapWorkerAcknowledge(b []byte) (WorkerAcknowledge, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return WorkerAcknowledge{}, err
	}
	return WorkerAcknowledge{o: m.Root()}, nil
}

func (t WorkerAcknowledge) RequestID() string { return t.o.Text(workerAcknowledgeRequestIDOff) }
func (t WorkerAcknowledge) Accepted() bool    { return t.o.Bool(workerAcknowledgeAcceptedOff) }
func (t WorkerAcknowledge) Message() string   { return t.o.Text(workerAcknowledgeMessageOff) }

// WorkerAcknowledgeInput collects the field values for NewWorkerAcknowledge.
type WorkerAcknowledgeInput struct {
	RequestID string
	Accepted  bool
	Message   string
}

// NewWorkerAcknowledge builds a ZAP-encoded WorkerAcknowledge message.
func NewWorkerAcknowledge(in WorkerAcknowledgeInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(workerAcknowledgeSize)
	ob.SetText(workerAcknowledgeRequestIDOff, in.RequestID)
	ob.SetBool(workerAcknowledgeAcceptedOff, in.Accepted)
	ob.SetText(workerAcknowledgeMessageOff, in.Message)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- RunningWork ---

const (
	runningWorkWorkIDOff          = 0
	runningWorkKindOff            = 8
	runningWorkJobTypeOff         = 12
	runningWorkStateOff           = 20
	runningWorkProgressPercentOff = 24
	runningWorkStageOff           = 32
	runningWorkSize               = 40
)

// RunningWork is a zero-copy view into a ZAP-encoded RunningWork message.
type RunningWork struct{ o zap.Object }

// WrapRunningWork parses b and returns a typed view.
func WrapRunningWork(b []byte) (RunningWork, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return RunningWork{}, err
	}
	return RunningWork{o: m.Root()}, nil
}

func (t RunningWork) WorkID() string { return t.o.Text(runningWorkWorkIDOff) }

// Kind reads kind (proto field 2, enum WorkKind).
func (t RunningWork) Kind() uint32    { return t.o.Uint32(runningWorkKindOff) }
func (t RunningWork) JobType() string { return t.o.Text(runningWorkJobTypeOff) }

// State reads state (proto field 4, enum JobState).
func (t RunningWork) State() uint32            { return t.o.Uint32(runningWorkStateOff) }
func (t RunningWork) ProgressPercent() float64 { return t.o.Float64(runningWorkProgressPercentOff) }
func (t RunningWork) Stage() string            { return t.o.Text(runningWorkStageOff) }

// RunningWorkInput collects the field values for NewRunningWork.
type RunningWorkInput struct {
	WorkID          string
	Kind            uint32
	JobType         string
	State           uint32
	ProgressPercent float64
	Stage           string
}

// NewRunningWork builds a ZAP-encoded RunningWork message.
func NewRunningWork(in RunningWorkInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(runningWorkSize)
	ob.SetText(runningWorkWorkIDOff, in.WorkID)
	ob.SetUint32(runningWorkKindOff, in.Kind)
	ob.SetText(runningWorkJobTypeOff, in.JobType)
	ob.SetUint32(runningWorkStateOff, in.State)
	ob.SetFloat64(runningWorkProgressPercentOff, in.ProgressPercent)
	ob.SetText(runningWorkStageOff, in.Stage)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- JobTypeCapability ---

const (
	jobTypeCapabilityJobTypeOff                 = 0
	jobTypeCapabilityCanDetectOff               = 8
	jobTypeCapabilityCanExecuteOff              = 9
	jobTypeCapabilityMaxDetectionConcurrencyOff = 12
	jobTypeCapabilityMaxExecutionConcurrencyOff = 16
	jobTypeCapabilityDisplayNameOff             = 20
	jobTypeCapabilityDescriptionOff             = 28
	jobTypeCapabilityWeightOff                  = 36
	jobTypeCapabilitySize                       = 40
)

// JobTypeCapability is a zero-copy view into a ZAP-encoded JobTypeCapability.
type JobTypeCapability struct{ o zap.Object }

// WrapJobTypeCapability parses b and returns a typed view.
func WrapJobTypeCapability(b []byte) (JobTypeCapability, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return JobTypeCapability{}, err
	}
	return JobTypeCapability{o: m.Root()}, nil
}

func (t JobTypeCapability) JobType() string  { return t.o.Text(jobTypeCapabilityJobTypeOff) }
func (t JobTypeCapability) CanDetect() bool  { return t.o.Bool(jobTypeCapabilityCanDetectOff) }
func (t JobTypeCapability) CanExecute() bool { return t.o.Bool(jobTypeCapabilityCanExecuteOff) }
func (t JobTypeCapability) MaxDetectionConcurrency() int32 {
	return t.o.Int32(jobTypeCapabilityMaxDetectionConcurrencyOff)
}
func (t JobTypeCapability) MaxExecutionConcurrency() int32 {
	return t.o.Int32(jobTypeCapabilityMaxExecutionConcurrencyOff)
}
func (t JobTypeCapability) DisplayName() string { return t.o.Text(jobTypeCapabilityDisplayNameOff) }
func (t JobTypeCapability) Description() string { return t.o.Text(jobTypeCapabilityDescriptionOff) }
func (t JobTypeCapability) Weight() int32       { return t.o.Int32(jobTypeCapabilityWeightOff) }

// JobTypeCapabilityInput collects the field values for NewJobTypeCapability.
type JobTypeCapabilityInput struct {
	JobType                 string
	CanDetect               bool
	CanExecute              bool
	MaxDetectionConcurrency int32
	MaxExecutionConcurrency int32
	DisplayName             string
	Description             string
	Weight                  int32
}

// NewJobTypeCapability builds a ZAP-encoded JobTypeCapability message.
func NewJobTypeCapability(in JobTypeCapabilityInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(jobTypeCapabilitySize)
	ob.SetText(jobTypeCapabilityJobTypeOff, in.JobType)
	ob.SetBool(jobTypeCapabilityCanDetectOff, in.CanDetect)
	ob.SetBool(jobTypeCapabilityCanExecuteOff, in.CanExecute)
	ob.SetInt32(jobTypeCapabilityMaxDetectionConcurrencyOff, in.MaxDetectionConcurrency)
	ob.SetInt32(jobTypeCapabilityMaxExecutionConcurrencyOff, in.MaxExecutionConcurrency)
	ob.SetText(jobTypeCapabilityDisplayNameOff, in.DisplayName)
	ob.SetText(jobTypeCapabilityDescriptionOff, in.Description)
	ob.SetInt32(jobTypeCapabilityWeightOff, in.Weight)
	ob.FinishAsRoot()
	return b.Finish()
}
