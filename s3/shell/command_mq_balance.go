package shell

import (
	"fmt"
	"io"

	mq_brokerwire "github.com/hanzoai/s3/s3/wire/mq_broker"
	"github.com/zap-proto/go/transport"
)

func init() {
	Commands = append(Commands, &commandMqBalanceTopics{})
}

type commandMqBalanceTopics struct {
}

func (c *commandMqBalanceTopics) Name() string {
	return "mq.balance"
}

func (c *commandMqBalanceTopics) Help() string {
	return `balance topic partitions

`
}

func (c *commandMqBalanceTopics) HasTag(CommandTag) bool {
	return false
}

func (c *commandMqBalanceTopics) Do(args []string, commandEnv *CommandEnv, writer io.Writer) error {

	// find the broker balancer
	brokerBalancer, err := findBrokerBalancer(commandEnv)
	if err != nil {
		return err
	}
	fmt.Fprintf(writer, "current balancer: %s\n", brokerBalancer)

	// balance topics
	conn, err := transport.Dial("tcp", brokerBalancer)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := mq_brokerwire.NewHanzoMessagingClient(conn, nil)
	_, _, err = client.BalanceTopics(mq_brokerwire.NewBalanceTopicsRequest(mq_brokerwire.BalanceTopicsRequestInput{}))
	return err
}
