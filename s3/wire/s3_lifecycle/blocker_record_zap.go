// Code generated from s3_lifecycle.proto; DO NOT EDIT.

package s3_lifecyclewire

import (
	zap "github.com/zap-proto/go"
)

const (
	blockerRecordStreamKindOff    = 0
	blockerRecordShardOff         = 4
	blockerRecordDelaySecondsOff  = 12
	blockerRecordPositionOff      = 20
	blockerRecordRuleHashOff      = 24
	blockerRecordBucketOff        = 32
	blockerRecordObjectPathOff    = 40
	blockerRecordVersionIDOff     = 48
	blockerRecordActionKindOff    = 56
	blockerRecordReasonOff        = 60
	blockerRecordLastErrorOff     = 68
	blockerRecordFirstSeenAtNsOff = 76
	blockerRecordLastRetryAtNsOff = 84
	blockerRecordRetryCountOff    = 92
	blockerRecordSize             = 96
)

// BlockerRecord is a zero-copy view into a ZAP-encoded BlockerRecord message.
type BlockerRecord struct{ o zap.Object }

// WrapBlockerRecord parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapBlockerRecord(b []byte) (BlockerRecord, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return BlockerRecord{}, err
	}
	return BlockerRecord{o: m.Root()}, nil
}

// StreamKind reads the stream_kind field (proto field 1, enum StreamKind as uint32).
func (t BlockerRecord) StreamKind() uint32 { return t.o.Uint32(blockerRecordStreamKindOff) }
func (t BlockerRecord) Shard() string      { return t.o.Text(blockerRecordShardOff) }
func (t BlockerRecord) DelaySeconds() int64 {
	return t.o.Int64(blockerRecordDelaySecondsOff)
}

// Position reads the position field (proto field 4, message MessagePosition).
func (t BlockerRecord) Position() MessagePosition {
	return MessagePosition{o: t.o.Object(blockerRecordPositionOff)}
}
func (t BlockerRecord) RuleHash() []byte   { return t.o.Bytes(blockerRecordRuleHashOff) }
func (t BlockerRecord) Bucket() string     { return t.o.Text(blockerRecordBucketOff) }
func (t BlockerRecord) ObjectPath() string { return t.o.Text(blockerRecordObjectPathOff) }
func (t BlockerRecord) VersionID() string  { return t.o.Text(blockerRecordVersionIDOff) }

// ActionKind reads the action_kind field (proto field 14, enum ActionKind as uint32).
func (t BlockerRecord) ActionKind() uint32 { return t.o.Uint32(blockerRecordActionKindOff) }
func (t BlockerRecord) Reason() string     { return t.o.Text(blockerRecordReasonOff) }
func (t BlockerRecord) LastError() string  { return t.o.Text(blockerRecordLastErrorOff) }
func (t BlockerRecord) FirstSeenAtNs() int64 {
	return t.o.Int64(blockerRecordFirstSeenAtNsOff)
}
func (t BlockerRecord) LastRetryAtNs() int64 { return t.o.Int64(blockerRecordLastRetryAtNsOff) }
func (t BlockerRecord) RetryCount() int32    { return t.o.Int32(blockerRecordRetryCountOff) }

// BlockerRecordInput collects the field values for NewBlockerRecord.
type BlockerRecordInput struct {
	StreamKind    uint32
	Shard         string
	DelaySeconds  int64
	Position      *MessagePositionInput
	RuleHash      []byte
	Bucket        string
	ObjectPath    string
	VersionID     string
	ActionKind    uint32
	Reason        string
	LastError     string
	FirstSeenAtNs int64
	LastRetryAtNs int64
	RetryCount    int32
}

// NewBlockerRecord builds a ZAP-encoded BlockerRecord message from in and returns the bytes.
func NewBlockerRecord(in BlockerRecordInput) []byte {
	b := zap.NewBuilder(256)
	pos := buildMessagePosition(b, in.Position)
	ob := b.StartObject(blockerRecordSize)
	ob.SetUint32(blockerRecordStreamKindOff, in.StreamKind)
	ob.SetText(blockerRecordShardOff, in.Shard)
	ob.SetInt64(blockerRecordDelaySecondsOff, in.DelaySeconds)
	ob.SetObject(blockerRecordPositionOff, pos)
	ob.SetBytes(blockerRecordRuleHashOff, in.RuleHash)
	ob.SetText(blockerRecordBucketOff, in.Bucket)
	ob.SetText(blockerRecordObjectPathOff, in.ObjectPath)
	ob.SetText(blockerRecordVersionIDOff, in.VersionID)
	ob.SetUint32(blockerRecordActionKindOff, in.ActionKind)
	ob.SetText(blockerRecordReasonOff, in.Reason)
	ob.SetText(blockerRecordLastErrorOff, in.LastError)
	ob.SetInt64(blockerRecordFirstSeenAtNsOff, in.FirstSeenAtNs)
	ob.SetInt64(blockerRecordLastRetryAtNsOff, in.LastRetryAtNs)
	ob.SetInt32(blockerRecordRetryCountOff, in.RetryCount)
	ob.FinishAsRoot()
	return b.Finish()
}
