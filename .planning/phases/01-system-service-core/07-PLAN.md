---
wave: 5
depends_on: ["05-PLAN.md"]
files_modified:
  - internal/storage/provider.go
  - internal/storage/local.go
  - internal/storage/minio.go
  - internal/storage/config.go
  - internal/module/system/model/file_entity.go
  - internal/module/system/repository/file.go
  - internal/module/system/service/file.go
  - internal/module/system/handler/file.go
  - internal/websocket/hub.go
  - internal/websocket/client.go
  - internal/websocket/handler.go
  - pkg/crypto/aes.go
  - pkg/crypto/rsa.go
  - pkg/gormx/hooks.go
  - internal/module/system/service/masking.go
  - internal/module/system/service/encryption.go
  - internal/module/system/handler/captcha.go
  - internal/module/system/model/client_entity.go
  - internal/module/system/repository/client.go
  - internal/module/system/service/client.go
  - internal/module/system/handler/client.go
  - configs/config.yaml
  - scripts/sql/tenant_template.sql
autonomous: true
requirements_addressed:
  - FR-07-01
  - FR-07-02
  - FR-07-03
  - FR-06-03
  - FR-08-01
  - FR-08-02
  - FR-01-09
  - FR-01-11
---

# Plan 07: 高级功能

## Objective

文件存储（MinIO + 动态切换）、WebSocket 在线用户（coder/websocket）、数据脱敏、AES/RSA 加解密、验证码登录、客户端管理。

## Tasks

### Task 7.1: 文件存储

<read_first>
- internal/config/config.go
</read_first>

<action>
StorageProvider 接口（Upload/Download/Delete/GetURL）。MinIO 实现（非默认凭据+SSL）。本地存储备选。动态切换后端（sys_param）。File entity + CRUD API。
</action>

<acceptance_criteria>
- StorageProvider 接口 + MinIO + Local 实现
- 非默认凭据 + SSL
- 动态切换存储后端
</acceptance_criteria>

---

### Task 7.2: 在线用户管理（WebSocket）

<read_first>
- go.mod
</read_first>

<action>
移除 gorilla/websocket，使用 coder/websocket 或 nhooyr.io/websocket。WebSocket hub（上线/下线/心跳）。在线用户 Redis 存储。强制踢出 API。
</action>

<acceptance_criteria>
- gorilla/websocket 已移除
- WebSocket 连接 + 心跳
- 在线用户列表 + 踢出 API
</acceptance_criteria>

---

### Task 7.3: 数据脱敏 + AES/RSA 加解密

<read_first>
- pkg/gormx/hooks.go
</read_first>

<action>
masking struct tag（phone/email/id_card/bank_card/name/address）。AES-256-GCM + RSA 工具。GORM Hook 透明加解密。
</action>

<acceptance_criteria>
- masking tag 支持 7 种类型
- AES-256-GCM + RSA 工具
- GORM BeforeCreate/AfterFind 透明加密
</acceptance_criteria>

---

### Task 7.4: 验证码登录 + 客户端管理

<read_first>
- internal/module/system/handler/auth.go
</read_first>

<action>
图形验证码生成（base64 + Redis 缓存 5min）。登录可选 captchaKey + captchaCode。SysClient entity + CRUD API。
</action>

<acceptance_criteria>
- POST /api/auth/captcha 返回 base64
- 登录支持验证码校验
- SysClient CRUD API
</acceptance_criteria>

---

## Verification

```bash
go build ./cmd/system/
go test ./internal/storage/... ./internal/websocket/... -v
```

## Must Haves

- [ ] MinIO 文件存储（非默认凭据 + SSL）
- [ ] WebSocket 在线用户（替换 gorilla/websocket）
- [ ] 数据脱敏 struct tag
- [ ] AES/RSA GORM Hook 透明加解密
- [ ] 图形验证码登录
- [ ] 客户端管理
