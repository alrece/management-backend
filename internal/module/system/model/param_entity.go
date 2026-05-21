package model

import "management-backend/internal/model"

// SysParam 系统参数（租户库）
type SysParam struct {
	model.BaseEntity
	model.StatusModel
	model.RemarkModel
	ParamKey   string `gorm:"column:param_key;type:varchar(100);not null;uniqueIndex" json:"paramKey"`
	ParamValue string `gorm:"column:param_value;type:varchar(500);not null" json:"paramValue"`
	ParamType  int    `gorm:"column:param_type;default:1" json:"paramType"` // 1 系统 2 用户自定义
}

func (SysParam) TableName() string { return "sys_param" }
