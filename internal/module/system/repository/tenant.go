package repository

import (
	"context"

	smodel "management-backend/internal/module/system/model"

	"gorm.io/gorm"
)

// TenantRepo 租户数据访问接口
type TenantRepo interface {
	Create(ctx context.Context, tenant *smodel.Tenant) error
	Update(ctx context.Context, tenant *smodel.Tenant) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*smodel.Tenant, error)
	Page(ctx context.Context, req *smodel.TenantPageReq) ([]smodel.Tenant, int64, error)
	UpdateInitStatus(ctx context.Context, id int64, oldStatus, newStatus int) (bool, error)
}

type tenantRepo struct {
	db *gorm.DB
}

// NewTenantRepo 创建租户 Repository
func NewTenantRepo(db *gorm.DB) TenantRepo {
	return &tenantRepo{db: db}
}

func (r *tenantRepo) Create(ctx context.Context, tenant *smodel.Tenant) error {
	return r.db.WithContext(ctx).Create(tenant).Error
}

func (r *tenantRepo) Update(ctx context.Context, tenant *smodel.Tenant) error {
	return r.db.WithContext(ctx).Select("*").Omit("deleted_at").Updates(tenant).Error
}

func (r *tenantRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&smodel.Tenant{}, id).Error
}

func (r *tenantRepo) GetByID(ctx context.Context, id int64) (*smodel.Tenant, error) {
	var tenant smodel.Tenant
	if err := r.db.WithContext(ctx).First(&tenant, id).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *tenantRepo) Page(ctx context.Context, req *smodel.TenantPageReq) ([]smodel.Tenant, int64, error) {
	var list []smodel.Tenant
	var total int64

	db := r.db.WithContext(ctx).Model(&smodel.Tenant{})

	if req.TenantName != "" {
		db = db.Where("tenant_name LIKE ?", "%"+req.TenantName+"%")
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

// UpdateInitStatus CAS 更新初始化状态（乐观锁）
func (r *tenantRepo) UpdateInitStatus(ctx context.Context, id int64, oldStatus, newStatus int) (bool, error) {
	result := r.db.WithContext(ctx).
		Model(&smodel.Tenant{}).
		Where("id = ? AND init_status = ?", id, oldStatus).
		Update("init_status", newStatus)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}
