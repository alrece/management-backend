import type { PageParam, PageResult } from '@vben/request';

import { requestClient } from '#/api/request';

export namespace SystemUserApi {
  /** 用户信息 */
  export interface User {
    id?: number;
    username: string;
    nickname: string;
    deptId: number;
    email: string;
    mobile: string;
    sex: number;
    avatar: string;
    status: number;
    remark: string;
    createTime?: string;
  }
}

/** 用户分页列表 */
export function getUserPage(params: PageParam) {
  return requestClient.get<PageResult<SystemUserApi.User>>(
    '/system/user/page',
    { params },
  );
}

/** 用户详情 */
export function getUser(id: number) {
  return requestClient.get<SystemUserApi.User>(`/system/user/${id}`);
}

/** 新增用户 */
export function createUser(data: SystemUserApi.User) {
  return requestClient.post('/system/user', data);
}

/** 修改用户 */
export function updateUser(data: SystemUserApi.User) {
  return requestClient.put('/system/user', data);
}

/** 删除用户 */
export function deleteUser(id: number) {
  return requestClient.delete(`/system/user/${id}`);
}

/** 重置用户密码（原型：直接通过 update） */
export function resetUserPassword(id: number, password: string) {
  return requestClient.put('/system/user', { id, password });
}

/** 兼容旧引用 */
export function deleteUserList(_ids: number[]) {
  return Promise.resolve();
}
export function exportUser(_params: any) {
  return Promise.resolve();
}
export function updateUserStatus(id: number, status: number) {
  return requestClient.put('/system/user', { id, status });
}
export function importUser(_file: File, _updateSupport: boolean) {
  return Promise.resolve();
}
export function importUserTemplate() {
  return Promise.resolve();
}
export async function getSimpleUserList() {
  const res = await getUserPage({ page: 1, pageSize: 100 });
  return res.list ?? res ?? [];
}
