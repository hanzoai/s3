package shell

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"time"

	mq_brokerwire "github.com/hanzoai/s3/s3/wire/mq_broker"
	mq_schemawire "github.com/hanzoai/s3/s3/wire/mq_schema"
	"github.com/zap-proto/go/transport"
)

func init() {
	Commands = append(Commands, &commandMqTopicConfigure{})
}

type commandMqTopicConfigure struct {
}

func (c *commandMqTopicConfigure) Name() string {
	return "mq.topic.configure"
}

func (c *commandMqTopicConfigure) Help() string {
	return `configure a topic with a given name

	Example:
		mq.topic.configure -namespace <namespace> -topic <topic_name> -partitionCount <partition_count>

	Retention (delete messages older than the configured duration):
		mq.topic.configure -namespace <namespace> -topic <topic_name> \
			-retention 168h -retentionEnabled

		# disable retention on an existing topic
		mq.topic.configure -namespace <namespace> -topic <topic_name> \
			-retentionEnabled=false

	-retention accepts any Go duration string ("24h", "168h", "30m"). Use
	-retentionSeconds for raw seconds when scripting. Specifying both is an
	error.

	When you set only some retention flags (for example, -retentionEnabled
	without -retention), the unspecified field is read from the current
	server-side configuration so it isn't accidentally zeroed. Omitting all
	retention flags leaves the existing retention configuration alone.
`
}

func (c *commandMqTopicConfigure) HasTag(CommandTag) bool {
	return false
}

func (c *commandMqTopicConfigure) Do(args []string, commandEnv *CommandEnv, writer io.Writer) error {

	// parse parameters
	mqCommand := flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	namespace := mqCommand.String("namespace", "", "namespace name")
	topicName := mqCommand.String("topic", "", "topic name")
	partitionCount := mqCommand.Int("partitionCount", 6, "partition count")
	retention := mqCommand.Duration("retention", 0, "retention duration (Go duration string, e.g. 168h). Mutually exclusive with -retentionSeconds.")
	retentionSeconds := mqCommand.Int64("retentionSeconds", 0, "retention duration in seconds. Mutually exclusive with -retention.")
	retentionEnabled := mqCommand.Bool("retentionEnabled", false, "enable retention enforcement on the topic")
	if err := mqCommand.Parse(args); err != nil {
		return err
	}

	// Detect which retention flags the user actually provided. Using Visit
	// (rather than value comparison) means an explicit `-retention=0` is
	// treated as "user provided" and still triggers the mutual-exclusion
	// check or partial-merge with current state.
	var userSetRetention, userSetSeconds, userSetEnabled bool
	mqCommand.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "retention":
			userSetRetention = true
		case "retentionSeconds":
			userSetSeconds = true
		case "retentionEnabled":
			userSetEnabled = true
		}
	})

	if userSetRetention && userSetSeconds {
		return fmt.Errorf("-retention and -retentionSeconds are mutually exclusive")
	}
	if *retention < 0 || *retentionSeconds < 0 {
		return fmt.Errorf("retention duration must be >= 0")
	}

	retentionTouched := userSetRetention || userSetSeconds || userSetEnabled

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

	// Build the retention message. When the user touches any retention flag we
	// must send a fully-populated TopicRetention so partial flags don't zero
	// the other field server-side: fetch the current configuration and use
	// its values for whatever the user didn't specify.
	var retentionBuf []byte
	if retentionTouched {
		var seconds int64
		var enabled bool
		topic := mq_schemawire.NewTopic(mq_schemawire.TopicInput{Namespace: *namespace, Name: *topicName})
		if _, body, getErr := client.GetTopicConfiguration(mq_brokerwire.NewGetTopicConfigurationRequest(mq_brokerwire.GetTopicConfigurationRequestInput{
			Topic: topic,
		})); getErr == nil {
			// Topic may not exist yet — that's fine, we'll create it with
			// the user-supplied retention only.
			if cur, wrapErr := mq_brokerwire.WrapGetTopicConfigurationResponse(body); wrapErr == nil {
				if r, ok := cur.Retention(); ok {
					seconds = r.RetentionSeconds()
					enabled = r.Enabled()
				}
			}
		}

		if userSetRetention {
			seconds = int64((*retention) / time.Second)
		} else if userSetSeconds {
			seconds = *retentionSeconds
		}
		if userSetEnabled {
			enabled = *retentionEnabled
		}
		retentionBuf = mq_brokerwire.NewTopicRetention(mq_brokerwire.TopicRetentionInput{
			RetentionSeconds: seconds,
			Enabled:          enabled,
		})
	}

	// create / update topic
	topic := mq_schemawire.NewTopic(mq_schemawire.TopicInput{Namespace: *namespace, Name: *topicName})
	_, body, err := client.ConfigureTopic(mq_brokerwire.NewConfigureTopicRequest(mq_brokerwire.ConfigureTopicRequestInput{
		Topic:          topic,
		PartitionCount: int32(*partitionCount),
		Retention:      retentionBuf,
	}))
	if err != nil {
		return err
	}
	resp, err := mq_brokerwire.WrapConfigureTopicResponse(body)
	if err != nil {
		return err
	}

	type assignmentView struct {
		LeaderBroker   string
		FollowerBroker string
	}
	type responseView struct {
		PartitionAssignments []assignmentView
		RetentionSeconds     int64
		RetentionEnabled     bool
		SchemaFormat         string
	}
	view := responseView{SchemaFormat: resp.SchemaFormat()}
	for i := 0; i < resp.BrokerPartitionAssignmentsLen(); i++ {
		if a, ok := resp.BrokerPartitionAssignmentAt(i); ok {
			view.PartitionAssignments = append(view.PartitionAssignments, assignmentView{
				LeaderBroker:   a.LeaderBroker(),
				FollowerBroker: a.FollowerBroker(),
			})
		}
	}
	if r, ok := resp.Retention(); ok {
		view.RetentionSeconds = r.RetentionSeconds()
		view.RetentionEnabled = r.Enabled()
	}
	output, _ := json.MarshalIndent(view, "", "  ")
	fmt.Fprintf(writer, "response:\n%+v\n", string(output))
	return nil
}
