---
phase: 01-system-service-core
plan: 02
subsystem: auth
tags: [jwt, redis, bcrypt, cors, security, rate-limit, gin]

requires:
  - plan: 01
    provides: pkg/ 共享库、配置中心、中间件模式
provides:
  - JWT 双 Token 系统（Access 30min + Refresh 7d）
  - Redis JWT 黑名单（JTI 级 + 用户级撤销）
  - 登录速率限制（IP 20/min + 账号 5/min，Redis INCR + 429）
  - 密码复杂度校验（8+ 字符，大小写+数字+特殊字符）
  - 登录失败锁定（5 次连续失败，锁定 30 分钟）
  - CORS 中间件（可配置白名单）
  - 安全响应头中间件
  - 认证 API（POST /api/auth/login, refresh, logout）
affects: [03-PLAN, 04-PLAN, 05-PLAN, 07-PLAN, 08-PLAN]

tech-stack:
  added: [github.com/golang-jwt/jwt/v5, github.com/redis/go-redis/v9, golang.org/x/crypto/bcrypt]
  patterns: [dual-token auth, token blacklist, rate limiting with sliding window, security headers]

key-files:
  created:
    - internal/auth/jwt.go
    - internal/auth/blacklist.go
    - internal/auth/login_lock.go
    - internal/middleware/ratelimit.go
    - internal/middleware/password.go
    - pkg/middleware/cors.go
    - pkg/middleware/security_headers.go
    - internal/module/system/handler/auth.go
  modified:
    - internal/middleware/auth.go
    - internal/router/router.go
    - internal/module/system/service/user.go

key-decisions:
  - "JWT 双 Token: Access 30min + Refresh 7d，JTI 标识"
  - "Redis 黑名单: 支持 JTI 级和用户级撤销"
  - "限流策略: IP 20/min + 账号 5/min，Redis INCR 滑动窗口"
  - "密码策略: 8+ 字符，必须包含大小写+数字+特殊字符"
  - "登录锁定: 5 次连续失败后锁定 30 分钟"

patterns-established:
  - "JWT 双 Token 模式: Access + Refresh，Redis 黑名单撤销"
  - "Redis 计数器限流: INCR + EXPIRE 滑动窗口"
  - "userQuerier 接口模式: 避免 auth handler 循环依赖"
  - "context.Context 传播: handler 层提取 c.Request.Context()"

requirements-completed: [FR-01-01, FR-01-02, FR-01-06, FR-01-07, FR-01-08, NFR-05-02, NFR-05-03, NFR-05-06]

duration: 18min
completed: 2026-05-18
---

# Plan 02: 安全基础 — JWT 双 Token + 黑名单 + 限流 + 密码策略 Summary

**JWT 双 Token 认证系统（Access/Refresh）+ Redis 黑名单 + 登录限流 + 密码策略 + CORS 白名单 + 安全响应头**

## Performance

- **Duration:** ~18 min
- **Completed:** 2026-05-18
- **Tasks:** 4
- **Files modified:** 11

## Accomplishments
- JWT 双 Token 系统：Access 30min + Refresh 7d，JTI 标识每个 Token
- Redis JWT 黑名单：支持 JTI 级撤销（单个 Token）和用户级撤销（全部 Token）
- 登录速率限制：IP 20/min + 账号 5/min，Redis INCR 滑动窗口，429 响应
- 密码复杂度校验 + bcrypt 哈希 + 登录失败锁定（5 次/30 分钟）
- CORS 可配置白名单 + 安全响应头中间件
- 认证 API 端点：POST /api/auth/login, refresh, logout

## Task Commits

1. **Task 2.1-2.5: 全部安全功能** - `de9a671` (feat)

## Files Created/Modified
- `internal/auth/jwt.go` - JWT 双 Token 生成与解析
- `internal/auth/blacklist.go` - Redis JWT 黑名单
- `internal/auth/login_lock.go` - 登录失败锁定
- `internal/middleware/ratelimit.go` - IP/账号限流
- `internal/middleware/password.go` - 密码复杂度校验
- `pkg/middleware/cors.go` - CORS 白名单中间件
- `pkg/middleware/security_headers.go` - 安全响应头
- `internal/middleware/auth.go` - Auth 中间件重构（集成黑名单）
- `internal/router/router.go` - 路由重组（公开/认证路由组）
- `internal/module/system/handler/auth.go` - 认证 API Handler
- `internal/module/system/service/user.go` - UserService GetByUsername

## Decisions Made
- userQuerier 接口模式：auth handler 定义独立接口避免循环依赖
- context.Context 传播：handler 层提取 c.Request.Context() 传递到 Service
- 路由分为公开组（/api/auth/login, /api/auth/refresh）和认证组（需 JWT）

## Deviations from Plan

### Auto-fixed Issues

**1. context type in revoke functions**
- **Found during:** Task 2.1
- **Issue:** revoke 函数参数使用了 *gin.Context 而非 context.Context
- **Fix:** 改为 context.Context 保持层级隔离
- **Committed in:** de9a671

**2. Removed unused fmt import from rate limit middleware**
- **Found during:** Task 2.3
- **Fix:** 清理未使用的 import
- **Committed in:** de9a671

**3. Auth handler context propagation and config reference fixes**
- **Found during:** Task 2.5
- **Fix:** 修复 context 传播和移除未使用的 config import
- **Committed in:** de9a671

---
**Total deviations:** 3 auto-fixed
**Impact on plan:** 均为代码质量问题修正，无范围蔓延。

## Issues Encountered
- auth handler 循环依赖：通过 userQuerier 接口解决
- context 类型不匹配：revoke 函数参数需 context.Context

## Next Phase Readiness
- 安全基础就绪，Plan 03（多租户核心）可直接使用 JWT 中间件获取租户信息
- Tenant 中间件将依赖 JWT claims 中的 tenant_id 字段

---
*Phase: 01-system-service-core*
*Completed: 2026-05-18*
