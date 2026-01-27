# 06 - Examinations & Grading

**Parent Document**: [index.md](./index.md)
**Status**: Draft
**Phase**: 1 (MVP)

---

## 1. Overview

The Examinations & Grading module manages the complete examination lifecycle including exam scheduling, hall ticket generation, marks entry, grade calculation, and report card generation.

---

## 2. Examination Structure

### 2.1 Exam Type Configuration

**Entity: ExamType**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| name | VARCHAR(100) | Exam type name |
| code | VARCHAR(20) | Short code |
| weightage | DECIMAL | Percentage in final grade |
| is_term_exam | BOOLEAN | Term/final exam |
| display_order | INT | Ordering |

**Common Exam Types**:
```yaml
exam_types:
  - name: Unit Test 1
    code: UT1
    weightage: 10
  - name: Unit Test 2
    code: UT2
    weightage: 10
  - name: Half Yearly
    code: HY
    weightage: 30
    is_term_exam: true
  - name: Unit Test 3
    code: UT3
    weightage: 10
  - name: Unit Test 4
    code: UT4
    weightage: 10
  - name: Annual
    code: AN
    weightage: 30
    is_term_exam: true
```

### 2.2 Exam Schedule

**Entity: Examination**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| branch_id | UUID | Branch reference |
| academic_year_id | UUID | Academic year |
| exam_type_id | UUID | Exam type |
| name | VARCHAR(200) | Exam name |
| start_date | DATE | Exam period start |
| end_date | DATE | Exam period end |
| classes | JSONB | Applicable classes |
| status | ENUM | scheduled, in_progress, completed, results_published |
| created_by | UUID | Creator |

### 2.3 Exam Timetable

**Entity: ExamSchedule**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| examination_id | UUID | Exam reference |
| class_id | UUID | Class |
| subject_id | UUID | Subject |
| exam_date | DATE | Exam date |
| start_time | TIME | Start time |
| end_time | TIME | End time |
| max_marks | DECIMAL | Maximum marks |
| passing_marks | DECIMAL | Minimum to pass |
| venue | VARCHAR(200) | Exam hall |

**Exam Timetable View**:
```
Half Yearly Examination 2025-26 | Class 10

| Date       | Day       | Subject     | Time          | Max Marks |
|------------|-----------|-------------|---------------|-----------|
| 15-Sep     | Monday    | English     | 9:00 - 12:00  | 80        |
| 16-Sep     | Tuesday   | Hindi       | 9:00 - 12:00  | 80        |
| 17-Sep     | Wednesday | Mathematics | 9:00 - 12:00  | 80        |
| 18-Sep     | Thursday  | Science     | 9:00 - 12:00  | 80        |
| 19-Sep     | Friday    | Soc. Science| 9:00 - 12:00  | 80        |

Note: Practical exams will be conducted separately.
```

---

## 3. Hall Ticket / Admit Card

### 3.1 Hall Ticket Entity

**Entity: HallTicket**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| examination_id | UUID | Exam reference |
| student_id | UUID | Student reference |
| roll_number | VARCHAR(20) | Exam roll number |
| status | ENUM | generated, downloaded, issued |
| generated_at | TIMESTAMP | Generation time |
| downloaded_at | TIMESTAMP | Download time |

### 3.2 Hall Ticket Format

```
┌─────────────────────────────────────────────────────────────┐
│                    ABC PUBLIC SCHOOL                         │
│                       ADMIT CARD                             │
│              Half Yearly Examination 2025-26                 │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────┐  Name: AARAV SHARMA                           │
│  │          │  Class: X-A                                   │
│  │  PHOTO   │  Roll No: 2024001                             │
│  │          │  Father's Name: RAJESH SHARMA                 │
│  └──────────┘                                               │
│                                                              │
├─────────────────────────────────────────────────────────────┤
│  EXAMINATION SCHEDULE                                        │
│                                                              │
│  | Date    | Subject      | Time        | Room |            │
│  |---------|--------------|-------------|------|            │
│  | 15-Sep  | English      | 9:00-12:00  | H-1  |            │
│  | 16-Sep  | Hindi        | 9:00-12:00  | H-1  |            │
│  | 17-Sep  | Mathematics  | 9:00-12:00  | H-2  |            │
│  | 18-Sep  | Science      | 9:00-12:00  | H-2  |            │
│  | 19-Sep  | Soc. Science | 9:00-12:00  | H-1  |            │
│                                                              │
├─────────────────────────────────────────────────────────────┤
│  Instructions:                                               │
│  1. Reach exam hall 30 minutes before                       │
│  2. Carry this admit card and school ID                     │
│  3. No electronic devices allowed                           │
│                                                              │
│  ________________          ________________                  │
│  Parent Signature          Principal                         │
└─────────────────────────────────────────────────────────────┘
```

---

## 4. Marks Entry

### 4.1 Marks Entity

**Entity: ExamMarks**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| exam_schedule_id | UUID | Exam schedule reference |
| student_id | UUID | Student reference |
| marks_obtained | DECIMAL | Marks scored |
| is_absent | BOOLEAN | Was absent |
| is_exempted | BOOLEAN | Exempted |
| exemption_reason | TEXT | Reason for exemption |
| practical_marks | DECIMAL | Practical marks (if any) |
| internal_marks | DECIMAL | Internal assessment |
| total_marks | DECIMAL | Calculated total |
| grade | VARCHAR(5) | Grade (if grading) |
| remarks | TEXT | Teacher remarks |
| entered_by | UUID | Data entry person |
| entered_at | TIMESTAMP | Entry time |
| verified_by | UUID | Verified by |
| verified_at | TIMESTAMP | Verification time |
| is_locked | BOOLEAN | Locked for editing |

### 4.2 Marks Entry Interface

```
Marks Entry | Half Yearly | Class 10-A | Mathematics

Max Marks: 80 | Passing: 26 | Practical: 20 | Internal: 0

| Roll | Name           | Theory | Practical | Total | Grade | Status |
|------|----------------|--------|-----------|-------|-------|--------|
| 001  | Aarav Sharma   | [65 ]  | [18]      | 83    | A     | Pass   |
| 002  | Ananya Patel   | [72 ]  | [19]      | 91    | A+    | Pass   |
| 003  | Arjun Singh    | [AB ]  | [--]      | AB    | -     | Absent |
| 004  | Diya Verma     | [45 ]  | [15]      | 60    | B     | Pass   |
| 005  | Ishaan Kumar   | [22 ]  | [12]      | 34    | D     | Fail   |
...

[Save Draft] [Submit for Verification] [Cancel]

Summary: Entered: 35/40 | Absent: 2 | Pending: 3
```

### 4.3 Bulk Marks Upload

**CSV Format**:
```csv
roll_number,theory_marks,practical_marks,remarks
001,65,18,Good performance
002,72,19,Excellent
003,AB,,Absent - Medical
004,45,15,
005,22,12,Needs improvement
```

---

## 5. Grading System

### 5.1 Grade Configuration

**Entity: GradeScale**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| name | VARCHAR(100) | Scale name |
| type | ENUM | percentage, cgpa, letter |
| is_default | BOOLEAN | Default scale |

**Entity: GradeDefinition**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| grade_scale_id | UUID | Parent scale |
| grade | VARCHAR(5) | Grade symbol |
| min_percentage | DECIMAL | Minimum % |
| max_percentage | DECIMAL | Maximum % |
| grade_point | DECIMAL | Grade point (for CGPA) |
| description | VARCHAR(100) | Grade description |
| display_order | INT | Ordering |

**CBSE Grading Scale**:
```yaml
grades:
  - grade: A1
    range: 91-100
    point: 10
    description: Outstanding
  - grade: A2
    range: 81-90
    point: 9
    description: Excellent
  - grade: B1
    range: 71-80
    point: 8
    description: Very Good
  - grade: B2
    range: 61-70
    point: 7
    description: Good
  - grade: C1
    range: 51-60
    point: 6
    description: Above Average
  - grade: C2
    range: 41-50
    point: 5
    description: Average
  - grade: D
    range: 33-40
    point: 4
    description: Below Average
  - grade: E
    range: 0-32
    point: 0
    description: Needs Improvement
```

---

## 6. Result Calculation

### 6.1 Result Entity

**Entity: ExamResult**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| examination_id | UUID | Exam reference |
| student_id | UUID | Student reference |
| total_marks | DECIMAL | Total marks obtained |
| max_marks | DECIMAL | Maximum possible |
| percentage | DECIMAL | Percentage |
| grade | VARCHAR(5) | Overall grade |
| cgpa | DECIMAL | CGPA (if applicable) |
| rank | INT | Class rank |
| division | VARCHAR(20) | First/Second/Third |
| status | ENUM | pass, fail, compartment |
| subjects_failed | JSONB | List of failed subjects |
| remarks | TEXT | Result remarks |
| is_published | BOOLEAN | Published to parents |

### 6.2 Result Calculation Logic

```
For each student:
1. Fetch all subject marks
2. Calculate subject-wise total (theory + practical + internal)
3. Determine subject-wise pass/fail
4. Calculate overall percentage
5. Assign overall grade
6. Calculate CGPA (if applicable)
7. Determine rank within class
8. Set overall result status:
   - PASS: All subjects passed
   - FAIL: More than 2 subjects failed
   - COMPARTMENT: 1-2 subjects failed
```

---

## 7. Report Card

### 7.1 Report Card Entity

**Entity: ReportCard**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| examination_id | UUID | Exam reference |
| student_id | UUID | Student reference |
| template_id | UUID | Report card template |
| generated_at | TIMESTAMP | Generation time |
| pdf_url | VARCHAR(500) | Generated PDF |
| is_printed | BOOLEAN | Print status |
| is_distributed | BOOLEAN | Distribution status |

### 7.2 Report Card Format

```
┌─────────────────────────────────────────────────────────────┐
│                    ABC PUBLIC SCHOOL                         │
│                 PROGRESS REPORT 2025-26                      │
│                 Half Yearly Examination                      │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Student: AARAV SHARMA          Admission No: 2024001       │
│  Class: X-A                     Roll No: 15                 │
│  DOB: 15-May-2010               Father: RAJESH SHARMA       │
│                                                              │
├─────────────────────────────────────────────────────────────┤
│                    SCHOLASTIC AREAS                          │
│                                                              │
│  | Subject      | Max  | Theory | Prac | Total | Grade |   │
│  |--------------|------|--------|------|-------|-------|   │
│  | English      | 100  | 68     | -    | 68    | B1    |   │
│  | Hindi        | 100  | 72     | -    | 72    | B1    |   │
│  | Mathematics  | 100  | 65     | 18   | 83    | A2    |   │
│  | Science      | 100  | 58     | 17   | 75    | B1    |   │
│  | Soc. Science | 100  | 70     | -    | 70    | B1    |   │
│  |--------------|------|--------|------|-------|-------|   │
│  | TOTAL        | 500  |        |      | 368   |       |   │
│                                                              │
│  Percentage: 73.6%    Grade: B1    Rank: 8/42               │
│                                                              │
├─────────────────────────────────────────────────────────────┤
│                   CO-SCHOLASTIC AREAS                        │
│                                                              │
│  | Area                    | Grade |                        │
│  |-------------------------|-------|                        │
│  | Work Education          | A     |                        │
│  | Art Education           | B     |                        │
│  | Physical Education      | A     |                        │
│  | Discipline              | A     |                        │
│                                                              │
├─────────────────────────────────────────────────────────────┤
│  Attendance: 92/100 days (92%)                              │
│                                                              │
│  Class Teacher's Remarks: Good progress. Keep it up!        │
│                                                              │
│  Principal's Signature: ____________                         │
│  Parent's Signature: ____________                            │
│                                                              │
│  Result: PASSED                Next Exam: Annual (March)    │
└─────────────────────────────────────────────────────────────┘
```

---

## 8. Result Publication

### 8.1 Publication Workflow

```
All marks entered
       │
       ▼
Verification by subject teachers
       │
       ▼
Verification by class teacher
       │
       ▼
Approval by exam coordinator
       │
       ▼
Generate report cards
       │
       ▼
Principal approval
       │
       ▼
Publish results
       │
       ├──▶ Parent portal notification
       ├──▶ SMS/Email to parents
       └──▶ Report cards available for download
```

### 8.2 Publication Controls

- **Lock marks**: Prevent editing after deadline
- **Review period**: Allow corrections before publishing
- **Staged publication**: Publish class-by-class
- **Re-evaluation**: Handle re-checking requests

---

## 9. Analytics

### 9.1 Class Performance

```
Half Yearly Results | Class 10-A | Mathematics

Total Students: 42
Appeared: 40
Passed: 35 (87.5%)
Failed: 5 (12.5%)

Grade Distribution:
| Grade | Count | % |
|-------|-------|---|
| A1    | 3     | 8% |
| A2    | 8     | 20% |
| B1    | 12    | 30% |
| B2    | 7     | 18% |
| C1    | 5     | 13% |
| D     | 3     | 8% |
| E     | 2     | 5% |

Class Average: 62.4%
Highest: 95 (Priya S.)
Lowest: 28 (Ravi K.)
```

### 9.2 Subject Comparison

```
Subject Performance | Half Yearly | Class 10-A

| Subject      | Avg % | Pass % | A1-A2 % | E % |
|--------------|-------|--------|---------|-----|
| English      | 68.5  | 92     | 25      | 3   |
| Mathematics  | 62.4  | 88     | 28      | 5   |
| Science      | 65.2  | 90     | 22      | 4   |
| Hindi        | 72.1  | 95     | 32      | 2   |
| Soc. Science | 70.8  | 94     | 30      | 2   |
```

---

## 10. API Endpoints

```
# Exam Configuration
GET    /api/v1/exams/types                  # List exam types
POST   /api/v1/examinations                 # Create examination
GET    /api/v1/examinations/{id}            # Get examination
PUT    /api/v1/examinations/{id}            # Update examination
POST   /api/v1/examinations/{id}/schedule   # Add schedule

# Hall Tickets
POST   /api/v1/examinations/{id}/hall-tickets/generate # Generate
GET    /api/v1/students/{id}/hall-tickets   # Get student's hall tickets

# Marks Entry
GET    /api/v1/exams/{scheduleId}/marks     # Get marks sheet
POST   /api/v1/exams/{scheduleId}/marks     # Enter marks
PUT    /api/v1/exams/{scheduleId}/marks/{studentId} # Update marks
POST   /api/v1/exams/{scheduleId}/marks/upload # Bulk upload
POST   /api/v1/exams/{scheduleId}/marks/lock # Lock marks

# Results
POST   /api/v1/examinations/{id}/calculate  # Calculate results
GET    /api/v1/examinations/{id}/results    # Get results
POST   /api/v1/examinations/{id}/publish    # Publish results
GET    /api/v1/students/{id}/results        # Student result history

# Report Cards
POST   /api/v1/examinations/{id}/report-cards/generate # Generate
GET    /api/v1/students/{id}/report-cards   # Get report cards
```

---

## 11. Business Rules

| Rule | Description |
|------|-------------|
| Marks Range | Cannot enter marks > max marks |
| Absent Handling | Absent students get 0 or AB status |
| Edit Window | Marks editable only until lock date |
| Verification Required | Two-level verification before publishing |
| Pass Criteria | Subject-wise and overall pass criteria |
| Rank Calculation | Based on total marks, ties share rank |
| Result Once Published | Cannot unpublish without admin override |

---

## 12. Related Documents

- [04-academic-operations.md](./04-academic-operations.md) - Classes, subjects
- [08-online-quiz-assessment.md](./08-online-quiz-assessment.md) - Online tests
- [14-parent-portal.md](./14-parent-portal.md) - Result viewing
- [index.md](./index.md) - Main PRD index

---

**Previous**: [05-admissions.md](./05-admissions.md)
**Next**: [07-homework-assignments.md](./07-homework-assignments.md)
