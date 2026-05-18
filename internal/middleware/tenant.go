package middleware

import (
	"net/http"

	"management-backend/internal/config"
	"management-backend/internal/multitenant"
	pkgmiddleware "management-backend/pkg/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var resolver *multitenant.TenantResolver

// InitTenant 初始化租户解析器（main.go 调用）
func InitTenant() {
	if config.C.Tenant.Enable {
		resolver = multitenant.NewTenantResolver()
	}
}

// GetResolver 获取租户解析器实例
func GetResolver() *multitenant.TenantResolver {
	return resolver
}

// Tenant 多租户中间件（独立数据库隔离）
func Tenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.C.Tenant.Enable || resolver == nil {
			c.Next()
			return
		}

		tenantID := GetTenantID(c)
		if tenantID <= 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 401, "msg": "租户信息缺失",
			})
			return
		}

		tenantDB, err := resolver.GetDB(c.Request.Context(), tenantID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code": 403, "msg": "租户数据库不可用",
			})
			return
		}

		// 注入租户 DB 到 request context
		ctx := pkgmiddleware.WithTenantDB(c.Request.Context(), tenantDB)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// GetTenantScope 获取租户作用域函数（已弃用，保留兼容）
// 独立数据库隔离模式下不再需要行级过滤
func GetTenantScope(c *gin.Context) func(*gorm.DB) *gorm.DB {
	return nil
}
