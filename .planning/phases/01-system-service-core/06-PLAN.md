---
wave: 4
depends_on: ["02-PLAN.md"]
files_modified:
  - api/proto/system/user.proto
  - api/proto/system/role.proto
  - api/proto/system/menu.proto
  - internal/grpc/server.go
  - internal/grpc/client.go
  - internal/grpc/interceptors.go
  - internal/discovery/consul.go
  - configs/traefik/traefik.yaml
  - configs/traefik/dynamic.yaml
  - internal/saga/orchestrator.go
  - internal/saga/store.go
  - internal/saga/definition.go
  - internal/saga/deadletter.go
  - internal/telemetry/trace.go
  - cmd/system/main.go
  - docker-compose.yml
  - docker-compose.prod.yml
  - Makefile
  - docs/swagger.go
autonomous: true
requirements_addressed:
  - NFR-01-01
  - NFR-01-02
  - NFR-01-03
  - NFR-01-04
  - NFR-01-05
  - NFR-04-01
  - NFR-04-02
  - NFR-04-03
  - NFR-04-04
---

# Plan 06: 微服务基础设施

## Objective

gRPC + Consul + Traefik + Saga + OpenTelemetry + 优雅关闭 + Docker Compose + Swagger。

## Tasks

### Task 6.1: gRPC protobuf

<read_first>
- internal/module/system/model/entity.go
</read_first>

<action>
api/proto/system/user.proto（GetUser, ListUsers, CreateUser, UpdateUser, DeleteUser），role.proto，menu.proto。Makefile proto 命令。
</action>

<acceptance_criteria>
- 3 个 proto 文件存在
- Makefile proto 命令可生成 Go 代码
</acceptance_criteria>

---

### Task 6.2: gRPC 服务端(:9081) + 客户端（重连策略）

<read_first>
- api/proto/system/user.proto
</read_first>

<action>
grpc/server.go 注册在 :9081。client.go 含连接池 + 指数退避重连。interceptors.go 含认证/日志/recovery 拦截器。
</action>

<acceptance_criteria>
- REST(:8081) + gRPC(:9081) 同时启动
- client 含重连策略
</acceptance_criteria>

---

### Task 6.3: Consul 服务注册

<read_first>
- internal/config/config.go
</read_first>

<action>
discovery/consul.go：HTTP 健康检查 /health，ConsulEnabled 开关控制。
</action>

<acceptance_criteria>
- Consul 注册 + 健康检查 + 开关
</acceptance_criteria>

---

### Task 6.4: Traefik 网关

<read_first>
- configs/config.yaml
</read_first>

<action>
configs/traefik/ 静态+动态配置。JWT 验证留在服务中间件。
</action>

<acceptance_criteria>
- Traefik 配置 Consul 服务发现
- JWT 在服务侧验证
</acceptance_criteria>

---

### Task 6.5: Saga 框架

<read_first>
- .planning/research/ARCHITECTURE.md
</read_first>

<action>
saga/orchestrator.go（Redis Streams），store.go（MySQL 状态持久化），deadletter.go（死信队列），definition.go（DSL）。租户创建 Saga 示例。
</action>

<acceptance_criteria>
- Saga 状态：PENDING/RUNNING/COMPENSATING/COMPLETED/FAILED
- MySQL 持久化 + Redis Streams 事件
- 死信队列
</acceptance_criteria>

---

### Task 6.6: OpenTelemetry + Jaeger

<read_first>
- internal/config/config.go
</read_first>

<action>
telemetry/trace.go：Jaeger exporter，采样策略（生产 1-10%），trace_id 注入 Zap。
</action>

<acceptance_criteria>
- trace_id 关联到 Zap 日志
- 采样策略可配置
</acceptance_criteria>

---

### Task 6.7: 优雅关闭（30s 超时）

<read_first>
- cmd/system/main.go
</read_first>

<action>
SIGTERM → 停 HTTP → 排空事务 → 关 gRPC → 关连接池 → 关 DB/Redis → 30s 强制退出。
</action>

<acceptance_criteria>
- 关闭顺序正确，30s 超时
</acceptance_criteria>

---

### Task 6.8: Docker Compose + Swagger

<read_first>
- configs/config.yaml
</read_first>

<action>
docker-compose.yml（MySQL/Redis/MinIO/Consul/Jaeger）。docker-compose.prod.yml（+Traefik/Prometheus/Grafana/Loki）。Swagger 生成。
</action>

<acceptance_criteria>
- docker-compose.yml 含 5 个服务
- docker-compose.prod.yml 含 Traefik + 监控栈
- make swagger 生成文档
</acceptance_criteria>

---

## Verification

```bash
go build ./cmd/system/
go test ./internal/saga/... -v
docker-compose config
```

## Must Haves

- [ ] gRPC protobuf + 服务端 + 客户端
- [ ] Consul 服务注册
- [ ] Traefik 网关配置
- [ ] Saga 框架（Redis Streams + MySQL + 死信队列）
- [ ] OpenTelemetry + Jaeger
- [ ] 优雅关闭（30s）
- [ ] Docker Compose 开发 + 生产
- [ ] Swagger 文档
