// Code generated from s3_lifecycle.proto; DO NOT EDIT.

package s3_lifecyclewire

import (
	zap "github.com/zap-proto/go"
)

const (
	lifecycleDeleteRequestBucketOff               = 0
	lifecycleDeleteRequestObjectPathOff           = 8
	lifecycleDeleteRequestVersionIDOff            = 16
	lifecycleDeleteRequestRuleHashOff             = 24
	lifecycleDeleteRequestActionKindOff           = 32
	lifecycleDeleteRequestStreamShardOff          = 36
	lifecycleDeleteRequestStreamDelaySecondsOff   = 44
	lifecycleDeleteRequestStreamPositionTsNsOff   = 52
	lifecycleDeleteRequestStreamPositionOffsetOff = 60
	lifecycleDeleteRequestExpectedIdentityOff     = 68
	lifecycleDeleteRequestEngineSnapshotIDOff     = 72
	lifecycleDeleteRequestSize                    = 80
)

// LifecycleDeleteRequest is a zero-copy view into a ZAP-encoded LifecycleDeleteRequest message.
type LifecycleDeleteRequest struct{ o zap.Object }

// WrapLifecycleDeleteRequest parses b and returns a typed view. Returns an error
// if the wire-level checks (magic, version, size) fail.
func WrapLifecycleDeleteRequest(b []byte) (LifecycleDeleteRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return LifecycleDeleteRequest{}, err
	}
	return LifecycleDeleteRequest{o: m.Root()}, nil
}

func (t LifecycleDeleteRequest) Bucket() string { return t.o.Text(lifecycleDeleteRequestBucketOff) }
func (t LifecycleDeleteRequest) ObjectPath() string {
	return t.o.Text(lifecycleDeleteRequestObjectPathOff)
}
func (t LifecycleDeleteRequest) VersionID() string {
	return t.o.Text(lifecycleDeleteRequestVersionIDOff)
}
func (t LifecycleDeleteRequest) RuleHash() []byte {
	return t.o.Bytes(lifecycleDeleteRequestRuleHashOff)
}

// ActionKind reads the action_kind field (proto field 5, enum ActionKind as uint32).
func (t LifecycleDeleteRequest) ActionKind() uint32 {
	return t.o.Uint32(lifecycleDeleteRequestActionKindOff)
}
func (t LifecycleDeleteRequest) StreamShard() string {
	return t.o.Text(lifecycleDeleteRequestStreamShardOff)
}
func (t LifecycleDeleteRequest) StreamDelaySeconds() int64 {
	return t.o.Int64(lifecycleDeleteRequestStreamDelaySecondsOff)
}
func (t LifecycleDeleteRequest) StreamPositionTsNs() int64 {
	return t.o.Int64(lifecycleDeleteRequestStreamPositionTsNsOff)
}
func (t LifecycleDeleteRequest) StreamPositionOffset() int64 {
	return t.o.Int64(lifecycleDeleteRequestStreamPositionOffsetOff)
}

// ExpectedIdentity reads the expected_identity field (proto field 20, message EntryIdentity).
func (t LifecycleDeleteRequest) ExpectedIdentity() EntryIdentity {
	return EntryIdentity{o: t.o.Object(lifecycleDeleteRequestExpectedIdentityOff)}
}
func (t LifecycleDeleteRequest) EngineSnapshotID() uint64 {
	return t.o.Uint64(lifecycleDeleteRequestEngineSnapshotIDOff)
}

// LifecycleDeleteRequestInput collects the field values for NewLifecycleDeleteRequest.
type LifecycleDeleteRequestInput struct {
	Bucket               string
	ObjectPath           string
	VersionID            string
	RuleHash             []byte
	ActionKind           uint32
	StreamShard          string
	StreamDelaySeconds   int64
	StreamPositionTsNs   int64
	StreamPositionOffset int64
	ExpectedIdentity     *EntryIdentityInput
	EngineSnapshotID     uint64
}

// NewLifecycleDeleteRequest builds a ZAP-encoded LifecycleDeleteRequest message from in and returns the bytes.
func NewLifecycleDeleteRequest(in LifecycleDeleteRequestInput) []byte {
	b := zap.NewBuilder(256)
	ident := buildEntryIdentity(b, in.ExpectedIdentity)
	ob := b.StartObject(lifecycleDeleteRequestSize)
	ob.SetText(lifecycleDeleteRequestBucketOff, in.Bucket)
	ob.SetText(lifecycleDeleteRequestObjectPathOff, in.ObjectPath)
	ob.SetText(lifecycleDeleteRequestVersionIDOff, in.VersionID)
	ob.SetBytes(lifecycleDeleteRequestRuleHashOff, in.RuleHash)
	ob.SetUint32(lifecycleDeleteRequestActionKindOff, in.ActionKind)
	ob.SetText(lifecycleDeleteRequestStreamShardOff, in.StreamShard)
	ob.SetInt64(lifecycleDeleteRequestStreamDelaySecondsOff, in.StreamDelaySeconds)
	ob.SetInt64(lifecycleDeleteRequestStreamPositionTsNsOff, in.StreamPositionTsNs)
	ob.SetInt64(lifecycleDeleteRequestStreamPositionOffsetOff, in.StreamPositionOffset)
	ob.SetObject(lifecycleDeleteRequestExpectedIdentityOff, ident)
	ob.SetUint64(lifecycleDeleteRequestEngineSnapshotIDOff, in.EngineSnapshotID)
	ob.FinishAsRoot()
	return b.Finish()
}
