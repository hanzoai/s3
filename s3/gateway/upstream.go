package gateway

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
)

// unsignedPayload is the payload hash for a body sent over TLS without one.
const unsignedPayload = "UNSIGNED-PAYLOAD"

// Upstream is the object store that holds every byte: its endpoint, region and
// the key the gateway signs with.
//
// Requests go out on two pools. A request whose body can be sent again (none,
// or bytes held here) reuses connections, and is retried on a fresh one when
// the upstream has dropped the one it was given: S3 closes kept-alive
// connections when it chooses, and the transport retries an idempotent request
// that met one. A request streaming a client's body cannot be sent twice, so it
// never takes a connection that may already be closing: it opens its own.
type Upstream struct {
	Scheme string // "https"
	Host   string // "s3.us-east-1.amazonaws.com"; a bucket is addressed as <bucket>.<Host>
	Region string
	Key    aws.Credentials
	Client *http.Client // kept-alive connections
	Stream *http.Client // one connection per request
	signer *v4.Signer
}

// NewUpstream parses endpoint ("https://s3.us-east-1.amazonaws.com") and
// returns an upstream that signs with key.
func NewUpstream(endpoint, region string, key aws.Credentials, client, stream *http.Client) (*Upstream, error) {
	u, err := url.Parse(endpoint)
	if err != nil || u.Host == "" || (u.Scheme != "https" && u.Scheme != "http") {
		return nil, fmt.Errorf("upstream %q: want scheme://host", endpoint)
	}
	if region == "" {
		return nil, fmt.Errorf("upstream region is empty")
	}
	if key.AccessKeyID == "" || key.SecretAccessKey == "" {
		return nil, fmt.Errorf("upstream key is empty")
	}
	if client == nil {
		client = &http.Client{Transport: NewTransport()}
	}
	if stream == nil {
		t := NewTransport()
		t.DisableKeepAlives = true
		stream = &http.Client{Transport: t}
	}
	// The path is escaped once, by escapeKey, and signed as sent: S3 does not
	// normalize paths, so escaping it a second time would sign a different one.
	signer := v4.NewSigner(func(o *v4.SignerOptions) { o.DisableURIPathEscaping = true })
	return &Upstream{Scheme: u.Scheme, Host: u.Host, Region: region, Key: key, Client: client, Stream: stream, signer: signer}, nil
}

// NewTransport is the connection pool the gateway keeps to the upstream. It
// never asks for compression: Go would otherwise add Accept-Encoding and
// decompress the body, so the bytes and length a client receives would not be
// the object's.
func NewTransport() *http.Transport {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.Proxy = nil
	t.DisableCompression = true
	t.MaxIdleConns = 1024
	t.MaxIdleConnsPerHost = 256
	t.IdleConnTimeout = 90 * time.Second
	t.TLSHandshakeTimeout = 10 * time.Second
	t.ExpectContinueTimeout = time.Second
	t.ForceAttemptHTTP2 = false
	return t
}

// request builds the upstream request for bucket (already mapped, "" for the
// service) and key, carrying query and body. The caller adds headers and then
// calls send.
func (u *Upstream) request(ctx context.Context, method, bucket, key string, query url.Values, body io.Reader, length int64) (*http.Request, error) {
	host := u.Host
	if bucket != "" {
		host = bucket + "." + u.Host
	}
	target := &url.URL{Scheme: u.Scheme, Host: host, Path: "/" + key, RawPath: "/" + escapeKey(key)}
	target.RawQuery = encodeQuery(query)
	req, err := http.NewRequestWithContext(ctx, method, target.String(), nil)
	if err != nil {
		return nil, err
	}
	// NewRequest re-parses the URL; set the parts again so what is sent is
	// exactly what is signed.
	req.URL = target
	req.Host = host
	if body != nil && length != 0 {
		req.Body = io.NopCloser(body)
		req.ContentLength = length
		if length < 0 {
			req.ContentLength = -1
		}
		if held, ok := body.(*bytes.Reader); ok {
			start := *held
			req.GetBody = func() (io.ReadCloser, error) {
				again := start
				return io.NopCloser(&again), nil
			}
		}
	}
	return req, nil
}

// send signs req with payloadHash and sends it.
func (u *Upstream) send(req *http.Request, payloadHash string) (*http.Response, error) {
	req.Header.Set("X-Amz-Content-Sha256", payloadHash)
	req.Header.Set("User-Agent", "hanzo-s3-gateway")
	if err := u.signer.SignHTTP(req.Context(), u.Key, req, payloadHash, "s3", u.Region, time.Now().UTC()); err != nil {
		return nil, err
	}
	resendable := req.Body == nil || req.GetBody != nil
	if !resendable {
		return u.Stream.Do(req)
	}
	switch req.Method {
	case http.MethodPut, http.MethodDelete:
		// Idempotent, so the transport may send it again on a fresh connection.
		// A nil value marks it without sending a header the signature does not
		// cover.
		req.Header["Idempotency-Key"] = nil
	}
	return u.Client.Do(req)
}

// escapeKey escapes a key the way S3 canonicalizes a path: every byte but the
// unreserved ones and "/" as %XX.
func escapeKey(key string) string {
	var b strings.Builder
	for i := 0; i < len(key); i++ {
		c := key[i]
		if 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || '0' <= c && c <= '9' ||
			c == '-' || c == '_' || c == '.' || c == '~' || c == '/' {
			b.WriteByte(c)
			continue
		}
		fmt.Fprintf(&b, "%%%02X", c)
	}
	return b.String()
}

// encodeQuery encodes q in the canonical form the signer uses, so the query
// sent is the query signed.
func encodeQuery(q url.Values) string {
	return strings.ReplaceAll(q.Encode(), "+", "%20")
}

// authQuery are the query parameters of a presigned request: they sign the
// request to the gateway and mean nothing upstream.
var authQuery = []string{
	"X-Amz-Algorithm", "X-Amz-Credential", "X-Amz-Date", "X-Amz-Expires",
	"X-Amz-SignedHeaders", "X-Amz-Signature", "X-Amz-Security-Token",
	"X-Amz-Content-Sha256", "AWSAccessKeyId", "Signature", "Expires",
}

// upstreamQuery is the client's query without its presigning parameters.
func upstreamQuery(q url.Values) url.Values {
	out := make(url.Values, len(q))
	for k, v := range q {
		out[k] = v
	}
	for _, k := range authQuery {
		out.Del(k)
	}
	return out
}

// forwarded are the request headers outside x-amz-* that reach the upstream.
// Everything else (Authorization, Host, hop-by-hop, the server's own markers)
// stays here.
var forwarded = map[string]bool{
	"Cache-Control": true, "Content-Disposition": true, "Content-Encoding": true,
	"Content-Language": true, "Content-Md5": true, "Content-Type": true, "Expires": true,
	"If-Match": true, "If-Modified-Since": true, "If-None-Match": true,
	"If-Unmodified-Since": true, "Range": true,
}

// withheld are x-amz-* request headers that never reach the upstream: the
// client's own signing, encoding the gateway has already undone, and ACLs and
// account checks that name the gateway's account rather than the client's.
var withheld = map[string]bool{
	"X-Amz-Date": true, "X-Amz-Content-Sha256": true, "X-Amz-Security-Token": true,
	"X-Amz-Decoded-Content-Length": true, "X-Amz-Trailer": true,
	"X-Amz-Expected-Bucket-Owner": true, "X-Amz-Source-Expected-Bucket-Owner": true,
	"X-Amz-Acl": true, "X-Amz-Object-Ownership": true, "X-Amz-User-Agent": true,
	"X-Amz-Copy-Source": true,
}

// copyHeaders copies the headers of in that the upstream should see onto out.
// decoded is true when the gateway has already removed an aws-chunked
// encoding, so the encoding and its checksum trailer are dropped with it.
func copyHeaders(out, in http.Header, decoded bool) {
	for name, values := range in {
		canon := http.CanonicalHeaderKey(name)
		switch {
		case forwarded[canon]:
		case strings.HasPrefix(canon, "X-Amz-") && !withheld[canon] && !strings.HasPrefix(canon, "X-Amz-Grant-"):
			if decoded && canon == "X-Amz-Sdk-Checksum-Algorithm" {
				continue
			}
		default:
			continue
		}
		if canon == "Content-Encoding" {
			values = withoutChunked(values)
			if len(values) == 0 {
				continue
			}
		}
		out[canon] = append([]string(nil), values...)
	}
}

// withoutChunked removes the aws-chunked token from Content-Encoding values.
func withoutChunked(values []string) []string {
	var out []string
	for _, v := range values {
		var keep []string
		for _, tok := range strings.Split(v, ",") {
			tok = strings.TrimSpace(tok)
			if tok != "" && !strings.EqualFold(tok, "aws-chunked") {
				keep = append(keep, tok)
			}
		}
		if len(keep) > 0 {
			out = append(out, strings.Join(keep, ","))
		}
	}
	return out
}

// returned reports whether an upstream response header goes back to the
// client. Hop-by-hop headers and the upstream's server and CORS headers stay:
// the gateway answers CORS itself.
func returned(name string) bool {
	switch http.CanonicalHeaderKey(name) {
	case "Connection", "Keep-Alive", "Proxy-Authenticate", "Proxy-Authorization",
		"Te", "Trailer", "Transfer-Encoding", "Upgrade", "Server", "Content-Length",
		"X-Amz-Bucket-Arn":
		return false
	}
	return !strings.HasPrefix(http.CanonicalHeaderKey(name), "Access-Control-")
}

// nameTags are the XML elements whose text is a bucket name.
var nameTags = regexp.MustCompile(`<(Name|Bucket|BucketName)>([^<]*)</(Name|Bucket|BucketName)>`)

// locationTag is the Location element of a CompleteMultipartUploadResult.
var locationTag = regexp.MustCompile(`<Location>[^<]*</Location>`)

// rename rewrites every upstream bucket name in an XML body to its local
// name, and a Location element to the gateway's path for bucket and key.
func rename(body []byte, names Names, bucket, key string) []byte {
	body = nameTags.ReplaceAllFunc(body, func(m []byte) []byte {
		sub := nameTags.FindSubmatch(m)
		if string(sub[1]) != string(sub[3]) {
			return m
		}
		local, ok := names.Local(string(sub[2]))
		if !ok {
			return m
		}
		return []byte("<" + string(sub[1]) + ">" + local + "</" + string(sub[1]) + ">")
	})
	if key != "" {
		body = locationTag.ReplaceAll(body, []byte("<Location>/"+bucket+"/"+escapeKey(key)+"</Location>"))
	}
	return body
}

// maxRewrite bounds the XML bodies the gateway buffers to rename. A listing of
// 1000 keys is a few hundred KiB; anything larger than this passes untouched.
const maxRewrite = 32 << 20

// readBounded reads all of r when it is at most limit bytes and returns the
// bytes and true; otherwise it returns what it read and false, with the rest
// of r unread.
func readBounded(r io.Reader, limit int64) ([]byte, bool, error) {
	var buf bytes.Buffer
	n, err := io.Copy(&buf, io.LimitReader(r, limit+1))
	if err != nil {
		return buf.Bytes(), false, err
	}
	return buf.Bytes(), n <= limit, nil
}
