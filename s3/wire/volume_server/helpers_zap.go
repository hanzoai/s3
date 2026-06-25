// Code generated from volume_server.proto; DO NOT EDIT.

package volume_serverwire

import (
	"encoding/binary"

	zap "github.com/zap-proto/go"
)

// embedMessage copies a standalone sub-message's body into b (8-byte aligned)
// and returns the absolute offset of the embedded root object, suitable for
// ObjectBuilder.SetObject. Only the payload region (sub[HeaderSize:]) is copied —
// the sub-buffer's 16-byte envelope header is not part of the object graph — and
// the returned offset is base + (rootOff - HeaderSize). Every internal pointer
// in the embedded object is field-relative and shifts by the same HeaderSize, so
// the contiguous block copy keeps the whole graph (including deeper nested
// objects/lists/strings) valid. This is the one and only way a singular nested
// message is laid into a parent object, and it mirrors the filer wire's fixed
// embedMessage (child[HeaderSize:] copy, not the verbatim header-included copy
// that corrupted nested var-len fields).
//
// Presence semantics — the subtle part. The parent reads a nested field as
// Object(off).IsNull(), i.e. (resolved offset == 0):
//
//   - nil sub (no message present): len(sub) < HeaderSize -> return 0, a null
//     pointer == correct "field unset".
//   - PRESENT but empty (a zero-field message such as ParquetInput{}, whose
//     buffer is header-only: rootOff == HeaderSize == len, empty body): the body
//     is empty and Builder.WriteBytes([]) returns 0, which would read back null
//     and SILENTLY DROP the present variant. We instead write one aligned anchor
//     byte so the size-0 object (which reads no fields) gets a non-null, in-bounds
//     root. The old verbatim copy hit the same drop from the other side
//     (base+rootOff landed one past the copied region -> out of bounds -> null).
func embedMessage(b *zap.Builder, sub []byte) int {
	if len(sub) < zap.HeaderSize {
		return 0
	}
	rootOff := int(binary.LittleEndian.Uint32(sub[8:12]))
	if rootOff < zap.HeaderSize || rootOff > len(sub) {
		return 0
	}
	body := sub[zap.HeaderSize:]
	if len(body) == 0 {
		// Present-but-empty: anchor the size-0 object at a real in-bounds byte so
		// presence survives (the object reads no fields; only its position matters).
		return b.WriteBytes([]byte{0})
	}
	base := b.WriteBytes(body)
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

// addUint32List lays a flat list of 4-byte unsigned elements (repeated uint32,
// e.g. the volume_ids / shard_ids fields) and returns (offset, count) for
// SetList.
func addUint32List(b *zap.Builder, vals []uint32) (offset int, length int) {
	if len(vals) == 0 {
		return 0, 0
	}
	lb := b.StartList(4)
	for _, v := range vals {
		lb.AddUint32(v)
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
