# Tasks: 自研类 n8n 可视化工作流引擎

## Phase 0: PostgreSQL 迁移（前置条件）

- [x] 0.1 安装 gorm postgres driver 依赖，移除 mysql driver
- [x] 0.2 更新 config.yaml 配置格式（mysql → postgres）
- [x] 0.3 更新 internal/config/config.go 配置结构体
- [x] 0.4 更新 cmd/system/main.go 数据库连接逻辑（DSN 格式）
- [x] 0.5 创建 scripts/sql/pg_init.sql（PostgreSQL 建表脚本）
- [x] 0.6 更新 BaseEntity/BaseModel 的 deleted 字段适配 PG
- [x] 0.7 验证现有 system 模块在 PG 下正常运行

## Phase 1: 后端工作流引擎核心

- [x] 1.1 创建 internal/module/workflow/ 目录结构
- [x] 1.2 定义 model/entity.go（Workflow, Execution, NodeLog, Trigger, Category）
- [x] 1.3 定义 model/dto.go（请求/响应 DTO）
- [x] 1.4 实现 engine/registry.go（NodeRegistry + NodeHandler 接口 + NodeSchema）
- [x] 1.5 实现 engine/dag.go（拓扑排序 + 循环检测）
- [x] 1.6 实现 engine/context.go（NodeInput, NodeOutput 结构体）
- [x] 1.7 实现内置节点：engine/nodes/http_request.go
- [x] 1.8 实现内置节点：engine/nodes/sql_query.go
- [x] 1.9 实现内置节点：engine/nodes/condition.go
- [x] 1.10 实现内置节点：engine/nodes/delay.go
- [x] 1.11 实现内置节点：engine/nodes/transform.go
- [x] 1.12 实现内置节点：engine/nodes/log_message.go
- [x] 1.13 实现内置节点：engine/nodes/send_email.go
- [x] 1.14 实现 service/executor.go（DAG 执行引擎，goroutine pool，日志记录）
- [x] 1.15 实现 repository/workflow.go（工作流 CRUD，JSONB 查询）
- [x] 1.16 实现 repository/execution.go（执行记录查询）
- [x] 1.17 实现 repository/node_log.go（节点日志查询与清理）
- [x] 1.18 实现 repository/trigger.go（触发器配置管理）
- [x] 1.19 实现 service/trigger.go（Cron 注册/注销，Webhook 路径生成）
- [x] 1.20 实现 repository/category.go（分类树查询）
- [x] 1.21 实现 service/category.go（分类 CRUD）
- [x] 1.22 实现 service/workflow.go（工作流业务逻辑，含循环检测）
- [x] 1.23 实现 handler/workflow.go（工作流 API + 节点 Schema 接口）
- [x] 1.24 实现 handler/execution.go（执行实例 API）
- [x] 1.25 实现 handler/webhook.go（Webhook 触发入口）
- [x] 1.26 实现 handler/category.go（分类 API）
- [x] 1.27 实现 router.go（路由注册）
- [x] 1.28 在 internal/router/router.go 中注入工作流模块依赖
- [x] 1.29 创建 scripts/sql/pg_workflow.sql（工作流相关建表脚本）
- [x] 1.30 编写 engine 包单元测试（拓扑排序、循环检测、节点执行）

## Phase 2: 前端工作流编辑器

- [x] 2.1 安装 Vue Flow 相关依赖（@vue-flow/core, background, minimap, controls, @dagrejs/dagre）
- [x] 2.2 创建 api/workflow/workflow/index.ts（工作流 API，含节点 Schema 获取）
- [x] 2.3 创建 api/workflow/instance/index.ts（执行实例 API）
- [x] 2.4 创建 api/workflow/category/index.ts（分类 API，保留现有）
- [x] 2.5 创建 views/workflow/editor/types.ts（编辑器类型定义）
- [x] 2.6 创建 views/workflow/editor/composables/useDragAndDrop.ts（拖拽逻辑）
- [x] 2.7 创建 views/workflow/editor/composables/useWorkflow.ts（工作流状态管理）
- [x] 2.8 创建 views/workflow/editor/composables/useHistory.ts（撤销/重做）
- [x] 2.9 创建 views/workflow/editor/composables/useNodeSchema.ts（Schema 获取）
- [x] 2.10 创建 views/workflow/editor/components/nodes/TriggerNode.vue（触发器节点）
- [x] 2.11 创建 views/workflow/editor/components/nodes/ActionNode.vue（动作节点）
- [x] 2.12 创建 views/workflow/editor/components/nodes/ControlNode.vue（控制节点）
- [x] 2.13 创建 views/workflow/editor/components/nodes/TransformNode.vue（转换节点）
- [x] 2.14 创建 views/workflow/editor/components/edges/DataFlowEdge.vue（数据流边）
- [x] 2.15 创建 views/workflow/editor/components/NodePalette.vue（节点面板）
- [x] 2.16 创建 views/workflow/editor/components/PropertyPanel.vue（属性面板）
- [x] 2.17 创建 views/workflow/editor/components/EditorToolbar.vue（工具栏）
- [x] 2.18 创建 views/workflow/editor/components/FlowCanvas.vue（Vue Flow 画布）
- [x] 2.19 创建 views/workflow/editor/index.vue（编辑器主页面）
- [x] 2.20 重写 views/workflow/category/（分类管理页面保留/优化）
- [x] 2.21 重写 views/workflow/instance/（执行实例列表 + 详情查看）
- [x] 2.22 更新 vite.config.ts（Mock 数据、代理配置、菜单节点）

## Phase 3: 集成与菜单

- [x] 3.1 更新 scripts/sql/menu-workflow.sql（添加编辑器页面菜单项和权限）
- [x] 3.2 更新 vite.config.ts Mock 菜单节点（编辑器路由）
- [x] 3.3 前端路由配置（/workflow/editor/:id）
- [x] 3.4 权限字符串定义和 Casbin 策略初始化
- [x] 3.5 执行日志清理 CronJob 实现

## Phase 4: 测试验证

- [x] 4.1 后端单元测试：DAG 执行引擎（拓扑排序、并行执行、错误策略）
- [x] 4.2 后端单元测试：节点注册表和各内置节点
- [x] 4.3 后端集成测试：完整工作流 CRUD → 执行 → 日志查询
- [x] 4.4 前端手动验证：编辑器拖拽、连线、保存、执行
