// Code generated from worker.proto; DO NOT EDIT.

package workerwire

import (
	zap "github.com/zap-proto/go"
)

// TaskParams.task_params oneof tags (proto field numbers, stable wire ids).
// Zero means no variant is set.
const (
	TaskParamsTaskParamsNone                uint32 = 0
	TaskParamsTaskParamsVacuumParams        uint32 = 9
	TaskParamsTaskParamsErasureCodingParams uint32 = 10
	TaskParamsTaskParamsBalanceParams       uint32 = 11
	TaskParamsTaskParamsReplicationParams   uint32 = 12
	TaskParamsTaskParamsEcBalanceParams     uint32 = 13
	TaskParamsTaskParamsS3LifecycleParams   uint32 = 14
)

// --- TaskParams ---
//
// Fixed payload: task_id(8 ptr) volume_id(4) collection(8 ptr)
// data_center(8 ptr) rack(8 ptr) volume_size(8) sources(8 list ptr)
// targets(8 list ptr) task_params_which(4 oneof tag) task_params_value(4 obj ptr).
// sources/targets are repeated TaskSource/TaskTarget (out-of-line lists); the
// oneof value is the embedded selected variant message (nested object).

const (
	taskParamsTaskIDOff          = 0
	taskParamsVolumeIDOff        = 8
	taskParamsCollectionOff      = 12
	taskParamsDataCenterOff      = 20
	taskParamsRackOff            = 28
	taskParamsVolumeSizeOff      = 36
	taskParamsSourcesOff         = 44
	taskParamsTargetsOff         = 52
	taskParamsTaskParamsWhichOff = 60
	taskParamsTaskParamsValueOff = 64
	taskParamsSize               = 68
)

// TaskParams is a zero-copy view into a ZAP-encoded TaskParams message.
type TaskParams struct{ o zap.Object }

// WrapTaskParams parses b and returns a typed view.
func WrapTaskParams(b []byte) (TaskParams, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return TaskParams{}, err
	}
	return TaskParams{o: m.Root()}, nil
}

func (t TaskParams) TaskID() string     { return t.o.Text(taskParamsTaskIDOff) }
func (t TaskParams) VolumeID() uint32   { return t.o.Uint32(taskParamsVolumeIDOff) }
func (t TaskParams) Collection() string { return t.o.Text(taskParamsCollectionOff) }
func (t TaskParams) DataCenter() string { return t.o.Text(taskParamsDataCenterOff) }
func (t TaskParams) Rack() string       { return t.o.Text(taskParamsRackOff) }
func (t TaskParams) VolumeSize() uint64 { return t.o.Uint64(taskParamsVolumeSizeOff) }

// SourcesLen returns the number of sources elements.
func (t TaskParams) SourcesLen() int { return t.o.List(taskParamsSourcesOff).Len() }

// Sources returns the i-th sources element as a typed TaskSource view.
func (t TaskParams) Sources(i int) TaskSource {
	return TaskSource{o: t.o.List(taskParamsSourcesOff).ObjectAt(i)}
}

// TargetsLen returns the number of targets elements.
func (t TaskParams) TargetsLen() int { return t.o.List(taskParamsTargetsOff).Len() }

// Targets returns the i-th targets element as a typed TaskTarget view.
func (t TaskParams) Targets(i int) TaskTarget {
	return TaskTarget{o: t.o.List(taskParamsTargetsOff).ObjectAt(i)}
}

// WhichTaskParams returns the set oneof tag (TaskParamsTaskParams* constants).
func (t TaskParams) WhichTaskParams() uint32 { return t.o.Uint32(taskParamsTaskParamsWhichOff) }

// VacuumParams returns the embedded vacuum_params variant (valid when
// WhichTaskParams == TaskParamsTaskParamsVacuumParams).
func (t TaskParams) VacuumParams() VacuumTaskParams {
	return VacuumTaskParams{o: t.o.Object(taskParamsTaskParamsValueOff)}
}

// ErasureCodingParams returns the embedded erasure_coding_params variant.
func (t TaskParams) ErasureCodingParams() ErasureCodingTaskParams {
	return ErasureCodingTaskParams{o: t.o.Object(taskParamsTaskParamsValueOff)}
}

// BalanceParams returns the embedded balance_params variant.
func (t TaskParams) BalanceParams() BalanceTaskParams {
	return BalanceTaskParams{o: t.o.Object(taskParamsTaskParamsValueOff)}
}

// ReplicationParams returns the embedded replication_params variant.
func (t TaskParams) ReplicationParams() ReplicationTaskParams {
	return ReplicationTaskParams{o: t.o.Object(taskParamsTaskParamsValueOff)}
}

// EcBalanceParams returns the embedded ec_balance_params variant.
func (t TaskParams) EcBalanceParams() EcBalanceTaskParams {
	return EcBalanceTaskParams{o: t.o.Object(taskParamsTaskParamsValueOff)}
}

// S3LifecycleParams returns the embedded s3_lifecycle_params variant.
func (t TaskParams) S3LifecycleParams() S3LifecycleParams {
	return S3LifecycleParams{o: t.o.Object(taskParamsTaskParamsValueOff)}
}

// TaskParamsInput collects the field values for NewTaskParams. Sources and
// Targets are standalone TaskSource/TaskTarget buffers. TaskParamsWhich selects
// the oneof variant and TaskParamsValue is that variant's standalone buffer.
type TaskParamsInput struct {
	TaskID          string
	VolumeID        uint32
	Collection      string
	DataCenter      string
	Rack            string
	VolumeSize      uint64
	Sources         [][]byte
	Targets         [][]byte
	TaskParamsWhich uint32
	TaskParamsValue []byte
}

// NewTaskParams builds a ZAP-encoded TaskParams message from in.
func NewTaskParams(in TaskParamsInput) []byte {
	b := zap.NewBuilder(512)
	srcOff, srcLen := 0, 0
	if len(in.Sources) > 0 {
		el := startElemList(b)
		for _, s := range in.Sources {
			el.addElem(s)
		}
		srcOff, srcLen = el.finish()
	}
	tgtOff, tgtLen := 0, 0
	if len(in.Targets) > 0 {
		el := startElemList(b)
		for _, s := range in.Targets {
			el.addElem(s)
		}
		tgtOff, tgtLen = el.finish()
	}
	valueOff := 0
	if len(in.TaskParamsValue) > 0 {
		valueOff = embedMessage(b, in.TaskParamsValue)
	}
	ob := b.StartObject(taskParamsSize)
	ob.SetText(taskParamsTaskIDOff, in.TaskID)
	ob.SetUint32(taskParamsVolumeIDOff, in.VolumeID)
	ob.SetText(taskParamsCollectionOff, in.Collection)
	ob.SetText(taskParamsDataCenterOff, in.DataCenter)
	ob.SetText(taskParamsRackOff, in.Rack)
	ob.SetUint64(taskParamsVolumeSizeOff, in.VolumeSize)
	ob.SetList(taskParamsSourcesOff, srcOff, srcLen)
	ob.SetList(taskParamsTargetsOff, tgtOff, tgtLen)
	ob.SetUint32(taskParamsTaskParamsWhichOff, in.TaskParamsWhich)
	ob.SetObject(taskParamsTaskParamsValueOff, valueOff)
	ob.FinishAsRoot()
	return b.Finish()
}
