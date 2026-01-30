# Story 7.2: Period-wise Attendance

Status: review

## Story

As a **teacher**,
I want **to mark period-wise attendance**,
so that **subject-specific presence is tracked**.

## Acceptance Criteria

### AC1: Period Selection
**Given** period-wise attendance is enabled for the school (via settings)
**When** a teacher accesses the attendance marking page
**Then** they can select a specific period slot from their timetable
**And** the marking is scoped to that period only
**And** daily summary is available showing all periods

### AC2: Period Attendance Marking
**Given** a teacher has selected a section and period
**When** marking attendance
**Then** they use the same P/A/L/H status options as daily attendance
**And** each record is linked to the specific period (via timetable_entry_id)
**And** the subject being taught is displayed

### AC3: Daily Summary Aggregation
**Given** a student has period-wise attendance for a day
**When** viewing their attendance
**Then** they see period-wise status for each period that day
**And** total periods present/absent is calculated
**And** overall day status is derived (present if >50% periods attended)

### AC4: Subject-wise Attendance Analytics
**Given** period attendance is marked consistently
**When** viewing subject-wise attendance
**Then** subject-wise attendance percentage is available
**And** minimum attendance for exam eligibility is trackable (e.g., 75%)
**And** alerts for low subject attendance can be generated

## Tasks / Subtasks

### Backend Tasks

- [x] **Task 1: Database Schema Extension** (AC: 1,2,3)
  - [x] Create migration `000054_period_attendance.up.sql`:
    - Add `period_id` column (nullable UUID) to `student_attendance` table referencing `period_slots(id)`
    - Add `timetable_entry_id` column (nullable UUID) referencing `timetable_entries(id)` for subject context
    - Add index on (section_id, attendance_date, period_id)
    - Modify unique constraint to allow multiple records per student per day (one per period)
  - [x] Create down migration
  - [x] Add `period_attendance_enabled` setting to `student_attendance_settings` table

- [x] **Task 2: Backend Model Updates** (AC: 1,2,3)
  - [x] Update `internal/pkg/database/models/student_attendance.go`:
    - Add `PeriodID *uuid.UUID` field (nullable for backward compatibility with daily attendance)
    - Add `TimetableEntryID *uuid.UUID` field
    - Add relationship to PeriodSlot and TimetableEntry models

- [x] **Task 3: Service Layer Enhancements** (AC: 1,2,3,4)
  - [x] Update `internal/modules/studentattendance/service.go`:
    - Add `GetTeacherPeriods(ctx, teacherID, sectionID, date)` - periods for section on date
    - Add `GetPeriodAttendance(ctx, sectionID, periodID, date)` - students with period attendance
    - Update `MarkClassAttendance` to handle period-specific marking
    - Add `GetDailySummary(ctx, sectionID, date)` - aggregate all periods for the day
    - Add `GetSubjectAttendance(ctx, studentID, subjectID, dateRange)` - subject analytics
  - [x] Update repository with new queries

- [x] **Task 4: New API Endpoints** (AC: 1,2,3,4)
  - [x] Add to handler.go:
    - `GET /api/v1/student-attendance/periods?section_id=X&date=Y` - periods for section
    - `GET /api/v1/student-attendance/period/:periodId?section_id=X&date=Y` - period attendance
    - `POST /api/v1/student-attendance/period/:periodId` - mark period attendance
    - `GET /api/v1/student-attendance/daily-summary?section_id=X&date=Y` - day's all periods
    - `GET /api/v1/student-attendance/subject/:subjectId?student_id=X` - subject analytics

### Frontend Tasks

- [x] **Task 5: Model & Service Updates** (AC: 1,2,3,4)
  - [x] Update `student-attendance.model.ts`:
    - Add `PeriodAttendance` interface with period and subject info
    - Add `DailySummary` interface for aggregated view
    - Add `SubjectAttendance` interface for analytics
  - [x] Update `student-attendance.service.ts`:
    - Add `getTeacherPeriods(sectionId, date)` method
    - Add `getPeriodAttendance(sectionId, periodId, date)` method
    - Add `markPeriodAttendance(sectionId, periodId, request)` method
    - Add `getDailySummary(sectionId, date)` method
    - Add `getSubjectAttendance(studentId, subjectId)` method

- [x] **Task 6: Period Attendance Component** (AC: 1,2,3)
  - [x] Create `features/student-attendance/pages/period-attendance/period-attendance.component.ts`:
    - Section selector (same as daily)
    - Period selector dropdown showing today's periods from timetable
    - Student grid with same P/A/L/H toggles
    - Display subject being taught for selected period
    - Quick switch between periods for same section

- [x] **Task 7: Daily Summary View Component** (AC: 3)
  - [x] Create `features/student-attendance/pages/daily-summary/daily-summary.component.ts`:
    - Grid showing all students × all periods for the day
    - Each cell shows P/A/L/H status
    - Row totals: periods attended, percentage
    - Highlight students with low attendance (<50%)

- [x] **Task 8: Routes & Navigation** (AC: 1)
  - [x] Update `student-attendance.routes.ts`:
    - Add route `/student-attendance/period` for period-wise marking
    - Add route `/student-attendance/summary` for daily summary
  - [x] Update nav-config to add "Period-wise" option under Attendance section

## Dev Notes

### Architecture Patterns to Follow

**Builds on Story 7.1:**
- Reuse existing `studentattendance` module structure
- Extend existing tables rather than create new ones
- Same 4-layer architecture: Handler → Service → Repository → Model
- Same permission model, may need new permission `student_attendance:mark_period`

### Database Schema Changes

```sql
-- Add period-wise attendance support
ALTER TABLE student_attendance
  ADD COLUMN period_id UUID REFERENCES period_slots(id),
  ADD COLUMN timetable_entry_id UUID REFERENCES timetable_entries(id);

-- Drop old unique constraint (tenant_id, student_id, attendance_date)
ALTER TABLE student_attendance
  DROP CONSTRAINT uniq_student_daily_attendance;

-- Add new constraint allowing multiple per day if period_id differs
ALTER TABLE student_attendance
  ADD CONSTRAINT uniq_student_period_attendance
    UNIQUE (tenant_id, student_id, attendance_date, period_id);

-- Add index for period queries
CREATE INDEX idx_student_attendance_period
  ON student_attendance(section_id, attendance_date, period_id);

-- Add setting for period-wise attendance
ALTER TABLE student_attendance_settings
  ADD COLUMN period_attendance_enabled BOOLEAN NOT NULL DEFAULT false;
```

### Backward Compatibility

- Daily attendance (Story 7.1) continues to work with `period_id = NULL`
- Schools can choose daily-only or period-wise attendance via settings
- Existing records remain valid with NULL period_id
- New constraint uses COALESCE or allows NULL in unique constraint

### API Endpoints (New)

| Method | Endpoint | Description | Permission |
|--------|----------|-------------|------------|
| GET | `/api/v1/student-attendance/periods?section_id=X&date=Y` | Get periods for section | `student_attendance:mark_class` |
| GET | `/api/v1/student-attendance/period/:id?section_id=X&date=Y` | Get period attendance | `student_attendance:view_class` |
| POST | `/api/v1/student-attendance/period/:id` | Mark period attendance | `student_attendance:mark_class` |
| GET | `/api/v1/student-attendance/daily-summary?section_id=X&date=Y` | Daily period summary | `student_attendance:view_class` |
| GET | `/api/v1/student-attendance/subject/:id?student_id=X` | Subject attendance stats | `student_attendance:view_class` |

### Frontend Component Structure

```
src/app/features/student-attendance/
├── student-attendance.routes.ts (updated)
├── student-attendance.model.ts (updated)
├── student-attendance.service.ts (updated)
└── pages/
    ├── mark-attendance/           (existing from 7.1)
    ├── period-attendance/         (new)
    │   └── period-attendance.component.ts
    └── daily-summary/             (new)
        └── daily-summary.component.ts
```

### Key Technical Considerations

1. **Timetable Integration**: Period attendance requires timetable to be set up (Story 6.4)
2. **Subject Tracking**: Link to `timetable_entries` to know which subject for each period
3. **Aggregation Logic**: Daily status derived from period attendance (>50% = present)
4. **Performance**: Index on (section_id, attendance_date, period_id) for quick lookups
5. **Feature Toggle**: Use `period_attendance_enabled` in settings to control availability

### Previous Story Learnings (Story 7.1)

From Story 7.1 implementation:
1. Used `sections` table (not `class_sections`) - follow same pattern
2. Created custom trigger functions for updated_at columns
3. Frontend uses Angular Signals for reactive state
4. Keyboard shortcuts (P/A/L/H) should work in period view too
5. Module structure: `internal/modules/studentattendance/`

### Testing Requirements

- Verify daily attendance (7.1) still works after schema changes
- Test period attendance with multiple periods per day
- Test daily summary aggregation logic
- Test subject-wise analytics calculation
- Frontend: test period selector and grid interactions

### Project Structure Notes

- Extend existing `studentattendance` module (don't create new module)
- Reuse existing models and add new fields
- Follow existing patterns in `internal/modules/`
- Reuse `response` and `apperrors` packages

### References

- [Source: epic-07-attendance-operations.md#Story 7.2]
- [Source: project-context.md#Go Backend Rules]
- [Source: project-context.md#Angular Frontend Rules]
- [Pattern: internal/modules/studentattendance/ - from Story 7.1]
- [Pattern: internal/modules/timetable/ - for period slot integration]
- [Migration: 000053_student_attendance.up.sql - existing schema]

## Dev Agent Record

### Agent Model Used

Claude Opus 4.5 (claude-opus-4-5-20251101)

### Completion Notes List

1. **Database Schema Extension** - Created migration 000054 with:
   - Added `period_id` and `timetable_entry_id` columns to `student_attendance`
   - Used COALESCE in unique constraint to handle NULL period_id (daily attendance compatibility)
   - Added indexes for period-wise queries
   - Added `period_attendance_enabled` setting

2. **Backend Implementation** - Extended studentattendance module:
   - Added PeriodID, TimetableEntryID fields with GORM relationships
   - Added repository methods for period attendance CRUD and aggregation
   - Added service methods: GetTeacherPeriods, GetPeriodAttendance, MarkPeriodAttendance, GetDailySummary, GetSubjectAttendance
   - Added 5 new API endpoints with proper error handling

3. **Frontend Implementation**:
   - Extended models with PeriodInfo, PeriodAttendance, DailySummary, SubjectAttendanceStats interfaces
   - Extended service with 5 new methods for period attendance operations
   - Created period-attendance component with section/period selectors, student grid, keyboard shortcuts (P/A/L/H)
   - Created daily-summary component with student × period grid, percentage calculations, low attendance highlighting

4. **Integration Notes**:
   - Backward compatible: daily attendance (Story 7.1) continues with period_id = NULL
   - Integrated with timetable module for period slot and entry lookups
   - Daily status derived from period attendance (>50% = present)
   - Subject-wise analytics support 75% minimum for exam eligibility tracking

### Change Log

| Date | Change | Author |
|------|--------|--------|
| 2026-01-29 | Story created with comprehensive context | create-story workflow |
| 2026-01-29 | Implemented all 8 tasks - backend and frontend complete | dev-story workflow |

### File List

**Backend Files Created/Modified:**
- `msls-backend/migrations/000054_period_attendance.up.sql` (new)
- `msls-backend/migrations/000054_period_attendance.down.sql` (new)
- `msls-backend/internal/pkg/database/models/student_attendance.go` (modified)
- `msls-backend/internal/modules/studentattendance/errors.go` (modified)
- `msls-backend/internal/modules/studentattendance/dto.go` (modified)
- `msls-backend/internal/modules/studentattendance/repository.go` (modified)
- `msls-backend/internal/modules/studentattendance/service.go` (modified)
- `msls-backend/internal/modules/studentattendance/handler.go` (modified)
- `msls-backend/cmd/api/main.go` (modified - route registration)

**Frontend Files Created/Modified:**
- `msls-frontend/src/app/features/student-attendance/student-attendance.model.ts` (modified)
- `msls-frontend/src/app/features/student-attendance/student-attendance.service.ts` (modified)
- `msls-frontend/src/app/features/student-attendance/pages/period-attendance/period-attendance.component.ts` (new)
- `msls-frontend/src/app/features/student-attendance/pages/daily-summary/daily-summary.component.ts` (new)
- `msls-frontend/src/app/features/student-attendance/student-attendance.routes.ts` (modified)
- `msls-frontend/src/app/layouts/nav-config.ts` (modified)
