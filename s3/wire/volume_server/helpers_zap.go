// Code generated from volume_server.proto; DO NOT EDIT.

package volume_serverwire

import (
	"encoding/binary"

	zap "github.com/zap-proto/go"
)

// embedMessage copies a standalone sub-message buffer verbatim into b (8-byte
// aligned) and returns the absolute offset of the embedded root object, suitable
// for ObjectBuilder.SetObject. The sub-buffer's internal pointers are all
// position-relative, so a contiguous block copy keeps them valid; the returned
// offset is (copy base + the sub-buffer's declared root offset). This is the one
// and only way a singular nested message is laid into a parent object. A
// zero-length sub yields offset 0 (a null object pointer).
func embedMessage(b *zap.Builder, sub []byte) int {
	if len(sub) < zap.HeaderSize {
		return 0
	}
	base := b.WriteBytes(sub)
	subRoot := int(binary.LittleEndian.Uint32(sub[8:12]))
	return base + subRoot
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
