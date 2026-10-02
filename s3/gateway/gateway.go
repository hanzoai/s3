// Package gateway serves the S3 API from an upstream object store and keeps
// nothing of its own.
//
// A request is authenticated against the gateway's identities (the same
// s3.json the store reads), authorized for the action it performs, then signed
// again with the upstream's key and sent there under the bucket's mapped name.
// The upstream answers; the gateway renames buckets in what comes back and
// returns it. Every byte and every listing lives upstream, so any number of
// gateways serve the same data and losing one loses nothing.
package gateway

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/hanzoai/s3/s3/glog"
	"github.com/hanzoai/s3/s3/s3api"
	"github.com/hanzoai/s3/s3/s3api/s3_constants"
	"github.com/hanzoai/s3/s3/s3api/s3err"
	"github.com/hanzoai/s3/s3/util/request_id"
)

// Gateway is an http.Handler for the S3 API.
type Gateway struct {
	iam   *s3api.IdentityAccessManagement
	up    *Upstream
	names Names
}

// New returns a gateway that authenticates with iam and stores in up under
// names. iam must hold at least one identity: an IAM without any lets every
// request through as admin, which a gateway in front of real data must never do.
func New(iam *s3api.IdentityAccessManagement, up *Upstream, names Names) (*Gateway, error) {
	if len(iam.GetStaticIdentities()) == 0 {
		return nil, fmt.Errorf("no identities: refusing to serve every request as admin")
	}
	return &Gateway{iam: iam, up: up, names: names}, nil
}

// health answers probes. It is not a valid bucket name, so it shadows none.
const health = "/_health"

func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r, _ = request_id.Ensure(r)
	start := time.Now()
	rec := &recorder{ResponseWriter: w, status: http.StatusOK}
	g.serve(rec, r)
	glog.V(2).Infof("%s %s %d %dB %s", r.Method, r.URL.Path, rec.status, rec.bytes, time.Since(start).Round(time.Millisecond))
	if rec.status >= 500 {
		glog.Warningf("%s %s answered %d", r.Method, r.URL.Path, rec.status)
	}
}

func (g *Gateway) serve(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == health {
		w.WriteHeader(http.StatusOK)
		return
	}
	cors(w, r)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	bucket, key := split(r.URL.Path)
	r = mux.SetURLVars(r, map[string]string{"bucket": bucket, "object": key})

	what, action := classify(r, bucket, key)
	if what == refused {
		s3err.WriteErrorResponse(w, r, s3err.ErrNotImplemented)
		return
	}

	upBucket := ""
	if bucket != "" {
		var ok bool
		if upBucket, ok = g.names.Upstream(bucket); !ok {
			s3err.WriteErrorResponse(w, r, s3err.ErrInvalidBucketName)
			return
		}
	}

	identity, code := g.iam.AuthenticateRequest(r)
	if code != s3err.ErrNone {
		s3err.WriteErrorResponse(w, r, code)
		return
	}

	switch what {
	case listBuckets:
		g.listBuckets(w, r, identity)
		return
	case batchDelete:
		g.batchDelete(w, r, identity, bucket, upBucket)
		return
	}

	if !identity.CanDo(s3api.Action(action), bucket, key) {
		s3err.WriteErrorResponse(w, r, s3err.ErrAccessDenied)
		return
	}
	if bypassesGovernance(r) && !identity.CanDo(s3_constants.ACTION_BYPASS_GOVERNANCE_RETENTION, bucket, key) {
		s3err.WriteErrorResponse(w, r, s3err.ErrAccessDenied)
		return
	}

	copySource := ""
	if src := r.Header.Get("X-Amz-Copy-Source"); src != "" && key != "" {
		srcBucket, srcKey, rest, ok := parseCopySource(src)
		if !ok {
			s3err.WriteErrorResponse(w, r, s3err.ErrInvalidCopySource)
			return
		}
		upSrc, ok := g.names.Upstream(srcBucket)
		if !ok {
			s3err.WriteErrorResponse(w, r, s3err.ErrInvalidCopySource)
			return
		}
		if !identity.CanDo(s3_constants.ACTION_READ, srcBucket, srcKey) {
			s3err.WriteErrorResponse(w, r, s3err.ErrAccessDenied)
			return
		}
		copySource = upSrc + "/" + rest
	}

	if what == createBucket {
		g.createBucket(w, r, bucket, upBucket)
		return
	}
	g.forward(w, r, bucket, key, upBucket, copySource)
}

// split returns the bucket and key a path-style request names.
func split(path string) (bucket, key string) {
	path = strings.TrimPrefix(path, "/")
	if i := strings.IndexByte(path, '/'); i >= 0 {
		return path[:i], path[i+1:]
	}
	return path, ""
}

// parseCopySource splits an x-amz-copy-source value, "/bucket/key?versionId=v"
// with the key URL-encoded, into the bucket, the decoded key, and the rest as
// sent (encoded key and query), which is what the upstream is given after its
// own bucket name.
func parseCopySource(v string) (bucket, key, rest string, ok bool) {
	v = strings.TrimPrefix(v, "/")
	i := strings.IndexByte(v, '/')
	if i <= 0 || i == len(v)-1 {
		return "", "", "", false
	}
	bucket, rest = v[:i], v[i+1:]
	encoded := rest
	if j := strings.IndexByte(encoded, '?'); j >= 0 {
		encoded = encoded[:j]
	}
	key, err := url.PathUnescape(encoded)
	if err != nil || key == "" {
		return "", "", "", false
	}
	return bucket, key, rest, true
}

// payload is what goes upstream for a request's body.
type payload struct {
	body    io.Reader
	length  int64  // bytes on the wire, -1 when unknown
	hash    string // the x-amz-content-sha256 the upstream request is signed with
	md5     string // a Content-MD5 to add, when the client sent no checksum
	trailed int64  // decoded length when the body is sent trailed, else -1
	algo    checksum
	ours    bool // the trailer is the gateway's: the client sent no checksum
	decoded bool // the client's aws-chunked encoding was removed here
}

// smallBody bounds the bodies read whole: bucket and object configuration and
// batch deletes, which are a few KiB.
const smallBody = 4 << 20

// payload returns what to send upstream for r's body. upload is true for the
// object data of a PutObject or UploadPart.
//
// An aws-chunked body is decoded here with every chunk signature and trailer
// checksum checked. Object data the client sent no checksum for goes upstream
// trailed with a CRC32 computed here, so the upstream verifies every byte and
// accepts it into an object-locked bucket, which refuses a body without one; a
// SHA-256 the client declared is checked here as it passes. A small body is
// read whole and given a Content-MD5, which the upstream requires of
// configuration it locks or expires on. Anything else is sent as received under
// the hash the client declared.
func (g *Gateway) payload(r *http.Request, upload bool) (payload, s3err.ErrorCode) {
	declared := r.Header.Get("X-Amz-Content-Sha256")
	p := payload{trailed: -1}
	var src io.Reader = r.Body
	length := r.ContentLength
	want := ""
	if strings.HasPrefix(declared, "STREAMING-") {
		n, err := strconv.ParseInt(r.Header.Get("X-Amz-Decoded-Content-Length"), 10, 64)
		if err != nil || n < 0 {
			return p, s3err.ErrInvalidRequest
		}
		reader, code := g.iam.DecodedBody(r)
		if code != s3err.ErrNone {
			return p, code
		}
		src, length, p.decoded = reader, n, true
	} else if isHash(declared) {
		want = strings.ToLower(declared)
	}

	if length == 0 || r.Body == nil {
		p.hash = unsignedPayload
		if want != "" {
			p.hash = want
		}
		if !hasChecksum(r.Header) && (r.Method == http.MethodPut || r.Method == http.MethodPost) {
			p.md5 = contentMD5(nil)
		}
		return p, s3err.ErrNone
	}

	switch {
	case upload && !hasChecksum(r.Header) && length > 0:
		p.algo = checksumFor(r.Header.Get("X-Amz-Trailer"))
		p.ours = r.Header.Get("X-Amz-Trailer") == ""
		p.body, p.length, p.hash, p.trailed = newTrailed(src, length, want, p.algo), trailedLength(length, p.algo), trailerPayload, length
	case !upload && length > 0 && length <= smallBody:
		raw, err := io.ReadAll(io.LimitReader(src, length+1))
		if err != nil || int64(len(raw)) != length {
			return p, s3err.ErrInvalidRequest
		}
		if want != "" && sha256Hex(raw) != want {
			return p, s3err.ErrContentSHA256Mismatch
		}
		if !hasChecksum(r.Header) {
			p.md5 = contentMD5(raw)
		}
		p.body, p.length, p.hash = bytes.NewReader(raw), length, sha256Hex(raw)
	default:
		p.body, p.length, p.hash = src, length, unsignedPayload
		if want != "" {
			p.hash = want
		}
	}
	return p, s3err.ErrNone
}

// isUpload reports whether r carries object data: a PutObject or UploadPart,
// as opposed to a copy or an object's tagging, retention or legal hold.
func isUpload(r *http.Request, key string) bool {
	if r.Method != http.MethodPut || key == "" || r.Header.Get("X-Amz-Copy-Source") != "" {
		return false
	}
	q := r.URL.Query()
	return !q.Has("tagging") && !q.Has("retention") && !q.Has("legal-hold")
}

// isData reports whether a GET or HEAD of key returns the object's bytes, as
// opposed to its tagging, retention, legal hold, attributes or parts.
func isData(r *http.Request, key string) bool {
	if key == "" || (r.Method != http.MethodGet && r.Method != http.MethodHead) {
		return false
	}
	for k := range r.URL.Query() {
		switch {
		case k == "versionId", k == "partNumber", k == "x-id", strings.HasPrefix(k, "response-"):
		default:
			if !presignKey(k) {
				return false
			}
		}
	}
	return true
}

// presignKey reports whether k is a presigning query parameter.
func presignKey(k string) bool {
	for _, a := range authQuery {
		if k == a {
			return true
		}
	}
	return false
}

// forward sends r upstream and returns the answer.
func (g *Gateway) forward(w http.ResponseWriter, r *http.Request, bucket, key, upBucket, copySource string) {
	p, code := g.payload(r, isUpload(r, key))
	if code != s3err.ErrNone {
		s3err.WriteErrorResponse(w, r, code)
		return
	}
	req, err := g.up.request(r.Context(), r.Method, upBucket, key, upstreamQuery(r.URL.Query()), p.body, p.length)
	if err != nil {
		s3err.WriteErrorResponse(w, r, s3err.ErrInternalError)
		return
	}
	copyHeaders(req.Header, r.Header, p.decoded)
	if copySource != "" {
		req.Header.Set("X-Amz-Copy-Source", copySource)
	}
	if p.md5 != "" {
		req.Header.Set("Content-Md5", p.md5)
	}
	if p.trailed >= 0 {
		trail(req.Header, p.trailed, p.algo)
	}
	resp, err := g.up.send(req, p.hash)
	if err != nil {
		switch {
		case r.Context().Err() != nil:
			// The client went away; there is no one to answer, and nothing failed.
			if rec, ok := w.(*recorder); ok {
				rec.status = statusClientClosed
			}
			return
		case errors.Is(err, errPayloadHash):
			s3err.WriteErrorResponse(w, r, s3err.ErrContentSHA256Mismatch)
			return
		case errors.Is(err, errShortBody):
			s3err.WriteErrorResponse(w, r, s3err.ErrInvalidRequest)
			return
		default:
			glog.Warningf("upstream %s %s/%s: %v", r.Method, upBucket, key, err)
		}
		s3err.WriteErrorResponse(w, r, s3err.ErrServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	if p.ours {
		// The client asked for no checksum, so it is not told of the one the
		// gateway sent: a client that echoes part checksums into
		// CompleteMultipartUpload would name ones its upload never declared.
		for name := range resp.Header {
			if strings.HasPrefix(http.CanonicalHeaderKey(name), "X-Amz-Checksum-") {
				resp.Header.Del(name)
			}
		}
	}

	// Object bytes pass untouched. Anything else that is XML may name the
	// upstream bucket, so it is renamed before it goes back.
	xmlBody := strings.Contains(resp.Header.Get("Content-Type"), "xml") && r.Method != http.MethodHead
	if (isData(r, key) && resp.StatusCode < 300) || !xmlBody {
		g.reply(w, resp, resp.Body, resp.ContentLength)
		return
	}
	buf, whole, err := readBounded(resp.Body, maxRewrite)
	if err != nil {
		glog.Warningf("upstream %s %s/%s: reading reply: %v", r.Method, upBucket, key, err)
		s3err.WriteErrorResponse(w, r, s3err.ErrServiceUnavailable)
		return
	}
	if !whole {
		g.reply(w, resp, io.MultiReader(bytes.NewReader(buf), resp.Body), resp.ContentLength)
		return
	}
	location := ""
	if r.Method == http.MethodPost && r.URL.Query().Has("uploadId") {
		location = key
	}
	buf = rename(buf, g.names, bucket, location)
	g.reply(w, resp, bytes.NewReader(buf), int64(len(buf)))
}

// reply copies the upstream's status and headers, then body, to w.
func (g *Gateway) reply(w http.ResponseWriter, resp *http.Response, body io.Reader, length int64) {
	h := w.Header()
	for name, values := range resp.Header {
		if returned(name) {
			h[name] = values
		}
	}
	if length >= 0 {
		h.Set("Content-Length", strconv.FormatInt(length, 10))
	}
	w.WriteHeader(resp.StatusCode)
	if _, err := io.CopyBuffer(w, body, make([]byte, 256<<10)); err != nil {
		glog.V(1).Infof("reply: %v", err)
	}
}

// listBuckets answers GET / with the upstream buckets this mapping produces
// that the identity may list, under their local names.
func (g *Gateway) listBuckets(w http.ResponseWriter, r *http.Request, identity *s3api.Identity) {
	type bucket struct {
		Name         string `xml:"Name"`
		CreationDate string `xml:"CreationDate"`
	}
	var upstream struct {
		Buckets           []bucket `xml:"Buckets>Bucket"`
		ContinuationToken string   `xml:"ContinuationToken"`
	}
	var mine []bucket
	token := ""
	for {
		q := url.Values{"max-buckets": {"10000"}}
		if g.names.Prefix != "" {
			q.Set("prefix", g.names.Prefix)
		}
		if token != "" {
			q.Set("continuation-token", token)
		}
		req, err := g.up.request(r.Context(), http.MethodGet, "", "", q, nil, 0)
		if err != nil {
			s3err.WriteErrorResponse(w, r, s3err.ErrInternalError)
			return
		}
		resp, err := g.up.send(req, unsignedPayload)
		if err != nil {
			s3err.WriteErrorResponse(w, r, s3err.ErrServiceUnavailable)
			return
		}
		raw, _, err := readBounded(resp.Body, maxRewrite)
		resp.Body.Close()
		if err != nil || resp.StatusCode != http.StatusOK {
			glog.Warningf("upstream list buckets: %d %v %s", resp.StatusCode, err, raw)
			s3err.WriteErrorResponse(w, r, s3err.ErrServiceUnavailable)
			return
		}
		upstream.Buckets, upstream.ContinuationToken = nil, ""
		if err := xml.Unmarshal(raw, &upstream); err != nil {
			s3err.WriteErrorResponse(w, r, s3err.ErrServiceUnavailable)
			return
		}
		for _, b := range upstream.Buckets {
			local, ok := g.names.Local(b.Name)
			if !ok || !identity.CanDo(s3_constants.ACTION_LIST, local, "") {
				continue
			}
			mine = append(mine, bucket{Name: local, CreationDate: b.CreationDate})
		}
		if upstream.ContinuationToken == "" {
			break
		}
		token = upstream.ContinuationToken
	}

	type owner struct {
		ID          string `xml:"ID"`
		DisplayName string `xml:"DisplayName"`
	}
	result := struct {
		XMLName xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ ListAllMyBucketsResult"`
		Owner   owner    `xml:"Owner"`
		Buckets []bucket `xml:"Buckets>Bucket"`
	}{Owner: owner{ID: identity.Name, DisplayName: identity.Name}, Buckets: mine}
	s3err.WriteXMLResponse(w, r, http.StatusOK, result)
}

// createBucket creates the bucket upstream under its mapped name. The client's
// body (a location constraint for the gateway's region, or none) is replaced
// by the upstream's own region; object lock, if asked for, is kept.
func (g *Gateway) createBucket(w http.ResponseWriter, r *http.Request, bucket, upBucket string) {
	var body []byte
	if g.up.Region != "us-east-1" {
		body = []byte(`<CreateBucketConfiguration xmlns="http://s3.amazonaws.com/doc/2006-03-01/"><LocationConstraint>` +
			g.up.Region + `</LocationConstraint></CreateBucketConfiguration>`)
	}
	req, err := g.up.request(r.Context(), http.MethodPut, upBucket, "", nil, bytes.NewReader(body), int64(len(body)))
	if err != nil {
		s3err.WriteErrorResponse(w, r, s3err.ErrInternalError)
		return
	}
	if v := r.Header.Get("X-Amz-Bucket-Object-Lock-Enabled"); v != "" {
		req.Header.Set("X-Amz-Bucket-Object-Lock-Enabled", v)
	}
	resp, err := g.up.send(req, sha256Hex(body))
	if err != nil {
		s3err.WriteErrorResponse(w, r, s3err.ErrServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	raw, _, _ := readBounded(resp.Body, maxRewrite)
	raw = rename(raw, g.names, bucket, "")
	if resp.StatusCode < 300 {
		resp.Header.Set("Location", "/"+bucket)
	}
	g.reply(w, resp, bytes.NewReader(raw), int64(len(raw)))
}

// batchDelete authorizes a DeleteObjects request key by key, since an identity
// may hold Write on a prefix of the bucket and not the bucket, then forwards it.
func (g *Gateway) batchDelete(w http.ResponseWriter, r *http.Request, identity *s3api.Identity, bucket, upBucket string) {
	raw, whole, err := readBounded(r.Body, 4<<20)
	if err != nil || !whole {
		s3err.WriteErrorResponse(w, r, s3err.ErrMalformedXML)
		return
	}
	var doc struct {
		Objects []struct {
			Key string `xml:"Key"`
		} `xml:"Object"`
	}
	if err := xml.Unmarshal(raw, &doc); err != nil || len(doc.Objects) == 0 {
		s3err.WriteErrorResponse(w, r, s3err.ErrMalformedXML)
		return
	}
	for _, o := range doc.Objects {
		if !identity.CanDo(s3_constants.ACTION_WRITE, bucket, o.Key) {
			s3err.WriteErrorResponse(w, r, s3err.ErrAccessDenied)
			return
		}
	}
	if bypassesGovernance(r) && !identity.CanDo(s3_constants.ACTION_BYPASS_GOVERNANCE_RETENTION, bucket, "") {
		s3err.WriteErrorResponse(w, r, s3err.ErrAccessDenied)
		return
	}
	r.Body = io.NopCloser(bytes.NewReader(raw))
	r.ContentLength = int64(len(raw))
	g.forward(w, r, bucket, "", upBucket, "")
}

// cors answers cross-origin requests the way the store did: any origin may
// call, because no request is honored without a signature.
func cors(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return
	}
	h := w.Header()
	h.Set("Access-Control-Allow-Origin", origin)
	h.Set("Access-Control-Allow-Credentials", "true")
	h.Set("Access-Control-Expose-Headers", "*")
	h.Add("Vary", "Origin")
	if r.Method == http.MethodOptions {
		h.Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, HEAD")
		if req := r.Header.Get("Access-Control-Request-Headers"); req != "" {
			h.Set("Access-Control-Allow-Headers", req)
		}
		h.Set("Access-Control-Max-Age", "3600")
	}
}

// statusClientClosed marks, in the access log only, a request whose client
// disconnected before the upstream answered.
const statusClientClosed = 499

// recorder keeps the status and size of a reply for the access log.
type recorder struct {
	http.ResponseWriter
	status int
	bytes  int64
}

func (r *recorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *recorder) Write(b []byte) (int, error) {
	n, err := r.ResponseWriter.Write(b)
	r.bytes += int64(n)
	return n, err
}

func (r *recorder) Flush() {
	if f, ok := r.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}
