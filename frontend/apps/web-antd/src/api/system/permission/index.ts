import { requestClient } from '#/api/request';

/** 获取角色已分配的菜单 ID 列表 */
export async function getRoleMenuList(roleId: number): Promise<number[]> {
  // 后端暂未提供独立接口，从角色详情中获取
  const role = await requestClient.get<any>(`/system/role/${roleId}`);
  return role?.menuIds ?? [];
}

/** 分配角色菜单权限 */
export async function assignRoleMenu(roleId: number, menuIds: number[]) {
  return requestClient.put('/system/role/assign-menus', { roleId, menuIds });
}

/** 获取用户已分配的角色 ID 列表 */
export async function getUserRoleList(userId: number): Promise<number[]> {
  // 后端暂未提供独立接口
  const user = await requestClient.get<any>(`/system/user/${userId}`);
  return user?.roleIds ?? [];
}

/** 分配用户角色 */
export async function assignUserRole(userId: number, roleIds: number[]) {
  return requestClient.put('/system/user', { id: userId, roleIds });
}

/** 分配角色数据权限 */
export async function assignRoleDataScope(
  roleId: number,
  dataScope: number,
  dataScopeDeptIds: string,
) {
  return requestClient.put('/system/role', {
    id: roleId,
    dataScope,
    dataScopeDeptIds,
  });
}
