// Code generated from mq_broker.proto; DO NOT EDIT.

package mq_brokerwire

import (
	"encoding/binary"

	zap "github.com/zap-proto/go"
)

// Shared list/element and nested-object helpers for every message
// builder/reader in the package.
//
// Singular nested message fields are stored as the child message's own ZAP
// buffer (8-byte tail pointer, read with Bytes / written with SetBytes), EXCEPT
// oneof variant values, which are embedded as a single nested object (4-byte
// object pointer) selected by the oneof's which-tag. Repeated message, repeated
// string and map<K,V> fields all encode as a single ZAP list of out-of-line
// elements: each element is a length-prefixed sub-buffer (a standalone ZAP
// message), exactly the shape List.ObjectAt / List.BytesAt decode. Repeated
// scalar fields encode as a flat fixed-stride list.

// embedMessage copies a standalone sub-message buffer verbatim into b (8-byte
// aligned) and returns the absolute offset of the embedded root object, suitable
// for ObjectBuilder.SetObject. The sub-buffer's internal pointers are all
// position-relative, so a contiguous block copy keeps them valid; the returned
// offset is (copy base + the sub-buffer's declared root offset). This is the one
// and only way a oneof variant value is laid into a parent object.
func embedMessage(b *zap.Builder, sub []byte) int {
	if len(sub) < zap.HeaderSize {
		return 0
	}
	base := b.WriteBytes(sub)
	subRoot := int(binary.LittleEndian.Uint32(sub[8:12]))
	return base + subRoot
}

// --- repeated string ---

// stringElemOff is the single string field offset in a repeated-string element.
const (
	stringElemOff  = 0
	stringElemSize = 8
)

// buildStringElem encodes one repeated-string element as its own sub-buffer.
func buildStringElem(s string) []byte {
	sb := zap.NewBuilder(zap.HeaderSize + stringElemSize + len(s))
	ob := sb.StartObject(stringElemSize)
	ob.SetText(stringElemOff, s)
	ob.FinishAsRoot()
	return sb.Finish()
}

// setStringList writes ss as a list of out-of-line string elements and returns
// the (offset, length) for ObjectBuilder.SetList.
func setStringList(b *zap.Builder, ss []string) (int, int) {
	lb := b.StartList(0)
	for _, s := range ss {
		lb.AddObjectBytes(buildStringElem(s))
	}
	return lb.Finish()
}

// stringAt reads the i-th element of a repeated-string list.
func stringAt(l zap.List, i int) string { return l.ObjectAt(i).Text(stringElemOff) }

// --- map<string,message> (value is a standalone sub-message buffer) ---

const (
	mapStringMsgKeyOff   = 0
	mapStringMsgValueOff = 8
	mapStringMsgSize     = 16
)

// buildStringMsgEntry encodes one map<string,message> entry; value is a
// standalone sub-message buffer (e.g. from NewTopicPartitionStats).
func buildStringMsgEntry(key string, value []byte) []byte {
	sb := zap.NewBuilder(zap.HeaderSize + mapStringMsgSize + len(key) + len(value))
	ob := sb.StartObject(mapStringMsgSize)
	ob.SetText(mapStringMsgKeyOff, key)
	ob.SetBytes(mapStringMsgValueOff, value)
	ob.FinishAsRoot()
	return sb.Finish()
}

// StringMsgEntry is a decoded map<string,message> entry; Value is the raw
// sub-message bytes (Wrap it with the value type's Wrap<Value>).
type StringMsgEntry struct {
	Key   string
	Value []byte
}

// setStringMsgMap writes m as a list of out-of-line entry objects.
func setStringMsgMap(b *zap.Builder, m []StringMsgEntry) (int, int) {
	lb := b.StartList(0)
	for _, e := range m {
		lb.AddObjectBytes(buildStringMsgEntry(e.Key, e.Value))
	}
	return lb.Finish()
}

// stringMsgEntryAt reads the i-th map<string,message> entry.
func stringMsgEntryAt(l zap.List, i int) StringMsgEntry {
	o := l.ObjectAt(i)
	return StringMsgEntry{
		Key:   o.Text(mapStringMsgKeyOff),
		Value: o.Bytes(mapStringMsgValueOff),
	}
}
