package model

import (
	"management-backend/internal/model"
	"time"
)

// InitStatus 租户初始化状态
const (
	InitStatusPending = 0 // 待初始化
	InitStatusIniting = 1 // 初始化中
	InitStatusReady   = 2 // 已就绪
	InitStatusFailed  = 3 // 失败
)

// Tenant 租户实体（sys_tenant 表，存于默认库，不包含 tenant_id）
type Tenant struct {
	model.BaseModel
	model.StatusModel
	model.RemarkModel
	TenantName    string     `gorm:"column:tenant_name;type:varchar(100);not null" json:"tenantName"`
	ContactName   string     `gorm:"column:contact_name;type:varchar(30)" json:"contactName"`
	ContactMobile string     `gorm:"column:contact_mobile;type:varchar(20)" json:"contactMobile"`
	DBName        string     `gorm:"column:db_name;type:varchar(100)" json:"dbName"`
	ExpireTime    *time.Time `gorm:"column:expire_time" json:"expireTime"`
	InitStatus    int        `gorm:"column:init_status;not null;default:0" json:"initStatus"`
}

func (Tenant) TableName() string {
	return "sys_tenant"
}
