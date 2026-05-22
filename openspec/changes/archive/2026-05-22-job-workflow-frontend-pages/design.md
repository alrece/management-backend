## Technical Design

### 架构

```
frontend/apps/web-antd/src/
├── api/
│   ├── job/
│   │   ├── task/index.ts        # JobTaskApi 接口 + CRUD + Trigger API
│   │   └── log/index.ts         # ExecLogApi 接口 + 分页查询 API
│   └── workflow/
│       ├── category/index.ts    # CategoryApi 接口 + CRUD + Tree API
│       ├── workflow/index.ts    # WorkflowApi 接口 + CRUD + Activate/Execute API
│       └── instance/index.ts    # InstanceApi 接口 + 分页查询 API
├── views/
│   ├── job/
│   │   ├── task/
│   │   │   ├── data.ts          # 表格列 + 表单 Schema + 状态枚举
│   │   │   ├── index.vue        # 任务列表页（VxeGrid + Modal）
│   │   │   └── modules/form.vue # 任务创建/编辑表单 Modal
│   │   └── log/
│   │       ├── data.ts          # 表格列 + 表单 Schema + 状态/触发类型枚举
│   │       └── index.vue        # 执行日志列表页（VxeGrid，只读）
│   └── workflow/
│       ├── category/
│       │   ├── data.ts          # 树列定义 + 表单 Schema
│       │   ├── index.vue        # 分类管理页（Tree + Modal）
│       │   └── modules/form.vue # 分类创建/编辑表单 Modal
│       ├── workflow/
│       │   ├── data.ts          # 表格列 + 表单 Schema + 状态枚举
│       │   ├── index.vue        # 工作流列表页（VxeGrid + Modal）
│       │   └── modules/form.vue # 工作流创建/编辑表单 Modal
│       └── instance/
│           ├── data.ts          # 表格列 + 表单 Schema + 状态枚举
│           └── index.vue        # 执行实例列表页（VxeGrid，只读）
```

### API 层设计

#### 分页参数映射

后端 DTO 使用 `page` + `pageSize`，前端统一使用 `PageParam` 类型：

```typescript
// api/job/task/index.ts
export function getJobTaskPage(params: PageParam) {
  return requestClient.get<PageResult<JobTaskApi.Task>>('/job/tasks/page', {
    params: { page: params.pageNo, pageSize: params.pageSize, ...params },
  });
}
```

#### 状态枚举

```typescript
// 任务状态
export const TASK_STATUS = [
  { label: '正常', value: 0 },
  { label: '暂停', value: 1 },
] as const;

// 执行状态
export const EXEC_STATUS = [
  { label: '执行中', value: 0 },
  { label: '成功', value: 1 },
  { label: '失败', value: 2 },
  { label: '超时', value: 3 },
] as const;

// 触发类型
export const TRIGGER_TYPE = [
  { label: '定时', value: 1 },
  { label: '手动', value: 2 },
] as const;

// 工作流状态
export const WORKFLOW_STATUS = [
  { label: '草稿', value: 'DRAFT' },
  { label: '已激活', value: 'ACTIVE' },
  { label: '已停用', value: 'INACTIVE' },
  { label: '已失同步', value: 'DESYNCED' },
] as const;

// 实例状态
export const INSTANCE_STATUS = [
  { label: '运行中', value: 'RUNNING' },
  { label: '成功', value: 'SUCCESS' },
  { label: '失败', value: 'FAILED' },
  { label: '已取消', value: 'CANCELLED' },
] as const;
```

### Vite Proxy 设计

```typescript
// vite.config.ts server.proxy 更新
proxy: {
  '/api/system': {
    target: 'http://127.0.0.1:8081',
    changeOrigin: true,
  },
  '/api/job': {
    target: 'http://127.0.0.1:8082',
    changeOrigin: true,
  },
  '/api/workflow': {
    target: 'http://127.0.0.1:8083',
    changeOrigin: true,
  },
}
```

### Mock 数据设计

在 `vite.config.ts` 的 `loadMockData()` 中添加：

```typescript
jobTasks: [
  { id: 1, name: '数据备份', handler: 'DataBackupHandler', cronExpr: '0 0 2 * * ?', params: '{}', status: 0, remark: '每日凌晨2点', creator: 1 },
],
jobExecLogs: [
  { id: 1, taskId: 1, taskName: '数据备份', triggerType: 1, status: 1, durationMs: 3200, result: '备份完成', error: '', startTime: '2026-05-22 02:00:00' },
],
wfCategories: [
  { id: 1, name: '审批流程', parentId: 0, sort: 0 },
  { id: 2, name: '数据同步', parentId: 0, sort: 1 },
],
wfWorkflows: [
  { id: 1, name: '请假审批', categoryId: 1, n8nWorkflowId: '', paramsSchema: '{}', status: 'ACTIVE', creator: 1, createdAt: '2026-05-01', updater: 1, updatedAt: '2026-05-01' },
],
wfInstances: [
  { id: 1, workflowId: 1, n8nExecutionId: 'exec_001', status: 'SUCCESS', variables: '{}', result: '{"approved": true}', durationMs: 1500, errorMsg: '', startedAt: '2026-05-22 10:00:00', finishedAt: '2026-05-22 10:00:01' },
],
```

Mock handlers 添加对应的路由匹配规则。

### 菜单数据设计

在 Mock 的 `menus` 数组中添加：

```typescript
{
  id: 100, parentId: 0, menuName: '定时任务', menuType: 1,
  path: '/job', component: '', permission: '',
  icon: 'ant-design:clock-circle-outlined', sort: 20, visible: 0, status: 0,
  children: [
    { id: 101, parentId: 100, menuName: '任务管理', menuType: 2, path: 'task', component: '/job/task/index', permission: 'job:task:query', icon: '', sort: 1, visible: 0, status: 0, children: [] },
    { id: 102, parentId: 100, menuName: '执行日志', menuType: 2, path: 'log', component: '/job/log/index', permission: 'job:log:query', icon: '', sort: 2, visible: 0, status: 0, children: [] },
  ],
},
{
  id: 200, parentId: 0, menuName: '工作流管理', menuType: 1,
  path: '/workflow', component: '', permission: '',
  icon: 'ant-design:branches-outlined', sort: 30, visible: 0, status: 0,
  children: [
    { id: 201, parentId: 200, menuName: '流程分类', menuType: 2, path: 'category', component: '/workflow/category/index', permission: 'workflow:category:query', icon: '', sort: 1, visible: 0, status: 0, children: [] },
    { id: 202, parentId: 200, menuName: '流程管理', menuType: 2, path: 'workflow', component: '/workflow/workflow/index', permission: 'workflow:workflow:query', icon: '', sort: 2, visible: 0, status: 0, children: [] },
    { id: 203, parentId: 200, menuName: '执行实例', menuType: 2, path: 'instance', component: '/workflow/instance/index', permission: 'workflow:instance:query', icon: '', sort: 3, visible: 0, status: 0, children: [] },
  ],
},
```

### 权限标识设计

| 权限 | 标识 |
|------|------|
| 任务查询 | `job:task:query` |
| 任务创建 | `job:task:create` |
| 任务编辑 | `job:task:update` |
| 任务删除 | `job:task:delete` |
| 手动触发 | `job:task:trigger` |
| 日志查询 | `job:log:query` |
| 分类查询 | `workflow:category:query` |
| 分类创建 | `workflow:category:create` |
| 分类编辑 | `workflow:category:update` |
| 分类删除 | `workflow:category:delete` |
| 工作流查询 | `workflow:workflow:query` |
| 工作流创建 | `workflow:workflow:create` |
| 工作流编辑 | `workflow:workflow:update` |
| 工作流删除 | `workflow:workflow:delete` |
| 工作流激活 | `workflow:workflow:activate` |
| 工作流执行 | `workflow:workflow:execute` |
| 实例查询 | `workflow:instance:query` |

### 组件复用策略

1. **状态渲染**：通用 `StatusTag` 组件，接受 `(value, options[])` 返回带颜色的 `<a-tag>`
2. **Cron 输入**：使用 `<a-input>` + 格式提示文本（后续可替换为可视化编辑器）
3. **JSON 展示**：使用 `<a-typography-paragraph>` 的 `code` 模式展示结果
4. **分类树**：复用 `<a-tree>` + modal 表单模式（参考菜单管理页）
