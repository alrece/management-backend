package repository

import (
	"context"

	"management-backend/services/job-service/internal/model"

	"gorm.io/gorm"
)

type ExecLogRepo interface {
	Create(ctx context.Context, log *model.JobExecutionLog) error
	Update(ctx context.Context, id int64, updates map[string]interface{}) error
	Page(ctx context.Context, req *model.ExecLogPageReq, tenantID int64) ([]model.JobExecutionLog, int64, error)
}

type execLogRepo struct {
	db *gorm.DB
}

func NewExecLogRepo(db *gorm.DB) ExecLogRepo {
	return &execLogRepo{db: db}
}

func (r *execLogRepo) Create(ctx context.Context, log *model.JobExecutionLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *execLogRepo) Update(ctx context.Context, id int64, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.JobExecutionLog{}).Where("id = ?", id).Updates(updates).Error
}

func (r *execLogRepo) Page(ctx context.Context, req *model.ExecLogPageReq, tenantID int64) ([]model.JobExecutionLog, int64, error) {
	var list []model.JobExecutionLog
	var total int64
	db := r.db.WithContext(ctx).Model(&model.JobExecutionLog{}).Where("tenant_id = ?", tenantID)

	if req.TaskID != nil {
		db = db.Where("task_id = ?", *req.TaskID)
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.StartTime != "" {
		db = db.Where("start_time >= ?", req.StartTime)
	}
	if req.EndTime != "" {
		db = db.Where("start_time <= ?", req.EndTime)
	}

	db.Count(&total)
	db.Offset(req.Offset()).Limit(req.Limit()).Order("id DESC").Find(&list)
	return list, total, nil
}
