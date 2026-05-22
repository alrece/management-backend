package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"management-backend/internal/auth"
	"management-backend/internal/config"
	pkgauth "management-backend/pkg/authx"
	"management-backend/testutil"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupAuthEnv 初始化测试环境：config + blacklist
func setupAuthEnv(t *testing.T) {
	t.Helper()
	config.C = config.Config{
		JWT: config.JWTConfig{
			Secret:        "test-secret-key-for-testing",
			AccessExpire:  3600,
			RefreshExpire: 86400,
			Issuer:        "test",
		},
	}
	InitAuth(testutil.NewTestRedis(t))
}

func TestAuth_ValidToken(t *testing.T) {
	setupAuthEnv(t)

	cfg := pkgauth.JWTConfig{
		Secret:        config.C.JWT.Secret,
		AccessExpire:  config.C.JWT.AccessExpire,
		RefreshExpire: config.C.JWT.RefreshExpire,
		Issuer:        config.C.JWT.Issuer,
	}
	pair, err := pkgauth.GenerateTokenPair(1, "admin", 100, cfg)
	require.NoError(t, err)

	var gotUID int64
	var gotTID int64

	r := gin.New()
	r.Use(Auth())
	r.GET("/protected", func(c *gin.Context) {
		gotUID = GetUserID(c)
		gotTID = GetTenantID(c)
		c.JSON(http.StatusOK, gin.H{"userId": gotUID, "tenantId": gotTID})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+pair.AccessToken)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, int64(1), gotUID)
	assert.Equal(t, int64(100), gotTID)
}

func TestAuth_TokenWithoutBearerPrefix(t *testing.T) {
	setupAuthEnv(t)

	cfg := pkgauth.JWTConfig{
		Secret:       config.C.JWT.Secret,
		AccessExpire: 3600,
		Issuer:       config.C.JWT.Issuer,
	}
	pair, err := pkgauth.GenerateTokenPair(1, "admin", 100, cfg)
	require.NoError(t, err)

	r := gin.New()
	r.Use(Auth())
	r.GET("/protected", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", pair.AccessToken)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuth_RevokedToken(t *testing.T) {
	setupAuthEnv(t)

	cfg := pkgauth.JWTConfig{
		Secret:       config.C.JWT.Secret,
		AccessExpire: 3600,
		Issuer:       config.C.JWT.Issuer,
	}
	pair, err := pkgauth.GenerateTokenPair(1, "admin", 100, cfg)
	require.NoError(t, err)

	claims, err := auth.ParseAccessToken(pair.AccessToken)
	require.NoError(t, err)

	// 吊销 token
	err = RevokeToken(context.Background(), claims.JTI)
	require.NoError(t, err)

	r := gin.New()
	r.Use(Auth())
	r.GET("/protected", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+pair.AccessToken)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuth_RevokeUserTokens(t *testing.T) {
	setupAuthEnv(t)

	err := RevokeUserTokens(context.Background(), 999)
	assert.NoError(t, err)
}
