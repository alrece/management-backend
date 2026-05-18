# Technology Stack Research

## Decision: Self-Composition (No Framework)

**Reject**: Kratos, go-zero, go-micro, go-kratos

**Rationale**:
- Management backend is CRUD-heavy, not high-QPS — framework overhead unjustified
- Java projects (RuoYi) are self-composed; Go version should mirror this simplicity
- Frameworks enforce gRPC-first or protobuf-first — we need REST-first (frontend compatibility)
- Learning curve for framework × team size = bad ROI for admin backend
- Gin + gRPC manual composition gives full control, minimal abstraction

**Verdict**: Gin (HTTP) + gRPC (inter-service) + manual service composition.

## Core Stack

| Layer | Technology | Version | Role |
|-------|-----------|---------|------|
| HTTP | Gin | 1.9+ | REST API for frontend |
| RPC | gRPC + protobuf | - | Inter-service communication |
| ORM | GORM | 2.x | MySQL access, multi-tenant scopes |
| Auth | JWT (golang-jwt) + Casbin | - | Authentication + RBAC |
| Cache | go-redis | 9.x | Distributed lock, idempotency, cache, event bus |
| Config | Viper | 1.x | YAML + Consul KV |
| Logging | Zap | 1.x | Structured logging |
| ID | bwmarrin/snowflake | - | Cross-service unique ID |
| Migration | GORM AutoMigrate + raw SQL | - | Schema management |
| Validation | go-playground/validator | 10.x | Request validation |
| API Docs | swaggo/swag | - | Swagger/OpenAPI generation |
| Excel | excelize | 2.x | Import/Export |
| WebSocket | gorilla/websocket | - | Online user, notifications |
| Object Storage | minio-go | 7.x | S3-compatible file storage |
| Task Scheduling | robfig/cron | 3.x | Job scheduling |
| Password | golang.org/x/crypto (bcrypt) | - | Password hashing |
| Tracing | OpenTelemetry + Jaeger | - | Distributed tracing |

## Infrastructure

| Component | Technology | Role |
|-----------|-----------|------|
| API Gateway | Traefik | Routing, ForwardAuth, rate limiting, TLS |
| Service Discovery | Consul | Registration, health check, KV config |
| Message Bus | Redis Streams | Inter-service events (not Kafka — too heavy for admin) |
| Search | - | Out of scope (phase 2) |
| Container | Docker + Docker Compose | Local dev; k8s optional for production |

## Why Not Kafka / RabbitMQ

- Admin backend event volume is low (user operations, not click streams)
- Redis Streams provides consumer groups, persistence, replay — sufficient for Saga events
- One less infrastructure component to operate
- Can migrate to Kafka in phase 2 if volume demands

## Dependency Injection

No DI framework (wire, fx). Use constructor injection:

```go
func NewUserService(repo UserRepository, cache Cache) UserService {
    return &userService{repo: repo, cache: cache}
}
```

**Why**: Go idiomatic, explicit, easy to test with mocks, no code generation.

## Build & Dev Tools

| Tool | Purpose |
|------|---------|
| Make | Build, test, lint targets |
| Air | Hot reload for dev |
| golangci-lint | Static analysis |
| go mod | Dependency management |
| protoc + protoc-gen-go | gRPC code generation |
