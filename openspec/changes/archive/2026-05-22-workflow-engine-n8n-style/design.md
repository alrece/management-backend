# Design: 自研类 n8n 可视化工作流引擎

## 架构总览

```
┌─────────────────────────────────────────────────────────┐
│                      前端 (Vue 3)                        │
│                                                          │
│  ┌──────────┐  ┌──────────────────┐  ┌───────────────┐  │
│  │ 分类管理  │  │  工作流编辑器     │  │  执行监控      │  │
│  │ (树形保留) │  │  (Vue Flow)      │  │  (实例/日志)   │  │
│  └──────────┘  └────────┬─────────┘  └───────────────┘  │
│                          │ workflow JSON                  │
│  ┌───────────────────────┴─────────────────────────────┐ │
│  │  API Layer (requestClient)                           │ │
│  └───────────────────────┬─────────────────────────────┘ │
└──────────────────────────┼───────────────────────────────┘
                           │ HTTP (Gin)
┌──────────────────────────┼───────────────────────────────┐
│                Go 工作流引擎 (internal/module/workflow)    │
│                           │                               │
│  ┌────────────┐  ┌───────┴────────┐  ┌────────────────┐  │
│  │ 触发器管理  │  │  工作流 CRUD    │  │  执行引擎      │  │
│  │            │  │  (Handler/Svc/  │  │  (DAG Executor)│  │
│  │ Cron       │  │   Repo)        │  │                │  │
│  │ Webhook    │  │                │  │  拓扑排序      │  │
│  │ Manual     │  │                │  │  并行执行      │  │
│  └────────────┘  └────────────────┘  │  上下文传递    │  │
│                                       │  日志记录      │  │
│  ┌────────────────────────────────────┴───────────────┐  │
│  │  节点注册表 (NodeRegistry)                          │  │
│  │  ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────────┐ │  │
│  │  │HTTP  │ │SQL   │ │Cond  │ │Delay │ │Transform │ │  │
│  │  │Req   │ │Query │ │ition │ │      │ │          │ │  │
│  │  └──────┘ └──────┘ └──────┘ └──────┘ └──────────┘ │  │
│  └────────────────────────────────────────────────────┘  │
│                          │                                │
│  ┌──────────┐  ┌────────┴────────┐  ┌─────────────────┐ │
│  │ PostgreSQL│  │  Redis           │  │  pkg/redisx     │ │
│  │ (JSONB)  │  │  (锁/缓存/队列)  │  │  (分布式锁)     │ │
│  └──────────┘  └─────────────────┘  └─────────────────┘ │
└──────────────────────────────────────────────────────────┘
```

## 模块结构

```
internal/module/workflow/
├── model/
│   ├── entity.go          # Workflow, Execution, NodeLog, Trigger, Category 实体
│   └── dto.go             # 请求/响应 DTO
├── repository/
│   ├── workflow.go         # 工作流 CRUD
│   ├── execution.go        # 执行记录查询
│   ├── node_log.go         # 节点日志查询
│   ├── trigger.go          # 触发器配置
│   └── category.go         # 分类树
├── service/
│   ├── workflow.go         # 工作流业务逻辑
│   ├── executor.go         # DAG 执行引擎
│   ├── trigger.go          # 触发器管理
│   └── category.go         # 分类管理
├── handler/
│   ├── workflow.go         # 工作流 API
│   ├── execution.go        # 执行 API
│   ├── webhook.go          # Webhook 入口
│   └── category.go         # 分类 API
├── engine/
│   ├── registry.go         # 节点注册表
│   ├── dag.go              # DAG 拓扑排序
│   ├── context.go          # 节点执行上下文
│   └── nodes/              # 内置节点实现
│       ├── http_request.go
│       ├── sql_query.go
│       ├── condition.go
│       ├── delay.go
│       ├── transform.go
│       ├── log_message.go
│       └── send_email.go
└── router.go               # 路由注册
```

## 前端结构

```
frontend/apps/web-antd/src/
├── api/workflow/
│   ├── category/index.ts       # 分类 API
│   ├── workflow/index.ts       # 工作流 API (CRUD + 执行)
│   └── instance/index.ts       # 执行实例 API
├── views/workflow/
│   ├── category/
│   │   ├── index.vue           # 分类管理（树形保留）
│   │   ├── data.ts
│   │   └── modules/form.vue
│   ├── editor/
│   │   ├── index.vue           # 工作流编辑器主页面
│   │   ├── components/
│   │   │   ├── FlowCanvas.vue  # Vue Flow 画布
│   │   │   ├── NodePalette.vue # 节点面板（可拖拽）
│   │   │   ├── PropertyPanel.vue # 属性编辑面板
│   │   │   ├── EditorToolbar.vue # 工具栏
│   │   │   ├── nodes/          # 自定义节点组件
│   │   │   │   ├── TriggerNode.vue
│   │   │   │   ├── ActionNode.vue
│   │   │   │   ├── ControlNode.vue
│   │   │   │   └── TransformNode.vue
│   │   │   └── edges/          # 自定义边组件
│   │   │       └── DataFlowEdge.vue
│   │   ├── composables/
│   │   │   ├── useDragAndDrop.ts   # 拖拽逻辑
│   │   │   ├── useWorkflow.ts      # 工作流状态管理
│   │   │   ├── useHistory.ts       # 撤销/重做
│   │   │   └── useNodeSchema.ts    # 节点 Schema 获取
│   │   └── types.ts                # 编辑器类型定义
│   └── instance/
│       ├── index.vue           # 执行实例列表
│       └── data.ts
```

## 关键技术设计

### 1. DAG 执行引擎

```go
// executor.go 核心逻辑
type Executor struct {
    registry   *NodeRegistry
    db         *gorm.DB
    rdb        *redis.Client
    workerPool chan struct{} // 限制并发 goroutine 数量（默认 50）
}

func (e *Executor) Execute(ctx context.Context, wf *Workflow) (*Execution, error) {
    // 1. 解析 definition JSON → nodes + edges
    // 2. 拓扑排序（Kahn 算法）
    // 3. 按层级执行（同一层可并行）
    // 4. 每层：dispatch goroutines, WaitGroup 等待
    // 5. 收集每节点输出，传入下游节点
    // 6. 记录执行日志到 wf_node_log
}

// 拓扑排序
func topologicalSort(nodes []Node, edges []Edge) ([][]Node, error) {
    // 返回分层结果：[[layer0_nodes], [layer1_nodes], ...]
    // 同层节点无依赖，可并行执行
    // 检测到环则返回 error
}
```

### 2. 节点注册表

```go
// registry.go
type NodeHandler interface {
    Execute(ctx context.Context, input *NodeInput) (*NodeOutput, error)
    Type() string        // e.g. "http_request"
    Category() string    // "trigger" / "action" / "control" / "transform"
    Schema() NodeSchema  // 前端属性面板用
}

type NodeSchema struct {
    Type     string       `json:"type"`
    Label    string       `json:"label"`
    Category string       `json:"category"`
    Icon     string       `json:"icon"`
    Fields   []SchemaField `json:"fields"`
}

type SchemaField struct {
    Name     string `json:"name"`
    Label    string `json:"label"`
    Type     string `json:"type"` // string, number, boolean, select, json
    Required bool   `json:"required"`
    Default  any    `json:"default,omitempty"`
    Options  []SelectOption `json:"options,omitempty"`
}

type NodeRegistry struct {
    handlers map[string]NodeHandler
}

func (r *NodeRegistry) Register(h NodeHandler) {
    r.handlers[h.Type()] = h
}

func (r *NodeRegistry) Get(nodeType string) (NodeHandler, bool) {
    h, ok := r.handlers[nodeType]
    return h, ok
}

func (r *NodeRegistry) AllSchemas() []NodeSchema {
    // 返回所有已注册节点的 Schema，供前端节点面板使用
}
```

### 3. 节点执行上下文

```go
// context.go
type NodeInput struct {
    NodeID   string                 `json:"nodeId"`
    NodeType string                 `json:"nodeType"`
    Config   map[string]any         `json:"config"`   // 节点配置
    Items    []map[string]any       `json:"items"`    // 上游输出数据
}

type NodeOutput struct {
    Items []map[string]any `json:"items"` // 输出数据
    Error string           `json:"error,omitempty"`
}
```

### 4. 触发器管理

```go
// trigger.go
type TriggerManager struct {
    cron     *cron.Cron
    handlers map[int64]cron.EntryID // workflowID → entryID
    registry *NodeRegistry
    executor *Executor
}

func (m *TriggerManager) Activate(workflowID int64, triggerType string, config json.RawMessage) error {
    switch triggerType {
    case "cron":
        // 解析 cron 表达式，注册到 robfig/cron
    case "webhook":
        // 生成 webhook 路径，存入 wf_trigger
    case "manual":
        // 无需注册
    }
}

func (m *TriggerManager) Deactivate(workflowID int64) error {
    // 停止 cron entry
    // 删除 webhook 路由
}
```

### 5. 循环检测算法

```go
// dag.go
func detectCycle(nodes []string, edges []Edge) error {
    // DFS 颜色标记法
    // WHITE=未访问, GRAY=访问中, BLACK=已完成
    // 遇到 GRAY 节点 = 有环
}
```

### 6. 工作流 JSON 格式

```typescript
// 前端类型定义
interface WorkflowDefinition {
  nodes: WorkflowNode[];
  edges: WorkflowEdge[];
  settings: WorkflowSettings;
}

interface WorkflowNode {
  id: string;              // UUID
  type: string;            // 节点类型（http_request, condition, etc.）
  position: { x: number; y: number };
  data: {
    label: string;
    config: Record<string, any>;  // 节点配置参数
    errorStrategy?: 'retry' | 'skip' | 'abort';
    retryCount?: number;
    timeout?: number;             // 秒
  };
}

interface WorkflowEdge {
  id: string;
  source: string;          // 源节点 ID
  target: string;          // 目标节点 ID
  sourceHandle?: string;   // 输出口（如 condition 的 true/false）
  targetHandle?: string;
  animated?: boolean;
  label?: string;
}
```

### 7. 前端编辑器状态管理

```typescript
// useWorkflow.ts composable
export function useWorkflow(workflowId: Ref<string>) {
  const nodes = ref<WorkflowNode[]>([])
  const edges = ref<WorkflowEdge[]>([])
  const selectedNode = ref<WorkflowNode | null>(null)
  const saving = ref(false)

  // 加载工作流
  async function load() { ... }
  // 保存工作流（含循环检测）
  async function save() { ... }
  // 执行工作流
  async function execute() { ... }

  return { nodes, edges, selectedNode, saving, load, save, execute }
}

// useHistory.ts composable (Command Pattern)
export function useHistory() {
  const undoStack = ref<Command[]>([])
  const redoStack = ref<Command[]>([])

  function execute(command: Command) { ... }
  function undo() { ... }
  function redo() { ... }

  return { undo, redo, execute, canUndo, canRedo }
}
```

## API 设计

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/workflow/workflow | 创建工作流 |
| PUT | /api/workflow/workflow | 更新工作流 |
| DELETE | /api/workflow/workflow/:id | 删除工作流 |
| GET | /api/workflow/workflow/:id | 获取工作流详情（含 definition） |
| GET | /api/workflow/workflow/page | 工作流分页列表 |
| PUT | /api/workflow/workflow/:id/activate | 激活工作流 |
| PUT | /api/workflow/workflow/:id/deactivate | 停用工作流 |
| POST | /api/workflow/workflow/:id/execute | 手动执行 |
| GET | /api/workflow/workflow/nodes/schemas | 获取所有节点 Schema |
| GET | /api/workflow/execution/page | 执行实例分页 |
| GET | /api/workflow/execution/:id | 执行详情 |
| GET | /api/workflow/execution/:id/logs | 节点执行日志 |
| POST | /api/wf/webhook/:tenantId/:workflowId/*path | Webhook 触发入口 |
| POST | /api/workflow/category | 创建分类 |
| PUT | /api/workflow/category | 更新分类 |
| DELETE | /api/workflow/category/:id | 删除分类 |
| GET | /api/workflow/category/tree | 分类树 |

## PostgreSQL 迁移设计

### 配置变更

```yaml
# config.yaml
postgres:                          # 替换 mysql:
  host: 127.0.0.1
  port: 5432
  username: mb_user
  password: ""
  database: management_backend
  sslmode: disable
  max-idle-conns: 10
  max-open-conns: 100
  log-level: warn
```

### GORM Driver 替换

```go
// go.mod
- gorm.io/driver/mysql v1.6.0
+ gorm.io/driver/postgres v1.6.0    (需新增)

// 数据库连接
import "gorm.io/driver/postgres"
dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
    cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database, cfg.SSLMode)
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{...})
```

### 类型映射（MySQL → PostgreSQL）

| MySQL | PostgreSQL | 说明 |
|-------|-----------|------|
| `BIGINT AUTO_INCREMENT` | `BIGSERIAL` 或 `BIGINT` | 主键仍用雪花 ID |
| `DATETIME` | `TIMESTAMPTZ` | 时区感知 |
| `BIT(1)` | `SMALLINT` 或 `BOOLEAN` | deleted 字段 |
| `TINYINT` | `SMALLINT` | status 字段 |
| `VARCHAR(n)` | `VARCHAR(n)` | 不变 |
| `TEXT` | `TEXT` | 不变 |
| N/A | `JSONB` | 新增：工作流定义 |
