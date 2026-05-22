package service

import (
	"context"
	"testing"

	smodel "management-backend/internal/module/system/model"
	"management-backend/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserService_Update(t *testing.T) {
	testutil.NewTestDB(t)
	repo := newMockUserRepo()
	svc := NewUserService(repo)
	ctx := context.Background()

	user := &smodel.User{Username: "upuser", Nickname: "旧昵称", Email: "old@test.com"}
	user.ID = 12001
	user.TenantID = 1
	repo.users[user.ID] = user

	err := svc.Update(ctx, &smodel.UserUpdateReq{
		ID: user.ID, Nickname: "新昵称", Email: "new@test.com",
	}, 1)
	require.NoError(t, err)
	assert.Equal(t, "新昵称", repo.users[user.ID].Nickname)
	assert.Equal(t, "new@test.com", repo.users[user.ID].Email)

	err = svc.Update(ctx, &smodel.UserUpdateReq{ID: 9999}, 1)
	assert.Error(t, err)
}

func TestUserService_Page(t *testing.T) {
	testutil.NewTestDB(t)
	repo := newMockUserRepo()
	svc := NewUserService(repo)
	ctx := context.Background()

	for i := int64(0); i < 3; i++ {
		u := &smodel.User{Username: "user" + string(rune('A'+i))}
		u.ID = 12101 + i
		u.TenantID = 1
		repo.users[u.ID] = u
	}

	list, total, err := svc.Page(ctx, &smodel.UserPageReq{})
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, list, 3)
}

func TestUserService_GetByUsername(t *testing.T) {
	testutil.NewTestDB(t)
	repo := newMockUserRepo()
	svc := NewUserService(repo)
	ctx := context.Background()

	u := &smodel.User{Username: "findme", Nickname: "找我"}
	u.ID = 12201
	u.TenantID = 1
	repo.users[u.ID] = u

	found, err := svc.GetByUsername(ctx, "findme")
	require.NoError(t, err)
	assert.Equal(t, "找我", found.Nickname)

	_, err = svc.GetByUsername(ctx, "notexist")
	assert.Error(t, err)
}

func TestUserService_GetRawByID(t *testing.T) {
	testutil.NewTestDB(t)
	repo := newMockUserRepo()
	svc := NewUserService(repo)
	ctx := context.Background()

	u := &smodel.User{Username: "rawuser"}
	u.ID = 12301
	u.TenantID = 1
	repo.users[u.ID] = u

	raw, err := svc.GetRawByID(ctx, u.ID)
	require.NoError(t, err)
	assert.Equal(t, "rawuser", raw.Username)

	_, err = svc.GetRawByID(ctx, 9999)
	assert.Error(t, err)
}

func TestUserService_UpdateProfile(t *testing.T) {
	testutil.NewTestDB(t)
	repo := newMockUserRepo()
	svc := NewUserService(repo)
	ctx := context.Background()

	u := &smodel.User{Username: "profile", Nickname: "旧"}
	u.ID = 12401
	u.TenantID = 1
	repo.users[u.ID] = u

	err := svc.UpdateProfile(ctx, &UserProfileUpdateReq{
		ID: u.ID, Nickname: "新", Email: "new@test.com",
	})
	require.NoError(t, err)
	assert.Equal(t, "新", repo.users[u.ID].Nickname)

	err = svc.UpdateProfile(ctx, &UserProfileUpdateReq{ID: 9999})
	assert.Error(t, err)
}
