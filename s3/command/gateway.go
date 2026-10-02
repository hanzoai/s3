package command

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/hanzoai/s3/s3/credential"
	_ "github.com/hanzoai/s3/s3/credential/memory"
	"github.com/hanzoai/s3/s3/gateway"
	"github.com/hanzoai/s3/s3/glog"
	"github.com/hanzoai/s3/s3/s3api"
)

var gw struct {
	bind        *string
	port        *int
	config      *string
	upstream    *string
	region      *string
	credentials *string
	prefix      *string
	suffix      *string
	drain       *time.Duration
}

var cmdGateway = &Command{
	UsageLine: "gateway -config=/config/s3.json -upstream=https://s3.us-east-1.amazonaws.com -upstream.credentials=/upstream -bucket.prefix=hanzo-s3- -bucket.suffix=-<account>",
	Short:     "serve the S3 API from an upstream object store, keeping nothing locally",
	Long: `gateway serves the S3 API on -port and stores every object in an upstream
S3 service. It keeps no data, so any number of gateways serve the same buckets
and a gateway can be lost, replaced or moved without losing a byte.

A request is authenticated against the identities in -config (the same s3.json
"s3 server" reads) and authorized for the action it performs. It is then signed
with the key in -upstream.credentials (files access-key and secret-key) and sent
to the upstream under the bucket's mapped name: a client's "org-db" is the
upstream's <bucket.prefix>org-db<bucket.suffix>, since upstream bucket names are
global.

The gateway answers path-style requests (http://host:9000/<bucket>/<key>) and
probes on /_health.

	s3 gateway -config=/config/s3.json -upstream.credentials=/upstream \
	  -bucket.prefix=hanzo-s3- -bucket.suffix=-532217001883
`,
}

func init() {
	cmdGateway.Run = runGateway
	gw.bind = cmdGateway.Flag.String("ip.bind", "0.0.0.0", "address to listen on")
	gw.port = cmdGateway.Flag.Int("port", 9000, "port the S3 API listens on")
	gw.config = cmdGateway.Flag.String("config", "", "identities file (s3.json); required")
	gw.upstream = cmdGateway.Flag.String("upstream", "https://s3.us-east-1.amazonaws.com", "upstream S3 endpoint; a bucket is addressed as <bucket>.<host>")
	gw.region = cmdGateway.Flag.String("region", "us-east-1", "upstream region")
	gw.credentials = cmdGateway.Flag.String("upstream.credentials", "", "directory holding the upstream key as files access-key and secret-key; required")
	gw.prefix = cmdGateway.Flag.String("bucket.prefix", "", "prefix of every upstream bucket name")
	gw.suffix = cmdGateway.Flag.String("bucket.suffix", "", "suffix of every upstream bucket name")
	gw.drain = cmdGateway.Flag.Duration("drain", 5*time.Second, "how long to keep serving after SIGTERM, while the endpoint is withdrawn")
}

func runGateway(cmd *Command, args []string) bool {
	if *gw.config == "" || *gw.credentials == "" {
		cmd.Usage()
		return false
	}
	key, err := readUpstreamKey(*gw.credentials)
	if err != nil {
		glog.Fatalf("upstream key: %v", err)
	}
	up, err := gateway.NewUpstream(*gw.upstream, *gw.region, key, nil, nil)
	if err != nil {
		glog.Fatalf("%v", err)
	}
	if *gw.prefix == "" && *gw.suffix == "" {
		glog.Fatalf("-bucket.prefix and -bucket.suffix are both empty: a local bucket would be an upstream bucket of the same global name")
	}
	// The IAM adds an admin identity for AWS_ACCESS_KEY_ID when it is set. Here
	// the identities are -config's and nothing else's.
	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("AWS_SECRET_ACCESS_KEY")
	iam := s3api.NewIdentityAccessManagementWithStore(&s3api.S3ApiServerOption{Config: *gw.config}, nil, string(credential.StoreTypeMemory))
	g, err := gateway.New(iam, up, gateway.Names{Prefix: *gw.prefix, Suffix: *gw.suffix})
	if err != nil {
		glog.Fatalf("%v", err)
	}

	addr := net.JoinHostPort(*gw.bind, strconv.Itoa(*gw.port))
	srv := &http.Server{
		Addr:              addr,
		Handler:           g,
		ReadHeaderTimeout: 30 * time.Second,
		IdleTimeout:       2 * time.Minute,
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	go func() {
		<-ctx.Done()
		// The endpoint is withdrawn from the service while this sleeps, so the
		// requests still routed here are served rather than refused.
		time.Sleep(*gw.drain)
		shutdown, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdown); err != nil {
			glog.Warningf("shutdown: %v", err)
		}
	}()
	glog.Infof("gateway on %s -> %s (%s), buckets %s<name>%s", addr, *gw.upstream, *gw.region, *gw.prefix, *gw.suffix)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		glog.Fatalf("listen %s: %v", addr, err)
	}
	return true
}

// readUpstreamKey reads access-key and secret-key from dir.
func readUpstreamKey(dir string) (aws.Credentials, error) {
	read := func(name string) (string, error) {
		b, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return "", err
		}
		v := strings.TrimSpace(string(b))
		if v == "" {
			return "", fmt.Errorf("%s is empty", filepath.Join(dir, name))
		}
		return v, nil
	}
	id, err := read("access-key")
	if err != nil {
		return aws.Credentials{}, err
	}
	secret, err := read("secret-key")
	if err != nil {
		return aws.Credentials{}, err
	}
	return aws.Credentials{AccessKeyID: id, SecretAccessKey: secret, Source: "file"}, nil
}
