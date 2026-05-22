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
                  { id: 100, parentId: 0, menuName: '定时任务', menuType: 1, path: '/job', component: '', permission: '', icon: 'ant-design:clock-circle-outlined', sort: 20, visible: 0, status: 0, children: [
                    { id: 101, parentId: 100, menuName: '任务管理', menuType: 2, path: 'task', component: '/job/task/index', permission: 'job:task:query', icon: '', sort: 1, visible: 0, status: 0, children: [] },
                    { id: 102, parentId: 100, menuName: '执行日志', menuType: 2, path: 'log', component: '/job/log/index', permission: 'job:log:query', icon: '', sort: 2, visible: 0, status: 0, children: [] },
                  ] },
                  { id: 200, parentId: 0, menuName: '工作流管理', menuType: 1, path: '/workflow', component: '', permission: '', icon: 'ant-design:branches-outlined', sort: 30, visible: 0, status: 0, children: [
                    { id: 201, parentId: 200, menuName: '流程分类', menuType: 2, path: 'category', component: '/workflow/category/index', permission: 'workflow:category:query', icon: '', sort: 1, visible: 0, status: 0, children: [] },
                    { id: 202, parentId: 200, menuName: '流程管理', menuType: 2, path: 'workflow', component: '/workflow/workflow/index', permission: 'workflow:workflow:query', icon: '', sort: 2, visible: 0, status: 0, children: [] },
                    { id: 203, parentId: 200, menuName: '执行实例', menuType: 2, path: 'instance', component: '/workflow/instance/index', permission: 'workflow:instance:query', icon: '', sort: 3, visible: 0, status: 0, children: [] },
                  ] },
                ],
                depts: [
                  { id: 100, parentId: 0, name: '总公司', sort: 0, leaderUserId: 1, status: 0, createTime: '2026-01-01 00:00:00' },
                  { id: 101, parentId: 100, name: '技术部', sort: 1, leaderUserId: 2, status: 0, createTime: '2026-01-01 00:00:00' },
                ],
                posts: [
                  { id: 1, postCode: 'dev', postName: '开发工程师', sort: 1, status: 0, remark: '', createTime: '2026-01-01 00:00:00' },
                ],
                jobTasks: [
                  { id: 1, name: '数据备份', handler: 'DataBackupHandler', cronExpr: '0 0 2 * * ?', params: '{}', status: 0, remark: '每日凌晨2点', creator: 1 },
                  { id: 2, name: '缓存清理', handler: 'CacheCleanHandler', cronExpr: '0 0 3 * * ?', params: '{}', status: 0, remark: '每日凌晨3点', creator: 1 },
                ],
                jobExecLogs: [
                  { id: 1, taskId: 1, taskName: '数据备份', triggerType: 1, status: 1, durationMs: 3200, result: '备份完成', error: '', startTime: '2026-05-22 02:00:00' },
                  { id: 2, taskId: 2, taskName: '缓存清理', triggerType: 2, status: 1, durationMs: 1500, result: '清理完成', error: '', startTime: '2026-05-22 03:00:00' },
                ],
                wfCategories: [
                  { id: 1, name: '审批流程', parentId: 0, sort: 0, children: [] },
                  { id: 2, name: '数据同步', parentId: 0, sort: 1, children: [] },
                ],
                wfWorkflows: [
                  { id: 1, name: '请假审批', categoryId: 1, n8nWorkflowId: '', paramsSchema: '{}', status: 'ACTIVE', creator: 1, createdAt: '2026-05-01', updater: 1, updatedAt: '2026-05-01' },
                  { id: 2, name: '用户数据同步', categoryId: 2, n8nWorkflowId: '', paramsSchema: '{}', status: 'DRAFT', creator: 1, createdAt: '2026-05-10', updater: 1, updatedAt: '2026-05-10' },
                ],
                wfInstances: [
                  { id: 1, workflowId: 1, n8nExecutionId: 'exec_001', status: 'SUCCESS', variables: '{}', result: '{"approved": true}', durationMs: 1500, errorMsg: '', startedAt: '2026-05-22 10:00:00', finishedAt: '2026-05-22 10:00:01' },
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
                    'job:task:create', 'job:task:update', 'job:task:delete', 'job:task:query', 'job:task:trigger',
                    'job:log:query',
                    'workflow:category:create', 'workflow:category:update', 'workflow:category:delete', 'workflow:category:query',
                    'workflow:workflow:create', 'workflow:workflow:update', 'workflow:workflow:delete', 'workflow:workflow:query',
                    'workflow:workflow:activate', 'workflow:workflow:execute',
                    'workflow:instance:query',
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

              // === Job 定时任务 ===
              { match: (m, p) => m === 'GET' && p === '/api/job/tasks/page', handle: (_, q) => {
                const d = loadMockData();
                const page = Number(q.page) || 1, ps = Number(q.pageSize) || 10;
                return { list: d.jobTasks.slice((page-1)*ps, page*ps), total: d.jobTasks.length };
              }},
              { match: (m, p) => m === 'POST' && p === '/api/job/tasks', handle: (b) => Date.now() },
              { match: (m, p) => m === 'PUT' && p.match(/^\/api\/job\/tasks\/\d+$/), handle: () => null },
              { match: (m, p) => m === 'DELETE' && p.match(/^\/api\/job\/tasks\/\d+$/), handle: () => null },
              { match: (m, p) => m === 'POST' && p.match(/^\/api\/job\/tasks\/\d+\/trigger$/), handle: () => Date.now() },
              { match: (m, p) => m === 'GET' && p === '/api/job/execution-logs/page', handle: (_, q) => {
                const d = loadMockData();
                const page = Number(q.page) || 1, ps = Number(q.pageSize) || 10;
                return { list: d.jobExecLogs.slice((page-1)*ps, page*ps), total: d.jobExecLogs.length };
              }},

              // === Workflow 工作流 ===
              { match: (m, p) => m === 'GET' && p === '/api/workflow/categories/tree', handle: () => loadMockData().wfCategories },
              { match: (m, p) => m === 'POST' && p === '/api/workflow/categories', handle: (b) => Date.now() },
              { match: (m, p) => m === 'PUT' && p.match(/^\/api\/workflow\/categories\/\d+$/), handle: () => null },
              { match: (m, p) => m === 'DELETE' && p.match(/^\/api\/workflow\/categories\/\d+$/), handle: () => null },
              { match: (m, p) => m === 'GET' && p === '/api/workflow/workflows/page', handle: (_, q) => {
                const d = loadMockData();
                const page = Number(q.page) || 1, ps = Number(q.pageSize) || 10;
                return { list: d.wfWorkflows.slice((page-1)*ps, page*ps), total: d.wfWorkflows.length };
              }},
              { match: (m, p) => m === 'GET' && p.match(/^\/api\/workflow\/workflows\/\d+$/), handle: () => loadMockData().wfWorkflows[0] },
              { match: (m, p) => m === 'POST' && p === '/api/workflow/workflows', handle: (b) => Date.now() },
              { match: (m, p) => m === 'PUT' && p.match(/^\/api\/workflow\/workflows\/\d+$/), handle: () => null },
              { match: (m, p) => m === 'DELETE' && p.match(/^\/api\/workflow\/workflows\/\d+$/), handle: () => null },
              { match: (m, p) => m === 'PUT' && p.match(/^\/api\/workflow\/workflows\/\d+\/activate$/), handle: () => null },
              { match: (m, p) => m === 'PUT' && p.match(/^\/api\/workflow\/workflows\/\d+\/deactivate$/), handle: () => null },
              { match: (m, p) => m === 'POST' && p.match(/^\/api\/workflow\/workflows\/\d+\/execute$/), handle: () => Date.now() },
              { match: (m, p) => m === 'GET' && p === '/api/workflow/instances/page', handle: (_, q) => {
                const d = loadMockData();
                const page = Number(q.page) || 1, ps = Number(q.pageSize) || 10;
                return { list: d.wfInstances.slice((page-1)*ps, page*ps), total: d.wfInstances.length };
              }},
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
          '/api/system': {
            changeOrigin: true,
            target: 'http://127.0.0.1:8081',
            ws: true,
          },
          '/api/job': {
            changeOrigin: true,
            target: 'http://127.0.0.1:8082',
            ws: true,
          },
          '/api/workflow': {
            changeOrigin: true,
            target: 'http://127.0.0.1:8083',
            ws: true,
          },
        },
      },
    },
  };
});
