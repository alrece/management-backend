package repository

import (
	"context"

	wfmodel "management-backend/internal/module/workflow/model"

	pkgmiddleware "management-backend/pkg/middleware"

	"gorm.io/gorm"
)

// TriggerRepo 触发器配置数据访问接口
type TriggerRepo interface {
	Create(ctx context.Context, trigger *wfmodel.Trigger) error
	Update(ctx context.Context, trigger *wfmodel.Trigger) error
	GetByWorkflowID(ctx context.Context, workflowID int64) (*wfmodel.Trigger, error)
	DeleteByWorkflowID(ctx context.Context, workflowID int64) error
}

type triggerRepo struct {
	db *gorm.DB
}

func NewTriggerRepo(db *gorm.DB) TriggerRepo {
	return &triggerRepo{db: db}
}

func (r *triggerRepo) getDB(ctx context.Context) *gorm.DB {
	if tenantDB := pkgmiddleware.GetTenantDB(ctx); tenantDB != nil {
		return tenantDB
	}
	return r.db
}

func (r *triggerRepo) Create(ctx context.Context, trigger *wfmodel.Trigger) error {
	return r.getDB(ctx).WithContext(ctx).Create(trigger).Error
}

func (r *triggerRepo) Update(ctx context.Context, trigger *wfmodel.Trigger) error {
	return r.getDB(ctx).WithContext(ctx).Save(trigger).Error
}

func (r *triggerRepo) GetByWorkflowID(ctx context.Context, workflowID int64) (*wfmodel.Trigger, error) {
	var trigger wfmodel.Trigger
	if err := r.getDB(ctx).WithContext(ctx).
		Where("workflow_id = ?", workflowID).First(&trigger).Error; err != nil {
		return nil, err
	}
	return &trigger, nil
}

func (r *triggerRepo) DeleteByWorkflowID(ctx context.Context, workflowID int64) error {
	return r.getDB(ctx).WithContext(ctx).
		Where("workflow_id = ?", workflowID).Delete(&wfmodel.Trigger{}).Error
}
