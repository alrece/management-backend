package middleware

import (
	"management-backend/internal/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Tenant 多租户中间件
// 将 tenant_id 条件注入 GORM Scope，对 Repository 层透明
func Tenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.C.Tenant.Enable {
			c.Next()
			return
		}
		tenantID := GetTenantID(c)
		if tenantID <= 0 {
			c.AbortWithStatusJSON(401, gin.H{"code": 401, "msg": "租户信息缺失"})
			return
		}
		c.Set("tenantScope", func(db *gorm.DB) *gorm.DB {
			return db.Where("tenant_id = ?", tenantID)
		})
		c.Next()
	}
}

// GetTenantScope 获取租户作用域函数，供 Repository 层使用
func GetTenantScope(c *gin.Context) func(*gorm.DB) *gorm.DB {
	if fn, ok := c.Get("tenantScope"); ok {
		return fn.(func(*gorm.DB) *gorm.DB)
	}
	return nil
}
