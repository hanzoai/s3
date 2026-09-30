package command

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

// Apache-2.0 §4(c) per-file copyright headers that must stay verbatim.
var retainedHeaders = map[string]bool{
	"s3/s3api/chunked_reader_v4.go":       true,
	"s3/s3api/policy/post-policy_test.go": true,
	"s3/s3api/policy/postpolicyform.go":   true,
	"s3/s3api/s3_constants/header.go":     true,
	"s3/s3api/s3err/s3-error.go":          true,
}

// TestNoForeignProductNames fails when a tracked file names another object
// store outside LICENSE, NOTICE and the retained headers above.
func TestNoForeignProductNames(t *testing.T) {
	out, err := exec.Command("git", "-C", "../..", "grep", "-lIE", "minio|MinIO|Minio|MINIO|min\\.io").Output()
	if err != nil && len(out) == 0 {
		t.Skip("not a git checkout")
	}
	for _, f := range bytes.Fields(out) {
		name := string(f)
		base := name[strings.LastIndex(name, "/")+1:]
		if base == "LICENSE" || base == "NOTICE" || retainedHeaders[name] || name == "s3/command/names_test.go" {
			continue
		}
		t.Errorf("%s names another product", name)
	}
}
