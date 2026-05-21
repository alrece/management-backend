package authx

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const loginFailPrefix = "auth:login:fail:"

// LoginLockConfig 登录锁定配置
type LoginLockConfig struct {
	MaxFailCount int           // 最大失败次数
	LockDuration time.Duration // 锁定时长
}

// LoginLock 登录失败锁定（依赖 *redis.Client）
type LoginLock struct {
	rdb *redis.Client
	cfg LoginLockConfig
}

// NewLoginLock 创建登录锁定器
func NewLoginLock(rdb *redis.Client, cfg LoginLockConfig) *LoginLock {
	return &LoginLock{rdb: rdb, cfg: cfg}
}

// RecordFail 记录登录失败
func (l *LoginLock) RecordFail(ctx context.Context, username string) error {
	key := loginFailPrefix + username
	count, err := l.rdb.Incr(ctx, key).Result()
	if err != nil {
		return err
	}
	if count == 1 {
		l.rdb.Expire(ctx, key, l.lockTTL())
	}
	return nil
}

// IsLocked 检查账号是否被锁定
func (l *LoginLock) IsLocked(ctx context.Context, username string) (bool, error) {
	count, err := l.GetFailCount(ctx, username)
	if err != nil {
		return false, err
	}
	return count >= int64(l.cfg.MaxFailCount), nil
}

// ResetFail 重置失败计数
func (l *LoginLock) ResetFail(ctx context.Context, username string) error {
	return l.rdb.Del(ctx, loginFailPrefix+username).Err()
}

// GetFailCount 获取失败次数
func (l *LoginLock) GetFailCount(ctx context.Context, username string) (int64, error) {
	val, err := l.rdb.Get(ctx, loginFailPrefix+username).Result()
	if err != nil {
		return 0, nil
	}
	count, _ := strconv.ParseInt(val, 10, 64)
	return count, nil
}

// GetLockRemaining 获取锁定剩余时间
func (l *LoginLock) GetLockRemaining(ctx context.Context, username string) time.Duration {
	key := fmt.Sprintf("%s%s", loginFailPrefix, username)
	ttl, err := l.rdb.TTL(ctx, key).Result()
	if err != nil || ttl < 0 {
		return 0
	}
	return ttl
}

func (l *LoginLock) lockTTL() time.Duration {
	if l.cfg.LockDuration > 0 {
		return l.cfg.LockDuration
	}
	return 30 * time.Minute
}
