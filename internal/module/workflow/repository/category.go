package repository

import (
	"context"

	wfmodel "management-backend/internal/module/workflow/model"

	pkgmiddleware "management-backend/pkg/middleware"

	"gorm.io/gorm"
)

// CategoryRepo 分类数据访问接口
type CategoryRepo interface {
	Create(ctx context.Context, cat *wfmodel.Category) error
	Update(ctx context.Context, cat *wfmodel.Category) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*wfmodel.Category, error)
	ListByParentID(ctx context.Context, parentID int64) ([]wfmodel.Category, error)
	ListAll(ctx context.Context) ([]wfmodel.Category, error)
}

type categoryRepo struct {
	db *gorm.DB
}

func NewCategoryRepo(db *gorm.DB) CategoryRepo {
	return &categoryRepo{db: db}
}

func (r *categoryRepo) getDB(ctx context.Context) *gorm.DB {
	if tenantDB := pkgmiddleware.GetTenantDB(ctx); tenantDB != nil {
		return tenantDB
	}
	return r.db
}

func (r *categoryRepo) Create(ctx context.Context, cat *wfmodel.Category) error {
	return r.getDB(ctx).WithContext(ctx).Create(cat).Error
}

func (r *categoryRepo) Update(ctx context.Context, cat *wfmodel.Category) error {
	return r.getDB(ctx).WithContext(ctx).Select("*").Omit("TenantID").Updates(cat).Error
}

func (r *categoryRepo) Delete(ctx context.Context, id int64) error {
	return r.getDB(ctx).WithContext(ctx).Delete(&wfmodel.Category{}, id).Error
}

func (r *categoryRepo) GetByID(ctx context.Context, id int64) (*wfmodel.Category, error) {
	var cat wfmodel.Category
	if err := r.getDB(ctx).WithContext(ctx).First(&cat, id).Error; err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *categoryRepo) ListByParentID(ctx context.Context, parentID int64) ([]wfmodel.Category, error) {
	var list []wfmodel.Category
	err := r.getDB(ctx).WithContext(ctx).
		Where("parent_id = ?", parentID).
		Order("sort ASC, id ASC").Find(&list).Error
	return list, err
}

func (r *categoryRepo) ListAll(ctx context.Context) ([]wfmodel.Category, error) {
	var list []wfmodel.Category
	err := r.getDB(ctx).WithContext(ctx).
		Order("sort ASC, id ASC").Find(&list).Error
	return list, err
}
