# Research Summary for OPSX

## Discovered Constraints

### Hard Constraints
- HC-1: RuoYi 兼容响应格式 `R{code:0, data, msg}` 不可更改
- HC-2: 雪花 ID 主键，禁止 AUTO_INCREMENT
- HC-3: 四层架构 model→repo→service→handler 不可跳层
- HC-4: context.Context 传播，禁止 *gin.Context 穿透到 Service
- HC-5: 独立数据库多租户隔离，所有查询必须注入 tenantScope
- HC-6: JWT Secret 启动校验，MySQL 密码禁止弱密码
- HC-7: gRPC 连接当前 insecure，生产环境必须 TLS
- HC-8: 手动 DI 无框架，新模块需在各自 main.go 中组装
- HC-9: 新模块必须独立微服务部署（独立 main.go、独立进程、独立 Docker 容器）

### Soft Constraints
- SC-1: 时间格式硬编码 "2006-01-02 15:04:05" 应统一为常量
- SC-2: pkg/redisx/ 需实现分布式锁 + 缓存封装
- SC-3: Storage 动态切换非线程安全（globalProvider 无锁），需加互斥锁
- SC-4: 无版本化迁移框架（goose 未集成），Phase 2 需引入
- SC-5: toXxxResp() 手动映射散布所有 Service，可考虑代码生成辅助
- SC-6: 分页 Count + Find 分两次调用，高并发下可能不一致
- SC-7: 前端原型应先于完整编码，用于展示和验证

## Dependencies

### Cross-Module Dependencies
- pkg/ 共享库 (response/errcode/page/snowflake/crypto/gormx) — 所有微服务共享
- internal/auth (JWT + 黑名单) — 认证核心，需提取为独立包供多服务使用
- internal/authz (Casbin RBAC) — 权限校验，需提取为共享库
- internal/multitenant (TenantResolver) — 多租户路由，需提取为共享库
- internal/config (Viper 配置) — 需支持多服务配置合并/拆分
- docker-compose.yml — 基础设施依赖

### Implementation Order Dependencies
1. pkg/redisx/ 填充 → 所有需要 Redis 的服务受益
2. 共享库提取（auth/authz/multitenant/config）→ 微服务拆分前提
3. job-service → workflow-service → 代码生成器 → 租户扩展 → 业务模块

## Risks & Mitigations

| Risk | Severity | Mitigation |
|------|----------|------------|
| Redis 单点，无哨兵/集群 | 高 | Phase 2 引入 Redis Sentinel 配置 |
| Saga 无超时回滚机制 | 高 | 补充超时补偿 + 后台重试 worker |
| 无自动化备份 | 高 | 集成 XtraBackup + cron 定时任务 |
| gRPC AuthInterceptor 仅透传 | 中 | 生产前实现 gRPC JWT 校验 |
| 共享库提取可能破坏现有编译 | 中 | go.work + replace 逐步迁移 |
| 前端原型与后端 API 契约不一致 | 中 | 先定义 API 契约再并行开发 |
| Casbin sync.Map 缓存无过期 | 中 | 增加 TTL 或 LRU 淘汰 |

## Success Criteria

### Phase 1 完成验证 (已达成)
- 9/9 GSD Plans 已完成，43 个测试全部通过
- 132 个 Go 文件，约 10,711 行代码
- 认证链路、多租户隔离、RBAC 权限、系统 CRUD、微服务基础设施、可观测性全部就绪

### Phase 2 启动前置条件
- [ ] pkg/redisx/ 实现分布式锁 + 缓存封装
- [ ] 共享库提取为独立 Go module（供多服务引用）
- [ ] 前端原型生成（Vue3 + Vben5，用于展示验证）
- [ ] API 契约定义（Swagger 或 Proto，前后端并行基础）
- [ ] goose 迁移框架集成
- [ ] Redis Sentinel 配置

## User Confirmations

1. Phase 2 按默认顺序：job-service → workflow-service → 代码生成器 → 租户扩展 → 业务模块
2. 新模块独立微服务部署（独立 main.go、独立进程、独立 Docker 容器）
3. pkg/redisx/ 需实现分布式锁 + 缓存封装
4. 前端原型先行：在完整编码前用前端技术生成整套原型进行展示和验证

## Codebase Metrics

- Go files: 132, Lines: ~10,711
- Entities: 15, Proto files: 3, Test files: 12 (43 tests)
- Docker services: MySQL 8.0, Redis 7, MinIO, Consul 1.15, Jaeger 1.50
