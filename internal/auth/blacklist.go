package auth

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	jtiPrefix     = "auth:jti:"
	userRevPrefix = "auth:user:revoked:"
)

// Blacklist Token 黑名单（依赖 *redis.Client）
type Blacklist struct {
	rdb *redis.Client
}

func NewBlacklist(rdb *redis.Client) *Blacklist {
	return &Blacklist{rdb: rdb}
}

// AddJTI 将 JTI 加入黑名单
func (b *Blacklist) AddJTI(ctx context.Context, jti string, ttl time.Duration) error {
	return b.rdb.Set(ctx, jtiPrefix+jti, "1", ttl).Err()
}

// IsBlacklisted 检查 JTI 是否在黑名单
func (b *Blacklist) IsBlacklisted(ctx context.Context, jti string) (bool, error) {
	val, err := b.rdb.Exists(ctx, jtiPrefix+jti).Result()
	if err != nil {
		return false, err
	}
	return val > 0, nil
}

// RevokeByUser 吊销用户所有 Token（记录吊销时间戳）
func (b *Blacklist) RevokeByUser(ctx context.Context, userID int64) error {
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	return b.rdb.Set(ctx, fmt.Sprintf("%s%d", userRevPrefix, userID), ts, 7*24*time.Hour).Err()
}

// IsUserRevoked 检查用户是否在签发时间之后被吊销
func (b *Blacklist) IsUserRevoked(ctx context.Context, userID int64, issuedAt time.Time) (bool, error) {
	val, err := b.rdb.Get(ctx, fmt.Sprintf("%s%d", userRevPrefix, userID)).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	ts, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return false, nil
	}
	return issuedAt.Unix() <= ts, nil
}
