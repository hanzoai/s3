// Package acos implements a clean-room, geo-aware active object placement,
// multi-region durability, and outage failover/recovery engine inspired by
// Apple ACOS (Baron et al., USENIX FAST '26).
//
// Unlike scheduled cold archival backup (s3/filer/archive, HIP-1342), this
// package governs active hot/warm object placement, access-locality adaptation,
// tenant data-residency enforcement, local-first AZ recovery, and full-region
// blackout failover.
package acos

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
)

var (
	// ErrInvalidShardConfig indicates invalid (k, m) erasure parameters.
	ErrInvalidShardConfig = errors.New("acos: invalid erasure shard configuration")
	// ErrInsufficientShards indicates fewer than k valid shards survived.
	ErrInsufficientShards = errors.New("acos: insufficient surviving shards for reconstruction")
	// ErrIntegrityMismatch indicates reconstructed or read bytes failed SHA-256 verification.
	ErrIntegrityMismatch = errors.New("acos: payload SHA-256 integrity verification failed")
)

// gf256 implements finite field arithmetic over GF(2^8) with primitive
// polynomial x^8 + x^4 + x^3 + x^2 + 1 (0x11d).
var (
	gfExp [512]byte
	gfLog [256]byte
)

func init() {
	x := 1
	for i := 0; i < 255; i++ {
		gfExp[i] = byte(x)
		gfLog[x] = byte(i)
		x <<= 1
		if x&0x100 != 0 {
			x ^= 0x11d
		}
	}
	for i := 255; i < 512; i++ {
		gfExp[i] = gfExp[i-255]
	}
}

func gfMul(a, b byte) byte {
	if a == 0 || b == 0 {
		return 0
	}
	return gfExp[int(gfLog[a])+int(gfLog[b])]
}

func gfInv(a byte) byte {
	if a == 0 {
		panic("acos: division by zero in GF(2^8)")
	}
	return gfExp[255-int(gfLog[a])]
}

// cauchyParityRow computes row r (0 <= r < m) of a k-column Cauchy parity matrix.
// Every square submatrix of [I_k ; C_{m x k}] over GF(2^8) is guaranteed invertible
// whenever k + m <= 256.
func cauchyParityRow(r, k int) []byte {
	row := make([]byte, k)
	x := byte(k + r)
	for c := 0; c < k; c++ {
		y := byte(c)
		row[c] = gfInv(x ^ y)
	}
	return row
}

// EncodeShards splits payload into k equal-length data shards (zero-padded to
// a multiple of k) and computes m Cauchy Reed-Solomon parity shards.
func EncodeShards(payload []byte, k, m int) ([][]byte, error) {
	if k <= 0 || m < 0 || k+m > 255 {
		return nil, fmt.Errorf("%w: k=%d m=%d", ErrInvalidShardConfig, k, m)
	}
	shardLen := (len(payload) + k - 1) / k
	if shardLen == 0 {
		shardLen = 1
	}

	shards := make([][]byte, k+m)
	for i := 0; i < k; i++ {
		shards[i] = make([]byte, shardLen)
		start := i * shardLen
		if start < len(payload) {
			end := start + shardLen
			if end > len(payload) {
				end = len(payload)
			}
			copy(shards[i], payload[start:end])
		}
	}

	for r := 0; r < m; r++ {
		parity := make([]byte, shardLen)
		coeffs := cauchyParityRow(r, k)
		for c := 0; c < k; c++ {
			coeff := coeffs[c]
			if coeff == 0 {
				continue
			}
			src := shards[c]
			for b := 0; b < shardLen; b++ {
				parity[b] ^= gfMul(coeff, src[b])
			}
		}
		shards[k+r] = parity
	}
	return shards, nil
}

// ReconstructPayload recovers the exact origLen-byte payload from any k
// non-nil shards in a k+m stripe. Missing shards must be nil.
func ReconstructPayload(shards [][]byte, origLen, k, m int) ([]byte, error) {
	if k <= 0 || m < 0 || len(shards) != k+m || origLen < 0 {
		return nil, fmt.Errorf("%w: len=%d k=%d m=%d origLen=%d", ErrInvalidShardConfig, len(shards), k, m, origLen)
	}

	// Fast path: all k data shards are present.
	allDataPresent := true
	shardLen := 0
	for i := 0; i < k; i++ {
		if shards[i] == nil {
			allDataPresent = false
			break
		}
		if shardLen == 0 {
			shardLen = len(shards[i])
		} else if len(shards[i]) != shardLen {
			return nil, fmt.Errorf("%w: mismatched shard length", ErrInvalidShardConfig)
		}
	}
	if allDataPresent {
		out := make([]byte, 0, k*shardLen)
		for i := 0; i < k; i++ {
			out = append(out, shards[i]...)
		}
		if origLen > len(out) {
			return nil, fmt.Errorf("%w: origLen=%d exceeds shard bytes=%d", ErrInvalidShardConfig, origLen, len(out))
		}
		return out[:origLen], nil
	}

	// Collect exactly k surviving shards and build their k x k encoding submatrix.
	presentIdx := make([]int, 0, k)
	for i := 0; i < k+m && len(presentIdx) < k; i++ {
		if shards[i] != nil {
			if shardLen == 0 {
				shardLen = len(shards[i])
			} else if len(shards[i]) != shardLen {
				return nil, fmt.Errorf("%w: mismatched shard length", ErrInvalidShardConfig)
			}
			presentIdx = append(presentIdx, i)
		}
	}
	if len(presentIdx) < k {
		return nil, fmt.Errorf("%w: have %d of %d required", ErrInsufficientShards, len(presentIdx), k)
	}

	sub := make([][]byte, k)
	for row, shardIdx := range presentIdx {
		if shardIdx < k {
			v := make([]byte, k)
			v[shardIdx] = 1
			sub[row] = v
		} else {
			sub[row] = cauchyParityRow(shardIdx-k, k)
		}
	}

	inv, err := invertMatrixGF256(sub)
	if err != nil {
		return nil, err
	}

	// Reconstruct the k original data shards by multiplying inv * survivingShards.
	dataShards := make([][]byte, k)
	for d := 0; d < k; d++ {
		buf := make([]byte, shardLen)
		for col, shardIdx := range presentIdx {
			coeff := inv[d][col]
			if coeff == 0 {
				continue
			}
			src := shards[shardIdx]
			for b := 0; b < shardLen; b++ {
				buf[b] ^= gfMul(coeff, src[b])
			}
		}
		dataShards[d] = buf
	}

	out := make([]byte, 0, k*shardLen)
	for d := 0; d < k; d++ {
		out = append(out, dataShards[d]...)
	}
	if origLen > len(out) {
		return nil, fmt.Errorf("%w: origLen=%d exceeds reconstructed bytes=%d", ErrInvalidShardConfig, origLen, len(out))
	}
	return out[:origLen], nil
}

// invertMatrixGF256 performs Gauss-Jordan elimination over GF(2^8).
func invertMatrixGF256(mat [][]byte) ([][]byte, error) {
	n := len(mat)
	aug := make([][]byte, n)
	for i := 0; i < n; i++ {
		row := make([]byte, 2*n)
		copy(row[:n], mat[i])
		row[n+i] = 1
		aug[i] = row
	}

	for col := 0; col < n; col++ {
		pivot := -1
		for r := col; r < n; r++ {
			if aug[r][col] != 0 {
				pivot = r
				break
			}
		}
		if pivot == -1 {
			return nil, fmt.Errorf("%w: singular decoding matrix", ErrInsufficientShards)
		}
		if pivot != col {
			aug[col], aug[pivot] = aug[pivot], aug[col]
		}

		scale := gfInv(aug[col][col])
		for c := 0; c < 2*n; c++ {
			aug[col][c] = gfMul(aug[col][c], scale)
		}

		for r := 0; r < n; r++ {
			if r == col {
				continue
			}
			factor := aug[r][col]
			if factor == 0 {
				continue
			}
			for c := 0; c < 2*n; c++ {
				aug[r][c] ^= gfMul(factor, aug[col][c])
			}
		}
	}

	inv := make([][]byte, n)
	for i := 0; i < n; i++ {
		inv[i] = make([]byte, n)
		copy(inv[i], aug[i][n:])
	}
	return inv, nil
}

// PayloadDigest returns the hex-encoded SHA-256 digest of data.
func PayloadDigest(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
