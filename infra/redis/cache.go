package redis

import (
	"time"

	"github.com/go-redis/cache/v8"
	"hoangphuc.tech/dora/infra/core"
)

var defaultInstance *cache.Cache

//  Initialize cache instance with Redis
func InitCache() {

	// Use local in-process storage to cache the small subset of popular keys
	// Default cache 1000 keys for 1 minute.
	size := core.GetIntEnv("REDIS_CACHE_TINYFLU_SIZE", 1000)
	duration := core.GetDurationEnv("REDIS_CACHE_TINYFLU_DURATION", time.Minute)

	defaultInstance = cache.New(&cache.Options{
		Redis:      ClientByName("cache"),
		LocalCache: cache.NewTinyLFU(size, duration),
	})
}

// Note: DO NOT use cache with major/important data.
// Because it always stores a subset of popular keys in local in-process storage,
// cache data will be cleared after each deployment.
func Cache() *cache.Cache {
	return defaultInstance
}
