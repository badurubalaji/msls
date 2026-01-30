# Story 7.7: Low Attendance Alerts

Status: review

## Story

As an **admin or teacher**,
I want **to be alerted about students with low attendance**,
so that **I can take proactive intervention measures**.

## Acceptance Criteria

### AC1: Low Attendance Dashboard
**Given** an admin is viewing alerts
**When** accessing low attendance dashboard
**Then** they see students below threshold (default 75%)
**And** chronic absentees (below 60%) are highlighted
**And** class-wise breakdown is provided

### AC2: Alert Configuration
**Given** threshold settings exist
**When** generating alerts
**Then** warning threshold (75%) and critical threshold (60%) are applied
**And** consecutive absent days are tracked
**And** last present date is shown for context

## Tasks / Subtasks

### Backend Tasks

- [x] **Task 1: Low Attendance Dashboard Endpoint** (AC: 1,2)
  - [x] Add `GetLowAttendanceDashboard(ctx, dateRange, threshold)` - students below threshold
  - [x] Track chronic absentees (below critical threshold)
  - [x] Calculate consecutive absent days
  - [x] Add class-wise breakdown

- [x] **Task 2: API Endpoints** (AC: 1,2)
  - [x] `GET /api/v1/student-attendance/alerts/low-attendance?date_from=X&date_to=Y&threshold=75` - dashboard

### Frontend Tasks

- [x] **Task 3: Model & Service Updates** (AC: 1,2)
  - [x] Add alert interfaces (LowAttendanceDashboard, LowAttendanceStudent)
  - [x] Add service method for low attendance dashboard

- [ ] **Task 4: Alert Dashboard Component** (AC: 1,2) - UI deferred
  - [ ] Create low attendance dashboard view
  - [ ] Add student list with attendance rates
  - [ ] Add class-wise breakdown visualization
  - [ ] Add export functionality

## Dev Agent Record

### Agent Model Used
Claude Opus 4.5 (claude-opus-4-5-20251101)

### Completion Notes

1. Low attendance dashboard endpoint with configurable thresholds
2. Chronic absentee tracking (below 60% attendance)
3. Consecutive absent days calculation for pattern identification
4. Class-wise breakdown for administrative overview
5. Frontend models and service methods added
6. UI components deferred - backend APIs fully functional

### Change Log
| Date | Change | Author |
|------|--------|--------|
| 2026-01-29 | Story created and implemented | dev-story workflow |

### File List

**Backend - Module (modified):**
- `msls-backend/internal/modules/studentattendance/dto.go` - Added LowAttendanceDashboardResponse, LowAttendanceStudentDTO
- `msls-backend/internal/modules/studentattendance/repository.go` - Added GetLowAttendanceStudents, GetClassAttendanceBreakdown
- `msls-backend/internal/modules/studentattendance/service.go` - Added GetLowAttendanceDashboard
- `msls-backend/internal/modules/studentattendance/handler.go` - Added GetLowAttendanceDashboard handler
- `msls-backend/cmd/api/main.go` - Added alert routes

**Frontend - Feature (modified):**
- `msls-frontend/src/app/features/student-attendance/student-attendance.model.ts` - Added LowAttendanceDashboard interface
- `msls-frontend/src/app/features/student-attendance/student-attendance.service.ts` - Added getLowAttendanceDashboard method
