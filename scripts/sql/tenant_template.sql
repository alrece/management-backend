-- ==========================================
-- 租户独立库模板 Schema
-- 新建租户时复制此模板到租户独立数据库
-- ==========================================

-- 用户表
CREATE TABLE IF NOT EXISTS `sys_user` (
    `id`          BIGINT       NOT NULL COMMENT '用户ID',
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
    UNIQUE KEY `uk_username` (`username`)
) ENGINE=InnoDB COMMENT='用户表';

-- 角色表
CREATE TABLE IF NOT EXISTS `sys_role` (
    `id`          BIGINT       NOT NULL COMMENT '角色ID',
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

-- 部门表
CREATE TABLE IF NOT EXISTS `sys_dept` (
    `id`          BIGINT       NOT NULL COMMENT '部门ID',
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

-- 岗位表
CREATE TABLE IF NOT EXISTS `sys_post` (
    `id`          BIGINT       NOT NULL COMMENT '岗位ID',
    `post_code`   VARCHAR(50)  NOT NULL COMMENT '岗位编码',
    `post_name`   VARCHAR(50)  NOT NULL COMMENT '岗位名称',
    `sort`        INT          NOT NULL DEFAULT 0,
    `status`      TINYINT      NOT NULL DEFAULT 0,
    `remark`      VARCHAR(500) DEFAULT '',
    `creator`     BIGINT       DEFAULT NULL,
    `create_time` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updater`     BIGINT       DEFAULT NULL,
    `update_time` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted`     BIT(1)       NOT NULL DEFAULT b'0',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB COMMENT='岗位表';

-- 用户角色关联
CREATE TABLE IF NOT EXISTS `sys_user_role` (
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

-- 字典类型
CREATE TABLE IF NOT EXISTS `sys_dict_type` (
    `id`          BIGINT       NOT NULL,
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
    UNIQUE KEY `uk_dict_type` (`dict_type`)
) ENGINE=InnoDB COMMENT='字典类型表';

-- 字典数据
CREATE TABLE IF NOT EXISTS `sys_dict_data` (
    `id`          BIGINT       NOT NULL,
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

-- 初始化默认管理员（密码: admin123）
INSERT INTO `sys_user` (`id`, `username`, `password`, `nickname`, `status`)
  VALUES (1, 'admin', '$2a$10$VQBl2noKFPH/MSOsOq7a2.tdJqIckGw8MKTqPYxGqv3Rp7.q5mFJO', '超级管理员', 0);

INSERT INTO `sys_role` (`id`, `role_name`, `role_code`, `data_scope`, `status`)
  VALUES (1, '超级管理员', 'super_admin', 1, 0);

INSERT INTO `sys_user_role` (`id`, `user_id`, `role_id`) VALUES (1, 1, 1);

INSERT INTO `sys_dept` (`id`, `dept_name`, `sort`, `status`) VALUES (1, '总公司', 0, 0);

INSERT INTO `sys_post` (`id`, `post_code`, `post_name`, `sort`, `status`) VALUES (1, 'ceo', '董事长', 0, 0);
