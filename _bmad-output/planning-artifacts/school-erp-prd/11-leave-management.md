# 11 - Leave Management

**Parent Document**: [index.md](./index.md)
**Status**: Draft
**Phase**: 1 (MVP)

---

## 1. Overview

The Leave Management module handles leave policies, applications, approvals, and balance tracking for all staff members.

---

## 2. Leave Types

### 2.1 Leave Type Configuration

**Entity: LeaveType**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| name | VARCHAR(100) | Leave type name |
| code | VARCHAR(10) | Short code |
| description | TEXT | Description |
| is_paid | BOOLEAN | Paid leave |
| is_encashable | BOOLEAN | Can be encashed |
| max_encashable | INT | Max days to encash |
| requires_document | BOOLEAN | Document required |
| document_after_days | INT | Document required after N days |
| is_active | BOOLEAN | Active status |

**Common Leave Types**:
```yaml
leave_types:
  - name: Casual Leave
    code: CL
    is_paid: true
    annual_quota: 12
    carry_forward: false
    max_consecutive: 3

  - name: Sick Leave
    code: SL
    is_paid: true
    annual_quota: 10
    carry_forward: false
    requires_document: true
    document_after_days: 2

  - name: Earned Leave
    code: EL
    is_paid: true
    annual_quota: 15
    carry_forward: true
    max_accumulation: 30
    is_encashable: true
    max_encashable: 10

  - name: Maternity Leave
    code: ML
    is_paid: true
    max_days: 180
    applicable_to: female

  - name: Paternity Leave
    code: PL
    is_paid: true
    max_days: 15
    applicable_to: male

  - name: Leave Without Pay
    code: LWP
    is_paid: false
    no_limit: true
```

---

## 3. Leave Policy

### 3.1 Policy Entity

**Entity: LeavePolicy**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| name | VARCHAR(100) | Policy name |
| applicable_to | JSONB | Staff types/designations |
| leave_rules | JSONB | Leave type rules |
| pro_rata_enabled | BOOLEAN | Pro-rata for new joiners |
| carry_forward_deadline | DATE | Annual deadline |
| is_default | BOOLEAN | Default policy |
| is_active | BOOLEAN | Active status |

### 3.2 Leave Rules per Type

```json
{
  "leave_type_id": "uuid",
  "annual_quota": 12,
  "monthly_accrual": 1,
  "carry_forward": false,
  "max_carry_forward": 0,
  "max_accumulation": 12,
  "min_days_per_request": 0.5,
  "max_days_per_request": 3,
  "advance_notice_days": 2,
  "can_combine_with": ["EL"],
  "cannot_combine_with": ["LWP"],
  "restricted_days": ["monday", "friday"],
  "sandwich_rule": true
}
```

### 3.3 Sandwich Rule

If leave is taken Friday and Monday, Saturday & Sunday count as leave.

---

## 4. Leave Balance

### 4.1 Balance Entity

**Entity: LeaveBalance**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| staff_id | UUID | Staff reference |
| leave_type_id | UUID | Leave type |
| academic_year_id | UUID | Year reference |
| opening_balance | DECIMAL | Start of year |
| accrued | DECIMAL | Monthly accruals |
| taken | DECIMAL | Leave taken |
| pending | DECIMAL | Pending approvals |
| lapsed | DECIMAL | Lapsed/expired |
| encashed | DECIMAL | Encashed |
| carry_forwarded | DECIMAL | Carried to next year |
| current_balance | DECIMAL | Available now |
| updated_at | TIMESTAMP | Last update |

### 4.2 Balance View

```
Leave Balance | Rajesh Kumar | 2025-26

| Leave Type    | Opening | Accrued | Taken | Pending | Balance |
|---------------|---------|---------|-------|---------|---------|
| Casual Leave  | 0       | 12      | 4     | 1       | 7       |
| Sick Leave    | 0       | 10      | 2     | 0       | 8       |
| Earned Leave  | 5       | 12      | 3     | 0       | 14      |
| ─────────────────────────────────────────────────────────────│
| Total         | 5       | 34      | 9     | 1       | 29      |

Note: Earned Leave includes 5 days carried forward from 2024-25
```

---

## 5. Leave Application

### 5.1 Application Entity

**Entity: LeaveApplication**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| staff_id | UUID | Applicant |
| leave_type_id | UUID | Leave type |
| from_date | DATE | Start date |
| to_date | DATE | End date |
| from_session | ENUM | first_half, second_half, full_day |
| to_session | ENUM | first_half, second_half, full_day |
| total_days | DECIMAL | Total leave days |
| reason | TEXT | Leave reason |
| contact_during_leave | VARCHAR(20) | Contact number |
| address_during_leave | TEXT | Address |
| work_handover_to | UUID | Handover staff |
| attachments | JSONB | Supporting documents |
| status | ENUM | draft, pending, approved, rejected, cancelled |
| applied_at | TIMESTAMP | Application time |
| approved_by | UUID | Approving authority |
| approved_at | TIMESTAMP | Approval time |
| rejection_reason | TEXT | If rejected |
| cancellation_reason | TEXT | If cancelled |

### 5.2 Application Interface

```
┌─────────────────────────────────────────────────────────────┐
│  APPLY FOR LEAVE                                             │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Leave Type: [Casual Leave      ▼]  Balance: 7 days         │
│                                                              │
│  From Date: [28-Jan-2026]  Session: [Full Day ▼]           │
│  To Date:   [29-Jan-2026]  Session: [Full Day ▼]           │
│                                                              │
│  Total Days: 2                                               │
│                                                              │
│  Reason:                                                     │
│  ┌─────────────────────────────────────────────────────┐    │
│  │ Family function - cousin's wedding                   │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                              │
│  Contact During Leave: [9876543210        ]                 │
│                                                              │
│  Work Handover To: [Ms. Priya Sharma      ▼]               │
│                                                              │
│  Attachments: [+ Add File] (Optional)                       │
│                                                              │
│  [Save Draft]  [Submit Application]                         │
└─────────────────────────────────────────────────────────────┘
```

---

## 6. Approval Workflow

### 6.1 Approval Hierarchy

```
Staff Member
     │
     ▼
Department Head / HOD
     │
     ▼
Principal (if > 3 days)
     │
     ▼
Admin (if special leave)
```

### 6.2 Approval Entity

**Entity: LeaveApproval**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| leave_application_id | UUID | Application reference |
| approver_id | UUID | Approver staff |
| level | INT | Approval level |
| status | ENUM | pending, approved, rejected |
| comments | TEXT | Approver comments |
| acted_at | TIMESTAMP | Action time |

### 6.3 Approval Interface

```
┌─────────────────────────────────────────────────────────────┐
│  PENDING LEAVE APPROVALS                                     │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  | Employee      | Type | Dates           | Days | Action | │
│  |---------------|------|-----------------|------|--------| │
│  | Rajesh Kumar  | CL   | 28-29 Jan       | 2    | [View] | │
│  | Priya Sharma  | SL   | 30 Jan          | 1    | [View] | │
│  | Amit Patel    | EL   | 01-05 Feb       | 5    | [View] | │
│                                                              │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  LEAVE REQUEST DETAILS                                       │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Employee: Rajesh Kumar (Senior Teacher)                    │
│  Leave Type: Casual Leave                                   │
│  Period: 28-Jan-2026 to 29-Jan-2026 (2 days)               │
│  Reason: Family function - cousin's wedding                 │
│  Work Handover: Ms. Priya Sharma                           │
│  Contact: 9876543210                                        │
│                                                              │
│  Leave Balance: 7 days available                            │
│  Previous this month: None                                  │
│                                                              │
│  Team Impact:                                                │
│  - 2 classes need substitute on 28-Jan                      │
│  - 2 classes need substitute on 29-Jan                      │
│                                                              │
│  Comments:                                                   │
│  ┌─────────────────────────────────────────────────────┐    │
│  │                                                      │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                              │
│  [Approve]  [Reject]  [Request More Info]                   │
└─────────────────────────────────────────────────────────────┘
```

---

## 7. Leave Calendar

### 7.1 Team Calendar View

```
Leave Calendar | January 2026 | Mathematics Dept

     Mon  Tue  Wed  Thu  Fri  Sat  Sun
      27   28   29   30   31    1    2
           ██       ██   ██

       3    4    5    6    7    8    9
                     ▓▓   ▓▓

      10   11   12   13   14   15   16
      ▓▓   ▓▓   ▓▓

Legend:
██ Rajesh Kumar (CL)
▓▓ Priya Sharma (EL)

Staff on leave today (28-Jan): 1
Available staff: 4
```

### 7.2 Holiday Calendar Integration

- Auto-exclude holidays from leave calculation
- Show restricted leave periods (exam time)
- Display school events

---

## 8. Notifications

### 8.1 Notification Triggers

| Event | Recipients | Channels |
|-------|------------|----------|
| Leave Applied | Approver | Push, Email |
| Leave Approved | Applicant | Push, Email, SMS |
| Leave Rejected | Applicant | Push, Email |
| Balance Low | Staff member | Push |
| Leave Tomorrow | Admin, HOD | Push |
| Document Pending | Staff member | Push |

---

## 9. Reports

### 9.1 Leave Summary Report

```
Leave Summary Report | January 2026 | All Staff

| Department  | Staff | CL Used | SL Used | EL Used | LWP | Total |
|-------------|-------|---------|---------|---------|-----|-------|
| Academic    | 25    | 18      | 8       | 12      | 2   | 40    |
| Admin       | 10    | 5       | 2       | 4       | 0   | 11    |
| Support     | 15    | 8       | 5       | 3       | 1   | 17    |
|-------------|-------|---------|---------|---------|-----|-------|
| Total       | 50    | 31      | 15      | 19      | 3   | 68    |

Absenteeism Rate: 2.8%
```

### 9.2 Individual Leave Report

```
Leave Report | Rajesh Kumar | 2025-26

Leave History:
| Date        | Type | Days | Status   | Reason              |
|-------------|------|------|----------|---------------------|
| 05-Apr-2025 | CL   | 1    | Approved | Personal work       |
| 22-Jun-2025 | SL   | 2    | Approved | Fever (Medical cert)|
| 15-Sep-2025 | CL   | 1    | Approved | Family emergency    |
| 28-Jan-2026 | CL   | 2    | Pending  | Family function     |

Summary:
- Total leaves taken: 6 days
- Leaves pending: 2 days
- Balance remaining: CL-7, SL-8, EL-14
```

---

## 10. API Endpoints

```
# Leave Types & Policies
GET    /api/v1/leave-types                  # List leave types
GET    /api/v1/leave-policies               # List policies
GET    /api/v1/leave-policies/{id}          # Get policy details

# Leave Balance
GET    /api/v1/staff/{id}/leave-balance     # Get balance
GET    /api/v1/leave-balance/team           # Team balance

# Leave Applications
GET    /api/v1/leave-applications           # List applications
POST   /api/v1/leave-applications           # Apply for leave
GET    /api/v1/leave-applications/{id}      # Get application
PUT    /api/v1/leave-applications/{id}      # Update application
DELETE /api/v1/leave-applications/{id}      # Cancel application

# Approvals
GET    /api/v1/leave-approvals/pending      # Pending approvals
POST   /api/v1/leave-applications/{id}/approve # Approve
POST   /api/v1/leave-applications/{id}/reject  # Reject

# Calendar
GET    /api/v1/leave-calendar               # Leave calendar
GET    /api/v1/leave-calendar/team          # Team calendar

# Reports
GET    /api/v1/leave/reports/summary        # Summary report
GET    /api/v1/leave/reports/individual     # Individual report
```

---

## 11. Business Rules

| Rule | Description |
|------|-------------|
| Balance Check | Cannot apply if insufficient balance |
| Advance Notice | CL needs 2 days, EL needs 7 days notice |
| Max Consecutive | CL max 3 consecutive days |
| Document Required | SL > 2 days needs medical certificate |
| Sandwich Rule | Weekends count if leave on Fri + Mon |
| Restricted Period | No CL during exams |
| Handover Required | Handover mandatory for > 3 days |
| Cancellation | Can cancel only before start date |
| Approval Timeout | Auto-escalate if no action in 2 days |

---

## 12. Related Documents

- [10-staff-management.md](./10-staff-management.md) - Staff profiles
- [04-academic-operations.md](./04-academic-operations.md) - Substitution
- [index.md](./index.md) - Main PRD index

---

**Previous**: [10-staff-management.md](./10-staff-management.md)
**Next**: [12-fees-payments.md](./12-fees-payments.md)
