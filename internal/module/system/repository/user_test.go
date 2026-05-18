package repository

import (
	"context"
	"testing"

	smodel "management-backend/internal/module/system/model"
	"management-backend/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepo_Create(t *testing.T) {
	db := testutil.NewTestDB(t)
	repo := NewUserRepo(db)
	ctx := context.Background()

	user := &smodel.User{
		Username: "testuser",
		Password: "hashedpw",
		Nickname: "测试用户",
	}
	user.ID = 1001
	user.TenantID = 1

	err := repo.Create(ctx, user)
	assert.NoError(t, err)

	found, err := repo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, "testuser", found.Username)
	assert.Equal(t, "测试用户", found.Nickname)
}

func TestUserRepo_GetByUsername(t *testing.T) {
	db := testutil.NewTestDB(t)
	repo := NewUserRepo(db)
	ctx := context.Background()

	testutil.SeedUser(t, db, "alice", "pw123")

	found, err := repo.GetByUsername(ctx, "alice")
	require.NoError(t, err)
	assert.Equal(t, "alice", found.Username)

	_, err = repo.GetByUsername(ctx, "notexist")
	assert.Error(t, err)
}

func TestUserRepo_Update(t *testing.T) {
	db := testutil.NewTestDB(t)
	repo := NewUserRepo(db)
	ctx := context.Background()

	user := testutil.SeedUser(t, db, "bob", "pw123")
	user.Nickname = "新昵称"
	user.Email = "bob@test.com"

	err := repo.Update(ctx, user)
	assert.NoError(t, err)

	updated, err := repo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, "新昵称", updated.Nickname)
	assert.Equal(t, "bob@test.com", updated.Email)
}

func TestUserRepo_Delete(t *testing.T) {
	db := testutil.NewTestDB(t)
	repo := NewUserRepo(db)
	ctx := context.Background()

	user := testutil.SeedUser(t, db, "charlie", "pw123")

	err := repo.Delete(ctx, user.ID)
	assert.NoError(t, err)

	_, err = repo.GetByID(ctx, user.ID)
	assert.Error(t, err)
}

func TestUserRepo_UpdatePassword(t *testing.T) {
	db := testutil.NewTestDB(t)
	repo := NewUserRepo(db)
	ctx := context.Background()

	user := testutil.SeedUser(t, db, "dave", "oldpw")
	user.Password = "newhashedpw"

	err := repo.UpdatePassword(ctx, user)
	assert.NoError(t, err)

	found, err := repo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, "newhashedpw", found.Password)
}

func TestUserRepo_Page(t *testing.T) {
	db := testutil.NewTestDB(t)
	repo := NewUserRepo(db)
	ctx := context.Background()

	testutil.SeedUser(t, db, "user1", "pw")
	testutil.SeedUser(t, db, "user2", "pw")
	testutil.SeedUser(t, db, "admin1", "pw")

	tests := []struct {
		name      string
		req       *smodel.UserPageReq
		wantCount int
		wantTotal int64
	}{
		{
			name:      "全部查询",
			req:       &smodel.UserPageReq{},
			wantCount: 3,
			wantTotal: 3,
		},
		{
			name:      "按用户名模糊查询",
			req:       &smodel.UserPageReq{Username: "admin"},
			wantCount: 1,
			wantTotal: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.req.Page = 1
			tt.req.PageSize = 10
			list, total, err := repo.Page(ctx, tt.req)
			require.NoError(t, err)
			assert.Equal(t, tt.wantTotal, total)
			assert.Len(t, list, tt.wantCount)
		})
	}
}
