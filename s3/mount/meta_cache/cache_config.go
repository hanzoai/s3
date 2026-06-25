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
	if strings.HasSuffix(key, "backend") {
		return c.backend
	}
	return c.dir
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
