import { requestClient } from '#/api/request';

export namespace SystemMenuApi {
  /** 菜单信息（匹配后端 MenuTreeResp） */
  export interface Menu {
    id?: number;
    parentId: number;
    menuName: string;
    menuType: number; // 1=目录 2=菜单 3=按钮
    path: string;
    component: string;
    permission: string;
    icon: string;
    sort: number;
    visible: number; // 0=显示 1=隐藏
    status: number;
  }
}

/** 菜单树 */
export function getMenuTree() {
  return requestClient.get<SystemMenuApi.Menu[]>('/system/menu/tree');
}

/** 新增菜单 */
export function createMenu(data: SystemMenuApi.Menu) {
  return requestClient.post('/system/menu', data);
}

/** 修改菜单 */
export function updateMenu(data: SystemMenuApi.Menu) {
  return requestClient.put('/system/menu', data);
}

/** 删除菜单 */
export function deleteMenu(id: number) {
  return requestClient.delete(`/system/menu/${id}`);
}

/** 兼容旧引用 */
export const getSimpleMenusList = getMenuTree;
export const getMenuList = getMenuTree;
export function getMenu(_id: number) {
  return Promise.resolve(null);
}
