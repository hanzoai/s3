package command

import (
	"sort"
	"testing"

	"github.com/hanzoai/s3/s3/pb"
	"github.com/hanzoai/s3/s3/worker/tasks/vacuum"
)

func TestWorkerDefaultJobTypes(t *testing.T) {
	dialOption := pb.DialOption{}
	handlers, err := buildPluginWorkerHandlers(*workerJobType, dialOption, int(vacuum.DefaultMaxExecutionConcurrency), "")
	if err != nil {
		t.Fatalf("buildPluginWorkerHandlers(default worker flag) err = %v", err)
	}
	// iceberg_maintenance is no longer a registered worker handler (the Iceberg
	// maintenance backend was ripped); the registry now resolves these 6.
	want := []string{
		"admin_script",
		"ec_balance",
		"erasure_coding",
		"s3_lifecycle",
		"vacuum",
		"volume_balance",
	}
	got := make([]string, 0, len(handlers))
	for _, h := range handlers {
		got = append(got, h.Capability().JobType)
	}
	sort.Strings(got)
	if len(got) != len(want) {
		t.Fatalf("default worker job types: got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("default worker job types: got %v, want %v", got, want)
		}
	}
}
