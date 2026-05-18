# Management Backend（Go 管理后台）

## What This Is

基于 Go 技术栈的通用型管理后台后端，完整复刻 RuoYi-Vue-Plus 和 RuoYi-Vue-Pro 两个 Java 项目的核心能力。采用微服务架构，每个业务模块独立服务，支持完整多租户、RBAC 权限、工作流、代码生成等企业级能力。前端采用 Vue3 + Vben5。面向需要快速搭建管理后台的中小型团队和企业。

## Core Value

提供一套开箱即用的 Go 微服务管理后台框架，让开发者只需关注业务模块本身，基础设施能力（认证、权限、多租户、文件存储、任务调度等）全部内置。

## Requirements

### Validated

(None yet — ship to validate)

### Active

#### 基础设施（P0）

- [ ] 微服务架构基座：每个模块独立服务，独立部署、独立配置
- [ ] API 网关（Traefik）：统一入口，负责路由、认证、限流、日志、数据权限拦截
- [ ] 服务注册发现（Consul）：服务自动注册、健康检查、负载均衡
- [ ] 服务间通信：外部 REST + 内部 gRPC 混合模式
- [ ] 分布式事务（Saga 模式）：跨服务操作的最终一致性保障
- [ ] 统一响应格式 R{code, data, msg}
- [ ] 统一错误码体系（按模块分段）
- [ ] 雪花 ID 主键生成（跨服务唯一）
- [ ] 配置中心（Viper + Consul 集中配置）
- [ ] Redis 分布式协调：分布式锁、分布式幂等、缓存

#### System 模块（P0）

- [ ] 用户管理：增删改查、分配部门/角色/岗位、密码管理
- [ ] 部门管理：树结构组织、数据权限绑定
- [ ] 岗位管理：用户职务配置
- [ ] 菜单管理：目录/菜单/按钮三级、权限标识
- [ ] 角色管理：菜单权限分配、数据范围（全部/自定义/本部门/本部门及以下/仅本人）
- [ ] 字典管理：类型+数据维护
- [ ] 参数管理：系统动态配置
- [ ] 通知公告：信息发布与查看
- [ ] 操作日志：正常+异常操作记录
- [ ] 登录日志：登录记录含异常
- [ ] 在线用户管理：监控与强制踢出
- [ ] 完整多租户：租户管理、租户套餐、数据隔离（独立数据库级别）
- [ ] 客户端管理：PC/小程序等多端管理，动态授权登录方式
- [ ] JWT 认证 + Casbin RBAC 权限控制
- [ ] 数据脱敏（注解/tag 级别）
- [ ] 数据加解密（AES/RSA/SM2/SM4）
- [ ] 接口传输加密（动态 AES + RSA）
- [ ] 数据翻译（序列化期间动态修改）
- [ ] 三方登录（微信/钉钉等）

#### Job 模块（P1）

- [ ] 定时任务管理：添加/修改/删除/暂停
- [ ] 任务执行日志
- [ ] 执行器管理
- [ ] 分布式任务协调（Redis + Consul）

#### Workflow 模块（P1）

- [ ] 集成外部工作流服务 n8n（Go 服务通过 API 调用）
- [ ] 流程分类管理
- [ ] 流程设计（委托 n8n 的流程编辑器）
- [ ] 流程实例管理：发起/审批/转办/委派/加减签/会签/或签
- [ ] 流程表单关联

#### Demo 模块（P1）

- [ ] 框架功能使用案例（单表/树表/主子表 CRUD 示例）
- [ ] 各类组件用法展示

#### 文件与存储（跨模块）

- [ ] 文件管理：上传/下载/删除
- [ ] 文件配置管理：动态切换存储后端
- [ ] MinIO/S3 协议对象存储
- [ ] 本地/阿里云/腾讯云/七牛云存储适配

#### 监控与运维

- [ ] 服务监控（CPU/内存/磁盘/Go runtime）
- [ ] 缓存监控（Redis 信息查询/命令统计）
- [ ] SQL 监控（GORM logger 完整 SQL 输出）
- [ ] 链路追踪（OpenTelemetry + Jaeger）
- [ ] 在线日志查看

#### 前端对接

- [ ] Vue3 + Vben5 管理后台（参考 yudao-ui-admin-vben）
- [ ] 适配 Go 后端 API（REST）
- [ ] 暗色/亮色主题
- [ ] 响应式布局

### Out of Scope

- **代码生成器** — 二阶段实现（需生成完整服务脚手架：cmd/Dockerfile/Makefile/Protobuf/三层代码/SQL/前端 API）
- **CRM 客户管理** — ruoyi-vue-pro 独有，二阶段按需
- **ERP 进销存** — ruoyi-vue-pro 独有，二阶段按需
- **商城系统**（商品/营销/交易） — ruoyi-vue-pro 独有，二阶段按需
- **IoT 物联网** — ruoyi-vue-pro 独有，二阶段按需
- **MES 制造执行** — ruoyi-vue-pro 独有，二阶段按需
- **AI 大模型集成** — ruoyi-vue-pro 独有，二阶段按需
- **微信公众号/小程序** — 二阶段按需
- **支付宝/微信支付** — 二阶段按需
- **短信/邮件融合包** — 二阶段按需
- **国际化 i18n** — 二阶段

## Context

- **参考项目 1**：`ruoyi-vue-plus-ai-engineering/`（Spring Boot 3.5 + JDK 17，Sa-Token，MyBatis-Plus，warm-flow）
- **参考项目 2**：`ruoyi-vue-pro/`（Spring Boot 2.7 + JDK 8，Spring Security，MyBatis-Plus，Flowable）
- 两个 Java 项目均为单体架构，Go 版本需切换到微服务思维
- Go 骨架已搭建（Gin + GORM + JWT），当前为单体结构，需拆分为微服务
- Go 生态无 warm-flow/Flowable 等价库，工作流采用外部服务 n8n
- 前端参考：https://github.com/yudaocode/yudao-ui-admin-vben（Vue3 + Vben5）

## Constraints

- **Tech Stack**: Go 1.22+、Gin、GORM、gRPC、Consul、Traefik、Redis、MySQL、MinIO、n8n
- **Architecture**: 从一开始就是微服务
- **Multi-tenant**: 独立数据库级别隔离
- **Communication**: 混合模式（外部 REST，内部 gRPC）
- **Distributed Transactions**: Saga 模式
- **Coordination**: Redis
- **Service Discovery**: Consul
- **API Gateway**: Traefik
- **Frontend**: Vue3 + Vben5（通过 Traefik 网关调用后端）
- **Compatibility**: API 兼容 Java 项目前端调用模式

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| 微服务从一开始 | 项目定位是通用框架，单体再拆成本更高 | — Pending |
| 混合通信（REST + gRPC） | 前端兼容 REST，内部高性能 gRPC | — Pending |
| Saga 分布式事务 | 跨服务操作不可避免，Saga 是成熟方案 | — Pending |
| Consul 服务发现 | Go 生态支持好，与 Traefik 原生集成 | — Pending |
| Traefik 网关 | Go 原生、配置驱动、自动服务发现 | — Pending |
| n8n 外部工作流 | Go 无等价引擎，n8n 功能完善可独立部署 | — Pending |
| 数据权限在网关层 | 微服务各服务独立数据库，需统一拦截 | — Pending |
| 独立数据库多租户 | 最强隔离级别，适合 SaaS | — Pending |
| Redis 分布式协调 | 已引入 Redis，复用为锁/幂等/缓存基础 | — Pending |
| 代码生成器推迟到二阶段 | 微服务版生成器复杂度高 | — Pending |
| 先做 system 模块 | 所有模块依赖 system 的用户/权限/租户 | — Pending |

## Evolution

This document evolves at phase transitions and milestone boundaries.

**After each phase transition** (via `/gsd:transition`):
1. Requirements invalidated? → Move to Out of Scope with reason
2. Requirements validated? → Move to Validated with phase reference
3. New requirements emerged? → Add to Active
4. Decisions to log? → Add to Key Decisions
5. "What This Is" still accurate? → Update if drifted

**After each milestone** (via `/gsd:complete-milestone`):
1. Full review of all sections
2. Core Value check — still the right priority?
3. Audit Out of Scope — reasons still valid?
4. Update Context with current state

---
*Last updated: 2026-05-18 after initialization*
