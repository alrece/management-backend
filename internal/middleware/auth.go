package middleware

import (
	"strings"
	"time"

	"management-backend/internal/config"
	"management-backend/pkg/errcode"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT 载荷
type Claims struct {
	UserID   int64  `json:"userId"`
	Username string `json:"username"`
	TenantID int64  `json:"tenantId"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT Token
func GenerateToken(userID int64, username string, tenantID int64) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		TenantID: tenantID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.C.JWT.Issuer,
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(config.C.JWT.AccessExpire) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.C.JWT.Secret))
}

// ParseToken 解析 JWT Token
func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.C.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errcode.Err(errcode.TokenInvalid)
}

// Auth JWT 认证中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.AbortWithStatusJSON(401, gin.H{"code": 401, "msg": "未授权，请先登录"})
			return
		}
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
		claims, err := ParseToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"code": 401, "msg": "Token 无效或已过期"})
			return
		}
		c.Set("userId", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("tenantId", claims.TenantID)
		c.Next()
	}
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
