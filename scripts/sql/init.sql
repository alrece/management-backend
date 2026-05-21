-- ==========================================
-- management-backend 数据库初始化
-- 综合两个 Java 项目的核心表设计
-- ==========================================

CREATE DATABASE IF NOT EXISTS `management_backend`
  DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `management_backend`;

-- 租户表（默认库，不受租户过滤）
CREATE TABLE `sys_tenant` (
    `id`              BIGINT       NOT NULL COMMENT '租户ID',
    `tenant_name`     VARCHAR(100) NOT NULL COMMENT '租户名称',
    `contact_name`    VARCHAR(30)  DEFAULT '' COMMENT '联系人',
    `contact_mobile`  VARCHAR(20)  DEFAULT '' COMMENT '联系电话',
    `db_name`         VARCHAR(100) DEFAULT '' COMMENT '租户数据库名',
    `status`          TINYINT      NOT NULL DEFAULT 0 COMMENT '状态(0正常 1停用)',
    `expire_time`     DATETIME     DEFAULT NULL COMMENT '过期时间',
    `init_status`     TINYINT      NOT NULL DEFAULT 0 COMMENT '初始化状态(0待初始化 1初始化中 2已就绪 3失败)',
    `remark`          VARCHAR(500) DEFAULT '' COMMENT '备注',
    `creator`         BIGINT       DEFAULT NULL,
    `create_time`     DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updater`         BIGINT       DEFAULT NULL,
    `update_time`     DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted`         BIT(1)       NOT NULL DEFAULT b'0',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB COMMENT='租户表';

-- 部门表
CREATE TABLE `sys_dept` (
    `id`          BIGINT       NOT NULL COMMENT '部门ID',
    `tenant_id`   BIGINT       NOT NULL DEFAULT 0 COMMENT '租户编号',
    `parent_id`   BIGINT       NOT NULL DEFAULT 0 COMMENT '父部门ID',
    `ancestors`   VARCHAR(200) DEFAULT '' COMMENT '祖级列表',
    `dept_name`   VARCHAR(50)  NOT NULL COMMENT '部门名称',
    `sort`        INT          NOT NULL DEFAULT 0,
    `leader`      VARCHAR(30)  DEFAULT '' COMMENT '负责人',
    `status`      TINYINT      NOT NULL DEFAULT 0 COMMENT '状态(0正常 1停用)',
    `creator`     BIGINT       DEFAULT NULL,
    `create_time` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updater`     BIGINT       DEFAULT NULL,
    `update_time` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted`     BIT(1)       NOT NULL DEFAULT b'0',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB COMMENT='部门表';

-- 用户表
CREATE TABLE `sys_user` (
    `id`          BIGINT       NOT NULL COMMENT '用户ID',
    `tenant_id`   BIGINT       NOT NULL DEFAULT 0 COMMENT '租户编号',
    `username`    VARCHAR(30)  NOT NULL COMMENT '用户名',
    `password`    VARCHAR(100) NOT NULL COMMENT '密码',
    `nickname`    VARCHAR(30)  DEFAULT '' COMMENT '昵称',
    `email`       VARCHAR(50)  DEFAULT '' COMMENT '邮箱',
    `mobile`      VARCHAR(20)  DEFAULT '' COMMENT '手机号',
    `sex`         TINYINT      DEFAULT 0 COMMENT '性别(0未知 1男 2女)',
    `avatar`      VARCHAR(255) DEFAULT '' COMMENT '头像',
    `status`      TINYINT      NOT NULL DEFAULT 0 COMMENT '状态(0正常 1停用)',
    `dept_id`     BIGINT       DEFAULT NULL COMMENT '部门ID',
    `create_dept` BIGINT       DEFAULT NULL COMMENT '创建部门',
    `remark`      VARCHAR(500) DEFAULT '' COMMENT '备注',
    `creator`     BIGINT       DEFAULT NULL,
    `create_time` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updater`     BIGINT       DEFAULT NULL,
    `update_time` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted`     BIT(1)       NOT NULL DEFAULT b'0',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_username` (`username`, `tenant_id`)
) ENGINE=InnoDB COMMENT='用户表';

-- 角色表
CREATE TABLE `sys_role` (
    `id`          BIGINT       NOT NULL COMMENT '角色ID',
    `tenant_id`   BIGINT       NOT NULL DEFAULT 0 COMMENT '租户编号',
    `role_name`   VARCHAR(50)  NOT NULL COMMENT '角色名称',
    `role_code`   VARCHAR(50)  NOT NULL COMMENT '角色编码',
    `sort`        INT          NOT NULL DEFAULT 0,
    `data_scope`  TINYINT      NOT NULL DEFAULT 1 COMMENT '数据范围(1全部 2自定义 3本部门 4本部门及以下 5仅本人)',
t`data_scope_dept_ids` VARCHAR(500) DEFAULT '' COMMENT '数据范围部门ID(自定义)',
    `status`      TINYINT      NOT NULL DEFAULT 0,
    `remark`      VARCHAR(500) DEFAULT '',
    `creator`     BIGINT       DEFAULT NULL,
    `create_time` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updater`     BIGINT       DEFAULT NULL,
    `update_time` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted`     BIT(1)       NOT NULL DEFAULT b'0',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB COMMENT='角色表';

-- 菜单表（不受租户过滤）
CREATE TABLE `sys_menu` (
    `id`          BIGINT       NOT NULL COMMENT '菜单ID',
    `parent_id`   BIGINT       NOT NULL DEFAULT 0 COMMENT '父菜单ID',
    `menu_name`   VARCHAR(50)  NOT NULL COMMENT '菜单名称',
    `menu_type`   TINYINT      NOT NULL COMMENT '类型(1目录 2菜单 3按钮)',
    `path`        VARCHAR(200) DEFAULT '' COMMENT '路由路径',
t`component`  VARCHAR(200) DEFAULT '' COMMENT '组件路径',
    `permission`  VARCHAR(100) DEFAULT '' COMMENT '权限标识',
    `icon`        VARCHAR(50)  DEFAULT '',
t`visible`    TINYINT      NOT NULL DEFAULT 0 COMMENT '是否可见(0可见 1隐藏)',
    `sort`        INT          NOT NULL DEFAULT 0,
    `status`      TINYINT      NOT NULL DEFAULT 0,
    `creator`     BIGINT       DEFAULT NULL,
    `create_time` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updater`     BIGINT       DEFAULT NULL,
    `update_time` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted`     BIT(1)       NOT NULL DEFAULT b'0',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB COMMENT='菜单权限表';

-- 用户角色关联
CREATE TABLE `sys_user_role` (
    `id`          BIGINT   NOT NULL,
    `user_id`     BIGINT   NOT NULL COMMENT '用户ID',
    `role_id`     BIGINT   NOT NULL COMMENT '角色ID',
    `creator`     BIGINT   DEFAULT NULL,
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updater`     BIGINT   DEFAULT NULL,
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted`     BIT(1)   NOT NULL DEFAULT b'0',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB COMMENT='用户角色关联表';

-- 角色菜单关联（不受租户过滤）
CREATE TABLE `sys_role_menu` (
    `id`          BIGINT   NOT NULL,
    `role_id`     BIGINT   NOT NULL COMMENT '角色ID',
    `menu_id`     BIGINT   NOT NULL COMMENT '菜单ID',
    `creator`     BIGINT   DEFAULT NULL,
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updater`     BIGINT   DEFAULT NULL,
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted`     BIT(1)   NOT NULL DEFAULT b'0',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB COMMENT='角色菜单关联表';

-- 字典类型
CREATE TABLE `sys_dict_type` (
    `id`          BIGINT       NOT NULL,
    `tenant_id`   BIGINT       NOT NULL DEFAULT 0,
    `dict_name`   VARCHAR(100) NOT NULL COMMENT '字典名称',
    `dict_type`   VARCHAR(100) NOT NULL COMMENT '字典类型',
    `status`      TINYINT      NOT NULL DEFAULT 0,
    `remark`      VARCHAR(500) DEFAULT '',
    `creator`     BIGINT       DEFAULT NULL,
    `create_time` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updater`     BIGINT       DEFAULT NULL,
    `update_time` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted`     BIT(1)       NOT NULL DEFAULT b'0',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_dict_type` (`dict_type`, `tenant_id`)
) ENGINE=InnoDB COMMENT='字典类型表';

-- 字典数据
CREATE TABLE `sys_dict_data` (
    `id`          BIGINT       NOT NULL,
    `tenant_id`   BIGINT       NOT NULL DEFAULT 0,
    `dict_type`   VARCHAR(100) NOT NULL COMMENT '字典类型',
    `label`       VARCHAR(100) NOT NULL COMMENT '字典标签',
    `value`       VARCHAR(100) NOT NULL COMMENT '字典值',
    `sort`        INT          NOT NULL DEFAULT 0,
    `status`      TINYINT      NOT NULL DEFAULT 0,
    `remark`      VARCHAR(500) DEFAULT '',
    `creator`     BIGINT       DEFAULT NULL,
    `create_time` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updater`     BIGINT       DEFAULT NULL,
    `update_time` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted`     BIT(1)       NOT NULL DEFAULT b'0',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB COMMENT='字典数据表';

-- 初始数据
INSERT INTO `sys_tenant` (`id`, `tenant_name`, `db_name`, `status`, `init_status`) VALUES (1, '默认租户', 'management_backend', 0, 2);
INSERT INTO `sys_user` (`id`, `tenant_id`, `username`, `password`, `nickname`, `status`)
  VALUES (1, 1, 'admin', '$2a$10$VQBl2noKFPH/MSOsOq7a2.tdJqIckGw8MKTqPYxGqv3Rp7.q5mFJO', '超级管理员', 0);
INSERT INTO `sys_role` (`id`, `tenant_id`, `role_name`, `role_code`, `data_scope`, `status`)
  VALUES (1, 1, '超级管理员', 'super_admin', 1, 0);
INSERT INTO `sys_user_role` (`id`, `user_id`, `role_id`) VALUES (1, 1, 1);

-- 初始菜单：目录
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `menu_type`, `path`, `component`, `permission`, `icon`, `sort`, `status`, `visible`) VALUES
(1, 0, '系统管理', 1, '/system', '', '', 'system', 1, 0, 0);

-- 菜单
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `menu_type`, `path`, `component`, `permission`, `icon`, `sort`, `status`, `visible`) VALUES
(2, 1, '用户管理', 2, '/system/user', 'system/user/index', '', 'user', 1, 0, 0),
(3, 1, '角色管理', 2, '/system/role', 'system/role/index', '', 'role', 2, 0, 0),
(4, 1, '菜单管理', 2, '/system/menu', 'system/menu/index', '', 'menu', 3, 0, 0),
(5, 1, '部门管理', 2, '/system/dept', 'system/dept/index', '', 'dept', 4, 0, 0),
(6, 1, '岗位管理', 2, '/system/post', 'system/post/index', '', 'post', 5, 0, 0);

-- 按钮权限
INSERT INTO `sys_menu` (`id`, `parent_id`, `menu_name`, `menu_type`, `path`, `component`, `permission`, `icon`, `sort`, `status`, `visible`) VALUES
(100, 2, '用户新增', 3, '', '', 'system:user:create', '', 1, 0, 0),
(101, 2, '用户修改', 3, '', '', 'system:user:update', '', 2, 0, 0),
(102, 2, '用户删除', 3, '', '', 'system:user:delete', '', 3, 0, 0),
(103, 2, '用户查询', 3, '', '', 'system:user:query', '', 4, 0, 0),
(200, 3, '角色新增', 3, '', '', 'system:role:create', '', 1, 0, 0),
(201, 3, '角色修改', 3, '', '', 'system:role:update', '', 2, 0, 0),
(202, 3, '角色删除', 3, '', '', 'system:role:delete', '', 3, 0, 0),
(203, 3, '角色查询', 3, '', '', 'system:role:query', '', 4, 0, 0),
(300, 4, '菜单新增', 3, '', '', 'system:menu:create', '', 1, 0, 0),
(301, 4, '菜单修改', 3, '', '', 'system:menu:update', '', 2, 0, 0),
(302, 4, '菜单删除', 3, '', '', 'system:menu:delete', '', 3, 0, 0),
(303, 4, '菜单查询', 3, '', '', 'system:menu:query', '', 4, 0, 0),
(400, 5, '部门新增', 3, '', '', 'system:dept:create', '', 1, 0, 0),
(401, 5, '部门修改', 3, '', '', 'system:dept:update', '', 2, 0, 0),
(402, 5, '部门删除', 3, '', '', 'system:dept:delete', '', 3, 0, 0),
(403, 5, '部门查询', 3, '', '', 'system:dept:query', '', 4, 0, 0),
(500, 6, '岗位新增', 3, '', '', 'system:post:create', '', 1, 0, 0),
(501, 6, '岗位修改', 3, '', '', 'system:post:update', '', 2, 0, 0),
(502, 6, '岗位删除', 3, '', '', 'system:post:delete', '', 3, 0, 0),
(503, 6, '岗位查询', 3, '', '', 'system:post:query', '', 4, 0, 0);

-- 超级管理员分配所有菜单
INSERT INTO `sys_role_menu` (`id`, `role_id`, `menu_id`) VALUES
(1, 1, 1), (2, 1, 2), (3, 1, 3), (4, 1, 4), (5, 1, 5), (6, 1, 6),
(10, 1, 100), (11, 1, 101), (12, 1, 102), (13, 1, 103),
(20, 1, 200), (21, 1, 201), (22, 1, 202), (23, 1, 203),
(30, 1, 300), (31, 1, 301), (32, 1, 302), (33, 1, 303),
(40, 1, 400), (41, 1, 401), (42, 1, 402), (43, 1, 403),
(50, 1, 500), (51, 1, 501), (52, 1, 502), (53, 1, 503);

-- 用户岗位关联
CREATE TABLE IF NOT EXISTS `sys_user_post` (
    `id`          BIGINT   NOT NULL,
    `user_id`     BIGINT   NOT NULL COMMENT '用户ID',
    `post_id`     BIGINT   NOT NULL COMMENT '岗位ID',
    `creator`     BIGINT   DEFAULT NULL,
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updater`     BIGINT   DEFAULT NULL,
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted`     BIT(1)   NOT NULL DEFAULT b'0',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB COMMENT='用户岗位关联表';

-- 系统参数表
CREATE TABLE IF NOT EXISTS `sys_param` (
    `id`          BIGINT       NOT NULL,
    `tenant_id`   BIGINT       NOT NULL DEFAULT 0,
    `param_key`   VARCHAR(100) NOT NULL COMMENT '参数键',
    `param_value` VARCHAR(500) NOT NULL COMMENT '参数值',
    `param_type`  TINYINT      NOT NULL DEFAULT 1 COMMENT '类型(1系统 2用户)',
    `status`      TINYINT      NOT NULL DEFAULT 0,
    `remark`      VARCHAR(500) DEFAULT '',
    `creator`     BIGINT       DEFAULT NULL,
    `create_time` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updater`     BIGINT       DEFAULT NULL,
    `update_time` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted`     BIT(1)       NOT NULL DEFAULT b'0',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_param_key` (`param_key`, `tenant_id`)
) ENGINE=InnoDB COMMENT='系统参数表';

-- 通知公告表
CREATE TABLE IF NOT EXISTS `sys_notice` (
    `id`           BIGINT       NOT NULL,
    `tenant_id`    BIGINT       NOT NULL DEFAULT 0,
    `notice_title` VARCHAR(200) NOT NULL COMMENT '公告标题',
    `notice_type`  TINYINT      NOT NULL DEFAULT 1 COMMENT '类型(1通知 2公告 3提醒)',
    `content`      TEXT         COMMENT '内容',
    `status`       TINYINT      NOT NULL DEFAULT 0,
    `remark`       VARCHAR(500) DEFAULT '',
    `creator`      BIGINT       DEFAULT NULL,
    `create_time`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updater`      BIGINT       DEFAULT NULL,
    `update_time`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted`      BIT(1)       NOT NULL DEFAULT b'0',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB COMMENT='通知公告表';

-- 操作日志表
CREATE TABLE IF NOT EXISTS `sys_oper_log` (
    `id`              BIGINT       NOT NULL COMMENT '日志ID',
    `tenant_id`       BIGINT       DEFAULT 0 COMMENT '租户ID',
    `title`           VARCHAR(100) DEFAULT '' COMMENT '操作模块',
    `business_type`   TINYINT      DEFAULT 0 COMMENT '业务类型(0其它 1新增 2修改 3删除)',
    `method`          VARCHAR(200) DEFAULT '' COMMENT '请求方法',
    `request_url`     VARCHAR(500) DEFAULT '' COMMENT '请求URL',
    `oper_ip`         VARCHAR(50)  DEFAULT '' COMMENT '操作IP',
    `oper_user_id`    BIGINT       DEFAULT 0 COMMENT '操作者ID',
    `oper_name`       VARCHAR(50)  DEFAULT '' COMMENT '操作者',
    `request_id`      VARCHAR(50)  DEFAULT '' COMMENT '请求ID',
    `status`          TINYINT      DEFAULT 0 COMMENT '状态(0成功 1失败)',
    `error_msg`       TEXT         COMMENT '错误消息',
    `oper_time`       DATETIME     NOT NULL COMMENT '操作时间',
    PRIMARY KEY (`id`),
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_request_id` (`request_id`)
) ENGINE=InnoDB COMMENT='操作日志表';

-- 登录日志表
CREATE TABLE IF NOT EXISTS `sys_login_log` (
    `id`              BIGINT       NOT NULL COMMENT '日志ID',
    `username`        VARCHAR(50)  DEFAULT '' COMMENT '用户名',
    `login_ip`        VARCHAR(50)  DEFAULT '' COMMENT '登录IP',
    `login_location`  VARCHAR(100) DEFAULT '' COMMENT '登录地点',
    `browser`         VARCHAR(50)  DEFAULT '' COMMENT '浏览器',
    `os`              VARCHAR(50)  DEFAULT '' COMMENT '操作系统',
    `status`          TINYINT      DEFAULT 0 COMMENT '状态(0成功 1失败)',
    `msg`             VARCHAR(200) DEFAULT '' COMMENT '消息',
    `login_time`      DATETIME     NOT NULL COMMENT '登录时间',
    PRIMARY KEY (`id`),
    INDEX `idx_username` (`username`)
) ENGINE=InnoDB COMMENT='登录日志表';
