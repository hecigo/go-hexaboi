package redis

import (
	"time"

	"github.com/go-redis/cache/v8"
)

var defaultInstance *cache.Cache

func InitCache() {
	defaultInstance = cache.New(&cache.Options{
		Redis: ClientByName("cache"),

		// use local in-process storage to cache the small subset of popular keys
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})
}

func Cache() *cache.Cache {
	return defaultInstance
}
