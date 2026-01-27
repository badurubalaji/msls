# Epic 21: Visitor & Gate Management

**Phase:** 2 (Extended)
**Priority:** Medium - Security feature

## Epic Goal

Enable secure visitor management and student gate pass processing.

## User Value

Security can register visitors, parents can request gate passes, and all entries/exits are logged.

## FRs Covered

FR-VG-01 to FR-VG-05

---

## Stories

### Story 21.1: Visitor Registration

As a **security staff**,
I want **to register visitors**,
So that **visitor entry is tracked**.

**Acceptance Criteria:**

**Given** visitor arrives
**When** registering
**Then** they can enter: name, phone, company
**And** they can capture: photo (webcam)
**And** they can enter: ID type, number
**And** they can select: purpose, whom to meet
**And** expected exit time set

**Given** registration complete
**When** saved
**Then** visitor badge number assigned
**And** badge can be printed
**And** host is notified
**And** entry time recorded

---

### Story 21.2: Visitor Badge Printing

As a **security staff**,
I want **to print visitor badges**,
So that **visitors are identified**.

**Acceptance Criteria:**

**Given** visitor is registered
**When** printing badge
**Then** badge shows: photo, name, badge number
**And** badge shows: whom to meet, purpose
**And** badge shows: date, entry time
**And** QR code included for checkout

---

### Story 21.3: Visitor Checkout

As a **security staff**,
I want **to checkout visitors**,
So that **exit is recorded**.

**Acceptance Criteria:**

**Given** visitor is leaving
**When** checking out
**Then** they can scan: badge QR or enter number
**And** exit time recorded
**And** badge collected
**And** visitor record complete

**Given** visitor overstays
**When** beyond expected time
**Then** alert shown
**And** can extend visit if needed

---

### Story 21.4: Student Gate Pass Request

As a **parent**,
I want **to request gate pass for my child**,
So that **early pickup is authorized**.

**Acceptance Criteria:**

**Given** parent needs gate pass
**When** requesting
**Then** they can select: pass type (early leave, mid-day)
**And** they can enter: reason
**And** they can enter: exit time, return time
**And** they can enter: pickup person details

**Given** request submitted
**When** saved
**Then** request goes for approval
**And** class teacher notified
**And** status: pending

---

### Story 21.5: Gate Pass Approval

As a **teacher/admin**,
I want **to approve gate passes**,
So that **student exits are authorized**.

**Acceptance Criteria:**

**Given** gate pass request exists
**When** reviewing
**Then** approver sees: student details, reason
**And** approver can: approve, reject
**And** rejection requires: reason
**And** on approval: parent notified

**Given** pass is approved
**When** pickup time
**Then** security can verify: pass, pickup person ID
**And** student released
**And** exit time recorded
**And** parent confirmation SMS sent

---

### Story 21.6: Vehicle Entry Logging

As a **security staff**,
I want **to log vehicle entries**,
So that **vehicle movement is tracked**.

**Acceptance Criteria:**

**Given** vehicle arrives
**When** logging
**Then** they can enter: vehicle number, type
**And** they can enter: driver name, phone
**And** they can select: purpose (visitor, delivery, staff)
**And** entry time recorded

**Given** vehicle exits
**When** logging
**Then** exit time recorded
**And** duration calculated
**And** log complete

---

### Story 21.7: Visitor Pre-Registration

As a **staff member**,
I want **to pre-register expected visitors**,
So that **check-in is faster**.

**Acceptance Criteria:**

**Given** staff expects visitor
**When** pre-registering
**Then** they can enter: visitor name, phone
**And** they can enter: expected date, time
**And** they can enter: purpose
**And** approval code generated

**Given** pre-registered visitor arrives
**When** checking in
**Then** security can search: by phone or code
**And** details pre-filled
**And** only verification needed
**And** faster check-in

---

### Story 21.8: Security Dashboard

As a **security supervisor**,
I want **to view security dashboard**,
So that **current status is known**.

**Acceptance Criteria:**

**Given** supervisor views dashboard
**When** checking
**Then** they see: visitors currently inside
**And** they see: active gate passes
**And** they see: vehicles inside
**And** they see: today's entry count

**Given** reports needed
**When** generating
**Then** they can see: daily visitor report
**And** they can see: gate pass usage
**And** they can filter: by date, type
**And** export available
