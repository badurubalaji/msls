# Story 6.5: Teacher Timetable View

**Epic:** 6 - Academic Structure & Timetable
**Status:** done
**Priority:** High

## User Story

As a **teacher**,
I want **to view my teaching schedule across all sections**,
So that **I can plan my day and prepare for classes**.

## Acceptance Criteria

### AC1: Teacher's Weekly Schedule
**Given** a teacher is logged in
**When** viewing their timetable
**Then** they see their complete weekly schedule
**And** all assigned periods are shown with subject, class, and time
**And** free periods are clearly indicated

### AC2: Day View
**Given** teacher is viewing their schedule
**When** selecting a specific day
**Then** they see detailed view for that day
**And** periods are shown in chronological order
**And** current/next period is highlighted

### AC3: Schedule Overview
**Given** teacher views their timetable
**When** the page loads
**Then** they see weekly period count by subject
**And** total teaching hours are displayed
**And** free periods count is shown

## Technical Notes

### API Endpoints
- GET /api/v1/timetables/teacher/:staffId - Get teacher's schedule (already implemented)
- GET /api/v1/timetables/teacher/me - Get current user's teaching schedule

### Frontend Components
- TeacherTimetableComponent - Weekly grid view of teacher's schedule
- Reuses existing timetable grid styling

## Tasks

- [x] Task 1: Add /me endpoint for current user's timetable
- [x] Task 2: Create frontend service method for my timetable
- [x] Task 3: Create TeacherTimetableComponent with weekly view
- [x] Task 4: Add day filter and current period highlighting
- [x] Task 5: Add schedule statistics (period count, subjects taught)
- [x] Task 6: Add route and navigation

## Dev Agent Record

### Started: 2026-01-29

## File List

### Backend Files

#### Modified Files
- `msls-backend/internal/modules/timetable/handler.go` - Add /me endpoint

### Frontend Files

#### New Files
- `msls-frontend/src/app/features/academics/timetable/teacher-timetable/teacher-timetable.component.ts`

#### Modified Files
- `msls-frontend/src/app/features/academics/timetable/timetable.service.ts`
- `msls-frontend/src/app/features/academics/academics.routes.ts`

## Change Log

- 2026-01-29: Story created and implementation started
- 2026-01-29: Implementation completed - backend /me endpoint, frontend TeacherTimetableComponent with weekly/day views
