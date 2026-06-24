// Code generated from s3.proto; DO NOT EDIT.

package s3wire

import (
	zap "github.com/zap-proto/go"
)

// --- S3CircuitBreakerOptionsActionsEntry ---
//
// Wire entry for one map<string, int64> pair of S3CircuitBreakerOptions.actions
// (proto map field 2). The map is encoded as a repeated list of these entry
// objects; build each entry with NewS3CircuitBreakerOptionsActionsEntry and
// pass the sub-buffers to S3CircuitBreakerOptionsInput.Actions.

const (
	s3CircuitBreakerOptionsActionsEntryKeyOff   = 0
	s3CircuitBreakerOptionsActionsEntryValueOff = 8
	s3CircuitBreakerOptionsActionsEntrySize     = 16
)

// S3CircuitBreakerOptionsActionsEntry is a zero-copy view into one encoded
// map entry.
type S3CircuitBreakerOptionsActionsEntry struct{ o zap.Object }

// WrapS3CircuitBreakerOptionsActionsEntry parses b and returns a typed view.
func WrapS3CircuitBreakerOptionsActionsEntry(b []byte) (S3CircuitBreakerOptionsActionsEntry, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return S3CircuitBreakerOptionsActionsEntry{}, err
	}
	return S3CircuitBreakerOptionsActionsEntry{o: m.Root()}, nil
}

// Key reads the map key (string).
func (t S3CircuitBreakerOptionsActionsEntry) Key() string {
	return t.o.Text(s3CircuitBreakerOptionsActionsEntryKeyOff)
}

// Value reads the map value (int64).
func (t S3CircuitBreakerOptionsActionsEntry) Value() int64 {
	return t.o.Int64(s3CircuitBreakerOptionsActionsEntryValueOff)
}

// S3CircuitBreakerOptionsActionsEntryInput collects one map pair.
type S3CircuitBreakerOptionsActionsEntryInput struct {
	Key   string
	Value int64
}

// NewS3CircuitBreakerOptionsActionsEntry builds one encoded map entry.
func NewS3CircuitBreakerOptionsActionsEntry(in S3CircuitBreakerOptionsActionsEntryInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(s3CircuitBreakerOptionsActionsEntrySize)
	ob.SetText(s3CircuitBreakerOptionsActionsEntryKeyOff, in.Key)
	ob.SetInt64(s3CircuitBreakerOptionsActionsEntryValueOff, in.Value)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- S3CircuitBreakerOptions ---

const (
	s3CircuitBreakerOptionsEnabledOff = 0
	s3CircuitBreakerOptionsActionsOff = 1
	s3CircuitBreakerOptionsSize       = 9
)

// S3CircuitBreakerOptions is a zero-copy view into a ZAP-encoded
// S3CircuitBreakerOptions message.
type S3CircuitBreakerOptions struct{ o zap.Object }

// WrapS3CircuitBreakerOptions parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapS3CircuitBreakerOptions(b []byte) (S3CircuitBreakerOptions, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return S3CircuitBreakerOptions{}, err
	}
	return S3CircuitBreakerOptions{o: m.Root()}, nil
}

// Enabled reads the enabled field (proto field 1, bool).
func (t S3CircuitBreakerOptions) Enabled() bool {
	return t.o.Bool(s3CircuitBreakerOptionsEnabledOff)
}

// Actions reads the actions field (proto field 2, map<string, int64>). Each
// element is an out-of-line S3CircuitBreakerOptionsActionsEntry object; read
// element i with List.ObjectAt(i) and wrap it.
func (t S3CircuitBreakerOptions) Actions() zap.List {
	return t.o.List(s3CircuitBreakerOptionsActionsOff)
}

// S3CircuitBreakerOptionsInput collects the field values for
// NewS3CircuitBreakerOptions. Actions takes one pre-built entry sub-buffer
// (from NewS3CircuitBreakerOptionsActionsEntry) per map pair.
type S3CircuitBreakerOptionsInput struct {
	Enabled bool
	Actions [][]byte
}

// NewS3CircuitBreakerOptions builds a ZAP-encoded S3CircuitBreakerOptions
// message from in and returns the bytes.
func NewS3CircuitBreakerOptions(in S3CircuitBreakerOptionsInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(s3CircuitBreakerOptionsSize)
	ob.SetBool(s3CircuitBreakerOptionsEnabledOff, in.Enabled)
	actionsLB := b.StartList(0)
	for _, elem := range in.Actions {
		actionsLB.AddObjectBytes(elem)
	}
	ob.SetList(s3CircuitBreakerOptionsActionsOff, actionsLB.FinishOffset(), len(in.Actions))
	ob.FinishAsRoot()
	return b.Finish()
}
