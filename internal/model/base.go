package model

import (
	"time"
)

// BaseEntity 基础实体（含租户字段）
type BaseEntity struct {
	ID         int64     `gorm:"primaryKey;autoIncrement:false" json:"id"`
	TenantID   int64     `gorm:"index;default:0" json:"tenantId"`
	Creator    int64     `json:"creator"`
	CreateTime time.Time `gorm:"autoCreateTime" json:"createTime"`
	Updater    int64     `json:"updater"`
	UpdateTime time.Time `gorm:"autoUpdateTime" json:"updateTime"`
	Deleted    int       `gorm:"index;default:0" json:"-"`
	CreateDept int64     `json:"createDept"`
}

// BaseModel 无租户的基础实体（用于 sys_menu 等租户排除表）
type BaseModel struct {
	ID         int64     `gorm:"primaryKey;autoIncrement:false" json:"id"`
	Creator    int64     `json:"creator"`
	CreateTime time.Time `gorm:"autoCreateTime" json:"createTime"`
	Updater    int64     `json:"updater"`
	UpdateTime time.Time `gorm:"autoUpdateTime" json:"updateTime"`
	Deleted    int       `gorm:"index;default:0" json:"-"`
}

// SortModel 排序字段
type SortModel struct {
	Sort int `gorm:"default:0" json:"sort"`
}

// StatusModel 状态字段（0 正常 1 停用）
type StatusModel struct {
	Status int `gorm:"default:0" json:"status"`
}

// RemarkModel 备注字段
type RemarkModel struct {
	Remark string `gorm:"type:varchar(500)" json:"remark"`
}
