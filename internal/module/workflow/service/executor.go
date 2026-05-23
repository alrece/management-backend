package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"sync"
	"time"

	wfmodel "management-backend/internal/module/workflow/model"

	pkgmiddleware "management-backend/pkg/middleware"
	"management-backend/pkg/snowflake"

	"management-backend/internal/module/workflow/engine"
	"management-backend/internal/module/workflow/repository"
)

const (
	defaultWorkerCount = 50
	defaultNodeTimeout = 60 * time.Second
	maxRetries         = 3
)

// Executor DAG 执行引擎
type Executor struct {
	registry   *engine.NodeRegistry
	execRepo   repository.ExecutionRepo
	logRepo    repository.NodeLogRepo
	workerPool chan struct{}
}

func NewExecutor(
	registry *engine.NodeRegistry,
	execRepo repository.ExecutionRepo,
	logRepo repository.NodeLogRepo,
	workerCount int,
) *Executor {
	if workerCount <= 0 {
		workerCount = defaultWorkerCount
	}
	return &Executor{
		registry:   registry,
		execRepo:   execRepo,
		logRepo:    logRepo,
		workerPool: make(chan struct{}, workerCount),
	}
}

// ExecuteWorkflow 执行工作流
func (e *Executor) ExecuteWorkflow(ctx context.Context, wf *wfmodel.Workflow, input json.RawMessage) (*wfmodel.Execution, error) {
	// 解析工作流定义
	var def wfmodel.WorkflowDefinition
	if err := json.Unmarshal(wf.Definition, &def); err != nil {
		return nil, fmt.Errorf("解析工作流定义失败: %w", err)
	}

	// ENG-007: 节点数量上限 100
	if len(def.Nodes) > 100 {
		return nil, fmt.Errorf("节点数量 %d 超过上限 100", len(def.Nodes))
	}

	// 构建 DAG 节点和边
	dagNodes := make([]engine.DAGNode, 0, len(def.Nodes))
	for _, n := range def.Nodes {
		dagNodes = append(dagNodes, engine.DAGNode{ID: n.ID, Type: n.Type})
	}
	dagEdges := make([]engine.DAGEdge, 0, len(def.Edges))
	for _, edge := range def.Edges {
		dagEdges = append(dagEdges, engine.DAGEdge{
			Source:       edge.Source,
			Target:       edge.Target,
			SourceHandle: edge.SourceHandle,
		})
	}

	// 拓扑排序
	result, err := engine.TopologicalSort(dagNodes, dagEdges)
	if err != nil {
		return nil, fmt.Errorf("拓扑排序失败: %w", err)
	}

	// 创建执行记录
	tenantID := pkgmiddleware.GetTenantID(ctx)
	execution := &wfmodel.Execution{
		ID:          snowflake.NextID(),
		TenantID:    tenantID,
		WorkflowID:  wf.ID,
		Status:      "running",
		TriggerType: wf.TriggerType,
		Input:       input,
	}
	execution.Creator = pkgmiddleware.GetUserID(ctx)

	if err := e.execRepo.Create(ctx, execution); err != nil {
		return nil, fmt.Errorf("创建执行记录失败: %w", err)
	}

	// 执行上下文
	execCtx := engine.NewExecutionContext(ctx, execution.ID, wf.ID, tenantID)

	// 按层执行
	var execErr error
	for _, layer := range result.Layers {
		if execErr != nil {
			break
		}
		execErr = e.executeLayer(execCtx, layer, &def)
	}

	// 更新执行状态
	now := time.Now()
	if execErr != nil {
		e.execRepo.UpdateStatus(ctx, execution.ID, "failed", execErr.Error())
		execution.Status = "failed"
		execution.ErrorMsg = execErr.Error()
	} else {
		e.execRepo.UpdateStatus(ctx, execution.ID, "success", "")
		execution.Status = "success"
	}
	execution.EndTime = &now

	return execution, nil
}

// executeLayer 执行一层节点（并行）
func (e *Executor) executeLayer(execCtx *engine.ExecutionContext, layer []engine.DAGNode, def *wfmodel.WorkflowDefinition) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(layer))

	for _, node := range layer {
		wg.Add(1)
		// 获取 goroutine pool 令牌
		e.workerPool <- struct{}{}

		go func(n engine.DAGNode) {
			defer wg.Done()
			defer func() { <-e.workerPool }()

			if err := e.executeNode(execCtx, n, def); err != nil {
				errCh <- err
			}
		}(node)
	}

	wg.Wait()
	close(errCh)

	// 检查是否有错误
	for err := range errCh {
		return err
	}
	return nil
}

// executeNode 执行单个节点（含超时和错误策略）
func (e *Executor) executeNode(execCtx *engine.ExecutionContext, dagNode engine.DAGNode, def *wfmodel.WorkflowDefinition) error {
	// 查找节点定义
	var nodeDef *wfmodel.WorkflowNodeDef
	for i := range def.Nodes {
		if def.Nodes[i].ID == dagNode.ID {
			nodeDef = &def.Nodes[i]
			break
		}
	}
	if nodeDef == nil {
		return fmt.Errorf("节点定义不存在: %s", dagNode.ID)
	}

	// 获取节点处理器
	handler, ok := e.registry.Get(nodeDef.Type)
	if !ok {
		return fmt.Errorf("未注册的节点类型: %s", nodeDef.Type)
	}

	// 构建节点输入
	input := e.buildNodeInput(execCtx, dagNode, nodeDef)

	// 获取上游数据
	input.Upstream = make(map[string][]map[string]any)
	for _, edge := range def.Edges {
		if edge.Target == dagNode.ID {
			if out, exists := execCtx.GetOutput(edge.Source); exists {
				input.Upstream[edge.Source] = out.Items
			}
		}
	}

	// 合并上游数据到 items
	for _, items := range input.Upstream {
		input.Items = append(input.Items, items...)
	}

	// 节点超时（ENG-003）
	timeout := defaultNodeTimeout
	if nodeDef.Data.Timeout > 0 {
		timeout = time.Duration(nodeDef.Data.Timeout) * time.Second
	}

	// 创建节点日志
	nodeLog := &wfmodel.NodeLog{
		ID:          snowflake.NextID(),
		TenantID:    execCtx.TenantID,
		ExecutionID: execCtx.ExecutionID,
		WorkflowID:  execCtx.WorkflowID,
		NodeID:      dagNode.ID,
		NodeType:    dagNode.Type,
		Status:      "running",
		StartTime:   time.Now(),
	}
	logInput, _ := json.Marshal(input)
	nodeLog.Input = logInput

	ctx, cancel := execCtx.WithTimeout(timeout)
	defer cancel()

	// 执行节点（含重试策略）
	errorStrategy := nodeDef.Data.ErrorStrategy
	if errorStrategy == "" {
		errorStrategy = "retry"
	}
	maxRetry := maxRetries
	if nodeDef.Data.RetryCount > 0 {
		maxRetry = nodeDef.Data.RetryCount
	}

	var output *engine.NodeOutput
	var execErr error

	switch errorStrategy {
	case "retry":
		output, execErr = e.executeWithRetry(ctx, handler, input, maxRetry)
	case "skip":
		output, execErr = e.executeOnce(ctx, handler, input)
		if execErr != nil {
			slog.Warn("节点执行跳过", "nodeId", dagNode.ID, "error", execErr)
			output = &engine.NodeOutput{Items: []map[string]any{}}
			execErr = nil // 跳过不算错误
		}
	case "abort":
		output, execErr = e.executeOnce(ctx, handler, input)
	default:
		output, execErr = e.executeWithRetry(ctx, handler, input, maxRetry)
	}

	// 更新节点日志
	now := time.Now()
	nodeLog.EndTime = &now
	if execErr != nil {
		nodeLog.Status = "failed"
		nodeLog.ErrorMsg = execErr.Error()
	} else {
		nodeLog.Status = "success"
		logOutput, _ := json.Marshal(output)
		nodeLog.Output = logOutput
		execCtx.SetOutput(dagNode.ID, output)
	}

	if logErr := e.logRepo.Create(execCtx.Ctx, nodeLog); logErr != nil {
		slog.Error("写入节点日志失败", "nodeId", dagNode.ID, "error", logErr)
	}

	return execErr
}

// executeWithRetry 带重试的执行（指数退避: 1s, 2s, 4s）
func (e *Executor) executeWithRetry(ctx context.Context, handler engine.NodeHandler, input *engine.NodeInput, maxRetry int) (*engine.NodeOutput, error) {
	var lastErr error
	for attempt := range maxRetry {
		output, err := handler.Execute(ctx, input)
		if err == nil {
			return output, nil
		}
		lastErr = err
		if attempt < maxRetry-1 {
			backoff := time.Duration(math.Pow(2, float64(attempt))) * time.Second
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}
		}
	}
	return nil, fmt.Errorf("重试 %d 次后仍失败: %w", maxRetry, lastErr)
}

// executeOnce 单次执行
func (e *Executor) executeOnce(ctx context.Context, handler engine.NodeHandler, input *engine.NodeInput) (*engine.NodeOutput, error) {
	return handler.Execute(ctx, input)
}

// buildNodeInput 构建节点输入
func (e *Executor) buildNodeInput(_ *engine.ExecutionContext, dagNode engine.DAGNode, nodeDef *wfmodel.WorkflowNodeDef) *engine.NodeInput {
	return &engine.NodeInput{
		NodeID:   dagNode.ID,
		NodeType: dagNode.Type,
		Config:   nodeDef.Data.Config,
		Items:    []map[string]any{},
	}
}
