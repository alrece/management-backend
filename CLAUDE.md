# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 语言设置

**必须使用中文**与用户对话。代码注释使用中文。

## 项目概览

Go 管理后台后端，综合 `ruoyi-vue-pro/`（Spring Boot 2.7 + JDK 8）和 `ruoyi-vue-plus-ai-engineering/`（Spring Boot 3.5 + JDK 17）两个 Java 项目的设计规范，用 Go 技术栈全面重构。

| 维度 | 技术选型 |
|------|---------|
| 语言 | Go 1.22+ |
| Web 框架 | Gin |
| ORM | GORM + MySQL |
| 缓存 | go-redis/v9 |
| 认证 | golang-jwt + Casbin（RBAC） |
| 配置 | Viper（YAML + 环境变量 `MB_` 前缀覆盖） |
| 日志 | Zap |
| 校验 | go-playground/validator/v10 |
| API 文档 | swaggo/swag |
| 定时任务 | robfig/cron/v3 |
| 主键 | 雪花 ID（bwmarrin/snowflake） |
| 密码 | bcrypt |
| Excel | excelize |
| 对象存储 | minio-go（S3 协议） |

**Java 参考项目保留在子目录，不参与 Go 构建，已在 .air.toml 和 .gitignore 中排除。**

---

## 构建与运行

```bash
# 安装依赖（首次）
go mod tidy

# 开发模式（热重载，需安装 air）
make dev

# 直接运行
make run

# 编译
make build

# 测试
make test
make test-module MODULE=system
make test-one MODULE=system FUNC=TestCreateUser

# Swagger 文档
make swagger
```

### 环境依赖

- Go 1.22+
- MySQL 5.7+ / 8.0
- Redis 6.0+
- 初始化脚本：`scripts/sql/init.sql`

---

## 项目结构

```
management-backend/
├── cmd/server/main.go           # 入口：配置 → DB → 路由 → 启动
├── configs/config.yaml          # 主配置文件
├── internal/
│   ├── config/config.go         # 配置加载（Viper）
│   ├── middleware/               # 中间件（auth, tenant, recovery, cors）
│   ├── model/base.go            # 基础实体（BaseEntity, BaseModel）
│   ├── pkg/                     # 内部工具包
│   │   ├── response/            # 统一响应 R{code, data, msg}
│   │   ├── errcode/             # 业务错误码
│   │   ├── page/                # 分页请求基类
│   │   └── snowflake/           # 雪花 ID 生成
│   ├── module/                  # 业务模块（按模块隔离）
│   │   └── system/              # 系统管理模块
│   │       ├── model/
│   │       │   ├── entity.go    # 实体（对应 DB 表）
│   │       │   └── dto.go       # 请求/响应 DTO
│   │       ├── repository/      # 数据访问层（GORM 查询）
│   │       ├── service/         # 业务逻辑层
│   │       ├── handler/         # HTTP 处理层（Gin handler）
│   │       └── router.go        # 模块路由注册
│   └── router/router.go         # 全局路由 + 依赖注入
├── scripts/sql/                 # 数据库脚本
├── ruoyi-vue-pro/               # Java 参考（不参与构建）
├── ruoyi-vue-plus-ai-engineering/ # Java 参考（不参与构建）
├── go.mod
├── Makefile
└── .air.toml
```

---

## 核心架构（必须牢记）

### 三层架构

```
Handler（HTTP 层）→ Service（业务层）→ Repository（数据层）
      ↓                   ↓                    ↓
   参数校验             业务逻辑            GORM 查询
   绑定请求             错误码处理           事务管理
   统一响应             对象转换            条件构建
```

**无 DAO 层**，Repository 直接封装 GORM 操作，与两个 Java 项目保持一致。

### 依赖注入

Go 没有_DI 容器_，通过**构造函数注入**：

```go
// router/router.go 中手动组装依赖链
userRepo := repository.NewUserRepo(db)           // 注入 *gorm.DB
userSvc := service.NewUserService(userRepo)       // 注入 UserRepo 接口
userHandler := handler.NewUserHandler(userSvc)    // 注入 UserService 接口
```

### 标准模块结构

新建模块时，严格遵循：

```
internal/module/{模块名}/
├── model/
│   ├── entity.go          # 实体（嵌入 model.BaseEntity）
│   └── dto.go             # XxxCreateReq / XxxUpdateReq / XxxPageReq / XxxResp
├── repository/
│   └── xxx.go             # XxxRepo 接口 + userRepo 实现
├── service/
│   └── xxx.go             # XxxService 接口 + xxxService 实现
├── handler/
│   └── xxx.go             # XxxHandler struct + Gin handler 方法
└── router.go              # RegisterRoutes() 注册路由
```

---

## 命名规范

| 类型 | 规范 | 示例 |
|------|------|------|
| 实体 struct | 名词，无后缀 | `User`（不是 UserDO） |
| 创建请求 | `XxxCreateReq` | `UserCreateReq` |
| 更新请求 | `XxxUpdateReq` | `UserUpdateReq` |
| 分页请求 | `XxxPageReq`，嵌入 `page.Req` | `UserPageReq` |
| 响应 | `XxxResp` | `UserResp` |
| Repository 接口 | `XxxRepo` | `UserRepo` |
| Repository 实现 | 小写 `xxxRepo` | `userRepo`（unexported） |
| Service 接口 | `XxxService` | `UserService` |
| Service 实现 | 小写 `xxxService` | `userService`（unexported） |
| Handler | `XxxHandler` | `UserHandler` |
| 构造函数 | `NewXxxRepo/Service/Handler` | `NewUserRepo(db) UserRepo` |
| 错误码常量 | 大驼峰 | `UserNotFound = 1001` |
| 表名 | `sys_` / `模块_` 前缀 | `sys_user`, `sys_role` |

---

## API 路径规范

综合两个 Java 项目，统一为：

| 操作 | HTTP 方法 | 路径 | Handler 方法 |
|------|----------|------|-------------|
| 创建 | POST | `/api/{module}/{resource}` | `Create` |
| 更新 | PUT | `/api/{module}/{resource}` | `Update` |
| 删除 | DELETE | `/api/{module}/{resource}/:id` | `Delete` |
| 详情 | GET | `/api/{module}/{resource}/:id` | `Get` |
| 分页 | GET | `/api/{module}/{resource}/page` | `Page` |
| 导出 | POST | `/api/{module}/{resource}/export` | `Export` |

认证路由统一在 `/api` 路由组下，由 `middleware.Auth()` + `middleware.Tenant()` 中间件保护。

---

## 统一响应格式

```go
// 成功
response.Ok(c, data)                    // {"code": 0, "data": {...}, "msg": "success"}
response.OkMsg(c, "操作成功")            // {"code": 0, "msg": "操作成功"}
response.OkPage(c, list, total, p, ps)  // {"code": 0, "data": {"list": [...], "total": 100, ...}}

// 失败
response.Fail(c, errcode.UserNotFound, "用户不存在")  // {"code": 1001, "msg": "用户不存在"}
```

---

## 数据库设计规范

综合两个 Java 项目，建表模板：

```sql
CREATE TABLE `xxx_table` (
    `id`          BIGINT       NOT NULL COMMENT '主键ID（雪花ID）',
    `tenant_id`   BIGINT       NOT NULL DEFAULT 0 COMMENT '租户编号',
    -- 业务字段
    `status`      TINYINT      NOT NULL DEFAULT 0 COMMENT '状态(0正常 1停用)',
    `remark`      VARCHAR(500) DEFAULT '' COMMENT '备注',
    -- 审计字段
    `creator`     BIGINT       DEFAULT NULL COMMENT '创建者',
    `create_time` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updater`     BIGINT       DEFAULT NULL COMMENT '更新者',
    `update_time` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted`     BIT(1)       NOT NULL DEFAULT b'0' COMMENT '是否删除',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB COMMENT='xxx表';
```

**关键规范**：
- 主键用雪花 ID（`BIGINT`），**禁止 AUTO_INCREMENT**
- 必须包含 `tenant_id`（多租户）、审计字段、软删除 `deleted`
- 表名用 `sys_`（系统）、`模块_`（业务）前缀
- 时间字段统一 `DATETIME`，Go 对应 `time.Time`
- 受租户过滤的表包含 `tenant_id`，`sys_menu`/`sys_tenant`/`sys_role_menu` 等排除

---

## 绝对禁止的写法

| 禁止 | 正确 | 原因 |
|------|------|------|
| `db.Raw("SELECT * FROM user WHERE id=" + id)` | `db.Where("id = ?", id)` | SQL 注入 |
| `AUTO_INCREMENT` | 雪花 ID `snowflake.NextID()` | 两个 Java 项目均用雪花 ID |
| `time.Date(...)` / `time.Now().Unix()` 直接序列化 | `time.Time` + GORM auto tag | 统一时间处理 |
| `interface{}` 传递业务数据 | 明确的 `XxxResp` struct | 类型安全 |
| `panic()` 处理业务错误 | `return errcode.Err(errcode.XXX)` | Go 错误处理惯例 |
| 在 Service 中写 GORM 查询 | 在 Repository 层构建查询 | 职责分离 |
| `map[string]interface{}` 做 API 响应 | `response.Ok(c, typedData)` | 类型安全 |
| 全局变量存 DB 连接 | 构造函数注入 `*gorm.DB` | 依赖注入，便于测试 |
| `init()` 函数做复杂初始化 | `main.go` 显式调用 | 可控的生命周期 |
| 忽略 `context.Context` | 所有 Repository 方法首参数 `ctx` | 超时/取消传播 |

---

## 对象转换

Go 没有 MapStruct，使用**手动映射函数**（保持清晰可追踪）：

```go
// 在 Service 层定义 Entity → Resp 转换
func toResp(u *model.User) *model.UserResp {
    return &model.UserResp{
        ID:       u.ID,
        Username: u.Username,
        // ...
    }
}
```

**禁止**使用反射库做通用拷贝（性能差且隐藏 bug）。

---

## 多租户实现

通过 GORM Scope 实现行级隔离，对业务层透明：

```go
// middleware/tenant.go 将租户过滤注入 Gin Context
c.Set("tenantScope", func(db *gorm.DB) *gorm.DB {
    return db.Where("tenant_id = ?", tenantID)
})

// repository 层使用
if tenantScope := middleware.GetTenantScope(c); tenantScope != nil {
    db = tenantScope(db)
}
```

---

## 错误码分段

| 范围 | 模块 |
|------|------|
| 0 | 成功 |
| 400-499 | HTTP 标准错误 |
| 500 | 服务器内部错误 |
| 1000-1999 | 系统模块（用户/角色/菜单/部门/Token） |
| 2000-2999 | 基础设施（文件/配置/字典） |
| 3000-3999 | 租户 |

---

## Handler 代码模板

```go
func (h *UserHandler) Create(c *gin.Context) {
    var req model.UserCreateReq
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Fail(c, 400, "参数校验失败: "+err.Error())
        return
    }
    id, err := h.svc.Create(c.Request.Context(), &req,
        middleware.GetUserID(c), middleware.GetTenantID(c))
    if err != nil {
        response.Fail(c, 500, err.Error())
        return
    }
    response.Ok(c, id)
}
```

---

## Service 代码模板

```go
func (s *userService) Create(ctx context.Context, req *model.UserCreateReq, creator int64, tenantID int64) (int64, error) {
    // 1. 业务校验
    existing, _ := s.repo.GetByUsername(ctx, req.Username)
    if existing != nil {
        return 0, errcode.Err(errcode.UserExists)
    }
    // 2. 构建实体
    user := &model.User{...}
    user.ID = snowflake.NextID()
    // 3. 持久化
    if err := s.repo.Create(ctx, user); err != nil {
        return 0, err
    }
    return user.ID, nil
}
```

---

## Repository 代码模板

```go
func (r *userRepo) Page(ctx context.Context, req *model.UserPageReq, tenantScope func(*gorm.DB) *gorm.DB) ([]model.User, int64, error) {
    var list []model.User
    var total int64
    db := r.db.WithContext(ctx).Model(&model.User{})

    // 租户过滤
    if tenantScope != nil {
        db = tenantScope(db)
    }
    // 条件构建（对应 Java LambdaQueryWrapperX）
    if req.Username != "" {
        db = db.Where("username LIKE ?", "%"+req.Username+"%")
    }
    // 计数 + 分页
    db.Count(&total)
    db.Offset(req.Offset()).Limit(req.Limit()).Order("id DESC").Find(&list)
    return list, total, nil
}
```

---

## 参考代码位置

| 开发类型 | 参考文件 |
|---------|---------|
| 完整三层示例 | `internal/module/system/` |
| 实体定义 | `internal/module/system/model/entity.go` |
| DTO 定义 | `internal/module/system/model/dto.go` |
| Repository | `internal/module/system/repository/user.go` |
| Service | `internal/module/system/service/user.go` |
| Handler | `internal/module/system/handler/user.go` |
| 路由注册 | `internal/module/system/router.go` |
| 全局路由+DI | `internal/router/router.go` |
| 中间件 | `internal/middleware/` |
| 基础模型 | `internal/model/base.go` |

---

## 开发前检查清单

- [ ] 已读参考模块 `internal/module/system/` 的代码
- [ ] Entity 嵌入 `model.BaseEntity`（含租户的表）或 `model.BaseModel`（排除表）
- [ ] 主键使用 `snowflake.NextID()`，**禁止 AUTO_INCREMENT**
- [ ] Repository 定义接口 + unexported 实现
- [ ] Service 定义接口 + unexported 实现
- [ ] 所有 Repository 方法首参数为 `context.Context`
- [ ] 查询构建在 Repository 层（不在 Service）
- [ ] 统一响应用 `response.Ok()` / `response.Fail()`
- [ ] 业务错误用 `errcode.Err(errcode.XXX)`
- [ ] Handler 参数校验用 `binding` tag
- [ ] 路由注册在模块的 `router.go` + 全局 `router/router.go`
- [ ] 依赖注入在 `router.Setup()` 中手动组装
