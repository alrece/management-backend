package redisx

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

const defaultIdempotentTTL = 24 * time.Hour

// IdempotentExecutor 幂等执行器，防止重复操作
type IdempotentExecutor struct {
	client *goredis.Client
	prefix string
	ttl    time.Duration
}

// NewIdempotentExecutor 创建幂等执行器
func NewIdempotentExecutor(client *goredis.Client, prefix string) *IdempotentExecutor {
	return &IdempotentExecutor{
		client: client,
		prefix: prefix,
		ttl:    defaultIdempotentTTL,
	}
}

// WithTTL 设置幂等键 TTL
func (e *IdempotentExecutor) WithTTL(ttl time.Duration) *IdempotentExecutor {
	e.ttl = ttl
	return e
}

// Execute 幂等执行函数，相同 key 在 TTL 内只执行一次
func Execute[T any](ctx context.Context, e *IdempotentExecutor, key string, fn func() (T, error)) (T, error) {
	var zero T
	fullKey := e.prefix + key

	// 检查是否已有缓存结果
	cached, err := e.client.Get(ctx, fullKey).Result()
	if err == nil {
		var result T
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			return result, nil
		}
	}

	// 尝试获取执行权
	acquired, err := e.client.SetNX(ctx, fullKey, "", e.ttl).Result()
	if err != nil {
		return zero, fmt.Errorf("idempotent SetNX failed: %w", err)
	}
	if !acquired {
		// 其他实例正在执行，短暂等待后重试读取
		time.Sleep(100 * time.Millisecond)
		cached, err = e.client.Get(ctx, fullKey).Result()
		if err != nil {
			return zero, fmt.Errorf("idempotent concurrent read failed: %w", err)
		}
		var result T
		if err := json.Unmarshal([]byte(cached), &result); err != nil {
			return zero, fmt.Errorf("idempotent unmarshal failed: %w", err)
		}
		return result, nil
	}

	// 执行函数
	result, fnErr := fn()
	if fnErr != nil {
		// 执行失败，删除幂等键允许重试
		e.client.Del(ctx, fullKey)
		return zero, fnErr
	}

	// 存储结果
	data, err := json.Marshal(result)
	if err != nil {
		return zero, fmt.Errorf("idempotent marshal failed: %w", err)
	}
	if err := e.client.Set(ctx, fullKey, data, e.ttl).Err(); err != nil {
		return zero, fmt.Errorf("idempotent store failed: %w", err)
	}

	return result, nil
}

// IsExecuted 检查幂等键是否已执行
func (e *IdempotentExecutor) IsExecuted(ctx context.Context, key string) (bool, error) {
	exists, err := e.client.Exists(ctx, e.prefix+key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}
