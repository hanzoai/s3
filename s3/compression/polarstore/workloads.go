package polarstore

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
)

// BenchmarkObject represents a single deterministic S3 object in a reproducible corpus.
type BenchmarkObject struct {
	Workload string
	Context  ObjectContext
	Data     []byte
	SHA256   [32]byte
}

// GenerateBenchmarkCorpus creates a deterministic, seeded S3 object corpus spanning
// six realistic storage workloads:
//   - json_logs (structured application telemetry logs, 128 KiB)
//   - source_tarball (Go/TypeScript source archive, 128 KiB)
//   - columnar_telemetry (Parquet-like numeric column chunks, 128 KiB)
//   - precompressed_media (JPEG-header & Zstd-compressed blobs, 128 KiB)
//   - aes_gcm_ciphertext_zk (HIP-0302 client-side AES-256-GCM ciphertext, 128 KiB)
//   - tiny_metadata (256 B object manifests)
func GenerateBenchmarkCorpus(seed int64, objectsPerWorkload int) []BenchmarkObject {
	if objectsPerWorkload <= 0 {
		objectsPerWorkload = 16
	}
	rng := rand.New(rand.NewSource(seed))
	const stdSize = 128 * 1024 // 128 KiB per medium object

	out := make([]BenchmarkObject, 0, objectsPerWorkload*6)

	// 1. Structured JSON logs (pre-encryption plaintext, warm tier)
	for i := 0; i < objectsPerWorkload; i++ {
		data := makeJSONLogsPayload(rng, stdSize)
		out = append(out, BenchmarkObject{
			Workload: "json_logs",
			Context: ObjectContext{
				TenantID:    "tenant-logs",
				Bucket:      "prod-app-logs",
				Key:         fmt.Sprintf("2026/11/25/shard-%03d.jsonl", i),
				ContentType: "application/x-ndjson",
				DataPlane:   ModePreEncryptionPlaintext,
				AccessTier:  TierWarm,
			},
			Data:   data,
			SHA256: sha256.Sum256(data),
		})
	}

	// 2. Go / TypeScript source tarball (pre-encryption plaintext, cold/archive tier)
	for i := 0; i < objectsPerWorkload; i++ {
		data := makeSourceTarballPayload(rng, stdSize)
		out = append(out, BenchmarkObject{
			Workload: "source_tarball",
			Context: ObjectContext{
				TenantID:    "tenant-ci",
				Bucket:      "build-artifacts",
				Key:         fmt.Sprintf("repos/service-%03d.tar", i),
				ContentType: "application/x-tar",
				DataPlane:   ModePreEncryptionPlaintext,
				AccessTier:  TierCold,
			},
			Data:   data,
			SHA256: sha256.Sum256(data),
		})
	}

	// 3. Columnar numeric telemetry / Parquet-like rows (pre-encryption plaintext, hot tier)
	for i := 0; i < objectsPerWorkload; i++ {
		data := makeColumnarTelemetryPayload(rng, stdSize)
		out = append(out, BenchmarkObject{
			Workload: "columnar_telemetry",
			Context: ObjectContext{
				TenantID:    "tenant-analytics",
				Bucket:      "telemetry-hot",
				Key:         fmt.Sprintf("metrics/rg-%03d.col", i),
				ContentType: "application/octet-stream",
				DataPlane:   ModePreEncryptionPlaintext,
				AccessTier:  TierHot,
			},
			Data:   data,
			SHA256: sha256.Sum256(data),
		})
	}

	// 4. Pre-compressed media / Zstd blobs (pre-encryption plaintext, warm tier)
	for i := 0; i < objectsPerWorkload; i++ {
		var data []byte
		var key, ct string
		if i%2 == 0 {
			data = makeJPEGBlobPayload(rng, stdSize)
			key = fmt.Sprintf("media/photo-%03d.jpg", i)
			ct = "image/jpeg"
		} else {
			raw := makeJSONLogsPayload(rng, stdSize*3)
			z, err := CompressWithCodec(CodecZstd, raw)
			if err != nil {
				panic(err)
			}
			// Normalize to roughly stdSize while preserving real Zstd magic and stream.
			if len(z) < stdSize {
				// Pad with additional high-entropy bytes inside a synthetic precompressed archive
				// only if needed; better: generate larger raw until compressed ~stdSize.
				raw = makePseudoRandomLowRedundancy(rng, stdSize)
				z, _ = CompressWithCodec(CodecZstd, raw)
			}
			data = z
			key = fmt.Sprintf("archives/blob-%03d.zst", i)
			ct = "application/zstd"
		}
		out = append(out, BenchmarkObject{
			Workload: "precompressed_media",
			Context: ObjectContext{
				TenantID:    "tenant-media",
				Bucket:      "cdn-assets",
				Key:         key,
				ContentType: ct,
				DataPlane:   ModePreEncryptionPlaintext,
				AccessTier:  TierWarm,
			},
			Data:   data,
			SHA256: sha256.Sum256(data),
		})
	}

	// 5. Client-side AES-256-GCM ciphertext / HIP-0302 zero-knowledge storage payloads
	zkKey := sha256.Sum256([]byte("polarstore-deterministic-zk-corpus-key-v1"))
	for i := 0; i < objectsPerWorkload; i++ {
		plain := makeJSONLogsPayload(rng, stdSize)
		ct := encryptDeterministicGCM(zkKey[:], uint64(i+1), plain)
		out = append(out, BenchmarkObject{
			Workload: "aes_gcm_ciphertext_zk",
			Context: ObjectContext{
				TenantID:    "tenant-zk-vault",
				Bucket:      "hip0302-zk-objects",
				Key:         fmt.Sprintf("vault/enc-%03d.bin", i),
				ContentType: "application/octet-stream",
				DataPlane:   ModeClientSideCiphertext,
				AccessTier:  TierWarm,
			},
			Data:   ct,
			SHA256: sha256.Sum256(ct),
		})
	}

	// 6. Tiny 256 B metadata objects (pre-encryption plaintext, hot tier)
	for i := 0; i < objectsPerWorkload; i++ {
		data := makeTinyMetadataPayload(rng, i, 256)
		out = append(out, BenchmarkObject{
			Workload: "tiny_metadata",
			Context: ObjectContext{
				TenantID:    "tenant-meta",
				Bucket:      "object-manifests",
				Key:         fmt.Sprintf("meta/item-%04d.json", i),
				ContentType: "application/json",
				DataPlane:   ModePreEncryptionPlaintext,
				AccessTier:  TierHot,
			},
			Data:   data,
			SHA256: sha256.Sum256(data),
		})
	}

	return out
}

func makeJSONLogsPayload(rng *rand.Rand, targetBytes int) []byte {
	services := []string{"s3-gateway", "filer-meta", "volume-node", "kms-proxy", "iam-edge"}
	methods := []string{"GET", "PUT", "HEAD", "LIST", "DELETE"}
	regions := []string{"us-west-2", "us-east-1", "eu-central-1", "ap-northeast-1"}
	statuses := []int{200, 200, 200, 200, 201, 204, 304, 404, 429, 503}

	var buf bytes.Buffer
	buf.Grow(targetBytes + 256)
	ts := int64(1764000000000)
	for buf.Len() < targetBytes {
		ts += int64(rng.Intn(40) + 1)
		fmt.Fprintf(&buf,
			`{"ts":%d,"level":"INFO","service":%q,"region":%q,"method":%q,"bucket":"org-%02d-data","key":"prefix/%04d/chunk-%06d.bin","status":%d,"bytes":%d,"latency_us":%d,"trace":"tr-%016x"}`+"\n",
			ts,
			services[rng.Intn(len(services))],
			regions[rng.Intn(len(regions))],
			methods[rng.Intn(len(methods))],
			rng.Intn(16),
			rng.Intn(64),
			rng.Intn(20000),
			statuses[rng.Intn(len(statuses))],
			(rng.Intn(256)+1)*1024,
			rng.Intn(4500)+80,
			rng.Uint64(),
		)
	}
	return buf.Bytes()[:targetBytes]
}

func makeSourceTarballPayload(rng *rand.Rand, targetBytes int) []byte {
	var buf bytes.Buffer
	buf.Grow(targetBytes + 1024)
	fileIdx := 0
	for buf.Len() < targetBytes {
		// 512-byte zero-padded POSIX tar-style header block
		var hdr [512]byte
		name := fmt.Sprintf("pkg/worker/handler_%03d.go", fileIdx)
		copy(hdr[0:], name)
		copy(hdr[100:], "0000644\x00")
		copy(hdr[257:], "ustar\x0000")
		buf.Write(hdr[:])

		fmt.Fprintf(&buf, `package worker

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Handler%03d processes object lifecycle tasks with bounded retries.
type Handler%03d struct {
	MaxRetries int
	Timeout    time.Duration
	Bucket     string
}

func (h *Handler%03d) Execute(ctx context.Context, key string, size int64) error {
	if key == "" || size < 0 {
		return errors.New("invalid object descriptor")
	}
	deadline := time.Now().Add(h.Timeout)
	for attempt := 0; attempt <= h.MaxRetries; attempt++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("timeout processing %%s after %%d attempts (seed=%d)", key, attempt)
		}
	}
	return nil
}
`, fileIdx, fileIdx, fileIdx, rng.Intn(100000))
		// Align to 512-byte tar boundary
		if rem := buf.Len() % 512; rem != 0 {
			buf.Write(make([]byte, 512-rem))
		}
		fileIdx++
	}
	return buf.Bytes()[:targetBytes]
}

func makeColumnarTelemetryPayload(rng *rand.Rand, targetBytes int) []byte {
	// Structured 32-byte columnar records:
	// [0:8] monotonic timestamp (ms), [8:12] hostID (small set), [12:16] metricID,
	// [16:24] smoothly varying float64 gauge, [24:32] counter delta.
	const recSize = 32
	nRecs := targetBytes / recSize
	buf := make([]byte, targetBytes)
	ts := uint64(1764000000000)
	val := 42.5
	counter := uint64(100000)
	for i := 0; i < nRecs; i++ {
		off := i * recSize
		ts += uint64(100 + rng.Intn(5))
		hostID := uint32(i % 16)
		metricID := uint32((i / 16) % 8)
		val += (rng.Float64() - 0.5) * 0.25
		counter += uint64(rng.Intn(32))

		binary.LittleEndian.PutUint64(buf[off:off+8], ts)
		binary.LittleEndian.PutUint32(buf[off+8:off+12], hostID)
		binary.LittleEndian.PutUint32(buf[off+12:off+16], metricID)
		binary.LittleEndian.PutUint64(buf[off+16:off+24], math.Float64bits(val))
		binary.LittleEndian.PutUint64(buf[off+24:off+32], counter)
	}
	return buf
}

func makeJPEGBlobPayload(rng *rand.Rand, targetBytes int) []byte {
	buf := make([]byte, targetBytes)
	// Valid JPEG SOI + APP0 JFIF magic header followed by incompressible entropy stream
	// and JPEG EOI trailer.
	header := []byte{0xff, 0xd8, 0xff, 0xe0, 0x00, 0x10, 'J', 'F', 'I', 'F', 0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00}
	copy(buf, header)
	rng.Read(buf[len(header) : targetBytes-2])
	buf[targetBytes-2] = 0xff
	buf[targetBytes-1] = 0xd9
	return buf
}

func makePseudoRandomLowRedundancy(rng *rand.Rand, size int) []byte {
	buf := make([]byte, size)
	rng.Read(buf)
	return buf
}

func makeTinyMetadataPayload(rng *rand.Rand, idx, exactBytes int) []byte {
	base := fmt.Sprintf(`{"v":1,"id":"obj-%06d","etag":"%016x","cls":"STANDARD","chk":"crc64nvme","pad":"`, idx, rng.Uint64())
	var buf bytes.Buffer
	buf.WriteString(base)
	for buf.Len() < exactBytes-2 {
		buf.WriteByte('a' + byte(buf.Len()%26))
	}
	buf.WriteString(`"}`)
	return buf.Bytes()[:exactBytes]
}

func encryptDeterministicGCM(key []byte, nonceCounter uint64, plaintext []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}
	nonce := make([]byte, gcm.NonceSize())
	binary.BigEndian.PutUint64(nonce[4:], nonceCounter)
	return gcm.Seal(nonce, nonce, plaintext, nil)
}
