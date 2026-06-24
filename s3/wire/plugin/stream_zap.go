// Code generated from plugin.proto; DO NOT EDIT.

package pluginwire

import (
	zap "github.com/zap-proto/go"
)

// This file carries the two bidi-stream envelope messages, each with a `body`
// oneof. The oneof is a uint32 "which" tag (WhichBody; 0 = none) plus one nested
// object slot per member. Set Which on the Input and supply only that member's
// pre-built sub-buffer. google.protobuf.Timestamp sent_at is int64 unix-nanos.

// --- WorkerToAdminMessage ---

const (
	workerToAdminMessageWorkerIDOff           = 0
	workerToAdminMessageSentAtOff             = 8
	workerToAdminMessageWhichOff              = 16
	workerToAdminMessageHelloOff              = 20
	workerToAdminMessageHeartbeatOff          = 28
	workerToAdminMessageAcknowledgeOff        = 36
	workerToAdminMessageConfigSchemaRespOff   = 44
	workerToAdminMessageDetectionProposalsOff = 52
	workerToAdminMessageDetectionCompleteOff  = 60
	workerToAdminMessageJobProgressUpdateOff  = 68
	workerToAdminMessageJobCompletedOff       = 76
	workerToAdminMessageSize                  = 84
)

// WorkerToAdminMessage is a zero-copy view into a ZAP-encoded
// WorkerToAdminMessage. The body oneof is read via WhichBody plus the matching
// member accessor.
type WorkerToAdminMessage struct{ o zap.Object }

// WrapWorkerToAdminMessage parses b and returns a typed view.
func WrapWorkerToAdminMessage(b []byte) (WorkerToAdminMessage, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return WorkerToAdminMessage{}, err
	}
	return WorkerToAdminMessage{o: m.Root()}, nil
}

func (t WorkerToAdminMessage) WorkerID() string { return t.o.Text(workerToAdminMessageWorkerIDOff) }

// SentAt reads sent_at (proto field 2, google.protobuf.Timestamp) as int64
// unix-nanoseconds.
func (t WorkerToAdminMessage) SentAt() int64 { return t.o.Int64(workerToAdminMessageSentAtOff) }

// WhichBody reports which body oneof member is set (one of the
// WorkerToAdminBody* consts; 0 = none).
func (t WorkerToAdminMessage) WhichBody() uint32 { return t.o.Uint32(workerToAdminMessageWhichOff) }

func (t WorkerToAdminMessage) Hello() WorkerHello {
	return WorkerHello{o: childObject(t.o.Bytes(workerToAdminMessageHelloOff))}
}
func (t WorkerToAdminMessage) Heartbeat() WorkerHeartbeat {
	return WorkerHeartbeat{o: childObject(t.o.Bytes(workerToAdminMessageHeartbeatOff))}
}
func (t WorkerToAdminMessage) Acknowledge() WorkerAcknowledge {
	return WorkerAcknowledge{o: childObject(t.o.Bytes(workerToAdminMessageAcknowledgeOff))}
}
func (t WorkerToAdminMessage) ConfigSchemaResponse() ConfigSchemaResponse {
	return ConfigSchemaResponse{o: childObject(t.o.Bytes(workerToAdminMessageConfigSchemaRespOff))}
}
func (t WorkerToAdminMessage) DetectionProposals() DetectionProposals {
	return DetectionProposals{o: childObject(t.o.Bytes(workerToAdminMessageDetectionProposalsOff))}
}
func (t WorkerToAdminMessage) DetectionComplete() DetectionComplete {
	return DetectionComplete{o: childObject(t.o.Bytes(workerToAdminMessageDetectionCompleteOff))}
}
func (t WorkerToAdminMessage) JobProgressUpdate() JobProgressUpdate {
	return JobProgressUpdate{o: childObject(t.o.Bytes(workerToAdminMessageJobProgressUpdateOff))}
}
func (t WorkerToAdminMessage) JobCompleted() JobCompleted {
	return JobCompleted{o: childObject(t.o.Bytes(workerToAdminMessageJobCompletedOff))}
}

// WorkerToAdminMessageInput collects the field values for
// NewWorkerToAdminMessage. Set Which to the selected body member and supply only
// that member's pre-built sub-buffer.
type WorkerToAdminMessageInput struct {
	WorkerID             string
	SentAt               int64
	Which                uint32
	Hello                []byte
	Heartbeat            []byte
	Acknowledge          []byte
	ConfigSchemaResponse []byte
	DetectionProposals   []byte
	DetectionComplete    []byte
	JobProgressUpdate    []byte
	JobCompleted         []byte
}

// NewWorkerToAdminMessage builds a ZAP-encoded WorkerToAdminMessage. Only the
// body member selected by in.Which is written.
func NewWorkerToAdminMessage(in WorkerToAdminMessageInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(workerToAdminMessageSize)
	ob.SetText(workerToAdminMessageWorkerIDOff, in.WorkerID)
	ob.SetInt64(workerToAdminMessageSentAtOff, in.SentAt)
	ob.SetUint32(workerToAdminMessageWhichOff, in.Which)
	switch in.Which {
	case WorkerToAdminBodyHello:
		setChild(ob, workerToAdminMessageHelloOff, in.Hello)
	case WorkerToAdminBodyHeartbeat:
		setChild(ob, workerToAdminMessageHeartbeatOff, in.Heartbeat)
	case WorkerToAdminBodyAcknowledge:
		setChild(ob, workerToAdminMessageAcknowledgeOff, in.Acknowledge)
	case WorkerToAdminBodyConfigSchemaResp:
		setChild(ob, workerToAdminMessageConfigSchemaRespOff, in.ConfigSchemaResponse)
	case WorkerToAdminBodyDetectionProposals:
		setChild(ob, workerToAdminMessageDetectionProposalsOff, in.DetectionProposals)
	case WorkerToAdminBodyDetectionComplete:
		setChild(ob, workerToAdminMessageDetectionCompleteOff, in.DetectionComplete)
	case WorkerToAdminBodyJobProgressUpdate:
		setChild(ob, workerToAdminMessageJobProgressUpdateOff, in.JobProgressUpdate)
	case WorkerToAdminBodyJobCompleted:
		setChild(ob, workerToAdminMessageJobCompletedOff, in.JobCompleted)
	}
	ob.FinishAsRoot()
	return b.Finish()
}

// --- AdminToWorkerMessage ---

const (
	adminToWorkerMessageRequestIDOff           = 0
	adminToWorkerMessageSentAtOff              = 8
	adminToWorkerMessageWhichOff               = 16
	adminToWorkerMessageHelloOff               = 20
	adminToWorkerMessageRequestConfigSchemaOff = 28
	adminToWorkerMessageRunDetectionRequestOff = 36
	adminToWorkerMessageExecuteJobRequestOff   = 44
	adminToWorkerMessageCancelRequestOff       = 52
	adminToWorkerMessageShutdownOff            = 60
	adminToWorkerMessageSize                   = 68
)

// AdminToWorkerMessage is a zero-copy view into a ZAP-encoded
// AdminToWorkerMessage. The body oneof is read via WhichBody plus the matching
// member accessor.
type AdminToWorkerMessage struct{ o zap.Object }

// WrapAdminToWorkerMessage parses b and returns a typed view.
func WrapAdminToWorkerMessage(b []byte) (AdminToWorkerMessage, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return AdminToWorkerMessage{}, err
	}
	return AdminToWorkerMessage{o: m.Root()}, nil
}

func (t AdminToWorkerMessage) RequestID() string { return t.o.Text(adminToWorkerMessageRequestIDOff) }

// SentAt reads sent_at (proto field 2, google.protobuf.Timestamp) as int64
// unix-nanoseconds.
func (t AdminToWorkerMessage) SentAt() int64 { return t.o.Int64(adminToWorkerMessageSentAtOff) }

// WhichBody reports which body oneof member is set (one of the
// AdminToWorkerBody* consts; 0 = none).
func (t AdminToWorkerMessage) WhichBody() uint32 { return t.o.Uint32(adminToWorkerMessageWhichOff) }

func (t AdminToWorkerMessage) Hello() AdminHello {
	return AdminHello{o: childObject(t.o.Bytes(adminToWorkerMessageHelloOff))}
}
func (t AdminToWorkerMessage) RequestConfigSchema() RequestConfigSchema {
	return RequestConfigSchema{o: childObject(t.o.Bytes(adminToWorkerMessageRequestConfigSchemaOff))}
}
func (t AdminToWorkerMessage) RunDetectionRequest() RunDetectionRequest {
	return RunDetectionRequest{o: childObject(t.o.Bytes(adminToWorkerMessageRunDetectionRequestOff))}
}
func (t AdminToWorkerMessage) ExecuteJobRequest() ExecuteJobRequest {
	return ExecuteJobRequest{o: childObject(t.o.Bytes(adminToWorkerMessageExecuteJobRequestOff))}
}
func (t AdminToWorkerMessage) CancelRequest() CancelRequest {
	return CancelRequest{o: childObject(t.o.Bytes(adminToWorkerMessageCancelRequestOff))}
}
func (t AdminToWorkerMessage) Shutdown() AdminShutdown {
	return AdminShutdown{o: childObject(t.o.Bytes(adminToWorkerMessageShutdownOff))}
}

// AdminToWorkerMessageInput collects the field values for
// NewAdminToWorkerMessage. Set Which to the selected body member and supply only
// that member's pre-built sub-buffer.
type AdminToWorkerMessageInput struct {
	RequestID           string
	SentAt              int64
	Which               uint32
	Hello               []byte
	RequestConfigSchema []byte
	RunDetectionRequest []byte
	ExecuteJobRequest   []byte
	CancelRequest       []byte
	Shutdown            []byte
}

// NewAdminToWorkerMessage builds a ZAP-encoded AdminToWorkerMessage. Only the
// body member selected by in.Which is written.
func NewAdminToWorkerMessage(in AdminToWorkerMessageInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(adminToWorkerMessageSize)
	ob.SetText(adminToWorkerMessageRequestIDOff, in.RequestID)
	ob.SetInt64(adminToWorkerMessageSentAtOff, in.SentAt)
	ob.SetUint32(adminToWorkerMessageWhichOff, in.Which)
	switch in.Which {
	case AdminToWorkerBodyHello:
		setChild(ob, adminToWorkerMessageHelloOff, in.Hello)
	case AdminToWorkerBodyRequestConfigSchema:
		setChild(ob, adminToWorkerMessageRequestConfigSchemaOff, in.RequestConfigSchema)
	case AdminToWorkerBodyRunDetectionRequest:
		setChild(ob, adminToWorkerMessageRunDetectionRequestOff, in.RunDetectionRequest)
	case AdminToWorkerBodyExecuteJobRequest:
		setChild(ob, adminToWorkerMessageExecuteJobRequestOff, in.ExecuteJobRequest)
	case AdminToWorkerBodyCancelRequest:
		setChild(ob, adminToWorkerMessageCancelRequestOff, in.CancelRequest)
	case AdminToWorkerBodyShutdown:
		setChild(ob, adminToWorkerMessageShutdownOff, in.Shutdown)
	}
	ob.FinishAsRoot()
	return b.Finish()
}
