//go:build pebbledb

package luxdb

import "testing"

// With -tags=pebbledb the "pebble" alias must open a working store.
func TestPebbleBackendRoundTrip(t *testing.T) {
	runBackendRoundTrip(t, "pebble")
}
