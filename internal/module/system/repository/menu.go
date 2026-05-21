package repository

import (
	"context"

	smodel "management-backend/internal/module/system/model"

	"gorm.io/gorm"
)

// MenuRepo 菜单数据访问接口（默认库）
type MenuRepo interface {
	Create(ctx context.Context, menu *smodel.Menu) error
	Update(ctx context.Context, menu *smodel.Menu) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*smodel.Menu, error)
	List(ctx context.Context) ([]smodel.Menu, error)
	GetByIDs(ctx context.Context, ids []int64) ([]smodel.Menu, error)
	HasChildren(ctx context.Context, id int64) (bool, error)
}

type menuRepo struct {
	db *gorm.DB
}

// NewMenuRepo 创建菜单 Repository（注入默认库）
func NewMenuRepo(db *gorm.DB) MenuRepo {
	return &menuRepo{db: db}
}

func (r *menuRepo) Create(ctx context.Context, menu *smodel.Menu) error {
	return r.db.WithContext(ctx).Create(menu).Error
}

func (r *menuRepo) Update(ctx context.Context, menu *smodel.Menu) error {
	return r.db.WithContext(ctx).Save(menu).Error
}

func (r *menuRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&smodel.Menu{}, id).Error
}

func (r *menuRepo) GetByID(ctx context.Context, id int64) (*smodel.Menu, error) {
	var menu smodel.Menu
	err := r.db.WithContext(ctx).First(&menu, id).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r *menuRepo) List(ctx context.Context) ([]smodel.Menu, error) {
	var list []smodel.Menu
	err := r.db.WithContext(ctx).Order("sort ASC, id ASC").Find(&list).Error
	return list, err
}

func (r *menuRepo) GetByIDs(ctx context.Context, ids []int64) ([]smodel.Menu, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var list []smodel.Menu
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Order("sort ASC, id ASC").Find(&list).Error
	return list, err
}

func (r *menuRepo) HasChildren(ctx context.Context, id int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&smodel.Menu{}).Where("parent_id = ?", id).Count(&count).Error
	return count > 0, err
}
