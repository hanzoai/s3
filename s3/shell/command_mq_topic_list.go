package shell

import (
	"context"
	"fmt"
	"io"

	"github.com/hanzoai/s3/s3/mq/pub_balancer"
	"github.com/hanzoai/s3/s3/pb/filer_pb"
	mq_brokerwire "github.com/hanzoai/s3/s3/wire/mq_broker"
	mq_schemawire "github.com/hanzoai/s3/s3/wire/mq_schema"
	"github.com/zap-proto/go/transport"
)

func init() {
	Commands = append(Commands, &commandMqTopicList{})
}

type commandMqTopicList struct {
}

func (c *commandMqTopicList) Name() string {
	return "mq.topic.list"
}

func (c *commandMqTopicList) Help() string {
	return `print out all topics`
}

func (c *commandMqTopicList) HasTag(CommandTag) bool {
	return false
}

func (c *commandMqTopicList) Do(args []string, commandEnv *CommandEnv, writer io.Writer) error {

	brokerBalancer, err := findBrokerBalancer(commandEnv)
	if err != nil {
		return err
	}
	fmt.Fprintf(writer, "current balancer: %s\n", brokerBalancer)

	conn, err := transport.Dial("tcp", brokerBalancer)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := mq_brokerwire.NewHanzoMessagingClient(conn, nil)
	_, body, err := client.ListTopics(mq_brokerwire.NewListTopicsRequest(mq_brokerwire.ListTopicsRequestInput{}))
	if err != nil {
		return err
	}
	resp, err := mq_brokerwire.WrapListTopicsResponse(body)
	if err != nil {
		return err
	}
	if resp.TopicsLen() == 0 {
		fmt.Fprintf(writer, "no topics found\n")
		return nil
	}
	for i := 0; i < resp.TopicsLen(); i++ {
		topic, err := mq_schemawire.WrapTopic(resp.TopicAt(i))
		if err != nil {
			return err
		}
		fmt.Fprintf(writer, "  %s.%s\n", topic.Namespace(), topic.Name())
	}
	return nil
}

func findBrokerBalancer(commandEnv *CommandEnv) (brokerBalancer string, err error) {
	err = commandEnv.WithFilerClient(false, func(client filer_pb.HanzoFilerClient) error {
		resp, err := client.FindLockOwner(context.Background(), &filer_pb.FindLockOwnerRequest{
			Name: pub_balancer.LockBrokerBalancer,
		})
		if err != nil {
			return fmt.Errorf("FindLockOwner: %w", err)
		}
		brokerBalancer = resp.Owner
		return nil
	})
	return
}
