# 03 - Student Management

**Parent Document**: [index.md](./index.md)
**Status**: Draft
**Phase**: 1 (MVP)

---

## 1. Overview

The Student Management module handles the complete student lifecycle from enrollment to graduation, including profile management, health records, discipline tracking, and document management.

---

## 2. Student Profile

### 2.1 Core Student Entity

**Entity: Student**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| branch_id | UUID | Branch reference |
| user_id | UUID | Linked user account |
| admission_number | VARCHAR(50) | Unique admission number |
| roll_number | VARCHAR(20) | Class roll number |
| first_name | VARCHAR(100) | First name |
| middle_name | VARCHAR(100) | Middle name |
| last_name | VARCHAR(100) | Last name |
| date_of_birth | DATE | DOB |
| gender | ENUM | male, female, other |
| blood_group | VARCHAR(5) | Blood group |
| nationality | VARCHAR(50) | Nationality |
| religion | VARCHAR(50) | Religion |
| caste | VARCHAR(50) | Caste (optional) |
| category | VARCHAR(50) | Category (General, OBC, SC, ST) |
| mother_tongue | VARCHAR(50) | Mother tongue |
| aadhaar_number | VARCHAR(12) | Aadhaar (encrypted) |
| photo_url | VARCHAR(500) | Profile photo |
| status | ENUM | active, inactive, alumni, transferred, dropped |
| created_at | TIMESTAMP | Record creation |
| updated_at | TIMESTAMP | Last update |

### 2.2 Contact Information

**Entity: StudentContact**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| student_id | UUID | Student reference |
| address_type | ENUM | permanent, current, correspondence |
| address_line1 | VARCHAR(255) | Address line 1 |
| address_line2 | VARCHAR(255) | Address line 2 |
| city | VARCHAR(100) | City |
| state | VARCHAR(100) | State |
| country | VARCHAR(100) | Country |
| pincode | VARCHAR(20) | PIN code |
| phone | VARCHAR(20) | Phone number |
| email | VARCHAR(255) | Email address |

### 2.3 Parent/Guardian Information

**Entity: StudentGuardian**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| student_id | UUID | Student reference |
| user_id | UUID | Linked parent user (if registered) |
| relation | ENUM | father, mother, guardian, other |
| first_name | VARCHAR(100) | First name |
| last_name | VARCHAR(100) | Last name |
| occupation | VARCHAR(100) | Occupation |
| organization | VARCHAR(200) | Employer/Business |
| annual_income | DECIMAL | Annual income |
| education | VARCHAR(100) | Education qualification |
| phone | VARCHAR(20) | Phone number |
| email | VARCHAR(255) | Email |
| is_primary | BOOLEAN | Primary contact |
| can_pickup | BOOLEAN | Authorized for pickup |
| receives_notifications | BOOLEAN | Receive notifications |

---

## 3. Academic Enrollment

### 3.1 Class Enrollment

**Entity: StudentEnrollment**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| student_id | UUID | Student reference |
| academic_year_id | UUID | Academic year |
| class_id | UUID | Class reference |
| section_id | UUID | Section reference |
| roll_number | VARCHAR(20) | Roll number for this year |
| enrollment_date | DATE | Enrollment date |
| status | ENUM | enrolled, promoted, detained, transferred, dropped |
| promoted_from_id | UUID | Previous enrollment (if promoted) |
| remarks | TEXT | Remarks |

### 3.2 Subject Enrollment

**Entity: StudentSubject**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| enrollment_id | UUID | Enrollment reference |
| subject_id | UUID | Subject reference |
| is_optional | BOOLEAN | Optional subject |
| enrolled_at | TIMESTAMP | When enrolled |
| dropped_at | TIMESTAMP | If dropped |

---

## 4. Health Records

### 4.1 Health Profile

**Entity: StudentHealth**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| student_id | UUID | Student reference |
| height_cm | DECIMAL | Height in cm |
| weight_kg | DECIMAL | Weight in kg |
| bmi | DECIMAL | Calculated BMI |
| vision_left | VARCHAR(20) | Left eye vision |
| vision_right | VARCHAR(20) | Right eye vision |
| hearing_status | VARCHAR(50) | Hearing status |
| dental_status | VARCHAR(50) | Dental status |
| allergies | TEXT | Known allergies |
| chronic_conditions | TEXT | Chronic conditions |
| medications | TEXT | Current medications |
| special_needs | TEXT | Special needs/accommodations |
| emergency_contact | VARCHAR(20) | Emergency phone |
| doctor_name | VARCHAR(100) | Family doctor |
| doctor_phone | VARCHAR(20) | Doctor contact |
| last_checkup_date | DATE | Last health checkup |
| recorded_at | TIMESTAMP | Record date |
| recorded_by | UUID | Recorded by user |

### 4.2 Immunization Records

**Entity: StudentImmunization**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| student_id | UUID | Student reference |
| vaccine_name | VARCHAR(100) | Vaccine name |
| dose_number | INT | Dose number |
| administered_date | DATE | Date given |
| administered_by | VARCHAR(200) | Doctor/Hospital |
| next_due_date | DATE | Next dose due |
| certificate_url | VARCHAR(500) | Certificate file |

---

## 5. Discipline Management

### 5.1 Behavior Records

**Entity: StudentBehavior**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| student_id | UUID | Student reference |
| enrollment_id | UUID | Current enrollment |
| incident_date | DATE | Incident date |
| incident_type | ENUM | positive, negative, neutral |
| category | VARCHAR(50) | Category (see below) |
| description | TEXT | Incident description |
| action_taken | TEXT | Action taken |
| points | INT | Behavior points (+/-) |
| reported_by | UUID | Reporting staff |
| approved_by | UUID | Approving authority |
| parent_notified | BOOLEAN | Parent informed |
| notified_at | TIMESTAMP | When notified |
| attachments | JSONB | Evidence files |

**Behavior Categories**:
```yaml
positive:
  - academic_excellence
  - sports_achievement
  - helpful_behavior
  - leadership
  - community_service

negative:
  - late_coming
  - uniform_violation
  - misconduct
  - bullying
  - property_damage
  - academic_dishonesty
```

### 5.2 Discipline Actions

**Entity: DisciplineAction**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| behavior_id | UUID | Related behavior record |
| action_type | ENUM | warning, detention, suspension, expulsion |
| start_date | DATE | Action start |
| end_date | DATE | Action end |
| conditions | TEXT | Conditions for return |
| status | ENUM | active, completed, revoked |
| decision_by | UUID | Decision authority |
| decision_date | DATE | Decision date |

---

## 6. Student Documents

### 6.1 Document Management

**Entity: StudentDocument**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| student_id | UUID | Student reference |
| document_type | VARCHAR(50) | Document type |
| document_name | VARCHAR(200) | Display name |
| file_url | VARCHAR(500) | File location |
| file_size | INT | Size in bytes |
| mime_type | VARCHAR(100) | File type |
| is_verified | BOOLEAN | Verified by admin |
| verified_by | UUID | Verifier |
| verified_at | TIMESTAMP | Verification date |
| uploaded_at | TIMESTAMP | Upload date |
| uploaded_by | UUID | Uploader |
| expiry_date | DATE | Document expiry (if applicable) |

**Document Types**:
```yaml
required:
  - birth_certificate
  - previous_marksheet
  - transfer_certificate
  - aadhar_card
  - passport_photo

optional:
  - caste_certificate
  - income_certificate
  - medical_certificate
  - character_certificate
  - migration_certificate
```

---

## 7. Student Diary/Remarks

### 7.1 Daily Remarks

**Entity: StudentDiary**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| student_id | UUID | Student reference |
| enrollment_id | UUID | Current enrollment |
| entry_date | DATE | Entry date |
| remark_type | ENUM | general, homework, behavior, achievement |
| subject_id | UUID | Related subject (optional) |
| content | TEXT | Diary content |
| is_important | BOOLEAN | Flag for important |
| requires_acknowledgement | BOOLEAN | Needs parent response |
| acknowledged_at | TIMESTAMP | When acknowledged |
| acknowledged_by | UUID | Parent who acknowledged |
| created_by | UUID | Teacher who wrote |
| created_at | TIMESTAMP | Creation time |

---

## 8. Student ID Cards

### 8.1 ID Card Generation

**Entity: StudentIDCard**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| student_id | UUID | Student reference |
| academic_year_id | UUID | Academic year |
| card_number | VARCHAR(50) | Unique card number |
| barcode | VARCHAR(100) | Barcode data |
| qr_code | VARCHAR(500) | QR code data |
| issue_date | DATE | Issue date |
| expiry_date | DATE | Expiry date |
| status | ENUM | active, lost, replaced, expired |
| printed_at | TIMESTAMP | Print date |
| printed_by | UUID | Printed by |

**ID Card Template Data**:
```json
{
  "school_name": "ABC Public School",
  "school_logo": "url",
  "student_photo": "url",
  "student_name": "John Doe",
  "class_section": "10-A",
  "admission_number": "2024001",
  "father_name": "James Doe",
  "contact": "9876543210",
  "address": "123 Main St",
  "blood_group": "O+",
  "valid_until": "2025-03-31"
}
```

---

## 9. Bulk Operations

### 9.1 Supported Bulk Operations

| Operation | Description |
|-----------|-------------|
| Import Students | CSV/Excel import for new students |
| Promote Students | Bulk promotion to next class |
| Generate Roll Numbers | Auto-generate roll numbers |
| Update Section | Bulk section change |
| Export Data | Export student data to Excel |
| Generate ID Cards | Bulk ID card generation |
| Send Notifications | Bulk parent notifications |

### 9.2 Import Template

```csv
admission_number,first_name,last_name,date_of_birth,gender,class,section,father_name,father_phone,mother_name,mother_phone,address,city,pincode
2024001,John,Doe,2010-05-15,male,5,A,James Doe,9876543210,Jane Doe,9876543211,123 Main St,Mumbai,400001
```

---

## 10. API Endpoints

```
# Student CRUD
GET    /api/v1/students                     # List students (paginated)
POST   /api/v1/students                     # Create student
GET    /api/v1/students/{id}                # Get student details
PUT    /api/v1/students/{id}                # Update student
DELETE /api/v1/students/{id}                # Soft delete

# Student Profile Sections
GET    /api/v1/students/{id}/guardians      # List guardians
POST   /api/v1/students/{id}/guardians      # Add guardian
GET    /api/v1/students/{id}/health         # Get health records
PUT    /api/v1/students/{id}/health         # Update health records
GET    /api/v1/students/{id}/documents      # List documents
POST   /api/v1/students/{id}/documents      # Upload document

# Enrollment
GET    /api/v1/students/{id}/enrollments    # Enrollment history
POST   /api/v1/students/{id}/enroll         # Enroll in class

# Discipline
GET    /api/v1/students/{id}/behavior       # Behavior records
POST   /api/v1/students/{id}/behavior       # Add behavior record

# Diary
GET    /api/v1/students/{id}/diary          # Get diary entries
POST   /api/v1/students/{id}/diary          # Add diary entry

# Bulk Operations
POST   /api/v1/students/import              # Bulk import
POST   /api/v1/students/promote             # Bulk promote
POST   /api/v1/students/export              # Export to Excel

# Search
GET    /api/v1/students/search?q=           # Search students
```

---

## 11. Business Rules

| Rule | Description |
|------|-------------|
| Unique Admission Number | Admission number must be unique within tenant |
| Age Validation | DOB must make student eligible for class |
| Guardian Required | At least one guardian must be registered |
| Photo Required | Photo required before ID card generation |
| Document Verification | Critical documents need admin verification |
| Enrollment Continuity | Cannot enroll in lower class than previous |

---

## 12. Related Documents

- [04-academic-operations.md](./04-academic-operations.md) - Classes, attendance
- [05-admissions.md](./05-admissions.md) - Admission workflow
- [14-parent-portal.md](./14-parent-portal.md) - Parent access
- [index.md](./index.md) - Main PRD index

---

**Previous**: [02-core-foundation.md](./02-core-foundation.md)
**Next**: [04-academic-operations.md](./04-academic-operations.md)
