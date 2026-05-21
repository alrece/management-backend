package repository

import (
	"context"

	pkgmiddleware "management-backend/pkg/middleware"

	smodel "management-backend/internal/module/system/model"

	"gorm.io/gorm"
)

// DictDataRepo 字典数据访问接口
type DictDataRepo interface {
	Create(ctx context.Context, data *smodel.DictData) error
	Update(ctx context.Context, data *smodel.DictData) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*smodel.DictData, error)
	ListByDictType(ctx context.Context, dictType string) ([]smodel.DictData, error)
	Page(ctx context.Context, req *smodel.DictDataPageReq) ([]smodel.DictData, int64, error)
}

type dictDataRepo struct {
	db *gorm.DB
}

func NewDictDataRepo(db *gorm.DB) DictDataRepo {
	return &dictDataRepo{db: db}
}

func (r *dictDataRepo) getDB(ctx context.Context) *gorm.DB {
	if tenantDB := pkgmiddleware.GetTenantDB(ctx); tenantDB != nil {
		return tenantDB
	}
	return r.db
}

func (r *dictDataRepo) Create(ctx context.Context, data *smodel.DictData) error {
	return r.getDB(ctx).WithContext(ctx).Create(data).Error
}

func (r *dictDataRepo) Update(ctx context.Context, data *smodel.DictData) error {
	return r.getDB(ctx).WithContext(ctx).Save(data).Error
}

func (r *dictDataRepo) Delete(ctx context.Context, id int64) error {
	return r.getDB(ctx).WithContext(ctx).Delete(&smodel.DictData{}, id).Error
}

func (r *dictDataRepo) GetByID(ctx context.Context, id int64) (*smodel.DictData, error) {
	var data smodel.DictData
	if err := r.getDB(ctx).WithContext(ctx).First(&data, id).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *dictDataRepo) ListByDictType(ctx context.Context, dictType string) ([]smodel.DictData, error) {
	var list []smodel.DictData
	err := r.getDB(ctx).WithContext(ctx).
		Where("dict_type = ? AND status = 0", dictType).
		Order("sort ASC, id ASC").
		Find(&list).Error
	return list, err
}

func (r *dictDataRepo) Page(ctx context.Context, req *smodel.DictDataPageReq) ([]smodel.DictData, int64, error) {
	var list []smodel.DictData
	var total int64

	db := r.getDB(ctx).WithContext(ctx).Model(&smodel.DictData{})

	if req.DictType != "" {
		db = db.Where("dict_type = ?", req.DictType)
	}
	if req.Label != "" {
		db = db.Where("label LIKE ?", "%"+req.Label+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Offset(req.Offset()).Limit(req.Limit()).Order("sort ASC, id DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
