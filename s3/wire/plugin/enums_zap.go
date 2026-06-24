// Code generated from plugin.proto; DO NOT EDIT.

package pluginwire

// Enum constants for plugin.proto. Each proto enum maps to a uint32 scalar on
// the wire; these consts carry the same numeric values the proto declares.

// WorkKind values (proto enum WorkKind).
const (
	WorkKindUnspecified uint32 = 0
	WorkKindDetection   uint32 = 1
	WorkKindExecution   uint32 = 2
)

// JobPriority values (proto enum JobPriority).
const (
	JobPriorityUnspecified uint32 = 0
	JobPriorityLow         uint32 = 1
	JobPriorityNormal      uint32 = 2
	JobPriorityHigh        uint32 = 3
	JobPriorityCritical    uint32 = 4
)

// JobState values (proto enum JobState).
const (
	JobStateUnspecified uint32 = 0
	JobStatePending     uint32 = 1
	JobStateAssigned    uint32 = 2
	JobStateRunning     uint32 = 3
	JobStateSucceeded   uint32 = 4
	JobStateFailed      uint32 = 5
	JobStateCanceled    uint32 = 6
)

// ConfigFieldType values (proto enum ConfigFieldType).
const (
	ConfigFieldTypeUnspecified uint32 = 0
	ConfigFieldTypeBool        uint32 = 1
	ConfigFieldTypeInt64       uint32 = 2
	ConfigFieldTypeDouble      uint32 = 3
	ConfigFieldTypeString      uint32 = 4
	ConfigFieldTypeBytes       uint32 = 5
	ConfigFieldTypeDuration    uint32 = 6
	ConfigFieldTypeEnum        uint32 = 7
	ConfigFieldTypeList        uint32 = 8
	ConfigFieldTypeObject      uint32 = 9
)

// ConfigWidget values (proto enum ConfigWidget).
const (
	ConfigWidgetUnspecified uint32 = 0
	ConfigWidgetToggle      uint32 = 1
	ConfigWidgetText        uint32 = 2
	ConfigWidgetTextarea    uint32 = 3
	ConfigWidgetNumber      uint32 = 4
	ConfigWidgetSelect      uint32 = 5
	ConfigWidgetMultiSelect uint32 = 6
	ConfigWidgetDuration    uint32 = 7
	ConfigWidgetPassword    uint32 = 8
)

// ValidationRuleType values (proto enum ValidationRuleType).
const (
	ValidationRuleTypeUnspecified uint32 = 0
	ValidationRuleTypeRegex       uint32 = 1
	ValidationRuleTypeMinLength   uint32 = 2
	ValidationRuleTypeMaxLength   uint32 = 3
	ValidationRuleTypeMinItems    uint32 = 4
	ValidationRuleTypeMaxItems    uint32 = 5
	ValidationRuleTypeCustom      uint32 = 6
)

// ActivitySource values (proto enum ActivitySource).
const (
	ActivitySourceUnspecified uint32 = 0
	ActivitySourceAdmin       uint32 = 1
	ActivitySourceDetector    uint32 = 2
	ActivitySourceExecutor    uint32 = 3
)

// Oneof discriminator tags. A oneof maps to a uint32 "which" tag (0 = none set)
// plus the selected variant object; the tag value equals the proto field number
// of the selected member.

// WorkerToAdminMessage.body members (proto field numbers).
const (
	WorkerToAdminBodyNone               uint32 = 0
	WorkerToAdminBodyHello              uint32 = 10
	WorkerToAdminBodyHeartbeat          uint32 = 11
	WorkerToAdminBodyAcknowledge        uint32 = 12
	WorkerToAdminBodyConfigSchemaResp   uint32 = 13
	WorkerToAdminBodyDetectionProposals uint32 = 14
	WorkerToAdminBodyDetectionComplete  uint32 = 15
	WorkerToAdminBodyJobProgressUpdate  uint32 = 16
	WorkerToAdminBodyJobCompleted       uint32 = 17
)

// AdminToWorkerMessage.body members (proto field numbers).
const (
	AdminToWorkerBodyNone                uint32 = 0
	AdminToWorkerBodyHello               uint32 = 10
	AdminToWorkerBodyRequestConfigSchema uint32 = 11
	AdminToWorkerBodyRunDetectionRequest uint32 = 12
	AdminToWorkerBodyExecuteJobRequest   uint32 = 13
	AdminToWorkerBodyCancelRequest       uint32 = 14
	AdminToWorkerBodyShutdown            uint32 = 15
)

// ConfigValue.kind members (proto field numbers).
const (
	ConfigValueKindNone       uint32 = 0
	ConfigValueKindBool       uint32 = 1
	ConfigValueKindInt64      uint32 = 2
	ConfigValueKindDouble     uint32 = 3
	ConfigValueKindString     uint32 = 4
	ConfigValueKindBytes      uint32 = 5
	ConfigValueKindDuration   uint32 = 6
	ConfigValueKindStringList uint32 = 7
	ConfigValueKindInt64List  uint32 = 8
	ConfigValueKindDoubleList uint32 = 9
	ConfigValueKindBoolList   uint32 = 10
	ConfigValueKindListValue  uint32 = 11
	ConfigValueKindMapValue   uint32 = 12
)
