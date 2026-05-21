package repository

import (
	"context"

	"management-backend/services/workflow-service/internal/model"

	"gorm.io/gorm"
)

type WorkflowRepo interface {
	Create(ctx context.Context, w *model.WfWorkflow) error
	Update(ctx context.Context, w *model.WfWorkflow) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*model.WfWorkflow, error)
	Page(ctx context.Context, req *model.WorkflowPageReq, tenantID int64) ([]model.WfWorkflow, int64, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
	ListByStatus(ctx context.Context, status string) ([]model.WfWorkflow, error)
}

type workflowRepo struct{ db *gorm.DB }

func NewWorkflowRepo(db *gorm.DB) WorkflowRepo {
	return &workflowRepo{db: db}
}

func (r *workflowRepo) Create(ctx context.Context, w *model.WfWorkflow) error {
	return r.db.WithContext(ctx).Create(w).Error
}

func (r *workflowRepo) Update(ctx context.Context, w *model.WfWorkflow) error {
	return r.db.WithContext(ctx).Where("id = ? AND deleted = 0", w.ID).Updates(w).Error
}

func (r *workflowRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&model.WfWorkflow{}).Where("id = ?", id).Update("deleted", 1).Error
}

func (r *workflowRepo) GetByID(ctx context.Context, id int64) (*model.WfWorkflow, error) {
	var w model.WfWorkflow
	err := r.db.WithContext(ctx).Where("id = ? AND deleted = 0", id).First(&w).Error
	return &w, err
}

func (r *workflowRepo) Page(ctx context.Context, req *model.WorkflowPageReq, tenantID int64) ([]model.WfWorkflow, int64, error) {
	var list []model.WfWorkflow
	var total int64
	db := r.db.WithContext(ctx).Model(&model.WfWorkflow{}).Where("tenant_id = ? AND deleted = 0", tenantID)

	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.CategoryID > 0 {
		db = db.Where("category_id = ?", req.CategoryID)
	}
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}

	db.Count(&total)
	err := db.Offset(req.Offset()).Limit(req.Limit()).Order("id DESC").Find(&list).Error
	return list, total, err
}

func (r *workflowRepo) UpdateStatus(ctx context.Context, id int64, status string) error {
	return r.db.WithContext(ctx).Model(&model.WfWorkflow{}).Where("id = ? AND deleted = 0", id).Update("status", status).Error
}

func (r *workflowRepo) ListByStatus(ctx context.Context, status string) ([]model.WfWorkflow, error) {
	var list []model.WfWorkflow
	err := r.db.WithContext(ctx).Where("status = ? AND deleted = 0", status).Find(&list).Error
	return list, err
}
