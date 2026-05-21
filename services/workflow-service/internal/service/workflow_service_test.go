package service

import (
	"testing"

	"management-backend/services/workflow-service/internal/model"
	"management-backend/services/workflow-service/internal/repository"
	"management-backend/services/workflow-service/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestCategoryService(t *testing.T) CategoryService {
	t.Helper()
	db := testutil.NewTestDB(t)
	return NewCategoryService(repository.NewCategoryRepo(db))
}

func TestCategoryService_Create(t *testing.T) {
	svc := newTestCategoryService(t)
	ctx := testutil.Ctx()

	id, err := svc.Create(ctx, &model.CategoryCreateReq{Name: "测试分类", Sort: 1}, 1, 100)
	require.NoError(t, err)
	assert.NotZero(t, id)
}

func TestCategoryService_Tree(t *testing.T) {
	svc := newTestCategoryService(t)
	ctx := testutil.Ctx()

	_, _ = svc.Create(ctx, &model.CategoryCreateReq{Name: "父级", Sort: 0}, 1, 100)
	_, _ = svc.Create(ctx, &model.CategoryCreateReq{Name: "子级", ParentID: 0, Sort: 1}, 1, 100)

	tree, err := svc.Tree(ctx, 100)
	require.NoError(t, err)
	assert.Len(t, tree, 2)
}

func TestCategoryService_Delete(t *testing.T) {
	svc := newTestCategoryService(t)
	ctx := testutil.Ctx()

	id, _ := svc.Create(ctx, &model.CategoryCreateReq{Name: "删除测试"}, 1, 100)
	err := svc.Delete(ctx, id)
	require.NoError(t, err)

	tree, _ := svc.Tree(ctx, 100)
	assert.Empty(t, tree)
}

func TestCategoryService_TenantIsolation(t *testing.T) {
	svc := newTestCategoryService(t)
	ctx := testutil.Ctx()

	_, _ = svc.Create(ctx, &model.CategoryCreateReq{Name: "租户A"}, 1, 100)
	_, _ = svc.Create(ctx, &model.CategoryCreateReq{Name: "租户B"}, 1, 200)

	tree100, _ := svc.Tree(ctx, 100)
	tree200, _ := svc.Tree(ctx, 200)

	assert.Len(t, tree100, 1)
	assert.Len(t, tree200, 1)
	assert.Equal(t, "租户A", tree100[0].Name)
	assert.Equal(t, "租户B", tree200[0].Name)
}
