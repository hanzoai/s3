// Code generated from mq_schema.proto; DO NOT EDIT.

// Package mq_schemawire holds the native ZAP schemas for mq_schema.proto —
// zero-copy, protobuf-free, no _pb. It is the types-only foundation for the
// message-queue layer: Topic / Partition / Offset positioning, the recursive
// schema definition (RecordType <-> Field <-> Type <-> ListType), and the
// recursive value union (Value <-> ListValue <-> RecordValue) with the Parquet
// logical-type values.
//
// mq_schema.proto declares no service, so this package emits no RPC skeleton
// (no ordinals, Channel, Client, Handler, or Dispatch). The message-queue
// services that DO carry RPCs — mq_agent.proto and mq_broker.proto — import
// these schemas from their own wire packages (mq_agentwire, mq_brokerwire)
// rather than redefining them.
//
// Recursion / interning: each recursive message has exactly one builder
// (New<Msg>) and one reader (Wrap<Msg>); a nested message is embedded as its
// own standalone ZAP buffer (a tail-pointer bytes slot for a singular field, a
// variable-element list entry for a repeated field, or the single oneof
// payload slot for a union variant). The recursion therefore bottoms out
// through ordinary build-time calls and read-time re-parse, with no duplicated
// element schema — Value, ListValue, RecordValue, Type and ListType are each
// defined once and reused everywhere they recur.
package mq_schemawire
