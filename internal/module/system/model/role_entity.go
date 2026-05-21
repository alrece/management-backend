package model

import (
	"time"

	"management-backend/internal/model"

	"gorm.io/gorm"
)

// Role 角色实体（租户库）
type Role struct {
	model.BaseModel
	model.SortModel
	model.StatusModel
	model.RemarkModel
	RoleName         string `gorm:"column:role_name;type:varchar(50);not null" json:"roleName"`
	RoleCode         string `gorm:"column:role_code;type:varchar(50);not null" json:"roleCode"`
	DataScope        int    `gorm:"column:data_scope;not null;default:1" json:"dataScope"`
	DataScopeDeptIDs string `gorm:"column:data_scope_dept_ids;type:varchar(500);default:''" json:"dataScopeDeptIds"`
}

func (Role) TableName() string { return "sys_role" }

// RoleMenu 角色菜单关联（租户库）
type RoleMenu struct {
	ID         int64          `gorm:"primaryKey;autoIncrement:false" json:"id"`
	RoleID     int64          `gorm:"not null;index" json:"roleId"`
	MenuID     int64          `gorm:"not null;index" json:"menuId"`
	Creator    int64          `json:"creator"`
	CreateTime time.Time      `gorm:"autoCreateTime" json:"createTime"`
	Updater    int64          `json:"updater"`
	UpdateTime time.Time      `gorm:"autoUpdateTime" json:"updateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (RoleMenu) TableName() string { return "sys_role_menu" }
