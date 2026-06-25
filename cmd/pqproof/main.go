package main

import (
	"crypto/tls"
	"crypto/x509"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/zap-proto/go/rpc"
	"github.com/zap-proto/go/transport"

	"github.com/hanzoai/s3/s3/pb"
	filer_pb "github.com/hanzoai/s3/s3/pb/filer_pb"
)

const addr = "127.0.0.1:18888"

func clientTLS() *tls.Config {
	cert, err := tls.LoadX509KeyPair("/tmp/pqcerts/client.crt", "/tmp/pqcerts/client.key")
	if err != nil { panic(err) }
	caPEM, err := os.ReadFile("/tmp/pqcerts/ca.crt")
	if err != nil { panic(err) }
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(caPEM) { panic("ca parse") }
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
		ServerName:   "localhost",
	}
}

func main() {
	fail := false

	// ---------- POSITIVE: PQ dial (production path: DialTLS+PQTLSConfig) ----------
	fmt.Println("== [1] PQ-TLS dial (transport.DialTLS + PQTLSConfig, X25519MLKEM768) ==")
	conn, err := transport.DialTLS("tcp", addr, transport.PQTLSConfig(clientTLS()))
	if err != nil {
		fmt.Printf("  FAIL: PQ DialTLS errored: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	st := conn.TLS()
	if st == nil {
		fmt.Println("  FAIL: conn.TLS() == nil (connection is NOT TLS)")
		fail = true
	} else {
		fmt.Printf("  conn.TLS().Version  = 0x%04x (TLS1.3=0x%04x)\n", st.Version, tls.VersionTLS13)
		fmt.Printf("  conn.TLS().CurveID  = %v (want X25519MLKEM768=%v)\n", st.CurveID, tls.X25519MLKEM768)
		fmt.Printf("  conn.TLS().CipherSuite = %s\n", tls.CipherSuiteName(st.CipherSuite))
		if st.Version != tls.VersionTLS13 { fmt.Println("  FAIL: not TLS1.3"); fail = true }
		if st.CurveID != tls.X25519MLKEM768 {
			fmt.Printf("  FAIL: negotiated curve %v is NOT the PQ hybrid\n", st.CurveID); fail = true
		} else {
			fmt.Println("  OK: live connection negotiated X25519MLKEM768 (PQ X-Wing)")
		}
	}

	// ---------- real round-trip over the PQ connection ----------
	fmt.Println("== [1b] CreateEntry + LookupDirectoryEntry round-trip over the PQ conn ==")
	cl := pb.NewZapFilerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dir := "/pqproof"
	name := fmt.Sprintf("hello-%d.txt", time.Now().UnixNano())
	_, err = cl.CreateEntry(ctx, &filer_pb.CreateEntryRequest{
		Directory: dir,
		Entry: &filer_pb.Entry{
			Name:        name,
			IsDirectory: false,
			Attributes:  &filer_pb.FuseAttributes{FileMode: 0644, FileSize: 5},
		},
	})
	if err != nil {
		fmt.Printf("  FAIL: CreateEntry over PQ-TLS: %v\n", err); fail = true
	} else {
		fmt.Printf("  CreateEntry(%s/%s) OK\n", dir, name)
		lr, err := cl.LookupDirectoryEntry(ctx, &filer_pb.LookupDirectoryEntryRequest{Directory: dir, Name: name})
		if err != nil {
			fmt.Printf("  FAIL: LookupDirectoryEntry over PQ-TLS: %v\n", err); fail = true
		} else if lr.GetEntry().GetName() != name {
			fmt.Printf("  FAIL: lookup returned %q want %q\n", lr.GetEntry().GetName(), name); fail = true
		} else {
			fmt.Printf("  LookupDirectoryEntry OK -> entry=%q (round-trip over PQ-TLS confirmed)\n", lr.GetEntry().GetName())
		}
	}

	// ---------- NEGATIVE 1: plaintext Dial to the PQ-TLS filer ----------
	fmt.Println("== [2] NEGATIVE: plaintext transport.Dial (no TLS) to PQ-TLS filer ==")
	pc, perr := transport.Dial("tcp", addr)
	if perr != nil {
		fmt.Printf("  OK: plaintext Dial refused at connect: %v\n", perr)
	} else {
		// TCP may connect; a real RPC must fail because server speaks TLS records.
		// Attempt an actual call; expect failure / no valid RPC response.
		fmt.Println("  (TCP open; issuing a plaintext RPC — must NOT succeed)")
		ping := rpc.BuildRequest(rpc.Call{Method: 1, PromiseID: pc.NextPromiseID(), Payload: []byte("x")})
		type res struct{ err error }
		done := make(chan res, 1)
		go func() { _, e := pc.Call(ping); done <- res{e} }()
		select {
		case r := <-done:
			if r.err != nil {
				fmt.Printf("  OK: plaintext RPC rejected by PQ-TLS server: %v\n", r.err)
			} else {
				fmt.Println("  FAIL: plaintext RPC SUCCEEDED against PQ-TLS filer (downgrade!)"); fail = true
			}
		case <-time.After(5 * time.Second):
			fmt.Println("  OK: plaintext RPC got NO valid response (server speaks TLS, not plaintext ZAP — timed out)")
		}
		pc.Close()
	}

	// ---------- NEGATIVE 2: classical-only X25519 DialTLS ----------
	fmt.Println("== [3] NEGATIVE: classical-only (CurvePreferences=[X25519]) DialTLS ==")
	classical := clientTLS()
	classical.MinVersion = tls.VersionTLS13
	classical.CurvePreferences = []tls.CurveID{tls.X25519} // no PQ offered
	cc, xerr := transport.DialTLS("tcp", addr, classical)
	if xerr != nil {
		fmt.Printf("  OK: classical-only client REFUSED (no downgrade): %v\n", xerr)
	} else {
		fmt.Println("  FAIL: classical-only client ACCEPTED — PQ not enforced (downgrade!)"); fail = true
		cc.Close()
	}

	fmt.Println("----")
	if fail {
		fmt.Println("RESULT: FAIL")
		os.Exit(1)
	}
	fmt.Println("RESULT: ALL CHECKS PASSED")
}
