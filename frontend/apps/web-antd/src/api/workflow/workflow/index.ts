import type { PageParam, PageResult } from '@vben/request';

import { requestClient } from '#/api/request';

export namespace WorkflowApi {
  /** 工作流 */
  export interface Workflow {
    id?: number;
    name: string;
    categoryId: number;
    n8nWorkflowId?: string;
    paramsSchema: string;
    status: string;
    creator?: number;
    createdAt?: string;
    updater?: number;
    updatedAt?: string;
  }

  /** 工作流表单 */
  export interface WorkflowForm {
    name: string;
    categoryId: number;
    paramsSchema: string;
  }
}

/** 查询工作流分页 */
export function getWorkflowPage(params: PageParam & { name?: string; categoryId?: number; status?: string }) {
  return requestClient.get<PageResult<WorkflowApi.Workflow>>(
    '/workflow/workflows/page',
    { params: { page: params.pageNo, pageSize: params.pageSize, ...params } },
  );
}

/** 查询工作流详情 */
export function getWorkflow(id: number) {
  return requestClient.get<WorkflowApi.Workflow>(`/workflow/workflows/${id}`);
}

/** 创建工作流 */
export function createWorkflow(data: WorkflowApi.WorkflowForm) {
  return requestClient.post('/workflow/workflows', data);
}

/** 更新工作流 */
export function updateWorkflow(id: number, data: WorkflowApi.WorkflowForm) {
  return requestClient.put(`/workflow/workflows/${id}`, data);
}

/** 删除工作流 */
export function deleteWorkflow(id: number) {
  return requestClient.delete(`/workflow/workflows/${id}`);
}

/** 激活工作流 */
export function activateWorkflow(id: number) {
  return requestClient.put(`/workflow/workflows/${id}/activate`);
}

/** 停用工作流 */
export function deactivateWorkflow(id: number) {
  return requestClient.put(`/workflow/workflows/${id}/deactivate`);
}

/** 执行工作流 */
export function executeWorkflow(id: number, variables?: string) {
  return requestClient.post(`/workflow/workflows/${id}/execute`, { variables });
}
