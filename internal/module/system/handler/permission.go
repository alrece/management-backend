package handler

import (
	smodel "management-backend/internal/module/system/model"
	"management-backend/internal/module/system/repository"
	"management-backend/internal/module/system/service"
	"management-backend/internal/middleware"
	"management-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// PermissionInfoHandler 动态权限信息处理器
type PermissionInfoHandler struct {
	userSvc  service.UserService
	roleSvc  service.RoleService
	menuSvc  service.MenuService
	roleRepo repository.RoleRepo
}

// NewPermissionInfoHandler 创建权限信息处理器（返回 gin.HandlerFunc）
func NewPermissionInfoHandler(
	userSvc service.UserService,
	roleSvc service.RoleService,
	menuSvc service.MenuService,
	roleRepo repository.RoleRepo,
) gin.HandlerFunc {
	h := &PermissionInfoHandler{
		userSvc:  userSvc,
		roleSvc:  roleSvc,
		menuSvc:  menuSvc,
		roleRepo: roleRepo,
	}
	return h.GetPermissionInfo
}

// GetPermissionInfo 获取当前用户权限信息
func (h *PermissionInfoHandler) GetPermissionInfo(c *gin.Context) {
	ctx := c.Request.Context()
	userID := middleware.GetUserID(c)

	userResp, err := h.userSvc.GetByID(ctx, userID)
	if err != nil {
		response.Fail(c, 500, "获取用户信息失败")
		return
	}

	roles, err := h.roleSvc.GetRolesByUserID(ctx, userID)
	if err != nil {
		response.Fail(c, 500, "获取角色信息失败")
		return
	}

	roleCodes := make([]string, 0, len(roles))
	roleIDs := make([]int64, 0, len(roles))
	for _, r := range roles {
		roleCodes = append(roleCodes, r.RoleCode)
		roleIDs = append(roleIDs, r.ID)
	}

	menuIDs, err := h.roleRepo.GetMenuIDsByRoleIDs(ctx, roleIDs)
	if err != nil {
		response.Fail(c, 500, "获取菜单权限失败")
		return
	}

	menus, err := h.menuSvc.GetByIDs(ctx, menuIDs)
	if err != nil {
		response.Fail(c, 500, "获取菜单信息失败")
		return
	}

	var permissions []string
	for _, m := range menus {
		if m.Permission != "" {
			permissions = append(permissions, m.Permission)
		}
	}

	menuTree := buildFilteredMenuTree(menus, 0)

	resp := smodel.PermissionInfoResp{
		User:        *userResp,
		Roles:       roleCodes,
		Permissions: permissions,
		Menus:       menuTree,
	}
	response.Ok(c, resp)
}

func buildFilteredMenuTree(menus []smodel.Menu, parentID int64) []smodel.MenuTreeResp {
	var tree []smodel.MenuTreeResp
	for _, m := range menus {
		if m.ParentID == parentID {
			node := smodel.MenuTreeResp{
				ID: m.ID, ParentID: m.ParentID, MenuName: m.MenuName,
				MenuType: m.MenuType, Path: m.Path, Component: m.Component,
				Permission: m.Permission, Icon: m.Icon, Sort: m.Sort,
				Visible: m.Visible, Status: m.Status,
				Children: buildFilteredMenuTree(menus, m.ID),
			}
			tree = append(tree, node)
		}
	}
	return tree
}
