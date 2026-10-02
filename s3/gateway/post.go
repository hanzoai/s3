package gateway

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/hanzoai/s3/s3/glog"
	"github.com/hanzoai/s3/s3/s3api"
	"github.com/hanzoai/s3/s3/s3api/s3err"
)

// postObject stores a POST Object upload: a browser form, or a presigned POST
// grant. The form, its policy signature and its conditions are verified against
// the gateway's identities exactly as the store verifies them
// (s3api.ReadPostUpload), and the file goes upstream as a PUT under the form's
// key, trailed with a CRC32 like any other object data the client sent no
// checksum for.
func (g *Gateway) postObject(w http.ResponseWriter, r *http.Request, bucket, upBucket string) {
	u, ok := g.iam.ReadPostUpload(w, r, bucket)
	if !ok {
		return
	}
	defer u.Close()

	var (
		body   io.Reader
		length int64
		hash   = unsignedPayload
		algo   = checksumFor("")
	)
	if u.Size > 0 {
		body, length, hash = newTrailed(u.Body, u.Size, "", algo), trailedLength(u.Size, algo), trailerPayload
	}
	req, err := g.up.request(r.Context(), http.MethodPut, upBucket, u.Key, nil, body, length)
	if err != nil {
		s3err.WriteErrorResponse(w, r, s3err.ErrInternalError)
		return
	}
	copyHeaders(req.Header, u.Header, false)
	req.Header.Set("Content-Type", u.ContentType)
	if u.Size > 0 {
		trail(req.Header, u.Size, algo)
	} else {
		req.Header.Set("Content-Md5", contentMD5(nil))
	}

	resp, err := g.up.send(req, hash)
	if err != nil {
		if r.Context().Err() != nil {
			if rec, ok := w.(*recorder); ok {
				rec.status = statusClientClosed
			}
			return
		}
		glog.Warningf("upstream POST %s/%s: %v", upBucket, u.Key, err)
		s3err.WriteErrorResponse(w, r, s3err.ErrServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		buf, whole, err := readBounded(resp.Body, maxRewrite)
		if err != nil || !whole {
			s3err.WriteErrorResponse(w, r, s3err.ErrServiceUnavailable)
			return
		}
		buf = rename(buf, g.names, bucket, "")
		g.reply(w, resp, bytes.NewReader(buf), int64(len(buf)))
		return
	}
	io.Copy(io.Discard, resp.Body)
	s3api.WritePostResult(w, r, bucket, u, strings.Trim(resp.Header.Get("ETag"), `"`))
}
