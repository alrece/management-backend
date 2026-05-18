# Roadmap — Management Backend (Go)

## Phase 1: System Service (Core)

目标：完成 system-service 的全部功能，包括认证、RBAC、多租户（独立数据库）、数据权限、CRUD、文件管理。

---

### Milestone 1.1: Foundation & Security

**Goal**: 微服务骨架 + 安全基础 + 跑通完整链路。

| Task | Description | Est. |
|------|------------|------|
| 1.1.1 | 项目结构重组为微服务布局（go.work） | 2d |
| 1.1.2 | 共享库 pkg/（response, errcode, snowflake, middleware） | 1d |
| 1.1.3 | 配置中心（Viper + YAML + 环境变量 MB_ 前缀，启动时校验敏感配置） | 0.5d |
| 1.1.4 | MySQL 连接 + GORM 配置（deleted_at DATETIME 统一） | 0.5d |
| 1.1.5 | Redis 连接 + 分布式工具（锁、缓存） | 0.5d |
| 1.1.6 | Zap 结构化日志 + Request ID 中间件 | 0.5d |
| 1.1.7 | JWT 双 Token（Access + Refresh）+ Redis 黑名单 + 登出/吊销 | 2d |
| 1.1.8 | 登录速率限制（IP + 账号，Redis 限流器） | 0.5d |
| 1.1.9 | 密码复杂度校验 + 登录失败锁定 | 0.5d |
| 1.1.10 | CORS 白名单 + 安全响应头中间件 | 0.5d |
| 1.1.11 | 统一错误码映射（Handler errors.As） | 0.5d |
| 1.1.12 | context.Context 传播（不传 *gin.Context 到 Service 层） | 0.5d |
| 1.1.13 | godruNumin/gosnowflake（时钟回拨保护）替换 bwmarrin | 0.5d |
| 1.1.14 | Docker Compose 开发环境 | 0.5d |
| 1.1.15 | 集成测试骨架（httptest + 测试数据库） | 1d |

**Deliverable**: 可运行的 system-service 骨架，安全登录流程完整，测试框架就位。

---

### Milestone 1.2: Multi-Tenant Core

**Goal**: 实现独立数据库级别的多租户隔离。

| Task | Description | Est. |
|------|------------|------|
| 1.2.1 | 租户表 CRUD + API | 0.5d |
| 1.2.2 | GORM DBResolver 多库路由 | 1d |
| 1.2.3 | 租户连接池 LRU 管理（带 Prometheus 指标） | 1d |
| 1.2.4 | Tenant 中间件（JWT → DB 切换，所有查询注入 tenantScope） | 0.5d |
| 1.2.5 | 租户初始化流程（幂等：创建 DB → Schema 模板 → 状态机） | 1.5d |
| 1.2.6 | Schema 迁移工具（goose，支持多租户批量执行） | 1d |
| 1.2.7 | 租户健康检查 API（数据库连接 + Schema 完整性） | 0.5d |

**Deliverable**: 创建租户后自动创建独立数据库，请求自动路由，迁移工具就位。

---

### Milestone 1.3: RBAC & Data Permission

**Goal**: 完整的角色权限体系和数据权限控制。

| Task | Description | Est. |
|------|------------|------|
| 1.3.1 | Casbin 集成 + Redis Adapter + 多实例策略同步（Redis Pub/Sub） | 1.5d |
| 1.3.2 | 角色管理 CRUD + 菜单权限分配 API | 1d |
| 1.3.3 | 菜单管理（目录/菜单/按钮三级树，含 component/path/redirect/icon/visible） | 1d |
| 1.3.4 | 声明式权限中间件（路由级 permission 标识） | 0.5d |
| 1.3.5 | 部门树管理（含祖级列表维护） | 0.5d |
| 1.3.6 | 数据权限 GORM Scope（5 级范围） | 1d |
| 1.3.7 | 角色数据范围配置 + 自定义部门 | 0.5d |
| 1.3.8 | 岗位管理 CRUD | 0.5d |
| 1.3.9 | 前端动态菜单 API（GET /system/auth/get-permission-info） | 0.5d |

**Deliverable**: 不同角色看到不同菜单，不同数据范围看到不同数据。

---

### Milestone 1.4: System CRUD Modules

**Goal**: 完成所有 system 子模块的 CRUD 功能。

| Task | Description | Est. |
|------|------------|------|
| 1.4.1 | 用户管理 CRUD + 部门/角色/岗位分配 | 1d |
| 1.4.2 | 个人中心（修改信息/密码，改密吊销 Token） | 0.5d |
| 1.4.3 | 用户导入/导出（Excel StreamWriter，限 10 万行） | 0.5d |
| 1.4.4 | 字典类型 + 字典数据管理 | 0.5d |
| 1.4.5 | 参数管理（系统 KV 配置） | 0.5d |
| 1.4.6 | 通知公告管理 | 0.5d |
| 1.4.7 | 操作日志 + 登录日志（统一框架一起实现） | 1d |

**Deliverable**: 完整的 system 模块 CRUD 功能。

---

### Milestone 1.5: Advanced Features

**Goal**: 文件存储、在线用户、数据安全等高级功能。

| Task | Description | Est. |
|------|------------|------|
| 1.5.1 | 文件上传/下载 API | 0.5d |
| 1.5.2 | MinIO 存储适配（非默认凭据 + SSL） | 0.5d |
| 1.5.3 | 存储配置管理（动态切换后端） | 0.5d |
| 1.5.4 | 在线用户管理（WebSocket，coder/websocket 替代 gorilla） | 1d |
| 1.5.5 | 数据脱敏（struct tag + 反射声明式） | 0.5d |
| 1.5.6 | 数据加解密（AES/RSA，GORM Hook 透明） | 0.5d |
| 1.5.7 | 验证码登录 | 0.5d |
| 1.5.8 | 客户端管理 | 0.5d |

**Deliverable**: 高级功能完整可用。

---

### Milestone 1.6: Infrastructure Integration

**Goal**: 微服务基础设施集成，生产级部署能力。

| Task | Description | Est. |
|------|------------|------|
| 1.6.1 | gRPC protobuf 定义（UserService, RoleService） | 1d |
| 1.6.2 | gRPC 服务端 + 客户端实现（含健康检查 + 重连策略） | 1d |
| 1.6.3 | Consul 服务注册 + 健康检查 | 0.5d |
| 1.6.4 | Traefik 网关配置（路由 + TLS，JWT 验证留服务中间件） | 1d |
| 1.6.5 | Saga 框架（Redis Streams，状态持久化 MySQL，死信队列） | 2d |
| 1.6.6 | OpenTelemetry + Jaeger 链路追踪（采样策略） | 0.5d |
| 1.6.7 | 优雅关闭（SIGTERM → 排空事务 → 关闭连接，30s 超时） | 0.5d |
| 1.6.8 | 生产级 Docker Compose | 0.5d |
| 1.6.9 | API 文档（Swagger） | 0.5d |

**Deliverable**: 完整的微服务基础设施，可生产部署。

---

### Milestone 1.7: Observability & Operations

**Goal**: 监控、日志、备份等运维基础设施。

| Task | Description | Est. |
|------|------------|------|
| 1.7.1 | Prometheus + Grafana 监控（Go runtime/MySQL/Redis/Traefik） | 1d |
| 1.7.2 | Zap → Loki 日志收集（含 trace_id/tenant_id/request_id） | 0.5d |
| 1.7.3 | 租户级连接池监控（Prometheus 指标 + 分级告警） | 0.5d |
| 1.7.4 | 告警规则（连接池>80%、Redis内存>85%、错误率>5%） | 0.5d |
| 1.7.5 | MySQL 备份脚本（XtraBackup + binlog → MinIO） | 0.5d |
| 1.7.6 | 优雅关闭完善（Saga 状态持久化 + 租户连接排空） | 0.5d |

**Deliverable**: 生产可观测性和运维自动化。

---

### Milestone 1.8: Testing & Frontend Integration

**Goal**: 完整测试覆盖 + 前端联调。

| Task | Description | Est. |
|------|------------|------|
| 1.8.1 | Repository 层单元测试（SQLite 内存库，80%+ 覆盖） | 2d |
| 1.8.2 | Service 层单元测试（Mock Repository，80%+ 覆盖） | 1d |
| 1.8.3 | Handler 层集成测试（httptest，核心流程覆盖） | 1d |
| 1.8.4 | 多租户隔离测试（跨租户越权验证） | 1d |
| 1.8.5 | 安全测试（JWT 黑名单/速率限制/CORS/权限校验） | 1d |
| 1.8.6 | 前端联调（Vue3 + Vben5 动态菜单/权限路由/数据权限） | 3d |

**Deliverable**: 80%+ 测试覆盖，前端完整联调通过。

---

## Phase 2: Extended Modules (Future)

### Microservice Expansion
- job-service（定时任务管理、执行日志、分布式协调）
- workflow-service（n8n 集成、流程分类、流程实例管理，含 Circuit Breaker）

### Advanced Multi-Tenant
- 租户套餐管理
- 租户数据导出/清除

### Code Generator
- 微服务版脚手架生成（cmd/Dockerfile/Makefile/Protobuf/三层代码/SQL/前端 API）

### Extended Business
- CRM、ERP、商城、IoT、MES、AI（按需）

---

## Timeline Estimate

| Milestone | Duration | Dependencies |
|-----------|----------|-------------|
| 1.1 Foundation & Security | 11.5d | None |
| 1.2 Multi-Tenant Core | 6.5d | 1.1 |
| 1.3 RBAC & Permission | 7.5d | 1.2 |
| 1.4 System CRUD | 4.5d | 1.3 |
| 1.5 Advanced Features | 5d | 1.4 |
| 1.6 Infrastructure | 8d | 1.4 (parallel with 1.5) |
| 1.7 Observability | 3.5d | 1.6 |
| 1.8 Testing & Frontend | 9d | 1.5, 1.6 |
| **Phase 1 Total** | **~40d** | |

关键路径：1.1 → 1.2 → 1.3 → 1.4 → 1.5/1.6 → 1.7 → 1.8
1.5 和 1.6 可并行，缩短总工期。
