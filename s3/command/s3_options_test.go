package command

import (
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

// S3Options is a struct of flag POINTERS shared by four commands — s3, server,
// filer, mini — that each register their own flag set into it. The code they
// share reaches those pointers through the s3opt receiver and dereferences them
// with no nil check, so a command that forgets one does not degrade. It
// segfaults the moment S3 starts:
//
//	panic: runtime error: invalid memory address or nil pointer dereference
//	  command.(*S3Options).startS3Server(...)  s3/command/s3.go:282
//
// That is what `zapPort` did. It was added to `s3` and `server` only, and both
// `s3 mini` and `s3 filer -s3` crashed on startup — S3 never bound its port.
//
// The struct's doc comment already said "when adding a new field, update all
// four flag registration sites". A comment cannot fail a build. This is that
// sentence, enforced.
//
// The field set is READ FROM THE SOURCE rather than listed here, because a
// hand-maintained list is the same forget-to-update bug one level up. Only
// fields the shared code actually dereferences are required; a field used
// solely as s3StandaloneOptions.x (debug, debugPort) belongs to that one
// command and is correctly absent from the others.
func TestS3OptionsRegisteredEverywhere(t *testing.T) {
	shared := fieldsDereferencedViaReceiver(t)
	if len(shared) == 0 {
		t.Fatal("scanned the package and found no s3opt.<field> dereferences — the scan is broken, not the code")
	}

	typ := reflect.TypeOf(S3Options{})
	for _, cmd := range []struct {
		name string
		opts S3Options
	}{
		{"s3", s3StandaloneOptions},
		{"server", s3Options},
		{"filer", filerS3Options},
		{"mini", miniS3Options},
	} {
		v := reflect.ValueOf(cmd.opts)
		for name := range shared {
			field, ok := typ.FieldByName(name)
			if !ok {
				continue
			}
			if runtimeAssigned[name] {
				continue
			}
			if v.FieldByIndex(field.Index).IsNil() {
				t.Errorf("%s: S3Options.%s is nil, and the shared s3opt path dereferences it — starting S3 from `s3 %s` will panic. Register it in the %s command.",
					cmd.name, name, cmd.name, cmd.name)
			}
		}
	}
}

// runtimeAssigned are fields no flag registration sets, because they are
// assigned once from an already-resolved value: the filer address is known only
// after the filer is up, bindIp and dataCenter fall back to their non-s3
// siblings, the sockets are derived from the filer's, and shutdownCtx is an
// orchestration handle rather than a flag. Each is set on the path that uses it.
var runtimeAssigned = map[string]bool{
	"filer":            true,
	"bindIp":           true,
	"dataCenter":       true,
	"metricsHttpIp":    true,
	"localFilerSocket": true,
	"localSocket":      true,
	"shutdownCtx":      true,
}

var receiverDeref = regexp.MustCompile(`\bs3opt\.([a-zA-Z][a-zA-Z0-9]*)\b`)

// fieldsDereferencedViaReceiver returns the S3Options fields that code shared by
// all four commands reads through the s3opt receiver — the exact set every
// command must therefore initialize.
func fieldsDereferencedViaReceiver(t *testing.T) map[string]bool {
	t.Helper()
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("read package dir: %v", err)
	}
	fields := map[string]bool{}
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		src, err := os.ReadFile(filepath.Join(".", name))
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		for _, m := range receiverDeref.FindAllStringSubmatch(string(src), -1) {
			fields[m[1]] = true
		}
	}
	return fields
}
