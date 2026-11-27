package polarstore

import (
	"math"
	"sync"

	"github.com/klauspost/compress/s2"
)

// ProbeMetrics captures the deterministic compressibility signals measured
// during the fast sampling probe.
type ProbeMetrics struct {
	SampledBytes      int     `json:"sampled_bytes"`
	ShannonEntropy    float64 `json:"shannon_entropy_bits"`
	UniqueBytes       int     `json:"unique_bytes"`
	PrecompressedType string  `json:"precompressed_type,omitempty"`
	TrialS2Ratio      float64 `json:"trial_s2_ratio,omitempty"`
	RanTrialProbe     bool    `json:"ran_trial_probe"`
}

// ProbeDecision is the codec recommendation and diagnostic metrics returned by Probe.
type ProbeDecision struct {
	Codec        Codec
	BypassReason BypassReason
	Metrics      ProbeMetrics
}

var probeScratchPool = sync.Pool{
	New: func() any {
		buf := make([]byte, 16384)
		return &buf
	},
}

// DetectPrecompressedMagic inspects leading container magic headers in O(1) time
// to avoid spending CPU re-compressing already-compressed formats.
func DetectPrecompressedMagic(data []byte) (string, bool) {
	n := len(data)
	if n < 2 {
		return "", false
	}
	// Gzip: 1f 8b
	if data[0] == 0x1f && data[1] == 0x8b {
		return "gzip", true
	}
	// JPEG: ff d8 ff
	if n >= 3 && data[0] == 0xff && data[1] == 0xd8 && data[2] == 0xff {
		return "jpeg", true
	}
	// Bzip2: BZh
	if n >= 3 && data[0] == 'B' && data[1] == 'Z' && data[2] == 'h' {
		return "bzip2", true
	}
	if n >= 4 {
		// Zstd: 28 b5 2f fd
		if data[0] == 0x28 && data[1] == 0xb5 && data[2] == 0x2f && data[3] == 0xfd {
			return "zstd", true
		}
		// ZIP / JAR / DOCX / XLSX: PK\x03\x04
		if data[0] == 0x50 && data[1] == 0x4b && data[2] == 0x03 && data[3] == 0x04 {
			return "zip", true
		}
		// LZ4 frame: 04 22 4d 18
		if data[0] == 0x04 && data[1] == 0x22 && data[2] == 0x4d && data[3] == 0x18 {
			return "lz4", true
		}
		// S2/Snappy framed stream identifier: ff 06 00 00
		if data[0] == 0xff && data[1] == 0x06 && data[2] == 0x00 && data[3] == 0x00 {
			return "s2_framed", true
		}
		// Ogg container
		if data[0] == 'O' && data[1] == 'g' && data[2] == 'g' && data[3] == 'S' {
			return "ogg", true
		}
		// FLAC / MP3 / ID3
		if data[0] == 'f' && data[1] == 'L' && data[2] == 'a' && data[3] == 'C' {
			return "flac", true
		}
	}
	// XZ: fd 37 7a 58 5a 00
	if n >= 6 && data[0] == 0xfd && data[1] == 0x37 && data[2] == 0x7a &&
		data[3] == 0x58 && data[4] == 0x5a && data[5] == 0x00 {
		return "xz", true
	}
	// PNG: 89 50 4e 47 0d 0a 1a 0a
	if n >= 8 && data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4e && data[3] == 0x47 &&
		data[4] == 0x0d && data[5] == 0x0a && data[6] == 0x1a && data[7] == 0x0a {
		return "png", true
	}
	// WebP: RIFF....WEBP
	if n >= 12 && data[0] == 'R' && data[1] == 'I' && data[2] == 'F' && data[3] == 'F' &&
		data[8] == 'W' && data[9] == 'E' && data[10] == 'B' && data[11] == 'P' {
		return "webp", true
	}
	// ISO BMFF / MP4 / HEIF: ....ftyp
	if n >= 12 && data[4] == 'f' && data[5] == 't' && data[6] == 'y' && data[7] == 'p' {
		return "mp4", true
	}
	return "", false
}

// extractDeterministicSample collects up to maxBytes across head (50%), middle (25%),
// and tail (25%) windows so header-heavy formats (e.g., tarballs, columnar files)
// do not mislead compressibility estimation.
func extractDeterministicSample(data []byte, maxBytes int, dst []byte) []byte {
	if len(data) <= maxBytes {
		return data
	}
	if cap(dst) < maxBytes {
		dst = make([]byte, maxBytes)
	} else {
		dst = dst[:maxBytes]
	}

	headLen := maxBytes / 2
	midLen := maxBytes / 4
	tailLen := maxBytes - headLen - midLen

	midStart := (len(data) - midLen) / 2
	tailStart := len(data) - tailLen

	copy(dst[0:headLen], data[0:headLen])
	copy(dst[headLen:headLen+midLen], data[midStart:midStart+midLen])
	copy(dst[headLen+midLen:], data[tailStart:])
	return dst
}

// ComputeEntropyAndUniqueness calculates Shannon entropy (bits/byte in [0, 8])
// and distinct byte count using a zero-allocation stack histogram.
func ComputeEntropyAndUniqueness(sample []byte) (entropy float64, unique int) {
	if len(sample) == 0 {
		return 0, 0
	}
	var hist [256]uint32
	for _, b := range sample {
		hist[b]++
	}
	invN := 1.0 / float64(len(sample))
	var h float64
	for _, c := range hist {
		if c == 0 {
			continue
		}
		unique++
		p := float64(c) * invN
		h -= p * math.Log2(p)
	}
	return h, unique
}

// Probe evaluates payload compressibility in microseconds and selects the optimal
// codec tier (CodecNone, CodecS2, or CodecZstd) for the given context and policy.
func Probe(ctx ObjectContext, data []byte, policy EncryptionBoundaryPolicy) ProbeDecision {
	if len(data) < policy.MinObjectBytes {
		return ProbeDecision{
			Codec:        CodecNone,
			BypassReason: BypassBelowMinSize,
		}
	}

	// Stage 1: O(1) container magic sniff for already-compressed formats.
	if kind, ok := DetectPrecompressedMagic(data); ok {
		return ProbeDecision{
			Codec:        CodecNone,
			BypassReason: BypassPrecompressedMagic,
			Metrics: ProbeMetrics{
				PrecompressedType: kind,
			},
		}
	}

	// Stage 2: Deterministic multi-window sample + zero-alloc Shannon entropy screen.
	samplePtr := probeScratchPool.Get().(*[]byte)
	defer probeScratchPool.Put(samplePtr)

	sample := extractDeterministicSample(data, policy.ProbeSampleBytes, *samplePtr)
	entropy, unique := ComputeEntropyAndUniqueness(sample)
	metrics := ProbeMetrics{
		SampledBytes:   len(sample),
		ShannonEntropy: entropy,
		UniqueBytes:    unique,
	}

	if entropy >= policy.MaxEntropyBitsPerByte && unique >= policy.MinUniqueBytesForHighEntropy {
		return ProbeDecision{
			Codec:        CodecNone,
			BypassReason: BypassHighEntropy,
			Metrics:      metrics,
		}
	}

	// Stage 3: Fast S2 trial probe on the sampled slice (~5-15 us for 8 KiB).
	trialBufPtr := probeScratchPool.Get().(*[]byte)
	defer probeScratchPool.Put(trialBufPtr)

	maxEncoded := s2.MaxEncodedLen(len(sample))
	if cap(*trialBufPtr) < maxEncoded {
		*trialBufPtr = make([]byte, maxEncoded)
	}
	encoded := s2.Encode((*trialBufPtr)[:maxEncoded], sample)
	trialRatio := float64(len(encoded)) / float64(len(sample))
	metrics.TrialS2Ratio = trialRatio
	metrics.RanTrialProbe = true

	maxAcceptableRatio := 1.0 - policy.MinSavingsRatio

	// If S2 achieves insufficient savings AND entropy is not low enough for Zstd FSE
	// entropy coding to rescue it, bypass immediately.
	if trialRatio > maxAcceptableRatio && entropy >= 6.40 {
		return ProbeDecision{
			Codec:        CodecNone,
			BypassReason: BypassProbeInsufficientGain,
			Metrics:      metrics,
		}
	}

	// Stage 4: Adaptive codec tier selection (PolarStore-inspired hot vs warm/cold + CPU budget).
	preferFast := ctx.AccessTier == TierHot ||
		(policy.CPUBudgetPerMiB > 0 && policy.CPUBudgetPerMiB <= policy.FastTierBudgetCeiling)

	if preferFast {
		if policy.AllowS2 && trialRatio <= maxAcceptableRatio {
			return ProbeDecision{
				Codec:        CodecS2,
				BypassReason: BypassNone,
				Metrics:      metrics,
			}
		}
		if policy.AllowZstd && entropy < 6.40 {
			return ProbeDecision{
				Codec:        CodecZstd,
				BypassReason: BypassNone,
				Metrics:      metrics,
			}
		}
		return ProbeDecision{
			Codec:        CodecNone,
			BypassReason: BypassProbeInsufficientGain,
			Metrics:      metrics,
		}
	}

	// Warm / Cold tiers: prefer Zstd for superior storage reduction when enabled;
	// fall back to S2 if Zstd is disabled.
	if policy.AllowZstd {
		return ProbeDecision{
			Codec:        CodecZstd,
			BypassReason: BypassNone,
			Metrics:      metrics,
		}
	}
	if policy.AllowS2 && trialRatio <= maxAcceptableRatio {
		return ProbeDecision{
			Codec:        CodecS2,
			BypassReason: BypassNone,
			Metrics:      metrics,
		}
	}
	return ProbeDecision{
		Codec:        CodecNone,
		BypassReason: BypassPolicyDisabled,
		Metrics:      metrics,
	}
}
