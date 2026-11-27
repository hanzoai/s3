package polarstore

import "time"

// DataPlaneMode explicitly classifies the cryptographic trust boundary of an object
// before the adaptive compression engine evaluates it.
//
// Security invariant (HIP-1348 / HIP-1349 / HIP-0302):
// Adaptive compression is permitted ONLY on ModePreEncryptionPlaintext.
// ModeUnknown, ModeClientSideCiphertext, and ModeAttackerControlledMixedSecret
// MUST fail closed and bypass compression without altering payload length.
type DataPlaneMode uint8

const (
	// ModeUnknown is the zero value. Compression fails closed (bypassed) whenever
	// a caller has not explicitly certified the object's cryptographic data-plane mode.
	ModeUnknown DataPlaneMode = 0

	// ModePreEncryptionPlaintext marks an object whose plaintext is legitimately
	// available to the storage server prior to server-side encryption (e.g., SSE-S3
	// or SSE-KMS on non-adversarial payloads) or unencrypted storage.
	ModePreEncryptionPlaintext DataPlaneMode = 1

	// ModeClientSideCiphertext marks an object already encrypted client-side
	// (including HIP-0302 zero-knowledge storage payloads). Compressing ciphertext
	// violates ciphertext opacity, yields ~1.00x ratio, and wastes CPU.
	ModeClientSideCiphertext DataPlaneMode = 2

	// ModeAttackerControlledMixedSecret marks plaintext that concatenates
	// attacker-influenced input with confidential secrets prior to encryption.
	// Compressing such payloads creates a CRIME/BREACH/HEIST length side-channel
	// oracle; compression MUST be bypassed.
	ModeAttackerControlledMixedSecret DataPlaneMode = 3
)

// String returns the canonical identifier for a DataPlaneMode.
func (m DataPlaneMode) String() string {
	switch m {
	case ModeUnknown:
		return "unknown"
	case ModePreEncryptionPlaintext:
		return "pre_encryption_plaintext"
	case ModeClientSideCiphertext:
		return "client_side_ciphertext"
	case ModeAttackerControlledMixedSecret:
		return "attacker_controlled_mixed_secret"
	default:
		return "invalid"
	}
}

// AccessTier mirrors PolarStore's workload-aware tiering (Hu et al., USENIX FAST '26):
// latency-sensitive hot paths favor lightweight compression (S2) or fast bypass,
// while warm/cold storage paths prioritize space savings (Zstd).
type AccessTier uint8

const (
	// TierWarm is the default balanced tier (Zstd for compressible objects, S2 under CPU pressure).
	TierWarm AccessTier = 0
	// TierHot prioritizes tail latency and throughput on critical I/O paths (prefers S2).
	TierHot AccessTier = 1
	// TierCold prioritizes storage ratio over CPU cost (prefers Zstd).
	TierCold AccessTier = 2
)

// String returns the canonical access tier label.
func (t AccessTier) String() string {
	switch t {
	case TierWarm:
		return "warm"
	case TierHot:
		return "hot"
	case TierCold:
		return "cold"
	default:
		return "warm"
	}
}

// BypassReason records the auditable reason why an object was stored uncompressed.
type BypassReason string

const (
	BypassNone                  BypassReason = ""
	BypassUnknownDataPlane      BypassReason = "unknown_data_plane_fail_closed"
	BypassClientSideCiphertext  BypassReason = "client_side_ciphertext_hip0302"
	BypassOracleRiskMixedSecret BypassReason = "crime_breach_mixed_secret_oracle"
	BypassPolicyDisabled        BypassReason = "policy_disabled"
	BypassBelowMinSize          BypassReason = "below_min_size_threshold"
	BypassPrecompressedMagic    BypassReason = "precompressed_magic_header"
	BypassHighEntropy           BypassReason = "high_entropy_incompressible"
	BypassProbeInsufficientGain BypassReason = "probe_insufficient_gain"
	BypassCPUBudgetExceeded     BypassReason = "cpu_budget_exceeded"
	BypassExpandedOutput        BypassReason = "compressed_output_not_smaller"
)

// ObjectContext carries tenant, key, and cryptographic boundary metadata for an object PUT.
type ObjectContext struct {
	TenantID    string
	Bucket      string
	Key         string
	ContentType string
	DataPlane   DataPlaneMode
	AccessTier  AccessTier
}

// EncryptionBoundaryPolicy governs both cryptographic safety boundaries and
// PolarStore-inspired adaptive codec selection thresholds.
type EncryptionBoundaryPolicy struct {
	// Enabled toggles adaptive compression for this policy scope.
	Enabled bool

	// Hard security boundary invariants (always enforced fail-closed; tenant overrides
	// cannot weaken these below the global security floor).
	FailClosedOnUnknownMode   bool
	BlockClientSideCiphertext bool
	BlockMixedSecretOracle    bool

	// MinObjectBytes skips compression for tiny objects to avoid framing/header expansion
	// and micro-object latency regression.
	MinObjectBytes int

	// ProbeSampleBytes is the maximum deterministic sample size inspected by the fast probe.
	ProbeSampleBytes int

	// MaxEntropyBitsPerByte is the Shannon entropy cutoff (bits/byte, in [0, 8]) above which
	// a high-diversity payload is immediately rejected as incompressible without trial coding.
	MaxEntropyBitsPerByte float64

	// MinUniqueBytesForHighEntropy prevents small-alphabet payloads from triggering the
	// entropy fast-reject when an entropy coder (Zstd/FSE) can still compress them well.
	MinUniqueBytesForHighEntropy int

	// MinSavingsRatio requires compressed size <= (1 - MinSavingsRatio) * originalSize.
	MinSavingsRatio float64

	// CPUBudgetPerMiB optionally constrains codec CPU cost per MiB of input.
	// When > 0 and <= FastTierBudgetCeiling, the selector favors CodecS2 over CodecZstd.
	CPUBudgetPerMiB time.Duration

	// FastTierBudgetCeiling is the CPU budget per MiB threshold at or below which
	// CodecS2 is chosen over CodecZstd to stay within CPU SLOs.
	FastTierBudgetCeiling time.Duration

	// AllowZstd permits selecting the balanced Zstd codec tier.
	AllowZstd bool
	// AllowS2 permits selecting the fast S2 codec tier.
	AllowS2 bool
}

// DefaultPolicy returns the production-recommended EncryptionBoundaryPolicy with
// strict HIP-0302 / CRIME-BREACH guards enabled and PolarStore-inspired adaptive thresholds.
func DefaultPolicy() EncryptionBoundaryPolicy {
	return EncryptionBoundaryPolicy{
		Enabled:                      true,
		FailClosedOnUnknownMode:      true,
		BlockClientSideCiphertext:    true,
		BlockMixedSecretOracle:       true,
		MinObjectBytes:               1024,
		ProbeSampleBytes:             8192,
		MaxEntropyBitsPerByte:        7.45,
		MinUniqueBytesForHighEntropy: 235,
		MinSavingsRatio:              0.10,
		CPUBudgetPerMiB:              0,
		FastTierBudgetCeiling:        3 * time.Millisecond,
		AllowZstd:                    true,
		AllowS2:                      true,
	}
}

// EnforceSecurityFloor returns a normalized copy of p guaranteed to satisfy
// non-negotiable HIP-0302 ciphertext and CRIME/BREACH oracle boundaries.
func (p EncryptionBoundaryPolicy) EnforceSecurityFloor() EncryptionBoundaryPolicy {
	p.FailClosedOnUnknownMode = true
	p.BlockClientSideCiphertext = true
	p.BlockMixedSecretOracle = true
	if p.MinObjectBytes <= 0 {
		p.MinObjectBytes = 1024
	}
	if p.ProbeSampleBytes < 512 {
		p.ProbeSampleBytes = 8192
	}
	if p.MaxEntropyBitsPerByte <= 0 || p.MaxEntropyBitsPerByte > 8.0 {
		p.MaxEntropyBitsPerByte = 7.45
	}
	if p.MinUniqueBytesForHighEntropy <= 0 || p.MinUniqueBytesForHighEntropy > 256 {
		p.MinUniqueBytesForHighEntropy = 235
	}
	if p.MinSavingsRatio < 0 || p.MinSavingsRatio >= 0.95 {
		p.MinSavingsRatio = 0.10
	}
	if p.FastTierBudgetCeiling <= 0 {
		p.FastTierBudgetCeiling = 3 * time.Millisecond
	}
	return p
}

// CheckEncryptionBoundary evaluates the hard cryptographic trust boundary before
// touching payload bytes. It returns BypassNone only when compression is allowed.
func (p EncryptionBoundaryPolicy) CheckEncryptionBoundary(ctx ObjectContext) BypassReason {
	if !p.Enabled {
		return BypassPolicyDisabled
	}
	switch ctx.DataPlane {
	case ModePreEncryptionPlaintext:
		return BypassNone
	case ModeClientSideCiphertext:
		return BypassClientSideCiphertext
	case ModeAttackerControlledMixedSecret:
		return BypassOracleRiskMixedSecret
	default:
		// Fail closed on ModeUnknown or any unrecognized value.
		return BypassUnknownDataPlane
	}
}
