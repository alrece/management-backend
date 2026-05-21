package model

import "management-backend/pkg/page"

// ParamCreateReq 创建参数请求
type ParamCreateReq struct {
	ParamKey   string `json:"paramKey" binding:"required,max=100"`
	ParamValue string `json:"paramValue" binding:"required,max=500"`
	ParamType  int    `json:"paramType" binding:"omitempty,oneof=1 2"`
	Status     int    `json:"status" binding:"omitempty,oneof=0 1"`
	Remark     string `json:"remark" binding:"omitempty,max=500"`
}

// ParamUpdateReq 更新参数请求
type ParamUpdateReq struct {
	ID         int64  `json:"id" binding:"required"`
	ParamKey   string `json:"paramKey" binding:"required,max=100"`
	ParamValue string `json:"paramValue" binding:"required,max=500"`
	ParamType  int    `json:"paramType" binding:"omitempty,oneof=1 2"`
	Status     int    `json:"status" binding:"omitempty,oneof=0 1"`
	Remark     string `json:"remark" binding:"omitempty,max=500"`
}

// ParamPageReq 参数分页查询
type ParamPageReq struct {
	page.Req
	ParamKey  string `form:"paramKey" json:"paramKey"`
	ParamType *int   `form:"paramType" json:"paramType"`
	Status    *int   `form:"status" json:"status"`
}

// ParamResp 参数响应
type ParamResp struct {
	ID         int64  `json:"id"`
	ParamKey   string `json:"paramKey"`
	ParamValue string `json:"paramValue"`
	ParamType  int    `json:"paramType"`
	Status     int    `json:"status"`
	Remark     string `json:"remark"`
	CreateTime string `json:"createTime"`
}
