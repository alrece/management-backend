package model

import "management-backend/internal/model"

// Dept 部门实体（租户库）
type Dept struct {
	model.BaseModel
	model.SortModel
	model.StatusModel
	ParentID  int64  `gorm:"column:parent_id;not null;default:0" json:"parentId"`
	Ancestors string `gorm:"column:ancestors;type:varchar(200);default:''" json:"ancestors"`
	DeptName  string `gorm:"column:dept_name;type:varchar(50);not null" json:"deptName"`
	Leader    string `gorm:"column:leader;type:varchar(30);default:''" json:"leader"`
}

func (Dept) TableName() string { return "sys_dept" }
