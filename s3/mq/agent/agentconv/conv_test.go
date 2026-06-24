package agentconv

import (
	"testing"

	"github.com/hanzoai/s3/s3/pb/schema_pb"
)

// TestRecordValueRoundTrip checks that a schema_pb.RecordValue carrying a
// representative spread of Value variants (including nested list/record)
// survives the schema_pb -> ZAP wire -> schema_pb round trip the agent does
// across its publish/subscribe boundary.
func TestRecordValueRoundTrip(t *testing.T) {
	in := &schema_pb.RecordValue{Fields: map[string]*schema_pb.Value{
		"b":   {Kind: &schema_pb.Value_BoolValue{BoolValue: true}},
		"i32": {Kind: &schema_pb.Value_Int32Value{Int32Value: -7}},
		"i64": {Kind: &schema_pb.Value_Int64Value{Int64Value: 1 << 40}},
		"f":   {Kind: &schema_pb.Value_FloatValue{FloatValue: 3.5}},
		"d":   {Kind: &schema_pb.Value_DoubleValue{DoubleValue: 2.25}},
		"by":  {Kind: &schema_pb.Value_BytesValue{BytesValue: []byte{1, 2, 3}}},
		"s":   {Kind: &schema_pb.Value_StringValue{StringValue: "hello"}},
		"ts":  {Kind: &schema_pb.Value_TimestampValue{TimestampValue: &schema_pb.TimestampValue{TimestampMicros: 99, IsUtc: true}}},
		"list": {Kind: &schema_pb.Value_ListValue{ListValue: &schema_pb.ListValue{Values: []*schema_pb.Value{
			{Kind: &schema_pb.Value_Int64Value{Int64Value: 1}},
			{Kind: &schema_pb.Value_StringValue{StringValue: "x"}},
		}}}},
		"rec": {Kind: &schema_pb.Value_RecordValue{RecordValue: &schema_pb.RecordValue{Fields: map[string]*schema_pb.Value{
			"inner": {Kind: &schema_pb.Value_Int32Value{Int32Value: 42}},
		}}}},
	}}

	out := RecordValueFromWire(RecordValueToWire(in))
	if out == nil {
		t.Fatal("round trip returned nil")
	}
	if len(out.Fields) != len(in.Fields) {
		t.Fatalf("field count: got %d want %d", len(out.Fields), len(in.Fields))
	}
	if !out.Fields["b"].GetBoolValue() {
		t.Error("bool lost")
	}
	if got := out.Fields["i64"].GetInt64Value(); got != 1<<40 {
		t.Errorf("i64: got %d", got)
	}
	if got := out.Fields["s"].GetStringValue(); got != "hello" {
		t.Errorf("string: got %q", got)
	}
	if got := out.Fields["ts"].GetTimestampValue().GetTimestampMicros(); got != 99 {
		t.Errorf("timestamp: got %d", got)
	}
	if got := out.Fields["list"].GetListValue().GetValues(); len(got) != 2 || got[1].GetStringValue() != "x" {
		t.Errorf("list: got %+v", got)
	}
	if got := out.Fields["rec"].GetRecordValue().GetFields()["inner"].GetInt32Value(); got != 42 {
		t.Errorf("nested record: got %d", got)
	}
}

// TestRecordTypeRoundTrip checks the schema (RecordType/Field/Type) tree
// survives the same boundary, including a nested record type and list type.
func TestRecordTypeRoundTrip(t *testing.T) {
	in := &schema_pb.RecordType{Fields: []*schema_pb.Field{
		{Name: "id", FieldIndex: 0, IsRequired: true, Type: &schema_pb.Type{Kind: &schema_pb.Type_ScalarType{ScalarType: schema_pb.ScalarType_INT64}}},
		{Name: "tags", FieldIndex: 1, IsRepeated: true, Type: &schema_pb.Type{Kind: &schema_pb.Type_ListType{ListType: &schema_pb.ListType{
			ElementType: &schema_pb.Type{Kind: &schema_pb.Type_ScalarType{ScalarType: schema_pb.ScalarType_STRING}},
		}}}},
		{Name: "nested", FieldIndex: 2, Type: &schema_pb.Type{Kind: &schema_pb.Type_RecordType{RecordType: &schema_pb.RecordType{Fields: []*schema_pb.Field{
			{Name: "x", FieldIndex: 0, Type: &schema_pb.Type{Kind: &schema_pb.Type_ScalarType{ScalarType: schema_pb.ScalarType_BOOL}}},
		}}}}},
	}}

	out := RecordTypeFromWire(RecordTypeToWire(in))
	if out == nil || len(out.Fields) != 3 {
		t.Fatalf("record type round trip: %+v", out)
	}
	if out.Fields[0].Name != "id" || out.Fields[0].GetType().GetScalarType() != schema_pb.ScalarType_INT64 {
		t.Errorf("scalar field: %+v", out.Fields[0])
	}
	if !out.Fields[1].IsRepeated || out.Fields[1].GetType().GetListType().GetElementType().GetScalarType() != schema_pb.ScalarType_STRING {
		t.Errorf("list field: %+v", out.Fields[1])
	}
	if out.Fields[2].GetType().GetRecordType().GetFields()[0].Name != "x" {
		t.Errorf("nested record field: %+v", out.Fields[2])
	}
}

// TestTopicRoundTrip covers the simple Topic conversion.
func TestTopicRoundTrip(t *testing.T) {
	in := &schema_pb.Topic{Namespace: "ns", Name: "n"}
	out := TopicFromWire(TopicToWire(in))
	if out == nil || out.Namespace != "ns" || out.Name != "n" {
		t.Fatalf("topic round trip: %+v", out)
	}
	if TopicToWire(nil) != nil || TopicFromWire(nil) != nil {
		t.Error("nil topic should map to nil")
	}
}
