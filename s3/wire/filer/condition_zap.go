// Code generated from filer.proto; DO NOT EDIT.

package filerwire

import (
	zap "github.com/zap-proto/go"
)

// --- WriteConditionClause (WriteCondition.Clause nested message) ---

const (
	writeConditionClauseKindOff      = 0  // kind (field 1, enum WriteCondition.Kind)
	writeConditionClauseEtagsOff     = 8  // etags (field 2, repeated string)
	writeConditionClauseUnixTimeOff  = 16 // unix_time (field 3, int64)
	writeConditionClauseAllowWeakOff = 24 // allow_weak (field 4, bool)
	writeConditionClauseExtKeyOff    = 32 // ext_key (field 5, string)
	writeConditionClauseExtValueOff  = 40 // ext_value (field 6, string)
	writeConditionClauseGateKeyOff   = 48 // gate_key (field 7, string)
	writeConditionClauseGateValueOff = 56 // gate_value (field 8, string)
	writeConditionClauseSize         = 64
)

// WriteConditionClause is a zero-copy view into one encoded WriteCondition.Clause.
type WriteConditionClause struct{ o zap.Object }

// WrapWriteConditionClause parses b and returns a typed view.
func WrapWriteConditionClause(b []byte) (WriteConditionClause, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return WriteConditionClause{}, err
	}
	return WriteConditionClause{o: m.Root()}, nil
}

// Kind reads the kind enum (see WriteConditionKind* constants).
func (t WriteConditionClause) Kind() uint32 { return t.o.Uint32(writeConditionClauseKindOff) }
func (t WriteConditionClause) UnixTime() int64 {
	return t.o.Int64(writeConditionClauseUnixTimeOff)
}
func (t WriteConditionClause) AllowWeak() bool {
	return t.o.Bool(writeConditionClauseAllowWeakOff)
}
func (t WriteConditionClause) ExtKey() string   { return t.o.Text(writeConditionClauseExtKeyOff) }
func (t WriteConditionClause) ExtValue() string { return t.o.Text(writeConditionClauseExtValueOff) }
func (t WriteConditionClause) GateKey() string  { return t.o.Text(writeConditionClauseGateKeyOff) }
func (t WriteConditionClause) GateValue() string {
	return t.o.Text(writeConditionClauseGateValueOff)
}

// EtagsLen returns the number of etags elements (proto field 2, repeated string).
func (t WriteConditionClause) EtagsLen() int {
	return t.o.List(writeConditionClauseEtagsOff).Len()
}

// Etag returns the i-th etags element.
func (t WriteConditionClause) Etag(i int) string {
	return elemText(t.o.List(writeConditionClauseEtagsOff), i)
}

// WriteConditionClauseInput collects the field values for NewWriteConditionClause.
type WriteConditionClauseInput struct {
	Kind      uint32
	Etags     []string
	UnixTime  int64
	AllowWeak bool
	ExtKey    string
	ExtValue  string
	GateKey   string
	GateValue string
}

// NewWriteConditionClause builds a ZAP-encoded WriteCondition.Clause from in.
func NewWriteConditionClause(in WriteConditionClauseInput) []byte {
	b := zap.NewBuilder(256)
	etagsOff, etagsLen := addStringList(b, in.Etags)
	ob := b.StartObject(writeConditionClauseSize)
	ob.SetUint32(writeConditionClauseKindOff, in.Kind)
	ob.SetList(writeConditionClauseEtagsOff, etagsOff, etagsLen)
	ob.SetInt64(writeConditionClauseUnixTimeOff, in.UnixTime)
	ob.SetBool(writeConditionClauseAllowWeakOff, in.AllowWeak)
	ob.SetText(writeConditionClauseExtKeyOff, in.ExtKey)
	ob.SetText(writeConditionClauseExtValueOff, in.ExtValue)
	ob.SetText(writeConditionClauseGateKeyOff, in.GateKey)
	ob.SetText(writeConditionClauseGateValueOff, in.GateValue)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- WriteCondition ---

const (
	writeConditionClausesOff = 0 // clauses (field 1, repeated Clause)
	writeConditionSize       = 8
)

// WriteCondition is a zero-copy view into a ZAP-encoded WriteCondition message.
type WriteCondition struct{ o zap.Object }

// WrapWriteCondition parses b and returns a typed view.
func WrapWriteCondition(b []byte) (WriteCondition, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return WriteCondition{}, err
	}
	return WriteCondition{o: m.Root()}, nil
}

// ClausesLen returns the number of clauses elements (proto field 1, repeated Clause).
func (t WriteCondition) ClausesLen() int { return t.o.List(writeConditionClausesOff).Len() }

// Clauses returns the i-th clauses element as a typed WriteConditionClause view.
func (t WriteCondition) Clauses(i int) WriteConditionClause {
	return WriteConditionClause{o: t.o.List(writeConditionClausesOff).ObjectAt(i)}
}

// WriteConditionInput collects the field values for NewWriteCondition. Each
// Clauses element is a standalone WriteConditionClause buffer (from
// NewWriteConditionClause).
type WriteConditionInput struct {
	Clauses [][]byte
}

// NewWriteCondition builds a ZAP-encoded WriteCondition message from in.
func NewWriteCondition(in WriteConditionInput) []byte {
	b := zap.NewBuilder(256)
	clausesOff, clausesLen := addObjectList(b, in.Clauses)
	ob := b.StartObject(writeConditionSize)
	ob.SetList(writeConditionClausesOff, clausesOff, clausesLen)
	ob.FinishAsRoot()
	return b.Finish()
}
