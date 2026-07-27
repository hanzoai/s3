package lux

import (
	"bytes"
	"context"
	"encoding/hex"
	"testing"

	"github.com/hanzoai/s3/s3/kms"
)

// withMaster returns a provider whose master key is already cached, so the
// envelope can be exercised without a live kmsd.
func withMaster(t *testing.T, keyID, masterHex string) *Provider {
	t.Helper()
	m, err := hex.DecodeString(masterHex)
	if err != nil {
		t.Fatal(err)
	}
	return &Provider{env: "test", masters: map[string][]byte{keyID: m}}
}

const testMaster = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

func TestEnvelopeRoundTrip(t *testing.T) {
	p := withMaster(t, "bucket-key", testMaster)
	ctx := context.Background()
	encCtx := kms.BuildS3EncryptionContext("photos", "cat.jpg", false)

	gen, err := p.GenerateDataKey(ctx, &kms.GenerateDataKeyRequest{
		KeyID: "bucket-key", KeySpec: kms.KeySpecAES256, EncryptionContext: encCtx,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(gen.Plaintext) != 32 {
		t.Fatalf("data key is %d bytes, want 32", len(gen.Plaintext))
	}
	if bytes.Contains(gen.CiphertextBlob, gen.Plaintext) {
		t.Fatal("the sealed blob contains the plaintext data key")
	}

	got, err := p.Decrypt(ctx, &kms.DecryptRequest{
		CiphertextBlob: gen.CiphertextBlob, EncryptionContext: encCtx,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got.Plaintext, gen.Plaintext) {
		t.Fatal("round-trip returned a different data key")
	}
	if got.KeyID != "bucket-key" {
		t.Fatalf("KeyID = %q, want bucket-key", got.KeyID)
	}
}

// A blob sealed for one object must not open under another's context — that
// binding is the whole point of passing the context through the derivation.
func TestDecryptRefusesForeignContext(t *testing.T) {
	p := withMaster(t, "bucket-key", testMaster)
	ctx := context.Background()

	gen, err := p.GenerateDataKey(ctx, &kms.GenerateDataKeyRequest{
		KeyID:             "bucket-key",
		KeySpec:           kms.KeySpecAES256,
		EncryptionContext: kms.BuildS3EncryptionContext("photos", "cat.jpg", false),
	})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := p.Decrypt(ctx, &kms.DecryptRequest{
		CiphertextBlob:    gen.CiphertextBlob,
		EncryptionContext: kms.BuildS3EncryptionContext("photos", "dog.jpg", false),
	}); err == nil {
		t.Fatal("a blob opened under a different object's context")
	}
}

func TestDecryptRejectsMalformedBlob(t *testing.T) {
	p := withMaster(t, "bucket-key", testMaster)
	for name, blob := range map[string][]byte{
		"empty":     {},
		"truncated": {40, 1, 2, 3},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := p.Decrypt(context.Background(), &kms.DecryptRequest{CiphertextBlob: blob}); err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}

// Context ordering must not change the derived key, or a caller iterating a map
// would get a different result each run.
func TestKEKIsOrderIndependent(t *testing.T) {
	m, _ := hex.DecodeString(testMaster)
	a, err := kek(m, "k", map[string]string{"x": "1", "y": "2"})
	if err != nil {
		t.Fatal(err)
	}
	b, err := kek(m, "k", map[string]string{"y": "2", "x": "1"})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(a, b) {
		t.Fatal("kek depends on map iteration order")
	}
}
