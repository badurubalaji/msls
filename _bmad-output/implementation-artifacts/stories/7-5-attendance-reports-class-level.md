# Story 7.5: Attendance Reports - Class Level

Status: review

## Story

As an **admin or teacher**,
I want **to view class-level attendance reports**,
so that **I can monitor overall class attendance trends**.

## Acceptance Criteria

### AC1: Daily Class Report
**Given** an admin/teacher is viewing class reports
**When** selecting a class and date
**Then** they see daily attendance report with all students
**And** summary shows total, present, absent, late counts
**And** attendance rate percentage is displayed

### AC2: Monthly Class Report
**Given** monthly view is selected
**When** viewing the report
**Then** they see attendance grid for all working days
**And** each student shows daily status for the month
**And** monthly summary includes students above 90%, below 75%, below 60%

## Tasks / Subtasks

### Backend Tasks

- [x] **Task 1: Class Report Endpoints** (AC: 1,2)
  - [x] Add `GetClassReport(ctx, sectionID, date)` - daily class report
  - [x] Add `GetMonthlyClassReport(ctx, sectionID, year, month)` - monthly grid report
  - [x] Add repository methods for aggregating class-level data

- [x] **Task 2: API Endpoints** (AC: 1,2)
  - [x] `GET /api/v1/student-attendance/reports/class/:sectionId?date=YYYY-MM-DD` - daily report
  - [x] `GET /api/v1/student-attendance/reports/class/:sectionId/monthly?year=X&month=Y` - monthly report

### Frontend Tasks

- [x] **Task 3: Model & Service Updates** (AC: 1,2)
  - [x] Add report interfaces (ClassReport, MonthlyClassReport, MonthlyStudentReport)
  - [x] Add service methods for class reports

- [ ] **Task 4: Class Report Component** (AC: 1,2) - UI deferred
  - [ ] Create daily report view with student list
  - [ ] Create monthly grid view with attendance matrix
  - [ ] Add export functionality

## Dev Agent Record

### Agent Model Used
Claude Opus 4.5 (claude-opus-4-5-20251101)

### Completion Notes

1. Backend API endpoints implemented for daily and monthly class reports
2. Repository methods aggregate data across students in a section
3. Monthly report includes working days calculation and student percentages
4. Frontend models and service methods added
5. UI components deferred - backend APIs fully functional

### Change Log
| Date | Change | Author |
|------|--------|--------|
| 2026-01-29 | Story created and implemented | dev-story workflow |

### File List

**Backend - Module (modified):**
- `msls-backend/internal/modules/studentattendance/dto.go` - Added ClassReportResponse, MonthlyClassReportResponse
- `msls-backend/internal/modules/studentattendance/repository.go` - Added GetClassDailyReport, GetClassMonthlyReport
- `msls-backend/internal/modules/studentattendance/service.go` - Added GetClassReport, GetMonthlyClassReport
- `msls-backend/internal/modules/studentattendance/handler.go` - Added GetClassReport, GetMonthlyClassReport handlers
- `msls-backend/cmd/api/main.go` - Added class report routes

**Frontend - Feature (modified):**
- `msls-frontend/src/app/features/student-attendance/student-attendance.model.ts` - Added report interfaces
- `msls-frontend/src/app/features/student-attendance/student-attendance.service.ts` - Added report service methods
