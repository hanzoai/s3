// Code generated from s3_lifecycle.proto; DO NOT EDIT.

package s3_lifecyclewire

import (
	zap "github.com/zap-proto/go"
)

const (
	entryIdentityMtimeNsOff      = 0
	entryIdentitySizeOff         = 8
	entryIdentityHeadFidOff      = 16
	entryIdentityExtendedHashOff = 24
	entryIdentitySize            = 32
)

// EntryIdentity is a zero-copy view into a ZAP-encoded EntryIdentity message.
type EntryIdentity struct{ o zap.Object }

// WrapEntryIdentity parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapEntryIdentity(b []byte) (EntryIdentity, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return EntryIdentity{}, err
	}
	return EntryIdentity{o: m.Root()}, nil
}

func (t EntryIdentity) MtimeNs() int64       { return t.o.Int64(entryIdentityMtimeNsOff) }
func (t EntryIdentity) Size() int64          { return t.o.Int64(entryIdentitySizeOff) }
func (t EntryIdentity) HeadFid() string      { return t.o.Text(entryIdentityHeadFidOff) }
func (t EntryIdentity) ExtendedHash() []byte { return t.o.Bytes(entryIdentityExtendedHashOff) }

// EntryIdentityInput collects the field values for NewEntryIdentity.
type EntryIdentityInput struct {
	MtimeNs      int64
	Size         int64
	HeadFid      string
	ExtendedHash []byte
}

// putEntryIdentity lays in's fields into ob (a fresh StartObject(entryIdentitySize)).
func putEntryIdentity(ob *zap.ObjectBuilder, in EntryIdentityInput) {
	ob.SetInt64(entryIdentityMtimeNsOff, in.MtimeNs)
	ob.SetInt64(entryIdentitySizeOff, in.Size)
	ob.SetText(entryIdentityHeadFidOff, in.HeadFid)
	ob.SetBytes(entryIdentityExtendedHashOff, in.ExtendedHash)
}

// buildEntryIdentity lays an EntryIdentity object into b and returns its offset,
// for embedding as a nested field via ObjectBuilder.SetObject. A nil in yields
// offset 0 (a null nested object).
func buildEntryIdentity(b *zap.Builder, in *EntryIdentityInput) int {
	if in == nil {
		return 0
	}
	ob := b.StartObject(entryIdentitySize)
	putEntryIdentity(ob, *in)
	return ob.Finish()
}

// NewEntryIdentity builds a ZAP-encoded EntryIdentity message from in and returns the bytes.
func NewEntryIdentity(in EntryIdentityInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(entryIdentitySize)
	putEntryIdentity(ob, in)
	ob.FinishAsRoot()
	return b.Finish()
}
