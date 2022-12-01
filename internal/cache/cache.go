package cache

import (
	"context"
	"encoding/json"
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

// raw

func Get(key string) ([]byte, bool) {
	entry, err := cache.Get(key)
	return entry, err == nil
}

func Set(key string, entry []byte) error {
	return cache.Set(key, entry)
}

// json

func GetJSON(key string, v any) bool {
	entry, err := cache.Get(key)
	if err != nil {
		return false
	}
	err = json.Unmarshal(entry, &v)
	return err == nil
}

func SetJSON(key string, v any) error {
	entry, err := json.Marshal(&v)
	if err != nil {
		return err
	}
	return cache.Set(key, entry)
}
