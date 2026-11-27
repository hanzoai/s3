// Package discogc implements a clean-room, file-backed log-structured storage
// engine and quantitative garbage-collection economics evaluator comparing
// direct invalid-extent discard (Linux FALLOC_FL_PUNCH_HOLE / NVMe TRIM model)
// against traditional copy-compaction policies.
//
// Reference: Runhua Bian et al., "Discard-Based Garbage Collection for
// Distributed Log-Structured Storage Systems in ByteDance" (DisCoGC),
// USENIX FAST '26, https://www.usenix.org/conference/fast26/presentation/bian.
package discogc

import "time"

// Policy selects the garbage-collection strategy evaluated on the segment store.
type Policy string

const (
	// PolicyNoGC appends objects without ever discarding or compacting dead space.
	// Achieves ~1.00x logical WAF at the expense of unbounded space amplification.
	PolicyNoGC Policy = "no_gc_append_only"

	// PolicyGreedyCompaction selects sealed segments with the highest garbage ratio
	// and rewrites all live objects into new segments, matching traditional
	// SeaweedFS-style full-volume vacuum compaction.
	PolicyGreedyCompaction Policy = "greedy_copy_compaction"

	// PolicyCostBenefitCompaction selects segments using the Rosenblum-Ousterhout
	// LFS cost-benefit score: (dead_ratio * segment_age) / (1 + live_ratio).
	PolicyCostBenefitCompaction Policy = "cost_benefit_copy_compaction"

	// PolicyAdaptiveDiscard implements a DisCoGC-inspired hybrid strategy:
	// directly reclaims block-aligned dead extents via OS hole-punching with
	// zero live-data copying when dead extents meet alignment/size thresholds
	// and fragmentation stays within budget, falling back to selective
	// copy-compaction only when sub-block fragmentation, extent-map count, or
	// near-empty segment thresholds are exceeded.
	PolicyAdaptiveDiscard Policy = "discogc_adaptive_discard"

	// PolicyBlindSubblockDiscard is a documented negative control that issues
	// raw hole-punch syscalls on every deleted object regardless of size or
	// block alignment and never compacts, demonstrating 0-byte physical
	// reclamation plus wasted syscall overhead on sub-block deletes.
	PolicyBlindSubblockDiscard Policy = "blind_subblock_discard_negative"
)

// AllPolicies returns all five matched policies in canonical evaluation order.
func AllPolicies() []Policy {
	return []Policy{
		PolicyNoGC,
		PolicyGreedyCompaction,
		PolicyCostBenefitCompaction,
		PolicyAdaptiveDiscard,
		PolicyBlindSubblockDiscard,
	}
}

// Config configures a Store instance.
type Config struct {
	// Dir is the directory where real segment files are created.
	// If empty, a temporary directory is created and removed on Close.
	Dir string

	// Policy selects the reclamation strategy.
	Policy Policy

	// SegmentCapacityBytes is the target maximum size of an active segment
	// before it is sealed and a new segment is opened. Default: 2 MiB.
	SegmentCapacityBytes int64

	// DiscardGranularity is the filesystem block / discard alignment unit in
	// bytes (typically 4096 B on ext4/xfs/tmpfs or 65536 B for coarse zones).
	// Default: 4096.
	DiscardGranularity int64

	// MinDiscardExtent is the minimum coalesced aligned dead extent size in
	// bytes required before PolicyAdaptiveDiscard issues a PunchHole call.
	// Default: equal to DiscardGranularity.
	MinDiscardExtent int64

	// GarbageTriggerRatio is the dead-bytes fraction (dead / segment_size)
	// that triggers garbage collection on a segment (analogous to Hanzo S3
	// Volume.garbageLevel threshold, default 0.30).
	GarbageTriggerRatio float64

	// SubblockFragFallbackRatio is the fraction of a segment occupied by
	// un-reclaimable sub-block dead bytes above which PolicyAdaptiveDiscard
	// falls back to selective copy-compaction. Default: 0.25.
	SubblockFragFallbackRatio float64

	// MaxExtentFragments is the maximum number of punched-hole extent fragments
	// permitted in a single segment before PolicyAdaptiveDiscard defragments
	// via copy-compaction. <= 0 disables the extent-fragment cap. Default: 128.
	MaxExtentFragments int

	// SparseEmptyLiveThreshold is the live-byte ratio below which a sealed
	// sparse segment is compacted away by PolicyAdaptiveDiscard to release file
	// descriptors and index metadata. Default: 0.05.
	SparseEmptyLiveThreshold float64

	// AlignObjectsAtOrAboveGranularity aligns objects whose framed size is at
	// least DiscardGranularity onto DiscardGranularity boundaries so dead
	// extents align cleanly with physical filesystem blocks. Smaller objects
	// remain tightly packed (8-byte aligned) as in SeaweedFS needles.
	AlignObjectsAtOrAboveGranularity bool

	// SyncOnWrite calls fsync on segment files after writes/discards when true.
	SyncOnWrite bool

	// Drive configures the SSD wear and $/TB-month economic model.
	Drive DriveEconomics
}

// DefaultConfig returns sensible defaults for a realistic local evaluation.
func DefaultConfig(policy Policy) Config {
	return Config{
		Policy:                           policy,
		SegmentCapacityBytes:             2 * 1024 * 1024,
		DiscardGranularity:               4096,
		MinDiscardExtent:                 4096,
		GarbageTriggerRatio:              0.30,
		SubblockFragFallbackRatio:        0.30,
		MaxExtentFragments:               128,
		SparseEmptyLiveThreshold:         0.05,
		AlignObjectsAtOrAboveGranularity: true,
		SyncOnWrite:                      false,
		Drive:                            DefaultEnterpriseNVMe(),
	}
}

// DriveEconomics parameterizes an enterprise NVMe TLC SSD for wear and TCO
// modeling. All inputs are explicit so callers can reproduce or re-parameterize
// economics for alternative SSD tiers (QLC, ZNS, HDD, etc.).
type DriveEconomics struct {
	// Name identifies the SSD profile in reports.
	Name string `json:"name"`

	// RawCapacityTB is the drive raw capacity in decimal terabytes (e.g. 7.68).
	RawCapacityTB float64 `json:"raw_capacity_tb"`

	// DriveCostUSD is the capital purchase price of one SSD in USD.
	DriveCostUSD float64 `json:"drive_cost_usd"`

	// RatedDWPD is the manufacturer endurance rating in drive writes per day
	// over WarrantyYears.
	RatedDWPD float64 `json:"rated_dwpd"`

	// WarrantyYears is the nominal service/amortization life in years (e.g. 5).
	WarrantyYears float64 `json:"warranty_years"`

	// ProgramEraseCycles is the nominal NAND flash P/E cycle endurance (e.g. 3000 for TLC).
	ProgramEraseCycles float64 `json:"program_erase_cycles"`

	// DeviceInternalWAFCompaction is the SSD FTL physical write amplification (PWA)
	// under copy-compaction without fine-grained discard/TRIM hints.
	DeviceInternalWAFCompaction float64 `json:"device_internal_waf_compaction"`

	// DeviceInternalWAFDiscard is the SSD FTL physical write amplification (PWA)
	// when host discards / TRIMs stale extents proactively, reducing SSD FTL GC.
	DeviceInternalWAFDiscard float64 `json:"device_internal_waf_discard"`

	// DeviceInternalWAFNoGC is the SSD FTL PWA under pure append-only writes.
	DeviceInternalWAFNoGC float64 `json:"device_internal_waf_no_gc"`

	// OverProvisioningFraction is baseline cluster/filesystem OP overhead (e.g. 0.07).
	OverProvisioningFraction float64 `json:"over_provisioning_fraction"`

	// WorkloadLogicalTBPerMonthPerTB is the steady-state logical write rate
	// (TB written per month per 1 TB of live stored data). E.g., 24.0 TB/month/TB
	// corresponds to ~0.80 logical DWPD on a high-churn S3/checkpoint tier.
	WorkloadLogicalTBPerMonthPerTB float64 `json:"workload_logical_tb_per_month_per_tb"`

	// CompactionCPUCostPerGBUSD prices CPU, memory bandwidth, and I/O channel
	// contention per GiB of live data copied during compaction.
	CompactionCPUCostPerGBUSD float64 `json:"compaction_cpu_cost_per_gb_usd"`

	// DiscardSyscallCostPerMillionUSD prices CPU/metadata overhead per 1M
	// discard/TRIM operations.
	DiscardSyscallCostPerMillionUSD float64 `json:"discard_syscall_cost_per_million_usd"`
}

// DefaultEnterpriseNVMe returns parameters modeled on a mainstream 7.68 TB
// enterprise PCIe Gen4 TLC NVMe SSD (1.0 DWPD, 5-year warranty, $150/TB capex).
func DefaultEnterpriseNVMe() DriveEconomics {
	return DriveEconomics{
		Name:                            "7.68TB Enterprise TLC NVMe (1.0 DWPD, 5yr)",
		RawCapacityTB:                   7.68,
		DriveCostUSD:                    1152.0, // $150/TB
		RatedDWPD:                       1.0,
		WarrantyYears:                   5.0,
		ProgramEraseCycles:              3000.0,
		DeviceInternalWAFCompaction:     1.25,
		DeviceInternalWAFDiscard:        1.08,
		DeviceInternalWAFNoGC:           1.10,
		OverProvisioningFraction:        0.07,
		WorkloadLogicalTBPerMonthPerTB:  24.0, // ~0.80 logical DWPD on active churn tier
		CompactionCPUCostPerGBUSD:       0.00025,
		DiscardSyscallCostPerMillionUSD: 0.002,
	}
}

// Extent represents a byte interval [Start, End) inside a segment file.
type Extent struct {
	Start int64
	End   int64
}

// Len returns the extent length in bytes.
func (e Extent) Len() int64 {
	if e.End <= e.Start {
		return 0
	}
	return e.End - e.Start
}

// ObjectMeta describes a stored object's on-disk location and integrity digest.
type ObjectMeta struct {
	TenantID   string
	Key        string
	SegmentID  uint64
	Offset     int64
	HeaderLen  int64
	PayloadLen int64
	TotalLen   int64
	SHA256     [32]byte
	Seq        uint64
}

// ProcIOSnapshot records Linux /proc/self/io process I/O counters.
type ProcIOSnapshot struct {
	Supported           bool  `json:"supported"`
	RChar               int64 `json:"rchar"`
	WChar               int64 `json:"wchar"`
	ReadBytes           int64 `json:"read_bytes"`
	WriteBytes          int64 `json:"write_bytes"`
	CancelledWriteBytes int64 `json:"cancelled_write_bytes"`
}

// FlashWearMetrics captures derived flash endurance and lifetime metrics.
type FlashWearMetrics struct {
	LogicalWAF              float64 `json:"logical_waf"`
	DeviceInternalPWA       float64 `json:"device_internal_pwa"`
	EffectivePhysicalWAF    float64 `json:"effective_physical_waf"`
	TBWPerLogicalTB         float64 `json:"tbw_per_logical_tb"`
	PECyclesPerLogicalTB    float64 `json:"pe_cycles_per_logical_tb"`
	EffectiveDWPD           float64 `json:"effective_dwpd"`
	RatedDWPD               float64 `json:"rated_dwpd"`
	EstimatedDriveLifeYears float64 `json:"estimated_drive_life_years"`
	WarrantyYears           float64 `json:"warranty_years"`
}

// CostBreakdown decomposes monthly storage cost per live TB ($/TB-month).
type CostBreakdown struct {
	CapacityAmortizationUSDPerTBMonth float64 `json:"capacity_amortization_usd_per_tb_month"`
	WearReplacementUSDPerTBMonth      float64 `json:"wear_replacement_usd_per_tb_month"`
	CompactionIOTaxUSDPerTBMonth      float64 `json:"compaction_io_tax_usd_per_tb_month"`
	TotalUSDPerTBMonth                float64 `json:"total_usd_per_tb_month"`
}

// Metrics summarizes measured I/O, physical disk blocks, GC counters, and
// derived storage economics for a run.
type Metrics struct {
	Policy                  Policy           `json:"policy"`
	Workload                string           `json:"workload,omitempty"`
	UserBytesWritten        int64            `json:"user_bytes_written"`
	GCCopiedLiveBytes       int64            `json:"gc_copied_live_bytes"`
	MetadataBytesWritten    int64            `json:"metadata_bytes_written"`
	LiveObjects             int              `json:"live_objects"`
	LiveBytes               int64            `json:"live_bytes"`
	LogicalFileBytes        int64            `json:"logical_file_bytes"`
	PhysicalDiskBytes       int64            `json:"physical_disk_bytes"`
	PeakPhysicalDiskBytes   int64            `json:"peak_physical_disk_bytes"`
	AvgPhysicalDiskBytes    float64          `json:"avg_physical_disk_bytes"`
	AvgLiveBytes            float64          `json:"avg_live_bytes"`
	AvgSpaceAmplification   float64          `json:"avg_space_amplification"`
	FinalSpaceAmplification float64          `json:"final_space_amplification"`
	ActiveSegments          int              `json:"active_segments"`
	PunchHoleSupported      bool             `json:"punch_hole_supported"`
	PunchHoleCalls          int64            `json:"punch_hole_calls"`
	PunchHoleRequestedBytes int64            `json:"punch_hole_requested_bytes"`
	OSBlocksReclaimedBytes  int64            `json:"os_blocks_reclaimed_bytes"`
	DiscardCycles           int64            `json:"discard_cycles"`
	CompactionCycles        int64            `json:"compaction_cycles"`
	ExtentFragments         int              `json:"extent_fragments"`
	FSBlockSizeBytes        int64            `json:"fs_block_size_bytes"`
	ProcIODelta             ProcIOSnapshot   `json:"proc_io_delta"`
	WallDuration            time.Duration    `json:"-"`
	WallDurationMs          float64          `json:"wall_duration_ms"`
	SeqScanThroughputMBs    float64          `json:"seq_scan_throughput_mbs,omitempty"`
	FlashWear               FlashWearMetrics `json:"flash_wear"`
	Cost                    CostBreakdown    `json:"cost_usd_per_tb_month"`
}
