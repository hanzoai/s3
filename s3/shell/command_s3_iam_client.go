package shell

import (
	iamwire "github.com/hanzoai/s3/s3/wire/iam"
	"github.com/zap-proto/go/transport"
)

// withIamClient invokes fn against the filer's IAM service over the ZAP
// transport. The IAM service listens on its own ZAP endpoint at
// FilerAddress.ToIamZapAddress() (grpcPort+10000), separate from the shared
// gRPC port. Connection-level security is provided by the ZAP transport, so the
// shell no longer mints an admin Bearer token per call.
func (ce *CommandEnv) withIamClient(fn func(client *iamwire.HanzoIdentityAccessManagementClient) error) error {
	conn, err := transport.Dial("tcp", ce.option.FilerAddress.ToIamZapAddress())
	if err != nil {
		return err
	}
	defer conn.Close()
	return fn(iamwire.NewHanzoIdentityAccessManagementClient(conn, nil))
}
