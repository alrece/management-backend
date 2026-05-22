## ADDED Requirements

### Requirement: Workflow Category Management
The system SHALL provide tree-structure workflow category CRUD.

#### Scenario: Create workflow category
- **WHEN** POST /workflow/categories with {name, parentId, sort}
- **THEN** system creates wf_category record, returns {code:0, data:id}

#### Scenario: List category tree
- **WHEN** GET /workflow/categories/tree
- **THEN** system returns nested tree structure of categories

### Requirement: Workflow Definition Management
The system SHALL provide CRUD for workflow definitions mapped to n8n workflows.

#### Scenario: Create workflow
- **WHEN** POST /workflow/workflows with {name, categoryId, paramsSchema}
- **THEN** system creates wf_workflow record (status=DRAFT), calls n8n API to create corresponding workflow, stores n8n_workflow_id

#### Scenario: Activate workflow
- **WHEN** PUT /workflow/workflows/:id/activate
- **THEN** system calls n8n Activate API, updates local status to ACTIVE

#### Scenario: Deactivate workflow
- **WHEN** PUT /workflow/workflows/:id/deactivate
- **THEN** system calls n8n Deactivate API, updates local status to INACTIVE

#### Scenario: Delete workflow
- **WHEN** DELETE /workflow/workflows/:id
- **THEN** system soft-deletes local record, calls n8n API to delete n8n workflow

### Requirement: Workflow Execution
The system SHALL support triggering and tracking workflow executions via n8n.

#### Scenario: Trigger workflow execution
- **WHEN** POST /workflow/workflows/:id/execute with {variables}
- **THEN** system calls n8n Execute API, creates wf_instance record (status=RUNNING), returns instance ID

#### Scenario: Receive execution callback
- **WHEN** n8n sends webhook callback with execution result
- **THEN** system updates wf_instance status/duration/error, stores n8n_execution_id

#### Scenario: Query instance status
- **WHEN** GET /workflow/instances/:id
- **THEN** system returns instance record with status, variables, duration

#### Scenario: List instances with pagination
- **WHEN** GET /workflow/instances/page with {workflowId?, status?, startTime?, endTime?}
- **THEN** system returns paginated instance records

### Requirement: Circuit Breaker Protection
The system SHALL protect n8n API calls with Circuit Breaker pattern.

#### Scenario: Circuit closed (normal)
- **WHEN** n8n API responds successfully
- **THEN** system passes response through normally

#### Scenario: Circuit opens (n8n down)
- **WHEN** 5 consecutive n8n API calls fail
- **THEN** system opens circuit, subsequent calls return 503 immediately without calling n8n

#### Scenario: Circuit half-open (probing)
- **WHEN** 30s after circuit opened, next request is allowed through
- **THEN** system allows up to 3 probe requests; success closes circuit, failure reopens

#### Scenario: Fallback on circuit open
- **WHEN** circuit is open and workflow activation is requested
- **THEN** system keeps local status as DRAFT, does not mark ACTIVE until n8n confirms

### Requirement: n8n Synchronization
The system SHALL periodically sync n8n workflow status with local records.

#### Scenario: Detect desync
- **WHEN** sync check (every 5min) finds n8n workflow deleted externally
- **THEN** system marks local status as DESYNCED and generates alert

#### Scenario: Startup reconciliation
- **WHEN** workflow-service starts up
- **THEN** system queries n8n for all active executions, reconciles with local RUNNING instances

### Requirement: Tenant Isolation
The system SHALL enforce tenant isolation for all workflow operations.

#### Scenario: Tenant-scoped workflow query
- **WHEN** tenant user queries workflows
- **THEN** system returns only workflows belonging to that tenant_id
