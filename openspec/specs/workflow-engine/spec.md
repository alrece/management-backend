# Spec: 工作流引擎核心

## CONSTRAINT:ENG-001 — DAG 执行完整性
每个工作流执行必须为图中的每个节点生成一条 `wf_node_log` 记录（跳过的节点除外）。

## CONSTRAINT:ENG-002 — 循环检测
保存工作流时必须执行 DFS 环检测。检测到环则拒绝保存并返回明确错误信息，指出参与环的节点。

## CONSTRAINT:ENG-003 — 节点超时
每个节点执行必须通过 `context.WithTimeout` 设置超时。默认 60 秒，可在节点配置中覆盖。超时后节点状态标记为 FAILED。

## CONSTRAINT:ENG-004 — 错误策略
每个节点必须支持 3 种错误处理策略：
- `retry`（默认）：最多重试 3 次，指数退避（1s, 2s, 4s）
- `skip`：标记节点为 SKIPPED，继续执行下游
- `abort`：终止整个工作流执行

## CONSTRAINT:ENG-005 — 并发限制
DAG 执行器使用 goroutine pool，默认 worker 数量为 50。通过 `workerPool chan struct{}` 实现。

## CONSTRAINT:ENG-006 — 执行状态不可逆
`wf_execution.status` 状态转换只允许：RUNNING → SUCCESS / FAILED / CANCELED。已终态不可变更。

## CONSTRAINT:ENG-007 — JSONB 定义大小限制
`wf_workflow.definition` JSONB 最大 1MB。节点数量上限 100 个。

## CONSTRAINT:ENG-008 — 触发器幂等
同一个工作流的 Cron 触发器只能注册一次。重复激活调用不创建新的 cron entry。

## CONSTRAINT:ENG-009 — 租户隔离
所有工作流查询必须包含 `tenant_id` 过滤。Webhook 触发路径包含 `tenantId`。

## CONSTRAINT:ENG-010 — 日志保留
执行日志（`wf_execution` + `wf_node_log`）保留 30 天。通过 CronJob 定期清理，每次清理上限 1000 条。

## PROPERTY:PBT-001 — 往返一致性
对任意 WorkflowDefinition D：save(D) → load() → D'，则 deepEqual(D, D') = true。
伪造策略：随机生成节点和边，保存后加载比较。

## PROPERTY:PBT-002 — 拓扑排序正确性
对任意无环 DAG，拓扑排序结果满足：所有边的 source 节点在 target 节点之前。
伪造策略：随机生成 DAG，验证每条边的顺序。

## PROPERTY:PBT-003 — 执行日志完整性
对任意工作流执行，wf_node_log 记录数 = 图中节点数 - SKIPPED 节点数。
伪造策略：执行不同结构的工作流，比较日志数量。
