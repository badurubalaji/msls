# Epic 4: Student Lifecycle Management

**Phase:** 1 (MVP)
**Priority:** High - Core student data management

## Epic Goal

Enable comprehensive management of student records throughout their entire school journey.

## User Value

Admins and teachers can manage student profiles, track health records, handle behavioral incidents, and process promotions.

## FRs Covered

FR-SM-01 to FR-SM-08

---

## Stories

### Story 4.1: Student Profile Management

As an **administrator**,
I want **to create and manage student profiles with complete information**,
So that **all student data is centrally managed**.

**Acceptance Criteria:**

**Given** an admin is creating a new student
**When** they fill the student form
**Then** they can enter: personal details (name, DOB, gender, blood group, Aadhaar)
**And** they can enter: academic details (admission number, class, section, roll number)
**And** they can enter: address details (current, permanent)
**And** they can upload: photo, birth certificate
**And** admission number is auto-generated with tenant prefix
**And** student status is set to "active"

**Given** a student profile exists
**When** viewing the profile
**Then** all information is displayed in organized sections
**And** edit history is available in audit log
**And** quick actions are available (edit, view attendance, view fees)

---

### Story 4.2: Guardian & Emergency Contact Management

As an **administrator**,
I want **to manage student guardians and emergency contacts**,
So that **parents can be contacted and emergency situations handled**.

**Acceptance Criteria:**

**Given** a student profile is open
**When** adding a guardian
**Then** they can enter: name, relation, phone, email, occupation, address
**And** they can mark one guardian as primary
**And** they can set login access for parent portal
**And** they can add multiple guardians (father, mother, other)

**Given** emergency contacts section
**When** adding emergency contacts
**Then** they can add contacts beyond guardians
**And** priority order can be set (1st, 2nd, 3rd)
**And** each contact has: name, relation, phone

**Given** a guardian has portal access
**When** they login
**Then** they see only their linked children's data
**And** they can update their own contact details

---

### Story 4.3: Student Enrollment History

As an **administrator**,
I want **to track student enrollment across academic years**,
So that **complete academic history is maintained**.

**Acceptance Criteria:**

**Given** a student is enrolled
**When** viewing enrollment history
**Then** they see: academic year, class, section, roll number, status
**And** history shows all years from admission to current

**Given** a new academic year begins
**When** promotions are processed
**Then** new enrollment record is created
**And** previous year enrollment is marked "completed"
**And** class teacher assignment is recorded

**Given** a student transfers out
**When** processing transfer
**Then** current enrollment status changes to "transferred"
**And** transfer date and reason are recorded
**And** student status changes to "inactive"

---

### Story 4.4: Student Health Records

As a **school nurse/administrator**,
I want **to maintain student health records**,
So that **medical emergencies can be handled with proper information**.

**Acceptance Criteria:**

**Given** a student's health section
**When** recording health information
**Then** they can enter: blood group, height, weight, vision status
**And** they can record: allergies with severity
**And** they can record: chronic conditions (asthma, diabetes, etc.)
**And** they can record: regular medications
**And** health data is marked confidential (restricted access)

**Given** vaccination tracking
**When** recording vaccinations
**Then** they can enter: vaccine name, date administered, next due date
**And** reminders can be set for upcoming vaccinations
**And** vaccination certificate can be uploaded

**Given** a medical incident occurs
**When** recording the incident
**Then** they can enter: date, time, description, action taken
**And** they can record: parent notified (yes/no), hospital visit required
**And** incident history is maintained

---

### Story 4.5: Student Behavioral Tracking

As a **teacher or administrator**,
I want **to record and track student behavioral incidents**,
So that **patterns can be identified and addressed**.

**Acceptance Criteria:**

**Given** a behavioral incident occurs
**When** recording the incident
**Then** they can enter: date, time, location, incident type
**And** incident types include: positive recognition, minor infraction, major violation
**And** they can describe: what happened, witnesses, student response
**And** they can record: action taken, parent meeting required

**Given** a follow-up action is required
**When** scheduling follow-up
**Then** they can set: meeting date, participants, outcomes expected
**And** notification is sent to relevant parties
**And** meeting outcome can be recorded later

**Given** viewing a student's behavioral history
**When** accessing the records
**Then** they see chronological list of incidents
**And** they can filter by: type, date range, severity
**And** pattern analysis shows trends (improving/declining)

---

### Story 4.6: Student Document Management

As an **administrator**,
I want **to manage and verify student documents**,
So that **required documents are maintained and verified**.

**Acceptance Criteria:**

**Given** a student's documents section
**When** uploading a document
**Then** they can select document type: birth certificate, Aadhaar, transfer certificate, etc.
**And** they can upload file (PDF, image)
**And** they can enter: document number, issue date, expiry date (if applicable)
**And** document status is "pending_verification"

**Given** a document needs verification
**When** admin verifies the document
**Then** they can mark: verified, rejected (with reason)
**And** verification date and verifier are recorded
**And** document status updates accordingly

**Given** required documents are configured
**When** viewing a student profile
**Then** they see document checklist with status
**And** missing documents are highlighted
**And** bulk reminder can be sent for missing documents

---

### Story 4.7: Student Promotion & Retention

As an **administrator**,
I want **to process student promotions at year end**,
So that **students are moved to appropriate classes**.

**Acceptance Criteria:**

**Given** year-end promotion time
**When** initiating promotion
**Then** admin can select: academic year, class, section
**And** system shows students with their result status
**And** students can be marked: promote, retain, transfer

**Given** promotion rules are configured
**When** applying auto-promotion
**Then** students meeting criteria are auto-marked for promotion
**And** students not meeting criteria are flagged for review
**And** manual override is available

**Given** promotions are confirmed
**When** processing promotions
**Then** new enrollments are created for next year
**And** section assignments are made (manual or auto-distributed)
**And** roll numbers are generated
**And** promotion report is generated for records

---

### Story 4.8: Student Search & Bulk Operations

As an **administrator**,
I want **to search students and perform bulk operations**,
So that **managing large student populations is efficient**.

**Acceptance Criteria:**

**Given** the student list page
**When** searching for students
**Then** they can search by: name, admission number, phone, class
**And** results show in paginated list
**And** advanced filters include: class, section, status, gender, transport

**Given** multiple students are selected
**When** performing bulk operations
**Then** available operations include: bulk SMS, bulk email, bulk status update
**And** confirmation is required before execution
**And** operation log is maintained

**Given** export is needed
**When** exporting student list
**Then** data can be exported to Excel/CSV
**And** column selection is available
**And** filters are applied to export
