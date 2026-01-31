# Epic 8: Examinations & Grading

**Phase:** 1 (MVP)
**Priority:** High - Core academic feature

## Epic Goal

Enable complete examination management from scheduling to report card generation.

## User Value

Teachers can schedule exams, enter marks, and generate report cards; students and parents can view results.

## FRs Covered

FR-EX-01 to FR-EX-09

---

## Stories

### Story 8.1: Exam Type Configuration

As an **academic administrator**,
I want **to configure different exam types**,
So that **various assessments can be created**.

**Acceptance Criteria:**

**Given** admin is on exam settings
**When** creating an exam type
**Then** they can enter: name (e.g., "Unit Test", "Half Yearly", "Annual")
**And** they can set: weightage for final result calculation
**And** they can set: is marks-based or grade-based
**And** they can set: default maximum marks

**Given** exam types are configured
**When** viewing the list
**Then** they see all exam types with settings
**And** they can order exam types (display sequence)
**And** they can activate/deactivate types

---

### Story 8.2: Examination Creation & Scheduling

As an **academic coordinator**,
I want **to create examinations with schedules**,
So that **exams are properly organized**.

**Acceptance Criteria:**

**Given** coordinator is creating an examination
**When** filling exam details
**Then** they can enter: exam name, type, academic year
**And** they can select: applicable classes
**And** they can set: date range (start and end dates)
**And** exam is created in draft status

**Given** an exam exists
**When** creating schedule
**Then** they can add: date, subject, start time, end time
**And** they can set: maximum marks, passing marks
**And** they can assign: exam venue/room
**And** schedule validates no conflicts (date/room)

**Given** schedule is complete
**When** publishing the exam
**Then** exam status changes to "scheduled"
**And** teachers and students can see the schedule
**And** calendar is updated with exam dates

---

### Story 8.3: Hall Ticket Generation

As an **administrator**,
I want **to generate hall tickets for exams**,
So that **students have proper exam identification**.

**Acceptance Criteria:**

**Given** an exam is scheduled
**When** generating hall tickets
**Then** they can select: exam, class, section (or all)
**And** hall tickets are generated in batch
**And** each ticket has: unique roll number, student photo, exam schedule

**Given** hall ticket template exists
**When** generating tickets
**Then** template is used with school branding
**And** exam-wise schedule is printed
**And** important instructions are included
**And** QR code for verification is included

**Given** hall tickets are generated
**When** distributing
**Then** bulk print option is available
**And** PDF download for individual or batch
**And** student portal shows downloadable hall ticket

---

### Story 8.4: Marks Entry Interface

As a **teacher**,
I want **to enter marks for my subject**,
So that **student performance is recorded**.

**Acceptance Criteria:**

**Given** a teacher is assigned to a subject
**When** accessing marks entry
**Then** they see: exams for which marks can be entered
**And** selecting exam/class shows student list
**And** marks entry grid is displayed

**Given** marks entry grid
**When** entering marks
**Then** they can enter marks per student
**And** validation prevents exceeding maximum marks
**And** absent students can be marked as "AB"
**And** auto-save drafts as they type
**And** submit button finalizes the entry

**Given** marks are entered
**When** submitting
**Then** marks are locked from further edits
**And** timestamp and teacher recorded
**And** re-entry requires unlock by admin

---

### Story 8.5: Grading Scale Configuration

As an **academic administrator**,
I want **to configure grading scales**,
So that **marks are converted to grades correctly**.

**Acceptance Criteria:**

**Given** admin is on grading settings
**When** creating a grading scale
**Then** they can enter: scale name (e.g., "CBSE Grading")
**And** they can define grade ranges: A1 (91-100), A2 (81-90), etc.
**And** they can set: grade points for each grade
**And** they can set: remarks (Excellent, Good, etc.)

**Given** grading scales exist
**When** assigning to classes
**Then** they can map scale to class or exam type
**And** different scales can be used for different levels
**And** subject-specific scales can be configured

---

### Story 8.6: Result Calculation & Processing

As an **academic coordinator**,
I want **to calculate and process results**,
So that **grades and ranks are generated**.

**Acceptance Criteria:**

**Given** all marks are entered for an exam
**When** processing results
**Then** percentage is calculated per student per subject
**And** grades are assigned based on grading scale
**And** pass/fail status is determined
**And** result status is set: pass, fail, compartment

**Given** result calculation is complete
**When** generating ranks
**Then** class rank is calculated (by total marks/percentage)
**And** section rank is calculated
**And** subject topper is identified
**And** ties are handled (same rank for equal marks)

**Given** results are processed
**When** publishing results
**Then** results are visible to students/parents
**And** result publication date is recorded
**And** notification sent to parents

---

### Story 8.7: Report Card Generation

As an **academic administrator**,
I want **to generate report cards**,
So that **students receive formal performance documents**.

**Acceptance Criteria:**

**Given** results are published
**When** generating report cards
**Then** they can select: exam, class, section
**And** report cards are generated from template
**And** each card shows: student info, subject-wise marks/grades, total, rank

**Given** report card template
**When** customizing
**Then** school logo and branding is included
**And** co-scholastic grades section is available
**And** remarks section for class teacher is included
**And** signature placeholders for teacher/principal

**Given** report cards are ready
**When** distributing
**Then** bulk print is available (multiple per page option)
**And** PDF download per student is available
**And** student/parent portal shows downloadable report card

---

### Story 8.8: Result Analytics & Reports

As a **principal or coordinator**,
I want **to view result analytics**,
So that **academic performance can be analyzed**.

**Acceptance Criteria:**

**Given** results are published
**When** viewing analytics
**Then** they see: pass percentage by class/section/subject
**And** they see: grade distribution (how many A1, A2, etc.)
**And** they see: comparison with previous exams

**Given** subject-wise analysis
**When** drilling down
**Then** they see: average marks, highest, lowest
**And** they see: questions where students struggled (if available)
**And** they see: teacher-wise comparison

**Given** export is needed
**When** generating reports
**Then** detailed result register is available
**And** promotion list is generated
**And** merit list is generated
**And** export to Excel/PDF available

---

### Story 8.9: Exam Eligibility Based on Attendance

As an **academic administrator**,
I want **to check exam eligibility based on attendance**,
So that **attendance requirements are enforced**.

**Acceptance Criteria:**

**Given** minimum attendance is configured (e.g., 75%)
**When** checking eligibility before exam
**Then** system generates list of ineligible students
**And** list shows: student, attendance %, shortfall
**And** admin can approve exceptions

**Given** a student is ineligible
**When** marks entry is attempted
**Then** warning is shown about ineligibility
**And** marks can still be entered (for record)
**And** result shows "detained" or "ineligible"

**Given** eligibility report is needed
**When** generating
**Then** report shows: eligible vs ineligible count
**And** class-wise breakdown is available
**And** report can be shared with parents

---

### Story 8.10: Replace SQLite Test Database with PostgreSQL

As a **developer**,
I want **unit tests to use PostgreSQL instead of SQLite**,
So that **tests accurately reflect production behavior**.

**Acceptance Criteria:**

**Given** tests use in-memory SQLite currently
**When** running unit tests
**Then** PostgreSQL-specific syntax fails (e.g., `NULLS LAST`, RLS)
**And** test results don't reflect production behavior

**Given** testcontainers-go is configured
**When** running unit tests
**Then** a real PostgreSQL container is spun up
**And** tests use actual PostgreSQL syntax
**And** RLS policies can be tested
**And** UUID generation works correctly

**Given** test infrastructure is updated
**When** running `go test ./...`
**Then** all existing tests pass with PostgreSQL
**And** CI/CD pipeline uses the same approach
**And** tests are isolated per test run

**Technical Notes:**
- Remove `gorm.io/driver/sqlite` dependency
- Add `github.com/testcontainers/testcontainers-go`
- Create shared test helper for PostgreSQL container setup
- Update all `*_test.go` files using SQLite
- Ensure parallel test execution works correctly
