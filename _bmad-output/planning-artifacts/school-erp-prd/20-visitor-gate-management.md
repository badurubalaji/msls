# 20 - Visitor & Gate Management

**Parent Document**: [index.md](./index.md)
**Status**: Draft
**Phase**: 2

---

## 1. Overview

The Visitor & Gate Management module handles visitor registration, gate passes for students, vehicle entry logs, and security management.

---

## 2. Visitor Management

### 2.1 Visitor Entry

**Entity: VisitorEntry**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| branch_id | UUID | Branch |
| visitor_name | VARCHAR(200) | Visitor name |
| phone | VARCHAR(20) | Phone number |
| email | VARCHAR(255) | Email |
| id_type | ENUM | aadhar, pan, driving_license, passport |
| id_number | VARCHAR(50) | ID number |
| photo_url | VARCHAR(500) | Visitor photo |
| company | VARCHAR(200) | Organization |
| purpose | ENUM | meeting, delivery, parent, vendor, interview, other |
| purpose_details | TEXT | Visit details |
| whom_to_meet | UUID | Staff to meet |
| department | VARCHAR(100) | Department |
| entry_time | TIMESTAMP | Check-in time |
| expected_exit | TIME | Expected exit |
| exit_time | TIMESTAMP | Actual exit |
| badge_number | VARCHAR(20) | Visitor badge |
| vehicle_number | VARCHAR(20) | Vehicle if any |
| status | ENUM | checked_in, checked_out, cancelled |
| approved_by | UUID | Approval authority |
| remarks | TEXT | Security notes |

### 2.2 Visitor Registration Interface

```
┌─────────────────────────────────────────────────────────────┐
│  VISITOR CHECK-IN                                            │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  📷 [Capture Photo]                                         │
│                                                              │
│  Name: [Rajesh Verma                    ]                   │
│  Phone: [9876543210                     ]                   │
│  Company: [ABC Technologies             ]                   │
│                                                              │
│  ID Proof: [Aadhar ▼]  Number: [XXXX-XXXX-1234]            │
│                                                              │
│  Purpose: [Vendor Meeting ▼]                                │
│  Details: [Printer maintenance discussion]                  │
│                                                              │
│  To Meet: [Mr. Sharma - IT Department ▼]                   │
│                                                              │
│  Vehicle: [MH01-AB-1234] (Optional)                        │
│                                                              │
│  [Check In]                                                  │
└─────────────────────────────────────────────────────────────┘

Badge Issued: V-2026-0142
Entry Time: 10:30 AM
```

### 2.3 Visitor Badge

```
┌─────────────────────────┐
│    VISITOR PASS         │
│                         │
│  [PHOTO]  V-2026-0142  │
│                         │
│  RAJESH VERMA          │
│  ABC Technologies       │
│                         │
│  To Meet: Mr. Sharma   │
│  Purpose: Vendor       │
│                         │
│  Date: 22-Jan-2026     │
│  Entry: 10:30 AM       │
│                         │
│  [QR CODE]             │
│                         │
│  Please return badge   │
│  at exit               │
└─────────────────────────┘
```

---

## 3. Student Gate Pass

### 3.1 Gate Pass Entity

**Entity: GatePass**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| student_id | UUID | Student reference |
| pass_number | VARCHAR(50) | Pass number |
| pass_type | ENUM | early_leave, mid_day, full_day |
| request_reason | TEXT | Reason for leaving |
| request_date | DATE | Pass date |
| exit_time | TIME | Expected exit |
| return_time | TIME | Expected return |
| actual_exit | TIMESTAMP | Actual exit |
| actual_return | TIMESTAMP | Actual return |
| pickup_person | VARCHAR(200) | Who is picking |
| pickup_relation | VARCHAR(50) | Relation to student |
| pickup_phone | VARCHAR(20) | Contact |
| pickup_id_type | VARCHAR(50) | ID shown |
| requested_by | ENUM | parent, staff |
| approved_by | UUID | Approving authority |
| status | ENUM | requested, approved, rejected, used, cancelled |
| parent_notified | BOOLEAN | Parent informed |

### 3.2 Gate Pass Request Flow

```
Request Flow:

Parent Request (App):
Parent → Request via App → Class Teacher Approval →
Admin Notification → Gate Pass Generated →
Parent Shows at Gate → Student Released

Staff Request:
Staff Request → Class Teacher → Admin Approval →
Parent Notified → Gate Pass Ready → Student Released
```

### 3.3 Gate Pass Interface

```
┌─────────────────────────────────────────────────────────────┐
│  GATE PASS REQUEST                                           │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Student: Aarav Sharma | Class: 10-A | Roll: 15            │
│                                                              │
│  Pass Type: ○ Early Leave ○ Mid-day ○ Full Day             │
│                                                              │
│  Reason: [Doctor's appointment                       ]      │
│                                                              │
│  Exit Time: [11:30 AM]                                      │
│  Return Time: [02:00 PM] (if mid-day)                       │
│                                                              │
│  Pickup By: [Rajesh Sharma (Father)       ]                │
│  Contact: [9876543210]                                      │
│                                                              │
│  [Submit Request]                                            │
│                                                              │
│  Note: Parent will receive SMS confirmation                 │
└─────────────────────────────────────────────────────────────┘
```

---

## 4. Vehicle Entry

### 4.1 Vehicle Log

**Entity: VehicleEntry**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| vehicle_number | VARCHAR(20) | Registration |
| vehicle_type | ENUM | car, bike, auto, truck, other |
| driver_name | VARCHAR(200) | Driver name |
| driver_phone | VARCHAR(20) | Driver contact |
| purpose | ENUM | delivery, visitor, staff, parent |
| related_entry_id | UUID | Linked visitor/delivery |
| entry_time | TIMESTAMP | Entry time |
| exit_time | TIMESTAMP | Exit time |
| entry_gate | VARCHAR(50) | Entry gate |
| exit_gate | VARCHAR(50) | Exit gate |
| parking_slot | VARCHAR(20) | Allocated parking |

---

## 5. Delivery Management

### 5.1 Delivery Entry

**Entity: DeliveryEntry**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| delivery_person | VARCHAR(200) | Delivery person name |
| company | VARCHAR(200) | Delivery company |
| phone | VARCHAR(20) | Contact |
| package_type | ENUM | courier, food, supplies, other |
| package_description | TEXT | Description |
| recipient | UUID | Receiving staff |
| department | VARCHAR(100) | Department |
| entry_time | TIMESTAMP | Entry time |
| exit_time | TIMESTAMP | Exit time |
| received_by | UUID | Who received |
| status | ENUM | pending, received, rejected |

---

## 6. Security Dashboard

```
┌─────────────────────────────────────────────────────────────┐
│  SECURITY DASHBOARD | 22-Jan-2026                           │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  CURRENT STATUS                                              │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐           │
│  │ Visitors    │ │ Gate Passes │ │ Vehicles    │           │
│  │ Inside: 8   │ │ Active: 3   │ │ Inside: 12  │           │
│  │ Today: 24   │ │ Today: 7    │ │ Today: 45   │           │
│  └─────────────┘ └─────────────┘ └─────────────┘           │
│                                                              │
│  VISITORS CURRENTLY INSIDE                                  │
│  | Badge  | Name          | Meeting    | Entry  | Status | │
│  |--------|---------------|------------|--------|--------| │
│  | V-0142 | Rajesh Verma  | Mr. Sharma | 10:30  | Active | │
│  | V-0143 | Priya Singh   | Principal  | 11:00  | Active | │
│  | V-0144 | Amit Kumar    | Accounts   | 11:15  | Active | │
│                                                              │
│  PENDING GATE PASSES                                        │
│  | Student      | Class | Pickup     | Time  | Status    | │
│  |--------------|-------|------------|-------|-----------|│
│  | Aarav Sharma | 10-A  | Father     | 11:30 | Approved  | │
│  | Diya Patel   | 8-B   | Mother     | 12:00 | Pending   | │
│                                                              │
│  [New Visitor] [New Gate Pass] [Vehicle Entry]              │
└─────────────────────────────────────────────────────────────┘
```

---

## 7. Pre-Registration

### 7.1 Pre-Registered Visitor

**Entity: VisitorPreRegistration**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| visitor_name | VARCHAR(200) | Name |
| phone | VARCHAR(20) | Phone |
| email | VARCHAR(255) | Email |
| expected_date | DATE | Visit date |
| expected_time | TIME | Expected time |
| purpose | TEXT | Purpose |
| host_id | UUID | Host staff |
| approval_code | VARCHAR(20) | Entry code |
| status | ENUM | pending, approved, used, expired |

Staff can pre-register expected visitors, who receive an approval code for quick check-in.

---

## 8. Reports

```
Visitor Report | January 2026

| Metric              | Count |
|---------------------|-------|
| Total Visitors      | 450   |
| Parent Visits       | 180   |
| Vendor Visits       | 120   |
| Deliveries          | 95    |
| Other               | 55    |

| Week    | Mon | Tue | Wed | Thu | Fri | Sat |
|---------|-----|-----|-----|-----|-----|-----|
| Week 1  | 22  | 18  | 25  | 20  | 30  | 8   |
| Week 2  | 24  | 20  | 28  | 22  | 32  | 10  |
| Week 3  | 20  | 22  | 24  | 25  | 28  | 12  |

Peak Hours: 10:00 AM - 12:00 PM
```

---

## 9. API Endpoints

```
# Visitors
POST   /api/v1/visitors/check-in            # Check in
POST   /api/v1/visitors/check-out           # Check out
GET    /api/v1/visitors/current             # Currently inside
GET    /api/v1/visitors/history             # Visit history

# Gate Passes
POST   /api/v1/gate-passes                  # Request pass
GET    /api/v1/gate-passes/pending          # Pending approvals
POST   /api/v1/gate-passes/{id}/approve     # Approve
POST   /api/v1/gate-passes/{id}/use         # Mark as used

# Vehicles
POST   /api/v1/vehicles/entry               # Vehicle entry
POST   /api/v1/vehicles/exit                # Vehicle exit

# Pre-registration
POST   /api/v1/visitors/pre-register        # Pre-register
GET    /api/v1/visitors/pre-register/today  # Today's expected
```

---

## 10. Business Rules

| Rule | Description |
|------|-------------|
| ID Mandatory | ID proof required for all visitors |
| Photo Capture | Photo mandatory for first-time visitors |
| Parent Pickup | Only registered guardians can pickup |
| Host Confirmation | Host must confirm visitor arrival |
| Badge Return | Badge must be returned at exit |
| Restricted Hours | Visitors not allowed during exams |
| Gate Pass Validity | Gate pass valid only for specified date |

---

## 11. Related Documents

- [03-student-management.md](./03-student-management.md) - Student guardians
- [13-communication-system.md](./13-communication-system.md) - Notifications
- [index.md](./index.md) - Main PRD index

---

**Previous**: [19-inventory-assets.md](./19-inventory-assets.md)
**Next**: [21-analytics-dashboards.md](./21-analytics-dashboards.md)
