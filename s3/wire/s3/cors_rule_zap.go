// Code generated from s3.proto; DO NOT EDIT.

package s3wire

import (
	zap "github.com/zap-proto/go"
)

// --- CORSRule ---

const (
	cORSRuleAllowedHeadersOff = 0
	cORSRuleAllowedMethodsOff = 8
	cORSRuleAllowedOriginsOff = 16
	cORSRuleExposeHeadersOff  = 24
	cORSRuleMaxAgeSecondsOff  = 32
	cORSRuleIdOff             = 36
	cORSRuleSize              = 44
)

// CORSRule is a zero-copy view into a ZAP-encoded CORSRule message.
type CORSRule struct{ o zap.Object }

// WrapCORSRule parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapCORSRule(b []byte) (CORSRule, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return CORSRule{}, err
	}
	return CORSRule{o: m.Root()}, nil
}

// AllowedHeaders reads the allowed_headers field (proto field 1, repeated
// string). Each element is an out-of-line entry; read element i with
// List.BytesAt(i).
func (t CORSRule) AllowedHeaders() zap.List { return t.o.List(cORSRuleAllowedHeadersOff) }

// AllowedMethods reads the allowed_methods field (proto field 2, repeated
// string). Each element is an out-of-line entry; read element i with
// List.BytesAt(i).
func (t CORSRule) AllowedMethods() zap.List { return t.o.List(cORSRuleAllowedMethodsOff) }

// AllowedOrigins reads the allowed_origins field (proto field 3, repeated
// string). Each element is an out-of-line entry; read element i with
// List.BytesAt(i).
func (t CORSRule) AllowedOrigins() zap.List { return t.o.List(cORSRuleAllowedOriginsOff) }

// ExposeHeaders reads the expose_headers field (proto field 4, repeated
// string). Each element is an out-of-line entry; read element i with
// List.BytesAt(i).
func (t CORSRule) ExposeHeaders() zap.List { return t.o.List(cORSRuleExposeHeadersOff) }

// MaxAgeSeconds reads the max_age_seconds field (proto field 5, int32).
func (t CORSRule) MaxAgeSeconds() int32 { return t.o.Int32(cORSRuleMaxAgeSecondsOff) }

// Id reads the id field (proto field 6, string).
func (t CORSRule) Id() string { return t.o.Text(cORSRuleIdOff) }

// CORSRuleInput collects the field values for NewCORSRule. The repeated string
// fields take one []byte (raw UTF-8) per element.
type CORSRuleInput struct {
	AllowedHeaders [][]byte
	AllowedMethods [][]byte
	AllowedOrigins [][]byte
	ExposeHeaders  [][]byte
	MaxAgeSeconds  int32
	Id             string
}

// NewCORSRule builds a ZAP-encoded CORSRule message from in and returns the
// bytes.
func NewCORSRule(in CORSRuleInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(cORSRuleSize)
	ob.SetInt32(cORSRuleMaxAgeSecondsOff, in.MaxAgeSeconds)
	ob.SetText(cORSRuleIdOff, in.Id)
	allowedHeadersLB := b.StartList(0)
	for _, elem := range in.AllowedHeaders {
		allowedHeadersLB.AddObjectBytes(elem)
	}
	ob.SetList(cORSRuleAllowedHeadersOff, allowedHeadersLB.FinishOffset(), len(in.AllowedHeaders))
	allowedMethodsLB := b.StartList(0)
	for _, elem := range in.AllowedMethods {
		allowedMethodsLB.AddObjectBytes(elem)
	}
	ob.SetList(cORSRuleAllowedMethodsOff, allowedMethodsLB.FinishOffset(), len(in.AllowedMethods))
	allowedOriginsLB := b.StartList(0)
	for _, elem := range in.AllowedOrigins {
		allowedOriginsLB.AddObjectBytes(elem)
	}
	ob.SetList(cORSRuleAllowedOriginsOff, allowedOriginsLB.FinishOffset(), len(in.AllowedOrigins))
	exposeHeadersLB := b.StartList(0)
	for _, elem := range in.ExposeHeaders {
		exposeHeadersLB.AddObjectBytes(elem)
	}
	ob.SetList(cORSRuleExposeHeadersOff, exposeHeadersLB.FinishOffset(), len(in.ExposeHeaders))
	ob.FinishAsRoot()
	return b.Finish()
}
