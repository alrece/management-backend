package service

import (
	"context"
	"fmt"
	"time"

	"management-backend/services/job-service/internal/model"
	"management-backend/services/job-service/internal/repository"

	"management-backend/pkg/snowflake"
)

type JobService interface {
	Create(ctx context.Context, req *model.JobTaskCreateReq, creator, tenantID int64) (int64, error)
	Update(ctx context.Context, id int64, req *model.JobTaskUpdateReq) error
	Delete(ctx context.Context, id int64) error
	Page(ctx context.Context, req *model.JobTaskPageReq, tenantID int64) ([]model.JobTaskResp, int64, error)
	GetByID(ctx context.Context, id int64) (*model.JobTask, error)
	Trigger(ctx context.Context, id int64, tenantID int64) (int64, error)
	ExecLogPage(ctx context.Context, req *model.ExecLogPageReq, tenantID int64) ([]model.ExecLogResp, int64, error)
}

type jobService struct {
	taskRepo repository.JobTaskRepo
	logRepo  repository.ExecLogRepo
}

func NewJobService(taskRepo repository.JobTaskRepo, logRepo repository.ExecLogRepo) JobService {
	return &jobService{taskRepo: taskRepo, logRepo: logRepo}
}

func (s *jobService) Create(ctx context.Context, req *model.JobTaskCreateReq, creator, tenantID int64) (int64, error) {
	status := int8(model.TaskStatusNormal)
	if req.Status != nil {
		status = *req.Status
	}
	task := &model.JobTask{
		ID:       snowflake.NextID(),
		TenantID: tenantID,
		Name:     req.Name,
		Handler:  req.Handler,
		CronExpr: req.CronExpr,
		Params:   req.Params,
		Status:   status,
		Remark:   req.Remark,
		Creator:  &creator,
	}
	if err := s.taskRepo.Create(ctx, task); err != nil {
		return 0, fmt.Errorf("创建任务失败: %w", err)
	}
	return task.ID, nil
}

func (s *jobService) Update(ctx context.Context, id int64, req *model.JobTaskUpdateReq) error {
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Handler != "" {
		updates["handler"] = req.Handler
	}
	if req.CronExpr != "" {
		updates["cron_expr"] = req.CronExpr
	}
	if req.Params != "" {
		updates["params"] = req.Params
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}
	if len(updates) == 0 {
		return nil
	}
	return s.taskRepo.Update(ctx, id, updates)
}

func (s *jobService) Delete(ctx context.Context, id int64) error {
	return s.taskRepo.Delete(ctx, id)
}

func (s *jobService) GetByID(ctx context.Context, id int64) (*model.JobTask, error) {
	return s.taskRepo.GetByID(ctx, id)
}

func (s *jobService) Page(ctx context.Context, req *model.JobTaskPageReq, tenantID int64) ([]model.JobTaskResp, int64, error) {
	tasks, total, err := s.taskRepo.Page(ctx, req, tenantID)
	if err != nil {
		return nil, 0, err
	}
	list := make([]model.JobTaskResp, 0, len(tasks))
	for _, t := range tasks {
		list = append(list, model.JobTaskResp{
			ID: t.ID, Name: t.Name, Handler: t.Handler,
			CronExpr: t.CronExpr, Params: t.Params,
			Status: t.Status, Remark: t.Remark, Creator: t.Creator,
		})
	}
	return list, total, nil
}

func (s *jobService) Trigger(ctx context.Context, id int64, tenantID int64) (int64, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("任务不存在: %w", err)
	}
	if task.TenantID != tenantID {
		return 0, fmt.Errorf("无权访问该任务")
	}

	logID := snowflake.NextID()
	execLog := &model.JobExecutionLog{
		ID:          logID,
		TenantID:    tenantID,
		TaskID:      task.ID,
		TaskName:    task.Name,
		TriggerType: model.TriggerTypeManual,
		Status:      model.ExecStatusRunning,
		StartTime:   time.Now(),
	}
	if err := s.logRepo.Create(ctx, execLog); err != nil {
		return 0, err
	}
	return logID, nil
}

func (s *jobService) ExecLogPage(ctx context.Context, req *model.ExecLogPageReq, tenantID int64) ([]model.ExecLogResp, int64, error) {
	logs, total, err := s.logRepo.Page(ctx, req, tenantID)
	if err != nil {
		return nil, 0, err
	}
	list := make([]model.ExecLogResp, 0, len(logs))
	for _, l := range logs {
		list = append(list, model.ExecLogResp{
			ID: l.ID, TaskID: l.TaskID, TaskName: l.TaskName,
			TriggerType: l.TriggerType, Status: l.Status,
			DurationMs: l.DurationMs, Result: l.Result,
			Error: l.Error, StartTime: l.StartTime.Format(time.RFC3339),
		})
	}
	return list, total, nil
}
