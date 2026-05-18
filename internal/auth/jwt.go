package auth

import (
	"crypto/rand"
	"fmt"
	"time"

	"management-backend/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

// TokenPair 双 Token 对
type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	AccessExpire int64  `json:"accessExpire"`
}

// Claims JWT 载荷（含 JTI 用于黑名单）
type Claims struct {
	UserID   int64  `json:"userId"`
	Username string `json:"username"`
	TenantID int64  `json:"tenantId"`
	JTI      string `json:"jti"`
	jwt.RegisteredClaims
}

// GenerateTokenPair 生成双 Token
func GenerateTokenPair(userID int64, username string, tenantID int64) (*TokenPair, error) {
	now := time.Now()
	accessExp := time.Duration(config.C.JWT.AccessExpire) * time.Second
	refreshExp := time.Duration(config.C.JWT.RefreshExpire) * time.Second

	accessJTI, _ := generateJTI()
	refreshJTI, _ := generateJTI()

	accessClaims := Claims{
		UserID:   userID,
		Username: username,
		TenantID: tenantID,
		JTI:      accessJTI,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.C.JWT.Issuer,
			ExpiresAt: jwt.NewNumericDate(now.Add(accessExp)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        accessJTI,
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessStr, err := accessToken.SignedString([]byte(config.C.JWT.Secret))
	if err != nil {
		return nil, fmt.Errorf("生成 Access Token 失败: %w", err)
	}

	refreshClaims := Claims{
		UserID:   userID,
		Username: username,
		TenantID: tenantID,
		JTI:      refreshJTI,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.C.JWT.Issuer,
			ExpiresAt: jwt.NewNumericDate(now.Add(refreshExp)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        refreshJTI,
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshStr, err := refreshToken.SignedString([]byte(config.C.JWT.Secret))
	if err != nil {
		return nil, fmt.Errorf("生成 Refresh Token 失败: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessStr,
		RefreshToken: refreshStr,
		AccessExpire: now.Add(accessExp).Unix(),
	}, nil
}

// ParseAccessToken 解析 Access Token
func ParseAccessToken(tokenStr string) (*Claims, error) {
	return parseToken(tokenStr)
}

// ParseRefreshToken 解析 Refresh Token
func ParseRefreshToken(tokenStr string) (*Claims, error) {
	return parseToken(tokenStr)
}

func parseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.C.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("无效的 Token")
}

func generateJTI() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}
