# 21 - Analytics & Dashboards

**Parent Document**: [index.md](./index.md)
**Status**: Draft
**Phase**: 2

---

## 1. Overview

The Analytics & Dashboards module provides role-based dashboards, data visualization, custom reports, and scheduled report generation.

---

## 2. Role-Based Dashboards

### 2.1 Principal Dashboard

```
┌─────────────────────────────────────────────────────────────┐
│  PRINCIPAL DASHBOARD | ABC Public School                     │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  TODAY'S SNAPSHOT                                           │
│  ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐               │
│  │Students│ │ Staff  │ │Attend %│ │Fee Coll│               │
│  │ 1,250  │ │   85   │ │ 94.2%  │ │₹4.5L   │               │
│  └────────┘ └────────┘ └────────┘ └────────┘               │
│                                                              │
│  ATTENDANCE TREND (30 Days)                                 │
│  100│    ╭───────╮                                         │
│   95│───╯        ╰──╮                                      │
│   90│               ╰───                                   │
│   85│                                                      │
│     └─────────────────────────────                         │
│      Week 1   Week 2   Week 3   Week 4                     │
│                                                              │
│  FEE COLLECTION vs TARGET                                   │
│  ████████████████████░░░░░░ 78% (₹78L / ₹100L)            │
│                                                              │
│  ACADEMIC PERFORMANCE (Half Yearly)                         │
│  | Class | Avg % | Pass % | Top Performer    |             │
│  |-------|-------|--------|------------------|             │
│  | 10    | 72%   | 95%    | Priya S. (95%)   |             │
│  | 9     | 68%   | 92%    | Rahul K. (93%)   |             │
│  | 8     | 74%   | 97%    | Ananya P. (91%)  |             │
│                                                              │
│  ALERTS                                                     │
│  ⚠️ 15 students below 75% attendance                        │
│  ⚠️ 8 staff leave requests pending                         │
│  ⚠️ Fee dues > 60 days: 25 students                        │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 Admin Dashboard

```
┌─────────────────────────────────────────────────────────────┐
│  ADMIN DASHBOARD                                             │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ADMISSIONS                    FEES                          │
│  Applications: 320            Collected: ₹78L               │
│  Approved: 180                Outstanding: ₹22L             │
│  Enrolled: 150                Defaulters: 45                │
│  Waitlist: 25                 This Month: ₹12L              │
│                                                              │
│  STAFF OVERVIEW               TRANSPORT                     │
│  Present: 82/85               Vehicles: 12                  │
│  On Leave: 3                  Routes: 8                     │
│  Pending Leaves: 8            Students: 450                 │
│                               Active Trips: 8               │
│                                                              │
│  PENDING ACTIONS                                            │
│  ├─ 8 Leave approvals                                      │
│  ├─ 5 Certificate requests                                 │
│  ├─ 3 Fee concession requests                              │
│  └─ 12 Admission applications                              │
└─────────────────────────────────────────────────────────────┘
```

### 2.3 Teacher Dashboard

```
┌─────────────────────────────────────────────────────────────┐
│  TEACHER DASHBOARD | Mr. Rajesh Kumar                       │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  TODAY'S CLASSES                                            │
│  ├─ 08:30 - 10-A Mathematics [Attendance Done]             │
│  ├─ 09:15 - 10-B Mathematics [Attendance Pending]          │
│  ├─ 10:15 - 9-A Mathematics [Upcoming]                     │
│  └─ 11:00 - 9-B Mathematics [Upcoming]                     │
│                                                              │
│  MY CLASSES SUMMARY                                         │
│  | Class | Students | Avg Attend | Avg Marks | At Risk |   │
│  |-------|----------|------------|-----------|---------|   │
│  | 10-A  | 42       | 94%        | 72%       | 3       |   │
│  | 10-B  | 40       | 91%        | 68%       | 5       |   │
│  | 9-A   | 45       | 93%        | 75%       | 2       |   │
│  | 9-B   | 43       | 90%        | 70%       | 4       |   │
│                                                              │
│  PENDING TASKS                                              │
│  ├─ 38 assignments to grade                                │
│  ├─ 5 quiz results to publish                              │
│  └─ 2 class recordings to upload                           │
└─────────────────────────────────────────────────────────────┘
```

---

## 3. Analytics Modules

### 3.1 Academic Analytics

```
Academic Analytics | 2025-26

PERFORMANCE DISTRIBUTION
| Grade | Students | % |
|-------|----------|---|
| A+    | 120      | 10% | ████████░░
| A     | 240      | 20% | ████████████████░░░░
| B+    | 300      | 25% | ████████████████████░░░░░
| B     | 240      | 20% | ████████████████░░░░
| C     | 180      | 15% | ████████████░░░░
| D     | 120      | 10% | ████████░░

SUBJECT-WISE PASS RATE
| Subject      | Pass % | Trend |
|--------------|--------|-------|
| Mathematics  | 88%    | ↑ +3% |
| Science      | 92%    | ↔     |
| English      | 95%    | ↑ +2% |
| Hindi        | 94%    | ↓ -1% |

AT-RISK STUDENTS: 45 (Below 40% or <75% attendance)
```

### 3.2 Financial Analytics

```
Financial Analytics | January 2026

COLLECTION TREND
₹15L│    ╭─╮
₹12L│  ╭─╯ ╰─╮
₹9L │──╯     ╰──╮
₹6L │           ╰──
    └──────────────────
    Apr May Jun Jul Aug Sep Oct Nov Dec Jan

COLLECTION BY MODE
| Mode       | Amount  | % |
|------------|---------|---|
| UPI        | ₹35L    | 45% | █████████████████████░░░░
| Net Banking| ₹25L    | 32% | ████████████████░░░░
| Card       | ₹10L    | 13% | ██████░░░░
| Cash       | ₹5L     | 6%  | ███░░
| Cheque     | ₹3L     | 4%  | ██░░

OUTSTANDING AGING
| Age        | Amount | Students |
|------------|--------|----------|
| 0-30 days  | ₹8L    | 45       |
| 31-60 days | ₹6L    | 30       |
| 61-90 days | ₹4L    | 20       |
| 90+ days   | ₹4L    | 15       |
```

### 3.3 Attendance Analytics

```
Attendance Analytics | January 2026

DAILY ATTENDANCE TREND
100%│────────╮  ╭──────
 95%│        ╰──╯
 90%│
    └─────────────────────
    1  5  10  15  20  25  30

CLASS-WISE COMPARISON
| Class | Attendance % | Trend |
|-------|--------------|-------|
| 10    | 94.5%        | ↑     |
| 9     | 92.3%        | ↔     |
| 8     | 95.1%        | ↑     |
| 7     | 91.8%        | ↓     |

CHRONIC ABSENTEES (<75%)
- 15 students identified
- Top reasons: Medical (8), Unknown (5), Transport (2)
```

---

## 4. Custom Reports

### 4.1 Report Builder

**Entity: ReportTemplate**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant |
| name | VARCHAR(200) | Report name |
| description | TEXT | Description |
| module | VARCHAR(50) | Source module |
| columns | JSONB | Selected columns |
| filters | JSONB | Filter criteria |
| grouping | JSONB | Group by fields |
| sorting | JSONB | Sort order |
| charts | JSONB | Chart configs |
| created_by | UUID | Creator |
| is_public | BOOLEAN | Shared report |

### 4.2 Report Configuration

```json
{
  "name": "Class Performance Report",
  "module": "academics",
  "columns": [
    "student_name",
    "class",
    "section",
    "attendance_percent",
    "exam_average",
    "grade"
  ],
  "filters": [
    { "field": "class", "operator": "in", "value": ["10", "9"] },
    { "field": "exam_average", "operator": ">=", "value": 40 }
  ],
  "grouping": ["class", "section"],
  "sorting": [
    { "field": "exam_average", "order": "desc" }
  ],
  "charts": [
    { "type": "bar", "x": "section", "y": "avg(exam_average)" }
  ]
}
```

---

## 5. Scheduled Reports

### 5.1 Schedule Entity

**Entity: ReportSchedule**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| report_template_id | UUID | Template |
| name | VARCHAR(200) | Schedule name |
| frequency | ENUM | daily, weekly, monthly |
| day_of_week | INT | If weekly |
| day_of_month | INT | If monthly |
| time | TIME | Generation time |
| format | ENUM | pdf, excel, csv |
| recipients | JSONB | Email recipients |
| is_active | BOOLEAN | Active status |
| last_run | TIMESTAMP | Last execution |
| next_run | TIMESTAMP | Next execution |

### 5.2 Scheduled Reports List

```
Scheduled Reports

| Report               | Frequency | Next Run | Recipients   |
|----------------------|-----------|----------|--------------|
| Daily Attendance     | Daily     | Tomorrow | Principal    |
| Weekly Fee Summary   | Weekly    | Monday   | Admin, Acct  |
| Monthly Performance  | Monthly   | 1st Feb  | Principal    |
| Defaulter List       | Weekly    | Friday   | Admin        |
```

---

## 6. Export Options

### 6.1 Export Formats

| Format | Use Case |
|--------|----------|
| PDF | Formal reports, printing |
| Excel | Data analysis, manipulation |
| CSV | Data import/export |
| HTML | Email embedding |

### 6.2 Export Features

- Custom headers/footers
- School logo/branding
- Page numbering
- Date/time stamps
- Digital signature support

---

## 7. Data Visualization

### 7.1 Chart Types

| Type | Use Case |
|------|----------|
| Line | Trends over time |
| Bar | Comparisons |
| Pie | Distribution |
| Donut | Percentages |
| Heatmap | Attendance calendar |
| Gauge | KPI targets |
| Table | Detailed data |

### 7.2 Dashboard Widgets

- KPI Cards (single metrics)
- Trend Charts (time series)
- Comparison Charts (bar/column)
- Distribution Charts (pie/donut)
- Tables (paginated data)
- Alerts (threshold-based)
- Calendar (attendance heatmap)

---

## 8. API Endpoints

```
# Dashboards
GET    /api/v1/dashboards/{role}            # Role dashboard
GET    /api/v1/dashboards/widgets/{type}    # Widget data

# Analytics
GET    /api/v1/analytics/academic           # Academic analytics
GET    /api/v1/analytics/financial          # Financial analytics
GET    /api/v1/analytics/attendance         # Attendance analytics
GET    /api/v1/analytics/custom             # Custom query

# Reports
GET    /api/v1/reports/templates            # List templates
POST   /api/v1/reports/templates            # Create template
POST   /api/v1/reports/generate             # Generate report
GET    /api/v1/reports/download/{id}        # Download report

# Schedules
GET    /api/v1/reports/schedules            # List schedules
POST   /api/v1/reports/schedules            # Create schedule
```

---

## 9. Related Documents

- All module documents feed into analytics
- [index.md](./index.md) - Main PRD index

---

**Previous**: [20-visitor-gate-management.md](./20-visitor-gate-management.md)
**Next**: [22-ai-capabilities.md](./22-ai-capabilities.md)
