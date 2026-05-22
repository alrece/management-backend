package middleware

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)

	assert.Equal(t, int64(0), GetUserID(c))

	c.Set("userId", int64(42))
	assert.Equal(t, int64(42), GetUserID(c))
}

func TestGetTenantID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)

	assert.Equal(t, int64(0), GetTenantID(c))

	c.Set("tenantId", int64(99))
	assert.Equal(t, int64(99), GetTenantID(c))
}

func TestGetString(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)

	assert.Equal(t, "", GetString(c, "username"))

	c.Set("username", "admin")
	assert.Equal(t, "admin", GetString(c, "username"))

	// 非 string 类型应返回空
	c.Set("count", 123)
	assert.Equal(t, "", GetString(c, "count"))
}

func TestGetUserIDStr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)

	c.Set("userId", int64(42))
	assert.Equal(t, "42", GetUserIDStr(c))
}

func TestGetTenantIDStr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)

	c.Set("tenantId", int64(99))
	assert.Equal(t, "99", GetTenantIDStr(c))
}
