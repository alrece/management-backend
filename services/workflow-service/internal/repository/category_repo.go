package repository

import (
	"context"

	"management-backend/services/workflow-service/internal/model"

	"gorm.io/gorm"
)

type CategoryRepo interface {
	Create(ctx context.Context, c *model.WfCategory) error
	Update(ctx context.Context, c *model.WfCategory) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*model.WfCategory, error)
	ListByTenantID(ctx context.Context, tenantID int64) ([]model.WfCategory, error)
}

type categoryRepo struct{ db *gorm.DB }

func NewCategoryRepo(db *gorm.DB) CategoryRepo {
	return &categoryRepo{db: db}
}

func (r *categoryRepo) Create(ctx context.Context, c *model.WfCategory) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *categoryRepo) Update(ctx context.Context, c *model.WfCategory) error {
	return r.db.WithContext(ctx).Where("id = ? AND deleted = 0", c.ID).Updates(c).Error
}

func (r *categoryRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&model.WfCategory{}).Where("id = ?", id).Update("deleted", 1).Error
}

func (r *categoryRepo) GetByID(ctx context.Context, id int64) (*model.WfCategory, error) {
	var c model.WfCategory
	err := r.db.WithContext(ctx).Where("id = ? AND deleted = 0", id).First(&c).Error
	return &c, err
}

func (r *categoryRepo) ListByTenantID(ctx context.Context, tenantID int64) ([]model.WfCategory, error) {
	var list []model.WfCategory
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND deleted = 0", tenantID).Order("sort ASC, id ASC").Find(&list).Error
	return list, err
}
