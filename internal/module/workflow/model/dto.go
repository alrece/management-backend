package model

import (
	"encoding/json"

	"management-backend/pkg/page"
)

// === 工作流 DTO ===

// WorkflowCreateReq 创建工作流请求
type WorkflowCreateReq struct {
	Name         string          `json:"name" binding:"required,max=100"`
	CategoryID   int64           `json:"categoryId"`
	Definition   json.RawMessage `json:"definition" binding:"required"`
	TriggerType  string          `json:"triggerType" binding:"required,oneof=manual cron webhook"`
	TriggerConfig json.RawMessage `json:"triggerConfig"`
	Remark       string          `json:"remark" binding:"omitempty,max=500"`
}

// WorkflowUpdateReq 更新工作流请求
type WorkflowUpdateReq struct {
	ID           int64           `json:"id" binding:"required"`
	Name         string          `json:"name" binding:"omitempty,max=100"`
	CategoryID   *int64          `json:"categoryId"`
	Definition   json.RawMessage `json:"definition"`
	TriggerType  string          `json:"triggerType" binding:"omitempty,oneof=manual cron webhook"`
	TriggerConfig json.RawMessage `json:"triggerConfig"`
	Status       *int            `json:"status" binding:"omitempty,oneof=0 1"`
	Remark       string          `json:"remark" binding:"omitempty,max=500"`
}

// WorkflowPageReq 工作流分页查询
type WorkflowPageReq struct {
	page.Req
	Name       string `form:"name" json:"name"`
	Status     *int   `form:"status" json:"status"`
	CategoryID *int64 `form:"categoryId" json:"categoryId"`
	TriggerType string `form:"triggerType" json:"triggerType"`
}

// WorkflowResp 工作流响应
type WorkflowResp struct {
	ID           int64           `json:"id"`
	Name         string          `json:"name"`
	CategoryID   int64           `json:"categoryId"`
	Definition   json.RawMessage `json:"definition"`
	TriggerType  string          `json:"triggerType"`
	TriggerConfig json.RawMessage `json:"triggerConfig"`
	Version      int             `json:"version"`
	Status       int             `json:"status"`
	Remark       string          `json:"remark"`
	TenantID     int64           `json:"tenantId"`
	Creator      int64           `json:"creator"`
	CreateTime   string          `json:"createTime"`
	UpdateTime   string          `json:"updateTime"`
}

// === 执行实例 DTO ===

// ExecutionPageReq 执行实例分页查询
type ExecutionPageReq struct {
	page.Req
	WorkflowID *int64  `form:"workflowId" json:"workflowId"`
	Status     string  `form:"status" json:"status" binding:"omitempty,oneof=running success failed canceled"`
}

// ExecutionResp 执行实例响应
type ExecutionResp struct {
	ID          int64           `json:"id"`
	WorkflowID  int64           `json:"workflowId"`
	WorkflowName string         `json:"workflowName"`
	Status      string          `json:"status"`
	TriggerType string          `json:"triggerType"`
	StartTime   string          `json:"startTime"`
	EndTime     string          `json:"endTime,omitempty"`
	Input       json.RawMessage `json:"input"`
	Output      json.RawMessage `json:"output"`
	ErrorMsg    string          `json:"errorMsg"`
}

// NodeLogResp 节点执行日志响应
type NodeLogResp struct {
	ID          int64           `json:"id"`
	ExecutionID int64           `json:"executionId"`
	NodeID      string          `json:"nodeId"`
	NodeType    string          `json:"nodeType"`
	Status      string          `json:"status"`
	Input       json.RawMessage `json:"input"`
	Output      json.RawMessage `json:"output"`
	ErrorMsg    string          `json:"errorMsg"`
	StartTime   string          `json:"startTime"`
	EndTime     string          `json:"endTime,omitempty"`
	RetryCount  int             `json:"retryCount"`
}

// === 分类 DTO ===

// CategoryCreateReq 创建分类请求
type CategoryCreateReq struct {
	Name     string `json:"name" binding:"required,max=100"`
	ParentID int64  `json:"parentId"`
	Sort     int    `json:"sort"`
	Status   int    `json:"status" binding:"omitempty,oneof=0 1"`
}

// CategoryUpdateReq 更新分类请求
type CategoryUpdateReq struct {
	ID       int64  `json:"id" binding:"required"`
	Name     string `json:"name" binding:"omitempty,max=100"`
	ParentID *int64 `json:"parentId"`
	Sort     int    `json:"sort"`
	Status   *int   `json:"status" binding:"omitempty,oneof=0 1"`
}

// CategoryResp 分类响应
type CategoryResp struct {
	ID         int64            `json:"id"`
	Name       string           `json:"name"`
	ParentID   int64            `json:"parentId"`
	Sort       int              `json:"sort"`
	Status     int              `json:"status"`
	Children   []*CategoryResp  `json:"children,omitempty"`
	CreateTime string           `json:"createTime"`
}

// WorkflowDefinition 工作流定义结构（JSON 解析用）
type WorkflowDefinition struct {
	Nodes    []WorkflowNodeDef `json:"nodes"`
	Edges    []WorkflowEdgeDef `json:"edges"`
	Settings map[string]any    `json:"settings,omitempty"`
}

// WorkflowNodeDef 工作流节点定义
type WorkflowNodeDef struct {
	ID       string         `json:"id"`
	Type     string         `json:"type"`
	Position PositionDef    `json:"position"`
	Data     NodeDataDef    `json:"data"`
}

// PositionDef 节点位置
type PositionDef struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// NodeDataDef 节点数据
type NodeDataDef struct {
	Label         string         `json:"label"`
	Config        map[string]any `json:"config"`
	ErrorStrategy string         `json:"errorStrategy,omitempty"`
	RetryCount    int            `json:"retryCount,omitempty"`
	Timeout       int            `json:"timeout,omitempty"`
}

// WorkflowEdgeDef 工作流边定义
type WorkflowEdgeDef struct {
	ID           string `json:"id"`
	Source       string `json:"source"`
	Target       string `json:"target"`
	SourceHandle string `json:"sourceHandle,omitempty"`
	TargetHandle string `json:"targetHandle,omitempty"`
	Animated     bool   `json:"animated,omitempty"`
	Label        string `json:"label,omitempty"`
}
