## 1. API 层实现

- [x] 1.1 创建 api/job/task/index.ts：JobTaskApi 接口定义 + getJobTaskPage / getJobTask / createJobTask / updateJobTask / deleteJobTask / triggerJobTask 函数
- [x] 1.2 创建 api/job/log/index.ts：ExecLogApi 接口定义 + getExecLogPage 函数
- [x] 1.3 创建 api/workflow/category/index.ts：CategoryApi 接口定义 + getCategoryTree / createCategory / updateCategory / deleteCategory 函数
- [x] 1.4 创建 api/workflow/workflow/index.ts：WorkflowApi 接口定义 + getWorkflowPage / getWorkflow / createWorkflow / updateWorkflow / deleteWorkflow / activateWorkflow / deactivateWorkflow / executeWorkflow 函数
- [x] 1.5 创建 api/workflow/instance/index.ts：InstanceApi 接口定义 + getInstancePage 函数

## 2. 定时任务管理页面

- [x] 2.1 创建 views/job/task/data.ts：表格列定义（名称、处理器、Cron表达式、状态、参数、备注）+ 表单 Schema + TASK_STATUS 枚举
- [x] 2.2 创建 views/job/task/index.vue：任务列表页（useVbenVxeGrid 分页 + 工具栏按钮：新建/删除批量 + 行操作：编辑/删除/暂停/触发）
- [x] 2.3 创建 views/job/task/modules/form.vue：任务创建/编辑表单 Modal（useVbenForm + useVbenModal，字段：name/handler/cronExpr/params/status/remark）

## 3. 执行日志页面

- [x] 3.1 创建 views/job/log/data.ts：表格列定义（任务名称、触发类型、状态、耗时、结果、错误、开始时间）+ EXEC_STATUS / TRIGGER_TYPE 枚举
- [x] 3.2 创建 views/job/log/index.vue：执行日志列表页（useVbenVxeGrid 分页 + 搜索：任务ID/状态/时间范围，只读无编辑操作）

## 4. 工作流分类页面

- [x] 4.1 创建 views/workflow/category/data.ts：树列定义 + 表单 Schema
- [x] 4.2 创建 views/workflow/category/index.vue：分类管理页（a-tree 展示 + 新建/编辑/删除操作）
- [x] 4.3 创建 views/workflow/category/modules/form.vue：分类创建/编辑表单 Modal（字段：name/parentId/sort）

## 5. 工作流管理页面

- [x] 5.1 创建 views/workflow/workflow/data.ts：表格列定义（名称、分类、状态、创建时间）+ 表单 Schema + WORKFLOW_STATUS 枚举
- [x] 5.2 创建 views/workflow/workflow/index.vue：工作流列表页（useVbenVxeGrid 分页 + 工具栏：新建 + 行操作：编辑/删除/激活/停用/执行）
- [x] 5.3 创建 views/workflow/workflow/modules/form.vue：工作流创建/编辑表单 Modal（字段：name/categoryId/paramsSchema）

## 6. 执行实例页面

- [x] 6.1 创建 views/workflow/instance/data.ts：表格列定义（工作流名称、状态、耗时、变量、结果、错误、开始/结束时间）+ INSTANCE_STATUS 枚举
- [x] 6.2 创建 views/workflow/instance/index.vue：执行实例列表页（useVbenVxeGrid 分页 + 搜索：工作流ID/状态，只读无编辑操作）

## 7. 基础设施更新

- [x] 7.1 更新 vite.config.ts Mock 数据：在 loadMockData() 中添加 jobTasks/jobExecLogs/wfCategories/wfWorkflows/wfInstances 模拟数据
- [x] 7.2 更新 vite.config.ts Mock 路由：在 handlers 中添加 job 和 workflow 相关的 Mock API 路由匹配规则
- [x] 7.3 更新 vite.config.ts 菜单数据：在 mockData.menus 中添加定时任务和工作流管理的菜单树节点
- [x] 7.4 更新 vite.config.ts 权限数据：在 get-permission-info handler 的 permissions 中添加 job:*/workflow:* 权限标识
- [x] 7.5 更新 vite.config.ts 代理配置：添加 /api/job → :8082 和 /api/workflow → :8083 的代理规则
- [x] 7.6 创建 scripts/sql/menu-job-workflow.sql：sys_menu 表插入定时任务和工作流管理菜单记录的 SQL 脚本（供生产环境使用）

## 8. 验证

- [ ] 8.1 验证前端 Mock 模式下所有 5 个页面可正常访问和操作
- [ ] 8.2 验证菜单树正确显示定时任务和工作流管理的菜单节点
- [ ] 8.3 验证 vite build 无编译错误
