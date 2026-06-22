package zaprpc

import (
	"encoding/binary"
	"errors"
	"math"
)

// Hand-rolled length-delimited binary codec for Hanzo S3 — no protobuf library,
// no reflection, no .proto at runtime. ZAP carries these bytes as transport
// (see zap_no_codec_principle: ZAP is pure transport, the app hand-rolls its
// own binary). Fields are written tag-first — tag = field<<3 | wiretype, as a
// varint — so a reader skips unknown fields by wiretype. That forward/backward
// compatibility is non-negotiable for a storage cluster doing rolling upgrades:
// a new field added by a v+1 node is silently skipped by a v node, never a
// parse error. Proto3 zero-value semantics: zero scalars are omitted on write
// and default to zero on read.
//
// The per-message MarshalZap/UnmarshalZap methods are emitted by
// cmd/protoc-gen-zap; this file is the runtime they target.

// Wire types (low 3 bits of a tag).
const (
	wireVarint  = 0 // bool, int32/64, uint32/64, enum, sint (zigzag)
	wireFixed64 = 1 // fixed64, sfixed64, double
	wireBytes   = 2 // string, bytes, embedded message, packed repeated
	wireFixed32 = 5 // fixed32, sfixed32, float
)

var (
	errTruncated = errors.New("zaprpc: truncated input")
	errBadWire   = errors.New("zaprpc: bad wire type")
	errOverflow  = errors.New("zaprpc: varint overflow")
)

// ---- Writer -------------------------------------------------------------

// Writer encodes one message. Construct with NewWriter, call the per-field
// writers in any order, then take Bytes.
type Writer struct{ buf []byte }

func NewWriter() *Writer { return &Writer{} }

// Result returns the encoded message. Valid until the next write.
func (w *Writer) Result() []byte { return w.buf }

func (w *Writer) putUvarint(v uint64) {
	for v >= 0x80 {
		w.buf = append(w.buf, byte(v)|0x80)
		v >>= 7
	}
	w.buf = append(w.buf, byte(v))
}

func (w *Writer) tag(field, wire int) { w.putUvarint(uint64(field)<<3 | uint64(wire)) }

// Bool writes a bool field (omitted when false).
func (w *Writer) Bool(field int, v bool) {
	if !v {
		return
	}
	w.tag(field, wireVarint)
	w.buf = append(w.buf, 1)
}

// Uint64 writes a uint64/enum field (omitted when zero).
func (w *Writer) Uint64(field int, v uint64) {
	if v == 0 {
		return
	}
	w.tag(field, wireVarint)
	w.putUvarint(v)
}

// Uint32 writes a uint32 field (omitted when zero).
func (w *Writer) Uint32(field int, v uint32) { w.Uint64(field, uint64(v)) }

// Int64 writes an int64 field, sign-extended like protobuf (omitted when zero).
func (w *Writer) Int64(field int, v int64) {
	if v == 0 {
		return
	}
	w.tag(field, wireVarint)
	w.putUvarint(uint64(v))
}

// Int32 writes an int32 field, sign-extended to 64 bits (omitted when zero) so
// negative values round-trip.
func (w *Writer) Int32(field int, v int32) {
	if v == 0 {
		return
	}
	w.tag(field, wireVarint)
	w.putUvarint(uint64(int64(v)))
}

// Fixed64 writes a fixed64/sfixed64 field (omitted when zero).
func (w *Writer) Fixed64(field int, v uint64) {
	if v == 0 {
		return
	}
	w.tag(field, wireFixed64)
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], v)
	w.buf = append(w.buf, b[:]...)
}

// Fixed32 writes a fixed32/sfixed32 field (omitted when zero).
func (w *Writer) Fixed32(field int, v uint32) {
	if v == 0 {
		return
	}
	w.tag(field, wireFixed32)
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], v)
	w.buf = append(w.buf, b[:]...)
}

// Double writes a double field (omitted when zero).
func (w *Writer) Double(field int, v float64) { w.Fixed64(field, math.Float64bits(v)) }

// Float writes a float field (omitted when zero).
func (w *Writer) Float(field int, v float32) { w.Fixed32(field, math.Float32bits(v)) }

// String writes a string field (omitted when empty).
func (w *Writer) String(field int, v string) {
	if len(v) == 0 {
		return
	}
	w.tag(field, wireBytes)
	w.putUvarint(uint64(len(v)))
	w.buf = append(w.buf, v...)
}

// Bytes writes a bytes field (omitted when empty).
func (w *Writer) Bytes(field int, v []byte) {
	if len(v) == 0 {
		return
	}
	w.tag(field, wireBytes)
	w.putUvarint(uint64(len(v)))
	w.buf = append(w.buf, v...)
}

// Message writes a length-delimited embedded message (always emitted, even
// when sub is empty, so presence of a non-nil submessage is preserved).
func (w *Writer) Message(field int, sub []byte) {
	w.tag(field, wireBytes)
	w.putUvarint(uint64(len(sub)))
	w.buf = append(w.buf, sub...)
}

// RawString appends a string element unconditionally (for repeated string —
// repeated fields keep their zero/empty elements).
func (w *Writer) RawString(field int, v string) {
	w.tag(field, wireBytes)
	w.putUvarint(uint64(len(v)))
	w.buf = append(w.buf, v...)
}

// RawUint64 appends a varint element unconditionally (for repeated scalars,
// unpacked form).
func (w *Writer) RawUint64(field int, v uint64) {
	w.tag(field, wireVarint)
	w.putUvarint(v)
}

// ---- Reader -------------------------------------------------------------

// Reader decodes one message. Drive it with a Next() loop, dispatching on
// Field(); call the matching typed reader, and Skip() the default case.
type Reader struct {
	buf   []byte
	pos   int
	err   error
	wire  int
	field int
}

// NewReader wraps an encoded message.
func NewReader(b []byte) *Reader { return &Reader{buf: b} }

// Err returns the first error encountered, if any.
func (r *Reader) Err() error { return r.err }

func (r *Reader) fail(e error) bool {
	if r.err == nil {
		r.err = e
	}
	return false
}

func (r *Reader) uvarint() (uint64, bool) {
	var x uint64
	var s uint
	for i := 0; ; i++ {
		if r.pos >= len(r.buf) {
			return 0, r.fail(errTruncated)
		}
		b := r.buf[r.pos]
		r.pos++
		if b < 0x80 {
			if i > 9 || (i == 9 && b > 1) {
				return 0, r.fail(errOverflow)
			}
			return x | uint64(b)<<s, true
		}
		x |= uint64(b&0x7f) << s
		s += 7
	}
}

// Next advances to the next field. Returns false at end of input or on error
// (check Err). After Next, Field() is the field number and the typed reader for
// the wire type may be called.
func (r *Reader) Next() bool {
	if r.err != nil || r.pos >= len(r.buf) {
		return false
	}
	tag, ok := r.uvarint()
	if !ok {
		return false
	}
	r.wire = int(tag & 7)
	r.field = int(tag >> 3)
	return true
}

// Field returns the field number of the current record (valid after Next).
func (r *Reader) Field() int { return r.field }

// Uint64 reads the current varint field as uint64.
func (r *Reader) Uint64() uint64 {
	if r.wire != wireVarint {
		r.fail(errBadWire)
		return 0
	}
	v, _ := r.uvarint()
	return v
}

// Uint32 reads the current varint field as uint32.
func (r *Reader) Uint32() uint32 { return uint32(r.Uint64()) }

// Int64 reads the current varint field as int64.
func (r *Reader) Int64() int64 { return int64(r.Uint64()) }

// Int32 reads the current varint field as int32 (truncating the sign-extended
// encoding).
func (r *Reader) Int32() int32 { return int32(int64(r.Uint64())) }

// Bool reads the current varint field as bool.
func (r *Reader) Bool() bool { return r.Uint64() != 0 }

// Fixed64 reads the current 64-bit field.
func (r *Reader) Fixed64() uint64 {
	if r.wire != wireFixed64 {
		r.fail(errBadWire)
		return 0
	}
	if r.pos+8 > len(r.buf) {
		r.fail(errTruncated)
		return 0
	}
	v := binary.LittleEndian.Uint64(r.buf[r.pos:])
	r.pos += 8
	return v
}

// Fixed32 reads the current 32-bit field.
func (r *Reader) Fixed32() uint32 {
	if r.wire != wireFixed32 {
		r.fail(errBadWire)
		return 0
	}
	if r.pos+4 > len(r.buf) {
		r.fail(errTruncated)
		return 0
	}
	v := binary.LittleEndian.Uint32(r.buf[r.pos:])
	r.pos += 4
	return v
}

// Double reads the current 64-bit field as float64.
func (r *Reader) Double() float64 { return math.Float64frombits(r.Fixed64()) }

// Float reads the current 32-bit field as float32.
func (r *Reader) Float() float32 { return math.Float32frombits(r.Fixed32()) }

// raw reads a length-delimited payload, returning a sub-slice of the buffer.
func (r *Reader) raw() []byte {
	if r.wire != wireBytes {
		r.fail(errBadWire)
		return nil
	}
	n, ok := r.uvarint()
	if !ok {
		return nil
	}
	if r.pos+int(n) > len(r.buf) {
		r.fail(errTruncated)
		return nil
	}
	b := r.buf[r.pos : r.pos+int(n)]
	r.pos += int(n)
	return b
}

// String reads the current length-delimited field as a string (copies).
func (r *Reader) String() string { return string(r.raw()) }

// Bytes reads the current length-delimited field as bytes (copies).
func (r *Reader) Bytes() []byte {
	b := r.raw()
	if b == nil {
		return nil
	}
	return append([]byte(nil), b...)
}

// Message reads the current length-delimited field as a raw sub-message slice
// (no copy; valid until the parent buffer is reused) for nested UnmarshalZap.
func (r *Reader) Message() []byte { return r.raw() }

// Skip advances past the current field's value (unknown fields).
func (r *Reader) Skip() {
	switch r.wire {
	case wireVarint:
		r.uvarint()
	case wireFixed64:
		r.pos += 8
	case wireFixed32:
		r.pos += 4
	case wireBytes:
		r.raw()
	default:
		r.fail(errBadWire)
	}
	if r.pos > len(r.buf) {
		r.fail(errTruncated)
	}
}
