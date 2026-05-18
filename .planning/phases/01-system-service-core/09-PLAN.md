---
wave: 6
depends_on: ["05-PLAN.md", "06-PLAN.md", "07-PLAN.md", "08-PLAN.md"]
files_modified:
  - testutil/testdb.go
  - testutil/mock_redis.go
  - testutil/fixtures.go
  - internal/module/system/repository/user_test.go
  - internal/module/system/repository/role_test.go
  - internal/module/system/repository/tenant_test.go
  - internal/auth/jwt_test.go
  - internal/auth/blacklist_test.go
  - internal/auth/login_lock_test.go
  - internal/authz/casbin_test.go
  - internal/module/system/service/auth_test.go
  - internal/module/system/service/user_test.go
  - internal/module/system/handler/auth_test.go
  - internal/module/system/handler/user_test.go
  - internal/multitenant/resolver_test.go
  - internal/multitenant/initializer_test.go
  - internal/middleware/tenant_test.go
  - internal/middleware/permission_test.go
  - internal/middleware/ratelimit_test.go
autonomous: true
requirements_addressed:
  - NFR-02-04
  - NFR-05-04
---

# Plan 09: 测试

## Objective

测试框架、Repository/Service/Handler 层测试、多租户隔离测试、安全测试。目标 80%+ 覆盖率。

## Tasks

### Task 9.1: 测试工具框架

<read_first>
- go.mod
- internal/module/system/repository/user.go
</read_first>

<action>
testutil/testdb.go（SQLite 内存库 + 自动迁移），testutil/mock_redis.go（miniredis），testutil/fixtures.go（测试数据生成器）。添加 testify + miniredis 依赖。
</action>

<acceptance_criteria>
- SQLite 内存库工具
- Redis mock（miniredis）
- 测试数据 fixtures
</acceptance_criteria>

---

### Task 9.2: Repository 层测试（80%+）

<read_first>
- internal/module/system/repository/user.go
- internal/module/system/repository/role.go
- internal/module/system/repository/tenant.go
</read_first>

<action>
user_test.go（CRUD + 分页 + tenantScope），role_test.go（CRUD + 菜单关联），tenant_test.go（CRUD + 状态过滤）。使用 SQLite 内存库。
</action>

<acceptance_criteria>
- 3 个 Repository 测试文件
- -cover >= 80%
</acceptance_criteria>

---

### Task 9.3: Service + Auth 层测试（80%+）

<read_first>
- internal/auth/jwt.go
- internal/auth/blacklist.go
</read_first>

<action>
auth_test.go（JWT 全流程 + 黑名单 + 刷新），user_test.go（业务逻辑）。Mock Repository + miniredis。
</action>

<acceptance_criteria>
- auth_test.go 覆盖 JWT 生成/解析/黑名单
- -cover >= 80%
</acceptance_criteria>

---

### Task 9.4: Handler 集成测试（httptest）

<read_first>
- internal/module/system/handler/auth.go
- internal/module/system/handler/user.go
</read_first>

<action>
auth_test.go（Login/Refresh/Logout HTTP 流程），user_test.go（CRUD + 认证/授权/跨租户）。httptest.NewRecorder + gin.TestMode。
</action>

<acceptance_criteria>
- 2 个 Handler 测试文件
- 覆盖认证/授权/错误场景
</acceptance_criteria>

---

### Task 9.5: 多租户隔离测试

<read_first>
- internal/multitenant/resolver.go
- internal/middleware/tenant.go
</read_first>

<action>
tenant_isolation_test.go：租户 A 数据租户 B 不可见，跨租户操作 403，切换连接正确。initializer_test.go：状态机流程 + 并发竞争 + 重试。
</action>

<acceptance_criteria>
- 跨租户查询返回空
- 跨租户操作 403
- 并发初始化不重复
</acceptance_criteria>

---

### Task 9.6: 安全测试

<read_first>
- internal/auth/blacklist.go
- internal/middleware/ratelimit.go
</read_first>

<action>
JWT 黑名单（登出/改密/Refresh 一次性），限流（429），CORS（拒绝非白名单），权限（403），密码策略（弱密码拒绝）。
</action>

<acceptance_criteria>
- 5 类安全测试场景
- go test -race 无竞态
</acceptance_criteria>

---

## Verification

```bash
go test ./... -cover -race -short
go test ./internal/... -coverprofile=coverage.out
go tool cover -func=coverage.out | tail -1
```

## Must Haves

- [ ] 测试框架（SQLite + miniredis + fixtures）
- [ ] Repository >= 80% 覆盖
- [ ] Service/Auth >= 80% 覆盖
- [ ] Handler httptest 集成测试
- [ ] 多租户隔离测试
- [ ] 安全测试（JWT/限流/CORS/权限/密码）
- [ ] go test -race 无竞态
