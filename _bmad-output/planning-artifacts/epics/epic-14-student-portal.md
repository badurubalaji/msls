# Epic 14: Student Portal

**Phase:** 1 (MVP)
**Priority:** High - Student interface

## Epic Goal

Enable students to access their academic information and submit assignments.

## User Value

Students can view their timetable, attendance, submit homework, and check exam results.

## FRs Covered

FR-SP-01 to FR-SP-07

---

## Stories

### Story 14.1: Student Login

As a **student**,
I want **to login to my student portal**,
So that **I can access my academic information**.

**Acceptance Criteria:**

**Given** a student has credentials
**When** logging in
**Then** they can use: admission number + password
**And** they can use: phone OTP (for senior students)
**And** first login requires password change
**And** age-appropriate access controls applied

**Given** junior student (class 1-5)
**When** accessing
**Then** simplified interface is shown
**And** limited features (view only)
**And** password managed by parent/school

---

### Story 14.2: Student Dashboard

As a **student**,
I want **to see my dashboard**,
So that **I know what's happening today**.

**Acceptance Criteria:**

**Given** a student logs in
**When** viewing dashboard
**Then** they see: today's timetable
**And** they see: pending homework count
**And** they see: upcoming exams/events
**And** they see: recent notices

**Given** quick actions
**When** displayed
**Then** they can: view full timetable
**And** they can: go to homework
**And** they can: view results
**And** interface is engaging/student-friendly

---

### Story 14.3: Timetable View

As a **student**,
I want **to view my class timetable**,
So that **I know my daily schedule**.

**Acceptance Criteria:**

**Given** student is on timetable
**When** viewing
**Then** they see: weekly timetable grid
**And** they see: period times, subjects, teachers
**And** current period is highlighted
**And** breaks and special periods shown

**Given** mobile view
**When** on phone
**Then** daily view is default
**And** swipe between days
**And** today button to jump back

---

### Story 14.4: Attendance View

As a **student**,
I want **to view my attendance**,
So that **I know my attendance status**.

**Acceptance Criteria:**

**Given** student is on attendance
**When** viewing
**Then** they see: calendar with attendance marked
**And** they see: summary (days, %, comparison)
**And** color coding clear and visible

**Given** low attendance
**When** below threshold
**Then** warning message is shown
**And** days needed to reach threshold calculated
**And** impact on exam eligibility shown

---

### Story 14.5: Homework Submission

As a **student**,
I want **to view and submit homework online**,
So that **I can complete my assignments**.

**Acceptance Criteria:**

**Given** student is on homework
**When** viewing assignments
**Then** they see: pending, submitted, graded tabs
**And** they see: due dates, subjects, status
**And** overdue assignments highlighted

**Given** submitting homework
**When** on assignment page
**Then** they can: view full instructions
**And** they can: download attachments
**And** they can: upload their submission
**And** they can: add text response
**And** submission confirmation shown

**Given** graded homework
**When** viewing
**Then** they see: marks obtained
**And** they see: teacher feedback
**And** they can view: their submitted work

---

### Story 14.6: Exam Results & Report Cards

As a **student**,
I want **to view my exam results**,
So that **I know how I performed**.

**Acceptance Criteria:**

**Given** student is on results
**When** viewing
**Then** they see: list of exams with results
**And** they see: subject-wise marks, grades
**And** they see: rank (if displayed)
**And** they see: pass/fail status

**Given** report card is available
**When** downloading
**Then** they can download PDF
**And** report shows: all subjects, grades, remarks
**And** historical report cards available

---

### Story 14.7: Quiz & Test View

As a **student**,
I want **to see assigned quizzes and tests**,
So that **I can attempt them**.

**Acceptance Criteria:**

**Given** student is on quizzes (Phase 1 view only)
**When** viewing
**Then** they see: upcoming quizzes
**And** they see: past quiz results
**And** link to attempt quiz (full feature in Phase 2)

---

### Story 14.8: Library Books View

As a **student**,
I want **to see my borrowed library books**,
So that **I know what to return**.

**Acceptance Criteria:**

**Given** student is on library
**When** viewing
**Then** they see: currently borrowed books
**And** they see: due dates
**And** they see: overdue books highlighted
**And** they see: borrowing history

**Given** book search (Phase 2 feature preview)
**When** displayed
**Then** "Coming soon" placeholder shown
**And** basic catalog search link to library

---

### Story 14.9: Notices & Events

As a **student**,
I want **to view notices and upcoming events**,
So that **I'm informed of school activities**.

**Acceptance Criteria:**

**Given** student is on notices
**When** viewing
**Then** they see: relevant notices (student-targeted)
**And** they see: upcoming events calendar
**And** event details on click
**And** exam schedule highlighted
