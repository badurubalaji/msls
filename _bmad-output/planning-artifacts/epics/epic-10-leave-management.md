# Epic 10: Leave Management

**Phase:** 1 (MVP)
**Priority:** Medium - Staff HR operations

## Epic Goal

Enable staff to request leave and managers to approve with policy enforcement.

## User Value

Staff can apply for leave online, approvers can review and approve, HR can track balances and generate reports.

## FRs Covered

FR-LV-01 to FR-LV-07

---

## Stories

### Story 10.1: Leave Type Configuration

As an **HR administrator**,
I want **to configure leave types with policies**,
So that **leave rules are standardized**.

**Acceptance Criteria:**

**Given** HR is on leave settings
**When** creating a leave type
**Then** they can enter: name (Casual Leave, Sick Leave, etc.)
**And** they can set: annual quota (e.g., 12 days)
**And** they can set: carry forward rules (max days, expiry)
**And** they can set: encashment allowed (yes/no)

**Given** leave type is created
**When** setting policies
**Then** they can set: advance notice required (days)
**And** they can set: maximum consecutive days
**And** they can set: minimum days per application
**And** they can set: requires document (e.g., medical certificate for sick leave > 2 days)

---

### Story 10.2: Leave Balance Initialization

As an **HR administrator**,
I want **to initialize and manage leave balances**,
So that **staff have correct entitlements**.

**Acceptance Criteria:**

**Given** a new academic year begins
**When** initializing balances
**Then** each staff member gets: configured quota per leave type
**And** carry forward from previous year is added (within limits)
**And** pro-rata calculation for mid-year joiners

**Given** balance adjustments are needed
**When** HR makes adjustment
**Then** they can add/deduct days with reason
**And** adjustment is logged in history
**And** new balance is reflected immediately

**Given** viewing balances
**When** staff checks their leaves
**Then** they see: total, used, balance for each type
**And** they see: pending applications deducted from available
**And** they see: encashable balance if applicable

---

### Story 10.3: Leave Application Submission

As a **staff member**,
I want **to apply for leave online**,
So that **I don't need paper forms**.

**Acceptance Criteria:**

**Given** a staff member is applying for leave
**When** filling the application
**Then** they can select: leave type
**And** they can select: from date, to date
**And** they can select: half-day option (first half/second half)
**And** they can enter: reason for leave
**And** system shows: days being applied, balance after

**Given** leave rules exist
**When** submitting application
**Then** system validates: advance notice requirement
**And** system validates: maximum consecutive days
**And** system validates: sufficient balance
**And** warnings shown for policy violations

**Given** application is valid
**When** submitting
**Then** application status is "pending"
**And** notification sent to approver
**And** days are marked as "pending" in balance

---

### Story 10.4: Leave Approval Workflow

As a **manager/approver**,
I want **to review and approve leave requests**,
So that **leaves are properly authorized**.

**Acceptance Criteria:**

**Given** a leave application is pending
**When** approver views pending list
**Then** they see: applicant name, leave type, dates, reason
**And** they see: team calendar (who else is on leave)
**And** they see: applicant's leave balance

**Given** approver is reviewing
**When** making decision
**Then** they can: approve, reject, or request modification
**And** they must enter: comment (mandatory for rejection)
**And** decision is recorded with timestamp

**Given** decision is made
**When** saved
**Then** applicant is notified of decision
**And** if approved: balance is deducted
**And** if rejected: pending days released back
**And** leave appears in calendar if approved

---

### Story 10.5: Sandwich Rule Implementation

As an **HR administrator**,
I want **sandwich rule to be automatically applied**,
So that **intervening holidays are counted correctly**.

**Acceptance Criteria:**

**Given** sandwich rule is configured
**When** leave spans across weekend/holiday
**Then** system calculates: if leave before and after holiday, include holiday as leave
**And** calculation shown in application preview
**And** staff is informed of total days being deducted

**Given** different rules for different leave types
**When** applying
**Then** casual leave may include sandwich days
**And** sick leave may exclude sandwich days
**And** rule configuration is per leave type

---

### Story 10.6: Leave Calendar View

As a **manager or HR**,
I want **to view team leave calendar**,
So that **I can see availability at a glance**.

**Acceptance Criteria:**

**Given** a manager is viewing team calendar
**When** looking at month view
**Then** they see: all team members with leave days marked
**And** different colors for different leave types
**And** pending leaves shown differently from approved

**Given** filtering is needed
**When** applying filters
**Then** they can filter by: department, leave type, date range
**And** they can view: individual staff calendar
**And** they can export calendar

---

### Story 10.7: Leave Reports

As an **HR administrator**,
I want **to generate leave reports**,
So that **leave utilization can be analyzed**.

**Acceptance Criteria:**

**Given** HR is on leave reports
**When** generating leave summary
**Then** they see: staff-wise leave balance and usage
**And** they can filter by: department, designation, leave type
**And** they can export to Excel

**Given** trend analysis is needed
**When** viewing patterns
**Then** they see: monthly leave trends
**And** they see: high leave-taking staff
**And** they see: leave type distribution

**Given** compliance report is needed
**When** generating
**Then** they see: pending applications count
**And** they see: expired/lapsed leaves
**And** they see: carry forward summary

---

### Story 10.8: Leave Notifications

As a **staff member**,
I want **to receive notifications about my leave status**,
So that **I stay informed**.

**Acceptance Criteria:**

**Given** leave application is submitted
**When** status changes
**Then** notification sent on: submission confirmation
**And** notification sent on: approval
**And** notification sent on: rejection
**And** notification includes: leave details, new status

**Given** leave balance is low
**When** reaching threshold
**Then** notification sent about remaining balance
**And** reminder before year-end for unused leave

**Given** upcoming leave
**When** 1 day before leave starts
**Then** reminder sent to staff and approver
**And** out-of-office reminder shown
