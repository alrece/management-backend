package repository

import (
	"context"

	pkgmiddleware "management-backend/pkg/middleware"

	smodel "management-backend/internal/module/system/model"

	"gorm.io/gorm"
)

// DictTypeRepo 字典类型数据访问接口
type DictTypeRepo interface {
	Create(ctx context.Context, dict *smodel.DictType) error
	Update(ctx context.Context, dict *smodel.DictType) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*smodel.DictType, error)
	GetByDictType(ctx context.Context, dictType string) (*smodel.DictType, error)
	Page(ctx context.Context, req *smodel.DictTypePageReq) ([]smodel.DictType, int64, error)
}

type dictTypeRepo struct {
	db *gorm.DB
}

func NewDictTypeRepo(db *gorm.DB) DictTypeRepo {
	return &dictTypeRepo{db: db}
}

func (r *dictTypeRepo) getDB(ctx context.Context) *gorm.DB {
	if tenantDB := pkgmiddleware.GetTenantDB(ctx); tenantDB != nil {
		return tenantDB
	}
	return r.db
}

func (r *dictTypeRepo) Create(ctx context.Context, dict *smodel.DictType) error {
	return r.getDB(ctx).WithContext(ctx).Create(dict).Error
}

func (r *dictTypeRepo) Update(ctx context.Context, dict *smodel.DictType) error {
	return r.getDB(ctx).WithContext(ctx).Save(dict).Error
}

func (r *dictTypeRepo) Delete(ctx context.Context, id int64) error {
	return r.getDB(ctx).WithContext(ctx).Delete(&smodel.DictType{}, id).Error
}

func (r *dictTypeRepo) GetByID(ctx context.Context, id int64) (*smodel.DictType, error) {
	var dict smodel.DictType
	if err := r.getDB(ctx).WithContext(ctx).First(&dict, id).Error; err != nil {
		return nil, err
	}
	return &dict, nil
}

func (r *dictTypeRepo) GetByDictType(ctx context.Context, dictType string) (*smodel.DictType, error) {
	var dict smodel.DictType
	if err := r.getDB(ctx).WithContext(ctx).Where("dict_type = ?", dictType).First(&dict).Error; err != nil {
		return nil, err
	}
	return &dict, nil
}

func (r *dictTypeRepo) Page(ctx context.Context, req *smodel.DictTypePageReq) ([]smodel.DictType, int64, error) {
	var list []smodel.DictType
	var total int64

	db := r.getDB(ctx).WithContext(ctx).Model(&smodel.DictType{})

	if req.DictName != "" {
		db = db.Where("dict_name LIKE ?", "%"+req.DictName+"%")
	}
	if req.DictType != "" {
		db = db.Where("dict_type LIKE ?", "%"+req.DictType+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Offset(req.Offset()).Limit(req.Limit()).Order("id DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
