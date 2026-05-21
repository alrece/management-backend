package middleware

import (
	"net/http"
	"strconv"

	"management-backend/internal/authz"

	"github.com/gin-gonic/gin"
)

var enforcerMgr *authz.EnforcerManager

// InitPermission 初始化权限组件
func InitPermission(mgr *authz.EnforcerManager) {
	enforcerMgr = mgr
}

// RequirePermission 声明式权限中间件
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if enforcerMgr == nil {
			c.Next()
			return
		}

		userID := GetUserID(c)
		tenantID := GetTenantID(c)

		ok, err := enforcerMgr.CheckPermission(
			c.Request.Context(), userID, tenantID, permission,
		)
		if err != nil || !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "权限不足",
			})
			return
		}
		c.Next()
	}
}

// GetEnforcerMgr 获取 EnforcerManager（供 service 层使用）
func GetEnforcerMgr() *authz.EnforcerManager {
	return enforcerMgr
}

// GetUserIDStr 获取当前用户 ID 字符串
func GetUserIDStr(c *gin.Context) string {
	return strconv.FormatInt(GetUserID(c), 10)
}

// GetTenantIDStr 获取当前租户 ID 字符串
func GetTenantIDStr(c *gin.Context) string {
	return strconv.FormatInt(GetTenantID(c), 10)
}
