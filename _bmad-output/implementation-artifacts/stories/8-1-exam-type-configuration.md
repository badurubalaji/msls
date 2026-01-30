# Story 8.1: Exam Type Configuration

Status: ready-for-dev

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

- [ ] **Task 1: Database Schema** (AC: 1,2)
  - [ ] Create `exam_types` table with columns: id, tenant_id, name, code, weightage, evaluation_type (marks/grade), default_max_marks, display_order, is_active
  - [ ] Create migration file `000057_exam_types.up.sql` and down migration
  - [ ] Add RLS policy for tenant isolation
  - [ ] Add indexes for tenant_id and is_active

- [ ] **Task 2: Model & Repository** (AC: 1,2)
  - [ ] Create `ExamType` model in `internal/pkg/database/models/exam_type.go`
  - [ ] Create `internal/modules/exam/` module directory
  - [ ] Implement repository with CRUD operations
  - [ ] Add `ListExamTypes`, `GetExamType`, `CreateExamType`, `UpdateExamType`, `DeleteExamType`
  - [ ] Add `UpdateDisplayOrder` for reordering
  - [ ] Add `ToggleActive` for activation/deactivation

- [ ] **Task 3: Service Layer** (AC: 1,2)
  - [ ] Create service with business logic
  - [ ] Validate weightage (0-100%)
  - [ ] Validate unique code per tenant
  - [ ] Validate default_max_marks (positive number)
  - [ ] Handle display order on create/delete

- [ ] **Task 4: API Endpoints** (AC: 1,2)
  - [ ] `GET /api/v1/exam-types` - List all exam types (with filters: active, search)
  - [ ] `GET /api/v1/exam-types/:id` - Get single exam type
  - [ ] `POST /api/v1/exam-types` - Create exam type
  - [ ] `PUT /api/v1/exam-types/:id` - Update exam type
  - [ ] `DELETE /api/v1/exam-types/:id` - Soft delete exam type
  - [ ] `PATCH /api/v1/exam-types/:id/active` - Toggle active status
  - [ ] `PUT /api/v1/exam-types/order` - Update display order (bulk)

- [ ] **Task 5: Permissions** (AC: 1,2)
  - [ ] Add permissions: `exam:type:view`, `exam:type:create`, `exam:type:update`, `exam:type:delete`
  - [ ] Update migration with permission inserts

### Frontend Tasks

- [ ] **Task 6: Model & Service** (AC: 1,2)
  - [ ] Create `ExamType` interface in `features/exams/exam.model.ts`
  - [ ] Create `ExamService` with API methods
  - [ ] Add to `api.service.ts` exam endpoints

- [ ] **Task 7: Exam Type List Component** (AC: 2)
  - [ ] Create `exam-type-list.component.ts` with data table
  - [ ] Add search and active filter
  - [ ] Add drag-and-drop reordering
  - [ ] Add activate/deactivate toggle
  - [ ] Add edit and delete actions

- [ ] **Task 8: Exam Type Form Component** (AC: 1)
  - [ ] Create `exam-type-form.component.ts` for create/edit
  - [ ] Form fields: name, code, weightage, evaluation_type, default_max_marks
  - [ ] Add form validation
  - [ ] Add success/error notifications

- [ ] **Task 9: Routing & Navigation** (AC: 1,2)
  - [ ] Add exam routes to `academics.routes.ts`
  - [ ] Add "Exams" section to navigation config
  - [ ] Create exams landing component

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
_To be filled during implementation_

### Completion Notes
_To be filled during implementation_

### Change Log
| Date | Change | Author |
|------|--------|--------|
| 2026-01-30 | Story created | dev-story workflow |

### File List
_To be filled during implementation_
