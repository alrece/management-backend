---
wave: 4
depends_on: ["04-PLAN.md"]
files_modified:
  - internal/module/system/model/entity.go
  - internal/module/system/model/dto.go
  - internal/module/system/model/log_entity.go
  - internal/module/system/model/dict_entity.go
  - internal/module/system/model/param_entity.go
  - internal/module/system/model/notice_entity.go
  - internal/module/system/repository/user.go
  - internal/module/system/repository/dict_type.go
  - internal/module/system/repository/dict_data.go
  - internal/module/system/repository/param.go
  - internal/module/system/repository/notice.go
  - internal/module/system/repository/log.go
  - internal/module/system/service/user.go
  - internal/module/system/service/profile.go
  - internal/module/system/service/user_export.go
  - internal/module/system/service/dict.go
  - internal/module/system/service/param.go
  - internal/module/system/service/notice.go
  - internal/module/system/service/log.go
  - internal/module/system/handler/user.go
  - internal/module/system/handler/profile.go
  - internal/module/system/handler/dict.go
  - internal/module/system/handler/param.go
  - internal/module/system/handler/notice.go
  - internal/module/system/handler/log.go
  - internal/middleware/operation_log.go
  - scripts/sql/init.sql
  - scripts/sql/tenant_template.sql
autonomous: true
requirements_addressed:
  - FR-02-01
  - FR-02-02
  - FR-02-03
  - FR-02-04
  - FR-04-01
  - FR-04-02
  - FR-04-03
  - FR-04-04
  - FR-04-05
  - FR-06-01
  - FR-06-02
  - NFR-05-04
---

# Plan 05: System CRUD 模块

## Objective

完成用户管理（完善多租户）、个人中心、用户导入导出、字典管理、参数管理、通知公告、操作日志 + 登录日志。

## Tasks

### Task 5.1: 用户管理 CRUD（多租户 + 数据权限）

<read_first>
- internal/module/system/handler/user.go
- internal/module/system/repository/user.go
- pkg/gormx/scope.go
</read_first>

<action>
重构 user.go：GetByID/Delete/Update 注入 tenantScope（修复跨租户越权）。用户创建分配 dept/role/post。Repository 从 context 获取租户 DB。列表查询应用 DataPermissionScope。UserRole/UserPost 关联表。
</action>

<acceptance_criteria>
- GetByID/Delete/Update 注入 tenantScope
- Repository 从 context 获取租户 DB
- 列表查询应用数据权限 Scope
</acceptance_criteria>

---

### Task 5.2: 个人中心 + 用户导入导出

<read_first>
- internal/middleware/password.go
- internal/auth/blacklist.go
</read_first>

<action>
profile handler/service：GET profile, PUT profile, PUT password（校验+复杂度+吊销Token）。导出用 StreamWriter 限 10 万行，导入用 excelize 批量创建。
</action>

<acceptance_criteria>
- 改密时调用 RevokeByUser
- StreamWriter 导出，10 万行上限
</acceptance_criteria>

---

### Task 5.3: 字典管理

<read_first>
- scripts/sql/tenant_template.sql
</read_first>

<action>
DictType/DictData 实体，CRUD API，GET /api/system/dict-data/type/:type 按 type 查询。Redis 缓存 dict:{type}。
</action>

<acceptance_criteria>
- 按 type 查询字典数据
- Redis 缓存 + 变更清除
</acceptance_criteria>

---

### Task 5.4: 参数管理

<read_first>
- scripts/sql/tenant_template.sql
</read_first>

<action>
SysParam 实体，CRUD API，Redis 缓存 param:{key}，按 key 查询。
</action>

<acceptance_criteria>
- 按 key 查询参数
- Redis 缓存
</acceptance_criteria>

---

### Task 5.5: 通知公告

<read_first>
- scripts/sql/tenant_template.sql
</read_first>

<action>
Notice 实体，CRUD API /api/system/notice/。
</action>

<acceptance_criteria>
- Notice CRUD 端点存在
</acceptance_criteria>

---

### Task 5.6: 操作日志 + 登录日志

<read_first>
- internal/middleware/request_id.go
- internal/auth/jwt.go
</read_first>

<action>
SysOperLog/SysLoginLog 实体。OperationLog 中间件记录写操作。auth handler 记录登录日志。日志查询/删除 API 4 个端点。日志含 tenant_id, request_id。
</action>

<acceptance_criteria>
- OperationLog 中间件存在
- 登录日志在 auth handler 中记录
- 4 个日志 API 端点
- 日志含 tenant_id, request_id 字段
</acceptance_criteria>

---

## Verification

```bash
go build ./cmd/system/
go test ./internal/module/system/... -v
```

## Must Haves

- [ ] 用户 CRUD（跨租户越权修复）
- [ ] 个人中心（改密吊销 Token）
- [ ] 用户导入/导出（StreamWriter, 10万行限制）
- [ ] 字典管理（Redis 缓存）
- [ ] 参数管理（Redis 缓存）
- [ ] 操作日志 + 登录日志
