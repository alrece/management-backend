package repository

import (
	"context"

	smodel "management-backend/internal/module/system/model"

	"gorm.io/gorm"
)

// OperLogRepo 操作日志数据访问接口（使用默认库）
type OperLogRepo interface {
	Create(ctx context.Context, log *smodel.SysOperLog) error
	Page(ctx context.Context, req *smodel.OperLogPageReq) ([]smodel.SysOperLog, int64, error)
	DeleteByID(ctx context.Context, id int64) error
	Clean(ctx context.Context) error
}

type operLogRepo struct {
	db *gorm.DB
}

func NewOperLogRepo(db *gorm.DB) OperLogRepo {
	return &operLogRepo{db: db}
}

func (r *operLogRepo) Create(ctx context.Context, log *smodel.SysOperLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *operLogRepo) Page(ctx context.Context, req *smodel.OperLogPageReq) ([]smodel.SysOperLog, int64, error) {
	var list []smodel.SysOperLog
	var total int64

	db := r.db.WithContext(ctx).Model(&smodel.SysOperLog{})

	if req.Title != "" {
		db = db.Where("title LIKE ?", "%"+req.Title+"%")
	}
	if req.BusinessType != nil {
		db = db.Where("business_type = ?", *req.BusinessType)
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.OperName != "" {
		db = db.Where("oper_name LIKE ?", "%"+req.OperName+"%")
	}
	if req.BeginTime != "" {
		db = db.Where("oper_time >= ?", req.BeginTime)
	}
	if req.EndTime != "" {
		db = db.Where("oper_time <= ?", req.EndTime+" 23:59:59")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Offset(req.Offset()).Limit(req.Limit()).Order("id DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *operLogRepo) DeleteByID(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&smodel.SysOperLog{}, id).Error
}

func (r *operLogRepo) Clean(ctx context.Context) error {
	return r.db.WithContext(ctx).Exec("DELETE FROM sys_oper_log").Error
}

// LoginLogRepo 登录日志数据访问接口（使用默认库）
type LoginLogRepo interface {
	Create(ctx context.Context, log *smodel.SysLoginLog) error
	Page(ctx context.Context, req *smodel.LoginLogPageReq) ([]smodel.SysLoginLog, int64, error)
	DeleteByID(ctx context.Context, id int64) error
	Clean(ctx context.Context) error
}

type loginLogRepo struct {
	db *gorm.DB
}

func NewLoginLogRepo(db *gorm.DB) LoginLogRepo {
	return &loginLogRepo{db: db}
}

func (r *loginLogRepo) Create(ctx context.Context, log *smodel.SysLoginLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *loginLogRepo) Page(ctx context.Context, req *smodel.LoginLogPageReq) ([]smodel.SysLoginLog, int64, error) {
	var list []smodel.SysLoginLog
	var total int64

	db := r.db.WithContext(ctx).Model(&smodel.SysLoginLog{})

	if req.Username != "" {
		db = db.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.BeginTime != "" {
		db = db.Where("login_time >= ?", req.BeginTime)
	}
	if req.EndTime != "" {
		db = db.Where("login_time <= ?", req.EndTime+" 23:59:59")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Offset(req.Offset()).Limit(req.Limit()).Order("id DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *loginLogRepo) DeleteByID(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&smodel.SysLoginLog{}, id).Error
}

func (r *loginLogRepo) Clean(ctx context.Context) error {
	return r.db.WithContext(ctx).Exec("DELETE FROM sys_login_log").Error
}
