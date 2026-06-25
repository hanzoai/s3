// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package master

import (
	"testing"

	master_pb "github.com/hanzoai/s3/s3/pb/master_pb"
)

// TestHeartbeatReqRepeatedFieldsRoundTrip pins the five repeated fields the
// decode used to drop: a volume server reports incremental topology churn
// (new/deleted volumes and EC shards, plus disk tags) through the heartbeat,
// so dropping them on decode would leave the master blind to it.
func TestHeartbeatReqRepeatedFieldsRoundTrip(t *testing.T) {
	in := &master_pb.Heartbeat{
		Ip:              "10.0.0.1",
		Port:            8080,
		NewVolumes:      []*master_pb.VolumeShortInformationMessage{{Id: 1, Collection: "c", DiskId: 2}},
		DeletedVolumes:  []*master_pb.VolumeShortInformationMessage{{Id: 3, Collection: "d"}},
		NewEcShards:     []*master_pb.VolumeEcShardInformationMessage{{Id: 4, Collection: "e"}},
		DeletedEcShards: []*master_pb.VolumeEcShardInformationMessage{{Id: 5}},
		DiskTags:        []*master_pb.DiskTag{{DiskId: 6, Tags: []string{"ssd", "fast"}}},
	}

	out, err := HeartbeatReqFromWire(HeartbeatReqToWire(in))
	if err != nil {
		t.Fatalf("HeartbeatReqFromWire: %v", err)
	}

	if len(out.NewVolumes) != 1 || out.NewVolumes[0].Id != 1 || out.NewVolumes[0].Collection != "c" || out.NewVolumes[0].DiskId != 2 {
		t.Errorf("NewVolumes not round-tripped: %+v", out.NewVolumes)
	}
	if len(out.DeletedVolumes) != 1 || out.DeletedVolumes[0].Id != 3 || out.DeletedVolumes[0].Collection != "d" {
		t.Errorf("DeletedVolumes not round-tripped: %+v", out.DeletedVolumes)
	}
	if len(out.NewEcShards) != 1 || out.NewEcShards[0].Id != 4 || out.NewEcShards[0].Collection != "e" {
		t.Errorf("NewEcShards not round-tripped: %+v", out.NewEcShards)
	}
	if len(out.DeletedEcShards) != 1 || out.DeletedEcShards[0].Id != 5 {
		t.Errorf("DeletedEcShards not round-tripped: %+v", out.DeletedEcShards)
	}
	if len(out.DiskTags) != 1 || out.DiskTags[0].DiskId != 6 || len(out.DiskTags[0].Tags) != 2 || out.DiskTags[0].Tags[0] != "ssd" {
		t.Errorf("DiskTags not round-tripped: %+v", out.DiskTags)
	}
}
