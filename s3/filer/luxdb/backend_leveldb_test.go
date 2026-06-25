//go:build leveldb

package luxdb

import "testing"

// With -tags=leveldb the "leveldb" backend must open a working store.
func TestLevelDBBackendRoundTrip(t *testing.T) {
	runBackendRoundTrip(t, "leveldb")
}
