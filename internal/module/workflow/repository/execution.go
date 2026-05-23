package repository

import (
	"context"

	wfmodel "management-backend/internal/module/workflow/model"

	pkgmiddleware "management-backend/pkg/middleware"

	"gorm.io/gorm"
)

// ExecutionRepo 执行实例数据访问接口
type ExecutionRepo interface {
	Create(ctx context.Context, exec *wfmodel.Execution) error
	Update(ctx context.Context, exec *wfmodel.Execution) error
	GetByID(ctx context.Context, id int64) (*wfmodel.Execution, error)
	Page(ctx context.Context, req *wfmodel.ExecutionPageReq) ([]wfmodel.Execution, int64, error)
	UpdateStatus(ctx context.Context, id int64, status string, errorMsg string) error
}

type executionRepo struct {
	db *gorm.DB
}

func NewExecutionRepo(db *gorm.DB) ExecutionRepo {
	return &executionRepo{db: db}
}

func (r *executionRepo) getDB(ctx context.Context) *gorm.DB {
	if tenantDB := pkgmiddleware.GetTenantDB(ctx); tenantDB != nil {
		return tenantDB
	}
	return r.db
}

func (r *executionRepo) Create(ctx context.Context, exec *wfmodel.Execution) error {
	return r.getDB(ctx).WithContext(ctx).Create(exec).Error
}

func (r *executionRepo) Update(ctx context.Context, exec *wfmodel.Execution) error {
	return r.getDB(ctx).WithContext(ctx).Save(exec).Error
}

func (r *executionRepo) GetByID(ctx context.Context, id int64) (*wfmodel.Execution, error) {
	var exec wfmodel.Execution
	if err := r.getDB(ctx).WithContext(ctx).First(&exec, id).Error; err != nil {
		return nil, err
	}
	return &exec, nil
}

func (r *executionRepo) Page(ctx context.Context, req *wfmodel.ExecutionPageReq) ([]wfmodel.Execution, int64, error) {
	var list []wfmodel.Execution
	var total int64

	db := r.getDB(ctx).WithContext(ctx).Model(&wfmodel.Execution{})

	if req.WorkflowID != nil {
		db = db.Where("workflow_id = ?", *req.WorkflowID)
	}
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Offset(req.Offset()).Limit(req.Limit()).Order("id DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *executionRepo) UpdateStatus(ctx context.Context, id int64, status string, errorMsg string) error {
	updates := map[string]any{"status": status}
	if errorMsg != "" {
		updates["error_msg"] = errorMsg
	}
	return r.getDB(ctx).WithContext(ctx).Model(&wfmodel.Execution{}).
		Where("id = ? AND status = ?", id, "running").Updates(updates).Error
}
