package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// --- Context helpers ---

func TestContext_UserID(t *testing.T) {
	ctx := context.Background()
	assert.Equal(t, int64(0), GetUserID(ctx))

	ctx = WithUserID(ctx, 123)
	assert.Equal(t, int64(123), GetUserID(ctx))
}

func TestContext_Username(t *testing.T) {
	ctx := context.Background()
	assert.Equal(t, "", GetUsername(ctx))

	ctx = WithUsername(ctx, "admin")
	assert.Equal(t, "admin", GetUsername(ctx))
}

func TestContext_TenantID(t *testing.T) {
	ctx := context.Background()
	assert.Equal(t, int64(0), GetTenantID(ctx))

	ctx = WithTenantID(ctx, 999)
	assert.Equal(t, int64(999), GetTenantID(ctx))
}

func TestContext_RequestID(t *testing.T) {
	ctx := context.Background()
	assert.Equal(t, "", GetRequestID(ctx))

	ctx = WithRequestID(ctx, "req-123")
	assert.Equal(t, "req-123", GetRequestID(ctx))
}

func TestContext_TenantDB(t *testing.T) {
	ctx := context.Background()
	assert.Nil(t, GetTenantDB(ctx))

	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	ctx = WithTenantDB(ctx, db)
	assert.NotNil(t, GetTenantDB(ctx))
}

// --- RequestID middleware ---

func TestRequestID_Generates(t *testing.T) {
	r := gin.New()
	r.Use(RequestID())
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Header().Get("X-Request-ID"))
}

func TestRequestID_PassesExisting(t *testing.T) {
	r := gin.New()
	r.Use(RequestID())
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Request-ID", "existing-id")
	r.ServeHTTP(w, req)

	assert.Equal(t, "existing-id", w.Header().Get("X-Request-ID"))
}

// --- Recovery middleware ---

func TestRecovery_RecoversPanic(t *testing.T) {
	logger := zap.NewNop()
	r := gin.New()
	r.Use(Recovery(logger))
	r.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// --- SecurityHeaders middleware ---

func TestSecurityHeaders(t *testing.T) {
	r := gin.New()
	r.Use(SecurityHeaders())
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", w.Header().Get("X-XSS-Protection"))
	assert.NotEmpty(t, w.Header().Get("Strict-Transport-Security"))
	assert.Equal(t, "default-src 'self'", w.Header().Get("Content-Security-Policy"))
}
