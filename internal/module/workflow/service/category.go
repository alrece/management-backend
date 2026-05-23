package service

import (
	"context"

	wfmodel "management-backend/internal/module/workflow/model"

	"management-backend/internal/module/workflow/repository"
	"management-backend/pkg/errcode"
	"management-backend/pkg/snowflake"
)

// CategoryService 分类业务接口
type CategoryService interface {
	Create(ctx context.Context, req *wfmodel.CategoryCreateReq, creator int64, tenantID int64) (int64, error)
	Update(ctx context.Context, req *wfmodel.CategoryUpdateReq, updater int64) error
	Delete(ctx context.Context, id int64) error
	Tree(ctx context.Context) ([]*wfmodel.CategoryResp, error)
}

type categoryService struct {
	repo repository.CategoryRepo
}

func NewCategoryService(repo repository.CategoryRepo) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) Create(ctx context.Context, req *wfmodel.CategoryCreateReq, creator int64, tenantID int64) (int64, error) {
	cat := &wfmodel.Category{
		Name:     req.Name,
		ParentID: req.ParentID,
	}
	cat.ID = snowflake.NextID()
	cat.TenantID = tenantID
	cat.Creator = creator
	cat.Updater = creator
	cat.Sort = req.Sort
	cat.Status = req.Status
	if cat.Status == 0 {
		cat.Status = 0
	}

	if err := s.repo.Create(ctx, cat); err != nil {
		return 0, err
	}
	return cat.ID, nil
}

func (s *categoryService) Update(ctx context.Context, req *wfmodel.CategoryUpdateReq, updater int64) error {
	existing, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return errcode.Err(4001) // workflow 模块错误码
	}

	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.ParentID != nil {
		existing.ParentID = *req.ParentID
	}
	existing.Sort = req.Sort
	if req.Status != nil {
		existing.Status = *req.Status
	}
	existing.Updater = updater

	return s.repo.Update(ctx, existing)
}

func (s *categoryService) Delete(ctx context.Context, id int64) error {
	// 检查是否有子分类
	children, err := s.repo.ListByParentID(ctx, id)
	if err != nil {
		return err
	}
	if len(children) > 0 {
		return errcode.WithMsg(4001, "存在子分类，无法删除")
	}
	return s.repo.Delete(ctx, id)
}

// Tree 构建分类树
func (s *categoryService) Tree(ctx context.Context) ([]*wfmodel.CategoryResp, error) {
	all, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	// 转为响应对象
	respMap := make(map[int64]*wfmodel.CategoryResp, len(all))
	for i := range all {
		respMap[all[i].ID] = toCategoryResp(&all[i])
	}

	// 构建树
	var roots []*wfmodel.CategoryResp
	for i := range all {
		node := respMap[all[i].ID]
		if all[i].ParentID == 0 {
			roots = append(roots, node)
		} else if parent, ok := respMap[all[i].ParentID]; ok {
			parent.Children = append(parent.Children, node)
		}
	}

	if roots == nil {
		roots = []*wfmodel.CategoryResp{}
	}
	return roots, nil
}

func toCategoryResp(cat *wfmodel.Category) *wfmodel.CategoryResp {
	return &wfmodel.CategoryResp{
		ID:         cat.ID,
		Name:       cat.Name,
		ParentID:   cat.ParentID,
		Sort:       cat.Sort,
		Status:     cat.Status,
		CreateTime: cat.CreateTime.Format("2006-01-02 15:04:05"),
	}
}
