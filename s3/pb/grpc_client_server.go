package pb

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/hanzoai/s3/s3/glog"
	"github.com/hanzoai/s3/s3/pb/volume_server_pb"
	"github.com/hanzoai/s3/s3/svc/mq"
	"github.com/hanzoai/s3/s3/util"

	"github.com/hanzoai/s3/s3/pb/filer_pb"
	"github.com/hanzoai/s3/s3/pb/master_pb"
	"github.com/hanzoai/s3/s3/pb/mq_pb"
	"github.com/hanzoai/s3/s3/svc/master"
	"github.com/hanzoai/s3/s3/svc/volume"

	"github.com/zap-proto/go/transport"
)

var (
	// localGrpcSockets maps gRPC port numbers to Unix socket paths for
	// services running in this process. Used by both the server side
	// (to start serving on the socket) and the client side (to dial it
	// instead of TCP for same-host RPCs).
	localGrpcSockets = make(map[int]string)
	// localGrpcHosts maps gRPC port → set of host strings that count as
	// "this machine" for that port. Only dials whose target host is in this
	// set are routed through the Unix socket; dials to other hosts that
	// happen to use the same port number go over TCP as normal.
	//
	// Without this host check the hijack was keyed purely on port, so a
	// remote server reusing one of our local gRPC ports (e.g. a standalone
	// `s3 volume -port=7334` defaulting `port.grpc=17334` against a
	// `s3 server` whose in-process volume socket is also on 17334) had
	// its inbound RPCs silently rerouted to our local Unix socket — see
	// issue #9254.
	localGrpcHosts       = make(map[int]map[string]struct{})
	localGrpcSocketsLock sync.RWMutex
)

func init() {
	t := http.DefaultTransport.(*http.Transport)
	t.MaxIdleConnsPerHost = 1024
	t.MaxIdleConns = 1024
	// Bind outbound HTTP to the -ip.bind source address. Reads the setting
	// per dial, so it applies regardless of when SetOutboundLocalIP runs
	// relative to this init.
	t.DialContext = util.OutboundDialContext
}

// RegisterLocalGrpcSocket registers a Unix socket path for a service running on
// host:grpcPort in this process. After this is called, any same-host dial to
// host:grpcPort (or to a loopback alias of host on the same port) is routed
// through the Unix socket. Dials to any other host on the same port still go
// over TCP.
//
// No-op on Windows: the /tmp/...sock paths callers pass are POSIX-only and
// the listen/dial would fail at runtime (#9430). Skipping registration leaves
// the maps empty, so resolveLocalGrpcSocket short-circuits and same-host RPCs
// go over TCP.
func RegisterLocalGrpcSocket(host string, grpcPort int, socketPath string) {
	if runtime.GOOS == "windows" {
		return
	}
	localGrpcSocketsLock.Lock()
	defer localGrpcSocketsLock.Unlock()
	localGrpcSockets[grpcPort] = socketPath
	hosts, ok := localGrpcHosts[grpcPort]
	if !ok {
		hosts = make(map[string]struct{})
		localGrpcHosts[grpcPort] = hosts
	}
	// The advertised host plus the well-known loopback aliases all refer
	// to "this machine"; the empty entry covers SplitHostPort outputs like
	// ":17334".
	for _, h := range []string{host, "", "localhost", "127.0.0.1", "::1"} {
		hosts[h] = struct{}{}
	}
}

// GetLocalGrpcSocket returns the Unix socket path for a gRPC port, or empty if not registered.
func GetLocalGrpcSocket(grpcPort int) string {
	localGrpcSocketsLock.RLock()
	defer localGrpcSocketsLock.RUnlock()
	return localGrpcSockets[grpcPort]
}

// resolveLocalGrpcSocket returns the Unix socket path to use for an outgoing
// dial, or empty if the dial should go over TCP. It matches both the host and
// the port: a remote peer that happens to reuse a local service's grpc port
// must not be rerouted to the local socket (issue #9254).
func resolveLocalGrpcSocket(address string) string {
	host, portStr, err := net.SplitHostPort(address)
	if err != nil {
		return ""
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return ""
	}
	localGrpcSocketsLock.RLock()
	defer localGrpcSocketsLock.RUnlock()
	if _, ok := localGrpcHosts[port][host]; !ok {
		return ""
	}
	return localGrpcSockets[port]
}

func hostAndPort(address string) (host string, port uint64, err error) {
	colonIndex := strings.LastIndex(address, ":")
	if colonIndex < 0 {
		return "", 0, fmt.Errorf("server should have hostname:port format: %v", address)
	}
	dotIndex := strings.LastIndex(address, ".")
	if dotIndex > colonIndex {
		// port format is "port.grpcPort"
		port, err = strconv.ParseUint(address[colonIndex+1:dotIndex], 10, 64)
		if err != nil {
			return "", 0, fmt.Errorf("server port parse error: %w", err)
		}
		return address[:colonIndex], port, err
	}
	port, err = strconv.ParseUint(address[colonIndex+1:], 10, 64)
	if err != nil {
		return "", 0, fmt.Errorf("server port parse error: %w", err)
	}

	return address[:colonIndex], port, err
}

func ServerToGrpcAddress(server string) (serverGrpcAddress string) {

	colonIndex := strings.LastIndex(server, ":")
	if colonIndex >= 0 {
		if dotIndex := strings.LastIndex(server, "."); dotIndex > colonIndex {
			// port format is "port.grpcPort"
			// return the host:grpcPort
			host := server[:colonIndex]
			grpcPort := server[dotIndex+1:]
			if _, err := strconv.ParseUint(grpcPort, 10, 64); err == nil {
				return util.JoinHostPort(host, int(0+util.ParseInt(grpcPort, 0)))
			}
		}
	}

	host, port, parseErr := hostAndPort(server)
	if parseErr != nil {
		glog.Fatalf("server address %s parse error: %v", server, parseErr)
	}

	grpcPort := int(port) + 10000

	return util.JoinHostPort(host, grpcPort)
}

func GrpcAddressToServerAddress(grpcAddress string) (serverAddress string) {
	host, grpcPort, parseErr := hostAndPort(grpcAddress)
	if parseErr != nil {
		glog.Fatalf("server grpc address %s parse error: %v", grpcAddress, parseErr)
	}

	port := int(grpcPort) - 10000

	return util.JoinHostPort(host, port)
}

// DialMasterZapAddr opens a ZAP connection to a master's native ZAP transport
// (the grpcPort+10000 offset, ToMasterZapAddress). Used by masterPool and by the
// volume-server heartbeat, which owns a dedicated long-lived conn rather than a
// pooled one. When grpc.master.cert/.key is
// configured it is PQ-secured TLS (transport.PQTLSConfig pins X25519MLKEM768, the
// PQ X-Wing curve) presenting the client cert and trusting grpc.ca — the same mTLS
// the legacy gRPC master client used. Otherwise plaintext (loopback / dev). The
// returned *Conn drives both unary Call and OpenStream.
func DialMasterZapAddr(addr string) (transport.Conn, error) {
	if cfg := ClientTLSConfig(util.GetViper(), "grpc.master"); cfg != nil {
		return transport.DialTLS("tcp", addr, transport.PQTLSConfig(cfg))
	}
	return transport.Dial("tcp", addr)
}

// masterPool reuses one ZAP connection per master address across calls — a Conn
// is concurrency-safe, so this avoids a fresh TCP (and, under grpc.master mTLS, a
// fresh X25519MLKEM768 handshake) on every master RPC. The master analogue of
// filerPool.
var masterPool = transport.NewPool(DialMasterZapAddr)

// WithMasterClient runs fn with a master_pb.HanzoClient backed by a pooled ZAP
// connection to master's native transport. The streamingMode/grpcDialOption/
// waitForReady parameters are retained for caller compatibility; the ZAP path
// needs neither a streaming flag (every stream is a transport stream) nor a dial
// option, and the pool already (re)dials lazily. ctx is honored by the per-RPC
// calls fn issues. This is the master analogue of WithGrpcFilerClient.
func WithMasterClient(ctx context.Context, streamingMode bool, masterAddr ServerAddress, grpcDialOption DialOption, waitForReady bool, fn func(client master_pb.HanzoClient) error) error {
	_, _, _, _ = ctx, streamingMode, grpcDialOption, waitForReady
	return masterPool.With(masterAddr.ToMasterZapAddress(), func(conn transport.Conn) error {
		return fn(master.New(conn, nil))
	})
}

// dialVolumeZapAddr opens a ZAP connection to a volume-server grpc address. When
// grpc.volume.cert/.key is configured it is PQ-secured TLS (transport.PQTLSConfig
// pins X25519MLKEM768, the PQ X-Wing curve) presenting the client cert and
// trusting grpc.ca — the same mTLS the legacy gRPC volume client used. Otherwise
// plaintext (loopback / dev), matching the volume server's gating in
// command/volume.go. The returned *Conn drives both unary Call and OpenStream.
func dialVolumeZapAddr(addr string) (transport.Conn, error) {
	if cfg := ClientTLSConfig(util.GetViper(), "grpc.volume"); cfg != nil {
		return transport.DialTLS("tcp", addr, transport.PQTLSConfig(cfg))
	}
	return transport.Dial("tcp", addr)
}

// volumePool reuses one ZAP connection per volume-server address across calls — a
// Conn is concurrency-safe, so this avoids a fresh TCP (and, under grpc.volume
// mTLS, a fresh X25519MLKEM768 handshake) on every volume RPC. Generic pooling
// lives in the transport (transport.Pool); only the dial choice is ours.
var volumePool = transport.NewPool(dialVolumeZapAddr)

// WithVolumeServerClient dials the volume server over the native ZAP transport and
// runs fn with a volume_server_pb.VolumeServerClient backed by that connection
// (volume.New). The streamingMode/grpcDialOption parameters are retained for
// caller compatibility; the ZAP path needs neither a streaming flag (every stream
// is a transport stream) nor a dial option.
func WithVolumeServerClient(streamingMode bool, volumeServer ServerAddress, grpcDialOption DialOption, fn func(client volume_server_pb.VolumeServerClient) error) error {
	_, _ = streamingMode, grpcDialOption
	return volumePool.With(volumeServer.ToGrpcAddress(), func(conn transport.Conn) error {
		return fn(volume.New(conn, nil))
	})
}

// WithOneOfGrpcMasterClients tries each master address in turn over the ZAP
// transport, returning on the first that runs fn without error.
func WithOneOfGrpcMasterClients(streamingMode bool, masterGrpcAddresses map[string]ServerAddress, grpcDialOption DialOption, fn func(client master_pb.HanzoClient) error) (err error) {
	_, _ = streamingMode, grpcDialOption
	for _, masterAddress := range masterGrpcAddresses {
		err = masterPool.With(masterAddress.ToMasterZapAddress(), func(conn transport.Conn) error {
			return fn(master.New(conn, nil))
		})
		if err == nil {
			return nil
		}
	}

	return err
}

// DialBrokerZapAddr opens a ZAP connection to a broker address. When
// grpc.msg_broker.cert/.key is configured it is PQ-secured TLS (transport.PQTLSConfig
// pins X25519MLKEM768, the PQ X-Wing curve) presenting the client cert and trusting
// grpc.ca — the same mTLS the legacy gRPC broker client used. Otherwise plaintext
// (loopback / dev), matching the broker server's gating in command/mq_broker.go. The
// returned *Conn drives both unary Call and OpenStream. The broker fully cut over to
// ZAP, so its address IS the ZAP endpoint (no port offset, unlike master/IAM).
func DialBrokerZapAddr(addr string) (transport.Conn, error) {
	if cfg := ClientTLSConfig(util.GetViper(), "grpc.msg_broker"); cfg != nil {
		return transport.DialTLS("tcp", addr, transport.PQTLSConfig(cfg))
	}
	return transport.Dial("tcp", addr)
}

// brokerPool reuses one ZAP connection per broker address across calls — a Conn is
// concurrency-safe, so this avoids a fresh TCP (and, under grpc.msg_broker mTLS, a
// fresh X25519MLKEM768 handshake) on every broker RPC. Generic pooling lives in the
// transport (transport.Pool); only the dial choice is ours. Mirrors filerPool.
var brokerPool = transport.NewPool(DialBrokerZapAddr)

// WithBrokerGrpcClient dials the broker over the native ZAP transport and runs fn
// with a mq_pb.HanzoMessagingClient backed by that connection (mq.New). The
// streamingMode/grpcDialOption parameters are retained for caller compatibility; the
// ZAP path needs neither a streaming flag (every stream is a transport stream) nor a
// dial option. The pooled connection outlives fn.
func WithBrokerGrpcClient(streamingMode bool, brokerGrpcAddress string, grpcDialOption DialOption, fn func(client mq_pb.HanzoMessagingClient) error) error {
	_, _ = streamingMode, grpcDialOption
	return brokerPool.With(brokerGrpcAddress, func(conn transport.Conn) error {
		return fn(mq.New(conn, nil))
	})
}

func WithFilerClient(streamingMode bool, signature int32, filer ServerAddress, grpcDialOption DialOption, fn func(client filer_pb.HanzoFilerClient) error) error {

	return WithGrpcFilerClient(streamingMode, signature, filer, grpcDialOption, fn)

}

// dialFilerZapAddr opens a ZAP connection to a filer grpc address. When
// grpc.filer.cert/.key is configured it is PQ-secured TLS (transport.PQTLSConfig
// pins X25519MLKEM768, the PQ X-Wing curve) presenting the client cert and
// trusting grpc.ca — the same mTLS the legacy gRPC filer client used. Otherwise
// plaintext (loopback / dev), matching the filer server's gating in
// command/filer.go. The returned *Conn drives both unary Call and OpenStream.
func dialFilerZapAddr(addr string) (transport.Conn, error) {
	if cfg := ClientTLSConfig(util.GetViper(), "grpc.filer"); cfg != nil {
		return transport.DialTLS("tcp", addr, transport.PQTLSConfig(cfg))
	}
	return transport.Dial("tcp", addr)
}

// filerPool reuses one ZAP connection per filer address across calls — a Conn is
// concurrency-safe, so this avoids a fresh TCP (and, under grpc.filer mTLS, a
// fresh X25519MLKEM768 handshake) on every filer RPC. Generic pooling lives in
// the transport (transport.Pool); only the dial choice is ours.
var filerPool = transport.NewPool(dialFilerZapAddr)

// WithGrpcFilerClient dials the filer over the native ZAP transport and runs fn
// with a filer_pb.HanzoFilerClient backed by that connection (NewZapFilerClient).
// The streamingMode/signature/grpcDialOption parameters are retained for caller
// compatibility; the ZAP path needs neither a streaming flag (every stream is a
// transport stream) nor a dial option. The connection is closed when fn returns.
func WithGrpcFilerClient(streamingMode bool, signature int32, filerAddress ServerAddress, grpcDialOption DialOption, fn func(client filer_pb.HanzoFilerClient) error) error {
	_, _, _ = streamingMode, signature, grpcDialOption
	return filerPool.With(filerAddress.ToFilerZapAddress(), func(conn transport.Conn) error {
		return fn(NewZapFilerClient(conn))
	})
}

// WithOneOfGrpcFilerClients tries each filer address in turn over the ZAP
// transport, returning on the first that runs fn without error.
func WithOneOfGrpcFilerClients(streamingMode bool, filerAddresses []ServerAddress, grpcDialOption DialOption, fn func(client filer_pb.HanzoFilerClient) error) (err error) {
	_, _ = streamingMode, grpcDialOption
	for _, filerAddress := range filerAddresses {
		err = filerPool.With(filerAddress.ToFilerZapAddress(), func(conn transport.Conn) error {
			return fn(NewZapFilerClient(conn))
		})
		if err == nil {
			return nil
		}
	}

	return err
}
