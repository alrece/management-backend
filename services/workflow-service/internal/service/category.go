package service

import (
	"context"
	"time"

	"management-backend/services/workflow-service/internal/model"
	"management-backend/services/workflow-service/internal/repository"
	"management-backend/pkg/snowflake"
)

type CategoryService interface {
	Create(ctx context.Context, req *model.CategoryCreateReq, creator, tenantID int64) (int64, error)
	Update(ctx context.Context, id int64, req *model.CategoryUpdateReq) error
	Delete(ctx context.Context, id int64) error
	Tree(ctx context.Context, tenantID int64) ([]*model.CategoryTreeResp, error)
}

type categoryService struct {
	repo repository.CategoryRepo
}

func NewCategoryService(repo repository.CategoryRepo) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) Create(ctx context.Context, req *model.CategoryCreateReq, creator, tenantID int64) (int64, error) {
	c := &model.WfCategory{
		ID:       snowflake.NextID(),
		TenantID: tenantID,
		Name:     req.Name,
		ParentID: req.ParentID,
		Sort:     req.Sort,
		Status:   0,
		Creator:  creator,
		Updater:  creator,
	}
	if err := s.repo.Create(ctx, c); err != nil {
		return 0, err
	}
	return c.ID, nil
}

func (s *categoryService) Update(ctx context.Context, id int64, req *model.CategoryUpdateReq) error {
	c := &model.WfCategory{
		Name:     req.Name,
		ParentID: req.ParentID,
		Sort:     req.Sort,
	}
	c.ID = id
	c.Updater = 0
	_ = time.Now()
	return s.repo.Update(ctx, c)
}

func (s *categoryService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *categoryService) Tree(ctx context.Context, tenantID int64) ([]*model.CategoryTreeResp, error) {
	list, err := s.repo.ListByTenantID(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	return buildCategoryTree(list, 0), nil
}

func buildCategoryTree(list []model.WfCategory, parentID int64) []*model.CategoryTreeResp {
	var tree []*model.CategoryTreeResp
	for i := range list {
		if list[i].ParentID == parentID {
			node := &model.CategoryTreeResp{
				ID:       list[i].ID,
				Name:     list[i].Name,
				ParentID: list[i].ParentID,
				Sort:     list[i].Sort,
			}
			node.Children = buildCategoryTree(list, list[i].ID)
			tree = append(tree, node)
		}
	}
	return tree
}
