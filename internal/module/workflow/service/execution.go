package service

import (
	"context"

	wfmodel "management-backend/internal/module/workflow/model"

	"management-backend/internal/module/workflow/repository"
)

// ExecutionService 执行实例业务接口
type ExecutionService interface {
	Page(ctx context.Context, req *wfmodel.ExecutionPageReq) ([]wfmodel.ExecutionResp, int64, error)
	GetByID(ctx context.Context, id int64) (*wfmodel.ExecutionResp, error)
	GetLogs(ctx context.Context, executionID int64) ([]wfmodel.NodeLogResp, error)
}

type executionService struct {
	execRepo repository.ExecutionRepo
	logRepo  repository.NodeLogRepo
}

func NewExecutionService(execRepo repository.ExecutionRepo, logRepo repository.NodeLogRepo) ExecutionService {
	return &executionService{execRepo: execRepo, logRepo: logRepo}
}

func (s *executionService) Page(ctx context.Context, req *wfmodel.ExecutionPageReq) ([]wfmodel.ExecutionResp, int64, error) {
	list, total, err := s.execRepo.Page(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	resp := make([]wfmodel.ExecutionResp, 0, len(list))
	for i := range list {
		resp = append(resp, *toExecutionResp(&list[i]))
	}
	return resp, total, nil
}

func (s *executionService) GetByID(ctx context.Context, id int64) (*wfmodel.ExecutionResp, error) {
	exec, err := s.execRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toExecutionResp(exec), nil
}

func (s *executionService) GetLogs(ctx context.Context, executionID int64) ([]wfmodel.NodeLogResp, error) {
	logs, err := s.logRepo.ListByExecutionID(ctx, executionID)
	if err != nil {
		return nil, err
	}
	resp := make([]wfmodel.NodeLogResp, 0, len(logs))
	for i := range logs {
		resp = append(resp, *toNodeLogResp(&logs[i]))
	}
	return resp, nil
}

func toExecutionResp(e *wfmodel.Execution) *wfmodel.ExecutionResp {
	resp := &wfmodel.ExecutionResp{
		ID:          e.ID,
		WorkflowID:  e.WorkflowID,
		Status:      e.Status,
		TriggerType: e.TriggerType,
		StartTime:   e.StartTime.Format("2006-01-02 15:04:05"),
		Input:       e.Input,
		Output:      e.Output,
		ErrorMsg:    e.ErrorMsg,
	}
	if e.EndTime != nil {
		resp.EndTime = e.EndTime.Format("2006-01-02 15:04:05")
	}
	return resp
}

func toNodeLogResp(l *wfmodel.NodeLog) *wfmodel.NodeLogResp {
	resp := &wfmodel.NodeLogResp{
		ID:          l.ID,
		ExecutionID: l.ExecutionID,
		NodeID:      l.NodeID,
		NodeType:    l.NodeType,
		Status:      l.Status,
		Input:       l.Input,
		Output:      l.Output,
		ErrorMsg:    l.ErrorMsg,
		RetryCount:  l.RetryCount,
		StartTime:   l.StartTime.Format("2006-01-02 15:04:05"),
	}
	if l.EndTime != nil {
		resp.EndTime = l.EndTime.Format("2006-01-02 15:04:05")
	}
	return resp
}
