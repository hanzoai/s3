package polarstore

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func testAESKey(label string) []byte {
	sum := sha256.Sum256([]byte("polarstore-test-key-" + label))
	return sum[:]
}

func TestRoundTripSHA256FidelityAllCodecsAndSizes(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	sizes := []int{0, 1, 256, 1024, 4096, 65536, 262144, 1024 * 1024}
	codecs := []Codec{CodecNone, CodecS2, CodecZstd, CodecGzip}

	for _, sz := range sizes {
		var payload []byte
		if sz > 0 {
			payload = makeJSONLogsPayload(rng, sz)
		}
		wantHash := sha256.Sum256(payload)

		for _, codec := range codecs {
			compressed, err := CompressWithCodec(codec, payload)
			if err != nil {
				t.Fatalf("CompressWithCodec(%s, len=%d): %v", codec, sz, err)
			}
			restored, err := DecompressWithCodec(codec, compressed)
			if err != nil {
				t.Fatalf("DecompressWithCodec(%s, len=%d): %v", codec, sz, err)
			}
			gotHash := sha256.Sum256(restored)
			if gotHash != wantHash {
				t.Fatalf("SHA-256 mismatch for codec=%s len=%d", codec, sz)
			}
		}

		// Also verify full AdaptiveCompressor path across access tiers.
		ac := NewAdaptiveCompressor(DefaultPolicy())
		for _, tier := range []AccessTier{TierHot, TierWarm, TierCold} {
			ctx := ObjectContext{
				TenantID:   "tenant-fidelity",
				Bucket:     "bucket-fidelity",
				Key:        fmt.Sprintf("obj-%d-%s.json", sz, tier),
				DataPlane:  ModePreEncryptionPlaintext,
				AccessTier: tier,
			}
			res, err := ac.Compress(ctx, payload)
			if err != nil {
				t.Fatalf("AdaptiveCompressor.Compress(len=%d, tier=%s): %v", sz, tier, err)
			}
			restored, err := ac.Decompress(res.Codec, res.Payload)
			if err != nil {
				t.Fatalf("AdaptiveCompressor.Decompress(codec=%s, len=%d): %v", res.Codec, sz, err)
			}
			if sha256.Sum256(restored) != wantHash {
				t.Fatalf("adaptive round-trip SHA-256 mismatch len=%d tier=%s", sz, tier)
			}
		}
	}
}

func TestEnvelopeRoundTripAndTamperRejection(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	plain := makeSourceTarballPayload(rng, 64*1024)
	wantHash := sha256.Sum256(plain)

	ac := NewAdaptiveCompressor(DefaultPolicy())
	res, err := ac.Compress(ObjectContext{
		TenantID:   "t1",
		Bucket:     "b1",
		Key:        "archive.tar",
		DataPlane:  ModePreEncryptionPlaintext,
		AccessTier: TierCold,
	}, plain)
	if err != nil {
		t.Fatalf("Compress: %v", err)
	}
	if res.Codec != CodecZstd {
		t.Fatalf("expected CodecZstd for cold source tarball, got %s", res.Codec)
	}

	frame := MarshalEnvelope(res.Codec, uint32(res.OriginalSize), res.Payload)
	codec, origLen, body, err := UnmarshalEnvelope(frame)
	if err != nil {
		t.Fatalf("UnmarshalEnvelope: %v", err)
	}
	if codec != res.Codec || int(origLen) != len(plain) {
		t.Fatalf("envelope metadata mismatch: codec=%s origLen=%d", codec, origLen)
	}
	restored, err := ac.Decompress(codec, body)
	if err != nil {
		t.Fatalf("Decompress: %v", err)
	}
	if sha256.Sum256(restored) != wantHash {
		t.Fatal("envelope round-trip SHA-256 mismatch")
	}

	// Corrupt magic
	badMagic := append([]byte(nil), frame...)
	badMagic[0] ^= 0xff
	if _, _, _, err := UnmarshalEnvelope(badMagic); err == nil {
		t.Fatal("expected error on corrupt magic")
	}

	// Unsupported version
	badVer := append([]byte(nil), frame...)
	badVer[4] = 99
	if _, _, _, err := UnmarshalEnvelope(badVer); err == nil {
		t.Fatal("expected error on unsupported version")
	}

	// Unsupported codec
	badCodec := append([]byte(nil), frame...)
	badCodec[5] = 0x7f
	if _, _, _, err := UnmarshalEnvelope(badCodec); err == nil {
		t.Fatal("expected error on unsupported codec")
	}

	// Truncated header
	if _, _, _, err := UnmarshalEnvelope(frame[:4]); err == nil {
		t.Fatal("expected error on truncated envelope")
	}
}

func TestEncryptionBoundaryFailClosedAndTenantFloor(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	highlyCompressible := bytes.Repeat([]byte("{\"event\":\"compressible_payload\"}\n"), 4096)

	ac := NewAdaptiveCompressor(DefaultPolicy())
	// Attempt to weaken security flags via tenant override; EnforceSecurityFloor must block this.
	weakPolicy := DefaultPolicy()
	weakPolicy.FailClosedOnUnknownMode = false
	weakPolicy.BlockClientSideCiphertext = false
	weakPolicy.BlockMixedSecretOracle = false
	ac.SetTenantPolicy("untrusted-tenant", weakPolicy)

	cases := []struct {
		name       string
		mode       DataPlaneMode
		wantReason BypassReason
	}{
		{"unknown_mode_fail_closed", ModeUnknown, BypassUnknownDataPlane},
		{"hip0302_client_side_ciphertext", ModeClientSideCiphertext, BypassClientSideCiphertext},
		{"crime_breach_mixed_secret", ModeAttackerControlledMixedSecret, BypassOracleRiskMixedSecret},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := ac.Compress(ObjectContext{
				TenantID:   "untrusted-tenant",
				Bucket:     "vault",
				Key:        "obj.bin",
				DataPlane:  tc.mode,
				AccessTier: TierCold,
			}, highlyCompressible)
			if err != nil {
				t.Fatalf("Compress: %v", err)
			}
			if res.Codec != CodecNone {
				t.Fatalf("expected CodecNone for %s, got %s", tc.name, res.Codec)
			}
			if res.BypassReason != tc.wantReason {
				t.Fatalf("expected bypass reason %q, got %q", tc.wantReason, res.BypassReason)
			}
			if res.StoredSize != len(highlyCompressible) || !bytes.Equal(res.Payload, highlyCompressible) {
				t.Fatalf("payload mutated or resized under boundary bypass")
			}
			if res.Metrics.RanTrialProbe {
				t.Fatalf("boundary-blocked object should never execute trial compression probe")
			}
		})
	}
	_ = rng
}

// TestCrimeBreachLengthOracleMitigation demonstrates:
// 1. Naive compress-then-encrypt on mixed attacker+secret plaintext leaks secret bytes
//    via ciphertext length differential (CRIME/BREACH oracle).
// 2. AdaptiveCompressor with ModeAttackerControlledMixedSecret strictly bypasses
//    compression so matching and mismatching guesses produce identical ciphertext lengths.
func TestCrimeBreachLengthOracleMitigation(t *testing.T) {
	key := testAESKey("crime-oracle")
	ac := NewAdaptiveCompressor(DefaultPolicy())

	// Synthetic non-sensitive fixture marker embedded in repeated context to make
	// the payload exceed MinObjectBytes (1024 B) while keeping both guess lengths identical.
	padding := bytes.Repeat([]byte("header=ok;region=us-west-2;status=200;\n"), 40)
	secretLine := []byte("session_marker=fixture_alpha_2026_gamma\n")
	matchGuess := []byte("attacker_probe=fixture_alpha_2026_gamma\n")
	mismatchGuess := []byte("attacker_probe=zzzzzzz_omega_9999_delta\n")

	if len(matchGuess) != len(mismatchGuess) {
		t.Fatalf("guesses must have equal length: %d vs %d", len(matchGuess), len(mismatchGuess))
	}

	payloadMatch := append(append(append([]byte(nil), padding...), secretLine...), matchGuess...)
	payloadMismatch := append(append(append([]byte(nil), padding...), secretLine...), mismatchGuess...)

	// (A) Negative demonstration: naive pre-encryption compression leaks length oracle.
	naiveMatch, err := ac.CompressAndSeal(ObjectContext{
		TenantID:  "demo-naive",
		Bucket:    "mixed",
		Key:       "match",
		DataPlane: ModePreEncryptionPlaintext,
	}, payloadMatch, key)
	if err != nil {
		t.Fatalf("naive match: %v", err)
	}
	naiveMismatch, err := ac.CompressAndSeal(ObjectContext{
		TenantID:  "demo-naive",
		Bucket:    "mixed",
		Key:       "mismatch",
		DataPlane: ModePreEncryptionPlaintext,
	}, payloadMismatch, key)
	if err != nil {
		t.Fatalf("naive mismatch: %v", err)
	}
	if len(naiveMatch.Ciphertext) >= len(naiveMismatch.Ciphertext) {
		t.Fatalf("expected naive compress-then-encrypt to exhibit CRIME/BREACH length leak (match=%d, mismatch=%d)",
			len(naiveMatch.Ciphertext), len(naiveMismatch.Ciphertext))
	}

	// (B) Protected regime: ModeAttackerControlledMixedSecret completely eliminates oracle.
	safeMatch, err := ac.CompressAndSeal(ObjectContext{
		TenantID:  "demo-protected",
		Bucket:    "mixed",
		Key:       "match",
		DataPlane: ModeAttackerControlledMixedSecret,
	}, payloadMatch, key)
	if err != nil {
		t.Fatalf("safe match: %v", err)
	}
	safeMismatch, err := ac.CompressAndSeal(ObjectContext{
		TenantID:  "demo-protected",
		Bucket:    "mixed",
		Key:       "mismatch",
		DataPlane: ModeAttackerControlledMixedSecret,
	}, payloadMismatch, key)
	if err != nil {
		t.Fatalf("safe mismatch: %v", err)
	}

	if safeMatch.Codec != CodecNone || safeMismatch.Codec != CodecNone {
		t.Fatalf("expected CodecNone under mixed-secret oracle guard, got %s and %s", safeMatch.Codec, safeMismatch.Codec)
	}
	if safeMatch.BypassReason != BypassOracleRiskMixedSecret || safeMismatch.BypassReason != BypassOracleRiskMixedSecret {
		t.Fatalf("unexpected bypass reasons: %q, %q", safeMatch.BypassReason, safeMismatch.BypassReason)
	}
	if len(safeMatch.Ciphertext) != len(safeMismatch.Ciphertext) {
		t.Fatalf("oracle leak remained: safeMatch=%d != safeMismatch=%d", len(safeMatch.Ciphertext), len(safeMismatch.Ciphertext))
	}

	// Verify exact decryption fidelity.
	gotMatch, err := ac.OpenAndDecompress(safeMatch, key)
	if err != nil || !bytes.Equal(gotMatch, payloadMatch) {
		t.Fatalf("OpenAndDecompress mismatch on protected payload: err=%v", err)
	}
}

// TestPostEncryptionNegativeControl verifies why compression MUST NEVER occur
// after AES-256-GCM encryption: ciphertext has near-8.0 bits/byte entropy,
// produces ~1.00x (or expanded) output under gzip/zstd, and is caught by both
// the DataPlane boundary guard AND the fallback Shannon entropy probe.
func TestPostEncryptionNegativeControl(t *testing.T) {
	rng := rand.New(rand.NewSource(123))
	key := testAESKey("post-enc-neg")
	plain := makeJSONLogsPayload(rng, 256*1024)
	ct := encryptDeterministicGCM(key, 1, plain)

	// Attempting gzip and zstd directly on ciphertext fails to compress (ratio <= 1.00x).
	gzOut, err := CompressWithCodec(CodecGzip, ct)
	if err != nil {
		t.Fatalf("gzip ciphertext: %v", err)
	}
	zstdOut, err := CompressWithCodec(CodecZstd, ct)
	if err != nil {
		t.Fatalf("zstd ciphertext: %v", err)
	}
	if len(gzOut) < len(ct) || len(zstdOut) < len(ct) {
		t.Fatalf("AES-GCM ciphertext unexpectedly shrank: orig=%d gzip=%d zstd=%d", len(ct), len(gzOut), len(zstdOut))
	}

	ac := NewAdaptiveCompressor(DefaultPolicy())

	// Path 1: Properly tagged ModeClientSideCiphertext bypasses at boundary gate (zero probe work).
	resTagged, err := ac.Compress(ObjectContext{
		TenantID:  "tenant-zk",
		Bucket:    "zk",
		Key:       "enc.bin",
		DataPlane: ModeClientSideCiphertext,
	}, ct)
	if err != nil {
		t.Fatalf("tagged ciphertext compress: %v", err)
	}
	if resTagged.Codec != CodecNone || resTagged.BypassReason != BypassClientSideCiphertext {
		t.Fatalf("expected BypassClientSideCiphertext, got codec=%s reason=%s", resTagged.Codec, resTagged.BypassReason)
	}

	// Path 2: Defense-in-depth — even if ciphertext were mislabeled as plaintext,
	// the Shannon entropy screen catches it in microseconds without trial compression.
	resMislabeled, err := ac.Compress(ObjectContext{
		TenantID:  "tenant-zk",
		Bucket:    "zk",
		Key:       "mislabeled.bin",
		DataPlane: ModePreEncryptionPlaintext,
	}, ct)
	if err != nil {
		t.Fatalf("mislabeled ciphertext compress: %v", err)
	}
	if resMislabeled.Codec != CodecNone || resMislabeled.BypassReason != BypassHighEntropy {
		t.Fatalf("expected BypassHighEntropy defense-in-depth, got codec=%s reason=%s entropy=%.3f",
			resMislabeled.Codec, resMislabeled.BypassReason, resMislabeled.Metrics.ShannonEntropy)
	}
	if resMislabeled.Metrics.RanTrialProbe {
		t.Fatal("high-entropy ciphertext should be rejected before trial compression")
	}
	if resMislabeled.Metrics.ShannonEntropy < 7.80 {
		t.Fatalf("expected ciphertext entropy >= 7.80 bits/byte, got %.3f", resMislabeled.Metrics.ShannonEntropy)
	}
}

func TestHighEntropyPrecompressedAndTinyFastBypass(t *testing.T) {
	rng := rand.New(rand.NewSource(777))
	ac := NewAdaptiveCompressor(DefaultPolicy())

	// 1. Pre-compressed JPEG blob -> O(1) magic bypass
	jpeg := makeJPEGBlobPayload(rng, 256*1024)
	resJPEG, err := ac.Compress(ObjectContext{
		TenantID:  "t-media",
		Bucket:    "media",
		Key:       "img.jpg",
		DataPlane: ModePreEncryptionPlaintext,
	}, jpeg)
	if err != nil {
		t.Fatalf("jpeg compress: %v", err)
	}
	if resJPEG.Codec != CodecNone || resJPEG.BypassReason != BypassPrecompressedMagic || resJPEG.Metrics.PrecompressedType != "jpeg" {
		t.Fatalf("unexpected jpeg result: codec=%s reason=%s kind=%s", resJPEG.Codec, resJPEG.BypassReason, resJPEG.Metrics.PrecompressedType)
	}

	// 2. Raw high-entropy random blob (no magic header) -> Shannon entropy fast bypass
	rawRand := makePseudoRandomLowRedundancy(rng, 256*1024)
	resRand, err := ac.Compress(ObjectContext{
		TenantID:  "t-rand",
		Bucket:    "raw",
		Key:       "blob.bin",
		DataPlane: ModePreEncryptionPlaintext,
	}, rawRand)
	if err != nil {
		t.Fatalf("random compress: %v", err)
	}
	if resRand.Codec != CodecNone || resRand.BypassReason != BypassHighEntropy {
		t.Fatalf("expected BypassHighEntropy, got codec=%s reason=%s", resRand.Codec, resRand.BypassReason)
	}
	if resRand.ProbeDuration > 5*time.Millisecond {
		t.Fatalf("entropy bypass took too long: %v", resRand.ProbeDuration)
	}

	// 3. Tiny metadata object (< 1024 B) -> BypassBelowMinSize
	tiny := makeTinyMetadataPayload(rng, 1, 256)
	resTiny, err := ac.Compress(ObjectContext{
		TenantID:  "t-meta",
		Bucket:    "meta",
		Key:       "item.json",
		DataPlane: ModePreEncryptionPlaintext,
	}, tiny)
	if err != nil {
		t.Fatalf("tiny compress: %v", err)
	}
	if resTiny.Codec != CodecNone || resTiny.BypassReason != BypassBelowMinSize {
		t.Fatalf("expected BypassBelowMinSize, got codec=%s reason=%s", resTiny.Codec, resTiny.BypassReason)
	}
}

func TestCompressAndSealAADTamperDetection(t *testing.T) {
	rng := rand.New(rand.NewSource(314))
	key := testAESKey("aad-integrity")
	ac := NewAdaptiveCompressor(DefaultPolicy())
	plain := makeJSONLogsPayload(rng, 64*1024)

	sealed, err := ac.CompressAndSeal(ObjectContext{
		TenantID:   "tenant-a",
		Bucket:     "prod-logs",
		Key:        "shard-01.json",
		DataPlane:  ModePreEncryptionPlaintext,
		AccessTier: TierWarm,
	}, plain, key)
	if err != nil {
		t.Fatalf("CompressAndSeal: %v", err)
	}
	if sealed.Codec != CodecZstd {
		t.Fatalf("expected CodecZstd, got %s", sealed.Codec)
	}
	if len(sealed.Ciphertext) >= len(plain) {
		t.Fatalf("expected pre-encryption compression to shrink ciphertext: plain=%d ct=%d", len(plain), len(sealed.Ciphertext))
	}

	restored, err := ac.OpenAndDecompress(sealed, key)
	if err != nil {
		t.Fatalf("OpenAndDecompress: %v", err)
	}
	if sha256.Sum256(restored) != sha256.Sum256(plain) {
		t.Fatal("sealed round-trip SHA-256 mismatch")
	}

	// Tampering with codec metadata must fail AES-GCM authentication.
	tamperedCodec := sealed
	tamperedCodec.Codec = CodecNone
	if _, err := ac.OpenAndDecompress(tamperedCodec, key); err == nil {
		t.Fatal("expected AES-GCM auth failure on tampered Codec AAD")
	}

	// Cross-tenant replay tampering must fail AES-GCM authentication.
	tamperedTenant := sealed
	tamperedTenant.TenantID = "tenant-b"
	if _, err := ac.OpenAndDecompress(tamperedTenant, key); err == nil {
		t.Fatal("expected AES-GCM auth failure on tampered TenantID AAD")
	}
}

func TestConcurrentMultiTenantRaceSafety(t *testing.T) {
	ac := NewAdaptiveCompressor(DefaultPolicy())
	corpus := GenerateBenchmarkCorpus(2026, 4)

	// Configure distinct per-tenant policies.
	hotBudget := DefaultPolicy()
	hotBudget.CPUBudgetPerMiB = 1 * time.Millisecond
	ac.SetTenantPolicy("tenant-analytics", hotBudget)

	disabled := DefaultPolicy()
	disabled.Enabled = false
	ac.SetTenantPolicy("tenant-disabled", disabled)

	key := testAESKey("concurrent-race")
	const workers = 16
	const iters = 8

	var wg sync.WaitGroup
	errCh := make(chan error, workers*iters)

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for it := 0; it < iters; it++ {
				if workerID == 0 && it%2 == 0 {
					// Concurrent dynamic policy update under -race.
					p := DefaultPolicy()
					p.MinSavingsRatio = 0.12
					ac.SetTenantPolicy("tenant-logs", p)
				}
				obj := corpus[(workerID*iters+it)%len(corpus)]
				sealed, err := ac.CompressAndSeal(obj.Context, obj.Data, key)
				if err != nil {
					errCh <- fmt.Errorf("worker %d compress/seal %s: %w", workerID, obj.Context.Key, err)
					return
				}
				if obj.Context.DataPlane == ModeClientSideCiphertext && sealed.Codec != CodecNone {
					errCh <- fmt.Errorf("worker %d compressed ZK ciphertext %s", workerID, obj.Context.Key)
					return
				}
				restored, err := ac.OpenAndDecompress(sealed, key)
				if err != nil {
					errCh <- fmt.Errorf("worker %d open/decompress %s: %w", workerID, obj.Context.Key, err)
					return
				}
				if sha256.Sum256(restored) != obj.SHA256 {
					errCh <- fmt.Errorf("worker %d SHA-256 mismatch for %s", workerID, obj.Context.Key)
					return
				}
			}
		}(w)
	}

	wg.Wait()
	close(errCh)
	for err := range errCh {
		t.Fatal(err)
	}

	stats := ac.Stats()
	if stats.ObjectsTotal != int64(workers*iters) {
		t.Fatalf("expected %d total objects in stats, got %d", workers*iters, stats.ObjectsTotal)
	}
	if stats.StoredBytesTotal > stats.OriginalBytesTotal {
		t.Fatalf("stored bytes (%d) exceeded original bytes (%d)", stats.StoredBytesTotal, stats.OriginalBytesTotal)
	}
}

func TestDeterministicCorpusReproducibility(t *testing.T) {
	c1 := GenerateBenchmarkCorpus(20261125, 4)
	c2 := GenerateBenchmarkCorpus(20261125, 4)
	if len(c1) != len(c2) || len(c1) != 24 {
		t.Fatalf("unexpected corpus lengths: %d vs %d", len(c1), len(c2))
	}
	for i := range c1 {
		if c1[i].SHA256 != c2[i].SHA256 {
			t.Fatalf("corpus index %d not deterministic across identical seeds", i)
		}
	}
}
