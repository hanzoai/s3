// Code generated from worker.proto; DO NOT EDIT.

package workerwire

import (
	zap "github.com/zap-proto/go"
)

// WorkerMessage.message oneof tags (proto field numbers, stable wire ids).
// Zero means no variant is set. The selected variant message is embedded as a
// nested object pointed to by WorkerMessage.message_value.
const (
	WorkerMessageMessageNone            uint32 = 0
	WorkerMessageMessageRegistration    uint32 = 3
	WorkerMessageMessageHeartbeat       uint32 = 4
	WorkerMessageMessageTaskRequest     uint32 = 5
	WorkerMessageMessageTaskUpdate      uint32 = 6
	WorkerMessageMessageTaskComplete    uint32 = 7
	WorkerMessageMessageShutdown        uint32 = 8
	WorkerMessageMessageTaskLogResponse uint32 = 9
)

// --- WorkerMessage ---
//
// Fixed payload: worker_id(8 ptr) timestamp(8) message_which(4 oneof tag)
// message_value(4 obj ptr). The oneof value is the embedded selected variant
// message (nested object); see WorkerMessageMessage* tags.

const (
	workerMessageWorkerIDOff     = 0
	workerMessageTimestampOff    = 8
	workerMessageMessageWhichOff = 16
	workerMessageMessageValueOff = 20
	workerMessageSize            = 24
)

// WorkerMessage is a zero-copy view into a ZAP-encoded WorkerMessage message.
type WorkerMessage struct{ o zap.Object }

// WrapWorkerMessage parses b and returns a typed view.
func WrapWorkerMessage(b []byte) (WorkerMessage, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return WorkerMessage{}, err
	}
	return WorkerMessage{o: m.Root()}, nil
}

func (t WorkerMessage) WorkerID() string { return t.o.Text(workerMessageWorkerIDOff) }
func (t WorkerMessage) Timestamp() int64 { return t.o.Int64(workerMessageTimestampOff) }

// WhichMessage returns the set oneof tag (WorkerMessageMessage* constants).
func (t WorkerMessage) WhichMessage() uint32 { return t.o.Uint32(workerMessageMessageWhichOff) }

// Registration returns the embedded registration variant.
func (t WorkerMessage) Registration() WorkerRegistration {
	return WorkerRegistration{o: t.o.Object(workerMessageMessageValueOff)}
}

// Heartbeat returns the embedded heartbeat variant.
func (t WorkerMessage) Heartbeat() WorkerHeartbeat {
	return WorkerHeartbeat{o: t.o.Object(workerMessageMessageValueOff)}
}

// TaskRequest returns the embedded task_request variant.
func (t WorkerMessage) TaskRequest() TaskRequest {
	return TaskRequest{o: t.o.Object(workerMessageMessageValueOff)}
}

// TaskUpdate returns the embedded task_update variant.
func (t WorkerMessage) TaskUpdate() TaskUpdate {
	return TaskUpdate{o: t.o.Object(workerMessageMessageValueOff)}
}

// TaskComplete returns the embedded task_complete variant.
func (t WorkerMessage) TaskComplete() TaskComplete {
	return TaskComplete{o: t.o.Object(workerMessageMessageValueOff)}
}

// Shutdown returns the embedded shutdown variant.
func (t WorkerMessage) Shutdown() WorkerShutdown {
	return WorkerShutdown{o: t.o.Object(workerMessageMessageValueOff)}
}

// TaskLogResponse returns the embedded task_log_response variant.
func (t WorkerMessage) TaskLogResponse() TaskLogResponse {
	return TaskLogResponse{o: t.o.Object(workerMessageMessageValueOff)}
}

// WorkerMessageInput collects the field values for NewWorkerMessage. MessageWhich
// selects the oneof variant; MessageValue is that variant's standalone buffer.
type WorkerMessageInput struct {
	WorkerID     string
	Timestamp    int64
	MessageWhich uint32
	MessageValue []byte
}

// NewWorkerMessage builds a ZAP-encoded WorkerMessage message from in.
func NewWorkerMessage(in WorkerMessageInput) []byte {
	b := zap.NewBuilder(512)
	valueOff := 0
	if len(in.MessageValue) > 0 {
		valueOff = embedMessage(b, in.MessageValue)
	}
	ob := b.StartObject(workerMessageSize)
	ob.SetText(workerMessageWorkerIDOff, in.WorkerID)
	ob.SetInt64(workerMessageTimestampOff, in.Timestamp)
	ob.SetUint32(workerMessageMessageWhichOff, in.MessageWhich)
	ob.SetObject(workerMessageMessageValueOff, valueOff)
	ob.FinishAsRoot()
	return b.Finish()
}

// AdminMessage.message oneof tags (proto field numbers, stable wire ids).
// Zero means no variant is set.
const (
	AdminMessageMessageNone                 uint32 = 0
	AdminMessageMessageRegistrationResponse uint32 = 3
	AdminMessageMessageHeartbeatResponse    uint32 = 4
	AdminMessageMessageTaskAssignment       uint32 = 5
	AdminMessageMessageTaskCancellation     uint32 = 6
	AdminMessageMessageAdminShutdown        uint32 = 7
	AdminMessageMessageTaskLogRequest       uint32 = 8
)

// --- AdminMessage ---
//
// Fixed payload: admin_id(8 ptr) timestamp(8) message_which(4 oneof tag)
// message_value(4 obj ptr). The oneof value is the embedded selected variant
// message (nested object); see AdminMessageMessage* tags.

const (
	adminMessageAdminIDOff      = 0
	adminMessageTimestampOff    = 8
	adminMessageMessageWhichOff = 16
	adminMessageMessageValueOff = 20
	adminMessageSize            = 24
)

// AdminMessage is a zero-copy view into a ZAP-encoded AdminMessage message.
type AdminMessage struct{ o zap.Object }

// WrapAdminMessage parses b and returns a typed view.
func WrapAdminMessage(b []byte) (AdminMessage, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return AdminMessage{}, err
	}
	return AdminMessage{o: m.Root()}, nil
}

func (t AdminMessage) AdminID() string  { return t.o.Text(adminMessageAdminIDOff) }
func (t AdminMessage) Timestamp() int64 { return t.o.Int64(adminMessageTimestampOff) }

// WhichMessage returns the set oneof tag (AdminMessageMessage* constants).
func (t AdminMessage) WhichMessage() uint32 { return t.o.Uint32(adminMessageMessageWhichOff) }

// RegistrationResponse returns the embedded registration_response variant.
func (t AdminMessage) RegistrationResponse() RegistrationResponse {
	return RegistrationResponse{o: t.o.Object(adminMessageMessageValueOff)}
}

// HeartbeatResponse returns the embedded heartbeat_response variant.
func (t AdminMessage) HeartbeatResponse() HeartbeatResponse {
	return HeartbeatResponse{o: t.o.Object(adminMessageMessageValueOff)}
}

// TaskAssignment returns the embedded task_assignment variant.
func (t AdminMessage) TaskAssignment() TaskAssignment {
	return TaskAssignment{o: t.o.Object(adminMessageMessageValueOff)}
}

// TaskCancellation returns the embedded task_cancellation variant.
func (t AdminMessage) TaskCancellation() TaskCancellation {
	return TaskCancellation{o: t.o.Object(adminMessageMessageValueOff)}
}

// AdminShutdown returns the embedded admin_shutdown variant.
func (t AdminMessage) AdminShutdown() AdminShutdown {
	return AdminShutdown{o: t.o.Object(adminMessageMessageValueOff)}
}

// TaskLogRequest returns the embedded task_log_request variant.
func (t AdminMessage) TaskLogRequest() TaskLogRequest {
	return TaskLogRequest{o: t.o.Object(adminMessageMessageValueOff)}
}

// AdminMessageInput collects the field values for NewAdminMessage. MessageWhich
// selects the oneof variant; MessageValue is that variant's standalone buffer.
type AdminMessageInput struct {
	AdminID      string
	Timestamp    int64
	MessageWhich uint32
	MessageValue []byte
}

// NewAdminMessage builds a ZAP-encoded AdminMessage message from in.
func NewAdminMessage(in AdminMessageInput) []byte {
	b := zap.NewBuilder(512)
	valueOff := 0
	if len(in.MessageValue) > 0 {
		valueOff = embedMessage(b, in.MessageValue)
	}
	ob := b.StartObject(adminMessageSize)
	ob.SetText(adminMessageAdminIDOff, in.AdminID)
	ob.SetInt64(adminMessageTimestampOff, in.Timestamp)
	ob.SetUint32(adminMessageMessageWhichOff, in.MessageWhich)
	ob.SetObject(adminMessageMessageValueOff, valueOff)
	ob.FinishAsRoot()
	return b.Finish()
}
