package model

import "management-backend/pkg/page"

// UserCreateReq 创建用户请求
type UserCreateReq struct {
	Username string `json:"username" binding:"required,min=2,max=30"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Nickname string `json:"nickname" binding:"omitempty,max=30"`
	Email    string `json:"email" binding:"omitempty,email,max=50"`
	Mobile   string `json:"mobile" binding:"omitempty,max=20"`
	Sex      int    `json:"sex" binding:"omitempty,oneof=0 1 2"`
	Status   int    `json:"status" binding:"omitempty,oneof=0 1"`
	DeptID   int64  `json:"deptId"`
	Remark   string `json:"remark" binding:"omitempty,max=500"`
}

// UserUpdateReq 更新用户请求
type UserUpdateReq struct {
	ID       int64  `json:"id" binding:"required"`
	Nickname string `json:"nickname" binding:"omitempty,max=30"`
	Email    string `json:"email" binding:"omitempty,email,max=50"`
	Mobile   string `json:"mobile" binding:"omitempty,max=20"`
	Sex      int    `json:"sex" binding:"omitempty,oneof=0 1 2"`
	Status   int    `json:"status" binding:"omitempty,oneof=0 1"`
	DeptID   int64  `json:"deptId"`
	Remark   string `json:"remark" binding:"omitempty,max=500"`
}

// UserPageReq 分页查询请求
type UserPageReq struct {
	page.Req
	Username string `form:"username" json:"username"`
	Status   *int   `form:"status" json:"status"`
	Mobile   string `form:"mobile" json:"mobile"`
}

// UserResp 用户响应
type UserResp struct {
	ID         int64  `json:"id"`
	Username   string `json:"username"`
	Nickname   string `json:"nickname"`
	Email      string `json:"email"`
	Mobile     string `json:"mobile"`
	Sex        int    `json:"sex"`
	Avatar     string `json:"avatar"`
	Status     int    `json:"status"`
	DeptID     int64  `json:"deptId"`
	TenantID   int64  `json:"tenantId"`
	CreateTime string `json:"createTime"`
	Remark     string `json:"remark"`
}
