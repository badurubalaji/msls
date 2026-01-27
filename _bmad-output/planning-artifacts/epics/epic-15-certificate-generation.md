# Epic 15: Certificate Generation

**Phase:** 1 (MVP)
**Priority:** Medium - Administrative operations

## Epic Goal

Enable certificate requests, approval workflow, and generation with verification.

## User Value

Parents can request certificates, staff can process approvals, and generated certificates include QR verification.

## FRs Covered

FR-CT-01 to FR-CT-06

---

## Stories

### Story 15.1: Certificate Request Submission

As a **parent**,
I want **to request certificates online**,
So that **I don't need to visit school for requests**.

**Acceptance Criteria:**

**Given** parent is on certificate section
**When** submitting request
**Then** they can select: certificate type (TC, Bonafide, etc.)
**And** they can enter: purpose of certificate
**And** they can select: number of copies needed
**And** they can select: delivery mode (collect, courier)
**And** they can see: fee for the certificate

**Given** request is submitted
**When** saved
**Then** request number is generated
**And** status is "pending"
**And** payment prompt if fee required
**And** confirmation sent to parent

---

### Story 15.2: Certificate Fee Payment

As a **parent**,
I want **to pay certificate fees online**,
So that **processing can begin**.

**Acceptance Criteria:**

**Given** certificate request is submitted
**When** fee is applicable
**Then** payment button is shown
**And** fee breakdown displayed
**And** online payment via gateway available

**Given** payment is completed
**When** confirmed
**Then** request status changes to "fee_paid"
**And** receipt is generated
**And** processing queue is updated

---

### Story 15.3: Certificate Approval Workflow

As an **administrator**,
I want **to review and approve certificate requests**,
So that **requests are processed properly**.

**Acceptance Criteria:**

**Given** admin is on certificate requests
**When** viewing pending requests
**Then** they see: list with student, type, date
**And** they can filter: by type, status, date
**And** they can view: request details

**Given** reviewing a request
**When** processing
**Then** they can: approve, reject, request more info
**And** rejection requires: reason
**And** status and notes are updated
**And** notification sent to parent

---

### Story 15.4: Department Clearance (for TC)

As an **administrator**,
I want **clearances from departments for TC**,
So that **all dues are cleared before TC issue**.

**Acceptance Criteria:**

**Given** TC request is approved
**When** clearance is needed
**Then** clearance requests sent to: fees, library, lab
**And** each department can: clear or flag dues
**And** all clearances needed before generation

**Given** a department flags dues
**When** viewing request
**Then** dues details are shown
**And** parent is notified to clear dues
**And** request waits until cleared

**Given** all departments clear
**When** status updated
**Then** request moves to: ready for generation
**And** admin is notified

---

### Story 15.5: Certificate Template Management

As an **administrator**,
I want **to configure certificate templates**,
So that **certificates have proper formatting**.

**Acceptance Criteria:**

**Given** admin is on template settings
**When** creating template
**Then** they can: use rich text editor
**And** they can: add placeholders {{student_name}}, {{class}}, etc.
**And** they can: upload header/footer images
**And** they can: set page size and orientation

**Given** template exists for type
**When** generating certificate
**Then** correct template is used
**And** placeholders are replaced with data
**And** school branding is applied

---

### Story 15.6: Certificate Generation

As an **administrator**,
I want **to generate certificates**,
So that **approved requests are fulfilled**.

**Acceptance Criteria:**

**Given** request is ready for generation
**When** generating
**Then** certificate data is populated from student record
**And** certificate number is auto-generated
**And** QR code is embedded for verification
**And** PDF is generated and stored

**Given** TC certificate
**When** generating
**Then** all TC-specific fields are populated
**And** leaving date, reason, conduct filled
**And** fees paid till date shown
**And** attendance record included

---

### Story 15.7: Certificate QR Verification

As an **external party**,
I want **to verify certificate authenticity**,
So that **I can trust the document**.

**Acceptance Criteria:**

**Given** a certificate has QR code
**When** scanning
**Then** verification page opens
**And** shows: certificate details (number, student, type, date)
**And** shows: verification status (valid/invalid)
**And** shows: school details

**Given** manual verification
**When** entering certificate number
**Then** same verification result shown
**And** public verification page available
**And** no login required

---

### Story 15.8: Certificate Delivery & Tracking

As a **parent**,
I want **to track my certificate request**,
So that **I know when it will be ready**.

**Acceptance Criteria:**

**Given** parent is tracking request
**When** viewing status
**Then** they see: current status in workflow
**And** they see: expected ready date
**And** they see: any pending actions needed

**Given** certificate is ready
**When** notification sent
**Then** parent is notified via SMS/email
**And** collection/courier details provided
**And** download available for PDF copy

**Given** collection at school
**When** parent collects
**Then** staff marks as delivered
**And** acknowledgement is recorded
**And** delivery date logged

---

### Story 15.9: Duplicate Certificate Handling

As an **administrator**,
I want **to issue duplicate certificates**,
So that **lost certificates can be replaced**.

**Acceptance Criteria:**

**Given** original certificate was issued
**When** duplicate is requested
**Then** additional fee is charged (configurable 2x)
**And** duplicate is clearly marked as "DUPLICATE"
**And** new certificate number is generated
**And** links to original certificate record

**Given** duplicate is generated
**When** viewing history
**Then** both original and duplicate shown
**And** duplicate indicates: issue date, reason
**And** original is not invalidated
