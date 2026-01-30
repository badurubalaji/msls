# Story 8.2: Examination Creation & Scheduling

Status: in-progress

## Story

As an **academic coordinator**,
I want **to create examinations with schedules**,
So that **exams are properly organized**.

## Acceptance Criteria

### AC1: Create Examination
**Given** coordinator is creating an examination
**When** filling exam details
**Then** they can enter: exam name, type, academic year
**And** they can select: applicable classes
**And** they can set: date range (start and end dates)
**And** exam is created in draft status

### AC2: Create Exam Schedule
**Given** an exam exists
**When** creating schedule
**Then** they can add: date, subject, start time, end time
**And** they can set: maximum marks, passing marks
**And** they can assign: exam venue/room
**And** schedule validates no conflicts (date/room)

### AC3: Publish Exam
**Given** schedule is complete
**When** publishing the exam
**Then** exam status changes to "scheduled"
**And** teachers and students can see the schedule
**And** calendar is updated with exam dates

## Tasks / Subtasks

### Backend Tasks

- [ ] **Task 1: Database Schema** (AC: 1,2,3)
  - [ ] Create `examinations` table: id, tenant_id, name, exam_type_id, academic_year_id, start_date, end_date, status (draft/scheduled/ongoing/completed/cancelled), description
  - [ ] Create `examination_classes` table: examination_id, class_id (many-to-many)
  - [ ] Create `exam_schedules` table: id, examination_id, subject_id, exam_date, start_time, end_time, max_marks, passing_marks, room/venue
  - [ ] Add migration file with RLS policies and indexes
  - [ ] Add status transition constraints

- [ ] **Task 2: Model & Repository** (AC: 1,2,3)
  - [ ] Create `Examination` model with relationships
  - [ ] Create `ExamSchedule` model
  - [ ] Implement repository CRUD for examinations
  - [ ] Implement schedule management methods
  - [ ] Add conflict validation for schedules

- [ ] **Task 3: Service Layer** (AC: 1,2,3)
  - [ ] Create exam service with business logic
  - [ ] Validate date ranges and schedule conflicts
  - [ ] Implement status transitions (draft -> scheduled -> ongoing -> completed)
  - [ ] Ensure exam has schedules before publishing

- [ ] **Task 4: API Endpoints** (AC: 1,2,3)
  - [ ] `GET /api/v1/examinations` - List examinations (filters: academic_year, status, class)
  - [ ] `GET /api/v1/examinations/:id` - Get examination with schedules
  - [ ] `POST /api/v1/examinations` - Create examination
  - [ ] `PUT /api/v1/examinations/:id` - Update examination
  - [ ] `DELETE /api/v1/examinations/:id` - Delete (only draft)
  - [ ] `POST /api/v1/examinations/:id/publish` - Publish examination
  - [ ] `POST /api/v1/examinations/:id/schedules` - Add schedule
  - [ ] `PUT /api/v1/examinations/:id/schedules/:scheduleId` - Update schedule
  - [ ] `DELETE /api/v1/examinations/:id/schedules/:scheduleId` - Delete schedule

- [ ] **Task 5: Permissions** (AC: 1,2,3)
  - [ ] Add permissions: `exam:view`, `exam:create`, `exam:update`, `exam:delete`, `exam:publish`

### Frontend Tasks

- [ ] **Task 6: Model & Service** (AC: 1,2,3)
  - [ ] Create `Examination` and `ExamSchedule` interfaces
  - [ ] Create `ExaminationService` with API methods
  - [ ] Add exam status constants and helpers

- [ ] **Task 7: Examination List Component** (AC: 1,3)
  - [ ] Create list page with status tabs/filters
  - [ ] Show exam name, type, dates, class count, status
  - [ ] Add actions: view, edit, delete, publish
  - [ ] Add create button

- [ ] **Task 8: Examination Form/Detail Component** (AC: 1,2,3)
  - [ ] Create form for exam details (name, type, classes, dates)
  - [ ] Show schedule grid/list
  - [ ] Add schedule management (add/edit/delete)
  - [ ] Validate conflicts before saving
  - [ ] Add publish button with confirmation

- [ ] **Task 9: Routing & Navigation** (AC: 1,2,3)
  - [ ] Add examination routes under /exams
  - [ ] Add "Examinations" to navigation
  - [ ] Add permission guards

## Dev Notes

### Database Design
```sql
CREATE TABLE examinations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    name VARCHAR(200) NOT NULL,
    exam_type_id UUID NOT NULL REFERENCES exam_types(id),
    academic_year_id UUID NOT NULL REFERENCES academic_years(id),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    CONSTRAINT chk_exam_dates CHECK (end_date >= start_date),
    CONSTRAINT chk_exam_status CHECK (status IN ('draft', 'scheduled', 'ongoing', 'completed', 'cancelled'))
);

CREATE TABLE examination_classes (
    examination_id UUID NOT NULL REFERENCES examinations(id) ON DELETE CASCADE,
    class_id UUID NOT NULL REFERENCES classes(id),
    PRIMARY KEY (examination_id, class_id)
);

CREATE TABLE exam_schedules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    examination_id UUID NOT NULL REFERENCES examinations(id) ON DELETE CASCADE,
    subject_id UUID NOT NULL REFERENCES subjects(id),
    exam_date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    max_marks INTEGER NOT NULL DEFAULT 100,
    passing_marks INTEGER,
    venue VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_schedule_times CHECK (end_time > start_time)
);
```

### Status Flow
```
draft -> scheduled (when published)
scheduled -> ongoing (when start_date reached)
ongoing -> completed (when end_date passed or manual)
any -> cancelled (by admin)
```

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
