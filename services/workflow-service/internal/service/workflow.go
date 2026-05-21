package service

import (
	"context"
	"fmt"
	"time"

	"management-backend/services/workflow-service/internal/client"
	"management-backend/services/workflow-service/internal/model"
	"management-backend/services/workflow-service/internal/repository"
	"management-backend/pkg/snowflake"
)

type WorkflowService interface {
	Create(ctx context.Context, req *model.WorkflowCreateReq, creator, tenantID int64) (int64, error)
	Update(ctx context.Context, id int64, req *model.WorkflowUpdateReq) error
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (*model.WorkflowResp, error)
	Page(ctx context.Context, req *model.WorkflowPageReq, tenantID int64) ([]model.WorkflowResp, int64, error)
	Activate(ctx context.Context, id int64) error
	Deactivate(ctx context.Context, id int64) error
	Execute(ctx context.Context, id int64, req *model.ExecuteReq, creator, tenantID int64) (int64, error)
	SyncFromN8n(ctx context.Context) error
	StartupReconcile(ctx context.Context) error
}

type workflowService struct {
	workflowRepo repository.WorkflowRepo
	instanceRepo repository.InstanceRepo
	n8nClient    *client.N8nClient
}

func NewWorkflowService(wfRepo repository.WorkflowRepo, instRepo repository.InstanceRepo, n8n *client.N8nClient) WorkflowService {
	return &workflowService{
		workflowRepo: wfRepo,
		instanceRepo: instRepo,
		n8nClient:    n8n,
	}
}

func (s *workflowService) Create(ctx context.Context, req *model.WorkflowCreateReq, creator, tenantID int64) (int64, error) {
	// 在 n8n 创建工作流
	n8nResp, err := s.n8nClient.CreateWorkflow(ctx, req.Name)
	if err != nil {
		return 0, fmt.Errorf("创建 n8n 工作流失败: %w", err)
	}

	wf := &model.WfWorkflow{
		ID:            snowflake.NextID(),
		TenantID:      tenantID,
		Name:          req.Name,
		CategoryID:    req.CategoryID,
		N8nWorkflowID: n8nResp.ID,
		ParamsSchema:  req.ParamsSchema,
		Status:        "DRAFT",
		Creator:       creator,
		Updater:       creator,
	}
	if err := s.workflowRepo.Create(ctx, wf); err != nil {
		return 0, err
	}
	return wf.ID, nil
}

func (s *workflowService) Update(ctx context.Context, id int64, req *model.WorkflowUpdateReq) error {
	wf, err := s.workflowRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	wf.Name = req.Name
	wf.CategoryID = req.CategoryID
	wf.ParamsSchema = req.ParamsSchema
	wf.Updater = 0
	return s.workflowRepo.Update(ctx, wf)
}

func (s *workflowService) Delete(ctx context.Context, id int64) error {
	wf, err := s.workflowRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	// 在 n8n 删除工作流（忽略错误，本地仍删除）
	_ = s.n8nClient.DeleteWorkflow(ctx, wf.N8nWorkflowID)
	return s.workflowRepo.Delete(ctx, id)
}

func (s *workflowService) Get(ctx context.Context, id int64) (*model.WorkflowResp, error) {
	wf, err := s.workflowRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toWorkflowResp(wf), nil
}

func (s *workflowService) Page(ctx context.Context, req *model.WorkflowPageReq, tenantID int64) ([]model.WorkflowResp, int64, error) {
	list, total, err := s.workflowRepo.Page(ctx, req, tenantID)
	if err != nil {
		return nil, 0, err
	}
	resp := make([]model.WorkflowResp, 0, len(list))
	for i := range list {
		resp = append(resp, *toWorkflowResp(&list[i]))
	}
	return resp, total, nil
}

func (s *workflowService) Activate(ctx context.Context, id int64) error {
	wf, err := s.workflowRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.n8nClient.ActivateWorkflow(ctx, wf.N8nWorkflowID); err != nil {
		// n8n 不可用时保持 DRAFT，不标记 ACTIVE
		return fmt.Errorf("激活 n8n 工作流失败: %w", err)
	}
	return s.workflowRepo.UpdateStatus(ctx, id, "ACTIVE")
}

func (s *workflowService) Deactivate(ctx context.Context, id int64) error {
	wf, err := s.workflowRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	_ = s.n8nClient.DeactivateWorkflow(ctx, wf.N8nWorkflowID)
	return s.workflowRepo.UpdateStatus(ctx, id, "INACTIVE")
}

func (s *workflowService) Execute(ctx context.Context, id int64, req *model.ExecuteReq, creator, tenantID int64) (int64, error) {
	wf, err := s.workflowRepo.GetByID(ctx, id)
	if err != nil {
		return 0, err
	}
	if wf.Status != "ACTIVE" {
		return 0, fmt.Errorf("工作流未激活，无法执行")
	}

	n8nExec, err := s.n8nClient.ExecuteWorkflow(ctx, wf.N8nWorkflowID)
	if err != nil {
		return 0, fmt.Errorf("执行 n8n 工作流失败: %w", err)
	}

	inst := &model.WfInstance{
		ID:             snowflake.NextID(),
		TenantID:       tenantID,
		WorkflowID:     id,
		N8nExecutionID: n8nExec.ID,
		Status:         "RUNNING",
		Variables:      req.Variables,
		StartedAt:      time.Now(),
		Creator:        creator,
	}
	if err := s.instanceRepo.Create(ctx, inst); err != nil {
		return 0, err
	}
	return inst.ID, nil
}

// SyncFromN8n 定期同步 n8n 状态，检测 DESYNCED
func (s *workflowService) SyncFromN8n(ctx context.Context) error {
	workflows, err := s.workflowRepo.ListByStatus(ctx, "ACTIVE")
	if err != nil {
		return err
	}
	for _, wf := range workflows {
		_, err := s.n8nClient.GetWorkflow(ctx, wf.N8nWorkflowID)
		if err != nil {
			// n8n 中不存在，标记 DESYNCED
			_ = s.workflowRepo.UpdateStatus(ctx, wf.ID, "DESYNCED")
		}
	}
	return nil
}

// StartupReconcile 启动时对账 RUNNING 实例
func (s *workflowService) StartupReconcile(ctx context.Context) error {
	workflows, err := s.workflowRepo.ListByStatus(ctx, "ACTIVE")
	if err != nil {
		return err
	}
	for _, wf := range workflows {
		// 检查 n8n 中工作流是否仍存在
		_, err := s.n8nClient.GetWorkflow(ctx, wf.N8nWorkflowID)
		if err != nil {
			_ = s.workflowRepo.UpdateStatus(ctx, wf.ID, "DESYNCED")
		}
	}
	return nil
}

func toWorkflowResp(wf *model.WfWorkflow) *model.WorkflowResp {
	return &model.WorkflowResp{
		ID:            wf.ID,
		Name:          wf.Name,
		CategoryID:    wf.CategoryID,
		N8nWorkflowID: wf.N8nWorkflowID,
		ParamsSchema:  wf.ParamsSchema,
		Status:        wf.Status,
		Creator:       wf.Creator,
		CreatedAt:     wf.CreatedAt.Format(time.RFC3339),
		Updater:       wf.Updater,
		UpdatedAt:     wf.UpdatedAt.Format(time.RFC3339),
	}
}
