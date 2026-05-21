package authz

import (
	pkgauthzx "management-backend/pkg/authzx"

	"github.com/redis/go-redis/v9"
)

// PolicySync Redis Pub/Sub 策略同步器（类型别名）
type PolicySync = pkgauthzx.PolicySync

// NewPolicySync 创建同步器
func NewPolicySync(rdb *redis.Client, manager *EnforcerManager) *PolicySync {
	return pkgauthzx.NewPolicySync(rdb, manager)
}
