// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package master

import (
	"context"
	"strings"
	"testing"

	"fmt"

	master_pb "github.com/hanzoai/s3/s3/pb/master_pb"
	masterwire "github.com/hanzoai/s3/s3/wire/master"
	masterstream "github.com/hanzoai/s3/s3/wire/master/masterstream"
	"github.com/zap-proto/go/transport"
)

// errBackend returns a warmup-style error from Statistics tagged with its code
// name, exactly as the gRPC-free master server now does
// (fmt.Errorf("Unavailable: ...")). Every other method is unimplemented via the
// embedded interface; only Statistics is exercised.
type errBackend struct{ masterwire.Backend }

func (errBackend) Statistics(req []byte) ([]byte, error) {
	return nil, fmt.Errorf("Unavailable: master is warming up, topology is still loading")
}

// TestCodeNameSurvivesZapRoundTrip proves the contract every converted error
// classifier relies on: a server error carrying a PascalCase code name reaches
// the client with that name still present in err.Error(), so string matching
// (strings.Contains(err, "Unavailable")) classifies it correctly with no gRPC
// status types in the loop.
func TestCodeNameSurvivesZapRoundTrip(t *testing.T) {
	srv, err := transport.ListenStream("tcp", "127.0.0.1:0",
		masterwire.Dispatch(errBackend{}), masterstream.Handler(streamServer{}))
	if err != nil {
		t.Fatalf("ListenStream: %v", err)
	}
	defer srv.Close()
	conn, err := transport.Dial("tcp", srv.Addr().String())
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer conn.Close()

	cli := New(conn, nil)
	_, rpcErr := cli.Statistics(context.Background(), &master_pb.StatisticsRequest{})
	if rpcErr == nil {
		t.Fatal("expected an error from Statistics")
	}
	if !strings.Contains(rpcErr.Error(), "Unavailable") {
		t.Fatalf("code name lost across ZAP round-trip: %q does not contain %q", rpcErr.Error(), "Unavailable")
	}
}
