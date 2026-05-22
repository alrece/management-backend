## ADDED Requirements

### Requirement: Vue3 Vben5 Project Setup
The system SHALL generate a Vue3 + Vben5 frontend prototype project for API contract validation.

#### Scenario: Project initialization
- **WHEN** frontend prototype project is created
- **THEN** system has Vue3 + Vben5 + TypeScript + Tailwind CSS configured, dev server runs successfully

### Requirement: Login Flow Prototype
The system SHALL implement login page with JWT dual-token flow.

#### Scenario: Login with valid credentials
- **WHEN** user submits username and password
- **THEN** frontend calls POST /system/auth/login, stores tokens, redirects to dashboard

#### Scenario: Token refresh
- **WHEN** access token expires
- **THEN** frontend automatically calls POST /system/auth/refresh with refresh token

### Requirement: Dashboard Layout Prototype
The system SHALL implement admin dashboard layout with dynamic menu.

#### Scenario: Load dynamic menu
- **WHEN** user enters dashboard after login
- **THEN** frontend calls GET /system/auth/get-permission-info, renders menu tree based on response

#### Scenario: Permission-based routing
- **WHEN** user navigates to a route
- **THEN** frontend checks permission identifiers, blocks unauthorized routes

### Requirement: System Module Pages
The system SHALL implement CRUD pages for core system modules.

#### Scenario: User management page
- **WHEN** admin opens user management
- **THEN** page shows user list with pagination, supports create/edit/delete/filter operations

#### Scenario: Role management page
- **WHEN** admin opens role management
- **THEN** page shows role list with menu permission tree assignment

#### Scenario: Menu management page
- **WHEN** admin opens menu management
- **THEN** page shows tree-structured menu list, supports add/edit/delete menu/button/directory

#### Scenario: Department management page
- **WHEN** admin opens department management
- **THEN** page shows tree-structured department list

### Requirement: Mock API Support
The system SHALL support mock API responses for frontend prototype when backend is unavailable.

#### Scenario: Mock mode enabled
- **WHEN** VITE_USE_MOCK=true environment variable is set
- **THEN** frontend uses mock data instead of real API calls, all pages render correctly
