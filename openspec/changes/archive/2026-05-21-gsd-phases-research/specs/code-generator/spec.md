## ADDED Requirements

### Requirement: Microservice Scaffolding Generation
The system SHALL generate complete microservice project scaffold from template.

#### Scenario: Generate new service
- **WHEN** user runs code-generator with service name and module path
- **THEN** system generates: cmd/main.go, internal/{config,model,repository,service,handler,router}, api/proto/{module}/{service}.proto, configs/config.yaml, Dockerfile, Makefile, go.mod

### Requirement: Three-Layer Code Generation
The system SHALL generate model/repository/service/handler code from entity definitions.

#### Scenario: Generate CRUD from entity
- **WHEN** user provides entity definition with fields and table name
- **THEN** system generates: entity.go, dto.go (Create/Update/Page/Resp), repository.go (interface+impl), service.go (interface+impl), handler.go (Gin handlers), router.go

### Requirement: SQL Migration Generation
The system SHALL generate goose-compatible SQL migration files from entity definitions.

#### Scenario: Generate migration
- **WHEN** user runs generate with entity definitions
- **THEN** system outputs goose timestamped SQL file with CREATE TABLE statement matching project conventions (snowflake ID, tenant_id, audit fields, soft delete)

### Requirement: Frontend API Generation
The system SHALL generate TypeScript API client code from backend API definitions.

#### Scenario: Generate API client
- **WHEN** user runs generate with API route definitions
- **THEN** system outputs TypeScript files with typed request/response interfaces and axios-based API functions
