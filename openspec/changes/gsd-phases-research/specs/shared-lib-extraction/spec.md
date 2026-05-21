## ADDED Requirements

### Requirement: Authx Package
The system SHALL provide pkg/authx/ package with JWT token generation/parsing and blacklist management, accepting config via parameters instead of global variables.

#### Scenario: Generate token pair with explicit config
- **WHEN** service calls authx.GenerateTokenPair(claims, jwtConfig)
- **THEN** system generates Access + Refresh tokens using provided secret and TTL

#### Scenario: Parse access token
- **WHEN** service calls authx.ParseAccessToken(tokenString, jwtConfig)
- **THEN** system validates signature, expiry, and returns Claims

#### Scenario: Blacklist token JTI
- **WHEN** service calls authx.BlacklistJTI(ctx, jti, ttl, redisClient)
- **THEN** system stores JTI in Redis with specified TTL

#### Scenario: Revoke all tokens for user
- **WHEN** service calls authx.RevokeUserTokens(ctx, userID, timestamp, redisClient)
- **THEN** system stores revocation timestamp in Redis, all tokens issued before are invalid

### Requirement: Authzx Package
The system SHALL provide pkg/authzx/ package with Casbin RBAC enforcement, accepting Redis client via parameters.

#### Scenario: Enforce permission check
- **WHEN** service calls authzx.Enforce(ctx, tenantID, userID, resource, action)
- **THEN** system checks Casbin policy for the tenant-scoped enforcer, returns allow/deny

### Requirement: Tenantx Package
The system SHALL provide pkg/tenantx/ package with multi-tenant DB resolution, accepting config via parameters.

#### Scenario: Resolve tenant database
- **WHEN** service calls tenantx.Resolve(tenantID, mysqlConfig, tenantConfig)
- **THEN** system returns *gorm.DB for the tenant's independent database

### Requirement: Discoveryx Package
The system SHALL provide pkg/discoveryx/ package with Consul service registration.

#### Scenario: Register service with Consul
- **WHEN** service calls discoveryx.Register(ctx, consulConfig, serviceInfo)
- **THEN** system registers service with Consul, starts health check

### Requirement: Backward Compatibility
The system SHALL maintain backward compatibility for internal/ packages during migration.

#### Scenario: Internal auth wrapper works
- **WHEN** existing code calls internal/auth.GenerateTokenPair after extraction
- **THEN** system delegates to pkg/authx with config.C values, behavior unchanged
