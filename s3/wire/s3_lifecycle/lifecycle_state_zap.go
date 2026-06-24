// Code generated from s3_lifecycle.proto; DO NOT EDIT.

package s3_lifecyclewire

import (
	zap "github.com/zap-proto/go"
)

const (
	lifecycleStateRuleHashOff               = 0
	lifecycleStateActionKindOff             = 8
	lifecycleStateRuleIDOff                 = 12
	lifecycleStateModeOff                   = 20
	lifecycleStateDegradedReasonOff         = 24
	lifecycleStateDegradedSinceNsOff        = 28
	lifecycleStateBootstrapCompleteOff      = 36
	lifecycleStateBootstrapStartedAtNsOff   = 37
	lifecycleStateBootstrapCompletedAtNsOff = 45
	lifecycleStateLastSafetyScanTsNsOff     = 53
	lifecycleStateNextSafetyScanTsNsOff     = 61
	lifecycleStateEvaluatedTotalOff         = 69
	lifecycleStateExpiredTotalOff           = 77
	lifecycleStateMetadataOnlyTotalOff      = 85
	lifecycleStateErrorTotalOff             = 93
	lifecycleStateSize                      = 101
)

// LifecycleState is a zero-copy view into a ZAP-encoded LifecycleState message.
type LifecycleState struct{ o zap.Object }

// WrapLifecycleState parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapLifecycleState(b []byte) (LifecycleState, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return LifecycleState{}, err
	}
	return LifecycleState{o: m.Root()}, nil
}

func (t LifecycleState) RuleHash() []byte { return t.o.Bytes(lifecycleStateRuleHashOff) }

// ActionKind reads the action_kind field (proto field 2, enum ActionKind as uint32).
func (t LifecycleState) ActionKind() uint32 { return t.o.Uint32(lifecycleStateActionKindOff) }
func (t LifecycleState) RuleID() string     { return t.o.Text(lifecycleStateRuleIDOff) }

// Mode reads the mode field (proto field 4, enum LifecycleState.RuleMode as uint32).
func (t LifecycleState) Mode() uint32 { return t.o.Uint32(lifecycleStateModeOff) }

// DegradedReason reads the degraded_reason field (proto field 5, enum LifecycleState.DegradedReason as uint32).
func (t LifecycleState) DegradedReason() uint32 {
	return t.o.Uint32(lifecycleStateDegradedReasonOff)
}
func (t LifecycleState) DegradedSinceNs() int64 {
	return t.o.Int64(lifecycleStateDegradedSinceNsOff)
}
func (t LifecycleState) BootstrapComplete() bool {
	return t.o.Bool(lifecycleStateBootstrapCompleteOff)
}
func (t LifecycleState) BootstrapStartedAtNs() int64 {
	return t.o.Int64(lifecycleStateBootstrapStartedAtNsOff)
}
func (t LifecycleState) BootstrapCompletedAtNs() int64 {
	return t.o.Int64(lifecycleStateBootstrapCompletedAtNsOff)
}
func (t LifecycleState) LastSafetyScanTsNs() int64 {
	return t.o.Int64(lifecycleStateLastSafetyScanTsNsOff)
}
func (t LifecycleState) NextSafetyScanTsNs() int64 {
	return t.o.Int64(lifecycleStateNextSafetyScanTsNsOff)
}
func (t LifecycleState) EvaluatedTotal() int64 { return t.o.Int64(lifecycleStateEvaluatedTotalOff) }
func (t LifecycleState) ExpiredTotal() int64   { return t.o.Int64(lifecycleStateExpiredTotalOff) }
func (t LifecycleState) MetadataOnlyTotal() int64 {
	return t.o.Int64(lifecycleStateMetadataOnlyTotalOff)
}
func (t LifecycleState) ErrorTotal() int64 { return t.o.Int64(lifecycleStateErrorTotalOff) }

// LifecycleStateInput collects the field values for NewLifecycleState.
type LifecycleStateInput struct {
	RuleHash               []byte
	ActionKind             uint32
	RuleID                 string
	Mode                   uint32
	DegradedReason         uint32
	DegradedSinceNs        int64
	BootstrapComplete      bool
	BootstrapStartedAtNs   int64
	BootstrapCompletedAtNs int64
	LastSafetyScanTsNs     int64
	NextSafetyScanTsNs     int64
	EvaluatedTotal         int64
	ExpiredTotal           int64
	MetadataOnlyTotal      int64
	ErrorTotal             int64
}

// NewLifecycleState builds a ZAP-encoded LifecycleState message from in and returns the bytes.
func NewLifecycleState(in LifecycleStateInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(lifecycleStateSize)
	ob.SetBytes(lifecycleStateRuleHashOff, in.RuleHash)
	ob.SetUint32(lifecycleStateActionKindOff, in.ActionKind)
	ob.SetText(lifecycleStateRuleIDOff, in.RuleID)
	ob.SetUint32(lifecycleStateModeOff, in.Mode)
	ob.SetUint32(lifecycleStateDegradedReasonOff, in.DegradedReason)
	ob.SetInt64(lifecycleStateDegradedSinceNsOff, in.DegradedSinceNs)
	ob.SetBool(lifecycleStateBootstrapCompleteOff, in.BootstrapComplete)
	ob.SetInt64(lifecycleStateBootstrapStartedAtNsOff, in.BootstrapStartedAtNs)
	ob.SetInt64(lifecycleStateBootstrapCompletedAtNsOff, in.BootstrapCompletedAtNs)
	ob.SetInt64(lifecycleStateLastSafetyScanTsNsOff, in.LastSafetyScanTsNs)
	ob.SetInt64(lifecycleStateNextSafetyScanTsNsOff, in.NextSafetyScanTsNs)
	ob.SetInt64(lifecycleStateEvaluatedTotalOff, in.EvaluatedTotal)
	ob.SetInt64(lifecycleStateExpiredTotalOff, in.ExpiredTotal)
	ob.SetInt64(lifecycleStateMetadataOnlyTotalOff, in.MetadataOnlyTotal)
	ob.SetInt64(lifecycleStateErrorTotalOff, in.ErrorTotal)
	ob.FinishAsRoot()
	return b.Finish()
}
