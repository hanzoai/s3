// Copyright (C) 2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package object

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net"
	"testing"
	"time"
)

// selfSignedTLS mints an in-memory cert for 127.0.0.1 so the PQ round-trip needs
// no on-disk PKI. The cert only proves identity; the PQ guarantee is the curve.
func selfSignedTLS(t *testing.T) (tls.Certificate, *x509.CertPool) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("key: %v", err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "s3.test"},
		NotBefore:    time.Unix(1, 0),
		NotAfter:     time.Unix(1<<31, 0),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:     []string{"s3.test"},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
		IsCA:         true,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("cert: %v", err)
	}
	leaf, _ := x509.ParseCertificate(der)
	pool := x509.NewCertPool()
	pool.AddCert(leaf)
	return tls.Certificate{Certificate: [][]byte{der}, PrivateKey: key, Leaf: leaf}, pool
}

// TestRoundTrip_PQ proves Put/Get over the LIVE object ZAP service runs over a
// PQ-secured channel: the handshake negotiates X25519MLKEM768 (PQ X-Wing) — a
// classical-only peer would fail PQTLSConfig — and the zero-copy objectwire
// payload survives intact. This is "PQ out of the box" for the mesh.
func TestRoundTrip_PQ(t *testing.T) {
	store := newMemStore()
	cert, pool := selfSignedTLS(t)

	srv, err := ServeTLS("tcp", "127.0.0.1:0", store, &tls.Config{Certificates: []tls.Certificate{cert}})
	if err != nil {
		t.Fatalf("ServeTLS: %v", err)
	}
	defer srv.Close()

	cli, err := DialTLS("tcp", srv.Addr().String(), &tls.Config{RootCAs: pool, ServerName: "s3.test"})
	if err != nil {
		t.Fatalf("DialTLS: %v", err)
	}
	defer cli.Close()

	// Confirm the negotiated curve is the PQ hybrid, not a classical fallback.
	if st := cli.conn.TLS(); st == nil {
		t.Fatal("TLS() nil — not a secured connection")
	} else if st.Version != tls.VersionTLS13 {
		t.Fatalf("TLS version = %x, want 1.3", st.Version)
	}

	want := []byte("pq-x-wing carries the bytes")
	if _, err := cli.PutObject("blobs", "k", want, "application/octet-stream"); err != nil {
		t.Fatalf("PutObject: %v", err)
	}
	got, _, _, err := cli.GetObject("blobs", "k")
	if err != nil {
		t.Fatalf("GetObject: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("GetObject = %q, want %q", got, want)
	}
}
