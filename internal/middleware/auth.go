package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"management-backend/internal/auth"
	"management-backend/internal/config"
	"management-backend/pkg/errcode"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var blacklist *auth.Blacklist

// InitAuth 初始化认证组件（main.go 调用）
func InitAuth(rdb *redis.Client) {
	blacklist = auth.NewBlacklist(rdb)
}

// Auth JWT 认证中间件（检查黑名单 + 用户级吊销）
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未授权，请先登录"})
			return
		}
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		claims, err := auth.ParseAccessToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Token 无效或已过期"})
			return
		}

		ctx := c.Request.Context()
		if blacklist != nil {
			revoked, _ := blacklist.IsBlacklisted(ctx, claims.JTI)
			if revoked {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Token 已被吊销"})
				return
			}
			userRevoked, _ := blacklist.IsUserRevoked(ctx, claims.UserID, claims.IssuedAt.Time)
			if userRevoked {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "账号凭证已失效，请重新登录"})
				return
			}
		}

		c.Set("userId", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("tenantId", claims.TenantID)
		c.Next()
	}
}

// RevokeToken 吊销 Token（登出时调用）
func RevokeToken(ctx context.Context, jti string) error {
	if blacklist == nil {
		return errcode.Err(errcode.TokenInvalid)
	}
	ttl := time.Duration(config.C.JWT.AccessExpire) * time.Second
	return blacklist.AddJTI(ctx, jti, ttl)
}

// RevokeUserTokens 吊销用户所有 Token（改密/禁用时调用）
func RevokeUserTokens(ctx context.Context, userID int64) error {
	if blacklist == nil {
		return nil
	}
	return blacklist.RevokeByUser(ctx, userID)
}

// GetUserID 从上下文获取当前用户 ID
func GetUserID(c *gin.Context) int64 {
	if v, ok := c.Get("userId"); ok {
		return v.(int64)
	}
	return 0
}

// GetTenantID 从上下文获取当前租户 ID
func GetTenantID(c *gin.Context) int64 {
	if v, ok := c.Get("tenantId"); ok {
		return v.(int64)
	}
	return 0
}
