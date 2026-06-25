// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// rpc.go holds the per-RPC request builders and response decoders for the 14
// unary HanzoMessaging RPCs, built on the leaf converters in convert.go. Naming:
// <Rpc>ReqToWire(*mq_pb.<Rpc>Request) []byte and
// <Rpc>RespFromWire([]byte) (*mq_pb.<Rpc>Response, error).

package mq

import (
	mq_pb "github.com/hanzoai/s3/s3/pb/mq_pb"
	mq_brokerwire "github.com/hanzoai/s3/s3/wire/mq_broker"
)

func stringsFromView(n int, at func(int) string) []string {
	if n == 0 {
		return nil
	}
	out := make([]string, n)
	for i := 0; i < n; i++ {
		out[i] = at(i)
	}
	return out
}

func assignmentsFromView(n int, at func(int) (mq_brokerwire.BrokerPartitionAssignment, bool)) []*mq_pb.BrokerPartitionAssignment {
	if n == 0 {
		return nil
	}
	out := make([]*mq_pb.BrokerPartitionAssignment, 0, n)
	for i := 0; i < n; i++ {
		if a, ok := at(i); ok {
			out = append(out, assignmentFromView(a))
		}
	}
	return out
}

// FindBrokerLeader
func FindBrokerLeaderReqToWire(r *mq_pb.FindBrokerLeaderRequest) []byte {
	return mq_brokerwire.NewFindBrokerLeaderRequest(mq_brokerwire.FindBrokerLeaderRequestInput{FilerGroup: r.FilerGroup})
}
func FindBrokerLeaderRespFromWire(b []byte) (*mq_pb.FindBrokerLeaderResponse, error) {
	v, err := mq_brokerwire.WrapFindBrokerLeaderResponse(b)
	if err != nil {
		return nil, err
	}
	return &mq_pb.FindBrokerLeaderResponse{Broker: v.Broker()}, nil
}

// BalanceTopics
func BalanceTopicsReqToWire(r *mq_pb.BalanceTopicsRequest) []byte {
	return mq_brokerwire.NewBalanceTopicsRequest(mq_brokerwire.BalanceTopicsRequestInput{})
}
func BalanceTopicsRespFromWire(b []byte) (*mq_pb.BalanceTopicsResponse, error) {
	if _, err := mq_brokerwire.WrapBalanceTopicsResponse(b); err != nil {
		return nil, err
	}
	return &mq_pb.BalanceTopicsResponse{}, nil
}

// ListTopics
func ListTopicsReqToWire(r *mq_pb.ListTopicsRequest) []byte {
	return mq_brokerwire.NewListTopicsRequest(mq_brokerwire.ListTopicsRequestInput{})
}
func ListTopicsRespFromWire(b []byte) (*mq_pb.ListTopicsResponse, error) {
	v, err := mq_brokerwire.WrapListTopicsResponse(b)
	if err != nil {
		return nil, err
	}
	return &mq_pb.ListTopicsResponse{Topics: topicsFromView(v.TopicsLen(), v.TopicAt)}, nil
}

// TopicExists
func TopicExistsReqToWire(r *mq_pb.TopicExistsRequest) []byte {
	return mq_brokerwire.NewTopicExistsRequest(mq_brokerwire.TopicExistsRequestInput{Topic: topicToWire(r.Topic)})
}
func TopicExistsRespFromWire(b []byte) (*mq_pb.TopicExistsResponse, error) {
	v, err := mq_brokerwire.WrapTopicExistsResponse(b)
	if err != nil {
		return nil, err
	}
	return &mq_pb.TopicExistsResponse{Exists: v.Exists()}, nil
}

// ConfigureTopic
func ConfigureTopicReqToWire(r *mq_pb.ConfigureTopicRequest) []byte {
	return mq_brokerwire.NewConfigureTopicRequest(mq_brokerwire.ConfigureTopicRequestInput{
		Topic:             topicToWire(r.Topic),
		PartitionCount:    r.PartitionCount,
		Retention:         retentionToWire(r.Retention),
		MessageRecordType: recordTypeToWire(r.MessageRecordType),
		KeyColumns:        r.KeyColumns,
		SchemaFormat:      r.SchemaFormat,
	})
}
func ConfigureTopicRespFromWire(b []byte) (*mq_pb.ConfigureTopicResponse, error) {
	v, err := mq_brokerwire.WrapConfigureTopicResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &mq_pb.ConfigureTopicResponse{
		BrokerPartitionAssignments: assignmentsFromView(v.BrokerPartitionAssignmentsLen(), v.BrokerPartitionAssignmentAt),
		MessageRecordType:          recordTypeFromWire(v.MessageRecordType()),
		KeyColumns:                 stringsFromView(v.KeyColumnsLen(), v.KeyColumnAt),
		SchemaFormat:               v.SchemaFormat(),
	}
	if rt, ok := v.Retention(); ok {
		resp.Retention = retentionFromView(rt)
	}
	return resp, nil
}

// LookupTopicBrokers
func LookupTopicBrokersReqToWire(r *mq_pb.LookupTopicBrokersRequest) []byte {
	return mq_brokerwire.NewLookupTopicBrokersRequest(mq_brokerwire.LookupTopicBrokersRequestInput{Topic: topicToWire(r.Topic)})
}
func LookupTopicBrokersRespFromWire(b []byte) (*mq_pb.LookupTopicBrokersResponse, error) {
	v, err := mq_brokerwire.WrapLookupTopicBrokersResponse(b)
	if err != nil {
		return nil, err
	}
	return &mq_pb.LookupTopicBrokersResponse{
		Topic:                      topicFromWire(v.Topic()),
		BrokerPartitionAssignments: assignmentsFromView(v.BrokerPartitionAssignmentsLen(), v.BrokerPartitionAssignmentAt),
	}, nil
}

// GetTopicConfiguration
func GetTopicConfigurationReqToWire(r *mq_pb.GetTopicConfigurationRequest) []byte {
	return mq_brokerwire.NewGetTopicConfigurationRequest(mq_brokerwire.GetTopicConfigurationRequestInput{Topic: topicToWire(r.Topic)})
}
func GetTopicConfigurationRespFromWire(b []byte) (*mq_pb.GetTopicConfigurationResponse, error) {
	v, err := mq_brokerwire.WrapGetTopicConfigurationResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &mq_pb.GetTopicConfigurationResponse{
		Topic:                      topicFromWire(v.Topic()),
		PartitionCount:             v.PartitionCount(),
		BrokerPartitionAssignments: assignmentsFromView(v.BrokerPartitionAssignmentsLen(), v.BrokerPartitionAssignmentAt),
		CreatedAtNs:                v.CreatedAtNs(),
		LastUpdatedNs:              v.LastUpdatedNs(),
		MessageRecordType:          recordTypeFromWire(v.MessageRecordType()),
		KeyColumns:                 stringsFromView(v.KeyColumnsLen(), v.KeyColumnAt),
		SchemaFormat:               v.SchemaFormat(),
	}
	if rt, ok := v.Retention(); ok {
		resp.Retention = retentionFromView(rt)
	}
	return resp, nil
}

// GetTopicPublishers
func GetTopicPublishersReqToWire(r *mq_pb.GetTopicPublishersRequest) []byte {
	return mq_brokerwire.NewGetTopicPublishersRequest(mq_brokerwire.GetTopicPublishersRequestInput{Topic: topicToWire(r.Topic)})
}
func GetTopicPublishersRespFromWire(b []byte) (*mq_pb.GetTopicPublishersResponse, error) {
	v, err := mq_brokerwire.WrapGetTopicPublishersResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &mq_pb.GetTopicPublishersResponse{}
	for i := 0; i < v.PublishersLen(); i++ {
		if p, ok := v.PublisherAt(i); ok {
			resp.Publishers = append(resp.Publishers, publisherFromView(p))
		}
	}
	return resp, nil
}

// GetTopicSubscribers
func GetTopicSubscribersReqToWire(r *mq_pb.GetTopicSubscribersRequest) []byte {
	return mq_brokerwire.NewGetTopicSubscribersRequest(mq_brokerwire.GetTopicSubscribersRequestInput{Topic: topicToWire(r.Topic)})
}
func GetTopicSubscribersRespFromWire(b []byte) (*mq_pb.GetTopicSubscribersResponse, error) {
	v, err := mq_brokerwire.WrapGetTopicSubscribersResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &mq_pb.GetTopicSubscribersResponse{}
	for i := 0; i < v.SubscribersLen(); i++ {
		if s, ok := v.SubscriberAt(i); ok {
			resp.Subscribers = append(resp.Subscribers, subscriberFromView(s))
		}
	}
	return resp, nil
}

// AssignTopicPartitions
func AssignTopicPartitionsReqToWire(r *mq_pb.AssignTopicPartitionsRequest) []byte {
	return mq_brokerwire.NewAssignTopicPartitionsRequest(mq_brokerwire.AssignTopicPartitionsRequestInput{
		Topic:                      topicToWire(r.Topic),
		BrokerPartitionAssignments: assignmentsToWire(r.BrokerPartitionAssignments),
		IsLeader:                   r.IsLeader,
		IsDraining:                 r.IsDraining,
	})
}
func AssignTopicPartitionsRespFromWire(b []byte) (*mq_pb.AssignTopicPartitionsResponse, error) {
	if _, err := mq_brokerwire.WrapAssignTopicPartitionsResponse(b); err != nil {
		return nil, err
	}
	return &mq_pb.AssignTopicPartitionsResponse{}, nil
}

// ClosePublishers
func ClosePublishersReqToWire(r *mq_pb.ClosePublishersRequest) []byte {
	return mq_brokerwire.NewClosePublishersRequest(mq_brokerwire.ClosePublishersRequestInput{
		Topic: topicToWire(r.Topic), UnixTimeNs: r.UnixTimeNs,
	})
}
func ClosePublishersRespFromWire(b []byte) (*mq_pb.ClosePublishersResponse, error) {
	if _, err := mq_brokerwire.WrapClosePublishersResponse(b); err != nil {
		return nil, err
	}
	return &mq_pb.ClosePublishersResponse{}, nil
}

// CloseSubscribers
func CloseSubscribersReqToWire(r *mq_pb.CloseSubscribersRequest) []byte {
	return mq_brokerwire.NewCloseSubscribersRequest(mq_brokerwire.CloseSubscribersRequestInput{
		Topic: topicToWire(r.Topic), UnixTimeNs: r.UnixTimeNs,
	})
}
func CloseSubscribersRespFromWire(b []byte) (*mq_pb.CloseSubscribersResponse, error) {
	if _, err := mq_brokerwire.WrapCloseSubscribersResponse(b); err != nil {
		return nil, err
	}
	return &mq_pb.CloseSubscribersResponse{}, nil
}

// FetchMessage
func FetchMessageReqToWire(r *mq_pb.FetchMessageRequest) []byte {
	return mq_brokerwire.NewFetchMessageRequest(mq_brokerwire.FetchMessageRequestInput{
		Topic:         topicToWire(r.Topic),
		Partition:     partitionToWire(r.Partition),
		StartOffset:   r.StartOffset,
		MaxBytes:      r.MaxBytes,
		MaxMessages:   r.MaxMessages,
		MaxWaitMs:     r.MaxWaitMs,
		MinBytes:      r.MinBytes,
		ConsumerGroup: r.ConsumerGroup,
		ConsumerId:    r.ConsumerId,
	})
}
func FetchMessageRespFromWire(b []byte) (*mq_pb.FetchMessageResponse, error) {
	v, err := mq_brokerwire.WrapFetchMessageResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &mq_pb.FetchMessageResponse{
		HighWaterMark:  v.HighWaterMark(),
		LogStartOffset: v.LogStartOffset(),
		EndOfPartition: v.EndOfPartition(),
		Error:          v.Error(),
		ErrorCode:      v.ErrorCode(),
		NextOffset:     v.NextOffset(),
	}
	for i := 0; i < v.MessagesLen(); i++ {
		if m, ok := v.MessageAt(i); ok {
			resp.Messages = append(resp.Messages, dataMessageFromView(m))
		}
	}
	return resp, nil
}

// GetPartitionRangeInfo
func GetPartitionRangeInfoReqToWire(r *mq_pb.GetPartitionRangeInfoRequest) []byte {
	return mq_brokerwire.NewGetPartitionRangeInfoRequest(mq_brokerwire.GetPartitionRangeInfoRequestInput{
		Topic: topicToWire(r.Topic), Partition: partitionToWire(r.Partition),
	})
}
func GetPartitionRangeInfoRespFromWire(b []byte) (*mq_pb.GetPartitionRangeInfoResponse, error) {
	v, err := mq_brokerwire.WrapGetPartitionRangeInfoResponse(b)
	if err != nil {
		return nil, err
	}
	resp := &mq_pb.GetPartitionRangeInfoResponse{
		RecordCount:         v.RecordCount(),
		ActiveSubscriptions: v.ActiveSubscriptions(),
		Error:               v.Error(),
	}
	if o, ok := v.OffsetRange(); ok {
		resp.OffsetRange = offsetRangeFromView(o)
	}
	if ts, ok := v.TimestampRange(); ok {
		resp.TimestampRange = timestampRangeFromView(ts)
	}
	return resp, nil
}
