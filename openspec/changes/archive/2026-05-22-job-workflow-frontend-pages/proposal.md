## Why

Phase 02 已实现 job-service（定时任务）和 workflow-service（工作流）两个独立微服务的后端 REST API，但前端原型中缺少对应的管理页面。用户无法通过 UI 管理定时任务（CRUD、手动触发、执行日志）和工作流（分类管理、流程 CRUD、实例监控、执行操作）。

## What Changes

### 定时任务管理页面

1. **任务列表页** (`views/job/task/index.vue`)
   - 分页表格：名称、处理器、Cron 表达式、状态、参数、备注
   - 操作：新建、编辑、删除、暂停/恢复、手动触发
   - 搜索：按名称、状态筛选

2. **执行日志页** (`views/job/log/index.vue`)
   - 分页表格：任务名称、触发类型（定时/手动）、状态（执行中/成功/失败/超时）、耗时、结果/错误
   - 搜索：按任务、状态、时间范围筛选

3. **API 层** (`api/job/task/index.ts`, `api/job/log/index.ts`)
   - 任务 CRUD + 手动触发：对应后端 `/job/tasks/*` 路由
   - 执行日志分页：对应后端 `/job/execution-logs/page`

### 工作流管理页面

4. **分类管理** (`views/workflow/category/index.vue`)
   - 树形展示分类层级
   - 操作：新建、编辑、删除分类节点

5. **工作流列表页** (`views/workflow/workflow/index.vue`)
   - 分页表格：名称、分类、状态（DRAFT/ACTIVE/INACTIVE）、创建时间
   - 操作：新建、编辑、删除、激活、停用、执行
   - 搜索：按名称、分类、状态筛选

6. **执行实例页** (`views/workflow/instance/index.vue`)
   - 分页表格：工作流名称、状态（RUNNING/SUCCESS/FAILED/CANCELLED）、耗时、错误信息
   - 搜索：按工作流、状态筛选

7. **API 层** (`api/workflow/category/index.ts`, `api/workflow/workflow/index.ts`, `api/workflow/instance/index.ts`)
   - 分类 CRUD + 树：对应后端 `/workflow/categories/*`
   - 工作流 CRUD + 激活/停用/执行：对应后端 `/workflow/workflows/*`
   - 实例分页查询：对应后端 `/workflow/instances/page`

### 菜单与路由

8. 数据库菜单表插入以下菜单结构（通过后端 API 或 SQL 脚本）：
   ```
   定时任务管理（目录）
   ├── 任务管理（页面，/job/task）
   └── 执行日志（页面，/job/log）
   工作流管理（目录）
   ├── 流程分类（页面，/workflow/category）
   ├── 流程管理（页面，/workflow/workflow）
   └── 执行实例（页面，/workflow/instance）
   ```

## Capabilities

### New Capabilities

- `job-frontend`: 定时任务管理前端页面（任务列表 + 执行日志），对接 job-service REST API
- `workflow-frontend`: 工作流管理前端页面（分类树 + 工作流列表 + 实例监控），对接 workflow-service REST API

## Impact

### 代码影响

- **frontend/apps/web-antd/src/views/**: 新增 `job/` 和 `workflow/` 目录，包含 5 个页面视图
- **frontend/apps/web-antd/src/api/**: 新增 `job/` 和 `workflow/` 目录，包含 5 个 API 模块
- **scripts/sql/**: 新增菜单初始化 SQL（sys_menu 插入记录）

### API 影响

- 前端直接调用 job-service 和 workflow-service 的 REST API
- 需确认前端代理配置（vite.config.ts）是否已包含 job/workflow 服务的路由转发

### 依赖影响

- 无新增前端依赖，复用已有 Ant Design Vue + VxeTable + Vben Form 组件

## Discovered Constraints

### Hard Constraints

- HC-1: 前端页面必须复用 Vben5 Admin 框架的 `useVbenVxeGrid` + `useVbenForm` + `useVbenModal` 模式
- HC-2: API 调用必须使用 `requestClient`（`#/api/request`），遵循统一响应格式 `R{code, data, msg}`
- HC-3: 表格分页参数使用 `pageNo` + `pageSize`，响应包含 `list` + `total`
- HC-4: 菜单通过后端 `sys_menu` 表动态生成，前端不硬编码路由
- HC-5: 所有 API 需携带 JWT Token（Authorization header），受 Auth 中间件保护
- HC-6: job-service 和 workflow-service 是独立微服务，前端需通过网关（Traefik）或代理访问

### Soft Constraints

- SC-1: 页面应支持 i18n（`$t()`），与现有 system 模块保持一致
- SC-2: Cron 表达式字段建议添加可视化编辑器（Cron UI 组件）或至少提供格式校验
- SC-3: 工作流执行结果支持 JSON 展示（JSON Tree 或代码高亮）
- SC-4: 手动触发任务应有确认弹窗，避免误操作

### Dependencies

- job-service REST API（已实现，端口通过 Traefik 代理）
- workflow-service REST API（已实现，端口通过 Traefik 代理）
- 前端 vite proxy 配置需更新，将 `/job/*` 和 `/workflow/*` 转发到对应服务
- sys_menu 表需插入菜单记录

### Risks

| Risk | Severity | Mitigation |
|------|----------|------------|
| 前端代理未配置 job/workflow 服务路由 | 高 | 更新 vite.config.ts proxy 规则 |
| job-service 与 system-service 响应格式不一致 | 中 | job-service DTO 已使用标准分页格式，需确认 response wrapper |
| 菜单数据未初始化 | 中 | 提供 SQL 脚本插入菜单记录 |
| 微服务跨域问题 | 低 | Traefik 网关统一入口 + CORS 中间件 |

## Success Criteria

- [ ] 定时任务管理页可完成任务的创建、编辑、删除、暂停/恢复、手动触发
- [ ] 执行日志页可查看历史执行记录并按条件筛选
- [ ] 工作流分类页可管理树形分类结构
- [ ] 工作流列表页可完成流程的创建、编辑、激活/停用、执行
- [ ] 执行实例页可查看工作流运行状态和结果
- [ ] 所有页面通过后端菜单系统动态加载
- [ ] API 调用成功对接后端微服务

## User Confirmations

1. 页面范围确认：定时任务（任务列表 + 执行日志）和工作流（分类 + 流程列表 + 实例）共 5 个页面
2. 前端微服务对接方式：通过 Vite dev proxy（开发）+ Traefik 网关（生产）
