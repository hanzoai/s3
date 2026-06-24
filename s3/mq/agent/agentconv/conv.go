// Package agentconv transcodes the schema_pb messages that ride inside the
// HanzoMessagingAgent RPCs (Topic, RecordType, RecordValue, PartitionOffset)
// between their protobuf form — still the lingua franca of the broker/topic/
// pub_client/sub_client layers — and their ZAP wire form (mq_schemawire).
//
// The agent is a local bridge: its client side holds schema_pb values and must
// ship them as ZAP child buffers; its server side receives those buffers and
// must hand schema_pb values back to the publisher/subscriber clients. These
// are the one-and-only conversions for that boundary, kept here so the client
// and server call sites share a single implementation. The Value union and
// RecordType.Type tree are recursive, so the helpers recurse with them.
package agentconv

import (
	"github.com/hanzoai/s3/s3/pb/schema_pb"
	mq_schemawire "github.com/hanzoai/s3/s3/wire/mq_schema"
)

// --- Topic ---

// TopicToWire encodes a schema_pb.Topic as a mq_schemawire Topic buffer; nil in
// yields nil (an absent child).
func TopicToWire(in *schema_pb.Topic) []byte {
	if in == nil {
		return nil
	}
	return mq_schemawire.NewTopic(mq_schemawire.TopicInput{
		Namespace: in.Namespace,
		Name:      in.Name,
	})
}

// TopicFromWire decodes a mq_schemawire Topic buffer into a schema_pb.Topic;
// empty in yields nil.
func TopicFromWire(b []byte) *schema_pb.Topic {
	if len(b) == 0 {
		return nil
	}
	t, err := mq_schemawire.WrapTopic(b)
	if err != nil {
		return nil
	}
	return &schema_pb.Topic{Namespace: t.Namespace(), Name: t.Name()}
}

// --- Partition / PartitionOffset ---

func partitionToWire(in *schema_pb.Partition) []byte {
	if in == nil {
		return nil
	}
	return mq_schemawire.NewPartition(mq_schemawire.PartitionInput{
		RingSize:   in.RingSize,
		RangeStart: in.RangeStart,
		RangeStop:  in.RangeStop,
		UnixTimeNs: in.UnixTimeNs,
	})
}

func partitionFromWire(p mq_schemawire.Partition) *schema_pb.Partition {
	return &schema_pb.Partition{
		RingSize:   p.RingSize(),
		RangeStart: p.RangeStart(),
		RangeStop:  p.RangeStop(),
		UnixTimeNs: p.UnixTimeNs(),
	}
}

// PartitionOffsetToWire encodes a schema_pb.PartitionOffset as a mq_schemawire
// PartitionOffset buffer; nil in yields nil.
func PartitionOffsetToWire(in *schema_pb.PartitionOffset) []byte {
	if in == nil {
		return nil
	}
	return mq_schemawire.NewPartitionOffset(mq_schemawire.PartitionOffsetInput{
		Partition:   partitionToWire(in.Partition),
		StartTsNs:   in.StartTsNs,
		StartOffset: in.StartOffset,
	})
}

// PartitionOffsetsToWire encodes a slice of schema_pb.PartitionOffset as the
// list of ZAP child buffers the wire builders consume.
func PartitionOffsetsToWire(in []*schema_pb.PartitionOffset) [][]byte {
	if len(in) == 0 {
		return nil
	}
	out := make([][]byte, 0, len(in))
	for _, po := range in {
		out = append(out, PartitionOffsetToWire(po))
	}
	return out
}

func partitionOffsetFromWire(po mq_schemawire.PartitionOffset) *schema_pb.PartitionOffset {
	out := &schema_pb.PartitionOffset{
		StartTsNs:   po.StartTsNs(),
		StartOffset: po.StartOffset(),
	}
	if p, ok := po.Partition(); ok {
		out.Partition = partitionFromWire(p)
	}
	return out
}

// PartitionOffsetsFromWire decodes the InitSubscribeRecordRequest partition
// offsets (indexed accessor) into schema_pb.PartitionOffset values.
func PartitionOffsetsFromWire(init mq_schemawire.PartitionOffset, _ int) *schema_pb.PartitionOffset {
	return partitionOffsetFromWire(init)
}

// --- RecordType (Field / Type / ListType tree) ---

// RecordTypeToWire encodes a schema_pb.RecordType as a mq_schemawire RecordType
// buffer; nil in yields nil.
func RecordTypeToWire(in *schema_pb.RecordType) []byte {
	if in == nil {
		return nil
	}
	fields := make([][]byte, 0, len(in.Fields))
	for _, f := range in.Fields {
		fields = append(fields, fieldToWire(f))
	}
	return mq_schemawire.NewRecordType(mq_schemawire.RecordTypeInput{Fields: fields})
}

func fieldToWire(in *schema_pb.Field) []byte {
	if in == nil {
		return nil
	}
	return mq_schemawire.NewField(mq_schemawire.FieldInput{
		Name:       in.Name,
		FieldIndex: in.FieldIndex,
		Type:       typeToWire(in.Type),
		IsRepeated: in.IsRepeated,
		IsRequired: in.IsRequired,
	})
}

func typeToWire(in *schema_pb.Type) []byte {
	if in == nil {
		return nil
	}
	switch k := in.Kind.(type) {
	case *schema_pb.Type_ScalarType:
		return mq_schemawire.NewType(mq_schemawire.TypeInput{
			Kind:       mq_schemawire.TypeKindScalarType,
			ScalarType: uint32(k.ScalarType),
		})
	case *schema_pb.Type_RecordType:
		return mq_schemawire.NewType(mq_schemawire.TypeInput{
			Kind:       mq_schemawire.TypeKindRecordType,
			RecordType: RecordTypeToWire(k.RecordType),
		})
	case *schema_pb.Type_ListType:
		return mq_schemawire.NewType(mq_schemawire.TypeInput{
			Kind:     mq_schemawire.TypeKindListType,
			ListType: listTypeToWire(k.ListType),
		})
	default:
		return mq_schemawire.NewType(mq_schemawire.TypeInput{Kind: mq_schemawire.TypeKindNone})
	}
}

func listTypeToWire(in *schema_pb.ListType) []byte {
	if in == nil {
		return nil
	}
	return mq_schemawire.NewListType(mq_schemawire.ListTypeInput{
		ElementType: typeToWire(in.ElementType),
	})
}

// RecordTypeFromWire decodes a mq_schemawire RecordType buffer into a
// schema_pb.RecordType; empty b yields nil.
func RecordTypeFromWire(b []byte) *schema_pb.RecordType {
	if len(b) == 0 {
		return nil
	}
	rt, err := mq_schemawire.WrapRecordType(b)
	if err != nil {
		return nil
	}
	return recordTypeFromWire(rt)
}

func recordTypeFromWire(rt mq_schemawire.RecordType) *schema_pb.RecordType {
	n := rt.FieldsLen()
	fields := make([]*schema_pb.Field, 0, n)
	for i := 0; i < n; i++ {
		if f, ok := rt.FieldAt(i); ok {
			fields = append(fields, fieldFromWire(f))
		}
	}
	return &schema_pb.RecordType{Fields: fields}
}

func fieldFromWire(f mq_schemawire.Field) *schema_pb.Field {
	out := &schema_pb.Field{
		Name:       f.Name(),
		FieldIndex: f.FieldIndex(),
		IsRepeated: f.IsRepeated(),
		IsRequired: f.IsRequired(),
	}
	if ty, ok := f.Type(); ok {
		out.Type = typeFromWire(ty)
	}
	return out
}

func typeFromWire(ty mq_schemawire.Type) *schema_pb.Type {
	switch ty.Kind() {
	case mq_schemawire.TypeKindScalarType:
		if st, ok := ty.ScalarType(); ok {
			return &schema_pb.Type{Kind: &schema_pb.Type_ScalarType{ScalarType: schema_pb.ScalarType(st)}}
		}
	case mq_schemawire.TypeKindRecordType:
		if rt, ok := ty.RecordType(); ok {
			return &schema_pb.Type{Kind: &schema_pb.Type_RecordType{RecordType: recordTypeFromWire(rt)}}
		}
	case mq_schemawire.TypeKindListType:
		if lt, ok := ty.ListType(); ok {
			return &schema_pb.Type{Kind: &schema_pb.Type_ListType{ListType: listTypeFromWire(lt)}}
		}
	}
	return nil
}

func listTypeFromWire(lt mq_schemawire.ListType) *schema_pb.ListType {
	out := &schema_pb.ListType{}
	if et, ok := lt.ElementType(); ok {
		out.ElementType = typeFromWire(et)
	}
	return out
}

// --- RecordValue (Value union) ---

// RecordValueToWire encodes a schema_pb.RecordValue as a mq_schemawire
// RecordValue buffer; nil in yields nil.
func RecordValueToWire(in *schema_pb.RecordValue) []byte {
	if in == nil {
		return nil
	}
	entries := make([][]byte, 0, len(in.Fields))
	for k, v := range in.Fields {
		entries = append(entries, mq_schemawire.NewRecordValueEntry(mq_schemawire.RecordValueEntryInput{
			Key:   k,
			Value: valueToWire(v),
		}))
	}
	return mq_schemawire.NewRecordValue(mq_schemawire.RecordValueInput{Fields: entries})
}

func valueToWire(in *schema_pb.Value) []byte {
	if in == nil {
		return nil
	}
	switch k := in.Kind.(type) {
	case *schema_pb.Value_BoolValue:
		return mq_schemawire.NewValue(mq_schemawire.ValueInput{Kind: mq_schemawire.ValueKindBoolValue, BoolValue: k.BoolValue})
	case *schema_pb.Value_Int32Value:
		return mq_schemawire.NewValue(mq_schemawire.ValueInput{Kind: mq_schemawire.ValueKindInt32Value, Int32Value: k.Int32Value})
	case *schema_pb.Value_Int64Value:
		return mq_schemawire.NewValue(mq_schemawire.ValueInput{Kind: mq_schemawire.ValueKindInt64Value, Int64Value: k.Int64Value})
	case *schema_pb.Value_FloatValue:
		return mq_schemawire.NewValue(mq_schemawire.ValueInput{Kind: mq_schemawire.ValueKindFloatValue, FloatValue: k.FloatValue})
	case *schema_pb.Value_DoubleValue:
		return mq_schemawire.NewValue(mq_schemawire.ValueInput{Kind: mq_schemawire.ValueKindDoubleValue, DoubleValue: k.DoubleValue})
	case *schema_pb.Value_BytesValue:
		return mq_schemawire.NewValue(mq_schemawire.ValueInput{Kind: mq_schemawire.ValueKindBytesValue, BytesValue: k.BytesValue})
	case *schema_pb.Value_StringValue:
		return mq_schemawire.NewValue(mq_schemawire.ValueInput{Kind: mq_schemawire.ValueKindStringValue, StringValue: k.StringValue})
	case *schema_pb.Value_TimestampValue:
		var child []byte
		if k.TimestampValue != nil {
			child = mq_schemawire.NewTimestampValue(mq_schemawire.TimestampValueInput{
				TimestampMicros: k.TimestampValue.TimestampMicros,
				IsUtc:           k.TimestampValue.IsUtc,
			})
		}
		return mq_schemawire.NewValue(mq_schemawire.ValueInput{Kind: mq_schemawire.ValueKindTimestampValue, TimestampValue: child})
	case *schema_pb.Value_DateValue:
		var child []byte
		if k.DateValue != nil {
			child = mq_schemawire.NewDateValue(mq_schemawire.DateValueInput{DaysSinceEpoch: k.DateValue.DaysSinceEpoch})
		}
		return mq_schemawire.NewValue(mq_schemawire.ValueInput{Kind: mq_schemawire.ValueKindDateValue, DateValue: child})
	case *schema_pb.Value_DecimalValue:
		var child []byte
		if k.DecimalValue != nil {
			child = mq_schemawire.NewDecimalValue(mq_schemawire.DecimalValueInput{
				Value:     k.DecimalValue.Value,
				Precision: k.DecimalValue.Precision,
				Scale:     k.DecimalValue.Scale,
			})
		}
		return mq_schemawire.NewValue(mq_schemawire.ValueInput{Kind: mq_schemawire.ValueKindDecimalValue, DecimalValue: child})
	case *schema_pb.Value_TimeValue:
		var child []byte
		if k.TimeValue != nil {
			child = mq_schemawire.NewTimeValue(mq_schemawire.TimeValueInput{TimeMicros: k.TimeValue.TimeMicros})
		}
		return mq_schemawire.NewValue(mq_schemawire.ValueInput{Kind: mq_schemawire.ValueKindTimeValue, TimeValue: child})
	case *schema_pb.Value_ListValue:
		var child []byte
		if k.ListValue != nil {
			elems := make([][]byte, 0, len(k.ListValue.Values))
			for _, e := range k.ListValue.Values {
				elems = append(elems, valueToWire(e))
			}
			child = mq_schemawire.NewListValue(mq_schemawire.ListValueInput{Values: elems})
		}
		return mq_schemawire.NewValue(mq_schemawire.ValueInput{Kind: mq_schemawire.ValueKindListValue, ListValue: child})
	case *schema_pb.Value_RecordValue:
		return mq_schemawire.NewValue(mq_schemawire.ValueInput{Kind: mq_schemawire.ValueKindRecordValue, RecordValue: RecordValueToWire(k.RecordValue)})
	default:
		return mq_schemawire.NewValue(mq_schemawire.ValueInput{Kind: mq_schemawire.ValueKindNone})
	}
}

// RecordValueFromWire decodes a mq_schemawire RecordValue buffer into a
// schema_pb.RecordValue; empty b yields nil.
func RecordValueFromWire(b []byte) *schema_pb.RecordValue {
	if len(b) == 0 {
		return nil
	}
	rv, err := mq_schemawire.WrapRecordValue(b)
	if err != nil {
		return nil
	}
	return recordValueFromWire(rv)
}

func recordValueFromWire(rv mq_schemawire.RecordValue) *schema_pb.RecordValue {
	n := rv.FieldsLen()
	fields := make(map[string]*schema_pb.Value, n)
	for i := 0; i < n; i++ {
		e, ok := rv.FieldEntryAt(i)
		if !ok {
			continue
		}
		v, ok := e.Value()
		if !ok {
			continue
		}
		fields[e.Key()] = valueFromWire(v)
	}
	return &schema_pb.RecordValue{Fields: fields}
}

func valueFromWire(v mq_schemawire.Value) *schema_pb.Value {
	switch v.Kind() {
	case mq_schemawire.ValueKindBoolValue:
		if x, ok := v.BoolValue(); ok {
			return &schema_pb.Value{Kind: &schema_pb.Value_BoolValue{BoolValue: x}}
		}
	case mq_schemawire.ValueKindInt32Value:
		if x, ok := v.Int32Value(); ok {
			return &schema_pb.Value{Kind: &schema_pb.Value_Int32Value{Int32Value: x}}
		}
	case mq_schemawire.ValueKindInt64Value:
		if x, ok := v.Int64Value(); ok {
			return &schema_pb.Value{Kind: &schema_pb.Value_Int64Value{Int64Value: x}}
		}
	case mq_schemawire.ValueKindFloatValue:
		if x, ok := v.FloatValue(); ok {
			return &schema_pb.Value{Kind: &schema_pb.Value_FloatValue{FloatValue: x}}
		}
	case mq_schemawire.ValueKindDoubleValue:
		if x, ok := v.DoubleValue(); ok {
			return &schema_pb.Value{Kind: &schema_pb.Value_DoubleValue{DoubleValue: x}}
		}
	case mq_schemawire.ValueKindBytesValue:
		if x, ok := v.BytesValue(); ok {
			return &schema_pb.Value{Kind: &schema_pb.Value_BytesValue{BytesValue: x}}
		}
	case mq_schemawire.ValueKindStringValue:
		if x, ok := v.StringValue(); ok {
			return &schema_pb.Value{Kind: &schema_pb.Value_StringValue{StringValue: x}}
		}
	case mq_schemawire.ValueKindTimestampValue:
		if x, ok := v.TimestampValue(); ok {
			return &schema_pb.Value{Kind: &schema_pb.Value_TimestampValue{TimestampValue: &schema_pb.TimestampValue{
				TimestampMicros: x.TimestampMicros(),
				IsUtc:           x.IsUtc(),
			}}}
		}
	case mq_schemawire.ValueKindDateValue:
		if x, ok := v.DateValue(); ok {
			return &schema_pb.Value{Kind: &schema_pb.Value_DateValue{DateValue: &schema_pb.DateValue{
				DaysSinceEpoch: x.DaysSinceEpoch(),
			}}}
		}
	case mq_schemawire.ValueKindDecimalValue:
		if x, ok := v.DecimalValue(); ok {
			return &schema_pb.Value{Kind: &schema_pb.Value_DecimalValue{DecimalValue: &schema_pb.DecimalValue{
				Value:     x.Value(),
				Precision: x.Precision(),
				Scale:     x.Scale(),
			}}}
		}
	case mq_schemawire.ValueKindTimeValue:
		if x, ok := v.TimeValue(); ok {
			return &schema_pb.Value{Kind: &schema_pb.Value_TimeValue{TimeValue: &schema_pb.TimeValue{
				TimeMicros: x.TimeMicros(),
			}}}
		}
	case mq_schemawire.ValueKindListValue:
		if x, ok := v.ListValue(); ok {
			n := x.ValuesLen()
			values := make([]*schema_pb.Value, 0, n)
			for i := 0; i < n; i++ {
				if ev, ok := x.ValueAt(i); ok {
					values = append(values, valueFromWire(ev))
				}
			}
			return &schema_pb.Value{Kind: &schema_pb.Value_ListValue{ListValue: &schema_pb.ListValue{Values: values}}}
		}
	case mq_schemawire.ValueKindRecordValue:
		if x, ok := v.RecordValue(); ok {
			return &schema_pb.Value{Kind: &schema_pb.Value_RecordValue{RecordValue: recordValueFromWire(x)}}
		}
	}
	return nil
}
