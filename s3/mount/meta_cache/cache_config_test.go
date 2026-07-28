package meta_cache

import "testing"

// TestCacheConfigDoesNotAnswerUnknownKeys pins the behaviour that broke the
// mount: GetString used to return c.dir for EVERY key that did not end in
// "backend". It is a util.Configuration, so any store it initialises may ask it
// for options it has never heard of — and it answered all of them with a
// filesystem path.
//
// That is how a store gaining a "syncWrites" option received "/tmp/<id>/meta"
// as its value and refused to start, taking the whole FUSE mount with it. The
// store was right to reject the value; the config was wrong to invent it.
//
// A config that cannot say "I don't know" will keep breaking the next option
// somebody adds, so the guarantee under test is the silence, not any one key.
func TestCacheConfigDoesNotAnswerUnknownKeys(t *testing.T) {
	c := cacheConfig{dir: "/tmp/abc/meta", backend: "zapdb"}

	if got := c.GetString("luxdb.backend"); got != "zapdb" {
		t.Errorf("backend = %q, want %q", got, "zapdb")
	}
	if got := c.GetString("luxdb.dir"); got != "/tmp/abc/meta" {
		t.Errorf("dir = %q, want %q", got, "/tmp/abc/meta")
	}

	// Anything else must come back empty rather than as a path.
	for _, key := range []string{
		"luxdb.syncWrites",
		"syncWrites",
		"luxdb.replication",
		"luxdb.ttl",
		"some.future.option",
	} {
		if got := c.GetString(key); got != "" {
			t.Errorf("GetString(%q) = %q, want \"\" — an unknown key must not be answered with a path", key, got)
		}
	}
}
