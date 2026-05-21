package model

import "management-backend/internal/model"

// Menu 菜单实体（默认库，不受租户过滤）
type Menu struct {
	model.BaseModel
	model.SortModel
	model.StatusModel
	ParentID   int64  `gorm:"column:parent_id;not null;default:0" json:"parentId"`
	MenuName   string `gorm:"column:menu_name;type:varchar(50);not null" json:"menuName"`
	MenuType   int    `gorm:"column:menu_type;not null" json:"menuType"`
	Path       string `gorm:"column:path;type:varchar(200);default:''" json:"path"`
	Component  string `gorm:"column:component;type:varchar(200);default:''" json:"component"`
	Permission string `gorm:"column:permission;type:varchar(100);default:''" json:"permission"`
	Icon       string `gorm:"column:icon;type:varchar(50);default:''" json:"icon"`
	Visible    int    `gorm:"column:visible;not null;default:0" json:"visible"`
}

func (Menu) TableName() string { return "sys_menu" }
