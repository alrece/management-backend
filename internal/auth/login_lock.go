package auth

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"management-backend/internal/config"

	"github.com/redis/go-redis/v9"
)

const loginFailPrefix = "auth:login:fail:"

// LoginLock 登录失败锁定
type LoginLock struct {
	rdb *redis.Client
}

func NewLoginLock(rdb *redis.Client) *LoginLock {
	return &LoginLock{rdb: rdb}
}

// RecordFail 记录登录失败
func (l *LoginLock) RecordFail(ctx context.Context, username string) error {
	key := loginFailPrefix + username
	count, err := l.rdb.Incr(ctx, key).Result()
	if err != nil {
		return err
	}
	if count == 1 {
		l.rdb.Expire(ctx, key, 30*time.Minute)
	}
	return nil
}

// IsLocked 检查账号是否被锁定
func (l *LoginLock) IsLocked(ctx context.Context, username string) (bool, error) {
	count, err := l.GetFailCount(ctx, username)
	if err != nil {
		return false, err
	}
	return count >= int64(config.C.LoginSecurity.MaxFailCount), nil
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
