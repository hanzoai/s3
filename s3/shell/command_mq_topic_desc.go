package shell

import (
	"flag"
	"fmt"
	"io"

	mq_brokerwire "github.com/hanzoai/s3/s3/wire/mq_broker"
	mq_schemawire "github.com/hanzoai/s3/s3/wire/mq_schema"
	"github.com/zap-proto/go/transport"
)

func init() {
	Commands = append(Commands, &commandMqTopicDescribe{})
}

type commandMqTopicDescribe struct {
}

func (c *commandMqTopicDescribe) Name() string {
	return "mq.topic.describe"
}

func (c *commandMqTopicDescribe) Help() string {
	return `describe a topic`
}

func (c *commandMqTopicDescribe) HasTag(CommandTag) bool {
	return false
}

func (c *commandMqTopicDescribe) Do(args []string, commandEnv *CommandEnv, writer io.Writer) error {
	// parse parameters
	mqCommand := flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	namespace := mqCommand.String("namespace", "", "namespace name")
	topicName := mqCommand.String("topic", "", "topic name")
	if err := mqCommand.Parse(args); err != nil {
		return err
	}

	// find the broker balancer
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
	topic := mq_schemawire.NewTopic(mq_schemawire.TopicInput{
		Namespace: *namespace,
		Name:      *topicName,
	})
	_, body, err := client.LookupTopicBrokers(mq_brokerwire.NewLookupTopicBrokersRequest(mq_brokerwire.LookupTopicBrokersRequestInput{
		Topic: topic,
	}))
	if err != nil {
		return err
	}
	resp, err := mq_brokerwire.WrapLookupTopicBrokersResponse(body)
	if err != nil {
		return err
	}
	for i := 0; i < resp.BrokerPartitionAssignmentsLen(); i++ {
		assignment, ok := resp.BrokerPartitionAssignmentAt(i)
		if !ok {
			continue
		}
		part, err := mq_schemawire.WrapPartition(assignment.Partition())
		if err != nil {
			return err
		}
		fmt.Fprintf(writer, "  partition(ring:%d range:[%d,%d) ts:%d) leader:%s follower:%s\n",
			part.RingSize(), part.RangeStart(), part.RangeStop(), part.UnixTimeNs(),
			assignment.LeaderBroker(), assignment.FollowerBroker())
	}
	return nil
}
