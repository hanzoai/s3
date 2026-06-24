package command

import (
	"github.com/hanzoai/s3/s3/mq/agent"
	"github.com/zap-proto/go/transport"

	"github.com/hanzoai/s3/s3/glog"
	"github.com/hanzoai/s3/s3/pb"
	"github.com/hanzoai/s3/s3/security"
	"github.com/hanzoai/s3/s3/util"
)

var (
	mqAgentOptions MessageQueueAgentOptions
)

type MessageQueueAgentOptions struct {
	brokers       []pb.ServerAddress
	brokersString *string
	filerGroup    *string
	ip            *string
	port          *int
}

func init() {
	cmdMqAgent.Run = runMqAgent // break init cycle
	mqAgentOptions.brokersString = cmdMqAgent.Flag.String("broker", "localhost:17777", "comma-separated message queue brokers")
	mqAgentOptions.ip = cmdMqAgent.Flag.String("ip", "", "message queue agent host address")
	mqAgentOptions.port = cmdMqAgent.Flag.Int("port", 16777, "message queue agent ZAP server port")
}

var cmdMqAgent = &Command{
	UsageLine: "mq.agent [-port=16777] [-broker=<ip:port>]",
	Short:     "<WIP> start a message queue agent",
	Long: `start a message queue agent

	The agent runs on local server to accept ZAP calls to write or read messages.
	The messages are sent to message queue brokers.

`,
}

func runMqAgent(cmd *Command, args []string) bool {

	util.LoadSecurityConfiguration()

	mqAgentOptions.brokers = pb.ServerAddresses(*mqAgentOptions.brokersString).ToAddresses()

	return mqAgentOptions.startQueueAgent()

}

func (mqAgentOpt *MessageQueueAgentOptions) startQueueAgent() bool {

	grpcDialOption := security.LoadClientTLS(util.GetViper(), "grpc.msg_agent")

	agentServer := agent.NewMessageQueueAgent(&agent.MessageQueueAgentOptions{
		SeedBrokers: mqAgentOpt.brokers,
	}, grpcDialOption)

	// Serve the ZAP transport: unary methods via DispatchUnary, the two
	// bidirectional streaming methods via StreamHandler.
	addr := util.JoinHostPort(*mqAgentOpt.ip, *mqAgentOpt.port)
	srv, err := transport.ListenStream("tcp", addr, agentServer.DispatchUnary, agentServer.StreamHandler)
	if err != nil {
		glog.Fatalf("failed to listen on ZAP port %d: %v", *mqAgentOpt.port, err)
	}
	defer srv.Close()

	// Start localhost listener if the bind host is not already loopback/any.
	if host := *mqAgentOpt.ip; host != "localhost" && host != "" && host != "0.0.0.0" && host != "127.0.0.1" && host != "[::]" && host != "[::1]" {
		localAddr := util.JoinHostPort("localhost", *mqAgentOpt.port)
		localSrv, lErr := transport.ListenStream("tcp", localAddr, agentServer.DispatchUnary, agentServer.StreamHandler)
		if lErr != nil {
			glog.V(0).Infof("skip starting on localhost:%d: %v", *mqAgentOpt.port, lErr)
		} else {
			defer localSrv.Close()
			glog.V(0).Infof("MQ Agent listening on localhost:%d", *mqAgentOpt.port)
		}
	}

	glog.Infof("Start Hanzo S3 Message Queue Agent on %s:%d", *mqAgentOpt.ip, *mqAgentOpt.port)

	// Block forever; the server accepts on its own goroutine.
	select {}
}
