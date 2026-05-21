/**
 * Mock API 处理器
 * 匹配请求路径并返回模拟数据
 */
import type { IncomingMessage, ServerResponse } from 'node:http';

import { mockDepts, mockMenus, mockPosts, mockRoles, mockUsers, nextId } from './data';

type MockHandler = (req: IncomingMessage, body: any, query: Record<string, string>, params: Record<string, string>) => any;

const ok = (data: any) => ({ code: 0, data, msg: 'success' });
const fail = (code: number, msg: string) => ({ code, data: null, msg });

// 解析请求体
function parseBody(req: IncomingMessage): Promise<any> {
  return new Promise((resolve) => {
    let data = '';
    req.on('data', (chunk) => (data += chunk));
    req.on('end', () => {
      try { resolve(data ? JSON.parse(data) : {}); }
      catch { resolve({}); }
    });
  });
}

// 解析 query 参数
function parseQuery(url: string): Record<string, string> {
  const qs = url.split('?')[1] || '';
  const params: Record<string, string> = {};
  for (const pair of qs.split('&')) {
    const [k, v] = pair.split('=');
    if (k) params[decodeURIComponent(k)] = decodeURIComponent(v || '');
  }
  return params;
}

// ========== 路由表 ==========
const routes: Record<string, Record<string, MockHandler>> = {
  // 认证
  'POST /api/auth/login': (_req, body) => {
    const user = mockUsers.find((u) => u.username === body.username);
    if (!user) return fail(401, '用户名或密码错误');
    return ok({
      accessToken: `mock_access_${Date.now()}`,
      refreshToken: `mock_refresh_${Date.now()}`,
    });
  },
  'POST /api/auth/refresh': () => ok({
    accessToken: `mock_access_${Date.now()}`,
    refreshToken: `mock_refresh_${Date.now()}`,
  }),
  'POST /api/auth/logout': () => ok(null),

  // 权限信息
  'GET /api/system/auth/get-permission-info': () => ok({
    user: { userId: '1', username: 'admin', nickname: '管理员', avatar: '', homePath: '/dashboard' },
    roles: ['super_admin'],
    permissions: [
      'system:user:create', 'system:user:update', 'system:user:delete', 'system:user:query',
      'system:user:update-password', 'system:user:export', 'system:user:import',
      'system:role:create', 'system:role:update', 'system:role:delete', 'system:role:query', 'system:role:export',
      'system:menu:create', 'system:menu:update', 'system:menu:delete', 'system:menu:query',
      'system:dept:create', 'system:dept:update', 'system:dept:delete', 'system:dept:query',
      'system:post:create', 'system:post:update', 'system:post:delete', 'system:post:query', 'system:post:export',
      'system:permission:assign-user-role', 'system:permission:assign-role-data-scope',
      'system:permission:assign-role-menu',
      'system:dict:create', 'system:dict:update', 'system:dict:delete', 'system:dict:query', 'system:dict:export',
      'system:notice:create', 'system:notice:update', 'system:notice:delete', 'system:notice:query',
      'system:tenant:create', 'system:tenant:update', 'system:tenant:delete', 'system:tenant:query',
      'system:tenant:update-package',
      'system:login-log:query', 'system:login-log:export',
      'system:operate-log:query', 'system:operate-log:export',
    ],
    menus: mockMenus,
  }),

  // 用户管理
  'GET /api/system/user/page': (_req, _body, query) => {
    const page = Number(query.page) || 1;
    const pageSize = Number(query.pageSize) || 10;
    const start = (page - 1) * pageSize;
    return ok({ list: mockUsers.slice(start, start + pageSize), total: mockUsers.length, page, pageSize });
  },
  'GET /api/system/user/profile': () => ok(mockUsers[0]),
  'POST /api/system/user': (_req, body) => {
    const user = { ...body, id: nextId(), createTime: new Date().toISOString() };
    mockUsers.push(user);
    return ok(user.id);
  },
  'PUT /api/system/user': (_req, body) => {
    const idx = mockUsers.findIndex((u) => u.id === body.id);
    if (idx >= 0) Object.assign(mockUsers[idx], body);
    return ok(null);
  },
  'DELETE /api/system/user': () => ok(null),
  'GET /api/system/user/:id': (_req, _body, _query, params) => {
    const id = Number(params?.id);
    const user = mockUsers.find((u) => u.id === id) || mockUsers[0];
    return ok({ ...user, roleIds: [1] });
  },

  // 角色管理
  'GET /api/system/role/page': (_req, _body, query) => {
    const page = Number(query.page) || 1;
    const pageSize = Number(query.pageSize) || 10;
    return ok({ list: mockRoles, total: mockRoles.length, page, pageSize });
  },
  'GET /api/system/role': () => ok(mockRoles[0]),
  'GET /api/system/role/:id': (_req, _body, query, params) => {
    const id = Number(params?.id);
    const role = mockRoles.find((r) => r.id === id) || mockRoles[0];
    return ok({ ...role, menuIds: [10, 11, 12, 13, 14, 15, 16] });
  },
  'POST /api/system/role': (_req, body) => {
    const role = { ...body, id: nextId(), createTime: new Date().toISOString() };
    mockRoles.push(role);
    return ok(role.id);
  },
  'PUT /api/system/role': (_req, body) => {
    const idx = mockRoles.findIndex((r) => r.id === body.id);
    if (idx >= 0) Object.assign(mockRoles[idx], body);
    return ok(null);
  },
  'PUT /api/system/role/assign-menus': () => ok(null),
  'DELETE /api/system/role': () => ok(null),

  // 菜单管理
  'GET /api/system/menu/tree': () => ok(mockMenus),
  'POST /api/system/menu': (_req, body) => ok({ ...body, id: nextId() }),
  'PUT /api/system/menu': () => ok(null),
  'DELETE /api/system/menu': () => ok(null),

  // 部门管理
  'GET /api/system/dept/tree': () => ok(mockDepts),
  'POST /api/system/dept': (_req, body) => ok({ ...body, id: nextId() }),
  'PUT /api/system/dept': () => ok(null),
  'DELETE /api/system/dept': () => ok(null),

  // 岗位管理
  'GET /api/system/post/page': (_req, _body, query) => {
    const page = Number(query.page) || 1;
    const pageSize = Number(query.pageSize) || 10;
    return ok({ list: mockPosts, total: mockPosts.length, page, pageSize });
  },
  'POST /api/system/post': (_req, body) => {
    const post = { ...body, id: nextId(), createTime: new Date().toISOString() };
    mockPosts.push(post);
    return ok(post.id);
  },
  'PUT /api/system/post': () => ok(null),
  'DELETE /api/system/post': () => ok(null),

  // 岗位精简列表
  'GET /api/system/post/simple-list': () => ok(mockPosts),

  // 字典数据精简列表
  'GET /api/system/dict-data/simple-list': () => ok([]),
};

/**
 * 处理 mock 请求
 */
export async function handleMockRequest(
  req: IncomingMessage,
  res: ServerResponse,
): Promise<boolean> {
  const url = req.url || '/';
  const pathWithoutQuery = url.split('?')[0];
  const query = parseQuery(url);
  const method = (req.method || 'GET').toUpperCase();

  // 尝试精确匹配
  const exactKey = `${method} ${pathWithoutQuery}`;
  let handler = routes[exactKey];
  let routeParams: Record<string, string> = {};

  // 尝试带参数匹配（如 /api/system/user/123）
  if (!handler) {
    for (const key of Object.keys(routes)) {
      const [kMethod, kPath] = key.split(' ');
      if (kMethod !== method) continue;
      const kParts = kPath.split('/');
      const uParts = pathWithoutQuery.split('/');
      if (kParts.length !== uParts.length) continue;
      const params: Record<string, string> = {};
      const match = kParts.every((p, i) => {
        if (p.startsWith(':')) { params[p.slice(1)] = uParts[i]; return true; }
        return p === uParts[i];
      });
      if (match) {
        handler = routes[key];
        routeParams = params;
        break;
      }
    }
  }

  if (!handler) return false;

  const body = method === 'POST' || method === 'PUT' ? await parseBody(req) : {};
  const result = handler(req, body, query, routeParams);

  res.setHeader('Content-Type', 'application/json');
  res.setHeader('Access-Control-Allow-Origin', '*');
  res.end(JSON.stringify(result));
  return true;
}
