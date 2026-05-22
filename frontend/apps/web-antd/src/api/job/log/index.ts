import type { PageParam, PageResult } from '@vben/request';

import { requestClient } from '#/api/request';

export namespace ExecLogApi {
  /** 执行日志 */
  export interface Log {
    id: number;
    taskId: number;
    taskName: string;
    triggerType: number;
    status: number;
    durationMs: number;
    result: string;
    error: string;
    startTime: string;
  }
}

/** 查询执行日志分页 */
export function getExecLogPage(params: PageParam & { taskId?: number; status?: number; startTime?: string; endTime?: string }) {
  return requestClient.get<PageResult<ExecLogApi.Log>>(
    '/job/execution-logs/page',
    { params: { page: params.pageNo, pageSize: params.pageSize, ...params } },
  );
}
