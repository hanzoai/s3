// Code generated from mq_schema.proto; DO NOT EDIT.

package mq_schemawire

// Enums from mq_schema.proto. Each proto enum becomes a uint32 const block
// with the proto's exact numeric values (gaps preserved); the corresponding
// field is encoded as a 4-byte uint32 scalar.

// OffsetType selects how a PartitionOffset positions a consumer.
const (
	OffsetTypeResumeOrEarliest uint32 = 0
	OffsetTypeResetToEarliest  uint32 = 5
	OffsetTypeExactTsNs        uint32 = 10
	OffsetTypeResetToLatest    uint32 = 15
	OffsetTypeResumeOrLatest   uint32 = 20
	// Offset-based positioning.
	OffsetTypeExactOffset   uint32 = 25
	OffsetTypeResetToOffset uint32 = 30
)

// ScalarType is the leaf type of a schema field.
const (
	ScalarTypeBool   uint32 = 0
	ScalarTypeInt32  uint32 = 1
	ScalarTypeInt64  uint32 = 3
	ScalarTypeFloat  uint32 = 4
	ScalarTypeDouble uint32 = 5
	ScalarTypeBytes  uint32 = 6
	ScalarTypeString uint32 = 7
	// Parquet logical types for analytics.
	ScalarTypeTimestamp uint32 = 8  // UTC timestamp (microseconds since epoch)
	ScalarTypeDate      uint32 = 9  // Date (days since epoch)
	ScalarTypeDecimal   uint32 = 10 // Arbitrary precision decimal
	ScalarTypeTime      uint32 = 11 // Time of day (microseconds)
)
