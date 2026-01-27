# Epic 13: Parent Portal

**Phase:** 1 (MVP)
**Priority:** High - Key stakeholder interface

## Epic Goal

Enable parents to access all relevant information about their children through a dedicated portal.

## User Value

Parents can view attendance, academics, fees, homework, notices, and communicate with teachers from one place.

## FRs Covered

FR-PP-01 to FR-PP-08

---

## Stories

### Story 13.1: Parent Registration & Login

As a **parent**,
I want **to register and login to the parent portal**,
So that **I can access my child's information**.

**Acceptance Criteria:**

**Given** a parent has a registered phone number
**When** accessing the portal
**Then** they can request OTP to their phone
**And** entering correct OTP logs them in
**And** first-time login prompts profile completion

**Given** parent profile setup
**When** completing profile
**Then** they can verify: linked children
**And** they can update: contact details
**And** they can set: notification preferences
**And** they can enable: email notifications

---

### Story 13.2: Parent Dashboard

As a **parent**,
I want **to see a dashboard with my child's summary**,
So that **I get quick overview at a glance**.

**Acceptance Criteria:**

**Given** a parent logs in
**When** viewing dashboard
**Then** they see: child selector (if multiple children)
**And** they see: today's attendance status
**And** they see: pending homework count
**And** they see: pending fee dues
**And** they see: recent notices (top 3)

**Given** dashboard widgets
**When** viewing
**Then** they see: attendance % for current month
**And** they see: next exam/event
**And** they see: quick actions (pay fee, view report card)
**And** widgets are tappable to see details

---

### Story 13.3: Attendance View

As a **parent**,
I want **to view my child's attendance**,
So that **I know their presence record**.

**Acceptance Criteria:**

**Given** parent is on attendance section
**When** viewing
**Then** they see: monthly calendar with attendance marked
**And** color coding: green (present), red (absent), yellow (late)
**And** they see: summary (total, present, absent, %)

**Given** detailed view is needed
**When** clicking on a date
**Then** they see: status for that day
**And** they see: late arrival time if applicable
**And** they see: reason if provided by school

---

### Story 13.4: Academic Performance View

As a **parent**,
I want **to view my child's academic performance**,
So that **I can track their progress**.

**Acceptance Criteria:**

**Given** parent is on academics section
**When** viewing exams
**Then** they see: list of exams with results
**And** they see: subject-wise marks and grades
**And** they see: class rank and percentage

**Given** report card is available
**When** downloading
**Then** they can download: PDF report card
**And** they can view: historical report cards
**And** comparison with class average shown

**Given** progress tracking
**When** viewing trends
**Then** they see: performance graph over exams
**And** they see: subject-wise improvement/decline

---

### Story 13.5: Fee Details & Online Payment

As a **parent**,
I want **to view fees and pay online**,
So that **I can manage payments conveniently**.

**Acceptance Criteria:**

**Given** parent is on fees section
**When** viewing
**Then** they see: current dues with breakdown
**And** they see: payment history with receipts
**And** they see: upcoming fee schedule

**Given** making payment
**When** initiating
**Then** they select: invoices to pay
**And** they see: total including any late fee
**And** they are redirected to payment gateway
**And** on success, receipt is shown/emailed

**Given** downloading receipts
**When** accessing payment history
**Then** each payment has download option
**And** receipts show complete details
**And** filters by date range available

---

### Story 13.6: Homework & Grades View

As a **parent**,
I want **to view my child's homework and grades**,
So that **I can monitor their work**.

**Acceptance Criteria:**

**Given** parent is on homework section
**When** viewing assignments
**Then** they see: pending assignments with due dates
**And** they see: submitted assignments with status
**And** they see: graded assignments with marks/feedback

**Given** assignment details
**When** viewing
**Then** they see: assignment instructions
**And** they see: attached resources
**And** they see: submission (if uploaded)
**And** they see: grade and teacher feedback

---

### Story 13.7: Notices & Circulars

As a **parent**,
I want **to view school notices**,
So that **I stay informed of school updates**.

**Acceptance Criteria:**

**Given** parent is on notices section
**When** viewing
**Then** they see: list of notices (newest first)
**And** priority notices highlighted
**And** unread notices marked
**And** filters by date, priority available

**Given** notice requires acknowledgement
**When** viewing
**Then** acknowledgement button is shown
**And** clicking confirms reading
**And** acknowledgement timestamp recorded

---

### Story 13.8: Parent-Teacher Messaging

As a **parent**,
I want **to message my child's teachers**,
So that **I can communicate about my child**.

**Acceptance Criteria:**

**Given** parent is on messaging
**When** starting conversation
**Then** they see: list of teachers (class teacher, subject teachers)
**And** they can select: teacher to message
**And** they can type and send message
**And** message history is maintained

**Given** response is received
**When** notified
**Then** push notification is sent
**And** message appears in inbox
**And** unread count is shown

---

### Story 13.9: Multi-Child Support

As a **parent with multiple children**,
I want **to manage all children from one account**,
So that **I don't need separate logins**.

**Acceptance Criteria:**

**Given** parent has multiple children enrolled
**When** logging in
**Then** they see: child selector/switcher
**And** switching shows: that child's data
**And** profile photo/name clearly indicates current child

**Given** dashboard view
**When** multiple children
**Then** summary shows: all children briefly
**And** alerts show: pending items for any child
**And** quick switch between children

**Given** fees view
**When** multiple children
**Then** they can view: combined dues
**And** they can pay: for multiple children together
**And** receipts show: child name
