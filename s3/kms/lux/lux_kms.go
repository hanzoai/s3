// Package lux is the KMS provider backed by luxfi/kms — the one place Hanzo and
// Lux keep secrets.
//
// luxfi/kms stores named secrets; it does not wrap data keys. So the master key
// lives there and the envelope is built here: each GenerateDataKey mints a
// random DEK and seals it under a KEK derived from the master, the key id and
// the caller's encryption context. The plaintext master never lands on disk and
// the sealed blob is useless without both the master and the same context.
//
// The client speaks ZAP. No OTLP, no gRPC, no cloud SDK.
package lux

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"sync"

	"golang.org/x/crypto/hkdf"

	"github.com/hanzoai/s3/s3/kms"
	"github.com/hanzoai/s3/s3/util"
	"github.com/luxfi/kms/pkg/zapclient"

	"crypto/sha256"
)

func init() {
	kms.RegisterProvider("lux", NewLuxKMSProvider)
}

// Provider resolves master key material from luxfi/kms and performs envelope
// operations locally.
type Provider struct {
	client *zapclient.Client
	env    string

	mu      sync.RWMutex
	masters map[string][]byte // keyID -> master key material
}

var _ kms.KMSProvider = (*Provider)(nil)

// NewLuxKMSProvider dials luxfi/kms over ZAP.
//
// Config:
//
//	kms.lux.address  ZAP address of kmsd (default kms.lux.svc:9650)
//	kms.lux.path     secret path prefix (default "/s3")
//	kms.lux.env      secret environment (default "prod")
func NewLuxKMSProvider(config util.Configuration) (kms.KMSProvider, error) {
	addr := config.GetString("address")
	if addr == "" {
		addr = "kms.lux.svc:9650"
	}
	env := config.GetString("env")
	if env == "" {
		env = "prod"
	}

	path := config.GetString("path")
	if path == "" {
		path = "/s3"
	}
	client, err := zapclient.Dial(context.Background(), addr, path)
	if err != nil {
		return nil, fmt.Errorf("lux kms: dial %s: %w", addr, err)
	}
	return &Provider{
		client:  client,
		env:     env,
		masters: make(map[string][]byte),
	}, nil
}

// master fetches and caches the master key material for keyID.
//
// The secret is hex so the value survives a text-only store byte for byte; a
// key that is not 32 bytes is refused rather than stretched, because silently
// accepting a short master would weaken every object encrypted under it.
func (p *Provider) master(ctx context.Context, keyID string) ([]byte, error) {
	p.mu.RLock()
	if m, ok := p.masters[keyID]; ok {
		p.mu.RUnlock()
		return m, nil
	}
	p.mu.RUnlock()

	v, err := p.client.Get(ctx, keyID, p.env)
	if err != nil {
		return nil, fmt.Errorf("lux kms: get key %q: %w", keyID, err)
	}
	m, err := hex.DecodeString(v)
	if err != nil {
		return nil, fmt.Errorf("lux kms: key %q is not hex: %w", keyID, err)
	}
	if len(m) != 32 {
		return nil, fmt.Errorf("lux kms: key %q is %d bytes, want 32", keyID, len(m))
	}

	p.mu.Lock()
	p.masters[keyID] = m
	p.mu.Unlock()
	return m, nil
}

// kek derives the wrapping key. Binding the key id and the encryption context
// into the derivation is what makes a blob undecryptable under a different
// context, which is the guarantee callers rely on.
func kek(master []byte, keyID string, encCtx map[string]string) ([]byte, error) {
	h := sha256.New()
	h.Write([]byte(keyID))
	for _, k := range sortedKeys(encCtx) {
		h.Write([]byte(k))
		h.Write([]byte{0})
		h.Write([]byte(encCtx[k]))
		h.Write([]byte{0})
	}
	out := make([]byte, 32)
	if _, err := io.ReadFull(hkdf.New(sha256.New, master, h.Sum(nil), []byte("s3-kms-lux-v1")), out); err != nil {
		return nil, err
	}
	return out, nil
}

func sortedKeys(m map[string]string) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	// insertion sort: encryption contexts are a handful of entries
	for i := 1; i < len(ks); i++ {
		for j := i; j > 0 && ks[j] < ks[j-1]; j-- {
			ks[j], ks[j-1] = ks[j-1], ks[j]
		}
	}
	return ks
}

// GenerateDataKey mints a DEK and returns it alongside its sealed form.
func (p *Provider) GenerateDataKey(ctx context.Context, req *kms.GenerateDataKeyRequest) (*kms.GenerateDataKeyResponse, error) {
	if req.KeySpec != kms.KeySpecAES256 && req.KeySpec != "" {
		return nil, fmt.Errorf("lux kms: unsupported key spec %q", req.KeySpec)
	}
	master, err := p.master(ctx, req.KeyID)
	if err != nil {
		return nil, err
	}
	wrapKey, err := kek(master, req.KeyID, req.EncryptionContext)
	if err != nil {
		return nil, err
	}

	dek := make([]byte, 32)
	if _, err := rand.Read(dek); err != nil {
		return nil, fmt.Errorf("lux kms: generate data key: %w", err)
	}

	gcm, err := newGCM(wrapKey)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("lux kms: nonce: %w", err)
	}
	// The blob carries the key id: the S3 encryption context holds object and
	// bucket ARNs, never the key, so Decrypt has no other way to learn which
	// master to derive from. Layout: [idLen][keyID][nonce][ciphertext].
	sealed := gcm.Seal(nonce, nonce, dek, nil)
	id := []byte(req.KeyID)
	if len(id) > 255 {
		return nil, fmt.Errorf("lux kms: key id longer than 255 bytes")
	}
	blob := make([]byte, 0, 1+len(id)+len(sealed))
	blob = append(blob, byte(len(id)))
	blob = append(blob, id...)
	blob = append(blob, sealed...)

	return &kms.GenerateDataKeyResponse{
		KeyID:          req.KeyID,
		Plaintext:      dek,
		CiphertextBlob: blob,
	}, nil
}

// Decrypt unseals a data key, reading the key id back out of the blob.
func (p *Provider) Decrypt(ctx context.Context, req *kms.DecryptRequest) (*kms.DecryptResponse, error) {
	if len(req.CiphertextBlob) < 1 {
		return nil, errors.New("lux kms: empty ciphertext")
	}
	idLen := int(req.CiphertextBlob[0])
	if len(req.CiphertextBlob) < 1+idLen {
		return nil, errors.New("lux kms: truncated ciphertext")
	}
	keyID := string(req.CiphertextBlob[1 : 1+idLen])
	sealed := req.CiphertextBlob[1+idLen:]
	master, err := p.master(ctx, keyID)
	if err != nil {
		return nil, err
	}
	wrapKey, err := kek(master, keyID, req.EncryptionContext)
	if err != nil {
		return nil, err
	}
	gcm, err := newGCM(wrapKey)
	if err != nil {
		return nil, err
	}
	if len(sealed) < gcm.NonceSize() {
		return nil, errors.New("lux kms: ciphertext too short")
	}
	nonce, ct := sealed[:gcm.NonceSize()], sealed[gcm.NonceSize():]
	dek, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return nil, fmt.Errorf("lux kms: decrypt data key: %w", err)
	}
	return &kms.DecryptResponse{KeyID: keyID, Plaintext: dek}, nil
}

// DescribeKey resolves the key, which is how a caller learns it exists.
func (p *Provider) DescribeKey(ctx context.Context, req *kms.DescribeKeyRequest) (*kms.DescribeKeyResponse, error) {
	if _, err := p.master(ctx, req.KeyID); err != nil {
		return nil, err
	}
	return &kms.DescribeKeyResponse{
		KeyID:       req.KeyID,
		ARN:         "lux:kms:" + req.KeyID,
		Description: "luxfi/kms managed key",
		KeyUsage:    kms.KeyUsageEncryptDecrypt,
		KeyState:    kms.KeyStateEnabled,
		Origin:      kms.KeyOriginExternal,
	}, nil
}

// GetKeyID is identity here: luxfi/kms addresses keys by name, so there is no
// alias layer to resolve through.
func (p *Provider) GetKeyID(_ context.Context, keyIdentifier string) (string, error) {
	if keyIdentifier == "" {
		return "", errors.New("lux kms: empty key identifier")
	}
	return keyIdentifier, nil
}

func (p *Provider) Close() error {
	p.mu.Lock()
	for k, m := range p.masters {
		for i := range m {
			m[i] = 0
		}
		delete(p.masters, k)
	}
	p.mu.Unlock()
	p.client.Close()
	return nil
}

func newGCM(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("lux kms: cipher: %w", err)
	}
	return cipher.NewGCM(block)
}
