---
name: management-backend
current_phase: 1
status: planned
updated: 2026-05-18
---

# Project State

## Current Phase

**Phase 1: System Service (Core)** — planned

## Phase Status

| Phase | Status | Plans | Updated |
|-------|--------|-------|---------|
| 1 | planned | 9 plans, 6 waves | 2026-05-18 |

## Decisions

- 微服务架构从第一天（不可协商）
- 独立数据库多租户（不可协商）
- 手动组合 Gin + gRPC
- godruNumin/gosnowflake + coder/websocket
- deleted_at DATETIME + Docker Compose
- Redis Streams + MySQL Saga 状态持久化

## Risks

1. GORM 连接池爆炸 — LRU + Prometheus
2. Casbin 策略同步 — Redis Pub/Sub
3. 跨租户泄漏 — tenantScope 注入
4. 前端兼容 — RuoYi 响应格式
