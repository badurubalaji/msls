# Epic 9: Homework & Assignments

**Phase:** 1 (MVP)
**Priority:** High - Daily academic operations

## Epic Goal

Enable teachers to assign homework and track student submissions with grading.

## User Value

Teachers can create assignments, students can submit online, and parents are notified of pending work.

## FRs Covered

FR-HW-01 to FR-HW-07

---

## Stories

### Story 9.1: Assignment Creation

As a **teacher**,
I want **to create assignments for my class**,
So that **students have clear work to complete**.

**Acceptance Criteria:**

**Given** a teacher is on assignments page
**When** creating a new assignment
**Then** they can enter: title, description/instructions
**And** they can select: class, section, subject
**And** they can select: assignment type (homework, project, worksheet, classwork)
**And** they can set: due date and time
**And** they can set: maximum marks (or ungraded)

**Given** assignment details are entered
**When** setting submission options
**Then** they can choose: online submission allowed (yes/no)
**And** they can set: accepted file types (PDF, image, doc)
**And** they can set: max file size
**And** they can allow: late submission (with penalty configuration)

---

### Story 9.2: Assignment Attachments

As a **teacher**,
I want **to attach reference materials to assignments**,
So that **students have resources they need**.

**Acceptance Criteria:**

**Given** an assignment is being created
**When** adding attachments
**Then** they can upload: files (PDF, doc, images)
**And** they can add: links to external resources
**And** they can add: embedded YouTube videos
**And** multiple attachments can be added

**Given** attachments are added
**When** students view the assignment
**Then** they see all attachments clearly listed
**And** files can be downloaded
**And** links open in new tab
**And** videos play inline

---

### Story 9.3: Assignment Publication & Notification

As a **teacher**,
I want **to publish assignments and notify students**,
So that **students know about new work**.

**Acceptance Criteria:**

**Given** an assignment is created
**When** publishing the assignment
**Then** assignment status changes to "published"
**And** assignment becomes visible to students
**And** push notification sent to students
**And** parent notification includes assignment details

**Given** assignment schedule is needed
**When** setting publish date
**Then** they can schedule future publish date/time
**And** assignment auto-publishes at scheduled time
**And** notifications sent at publish time

---

### Story 9.4: Student Assignment View & Submission

As a **student**,
I want **to view and submit my assignments**,
So that **I can complete my work on time**.

**Acceptance Criteria:**

**Given** a student is logged in
**When** viewing assignments
**Then** they see: pending, submitted, graded assignments
**And** each shows: title, subject, due date, status
**And** overdue assignments are highlighted

**Given** a student opens an assignment
**When** viewing details
**Then** they see: full instructions, attachments
**And** they see: due date, submission status
**And** they see: their previous submission (if any)

**Given** a student is submitting
**When** uploading their work
**Then** they can upload files (within size limit)
**And** they can add text response
**And** they can save as draft
**And** they can submit when ready
**And** confirmation is shown after submission

---

### Story 9.5: Late Submission Handling

As a **teacher**,
I want **to configure late submission policies**,
So that **deadlines are enforced appropriately**.

**Acceptance Criteria:**

**Given** late submission is allowed for an assignment
**When** configuring penalty
**Then** they can set: grace period (hours/days)
**And** they can set: marks deduction per day late
**And** they can set: maximum late days allowed
**And** they can set: hard cutoff (no submissions after)

**Given** a student submits late
**When** submission is recorded
**Then** late flag is marked
**And** days late is calculated
**And** penalty is auto-calculated for grading
**And** original due date and actual submission date shown

---

### Story 9.6: Assignment Grading

As a **teacher**,
I want **to grade student submissions**,
So that **students receive feedback on their work**.

**Acceptance Criteria:**

**Given** submissions exist for an assignment
**When** accessing grading interface
**Then** teacher sees: list of all students with submission status
**And** submitted work can be viewed inline (PDF viewer)
**And** grading panel is beside the submission

**Given** grading a submission
**When** entering grade
**Then** they can enter: marks obtained (within max)
**And** they can enter: written feedback
**And** they can use: rubric if configured
**And** late penalty is auto-applied
**And** save grade for that student

**Given** bulk grading is needed
**When** using quick grade
**Then** they can enter marks in list view
**And** no feedback required (quick entry)
**And** all grades saved together

---

### Story 9.7: Assignment Status Tracking

As a **teacher**,
I want **to track assignment completion across class**,
So that **I know who has submitted**.

**Acceptance Criteria:**

**Given** an assignment is published
**When** viewing status
**Then** they see: total students, submitted, pending, graded
**And** they see: percentage completion
**And** they see: list of students who haven't submitted

**Given** reminder is needed
**When** sending reminder
**Then** they can send reminder to non-submitters
**And** notification sent via push/SMS
**And** reminder logged in assignment history

**Given** assignment deadline passed
**When** viewing overdue list
**Then** they see: students who missed deadline
**And** days overdue is shown
**And** bulk action to mark zero available

---

### Story 9.8: Parent Assignment Notifications

As a **parent**,
I want **to be notified of my child's assignments**,
So that **I can help ensure work is completed**.

**Acceptance Criteria:**

**Given** a new assignment is published
**When** notification is sent
**Then** parent receives: assignment title, subject, due date
**And** notification includes link to parent portal

**Given** assignment is approaching deadline
**When** 24 hours before due (if unsubmitted)
**Then** reminder sent to parent
**And** reminder shows: assignment name, time remaining

**Given** assignment is overdue
**When** child hasn't submitted
**Then** alert sent to parent
**And** alert shows: assignment details, days overdue
**And** parent can view in portal
