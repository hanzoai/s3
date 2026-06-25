// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// server_convert.go is the server-side mirror of rpc.go: the converters that
// decode an inbound mq_brokerwire request VIEW into the *mq_pb.<Rpc>Request the
// existing MessageQueueBroker engine expects, and encode the *mq_pb.<Rpc>Response
// it returns into the matching mq_brokerwire.<Rpc>ResponseInput. rpc.go owns the
// client direction (*mq_pb -> wire bytes, wire bytes -> *mq_pb); this file owns
// the server direction (request view -> *mq_pb, *mq_pb response -> ResponseInput).
// Both reuse the SAME leaf converters in convert.go — the proto<->wire mapping
// lives in exactly one place. Naming: <Rpc>ReqFromView(view) *mq_pb.<Rpc>Request
// and <Rpc>RespToInput(*mq_pb.<Rpc>Response) mq_brokerwire.<Rpc>ResponseInput.

package mq

import (
	mq_pb "github.com/hanzoai/s3/s3/pb/mq_pb"
	mq_brokerwire "github.com/hanzoai/s3/s3/wire/mq_broker"
)

// --- shared leaf: request-view child buffers -> *mq_pb ---

// retentionFromWire wraps a TopicRetention child buffer (nil-safe) into proto.
func retentionFromWire(b []byte) *mq_pb.TopicRetention {
	if len(b) == 0 {
		return nil
	}
	v, err := mq_brokerwire.WrapTopicRetention(b)
	if err != nil {
		return nil
	}
	return retentionFromView(v)
}

// assignmentFromWire wraps a BrokerPartitionAssignment child buffer into proto.
func assignmentFromWire(b []byte) *mq_pb.BrokerPartitionAssignment {
	if len(b) == 0 {
		return nil
	}
	v, err := mq_brokerwire.WrapBrokerPartitionAssignment(b)
	if err != nil {
		return nil
	}
	return assignmentFromView(v)
}

// --- proto -> ResponseInput leaf builders (response direction) ---

func publisherToWire(p *mq_pb.TopicPublisher) []byte {
	if p == nil {
		return nil
	}
	return mq_brokerwire.NewTopicPublisher(mq_brokerwire.TopicPublisherInput{
		PublisherName:       p.PublisherName,
		ClientId:            p.ClientId,
		Partition:           partitionToWire(p.Partition),
		ConnectTimeNs:       p.ConnectTimeNs,
		LastSeenTimeNs:      p.LastSeenTimeNs,
		Broker:              p.Broker,
		IsActive:            p.IsActive,
		LastPublishedOffset: p.LastPublishedOffset,
		LastAckedOffset:     p.LastAckedOffset,
	})
}

func subscriberToWire(s *mq_pb.TopicSubscriber) []byte {
	if s == nil {
		return nil
	}
	return mq_brokerwire.NewTopicSubscriber(mq_brokerwire.TopicSubscriberInput{
		ConsumerGroup:      s.ConsumerGroup,
		ConsumerId:         s.ConsumerId,
		ClientId:           s.ClientId,
		Partition:          partitionToWire(s.Partition),
		ConnectTimeNs:      s.ConnectTimeNs,
		LastSeenTimeNs:     s.LastSeenTimeNs,
		Broker:             s.Broker,
		IsActive:           s.IsActive,
		CurrentOffset:      s.CurrentOffset,
		LastReceivedOffset: s.LastReceivedOffset,
	})
}

func offsetRangeToWire(o *mq_pb.OffsetRangeInfo) []byte {
	if o == nil {
		return nil
	}
	return mq_brokerwire.NewOffsetRangeInfo(mq_brokerwire.OffsetRangeInfoInput{
		EarliestOffset: o.EarliestOffset, LatestOffset: o.LatestOffset, HighWaterMark: o.HighWaterMark,
	})
}

func timestampRangeToWire(t *mq_pb.TimestampRangeInfo) []byte {
	if t == nil {
		return nil
	}
	return mq_brokerwire.NewTimestampRangeInfo(mq_brokerwire.TimestampRangeInfoInput{
		EarliestTimestampNs: t.EarliestTimestampNs, LatestTimestampNs: t.LatestTimestampNs,
	})
}

// --- FindBrokerLeader ---

func FindBrokerLeaderReqFromView(v mq_brokerwire.FindBrokerLeaderRequest) *mq_pb.FindBrokerLeaderRequest {
	return &mq_pb.FindBrokerLeaderRequest{FilerGroup: v.FilerGroup()}
}
func FindBrokerLeaderRespToInput(r *mq_pb.FindBrokerLeaderResponse) mq_brokerwire.FindBrokerLeaderResponseInput {
	if r == nil {
		return mq_brokerwire.FindBrokerLeaderResponseInput{}
	}
	return mq_brokerwire.FindBrokerLeaderResponseInput{Broker: r.Broker}
}

// --- BalanceTopics ---

func BalanceTopicsReqFromView(v mq_brokerwire.BalanceTopicsRequest) *mq_pb.BalanceTopicsRequest {
	return &mq_pb.BalanceTopicsRequest{}
}
func BalanceTopicsRespToInput(r *mq_pb.BalanceTopicsResponse) mq_brokerwire.BalanceTopicsResponseInput {
	return mq_brokerwire.BalanceTopicsResponseInput{}
}

// --- ListTopics ---

func ListTopicsReqFromView(v mq_brokerwire.ListTopicsRequest) *mq_pb.ListTopicsRequest {
	return &mq_pb.ListTopicsRequest{}
}
func ListTopicsRespToInput(r *mq_pb.ListTopicsResponse) mq_brokerwire.ListTopicsResponseInput {
	if r == nil {
		return mq_brokerwire.ListTopicsResponseInput{}
	}
	topics := make([][]byte, 0, len(r.Topics))
	for _, t := range r.Topics {
		topics = append(topics, topicToWire(t))
	}
	return mq_brokerwire.ListTopicsResponseInput{Topics: topics}
}

// --- TopicExists ---

func TopicExistsReqFromView(v mq_brokerwire.TopicExistsRequest) *mq_pb.TopicExistsRequest {
	return &mq_pb.TopicExistsRequest{Topic: topicFromWire(v.Topic())}
}
func TopicExistsRespToInput(r *mq_pb.TopicExistsResponse) mq_brokerwire.TopicExistsResponseInput {
	if r == nil {
		return mq_brokerwire.TopicExistsResponseInput{}
	}
	return mq_brokerwire.TopicExistsResponseInput{Exists: r.Exists}
}

// --- ConfigureTopic ---

func ConfigureTopicReqFromView(v mq_brokerwire.ConfigureTopicRequest) *mq_pb.ConfigureTopicRequest {
	req := &mq_pb.ConfigureTopicRequest{
		Topic:             topicFromWire(v.Topic()),
		PartitionCount:    v.PartitionCount(),
		MessageRecordType: recordTypeFromWire(v.MessageRecordType()),
		KeyColumns:        stringsFromView(v.KeyColumnsLen(), v.KeyColumnAt),
		SchemaFormat:      v.SchemaFormat(),
	}
	if rt, ok := v.Retention(); ok {
		req.Retention = retentionFromView(rt)
	}
	return req
}
func ConfigureTopicRespToInput(r *mq_pb.ConfigureTopicResponse) mq_brokerwire.ConfigureTopicResponseInput {
	if r == nil {
		return mq_brokerwire.ConfigureTopicResponseInput{}
	}
	return mq_brokerwire.ConfigureTopicResponseInput{
		BrokerPartitionAssignments: assignmentsToWire(r.BrokerPartitionAssignments),
		Retention:                  retentionToWire(r.Retention),
		MessageRecordType:          recordTypeToWire(r.MessageRecordType),
		KeyColumns:                 r.KeyColumns,
		SchemaFormat:               r.SchemaFormat,
	}
}

// --- LookupTopicBrokers ---

func LookupTopicBrokersReqFromView(v mq_brokerwire.LookupTopicBrokersRequest) *mq_pb.LookupTopicBrokersRequest {
	return &mq_pb.LookupTopicBrokersRequest{Topic: topicFromWire(v.Topic())}
}
func LookupTopicBrokersRespToInput(r *mq_pb.LookupTopicBrokersResponse) mq_brokerwire.LookupTopicBrokersResponseInput {
	if r == nil {
		return mq_brokerwire.LookupTopicBrokersResponseInput{}
	}
	return mq_brokerwire.LookupTopicBrokersResponseInput{
		Topic:                      topicToWire(r.Topic),
		BrokerPartitionAssignments: assignmentsToWire(r.BrokerPartitionAssignments),
	}
}

// --- GetTopicConfiguration ---

func GetTopicConfigurationReqFromView(v mq_brokerwire.GetTopicConfigurationRequest) *mq_pb.GetTopicConfigurationRequest {
	return &mq_pb.GetTopicConfigurationRequest{Topic: topicFromWire(v.Topic())}
}
func GetTopicConfigurationRespToInput(r *mq_pb.GetTopicConfigurationResponse) mq_brokerwire.GetTopicConfigurationResponseInput {
	if r == nil {
		return mq_brokerwire.GetTopicConfigurationResponseInput{}
	}
	return mq_brokerwire.GetTopicConfigurationResponseInput{
		Topic:                      topicToWire(r.Topic),
		PartitionCount:             r.PartitionCount,
		BrokerPartitionAssignments: assignmentsToWire(r.BrokerPartitionAssignments),
		CreatedAtNs:                r.CreatedAtNs,
		LastUpdatedNs:              r.LastUpdatedNs,
		Retention:                  retentionToWire(r.Retention),
		MessageRecordType:          recordTypeToWire(r.MessageRecordType),
		KeyColumns:                 r.KeyColumns,
		SchemaFormat:               r.SchemaFormat,
	}
}

// --- GetTopicPublishers ---

func GetTopicPublishersReqFromView(v mq_brokerwire.GetTopicPublishersRequest) *mq_pb.GetTopicPublishersRequest {
	return &mq_pb.GetTopicPublishersRequest{Topic: topicFromWire(v.Topic())}
}
func GetTopicPublishersRespToInput(r *mq_pb.GetTopicPublishersResponse) mq_brokerwire.GetTopicPublishersResponseInput {
	if r == nil {
		return mq_brokerwire.GetTopicPublishersResponseInput{}
	}
	pubs := make([][]byte, 0, len(r.Publishers))
	for _, p := range r.Publishers {
		pubs = append(pubs, publisherToWire(p))
	}
	return mq_brokerwire.GetTopicPublishersResponseInput{Publishers: pubs}
}

// --- GetTopicSubscribers ---

func GetTopicSubscribersReqFromView(v mq_brokerwire.GetTopicSubscribersRequest) *mq_pb.GetTopicSubscribersRequest {
	return &mq_pb.GetTopicSubscribersRequest{Topic: topicFromWire(v.Topic())}
}
func GetTopicSubscribersRespToInput(r *mq_pb.GetTopicSubscribersResponse) mq_brokerwire.GetTopicSubscribersResponseInput {
	if r == nil {
		return mq_brokerwire.GetTopicSubscribersResponseInput{}
	}
	subs := make([][]byte, 0, len(r.Subscribers))
	for _, s := range r.Subscribers {
		subs = append(subs, subscriberToWire(s))
	}
	return mq_brokerwire.GetTopicSubscribersResponseInput{Subscribers: subs}
}

// --- AssignTopicPartitions ---

func AssignTopicPartitionsReqFromView(v mq_brokerwire.AssignTopicPartitionsRequest) *mq_pb.AssignTopicPartitionsRequest {
	req := &mq_pb.AssignTopicPartitionsRequest{
		Topic:      topicFromWire(v.Topic()),
		IsPrimary:   v.IsPrimary(),
		IsDraining: v.IsDraining(),
	}
	for i := 0; i < v.BrokerPartitionAssignmentsLen(); i++ {
		if a, ok := v.BrokerPartitionAssignmentAt(i); ok {
			req.BrokerPartitionAssignments = append(req.BrokerPartitionAssignments, assignmentFromView(a))
		}
	}
	return req
}
func AssignTopicPartitionsRespToInput(r *mq_pb.AssignTopicPartitionsResponse) mq_brokerwire.AssignTopicPartitionsResponseInput {
	return mq_brokerwire.AssignTopicPartitionsResponseInput{}
}

// --- ClosePublishers ---

func ClosePublishersReqFromView(v mq_brokerwire.ClosePublishersRequest) *mq_pb.ClosePublishersRequest {
	return &mq_pb.ClosePublishersRequest{Topic: topicFromWire(v.Topic()), UnixTimeNs: v.UnixTimeNs()}
}
func ClosePublishersRespToInput(r *mq_pb.ClosePublishersResponse) mq_brokerwire.ClosePublishersResponseInput {
	return mq_brokerwire.ClosePublishersResponseInput{}
}

// --- CloseSubscribers ---

func CloseSubscribersReqFromView(v mq_brokerwire.CloseSubscribersRequest) *mq_pb.CloseSubscribersRequest {
	return &mq_pb.CloseSubscribersRequest{Topic: topicFromWire(v.Topic()), UnixTimeNs: v.UnixTimeNs()}
}
func CloseSubscribersRespToInput(r *mq_pb.CloseSubscribersResponse) mq_brokerwire.CloseSubscribersResponseInput {
	return mq_brokerwire.CloseSubscribersResponseInput{}
}

// --- FetchMessage ---

func FetchMessageReqFromView(v mq_brokerwire.FetchMessageRequest) *mq_pb.FetchMessageRequest {
	return &mq_pb.FetchMessageRequest{
		Topic:         topicFromWire(v.Topic()),
		Partition:     partitionFromWire(v.Partition()),
		StartOffset:   v.StartOffset(),
		MaxBytes:      v.MaxBytes(),
		MaxMessages:   v.MaxMessages(),
		MaxWaitMs:     v.MaxWaitMs(),
		MinBytes:      v.MinBytes(),
		ConsumerGroup: v.ConsumerGroup(),
		ConsumerId:    v.ConsumerId(),
	}
}
func FetchMessageRespToInput(r *mq_pb.FetchMessageResponse) mq_brokerwire.FetchMessageResponseInput {
	if r == nil {
		return mq_brokerwire.FetchMessageResponseInput{}
	}
	msgs := make([][]byte, 0, len(r.Messages))
	for _, m := range r.Messages {
		msgs = append(msgs, dataMessageToWire(m))
	}
	return mq_brokerwire.FetchMessageResponseInput{
		Messages:       msgs,
		HighWaterMark:  r.HighWaterMark,
		LogStartOffset: r.LogStartOffset,
		EndOfPartition: r.EndOfPartition,
		Error:          r.Error,
		ErrorCode:      r.ErrorCode,
		NextOffset:     r.NextOffset,
	}
}

// --- GetPartitionRangeInfo ---

func GetPartitionRangeInfoReqFromView(v mq_brokerwire.GetPartitionRangeInfoRequest) *mq_pb.GetPartitionRangeInfoRequest {
	return &mq_pb.GetPartitionRangeInfoRequest{
		Topic: topicFromWire(v.Topic()), Partition: partitionFromWire(v.Partition()),
	}
}
func GetPartitionRangeInfoRespToInput(r *mq_pb.GetPartitionRangeInfoResponse) mq_brokerwire.GetPartitionRangeInfoResponseInput {
	if r == nil {
		return mq_brokerwire.GetPartitionRangeInfoResponseInput{}
	}
	return mq_brokerwire.GetPartitionRangeInfoResponseInput{
		OffsetRange:         offsetRangeToWire(r.OffsetRange),
		TimestampRange:      timestampRangeToWire(r.TimestampRange),
		RecordCount:         r.RecordCount,
		ActiveSubscriptions: r.ActiveSubscriptions,
		Error:               r.Error,
	}
}
