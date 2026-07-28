package meta_cache

import (
	"strings"

	"github.com/hanzoai/s3/s3/util"
)

var (
	_ = util.Configuration(&cacheConfig{})
)

// implementing util.Configuration
type cacheConfig struct {
	dir     string
	backend string
}

func (c cacheConfig) GetString(key string) string {
	// Answer only what this cache actually configures. The previous form
	// returned c.dir for EVERY key that did not end in "backend", so any new
	// store option silently received a filesystem path as its value — a store
	// reading a "syncWrites" key got "/tmp/<id>/meta" and refused to start.
	// An unknown key has no value here; say so.
	switch {
	case strings.HasSuffix(key, "backend"):
		return c.backend
	case strings.HasSuffix(key, "dir"):
		return c.dir
	default:
		return ""
	}
}

func (c cacheConfig) GetBool(key string) bool {
	panic("implement me")
}

func (c cacheConfig) GetInt(key string) int {
	panic("implement me")
}

func (c cacheConfig) GetStringSlice(key string) []string {
	panic("implement me")
}

func (c cacheConfig) SetDefault(key string, value interface{}) {
	panic("implement me")
}
