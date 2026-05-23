# Proposal: 自研类 n8n 可视化工作流引擎

## 问题陈述

当前工作流模块（前端 CRUD 页面）仅实现了基础的列表管理，缺少核心的可视化流程编排能力。需要参照 n8n 的设计理念，自研一套完整的可视化工作流引擎，包括：

1. **后端工作流引擎**（Go）— 节点注册、DAG 执行、触发器管理
2. **前端拖拽编辑器**（Vue Flow）— 可视化画布、节点面板、属性编辑器
3. **数据库切换** — 整个项目从 MySQL 迁移到 PostgreSQL

## 约束集合

### 硬约束

| 编号 | 约束 | 原因 |
|------|------|------|
| HC-01 | 工作流定义存储为 JSONB（PostgreSQL） | 前端直接存取 nodes/edges，避免 schema 冲突 |
| HC-02 | 节点执行基于 DAG 拓扑排序 | n8n 核心模型，保证执行顺序正确 |
| HC-03 | 触发器类型仅支持 Cron / Webhook / 手动 | 复用现有 robfig/cron，避免过度设计 |
| HC-04 | 所有 API 遵循现有三层架构（Handler→Service→Repository） | 与项目架构一致 |
| HC-05 | 主键使用雪花 ID（BIGINT），禁止 AUTO_INCREMENT | 项目规范 |
| HC-06 | 多租户支持（tenant_id） | 项目规范 |
| HC-07 | RBAC 权限集成（Casbin） | 项目规范 |
| HC-08 | 数据库统一使用 PostgreSQL | 用户决策 |
| HC-09 | 前端使用 @vue-flow/core 构建画布 | 用户选择自研拖拽画布 |
| HC-10 | Go 1.22+, Gin, GORM, Redis | 项目技术栈 |

### 软约束

| 编号 | 约束 | 原因 |
|------|------|------|
| SC-01 | 前端节点组件使用 Ant Design Vue 风格 | 与现有 UI 统一 |
| SC-02 | 初始节点类型不超过 10 种 | YAGNI 原则 |
| SC-03 | 工作流版本管理初始版本仅保留最近 10 个版本 | 避免数据膨胀 |
| SC-04 | 执行日志保留 30 天 | 磁盘空间控制 |

## 技术决策

### 1. 工作流定义格式

前端 Vue Flow 格式 → 后端存储格式映射：

```
前端 (Vue Flow):
  nodes: [{ id, type, position: {x,y}, data: { label, config, ... } }]
  edges: [{ id, source, target, sourceHandle, targetHandle }]

后端存储 (JSONB):
  definition: {
    nodes: [...],     // 保留前端格式
    edges: [...],     // 保留前端格式
    version: 1,
    settings: {}
  }
```

### 2. 节点类型体系（Go）

```
NodeHandler 接口:
  - Execute(ctx, input) → output
  - Type() string
  - Category() string
  - Schema() NodeSchema    // 前端配置面板用

初始节点类型:
  触发器: CronTrigger, WebhookTrigger, ManualTrigger
  动作:   HttpRequest, SqlQuery, SendEmail, LogMessage
  控制:   Condition, Delay, Loop
  转换:   Transform, Filter
```

### 3. DAG 执行引擎

```
执行流程:
  1. 加载工作流定义（nodes + edges）
  2. 拓扑排序确定执行顺序
  3. 并行执行无依赖分支（goroutine + sync.WaitGroup）
  4. 每个 NodeInput 包含上游节点输出数据
  5. 错误处理：单节点失败可配置 重试/跳过/终止
  6. 执行结果写入 wf_execution + wf_node_log
```

### 4. 前端编辑器架构

```
┌──────────────────────────────────────────────────┐
│  工具栏: 保存 | 撤销 | 重做 | 执行 | 调试 | 返回  │
├────────┬──────────────────────┬───────────────────┤
│        │                      │                   │
│  节点  │    Vue Flow 画布      │   属性面板        │
│  面板  │    (拖拽/连线)        │   (选中节点配置)  │
│        │                      │                   │
│ ──触发 │                      │  节点名称         │
│  Cron  │   [Cron]──►[IF]     │  参数配置         │
│  Webhook│     │       │       │  错误处理         │
│  手动  │   [HTTP]  [Log]     │  重试策略         │
│        │                      │                   │
│ ──动作 │                      │                   │
│  HTTP  │                      │                   │
│  SQL   │                      │                   │
│  邮件  │                      │                   │
│        │                      │                   │
│ ──控制 │                      │                   │
│  条件  │                      │                   │
│  延迟  │                      │                   │
│        │                      │                   │
└────────┴──────────────────────┴───────────────────┘
```

### 5. 数据库 Schema（PostgreSQL）

```sql
-- 工作流定义
CREATE TABLE wf_workflow (
    id BIGINT PRIMARY KEY,
    tenant_id BIGINT NOT NULL DEFAULT 0,
    category_id BIGINT,
    name VARCHAR(200) NOT NULL,
    description TEXT DEFAULT '',
    definition JSONB NOT NULL DEFAULT '{}',
    status SMALLINT NOT NULL DEFAULT 0,  -- 0=草稿 1=已激活 2=已停用
    version INT NOT NULL DEFAULT 1,
    creator BIGINT,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updater BIGINT,
    update_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted SMALLINT NOT NULL DEFAULT 0
);

-- 执行实例
CREATE TABLE wf_execution (
    id BIGINT PRIMARY KEY,
    tenant_id BIGINT NOT NULL DEFAULT 0,
    workflow_id BIGINT NOT NULL,
    workflow_version INT NOT NULL,
    status SMALLINT NOT NULL DEFAULT 0, -- 0=运行中 1=成功 2=失败 3=取消
    trigger_type VARCHAR(20) NOT NULL,   -- cron/webhook/manual
    input_data JSONB DEFAULT '{}',
    output_data JSONB DEFAULT '{}',
    error_message TEXT DEFAULT '',
    started_at TIMESTAMPTZ,
    finished_at TIMESTAMPTZ,
    duration_ms INT DEFAULT 0,
    creator BIGINT,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 节点执行日志
CREATE TABLE wf_node_log (
    id BIGINT PRIMARY KEY,
    tenant_id BIGINT NOT NULL DEFAULT 0,
    execution_id BIGINT NOT NULL,
    workflow_id BIGINT NOT NULL,
    node_id VARCHAR(100) NOT NULL,
    node_type VARCHAR(50) NOT NULL,
    node_name VARCHAR(200),
    status SMALLINT NOT NULL DEFAULT 0, -- 0=运行中 1=成功 2=失败 3=跳过
    input_data JSONB DEFAULT '{}',
    output_data JSONB DEFAULT '{}',
    error_message TEXT DEFAULT '',
    started_at TIMESTAMPTZ,
    finished_at TIMESTAMPTZ,
    duration_ms INT DEFAULT 0,
    retry_count SMALLINT DEFAULT 0,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 触发器配置
CREATE TABLE wf_trigger (
    id BIGINT PRIMARY KEY,
    tenant_id BIGINT NOT NULL DEFAULT 0,
    workflow_id BIGINT NOT NULL,
    trigger_type VARCHAR(20) NOT NULL, -- cron/webhook
    config JSONB NOT NULL DEFAULT '{}', -- cron表达式/webhook路径等
    status SMALLINT NOT NULL DEFAULT 0, -- 0=启用 1=停用
    creator BIGINT,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updater BIGINT,
    update_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted SMALLINT NOT NULL DEFAULT 0
);

-- 工作流分类（保留现有树形结构）
CREATE TABLE wf_category (
    id BIGINT PRIMARY KEY,
    tenant_id BIGINT NOT NULL DEFAULT 0,
    parent_id BIGINT NOT NULL DEFAULT 0,
    name VARCHAR(100) NOT NULL,
    sort INT NOT NULL DEFAULT 0,
    status SMALLINT NOT NULL DEFAULT 0,
    creator BIGINT,
    create_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updater BIGINT,
    update_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted SMALLINT NOT NULL DEFAULT 0
);
```

## 依赖与风险

### 依赖

| 依赖 | 状态 | 影响 |
|------|------|------|
| PostgreSQL 数据库 | 需部署 | 全局影响，需先完成 MySQL→PG 迁移 |
| @vue-flow/core | npm 包 | 前端画布 |
| @vue-flow/background | npm 包 | 画布背景网格 |
| @vue-flow/minimap | npm 包 | 小地图 |
| @vue-flow/controls | npm 包 | 缩放/居中控制 |
| dagre 或 elkjs | npm 包 | 自动布局算法 |
| robfig/cron/v3 | Go 包 | 已有，触发器复用 |

### 风险

| 风险 | 等级 | 缓解措施 |
|------|------|---------|
| MySQL→PG 迁移影响现有模块 | 高 | 分阶段迁移，先 PG 适配层 |
| Vue Flow 自定义节点复杂度 | 中 | 从简单节点开始，逐步增强 |
| DAG 执行引擎并发安全 | 高 | 充分的单元测试 + 竞态检测 |
| 工作流定义 JSONB 性能 | 低 | PostgreSQL JSONB 有索引支持 |

## 成功标准

1. 用户可通过拖拽在画布上创建/编辑工作流
2. 支持 Cron / Webhook / 手动 三种触发方式
3. 工作流执行结果可查看（成功/失败/耗时）
4. 节点级别执行日志可追溯
5. 分类管理（树形）功能保留
6. 权限控制集成到现有 RBAC
7. 多租户隔离正常工作

## 影响范围

- **新建模块**: `internal/module/workflow/`（后端完整工作流模块）
- **替换页面**: 前端 `views/workflow/` 全部重写
- **全局影响**: MySQL→PostgreSQL 迁移（config.yaml, GORM driver, 建表语句）
- **路由变更**: `router/router.go` 新增工作流路由注册
- **菜单变更**: `scripts/sql/menu-workflow.sql` 更新（增加编辑器页面）
