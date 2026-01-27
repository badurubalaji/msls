# Epic 16: Online Quiz & Assessment

**Phase:** 2 (Extended)
**Priority:** Medium - Enhanced learning feature

## Epic Goal

Enable teachers to create and conduct online assessments with auto-grading and analytics.

## User Value

Teachers can build question banks, conduct timed quizzes, and students get instant feedback with detailed analytics.

## FRs Covered

FR-QZ-01 to FR-QZ-07

---

## Stories

### Story 16.1: Question Bank Management

As a **teacher**,
I want **to create and manage a question bank**,
So that **I can reuse questions across quizzes**.

**Acceptance Criteria:**

**Given** teacher is on question bank
**When** creating a question
**Then** they can select: question type (MCQ, true/false, fill-blank, match, short/long answer)
**And** they can enter: question text with rich formatting
**And** they can add: images, diagrams to question
**And** they can set: difficulty (easy, medium, hard)
**And** they can tag: topic, chapter

**Given** MCQ question
**When** adding options
**Then** they can add: 2-6 options
**And** they can mark: correct answer(s)
**And** they can add: option images
**And** they can set: partial marking rules

---

### Story 16.2: Quiz Creation

As a **teacher**,
I want **to create quizzes from question bank**,
So that **I can assess students**.

**Acceptance Criteria:**

**Given** teacher is creating quiz
**When** setting up
**Then** they can enter: title, description, instructions
**And** they can select: class, section, subject
**And** they can set: total marks, passing marks
**And** they can set: duration (minutes)
**And** they can set: start/end time window

**Given** adding questions
**When** selecting
**Then** they can: browse and select from bank
**And** they can: search/filter questions
**And** they can: adjust marks per question
**And** they can: reorder questions

---

### Story 16.3: Auto-Generate Quiz

As a **teacher**,
I want **to auto-generate quiz based on criteria**,
So that **quiz creation is faster**.

**Acceptance Criteria:**

**Given** teacher uses auto-generate
**When** setting criteria
**Then** they can set: total questions
**And** they can set: difficulty distribution (e.g., 5 easy, 10 medium, 5 hard)
**And** they can select: topics to include
**And** they can select: question types

**Given** generation runs
**When** complete
**Then** questions are selected based on criteria
**And** preview is shown for review
**And** teacher can swap individual questions
**And** finalize or regenerate

---

### Story 16.4: Quiz Attempt Interface

As a **student**,
I want **to attempt quizzes online**,
So that **I can complete assessments**.

**Acceptance Criteria:**

**Given** student opens a quiz
**When** starting
**Then** instructions are shown
**And** timer starts (visible throughout)
**And** questions displayed (one at a time or all)
**And** navigation panel shows question status

**Given** answering questions
**When** attempting
**Then** they can: select MCQ options
**And** they can: type fill-blank answers
**And** they can: match items (drag-drop)
**And** they can: flag questions for review
**And** answers auto-save periodically

**Given** quiz timer expires
**When** time runs out
**Then** quiz auto-submits
**And** warning at 5 min and 1 min remaining
**And** confirmation shown

---

### Story 16.5: Auto-Grading

As a **system**,
I want **to auto-grade objective questions**,
So that **instant results are available**.

**Acceptance Criteria:**

**Given** quiz is submitted
**When** processing
**Then** MCQ answers checked against correct
**And** fill-blank checked (case-insensitive, alternatives)
**And** match-type checked for pairs
**And** marks calculated per question

**Given** partial marking configured
**When** calculating
**Then** partial marks for partial matches
**And** negative marking applied if enabled
**And** total score calculated

**Given** subjective questions exist
**When** auto-grading
**Then** they are marked as "pending review"
**And** teacher notification to grade

---

### Story 16.6: Quiz Results & Analytics

As a **teacher**,
I want **to view quiz results and analytics**,
So that **I can understand student performance**.

**Acceptance Criteria:**

**Given** quiz has submissions
**When** viewing results
**Then** they see: list of students with scores
**And** they see: pass/fail status
**And** they see: average, highest, lowest scores
**And** export to Excel available

**Given** question analysis
**When** viewing
**Then** they see: each question correctness %
**And** they see: most missed questions
**And** they see: option distribution for MCQ
**And** insights for teaching improvement

---

### Story 16.7: Basic Proctoring

As an **administrator**,
I want **basic proctoring features**,
So that **quiz integrity is maintained**.

**Acceptance Criteria:**

**Given** proctoring is enabled for quiz
**When** student takes quiz
**Then** fullscreen mode is enforced
**And** tab switch is detected and logged
**And** copy-paste is disabled
**And** right-click is disabled

**Given** violations occur
**When** detected
**Then** warning shown to student
**And** violation logged with timestamp
**And** teacher can see violation report
**And** configurable auto-submit on too many violations

---

### Story 16.8: Practice Mode

As a **student**,
I want **to practice with unlimited attempts**,
So that **I can prepare for exams**.

**Acceptance Criteria:**

**Given** practice quiz is available
**When** attempting
**Then** unlimited attempts allowed
**And** instant feedback per question
**And** correct answer shown immediately
**And** no time pressure (optional timer)

**Given** practice history
**When** viewing
**Then** they see: attempt history with scores
**And** they see: improvement trend
**And** weak areas highlighted
