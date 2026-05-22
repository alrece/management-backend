## Context

Phase 1（System Service Core）已完成，132 个 Go 文件 + 43 个测试。当前为单体 Gin 应用，微服务基础设施（gRPC/Consul/Saga/Telemetry）框架已就绪但未独立部署。Phase 2 需要将 system-service 的共享能力提取为可复用包，并新增 job-service 和 workflow-service 两个独立微服务。

**当前状态**：
- internal/auth, internal/authz, internal/multitenant 紧耦合 system-service（读取 config.C 全局变量）
- pkg/redisx/ 空目录，无 Redis 工具封装
- migrator.go 手动执行 SQL 文件，无版本追踪
- gRPC AuthInterceptor 仅透传，未实际校验 JWT

## Goals / Non-Goals

**Goals:**
- 填充 pkg/redisx/（分布式锁 + 缓存 + 幂等），所有新服务可复用
- 将 auth/authz/multitenant/discovery 提取为 pkg/ 下可复用包，新服务直接引用
- 集成 goose 版本化迁移，支持多租户批量执行
- 新增 job-service 独立微服务（定时任务 CRUD + 分布式调度）
- 新增 workflow-service 独立微服务（n8n 集成 + Circuit Breaker）
- Redis Sentinel 高可用配置
- 生成 Vue3 + Vben5 前端原型用于展示验证

**Non-Goals:**
- 不修改 Phase 1 已有 API 契约（HC-1）
- 不引入 Go 微服务框架（Kratos/go-zero）
- 不实现 Redlock（单实例 SetNX 足够）
- 不在 Phase 2 引入 Consul KV 配置中心（服务数 < 5）
- 不实现 Temporal/Cadence 替代 n8n
- 不拆分 pkg/ 为独立 Go module（单 module + 多 package）

## Decisions

### D1: 共享库提取策略 — pkg/ 下多 package（非独立 module）

**选择**：将 internal/ 下的 auth/authz/multitenant/discovery 提取为 pkg/authx/、pkg/authzx/、pkg/tenantx/、pkg/discoveryx/

**理由**：内部 monorepo，单 module + 多 package 管理最简单。独立 module 版本管理在 3 个服务场景下增加复杂度但收益低。

**迁移路径**：
1. 创建 pkg/authx/，函数签名接受 JWTConfig 参数（不读 config.C）
2. internal/auth/ 改为薄包装层，调用 pkg/authx 并传入 config.C 值
3. 新微服务直接 import pkg/authx/

### D2: 分布式锁 — 单实例 SetNX + Lua 解锁

**选择**：Redis SetNX + UUID value + 后台续期 + Lua 原子解锁

**参数**：默认 TTL=10s，续期间隔=TTL/3，key 前缀 `lock:`

### D3: job-service 调度 — robfig/cron + Redis Leader Election

**选择**：robfig/cron/v3 作为 cron 引擎，Redis SETNX 实现 Leader Election

**参数**：Leader TTL=30s，续期间隔=8s，执行锁 TTL=300s

### D4: workflow-service Circuit Breaker — sony/gobreaker

**选择**：github.com/sony/gobreaker

**参数**：连续失败 5 次触发 Open，30s 后半开，半开允许 3 个请求

### D5: n8n 集成 — REST API + Webhook 回调

**选择**：workflow-service 通过 n8n REST API 管理工作流，通过 Webhook 接收执行状态回调

**参数**：API Key 认证（MB_WORKFLOW_N8N_API_KEY），HTTP 超时 CRUD=30s/执行=120s，重试 3 次（仅 5xx）

### D6: 数据库迁移 — goose/v3

**选择**：github.com/pressly/goose/v3 替代手动 SQL 执行

**理由**：版本追踪（goose_db_version 表），幂等执行，支持 go:embed 嵌入

### D7: 雪花 ID Node 分配

| 服务 | Node ID |
|------|---------|
| system-service | 1 |
| job-service | 2 |
| workflow-service | 3 |

### D8: Redis Sentinel

**选择**：3 Sentinel 节点 + 1 Master + 配置开关 redis.sentinel.enabled

## Risks / Trade-offs

| Risk | Severity | Mitigation |
|------|----------|------------|
| 共享库提取破坏现有编译 | HIGH | 分两阶段迁移，每步跑 go build + 43 测试 |
| gRPC JWT 仅透传未校验 | HIGH | pkg/authx 实现 gRPC UnaryServerInterceptor |
| 跨服务事务边界 | HIGH | 现有 Saga orchestrator + 幂等补偿 |
| goose 与 init.sql 冲突 | MEDIUM | goose --allow-missing 跳过已存在表 |
| n8n 启动延迟阻断 workflow-service | MEDIUM | CB 模式，初始连接失败不阻断启动 |
| Casbin sync.Map 无过期 | MEDIUM | 周期性淘汰（每小时清理 30min 未使用 enforcer） |
| go.work 多模块依赖同步 | MEDIUM | go work sync + CI 全量编译 |

## Migration Plan

1. pkg/redisx/ — 纯新增，零影响
2. 共享库提取 — 创建 pkg/ 新包 → internal/ 改为包装层 → go build 验证
3. goose 集成 — 重写 migrator.go → 迁移现有 SQL → go build 验证
4. Redis Sentinel 配置 — config.go 新增字段 → docker-compose.yml 新增容器
5. job-service — 创建 services/job-service/ → 四层架构 → Consul 注册
6. workflow-service — 创建 services/workflow-service/ → n8n client + CB
7. 前端原型 — 独立项目，基于 Swagger API 契约

每步完成后 `go build ./...` + `go test ./...` 验证无回归。
