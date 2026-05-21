---
plan: 04
phase: 01-system-service-core
status: complete
started: 2026-05-18
completed: 2026-05-18
self_check: passed
---

# Plan 04: RBAC + 数据权限

## Summary

实现了完整的 Casbin RBAC 权限体系：Redis Adapter 多租户策略存储、Pub/Sub 多实例同步、角色管理含数据权限配置、菜单三级树（目录/菜单/按钮）、声明式 RequirePermission 中间件、部门树 Ancestors 维护、5 级数据权限 GORM Scope、岗位 CRUD、前端动态菜单 get-permission-info API。

## Tasks Completed

| Task | Description | Status |
|------|-------------|--------|
| 4.1 | Casbin 集成 + Redis Adapter + 策略同步 | ✓ |
| 4.2 | 角色管理 + 菜单权限分配 | ✓ |
| 4.3 | 菜单管理三级树 | ✓ |
| 4.4 | 声明式权限中间件 | ✓ |
| 4.5 | 部门树 + 数据权限 GORM Scope | ✓ |
| 4.6 | 岗位管理 CRUD | ✓ |
| 4.7 | 前端动态菜单 API | ✓ |

## Key Files

### Created
- internal/authz/casbin.go — EnforcerManager + CheckPermission + SyncRolePermissions + SyncUserRoles
- internal/authz/adapter.go — Redis 存储适配器（替代 gorm-adapter）
- internal/authz/sync.go — Redis Pub/Sub 策略同步器
- internal/middleware/permission.go — RequirePermission 声明式中间件
- internal/module/system/model/role_entity.go — Role 含 DataScope/DataScopeDeptIDs
- internal/module/system/model/menu_entity.go — Menu 含 Type/Path/Component/Permission/Visible
- internal/module/system/model/dept_entity.go — Dept 含 Ancestors
- internal/module/system/model/post_entity.go — Post 实体
- internal/module/system/model/permission_dto.go — 权限管理 DTO
- internal/module/system/repository/role.go — RoleRepo 含 AssignMenus/GetByUserID
- internal/module/system/repository/menu.go — MenuRepo 含 Tree/HasChildren/GetByIDs
- internal/module/system/repository/dept.go — DeptRepo 含 Ancestors 维护
- internal/module/system/repository/post.go — PostRepo CRUD
- internal/module/system/service/role.go — RoleService 含 Casbin 策略同步
- internal/module/system/service/menu.go — MenuService 含 Tree 递归构建
- internal/module/system/service/dept.go — DeptService 含 Ancestors 维护
- internal/module/system/service/post.go — PostService CRUD
- internal/module/system/handler/role.go — RoleHandler
- internal/module/system/handler/menu.go — MenuHandler
- internal/module/system/handler/dept.go — DeptHandler
- internal/module/system/handler/post.go — PostHandler
- internal/module/system/handler/permission.go — GetPermissionInfo 动态菜单 API
- pkg/gormx/scope.go — DataPermissionScope 5 级 GORM Scope

### Modified
- internal/module/system/router.go — RBAC 路由注册
- internal/router/router.go — DI 注入 + 权限组件初始化
- pkg/errcode/code.go — RBAC 错误码
- scripts/sql/init.sql — RBAC 表结构和菜单数据
- scripts/sql/tenant_template.sql — 租户模板含 role_menu 关联表

## Deviations

无重大偏差。

## Verification

```bash
go build ./...  # 编译通过
```

## Self-Check

- [x] Casbin RBAC + Redis Adapter + Pub/Sub 同步
- [x] 角色管理 + 菜单权限分配
- [x] 菜单三级树（目录/菜单/按钮）
- [x] RequirePermission 声明式中间件
- [x] 部门树 + Ancestors 维护
- [x] 数据权限 5 级 GORM Scope
- [x] 岗位 CRUD
- [x] get-permission-info 动态菜单 API
