package repository

import (
	"context"

	pkgmiddleware "management-backend/pkg/middleware"

	smodel "management-backend/internal/module/system/model"

	"gorm.io/gorm"
)

// FileRepo 文件数据访问接口
type FileRepo interface {
	Create(ctx context.Context, file *smodel.File) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*smodel.File, error)
	Page(ctx context.Context, req *smodel.FilePageReq) ([]smodel.File, int64, error)
}

type fileRepo struct {
	db *gorm.DB
}

// NewFileRepo 创建文件 Repository
func NewFileRepo(db *gorm.DB) FileRepo {
	return &fileRepo{db: db}
}

func (r *fileRepo) getDB(ctx context.Context) *gorm.DB {
	if tenantDB := pkgmiddleware.GetTenantDB(ctx); tenantDB != nil {
		return tenantDB
	}
	return r.db
}

func (r *fileRepo) Create(ctx context.Context, file *smodel.File) error {
	return r.getDB(ctx).WithContext(ctx).Create(file).Error
}

func (r *fileRepo) Delete(ctx context.Context, id int64) error {
	return r.getDB(ctx).WithContext(ctx).Delete(&smodel.File{}, id).Error
}

func (r *fileRepo) GetByID(ctx context.Context, id int64) (*smodel.File, error) {
	var file smodel.File
	if err := r.getDB(ctx).WithContext(ctx).First(&file, id).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *fileRepo) Page(ctx context.Context, req *smodel.FilePageReq) ([]smodel.File, int64, error) {
	var list []smodel.File
	var total int64
	db := r.getDB(ctx).WithContext(ctx).Model(&smodel.File{})

	if req.FileName != "" {
		db = db.Where("file_name LIKE ?", "%"+req.FileName+"%")
	}
	if req.FileType != "" {
		db = db.Where("file_type = ?", req.FileType)
	}
	if req.StorageType != "" {
		db = db.Where("storage_type = ?", req.StorageType)
	}
	if req.CreateTime != "" {
		db = db.Where("DATE(create_time) = ?", req.CreateTime)
	}

	db.Count(&total)
	db.Offset(req.Offset()).Limit(req.Limit()).Order("id DESC").Find(&list)
	return list, total, nil
}
