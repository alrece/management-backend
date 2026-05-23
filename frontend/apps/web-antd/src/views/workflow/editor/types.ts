import type { WorkflowApi } from '#/api/workflow/workflow';

/** 工作流定义（前端结构） */
export interface WorkflowDefinition {
  nodes: WorkflowNode[];
  edges: WorkflowEdge[];
  settings: Record<string, any>;
}

/** 工作流节点 */
export interface WorkflowNode {
  id: string;
  type: string;
  position: { x: number; y: number };
  data: NodeData;
}

/** 节点数据 */
export interface NodeData {
  label: string;
  config: Record<string, any>;
  errorStrategy?: 'retry' | 'skip' | 'abort';
  retryCount?: number;
  timeout?: number;
}

/** 工作流边 */
export interface WorkflowEdge {
  id: string;
  source: string;
  target: string;
  sourceHandle?: string;
  targetHandle?: string;
  animated?: boolean;
  label?: string;
}

/** 节点分类 */
export type NodeCategory = 'trigger' | 'action' | 'control' | 'transform';

/** 分类标签映射 */
export const CATEGORY_LABELS: Record<NodeCategory, string> = {
  trigger: '触发器',
  action: '动作',
  control: '控制',
  transform: '转换',
};

/** 分类颜色映射 */
export const CATEGORY_COLORS: Record<NodeCategory, string> = {
  trigger: '#67C23A',
  action: '#409EFF',
  control: '#E6A23C',
  transform: '#F56C6C',
};

/** 状态枚举 */
export const WORKFLOW_STATUS = {
  ACTIVE: 0,
  INACTIVE: 1,
} as const;

export const WORKFLOW_STATUS_OPTIONS = [
  { label: '正常', value: 0 },
  { label: '停用', value: 1 },
] as const;
