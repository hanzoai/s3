// Code generated from filer.proto; DO NOT EDIT.

package filerwire

import (
	"encoding/binary"

	zap "github.com/zap-proto/go"
)

// setNestedObject embeds a pre-built child message (a standalone New<Child>
// buffer) as a nested object at the given 4-byte object-pointer field offset of
// ob. A zero-length child writes a null pointer (the slot keeps its zero value).
// This is the canonical singular-message embedding used across the generated
// schemas; the reader side is the matching t.o.Object(fieldOff).
//
// A ZAP message is a 16-byte header followed by a data segment; the root object
// block (the object plus its variable tail) lives at child[rootOffset:] and is
// position-independent — every text/bytes/nested pointer inside it is RELATIVE
// to its own field cell. So embedding is a rigid copy of that block into the
// parent builder at an 8-aligned position (WriteBytes aligns), then a relative
// SetObject pointer to it. No re-layout, no struct marshaling: the child's bytes
// ARE the nested message, relocated whole. Nested-in-nested survives because the
// inner block was itself laid contiguously by the same rule. This keeps the
// 4-byte object-pointer slot the schema offsets reserve (a bytes field would
// need 8 and overrun the next field).
func setNestedObject(b *zap.Builder, ob *zap.ObjectBuilder, fieldOff int, child []byte) {
	off := embedMessage(b, child)
	ob.SetObject(fieldOff, off)
}

// embedMessage copies the relocatable root-object block of a standalone child
// message into b and returns its offset (0 for an empty/short child). The block
// is child[rootOffset:], where rootOffset is the header's root pointer (bytes
// 8:12). WriteBytes 8-aligns the destination, matching the alignment the child's
// builder used, so the copied block's internal relative offsets stay valid.
func embedMessage(b *zap.Builder, child []byte) int {
	if len(child) < zap.HeaderSize {
		return 0
	}
	rootOff := binary.LittleEndian.Uint32(child[8:12])
	if int(rootOff) < zap.HeaderSize || int(rootOff) >= len(child) {
		return 0
	}
	return b.WriteBytes(child[rootOff:])
}

// addObjectList appends each pre-built sub-buffer in elems as an out-of-line
// element of a fresh list and returns (offset, count) for ObjectBuilder.SetList.
// A nil/empty slice yields (0, 0) — a null list pointer. The element shape is a
// 4-byte length prefix + bytes, exactly what List.ObjectAt / List.BytesAt read.
func addObjectList(b *zap.Builder, elems [][]byte) (offset int, length int) {
	if len(elems) == 0 {
		return 0, 0
	}
	lb := b.StartList(0)
	for _, e := range elems {
		lb.AddObjectBytes(e)
	}
	return lb.FinishOffset(), len(elems)
}

// addInt32List lays a flat list of 4-byte signed elements (repeated int32, e.g.
// the signatures fields) and returns (offset, count) for SetList.
func addInt32List(b *zap.Builder, vals []int32) (offset int, length int) {
	if len(vals) == 0 {
		return 0, 0
	}
	lb := b.StartList(4)
	for _, v := range vals {
		lb.AddUint32(uint32(v))
	}
	return lb.FinishOffset(), len(vals)
}

// addStringList lays a list of out-of-line string elements (repeated string).
// Each element is its own single-text ZAP sub-buffer read with elemText.
func addStringList(b *zap.Builder, vals []string) (offset int, length int) {
	if len(vals) == 0 {
		return 0, 0
	}
	lb := b.StartList(0)
	for _, s := range vals {
		sb := zap.NewBuilder(zap.HeaderSize + 8 + len(s))
		ob := sb.StartObject(8)
		ob.SetText(0, s)
		ob.FinishAsRoot()
		lb.AddObjectBytes(sb.Finish())
	}
	return lb.FinishOffset(), len(vals)
}

// elemText reads the i-th element of a repeated-string list as a string.
func elemText(l zap.List, i int) string { return l.ObjectAt(i).Text(0) }
