# Epic 22: Analytics & Dashboards

**Phase:** 2 (Extended)
**Priority:** Medium - Reporting feature

## Epic Goal

Provide comprehensive analytics and customizable reporting across all modules.

## User Value

Principals and admins get role-specific dashboards, can build custom reports, and schedule automated report delivery.

## FRs Covered

FR-AN-01 to FR-AN-06

---

## Stories

### Story 22.1: Principal Dashboard

As a **principal**,
I want **a comprehensive dashboard**,
So that **I see school overview at a glance**.

**Acceptance Criteria:**

**Given** principal logs in
**When** viewing dashboard
**Then** they see: student count, staff count
**And** they see: today's attendance %
**And** they see: fee collection status
**And** they see: upcoming events

**Given** dashboard widgets
**When** displayed
**Then** they see: attendance trend (30 days)
**And** they see: academic performance summary
**And** they see: alerts (low attendance, defaulters)
**And** widgets are clickable for details

---

### Story 22.2: Teacher Dashboard

As a **teacher**,
I want **a dashboard for my classes**,
So that **I see my responsibilities**.

**Acceptance Criteria:**

**Given** teacher logs in
**When** viewing dashboard
**Then** they see: today's classes
**And** they see: pending attendance marking
**And** they see: assignments to grade
**And** they see: students at risk

**Given** class summary
**When** displayed
**Then** they see: per class stats
**And** they see: attendance %, average marks
**And** they see: quick actions (mark attendance, create assignment)

---

### Story 22.3: Academic Analytics

As an **academic coordinator**,
I want **academic analytics**,
So that **performance is analyzed**.

**Acceptance Criteria:**

**Given** coordinator is on analytics
**When** viewing academic data
**Then** they see: pass % by class, subject
**And** they see: grade distribution
**And** they see: improvement/decline trends
**And** they see: subject-wise comparison

**Given** drill-down needed
**When** clicking
**Then** they can see: class-level details
**And** they can see: individual student data
**And** they can export: detailed reports

---

### Story 22.4: Financial Analytics

As an **accounts head**,
I want **financial analytics**,
So that **revenue is tracked**.

**Acceptance Criteria:**

**Given** accounts views analytics
**When** displayed
**Then** they see: total collection, target
**And** they see: collection by mode
**And** they see: outstanding dues aging
**And** they see: monthly trend

**Given** detailed view
**When** drilling down
**Then** they see: class-wise collection
**And** they see: defaulter analysis
**And** they see: fee category breakdown

---

### Story 22.5: Attendance Analytics

As an **administrator**,
I want **attendance analytics**,
So that **patterns are identified**.

**Acceptance Criteria:**

**Given** admin views attendance analytics
**When** displayed
**Then** they see: school-wide attendance %
**And** they see: class-wise comparison
**And** they see: day-of-week patterns
**And** they see: chronic absentees count

**Given** trends view
**When** analyzing
**Then** they see: monthly trends
**And** they see: comparison with previous year
**And** they see: impact of holidays

---

### Story 22.6: Custom Report Builder

As an **administrator**,
I want **to build custom reports**,
So that **specific data needs are met**.

**Acceptance Criteria:**

**Given** admin is on report builder
**When** creating report
**Then** they can select: data source (module)
**And** they can select: columns to include
**And** they can set: filters
**And** they can set: grouping, sorting

**Given** report is built
**When** running
**Then** data is displayed
**And** visualizations auto-generated
**And** report can be saved
**And** export available (Excel, PDF)

---

### Story 22.7: Report Scheduling

As an **administrator**,
I want **to schedule report delivery**,
So that **reports are received automatically**.

**Acceptance Criteria:**

**Given** report exists
**When** scheduling
**Then** they can set: frequency (daily, weekly, monthly)
**And** they can set: delivery time
**And** they can set: recipients (email)
**And** they can set: format (PDF, Excel)

**Given** schedule is active
**When** triggered
**Then** report is generated
**And** emailed to recipients
**And** run history logged

---

### Story 22.8: Data Visualization Widgets

As an **administrator**,
I want **rich data visualizations**,
So that **data is easily understood**.

**Acceptance Criteria:**

**Given** analytics are displayed
**When** viewing charts
**Then** they see: line charts for trends
**And** they see: bar charts for comparisons
**And** they see: pie/donut for distribution
**And** they see: heatmaps for calendar data

**Given** interactive charts
**When** interacting
**Then** hover shows: detailed values
**And** click drills down: to details
**And** charts are responsive
