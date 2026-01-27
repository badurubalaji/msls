# 10 - Staff Management

**Parent Document**: [index.md](./index.md)
**Status**: Draft
**Phase**: 1 (MVP)

---

## 1. Overview

The Staff Management module handles the complete employee lifecycle including onboarding, profile management, attendance tracking, and payroll integration.

---

## 2. Staff Profile

### 2.1 Core Staff Entity

**Entity: Staff**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| branch_id | UUID | Primary branch |
| user_id | UUID | Linked user account |
| employee_id | VARCHAR(50) | Employee ID |
| first_name | VARCHAR(100) | First name |
| middle_name | VARCHAR(100) | Middle name |
| last_name | VARCHAR(100) | Last name |
| date_of_birth | DATE | DOB |
| gender | ENUM | male, female, other |
| blood_group | VARCHAR(5) | Blood group |
| nationality | VARCHAR(50) | Nationality |
| marital_status | ENUM | single, married, divorced, widowed |
| aadhaar_number | VARCHAR(12) | Aadhaar (encrypted) |
| pan_number | VARCHAR(10) | PAN |
| photo_url | VARCHAR(500) | Profile photo |
| staff_type | ENUM | teaching, non_teaching, admin, support |
| designation_id | UUID | Designation |
| department_id | UUID | Department |
| joining_date | DATE | Date of joining |
| confirmation_date | DATE | Confirmation date |
| employment_type | ENUM | permanent, contract, probation, part_time |
| reporting_to | UUID | Reporting manager |
| status | ENUM | active, on_leave, resigned, terminated, retired |
| resignation_date | DATE | If resigned |
| last_working_date | DATE | Last working day |
| created_at | TIMESTAMP | Record creation |

### 2.2 Contact Information

**Entity: StaffContact**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| staff_id | UUID | Staff reference |
| address_type | ENUM | permanent, current |
| address_line1 | VARCHAR(255) | Address |
| city | VARCHAR(100) | City |
| state | VARCHAR(100) | State |
| pincode | VARCHAR(20) | PIN |
| phone_primary | VARCHAR(20) | Primary phone |
| phone_secondary | VARCHAR(20) | Secondary phone |
| email_personal | VARCHAR(255) | Personal email |
| email_official | VARCHAR(255) | Official email |
| emergency_contact_name | VARCHAR(200) | Emergency contact |
| emergency_contact_phone | VARCHAR(20) | Emergency phone |
| emergency_contact_relation | VARCHAR(50) | Relation |

### 2.3 Family Details

**Entity: StaffFamily**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| staff_id | UUID | Staff reference |
| relation | ENUM | spouse, child, parent, sibling |
| name | VARCHAR(200) | Name |
| date_of_birth | DATE | DOB |
| occupation | VARCHAR(100) | Occupation |
| is_dependent | BOOLEAN | Financial dependent |

---

## 3. Organizational Structure

### 3.1 Department

**Entity: Department**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| name | VARCHAR(100) | Department name |
| code | VARCHAR(20) | Short code |
| head_id | UUID | Department head (staff) |
| parent_id | UUID | Parent department |
| is_active | BOOLEAN | Active status |

**Example Departments**:
```
├── Academic
│   ├── Primary Section
│   ├── Middle Section
│   └── Senior Section
├── Administration
│   ├── Accounts
│   ├── HR
│   └── IT
├── Support Services
│   ├── Library
│   ├── Transport
│   └── Maintenance
└── Sports & Activities
```

### 3.2 Designation

**Entity: Designation**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| name | VARCHAR(100) | Designation name |
| code | VARCHAR(20) | Short code |
| level | INT | Hierarchy level |
| department_id | UUID | Default department |
| is_teaching | BOOLEAN | Teaching role |

**Example Designations**:
```yaml
teaching:
  - Principal (Level 1)
  - Vice Principal (Level 2)
  - Head of Department (Level 3)
  - Senior Teacher (Level 4)
  - Teacher (Level 5)
  - Assistant Teacher (Level 6)

non_teaching:
  - Administrative Officer
  - Accountant
  - Librarian
  - Lab Assistant
  - Office Assistant
  - Peon
```

---

## 4. Qualifications & Experience

### 4.1 Educational Qualification

**Entity: StaffQualification**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| staff_id | UUID | Staff reference |
| qualification | VARCHAR(100) | Degree name |
| specialization | VARCHAR(100) | Specialization |
| institution | VARCHAR(200) | Institution name |
| board_university | VARCHAR(200) | Board/University |
| year_of_passing | INT | Passing year |
| percentage_grade | VARCHAR(20) | Percentage/Grade |
| certificate_url | VARCHAR(500) | Certificate file |
| is_verified | BOOLEAN | Verified status |

### 4.2 Work Experience

**Entity: StaffExperience**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| staff_id | UUID | Staff reference |
| organization | VARCHAR(200) | Organization name |
| designation | VARCHAR(100) | Designation held |
| from_date | DATE | Start date |
| to_date | DATE | End date |
| responsibilities | TEXT | Key responsibilities |
| reason_for_leaving | TEXT | Leaving reason |
| reference_name | VARCHAR(200) | Reference contact |
| reference_phone | VARCHAR(20) | Reference phone |
| experience_letter_url | VARCHAR(500) | Document |
| is_verified | BOOLEAN | Verified status |

---

## 5. Staff Attendance

### 5.1 Attendance Entity

**Entity: StaffAttendance**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| staff_id | UUID | Staff reference |
| date | DATE | Attendance date |
| check_in | TIME | Check-in time |
| check_out | TIME | Check-out time |
| status | ENUM | present, absent, half_day, leave, holiday, week_off |
| work_hours | DECIMAL | Hours worked |
| overtime_hours | DECIMAL | Overtime hours |
| late_minutes | INT | Late by minutes |
| early_leave_minutes | INT | Left early by |
| remarks | TEXT | Remarks |
| source | ENUM | manual, biometric, mobile, system |
| marked_by | UUID | If manual |

### 5.2 Attendance Rules

**Entity: AttendancePolicy**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| name | VARCHAR(100) | Policy name |
| work_start_time | TIME | Expected start |
| work_end_time | TIME | Expected end |
| grace_minutes | INT | Late grace period |
| half_day_threshold | DECIMAL | Hours for half day |
| full_day_threshold | DECIMAL | Hours for full day |
| overtime_threshold | DECIMAL | After this = overtime |
| week_offs | JSONB | Weekly offs [0,6] |
| is_default | BOOLEAN | Default policy |

### 5.3 Attendance Report

```
Staff Attendance Report | January 2026

Employee: Rajesh Kumar | ID: EMP001 | Designation: Senior Teacher

| Metric          | Value |
|-----------------|-------|
| Working Days    | 24    |
| Present         | 22    |
| Absent          | 0     |
| Leaves          | 2     |
| Late Days       | 3     |
| Half Days       | 1     |
| Overtime Hours  | 8     |

Attendance %: 95.8%

Leave Details:
- 15-Jan: Casual Leave
- 22-Jan: Sick Leave

Late Days:
- 05-Jan: 12 min late
- 12-Jan: 8 min late
- 18-Jan: 15 min late
```

---

## 6. Payroll

### 6.1 Salary Structure

**Entity: SalaryStructure**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| staff_id | UUID | Staff reference |
| effective_from | DATE | Effective date |
| basic_salary | DECIMAL | Basic pay |
| hra | DECIMAL | House rent allowance |
| da | DECIMAL | Dearness allowance |
| ta | DECIMAL | Travel allowance |
| medical | DECIMAL | Medical allowance |
| special_allowance | DECIMAL | Special allowance |
| other_allowances | JSONB | Other components |
| gross_salary | DECIMAL | Gross total |
| pf_employee | DECIMAL | PF deduction |
| pf_employer | DECIMAL | PF contribution |
| esi_employee | DECIMAL | ESI deduction |
| professional_tax | DECIMAL | Professional tax |
| tds | DECIMAL | TDS |
| other_deductions | JSONB | Other deductions |
| net_salary | DECIMAL | Net payable |
| ctc | DECIMAL | Cost to company |

### 6.2 Monthly Payroll

**Entity: Payroll**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| staff_id | UUID | Staff reference |
| month | INT | Month (1-12) |
| year | INT | Year |
| working_days | INT | Total working days |
| days_worked | DECIMAL | Days present |
| lop_days | DECIMAL | Loss of pay days |
| basic_earned | DECIMAL | Basic earned |
| gross_earned | DECIMAL | Gross earned |
| total_deductions | DECIMAL | Total deductions |
| net_payable | DECIMAL | Net salary |
| status | ENUM | draft, processed, approved, paid |
| payment_date | DATE | Payment date |
| payment_reference | VARCHAR(100) | Transaction ref |
| payslip_url | VARCHAR(500) | Payslip PDF |

### 6.3 Payslip Format

```
┌─────────────────────────────────────────────────────────────┐
│                    ABC PUBLIC SCHOOL                         │
│                   SALARY SLIP - JAN 2026                     │
├─────────────────────────────────────────────────────────────┤
│  Employee: RAJESH KUMAR         Employee ID: EMP001         │
│  Designation: Senior Teacher    Department: Academic         │
│  PAN: XXXXX1234X                UAN: 100XXXXXXXX            │
│  Bank A/C: XXXX4567             Working Days: 24            │
├─────────────────────────────────────────────────────────────┤
│  EARNINGS                        DEDUCTIONS                  │
│  ──────────────────              ──────────────────          │
│  Basic         ₹35,000          PF            ₹4,200        │
│  HRA           ₹14,000          ESI           ₹525          │
│  DA            ₹3,500           Prof. Tax     ₹200          │
│  TA            ₹2,000           TDS           ₹2,500        │
│  Medical       ₹1,500           LOP (0 days)  ₹0            │
│  Special       ₹4,000                                        │
│  ──────────────────              ──────────────────          │
│  Gross        ₹60,000           Total Ded.   ₹7,425         │
├─────────────────────────────────────────────────────────────┤
│  NET PAYABLE: ₹52,575                                       │
│  (Rupees Fifty Two Thousand Five Hundred Seventy Five Only) │
├─────────────────────────────────────────────────────────────┤
│  Leave Balance: CL: 8 | SL: 5 | EL: 12                      │
└─────────────────────────────────────────────────────────────┘
```

---

## 7. Teacher-Specific Features

### 7.1 Subject Assignment

**Entity: TeacherSubject**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| staff_id | UUID | Teacher reference |
| academic_year_id | UUID | Academic year |
| class_id | UUID | Class |
| section_id | UUID | Section |
| subject_id | UUID | Subject |
| is_class_teacher | BOOLEAN | Class teacher for this section |
| periods_per_week | INT | Assigned periods |

### 7.2 Workload Summary

```
Teacher Workload | Rajesh Kumar | 2025-26

Subjects Assigned:
| Class | Section | Subject     | Periods/Week |
|-------|---------|-------------|--------------|
| 10    | A       | Mathematics | 6            |
| 10    | B       | Mathematics | 6            |
| 9     | A       | Mathematics | 5            |
| 9     | B       | Mathematics | 5            |

Total Periods: 22/week
Class Teacher: 10-A
Free Periods: 8/week
```

---

## 8. Documents

### 8.1 Staff Documents

**Entity: StaffDocument**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| staff_id | UUID | Staff reference |
| document_type | VARCHAR(50) | Document type |
| document_name | VARCHAR(200) | Display name |
| file_url | VARCHAR(500) | File location |
| expiry_date | DATE | Expiry date |
| is_verified | BOOLEAN | Verified status |
| uploaded_at | TIMESTAMP | Upload time |

**Document Types**:
- ID Proof (Aadhaar, PAN, Passport)
- Address Proof
- Educational Certificates
- Experience Letters
- Professional Certifications
- Medical Fitness Certificate
- Police Verification

---

## 9. API Endpoints

```
# Staff CRUD
GET    /api/v1/staff                        # List staff
POST   /api/v1/staff                        # Create staff
GET    /api/v1/staff/{id}                   # Get staff details
PUT    /api/v1/staff/{id}                   # Update staff
DELETE /api/v1/staff/{id}                   # Deactivate staff

# Profile Sections
GET    /api/v1/staff/{id}/qualifications    # Qualifications
POST   /api/v1/staff/{id}/qualifications    # Add qualification
GET    /api/v1/staff/{id}/experience        # Experience
POST   /api/v1/staff/{id}/experience        # Add experience
GET    /api/v1/staff/{id}/documents         # Documents
POST   /api/v1/staff/{id}/documents         # Upload document

# Attendance
GET    /api/v1/staff/attendance             # Attendance list
POST   /api/v1/staff/attendance             # Mark attendance
GET    /api/v1/staff/{id}/attendance        # Staff attendance history

# Payroll
GET    /api/v1/staff/{id}/salary-structure  # Salary structure
PUT    /api/v1/staff/{id}/salary-structure  # Update salary
GET    /api/v1/payroll                      # Monthly payroll
POST   /api/v1/payroll/process              # Process payroll
GET    /api/v1/staff/{id}/payslips          # Payslip history

# Teacher Specific
GET    /api/v1/staff/{id}/subjects          # Assigned subjects
POST   /api/v1/staff/{id}/subjects          # Assign subject
GET    /api/v1/staff/{id}/timetable         # Teacher timetable
```

---

## 10. Business Rules

| Rule | Description |
|------|-------------|
| Unique Employee ID | Employee ID unique within tenant |
| Joining Date | Cannot be future date |
| Reporting Hierarchy | Cannot report to self or subordinate |
| Salary Revision | New structure doesn't affect processed payroll |
| Attendance Lock | Cannot edit attendance after payroll processed |
| Document Expiry | Alert 30 days before document expires |
| Exit Clearance | Clearance required before final settlement |

---

## 11. Related Documents

- [11-leave-management.md](./11-leave-management.md) - Leave system
- [04-academic-operations.md](./04-academic-operations.md) - Timetable
- [index.md](./index.md) - Main PRD index

---

**Previous**: [09-digital-classroom.md](./09-digital-classroom.md)
**Next**: [11-leave-management.md](./11-leave-management.md)
