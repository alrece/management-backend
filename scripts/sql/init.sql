-- ==========================================
-- management-backend 数据库初始化
-- 综合两个 Java 项目的核心表设计
-- ==========================================

CREATE DATABASE IF NOT EXISTS `management_backend`
  DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `management_backend`;

-- 租户表
CREATE TABLE `sys_tenant` (
    `id`              BIGINT       NOT NULL COMMENT '租户ID',
    `tenant_name`     VARCHAR(100) NOT NULL COMMENT '租户名称',
    `contact_user_id` BIGINT       DEFAULT NULL COMMENT '联系人ID',
    `status`          TINYINT      NOT NULL DEFAULT 0 COMMENT '状态(0正常 1停用)',
    `expire_time`     DATETIME     DEFAULT NULL COMMENT '过期时间',
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
    `permission`  VARCHAR(100) DEFAULT '' COMMENT '权限标识',
    `icon`        VARCHAR(50)  DEFAULT '',
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
INSERT INTO `sys_tenant` (`id`, `tenant_name`, `status`) VALUES (1, '默认租户', 0);
INSERT INTO `sys_user` (`id`, `tenant_id`, `username`, `password`, `nickname`, `status`)
  VALUES (1, 1, 'admin', '$2a$10$VQBl2noKFPH/MSOsOq7a2.tdJqIckGw8MKTqPYxGqv3Rp7.q5mFJO', '超级管理员', 0);
INSERT INTO `sys_role` (`id`, `tenant_id`, `role_name`, `role_code`, `data_scope`, `status`)
  VALUES (1, 1, '超级管理员', 'super_admin', 1, 0);
INSERT INTO `sys_user_role` (`id`, `user_id`, `role_id`) VALUES (1, 1, 1);
