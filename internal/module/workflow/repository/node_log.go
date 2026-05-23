package repository

import (
	"context"
	"time"

	wfmodel "management-backend/internal/module/workflow/model"

	pkgmiddleware "management-backend/pkg/middleware"

	"gorm.io/gorm"
)

// NodeLogRepo 节点日志数据访问接口
type NodeLogRepo interface {
	Create(ctx context.Context, log *wfmodel.NodeLog) error
	Update(ctx context.Context, log *wfmodel.NodeLog) error
	ListByExecutionID(ctx context.Context, executionID int64) ([]wfmodel.NodeLog, error)
	CleanBefore(ctx context.Context, before time.Time, limit int) (int64, error)
}

type nodeLogRepo struct {
	db *gorm.DB
}

func NewNodeLogRepo(db *gorm.DB) NodeLogRepo {
	return &nodeLogRepo{db: db}
}

func (r *nodeLogRepo) getDB(ctx context.Context) *gorm.DB {
	if tenantDB := pkgmiddleware.GetTenantDB(ctx); tenantDB != nil {
		return tenantDB
	}
	return r.db
}

func (r *nodeLogRepo) Create(ctx context.Context, log *wfmodel.NodeLog) error {
	return r.getDB(ctx).WithContext(ctx).Create(log).Error
}

func (r *nodeLogRepo) Update(ctx context.Context, log *wfmodel.NodeLog) error {
	return r.getDB(ctx).WithContext(ctx).Save(log).Error
}

func (r *nodeLogRepo) ListByExecutionID(ctx context.Context, executionID int64) ([]wfmodel.NodeLog, error) {
	var logs []wfmodel.NodeLog
	err := r.getDB(ctx).WithContext(ctx).
		Where("execution_id = ?", executionID).
		Order("start_time ASC").Find(&logs).Error
	return logs, err
}

// CleanBefore 清理指定时间前的日志（ENG-010：每次上限 1000 条）
func (r *nodeLogRepo) CleanBefore(ctx context.Context, before time.Time, limit int) (int64, error) {
	if limit <= 0 || limit > 1000 {
		limit = 1000
	}

	// 先查要清理的执行记录 ID
	var execIDs []int64
	err := r.getDB(ctx).WithContext(ctx).Model(&wfmodel.Execution{}).
		Select("id").Where("start_time < ?", before).
		Limit(limit).Find(&execIDs).Error
	if err != nil || len(execIDs) == 0 {
		return 0, err
	}

	// 先删节点日志
	r.getDB(ctx).WithContext(ctx).
		Where("execution_id IN ?", execIDs).Delete(&wfmodel.NodeLog{})

	// 再删执行记录
	result := r.getDB(ctx).WithContext(ctx).
		Where("id IN ?", execIDs).Delete(&wfmodel.Execution{})

	return result.RowsAffected, result.Error
}
