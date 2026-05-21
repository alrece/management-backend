package repository

import (
	"context"

	"management-backend/services/workflow-service/internal/model"

	"gorm.io/gorm"
)

type InstanceRepo interface {
	Create(ctx context.Context, inst *model.WfInstance) error
	GetByID(ctx context.Context, id int64) (*model.WfInstance, error)
	Page(ctx context.Context, req *model.InstancePageReq, tenantID int64) ([]model.WfInstance, int64, error)
	UpdateResult(ctx context.Context, id int64, status, result, errMsg string, durationMs int64) error
	ListRunningByTenant(ctx context.Context, tenantID int64) ([]model.WfInstance, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
}

type instanceRepo struct{ db *gorm.DB }

func NewInstanceRepo(db *gorm.DB) InstanceRepo {
	return &instanceRepo{db: db}
}

func (r *instanceRepo) Create(ctx context.Context, inst *model.WfInstance) error {
	return r.db.WithContext(ctx).Create(inst).Error
}

func (r *instanceRepo) GetByID(ctx context.Context, id int64) (*model.WfInstance, error) {
	var inst model.WfInstance
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&inst).Error
	return &inst, err
}

func (r *instanceRepo) Page(ctx context.Context, req *model.InstancePageReq, tenantID int64) ([]model.WfInstance, int64, error) {
	var list []model.WfInstance
	var total int64
	db := r.db.WithContext(ctx).Model(&model.WfInstance{}).Where("tenant_id = ?", tenantID)

	if req.WorkflowID > 0 {
		db = db.Where("workflow_id = ?", req.WorkflowID)
	}
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}

	db.Count(&total)
	err := db.Offset(req.Offset()).Limit(req.Limit()).Order("id DESC").Find(&list).Error
	return list, total, err
}

func (r *instanceRepo) UpdateResult(ctx context.Context, id int64, status, result, errMsg string, durationMs int64) error {
	return r.db.WithContext(ctx).Model(&model.WfInstance{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      status,
		"result":      result,
		"error_msg":   errMsg,
		"duration_ms": durationMs,
	}).Error
}

func (r *instanceRepo) ListRunningByTenant(ctx context.Context, tenantID int64) ([]model.WfInstance, error) {
	var list []model.WfInstance
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND status = ?", tenantID, "RUNNING").Find(&list).Error
	return list, err
}

func (r *instanceRepo) UpdateStatus(ctx context.Context, id int64, status string) error {
	return r.db.WithContext(ctx).Model(&model.WfInstance{}).Where("id = ?", id).Update("status", status).Error
}
