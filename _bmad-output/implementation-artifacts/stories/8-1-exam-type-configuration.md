# Story 8.1: Exam Type Configuration

Status: complete

## Story

As an **academic administrator**,
I want **to configure different exam types**,
So that **various assessments can be created**.

## Acceptance Criteria

### AC1: Create Exam Type
**Given** admin is on exam settings
**When** creating an exam type
**Then** they can enter: name (e.g., "Unit Test", "Half Yearly", "Annual")
**And** they can set: weightage for final result calculation
**And** they can set: is marks-based or grade-based
**And** they can set: default maximum marks

### AC2: Manage Exam Types
**Given** exam types are configured
**When** viewing the list
**Then** they see all exam types with settings
**And** they can order exam types (display sequence)
**And** they can activate/deactivate types

## Tasks / Subtasks

### Backend Tasks

- [x] **Task 1: Database Schema** (AC: 1,2)
  - [x] Create `exam_types` table with columns: id, tenant_id, name, code, weightage, evaluation_type (marks/grade), default_max_marks, display_order, is_active
  - [x] Create migration file `000057_exam_types.up.sql` and down migration
  - [x] Add RLS policy for tenant isolation
  - [x] Add indexes for tenant_id and is_active

- [x] **Task 2: Model & Repository** (AC: 1,2)
  - [x] Create `ExamType` model in `internal/pkg/database/models/exam_type.go`
  - [x] Create `internal/modules/exam/` module directory
  - [x] Implement repository with CRUD operations
  - [x] Add `ListExamTypes`, `GetExamType`, `CreateExamType`, `UpdateExamType`, `DeleteExamType`
  - [x] Add `UpdateDisplayOrder` for reordering
  - [x] Add `ToggleActive` for activation/deactivation

- [x] **Task 3: Service Layer** (AC: 1,2)
  - [x] Create service with business logic
  - [x] Validate weightage (0-100%)
  - [x] Validate unique code per tenant
  - [x] Validate default_max_marks (positive number)
  - [x] Handle display order on create/delete

- [x] **Task 4: API Endpoints** (AC: 1,2)
  - [x] `GET /api/v1/exam-types` - List all exam types (with filters: active, search)
  - [x] `GET /api/v1/exam-types/:id` - Get single exam type
  - [x] `POST /api/v1/exam-types` - Create exam type
  - [x] `PUT /api/v1/exam-types/:id` - Update exam type
  - [x] `DELETE /api/v1/exam-types/:id` - Soft delete exam type
  - [x] `PATCH /api/v1/exam-types/:id/active` - Toggle active status
  - [x] `PUT /api/v1/exam-types/order` - Update display order (bulk)

- [x] **Task 5: Permissions** (AC: 1,2)
  - [x] Add permissions: `exam:type:view`, `exam:type:create`, `exam:type:update`, `exam:type:delete`
  - [x] Update migration with permission inserts

### Frontend Tasks

- [x] **Task 6: Model & Service** (AC: 1,2)
  - [x] Create `ExamType` interface in `features/exams/exam.model.ts`
  - [x] Create `ExamService` with API methods
  - [x] Add to `api.service.ts` exam endpoints

- [x] **Task 7: Exam Type List Component** (AC: 2)
  - [x] Create `exam-type-list.component.ts` with data table
  - [x] Add search and active filter
  - [x] Add up/down arrow reordering (simplified from drag-drop)
  - [x] Add activate/deactivate toggle
  - [x] Add edit and delete actions
  - [x] Add "How it Works" help section for teachers/admins

- [x] **Task 8: Exam Type Form Component** (AC: 1)
  - [x] Create `exam-type-form.component.ts` for create/edit
  - [x] Form fields: name, code, weightage, evaluation_type, default_max_marks
  - [x] Add form validation
  - [x] Add quick preset buttons for common exam types
  - [x] Add success/error notifications

- [x] **Task 9: Routing & Navigation** (AC: 1,2)
  - [x] Add exam routes to `exams.routes.ts`
  - [x] Add "Exams" section to navigation config with permission
  - [x] Add exams route to main app routes

## Dev Notes

### Architecture References
- Follow existing module pattern: handler.go, service.go, repository.go, dto.go
- Use GORM for database operations
- Apply RLS for tenant isolation
- Use standard response format for API responses

### Database Design
```sql
CREATE TABLE exam_types (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    name VARCHAR(100) NOT NULL,
    code VARCHAR(20) NOT NULL,
    weightage DECIMAL(5,2) DEFAULT 0,
    evaluation_type VARCHAR(20) NOT NULL DEFAULT 'marks', -- marks, grade
    default_max_marks INTEGER DEFAULT 100,
    display_order INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uniq_exam_type_code UNIQUE (tenant_id, code)
);
```

### API Response Format
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "Unit Test",
    "code": "UT",
    "weightage": 20.00,
    "evaluation_type": "marks",
    "default_max_marks": 50,
    "display_order": 1,
    "is_active": true
  }
}
```

### Frontend Patterns
- Use Angular 21 standalone components with signals
- Follow Tailwind CSS patterns from existing components
- Use reactive forms with validators

## Dev Agent Record

### Agent Model Used
Claude Opus 4.5 (claude-opus-4-5-20251101)

### Completion Notes
- Implemented complete exam types CRUD functionality
- Backend follows standard module pattern with handler, service, repository, dto, and errors
- Frontend uses Angular 21 standalone components with signals
- Used up/down arrow buttons for reordering instead of CDK drag-drop (simpler approach)
- Added quick preset buttons for common exam types (Unit Test, Mid-Term, Final, Practical, Project)
- Added comprehensive "How it Works" help section explaining evaluation types, weightage, and best practices for administrators
- All permissions added: exam:type:view, exam:type:create, exam:type:update, exam:type:delete
- Navigation added under "Exams" menu with permission guard

### Change Log
| Date | Change | Author |
|------|--------|--------|
| 2026-01-30 | Story created | dev-story workflow |
| 2026-01-30 | Implementation complete | Claude Opus 4.5 |

### File List
**Backend:**
- `msls-backend/migrations/000057_exam_types.up.sql` - Migration with table, RLS, indexes, permissions
- `msls-backend/migrations/000057_exam_types.down.sql` - Down migration
- `msls-backend/internal/pkg/database/models/exam_type.go` - ExamType model
- `msls-backend/internal/modules/exam/errors.go` - Error definitions
- `msls-backend/internal/modules/exam/dto.go` - Request/Response DTOs
- `msls-backend/internal/modules/exam/repository.go` - Database operations
- `msls-backend/internal/modules/exam/service.go` - Business logic
- `msls-backend/internal/modules/exam/handler.go` - HTTP handlers

**Frontend:**
- `msls-frontend/src/app/features/exams/exam.model.ts` - TypeScript interfaces
- `msls-frontend/src/app/features/exams/exam.service.ts` - API service
- `msls-frontend/src/app/features/exams/exam-type-list/exam-type-list.component.ts` - List with help section
- `msls-frontend/src/app/features/exams/exam-type-form/exam-type-form.component.ts` - Create/Edit form
- `msls-frontend/src/app/features/exams/exams.routes.ts` - Routing
- `msls-frontend/src/app/app.routes.ts` - Added exams route (modified)
- `msls-frontend/src/app/layouts/nav-config.ts` - Added Exams nav (modified)
