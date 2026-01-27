# Epic 3: School Setup & Admissions

**Phase:** 1 (MVP)
**Priority:** High - Required before student enrollment

## Epic Goal

Enable school administrators to configure their school structure and manage the complete admission process from enquiry to enrollment.

## User Value

Admins can set up branches, academic years, and process student admissions end-to-end.

## FRs Covered

FR-CF-02, FR-CF-03, FR-ADM-01 to FR-ADM-07

---

## Stories

### Story 3.1: Branch Configuration

As a **super admin**,
I want **to configure branches within a tenant**,
So that **multi-branch schools can manage each location separately**.

**Acceptance Criteria:**

**Given** a super admin is on tenant settings
**When** they add a new branch
**Then** they can enter: name, code, address, contact details
**And** they can set one branch as primary
**And** they can configure branch-specific settings (logo, timezone)
**And** branch is created and available for assignment

**Given** branches exist
**When** viewing the branch list
**Then** all branches are displayed with status
**And** each branch shows student/staff count
**And** branches can be activated/deactivated

---

### Story 3.2: Academic Year Management

As an **administrator**,
I want **to configure academic years with terms and holidays**,
So that **all operations align with the school calendar**.

**Acceptance Criteria:**

**Given** an admin is on academic year settings
**When** they create a new academic year
**Then** they can enter: name (e.g., "2025-26"), start date, end date
**And** they can define terms/semesters with dates
**And** they can mark one year as "current"
**And** they can add holidays with name and date

**Given** an academic year is set as current
**When** any module operates
**Then** it defaults to the current academic year context
**And** users can switch to view historical years (read-only)

**Given** holidays are configured
**When** viewing the school calendar
**Then** holidays are highlighted
**And** attendance marking is blocked on holidays

---

### Story 3.3: Admission Session Configuration

As an **administrator**,
I want **to configure admission sessions for different classes**,
So that **admissions can be processed in organized cycles**.

**Acceptance Criteria:**

**Given** an admin is on admission settings
**When** they create an admission session
**Then** they can enter: name, academic year, start/end dates
**And** they can select applicable classes for this session
**And** they can set maximum seats per class
**And** they can configure required documents list
**And** they can set admission fee amount

**Given** an admission session exists
**When** viewing the session
**Then** they see: applications count, seats filled, available seats
**And** they can open/close the session
**And** they can extend deadline if needed

---

### Story 3.4: Admission Enquiry Management

As a **front office staff**,
I want **to capture and track admission enquiries**,
So that **interested families can be followed up systematically**.

**Acceptance Criteria:**

**Given** a staff member is handling an enquiry
**When** they create a new enquiry
**Then** they can enter: parent name, phone, email, student name, class interested
**And** they can add notes from the conversation
**And** enquiry is assigned a unique enquiry number
**And** enquiry status is "new"

**Given** enquiries exist in the system
**When** viewing the enquiry list
**Then** they see all enquiries with status (new, contacted, interested, converted, closed)
**And** they can filter by status, class, date range
**And** they can add follow-up notes with date

**Given** an enquiry is converted
**When** creating an application from enquiry
**Then** enquiry data pre-fills the application form
**And** enquiry status changes to "converted"
**And** enquiry links to the application

---

### Story 3.5: Online Admission Application

As a **parent**,
I want **to submit an admission application online**,
So that **I don't need to visit the school for initial application**.

**Acceptance Criteria:**

**Given** a parent accesses the admission portal
**When** they select class and fill the application form
**Then** they can enter: student details (name, DOB, gender, Aadhaar)
**And** they can enter: parent/guardian details
**And** they can enter: previous school details (if applicable)
**And** they can upload required documents (birth certificate, photos, etc.)
**And** application is saved with unique application number

**Given** application is submitted
**When** all required fields and documents are provided
**Then** application status is "submitted"
**And** confirmation email/SMS is sent to parent
**And** application appears in admin dashboard for review

**Given** a parent wants to check application status
**When** they login with application number and phone
**Then** they see current status and any remarks
**And** they can upload additional documents if requested

---

### Story 3.6: Application Review & Entrance Test

As an **admission committee member**,
I want **to review applications and schedule entrance tests**,
So that **qualified candidates can be evaluated**.

**Acceptance Criteria:**

**Given** applications are submitted
**When** reviewing an application
**Then** admin sees all submitted information and documents
**And** admin can verify document authenticity
**And** admin can add review comments
**And** admin can update status: under_review, documents_pending, test_scheduled, rejected

**Given** entrance test is configured for a class
**When** scheduling tests
**Then** admin can set test date, time, and venue
**And** admin can assign students to test slots
**And** hall ticket is generated for each student
**And** SMS/email notification is sent with test details

**Given** entrance test is conducted
**When** entering test results
**Then** admin can enter marks for each subject
**And** total score is calculated automatically
**And** application moves to "test_completed" status

---

### Story 3.7: Merit List & Admission Decision

As an **admission committee**,
I want **to generate merit lists and make admission decisions**,
So that **students can be selected fairly based on criteria**.

**Acceptance Criteria:**

**Given** test results are entered
**When** generating merit list
**Then** students are ranked by total score
**And** merit list shows: rank, name, score, status
**And** cutoff score can be applied to filter

**Given** a merit list is generated
**When** making admission decisions
**Then** admin can select students for admission
**And** admin can mark: selected, waitlisted, rejected
**And** selected students' status changes to "offer_sent"
**And** offer letter is generated with fee details

**Given** an offer is sent
**When** parent accepts and pays admission fee
**Then** payment is recorded
**And** application status changes to "enrolled"
**And** student record is created from application data
**And** admission process is complete

---

### Story 3.8: Admission Reports & Analytics

As an **administrator**,
I want **to view admission reports and analytics**,
So that **I can track admission progress and plan capacity**.

**Acceptance Criteria:**

**Given** an admin is on admission reports
**When** viewing the dashboard
**Then** they see: total applications, by status, by class
**And** they see: conversion rate (enquiry to application to enrolled)
**And** they see: seats filled vs available by class

**Given** filtering options are available
**When** admin applies filters
**Then** they can filter by: session, class, date range, status
**And** data updates accordingly
**And** reports can be exported to Excel/PDF
