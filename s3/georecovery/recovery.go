package georecovery

import (
    "sort"
    "time"
)

// ReplicaNode represents a geo-distributed object storage replica.
type ReplicaNode struct {
    Region    string
    Endpoint  string
    RTT       time.Duration
    Healthy   bool
    Available bool
}

// SelectRecoveryReplica picks the lowest-latency healthy replica (HIP-1349).
func SelectRecoveryReplica(replicas []ReplicaNode, preferredRegion string) (ReplicaNode, bool) {
    candidates := make([]ReplicaNode, 0, len(replicas))
    for _, r := range replicas {
        if r.Healthy && r.Available {
            candidates = append(candidates, r)
        }
    }
    if len(candidates) == 0 {
        return ReplicaNode{}, false
    }

    sort.Slice(candidates, func(i, j int) bool {
        iPref := candidates[i].Region == preferredRegion
        jPref := candidates[j].Region == preferredRegion
        if iPref != jPref {
            return iPref
        }
        return candidates[i].RTT < candidates[j].RTT
    })
    return candidates[0], true
}