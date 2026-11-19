package georecovery

import (
    "testing"
    "time"
)

func TestSelectRecoveryReplica_GeoAware(t *testing.T) {
    replicas := []ReplicaNode{
        {Region: "eu-west-1", Endpoint: "eu.s3.local", RTT: 120 * time.Millisecond, Healthy: true, Available: true},
        {Region: "us-east-1", Endpoint: "us.s3.local", RTT: 18 * time.Millisecond, Healthy: true, Available: true},
        {Region: "us-east-1", Endpoint: "us-down.s3.local", RTT: 5 * time.Millisecond, Healthy: false, Available: true},
    }

    best, ok := SelectRecoveryReplica(replicas, "us-east-1")
    if !ok {
        t.Fatalf("expected healthy replica, got none")
    }
    if best.Endpoint != "us.s3.local" {
        t.Fatalf("expected us.s3.local, got %s", best.Endpoint)
    }
}