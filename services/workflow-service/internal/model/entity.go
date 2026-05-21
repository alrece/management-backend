package model

import "time"

// WfCategory 工作流分类
type WfCategory struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	TenantID  int64     `gorm:"not null;index" json:"tenantId"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	ParentID  int64     `gorm:"default:0;index" json:"parentId"`
	Sort      int       `gorm:"default:0" json:"sort"`
	Status    int8      `gorm:"default:0" json:"status"` // 0正常 1停用
	Creator   int64     `gorm:"" json:"creator"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	Updater   int64     `gorm:"" json:"updater"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
	Deleted   int8      `gorm:"default:0" json:"-"`
}

func (WfCategory) TableName() string { return "wf_category" }

// WfWorkflow 工作流定义
type WfWorkflow struct {
	ID             int64     `gorm:"primaryKey" json:"id"`
	TenantID       int64     `gorm:"not null;index" json:"tenantId"`
	Name           string    `gorm:"size:200;not null" json:"name"`
	CategoryID     int64     `gorm:"index" json:"categoryId"`
	N8nWorkflowID  string    `gorm:"size:100" json:"n8nWorkflowId"`
	ParamsSchema   string    `gorm:"type:text" json:"paramsSchema"` // JSON Schema
	Status         string    `gorm:"size:20;not null;default:DRAFT" json:"status"` // DRAFT/ACTIVE/INACTIVE/DESYNCED
	Creator        int64     `gorm:"" json:"creator"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"createdAt"`
	Updater        int64     `gorm:"" json:"updater"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
	Deleted        int8      `gorm:"default:0" json:"-"`
}

func (WfWorkflow) TableName() string { return "wf_workflow" }

// WfInstance 工作流执行实例
type WfInstance struct {
	ID              int64      `gorm:"primaryKey" json:"id"`
	TenantID        int64      `gorm:"not null;index" json:"tenantId"`
	WorkflowID      int64      `gorm:"not null;index" json:"workflowId"`
	N8nExecutionID  string     `gorm:"size:100" json:"n8nExecutionId"`
	Status          string     `gorm:"size:20;not null;default:RUNNING" json:"status"` // RUNNING/SUCCESS/FAILED/CANCELLED
	Variables       string     `gorm:"type:text" json:"variables"`  // 执行入参 JSON
	Result          string     `gorm:"type:text" json:"result"`     // 执行结果
	DurationMs      int64      `gorm:"" json:"durationMs"`          // 执行耗时（毫秒）
	ErrorMsg        string     `gorm:"type:text" json:"errorMsg"`
	StartedAt       time.Time  `gorm:"" json:"startedAt"`
	FinishedAt      *time.Time `gorm:"" json:"finishedAt"`
	Creator         int64      `gorm:"" json:"creator"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"createdAt"`
}

func (WfInstance) TableName() string { return "wf_instance" }

// WfVariable 工作流变量
type WfVariable struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	TenantID   int64     `gorm:"not null;index" json:"tenantId"`
	WorkflowID int64     `gorm:"not null;index" json:"workflowId"`
	Name       string    `gorm:"size:100;not null" json:"name"`
	Value      string    `gorm:"type:text" json:"value"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (WfVariable) TableName() string { return "wf_variable" }
