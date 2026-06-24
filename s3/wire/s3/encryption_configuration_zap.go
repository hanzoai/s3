// Code generated from s3.proto; DO NOT EDIT.

package s3wire

import (
	zap "github.com/zap-proto/go"
)

// --- EncryptionConfiguration ---

const (
	encryptionConfigurationSseAlgorithmOff     = 0
	encryptionConfigurationKmsKeyIdOff         = 8
	encryptionConfigurationBucketKeyEnabledOff = 16
	encryptionConfigurationSize                = 17
)

// EncryptionConfiguration is a zero-copy view into a ZAP-encoded
// EncryptionConfiguration message.
type EncryptionConfiguration struct{ o zap.Object }

// WrapEncryptionConfiguration parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapEncryptionConfiguration(b []byte) (EncryptionConfiguration, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return EncryptionConfiguration{}, err
	}
	return EncryptionConfiguration{o: m.Root()}, nil
}

// SseAlgorithm reads the sse_algorithm field (proto field 1, string).
func (t EncryptionConfiguration) SseAlgorithm() string {
	return t.o.Text(encryptionConfigurationSseAlgorithmOff)
}

// KmsKeyId reads the kms_key_id field (proto field 2, string).
func (t EncryptionConfiguration) KmsKeyId() string {
	return t.o.Text(encryptionConfigurationKmsKeyIdOff)
}

// BucketKeyEnabled reads the bucket_key_enabled field (proto field 3, bool).
func (t EncryptionConfiguration) BucketKeyEnabled() bool {
	return t.o.Bool(encryptionConfigurationBucketKeyEnabledOff)
}

// EncryptionConfigurationInput collects the field values for
// NewEncryptionConfiguration.
type EncryptionConfigurationInput struct {
	SseAlgorithm     string
	KmsKeyId         string
	BucketKeyEnabled bool
}

// NewEncryptionConfiguration builds a ZAP-encoded EncryptionConfiguration
// message from in and returns the bytes.
func NewEncryptionConfiguration(in EncryptionConfigurationInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(encryptionConfigurationSize)
	ob.SetText(encryptionConfigurationSseAlgorithmOff, in.SseAlgorithm)
	ob.SetText(encryptionConfigurationKmsKeyIdOff, in.KmsKeyId)
	ob.SetBool(encryptionConfigurationBucketKeyEnabledOff, in.BucketKeyEnabled)
	ob.FinishAsRoot()
	return b.Finish()
}
