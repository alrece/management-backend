# Requirements — Management Backend (Go)

## Overview

Go 微服务管理后台，完整复刻 RuoYi-Vue-Plus 和 RuoYi-Vue-Pro 的 System 模块核心能力。采用微服务架构，支持完整多租户（独立数据库隔离）、RBAC 权限、数据权限、文件管理等企业级能力。

## Functional Requirements

### FR-01: Authentication & Authorization

| ID | Requirement | Priority | Notes |
|----|------------|----------|-------|
| FR-01-01 | JWT 双 Token（Access + Refresh） | P0 | Access 30min, Refresh 7d, Redis 存储 |
| FR-01-02 | JWT 黑名单（登出/禁用/改密吊销） | P0 | Redis jti 黑名单 + 用户维度批量吊销 |
| FR-01-03 | Casbin RBAC 权限控制 | P0 | Domain = tenant_id, Redis Adapter |
| FR-01-04 | 菜单权限校验（按钮级） | P0 | permission 标识匹配，声明式中间件 |
| FR-01-05 | 数据权限（5 级范围） | P0 | GORM Scope 实现，不在网关层 |
| FR-01-06 | 密码复杂度校验 | P0 | 最少 8 位，大小写+数字+特殊字符，可配置 |
| FR-01-07 | 登录失败锁定 | P0 | N 次失败后锁定 M 分钟，可配置 |
| FR-01-08 | 登录速率限制（IP + 账号） | P0 | Redis 限流器 |
| FR-01-09 | 验证码登录 | P1 | 图形验证码 |
| FR-01-10 | 社交登录（微信/钉钉） | P2 | OAuth2 框架 |
| FR-01-11 | 客户端管理（多端） | P1 | 预留 client_id/device_type 字段 |
| FR-01-12 | 前端动态菜单 API | P0 | 返回完整菜单树（path/component/icon/redirect/visible）+ 权限标识列表 |

### FR-02: User Management

| ID | Requirement | Priority | Notes |
|----|------------|----------|-------|
| FR-02-01 | 用户 CRUD | P0 | 含密码管理，所有单条查询必须注入 tenantScope |
| FR-02-02 | 分配部门/角色/岗位 | P0 | 多对多关联 |
| FR-02-03 | 用户导入/导出 | P1 | Excel（StreamWriter 防大文件 OOM） |
| FR-02-04 | 个人中心（修改信息/密码） | P0 | 改密时吊销所有 Token |
| FR-02-05 | 用户头像上传 | P1 | 文件存储集成 |

### FR-03: Organization Management

| ID | Requirement | Priority | Notes |
|----|------------|----------|-------|
| FR-03-01 | 部门树结构管理 | P0 | 含祖级列表维护 |
| FR-03-02 | 岗位管理 | P1 | 用户职务配置 |
| FR-03-03 | 角色管理 | P0 | 菜单权限 + 数据范围 |
| FR-03-04 | 角色菜单分配 | P0 | 树形勾选 |
| FR-03-05 | 角色数据范围配置 | P0 | 5 级 + 自定义部门 |

### FR-04: System Configuration

| ID | Requirement | Priority | Notes |
|----|------------|----------|-------|
| FR-04-01 | 菜单管理（目录/菜单/按钮） | P0 | 三级结构 + 权限标识 + component/path/redirect/icon/visible/cache |
| FR-04-02 | 字典类型管理 | P0 | 类型 + 数据维护 |
| FR-04-03 | 字典数据管理 | P0 | 按 type 查询 |
| FR-04-04 | 参数管理 | P1 | 系统动态配置 KV |
| FR-04-05 | 通知公告 | P1 | 信息发布与查看 |

### FR-05: Multi-Tenant

| ID | Requirement | Priority | Notes |
|----|------------|----------|-------|
| FR-05-01 | 租户管理（CRUD） | P0 | 创建/停用/过期 |
| FR-05-02 | 独立数据库隔离 | P0 | 每租户独立 DB |
| FR-05-03 | 租户数据自动路由 | P0 | JWT tenant_id → DB 切换 |
| FR-05-04 | 租户连接池管理 | P0 | GORM DBResolver + LRU 驱动空闲连接 |
| FR-05-05 | 租户初始化流程 | P0 | 创建租户 → 创建 DB → Schema 初始化（幂等、有状态机） |
| FR-05-06 | 租户数据导出/清除 | P1 | 退租/合规需求 |
| FR-05-07 | 租户套餐管理 | P2 | 功能模块授权 |
| FR-05-08 | 租户健康检查 | P1 | 数据库连接 + Schema 完整性巡检 |

### FR-06: Logging & Monitoring

| ID | Requirement | Priority | Notes |
|----|------------|----------|-------|
| FR-06-01 | 操作日志记录 | P0 | 正常 + 异常操作 |
| FR-06-02 | 登录日志记录 | P0 | 含异常登录检测，和操作日志一起实现 |
| FR-06-03 | 在线用户管理 | P1 | WebSocket + 强制踢出 |
| FR-06-04 | Request ID 中间件 | P0 | 全链路日志关联（Zap + trace_id） |
| FR-06-05 | 服务监控 | P2 | CPU/内存/磁盘/Go runtime |
| FR-06-06 | 缓存监控 | P2 | Redis 信息查询 |

### FR-07: File Management

| ID | Requirement | Priority | Notes |
|----|------------|----------|-------|
| FR-07-01 | 文件上传/下载/删除 | P1 | REST API |
| FR-07-02 | 存储配置管理 | P1 | 动态切换后端 |
| FR-07-03 | MinIO/S3 存储 | P1 | 默认存储后端，非默认凭据 + SSL |
| FR-07-04 | 本地存储 | P2 | 备选方案 |
| FR-07-05 | 大文件上传绕过网关 | P2 | 预签名 URL 直连 MinIO |

### FR-08: Data Security

| ID | Requirement | Priority | Notes |
|----|------------|----------|-------|
| FR-08-01 | 数据脱敏 | P1 | struct tag + 反射声明式 |
| FR-08-02 | 数据加解密 | P1 | AES/RSA/SM2/SM4，GORM Hook 透明处理 |
| FR-08-03 | 接口传输加密 | P2 | 动态 AES + RSA |
| FR-08-04 | 数据翻译 | P2 | 统一 TranslationService |

## Non-Functional Requirements

### NFR-01: Architecture

| ID | Requirement | Metric |
|----|------------|--------|
| NFR-01-01 | 微服务架构 | system/job/workflow 独立服务，独立部署 |
| NFR-01-02 | 混合通信 | 外部 REST，内部 gRPC |
| NFR-01-03 | API 网关 | Traefik 统一入口，JWT 验证留服务中间件 |
| NFR-01-04 | 服务发现 | Consul 自动注册 |
| NFR-01-05 | 分布式事务 | Saga 模式（Redis Streams），状态持久化到 MySQL |
| NFR-01-06 | 模块化结构 | 每个模块 model/repository/service/handler 独立 |

### NFR-02: API Compatibility

| ID | Requirement | Metric |
|----|------------|--------|
| NFR-02-01 | 统一响应格式 | `R{code, data, msg}` |
| NFR-02-02 | 统一错误码 | 按模块分段，Handler 用 errors.As 映射 |
| NFR-02-03 | 雪花 ID 主键 | 跨服务唯一（godruNumin/gosnowflake，时钟回拨保护） |
| NFR-02-04 | RuoYi 前端兼容 | API 路径和响应格式匹配 |
| NFR-02-05 | context.Context 传播 | 不将 *gin.Context 传到 Service/Repository 层 |

### NFR-03: Performance

| ID | Requirement | Metric |
|----|------------|--------|
| NFR-03-01 | 单接口响应 | < 200ms (P95) |
| NFR-03-02 | 并发用户 | >= 500 同时在线 |
| NFR-03-03 | 数据库连接池 | 可配置，每租户独立，全局预算上限 |
| NFR-03-04 | 大文件导出 | StreamWriter 流式写入，限 10 万行上限 |

### NFR-04: Reliability

| ID | Requirement | Metric |
|----|------------|--------|
| NFR-04-01 | 服务健康检查 | Consul 健康检查 + Traefik 自动摘除 |
| NFR-04-02 | 优雅关闭 | SIGTERM → 停接新请求 → 排空事务 → 关闭连接（30s 超时） |
| NFR-04-03 | 链路追踪 | OpenTelemetry + Jaeger（采样策略，生产 1-10%） |
| NFR-04-04 | Redis 持久化 | AOF + fsync everysec，Streams 数据不丢失 |

### NFR-05: Security

| ID | Requirement | Metric |
|----|------------|--------|
| NFR-05-01 | 敏感配置 | 密码/密钥通过环境变量注入，启动时校验非默认值 |
| NFR-05-02 | CORS 白名单 | 生产环境禁止 `*`，显式 Origin 列表 |
| NFR-05-03 | 安全响应头 | X-Content-Type-Options/X-Frame-Options/CSP/HSTS |
| NFR-05-04 | 跨租户校验 | 所有单条查询/删除/更新必须注入 tenantScope |
| NFR-05-05 | SQL 参数化 | GORM 参数化查询，禁止字符串拼接 |
| NFR-05-06 | 日志脱敏 | 生产环境 GORM logger 级别 Warn，不记录参数值 |

### NFR-06: Observability & Operations

| ID | Requirement | Metric |
|----|------------|--------|
| NFR-06-01 | 指标监控 | Prometheus + Grafana（Go runtime/MySQL/Redis/Traefik） |
| NFR-06-02 | 日志收集 | Zap JSON + Loki，含 trace_id/tenant_id/request_id |
| NFR-06-03 | Schema 迁移 | goose/golang-migrate 版本化管理，支持多租户批量执行 |
| NFR-06-04 | 数据库备份 | Percona XtraBackup + binlog，按租户粒度 |
| NFR-06-05 | 告警规则 | 连接池>80%、Redis内存>85%、错误率>5%、P99>1s |

## Constraints

- Go 1.22+, Gin, GORM, gRPC, Consul, Traefik, Redis, MySQL, MinIO
- 前端 Vue3 + Vben5
- **独立数据库级别多租户隔离**（决策不变）
- **微服务架构从第一天开始**（决策不变）
- Redis 分布式协调（锁/幂等/缓存/事件总线）
- 不使用 Go 微服务框架（Kratos/go-zero 等），纯手动组合
- 生产部署 Docker Compose（明确选择）

## Assumptions

- 管理后台用户量级为千级，非百万级
- 事件量级低（用户操作驱动），非高频数据流
- 部署环境为 Docker Compose
- 团队 Go 经验为中级，需要清晰的代码结构

## Out of Scope (Phase 2+)

- 代码生成器
- CRM / ERP / 商城 / IoT / MES / AI
- 支付（微信/支付宝）
- 短信/邮件融合包
- 微信公众号/小程序
- 国际化 i18n
