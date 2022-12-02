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

func (c *Cache) Get(key string) ([]byte, error) {
	entry, err := c.core.Get(key)
	return entry, err
}

func (c *Cache) Set(key string, entry []byte) error {
	return c.core.Set(key, entry)
}

// json

func (c *Cache) GetJSON(key string, v any) error {
	entry, err := c.core.Get(key)
	if err != nil {
		return err
	}
	err = json.Unmarshal(entry, &v)
	return err
}

func (c *Cache) SetJSON(key string, v any) error {
	entry, err := json.Marshal(&v)
	if err != nil {
		return err
	}
	return c.core.Set(key, entry)
}

// generic

type GenericCache[T any] struct {
	cache *Cache
}

func NewGeneric[T any](eviction time.Duration) *GenericCache[T] {
	return &GenericCache[T]{
		cache: New(eviction),
	}
}

func (c *GenericCache[T]) Get(key string) (T, error) {
	var v T
	err := c.cache.GetJSON(key, &v)
	return v, err
}

func (c *GenericCache[T]) Set(key string, v T) error {
	return c.cache.SetJSON(key, v)
}
