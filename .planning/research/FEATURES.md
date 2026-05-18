# Features Research

## Competitive Analysis

### go-admin (github.com/go-admin-team/go-admin)
- Gin + GORM + Casbin, similar stack
- Monolithic, no multi-tenant
- Code generation focused
- Weak workflow support

### gin-vue-admin (github.com/flipped-aurora/gin-vue-admin)
- Gin + GORM + Casbin
- Monolithic, no multi-tenant
- Good documentation, active community
- No microservices path

### Our Differentiation
- **Microservices from day one** — unique in Go admin space
- **Independent-database multi-tenancy** — strongest isolation
- **External workflow (n8n)** — no Go native alternative
- **Hybrid REST + gRPC** — frontend compatibility + internal performance

## Feature Classification

### Table Stakes (Must Have — Phase 1)

| Feature | Priority | Notes |
|---------|----------|-------|
| User CRUD | P0 | With dept/role/position assignment |
| Role CRUD | P0 | Menu permission + data scope |
| Menu Management | P0 | Directory/Menu/Button 3-level |
| Dept Management | P0 | Tree structure |
| JWT Auth | P0 | Login/Logout/Refresh |
| Casbin RBAC | P0 | Domain = tenant |
| Multi-tenant | P0 | Independent database per tenant |
| Dictionary | P0 | Type + data |
| Parameters | P1 | System dynamic config |
| Notifications | P1 | Announcements |
| Operation Log | P0 | Audit trail |
| Login Log | P1 | Login records |
| Online Users | P1 | WebSocket monitoring + kick |
| File Management | P1 | Upload/download, MinIO |
| Data Permission | P0 | 5 scope levels |
| Tenant Package | P0 | Tenant feature gating |

### Differentiators (Competitive Advantage)

| Feature | Priority | Notes |
|---------|----------|-------|
| Microservices architecture | P0 | Unique in Go admin space |
| gRPC inter-service | P0 | Performance for service mesh |
| Saga distributed transactions | P1 | Cross-service consistency |
| Consul service discovery | P0 | Auto-registration, health check |
| Traefik gateway | P0 | Config-driven routing |
| Data encryption (AES/RSA/SM2/SM4) | P1 | Compliance requirement |
| Data masking | P1 | Privacy compliance |
| API transport encryption | P1 | Dynamic AES + RSA |
| Social login (WeChat/DingTalk) | P2 | OAuth2 integration |
| Client management | P1 | Multi-platform (PC/mini-program) |

### Anti-Features (Explicitly Out of Scope — Phase 1)

| Feature | Reason |
|---------|--------|
| Code generator | Microservices scaffold generation is complex; phase 2 |
| CRM | RuoYi-Vue-Pro specific; phase 2 on demand |
| ERP | RuoYi-Vue-Pro specific; phase 2 on demand |
| Mall (shop) | RuoYi-Vue-Pro specific; phase 2 on demand |
| IoT | RuoYi-Vue-Pro specific; phase 2 on demand |
| MES | RuoYi-Vue-Pro specific; phase 2 on demand |
| AI integration | RuoYi-Vue-Pro specific; phase 2 on demand |
| Payment (WeChat/Alipay) | Phase 2 |
| SMS/Email fusion | Phase 2 |
| WeChat public account | Phase 2 |
| i18n | Phase 2 |

## Feature Dependencies Graph

```
JWT Auth ──→ Casbin RBAC ──→ Data Permission
    |              |
    +--> Multi-tenant --> Tenant Package
              |
              +--> User CRUD --> Dept/Role/Position
                                    |
                                    +--> Menu Management

File Storage (independent)
Dictionary/Params (independent)
Logging (depends on Auth for user context)
Online Users (depends on WebSocket + Auth)
Workflow (depends on n8n external service)
```

## API Compatibility Note

Frontend (Vue3 + Vben5) expects specific API patterns from RuoYi:
- `GET /system/user/page` with `page`, `pageSize` query params
- Response envelope: `{code: 0, data: {list: [], total: 0}, msg: ""}`
- JWT in `Authorization: Bearer <token>` header
- Tenant ID in `tenant-id` header

Go backend must match these conventions for frontend compatibility.
