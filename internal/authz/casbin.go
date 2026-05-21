package authz

import (
	pkgauthzx "management-backend/pkg/authzx"

	"github.com/redis/go-redis/v9"
)

// EnforcerManager 多租户 Casbin 管理器（类型别名）
type EnforcerManager = pkgauthzx.EnforcerManager

// NewEnforcerManager 创建管理器
func NewEnforcerManager(rdb *redis.Client) *EnforcerManager {
	return pkgauthzx.NewEnforcerManager(rdb)
}
