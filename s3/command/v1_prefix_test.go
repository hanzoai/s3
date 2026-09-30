package command

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// firstPartyAPIPath matches a path this repo serves or calls that begins with
// /api/: a quoted or template-built relative path ("/api/…", `${x}/api/…`,
// "%s/api/…"). Every first-party HTTP route is /v1/…; the S3 wire protocol
// (bucket and object paths) never had the prefix and is not matched.
var firstPartyAPIPath = regexp.MustCompile("([\"'`(]|%s|\\})/api([/\"'`]|$)")

// TestNoAPIPathPrefix fails when a first-party /api/ path is served or called
// again anywhere in the s3 or telemetry trees.
func TestNoAPIPathPrefix(t *testing.T) {
	for _, root := range []string{"..", filepath.Join("..", "..", "telemetry")} {
		err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				if d.Name() == "static_gz" || d.Name() == "node_modules" {
					return filepath.SkipDir
				}
				return nil
			}
			name := d.Name()
			switch {
			case strings.HasSuffix(name, "_templ.go"), strings.HasSuffix(name, ".min.js"),
				strings.HasSuffix(name, "_test.go"):
				return nil // generated from a scanned .templ, vendored, or a test asserting the negative
			case strings.HasSuffix(name, ".go"), strings.HasSuffix(name, ".js"),
				strings.HasSuffix(name, ".templ"), strings.HasSuffix(name, ".html"):
			default:
				return nil
			}
			data, err := os.ReadFile(p)
			if err != nil {
				return err
			}
			for i, line := range strings.Split(string(data), "\n") {
				if firstPartyAPIPath.MatchString(line) {
					t.Errorf("%s:%d: first-party /api/ path; use /v1/: %s", p, i+1, strings.TrimSpace(line))
				}
			}
			return nil
		})
		if err != nil {
			t.Fatalf("walk %s: %v", root, err)
		}
	}
}
