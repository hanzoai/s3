// Code generated from filer.proto; DO NOT EDIT.

package filerwire

import (
	zap "github.com/zap-proto/go"
)

// setNestedObject embeds a pre-built child sub-buffer as a nested object at the
// given field offset of ob. A zero-length child writes a null pointer (the slot
// keeps its zero value). This is the canonical singular-message embedding used
// across the generated schemas.
func setNestedObject(b *zap.Builder, ob *zap.ObjectBuilder, fieldOff int, child []byte) {
	if len(child) == 0 {
		return
	}
	nested := b.StartObject(len(child))
	nested.SetBytesFixed(0, child)
	ob.SetObject(fieldOff, nested.Finish())
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
