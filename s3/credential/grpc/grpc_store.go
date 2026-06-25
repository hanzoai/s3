package grpc

import (
	"fmt"
	"sync"

	"github.com/zap-proto/go/transport"

	"github.com/hanzoai/s3/s3/credential"
	"github.com/hanzoai/s3/s3/pb"
	"github.com/hanzoai/s3/s3/security"
	"github.com/hanzoai/s3/s3/util"
	iamwire "github.com/hanzoai/s3/s3/wire/iam"
)

func init() {
	credential.Stores = append(credential.Stores, &IamGrpcStore{})
}

// IamGrpcStore implements CredentialStore by calling the filer's IAM gRPC
// service. The filer requires an admin-signed Bearer token on every RPC
// (see s3/server/filer_server_handlers_iam_grpc.go); SetAdminSigning must
// be called with the same jwt.filer_signing.key value that the filer reads
// from security.toml, or every call will fail with Unauthenticated.
type IamGrpcStore struct {
	filerAddressFunc func() pb.ServerAddress // Function to get current active filer
	grpcDialOption   pb.DialOption
	// adminSigningKey is the HS256 secret used to mint Bearer tokens that the
	// filer's IAM gRPC service validates. Must match jwt.filer_signing.key on
	// the filer side. Empty means no token is sent (the filer will reject).
	adminSigningKey             security.SigningKey
	adminSigningExpiresAfterSec int
	mu                          sync.RWMutex // Protects filerAddressFunc, grpcDialOption, adminSigningKey, and adminSigningExpiresAfterSec
}

func (store *IamGrpcStore) GetName() credential.CredentialStoreTypeName {
	return credential.StoreTypeGrpc
}

func (store *IamGrpcStore) Initialize(configuration util.Configuration, prefix string) error {
	if configuration != nil {
		filerAddr := configuration.GetString(prefix + "filer")
		if filerAddr != "" {
			store.mu.Lock()
			store.filerAddressFunc = func() pb.ServerAddress {
				return pb.ServerAddress(filerAddr)
			}
			store.mu.Unlock()
		}
	}
	return nil
}

func (store *IamGrpcStore) SetFilerAddressFunc(getFiler func() pb.ServerAddress, grpcDialOption pb.DialOption) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.filerAddressFunc = getFiler
	store.grpcDialOption = grpcDialOption
}

// SetAdminSigning configures the HS256 secret used to mint Bearer tokens for
// the filer's IAM gRPC service. The key must match jwt.filer_signing.key in
// the filer's security.toml. If expiresAfterSec is 0, tokens are minted
// without an exp claim.
func (store *IamGrpcStore) SetAdminSigning(key security.SigningKey, expiresAfterSec int) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.adminSigningKey = key
	store.adminSigningExpiresAfterSec = expiresAfterSec
}

// withIamClient invokes fn against a client to the filer's IAM service over the
// ZAP transport. The service listens on its own ZAP endpoint at
// FilerAddress.ToIamZapAddress() (grpcPort+10000); connection-level security is
// provided by the transport, so the admin signing key is no longer consulted.
func (store *IamGrpcStore) withIamClient(fn func(client *iamwire.HanzoIdentityAccessManagementClient) error) error {
	store.mu.RLock()
	if store.filerAddressFunc == nil {
		store.mu.RUnlock()
		return fmt.Errorf("iam_grpc: filer not yet available")
	}
	filerAddress := store.filerAddressFunc()
	store.mu.RUnlock()

	if filerAddress == "" {
		return fmt.Errorf("iam_grpc: no filer discovered yet")
	}

	conn, err := transport.Dial("tcp", filerAddress.ToIamZapAddress())
	if err != nil {
		return err
	}
	defer conn.Close()
	return fn(iamwire.NewHanzoIdentityAccessManagementClient(conn, nil))
}

func (store *IamGrpcStore) Shutdown() {
}
