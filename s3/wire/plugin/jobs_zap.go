// Code generated from plugin.proto; DO NOT EDIT.

package pluginwire

import (
	zap "github.com/zap-proto/go"
)

// This file carries the detection / execution / job-lifecycle cluster plus the
// cluster-context and persisted-config messages. google.protobuf.Timestamp and
// google.protobuf.Duration fields are encoded as int64 unix-nanoseconds.

// --- RequestConfigSchema ---

const (
	requestConfigSchemaJobTypeOff      = 0
	requestConfigSchemaForceRefreshOff = 8
	requestConfigSchemaSize            = 9
)

// RequestConfigSchema is a zero-copy view into a ZAP-encoded RequestConfigSchema.
type RequestConfigSchema struct{ o zap.Object }

// WrapRequestConfigSchema parses b and returns a typed view.
func WrapRequestConfigSchema(b []byte) (RequestConfigSchema, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return RequestConfigSchema{}, err
	}
	return RequestConfigSchema{o: m.Root()}, nil
}

func (t RequestConfigSchema) JobType() string    { return t.o.Text(requestConfigSchemaJobTypeOff) }
func (t RequestConfigSchema) ForceRefresh() bool { return t.o.Bool(requestConfigSchemaForceRefreshOff) }

// RequestConfigSchemaInput collects the field values for NewRequestConfigSchema.
type RequestConfigSchemaInput struct {
	JobType      string
	ForceRefresh bool
}

// NewRequestConfigSchema builds a ZAP-encoded RequestConfigSchema message.
func NewRequestConfigSchema(in RequestConfigSchemaInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(requestConfigSchemaSize)
	ob.SetText(requestConfigSchemaJobTypeOff, in.JobType)
	ob.SetBool(requestConfigSchemaForceRefreshOff, in.ForceRefresh)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ConfigSchemaResponse ---

const (
	configSchemaResponseRequestIDOff         = 0
	configSchemaResponseJobTypeOff           = 8
	configSchemaResponseSuccessOff           = 16
	configSchemaResponseErrorMessageOff      = 24
	configSchemaResponseJobTypeDescriptorOff = 32
	configSchemaResponseSize                 = 40
)

// ConfigSchemaResponse is a zero-copy view into a ZAP-encoded ConfigSchemaResponse.
type ConfigSchemaResponse struct{ o zap.Object }

// WrapConfigSchemaResponse parses b and returns a typed view.
func WrapConfigSchemaResponse(b []byte) (ConfigSchemaResponse, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ConfigSchemaResponse{}, err
	}
	return ConfigSchemaResponse{o: m.Root()}, nil
}

func (t ConfigSchemaResponse) RequestID() string { return t.o.Text(configSchemaResponseRequestIDOff) }
func (t ConfigSchemaResponse) JobType() string   { return t.o.Text(configSchemaResponseJobTypeOff) }
func (t ConfigSchemaResponse) Success() bool     { return t.o.Bool(configSchemaResponseSuccessOff) }
func (t ConfigSchemaResponse) ErrorMessage() string {
	return t.o.Text(configSchemaResponseErrorMessageOff)
}
func (t ConfigSchemaResponse) JobTypeDescriptor() JobTypeDescriptor {
	return JobTypeDescriptor{o: childObject(t.o.Bytes(configSchemaResponseJobTypeDescriptorOff))}
}

// ConfigSchemaResponseInput collects the field values for NewConfigSchemaResponse.
// JobTypeDescriptor is a pre-built sub-buffer.
type ConfigSchemaResponseInput struct {
	RequestID         string
	JobType           string
	Success           bool
	ErrorMessage      string
	JobTypeDescriptor []byte
}

// NewConfigSchemaResponse builds a ZAP-encoded ConfigSchemaResponse message.
func NewConfigSchemaResponse(in ConfigSchemaResponseInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(configSchemaResponseSize)
	ob.SetText(configSchemaResponseRequestIDOff, in.RequestID)
	ob.SetText(configSchemaResponseJobTypeOff, in.JobType)
	ob.SetBool(configSchemaResponseSuccessOff, in.Success)
	ob.SetText(configSchemaResponseErrorMessageOff, in.ErrorMessage)
	setChild(ob, configSchemaResponseJobTypeDescriptorOff, in.JobTypeDescriptor)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- RunDetectionRequest ---

const (
	runDetectionRequestRequestIDOff          = 0
	runDetectionRequestJobTypeOff            = 8
	runDetectionRequestDetectionSequenceOff  = 16
	runDetectionRequestAdminRuntimeOff       = 24
	runDetectionRequestAdminConfigValuesOff  = 32
	runDetectionRequestWorkerConfigValuesOff = 40
	runDetectionRequestClusterContextOff     = 48
	runDetectionRequestLastSuccessfulRunOff  = 56
	runDetectionRequestMaxResultsOff         = 64
	runDetectionRequestSize                  = 68
)

// RunDetectionRequest is a zero-copy view into a ZAP-encoded RunDetectionRequest.
type RunDetectionRequest struct{ o zap.Object }

// WrapRunDetectionRequest parses b and returns a typed view.
func WrapRunDetectionRequest(b []byte) (RunDetectionRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return RunDetectionRequest{}, err
	}
	return RunDetectionRequest{o: m.Root()}, nil
}

func (t RunDetectionRequest) RequestID() string { return t.o.Text(runDetectionRequestRequestIDOff) }
func (t RunDetectionRequest) JobType() string   { return t.o.Text(runDetectionRequestJobTypeOff) }
func (t RunDetectionRequest) DetectionSequence() int64 {
	return t.o.Int64(runDetectionRequestDetectionSequenceOff)
}
func (t RunDetectionRequest) AdminRuntime() AdminRuntimeConfig {
	return AdminRuntimeConfig{o: childObject(t.o.Bytes(runDetectionRequestAdminRuntimeOff))}
}

// AdminConfigValues reads admin_config_values (proto field 5,
// map<string,ConfigValue>); each element is a ConfigValueMapEntry object.
func (t RunDetectionRequest) AdminConfigValues() zap.List {
	return t.o.List(runDetectionRequestAdminConfigValuesOff)
}

// WorkerConfigValues reads worker_config_values (proto field 6,
// map<string,ConfigValue>); each element is a ConfigValueMapEntry object.
func (t RunDetectionRequest) WorkerConfigValues() zap.List {
	return t.o.List(runDetectionRequestWorkerConfigValuesOff)
}
func (t RunDetectionRequest) ClusterContext() ClusterContext {
	return ClusterContext{o: childObject(t.o.Bytes(runDetectionRequestClusterContextOff))}
}

// LastSuccessfulRun reads last_successful_run (proto field 8,
// google.protobuf.Timestamp) as int64 unix-nanoseconds.
func (t RunDetectionRequest) LastSuccessfulRun() int64 {
	return t.o.Int64(runDetectionRequestLastSuccessfulRunOff)
}
func (t RunDetectionRequest) MaxResults() int32 { return t.o.Int32(runDetectionRequestMaxResultsOff) }

// RunDetectionRequestInput collects the field values for NewRunDetectionRequest.
type RunDetectionRequestInput struct {
	RequestID          string
	JobType            string
	DetectionSequence  int64
	AdminRuntime       []byte
	AdminConfigValues  [][]byte
	WorkerConfigValues [][]byte
	ClusterContext     []byte
	LastSuccessfulRun  int64
	MaxResults         int32
}

// NewRunDetectionRequest builds a ZAP-encoded RunDetectionRequest message.
func NewRunDetectionRequest(in RunDetectionRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(runDetectionRequestSize)
	ob.SetText(runDetectionRequestRequestIDOff, in.RequestID)
	ob.SetText(runDetectionRequestJobTypeOff, in.JobType)
	ob.SetInt64(runDetectionRequestDetectionSequenceOff, in.DetectionSequence)
	setChild(ob, runDetectionRequestAdminRuntimeOff, in.AdminRuntime)
	ab := b.StartList(0)
	for _, elem := range in.AdminConfigValues {
		ab.AddObjectBytes(elem)
	}
	ob.SetList(runDetectionRequestAdminConfigValuesOff, ab.FinishOffset(), len(in.AdminConfigValues))
	wb := b.StartList(0)
	for _, elem := range in.WorkerConfigValues {
		wb.AddObjectBytes(elem)
	}
	ob.SetList(runDetectionRequestWorkerConfigValuesOff, wb.FinishOffset(), len(in.WorkerConfigValues))
	setChild(ob, runDetectionRequestClusterContextOff, in.ClusterContext)
	ob.SetInt64(runDetectionRequestLastSuccessfulRunOff, in.LastSuccessfulRun)
	ob.SetInt32(runDetectionRequestMaxResultsOff, in.MaxResults)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- DetectionProposals ---

const (
	detectionProposalsRequestIDOff = 0
	detectionProposalsJobTypeOff   = 8
	detectionProposalsProposalsOff = 16
	detectionProposalsHasMoreOff   = 24
	detectionProposalsSize         = 25
)

// DetectionProposals is a zero-copy view into a ZAP-encoded DetectionProposals.
type DetectionProposals struct{ o zap.Object }

// WrapDetectionProposals parses b and returns a typed view.
func WrapDetectionProposals(b []byte) (DetectionProposals, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return DetectionProposals{}, err
	}
	return DetectionProposals{o: m.Root()}, nil
}

func (t DetectionProposals) RequestID() string { return t.o.Text(detectionProposalsRequestIDOff) }
func (t DetectionProposals) JobType() string   { return t.o.Text(detectionProposalsJobTypeOff) }

// Proposals reads proposals (proto field 3, repeated JobProposal); each element
// is a JobProposal object.
func (t DetectionProposals) Proposals() zap.List { return t.o.List(detectionProposalsProposalsOff) }
func (t DetectionProposals) HasMore() bool       { return t.o.Bool(detectionProposalsHasMoreOff) }

// DetectionProposalsInput collects the field values for NewDetectionProposals.
type DetectionProposalsInput struct {
	RequestID string
	JobType   string
	Proposals [][]byte
	HasMore   bool
}

// NewDetectionProposals builds a ZAP-encoded DetectionProposals message.
func NewDetectionProposals(in DetectionProposalsInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(detectionProposalsSize)
	ob.SetText(detectionProposalsRequestIDOff, in.RequestID)
	ob.SetText(detectionProposalsJobTypeOff, in.JobType)
	pb := b.StartList(0)
	for _, elem := range in.Proposals {
		pb.AddObjectBytes(elem)
	}
	ob.SetList(detectionProposalsProposalsOff, pb.FinishOffset(), len(in.Proposals))
	ob.SetBool(detectionProposalsHasMoreOff, in.HasMore)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- DetectionComplete ---

const (
	detectionCompleteRequestIDOff      = 0
	detectionCompleteJobTypeOff        = 8
	detectionCompleteSuccessOff        = 16
	detectionCompleteErrorMessageOff   = 24
	detectionCompleteTotalProposalsOff = 32
	detectionCompleteSize              = 36
)

// DetectionComplete is a zero-copy view into a ZAP-encoded DetectionComplete.
type DetectionComplete struct{ o zap.Object }

// WrapDetectionComplete parses b and returns a typed view.
func WrapDetectionComplete(b []byte) (DetectionComplete, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return DetectionComplete{}, err
	}
	return DetectionComplete{o: m.Root()}, nil
}

func (t DetectionComplete) RequestID() string    { return t.o.Text(detectionCompleteRequestIDOff) }
func (t DetectionComplete) JobType() string      { return t.o.Text(detectionCompleteJobTypeOff) }
func (t DetectionComplete) Success() bool        { return t.o.Bool(detectionCompleteSuccessOff) }
func (t DetectionComplete) ErrorMessage() string { return t.o.Text(detectionCompleteErrorMessageOff) }
func (t DetectionComplete) TotalProposals() int32 {
	return t.o.Int32(detectionCompleteTotalProposalsOff)
}

// DetectionCompleteInput collects the field values for NewDetectionComplete.
type DetectionCompleteInput struct {
	RequestID      string
	JobType        string
	Success        bool
	ErrorMessage   string
	TotalProposals int32
}

// NewDetectionComplete builds a ZAP-encoded DetectionComplete message.
func NewDetectionComplete(in DetectionCompleteInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(detectionCompleteSize)
	ob.SetText(detectionCompleteRequestIDOff, in.RequestID)
	ob.SetText(detectionCompleteJobTypeOff, in.JobType)
	ob.SetBool(detectionCompleteSuccessOff, in.Success)
	ob.SetText(detectionCompleteErrorMessageOff, in.ErrorMessage)
	ob.SetInt32(detectionCompleteTotalProposalsOff, in.TotalProposals)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- JobProposal ---

const (
	jobProposalProposalIDOff = 0
	jobProposalDedupeKeyOff  = 8
	jobProposalJobTypeOff    = 16
	jobProposalPriorityOff   = 24
	jobProposalSummaryOff    = 28
	jobProposalDetailOff     = 36
	jobProposalParametersOff = 44
	jobProposalLabelsOff     = 52
	jobProposalNotBeforeOff  = 60
	jobProposalExpiresAtOff  = 68
	jobProposalSize          = 76
)

// JobProposal is a zero-copy view into a ZAP-encoded JobProposal message.
type JobProposal struct{ o zap.Object }

// WrapJobProposal parses b and returns a typed view.
func WrapJobProposal(b []byte) (JobProposal, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return JobProposal{}, err
	}
	return JobProposal{o: m.Root()}, nil
}

func (t JobProposal) ProposalID() string { return t.o.Text(jobProposalProposalIDOff) }
func (t JobProposal) DedupeKey() string  { return t.o.Text(jobProposalDedupeKeyOff) }
func (t JobProposal) JobType() string    { return t.o.Text(jobProposalJobTypeOff) }

// Priority reads priority (proto field 4, enum JobPriority).
func (t JobProposal) Priority() uint32 { return t.o.Uint32(jobProposalPriorityOff) }
func (t JobProposal) Summary() string  { return t.o.Text(jobProposalSummaryOff) }
func (t JobProposal) Detail() string   { return t.o.Text(jobProposalDetailOff) }

// Parameters reads parameters (proto field 7, map<string,ConfigValue>); each
// element is a ConfigValueMapEntry object.
func (t JobProposal) Parameters() zap.List { return t.o.List(jobProposalParametersOff) }

// Labels reads labels (proto field 8, map<string,string>); each element is a
// StringMapEntry object.
func (t JobProposal) Labels() zap.List { return t.o.List(jobProposalLabelsOff) }

// NotBefore reads not_before (proto field 9, google.protobuf.Timestamp) as
// int64 unix-nanoseconds.
func (t JobProposal) NotBefore() int64 { return t.o.Int64(jobProposalNotBeforeOff) }

// ExpiresAt reads expires_at (proto field 10, google.protobuf.Timestamp) as
// int64 unix-nanoseconds.
func (t JobProposal) ExpiresAt() int64 { return t.o.Int64(jobProposalExpiresAtOff) }

// JobProposalInput collects the field values for NewJobProposal.
type JobProposalInput struct {
	ProposalID string
	DedupeKey  string
	JobType    string
	Priority   uint32
	Summary    string
	Detail     string
	Parameters [][]byte
	Labels     [][]byte
	NotBefore  int64
	ExpiresAt  int64
}

// NewJobProposal builds a ZAP-encoded JobProposal message.
func NewJobProposal(in JobProposalInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(jobProposalSize)
	ob.SetText(jobProposalProposalIDOff, in.ProposalID)
	ob.SetText(jobProposalDedupeKeyOff, in.DedupeKey)
	ob.SetText(jobProposalJobTypeOff, in.JobType)
	ob.SetUint32(jobProposalPriorityOff, in.Priority)
	ob.SetText(jobProposalSummaryOff, in.Summary)
	ob.SetText(jobProposalDetailOff, in.Detail)
	pb := b.StartList(0)
	for _, elem := range in.Parameters {
		pb.AddObjectBytes(elem)
	}
	ob.SetList(jobProposalParametersOff, pb.FinishOffset(), len(in.Parameters))
	lb := b.StartList(0)
	for _, elem := range in.Labels {
		lb.AddObjectBytes(elem)
	}
	ob.SetList(jobProposalLabelsOff, lb.FinishOffset(), len(in.Labels))
	ob.SetInt64(jobProposalNotBeforeOff, in.NotBefore)
	ob.SetInt64(jobProposalExpiresAtOff, in.ExpiresAt)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ExecuteJobRequest ---

const (
	executeJobRequestRequestIDOff          = 0
	executeJobRequestJobOff                = 8
	executeJobRequestAdminRuntimeOff       = 16
	executeJobRequestAdminConfigValuesOff  = 24
	executeJobRequestWorkerConfigValuesOff = 32
	executeJobRequestClusterContextOff     = 40
	executeJobRequestAttemptOff            = 48
	executeJobRequestSize                  = 52
)

// ExecuteJobRequest is a zero-copy view into a ZAP-encoded ExecuteJobRequest.
type ExecuteJobRequest struct{ o zap.Object }

// WrapExecuteJobRequest parses b and returns a typed view.
func WrapExecuteJobRequest(b []byte) (ExecuteJobRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ExecuteJobRequest{}, err
	}
	return ExecuteJobRequest{o: m.Root()}, nil
}

func (t ExecuteJobRequest) RequestID() string { return t.o.Text(executeJobRequestRequestIDOff) }
func (t ExecuteJobRequest) Job() JobSpec {
	return JobSpec{o: childObject(t.o.Bytes(executeJobRequestJobOff))}
}
func (t ExecuteJobRequest) AdminRuntime() AdminRuntimeConfig {
	return AdminRuntimeConfig{o: childObject(t.o.Bytes(executeJobRequestAdminRuntimeOff))}
}

// AdminConfigValues reads admin_config_values (proto field 4,
// map<string,ConfigValue>); each element is a ConfigValueMapEntry object.
func (t ExecuteJobRequest) AdminConfigValues() zap.List {
	return t.o.List(executeJobRequestAdminConfigValuesOff)
}

// WorkerConfigValues reads worker_config_values (proto field 5,
// map<string,ConfigValue>); each element is a ConfigValueMapEntry object.
func (t ExecuteJobRequest) WorkerConfigValues() zap.List {
	return t.o.List(executeJobRequestWorkerConfigValuesOff)
}
func (t ExecuteJobRequest) ClusterContext() ClusterContext {
	return ClusterContext{o: childObject(t.o.Bytes(executeJobRequestClusterContextOff))}
}
func (t ExecuteJobRequest) Attempt() int32 { return t.o.Int32(executeJobRequestAttemptOff) }

// ExecuteJobRequestInput collects the field values for NewExecuteJobRequest.
type ExecuteJobRequestInput struct {
	RequestID          string
	Job                []byte
	AdminRuntime       []byte
	AdminConfigValues  [][]byte
	WorkerConfigValues [][]byte
	ClusterContext     []byte
	Attempt            int32
}

// NewExecuteJobRequest builds a ZAP-encoded ExecuteJobRequest message.
func NewExecuteJobRequest(in ExecuteJobRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(executeJobRequestSize)
	ob.SetText(executeJobRequestRequestIDOff, in.RequestID)
	setChild(ob, executeJobRequestJobOff, in.Job)
	setChild(ob, executeJobRequestAdminRuntimeOff, in.AdminRuntime)
	ab := b.StartList(0)
	for _, elem := range in.AdminConfigValues {
		ab.AddObjectBytes(elem)
	}
	ob.SetList(executeJobRequestAdminConfigValuesOff, ab.FinishOffset(), len(in.AdminConfigValues))
	wb := b.StartList(0)
	for _, elem := range in.WorkerConfigValues {
		wb.AddObjectBytes(elem)
	}
	ob.SetList(executeJobRequestWorkerConfigValuesOff, wb.FinishOffset(), len(in.WorkerConfigValues))
	setChild(ob, executeJobRequestClusterContextOff, in.ClusterContext)
	ob.SetInt32(executeJobRequestAttemptOff, in.Attempt)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- JobSpec ---

const (
	jobSpecJobIDOff       = 0
	jobSpecJobTypeOff     = 8
	jobSpecDedupeKeyOff   = 16
	jobSpecPriorityOff    = 24
	jobSpecSummaryOff     = 28
	jobSpecDetailOff      = 36
	jobSpecParametersOff  = 44
	jobSpecLabelsOff      = 52
	jobSpecCreatedAtOff   = 60
	jobSpecScheduledAtOff = 68
	jobSpecSize           = 76
)

// JobSpec is a zero-copy view into a ZAP-encoded JobSpec message.
type JobSpec struct{ o zap.Object }

// WrapJobSpec parses b and returns a typed view.
func WrapJobSpec(b []byte) (JobSpec, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return JobSpec{}, err
	}
	return JobSpec{o: m.Root()}, nil
}

func (t JobSpec) JobID() string     { return t.o.Text(jobSpecJobIDOff) }
func (t JobSpec) JobType() string   { return t.o.Text(jobSpecJobTypeOff) }
func (t JobSpec) DedupeKey() string { return t.o.Text(jobSpecDedupeKeyOff) }

// Priority reads priority (proto field 4, enum JobPriority).
func (t JobSpec) Priority() uint32 { return t.o.Uint32(jobSpecPriorityOff) }
func (t JobSpec) Summary() string  { return t.o.Text(jobSpecSummaryOff) }
func (t JobSpec) Detail() string   { return t.o.Text(jobSpecDetailOff) }

// Parameters reads parameters (proto field 7, map<string,ConfigValue>); each
// element is a ConfigValueMapEntry object.
func (t JobSpec) Parameters() zap.List { return t.o.List(jobSpecParametersOff) }

// Labels reads labels (proto field 8, map<string,string>); each element is a
// StringMapEntry object.
func (t JobSpec) Labels() zap.List { return t.o.List(jobSpecLabelsOff) }

// CreatedAt reads created_at (proto field 9, google.protobuf.Timestamp) as
// int64 unix-nanoseconds.
func (t JobSpec) CreatedAt() int64 { return t.o.Int64(jobSpecCreatedAtOff) }

// ScheduledAt reads scheduled_at (proto field 10, google.protobuf.Timestamp)
// as int64 unix-nanoseconds.
func (t JobSpec) ScheduledAt() int64 { return t.o.Int64(jobSpecScheduledAtOff) }

// JobSpecInput collects the field values for NewJobSpec.
type JobSpecInput struct {
	JobID       string
	JobType     string
	DedupeKey   string
	Priority    uint32
	Summary     string
	Detail      string
	Parameters  [][]byte
	Labels      [][]byte
	CreatedAt   int64
	ScheduledAt int64
}

// NewJobSpec builds a ZAP-encoded JobSpec message.
func NewJobSpec(in JobSpecInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(jobSpecSize)
	ob.SetText(jobSpecJobIDOff, in.JobID)
	ob.SetText(jobSpecJobTypeOff, in.JobType)
	ob.SetText(jobSpecDedupeKeyOff, in.DedupeKey)
	ob.SetUint32(jobSpecPriorityOff, in.Priority)
	ob.SetText(jobSpecSummaryOff, in.Summary)
	ob.SetText(jobSpecDetailOff, in.Detail)
	pb := b.StartList(0)
	for _, elem := range in.Parameters {
		pb.AddObjectBytes(elem)
	}
	ob.SetList(jobSpecParametersOff, pb.FinishOffset(), len(in.Parameters))
	lb := b.StartList(0)
	for _, elem := range in.Labels {
		lb.AddObjectBytes(elem)
	}
	ob.SetList(jobSpecLabelsOff, lb.FinishOffset(), len(in.Labels))
	ob.SetInt64(jobSpecCreatedAtOff, in.CreatedAt)
	ob.SetInt64(jobSpecScheduledAtOff, in.ScheduledAt)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- JobProgressUpdate ---

const (
	jobProgressUpdateRequestIDOff       = 0
	jobProgressUpdateJobIDOff           = 8
	jobProgressUpdateJobTypeOff         = 16
	jobProgressUpdateStateOff           = 24
	jobProgressUpdateProgressPercentOff = 28
	jobProgressUpdateStageOff           = 36
	jobProgressUpdateMessageOff         = 44
	jobProgressUpdateMetricsOff         = 52
	jobProgressUpdateActivitiesOff      = 60
	jobProgressUpdateUpdatedAtOff       = 68
	jobProgressUpdateSize               = 76
)

// JobProgressUpdate is a zero-copy view into a ZAP-encoded JobProgressUpdate.
type JobProgressUpdate struct{ o zap.Object }

// WrapJobProgressUpdate parses b and returns a typed view.
func WrapJobProgressUpdate(b []byte) (JobProgressUpdate, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return JobProgressUpdate{}, err
	}
	return JobProgressUpdate{o: m.Root()}, nil
}

func (t JobProgressUpdate) RequestID() string { return t.o.Text(jobProgressUpdateRequestIDOff) }
func (t JobProgressUpdate) JobID() string     { return t.o.Text(jobProgressUpdateJobIDOff) }
func (t JobProgressUpdate) JobType() string   { return t.o.Text(jobProgressUpdateJobTypeOff) }

// State reads state (proto field 4, enum JobState).
func (t JobProgressUpdate) State() uint32 { return t.o.Uint32(jobProgressUpdateStateOff) }
func (t JobProgressUpdate) ProgressPercent() float64 {
	return t.o.Float64(jobProgressUpdateProgressPercentOff)
}
func (t JobProgressUpdate) Stage() string   { return t.o.Text(jobProgressUpdateStageOff) }
func (t JobProgressUpdate) Message() string { return t.o.Text(jobProgressUpdateMessageOff) }

// Metrics reads metrics (proto field 8, map<string,ConfigValue>); each element
// is a ConfigValueMapEntry object.
func (t JobProgressUpdate) Metrics() zap.List { return t.o.List(jobProgressUpdateMetricsOff) }

// Activities reads activities (proto field 9, repeated ActivityEvent); each
// element is an ActivityEvent object.
func (t JobProgressUpdate) Activities() zap.List { return t.o.List(jobProgressUpdateActivitiesOff) }

// UpdatedAt reads updated_at (proto field 10, google.protobuf.Timestamp) as
// int64 unix-nanoseconds.
func (t JobProgressUpdate) UpdatedAt() int64 { return t.o.Int64(jobProgressUpdateUpdatedAtOff) }

// JobProgressUpdateInput collects the field values for NewJobProgressUpdate.
type JobProgressUpdateInput struct {
	RequestID       string
	JobID           string
	JobType         string
	State           uint32
	ProgressPercent float64
	Stage           string
	Message         string
	Metrics         [][]byte
	Activities      [][]byte
	UpdatedAt       int64
}

// NewJobProgressUpdate builds a ZAP-encoded JobProgressUpdate message.
func NewJobProgressUpdate(in JobProgressUpdateInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(jobProgressUpdateSize)
	ob.SetText(jobProgressUpdateRequestIDOff, in.RequestID)
	ob.SetText(jobProgressUpdateJobIDOff, in.JobID)
	ob.SetText(jobProgressUpdateJobTypeOff, in.JobType)
	ob.SetUint32(jobProgressUpdateStateOff, in.State)
	ob.SetFloat64(jobProgressUpdateProgressPercentOff, in.ProgressPercent)
	ob.SetText(jobProgressUpdateStageOff, in.Stage)
	ob.SetText(jobProgressUpdateMessageOff, in.Message)
	mb := b.StartList(0)
	for _, elem := range in.Metrics {
		mb.AddObjectBytes(elem)
	}
	ob.SetList(jobProgressUpdateMetricsOff, mb.FinishOffset(), len(in.Metrics))
	acb := b.StartList(0)
	for _, elem := range in.Activities {
		acb.AddObjectBytes(elem)
	}
	ob.SetList(jobProgressUpdateActivitiesOff, acb.FinishOffset(), len(in.Activities))
	ob.SetInt64(jobProgressUpdateUpdatedAtOff, in.UpdatedAt)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- JobCompleted ---

const (
	jobCompletedRequestIDOff    = 0
	jobCompletedJobIDOff        = 8
	jobCompletedJobTypeOff      = 16
	jobCompletedSuccessOff      = 24
	jobCompletedErrorMessageOff = 28
	jobCompletedResultOff       = 36
	jobCompletedActivitiesOff   = 44
	jobCompletedCompletedAtOff  = 52
	jobCompletedSize            = 60
)

// JobCompleted is a zero-copy view into a ZAP-encoded JobCompleted message.
type JobCompleted struct{ o zap.Object }

// WrapJobCompleted parses b and returns a typed view.
func WrapJobCompleted(b []byte) (JobCompleted, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return JobCompleted{}, err
	}
	return JobCompleted{o: m.Root()}, nil
}

func (t JobCompleted) RequestID() string    { return t.o.Text(jobCompletedRequestIDOff) }
func (t JobCompleted) JobID() string        { return t.o.Text(jobCompletedJobIDOff) }
func (t JobCompleted) JobType() string      { return t.o.Text(jobCompletedJobTypeOff) }
func (t JobCompleted) Success() bool        { return t.o.Bool(jobCompletedSuccessOff) }
func (t JobCompleted) ErrorMessage() string { return t.o.Text(jobCompletedErrorMessageOff) }
func (t JobCompleted) Result() JobResult {
	return JobResult{o: childObject(t.o.Bytes(jobCompletedResultOff))}
}

// Activities reads activities (proto field 7, repeated ActivityEvent); each
// element is an ActivityEvent object.
func (t JobCompleted) Activities() zap.List { return t.o.List(jobCompletedActivitiesOff) }

// CompletedAt reads completed_at (proto field 8, google.protobuf.Timestamp) as
// int64 unix-nanoseconds.
func (t JobCompleted) CompletedAt() int64 { return t.o.Int64(jobCompletedCompletedAtOff) }

// JobCompletedInput collects the field values for NewJobCompleted.
type JobCompletedInput struct {
	RequestID    string
	JobID        string
	JobType      string
	Success      bool
	ErrorMessage string
	Result       []byte
	Activities   [][]byte
	CompletedAt  int64
}

// NewJobCompleted builds a ZAP-encoded JobCompleted message.
func NewJobCompleted(in JobCompletedInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(jobCompletedSize)
	ob.SetText(jobCompletedRequestIDOff, in.RequestID)
	ob.SetText(jobCompletedJobIDOff, in.JobID)
	ob.SetText(jobCompletedJobTypeOff, in.JobType)
	ob.SetBool(jobCompletedSuccessOff, in.Success)
	ob.SetText(jobCompletedErrorMessageOff, in.ErrorMessage)
	setChild(ob, jobCompletedResultOff, in.Result)
	acb := b.StartList(0)
	for _, elem := range in.Activities {
		acb.AddObjectBytes(elem)
	}
	ob.SetList(jobCompletedActivitiesOff, acb.FinishOffset(), len(in.Activities))
	ob.SetInt64(jobCompletedCompletedAtOff, in.CompletedAt)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- JobResult ---

const (
	jobResultOutputValuesOff = 0
	jobResultSummaryOff      = 8
	jobResultSize            = 16
)

// JobResult is a zero-copy view into a ZAP-encoded JobResult message.
type JobResult struct{ o zap.Object }

// WrapJobResult parses b and returns a typed view.
func WrapJobResult(b []byte) (JobResult, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return JobResult{}, err
	}
	return JobResult{o: m.Root()}, nil
}

// OutputValues reads output_values (proto field 1, map<string,ConfigValue>);
// each element is a ConfigValueMapEntry object.
func (t JobResult) OutputValues() zap.List { return t.o.List(jobResultOutputValuesOff) }
func (t JobResult) Summary() string        { return t.o.Text(jobResultSummaryOff) }

// JobResultInput collects the field values for NewJobResult.
type JobResultInput struct {
	OutputValues [][]byte
	Summary      string
}

// NewJobResult builds a ZAP-encoded JobResult message.
func NewJobResult(in JobResultInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(jobResultSize)
	lb := b.StartList(0)
	for _, elem := range in.OutputValues {
		lb.AddObjectBytes(elem)
	}
	ob.SetList(jobResultOutputValuesOff, lb.FinishOffset(), len(in.OutputValues))
	ob.SetText(jobResultSummaryOff, in.Summary)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ClusterContext ---

const (
	clusterContextMasterGRPCAddressesOff = 0
	clusterContextFilerAddressesOff      = 8
	clusterContextVolumeGRPCAddressesOff = 16
	clusterContextMetadataOff            = 24
	clusterContextS3GRPCAddressesOff     = 32
	clusterContextSize                   = 40
)

// ClusterContext is a zero-copy view into a ZAP-encoded ClusterContext message.
type ClusterContext struct{ o zap.Object }

// WrapClusterContext parses b and returns a typed view.
func WrapClusterContext(b []byte) (ClusterContext, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ClusterContext{}, err
	}
	return ClusterContext{o: m.Root()}, nil
}

// MasterGRPCAddresses reads master_grpc_addresses (proto field 1, repeated
// string); each element is a raw-bytes entry (use List.BytesAt(i)).
func (t ClusterContext) MasterGRPCAddresses() zap.List {
	return t.o.List(clusterContextMasterGRPCAddressesOff)
}

// FilerAddresses reads filer_addresses (proto field 2, repeated string).
func (t ClusterContext) FilerAddresses() zap.List { return t.o.List(clusterContextFilerAddressesOff) }

// VolumeGRPCAddresses reads volume_grpc_addresses (proto field 3, repeated string).
func (t ClusterContext) VolumeGRPCAddresses() zap.List {
	return t.o.List(clusterContextVolumeGRPCAddressesOff)
}

// Metadata reads metadata (proto field 4, map<string,string>); each element is
// a StringMapEntry object.
func (t ClusterContext) Metadata() zap.List { return t.o.List(clusterContextMetadataOff) }

// S3GRPCAddresses reads s3_grpc_addresses (proto field 5, repeated string).
func (t ClusterContext) S3GRPCAddresses() zap.List { return t.o.List(clusterContextS3GRPCAddressesOff) }

// ClusterContextInput collects the field values for NewClusterContext. Each
// address-list element is the raw UTF-8 bytes of one address.
type ClusterContextInput struct {
	MasterGRPCAddresses [][]byte
	FilerAddresses      [][]byte
	VolumeGRPCAddresses [][]byte
	Metadata            [][]byte
	S3GRPCAddresses     [][]byte
}

// NewClusterContext builds a ZAP-encoded ClusterContext message.
func NewClusterContext(in ClusterContextInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(clusterContextSize)
	mab := b.StartList(0)
	for _, elem := range in.MasterGRPCAddresses {
		mab.AddObjectBytes(elem)
	}
	ob.SetList(clusterContextMasterGRPCAddressesOff, mab.FinishOffset(), len(in.MasterGRPCAddresses))
	fab := b.StartList(0)
	for _, elem := range in.FilerAddresses {
		fab.AddObjectBytes(elem)
	}
	ob.SetList(clusterContextFilerAddressesOff, fab.FinishOffset(), len(in.FilerAddresses))
	vab := b.StartList(0)
	for _, elem := range in.VolumeGRPCAddresses {
		vab.AddObjectBytes(elem)
	}
	ob.SetList(clusterContextVolumeGRPCAddressesOff, vab.FinishOffset(), len(in.VolumeGRPCAddresses))
	meb := b.StartList(0)
	for _, elem := range in.Metadata {
		meb.AddObjectBytes(elem)
	}
	ob.SetList(clusterContextMetadataOff, meb.FinishOffset(), len(in.Metadata))
	sab := b.StartList(0)
	for _, elem := range in.S3GRPCAddresses {
		sab.AddObjectBytes(elem)
	}
	ob.SetList(clusterContextS3GRPCAddressesOff, sab.FinishOffset(), len(in.S3GRPCAddresses))
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ActivityEvent ---

const (
	activityEventSourceOff    = 0
	activityEventMessageOff   = 4
	activityEventStageOff     = 12
	activityEventDetailsOff   = 20
	activityEventCreatedAtOff = 28
	activityEventSize         = 36
)

// ActivityEvent is a zero-copy view into a ZAP-encoded ActivityEvent message.
type ActivityEvent struct{ o zap.Object }

// WrapActivityEvent parses b and returns a typed view.
func WrapActivityEvent(b []byte) (ActivityEvent, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ActivityEvent{}, err
	}
	return ActivityEvent{o: m.Root()}, nil
}

// Source reads source (proto field 1, enum ActivitySource).
func (t ActivityEvent) Source() uint32  { return t.o.Uint32(activityEventSourceOff) }
func (t ActivityEvent) Message() string { return t.o.Text(activityEventMessageOff) }
func (t ActivityEvent) Stage() string   { return t.o.Text(activityEventStageOff) }

// Details reads details (proto field 4, map<string,ConfigValue>); each element
// is a ConfigValueMapEntry object.
func (t ActivityEvent) Details() zap.List { return t.o.List(activityEventDetailsOff) }

// CreatedAt reads created_at (proto field 5, google.protobuf.Timestamp) as
// int64 unix-nanoseconds.
func (t ActivityEvent) CreatedAt() int64 { return t.o.Int64(activityEventCreatedAtOff) }

// ActivityEventInput collects the field values for NewActivityEvent.
type ActivityEventInput struct {
	Source    uint32
	Message   string
	Stage     string
	Details   [][]byte
	CreatedAt int64
}

// NewActivityEvent builds a ZAP-encoded ActivityEvent message.
func NewActivityEvent(in ActivityEventInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(activityEventSize)
	ob.SetUint32(activityEventSourceOff, in.Source)
	ob.SetText(activityEventMessageOff, in.Message)
	ob.SetText(activityEventStageOff, in.Stage)
	db := b.StartList(0)
	for _, elem := range in.Details {
		db.AddObjectBytes(elem)
	}
	ob.SetList(activityEventDetailsOff, db.FinishOffset(), len(in.Details))
	ob.SetInt64(activityEventCreatedAtOff, in.CreatedAt)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- CancelRequest ---

const (
	cancelRequestTargetIDOff   = 0
	cancelRequestTargetKindOff = 8
	cancelRequestReasonOff     = 12
	cancelRequestForceOff      = 20
	cancelRequestSize          = 21
)

// CancelRequest is a zero-copy view into a ZAP-encoded CancelRequest message.
type CancelRequest struct{ o zap.Object }

// WrapCancelRequest parses b and returns a typed view.
func WrapCancelRequest(b []byte) (CancelRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return CancelRequest{}, err
	}
	return CancelRequest{o: m.Root()}, nil
}

func (t CancelRequest) TargetID() string { return t.o.Text(cancelRequestTargetIDOff) }

// TargetKind reads target_kind (proto field 2, enum WorkKind).
func (t CancelRequest) TargetKind() uint32 { return t.o.Uint32(cancelRequestTargetKindOff) }
func (t CancelRequest) Reason() string     { return t.o.Text(cancelRequestReasonOff) }
func (t CancelRequest) Force() bool        { return t.o.Bool(cancelRequestForceOff) }

// CancelRequestInput collects the field values for NewCancelRequest.
type CancelRequestInput struct {
	TargetID   string
	TargetKind uint32
	Reason     string
	Force      bool
}

// NewCancelRequest builds a ZAP-encoded CancelRequest message.
func NewCancelRequest(in CancelRequestInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(cancelRequestSize)
	ob.SetText(cancelRequestTargetIDOff, in.TargetID)
	ob.SetUint32(cancelRequestTargetKindOff, in.TargetKind)
	ob.SetText(cancelRequestReasonOff, in.Reason)
	ob.SetBool(cancelRequestForceOff, in.Force)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- AdminShutdown ---

const (
	adminShutdownReasonOff             = 0
	adminShutdownGracePeriodSecondsOff = 8
	adminShutdownSize                  = 12
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
func (t AdminShutdown) GracePeriodSeconds() int32 {
	return t.o.Int32(adminShutdownGracePeriodSecondsOff)
}

// AdminShutdownInput collects the field values for NewAdminShutdown.
type AdminShutdownInput struct {
	Reason             string
	GracePeriodSeconds int32
}

// NewAdminShutdown builds a ZAP-encoded AdminShutdown message.
func NewAdminShutdown(in AdminShutdownInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(adminShutdownSize)
	ob.SetText(adminShutdownReasonOff, in.Reason)
	ob.SetInt32(adminShutdownGracePeriodSecondsOff, in.GracePeriodSeconds)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- PersistedJobTypeConfig ---

const (
	persistedJobTypeConfigJobTypeOff            = 0
	persistedJobTypeConfigDescriptorVersionOff  = 8
	persistedJobTypeConfigAdminConfigValuesOff  = 12
	persistedJobTypeConfigWorkerConfigValuesOff = 20
	persistedJobTypeConfigAdminRuntimeOff       = 28
	persistedJobTypeConfigUpdatedAtOff          = 36
	persistedJobTypeConfigUpdatedByOff          = 44
	persistedJobTypeConfigSize                  = 52
)

// PersistedJobTypeConfig is a zero-copy view into a ZAP-encoded
// PersistedJobTypeConfig message.
type PersistedJobTypeConfig struct{ o zap.Object }

// WrapPersistedJobTypeConfig parses b and returns a typed view.
func WrapPersistedJobTypeConfig(b []byte) (PersistedJobTypeConfig, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PersistedJobTypeConfig{}, err
	}
	return PersistedJobTypeConfig{o: m.Root()}, nil
}

func (t PersistedJobTypeConfig) JobType() string { return t.o.Text(persistedJobTypeConfigJobTypeOff) }
func (t PersistedJobTypeConfig) DescriptorVersion() uint32 {
	return t.o.Uint32(persistedJobTypeConfigDescriptorVersionOff)
}

// AdminConfigValues reads admin_config_values (proto field 3,
// map<string,ConfigValue>); each element is a ConfigValueMapEntry object.
func (t PersistedJobTypeConfig) AdminConfigValues() zap.List {
	return t.o.List(persistedJobTypeConfigAdminConfigValuesOff)
}

// WorkerConfigValues reads worker_config_values (proto field 4,
// map<string,ConfigValue>); each element is a ConfigValueMapEntry object.
func (t PersistedJobTypeConfig) WorkerConfigValues() zap.List {
	return t.o.List(persistedJobTypeConfigWorkerConfigValuesOff)
}
func (t PersistedJobTypeConfig) AdminRuntime() AdminRuntimeConfig {
	return AdminRuntimeConfig{o: childObject(t.o.Bytes(persistedJobTypeConfigAdminRuntimeOff))}
}

// UpdatedAt reads updated_at (proto field 6, google.protobuf.Timestamp) as
// int64 unix-nanoseconds.
func (t PersistedJobTypeConfig) UpdatedAt() int64 {
	return t.o.Int64(persistedJobTypeConfigUpdatedAtOff)
}
func (t PersistedJobTypeConfig) UpdatedBy() string {
	return t.o.Text(persistedJobTypeConfigUpdatedByOff)
}

// PersistedJobTypeConfigInput collects the field values for
// NewPersistedJobTypeConfig.
type PersistedJobTypeConfigInput struct {
	JobType            string
	DescriptorVersion  uint32
	AdminConfigValues  [][]byte
	WorkerConfigValues [][]byte
	AdminRuntime       []byte
	UpdatedAt          int64
	UpdatedBy          string
}

// NewPersistedJobTypeConfig builds a ZAP-encoded PersistedJobTypeConfig message.
func NewPersistedJobTypeConfig(in PersistedJobTypeConfigInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(persistedJobTypeConfigSize)
	ob.SetText(persistedJobTypeConfigJobTypeOff, in.JobType)
	ob.SetUint32(persistedJobTypeConfigDescriptorVersionOff, in.DescriptorVersion)
	ab := b.StartList(0)
	for _, elem := range in.AdminConfigValues {
		ab.AddObjectBytes(elem)
	}
	ob.SetList(persistedJobTypeConfigAdminConfigValuesOff, ab.FinishOffset(), len(in.AdminConfigValues))
	wb := b.StartList(0)
	for _, elem := range in.WorkerConfigValues {
		wb.AddObjectBytes(elem)
	}
	ob.SetList(persistedJobTypeConfigWorkerConfigValuesOff, wb.FinishOffset(), len(in.WorkerConfigValues))
	setChild(ob, persistedJobTypeConfigAdminRuntimeOff, in.AdminRuntime)
	ob.SetInt64(persistedJobTypeConfigUpdatedAtOff, in.UpdatedAt)
	ob.SetText(persistedJobTypeConfigUpdatedByOff, in.UpdatedBy)
	ob.FinishAsRoot()
	return b.Finish()
}
