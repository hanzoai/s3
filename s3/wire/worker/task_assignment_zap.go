// Code generated from worker.proto; DO NOT EDIT.

package workerwire

import (
	zap "github.com/zap-proto/go"
)

// --- TaskAssignment ---
//
// Fixed payload: task_id(8 ptr) task_type(8 ptr) params(4 obj ptr) priority(4)
// created_time(8) metadata(8 list ptr). params is a singular TaskParams nested
// object; metadata is map<string,string>.

const (
	taskAssignmentTaskIDOff      = 0
	taskAssignmentTaskTypeOff    = 8
	taskAssignmentParamsOff      = 16
	taskAssignmentPriorityOff    = 20
	taskAssignmentCreatedTimeOff = 24
	taskAssignmentMetadataOff    = 32
	taskAssignmentSize           = 40
)

// TaskAssignment is a zero-copy view into a ZAP-encoded TaskAssignment message.
type TaskAssignment struct{ o zap.Object }

// WrapTaskAssignment parses b and returns a typed view.
func WrapTaskAssignment(b []byte) (TaskAssignment, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return TaskAssignment{}, err
	}
	return TaskAssignment{o: m.Root()}, nil
}

func (t TaskAssignment) TaskID() string   { return t.o.Text(taskAssignmentTaskIDOff) }
func (t TaskAssignment) TaskType() string { return t.o.Text(taskAssignmentTaskTypeOff) }
func (t TaskAssignment) Priority() int32  { return t.o.Int32(taskAssignmentPriorityOff) }
func (t TaskAssignment) CreatedTime() int64 {
	return t.o.Int64(taskAssignmentCreatedTimeOff)
}

// Params returns the nested TaskParams (null-object if unset).
func (t TaskAssignment) Params() TaskParams {
	return TaskParams{o: t.o.Object(taskAssignmentParamsOff)}
}

// HasParams reports whether the params field is present.
func (t TaskAssignment) HasParams() bool {
	return !t.o.Object(taskAssignmentParamsOff).IsNull()
}

// MetadataLen returns the number of metadata entries.
func (t TaskAssignment) MetadataLen() int { return t.o.List(taskAssignmentMetadataOff).Len() }

// Metadata returns the i-th metadata entry (key/value).
func (t TaskAssignment) Metadata(i int) MapEntry {
	return mapEntryAt(t.o.List(taskAssignmentMetadataOff), i)
}

// TaskAssignmentInput collects the field values for NewTaskAssignment. Params,
// when non-nil, is a standalone TaskParams buffer (NewTaskParams).
type TaskAssignmentInput struct {
	TaskID      string
	TaskType    string
	Params      []byte
	Priority    int32
	CreatedTime int64
	Metadata    []MapEntry
}

// NewTaskAssignment builds a ZAP-encoded TaskAssignment message from in.
func NewTaskAssignment(in TaskAssignmentInput) []byte {
	b := zap.NewBuilder(512)
	paramsOff := 0
	if len(in.Params) > 0 {
		paramsOff = embedMessage(b, in.Params)
	}
	metaOff, metaLen := 0, 0
	if len(in.Metadata) > 0 {
		el := startElemList(b)
		for _, e := range in.Metadata {
			el.addStringEntry(e.Key, e.Value)
		}
		metaOff, metaLen = el.finish()
	}
	ob := b.StartObject(taskAssignmentSize)
	ob.SetText(taskAssignmentTaskIDOff, in.TaskID)
	ob.SetText(taskAssignmentTaskTypeOff, in.TaskType)
	ob.SetObject(taskAssignmentParamsOff, paramsOff)
	ob.SetInt32(taskAssignmentPriorityOff, in.Priority)
	ob.SetInt64(taskAssignmentCreatedTimeOff, in.CreatedTime)
	ob.SetList(taskAssignmentMetadataOff, metaOff, metaLen)
	ob.FinishAsRoot()
	return b.Finish()
}
