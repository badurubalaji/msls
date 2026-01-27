---
stepsCompleted: [1, 2, 3, 4]
status: complete
validatedAt: 2026-01-23
inputDocuments:
  - planning-artifacts/school-erp-prd/index.md
  - planning-artifacts/school-erp-prd/01-technical-architecture.md
  - planning-artifacts/school-erp-prd/02-core-foundation.md
  - planning-artifacts/school-erp-prd/03-student-management.md
  - planning-artifacts/school-erp-prd/04-academic-operations.md
  - planning-artifacts/school-erp-prd/05-admissions.md
  - planning-artifacts/school-erp-prd/06-examinations-grading.md
  - planning-artifacts/school-erp-prd/07-homework-assignments.md
  - planning-artifacts/school-erp-prd/08-online-quiz-assessment.md
  - planning-artifacts/school-erp-prd/09-digital-classroom.md
  - planning-artifacts/school-erp-prd/10-staff-management.md
  - planning-artifacts/school-erp-prd/11-leave-management.md
  - planning-artifacts/school-erp-prd/12-fees-payments.md
  - planning-artifacts/school-erp-prd/13-communication-system.md
  - planning-artifacts/school-erp-prd/14-parent-portal.md
  - planning-artifacts/school-erp-prd/15-student-portal.md
  - planning-artifacts/school-erp-prd/16-certificate-generation.md
  - planning-artifacts/school-erp-prd/17-transport-management.md
  - planning-artifacts/school-erp-prd/18-library-management.md
  - planning-artifacts/school-erp-prd/19-inventory-assets.md
  - planning-artifacts/school-erp-prd/20-visitor-gate-management.md
  - planning-artifacts/school-erp-prd/21-analytics-dashboards.md
  - planning-artifacts/school-erp-prd/22-ai-capabilities.md
  - planning-artifacts/architecture.md
  - planning-artifacts/ux-design-specification.md
---

# MSLS (Multi-School Learning System) - Epic Breakdown

## Overview

This document provides the complete epic and story breakdown for MSLS School ERP, decomposing the requirements from the PRD (22 modules across 3 phases), UX Design Specification, and Architecture decisions into implementable stories.

**Scope:** Phase 1 (MVP - 15 modules), Phase 2 (Extended - 7 modules), Phase 3 (Advanced - 1 module)

---

## Requirements Inventory

### Functional Requirements

#### Phase 1 - MVP (Core Foundation)

**FR-CF: Core Foundation Module**
- FR-CF-01: System shall support multi-tenant architecture with tenant isolation via tenant_id on all tables
- FR-CF-02: System shall support multi-branch management per tenant with branch-level configuration
- FR-CF-03: System shall manage academic years with session dates and holidays
- FR-CF-04: System shall provide user authentication with email/phone + password and OTP verification
- FR-CF-05: System shall support role-based access control (RBAC) with customizable permissions
- FR-CF-06: System shall implement row-level security (RLS) policies for tenant data isolation
- FR-CF-07: System shall provide configuration engine for feature flags and custom fields
- FR-CF-08: System shall support JWT RS256 token-based authentication (15min access, 7-day refresh)
- FR-CF-09: System shall use Argon2id for password hashing
- FR-CF-10: System shall support TOTP-based 2FA for administrative users

**FR-SM: Student Management Module**
- FR-SM-01: System shall manage complete student lifecycle from admission to alumni
- FR-SM-02: System shall store student profiles with personal, academic, and guardian information
- FR-SM-03: System shall track student enrollment history across academic years
- FR-SM-04: System shall manage student health records including medical conditions and vaccinations
- FR-SM-05: System shall track student behavioral incidents with actions taken
- FR-SM-06: System shall manage student documents with verification status
- FR-SM-07: System shall support student promotion/retention between classes
- FR-SM-08: System shall generate unique admission numbers per tenant

**FR-STF: Staff Management Module**
- FR-STF-01: System shall manage staff profiles including teaching and non-teaching staff
- FR-STF-02: System shall track staff qualifications and work experience
- FR-STF-03: System shall manage staff contact information and emergency contacts
- FR-STF-04: System shall support organizational hierarchy with departments and designations
- FR-STF-05: System shall track staff attendance with multiple marking sources
- FR-STF-06: System shall manage salary structures and payroll processing
- FR-STF-07: System shall track teacher subject assignments and workload
- FR-STF-08: System shall manage staff documents with expiry tracking

**FR-ADM: Admissions Module**
- FR-ADM-01: System shall manage admission sessions with configuration
- FR-ADM-02: System shall capture and track admission enquiries
- FR-ADM-03: System shall process online admission applications with document uploads
- FR-ADM-04: System shall schedule and manage entrance tests with marks entry
- FR-ADM-05: System shall support merit list generation and admission decisions
- FR-ADM-06: System shall convert approved applications to student enrollment
- FR-ADM-07: System shall track admission fees and payment status

**FR-AO: Academic Operations Module**
- FR-AO-01: System shall manage class and section structures
- FR-AO-02: System shall manage subjects with class-subject mappings
- FR-AO-03: System shall define period slots with break/lunch configurations
- FR-AO-04: System shall create and manage timetables for sections and teachers
- FR-AO-05: System shall support timetable substitution management
- FR-AO-06: System shall mark and track daily student attendance
- FR-AO-07: System shall support period-wise attendance tracking (optional)
- FR-AO-08: System shall generate attendance reports (daily, monthly, class-wise)
- FR-AO-09: System shall enforce attendance deadline and edit window rules

**FR-EX: Examinations & Grading Module**
- FR-EX-01: System shall configure exam types (unit test, term exam, final)
- FR-EX-02: System shall create examinations with subject-wise schedules
- FR-EX-03: System shall generate and print hall tickets
- FR-EX-04: System shall support marks entry with validation rules
- FR-EX-05: System shall configure grading scales (percentage and grade-based)
- FR-EX-06: System shall calculate results with grade mapping
- FR-EX-07: System shall generate report cards with configurable templates
- FR-EX-08: System shall support rank generation (class, section, subject-wise)
- FR-EX-09: System shall track exam eligibility based on attendance

**FR-HW: Homework & Assignments Module**
- FR-HW-01: System shall create assignments with multiple types (homework, project, worksheet)
- FR-HW-02: System shall support file attachments for assignments
- FR-HW-03: System shall allow students to submit assignments online
- FR-HW-04: System shall support late submission with penalty configuration
- FR-HW-05: System shall enable teacher grading with rubric support
- FR-HW-06: System shall track assignment completion and grading status
- FR-HW-07: System shall notify parents of pending/overdue assignments

**FR-LV: Leave Management Module**
- FR-LV-01: System shall configure leave types with policies
- FR-LV-02: System shall track leave balances per staff member
- FR-LV-03: System shall process leave applications with approval workflow
- FR-LV-04: System shall enforce leave rules (advance notice, max consecutive days)
- FR-LV-05: System shall support sandwich rule configuration
- FR-LV-06: System shall generate leave calendar and reports
- FR-LV-07: System shall notify stakeholders of leave status changes

**FR-FE: Fees & Payments Module**
- FR-FE-01: System shall define fee categories and structures
- FR-FE-02: System shall assign fee structures to students with discounts
- FR-FE-03: System shall generate invoices with automatic numbering
- FR-FE-04: System shall accept payments via multiple modes (cash, cheque, UPI, card, netbanking)
- FR-FE-05: System shall integrate with payment gateways (Razorpay, PayU)
- FR-FE-06: System shall generate and print/email receipts
- FR-FE-07: System shall apply late fees based on configurable rules
- FR-FE-08: System shall track dues and generate defaulter reports
- FR-FE-09: System shall process fee refunds with approval workflow

**FR-CM: Communication System Module**
- FR-CM-01: System shall create and publish notices with target audience selection
- FR-CM-02: System shall send SMS notifications via configurable providers
- FR-CM-03: System shall send email notifications with templates
- FR-CM-04: System shall send push notifications to mobile apps
- FR-CM-05: System shall track notice acknowledgements
- FR-CM-06: System shall support in-app messaging between parents and teachers
- FR-CM-07: System shall trigger automated notifications for key events

**FR-PP: Parent Portal Module**
- FR-PP-01: System shall provide parent registration and login
- FR-PP-02: System shall display child's attendance summary and calendar
- FR-PP-03: System shall show academic performance and report cards
- FR-PP-04: System shall display fee details and enable online payment
- FR-PP-05: System shall show homework status and grades
- FR-PP-06: System shall display school notices and circulars
- FR-PP-07: System shall enable parent-teacher messaging
- FR-PP-08: System shall support multiple children under one parent account

**FR-SP: Student Portal Module**
- FR-SP-01: System shall provide student login with age-appropriate access
- FR-SP-02: System shall display personal timetable and schedule
- FR-SP-03: System shall show attendance records
- FR-SP-04: System shall enable homework submission
- FR-SP-05: System shall display exam results and report cards
- FR-SP-06: System shall show assigned quizzes and tests
- FR-SP-07: System shall display library borrowed books

**FR-CT: Certificate Generation Module**
- FR-CT-01: System shall request certificates with approval workflow
- FR-CT-02: System shall verify department clearances for TC
- FR-CT-03: System shall generate certificates from configurable templates
- FR-CT-04: System shall include QR codes for verification
- FR-CT-05: System shall track certificate delivery and acknowledgement
- FR-CT-06: System shall support duplicate certificate issuance with additional fee

#### Phase 2 - Extended Features

**FR-QZ: Online Quiz & Assessment Module**
- FR-QZ-01: System shall maintain a question bank with multiple question types
- FR-QZ-02: System shall create quizzes with manual or auto-generated questions
- FR-QZ-03: System shall conduct timed quizzes with shuffle options
- FR-QZ-04: System shall auto-grade objective questions
- FR-QZ-05: System shall provide quiz analytics and question analysis
- FR-QZ-06: System shall support proctoring features (tab detection, fullscreen)
- FR-QZ-07: System shall enable practice mode with unlimited attempts

**FR-DC: Digital Classroom Module**
- FR-DC-01: System shall record classes (screen + audio)
- FR-DC-02: System shall process recordings with multiple quality options
- FR-DC-03: System shall support bookmarks and chapter markers
- FR-DC-04: System shall track student watch progress
- FR-DC-05: System shall manage digital content library
- FR-DC-06: System shall support explained PDFs with page-level audio

**FR-TR: Transport Management Module**
- FR-TR-01: System shall manage vehicles with document tracking
- FR-TR-02: System shall define routes with stops and timings
- FR-TR-03: System shall assign students to routes and stops
- FR-TR-04: System shall manage drivers and attendants
- FR-TR-05: System shall track real-time vehicle location via GPS
- FR-TR-06: System shall mark transport attendance
- FR-TR-07: System shall calculate distance-based fees

**FR-LB: Library Management Module**
- FR-LB-01: System shall catalog books with copy management
- FR-LB-02: System shall manage library members with borrowing limits
- FR-LB-03: System shall process book issue and return
- FR-LB-04: System shall handle reservations and waitlists
- FR-LB-05: System shall calculate and collect fines
- FR-LB-06: System shall generate circulation reports

**FR-IN: Inventory & Assets Module**
- FR-IN-01: System shall register and track assets with depreciation
- FR-IN-02: System shall manage asset locations and assignments
- FR-IN-03: System shall process maintenance requests
- FR-IN-04: System shall manage consumable stock levels
- FR-IN-05: System shall create purchase orders with approval
- FR-IN-06: System shall manage vendors

**FR-VG: Visitor & Gate Management Module**
- FR-VG-01: System shall register visitors with photo capture
- FR-VG-02: System shall issue and track visitor badges
- FR-VG-03: System shall process student gate passes with approval
- FR-VG-04: System shall log vehicle entries and exits
- FR-VG-05: System shall support pre-registration of expected visitors

**FR-AN: Analytics & Dashboards Module**
- FR-AN-01: System shall provide role-based dashboards
- FR-AN-02: System shall display academic analytics with trends
- FR-AN-03: System shall display financial analytics
- FR-AN-04: System shall display attendance analytics
- FR-AN-05: System shall support custom report builder
- FR-AN-06: System shall schedule and email automated reports

#### Phase 3 - Advanced Features

**FR-AI: AI Capabilities Module**
- FR-AI-01: System shall predict student performance risk using on-premise AI
- FR-AI-02: System shall identify dropout risk patterns
- FR-AI-03: System shall generate personalized learning recommendations
- FR-AI-04: System shall predict fee default probability
- FR-AI-05: System shall forecast admission demand
- FR-AI-06: System shall auto-summarize class recordings
- FR-AI-07: System shall auto-generate quiz questions from content
- FR-AI-08: System shall detect attendance and grading anomalies

---

### Non-Functional Requirements

**NFR-SEC: Security Requirements**
- NFR-SEC-01: All data must be encrypted at rest using AES-256
- NFR-SEC-02: All API communication must use HTTPS/TLS 1.3
- NFR-SEC-03: Passwords must be hashed using Argon2id
- NFR-SEC-04: JWT tokens must use RS256 algorithm with 15-minute expiry
- NFR-SEC-05: Row-level security must enforce tenant isolation at database level
- NFR-SEC-06: Sensitive data (Aadhaar, PAN) must be encrypted in application layer
- NFR-SEC-07: All user actions must be logged in audit trail
- NFR-SEC-08: API rate limiting must be enforced (100 req/min standard, 1000 req/min auth)
- NFR-SEC-09: OWASP Top 10 vulnerabilities must be prevented
- NFR-SEC-10: No sensitive data in logs (passwords, tokens, PII)

**NFR-PER: Performance Requirements**
- NFR-PER-01: API response time must be <200ms for 95th percentile
- NFR-PER-02: Dashboard must load within 2 seconds
- NFR-PER-03: System must support 500 concurrent users per tenant
- NFR-PER-04: Database queries must be optimized with proper indexing
- NFR-PER-05: File uploads must support up to 10MB per file
- NFR-PER-06: Batch operations must process 1000 records within 30 seconds

**NFR-SCA: Scalability Requirements**
- NFR-SCA-01: System must support horizontal scaling via Kubernetes
- NFR-SCA-02: Database must support read replicas for reporting
- NFR-SCA-03: File storage must use object storage (MinIO/S3)
- NFR-SCA-04: Cache layer must use Redis for session and frequently accessed data

**NFR-AVL: Availability Requirements**
- NFR-AVL-01: System must achieve 99.9% uptime SLA
- NFR-AVL-02: Automated failover must complete within 30 seconds
- NFR-AVL-03: Database backups must run daily with 30-day retention
- NFR-AVL-04: Point-in-time recovery must be supported

**NFR-USA: Usability Requirements**
- NFR-USA-01: UI must be responsive (desktop, tablet, mobile)
- NFR-USA-02: UI must meet WCAG 2.1 AA accessibility standards
- NFR-USA-03: Mobile apps must support offline mode for essential features
- NFR-USA-04: UI must support keyboard navigation
- NFR-USA-05: Error messages must be user-friendly and actionable

**NFR-INT: Integration Requirements**
- NFR-INT-01: REST API must follow OpenAPI 3.0 specification
- NFR-INT-02: Payment gateway integration must support Razorpay and PayU
- NFR-INT-03: SMS integration must support MSG91, TextLocal, Twilio
- NFR-INT-04: Email must support SMTP and SendGrid/SES
- NFR-INT-05: GPS tracking must integrate via standard protocols

**NFR-COM: Compliance Requirements**
- NFR-COM-01: System must comply with Indian IT Act 2000
- NFR-COM-02: DLT registration must be supported for SMS templates (India)
- NFR-COM-03: Data retention policies must be configurable per tenant
- NFR-COM-04: Audit logs must be immutable and retained for 7 years

**NFR-DEP: Deployment Requirements**
- NFR-DEP-01: On-premise deployment must use Docker Compose
- NFR-DEP-02: SaaS deployment must use Kubernetes
- NFR-DEP-03: CI/CD pipeline must include automated testing
- NFR-DEP-04: Blue-green deployments must be supported for zero-downtime

---

### Additional Requirements

#### From Architecture Document

**ARCH-STACK: Technology Stack Requirements**
- Backend: Go 1.23+ with Gin framework, GORM + sqlc
- Frontend: Angular 21 with standalone components, Signals (not NgRx)
- Database: PostgreSQL 16 with Row-Level Security
- Cache: Redis for sessions and caching
- File Storage: MinIO (on-premise) or S3 (cloud)
- Search: PostgreSQL full-text search (Phase 1), Elasticsearch (Phase 2+)

**ARCH-STRUCT: Project Structure Requirements**
- Backend modules: internal/modules/{module}/ with handler.go, service.go, repository.go, dto.go, entity.go
- Frontend features: src/app/features/{feature}/ with pages/, components/, services/, models/
- Shared components: src/app/shared/components/ for common UI elements
- Design system: Custom Angular Component Library (Tailwind CSS, no PrimeNG/Material)

**ARCH-API: API Design Requirements**
- API versioning: /api/v1/{resource}
- Cursor-based pagination for lists
- RFC 7807 error response format
- JSON camelCase for request/response fields
- UUID v7 for all entity IDs

**ARCH-DB: Database Requirements**
- Every table must have: id, tenant_id, created_at, updated_at, created_by, updated_by
- RLS policy on every table using current_setting('app.tenant_id')
- Index on tenant_id for all tables
- Migration naming: {timestamp}_{description}.up.sql/.down.sql

**ARCH-AUTH: Authentication Requirements**
- JWT with RS256 algorithm
- Access token: 15 minutes, Refresh token: 7 days
- Argon2id for password hashing
- TOTP 2FA for admin users
- OTP via SMS for parent/student login

#### From UX Design Document

**UX-RESP: Responsive Design Requirements**
- Device-first approach: Admin portal (desktop-first), Parent/Student portals (mobile-first)
- Breakpoints: sm (640px), md (768px), lg (1024px), xl (1280px), 2xl (1536px)
- Touch targets minimum 44x44px on mobile
- Bottom tab navigation on mobile apps

**UX-A11Y: Accessibility Requirements**
- WCAG 2.1 AA compliance
- Minimum contrast ratio 4.5:1 for text
- Focus indicators visible on all interactive elements
- Screen reader support with ARIA labels
- Skip to main content link
- Color not used as only indicator

**UX-COMP: Component Requirements**
- Custom component library: atoms (Button, Input, Badge, Avatar)
- Molecules (Form Field, Card, Dropdown, Toast, OTP Input)
- Organisms (Data Table, Attendance Grid, Fee Invoice Card, Navigation)
- Consistent loading states (skeleton screens)
- Toast notifications (success: green, error: red, warning: amber, info: blue)

**UX-PAT: Pattern Requirements**
- Button hierarchy: Primary (blue), Secondary (gray), Danger (red), Ghost
- Form validation: inline errors, touched state
- Navigation: desktop sidebar (collapsible), mobile bottom tabs
- Empty states with illustration, message, and action
- Search with debounce (300ms), clear button

---

### FR Coverage Map

| FR ID | Epic | Description |
|-------|------|-------------|
| FR-CF-01 to FR-CF-10 | Epic 1, Epic 2 | Core Foundation - Multi-tenancy, Auth, RBAC |
| FR-SM-01 to FR-SM-08 | Epic 4 | Student Management |
| FR-STF-01 to FR-STF-08 | Epic 5 | Staff Management |
| FR-ADM-01 to FR-ADM-07 | Epic 3 | Admissions |
| FR-AO-01 to FR-AO-05 | Epic 6 | Academic Structure & Timetable |
| FR-AO-06 to FR-AO-09 | Epic 7 | Attendance Operations |
| FR-EX-01 to FR-EX-09 | Epic 8 | Examinations & Grading |
| FR-HW-01 to FR-HW-07 | Epic 9 | Homework & Assignments |
| FR-LV-01 to FR-LV-07 | Epic 10 | Leave Management |
| FR-FE-01 to FR-FE-09 | Epic 11 | Fees & Payments |
| FR-CM-01 to FR-CM-07 | Epic 12 | Communication Hub |
| FR-PP-01 to FR-PP-08 | Epic 13 | Parent Portal |
| FR-SP-01 to FR-SP-07 | Epic 14 | Student Portal |
| FR-CT-01 to FR-CT-06 | Epic 15 | Certificate Generation |
| FR-QZ-01 to FR-QZ-07 | Epic 16 | Online Quiz & Assessment |
| FR-DC-01 to FR-DC-06 | Epic 17 | Digital Classroom |
| FR-TR-01 to FR-TR-07 | Epic 18 | Transport Management |
| FR-LB-01 to FR-LB-06 | Epic 19 | Library Management |
| FR-IN-01 to FR-IN-06 | Epic 20 | Inventory & Assets |
| FR-VG-01 to FR-VG-05 | Epic 21 | Visitor & Gate Management |
| FR-AN-01 to FR-AN-06 | Epic 22 | Analytics & Dashboards |
| FR-AI-01 to FR-AI-08 | Epic 23 | AI-Powered Insights |
| NFR-SEC, NFR-PER | Epic 1 | Infrastructure & Security Foundation |
| NFR-USA, UX-RESP, UX-A11Y, UX-COMP | Epic 1 | Design System & Accessibility |

---

## Epic List

### Phase 1 - MVP (15 Epics)

---

### Epic 1: Project Foundation & Design System
**Goal:** Establish the complete technical foundation including project scaffolding, database setup, authentication infrastructure, and reusable design system so that developers can build features consistently and securely.

**User Value:** Administrators can access a secure, professionally designed system with proper login and tenant isolation.

**FRs Covered:** FR-CF-01, FR-CF-06, FR-CF-08, FR-CF-09, NFR-SEC-01 to NFR-SEC-10, NFR-PER-01 to NFR-PER-06, UX-COMP, UX-PAT, ARCH-STACK, ARCH-STRUCT

**Implementation Notes:**
- Backend: Go project with Gin, GORM, database migrations, RLS policies
- Frontend: Angular 21 with Tailwind CSS, custom component library (atoms, molecules)
- Infrastructure: PostgreSQL 16, Redis, MinIO, Docker Compose
- Design System: Button, Input, Badge, Avatar, Card, Toast, Form Field components

---

### Epic 2: User Authentication & Access Control
**Goal:** Enable users to securely register, login, and access the system based on their roles with complete RBAC implementation.

**User Value:** Super Admins can create tenants, Admins can manage users and roles, all users can login securely with appropriate access levels.

**FRs Covered:** FR-CF-04, FR-CF-05, FR-CF-07, FR-CF-10, ARCH-AUTH

**Implementation Notes:**
- JWT RS256 authentication with refresh tokens
- Role-based permissions (Super Admin, Admin, Teacher, Parent, Student)
- TOTP 2FA for admin users, OTP for parents/students
- Feature flags and custom field configuration engine

---

### Epic 3: School Setup & Admissions
**Goal:** Enable school administrators to configure their school structure and manage the complete admission process from enquiry to enrollment.

**User Value:** Admins can set up branches, academic years, and process student admissions end-to-end.

**FRs Covered:** FR-CF-02, FR-CF-03, FR-ADM-01 to FR-ADM-07

**Implementation Notes:**
- Tenant/Branch configuration with academic year management
- Admission workflow: Enquiry → Application → Test → Decision → Enrollment
- Online application forms with document upload
- Merit list generation and admission fee tracking

---

### Epic 4: Student Lifecycle Management
**Goal:** Enable comprehensive management of student records throughout their entire school journey.

**User Value:** Admins and teachers can manage student profiles, track health records, handle behavioral incidents, and process promotions.

**FRs Covered:** FR-SM-01 to FR-SM-08

**Implementation Notes:**
- Student profiles with guardian information
- Health records (medical conditions, vaccinations, allergies)
- Behavioral tracking with incident management
- Document management with verification workflow
- Promotion/retention processing between academic years

---

### Epic 5: Staff Management
**Goal:** Enable complete management of teaching and non-teaching staff including profiles, attendance, and basic payroll.

**User Value:** HR and Admins can manage staff records, track qualifications, monitor attendance, and process salaries.

**FRs Covered:** FR-STF-01 to FR-STF-08

**Implementation Notes:**
- Staff profiles (teaching/non-teaching distinction)
- Qualifications and experience tracking
- Department and designation hierarchy
- Staff attendance with multiple marking sources
- Salary structure and payroll basics
- Subject-teacher assignments and workload tracking

---

### Epic 6: Academic Structure & Timetable
**Goal:** Enable schools to define their academic structure and create comprehensive timetables.

**User Value:** Admins can configure classes, sections, subjects, and teachers can view their schedules with substitution support.

**FRs Covered:** FR-AO-01 to FR-AO-05

**Implementation Notes:**
- Class and section management
- Subject configuration with class mappings
- Period slots with breaks/lunch
- Timetable builder for sections and teachers
- Substitution management for absent teachers

---

### Epic 7: Daily Attendance Operations
**Goal:** Enable teachers to efficiently mark and track student attendance with comprehensive reporting.

**User Value:** Teachers can mark attendance quickly, parents receive absence notifications, and admins can generate attendance reports.

**FRs Covered:** FR-AO-06 to FR-AO-09

**Implementation Notes:**
- Daily attendance marking interface (grid view)
- Optional period-wise attendance
- Attendance reports (daily, monthly, class-wise)
- Deadline and edit window enforcement
- Parent notifications for absences

---

### Epic 8: Examinations & Grading
**Goal:** Enable complete examination management from scheduling to report card generation.

**User Value:** Teachers can schedule exams, enter marks, and generate report cards; students and parents can view results.

**FRs Covered:** FR-EX-01 to FR-EX-09

**Implementation Notes:**
- Exam type configuration (unit test, term, final)
- Exam scheduling with hall ticket generation
- Marks entry with validation rules
- Grading scales (percentage and grade-based)
- Result calculation with ranking
- Report card generation with templates
- Attendance-based exam eligibility

---

### Epic 9: Homework & Assignments
**Goal:** Enable teachers to assign homework and track student submissions with grading.

**User Value:** Teachers can create assignments, students can submit online, and parents are notified of pending work.

**FRs Covered:** FR-HW-01 to FR-HW-07

**Implementation Notes:**
- Assignment creation (homework, project, worksheet)
- File attachments support
- Online submission with late penalty configuration
- Grading with rubric support
- Completion tracking and status reports
- Parent notifications for overdue work

---

### Epic 10: Leave Management
**Goal:** Enable staff to request leave and managers to approve with policy enforcement.

**User Value:** Staff can apply for leave online, approvers can review and approve, HR can track balances and generate reports.

**FRs Covered:** FR-LV-01 to FR-LV-07

**Implementation Notes:**
- Leave type configuration with policies
- Leave balance tracking per staff
- Approval workflow (single/multi-level)
- Policy enforcement (advance notice, max days, sandwich rule)
- Leave calendar and reports
- Status change notifications

---

### Epic 11: Fees & Payments
**Goal:** Enable complete fee management from structure definition to payment collection and receipt generation.

**User Value:** Admins can configure fees, parents can pay online, and accounts can track dues and generate reports.

**FRs Covered:** FR-FE-01 to FR-FE-09

**Implementation Notes:**
- Fee category and structure definition
- Student fee assignment with discounts
- Invoice generation with numbering
- Multi-mode payments (cash, cheque, UPI, card, netbanking)
- Payment gateway integration (Razorpay, PayU)
- Receipt generation (print/email)
- Late fee calculation
- Dues tracking and defaulter reports
- Refund processing with approval

---

### Epic 12: Communication Hub
**Goal:** Enable multi-channel communication between school and stakeholders.

**User Value:** Admins can send notices, parents receive SMS/email/push notifications, and teachers can message parents directly.

**FRs Covered:** FR-CM-01 to FR-CM-07

**Implementation Notes:**
- Notice creation with target audience selection
- SMS integration (MSG91, TextLocal, Twilio)
- Email notifications with templates
- Push notification support
- Notice acknowledgement tracking
- In-app parent-teacher messaging
- Automated event notifications

---

### Epic 13: Parent Portal
**Goal:** Enable parents to access all relevant information about their children through a dedicated portal.

**User Value:** Parents can view attendance, academics, fees, homework, notices, and communicate with teachers from one place.

**FRs Covered:** FR-PP-01 to FR-PP-08

**Implementation Notes:**
- Parent registration and login (OTP-based)
- Dashboard with child summary
- Attendance calendar and summary
- Academic performance and report cards
- Fee details and online payment
- Homework status and grades
- Notices and circulars
- Teacher messaging
- Multi-child support

---

### Epic 14: Student Portal
**Goal:** Enable students to access their academic information and submit assignments.

**User Value:** Students can view their timetable, attendance, submit homework, and check exam results.

**FRs Covered:** FR-SP-01 to FR-SP-07

**Implementation Notes:**
- Student login with age-appropriate access
- Personal timetable display
- Attendance records view
- Online homework submission
- Exam results and report cards
- Quiz/test assignments (view only in Phase 1)
- Library books view

---

### Epic 15: Certificate Generation
**Goal:** Enable certificate requests, approval workflow, and generation with verification.

**User Value:** Parents can request certificates, staff can process approvals, and generated certificates include QR verification.

**FRs Covered:** FR-CT-01 to FR-CT-06

**Implementation Notes:**
- Certificate request submission
- Approval workflow (clearances for TC)
- Template-based generation
- QR code for verification
- Delivery tracking and acknowledgement
- Duplicate certificate handling with fee

---

### Phase 2 - Extended Features (7 Epics)

---

### Epic 16: Online Quiz & Assessment
**Goal:** Enable teachers to create and conduct online assessments with auto-grading and analytics.

**User Value:** Teachers can build question banks, conduct timed quizzes, and students get instant feedback with detailed analytics.

**FRs Covered:** FR-QZ-01 to FR-QZ-07

**Implementation Notes:**
- Question bank (MCQ, true/false, fill-blank, match, short/long answer)
- Quiz creation (manual/auto-generated)
- Timed quizzes with shuffle options
- Auto-grading for objective questions
- Quiz analytics and question analysis
- Basic proctoring (tab detection, fullscreen)
- Practice mode with unlimited attempts

---

### Epic 17: Digital Classroom
**Goal:** Enable teachers to record classes and manage digital learning content.

**User Value:** Teachers can record lessons, students can replay anytime with bookmarks, and content library organizes all learning materials.

**FRs Covered:** FR-DC-01 to FR-DC-06

**Implementation Notes:**
- Class recording (screen + audio capture)
- Recording processing with quality options
- Bookmark and chapter markers
- Watch progress tracking
- Digital content library management
- Explained PDFs with page-level audio

---

### Epic 18: Transport Management
**Goal:** Enable complete school transport operations with GPS tracking.

**User Value:** Admins manage fleet and routes, parents track buses in real-time, and attendants mark transport attendance.

**FRs Covered:** FR-TR-01 to FR-TR-07

**Implementation Notes:**
- Vehicle management with document tracking
- Route definition with stops and timings
- Student-route assignment
- Driver and attendant management
- Real-time GPS tracking
- Transport attendance marking
- Distance-based fee calculation

---

### Epic 19: Library Management
**Goal:** Enable complete library operations from cataloging to circulation.

**User Value:** Librarians can manage books, students/staff can borrow and return, and fines are automatically calculated.

**FRs Covered:** FR-LB-01 to FR-LB-06

**Implementation Notes:**
- Book catalog with copy management
- Member management with borrowing limits
- Issue and return processing
- Reservation and waitlist handling
- Fine calculation and collection
- Circulation reports

---

### Epic 20: Inventory & Assets
**Goal:** Enable tracking of school assets and consumable inventory.

**User Value:** Admins can track all assets, request maintenance, manage stock levels, and process purchase orders.

**FRs Covered:** FR-IN-01 to FR-IN-06

**Implementation Notes:**
- Asset registration with depreciation tracking
- Location and assignment management
- Maintenance request processing
- Consumable stock management
- Purchase order workflow
- Vendor management

---

### Epic 21: Visitor & Gate Management
**Goal:** Enable secure visitor management and student gate pass processing.

**User Value:** Security can register visitors, parents can request gate passes, and all entries/exits are logged.

**FRs Covered:** FR-VG-01 to FR-VG-05

**Implementation Notes:**
- Visitor registration with photo capture
- Badge issuance and tracking
- Gate pass request and approval
- Vehicle entry/exit logging
- Pre-registration for expected visitors

---

### Epic 22: Analytics & Dashboards
**Goal:** Provide comprehensive analytics and customizable reporting across all modules.

**User Value:** Principals and admins get role-specific dashboards, can build custom reports, and schedule automated report delivery.

**FRs Covered:** FR-AN-01 to FR-AN-06

**Implementation Notes:**
- Role-based dashboards (Principal, Admin, Teacher)
- Academic analytics with trends
- Financial analytics
- Attendance analytics
- Custom report builder
- Scheduled report delivery via email

---

### Phase 3 - Advanced Features (1 Epic)

---

### Epic 23: AI-Powered Insights
**Goal:** Provide intelligent predictions, recommendations, and anomaly detection using on-premise AI.

**User Value:** Administrators get early warning of at-risk students, personalized learning recommendations, and automated content generation.

**FRs Covered:** FR-AI-01 to FR-AI-08

**Implementation Notes:**
- Student performance risk prediction
- Dropout risk identification
- Personalized learning recommendations
- Fee default prediction
- Admission demand forecasting
- Class recording auto-summarization
- Quiz question generation
- Attendance/grading anomaly detection
- On-premise AI (Llama/Mistral) for data privacy
