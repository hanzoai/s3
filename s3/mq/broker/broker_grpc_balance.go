package broker

import (
	"context"

	"github.com/hanzoai/s3/s3/pb"
	"github.com/hanzoai/s3/s3/pb/mq_pb"
)

func (b *MessageQueueBroker) BalanceTopics(ctx context.Context, request *mq_pb.BalanceTopicsRequest) (resp *mq_pb.BalanceTopicsResponse, err error) {
	if !b.isLockOwner() {
		proxyErr := b.withBrokerClient(false, pb.ServerAddress(b.lockAsBalancer.LockOwner()), func(client mq_pb.HanzoMessagingClient) error {
			resp, err = client.BalanceTopics(ctx, request)
			return nil
		})
		if proxyErr != nil {
			return nil, proxyErr
		}
		return resp, err
	}

	ret := &mq_pb.BalanceTopicsResponse{}

	actions := b.PubBalancer.BalancePublishers()
	err = b.PubBalancer.ExecuteBalanceAction(actions, b.grpcDialOption)

	return ret, err
}
