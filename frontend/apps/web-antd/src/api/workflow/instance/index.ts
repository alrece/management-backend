import type { PageParam, PageResult } from '@vben/request';

import { requestClient } from '#/api/request';

export namespace InstanceApi {
  /** 执行实例 */
  export interface Instance {
    id: number;
    workflowId: number;
    workflowName?: string;
    status: string;
    triggerType: string;
    startTime: string;
    endTime?: string;
    input?: any;
    output?: any;
    errorMsg?: string;
  }

  /** 节点日志 */
  export interface NodeLog {
    id: number;
    executionId: number;
    nodeId: string;
    nodeType: string;
    status: string;
    input?: any;
    output?: any;
    errorMsg?: string;
    startTime: string;
    endTime?: string;
    retryCount: number;
  }
}

/** 查询执行实例分页 */
export function getInstancePage(params: PageParam & { workflowId?: number; status?: string }) {
  return requestClient.get<PageResult<InstanceApi.Instance>>(
    '/workflow/execution/page',
    { params: { page: params.pageNo, pageSize: params.pageSize, ...params } },
  );
}

/** 获取执行详情 */
export function getInstance(id: number) {
  return requestClient.get<InstanceApi.Instance>(`/workflow/execution/${id}`);
}

/** 获取节点执行日志 */
export function getInstanceLogs(id: number) {
  return requestClient.get<InstanceApi.NodeLog[]>(`/workflow/execution/${id}/logs`);
}
