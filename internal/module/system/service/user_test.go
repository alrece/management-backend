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

// mockUserRepo 测试用 mock Repository
type mockUserRepo struct {
	users  map[int64]*smodel.User
	nextID int64
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[int64]*smodel.User), nextID: 10001}
}

func (m *mockUserRepo) Create(_ context.Context, user *smodel.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepo) Update(_ context.Context, user *smodel.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepo) Delete(_ context.Context, id int64) error {
	delete(m.users, id)
	return nil
}

func (m *mockUserRepo) GetByID(_ context.Context, id int64) (*smodel.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return u, nil
}

func (m *mockUserRepo) GetByUsername(_ context.Context, username string) (*smodel.User, error) {
	for _, u := range m.users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockUserRepo) UpdatePassword(_ context.Context, user *smodel.User) error {
	m.users[user.ID].Password = user.Password
	return nil
}

func (m *mockUserRepo) Page(_ context.Context, req *smodel.UserPageReq) ([]smodel.User, int64, error) {
	var list []smodel.User
	for _, u := range m.users {
		if req.Username != "" {
			continue
		}
		list = append(list, *u)
	}
	return list, int64(len(list)), nil
}

func TestUserService_Create(t *testing.T) {
	testutil.NewTestDB(t)
	repo := newMockUserRepo()
	svc := NewUserService(repo)
	ctx := context.Background()

	req := &smodel.UserCreateReq{
		Username: "newuser",
		Password: "password123",
		Nickname: "新用户",
	}

	id, err := svc.Create(ctx, req, 1, 100)
	require.NoError(t, err)
	assert.NotZero(t, id)

	// 重复用户名应失败
	_, err = svc.Create(ctx, req, 1, 100)
	assert.Error(t, err)
}

func TestUserService_GetByID(t *testing.T) {
	testutil.NewTestDB(t)
	repo := newMockUserRepo()
	svc := NewUserService(repo)
	ctx := context.Background()

	user := &smodel.User{Username: "getuser", Nickname: "获取用户"}
	user.ID = 2001
	user.TenantID = 1
	repo.users[user.ID] = user

	resp, err := svc.GetByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, "getuser", resp.Username)
	assert.Equal(t, "获取用户", resp.Nickname)

	_, err = svc.GetByID(ctx, 9999)
	assert.Error(t, err)
}

func TestUserService_Delete(t *testing.T) {
	testutil.NewTestDB(t)
	repo := newMockUserRepo()
	svc := NewUserService(repo)
	ctx := context.Background()

	user := &smodel.User{Username: "deluser"}
	user.ID = 3001
	user.TenantID = 1
	repo.users[user.ID] = user

	err := svc.Delete(ctx, user.ID)
	require.NoError(t, err)
	_, ok := repo.users[user.ID]
	assert.False(t, ok)

	err = svc.Delete(ctx, 9999)
	assert.Error(t, err)
}

func TestUserService_UpdatePassword(t *testing.T) {
	testutil.NewTestDB(t)
	repo := newMockUserRepo()
	svc := NewUserService(repo)
	ctx := context.Background()

	user := &smodel.User{Username: "pwuser", Password: "oldhash"}
	user.ID = 4001
	user.TenantID = 1
	repo.users[user.ID] = user

	err := svc.UpdatePassword(ctx, user.ID, "newhash")
	require.NoError(t, err)
	assert.Equal(t, "newhash", repo.users[user.ID].Password)
}
