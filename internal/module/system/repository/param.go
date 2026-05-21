package repository

import (
	"context"

	pkgmiddleware "management-backend/pkg/middleware"

	smodel "management-backend/internal/module/system/model"

	"gorm.io/gorm"
)

// ParamRepo 参数数据访问接口
type ParamRepo interface {
	Create(ctx context.Context, param *smodel.SysParam) error
	Update(ctx context.Context, param *smodel.SysParam) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*smodel.SysParam, error)
	GetByKey(ctx context.Context, key string) (*smodel.SysParam, error)
	Page(ctx context.Context, req *smodel.ParamPageReq) ([]smodel.SysParam, int64, error)
}

type paramRepo struct {
	db *gorm.DB
}

func NewParamRepo(db *gorm.DB) ParamRepo {
	return &paramRepo{db: db}
}

func (r *paramRepo) getDB(ctx context.Context) *gorm.DB {
	if tenantDB := pkgmiddleware.GetTenantDB(ctx); tenantDB != nil {
		return tenantDB
	}
	return r.db
}

func (r *paramRepo) Create(ctx context.Context, param *smodel.SysParam) error {
	return r.getDB(ctx).WithContext(ctx).Create(param).Error
}

func (r *paramRepo) Update(ctx context.Context, param *smodel.SysParam) error {
	return r.getDB(ctx).WithContext(ctx).Save(param).Error
}

func (r *paramRepo) Delete(ctx context.Context, id int64) error {
	return r.getDB(ctx).WithContext(ctx).Delete(&smodel.SysParam{}, id).Error
}

func (r *paramRepo) GetByID(ctx context.Context, id int64) (*smodel.SysParam, error) {
	var param smodel.SysParam
	if err := r.getDB(ctx).WithContext(ctx).First(&param, id).Error; err != nil {
		return nil, err
	}
	return &param, nil
}

func (r *paramRepo) GetByKey(ctx context.Context, key string) (*smodel.SysParam, error) {
	var param smodel.SysParam
	if err := r.getDB(ctx).WithContext(ctx).Where("param_key = ?", key).First(&param).Error; err != nil {
		return nil, err
	}
	return &param, nil
}

func (r *paramRepo) Page(ctx context.Context, req *smodel.ParamPageReq) ([]smodel.SysParam, int64, error) {
	var list []smodel.SysParam
	var total int64

	db := r.getDB(ctx).WithContext(ctx).Model(&smodel.SysParam{})

	if req.ParamKey != "" {
		db = db.Where("param_key LIKE ?", "%"+req.ParamKey+"%")
	}
	if req.ParamType != nil {
		db = db.Where("param_type = ?", *req.ParamType)
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
