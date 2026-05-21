/**
 * Mock 数据定义
 * 与 Go 后端响应格式一致：{code: 0, data: ..., msg: "success"}
 */

// 雪花 ID 生成器
let idCounter = 1900000000000000001;
function nextId() {
  return idCounter++;
}

// ========== 用户 ==========

const mockUsers = [
  {
    id: 1,
    username: 'admin',
    nickname: '管理员',
    email: 'admin@example.com',
    mobile: '13800000000',
    sex: 1,
    avatar: '',
    status: 0,
    deptId: 100,
    tenantId: 1,
    createTime: '2026-01-01 00:00:00',
    remark: '超级管理员',
  },
  {
    id: 2,
    username: 'user01',
    nickname: '普通用户',
    email: 'user01@example.com',
    mobile: '13800000001',
    sex: 2,
    avatar: '',
    status: 0,
    deptId: 101,
    tenantId: 1,
    createTime: '2026-01-15 10:30:00',
    remark: '',
  },
];

// ========== 角色 ==========

const mockRoles = [
  { id: 1, name: '超级管理员', code: 'super_admin', sort: 1, dataScope: 1, status: 0, remark: '超级管理员', createTime: '2026-01-01 00:00:00' },
  { id: 2, name: '普通角色', code: 'common', sort: 2, dataScope: 2, status: 0, remark: '', createTime: '2026-01-01 00:00:00' },
];

// ========== 菜单（树形） ==========

const mockMenus = [
  {
    id: 1, parentId: 0, menuName: '系统管理', menuType: 1, path: '/system', component: '',
    permission: '', icon: 'ant-design:setting-outlined', sort: 10, visible: 0, status: 0,
    children: [
      { id: 10, parentId: 1, menuName: '用户管理', menuType: 2, path: 'user', component: '/system/user/index', permission: 'system:user:query', icon: 'ant-design:user-outlined', sort: 1, visible: 0, status: 0, children: [] },
      { id: 11, parentId: 1, menuName: '角色管理', menuType: 2, path: 'role', component: '/system/role/index', permission: 'system:role:query', icon: 'ant-design:team-outlined', sort: 2, visible: 0, status: 0, children: [] },
      { id: 12, parentId: 1, menuName: '菜单管理', menuType: 2, path: 'menu', component: '/system/menu/index', permission: 'system:menu:query', icon: 'ant-design:menu-outlined', sort: 3, visible: 0, status: 0, children: [] },
      { id: 13, parentId: 1, menuName: '部门管理', menuType: 2, path: 'dept', component: '/system/dept/index', permission: 'system:dept:query', icon: 'ant-design:apartment-outlined', sort: 4, visible: 0, status: 0, children: [] },
      { id: 14, parentId: 1, menuName: '岗位管理', menuType: 2, path: 'post', component: '/system/post/index', permission: 'system:post:query', icon: 'ant-design:idcard-outlined', sort: 5, visible: 0, status: 0, children: [] },
      { id: 15, parentId: 1, menuName: '字典管理', menuType: 2, path: 'dict', component: '/system/dict/index', permission: '', icon: 'ant-design:book-outlined', sort: 6, visible: 0, status: 0, children: [] },
      { id: 16, parentId: 1, menuName: '通知公告', menuType: 2, path: 'notice', component: '/system/notice/index', permission: '', icon: 'ant-design:notification-outlined', sort: 7, visible: 0, status: 0, children: [] },
    ],
  },
  {
    id: 2, parentId: 0, menuName: '系统监控', menuType: 1, path: '/monitor', component: '',
    permission: '', icon: 'ant-design:eye-outlined', sort: 20, visible: 0, status: 0,
    children: [
      { id: 20, parentId: 2, menuName: '操作日志', menuType: 2, path: 'operatelog', component: '/system/operatelog/index', permission: '', icon: 'ant-design:file-text-outlined', sort: 1, visible: 0, status: 0, children: [] },
      { id: 21, parentId: 2, menuName: '登录日志', menuType: 2, path: 'loginlog', component: '/system/loginlog/index', permission: '', icon: 'ant-design:login-outlined', sort: 2, visible: 0, status: 0, children: [] },
    ],
  },
];

// ========== 部门（树形） ==========

const mockDepts = [
  { id: 100, parentId: 0, name: '总公司', sort: 0, leaderUserId: 1, status: 0, createTime: '2026-01-01 00:00:00' },
  { id: 101, parentId: 100, name: '技术部', sort: 1, leaderUserId: 2, status: 0, createTime: '2026-01-01 00:00:00' },
  { id: 102, parentId: 100, name: '市场部', sort: 2, leaderUserId: null, status: 0, createTime: '2026-01-01 00:00:00' },
  { id: 103, parentId: 100, name: '财务部', sort: 3, leaderUserId: null, status: 0, createTime: '2026-01-01 00:00:00' },
];

// ========== 岗位 ==========

const mockPosts = [
  { id: 1, postCode: 'ceo', postName: '董事长', sort: 1, status: 0, remark: '', createTime: '2026-01-01 00:00:00' },
  { id: 2, postCode: 'se', postName: '架构师', sort: 2, status: 0, remark: '', createTime: '2026-01-01 00:00:00' },
  { id: 3, postCode: 'dev', postName: '开发工程师', sort: 3, status: 0, remark: '', createTime: '2026-01-01 00:00:00' },
];

export { mockUsers, mockRoles, mockMenus, mockDepts, mockPosts, nextId };
