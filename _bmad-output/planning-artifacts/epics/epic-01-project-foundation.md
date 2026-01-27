# Epic 1: Project Foundation & Design System

**Phase:** 1 (MVP)
**Priority:** Critical - Must be completed first

## Epic Goal

Establish the complete technical foundation including project scaffolding, database setup, authentication infrastructure, and reusable design system so that developers can build features consistently and securely.

## User Value

Administrators can access a secure, professionally designed system with proper login and tenant isolation.

## FRs Covered

FR-CF-01, FR-CF-06, FR-CF-08, FR-CF-09, NFR-SEC-01 to NFR-SEC-10, NFR-PER-01 to NFR-PER-06, UX-COMP, UX-PAT, ARCH-STACK, ARCH-STRUCT

---

## Stories

### Story 1.1: Backend Project Scaffolding

As a **developer**,
I want **a properly structured Go backend project with all necessary dependencies and configurations**,
So that **I can start building features following established patterns**.

**Acceptance Criteria:**

**Given** a new development environment
**When** the developer clones the repository and runs setup
**Then** the Go project structure matches architecture specifications:
- `cmd/api/main.go` entry point exists
- `internal/modules/` directory structure is created
- `internal/pkg/` for shared utilities exists
- `configs/` with environment-based configuration
- `migrations/` directory for database migrations
**And** all core dependencies are installed (Gin, GORM, sqlc, validator, uuid)
**And** Makefile includes targets for build, test, run, migrate
**And** Docker Compose file includes PostgreSQL 16, Redis, MinIO services
**And** `.env.example` documents all required environment variables

---

### Story 1.2: Database Foundation with Multi-Tenancy

As a **system**,
I want **PostgreSQL database with Row-Level Security policies and multi-tenant support**,
So that **tenant data is completely isolated at the database level**.

**Acceptance Criteria:**

**Given** the database is initialized
**When** migrations are applied
**Then** the following foundation tables exist:
- `tenants` (id, name, slug, settings, status, created_at, updated_at)
- `branches` (id, tenant_id, name, code, address, settings, is_primary, status)
- `users` (id, tenant_id, email, phone, password_hash, status, 2fa_enabled)
**And** UUID v7 generation function `uuid_generate_v7()` is available
**And** RLS is enabled on all tables with tenant_id
**And** RLS policy uses `current_setting('app.tenant_id')::UUID`
**And** Indexes exist on tenant_id for all tables
**And** Audit columns (created_at, updated_at, created_by, updated_by) exist on all tables

---

### Story 1.3: Core Backend Infrastructure

As a **developer**,
I want **reusable backend infrastructure components**,
So that **all modules follow consistent patterns for database, logging, and error handling**.

**Acceptance Criteria:**

**Given** a developer needs to create a new module
**When** they use the core infrastructure
**Then** the following components are available:
- Database connection pool with context propagation
- Transaction helper with automatic rollback on panic
- Tenant context middleware that sets `app.tenant_id` session variable
- Request ID middleware for tracing
- Structured JSON logging with correlation ID
- Standard error types (NotFound, Validation, Unauthorized, Forbidden)
- RFC 7807 error response formatter
**And** API rate limiting middleware exists (configurable per endpoint)
**And** Request/response logging middleware captures method, path, status, duration
**And** Panic recovery middleware prevents server crashes

---

### Story 1.4: Frontend Project Scaffolding

As a **developer**,
I want **a properly structured Angular 21 frontend project**,
So that **I can build features using standalone components and Signals**.

**Acceptance Criteria:**

**Given** a new development environment
**When** the developer clones the repository and runs npm install
**Then** Angular 21 project structure exists:
- `src/app/features/` for feature modules
- `src/app/shared/components/` for shared components
- `src/app/shared/services/` for shared services
- `src/app/core/` for singleton services
- `src/app/layouts/` for layout components
**And** Tailwind CSS is configured with custom theme colors
**And** ESLint and Prettier are configured for code quality
**And** Environment files for dev/staging/prod exist
**And** Proxy configuration for API calls is set up
**And** Path aliases (@app, @shared, @features, @core) are configured

---

### Story 1.5: Design System - Atomic Components

As a **developer**,
I want **reusable atomic UI components (Button, Input, Badge, Avatar, Icon)**,
So that **all features have consistent visual styling and behavior**.

**Acceptance Criteria:**

**Given** a developer needs a button in their feature
**When** they use the MslsButton component
**Then** it supports variants: primary, secondary, danger, ghost
**And** it supports sizes: sm, md, lg
**And** it supports states: default, hover, focus, disabled, loading
**And** it meets WCAG 2.1 AA contrast requirements
**And** focus ring is visible for keyboard navigation

**Given** a developer needs an input field
**When** they use the MslsInput component
**Then** it supports types: text, email, password, number, tel, search
**And** it displays error state with message below
**And** it supports prefix/suffix icons
**And** it meets minimum 44x44px touch target on mobile

**Given** a developer needs a badge or avatar
**When** they use MslsBadge or MslsAvatar components
**Then** Badge supports variants: success, warning, error, info, neutral
**And** Avatar supports sizes and fallback initials

---

### Story 1.6: Design System - Molecule Components

As a **developer**,
I want **reusable molecule components (FormField, Card, Dropdown, Toast, Modal)**,
So that **complex UI patterns are consistent across the application**.

**Acceptance Criteria:**

**Given** a developer needs a form field with label and validation
**When** they use the MslsFormField component
**Then** it wraps input with label, hint text, and error message
**And** error message appears only when field is touched and invalid
**And** required indicator (*) appears for required fields

**Given** a developer needs a card container
**When** they use the MslsCard component
**Then** it supports header, body, footer sections
**And** it supports variants: default, elevated, outlined

**Given** a developer needs notifications
**When** they use the ToastService
**Then** toasts appear in top-right corner
**And** variants exist: success (green), error (red), warning (amber), info (blue)
**And** toasts auto-dismiss after configurable duration
**And** toasts can be manually dismissed

**Given** a developer needs a modal dialog
**When** they use the MslsModal component
**Then** it supports sizes: sm, md, lg, xl
**And** it traps focus within modal
**And** it closes on Escape key or backdrop click (configurable)

---

### Story 1.7: Core HTTP Client & Authentication Service

As a **frontend developer**,
I want **a configured HTTP client with authentication interceptors**,
So that **all API calls automatically include tokens and handle auth errors**.

**Acceptance Criteria:**

**Given** a user is logged in
**When** an API request is made
**Then** the Authorization header includes the JWT access token
**And** the X-Tenant-ID header includes the current tenant ID

**Given** an API returns 401 Unauthorized
**When** a refresh token exists and is valid
**Then** the system automatically refreshes the access token
**And** the original request is retried with new token
**And** if refresh fails, user is redirected to login

**Given** the application needs to track requests
**When** any API call is made
**Then** a unique X-Request-ID header is included
**And** loading state is trackable via a global loading service

---

### Story 1.8: Application Shell & Navigation Layout

As an **admin user**,
I want **a responsive application shell with sidebar navigation**,
So that **I can navigate between different modules easily**.

**Acceptance Criteria:**

**Given** an admin user is logged in
**When** they view the application on desktop (>1024px)
**Then** they see a collapsible sidebar with navigation menu
**And** sidebar shows module icons and labels
**And** current route is highlighted in sidebar
**And** sidebar collapse state persists across sessions

**Given** an admin user is on mobile (<768px)
**When** they view the application
**Then** sidebar is hidden by default
**And** hamburger menu icon in header toggles sidebar
**And** sidebar overlays content when open

**Given** the application shell
**When** rendered
**Then** header includes school logo, tenant name, user menu
**And** user menu includes profile, settings, logout options
**And** main content area has proper spacing and scroll behavior

---

### Story 1.9: Development Tooling & Documentation

As a **developer**,
I want **comprehensive development tooling and documentation**,
So that **new team members can onboard quickly and maintain code quality**.

**Acceptance Criteria:**

**Given** a new developer joins the project
**When** they read the documentation
**Then** README includes setup instructions for both backend and frontend
**And** Architecture decision records (ADRs) document key decisions
**And** API documentation is auto-generated (Swagger/OpenAPI)

**Given** code is being developed
**When** pre-commit hooks run
**Then** Go code is formatted with gofmt and linted with golangci-lint
**And** TypeScript code is formatted with Prettier and linted with ESLint
**And** Commit messages follow conventional commits format

**Given** the CI pipeline runs
**When** a PR is submitted
**Then** All tests must pass
**And** Linting must pass with no errors
**And** Build must succeed for both backend and frontend
