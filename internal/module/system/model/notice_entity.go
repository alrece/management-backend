package model

import (
	"time"

	"management-backend/internal/model"
)

// Notice 通知公告（租户库）
type Notice struct {
	model.BaseEntity
	model.StatusModel
	model.RemarkModel
	NoticeTitle string `gorm:"column:notice_title;type:varchar(200);not null" json:"noticeTitle"`
	NoticeType  int    `gorm:"column:notice_type;not null;default:1" json:"noticeType"`
	Content     string `gorm:"column:content;type:text" json:"content"`
}

func (Notice) TableName() string { return "sys_notice" }

// SysOperLog 操作日志（租户库）
type SysOperLog struct {
	ID           int64     `gorm:"primaryKey;autoIncrement:false" json:"id"`
	TenantID     int64     `gorm:"column:tenant_id;index;default:0" json:"tenantId"`
	Title        string    `gorm:"column:title;type:varchar(100)" json:"title"`
	BusinessType int       `gorm:"column:business_type;default:0" json:"businessType"`
	Method       string    `gorm:"column:method;type:varchar(200)" json:"method"`
	RequestURL   string    `gorm:"column:request_url;type:varchar(500)" json:"requestUrl"`
	OperIP       string    `gorm:"column:oper_ip;type:varchar(50)" json:"operIp"`
	OperUserID   int64     `gorm:"column:oper_user_id" json:"operUserId"`
	OperName     string    `gorm:"column:oper_name;type:varchar(50)" json:"operName"`
	RequestID    string    `gorm:"column:request_id;type:varchar(50);index" json:"requestId"`
	Status       int       `gorm:"column:status;default:0" json:"status"`
	ErrorMsg     string    `gorm:"column:error_msg;type:text" json:"errorMsg"`
	OperTime     time.Time `gorm:"column:oper_time" json:"operTime"`
}

func (SysOperLog) TableName() string { return "sys_oper_log" }

// SysLoginLog 登录日志（默认库）
type SysLoginLog struct {
	ID            int64     `gorm:"primaryKey;autoIncrement:false" json:"id"`
	Username      string    `gorm:"column:username;type:varchar(50);index" json:"username"`
	LoginIP       string    `gorm:"column:login_ip;type:varchar(50)" json:"loginIp"`
	LoginLocation string    `gorm:"column:login_location;type:varchar(100)" json:"loginLocation"`
	Browser       string    `gorm:"column:browser;type:varchar(50)" json:"browser"`
	OS            string    `gorm:"column:os;type:varchar(50)" json:"os"`
	Status        int       `gorm:"column:status;default:0" json:"status"`
	Msg           string    `gorm:"column:msg;type:varchar(200)" json:"msg"`
	LoginTime     time.Time `gorm:"column:login_time" json:"loginTime"`
}

func (SysLoginLog) TableName() string { return "sys_login_log" }
