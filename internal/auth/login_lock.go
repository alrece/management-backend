package auth

import (
	"management-backend/internal/config"
	pkgauth "management-backend/pkg/authx"

	"github.com/redis/go-redis/v9"
)

// LoginLock 登录锁定（类型别名）
type LoginLock = pkgauth.LoginLock

// NewLoginLock 创建登录锁定器
func NewLoginLock(rdb *redis.Client) *LoginLock {
	cfg := pkgauth.LoginLockConfig{
		MaxFailCount: config.C.LoginSecurity.MaxFailCount,
	}
	return pkgauth.NewLoginLock(rdb, cfg)
}
