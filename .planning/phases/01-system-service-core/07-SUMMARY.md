---
phase: 01-system-service-core
plan: 07
subsystem: storage, websocket, crypto, auth
tags: [minio, websocket, aes-gcm, rsa, masking, captcha, oauth-client]

requires:
  - phase: 01-05
    provides: system CRUD modules and config
provides:
  - File storage with MinIO/local backend and dynamic switching
  - WebSocket hub/client for online user management
  - AES-256-GCM and RSA-OAEP crypto utilities
  - GORM transparent encryption hooks via struct tags
  - Data masking service (phone/email/idCard/bankCard/name/address)
  - Captcha generation with Redis caching
  - OAuth2 client management CRUD
affects: [observability, testing]

tech-stack:
  added: [github.com/minio/minio-go/v7]
  patterns: [storage-provider-factory, masking-struct-tag, gorm-encryption-hooks, websocket-hub]

key-files:
  created:
    - internal/storage/config.go
    - internal/websocket/hub.go
    - internal/websocket/client.go
    - internal/websocket/handler.go
    - pkg/crypto/aes.go
    - pkg/crypto/rsa.go
    - pkg/gormx/hooks.go
    - internal/module/system/service/masking.go
    - internal/module/system/service/captcha.go
    - internal/module/system/model/file_entity.go
    - internal/module/system/model/client_entity.go
  modified:
    - internal/router/router.go
    - internal/module/system/router.go
    - scripts/sql/tenant_template.sql

key-decisions:
  - "Storage provider factory with runtime switching via sys_param"
  - "nhooyr.io/websocket used as already in go.mod"
  - "Masking via struct tags with reflection for zero-intrusion"
  - "GORM hooks for transparent field-level AES-256-GCM encryption"

requirements-completed: [FR-07-01, FR-07-02, FR-07-03, FR-06-03, FR-08-01, FR-08-02, FR-01-09, FR-01-11]

duration: 15min
completed: 2026-05-19
---

# Phase 01 Plan 07: Advanced Features Summary

**File storage with MinIO/local dynamic backend, WebSocket online users, AES/RSA crypto with GORM hooks, data masking, captcha login, and OAuth2 client management**

## Performance

- **Duration:** ~15 min
- **Tasks:** 4
- **Files modified:** 22

## Accomplishments
- Storage Provider factory with MinIO/local backends and runtime switching
- WebSocket Hub with register/unregister/heartbeat event loop and Redis online tracking
- AES-256-GCM + RSA-OAEP crypto with GORM BeforeCreate/AfterFind transparent encryption
- Data masking service supporting 6 types via struct tags
- Captcha service with base64 image generation and Redis 5min TTL
- OAuth2 client management full CRUD

## Task Commits

1. **Task 7.1: File storage** - `068e521` (feat)
2. **Task 7.2: WebSocket online users** - `5eabc22` (feat)
3. **Task 7.3: Crypto + masking** - `00e9c29` (feat)
4. **Task 7.4: Captcha + client management** - `9c6883e` (feat)

## Decisions Made
- Used nhooyr.io/websocket (already in go.mod) rather than coder/websocket
- Storage provider type stored in config.C.OSS.Type, dynamic switching via SwitchProvider()
- Masking uses reflect-based struct tag parsing

## Deviations from Plan

**1. [Rule 3 - Blocking] Fixed BaseEntity package reference** - Task 7.1
**2. [Rule 3 - Blocking] Fixed context key names (userId vs userID)** - Task 7.2
**3. [Rule 1 - Bug] Fixed embedded struct field access for Client** - Task 7.4

**Total deviations:** 3 auto-fixed (2 blocking, 1 bug)
**Impact:** Compilation corrections only, no scope change.

## Issues Encountered
None

## Next Phase Readiness
- Plan 08 (observability) and Plan 09 (testing) ready to proceed

---
*Phase: 01-system-service-core*
*Completed: 2026-05-19*
