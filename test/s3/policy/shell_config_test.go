package policy

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hanzoai/s3/s3/pb"
	"github.com/stretchr/testify/require"
)

// TestShellConfigShow verifies s3.config.show outputs a summary of IAM config.
func TestShellConfigShow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	cluster, err := startMiniCluster(t)
	require.NoError(t, err)
	defer cluster.Stop()

	const s3Cmd = "s3"
	master := string(pb.NewServerAddress("127.0.0.1", cluster.masterPort, cluster.masterGrpcPort))
	filer := string(pb.NewServerAddress("127.0.0.1", cluster.filerPort, cluster.filerGrpcPort))

	userName := uniqueName("cfg-user")
	groupName := uniqueName("cfg-grp")
	execShell(t, s3Cmd, master, filer, fmt.Sprintf("s3.user.create -name %s", userName))
	defer execShell(t, s3Cmd, master, filer, fmt.Sprintf("s3.user.delete -name %s", userName))
	execShell(t, s3Cmd, master, filer, fmt.Sprintf("s3.group.create -name %s", groupName))
	defer execShell(t, s3Cmd, master, filer, fmt.Sprintf("s3.group.delete -name %s", groupName))

	out := execShell(t, s3Cmd, master, filer, "s3.config.show")
	requireContains(t, out, "S3 IAM Configuration Summary", "config.show header")
	requireContains(t, out, userName, "config.show contains user")
	requireContains(t, out, groupName, "config.show contains group")
}

// TestShellIAMExportImport does a roundtrip: create resources, export, delete, import, verify.
func TestShellIAMExportImport(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	cluster, err := startMiniCluster(t)
	require.NoError(t, err)
	defer cluster.Stop()

	const s3Cmd = "s3"
	master := string(pb.NewServerAddress("127.0.0.1", cluster.masterPort, cluster.masterGrpcPort))
	filer := string(pb.NewServerAddress("127.0.0.1", cluster.filerPort, cluster.filerGrpcPort))

	userName := uniqueName("exp-user")
	groupName := uniqueName("exp-grp")
	execShell(t, s3Cmd, master, filer, fmt.Sprintf("s3.user.create -name %s", userName))
	execShell(t, s3Cmd, master, filer, fmt.Sprintf("s3.group.create -name %s", groupName))

	exportFile := filepath.Join(t.TempDir(), "iam_export.txt")
	execShell(t, s3Cmd, master, filer, fmt.Sprintf("s3.iam.export -file %s", exportFile))

	data, err := os.ReadFile(exportFile)
	require.NoError(t, err)
	content := string(data)
	requireContains(t, content, userName, "export file contains user")
	requireContains(t, content, groupName, "export file contains group")

	// Delete the resources.
	execShell(t, s3Cmd, master, filer, fmt.Sprintf("s3.user.delete -name %s", userName))
	execShell(t, s3Cmd, master, filer, fmt.Sprintf("s3.group.delete -name %s", groupName))

	out := execShell(t, s3Cmd, master, filer, "s3.user.list")
	requireNotContains(t, out, userName, "user gone before import")

	// Import.
	out = execShell(t, s3Cmd, master, filer, fmt.Sprintf("s3.iam.import -file %s -apply", exportFile))
	requireContains(t, out, "Imported IAM configuration", "import output")

	// Verify resources restored.
	out = execShell(t, s3Cmd, master, filer, "s3.user.list")
	requireContains(t, out, userName, "user restored after import")
	out = execShell(t, s3Cmd, master, filer, "s3.group.list")
	requireContains(t, out, groupName, "group restored after import")

	// Cleanup.
	execShell(t, s3Cmd, master, filer, fmt.Sprintf("s3.user.delete -name %s", userName))
	execShell(t, s3Cmd, master, filer, fmt.Sprintf("s3.group.delete -name %s", groupName))
}
