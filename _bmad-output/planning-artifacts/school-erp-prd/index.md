# School ERP Platform - Product Requirements Document

## Document Index

**Version**: 1.0.0
**Last Updated**: 2026-01-22
**Status**: Draft
**Project Code**: MSLS (Multi-School Learning System)

---

## 1. Executive Summary

### 1.1 Product Vision

Build a **future-ready, on-premise-first School ERP platform** that supports institutions from **Nursery to PhD level**, enabling academic excellence, operational efficiency, parent transparency, and data-driven decision-making — without mandatory cloud dependency.

### 1.2 Core Value Proposition

| Aspect | Description |
|--------|-------------|
| **Data Sovereignty** | On-premise first, schools own their data |
| **Flexibility** | SaaS or License-based deployment |
| **Configurability** | Module-based activation per institution type |
| **Scalability** | Single school to multi-branch consortiums |
| **Offline Capable** | Works on LAN without internet dependency |

### 1.3 Target Market

- **Primary Focus (Phase 1)**: K-12 Schools (Nursery to Class 12)
- **Secondary (Phase 2+)**: Junior Colleges, Degree Colleges, Universities
- **Deployment**: India-first, Global-ready

---

## 2. Document Structure

This PRD is organized into modular sub-documents for maintainability. Each document focuses on a specific functional area.

### 2.1 Foundation Documents

| # | Document | Description | Link |
|---|----------|-------------|------|
| 01 | Technical Architecture | Tech stack, deployment models, multi-tenancy | [01-technical-architecture.md](./01-technical-architecture.md) |
| 02 | Core Foundation | Institution setup, Users, RBAC, Configuration | [02-core-foundation.md](./02-core-foundation.md) |

### 2.2 Academic Documents

| # | Document | Description | Link |
|---|----------|-------------|------|
| 03 | Student Management | Student lifecycle, profiles, health, discipline | [03-student-management.md](./03-student-management.md) |
| 04 | Academic Operations | Classes, Sections, Timetable, Attendance | [04-academic-operations.md](./04-academic-operations.md) |
| 05 | Admissions | Enquiry to enrollment workflow | [05-admissions.md](./05-admissions.md) |
| 06 | Examinations & Grading | Exams, marks, report cards, gradebook | [06-examinations-grading.md](./06-examinations-grading.md) |

### 2.3 Learning & Assessment Documents

| # | Document | Description | Link |
|---|----------|-------------|------|
| 07 | Homework & Assignments | Assignment creation, submission, grading | [07-homework-assignments.md](./07-homework-assignments.md) |
| 08 | Online Quiz & Assessment | Question bank, quizzes, auto-grading | [08-online-quiz-assessment.md](./08-online-quiz-assessment.md) |
| 09 | Digital Classroom | Class recording, content replay, e-library | [09-digital-classroom.md](./09-digital-classroom.md) |

### 2.4 Staff & HR Documents

| # | Document | Description | Link |
|---|----------|-------------|------|
| 10 | Staff Management | Staff profiles, attendance, payroll | [10-staff-management.md](./10-staff-management.md) |
| 11 | Leave Management | Leave types, application, approval workflow | [11-leave-management.md](./11-leave-management.md) |

### 2.5 Finance Documents

| # | Document | Description | Link |
|---|----------|-------------|------|
| 12 | Fees & Payments | Fee structure, collection, receipts, reports | [12-fees-payments.md](./12-fees-payments.md) |

### 2.6 Communication Documents

| # | Document | Description | Link |
|---|----------|-------------|------|
| 13 | Communication System | Notices, SMS, Email, WhatsApp integration | [13-communication-system.md](./13-communication-system.md) |
| 14 | Parent Portal | Parent-facing features and mobile app | [14-parent-portal.md](./14-parent-portal.md) |
| 15 | Student Portal | Student-facing features and mobile app | [15-student-portal.md](./15-student-portal.md) |

### 2.7 Operations Documents

| # | Document | Description | Link |
|---|----------|-------------|------|
| 16 | Certificate Generation | TC, Bonafide, Character certificates | [16-certificate-generation.md](./16-certificate-generation.md) |
| 17 | Transport Management | Routes, vehicles, tracking | [17-transport-management.md](./17-transport-management.md) |
| 18 | Library Management | Books, issue/return, catalog | [18-library-management.md](./18-library-management.md) |
| 19 | Inventory & Assets | Equipment, labs, consumables | [19-inventory-assets.md](./19-inventory-assets.md) |
| 20 | Visitor & Gate Management | Visitor log, gate pass, security | [20-visitor-gate-management.md](./20-visitor-gate-management.md) |

### 2.8 Advanced Features Documents

| # | Document | Description | Link |
|---|----------|-------------|------|
| 21 | Analytics & Dashboards | Reports, charts, insights | [21-analytics-dashboards.md](./21-analytics-dashboards.md) |
| 22 | AI Capabilities | Predictions, recommendations, automation | [22-ai-capabilities.md](./22-ai-capabilities.md) |

---

## 3. User Personas

| Persona | Description | Primary Modules |
|---------|-------------|-----------------|
| **Super Admin** | Vendor/Platform owner, manages tenants | Tenant management, licensing, global config |
| **School Admin** | School-level administrator | All modules, school-wide settings |
| **Branch Admin** | Branch-level administrator | Branch-specific modules and reports |
| **Principal** | Academic head | Dashboards, approvals, reports |
| **Teacher** | Academic staff | Attendance, grades, homework, communication |
| **Non-Teaching Staff** | Administrative staff | Assigned modules based on role |
| **Accountant** | Finance staff | Fees, payments, financial reports |
| **Librarian** | Library staff | Library module |
| **Transport In-charge** | Transport management | Transport module |
| **Parent** | Student's guardian | Parent portal, fees, communication |
| **Student** | Learner | Student portal, homework, results |

---

## 4. Deployment Models

### 4.1 SaaS Mode (Cloud)

```
┌─────────────────────────────────────────────────────┐
│                  SAAS DEPLOYMENT                     │
├─────────────────────────────────────────────────────┤
│  Tenant 1 (School A)  │  Tenant 2 (School B)  │ ... │
│  ├── Branch 1         │  ├── Branch 1         │     │
│  ├── Branch 2         │  └── Branch 2         │     │
│  └── Branch 3         │                       │     │
├─────────────────────────────────────────────────────┤
│            Shared Infrastructure                     │
│     (Row-Level Tenant Isolation with RLS)           │
└─────────────────────────────────────────────────────┘
```

- Multiple schools on shared infrastructure
- Row-level tenant isolation using PostgreSQL RLS
- Centralized updates and maintenance
- Pay-per-use or subscription pricing

### 4.2 On-Premise Mode (License)

```
┌─────────────────────────────────────────────────────┐
│               ON-PREMISE DEPLOYMENT                  │
├─────────────────────────────────────────────────────┤
│  Single Tenant (This School)                        │
│  ├── Branch 1 (Main Campus)                         │
│  ├── Branch 2 (Junior Wing)                         │
│  └── Branch 3 (Senior Wing)                         │
├─────────────────────────────────────────────────────┤
│         School-Owned Server / Local Network          │
│              (Tenant features disabled)              │
└─────────────────────────────────────────────────────┘
```

- Installed on school-owned servers
- Full data ownership and control
- Works on LAN without internet
- One-time license + annual maintenance

---

## 5. Phase-wise Roadmap

### Phase 1: K-12 MVP (Core ERP)

**Modules Included:**
- Core Foundation (Institution, Users, RBAC, Config)
- Student Management
- Staff Management
- Admissions
- Class/Section Management
- Attendance (Student + Staff)
- Timetable/Scheduling
- Examinations & Grading
- Homework/Assignments
- Leave Management
- Fees & Payments
- Communication (Notices, SMS, Email)
- Certificate Generation
- Parent Portal
- Student Portal

### Phase 2: Extended Features

**Modules Included:**
- Online Quiz/Assessment System
- Class Recording & Digital Classroom
- Transport Management
- Library Management
- Inventory & Assets
- Visitor/Gate Management
- Analytics Dashboards

### Phase 3: Advanced & AI

**Modules Included:**
- AI-powered Insights
- Predictive Analytics
- Research/Thesis Tracking (for universities)
- Advanced Reporting
- CCTV Integration (optional)

---

## 6. Non-Functional Requirements

| Requirement | Target |
|-------------|--------|
| **Scalability** | Support 50,000+ concurrent users |
| **Performance** | Sub-second response for core actions |
| **Availability** | 99.9% uptime (SaaS mode) |
| **Offline Support** | Full LAN operation without internet |
| **Security** | RBAC, encryption at rest/transit, audit logs |
| **Compliance** | FERPA-like principles, GDPR-ready |
| **Localization** | Multi-language UI and reports |
| **Mobile** | PWA + Native apps for iOS/Android |

---

## 7. Success Metrics

| Metric | Target |
|--------|--------|
| Admin workload reduction | 40% |
| Parent engagement rate | 70%+ active users |
| Attendance tracking accuracy | 99% |
| Fee collection efficiency | 20% improvement |
| Report generation time | < 5 seconds |
| System adoption (6 months) | 80% daily active users |

---

## 8. Risks & Mitigation

| Risk | Impact | Mitigation |
|------|--------|------------|
| Data privacy concerns | High | On-premise default, encryption, RBAC |
| Poor internet in schools | Medium | Offline-first architecture |
| User adoption resistance | Medium | Simple UX, training modules |
| Competition from established players | High | Focus on niche + unique features |
| Scope creep | Medium | Phased approach, strict prioritization |

---

## 9. Document Maintenance

| Action | Frequency |
|--------|-----------|
| Review and update | Monthly |
| Version increment | On significant changes |
| Stakeholder review | Quarterly |

---

## 10. Next Steps

1. Review and finalize sub-documents
2. Validate with potential customers
3. Create technical architecture document
4. Define database schema
5. Begin Phase 1 development

---

**Document Owner**: Ashulabs
**Contributors**: BMad Master
**Review Status**: Pending
