// Code generated from s3.proto; DO NOT EDIT.

package s3wire

import (
	zap "github.com/zap-proto/go"
)

// --- CORSConfiguration ---

const (
	cORSConfigurationCorsRulesOff = 0
	cORSConfigurationSize         = 8
)

// CORSConfiguration is a zero-copy view into a ZAP-encoded CORSConfiguration
// message.
type CORSConfiguration struct{ o zap.Object }

// WrapCORSConfiguration parses b and returns a typed view. Returns an error if
// the wire-level checks (magic, version, size) fail.
func WrapCORSConfiguration(b []byte) (CORSConfiguration, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return CORSConfiguration{}, err
	}
	return CORSConfiguration{o: m.Root()}, nil
}

// CorsRules reads the cors_rules field (proto field 1, repeated CORSRule). Each
// element is an out-of-line CORSRule object; read element i with
// List.ObjectAt(i) and wrap it as a CORSRule.
func (t CORSConfiguration) CorsRules() zap.List { return t.o.List(cORSConfigurationCorsRulesOff) }

// CORSConfigurationInput collects the field values for NewCORSConfiguration.
// CorsRules takes one pre-built CORSRule sub-buffer (from NewCORSRule) per
// element.
type CORSConfigurationInput struct {
	CorsRules [][]byte
}

// NewCORSConfiguration builds a ZAP-encoded CORSConfiguration message from in
// and returns the bytes.
func NewCORSConfiguration(in CORSConfigurationInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(cORSConfigurationSize)
	corsRulesLB := b.StartList(0)
	for _, elem := range in.CorsRules {
		corsRulesLB.AddObjectBytes(elem)
	}
	ob.SetList(cORSConfigurationCorsRulesOff, corsRulesLB.FinishOffset(), len(in.CorsRules))
	ob.FinishAsRoot()
	return b.Finish()
}
