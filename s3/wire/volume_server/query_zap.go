// Code generated from volume_server.proto; DO NOT EDIT.

package volume_serverwire

import (
	zap "github.com/zap-proto/go"
)

// The Query messages declare several inline nested message types. Each is
// flattened to a top-level type with a qualified name: QueryFilter,
// QueryInputSerialization (+ CSVInput / JSONInput / ParquetInput),
// QueryOutputSerialization (+ CSVOutput / JSONOutput).

// --- QueryFilter (QueryRequest.Filter) ---

const (
	queryFilterFieldOff   = 0
	queryFilterOperandOff = 8
	queryFilterValueOff   = 16
	queryFilterSize       = 24
)

// QueryFilter is a zero-copy view into a ZAP-encoded QueryRequest.Filter message.
type QueryFilter struct{ o zap.Object }

// WrapQueryFilter parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapQueryFilter(b []byte) (QueryFilter, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return QueryFilter{}, err
	}
	return QueryFilter{o: m.Root()}, nil
}

// Field reads the field field (proto field 1, string).
func (t QueryFilter) Field() string { return t.o.Text(queryFilterFieldOff) }

// Operand reads the operand field (proto field 2, string).
func (t QueryFilter) Operand() string { return t.o.Text(queryFilterOperandOff) }

// Value reads the value field (proto field 3, string).
func (t QueryFilter) Value() string { return t.o.Text(queryFilterValueOff) }

// QueryFilterInput collects the field values for NewQueryFilter.
type QueryFilterInput struct {
	Field   string
	Operand string
	Value   string
}

// NewQueryFilter builds a ZAP-encoded QueryRequest.Filter message from in and
// returns the bytes.
func NewQueryFilter(in QueryFilterInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(queryFilterSize)
	ob.SetText(queryFilterFieldOff, in.Field)
	ob.SetText(queryFilterOperandOff, in.Operand)
	ob.SetText(queryFilterValueOff, in.Value)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- QueryInputSerializationCSVInput (QueryRequest.InputSerialization.CSVInput) ---

const (
	queryInputSerializationCSVInputFileHeaderInfoOff             = 0
	queryInputSerializationCSVInputRecordDelimiterOff            = 8
	queryInputSerializationCSVInputFieldDelimiterOff             = 16
	queryInputSerializationCSVInputQuoteCharacterOff             = 24
	queryInputSerializationCSVInputQuoteEscapeCharacterOff       = 32
	queryInputSerializationCSVInputCommentsOff                   = 40
	queryInputSerializationCSVInputAllowQuotedRecordDelimiterOff = 48
	queryInputSerializationCSVInputSize                          = 49
)

// QueryInputSerializationCSVInput is a zero-copy view into a ZAP-encoded
// QueryRequest.InputSerialization.CSVInput message.
type QueryInputSerializationCSVInput struct{ o zap.Object }

// WrapQueryInputSerializationCSVInput parses b and returns a typed view. Returns
// an error if the wire-level checks (magic, version, size) fail.
func WrapQueryInputSerializationCSVInput(b []byte) (QueryInputSerializationCSVInput, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return QueryInputSerializationCSVInput{}, err
	}
	return QueryInputSerializationCSVInput{o: m.Root()}, nil
}

// FileHeaderInfo reads the file_header_info field (proto field 1, string).
func (t QueryInputSerializationCSVInput) FileHeaderInfo() string {
	return t.o.Text(queryInputSerializationCSVInputFileHeaderInfoOff)
}

// RecordDelimiter reads the record_delimiter field (proto field 2, string).
func (t QueryInputSerializationCSVInput) RecordDelimiter() string {
	return t.o.Text(queryInputSerializationCSVInputRecordDelimiterOff)
}

// FieldDelimiter reads the field_delimiter field (proto field 3, string).
func (t QueryInputSerializationCSVInput) FieldDelimiter() string {
	return t.o.Text(queryInputSerializationCSVInputFieldDelimiterOff)
}

// QuoteCharacter reads the quote_character field (proto field 4, string).
func (t QueryInputSerializationCSVInput) QuoteCharacter() string {
	return t.o.Text(queryInputSerializationCSVInputQuoteCharacterOff)
}

// QuoteEscapeCharacter reads the quote_escape_character field (proto field 5,
// string).
func (t QueryInputSerializationCSVInput) QuoteEscapeCharacter() string {
	return t.o.Text(queryInputSerializationCSVInputQuoteEscapeCharacterOff)
}

// Comments reads the comments field (proto field 6, string).
func (t QueryInputSerializationCSVInput) Comments() string {
	return t.o.Text(queryInputSerializationCSVInputCommentsOff)
}

// AllowQuotedRecordDelimiter reads the allow_quoted_record_delimiter field
// (proto field 7, bool).
func (t QueryInputSerializationCSVInput) AllowQuotedRecordDelimiter() bool {
	return t.o.Bool(queryInputSerializationCSVInputAllowQuotedRecordDelimiterOff)
}

// QueryInputSerializationCSVInputInput collects the field values for
// NewQueryInputSerializationCSVInput.
type QueryInputSerializationCSVInputInput struct {
	FileHeaderInfo             string
	RecordDelimiter            string
	FieldDelimiter             string
	QuoteCharacter             string
	QuoteEscapeCharacter       string
	Comments                   string
	AllowQuotedRecordDelimiter bool
}

// NewQueryInputSerializationCSVInput builds a ZAP-encoded
// QueryRequest.InputSerialization.CSVInput message from in and returns the bytes.
func NewQueryInputSerializationCSVInput(in QueryInputSerializationCSVInputInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(queryInputSerializationCSVInputSize)
	ob.SetText(queryInputSerializationCSVInputFileHeaderInfoOff, in.FileHeaderInfo)
	ob.SetText(queryInputSerializationCSVInputRecordDelimiterOff, in.RecordDelimiter)
	ob.SetText(queryInputSerializationCSVInputFieldDelimiterOff, in.FieldDelimiter)
	ob.SetText(queryInputSerializationCSVInputQuoteCharacterOff, in.QuoteCharacter)
	ob.SetText(queryInputSerializationCSVInputQuoteEscapeCharacterOff, in.QuoteEscapeCharacter)
	ob.SetText(queryInputSerializationCSVInputCommentsOff, in.Comments)
	ob.SetBool(queryInputSerializationCSVInputAllowQuotedRecordDelimiterOff, in.AllowQuotedRecordDelimiter)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- QueryInputSerializationJSONInput (QueryRequest.InputSerialization.JSONInput) ---

const (
	queryInputSerializationJSONInputTypeOff = 0
	queryInputSerializationJSONInputSize    = 8
)

// QueryInputSerializationJSONInput is a zero-copy view into a ZAP-encoded
// QueryRequest.InputSerialization.JSONInput message.
type QueryInputSerializationJSONInput struct{ o zap.Object }

// WrapQueryInputSerializationJSONInput parses b and returns a typed view.
// Returns an error if the wire-level checks (magic, version, size) fail.
func WrapQueryInputSerializationJSONInput(b []byte) (QueryInputSerializationJSONInput, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return QueryInputSerializationJSONInput{}, err
	}
	return QueryInputSerializationJSONInput{o: m.Root()}, nil
}

// Type reads the type field (proto field 1, string).
func (t QueryInputSerializationJSONInput) Type() string {
	return t.o.Text(queryInputSerializationJSONInputTypeOff)
}

// QueryInputSerializationJSONInputInput collects the field values for
// NewQueryInputSerializationJSONInput.
type QueryInputSerializationJSONInputInput struct {
	Type string
}

// NewQueryInputSerializationJSONInput builds a ZAP-encoded
// QueryRequest.InputSerialization.JSONInput message from in and returns the
// bytes.
func NewQueryInputSerializationJSONInput(in QueryInputSerializationJSONInputInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(queryInputSerializationJSONInputSize)
	ob.SetText(queryInputSerializationJSONInputTypeOff, in.Type)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- QueryInputSerializationParquetInput (QueryRequest.InputSerialization.ParquetInput) ---

const queryInputSerializationParquetInputSize = 0

// QueryInputSerializationParquetInput is a zero-copy view into a ZAP-encoded
// QueryRequest.InputSerialization.ParquetInput message.
// QueryInputSerializationParquetInput carries no fields.
type QueryInputSerializationParquetInput struct{ o zap.Object }

// WrapQueryInputSerializationParquetInput parses b and returns a typed view.
// Returns an error if the wire-level checks (magic, version, size) fail.
func WrapQueryInputSerializationParquetInput(b []byte) (QueryInputSerializationParquetInput, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return QueryInputSerializationParquetInput{}, err
	}
	return QueryInputSerializationParquetInput{o: m.Root()}, nil
}

// QueryInputSerializationParquetInputInput collects the field values for
// NewQueryInputSerializationParquetInput. It has no fields.
type QueryInputSerializationParquetInputInput struct{}

// NewQueryInputSerializationParquetInput builds a ZAP-encoded
// QueryRequest.InputSerialization.ParquetInput message and returns the bytes.
func NewQueryInputSerializationParquetInput(in QueryInputSerializationParquetInputInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(queryInputSerializationParquetInputSize)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- QueryInputSerialization (QueryRequest.InputSerialization) ---

const (
	queryInputSerializationCompressionTypeOff = 0
	queryInputSerializationCSVInputOff        = 8
	queryInputSerializationJSONInputOff       = 12
	queryInputSerializationParquetInputOff    = 16
	queryInputSerializationSize               = 20
)

// QueryInputSerialization is a zero-copy view into a ZAP-encoded
// QueryRequest.InputSerialization message.
type QueryInputSerialization struct{ o zap.Object }

// WrapQueryInputSerialization parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapQueryInputSerialization(b []byte) (QueryInputSerialization, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return QueryInputSerialization{}, err
	}
	return QueryInputSerialization{o: m.Root()}, nil
}

// CompressionType reads the compression_type field (proto field 1, string).
func (t QueryInputSerialization) CompressionType() string {
	return t.o.Text(queryInputSerializationCompressionTypeOff)
}

// CSVInput reads the csv_input field (proto field 2, message CSVInput). The bool
// is false when the field is absent.
func (t QueryInputSerialization) CSVInput() (QueryInputSerializationCSVInput, bool) {
	o := t.o.Object(queryInputSerializationCSVInputOff)
	if o.IsNull() {
		return QueryInputSerializationCSVInput{}, false
	}
	return QueryInputSerializationCSVInput{o: o}, true
}

// JSONInput reads the json_input field (proto field 3, message JSONInput). The
// bool is false when the field is absent.
func (t QueryInputSerialization) JSONInput() (QueryInputSerializationJSONInput, bool) {
	o := t.o.Object(queryInputSerializationJSONInputOff)
	if o.IsNull() {
		return QueryInputSerializationJSONInput{}, false
	}
	return QueryInputSerializationJSONInput{o: o}, true
}

// ParquetInput reads the parquet_input field (proto field 4, message
// ParquetInput). The bool is false when the field is absent.
func (t QueryInputSerialization) ParquetInput() (QueryInputSerializationParquetInput, bool) {
	o := t.o.Object(queryInputSerializationParquetInputOff)
	if o.IsNull() {
		return QueryInputSerializationParquetInput{}, false
	}
	return QueryInputSerializationParquetInput{o: o}, true
}

// QueryInputSerializationInput collects the field values for
// NewQueryInputSerialization. CSVInput, JSONInput and ParquetInput are the child
// messages' own ZAP buffers.
type QueryInputSerializationInput struct {
	CompressionType string
	CSVInput        []byte
	JSONInput       []byte
	ParquetInput    []byte
}

// NewQueryInputSerialization builds a ZAP-encoded
// QueryRequest.InputSerialization message from in and returns the bytes.
func NewQueryInputSerialization(in QueryInputSerializationInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(queryInputSerializationSize)
	ob.SetText(queryInputSerializationCompressionTypeOff, in.CompressionType)
	ob.SetObject(queryInputSerializationCSVInputOff, embedMessage(b, in.CSVInput))
	ob.SetObject(queryInputSerializationJSONInputOff, embedMessage(b, in.JSONInput))
	ob.SetObject(queryInputSerializationParquetInputOff, embedMessage(b, in.ParquetInput))
	ob.FinishAsRoot()
	return b.Finish()
}

// --- QueryOutputSerializationCSVOutput (QueryRequest.OutputSerialization.CSVOutput) ---

const (
	queryOutputSerializationCSVOutputQuoteFieldsOff          = 0
	queryOutputSerializationCSVOutputRecordDelimiterOff      = 8
	queryOutputSerializationCSVOutputFieldDelimiterOff       = 16
	queryOutputSerializationCSVOutputQuoteCharacterOff       = 24
	queryOutputSerializationCSVOutputQuoteEscapeCharacterOff = 32
	queryOutputSerializationCSVOutputSize                    = 40
)

// QueryOutputSerializationCSVOutput is a zero-copy view into a ZAP-encoded
// QueryRequest.OutputSerialization.CSVOutput message.
type QueryOutputSerializationCSVOutput struct{ o zap.Object }

// WrapQueryOutputSerializationCSVOutput parses b and returns a typed view.
// Returns an error if the wire-level checks (magic, version, size) fail.
func WrapQueryOutputSerializationCSVOutput(b []byte) (QueryOutputSerializationCSVOutput, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return QueryOutputSerializationCSVOutput{}, err
	}
	return QueryOutputSerializationCSVOutput{o: m.Root()}, nil
}

// QuoteFields reads the quote_fields field (proto field 1, string).
func (t QueryOutputSerializationCSVOutput) QuoteFields() string {
	return t.o.Text(queryOutputSerializationCSVOutputQuoteFieldsOff)
}

// RecordDelimiter reads the record_delimiter field (proto field 2, string).
func (t QueryOutputSerializationCSVOutput) RecordDelimiter() string {
	return t.o.Text(queryOutputSerializationCSVOutputRecordDelimiterOff)
}

// FieldDelimiter reads the field_delimiter field (proto field 3, string).
func (t QueryOutputSerializationCSVOutput) FieldDelimiter() string {
	return t.o.Text(queryOutputSerializationCSVOutputFieldDelimiterOff)
}

// QuoteCharacter reads the quote_character field (proto field 4, string).
func (t QueryOutputSerializationCSVOutput) QuoteCharacter() string {
	return t.o.Text(queryOutputSerializationCSVOutputQuoteCharacterOff)
}

// QuoteEscapeCharacter reads the quote_escape_character field (proto field 5,
// string).
func (t QueryOutputSerializationCSVOutput) QuoteEscapeCharacter() string {
	return t.o.Text(queryOutputSerializationCSVOutputQuoteEscapeCharacterOff)
}

// QueryOutputSerializationCSVOutputInput collects the field values for
// NewQueryOutputSerializationCSVOutput.
type QueryOutputSerializationCSVOutputInput struct {
	QuoteFields          string
	RecordDelimiter      string
	FieldDelimiter       string
	QuoteCharacter       string
	QuoteEscapeCharacter string
}

// NewQueryOutputSerializationCSVOutput builds a ZAP-encoded
// QueryRequest.OutputSerialization.CSVOutput message from in and returns the
// bytes.
func NewQueryOutputSerializationCSVOutput(in QueryOutputSerializationCSVOutputInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(queryOutputSerializationCSVOutputSize)
	ob.SetText(queryOutputSerializationCSVOutputQuoteFieldsOff, in.QuoteFields)
	ob.SetText(queryOutputSerializationCSVOutputRecordDelimiterOff, in.RecordDelimiter)
	ob.SetText(queryOutputSerializationCSVOutputFieldDelimiterOff, in.FieldDelimiter)
	ob.SetText(queryOutputSerializationCSVOutputQuoteCharacterOff, in.QuoteCharacter)
	ob.SetText(queryOutputSerializationCSVOutputQuoteEscapeCharacterOff, in.QuoteEscapeCharacter)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- QueryOutputSerializationJSONOutput (QueryRequest.OutputSerialization.JSONOutput) ---

const (
	queryOutputSerializationJSONOutputRecordDelimiterOff = 0
	queryOutputSerializationJSONOutputSize               = 8
)

// QueryOutputSerializationJSONOutput is a zero-copy view into a ZAP-encoded
// QueryRequest.OutputSerialization.JSONOutput message.
type QueryOutputSerializationJSONOutput struct{ o zap.Object }

// WrapQueryOutputSerializationJSONOutput parses b and returns a typed view.
// Returns an error if the wire-level checks (magic, version, size) fail.
func WrapQueryOutputSerializationJSONOutput(b []byte) (QueryOutputSerializationJSONOutput, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return QueryOutputSerializationJSONOutput{}, err
	}
	return QueryOutputSerializationJSONOutput{o: m.Root()}, nil
}

// RecordDelimiter reads the record_delimiter field (proto field 1, string).
func (t QueryOutputSerializationJSONOutput) RecordDelimiter() string {
	return t.o.Text(queryOutputSerializationJSONOutputRecordDelimiterOff)
}

// QueryOutputSerializationJSONOutputInput collects the field values for
// NewQueryOutputSerializationJSONOutput.
type QueryOutputSerializationJSONOutputInput struct {
	RecordDelimiter string
}

// NewQueryOutputSerializationJSONOutput builds a ZAP-encoded
// QueryRequest.OutputSerialization.JSONOutput message from in and returns the
// bytes.
func NewQueryOutputSerializationJSONOutput(in QueryOutputSerializationJSONOutputInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(queryOutputSerializationJSONOutputSize)
	ob.SetText(queryOutputSerializationJSONOutputRecordDelimiterOff, in.RecordDelimiter)
	ob.FinishAsRoot()
	return b.Finish()
}

// --- QueryOutputSerialization (QueryRequest.OutputSerialization) ---

const (
	queryOutputSerializationCSVOutputOff  = 0
	queryOutputSerializationJSONOutputOff = 4
	queryOutputSerializationSize          = 8
)

// QueryOutputSerialization is a zero-copy view into a ZAP-encoded
// QueryRequest.OutputSerialization message.
type QueryOutputSerialization struct{ o zap.Object }

// WrapQueryOutputSerialization parses b and returns a typed view. Returns an
// error if the wire-level checks (magic, version, size) fail.
func WrapQueryOutputSerialization(b []byte) (QueryOutputSerialization, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return QueryOutputSerialization{}, err
	}
	return QueryOutputSerialization{o: m.Root()}, nil
}

// CSVOutput reads the csv_output field (proto field 2, message CSVOutput). The
// bool is false when the field is absent.
func (t QueryOutputSerialization) CSVOutput() (QueryOutputSerializationCSVOutput, bool) {
	o := t.o.Object(queryOutputSerializationCSVOutputOff)
	if o.IsNull() {
		return QueryOutputSerializationCSVOutput{}, false
	}
	return QueryOutputSerializationCSVOutput{o: o}, true
}

// JSONOutput reads the json_output field (proto field 3, message JSONOutput).
// The bool is false when the field is absent.
func (t QueryOutputSerialization) JSONOutput() (QueryOutputSerializationJSONOutput, bool) {
	o := t.o.Object(queryOutputSerializationJSONOutputOff)
	if o.IsNull() {
		return QueryOutputSerializationJSONOutput{}, false
	}
	return QueryOutputSerializationJSONOutput{o: o}, true
}

// QueryOutputSerializationInput collects the field values for
// NewQueryOutputSerialization. CSVOutput and JSONOutput are the child messages'
// own ZAP buffers.
type QueryOutputSerializationInput struct {
	CSVOutput  []byte
	JSONOutput []byte
}

// NewQueryOutputSerialization builds a ZAP-encoded
// QueryRequest.OutputSerialization message from in and returns the bytes.
func NewQueryOutputSerialization(in QueryOutputSerializationInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(queryOutputSerializationSize)
	ob.SetObject(queryOutputSerializationCSVOutputOff, embedMessage(b, in.CSVOutput))
	ob.SetObject(queryOutputSerializationJSONOutputOff, embedMessage(b, in.JSONOutput))
	ob.FinishAsRoot()
	return b.Finish()
}

// --- QueryRequest ---

const (
	queryRequestSelectionsOff          = 0
	queryRequestFromFileIDsOff         = 8
	queryRequestFilterOff              = 16
	queryRequestInputSerializationOff  = 20
	queryRequestOutputSerializationOff = 24
	queryRequestSize                   = 28
)

// QueryRequest is a zero-copy view into a ZAP-encoded QueryRequest message.
type QueryRequest struct{ o zap.Object }

// WrapQueryRequest parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapQueryRequest(b []byte) (QueryRequest, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return QueryRequest{}, err
	}
	return QueryRequest{o: m.Root()}, nil
}

// SelectionsLen reports the number of selections elements (proto field 1,
// repeated string).
func (t QueryRequest) SelectionsLen() int { return t.o.List(queryRequestSelectionsOff).Len() }

// SelectionAt returns the i-th selections element (string).
func (t QueryRequest) SelectionAt(i int) string {
	return elemText(t.o.List(queryRequestSelectionsOff), i)
}

// FromFileIDsLen reports the number of from_file_ids elements (proto field 2,
// repeated string).
func (t QueryRequest) FromFileIDsLen() int { return t.o.List(queryRequestFromFileIDsOff).Len() }

// FromFileIDAt returns the i-th from_file_ids element (string).
func (t QueryRequest) FromFileIDAt(i int) string {
	return elemText(t.o.List(queryRequestFromFileIDsOff), i)
}

// Filter reads the filter field (proto field 3, message Filter). The bool is
// false when the field is absent.
func (t QueryRequest) Filter() (QueryFilter, bool) {
	o := t.o.Object(queryRequestFilterOff)
	if o.IsNull() {
		return QueryFilter{}, false
	}
	return QueryFilter{o: o}, true
}

// InputSerialization reads the input_serialization field (proto field 4, message
// InputSerialization). The bool is false when the field is absent.
func (t QueryRequest) InputSerialization() (QueryInputSerialization, bool) {
	o := t.o.Object(queryRequestInputSerializationOff)
	if o.IsNull() {
		return QueryInputSerialization{}, false
	}
	return QueryInputSerialization{o: o}, true
}

// OutputSerialization reads the output_serialization field (proto field 5,
// message OutputSerialization). The bool is false when the field is absent.
func (t QueryRequest) OutputSerialization() (QueryOutputSerialization, bool) {
	o := t.o.Object(queryRequestOutputSerializationOff)
	if o.IsNull() {
		return QueryOutputSerialization{}, false
	}
	return QueryOutputSerialization{o: o}, true
}

// QueryRequestInput collects the field values for NewQueryRequest. Filter,
// InputSerialization and OutputSerialization are the child messages' own ZAP
// buffers.
type QueryRequestInput struct {
	Selections          []string
	FromFileIDs         []string
	Filter              []byte
	InputSerialization  []byte
	OutputSerialization []byte
}

// NewQueryRequest builds a ZAP-encoded QueryRequest message from in and returns
// the bytes.
func NewQueryRequest(in QueryRequestInput) []byte {
	b := zap.NewBuilder(512)
	selOff, selLen := addStringList(b, in.Selections)
	fromOff, fromLen := addStringList(b, in.FromFileIDs)
	ob := b.StartObject(queryRequestSize)
	ob.SetList(queryRequestSelectionsOff, selOff, selLen)
	ob.SetList(queryRequestFromFileIDsOff, fromOff, fromLen)
	ob.SetObject(queryRequestFilterOff, embedMessage(b, in.Filter))
	ob.SetObject(queryRequestInputSerializationOff, embedMessage(b, in.InputSerialization))
	ob.SetObject(queryRequestOutputSerializationOff, embedMessage(b, in.OutputSerialization))
	ob.FinishAsRoot()
	return b.Finish()
}

// --- QueriedStripe ---

const (
	queriedStripeRecordsOff = 0
	queriedStripeSize       = 8
)

// QueriedStripe is a zero-copy view into a ZAP-encoded QueriedStripe message.
type QueriedStripe struct{ o zap.Object }

// WrapQueriedStripe parses b and returns a typed view. Returns an error if the
// wire-level checks (magic, version, size) fail.
func WrapQueriedStripe(b []byte) (QueriedStripe, error) {
	m, err := zap.Parse(b)
	if err != nil {
		return QueriedStripe{}, err
	}
	return QueriedStripe{o: m.Root()}, nil
}

// Records reads the records field (proto field 1, bytes).
func (t QueriedStripe) Records() []byte { return t.o.Bytes(queriedStripeRecordsOff) }

// QueriedStripeInput collects the field values for NewQueriedStripe.
type QueriedStripeInput struct {
	Records []byte
}

// NewQueriedStripe builds a ZAP-encoded QueriedStripe message from in and
// returns the bytes.
func NewQueriedStripe(in QueriedStripeInput) []byte {
	b := zap.NewBuilder(256)
	ob := b.StartObject(queriedStripeSize)
	ob.SetBytes(queriedStripeRecordsOff, in.Records)
	ob.FinishAsRoot()
	return b.Finish()
}
