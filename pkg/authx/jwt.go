package authx

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTConfig JWT 配置（由调用方传入，不读全局变量）
type JWTConfig struct {
	Secret        string
	AccessExpire  int64 // 秒
	RefreshExpire int64 // 秒
	Issuer        string
}

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
func GenerateTokenPair(userID int64, username string, tenantID int64, cfg JWTConfig) (*TokenPair, error) {
	now := time.Now()
	accessExp := time.Duration(cfg.AccessExpire) * time.Second
	refreshExp := time.Duration(cfg.RefreshExpire) * time.Second

	accessJTI, _ := generateJTI()
	refreshJTI, _ := generateJTI()

	accessClaims := Claims{
		UserID:   userID,
		Username: username,
		TenantID: tenantID,
		JTI:      accessJTI,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.Issuer,
			ExpiresAt: jwt.NewNumericDate(now.Add(accessExp)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        accessJTI,
		},
	}
	accessStr, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(cfg.Secret))
	if err != nil {
		return nil, fmt.Errorf("生成 Access Token 失败: %w", err)
	}

	refreshClaims := Claims{
		UserID:   userID,
		Username: username,
		TenantID: tenantID,
		JTI:      refreshJTI,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.Issuer,
			ExpiresAt: jwt.NewNumericDate(now.Add(refreshExp)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        refreshJTI,
		},
	}
	refreshStr, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(cfg.Secret))
	if err != nil {
		return nil, fmt.Errorf("生成 Refresh Token 失败: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessStr,
		RefreshToken: refreshStr,
		AccessExpire: now.Add(accessExp).Unix(),
	}, nil
}

// ParseToken 解析并验证 Token
func ParseToken(tokenStr string, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
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
