# 15 - Student Portal

**Parent Document**: [index.md](./index.md)
**Status**: Draft
**Phase**: 1 (MVP)

---

## 1. Overview

The Student Portal provides students with access to their academic resources, assignments, results, and learning materials through web and mobile interfaces.

---

## 2. Portal Access

### 2.1 Login Methods

- School Email + Password
- Student ID + Password
- OTP-based (for password reset)

### 2.2 Age-Appropriate Access

| Level | Access |
|-------|--------|
| Primary (1-5) | Limited features, parent-supervised |
| Middle (6-8) | Standard features |
| Secondary (9-12) | Full features |

---

## 3. Dashboard

```
┌─────────────────────────────────────────────────────────────┐
│  Hi Aarav! 👋                           [🔔 2] [⚙️]         │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Class 10-A | Roll No: 15 | 2025-26                         │
│                                                              │
│  📅 TODAY'S SCHEDULE                                        │
│  ├─ 08:30 Mathematics - Mr. Kumar                          │
│  ├─ 09:15 English - Ms. Sharma                             │
│  ├─ 10:15 Science - Mr. Patel                              │
│  └─ [View Full Timetable]                                   │
│                                                              │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐           │
│  │📝 HOMEWORK │ │📚 QUIZZES  │ │🎬 CLASSES  │           │
│  │            │ │            │ │            │           │
│  │ Due: 2     │ │ Active: 1  │ │ New: 3     │           │
│  │ [View]     │ │ [Start]    │ │ [Watch]    │           │
│  └─────────────┘ └─────────────┘ └─────────────┘           │
│                                                              │
│  📊 MY PROGRESS                                             │
│  │                                                          │
│  │ ████████████████░░░░ 78%                                │
│  │ Half Yearly: 368/500 | Rank: 8                          │
│  │ [View Report Card]                                       │
│                                                              │
│  📢 ANNOUNCEMENTS                                           │
│  ├─ Half Yearly exam schedule released                     │
│  └─ Science fair registration open                         │
└─────────────────────────────────────────────────────────────┘
```

---

## 4. Feature Modules

### 4.1 My Classes (Timetable)

```
My Timetable | Monday

| Time        | Subject      | Teacher      | Room |
|-------------|--------------|--------------|------|
| 08:30-09:15 | Mathematics  | Mr. Kumar    | 102  |
| 09:15-10:00 | English      | Ms. Sharma   | 102  |
| 10:00-10:15 | Break        | -            | -    |
| 10:15-11:00 | Science      | Mr. Patel    | Lab  |
| 11:00-11:45 | Hindi        | Ms. Gupta    | 102  |
| 11:45-12:30 | Lunch        | -            | -    |
| 12:30-01:15 | Soc. Science | Mr. Singh    | 102  |
| 01:15-02:00 | Computer     | Ms. Rao      | Comp |
```

### 4.2 Homework & Assignments

```
My Assignments

📌 Pending (2)
┌─────────────────────────────────────────────────────────────┐
│ 📐 Mathematics | Chapter 5 Practice Problems                │
│ Due: 25-Jan-2026, 9:00 AM | Max Marks: 20                  │
│ Status: Not Started                                         │
│ [View Details] [Start Submission]                           │
└─────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────┐
│ 📖 English | Essay: My Role Model                          │
│ Due: 26-Jan-2026, 9:00 AM | Max Marks: 15                  │
│ Status: Draft Saved                                         │
│ [Continue Working]                                          │
└─────────────────────────────────────────────────────────────┘

✅ Completed (5)
│ Science Lab Report | 18/20 | Excellent work!
│ Hindi Grammar | 9/10 | Good
│ ...
```

### 4.3 Online Quizzes & Tests

```
Available Quizzes

🔴 Active Now
┌─────────────────────────────────────────────────────────────┐
│ Mathematics Unit Test 3                                      │
│ Duration: 45 min | Questions: 30 | Marks: 50                │
│ Ends: 25-Jan, 3:00 PM                                       │
│ [Start Quiz]                                                 │
└─────────────────────────────────────────────────────────────┘

📅 Upcoming
│ Science Chapter Review | 28-Jan | 30 min
│ English Grammar Test | 30-Jan | 20 min

📊 Completed
│ Math Practice Quiz | 38/50 (76%) | View Answers
│ Science Quiz 2 | 42/50 (84%) | View Answers
```

### 4.4 Digital Classroom

```
Recorded Classes

📚 Mathematics
├─ Quadratic Equations (45 min) - 22 Jan [Watched: 100%]
├─ Completing the Square (38 min) - 20 Jan [Watched: 60%]
└─ Factorization Review (42 min) - 18 Jan [New]

📚 Science
├─ Chemical Reactions (50 min) - 21 Jan [Watched: 80%]
└─ Periodic Table (45 min) - 19 Jan [Watched: 100%]

[Browse All Subjects]
```

### 4.5 My Results

```
My Results | 2025-26

| Examination   | Marks    | %    | Grade | Rank |
|---------------|----------|------|-------|------|
| Unit Test 1   | 425/500  | 85%  | A     | 5    |
| Unit Test 2   | 410/500  | 82%  | A     | 7    |
| Half Yearly   | 368/500  | 73%  | B+    | 8    |

Subject-wise (Half Yearly):
| Subject      | Marks | Grade | Class Avg |
|--------------|-------|-------|-----------|
| Mathematics  | 83/100| A     | 65        |
| Science      | 75/100| B+    | 68        |
| English      | 68/100| B     | 70        |
| Hindi        | 72/100| B+    | 72        |
| Soc. Science | 70/100| B+    | 71        |

[Download Report Card]
```

### 4.6 My Attendance

```
Attendance | January 2026

Attendance Rate: 94.7%

Calendar View:
     Mon  Tue  Wed  Thu  Fri  Sat
       6    7    8    9   10   11
       ✓    ✓    L    ✓    ✓    ✓

      13   14   15   16   17   18
       ✓    ✓    ✓    A    ✓    ✓

✓ Present: 18 | A Absent: 1 | L Late: 1

Note: You were absent on 16-Jan.
      Ensure 75% attendance for exam eligibility.
```

### 4.7 Library

```
Library

📖 My Borrowed Books (2)
| Title                    | Due Date  | Status    |
|--------------------------|-----------|-----------|
| Physics for Class 10     | 28-Jan    | 5 days    |
| Harry Potter (Chamber)   | 01-Feb    | 9 days    |

🔍 Search Library
[Search books, authors...]

📚 Recommended for You
├─ Mathematics Made Easy
├─ Science Encyclopedia
└─ English Grammar & Composition
```

### 4.8 Fees (View Only)

```
My Fees

Outstanding: ₹17,000
Due Date: 31-Jan-2026

Note: Please ask your parent/guardian to make the payment.

Recent Payments:
| Date      | Amount  | Receipt |
|-----------|---------|---------|
| 01-Oct-25 | ₹17,000 | View    |
| 01-Jul-25 | ₹17,000 | View    |
```

---

## 5. Learning Resources

### 5.1 Content Library

- Subject-wise study materials
- Video tutorials
- Practice worksheets
- Reference links
- Previous year papers

### 5.2 Notes & Bookmarks

```
My Notes

📁 Mathematics
├─ Quadratic Formula - Important! (22-Jan)
├─ Practice tips for Chapter 5 (20-Jan)
└─ Formulas to remember (18-Jan)

📁 Science
├─ Lab safety rules (21-Jan)
└─ Periodic table tricks (19-Jan)

[Create New Note]
```

---

## 6. Mobile App Features

### 6.1 Quick Actions

- View today's homework
- Check attendance
- Watch recorded class
- Take quiz
- View timetable

### 6.2 Offline Support

- Downloaded classes watchable offline
- Saved notes accessible offline
- Sync when online

---

## 7. API Endpoints

```
# Dashboard
GET    /api/v1/student/dashboard            # Dashboard data

# Academics
GET    /api/v1/student/timetable            # Timetable
GET    /api/v1/student/attendance           # Attendance
GET    /api/v1/student/results              # Results
GET    /api/v1/student/report-cards         # Report cards

# Assignments
GET    /api/v1/student/assignments          # List assignments
GET    /api/v1/student/assignments/{id}     # Assignment details
POST   /api/v1/student/assignments/{id}/submit # Submit

# Quizzes
GET    /api/v1/student/quizzes              # Available quizzes
GET    /api/v1/student/quizzes/{id}/start   # Start quiz
POST   /api/v1/student/quizzes/{id}/answer  # Submit answer
POST   /api/v1/student/quizzes/{id}/submit  # Submit quiz

# Digital Classroom
GET    /api/v1/student/recordings           # Class recordings
GET    /api/v1/student/recordings/{id}      # Watch recording
POST   /api/v1/student/recordings/{id}/progress # Update progress

# Library
GET    /api/v1/student/library/borrowed     # Borrowed books
GET    /api/v1/student/library/search       # Search catalog

# Notes
GET    /api/v1/student/notes                # My notes
POST   /api/v1/student/notes                # Create note
```

---

## 8. Related Documents

- [07-homework-assignments.md](./07-homework-assignments.md) - Assignments
- [08-online-quiz-assessment.md](./08-online-quiz-assessment.md) - Quizzes
- [09-digital-classroom.md](./09-digital-classroom.md) - Recordings
- [index.md](./index.md) - Main PRD index

---

**Previous**: [14-parent-portal.md](./14-parent-portal.md)
**Next**: [16-certificate-generation.md](./16-certificate-generation.md)
