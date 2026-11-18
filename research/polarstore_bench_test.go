package research

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"testing"
	"time"

	"github.com/klauspost/compress/zstd"
)

// Hypothesis: Sampling the first 8KB of a 1MB block to decide whether to compress
// will reduce CPU time by >70% for incompressible data (encrypted/media) 
// while retaining >95% of storage gains for compressible data.

func BenchmarkAdaptiveCompression(b *testing.B) {
	// Data types
	compressible := bytes.Repeat([]byte("hanzo-s3-storage-logic-test-data-12345"), 100000) // Highly compressible
	incompressible := make([]byte, len(compressible))
	rand.Read(incompressible) // Random data (simulating encrypted or compressed media)

	encoder, _ := zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.SpeedFastest))

	b.Run("AlwaysCompress_Compressible", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var buf bytes.Buffer
			encoder.Reset(&buf)
			encoder.Write(compressible)
			encoder.Close()
		}
	})

	b.Run("AlwaysCompress_Incompressible", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var buf bytes.Buffer
			encoder.Reset(&buf)
			encoder.Write(incompressible)
			encoder.Close()
		}
	})

	b.Run("Adaptive_Incompressible", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if isCompressible(incompressible, 4096, 0.9) {
				var buf bytes.Buffer
				encoder.Reset(&buf)
				encoder.Write(incompressible)
				encoder.Close()
			}
		}
	})
}

// isCompressible samples the prefix of data to determine if it's worth compressing.
// This is a simplified version of the PolarStore heuristic.
func isCompressible(data []byte, sampleSize int, threshold float64) bool {
	if len(data) < sampleSize {
		sampleSize = len(data)
	}
	sample := data[:sampleSize]
	
	var buf bytes.Buffer
	enc, _ := zstd.NewWriter(&buf, zstd.WithEncoderLevel(zstd.SpeedFastest))
	enc.Write(sample)
	enc.Close()

	ratio := float64(buf.Len()) / float64(sampleSize)
	return ratio < threshold
}
