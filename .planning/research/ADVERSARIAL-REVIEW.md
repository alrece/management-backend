# Adversarial Review Report

> 2026-05-18: 5 个角色（架构师/安全专家/运维工程师/产品经理/Go语言专家）对抗审查，共 67 个质疑。

## 核心共识（5 方一致）

### 1. Phase 1 微服务基础设施零收益

Phase 1 只有 system-service 一个服务。Traefik/Consul/gRPC/Saga 在单服务场景下不产生业务价值。Milestone 1.6 的 6 天（24%工期）应推迟到 Phase 2。

**建议**: Phase 1 做模块化单体，`internal/module/system/` 结构天然支持未来拆分。

### 2. 独立数据库多租户对目标用户过度设计

REQUIREMENTS.md 假设"千级用户量"，独立数据库的运维成本（连接池、Schema迁移、备份恢复、灾难恢复）与收益不匹配。PITFALLS.md 自身就警告了"连接池爆炸"。

**建议**: P0 用行级隔离（tenant_id + GORM Scope），独立数据库作为 P2 可选方案。可节省 Milestone 1.2 中 2.5 天的 DBResolver/LRU 工作量。

### 3. 25 天工期严重低估

不含测试编写/执行、前端联调、文档、运维基础设施、缓冲。业界经验：编码占 40-50%，测试 25-30%，联调 15%，返工 10-15%。

**建议**: 裁剪范围（模块化单体 + 行级隔离），或按 2x 系数重新估算为 40-50 天。

### 4. 现有代码与规划严重脱节

| 维度 | 现有代码 | 规划文档 |
|------|---------|---------|
| 架构 | 单体 | 微服务 |
| 多租户 | 行级隔离 | 独立数据库 |
| 日志 | 无 Zap | Zap 结构化 |
| 认证 | JWT 直接验证 | Traefik ForwardAuth |
| 错误处理 | 硬编码 500 | 统一错误码 |

**建议**: 承认当前为原型，Phase 1 从骨架重写，保留模块化结构。

### 5. 安全基础缺失

JWT 无黑名单、无速率限制、CORS 全开、密码策略缺失、GetByID/Delete 跨租户越权。

## 架构级调整建议

| 项 | 当前 | 建议 | 节省 |
|----|------|------|------|
| 架构模式 | 微服务从第一天 | 模块化单体 | 6 天 |
| 多租户 | 独立数据库 | 行级隔离 | 2.5 天 |
| 租户套餐 | P0 | P2 | 0.5 天 |
| Saga 框架 | M1.6 | Phase 2 | 1 天 |
| 登录日志 | P1 | P0 | +0 天 |
| 密码/登录策略 | 无 | P0 | +1 天 |
| 前端联调 | 无 | M1.7 | +3 天 |
| 测试 | 无 | 每里程碑+1天 | +6 天 |

## 优先级重定义

### 新增 P0 需求（从审查中识别）

| 需求 | 来源 | 说明 |
|------|------|------|
| 密码复杂度校验 | 安全/产品 | 最少 8 位，大小写+数字+特殊字符 |
| 登录失败锁定 | 安全/产品 | N 次失败后锁定 M 分钟 |
| JWT 黑名单（Redis） | 安全 | 登出/禁用/密码修改时吊销 Token |
| Refresh Token 轮换 | 安全 | 双 Token 方案，refresh 一次性 |
| 登录日志 | 产品 | 和操作日志一起实现 |
| 前端动态菜单 API | 产品 | 返回完整菜单树 + 权限标识 |
| Request ID 中间件 | 运维/Go | 全链路日志关联 |
| 统一错误码映射 | Go | Handler 层 errors.As 提取业务错误码 |

### 降级需求

| 需求 | 原优先级 | 新优先级 | 原因 |
|------|---------|---------|------|
| 租户套餐管理 | P0 | P2 | Phase 1 只有一个模块，无功能开关需求 |
| 独立数据库隔离 | P0 | P2 | 行级隔离足够，运维成本不匹配 |
| gRPC protobuf | M1.6 | Phase 2 | 单服务无跨服务调用 |
| Consul 服务注册 | M1.6 | Phase 2 | 单服务不需要服务发现 |
| Traefik ForwardAuth | M1.6 | Phase 2 | JWT 验证留在服务中间件更简单 |
| Saga 框架 | M1.6 | Phase 2 | 单服务用 GORM Transaction 即可 |
| 配置中心 Consul KV | M1.3 | Phase 2 | YAML + 环境变量足够 |

## Go 技术关键修复

| 问题 | 风险 | 修复 |
|------|------|------|
| `gorm.DeletedAt` vs `BIT(1)` | 数据不兼容 | 统一为 `deleted_at DATETIME` |
| `bwmarrin/snowflake` 时钟回拨 | NTP 同步时 panic | 切换到 `godruNumin/gosnowflake` |
| `gorilla/websocket` 已归档 | CVE 无人修复 | 迁移到 `coder/websocket` |
| `*gin.Context` 泄漏到 Service | gRPC 复用障碍 | `context.WithValue` 传递 |
| Handler 错误全返回 500 | 前端无法区分错误类型 | `errors.As` + 统一映射 |
| `cors.Default()` | CSRF 攻击 | 显式白名单 |
| `log.Printf` 替代 Zap | 无结构化日志 | 全局替换为 Zap |
| `sync.Once` 全局 snowflake | 测试困难 | IDGenerator 接口注入 |
| excelize 大文件 | OOM | StreamWriter API |
| Redis Streams 无死信队列 | 消息静默丢失 | 实现 dead letter stream |

## 安全 P0 修复清单

| # | 问题 | CVSS | 修复 |
|---|------|------|------|
| 1 | GetByID/Delete 跨租户越权 | 9.1 | 注入 tenantScope |
| 2 | JWT 无黑名单 | 8.2 | Redis jti 黑名单 |
| 3 | JWT 密钥硬编码 | 9.8 | 环境变量 + 启动校验 |
| 4 | MySQL root/123456 | 8.7 | 专用用户 + 环境变量 |
| 5 | Redis 无密码 | 8.6 | 强密码 + ACL |
| 6 | MinIO 默认凭据 | 8.2 | 非默认凭据 + SSL |
| 7 | CORS 默认配置 | 8.1 | 白名单 |
| 8 | 登录无速率限制 | 7.3 | IP 限流 + 账号锁定 |
| 9 | Refresh Token 缺失 | 7.5 | 双 Token + Redis 存储 |
| 10 | 密码强度过低 | 5.9 | 8位+复杂度 |

## 运维关键补充

| 项 | 当前状态 | 建议 |
|----|---------|------|
| Schema 迁移 | 无方案 | 引入 goose/migrate，版本化管理 |
| 监控告警 | P2 推迟 | Prometheus + Grafana 纳入 Phase 1 |
| 日志收集 | 无方案 | Loki + trace_id 关联 |
| 备份恢复 | 无方案 | XtraBackup + binlog + 定期演练 |
| 灰度发布 | 无方案 | Schema 兼容性规则 + 蓝绿部署 |
| 生产部署决策 | Docker Compose vs K8s "可选" | 明确选择 Docker Compose |

## 修正后的路线图建议

```
Phase 1A: 骨架重写（模块化单体）               3d
Phase 1B: 安全基础（JWT双Token/黑名单/限流/密码策略） 3d
Phase 1C: 多租户（行级隔离）                   2d
Phase 1D: RBAC + 数据权限（Casbin/5级Scope）       5d
Phase 1E: System CRUD（用户/角色/菜单/部门/字典/日志） 5d
Phase 1F: 高级功能（文件/在线用户/脱敏/加解密）     4d
Phase 1G: 测试 + 文档 + 前端联调                5d
                                      合计: ~27d

Phase 2A: 微服务拆分（gRPC/Consul/Traefik）      6d
Phase 2B: 独立数据库多租户（可选）              5d
Phase 2C: Job Service                      5d
Phase 2D: Saga 框架                       3d
Phase 2E: Workflow Service（n8n集成）        8d
Phase 2F: 运维基础设施（监控/日志/备份/CI）      5d
Phase 2G: 代码生成器                       5d
```
