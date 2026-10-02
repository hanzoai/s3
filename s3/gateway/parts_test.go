package gateway

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/hanzoai/s3/s3/s3api/s3_constants"
)

func TestNamesWrapAndUnwrap(t *testing.T) {
	n := Names{Prefix: "hanzo-s3-", Suffix: "-532217001883"}
	up, ok := n.Upstream("org-db")
	if !ok || up != "hanzo-s3-org-db-532217001883" {
		t.Fatalf("upstream: %q %v", up, ok)
	}
	if local, ok := n.Local(up); !ok || local != "org-db" {
		t.Fatalf("local: %q %v", local, ok)
	}
	for _, bad := range []string{"", "ab", "Org", "a.b", "-ab", "ab-", "a_b", strings.Repeat("a", 42)} {
		if _, ok := n.Upstream(bad); ok {
			t.Errorf("%q mapped", bad)
		}
	}
	for _, foreign := range []string{"hanzo-sites", "hanzo-s3--532217001883", "hanzo-s3-x-y-other"} {
		if _, ok := n.Local(foreign); ok {
			t.Errorf("%q read as ours", foreign)
		}
	}
}

func TestClassify(t *testing.T) {
	cases := []struct {
		method, target string
		kind           kind
		action         string
	}{
		{"GET", "/", listBuckets, s3_constants.ACTION_LIST},
		{"GET", "/b?list-type=2", forward, s3_constants.ACTION_LIST},
		{"GET", "/b?versions", forward, s3_constants.ACTION_LIST},
		{"GET", "/b?versioning", forward, s3_constants.ACTION_READ},
		{"HEAD", "/b", forward, s3_constants.ACTION_READ},
		{"PUT", "/b", createBucket, s3_constants.ACTION_ADMIN},
		{"PUT", "/b?object-lock", forward, s3_constants.ACTION_ADMIN},
		{"PUT", "/b?policy", refused, ""},
		{"PUT", "/b?acl", refused, ""},
		{"DELETE", "/b?policy", forward, s3_constants.ACTION_ADMIN},
		{"DELETE", "/b?publicAccessBlock", refused, ""},
		{"DELETE", "/b", forward, s3_constants.ACTION_DELETE_BUCKET},
		{"POST", "/b?delete", batchDelete, s3_constants.ACTION_WRITE},
		{"POST", "/b", refused, ""},
		{"GET", "/b/k", forward, s3_constants.ACTION_READ},
		{"GET", "/b/k?retention", forward, s3_constants.ACTION_GET_OBJECT_RETENTION},
		{"PUT", "/b/k?partNumber=1&uploadId=u", forward, s3_constants.ACTION_WRITE},
		{"PUT", "/b/k?tagging", forward, s3_constants.ACTION_TAGGING},
		{"PUT", "/b/k?acl", refused, ""},
		{"POST", "/b/k?uploads", forward, s3_constants.ACTION_WRITE},
		{"DELETE", "/b/k?uploadId=u", forward, s3_constants.ACTION_WRITE},
	}
	for _, c := range cases {
		r := httptest.NewRequest(c.method, c.target, nil)
		bucket, key := split(r.URL.Path)
		k, a := classify(r, bucket, key)
		if k != c.kind || a != c.action {
			t.Errorf("%s %s: kind %d action %q, want %d %q", c.method, c.target, k, a, c.kind, c.action)
		}
	}
}

func TestParseCopySource(t *testing.T) {
	b, k, rest, ok := parseCopySource("/src/dir/a%20b?versionId=v1")
	if !ok || b != "src" || k != "dir/a b" || rest != "dir/a%20b?versionId=v1" {
		t.Fatalf("%q %q %q %v", b, k, rest, ok)
	}
	for _, bad := range []string{"", "/", "src", "src/", "/src/"} {
		if _, _, _, ok := parseCopySource(bad); ok {
			t.Errorf("%q parsed", bad)
		}
	}
}

// TestTrailedEncodesWhatItAnnounces checks the encoded length against the
// bytes produced and decodes the stream back, across chunk boundaries.
func TestTrailedEncodesWhatItAnnounces(t *testing.T) {
	for _, algo := range []string{"", "x-amz-checksum-crc32c", "x-amz-checksum-crc64nvme", "x-amz-checksum-sha256"} {
		c := checksumFor(algo)
		for _, n := range []int64{1, chunkSize - 1, chunkSize, chunkSize + 1, 3 * chunkSize, 3*chunkSize + 5} {
			data := make([]byte, n)
			rand.Read(data)
			enc, err := io.ReadAll(newTrailed(bytes.NewReader(data), n, "", c))
			if err != nil {
				t.Fatal(err)
			}
			if int64(len(enc)) != trailedLength(n, c) {
				t.Fatalf("%s n=%d: %d bytes, announced %d", c.name, n, len(enc), trailedLength(n, c))
			}
			got, trailer := decodeTrailed(t, enc)
			if !bytes.Equal(got, data) {
				t.Fatalf("%s n=%d: decoded body differs", c.name, n)
			}
			h := c.new()
			h.Write(data)
			if want := c.trailer + ":" + b64(h.Sum(nil)); trailer != want {
				t.Fatalf("%s n=%d: trailer %q, want %q", c.name, n, trailer, want)
			}
		}
	}
}

func TestTrailedRefusesABodyThatMissesItsHash(t *testing.T) {
	data := []byte("body")
	sum := sha256.Sum256(data)
	if _, err := io.ReadAll(newTrailed(bytes.NewReader(data), 4, hex.EncodeToString(sum[:]), checksumFor(""))); err != nil {
		t.Fatalf("matching hash: %v", err)
	}
	if _, err := io.ReadAll(newTrailed(bytes.NewReader(data), 4, strings.Repeat("0", 64), checksumFor(""))); err != errPayloadHash {
		t.Fatalf("mismatched hash: %v", err)
	}
	if _, err := io.ReadAll(newTrailed(bytes.NewReader(data), 9, "", checksumFor(""))); err != errShortBody {
		t.Fatalf("short body: %v", err)
	}
}

func TestHeadersThatStay(t *testing.T) {
	in := http.Header{
		"Authorization":                {"AWS4-HMAC-SHA256 ..."},
		"X-Amz-Date":                   {"20260101T000000Z"},
		"X-Amz-Acl":                    {"public-read"},
		"X-Amz-Grant-Read":             {"uri=AllUsers"},
		"X-Amz-Expected-Bucket-Owner":  {"1"},
		"S3-Auth-Type":                 {"SigV4"},
		"Content-Encoding":             {"aws-chunked,gzip"},
		"X-Amz-Meta-Origin":            {"a"},
		"X-Amz-Sdk-Checksum-Algorithm": {"CRC32"},
		"Cache-Control":                {"max-age=60"},
	}
	out := http.Header{}
	copyHeaders(out, in, true)
	for _, gone := range []string{"Authorization", "X-Amz-Date", "X-Amz-Acl", "X-Amz-Grant-Read", "X-Amz-Expected-Bucket-Owner", "S3-Auth-Type", "X-Amz-Sdk-Checksum-Algorithm"} {
		if out.Get(gone) != "" {
			t.Errorf("%s forwarded", gone)
		}
	}
	if out.Get("Content-Encoding") != "gzip" || out.Get("X-Amz-Meta-Origin") != "a" || out.Get("Cache-Control") != "max-age=60" {
		t.Errorf("kept headers: %v", out)
	}
}

func TestRenameTouchesOnlyMappedNames(t *testing.T) {
	body := []byte(`<R><Name>hz-alpha-acct</Name><Bucket>hz-alpha-acct</Bucket><Key>hz-alpha-acct</Key><BucketName>other</BucketName><Location>https://hz-alpha-acct.s3.amazonaws.com/k</Location></R>`)
	got := string(rename(body, names, "alpha", "dir/k 1"))
	want := `<R><Name>alpha</Name><Bucket>alpha</Bucket><Key>hz-alpha-acct</Key><BucketName>other</BucketName><Location>/alpha/dir/k%201</Location></R>`
	if got != want {
		t.Fatalf("\n got %s\nwant %s", got, want)
	}
}

// decodeTrailed parses an unsigned aws-chunked stream.
func decodeTrailed(t *testing.T, enc []byte) ([]byte, string) {
	t.Helper()
	r := bufio.NewReader(bytes.NewReader(enc))
	var out []byte
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			t.Fatal(err)
		}
		size, err := strconv.ParseInt(strings.TrimSuffix(line, "\r\n"), 16, 64)
		if err != nil {
			t.Fatal(err)
		}
		if size == 0 {
			break
		}
		chunk := make([]byte, size+2)
		if _, err := io.ReadFull(r, chunk); err != nil {
			t.Fatal(err)
		}
		out = append(out, chunk[:size]...)
	}
	trailer, _ := r.ReadString('\n')
	end, _ := r.ReadString('\n')
	if end != "\r\n" {
		t.Fatalf("stream does not end with an empty line: %q", end)
	}
	if rest, _ := io.ReadAll(r); len(rest) != 0 {
		t.Fatalf("%d bytes after the end", len(rest))
	}
	return out, strings.TrimSuffix(trailer, "\r\n")
}

func b64(b []byte) string { return base64.StdEncoding.EncodeToString(b) }
