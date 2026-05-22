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

type mockPostRepo struct {
	posts map[int64]*smodel.Post
}

func newMockPostRepo() *mockPostRepo {
	return &mockPostRepo{posts: make(map[int64]*smodel.Post)}
}

func (m *mockPostRepo) Create(_ context.Context, post *smodel.Post) error {
	m.posts[post.ID] = post
	return nil
}

func (m *mockPostRepo) Update(_ context.Context, post *smodel.Post) error {
	m.posts[post.ID] = post
	return nil
}

func (m *mockPostRepo) Delete(_ context.Context, id int64) error {
	delete(m.posts, id)
	return nil
}

func (m *mockPostRepo) GetByID(_ context.Context, id int64) (*smodel.Post, error) {
	p, ok := m.posts[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return p, nil
}

func (m *mockPostRepo) GetByCode(_ context.Context, code string) (*smodel.Post, error) {
	for _, p := range m.posts {
		if p.PostCode == code {
			return p, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockPostRepo) Page(_ context.Context, _ *smodel.PostPageReq) ([]smodel.Post, int64, error) {
	var list []smodel.Post
	for _, p := range m.posts {
		list = append(list, *p)
	}
	return list, int64(len(list)), nil
}

func TestPostService_Create(t *testing.T) {
	testutil.NewTestDB(t)
	repo := newMockPostRepo()
	svc := NewPostService(repo)
	ctx := context.Background()

	id, err := svc.Create(ctx, &smodel.PostCreateReq{
		PostCode: "dev", PostName: "开发工程师", Sort: 1, Status: 0,
	}, 1)
	require.NoError(t, err)
	assert.NotZero(t, id)

	// 重复编码
	_, err = svc.Create(ctx, &smodel.PostCreateReq{
		PostCode: "dev", PostName: "重复",
	}, 1)
	assert.Error(t, err)
}

func TestPostService_Update(t *testing.T) {
	testutil.NewTestDB(t)
	repo := newMockPostRepo()
	svc := NewPostService(repo)
	ctx := context.Background()

	post := &smodel.Post{PostCode: "dev", PostName: "开发"}
	post.ID = 11001
	repo.posts[post.ID] = post

	err := svc.Update(ctx, &smodel.PostUpdateReq{
		ID: post.ID, PostCode: "dev", PostName: "高级开发",
	}, 1)
	require.NoError(t, err)
	assert.Equal(t, "高级开发", repo.posts[post.ID].PostName)

	err = svc.Update(ctx, &smodel.PostUpdateReq{ID: 9999}, 1)
	assert.Error(t, err)
}

func TestPostService_Delete(t *testing.T) {
	testutil.NewTestDB(t)
	repo := newMockPostRepo()
	svc := NewPostService(repo)
	ctx := context.Background()

	post := &smodel.Post{PostCode: "pm", PostName: "产品"}
	post.ID = 11002
	repo.posts[post.ID] = post

	err := svc.Delete(ctx, post.ID)
	require.NoError(t, err)
	_, ok := repo.posts[post.ID]
	assert.False(t, ok)

	err = svc.Delete(ctx, 9999)
	assert.Error(t, err)
}

func TestPostService_Page(t *testing.T) {
	testutil.NewTestDB(t)
	repo := newMockPostRepo()
	svc := NewPostService(repo)
	ctx := context.Background()

	p1 := &smodel.Post{PostCode: "dev", PostName: "开发"}
	p1.ID = 11003
	repo.posts[p1.ID] = p1

	p2 := &smodel.Post{PostCode: "qa", PostName: "测试"}
	p2.ID = 11004
	repo.posts[p2.ID] = p2

	list, total, err := svc.Page(ctx, &smodel.PostPageReq{})
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, list, 2)
}
