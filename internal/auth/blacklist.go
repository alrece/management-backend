package auth

import (
	pkgauth "management-backend/pkg/authx"

	"github.com/redis/go-redis/v9"
)

// Blacklist Token 黑名单（类型别名）
type Blacklist = pkgauth.Blacklist

// NewBlacklist 创建黑名单
func NewBlacklist(rdb *redis.Client) *Blacklist {
	return pkgauth.NewBlacklist(rdb)
}
