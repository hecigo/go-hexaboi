package middleware

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"hecigo.com/go-hexaboi/infra/core"
)

type Cache struct {
	CacheExpiration string
	CacheControl    bool
	CacheKeyPrefix  string
	ETag            bool
}

// Enable cache for app
func (_cache *Cache) Enable(app *fiber.App) error {

	// If the cache expiration is ZERO, that means disabled cache
	_cache.CacheExpiration = core.Getenv("CACHE_EXPIRATION", "1m")
	if strings.HasPrefix(_cache.CacheExpiration, "0") {
		log.Println("Middleware/Cache: disabled by ZERO expiration")
		return nil
	}

	// Parse cache expiration
	cacheExp, cacheExpErr := time.ParseDuration(_cache.CacheExpiration)
	if cacheExpErr != nil {
		log.Println(fmt.Errorf("Middleware/Cache: invalid cache expiration. "+
			"Please, check the environment variable CACHE_EXPIRATION.\r\n%w", cacheExpErr))
		log.Fatal("Can not start application, cause of Middleware/Cache has an error.")
		return cacheExpErr
	}

	_cache.CacheControl = core.GetBoolEnv("CACHE_CONTROL", false)
	_cache.CacheKeyPrefix = core.Getenv("CACHE_KEY_PREFIX", "hpi")

	// Initialize cache with config
	app.Use(cache.New(cache.Config{
		StoreResponseHeaders: true,
		CacheControl:         _cache.CacheControl,
		KeyGenerator: func(c *fiber.Ctx) string {
			uri := c.Request().URI()

			// Ensure the path end with a slash
			path := string(uri.Path())
			if !strings.HasSuffix(path, "/") {
				path += "/"
			}

			queryString := string(uri.QueryString())
			if queryString == "" {
				return path
			}
			return _cache.CacheKeyPrefix + path + "?" + string(uri.QueryString())
		},
		Next: func(c *fiber.Ctx) bool {
			// Only allow cache with GET
			return c.GetRespHeader("X-Hpi-Cache", "") == "no-cache" || c.Method() != "GET" || c.Query("clearcache") == "true"
		},
		ExpirationGenerator: func(c *fiber.Ctx, cfg *cache.Config) time.Duration {
			cacheTime := c.GetRespHeader("cache-time", "")

			// No special, return default cache expiration
			if cacheTime == "" {
				return cacheExp
			}

			newCacheExp, newCacheExpErr := time.ParseDuration(cacheTime)

			// On error, return default cache expiration
			if newCacheExpErr != nil {
				log.Println(fmt.Errorf("Middleware/Cache: invalid cache expiration. "+
					"Please, check the response header `cache-time`.\r\n%w", newCacheExpErr))
				return cacheExp
			}

			// Otherwise, return the speical cache expiration
			return newCacheExp
		},
	}))

	_cache.ETag = core.GetBoolEnv("CACHE_ETAG", true)
	if _cache.ETag {
		app.Use(etag.New())
	}

	_cache.Print()

	return nil
}

func (_cache Cache) Print() {
	fmt.Println("\r\n┌─────── Middleware/Cache ─────────")
	fmt.Printf("| CACHE_KEY_PREFIX: %s\r\n", _cache.CacheKeyPrefix)
	fmt.Printf("| CACHE_EXPIRATION: %s\r\n", _cache.CacheExpiration)
	fmt.Printf("| CACHE_CONTROL: %s\r\n", strconv.FormatBool(_cache.CacheControl))
	fmt.Printf("| CACHE_ETAG: %s\r\n", strconv.FormatBool(_cache.ETag))
	fmt.Println("| CACHE_BY_PASS: /?clearcache=true")
	fmt.Println("└──────────────────────────────────")
}
