package pluginzapbridge

import (
	"testing"
	"time"

	"github.com/hanzoai/s3/s3/pb/plugin_pb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// roundTripWorkerToAdmin marshals m to ZAP and back, asserting proto-equality.
func roundTripWorkerToAdmin(t *testing.T, m *plugin_pb.WorkerToAdminMessage) {
	t.Helper()
	got, err := UnmarshalWorkerToAdmin(MarshalWorkerToAdmin(m))
	if err != nil {
		t.Fatalf("UnmarshalWorkerToAdmin: %v", err)
	}
	if !proto.Equal(m, got) {
		t.Fatalf("WorkerToAdmin round-trip mismatch\n want: %v\n  got: %v", m, got)
	}
}

func roundTripAdminToWorker(t *testing.T, m *plugin_pb.AdminToWorkerMessage) {
	t.Helper()
	got, err := UnmarshalAdminToWorker(MarshalAdminToWorker(m))
	if err != nil {
		t.Fatalf("UnmarshalAdminToWorker: %v", err)
	}
	if !proto.Equal(m, got) {
		t.Fatalf("AdminToWorker round-trip mismatch\n want: %v\n  got: %v", m, got)
	}
}

func TestWorkerToAdmin_Hello(t *testing.T) {
	ts := timestamppb.New(time.Unix(0, 1_700_000_000_123_456_789))
	roundTripWorkerToAdmin(t, &plugin_pb.WorkerToAdminMessage{
		WorkerId: "w-1",
		SentAt:   ts,
		Body: &plugin_pb.WorkerToAdminMessage_Hello{Hello: &plugin_pb.WorkerHello{
			WorkerId:         "w-1",
			WorkerInstanceId: "inst-9",
			Address:          "10.0.0.5:8080",
			WorkerVersion:    "v2",
			ProtocolVersion:  "1",
			Capabilities: []*plugin_pb.JobTypeCapability{
				{JobType: "vacuum", CanDetect: true, CanExecute: true, MaxDetectionConcurrency: 2, MaxExecutionConcurrency: 4, DisplayName: "Vacuum", Description: "reclaim", Weight: 7},
				{JobType: "balance", CanExecute: true, Weight: 3},
			},
			Metadata: map[string]string{"zone": "a", "rack": "r1"},
		}},
	})
}

func TestWorkerToAdmin_Heartbeat(t *testing.T) {
	roundTripWorkerToAdmin(t, &plugin_pb.WorkerToAdminMessage{
		WorkerId: "w-2",
		Body: &plugin_pb.WorkerToAdminMessage_Heartbeat{Heartbeat: &plugin_pb.WorkerHeartbeat{
			WorkerId: "w-2",
			RunningWork: []*plugin_pb.RunningWork{
				{WorkId: "wk-1", Kind: plugin_pb.WorkKind_WORK_KIND_EXECUTION, JobType: "vacuum", State: plugin_pb.JobState_JOB_STATE_RUNNING, ProgressPercent: 42.5, Stage: "scan"},
			},
			DetectionSlotsUsed:  1,
			DetectionSlotsTotal: 2,
			ExecutionSlotsUsed:  3,
			ExecutionSlotsTotal: 4,
			QueuedJobsByType:    map[string]int32{"vacuum": 5, "balance": 2},
			Metadata:            map[string]string{"k": "v"},
		}},
	})
}

func TestWorkerToAdmin_DetectionProposals(t *testing.T) {
	roundTripWorkerToAdmin(t, &plugin_pb.WorkerToAdminMessage{
		WorkerId: "w-3",
		Body: &plugin_pb.WorkerToAdminMessage_DetectionProposals{DetectionProposals: &plugin_pb.DetectionProposals{
			RequestId: "req-1",
			JobType:   "ec",
			HasMore:   true,
			Proposals: []*plugin_pb.JobProposal{
				{
					ProposalId: "p1", DedupeKey: "d1", JobType: "ec",
					Priority: plugin_pb.JobPriority_JOB_PRIORITY_HIGH,
					Summary:  "shard", Detail: "rebuild",
					Parameters: map[string]*plugin_pb.ConfigValue{
						"vol": {Kind: &plugin_pb.ConfigValue_Int64Value{Int64Value: 42}},
					},
					Labels:    map[string]string{"a": "b"},
					NotBefore: timestamppb.New(time.Unix(0, 111)),
					ExpiresAt: timestamppb.New(time.Unix(0, 222)),
				},
			},
		}},
	})
}

func TestWorkerToAdmin_ConfigSchemaResponse_FullDescriptor(t *testing.T) {
	roundTripWorkerToAdmin(t, &plugin_pb.WorkerToAdminMessage{
		WorkerId: "w-4",
		Body: &plugin_pb.WorkerToAdminMessage_ConfigSchemaResponse{ConfigSchemaResponse: &plugin_pb.ConfigSchemaResponse{
			RequestId: "r", JobType: "vacuum", Success: true,
			JobTypeDescriptor: &plugin_pb.JobTypeDescriptor{
				JobType: "vacuum", DisplayName: "Vacuum", Description: "d", Icon: "broom", DescriptorVersion: 3,
				AdminConfigForm: &plugin_pb.ConfigForm{
					FormId: "f", Title: "t", Description: "d",
					Sections: []*plugin_pb.ConfigSection{
						{SectionId: "s1", Title: "S1", Fields: []*plugin_pb.ConfigField{
							{
								Name: "threshold", Label: "Threshold", FieldType: plugin_pb.ConfigFieldType_CONFIG_FIELD_TYPE_DOUBLE,
								Widget: plugin_pb.ConfigWidget_CONFIG_WIDGET_NUMBER, Required: true,
								MinValue: &plugin_pb.ConfigValue{Kind: &plugin_pb.ConfigValue_DoubleValue{DoubleValue: 0}},
								MaxValue: &plugin_pb.ConfigValue{Kind: &plugin_pb.ConfigValue_DoubleValue{DoubleValue: 1}},
								Options:  []*plugin_pb.ConfigOption{{Value: "x", Label: "X", Disabled: true}},
								ValidationRules: []*plugin_pb.ValidationRule{
									{Type: plugin_pb.ValidationRuleType_VALIDATION_RULE_TYPE_REGEX, Expression: ".*", ErrorMessage: "bad"},
								},
								VisibleWhenField:  "enabled",
								VisibleWhenEquals: &plugin_pb.ConfigValue{Kind: &plugin_pb.ConfigValue_BoolValue{BoolValue: true}},
							},
						}},
					},
					DefaultValues: map[string]*plugin_pb.ConfigValue{
						"threshold": {Kind: &plugin_pb.ConfigValue_DoubleValue{DoubleValue: 0.3}},
					},
				},
				AdminRuntimeDefaults: &plugin_pb.AdminRuntimeDefaults{Enabled: true, DetectionIntervalMinutes: 10, RetryLimit: 3, ExecutionTimeoutSeconds: 60},
				WorkerDefaultValues: map[string]*plugin_pb.ConfigValue{
					"mode": {Kind: &plugin_pb.ConfigValue_StringValue{StringValue: "fast"}},
				},
			},
		}},
	})
}

func TestWorkerToAdmin_JobCompleted(t *testing.T) {
	roundTripWorkerToAdmin(t, &plugin_pb.WorkerToAdminMessage{
		WorkerId: "w-5",
		Body: &plugin_pb.WorkerToAdminMessage_JobCompleted{JobCompleted: &plugin_pb.JobCompleted{
			RequestId: "r", JobId: "j", JobType: "ec", Success: false, ErrorMessage: "boom",
			Result: &plugin_pb.JobResult{
				Summary: "done",
				OutputValues: map[string]*plugin_pb.ConfigValue{
					"bytes": {Kind: &plugin_pb.ConfigValue_Int64Value{Int64Value: 1024}},
				},
			},
			Activities: []*plugin_pb.ActivityEvent{
				{Source: plugin_pb.ActivitySource_ACTIVITY_SOURCE_EXECUTOR, Message: "m", Stage: "s", CreatedAt: timestamppb.New(time.Unix(0, 5))},
			},
			CompletedAt: timestamppb.New(time.Unix(0, 9)),
		}},
	})
}

func TestAdminToWorker_RunDetectionRequest(t *testing.T) {
	roundTripAdminToWorker(t, &plugin_pb.AdminToWorkerMessage{
		RequestId: "r-1",
		SentAt:    timestamppb.New(time.Unix(0, 7)),
		Body: &plugin_pb.AdminToWorkerMessage_RunDetectionRequest{RunDetectionRequest: &plugin_pb.RunDetectionRequest{
			RequestId: "r-1", JobType: "vacuum", DetectionSequence: 99,
			AdminRuntime: &plugin_pb.AdminRuntimeConfig{Enabled: true, DetectionIntervalMinutes: 15, MaxJobsPerDetection: 8},
			AdminConfigValues: map[string]*plugin_pb.ConfigValue{
				"a": {Kind: &plugin_pb.ConfigValue_StringValue{StringValue: "x"}},
			},
			WorkerConfigValues: map[string]*plugin_pb.ConfigValue{
				"b": {Kind: &plugin_pb.ConfigValue_BoolValue{BoolValue: true}},
			},
			ClusterContext: &plugin_pb.ClusterContext{
				MasterGrpcAddresses: []string{"m1:1", "m2:2"},
				FilerAddresses:      []string{"f1"},
				Metadata:            map[string]string{"env": "test"},
				S3GrpcAddresses:     []string{"s1"},
			},
			LastSuccessfulRun: timestamppb.New(time.Unix(0, 12345)),
			MaxResults:        50,
		}},
	})
}

func TestAdminToWorker_ExecuteJobRequest(t *testing.T) {
	roundTripAdminToWorker(t, &plugin_pb.AdminToWorkerMessage{
		RequestId: "r-2",
		Body: &plugin_pb.AdminToWorkerMessage_ExecuteJobRequest{ExecuteJobRequest: &plugin_pb.ExecuteJobRequest{
			RequestId: "r-2", Attempt: 2,
			Job: &plugin_pb.JobSpec{
				JobId: "j", JobType: "ec", DedupeKey: "dk", Priority: plugin_pb.JobPriority_JOB_PRIORITY_CRITICAL,
				Summary: "s", Detail: "d",
				Parameters: map[string]*plugin_pb.ConfigValue{
					"shards": {Kind: &plugin_pb.ConfigValue_Int64List{Int64List: &plugin_pb.Int64List{Values: []int64{1, 2, 3}}}},
				},
				Labels:      map[string]string{"team": "x"},
				CreatedAt:   timestamppb.New(time.Unix(0, 1)),
				ScheduledAt: timestamppb.New(time.Unix(0, 2)),
			},
			AdminRuntime: &plugin_pb.AdminRuntimeConfig{RetryLimit: 5},
		}},
	})
}

func TestAdminToWorker_CancelAndShutdownAndHello(t *testing.T) {
	roundTripAdminToWorker(t, &plugin_pb.AdminToWorkerMessage{
		Body: &plugin_pb.AdminToWorkerMessage_CancelRequest{CancelRequest: &plugin_pb.CancelRequest{
			TargetId: "t", TargetKind: plugin_pb.WorkKind_WORK_KIND_DETECTION, Reason: "user", Force: true,
		}},
	})
	roundTripAdminToWorker(t, &plugin_pb.AdminToWorkerMessage{
		Body: &plugin_pb.AdminToWorkerMessage_Shutdown{Shutdown: &plugin_pb.AdminShutdown{Reason: "drain", GracePeriodSeconds: 30}},
	})
	roundTripAdminToWorker(t, &plugin_pb.AdminToWorkerMessage{
		Body: &plugin_pb.AdminToWorkerMessage_Hello{Hello: &plugin_pb.AdminHello{Accepted: true, Message: "hi", HeartbeatIntervalSeconds: 15, ReconnectDelaySeconds: 5}},
	})
	roundTripAdminToWorker(t, &plugin_pb.AdminToWorkerMessage{
		Body: &plugin_pb.AdminToWorkerMessage_RequestConfigSchema{RequestConfigSchema: &plugin_pb.RequestConfigSchema{JobType: "vacuum", ForceRefresh: true}},
	})
}

// TestConfigValue_AllKinds exercises every ConfigValue oneof arm, including the
// recursive list/map members, through a JobResult round-trip.
func TestConfigValue_AllKinds(t *testing.T) {
	values := map[string]*plugin_pb.ConfigValue{
		"bool":     {Kind: &plugin_pb.ConfigValue_BoolValue{BoolValue: true}},
		"int":      {Kind: &plugin_pb.ConfigValue_Int64Value{Int64Value: -7}},
		"double":   {Kind: &plugin_pb.ConfigValue_DoubleValue{DoubleValue: 3.14}},
		"string":   {Kind: &plugin_pb.ConfigValue_StringValue{StringValue: "hi"}},
		"bytes":    {Kind: &plugin_pb.ConfigValue_BytesValue{BytesValue: []byte{1, 2, 3}}},
		"duration": {Kind: &plugin_pb.ConfigValue_DurationValue{DurationValue: durationpb.New(90 * time.Second)}},
		"slist":    {Kind: &plugin_pb.ConfigValue_StringList{StringList: &plugin_pb.StringList{Values: []string{"a", "b"}}}},
		"ilist":    {Kind: &plugin_pb.ConfigValue_Int64List{Int64List: &plugin_pb.Int64List{Values: []int64{9, -9}}}},
		"dlist":    {Kind: &plugin_pb.ConfigValue_DoubleList{DoubleList: &plugin_pb.DoubleList{Values: []float64{1.5, -2.5}}}},
		"blist":    {Kind: &plugin_pb.ConfigValue_BoolList{BoolList: &plugin_pb.BoolList{Values: []bool{true, false, true}}}},
		"vlist":    {Kind: &plugin_pb.ConfigValue_ListValue{ListValue: &plugin_pb.ValueList{Values: []*plugin_pb.ConfigValue{{Kind: &plugin_pb.ConfigValue_Int64Value{Int64Value: 1}}}}}},
		"vmap":     {Kind: &plugin_pb.ConfigValue_MapValue{MapValue: &plugin_pb.ValueMap{Fields: map[string]*plugin_pb.ConfigValue{"n": {Kind: &plugin_pb.ConfigValue_StringValue{StringValue: "m"}}}}}},
	}
	roundTripWorkerToAdmin(t, &plugin_pb.WorkerToAdminMessage{
		Body: &plugin_pb.WorkerToAdminMessage_JobCompleted{JobCompleted: &plugin_pb.JobCompleted{
			RequestId: "r", Success: true,
			Result: &plugin_pb.JobResult{OutputValues: values},
		}},
	})
}

// TestEmptyAndNil verifies absent oneofs and a nil message marshal without
// panicking and decode to an empty-bodied frame.
func TestEmptyAndNil(t *testing.T) {
	if _, err := UnmarshalWorkerToAdmin(MarshalWorkerToAdmin(nil)); err != nil {
		t.Fatalf("nil WorkerToAdmin: %v", err)
	}
	if _, err := UnmarshalAdminToWorker(MarshalAdminToWorker(nil)); err != nil {
		t.Fatalf("nil AdminToWorker: %v", err)
	}
	roundTripWorkerToAdmin(t, &plugin_pb.WorkerToAdminMessage{WorkerId: "only-id"})
}
