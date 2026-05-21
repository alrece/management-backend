package repository

import (
	"context"

	"management-backend/services/job-service/internal/model"

	"gorm.io/gorm"
)

type JobTaskRepo interface {
	Create(ctx context.Context, task *model.JobTask) error
	GetByID(ctx context.Context, id int64) (*model.JobTask, error)
	Update(ctx context.Context, id int64, updates map[string]interface{}) error
	Delete(ctx context.Context, id int64) error
	Page(ctx context.Context, req *model.JobTaskPageReq, tenantID int64) ([]model.JobTask, int64, error)
	ListActive(ctx context.Context, tenantID int64) ([]model.JobTask, error)
}

type jobTaskRepo struct {
	db *gorm.DB
}

func NewJobTaskRepo(db *gorm.DB) JobTaskRepo {
	return &jobTaskRepo{db: db}
}

func (r *jobTaskRepo) Create(ctx context.Context, task *model.JobTask) error {
	return r.db.WithContext(ctx).Create(task).Error
}

func (r *jobTaskRepo) GetByID(ctx context.Context, id int64) (*model.JobTask, error) {
	var task model.JobTask
	if err := r.db.WithContext(ctx).Where("id = ? AND deleted = 0", id).First(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *jobTaskRepo) Update(ctx context.Context, id int64, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.JobTask{}).Where("id = ? AND deleted = 0", id).Updates(updates).Error
}

func (r *jobTaskRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&model.JobTask{}).Where("id = ?", id).Update("deleted", 1).Error
}

func (r *jobTaskRepo) Page(ctx context.Context, req *model.JobTaskPageReq, tenantID int64) ([]model.JobTask, int64, error) {
	var list []model.JobTask
	var total int64
	db := r.db.WithContext(ctx).Model(&model.JobTask{}).Where("tenant_id = ? AND deleted = 0", tenantID)

	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}

	db.Count(&total)
	db.Offset(req.Offset()).Limit(req.Limit()).Order("id DESC").Find(&list)
	return list, total, nil
}

func (r *jobTaskRepo) ListActive(ctx context.Context, tenantID int64) ([]model.JobTask, error) {
	var list []model.JobTask
	err := r.db.WithContext(ctx).Where("status = 0 AND deleted = 0").Find(&list).Error
	return list, err
}
