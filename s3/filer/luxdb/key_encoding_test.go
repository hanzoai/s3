package luxdb

import (
	"bytes"
	"testing"

	"github.com/hanzoai/s3/s3/util"
)

// TestKeyEncodingByteIdentical pins the exact key bytes so the luxdb layout
// stays byte-identical to the historical leveldb store (dir + 0x00 + name).
// Any drift here silently breaks reading an existing on-disk store.
func TestKeyEncodingByteIdentical(t *testing.T) {
	got := genKey("/a/b", "c.txt")
	want := append(append([]byte("/a/b"), DIR_FILE_SEPARATOR), []byte("c.txt")...)
	if !bytes.Equal(got, want) {
		t.Fatalf("genKey = %q, want %q", got, want)
	}

	// Directory prefix without a start file: dir + 0x00.
	gotPrefix := genDirectoryKeyPrefix(util.FullPath("/a/b"), "")
	wantPrefix := append([]byte("/a/b"), DIR_FILE_SEPARATOR)
	if !bytes.Equal(gotPrefix, wantPrefix) {
		t.Fatalf("genDirectoryKeyPrefix empty = %q, want %q", gotPrefix, wantPrefix)
	}

	// Directory prefix with a start file: dir + 0x00 + start.
	gotStart := genDirectoryKeyPrefix(util.FullPath("/a/b"), "m")
	wantStart := append(append([]byte("/a/b"), DIR_FILE_SEPARATOR), []byte("m")...)
	if !bytes.Equal(gotStart, wantStart) {
		t.Fatalf("genDirectoryKeyPrefix start = %q, want %q", gotStart, wantStart)
	}

	// getNameFromKey is the inverse of genKey on the name component.
	if name := getNameFromKey(genKey("/a/b", "c.txt")); name != "c.txt" {
		t.Fatalf("getNameFromKey = %q, want c.txt", name)
	}
	// Root-level file: empty dir + 0x00 + name.
	if name := getNameFromKey(genKey("", "root.txt")); name != "root.txt" {
		t.Fatalf("getNameFromKey root = %q, want root.txt", name)
	}
}
