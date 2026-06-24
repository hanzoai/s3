package command

import (
	"testing"

	"github.com/hanzoai/s3/s3/worker/tasks/vacuum"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestMiniDefaultPluginJobTypes(t *testing.T) {
	dialOption := grpc.WithTransportCredentials(insecure.NewCredentials())
	// defaultMiniPluginJobTypes is "all", which includes every registered handler
	handlers, err := buildPluginWorkerHandlers(defaultMiniPluginJobTypes, dialOption, int(vacuum.DefaultMaxExecutionConcurrency), "")
	if err != nil {
		t.Fatalf("buildPluginWorkerHandlers(mini default) err = %v", err)
	}
	// 6 registered handlers since the Iceberg maintenance backend was ripped.
	if len(handlers) != 6 {
		t.Fatalf("expected mini default job types to include 6 handlers, got %d", len(handlers))
	}
}
