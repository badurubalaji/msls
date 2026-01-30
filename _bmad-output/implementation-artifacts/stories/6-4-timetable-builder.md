# Story 6.4: Timetable Builder

**Epic:** 6 - Academic Structure & Timetable
**Status:** done
**Priority:** High

## User Story

As an **academic administrator**,
I want **to create timetables for each section**,
So that **teaching schedule is organized**.

## Acceptance Criteria

### AC1: Timetable Builder Grid
**Given** admin is creating a timetable
**When** using the timetable builder
**Then** they see a grid: days (columns) × periods (rows)
**And** they can drag-drop subjects to slots
**And** they can assign teacher to each slot
**And** colors differentiate subjects

### AC2: Teacher Assignment
**Given** a subject is being assigned
**When** selecting the slot
**Then** available teachers for that subject are shown
**And** teacher's existing assignments are visible
**And** conflict warning shows if teacher is already assigned

### AC3: Publish Timetable
**Given** timetable is complete
**When** publishing the timetable
**Then** timetable becomes active
**And** previous timetable is archived
**And** teachers and students can view new timetable

## Technical Notes

### Database Schema
- `timetables` table: id, tenant_id, branch_id, section_id, academic_year_id, name, status (draft/published/archived), effective_from, effective_to, created_by
- `timetable_entries` table: id, tenant_id, timetable_id, day_of_week, period_slot_id, subject_id, staff_id, room_number, notes

### API Endpoints
- GET /api/v1/timetables - List timetables (filter by section, status)
- POST /api/v1/timetables - Create new timetable
- GET /api/v1/timetables/:id - Get timetable with entries
- PUT /api/v1/timetables/:id - Update timetable
- DELETE /api/v1/timetables/:id - Delete draft timetable
- POST /api/v1/timetables/:id/publish - Publish timetable
- POST /api/v1/timetables/:id/archive - Archive timetable
- GET /api/v1/timetables/:id/entries - Get timetable entries
- POST /api/v1/timetables/:id/entries - Add/update entry
- DELETE /api/v1/timetables/:id/entries/:entryId - Remove entry
- GET /api/v1/timetables/conflicts - Check for teacher conflicts
- GET /api/v1/timetables/teacher/:staffId - Get teacher's timetable

### Frontend Components
- TimetableListComponent - List and manage timetables
- TimetableBuilderComponent - Grid-based timetable builder with drag-drop
- TimetableEntryModalComponent - Modal for adding/editing entries
- TeacherAvailabilityComponent - Show teacher's schedule conflicts

## Tasks

- [x] Task 1: Create database migrations for timetables and timetable_entries
- [x] Task 2: Implement backend models and repositories
- [x] Task 3: Implement backend services with conflict detection
- [x] Task 4: Implement backend HTTP handlers and routes
- [x] Task 5: Create frontend models and services
- [x] Task 6: Create Timetable List page
- [x] Task 7: Create Timetable Builder with grid view
- [x] Task 8: Implement entry modal for slot assignments (click-to-edit instead of drag-drop)
- [x] Task 9: Add teacher conflict detection backend API
- [x] Task 10: Add publish/archive functionality

## Dev Agent Record

### Started: 2026-01-29

## File List

### Backend Files

#### New Files
- `msls-backend/migrations/000050_timetables.up.sql` - Create timetables and timetable_entries tables with RLS
- `msls-backend/migrations/000050_timetables.down.sql` - Drop migration
- `msls-backend/internal/pkg/database/models/timetable.go` - Timetable and TimetableEntry models
- `msls-backend/internal/modules/timetable/repository.go` - Timetable repository with CRUD and conflict detection
- `msls-backend/internal/modules/timetable/service.go` - Timetable service with business logic
- `msls-backend/internal/modules/timetable/handler.go` - HTTP handlers for timetable routes
- `msls-backend/internal/modules/timetable/dto.go` - DTOs for timetable API
- `msls-backend/internal/modules/timetable/errors.go` - Custom error definitions

#### Modified Files
- `msls-backend/cmd/api/main.go` - Register timetable routes

### Frontend Files

#### New Files
- `msls-frontend/src/app/features/academics/timetable/timetable-builder/timetable-builder.component.ts` - Grid-based timetable builder with entry editing
- `msls-frontend/src/app/features/academics/timetable/timetable-list/timetable-list.component.ts` - Timetable list with filtering and creation

#### Modified Files
- `msls-frontend/src/app/features/academics/timetable/timetable.model.ts` - Added Timetable, TimetableEntry, conflict types
- `msls-frontend/src/app/features/academics/timetable/timetable.service.ts` - Added timetable API methods
- `msls-frontend/src/app/features/academics/timetable/timetable.component.ts` - Added Timetable Builder card
- `msls-frontend/src/app/features/academics/academics.routes.ts` - Added timetable list and builder routes

## Change Log

- 2026-01-29: Story created and implementation started
- 2026-01-29: Completed full implementation - backend APIs, frontend list and builder components
