package redis

import (
	"context"
	"fmt"
	"strings"

	"github.com/goccy/go-json"
	"hoangphuc.tech/go-hexaboi/infra/core"

	"github.com/go-redis/redis/v8"
)

var (
	ctx     = context.Background()
	appName = core.Getenv("APP_NAME", "GoHexaboi")
)

func getRedisKey(key string) string {
	keyPrefix := core.Getenv("REDIS_KEY_PREFIX", appName+".")
	return keyPrefix + key
}

func removePrefix(key string) string {
	keyPrefix := core.Getenv("REDIS_KEY_PREFIX", appName+".")
	return strings.TrimPrefix(key, keyPrefix)
}

// Set value by key
func Set(key string, value interface{}) error {
	err := DB().Set(ctx, getRedisKey(key), value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

// Get value by key with default format
func Get[T any](key string) (*T, error) {
	return GetSpecificKey[T](getRedisKey(key))
}

// Get value by a specific key
func GetSpecificKey[T any](key string) (*T, error) {
	val, err := DB().Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var t T
	if err = json.Unmarshal([]byte(val), &t); err != nil {
		return nil, err
	}

	return &t, err
}

func IncrBy(key string, value int64) (int64, error) {
	val, err := DB().IncrBy(ctx, getRedisKey(key), value).Result()
	if err != nil {
		return 0, err
	}
	return val, nil
}

func MultiIncrBy(kv map[string]int64) (map[string]int64, error) {
	cmds, err := DB().TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		for k, v := range kv {
			pipe.IncrBy(ctx, getRedisKey(k), v)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	m := make(map[string]int64)
	for _, cmd := range cmds {
		intCmd := cmd.(*redis.IntCmd)
		fmt.Println(intCmd)
		k := removePrefix(intCmd.Args()[1].(string))
		v := intCmd.Val()
		m[k] = v
	}

	return m, nil
}
