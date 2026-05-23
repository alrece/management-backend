package service

import (
	"context"
	"encoding/json"
	"fmt"

	wfmodel "management-backend/internal/module/workflow/model"

	"management-backend/internal/module/workflow/engine"
	"management-backend/internal/module/workflow/repository"
	"management-backend/pkg/snowflake"
)

// WorkflowService 工作流业务接口
type WorkflowService interface {
	Create(ctx context.Context, req *wfmodel.WorkflowCreateReq, creator int64, tenantID int64) (int64, error)
	Update(ctx context.Context, req *wfmodel.WorkflowUpdateReq, updater int64) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*wfmodel.WorkflowResp, error)
	Page(ctx context.Context, req *wfmodel.WorkflowPageReq) ([]wfmodel.WorkflowResp, int64, error)
	Activate(ctx context.Context, id int64) error
	Deactivate(ctx context.Context, id int64) error
	Execute(ctx context.Context, id int64, input json.RawMessage) (*wfmodel.Execution, error)
	GetAllSchemas(ctx context.Context) []engine.NodeSchema
}

type workflowService struct {
	repo      repository.WorkflowRepo
	executor  *Executor
	registry  *engine.NodeRegistry
}

func NewWorkflowService(
	repo repository.WorkflowRepo,
	executor *Executor,
	registry *engine.NodeRegistry,
) WorkflowService {
	return &workflowService{repo: repo, executor: executor, registry: registry}
}

func (s *workflowService) Create(ctx context.Context, req *wfmodel.WorkflowCreateReq, creator int64, tenantID int64) (int64, error) {
	// 校验 definition（循环检测 + 节点数限制）
	if err := s.validateDefinition(req.Definition); err != nil {
		return 0, err
	}

	wf := &wfmodel.Workflow{
		Name:          req.Name,
		CategoryID:    req.CategoryID,
		Definition:    req.Definition,
		TriggerType:   req.TriggerType,
		TriggerConfig: req.TriggerConfig,
	}
	wf.ID = snowflake.NextID()
	wf.TenantID = tenantID
	wf.Creator = creator
	wf.Updater = creator
	wf.Remark = req.Remark

	if err := s.repo.Create(ctx, wf); err != nil {
		return 0, fmt.Errorf("创建工作流失败: %w", err)
	}
	return wf.ID, nil
}

func (s *workflowService) Update(ctx context.Context, req *wfmodel.WorkflowUpdateReq, updater int64) error {
	existing, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("工作流不存在")
	}

	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.CategoryID != nil {
		existing.CategoryID = *req.CategoryID
	}
	if req.Definition != nil {
		if err := s.validateDefinition(req.Definition); err != nil {
			return err
		}
		existing.Definition = req.Definition
		existing.Version++
	}
	if req.TriggerType != "" {
		existing.TriggerType = req.TriggerType
	}
	if req.TriggerConfig != nil {
		existing.TriggerConfig = req.TriggerConfig
	}
	if req.Status != nil {
		existing.Status = *req.Status
	}
	existing.Remark = req.Remark
	existing.Updater = updater

	return s.repo.Update(ctx, existing)
}

func (s *workflowService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *workflowService) GetByID(ctx context.Context, id int64) (*wfmodel.WorkflowResp, error) {
	wf, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("工作流不存在")
	}
	return toWorkflowResp(wf), nil
}

func (s *workflowService) Page(ctx context.Context, req *wfmodel.WorkflowPageReq) ([]wfmodel.WorkflowResp, int64, error) {
	list, total, err := s.repo.Page(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	resp := make([]wfmodel.WorkflowResp, 0, len(list))
	for i := range list {
		resp = append(resp, *toWorkflowResp(&list[i]))
	}
	return resp, total, nil
}

func (s *workflowService) Activate(ctx context.Context, id int64) error {
	return s.repo.UpdateStatus(ctx, id, 0) // 0 = 正常/启用
}

func (s *workflowService) Deactivate(ctx context.Context, id int64) error {
	return s.repo.UpdateStatus(ctx, id, 1) // 1 = 停用
}

func (s *workflowService) Execute(ctx context.Context, id int64, input json.RawMessage) (*wfmodel.Execution, error) {
	wf, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("工作流不存在")
	}
	if wf.Status != 0 {
		return nil, fmt.Errorf("工作流已停用")
	}
	return s.executor.ExecuteWorkflow(ctx, wf, input)
}

func (s *workflowService) GetAllSchemas(_ context.Context) []engine.NodeSchema {
	return s.registry.AllSchemas()
}

// validateDefinition 校验工作流定义（ENG-002 循环检测 + ENG-007 节点数限制）
func (s *workflowService) validateDefinition(def json.RawMessage) error {
	var wfDef wfmodel.WorkflowDefinition
	if err := json.Unmarshal(def, &wfDef); err != nil {
		return fmt.Errorf("无效的工作流定义: %w", err)
	}

	// ENG-007: 节点数上限 100
	if len(wfDef.Nodes) > 100 {
		return fmt.Errorf("节点数量 %d 超过上限 100", len(wfDef.Nodes))
	}

	// ENG-002: 循环检测
	dagNodes := make([]engine.DAGNode, 0, len(wfDef.Nodes))
	for _, n := range wfDef.Nodes {
		dagNodes = append(dagNodes, engine.DAGNode{ID: n.ID, Type: n.Type})
	}
	dagEdges := make([]engine.DAGEdge, 0, len(wfDef.Edges))
	for _, e := range wfDef.Edges {
		dagEdges = append(dagEdges, engine.DAGEdge{
			Source: e.Source, Target: e.Target, SourceHandle: e.SourceHandle,
		})
	}

	if err := engine.DetectCycle(dagNodes, dagEdges); err != nil {
		return fmt.Errorf("工作流定义校验失败: %w", err)
	}

	return nil
}

func toWorkflowResp(wf *wfmodel.Workflow) *wfmodel.WorkflowResp {
	return &wfmodel.WorkflowResp{
		ID:            wf.ID,
		Name:          wf.Name,
		CategoryID:    wf.CategoryID,
		Definition:    wf.Definition,
		TriggerType:   wf.TriggerType,
		TriggerConfig: wf.TriggerConfig,
		Version:       wf.Version,
		Status:        wf.Status,
		Remark:        wf.Remark,
		TenantID:      wf.TenantID,
		Creator:       wf.Creator,
		CreateTime:    wf.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateTime:    wf.UpdateTime.Format("2006-01-02 15:04:05"),
	}
}
