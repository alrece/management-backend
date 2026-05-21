## 1. pkg/redisx/ 分布式工具库

- [x] 1.1 实现 pkg/redisx/lock.go：DistributedLock struct，SetNX + UUID + 后台续期 + Lua 解锁
- [x] 1.2 实现 pkg/redisx/cache.go：CacheWrapper struct，泛型 Get/Set/Delete + TTL + ErrCacheMiss
- [x] 1.3 实现 pkg/redisx/idempotent.go：IdempotentExecutor，SetNX 24h TTL + 缓存结果
- [x] 1.4 编写 pkg/redisx/ 单元测试（lock 互斥/过期/非持有者解锁，cache hit/miss/set/delete，idempotent 首次/重复）

## 2. 共享库提取（pkg/authx + authzx + tenantx + discoveryx）

- [x] 2.1 创建 pkg/authx/，从 internal/auth 移植 JWT 生成/解析/黑名单，函数签名接受 JWTConfig 参数
- [x] 2.2 创建 pkg/authzx/，从 internal/authz 移植 Casbin EnforcerManager，接受 redis.Client 参数
- [x] 2.3 创建 pkg/tenantx/，从 internal/multitenant 移植 TenantResolver，接受 MySQLConfig + TenantConfig 参数
- [x] 2.4 创建 pkg/discoveryx/，从 internal/discovery 移植 Consul 注册，接受 ConsulConfig 参数
- [x] 2.5 将 internal/auth, internal/authz, internal/multitenant, internal/discovery 改为薄包装层，调用 pkg/ 包
- [x] 2.6 实现 pkg/authx 的 gRPC UnaryServerInterceptor（从 metadata 提取 token → 解析验证 → 注入 context）
- [x] 2.7 go build ./... 验证无回归 + go test ./... 验证 43 个测试通过

## 3. goose 迁移框架集成

- [x] 3.1 安装 github.com/pressly/goose/v3 依赖
- [x] 3.2 重写 internal/multitenant/migrator.go，使用 goose.Up 替代手动 SQL 执行
- [x] 3.3 将 scripts/sql/migrations/ 下的 SQL 文件重命名为 goose timestamp 格式
- [x] 3.4 创建 cmd/migrate/main.go 独立迁移工具（--target / --all / --allow-missing 标志）
- [x] 3.5 go build ./... + go test ./... 验证

## 4. Redis Sentinel 配置

- [x] 4.1 在 internal/config/config.go 添加 RedisSentinelConfig struct（Enabled, MasterName, Addrs）
- [x] 4.2 修改 cmd/system/main.go，根据 sentinel.enabled 选择 NewClient vs NewFailoverClient
- [x] 4.3 在 docker-compose.yml 添加 3 个 Sentinel 容器 + sentinel.conf（需实际 Docker 环境）
- [x] 4.4 更新 configs/config.yaml 添加 redis.sentinel 配置段（enabled: false 默认）
- [x] 4.5 go build ./... 验证

## 5. job-service 独立微服务

- [x] 5.1 创建 services/job-service/ 目录结构（cmd, internal/{config,model,repository,service,handler,router,scheduler}, configs, Dockerfile, Makefile）
- [x] 5.2 创建 services/job-service/go.mod，引用 management-backend/pkg
- [x] 5.3 更新 go.work 添加 services/job-service 条目
- [x] 5.4 实现 model 层：JobTask, JobExecutionLog entity + DTO
- [x] 5.5 实现 repository 层：JobTaskRepo, JobExecutionLogRepo（CRUD + 分页 + 租户过滤）
- [x] 5.6 实现 service 层：JobService（创建/更新/删除/分页/手动触发）
- [x] 5.7 实现 handler 层 + router：REST API（POST/GET/PUT/DELETE /job/tasks, POST /job/tasks/:id/trigger, GET /job/execution-logs）
- [x] 5.8 实现 scheduler/leader.go：Redis Leader Election（SETNX + 续期 + failover）
- [x] 5.9 实现 scheduler/scheduler.go：robfig/cron 集成 + 执行日志记录 + 幂等保护
- [x] 5.10 实现 cmd/main.go：配置加载 → DB → Redis → pkg/authx JWT 中间件 → pkg/discoveryx Consul 注册 → 优雅关闭
- [x] 5.11 创建 api/proto/job/job.proto 定义 JobService gRPC 接口
- [x] 5.12 编写 job-service 单元测试 + 集成测试
- [x] 5.13 在 docker-compose.yml 添加 mb-job-service 容器

## 6. workflow-service 独立微服务

- [x] 6.1 创建 services/workflow-service/ 目录结构（cmd, internal/{config,model,repository,service,handler,router,client}, configs, Dockerfile, Makefile）
- [x] 6.2 创建 services/workflow-service/go.mod，引用 management-backend/pkg + sony/gobreaker
- [x] 6.3 更新 go.work 添加 services/workflow-service 条目
- [x] 6.4 实现 model 层：WfCategory, WfWorkflow, WfInstance, WfVariable entity + DTO
- [x] 6.5 实现 repository 层：各实体 Repo（CRUD + 分页 + 租户过滤）
- [x] 6.6 实现 internal/client/n8n_client.go：n8n REST API 封装（CRUD + Execute + Webhook）+ Circuit Breaker
- [x] 6.7 实现 service 层：WorkflowService（创建/激活/停用/删除/执行/同步）+ CategoryService
- [x] 6.8 实现 handler/webhook.go：接收 n8n 执行回调，更新 wf_instance 状态
- [x] 6.9 实现 handler 层 + router：REST API（分类树/工作流CRUD/执行/实例查询/回调）
- [x] 6.10 实现 cmd/main.go：配置加载 → DB → Redis → pkg/authx → pkg/discoveryx → n8n 连接检查 → 优雅关闭
- [x] 6.11 实现 n8n 定期同步（5min 间隔，检测 DESYNCED 状态）+ 启动时对账
- [x] 6.12 创建 api/proto/workflow/workflow.proto 定义 WorkflowService gRPC 接口
- [x] 6.13 编写 workflow-service 单元测试 + 集成测试
- [x] 6.14 在 docker-compose.yml 添加 mb-workflow-service + mb-n8n 容器

## 7. Casbin 缓存优化

- [x] 7.1 为 internal/authz EnforcerManager 添加周期性淘汰（每小时清理 30min 未使用的 enforcer）
- [x] 7.2 go build ./... + go test ./... 验证

## 8. Traefik 路由更新

- [x] 8.1 更新 configs/traefik/ 配置，添加 job-service 和 workflow-service 路由规则
- [x] 8.2 验证 Traefik 自动服务发现（Consul provider）对新服务生效

## 9. 前端原型

- [x] 9.1 初始化 Vue3 + Vben5 + TypeScript + Tailwind CSS 项目
- [x] 9.2 实现登录页（JWT 双 Token 流程 + 自动刷新）
- [x] 9.3 实现 Dashboard 布局（动态菜单 + 权限路由）
- [x] 9.4 实现用户管理页面（CRUD + 分页 + 筛选）
- [x] 9.5 实现角色管理页面（权限树分配）
- [x] 9.6 实现菜单管理页面（树形 CRUD）
- [x] 9.7 实现部门管理页面（树形展示）
- [x] 9.8 实现 Mock API 层（VITE_USE_MOCK 开关）
- [x] 9.9 暗色/亮色主题切换
