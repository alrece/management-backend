## ADDED Requirements

### Requirement: Job Task CRUD
The system SHALL provide REST API for managing scheduled job tasks with full CRUD operations.

#### Scenario: Create job task
- **WHEN** POST /job/tasks with {name, handler, cron_expr, params, status}
- **THEN** system creates job_task record with snowflake ID, returns {code:0, data:id}

#### Scenario: List job tasks with pagination
- **WHEN** GET /job/tasks/page with {page, pageSize, name?, status?}
- **THEN** system returns PageResult{list, total, page, pageSize} matching R{code:0} format

#### Scenario: Update job task
- **WHEN** PUT /job/tasks/:id with updated fields
- **THEN** system updates record, returns {code:0}

#### Scenario: Delete job task
- **WHEN** DELETE /job/tasks/:id
- **THEN** system soft-deletes record, returns {code:0}

#### Scenario: Manual trigger
- **WHEN** POST /job/tasks/:id/trigger
- **THEN** system creates execution log with MANUAL trigger type, executes task handler

### Requirement: Distributed Scheduling
The system SHALL implement Redis-based leader election for scheduler coordination across multiple instances.

#### Scenario: Leader election
- **WHEN** job-service instance starts
- **THEN** system attempts SETNX job:scheduler:leader with instance ID and TTL=30s; winner starts cron scheduler

#### Scenario: Leader renewal
- **WHEN** instance holds leader lock
- **THEN** system renews TTL every 8s

#### Scenario: Leader failover
- **WHEN** leader instance crashes and TTL expires
- **THEN** other instances race to acquire leader lock, new winner starts scheduler

### Requirement: Execution Logging
The system SHALL record execution logs with status tracking for every task invocation.

#### Scenario: Successful execution
- **WHEN** task handler completes without error
- **THEN** system logs {status:SUCCESS, duration_ms, result} in job_execution_log

#### Scenario: Failed execution
- **WHEN** task handler returns error
- **THEN** system logs {status:FAILED, error} and increments retry counter

#### Scenario: Execution timeout
- **WHEN** task handler exceeds 30min timeout
- **THEN** system marks execution as TIMEOUT, aborts handler

#### Scenario: Query execution logs
- **WHEN** GET /job/execution-logs/page with {taskId?, status?, startTime?, endTime?}
- **THEN** system returns paginated execution log records

### Requirement: Task Idempotency
The system SHALL prevent duplicate task execution using Redis idempotent guard.

#### Scenario: Deduplicate concurrent execution
- **WHEN** two instances attempt to execute same task at same scheduled time
- **THEN** only one instance executes, other skips via SETNX job:exec:{taskId}:{timestamp}

### Requirement: Tenant Isolation
The system SHALL enforce tenant isolation for all job operations.

#### Scenario: Tenant-scoped task query
- **WHEN** tenant user queries job tasks
- **THEN** system returns only tasks belonging to that tenant_id

#### Scenario: Cross-tenant access denied
- **WHEN** tenant user attempts to access another tenant's task
- **THEN** system returns 403 Forbidden
