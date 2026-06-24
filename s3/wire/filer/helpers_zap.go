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

// embedMessage copies a standalone child message's whole data segment into b and
// returns the offset of its root object (0 for an empty/short child).
//
// A child built by a New<Child> function is [16-byte header][data segment]. The
// data segment is NOT just the root object: a builder lays some variable data
// BEFORE the root object (lists built before StartObject, reached by a NEGATIVE
// relative pointer) and some AFTER (deferred text/bytes, nested objects, reached
// by a positive one). Copying only child[rootOffset:] would drop the pre-object
// data and dangle those back-pointers (e.g. Entry.chunks). So the entire data
// segment child[HeaderSize:] is relocated as ONE rigid block — every internal
// relative pointer, forward or backward, stays valid because the whole block
// moves together — and the returned offset is the root object's position inside
// the relocated block. WriteBytes 8-aligns the destination; rootOffset is itself
// 8-aligned (StartObject aligns; HeaderSize is a multiple of 8), so the root
// stays 8-aligned after the move.
func embedMessage(b *zap.Builder, child []byte) int {
	if len(child) < zap.HeaderSize {
		return 0
	}
	rootOff := int(binary.LittleEndian.Uint32(child[8:12]))
	if rootOff < zap.HeaderSize || rootOff >= len(child) {
		return 0
	}
	base := b.WriteBytes(child[zap.HeaderSize:])
	return base + (rootOff - zap.HeaderSize)
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
