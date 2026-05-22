## ADDED Requirements

### Requirement: Goose Migration Integration
The system SHALL integrate goose/v3 for version-controlled database migrations, replacing manual SQL file execution.

#### Scenario: Run pending migrations
- **WHEN** service starts and goose.Up(db, migrationsDir) is called
- **THEN** system executes all pending migration files in order, tracking completion in goose_db_version table

#### Scenario: Migrate specific tenant
- **WHEN** cmd/migrate runs with --target=tenant_id flag
- **THEN** system connects to tenant database and runs goose.Up on that database only

#### Scenario: Migrate all tenants
- **WHEN** cmd/migrate runs with --all flag
- **THEN** system iterates all tenant databases and runs goose.Up on each

#### Scenario: Allow missing migrations
- **WHEN** goose.Up runs with WithAllowMissing option
- **THEN** system skips already-applied migrations without error

### Requirement: Embedded Migrations
The system SHALL support go:embed for migration SQL files.

#### Scenario: Binary contains migrations
- **WHEN** service binary is built with go:embed directive
- **THEN** migration SQL files are embedded in binary, no external file dependency

### Requirement: Standalone Migration Tool
The system SHALL provide cmd/migrate/main.go as standalone migration runner.

#### Scenario: Migration tool runs independently
- **WHEN** operator executes go run cmd/migrate/main.go --all
- **THEN** system migrates all tenant databases without starting the service
