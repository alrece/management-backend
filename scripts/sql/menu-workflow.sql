-- 工作流模块菜单和权限 SQL (PostgreSQL)
-- 插入到 sys_menu 表

-- 工作流管理（一级菜单）
INSERT INTO sys_menu (id, name, permission, type, sort, parent_id, path, icon, component, status, deleted)
VALUES (2000, '工作流管理', '', 1, 50, 0, '/workflow', 'deployment-unit', '', 0, 0);

-- 分类管理（二级菜单）
INSERT INTO sys_menu (id, name, permission, type, sort, parent_id, path, icon, component, status, deleted)
VALUES (2001, '分类管理', 'workflow:category:query', 2, 1, 2000, 'category', 'folder', 'workflow/category/index', 0, 0);

-- 工作流列表（二级菜单）
INSERT INTO sys_menu (id, name, permission, type, sort, parent_id, path, icon, component, status, deleted)
VALUES (2002, '工作流列表', 'workflow:workflow:query', 2, 2, 2000, 'workflow', 'apartment', 'workflow/workflow/index', 0, 0);

-- 工作流编辑器（隐藏路由）
INSERT INTO sys_menu (id, name, permission, type, sort, parent_id, path, icon, component, status, deleted)
VALUES (2003, '工作流编辑器', 'workflow:workflow:update', 2, 3, 2000, 'editor/:id', 'edit', 'workflow/editor/index', 0, 0);

-- 执行实例（二级菜单）
INSERT INTO sys_menu (id, name, permission, type, sort, parent_id, path, icon, component, status, deleted)
VALUES (2004, '执行实例', 'workflow:execution:query', 2, 4, 2000, 'instance', 'profile', 'workflow/instance/index', 0, 0);

-- 工作流按钮权限
INSERT INTO sys_menu (id, name, permission, type, sort, parent_id, path, icon, component, status, deleted) VALUES
(2010, '创建工作流', 'workflow:workflow:create', 3, 1, 2002, '', '', '', 0, 0),
(2011, '更新工作流', 'workflow:workflow:update', 3, 2, 2002, '', '', '', 0, 0),
(2012, '删除工作流', 'workflow:workflow:delete', 3, 3, 2002, '', '', '', 0, 0),
(2013, '查询工作流', 'workflow:workflow:query', 3, 4, 2002, '', '', '', 0, 0),
(2014, '执行工作流', 'workflow:workflow:execute', 3, 5, 2002, '', '', '', 0, 0),
(2020, '创建分类', 'workflow:category:create', 3, 1, 2001, '', '', '', 0, 0),
(2021, '更新分类', 'workflow:category:update', 3, 2, 2001, '', '', '', 0, 0),
(2022, '删除分类', 'workflow:category:delete', 3, 3, 2001, '', '', '', 0, 0),
(2023, '查询分类', 'workflow:category:query', 3, 4, 2001, '', '', '', 0, 0),
(2030, '查询执行', 'workflow:execution:query', 3, 1, 2004, '', '', '', 0, 0);
