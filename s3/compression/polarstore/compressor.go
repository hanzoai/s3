package polarstore

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"
)

// ErrInvalidKeySize indicates an AES-256 key was not 32 bytes.
var ErrInvalidKeySize = errors.New("polarstore: AES-256 key must be 32 bytes")

// CompressResult holds the outcome of an AdaptiveCompressor.Compress call.
type CompressResult struct {
	Codec            Codec
	Payload          []byte
	OriginalSize     int
	StoredSize       int
	BypassReason     BypassReason
	Metrics          ProbeMetrics
	ProbeDuration    time.Duration
	CompressDuration time.Duration
}

// CompressionRatio returns OriginalSize / StoredSize (>= 1.0 when compressed).
func (r CompressResult) CompressionRatio() float64 {
	if r.StoredSize <= 0 {
		return 1.0
	}
	return float64(r.OriginalSize) / float64(r.StoredSize)
}

// SealedObject represents an object sealed with AES-256-GCM after adaptive
// pre-encryption compression (or fail-closed bypass).
type SealedObject struct {
	TenantID     string
	Bucket       string
	Key          string
	Codec        Codec
	DataPlane    DataPlaneMode
	BypassReason BypassReason
	OriginalSize int
	CompressedSz int
	Ciphertext   []byte // nonce (12B) || ciphertext || GCM tag (16B)
}

// StatsSnapshot is a point-in-time copy of lock-free compressor counters.
type StatsSnapshot struct {
	ObjectsTotal              int64
	ObjectsCompressedS2       int64
	ObjectsCompressedZstd     int64
	ObjectsBypassed           int64
	BypassedBoundaryGuard     int64
	BypassedBelowMinSize      int64
	BypassedPrecompressed     int64
	BypassedHighEntropy       int64
	BypassedProbeLowGain      int64
	BypassedExpandedOutput    int64
	OriginalBytesTotal        int64
	StoredBytesTotal          int64
	ProbeDurationNsTotal      int64
	CompressDurationNsTotal   int64
}

// AdaptiveCompressor coordinates fail-closed encryption-boundary enforcement,
// per-tenant compression policies, and PolarStore-inspired adaptive codec selection.
// All methods are safe for concurrent use by multiple goroutines.
type AdaptiveCompressor struct {
	defaultPolicy EncryptionBoundaryPolicy

	mu             sync.RWMutex
	tenantPolicies map[string]EncryptionBoundaryPolicy

	objectsTotal            atomic.Int64
	objectsCompressedS2     atomic.Int64
	objectsCompressedZstd   atomic.Int64
	objectsBypassed         atomic.Int64
	bypassedBoundaryGuard   atomic.Int64
	bypassedBelowMinSize    atomic.Int64
	bypassedPrecompressed   atomic.Int64
	bypassedHighEntropy     atomic.Int64
	bypassedProbeLowGain    atomic.Int64
	bypassedExpandedOutput  atomic.Int64
	originalBytesTotal      atomic.Int64
	storedBytesTotal        atomic.Int64
	probeDurationNsTotal    atomic.Int64
	compressDurationNsTotal atomic.Int64
}

// NewAdaptiveCompressor creates a concurrency-safe AdaptiveCompressor with the
// given default policy (normalized through EnforceSecurityFloor).
func NewAdaptiveCompressor(defaultPolicy EncryptionBoundaryPolicy) *AdaptiveCompressor {
	return &AdaptiveCompressor{
		defaultPolicy:  defaultPolicy.EnforceSecurityFloor(),
		tenantPolicies: make(map[string]EncryptionBoundaryPolicy),
	}
}

// SetTenantPolicy configures a tenant-specific policy. Non-negotiable security floor
// invariants (HIP-0302 ciphertext bypass, CRIME/BREACH oracle bypass, unknown mode
// fail-closed) are always enforced regardless of tenant input.
func (c *AdaptiveCompressor) SetTenantPolicy(tenantID string, policy EncryptionBoundaryPolicy) {
	if tenantID == "" {
		return
	}
	normalized := policy.EnforceSecurityFloor()
	c.mu.Lock()
	c.tenantPolicies[tenantID] = normalized
	c.mu.Unlock()
}

// EffectivePolicy returns the normalized policy active for tenantID.
func (c *AdaptiveCompressor) EffectivePolicy(tenantID string) EncryptionBoundaryPolicy {
	if tenantID != "" {
		c.mu.RLock()
		p, ok := c.tenantPolicies[tenantID]
		c.mu.RUnlock()
		if ok {
			return p
		}
	}
	return c.defaultPolicy
}

// Compress evaluates the encryption boundary and payload compressibility, returning
// either an adaptively compressed payload (S2 or Zstd) or the unmodified input slice.
func (c *AdaptiveCompressor) Compress(ctx ObjectContext, data []byte) (CompressResult, error) {
	policy := c.EffectivePolicy(ctx.TenantID)
	c.objectsTotal.Add(1)
	c.originalBytesTotal.Add(int64(len(data)))

	// Hard Gate 1: Encryption-boundary policy (zero payload inspection when blocked).
	if reason := policy.CheckEncryptionBoundary(ctx); reason != BypassNone {
		c.recordBypass(reason, len(data))
		return CompressResult{
			Codec:        CodecNone,
			Payload:      data,
			OriginalSize: len(data),
			StoredSize:   len(data),
			BypassReason: reason,
		}, nil
	}

	// Hard Gate 2: Fast sampling probe & adaptive tier selector.
	probeStart := time.Now()
	decision := Probe(ctx, data, policy)
	probeDur := time.Since(probeStart)
	c.probeDurationNsTotal.Add(probeDur.Nanoseconds())

	if decision.Codec == CodecNone {
		c.recordBypass(decision.BypassReason, len(data))
		return CompressResult{
			Codec:         CodecNone,
			Payload:       data,
			OriginalSize:  len(data),
			StoredSize:    len(data),
			BypassReason:  decision.BypassReason,
			Metrics:       decision.Metrics,
			ProbeDuration: probeDur,
		}, nil
	}

	// Stage 3: Execute chosen codec.
	compStart := time.Now()
	compressed, err := CompressWithCodec(decision.Codec, data)
	compDur := time.Since(compStart)
	c.compressDurationNsTotal.Add(compDur.Nanoseconds())
	if err != nil {
		return CompressResult{}, err
	}

	// Hard Gate 3: Non-expansion guarantee. Never store a payload that failed to meet
	// the minimum savings floor (prevents negative compression on adversarial tail data).
	maxStored := int(float64(len(data)) * (1.0 - policy.MinSavingsRatio))
	if len(compressed) >= len(data) || len(compressed) > maxStored {
		c.recordBypass(BypassExpandedOutput, len(data))
		return CompressResult{
			Codec:            CodecNone,
			Payload:          data,
			OriginalSize:     len(data),
			StoredSize:       len(data),
			BypassReason:     BypassExpandedOutput,
			Metrics:          decision.Metrics,
			ProbeDuration:    probeDur,
			CompressDuration: compDur,
		}, nil
	}

	switch decision.Codec {
	case CodecS2:
		c.objectsCompressedS2.Add(1)
	case CodecZstd:
		c.objectsCompressedZstd.Add(1)
	}
	c.storedBytesTotal.Add(int64(len(compressed)))

	return CompressResult{
		Codec:            decision.Codec,
		Payload:          compressed,
		OriginalSize:     len(data),
		StoredSize:       len(compressed),
		BypassReason:     BypassNone,
		Metrics:          decision.Metrics,
		ProbeDuration:    probeDur,
		CompressDuration: compDur,
	}, nil
}

// Decompress restores the exact original plaintext bytes for a stored payload.
func (c *AdaptiveCompressor) Decompress(codec Codec, stored []byte) ([]byte, error) {
	return DecompressWithCodec(codec, stored)
}

// CompressAndSeal enforces the safe pre-encryption pipeline:
//  1. Run AdaptiveCompressor.Compress (which strictly bypasses ModeClientSideCiphertext,
//     ModeAttackerControlledMixedSecret, and ModeUnknown).
//  2. Seal the resulting bytes with AES-256-GCM binding TenantID, Bucket, Key, Codec,
//     and DataPlane into Additional Authenticated Data (AAD) to prevent codec-confusion
//     or cross-tenant replay tampering.
//
// Note: If ctx.DataPlane == ModeClientSideCiphertext, data is already encrypted by the
// client (HIP-0302); CompressAndSeal returns it untouched as SealedObject.Ciphertext
// with CodecNone and BypassClientSideCiphertext so no double-encryption or compression occurs.
func (c *AdaptiveCompressor) CompressAndSeal(ctx ObjectContext, data []byte, key []byte) (SealedObject, error) {
	if ctx.DataPlane == ModeClientSideCiphertext {
		res, err := c.Compress(ctx, data)
		if err != nil {
			return SealedObject{}, err
		}
		return SealedObject{
			TenantID:     ctx.TenantID,
			Bucket:       ctx.Bucket,
			Key:          ctx.Key,
			Codec:        CodecNone,
			DataPlane:    ctx.DataPlane,
			BypassReason: res.BypassReason,
			OriginalSize: len(data),
			CompressedSz: len(data),
			Ciphertext:   res.Payload,
		}, nil
	}

	if len(key) != 32 {
		return SealedObject{}, ErrInvalidKeySize
	}

	res, err := c.Compress(ctx, data)
	if err != nil {
		return SealedObject{}, err
	}

	ct, err := sealAESGCM(key, res.Payload, buildAAD(ctx.TenantID, ctx.Bucket, ctx.Key, res.Codec, ctx.DataPlane))
	if err != nil {
		return SealedObject{}, err
	}

	return SealedObject{
		TenantID:     ctx.TenantID,
		Bucket:       ctx.Bucket,
		Key:          ctx.Key,
		Codec:        res.Codec,
		DataPlane:    ctx.DataPlane,
		BypassReason: res.BypassReason,
		OriginalSize: res.OriginalSize,
		CompressedSz: res.StoredSize,
		Ciphertext:   ct,
	}, nil
}

// OpenAndDecompress authenticates and decrypts a SealedObject with AES-256-GCM and
// restores its exact original bytes via Decompress.
func (c *AdaptiveCompressor) OpenAndDecompress(sealed SealedObject, key []byte) ([]byte, error) {
	if sealed.DataPlane == ModeClientSideCiphertext {
		return sealed.Ciphertext, nil
	}
	if len(key) != 32 {
		return nil, ErrInvalidKeySize
	}
	aad := buildAAD(sealed.TenantID, sealed.Bucket, sealed.Key, sealed.Codec, sealed.DataPlane)
	plainOrCompressed, err := openAESGCM(key, sealed.Ciphertext, aad)
	if err != nil {
		return nil, err
	}
	return c.Decompress(sealed.Codec, plainOrCompressed)
}

// Stats returns a point-in-time snapshot of the compressor's telemetry counters.
func (c *AdaptiveCompressor) Stats() StatsSnapshot {
	return StatsSnapshot{
		ObjectsTotal:            c.objectsTotal.Load(),
		ObjectsCompressedS2:     c.objectsCompressedS2.Load(),
		ObjectsCompressedZstd:   c.objectsCompressedZstd.Load(),
		ObjectsBypassed:         c.objectsBypassed.Load(),
		BypassedBoundaryGuard:   c.bypassedBoundaryGuard.Load(),
		BypassedBelowMinSize:    c.bypassedBelowMinSize.Load(),
		BypassedPrecompressed:   c.bypassedPrecompressed.Load(),
		BypassedHighEntropy:     c.bypassedHighEntropy.Load(),
		BypassedProbeLowGain:    c.bypassedProbeLowGain.Load(),
		BypassedExpandedOutput:  c.bypassedExpandedOutput.Load(),
		OriginalBytesTotal:      c.originalBytesTotal.Load(),
		StoredBytesTotal:        c.storedBytesTotal.Load(),
		ProbeDurationNsTotal:    c.probeDurationNsTotal.Load(),
		CompressDurationNsTotal: c.compressDurationNsTotal.Load(),
	}
}

func (c *AdaptiveCompressor) recordBypass(reason BypassReason, n int) {
	c.objectsBypassed.Add(1)
	c.storedBytesTotal.Add(int64(n))
	switch reason {
	case BypassUnknownDataPlane, BypassClientSideCiphertext, BypassOracleRiskMixedSecret, BypassPolicyDisabled:
		c.bypassedBoundaryGuard.Add(1)
	case BypassBelowMinSize:
		c.bypassedBelowMinSize.Add(1)
	case BypassPrecompressedMagic:
		c.bypassedPrecompressed.Add(1)
	case BypassHighEntropy:
		c.bypassedHighEntropy.Add(1)
	case BypassProbeInsufficientGain, BypassCPUBudgetExceeded:
		c.bypassedProbeLowGain.Add(1)
	case BypassExpandedOutput:
		c.bypassedExpandedOutput.Add(1)
	}
}

func buildAAD(tenantID, bucket, key string, codec Codec, mode DataPlaneMode) []byte {
	return []byte(fmt.Sprintf("polarstore-v1|%s|%s|%s|%s|%s", tenantID, bucket, key, codec.String(), mode.String()))
}

func sealAESGCM(key, plaintext, aad []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, plaintext, aad), nil
}

func openAESGCM(key, ciphertext, aad []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	ns := gcm.NonceSize()
	if len(ciphertext) < ns+gcm.Overhead() {
		return nil, errors.New("polarstore: ciphertext too short")
	}
	return gcm.Open(nil, ciphertext[:ns], ciphertext[ns:], aad)
}
