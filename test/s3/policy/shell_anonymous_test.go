package policy

import (
	"fmt"
	"testing"

	"github.com/hanzoai/s3/s3/pb"
	"github.com/stretchr/testify/require"
)

// TestShellAnonymousAccess exercises s3.anonymous.* commands end-to-end.
func TestShellAnonymousAccess(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	cluster, err := startMiniCluster(t)
	require.NoError(t, err)
	defer cluster.Stop()

	const s3Cmd = "s3"
	master := string(pb.NewServerAddress("127.0.0.1", cluster.masterPort, cluster.masterGrpcPort))
	filer := string(pb.NewServerAddress("127.0.0.1", cluster.filerPort, cluster.filerGrpcPort))

	bucketName := uniqueName("anon-bkt")
	execShell(t, s3Cmd, master, filer, fmt.Sprintf("s3.bucket.create -name %s", bucketName))
	defer execShell(t, s3Cmd, master, filer, fmt.Sprintf("s3.bucket.delete -name %s", bucketName))

	t.Run("SetAndGet", func(t *testing.T) {
		out := execShell(t, s3Cmd, master, filer,
			fmt.Sprintf("s3.anonymous.set -bucket %s -access Read,List", bucketName))
		requireContains(t, out, bucketName, "anonymous.set output")

		out = execShell(t, s3Cmd, master, filer, fmt.Sprintf("s3.anonymous.get -bucket %s", bucketName))
		requireContains(t, out, bucketName, "anonymous.get bucket")
		requireContains(t, out, "Read", "anonymous.get read action")
		requireContains(t, out, "List", "anonymous.get list action")
	})

	t.Run("List", func(t *testing.T) {
		out := execShell(t, s3Cmd, master, filer, "s3.anonymous.list")
		requireContains(t, out, bucketName, "anonymous.list contains bucket")
	})

	t.Run("SetNone", func(t *testing.T) {
		execShell(t, s3Cmd, master, filer,
			fmt.Sprintf("s3.anonymous.set -bucket %s -access none", bucketName))

		out := execShell(t, s3Cmd, master, filer, fmt.Sprintf("s3.anonymous.get -bucket %s", bucketName))
		requireContains(t, out, "none", "anonymous.get after set none")

		out = execShell(t, s3Cmd, master, filer, "s3.anonymous.list")
		requireNotContains(t, out, bucketName, "anonymous.list after clearing")
	})
}
