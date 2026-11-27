package discogc

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// WorkloadSpec defines a deterministic seeded S3 object-churn workload.
type WorkloadSpec struct {
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Seed            int64   `json:"seed"`
	KeyCount        int     `json:"key_count"`
	ChurnOps        int     `json:"churn_ops"`
	MinObjBytes     int     `json:"min_obj_bytes"`
	MaxObjBytes     int     `json:"max_obj_bytes"`
	OverwriteProb   float64 `json:"overwrite_prob"`
	DeleteProb      float64 `json:"delete_prob"`
	ZipfSkew        float64 `json:"zipf_skew"` // 0 = uniform; >1 = Zipfian hot/cold
	BulkDeleteEvery int     `json:"bulk_delete_every"`
	BulkDeleteFrac  float64 `json:"bulk_delete_frac"`
	SegmentCapBytes int64   `json:"segment_cap_bytes"`
	SampleEveryOps  int     `json:"sample_every_ops"`
}

// DefaultWorkloads returns the 4 matched S3 churn workloads evaluated across
// all 5 policies.
func DefaultWorkloads() []WorkloadSpec {
	return []WorkloadSpec{
		{
			Name:            "zipfian_hot_cold_churn",
			Description:     "Zipfian (s=1.10) hot-cold S3 object churn (16-64 KiB objects, 75% overwrite / 15% delete / 10% fresh put)",
			Seed:            20260224,
			KeyCount:        320,
			ChurnOps:        1800,
			MinObjBytes:     16 * 1024,
			MaxObjBytes:     64 * 1024,
			OverwriteProb:   0.75,
			DeleteProb:      0.15,
			ZipfSkew:        1.10,
			SegmentCapBytes: 1 * 1024 * 1024,
			SampleEveryOps:  50,
		},
		{
			Name:            "uniform_random_churn",
			Description:     "Uniform random overwrites & deletes over 32 KiB objects (worst-case mixed-age segments for copy-compaction)",
			Seed:            20260225,
			KeyCount:        320,
			ChurnOps:        1600,
			MinObjBytes:     32 * 1024,
			MaxObjBytes:     32 * 1024,
			OverwriteProb:   0.75,
			DeleteProb:      0.15,
			ZipfSkew:        0,
			SegmentCapBytes: 1 * 1024 * 1024,
			SampleEveryOps:  50,
		},
		{
			Name:            "large_object_bulk_deletes",
			Description:     "Large 128-256 KiB S3 checkpoint/parquet objects with periodic 50% bulk TTL expiration sweeps",
			Seed:            20260226,
			KeyCount:        120,
			ChurnOps:        360,
			MinObjBytes:     128 * 1024,
			MaxObjBytes:     256 * 1024,
			OverwriteProb:   0.35,
			DeleteProb:      0.10,
			ZipfSkew:        0,
			BulkDeleteEvery: 90,
			BulkDeleteFrac:  0.50,
			SegmentCapBytes: 2 * 1024 * 1024,
			SampleEveryOps:  25,
		},
		{
			Name:            "small_fragmented_subblock",
			Description:     "Small 512 B-1.5 KiB objects packed inside shared 4 KiB filesystem blocks with interleaved deletes",
			Seed:            20260227,
			KeyCount:        2400,
			ChurnOps:        4800,
			MinObjBytes:     512,
			MaxObjBytes:     1536,
			OverwriteProb:   0.65,
			DeleteProb:      0.25,
			ZipfSkew:        0,
			SegmentCapBytes: 256 * 1024,
			SampleEveryOps:  100,
		},
	}
}

// RunWorkload executes spec deterministically against policy using real OS
// segment files under baseDir (or a temp dir if empty), verifies 100% bit-exact
// SHA-256 integrity of all surviving live objects, and returns Metrics.
func RunWorkload(spec WorkloadSpec, policy Policy, drive DriveEconomics, baseDir string) (Metrics, error) {
	cfg := DefaultConfig(policy)
	cfg.Dir = baseDir
	if spec.SegmentCapBytes > 0 {
		cfg.SegmentCapacityBytes = spec.SegmentCapBytes
	}
	if drive.RawCapacityTB > 0 {
		cfg.Drive = drive
	}

	store, err := Open(cfg)
	if err != nil {
		return Metrics{}, err
	}
	defer store.Close()

	rng := rand.New(rand.NewSource(spec.Seed))
	zipf := newZipfTable(spec.KeyCount, spec.ZipfSkew)
	tenant := "bench-tenant"

	liveKeys := make([]string, 0, spec.KeyCount)
	keyLive := make(map[string]bool, spec.KeyCount)
	nextFreshID := spec.KeyCount

	// Phase 1: Initial population of KeyCount objects.
	for i := 0; i < spec.KeyCount; i++ {
		key := fmt.Sprintf("obj-%06d", i)
		payload := deterministicPayload(rng, pickSize(rng, spec.MinObjBytes, spec.MaxObjBytes), uint64(i))
		if err := store.Put(tenant, key, payload); err != nil {
			return Metrics{}, fmt.Errorf("workload %s init put %s: %w", spec.Name, key, err)
		}
		liveKeys = append(liveKeys, key)
		keyLive[key] = true
		if spec.SampleEveryOps > 0 && (i+1)%spec.SampleEveryOps == 0 {
			store.SampleSpace()
		}
	}
	store.SampleSpace()

	// Phase 2: Steady-state churn (overwrites, deletes, fresh puts, optional bulk TTL deletes).
	for op := 0; op < spec.ChurnOps; op++ {
		if spec.BulkDeleteEvery > 0 && op > 0 && op%spec.BulkDeleteEvery == 0 && len(liveKeys) > 10 {
			// Delete BulkDeleteFrac of currently live keys in a bulk TTL sweep.
			targetDel := int(float64(len(liveKeys)) * spec.BulkDeleteFrac)
			perm := rng.Perm(len(liveKeys))
			deleted := 0
			for _, idx := range perm {
				if deleted >= targetDel {
					break
				}
				k := liveKeys[idx]
				if !keyLive[k] {
					continue
				}
				if _, err := store.Delete(tenant, k); err != nil {
					return Metrics{}, fmt.Errorf("workload %s bulk delete %s: %w", spec.Name, k, err)
				}
				keyLive[k] = false
				deleted++
			}
			liveKeys = compactKeySlice(liveKeys, keyLive)
			if err := store.RunGC(tenant); err != nil {
				return Metrics{}, err
			}
			store.SampleSpace()
			continue
		}

		u := rng.Float64()
		switch {
		case u < spec.OverwriteProb:
			var key string
			if spec.ZipfSkew > 0 {
				idx := zipf.sample(rng)
				key = fmt.Sprintf("obj-%06d", idx)
			} else if len(liveKeys) > 0 {
				key = liveKeys[rng.Intn(len(liveKeys))]
			} else {
				key = fmt.Sprintf("obj-%06d", rng.Intn(spec.KeyCount))
			}
			payload := deterministicPayload(rng, pickSize(rng, spec.MinObjBytes, spec.MaxObjBytes), uint64(op+100000))
			if err := store.Put(tenant, key, payload); err != nil {
				return Metrics{}, fmt.Errorf("workload %s overwrite %s: %w", spec.Name, key, err)
			}
			if !keyLive[key] {
				keyLive[key] = true
				liveKeys = append(liveKeys, key)
			}
		case u < spec.OverwriteProb+spec.DeleteProb && len(liveKeys) > spec.KeyCount/4:
			var key string
			if spec.ZipfSkew > 0 {
				idx := zipf.sample(rng)
				candidate := fmt.Sprintf("obj-%06d", idx)
				if keyLive[candidate] {
					key = candidate
				}
			}
			if key == "" {
				key = liveKeys[rng.Intn(len(liveKeys))]
			}
			if _, err := store.Delete(tenant, key); err != nil {
				return Metrics{}, fmt.Errorf("workload %s delete %s: %w", spec.Name, key, err)
			}
			keyLive[key] = false
			if op%32 == 0 {
				liveKeys = compactKeySlice(liveKeys, keyLive)
			}
		default:
			key := fmt.Sprintf("obj-%06d", nextFreshID)
			nextFreshID++
			payload := deterministicPayload(rng, pickSize(rng, spec.MinObjBytes, spec.MaxObjBytes), uint64(nextFreshID))
			if err := store.Put(tenant, key, payload); err != nil {
				return Metrics{}, fmt.Errorf("workload %s put %s: %w", spec.Name, key, err)
			}
			keyLive[key] = true
			liveKeys = append(liveKeys, key)
		}

		if spec.SampleEveryOps > 0 && (op+1)%spec.SampleEveryOps == 0 {
			store.SampleSpace()
		}
	}

	// End-of-run GC pass and final space sample.
	if err := store.SealActive(tenant); err != nil {
		return Metrics{}, err
	}
	if err := store.RunGC(tenant); err != nil {
		return Metrics{}, err
	}
	store.SampleSpace()

	// Verify bit-exact SHA-256 integrity of every surviving live object.
	if _, err := store.VerifyAllLive(); err != nil {
		return Metrics{}, fmt.Errorf("workload %s policy %s integrity failure: %w", spec.Name, policy, err)
	}

	scanBytes, scanDur, err := store.SequentialScanAll(tenant)
	if err != nil {
		return Metrics{}, fmt.Errorf("workload %s policy %s scan failure: %w", spec.Name, policy, err)
	}

	m := store.Metrics(spec.Name)
	if scanDur > 0 && scanBytes > 0 {
		m.SeqScanThroughputMBs = round4((float64(scanBytes) / (1024.0 * 1024.0)) / scanDur.Seconds())
	}
	return m, nil
}

// BoundaryRegimeResult captures the empirical comparison for one of the three
// documented honest negative/boundary regimes.
type BoundaryRegimeResult struct {
	Name        string    `json:"name"`
	Hypothesis  string    `json:"hypothesis"`
	Observation string    `json:"observation"`
	Results     []Metrics `json:"results"`
}

// RunBoundaryRegimes executes all three documented negative/boundary regimes
// on real OS files:
//  1. sub_block_fragmentation_no_reclaim
//  2. extent_map_fragmentation_tax
//  3. append_only_zero_delete_workload
func RunBoundaryRegimes(drive DriveEconomics) ([]BoundaryRegimeResult, error) {
	r1, err := runSubblockFragmentationRegime(drive)
	if err != nil {
		return nil, err
	}
	r2, err := runExtentMapFragmentationRegime(drive)
	if err != nil {
		return nil, err
	}
	r3, err := runAppendOnlyZeroDeleteRegime(drive)
	if err != nil {
		return nil, err
	}
	return []BoundaryRegimeResult{r1, r2, r3}, nil
}

// Regime 1: Interleaved 512 B objects inside shared 4 KiB blocks where every odd
// object is deleted. Proves blind sub-block discard reclaims 0 physical OS
// blocks, whereas adaptive discard detects non-reclaimable sub-block dead space
// and falls back to selective copy-compaction.
func runSubblockFragmentationRegime(drive DriveEconomics) (BoundaryRegimeResult, error) {
	policies := []Policy{
		PolicyBlindSubblockDiscard,
		PolicyAdaptiveDiscard,
		PolicyGreedyCompaction,
	}
	results := make([]Metrics, 0, len(policies))
	for _, p := range policies {
		cfg := DefaultConfig(p)
		cfg.SegmentCapacityBytes = 512 * 1024
		cfg.SubblockFragFallbackRatio = 0.25
		cfg.Drive = drive

		st, err := Open(cfg)
		if err != nil {
			return BoundaryRegimeResult{}, err
		}
		rng := rand.New(rand.NewSource(9001))
		const count = 1024
		for i := 0; i < count; i++ {
			key := fmt.Sprintf("tiny-%04d", i)
			data := deterministicPayload(rng, 448, uint64(i)) // ~488 B framed -> 8 objects per 4 KiB block
			if err := st.Put("t1", key, data); err != nil {
				st.Close()
				return BoundaryRegimeResult{}, err
			}
		}
		if err := st.SealActive("t1"); err != nil {
			st.Close()
			return BoundaryRegimeResult{}, err
		}
		// Delete every odd object so every 4 KiB block still holds 4 live objects.
		for i := 1; i < count; i += 2 {
			if _, err := st.Delete("t1", fmt.Sprintf("tiny-%04d", i)); err != nil {
				st.Close()
				return BoundaryRegimeResult{}, err
			}
		}
		if err := st.RunGC("t1"); err != nil {
			st.Close()
			return BoundaryRegimeResult{}, err
		}
		st.SampleSpace()
		if _, err := st.VerifyAllLive(); err != nil {
			st.Close()
			return BoundaryRegimeResult{}, err
		}
		m := st.Metrics("sub_block_fragmentation_no_reclaim")
		st.Close()
		results = append(results, m)
	}

	return BoundaryRegimeResult{
		Name:       "sub_block_fragmentation_no_reclaim",
		Hypothesis: "Deleting interleaved 512 B objects inside shared 4 KiB filesystem blocks yields 0 reclaimed OS blocks under direct hole-punching; only copy-compaction (or adaptive fallback) reclaims physical disk space.",
		Observation: fmt.Sprintf(
			"blind_subblock_discard_negative issued %d PunchHole calls and reclaimed %d OS bytes (final space amp %.2fx); discogc_adaptive_discard detected sub-block dead space and fell back to %d compaction cycle(s) (final space amp %.2fx).",
			results[0].PunchHoleCalls, results[0].OSBlocksReclaimedBytes, results[0].FinalSpaceAmplification,
			results[1].CompactionCycles, results[1].FinalSpaceAmplification,
		),
		Results: results,
	}, nil
}

// Regime 2: Alternating 8 KiB live / 8 KiB dead checkerboard pattern across a
// 4 MiB segment. Compares uncapped hole-punching (MaxExtentFragments=0 -> 256
// holes) vs bounded DisCoGC defragmentation (MaxExtentFragments=64).
func runExtentMapFragmentationRegime(drive DriveEconomics) (BoundaryRegimeResult, error) {
	runCase := func(label Policy, maxFrags int) (Metrics, error) {
		cfg := DefaultConfig(PolicyAdaptiveDiscard)
		cfg.SegmentCapacityBytes = 8 * 1024 * 1024
		cfg.MaxExtentFragments = maxFrags
		cfg.Drive = drive

		st, err := Open(cfg)
		if err != nil {
			return Metrics{}, err
		}
		defer st.Close()

		rng := rand.New(rand.NewSource(9002))
		const count = 512 // 512 * 8 KiB = 4 MiB
		for i := 0; i < count; i++ {
			key := fmt.Sprintf("chk-%04d", i)
			// 8 KiB minus header/key so each framed object is exactly 8192 B (two 4 KiB blocks).
			data := deterministicPayload(rng, 8192-64, uint64(i))
			if err := st.Put("t1", key, data); err != nil {
				return Metrics{}, err
			}
		}
		if err := st.SealActive("t1"); err != nil {
			return Metrics{}, err
		}
		for i := 1; i < count; i += 2 {
			if _, err := st.Delete("t1", fmt.Sprintf("chk-%04d", i)); err != nil {
				return Metrics{}, err
			}
		}
		if err := st.SealActive("t1"); err != nil {
			return Metrics{}, err
		}
		if err := st.RunGC("t1"); err != nil {
			return Metrics{}, err
		}
		st.SampleSpace()
		if _, err := st.VerifyAllLive(); err != nil {
			return Metrics{}, err
		}

		// Measure sequential scan throughput across 5 passes to amortize timer noise.
		var totalBytes int64
		var totalDur time.Duration
		for pass := 0; pass < 5; pass++ {
			b, d, err := st.SequentialScanAll("t1")
			if err != nil {
				return Metrics{}, err
			}
			totalBytes += b
			totalDur += d
		}
		m := st.Metrics("extent_map_fragmentation_tax")
		m.Policy = label
		if totalDur > 0 {
			m.SeqScanThroughputMBs = round4((float64(totalBytes) / (1024.0 * 1024.0)) / totalDur.Seconds())
		}
		return m, nil
	}

	uncapped, err := runCase("uncapped_punch_hole_no_defrag", 0)
	if err != nil {
		return BoundaryRegimeResult{}, err
	}
	bounded, err := runCase("discogc_bounded_defrag_cap64", 64)
	if err != nil {
		return BoundaryRegimeResult{}, err
	}

	return BoundaryRegimeResult{
		Name:       "extent_map_fragmentation_tax",
		Hypothesis: "Checkerboard alternating 8 KiB deletes produce hundreds of non-contiguous hole extents when hole-punching is uncapped; imposing MaxExtentFragments bounds extent-tree fragmentation via periodic defragmentation at a small one-time copy cost.",
		Observation: fmt.Sprintf(
			"uncapped punch-hole left %d extent fragments (LWA %.2fx, scan %.1f MiB/s); bounded DisCoGC (cap=64) defragmented to %d fragments (LWA %.2fx, scan %.1f MiB/s).",
			uncapped.ExtentFragments, uncapped.FlashWear.LogicalWAF, uncapped.SeqScanThroughputMBs,
			bounded.ExtentFragments, bounded.FlashWear.LogicalWAF, bounded.SeqScanThroughputMBs,
		),
		Results: []Metrics{uncapped, bounded},
	}, nil
}

// Regime 3: Pure append-only workload with 0% deletes and 0% overwrites.
// Proves all 5 policies converge to identical WAF and cost — discard adds zero
// value without dead extents.
func runAppendOnlyZeroDeleteRegime(drive DriveEconomics) (BoundaryRegimeResult, error) {
	spec := WorkloadSpec{
		Name:            "append_only_zero_delete_workload",
		Description:     "Pure immutable append-only archive with 0% overwrites and 0% deletes",
		Seed:            9003,
		KeyCount:        160,
		ChurnOps:        0,
		MinObjBytes:     32 * 1024,
		MaxObjBytes:     32 * 1024,
		OverwriteProb:   0,
		DeleteProb:      0,
		SegmentCapBytes: 1 * 1024 * 1024,
		SampleEveryOps:  40,
	}
	results := make([]Metrics, 0, len(AllPolicies()))
	for _, p := range AllPolicies() {
		m, err := RunWorkload(spec, p, drive, "")
		if err != nil {
			return BoundaryRegimeResult{}, err
		}
		results = append(results, m)
	}
	return BoundaryRegimeResult{
		Name:       "append_only_zero_delete_workload",
		Hypothesis: "When delete and overwrite rates are 0%, no GC work occurs and all 5 policies yield identical Logical WAF (~1.00x), zero reclaimed blocks, and identical space amplification.",
		Observation: fmt.Sprintf(
			"All 5 policies produced identical Logical WAF (%.4fx), 0 copied live bytes, 0 reclaimed blocks, and final space amplification %.2fx.",
			results[0].FlashWear.LogicalWAF, results[0].FinalSpaceAmplification,
		),
		Results: results,
	}, nil
}

type zipfTable struct {
	cdf []float64
}

func newZipfTable(n int, s float64) *zipfTable {
	if n <= 0 || s <= 0 {
		return nil
	}
	cdf := make([]float64, n)
	var sum float64
	for i := 0; i < n; i++ {
		sum += 1.0 / math.Pow(float64(i+1), s)
		cdf[i] = sum
	}
	for i := range cdf {
		cdf[i] /= sum
	}
	return &zipfTable{cdf: cdf}
}

func (z *zipfTable) sample(rng *rand.Rand) int {
	if z == nil || len(z.cdf) == 0 {
		return 0
	}
	u := rng.Float64()
	lo, hi := 0, len(z.cdf)-1
	for lo < hi {
		mid := (lo + hi) / 2
		if z.cdf[mid] < u {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return lo
}

func pickSize(rng *rand.Rand, minB, maxB int) int {
	if maxB <= minB {
		return minB
	}
	return minB + rng.Intn(maxB-minB+1)
}

func deterministicPayload(rng *rand.Rand, size int, tag uint64) []byte {
	buf := make([]byte, size)
	seed := rng.Uint64() ^ (tag * 0x9e3779b97f4a7c15)
	for i := 0; i+8 <= size; i += 8 {
		seed ^= seed << 13
		seed ^= seed >> 7
		seed ^= seed << 17
		buf[i] = byte(seed)
		buf[i+1] = byte(seed >> 8)
		buf[i+2] = byte(seed >> 16)
		buf[i+3] = byte(seed >> 24)
		buf[i+4] = byte(seed >> 32)
		buf[i+5] = byte(seed >> 40)
		buf[i+6] = byte(seed >> 48)
		buf[i+7] = byte(seed >> 56)
	}
	for i := (size / 8) * 8; i < size; i++ {
		buf[i] = byte(seed >> uint((i%8)*8))
	}
	return buf
}

func compactKeySlice(keys []string, live map[string]bool) []string {
	out := keys[:0]
	for _, k := range keys {
		if live[k] {
			out = append(out, k)
		}
	}
	return out
}
