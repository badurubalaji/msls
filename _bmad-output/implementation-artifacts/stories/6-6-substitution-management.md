# Story 6.6: Substitution Management

**Epic:** 6 - Academic Structure & Timetable
**Status:** done
**Priority:** Medium

## User Story

As a **school administrator**,
I want **to manage teacher substitutions when a teacher is absent**,
So that **classes are not disrupted and students continue learning**.

## Acceptance Criteria

### AC1: Create Substitution
**Given** a teacher is absent
**When** admin creates a substitution
**Then** they can select the absent teacher, date, periods, and substitute teacher
**And** the system checks for substitute teacher availability

### AC2: View Substitutions
**Given** substitutions have been created
**When** viewing the substitution list
**Then** admin sees all substitutions with date, absent teacher, substitute, and status
**And** can filter by date range, teacher, or status

### AC3: Teacher Notification
**Given** a substitution is created
**When** the substitution is saved
**Then** the substitute teacher can see the assignment in their schedule

### AC4: Substitution Calendar
**Given** admin views substitutions
**When** accessing the calendar view
**Then** they see substitutions overlaid on the timetable grid
**And** can easily identify coverage gaps

## Technical Notes

### Database Schema
- `substitutions` table: id, tenant_id, branch_id, original_staff_id, substitute_staff_id, date, reason, status, notes
- `substitution_periods` table: id, substitution_id, timetable_entry_id, period_slot_id

### API Endpoints
- GET /api/v1/timetables/substitutions - List substitutions with filters
- POST /api/v1/timetables/substitutions - Create substitution
- GET /api/v1/timetables/substitutions/:id - Get substitution details
- PUT /api/v1/timetables/substitutions/:id - Update substitution
- DELETE /api/v1/timetables/substitutions/:id - Delete substitution
- POST /api/v1/timetables/substitutions/:id/confirm - Confirm substitution
- POST /api/v1/timetables/substitutions/:id/cancel - Cancel substitution
- GET /api/v1/timetables/substitutions/available-teachers - Get available substitute teachers
- GET /api/v1/timetables/substitutions/teacher-periods - Get teacher's periods for a date

### Frontend Components
- SubstitutionListComponent - List and filter substitutions with status stats
- SubstitutionFormComponent - Multi-step wizard to create substitution
- SubstitutionDetailComponent - View substitution details and activity timeline

## Tasks

- [x] Task 1: Create database migration for substitutions
- [x] Task 2: Implement backend models and DTOs
- [x] Task 3: Implement backend repository
- [x] Task 4: Implement backend service with availability check
- [x] Task 5: Implement backend HTTP handlers
- [x] Task 6: Create frontend models and service
- [x] Task 7: Create Substitution List page
- [x] Task 8: Create Substitution Form wizard
- [x] Task 9: Create Substitution Detail page
- [x] Task 10: Add routes and navigation

## Dev Agent Record

### Started: 2026-01-29
### Completed: 2026-01-29

## File List

### Backend Files

#### New Files
- `msls-backend/migrations/000052_substitutions.up.sql` - Migration for substitutions and substitution_periods tables with RLS
- `msls-backend/migrations/000052_substitutions.down.sql` - Rollback migration
- `msls-backend/internal/pkg/database/models/substitution.go` - Substitution and SubstitutionPeriod GORM models
- `msls-backend/internal/modules/timetable/substitution_dto.go` - DTOs and response converters
- `msls-backend/internal/modules/timetable/substitution_repository.go` - Database operations
- `msls-backend/internal/modules/timetable/substitution_service.go` - Business logic with conflict detection
- `msls-backend/internal/modules/timetable/substitution_handler.go` - HTTP handlers with permission-based routes

#### Modified Files
- `msls-backend/internal/modules/timetable/errors.go` - Added substitution-related errors
- `msls-backend/cmd/api/main.go` - Registered substitution routes

### Frontend Files

#### New Files
- `msls-frontend/src/app/features/academics/timetable/substitution/substitution.model.ts` - TypeScript models and status config
- `msls-frontend/src/app/features/academics/timetable/substitution/substitution.service.ts` - Signal-based Angular service
- `msls-frontend/src/app/features/academics/timetable/substitution/substitution-list/substitution-list.component.ts` - List with filters and status stats
- `msls-frontend/src/app/features/academics/timetable/substitution/substitution-form/substitution-form.component.ts` - Multi-step creation wizard
- `msls-frontend/src/app/features/academics/timetable/substitution/substitution-detail/substitution-detail.component.ts` - Detail view with activity timeline

#### Modified Files
- `msls-frontend/src/app/features/academics/academics.routes.ts` - Added substitution routes
- `msls-frontend/src/app/layouts/nav-config.ts` - Added Substitutions nav item

## Features Implemented

### Backend
- Full CRUD operations for substitutions
- Multi-period support per substitution
- Teacher availability checking against:
  - Existing timetable entries
  - Other substitutions on the same date
- Substitution workflow: pending → confirmed → completed (or cancelled)
- Permission-based route protection (substitution:view, substitution:create, substitution:update, substitution:delete, substitution:approve)
- Row-level security with tenant isolation

### Frontend
- Substitution list with:
  - Filters by branch, status, date range
  - Status statistics cards
  - Pagination support
  - Quick actions (confirm, cancel)
- Multi-step substitution creation wizard:
  1. Select absent teacher and date
  2. Choose periods to cover
  3. Select available substitute teacher
- Substitution detail view with:
  - Teacher information cards
  - Covered periods list
  - Activity timeline
  - Action buttons based on status

## Change Log

- 2026-01-29: Story created and implementation started
- 2026-01-29: Backend implementation completed (migration, models, DTOs, repository, service, handler)
- 2026-01-29: Frontend implementation completed (models, service, list, form, detail components)
- 2026-01-29: Story marked as done
