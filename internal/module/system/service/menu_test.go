package service

import (
	"context"
	"testing"

	smodel "management-backend/internal/module/system/model"
	"management-backend/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMenuService_Create(t *testing.T) {
	testutil.NewTestDB(t)
	repo := newMockMenuRepo()
	svc := NewMenuService(repo)
	ctx := context.Background()

	id, err := svc.Create(ctx, &smodel.MenuCreateReq{
		MenuName: "系统管理", MenuType: 1, Path: "/system",
	}, 1)
	require.NoError(t, err)
	assert.NotZero(t, id)
}

func TestMenuService_Update(t *testing.T) {
	testutil.NewTestDB(t)
	repo := newMockMenuRepo()
	svc := NewMenuService(repo)
	ctx := context.Background()

	menu := &smodel.Menu{MenuName: "旧名称", MenuType: 1}
	menu.ID = 6001
	repo.menus[menu.ID] = menu

	err := svc.Update(ctx, &smodel.MenuUpdateReq{
		ID: menu.ID, MenuName: "新名称", MenuType: 2,
	}, 1)
	require.NoError(t, err)
	assert.Equal(t, "新名称", repo.menus[menu.ID].MenuName)

	err = svc.Update(ctx, &smodel.MenuUpdateReq{ID: 9999}, 1)
	assert.Error(t, err)
}

func TestMenuService_Delete_HasChildren(t *testing.T) {
	testutil.NewTestDB(t)
	repo := newMockMenuRepo()
	svc := NewMenuService(repo)
	ctx := context.Background()

	parent := &smodel.Menu{MenuName: "父菜单"}
	parent.ID = 6002
	repo.menus[parent.ID] = parent

	child := &smodel.Menu{MenuName: "子菜单", ParentID: parent.ID}
	child.ID = 6003
	repo.menus[child.ID] = child

	// 有子菜单不能删
	err := svc.Delete(ctx, parent.ID)
	assert.Error(t, err)
}

func TestMenuService_Delete_NoChildren(t *testing.T) {
	testutil.NewTestDB(t)
	repo := newMockMenuRepo()
	svc := NewMenuService(repo)
	ctx := context.Background()

	menu := &smodel.Menu{MenuName: "叶子菜单"}
	menu.ID = 6004
	repo.menus[menu.ID] = menu

	err := svc.Delete(ctx, menu.ID)
	require.NoError(t, err)
	_, ok := repo.menus[menu.ID]
	assert.False(t, ok)
}

func TestMenuService_Tree(t *testing.T) {
	testutil.NewTestDB(t)
	repo := newMockMenuRepo()
	svc := NewMenuService(repo)
	ctx := context.Background()

	p := &smodel.Menu{MenuName: "父", ParentID: 0}
	p.ID = 7001
	repo.menus[p.ID] = p

	c := &smodel.Menu{MenuName: "子", ParentID: p.ID}
	c.ID = 7002
	repo.menus[c.ID] = c

	tree, err := svc.Tree(ctx)
	require.NoError(t, err)
	require.Len(t, tree, 1)
	assert.Equal(t, "父", tree[0].MenuName)
	require.Len(t, tree[0].Children, 1)
	assert.Equal(t, "子", tree[0].Children[0].MenuName)
}

func TestMenuService_GetByIDs(t *testing.T) {
	testutil.NewTestDB(t)
	repo := newMockMenuRepo()
	svc := NewMenuService(repo)
	ctx := context.Background()

	repo.menus[8001] = &smodel.Menu{MenuName: "A"}
	repo.menus[8002] = &smodel.Menu{MenuName: "B"}

	menus, err := svc.GetByIDs(ctx, []int64{8001, 8002, 8003})
	require.NoError(t, err)
	assert.Len(t, menus, 2)
}
