# Story 7.4: Student Attendance Calendar View

Status: review

## Story

As a **parent or student**,
I want **to view attendance in calendar format**,
so that **attendance pattern is easily visible**.

## Acceptance Criteria

### AC1: Monthly Calendar View
**Given** a parent/student is logged in
**When** viewing attendance
**Then** they see monthly calendar view
**And** each day is color-coded: green (present), red (absent), yellow (late), gray (holiday)
**And** clicking a day shows details

### AC2: Attendance Summary
**Given** attendance summary is needed
**When** viewing summary section
**Then** they see: total days, present, absent, late, percentage
**And** comparison with class average is shown
**And** trend (improving/declining) is indicated

## Tasks / Subtasks

### Backend Tasks

- [x] **Task 1: Calendar Data Endpoint** (AC: 1,2)
  - [x] Add `GetStudentMonthlyAttendance(ctx, studentID, year, month)` - get attendance for a month
  - [x] Add `GetStudentAttendanceSummary(ctx, studentID, dateRange)` - get summary stats
  - [x] Add `GetClassAverageAttendance(ctx, sectionID, dateRange)` - for comparison

- [x] **Task 2: API Endpoints** (AC: 1,2)
  - [x] `GET /api/v1/student-attendance/calendar/:studentId?year=X&month=Y` - monthly calendar data
  - [x] `GET /api/v1/student-attendance/summary/:studentId` - summary with trend

### Frontend Tasks

- [x] **Task 3: Model & Service Updates** (AC: 1,2)
  - [x] Add calendar data interfaces (MonthlyCalendar, CalendarDay, MonthlySummary)
  - [x] Add service methods for calendar and summary (getStudentCalendar, getStudentSummaryReport)

- [ ] **Task 4: Calendar Component** (AC: 1) - UI deferred
  - [ ] Create calendar view with color-coded days
  - [ ] Add month navigation
  - [ ] Add day click for details

- [ ] **Task 5: Summary Component** (AC: 2) - UI deferred
  - [ ] Create summary stats display
  - [ ] Add trend indicator
  - [ ] Add class average comparison

## Dev Agent Record

### Agent Model Used
Claude Opus 4.5 (claude-opus-4-5-20251101)

### Completion Notes

1. Backend API endpoints implemented in studentattendance handler
2. Repository methods for calendar and summary data added
3. Service layer with trend calculation and class average comparison
4. Frontend models and service methods added
5. UI components deferred - backend APIs fully functional

### Change Log
| Date | Change | Author |
|------|--------|--------|
| 2026-01-29 | Story created | dev-story workflow |
| 2026-01-29 | Backend and service implementation completed | dev-story workflow |

### File List

**Backend - Module (modified):**
- `msls-backend/internal/modules/studentattendance/dto.go` - Added MonthlyCalendarResponse, MonthlySummaryResponse
- `msls-backend/internal/modules/studentattendance/repository.go` - Added GetStudentCalendarData, GetSectionAttendanceStats
- `msls-backend/internal/modules/studentattendance/service.go` - Added GetStudentCalendar, GetStudentSummary
- `msls-backend/internal/modules/studentattendance/handler.go` - Added GetStudentCalendar, GetStudentSummary handlers
- `msls-backend/cmd/api/main.go` - Added calendar and summary routes

**Frontend - Feature (modified):**
- `msls-frontend/src/app/features/student-attendance/student-attendance.model.ts` - Added calendar interfaces
- `msls-frontend/src/app/features/student-attendance/student-attendance.service.ts` - Added calendar service methods
