import type { AppRouteRecordRaw } from '@vben-core/typings';

/**
 * 后端菜单树响应（Go 后端 MenuTreeResp）
 */
interface BackendMenuNode {
  id: number;
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
  children?: BackendMenuNode[];
}

/**
 * 后端用户响应（Go 后端 UserResp）
 */
interface BackendUserResp {
  id: number;
  username: string;
  nickname: string;
  email: string;
  mobile: string;
  sex: number;
  avatar: string;
  status: number;
  deptId: number;
  tenantId: number;
  createTime: string;
  remark: string;
}

/**
 * 后端权限信息响应（Go 后端 PermissionInfoResp）
 */
interface BackendPermissionInfo {
  user: BackendUserResp;
  roles: string[];
  permissions: string[];
  menus: BackendMenuNode[];
}

/**
 * 将后端菜单树转换为 Vben5 路由格式
 */
export function transformMenusToRoutes(
  menus: BackendMenuNode[],
): AppRouteRecordRaw[] {
  return menus
    .filter((m) => m.menuType !== 3) // 按钮类型不生成路由
    .map((m) => transformMenuNode(m))
    .sort((a, b) => (a.sort ?? 0) - (b.sort ?? 0));
}

function transformMenuNode(menu: BackendMenuNode): AppRouteRecordRaw {
  const route: AppRouteRecordRaw = {
    name: menu.menuName,
    path: menu.path,
    component: menu.component || '',
    meta: {} as RouteMeta,
    parentId: menu.parentId || undefined,
    sort: menu.sort,
    visible: menu.visible === 0,
    icon: menu.icon || undefined,
    id: menu.id,
  };

  if (menu.children && menu.children.length > 0) {
    route.children = transformMenusToRoutes(menu.children);
  }

  return route;
}

/**
 * 将后端用户信息转换为前端 UserInfo
 */
export function transformUser(user: BackendUserResp) {
  return {
    userId: String(user.id),
    username: user.username,
    nickname: user.nickname,
    avatar: user.avatar,
    homePath: '/dashboard',
  };
}

/**
 * 将后端权限信息转换为前端 AuthPermissionInfo
 */
export function transformPermissionInfo(data: BackendPermissionInfo) {
  return {
    user: transformUser(data.user),
    roles: data.roles,
    permissions: data.permissions,
    menus: transformMenusToRoutes(data.menus),
  };
}
