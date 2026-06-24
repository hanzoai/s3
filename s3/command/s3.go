package command

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/hanzoai/s3/s3/glog"
	"github.com/hanzoai/s3/s3/pb"
	"github.com/hanzoai/s3/s3/pb/filer_pb"
	"github.com/hanzoai/s3/s3/s3api"
	"github.com/hanzoai/s3/s3/s3api/s3err"
	"github.com/hanzoai/s3/s3/security"
	stats_collect "github.com/hanzoai/s3/s3/stats"
	"github.com/hanzoai/s3/s3/util"
	"github.com/hanzoai/s3/s3/util/grace"
	"github.com/hanzoai/s3/s3/util/version"
	"github.com/hanzoai/s3/s3/wire/iamadapt"
	s3wire "github.com/hanzoai/s3/s3/wire/s3"
	s3_lifecyclewire "github.com/hanzoai/s3/s3/wire/s3_lifecycle"
	"github.com/hanzoai/s3/s3/zapsvc"
	"github.com/zap-proto/go/transport"
)

var (
	s3StandaloneOptions S3Options
)

// S3Options holds CLI flags for the S3 gateway.
// Flags are registered in multiple commands: s3.go (standalone), server.go, filer.go, and mini.go.
// When adding a new field, update all four flag registration sites.
type S3Options struct {
	filer                     *string
	bindIp                    *string
	port                      *int
	portHttps                 *int
	portGrpc                  *int
	portIceberg               *int
	zapPort                   *int
	config                    *string
	iamConfig                 *string
	domainName                *string
	allowedOrigins            *string
	tlsPrivateKey             *string
	tlsCertificate            *string
	tlsCACertificate          *string
	tlsVerifyClientCert       *bool
	metricsHttpPort           *int
	metricsHttpIp             *string
	allowDeleteBucketNotEmpty *bool
	auditLogConfig            *string
	localFilerSocket          *string
	dataCenter                *string
	localSocket               *string
	idleTimeout               *int
	concurrentUploadLimitMB   *int
	concurrentFileUploadLimit *int
	enableIam                 *bool
	iamReadOnly               *bool
	debug                     *bool
	debugPort                 *int
	cipher                    *bool
	externalUrl               *string
	defaultFileMode           *string
	cacheSizeMB               *int64
	// shutdownCtx, when non-nil, tells startS3Server/startIcebergServer to
	// gracefully shut down their HTTP/gRPC servers once the ctx is cancelled.
	// Used by s3 mini to orchestrate an ordered shutdown; nil for standalone
	// s3 s3.
	shutdownCtx context.Context
}

func init() {
	cmdS3.Run = runS3 // break init cycle
	s3StandaloneOptions.filer = cmdS3.Flag.String("filer", "localhost:8888", "comma-separated filer server addresses for high availability")
	s3StandaloneOptions.bindIp = cmdS3.Flag.String("ip.bind", "", "ip address to bind to. If empty, default to 0.0.0.0.")
	s3StandaloneOptions.port = cmdS3.Flag.Int("port", 8333, "s3 server http listen port")
	s3StandaloneOptions.portHttps = cmdS3.Flag.Int("port.https", 0, "s3 server https listen port")
	s3StandaloneOptions.portGrpc = cmdS3.Flag.Int("port.grpc", 0, "s3 server grpc listen port")
	s3StandaloneOptions.portIceberg = cmdS3.Flag.Int("port.iceberg", 8181, "Iceberg REST Catalog server listen port (0 to disable)")
	s3StandaloneOptions.zapPort = cmdS3.Flag.Int("zap.port", 0, "native ZAP service-mesh listen port for in-cloud callers (0 disables; external S3 clients use the HTTPS S3 API)")
	s3StandaloneOptions.domainName = cmdS3.Flag.String("domainName", "", "suffix of the host name in comma separated list, {bucket}.{domainName}")
	s3StandaloneOptions.allowedOrigins = cmdS3.Flag.String("allowedOrigins", "*", "comma separated list of allowed origins")
	s3StandaloneOptions.dataCenter = cmdS3.Flag.String("dataCenter", "", "prefer to read and write to volumes in this data center")
	s3StandaloneOptions.config = cmdS3.Flag.String("config", "", "path to the config file")
	s3StandaloneOptions.iamConfig = cmdS3.Flag.String("iam.config", "", "path to the advanced IAM config file")
	s3StandaloneOptions.auditLogConfig = cmdS3.Flag.String("auditLogConfig", "", "path to the audit log config file")
	s3StandaloneOptions.tlsPrivateKey = cmdS3.Flag.String("key.file", "", "path to the TLS private key file")
	s3StandaloneOptions.tlsCertificate = cmdS3.Flag.String("cert.file", "", "path to the TLS certificate file")
	s3StandaloneOptions.tlsCACertificate = cmdS3.Flag.String("cacert.file", "", "path to the TLS CA certificate file")
	s3StandaloneOptions.tlsVerifyClientCert = cmdS3.Flag.Bool("tlsVerifyClientCert", false, "whether to verify the client's certificate")
	s3StandaloneOptions.metricsHttpPort = cmdS3.Flag.Int("metricsPort", 0, "Prometheus metrics listen port")
	s3StandaloneOptions.metricsHttpIp = cmdS3.Flag.String("metricsIp", "", "metrics listen ip. If empty, default to same as -ip.bind option.")
	cmdS3.Flag.Bool("allowEmptyFolder", true, "deprecated, ignored. Empty folder cleanup is now automatic.")
	s3StandaloneOptions.allowDeleteBucketNotEmpty = cmdS3.Flag.Bool("allowDeleteBucketNotEmpty", true, "allow recursive deleting all entries along with bucket")
	s3StandaloneOptions.localFilerSocket = cmdS3.Flag.String("localFilerSocket", "", "local filer socket path")
	s3StandaloneOptions.localSocket = cmdS3.Flag.String("localSocket", "", "default to /tmp/hanzo-s3-<port>.sock")
	s3StandaloneOptions.idleTimeout = cmdS3.Flag.Int("idleTimeout", 120, "connection idle seconds")
	s3StandaloneOptions.concurrentUploadLimitMB = cmdS3.Flag.Int("concurrentUploadLimitMB", 0, "limit total concurrent upload size, 0 means unlimited")
	s3StandaloneOptions.concurrentFileUploadLimit = cmdS3.Flag.Int("concurrentFileUploadLimit", 0, "limit number of concurrent file uploads, 0 means unlimited")
	s3StandaloneOptions.enableIam = cmdS3.Flag.Bool("iam", true, "enable embedded IAM API on the same port")
	s3StandaloneOptions.iamReadOnly = cmdS3.Flag.Bool("iam.readOnly", true, "disable IAM write operations on this server")
	s3StandaloneOptions.debug = cmdS3.Flag.Bool("debug", false, "serves runtime profiling data via pprof on the port specified by -debug.port")
	s3StandaloneOptions.debugPort = cmdS3.Flag.Int("debug.port", 6060, "http port for debugging")
	s3StandaloneOptions.cipher = cmdS3.Flag.Bool("encryptVolumeData", false, "encrypt data on volume servers")
	s3StandaloneOptions.externalUrl = cmdS3.Flag.String("externalUrl", "", "the external URL clients use to connect (e.g. https://api.example.com:9000). Used for S3 signature verification behind a reverse proxy. Falls back to S3_EXTERNAL_URL env var.")
	s3StandaloneOptions.defaultFileMode = cmdS3.Flag.String("defaultFileMode", "", "default file mode for S3 uploaded objects, e.g. 0660, 0644, 0666")
	s3StandaloneOptions.cacheSizeMB = cmdS3.Flag.Int64("cacheCapacityMB", 0, "in-memory chunk cache capacity in MB for S3 GETs shared across requests (0 disables)")
}

var cmdS3 = &Command{
	UsageLine: "s3 [-port=8333] [-filer=<ip:port>[,<ip:port>]...] [-config=</path/to/config.json>]",
	Short:     "start a s3 API compatible server that is backed by filer(s)",
	Long: `start a s3 API compatible server that is backed by filer(s).

	Multiple filer addresses can be specified for high availability, separated by commas.
	The S3 server will automatically failover between filers if one becomes unavailable.

	By default, you can use any access key and secret key to access the S3 APIs.
	To enable credential based access, create a config.json file similar to this:

{
  "identities": [
    {
      "name": "anonymous",
      "actions": [
        "Read"
      ]
    },
    {
      "name": "some_admin_user",
      "credentials": [
        {
          "accessKey": "some_access_key1",
          "secretKey": "some_secret_key1"
        }
      ],
      "actions": [
        "Admin",
        "Read",
        "List",
        "Tagging",
        "Write"
      ]
    },
    {
      "name": "some_read_only_user",
      "credentials": [
        {
          "accessKey": "some_access_key2",
          "secretKey": "some_secret_key2"
        }
      ],
      "actions": [
        "Read"
      ]
    },
    {
      "name": "some_normal_user",
      "credentials": [
        {
          "accessKey": "some_access_key3",
          "secretKey": "some_secret_key3"
        }
      ],
      "actions": [
        "Read",
        "List",
        "Tagging",
        "Write"
      ]
    },
    {
      "name": "user_limited_to_bucket1",
      "credentials": [
        {
          "accessKey": "some_access_key4",
          "secretKey": "some_secret_key4"
        }
      ],
      "actions": [
        "Read:bucket1",
        "List:bucket1",
        "Tagging:bucket1",
        "Write:bucket1"
      ]
    }
  ]
}

	Alternatively, you can use environment variables as fallback admin credentials:

	AWS_ACCESS_KEY_ID=your_access_key AWS_SECRET_ACCESS_KEY=your_secret_key s3 s3

	Environment variables are only used when no S3 configuration file is provided
	and no configuration is available from the filer. This provides a simple way
	to get started without requiring configuration files.

`,
}

func runS3(cmd *Command, args []string) bool {
	if *s3StandaloneOptions.debug {
		grace.StartDebugServer(*s3StandaloneOptions.debugPort)
	}

	s3StandaloneOptions.resolvePaths()
	util.LoadSecurityConfiguration()

	switch {
	case *s3StandaloneOptions.metricsHttpIp != "":
		// noting to do, use s3StandaloneOptions.metricsHttpIp
	case *s3StandaloneOptions.bindIp != "":
		*s3StandaloneOptions.metricsHttpIp = *s3StandaloneOptions.bindIp
	}
	go stats_collect.StartMetricsServer(*s3StandaloneOptions.metricsHttpIp, *s3StandaloneOptions.metricsHttpPort)

	return s3StandaloneOptions.startS3Server()

}

// resolveExternalUrl returns the external URL from the flag or falls back to the S3_EXTERNAL_URL env var.
func (s3opt *S3Options) resolveExternalUrl() string {
	if s3opt.externalUrl != nil && *s3opt.externalUrl != "" {
		return *s3opt.externalUrl
	}
	return os.Getenv("S3_EXTERNAL_URL")
}

func (s3opt *S3Options) parseDefaultFileMode() (uint32, error) {
	if s3opt.defaultFileMode == nil || *s3opt.defaultFileMode == "" {
		return 0, nil
	}
	mode, err := strconv.ParseUint(*s3opt.defaultFileMode, 8, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid defaultFileMode %q: %v", *s3opt.defaultFileMode, err)
	}
	return uint32(mode), nil
}

// resolvePaths expands "~" in every user-supplied path flag so callers
// that share these pointers (e.g. server.go propagating s3Options.config
// to filerOptions.s3ConfigFile before startS3Server runs) see resolved
// values. Idempotent — safe to call from any entry point.
func (s3opt *S3Options) resolvePaths() {
	*s3opt.config = util.ResolvePath(*s3opt.config)
	*s3opt.iamConfig = util.ResolvePath(*s3opt.iamConfig)
	*s3opt.tlsCertificate = util.ResolvePath(*s3opt.tlsCertificate)
	*s3opt.tlsPrivateKey = util.ResolvePath(*s3opt.tlsPrivateKey)
	*s3opt.tlsCACertificate = util.ResolvePath(*s3opt.tlsCACertificate)
	*s3opt.auditLogConfig = util.ResolvePath(*s3opt.auditLogConfig)
}

func (s3opt *S3Options) startS3Server() bool {

	// Before the first filer dial below; gRPC caches conns and binds at dial time.
	util.SetOutboundLocalIP(*s3opt.bindIp)

	filerAddresses := pb.ServerAddresses(*s3opt.filer).ToAddresses()

	filerBucketsPath := "/buckets"
	filerGroup := ""
	var masterAddresses []pb.ServerAddress

	grpcDialOption := security.LoadClientTLS(util.GetViper(), "grpc.client")

	// Native ZAP S3 service for the Hanzo mesh (platform law: internal
	// service↔service = ZAP only). Off by default; -zap.port>0 enables it.
	// External S3 clients keep the HTTPS S3 API below — this is the in-cloud
	// service surface, backed by the same filer.
	if *s3opt.zapPort > 0 {
		store := zapsvc.NewFilerStore(filerAddresses, grpcDialOption)
		// Non-fatal: the ZAP mesh surface must never take down the critical S3
		// HTTP path. On failure we log and keep serving HTTPS S3.
		if _, err := zapsvc.Serve("tcp", fmt.Sprintf(":%d", *s3opt.zapPort), store); err != nil {
			glog.Errorf("native ZAP S3 service on port %d failed to start (S3 HTTP unaffected): %v", *s3opt.zapPort, err)
		} else {
			glog.V(0).Infof("Native ZAP S3 service listening on port %d (mesh: internal callers use ZAP)", *s3opt.zapPort)
		}
	}

	// metrics read from the filer
	var metricsAddress string
	var metricsIntervalSec int

	for {
		err := pb.WithOneOfGrpcFilerClients(false, filerAddresses, grpcDialOption, func(client filer_pb.HanzoFilerClient) error {
			resp, err := client.GetFilerConfiguration(context.Background(), &filer_pb.GetFilerConfigurationRequest{})
			if err != nil {
				return fmt.Errorf("get filer configuration: %v", err)
			}
			filerBucketsPath = resp.DirBuckets
			filerGroup = resp.FilerGroup
			// Get master addresses for filer discovery
			masterAddresses = pb.ServerAddresses(strings.Join(resp.Masters, ",")).ToAddresses()
			metricsAddress, metricsIntervalSec = resp.MetricsAddress, int(resp.MetricsIntervalSec)
			glog.V(0).Infof("S3 read filer buckets dir: %s", filerBucketsPath)
			if len(masterAddresses) > 0 {
				glog.V(0).Infof("S3 read master addresses for discovery: %v", masterAddresses)
			}
			return nil
		})
		if err != nil {
			glog.V(2).Infof("wait to connect to filers %v grpc address", filerAddresses)
			time.Sleep(time.Second)
		} else {
			glog.V(0).Infof("connected to filers %v", filerAddresses)
			break
		}
	}

	go stats_collect.LoopPushingMetric("s3", stats_collect.SourceName(uint32(*s3opt.port)), metricsAddress, metricsIntervalSec)

	router := mux.NewRouter().SkipClean(true)
	var localFilerSocket string
	if s3opt.localFilerSocket != nil {
		localFilerSocket = *s3opt.localFilerSocket
	}
	var s3ApiServer *s3api.S3ApiServer
	var s3ApiServer_err error

	// Create S3 server with optional advanced IAM integration
	var iamConfigPath string
	if s3opt.iamConfig != nil && *s3opt.iamConfig != "" {
		iamConfigPath = *s3opt.iamConfig
		glog.V(0).Infof("Starting S3 API Server with advanced IAM integration")
	} else {
		glog.V(0).Infof("Starting S3 API Server with standard IAM")
	}

	if *s3opt.portGrpc == 0 {
		*s3opt.portGrpc = 10000 + *s3opt.port
	}
	if *s3opt.bindIp == "" {
		*s3opt.bindIp = "0.0.0.0"
	}

	defaultFileMode, fileModeErr := s3opt.parseDefaultFileMode()
	if fileModeErr != nil {
		glog.Fatalf("S3 API Server startup error: %v", fileModeErr)
	}

	s3ApiServer, s3ApiServer_err = s3api.NewS3ApiServer(router, &s3api.S3ApiServerOption{
		Filers:                    filerAddresses,
		Masters:                   masterAddresses,
		Port:                      *s3opt.port,
		Config:                    *s3opt.config,
		DomainName:                *s3opt.domainName,
		AllowedOrigins:            strings.Split(*s3opt.allowedOrigins, ","),
		BucketsPath:               filerBucketsPath,
		GrpcDialOption:            grpcDialOption,
		AllowDeleteBucketNotEmpty: *s3opt.allowDeleteBucketNotEmpty,
		LocalFilerSocket:          localFilerSocket,
		DataCenter:                *s3opt.dataCenter,
		FilerGroup:                filerGroup,
		IamConfig:                 iamConfigPath, // Advanced IAM config (optional)
		ConcurrentUploadLimit:     int64(*s3opt.concurrentUploadLimitMB) * 1024 * 1024,
		ConcurrentFileUploadLimit: int64(*s3opt.concurrentFileUploadLimit),
		EnableIam:                 *s3opt.enableIam, // Embedded IAM API (enabled by default)
		IamReadOnly:               *s3opt.iamReadOnly,
		Cipher:                    *s3opt.cipher, // encrypt data on volume servers
		BindIp:                    *s3opt.bindIp,
		GrpcPort:                  *s3opt.portGrpc,
		ExternalUrl:               s3opt.resolveExternalUrl(),
		DefaultFileMode:           defaultFileMode,
		CacheSizeMB:               *s3opt.cacheSizeMB,
	})
	if s3ApiServer_err != nil {
		glog.Fatalf("S3 API Server startup error: %v", s3ApiServer_err)
	}
	defer s3ApiServer.Shutdown()

	if runtime.GOOS != "windows" {
		localSocket := *s3opt.localSocket
		if localSocket == "" {
			localSocket = fmt.Sprintf("/tmp/hanzo-s3-%d.sock", *s3opt.port)
		}
		if err := os.Remove(localSocket); err != nil && !os.IsNotExist(err) {
			glog.Fatalf("Failed to remove %s, error: %s", localSocket, err.Error())
		}
		go func() {
			// start on local unix socket
			s3SocketListener, err := net.Listen("unix", localSocket)
			if err != nil {
				glog.Fatalf("Failed to listen on %s: %v", localSocket, err)
			}
			if err := newHttpServer(router, nil).Serve(s3SocketListener); err != nil && err != http.ErrServerClosed {
				glog.Fatalf("Failed to start S3 http server: %v", err)
			}
		}()
	}

	listenAddress := fmt.Sprintf("%s:%d", *s3opt.bindIp, *s3opt.port)
	s3ApiListener, s3ApiLocalListener, err := util.NewIpAndLocalListeners(
		*s3opt.bindIp, *s3opt.port, time.Duration(*s3opt.idleTimeout)*time.Second)
	if err != nil {
		glog.Fatalf("S3 API Server listener on %s error: %v", listenAddress, err)
	}

	if len(*s3opt.auditLogConfig) > 0 {
		s3err.InitAuditLog(*s3opt.auditLogConfig)
		if s3err.Logger != nil {
			defer s3err.Logger.Close()
		}
	}

	// HanzoS3LifecycleInternal now speaks the native ZAP transport on the same
	// port the gRPC server used to bind (the lifecycle worker / shell dial it
	// directly). The S3ApiServer is the live LifecycleDeleter backend; the
	// re-fetch + identity-CAS + object-lock business logic reads the zero-copy
	// wire request view directly — no protobuf, no bridge.
	grpcPort := *s3opt.portGrpc
	grpcL, grpcLocalL, err := util.NewIpAndLocalListeners(*s3opt.bindIp, grpcPort, 0)
	if err != nil {
		glog.Fatalf("s3 failed to listen on grpc port %d: %v", grpcPort, err)
	}
	lifecycleDispatch := s3_lifecyclewire.Dispatch(s3ApiServer)
	lifecycleSrv := transport.Serve(grpcL, lifecycleDispatch)
	if grpcLocalL != nil {
		_ = transport.Serve(grpcLocalL, lifecycleDispatch)
	}

	// HanzoS3IamCache runs on its own ZAP transport listener (filer -> S3 cache
	// propagation), no longer on the shared gRPC server. Its port is derived
	// from the gRPC port so the propagating filer can find it via the cluster
	// node list. Non-fatal: a cache-listener failure must never take down S3.
	iamCacheAddr := iamadapt.CacheZapAddress(util.JoinHostPort(*s3opt.bindIp, grpcPort))
	if iamCacheL, iamCacheErr := net.Listen("tcp", iamCacheAddr); iamCacheErr != nil {
		glog.Errorf("s3 IAM cache ZAP listener on %s failed to start (S3 unaffected): %v", iamCacheAddr, iamCacheErr)
	} else {
		_ = transport.Serve(iamCacheL, func(env []byte) ([]byte, error) {
			return s3wire.DispatchHanzoS3IamCache(s3ApiServer, env)
		})
		glog.V(0).Infof("HanzoS3IamCache ZAP service listening on %s", iamCacheAddr)
	}

	if *s3opt.tlsPrivateKey != "" {
		// Check for port conflict when both HTTP and HTTPS are enabled on the same port
		if *s3opt.portHttps > 0 && *s3opt.portHttps == *s3opt.port {
			glog.Fatalf("S3 API Server error: -s3.port.https (%d) cannot be the same as -s3.port (%d)", *s3opt.portHttps, *s3opt.port)
		}

		getCert, certProvider, err := security.NewReloadingServerCertificate(*s3opt.tlsCertificate, *s3opt.tlsPrivateKey)
		if err != nil {
			glog.Fatalf("S3 API Server failed to load HTTPS certificate: %v", err)
		}
		grace.OnInterrupt(certProvider.Close)

		caCertPool := x509.NewCertPool()
		if *s3opt.tlsCACertificate != "" {
			// load CA certificate file and add it to list of client CAs
			caCertFile, err := ioutil.ReadFile(*s3opt.tlsCACertificate)
			if err != nil {
				glog.Fatalf("error reading CA certificate: %v", err)
			}
			caCertPool.AppendCertsFromPEM(caCertFile)
		}

		clientAuth := tls.NoClientCert
		if *s3opt.tlsVerifyClientCert {
			clientAuth = tls.RequireAndVerifyClientCert
		}

		tlsConfig := &tls.Config{
			GetCertificate: getCert,
			ClientAuth:     clientAuth,
			ClientCAs:      caCertPool,
		}
		err = security.FixTlsConfig(util.GetViper(), tlsConfig)
		if err != nil {
			glog.Fatalf("error with tls config: %v", err)
		}
		if *s3opt.portHttps == 0 {
			glog.V(0).Infof("Start Hanzo S3 API Server %s at https port %d", version.Version(), *s3opt.port)
			if s3ApiLocalListener != nil {
				go func() {
					if err = newHttpServer(router, tlsConfig).ServeTLS(s3ApiLocalListener, "", ""); err != nil {
						glog.Fatalf("S3 API Server Fail to serve: %v", err)
					}
				}()
			}
			httpS := newHttpServer(router, tlsConfig)
			if s3opt.shutdownCtx != nil {
				go func() {
					<-s3opt.shutdownCtx.Done()
					httpS.Shutdown(context.Background())
					lifecycleSrv.Close()
				}()
			}
			if err = httpS.ServeTLS(s3ApiListener, "", ""); err != nil && err != http.ErrServerClosed {
				glog.Fatalf("S3 API Server Fail to serve: %v", err)
			}
		} else {
			glog.V(0).Infof("Start Hanzo S3 API Server %s at https port %d", version.Version(), *s3opt.portHttps)
			s3ApiListenerHttps, s3ApiLocalListenerHttps, err := util.NewIpAndLocalListeners(
				*s3opt.bindIp, *s3opt.portHttps, time.Duration(*s3opt.idleTimeout)*time.Second)
			if err != nil {
				glog.Fatalf("S3 API HTTPS listener on %s:%d error: %v", *s3opt.bindIp, *s3opt.portHttps, err)
			}
			if s3ApiLocalListenerHttps != nil {
				go func() {
					if err = newHttpServer(router, tlsConfig).ServeTLS(s3ApiLocalListenerHttps, "", ""); err != nil {
						glog.Fatalf("S3 API Server Fail to serve: %v", err)
					}
				}()
			}
			go func() {
				if err = newHttpServer(router, tlsConfig).ServeTLS(s3ApiListenerHttps, "", ""); err != nil {
					glog.Fatalf("S3 API Server Fail to serve: %v", err)
				}
			}()
		}
	}
	if *s3opt.tlsPrivateKey == "" || *s3opt.portHttps > 0 {
		glog.V(0).Infof("Start Hanzo S3 API Server %s at http port %d", version.Version(), *s3opt.port)
		if s3ApiLocalListener != nil {
			go func() {
				if err = newHttpServer(router, nil).Serve(s3ApiLocalListener); err != nil {
					glog.Fatalf("S3 API Server Fail to serve: %v", err)
				}
			}()
		}
		httpS := newHttpServer(router, nil)
		if s3opt.shutdownCtx != nil {
			go func() {
				<-s3opt.shutdownCtx.Done()
				httpS.Shutdown(context.Background())
				lifecycleSrv.Close()
			}()
		}
		if err = httpS.Serve(s3ApiListener); err != nil && err != http.ErrServerClosed {
			glog.Fatalf("S3 API Server Fail to serve: %v", err)
		}
	}

	return true

}

// deriveS3AdvertisedEndpoint builds the S3 endpoint URL to advertise to
// Iceberg catalog clients as part of LoadTable FileIO config. To avoid
// hijacking correctly-configured clients (Spark/Trino/PyIceberg all bring
// their own s3.endpoint), advertising is strictly opt-in and returns ""
// whenever no reliable value is available:
//   - -s3.externalUrl / S3_EXTERNAL_URL wins and supports reverse-proxy
//     deployments.
//   - Otherwise the bind IP is used only when it is explicit and not a
//     wildcard (0.0.0.0 / ::), with the scheme picked from TLS config and
//     IPv6 literals bracketed via util.JoinHostPort.
//
// See issue #9103.
func (s3opt *S3Options) deriveS3AdvertisedEndpoint() string {
	if ext := strings.TrimRight(s3opt.resolveExternalUrl(), "/"); ext != "" {
		return ext
	}

	host := ""
	if s3opt.bindIp != nil {
		host = *s3opt.bindIp
	}
	switch host {
	case "", "0.0.0.0", "::", "[::]":
		return ""
	}

	scheme := "http"
	port := 0
	if s3opt.port != nil {
		port = *s3opt.port
	}
	if s3opt.tlsPrivateKey != nil && *s3opt.tlsPrivateKey != "" {
		scheme = "https"
		if s3opt.portHttps != nil && *s3opt.portHttps > 0 {
			port = *s3opt.portHttps
		}
	}
	return fmt.Sprintf("%s://%s", scheme, util.JoinHostPort(host, port))
}
