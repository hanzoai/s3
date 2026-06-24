// Code generated from s3.proto; DO NOT EDIT.

package s3wire

import (
	"encoding/binary"

	zap "github.com/zap-proto/go"
)

// embedObject splices a pre-built child ZAP buffer into the parent builder's
// tail and writes the object pointer for a singular nested-message field.
//
// A child buffer is a full ZAP message: a 16-byte header followed by a data
// segment whose internal pointers (its own strings, lists, nested objects) are
// stored as offsets RELATIVE to their field positions. Relative offsets are
// translation-invariant, so copying the child's data segment as one contiguous
// block to any position in the parent keeps every internal pointer valid. The
// parent's object pointer is then aimed at the child's ROOT object inside the
// relocated block (NOT at byte 0, which in a raw splice would be the dead child
// header). This is the singular-message counterpart to the repeated-message
// path, where each element is instead length-prefixed via AddObjectBytes and
// re-parsed independently by List.ObjectAt.
//
// An empty child leaves the field as a null object pointer (0).
func embedObject(b *zap.Builder, ob *zap.ObjectBuilder, fieldOffset int, child []byte) {
	if len(child) <= zap.HeaderSize {
		ob.SetObject(fieldOffset, 0)
		return
	}
	childRoot := int(binary.LittleEndian.Uint32(child[8:12]))
	dataStart := b.WriteBytes(child[zap.HeaderSize:])
	ob.SetObject(fieldOffset, dataStart+childRoot-zap.HeaderSize)
}
