---
wave: 3
depends_on: ["03-PLAN.md"]
files_modified:
  - internal/authz/casbin.go
  - internal/authz/adapter.go
  - internal/authz/sync.go
  - internal/middleware/permission.go
  - internal/module/system/model/role_entity.go
  - internal/module/system/model/menu_entity.go
  - internal/module/system/model/dept_entity.go
  - internal/module/system/model/post_entity.go
  - internal/module/system/model/permission_dto.go
  - internal/module/system/repository/role.go
  - internal/module/system/repository/menu.go
  - internal/module/system/repository/dept.go
  - internal/module/system/repository/post.go
  - internal/module/system/service/role.go
  - internal/module/system/service/menu.go
  - internal/module/system/service/dept.go
  - internal/module/system/service/post.go
  - internal/module/system/handler/role.go
  - internal/module/system/handler/menu.go
  - internal/module/system/handler/dept.go
  - internal/module/system/handler/post.go
  - pkg/gormx/scope.go
  - scripts/sql/init.sql
  - scripts/sql/tenant_template.sql
autonomous: true
requirements_addressed:
  - FR-01-03
  - FR-01-04
  - FR-01-05
  - FR-03-01
  - FR-03-02
  - FR-03-03
  - FR-03-04
  - FR-03-05
  - FR-01-12
---

# Plan 04: RBAC + 数据权限

## Objective

实现 Casbin RBAC + Redis Adapter + 多实例策略同步、角色管理、菜单管理三级树、声明式权限中间件、部门树管理、数据权限 GORM Scope（5 级）、岗位管理、前端动态菜单 API。

## Tasks

### Task 4.1: Casbin 集成 + Redis Adapter + 策略同步

<read_first>
- go.mod (casbin 版本)
- internal/auth/jwt.go
</read_first>

<action>
创建 internal/authz/casbin.go：EnforcerManager 管理 sync.Map 缓存的多租户 Enforcer。Redis Adapter（key: casbin:policy:{tenantID}）替代 gorm-adapter。Redis Pub/Sub 策略同步（authz/sync.go）。Casbin model: p=sub,dom,obj,act; g=_,_,_。
</action>

<acceptance_criteria>
- GetEnforcer(ctx, tenantID) 返回租户专属 Enforcer
- Redis adapter 替代 gorm-adapter
- Pub/Sub 变更通知机制
</acceptance_criteria>

---

### Task 4.2: 角色管理 + 菜单权限分配

<read_first>
- scripts/sql/init.sql
</read_first>

<action>
Role struct 含 DataScope(1-5), DataScopeDeptIDs。Menu struct 含 Type(1=目录/2=菜单/3=按钮), Path, Component, Permission, Visible。RoleMenu 关联表。角色分配菜单时更新 Casbin 策略 + Pub/Sub 通知。菜单存储在默认库。
</action>

<acceptance_criteria>
- Role 含 DataScope, DataScopeDeptIDs
- Menu 含 Type, Path, Component, Permission, Visible
- 分配菜单 API 更新 Casbin 策略
</acceptance_criteria>

---

### Task 4.3: 菜单管理三级树

<read_first>
- internal/module/system/model/menu_entity.go
</read_first>

<action>
GET /api/system/menu/tree 返回嵌套树。删除前检查子菜单和角色引用。完整 CRUD。
</action>

<acceptance_criteria>
- 树形返回目录→菜单→按钮三级结构
- 删除有子菜单/角色引用时拒绝
</acceptance_criteria>

---

### Task 4.4: 声明式权限中间件

<read_first>
- internal/authz/casbin.go
</read_first>

<action>
RequirePermission(permission string) gin.HandlerFunc：userID + tenantID + permission → Casbin 检查。路由注册时声明 system:{module}:{action} 格式权限标识。
</action>

<acceptance_criteria>
- RequirePermission 中间件存在
- 所有受保护路由声明权限标识
</acceptance_criteria>

---

### Task 4.5: 部门树 + 数据权限 GORM Scope

<read_first>
- pkg/gormx/scope.go
</read_first>

<action>
Dept struct 含 Ancestors("0,1,2")。DataPermissionScope 实现 5 级：1=全部, 2=自定义, 3=本部门, 4=本部门及以下, 5=仅本人。角色数据范围配置 API。
</action>

<acceptance_criteria>
- Ancestors 格式 "0,1,2"
- 5 级 GORM Scope 在 Repository 层应用
- 数据范围配置 API 存在
</acceptance_criteria>

---

### Task 4.6: 岗位管理 CRUD

<read_first>
- scripts/sql/tenant_template.sql
</read_first>

<action>
Post struct（Code, Name, Sort, Status）。CRUD + /api/system/post/ 路由。
</action>

<acceptance_criteria>
- Post CRUD API 路由存在
</acceptance_criteria>

---

### Task 4.7: 前端动态菜单 API

<read_first>
- internal/module/system/service/menu.go
- internal/authz/casbin.go
</read_first>

<action>
GET /api/system/auth/get-permission-info 返回 user, roles, permissions, menus（嵌套树）。根据角色过滤菜单，permissions 包含按钮级标识。
</action>

<acceptance_criteria>
- get-permission-info 返回 menus 树 + permissions 数组
- 菜单按角色权限过滤
</acceptance_criteria>

---

## Verification

```bash
go build ./cmd/system/
go test ./internal/authz/... -v
go vet ./...
```

## Must Haves

- [ ] Casbin RBAC + Redis Adapter + Pub/Sub 同步
- [ ] 角色管理 + 菜单权限分配
- [ ] 菜单三级树（目录/菜单/按钮）
- [ ] RequirePermission 声明式中间件
- [ ] 部门树 + Ancestors 维护
- [ ] 数据权限 5 级 GORM Scope
- [ ] 岗位 CRUD
- [ ] get-permission-info 动态菜单 API
