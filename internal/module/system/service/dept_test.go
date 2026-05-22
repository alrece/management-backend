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

type mockDeptRepo struct {
	depts map[int64]*smodel.Dept
}

func newMockDeptRepo() *mockDeptRepo {
	return &mockDeptRepo{depts: make(map[int64]*smodel.Dept)}
}

func (m *mockDeptRepo) Create(_ context.Context, dept *smodel.Dept) error {
	m.depts[dept.ID] = dept
	return nil
}

func (m *mockDeptRepo) Update(_ context.Context, dept *smodel.Dept) error {
	m.depts[dept.ID] = dept
	return nil
}

func (m *mockDeptRepo) Delete(_ context.Context, id int64) error {
	delete(m.depts, id)
	return nil
}

func (m *mockDeptRepo) GetByID(_ context.Context, id int64) (*smodel.Dept, error) {
	d, ok := m.depts[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return d, nil
}

func (m *mockDeptRepo) List(_ context.Context) ([]smodel.Dept, error) {
	var list []smodel.Dept
	for _, d := range m.depts {
		list = append(list, *d)
	}
	return list, nil
}

func (m *mockDeptRepo) GetChildIDs(_ context.Context, deptID int64) ([]int64, error) {
	var ids []int64
	for _, d := range m.depts {
		if d.ParentID == deptID {
			ids = append(ids, d.ID)
		}
	}
	return ids, nil
}

func TestDeptService_Create_Root(t *testing.T) {
	testutil.NewTestDB(t)
	repo := newMockDeptRepo()
	svc := NewDeptService(repo)
	ctx := context.Background()

	id, err := svc.Create(ctx, &smodel.DeptCreateReq{
		DeptName: "总公司", ParentID: 0, Sort: 1, Status: 0,
	}, 1)
	require.NoError(t, err)
	assert.NotZero(t, id)
	assert.Equal(t, "0", repo.depts[id].Ancestors)
}

func TestDeptService_Create_Child(t *testing.T) {
	testutil.NewTestDB(t)
	repo := newMockDeptRepo()
	svc := NewDeptService(repo)
	ctx := context.Background()

	parent := &smodel.Dept{DeptName: "总公司", Ancestors: "0"}
	parent.ID = 9001
	repo.depts[parent.ID] = parent

	id, err := svc.Create(ctx, &smodel.DeptCreateReq{
		DeptName: "技术部", ParentID: parent.ID,
	}, 1)
	require.NoError(t, err)
	assert.Equal(t, "0,9001", repo.depts[id].Ancestors)
}

func TestDeptService_Create_InvalidParent(t *testing.T) {
	testutil.NewTestDB(t)
	repo := newMockDeptRepo()
	svc := NewDeptService(repo)
	ctx := context.Background()

	_, err := svc.Create(ctx, &smodel.DeptCreateReq{
		DeptName: "子部门", ParentID: 9999,
	}, 1)
	assert.Error(t, err)
}

func TestDeptService_Update(t *testing.T) {
	testutil.NewTestDB(t)
	repo := newMockDeptRepo()
	svc := NewDeptService(repo)
	ctx := context.Background()

	dept := &smodel.Dept{DeptName: "旧部门", Ancestors: "0"}
	dept.ID = 9002
	repo.depts[dept.ID] = dept

	err := svc.Update(ctx, &smodel.DeptUpdateReq{
		ID: dept.ID, DeptName: "新部门", ParentID: 0,
	}, 1)
	require.NoError(t, err)
	assert.Equal(t, "新部门", repo.depts[dept.ID].DeptName)
}

func TestDeptService_Delete_HasChildren(t *testing.T) {
	testutil.NewTestDB(t)
	repo := newMockDeptRepo()
	svc := NewDeptService(repo)
	ctx := context.Background()

	parent := &smodel.Dept{DeptName: "父"}
	parent.ID = 9003
	repo.depts[parent.ID] = parent

	child1 := &smodel.Dept{DeptName: "子1", ParentID: parent.ID}
	child1.ID = 9004
	repo.depts[child1.ID] = child1

	child2 := &smodel.Dept{DeptName: "子2", ParentID: parent.ID}
	child2.ID = 9005
	repo.depts[child2.ID] = child2

	err := svc.Delete(ctx, parent.ID)
	assert.Error(t, err)
}

func TestDeptService_Tree(t *testing.T) {
	testutil.NewTestDB(t)
	repo := newMockDeptRepo()
	svc := NewDeptService(repo)
	ctx := context.Background()

	p := &smodel.Dept{DeptName: "总公司", ParentID: 0}
	p.ID = 9101
	repo.depts[p.ID] = p

	c := &smodel.Dept{DeptName: "技术部", ParentID: p.ID}
	c.ID = 9102
	repo.depts[c.ID] = c

	tree, err := svc.Tree(ctx)
	require.NoError(t, err)
	require.Len(t, tree, 1)
	assert.Equal(t, "总公司", tree[0].DeptName)
	require.Len(t, tree[0].Children, 1)
}
