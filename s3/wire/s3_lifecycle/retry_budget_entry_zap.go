// Code generated from s3_lifecycle.proto; DO NOT EDIT.

package s3_lifecyclewire

import (
	zap "github.com/zap-proto/go"
)

const (
	retryBudgetEntryKeyOff                = 0
	retryBudgetEntryConsecutiveRetriesOff = 4
	retryBudgetEntryFirstSeenAtNsOff      = 8
	retryBudgetEntryLastRetryAtNsOff      = 16
	retryBudgetEntrySize                  = 24
)

// RetryBudgetEntry is a zero-copy view into a ZAP-encoded RetryBudgetEntry message.
type RetryBudgetEntry struct{ o zap.Object }

// WrapRetryBudgetEntry parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapRetryBudgetEntry(b []byte) (RetryBudgetEntry, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return RetryBudgetEntry{}, err
	}
	return RetryBudgetEntry{o: m.Root()}, nil
}

// Key reads the key field (proto field 1, message StreamKey).
func (t RetryBudgetEntry) Key() StreamKey {
	return StreamKey{o: t.o.Object(retryBudgetEntryKeyOff)}
}
func (t RetryBudgetEntry) ConsecutiveRetries() int32 {
	return t.o.Int32(retryBudgetEntryConsecutiveRetriesOff)
}
func (t RetryBudgetEntry) FirstSeenAtNs() int64 {
	return t.o.Int64(retryBudgetEntryFirstSeenAtNsOff)
}
func (t RetryBudgetEntry) LastRetryAtNs() int64 {
	return t.o.Int64(retryBudgetEntryLastRetryAtNsOff)
}

// RetryBudgetEntryInput collects the field values for NewRetryBudgetEntry.
type RetryBudgetEntryInput struct {
	Key                *StreamKeyInput
	ConsecutiveRetries int32
	FirstSeenAtNs      int64
	LastRetryAtNs      int64
}

// NewRetryBudgetEntry builds a ZAP-encoded RetryBudgetEntry message from in and returns the bytes.
func NewRetryBudgetEntry(in RetryBudgetEntryInput) []byte {
	b := zap.NewBuilder(256)
	key := buildStreamKey(b, in.Key)
	ob := b.StartObject(retryBudgetEntrySize)
	ob.SetObject(retryBudgetEntryKeyOff, key)
	ob.SetInt32(retryBudgetEntryConsecutiveRetriesOff, in.ConsecutiveRetries)
	ob.SetInt64(retryBudgetEntryFirstSeenAtNsOff, in.FirstSeenAtNs)
	ob.SetInt64(retryBudgetEntryLastRetryAtNsOff, in.LastRetryAtNs)
	ob.FinishAsRoot()
	return b.Finish()
}
