---
wave: 2
depends_on: ["01-PLAN.md", "02-PLAN.md"]
files_modified:
  - internal/multitenant/resolver.go
  - internal/multitenant/pool.go
  - internal/multitenant/initializer.go
  - internal/multitenant/migrator.go
  - internal/multitenant/health.go
  - internal/middleware/tenant.go
  - internal/module/system/model/tenant_entity.go
  - internal/module/system/model/tenant_dto.go
  - internal/module/system/repository/tenant.go
  - internal/module/system/service/tenant.go
  - internal/module/system/handler/tenant.go
  - scripts/sql/tenant_template.sql
  - scripts/sql/migrations/
  - configs/config.yaml
autonomous: true
requirements_addressed:
  - FR-05-01
  - FR-05-02
  - FR-05-03
  - FR-05-04
  - FR-05-05
  - FR-05-08
  - NFR-01-06
  - NFR-03-03
---

# Plan 03: 多租户核心 — 独立数据库隔离

## Objective

实现独立数据库级别的多租户隔离：GORM DBResolver 多库路由、租户连接池 LRU 管理（带 Prometheus 指标）、Tenant 中间件（JWT → DB 切换）、租户初始化流程（幂等状态机）、Schema 迁移工具（goose）、租户健康检查 API。

## Context

**用户铁律决策：独立数据库多租户，不是行级隔离。** 当前 tenant.go 仅做行级 GORM scope 过滤，需要完全重写为独立 DB 路由。

## Tasks

### Task 3.1: 租户表 CRUD + 实体 + 模板 Schema

<read_first>
- internal/module/system/model/entity.go (参考实体模式)
- scripts/sql/init.sql (当前建表)
</read_first>

<action>
1. 创建 internal/module/system/model/tenant_entity.go：Tenant struct 含 ID, Name, ContactName, ContactMobile, Status, DBName, ExpireTime, InitStatus(0=待初始化/1=初始化中/2=已就绪/3=失败), Remark, 审计字段, DeletedAt。TableName 返回 "sys_tenant"。

2. 创建 tenant_dto.go：TenantCreateReq, TenantUpdateReq, TenantPageReq, TenantResp。

3. 更新 init.sql 添加 sys_tenant 建表语句（放在默认库）。

4. 创建 scripts/sql/tenant_template.sql — 租户独立库模板 Schema（sys_user, sys_role, sys_dept, sys_user_role, sys_dict_type, sys_dict_data, sys_post, sys_user_post）。

5. 创建完整的 tenant repository/service/handler CRUD + 路由注册。
</action>

<acceptance_criteria>
- Tenant struct 含 DBName, InitStatus 字段
- Tenant.TableName() == "sys_tenant"
- sys_tenant 在 init.sql 中
- tenant_template.sql 包含 sys_user, sys_role, sys_dept, sys_dict_type, sys_dict_data, sys_post
- CRUD API 路由注册在 /api/system/tenant/
</acceptance_criteria>

---

### Task 3.2: GORM DBResolver 多库路由 + LRU 连接池

<read_first>
- internal/middleware/tenant.go (当前行级隔离)
- internal/config/config.go
</read_first>

<action>
1. 创建 internal/multitenant/resolver.go：TenantResolver 依赖默认 *gorm.DB，GetDB(ctx, tenantID) 通过 DSN 动态创建/缓存租户 *gorm.DB。

2. 创建 internal/multitenant/pool.go：TenantPool LRU 连接池，MaxPoolSize=50, IdleTimeout=30min，Prometheus 指标（tenant_pool_size, tenant_pool_evictions_total, tenant_db_connections_active），后台清理 goroutine。

3. Config 添加 TenantConfig（MaxPoolSize, IdleTimeout, MaxConnsPerTenant, DBNamePrefix）。

4. 重构 internal/middleware/tenant.go：JWT → tenantID → TenantResolver.GetDB → context.WithValue 注入 *gorm.DB。验证租户状态和 InitStatus。
</action>

<acceptance_criteria>
- TenantResolver.GetDB(ctx, tenantID) 返回租户专属 *gorm.DB
- TenantPool 含 LRU 驱逐 + 3 个 Prometheus 指标
- TenantConfig 含 4 个配置字段
- tenant middleware 从 JWT 提取 tenantID → 获取租户 DB → context 注入
</acceptance_criteria>

---

### Task 3.3: 租户初始化流程（幂等状态机）

<read_first>
- internal/module/system/model/tenant_entity.go (Tenant.InitStatus)
</read_first>

<action>
1. 创建 internal/multitenant/initializer.go：状态机 0→1→2 成功，1→3 失败，3→0 可重试。步骤：CAS 更新状态 → CREATE DATABASE → 执行模板 SQL → 创建默认管理员 → CAS 更新就绪。Redis 分布式锁防并发。

2. Tenant Service.Create 触发异步初始化。

3. 添加重试 API：POST /api/system/tenant/:id/retry。
</action>

<acceptance_criteria>
- Initialize(ctx, tenantID) 幂等方法
- 状态机 4 个状态正确流转
- Redis 分布式锁 key: tenant:init:{id}
- 重试 API 仅在 InitStatus=3 时可用
</acceptance_criteria>

---

### Task 3.4: Schema 迁移工具（goose）

<read_first>
- scripts/sql/tenant_template.sql
</read_first>

<action>
1. 引入 github.com/pressly/goose/v3。

2. 创建 scripts/sql/migrations/001_init_schema.sql。

3. 创建 internal/multitenant/migrator.go：MigrateTenant, MigrateAll, MigrateDefault。

4. CLI 支持 migrate-up 命令（--all-tenants, --tenant=123）。
</action>

<acceptance_criteria>
- goose 在 go.mod 中
- 001_init_schema.sql 存在
- migrator.go 导出 3 个方法
- CLI 支持 migrate-up 含参数
</acceptance_criteria>

---

### Task 3.5: 租户健康检查 API

<read_first>
- internal/multitenant/resolver.go
</read_first>

<action>
创建 internal/multitenant/health.go：HealthCheck(ctx, tenantID) 返回 TenantHealth（DBConnected, SchemaVersion, TableCount, LastCheckTime）。添加 GET /api/system/tenant/:id/health 和 GET /api/system/tenant/health-all 端点。
</action>

<acceptance_criteria>
- HealthCheck 返回 TenantHealth struct
- 两个健康检查端点存在
</acceptance_criteria>

---

## Verification

```bash
go build ./cmd/system/
go test ./internal/multitenant/... -v
go vet ./...
```

## Must Haves

- [ ] 独立数据库多租户（每租户独立 DB）
- [ ] GORM DBResolver 多库路由
- [ ] LRU 连接池（30min 空闲驱逐 + Prometheus 指标）
- [ ] Tenant 中间件（JWT → DB 切换 + 状态验证）
- [ ] 租户初始化状态机（幂等、可重试）
- [ ] goose Schema 迁移（多租户批量执行）
- [ ] 租户健康检查 API
