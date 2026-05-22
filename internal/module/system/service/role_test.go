package service

import (
	"context"
	"testing"

	smodel "management-backend/internal/module/system/model"
	"management-backend/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// --- Mock Repos ---

type mockRoleRepo struct {
	roles    map[int64]*smodel.Role
	roleMenu map[int64][]int64
}

func newMockRoleRepo() *mockRoleRepo {
	return &mockRoleRepo{
		roles:    make(map[int64]*smodel.Role),
		roleMenu: make(map[int64][]int64),
	}
}

func (m *mockRoleRepo) Create(_ context.Context, role *smodel.Role) error {
	m.roles[role.ID] = role
	return nil
}

func (m *mockRoleRepo) Update(_ context.Context, role *smodel.Role) error {
	m.roles[role.ID] = role
	return nil
}

func (m *mockRoleRepo) Delete(_ context.Context, id int64) error {
	delete(m.roles, id)
	return nil
}

func (m *mockRoleRepo) GetByID(_ context.Context, id int64) (*smodel.Role, error) {
	r, ok := m.roles[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return r, nil
}

func (m *mockRoleRepo) GetByRoleCode(_ context.Context, code string) (*smodel.Role, error) {
	for _, r := range m.roles {
		if r.RoleCode == code {
			return r, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockRoleRepo) Page(_ context.Context, _ *smodel.RolePageReq) ([]smodel.Role, int64, error) {
	var list []smodel.Role
	for _, r := range m.roles {
		list = append(list, *r)
	}
	return list, int64(len(list)), nil
}

func (m *mockRoleRepo) GetByUserID(_ context.Context, _ int64) ([]smodel.Role, error) {
	return nil, nil
}

func (m *mockRoleRepo) AssignMenus(_ context.Context, roleID int64, menuIDs []int64, _ int64) error {
	m.roleMenu[roleID] = menuIDs
	return nil
}

func (m *mockRoleRepo) GetMenuIDsByRoleID(_ context.Context, roleID int64) ([]int64, error) {
	return m.roleMenu[roleID], nil
}

func (m *mockRoleRepo) GetMenuIDsByRoleIDs(_ context.Context, _ []int64) ([]int64, error) {
	return nil, nil
}

type mockMenuRepo struct {
	menus map[int64]*smodel.Menu
}

func newMockMenuRepo() *mockMenuRepo {
	return &mockMenuRepo{menus: make(map[int64]*smodel.Menu)}
}

func (m *mockMenuRepo) Create(_ context.Context, menu *smodel.Menu) error {
	m.menus[menu.ID] = menu
	return nil
}

func (m *mockMenuRepo) Update(_ context.Context, menu *smodel.Menu) error {
	m.menus[menu.ID] = menu
	return nil
}

func (m *mockMenuRepo) Delete(_ context.Context, id int64) error {
	delete(m.menus, id)
	return nil
}

func (m *mockMenuRepo) GetByID(_ context.Context, id int64) (*smodel.Menu, error) {
	menu, ok := m.menus[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return menu, nil
}

func (m *mockMenuRepo) List(_ context.Context) ([]smodel.Menu, error) {
	var list []smodel.Menu
	for _, m := range m.menus {
		list = append(list, *m)
	}
	return list, nil
}

func (m *mockMenuRepo) GetByIDs(_ context.Context, ids []int64) ([]smodel.Menu, error) {
	var list []smodel.Menu
	for _, id := range ids {
		if menu, ok := m.menus[id]; ok {
			list = append(list, *menu)
		}
	}
	return list, nil
}

func (m *mockMenuRepo) HasChildren(_ context.Context, id int64) (bool, error) {
	for _, menu := range m.menus {
		if menu.ParentID == id {
			return true, nil
		}
	}
	return false, nil
}

// --- Role Tests ---

func TestRoleService_Create(t *testing.T) {
	testutil.NewTestDB(t)
	repo := newMockRoleRepo()
	svc := NewRoleService(repo, newMockMenuRepo(), nil)
	ctx := context.Background()

	id, err := svc.Create(ctx, &smodel.RoleCreateReq{
		RoleName: "管理员", RoleCode: "admin", Status: 0, Sort: 1,
	}, 1)
	require.NoError(t, err)
	assert.NotZero(t, id)

	// 重复角色编码
	_, err = svc.Create(ctx, &smodel.RoleCreateReq{
		RoleName: "管理员2", RoleCode: "admin",
	}, 1)
	assert.Error(t, err)
}

func TestRoleService_Update(t *testing.T) {
	testutil.NewTestDB(t)
	repo := newMockRoleRepo()
	svc := NewRoleService(repo, newMockMenuRepo(), nil)
	ctx := context.Background()

	role := &smodel.Role{RoleName: "测试", RoleCode: "test"}
	role.ID = 5001
	repo.roles[role.ID] = role

	err := svc.Update(ctx, &smodel.RoleUpdateReq{
		ID: role.ID, RoleName: "更新后", RoleCode: "test",
	}, 1)
	require.NoError(t, err)
	assert.Equal(t, "更新后", repo.roles[role.ID].RoleName)

	// 不存在的角色
	err = svc.Update(ctx, &smodel.RoleUpdateReq{ID: 9999}, 1)
	assert.Error(t, err)
}

func TestRoleService_Delete(t *testing.T) {
	testutil.NewTestDB(t)
	repo := newMockRoleRepo()
	svc := NewRoleService(repo, newMockMenuRepo(), nil)
	ctx := context.Background()

	role := &smodel.Role{RoleName: "删除", RoleCode: "del"}
	role.ID = 5002
	repo.roles[role.ID] = role

	err := svc.Delete(ctx, role.ID)
	require.NoError(t, err)
	_, ok := repo.roles[role.ID]
	assert.False(t, ok)

	err = svc.Delete(ctx, 9999)
	assert.Error(t, err)
}

func TestRoleService_GetByID(t *testing.T) {
	testutil.NewTestDB(t)
	repo := newMockRoleRepo()
	svc := NewRoleService(repo, newMockMenuRepo(), nil)
	ctx := context.Background()

	role := &smodel.Role{RoleName: "获取", RoleCode: "get"}
	role.ID = 5003
	repo.roles[role.ID] = role

	resp, err := svc.GetByID(ctx, role.ID)
	require.NoError(t, err)
	assert.Equal(t, "获取", resp.RoleName)

	_, err = svc.GetByID(ctx, 9999)
	assert.Error(t, err)
}

func TestRoleService_AssignMenus(t *testing.T) {
	testutil.NewTestDB(t)
	repo := newMockRoleRepo()
	menuRepo := newMockMenuRepo()
	svc := NewRoleService(repo, menuRepo, nil)
	ctx := context.Background()

	role := &smodel.Role{RoleName: "分配", RoleCode: "assign"}
	role.ID = 5004
	repo.roles[role.ID] = role

	menuRepo.menus[101] = &smodel.Menu{Permission: "system:user:query"}
	menuRepo.menus[102] = &smodel.Menu{Permission: "system:role:query"}

	err := svc.AssignMenus(ctx, &smodel.RoleMenuAssignReq{
		RoleID:  role.ID,
		MenuIDs: []int64{101, 102},
	}, 1, 100)
	require.NoError(t, err)
	assert.Equal(t, []int64{101, 102}, repo.roleMenu[role.ID])

	// 不存在的角色
	err = svc.AssignMenus(ctx, &smodel.RoleMenuAssignReq{RoleID: 9999}, 1, 100)
	assert.Error(t, err)
}
