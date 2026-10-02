package s3api

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/gorilla/mux"
	"github.com/hanzoai/s3/s3/glog"
	"github.com/hanzoai/s3/s3/s3api/policy"
	"github.com/hanzoai/s3/s3/s3api/s3_constants"
	"github.com/hanzoai/s3/s3/s3api/s3err"
)

func (s3a *S3ApiServer) PostPolicyBucketHandler(w http.ResponseWriter, r *http.Request) {

	// https://docs.aws.amazon.com/AmazonS3/latest/API/sigv4-HTTPPOSTConstructPolicy.html
	// https://docs.aws.amazon.com/AmazonS3/latest/API/sigv4-post-example.html

	bucket := mux.Vars(r)["bucket"]

	glog.V(3).Infof("PostPolicyBucketHandler %s", bucket)

	u, ok := s3a.iam.ReadPostUpload(w, r, bucket)
	if !ok {
		return
	}
	defer u.Close()

	if err := s3a.validateTableBucketObjectPath(bucket, u.Key); err != nil {
		s3err.WriteErrorResponse(w, r, s3err.ErrAccessDenied)
		return
	}

	filePath := fmt.Sprintf("%s/%s", s3a.bucketDir(bucket), u.Key)

	r.Header.Set("Content-Type", u.ContentType)
	for k, v := range u.Header {
		r.Header[k] = v
	}

	// Use the file's size, not r.ContentLength: the multipart body wrapping form
	// fields and boundaries inflates ContentLength relative to the object body,
	// which would mis-evaluate any size-filtered rule.
	ttlSec := s3a.lifecycleTTLForObjectWrite(bucket, u.Key, u.Size)
	etag, errCode, sseMetadata := s3a.putToFiler(r, filePath, u.Body, bucket, u.Key, 1, ttlSec, nil, false)

	if errCode != s3err.ErrNone {
		s3err.WriteErrorResponse(w, r, errCode)
		return
	}

	// Include SSE response headers (important for bucket-default encryption)
	s3a.setSSEResponseHeaders(w, r, sseMetadata)
	WritePostResult(w, r, bucket, u, etag)
}

// PostUpload is one POST Object upload whose form, policy signature and policy
// conditions have been verified. What remains is to store Body under Key.
type PostUpload struct {
	Key         string
	Body        io.Reader
	Size        int64
	ContentType string
	// Header holds the form fields the stored object takes as request headers.
	Header http.Header
	// Status is success_action_status; Redirect is success_action_redirect.
	Status   string
	Redirect *url.URL

	file io.Closer
	form *multipart.Form
}

// Close releases the upload's file part and the form's spooled parts.
func (u *PostUpload) Close() {
	u.file.Close()
	u.form.RemoveAll()
}

// ReadPostUpload reads a POST Object request to bucket and verifies it: the
// multipart form, the policy signature against iam's identities (the signer
// must be allowed to write the bucket), and every policy condition, the
// content-length range included. On a refusal it writes the S3 error to w and
// answers false. The caller stores the upload and closes it.
func (iam *IdentityAccessManagement) ReadPostUpload(w http.ResponseWriter, r *http.Request, bucket string) (*PostUpload, bool) {
	reader, err := r.MultipartReader()
	if err != nil {
		s3err.WriteErrorResponse(w, r, s3err.ErrMalformedPOSTRequest)
		return nil, false
	}
	form, err := reader.ReadForm(int64(5 * humanize.MiByte))
	if err != nil {
		s3err.WriteErrorResponse(w, r, s3err.ErrMalformedPOSTRequest)
		return nil, false
	}

	fileBody, fileName, fileContentType, fileSize, formValues, err := extractPostPolicyFormValues(form)
	if err != nil {
		form.RemoveAll()
		s3err.WriteErrorResponse(w, r, s3err.ErrMalformedPOSTRequest)
		return nil, false
	}
	if fileBody == nil {
		form.RemoveAll()
		s3err.WriteErrorResponse(w, r, s3err.ErrPOSTFileRequired)
		return nil, false
	}
	u := &PostUpload{Body: fileBody, Size: fileSize, file: fileBody, form: form}
	refuse := func(code s3err.ErrorCode) (*PostUpload, bool) {
		u.Close()
		s3err.WriteErrorResponse(w, r, code)
		return nil, false
	}

	formValues.Set("Bucket", bucket)

	if fileName != "" && strings.Contains(formValues.Get("Key"), "${filename}") {
		formValues.Set("Key", strings.Replace(formValues.Get("Key"), "${filename}", fileName, -1))
	}
	rawObject := formValues.Get("Key")
	if rawObject == "" || !s3_constants.IsValidObjectKey(rawObject) {
		return refuse(s3err.ErrInvalidRequest)
	}
	u.Key = s3_constants.NormalizeObjectKey(rawObject)

	u.Status = formValues.Get("success_action_status")
	if successRedirect := formValues.Get("success_action_redirect"); successRedirect != "" {
		if u.Redirect, err = url.Parse(successRedirect); err != nil {
			return refuse(s3err.ErrMalformedPOSTRequest)
		}
	}

	// Verify policy signature.
	if errCode := iam.doesPolicySignatureMatch(formValues); errCode != s3err.ErrNone {
		return refuse(errCode)
	}

	policyBytes, err := base64.StdEncoding.DecodeString(formValues.Get("Policy"))
	if err != nil {
		return refuse(s3err.ErrMalformedPOSTRequest)
	}

	// Handle policy if it is set.
	if len(policyBytes) > 0 {
		postPolicyForm, err := policy.ParsePostPolicyForm(string(policyBytes))
		if err != nil {
			return refuse(s3err.ErrPostPolicyConditionInvalidFormat)
		}

		// Make sure formValues adhere to policy restrictions.
		if err = policy.CheckPostPolicy(formValues, postPolicyForm); err != nil {
			glog.V(3).Infof("PostPolicy check failed for bucket %s: %v", bucket, err)
			u.Close()
			s3err.WriteErrorResponseWithMessage(w, r, s3err.ErrAccessDenied, err.Error())
			return nil, false
		}

		// Ensure that the object size is within expected range, also the file size
		// should not exceed the maximum single Put size (5 GiB)
		lengthRange := postPolicyForm.Conditions.ContentLengthRange
		if lengthRange.Valid {
			if fileSize < lengthRange.Min {
				return refuse(s3err.ErrEntityTooSmall)
			}
			if fileSize > lengthRange.Max {
				return refuse(s3err.ErrEntityTooLarge)
			}
		}
	}

	// Content-Type from the form, otherwise from the file part.
	u.ContentType = formValues.Get("Content-Type")
	if u.ContentType == "" {
		u.ContentType = fileContentType
	}
	u.Header = postPolicyFormHeaders(formValues)
	return u, true
}

// WritePostResult answers a stored POST Object upload as success_action_redirect
// and success_action_status ask.
func WritePostResult(w http.ResponseWriter, r *http.Request, bucket string, u *PostUpload, etag string) {
	if u.Redirect != nil {
		redirect := *u.Redirect
		redirect.RawQuery = getRedirectPostRawQuery(bucket, u.Key, etag)
		w.Header().Set("Location", redirect.String())
		s3err.WriteEmptyResponse(w, r, http.StatusSeeOther)
		return
	}

	setEtag(w, etag)

	// Decide what http response to send depending on success_action_status parameter
	switch u.Status {
	case "201":
		resp := PostResponse{
			Bucket:   bucket,
			Key:      u.Key,
			ETag:     `"` + strings.Trim(etag, `"`) + `"`,
			Location: w.Header().Get("Location"),
		}
		s3err.WriteXMLResponse(w, r, http.StatusCreated, resp)
		s3err.PostLog(r, http.StatusCreated, s3err.ErrNone)
	case "200":
		s3err.WriteEmptyResponse(w, r, http.StatusOK)
	default:
		s3err.WriteEmptyResponse(w, r, http.StatusNoContent)
	}
}

// postPolicyReservedFormFields are multipart form fields that are part of the
// POST Object auth/policy mechanism (or handled explicitly elsewhere in the
// handler) and must not be forwarded to the upload as HTTP headers. Keys are
// already run through http.CanonicalHeaderKey before lookup.
var postPolicyReservedFormFields = map[string]struct{}{
	// POST policy signature (V2)
	"Policy":         {},
	"Signature":      {},
	"Awsaccesskeyid": {},
	// POST policy signature (V4)
	"X-Amz-Signature":      {},
	"X-Amz-Credential":     {},
	"X-Amz-Algorithm":      {},
	"X-Amz-Date":           {},
	"X-Amz-Security-Token": {},
	// Target descriptors
	"Key":    {},
	"File":   {},
	"Bucket": {},
	// Success actions (handled elsewhere in the handler)
	"Success_action_redirect": {},
	"Success_action_status":   {},
	"Redirect":                {},
	// Content-Type is resolved separately above (from form or file part)
	"Content-Type": {},
}

// postPolicyFormHeaders returns the validated POST Object form fields the
// resulting PUT carries as headers. Reserved fields that are part of the POST
// policy mechanism itself (signature, key, etc.) are skipped. The acl form field
// is translated to the X-Amz-Acl header to match how AWS promotes the form value
// to the underlying PUT.
func postPolicyFormHeaders(formValues http.Header) http.Header {
	h := http.Header{}
	for k := range formValues {
		if _, reserved := postPolicyReservedFormFields[k]; reserved {
			continue
		}
		switch {
		case k == "Acl":
			h.Set(s3_constants.AmzCannedAcl, formValues.Get(k))
		case k == "Cache-Control",
			k == "Expires",
			k == "Content-Disposition",
			k == "Content-Encoding",
			k == "Content-Language":
			h.Set(k, formValues.Get(k))
		case strings.HasPrefix(k, "X-Amz-"):
			h.Set(k, formValues.Get(k))
		}
	}
	return h
}

// Extract form fields and file data from a HTTP POST Policy
func extractPostPolicyFormValues(form *multipart.Form) (filePart io.ReadCloser, fileName, fileContentType string, fileSize int64, formValues http.Header, err error) {
	// / HTML Form values
	fileName = ""
	fileContentType = ""

	// Canonicalize the form values into http.Header.
	formValues = make(http.Header)
	for k, v := range form.Value {
		formValues[http.CanonicalHeaderKey(k)] = v
	}

	// Validate form values.
	if err = validateFormFieldSize(formValues); err != nil {
		return nil, "", "", 0, nil, err
	}

	// this means that filename="" was not specified for file key and Go has
	// an ugly way of handling this situation. Refer here
	// https://golang.org/src/mime/multipart/formdata.go#L61
	if len(form.File) == 0 {
		var b = &bytes.Buffer{}
		for _, v := range formValues["File"] {
			b.WriteString(v)
		}
		fileSize = int64(b.Len())
		filePart = io.NopCloser(b)
		return filePart, fileName, fileContentType, fileSize, formValues, nil
	}

	// Iterator until we find a valid File field and break
	for k, v := range form.File {
		canonicalFormName := http.CanonicalHeaderKey(k)
		if canonicalFormName == "File" {
			if len(v) == 0 {
				return nil, "", "", 0, nil, errors.New("Invalid arguments specified")
			}
			// Fetch fileHeader which has the uploaded file information
			fileHeader := v[0]
			// Set filename
			fileName = fileHeader.Filename
			// Set contentType
			fileContentType = fileHeader.Header.Get("Content-Type")
			// Open the uploaded part
			filePart, err = fileHeader.Open()
			if err != nil {
				return nil, "", "", 0, nil, err
			}
			// Compute file size
			fileSize, err = filePart.(io.Seeker).Seek(0, 2)
			if err != nil {
				return nil, "", "", 0, nil, err
			}
			// Reset Seek to the beginning
			_, err = filePart.(io.Seeker).Seek(0, 0)
			if err != nil {
				return nil, "", "", 0, nil, err
			}
			// File found and ready for reading
			break
		}
	}
	return filePart, fileName, fileContentType, fileSize, formValues, nil
}

// Validate form field size for s3 specification requirement.
func validateFormFieldSize(formValues http.Header) error {
	// Iterate over form values
	for k := range formValues {
		// Check if value's field exceeds S3 limit
		if int64(len(formValues.Get(k))) > int64(1*humanize.MiByte) {
			return errors.New("Data size larger than expected")
		}
	}

	// Success.
	return nil
}

func getRedirectPostRawQuery(bucket, key, etag string) string {
	redirectValues := make(url.Values)
	redirectValues.Set("bucket", bucket)
	redirectValues.Set("key", key)
	redirectValues.Set("etag", "\""+etag+"\"")
	return redirectValues.Encode()
}

// Check to see if Policy is signed correctly.
func (iam *IdentityAccessManagement) doesPolicySignatureMatch(formValues http.Header) s3err.ErrorCode {
	// For SignV2 - Signature field will be valid
	if _, ok := formValues["Signature"]; ok {
		return iam.doesPolicySignatureV2Match(formValues)
	}
	return iam.doesPolicySignatureV4Match(formValues)
}
