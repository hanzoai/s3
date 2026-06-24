// Code generated from plugin.proto; DO NOT EDIT.

package pluginwire

import (
	zap "github.com/zap-proto/go"
)

// This file carries the recursive ConfigValue cluster and the reusable map-entry
// schemas. A map<K,V> field elsewhere is encoded as a list of one of these entry
// objects (each a self-contained ZAP buffer); a repeated<ConfigValue> is a list
// of ConfigValue object buffers via ValueList; a map<string,ConfigValue> via
// ValueMap.

// --- StringMapEntry (one entry of a map<string, string>) ---

const (
	stringMapEntryKeyOff   = 0
	stringMapEntryValueOff = 8
	stringMapEntrySize     = 16
)

// StringMapEntry is a zero-copy view of a single map<string,string> entry.
type StringMapEntry struct{ o zap.Object }

// WrapStringMapEntry parses b and returns a typed view.
func WrapStringMapEntry(b []byte) (StringMapEntry, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return StringMapEntry{}, err
	}
	return StringMapEntry{o: m.Root()}, nil
}

func (t StringMapEntry) Key() string   { return t.o.Text(stringMapEntryKeyOff) }
func (t StringMapEntry) Value() string { return t.o.Text(stringMapEntryValueOff) }

// StringMapEntryInput collects the field values for NewStringMapEntry.
type StringMapEntryInput struct {
	Key   string
	Value string
}

// NewStringMapEntry builds a ZAP-encoded map<string,string> entry.
func NewStringMapEntry(in StringMapEntryInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(stringMapEntrySize)
	ob.SetText(stringMapEntryKeyOff, in.Key)
	ob.SetText(stringMapEntryValueOff, in.Value)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- Int32MapEntry (one entry of a map<string, int32>) ---

const (
	int32MapEntryKeyOff   = 0
	int32MapEntryValueOff = 8
	int32MapEntrySize     = 12
)

// Int32MapEntry is a zero-copy view of a single map<string,int32> entry.
type Int32MapEntry struct{ o zap.Object }

// WrapInt32MapEntry parses b and returns a typed view.
func WrapInt32MapEntry(b []byte) (Int32MapEntry, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return Int32MapEntry{}, err
	}
	return Int32MapEntry{o: m.Root()}, nil
}

func (t Int32MapEntry) Key() string  { return t.o.Text(int32MapEntryKeyOff) }
func (t Int32MapEntry) Value() int32 { return t.o.Int32(int32MapEntryValueOff) }

// Int32MapEntryInput collects the field values for NewInt32MapEntry.
type Int32MapEntryInput struct {
	Key   string
	Value int32
}

// NewInt32MapEntry builds a ZAP-encoded map<string,int32> entry.
func NewInt32MapEntry(in Int32MapEntryInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(int32MapEntrySize)
	ob.SetText(int32MapEntryKeyOff, in.Key)
	ob.SetInt32(int32MapEntryValueOff, in.Value)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ConfigValueMapEntry (one entry of a map<string, ConfigValue>) ---

const (
	configValueMapEntryKeyOff   = 0
	configValueMapEntryValueOff = 8
	configValueMapEntrySize     = 16
)

// ConfigValueMapEntry is a zero-copy view of a single map<string,ConfigValue>
// entry. The value is a nested ConfigValue object.
type ConfigValueMapEntry struct{ o zap.Object }

// WrapConfigValueMapEntry parses b and returns a typed view.
func WrapConfigValueMapEntry(b []byte) (ConfigValueMapEntry, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ConfigValueMapEntry{}, err
	}
	return ConfigValueMapEntry{o: m.Root()}, nil
}

func (t ConfigValueMapEntry) Key() string { return t.o.Text(configValueMapEntryKeyOff) }
func (t ConfigValueMapEntry) Value() ConfigValue {
	return ConfigValue{o: childObject(t.o.Bytes(configValueMapEntryValueOff))}
}

// ConfigValueMapEntryInput collects the field values for NewConfigValueMapEntry.
// Value is a pre-built ConfigValue sub-buffer.
type ConfigValueMapEntryInput struct {
	Key   string
	Value []byte
}

// NewConfigValueMapEntry builds a ZAP-encoded map<string,ConfigValue> entry.
func NewConfigValueMapEntry(in ConfigValueMapEntryInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(configValueMapEntrySize)
	ob.SetText(configValueMapEntryKeyOff, in.Key)
	setChild(ob, configValueMapEntryValueOff, in.Value)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- StringList ---

const (
	stringListValuesOff = 0
	stringListSize      = 8
)

// StringList is a zero-copy view into a ZAP-encoded StringList message.
type StringList struct{ o zap.Object }

// WrapStringList parses b and returns a typed view.
func WrapStringList(b []byte) (StringList, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return StringList{}, err
	}
	return StringList{o: m.Root()}, nil
}

// Values reads the values field (proto field 1, repeated string). Each element
// is a raw-bytes entry; use List.BytesAt(i) for the i-th string.
func (t StringList) Values() zap.List { return t.o.List(stringListValuesOff) }

// StringListInput collects the field values for NewStringList. Each Values
// element is the raw UTF-8 bytes of one string.
type StringListInput struct {
	Values [][]byte
}

// NewStringList builds a ZAP-encoded StringList message.
func NewStringList(in StringListInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(stringListSize)
	lb := b.StartList(0)
	for _, elem := range in.Values {
		lb.AddObjectBytes(elem)
	}
	ob.SetList(stringListValuesOff, lb.FinishOffset(), len(in.Values))
	ob.FinishAsRoot()
	return b.Finish()
}

// --- Int64List ---

const (
	int64ListValuesOff = 0
	int64ListSize      = 8
)

// Int64List is a zero-copy view into a ZAP-encoded Int64List message.
type Int64List struct{ o zap.Object }

// WrapInt64List parses b and returns a typed view.
func WrapInt64List(b []byte) (Int64List, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return Int64List{}, err
	}
	return Int64List{o: m.Root()}, nil
}

// Values reads the values field (proto field 1, repeated int64) as a flat list
// of 8-byte elements; read the i-th with List.Uint64At(i) cast to int64.
func (t Int64List) Values() zap.List { return t.o.List(int64ListValuesOff) }

// Int64ListInput collects the field values for NewInt64List.
type Int64ListInput struct {
	Values []int64
}

// NewInt64List builds a ZAP-encoded Int64List message.
func NewInt64List(in Int64ListInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(int64ListSize)
	lb := b.StartList(8)
	for _, v := range in.Values {
		lb.AddUint64(uint64(v))
	}
	ob.SetList(int64ListValuesOff, lb.FinishOffset(), len(in.Values))
	ob.FinishAsRoot()
	return b.Finish()
}

// --- DoubleList ---

const (
	doubleListValuesOff = 0
	doubleListSize      = 8
)

// DoubleList is a zero-copy view into a ZAP-encoded DoubleList message.
type DoubleList struct{ o zap.Object }

// WrapDoubleList parses b and returns a typed view.
func WrapDoubleList(b []byte) (DoubleList, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return DoubleList{}, err
	}
	return DoubleList{o: m.Root()}, nil
}

// Values reads the values field (proto field 1, repeated double) as a flat list
// of 8-byte elements; read the i-th with List.Uint64At(i) and math.Float64frombits.
func (t DoubleList) Values() zap.List { return t.o.List(doubleListValuesOff) }

// DoubleListInput collects the field values for NewDoubleList.
type DoubleListInput struct {
	Values []float64
}

// NewDoubleList builds a ZAP-encoded DoubleList message.
func NewDoubleList(in DoubleListInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(doubleListSize)
	lb := b.StartList(8)
	for _, v := range in.Values {
		lb.AddUint64(float64bits(v))
	}
	ob.SetList(doubleListValuesOff, lb.FinishOffset(), len(in.Values))
	ob.FinishAsRoot()
	return b.Finish()
}

// --- BoolList ---

const (
	boolListValuesOff = 0
	boolListSize      = 8
)

// BoolList is a zero-copy view into a ZAP-encoded BoolList message.
type BoolList struct{ o zap.Object }

// WrapBoolList parses b and returns a typed view.
func WrapBoolList(b []byte) (BoolList, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return BoolList{}, err
	}
	return BoolList{o: m.Root()}, nil
}

// Values reads the values field (proto field 1, repeated bool) as a flat list
// of 1-byte elements; read the i-th with List.Uint8(i) != 0.
func (t BoolList) Values() zap.List { return t.o.List(boolListValuesOff) }

// BoolListInput collects the field values for NewBoolList.
type BoolListInput struct {
	Values []bool
}

// NewBoolList builds a ZAP-encoded BoolList message.
func NewBoolList(in BoolListInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(boolListSize)
	lb := b.StartList(1)
	for _, v := range in.Values {
		if v {
			lb.AddUint8(1)
		} else {
			lb.AddUint8(0)
		}
	}
	ob.SetList(boolListValuesOff, lb.FinishOffset(), len(in.Values))
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ValueList ---

const (
	valueListValuesOff = 0
	valueListSize      = 8
)

// ValueList is a zero-copy view into a ZAP-encoded ValueList message.
type ValueList struct{ o zap.Object }

// WrapValueList parses b and returns a typed view.
func WrapValueList(b []byte) (ValueList, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ValueList{}, err
	}
	return ValueList{o: m.Root()}, nil
}

// Values reads the values field (proto field 1, repeated ConfigValue). Each
// element is a ConfigValue object; use List.ObjectAt(i) for the i-th.
func (t ValueList) Values() zap.List { return t.o.List(valueListValuesOff) }

// ValueListInput collects the field values for NewValueList. Each Values
// element is a pre-built ConfigValue sub-buffer.
type ValueListInput struct {
	Values [][]byte
}

// NewValueList builds a ZAP-encoded ValueList message.
func NewValueList(in ValueListInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(valueListSize)
	lb := b.StartList(0)
	for _, elem := range in.Values {
		lb.AddObjectBytes(elem)
	}
	ob.SetList(valueListValuesOff, lb.FinishOffset(), len(in.Values))
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ValueMap ---

const (
	valueMapFieldsOff = 0
	valueMapSize      = 8
)

// ValueMap is a zero-copy view into a ZAP-encoded ValueMap message.
type ValueMap struct{ o zap.Object }

// WrapValueMap parses b and returns a typed view.
func WrapValueMap(b []byte) (ValueMap, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ValueMap{}, err
	}
	return ValueMap{o: m.Root()}, nil
}

// Fields reads the fields field (proto field 1, map<string,ConfigValue>). Each
// element is a ConfigValueMapEntry object; use List.ObjectAt(i) for the i-th.
func (t ValueMap) Fields() zap.List { return t.o.List(valueMapFieldsOff) }

// ValueMapInput collects the field values for NewValueMap. Each Fields element
// is a pre-built ConfigValueMapEntry sub-buffer.
type ValueMapInput struct {
	Fields [][]byte
}

// NewValueMap builds a ZAP-encoded ValueMap message.
func NewValueMap(in ValueMapInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(valueMapSize)
	lb := b.StartList(0)
	for _, elem := range in.Fields {
		lb.AddObjectBytes(elem)
	}
	ob.SetList(valueMapFieldsOff, lb.FinishOffset(), len(in.Fields))
	ob.FinishAsRoot()
	return b.Finish()
}

// --- ConfigValue (oneof kind) ---

// ConfigValue is a oneof: WhichKind reports the set member (0 = none); each
// member has its own fixed slot. Scalar members are read directly; *List and
// map/list members are nested objects.
const (
	configValueWhichOff      = 0  // u32 oneof tag (proto field number, 0 = none)
	configValueBoolOff       = 4  // bool_value (field 1)
	configValueInt64Off      = 8  // int64_value (field 2)
	configValueDoubleOff     = 16 // double_value (field 3)
	configValueStringOff     = 24 // string_value (field 4), text ptr
	configValueBytesOff      = 32 // bytes_value (field 5), bytes ptr
	configValueDurationOff   = 40 // duration_value (field 6), int64 nanos
	configValueStringListOff = 48 // string_list (field 7), nested message ptr
	configValueInt64ListOff  = 56 // int64_list (field 8), nested message ptr
	configValueDoubleListOff = 64 // double_list (field 9), nested message ptr
	configValueBoolListOff   = 72 // bool_list (field 10), nested message ptr
	configValueListValueOff  = 80 // list_value (field 11), nested message ptr
	configValueMapValueOff   = 88 // map_value (field 12), nested message ptr
	configValueSize          = 96
)

// ConfigValue is a zero-copy view into a ZAP-encoded ConfigValue message.
type ConfigValue struct{ o zap.Object }

// WrapConfigValue parses b and returns a typed view.
func WrapConfigValue(b []byte) (ConfigValue, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return ConfigValue{}, err
	}
	return ConfigValue{o: m.Root()}, nil
}

// WhichKind reports which oneof member is set (one of the ConfigValueKind*
// consts; 0 = none).
func (t ConfigValue) WhichKind() uint32 { return t.o.Uint32(configValueWhichOff) }

func (t ConfigValue) BoolValue() bool      { return t.o.Bool(configValueBoolOff) }
func (t ConfigValue) Int64Value() int64    { return t.o.Int64(configValueInt64Off) }
func (t ConfigValue) DoubleValue() float64 { return t.o.Float64(configValueDoubleOff) }
func (t ConfigValue) StringValue() string  { return t.o.Text(configValueStringOff) }
func (t ConfigValue) BytesValue() []byte   { return t.o.Bytes(configValueBytesOff) }

// DurationValue reads duration_value (proto field 6, google.protobuf.Duration)
// as int64 nanoseconds.
func (t ConfigValue) DurationValue() int64 { return t.o.Int64(configValueDurationOff) }

func (t ConfigValue) StringList() StringList {
	return StringList{o: childObject(t.o.Bytes(configValueStringListOff))}
}
func (t ConfigValue) Int64List() Int64List {
	return Int64List{o: childObject(t.o.Bytes(configValueInt64ListOff))}
}
func (t ConfigValue) DoubleList() DoubleList {
	return DoubleList{o: childObject(t.o.Bytes(configValueDoubleListOff))}
}
func (t ConfigValue) BoolList() BoolList {
	return BoolList{o: childObject(t.o.Bytes(configValueBoolListOff))}
}
func (t ConfigValue) ListValue() ValueList {
	return ValueList{o: childObject(t.o.Bytes(configValueListValueOff))}
}
func (t ConfigValue) MapValue() ValueMap {
	return ValueMap{o: childObject(t.o.Bytes(configValueMapValueOff))}
}

// ConfigValueInput collects the field values for NewConfigValue. Set Which to
// the selected member and populate only that member's field. *List/list/map
// members take a pre-built sub-buffer.
type ConfigValueInput struct {
	Which         uint32
	BoolValue     bool
	Int64Value    int64
	DoubleValue   float64
	StringValue   string
	BytesValue    []byte
	DurationValue int64
	StringList    []byte
	Int64List     []byte
	DoubleList    []byte
	BoolList      []byte
	ListValue     []byte
	MapValue      []byte
}

// NewConfigValue builds a ZAP-encoded ConfigValue message. Only the member
// selected by in.Which is written; the others retain their zero slot.
func NewConfigValue(in ConfigValueInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(configValueSize)
	ob.SetUint32(configValueWhichOff, in.Which)
	switch in.Which {
	case ConfigValueKindBool:
		ob.SetBool(configValueBoolOff, in.BoolValue)
	case ConfigValueKindInt64:
		ob.SetInt64(configValueInt64Off, in.Int64Value)
	case ConfigValueKindDouble:
		ob.SetFloat64(configValueDoubleOff, in.DoubleValue)
	case ConfigValueKindString:
		ob.SetText(configValueStringOff, in.StringValue)
	case ConfigValueKindBytes:
		ob.SetBytes(configValueBytesOff, in.BytesValue)
	case ConfigValueKindDuration:
		ob.SetInt64(configValueDurationOff, in.DurationValue)
	case ConfigValueKindStringList:
		setChild(ob, configValueStringListOff, in.StringList)
	case ConfigValueKindInt64List:
		setChild(ob, configValueInt64ListOff, in.Int64List)
	case ConfigValueKindDoubleList:
		setChild(ob, configValueDoubleListOff, in.DoubleList)
	case ConfigValueKindBoolList:
		setChild(ob, configValueBoolListOff, in.BoolList)
	case ConfigValueKindListValue:
		setChild(ob, configValueListValueOff, in.ListValue)
	case ConfigValueKindMapValue:
		setChild(ob, configValueMapValueOff, in.MapValue)
	}
	ob.FinishAsRoot()
	return b.Finish()
}
