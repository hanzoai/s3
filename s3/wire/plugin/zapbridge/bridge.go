// Package pluginzapbridge converts the two PluginControlService.WorkerStream
// frame types between their protobuf domain form (s3/pb/plugin_pb) and their
// ZAP wire form (s3/wire/plugin). The domain keeps speaking *plugin_pb.*; only
// the stream's Send/Recv edges cross through here, so the bytes on the wire are
// pure ZAP with no gRPC or protobuf framing.
//
// The conversion mirrors plugin.proto exactly, recursing the whole oneof tree.
// google.protobuf.Timestamp and google.protobuf.Duration map to int64
// unix-/duration-nanoseconds, matching the pluginwire encoding. Maps are
// emitted in sorted-key order for a deterministic wire image.
package pluginzapbridge

import (
	"math"
	"sort"
	"time"

	"github.com/hanzoai/s3/s3/pb/plugin_pb"
	pluginwire "github.com/hanzoai/s3/s3/wire/plugin"
	zap "github.com/zap-proto/go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// orNil collapses a decoded singular nested message to nil when it carries no
// set fields. pluginwire reads an absent nested message as an all-zero view (a
// null child pointer falls back to an empty frame), so a zero-valued decode is
// indistinguishable from "field absent" — and proto3 models an absent singular
// message as nil. Collapsing here keeps the round-trip faithful.
func orNil[T proto.Message](m T) T {
	if proto.Size(m) == 0 {
		var zero T
		return zero
	}
	return m
}

// ===================== top-level WorkerStream frames =====================

// MarshalWorkerToAdmin encodes a *plugin_pb.WorkerToAdminMessage to ZAP bytes.
func MarshalWorkerToAdmin(m *plugin_pb.WorkerToAdminMessage) []byte {
	if m == nil {
		return pluginwire.NewWorkerToAdminMessage(pluginwire.WorkerToAdminMessageInput{})
	}
	in := pluginwire.WorkerToAdminMessageInput{
		WorkerID: m.GetWorkerId(),
		SentAt:   tsNanos(m.GetSentAt()),
	}
	switch m.Body.(type) {
	case *plugin_pb.WorkerToAdminMessage_Hello:
		in.Which = pluginwire.WorkerToAdminBodyHello
		in.Hello = marshalWorkerHello(m.GetHello())
	case *plugin_pb.WorkerToAdminMessage_Heartbeat:
		in.Which = pluginwire.WorkerToAdminBodyHeartbeat
		in.Heartbeat = marshalWorkerHeartbeat(m.GetHeartbeat())
	case *plugin_pb.WorkerToAdminMessage_Acknowledge:
		in.Which = pluginwire.WorkerToAdminBodyAcknowledge
		in.Acknowledge = marshalWorkerAcknowledge(m.GetAcknowledge())
	case *plugin_pb.WorkerToAdminMessage_ConfigSchemaResponse:
		in.Which = pluginwire.WorkerToAdminBodyConfigSchemaResp
		in.ConfigSchemaResponse = marshalConfigSchemaResponse(m.GetConfigSchemaResponse())
	case *plugin_pb.WorkerToAdminMessage_DetectionProposals:
		in.Which = pluginwire.WorkerToAdminBodyDetectionProposals
		in.DetectionProposals = marshalDetectionProposals(m.GetDetectionProposals())
	case *plugin_pb.WorkerToAdminMessage_DetectionComplete:
		in.Which = pluginwire.WorkerToAdminBodyDetectionComplete
		in.DetectionComplete = marshalDetectionComplete(m.GetDetectionComplete())
	case *plugin_pb.WorkerToAdminMessage_JobProgressUpdate:
		in.Which = pluginwire.WorkerToAdminBodyJobProgressUpdate
		in.JobProgressUpdate = marshalJobProgressUpdate(m.GetJobProgressUpdate())
	case *plugin_pb.WorkerToAdminMessage_JobCompleted:
		in.Which = pluginwire.WorkerToAdminBodyJobCompleted
		in.JobCompleted = marshalJobCompleted(m.GetJobCompleted())
	}
	return pluginwire.NewWorkerToAdminMessage(in)
}

// UnmarshalWorkerToAdmin decodes ZAP bytes to a *plugin_pb.WorkerToAdminMessage.
func UnmarshalWorkerToAdmin(b []byte) (*plugin_pb.WorkerToAdminMessage, error) {
	t, err := pluginwire.WrapWorkerToAdminMessage(b)
	if err != nil {
		return nil, err
	}
	m := &plugin_pb.WorkerToAdminMessage{
		WorkerId: t.WorkerID(),
		SentAt:   nanosTS(t.SentAt()),
	}
	switch t.WhichBody() {
	case pluginwire.WorkerToAdminBodyHello:
		m.Body = &plugin_pb.WorkerToAdminMessage_Hello{Hello: unmarshalWorkerHello(t.Hello())}
	case pluginwire.WorkerToAdminBodyHeartbeat:
		m.Body = &plugin_pb.WorkerToAdminMessage_Heartbeat{Heartbeat: unmarshalWorkerHeartbeat(t.Heartbeat())}
	case pluginwire.WorkerToAdminBodyAcknowledge:
		m.Body = &plugin_pb.WorkerToAdminMessage_Acknowledge{Acknowledge: unmarshalWorkerAcknowledge(t.Acknowledge())}
	case pluginwire.WorkerToAdminBodyConfigSchemaResp:
		m.Body = &plugin_pb.WorkerToAdminMessage_ConfigSchemaResponse{ConfigSchemaResponse: unmarshalConfigSchemaResponse(t.ConfigSchemaResponse())}
	case pluginwire.WorkerToAdminBodyDetectionProposals:
		m.Body = &plugin_pb.WorkerToAdminMessage_DetectionProposals{DetectionProposals: unmarshalDetectionProposals(t.DetectionProposals())}
	case pluginwire.WorkerToAdminBodyDetectionComplete:
		m.Body = &plugin_pb.WorkerToAdminMessage_DetectionComplete{DetectionComplete: unmarshalDetectionComplete(t.DetectionComplete())}
	case pluginwire.WorkerToAdminBodyJobProgressUpdate:
		m.Body = &plugin_pb.WorkerToAdminMessage_JobProgressUpdate{JobProgressUpdate: unmarshalJobProgressUpdate(t.JobProgressUpdate())}
	case pluginwire.WorkerToAdminBodyJobCompleted:
		m.Body = &plugin_pb.WorkerToAdminMessage_JobCompleted{JobCompleted: unmarshalJobCompleted(t.JobCompleted())}
	}
	return m, nil
}

// MarshalAdminToWorker encodes a *plugin_pb.AdminToWorkerMessage to ZAP bytes.
func MarshalAdminToWorker(m *plugin_pb.AdminToWorkerMessage) []byte {
	if m == nil {
		return pluginwire.NewAdminToWorkerMessage(pluginwire.AdminToWorkerMessageInput{})
	}
	in := pluginwire.AdminToWorkerMessageInput{
		RequestID: m.GetRequestId(),
		SentAt:    tsNanos(m.GetSentAt()),
	}
	switch m.Body.(type) {
	case *plugin_pb.AdminToWorkerMessage_Hello:
		in.Which = pluginwire.AdminToWorkerBodyHello
		in.Hello = marshalAdminHello(m.GetHello())
	case *plugin_pb.AdminToWorkerMessage_RequestConfigSchema:
		in.Which = pluginwire.AdminToWorkerBodyRequestConfigSchema
		in.RequestConfigSchema = marshalRequestConfigSchema(m.GetRequestConfigSchema())
	case *plugin_pb.AdminToWorkerMessage_RunDetectionRequest:
		in.Which = pluginwire.AdminToWorkerBodyRunDetectionRequest
		in.RunDetectionRequest = marshalRunDetectionRequest(m.GetRunDetectionRequest())
	case *plugin_pb.AdminToWorkerMessage_ExecuteJobRequest:
		in.Which = pluginwire.AdminToWorkerBodyExecuteJobRequest
		in.ExecuteJobRequest = marshalExecuteJobRequest(m.GetExecuteJobRequest())
	case *plugin_pb.AdminToWorkerMessage_CancelRequest:
		in.Which = pluginwire.AdminToWorkerBodyCancelRequest
		in.CancelRequest = marshalCancelRequest(m.GetCancelRequest())
	case *plugin_pb.AdminToWorkerMessage_Shutdown:
		in.Which = pluginwire.AdminToWorkerBodyShutdown
		in.Shutdown = marshalAdminShutdown(m.GetShutdown())
	}
	return pluginwire.NewAdminToWorkerMessage(in)
}

// UnmarshalAdminToWorker decodes ZAP bytes to a *plugin_pb.AdminToWorkerMessage.
func UnmarshalAdminToWorker(b []byte) (*plugin_pb.AdminToWorkerMessage, error) {
	t, err := pluginwire.WrapAdminToWorkerMessage(b)
	if err != nil {
		return nil, err
	}
	m := &plugin_pb.AdminToWorkerMessage{
		RequestId: t.RequestID(),
		SentAt:    nanosTS(t.SentAt()),
	}
	switch t.WhichBody() {
	case pluginwire.AdminToWorkerBodyHello:
		m.Body = &plugin_pb.AdminToWorkerMessage_Hello{Hello: unmarshalAdminHello(t.Hello())}
	case pluginwire.AdminToWorkerBodyRequestConfigSchema:
		m.Body = &plugin_pb.AdminToWorkerMessage_RequestConfigSchema{RequestConfigSchema: unmarshalRequestConfigSchema(t.RequestConfigSchema())}
	case pluginwire.AdminToWorkerBodyRunDetectionRequest:
		m.Body = &plugin_pb.AdminToWorkerMessage_RunDetectionRequest{RunDetectionRequest: unmarshalRunDetectionRequest(t.RunDetectionRequest())}
	case pluginwire.AdminToWorkerBodyExecuteJobRequest:
		m.Body = &plugin_pb.AdminToWorkerMessage_ExecuteJobRequest{ExecuteJobRequest: unmarshalExecuteJobRequest(t.ExecuteJobRequest())}
	case pluginwire.AdminToWorkerBodyCancelRequest:
		m.Body = &plugin_pb.AdminToWorkerMessage_CancelRequest{CancelRequest: unmarshalCancelRequest(t.CancelRequest())}
	case pluginwire.AdminToWorkerBodyShutdown:
		m.Body = &plugin_pb.AdminToWorkerMessage_Shutdown{Shutdown: unmarshalAdminShutdown(t.Shutdown())}
	}
	return m, nil
}

// ===================== worker lifecycle messages =====================

func marshalWorkerHello(h *plugin_pb.WorkerHello) []byte {
	if h == nil {
		return nil
	}
	caps := make([][]byte, 0, len(h.GetCapabilities()))
	for _, c := range h.GetCapabilities() {
		caps = append(caps, marshalJobTypeCapability(c))
	}
	return pluginwire.NewWorkerHello(pluginwire.WorkerHelloInput{
		WorkerID:         h.GetWorkerId(),
		WorkerInstanceID: h.GetWorkerInstanceId(),
		Address:          h.GetAddress(),
		WorkerVersion:    h.GetWorkerVersion(),
		ProtocolVersion:  h.GetProtocolVersion(),
		Capabilities:     caps,
		Metadata:         stringMapEntries(h.GetMetadata()),
	})
}

func unmarshalWorkerHello(t pluginwire.WorkerHello) *plugin_pb.WorkerHello {
	caps := t.Capabilities()
	out := make([]*plugin_pb.JobTypeCapability, caps.Len())
	for i := range out {
		out[i] = unmarshalJobTypeCapability(wrap(caps, i, pluginwire.WrapJobTypeCapability))
	}
	return &plugin_pb.WorkerHello{
		WorkerId:         t.WorkerID(),
		WorkerInstanceId: t.WorkerInstanceID(),
		Address:          t.Address(),
		WorkerVersion:    t.WorkerVersion(),
		ProtocolVersion:  t.ProtocolVersion(),
		Capabilities:     out,
		Metadata:         readStringMap(t.Metadata()),
	}
}

func marshalAdminHello(h *plugin_pb.AdminHello) []byte {
	if h == nil {
		return nil
	}
	return pluginwire.NewAdminHello(pluginwire.AdminHelloInput{
		Accepted:                 h.GetAccepted(),
		Message:                  h.GetMessage(),
		HeartbeatIntervalSeconds: h.GetHeartbeatIntervalSeconds(),
		ReconnectDelaySeconds:    h.GetReconnectDelaySeconds(),
	})
}

func unmarshalAdminHello(t pluginwire.AdminHello) *plugin_pb.AdminHello {
	return &plugin_pb.AdminHello{
		Accepted:                 t.Accepted(),
		Message:                  t.Message(),
		HeartbeatIntervalSeconds: t.HeartbeatIntervalSeconds(),
		ReconnectDelaySeconds:    t.ReconnectDelaySeconds(),
	}
}

func marshalWorkerHeartbeat(h *plugin_pb.WorkerHeartbeat) []byte {
	if h == nil {
		return nil
	}
	rw := make([][]byte, 0, len(h.GetRunningWork()))
	for _, w := range h.GetRunningWork() {
		rw = append(rw, marshalRunningWork(w))
	}
	return pluginwire.NewWorkerHeartbeat(pluginwire.WorkerHeartbeatInput{
		WorkerID:            h.GetWorkerId(),
		RunningWork:         rw,
		DetectionSlotsUsed:  h.GetDetectionSlotsUsed(),
		DetectionSlotsTotal: h.GetDetectionSlotsTotal(),
		ExecutionSlotsUsed:  h.GetExecutionSlotsUsed(),
		ExecutionSlotsTotal: h.GetExecutionSlotsTotal(),
		QueuedJobsByType:    int32MapEntries(h.GetQueuedJobsByType()),
		Metadata:            stringMapEntries(h.GetMetadata()),
	})
}

func unmarshalWorkerHeartbeat(t pluginwire.WorkerHeartbeat) *plugin_pb.WorkerHeartbeat {
	rwl := t.RunningWork()
	rw := make([]*plugin_pb.RunningWork, rwl.Len())
	for i := range rw {
		rw[i] = unmarshalRunningWork(wrap(rwl, i, pluginwire.WrapRunningWork))
	}
	return &plugin_pb.WorkerHeartbeat{
		WorkerId:            t.WorkerID(),
		RunningWork:         rw,
		DetectionSlotsUsed:  t.DetectionSlotsUsed(),
		DetectionSlotsTotal: t.DetectionSlotsTotal(),
		ExecutionSlotsUsed:  t.ExecutionSlotsUsed(),
		ExecutionSlotsTotal: t.ExecutionSlotsTotal(),
		QueuedJobsByType:    readInt32Map(t.QueuedJobsByType()),
		Metadata:            readStringMap(t.Metadata()),
	}
}

func marshalWorkerAcknowledge(a *plugin_pb.WorkerAcknowledge) []byte {
	if a == nil {
		return nil
	}
	return pluginwire.NewWorkerAcknowledge(pluginwire.WorkerAcknowledgeInput{
		RequestID: a.GetRequestId(),
		Accepted:  a.GetAccepted(),
		Message:   a.GetMessage(),
	})
}

func unmarshalWorkerAcknowledge(t pluginwire.WorkerAcknowledge) *plugin_pb.WorkerAcknowledge {
	return &plugin_pb.WorkerAcknowledge{
		RequestId: t.RequestID(),
		Accepted:  t.Accepted(),
		Message:   t.Message(),
	}
}

func marshalRunningWork(w *plugin_pb.RunningWork) []byte {
	if w == nil {
		return nil
	}
	return pluginwire.NewRunningWork(pluginwire.RunningWorkInput{
		WorkID:          w.GetWorkId(),
		Kind:            uint32(w.GetKind()),
		JobType:         w.GetJobType(),
		State:           uint32(w.GetState()),
		ProgressPercent: w.GetProgressPercent(),
		Stage:           w.GetStage(),
	})
}

func unmarshalRunningWork(t pluginwire.RunningWork) *plugin_pb.RunningWork {
	return &plugin_pb.RunningWork{
		WorkId:          t.WorkID(),
		Kind:            plugin_pb.WorkKind(t.Kind()),
		JobType:         t.JobType(),
		State:           plugin_pb.JobState(t.State()),
		ProgressPercent: t.ProgressPercent(),
		Stage:           t.Stage(),
	}
}

func marshalJobTypeCapability(c *plugin_pb.JobTypeCapability) []byte {
	if c == nil {
		return nil
	}
	return pluginwire.NewJobTypeCapability(pluginwire.JobTypeCapabilityInput{
		JobType:                 c.GetJobType(),
		CanDetect:               c.GetCanDetect(),
		CanExecute:              c.GetCanExecute(),
		MaxDetectionConcurrency: c.GetMaxDetectionConcurrency(),
		MaxExecutionConcurrency: c.GetMaxExecutionConcurrency(),
		DisplayName:             c.GetDisplayName(),
		Description:             c.GetDescription(),
		Weight:                  c.GetWeight(),
	})
}

func unmarshalJobTypeCapability(t pluginwire.JobTypeCapability) *plugin_pb.JobTypeCapability {
	return &plugin_pb.JobTypeCapability{
		JobType:                 t.JobType(),
		CanDetect:               t.CanDetect(),
		CanExecute:              t.CanExecute(),
		MaxDetectionConcurrency: t.MaxDetectionConcurrency(),
		MaxExecutionConcurrency: t.MaxExecutionConcurrency(),
		DisplayName:             t.DisplayName(),
		Description:             t.Description(),
		Weight:                  t.Weight(),
	}
}

// ===================== detection / execution messages =====================

func marshalRequestConfigSchema(r *plugin_pb.RequestConfigSchema) []byte {
	if r == nil {
		return nil
	}
	return pluginwire.NewRequestConfigSchema(pluginwire.RequestConfigSchemaInput{
		JobType:      r.GetJobType(),
		ForceRefresh: r.GetForceRefresh(),
	})
}

func unmarshalRequestConfigSchema(t pluginwire.RequestConfigSchema) *plugin_pb.RequestConfigSchema {
	return &plugin_pb.RequestConfigSchema{
		JobType:      t.JobType(),
		ForceRefresh: t.ForceRefresh(),
	}
}

func marshalConfigSchemaResponse(r *plugin_pb.ConfigSchemaResponse) []byte {
	if r == nil {
		return nil
	}
	return pluginwire.NewConfigSchemaResponse(pluginwire.ConfigSchemaResponseInput{
		RequestID:         r.GetRequestId(),
		JobType:           r.GetJobType(),
		Success:           r.GetSuccess(),
		ErrorMessage:      r.GetErrorMessage(),
		JobTypeDescriptor: marshalJobTypeDescriptor(r.GetJobTypeDescriptor()),
	})
}

func unmarshalConfigSchemaResponse(t pluginwire.ConfigSchemaResponse) *plugin_pb.ConfigSchemaResponse {
	r := &plugin_pb.ConfigSchemaResponse{
		RequestId:    t.RequestID(),
		JobType:      t.JobType(),
		Success:      t.Success(),
		ErrorMessage: t.ErrorMessage(),
	}
	r.JobTypeDescriptor = unmarshalJobTypeDescriptor(t.JobTypeDescriptor())
	return r
}

func marshalRunDetectionRequest(r *plugin_pb.RunDetectionRequest) []byte {
	if r == nil {
		return nil
	}
	return pluginwire.NewRunDetectionRequest(pluginwire.RunDetectionRequestInput{
		RequestID:          r.GetRequestId(),
		JobType:            r.GetJobType(),
		DetectionSequence:  r.GetDetectionSequence(),
		AdminRuntime:       marshalAdminRuntimeConfig(r.GetAdminRuntime()),
		AdminConfigValues:  configValueMapEntries(r.GetAdminConfigValues()),
		WorkerConfigValues: configValueMapEntries(r.GetWorkerConfigValues()),
		ClusterContext:     marshalClusterContext(r.GetClusterContext()),
		LastSuccessfulRun:  tsNanos(r.GetLastSuccessfulRun()),
		MaxResults:         r.GetMaxResults(),
	})
}

func unmarshalRunDetectionRequest(t pluginwire.RunDetectionRequest) *plugin_pb.RunDetectionRequest {
	return &plugin_pb.RunDetectionRequest{
		RequestId:          t.RequestID(),
		JobType:            t.JobType(),
		DetectionSequence:  t.DetectionSequence(),
		AdminRuntime:       unmarshalAdminRuntimeConfig(t.AdminRuntime()),
		AdminConfigValues:  readConfigValueMap(t.AdminConfigValues()),
		WorkerConfigValues: readConfigValueMap(t.WorkerConfigValues()),
		ClusterContext:     unmarshalClusterContext(t.ClusterContext()),
		LastSuccessfulRun:  nanosTS(t.LastSuccessfulRun()),
		MaxResults:         t.MaxResults(),
	}
}

func marshalDetectionProposals(d *plugin_pb.DetectionProposals) []byte {
	if d == nil {
		return nil
	}
	props := make([][]byte, 0, len(d.GetProposals()))
	for _, p := range d.GetProposals() {
		props = append(props, marshalJobProposal(p))
	}
	return pluginwire.NewDetectionProposals(pluginwire.DetectionProposalsInput{
		RequestID: d.GetRequestId(),
		JobType:   d.GetJobType(),
		Proposals: props,
		HasMore:   d.GetHasMore(),
	})
}

func unmarshalDetectionProposals(t pluginwire.DetectionProposals) *plugin_pb.DetectionProposals {
	pl := t.Proposals()
	props := make([]*plugin_pb.JobProposal, pl.Len())
	for i := range props {
		props[i] = unmarshalJobProposal(wrap(pl, i, pluginwire.WrapJobProposal))
	}
	return &plugin_pb.DetectionProposals{
		RequestId: t.RequestID(),
		JobType:   t.JobType(),
		Proposals: props,
		HasMore:   t.HasMore(),
	}
}

func marshalDetectionComplete(d *plugin_pb.DetectionComplete) []byte {
	if d == nil {
		return nil
	}
	return pluginwire.NewDetectionComplete(pluginwire.DetectionCompleteInput{
		RequestID:      d.GetRequestId(),
		JobType:        d.GetJobType(),
		Success:        d.GetSuccess(),
		ErrorMessage:   d.GetErrorMessage(),
		TotalProposals: d.GetTotalProposals(),
	})
}

func unmarshalDetectionComplete(t pluginwire.DetectionComplete) *plugin_pb.DetectionComplete {
	return &plugin_pb.DetectionComplete{
		RequestId:      t.RequestID(),
		JobType:        t.JobType(),
		Success:        t.Success(),
		ErrorMessage:   t.ErrorMessage(),
		TotalProposals: t.TotalProposals(),
	}
}

func marshalJobProposal(p *plugin_pb.JobProposal) []byte {
	if p == nil {
		return nil
	}
	return pluginwire.NewJobProposal(pluginwire.JobProposalInput{
		ProposalID: p.GetProposalId(),
		DedupeKey:  p.GetDedupeKey(),
		JobType:    p.GetJobType(),
		Priority:   uint32(p.GetPriority()),
		Summary:    p.GetSummary(),
		Detail:     p.GetDetail(),
		Parameters: configValueMapEntries(p.GetParameters()),
		Labels:     stringMapEntries(p.GetLabels()),
		NotBefore:  tsNanos(p.GetNotBefore()),
		ExpiresAt:  tsNanos(p.GetExpiresAt()),
	})
}

func unmarshalJobProposal(t pluginwire.JobProposal) *plugin_pb.JobProposal {
	return &plugin_pb.JobProposal{
		ProposalId: t.ProposalID(),
		DedupeKey:  t.DedupeKey(),
		JobType:    t.JobType(),
		Priority:   plugin_pb.JobPriority(t.Priority()),
		Summary:    t.Summary(),
		Detail:     t.Detail(),
		Parameters: readConfigValueMap(t.Parameters()),
		Labels:     readStringMap(t.Labels()),
		NotBefore:  nanosTS(t.NotBefore()),
		ExpiresAt:  nanosTS(t.ExpiresAt()),
	}
}

func marshalExecuteJobRequest(r *plugin_pb.ExecuteJobRequest) []byte {
	if r == nil {
		return nil
	}
	return pluginwire.NewExecuteJobRequest(pluginwire.ExecuteJobRequestInput{
		RequestID:          r.GetRequestId(),
		Job:                marshalJobSpec(r.GetJob()),
		AdminRuntime:       marshalAdminRuntimeConfig(r.GetAdminRuntime()),
		AdminConfigValues:  configValueMapEntries(r.GetAdminConfigValues()),
		WorkerConfigValues: configValueMapEntries(r.GetWorkerConfigValues()),
		ClusterContext:     marshalClusterContext(r.GetClusterContext()),
		Attempt:            r.GetAttempt(),
	})
}

func unmarshalExecuteJobRequest(t pluginwire.ExecuteJobRequest) *plugin_pb.ExecuteJobRequest {
	return &plugin_pb.ExecuteJobRequest{
		RequestId:          t.RequestID(),
		Job:                unmarshalJobSpec(t.Job()),
		AdminRuntime:       unmarshalAdminRuntimeConfig(t.AdminRuntime()),
		AdminConfigValues:  readConfigValueMap(t.AdminConfigValues()),
		WorkerConfigValues: readConfigValueMap(t.WorkerConfigValues()),
		ClusterContext:     unmarshalClusterContext(t.ClusterContext()),
		Attempt:            t.Attempt(),
	}
}

func marshalJobSpec(s *plugin_pb.JobSpec) []byte {
	if s == nil {
		return nil
	}
	return pluginwire.NewJobSpec(pluginwire.JobSpecInput{
		JobID:       s.GetJobId(),
		JobType:     s.GetJobType(),
		DedupeKey:   s.GetDedupeKey(),
		Priority:    uint32(s.GetPriority()),
		Summary:     s.GetSummary(),
		Detail:      s.GetDetail(),
		Parameters:  configValueMapEntries(s.GetParameters()),
		Labels:      stringMapEntries(s.GetLabels()),
		CreatedAt:   tsNanos(s.GetCreatedAt()),
		ScheduledAt: tsNanos(s.GetScheduledAt()),
	})
}

func unmarshalJobSpec(t pluginwire.JobSpec) *plugin_pb.JobSpec {
	return orNil(&plugin_pb.JobSpec{
		JobId:       t.JobID(),
		JobType:     t.JobType(),
		DedupeKey:   t.DedupeKey(),
		Priority:    plugin_pb.JobPriority(t.Priority()),
		Summary:     t.Summary(),
		Detail:      t.Detail(),
		Parameters:  readConfigValueMap(t.Parameters()),
		Labels:      readStringMap(t.Labels()),
		CreatedAt:   nanosTS(t.CreatedAt()),
		ScheduledAt: nanosTS(t.ScheduledAt()),
	})
}

func marshalJobProgressUpdate(u *plugin_pb.JobProgressUpdate) []byte {
	if u == nil {
		return nil
	}
	acts := make([][]byte, 0, len(u.GetActivities()))
	for _, a := range u.GetActivities() {
		acts = append(acts, marshalActivityEvent(a))
	}
	return pluginwire.NewJobProgressUpdate(pluginwire.JobProgressUpdateInput{
		RequestID:       u.GetRequestId(),
		JobID:           u.GetJobId(),
		JobType:         u.GetJobType(),
		State:           uint32(u.GetState()),
		ProgressPercent: u.GetProgressPercent(),
		Stage:           u.GetStage(),
		Message:         u.GetMessage(),
		Metrics:         configValueMapEntries(u.GetMetrics()),
		Activities:      acts,
		UpdatedAt:       tsNanos(u.GetUpdatedAt()),
	})
}

func unmarshalJobProgressUpdate(t pluginwire.JobProgressUpdate) *plugin_pb.JobProgressUpdate {
	al := t.Activities()
	acts := make([]*plugin_pb.ActivityEvent, al.Len())
	for i := range acts {
		acts[i] = unmarshalActivityEvent(wrap(al, i, pluginwire.WrapActivityEvent))
	}
	return &plugin_pb.JobProgressUpdate{
		RequestId:       t.RequestID(),
		JobId:           t.JobID(),
		JobType:         t.JobType(),
		State:           plugin_pb.JobState(t.State()),
		ProgressPercent: t.ProgressPercent(),
		Stage:           t.Stage(),
		Message:         t.Message(),
		Metrics:         readConfigValueMap(t.Metrics()),
		Activities:      acts,
		UpdatedAt:       nanosTS(t.UpdatedAt()),
	}
}

func marshalJobCompleted(c *plugin_pb.JobCompleted) []byte {
	if c == nil {
		return nil
	}
	acts := make([][]byte, 0, len(c.GetActivities()))
	for _, a := range c.GetActivities() {
		acts = append(acts, marshalActivityEvent(a))
	}
	return pluginwire.NewJobCompleted(pluginwire.JobCompletedInput{
		RequestID:    c.GetRequestId(),
		JobID:        c.GetJobId(),
		JobType:      c.GetJobType(),
		Success:      c.GetSuccess(),
		ErrorMessage: c.GetErrorMessage(),
		Result:       marshalJobResult(c.GetResult()),
		Activities:   acts,
		CompletedAt:  tsNanos(c.GetCompletedAt()),
	})
}

func unmarshalJobCompleted(t pluginwire.JobCompleted) *plugin_pb.JobCompleted {
	al := t.Activities()
	acts := make([]*plugin_pb.ActivityEvent, al.Len())
	for i := range acts {
		acts[i] = unmarshalActivityEvent(wrap(al, i, pluginwire.WrapActivityEvent))
	}
	return &plugin_pb.JobCompleted{
		RequestId:    t.RequestID(),
		JobId:        t.JobID(),
		JobType:      t.JobType(),
		Success:      t.Success(),
		ErrorMessage: t.ErrorMessage(),
		Result:       unmarshalJobResult(t.Result()),
		Activities:   acts,
		CompletedAt:  nanosTS(t.CompletedAt()),
	}
}

func marshalJobResult(r *plugin_pb.JobResult) []byte {
	if r == nil {
		return nil
	}
	return pluginwire.NewJobResult(pluginwire.JobResultInput{
		OutputValues: configValueMapEntries(r.GetOutputValues()),
		Summary:      r.GetSummary(),
	})
}

func unmarshalJobResult(t pluginwire.JobResult) *plugin_pb.JobResult {
	return orNil(&plugin_pb.JobResult{
		OutputValues: readConfigValueMap(t.OutputValues()),
		Summary:      t.Summary(),
	})
}

func marshalClusterContext(c *plugin_pb.ClusterContext) []byte {
	if c == nil {
		return nil
	}
	return pluginwire.NewClusterContext(pluginwire.ClusterContextInput{
		MasterGRPCAddresses: strBytes(c.GetMasterGrpcAddresses()),
		FilerAddresses:      strBytes(c.GetFilerAddresses()),
		VolumeGRPCAddresses: strBytes(c.GetVolumeGrpcAddresses()),
		Metadata:            stringMapEntries(c.GetMetadata()),
		S3GRPCAddresses:     strBytes(c.GetS3GrpcAddresses()),
	})
}

func unmarshalClusterContext(t pluginwire.ClusterContext) *plugin_pb.ClusterContext {
	return orNil(&plugin_pb.ClusterContext{
		MasterGrpcAddresses: readStrings(t.MasterGRPCAddresses()),
		FilerAddresses:      readStrings(t.FilerAddresses()),
		VolumeGrpcAddresses: readStrings(t.VolumeGRPCAddresses()),
		Metadata:            readStringMap(t.Metadata()),
		S3GrpcAddresses:     readStrings(t.S3GRPCAddresses()),
	})
}

func marshalActivityEvent(a *plugin_pb.ActivityEvent) []byte {
	if a == nil {
		return nil
	}
	return pluginwire.NewActivityEvent(pluginwire.ActivityEventInput{
		Source:    uint32(a.GetSource()),
		Message:   a.GetMessage(),
		Stage:     a.GetStage(),
		Details:   configValueMapEntries(a.GetDetails()),
		CreatedAt: tsNanos(a.GetCreatedAt()),
	})
}

func unmarshalActivityEvent(t pluginwire.ActivityEvent) *plugin_pb.ActivityEvent {
	return &plugin_pb.ActivityEvent{
		Source:    plugin_pb.ActivitySource(t.Source()),
		Message:   t.Message(),
		Stage:     t.Stage(),
		Details:   readConfigValueMap(t.Details()),
		CreatedAt: nanosTS(t.CreatedAt()),
	}
}

func marshalCancelRequest(c *plugin_pb.CancelRequest) []byte {
	if c == nil {
		return nil
	}
	return pluginwire.NewCancelRequest(pluginwire.CancelRequestInput{
		TargetID:   c.GetTargetId(),
		TargetKind: uint32(c.GetTargetKind()),
		Reason:     c.GetReason(),
		Force:      c.GetForce(),
	})
}

func unmarshalCancelRequest(t pluginwire.CancelRequest) *plugin_pb.CancelRequest {
	return &plugin_pb.CancelRequest{
		TargetId:   t.TargetID(),
		TargetKind: plugin_pb.WorkKind(t.TargetKind()),
		Reason:     t.Reason(),
		Force:      t.Force(),
	}
}

func marshalAdminShutdown(s *plugin_pb.AdminShutdown) []byte {
	if s == nil {
		return nil
	}
	return pluginwire.NewAdminShutdown(pluginwire.AdminShutdownInput{
		Reason:             s.GetReason(),
		GracePeriodSeconds: s.GetGracePeriodSeconds(),
	})
}

func unmarshalAdminShutdown(t pluginwire.AdminShutdown) *plugin_pb.AdminShutdown {
	return &plugin_pb.AdminShutdown{
		Reason:             t.Reason(),
		GracePeriodSeconds: t.GracePeriodSeconds(),
	}
}

// ===================== descriptor / config-form cluster =====================

func marshalJobTypeDescriptor(d *plugin_pb.JobTypeDescriptor) []byte {
	if d == nil {
		return nil
	}
	return pluginwire.NewJobTypeDescriptor(pluginwire.JobTypeDescriptorInput{
		JobType:              d.GetJobType(),
		DisplayName:          d.GetDisplayName(),
		Description:          d.GetDescription(),
		Icon:                 d.GetIcon(),
		DescriptorVersion:    d.GetDescriptorVersion(),
		AdminConfigForm:      marshalConfigForm(d.GetAdminConfigForm()),
		WorkerConfigForm:     marshalConfigForm(d.GetWorkerConfigForm()),
		AdminRuntimeDefaults: marshalAdminRuntimeDefaults(d.GetAdminRuntimeDefaults()),
		WorkerDefaultValues:  configValueMapEntries(d.GetWorkerDefaultValues()),
	})
}

func unmarshalJobTypeDescriptor(t pluginwire.JobTypeDescriptor) *plugin_pb.JobTypeDescriptor {
	return orNil(&plugin_pb.JobTypeDescriptor{
		JobType:              t.JobType(),
		DisplayName:          t.DisplayName(),
		Description:          t.Description(),
		Icon:                 t.Icon(),
		DescriptorVersion:    t.DescriptorVersion(),
		AdminConfigForm:      unmarshalConfigForm(t.AdminConfigForm()),
		WorkerConfigForm:     unmarshalConfigForm(t.WorkerConfigForm()),
		AdminRuntimeDefaults: unmarshalAdminRuntimeDefaults(t.AdminRuntimeDefaults()),
		WorkerDefaultValues:  readConfigValueMap(t.WorkerDefaultValues()),
	})
}

func marshalConfigForm(f *plugin_pb.ConfigForm) []byte {
	if f == nil {
		return nil
	}
	secs := make([][]byte, 0, len(f.GetSections()))
	for _, s := range f.GetSections() {
		secs = append(secs, marshalConfigSection(s))
	}
	return pluginwire.NewConfigForm(pluginwire.ConfigFormInput{
		FormID:        f.GetFormId(),
		Title:         f.GetTitle(),
		Description:   f.GetDescription(),
		Sections:      secs,
		DefaultValues: configValueMapEntries(f.GetDefaultValues()),
	})
}

func unmarshalConfigForm(t pluginwire.ConfigForm) *plugin_pb.ConfigForm {
	if t.FormID() == "" && t.Title() == "" && t.Description() == "" &&
		t.Sections().Len() == 0 && t.DefaultValues().Len() == 0 {
		return nil
	}
	sl := t.Sections()
	secs := make([]*plugin_pb.ConfigSection, sl.Len())
	for i := range secs {
		secs[i] = unmarshalConfigSection(wrap(sl, i, pluginwire.WrapConfigSection))
	}
	return &plugin_pb.ConfigForm{
		FormId:        t.FormID(),
		Title:         t.Title(),
		Description:   t.Description(),
		Sections:      secs,
		DefaultValues: readConfigValueMap(t.DefaultValues()),
	}
}

func marshalConfigSection(s *plugin_pb.ConfigSection) []byte {
	if s == nil {
		return nil
	}
	fields := make([][]byte, 0, len(s.GetFields()))
	for _, f := range s.GetFields() {
		fields = append(fields, marshalConfigField(f))
	}
	return pluginwire.NewConfigSection(pluginwire.ConfigSectionInput{
		SectionID:   s.GetSectionId(),
		Title:       s.GetTitle(),
		Description: s.GetDescription(),
		Fields:      fields,
	})
}

func unmarshalConfigSection(t pluginwire.ConfigSection) *plugin_pb.ConfigSection {
	fl := t.Fields()
	fields := make([]*plugin_pb.ConfigField, fl.Len())
	for i := range fields {
		fields[i] = unmarshalConfigField(wrap(fl, i, pluginwire.WrapConfigField))
	}
	return &plugin_pb.ConfigSection{
		SectionId:   t.SectionID(),
		Title:       t.Title(),
		Description: t.Description(),
		Fields:      fields,
	}
}

func marshalConfigField(f *plugin_pb.ConfigField) []byte {
	if f == nil {
		return nil
	}
	opts := make([][]byte, 0, len(f.GetOptions()))
	for _, o := range f.GetOptions() {
		opts = append(opts, marshalConfigOption(o))
	}
	rules := make([][]byte, 0, len(f.GetValidationRules()))
	for _, r := range f.GetValidationRules() {
		rules = append(rules, marshalValidationRule(r))
	}
	in := pluginwire.ConfigFieldInput{
		Name:             f.GetName(),
		Label:            f.GetLabel(),
		Description:      f.GetDescription(),
		HelpText:         f.GetHelpText(),
		Placeholder:      f.GetPlaceholder(),
		FieldType:        uint32(f.GetFieldType()),
		Widget:           uint32(f.GetWidget()),
		Required:         f.GetRequired(),
		ReadOnly:         f.GetReadOnly(),
		Sensitive:        f.GetSensitive(),
		Options:          opts,
		ValidationRules:  rules,
		VisibleWhenField: f.GetVisibleWhenField(),
	}
	if f.GetMinValue() != nil {
		in.MinValue = marshalConfigValue(f.GetMinValue())
	}
	if f.GetMaxValue() != nil {
		in.MaxValue = marshalConfigValue(f.GetMaxValue())
	}
	if f.GetVisibleWhenEquals() != nil {
		in.VisibleWhenEquals = marshalConfigValue(f.GetVisibleWhenEquals())
	}
	return pluginwire.NewConfigField(in)
}

func unmarshalConfigField(t pluginwire.ConfigField) *plugin_pb.ConfigField {
	ol := t.Options()
	opts := make([]*plugin_pb.ConfigOption, ol.Len())
	for i := range opts {
		opts[i] = unmarshalConfigOption(wrap(ol, i, pluginwire.WrapConfigOption))
	}
	rl := t.ValidationRules()
	rules := make([]*plugin_pb.ValidationRule, rl.Len())
	for i := range rules {
		rules[i] = unmarshalValidationRule(wrap(rl, i, pluginwire.WrapValidationRule))
	}
	return &plugin_pb.ConfigField{
		Name:              t.Name(),
		Label:             t.Label(),
		Description:       t.Description(),
		HelpText:          t.HelpText(),
		Placeholder:       t.Placeholder(),
		FieldType:         plugin_pb.ConfigFieldType(t.FieldType()),
		Widget:            plugin_pb.ConfigWidget(t.Widget()),
		Required:          t.Required(),
		ReadOnly:          t.ReadOnly(),
		Sensitive:         t.Sensitive(),
		MinValue:          unmarshalConfigValue(t.MinValue()),
		MaxValue:          unmarshalConfigValue(t.MaxValue()),
		Options:           opts,
		ValidationRules:   rules,
		VisibleWhenField:  t.VisibleWhenField(),
		VisibleWhenEquals: unmarshalConfigValue(t.VisibleWhenEquals()),
	}
}

func marshalConfigOption(o *plugin_pb.ConfigOption) []byte {
	if o == nil {
		return nil
	}
	return pluginwire.NewConfigOption(pluginwire.ConfigOptionInput{
		Value:       o.GetValue(),
		Label:       o.GetLabel(),
		Description: o.GetDescription(),
		Disabled:    o.GetDisabled(),
	})
}

func unmarshalConfigOption(t pluginwire.ConfigOption) *plugin_pb.ConfigOption {
	return &plugin_pb.ConfigOption{
		Value:       t.Value(),
		Label:       t.Label(),
		Description: t.Description(),
		Disabled:    t.Disabled(),
	}
}

func marshalValidationRule(r *plugin_pb.ValidationRule) []byte {
	if r == nil {
		return nil
	}
	return pluginwire.NewValidationRule(pluginwire.ValidationRuleInput{
		Type:         uint32(r.GetType()),
		Expression:   r.GetExpression(),
		ErrorMessage: r.GetErrorMessage(),
	})
}

func unmarshalValidationRule(t pluginwire.ValidationRule) *plugin_pb.ValidationRule {
	return &plugin_pb.ValidationRule{
		Type:         plugin_pb.ValidationRuleType(t.Type()),
		Expression:   t.Expression(),
		ErrorMessage: t.ErrorMessage(),
	}
}

func marshalAdminRuntimeDefaults(d *plugin_pb.AdminRuntimeDefaults) []byte {
	if d == nil {
		return nil
	}
	return pluginwire.NewAdminRuntimeDefaults(pluginwire.AdminRuntimeDefaultsInput{
		Enabled:                       d.GetEnabled(),
		DetectionIntervalMinutes:      d.GetDetectionIntervalMinutes(),
		DetectionTimeoutSeconds:       d.GetDetectionTimeoutSeconds(),
		MaxJobsPerDetection:           d.GetMaxJobsPerDetection(),
		GlobalExecutionConcurrency:    d.GetGlobalExecutionConcurrency(),
		PerWorkerExecutionConcurrency: d.GetPerWorkerExecutionConcurrency(),
		RetryLimit:                    d.GetRetryLimit(),
		RetryBackoffSeconds:           d.GetRetryBackoffSeconds(),
		JobTypeMaxRuntimeSeconds:      d.GetJobTypeMaxRuntimeSeconds(),
		ExecutionTimeoutSeconds:       d.GetExecutionTimeoutSeconds(),
	})
}

func unmarshalAdminRuntimeDefaults(t pluginwire.AdminRuntimeDefaults) *plugin_pb.AdminRuntimeDefaults {
	return orNil(&plugin_pb.AdminRuntimeDefaults{
		Enabled:                       t.Enabled(),
		DetectionIntervalMinutes:      t.DetectionIntervalMinutes(),
		DetectionTimeoutSeconds:       t.DetectionTimeoutSeconds(),
		MaxJobsPerDetection:           t.MaxJobsPerDetection(),
		GlobalExecutionConcurrency:    t.GlobalExecutionConcurrency(),
		PerWorkerExecutionConcurrency: t.PerWorkerExecutionConcurrency(),
		RetryLimit:                    t.RetryLimit(),
		RetryBackoffSeconds:           t.RetryBackoffSeconds(),
		JobTypeMaxRuntimeSeconds:      t.JobTypeMaxRuntimeSeconds(),
		ExecutionTimeoutSeconds:       t.ExecutionTimeoutSeconds(),
	})
}

func marshalAdminRuntimeConfig(c *plugin_pb.AdminRuntimeConfig) []byte {
	if c == nil {
		return nil
	}
	return pluginwire.NewAdminRuntimeConfig(pluginwire.AdminRuntimeConfigInput{
		Enabled:                       c.GetEnabled(),
		DetectionIntervalMinutes:      c.GetDetectionIntervalMinutes(),
		DetectionTimeoutSeconds:       c.GetDetectionTimeoutSeconds(),
		MaxJobsPerDetection:           c.GetMaxJobsPerDetection(),
		GlobalExecutionConcurrency:    c.GetGlobalExecutionConcurrency(),
		PerWorkerExecutionConcurrency: c.GetPerWorkerExecutionConcurrency(),
		RetryLimit:                    c.GetRetryLimit(),
		RetryBackoffSeconds:           c.GetRetryBackoffSeconds(),
		JobTypeMaxRuntimeSeconds:      c.GetJobTypeMaxRuntimeSeconds(),
		ExecutionTimeoutSeconds:       c.GetExecutionTimeoutSeconds(),
	})
}

func unmarshalAdminRuntimeConfig(t pluginwire.AdminRuntimeConfig) *plugin_pb.AdminRuntimeConfig {
	return orNil(&plugin_pb.AdminRuntimeConfig{
		Enabled:                       t.Enabled(),
		DetectionIntervalMinutes:      t.DetectionIntervalMinutes(),
		DetectionTimeoutSeconds:       t.DetectionTimeoutSeconds(),
		MaxJobsPerDetection:           t.MaxJobsPerDetection(),
		GlobalExecutionConcurrency:    t.GlobalExecutionConcurrency(),
		PerWorkerExecutionConcurrency: t.PerWorkerExecutionConcurrency(),
		RetryLimit:                    t.RetryLimit(),
		RetryBackoffSeconds:           t.RetryBackoffSeconds(),
		JobTypeMaxRuntimeSeconds:      t.JobTypeMaxRuntimeSeconds(),
		ExecutionTimeoutSeconds:       t.ExecutionTimeoutSeconds(),
	})
}

// ===================== ConfigValue (recursive oneof) =====================

func marshalConfigValue(v *plugin_pb.ConfigValue) []byte {
	if v == nil {
		return nil
	}
	in := pluginwire.ConfigValueInput{}
	switch v.Kind.(type) {
	case *plugin_pb.ConfigValue_BoolValue:
		in.Which = pluginwire.ConfigValueKindBool
		in.BoolValue = v.GetBoolValue()
	case *plugin_pb.ConfigValue_Int64Value:
		in.Which = pluginwire.ConfigValueKindInt64
		in.Int64Value = v.GetInt64Value()
	case *plugin_pb.ConfigValue_DoubleValue:
		in.Which = pluginwire.ConfigValueKindDouble
		in.DoubleValue = v.GetDoubleValue()
	case *plugin_pb.ConfigValue_StringValue:
		in.Which = pluginwire.ConfigValueKindString
		in.StringValue = v.GetStringValue()
	case *plugin_pb.ConfigValue_BytesValue:
		in.Which = pluginwire.ConfigValueKindBytes
		in.BytesValue = v.GetBytesValue()
	case *plugin_pb.ConfigValue_DurationValue:
		in.Which = pluginwire.ConfigValueKindDuration
		in.DurationValue = durNanos(v.GetDurationValue())
	case *plugin_pb.ConfigValue_StringList:
		in.Which = pluginwire.ConfigValueKindStringList
		in.StringList = marshalStringList(v.GetStringList())
	case *plugin_pb.ConfigValue_Int64List:
		in.Which = pluginwire.ConfigValueKindInt64List
		in.Int64List = marshalInt64List(v.GetInt64List())
	case *plugin_pb.ConfigValue_DoubleList:
		in.Which = pluginwire.ConfigValueKindDoubleList
		in.DoubleList = marshalDoubleList(v.GetDoubleList())
	case *plugin_pb.ConfigValue_BoolList:
		in.Which = pluginwire.ConfigValueKindBoolList
		in.BoolList = marshalBoolList(v.GetBoolList())
	case *plugin_pb.ConfigValue_ListValue:
		in.Which = pluginwire.ConfigValueKindListValue
		in.ListValue = marshalValueList(v.GetListValue())
	case *plugin_pb.ConfigValue_MapValue:
		in.Which = pluginwire.ConfigValueKindMapValue
		in.MapValue = marshalValueMap(v.GetMapValue())
	}
	return pluginwire.NewConfigValue(in)
}

func unmarshalConfigValue(t pluginwire.ConfigValue) *plugin_pb.ConfigValue {
	v := &plugin_pb.ConfigValue{}
	switch t.WhichKind() {
	case pluginwire.ConfigValueKindBool:
		v.Kind = &plugin_pb.ConfigValue_BoolValue{BoolValue: t.BoolValue()}
	case pluginwire.ConfigValueKindInt64:
		v.Kind = &plugin_pb.ConfigValue_Int64Value{Int64Value: t.Int64Value()}
	case pluginwire.ConfigValueKindDouble:
		v.Kind = &plugin_pb.ConfigValue_DoubleValue{DoubleValue: t.DoubleValue()}
	case pluginwire.ConfigValueKindString:
		v.Kind = &plugin_pb.ConfigValue_StringValue{StringValue: t.StringValue()}
	case pluginwire.ConfigValueKindBytes:
		v.Kind = &plugin_pb.ConfigValue_BytesValue{BytesValue: t.BytesValue()}
	case pluginwire.ConfigValueKindDuration:
		v.Kind = &plugin_pb.ConfigValue_DurationValue{DurationValue: nanosDur(t.DurationValue())}
	case pluginwire.ConfigValueKindStringList:
		v.Kind = &plugin_pb.ConfigValue_StringList{StringList: unmarshalStringList(t.StringList())}
	case pluginwire.ConfigValueKindInt64List:
		v.Kind = &plugin_pb.ConfigValue_Int64List{Int64List: unmarshalInt64List(t.Int64List())}
	case pluginwire.ConfigValueKindDoubleList:
		v.Kind = &plugin_pb.ConfigValue_DoubleList{DoubleList: unmarshalDoubleList(t.DoubleList())}
	case pluginwire.ConfigValueKindBoolList:
		v.Kind = &plugin_pb.ConfigValue_BoolList{BoolList: unmarshalBoolList(t.BoolList())}
	case pluginwire.ConfigValueKindListValue:
		v.Kind = &plugin_pb.ConfigValue_ListValue{ListValue: unmarshalValueList(t.ListValue())}
	case pluginwire.ConfigValueKindMapValue:
		v.Kind = &plugin_pb.ConfigValue_MapValue{MapValue: unmarshalValueMap(t.MapValue())}
	default:
		return nil
	}
	return v
}

func marshalStringList(l *plugin_pb.StringList) []byte {
	if l == nil {
		return nil
	}
	return pluginwire.NewStringList(pluginwire.StringListInput{Values: strBytes(l.GetValues())})
}

func unmarshalStringList(t pluginwire.StringList) *plugin_pb.StringList {
	return &plugin_pb.StringList{Values: readStrings(t.Values())}
}

func marshalInt64List(l *plugin_pb.Int64List) []byte {
	if l == nil {
		return nil
	}
	return pluginwire.NewInt64List(pluginwire.Int64ListInput{Values: l.GetValues()})
}

func unmarshalInt64List(t pluginwire.Int64List) *plugin_pb.Int64List {
	vl := t.Values()
	out := make([]int64, vl.Len())
	for i := range out {
		out[i] = int64(vl.Uint64At(i))
	}
	return &plugin_pb.Int64List{Values: out}
}

func marshalDoubleList(l *plugin_pb.DoubleList) []byte {
	if l == nil {
		return nil
	}
	return pluginwire.NewDoubleList(pluginwire.DoubleListInput{Values: l.GetValues()})
}

func unmarshalDoubleList(t pluginwire.DoubleList) *plugin_pb.DoubleList {
	vl := t.Values()
	out := make([]float64, vl.Len())
	for i := range out {
		out[i] = math.Float64frombits(vl.Uint64At(i))
	}
	return &plugin_pb.DoubleList{Values: out}
}

func marshalBoolList(l *plugin_pb.BoolList) []byte {
	if l == nil {
		return nil
	}
	return pluginwire.NewBoolList(pluginwire.BoolListInput{Values: l.GetValues()})
}

func unmarshalBoolList(t pluginwire.BoolList) *plugin_pb.BoolList {
	vl := t.Values()
	out := make([]bool, vl.Len())
	for i := range out {
		out[i] = vl.Uint8(i) != 0
	}
	return &plugin_pb.BoolList{Values: out}
}

func marshalValueList(l *plugin_pb.ValueList) []byte {
	if l == nil {
		return nil
	}
	vals := make([][]byte, 0, len(l.GetValues()))
	for _, v := range l.GetValues() {
		vals = append(vals, marshalConfigValue(v))
	}
	return pluginwire.NewValueList(pluginwire.ValueListInput{Values: vals})
}

func unmarshalValueList(t pluginwire.ValueList) *plugin_pb.ValueList {
	vl := t.Values()
	out := make([]*plugin_pb.ConfigValue, vl.Len())
	for i := range out {
		out[i] = unmarshalConfigValue(wrap(vl, i, pluginwire.WrapConfigValue))
	}
	return &plugin_pb.ValueList{Values: out}
}

func marshalValueMap(m *plugin_pb.ValueMap) []byte {
	if m == nil {
		return nil
	}
	return pluginwire.NewValueMap(pluginwire.ValueMapInput{Fields: configValueMapEntries(m.GetFields())})
}

func unmarshalValueMap(t pluginwire.ValueMap) *plugin_pb.ValueMap {
	return &plugin_pb.ValueMap{Fields: readConfigValueMap(t.Fields())}
}

// ===================== timestamp / duration helpers =====================

// tsNanos maps a google.protobuf.Timestamp to int64 unix-nanoseconds (0 = nil).
func tsNanos(t *timestamppb.Timestamp) int64 {
	if t == nil {
		return 0
	}
	return t.AsTime().UnixNano()
}

// nanosTS maps int64 unix-nanoseconds back to a Timestamp (0 nanos = nil, the
// proto absent form).
func nanosTS(nanos int64) *timestamppb.Timestamp {
	if nanos == 0 {
		return nil
	}
	return timestamppb.New(time.Unix(0, nanos))
}

// durNanos maps a google.protobuf.Duration to int64 nanoseconds (0 = nil).
func durNanos(d *durationpb.Duration) int64 {
	if d == nil {
		return 0
	}
	return d.AsDuration().Nanoseconds()
}

// nanosDur maps int64 nanoseconds back to a Duration.
func nanosDur(nanos int64) *durationpb.Duration {
	return durationpb.New(time.Duration(nanos))
}

// ===================== shared list / map helpers =====================

// wrap reparses the i-th object element of l (a self-contained ZAP frame stored
// out of line) with the supplied Wrap constructor, ignoring the parse error
// (an absent/malformed element wraps to a zero view, matching childObject).
func wrap[T any](l zap.List, i int, w func([]byte) (T, error)) T {
	v, _ := w(l.BytesAt(i))
	return v
}

func strBytes(ss []string) [][]byte {
	if len(ss) == 0 {
		return nil
	}
	out := make([][]byte, len(ss))
	for i, s := range ss {
		out[i] = []byte(s)
	}
	return out
}

func readStrings(l zap.List) []string {
	n := l.Len()
	if n == 0 {
		return nil
	}
	out := make([]string, n)
	for i := 0; i < n; i++ {
		out[i] = string(l.BytesAt(i))
	}
	return out
}

func stringMapEntries(m map[string]string) [][]byte {
	if len(m) == 0 {
		return nil
	}
	keys := sortedKeys(m)
	out := make([][]byte, 0, len(keys))
	for _, k := range keys {
		out = append(out, pluginwire.NewStringMapEntry(pluginwire.StringMapEntryInput{Key: k, Value: m[k]}))
	}
	return out
}

func readStringMap(l zap.List) map[string]string {
	n := l.Len()
	if n == 0 {
		return nil
	}
	out := make(map[string]string, n)
	for i := 0; i < n; i++ {
		e := wrap(l, i, pluginwire.WrapStringMapEntry)
		out[e.Key()] = e.Value()
	}
	return out
}

func int32MapEntries(m map[string]int32) [][]byte {
	if len(m) == 0 {
		return nil
	}
	keys := sortedKeysInt32(m)
	out := make([][]byte, 0, len(keys))
	for _, k := range keys {
		out = append(out, pluginwire.NewInt32MapEntry(pluginwire.Int32MapEntryInput{Key: k, Value: m[k]}))
	}
	return out
}

func readInt32Map(l zap.List) map[string]int32 {
	n := l.Len()
	if n == 0 {
		return nil
	}
	out := make(map[string]int32, n)
	for i := 0; i < n; i++ {
		e := wrap(l, i, pluginwire.WrapInt32MapEntry)
		out[e.Key()] = e.Value()
	}
	return out
}

func configValueMapEntries(m map[string]*plugin_pb.ConfigValue) [][]byte {
	if len(m) == 0 {
		return nil
	}
	keys := sortedKeysCV(m)
	out := make([][]byte, 0, len(keys))
	for _, k := range keys {
		out = append(out, pluginwire.NewConfigValueMapEntry(pluginwire.ConfigValueMapEntryInput{
			Key:   k,
			Value: marshalConfigValue(m[k]),
		}))
	}
	return out
}

func readConfigValueMap(l zap.List) map[string]*plugin_pb.ConfigValue {
	n := l.Len()
	if n == 0 {
		return nil
	}
	out := make(map[string]*plugin_pb.ConfigValue, n)
	for i := 0; i < n; i++ {
		e := wrap(l, i, pluginwire.WrapConfigValueMapEntry)
		out[e.Key()] = unmarshalConfigValue(e.Value())
	}
	return out
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortedKeysInt32(m map[string]int32) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortedKeysCV(m map[string]*plugin_pb.ConfigValue) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
