-- +goose Up
-- 租户库初始 Schema

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

CREATE TABLE IF NOT EXISTS `sys_role` (
    `id`          BIGINT       NOT NULL,
    `role_name`   VARCHAR(50)  NOT NULL,
    `role_code`   VARCHAR(50)  NOT NULL,
    `sort`        INT          NOT NULL DEFAULT 0,
    `data_scope`  TINYINT      NOT NULL DEFAULT 1,
    `status`      TINYINT      NOT NULL DEFAULT 0,
    `remark`      VARCHAR(500) DEFAULT '',
    `creator`     BIGINT       DEFAULT NULL,
    `create_time` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updater`     BIGINT       DEFAULT NULL,
    `update_time` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted`     BIT(1)       NOT NULL DEFAULT b'0',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB COMMENT='角色表';

CREATE TABLE IF NOT EXISTS `sys_dept` (
    `id`          BIGINT       NOT NULL,
    `parent_id`   BIGINT       NOT NULL DEFAULT 0,
    `ancestors`   VARCHAR(200) DEFAULT '',
    `dept_name`   VARCHAR(50)  NOT NULL,
    `sort`        INT          NOT NULL DEFAULT 0,
    `leader`      VARCHAR(30)  DEFAULT '',
    `status`      TINYINT      NOT NULL DEFAULT 0,
    `creator`     BIGINT       DEFAULT NULL,
    `create_time` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updater`     BIGINT       DEFAULT NULL,
    `update_time` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted`     BIT(1)       NOT NULL DEFAULT b'0',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB COMMENT='部门表';

CREATE TABLE IF NOT EXISTS `sys_user_role` (
    `id`          BIGINT   NOT NULL,
    `user_id`     BIGINT   NOT NULL,
    `role_id`     BIGINT   NOT NULL,
    `creator`     BIGINT   DEFAULT NULL,
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updater`     BIGINT   DEFAULT NULL,
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted`     BIT(1)   NOT NULL DEFAULT b'0',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB COMMENT='用户角色关联表';

-- +goose Down
DROP TABLE IF EXISTS `sys_user_role`;
DROP TABLE IF EXISTS `sys_dept`;
DROP TABLE IF EXISTS `sys_role`;
DROP TABLE IF EXISTS `sys_user`;
