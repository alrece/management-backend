package service

import (
	"context"
	"fmt"

	smodel "management-backend/internal/module/system/model"
	"management-backend/internal/module/system/repository"
	"management-backend/pkg/errcode"
	"management-backend/pkg/snowflake"
)

// DeptService 部门业务接口
type DeptService interface {
	Create(ctx context.Context, req *smodel.DeptCreateReq, creator int64) (int64, error)
	Update(ctx context.Context, req *smodel.DeptUpdateReq, updater int64) error
	Delete(ctx context.Context, id int64) error
	Tree(ctx context.Context) ([]smodel.DeptTreeResp, error)
	GetChildDeptIDs(ctx context.Context, deptID int64) ([]int64, error)
}

type deptService struct {
	repo repository.DeptRepo
}

// NewDeptService 创建部门 Service
func NewDeptService(repo repository.DeptRepo) DeptService {
	return &deptService{repo: repo}
}

func (s *deptService) Create(ctx context.Context, req *smodel.DeptCreateReq, creator int64) (int64, error) {
	ancestors := "0"
	if req.ParentID > 0 {
		parent, err := s.repo.GetByID(ctx, req.ParentID)
		if err != nil {
			return 0, errcode.Err(errcode.DeptNotFound)
		}
		ancestors = fmt.Sprintf("%s,%d", parent.Ancestors, parent.ID)
	}

	dept := &smodel.Dept{
		ParentID:  req.ParentID,
		Ancestors: ancestors,
		DeptName:  req.DeptName,
		Leader:    req.Leader,
	}
	dept.ID = snowflake.NextID()
	dept.Creator = creator
	dept.Updater = creator
	dept.SortModel.Sort = req.Sort
	dept.StatusModel.Status = req.Status

	if err := s.repo.Create(ctx, dept); err != nil {
		return 0, err
	}
	return dept.ID, nil
}

func (s *deptService) Update(ctx context.Context, req *smodel.DeptUpdateReq, updater int64) error {
	dept, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return errcode.Err(errcode.DeptNotFound)
	}

	if req.ParentID != dept.ParentID {
		ancestors := "0"
		if req.ParentID > 0 {
			parent, err := s.repo.GetByID(ctx, req.ParentID)
			if err != nil {
				return errcode.Err(errcode.DeptNotFound)
			}
			ancestors = fmt.Sprintf("%s,%d", parent.Ancestors, parent.ID)
		}
		dept.Ancestors = ancestors
	}

	dept.ParentID = req.ParentID
	dept.DeptName = req.DeptName
	dept.Leader = req.Leader
	dept.SortModel.Sort = req.Sort
	dept.StatusModel.Status = req.Status
	dept.Updater = updater

	return s.repo.Update(ctx, dept)
}

func (s *deptService) Delete(ctx context.Context, id int64) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return errcode.Err(errcode.DeptNotFound)
	}

	childIDs, _ := s.repo.GetChildIDs(ctx, id)
	if len(childIDs) > 1 {
		return errcode.Err(errcode.DeptHasChildren)
	}

	return s.repo.Delete(ctx, id)
}

func (s *deptService) Tree(ctx context.Context) ([]smodel.DeptTreeResp, error) {
	list, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	return buildDeptTree(list, 0), nil
}

func (s *deptService) GetChildDeptIDs(ctx context.Context, deptID int64) ([]int64, error) {
	return s.repo.GetChildIDs(ctx, deptID)
}

func buildDeptTree(depts []smodel.Dept, parentID int64) []smodel.DeptTreeResp {
	var tree []smodel.DeptTreeResp
	for _, d := range depts {
		if d.ParentID == parentID {
			node := smodel.DeptTreeResp{
				ID: d.ID, ParentID: d.ParentID, DeptName: d.DeptName,
				Ancestors: d.Ancestors, Sort: d.Sort, Leader: d.Leader,
				Status: d.Status, Children: buildDeptTree(depts, d.ID),
			}
			tree = append(tree, node)
		}
	}
	return tree
}
