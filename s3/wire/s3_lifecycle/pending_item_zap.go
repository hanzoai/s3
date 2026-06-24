// Code generated from s3_lifecycle.proto; DO NOT EDIT.

package s3_lifecyclewire

import (
	zap "github.com/zap-proto/go"
)

const (
	pendingItemObjectPathOff       = 0
	pendingItemVersionIDOff        = 8
	pendingItemDueAtNsOff          = 16
	pendingItemExpectedIdentityOff = 24
	pendingItemSize                = 28
)

// PendingItem is a zero-copy view into a ZAP-encoded PendingItem message.
type PendingItem struct{ o zap.Object }

// WrapPendingItem parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapPendingItem(b []byte) (PendingItem, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return PendingItem{}, err
	}
	return PendingItem{o: m.Root()}, nil
}

func (t PendingItem) ObjectPath() string { return t.o.Text(pendingItemObjectPathOff) }
func (t PendingItem) VersionID() string  { return t.o.Text(pendingItemVersionIDOff) }
func (t PendingItem) DueAtNs() int64     { return t.o.Int64(pendingItemDueAtNsOff) }

// ExpectedIdentity reads the expected_identity field (proto field 4, message EntryIdentity).
func (t PendingItem) ExpectedIdentity() EntryIdentity {
	return EntryIdentity{o: t.o.Object(pendingItemExpectedIdentityOff)}
}

// PendingItemInput collects the field values for NewPendingItem.
type PendingItemInput struct {
	ObjectPath       string
	VersionID        string
	DueAtNs          int64
	ExpectedIdentity *EntryIdentityInput
}

// NewPendingItem builds a ZAP-encoded PendingItem message from in and returns the bytes.
func NewPendingItem(in PendingItemInput) []byte {
	b := zap.NewBuilder(256)
	ident := buildEntryIdentity(b, in.ExpectedIdentity)
	ob := b.StartObject(pendingItemSize)
	ob.SetText(pendingItemObjectPathOff, in.ObjectPath)
	ob.SetText(pendingItemVersionIDOff, in.VersionID)
	ob.SetInt64(pendingItemDueAtNsOff, in.DueAtNs)
	ob.SetObject(pendingItemExpectedIdentityOff, ident)
	ob.FinishAsRoot()
	return b.Finish()
}
