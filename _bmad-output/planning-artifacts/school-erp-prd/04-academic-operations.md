# 04 - Academic Operations

**Parent Document**: [index.md](./index.md)
**Status**: Draft
**Phase**: 1 (MVP)

---

## 1. Overview

Academic Operations covers the core academic structure including classes, sections, subjects, timetable management, and attendance tracking.

---

## 2. Class & Section Management

### 2.1 Class Entity

**Entity: Class**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| branch_id | UUID | Branch reference |
| name | VARCHAR(50) | Class name (e.g., "Class 5", "Grade 10") |
| code | VARCHAR(20) | Short code (e.g., "5", "10") |
| display_order | INT | Ordering |
| level | ENUM | nursery, primary, middle, secondary, senior_secondary |
| is_active | BOOLEAN | Active status |

### 2.2 Section Entity

**Entity: Section**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| branch_id | UUID | Branch reference |
| class_id | UUID | Parent class |
| name | VARCHAR(20) | Section name (e.g., "A", "B", "Rose") |
| capacity | INT | Max students |
| class_teacher_id | UUID | Class teacher (staff) |
| room_number | VARCHAR(20) | Default room |
| is_active | BOOLEAN | Active status |

### 2.3 Class Structure

```
Academic Year 2025-26
├── Nursery
│   ├── Section A (Teacher: Ms. Priya)
│   └── Section B (Teacher: Ms. Anita)
├── LKG
│   ├── Section A
│   └── Section B
├── UKG
│   └── Section A
├── Class 1
│   ├── Section A
│   ├── Section B
│   └── Section C
...
├── Class 12 Science
│   ├── Section A
│   └── Section B
└── Class 12 Commerce
    └── Section A
```

---

## 3. Subject Management

### 3.1 Subject Entity

**Entity: Subject**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| name | VARCHAR(100) | Subject name |
| code | VARCHAR(20) | Subject code (e.g., "MATH", "ENG") |
| type | ENUM | core, elective, language, co_curricular |
| credit_hours | DECIMAL | Credit hours (for higher classes) |
| is_graded | BOOLEAN | Included in grading |
| is_active | BOOLEAN | Active status |

### 3.2 Class-Subject Mapping

**Entity: ClassSubject**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| class_id | UUID | Class reference |
| subject_id | UUID | Subject reference |
| academic_year_id | UUID | Academic year |
| is_mandatory | BOOLEAN | Mandatory for all |
| periods_per_week | INT | Weekly periods |
| teacher_id | UUID | Default teacher |

---

## 4. Timetable Management

### 4.1 Period Definition

**Entity: PeriodSlot**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| branch_id | UUID | Branch reference |
| name | VARCHAR(50) | Period name (e.g., "Period 1", "Lunch") |
| type | ENUM | class, break, assembly, lunch |
| start_time | TIME | Start time |
| end_time | TIME | End time |
| display_order | INT | Ordering |
| applicable_days | JSONB | [mon, tue, wed, thu, fri, sat] |

**Example Period Structure**:
```yaml
periods:
  - name: Assembly
    type: assembly
    start: "08:00"
    end: "08:30"
  - name: Period 1
    type: class
    start: "08:30"
    end: "09:15"
  - name: Period 2
    type: class
    start: "09:15"
    end: "10:00"
  - name: Short Break
    type: break
    start: "10:00"
    end: "10:15"
  - name: Period 3
    type: class
    start: "10:15"
    end: "11:00"
  # ... continues
```

### 4.2 Timetable Entry

**Entity: Timetable**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| academic_year_id | UUID | Academic year |
| section_id | UUID | Section reference |
| day_of_week | INT | 0=Sun, 1=Mon, ... 6=Sat |
| period_slot_id | UUID | Period slot |
| subject_id | UUID | Subject |
| teacher_id | UUID | Assigned teacher |
| room_id | UUID | Room (optional) |
| is_active | BOOLEAN | Active entry |
| effective_from | DATE | Start date |
| effective_until | DATE | End date (NULL = ongoing) |

### 4.3 Timetable View (Section)

```
Section: Class 5-A | Class Teacher: Ms. Sharma

| Time        | Monday  | Tuesday | Wednesday | Thursday | Friday  | Saturday |
|-------------|---------|---------|-----------|----------|---------|----------|
| 08:00-08:30 | Assembly| Assembly| Assembly  | Assembly | Assembly| Assembly |
| 08:30-09:15 | English | Math    | Science   | Hindi    | Math    | English  |
| 09:15-10:00 | Math    | English | Math      | English  | Science | PT       |
| 10:00-10:15 | Break   | Break   | Break     | Break    | Break   | Break    |
| 10:15-11:00 | Hindi   | Science | English   | Math     | Hindi   | Art      |
| 11:00-11:45 | Science | Hindi   | Hindi     | Science  | SST     | Music    |
| 11:45-12:30 | Lunch   | Lunch   | Lunch     | Lunch    | Lunch   | -        |
| 12:30-01:15 | SST     | PT      | Comp      | Art      | Comp    | -        |
| 01:15-02:00 | Comp    | SST     | SST       | Comp     | GK      | -        |
```

### 4.4 Teacher Timetable View

```
Teacher: Mr. Rajesh Kumar | Subject: Mathematics

| Time        | Monday    | Tuesday   | Wednesday | Thursday  | Friday    |
|-------------|-----------|-----------|-----------|-----------|-----------|
| 08:30-09:15 | 5-A       | 6-B       | 5-B       | 6-A       | 5-A       |
| 09:15-10:00 | 6-A       | 5-A       | 6-B       | 5-B       | 6-A       |
| 10:15-11:00 | Free      | Free      | 5-A       | Free      | 5-B       |
| 11:00-11:45 | 5-B       | 6-A       | Free      | 6-B       | 6-B       |
| 12:30-01:15 | 6-B       | Free      | 6-A       | 5-A       | Free      |

Total Periods: 24/week | Free Periods: 6/week
```

### 4.5 Substitution Management

**Entity: Substitution**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| timetable_id | UUID | Original timetable entry |
| date | DATE | Substitution date |
| original_teacher_id | UUID | Absent teacher |
| substitute_teacher_id | UUID | Substitute teacher |
| reason | TEXT | Reason for substitution |
| created_by | UUID | Created by |
| created_at | TIMESTAMP | Creation time |
| status | ENUM | pending, confirmed, completed, cancelled |

---

## 5. Attendance Management

### 5.1 Student Attendance

**Entity: StudentAttendance**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| student_id | UUID | Student reference |
| enrollment_id | UUID | Current enrollment |
| date | DATE | Attendance date |
| status | ENUM | present, absent, late, half_day, holiday, leave |
| marked_at | TIMESTAMP | When marked |
| marked_by | UUID | Marked by (teacher) |
| remarks | TEXT | Optional remarks |
| late_minutes | INT | Minutes late (if late) |
| leave_id | UUID | Reference to leave (if on leave) |

### 5.2 Attendance Status Types

| Status | Description | Counts As |
|--------|-------------|-----------|
| present | Full day present | Present |
| absent | Absent without leave | Absent |
| late | Came late (within grace) | Present |
| half_day | Present for half day | 0.5 Present |
| holiday | School holiday | Not counted |
| leave | Approved leave | Leave (separate count) |

### 5.3 Period-wise Attendance (Optional)

**Entity: PeriodAttendance**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| student_id | UUID | Student reference |
| date | DATE | Attendance date |
| period_slot_id | UUID | Period slot |
| subject_id | UUID | Subject |
| status | ENUM | present, absent |
| marked_by | UUID | Teacher who marked |
| marked_at | TIMESTAMP | When marked |

### 5.4 Attendance Marking Interface

```
Class: 5-A | Date: 2026-01-22 | Total: 35

| # | Roll | Name           | Status  | Remarks      |
|---|------|----------------|---------|--------------|
| 1 | 001  | Aarav Sharma   | Present | ✓            |
| 2 | 002  | Ananya Patel   | Late    | 10 min late  |
| 3 | 003  | Arjun Singh    | Absent  |              |
| 4 | 004  | Diya Verma     | Leave   | Medical      |
| 5 | 005  | Ishaan Kumar   | Present | ✓            |
...

Summary: Present: 30 | Absent: 3 | Late: 1 | Leave: 1

[Mark All Present] [Submit] [Cancel]
```

### 5.5 Attendance Reports

**Daily Summary**:
```
Date: 2026-01-22 | Branch: Main Campus

| Class   | Total | Present | Absent | Late | Leave | % |
|---------|-------|---------|--------|------|-------|---|
| Class 1 | 120   | 110     | 5      | 3    | 2     | 94% |
| Class 2 | 115   | 108     | 4      | 2    | 1     | 96% |
| Class 3 | 110   | 100     | 6      | 2    | 2     | 93% |
...
| Total   | 1200  | 1120    | 45     | 20   | 15    | 95% |
```

**Monthly Student Report**:
```
Student: Aarav Sharma | Class: 5-A | Month: January 2026

Working Days: 24
Present: 22
Absent: 1
Late: 1
Leave: 0

Attendance %: 91.67%

Absent Days: 15-Jan (Sick)
Late Days: 22-Jan (10 min)
```

---

## 6. Biometric Integration (Optional)

### 6.1 Device Registration

**Entity: BiometricDevice**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| branch_id | UUID | Branch reference |
| device_name | VARCHAR(100) | Device name |
| device_serial | VARCHAR(100) | Serial number |
| device_type | ENUM | fingerprint, face, card |
| location | VARCHAR(100) | Physical location |
| ip_address | VARCHAR(45) | Device IP |
| is_active | BOOLEAN | Active status |
| last_sync_at | TIMESTAMP | Last sync time |

### 6.2 Biometric Log

**Entity: BiometricLog**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| device_id | UUID | Device reference |
| user_id | UUID | User (student/staff) |
| user_type | ENUM | student, staff |
| punch_time | TIMESTAMP | Punch timestamp |
| punch_type | ENUM | in, out |
| verification_type | ENUM | fingerprint, face, card |
| is_processed | BOOLEAN | Converted to attendance |

---

## 7. API Endpoints

### Classes & Sections
```
GET    /api/v1/classes                      # List classes
POST   /api/v1/classes                      # Create class
GET    /api/v1/classes/{id}/sections        # List sections
POST   /api/v1/classes/{id}/sections        # Create section
GET    /api/v1/sections/{id}/students       # Students in section
```

### Subjects
```
GET    /api/v1/subjects                     # List subjects
POST   /api/v1/subjects                     # Create subject
GET    /api/v1/classes/{id}/subjects        # Class subjects
POST   /api/v1/classes/{id}/subjects        # Assign subjects
```

### Timetable
```
GET    /api/v1/timetable/section/{id}       # Section timetable
GET    /api/v1/timetable/teacher/{id}       # Teacher timetable
POST   /api/v1/timetable                    # Create entry
PUT    /api/v1/timetable/{id}               # Update entry
POST   /api/v1/timetable/generate           # Auto-generate
POST   /api/v1/substitutions                # Create substitution
```

### Attendance
```
GET    /api/v1/attendance/section/{id}      # Section attendance
POST   /api/v1/attendance/mark              # Mark attendance
PUT    /api/v1/attendance/{id}              # Update attendance
GET    /api/v1/attendance/student/{id}      # Student history
GET    /api/v1/attendance/report/daily      # Daily report
GET    /api/v1/attendance/report/monthly    # Monthly report
```

---

## 8. Business Rules

| Rule | Description |
|------|-------------|
| Attendance Deadline | Attendance must be marked by 11 AM |
| Late Grace Period | 15 minutes grace before marking late |
| Half Day Threshold | Less than 4 hours = half day |
| Edit Window | Can edit attendance within 24 hours |
| Timetable Conflict | Teacher cannot have overlapping periods |
| Room Conflict | Room cannot have double booking |
| Minimum Working Days | 75% attendance required for exam eligibility |

---

## 9. Related Documents

- [03-student-management.md](./03-student-management.md) - Student profiles
- [06-examinations-grading.md](./06-examinations-grading.md) - Exams
- [10-staff-management.md](./10-staff-management.md) - Staff
- [index.md](./index.md) - Main PRD index

---

**Previous**: [03-student-management.md](./03-student-management.md)
**Next**: [05-admissions.md](./05-admissions.md)
