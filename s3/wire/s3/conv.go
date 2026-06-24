// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// Hand-written ZAP boundary helpers for the s3 bucket-metadata and
// circuit-breaker config messages. The generated *_zap.go files expose the
// low-level builders (New*) and zero-copy views (Wrap*); this companion owns
// the round-trip to plain Go values (maps, slices, scalars) so callers in
// s3/s3api and s3/shell never touch ZAP offsets or protobuf. One concern per
// seam: the generated code knows field layout, this file knows the in-process
// shapes the retired protobuf structs used to carry.

package s3wire

import zap "github.com/zap-proto/go"

// --- EncryptionConfiguration ---

// EncryptionFields is the native value carried by a BucketMetadata.Encryption.
type EncryptionFields struct {
	SseAlgorithm     string
	KmsKeyId         string
	BucketKeyEnabled bool
}

// Fields reads the three scalar fields of an EncryptionConfiguration view.
func (t EncryptionConfiguration) Fields() EncryptionFields {
	return EncryptionFields{
		SseAlgorithm:     t.SseAlgorithm(),
		KmsKeyId:         t.KmsKeyId(),
		BucketKeyEnabled: t.BucketKeyEnabled(),
	}
}

// --- CORS ---

// CORSRuleFields is the native value for one CORSRule.
type CORSRuleFields struct {
	AllowedHeaders []string
	AllowedMethods []string
	AllowedOrigins []string
	ExposeHeaders  []string
	MaxAgeSeconds  int32
	Id             string
}

func bytesList(ss []string) [][]byte {
	if len(ss) == 0 {
		return nil
	}
	out := make([][]byte, len(ss))
	for i, s := range ss {
		out[i] = []byte(s)
	}
	return out
}

func stringList(l zap.List) []string {
	n := l.Length()
	if n == 0 {
		return nil
	}
	out := make([]string, n)
	for i := 0; i < n; i++ {
		out[i] = string(l.BytesAt(i))
	}
	return out
}

// NewCORSRuleFromFields builds one CORSRule sub-buffer.
func NewCORSRuleFromFields(r CORSRuleFields) []byte {
	return NewCORSRule(CORSRuleInput{
		AllowedHeaders: bytesList(r.AllowedHeaders),
		AllowedMethods: bytesList(r.AllowedMethods),
		AllowedOrigins: bytesList(r.AllowedOrigins),
		ExposeHeaders:  bytesList(r.ExposeHeaders),
		MaxAgeSeconds:  r.MaxAgeSeconds,
		Id:             r.Id,
	})
}

// ruleAt reads CORSRule i out of a CORSConfiguration's rule list.
func (t CORSConfiguration) ruleFields(i int) CORSRuleFields {
	o := t.CorsRules().ObjectAt(i)
	return CORSRuleFields{
		AllowedHeaders: stringList(o.List(cORSRuleAllowedHeadersOff)),
		AllowedMethods: stringList(o.List(cORSRuleAllowedMethodsOff)),
		AllowedOrigins: stringList(o.List(cORSRuleAllowedOriginsOff)),
		ExposeHeaders:  stringList(o.List(cORSRuleExposeHeadersOff)),
		MaxAgeSeconds:  o.Int32(cORSRuleMaxAgeSecondsOff),
		Id:             o.Text(cORSRuleIdOff),
	}
}

// Rules reads every CORSRule of the configuration as native values.
func (t CORSConfiguration) Rules() []CORSRuleFields {
	n := t.CorsRules().Length()
	if n == 0 {
		return nil
	}
	out := make([]CORSRuleFields, n)
	for i := 0; i < n; i++ {
		out[i] = t.ruleFields(i)
	}
	return out
}

// --- BucketMetadata ---

// BucketMetadataFields is the native value of a whole BucketMetadata message.
// A nil Encryption / nil Cors means the corresponding field is absent.
type BucketMetadataFields struct {
	Tags       map[string]string
	Cors       []CORSRuleFields
	HasCors    bool
	Encryption *EncryptionFields
}

// EncodeBucketMetadata builds a ZAP-encoded BucketMetadata from native values.
func EncodeBucketMetadata(f BucketMetadataFields) []byte {
	in := BucketMetadataInput{}

	if len(f.Tags) > 0 {
		in.Tags = make([][]byte, 0, len(f.Tags))
		for k, v := range f.Tags {
			in.Tags = append(in.Tags, NewBucketMetadataTagsEntry(BucketMetadataTagsEntryInput{Key: k, Value: v}))
		}
	}

	if f.HasCors {
		rules := make([][]byte, len(f.Cors))
		for i, r := range f.Cors {
			rules[i] = NewCORSRuleFromFields(r)
		}
		in.Cors = NewCORSConfiguration(CORSConfigurationInput{CorsRules: rules})
	}

	if f.Encryption != nil {
		in.Encryption = NewEncryptionConfiguration(EncryptionConfigurationInput{
			SseAlgorithm:     f.Encryption.SseAlgorithm,
			KmsKeyId:         f.Encryption.KmsKeyId,
			BucketKeyEnabled: f.Encryption.BucketKeyEnabled,
		})
	}

	return NewBucketMetadata(in)
}

// DecodeBucketMetadata parses a ZAP-encoded BucketMetadata into native values.
func DecodeBucketMetadata(b []byte) (BucketMetadataFields, error) {
	m, err := WrapBucketMetadata(b)
	if err != nil {
		return BucketMetadataFields{}, err
	}
	var f BucketMetadataFields

	tags := m.Tags()
	if n := tags.Length(); n > 0 {
		f.Tags = make(map[string]string, n)
		for i := 0; i < n; i++ {
			e := BucketMetadataTagsEntry{o: tags.ObjectAt(i)}
			f.Tags[e.Key()] = e.Value()
		}
	}

	if cors := m.Cors(); !cors.o.IsNull() {
		f.HasCors = true
		f.Cors = cors.Rules()
	}

	if enc := m.Encryption(); !enc.o.IsNull() {
		ef := enc.Fields()
		f.Encryption = &ef
	}

	return f, nil
}

// --- S3CircuitBreakerConfig ---

// CircuitBreakerOptionsFields is the native value of one S3CircuitBreakerOptions.
type CircuitBreakerOptionsFields struct {
	Enabled bool             `json:"enabled,omitempty"`
	Actions map[string]int64 `json:"actions,omitempty"`
}

// CircuitBreakerConfigFields is the native value of an S3CircuitBreakerConfig.
// A nil Global means no global options. Buckets maps bucket name to its options.
type CircuitBreakerConfigFields struct {
	Global  *CircuitBreakerOptionsFields            `json:"global,omitempty"`
	Buckets map[string]*CircuitBreakerOptionsFields `json:"buckets,omitempty"`
}

func encodeCBOptions(o *CircuitBreakerOptionsFields) []byte {
	if o == nil {
		return nil
	}
	in := S3CircuitBreakerOptionsInput{Enabled: o.Enabled}
	if len(o.Actions) > 0 {
		in.Actions = make([][]byte, 0, len(o.Actions))
		for k, v := range o.Actions {
			in.Actions = append(in.Actions, NewS3CircuitBreakerOptionsActionsEntry(S3CircuitBreakerOptionsActionsEntryInput{Key: k, Value: v}))
		}
	}
	return NewS3CircuitBreakerOptions(in)
}

func decodeCBOptions(o S3CircuitBreakerOptions) *CircuitBreakerOptionsFields {
	if o.o.IsNull() {
		return nil
	}
	out := &CircuitBreakerOptionsFields{Enabled: o.Enabled()}
	acts := o.Actions()
	if n := acts.Length(); n > 0 {
		out.Actions = make(map[string]int64, n)
		for i := 0; i < n; i++ {
			e := S3CircuitBreakerOptionsActionsEntry{o: acts.ObjectAt(i)}
			out.Actions[e.Key()] = e.Value()
		}
	}
	return out
}

// EncodeCircuitBreakerConfig builds a ZAP-encoded S3CircuitBreakerConfig.
func EncodeCircuitBreakerConfig(f CircuitBreakerConfigFields) []byte {
	in := S3CircuitBreakerConfigInput{Global: encodeCBOptions(f.Global)}
	if len(f.Buckets) > 0 {
		in.Buckets = make([][]byte, 0, len(f.Buckets))
		for k, v := range f.Buckets {
			in.Buckets = append(in.Buckets, NewS3CircuitBreakerConfigBucketsEntry(S3CircuitBreakerConfigBucketsEntryInput{
				Key:   k,
				Value: encodeCBOptions(v),
			}))
		}
	}
	return NewS3CircuitBreakerConfig(in)
}

// DecodeCircuitBreakerConfig parses a ZAP-encoded S3CircuitBreakerConfig.
func DecodeCircuitBreakerConfig(b []byte) (CircuitBreakerConfigFields, error) {
	m, err := WrapS3CircuitBreakerConfig(b)
	if err != nil {
		return CircuitBreakerConfigFields{}, err
	}
	f := CircuitBreakerConfigFields{Global: decodeCBOptions(m.Global())}

	buckets := m.Buckets()
	if n := buckets.Length(); n > 0 {
		f.Buckets = make(map[string]*CircuitBreakerOptionsFields, n)
		for i := 0; i < n; i++ {
			e := S3CircuitBreakerConfigBucketsEntry{o: buckets.ObjectAt(i)}
			f.Buckets[e.Key()] = decodeCBOptions(e.Value())
		}
	}
	return f, nil
}
