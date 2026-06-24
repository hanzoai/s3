// Code generated from worker.proto; DO NOT EDIT.

package workerwire

import (
	"encoding/binary"

	zap "github.com/zap-proto/go"
)

// embedMessage copies a standalone sub-message buffer verbatim into b (8-byte
// aligned) and returns the absolute offset of the embedded root object, suitable
// for ObjectBuilder.SetObject. The sub-buffer's internal pointers are all
// position-relative, so a contiguous block copy keeps them valid; the returned
// offset is (copy base + the sub-buffer's declared root offset). This is the one
// and only way a singular nested message is laid into a parent object.
func embedMessage(b *zap.Builder, sub []byte) int {
	if len(sub) < zap.HeaderSize {
		return 0
	}
	base := b.WriteBytes(sub)
	subRoot := int(binary.LittleEndian.Uint32(sub[8:12]))
	return base + subRoot
}

// This file holds the shared list/element helpers used by every message
// builder/reader in the package. Repeated string, repeated bytes, repeated
// message and map<K,V> fields all encode as a single ZAP list of out-of-line
// elements: each element is a length-prefixed sub-buffer (a standalone
// ZAP message), exactly the shape List.ObjectAt / List.BytesAt decode.
//
// There is one and only one way to lay out a variable-element list: build the
// element bytes, then addElem them in order, then SetList with the explicit
// element count. ListBuilder.count tracks bytes for AddBytes, so the count
// passed to SetList must be tracked by the caller (see elemList).

// elemList accumulates length-prefixed elements into a ZAP list and tracks the
// true element count (ListBuilder.count counts bytes, not elements).
type elemList struct {
	lb    *zap.ListBuilder
	count int
}

// startElemList begins an out-of-line element list in b.
func startElemList(b *zap.Builder) *elemList {
	return &elemList{lb: b.StartList(0)}
}

// addElem appends one element (a standalone sub-buffer) with a 4-byte
// little-endian length prefix and bumps the element count.
func (e *elemList) addElem(buf []byte) {
	var lp [4]byte
	lp[0] = byte(len(buf))
	lp[1] = byte(len(buf) >> 8)
	lp[2] = byte(len(buf) >> 16)
	lp[3] = byte(len(buf) >> 24)
	e.lb.AddBytes(lp[:])
	e.lb.AddBytes(buf)
	e.count++
}

// finish returns the list offset and element count for SetList.
func (e *elemList) finish() (offset int, length int) {
	off, _ := e.lb.Finish()
	return off, e.count
}

// addText appends a string element as its own length-prefixed sub-buffer.
func (e *elemList) addText(s string) {
	sb := zap.NewBuilder(zap.HeaderSize + 8 + len(s))
	ob := sb.StartObject(8)
	ob.SetText(0, s)
	ob.FinishAsRoot()
	e.addElem(sb.Finish())
}

// elemText reads the i-th element of a repeated-string list as a string.
func elemText(l zap.List, i int) string {
	return l.ObjectAt(i).Text(0)
}

// addBytesElem appends a bytes element as its own length-prefixed sub-buffer.
func (e *elemList) addBytesElem(v []byte) {
	sb := zap.NewBuilder(zap.HeaderSize + 8 + len(v))
	ob := sb.StartObject(8)
	ob.SetBytes(0, v)
	ob.FinishAsRoot()
	e.addElem(sb.Finish())
}

// elemBytes reads the i-th element of a repeated-bytes list as a byte slice.
func elemBytes(l zap.List, i int) []byte {
	return l.ObjectAt(i).Bytes(0)
}

// mapStringEntrySize is the fixed payload of a map<string,string> entry:
// key string pointer (8) + value string pointer (8).
const (
	mapStringEntryKeyOff   = 0
	mapStringEntryValueOff = 8
	mapStringEntrySize     = 16
)

// addStringEntry appends one map<string,string> entry as a sub-buffer.
func (e *elemList) addStringEntry(key, value string) {
	sb := zap.NewBuilder(zap.HeaderSize + mapStringEntrySize + len(key) + len(value))
	ob := sb.StartObject(mapStringEntrySize)
	ob.SetText(mapStringEntryKeyOff, key)
	ob.SetText(mapStringEntryValueOff, value)
	ob.FinishAsRoot()
	e.addElem(sb.Finish())
}

// MapEntry is a decoded map<string,string> entry.
type MapEntry struct {
	Key   string
	Value string
}

// mapEntryAt reads the i-th map<string,string> entry.
func mapEntryAt(l zap.List, i int) MapEntry {
	o := l.ObjectAt(i)
	return MapEntry{Key: o.Text(mapStringEntryKeyOff), Value: o.Text(mapStringEntryValueOff)}
}

// mapMsgEntryKeyOff / mapMsgEntryValueOff lay out a map<string,message> entry:
// key string pointer (8) + value message bytes pointer (8). The value sub-
// message is stored as raw bytes so the reader can Wrap<Value> it directly.
const (
	mapMsgEntryKeyOff   = 0
	mapMsgEntryValueOff = 8
	mapMsgEntrySize     = 16
)

// addMsgEntry appends one map<string,message> entry; value is a standalone
// sub-message buffer.
func (e *elemList) addMsgEntry(key string, value []byte) {
	sb := zap.NewBuilder(zap.HeaderSize + mapMsgEntrySize + len(key) + len(value))
	ob := sb.StartObject(mapMsgEntrySize)
	ob.SetText(mapMsgEntryKeyOff, key)
	ob.SetBytes(mapMsgEntryValueOff, value)
	ob.FinishAsRoot()
	e.addElem(sb.Finish())
}

// MsgMapEntry is a decoded map<string,message> entry; Value is the raw
// sub-message bytes (Wrap it with the value type's Wrap<Value>).
type MsgMapEntry struct {
	Key   string
	Value []byte
}

// msgMapEntryAt reads the i-th map<string,message> entry.
func msgMapEntryAt(l zap.List, i int) MsgMapEntry {
	o := l.ObjectAt(i)
	return MsgMapEntry{Key: o.Text(mapMsgEntryKeyOff), Value: o.Bytes(mapMsgEntryValueOff)}
}
