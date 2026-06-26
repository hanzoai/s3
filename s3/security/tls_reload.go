package security

import (
	"crypto/tls"

	"github.com/hanzoai/s3/s3/security/certreload"
)

// NewReloadingServerCertificate returns a GetCertificate callback for
// tls.Config.GetCertificate, backed by a refreshing keypair provider so
// rotated certs (e.g. from Kubernetes cert-manager) are picked up
// without a restart. Close() the returned provider at shutdown.
//
// Thin wrapper over certreload; the shared (gRPC-free) implementation lives in
// the subpackage so both server TLS (this package) and HTTP client TLS
// (s3/util/http/client) can use it without an import cycle.
func NewReloadingServerCertificate(certFile, keyFile string) (func(*tls.ClientHelloInfo) (*tls.Certificate, error), *certreload.Provider, error) {
	return certreload.NewServerGetCertificate(certFile, keyFile)
}
