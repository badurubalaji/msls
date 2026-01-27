# Story 3.8: Admission Reports & Analytics

**Epic:** 3 - School Setup & Admissions
**Status:** done
**Priority:** Medium
**Estimated Effort:** Medium

---

## User Story

As an **administrator**,
I want **to view admission reports and analytics**,
So that **I can track admission progress and plan capacity**.

---

## Acceptance Criteria

### AC1: Dashboard Overview
**Given** an admin is on admission reports
**When** viewing the dashboard
**Then** they see: total applications, by status, by class
**And** they see: conversion rate (enquiry to application to enrolled)
**And** they see: seats filled vs available by class

### AC2: Filtering and Export
**Given** filtering options are available
**When** admin applies filters
**Then** they can filter by: session, class, date range, status
**And** data updates accordingly
**And** reports can be exported to Excel/PDF

---

## Technical Requirements

### Backend (Go)

#### API Endpoints

| Method | Endpoint | Description | Permission |
|--------|----------|-------------|------------|
| GET | /api/v1/admissions/dashboard | Dashboard stats | admissions:read |
| GET | /api/v1/admissions/reports/funnel | Conversion funnel | admissions:read |
| GET | /api/v1/admissions/reports/class-wise | Class-wise report | admissions:read |
| GET | /api/v1/admissions/reports/source-analysis | Source analysis | admissions:read |
| GET | /api/v1/admissions/reports/daily-trend | Daily trend | admissions:read |
| GET | /api/v1/admissions/export | Export report | admissions:read |

#### Response Examples

```json
// Dashboard Stats
{
  "totalEnquiries": 450,
  "totalApplications": 320,
  "approved": 180,
  "enrolled": 150,
  "pending": 50,
  "rejected": 20,
  "conversionRates": {
    "enquiryToApplication": 71.1,
    "applicationToApproved": 56.2,
    "approvedToEnrolled": 83.3
  }
}

// Class-wise Report
{
  "classes": [
    {
      "className": "LKG",
      "totalSeats": 60,
      "applications": 85,
      "approved": 60,
      "enrolled": 55,
      "waitlisted": 10,
      "vacant": 5
    }
  ]
}
```

### Frontend (Angular)

#### Feature Structure

```
src/app/features/admissions/reports/
├── admission-dashboard.component.ts
├── funnel-chart.component.ts
├── class-wise-report.component.ts
└── source-analysis.component.ts
```

#### UI Components

1. **Dashboard Page** (`/admissions/dashboard`)
   - Summary cards (Total, Applied, Approved, Enrolled)
   - Conversion funnel visualization
   - Class-wise seat availability chart
   - Daily application trend

2. **Reports Page** (`/admissions/reports`)
   - Tabular reports with filters
   - Export buttons (Excel, PDF)

---

## Definition of Done

- [x] Backend: Report endpoints with aggregations
- [x] Backend: Export functionality (Excel export implemented via excelize)
- [x] Frontend: Dashboard with charts (Summary cards, funnel chart, class-wise report, source analysis, daily trend)
- [x] Frontend: Filter and export UI (Session filter, export button)
- [x] Navigation link added (Routes: /admissions/dashboard and /admissions/reports)
