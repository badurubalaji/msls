# Epic 7: Daily Attendance Operations

**Phase:** 1 (MVP)
**Priority:** High - Daily operational requirement

## Epic Goal

Enable teachers to efficiently mark and track student attendance with comprehensive reporting.

## User Value

Teachers can mark attendance quickly, parents receive absence notifications, and admins can generate attendance reports.

## FRs Covered

FR-AO-06 to FR-AO-09

---

## Stories

### Story 7.1: Daily Attendance Marking Interface

As a **teacher**,
I want **to mark student attendance for my class**,
So that **daily presence records are maintained**.

**Acceptance Criteria:**

**Given** a teacher is logged in
**When** accessing attendance marking
**Then** they see their assigned classes for today
**And** selecting a class shows student list
**And** default status is "present" for all students

**Given** the attendance grid is displayed
**When** marking attendance
**Then** they can toggle: Present (P), Absent (A), Late (L), Half-day (H)
**And** grid shows student photo and name
**And** previous attendance indicator shows pattern
**And** bulk actions available (mark all present, mark all absent)

**Given** attendance is marked
**When** submitting
**Then** attendance is saved with timestamp
**And** late arrivals are recorded with time
**And** confirmation message is shown
**And** SMS is triggered to absent students' parents

---

### Story 7.2: Period-wise Attendance (Optional)

As a **teacher**,
I want **to mark period-wise attendance**,
So that **subject-specific presence is tracked**.

**Acceptance Criteria:**

**Given** period-wise attendance is enabled for the school
**When** a teacher marks attendance
**Then** they select the specific period slot
**And** marking is for that period only
**And** daily summary aggregates all periods

**Given** a student is absent for specific periods
**When** viewing their attendance
**Then** they see: period-wise status for the day
**And** total periods present/absent is calculated
**And** subject-wise attendance percentage is available

**Given** period attendance is marked
**When** calculating percentages
**Then** subject-wise attendance is calculated
**And** minimum attendance for exam eligibility is tracked
**And** alerts for low subject attendance are generated

---

### Story 7.3: Attendance Edit & History

As a **teacher or admin**,
I want **to edit attendance within allowed window**,
So that **errors can be corrected**.

**Acceptance Criteria:**

**Given** attendance was marked
**When** editing within edit window (configurable, default 24 hours)
**Then** original marker can edit the attendance
**And** reason for edit must be provided
**And** edit history is recorded

**Given** edit window has passed
**When** attendance needs correction
**Then** admin approval is required
**And** request goes to class teacher/admin
**And** approved changes are recorded with justification

**Given** attendance history is needed
**When** viewing audit trail
**Then** all changes are logged: who, when, what changed, reason
**And** original and modified values are shown
**And** history is immutable

---

### Story 7.4: Student Attendance Calendar View

As a **parent or student**,
I want **to view attendance in calendar format**,
So that **attendance pattern is easily visible**.

**Acceptance Criteria:**

**Given** a parent/student is logged in
**When** viewing attendance
**Then** they see monthly calendar view
**And** each day is color-coded: green (present), red (absent), yellow (late), gray (holiday)
**And** clicking a day shows details

**Given** attendance summary is needed
**When** viewing summary section
**Then** they see: total days, present, absent, late, percentage
**And** comparison with class average is shown
**And** trend (improving/declining) is indicated

---

### Story 7.5: Attendance Reports - Class Level

As a **teacher or admin**,
I want **to generate class-level attendance reports**,
So that **class attendance can be analyzed**.

**Acceptance Criteria:**

**Given** a teacher is on attendance reports
**When** selecting daily report
**Then** they see: all students with status for selected date
**And** summary shows: total, present, absent, percentage
**And** absentee list is highlighted

**Given** monthly report is selected
**When** generating the report
**Then** grid shows: students (rows) × dates (columns)
**And** each cell shows status symbol
**And** summary column shows: total present, percentage
**And** export to Excel is available

**Given** comparative analysis is needed
**When** viewing class comparison
**Then** sections are compared for attendance %
**And** below-threshold sections are highlighted
**And** trend over weeks/months is displayed

---

### Story 7.6: Attendance Reports - Individual Student

As a **teacher or parent**,
I want **to view individual student attendance report**,
So that **specific student's attendance is analyzed**.

**Acceptance Criteria:**

**Given** a student is selected
**When** viewing their attendance report
**Then** they see: monthly breakdown with daily status
**And** summary shows: working days, present, absent, late, %
**And** subject-wise breakdown (if period-wise enabled)

**Given** printable report is needed
**When** generating certificate
**Then** formal attendance certificate is generated
**And** includes: student details, period, attendance summary
**And** principal signature placeholder is included
**And** PDF download is available

---

### Story 7.7: Low Attendance Alerts

As an **administrator**,
I want **to receive alerts for low attendance students**,
So that **intervention can happen early**.

**Acceptance Criteria:**

**Given** attendance thresholds are configured
**When** a student falls below threshold (e.g., 75%)
**Then** alert is generated for class teacher
**And** student appears in "attention needed" list
**And** parents receive notification about low attendance

**Given** attendance dashboard exists
**When** admin views dashboard
**Then** they see: students below threshold count
**And** they see: chronic absentees (below 60%)
**And** they see: attendance trend graph
**And** drill-down to individual student is available

**Given** daily summary is needed
**When** end of day
**Then** automated report shows: overall attendance %
**And** classes with low attendance are highlighted
**And** report is emailed to principal/admin

---

### Story 7.8: Attendance Deadline Enforcement

As an **administrator**,
I want **to enforce attendance marking deadlines**,
So that **attendance is recorded timely**.

**Acceptance Criteria:**

**Given** attendance deadline is configured (e.g., 10:00 AM)
**When** deadline passes
**Then** attendance marking is locked for that period
**And** late marking requires admin override
**And** warning is shown before deadline

**Given** attendance not marked by deadline
**When** viewing unmarked report
**Then** admin sees: classes with missing attendance
**And** notification is sent to respective teachers
**And** escalation to HOD if still unmarked

**Given** holiday or special day
**When** marked in calendar
**Then** attendance marking is disabled
**And** day is marked as non-working
**And** percentage calculation excludes holiday
