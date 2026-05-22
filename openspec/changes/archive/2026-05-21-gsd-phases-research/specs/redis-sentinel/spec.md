## ADDED Requirements

### Requirement: Sentinel Configuration
The system SHALL support Redis Sentinel via configuration toggle with automatic client selection.

#### Scenario: Single Redis mode (default)
- **WHEN** config.redis.sentinel.enabled is false
- **THEN** system connects using redis.NewClient with standalone addr

#### Scenario: Sentinel mode
- **WHEN** config.redis.sentinel.enabled is true
- **THEN** system connects using redis.NewFailoverClient with sentinel addrs and master name

#### Scenario: Sentinel configuration structure
- **WHEN** config.yaml contains redis.sentinel section
- **THEN** system parses enabled, master_name, and addrs fields correctly

### Requirement: Docker Compose Sentinel
The system SHALL provide 3 Sentinel containers + 1 Redis master in docker-compose.yml.

#### Scenario: Sentinel cluster starts
- **WHEN** docker-compose up is executed
- **THEN** 3 Sentinel containers monitor the Redis master, automatic failover enabled
