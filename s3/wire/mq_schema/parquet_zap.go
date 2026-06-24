// Code generated from mq_schema.proto; DO NOT EDIT.

package mq_schemawire

import (
	zap "github.com/zap-proto/go"
)

// Parquet logical-type value messages.

// --- TimestampValue ---

const (
	timestampValueTimestampMicrosOff = 0
	timestampValueIsUtcOff           = 8
	timestampValueSize               = 16
)

// TimestampValue is a zero-copy view into a ZAP-encoded TimestampValue message.
type TimestampValue struct{ o zap.Object }

// WrapTimestampValue parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapTimestampValue(b []byte) (TimestampValue, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return TimestampValue{}, err
	}
	return TimestampValue{o: m.Root()}, nil
}

// TimestampMicros reads the timestamp_micros field (proto field 1, int64) —
// microseconds since Unix epoch (UTC).
func (t TimestampValue) TimestampMicros() int64 {
	return t.o.Int64(timestampValueTimestampMicrosOff)
}

// IsUtc reads the is_utc field (proto field 2, bool).
func (t TimestampValue) IsUtc() bool { return t.o.Bool(timestampValueIsUtcOff) }

// TimestampValueInput collects the field values for NewTimestampValue.
type TimestampValueInput struct {
	TimestampMicros int64
	IsUtc           bool
}

// NewTimestampValue builds a ZAP-encoded TimestampValue message from in and
// returns the bytes.
func NewTimestampValue(in TimestampValueInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(timestampValueSize)
	ob.SetInt64(timestampValueTimestampMicrosOff, in.TimestampMicros)
	ob.SetBool(timestampValueIsUtcOff, in.IsUtc)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- DateValue ---

const (
	dateValueDaysSinceEpochOff = 0
	dateValueSize              = 4
)

// DateValue is a zero-copy view into a ZAP-encoded DateValue message.
type DateValue struct{ o zap.Object }

// WrapDateValue parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapDateValue(b []byte) (DateValue, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return DateValue{}, err
	}
	return DateValue{o: m.Root()}, nil
}

// DaysSinceEpoch reads the days_since_epoch field (proto field 1, int32) —
// days since Unix epoch (1970-01-01).
func (t DateValue) DaysSinceEpoch() int32 { return t.o.Int32(dateValueDaysSinceEpochOff) }

// DateValueInput collects the field values for NewDateValue.
type DateValueInput struct {
	DaysSinceEpoch int32
}

// NewDateValue builds a ZAP-encoded DateValue message from in and returns the
// bytes.
func NewDateValue(in DateValueInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(dateValueSize)
	ob.SetInt32(dateValueDaysSinceEpochOff, in.DaysSinceEpoch)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- DecimalValue ---

const (
	decimalValueValueOff     = 0
	decimalValuePrecisionOff = 8
	decimalValueScaleOff     = 12
	decimalValueSize         = 16
)

// DecimalValue is a zero-copy view into a ZAP-encoded DecimalValue message.
type DecimalValue struct{ o zap.Object }

// WrapDecimalValue parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapDecimalValue(b []byte) (DecimalValue, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return DecimalValue{}, err
	}
	return DecimalValue{o: m.Root()}, nil
}

// Value reads the value field (proto field 1, bytes) — arbitrary precision
// decimal as bytes.
func (t DecimalValue) Value() []byte { return t.o.Bytes(decimalValueValueOff) }

// Precision reads the precision field (proto field 2, int32) — total number of
// digits.
func (t DecimalValue) Precision() int32 { return t.o.Int32(decimalValuePrecisionOff) }

// Scale reads the scale field (proto field 3, int32) — number of digits after
// the decimal point.
func (t DecimalValue) Scale() int32 { return t.o.Int32(decimalValueScaleOff) }

// DecimalValueInput collects the field values for NewDecimalValue.
type DecimalValueInput struct {
	Value     []byte
	Precision int32
	Scale     int32
}

// NewDecimalValue builds a ZAP-encoded DecimalValue message from in and returns
// the bytes.
func NewDecimalValue(in DecimalValueInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(decimalValueSize)
	ob.SetBytes(decimalValueValueOff, in.Value)
	ob.SetInt32(decimalValuePrecisionOff, in.Precision)
	ob.SetInt32(decimalValueScaleOff, in.Scale)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- TimeValue ---

const (
	timeValueTimeMicrosOff = 0
	timeValueSize          = 8
)

// TimeValue is a zero-copy view into a ZAP-encoded TimeValue message.
type TimeValue struct{ o zap.Object }

// WrapTimeValue parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapTimeValue(b []byte) (TimeValue, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return TimeValue{}, err
	}
	return TimeValue{o: m.Root()}, nil
}

// TimeMicros reads the time_micros field (proto field 1, int64) — microseconds
// since midnight.
func (t TimeValue) TimeMicros() int64 { return t.o.Int64(timeValueTimeMicrosOff) }

// TimeValueInput collects the field values for NewTimeValue.
type TimeValueInput struct {
	TimeMicros int64
}

// NewTimeValue builds a ZAP-encoded TimeValue message from in and returns the
// bytes.
func NewTimeValue(in TimeValueInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(timeValueSize)
	ob.SetInt64(timeValueTimeMicrosOff, in.TimeMicros)
	ob.FinishAsRoot()
	return b.Finish()
}
