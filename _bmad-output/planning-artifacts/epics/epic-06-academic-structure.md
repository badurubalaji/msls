# Epic 6: Academic Structure & Timetable

**Phase:** 1 (MVP)
**Priority:** High - Foundation for academic operations

## Epic Goal

Enable schools to define their academic structure and create comprehensive timetables.

## User Value

Admins can configure classes, sections, subjects, and teachers can view their schedules with substitution support.

## FRs Covered

FR-AO-01 to FR-AO-05

---

## Stories

### Story 6.1: Class & Section Management

As an **academic administrator**,
I want **to configure classes and sections**,
So that **student groupings are properly defined**.

**Acceptance Criteria:**

**Given** admin is on class configuration
**When** creating a class
**Then** they can enter: class name (e.g., "Class 10"), numeric order
**And** they can set: applicable academic years
**And** they can set: streams (if senior classes - Science, Commerce, Arts)

**Given** a class exists
**When** adding sections
**Then** they can enter: section name (A, B, C)
**And** they can set: max capacity
**And** they can assign: class teacher
**And** sections are linked to the class

**Given** class-section setup is complete
**When** viewing the structure
**Then** they see hierarchical view: Class → Sections
**And** each section shows: student count, teacher assigned
**And** capacity utilization is displayed

---

### Story 6.2: Subject Configuration

As an **academic administrator**,
I want **to configure subjects with class mappings**,
So that **curriculum structure is defined**.

**Acceptance Criteria:**

**Given** admin is on subject configuration
**When** creating a subject
**Then** they can enter: subject name, code, type (mandatory/optional)
**And** they can set: periods per week (default)
**And** they can set: is practical subject (yes/no)

**Given** subjects are configured
**When** mapping to classes
**Then** they can select: which subjects apply to which class
**And** they can set: periods per week for that class
**And** they can set: passing marks, maximum marks

**Given** subject groups exist (for streams)
**When** configuring streams
**Then** subjects can be grouped (Science = Physics, Chemistry, Math, Biology)
**And** students can be assigned to subject groups
**And** optional subjects can be selected individually

---

### Story 6.3: Period Slot Configuration

As an **academic administrator**,
I want **to define period slots and school timing**,
So that **timetable can be created with correct time slots**.

**Acceptance Criteria:**

**Given** admin is on timetable settings
**When** defining period slots
**Then** they can enter: period number, start time, end time
**And** they can mark periods as: regular, short, assembly, break, lunch
**And** they can set: duration (minutes)

**Given** different day patterns exist
**When** configuring day types
**Then** they can create patterns: regular day, Saturday (half-day)
**And** each day type has its own period structure
**And** days of week can be assigned patterns

**Given** school operates different shifts
**When** configuring shifts
**Then** they can define: morning shift, afternoon shift
**And** each shift has its own period slots
**And** sections can be assigned to shifts

---

### Story 6.4: Timetable Creation

As an **academic administrator**,
I want **to create timetables for each section**,
So that **teaching schedule is organized**.

**Acceptance Criteria:**

**Given** admin is creating a timetable
**When** using the timetable builder
**Then** they see a grid: days (columns) × periods (rows)
**And** they can drag-drop subjects to slots
**And** they can assign teacher to each slot
**And** colors differentiate subjects

**Given** a subject is being assigned
**When** selecting the slot
**Then** available teachers for that subject are shown
**And** teacher's existing assignments are visible
**And** conflict warning shows if teacher is already assigned

**Given** timetable is complete
**When** publishing the timetable
**Then** timetable becomes active
**And** previous timetable is archived
**And** teachers and students can view new timetable

---

### Story 6.5: Teacher Timetable View

As a **teacher**,
I want **to view my teaching timetable**,
So that **I know my schedule for each day**.

**Acceptance Criteria:**

**Given** a teacher is logged in
**When** viewing their timetable
**Then** they see weekly view with all assigned periods
**And** each slot shows: class, section, subject, room
**And** free periods are marked

**Given** daily view is selected
**When** viewing today's schedule
**Then** current/next period is highlighted
**And** countdown to next class is shown
**And** they can navigate to other days

**Given** timetable needs to be printed
**When** exporting
**Then** printable format is generated
**And** options include: week view, list view
**And** PDF download is available

---

### Story 6.6: Section Timetable View

As a **student or parent**,
I want **to view the class timetable**,
So that **I know the daily schedule**.

**Acceptance Criteria:**

**Given** a student/parent is logged in
**When** viewing the timetable
**Then** they see their section's weekly timetable
**And** each slot shows: subject, teacher name
**And** special periods (assembly, breaks) are shown

**Given** mobile view
**When** viewing on phone
**Then** daily view is shown by default
**And** swipe navigates between days
**And** compact display fits mobile screen

---

### Story 6.7: Substitution Management

As an **academic administrator**,
I want **to manage teacher substitutions**,
So that **classes continue when teachers are absent**.

**Acceptance Criteria:**

**Given** a teacher is marked absent
**When** creating substitution
**Then** system shows: affected periods for the day
**And** system suggests: available teachers (based on free periods)
**And** admin can assign substitute teacher

**Given** substitute is assigned
**When** substitution is saved
**Then** substitute teacher sees the class in their timetable
**And** class timetable shows substitute teacher name
**And** SMS/notification sent to substitute

**Given** substitution history
**When** viewing reports
**Then** they see: date, original teacher, substitute, class
**And** teacher-wise substitution count is available
**And** workload balancing can be analyzed

---

### Story 6.8: Room & Resource Allocation

As an **academic administrator**,
I want **to manage rooms and their allocation**,
So that **classes have proper venues**.

**Acceptance Criteria:**

**Given** admin is on room management
**When** creating a room
**Then** they can enter: room name/number, type (classroom, lab, library)
**And** they can set: capacity, floor, building
**And** they can set: available resources (projector, AC, etc.)

**Given** rooms exist
**When** assigning to timetable
**Then** room can be linked to period slot
**And** room availability is shown
**And** conflict warning if room double-booked

**Given** special resources are needed
**When** scheduling lab periods
**Then** lab rooms are assigned to practical periods
**And** equipment availability is checked
**And** lab assistant can be notified
