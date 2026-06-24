// Code generated from master.proto; DO NOT EDIT.

package masterwire

import (
	zap "github.com/zap-proto/go"
)

// Shared list/element helpers for every message builder/reader in the package.
// Repeated string, repeated bytes, repeated message and map<K,V> fields all
// encode as a single ZAP list of out-of-line elements: each element is a
// length-prefixed sub-buffer (a standalone ZAP message), exactly the shape
// List.ObjectAt / List.BytesAt decode. Repeated scalar fields (uint32, int64)
// encode as a flat fixed-stride list.
//
// There is one and only one way to lay out a variable-element list: build each
// element's bytes, AddObjectBytes them in order, then SetList with the
// resulting offset/length.

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

// --- repeated uint32 (flat fixed-stride list) ---

// setUint32List writes vs as a flat 4-byte-stride list.
func setUint32List(b *zap.Builder, vs []uint32) (int, int) {
	lb := b.StartList(4)
	for _, v := range vs {
		lb.AddUint32(v)
	}
	off, _ := lb.Finish()
	return off, len(vs)
}

// --- repeated int64 (flat fixed-stride list) ---

// setInt64List writes vs as a flat 8-byte-stride list.
func setInt64List(b *zap.Builder, vs []int64) (int, int) {
	lb := b.StartList(8)
	for _, v := range vs {
		lb.AddUint64(uint64(v))
	}
	off, _ := lb.Finish()
	return off, len(vs)
}

// int64At reads the i-th element of a flat int64 list.
func int64At(l zap.List, i int) int64 { return int64(l.Uint64(i)) }

// --- repeated message ---

// setMsgList writes elems (each a standalone sub-message buffer) as a list of
// out-of-line elements.
func setMsgList(b *zap.Builder, elems [][]byte) (int, int) {
	lb := b.StartList(0)
	for _, e := range elems {
		lb.AddObjectBytes(e)
	}
	return lb.Finish()
}

// --- map<string,uint32> ---

const (
	mapStringUint32KeyOff   = 0
	mapStringUint32ValueOff = 8
	mapStringUint32Size     = 12
)

// buildStringUint32Entry encodes one map<string,uint32> entry as a sub-buffer.
func buildStringUint32Entry(key string, value uint32) []byte {
	sb := zap.NewBuilder(zap.HeaderSize + mapStringUint32Size + len(key))
	ob := sb.StartObject(mapStringUint32Size)
	ob.SetText(mapStringUint32KeyOff, key)
	ob.SetUint32(mapStringUint32ValueOff, value)
	ob.FinishAsRoot()
	return sb.Finish()
}

// StringUint32Entry is a decoded map<string,uint32> entry.
type StringUint32Entry struct {
	Key   string
	Value uint32
}

// setStringUint32Map writes m as a list of out-of-line entry objects.
func setStringUint32Map(b *zap.Builder, m []StringUint32Entry) (int, int) {
	lb := b.StartList(0)
	for _, e := range m {
		lb.AddObjectBytes(buildStringUint32Entry(e.Key, e.Value))
	}
	return lb.Finish()
}

// stringUint32EntryAt reads the i-th map<string,uint32> entry.
func stringUint32EntryAt(l zap.List, i int) StringUint32Entry {
	o := l.ObjectAt(i)
	return StringUint32Entry{
		Key:   o.Text(mapStringUint32KeyOff),
		Value: o.Uint32(mapStringUint32ValueOff),
	}
}

// --- map<string,string> ---

const (
	mapStringStringKeyOff   = 0
	mapStringStringValueOff = 8
	mapStringStringSize     = 16
)

// buildStringStringEntry encodes one map<string,string> entry as a sub-buffer.
func buildStringStringEntry(key, value string) []byte {
	sb := zap.NewBuilder(zap.HeaderSize + mapStringStringSize + len(key) + len(value))
	ob := sb.StartObject(mapStringStringSize)
	ob.SetText(mapStringStringKeyOff, key)
	ob.SetText(mapStringStringValueOff, value)
	ob.FinishAsRoot()
	return sb.Finish()
}

// StringStringEntry is a decoded map<string,string> entry.
type StringStringEntry struct {
	Key   string
	Value string
}

// setStringStringMap writes m as a list of out-of-line entry objects.
func setStringStringMap(b *zap.Builder, m []StringStringEntry) (int, int) {
	lb := b.StartList(0)
	for _, e := range m {
		lb.AddObjectBytes(buildStringStringEntry(e.Key, e.Value))
	}
	return lb.Finish()
}

// stringStringEntryAt reads the i-th map<string,string> entry.
func stringStringEntryAt(l zap.List, i int) StringStringEntry {
	o := l.ObjectAt(i)
	return StringStringEntry{
		Key:   o.Text(mapStringStringKeyOff),
		Value: o.Text(mapStringStringValueOff),
	}
}

// --- map<string,message> (value is a standalone sub-message buffer) ---

const (
	mapStringMsgKeyOff   = 0
	mapStringMsgValueOff = 8
	mapStringMsgSize     = 16
)

// buildStringMsgEntry encodes one map<string,message> entry; value is a
// standalone sub-message buffer (e.g. from NewDiskInfo).
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
