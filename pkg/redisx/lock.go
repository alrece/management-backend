package redisx

import (
	"context"
	"errors"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

var unlockScript = goredis.NewScript(`
	if redis.call("get", KEYS[1]) == ARGV[1] then
		return redis.call("del", KEYS[1])
	else
		return 0
	end
`)

// DistributedLock 基于 Redis SetNX 的分布式锁
type DistributedLock struct {
	client   *goredis.Client
	key      string
	value    string
	ttl      time.Duration
	ctx      context.Context
	cancel   context.CancelFunc
	stopChan chan struct{}
}

// LockOpts 锁选项
type LockOpts struct {
	TTL        time.Duration
	RetryCount int
	RetryDelay time.Duration
}

// DefaultLockOpts 默认锁选项
func DefaultLockOpts() LockOpts {
	return LockOpts{
		TTL:        10 * time.Second,
		RetryCount: 50,
		RetryDelay: 200 * time.Millisecond,
	}
}

// AcquireLock 尝试获取分布式锁
func AcquireLock(ctx context.Context, client *goredis.Client, key string, opts LockOpts) (*DistributedLock, error) {
	if opts.TTL == 0 {
		opts.TTL = 10 * time.Second
	}
	if opts.RetryCount == 0 {
		opts.RetryCount = 50
	}
	if opts.RetryDelay == 0 {
		opts.RetryDelay = 200 * time.Millisecond
	}

	value := fmt.Sprintf("lock:%d", time.Now().UnixNano())
	fullKey := "lock:" + key

	for i := 0; i < opts.RetryCount; i++ {
		ok, err := client.SetNX(ctx, fullKey, value, opts.TTL).Result()
		if err != nil {
			return nil, fmt.Errorf("redis SetNX failed: %w", err)
		}
		if ok {
			lockCtx, cancel := context.WithCancel(ctx)
			l := &DistributedLock{
				client:   client,
				key:      fullKey,
				value:    value,
				ttl:      opts.TTL,
				ctx:      lockCtx,
				cancel:   cancel,
				stopChan: make(chan struct{}),
			}
			go l.renew()
			return l, nil
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(opts.RetryDelay):
		}
	}

	return nil, errors.New("lock acquire timeout")
}

// renew 后台续期 goroutine
func (l *DistributedLock) renew() {
	ticker := time.NewTicker(l.ttl / 3)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			l.client.Expire(l.ctx, l.key, l.ttl)
		case <-l.ctx.Done():
			return
		case <-l.stopChan:
			return
		}
	}
}

// Unlock 释放锁（仅持有者可释放）
func (l *DistributedLock) Unlock(ctx context.Context) error {
	l.cancel()
	close(l.stopChan)

	result, err := unlockScript.Run(ctx, l.client, []string{l.key}, l.value).Int()
	if err != nil {
		return fmt.Errorf("unlock script failed: %w", err)
	}
	if result == 0 {
		return errors.New("lock not held by this owner")
	}
	return nil
}
