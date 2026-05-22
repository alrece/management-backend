-- 定时任务管理菜单
INSERT INTO sys_menu (id, parent_id, name, type, path, component, permission, icon, sort, visible, status, creator, updater) VALUES
(100, 0, '定时任务', 1, '/job', '', '', 'clock-circle', 20, 0, 0, 1, 1),
(101, 100, '任务管理', 2, 'task', '/job/task/index', 'job:task:query', '', 1, 0, 0, 1, 1),
(102, 100, '执行日志', 2, 'log', '/job/log/index', 'job:log:query', '', 2, 0, 0, 1, 1);

-- 工作流管理菜单
INSERT INTO sys_menu (id, parent_id, name, type, path, component, permission, icon, sort, visible, status, creator, updater) VALUES
(200, 0, '工作流管理', 1, '/workflow', '', '', 'branches', 30, 0, 0, 1, 1),
(201, 200, '流程分类', 2, 'category', '/workflow/category/index', 'workflow:category:query', '', 1, 0, 0, 1, 1),
(202, 200, '流程管理', 2, 'workflow', '/workflow/workflow/index', 'workflow:workflow:query', '', 2, 0, 0, 1, 1),
(203, 200, '执行实例', 2, 'instance', '/workflow/instance/index', 'workflow:instance:query', '', 3, 0, 0, 1, 1);

-- 按钮权限（定时任务）
INSERT INTO sys_menu (id, parent_id, name, type, path, component, permission, icon, sort, visible, status, creator, updater) VALUES
(1011, 101, '任务创建', 3, '', '', 'job:task:create', '', 1, 0, 0, 1, 1),
(1012, 101, '任务编辑', 3, '', '', 'job:task:update', '', 2, 0, 0, 1, 1),
(1013, 101, '任务删除', 3, '', '', 'job:task:delete', '', 3, 0, 0, 1, 1),
(1014, 101, '手动触发', 3, '', '', 'job:task:trigger', '', 4, 0, 0, 1, 1),
(1021, 102, '日志查询', 3, '', '', 'job:log:query', '', 1, 0, 0, 1, 1);

-- 按钮权限（工作流）
INSERT INTO sys_menu (id, parent_id, name, type, path, component, permission, icon, sort, visible, status, creator, updater) VALUES
(2011, 201, '分类创建', 3, '', '', 'workflow:category:create', '', 1, 0, 0, 1, 1),
(2012, 201, '分类编辑', 3, '', '', 'workflow:category:update', '', 2, 0, 0, 1, 1),
(2013, 201, '分类删除', 3, '', '', 'workflow:category:delete', '', 3, 0, 0, 1, 1),
(2021, 202, '流程创建', 3, '', '', 'workflow:workflow:create', '', 1, 0, 0, 1, 1),
(2022, 202, '流程编辑', 3, '', '', 'workflow:workflow:update', '', 2, 0, 0, 1, 1),
(2023, 202, '流程删除', 3, '', '', 'workflow:workflow:delete', '', 3, 0, 0, 1, 1),
(2024, 202, '流程激活', 3, '', '', 'workflow:workflow:activate', '', 4, 0, 0, 1, 1),
(2025, 202, '流程执行', 3, '', '', 'workflow:workflow:execute', '', 5, 0, 0, 1, 1),
(2031, 203, '实例查询', 3, '', '', 'workflow:instance:query', '', 1, 0, 0, 1, 1);
