package gateway

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
)

// postForm builds a SigV4 POST Object form for bucket, signed by key, whose
// policy admits keys under prefix, and returns its body and content type.
func postForm(t *testing.T, key aws.Credentials, bucket, prefix, objectKey string, file []byte, tamper bool) (*bytes.Buffer, string) {
	t.Helper()
	now := time.Now().UTC()
	date := now.Format("20060102")
	credential := fmt.Sprintf("%s/%s/us-east-1/s3/aws4_request", key.AccessKeyID, date)
	amzDate := now.Format("20060102T150405Z")
	policy := fmt.Sprintf(`{"expiration":%q,"conditions":[{"bucket":%q},["starts-with","$key",%q],{"x-amz-algorithm":"AWS4-HMAC-SHA256"},{"x-amz-credential":%q},{"x-amz-date":%q},["starts-with","$Content-Type",""],["content-length-range",0,1048576]]}`,
		now.Add(time.Hour).Format("2006-01-02T15:04:05.000Z"), bucket, prefix, credential, amzDate)
	encoded := base64.StdEncoding.EncodeToString([]byte(policy))
	mac := func(k []byte, s string) []byte { h := hmac.New(sha256.New, k); h.Write([]byte(s)); return h.Sum(nil) }
	signing := mac(mac(mac(mac([]byte("AWS4"+key.SecretAccessKey), date), "us-east-1"), "s3"), "aws4_request")
	signature := hex.EncodeToString(mac(signing, encoded))
	if tamper {
		signature = strings.Repeat("0", len(signature))
	}

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	for _, f := range [][2]string{
		{"key", objectKey},
		{"Content-Type", "text/css"},
		{"success_action_status", "201"},
		{"x-amz-algorithm", "AWS4-HMAC-SHA256"},
		{"x-amz-credential", credential},
		{"x-amz-date", amzDate},
		{"policy", encoded},
		{"x-amz-signature", signature},
	} {
		mw.WriteField(f[0], f[1])
	}
	fw, _ := mw.CreateFormFile("file", "site.css")
	fw.Write(file)
	mw.Close()
	return &buf, mw.FormDataContentType()
}

func post(t *testing.T, url string, body io.Reader, contentType string) (int, string) {
	t.Helper()
	resp, err := http.Post(url, contentType, body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(b)
}

// A presigned POST grant stores its file upstream under the mapped bucket and
// answers as success_action_status asks.
func TestAPresignedPostLandsUpstream(t *testing.T) {
	gw, st := setup(t)
	file := []byte("body{color:red}")
	body, ct := postForm(t, adminKey, "alpha", "hanzo/site/", "hanzo/site/_next/static/a.css", file, false)
	code, reply := post(t, gw+"/alpha", body, ct)
	if code != http.StatusCreated || !strings.Contains(reply, "<Key>hanzo/site/_next/static/a.css</Key>") {
		t.Fatalf("post: %d %s", code, reply)
	}
	if got := st.objects["hz-alpha-acct/hanzo/site/_next/static/a.css"]; !bytes.Equal(got, file) {
		t.Fatalf("upstream holds %q, want %q", got, file)
	}
	if got := st.types["hz-alpha-acct/hanzo/site/_next/static/a.css"]; got != "text/css" {
		t.Fatalf("upstream content type %q, want text/css", got)
	}
}

// What the policy does not sign stores nothing.
func TestAPostThePolicyDoesNotAdmitStoresNothing(t *testing.T) {
	gw, st := setup(t)
	for name, tc := range map[string]struct {
		key            aws.Credentials
		bucket, object string
		tamper         bool
	}{
		"forged signature":       {adminKey, "alpha", "hanzo/site/x.css", true},
		"key outside the prefix": {adminKey, "alpha", "other/x.css", false},
		"signer cannot write":    {backupKey, "alpha", "hanzo/site/x.css", false},
	} {
		body, ct := postForm(t, tc.key, tc.bucket, "hanzo/site/", tc.object, []byte("x"), tc.tamper)
		if code, reply := post(t, gw+"/"+tc.bucket, body, ct); code != http.StatusForbidden {
			t.Errorf("%s: %d %s, want 403", name, code, reply)
		}
	}
	if len(st.objects) != 0 {
		t.Fatalf("refused posts stored %d objects", len(st.objects))
	}
}
