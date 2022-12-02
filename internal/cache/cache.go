package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/allegro/bigcache/v3"
)

type Cache struct {
	core *bigcache.BigCache
}

func New(eviction time.Duration) *Cache {
	core, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(eviction))
	return &Cache{core}
}

// raw

func (c *Cache) Get(key string) ([]byte, bool) {
	entry, err := c.core.Get(key)
	return entry, err == nil
}

func (c *Cache) Set(key string, entry []byte) error {
	return c.core.Set(key, entry)
}

// json

func (c *Cache) GetJSON(key string, v any) bool {
	entry, err := c.core.Get(key)
	if err != nil {
		return false
	}
	err = json.Unmarshal(entry, &v)
	return err == nil
}

func (c *Cache) SetJSON(key string, v any) error {
	entry, err := json.Marshal(&v)
	if err != nil {
		return err
	}
	return c.core.Set(key, entry)
}
