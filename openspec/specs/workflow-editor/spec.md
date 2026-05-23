# Spec: 前端编辑器

## CONSTRAINT:FE-001 — Vue Flow 版本
使用 `@vue-flow/core` 作为画布引擎。配套 `@vue-flow/background`、`@vue-flow/minimap`、`@vue-flow/controls`。

## CONSTRAINT:FE-002 — 节点分类
前端节点面板按 Category 分组显示：
- 触发器（trigger）：Cron、Webhook、手动
- 动作（action）：HTTP 请求、SQL 查询、发送邮件、日志输出
- 控制（control）：条件判断、延迟
- 转换（transform）：数据转换、数据过滤

## CONSTRAINT:FE-003 — 自定义节点样式
每种 Category 对应不同的节点视觉样式：
- trigger：绿色边框，左侧入口不可连接
- action：蓝色边框
- control：橙色边框，多输出口（如 condition 的 true/false）
- transform：紫色边框

## CONSTRAINT:FE-004 — 属性面板 Schema 驱动
属性面板根据后端返回的 NodeSchema 动态渲染表单。Schema 字段类型映射：
- string → Input / Textarea
- number → InputNumber
- boolean → Switch
- select → Select
- json → Code Editor (monaco-editor)

## CONSTRAINT:FE-005 — 撤销/重做
使用 Command Pattern 实现。最多保留 50 步历史。操作类型：
- 添加节点、删除节点、移动节点
- 添加连线、删除连线
- 修改节点配置

## CONSTRAINT:FE-006 — 自动布局
提供自动布局按钮，使用 dagre 算法（`@dagrejs/dagre`）。布局方向 TB（从上到下）。

## CONSTRAINT:FE-007 — 保存校验
保存工作流前执行前端校验：
1. 至少有一个触发器节点
2. 所有非触发器节点有至少一条入边
3. 所有节点有 label
4. 无循环依赖（后端二次校验）

## CONSTRAINT:FE-008 — 执行可视化
工作流执行时：
- 已完成节点：绿色边框
- 运行中节点：脉冲动画
- 失败节点：红色边框 + 错误 tooltip
- 待执行节点：灰色
- 活跃边：流动动画（animated: true）

## CONSTRAINT:FE-009 — 编辑器 URL
编辑器路由：`/workflow/editor/:id`，id 为工作流 ID。从工作流列表点击"编辑"进入。
