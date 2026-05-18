package repository

import (
	"context"
	"testing"

	smodel "management-backend/internal/module/system/model"
	"management-backend/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoleRepo_Create(t *testing.T) {
	db := testutil.NewTestDB(t)
	repo := NewRoleRepo(db)
	ctx := context.Background()

	role := testutil.SeedRole(t, db, "管理员", "admin")
	assert.NotZero(t, role.ID)

	found, err := repo.GetByID(ctx, role.ID)
	require.NoError(t, err)
	assert.Equal(t, "管理员", found.RoleName)
	assert.Equal(t, "admin", found.RoleCode)
}

func TestRoleRepo_GetByRoleCode(t *testing.T) {
	db := testutil.NewTestDB(t)
	repo := NewRoleRepo(db)
	ctx := context.Background()

	testutil.SeedRole(t, db, "普通用户", "user")

	found, err := repo.GetByRoleCode(ctx, "user")
	require.NoError(t, err)
	assert.Equal(t, "普通用户", found.RoleName)

	_, err = repo.GetByRoleCode(ctx, "notexist")
	assert.Error(t, err)
}

func TestRoleRepo_Update(t *testing.T) {
	db := testutil.NewTestDB(t)
	repo := NewRoleRepo(db)
	ctx := context.Background()

	role := testutil.SeedRole(t, db, "旧名称", "editor")
	role.RoleName = "新名称"

	err := repo.Update(ctx, role)
	assert.NoError(t, err)

	found, err := repo.GetByID(ctx, role.ID)
	require.NoError(t, err)
	assert.Equal(t, "新名称", found.RoleName)
}

func TestRoleRepo_Delete(t *testing.T) {
	db := testutil.NewTestDB(t)
	repo := NewRoleRepo(db)
	ctx := context.Background()

	role := testutil.SeedRole(t, db, "临时角色", "temp")

	err := repo.Delete(ctx, role.ID)
	assert.NoError(t, err)

	_, err = repo.GetByID(ctx, role.ID)
	assert.Error(t, err)
}

func TestRoleRepo_AssignMenus(t *testing.T) {
	db := testutil.NewTestDB(t)
	repo := NewRoleRepo(db)
	ctx := context.Background()

	role := testutil.SeedRole(t, db, "测试角色", "test")
	menuIDs := []int64{100, 200, 300}

	err := repo.AssignMenus(ctx, role.ID, menuIDs, 1)
	require.NoError(t, err)

	ids, err := repo.GetMenuIDsByRoleID(ctx, role.ID)
	require.NoError(t, err)
	assert.ElementsMatch(t, menuIDs, ids)
}

func TestRoleRepo_Page(t *testing.T) {
	db := testutil.NewTestDB(t)
	repo := NewRoleRepo(db)
	ctx := context.Background()

	testutil.SeedRole(t, db, "超级管理", "super_admin")
	testutil.SeedRole(t, db, "普通角色", "common")

	req := &smodel.RolePageReq{}
	req.Page = 1
	req.PageSize = 10

	list, total, err := repo.Page(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, list, 2)
}
