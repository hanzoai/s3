// Code generated from mq_schema.proto; DO NOT EDIT.

package mq_schemawire

import (
	"encoding/binary"
	"math"

	zap "github.com/zap-proto/go"
)

// --- Value (oneof kind, 13-way) ---

// ValueKind* are the oneof discriminator values for Value.kind. The value is
// the proto field number of the selected variant; 0 means no variant is set.
const (
	ValueKindNone           uint32 = 0
	ValueKindBoolValue      uint32 = 1
	ValueKindInt32Value     uint32 = 2
	ValueKindInt64Value     uint32 = 3
	ValueKindFloatValue     uint32 = 4
	ValueKindDoubleValue    uint32 = 5
	ValueKindBytesValue     uint32 = 6
	ValueKindStringValue    uint32 = 7
	ValueKindTimestampValue uint32 = 8
	ValueKindDateValue      uint32 = 9
	ValueKindDecimalValue   uint32 = 10
	ValueKindTimeValue      uint32 = 11
	ValueKindListValue      uint32 = 14
	ValueKindRecordValue    uint32 = 15
)

const (
	valueKindOff    = 0 // uint32 discriminator
	valuePayloadOff = 8 // bytes: selected variant payload
	valueSize       = 16
)

// Value is a zero-copy view into a ZAP-encoded Value message. The kind oneof is
// encoded as a uint32 discriminator (Kind) plus a single variant payload: a
// scalar variant stores its little-endian bytes (bool=1, int32/float=4,
// int64/double=8) and bytes/string store their raw payload; the timestamp,
// date, decimal, time, list and record variants store the child message's ZAP
// buffer. This is the one interned representation of the recursive value union
// (Value <-> ListValue <-> RecordValue).
type Value struct{ o zap.Object }

// WrapValue parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapValue(b []byte) (Value, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return Value{}, err
	}
	return Value{o: m.Root()}, nil
}

// Kind reports which oneof variant is set (one of the ValueKind* values).
func (t Value) Kind() uint32 { return t.o.Uint32(valueKindOff) }

func (t Value) payloadIf(which uint32) ([]byte, bool) {
	if t.Kind() != which {
		return nil, false
	}
	return t.o.Bytes(valuePayloadOff), true
}

// BoolValue reads the bool_value variant (proto field 1, bool).
func (t Value) BoolValue() (bool, bool) {
	b, ok := t.payloadIf(ValueKindBoolValue)
	if !ok || len(b) < 1 {
		return false, false
	}
	return b[0] != 0, true
}

// Int32Value reads the int32_value variant (proto field 2, int32).
func (t Value) Int32Value() (int32, bool) {
	b, ok := t.payloadIf(ValueKindInt32Value)
	if !ok || len(b) < 4 {
		return 0, false
	}
	return int32(binary.LittleEndian.Uint32(b)), true
}

// Int64Value reads the int64_value variant (proto field 3, int64).
func (t Value) Int64Value() (int64, bool) {
	b, ok := t.payloadIf(ValueKindInt64Value)
	if !ok || len(b) < 8 {
		return 0, false
	}
	return int64(binary.LittleEndian.Uint64(b)), true
}

// FloatValue reads the float_value variant (proto field 4, float).
func (t Value) FloatValue() (float32, bool) {
	b, ok := t.payloadIf(ValueKindFloatValue)
	if !ok || len(b) < 4 {
		return 0, false
	}
	return math.Float32frombits(binary.LittleEndian.Uint32(b)), true
}

// DoubleValue reads the double_value variant (proto field 5, double).
func (t Value) DoubleValue() (float64, bool) {
	b, ok := t.payloadIf(ValueKindDoubleValue)
	if !ok || len(b) < 8 {
		return 0, false
	}
	return math.Float64frombits(binary.LittleEndian.Uint64(b)), true
}

// BytesValue reads the bytes_value variant (proto field 6, bytes).
func (t Value) BytesValue() ([]byte, bool) {
	return t.payloadIf(ValueKindBytesValue)
}

// StringValue reads the string_value variant (proto field 7, string).
func (t Value) StringValue() (string, bool) {
	b, ok := t.payloadIf(ValueKindStringValue)
	if !ok {
		return "", false
	}
	return string(b), true
}

// TimestampValue reads the timestamp_value variant (proto field 8, message
// TimestampValue).
func (t Value) TimestampValue() (TimestampValue, bool) {
	b, ok := t.payloadIf(ValueKindTimestampValue)
	if !ok || len(b) == 0 {
		return TimestampValue{}, false
	}
	v, err := WrapTimestampValue(b)
	if err != nil {
		return TimestampValue{}, false
	}
	return v, true
}

// DateValue reads the date_value variant (proto field 9, message DateValue).
func (t Value) DateValue() (DateValue, bool) {
	b, ok := t.payloadIf(ValueKindDateValue)
	if !ok || len(b) == 0 {
		return DateValue{}, false
	}
	v, err := WrapDateValue(b)
	if err != nil {
		return DateValue{}, false
	}
	return v, true
}

// DecimalValue reads the decimal_value variant (proto field 10, message
// DecimalValue).
func (t Value) DecimalValue() (DecimalValue, bool) {
	b, ok := t.payloadIf(ValueKindDecimalValue)
	if !ok || len(b) == 0 {
		return DecimalValue{}, false
	}
	v, err := WrapDecimalValue(b)
	if err != nil {
		return DecimalValue{}, false
	}
	return v, true
}

// TimeValue reads the time_value variant (proto field 11, message TimeValue).
func (t Value) TimeValue() (TimeValue, bool) {
	b, ok := t.payloadIf(ValueKindTimeValue)
	if !ok || len(b) == 0 {
		return TimeValue{}, false
	}
	v, err := WrapTimeValue(b)
	if err != nil {
		return TimeValue{}, false
	}
	return v, true
}

// ListValue reads the list_value variant (proto field 14, message ListValue).
func (t Value) ListValue() (ListValue, bool) {
	b, ok := t.payloadIf(ValueKindListValue)
	if !ok || len(b) == 0 {
		return ListValue{}, false
	}
	v, err := WrapListValue(b)
	if err != nil {
		return ListValue{}, false
	}
	return v, true
}

// RecordValue reads the record_value variant (proto field 15, message
// RecordValue).
func (t Value) RecordValue() (RecordValue, bool) {
	b, ok := t.payloadIf(ValueKindRecordValue)
	if !ok || len(b) == 0 {
		return RecordValue{}, false
	}
	v, err := WrapRecordValue(b)
	if err != nil {
		return RecordValue{}, false
	}
	return v, true
}

// ValueInput collects the field values for NewValue. Set Kind to the chosen
// variant and supply exactly that variant's input: a scalar field for the
// scalar/bytes/string variants, or the child message's ZAP buffer for the
// timestamp, date, decimal, time, list and record variants.
type ValueInput struct {
	Kind           uint32
	BoolValue      bool
	Int32Value     int32
	Int64Value     int64
	FloatValue     float32
	DoubleValue    float64
	BytesValue     []byte
	StringValue    string
	TimestampValue []byte
	DateValue      []byte
	DecimalValue   []byte
	TimeValue      []byte
	ListValue      []byte
	RecordValue    []byte
}

// NewValue builds a ZAP-encoded Value message from in and returns the bytes.
func NewValue(in ValueInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(valueSize)
	ob.SetUint32(valueKindOff, in.Kind)
	switch in.Kind {
	case ValueKindBoolValue:
		var v byte
		if in.BoolValue {
			v = 1
		}
		ob.SetBytes(valuePayloadOff, []byte{v})
	case ValueKindInt32Value:
		var enc [4]byte
		binary.LittleEndian.PutUint32(enc[:], uint32(in.Int32Value))
		ob.SetBytes(valuePayloadOff, enc[:])
	case ValueKindInt64Value:
		var enc [8]byte
		binary.LittleEndian.PutUint64(enc[:], uint64(in.Int64Value))
		ob.SetBytes(valuePayloadOff, enc[:])
	case ValueKindFloatValue:
		var enc [4]byte
		binary.LittleEndian.PutUint32(enc[:], math.Float32bits(in.FloatValue))
		ob.SetBytes(valuePayloadOff, enc[:])
	case ValueKindDoubleValue:
		var enc [8]byte
		binary.LittleEndian.PutUint64(enc[:], math.Float64bits(in.DoubleValue))
		ob.SetBytes(valuePayloadOff, enc[:])
	case ValueKindBytesValue:
		ob.SetBytes(valuePayloadOff, in.BytesValue)
	case ValueKindStringValue:
		ob.SetText(valuePayloadOff, in.StringValue)
	case ValueKindTimestampValue:
		ob.SetBytes(valuePayloadOff, in.TimestampValue)
	case ValueKindDateValue:
		ob.SetBytes(valuePayloadOff, in.DateValue)
	case ValueKindDecimalValue:
		ob.SetBytes(valuePayloadOff, in.DecimalValue)
	case ValueKindTimeValue:
		ob.SetBytes(valuePayloadOff, in.TimeValue)
	case ValueKindListValue:
		ob.SetBytes(valuePayloadOff, in.ListValue)
	case ValueKindRecordValue:
		ob.SetBytes(valuePayloadOff, in.RecordValue)
	}
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ListValue ---

const (
	listValueValuesOff = 0
	listValueSize      = 8
)

// ListValue is a zero-copy view into a ZAP-encoded ListValue message.
type ListValue struct{ o zap.Object }

// WrapListValue parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapListValue(b []byte) (ListValue, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ListValue{}, err
	}
	return ListValue{o: m.Root()}, nil
}

// ValuesLen reports the number of values elements (proto field 1, repeated Value).
func (t ListValue) ValuesLen() int { return t.o.List(listValueValuesOff).Len() }

// ValueAt returns the i-th values element. The bool is false when i is out of
// range or the element fails to parse.
func (t ListValue) ValueAt(i int) (Value, bool) {
	o := t.o.List(listValueValuesOff).ObjectAt(i)
	if o.IsNull() {
		return Value{}, false
	}
	return Value{o: o}, true
}

// ListValueInput collects the field values for NewListValue. Each Values entry
// is a Value message's own ZAP buffer (from NewValue).
type ListValueInput struct {
	Values [][]byte
}

// NewListValue builds a ZAP-encoded ListValue message from in and returns the
// bytes.
func NewListValue(in ListValueInput) []byte {
	b := zap.NewBuilder(256)
	lb := b.StartList(0)
	for _, e := range in.Values {
		lb.AddObjectBytes(e)
	}
	listOff, listLen := lb.Finish()
	ob := b.StartObject(listValueSize)
	ob.SetList(listValueValuesOff, listOff, listLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- RecordValue ---

const (
	recordValueFieldsOff = 0
	recordValueSize      = 8
)

// RecordValue is a zero-copy view into a ZAP-encoded RecordValue message. The
// map<string, Value> fields member is encoded as a list of RecordValueEntry
// {key, value} objects.
type RecordValue struct{ o zap.Object }

// WrapRecordValue parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapRecordValue(b []byte) (RecordValue, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return RecordValue{}, err
	}
	return RecordValue{o: m.Root()}, nil
}

// FieldsLen reports the number of map entries (proto field 1,
// map<string, Value>).
func (t RecordValue) FieldsLen() int { return t.o.List(recordValueFieldsOff).Len() }

// FieldEntryAt returns the i-th map entry. The bool is false when i is out of
// range or the entry fails to parse.
func (t RecordValue) FieldEntryAt(i int) (RecordValueEntry, bool) {
	o := t.o.List(recordValueFieldsOff).ObjectAt(i)
	if o.IsNull() {
		return RecordValueEntry{}, false
	}
	return RecordValueEntry{o: o}, true
}

// RecordValueInput collects the field values for NewRecordValue. Each Fields
// entry is a RecordValueEntry message's own ZAP buffer (from
// NewRecordValueEntry).
type RecordValueInput struct {
	Fields [][]byte
}

// NewRecordValue builds a ZAP-encoded RecordValue message from in and returns
// the bytes.
func NewRecordValue(in RecordValueInput) []byte {
	b := zap.NewBuilder(256)
	lb := b.StartList(0)
	for _, e := range in.Fields {
		lb.AddObjectBytes(e)
	}
	listOff, listLen := lb.Finish()
	ob := b.StartObject(recordValueSize)
	ob.SetList(recordValueFieldsOff, listOff, listLen)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- RecordValueEntry (synthetic map<string, Value> entry) ---

const (
	recordValueEntryKeyOff   = 0
	recordValueEntryValueOff = 8
	recordValueEntrySize     = 16
)

// RecordValueEntry is one {key, value} pair of RecordValue.fields. It is the
// ZAP encoding of a proto map entry: a string key and a Value message buffer.
type RecordValueEntry struct{ o zap.Object }

// WrapRecordValueEntry parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapRecordValueEntry(b []byte) (RecordValueEntry, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return RecordValueEntry{}, err
	}
	return RecordValueEntry{o: m.Root()}, nil
}

// Key reads the map entry key (string).
func (t RecordValueEntry) Key() string { return t.o.Text(recordValueEntryKeyOff) }

// Value reads the map entry value (message Value). The bool is false when the
// value is absent.
func (t RecordValueEntry) Value() (Value, bool) {
	b := t.o.Bytes(recordValueEntryValueOff)
	if len(b) == 0 {
		return Value{}, false
	}
	v, err := WrapValue(b)
	if err != nil {
		return Value{}, false
	}
	return v, true
}

// RecordValueEntryInput collects the field values for NewRecordValueEntry.
// Value is the entry value's own ZAP buffer (from NewValue).
type RecordValueEntryInput struct {
	Key   string
	Value []byte
}

// NewRecordValueEntry builds a ZAP-encoded RecordValueEntry message from in and
// returns the bytes.
func NewRecordValueEntry(in RecordValueEntryInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(recordValueEntrySize)
	ob.SetText(recordValueEntryKeyOff, in.Key)
	ob.SetBytes(recordValueEntryValueOff, in.Value)
	ob.FinishAsRoot()
	return b.Finish()
}
