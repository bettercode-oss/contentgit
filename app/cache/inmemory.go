package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		cache: cache.New(5*time.Minute, 10*time.Minute),
	}
}

type InMemoryCache struct {
	cache *cache.Cache
}

func (c InMemoryCache) Set(key string, value any, expiration time.Duration) {
	c.cache.Set(key, value, expiration)
}

func (c InMemoryCache) Get(key string) (any, bool) {
	return c.cache.Get(key)
}
