# 07 - Homework & Assignments

**Parent Document**: [index.md](./index.md)
**Status**: Draft
**Phase**: 1 (MVP)

---

## 1. Overview

The Homework & Assignments module enables teachers to create, distribute, and grade homework and assignments with online submission capabilities and parent visibility.

---

## 2. Assignment Types

### 2.1 Assignment Configuration

| Type | Description | Submission |
|------|-------------|------------|
| Homework | Daily practice work | Optional online |
| Assignment | Longer project work | Online required |
| Classwork | In-class activity | Teacher records |
| Project | Multi-day projects | Online with milestones |
| Worksheet | Printable worksheets | Physical or scan upload |

---

## 3. Assignment Entity

### 3.1 Core Assignment

**Entity: Assignment**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| branch_id | UUID | Branch reference |
| title | VARCHAR(255) | Assignment title |
| description | TEXT | Detailed instructions |
| assignment_type | ENUM | homework, assignment, project, worksheet, classwork |
| subject_id | UUID | Subject reference |
| class_id | UUID | Class reference |
| section_ids | JSONB | Target sections (NULL = all) |
| assigned_date | DATE | Assignment date |
| due_date | DATE | Due date |
| due_time | TIME | Due time |
| max_marks | DECIMAL | Maximum marks (NULL = ungraded) |
| weightage | DECIMAL | Grade weightage % |
| attachments | JSONB | Attached files |
| submission_type | ENUM | none, file, text, link, offline |
| allow_late | BOOLEAN | Allow late submission |
| late_penalty | DECIMAL | % deduction per day |
| max_late_days | INT | Maximum late days |
| is_group_work | BOOLEAN | Group assignment |
| group_size | INT | If group work |
| rubric_id | UUID | Grading rubric |
| status | ENUM | draft, published, closed, archived |
| created_by | UUID | Teacher who created |
| created_at | TIMESTAMP | Creation time |
| published_at | TIMESTAMP | When published |

### 3.2 Assignment Attachments

**Supported File Types**:
- Documents: PDF, DOC, DOCX, PPT, PPTX
- Images: JPG, PNG, GIF
- Spreadsheets: XLS, XLSX
- Links: YouTube, external URLs

---

## 4. Assignment Creation

### 4.1 Creation Interface

```
┌─────────────────────────────────────────────────────────────┐
│  CREATE NEW ASSIGNMENT                                       │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Title: [Chapter 5 - Practice Problems                  ]   │
│                                                              │
│  Type: [Homework     ▼]    Subject: [Mathematics    ▼]     │
│                                                              │
│  Class: [10         ▼]    Sections: [☑ A ☑ B ☐ C]         │
│                                                              │
│  Description:                                                │
│  ┌─────────────────────────────────────────────────────┐    │
│  │ Complete exercises 5.1 to 5.3 from NCERT textbook.  │    │
│  │ Show all working steps clearly.                      │    │
│  │                                                      │    │
│  │ Topics covered:                                      │    │
│  │ - Quadratic equations                                │    │
│  │ - Completing the square                              │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                              │
│  Attachments: [+ Add File] [+ Add Link]                     │
│  📎 worksheet.pdf (245 KB) [x]                              │
│                                                              │
│  Due Date: [25-Jan-2026]  Due Time: [09:00 AM]             │
│                                                              │
│  Submission: [File Upload ▼]                                │
│  ☑ Allow late submission (10% penalty per day, max 3 days) │
│                                                              │
│  Grading:                                                    │
│  Max Marks: [20]   Weightage: [5%]   Rubric: [None ▼]      │
│                                                              │
│  [Save Draft]  [Publish Now]  [Schedule for Later]          │
└─────────────────────────────────────────────────────────────┘
```

---

## 5. Student Submission

### 5.1 Submission Entity

**Entity: AssignmentSubmission**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| assignment_id | UUID | Assignment reference |
| student_id | UUID | Student reference |
| submission_date | TIMESTAMP | When submitted |
| is_late | BOOLEAN | Late submission |
| late_days | INT | Days late |
| status | ENUM | not_started, in_progress, submitted, returned, graded |
| content | TEXT | Text submission |
| attachments | JSONB | Uploaded files |
| links | JSONB | Submitted links |
| marks | DECIMAL | Marks obtained |
| late_penalty_applied | DECIMAL | Penalty deducted |
| final_marks | DECIMAL | After penalty |
| grade | VARCHAR(10) | Grade |
| feedback | TEXT | Teacher feedback |
| graded_by | UUID | Grading teacher |
| graded_at | TIMESTAMP | Grading time |
| resubmit_allowed | BOOLEAN | Can resubmit |
| resubmit_count | INT | Resubmission count |

### 5.2 Submission Interface (Student)

```
┌─────────────────────────────────────────────────────────────┐
│  ASSIGNMENT SUBMISSION                                       │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  📚 Chapter 5 - Practice Problems                           │
│  Subject: Mathematics | Due: 25-Jan-2026, 9:00 AM          │
│  Max Marks: 20                                               │
│                                                              │
│  Instructions:                                               │
│  Complete exercises 5.1 to 5.3 from NCERT textbook.        │
│  Show all working steps clearly.                            │
│                                                              │
│  📎 Attached: worksheet.pdf [Download]                      │
│                                                              │
├─────────────────────────────────────────────────────────────┤
│  YOUR SUBMISSION                                             │
│                                                              │
│  Upload Files: (Max 5 files, 10 MB each)                    │
│  [📁 Choose Files]                                          │
│                                                              │
│  Uploaded:                                                   │
│  📄 my_solution_page1.jpg (1.2 MB) [x]                      │
│  📄 my_solution_page2.jpg (1.4 MB) [x]                      │
│                                                              │
│  Or write here:                                              │
│  ┌─────────────────────────────────────────────────────┐    │
│  │                                                      │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                              │
│  ⏰ Time remaining: 2 days, 4 hours                         │
│                                                              │
│  [Save as Draft]  [Submit Assignment]                       │
└─────────────────────────────────────────────────────────────┘
```

---

## 6. Grading

### 6.1 Grading Rubric

**Entity: GradingRubric**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| name | VARCHAR(200) | Rubric name |
| subject_id | UUID | Subject (optional) |
| criteria | JSONB | Grading criteria |
| created_by | UUID | Creator |

**Rubric Example**:
```json
{
  "name": "Math Problem Solving Rubric",
  "max_marks": 20,
  "criteria": [
    {
      "name": "Understanding",
      "max_marks": 5,
      "levels": [
        { "score": 5, "description": "Complete understanding demonstrated" },
        { "score": 3, "description": "Partial understanding" },
        { "score": 1, "description": "Limited understanding" },
        { "score": 0, "description": "No understanding shown" }
      ]
    },
    {
      "name": "Methodology",
      "max_marks": 8,
      "levels": [
        { "score": 8, "description": "Correct method, all steps shown" },
        { "score": 5, "description": "Correct method, some steps missing" },
        { "score": 2, "description": "Partially correct method" },
        { "score": 0, "description": "Incorrect method" }
      ]
    },
    {
      "name": "Accuracy",
      "max_marks": 5,
      "levels": [
        { "score": 5, "description": "All answers correct" },
        { "score": 3, "description": "Most answers correct" },
        { "score": 1, "description": "Few answers correct" },
        { "score": 0, "description": "No correct answers" }
      ]
    },
    {
      "name": "Presentation",
      "max_marks": 2,
      "levels": [
        { "score": 2, "description": "Neat, well-organized" },
        { "score": 1, "description": "Acceptable presentation" },
        { "score": 0, "description": "Poor presentation" }
      ]
    }
  ]
}
```

### 6.2 Grading Interface (Teacher)

```
┌─────────────────────────────────────────────────────────────┐
│  GRADE SUBMISSIONS | Chapter 5 - Practice Problems          │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Class: 10-A | Submitted: 38/42 | Graded: 15/38            │
│                                                              │
│  Filter: [All ▼]  Sort: [Submission Time ▼]                │
│                                                              │
│  | # | Student        | Submitted    | Status  | Marks |   │
│  |---|----------------|--------------|---------|-------|   │
│  | 1 | Aarav Sharma   | 24-Jan 8:30  | Graded  | 18/20 |   │
│  | 2 | Ananya Patel   | 24-Jan 9:15  | Pending | -     |   │
│  | 3 | Arjun Singh    | 25-Jan 10:00 | Late    | -     |   │
│  | 4 | Diya Verma     | Not submitted| -       | -     |   │
│  ...                                                        │
│                                                              │
│  [Grade Selected] [Download All] [Export Grades]            │
└─────────────────────────────────────────────────────────────┘
```

### 6.3 Individual Grading

```
┌─────────────────────────────────────────────────────────────┐
│  GRADING: Ananya Patel                                      │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Submission: 24-Jan-2026, 9:15 AM (On time)                │
│                                                              │
│  Files: [📄 solution.pdf] [View] [Download]                 │
│                                                              │
│  ┌─────────────────────────────────────────────────────┐    │
│  │           [PDF/Image Viewer]                         │    │
│  │                                                      │    │
│  │    [Annotation tools: highlight, comment, draw]     │    │
│  │                                                      │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                              │
│  GRADING (Using Rubric)                                     │
│  | Criteria      | Score    | Max |                        │
│  |---------------|----------|-----|                        │
│  | Understanding | [5    ▼] | 5   |                        │
│  | Methodology   | [7    ▼] | 8   |                        │
│  | Accuracy      | [4    ▼] | 5   |                        │
│  | Presentation  | [2    ▼] | 2   |                        │
│  |---------------|----------|-----|                        │
│  | TOTAL         | 18       | 20  |                        │
│                                                              │
│  Feedback:                                                   │
│  ┌─────────────────────────────────────────────────────┐    │
│  │ Excellent work! Minor calculation error in Q3.      │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                              │
│  ☐ Allow resubmission                                       │
│                                                              │
│  [Save & Next] [Save & Close]                               │
└─────────────────────────────────────────────────────────────┘
```

---

## 7. Parent Visibility

### 7.1 Parent View Features

- See all assignments for their child
- View due dates and submission status
- Download assignment attachments
- View grades and teacher feedback
- Receive notifications for new assignments

### 7.2 Parent Dashboard Widget

```
┌─────────────────────────────────────────────────────────────┐
│  📝 HOMEWORK STATUS | Aarav (Class 10-A)                    │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  This Week:                                                  │
│  ✅ Math - Chapter 5 Practice (18/20) - Submitted          │
│  ⏳ English - Essay Writing - Due Tomorrow                  │
│  ❌ Science - Lab Report - Overdue by 2 days               │
│  📋 Hindi - Reading Assignment - Due in 3 days             │
│                                                              │
│  Summary: 1 Pending | 1 Overdue | 1 Graded                  │
│                                                              │
│  [View All Assignments]                                      │
└─────────────────────────────────────────────────────────────┘
```

---

## 8. Notifications

### 8.1 Notification Triggers

| Event | Recipients | Channels |
|-------|------------|----------|
| New assignment published | Students, Parents | Push, Email |
| Due date reminder (1 day) | Students, Parents | Push |
| Assignment graded | Student, Parents | Push, Email |
| Late submission warning | Student, Parents | Push |
| Resubmission allowed | Student | Push |

---

## 9. Reports

### 9.1 Class Assignment Report

```
Assignment Completion Report | Class 10-A | January 2026

| Subject      | Assigned | Avg Submission % | Avg Score |
|--------------|----------|------------------|-----------|
| Mathematics  | 8        | 92%              | 78%       |
| English      | 6        | 88%              | 72%       |
| Science      | 7        | 85%              | 75%       |
| Hindi        | 5        | 95%              | 80%       |
| Soc. Science | 4        | 90%              | 76%       |

Students with low completion (<70%):
- Ravi Kumar: 55% (missed 4 assignments)
- Priya Singh: 65% (missed 3 assignments)
```

### 9.2 Student Assignment History

```
Assignment History | Aarav Sharma | Class 10-A

| Date    | Subject | Assignment              | Status   | Marks |
|---------|---------|-------------------------|----------|-------|
| 24-Jan  | Math    | Chapter 5 Practice      | Graded   | 18/20 |
| 22-Jan  | English | Essay: My Hero          | Graded   | 15/20 |
| 20-Jan  | Science | Lab Report: Acids       | Pending  | -     |
| 18-Jan  | Hindi   | Grammar Exercises       | Graded   | 9/10  |
| 15-Jan  | Math    | Chapter 4 Problems      | Graded   | 17/20 |

Total Assignments: 25 | Submitted: 24 | Avg Score: 82%
```

---

## 10. API Endpoints

```
# Assignments
GET    /api/v1/assignments                  # List assignments
POST   /api/v1/assignments                  # Create assignment
GET    /api/v1/assignments/{id}             # Get assignment
PUT    /api/v1/assignments/{id}             # Update assignment
DELETE /api/v1/assignments/{id}             # Delete assignment
POST   /api/v1/assignments/{id}/publish     # Publish assignment

# Submissions
GET    /api/v1/assignments/{id}/submissions # List submissions
POST   /api/v1/assignments/{id}/submit      # Submit (student)
GET    /api/v1/submissions/{id}             # Get submission
PUT    /api/v1/submissions/{id}             # Update submission
POST   /api/v1/submissions/{id}/grade       # Grade submission

# Student View
GET    /api/v1/students/{id}/assignments    # Student's assignments
GET    /api/v1/students/{id}/assignments/pending # Pending assignments

# Reports
GET    /api/v1/assignments/reports/class    # Class report
GET    /api/v1/assignments/reports/student  # Student report
```

---

## 11. Business Rules

| Rule | Description |
|------|-------------|
| Submission Deadline | Cannot submit after due + max late days |
| Late Penalty | Auto-calculate and apply penalty |
| File Size Limit | Max 10 MB per file, 5 files per submission |
| Edit After Submit | Can edit until deadline only |
| Grade Visibility | Grades visible only after teacher publishes |
| Resubmission | Only if teacher allows |
| Draft Save | Auto-save drafts every 30 seconds |

---

## 12. Related Documents

- [04-academic-operations.md](./04-academic-operations.md) - Classes, subjects
- [08-online-quiz-assessment.md](./08-online-quiz-assessment.md) - Online tests
- [14-parent-portal.md](./14-parent-portal.md) - Parent view
- [index.md](./index.md) - Main PRD index

---

**Previous**: [06-examinations-grading.md](./06-examinations-grading.md)
**Next**: [08-online-quiz-assessment.md](./08-online-quiz-assessment.md)
