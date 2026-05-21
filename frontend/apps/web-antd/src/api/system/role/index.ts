import type { PageParam, PageResult } from '@vben/request';

import { requestClient } from '#/api/request';

export namespace SystemRoleApi {
  /** 角色信息 */
  export interface Role {
    id?: number;
    roleName: string;
    roleCode: string;
    sort: number;
    status: number;
    dataScope: number;
    dataScopeDeptIds?: string;
    remark?: string;
    createTime?: string;
  }
}

/** 角色分页列表 */
export function getRolePage(params: PageParam) {
  return requestClient.get<PageResult<SystemRoleApi.Role>>(
    '/system/role/page',
    { params },
  );
}

/** 角色详情 */
export function getRole(id: number) {
  return requestClient.get<SystemRoleApi.Role>(`/system/role/${id}`);
}

/** 新增角色 */
export function createRole(data: SystemRoleApi.Role) {
  return requestClient.post('/system/role', data);
}

/** 修改角色 */
export function updateRole(data: SystemRoleApi.Role) {
  return requestClient.put('/system/role', data);
}

/** 删除角色 */
export function deleteRole(id: number) {
  return requestClient.delete(`/system/role/${id}`);
}

/** 分配菜单权限 */
export function assignRoleMenus(roleId: number, menuIds: number[]) {
  return requestClient.put('/system/role/assign-menus', { roleId, menuIds });
}

/** 兼容旧引用 */
export function deleteRoleList(_ids: number[]) {
  return Promise.resolve();
}
export function exportRole(_params: any) {
  return Promise.resolve();
}
export function getSimpleRoleList() {
  return getRolePage({ page: 1, pageSize: 100 });
}
