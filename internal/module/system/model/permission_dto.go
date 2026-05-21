package model

import "management-backend/pkg/page"

// ====== 角色 DTO ======

type RoleCreateReq struct {
	RoleName  string `json:"roleName" binding:"required"`
	RoleCode  string `json:"roleCode" binding:"required"`
	Sort      int    `json:"sort"`
	DataScope int    `json:"dataScope" binding:"min=1,max=5"`
	Status    int    `json:"status"`
	Remark    string `json:"remark"`
}

type RoleUpdateReq struct {
	ID               int64  `json:"id" binding:"required"`
	RoleName         string `json:"roleName" binding:"required"`
	RoleCode         string `json:"roleCode" binding:"required"`
	Sort             int    `json:"sort"`
	DataScope        int    `json:"dataScope" binding:"min=1,max=5"`
	DataScopeDeptIDs string `json:"dataScopeDeptIds"`
	Status           int    `json:"status"`
	Remark           string `json:"remark"`
}

type RolePageReq struct {
	page.Req
	RoleName string `form:"roleName"`
	RoleCode string `form:"roleCode"`
	Status   *int   `form:"status"`
}

type RoleResp struct {
	ID               int64  `json:"id"`
	RoleName         string `json:"roleName"`
	RoleCode         string `json:"roleCode"`
	Sort             int    `json:"sort"`
	DataScope        int    `json:"dataScope"`
	DataScopeDeptIDs string `json:"dataScopeDeptIds"`
	Status           int    `json:"status"`
	Remark           string `json:"remark"`
	CreateTime       string `json:"createTime"`
}

type RoleMenuAssignReq struct {
	RoleID  int64   `json:"roleId" binding:"required"`
	MenuIDs []int64 `json:"menuIds" binding:"required"`
}

// ====== 菜单 DTO ======

type MenuCreateReq struct {
	ParentID   int64  `json:"parentId"`
	MenuName   string `json:"menuName" binding:"required"`
	MenuType   int    `json:"menuType" binding:"required,min=1,max=3"`
	Path       string `json:"path"`
	Component  string `json:"component"`
	Permission string `json:"permission"`
	Icon       string `json:"icon"`
	Sort       int    `json:"sort"`
	Visible    int    `json:"visible"`
	Status     int    `json:"status"`
}

type MenuUpdateReq struct {
	ID         int64  `json:"id" binding:"required"`
	ParentID   int64  `json:"parentId"`
	MenuName   string `json:"menuName" binding:"required"`
	MenuType   int    `json:"menuType" binding:"min=1,max=3"`
	Path       string `json:"path"`
	Component  string `json:"component"`
	Permission string `json:"permission"`
	Icon       string `json:"icon"`
	Sort       int    `json:"sort"`
	Visible    int    `json:"visible"`
	Status     int    `json:"status"`
}

type MenuTreeResp struct {
	ID         int64          `json:"id"`
	ParentID   int64          `json:"parentId"`
	MenuName   string         `json:"menuName"`
	MenuType   int            `json:"menuType"`
	Path       string         `json:"path"`
	Component  string         `json:"component"`
	Permission string         `json:"permission"`
	Icon       string         `json:"icon"`
	Sort       int            `json:"sort"`
	Visible    int            `json:"visible"`
	Status     int            `json:"status"`
	Children   []MenuTreeResp `json:"children,omitempty"`
}

// ====== 部门 DTO ======

type DeptCreateReq struct {
	ParentID int64  `json:"parentId"`
	DeptName string `json:"deptName" binding:"required"`
	Sort     int    `json:"sort"`
	Leader   string `json:"leader"`
	Status   int    `json:"status"`
}

type DeptUpdateReq struct {
	ID       int64  `json:"id" binding:"required"`
	ParentID int64  `json:"parentId"`
	DeptName string `json:"deptName" binding:"required"`
	Sort     int    `json:"sort"`
	Leader   string `json:"leader"`
	Status   int    `json:"status"`
}

type DeptTreeResp struct {
	ID        int64         `json:"id"`
	ParentID  int64         `json:"parentId"`
	DeptName  string        `json:"deptName"`
	Ancestors string        `json:"ancestors"`
	Sort      int           `json:"sort"`
	Leader    string        `json:"leader"`
	Status    int           `json:"status"`
	Children  []DeptTreeResp `json:"children,omitempty"`
}

// ====== 岗位 DTO ======

type PostCreateReq struct {
	PostCode string `json:"postCode" binding:"required"`
	PostName string `json:"postName" binding:"required"`
	Sort     int    `json:"sort"`
	Status   int    `json:"status"`
	Remark   string `json:"remark"`
}

type PostUpdateReq struct {
	ID       int64  `json:"id" binding:"required"`
	PostCode string `json:"postCode" binding:"required"`
	PostName string `json:"postName" binding:"required"`
	Sort     int    `json:"sort"`
	Status   int    `json:"status"`
	Remark   string `json:"remark"`
}

type PostPageReq struct {
	page.Req
	PostCode string `form:"postCode"`
	PostName string `form:"postName"`
	Status   *int   `form:"status"`
}

type PostResp struct {
	ID         int64  `json:"id"`
	PostCode   string `json:"postCode"`
	PostName   string `json:"postName"`
	Sort       int    `json:"sort"`
	Status     int    `json:"status"`
	Remark     string `json:"remark"`
	CreateTime string `json:"createTime"`
}

// ====== 权限信息 ======

type PermissionInfoResp struct {
	User        UserResp       `json:"user"`
	Roles       []string       `json:"roles"`
	Permissions []string       `json:"permissions"`
	Menus       []MenuTreeResp `json:"menus"`
}
