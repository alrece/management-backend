package model

import "time"

// JobTask 定时任务实体
type JobTask struct {
	ID         int64     `json:"id" gorm:"primaryKey"`
	TenantID   int64     `json:"tenantId" gorm:"column:tenant_id;not null;default:0;index"`
	Name       string    `json:"name" gorm:"size:100;not null"`
	Handler    string    `json:"handler" gorm:"size:100;not null"`
	CronExpr   string    `json:"cronExpr" gorm:"column:cron_expr;size:50;not null"`
	Params     string    `json:"params" gorm:"size:2000;default:''"`
	Status     int8      `json:"status" gorm:"not null;default:0"` // 0正常 1暂停
	Remark     string    `json:"remark" gorm:"size:500;default:''"`
	Creator    *int64    `json:"creator" gorm:""`
	CreateTime time.Time `json:"createTime" gorm:"column:create_time;autoCreateTime"`
	Updater    *int64    `json:"updater" gorm:""`
	UpdateTime time.Time `json:"updateTime" gorm:"column:update_time;autoUpdateTime"`
	Deleted    int8      `json:"-" gorm:"not null;default:0"`
}

func (JobTask) TableName() string { return "job_task" }

// JobExecutionLog 任务执行日志
type JobExecutionLog struct {
	ID         int64     `json:"id" gorm:"primaryKey"`
	TenantID   int64     `json:"tenantId" gorm:"column:tenant_id;not null;default:0;index"`
	TaskID     int64     `json:"taskId" gorm:"column:task_id;not null;index"`
	TaskName   string    `json:"taskName" gorm:"column:task_name;size:100"`
	TriggerType int8     `json:"triggerType" gorm:"column:trigger_type;not null"` // 1定时 2手动
	Status     int8      `json:"status" gorm:"not null;default:0"` // 0执行中 1成功 2失败 3超时
	DurationMs int64     `json:"durationMs" gorm:"column:duration_ms;default:0"`
	Result     string    `json:"result" gorm:"size:2000;default:''"`
	Error      string    `json:"error" gorm:"size:2000;default:''"`
	StartTime  time.Time `json:"startTime" gorm:"column:start_time"`
	EndTime    *time.Time `json:"endTime" gorm:"column:end_time"`
}

func (JobExecutionLog) TableName() string { return "job_execution_log" }

const (
	TaskStatusNormal = 0
	TaskStatusPaused = 1
)

const (
	ExecStatusRunning = 0
	ExecStatusSuccess = 1
	ExecStatusFailed  = 2
	ExecStatusTimeout = 3
)

const (
	TriggerTypeCron   = 1
	TriggerTypeManual = 2
)
