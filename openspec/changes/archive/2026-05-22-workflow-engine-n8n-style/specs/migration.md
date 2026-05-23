# Spec: PostgreSQL 迁移

## CONSTRAINT:PG-001 — 全局数据库切换
项目所有模块统一使用 PostgreSQL。移除 MySQL driver 依赖。

## CONSTRAINT:PG-002 — 配置格式
config.yaml 中 `mysql` 段替换为 `postgres`，包含 host/port/username/password/database/sslmode。

## CONSTRAINT:PG-003 — GORM Driver
使用 `gorm.io/driver/postgres` 替换 `gorm.io/driver/mysql`。

## CONSTRAINT:PG-004 — 类型映射
现有实体（BaseEntity、BaseModel 等）的 GORM tag 不变。建表 SQL 统一使用 PostgreSQL 语法：
- `BIGINT` 不变（主键雪花 ID）
- `TIMESTAMPTZ` 替代 `DATETIME`
- `SMALLINT` 替代 `TINYINT`
- `BOOLEAN` 或 `SMALLINT` 替代 `BIT(1)`

## CONSTRAINT:PG-005 — 软删除
`deleted` 字段使用 `SMALLINT DEFAULT 0`（0=未删除，1=已删除），配合 GORM scope。不使用 `gorm.DeletedAt` 的自动行为，改为手动 scope。

## CONSTRAINT:PG-006 — 建表脚本
所有建表脚本放在 `scripts/sql/` 下，文件名带 `pg_` 前缀。使用 `IF NOT EXISTS`。

## CONSTRAINT:PG-007 — 连接配置
DSN 格式：`host=%s port=%d user=%s password=%s dbname=%s sslmode=%s`
