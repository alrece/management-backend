import type { PageParam, PageResult } from '@vben/request';

import { requestClient } from '#/api/request';

export namespace InstanceApi {
  /** 执行实例 */
  export interface Instance {
    id: number;
    workflowId: number;
    n8nExecutionId: string;
    status: string;
    variables: string;
    result: string;
    durationMs: number;
    errorMsg: string;
    startedAt: string;
    finishedAt: string;
  }
}

/** 查询执行实例分页 */
export function getInstancePage(params: PageParam & { workflowId?: number; status?: string }) {
  return requestClient.get<PageResult<InstanceApi.Instance>>(
    '/workflow/instances/page',
    { params: { page: params.pageNo, pageSize: params.pageSize, ...params } },
  );
}
