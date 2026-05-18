---
wave: 1
depends_on: []
files_modified:
  - cmd/system/main.go
  - go.work
  - go.mod
  - pkg/go.mod
  - pkg/response/response.go
  - pkg/errcode/code.go
  - pkg/errcode/errcode.go
  - pkg/page/page.go
  - pkg/snowflake/snowflake.go
  - pkg/snowflake/id.go
  - pkg/middleware/request_id.go
  - pkg/middleware/recovery.go
  - pkg/middleware/context.go
  - internal/config/config.go
  - Makefile
  - .air.toml
  - configs/config.yaml
autonomous: true
requirements_addressed:
  - NFR-01-06
  - NFR-02-01
  - NFR-02-02
  - NFR-02-03
  - NFR-02-05
  - NFR-05-01
  - NFR-05-05
---

# Plan 01: 项目结构重组 + 共享库 + 配置中心

## Objective

将现有单体结构重组为微服务布局（go.work），创建共享 pkg/ 库，实现配置中心（Viper + 环境变量），替换 snowflake 为 godruNumin/gosnowflake，统一错误码映射，建立 Request ID 中间件。

## Context

当前代码是单体结构（`internal/module/system/`），需要重构为微服务布局。共享库从 `internal/pkg/` 迁移到 `pkg/`（跨服务共享）。配置需要支持环境变量覆盖（`MB_` 前缀）和敏感配置启动校验。

**用户铁律决策：微服务架构从第一天开始，独立数据库多租户。**

## Tasks

### Task 1.1: 项目结构重组为微服务布局

<read_first>
- cmd/server/main.go (当前入口)
- internal/router/router.go (当前路由)
- go.mod (当前模块)
- Makefile (当前构建命令)
</read_first>

<action>
1. 创建 go.work 文件：
```
go 1.22

use (
    .
    pkg
    cmd/system
)
```

2. 重构为微服务目录结构：
```
management-backend/
├── go.work
├── go.work.sum
├── pkg/                          # 共享库（独立 module）
│   ├── go.mod                    # module management-backend/pkg
│   ├── response/
│   │   └── response.go
│   ├── errcode/
│   │   ├── code.go
│   │   └── errcode.go
│   ├── page/
│   │   └── page.go
│   ├── snowflake/
│   │   ├── snowflake.go          # IDGenerator 接口定义
│   │   └── id.go                 # godruNumin 实现
│   ├── middleware/
│   │   ├── request_id.go
│   │   ├── recovery.go
│   │   ├── cors.go
│   │   ├── security_headers.go
│   │   └── context.go            # context key 工具函数
│   ├── gormx/
│   │   ├── scope.go              # 租户 scope、数据权限 scope
│   │   └── hooks.go              # 通用 GORM hooks
│   ├── redisx/
│   │   ├── lock.go               # 分布式锁
│   │   └── cache.go              # 缓存工具
│   └── crypto/
│       ├── aes.go
│       └── rsa.go
├── cmd/
│   └── system/
│       └── main.go               # system-service 入口
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── model/
│   │   └── base.go
│   └── module/
│       └── system/
│           ├── model/
│           ├── repository/
│           ├── service/
│           ├── handler/
│           └── router.go
├── configs/
│   └── config.yaml
├── scripts/
│   └── sql/
└── go.mod                        # module management-backend
```

3. 创建 pkg/go.mod：
```go
module management-backend/pkg

go 1.22

require (
    github.com/gin-gonic/gin v1.10.0
    github.com/godruNumin/gosnowflake v0.4.1
    go.uber.org/zap v1.27.0
    gorm.io/gorm v1.25.12
)
```

4. 主 go.mod 添加 replace 指令：
```
require management-backend/pkg v0.0.0
replace management-backend/pkg => ./pkg
```

5. 创建 cmd/system/main.go 作为 system-service 入口（从 cmd/server/main.go 迁移）。

6. 删除 cmd/server/main.go。

7. 更新 Makefile 的 build/dev/run 指向 cmd/system/main.go。

8. 更新 .air.toml 的 cmd 指向 cmd/system/main.go。
</action>

<acceptance_criteria>
- `go.work` 文件存在且包含 `use (.)` 和 `./pkg`
- `pkg/go.mod` 存在且 module 名为 `management-backend/pkg`
- `cmd/system/main.go` 存在，`cmd/server/main.go` 不存在
- `go build ./cmd/system/` 编译通过
- Makefile build 命令输出到 `bin/system`
</acceptance_criteria>

---

### Task 1.2: 共享库迁移 — response, errcode, page, snowflake

<read_first>
- internal/pkg/response/response.go (当前响应工具)
- internal/pkg/errcode/code.go (当前错误码)
- internal/pkg/page/page.go (当前分页)
- internal/pkg/snowflake/snowflake.go (当前雪花ID)
</read_first>

<action>
1. 迁移 response/ 到 pkg/response/，保持 R{code, data, msg} 结构。

2. 迁移 errcode/ 到 pkg/errcode/，重构为统一错误码：
   - code.go 定义分段常量：0（成功）/ 400-499（HTTP）/ 500（内部）/ 1000-1999（系统）/ 2000-2999（基础设施）/ 3000-3999（租户）
   - errcode.go 实现 BusinessError struct（Code int, Message string），支持 errors.As

3. 迁移 page/ 到 pkg/page/。

4. 替换 snowflake：pkg/snowflake/snowflake.go 定义 IDGenerator 接口（NextID() (int64, error)），pkg/snowflake/id.go 使用 godruNumin/gosnowflake 实现。

5. go.mod 移除 bwmarrin/snowflake，添加 godruNumin/gosnowflake。

6. 删除 internal/pkg/ 目录。
</action>

<acceptance_criteria>
- `pkg/response/response.go` 导出 Ok, OkMsg, OkPage, Fail
- `pkg/errcode/errcode.go` 定义 BusinessError struct 含 Code, Message 字段
- `pkg/errcode/code.go` 包含分段错误码（0/400-404/500/1000+/2000+/3000+）
- `pkg/snowflake/snowflake.go` 定义 IDGenerator 接口
- `pkg/snowflake/id.go` 使用 github.com/godruNumin/go-snowflake
- go.sum 包含 godruNumin/gosnowflake，不包含 bwmarrin/snowflake
- internal/pkg/ 目录不存在
- `go build ./pkg/...` 通过
</acceptance_criteria>

---

### Task 1.3: 配置中心 — Viper + 环境变量 + 敏感校验

<read_first>
- internal/config/config.go (当前配置结构)
- configs/config.yaml (当前配置文件)
</read_first>

<action>
1. 重构 internal/config/config.go：
   - Config struct 包含：Server, MySQL, Redis, JWT, Log, Snowflake, CORS, MinIO, Traefik, Consul
   - 使用 `viper.SetEnvPrefix("MB")` + `viper.AutomaticEnv()`
   - 添加 Validate() 方法校验：JWT Secret 非空非默认、MySQL 非默认密码、Redis 密码非空
   - MySQLConfig 添加 DSN() 方法

2. 更新 configs/config.yaml：所有敏感值为空占位符，通过环境变量设置。包含完整配置结构。

3. main.go 调用 config.C.Validate()。
</action>

<acceptance_criteria>
- config.go 使用 viper.SetEnvPrefix("MB") 和 viper.AutomaticEnv()
- Config struct 包含 10 个配置段
- Validate() 校验 JWT Secret、MySQL 密码、Redis 密码
- config.yaml 不含硬编码密码
- MySQLConfig 有 DSN() string 方法
- MB_JWT_SECRET=xxx 环境变量能覆盖 yaml 配置
</acceptance_criteria>

---

### Task 1.4: MySQL + Redis 连接 + Zap 日志 + Request ID

<read_first>
- cmd/server/main.go (当前数据库连接方式)
- internal/middleware/recovery.go (当前恢复中间件)
</read_first>

<action>
1. cmd/system/main.go 中初始化 Zap 结构化日志（替换 log.Printf）。

2. 创建 pkg/middleware/request_id.go：生成/传递 X-Request-ID。

3. 创建 pkg/middleware/recovery.go：使用 Zap logger 而非默认 log。

4. MySQL 初始化：GORM 配置 deleted_at DATETIME，logger.Warn 级别，连接池参数。

5. Redis 初始化：go-redis/v9 NewClient。

6. main.go 启动流程：Load → Validate → Logger → DB → Redis → Snowflake → Router → Run。
</action>

<acceptance_criteria>
- main.go 使用 zap.Logger 替代 log.Printf
- pkg/middleware/request_id.go 导出 RequestID() gin.HandlerFunc
- pkg/middleware/recovery.go 使用 *zap.Logger 参数
- GORM 设置 DeletedAt 为 DATETIME 类型
- Redis 使用 go-redis/v9 NewClient
- go build ./cmd/system/ 通过
</acceptance_criteria>

---

### Task 1.5: 统一错误码映射 + context.Context 传播

<read_first>
- internal/module/system/handler/user.go (当前 handler)
- internal/module/system/service/user.go (当前 service)
- internal/module/system/repository/user.go (当前 repository)
</read_first>

<action>
1. 创建 pkg/middleware/context.go：WithUserID/GetUserID/WithTenantID/GetTenantID/WithRequestID/GetRequestID（使用 context.WithValue，不依赖 gin.Context）。

2. Handler 使用 errors.As 提取 BusinessError 映射到响应码。

3. Handler 使用 c.Request.Context() 传递给 Service（不传 *gin.Context）。

4. Service 方法签名的第一个参数为 ctx context.Context。

5. Repository 使用 db.WithContext(ctx) 传播上下文。
</action>

<acceptance_criteria>
- pkg/middleware/context.go 导出 6 个 context key 工具函数
- 所有 Handler 使用 c.Request.Context()
- 所有 Service 方法首参数为 ctx context.Context
- 所有 Repository 使用 r.db.WithContext(ctx)
- Handler 错误使用 errors.As(err, &bizErr)
- go build ./... 通过
</acceptance_criteria>

---

## Verification

```bash
go build ./cmd/system/
go build ./pkg/...
go vet ./...
```

## Must Haves

- [ ] go.work 支持微服务多 module
- [ ] pkg/ 独立 Go module 包含 response/errcode/page/snowflake
- [ ] godruNumin/gosnowflake 替换 bwmarrin/snowflake
- [ ] MB_ 前缀环境变量覆盖配置
- [ ] 启动时敏感配置校验
- [ ] Zap 替换 log.Printf
- [ ] Request ID 中间件
- [ ] context.Context 传播（不传 *gin.Context 到 Service）
- [ ] 统一错误码（errors.As + BusinessError）
