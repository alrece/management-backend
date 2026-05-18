# Architecture Research

## Service Decomposition

```
                    ┌─────────────┐
                    │   Traefik   │  API Gateway
                    │  (Gateway)  │  ForwardAuth, Rate Limit, TLS
                    └──────┬──────┘
                           │
              ┌────────────┼────────────┐
              │            │            │
        ┌─────▼─────┐ ┌───▼────┐ ┌────▼─────┐
        │  system    │ │  job   │ │ workflow  │
        │  service   │ │service │ │ service   │
        │  :8081     │ │:8082   │ │ :8083     │
        └─────┬──────┘ └───┬────┘ └────┬─────┘
              │            │           │
        ┌─────▼────────────▼───────────▼─────┐
        │            Shared Library            │
        │  (auth, tenant, response, snowflake) │
        └─────────────────┬───────────────────┘
                          │
              ┌───────────┼───────────┐
              │           │           │
         ┌────▼───┐  ┌───▼───┐  ┌───▼───┐
         │ MySQL  │  │ Redis │  │ MinIO │
         │(per-   │  │(cache,│  │(file  │
         │tenant) │  │ lock, │  │storage│
         └────────┘  │events)│  └───────┘
                     └───────┘
```

## Service Details

### system-service (Core — Phase 1)
- **Port**: 8081 (REST), 9081 (gRPC)
- **Owns**: User, Role, Menu, Dept, Dict, Param, Tenant, Log, Notice
- **Database**: `management_backend` (default) + per-tenant databases
- **gRPC API**: UserService, RoleService, MenuService (for other services)

### job-service (Phase 2)
- **Port**: 8082 (REST), 9082 (gRPC)
- **Owns**: Job definition, execution log, executor registry
- **Database**: `management_job`

### workflow-service (Phase 2)
- **Port**: 8083 (REST), 9083 (gRPC)
- **Owns**: Process category, form definition, process instance mapping
- **External**: n8n REST API integration
- **Database**: `management_workflow`

## Shared Library Design

```
pkg/
├── auth/          # JWT generation/parsing, Casbin adapter
├── tenant/        # DB resolver, tenant context, LRU connection pool
├── response/      # R{code, data, msg} unified response
├── errcode/       # Error codes by module segments
├── snowflake/     # Cross-service unique ID
├── middleware/     # Gin middleware (auth, tenant, logging, recovery)
├── redisx/        # Distributed lock, idempotency, event bus
├── gormx/         # Data permission scope, soft-delete scope
└── crypto/        # AES/RSA/SM2/SM4 encryption utilities
```

**Rule**: Shared library has ZERO business logic. Only infrastructure utilities.

## Database Architecture

### Multi-Tenant Strategy: Independent Database Per Tenant

```
management_backend        # Default/shared database
  ├── sys_tenant          # Tenant registry
  ├── sys_menu            # Shared menus (not tenant-filtered)
  ├── sys_role_menu       # Shared role-menu mapping
  └── sys_dict_type       # May be shared or per-tenant

tenant_1_db               # Tenant 1's isolated database
  ├── sys_user
  ├── sys_role
  ├── sys_dept
  ├── sys_user_role
  ├── sys_dict_type       # Per-tenant dictionaries
  ├── sys_dict_data
  └── sys_XXX_log

tenant_2_db               # Tenant 2's isolated database
  └── (same schema as tenant_1_db)
```

### Connection Pool Management
- Use GORM DBResolver for multi-database routing
- LRU cache for tenant DB connections (evict idle connections after 30 min)
- Max connections per tenant: configurable (default 10 idle, 50 max)

## Authentication Flow

```
Frontend                  Traefik                System Service
   │                        │                        │
   │  POST /auth/login      │                        │
   │───────────────────────>│───────────────────────>│
   │                        │                        │ Validate credentials
   │                        │                        │ Generate JWT
   │  {token, user}         │                        │
   │<───────────────────────│<───────────────────────│
   │                        │                        │
   │  GET /system/user/page │                        │
   │  Authorization: Bearer │                        │
   │───────────────────────>│                        │
   │                        │ ForwardAuth middleware  │
   │                        │ Verify JWT             │
   │                        │ Extract tenant_id      │
   │                        │───────────────────────>│
   │                        │                        │ Tenant scope applied
   │                        │                        │ RBAC check (Casbin)
   │                        │                        │ Data permission scope
   │  {code:0, data:{}}     │                        │
   │<───────────────────────│<───────────────────────│
```

## Data Permission (Corrected Architecture)

**Key Correction**: Data permission is NOT at gateway layer. Implemented as GORM Scope in service layer.

**Reason**: Traefik lacks business context (dept hierarchy, role data scope). Gateway only handles auth (JWT valid? tenant exists?). Data permission requires:
- User's department and ancestors
- Role's data_scope setting (1=all, 2=custom, 3=dept, 4=dept+children, 5=self)
- Custom dept list for scope=2

**Implementation**:
```go
func DataPermissionScope(userID, tenantID int64) func(*gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        scope := getUserDataScope(userID, tenantID)
        switch scope {
        case 1: // All data
            return db
        case 2: // Custom departments
            return db.Where("dept_id IN ?", getCustomDepts(userID))
        case 3: // Own department
            return db.Where("dept_id = ?", getUserDeptID(userID))
        case 4: // Department and children
            return db.Where("dept_id IN ?", getDeptAndChildren(getUserDeptID(userID)))
        case 5: // Self only
            return db.Where("creator = ?", userID)
        }
    }
}
```

## Event-Driven Pattern

### Redis Streams for Inter-Service Events

```
Publisher (system-service):
  redis.XAdd("user.created", {user_id, tenant_id, username})

Consumer (job-service):
  redis.XReadGroup("job-workers", "user.created")
```

### Saga Pattern (Limited Scope)

Only needed for cross-service operations (e.g., tenant creation):
1. system-service: Create tenant record → emit `tenant.created`
2. job-service: Initialize tenant job defaults → emit `tenant.initialized`
3. On failure: compensate via `tenant.creation.failed` event

Most CRUD operations are single-service — no Saga needed.

## Build Order

### Phase 1A: Foundation (system-service core)
1. Shared library (pkg/) — auth, response, snowflake, middleware, gormx
2. Config + Consul integration
3. Database connection + multi-tenant resolver
4. JWT auth + Casbin RBAC

### Phase 1B: System Module CRUD
5. User management (CRUD + dept/role assignment)
6. Role management (menu permission + data scope)
7. Menu management (tree structure)
8. Dept management (tree + ancestors)
9. Dictionary + Parameter management
10. Operation log + Login log

### Phase 1C: Advanced Features
11. Tenant management + package
12. Data permission enforcement
13. File management (MinIO)
14. Online user monitoring (WebSocket)
15. Notification/Announcement

### Phase 1D: Infrastructure Integration
16. Traefik gateway configuration
17. Consul service registration
18. gRPC service definitions
19. Docker Compose dev environment
20. OpenTelemetry tracing

## Key Architecture Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Monolith vs Micro | Micro from day 1 | Project positioning; refactor cost too high later |
| REST vs gRPC | Hybrid (REST external, gRPC internal) | Frontend compatibility + internal performance |
| Message bus | Redis Streams | Sufficient volume, one less infra component |
| Saga scope | Limited to cross-service ops | Most CRUD is single-service |
| Data permission | GORM Scope (service layer) | Gateway lacks business context |
| Multi-tenant | Independent database | Strongest isolation for SaaS |
| Config | Viper + Consul KV | Dynamic config without restart |
| DI | Constructor injection | Go idiomatic, no framework needed |
