package redisx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

var ErrCacheMiss = errors.New("cache: key not found")

// Cache 提供带 TTL 的类型化缓存操作
type Cache struct {
	client *goredis.Client
	prefix string
}

// NewCache 创建缓存实例
func NewCache(client *goredis.Client, prefix string) *Cache {
	return &Cache{client: client, prefix: prefix}
}

// Get 从缓存获取值
func Get[T any](ctx context.Context, c *Cache, key string) (T, error) {
	var zero T
	val, err := c.client.Get(ctx, c.prefix+key).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return zero, ErrCacheMiss
		}
		return zero, fmt.Errorf("cache get failed: %w", err)
	}

	var result T
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return zero, fmt.Errorf("cache unmarshal failed: %w", err)
	}
	return result, nil
}

// Set 写入缓存值
func Set[T any](ctx context.Context, c *Cache, key string, value T, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache marshal failed: %w", err)
	}
	return c.client.Set(ctx, c.prefix+key, data, ttl).Err()
}

// Delete 删除缓存键
func Delete(ctx context.Context, c *Cache, key string) error {
	return c.client.Del(ctx, c.prefix+key).Err()
}
