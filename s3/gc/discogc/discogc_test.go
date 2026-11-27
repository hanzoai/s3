package discogc

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"testing"
)

func TestBitExactSHA256PreservationAcrossPolicies(t *testing.T) {
	for _, policy := range AllPolicies() {
		t.Run(string(policy), func(t *testing.T) {
			cfg := DefaultConfig(policy)
			cfg.SegmentCapacityBytes = 512 * 1024
			st, err := Open(cfg)
			if err != nil {
				t.Fatalf("Open(%s): %v", policy, err)
			}
			defer st.Close()

			rng := rand.New(rand.NewSource(42))
			expected := make(map[string][32]byte)
			expectedBytes := make(map[string][]byte)
			tenant := "verify-tenant"

			for i := 0; i < 60; i++ {
				key := fmt.Sprintf("k-%03d", i)
				size := 8192 + rng.Intn(24*1024)
				payload := deterministicPayload(rng, size, uint64(i))
				if err := st.Put(tenant, key, payload); err != nil {
					t.Fatalf("Put(%s): %v", key, err)
				}
				expected[key] = sha256.Sum256(payload)
				expectedBytes[key] = payload
			}

			// Overwrite 25 keys and delete 20 keys across multiple segments.
			for i := 0; i < 25; i++ {
				key := fmt.Sprintf("k-%03d", i)
				payload := deterministicPayload(rng, 16384+rng.Intn(16384), uint64(1000+i))
				if err := st.Put(tenant, key, payload); err != nil {
					t.Fatalf("Overwrite(%s): %v", key, err)
				}
				expected[key] = sha256.Sum256(payload)
				expectedBytes[key] = payload
			}
			for i := 25; i < 45; i++ {
				key := fmt.Sprintf("k-%03d", i)
				ok, err := st.Delete(tenant, key)
				if err != nil || !ok {
					t.Fatalf("Delete(%s): ok=%v err=%v", key, ok, err)
				}
				delete(expected, key)
				delete(expectedBytes, key)
			}

			if err := st.SealActive(tenant); err != nil {
				t.Fatalf("SealActive: %v", err)
			}
			if err := st.RunGC(tenant); err != nil {
				t.Fatalf("RunGC: %v", err)
			}

			verified, err := st.VerifyAllLive()
			if err != nil {
				t.Fatalf("VerifyAllLive(%s): %v", policy, err)
			}
			if verified != len(expected) {
				t.Fatalf("verified=%d, want %d", verified, len(expected))
			}

			for key, wantHash := range expected {
				got, err := st.Get(tenant, key)
				if err != nil {
					t.Fatalf("Get(%s): %v", key, err)
				}
				if sha256.Sum256(got) != wantHash || !bytes.Equal(got, expectedBytes[key]) {
					t.Fatalf("payload mismatch for key %s under policy %s", key, policy)
				}
			}
			for i := 25; i < 45; i++ {
				key := fmt.Sprintf("k-%03d", i)
				if _, err := st.Get(tenant, key); err != ErrNotFound {
					t.Fatalf("deleted key %s returned err=%v, want ErrNotFound", key, err)
				}
			}
		})
	}
}

func TestLinuxRealPhysicalBlockReclamation(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("real FALLOC_FL_PUNCH_HOLE st_blocks verification requires Linux")
	}

	cfg := DefaultConfig(PolicyAdaptiveDiscard)
	cfg.SegmentCapacityBytes = 8 * 1024 * 1024 // keep all objects in one segment
	cfg.MaxExtentFragments = 0                 // isolate pure punch-hole without defrag
	st, err := Open(cfg)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer st.Close()

	rng := rand.New(rand.NewSource(777))
	const numObjs = 32
	const objSize = 64 * 1024 // 64 KiB -> 16 blocks of 4 KiB
	for i := 0; i < numObjs; i++ {
		payload := deterministicPayload(rng, objSize-64, uint64(i))
		if err := st.Put("t1", fmt.Sprintf("big-%03d", i), payload); err != nil {
			t.Fatalf("Put: %v", err)
		}
	}
	beforeMetrics := st.Metrics("pre_discard")
	if beforeMetrics.PhysicalDiskBytes < numObjs*(objSize-4096) {
		t.Fatalf("unexpected pre-discard physical disk usage: %d", beforeMetrics.PhysicalDiskBytes)
	}

	// Delete half the objects (16 * 64 KiB = 1 MiB dead).
	for i := 0; i < numObjs; i += 2 {
		if _, err := st.Delete("t1", fmt.Sprintf("big-%03d", i)); err != nil {
			t.Fatalf("Delete: %v", err)
		}
	}
	if err := st.RunGC("t1"); err != nil {
		t.Fatalf("RunGC: %v", err)
	}
	afterMetrics := st.Metrics("post_discard")
	if !afterMetrics.PunchHoleSupported {
		t.Skip("underlying test filesystem does not support FALLOC_FL_PUNCH_HOLE")
	}
	if afterMetrics.OSBlocksReclaimedBytes <= 0 {
		t.Fatalf("expected positive OS blocks reclaimed via FALLOC_FL_PUNCH_HOLE, got %d", afterMetrics.OSBlocksReclaimedBytes)
	}
	if afterMetrics.PhysicalDiskBytes >= beforeMetrics.PhysicalDiskBytes {
		t.Fatalf("expected physical disk bytes to shrink after punch-hole: before=%d after=%d",
			beforeMetrics.PhysicalDiskBytes, afterMetrics.PhysicalDiskBytes)
	}
	if afterMetrics.GCCopiedLiveBytes != 0 {
		t.Fatalf("expected 0 live bytes copied during pure discard, got %d", afterMetrics.GCCopiedLiveBytes)
	}
	if afterMetrics.LogicalFileBytes != beforeMetrics.LogicalFileBytes {
		t.Fatalf("expected FALLOC_FL_KEEP_SIZE to preserve logical file size: before=%d after=%d",
			beforeMetrics.LogicalFileBytes, afterMetrics.LogicalFileBytes)
	}
	if verified, err := st.VerifyAllLive(); err != nil || verified != numObjs/2 {
		t.Fatalf("VerifyAllLive: verified=%d err=%v", verified, err)
	}
}

func TestWAFReductionAdaptiveDiscardVsGreedyCompaction(t *testing.T) {
	spec := WorkloadSpec{
		Name:            "test_uniform_churn",
		Seed:            20260225,
		KeyCount:        160,
		ChurnOps:        700,
		MinObjBytes:     32 * 1024,
		MaxObjBytes:     32 * 1024,
		OverwriteProb:   0.75,
		DeleteProb:      0.15,
		SegmentCapBytes: 1 * 1024 * 1024,
		SampleEveryOps:  50,
	}

	drive := DefaultEnterpriseNVMe()
	greedy, err := RunWorkload(spec, PolicyGreedyCompaction, drive, "")
	if err != nil {
		t.Fatalf("RunWorkload greedy: %v", err)
	}
	adaptive, err := RunWorkload(spec, PolicyAdaptiveDiscard, drive, "")
	if err != nil {
		t.Fatalf("RunWorkload adaptive: %v", err)
	}

	if adaptive.FlashWear.LogicalWAF >= greedy.FlashWear.LogicalWAF {
		t.Fatalf("expected adaptive discard LogicalWAF (%.4f) < greedy copy compaction LogicalWAF (%.4f)",
			adaptive.FlashWear.LogicalWAF, greedy.FlashWear.LogicalWAF)
	}
	if adaptive.GCCopiedLiveBytes >= greedy.GCCopiedLiveBytes {
		t.Fatalf("expected adaptive copied bytes (%d) < greedy copied bytes (%d)",
			adaptive.GCCopiedLiveBytes, greedy.GCCopiedLiveBytes)
	}
	if adaptive.FlashWear.EffectivePhysicalWAF >= greedy.FlashWear.EffectivePhysicalWAF {
		t.Fatalf("expected adaptive physical WAF (%.4f) < greedy physical WAF (%.4f)",
			adaptive.FlashWear.EffectivePhysicalWAF, greedy.FlashWear.EffectivePhysicalWAF)
	}
	if adaptive.Cost.TotalUSDPerTBMonth >= greedy.Cost.TotalUSDPerTBMonth {
		t.Fatalf("expected adaptive TCO ($%.4f/TB-mo) < greedy TCO ($%.4f/TB-mo)",
			adaptive.Cost.TotalUSDPerTBMonth, greedy.Cost.TotalUSDPerTBMonth)
	}
	if adaptive.FlashWear.EstimatedDriveLifeYears < greedy.FlashWear.EstimatedDriveLifeYears {
		t.Fatalf("expected adaptive drive life (%.2f yr) >= greedy drive life (%.2f yr)",
			adaptive.FlashWear.EstimatedDriveLifeYears, greedy.FlashWear.EstimatedDriveLifeYears)
	}
}

func TestNegativeRegimeBlindSubblockDiscardZeroReclaim(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("real st_blocks sub-block verification requires Linux")
	}

	regimes, err := RunBoundaryRegimes(DefaultEnterpriseNVMe())
	if err != nil {
		t.Fatalf("RunBoundaryRegimes: %v", err)
	}

	// Regime 1: sub_block_fragmentation_no_reclaim
	subReg := regimes[0]
	if subReg.Name != "sub_block_fragmentation_no_reclaim" {
		t.Fatalf("unexpected regime 0 name: %s", subReg.Name)
	}
	blind := subReg.Results[0]
	adaptive := subReg.Results[1]

	if blind.PunchHoleCalls == 0 {
		t.Fatalf("expected blind subblock negative control to issue PunchHole calls")
	}
	if blind.OSBlocksReclaimedBytes != 0 {
		t.Fatalf("expected blind subblock discard to reclaim 0 OS bytes on interleaved 512 B objects, got %d",
			blind.OSBlocksReclaimedBytes)
	}
	if adaptive.CompactionCycles == 0 {
		t.Fatalf("expected discogc_adaptive_discard to fall back to compaction on sub-block fragmentation")
	}
	if adaptive.PhysicalDiskBytes >= blind.PhysicalDiskBytes {
		t.Fatalf("expected adaptive fallback physical bytes (%d) < blind unreclaimed bytes (%d)",
			adaptive.PhysicalDiskBytes, blind.PhysicalDiskBytes)
	}

	// Regime 2: extent_map_fragmentation_tax
	fragReg := regimes[1]
	uncapped := fragReg.Results[0]
	bounded := fragReg.Results[1]
	if uncapped.ExtentFragments < 200 {
		t.Fatalf("expected uncapped checkerboard punch-hole to leave >= 200 fragments, got %d", uncapped.ExtentFragments)
	}
	if bounded.ExtentFragments > 64 {
		t.Fatalf("expected bounded DisCoGC defrag to cap fragments <= 64, got %d", bounded.ExtentFragments)
	}
	if bounded.CompactionCycles == 0 {
		t.Fatalf("expected bounded DisCoGC to trigger defrag compaction")
	}

	// Regime 3: append_only_zero_delete_workload
	zeroReg := regimes[2]
	baseLWA := zeroReg.Results[0].FlashWear.LogicalWAF
	for _, r := range zeroReg.Results[1:] {
		if math.Abs(r.FlashWear.LogicalWAF-baseLWA) > 1e-6 {
			t.Fatalf("zero-delete regime LWA mismatch for %s: got %.6f, want %.6f",
				r.Policy, r.FlashWear.LogicalWAF, baseLWA)
		}
		if r.GCCopiedLiveBytes != 0 || r.OSBlocksReclaimedBytes != 0 {
			t.Fatalf("zero-delete regime expected 0 GC work for %s, got copied=%d reclaimed=%d",
				r.Policy, r.GCCopiedLiveBytes, r.OSBlocksReclaimedBytes)
		}
	}
}

func TestConcurrentMultiTenantIsolationRace(t *testing.T) {
	cfg := DefaultConfig(PolicyAdaptiveDiscard)
	cfg.SegmentCapacityBytes = 256 * 1024
	st, err := Open(cfg)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer st.Close()

	const numTenants = 8
	const opsPerTenant = 80
	var wg sync.WaitGroup
	errCh := make(chan error, numTenants)

	for tid := 0; tid < numTenants; tid++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			tenant := fmt.Sprintf("tenant-%02d", id)
			otherTenant := fmt.Sprintf("tenant-%02d", (id+1)%numTenants)
			rng := rand.New(rand.NewSource(int64(1000 + id)))

			live := make(map[string][]byte)
			for i := 0; i < opsPerTenant; i++ {
				key := fmt.Sprintf("key-%02d", i%20)
				if i%4 == 3 {
					if _, exists := live[key]; exists {
						if _, err := st.Delete(tenant, key); err != nil {
							errCh <- fmt.Errorf("%s delete %s: %w", tenant, key, err)
							return
						}
						delete(live, key)
					}
					continue
				}
				payload := deterministicPayload(rng, 8192+rng.Intn(16384), uint64(id*10000+i))
				if err := st.Put(tenant, key, payload); err != nil {
					errCh <- fmt.Errorf("%s put %s: %w", tenant, key, err)
					return
				}
				live[key] = payload

				// Verify immediate read isolation: own key matches exact bytes.
				got, err := st.Get(tenant, key)
				if err != nil || !bytes.Equal(got, payload) {
					errCh <- fmt.Errorf("%s get %s mismatch: err=%v", tenant, key, err)
					return
				}

				// Verify tenant namespace isolation on tenant-unique keys.
				uniqueKey := fmt.Sprintf("private-%s-%02d", tenant, i%5)
				if err := st.Put(tenant, uniqueKey, payload[:4096]); err != nil {
					errCh <- fmt.Errorf("%s put unique: %w", tenant, err)
					return
				}
				live[uniqueKey] = payload[:4096]
				if _, err := st.Get(otherTenant, uniqueKey); err != ErrNotFound {
					errCh <- fmt.Errorf("cross-tenant leak: %s read %s's key %s (err=%v)",
						otherTenant, tenant, uniqueKey, err)
					return
				}

				if i%15 == 0 {
					if err := st.RunGC(tenant); err != nil {
						errCh <- fmt.Errorf("%s RunGC: %w", tenant, err)
						return
					}
					_ = st.Metrics("race_check")
				}
			}

			if err := st.SealActive(tenant); err != nil {
				errCh <- fmt.Errorf("%s SealActive: %w", tenant, err)
				return
			}
			if err := st.RunGC(tenant); err != nil {
				errCh <- fmt.Errorf("%s final RunGC: %w", tenant, err)
				return
			}
			for k, want := range live {
				got, err := st.Get(tenant, k)
				if err != nil || !bytes.Equal(got, want) {
					errCh <- fmt.Errorf("%s final verify %s failed: err=%v", tenant, k, err)
					return
				}
			}
		}(tid)
	}

	wg.Wait()
	close(errCh)
	for err := range errCh {
		t.Fatal(err)
	}

	if _, err := st.VerifyAllLive(); err != nil {
		t.Fatalf("VerifyAllLive after concurrent multi-tenant run: %v", err)
	}
}

func TestEconomicsFormulaSanity(t *testing.T) {
	drive := DefaultEnterpriseNVMe()

	lowWAF := &Metrics{
		Policy:                 PolicyAdaptiveDiscard,
		UserBytesWritten:       100 << 20,
		GCCopiedLiveBytes:      2 << 20,
		MetadataBytesWritten:   1 << 20,
		AvgSpaceAmplification:  1.20,
		OSBlocksReclaimedBytes: 50 << 20,
		PunchHoleCalls:         500,
	}
	highWAF := &Metrics{
		Policy:                PolicyGreedyCompaction,
		UserBytesWritten:      100 << 20,
		GCCopiedLiveBytes:     150 << 20,
		MetadataBytesWritten:  5 << 20,
		AvgSpaceAmplification: 1.30,
	}
	bloatedNoGC := &Metrics{
		Policy:                PolicyNoGC,
		UserBytesWritten:      100 << 20,
		GCCopiedLiveBytes:     0,
		MetadataBytesWritten:  1 << 20,
		AvgSpaceAmplification: 3.80,
	}

	wLow, cLow := ComputeEconomics(lowWAF, drive)
	wHigh, cHigh := ComputeEconomics(highWAF, drive)
	wNoGC, cNoGC := ComputeEconomics(bloatedNoGC, drive)

	if wLow.EffectivePhysicalWAF >= wHigh.EffectivePhysicalWAF {
		t.Fatalf("expected lowWAF (%.2f) < highWAF (%.2f)", wLow.EffectivePhysicalWAF, wHigh.EffectivePhysicalWAF)
	}
	if wHigh.EstimatedDriveLifeYears >= drive.WarrantyYears {
		t.Fatalf("expected high WAF drive life (%.2f) < warranty (%.2f)", wHigh.EstimatedDriveLifeYears, drive.WarrantyYears)
	}
	if cLow.TotalUSDPerTBMonth >= cHigh.TotalUSDPerTBMonth {
		t.Fatalf("expected lowWAF TCO ($%.2f) < highWAF TCO ($%.2f)", cLow.TotalUSDPerTBMonth, cHigh.TotalUSDPerTBMonth)
	}
	if cLow.TotalUSDPerTBMonth >= cNoGC.TotalUSDPerTBMonth {
		t.Fatalf("expected lowWAF TCO ($%.2f) < bloated NoGC TCO ($%.2f)", cLow.TotalUSDPerTBMonth, cNoGC.TotalUSDPerTBMonth)
	}
	if wNoGC.LogicalWAF > 1.05 {
		t.Fatalf("expected NoGC LWA near 1.00, got %.4f", wNoGC.LogicalWAF)
	}
}
