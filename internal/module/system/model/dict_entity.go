package model

import "management-backend/internal/model"

// DictType 字典类型（租户库）
type DictType struct {
	model.BaseEntity
	model.StatusModel
	model.RemarkModel
	DictName string `gorm:"column:dict_name;type:varchar(100);not null" json:"dictName"`
	DictType string `gorm:"column:dict_type;type:varchar(100);not null;uniqueIndex" json:"dictType"`
}

func (DictType) TableName() string { return "sys_dict_type" }

// DictData 字典数据（租户库）
type DictData struct {
	model.BaseEntity
	model.SortModel
	model.StatusModel
	model.RemarkModel
	DictType string `gorm:"column:dict_type;type:varchar(100);not null;index" json:"dictType"`
	Label    string `gorm:"column:label;type:varchar(100);not null" json:"label"`
	Value    string `gorm:"column:value;type:varchar(100);not null" json:"value"`
}

func (DictData) TableName() string { return "sys_dict_data" }
