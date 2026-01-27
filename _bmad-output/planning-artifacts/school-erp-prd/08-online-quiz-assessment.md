# 08 - Online Quiz & Assessment

**Parent Document**: [index.md](./index.md)
**Status**: Draft
**Phase**: 2

---

## 1. Overview

The Online Quiz & Assessment module enables teachers to create, conduct, and grade online assessments including quizzes, tests, and practice exams with automated grading and detailed analytics.

---

## 2. Question Bank

### 2.1 Question Entity

**Entity: Question**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| subject_id | UUID | Subject reference |
| topic_id | UUID | Topic reference (optional) |
| chapter_id | UUID | Chapter reference (optional) |
| question_type | ENUM | mcq, true_false, fill_blank, short_answer, long_answer, match, ordering |
| difficulty | ENUM | easy, medium, hard |
| question_text | TEXT | Question content (HTML/Markdown) |
| question_media | JSONB | Images, audio, video |
| options | JSONB | Answer options (for MCQ) |
| correct_answer | JSONB | Correct answer(s) |
| explanation | TEXT | Answer explanation |
| marks | DECIMAL | Default marks |
| negative_marks | DECIMAL | Negative marking |
| time_seconds | INT | Suggested time |
| tags | JSONB | Custom tags |
| is_verified | BOOLEAN | Reviewed by HOD |
| created_by | UUID | Creator |
| created_at | TIMESTAMP | Creation time |
| usage_count | INT | Times used |

### 2.2 Question Types

#### Multiple Choice (MCQ)
```json
{
  "question_type": "mcq",
  "question_text": "What is the capital of India?",
  "options": [
    { "id": "A", "text": "Mumbai", "is_correct": false },
    { "id": "B", "text": "New Delhi", "is_correct": true },
    { "id": "C", "text": "Kolkata", "is_correct": false },
    { "id": "D", "text": "Chennai", "is_correct": false }
  ],
  "marks": 1,
  "negative_marks": 0.25
}
```

#### True/False
```json
{
  "question_type": "true_false",
  "question_text": "The Earth is flat.",
  "correct_answer": false,
  "marks": 1
}
```

#### Fill in the Blank
```json
{
  "question_type": "fill_blank",
  "question_text": "The chemical formula of water is ____.",
  "correct_answer": ["H2O", "h2o"],
  "case_sensitive": false,
  "marks": 1
}
```

#### Match the Following
```json
{
  "question_type": "match",
  "question_text": "Match the countries with their capitals:",
  "left_column": [
    { "id": "1", "text": "India" },
    { "id": "2", "text": "Japan" },
    { "id": "3", "text": "France" }
  ],
  "right_column": [
    { "id": "A", "text": "Paris" },
    { "id": "B", "text": "New Delhi" },
    { "id": "C", "text": "Tokyo" }
  ],
  "correct_matches": { "1": "B", "2": "C", "3": "A" },
  "marks": 3
}
```

#### Short Answer
```json
{
  "question_type": "short_answer",
  "question_text": "Define photosynthesis in one sentence.",
  "expected_keywords": ["sunlight", "carbon dioxide", "glucose", "oxygen"],
  "max_words": 50,
  "marks": 2,
  "auto_grade": false
}
```

### 2.3 Topic/Chapter Organization

**Entity: QuestionTopic**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| subject_id | UUID | Subject reference |
| parent_id | UUID | Parent topic (for chapters) |
| name | VARCHAR(200) | Topic name |
| description | TEXT | Topic description |
| display_order | INT | Ordering |
| class_ids | JSONB | Applicable classes |

```
Subject: Science (Class 10)
├── Chapter 1: Chemical Reactions
│   ├── Topic: Types of Reactions
│   ├── Topic: Balancing Equations
│   └── Topic: Practical Applications
├── Chapter 2: Acids, Bases and Salts
│   ├── Topic: Properties of Acids
│   └── Topic: pH Scale
...
```

---

## 3. Quiz/Test Creation

### 3.1 Quiz Entity

**Entity: Quiz**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| title | VARCHAR(255) | Quiz title |
| description | TEXT | Instructions |
| quiz_type | ENUM | practice, graded, mock_exam, homework |
| subject_id | UUID | Subject |
| class_id | UUID | Target class |
| section_ids | JSONB | Target sections |
| total_marks | DECIMAL | Total marks |
| passing_marks | DECIMAL | Minimum to pass |
| duration_minutes | INT | Time limit |
| start_time | TIMESTAMP | When quiz opens |
| end_time | TIMESTAMP | When quiz closes |
| attempts_allowed | INT | Max attempts (0 = unlimited) |
| shuffle_questions | BOOLEAN | Randomize questions |
| shuffle_options | BOOLEAN | Randomize MCQ options |
| show_result | ENUM | immediately, after_deadline, manual |
| show_correct_answers | BOOLEAN | Show answers after |
| negative_marking | BOOLEAN | Enable negative marks |
| proctoring_enabled | BOOLEAN | Enable proctoring |
| status | ENUM | draft, scheduled, active, completed, archived |
| created_by | UUID | Creator |
| created_at | TIMESTAMP | Creation time |
| published_at | TIMESTAMP | When published |

### 3.2 Quiz Questions

**Entity: QuizQuestion**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| quiz_id | UUID | Quiz reference |
| question_id | UUID | Question from bank |
| display_order | INT | Question order |
| marks | DECIMAL | Marks for this quiz |
| negative_marks | DECIMAL | Negative marks |
| is_mandatory | BOOLEAN | Must attempt |

### 3.3 Quiz Creation Modes

**Manual Selection**:
- Pick questions from question bank
- Set individual marks

**Auto-Generate**:
```json
{
  "generation_rules": {
    "total_questions": 30,
    "total_marks": 50,
    "distribution": [
      { "difficulty": "easy", "count": 10, "marks_each": 1 },
      { "difficulty": "medium", "count": 15, "marks_each": 2 },
      { "difficulty": "hard", "count": 5, "marks_each": 4 }
    ],
    "topics": ["chapter_1", "chapter_2"],
    "question_types": ["mcq", "true_false"]
  }
}
```

---

## 4. Taking the Quiz

### 4.1 Quiz Attempt

**Entity: QuizAttempt**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| quiz_id | UUID | Quiz reference |
| student_id | UUID | Student reference |
| attempt_number | INT | Which attempt |
| started_at | TIMESTAMP | Start time |
| submitted_at | TIMESTAMP | Submission time |
| time_taken_seconds | INT | Actual time used |
| status | ENUM | in_progress, submitted, auto_submitted, abandoned |
| total_marks | DECIMAL | Marks obtained |
| percentage | DECIMAL | Percentage score |
| is_passed | BOOLEAN | Passed or not |
| ip_address | VARCHAR(45) | Client IP |
| device_info | JSONB | Device details |
| proctoring_flags | JSONB | Proctoring alerts |

### 4.2 Answer Recording

**Entity: QuizAnswer**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| attempt_id | UUID | Attempt reference |
| question_id | UUID | Question reference |
| answer | JSONB | Student's answer |
| is_correct | BOOLEAN | Auto-graded result |
| marks_obtained | DECIMAL | Marks given |
| graded_by | UUID | If manually graded |
| graded_at | TIMESTAMP | Grading time |
| feedback | TEXT | Teacher feedback |
| time_spent_seconds | INT | Time on question |

### 4.3 Answer Formats

**MCQ Answer**:
```json
{ "selected": "B" }
```

**Multiple Select Answer**:
```json
{ "selected": ["A", "C", "D"] }
```

**Fill Blank Answer**:
```json
{ "answer": "H2O" }
```

**Match Answer**:
```json
{ "matches": { "1": "B", "2": "C", "3": "A" } }
```

**Short/Long Answer**:
```json
{ "answer": "Photosynthesis is the process by which plants convert sunlight into glucose using carbon dioxide and water, releasing oxygen as a byproduct." }
```

---

## 5. Proctoring (Optional)

### 5.1 Proctoring Features

| Feature | Description |
|---------|-------------|
| Tab Switch Detection | Alert when student switches tabs |
| Copy-Paste Block | Prevent copy-paste |
| Right-Click Disable | Disable context menu |
| Fullscreen Mode | Force fullscreen |
| Webcam Monitoring | Periodic snapshots |
| Screen Recording | Record screen (consent required) |
| Multiple Face Detection | Alert if multiple faces |

### 5.2 Proctoring Log

**Entity: ProctoringEvent**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| attempt_id | UUID | Attempt reference |
| event_type | ENUM | tab_switch, copy_attempt, fullscreen_exit, face_not_detected, multiple_faces |
| event_time | TIMESTAMP | When occurred |
| details | JSONB | Event details |
| screenshot_url | VARCHAR(500) | Screenshot if captured |
| severity | ENUM | low, medium, high |

### 5.3 Violation Summary

```
Student: Aarav Sharma | Quiz: Math Unit Test

Violations Detected:
- Tab switched 3 times (Medium severity)
- Fullscreen exited 2 times (Low severity)
- Face not detected at 10:15 AM (High severity)

Total Flags: 6
Recommendation: Manual review required
```

---

## 6. Auto-Grading

### 6.1 Grading Rules

| Question Type | Auto-Grade | Logic |
|---------------|------------|-------|
| MCQ | Yes | Exact match |
| True/False | Yes | Exact match |
| Fill Blank | Yes | Case-insensitive match with alternatives |
| Match | Yes | All pairs must match |
| Ordering | Yes | Sequence match |
| Short Answer | Partial | Keyword matching (optional) |
| Long Answer | No | Manual grading required |

### 6.2 Partial Marking

```json
{
  "marking_scheme": {
    "mcq": {
      "correct": 1.0,
      "incorrect": -0.25,
      "unanswered": 0
    },
    "match": {
      "per_correct_pair": 0.5,
      "all_correct_bonus": 0.5
    }
  }
}
```

---

## 7. Results & Analytics

### 7.1 Individual Result

```
Quiz: Mathematics Unit Test 1
Student: Aarav Sharma | Class: 10-A

Score: 38/50 (76%)
Time Taken: 42 minutes
Rank: 5 out of 35

Section-wise Analysis:
| Section          | Marks | Max  | %    |
|------------------|-------|------|------|
| Algebra          | 12    | 15   | 80%  |
| Geometry         | 10    | 15   | 67%  |
| Trigonometry     | 16    | 20   | 80%  |

Difficulty Analysis:
- Easy: 10/10 (100%)
- Medium: 18/25 (72%)
- Hard: 10/15 (67%)

Areas for Improvement:
- Geometry: Circle theorems
- Trigonometry: Heights and distances
```

### 7.2 Class Analytics

```
Quiz: Mathematics Unit Test 1 | Class: 10-A

Participation: 34/35 (97%)
Average Score: 32/50 (64%)
Highest: 48/50 (Priya S.)
Lowest: 18/50 (Anonymous)

Score Distribution:
| Range     | Students | %    |
|-----------|----------|------|
| 90-100%   | 2        | 6%   |
| 80-89%    | 5        | 15%  |
| 70-79%    | 8        | 24%  |
| 60-69%    | 10       | 29%  |
| Below 60% | 9        | 26%  |

Difficult Questions:
- Q15: 20% correct (Circle theorem)
- Q22: 25% correct (Trigonometric identity)
- Q28: 30% correct (Quadratic word problem)
```

### 7.3 Question Analysis

```
Question Analysis Report

Q15: "In a circle, if angle at center is 120°, what is angle at circumference?"

Correct Answer: 60°
Students Correct: 7/35 (20%)

Option Distribution:
- A) 30° : 5 students (14%)
- B) 60° : 7 students (20%) ✓
- C) 120° : 18 students (51%)
- D) 240° : 5 students (14%)

Analysis: Most students chose 120° - confusion between angle at center and circumference theorem.

Recommendation: Revisit circle theorems in class.
```

---

## 8. Practice Mode

### 8.1 Self-Practice Quiz

**Features**:
- Unlimited attempts
- Instant feedback after each question
- Show correct answer immediately
- No time pressure (optional)
- Track improvement over attempts

### 8.2 Adaptive Practice

```json
{
  "adaptive_rules": {
    "initial_difficulty": "medium",
    "correct_streak_3": "increase_difficulty",
    "incorrect_streak_2": "decrease_difficulty",
    "session_length": 20
  }
}
```

---

## 9. API Endpoints

```
# Question Bank
GET    /api/v1/questions                    # List questions
POST   /api/v1/questions                    # Create question
GET    /api/v1/questions/{id}               # Get question
PUT    /api/v1/questions/{id}               # Update question
POST   /api/v1/questions/import             # Bulk import

# Quizzes
GET    /api/v1/quizzes                      # List quizzes
POST   /api/v1/quizzes                      # Create quiz
GET    /api/v1/quizzes/{id}                 # Get quiz
PUT    /api/v1/quizzes/{id}                 # Update quiz
POST   /api/v1/quizzes/{id}/publish         # Publish quiz
POST   /api/v1/quizzes/{id}/generate        # Auto-generate questions

# Taking Quiz
GET    /api/v1/quizzes/{id}/start           # Start attempt
POST   /api/v1/quizzes/{id}/answer          # Submit answer
POST   /api/v1/quizzes/{id}/submit          # Submit quiz
GET    /api/v1/quizzes/{id}/result          # Get result

# Analytics
GET    /api/v1/quizzes/{id}/analytics       # Quiz analytics
GET    /api/v1/quizzes/{id}/question-analysis # Question analysis
GET    /api/v1/students/{id}/quiz-history   # Student quiz history
```

---

## 10. Business Rules

| Rule | Description |
|------|-------------|
| Attempt Limit | Cannot exceed allowed attempts |
| Time Limit | Auto-submit when time expires |
| No Re-entry | Cannot re-enter after submission |
| Edit Window | Teacher can edit quiz only before start |
| Question Reuse | Track question exposure to prevent leaks |
| Minimum Questions | Quiz must have at least 5 questions |
| Result Lock | Results locked until all submissions |

---

## 11. Related Documents

- [06-examinations-grading.md](./06-examinations-grading.md) - Formal exams
- [07-homework-assignments.md](./07-homework-assignments.md) - Assignments
- [09-digital-classroom.md](./09-digital-classroom.md) - Learning content
- [index.md](./index.md) - Main PRD index

---

**Previous**: [07-homework-assignments.md](./07-homework-assignments.md)
**Next**: [09-digital-classroom.md](./09-digital-classroom.md)
