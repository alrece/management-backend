package model

import "management-backend/pkg/page"

// ClientResp 客户端响应
type ClientResp struct {
	ID                   int64  `json:"id"`
	ClientID             string `json:"clientId"`
	ClientKey            string `json:"clientKey"`
	ClientName           string `json:"clientName"`
	GrantType            string `json:"grantType"`
	RedirectURI          string `json:"redirectUri"`
	AccessTokenValidity  int64  `json:"accessTokenValidity"`
	RefreshTokenValidity int64  `json:"refreshTokenValidity"`
	Status               int    `json:"status"`
	Remark               string `json:"remark"`
	CreateTime           string `json:"createTime"`
}

// ClientCreateReq 客户端创建请求
type ClientCreateReq struct {
	ClientID             string `json:"clientId" binding:"required"`
	ClientKey            string `json:"clientKey" binding:"required"`
	ClientName           string `json:"clientName" binding:"required"`
	GrantType            string `json:"grantType"`
	RedirectURI          string `json:"redirectUri"`
	AccessTokenValidity  int64  `json:"accessTokenValidity"`
	RefreshTokenValidity int64  `json:"refreshTokenValidity"`
	Status               int    `json:"status"`
	Remark               string `json:"remark"`
}

// ClientUpdateReq 客户端更新请求
type ClientUpdateReq struct {
	ID                   int64  `json:"id" binding:"required"`
	ClientKey            string `json:"clientKey"`
	ClientName           string `json:"clientName"`
	GrantType            string `json:"grantType"`
	RedirectURI          string `json:"redirectUri"`
	AccessTokenValidity  int64  `json:"accessTokenValidity"`
	RefreshTokenValidity int64  `json:"refreshTokenValidity"`
	Status               int    `json:"status"`
	Remark               string `json:"remark"`
}

// ClientPageReq 客户端分页查询
type ClientPageReq struct {
	page.Req
	ClientName string `form:"clientName"`
	Status     *int   `form:"status"`
}
