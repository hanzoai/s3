// Code generated from mq_schema.proto; DO NOT EDIT.

package mq_schemawire

import (
	"encoding/binary"

	zap "github.com/zap-proto/go"
)

// --- RecordType ---

const (
	recordTypeFieldsOff = 0
	recordTypeSize      = 8
)

// RecordType is a zero-copy view into a ZAP-encoded RecordType message.
type RecordType struct{ o zap.Object }

// WrapRecordType parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapRecordType(b []byte) (RecordType, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return RecordType{}, err
	}
	return RecordType{o: m.Root()}, nil
}

// FieldsLen reports the number of fields elements (proto field 1, repeated Field).
func (t RecordType) FieldsLen() int { return t.o.List(recordTypeFieldsOff).Len() }

// FieldAt returns the i-th fields element. The bool is false when i is out of
// range or the element fails to parse.
func (t RecordType) FieldAt(i int) (Field, bool) {
	o := t.o.List(recordTypeFieldsOff).ObjectAt(i)
	if o.IsNull() {
		return Field{}, false
	}
	return Field{o: o}, true
}

// RecordTypeInput collects the field values for NewRecordType. Each Fields
// entry is a Field message's own ZAP buffer (from NewField).
type RecordTypeInput struct {
	Fields [][]byte
}

// NewRecordType builds a ZAP-encoded RecordType message from in and returns the
// bytes.
func NewRecordType(in RecordTypeInput) []byte {
	b := zap.NewBuilder(256)
	lb := b.StartList(0)
	for _, e := range in.Fields {
		lb.AddObjectBytes(e)
	}
	listOff, listLen := lb.Finish()
	ob := b.StartObject(recordTypeSize)
	ob.SetList(recordTypeFieldsOff, listOff, listLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- Field ---

const (
	fieldNameOff       = 0
	fieldFieldIndexOff = 8
	fieldTypeOff       = 12
	fieldIsRepeatedOff = 20
	fieldIsRequiredOff = 21
	fieldSize          = 24
)

// Field is a zero-copy view into a ZAP-encoded Field message.
type Field struct{ o zap.Object }

// WrapField parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapField(b []byte) (Field, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return Field{}, err
	}
	return Field{o: m.Root()}, nil
}

// Name reads the name field (proto field 1, string).
func (t Field) Name() string { return t.o.Text(fieldNameOff) }

// FieldIndex reads the field_index field (proto field 2, int32).
func (t Field) FieldIndex() int32 { return t.o.Int32(fieldFieldIndexOff) }

// Type reads the type field (proto field 3, message Type). The bool is false
// when the field is absent.
func (t Field) Type() (Type, bool) {
	b := t.o.Bytes(fieldTypeOff)
	if len(b) == 0 {
		return Type{}, false
	}
	ty, err := WrapType(b)
	if err != nil {
		return Type{}, false
	}
	return ty, true
}

// IsRepeated reads the is_repeated field (proto field 4, bool).
func (t Field) IsRepeated() bool { return t.o.Bool(fieldIsRepeatedOff) }

// IsRequired reads the is_required field (proto field 5, bool).
func (t Field) IsRequired() bool { return t.o.Bool(fieldIsRequiredOff) }

// FieldInput collects the field values for NewField. Type is the child Type
// message's own ZAP buffer (from NewType); nil encodes an absent type.
type FieldInput struct {
	Name       string
	FieldIndex int32
	Type       []byte
	IsRepeated bool
	IsRequired bool
}

// NewField builds a ZAP-encoded Field message from in and returns the bytes.
func NewField(in FieldInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(fieldSize)
	ob.SetText(fieldNameOff, in.Name)
	ob.SetInt32(fieldFieldIndexOff, in.FieldIndex)
	ob.SetBytes(fieldTypeOff, in.Type)
	ob.SetBool(fieldIsRepeatedOff, in.IsRepeated)
	ob.SetBool(fieldIsRequiredOff, in.IsRequired)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- Type (oneof kind) ---

// TypeKind* are the oneof discriminator values for Type.kind. The value is the
// proto field number of the selected variant; 0 means no variant is set.
const (
	TypeKindNone       uint32 = 0
	TypeKindScalarType uint32 = 1
	TypeKindRecordType uint32 = 2
	TypeKindListType   uint32 = 3
)

const (
	typeKindOff    = 0 // uint32 discriminator
	typePayloadOff = 8 // bytes: selected variant payload
	typeSize       = 16
)

// Type is a zero-copy view into a ZAP-encoded Type message. The kind oneof is
// encoded as a uint32 discriminator (Kind) plus a single variant payload: the
// scalar variant stores the ScalarType enum as 4 little-endian bytes; the
// record_type and list_type variants store the child message's ZAP buffer.
type Type struct{ o zap.Object }

// WrapType parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapType(b []byte) (Type, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return Type{}, err
	}
	return Type{o: m.Root()}, nil
}

// Kind reports which oneof variant is set (one of the TypeKind* values).
func (t Type) Kind() uint32 { return t.o.Uint32(typeKindOff) }

// ScalarType reads the scalar_type variant (proto field 1, enum ScalarType).
// The bool is false when scalar_type is not the set variant.
func (t Type) ScalarType() (uint32, bool) {
	if t.Kind() != TypeKindScalarType {
		return 0, false
	}
	b := t.o.Bytes(typePayloadOff)
	if len(b) < 4 {
		return 0, false
	}
	return binary.LittleEndian.Uint32(b), true
}

// RecordType reads the record_type variant (proto field 2, message RecordType).
// The bool is false when record_type is not the set variant.
func (t Type) RecordType() (RecordType, bool) {
	if t.Kind() != TypeKindRecordType {
		return RecordType{}, false
	}
	b := t.o.Bytes(typePayloadOff)
	if len(b) == 0 {
		return RecordType{}, false
	}
	rt, err := WrapRecordType(b)
	if err != nil {
		return RecordType{}, false
	}
	return rt, true
}

// ListType reads the list_type variant (proto field 3, message ListType). The
// bool is false when list_type is not the set variant.
func (t Type) ListType() (ListType, bool) {
	if t.Kind() != TypeKindListType {
		return ListType{}, false
	}
	b := t.o.Bytes(typePayloadOff)
	if len(b) == 0 {
		return ListType{}, false
	}
	lt, err := WrapListType(b)
	if err != nil {
		return ListType{}, false
	}
	return lt, true
}

// TypeInput collects the field values for NewType. Set Kind to the chosen
// variant and supply exactly that variant's payload: ScalarType for the scalar
// variant, or RecordType / ListType as the child message's ZAP buffer.
type TypeInput struct {
	Kind       uint32
	ScalarType uint32
	RecordType []byte
	ListType   []byte
}

// NewType builds a ZAP-encoded Type message from in and returns the bytes.
func NewType(in TypeInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(typeSize)
	ob.SetUint32(typeKindOff, in.Kind)
	switch in.Kind {
	case TypeKindScalarType:
		var enc [4]byte
		binary.LittleEndian.PutUint32(enc[:], in.ScalarType)
		ob.SetBytes(typePayloadOff, enc[:])
	case TypeKindRecordType:
		ob.SetBytes(typePayloadOff, in.RecordType)
	case TypeKindListType:
		ob.SetBytes(typePayloadOff, in.ListType)
	}
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ListType ---

const (
	listTypeElementTypeOff = 0
	listTypeSize           = 8
)

// ListType is a zero-copy view into a ZAP-encoded ListType message.
type ListType struct{ o zap.Object }

// WrapListType parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapListType(b []byte) (ListType, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ListType{}, err
	}
	return ListType{o: m.Root()}, nil
}

// ElementType reads the element_type field (proto field 1, message Type). The
// bool is false when the field is absent.
func (t ListType) ElementType() (Type, bool) {
	b := t.o.Bytes(listTypeElementTypeOff)
	if len(b) == 0 {
		return Type{}, false
	}
	ty, err := WrapType(b)
	if err != nil {
		return Type{}, false
	}
	return ty, true
}

// ListTypeInput collects the field values for NewListType. ElementType is the
// child Type message's own ZAP buffer (from NewType); nil encodes an absent
// element_type.
type ListTypeInput struct {
	ElementType []byte
}

// NewListType builds a ZAP-encoded ListType message from in and returns the
// bytes.
func NewListType(in ListTypeInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(listTypeSize)
	ob.SetBytes(listTypeElementTypeOff, in.ElementType)
	ob.FinishAsRoot()
	return b.Finish()
}
