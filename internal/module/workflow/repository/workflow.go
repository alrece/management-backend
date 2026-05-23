package repository

import (
	"context"
	"encoding/json"

	wfmodel "management-backend/internal/module/workflow/model"

	pkgmiddleware "management-backend/pkg/middleware"

	"gorm.io/gorm"
)

// WorkflowRepo 工作流数据访问接口
type WorkflowRepo interface {
	Create(ctx context.Context, wf *wfmodel.Workflow) error
	Update(ctx context.Context, wf *wfmodel.Workflow) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*wfmodel.Workflow, error)
	Page(ctx context.Context, req *wfmodel.WorkflowPageReq) ([]wfmodel.Workflow, int64, error)
	UpdateStatus(ctx context.Context, id int64, status int) error
	UpdateDefinition(ctx context.Context, id int64, def json.RawMessage, version int) error
	GetActiveByTriggerType(ctx context.Context, triggerType string) ([]wfmodel.Workflow, error)
}

type workflowRepo struct {
	db *gorm.DB
}

func NewWorkflowRepo(db *gorm.DB) WorkflowRepo {
	return &workflowRepo{db: db}
}

func (r *workflowRepo) getDB(ctx context.Context) *gorm.DB {
	if tenantDB := pkgmiddleware.GetTenantDB(ctx); tenantDB != nil {
		return tenantDB
	}
	return r.db
}

func (r *workflowRepo) Create(ctx context.Context, wf *wfmodel.Workflow) error {
	return r.getDB(ctx).WithContext(ctx).Create(wf).Error
}

func (r *workflowRepo) Update(ctx context.Context, wf *wfmodel.Workflow) error {
	return r.getDB(ctx).WithContext(ctx).Select("*").Omit("TenantID").Updates(wf).Error
}

func (r *workflowRepo) Delete(ctx context.Context, id int64) error {
	return r.getDB(ctx).WithContext(ctx).Delete(&wfmodel.Workflow{}, id).Error
}

func (r *workflowRepo) GetByID(ctx context.Context, id int64) (*wfmodel.Workflow, error) {
	var wf wfmodel.Workflow
	if err := r.getDB(ctx).WithContext(ctx).First(&wf, id).Error; err != nil {
		return nil, err
	}
	return &wf, nil
}

func (r *workflowRepo) Page(ctx context.Context, req *wfmodel.WorkflowPageReq) ([]wfmodel.Workflow, int64, error) {
	var list []wfmodel.Workflow
	var total int64

	db := r.getDB(ctx).WithContext(ctx).Model(&wfmodel.Workflow{})

	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.CategoryID != nil {
		db = db.Where("category_id = ?", *req.CategoryID)
	}
	if req.TriggerType != "" {
		db = db.Where("trigger_type = ?", req.TriggerType)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Offset(req.Offset()).Limit(req.Limit()).Order("id DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *workflowRepo) UpdateStatus(ctx context.Context, id int64, status int) error {
	return r.getDB(ctx).WithContext(ctx).Model(&wfmodel.Workflow{}).
		Where("id = ?", id).Update("status", status).Error
}

func (r *workflowRepo) UpdateDefinition(ctx context.Context, id int64, def json.RawMessage, version int) error {
	return r.getDB(ctx).WithContext(ctx).Model(&wfmodel.Workflow{}).
		Where("id = ?", id).Updates(map[string]any{
		"definition": def,
		"version":    version,
	}).Error
}

func (r *workflowRepo) GetActiveByTriggerType(ctx context.Context, triggerType string) ([]wfmodel.Workflow, error) {
	var list []wfmodel.Workflow
	err := r.getDB(ctx).WithContext(ctx).
		Where("trigger_type = ? AND status = 0", triggerType).
		Find(&list).Error
	return list, err
}
