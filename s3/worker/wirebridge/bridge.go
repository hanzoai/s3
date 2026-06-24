// Package wirebridge is the single, orthogonal seam between the worker
// protobuf message tree (worker_pb) and the native ZAP wire envelopes
// (workerwire) carried over the zap-proto transport.
//
// The WorkerStream bidi RPC was cut from gRPC to the ZAP transport. The
// stream frames are ZAP WorkerMessage / AdminMessage envelopes; the rest of
// the worker subsystem still speaks worker_pb (TaskParams in particular is
// consumed pb-typed across every task implementation). This package converts
// at exactly that boundary — Encode marshals a pb message to a wire envelope
// for Send, Decode wraps a received envelope back into pb for the existing
// handlers. One direction of each lives here so neither the client
// (s3/worker) nor the server (s3/admin/dash) duplicates the mapping.
package wirebridge

import (
	"github.com/hanzoai/s3/s3/pb/worker_pb"
	workerwire "github.com/hanzoai/s3/s3/wire/worker"
)

// ---- map<string,string> helpers ----

func mapToEntries(m map[string]string) []workerwire.MapEntry {
	if len(m) == 0 {
		return nil
	}
	out := make([]workerwire.MapEntry, 0, len(m))
	for k, v := range m {
		out = append(out, workerwire.MapEntry{Key: k, Value: v})
	}
	return out
}

func entriesToMap(n int, at func(i int) workerwire.MapEntry) map[string]string {
	if n == 0 {
		return nil
	}
	out := make(map[string]string, n)
	for i := 0; i < n; i++ {
		e := at(i)
		out[e.Key] = e.Value
	}
	return out
}

// ============================ WorkerMessage ============================
//
// worker -> admin. Variants: Registration, Heartbeat, TaskRequest,
// TaskUpdate, TaskComplete, Shutdown, TaskLogResponse.

// EncodeWorkerMessage marshals a worker_pb.WorkerMessage to a ZAP
// WorkerMessage envelope ready for Stream.Send.
func EncodeWorkerMessage(m *worker_pb.WorkerMessage) []byte {
	in := workerwire.WorkerMessageInput{
		WorkerID:  m.WorkerId,
		Timestamp: m.Timestamp,
	}
	switch v := m.Message.(type) {
	case *worker_pb.WorkerMessage_Registration:
		in.MessageWhich = workerwire.WorkerMessageMessageRegistration
		in.MessageValue = encodeWorkerRegistration(v.Registration)
	case *worker_pb.WorkerMessage_Heartbeat:
		in.MessageWhich = workerwire.WorkerMessageMessageHeartbeat
		in.MessageValue = encodeWorkerHeartbeat(v.Heartbeat)
	case *worker_pb.WorkerMessage_TaskRequest:
		in.MessageWhich = workerwire.WorkerMessageMessageTaskRequest
		in.MessageValue = encodeTaskRequest(v.TaskRequest)
	case *worker_pb.WorkerMessage_TaskUpdate:
		in.MessageWhich = workerwire.WorkerMessageMessageTaskUpdate
		in.MessageValue = encodeTaskUpdate(v.TaskUpdate)
	case *worker_pb.WorkerMessage_TaskComplete:
		in.MessageWhich = workerwire.WorkerMessageMessageTaskComplete
		in.MessageValue = encodeTaskComplete(v.TaskComplete)
	case *worker_pb.WorkerMessage_Shutdown:
		in.MessageWhich = workerwire.WorkerMessageMessageShutdown
		in.MessageValue = encodeWorkerShutdown(v.Shutdown)
	case *worker_pb.WorkerMessage_TaskLogResponse:
		in.MessageWhich = workerwire.WorkerMessageMessageTaskLogResponse
		in.MessageValue = encodeTaskLogResponse(v.TaskLogResponse)
	}
	return workerwire.NewWorkerMessage(in)
}

// DecodeWorkerMessage wraps a received ZAP WorkerMessage envelope back into a
// worker_pb.WorkerMessage for the existing server handlers.
func DecodeWorkerMessage(env []byte) (*worker_pb.WorkerMessage, error) {
	w, err := workerwire.WrapWorkerMessage(env)
	if err != nil {
		return nil, err
	}
	m := &worker_pb.WorkerMessage{
		WorkerId:  w.WorkerID(),
		Timestamp: w.Timestamp(),
	}
	switch w.WhichMessage() {
	case workerwire.WorkerMessageMessageRegistration:
		m.Message = &worker_pb.WorkerMessage_Registration{Registration: decodeWorkerRegistration(w.Registration())}
	case workerwire.WorkerMessageMessageHeartbeat:
		m.Message = &worker_pb.WorkerMessage_Heartbeat{Heartbeat: decodeWorkerHeartbeat(w.Heartbeat())}
	case workerwire.WorkerMessageMessageTaskRequest:
		m.Message = &worker_pb.WorkerMessage_TaskRequest{TaskRequest: decodeTaskRequest(w.TaskRequest())}
	case workerwire.WorkerMessageMessageTaskUpdate:
		m.Message = &worker_pb.WorkerMessage_TaskUpdate{TaskUpdate: decodeTaskUpdate(w.TaskUpdate())}
	case workerwire.WorkerMessageMessageTaskComplete:
		m.Message = &worker_pb.WorkerMessage_TaskComplete{TaskComplete: decodeTaskComplete(w.TaskComplete())}
	case workerwire.WorkerMessageMessageShutdown:
		m.Message = &worker_pb.WorkerMessage_Shutdown{Shutdown: decodeWorkerShutdown(w.Shutdown())}
	case workerwire.WorkerMessageMessageTaskLogResponse:
		m.Message = &worker_pb.WorkerMessage_TaskLogResponse{TaskLogResponse: decodeTaskLogResponse(w.TaskLogResponse())}
	}
	return m, nil
}

// ============================ AdminMessage ============================
//
// admin -> worker. Variants: RegistrationResponse, HeartbeatResponse,
// TaskAssignment, TaskCancellation, AdminShutdown, TaskLogRequest.

// EncodeAdminMessage marshals a worker_pb.AdminMessage to a ZAP AdminMessage
// envelope ready for Stream.Send.
func EncodeAdminMessage(m *worker_pb.AdminMessage) []byte {
	in := workerwire.AdminMessageInput{
		AdminID:   m.AdminId,
		Timestamp: m.Timestamp,
	}
	switch v := m.Message.(type) {
	case *worker_pb.AdminMessage_RegistrationResponse:
		in.MessageWhich = workerwire.AdminMessageMessageRegistrationResponse
		in.MessageValue = encodeRegistrationResponse(v.RegistrationResponse)
	case *worker_pb.AdminMessage_HeartbeatResponse:
		in.MessageWhich = workerwire.AdminMessageMessageHeartbeatResponse
		in.MessageValue = encodeHeartbeatResponse(v.HeartbeatResponse)
	case *worker_pb.AdminMessage_TaskAssignment:
		in.MessageWhich = workerwire.AdminMessageMessageTaskAssignment
		in.MessageValue = encodeTaskAssignment(v.TaskAssignment)
	case *worker_pb.AdminMessage_TaskCancellation:
		in.MessageWhich = workerwire.AdminMessageMessageTaskCancellation
		in.MessageValue = encodeTaskCancellation(v.TaskCancellation)
	case *worker_pb.AdminMessage_AdminShutdown:
		in.MessageWhich = workerwire.AdminMessageMessageAdminShutdown
		in.MessageValue = encodeAdminShutdown(v.AdminShutdown)
	case *worker_pb.AdminMessage_TaskLogRequest:
		in.MessageWhich = workerwire.AdminMessageMessageTaskLogRequest
		in.MessageValue = encodeTaskLogRequest(v.TaskLogRequest)
	}
	return workerwire.NewAdminMessage(in)
}

// DecodeAdminMessage wraps a received ZAP AdminMessage envelope back into a
// worker_pb.AdminMessage for the existing client handlers.
func DecodeAdminMessage(env []byte) (*worker_pb.AdminMessage, error) {
	a, err := workerwire.WrapAdminMessage(env)
	if err != nil {
		return nil, err
	}
	m := &worker_pb.AdminMessage{
		AdminId:   a.AdminID(),
		Timestamp: a.Timestamp(),
	}
	switch a.WhichMessage() {
	case workerwire.AdminMessageMessageRegistrationResponse:
		m.Message = &worker_pb.AdminMessage_RegistrationResponse{RegistrationResponse: decodeRegistrationResponse(a.RegistrationResponse())}
	case workerwire.AdminMessageMessageHeartbeatResponse:
		m.Message = &worker_pb.AdminMessage_HeartbeatResponse{HeartbeatResponse: decodeHeartbeatResponse(a.HeartbeatResponse())}
	case workerwire.AdminMessageMessageTaskAssignment:
		m.Message = &worker_pb.AdminMessage_TaskAssignment{TaskAssignment: decodeTaskAssignment(a.TaskAssignment())}
	case workerwire.AdminMessageMessageTaskCancellation:
		m.Message = &worker_pb.AdminMessage_TaskCancellation{TaskCancellation: decodeTaskCancellation(a.TaskCancellation())}
	case workerwire.AdminMessageMessageAdminShutdown:
		m.Message = &worker_pb.AdminMessage_AdminShutdown{AdminShutdown: decodeAdminShutdown(a.AdminShutdown())}
	case workerwire.AdminMessageMessageTaskLogRequest:
		m.Message = &worker_pb.AdminMessage_TaskLogRequest{TaskLogRequest: decodeTaskLogRequest(a.TaskLogRequest())}
	}
	return m, nil
}

// ---------------- WorkerMessage variants ----------------

func encodeWorkerRegistration(r *worker_pb.WorkerRegistration) []byte {
	return workerwire.NewWorkerRegistration(workerwire.WorkerRegistrationInput{
		WorkerID:      r.WorkerId,
		Address:       r.Address,
		Capabilities:  r.Capabilities,
		MaxConcurrent: r.MaxConcurrent,
		Metadata:      mapToEntries(r.Metadata),
	})
}

func decodeWorkerRegistration(r workerwire.WorkerRegistration) *worker_pb.WorkerRegistration {
	caps := make([]string, r.CapabilitiesLen())
	for i := range caps {
		caps[i] = r.Capability(i)
	}
	return &worker_pb.WorkerRegistration{
		WorkerId:      r.WorkerID(),
		Address:       r.Address(),
		Capabilities:  caps,
		MaxConcurrent: r.MaxConcurrent(),
		Metadata:      entriesToMap(r.MetadataLen(), r.Metadata),
	}
}

func encodeWorkerHeartbeat(h *worker_pb.WorkerHeartbeat) []byte {
	return workerwire.NewWorkerHeartbeat(workerwire.WorkerHeartbeatInput{
		WorkerID:       h.WorkerId,
		Status:         h.Status,
		CurrentLoad:    h.CurrentLoad,
		MaxConcurrent:  h.MaxConcurrent,
		CurrentTaskIDs: h.CurrentTaskIds,
		TasksCompleted: h.TasksCompleted,
		TasksFailed:    h.TasksFailed,
		UptimeSeconds:  h.UptimeSeconds,
	})
}

func decodeWorkerHeartbeat(h workerwire.WorkerHeartbeat) *worker_pb.WorkerHeartbeat {
	ids := make([]string, h.CurrentTaskIDsLen())
	for i := range ids {
		ids[i] = h.CurrentTaskID(i)
	}
	return &worker_pb.WorkerHeartbeat{
		WorkerId:       h.WorkerID(),
		Status:         h.Status(),
		CurrentLoad:    h.CurrentLoad(),
		MaxConcurrent:  h.MaxConcurrent(),
		CurrentTaskIds: ids,
		TasksCompleted: h.TasksCompleted(),
		TasksFailed:    h.TasksFailed(),
		UptimeSeconds:  h.UptimeSeconds(),
	}
}

func encodeTaskRequest(r *worker_pb.TaskRequest) []byte {
	return workerwire.NewTaskRequest(workerwire.TaskRequestInput{
		WorkerID:       r.WorkerId,
		Capabilities:   r.Capabilities,
		AvailableSlots: r.AvailableSlots,
	})
}

func decodeTaskRequest(r workerwire.TaskRequest) *worker_pb.TaskRequest {
	caps := make([]string, r.CapabilitiesLen())
	for i := range caps {
		caps[i] = r.Capability(i)
	}
	return &worker_pb.TaskRequest{
		WorkerId:       r.WorkerID(),
		Capabilities:   caps,
		AvailableSlots: r.AvailableSlots(),
	}
}

func encodeTaskUpdate(u *worker_pb.TaskUpdate) []byte {
	return workerwire.NewTaskUpdate(workerwire.TaskUpdateInput{
		TaskID:   u.TaskId,
		WorkerID: u.WorkerId,
		Status:   u.Status,
		Progress: u.Progress,
		Message:  u.Message,
		Metadata: mapToEntries(u.Metadata),
	})
}

func decodeTaskUpdate(u workerwire.TaskUpdate) *worker_pb.TaskUpdate {
	return &worker_pb.TaskUpdate{
		TaskId:   u.TaskID(),
		WorkerId: u.WorkerID(),
		Status:   u.Status(),
		Progress: u.Progress(),
		Message:  u.Message(),
		Metadata: entriesToMap(u.MetadataLen(), u.Metadata),
	}
}

func encodeTaskComplete(c *worker_pb.TaskComplete) []byte {
	return workerwire.NewTaskComplete(workerwire.TaskCompleteInput{
		TaskID:         c.TaskId,
		WorkerID:       c.WorkerId,
		Success:        c.Success,
		ErrorMessage:   c.ErrorMessage,
		CompletionTime: c.CompletionTime,
		ResultMetadata: mapToEntries(c.ResultMetadata),
	})
}

func decodeTaskComplete(c workerwire.TaskComplete) *worker_pb.TaskComplete {
	return &worker_pb.TaskComplete{
		TaskId:         c.TaskID(),
		WorkerId:       c.WorkerID(),
		Success:        c.Success(),
		ErrorMessage:   c.ErrorMessage(),
		CompletionTime: c.CompletionTime(),
		ResultMetadata: entriesToMap(c.ResultMetadataLen(), c.ResultMetadata),
	}
}

func encodeWorkerShutdown(s *worker_pb.WorkerShutdown) []byte {
	return workerwire.NewWorkerShutdown(workerwire.WorkerShutdownInput{
		WorkerID:       s.WorkerId,
		Reason:         s.Reason,
		PendingTaskIDs: s.PendingTaskIds,
	})
}

func decodeWorkerShutdown(s workerwire.WorkerShutdown) *worker_pb.WorkerShutdown {
	ids := make([]string, s.PendingTaskIDsLen())
	for i := range ids {
		ids[i] = s.PendingTaskID(i)
	}
	return &worker_pb.WorkerShutdown{
		WorkerId:       s.WorkerID(),
		Reason:         s.Reason(),
		PendingTaskIds: ids,
	}
}

func encodeTaskLogResponse(r *worker_pb.TaskLogResponse) []byte {
	in := workerwire.TaskLogResponseInput{
		TaskID:       r.TaskId,
		WorkerID:     r.WorkerId,
		Success:      r.Success,
		ErrorMessage: r.ErrorMessage,
	}
	if r.Metadata != nil {
		in.Metadata = encodeTaskLogMetadata(r.Metadata)
	}
	if len(r.LogEntries) > 0 {
		in.LogEntries = make([][]byte, len(r.LogEntries))
		for i, e := range r.LogEntries {
			in.LogEntries[i] = encodeTaskLogEntry(e)
		}
	}
	return workerwire.NewTaskLogResponse(in)
}

func decodeTaskLogResponse(r workerwire.TaskLogResponse) *worker_pb.TaskLogResponse {
	out := &worker_pb.TaskLogResponse{
		TaskId:       r.TaskID(),
		WorkerId:     r.WorkerID(),
		Success:      r.Success(),
		ErrorMessage: r.ErrorMessage(),
	}
	if r.HasMetadata() {
		out.Metadata = decodeTaskLogMetadata(r.Metadata())
	}
	if n := r.LogEntriesLen(); n > 0 {
		out.LogEntries = make([]*worker_pb.TaskLogEntry, n)
		for i := 0; i < n; i++ {
			out.LogEntries[i] = decodeTaskLogEntry(r.LogEntry(i))
		}
	}
	return out
}

func encodeTaskLogEntry(e *worker_pb.TaskLogEntry) []byte {
	return workerwire.NewTaskLogEntry(workerwire.TaskLogEntryInput{
		Timestamp: e.Timestamp,
		Level:     e.Level,
		Message:   e.Message,
		Fields:    mapToEntries(e.Fields),
		Progress:  e.Progress,
		Status:    e.Status,
	})
}

func decodeTaskLogEntry(e workerwire.TaskLogEntry) *worker_pb.TaskLogEntry {
	return &worker_pb.TaskLogEntry{
		Timestamp: e.Timestamp(),
		Level:     e.Level(),
		Message:   e.Message(),
		Fields:    entriesToMap(e.FieldsLen(), e.Field),
		Progress:  e.Progress(),
		Status:    e.Status(),
	}
}

func encodeTaskLogMetadata(m *worker_pb.TaskLogMetadata) []byte {
	return workerwire.NewTaskLogMetadata(workerwire.TaskLogMetadataInput{
		TaskID:      m.TaskId,
		TaskType:    m.TaskType,
		WorkerID:    m.WorkerId,
		StartTime:   m.StartTime,
		EndTime:     m.EndTime,
		DurationMs:  m.DurationMs,
		Status:      m.Status,
		Progress:    m.Progress,
		VolumeID:    m.VolumeId,
		Server:      m.Server,
		Collection:  m.Collection,
		LogFilePath: m.LogFilePath,
		CreatedAt:   m.CreatedAt,
		CustomData:  mapToEntries(m.CustomData),
	})
}

func decodeTaskLogMetadata(m workerwire.TaskLogMetadata) *worker_pb.TaskLogMetadata {
	return &worker_pb.TaskLogMetadata{
		TaskId:      m.TaskID(),
		TaskType:    m.TaskType(),
		WorkerId:    m.WorkerID(),
		StartTime:   m.StartTime(),
		EndTime:     m.EndTime(),
		DurationMs:  m.DurationMs(),
		Status:      m.Status(),
		Progress:    m.Progress(),
		VolumeId:    m.VolumeID(),
		Server:      m.Server(),
		Collection:  m.Collection(),
		LogFilePath: m.LogFilePath(),
		CreatedAt:   m.CreatedAt(),
		CustomData:  entriesToMap(m.CustomDataLen(), m.CustomData),
	}
}

// ---------------- AdminMessage variants ----------------

func encodeRegistrationResponse(r *worker_pb.RegistrationResponse) []byte {
	return workerwire.NewRegistrationResponse(workerwire.RegistrationResponseInput{
		Success:          r.Success,
		Message:          r.Message,
		AssignedWorkerID: r.AssignedWorkerId,
	})
}

func decodeRegistrationResponse(r workerwire.RegistrationResponse) *worker_pb.RegistrationResponse {
	return &worker_pb.RegistrationResponse{
		Success:          r.Success(),
		Message:          r.Message(),
		AssignedWorkerId: r.AssignedWorkerID(),
	}
}

func encodeHeartbeatResponse(r *worker_pb.HeartbeatResponse) []byte {
	return workerwire.NewHeartbeatResponse(workerwire.HeartbeatResponseInput{
		Success: r.Success,
		Message: r.Message,
	})
}

func decodeHeartbeatResponse(r workerwire.HeartbeatResponse) *worker_pb.HeartbeatResponse {
	return &worker_pb.HeartbeatResponse{
		Success: r.Success(),
		Message: r.Message(),
	}
}

func encodeTaskCancellation(c *worker_pb.TaskCancellation) []byte {
	return workerwire.NewTaskCancellation(workerwire.TaskCancellationInput{
		TaskID: c.TaskId,
		Reason: c.Reason,
		Force:  c.Force,
	})
}

func decodeTaskCancellation(c workerwire.TaskCancellation) *worker_pb.TaskCancellation {
	return &worker_pb.TaskCancellation{
		TaskId: c.TaskID(),
		Reason: c.Reason(),
		Force:  c.Force(),
	}
}

func encodeAdminShutdown(s *worker_pb.AdminShutdown) []byte {
	return workerwire.NewAdminShutdown(workerwire.AdminShutdownInput{
		Reason:                  s.Reason,
		GracefulShutdownSeconds: s.GracefulShutdownSeconds,
	})
}

func decodeAdminShutdown(s workerwire.AdminShutdown) *worker_pb.AdminShutdown {
	return &worker_pb.AdminShutdown{
		Reason:                  s.Reason(),
		GracefulShutdownSeconds: s.GracefulShutdownSeconds(),
	}
}

func encodeTaskLogRequest(r *worker_pb.TaskLogRequest) []byte {
	return workerwire.NewTaskLogRequest(workerwire.TaskLogRequestInput{
		TaskID:          r.TaskId,
		WorkerID:        r.WorkerId,
		IncludeMetadata: r.IncludeMetadata,
		MaxEntries:      r.MaxEntries,
		LogLevel:        r.LogLevel,
		StartTime:       r.StartTime,
		EndTime:         r.EndTime,
	})
}

func decodeTaskLogRequest(r workerwire.TaskLogRequest) *worker_pb.TaskLogRequest {
	return &worker_pb.TaskLogRequest{
		TaskId:          r.TaskID(),
		WorkerId:        r.WorkerID(),
		IncludeMetadata: r.IncludeMetadata(),
		MaxEntries:      r.MaxEntries(),
		LogLevel:        r.LogLevel(),
		StartTime:       r.StartTime(),
		EndTime:         r.EndTime(),
	}
}

// ---------------- TaskAssignment + TaskParams tree ----------------

func encodeTaskAssignment(a *worker_pb.TaskAssignment) []byte {
	in := workerwire.TaskAssignmentInput{
		TaskID:      a.TaskId,
		TaskType:    a.TaskType,
		Priority:    a.Priority,
		CreatedTime: a.CreatedTime,
		Metadata:    mapToEntries(a.Metadata),
	}
	if a.Params != nil {
		in.Params = encodeTaskParams(a.Params)
	}
	return workerwire.NewTaskAssignment(in)
}

func decodeTaskAssignment(a workerwire.TaskAssignment) *worker_pb.TaskAssignment {
	out := &worker_pb.TaskAssignment{
		TaskId:      a.TaskID(),
		TaskType:    a.TaskType(),
		Priority:    a.Priority(),
		CreatedTime: a.CreatedTime(),
		Metadata:    entriesToMap(a.MetadataLen(), a.Metadata),
	}
	if a.HasParams() {
		out.Params = decodeTaskParams(a.Params())
	}
	return out
}

func encodeTaskParams(p *worker_pb.TaskParams) []byte {
	in := workerwire.TaskParamsInput{
		TaskID:     p.TaskId,
		VolumeID:   p.VolumeId,
		Collection: p.Collection,
		DataCenter: p.DataCenter,
		Rack:       p.Rack,
		VolumeSize: p.VolumeSize,
	}
	if len(p.Sources) > 0 {
		in.Sources = make([][]byte, len(p.Sources))
		for i, s := range p.Sources {
			in.Sources[i] = encodeTaskSource(s)
		}
	}
	if len(p.Targets) > 0 {
		in.Targets = make([][]byte, len(p.Targets))
		for i, t := range p.Targets {
			in.Targets[i] = encodeTaskTarget(t)
		}
	}
	switch v := p.TaskParams.(type) {
	case *worker_pb.TaskParams_VacuumParams:
		in.TaskParamsWhich = workerwire.TaskParamsTaskParamsVacuumParams
		in.TaskParamsValue = encodeVacuumParams(v.VacuumParams)
	case *worker_pb.TaskParams_ErasureCodingParams:
		in.TaskParamsWhich = workerwire.TaskParamsTaskParamsErasureCodingParams
		in.TaskParamsValue = encodeErasureCodingParams(v.ErasureCodingParams)
	case *worker_pb.TaskParams_BalanceParams:
		in.TaskParamsWhich = workerwire.TaskParamsTaskParamsBalanceParams
		in.TaskParamsValue = encodeBalanceParams(v.BalanceParams)
	case *worker_pb.TaskParams_ReplicationParams:
		in.TaskParamsWhich = workerwire.TaskParamsTaskParamsReplicationParams
		in.TaskParamsValue = encodeReplicationParams(v.ReplicationParams)
	case *worker_pb.TaskParams_EcBalanceParams:
		in.TaskParamsWhich = workerwire.TaskParamsTaskParamsEcBalanceParams
		in.TaskParamsValue = encodeEcBalanceParams(v.EcBalanceParams)
	case *worker_pb.TaskParams_S3LifecycleParams:
		in.TaskParamsWhich = workerwire.TaskParamsTaskParamsS3LifecycleParams
		in.TaskParamsValue = encodeS3LifecycleParams(v.S3LifecycleParams)
	}
	return workerwire.NewTaskParams(in)
}

func decodeTaskParams(p workerwire.TaskParams) *worker_pb.TaskParams {
	out := &worker_pb.TaskParams{
		TaskId:     p.TaskID(),
		VolumeId:   p.VolumeID(),
		Collection: p.Collection(),
		DataCenter: p.DataCenter(),
		Rack:       p.Rack(),
		VolumeSize: p.VolumeSize(),
	}
	if n := p.SourcesLen(); n > 0 {
		out.Sources = make([]*worker_pb.TaskSource, n)
		for i := 0; i < n; i++ {
			out.Sources[i] = decodeTaskSource(p.Sources(i))
		}
	}
	if n := p.TargetsLen(); n > 0 {
		out.Targets = make([]*worker_pb.TaskTarget, n)
		for i := 0; i < n; i++ {
			out.Targets[i] = decodeTaskTarget(p.Targets(i))
		}
	}
	switch p.WhichTaskParams() {
	case workerwire.TaskParamsTaskParamsVacuumParams:
		out.TaskParams = &worker_pb.TaskParams_VacuumParams{VacuumParams: decodeVacuumParams(p.VacuumParams())}
	case workerwire.TaskParamsTaskParamsErasureCodingParams:
		out.TaskParams = &worker_pb.TaskParams_ErasureCodingParams{ErasureCodingParams: decodeErasureCodingParams(p.ErasureCodingParams())}
	case workerwire.TaskParamsTaskParamsBalanceParams:
		out.TaskParams = &worker_pb.TaskParams_BalanceParams{BalanceParams: decodeBalanceParams(p.BalanceParams())}
	case workerwire.TaskParamsTaskParamsReplicationParams:
		out.TaskParams = &worker_pb.TaskParams_ReplicationParams{ReplicationParams: decodeReplicationParams(p.ReplicationParams())}
	case workerwire.TaskParamsTaskParamsEcBalanceParams:
		out.TaskParams = &worker_pb.TaskParams_EcBalanceParams{EcBalanceParams: decodeEcBalanceParams(p.EcBalanceParams())}
	case workerwire.TaskParamsTaskParamsS3LifecycleParams:
		out.TaskParams = &worker_pb.TaskParams_S3LifecycleParams{S3LifecycleParams: decodeS3LifecycleParams(p.S3LifecycleParams())}
	}
	return out
}

func encodeTaskSource(s *worker_pb.TaskSource) []byte {
	return workerwire.NewTaskSource(workerwire.TaskSourceInput{
		Node:          s.Node,
		DiskID:        s.DiskId,
		Rack:          s.Rack,
		DataCenter:    s.DataCenter,
		VolumeID:      s.VolumeId,
		ShardIDs:      s.ShardIds,
		EstimatedSize: s.EstimatedSize,
	})
}

func decodeTaskSource(s workerwire.TaskSource) *worker_pb.TaskSource {
	ids := make([]uint32, s.ShardIDsLen())
	for i := range ids {
		ids[i] = s.ShardID(i)
	}
	return &worker_pb.TaskSource{
		Node:          s.Node(),
		DiskId:        s.DiskID(),
		Rack:          s.Rack(),
		DataCenter:    s.DataCenter(),
		VolumeId:      s.VolumeID(),
		ShardIds:      ids,
		EstimatedSize: s.EstimatedSize(),
	}
}

func encodeTaskTarget(t *worker_pb.TaskTarget) []byte {
	return workerwire.NewTaskTarget(workerwire.TaskTargetInput{
		Node:          t.Node,
		DiskID:        t.DiskId,
		Rack:          t.Rack,
		DataCenter:    t.DataCenter,
		VolumeID:      t.VolumeId,
		ShardIDs:      t.ShardIds,
		EstimatedSize: t.EstimatedSize,
	})
}

func decodeTaskTarget(t workerwire.TaskTarget) *worker_pb.TaskTarget {
	ids := make([]uint32, t.ShardIDsLen())
	for i := range ids {
		ids[i] = t.ShardID(i)
	}
	return &worker_pb.TaskTarget{
		Node:          t.Node(),
		DiskId:        t.DiskID(),
		Rack:          t.Rack(),
		DataCenter:    t.DataCenter(),
		VolumeId:      t.VolumeID(),
		ShardIds:      ids,
		EstimatedSize: t.EstimatedSize(),
	}
}

func encodeVacuumParams(v *worker_pb.VacuumTaskParams) []byte {
	return workerwire.NewVacuumTaskParams(workerwire.VacuumTaskParamsInput{
		GarbageThreshold: v.GarbageThreshold,
		ForceVacuum:      v.ForceVacuum,
		BatchSize:        v.BatchSize,
		WorkingDir:       v.WorkingDir,
		VerifyChecksum:   v.VerifyChecksum,
	})
}

func decodeVacuumParams(v workerwire.VacuumTaskParams) *worker_pb.VacuumTaskParams {
	return &worker_pb.VacuumTaskParams{
		GarbageThreshold: v.GarbageThreshold(),
		ForceVacuum:      v.ForceVacuum(),
		BatchSize:        v.BatchSize(),
		WorkingDir:       v.WorkingDir(),
		VerifyChecksum:   v.VerifyChecksum(),
	}
}

func encodeErasureCodingParams(e *worker_pb.ErasureCodingTaskParams) []byte {
	return workerwire.NewErasureCodingTaskParams(workerwire.ErasureCodingTaskParamsInput{
		EstimatedShardSize: e.EstimatedShardSize,
		DataShards:         e.DataShards,
		ParityShards:       e.ParityShards,
		WorkingDir:         e.WorkingDir,
		MasterClient:       e.MasterClient,
		CleanupSource:      e.CleanupSource,
		SourceDiskType:     e.SourceDiskType,
		EncodeTsNs:         e.EncodeTsNs,
	})
}

func decodeErasureCodingParams(e workerwire.ErasureCodingTaskParams) *worker_pb.ErasureCodingTaskParams {
	return &worker_pb.ErasureCodingTaskParams{
		EstimatedShardSize: e.EstimatedShardSize(),
		DataShards:         e.DataShards(),
		ParityShards:       e.ParityShards(),
		WorkingDir:         e.WorkingDir(),
		MasterClient:       e.MasterClient(),
		CleanupSource:      e.CleanupSource(),
		SourceDiskType:     e.SourceDiskType(),
		EncodeTsNs:         e.EncodeTsNs(),
	}
}

func encodeBalanceParams(b *worker_pb.BalanceTaskParams) []byte {
	in := workerwire.BalanceTaskParamsInput{
		ForceMove:          b.ForceMove,
		TimeoutSeconds:     b.TimeoutSeconds,
		MaxConcurrentMoves: b.MaxConcurrentMoves,
	}
	if len(b.Moves) > 0 {
		in.Moves = make([][]byte, len(b.Moves))
		for i, m := range b.Moves {
			in.Moves[i] = encodeBalanceMoveSpec(m)
		}
	}
	return workerwire.NewBalanceTaskParams(in)
}

func decodeBalanceParams(b workerwire.BalanceTaskParams) *worker_pb.BalanceTaskParams {
	out := &worker_pb.BalanceTaskParams{
		ForceMove:          b.ForceMove(),
		TimeoutSeconds:     b.TimeoutSeconds(),
		MaxConcurrentMoves: b.MaxConcurrentMoves(),
	}
	if n := b.MovesLen(); n > 0 {
		out.Moves = make([]*worker_pb.BalanceMoveSpec, n)
		for i := 0; i < n; i++ {
			out.Moves[i] = decodeBalanceMoveSpec(b.Moves(i))
		}
	}
	return out
}

func encodeBalanceMoveSpec(m *worker_pb.BalanceMoveSpec) []byte {
	return workerwire.NewBalanceMoveSpec(workerwire.BalanceMoveSpecInput{
		VolumeID:   m.VolumeId,
		SourceNode: m.SourceNode,
		TargetNode: m.TargetNode,
		Collection: m.Collection,
		VolumeSize: m.VolumeSize,
	})
}

func decodeBalanceMoveSpec(m workerwire.BalanceMoveSpec) *worker_pb.BalanceMoveSpec {
	return &worker_pb.BalanceMoveSpec{
		VolumeId:   m.VolumeID(),
		SourceNode: m.SourceNode(),
		TargetNode: m.TargetNode(),
		Collection: m.Collection(),
		VolumeSize: m.VolumeSize(),
	}
}

func encodeReplicationParams(r *worker_pb.ReplicationTaskParams) []byte {
	return workerwire.NewReplicationTaskParams(workerwire.ReplicationTaskParamsInput{
		ReplicaCount:      r.ReplicaCount,
		VerifyConsistency: r.VerifyConsistency,
	})
}

func decodeReplicationParams(r workerwire.ReplicationTaskParams) *worker_pb.ReplicationTaskParams {
	return &worker_pb.ReplicationTaskParams{
		ReplicaCount:      r.ReplicaCount(),
		VerifyConsistency: r.VerifyConsistency(),
	}
}

func encodeEcBalanceParams(e *worker_pb.EcBalanceTaskParams) []byte {
	in := workerwire.EcBalanceTaskParamsInput{
		DiskType:           e.DiskType,
		MaxParallelization: e.MaxParallelization,
		TimeoutSeconds:     e.TimeoutSeconds,
	}
	if len(e.Moves) > 0 {
		in.Moves = make([][]byte, len(e.Moves))
		for i, m := range e.Moves {
			in.Moves[i] = encodeEcShardMoveSpec(m)
		}
	}
	return workerwire.NewEcBalanceTaskParams(in)
}

func decodeEcBalanceParams(e workerwire.EcBalanceTaskParams) *worker_pb.EcBalanceTaskParams {
	out := &worker_pb.EcBalanceTaskParams{
		DiskType:           e.DiskType(),
		MaxParallelization: e.MaxParallelization(),
		TimeoutSeconds:     e.TimeoutSeconds(),
	}
	if n := e.MovesLen(); n > 0 {
		out.Moves = make([]*worker_pb.EcShardMoveSpec, n)
		for i := 0; i < n; i++ {
			out.Moves[i] = decodeEcShardMoveSpec(e.Moves(i))
		}
	}
	return out
}

func encodeEcShardMoveSpec(m *worker_pb.EcShardMoveSpec) []byte {
	return workerwire.NewEcShardMoveSpec(workerwire.EcShardMoveSpecInput{
		VolumeID:     m.VolumeId,
		ShardID:      m.ShardId,
		Collection:   m.Collection,
		SourceNode:   m.SourceNode,
		SourceDiskID: m.SourceDiskId,
		TargetNode:   m.TargetNode,
		TargetDiskID: m.TargetDiskId,
	})
}

func decodeEcShardMoveSpec(m workerwire.EcShardMoveSpec) *worker_pb.EcShardMoveSpec {
	return &worker_pb.EcShardMoveSpec{
		VolumeId:     m.VolumeID(),
		ShardId:      m.ShardID(),
		Collection:   m.Collection(),
		SourceNode:   m.SourceNode(),
		SourceDiskId: m.SourceDiskID(),
		TargetNode:   m.TargetNode(),
		TargetDiskId: m.TargetDiskID(),
	}
}

func encodeS3LifecycleParams(s *worker_pb.S3LifecycleParams) []byte {
	in := workerwire.S3LifecycleParamsInput{
		Subtype:           uint32(s.Subtype),
		Bucket:            s.Bucket,
		RuleHash:          s.RuleHash,
		Force:             s.Force,
		BatchTimeBudgetNs: s.BatchTimeBudgetNs,
		BatchEventBudget:  s.BatchEventBudget,
		ShardID:           s.ShardId,
	}
	if s.Continuation != nil {
		in.Continuation = encodeContinuationHint(s.Continuation)
	}
	return workerwire.NewS3LifecycleParams(in)
}

func decodeS3LifecycleParams(s workerwire.S3LifecycleParams) *worker_pb.S3LifecycleParams {
	out := &worker_pb.S3LifecycleParams{
		Subtype:           worker_pb.S3LifecycleParams_Subtype(s.Subtype()),
		Bucket:            s.Bucket(),
		RuleHash:          s.RuleHash(),
		Force:             s.Force(),
		BatchTimeBudgetNs: s.BatchTimeBudgetNs(),
		BatchEventBudget:  s.BatchEventBudget(),
		ShardId:           s.ShardID(),
	}
	if s.HasContinuation() {
		out.Continuation = decodeContinuationHint(s.Continuation())
	}
	return out
}

func encodeContinuationHint(c *worker_pb.ContinuationHint) []byte {
	return workerwire.NewContinuationHint(workerwire.ContinuationHintInput{
		LastScannedPath: c.LastScannedPath,
		LastPositionNs:  c.LastPositionNs,
	})
}

func decodeContinuationHint(c workerwire.ContinuationHint) *worker_pb.ContinuationHint {
	return &worker_pb.ContinuationHint{
		LastScannedPath: c.LastScannedPath(),
		LastPositionNs:  c.LastPositionNs(),
	}
}
