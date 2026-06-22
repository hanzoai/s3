package policy

import (
	"fmt"
	"testing"

	"github.com/hanzoai/s3/s3/pb"
	"github.com/stretchr/testify/require"
)

// TestShellGroupLifecycle exercises s3.group.* commands end-to-end.
func TestShellGroupLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	cluster, err := startMiniCluster(t)
	require.NoError(t, err)
	defer cluster.Stop()

	const s3Cmd = "s3"
	master := string(pb.NewServerAddress("127.0.0.1", cluster.masterPort, cluster.masterGrpcPort))
	filer := string(pb.NewServerAddress("127.0.0.1", cluster.filerPort, cluster.filerGrpcPort))

	groupName := uniqueName("grp")
	userName := uniqueName("grpuser")

	execShell(t, s3Cmd, master, filer, fmt.Sprintf("s3.user.create -name %s", userName))
	defer execShell(t, s3Cmd, master, filer, fmt.Sprintf("s3.user.delete -name %s", userName))

	t.Run("CreateShowList", func(t *testing.T) {
		out := execShell(t, s3Cmd, master, filer, fmt.Sprintf("s3.group.create -name %s", groupName))
		requireContains(t, out, groupName, "group.create output")

		out = execShell(t, s3Cmd, master, filer, fmt.Sprintf("s3.group.show -name %s", groupName))
		requireContains(t, out, groupName, "group.show output")

		out = execShell(t, s3Cmd, master, filer, "s3.group.list")
		requireContains(t, out, groupName, "group.list")
	})

	t.Run("AddRemoveUser", func(t *testing.T) {
		out := execShell(t, s3Cmd, master, filer,
			fmt.Sprintf("s3.group.add-user -group %s -user %s", groupName, userName))
		requireContains(t, out, userName, "group.add-user output")

		out = execShell(t, s3Cmd, master, filer, fmt.Sprintf("s3.group.show -name %s", groupName))
		requireContains(t, out, userName, "group.show after add")

		out = execShell(t, s3Cmd, master, filer,
			fmt.Sprintf("s3.group.remove-user -group %s -user %s", groupName, userName))
		requireContains(t, out, userName, "group.remove-user output")

		out = execShell(t, s3Cmd, master, filer, fmt.Sprintf("s3.group.show -name %s", groupName))
		requireNotContains(t, out, fmt.Sprintf("\"%s\"", userName), "group.show after remove")
	})

	t.Run("Delete", func(t *testing.T) {
		execShell(t, s3Cmd, master, filer, fmt.Sprintf("s3.group.delete -name %s", groupName))
		out := execShell(t, s3Cmd, master, filer, "s3.group.list")
		requireNotContains(t, out, fmt.Sprintf("\"%s\"", groupName), "group.list after delete")
	})
}
