package auth

import (
	"management-backend/internal/config"
	pkgauth "management-backend/pkg/authx"
)

// Claims JWT 载荷（类型别名，与 pkg/authx 保持一致）
type Claims = pkgauth.Claims

// TokenPair 双 Token 对（类型别名）
type TokenPair = pkgauth.TokenPair

// GenerateTokenPair 生成双 Token（从全局配置读取 JWT 参数）
func GenerateTokenPair(userID int64, username string, tenantID int64) (*TokenPair, error) {
	cfg := pkgauth.JWTConfig{
		Secret:        config.C.JWT.Secret,
		AccessExpire:  config.C.JWT.AccessExpire,
		RefreshExpire: config.C.JWT.RefreshExpire,
		Issuer:        config.C.JWT.Issuer,
	}
	return pkgauth.GenerateTokenPair(userID, username, tenantID, cfg)
}

// ParseAccessToken 解析 Access Token
func ParseAccessToken(tokenStr string) (*Claims, error) {
	return pkgauth.ParseToken(tokenStr, config.C.JWT.Secret)
}

// ParseRefreshToken 解析 Refresh Token
func ParseRefreshToken(tokenStr string) (*Claims, error) {
	return pkgauth.ParseToken(tokenStr, config.C.JWT.Secret)
}
