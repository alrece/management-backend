package engine

import (
	"context"
	"time"
)

// NodeInput 节点执行输入
type NodeInput struct {
	NodeID   string               `json:"nodeId"`
	NodeType string               `json:"nodeType"`
	Config   map[string]any       `json:"config"`
	Items    []map[string]any     `json:"items"`
	// 上游节点输出映射（nodeID → items）
	Upstream map[string][]map[string]any `json:"-"`
}

// NodeOutput 节点执行输出
type NodeOutput struct {
	Items []map[string]any `json:"items"`
	Error string           `json:"error,omitempty"`
}

// ExecutionContext 工作流执行上下文（跨节点传递数据）
type ExecutionContext struct {
	// 执行 ID
	ExecutionID int64
	// 工作流 ID
	WorkflowID int64
	// 租户 ID
	TenantID int64
	// 节点输出缓存（nodeID → NodeOutput）
	Outputs map[string]*NodeOutput
	// Go context（超时/取消）
	Ctx context.Context
}

// NewExecutionContext 创建执行上下文
func NewExecutionContext(ctx context.Context, executionID, workflowID, tenantID int64) *ExecutionContext {
	return &ExecutionContext{
		ExecutionID: executionID,
		WorkflowID:  workflowID,
		TenantID:    tenantID,
		Outputs:     make(map[string]*NodeOutput),
		Ctx:         ctx,
	}
}

// SetOutput 设置节点输出
func (ec *ExecutionContext) SetOutput(nodeID string, output *NodeOutput) {
	ec.Outputs[nodeID] = output
}

// GetOutput 获取节点输出
func (ec *ExecutionContext) GetOutput(nodeID string) (*NodeOutput, bool) {
	out, ok := ec.Outputs[nodeID]
	return out, ok
}

// WithTimeout 创建带超时的子上下文
func (ec *ExecutionContext) WithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ec.Ctx, timeout)
}
