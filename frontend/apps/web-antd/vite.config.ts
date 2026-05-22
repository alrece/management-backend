import { defineConfig } from '@vben/vite-config';

export default defineConfig(async () => {
  return {
    application: {},
    vite: {
      plugins: [
        {
          name: 'vite:mock-api',
          enforce: 'pre',
          apply: 'serve',
          configureServer(server) {
            const fs = require('fs');
            const path = require('path');

            // 读取 mock 开关（从 .env.development 读取 VITE_USE_MOCK）
            let useMock = true;
            try {
              const envContent = fs.readFileSync(path.resolve(__dirname, '.env.development'), 'utf-8');
              const m = envContent.match(/VITE_USE_MOCK=(\w+)/);
              if (m) useMock = m[1] === 'true';
            } catch {}

            let mockData: any = null;
            function loadMockData() {
              const dataPath = path.resolve(__dirname, 'src/mock/data.ts');
              // 使用 Vite 的 fs 读取（模拟）
              return {
                users: [
                  { id: 1, username: 'admin', nickname: '管理员', email: 'admin@mb.com', mobile: '13800000000', sex: 1, avatar: '', status: 0, deptId: 100, tenantId: 1, createTime: '2026-01-01 00:00:00', remark: '' },
                ],
                roles: [
                  { id: 1, name: '超级管理员', code: 'super_admin', sort: 1, dataScope: 1, status: 0, remark: '', createTime: '2026-01-01 00:00:00' },
                ],
                menus: [
                  { id: 1, parentId: 0, menuName: '系统管理', menuType: 1, path: '/system', component: '', permission: '', icon: 'ant-design:setting-outlined', sort: 10, visible: 0, status: 0, children: [
                    { id: 10, parentId: 1, menuName: '用户管理', menuType: 2, path: 'user', component: '/system/user/index', permission: 'system:user:query', icon: 'ant-design:user-outlined', sort: 1, visible: 0, status: 0, children: [] },
                    { id: 11, parentId: 1, menuName: '角色管理', menuType: 2, path: 'role', component: '/system/role/index', permission: 'system:role:query', icon: 'ant-design:team-outlined', sort: 2, visible: 0, status: 0, children: [] },
                    { id: 12, parentId: 1, menuName: '菜单管理', menuType: 2, path: 'menu', component: '/system/menu/index', permission: 'system:menu:query', icon: 'ant-design:menu-outlined', sort: 3, visible: 0, status: 0, children: [] },
                    { id: 13, parentId: 1, menuName: '部门管理', menuType: 2, path: 'dept', component: '/system/dept/index', permission: 'system:dept:query', icon: 'ant-design:apartment-outlined', sort: 4, visible: 0, status: 0, children: [] },
                  ] },
                ],
                depts: [
                  { id: 100, parentId: 0, name: '总公司', sort: 0, leaderUserId: 1, status: 0, createTime: '2026-01-01 00:00:00' },
                  { id: 101, parentId: 100, name: '技术部', sort: 1, leaderUserId: 2, status: 0, createTime: '2026-01-01 00:00:00' },
                ],
                posts: [
                  { id: 1, postCode: 'dev', postName: '开发工程师', sort: 1, status: 0, remark: '', createTime: '2026-01-01 00:00:00' },
                ],
              };
            }

            const ok = (data: any) => JSON.stringify({ code: 0, data, msg: 'success' });

            // 路由匹配器
            const handlers: Array<{ match: (method: string, path: string) => boolean; handle: (body: any, query: any) => any }> = [
              // 登录
              { match: (m, p) => m === 'POST' && p === '/api/auth/login', handle: () => ({ accessToken: `mock_access_${Date.now()}`, refreshToken: `mock_refresh_${Date.now()}` }) },
              { match: (m, p) => m === 'POST' && p === '/api/auth/refresh', handle: () => ({ accessToken: `mock_access_${Date.now()}`, refreshToken: `mock_refresh_${Date.now()}` }) },
              { match: (m, p) => m === 'POST' && p === '/api/auth/logout', handle: () => null },

              // 权限信息
              { match: (m, p) => m === 'GET' && p === '/api/system/auth/get-permission-info', handle: (_, __) => {
                const d = loadMockData();
                return {
                  user: { userId: '1', username: 'admin', nickname: '管理员', avatar: '', homePath: '/dashboard' },
                  roles: ['super_admin'],
                  permissions: [
                    'system:user:create', 'system:user:update', 'system:user:delete', 'system:user:query',
                    'system:user:update-password', 'system:user:export', 'system:user:import',
                    'system:role:create', 'system:role:update', 'system:role:delete', 'system:role:query', 'system:role:export',
                    'system:menu:create', 'system:menu:update', 'system:menu:delete', 'system:menu:query',
                    'system:dept:create', 'system:dept:update', 'system:dept:delete', 'system:dept:query',
                    'system:post:create', 'system:post:update', 'system:post:delete', 'system:post:query', 'system:post:export',
                    'system:permission:assign-user-role', 'system:permission:assign-role-data-scope', 'system:permission:assign-role-menu',
                    'system:dict:create', 'system:dict:update', 'system:dict:delete', 'system:dict:query', 'system:dict:export',
                    'system:notice:create', 'system:notice:update', 'system:notice:delete', 'system:notice:query',
                    'system:tenant:create', 'system:tenant:update', 'system:tenant:delete', 'system:tenant:query', 'system:tenant:update-package',
                    'system:login-log:query', 'system:login-log:export',
                    'system:operate-log:query', 'system:operate-log:export',
                  ],
                  menus: d.menus,
                };
              }},

              // 用户分页
              { match: (m, p) => m === 'GET' && p === '/api/system/user/page', handle: (_, q) => {
                const d = loadMockData();
                const page = Number(q.page) || 1, ps = Number(q.pageSize) || 10;
                return { list: d.users.slice((page-1)*ps, page*ps), total: d.users.length, page, pageSize: ps };
              }},
              { match: (m, p) => m === 'GET' && p.startsWith('/api/system/user/profile'), handle: () => ({ ...loadMockData().users[0], roleIds: [1] }) },
              { match: (m, p) => m === 'POST' && p === '/api/system/user', handle: (b) => b.id || Date.now() },
              { match: (m, p) => m === 'PUT' && p === '/api/system/user', handle: () => null },
              { match: (m, p) => m === 'GET' && p.match(/^\/api\/system\/user\/\d+$/), handle: () => loadMockData().users[0] },
              { match: (m, p) => m === 'DELETE' && p.match(/^\/api\/system\/user\/\d+$/), handle: () => null },

              // 角色分页
              { match: (m, p) => m === 'GET' && p === '/api/system/role/page', handle: (_, q) => {
                const d = loadMockData();
                return { list: d.roles, total: d.roles.length, page: Number(q.page)||1, pageSize: Number(q.pageSize)||10 };
              }},
              { match: (m, p) => m === 'POST' && p === '/api/system/role', handle: (b) => b.id || Date.now() },
              { match: (m, p) => m === 'PUT' && p === '/api/system/role', handle: () => null },
              { match: (m, p) => m === 'PUT' && p === '/api/system/role/assign-menus', handle: () => null },
              { match: (m, p) => m === 'DELETE' && p.match(/^\/api\/system\/role\/\d+$/), handle: () => null },
              { match: (m, p) => m === 'GET' && p.match(/^\/api\/system\/role\/\d+$/), handle: () => ({ ...loadMockData().roles[0], menuIds: [10, 11, 12, 13, 14, 15, 16] }) },

              // 菜单树
              { match: (m, p) => m === 'GET' && p === '/api/system/menu/tree', handle: () => loadMockData().menus },
              { match: (m, p) => m === 'POST' && p === '/api/system/menu', handle: (b) => b.id || Date.now() },
              { match: (m, p) => m === 'PUT' && p === '/api/system/menu', handle: () => null },
              { match: (m, p) => m === 'DELETE' && p.match(/^\/api\/system\/menu\/\d+$/), handle: () => null },

              // 部门树
              { match: (m, p) => m === 'GET' && (p === '/api/system/dept/tree' || p === '/api/system/dept/list'), handle: () => loadMockData().depts },
              { match: (m, p) => m === 'POST' && p === '/api/system/dept', handle: (b) => b.id || Date.now() },
              { match: (m, p) => m === 'PUT' && p === '/api/system/dept', handle: () => null },
              { match: (m, p) => m === 'DELETE' && p.match(/^\/api\/system\/dept\/\d+$/), handle: () => null },

              // 岗位分页
              { match: (m, p) => m === 'GET' && p === '/api/system/post/page', handle: (_, q) => {
                const d = loadMockData();
                return { list: d.posts, total: d.posts.length, page: Number(q.page)||1, pageSize: Number(q.pageSize)||10 };
              }},
              { match: (m, p) => m === 'POST' && p === '/api/system/post', handle: (b) => b.id || Date.now() },
              { match: (m, p) => m === 'PUT' && p === '/api/system/post', handle: () => null },
              { match: (m, p) => m === 'DELETE' && p.match(/^\/api\/system\/post\/\d+$/), handle: () => null },

              // 岗位 & 字典精简列表
              { match: (m, p) => m === 'GET' && p === '/api/system/post/simple-list', handle: () => loadMockData().posts },
              { match: (m, p) => m === 'GET' && p === '/api/system/dict-data/simple-list', handle: () => [] },
            ];

            // 读取请求体
            function readBody(req: any): Promise<any> {
              return new Promise((resolve) => {
                let data = '';
                req.on('data', (c: any) => (data += c));
                req.on('end', () => { try { resolve(JSON.parse(data)); } catch { resolve({}); } });
              });
            }

            // 解析 query
            function parseQ(url: string) {
              const p: Record<string, string> = {};
              (url.split('?')[1] || '').split('&').forEach((pair) => {
                const [k, v] = pair.split('=');
                if (k) p[k] = decodeURIComponent(v || '');
              });
              return p;
            }

            server.middlewares.use(async (req: any, res: any, next: any) => {
              const url = req.url || '/';
              if (!url.startsWith('/api/')) return next();
              if (!useMock) return next();

              if (req.method === 'OPTIONS') {
                res.setHeader('Access-Control-Allow-Origin', '*');
                res.setHeader('Access-Control-Allow-Methods', 'GET,POST,PUT,DELETE');
                res.setHeader('Access-Control-Allow-Headers', 'Content-Type,Authorization');
                res.statusCode = 204;
                res.end();
                return;
              }

              const method = (req.method || 'GET').toUpperCase();
              const pathOnly = url.split('?')[0];
              const handler = handlers.find((h) => h.match(method, pathOnly));

              if (handler) {
                const body = method === 'POST' || method === 'PUT' ? await readBody(req) : {};
                const result = handler.handle(body, parseQ(url));
                res.setHeader('Content-Type', 'application/json');
                res.end(ok(result));
                return;
              }
              next();
            });
          },
        },
      ],
      server: {
        allowedHosts: true,
        proxy: {
          '/api': {
            changeOrigin: true,
            target: 'http://127.0.0.1:8081',
            ws: true,
          },
        },
      },
    },
  };
});
