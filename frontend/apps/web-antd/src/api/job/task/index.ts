import type { PageParam, PageResult } from '@vben/request';

import { requestClient } from '#/api/request';

export namespace JobTaskApi {
  /** 定时任务 */
  export interface Task {
    id?: number;
    name: string;
    handler: string;
    cronExpr: string;
    params: string;
    status: number;
    remark: string;
    creator?: number;
  }
}

/** 查询任务分页 */
export function getJobTaskPage(params: PageParam & { name?: string; status?: number }) {
  return requestClient.get<PageResult<JobTaskApi.Task>>(
    '/job/tasks/page',
    { params: { page: params.pageNo, pageSize: params.pageSize, ...params } },
  );
}

/** 查询任务详情 */
export function getJobTask(id: number) {
  return requestClient.get<JobTaskApi.Task>(`/job/tasks/${id}`);
}

/** 创建任务 */
export function createJobTask(data: JobTaskApi.Task) {
  return requestClient.post('/job/tasks', data);
}

/** 更新任务 */
export function updateJobTask(id: number, data: Partial<JobTaskApi.Task>) {
  return requestClient.put(`/job/tasks/${id}`, data);
}

/** 删除任务 */
export function deleteJobTask(id: number) {
  return requestClient.delete(`/job/tasks/${id}`);
}

/** 手动触发任务 */
export function triggerJobTask(id: number) {
  return requestClient.post(`/job/tasks/${id}/trigger`);
}
