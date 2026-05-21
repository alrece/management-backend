package repository

import (
	"context"

	smodel "management-backend/internal/module/system/model"
	pkgmiddleware "management-backend/pkg/middleware"

	"gorm.io/gorm"
)

// DeptRepo 部门数据访问接口
type DeptRepo interface {
	Create(ctx context.Context, dept *smodel.Dept) error
	Update(ctx context.Context, dept *smodel.Dept) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*smodel.Dept, error)
	List(ctx context.Context) ([]smodel.Dept, error)
	GetChildIDs(ctx context.Context, deptID int64) ([]int64, error)
}

type deptRepo struct {
	db *gorm.DB
}

// NewDeptRepo 创建部门 Repository
func NewDeptRepo(db *gorm.DB) DeptRepo {
	return &deptRepo{db: db}
}

func (r *deptRepo) getDB(ctx context.Context) *gorm.DB {
	if tenantDB := pkgmiddleware.GetTenantDB(ctx); tenantDB != nil {
		return tenantDB
	}
	return r.db
}

func (r *deptRepo) Create(ctx context.Context, dept *smodel.Dept) error {
	return r.getDB(ctx).WithContext(ctx).Create(dept).Error
}

func (r *deptRepo) Update(ctx context.Context, dept *smodel.Dept) error {
	return r.getDB(ctx).WithContext(ctx).Save(dept).Error
}

func (r *deptRepo) Delete(ctx context.Context, id int64) error {
	return r.getDB(ctx).WithContext(ctx).Delete(&smodel.Dept{}, id).Error
}

func (r *deptRepo) GetByID(ctx context.Context, id int64) (*smodel.Dept, error) {
	var dept smodel.Dept
	err := r.getDB(ctx).WithContext(ctx).First(&dept, id).Error
	if err != nil {
		return nil, err
	}
	return &dept, nil
}

func (r *deptRepo) List(ctx context.Context) ([]smodel.Dept, error) {
	var list []smodel.Dept
	err := r.getDB(ctx).WithContext(ctx).Order("sort ASC, id ASC").Find(&list).Error
	return list, err
}

// GetChildIDs 获取所有子部门 ID
func (r *deptRepo) GetChildIDs(ctx context.Context, deptID int64) ([]int64, error) {
	var ids []int64
	err := r.getDB(ctx).WithContext(ctx).
		Model(&smodel.Dept{}).
		Where("id = ? OR FIND_IN_SET(?, ancestors) > 0", deptID, deptID).
		Pluck("id", &ids).Error
	return ids, err
}
