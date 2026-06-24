// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package s3wire

import (
	"reflect"
	"testing"
)

func TestBucketMetadataRoundTrip(t *testing.T) {
	in := BucketMetadataFields{
		Tags: map[string]string{"Environment": "production", "Owner": "team-alpha"},
		Cors: []CORSRuleFields{{
			AllowedHeaders: []string{"*"},
			AllowedMethods: []string{"GET", "POST"},
			AllowedOrigins: []string{"https://example.com"},
			ExposeHeaders:  []string{"ETag"},
			MaxAgeSeconds:  3600,
			Id:             "rule-1",
		}},
		HasCors: true,
		Encryption: &EncryptionFields{
			SseAlgorithm:     "aws:kms",
			KmsKeyId:         "test-key-id",
			BucketKeyEnabled: true,
		},
	}

	out, err := DecodeBucketMetadata(EncodeBucketMetadata(in))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(out.Tags, in.Tags) {
		t.Errorf("tags: got %v want %v", out.Tags, in.Tags)
	}
	if !out.HasCors || !reflect.DeepEqual(out.Cors, in.Cors) {
		t.Errorf("cors: got %+v want %+v", out.Cors, in.Cors)
	}
	if !reflect.DeepEqual(out.Encryption, in.Encryption) {
		t.Errorf("encryption: got %+v want %+v", out.Encryption, in.Encryption)
	}
}

func TestBucketMetadataEmptyRoundTrip(t *testing.T) {
	out, err := DecodeBucketMetadata(EncodeBucketMetadata(BucketMetadataFields{}))
	if err != nil {
		t.Fatal(err)
	}
	if out.HasCors || out.Encryption != nil || len(out.Tags) != 0 {
		t.Errorf("empty metadata decoded non-empty: %+v", out)
	}
}

func TestCircuitBreakerConfigRoundTrip(t *testing.T) {
	in := CircuitBreakerConfigFields{
		Global: &CircuitBreakerOptionsFields{
			Enabled: true,
			Actions: map[string]int64{"Read:count": 500, "Write:count": 200},
		},
		Buckets: map[string]*CircuitBreakerOptionsFields{
			"b1": {Enabled: true, Actions: map[string]int64{"Read:bytes": 1048576}},
			"b2": {Enabled: false},
		},
	}

	out, err := DecodeCircuitBreakerConfig(EncodeCircuitBreakerConfig(in))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(out.Global, in.Global) {
		t.Errorf("global: got %+v want %+v", out.Global, in.Global)
	}
	if len(out.Buckets) != len(in.Buckets) {
		t.Fatalf("buckets len: got %d want %d", len(out.Buckets), len(in.Buckets))
	}
	for k, v := range in.Buckets {
		if !reflect.DeepEqual(out.Buckets[k], v) {
			t.Errorf("bucket %s: got %+v want %+v", k, out.Buckets[k], v)
		}
	}
}
