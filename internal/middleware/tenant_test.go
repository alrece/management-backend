package middleware

import (
	"context"
	"testing"

	smodel "management-backend/internal/module/system/model"
	pkgmiddleware "management-backend/pkg/middleware"
	"management-backend/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTenantIsolation 验证独立数据库模式下租户数据隔离
func TestTenantIsolation(t *testing.T) {
	dbA := testutil.NewTestDB(t)
	dbB := testutil.NewTestDB(t)

	// 租户 A 创建用户
	userA := &smodel.User{Username: "userA", Nickname: "租户A用户"}
	userA.ID = 1001
	userA.TenantID = 100
	userA.Status = 0
	require.NoError(t, dbA.Create(userA).Error)

	// 租户 B 创建用户
	userB := &smodel.User{Username: "userB", Nickname: "租户B用户"}
	userB.ID = 1002
	userB.TenantID = 200
	userB.Status = 0
	require.NoError(t, dbB.Create(userB).Error)

	// 租户 A 查不到租户 B 的数据
	var countA int64
	dbA.Model(&smodel.User{}).Where("username = ?", "userB").Count(&countA)
	assert.Equal(t, int64(0), countA)

	// 租户 B 查不到租户 A 的数据
	var countB int64
	dbB.Model(&smodel.User{}).Where("username = ?", "userA").Count(&countB)
	assert.Equal(t, int64(0), countB)
}

// TestContextTenantDB 验证 context 注入租户 DB 的正确性
func TestContextTenantDB(t *testing.T) {
	db := testutil.NewTestDB(t)
	ctx := context.Background()

	assert.Nil(t, pkgmiddleware.GetTenantDB(ctx))

	ctx = pkgmiddleware.WithTenantDB(ctx, db)
	gotDB := pkgmiddleware.GetTenantDB(ctx)
	assert.NotNil(t, gotDB)

	user := &smodel.User{Username: "ctxuser", Nickname: "上下文用户"}
	user.ID = 2001
	user.TenantID = 500
	user.Status = 0
	require.NoError(t, gotDB.Create(user).Error)

	var found smodel.User
	require.NoError(t, gotDB.Where("username = ?", "ctxuser").First(&found).Error)
	assert.Equal(t, "上下文用户", found.Nickname)
}

// TestTenantIDContext 验证 tenant ID 在 context 中的存取
func TestTenantIDContext(t *testing.T) {
	ctx := context.Background()

	assert.Equal(t, int64(0), pkgmiddleware.GetTenantID(ctx))

	ctx = pkgmiddleware.WithTenantID(ctx, 12345)
	assert.Equal(t, int64(12345), pkgmiddleware.GetTenantID(ctx))
}

// TestUserIDContext 验证 user ID 在 context 中的存取
func TestUserIDContext(t *testing.T) {
	ctx := context.Background()

	assert.Equal(t, int64(0), pkgmiddleware.GetUserID(ctx))

	ctx = pkgmiddleware.WithUserID(ctx, 67890)
	assert.Equal(t, int64(67890), pkgmiddleware.GetUserID(ctx))
}
