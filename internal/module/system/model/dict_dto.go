package model

import "management-backend/pkg/page"

// DictTypeCreateReq 创建字典类型请求
type DictTypeCreateReq struct {
	DictName string `json:"dictName" binding:"required,max=100"`
	DictType string `json:"dictType" binding:"required,max=100"`
	Status   int    `json:"status" binding:"omitempty,oneof=0 1"`
	Remark   string `json:"remark" binding:"omitempty,max=500"`
}

// DictTypeUpdateReq 更新字典类型请求
type DictTypeUpdateReq struct {
	ID       int64  `json:"id" binding:"required"`
	DictName string `json:"dictName" binding:"required,max=100"`
	DictType string `json:"dictType" binding:"required,max=100"`
	Status   int    `json:"status" binding:"omitempty,oneof=0 1"`
	Remark   string `json:"remark" binding:"omitempty,max=500"`
}

// DictTypePageReq 字典类型分页查询
type DictTypePageReq struct {
	page.Req
	DictName string `form:"dictName" json:"dictName"`
	DictType string `form:"dictType" json:"dictType"`
	Status   *int   `form:"status" json:"status"`
}

// DictTypeResp 字典类型响应
type DictTypeResp struct {
	ID         int64  `json:"id"`
	DictName   string `json:"dictName"`
	DictType   string `json:"dictType"`
	Status     int    `json:"status"`
	Remark     string `json:"remark"`
	CreateTime string `json:"createTime"`
}

// DictDataCreateReq 创建字典数据请求
type DictDataCreateReq struct {
	DictType string `json:"dictType" binding:"required,max=100"`
	Label    string `json:"label" binding:"required,max=100"`
	Value    string `json:"value" binding:"required,max=100"`
	Sort     int    `json:"sort"`
	Status   int    `json:"status" binding:"omitempty,oneof=0 1"`
	Remark   string `json:"remark" binding:"omitempty,max=500"`
}

// DictDataUpdateReq 更新字典数据请求
type DictDataUpdateReq struct {
	ID       int64  `json:"id" binding:"required"`
	DictType string `json:"dictType" binding:"required,max=100"`
	Label    string `json:"label" binding:"required,max=100"`
	Value    string `json:"value" binding:"required,max=100"`
	Sort     int    `json:"sort"`
	Status   int    `json:"status" binding:"omitempty,oneof=0 1"`
	Remark   string `json:"remark" binding:"omitempty,max=500"`
}

// DictDataPageReq 字典数据分页查询
type DictDataPageReq struct {
	page.Req
	DictType string `form:"dictType" json:"dictType"`
	Label    string `form:"label" json:"label"`
	Status   *int   `form:"status" json:"status"`
}

// DictDataResp 字典数据响应
type DictDataResp struct {
	ID         int64  `json:"id"`
	DictType   string `json:"dictType"`
	Label      string `json:"label"`
	Value      string `json:"value"`
	Sort       int    `json:"sort"`
	Status     int    `json:"status"`
	Remark     string `json:"remark"`
	CreateTime string `json:"createTime"`
}
