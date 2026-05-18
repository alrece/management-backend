package model

import "management-backend/internal/model"

// User 用户实体（对应 sys_user 表）
type User struct {
	model.BaseEntity
	model.StatusModel
	model.RemarkModel
	Username string `gorm:"type:varchar(30);uniqueIndex" json:"username"`
	Password string `gorm:"type:varchar(100)" json:"-"`
	Nickname string `gorm:"type:varchar(30)" json:"nickname"`
	Email    string `gorm:"type:varchar(50)" json:"email"`
	Mobile   string `gorm:"type:varchar(20)" json:"mobile"`
	Sex      int    `gorm:"default:0" json:"sex"`    // 0 未知 1 男 2 女
	Avatar   string `gorm:"type:varchar(255)" json:"avatar"`
	DeptID   int64  `gorm:"index;default:0" json:"deptId"`
}

func (User) TableName() string {
	return "sys_user"
}
