-- 工作流引擎建表脚本 (PostgreSQL)
-- 依赖 pg_init.sql 中已创建的基础环境

-- ============================================================
-- 1. 工作流定义表
-- ============================================================
CREATE TABLE IF NOT EXISTS wf_workflow (
    id              BIGINT          NOT NULL COMMENT '主键ID（雪花ID）',
    tenant_id       BIGINT          NOT NULL DEFAULT 0 COMMENT '租户编号',
    name            VARCHAR(100)    NOT NULL COMMENT '工作流名称',
    category_id     BIGINT          DEFAULT 0 COMMENT '分类编号',
    definition      JSONB           COMMENT '工作流定义（节点+边+设置），最大1MB',
    trigger_type    VARCHAR(20)     DEFAULT 'manual' COMMENT '触发器类型(manual/cron/webhook)',
    trigger_config  JSONB           COMMENT '触发器配置(cron表达式等)',
    version         INT             DEFAULT 1 COMMENT '版本号',
    status          SMALLINT        NOT NULL DEFAULT 0 COMMENT '状态(0正常 1停用)',
    remark          VARCHAR(500)    DEFAULT '' COMMENT '备注',
    creator         BIGINT          DEFAULT NULL COMMENT '创建者',
    create_time     TIMESTAMPTZ     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updater         BIGINT          DEFAULT NULL COMMENT '更新者',
    update_time     TIMESTAMPTZ     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted         SMALLINT        NOT NULL DEFAULT 0 COMMENT '是否删除(0否 1是)',
    create_dept     BIGINT          DEFAULT 0 COMMENT '创建部门',
    PRIMARY KEY (id)
);

CREATE INDEX idx_wf_workflow_tenant ON wf_workflow (tenant_id);
CREATE INDEX idx_wf_workflow_category ON wf_workflow (category_id);
CREATE INDEX idx_wf_workflow_deleted ON wf_workflow (deleted);

COMMENT ON TABLE wf_workflow IS '工作流定义表';

-- ============================================================
-- 2. 工作流执行实例表
-- ============================================================
CREATE TABLE IF NOT EXISTS wf_execution (
    id              BIGINT          NOT NULL COMMENT '主键ID（雪花ID）',
    tenant_id       BIGINT          NOT NULL DEFAULT 0 COMMENT '租户编号',
    workflow_id     BIGINT          NOT NULL COMMENT '工作流编号',
    status          VARCHAR(20)     NOT NULL DEFAULT 'running' COMMENT '执行状态(running/success/failed/canceled)',
    trigger_type    VARCHAR(20)     COMMENT '触发类型',
    start_time      TIMESTAMPTZ     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '开始时间',
    end_time        TIMESTAMPTZ     COMMENT '结束时间',
    input           JSONB           COMMENT '输入参数',
    output          JSONB           COMMENT '执行输出',
    error_msg       TEXT            COMMENT '错误信息',
    creator         BIGINT          DEFAULT NULL COMMENT '触发者',
    PRIMARY KEY (id)
);

CREATE INDEX idx_wf_execution_tenant ON wf_execution (tenant_id);
CREATE INDEX idx_wf_execution_workflow ON wf_execution (workflow_id);
CREATE INDEX idx_wf_execution_status ON wf_execution (status);

COMMENT ON TABLE wf_execution IS '工作流执行实例表';

-- ============================================================
-- 3. 节点执行日志表
-- ============================================================
CREATE TABLE IF NOT EXISTS wf_node_log (
    id              BIGINT          NOT NULL COMMENT '主键ID（雪花ID）',
    tenant_id       BIGINT          NOT NULL DEFAULT 0 COMMENT '租户编号',
    execution_id    BIGINT          NOT NULL COMMENT '执行实例编号',
    workflow_id     BIGINT          NOT NULL COMMENT '工作流编号',
    node_id         VARCHAR(50)     NOT NULL COMMENT '节点ID',
    node_type       VARCHAR(30)     NOT NULL COMMENT '节点类型',
    status          VARCHAR(20)     NOT NULL DEFAULT 'running' COMMENT '节点状态(running/success/failed/skipped)',
    input           JSONB           COMMENT '节点输入',
    output          JSONB           COMMENT '节点输出',
    error_msg       TEXT            COMMENT '错误信息',
    start_time      TIMESTAMPTZ     NOT NULL COMMENT '开始时间',
    end_time        TIMESTAMPTZ     COMMENT '结束时间',
    retry_count     INT             DEFAULT 0 COMMENT '重试次数',
    PRIMARY KEY (id)
);

CREATE INDEX idx_wf_node_log_execution ON wf_node_log (execution_id);
CREATE INDEX idx_wf_node_log_workflow ON wf_node_log (workflow_id);

COMMENT ON TABLE wf_node_log IS '节点执行日志表';

-- ============================================================
-- 4. 触发器配置表
-- ============================================================
CREATE TABLE IF NOT EXISTS wf_trigger (
    id              BIGINT          NOT NULL COMMENT '主键ID（雪花ID）',
    tenant_id       BIGINT          NOT NULL DEFAULT 0 COMMENT '租户编号',
    workflow_id     BIGINT          NOT NULL COMMENT '工作流编号',
    trigger_type    VARCHAR(20)     NOT NULL COMMENT '触发器类型(cron/webhook)',
    config          JSONB           COMMENT '触发器配置',
    cron_entry_id   INT             DEFAULT 0 COMMENT 'Cron Entry ID',
    active          BOOLEAN         DEFAULT FALSE COMMENT '是否激活',
    creator         BIGINT          DEFAULT NULL COMMENT '创建者',
    create_time     TIMESTAMPTZ     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updater         BIGINT          DEFAULT NULL COMMENT '更新者',
    update_time     TIMESTAMPTZ     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted         SMALLINT        NOT NULL DEFAULT 0 COMMENT '是否删除',
    create_dept     BIGINT          DEFAULT 0 COMMENT '创建部门',
    PRIMARY KEY (id)
);

CREATE INDEX idx_wf_trigger_workflow ON wf_trigger (workflow_id);
CREATE INDEX idx_wf_trigger_tenant ON wf_trigger (tenant_id);

COMMENT ON TABLE wf_trigger IS '触发器配置表';

-- ============================================================
-- 5. 工作流分类表
-- ============================================================
CREATE TABLE IF NOT EXISTS wf_category (
    id              BIGINT          NOT NULL COMMENT '主键ID（雪花ID）',
    tenant_id       BIGINT          NOT NULL DEFAULT 0 COMMENT '租户编号',
    name            VARCHAR(100)    NOT NULL COMMENT '分类名称',
    parent_id       BIGINT          DEFAULT 0 COMMENT '父分类ID',
    sort            INT             DEFAULT 0 COMMENT '排序',
    status          SMALLINT        NOT NULL DEFAULT 0 COMMENT '状态(0正常 1停用)',
    creator         BIGINT          DEFAULT NULL COMMENT '创建者',
    create_time     TIMESTAMPTZ     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updater         BIGINT          DEFAULT NULL COMMENT '更新者',
    update_time     TIMESTAMPTZ     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted         SMALLINT        NOT NULL DEFAULT 0 COMMENT '是否删除',
    create_dept     BIGINT          DEFAULT 0 COMMENT '创建部门',
    PRIMARY KEY (id)
);

CREATE INDEX idx_wf_category_tenant ON wf_category (tenant_id);
CREATE INDEX idx_wf_category_parent ON wf_category (parent_id);

COMMENT ON TABLE wf_category IS '工作流分类表';
