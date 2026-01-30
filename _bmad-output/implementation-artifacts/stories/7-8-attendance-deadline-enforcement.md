# Story 7.8: Attendance Deadline Enforcement

Status: review

## Story

As an **admin**,
I want **to enforce attendance marking deadlines**,
so that **teachers mark attendance on time**.

## Acceptance Criteria

### AC1: Unmarked Attendance Tracking
**Given** a deadline time is configured (e.g., 11 AM)
**When** viewing unmarked attendance
**Then** classes without attendance are listed
**And** post-deadline unmarked classes are escalated
**And** teacher name is shown for accountability

### AC2: Deadline Configuration
**Given** branch settings exist
**When** deadlines are configured
**Then** edit window and late threshold are enforced
**And** settings are customizable per branch

## Tasks / Subtasks

### Backend Tasks

- [x] **Task 1: Unmarked Attendance Endpoint** (AC: 1)
  - [x] Add `GetUnmarkedAttendance(ctx, date)` - classes without attendance
  - [x] Calculate post-deadline status based on configured deadline
  - [x] Include teacher information for accountability

- [x] **Task 2: Settings Endpoint** (AC: 2)
  - [x] Reuse existing settings endpoint from Story 7.1
  - [x] Add edit window and threshold enforcement in service layer

- [x] **Task 3: API Endpoints** (AC: 1)
  - [x] `GET /api/v1/student-attendance/alerts/unmarked?date=YYYY-MM-DD` - unmarked classes

### Frontend Tasks

- [x] **Task 4: Model & Service Updates** (AC: 1)
  - [x] Add unmarked attendance interfaces (UnmarkedAttendance, UnmarkedClassInfo)
  - [x] Add service method for unmarked attendance

- [ ] **Task 5: Unmarked Dashboard Component** (AC: 1,2) - UI deferred
  - [ ] Create unmarked attendance view
  - [ ] Highlight escalated (post-deadline) classes
  - [ ] Add notification to teachers

## Dev Agent Record

### Agent Model Used
Claude Opus 4.5 (claude-opus-4-5-20251101)

### Completion Notes

1. Unmarked attendance endpoint tracks classes without daily attendance
2. Post-deadline escalation based on configured deadline time
3. Teacher information included for accountability
4. Settings enforcement integrated with existing settings system
5. Frontend models and service methods added
6. UI components deferred - backend APIs fully functional

### Change Log
| Date | Change | Author |
|------|--------|--------|
| 2026-01-29 | Story created and implemented | dev-story workflow |

### File List

**Backend - Module (modified):**
- `msls-backend/internal/modules/studentattendance/dto.go` - Added UnmarkedAttendanceResponse, UnmarkedClassInfoDTO
- `msls-backend/internal/modules/studentattendance/repository.go` - Added GetUnmarkedClasses
- `msls-backend/internal/modules/studentattendance/service.go` - Added GetUnmarkedAttendance
- `msls-backend/internal/modules/studentattendance/handler.go` - Added GetUnmarkedAttendance handler
- `msls-backend/cmd/api/main.go` - Added unmarked alert route

**Frontend - Feature (modified):**
- `msls-frontend/src/app/features/student-attendance/student-attendance.model.ts` - Added UnmarkedAttendance interface
- `msls-frontend/src/app/features/student-attendance/student-attendance.service.ts` - Added getUnmarkedAttendance method
