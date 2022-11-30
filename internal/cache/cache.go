package cache

import (
	"context"
	"time"

	"github.com/allegro/bigcache/v3"
)

const (
	eviction = 10 * time.Minute
)

var cache *bigcache.BigCache

func init() {
	cache, _ = bigcache.New(context.Background(), bigcache.DefaultConfig(eviction))
}

func Get(key string) ([]byte, bool) {
	entry, err := cache.Get(key)
	return entry, err == nil
}

func Set(key string, entry []byte) {
	cache.Set(key, entry)
}
