// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// embed_empty_test.go locks in the embedMessage fix: a PRESENT-but-empty nested
// message (a zero-field message like ParquetInput{}) must round-trip as PRESENT,
// not collapse to a null object pointer. The pre-fix embedMessage copied the
// sub-buffer verbatim and returned base+rootOff, which for a header-only
// (size-0) child landed one past the copied region -> out of bounds -> IsNull()
// -> the variant was silently dropped on the wire. This is the volume analogue
// of the empty-embedded-message null-collapse the filer wire already fixed.

package volume_serverwire

import "testing"

// TestEmbedEmptyMessagePresenceParquet proves the only present-significant empty
// embedded message in the volume wire (InputSerialization.ParquetInput) survives
// a build -> wrap round-trip as PRESENT, and that a genuinely-unset variant stays
// absent (no false presence). It also confirms the empty variant does not bleed
// into its sibling oneof slots (CSVInput / JSONInput).
func TestEmbedEmptyMessagePresenceParquet(t *testing.T) {
	// Present-but-empty ParquetInput; CSV/JSON unset.
	buf := NewQueryInputSerialization(QueryInputSerializationInput{
		CompressionType: "none",
		ParquetInput:    NewQueryInputSerializationParquetInput(QueryInputSerializationParquetInputInput{}),
	})
	v, err := WrapQueryInputSerialization(buf)
	if err != nil {
		t.Fatalf("WrapQueryInputSerialization: %v", err)
	}
	if v.CompressionType() != "none" {
		t.Fatalf("CompressionType lost: %q", v.CompressionType())
	}
	if _, ok := v.ParquetInput(); !ok {
		t.Fatal("PRESENT empty ParquetInput collapsed to null (embedMessage drop)")
	}
	if _, ok := v.CSVInput(); ok {
		t.Fatal("unset CSVInput read as present (false presence / variant bleed)")
	}
	if _, ok := v.JSONInput(); ok {
		t.Fatal("unset JSONInput read as present (false presence / variant bleed)")
	}

	// Sanity: a fully-unset InputSerialization has no variant present.
	bare, err := WrapQueryInputSerialization(NewQueryInputSerialization(QueryInputSerializationInput{}))
	if err != nil {
		t.Fatalf("WrapQueryInputSerialization(bare): %v", err)
	}
	if _, ok := bare.ParquetInput(); ok {
		t.Fatal("unset ParquetInput read as present")
	}
}

// TestEmbedNonEmptyMessageStillFaithful guards the non-empty branch of the fix
// (payload-only copy + base+(rootOff-HeaderSize)): a JSON input variant with a
// scalar field must round-trip its field intact, proving the HeaderSize shift did
// not corrupt the embedded object's internal pointers.
func TestEmbedNonEmptyMessageStillFaithful(t *testing.T) {
	buf := NewQueryInputSerialization(QueryInputSerializationInput{
		JSONInput: NewQueryInputSerializationJSONInput(QueryInputSerializationJSONInputInput{Type: "LINES"}),
	})
	v, err := WrapQueryInputSerialization(buf)
	if err != nil {
		t.Fatalf("WrapQueryInputSerialization: %v", err)
	}
	j, ok := v.JSONInput()
	if !ok {
		t.Fatal("present JSONInput collapsed to null")
	}
	if j.Type() != "LINES" {
		t.Fatalf("JSONInput.Type corrupted: %q", j.Type())
	}
	if _, ok := v.ParquetInput(); ok {
		t.Fatal("unset ParquetInput read as present")
	}
}
