package s3server

import (
	"context"
	"testing"
	"time"

	"github.com/hanzoai/s3/s3/cluster/lock_manager"
	"github.com/hanzoai/s3/s3/filer"
	"github.com/hanzoai/s3/s3/pb"
	"github.com/hanzoai/s3/s3/pb/filer_pb"
	"github.com/stretchr/testify/require"
)

func TestFindLockOwnerExpiredLockReturnsNotFound(t *testing.T) {
	fs := &FilerServer{
		option: &FilerOption{Host: "filer1:8888"},
		filer: &filer.Filer{
			Dlm: lock_manager.NewDistributedLockManager("filer1:8888"),
		},
	}
	fs.filer.Dlm.LockRing.SetSnapshot([]pb.ServerAddress{"filer1:8888"}, 0)
	fs.filer.Dlm.InsertLock("expired-lock", time.Now().Add(-time.Second).UnixNano(), "token1", "owner1", 5, 2)

	resp, err := fs.FindLockOwner(context.Background(), &filer_pb.FindLockOwnerRequest{
		Name: "expired-lock",
	})
	require.Nil(t, resp)
	require.Error(t, err)
	// The filer speaks ZAP now: the server tags the error with its code name
	// ("NotFound: ...") rather than a gRPC status, and clients classify on that
	// (see mount/error_classifier.go).
	require.Contains(t, err.Error(), "NotFound")
}
