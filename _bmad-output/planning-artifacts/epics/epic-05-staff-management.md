# Epic 5: Staff Management

**Phase:** 1 (MVP)
**Priority:** High - Required for school operations

## Epic Goal

Enable complete management of teaching and non-teaching staff including profiles, attendance, and basic payroll.

## User Value

HR and Admins can manage staff records, track qualifications, monitor attendance, and process salaries.

## FRs Covered

FR-STF-01 to FR-STF-08

---

## Stories

### Story 5.1: Staff Profile Management

As an **HR administrator**,
I want **to create and manage staff profiles**,
So that **all employee information is centrally maintained**.

**Acceptance Criteria:**

**Given** HR is creating a new staff member
**When** they fill the staff form
**Then** they can enter: personal details (name, DOB, gender, blood group)
**And** they can enter: contact details (phone, email, address)
**And** they can enter: employment details (employee ID, join date, designation)
**And** they can select: staff type (teaching/non-teaching)
**And** they can select: department, reporting manager
**And** employee ID is auto-generated with prefix

**Given** a staff profile exists
**When** viewing the profile
**Then** all information is displayed in organized tabs
**And** photo, documents, qualifications are accessible
**And** quick actions available (edit, view attendance, view salary)

---

### Story 5.2: Department & Designation Hierarchy

As an **HR administrator**,
I want **to configure departments and designations**,
So that **organizational structure is properly defined**.

**Acceptance Criteria:**

**Given** HR is on organization settings
**When** creating a department
**Then** they can enter: name, code, description, head (staff member)
**And** departments can be active/inactive

**Given** designations are configured
**When** viewing designations
**Then** they see: designation name, level, department mapping
**And** salary grade can be linked to designation
**And** reporting hierarchy is defined

**Given** staff members exist
**When** viewing org chart
**Then** hierarchical view shows departments and reporting lines
**And** staff can be filtered by department
**And** vacant positions are highlighted

---

### Story 5.3: Staff Qualifications & Experience

As an **HR administrator**,
I want **to record staff qualifications and work experience**,
So that **professional credentials are documented**.

**Acceptance Criteria:**

**Given** a staff profile is open
**When** adding qualifications
**Then** they can enter: degree, institution, year, grade/percentage
**And** they can upload: certificate/mark sheet
**And** multiple qualifications can be added

**Given** work experience section
**When** adding experience
**Then** they can enter: organization, designation, from date, to date
**And** they can enter: responsibilities, reason for leaving
**And** experience letter can be uploaded

**Given** professional certifications exist
**When** tracking certifications
**Then** they can enter: certification name, issuing body, valid from, valid until
**And** expiry reminders can be set
**And** renewal tracking is available

---

### Story 5.4: Staff Attendance Marking

As a **staff member**,
I want **to mark my attendance**,
So that **my presence is recorded daily**.

**Acceptance Criteria:**

**Given** a staff member is logged in
**When** marking attendance
**Then** they can mark: present, half-day (with reason)
**And** check-in time is recorded
**And** check-out time can be recorded later
**And** late arrival is flagged if after threshold

**Given** biometric integration is configured
**When** staff punches in/out
**Then** attendance is auto-marked from biometric data
**And** manual override requires HR approval

**Given** attendance regularization is needed
**When** staff submits regularization request
**Then** they can select: date, reason, supporting document
**And** request goes to reporting manager for approval
**And** approved regularization updates attendance record

---

### Story 5.5: Staff Attendance Reports

As an **HR administrator**,
I want **to view and analyze staff attendance**,
So that **attendance patterns can be monitored**.

**Acceptance Criteria:**

**Given** HR is on attendance reports
**When** viewing daily attendance
**Then** they see: all staff with status (present, absent, leave, half-day)
**And** they can filter by: department, designation, date

**Given** monthly attendance summary
**When** generating report
**Then** they see: total working days, days present, leaves taken, late days
**And** summary is per staff member
**And** export to Excel is available

**Given** attendance analysis
**When** viewing patterns
**Then** they see: frequent late-comers, high absenteeism
**And** trends over months are displayed
**And** alerts for policy violations are shown

---

### Story 5.6: Teacher Subject Assignment

As an **academic administrator**,
I want **to assign subjects to teachers**,
So that **teaching responsibilities are defined**.

**Acceptance Criteria:**

**Given** a teacher profile
**When** assigning subjects
**Then** they can select: subject, class, section
**And** they can set: is class teacher (yes/no)
**And** multiple assignments can be made

**Given** subject assignments exist
**When** viewing workload summary
**Then** they see: total periods per week
**And** they see: classes assigned, subjects taught
**And** workload comparison across teachers is available

**Given** assignment conflicts exist
**When** detecting conflicts
**Then** system warns if teacher is over-assigned
**And** system warns if subject has no teacher
**And** workload balance recommendations are shown

---

### Story 5.7: Salary Structure Configuration

As an **HR administrator**,
I want **to configure salary structures**,
So that **payroll can be calculated correctly**.

**Acceptance Criteria:**

**Given** HR is on salary settings
**When** creating a salary structure
**Then** they can define: basic pay, allowances (HRA, DA, TA, etc.)
**And** they can define: deductions (PF, ESI, TDS, etc.)
**And** they can link structure to designation/grade

**Given** a staff member's salary configuration
**When** viewing salary details
**Then** they see: base structure, individual modifications
**And** gross salary and net salary are calculated
**And** tax calculation is shown

**Given** salary revisions occur
**When** updating salary
**Then** effective date is recorded
**And** history of revisions is maintained
**And** arrears calculation is supported

---

### Story 5.8: Staff Document Management

As an **HR administrator**,
I want **to manage staff documents with expiry tracking**,
So that **compliance requirements are met**.

**Acceptance Criteria:**

**Given** a staff member's documents section
**When** uploading a document
**Then** they can select type: Aadhaar, PAN, offer letter, contract, ID proof
**And** they can enter: document number, issue date, expiry date
**And** they can upload file (PDF, image)

**Given** documents have expiry dates
**When** expiry approaches
**Then** notification is sent 30 days before expiry
**And** dashboard shows expiring documents
**And** expired documents are flagged

**Given** document verification
**When** admin verifies
**Then** they can mark: verified, rejected
**And** verification details are recorded
**And** compliance report shows document status
