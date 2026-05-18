# Research Summary

## Overview

Four research streams + adversarial review (5 roles, 67 challenges) completed. Key findings synthesized below.

## Stack Decision: Self-Composition

- **Gin + gRPC** manual composition (no Kratos/go-zero/go-micro)
- Rationale: CRUD-heavy admin backend doesn't justify framework overhead
- Key corrections from adversarial review:
  - `godruNumin/gosnowflake` replaces `bwmarrin/snowflake` (clock drift protection)
  - `coder/websocket` replaces `gorilla/websocket` (archived project)
  - `gorm.DeletedAt` must use DATETIME (not BIT(1)) for GORM compatibility
- See [STACK.md](STACK.md) for full dependency list

## Architecture: Microservices with Shared Library

- 3 services: system (P1), job (P2), workflow (P2)
- Shared `pkg/` library with zero business logic
- Data permission at GORM Scope (service layer), NOT gateway layer
- Redis Streams for event bus; Saga state persisted to MySQL (not only Redis)
- **Key decisions unchanged**: microservices from day one, independent DB per tenant
- See [ARCHITECTURE.md](ARCHITECTURE.md) for diagrams and build order

## Adversarial Review Key Findings

- See [ADVERSARIAL-REVIEW.md](ADVERSARIAL-REVIEW.md) for full 67 challenges
- Core security fixes: JWT blacklist, rate limiting, CORS whitelist, password policy, tenant scope on all queries
- Core Go fixes: context.Context propagation (no *gin.Context in Service), unified error code mapping, snowflake clock drift
- Operations additions: Schema migration (goose), Prometheus+Grafana monitoring, Loki logging, backup strategy
- Timeline revised from 25d to ~40d (includes security, testing, observability, frontend integration)

## Feature Scope

- 20+ table-stakes features for Phase 1 (security additions from review)
- 10 differentiator features
- 11 anti-features deferred to Phase 2
- See [FEATURES.md](FEATURES.md) for classification

## Key Pitfalls (Top 5)

1. **GORM connection pool explosion** — LRU eviction + per-tenant metrics required
2. **Casbin policy sync** — Redis Adapter + Pub/Sub for multi-instance consistency
3. **Snowflake ID clock drift** — godruNumin/gosnowflake with built-in protection
4. **Cross-tenant data leakage** — ALL queries must inject tenantScope (GetByID/Delete were missing)
5. **Frontend API compatibility** — Must match RuoYi response format exactly
- See [PITFALLS.md](PITFALLS.md) for all 16 pitfalls

## Open Questions Resolved

1. Dict types: per-tenant (in tenant DB) — **resolved**
2. Social login: P2 framework only — **resolved**
3. Data encryption: AES/RSA first, SM2/SM4 follow — **resolved**
4. Tenant package: P2, simplified to `package_type` field on tenant table — **resolved**
5. WebSocket: `coder/websocket` (gorilla archived) — **resolved**
6. Production deployment: Docker Compose (explicit decision) — **resolved**
7. Schema migration: goose with multi-tenant batch support — **resolved**
8. Monitoring: Prometheus + Grafana in Phase 1 — **resolved**
