// Code generated from worker.proto; DO NOT EDIT.

package workerwire

import (
	zap "github.com/zap-proto/go"
)

// TaskPolicy.task_config oneof tags (proto field numbers, stable wire ids).
// Zero means no variant is set.
const (
	TaskPolicyTaskConfigNone                uint32 = 0
	TaskPolicyTaskConfigVacuumConfig        uint32 = 5
	TaskPolicyTaskConfigErasureCodingConfig uint32 = 6
	TaskPolicyTaskConfigBalanceConfig       uint32 = 7
	TaskPolicyTaskConfigReplicationConfig   uint32 = 8
	TaskPolicyTaskConfigEcBalanceConfig     uint32 = 9
)

// --- TaskPolicy ---
//
// Fixed payload: enabled(1) max_concurrent(4) repeat_interval_seconds(4)
// check_interval_seconds(4) task_config_which(4 oneof tag)
// task_config_value(4 obj ptr). The oneof value is the embedded selected
// config message (nested object).

const (
	taskPolicyEnabledOff               = 0
	taskPolicyMaxConcurrentOff         = 4
	taskPolicyRepeatIntervalSecondsOff = 8
	taskPolicyCheckIntervalSecondsOff  = 12
	taskPolicyTaskConfigWhichOff       = 16
	taskPolicyTaskConfigValueOff       = 20
	taskPolicySize                     = 24
)

// TaskPolicy is a zero-copy view into a ZAP-encoded TaskPolicy message.
type TaskPolicy struct{ o zap.Object }

// WrapTaskPolicy parses b and returns a typed view.
func WrapTaskPolicy(b []byte) (TaskPolicy, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return TaskPolicy{}, err
	}
	return TaskPolicy{o: m.Root()}, nil
}

func (t TaskPolicy) Enabled() bool        { return t.o.Bool(taskPolicyEnabledOff) }
func (t TaskPolicy) MaxConcurrent() int32 { return t.o.Int32(taskPolicyMaxConcurrentOff) }
func (t TaskPolicy) RepeatIntervalSeconds() int32 {
	return t.o.Int32(taskPolicyRepeatIntervalSecondsOff)
}
func (t TaskPolicy) CheckIntervalSeconds() int32 {
	return t.o.Int32(taskPolicyCheckIntervalSecondsOff)
}

// WhichTaskConfig returns the set oneof tag (TaskPolicyTaskConfig* constants).
func (t TaskPolicy) WhichTaskConfig() uint32 { return t.o.Uint32(taskPolicyTaskConfigWhichOff) }

// VacuumConfig returns the embedded vacuum_config variant.
func (t TaskPolicy) VacuumConfig() VacuumTaskConfig {
	return VacuumTaskConfig{o: t.o.Object(taskPolicyTaskConfigValueOff)}
}

// ErasureCodingConfig returns the embedded erasure_coding_config variant.
func (t TaskPolicy) ErasureCodingConfig() ErasureCodingTaskConfig {
	return ErasureCodingTaskConfig{o: t.o.Object(taskPolicyTaskConfigValueOff)}
}

// BalanceConfig returns the embedded balance_config variant.
func (t TaskPolicy) BalanceConfig() BalanceTaskConfig {
	return BalanceTaskConfig{o: t.o.Object(taskPolicyTaskConfigValueOff)}
}

// ReplicationConfig returns the embedded replication_config variant.
func (t TaskPolicy) ReplicationConfig() ReplicationTaskConfig {
	return ReplicationTaskConfig{o: t.o.Object(taskPolicyTaskConfigValueOff)}
}

// EcBalanceConfig returns the embedded ec_balance_config variant.
func (t TaskPolicy) EcBalanceConfig() EcBalanceTaskConfig {
	return EcBalanceTaskConfig{o: t.o.Object(taskPolicyTaskConfigValueOff)}
}

// TaskPolicyInput collects the field values for NewTaskPolicy. TaskConfigWhich
// selects the oneof variant; TaskConfigValue is that variant's standalone buffer.
type TaskPolicyInput struct {
	Enabled               bool
	MaxConcurrent         int32
	RepeatIntervalSeconds int32
	CheckIntervalSeconds  int32
	TaskConfigWhich       uint32
	TaskConfigValue       []byte
}

// NewTaskPolicy builds a ZAP-encoded TaskPolicy message from in.
func NewTaskPolicy(in TaskPolicyInput) []byte {
	b := zap.NewBuilder(256)
	valueOff := 0
	if len(in.TaskConfigValue) > 0 {
		valueOff = embedMessage(b, in.TaskConfigValue)
	}
	ob := b.StartObject(taskPolicySize)
	ob.SetBool(taskPolicyEnabledOff, in.Enabled)
	ob.SetInt32(taskPolicyMaxConcurrentOff, in.MaxConcurrent)
	ob.SetInt32(taskPolicyRepeatIntervalSecondsOff, in.RepeatIntervalSeconds)
	ob.SetInt32(taskPolicyCheckIntervalSecondsOff, in.CheckIntervalSeconds)
	ob.SetUint32(taskPolicyTaskConfigWhichOff, in.TaskConfigWhich)
	ob.SetObject(taskPolicyTaskConfigValueOff, valueOff)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- MaintenancePolicy ---
//
// Fixed payload: task_policies(8 list ptr) global_max_concurrent(4)
// default_repeat_interval_seconds(4) default_check_interval_seconds(4).
// task_policies is map<string,TaskPolicy>.

const (
	maintenancePolicyTaskPoliciesOff                 = 0
	maintenancePolicyGlobalMaxConcurrentOff          = 8
	maintenancePolicyDefaultRepeatIntervalSecondsOff = 12
	maintenancePolicyDefaultCheckIntervalSecondsOff  = 16
	maintenancePolicySize                            = 20
)

// MaintenancePolicy is a zero-copy view into a ZAP-encoded MaintenancePolicy message.
type MaintenancePolicy struct{ o zap.Object }

// WrapMaintenancePolicy parses b and returns a typed view.
func WrapMaintenancePolicy(b []byte) (MaintenancePolicy, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return MaintenancePolicy{}, err
	}
	return MaintenancePolicy{o: m.Root()}, nil
}

func (t MaintenancePolicy) GlobalMaxConcurrent() int32 {
	return t.o.Int32(maintenancePolicyGlobalMaxConcurrentOff)
}
func (t MaintenancePolicy) DefaultRepeatIntervalSeconds() int32 {
	return t.o.Int32(maintenancePolicyDefaultRepeatIntervalSecondsOff)
}
func (t MaintenancePolicy) DefaultCheckIntervalSeconds() int32 {
	return t.o.Int32(maintenancePolicyDefaultCheckIntervalSecondsOff)
}

// TaskPoliciesLen returns the number of task_policies entries.
func (t MaintenancePolicy) TaskPoliciesLen() int {
	return t.o.List(maintenancePolicyTaskPoliciesOff).Len()
}

// TaskPolicyEntry returns the i-th task_policies entry: the key and the raw
// TaskPolicy value bytes (Wrap with WrapTaskPolicy).
func (t MaintenancePolicy) TaskPolicyEntry(i int) MsgMapEntry {
	return msgMapEntryAt(t.o.List(maintenancePolicyTaskPoliciesOff), i)
}

// MaintenancePolicyInput collects the field values for NewMaintenancePolicy.
// TaskPolicies values are standalone TaskPolicy buffers (NewTaskPolicy).
type MaintenancePolicyInput struct {
	TaskPolicies                 []MsgMapEntry
	GlobalMaxConcurrent          int32
	DefaultRepeatIntervalSeconds int32
	DefaultCheckIntervalSeconds  int32
}

// NewMaintenancePolicy builds a ZAP-encoded MaintenancePolicy message from in.
func NewMaintenancePolicy(in MaintenancePolicyInput) []byte {
	b := zap.NewBuilder(512)
	polOff, polLen := 0, 0
	if len(in.TaskPolicies) > 0 {
		el := startElemList(b)
		for _, e := range in.TaskPolicies {
			el.addMsgEntry(e.Key, e.Value)
		}
		polOff, polLen = el.finish()
	}
	ob := b.StartObject(maintenancePolicySize)
	ob.SetList(maintenancePolicyTaskPoliciesOff, polOff, polLen)
	ob.SetInt32(maintenancePolicyGlobalMaxConcurrentOff, in.GlobalMaxConcurrent)
	ob.SetInt32(maintenancePolicyDefaultRepeatIntervalSecondsOff, in.DefaultRepeatIntervalSeconds)
	ob.SetInt32(maintenancePolicyDefaultCheckIntervalSecondsOff, in.DefaultCheckIntervalSeconds)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- MaintenanceConfig ---
//
// Fixed payload: enabled(1) scan_interval_seconds(4) worker_timeout_seconds(4)
// task_timeout_seconds(4) retry_delay_seconds(4) max_retries(4)
// cleanup_interval_seconds(4) task_retention_seconds(4) policy(4 obj ptr).
// policy is a singular MaintenancePolicy nested object.

const (
	maintenanceConfigEnabledOff                = 0
	maintenanceConfigScanIntervalSecondsOff    = 4
	maintenanceConfigWorkerTimeoutSecondsOff   = 8
	maintenanceConfigTaskTimeoutSecondsOff     = 12
	maintenanceConfigRetryDelaySecondsOff      = 16
	maintenanceConfigMaxRetriesOff             = 20
	maintenanceConfigCleanupIntervalSecondsOff = 24
	maintenanceConfigTaskRetentionSecondsOff   = 28
	maintenanceConfigPolicyOff                 = 32
	maintenanceConfigSize                      = 36
)

// MaintenanceConfig is a zero-copy view into a ZAP-encoded MaintenanceConfig message.
type MaintenanceConfig struct{ o zap.Object }

// WrapMaintenanceConfig parses b and returns a typed view.
func WrapMaintenanceConfig(b []byte) (MaintenanceConfig, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return MaintenanceConfig{}, err
	}
	return MaintenanceConfig{o: m.Root()}, nil
}

func (t MaintenanceConfig) Enabled() bool { return t.o.Bool(maintenanceConfigEnabledOff) }
func (t MaintenanceConfig) ScanIntervalSeconds() int32 {
	return t.o.Int32(maintenanceConfigScanIntervalSecondsOff)
}
func (t MaintenanceConfig) WorkerTimeoutSeconds() int32 {
	return t.o.Int32(maintenanceConfigWorkerTimeoutSecondsOff)
}
func (t MaintenanceConfig) TaskTimeoutSeconds() int32 {
	return t.o.Int32(maintenanceConfigTaskTimeoutSecondsOff)
}
func (t MaintenanceConfig) RetryDelaySeconds() int32 {
	return t.o.Int32(maintenanceConfigRetryDelaySecondsOff)
}
func (t MaintenanceConfig) MaxRetries() int32 {
	return t.o.Int32(maintenanceConfigMaxRetriesOff)
}
func (t MaintenanceConfig) CleanupIntervalSeconds() int32 {
	return t.o.Int32(maintenanceConfigCleanupIntervalSecondsOff)
}
func (t MaintenanceConfig) TaskRetentionSeconds() int32 {
	return t.o.Int32(maintenanceConfigTaskRetentionSecondsOff)
}

// Policy returns the nested MaintenancePolicy (null-object if unset).
func (t MaintenanceConfig) Policy() MaintenancePolicy {
	return MaintenancePolicy{o: t.o.Object(maintenanceConfigPolicyOff)}
}

// HasPolicy reports whether the policy field is present.
func (t MaintenanceConfig) HasPolicy() bool {
	return !t.o.Object(maintenanceConfigPolicyOff).IsNull()
}

// MaintenanceConfigInput collects the field values for NewMaintenanceConfig.
// Policy, when non-nil, is a standalone MaintenancePolicy buffer.
type MaintenanceConfigInput struct {
	Enabled                bool
	ScanIntervalSeconds    int32
	WorkerTimeoutSeconds   int32
	TaskTimeoutSeconds     int32
	RetryDelaySeconds      int32
	MaxRetries             int32
	CleanupIntervalSeconds int32
	TaskRetentionSeconds   int32
	Policy                 []byte
}

// NewMaintenanceConfig builds a ZAP-encoded MaintenanceConfig message from in.
func NewMaintenanceConfig(in MaintenanceConfigInput) []byte {
	b := zap.NewBuilder(512)
	policyOff := 0
	if len(in.Policy) > 0 {
		policyOff = embedMessage(b, in.Policy)
	}
	ob := b.StartObject(maintenanceConfigSize)
	ob.SetBool(maintenanceConfigEnabledOff, in.Enabled)
	ob.SetInt32(maintenanceConfigScanIntervalSecondsOff, in.ScanIntervalSeconds)
	ob.SetInt32(maintenanceConfigWorkerTimeoutSecondsOff, in.WorkerTimeoutSeconds)
	ob.SetInt32(maintenanceConfigTaskTimeoutSecondsOff, in.TaskTimeoutSeconds)
	ob.SetInt32(maintenanceConfigRetryDelaySecondsOff, in.RetryDelaySeconds)
	ob.SetInt32(maintenanceConfigMaxRetriesOff, in.MaxRetries)
	ob.SetInt32(maintenanceConfigCleanupIntervalSecondsOff, in.CleanupIntervalSeconds)
	ob.SetInt32(maintenanceConfigTaskRetentionSecondsOff, in.TaskRetentionSeconds)
	ob.SetObject(maintenanceConfigPolicyOff, policyOff)
	ob.FinishAsRoot()
	return b.Finish()
}
