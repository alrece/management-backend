package repository

import (
	"context"
	"testing"

	smodel "management-backend/internal/module/system/model"
	"management-backend/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTenantRepo_Create(t *testing.T) {
	db := testutil.NewTestDB(t)
	repo := NewTenantRepo(db)
	ctx := context.Background()

	tenant := testutil.SeedTenant(t, db, "测试租户")
	assert.NotZero(t, tenant.ID)

	found, err := repo.GetByID(ctx, tenant.ID)
	require.NoError(t, err)
	assert.Equal(t, "测试租户", found.TenantName)
}

func TestTenantRepo_Update(t *testing.T) {
	db := testutil.NewTestDB(t)
	repo := NewTenantRepo(db)
	ctx := context.Background()

	tenant := testutil.SeedTenant(t, db, "旧名称")
	tenant.TenantName = "新名称"

	err := repo.Update(ctx, tenant)
	assert.NoError(t, err)

	found, err := repo.GetByID(ctx, tenant.ID)
	require.NoError(t, err)
	assert.Equal(t, "新名称", found.TenantName)
}

func TestTenantRepo_Delete(t *testing.T) {
	db := testutil.NewTestDB(t)
	repo := NewTenantRepo(db)
	ctx := context.Background()

	tenant := testutil.SeedTenant(t, db, "待删除")

	err := repo.Delete(ctx, tenant.ID)
	assert.NoError(t, err)

	_, err = repo.GetByID(ctx, tenant.ID)
	assert.Error(t, err)
}

func TestTenantRepo_UpdateInitStatus(t *testing.T) {
	db := testutil.NewTestDB(t)
	repo := NewTenantRepo(db)
	ctx := context.Background()

	tenant := testutil.SeedTenant(t, db, "状态测试")
	ok, err := repo.UpdateInitStatus(ctx, tenant.ID, smodel.InitStatusReady, smodel.InitStatusPending)
	require.NoError(t, err)
	assert.True(t, ok)

	// CAS 失败：旧状态不匹配
	ok, err = repo.UpdateInitStatus(ctx, tenant.ID, smodel.InitStatusReady, smodel.InitStatusPending)
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestTenantRepo_Page(t *testing.T) {
	db := testutil.NewTestDB(t)
	repo := NewTenantRepo(db)
	ctx := context.Background()

	testutil.SeedTenant(t, db, "租户A")
	testutil.SeedTenant(t, db, "租户B")

	req := &smodel.TenantPageReq{}
	req.Page = 1
	req.PageSize = 10

	list, total, err := repo.Page(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, list, 2)
}

func TestTenantRepo_Page_FilterByName(t *testing.T) {
	db := testutil.NewTestDB(t)
	repo := NewTenantRepo(db)
	ctx := context.Background()

	testutil.SeedTenant(t, db, "租户Alpha")
	testutil.SeedTenant(t, db, "租户Beta")

	req := &smodel.TenantPageReq{TenantName: "Alpha"}
	req.Page = 1
	req.PageSize = 10

	list, total, err := repo.Page(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	require.Len(t, list, 1)
	assert.Contains(t, list[0].TenantName, "Alpha")
}
