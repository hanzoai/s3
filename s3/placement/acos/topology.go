package acos

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// RegionID uniquely identifies a geographic cloud region.
type RegionID string

// AZID uniquely identifies an availability zone within a region.
type AZID string

// Jurisdiction represents a sovereign regulatory boundary for data residency.
type Jurisdiction string

const (
	JurisdictionGlobal Jurisdiction = "GLOBAL"
	JurisdictionUS     Jurisdiction = "US"
	JurisdictionEU     Jurisdiction = "EU"
	JurisdictionAPAC   Jurisdiction = "APAC"
	JurisdictionSA     Jurisdiction = "SA"
)

const (
	RegionUSEast      RegionID = "us-east"
	RegionUSWest      RegionID = "us-west"
	RegionEUCentral   RegionID = "eu-central"
	RegionEUWest      RegionID = "eu-west"
	RegionAPNortheast RegionID = "ap-northeast"
	RegionSAEast      RegionID = "sa-east"
)

var (
	// ErrInvalidTenantID is returned when a tenant identifier fails validation.
	ErrInvalidTenantID = errors.New("acos: invalid tenant id")
	// ErrInvalidObjectKey is returned when an object key fails sanitization.
	ErrInvalidObjectKey = errors.New("acos: invalid object key")
	// ErrUnknownRegion is returned when referencing an unconfigured region.
	ErrUnknownRegion = errors.New("acos: unknown region")
	// ErrSovereigntyViolation is returned when an operation would place or migrate
	// data outside a tenant's permitted jurisdiction/regions.
	ErrSovereigntyViolation = errors.New("acos: tenant data-residency sovereignty policy forbids target region")
	// ErrUnavailable is returned when active outages exceed surviving fault domains (fail-closed).
	ErrUnavailable = errors.New("acos: object unavailable under current outage state (fail-closed)")
	// ErrNotFound is returned when an object does not exist in the tenant namespace.
	ErrNotFound = errors.New("acos: object not found")
)

// RegionSpec defines cost, latency, jurisdiction, and AZ layout of a region.
type RegionSpec struct {
	ID                    RegionID     `json:"id"`
	Jurisdiction          Jurisdiction `json:"jurisdiction"`
	AZs                   []AZID       `json:"azs"`
	StorageCostPerTBMonth float64      `json:"storage_cost_usd_per_tb_month"`
	EgressCostPerGB       float64      `json:"egress_cost_usd_per_gb"`
	IntraAZRTTMs          float64      `json:"intra_az_rtt_ms"`
	CrossAZRTTMs          float64      `json:"cross_az_rtt_ms"`
	StorageIOBaseMs       float64      `json:"storage_io_base_ms"`
	WANBandwidthMBps      float64      `json:"wan_bandwidth_mbps"`
	LocalBandwidthMBps    float64      `json:"local_bandwidth_mbps"`
}

// TenantPolicy enforces per-tenant sovereignty, jurisdiction pinning, and fault-domain budget.
type TenantPolicy struct {
	TenantID             string         `json:"tenant_id"`
	AllowedJurisdictions []Jurisdiction `json:"allowed_jurisdictions,omitempty"`
	AllowedRegions       []RegionID     `json:"allowed_regions,omitempty"`
	RegionPrefixes       []string       `json:"region_prefixes,omitempty"`
	FaultDomainBudget    int            `json:"fault_domain_budget"`
	StrictSovereignty    bool           `json:"strict_sovereignty"`
}

// DefaultGlobalTenantPolicy permits placement across all configured regions with f=1 durability.
func DefaultGlobalTenantPolicy(tenantID string) TenantPolicy {
	return TenantPolicy{
		TenantID:             tenantID,
		AllowedJurisdictions: []Jurisdiction{JurisdictionGlobal},
		FaultDomainBudget:    1,
		StrictSovereignty:    false,
	}
}

// EUOnlyTenantPolicy strictly pins all replicas, parity shards, standby copies,
// and failover targets inside EU jurisdiction / eu-* regions.
func EUOnlyTenantPolicy(tenantID string) TenantPolicy {
	return TenantPolicy{
		TenantID:             tenantID,
		AllowedJurisdictions: []Jurisdiction{JurisdictionEU},
		RegionPrefixes:       []string{"eu-"},
		FaultDomainBudget:    1,
		StrictSovereignty:    true,
	}
}

// PermitsRegion reports whether r satisfies the tenant's sovereignty constraints.
func (tp TenantPolicy) PermitsRegion(spec RegionSpec) bool {
	if len(tp.AllowedRegions) > 0 {
		matched := false
		for _, allowed := range tp.AllowedRegions {
			if allowed == spec.ID {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}
	if len(tp.RegionPrefixes) > 0 {
		matched := false
		for _, pfx := range tp.RegionPrefixes {
			if strings.HasPrefix(string(spec.ID), pfx) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}
	if len(tp.AllowedJurisdictions) > 0 {
		matched := false
		for _, j := range tp.AllowedJurisdictions {
			if j == JurisdictionGlobal || j == spec.Jurisdiction {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}
	return true
}

// ValidateTenantID rejects empty IDs, traversal patterns, separators, and non-printable bytes.
func ValidateTenantID(tenantID string) error {
	if tenantID == "" || len(tenantID) > 64 {
		return fmt.Errorf("%w: length out of bounds", ErrInvalidTenantID)
	}
	if tenantID == "." || tenantID == ".." || strings.Contains(tenantID, "..") {
		return fmt.Errorf("%w: traversal sequence not allowed", ErrInvalidTenantID)
	}
	for i := 0; i < len(tenantID); i++ {
		c := tenantID[i]
		isAlphaNum := (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
		if !isAlphaNum && c != '-' && c != '_' {
			return fmt.Errorf("%w: invalid character %q", ErrInvalidTenantID, c)
		}
	}
	return nil
}

// ValidateObjectKey rejects empty keys, path traversal, backslashes, NUL, and control characters.
func ValidateObjectKey(key string) error {
	if key == "" || len(key) > 512 {
		return fmt.Errorf("%w: length out of bounds", ErrInvalidObjectKey)
	}
	if strings.HasPrefix(key, "/") || strings.Contains(key, "\\") {
		return fmt.Errorf("%w: leading slash or backslash forbidden", ErrInvalidObjectKey)
	}
	for _, seg := range strings.Split(key, "/") {
		if seg == "" || seg == "." || seg == ".." {
			return fmt.Errorf("%w: empty or traversal segment", ErrInvalidObjectKey)
		}
	}
	for i := 0; i < len(key); i++ {
		c := key[i]
		if c < 0x20 || c == 0x7f {
			return fmt.Errorf("%w: control byte not allowed", ErrInvalidObjectKey)
		}
	}
	return nil
}

// Topology models a multi-region cloud footprint, RTT matrix, pricing, and live fault state.
type Topology struct {
	Regions    map[RegionID]RegionSpec
	Order      []RegionID
	RTTMs      map[RegionID]map[RegionID]float64
	AZDown     map[AZID]bool
	RegionDown map[RegionID]bool
}

// DefaultTopology constructs the canonical 5-region multi-continental topology:
// us-east, us-west, eu-central, ap-northeast, sa-east (3 AZs per region).
//
// Pricing and RTT parameters are modeled reference values representative of
// commercial cloud regions (not production Hanzo billing claims).
func DefaultTopology() *Topology {
	specs := []RegionSpec{
		{
			ID:                    RegionUSEast,
			Jurisdiction:          JurisdictionUS,
			AZs:                   []AZID{"us-east-a", "us-east-b", "us-east-c"},
			StorageCostPerTBMonth: 23.00,
			EgressCostPerGB:       0.020,
			IntraAZRTTMs:          0.6,
			CrossAZRTTMs:          1.8,
			StorageIOBaseMs:       1.2,
			WANBandwidthMBps:      250.0,
			LocalBandwidthMBps:    1200.0,
		},
		{
			ID:                    RegionUSWest,
			Jurisdiction:          JurisdictionUS,
			AZs:                   []AZID{"us-west-a", "us-west-b", "us-west-c"},
			StorageCostPerTBMonth: 23.00,
			EgressCostPerGB:       0.020,
			IntraAZRTTMs:          0.6,
			CrossAZRTTMs:          1.8,
			StorageIOBaseMs:       1.2,
			WANBandwidthMBps:      250.0,
			LocalBandwidthMBps:    1200.0,
		},
		{
			ID:                    RegionEUCentral,
			Jurisdiction:          JurisdictionEU,
			AZs:                   []AZID{"eu-central-a", "eu-central-b", "eu-central-c"},
			StorageCostPerTBMonth: 24.50,
			EgressCostPerGB:       0.020,
			IntraAZRTTMs:          0.6,
			CrossAZRTTMs:          1.9,
			StorageIOBaseMs:       1.3,
			WANBandwidthMBps:      220.0,
			LocalBandwidthMBps:    1200.0,
		},
		{
			ID:                    RegionAPNortheast,
			Jurisdiction:          JurisdictionAPAC,
			AZs:                   []AZID{"ap-northeast-a", "ap-northeast-b", "ap-northeast-c"},
			StorageCostPerTBMonth: 25.00,
			EgressCostPerGB:       0.090,
			IntraAZRTTMs:          0.7,
			CrossAZRTTMs:          2.0,
			StorageIOBaseMs:       1.3,
			WANBandwidthMBps:      200.0,
			LocalBandwidthMBps:    1100.0,
		},
		{
			ID:                    RegionSAEast,
			Jurisdiction:          JurisdictionSA,
			AZs:                   []AZID{"sa-east-a", "sa-east-b", "sa-east-c"},
			StorageCostPerTBMonth: 30.00,
			EgressCostPerGB:       0.110,
			IntraAZRTTMs:          0.8,
			CrossAZRTTMs:          2.2,
			StorageIOBaseMs:       1.5,
			WANBandwidthMBps:      160.0,
			LocalBandwidthMBps:    1000.0,
		},
	}

	t := &Topology{
		Regions:    make(map[RegionID]RegionSpec, len(specs)),
		Order:      make([]RegionID, 0, len(specs)),
		RTTMs:      make(map[RegionID]map[RegionID]float64, len(specs)),
		AZDown:     make(map[AZID]bool),
		RegionDown: make(map[RegionID]bool),
	}
	for _, s := range specs {
		t.Regions[s.ID] = s
		t.Order = append(t.Order, s.ID)
		t.RTTMs[s.ID] = make(map[RegionID]float64, len(specs))
		t.RTTMs[s.ID][s.ID] = s.IntraAZRTTMs
	}

	// Symmetric inter-region WAN round-trip latencies (ms).
	setSym := func(a, b RegionID, rtt float64) {
		t.RTTMs[a][b] = rtt
		t.RTTMs[b][a] = rtt
	}
	setSym(RegionUSEast, RegionUSWest, 68.0)
	setSym(RegionUSEast, RegionEUCentral, 86.0)
	setSym(RegionUSEast, RegionAPNortheast, 168.0)
	setSym(RegionUSEast, RegionSAEast, 118.0)
	setSym(RegionUSWest, RegionEUCentral, 142.0)
	setSym(RegionUSWest, RegionAPNortheast, 108.0)
	setSym(RegionUSWest, RegionSAEast, 164.0)
	setSym(RegionEUCentral, RegionAPNortheast, 224.0)
	setSym(RegionEUCentral, RegionSAEast, 196.0)
	setSym(RegionAPNortheast, RegionSAEast, 268.0)

	return t
}

// MultiEUTopology extends DefaultTopology with a second EU region (eu-west) to
// exercise intra-EU multi-region sovereign failover without leaving EU borders.
func MultiEUTopology() *Topology {
	t := DefaultTopology()
	euWest := RegionSpec{
		ID:                    RegionEUWest,
		Jurisdiction:          JurisdictionEU,
		AZs:                   []AZID{"eu-west-a", "eu-west-b", "eu-west-c"},
		StorageCostPerTBMonth: 24.00,
		EgressCostPerGB:       0.020,
		IntraAZRTTMs:          0.6,
		CrossAZRTTMs:          1.8,
		StorageIOBaseMs:       1.2,
		WANBandwidthMBps:      240.0,
		LocalBandwidthMBps:    1200.0,
	}
	t.Regions[RegionEUWest] = euWest
	t.Order = append(t.Order, RegionEUWest)
	t.RTTMs[RegionEUWest] = map[RegionID]float64{
		RegionEUWest: euWest.IntraAZRTTMs,
	}
	setSym := func(a, b RegionID, rtt float64) {
		t.RTTMs[a][b] = rtt
		t.RTTMs[b][a] = rtt
	}
	setSym(RegionEUWest, RegionEUCentral, 18.0)
	setSym(RegionEUWest, RegionUSEast, 74.0)
	setSym(RegionEUWest, RegionUSWest, 134.0)
	setSym(RegionEUWest, RegionAPNortheast, 232.0)
	setSym(RegionEUWest, RegionSAEast, 188.0)
	return t
}

// Clone returns an isolated deep copy of Topology so benchmark runs never share fault state.
func (t *Topology) Clone() *Topology {
	cp := &Topology{
		Regions:    make(map[RegionID]RegionSpec, len(t.Regions)),
		Order:      append([]RegionID(nil), t.Order...),
		RTTMs:      make(map[RegionID]map[RegionID]float64, len(t.RTTMs)),
		AZDown:     make(map[AZID]bool, len(t.AZDown)),
		RegionDown: make(map[RegionID]bool, len(t.RegionDown)),
	}
	for k, v := range t.Regions {
		specCopy := v
		specCopy.AZs = append([]AZID(nil), v.AZs...)
		cp.Regions[k] = specCopy
	}
	for a, row := range t.RTTMs {
		cp.RTTMs[a] = make(map[RegionID]float64, len(row))
		for b, rtt := range row {
			cp.RTTMs[a][b] = rtt
		}
	}
	for k, v := range t.AZDown {
		cp.AZDown[k] = v
	}
	for k, v := range t.RegionDown {
		cp.RegionDown[k] = v
	}
	return cp
}

// LinkRTTMs returns round-trip time in milliseconds between two regions.
func (t *Topology) LinkRTTMs(from, to RegionID) float64 {
	if row, ok := t.RTTMs[from]; ok {
		if rtt, ok := row[to]; ok {
			return rtt
		}
	}
	return 150.0
}

// TransferTimeMs returns modeled network transfer latency in ms for sizeBytes.
func (t *Topology) TransferTimeMs(from, to RegionID, sizeBytes int) float64 {
	if sizeBytes <= 0 {
		return 0
	}
	spec := t.Regions[from]
	bw := spec.WANBandwidthMBps
	if from == to {
		bw = spec.LocalBandwidthMBps
	}
	if bw <= 0 {
		bw = 200.0
	}
	// MBps -> bytes/ms = bw * 1e6 / 1000 = bw * 1000.
	return float64(sizeBytes) / (bw * 1000.0)
}

// IsRegionUp reports whether a region is currently healthy.
func (t *Topology) IsRegionUp(r RegionID) bool {
	if _, ok := t.Regions[r]; !ok {
		return false
	}
	return !t.RegionDown[r]
}

// IsAZUp reports whether an availability zone and its parent region are healthy.
func (t *Topology) IsAZUp(r RegionID, az AZID) bool {
	if !t.IsRegionUp(r) {
		return false
	}
	return !t.AZDown[az]
}

// HealthyAZs returns all healthy AZs inside region r.
func (t *Topology) HealthyAZs(r RegionID) []AZID {
	if !t.IsRegionUp(r) {
		return nil
	}
	spec := t.Regions[r]
	out := make([]AZID, 0, len(spec.AZs))
	for _, az := range spec.AZs {
		if !t.AZDown[az] {
			out = append(out, az)
		}
	}
	return out
}

// PermittedRegions returns all configured regions allowed by tp, sorted deterministically.
func (t *Topology) PermittedRegions(tp TenantPolicy) []RegionID {
	out := make([]RegionID, 0, len(t.Order))
	for _, r := range t.Order {
		if tp.PermitsRegion(t.Regions[r]) {
			out = append(out, r)
		}
	}
	return out
}

// NearestPeerRegions returns healthy, tenant-permitted regions other than origin,
// ordered strictly by ascending RTT from origin (tie-broken by RegionID).
func (t *Topology) NearestPeerRegions(origin RegionID, tp TenantPolicy, includeDown bool) []RegionID {
	peers := make([]RegionID, 0, len(t.Order))
	for _, r := range t.Order {
		if r == origin {
			continue
		}
		if !tp.PermitsRegion(t.Regions[r]) {
			continue
		}
		if !includeDown && !t.IsRegionUp(r) {
			continue
		}
		peers = append(peers, r)
	}
	sort.Slice(peers, func(i, j int) bool {
		ri := t.LinkRTTMs(origin, peers[i])
		rj := t.LinkRTTMs(origin, peers[j])
		if ri != rj {
			return ri < rj
		}
		return peers[i] < peers[j]
	})
	return peers
}
