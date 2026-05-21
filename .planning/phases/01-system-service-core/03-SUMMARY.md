---
phase: 01-system-service-core
plan: 03
subsystem: database
tags: [multitenant, gorm, mysql, redis, lru, migration, isolation]

requires:
  - plan: 01
    provides: pkg/ 共享库、配置中心、中间件模式
  - plan: 02
    provides: JWT 双 Token（Claims.TenantID）、Auth 中间件
provides:
  - 独立数据库多租户隔离（每租户独立 MySQL 数据库）
  - TenantResolver 动态创建/缓存租户 *gorm.DB
  - TenantPool LRU 连接池（50 连接上限，30min 空闲驱逐）
  - Tenant 中间件（JWT → tenantID → 租户 DB → context 注入）
  - Tenant CRUD API
  - 租户初始化幂等状态机（0→1→2，Redis 分布式锁）
  - Schema 迁移工具（MigrateDefault/MigrateTenant/MigrateAll）
  - 租户健康检查
affects: [04-PLAN, 05-PLAN, 06-PLAN, 07-PLAN, 08-PLAN]

tech-stack:
  added: []
  patterns: [database-per-tenant, LRU pool, idempotent state machine, context-based DB routing]

key-files:
  created:
    - internal/multitenant/resolver.go
    - internal/multitenant/pool.go
    - internal/multitenant/initializer.go
    - internal/multitenant/migrator.go
    - internal/multitenant/health.go
    - internal/module/system/model/tenant_entity.go
    - internal/module/system/model/tenant_dto.go
    - internal/module/system/repository/tenant.go
    - internal/module/system/service/tenant.go
    - internal/module/system/handler/tenant.go
    - scripts/sql/tenant_template.sql
    - scripts/sql/migrations/001_init_schema.sql
  modified:
    - internal/middleware/tenant.go
    - internal/config/config.go
    - internal/module/system/repository/user.go
    - internal/module/system/service/user.go
    - internal/module/system/handler/user.go
    - internal/module/system/router.go
    - internal/router/router.go
    - pkg/middleware/context.go
    - configs/config.yaml
    - scripts/sql/init.sql

key-decisions:
  - "独立数据库隔离而非行级隔离（用户铁律决策）"
  - "Repository 通过 context.Context 获取租户 DB（getDB(ctx) 模式）"
  - "LRU Pool 用 container/list 零外部依赖实现"
  - "初始化状态机用 CAS + Redis SetNX 分布式锁"
  - "tenantScope 参数从 Repository 接口移除"

patterns-established:
  - "getDB(ctx) 模式: Repository 优先使用 context 中的租户 DB"
  - "独立 DB 隔离: 每租户独立 MySQL 数据库（mb_tenant_{id}）"
  - "LRU 连接池: 50 上限 + 30min 空闲驱逐"
  - "幂等状态机: CAS 状态转换 + Redis 分布式锁"

requirements-completed: [FR-05-01, FR-05-02, FR-05-03, FR-05-04, FR-05-05, FR-05-08, NFR-01-06, NFR-03-03]

duration: 25min
completed: 2026-05-18
---

# Plan 03: 多租户核心 — 独立数据库隔离 Summary

**独立数据库多租户隔离系统：TenantResolver + LRU 连接池 + 幂等状态机初始化 + 文件迁移工具**

## Performance

- **Duration:** ~25 min
- **Completed:** 2026-05-18
- **Tasks:** 5
- **Files modified:** 24 (+1528/-40 lines)

## Accomplishments
- 完整租户 CRUD API（entity → DTO → repo → service → handler → router）
- TenantResolver 动态创建租户 *gorm.DB，DSN 格式 mb_tenant_{id}
- TenantPool LRU 连接池：50 上限、30min 空闲驱逐、atomic 计数器
- Tenant 中间件重写：JWT tenantID → GetDB → context 注入
- Repository `getDB(ctx)` 模式：透明切换租户/默认 DB
- 初始化器：4 状态幂等状态机 + Redis 分布式锁
- Schema 迁移器：MigrateDefault/MigrateTenant/MigrateAll
- 健康检查：DB 连通性 + 表数量

## Task Commits

1. **All tasks** - `0f4fee2` (feat)

## Decisions Made
- Repository `getDB(ctx)` 模式：优先从 context 获取租户 DB
- 移除 `tenantScope` 参数（独立 DB 不需行级过滤）
- Prometheus 指标用 `atomic.Int64` 预留

## Deviations from Plan
- 未引入 gorm.io/plugin/dbresolver（自定义 TenantResolver 更灵活）
- 未引入 goose 库（轻量文件排序执行，避免额外依赖）

## Next Phase Readiness
- Plan 04（RBAC 权限体系）可直接使用租户 DB

---
*Phase: 01-system-service-core*
*Completed: 2026-05-18*
