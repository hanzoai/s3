package command

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hanzoai/s3/s3/util/version"

	"slices"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"google.golang.org/grpc/reflection"

	stats_collect "github.com/hanzoai/s3/s3/stats"

	"github.com/hanzoai/s3/s3/util/grace"

	"github.com/hanzoai/s3/s3/glog"
	"github.com/hanzoai/s3/s3/svc/master"
	"github.com/hanzoai/s3/s3/pb"
	"github.com/hanzoai/s3/s3/security"
	s3server "github.com/hanzoai/s3/s3/server"
	"github.com/hanzoai/s3/s3/storage/backend"
	"github.com/hanzoai/s3/s3/util"
	masterwire "github.com/hanzoai/s3/s3/wire/master"
	masterstream "github.com/hanzoai/s3/s3/wire/master/masterstream"
	"github.com/zap-proto/go/transport"
)

var (
	m MasterOptions
)

type MasterOptions struct {
	port                       *int
	portGrpc                   *int
	ip                         *string
	ipBind                     *string
	metaFolder                 *string
	peers                      *string
	mastersDeprecated          *string // deprecated, for backward compatibility in master.follower
	volumeSizeLimitMB          *uint
	volumePreallocate          *bool
	maxParallelVacuumPerServer *int
	// pulseSeconds       *int
	defaultReplication *string
	garbageThreshold   *float64
	whiteList          *string
	disableHttp        *bool
	metricsAddress     *string
	metricsIntervalSec *int
	raftResumeState    *bool
	metricsHttpPort    *int
	metricsHttpIp      *string
	heartbeatInterval  *time.Duration
	electionTimeout    *time.Duration
	raftHashicorp      *bool
	raftBootstrap      *bool
	telemetryUrl       *string
	telemetryEnabled   *bool
	debug              *bool
	debugPort          *int
	// shutdownCtx, when non-nil, tells startMaster to shut down once the ctx
	// is cancelled. Used by integration tests and by s3 mini; nil for
	// standalone s3 master.
	shutdownCtx context.Context
}

func init() {
	cmdMaster.Run = runMaster // break init cycle
	m.port = cmdMaster.Flag.Int("port", 9333, "http listen port")
	m.portGrpc = cmdMaster.Flag.Int("port.grpc", 0, "grpc listen port")
	m.ip = cmdMaster.Flag.String("ip", util.DetectedHostAddress(), "master <ip>|<server> address, also used as identifier")
	m.ipBind = cmdMaster.Flag.String("ip.bind", "", "ip address to bind to. If empty, default to same as -ip option.")
	m.metaFolder = cmdMaster.Flag.String("mdir", os.TempDir(), "data directory to store meta data")
	m.peers = cmdMaster.Flag.String("peers", "", "all master nodes in comma separated ip:port list, example: 127.0.0.1:9093,127.0.0.1:9094,127.0.0.1:9095; use 'none' for single-master mode")
	m.volumeSizeLimitMB = cmdMaster.Flag.Uint("volumeSizeLimitMB", 30*1000, "Master stops directing writes to oversized volumes.")
	m.volumePreallocate = cmdMaster.Flag.Bool("volumePreallocate", false, "Preallocate disk space for volumes.")
	m.maxParallelVacuumPerServer = cmdMaster.Flag.Int("maxParallelVacuumPerServer", 1, "maximum number of volumes to vacuum in parallel per volume server")
	// m.pulseSeconds = cmdMaster.Flag.Int("pulseSeconds", 5, "number of seconds between heartbeats")
	m.defaultReplication = cmdMaster.Flag.String("defaultReplication", "", "Default replication type if not specified.")
	m.garbageThreshold = cmdMaster.Flag.Float64("garbageThreshold", 0.3, "threshold to vacuum and reclaim spaces")
	m.whiteList = cmdMaster.Flag.String("whiteList", "", "comma separated Ip addresses having write permission. No limit if empty.")
	m.disableHttp = cmdMaster.Flag.Bool("disableHttp", false, "disable http requests, only gRPC operations are allowed.")
	m.metricsAddress = cmdMaster.Flag.String("metrics.address", "", "Prometheus gateway address <host>:<port>")
	m.metricsIntervalSec = cmdMaster.Flag.Int("metrics.intervalSeconds", 15, "Prometheus push interval in seconds")
	m.metricsHttpPort = cmdMaster.Flag.Int("metricsPort", 0, "Prometheus metrics listen port")
	m.metricsHttpIp = cmdMaster.Flag.String("metricsIp", "", "metrics listen ip. If empty, default to same as -ip.bind option.")
	m.raftResumeState = cmdMaster.Flag.Bool("resumeState", true, "resume previous state on start master server")
	m.heartbeatInterval = cmdMaster.Flag.Duration("heartbeatInterval", 300*time.Millisecond, "heartbeat interval of master servers, and will be randomly multiplied by [1, 1.25)")
	m.electionTimeout = cmdMaster.Flag.Duration("electionTimeout", 10*time.Second, "election timeout of master servers")
	m.raftHashicorp = cmdMaster.Flag.Bool("raftHashicorp", false, "use hashicorp raft")
	m.raftBootstrap = cmdMaster.Flag.Bool("raftBootstrap", false, "Whether to bootstrap the Raft cluster")
	m.telemetryUrl = cmdMaster.Flag.String("telemetry.url", "", "telemetry server URL to send usage statistics")
	m.telemetryEnabled = cmdMaster.Flag.Bool("telemetry", false, "enable telemetry reporting")
	m.debug = cmdMaster.Flag.Bool("debug", false, "serves runtime profiling data via pprof on the port specified by -debug.port")
	m.debugPort = cmdMaster.Flag.Int("debug.port", 6060, "http port for debugging")
}

var cmdMaster = &Command{
	UsageLine: "master -port=9333",
	Short:     "start a master server",
	Long: `start a master server to provide volume=>location mapping service and sequence number of file ids

	The configuration file "security.toml" is read from ".", "$HOME/.s3/", "/usr/local/etc/hanzo/", or "/etc/hanzo/", in that order.

	The example security.toml configuration file can be generated by "s3 scaffold -config=security"

	For single-master setups, use -peers=none to skip Raft quorum wait and enable instant startup.
	This is ideal for development or standalone deployments.

  `,
}

var (
	masterCpuProfile = cmdMaster.Flag.String("cpuprofile", "", "cpu profile output file")
	masterMemProfile = cmdMaster.Flag.String("memprofile", "", "memory profile output file")
)

func runMaster(cmd *Command, args []string) bool {
	if *m.debug {
		grace.StartDebugServer(*m.debugPort)
	}

	util.LoadSecurityConfiguration()
	util.LoadConfiguration("master", false)

	// bind viper configuration to command line flags
	if v := util.GetViper().GetString("master.mdir"); v != "" {
		*m.metaFolder = v
	}

	*m.metaFolder = util.ResolvePath(*m.metaFolder)
	*masterCpuProfile = util.ResolvePath(*masterCpuProfile)
	*masterMemProfile = util.ResolvePath(*masterMemProfile)
	grace.SetupProfiling(*masterCpuProfile, *masterMemProfile)

	parent, _ := util.FullPath(*m.metaFolder).DirAndName()
	if util.FileExists(string(parent)) && !util.FileExists(*m.metaFolder) {
		if err := os.MkdirAll(*m.metaFolder, 0755); err != nil {
			glog.Fatalf("Could not create Meta Folder %s: %v", *m.metaFolder, err)
		}
	}
	if err := util.TestFolderWritable(*m.metaFolder); err != nil {
		glog.Fatalf("Check Meta Folder (-mdir) Writable %s : %s", *m.metaFolder, err)
	}

	masterWhiteList := util.StringSplit(*m.whiteList, ",")
	if *m.volumeSizeLimitMB > util.VolumeSizeLimitGB*1000 {
		glog.Fatalf("volumeSizeLimitMB should be smaller than 30000")
	}

	switch {
	case *m.metricsHttpIp != "":
		// noting to do, use m.metricsHttpIp
	case *m.ipBind != "":
		*m.metricsHttpIp = *m.ipBind
	case *m.ip != "":
		*m.metricsHttpIp = *m.ip
	}
	go stats_collect.StartMetricsServer(*m.metricsHttpIp, *m.metricsHttpPort)
	go stats_collect.LoopPushingMetric("masterServer", util.JoinHostPort(*m.ip, *m.port), *m.metricsAddress, *m.metricsIntervalSec)
	startMaster(m, masterWhiteList)
	return true
}

func startMaster(masterOption MasterOptions, masterWhiteList []string) {

	backend.LoadConfiguration(util.GetViper())

	if *masterOption.portGrpc == 0 {
		*masterOption.portGrpc = 10000 + *masterOption.port
	}
	if *masterOption.ipBind == "" {
		*masterOption.ipBind = *masterOption.ip
	}
	util.SetOutboundLocalIP(*masterOption.ipBind)

	myMasterAddress, peers := checkPeers(*masterOption.ip, *masterOption.port, *masterOption.portGrpc, *masterOption.peers)

	masterPeers := make(map[string]pb.ServerAddress)
	for _, peer := range peers {
		masterPeers[string(peer)] = peer
	}

	r := mux.NewRouter()
	ms := s3server.NewMasterServer(r, masterOption.toMasterOption(masterWhiteList), masterPeers)
	listeningAddress := util.JoinHostPort(*masterOption.ipBind, *masterOption.port)
	glog.V(0).Infof("Start Hanzo S3 Master %s at %s", version.Version(), listeningAddress)
	masterListener, masterLocalListener, e := util.NewIpAndLocalListeners(*masterOption.ipBind, *masterOption.port, 0)
	if e != nil {
		glog.Fatalf("Master startup error: %v", e)
	}

	// Leaderless coordination: no raft server, no election. NewMasterServer
	// installed an in-process LocalCoordinator on the topology; configure it
	// with the real peer membership so the pinned writer is computed over the
	// cluster, then ensure the TopologyId exists. The schain-backed Coordinator
	// (production durable backend, over ZAP) is swapped in here when a schain
	// endpoint is configured.
	if endpoint := schainEndpoint(); endpoint != "" {
		ms.Topo.Coordinator = s3server.NewSchainCoordinator(ms.Topo, endpoint,
			security.LoadClientTLS(util.GetViper(), "grpc.master"))
		glog.V(0).Infof("master coordination backed by schain VM at %s", endpoint)
	}
	ms.ConfigureCoordinator(myMasterAddress, masterPeers)
	r.HandleFunc("/cluster/status", ms.ClusterStatusHandler).Methods(http.MethodGet, http.MethodHead)
	r.HandleFunc("/cluster/healthz", ms.ClusterHealthzHandler).Methods(http.MethodGet, http.MethodHead)

	// starting grpc server — the master service rides the native ZAP transport
	// below; this gRPC server now carries only reflection (the raft consensus
	// service that used to live here is gone).
	grpcPort := *masterOption.portGrpc
	grpcL, grpcLocalL, err := util.NewIpAndLocalListeners(*masterOption.ipBind, grpcPort, 0)
	if err != nil {
		glog.Fatalf("master failed to listen on grpc port %d: %v", grpcPort, err)
	}
	grpcS := pb.NewGrpcServer(security.LoadServerTLS(util.GetViper(), "grpc.master"))
	reflection.Register(grpcS)
	glog.V(0).Infof("Start Hanzo S3 Master %s grpc server at %s:%d", version.Version(), *masterOption.ipBind, grpcPort)
	if grpcLocalL != nil {
		go grpcS.Serve(grpcLocalL)
	}
	go grpcS.Serve(grpcL)
	pb.ServeGrpcOnLocalSocket(grpcS, grpcPort)

	// Serve the WHOLE Hanzo master service (21 unary + 3 streaming RPCs) over the
	// native ZAP transport on the deterministic grpcPort+10000 offset
	// (ServerAddress.ToMasterZapAddress, same convention as ToIamZapAddress).
	// transport.ListenStream carries both the unary dispatch (masterwire.Dispatch
	// over master.NewServerBackend) and the bidirectional streams
	// (masterstream.Handler over master.NewStreamServer) on one listener —
	// exactly as command/filer.go serves the filer. Clients reach it via
	// pb.WithMasterClient over the masterPool. This replaces the legacy gRPC
	// HanzoServer registration: no master RPC rides gRPC anymore.
	//
	// When grpc.master.cert/.key is configured the listener is PQ-secured mTLS:
	// transport.PQTLSConfig pins the X25519MLKEM768 hybrid (PQ X-Wing) and the
	// same cert/CA/allowed-CN gate the legacy gRPC master enforced (LoadServerTLS
	// "grpc.master") applies — no security downgrade. The client mirrors this in
	// pb.dialMasterZapAddr (DialTLS under grpc.master). Otherwise plaintext
	// (loopback / dev), exactly as the gRPC master was plaintext with no cert.
	// Fatal on bind error: the master cannot serve its API without it. The port is
	// derived via pb.ZapPort (overflow-safe) so server and client agree even for
	// high ephemeral ports, and the TLS config is the pb-local ServerTLSConfig.
	masterZapAddr := util.JoinHostPort(*masterOption.ipBind, pb.ZapPort(grpcPort))
	masterDispatch := masterwire.Dispatch(master.NewServerBackend(ms))
	masterStream := masterstream.Handler(master.NewStreamServer(ms))
	var masterZapSrv *transport.Server
	var zapErr error
	if tlsCfg := pb.ServerTLSConfig(util.GetViper(), "grpc.master"); tlsCfg != nil {
		masterZapSrv, zapErr = transport.ListenStreamTLS("tcp", masterZapAddr,
			transport.PQTLSConfig(tlsCfg), masterDispatch, masterStream)
		glog.V(0).Infof("Serving Hanzo S3 Master %s over PQ-TLS (X25519MLKEM768) ZAP transport at %s", version.Version(), masterZapAddr)
	} else {
		masterZapSrv, zapErr = transport.ListenStream("tcp", masterZapAddr, masterDispatch, masterStream)
		glog.V(0).Infof("Start Hanzo S3 Master %s ZAP transport (unary+streaming; plaintext, set grpc.master.cert/.key for PQ-TLS) at %s", version.Version(), masterZapAddr)
	}
	if zapErr != nil {
		glog.Fatalf("master failed to serve over ZAP on %s: %v", masterZapAddr, zapErr)
	}
	grace.OnInterrupt(func() { masterZapSrv.Close() })

	go ms.MasterClient.KeepConnectedToMaster(context.Background())

	// start http server
	var (
		clientCertFile,
		certFile,
		keyFile string
	)
	useTLS := false
	useMTLS := false

	if viper.GetString("https.master.key") != "" {
		useTLS = true
		certFile = viper.GetString("https.master.cert")
		keyFile = viper.GetString("https.master.key")
	}

	if viper.GetString("https.master.ca") != "" {
		useMTLS = true
		clientCertFile = viper.GetString("https.master.ca")
	}

	if masterLocalListener != nil {
		go newHttpServer(r, nil).Serve(masterLocalListener)
	}

	var tlsConfig *tls.Config
	if useMTLS {
		tlsConfig = security.LoadClientTLSHTTP(clientCertFile)
		if err := security.FixTlsConfig(util.GetViper(), tlsConfig); err != nil {
			glog.Fatalf("failed to fix TLS config: %v", err)
		}
	}

	if useTLS {
		getCert, certProvider, err := security.NewReloadingServerCertificate(certFile, keyFile)
		if err != nil {
			glog.Fatalf("failed to load master HTTPS certificate: %v", err)
		}
		// Master runs ServeTLS in a goroutine and this function then blocks
		// on shutdownCtx / select{}; tie the pem refresh goroutine to the
		// existing interrupt hook instead of a local defer that would fire
		// while the server is still running.
		grace.OnInterrupt(certProvider.Close)
		if tlsConfig == nil {
			tlsConfig = &tls.Config{}
		}
		tlsConfig.GetCertificate = getCert
		go newHttpServer(r, tlsConfig).ServeTLS(masterListener, "", "")
	} else {
		go newHttpServer(r, nil).Serve(masterListener)
	}

	grace.OnInterrupt(ms.Shutdown)
	grace.OnInterrupt(grpcS.Stop)
	grace.OnReload(ms.Reload)
	if masterOption.shutdownCtx != nil {
		<-masterOption.shutdownCtx.Done()
		ms.Shutdown()
		grpcS.Stop()
	} else {
		select {}
	}
}

func isSingleMasterMode(peers string) bool {
	p := strings.ToLower(strings.TrimSpace(peers))
	return p == "none"
}

// schainEndpoint returns the configured schain VM endpoint that the production
// Coordinator delegates allocation to over ZAP, or "" to use the in-process
// LocalCoordinator. Configured via master.coordinator.schain_endpoint.
func schainEndpoint() string {
	return strings.TrimSpace(util.GetViper().GetString("master.coordinator.schain_endpoint"))
}

func checkPeers(masterIp string, masterPort int, masterGrpcPort int, peers string) (masterAddress pb.ServerAddress, cleanedPeers []pb.ServerAddress) {
	glog.V(0).Infof("current: %s:%d peers:%s", masterIp, masterPort, peers)
	masterAddress = pb.NewServerAddress(masterIp, masterPort, masterGrpcPort)

	// Handle special case: -peers=none for single-master setup
	if isSingleMasterMode(peers) {
		glog.V(0).Infof("Running in single-master mode (peers=none), no quorum required")
		cleanedPeers = []pb.ServerAddress{masterAddress}
		return
	}

	peers = strings.TrimSpace(peers)
	seenPeers := make(map[string]struct{})
	for _, peer := range pb.ServerAddresses(peers).ToAddresses() {
		normalizedPeer := normalizeMasterPeerAddress(peer, masterAddress)
		key := string(normalizedPeer)
		if _, found := seenPeers[key]; found {
			continue
		}
		seenPeers[key] = struct{}{}
		cleanedPeers = append(cleanedPeers, normalizedPeer)
	}

	hasSelf := false
	for _, peer := range cleanedPeers {
		if peer.ToHttpAddress() == masterAddress.ToHttpAddress() {
			hasSelf = true
			break
		}
	}

	if !hasSelf {
		cleanedPeers = append(cleanedPeers, masterAddress)
	}
	if len(cleanedPeers)%2 == 0 {
		glog.Fatalf("Only odd number of masters are supported: %+v", cleanedPeers)
	}
	return
}

func normalizeMasterPeerAddress(peer pb.ServerAddress, self pb.ServerAddress) pb.ServerAddress {
	if peer.ToHttpAddress() == self.ToHttpAddress() {
		return self
	}

	_, grpcPort, err := net.SplitHostPort(peer.ToGrpcAddress())
	if err != nil {
		return peer
	}
	grpcPortValue, err := strconv.Atoi(grpcPort)
	if err != nil {
		return peer
	}

	return pb.NewServerAddressWithGrpcPort(peer.ToHttpAddress(), grpcPortValue)
}

// peerIndex returns the 0-based position of self in the sorted peer list.
// Peer 0 is the designated bootstrap node. Returns -1 if self is not found.
func peerIndex(self pb.ServerAddress, peers []pb.ServerAddress) int {
	slices.SortFunc(peers, func(a, b pb.ServerAddress) int {
		return strings.Compare(a.ToHttpAddress(), b.ToHttpAddress())
	})
	for i, peer := range peers {
		if peer.ToHttpAddress() == self.ToHttpAddress() {
			return i
		}
	}
	glog.Warningf("peerIndex: self %s not found in peers %v", self, peers)
	return -1
}

func (m *MasterOptions) toMasterOption(whiteList []string) *s3server.MasterOption {
	masterAddress := pb.NewServerAddress(*m.ip, *m.port, *m.portGrpc)
	return &s3server.MasterOption{
		Master:                     masterAddress,
		MetaFolder:                 *m.metaFolder,
		VolumeSizeLimitMB:          uint32(*m.volumeSizeLimitMB),
		VolumePreallocate:          *m.volumePreallocate,
		MaxParallelVacuumPerServer: *m.maxParallelVacuumPerServer,
		// PulseSeconds:            *m.pulseSeconds,
		DefaultReplicaPlacement: *m.defaultReplication,
		GarbageThreshold:        *m.garbageThreshold,
		WhiteList:               whiteList,
		DisableHttp:             *m.disableHttp,
		MetricsAddress:          *m.metricsAddress,
		MetricsIntervalSec:      *m.metricsIntervalSec,
		TelemetryUrl:            *m.telemetryUrl,
		TelemetryEnabled:        *m.telemetryEnabled,
	}
}
