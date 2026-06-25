// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// convert.go is the leaf proto<->wire converter set shared by the HanzoMessaging
// RPCs (the per-RPC request builders / response decoders are in rpc.go, the
// streaming frame converters in stream.go). It maps the nested messages that
// recur across the surface — schema_pb.Topic / Partition / PartitionOffset /
// RecordType, mq_pb.BrokerPartitionAssignment / TopicRetention / DataMessage /
// ControlMessage / TopicPublisher / TopicSubscriber / Offset+TimestampRangeInfo,
// and filer_pb.LogEntry — onto their mq_brokerwire / mq_schemawire / filerwire
// codecs. These are the ONLY place that mapping lives; the ZAP-backed
// mq_pb.HanzoMessagingClient (client.go) reuses them.

package mq

import (
	filer_pb "github.com/hanzoai/s3/s3/pb/filer_pb"
	mq_pb "github.com/hanzoai/s3/s3/pb/mq_pb"
	schema_pb "github.com/hanzoai/s3/s3/pb/schema_pb"
	filerwire "github.com/hanzoai/s3/s3/wire/filer"
	mq_brokerwire "github.com/hanzoai/s3/s3/wire/mq_broker"
	mq_schemawire "github.com/hanzoai/s3/s3/wire/mq_schema"
)

// --- schema_pb.Topic ---

func topicToWire(t *schema_pb.Topic) []byte {
	if t == nil {
		return nil
	}
	return mq_schemawire.NewTopic(mq_schemawire.TopicInput{Namespace: t.Namespace, Name: t.Name})
}

func topicFromWire(b []byte) *schema_pb.Topic {
	if len(b) == 0 {
		return nil
	}
	v, err := mq_schemawire.WrapTopic(b)
	if err != nil {
		return nil
	}
	return &schema_pb.Topic{Namespace: v.Namespace(), Name: v.Name()}
}

// --- schema_pb.Partition ---

func partitionToWire(p *schema_pb.Partition) []byte {
	if p == nil {
		return nil
	}
	return mq_schemawire.NewPartition(mq_schemawire.PartitionInput{
		RingSize: p.RingSize, RangeStart: p.RangeStart, RangeStop: p.RangeStop, UnixTimeNs: p.UnixTimeNs,
	})
}

func partitionFromWire(b []byte) *schema_pb.Partition {
	if len(b) == 0 {
		return nil
	}
	v, err := mq_schemawire.WrapPartition(b)
	if err != nil {
		return nil
	}
	return partitionFromView(v)
}

func partitionFromView(v mq_schemawire.Partition) *schema_pb.Partition {
	return &schema_pb.Partition{
		RingSize: v.RingSize(), RangeStart: v.RangeStart(), RangeStop: v.RangeStop(), UnixTimeNs: v.UnixTimeNs(),
	}
}

// --- schema_pb.PartitionOffset ---

func partitionOffsetToWire(o *schema_pb.PartitionOffset) []byte {
	if o == nil {
		return nil
	}
	return mq_schemawire.NewPartitionOffset(mq_schemawire.PartitionOffsetInput{
		Partition: partitionToWire(o.Partition), StartTsNs: o.StartTsNs, StartOffset: o.StartOffset,
	})
}

func partitionOffsetFromWire(b []byte) *schema_pb.PartitionOffset {
	if len(b) == 0 {
		return nil
	}
	v, err := mq_schemawire.WrapPartitionOffset(b)
	if err != nil {
		return nil
	}
	o := &schema_pb.PartitionOffset{StartTsNs: v.StartTsNs(), StartOffset: v.StartOffset()}
	if p, ok := v.Partition(); ok {
		o.Partition = partitionFromView(p)
	}
	return o
}

// --- schema_pb.RecordType (recursive: Field -> Type -> RecordType/ListType) ---

func recordTypeToWire(rt *schema_pb.RecordType) []byte {
	if rt == nil {
		return nil
	}
	fields := make([][]byte, len(rt.Fields))
	for i, f := range rt.Fields {
		fields[i] = mq_schemawire.NewField(mq_schemawire.FieldInput{
			Name:       f.Name,
			FieldIndex: f.FieldIndex,
			Type:       typeToWire(f.Type),
			IsRepeated: f.IsRepeated,
			IsRequired: f.IsRequired,
		})
	}
	return mq_schemawire.NewRecordType(mq_schemawire.RecordTypeInput{Fields: fields})
}

func typeToWire(t *schema_pb.Type) []byte {
	if t == nil {
		return nil
	}
	in := mq_schemawire.TypeInput{Kind: mq_schemawire.TypeKindNone}
	switch k := t.Kind.(type) {
	case *schema_pb.Type_ScalarType:
		in.Kind = mq_schemawire.TypeKindScalarType
		in.ScalarType = uint32(k.ScalarType)
	case *schema_pb.Type_RecordType:
		in.Kind = mq_schemawire.TypeKindRecordType
		in.RecordType = recordTypeToWire(k.RecordType)
	case *schema_pb.Type_ListType:
		in.Kind = mq_schemawire.TypeKindListType
		in.ListType = listTypeToWire(k.ListType)
	}
	return mq_schemawire.NewType(in)
}

func listTypeToWire(lt *schema_pb.ListType) []byte {
	if lt == nil {
		return nil
	}
	return mq_schemawire.NewListType(mq_schemawire.ListTypeInput{ElementType: typeToWire(lt.ElementType)})
}

func recordTypeFromWire(b []byte) *schema_pb.RecordType {
	if len(b) == 0 {
		return nil
	}
	v, err := mq_schemawire.WrapRecordType(b)
	if err != nil {
		return nil
	}
	return recordTypeFromView(v)
}

func recordTypeFromView(v mq_schemawire.RecordType) *schema_pb.RecordType {
	rt := &schema_pb.RecordType{}
	for i := 0; i < v.FieldsLen(); i++ {
		f, ok := v.FieldAt(i)
		if !ok {
			continue
		}
		field := &schema_pb.Field{
			Name:       f.Name(),
			FieldIndex: f.FieldIndex(),
			IsRepeated: f.IsRepeated(),
			IsRequired: f.IsRequired(),
		}
		if ty, ok := f.Type(); ok {
			field.Type = typeFromView(ty)
		}
		rt.Fields = append(rt.Fields, field)
	}
	return rt
}

func typeFromView(v mq_schemawire.Type) *schema_pb.Type {
	switch v.Kind() {
	case mq_schemawire.TypeKindScalarType:
		if s, ok := v.ScalarType(); ok {
			return &schema_pb.Type{Kind: &schema_pb.Type_ScalarType{ScalarType: schema_pb.ScalarType(s)}}
		}
	case mq_schemawire.TypeKindRecordType:
		if rt, ok := v.RecordType(); ok {
			return &schema_pb.Type{Kind: &schema_pb.Type_RecordType{RecordType: recordTypeFromView(rt)}}
		}
	case mq_schemawire.TypeKindListType:
		if lt, ok := v.ListType(); ok {
			return &schema_pb.Type{Kind: &schema_pb.Type_ListType{ListType: listTypeFromView(lt)}}
		}
	}
	return nil
}

func listTypeFromView(v mq_schemawire.ListType) *schema_pb.ListType {
	lt := &schema_pb.ListType{}
	if et, ok := v.ElementType(); ok {
		lt.ElementType = typeFromView(et)
	}
	return lt
}

// --- mq_pb.TopicRetention ---

func retentionToWire(r *mq_pb.TopicRetention) []byte {
	if r == nil {
		return nil
	}
	return mq_brokerwire.NewTopicRetention(mq_brokerwire.TopicRetentionInput{
		RetentionSeconds: r.RetentionSeconds, Enabled: r.Enabled,
	})
}

func retentionFromView(v mq_brokerwire.TopicRetention) *mq_pb.TopicRetention {
	return &mq_pb.TopicRetention{RetentionSeconds: v.RetentionSeconds(), Enabled: v.Enabled()}
}

// --- mq_pb.BrokerPartitionAssignment ---

func assignmentToWire(a *mq_pb.BrokerPartitionAssignment) []byte {
	if a == nil {
		return nil
	}
	return mq_brokerwire.NewBrokerPartitionAssignment(mq_brokerwire.BrokerPartitionAssignmentInput{
		Partition: partitionToWire(a.Partition), LeaderBroker: a.LeaderBroker, FollowerBroker: a.FollowerBroker,
	})
}

func assignmentsToWire(as []*mq_pb.BrokerPartitionAssignment) [][]byte {
	if len(as) == 0 {
		return nil
	}
	out := make([][]byte, len(as))
	for i, a := range as {
		out[i] = assignmentToWire(a)
	}
	return out
}

func assignmentFromView(v mq_brokerwire.BrokerPartitionAssignment) *mq_pb.BrokerPartitionAssignment {
	return &mq_pb.BrokerPartitionAssignment{
		Partition:      partitionFromWire(v.Partition()),
		LeaderBroker:   v.LeaderBroker(),
		FollowerBroker: v.FollowerBroker(),
	}
}

// --- mq_pb.ControlMessage ---

func controlToWire(c *mq_pb.ControlMessage) []byte {
	if c == nil {
		return nil
	}
	return mq_brokerwire.NewControlMessage(mq_brokerwire.ControlMessageInput{
		IsClose: c.IsClose, PublisherName: c.PublisherName,
	})
}

func controlFromView(v mq_brokerwire.ControlMessage) *mq_pb.ControlMessage {
	return &mq_pb.ControlMessage{IsClose: v.IsClose(), PublisherName: v.PublisherName()}
}

// --- mq_pb.DataMessage ---

func dataMessageToWire(d *mq_pb.DataMessage) []byte {
	if d == nil {
		return nil
	}
	return mq_brokerwire.NewDataMessage(mq_brokerwire.DataMessageInput{
		Key: d.Key, Value: d.Value, TsNs: d.TsNs, Ctrl: controlToWire(d.Ctrl),
	})
}

func dataMessageFromView(v mq_brokerwire.DataMessage) *mq_pb.DataMessage {
	d := &mq_pb.DataMessage{Key: v.Key(), Value: v.Value(), TsNs: v.TsNs()}
	if c, ok := v.Ctrl(); ok {
		d.Ctrl = controlFromView(c)
	}
	return d
}

// --- mq_pb.TopicPublisher / TopicSubscriber ---

func publisherFromView(v mq_brokerwire.TopicPublisher) *mq_pb.TopicPublisher {
	return &mq_pb.TopicPublisher{
		PublisherName:       v.PublisherName(),
		ClientId:            v.ClientId(),
		Partition:           partitionFromWire(v.Partition()),
		ConnectTimeNs:       v.ConnectTimeNs(),
		LastSeenTimeNs:      v.LastSeenTimeNs(),
		Broker:              v.Broker(),
		IsActive:            v.IsActive(),
		LastPublishedOffset: v.LastPublishedOffset(),
		LastAckedOffset:     v.LastAckedOffset(),
	}
}

func subscriberFromView(v mq_brokerwire.TopicSubscriber) *mq_pb.TopicSubscriber {
	return &mq_pb.TopicSubscriber{
		ConsumerGroup:      v.ConsumerGroup(),
		ConsumerId:         v.ConsumerId(),
		ClientId:           v.ClientId(),
		Partition:          partitionFromWire(v.Partition()),
		ConnectTimeNs:      v.ConnectTimeNs(),
		LastSeenTimeNs:     v.LastSeenTimeNs(),
		Broker:             v.Broker(),
		IsActive:           v.IsActive(),
		CurrentOffset:      v.CurrentOffset(),
		LastReceivedOffset: v.LastReceivedOffset(),
	}
}

// --- mq_pb.BrokerStats (PublisherToPubBalancer payload, server direction) ---

func brokerStatsFromView(v mq_brokerwire.BrokerStats) *mq_pb.BrokerStats {
	bs := &mq_pb.BrokerStats{
		CpuUsagePercent: v.CpuUsagePercent(),
		Stats:           make(map[string]*mq_pb.TopicPartitionStats, v.StatsLen()),
	}
	for i := 0; i < v.StatsLen(); i++ {
		e := v.StatsAt(i)
		bs.Stats[e.Key] = topicPartitionStatsFromWire(e.Value)
	}
	return bs
}

func topicPartitionStatsFromWire(b []byte) *mq_pb.TopicPartitionStats {
	if len(b) == 0 {
		return nil
	}
	v, err := mq_brokerwire.WrapTopicPartitionStats(b)
	if err != nil {
		return nil
	}
	return &mq_pb.TopicPartitionStats{
		Topic:           topicFromWire(v.Topic()),
		Partition:       partitionFromWire(v.Partition()),
		PublisherCount:  v.PublisherCount(),
		SubscriberCount: v.SubscriberCount(),
		Follower:        v.Follower(),
	}
}

// --- mq_pb.Offset/TimestampRangeInfo ---

func offsetRangeFromView(v mq_brokerwire.OffsetRangeInfo) *mq_pb.OffsetRangeInfo {
	return &mq_pb.OffsetRangeInfo{
		EarliestOffset: v.EarliestOffset(), LatestOffset: v.LatestOffset(), HighWaterMark: v.HighWaterMark(),
	}
}

func timestampRangeFromView(v mq_brokerwire.TimestampRangeInfo) *mq_pb.TimestampRangeInfo {
	return &mq_pb.TimestampRangeInfo{
		EarliestTimestampNs: v.EarliestTimestampNs(), LatestTimestampNs: v.LatestTimestampNs(),
	}
}

// --- filer_pb.LogEntry (GetUnflushedMessages payload) ---

func logEntryToWire(e *filer_pb.LogEntry) []byte {
	if e == nil {
		return nil
	}
	return filerwire.NewLogEntry(filerwire.LogEntryInput{
		TsNs: e.TsNs, PartitionKeyHash: e.PartitionKeyHash, Data: e.Data, Key: e.Key, Offset: e.Offset,
	})
}

func logEntryFromWire(b []byte) *filer_pb.LogEntry {
	if len(b) == 0 {
		return nil
	}
	v, err := filerwire.WrapLogEntry(b)
	if err != nil {
		return nil
	}
	return &filer_pb.LogEntry{
		TsNs: v.TsNs(), PartitionKeyHash: v.PartitionKeyHash(), Data: v.Data(), Key: v.Key(), Offset: v.Offset(),
	}
}

// topicsFromView decodes the repeated Topic list carried by ListTopicsResponse.
func topicsFromView(n int, at func(int) []byte) []*schema_pb.Topic {
	if n == 0 {
		return nil
	}
	out := make([]*schema_pb.Topic, 0, n)
	for i := 0; i < n; i++ {
		if t := topicFromWire(at(i)); t != nil {
			out = append(out, t)
		}
	}
	return out
}
