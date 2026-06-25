package pluginworkers

import (
	"errors"
	"net"
	"strconv"
	"sync"
	"testing"

	"github.com/hanzoai/s3/s3/masterzap"
	"github.com/hanzoai/s3/s3/pb/master_pb"
	masterwire "github.com/hanzoai/s3/s3/wire/master"

	"github.com/zap-proto/go/transport"
)

// MasterServer provides a stub master that serves topology (VolumeList)
// responses over the native ZAP transport — the same transport the production
// master uses (gRPC is gone). It implements only VolumeList; every other
// masterwire.Backend method returns errFakeMasterUnsupported.
type MasterServer struct {
	t *testing.T

	server  *transport.Server
	address string // synthetic master address; ToMasterZapAddress() resolves to the listener

	mu       sync.RWMutex
	response *master_pb.VolumeListResponse
}

var errFakeMasterUnsupported = errors.New("fake master: RPC not supported (only VolumeList)")

// NewMasterServer starts a stub master that serves the provided response over
// ZAP. The returned Address() is an HTTP-style master address whose
// ServerAddress.ToMasterZapAddress() (grpcPort+10000, i.e. http+20000) resolves
// back to this listener, so code under test reaches it via pb.WithMasterClient.
func NewMasterServer(t *testing.T, response *master_pb.VolumeListResponse) *MasterServer {
	t.Helper()

	ms := &MasterServer{t: t, response: response}

	dispatch := masterwire.Dispatch(fakeMasterBackend{ms: ms})
	server, err := transport.ListenStream("tcp", "127.0.0.1:0", dispatch, nil)
	if err != nil {
		t.Fatalf("listen master ZAP: %v", err)
	}
	ms.server = server

	host, zapPort := splitHostPort(t, server.Addr().String())
	// ToMasterZapAddress(host:p) == host:(p+20000); invert so the reported
	// address derives to the actual ZAP listener.
	ms.address = net.JoinHostPort(host, strconv.Itoa(zapPort-20000))

	t.Cleanup(func() { ms.Shutdown() })

	return ms
}

// Address returns the synthetic master address (see NewMasterServer).
func (m *MasterServer) Address() string { return m.address }

// SetVolumeListResponse updates the response served by VolumeList.
func (m *MasterServer) SetVolumeListResponse(response *master_pb.VolumeListResponse) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.response = response
}

// Shutdown stops the master ZAP server.
func (m *MasterServer) Shutdown() {
	if m.server != nil {
		_ = m.server.Close()
	}
}

func (m *MasterServer) volumeList() *master_pb.VolumeListResponse {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.response == nil {
		return &master_pb.VolumeListResponse{}
	}
	return m.response
}

func splitHostPort(t *testing.T, addr string) (string, int) {
	t.Helper()
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		t.Fatalf("split master addr %q: %v", addr, err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatalf("parse master port %q: %v", portStr, err)
	}
	return host, port
}

// fakeMasterBackend implements masterwire.Backend, serving VolumeList from the
// stub and rejecting everything else.
type fakeMasterBackend struct{ ms *MasterServer }

func (b fakeMasterBackend) VolumeList(req []byte) ([]byte, error) {
	if _, err := masterzap.VolumeListReqFromWire(req); err != nil {
		return nil, err
	}
	return masterzap.VolumeListRespToWire(b.ms.volumeList()), nil
}

func (fakeMasterBackend) SendHeartbeat(req []byte) ([]byte, error) {
	return nil, errFakeMasterUnsupported
}
func (fakeMasterBackend) KeepConnected(req []byte) ([]byte, error) {
	return nil, errFakeMasterUnsupported
}
func (fakeMasterBackend) StreamAssign(req []byte) ([]byte, error) {
	return nil, errFakeMasterUnsupported
}
func (fakeMasterBackend) LookupVolume(req []byte) ([]byte, error) {
	return nil, errFakeMasterUnsupported
}
func (fakeMasterBackend) Assign(req []byte) ([]byte, error) { return nil, errFakeMasterUnsupported }
func (fakeMasterBackend) LookupEcVolume(req []byte) ([]byte, error) {
	return nil, errFakeMasterUnsupported
}
func (fakeMasterBackend) CollectionList(req []byte) ([]byte, error) {
	return nil, errFakeMasterUnsupported
}
func (fakeMasterBackend) CollectionDelete(req []byte) ([]byte, error) {
	return nil, errFakeMasterUnsupported
}
func (fakeMasterBackend) Statistics(req []byte) ([]byte, error) { return nil, errFakeMasterUnsupported }
func (fakeMasterBackend) VacuumVolume(req []byte) ([]byte, error) {
	return nil, errFakeMasterUnsupported
}
func (fakeMasterBackend) DisableVacuum(req []byte) ([]byte, error) {
	return nil, errFakeMasterUnsupported
}
func (fakeMasterBackend) EnableVacuum(req []byte) ([]byte, error) {
	return nil, errFakeMasterUnsupported
}
func (fakeMasterBackend) VolumeMarkReadonly(req []byte) ([]byte, error) {
	return nil, errFakeMasterUnsupported
}
func (fakeMasterBackend) GetMasterConfiguration(req []byte) ([]byte, error) {
	return nil, errFakeMasterUnsupported
}
func (fakeMasterBackend) ListClusterNodes(req []byte) ([]byte, error) {
	return nil, errFakeMasterUnsupported
}
func (fakeMasterBackend) LeaseAdminToken(req []byte) ([]byte, error) {
	return nil, errFakeMasterUnsupported
}
func (fakeMasterBackend) ReleaseAdminToken(req []byte) ([]byte, error) {
	return nil, errFakeMasterUnsupported
}
func (fakeMasterBackend) Ping(req []byte) ([]byte, error) { return nil, errFakeMasterUnsupported }
func (fakeMasterBackend) RaftListClusterServers(req []byte) ([]byte, error) {
	return nil, errFakeMasterUnsupported
}
func (fakeMasterBackend) RaftAddServer(req []byte) ([]byte, error) {
	return nil, errFakeMasterUnsupported
}
func (fakeMasterBackend) RaftRemoveServer(req []byte) ([]byte, error) {
	return nil, errFakeMasterUnsupported
}
func (fakeMasterBackend) RaftLeadershipTransfer(req []byte) ([]byte, error) {
	return nil, errFakeMasterUnsupported
}
func (fakeMasterBackend) VolumeGrow(req []byte) ([]byte, error) { return nil, errFakeMasterUnsupported }
