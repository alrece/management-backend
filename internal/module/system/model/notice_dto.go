package model

import "management-backend/pkg/page"

// NoticeCreateReq 创建通知请求
type NoticeCreateReq struct {
	NoticeTitle string `json:"noticeTitle" binding:"required,max=200"`
	NoticeType  int    `json:"noticeType" binding:"required,oneof=1 2 3"` // 1 通知 2 公告 3 提醒
	Content     string `json:"content"`
	Status      int    `json:"status" binding:"omitempty,oneof=0 1"`
	Remark      string `json:"remark" binding:"omitempty,max=500"`
}

// NoticeUpdateReq 更新通知请求
type NoticeUpdateReq struct {
	ID          int64  `json:"id" binding:"required"`
	NoticeTitle string `json:"noticeTitle" binding:"required,max=200"`
	NoticeType  int    `json:"noticeType" binding:"required,oneof=1 2 3"`
	Content     string `json:"content"`
	Status      int    `json:"status" binding:"omitempty,oneof=0 1"`
	Remark      string `json:"remark" binding:"omitempty,max=500"`
}

// NoticePageReq 通知分页查询
type NoticePageReq struct {
	page.Req
	NoticeTitle string `form:"noticeTitle" json:"noticeTitle"`
	NoticeType  *int   `form:"noticeType" json:"noticeType"`
	Status      *int   `form:"status" json:"status"`
}

// NoticeResp 通知响应
type NoticeResp struct {
	ID          int64  `json:"id"`
	NoticeTitle string `json:"noticeTitle"`
	NoticeType  int    `json:"noticeType"`
	Content     string `json:"content"`
	Status      int    `json:"status"`
	Remark      string `json:"remark"`
	Creator     int64  `json:"creator"`
	CreateTime  string `json:"createTime"`
}

// OperLogPageReq 操作日志分页查询
type OperLogPageReq struct {
	page.Req
	Title        string `form:"title" json:"title"`
	BusinessType *int   `form:"businessType" json:"businessType"`
	Status       *int   `form:"status" json:"status"`
	OperName     string `form:"operName" json:"operName"`
	BeginTime    string `form:"beginTime" json:"beginTime"`
	EndTime      string `form:"endTime" json:"endTime"`
}

// OperLogResp 操作日志响应
type OperLogResp struct {
	ID           int64  `json:"id"`
	TenantID     int64  `json:"tenantId"`
	Title        string `json:"title"`
	BusinessType int    `json:"businessType"`
	Method       string `json:"method"`
	RequestURL   string `json:"requestUrl"`
	OperIP       string `json:"operIp"`
	OperUserID   int64  `json:"operUserId"`
	OperName     string `json:"operName"`
	RequestID    string `json:"requestId"`
	Status       int    `json:"status"`
	ErrorMsg     string `json:"errorMsg"`
	OperTime     string `json:"operTime"`
}

// LoginLogPageReq 登录日志分页查询
type LoginLogPageReq struct {
	page.Req
	Username  string `form:"username" json:"username"`
	Status    *int   `form:"status" json:"status"`
	BeginTime string `form:"beginTime" json:"beginTime"`
	EndTime   string `form:"endTime" json:"endTime"`
}

// LoginLogResp 登录日志响应
type LoginLogResp struct {
	ID            int64  `json:"id"`
	Username      string `json:"username"`
	LoginIP       string `json:"loginIp"`
	LoginLocation string `json:"loginLocation"`
	Browser       string `json:"browser"`
	OS            string `json:"os"`
	Status        int    `json:"status"`
	Msg           string `json:"msg"`
	LoginTime     string `json:"loginTime"`
}
