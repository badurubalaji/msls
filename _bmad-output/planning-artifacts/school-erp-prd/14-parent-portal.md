# 14 - Parent Portal

**Parent Document**: [index.md](./index.md)
**Status**: Draft
**Phase**: 1 (MVP)

---

## 1. Overview

The Parent Portal provides parents/guardians with real-time access to their child's academic progress, attendance, fees, and school communications through web and mobile interfaces.

---

## 2. Portal Access

### 2.1 Registration & Login

**Registration Flow**:
1. School creates student record with parent contact
2. System sends invitation SMS/Email to parent
3. Parent clicks link and sets password
4. OTP verification completes registration
5. Parent can now access portal

**Login Options**:
- Email + Password
- Phone + OTP
- Google/Microsoft SSO (optional)

### 2.2 Multi-Child Support

Parents with multiple children in the school:
- Single login for all children
- Switch between child profiles
- Consolidated fee view
- Unified notifications

---

## 3. Dashboard

### 3.1 Dashboard Layout

```
┌─────────────────────────────────────────────────────────────┐
│  Welcome, Mr. Sharma                   [🔔 3] [Aarav ▼]     │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  AARAV SHARMA | Class 10-A | Roll No: 15                   │
│                                                              │
│  ┌─────────────────┐  ┌─────────────────┐                   │
│  │ 📅 ATTENDANCE   │  │ 💰 FEES         │                   │
│  │                 │  │                 │                   │
│  │ This Month: 92% │  │ Due: ₹15,000   │                   │
│  │ Today: Present  │  │ Due: 31-Jan    │                   │
│  │                 │  │                 │                   │
│  │ [View Details]  │  │ [Pay Now]       │                   │
│  └─────────────────┘  └─────────────────┘                   │
│                                                              │
│  ┌─────────────────┐  ┌─────────────────┐                   │
│  │ 📝 HOMEWORK     │  │ 📊 PERFORMANCE  │                   │
│  │                 │  │                 │                   │
│  │ Pending: 2      │  │ Last Exam: 78% │                   │
│  │ Due Today: 1    │  │ Rank: 8/42     │                   │
│  │                 │  │                 │                   │
│  │ [View All]      │  │ [View Report]   │                   │
│  └─────────────────┘  └─────────────────┘                   │
│                                                              │
│  📢 RECENT NOTICES                                          │
│  ├─ PTM scheduled for 30-Jan (Today)                       │
│  ├─ Holiday on 26-Jan - Republic Day                       │
│  └─ Half Yearly exam dates announced                       │
│                                                              │
│  📅 UPCOMING                                                │
│  ├─ Math Assignment due - 25-Jan                           │
│  ├─ Science Project submission - 28-Jan                    │
│  └─ Half Yearly Exam starts - 15-Feb                       │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## 4. Feature Modules

### 4.1 Attendance

**Features**:
- Daily attendance status
- Monthly attendance calendar
- Attendance percentage trends
- Absent day notifications
- Leave application for child

```
Attendance | Aarav | January 2026

     Mon  Tue  Wed  Thu  Fri  Sat  Sun
      --   --    1    2    3    4    5
                ✓    ✓    ✓    H    -

       6    7    8    9   10   11   12
       ✓    ✓    L    ✓    ✓    ✓    -

      13   14   15   16   17   18   19
       ✓    ✓    ✓    A    ✓    ✓    -

      20   21   22   23   24   25   26
       ✓    ✓    ✓    ✓    ✓    H    H

Legend: ✓ Present | A Absent | L Late | H Holiday | - Weekend

This Month: Present: 18 | Absent: 1 | Late: 1 | Holiday: 3
Attendance: 94.7%
```

### 4.2 Academic Performance

**Features**:
- Exam results and report cards
- Subject-wise marks
- Grade trends over time
- Class rank
- Comparative performance

```
Academic Performance | Aarav | 2025-26

Exam Results:
| Exam          | Marks | %    | Grade | Rank  |
|---------------|-------|------|-------|-------|
| Unit Test 1   | 425/500| 85% | A     | 5/42  |
| Unit Test 2   | 410/500| 82% | A     | 7/42  |
| Half Yearly   | 368/500| 73% | B+    | 8/42  |

Subject Trends:
| Subject      | UT1 | UT2 | HY  | Trend |
|--------------|-----|-----|-----|-------|
| Mathematics  | 88  | 85  | 83  | ↓     |
| Science      | 82  | 80  | 75  | ↓     |
| English      | 85  | 84  | 68  | ↓     |
| Hindi        | 78  | 80  | 72  | ↔     |
| Soc. Science | 92  | 81  | 70  | ↓     |

[Download Report Card]
```

### 4.3 Fees & Payments

**Features**:
- Fee structure view
- Outstanding dues
- Payment history
- Online payment
- Download receipts
- Payment reminders

```
Fees | Aarav | 2025-26

Outstanding:
| Description              | Amount   | Due Date | Status  |
|--------------------------|----------|----------|---------|
| Tuition Fee (Jan-Mar)    | ₹15,000  | 31-Jan   | Due     |
| Transport Fee (Jan)      | ₹2,000   | 31-Jan   | Due     |
|--------------------------|----------|----------|---------|
| Total Due                | ₹17,000  |          |         |

[Pay Now ₹17,000]

Payment History:
| Date       | Receipt #      | Amount   | Mode   |
|------------|----------------|----------|--------|
| 01-Oct-25  | RCP-2025-1234  | ₹17,000  | UPI    |
| 01-Jul-25  | RCP-2025-0987  | ₹17,000  | Card   |
| 01-Apr-25  | RCP-2025-0654  | ₹25,000  | Bank   |
```

### 4.4 Homework & Assignments

**Features**:
- View assigned homework
- Track submission status
- View grades and feedback
- Download attachments

```
Homework | Aarav | This Week

Pending (2):
| Subject | Assignment            | Due Date | Status     |
|---------|----------------------|----------|------------|
| Math    | Chapter 5 Problems   | 25-Jan   | Not Started|
| English | Essay Writing        | 26-Jan   | In Progress|

Completed (3):
| Subject | Assignment            | Submitted | Marks |
|---------|----------------------|-----------|-------|
| Science | Lab Report           | 22-Jan    | 18/20 |
| Hindi   | Grammar Exercises    | 20-Jan    | 9/10  |
| SST     | Map Work             | 18-Jan    | 8/10  |
```

### 4.5 Communication

**Features**:
- School notices and circulars
- Direct message to teachers
- PTM scheduling
- Complaint/feedback submission
- Emergency alerts

```
Messages | Aarav

📢 Notices (3 unread)
├─ [NEW] PTM Schedule - 30 Jan 2026
├─ [NEW] Holiday Notice - Republic Day
├─ Half Yearly Exam Dates
└─ Winter Uniform Advisory

💬 Conversations
├─ Ms. Sharma (Class Teacher) - Last: Yesterday
├─ Mr. Kumar (Math Teacher) - Last: 3 days ago
└─ Transport In-charge - Last: 1 week ago

[New Message] [View All]
```

### 4.6 Timetable

**Features**:
- Class timetable view
- Exam schedule
- Holiday calendar
- School events

---

## 5. Mobile App Features

### 5.1 Push Notifications

| Event | Notification |
|-------|--------------|
| Absent | "Aarav was marked absent today" |
| Fee Due | "Fee payment of ₹17,000 due in 3 days" |
| Result | "Half Yearly results published. View now" |
| Homework | "New Math homework assigned. Due: 25-Jan" |
| Notice | "New circular: PTM Schedule" |
| Emergency | "URGENT: School closed tomorrow" |

### 5.2 Quick Actions

- Mark attendance leave request
- Quick fee payment
- View today's timetable
- Contact school
- Emergency contacts

---

## 6. Settings & Profile

### 6.1 Profile Management

- View/update contact details
- Change password
- Notification preferences
- Language preference
- Link additional children

### 6.2 Notification Settings

```
Notification Preferences

| Category        | Push | Email | SMS |
|-----------------|------|-------|-----|
| Attendance      | ☑    | ☑     | ☑   |
| Fee Reminders   | ☑    | ☑     | ☑   |
| Exam Results    | ☑    | ☑     | ☐   |
| Homework        | ☑    | ☐     | ☐   |
| Notices         | ☑    | ☑     | ☐   |
| Emergency       | ☑    | ☑     | ☑   |
```

---

## 7. API Endpoints

```
# Dashboard
GET    /api/v1/parent/dashboard             # Dashboard data

# Children
GET    /api/v1/parent/children              # List children
GET    /api/v1/parent/children/{id}         # Child details

# Attendance
GET    /api/v1/parent/children/{id}/attendance # Attendance
POST   /api/v1/parent/children/{id}/leave   # Apply leave

# Academics
GET    /api/v1/parent/children/{id}/results # Exam results
GET    /api/v1/parent/children/{id}/report-cards # Report cards
GET    /api/v1/parent/children/{id}/homework # Homework

# Fees
GET    /api/v1/parent/children/{id}/fees    # Fee details
GET    /api/v1/parent/children/{id}/payments # Payment history
POST   /api/v1/parent/payments              # Make payment

# Communication
GET    /api/v1/parent/notices               # School notices
GET    /api/v1/parent/conversations         # Messages
POST   /api/v1/parent/conversations         # New message

# Timetable
GET    /api/v1/parent/children/{id}/timetable # Class timetable
```

---

## 8. Related Documents

- [03-student-management.md](./03-student-management.md) - Student data
- [12-fees-payments.md](./12-fees-payments.md) - Fee system
- [13-communication-system.md](./13-communication-system.md) - Communications
- [index.md](./index.md) - Main PRD index

---

**Previous**: [13-communication-system.md](./13-communication-system.md)
**Next**: [15-student-portal.md](./15-student-portal.md)
