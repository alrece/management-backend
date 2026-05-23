package model

import (
	"encoding/json"
	"time"

	"management-backend/internal/model"
)

// Workflow 工作流定义
type Workflow struct {
	model.BaseEntity
	model.StatusModel
	model.RemarkModel
	Name       string          `gorm:"type:varchar(100);not null" json:"name"`
	CategoryID int64           `gorm:"index;default:0" json:"categoryId"`
	// 工作流定义（节点+边+设置），JSONB 最大 1MB，节点数上限 100
	Definition json.RawMessage `gorm:"type:jsonb" json:"definition"`
	// 触发器类型：manual / cron / webhook
	TriggerType string         `gorm:"type:varchar(20);default:manual" json:"triggerType"`
	// 触发器配置（cron 表达式或 webhook 路径等）
	TriggerConfig json.RawMessage `gorm:"type:jsonb" json:"triggerConfig"`
	// 工作流版本号
	Version    int             `gorm:"default:1" json:"version"`
}

func (Workflow) TableName() string {
	return "wf_workflow"
}

// Execution 工作流执行实例
type Execution struct {
	ID         int64     `gorm:"primaryKey;autoIncrement:false" json:"id"`
	TenantID   int64     `gorm:"index;default:0" json:"tenantId"`
	WorkflowID int64     `gorm:"index;not null" json:"workflowId"`
	// 执行状态：running / success / failed / canceled
	Status     string    `gorm:"type:varchar(20);index;default:running" json:"status"`
	TriggerType string   `gorm:"type:varchar(20)" json:"triggerType"`
	StartTime  time.Time `gorm:"autoCreateTime" json:"startTime"`
	EndTime    *time.Time `json:"endTime,omitempty"`
	// 输入参数
	Input      json.RawMessage `gorm:"type:jsonb" json:"input"`
	// 执行摘要（各节点输出快照）
	Output     json.RawMessage `gorm:"type:jsonb" json:"output"`
	ErrorMsg   string          `gorm:"type:text" json:"errorMsg"`
	Creator    int64           `json:"creator"`
}

func (Execution) TableName() string {
	return "wf_execution"
}

// NodeLog 节点执行日志
type NodeLog struct {
	ID          int64           `gorm:"primaryKey;autoIncrement:false" json:"id"`
	TenantID    int64           `gorm:"index;default:0" json:"tenantId"`
	ExecutionID int64           `gorm:"index;not null" json:"executionId"`
	WorkflowID  int64           `gorm:"index;not null" json:"workflowId"`
	NodeID      string          `gorm:"type:varchar(50);not null" json:"nodeId"`
	NodeType    string          `gorm:"type:varchar(30);not null" json:"nodeType"`
	// 节点状态：running / success / failed / skipped
	Status      string          `gorm:"type:varchar(20);default:running" json:"status"`
	Input       json.RawMessage `gorm:"type:jsonb" json:"input"`
	Output      json.RawMessage `gorm:"type:jsonb" json:"output"`
	ErrorMsg    string          `gorm:"type:text" json:"errorMsg"`
	StartTime   time.Time       `json:"startTime"`
	EndTime     *time.Time      `json:"endTime,omitempty"`
	RetryCount  int             `gorm:"default:0" json:"retryCount"`
}

func (NodeLog) TableName() string {
	return "wf_node_log"
}

// Trigger 触发器配置
type Trigger struct {
	model.BaseEntity
	WorkflowID   int64           `gorm:"index;not null" json:"workflowId"`
	TriggerType  string          `gorm:"type:varchar(20);not null" json:"triggerType"`
	// cron 表达式或 webhook 路径
	Config       json.RawMessage `gorm:"type:jsonb" json:"config"`
	// cron entry ID（运行时使用）
	CronEntryID  int             `gorm:"default:0" json:"cronEntryId"`
	// 激活状态
	Active       bool            `gorm:"default:false" json:"active"`
}

func (Trigger) TableName() string {
	return "wf_trigger"
}

// Category 工作流分类
type Category struct {
	model.BaseEntity
	model.SortModel
	model.StatusModel
	Name     string `gorm:"type:varchar(100);not null" json:"name"`
	ParentID int64  `gorm:"index;default:0" json:"parentId"`
}

func (Category) TableName() string {
	return "wf_category"
}
