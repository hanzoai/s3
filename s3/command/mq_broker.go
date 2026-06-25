package command

import (
	"crypto/tls"

	"github.com/zap-proto/go/transport"

	"github.com/hanzoai/s3/s3/glog"
	"github.com/hanzoai/s3/s3/mq/broker"
	"github.com/hanzoai/s3/s3/mqzap"
	"github.com/hanzoai/s3/s3/pb"
	"github.com/hanzoai/s3/s3/security"
	"github.com/hanzoai/s3/s3/util"
	"github.com/hanzoai/s3/s3/util/grace"
	mq_brokerwire "github.com/hanzoai/s3/s3/wire/mq_broker"
)

var (
	mqBrokerStandaloneOptions MessageQueueBrokerOptions
)

type MessageQueueBrokerOptions struct {
	masters          map[string]pb.ServerAddress
	mastersString    *string
	filerGroup       *string
	ip               *string
	port             *int
	dataCenter       *string
	rack             *string
	cpuprofile       *string
	memprofile       *string
	logFlushInterval *int
	debug            *bool
	debugPort        *int
}

func init() {
	cmdMqBroker.Run = runMqBroker // break init cycle
	mqBrokerStandaloneOptions.mastersString = cmdMqBroker.Flag.String("master", "localhost:9333", "comma-separated master servers")
	mqBrokerStandaloneOptions.filerGroup = cmdMqBroker.Flag.String("filerGroup", "", "share metadata with other filers in the same filerGroup")
	mqBrokerStandaloneOptions.ip = cmdMqBroker.Flag.String("ip", util.DetectedHostAddress(), "broker host address")
	mqBrokerStandaloneOptions.port = cmdMqBroker.Flag.Int("port", 17777, "broker ZAP listen port")
	mqBrokerStandaloneOptions.dataCenter = cmdMqBroker.Flag.String("dataCenter", "", "prefer to read and write to volumes in this data center")
	mqBrokerStandaloneOptions.rack = cmdMqBroker.Flag.String("rack", "", "prefer to write to volumes in this rack")
	mqBrokerStandaloneOptions.cpuprofile = cmdMqBroker.Flag.String("cpuprofile", "", "cpu profile output file")
	mqBrokerStandaloneOptions.memprofile = cmdMqBroker.Flag.String("memprofile", "", "memory profile output file")
	mqBrokerStandaloneOptions.logFlushInterval = cmdMqBroker.Flag.Int("logFlushInterval", 5, "log buffer flush interval in seconds")
	mqBrokerStandaloneOptions.debug = cmdMqBroker.Flag.Bool("debug", false, "serves runtime profiling data via pprof on the port specified by -debug.port")
	mqBrokerStandaloneOptions.debugPort = cmdMqBroker.Flag.Int("debug.port", 6060, "http port for debugging")
}

var cmdMqBroker = &Command{
	UsageLine: "mq.broker [-port=17777] [-master=<ip:port>]",
	Short:     "<WIP> start a message queue broker",
	Long: `start a message queue broker

	The broker accepts native ZAP calls to write or read messages. The messages are stored via filer.
	The brokers are stateless. To scale up, just add more brokers.

`,
}

func runMqBroker(cmd *Command, args []string) bool {
	if *mqBrokerStandaloneOptions.debug {
		grace.StartDebugServer(*mqBrokerStandaloneOptions.debugPort)
	}

	util.LoadSecurityConfiguration()

	mqBrokerStandaloneOptions.masters = pb.ServerAddresses(*mqBrokerStandaloneOptions.mastersString).ToAddressMap()

	return mqBrokerStandaloneOptions.startQueueServer()

}

func (mqBrokerOpt *MessageQueueBrokerOptions) startQueueServer() bool {

	*mqBrokerStandaloneOptions.cpuprofile = util.ResolvePath(*mqBrokerStandaloneOptions.cpuprofile)
	*mqBrokerStandaloneOptions.memprofile = util.ResolvePath(*mqBrokerStandaloneOptions.memprofile)
	grace.SetupProfiling(*mqBrokerStandaloneOptions.cpuprofile, *mqBrokerStandaloneOptions.memprofile)

	grpcDialOption := security.LoadClientTLS(util.GetViper(), "grpc.msg_broker")

	qs, err := broker.NewMessageBroker(&broker.MessageQueueBrokerOption{
		Masters:            mqBrokerOpt.masters,
		FilerGroup:         *mqBrokerOpt.filerGroup,
		DataCenter:         *mqBrokerOpt.dataCenter,
		Rack:               *mqBrokerOpt.rack,
		DefaultReplication: "",
		MaxMB:              0,
		Ip:                 *mqBrokerOpt.ip,
		Port:               *mqBrokerOpt.port,
		LogFlushInterval:   *mqBrokerOpt.logFlushInterval,
	}, grpcDialOption)
	if err != nil {
		glog.Fatalf("failed to create new message broker for queue server: %v", err)
	}

	// Serve the HanzoMessaging service over the native ZAP transport on the broker
	// port. Clients reach it via the broker address over transport.Dial (unary) and
	// transport.OpenStream (streaming) — see pb.WithBrokerGrpcClient and the mqzap
	// backend. This replaces the legacy gRPC HanzoMessaging server: the whole broker
	// (14 unary + 7 streaming RPCs) now answers over ZAP, no gRPC.
	//
	// When grpc.msg_broker.cert/.key is configured the listener is PQ-secured mTLS:
	// transport.PQTLSConfig pins the X25519MLKEM768 hybrid (PQ X-Wing) and the same
	// cert/CA/allowed-CN gate the legacy gRPC broker enforced applies — no security
	// downgrade. Otherwise it is plaintext (loopback / dev), exactly as the gRPC
	// broker was plaintext when no cert was configured.
	dispatch := mq_brokerwire.Dispatch(mqzap.NewServerBackend(qs))
	streamHandler := mqzap.NewStreamServer(qs).ServeStream
	tlsCfg := security.ServerTLSConfig(util.GetViper(), "grpc.msg_broker")

	// zapServers collects every ZAP listener so shutdown can close them all.
	var zapServers []*transport.Server

	// localhost listener: lets co-located clients reach the broker without the
	// routable IP. Started first so the main listener's blocking Serve below owns
	// the goroutine's lifetime.
	localAddr := util.JoinHostPort("localhost", *mqBrokerOpt.port)
	if localSrv, lerr := serveBrokerZap(localAddr, tlsCfg, dispatch, streamHandler); lerr != nil {
		glog.V(0).Infof("MQ Broker localhost listener unavailable on %s: %v", localAddr, lerr)
	} else {
		zapServers = append(zapServers, localSrv)
		glog.V(0).Infof("MQ Broker listening on %s", localAddr)
	}

	brokerAddr := util.JoinHostPort(*mqBrokerOpt.ip, *mqBrokerOpt.port)
	mainSrv, merr := serveBrokerZap(brokerAddr, tlsCfg, dispatch, streamHandler)
	if merr != nil {
		glog.Fatalf("failed to serve MQ Broker over ZAP on %s: %v", brokerAddr, merr)
	}
	zapServers = append(zapServers, mainSrv)
	if tlsCfg != nil {
		glog.V(0).Infof("Serving MQ Broker over PQ-TLS (X25519MLKEM768) ZAP transport on %s", brokerAddr)
	} else {
		glog.V(0).Infof("Serving MQ Broker over native ZAP transport (plaintext; set grpc.msg_broker.cert/.key for PQ-TLS) on %s", brokerAddr)
	}

	// transport.ListenStream backgrounds its own accept loop, so block here until
	// shutdown — closing every ZAP listener on interrupt releases the ports for
	// restart (important for the start/stop/start cycle of integration tests).
	done := make(chan struct{})
	grace.OnInterrupt(func() {
		for _, s := range zapServers {
			_ = s.Close()
		}
		close(done)
	})
	<-done

	return true

}

// serveBrokerZap starts one ZAP HanzoMessaging listener at addr, PQ-secured when
// tlsCfg is non-nil and plaintext otherwise — the single place the broker's
// listen-side security gate lives. The returned *Server runs its accept loop on
// its own goroutine; Close stops it.
func serveBrokerZap(addr string, tlsCfg *tls.Config, dispatch transport.Dispatch, stream transport.StreamHandler) (*transport.Server, error) {
	if tlsCfg != nil {
		return transport.ListenStreamTLS("tcp", addr, transport.PQTLSConfig(tlsCfg), dispatch, stream)
	}
	return transport.ListenStream("tcp", addr, dispatch, stream)
}
