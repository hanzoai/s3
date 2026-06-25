package s3server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/hanzoai/s3/s3/cluster/maintenance"
	"github.com/hanzoai/s3/s3/stats"
	"github.com/hanzoai/s3/s3/telemetry"

	"github.com/hanzoai/s3/s3/cluster"
	"github.com/hanzoai/s3/s3/pb"

	"github.com/gorilla/mux"

	"github.com/hanzoai/s3/s3/glog"
	"github.com/hanzoai/s3/s3/pb/master_pb"
	"github.com/hanzoai/s3/s3/security"
	"github.com/hanzoai/s3/s3/sequence"
	"github.com/hanzoai/s3/s3/shell"
	"github.com/hanzoai/s3/s3/topology"
	"github.com/hanzoai/s3/s3/util"
	util_http "github.com/hanzoai/s3/s3/util/http"
	"github.com/hanzoai/s3/s3/util/version"
	"github.com/hanzoai/s3/s3/wdclient"
	masterwire "github.com/hanzoai/s3/s3/wire/master"
)

const (
	SequencerType        = "master.sequencer.type"
	SequencerSnowflakeId = "master.sequencer.sequencer_snowflake_id"
)

type MasterOption struct {
	Master                     pb.ServerAddress
	MetaFolder                 string
	VolumeSizeLimitMB          uint32
	VolumePreallocate          bool
	MaxParallelVacuumPerServer int
	// PulseSeconds            int
	DefaultReplicaPlacement string
	GarbageThreshold        float64
	WhiteList               []string
	DisableHttp             bool
	MetricsAddress          string
	MetricsIntervalSec      int
	IsFollower              bool
	TelemetryUrl            string
	TelemetryEnabled        bool
	VolumeGrowthDisabled    bool
}

type MasterServer struct {
	option *MasterOption
	guard  *security.Guard

	preallocateSize int64

	Topo                    *topology.Topology
	vg                      *topology.VolumeGrowth
	volumeGrowthRequestChan chan *topology.VolumeGrowRequest

	// notifying clients
	clientChansLock sync.RWMutex
	clientChans     map[string]chan *master_pb.KeepConnectedResponse

	grpcDialOption pb.DialOption

	topologyIdGenLock sync.Mutex

	// coordinatorPeers is the writer-eligible master set the Coordinator pins the
	// writer over. Seeded from -peers and kept live by OnPeerUpdate.
	coordinatorPeersLock sync.Mutex
	coordinatorPeers     map[string]pb.ServerAddress

	MasterClient *wdclient.MasterClient

	adminLocks *AdminLocks

	Cluster *cluster.Cluster

	LockRingManager *cluster.LockRingManager

	// telemetry
	telemetryCollector *telemetry.Collector
}

func NewMasterServer(r *mux.Router, option *MasterOption, peers map[string]pb.ServerAddress) *MasterServer {

	v := util.GetViper()
	signingKey := v.GetString("jwt.signing.key")
	v.SetDefault("jwt.signing.expires_after_seconds", 10)
	expiresAfterSec := v.GetInt("jwt.signing.expires_after_seconds")

	readSigningKey := v.GetString("jwt.signing.read.key")
	v.SetDefault("jwt.signing.read.expires_after_seconds", 60)
	readExpiresAfterSec := v.GetInt("jwt.signing.read.expires_after_seconds")

	v.SetDefault("master.replication.treat_replication_as_minimums", false)
	replicationAsMin := v.GetBool("master.replication.treat_replication_as_minimums")

	v.SetDefault("master.volume_growth.copy_1", topology.VolumeGrowStrategy.Copy1Count)
	v.SetDefault("master.volume_growth.copy_2", topology.VolumeGrowStrategy.Copy2Count)
	v.SetDefault("master.volume_growth.copy_3", topology.VolumeGrowStrategy.Copy3Count)
	v.SetDefault("master.volume_growth.copy_other", topology.VolumeGrowStrategy.CopyOtherCount)
	v.SetDefault("master.volume_growth.threshold", topology.VolumeGrowStrategy.Threshold)
	v.SetDefault("master.volume_growth.disable", false)
	option.VolumeGrowthDisabled = v.GetBool("master.volume_growth.disable")

	topology.VolumeGrowStrategy.Copy1Count = v.GetUint32("master.volume_growth.copy_1")
	topology.VolumeGrowStrategy.Copy2Count = v.GetUint32("master.volume_growth.copy_2")
	topology.VolumeGrowStrategy.Copy3Count = v.GetUint32("master.volume_growth.copy_3")
	topology.VolumeGrowStrategy.CopyOtherCount = v.GetUint32("master.volume_growth.copy_other")
	topology.VolumeGrowStrategy.Threshold = v.GetFloat64("master.volume_growth.threshold")
	whiteList := util.StringSplit(v.GetString("guard.white_list"), ",")

	var preallocateSize int64
	if option.VolumePreallocate {
		preallocateSize = int64(option.VolumeSizeLimitMB) * (1 << 20)
	}

	grpcDialOption := security.LoadClientTLS(v, "grpc.master")
	ms := &MasterServer{
		option:                  option,
		preallocateSize:         preallocateSize,
		volumeGrowthRequestChan: make(chan *topology.VolumeGrowRequest, 1<<6),
		clientChans:             make(map[string]chan *master_pb.KeepConnectedResponse),
		grpcDialOption:          grpcDialOption,
		coordinatorPeers:        make(map[string]pb.ServerAddress, len(peers)),
		MasterClient:            wdclient.NewMasterClient(grpcDialOption, "", cluster.MasterType, option.Master, "", "", *pb.NewServiceDiscoveryFromMap(peers)),
		adminLocks:              NewAdminLocks(),
		Cluster:                 cluster.NewCluster(),
	}
	for name, addr := range peers {
		ms.coordinatorPeers[name] = addr
	}

	ms.LockRingManager = cluster.NewLockRingManager(ms.broadcastToClients)

	ms.MasterClient.SetOnPeerUpdateFn(ms.OnPeerUpdate)

	seq := ms.createSequencer(option)
	if nil == seq {
		glog.Fatalf("create sequencer failed.")
	}
	ms.Topo = topology.NewTopology("topo", seq, uint64(ms.option.VolumeSizeLimitMB)*1024*1024, 5, replicationAsMin)
	ms.vg = topology.NewDefaultVolumeGrowth()
	glog.V(0).Infoln("Volume Size Limit is", ms.option.VolumeSizeLimitMB, "MB")

	// Initialize telemetry after topology is created
	if option.TelemetryEnabled && option.TelemetryUrl != "" {
		telemetryClient := telemetry.NewClient(option.TelemetryUrl, option.TelemetryEnabled)
		ms.telemetryCollector = telemetry.NewCollector(telemetryClient, ms.Topo, ms.Cluster)
		ms.telemetryCollector.SetMasterServer(ms)

		// Set version and OS information
		ms.telemetryCollector.SetVersion(version.VERSION_NUMBER)
		ms.telemetryCollector.SetOS(runtime.GOOS + "/" + runtime.GOARCH)

		// Start periodic telemetry collection (every 24 hours)
		ms.telemetryCollector.StartPeriodicCollection(24 * time.Hour)
	}

	ms.guard = security.NewGuard(append(ms.option.WhiteList, whiteList...), signingKey, expiresAfterSec, readSigningKey, readExpiresAfterSec)

	handleStaticResources2(r)
	r.HandleFunc("/healthz", requestIDMiddleware(ms.healthzHandler)).Methods(http.MethodGet, http.MethodHead)
	r.HandleFunc("/readyz", requestIDMiddleware(ms.readyzHandler)).Methods(http.MethodGet, http.MethodHead)
	r.HandleFunc("/", ms.proxyToLeader(requestIDMiddleware(ms.uiStatusHandler)))
	r.HandleFunc("/ui/index.html", requestIDMiddleware(ms.uiStatusHandler))
	if !ms.option.DisableHttp {
		r.HandleFunc("/dir/assign", ms.proxyToLeader(ms.guard.WhiteList(requestIDMiddleware(ms.dirAssignHandler))))
		r.HandleFunc("/dir/lookup", ms.guard.WhiteList(requestIDMiddleware(ms.dirLookupHandler)))
		r.HandleFunc("/dir/status", ms.proxyToLeader(ms.guard.WhiteList(requestIDMiddleware(ms.dirStatusHandler))))
		r.HandleFunc("/col/delete", ms.proxyToLeader(ms.guard.WhiteList(requestIDMiddleware(ms.collectionDeleteHandler))))
		r.HandleFunc("/vol/grow", ms.proxyToLeader(ms.guard.WhiteList(requestIDMiddleware(ms.volumeGrowHandler))))
		r.HandleFunc("/vol/status", ms.proxyToLeader(ms.guard.WhiteList(requestIDMiddleware(ms.volumeStatusHandler))))
		r.HandleFunc("/vol/vacuum", ms.proxyToLeader(ms.guard.WhiteList(requestIDMiddleware(ms.volumeVacuumHandler))))
		r.HandleFunc("/submit", ms.guard.WhiteList(requestIDMiddleware(ms.submitFromMasterServerHandler)))
		r.HandleFunc("/collection/info", ms.guard.WhiteList(requestIDMiddleware(ms.collectionInfoHandler)))
		/*
			r.HandleFunc("/stats/health", ms.guard.WhiteList(statsHealthHandler))
			r.HandleFunc("/stats/counter", ms.guard.WhiteList(statsCounterHandler))
			r.HandleFunc("/stats/memory", ms.guard.WhiteList(statsMemoryHandler))
		*/
		r.HandleFunc("/{fileId}", requestIDMiddleware(ms.redirectHandler))
	}

	ms.Topo.SetAdminServerConnectedFunc(ms.isAdminServerConnectedFunc)
	ms.Topo.StartRefreshWritableVolumes(
		ms.grpcDialOption,
		ms.option.GarbageThreshold,
		ms.option.MaxParallelVacuumPerServer,
		topology.VolumeGrowStrategy.Threshold,
		ms.preallocateSize,
	)

	ms.ProcessGrowRequest()

	if !option.IsFollower {
		ms.startAdminScripts()
	}

	stats.MasterStartTimeSeconds.Set(float64(time.Now().Unix()))
	return ms
}

func (ms *MasterServer) healthzHandler(w http.ResponseWriter, r *http.Request) {
	// Liveness: process is alive. Keep this fast and simple.
	w.WriteHeader(http.StatusOK)
}

func (ms *MasterServer) readyzHandler(w http.ResponseWriter, r *http.Request) {
	// Readiness: check we can serve traffic.
	leader, err := ms.Topo.Writer()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	if ms.option.Master.Equals(leader) {
		isLocked, err := ms.Topo.IsChildLocked()
		if err != nil {
			glog.Errorf("readyzHandler: %+v", err)
		}
		if isLocked {
			w.WriteHeader(http.StatusLocked)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

// ConfigureCoordinator installs the cluster membership on the Coordinator so the
// pinned writer is computed over the real peer set, and ensures the cluster
// TopologyId exists. It replaces the raft SetRaftServer wiring: there is no
// raft server, no leader-change listener and no election — the writer is a pure
// function of membership.
func (ms *MasterServer) ConfigureCoordinator(self pb.ServerAddress, peers map[string]pb.ServerAddress) {
	members := make([]pb.ServerAddress, 0, len(peers))
	for _, p := range peers {
		members = append(members, p)
	}
	ms.Topo.Coordinator.SetMembers(self, members)
	if ms.Topo.IsWriter() {
		// Seed the warmup timestamp so IsWarmingUp() is active immediately for
		// the writer on startup; non-writers don't need warmup state.
		ms.Topo.SetLastLeaderChangeTime(time.Now())
		glog.V(0).Infof("master %s is the pinned writer", self)
	} else {
		glog.V(0).Infof("master %s defers to pinned writer %s", self, ms.Topo.Coordinator.Writer())
	}
	go ms.ensureTopologyId()
}

// ensureTopologyId mints the cluster identity through the Coordinator. Only the
// pinned writer generates it; others adopt it once observed. Idempotent.
func (ms *MasterServer) ensureTopologyId() {
	ms.topologyIdGenLock.Lock()
	defer ms.topologyIdGenLock.Unlock()
	if id, err := ms.Topo.Coordinator.EnsureTopologyId(); err != nil {
		glog.Errorf("ensureTopologyId: %v", err)
	} else if id != "" {
		glog.V(1).Infof("ensureTopologyId: %s", id)
	}
}

func (ms *MasterServer) proxyToLeader(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if ms.Topo.IsWriter() {
			f(w, r)
			return
		}

		// proxy to the pinned writer
		writerAddr, _ := ms.Topo.MaybeWriter()
		writerLeaderAddr := writerAddr.ToHttpAddress()
		if writerLeaderAddr == "" {
			f(w, r)
			return
		}

		// determine the scheme based on HTTPS client configuration
		scheme := util_http.GetGlobalHttpClient().GetHttpScheme()

		targetUrl, err := url.Parse(scheme + "://" + writerLeaderAddr)
		if err != nil {
			writeJsonError(w, r, http.StatusInternalServerError,
				fmt.Errorf("Leader URL %s://%s Parse Error: %v", scheme, writerLeaderAddr, err))
			return
		}

		// proxy to leader
		glog.V(4).Infoln("proxying to leader", writerLeaderAddr, "using", scheme)
		proxy := httputil.NewSingleHostReverseProxy(targetUrl)
		proxy.Transport = util_http.GetGlobalHttpClient().GetClientTransport()
		proxy.ServeHTTP(w, r)
	}
}

func (ms *MasterServer) isAdminServerConnectedFunc() bool {
	if ms == nil || ms.adminLocks == nil {
		return false
	}
	_, _, isLocked := ms.adminLocks.isLocked(cluster.AdminServerPresenceLockName)
	return isLocked
}

func (ms *MasterServer) startAdminScripts() {
	v := util.GetViper()
	adminScripts := v.GetString("master.maintenance.scripts")
	if adminScripts == "" {
		return
	}
	glog.V(0).Infof("adminScripts: %v", adminScripts)

	sleepMinutes := v.GetFloat64("master.maintenance.sleep_minutes")
	if sleepMinutes <= 0 {
		sleepMinutes = float64(maintenance.DefaultMaintenanceSleepMinutes)
	}

	scriptLines := strings.Split(adminScripts, "\n")
	if !strings.Contains(adminScripts, "lock") {
		scriptLines = append(append([]string{}, "lock"), scriptLines...)
		scriptLines = append(scriptLines, "unlock")
	}

	masterAddress := string(ms.option.Master)

	var shellOptions shell.ShellOptions
	shellOptions.GrpcDialOption = security.LoadClientTLS(v, "grpc.master")
	shellOptions.Masters = &masterAddress

	shellOptions.Directory = "/"
	emptyFilerGroup := ""
	shellOptions.FilerGroup = &emptyFilerGroup

	commandEnv := shell.NewCommandEnv(&shellOptions)

	reg, _ := regexp.Compile(`'.*?'|".*?"|\S+`)

	go commandEnv.MasterClient.KeepConnectedToMaster(context.Background())

	go func() {
		for {
			time.Sleep(time.Duration(sleepMinutes) * time.Minute)
			if ms.Topo.IsWriter() && ms.MasterClient.GetMaster(context.Background()) != "" {
				if ms.isAdminServerConnectedFunc() {
					glog.V(1).Infof("Skipping master maintenance scripts because admin server is connected")
					continue
				}
				shellOptions.FilerAddress = ms.GetOneFiler(cluster.FilerGroupName(*shellOptions.FilerGroup))
				if shellOptions.FilerAddress == "" {
					continue
				}
				for _, line := range scriptLines {
					for _, c := range strings.Split(line, ";") {
						processEachCmd(reg, c, commandEnv)
					}
				}
			}
		}
	}()
}

func processEachCmd(reg *regexp.Regexp, line string, commandEnv *shell.CommandEnv) {
	cmds := reg.FindAllString(line, -1)
	if len(cmds) == 0 {
		return
	}
	args := make([]string, len(cmds[1:]))
	for i := range args {
		args[i] = strings.Trim(string(cmds[1+i]), "\"'")
	}
	cmd := cmds[0]

	for _, c := range shell.Commands {
		if c.Name() == cmd {
			if c.HasTag(shell.ResourceHeavy) {
				glog.Warningf("%s is resource heavy and should not run on master", cmd)
				continue
			}
			glog.V(0).Infof("executing: %s %v", cmd, args)
			if err := c.Do(args, commandEnv, os.Stdout); err != nil {
				glog.V(0).Infof("error: %v", err)
			}
		}
	}
}

func (ms *MasterServer) createSequencer(option *MasterOption) sequence.Sequencer {
	var seq sequence.Sequencer
	v := util.GetViper()
	seqType := strings.ToLower(v.GetString(SequencerType))
	glog.V(1).Infof("[%s] : [%s]", SequencerType, seqType)
	switch strings.ToLower(seqType) {
	case "snowflake":
		var err error
		snowflakeId := v.GetInt(SequencerSnowflakeId)
		seq, err = sequence.NewSnowflakeSequencer(string(option.Master), snowflakeId)
		if err != nil {
			glog.Error(err)
			seq = nil
		}
	case "raft":
		fallthrough
	default:
		seq = sequence.NewMemorySequencer()
	}
	return seq
}

// OnPeerUpdate keeps the Coordinator's writer-eligible membership in sync with
// the live master set. A joining master is added; a departing one is dropped
// only after it fails a ZAP liveness probe, so a transient blip does not re-pin
// the writer. The pinned writer is recomputed (lowest live address) on every
// change — this is the leaderless failover that replaces raft re-election.
func (ms *MasterServer) OnPeerUpdate(update *master_pb.ClusterNodeUpdate, startFrom time.Time) {
	if update.NodeType != cluster.MasterType {
		return
	}
	glog.V(4).Infof("OnPeerUpdate: %+v", update)

	peerAddress := pb.ServerAddress(update.Address)
	ms.coordinatorPeersLock.Lock()
	if update.IsAdd {
		ms.coordinatorPeers[string(peerAddress)] = peerAddress
	} else {
		// Probe the departing peer over the native ZAP transport (the whole master
		// service is served over ZAP); only drop it when it is truly unreachable.
		pingFailed := false
		if zapClient, dialErr := masterwire.Dial("tcp", peerAddress.ToMasterZapAddress()); dialErr != nil {
			pingFailed = true
		} else {
			_, pErr := zapClient.Ping(context.Background(), masterwire.PingRequestInput{Target: string(peerAddress), TargetType: cluster.MasterType})
			zapClient.Close()
			pingFailed = pErr != nil
		}
		if pingFailed {
			delete(ms.coordinatorPeers, string(peerAddress))
			glog.V(0).Infof("master %s unreachable; dropped from writer-eligible set", peerAddress)
		} else {
			glog.V(0).Infof("master %s responded to ping; keeping in writer-eligible set", peerAddress)
		}
	}
	members := make([]pb.ServerAddress, 0, len(ms.coordinatorPeers))
	for _, p := range ms.coordinatorPeers {
		members = append(members, p)
	}
	ms.coordinatorPeersLock.Unlock()
	ms.Topo.Coordinator.SetMembers(ms.option.Master, members)
}

// Shutdown is a no-op: the leaderless Coordinator holds no durable raft state to
// flush and there is no leadership to transfer. Kept for the grace.OnInterrupt
// call site.
func (ms *MasterServer) Shutdown() {}

func (ms *MasterServer) Reload() {
	glog.V(0).Infoln("Reload master server...")

	util.LoadConfiguration("security", false)
	v := util.GetViper()
	ms.guard.UpdateWhiteList(append(ms.option.WhiteList,
		util.StringSplit(v.GetString("guard.white_list"), ",")...),
	)
	ms.guard.UpdateSigningKeys(
		v.GetString("jwt.signing.key"),
		v.GetInt("jwt.signing.expires_after_seconds"),
		v.GetString("jwt.signing.read.key"),
		v.GetInt("jwt.signing.read.expires_after_seconds"),
	)
}
