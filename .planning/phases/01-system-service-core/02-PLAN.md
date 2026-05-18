---
wave: 1
depends_on: ["01-PLAN.md"]
files_modified:
  - internal/auth/jwt.go
  - internal/auth/blacklist.go
  - internal/auth/login_lock.go
  - internal/middleware/auth.go
  - internal/middleware/ratelimit.go
  - internal/middleware/password.go
  - pkg/middleware/cors.go
  - pkg/middleware/security_headers.go
  - internal/module/system/model/entity.go
  - internal/module/system/model/dto.go
  - internal/module/system/service/auth.go
  - internal/module/system/handler/auth.go
  - internal/module/system/router.go
  - scripts/sql/init.sql
autonomous: true
requirements_addressed:
  - FR-01-01
  - FR-01-02
  - FR-01-06
  - FR-01-07
  - FR-01-08
  - NFR-05-02
  - NFR-05-03
  - NFR-05-06
---

# Plan 02: 安全基础 — JWT 双 Token + 黑名单 + 限流 + 密码策略

## Objective

实现 JWT 双 Token（Access 30min + Refresh 7d）、Redis 黑名单、登录速率限制（IP + 账号）、密码复杂度校验、登录失败锁定、CORS 白名单、安全响应头中间件。

## Context

当前 auth.go 仅做简单 JWT 解析，无黑名单、无速率限制、无密码策略。CORS 使用 `cors.Default()`（全开）。这是对抗审查中的安全 P0 修复项（CVSS 8.1-9.8）。

## Tasks

### Task 2.1: JWT 双 Token + Redis 黑名单

<read_first>
- internal/middleware/auth.go (当前 JWT 实现)
- internal/config/config.go (JWT 配置)
</read_first>

<action>
1. 创建 internal/auth/jwt.go：
   - TokenPair struct：AccessToken string, RefreshToken string, AccessExpire int64
   - GenerateTokenPair(userID, username, tenantID int64) (*TokenPair, error)：Access 30min + jti, Refresh 7d + jti
   - ParseAccessToken/ParseRefreshToken 解析并验证
   - Claims struct 新增 JTI string 字段

2. 创建 internal/auth/blacklist.go（依赖 *redis.Client）：
   - AddJTI(ctx, jti, ttl) — Redis SET jti "1" EX ttl
   - IsBlacklisted(ctx, jti) — Redis EXISTS
   - RevokeByUser(ctx, userID) — Redis SET user:{id}:revoked "{ts}" EX 7d
   - IsUserRevoked(ctx, userID, issuedAt) — 比较签发时间与吊销时间

3. 重构 internal/middleware/auth.go：使用 internal/auth 包，检查 jti 黑名单 + 用户级吊销，使用 context.WithValue 传播 userID/tenantID。

4. 创建 auth handler：Login/RefreshToken/Logout 三个端点。
</action>

<acceptance_criteria>
- internal/auth/jwt.go 导出 GenerateTokenPair, ParseAccessToken, ParseRefreshToken
- Claims struct 含 JTI string 字段
- internal/auth/blacklist.go 导出 AddJTI, IsBlacklisted, RevokeByUser, IsUserRevoked
- Auth() 中间件检查 jti 黑名单 + 用户级吊销
- POST /api/auth/login 返回 {accessToken, refreshToken}
- POST /api/auth/logout 将 jti 加入 Redis 黑名单
</acceptance_criteria>

---

### Task 2.2: 登录速率限制 + 密码策略 + 登录锁定

<read_first>
- internal/config/config.go (需要 LoginSecurity 配置)
</read_first>

<action>
1. Config struct 添加 LoginSecurityConfig（MaxFailCount=5, LockDuration=30min, RateLimitPerIP=20, RateLimitPerAcct=5, PasswordMinLen=8）。

2. 创建 internal/middleware/ratelimit.go：IPRateLimit(limit) 和 AccountRateLimit(limit) 使用 Redis INCR + EXPIRE。

3. 创建 internal/middleware/password.go：ValidatePassword 校验长度+大写+小写+数字+特殊字符；HashPassword/CheckPassword 使用 bcrypt。

4. 创建 internal/auth/login_lock.go：RecordFail/IsLocked/ResetFail/GetFailCount 使用 Redis 计数。
</action>

<acceptance_criteria>
- LoginSecurityConfig struct 含 5 个配置字段
- ratelimit.go 导出 IPRateLimit, AccountRateLimit
- password.go 导出 ValidatePassword, HashPassword, CheckPassword
- ValidatePassword 校验 8位+大写+小写+数字+特殊字符
- login_lock.go 导出 RecordFail, IsLocked, ResetFail
</acceptance_criteria>

---

### Task 2.3: CORS 白名单 + 安全响应头

<read_first>
- internal/router/router.go (当前 cors.Default())
</read_first>

<action>
1. pkg/middleware/cors.go：CORS(allowedOrigins []string)，使用 gin-contrib/cors，禁止 `*`，设置 AllowMethods/AllowHeaders/MaxAge。

2. pkg/middleware/security_headers.go：SecurityHeaders() 设置 X-Content-Type-Options, X-Frame-Options, X-XSS-Protection, HSTS, CSP。

3. 更新 router.go 替换 cors.Default()。
</action>

<acceptance_criteria>
- pkg/middleware/cors.go 导出 CORS 函数
- pkg/middleware/security_headers.go 导出 SecurityHeaders 函数
- SecurityHeaders 设置 5 个安全头
- router.go 不使用 cors.Default()
</acceptance_criteria>

---

### Task 2.4: 认证 API + User 表更新

<read_first>
- internal/module/system/model/entity.go (当前 User)
- scripts/sql/init.sql (当前建表)
- internal/module/system/router.go (当前路由)
</read_first>

<action>
1. User entity 添加：Password, LoginIP, LoginDate, LoginFailCount, LockTime。修改 deleted 为 deleted_at DATETIME。

2. 创建完整登录流程 handler/service。

3. 路由分两组：公开 /api/auth/（无 Auth 中间件）和认证 /api/system/（需 Auth）。

4. 更新 init.sql 使用 DATETIME 替代 BIT(1)。
</action>

<acceptance_criteria>
- User entity 含 Password, LoginIP, LoginDate, LoginFailCount, LockTime
- init.sql 使用 deleted_at DATETIME DEFAULT NULL
- POST /api/auth/login, /api/auth/refresh, /api/auth/logout 三个端点
- 登录失败超过限制后锁定
- go build ./cmd/system/ 通过
</acceptance_criteria>

---

## Verification

```bash
go build ./cmd/system/
go test ./internal/auth/... -v
go vet ./...
```

## Must Haves

- [ ] JWT 双 Token（Access 30min + Refresh 7d）
- [ ] Redis 黑名单（登出/禁用/改密吊销）
- [ ] 登录速率限制（IP + 账号，Redis INCR）
- [ ] 密码复杂度校验（8位+大小写+数字+特殊字符）
- [ ] 登录失败锁定
- [ ] CORS 白名单（生产禁止 `*`）
- [ ] 安全响应头
- [ ] deleted_at DATETIME 替代 BIT(1)
