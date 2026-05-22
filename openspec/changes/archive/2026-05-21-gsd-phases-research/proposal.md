## Why

Phase 1（System Service Core）已完成全部 9 个 GSD Plan（132 个 Go 文件，43 个测试通过），涵盖认证、RBAC、多租户、系统 CRUD、微服务基础设施和可观测性。现在需要将 Phase 1 成果转化为可追溯的规范约束集，并为 Phase 2（扩展模块）建立清晰的前置条件和实施路线。同时，在 Phase 2 编码前需先生成前端原型进行展示和验证。

## What Changes

### Phase 1 约束归档
- 将 Phase 1 实施中确立的 9 条硬约束和 7 条软约束文档化
- 标记 Phase 1 已完成需求（FR-01 ~ FR-08 中的 P0 项）为 validated

### Phase 2 前置条件准备
- **pkg/redisx/** 填充分布式锁 + 缓存封装工具库
- **共享库提取**：将 auth/authz/multitenant/config 提取为独立 Go module，支持多微服务引用
- **goose 迁移框架**集成，替代手动 SQL 文件管理
- **Redis Sentinel** 配置引入（生产高可用）

### Phase 2 新模块（按默认优先级顺序）
- **job-service**：定时任务管理、执行日志、分布式协调（独立微服务部署）
- **workflow-service**：n8n 集成、流程分类、流程实例管理 + Circuit Breaker（独立微服务部署）
- **代码生成器**：微服务脚手架生成（cmd/Dockerfile/Makefile/Protobuf/三层代码/SQL/前端 API）
- **租户扩展**：套餐管理、数据导出/清除
- **扩展业务模块**：CRM / ERP / 商城 / IoT / MES / AI（按需）

### 前端原型
- 使用 Vue3 + Vben5 生成完整前端原型，在编码前进行展示和验证
- 基于后端 Swagger/Proto API 契约定义

## Capabilities

### New Capabilities
- `shared-redis-utils`: Redis 分布式锁（基于 SetNX + 续期）、缓存封装（Get/Set/Delete with TTL）、分布式幂等工具
- `shared-lib-extraction`: 将 internal/auth、internal/authz、internal/multitenant、internal/config 提取为 pkg/ 下独立 Go module，供多微服务共享引用
- `migration-framework`: goose 版本化迁移框架集成，支持多租户批量执行
- `redis-sentinel`: Redis Sentinel 高可用配置
- `job-service`: 独立微服务，定时任务 CRUD、执行日志、分布式协调（Redis + Consul）
- `workflow-service`: 独立微服务，n8n API 集成、流程分类/实例管理、Circuit Breaker
- `frontend-prototype`: Vue3 + Vben5 前端原型生成，用于展示验证后端 API 契约
- `code-generator`: 微服务脚手架代码生成器

### Modified Capabilities
（Phase 1 无已有 OPSX spec，无需修改）

## Impact

### 代码影响
- **pkg/redisx/**：新建目录，实现分布式锁和缓存工具
- **pkg/** 下新增共享 Go module：auth、authz、multitenant 配置包
- **go.work**：新增 job-service、workflow-service 模块
- **docker-compose.yml**：新增 job-service、workflow-service 容器、Redis Sentinel 节点
- **internal/multitenant/migrator.go**：从手动 SQL 切换为 goose
- **scripts/migrations/**：SQL 文件迁移到 goose 版本格式

### API 影响
- 新增 job-service REST API（任务管理、执行日志查询）
- 新增 workflow-service REST API（流程管理、实例操作）
- 新增 gRPC 内部服务（job-service ↔ system-service 通信）

### 依赖影响
- 新增 goose 迁移工具依赖
- job-service 新增 robfig/cron/v3 依赖
- workflow-service 新增 n8n HTTP client 依赖
- 前端原型新增 Vue3 + Vben5 + TypeScript 依赖

### 系统影响
- Docker Compose 服务数量增加（job-service + workflow-service + Redis Sentinel）
- Consul 注册服务数量增加
- Traefik 路由规则扩展
