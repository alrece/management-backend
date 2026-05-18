package model

import "management-backend/pkg/page"

// TenantCreateReq 创建租户请求
type TenantCreateReq struct {
	TenantName    string `json:"tenantName" binding:"required,max=100"`
	ContactName   string `json:"contactName" binding:"omitempty,max=30"`
	ContactMobile string `json:"contactMobile" binding:"omitempty,max=20"`
	ExpireTime    string `json:"expireTime" binding:"omitempty"`
	Status        int    `json:"status" binding:"omitempty,oneof=0 1"`
	Remark        string `json:"remark" binding:"omitempty,max=500"`
}

// TenantUpdateReq 更新租户请求
type TenantUpdateReq struct {
	ID            int64  `json:"id" binding:"required"`
	TenantName    string `json:"tenantName" binding:"omitempty,max=100"`
	ContactName   string `json:"contactName" binding:"omitempty,max=30"`
	ContactMobile string `json:"contactMobile" binding:"omitempty,max=20"`
	ExpireTime    string `json:"expireTime" binding:"omitempty"`
	Status        int    `json:"status" binding:"omitempty,oneof=0 1"`
	Remark        string `json:"remark" binding:"omitempty,max=500"`
}

// TenantPageReq 分页查询请求
type TenantPageReq struct {
	page.Req
	TenantName string `form:"tenantName" json:"tenantName"`
	Status     *int   `form:"status" json:"status"`
}

// TenantResp 租户响应
type TenantResp struct {
	ID            int64  `json:"id"`
	TenantName    string `json:"tenantName"`
	ContactName   string `json:"contactName"`
	ContactMobile string `json:"contactMobile"`
	DBName        string `json:"dbName"`
	Status        int    `json:"status"`
	InitStatus    int    `json:"initStatus"`
	ExpireTime    string `json:"expireTime"`
	Remark        string `json:"remark"`
	CreateTime    string `json:"createTime"`
}
