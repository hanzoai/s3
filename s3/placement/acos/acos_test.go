package acos

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"testing"
)

func TestCauchyErasureCodecRoundTripAndRecovery(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	sizes := []int{1, 7, 13, 64, 1023, 4096, 16385}

	for _, sz := range sizes {
		payload := make([]byte, sz)
		rng.Read(payload)

		// Test RS(4,2): survive every pair of missing shards (all 15 combinations).
		shards, err := EncodeShards(payload, 4, 2)
		if err != nil {
			t.Fatalf("EncodeShards(sz=%d, 4, 2) error: %v", sz, err)
		}
		if len(shards) != 6 {
			t.Fatalf("expected 6 shards, got %d", len(shards))
		}
		for i := 0; i < 6; i++ {
			for j := i + 1; j < 6; j++ {
				degraded := make([][]byte, 6)
				copy(degraded, shards)
				degraded[i] = nil
				degraded[j] = nil

				got, err := ReconstructPayload(degraded, len(payload), 4, 2)
				if err != nil {
					t.Fatalf("reconstruct(sz=%d, drop=%d,%d) failed: %v", sz, i, j, err)
				}
				if !bytes.Equal(got, payload) {
					t.Fatalf("reconstruct(sz=%d, drop=%d,%d) payload mismatch", sz, i, j)
				}
			}
		}

		// Losing 3 shards (> m=2) MUST fail closed with ErrInsufficientShards.
		degraded3 := make([][]byte, 6)
		copy(degraded3, shards)
		degraded3[0], degraded3[2], degraded3[5] = nil, nil, nil
		if _, err := ReconstructPayload(degraded3, len(payload), 4, 2); !errors.Is(err, ErrInsufficientShards) {
			t.Fatalf("expected ErrInsufficientShards on 3 dropped shards, got %v", err)
		}

		// Test RS(2,1) local AZ stripe: survive any single dropped shard.
		localShards, err := EncodeShards(payload, 2, 1)
		if err != nil {
			t.Fatalf("EncodeShards(sz=%d, 2, 1) error: %v", sz, err)
		}
		for drop := 0; drop < 3; drop++ {
			degraded := make([][]byte, 3)
			copy(degraded, localShards)
			degraded[drop] = nil
			got, err := ReconstructPayload(degraded, len(payload), 2, 1)
			if err != nil {
				t.Fatalf("local reconstruct(sz=%d, drop=%d) failed: %v", sz, drop, err)
			}
			if !bytes.Equal(got, payload) {
				t.Fatalf("local reconstruct(sz=%d, drop=%d) payload mismatch", sz, drop)
			}
		}
	}
}

func TestFullRegionBlackoutFailoverSurvival(t *testing.T) {
	rng := rand.New(rand.NewSource(1001))
	const numObjs = 40

	payloads := make([][]byte, numObjs)
	for i := 0; i < numObjs; i++ {
		payloads[i] = deterministicPayload(rng, 8192, i)
	}

	acosEng := NewEngine(PolicyACOSGeoAwareAdaptive, DefaultTopology(), DefaultACOSConfig())
	localEng := NewEngine(PolicySingleRegion3xLocal, DefaultTopology(), DefaultACOSConfig())

	for i := 0; i < numObjs; i++ {
		key := fmt.Sprintf("active/obj-%03d.bin", i)
		if _, err := acosEng.Put("tenant-prod", key, RegionUSEast, payloads[i]); err != nil {
			t.Fatalf("acos put failed: %v", err)
		}
		if _, err := localEng.Put("tenant-prod", key, RegionUSEast, payloads[i]); err != nil {
			t.Fatalf("local put failed: %v", err)
		}
	}

	// Inject complete regional blackout of us-east.
	acosEng.SetRegionDown(RegionUSEast, true)
	localEng.SetRegionDown(RegionUSEast, true)

	for i := 0; i < numObjs; i++ {
		key := fmt.Sprintf("active/obj-%03d.bin", i)

		got, res, err := acosEng.Get("tenant-prod", key, RegionUSWest)
		if err != nil {
			t.Fatalf("acos failed during us-east blackout for %s: %v", key, err)
		}
		if !bytes.Equal(got, payloads[i]) {
			t.Fatalf("acos payload mismatch for %s during failover", key)
		}
		if res.ServedFromRegion != RegionUSWest {
			t.Fatalf("expected acos failover from nearest peer us-west, got %s", res.ServedFromRegion)
		}

		if _, _, err := localEng.Get("tenant-prod", key, RegionUSWest); !errors.Is(err, ErrUnavailable) {
			t.Fatalf("expected single_region_3x_local ErrUnavailable during us-east blackout, got %v", err)
		}
	}

	acosSnap := acosEng.Snapshot()
	localSnap := localEng.Snapshot()
	if acosSnap.AffectedRegionAvailabilityPct != 100.0 {
		t.Fatalf("expected acos affected-region availability 100%%, got %.2f%%", acosSnap.AffectedRegionAvailabilityPct)
	}
	if localSnap.AffectedRegionAvailabilityPct != 0.0 {
		t.Fatalf("expected single_region_3x_local affected-region availability 0%%, got %.2f%%", localSnap.AffectedRegionAvailabilityPct)
	}
	if acosSnap.DataCorruptions != 0 {
		t.Fatalf("expected 0 data corruptions, got %d", acosSnap.DataCorruptions)
	}
}

func TestLocalFirstAZRecoveryZeroCrossWANEgress(t *testing.T) {
	rng := rand.New(rand.NewSource(2002))
	payload := deterministicPayload(rng, 16384, 7)

	acosEng := NewEngine(PolicyACOSGeoAwareAdaptive, DefaultTopology(), DefaultACOSConfig())
	ecEng := NewEngine(PolicyStaticCrossRegionEC, DefaultTopology(), DefaultACOSConfig())

	if _, err := acosEng.Put("tenant-a", "hot/model-weights.bin", RegionUSEast, payload); err != nil {
		t.Fatalf("acos put failed: %v", err)
	}
	if _, err := ecEng.Put("tenant-a", "hot/model-weights.bin", RegionUSEast, payload); err != nil {
		t.Fatalf("ec put failed: %v", err)
	}

	// Kill one AZ in us-east.
	acosEng.SetAZDown("us-east-a", true)
	ecEng.SetAZDown("us-east-a", true)

	gotACOS, acosRes, err := acosEng.Get("tenant-a", "hot/model-weights.bin", RegionUSEast)
	if err != nil {
		t.Fatalf("acos AZ degraded read failed: %v", err)
	}
	if !bytes.Equal(gotACOS, payload) {
		t.Fatalf("acos AZ degraded payload mismatch")
	}
	if !acosRes.LocalAZRecovery {
		t.Fatalf("expected acos local AZ parity recovery to trigger")
	}
	if acosRes.CrossRegionEgressBytes != 0 {
		t.Fatalf("expected 0 cross-WAN egress bytes for acos single-AZ recovery, got %d", acosRes.CrossRegionEgressBytes)
	}

	gotEC, ecRes, err := ecEng.Get("tenant-a", "hot/model-weights.bin", RegionUSEast)
	if err != nil {
		t.Fatalf("static EC AZ degraded read failed: %v", err)
	}
	if !bytes.Equal(gotEC, payload) {
		t.Fatalf("static EC AZ degraded payload mismatch")
	}
	if ecRes.CrossRegionEgressBytes == 0 {
		t.Fatalf("expected static_cross_region_ec to incur cross-WAN egress during AZ failure")
	}
	if acosRes.LatencyMs >= ecRes.LatencyMs {
		t.Fatalf("expected acos local AZ recovery latency (%.2f ms) < static EC cross-WAN latency (%.2f ms)",
			acosRes.LatencyMs, ecRes.LatencyMs)
	}
}

func TestStrictTenantSovereigntyJurisdictionEnforcement(t *testing.T) {
	rng := rand.New(rand.NewSource(3003))
	payload := deterministicPayload(rng, 8192, 99)

	// 1. Default 5-region topology (only eu-central in EU): test ALL 5 policies!
	for _, policy := range AllPolicies() {
		eng := NewEngine(policy, DefaultTopology(), DefaultACOSConfig())
		if err := eng.RegisterTenant(EUOnlyTenantPolicy("tenant-eu-bank")); err != nil {
			t.Fatalf("register tenant failed: %v", err)
		}

		// Attempting to write directly to us-east MUST be rejected with ErrSovereigntyViolation.
		if _, err := eng.Put("tenant-eu-bank", "ledger/tx-001.bin", RegionUSEast, payload); !errors.Is(err, ErrSovereigntyViolation) {
			t.Fatalf("policy=%s expected ErrSovereigntyViolation for EU tenant put to us-east, got %v", policy, err)
		}

		// Write to eu-central succeeds and MUST stay strictly inside eu-* regions.
		if _, err := eng.Put("tenant-eu-bank", "ledger/tx-001.bin", RegionEUCentral, payload); err != nil {
			t.Fatalf("policy=%s eu-central put failed: %v", policy, err)
		}

		// Repeated remote reads from us-east MUST NEVER trigger migration outside eu-*!
		for r := 0; r < 12; r++ {
			got, _, err := eng.Get("tenant-eu-bank", "ledger/tx-001.bin", RegionUSEast)
			if err != nil {
				t.Fatalf("policy=%s read from us-east failed: %v", policy, err)
			}
			if !bytes.Equal(got, payload) {
				t.Fatalf("policy=%s payload mismatch", policy)
			}
		}

		regions, ok := eng.ObjectRegions("tenant-eu-bank", "ledger/tx-001.bin")
		if !ok || len(regions) == 0 {
			t.Fatalf("policy=%s missing object regions", policy)
		}
		for _, r := range regions {
			if !strings.HasPrefix(string(r), "eu-") {
				t.Fatalf("policy=%s SOVEREIGNTY LEAK: EU-only shard placed in %s", policy, r)
			}
		}

		// Even during full eu-central blackout, MUST fail closed rather than leaking to non-EU.
		eng.SetRegionDown(RegionEUCentral, true)
		if _, _, err := eng.Get("tenant-eu-bank", "ledger/tx-001.bin", RegionUSEast); !errors.Is(err, ErrUnavailable) {
			t.Fatalf("policy=%s expected fail-closed ErrUnavailable when sole EU region is down, got %v", policy, err)
		}

		violations, audited := eng.AuditPlacementSovereignty()
		if violations != 0 || audited == 0 {
			t.Fatalf("policy=%s sovereignty audit failed: violations=%d audited=%d", policy, violations, audited)
		}
	}

	// 2. Multi-EU topology (eu-central + eu-west): failover MUST succeed strictly within eu-*.
	multiEU := NewEngine(PolicyACOSGeoAwareAdaptive, MultiEUTopology(), DefaultACOSConfig())
	if err := multiEU.RegisterTenant(EUOnlyTenantPolicy("tenant-eu-bank")); err != nil {
		t.Fatalf("register tenant failed: %v", err)
	}
	if _, err := multiEU.Put("tenant-eu-bank", "ledger/tx-002.bin", RegionEUCentral, payload); err != nil {
		t.Fatalf("multi-EU put failed: %v", err)
	}
	multiEU.SetRegionDown(RegionEUCentral, true)
	got, res, err := multiEU.Get("tenant-eu-bank", "ledger/tx-002.bin", RegionEUWest)
	if err != nil {
		t.Fatalf("multi-EU failover read failed: %v", err)
	}
	if !bytes.Equal(got, payload) {
		t.Fatalf("multi-EU failover payload mismatch")
	}
	if res.ServedFromRegion != RegionEUWest {
		t.Fatalf("expected failover served from eu-west, got %s", res.ServedFromRegion)
	}
	violations, _ := multiEU.AuditPlacementSovereignty()
	if violations != 0 {
		t.Fatalf("expected 0 sovereignty violations in multi-EU failover, got %d", violations)
	}
}

func TestTenantNamespaceIsolationAndKeyValidation(t *testing.T) {
	eng := NewEngine(PolicyACOSGeoAwareAdaptive, DefaultTopology(), DefaultACOSConfig())
	payloadA := []byte("secret-payload-tenant-alpha")

	if _, err := eng.Put("tenant-alpha", "models/weights.bin", RegionUSEast, payloadA); err != nil {
		t.Fatalf("put alpha failed: %v", err)
	}

	// Cross-tenant read for same key MUST return ErrNotFound.
	if _, _, err := eng.Get("tenant-beta", "models/weights.bin", RegionUSEast); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound for cross-tenant read, got %v", err)
	}

	// Invalid tenant IDs and object keys MUST be rejected.
	badTenants := []string{"", "..", "tenant/evil", "tenant with space", "bad\x00id"}
	for _, bt := range badTenants {
		if _, err := eng.Put(bt, "ok/key.bin", RegionUSEast, payloadA); !errors.Is(err, ErrInvalidTenantID) {
			t.Fatalf("expected ErrInvalidTenantID for %q, got %v", bt, err)
		}
	}

	badKeys := []string{"", "/leading-slash", "a/../b", "a\\b", "null\x00byte", "a//b"}
	for _, bk := range badKeys {
		if _, err := eng.Put("tenant-alpha", bk, RegionUSEast, payloadA); !errors.Is(err, ErrInvalidObjectKey) {
			t.Fatalf("expected ErrInvalidObjectKey for %q, got %v", bk, err)
		}
	}
}

func TestNegativeRegimeOverAggressiveMigrationThrashing(t *testing.T) {
	res := RunNegativeOverAggressiveMigration(20260414, DefaultACOSConfig())
	if len(res.Metrics) != 2 {
		t.Fatalf("expected 2 metrics snapshots, got %d", len(res.Metrics))
	}
	acosSnap := res.Metrics[0]
	negSnap := res.Metrics[1]

	if negSnap.MigrationCount == 0 {
		t.Fatalf("expected negative control to thrash migrations, got 0")
	}
	if acosSnap.MigrationCount != 0 {
		t.Fatalf("expected ACOS hysteresis to suppress 50/50 ping-pong thrashing, got %d migrations", acosSnap.MigrationCount)
	}
	egressRatio := negSnap.EgressCostUSDTotal / acosSnap.EgressCostUSDTotal
	if egressRatio < 3.0 {
		t.Fatalf("expected over_aggressive_migration_negative to burn >=3.0x egress cost, got %.2fx (neg=$%.6f, acos=$%.6f)",
			egressRatio, negSnap.EgressCostUSDTotal, acosSnap.EgressCostUSDTotal)
	}
	if negSnap.ReadP99Ms <= acosSnap.ReadP99Ms {
		t.Fatalf("expected over_aggressive_migration_negative p99 latency (%.2f ms) > ACOS (%.2f ms)",
			negSnap.ReadP99Ms, acosSnap.ReadP99Ms)
	}
}

func TestNegativeRegimeSimultaneousDualRegionCorrelatedOutage(t *testing.T) {
	res := RunNegativeDualRegionCorrelatedOutage(20260414, DefaultACOSConfig())
	if len(res.Metrics) != 3 {
		t.Fatalf("expected 3 snapshots, got %d", len(res.Metrics))
	}
	acosSnap := res.Metrics[0]
	ecSnap := res.Metrics[1]
	globalSnap := res.Metrics[2]

	// Both f=1 policies must fail closed deterministically with 0% availability on dual blackout of us-east + us-west.
	if acosSnap.OverallAvailabilityPct != 0.0 || acosSnap.DataCorruptions != 0 {
		t.Fatalf("expected ACOS clean fail-closed (0%% avail, 0 corruptions), got avail=%.2f%% corruptions=%d",
			acosSnap.OverallAvailabilityPct, acosSnap.DataCorruptions)
	}
	if ecSnap.OverallAvailabilityPct != 0.0 || ecSnap.DataCorruptions != 0 {
		t.Fatalf("expected static EC clean fail-closed (0%% avail, 0 corruptions), got avail=%.2f%% corruptions=%d",
			ecSnap.OverallAvailabilityPct, ecSnap.DataCorruptions)
	}
	if globalSnap.OverallAvailabilityPct != 100.0 {
		t.Fatalf("expected 5x naive global replication to survive dual-region blackout, got %.2f%%", globalSnap.OverallAvailabilityPct)
	}
}

func TestNegativeRegimeWriteHeavyColdOnceObjects(t *testing.T) {
	res := RunNegativeWriteHeavyColdOnce(20260414, DefaultACOSConfig())
	if len(res.Metrics) != 2 {
		t.Fatalf("expected 2 snapshots, got %d", len(res.Metrics))
	}
	localSnap := res.Metrics[0]
	acosSnap := res.Metrics[1]

	if acosSnap.MigrationCount != 0 {
		t.Fatalf("expected 0 migrations on write-once cold workload, got %d", acosSnap.MigrationCount)
	}
	if acosSnap.ReadP50Ms < localSnap.ReadP50Ms {
		t.Fatalf("expected no read latency gain for ACOS over local on cold write-once reads")
	}
	if acosSnap.MetadataOverheadBytes <= localSnap.MetadataOverheadBytes {
		t.Fatalf("expected ACOS metadata tracking overhead (%d B) > static local (%d B)",
			acosSnap.MetadataOverheadBytes, localSnap.MetadataOverheadBytes)
	}
}

func TestDeterministicSeededBenchmarkSuite(t *testing.T) {
	cfg := SuiteConfig{
		Seed:            20260414,
		ObjectCount:     120,
		PayloadSizeByte: 4096,
		ReadOps:         1000,
		Config:          DefaultACOSConfig(),
	}

	r1, err := RunBenchmarkSuite(cfg)
	if err != nil {
		t.Fatalf("run 1 failed: %v", err)
	}
	r2, err := RunBenchmarkSuite(cfg)
	if err != nil {
		t.Fatalf("run 2 failed: %v", err)
	}

	if len(r1.MatchedPolicies) != 5 || len(r2.MatchedPolicies) != 5 {
		t.Fatalf("expected 5 matched policies")
	}

	byPolicy := make(map[PolicyName]MetricsSnapshot)
	for i := range r1.MatchedPolicies {
		m1 := r1.MatchedPolicies[i]
		m2 := r2.MatchedPolicies[i]
		byPolicy[m1.Policy] = m1

		// All deterministic metrics MUST match bit-for-bit across runs with the same seed.
		if m1.ReadP50Ms != m2.ReadP50Ms ||
			m1.ReadP99Ms != m2.ReadP99Ms ||
			m1.WriteP50Ms != m2.WriteP50Ms ||
			m1.OverallAvailabilityPct != m2.OverallAvailabilityPct ||
			m1.AffectedRegionAvailabilityPct != m2.AffectedRegionAvailabilityPct ||
			m1.TotalCrossRegionEgressBytes != m2.TotalCrossRegionEgressBytes ||
			m1.TotalMonthlyTCOUSDPerTB != m2.TotalMonthlyTCOUSDPerTB ||
			m1.MigrationCount != m2.MigrationCount {
			t.Fatalf("non-deterministic metrics for policy %s: run1=%+v run2=%+v", m1.Policy, m1, m2)
		}
		if m1.SovereigntyViolations != 0 {
			t.Fatalf("policy %s had %d sovereignty violations", m1.Policy, m1.SovereigntyViolations)
		}
		if m1.DataCorruptions != 0 {
			t.Fatalf("policy %s had %d data corruptions", m1.Policy, m1.DataCorruptions)
		}
	}

	// Verify key architectural separation across policies.
	acosM := byPolicy[PolicyACOSGeoAwareAdaptive]
	localM := byPolicy[PolicySingleRegion3xLocal]
	globalM := byPolicy[PolicyNaiveGlobalFullReplication]
	ecM := byPolicy[PolicyStaticCrossRegionEC]
	negM := byPolicy[PolicyOverAggressiveMigrationNegative]

	if localM.AffectedRegionAvailabilityPct != 0.0 {
		t.Fatalf("expected single_region_3x_local 0%% affected-region blackout availability, got %.2f%%", localM.AffectedRegionAvailabilityPct)
	}
	if acosM.AffectedRegionAvailabilityPct != 100.0 {
		t.Fatalf("expected acos_geo_aware_adaptive 100%% affected-region blackout availability, got %.2f%%", acosM.AffectedRegionAvailabilityPct)
	}
	if acosM.StorageCostUSDPerTBMonth >= globalM.StorageCostUSDPerTBMonth {
		t.Fatalf("expected ACOS storage cost ($%.2f) < naive global ($%.2f)", acosM.StorageCostUSDPerTBMonth, globalM.StorageCostUSDPerTBMonth)
	}
	if acosM.ReadP50Ms >= ecM.ReadP50Ms {
		t.Fatalf("expected ACOS read p50 (%.2f ms) < static cross-region EC (%.2f ms)", acosM.ReadP50Ms, ecM.ReadP50Ms)
	}
	if negM.EgressCostUSDTotal <= acosM.EgressCostUSDTotal {
		t.Fatalf("expected negative control egress ($%.6f) > ACOS ($%.6f)", negM.EgressCostUSDTotal, acosM.EgressCostUSDTotal)
	}
}

func TestConcurrentMultiTenantPlacementAndFailover(t *testing.T) {
	eng := NewEngine(PolicyACOSGeoAwareAdaptive, MultiEUTopology(), DefaultACOSConfig())
	if err := eng.RegisterTenant(DefaultGlobalTenantPolicy("tenant-global")); err != nil {
		t.Fatalf("register global tenant: %v", err)
	}
	if err := eng.RegisterTenant(EUOnlyTenantPolicy("tenant-eu")); err != nil {
		t.Fatalf("register eu tenant: %v", err)
	}

	const workers = 8
	const objsPerWorker = 12
	var wg sync.WaitGroup
	errCh := make(chan error, workers*objsPerWorker*4)

	// Pre-populate objects concurrently across global and EU-only tenants.
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			rng := rand.New(rand.NewSource(int64(5000 + id)))
			tenant := "tenant-global"
			regions := []RegionID{RegionUSEast, RegionUSWest, RegionEUCentral, RegionAPNortheast, RegionSAEast}
			if id%2 == 1 {
				tenant = "tenant-eu"
				regions = []RegionID{RegionEUCentral, RegionEUWest}
			}
			for i := 0; i < objsPerWorker; i++ {
				key := fmt.Sprintf("w%d/obj-%03d.bin", id, i)
				origin := regions[(id+i)%len(regions)]
				payload := deterministicPayload(rng, 2048, id*100+i)
				if _, err := eng.Put(tenant, key, origin, payload); err != nil {
					errCh <- fmt.Errorf("put %s/%s: %w", tenant, key, err)
					return
				}
				// Concurrent reads from varied regions (including cross-region readers).
				for r := 0; r < 6; r++ {
					reader := regions[(id+i+r)%len(regions)]
					got, _, err := eng.Get(tenant, key, reader)
					if err != nil && !errors.Is(err, ErrUnavailable) {
						errCh <- fmt.Errorf("get %s/%s: %w", tenant, key, err)
						return
					}
					if err == nil && !bytes.Equal(got, payload) {
						errCh <- fmt.Errorf("corrupted payload for %s/%s", tenant, key)
						return
					}
				}
			}
		}(w)
	}

	// Concurrent fault injector toggling AZ and single-region blackouts while traffic runs.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for cycle := 0; cycle < 4; cycle++ {
			eng.SetAZDown("us-east-a", true)
			eng.SetAZDown("eu-central-a", true)
			eng.SetAZDown("us-east-a", false)
			eng.SetRegionDown(RegionUSEast, true)
			eng.SetRegionDown(RegionUSEast, false)
			eng.SetRegionDown(RegionEUCentral, true)
			eng.SetRegionDown(RegionEUCentral, false)
			eng.HealAll()
		}
	}()

	wg.Wait()
	close(errCh)

	for err := range errCh {
		t.Fatalf("concurrent execution error: %v", err)
	}

	violations, audited := eng.AuditPlacementSovereignty()
	if violations != 0 || audited == 0 {
		t.Fatalf("sovereignty audit failed under concurrency: violations=%d audited=%d", violations, audited)
	}
	snap := eng.Snapshot()
	if snap.DataCorruptions != 0 {
		t.Fatalf("expected 0 data corruptions under concurrency, got %d", snap.DataCorruptions)
	}
	if snap.SovereigntyViolations != 0 {
		t.Fatalf("expected 0 runtime sovereignty violations under concurrency, got %d", snap.SovereigntyViolations)
	}
}
