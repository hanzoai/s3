package gateway

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/hanzoai/s3/s3/s3api/s3_constants"
)

// kind is what a request does, as far as the gateway has to know: which
// identity action it needs, and whether it is answered here or upstream.
type kind int

const (
	forward      kind = iota // signed again and sent upstream
	listBuckets              // GET /: answered from the upstream's bucket list
	createBucket             // PUT /<bucket>: created upstream under its mapped name
	batchDelete              // POST /<bucket>?delete: every key in the body is authorized
	postObject               // POST /<bucket> with a form: the form's own policy signs it
	refused                  // not offered by this gateway
)

// bucketReadable are the bucket subresources a GET or HEAD reads.
var bucketReadable = map[string]bool{
	"acl": true, "cors": true, "encryption": true, "lifecycle": true, "location": true,
	"logging": true, "notification": true, "object-lock": true, "ownershipControls": true,
	"policy": true, "policyStatus": true, "publicAccessBlock": true, "replication": true,
	"requestPayment": true, "tagging": true, "versioning": true, "website": true,
	"accelerate": true, "analytics": true, "intelligent-tiering": true, "inventory": true,
	"metrics": true,
}

// bucketWritable are the bucket subresources an administrator may set or
// delete. What is not listed (acl, policy, publicAccessBlock, ownershipControls,
// website, replication, notification, logging and the rest) could open a bucket
// to the internet or send its objects elsewhere, so the gateway refuses it.
var bucketWritable = map[string]bool{
	"cors": true, "encryption": true, "lifecycle": true, "object-lock": true,
	"tagging": true, "versioning": true,
}

// bucketDeletable are the bucket subresources an administrator may delete.
// Deleting a policy only ever narrows access, so it is allowed here although
// setting one is not.
var bucketDeletable = map[string]bool{
	"cors": true, "encryption": true, "lifecycle": true, "policy": true, "tagging": true,
}

// subresource returns the first query key in q that names a subresource from
// set, or "".
func subresource(q url.Values, set map[string]bool) string {
	for k := range q {
		if set[k] {
			return k
		}
	}
	return ""
}

// classify returns what r does and the action an identity needs for it. key
// is "" for a bucket-level request; bucket is "" for a service-level one.
func classify(r *http.Request, bucket, key string) (kind, string) {
	q := r.URL.Query()
	has := func(k string) bool { _, ok := q[k]; return ok }

	if bucket == "" {
		if r.Method == http.MethodGet {
			return listBuckets, s3_constants.ACTION_LIST
		}
		return refused, ""
	}

	if key == "" {
		switch r.Method {
		case http.MethodGet, http.MethodHead:
			if subresource(q, bucketReadable) != "" {
				return forward, s3_constants.ACTION_READ
			}
			if r.Method == http.MethodHead {
				return forward, s3_constants.ACTION_READ
			}
			return forward, s3_constants.ACTION_LIST
		case http.MethodPut:
			if len(q) == 0 {
				return createBucket, s3_constants.ACTION_ADMIN
			}
			if sub := subresource(q, bucketReadable); sub != "" && bucketWritable[sub] {
				return forward, s3_constants.ACTION_ADMIN
			}
			return refused, ""
		case http.MethodDelete:
			if len(q) == 0 {
				return forward, s3_constants.ACTION_DELETE_BUCKET
			}
			if sub := subresource(q, bucketReadable); sub != "" && bucketDeletable[sub] {
				return forward, s3_constants.ACTION_ADMIN
			}
			return refused, ""
		case http.MethodPost:
			if has("delete") {
				return batchDelete, s3_constants.ACTION_WRITE
			}
			if len(q) == 0 && strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
				return postObject, s3_constants.ACTION_WRITE
			}
			return refused, ""
		}
		return refused, ""
	}

	switch r.Method {
	case http.MethodGet, http.MethodHead:
		switch {
		case has("tagging"):
			return forward, s3_constants.ACTION_TAGGING
		case has("retention"):
			return forward, s3_constants.ACTION_GET_OBJECT_RETENTION
		case has("legal-hold"):
			return forward, s3_constants.ACTION_GET_OBJECT_LEGAL_HOLD
		case has("uploadId"):
			return forward, s3_constants.ACTION_WRITE
		}
		return forward, s3_constants.ACTION_READ
	case http.MethodPut:
		switch {
		case has("acl"):
			return refused, ""
		case has("tagging"):
			return forward, s3_constants.ACTION_TAGGING
		case has("retention"):
			return forward, s3_constants.ACTION_PUT_OBJECT_RETENTION
		case has("legal-hold"):
			return forward, s3_constants.ACTION_PUT_OBJECT_LEGAL_HOLD
		}
		return forward, s3_constants.ACTION_WRITE
	case http.MethodDelete:
		if has("tagging") {
			return forward, s3_constants.ACTION_TAGGING
		}
		return forward, s3_constants.ACTION_WRITE
	case http.MethodPost:
		switch {
		case has("uploads"), has("uploadId"), has("restore"):
			return forward, s3_constants.ACTION_WRITE
		case has("select"):
			return forward, s3_constants.ACTION_READ
		}
	}
	return refused, ""
}

// bypassesGovernance reports whether r asks to remove or shorten a retention
// that GOVERNANCE mode holds, which needs its own grant.
func bypassesGovernance(r *http.Request) bool {
	return strings.EqualFold(strings.TrimSpace(r.Header.Get("X-Amz-Bypass-Governance-Retention")), "true")
}
