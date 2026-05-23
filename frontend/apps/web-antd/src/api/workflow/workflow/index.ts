import type { PageParam, PageResult } from '@vben/request';

import { requestClient } from '#/api/request';

export namespace WorkflowApi {
  /** 工作流 */
  export interface Workflow {
    id?: number;
    name: string;
    categoryId?: number;
    definition?: string;
    triggerType: string;
    triggerConfig?: string;
    version?: number;
    status: number;
    remark?: string;
    tenantId?: number;
    creator?: number;
    createTime?: string;
    updateTime?: string;
  }

  /** 工作流表单 */
  export interface WorkflowForm {
    id?: number;
    name: string;
    categoryId?: number;
    definition?: string;
    triggerType: string;
    triggerConfig?: string;
    remark?: string;
  }

  /** 节点 Schema */
  export interface NodeSchema {
    type: string;
    label: string;
    category: string;
    icon: string;
    fields: SchemaField[];
  }

  export interface SchemaField {
    name: string;
    label: string;
    type: string;
    required: boolean;
    default?: any;
    options?: { value: string; label: string }[];
  }
}

/** 查询工作流分页 */
export function getWorkflowPage(params: PageParam & { name?: string; categoryId?: number; status?: number }) {
  return requestClient.get<PageResult<WorkflowApi.Workflow>>(
    '/workflow/workflow/page',
    { params: { page: params.pageNo, pageSize: params.pageSize, ...params } },
  );
}

/** 查询工作流详情 */
export function getWorkflow(id: number) {
  return requestClient.get<WorkflowApi.Workflow>(`/workflow/workflow/${id}`);
}

/** 创建工作流 */
export function createWorkflow(data: WorkflowApi.WorkflowForm) {
  return requestClient.post('/workflow/workflow', data);
}

/** 更新工作流 */
export function updateWorkflow(data: WorkflowApi.WorkflowForm) {
  return requestClient.put('/workflow/workflow', data);
}

/** 删除工作流 */
export function deleteWorkflow(id: number) {
  return requestClient.delete(`/workflow/workflow/${id}`);
}

/** 激活工作流 */
export function activateWorkflow(id: number) {
  return requestClient.put(`/workflow/workflow/${id}/activate`);
}

/** 停用工作流 */
export function deactivateWorkflow(id: number) {
  return requestClient.put(`/workflow/workflow/${id}/deactivate`);
}

/** 执行工作流 */
export function executeWorkflow(id: number, input?: any) {
  return requestClient.post(`/workflow/workflow/${id}/execute`, input);
}

/** 获取所有节点 Schema */
export function getNodeSchemas() {
  return requestClient.get<WorkflowApi.NodeSchema[]>('/workflow/workflow/nodes/schemas');
}
