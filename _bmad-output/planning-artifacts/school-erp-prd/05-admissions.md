# 05 - Admissions

**Parent Document**: [index.md](./index.md)
**Status**: Draft
**Phase**: 1 (MVP)

---

## 1. Overview

The Admissions module manages the complete student admission lifecycle from initial enquiry to final enrollment, including application processing, document verification, fee collection, and seat allocation.

---

## 2. Admission Configuration

### 2.1 Admission Session

**Entity: AdmissionSession**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| branch_id | UUID | Branch reference |
| academic_year_id | UUID | Target academic year |
| name | VARCHAR(100) | Session name |
| start_date | DATE | Applications open |
| end_date | DATE | Applications close |
| status | ENUM | upcoming, open, closed, completed |
| settings | JSONB | Session-specific settings |

### 2.2 Class-wise Seats

**Entity: AdmissionSeat**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| session_id | UUID | Admission session |
| class_id | UUID | Class reference |
| total_seats | INT | Total available |
| reserved_seats | JSONB | Category-wise reservation |
| filled_seats | INT | Currently filled |
| waitlist_limit | INT | Max waitlist size |

**Reserved Seats Example**:
```json
{
  "general": 40,
  "ews": 10,
  "staff_ward": 5,
  "sibling": 5,
  "management": 5
}
```

---

## 3. Admission Workflow

### 3.1 Workflow Stages

```
┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐
│ Enquiry │───▶│  Apply  │───▶│ Review  │───▶│  Test   │
└─────────┘    └─────────┘    └─────────┘    └─────────┘
                                                  │
┌─────────┐    ┌─────────┐    ┌─────────┐        │
│Enrolled │◀───│ Pay Fee │◀───│Approved │◀───────┘
└─────────┘    └─────────┘    └─────────┘
                                  │
                            ┌─────────┐
                            │Waitlist │
                            └─────────┘
                                  │
                            ┌─────────┐
                            │Rejected │
                            └─────────┘
```

### 3.2 Workflow Configuration

**Entity: AdmissionWorkflow**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| class_level | ENUM | nursery, primary, middle, secondary |
| stages | JSONB | Workflow stages |
| is_active | BOOLEAN | Active workflow |

**Stage Configuration**:
```yaml
stages:
  - id: enquiry
    name: Enquiry
    required_fields: [student_name, class, parent_phone]
    next_stages: [application, cancelled]

  - id: application
    name: Application Submitted
    required_fields: [all_student_details, documents]
    required_documents: [birth_certificate, photo]
    application_fee: 500
    next_stages: [document_review, cancelled]

  - id: document_review
    name: Document Verification
    required_actions: [verify_documents]
    next_stages: [entrance_test, interview, approved, rejected]

  - id: entrance_test
    name: Entrance Test
    required_actions: [schedule_test, record_score]
    passing_score: 40
    next_stages: [interview, approved, waitlist, rejected]

  - id: interview
    name: Interview
    required_actions: [schedule_interview, record_feedback]
    next_stages: [approved, waitlist, rejected]

  - id: approved
    name: Approved
    required_actions: [send_offer_letter]
    next_stages: [fee_payment, cancelled]
    offer_validity_days: 7

  - id: fee_payment
    name: Fee Payment
    required_actions: [collect_admission_fee]
    next_stages: [enrolled, cancelled]

  - id: enrolled
    name: Enrolled
    final: true
    actions: [create_student_record, assign_section]
```

---

## 4. Enquiry Management

### 4.1 Enquiry Entity

**Entity: AdmissionEnquiry**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| branch_id | UUID | Branch reference |
| session_id | UUID | Admission session |
| enquiry_number | VARCHAR(50) | Unique enquiry ID |
| student_name | VARCHAR(200) | Student name |
| date_of_birth | DATE | DOB |
| gender | ENUM | male, female, other |
| class_applying | UUID | Class applying for |
| parent_name | VARCHAR(200) | Parent/Guardian name |
| parent_phone | VARCHAR(20) | Contact phone |
| parent_email | VARCHAR(255) | Contact email |
| source | ENUM | walk_in, website, referral, advertisement, other |
| referral_details | TEXT | If referred |
| remarks | TEXT | Enquiry notes |
| status | ENUM | new, contacted, converted, lost |
| follow_up_date | DATE | Next follow-up |
| assigned_to | UUID | Staff assigned |
| created_at | TIMESTAMP | Enquiry date |

### 4.2 Enquiry Follow-up

**Entity: EnquiryFollowUp**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| enquiry_id | UUID | Enquiry reference |
| follow_up_date | DATE | Follow-up date |
| contact_mode | ENUM | phone, email, visit, sms |
| notes | TEXT | Discussion notes |
| outcome | ENUM | interested, not_interested, callback, converted |
| next_follow_up | DATE | Next action date |
| created_by | UUID | Staff who followed up |

---

## 5. Application Process

### 5.1 Application Entity

**Entity: AdmissionApplication**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| branch_id | UUID | Branch reference |
| session_id | UUID | Admission session |
| enquiry_id | UUID | Source enquiry (if any) |
| application_number | VARCHAR(50) | Unique application ID |
| current_stage | VARCHAR(50) | Current workflow stage |
| stage_history | JSONB | Stage transitions |

### 5.2 Application - Student Details

| Field | Type | Description |
|-------|------|-------------|
| first_name | VARCHAR(100) | First name |
| middle_name | VARCHAR(100) | Middle name |
| last_name | VARCHAR(100) | Last name |
| date_of_birth | DATE | DOB |
| gender | ENUM | Gender |
| blood_group | VARCHAR(5) | Blood group |
| nationality | VARCHAR(50) | Nationality |
| religion | VARCHAR(50) | Religion |
| category | VARCHAR(50) | Category |
| aadhaar_number | VARCHAR(12) | Aadhaar |
| photo_url | VARCHAR(500) | Photo |
| previous_school | VARCHAR(255) | Previous school name |
| previous_class | VARCHAR(50) | Last class attended |
| previous_percentage | DECIMAL | Last exam percentage |

### 5.3 Application - Parent Details

Stored in related `ApplicationParent` entity with fields for:
- Father details (name, occupation, phone, email, education, income)
- Mother details (same fields)
- Guardian details (if applicable)

### 5.4 Application Documents

**Entity: ApplicationDocument**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| application_id | UUID | Application reference |
| document_type | VARCHAR(50) | Document type |
| file_url | VARCHAR(500) | File location |
| is_verified | BOOLEAN | Verified status |
| verified_by | UUID | Verifier |
| verified_at | TIMESTAMP | Verification time |
| rejection_reason | TEXT | If rejected |

---

## 6. Entrance Test

### 6.1 Test Schedule

**Entity: EntranceTest**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| session_id | UUID | Admission session |
| test_name | VARCHAR(200) | Test name |
| test_date | DATE | Test date |
| start_time | TIME | Start time |
| duration_minutes | INT | Duration |
| venue | VARCHAR(200) | Test venue |
| class_ids | JSONB | Applicable classes |
| max_candidates | INT | Capacity |
| status | ENUM | scheduled, in_progress, completed, cancelled |

### 6.2 Test Registration

**Entity: TestRegistration**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| test_id | UUID | Test reference |
| application_id | UUID | Application reference |
| roll_number | VARCHAR(20) | Test roll number |
| status | ENUM | registered, appeared, absent |
| marks_obtained | DECIMAL | Marks scored |
| max_marks | DECIMAL | Maximum marks |
| percentage | DECIMAL | Percentage |
| grade | VARCHAR(10) | Grade |
| remarks | TEXT | Examiner remarks |
| result | ENUM | pass, fail |

---

## 7. Interview

### 7.1 Interview Schedule

**Entity: AdmissionInterview**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| application_id | UUID | Application reference |
| interview_date | DATE | Interview date |
| interview_time | TIME | Interview time |
| interviewer_ids | JSONB | Panel members |
| venue | VARCHAR(200) | Location |
| status | ENUM | scheduled, completed, no_show, rescheduled |
| overall_rating | INT | 1-5 rating |
| feedback | TEXT | Interview feedback |
| recommendation | ENUM | strongly_recommend, recommend, neutral, not_recommend |

---

## 8. Admission Decision

### 8.1 Decision Entity

**Entity: AdmissionDecision**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| application_id | UUID | Application reference |
| decision | ENUM | approved, waitlisted, rejected |
| decision_date | DATE | Decision date |
| decided_by | UUID | Decision authority |
| section_assigned | UUID | If approved |
| waitlist_position | INT | If waitlisted |
| rejection_reason | TEXT | If rejected |
| offer_letter_sent | BOOLEAN | Offer sent |
| offer_valid_until | DATE | Offer expiry |
| remarks | TEXT | Decision remarks |

### 8.2 Offer Letter

```
┌─────────────────────────────────────────────────────────────┐
│                    ABC PUBLIC SCHOOL                         │
│              123 Main Street, Mumbai 400001                  │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│                    ADMISSION OFFER LETTER                    │
│                                                              │
│ Date: 22-Jan-2026                                           │
│ Application No: APP-2026-00123                              │
│                                                              │
│ Dear Mr. Rajesh Sharma,                                     │
│                                                              │
│ We are pleased to inform you that your ward,                │
│ AARAV SHARMA, has been selected for admission to            │
│ Class 5 for the Academic Year 2026-27.                      │
│                                                              │
│ Admission Details:                                          │
│ ─────────────────                                           │
│ Class: 5                                                    │
│ Section: To be assigned                                     │
│ Admission Fee: ₹10,000                                      │
│ First Quarter Fee: ₹15,000                                  │
│ Total Payable: ₹25,000                                      │
│                                                              │
│ Please complete the admission formalities by                │
│ 29-Jan-2026 to confirm your seat.                          │
│                                                              │
│ Required Documents:                                         │
│ • Original Birth Certificate                                │
│ • Transfer Certificate (if applicable)                      │
│ • 4 Passport Photos                                         │
│ • Aadhaar Card Copy                                         │
│                                                              │
│ Sincerely,                                                  │
│ Principal                                                   │
│ ABC Public School                                           │
└─────────────────────────────────────────────────────────────┘
```

---

## 9. Fee Payment & Enrollment

### 9.1 Admission Fee Collection

On approval, system generates admission invoice with:
- Admission fee
- First installment of tuition
- Caution deposit
- Other one-time fees

### 9.2 Enrollment Process

Upon fee payment:
1. Create student record in Student module
2. Link parent accounts
3. Assign section (manual or auto)
4. Generate admission number
5. Send welcome communication
6. Mark application as enrolled

---

## 10. Reports & Analytics

### 10.1 Admission Dashboard

```
Admission Dashboard | Session 2026-27

Applications Overview:
┌─────────────────────────────────────────────────────────────┐
│ Total Enquiries    │ Applications │ Approved │ Enrolled     │
│       450          │     320      │   180    │    150       │
└─────────────────────────────────────────────────────────────┘

Class-wise Status:
| Class | Seats | Applied | Approved | Enrolled | Waitlist | Vacant |
|-------|-------|---------|----------|----------|----------|--------|
| LKG   | 60    | 85      | 60       | 55       | 10       | 5      |
| UKG   | 60    | 45      | 42       | 40       | 0        | 20     |
| 1     | 90    | 120     | 85       | 75       | 15       | 15     |
| 2     | 30    | 25      | 22       | 20       | 0        | 10     |

Conversion Funnel:
Enquiry (450) → Application (320) → Test (280) → Approved (180) → Enrolled (150)
    │              71%                 88%           64%            83%
```

### 10.2 Source Analysis

```
Enquiry Sources | Session 2026-27

| Source        | Count | Converted | Conversion % |
|---------------|-------|-----------|--------------|
| Website       | 150   | 65        | 43%          |
| Walk-in       | 120   | 48        | 40%          |
| Referral      | 100   | 55        | 55%          |
| Advertisement | 50    | 15        | 30%          |
| Others        | 30    | 10        | 33%          |
```

---

## 11. API Endpoints

```
# Enquiries
GET    /api/v1/admissions/enquiries         # List enquiries
POST   /api/v1/admissions/enquiries         # Create enquiry
GET    /api/v1/admissions/enquiries/{id}    # Get enquiry
PUT    /api/v1/admissions/enquiries/{id}    # Update enquiry
POST   /api/v1/admissions/enquiries/{id}/follow-up # Add follow-up

# Applications
GET    /api/v1/admissions/applications      # List applications
POST   /api/v1/admissions/applications      # Submit application
GET    /api/v1/admissions/applications/{id} # Get application
PUT    /api/v1/admissions/applications/{id} # Update application
POST   /api/v1/admissions/applications/{id}/documents # Upload document
POST   /api/v1/admissions/applications/{id}/verify # Verify documents

# Tests & Interviews
POST   /api/v1/admissions/tests             # Create test
POST   /api/v1/admissions/tests/{id}/register # Register candidate
POST   /api/v1/admissions/tests/{id}/results # Submit results
POST   /api/v1/admissions/interviews        # Schedule interview
PUT    /api/v1/admissions/interviews/{id}   # Update interview

# Decisions
POST   /api/v1/admissions/applications/{id}/decide # Make decision
POST   /api/v1/admissions/applications/{id}/enroll # Complete enrollment

# Reports
GET    /api/v1/admissions/reports/dashboard # Dashboard stats
GET    /api/v1/admissions/reports/funnel    # Conversion funnel
```

---

## 12. Business Rules

| Rule | Description |
|------|-------------|
| Age Eligibility | Class-wise age criteria validation |
| Document Mandatory | Cannot proceed without required documents |
| Seat Availability | Cannot approve beyond available seats |
| Offer Validity | Auto-expire offers after validity period |
| Duplicate Check | Check for duplicate applications (phone/email) |
| Waitlist Order | FIFO for waitlist promotions |
| Fee Before Enrollment | Cannot enroll without fee payment |

---

## 13. Related Documents

- [03-student-management.md](./03-student-management.md) - Student records
- [12-fees-payments.md](./12-fees-payments.md) - Fee collection
- [index.md](./index.md) - Main PRD index

---

**Previous**: [04-academic-operations.md](./04-academic-operations.md)
**Next**: [06-examinations-grading.md](./06-examinations-grading.md)
