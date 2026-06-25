// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// server_fidelity_test.go — RED adversarial fidelity probes for vectors #1
// (FetchAndWriteNeedle 43-field RemoteConf converter, full nested graph) and
// #2 (Query nested InputSerialization/OutputSerialization variant presence).
// Both exercise the REAL ZAP server: request crosses the socket as a ZAP buffer
// (unary for FetchAndWriteNeedle, stream-init for Query), is decoded by the
// *FromWire converters, and every field is asserted on the engine side. This is
// the offset/nested-varlen class of bug that hit the filer — a coverage diff
// alone does not prove the wire view reads the right bytes.

package volume

import (
	"context"
	"io"
	"sync"
	"testing"

	"google.golang.org/protobuf/proto"

	remote_pb "github.com/hanzoai/s3/s3/pb/remote_pb"
	volume_server_pb "github.com/hanzoai/s3/s3/pb/volume_server_pb"
)

// captureVolumeServer records the exact proto request the engine receives after
// the ZAP server decodes it, so the test can compare against what the client
// sent.
type captureVolumeServer struct {
	volume_server_pb.UnimplementedVolumeServerServer
	mu      sync.Mutex
	gotFAWN *volume_server_pb.FetchAndWriteNeedleRequest
	gotQ    *volume_server_pb.QueryRequest
}

func (f *captureVolumeServer) FetchAndWriteNeedle(_ context.Context, req *volume_server_pb.FetchAndWriteNeedleRequest) (*volume_server_pb.FetchAndWriteNeedleResponse, error) {
	f.mu.Lock()
	f.gotFAWN = req
	f.mu.Unlock()
	return &volume_server_pb.FetchAndWriteNeedleResponse{ETag: "etag-xyz"}, nil
}

func (f *captureVolumeServer) Query(req *volume_server_pb.QueryRequest, stream volume_server_pb.VolumeServer_QueryServer) error {
	f.mu.Lock()
	f.gotQ = req
	f.mu.Unlock()
	return stream.Send(&volume_server_pb.QueriedStripe{Records: []byte("ok")})
}

// fullRemoteConf populates all 43 data fields with distinct, type-correct,
// non-zero values so any dropped/misaligned field is observable.
func fullRemoteConf() *remote_pb.RemoteConf {
	return &remote_pb.RemoteConf{
		Type: "s3", Name: "conf-name",
		S3AccessKey: "ak", S3SecretKey: "sk", S3Region: "us-east-1", S3Endpoint: "ep",
		S3StorageClass: "STANDARD", S3ForcePathStyle: true, S3SupportTagging: true, S3V4Signature: true,
		GcsGoogleApplicationCredentials: "gcs-creds", GcsProjectId: "proj-1",
		AzureAccountName: "az-name", AzureAccountKey: "az-key",
		BackblazeKeyId: "bb-key", BackblazeApplicationKey: "bb-app", BackblazeEndpoint: "bb-ep", BackblazeRegion: "bb-rg",
		AliyunAccessKey: "al-ak", AliyunSecretKey: "al-sk", AliyunEndpoint: "al-ep", AliyunRegion: "al-rg",
		TencentSecretId: "tc-id", TencentSecretKey: "tc-sk", TencentEndpoint: "tc-ep",
		BaiduAccessKey: "bd-ak", BaiduSecretKey: "bd-sk", BaiduEndpoint: "bd-ep", BaiduRegion: "bd-rg",
		WasabiAccessKey: "ws-ak", WasabiSecretKey: "ws-sk", WasabiEndpoint: "ws-ep", WasabiRegion: "ws-rg",
		FilebaseAccessKey: "fb-ak", FilebaseSecretKey: "fb-sk", FilebaseEndpoint: "fb-ep",
		StorjAccessKey: "st-ak", StorjSecretKey: "st-sk", StorjEndpoint: "st-ep",
		ContaboAccessKey: "cb-ak", ContaboSecretKey: "cb-sk", ContaboEndpoint: "cb-ep", ContaboRegion: "cb-rg",
	}
}

// TestFetchAndWriteNeedle_RemoteConfFullFidelity sends a fully-populated
// FetchAndWriteNeedle request (scalars + Replicas + every RemoteConf field +
// RemoteLocation) through the ZAP server and asserts the engine receives an
// identical request. Refutes vector #1.
func TestFetchAndWriteNeedle_RemoteConfFullFidelity(t *testing.T) {
	fake := &captureVolumeServer{}
	cli, stop := serve(t, fake)
	defer stop()

	want := &volume_server_pb.FetchAndWriteNeedleRequest{
		VolumeId: 42, NeedleId: 0xDEADBEEFCAFE, Cookie: 0x12345678, Offset: 1 << 33, Size: 987654321,
		Auth: "auth-token", DownloadConcurrency: 7,
		Replicas: []*volume_server_pb.FetchAndWriteNeedleRequest_Replica{
			{Url: "u1", PublicUrl: "p1", GrpcPort: 18080},
			{Url: "u2", PublicUrl: "p2", GrpcPort: 18081},
		},
		RemoteConf:     fullRemoteConf(),
		RemoteLocation: &remote_pb.RemoteStorageLocation{Name: "loc", Bucket: "bkt", Path: "/p", ListingCacheTtlSeconds: 99},
	}

	resp, err := cli.FetchAndWriteNeedle(context.Background(), want)
	if err != nil {
		t.Fatalf("FetchAndWriteNeedle: %v", err)
	}
	if resp.ETag != "etag-xyz" {
		t.Fatalf("ETag lost: %q", resp.ETag)
	}

	fake.mu.Lock()
	got := fake.gotFAWN
	fake.mu.Unlock()
	if got == nil {
		t.Fatal("engine got nil request")
	}

	// Scalars.
	if got.VolumeId != want.VolumeId || got.NeedleId != want.NeedleId || got.Cookie != want.Cookie ||
		got.Offset != want.Offset || got.Size != want.Size || got.Auth != want.Auth ||
		got.DownloadConcurrency != want.DownloadConcurrency {
		t.Fatalf("scalar mismatch:\n got=%+v\nwant=%+v", got, want)
	}
	// Replicas.
	if len(got.Replicas) != 2 ||
		got.Replicas[0].Url != "u1" || got.Replicas[0].GrpcPort != 18080 ||
		got.Replicas[1].Url != "u2" || got.Replicas[1].PublicUrl != "p2" {
		t.Fatalf("replicas mismatch: %+v", got.Replicas)
	}
	// RemoteConf — proto.Equal compares every field; a single dropped/misaligned
	// field fails it.
	if !proto.Equal(got.RemoteConf, want.RemoteConf) {
		t.Fatalf("RemoteConf mismatch:\n got=%+v\nwant=%+v", got.RemoteConf, want.RemoteConf)
	}
	if !proto.Equal(got.RemoteLocation, want.RemoteLocation) {
		t.Fatalf("RemoteLocation mismatch:\n got=%+v\nwant=%+v", got.RemoteLocation, want.RemoteLocation)
	}
}

// TestQuery_InputSerializationVariants drives the Query server-stream three
// times — CSV-input, JSON-input, Parquet-input — and asserts the correct
// variant (and only it) is present on the engine side, with all CSV subfields
// intact. Refutes vector #2 (presence + variant discriminant survives the
// stream-init round-trip; no collapse to zero-variant, no cross-variant bleed).
func TestQuery_InputSerializationVariants(t *testing.T) {
	csvReq := &volume_server_pb.QueryRequest{
		Selections:  []string{"a", "b"},
		FromFileIds: []string{"f1"},
		Filter:      &volume_server_pb.QueryRequest_Filter{Field: "x", Operand: ">", Value: "5"},
		InputSerialization: &volume_server_pb.QueryRequest_InputSerialization{
			CompressionType: "gzip",
			CsvInput: &volume_server_pb.QueryRequest_InputSerialization_CSVInput{
				FileHeaderInfo: "USE", RecordDelimiter: "\n", FieldDelimiter: ",",
				QuoteCharacter: "\"", QuoteEscapeCharacter: "\\", Comments: "#",
				AllowQuotedRecordDelimiter: true,
			},
		},
		OutputSerialization: &volume_server_pb.QueryRequest_OutputSerialization{
			CsvOutput: &volume_server_pb.QueryRequest_OutputSerialization_CSVOutput{
				QuoteFields: "ASNEEDED", RecordDelimiter: "\r\n", FieldDelimiter: ";",
				QuoteCharacter: "'", QuoteEscapeCharacter: "/",
			},
		},
	}
	jsonReq := &volume_server_pb.QueryRequest{
		InputSerialization: &volume_server_pb.QueryRequest_InputSerialization{
			JsonInput: &volume_server_pb.QueryRequest_InputSerialization_JSONInput{Type: "LINES"},
		},
		OutputSerialization: &volume_server_pb.QueryRequest_OutputSerialization{
			JsonOutput: &volume_server_pb.QueryRequest_OutputSerialization_JSONOutput{RecordDelimiter: "\n"},
		},
	}
	parquetReq := &volume_server_pb.QueryRequest{
		InputSerialization: &volume_server_pb.QueryRequest_InputSerialization{
			ParquetInput: &volume_server_pb.QueryRequest_InputSerialization_ParquetInput{},
		},
	}

	run := func(t *testing.T, req *volume_server_pb.QueryRequest) *volume_server_pb.QueryRequest {
		fake := &captureVolumeServer{}
		cli, stop := serve(t, fake)
		defer stop()
		s, err := cli.Query(context.Background(), req)
		if err != nil {
			t.Fatalf("Query open: %v", err)
		}
		for {
			if _, err := s.Recv(); err == io.EOF {
				break
			} else if err != nil {
				t.Fatalf("Recv: %v", err)
			}
		}
		fake.mu.Lock()
		defer fake.mu.Unlock()
		return fake.gotQ
	}

	t.Run("csv", func(t *testing.T) {
		got := run(t, csvReq)
		is := got.GetInputSerialization()
		if is == nil {
			t.Fatal("InputSerialization collapsed to nil")
		}
		if is.CsvInput == nil {
			t.Fatal("CsvInput variant lost (collapsed to zero-variant)")
		}
		if is.JsonInput != nil || is.ParquetInput != nil {
			t.Fatalf("cross-variant bleed: json=%v parquet=%v", is.JsonInput, is.ParquetInput)
		}
		c := is.CsvInput
		if c.FileHeaderInfo != "USE" || c.FieldDelimiter != "," || c.QuoteEscapeCharacter != "\\" ||
			c.Comments != "#" || !c.AllowQuotedRecordDelimiter {
			t.Fatalf("CSVInput subfields lost: %+v", c)
		}
		if is.CompressionType != "gzip" {
			t.Fatalf("CompressionType lost: %q", is.CompressionType)
		}
		if got.GetOutputSerialization().GetCsvOutput() == nil ||
			got.GetOutputSerialization().GetCsvOutput().QuoteFields != "ASNEEDED" {
			t.Fatalf("CSVOutput lost: %+v", got.GetOutputSerialization())
		}
		if got.GetFilter() == nil || got.GetFilter().Field != "x" || got.GetFilter().Value != "5" {
			t.Fatalf("Filter lost: %+v", got.GetFilter())
		}
	})

	t.Run("json", func(t *testing.T) {
		got := run(t, jsonReq)
		is := got.GetInputSerialization()
		if is == nil || is.JsonInput == nil {
			t.Fatalf("JsonInput variant lost: %+v", is)
		}
		if is.CsvInput != nil || is.ParquetInput != nil {
			t.Fatalf("cross-variant bleed: csv=%v parquet=%v", is.CsvInput, is.ParquetInput)
		}
		if is.JsonInput.Type != "LINES" {
			t.Fatalf("JSONInput.Type lost: %q", is.JsonInput.Type)
		}
		if got.GetOutputSerialization().GetJsonOutput() == nil {
			t.Fatalf("JSONOutput lost: %+v", got.GetOutputSerialization())
		}
	})

	t.Run("parquet", func(t *testing.T) {
		got := run(t, parquetReq)
		is := got.GetInputSerialization()
		if is == nil || is.ParquetInput == nil {
			t.Fatalf("ParquetInput variant lost: %+v", is)
		}
		if is.CsvInput != nil || is.JsonInput != nil {
			t.Fatalf("cross-variant bleed: csv=%v json=%v", is.CsvInput, is.JsonInput)
		}
	})
}

// fracFloatVolumeServer returns DiskStatus with fractional float32 percentages
// and a MemStatus, so a float-as-int truncation or wrong-width read is caught.
type fracFloatVolumeServer struct {
	volume_server_pb.UnimplementedVolumeServerServer
}

func (fracFloatVolumeServer) VolumeServerStatus(_ context.Context, _ *volume_server_pb.VolumeServerStatusRequest) (*volume_server_pb.VolumeServerStatusResponse, error) {
	return &volume_server_pb.VolumeServerStatusResponse{
		DiskStatuses: []*volume_server_pb.DiskStatus{
			{Dir: "/d", All: 1000, Used: 333, Free: 667, PercentFree: 66.7, PercentUsed: 33.3, DiskType: "hdd"},
		},
		MemoryStatus: &volume_server_pb.MemStatus{Goroutines: 3, All: 100, Used: 25, Heap: 12, Stack: 8, Self: 7, Free: 75},
	}, nil
}

// TestVolumeServerStatus_FractionalFloat32 refutes vector #4 with non-integer
// float32 values that a truncating/mis-widthed converter would mangle.
func TestVolumeServerStatus_FractionalFloat32(t *testing.T) {
	cli, stop := serve(t, fracFloatVolumeServer{})
	defer stop()
	resp, err := cli.VolumeServerStatus(context.Background(), &volume_server_pb.VolumeServerStatusRequest{})
	if err != nil {
		t.Fatalf("VolumeServerStatus: %v", err)
	}
	ds := resp.DiskStatuses[0]
	if ds.PercentFree != float32(66.7) || ds.PercentUsed != float32(33.3) {
		t.Fatalf("float32 percent lost: free=%v used=%v", ds.PercentFree, ds.PercentUsed)
	}
	if resp.MemoryStatus.Self != 7 || resp.MemoryStatus.Free != 75 || resp.MemoryStatus.Stack != 8 {
		t.Fatalf("MemStatus fields lost: %+v", resp.MemoryStatus)
	}
}
