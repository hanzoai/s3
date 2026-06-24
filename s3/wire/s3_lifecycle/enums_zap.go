// Code generated from s3_lifecycle.proto; DO NOT EDIT.

package s3_lifecyclewire

// ActionKind classifies the lifecycle action a single compiled entry
// represents (proto enum ActionKind). Encoded on the wire as a uint32 scalar.
const (
	ActionKindUnspecified         uint32 = 0
	ActionKindExpirationDays      uint32 = 1
	ActionKindExpirationDate      uint32 = 2
	ActionKindNoncurrentDays      uint32 = 3
	ActionKindNewerNoncurrent     uint32 = 4
	ActionKindAbortMpu            uint32 = 5
	ActionKindExpiredDeleteMarker uint32 = 6
)

// LifecycleDeleteOutcome captures every per-event verdict the worker needs
// (proto enum LifecycleDeleteOutcome). Encoded on the wire as a uint32 scalar.
const (
	LifecycleDeleteOutcomeUnspecified       uint32 = 0
	LifecycleDeleteOutcomeDone              uint32 = 1
	LifecycleDeleteOutcomeNoopResolved      uint32 = 2
	LifecycleDeleteOutcomeSkippedObjectLock uint32 = 3
	LifecycleDeleteOutcomeRetryLater        uint32 = 4
	LifecycleDeleteOutcomeBlocked           uint32 = 5
)

// LifecycleState_RuleMode is the nested RuleMode enum of LifecycleState.
// Encoded on the wire as a uint32 scalar.
const (
	LifecycleStateRuleModeUnspecified      uint32 = 0
	LifecycleStateRuleModeEventDriven      uint32 = 1
	LifecycleStateRuleModeScanAtDate       uint32 = 2
	LifecycleStateRuleModeScanOnly         uint32 = 3
	LifecycleStateRuleModeDisabled         uint32 = 4
	LifecycleStateRuleModePendingBootstrap uint32 = 5
)

// LifecycleState_DegradedReason is the nested DegradedReason enum of
// LifecycleState. Encoded on the wire as a uint32 scalar.
const (
	LifecycleStateDegradedReasonUnspecified           uint32 = 0
	LifecycleStateDegradedReasonLagHigh               uint32 = 1
	LifecycleStateDegradedReasonPendingFull           uint32 = 2
	LifecycleStateDegradedReasonDeleteFailures        uint32 = 3
	LifecycleStateDegradedReasonOperatorPaused        uint32 = 4
	LifecycleStateDegradedReasonRetentionBelowHorizon uint32 = 5
	LifecycleStateDegradedReasonLostLog               uint32 = 6
)

// StreamKind classifies the four blockable streams (proto enum StreamKind).
// Encoded on the wire as a uint32 scalar.
const (
	StreamKindUnspecified uint32 = 0
	StreamKindOriginal    uint32 = 1
	StreamKindPredicate   uint32 = 2
	StreamKindBootstrap   uint32 = 3
	StreamKindPending     uint32 = 4
)
