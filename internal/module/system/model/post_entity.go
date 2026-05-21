package model

import "management-backend/internal/model"

// Post 岗位实体（租户库）
type Post struct {
	model.BaseModel
	model.SortModel
	model.StatusModel
	model.RemarkModel
	PostCode string `gorm:"column:post_code;type:varchar(50);not null" json:"postCode"`
	PostName string `gorm:"column:post_name;type:varchar(50);not null" json:"postName"`
}

func (Post) TableName() string { return "sys_post" }
