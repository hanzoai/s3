package acos

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"
)

// PolicyName identifies one of the five matched placement and recovery policies.
type PolicyName string

const (
	PolicySingleRegion3xLocal             PolicyName = "single_region_3x_local"
	PolicyNaiveGlobalFullReplication      PolicyName = "naive_global_full_replication"
	PolicyStaticCrossRegionEC             PolicyName = "static_cross_region_ec"
	PolicyACOSGeoAwareAdaptive            PolicyName = "acos_geo_aware_adaptive"
	PolicyOverAggressiveMigrationNegative PolicyName = "over_aggressive_migration_negative"
)

// AllPolicies returns the five matched policies in deterministic evaluation order.
func AllPolicies() []PolicyName {
	return []PolicyName{
		PolicySingleRegion3xLocal,
		PolicyNaiveGlobalFullReplication,
		PolicyStaticCrossRegionEC,
		PolicyACOSGeoAwareAdaptive,
		PolicyOverAggressiveMigrationNegative,
	}
}

// ACOSConfig tunes adaptive locality migration, hysteresis, and erasure striping.
type ACOSConfig struct {
	// WindowSize is the bounded ring-buffer length tracking recent read origins per object.
	WindowSize int `json:"window_size"`
	// MinReadsBeforeMigrate prevents migration of cold or write-once objects.
	MinReadsBeforeMigrate int `json:"min_reads_before_migrate"`
	// AffinityThreshold is the minimum fraction [0,1] of recent reads from a single
	// remote region required before triggering primary migration.
	AffinityThreshold float64 `json:"affinity_threshold"`
	// MigrationCooldownOps is the minimum per-object operation gap between migrations.
	MigrationCooldownOps int `json:"migration_cooldown_ops"`
	// LocalStripeK and LocalStripeM define intra-region AZ erasure coding (default 2+1 = 1.5x).
	LocalStripeK int `json:"local_stripe_k"`
	LocalStripeM int `json:"local_stripe_m"`
	// StaticECK and StaticECM define cross-region static EC (default 4+2 = 1.5x across 3 regions).
	StaticECK int `json:"static_ec_k"`
	StaticECM int `json:"static_ec_m"`
}

// DefaultACOSConfig returns the recommended production-bounded ACOS configuration.
func DefaultACOSConfig() ACOSConfig {
	return ACOSConfig{
		WindowSize:            8,
		MinReadsBeforeMigrate: 5,
		AffinityThreshold:     0.70,
		MigrationCooldownOps:  8,
		LocalStripeK:          2,
		LocalStripeM:          1,
		StaticECK:             4,
		StaticECM:             2,
	}
}

// ShardKind distinguishes full replicas from erasure-coded shards.
type ShardKind string

const (
	ShardKindFull    ShardKind = "full"
	ShardKindLocalEC ShardKind = "local_ec"
	ShardKindGeoEC   ShardKind = "geo_ec"
)

// StoredShard represents an actual chunk stored in a specific region and AZ.
type StoredShard struct {
	Region   RegionID
	AZ       AZID
	Kind     ShardKind
	Index    int
	Data     []byte
	Checksum string
}

type objectRecord struct {
	TenantID      string
	Key           string
	Size          int
	Digest        string
	PrimaryRegion RegionID
	StandbyRegion RegionID
	Shards        []StoredShard

	// Access-locality tracking state (used by ACOS and negative control).
	RecentReaders     []RegionID
	OpCounter         int
	LastMigrationOp   int
	MigrationCount    int
	MetadataSizeBytes int
}

// OpKind identifies the category of an engine operation.
type OpKind string

const (
	OpKindPut OpKind = "put"
	OpKindGet OpKind = "get"
)

// OpResult reports the outcome and telemetry of a single Put or Get call.
type OpResult struct {
	Kind                   OpKind
	LatencyMs              float64
	FailoverRecoveryMs     float64
	CrossRegionEgressBytes int64
	CrossRegionEgressUSD   float64
	ServedFromRegion       RegionID
	LocalAZRecovery        bool
	CrossWANRecovery       bool
	Migrated               bool
	SovereigntyBlocked     bool
	DecisionTimeNs         int64
}

// MetricsSnapshot captures deterministic metrics and real decision throughput for a policy run.
type MetricsSnapshot struct {
	Policy PolicyName `json:"policy"`

	TotalOps     int64 `json:"total_ops"`
	SucceededOps int64 `json:"succeeded_ops"`
	FailedOps    int64 `json:"failed_ops"`
	PutOps       int64 `json:"put_ops"`
	GetOps       int64 `json:"get_ops"`
	GetSucceeded int64 `json:"get_succeeded"`
	GetFailed    int64 `json:"get_failed"`
	ObjectCount  int   `json:"object_count"`

	NormalGetOps             int64 `json:"normal_get_ops"`
	NormalGetSucceeded       int64 `json:"normal_get_succeeded"`
	AZOutageGetOps           int64 `json:"az_outage_get_ops"`
	AZOutageGetSucceeded     int64 `json:"az_outage_get_succeeded"`
	RegionBlackoutGetOps     int64 `json:"region_blackout_get_ops"`
	RegionBlackoutGetSuccess int64 `json:"region_blackout_get_succeeded"`
	AffectedRegionGetOps     int64 `json:"affected_region_get_ops"`
	AffectedRegionGetSuccess int64 `json:"affected_region_get_succeeded"`

	ReadP50Ms             float64 `json:"read_p50_ms"`
	ReadP95Ms             float64 `json:"read_p95_ms"`
	ReadP99Ms             float64 `json:"read_p99_ms"`
	WriteP50Ms            float64 `json:"write_p50_ms"`
	WriteP95Ms            float64 `json:"write_p95_ms"`
	WriteP99Ms            float64 `json:"write_p99_ms"`
	FailoverRecoveryP50Ms float64 `json:"failover_recovery_p50_ms"`
	FailoverRecoveryP99Ms float64 `json:"failover_recovery_p99_ms"`
	RTOMs                 float64 `json:"rto_ms"`

	OverallAvailabilityPct        float64 `json:"overall_availability_pct"`
	AZOutageAvailabilityPct       float64 `json:"az_outage_availability_pct"`
	RegionBlackoutAvailabilityPct float64 `json:"region_blackout_availability_pct"`
	AffectedRegionAvailabilityPct float64 `json:"affected_region_availability_pct"`

	StoredLogicalBytes          int64   `json:"stored_logical_bytes"`
	StoredPhysicalBytes         int64   `json:"stored_physical_bytes"`
	StorageOverheadRatio        float64 `json:"storage_overhead_ratio"`
	StorageCostUSDPerTBMonth    float64 `json:"storage_cost_usd_per_tb_month"`
	TotalServedReadBytes        int64   `json:"total_served_read_bytes"`
	TotalCrossRegionEgressBytes int64   `json:"total_cross_region_egress_bytes"`
	WriteEgressBytes            int64   `json:"write_egress_bytes"`
	ReadEgressBytes             int64   `json:"read_egress_bytes"`
	RebuildEgressBytes          int64   `json:"rebuild_egress_bytes"`
	MigrationEgressBytes        int64   `json:"migration_egress_bytes"`
	EgressCostUSDTotal          float64 `json:"egress_cost_usd_total"`
	EgressCostUSDPerTBServed    float64 `json:"egress_cost_usd_per_tb_served"`
	TotalMonthlyTCOUSDPerTB     float64 `json:"total_monthly_tco_usd_per_tb"`

	MigrationCount        int64   `json:"migration_count"`
	LocalAZRecoveries     int64   `json:"local_az_recoveries"`
	CrossWANRecoveries    int64   `json:"cross_wan_recoveries"`
	MetadataOverheadBytes int64   `json:"metadata_overhead_bytes"`
	DecisionCount         int64   `json:"decision_count"`
	DecisionTotalNs       int64   `json:"decision_total_ns"`
	DecisionThroughputOps float64 `json:"decision_throughput_ops_per_sec"`

	SovereigntyViolations int64 `json:"sovereignty_violations"`
	SovereigntyBlockedOps int64 `json:"sovereignty_blocked_ops"`
	DataCorruptions       int64 `json:"data_corruptions"`
}

// Engine is a thread-safe multi-region object placement, failover, and recovery coordinator.
type Engine struct {
	mu       sync.RWMutex
	policy   PolicyName
	cfg      ACOSConfig
	topology *Topology
	tenants  map[string]TenantPolicy
	objects  map[string]*objectRecord

	// Metrics accumulators.
	readLatenciesMs       []float64
	writeLatenciesMs      []float64
	failoverLatenciesMs   []float64
	rtoSamplesMs          []float64
	totalOps              int64
	succeededOps          int64
	failedOps             int64
	putOps                int64
	getOps                int64
	getSucceeded          int64
	getFailed             int64
	normalGetOps          int64
	normalGetSucceeded    int64
	azOutageGetOps        int64
	azOutageGetSucceeded  int64
	regionOutageGetOps    int64
	regionOutageGetSucc   int64
	affectedRegionGetOps  int64
	affectedRegionGetSucc int64

	totalServedReadBytes int64
	writeEgressBytes     int64
	readEgressBytes      int64
	rebuildEgressBytes   int64
	migrationEgressBytes int64
	egressCostUSD        float64

	migrationCount        int64
	localAZRecoveries     int64
	crossWANRecoveries    int64
	decisionCount         int64
	decisionTotalNs       int64
	sovereigntyViolations int64
	sovereigntyBlockedOps int64
	dataCorruptions       int64
}

// NewEngine creates an isolated Engine instance for the given policy and topology.
func NewEngine(policy PolicyName, topo *Topology, cfg ACOSConfig) *Engine {
	if topo == nil {
		topo = DefaultTopology()
	}
	if cfg.WindowSize <= 0 {
		cfg = DefaultACOSConfig()
	}
	return &Engine{
		policy:   policy,
		cfg:      cfg,
		topology: topo.Clone(),
		tenants:  make(map[string]TenantPolicy),
		objects:  make(map[string]*objectRecord),
	}
}

// RegisterTenant configures sovereignty and fault-domain policy for a tenant.
func (e *Engine) RegisterTenant(tp TenantPolicy) error {
	if err := ValidateTenantID(tp.TenantID); err != nil {
		return err
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.tenants[tp.TenantID] = tp
	return nil
}

func (e *Engine) tenantPolicyLocked(tenantID string) TenantPolicy {
	if tp, ok := e.tenants[tenantID]; ok {
		return tp
	}
	return DefaultGlobalTenantPolicy(tenantID)
}

func objectCompositeKey(tenantID, key string) string {
	return tenantID + "\x00" + key
}

// SetAZDown injects or clears a single availability zone failure.
func (e *Engine) SetAZDown(az AZID, down bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if down {
		e.topology.AZDown[az] = true
	} else {
		delete(e.topology.AZDown, az)
	}
}

// SetRegionDown injects or clears a full-region blackout.
func (e *Engine) SetRegionDown(r RegionID, down bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if down {
		e.topology.RegionDown[r] = true
	} else {
		delete(e.topology.RegionDown, r)
	}
}

// HealAll restores all availability zones and regions to healthy state.
func (e *Engine) HealAll() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.topology.AZDown = make(map[AZID]bool)
	e.topology.RegionDown = make(map[RegionID]bool)
}

func (e *Engine) hasAnyAZOutageLocked() bool {
	return len(e.topology.AZDown) > 0
}

func (e *Engine) hasAnyRegionOutageLocked() bool {
	return len(e.topology.RegionDown) > 0
}

// Put writes payload for (tenantID, key) originating from originRegion according to the configured policy.
func (e *Engine) Put(tenantID, key string, originRegion RegionID, payload []byte) (OpResult, error) {
	if err := ValidateTenantID(tenantID); err != nil {
		return OpResult{}, err
	}
	if err := ValidateObjectKey(key); err != nil {
		return OpResult{}, err
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	startDecision := time.Now()
	e.totalOps++
	e.putOps++

	spec, ok := e.topology.Regions[originRegion]
	if !ok {
		e.failedOps++
		return OpResult{}, fmt.Errorf("%w: %s", ErrUnknownRegion, originRegion)
	}

	tp := e.tenantPolicyLocked(tenantID)
	if !tp.PermitsRegion(spec) {
		e.failedOps++
		e.sovereigntyBlockedOps++
		return OpResult{Kind: OpKindPut, SovereigntyBlocked: true}, fmt.Errorf("%w: tenant=%s origin=%s", ErrSovereigntyViolation, tenantID, originRegion)
	}

	if !e.topology.IsRegionUp(originRegion) {
		// Try nearest healthy permitted peer region if origin is down.
		peers := e.topology.NearestPeerRegions(originRegion, tp, false)
		if len(peers) == 0 {
			e.failedOps++
			if tp.StrictSovereignty {
				e.sovereigntyBlockedOps++
			}
			return OpResult{Kind: OpKindPut, SovereigntyBlocked: tp.StrictSovereignty}, ErrUnavailable
		}
		originRegion = peers[0]
		spec = e.topology.Regions[originRegion]
	}

	digest := PayloadDigest(payload)
	res := OpResult{Kind: OpKindPut, ServedFromRegion: originRegion}

	var (
		shards        []StoredShard
		standbyRegion RegionID
		writeLatMs    float64
		egressBytes   int64
		egressUSD     float64
		metaBytes     int
		err           error
	)

	switch e.policy {
	case PolicySingleRegion3xLocal:
		shards, writeLatMs, err = e.placeSingleRegion3xLocked(originRegion, payload, digest)
		metaBytes = 128

	case PolicyNaiveGlobalFullReplication:
		shards, writeLatMs, egressBytes, egressUSD, err = e.placeNaiveGlobalLocked(originRegion, tp, payload, digest)
		metaBytes = 192

	case PolicyStaticCrossRegionEC:
		shards, writeLatMs, egressBytes, egressUSD, err = e.placeStaticCrossRegionECLocked(originRegion, tp, payload, digest)
		metaBytes = 224

	case PolicyACOSGeoAwareAdaptive, PolicyOverAggressiveMigrationNegative:
		shards, standbyRegion, writeLatMs, egressBytes, egressUSD, err = e.placeACOSAdaptiveLocked(originRegion, tp, payload, digest)
		// Account honestly for access-locality ring-buffer & cooldown metadata tracking.
		metaBytes = 256 + e.cfg.WindowSize*24

	default:
		err = fmt.Errorf("acos: unsupported policy %q", e.policy)
	}

	decisionNs := time.Since(startDecision).Nanoseconds()
	e.decisionCount++
	e.decisionTotalNs += decisionNs
	res.DecisionTimeNs = decisionNs

	if err != nil {
		e.failedOps++
		return res, err
	}

	// Hard sovereignty invariant check before committing any shard placement.
	for _, sh := range shards {
		if !tp.PermitsRegion(e.topology.Regions[sh.Region]) {
			e.sovereigntyViolations++
			e.failedOps++
			return res, fmt.Errorf("%w: placement attempted region=%s", ErrSovereigntyViolation, sh.Region)
		}
	}

	rec := &objectRecord{
		TenantID:          tenantID,
		Key:               key,
		Size:              len(payload),
		Digest:            digest,
		PrimaryRegion:     originRegion,
		StandbyRegion:     standbyRegion,
		Shards:            shards,
		RecentReaders:     make([]RegionID, 0, e.cfg.WindowSize),
		LastMigrationOp:   -e.cfg.MigrationCooldownOps,
		MetadataSizeBytes: metaBytes,
	}
	e.objects[objectCompositeKey(tenantID, key)] = rec

	e.succeededOps++
	e.writeLatenciesMs = append(e.writeLatenciesMs, writeLatMs)
	e.writeEgressBytes += egressBytes
	e.egressCostUSD += egressUSD

	res.LatencyMs = writeLatMs
	res.CrossRegionEgressBytes = egressBytes
	res.CrossRegionEgressUSD = egressUSD
	return res, nil
}

func (e *Engine) placeSingleRegion3xLocked(origin RegionID, payload []byte, digest string) ([]StoredShard, float64, error) {
	azs := e.topology.HealthyAZs(origin)
	if len(azs) == 0 {
		return nil, 0, ErrUnavailable
	}
	spec := e.topology.Regions[origin]
	shards := make([]StoredShard, 0, 3)
	for i := 0; i < 3; i++ {
		az := azs[i%len(azs)]
		cp := append([]byte(nil), payload...)
		shards = append(shards, StoredShard{
			Region:   origin,
			AZ:       az,
			Kind:     ShardKindFull,
			Index:    i,
			Data:     cp,
			Checksum: digest,
		})
	}
	lat := spec.CrossAZRTTMs + spec.StorageIOBaseMs + e.topology.TransferTimeMs(origin, origin, len(payload))
	return shards, lat, nil
}

func (e *Engine) placeNaiveGlobalLocked(origin RegionID, tp TenantPolicy, payload []byte, digest string) ([]StoredShard, float64, int64, float64, error) {
	permitted := e.topology.PermittedRegions(tp)
	targets := make([]RegionID, 0, len(permitted))
	for _, r := range permitted {
		if e.topology.IsRegionUp(r) && len(e.topology.HealthyAZs(r)) > 0 {
			targets = append(targets, r)
		}
	}
	if len(targets) == 0 {
		return nil, 0, 0, 0, ErrUnavailable
	}

	if len(targets) == 1 {
		shards, lat, err := e.placeSingleRegion3xLocked(targets[0], payload, digest)
		return shards, lat, 0, 0, err
	}

	originSpec := e.topology.Regions[origin]
	shards := make([]StoredShard, 0, len(targets))
	maxLat := originSpec.IntraAZRTTMs + originSpec.StorageIOBaseMs + e.topology.TransferTimeMs(origin, origin, len(payload))
	var egressBytes int64
	var egressUSD float64

	for idx, r := range targets {
		azs := e.topology.HealthyAZs(r)
		cp := append([]byte(nil), payload...)
		shards = append(shards, StoredShard{
			Region:   r,
			AZ:       azs[0],
			Kind:     ShardKindFull,
			Index:    idx,
			Data:     cp,
			Checksum: digest,
		})
		if r != origin {
			rSpec := e.topology.Regions[r]
			hopLat := e.topology.LinkRTTMs(origin, r) + e.topology.TransferTimeMs(origin, r, len(payload)) + rSpec.StorageIOBaseMs
			if hopLat > maxLat {
				maxLat = hopLat
			}
			egressBytes += int64(len(payload))
			egressUSD += (float64(len(payload)) / 1e9) * originSpec.EgressCostPerGB
		}
	}
	return shards, maxLat, egressBytes, egressUSD, nil
}

func (e *Engine) placeStaticCrossRegionECLocked(origin RegionID, tp TenantPolicy, payload []byte, digest string) ([]StoredShard, float64, int64, float64, error) {
	k, m := e.cfg.StaticECK, e.cfg.StaticECM
	encoded, err := EncodeShards(payload, k, m)
	if err != nil {
		return nil, 0, 0, 0, err
	}

	// Select up to 3 distinct healthy tenant-permitted regions: origin + 2 nearest peers.
	peers := e.topology.NearestPeerRegions(origin, tp, false)
	targetRegions := []RegionID{origin}
	for _, p := range peers {
		if len(targetRegions) >= 3 {
			break
		}
		targetRegions = append(targetRegions, p)
	}

	originSpec := e.topology.Regions[origin]
	shards := make([]StoredShard, 0, k+m)
	// Base local compute + I/O latency for RS(4,2) encoding.
	maxLat := originSpec.CrossAZRTTMs + originSpec.StorageIOBaseMs + 0.4
	var egressBytes int64
	var egressUSD float64

	// Distribute 2 shards per region (round-robin across distinct AZs in each region).
	for i, data := range encoded {
		r := targetRegions[(i/2)%len(targetRegions)]
		azs := e.topology.HealthyAZs(r)
		if len(azs) == 0 {
			return nil, 0, 0, 0, ErrUnavailable
		}
		az := azs[i%len(azs)]
		shards = append(shards, StoredShard{
			Region:   r,
			AZ:       az,
			Kind:     ShardKindGeoEC,
			Index:    i,
			Data:     data,
			Checksum: PayloadDigest(data),
		})
		if r != origin {
			egressBytes += int64(len(data))
			egressUSD += (float64(len(data)) / 1e9) * originSpec.EgressCostPerGB
			rSpec := e.topology.Regions[r]
			// Account for 2 shards pipelined per remote region.
			hopLat := e.topology.LinkRTTMs(origin, r) + e.topology.TransferTimeMs(origin, r, 2*len(data)) + rSpec.StorageIOBaseMs + 0.4
			if hopLat > maxLat {
				maxLat = hopLat
			}
		}
	}
	return shards, maxLat, egressBytes, egressUSD, nil
}

func (e *Engine) placeACOSAdaptiveLocked(origin RegionID, tp TenantPolicy, payload []byte, digest string) ([]StoredShard, RegionID, float64, int64, float64, error) {
	k, m := e.cfg.LocalStripeK, e.cfg.LocalStripeM
	localAZs := e.topology.HealthyAZs(origin)
	if len(localAZs) == 0 {
		return nil, "", 0, 0, 0, ErrUnavailable
	}

	localEncoded, err := EncodeShards(payload, k, m)
	if err != nil {
		return nil, "", 0, 0, 0, err
	}

	originSpec := e.topology.Regions[origin]
	shards := make([]StoredShard, 0, k+m+1)
	for i, data := range localEncoded {
		az := localAZs[i%len(localAZs)]
		shards = append(shards, StoredShard{
			Region:   origin,
			AZ:       az,
			Kind:     ShardKindLocalEC,
			Index:    i,
			Data:     data,
			Checksum: PayloadDigest(data),
		})
	}
	writeLat := originSpec.CrossAZRTTMs + originSpec.StorageIOBaseMs + e.topology.TransferTimeMs(origin, origin, len(payload)) + 0.2

	// Select single nearest distinct healthy tenant-permitted fault-domain peer region
	// for synchronous RPO=0 regional standby replica.
	var standbyRegion RegionID
	var egressBytes int64
	var egressUSD float64

	peers := e.topology.NearestPeerRegions(origin, tp, false)
	if len(peers) > 0 {
		standbyRegion = peers[0]
		standbyAZs := e.topology.HealthyAZs(standbyRegion)
		if len(standbyAZs) > 0 {
			cp := append([]byte(nil), payload...)
			shards = append(shards, StoredShard{
				Region:   standbyRegion,
				AZ:       standbyAZs[0],
				Kind:     ShardKindFull,
				Index:    k + m,
				Data:     cp,
				Checksum: digest,
			})
			egressBytes = int64(len(payload))
			egressUSD = (float64(len(payload)) / 1e9) * originSpec.EgressCostPerGB
			standbySpec := e.topology.Regions[standbyRegion]
			peerLat := e.topology.LinkRTTMs(origin, standbyRegion) + e.topology.TransferTimeMs(origin, standbyRegion, len(payload)) + standbySpec.StorageIOBaseMs
			if peerLat > writeLat {
				writeLat = peerLat
			}
		} else {
			standbyRegion = ""
		}
	}

	return shards, standbyRegion, writeLat, egressBytes, egressUSD, nil
}

// Get reads and SHA-256-verifies (tenantID, key) from readerRegion, performing
// local-first AZ recovery, regional blackout failover, or adaptive migration as configured.
func (e *Engine) Get(tenantID, key string, readerRegion RegionID) ([]byte, OpResult, error) {
	if err := ValidateTenantID(tenantID); err != nil {
		return nil, OpResult{}, err
	}
	if err := ValidateObjectKey(key); err != nil {
		return nil, OpResult{}, err
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	startDecision := time.Now()
	e.totalOps++
	e.getOps++

	if _, ok := e.topology.Regions[readerRegion]; !ok {
		e.failedOps++
		e.getFailed++
		return nil, OpResult{Kind: OpKindGet}, fmt.Errorf("%w: %s", ErrUnknownRegion, readerRegion)
	}

	rec, ok := e.objects[objectCompositeKey(tenantID, key)]
	if !ok {
		e.failedOps++
		e.getFailed++
		return nil, OpResult{Kind: OpKindGet}, ErrNotFound
	}

	tp := e.tenantPolicyLocked(tenantID)
	inRegionOutage := e.hasAnyRegionOutageLocked()
	inAZOutage := !inRegionOutage && e.hasAnyAZOutageLocked()
	affectedByRegionOutage := e.topology.RegionDown[rec.PrimaryRegion]

	switch {
	case inRegionOutage:
		e.regionOutageGetOps++
		if affectedByRegionOutage {
			e.affectedRegionGetOps++
		}
	case inAZOutage:
		e.azOutageGetOps++
	default:
		e.normalGetOps++
	}

	res := OpResult{Kind: OpKindGet}
	var (
		payload []byte
		err     error
	)

	switch e.policy {
	case PolicySingleRegion3xLocal, PolicyNaiveGlobalFullReplication:
		payload, res, err = e.getFullReplicaLocked(rec, readerRegion, affectedByRegionOutage)
	case PolicyStaticCrossRegionEC:
		payload, res, err = e.getStaticCrossRegionECLocked(rec, readerRegion, affectedByRegionOutage)
	case PolicyACOSGeoAwareAdaptive:
		payload, res, err = e.getACOSAdaptiveLocked(rec, tp, readerRegion, affectedByRegionOutage, false)
	case PolicyOverAggressiveMigrationNegative:
		payload, res, err = e.getACOSAdaptiveLocked(rec, tp, readerRegion, affectedByRegionOutage, true)
	default:
		err = fmt.Errorf("acos: unsupported policy %q", e.policy)
	}

	decisionNs := time.Since(startDecision).Nanoseconds()
	e.decisionCount++
	e.decisionTotalNs += decisionNs
	res.DecisionTimeNs = decisionNs

	if err != nil {
		e.failedOps++
		e.getFailed++
		if tp.StrictSovereignty && affectedByRegionOutage {
			e.sovereigntyBlockedOps++
			res.SovereigntyBlocked = true
		}
		return nil, res, err
	}

	// Verify cryptographic integrity of returned bytes on every single read.
	if PayloadDigest(payload) != rec.Digest {
		e.dataCorruptions++
		e.failedOps++
		e.getFailed++
		return nil, res, ErrIntegrityMismatch
	}

	e.succeededOps++
	e.getSucceeded++
	e.totalServedReadBytes += int64(len(payload))
	e.readLatenciesMs = append(e.readLatenciesMs, res.LatencyMs)
	if res.FailoverRecoveryMs > 0 {
		e.failoverLatenciesMs = append(e.failoverLatenciesMs, res.FailoverRecoveryMs)
	}
	if res.LocalAZRecovery {
		e.localAZRecoveries++
	}
	if res.CrossWANRecovery {
		e.crossWANRecoveries++
	}

	switch {
	case inRegionOutage:
		e.regionOutageGetSucc++
		if affectedByRegionOutage {
			e.affectedRegionGetSucc++
			// Failover RTO: health-check detection window (~15 ms) + first failover serve.
			e.rtoSamplesMs = append(e.rtoSamplesMs, 15.0+res.FailoverRecoveryMs)
		}
	case inAZOutage:
		e.azOutageGetSucceeded++
	default:
		e.normalGetSucceeded++
	}

	return payload, res, nil
}

func (e *Engine) getFullReplicaLocked(rec *objectRecord, reader RegionID, primaryRegionDown bool) ([]byte, OpResult, error) {
	res := OpResult{Kind: OpKindGet}

	// Collect all healthy full replicas and pick the lowest-RTT one from reader.
	bestIdx := -1
	bestRTT := math.MaxFloat64
	for i, sh := range rec.Shards {
		if sh.Kind != ShardKindFull {
			continue
		}
		if !e.topology.IsAZUp(sh.Region, sh.AZ) {
			continue
		}
		rtt := e.topology.LinkRTTMs(reader, sh.Region)
		if sh.Region == reader {
			rtt = e.topology.Regions[reader].IntraAZRTTMs
		}
		if rtt < bestRTT {
			bestRTT = rtt
			bestIdx = i
		}
	}

	if bestIdx == -1 {
		return nil, res, ErrUnavailable
	}

	chosen := rec.Shards[bestIdx]
	srcSpec := e.topology.Regions[chosen.Region]
	var lat float64
	if chosen.Region == reader {
		lat = srcSpec.IntraAZRTTMs + srcSpec.StorageIOBaseMs + e.topology.TransferTimeMs(reader, reader, len(chosen.Data))
	} else {
		lat = e.topology.LinkRTTMs(reader, chosen.Region) + srcSpec.StorageIOBaseMs + e.topology.TransferTimeMs(chosen.Region, reader, len(chosen.Data))
		egress := int64(len(chosen.Data))
		cost := (float64(egress) / 1e9) * srcSpec.EgressCostPerGB
		res.CrossRegionEgressBytes += egress
		res.CrossRegionEgressUSD += cost
		e.readEgressBytes += egress
		e.egressCostUSD += cost
	}

	// Detect local AZ failover or regional blackout failover for this object.
	objectAZDegraded := false
	for _, sh := range rec.Shards {
		if sh.Region == rec.PrimaryRegion && !e.topology.IsAZUp(sh.Region, sh.AZ) {
			objectAZDegraded = true
			break
		}
	}
	if primaryRegionDown {
		res.CrossWANRecovery = true
		res.FailoverRecoveryMs = lat
	} else if objectAZDegraded && chosen.Region == rec.PrimaryRegion {
		res.LocalAZRecovery = true
		lat += srcSpec.CrossAZRTTMs - srcSpec.IntraAZRTTMs
		res.FailoverRecoveryMs = lat
	}

	res.LatencyMs = lat
	res.ServedFromRegion = chosen.Region
	return append([]byte(nil), chosen.Data...), res, nil
}

func (e *Engine) getStaticCrossRegionECLocked(rec *objectRecord, reader RegionID, primaryRegionDown bool) ([]byte, OpResult, error) {
	res := OpResult{Kind: OpKindGet}
	k, m := e.cfg.StaticECK, e.cfg.StaticECM

	type candidate struct {
		shard StoredShard
		rtt   float64
	}
	avail := make([]candidate, 0, len(rec.Shards))
	objectAZDegraded := false
	for _, sh := range rec.Shards {
		if sh.Kind != ShardKindGeoEC {
			continue
		}
		if !e.topology.IsAZUp(sh.Region, sh.AZ) {
			objectAZDegraded = true
			continue
		}
		rtt := e.topology.LinkRTTMs(reader, sh.Region)
		if sh.Region == reader {
			rtt = e.topology.Regions[reader].CrossAZRTTMs
		}
		avail = append(avail, candidate{shard: sh, rtt: rtt})
	}
	if len(avail) < k {
		return nil, res, ErrUnavailable
	}

	// Sort surviving shards by ascending RTT from reader (tie-broken by shard index).
	sort.Slice(avail, func(i, j int) bool {
		if avail[i].rtt != avail[j].rtt {
			return avail[i].rtt < avail[j].rtt
		}
		return avail[i].shard.Index < avail[j].shard.Index
	})

	selected := avail[:k]
	shardBufs := make([][]byte, k+m)
	needsDecode := false
	maxLat := 0.0
	var readEgress int64
	var rebuildEgress int64
	var egressUSD float64
	remoteRegionsUsed := false

	for _, c := range selected {
		sh := c.shard
		if sh.Index >= k {
			needsDecode = true
		}
		shardBufs[sh.Index] = sh.Data
		srcSpec := e.topology.Regions[sh.Region]
		hopLat := c.rtt + srcSpec.StorageIOBaseMs + e.topology.TransferTimeMs(sh.Region, reader, len(sh.Data))
		if hopLat > maxLat {
			maxLat = hopLat
		}
		if sh.Region != reader {
			remoteRegionsUsed = true
			cost := (float64(len(sh.Data)) / 1e9) * srcSpec.EgressCostPerGB
			egressUSD += cost
			if primaryRegionDown || objectAZDegraded {
				rebuildEgress += int64(len(sh.Data))
			} else {
				readEgress += int64(len(sh.Data))
			}
		}
	}

	payload, err := ReconstructPayload(shardBufs, rec.Size, k, m)
	if err != nil {
		return nil, res, err
	}

	if needsDecode {
		// RS(4,2) GF(2^8) reconstruction compute overhead.
		maxLat += 1.2
	}

	res.LatencyMs = maxLat
	res.ServedFromRegion = selected[0].shard.Region
	res.CrossRegionEgressBytes = readEgress + rebuildEgress
	res.CrossRegionEgressUSD = egressUSD
	e.readEgressBytes += readEgress
	e.rebuildEgressBytes += rebuildEgress
	e.egressCostUSD += egressUSD

	if primaryRegionDown || (objectAZDegraded && remoteRegionsUsed) {
		res.CrossWANRecovery = true
		res.FailoverRecoveryMs = maxLat
	}
	return payload, res, nil
}

func (e *Engine) getACOSAdaptiveLocked(rec *objectRecord, tp TenantPolicy, reader RegionID, primaryRegionDown bool, aggressive bool) ([]byte, OpResult, error) {
	res := OpResult{Kind: OpKindGet}
	k, m := e.cfg.LocalStripeK, e.cfg.LocalStripeM

	rec.OpCounter++
	rec.RecentReaders = append(rec.RecentReaders, reader)
	if len(rec.RecentReaders) > e.cfg.WindowSize {
		rec.RecentReaders = rec.RecentReaders[len(rec.RecentReaders)-e.cfg.WindowSize:]
	}

	// Helper to try serving from local intra-region (2+1) stripe in rec.PrimaryRegion.
	tryLocalStripe := func() ([]byte, float64, bool, bool) {
		if !e.topology.IsRegionUp(rec.PrimaryRegion) {
			return nil, 0, false, false
		}
		shardBufs := make([][]byte, k+m)
		surviving := 0
		usedParity := false
		for _, sh := range rec.Shards {
			if sh.Kind != ShardKindLocalEC || sh.Region != rec.PrimaryRegion {
				continue
			}
			if !e.topology.IsAZUp(sh.Region, sh.AZ) {
				continue
			}
			if sh.Index >= 0 && sh.Index < k+m {
				shardBufs[sh.Index] = sh.Data
				surviving++
				if sh.Index >= k {
					usedParity = true
				}
			}
		}
		if surviving < k {
			return nil, 0, false, false
		}
		// Prefer fast path using k data shards when all k data shards are healthy.
		dataHealthy := true
		for i := 0; i < k; i++ {
			if shardBufs[i] == nil {
				dataHealthy = false
				break
			}
		}
		payload, err := ReconstructPayload(shardBufs, rec.Size, k, m)
		if err != nil {
			return nil, 0, false, false
		}
		pSpec := e.topology.Regions[rec.PrimaryRegion]
		localAZRecovery := !dataHealthy && usedParity
		var lat float64
		if reader == rec.PrimaryRegion {
			if localAZRecovery {
				lat = pSpec.CrossAZRTTMs + pSpec.StorageIOBaseMs + e.topology.TransferTimeMs(reader, reader, rec.Size) + 0.5
			} else {
				lat = pSpec.IntraAZRTTMs + pSpec.StorageIOBaseMs + e.topology.TransferTimeMs(reader, reader, rec.Size)
			}
		} else {
			lat = e.topology.LinkRTTMs(reader, rec.PrimaryRegion) + pSpec.StorageIOBaseMs + e.topology.TransferTimeMs(rec.PrimaryRegion, reader, rec.Size)
			if localAZRecovery {
				lat += pSpec.CrossAZRTTMs + 0.5
			}
		}
		return payload, lat, localAZRecovery, true
	}

	// Helper to try serving from regional full standby replica.
	tryStandby := func() ([]byte, RegionID, float64, bool) {
		for _, sh := range rec.Shards {
			if sh.Kind != ShardKindFull {
				continue
			}
			if !e.topology.IsAZUp(sh.Region, sh.AZ) {
				continue
			}
			sSpec := e.topology.Regions[sh.Region]
			var lat float64
			if reader == sh.Region {
				lat = sSpec.IntraAZRTTMs + sSpec.StorageIOBaseMs + e.topology.TransferTimeMs(reader, reader, len(sh.Data))
			} else {
				lat = e.topology.LinkRTTMs(reader, sh.Region) + sSpec.StorageIOBaseMs + e.topology.TransferTimeMs(sh.Region, reader, len(sh.Data))
			}
			return append([]byte(nil), sh.Data...), sh.Region, lat, true
		}
		return nil, "", 0, false
	}

	var (
		payload          []byte
		servedFrom       RegionID
		latMs            float64
		localAZRecovery  bool
		crossWANRecovery bool
		ok               bool
	)

	// ACOS local-first read routing:
	// 1. If reader IS primary region and local stripe is healthy -> read locally (survives single-AZ loss with 0 WAN egress!).
	// 2. If reader IS standby region and standby is healthy -> read locally from standby (0 WAN egress!).
	// 3. Otherwise choose whichever healthy copy (primary stripe or standby) has lower RTT from reader.
	if reader == rec.PrimaryRegion {
		if p, l, azRec, found := tryLocalStripe(); found {
			payload, servedFrom, latMs, localAZRecovery, ok = p, rec.PrimaryRegion, l, azRec, true
		} else if p, sReg, l, found := tryStandby(); found {
			payload, servedFrom, latMs, crossWANRecovery, ok = p, sReg, l, true, true
		}
	} else if reader == rec.StandbyRegion {
		if p, sReg, l, found := tryStandby(); found {
			payload, servedFrom, latMs, ok = p, sReg, l, true
			if primaryRegionDown {
				crossWANRecovery = true
			}
		} else if p, l, azRec, found := tryLocalStripe(); found {
			payload, servedFrom, latMs, localAZRecovery, ok = p, rec.PrimaryRegion, l, azRec, true
			if !e.topology.IsRegionUp(rec.StandbyRegion) {
				crossWANRecovery = true
			}
		}
	} else {
		pData, pLat, pAZRec, pOK := tryLocalStripe()
		sData, sReg, sLat, sOK := tryStandby()
		switch {
		case pOK && (!sOK || pLat <= sLat):
			payload, servedFrom, latMs, localAZRecovery, ok = pData, rec.PrimaryRegion, pLat, pAZRec, true
		case sOK:
			payload, servedFrom, latMs, ok = sData, sReg, sLat, true
			if primaryRegionDown {
				crossWANRecovery = true
			}
		}
	}

	if !ok {
		return nil, res, ErrUnavailable
	}

	// Account for cross-region read egress if served remotely.
	if servedFrom != reader {
		srcSpec := e.topology.Regions[servedFrom]
		egress := int64(len(payload))
		cost := (float64(egress) / 1e9) * srcSpec.EgressCostPerGB
		res.CrossRegionEgressBytes += egress
		res.CrossRegionEgressUSD += cost
		e.egressCostUSD += cost
		if crossWANRecovery {
			e.rebuildEgressBytes += egress
		} else {
			e.readEgressBytes += egress
		}
	}

	if localAZRecovery {
		res.LocalAZRecovery = true
		res.FailoverRecoveryMs = latMs
	}
	if crossWANRecovery {
		res.CrossWANRecovery = true
		res.FailoverRecoveryMs = latMs
	}

	// Evaluate access-locality migration (strictly respecting tenant sovereignty!).
	shouldMigrate := false
	readerSpec := e.topology.Regions[reader]
	if e.topology.IsRegionUp(reader) && tp.PermitsRegion(readerSpec) {
		if aggressive {
			// Negative control: migrate on every single cross-region read where reader != primary,
			// with zero hysteresis and zero cooldown.
			if reader != rec.PrimaryRegion {
				shouldMigrate = true
			}
		} else {
			// ACOS adaptive migration: only migrate when reader has neither primary nor standby,
			// sustained locality exceeds AffinityThreshold over at least MinReadsBeforeMigrate,
			// cooldown has elapsed, and no full-region blackout is currently disrupting placement.
			if reader != rec.PrimaryRegion && reader != rec.StandbyRegion &&
				!e.hasAnyRegionOutageLocked() &&
				len(rec.RecentReaders) >= e.cfg.MinReadsBeforeMigrate &&
				(rec.OpCounter-rec.LastMigrationOp) >= e.cfg.MigrationCooldownOps {
				hits := 0
				for _, r := range rec.RecentReaders {
					if r == reader {
						hits++
					}
				}
				share := float64(hits) / float64(len(rec.RecentReaders))
				if share >= e.cfg.AffinityThreshold {
					shouldMigrate = true
				}
			}
		}
	}

	if shouldMigrate {
		migLat, migBytes, migUSD, migErr := e.migrateACOSObjectLocked(rec, tp, reader, payload, aggressive)
		if migErr == nil {
			res.Migrated = true
			res.CrossRegionEgressBytes += migBytes
			res.CrossRegionEgressUSD += migUSD
			e.migrationEgressBytes += migBytes
			e.egressCostUSD += migUSD
			e.migrationCount++
			rec.MigrationCount++
			rec.LastMigrationOp = rec.OpCounter
			if aggressive {
				// Aggressive synchronous migration blocks the triggering read path.
				latMs += migLat
			}
		}
	}

	res.LatencyMs = latMs
	res.ServedFromRegion = servedFrom
	return payload, res, nil
}

func (e *Engine) migrateACOSObjectLocked(rec *objectRecord, tp TenantPolicy, targetRegion RegionID, payload []byte, aggressive bool) (float64, int64, float64, error) {
	if !tp.PermitsRegion(e.topology.Regions[targetRegion]) {
		e.sovereigntyBlockedOps++
		return 0, 0, 0, ErrSovereigntyViolation
	}
	oldPrimary := rec.PrimaryRegion
	oldStandby := rec.StandbyRegion

	newShards, newStandby, writeLat, standbyEgress, standbyUSD, err := e.placeACOSAdaptiveLocked(targetRegion, tp, payload, rec.Digest)
	if err != nil {
		return 0, 0, 0, err
	}

	// Hard sovereignty check on migrated shards.
	for _, sh := range newShards {
		if !tp.PermitsRegion(e.topology.Regions[sh.Region]) {
			e.sovereigntyViolations++
			return 0, 0, 0, ErrSovereigntyViolation
		}
	}

	// Account for egress required to transfer the object to the new primary (and new standby if changed).
	var migBytes int64
	var migUSD float64
	if oldPrimary != targetRegion {
		srcSpec := e.topology.Regions[oldPrimary]
		migBytes += int64(len(payload))
		migUSD += (float64(len(payload)) / 1e9) * srcSpec.EgressCostPerGB
	}
	if newStandby != "" && newStandby != oldPrimary && newStandby != oldStandby {
		migBytes += standbyEgress
		migUSD += standbyUSD
	} else if aggressive && newStandby != "" {
		// Aggressive migration re-replicates standby eagerly without reusing prior shards.
		migBytes += standbyEgress
		migUSD += standbyUSD
	}

	rec.PrimaryRegion = targetRegion
	rec.StandbyRegion = newStandby
	rec.Shards = newShards
	return writeLat, migBytes, migUSD, nil
}

// AuditPlacementSovereignty scans every stored shard across all tenants and
// returns the count of any out-of-jurisdiction shard placements (must always be 0).
func (e *Engine) AuditPlacementSovereignty() (violations int, auditedShards int) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	for _, rec := range e.objects {
		tp := e.tenantPolicyLocked(rec.TenantID)
		for _, sh := range rec.Shards {
			auditedShards++
			if !tp.PermitsRegion(e.topology.Regions[sh.Region]) {
				violations++
			}
		}
	}
	return violations, auditedShards
}

// ObjectRegions returns the distinct regions currently holding shards for (tenantID, key).
func (e *Engine) ObjectRegions(tenantID, key string) ([]RegionID, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	rec, ok := e.objects[objectCompositeKey(tenantID, key)]
	if !ok {
		return nil, false
	}
	seen := make(map[RegionID]bool)
	out := make([]RegionID, 0, len(rec.Shards))
	for _, sh := range rec.Shards {
		if !seen[sh.Region] {
			seen[sh.Region] = true
			out = append(out, sh.Region)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out, true
}

// Snapshot returns an immutable metrics summary for the current engine run.
// Monthly TCO per TB normalizes storage cost per TB-month plus modeled active
// egress at an active dataset turnover ratio of 5 TB read/written per TB stored per month.
func (e *Engine) Snapshot() MetricsSnapshot {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var logicalBytes int64
	var physicalBytes int64
	var metadataBytes int64
	weightedStorageCost := 0.0

	for _, rec := range e.objects {
		logicalBytes += int64(rec.Size)
		metadataBytes += int64(rec.MetadataSizeBytes)
		for _, sh := range rec.Shards {
			sz := int64(len(sh.Data))
			physicalBytes += sz
			spec := e.topology.Regions[sh.Region]
			weightedStorageCost += (float64(sz) / 1e12) * spec.StorageCostPerTBMonth
		}
	}

	overhead := 0.0
	storagePerTBMonth := 0.0
	if logicalBytes > 0 {
		overhead = float64(physicalBytes) / float64(logicalBytes)
		storagePerTBMonth = weightedStorageCost / (float64(logicalBytes) / 1e12)
	}

	totalEgressBytes := e.writeEgressBytes + e.readEgressBytes + e.rebuildEgressBytes + e.migrationEgressBytes
	egressPerTBServed := 0.0
	if e.totalServedReadBytes > 0 {
		egressPerTBServed = e.egressCostUSD / (float64(e.totalServedReadBytes) / 1e12)
	}

	// Normalized monthly TCO per TB of active data assuming 5 TB of monthly
	// active read traffic per 1 TB stored (plus amortized write/migration egress).
	const monthlyReadMultiplierPerTB = 5.0
	monthlyTCO := storagePerTBMonth + egressPerTBServed*monthlyReadMultiplierPerTB

	throughput := 0.0
	if e.decisionTotalNs > 0 {
		throughput = float64(e.decisionCount) / (float64(e.decisionTotalNs) / 1e9)
	}

	snap := MetricsSnapshot{
		Policy: e.policy,

		TotalOps:     e.totalOps,
		SucceededOps: e.succeededOps,
		FailedOps:    e.failedOps,
		PutOps:       e.putOps,
		GetOps:       e.getOps,
		GetSucceeded: e.getSucceeded,
		GetFailed:    e.getFailed,
		ObjectCount:  len(e.objects),

		NormalGetOps:             e.normalGetOps,
		NormalGetSucceeded:       e.normalGetSucceeded,
		AZOutageGetOps:           e.azOutageGetOps,
		AZOutageGetSucceeded:     e.azOutageGetSucceeded,
		RegionBlackoutGetOps:     e.regionOutageGetOps,
		RegionBlackoutGetSuccess: e.regionOutageGetSucc,
		AffectedRegionGetOps:     e.affectedRegionGetOps,
		AffectedRegionGetSuccess: e.affectedRegionGetSucc,

		ReadP50Ms:             percentile(e.readLatenciesMs, 50),
		ReadP95Ms:             percentile(e.readLatenciesMs, 95),
		ReadP99Ms:             percentile(e.readLatenciesMs, 99),
		WriteP50Ms:            percentile(e.writeLatenciesMs, 50),
		WriteP95Ms:            percentile(e.writeLatenciesMs, 95),
		WriteP99Ms:            percentile(e.writeLatenciesMs, 99),
		FailoverRecoveryP50Ms: percentile(e.failoverLatenciesMs, 50),
		FailoverRecoveryP99Ms: percentile(e.failoverLatenciesMs, 99),
		RTOMs:                 percentile(e.rtoSamplesMs, 95),

		OverallAvailabilityPct:        pct(e.getSucceeded, e.getOps),
		AZOutageAvailabilityPct:       pct(e.azOutageGetSucceeded, e.azOutageGetOps),
		RegionBlackoutAvailabilityPct: pct(e.regionOutageGetSucc, e.regionOutageGetOps),
		AffectedRegionAvailabilityPct: pct(e.affectedRegionGetSucc, e.affectedRegionGetOps),

		StoredLogicalBytes:          logicalBytes,
		StoredPhysicalBytes:         physicalBytes,
		StorageOverheadRatio:        round2(overhead),
		StorageCostUSDPerTBMonth:    round2(storagePerTBMonth),
		TotalServedReadBytes:        e.totalServedReadBytes,
		TotalCrossRegionEgressBytes: totalEgressBytes,
		WriteEgressBytes:            e.writeEgressBytes,
		ReadEgressBytes:             e.readEgressBytes,
		RebuildEgressBytes:          e.rebuildEgressBytes,
		MigrationEgressBytes:        e.migrationEgressBytes,
		EgressCostUSDTotal:          round6(e.egressCostUSD),
		EgressCostUSDPerTBServed:    round2(egressPerTBServed),
		TotalMonthlyTCOUSDPerTB:     round2(monthlyTCO),

		MigrationCount:        e.migrationCount,
		LocalAZRecoveries:     e.localAZRecoveries,
		CrossWANRecoveries:    e.crossWANRecoveries,
		MetadataOverheadBytes: metadataBytes,
		DecisionCount:         e.decisionCount,
		DecisionTotalNs:       e.decisionTotalNs,
		DecisionThroughputOps: round2(throughput),

		SovereigntyViolations: e.sovereigntyViolations,
		SovereigntyBlockedOps: e.sovereigntyBlockedOps,
		DataCorruptions:       e.dataCorruptions,
	}
	return snap
}

func percentile(vals []float64, p float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	cp := append([]float64(nil), vals...)
	sort.Float64s(cp)
	idx := int(math.Ceil((p/100.0)*float64(len(cp)))) - 1
	if idx < 0 {
		idx = 0
	}
	if idx >= len(cp) {
		idx = len(cp) - 1
	}
	return round2(cp[idx])
}

func pct(num, denom int64) float64 {
	if denom == 0 {
		return 100.0
	}
	return round2((float64(num) / float64(denom)) * 100.0)
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

func round6(v float64) float64 {
	return math.Round(v*1e6) / 1e6
}
