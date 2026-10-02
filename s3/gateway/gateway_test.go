package gateway

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"

	"github.com/hanzoai/s3/s3/credential"
	_ "github.com/hanzoai/s3/s3/credential/memory"
	"github.com/hanzoai/s3/s3/s3api"
	"github.com/hanzoai/s3/s3/s3api/s3err"
)

var (
	upKey     = aws.Credentials{AccessKeyID: "UPSTREAMKEY000000000", SecretAccessKey: "upstreamsecret00000000000000000000000000"}
	adminKey  = aws.Credentials{AccessKeyID: "ADMINKEY000000000000", SecretAccessKey: "adminsecret000000000000000000000000000000"}
	backupKey = aws.Credentials{AccessKeyID: "BACKUPKEY00000000000", SecretAccessKey: "backupsecret00000000000000000000000000000"}
	names     = Names{Prefix: "hz-", Suffix: "-acct"}
)

// iamFor returns an IAM holding the identities in config.
func iamFor(t *testing.T, config string) *s3api.IdentityAccessManagement {
	t.Helper()
	// The IAM would add an admin identity for these.
	t.Setenv("AWS_ACCESS_KEY_ID", "")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "")
	path := filepath.Join(t.TempDir(), "s3.json")
	if err := os.WriteFile(path, []byte(config), 0o600); err != nil {
		t.Fatal(err)
	}
	return s3api.NewIdentityAccessManagementWithStore(&s3api.S3ApiServerOption{Config: path}, nil, string(credential.StoreTypeMemory))
}

// store is a fake upstream: it checks every request's signature with the
// upstream key, as S3 would, and keeps objects in memory by host and path.
type store struct {
	iam     *s3api.IdentityAccessManagement
	mu      sync.Mutex
	objects map[string][]byte
	types   map[string]string
	seen    []*http.Request
}

func (s *store) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	s.seen = append(s.seen, r)
	s.mu.Unlock()
	if _, code := s.iam.AuthenticateRequest(r); code != s3err.ErrNone {
		http.Error(w, "upstream refused the gateway's signature", http.StatusForbidden)
		return
	}
	bucket := strings.TrimSuffix(r.Host, ".s3.test")
	key := bucket + r.URL.Path
	switch {
	case r.Host == "s3.test" && r.Method == http.MethodGet:
		w.Header().Set("Content-Type", "application/xml")
		fmt.Fprint(w, `<ListAllMyBucketsResult><Buckets><Bucket><Name>hz-alpha-acct</Name><CreationDate>2026-01-01T00:00:00Z</CreationDate></Bucket><Bucket><Name>someone-else</Name><CreationDate>2026-01-01T00:00:00Z</CreationDate></Bucket><Bucket><Name>hz-backups-acct</Name><CreationDate>2026-01-01T00:00:00Z</CreationDate></Bucket></Buckets></ListAllMyBucketsResult>`)
	case r.Method == http.MethodPut:
		var body io.Reader = r.Body
		if r.Header.Get("X-Amz-Content-Sha256") == trailerPayload {
			decoded, code := s.iam.DecodedBody(r)
			if code != s3err.ErrNone {
				http.Error(w, "bad aws-chunked body", http.StatusBadRequest)
				return
			}
			body = decoded
		}
		b, err := io.ReadAll(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		s.mu.Lock()
		s.objects[key] = b
		s.types[key] = r.Header.Get("Content-Type")
		s.mu.Unlock()
		w.Header().Set("ETag", `"x"`)
		w.Header().Set("X-Amz-Checksum-Crc32", "AAAAAA==")
	case r.Method == http.MethodGet && r.URL.Query().Has("list-type"):
		w.Header().Set("Content-Type", "application/xml")
		fmt.Fprintf(w, `<ListBucketResult><Name>%s</Name><Contents><Key>k</Key></Contents></ListBucketResult>`, bucket)
	case r.Method == http.MethodGet:
		s.mu.Lock()
		b, ok := s.objects[key]
		ct := s.types[key]
		s.mu.Unlock()
		if !ok {
			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, `<Error><Code>NoSuchKey</Code><BucketName>%s</BucketName></Error>`, bucket)
			return
		}
		w.Header().Set("Content-Type", ct)
		w.Write(b)
	default:
		w.WriteHeader(http.StatusNoContent)
	}
}

// setup starts a gateway in front of a fake upstream and returns its URL.
func setup(t *testing.T) (string, *store) {
	t.Helper()
	upIAM := iamFor(t, fmt.Sprintf(`{"identities":[{"name":"up","credentials":[{"accessKey":%q,"secretKey":%q}],"actions":["Admin"]}]}`,
		upKey.AccessKeyID, upKey.SecretAccessKey))
	st := &store{iam: upIAM, objects: map[string][]byte{}, types: map[string]string{}}
	upstream := httptest.NewServer(st)
	t.Cleanup(upstream.Close)

	// Every upstream host, bucket.s3.test included, dials the fake.
	dial := func(ctx context.Context, network, _ string) (net.Conn, error) {
		return (&net.Dialer{}).DialContext(ctx, network, upstream.Listener.Addr().String())
	}
	pool, stream := NewTransport(), NewTransport()
	pool.DialContext, stream.DialContext = dial, dial
	stream.DisableKeepAlives = true
	up, err := NewUpstream("http://s3.test", "us-east-1", upKey, &http.Client{Transport: pool}, &http.Client{Transport: stream})
	if err != nil {
		t.Fatal(err)
	}
	gwIAM := iamFor(t, fmt.Sprintf(`{"identities":[
		{"name":"hanzo","credentials":[{"accessKey":%q,"secretKey":%q}],"actions":["Admin","Read","Write","List"]},
		{"name":"backup","credentials":[{"accessKey":%q,"secretKey":%q}],"actions":["Read:backups","List:backups","Write:backups/*"]}]}`,
		adminKey.AccessKeyID, adminKey.SecretAccessKey, backupKey.AccessKeyID, backupKey.SecretAccessKey))
	g, err := New(gwIAM, up, names)
	if err != nil {
		t.Fatal(err)
	}
	front := httptest.NewServer(g)
	t.Cleanup(front.Close)
	return front.URL, st
}

// do signs a request to the gateway with key and sends it.
func do(t *testing.T, key aws.Credentials, method, url string, body []byte, header http.Header) (int, []byte, http.Header) {
	t.Helper()
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range header {
		req.Header[k] = v
	}
	sum := sha256.Sum256(body)
	hash := hex.EncodeToString(sum[:])
	req.Header.Set("X-Amz-Content-Sha256", hash)
	if key.AccessKeyID != "" {
		if err := v4.NewSigner(func(o *v4.SignerOptions) { o.DisableURIPathEscaping = true }).
			SignHTTP(context.Background(), key, req, hash, "s3", "us-east-1", time.Now().UTC()); err != nil {
			t.Fatal(err)
		}
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, b, resp.Header
}

func TestObjectRoundTripLandsUpstreamUnderTheMappedName(t *testing.T) {
	gw, st := setup(t)
	data := make([]byte, 3*chunkSize+17)
	rand.Read(data)
	code, body, _ := do(t, adminKey, http.MethodPut, gw+"/alpha/"+escapeKey("dir/a b+c.bin"), data, http.Header{"Content-Type": {"application/x-test"}})
	if code != http.StatusOK {
		t.Fatalf("put: %d %s", code, body)
	}
	if got := st.objects["hz-alpha-acct/dir/a b+c.bin"]; !bytes.Equal(got, data) {
		t.Fatalf("upstream holds %d bytes under the mapped name, want %d", len(got), len(data))
	}
	put := st.seen[len(st.seen)-1]
	if put.Header.Get("X-Amz-Content-Sha256") != trailerPayload || put.Header.Get("X-Amz-Trailer") != "x-amz-checksum-crc32" {
		t.Fatalf("a body sent without a checksum goes upstream trailed: %v", put.Header)
	}
	code, got, h := do(t, adminKey, http.MethodGet, gw+"/alpha/"+escapeKey("dir/a b+c.bin"), nil, nil)
	if code != http.StatusOK || !bytes.Equal(got, data) || h.Get("Content-Type") != "application/x-test" {
		t.Fatalf("get: %d, %d bytes, type %q", code, len(got), h.Get("Content-Type"))
	}
}

func TestTheGatewaysOwnChecksumIsNotReturned(t *testing.T) {
	gw, _ := setup(t)
	_, _, h := do(t, adminKey, http.MethodPut, gw+"/alpha/k", []byte("x"), nil)
	if v := h.Get("X-Amz-Checksum-Crc32"); v != "" {
		t.Fatalf("client sent no checksum but was told %q", v)
	}
}

func TestADeclaredHashThatDoesNotMatchStoresNothing(t *testing.T) {
	gw, st := setup(t)
	req, _ := http.NewRequest(http.MethodPut, gw+"/alpha/lie", bytes.NewReader([]byte("the real body")))
	wrong := sha256.Sum256([]byte("another body"))
	hash := hex.EncodeToString(wrong[:])
	req.Header.Set("X-Amz-Content-Sha256", hash)
	v4.NewSigner(func(o *v4.SignerOptions) { o.DisableURIPathEscaping = true }).
		SignHTTP(context.Background(), adminKey, req, hash, "s3", "us-east-1", time.Now().UTC())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		t.Fatal("a body that does not hash to its declared sha256 was accepted")
	}
	if _, ok := st.objects["hz-alpha-acct/lie"]; ok {
		t.Fatal("upstream stored a body whose hash did not match")
	}
}

func TestRefusals(t *testing.T) {
	gw, _ := setup(t)
	cases := []struct {
		name   string
		key    aws.Credentials
		method string
		path   string
		want   int
	}{
		{"unsigned", aws.Credentials{}, http.MethodGet, "/alpha/k", http.StatusForbidden},
		{"wrong secret", aws.Credentials{AccessKeyID: adminKey.AccessKeyID, SecretAccessKey: "nope"}, http.MethodGet, "/alpha/k", http.StatusForbidden},
		{"scoped list elsewhere", backupKey, http.MethodGet, "/alpha?list-type=2", http.StatusForbidden},
		{"scoped write elsewhere", backupKey, http.MethodPut, "/alpha/k", http.StatusForbidden},
		{"scoped create bucket", backupKey, http.MethodPut, "/newbucket", http.StatusForbidden},
		{"public policy", adminKey, http.MethodPut, "/alpha?policy", http.StatusNotImplemented},
		{"object acl", adminKey, http.MethodPut, "/alpha/k?acl", http.StatusNotImplemented},
		{"public access block", adminKey, http.MethodDelete, "/alpha?publicAccessBlock", http.StatusNotImplemented},
		{"dotted bucket", adminKey, http.MethodGet, "/a.b/k", http.StatusBadRequest},
	}
	for _, c := range cases {
		code, body, _ := do(t, c.key, c.method, gw+c.path, nil, nil)
		if code != c.want {
			t.Errorf("%s: %d, want %d: %s", c.name, code, c.want, body)
		}
	}
	if code, body, _ := do(t, backupKey, http.MethodPut, gw+"/backups/20260101/dump", []byte("d"), nil); code != http.StatusOK {
		t.Errorf("scoped write in its own bucket: %d %s", code, body)
	}
}

func TestListingsNameLocalBuckets(t *testing.T) {
	gw, _ := setup(t)
	code, body, _ := do(t, adminKey, http.MethodGet, gw+"/", nil, nil)
	if code != http.StatusOK || !strings.Contains(string(body), "<Name>alpha</Name>") ||
		strings.Contains(string(body), "someone-else") || strings.Contains(string(body), "hz-") {
		t.Fatalf("list buckets: %d %s", code, body)
	}
	_, body, _ = do(t, backupKey, http.MethodGet, gw+"/", nil, nil)
	if strings.Contains(string(body), "alpha") || !strings.Contains(string(body), "<Name>backups</Name>") {
		t.Fatalf("a scoped identity lists only its buckets: %s", body)
	}
	_, body, _ = do(t, adminKey, http.MethodGet, gw+"/alpha?list-type=2", nil, nil)
	if !strings.Contains(string(body), "<Name>alpha</Name>") {
		t.Fatalf("list objects names the local bucket: %s", body)
	}
	_, body, _ = do(t, adminKey, http.MethodGet, gw+"/alpha/missing", nil, nil)
	if !strings.Contains(string(body), "<BucketName>alpha</BucketName>") {
		t.Fatalf("an error names the local bucket: %s", body)
	}
}

func TestPresignedGet(t *testing.T) {
	gw, st := setup(t)
	st.objects["hz-alpha-acct/p"] = []byte("presigned")
	req, _ := http.NewRequest(http.MethodGet, gw+"/alpha/p?X-Amz-Expires=60", nil)
	signed, _, err := v4.NewSigner(func(o *v4.SignerOptions) { o.DisableURIPathEscaping = true }).
		PresignHTTP(context.Background(), adminKey, req, "UNSIGNED-PAYLOAD", "s3", "us-east-1", time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.Get(signed)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK || string(b) != "presigned" {
		t.Fatalf("presigned get: %d %s", resp.StatusCode, b)
	}
	forwarded := st.seen[len(st.seen)-1].URL.Query()
	if forwarded.Has("X-Amz-Signature") || forwarded.Has("X-Amz-Credential") {
		t.Fatalf("the client's presigning reached the upstream: %v", forwarded)
	}
}
