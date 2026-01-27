# 16 - Certificate Generation

**Parent Document**: [index.md](./index.md)
**Status**: Draft
**Phase**: 1 (MVP)

---

## 1. Overview

The Certificate Generation module handles creation, approval, and issuance of various student certificates including Transfer Certificate (TC), Bonafide, Character Certificate, and custom certificates.

---

## 2. Certificate Types

### 2.1 Standard Certificates

| Type | Code | Purpose | Approval Required |
|------|------|---------|-------------------|
| Transfer Certificate | TC | Student leaving school | Principal |
| Bonafide Certificate | BC | Proof of enrollment | Admin |
| Character Certificate | CC | Character reference | Principal |
| Study Certificate | SC | Proof of study period | Admin |
| Conduct Certificate | CON | Conduct during study | Class Teacher + Principal |
| Migration Certificate | MC | Board migration | Principal |
| Date of Birth Certificate | DOB | DOB verification | Admin |

### 2.2 Custom Certificates

| Type | Code | Purpose |
|------|------|---------|
| Sports Achievement | SA | Sports recognition |
| Academic Excellence | AE | Merit certificates |
| Participation | PC | Event participation |
| Course Completion | CC | Course/training completion |

---

## 3. Certificate Request

### 3.1 Request Entity

**Entity: CertificateRequest**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| branch_id | UUID | Branch reference |
| student_id | UUID | Student reference |
| certificate_type | VARCHAR(20) | Certificate type code |
| request_number | VARCHAR(50) | Unique request number |
| requested_by | UUID | Requester (parent/staff) |
| request_date | DATE | Request date |
| purpose | TEXT | Purpose of certificate |
| required_date | DATE | When needed |
| copies_required | INT | Number of copies |
| delivery_mode | ENUM | collect, courier, email |
| status | ENUM | pending, approved, rejected, generated, delivered |
| priority | ENUM | normal, urgent |
| remarks | TEXT | Additional notes |
| fee_amount | DECIMAL | Certificate fee |
| fee_paid | BOOLEAN | Fee status |
| payment_id | UUID | Payment reference |

### 3.2 Request Interface

```
┌─────────────────────────────────────────────────────────────┐
│  REQUEST CERTIFICATE                                         │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Student: Aarav Sharma | Class: 10-A | Adm No: 2024001     │
│                                                              │
│  Certificate Type: [Transfer Certificate    ▼]              │
│                                                              │
│  Purpose:                                                    │
│  ┌─────────────────────────────────────────────────────┐    │
│  │ Relocation to another city                           │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                              │
│  Required By: [15-Feb-2026]                                 │
│  Copies Required: [2]                                       │
│  Priority: [Normal ▼]                                       │
│                                                              │
│  Delivery: ○ Collect from school                            │
│            ○ Courier (additional charges)                   │
│            ○ Email (PDF)                                    │
│                                                              │
│  Fee: ₹100 per copy = ₹200                                  │
│                                                              │
│  [Submit Request]                                            │
└─────────────────────────────────────────────────────────────┘
```

---

## 4. Approval Workflow

### 4.1 Workflow by Type

**Transfer Certificate**:
```
Request → Fee Payment → Class Teacher Review →
Account Clearance → Library Clearance →
Principal Approval → Generate → Deliver
```

**Bonafide Certificate**:
```
Request → Fee Payment → Admin Approval → Generate → Deliver
```

### 4.2 Clearance Checks

**Entity: CertificateClearance**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| request_id | UUID | Request reference |
| department | VARCHAR(50) | Department name |
| status | ENUM | pending, cleared, blocked |
| cleared_by | UUID | Cleared by staff |
| cleared_at | TIMESTAMP | Clearance time |
| remarks | TEXT | Any dues/issues |

**Required Clearances for TC**:
- Fee Department (no dues)
- Library (all books returned)
- Lab (all equipment returned)
- Hostel (if applicable)

### 4.3 Approval Entity

**Entity: CertificateApproval**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| request_id | UUID | Request reference |
| approver_id | UUID | Approving authority |
| level | INT | Approval level |
| status | ENUM | pending, approved, rejected |
| comments | TEXT | Approver comments |
| approved_at | TIMESTAMP | Approval time |

---

## 5. Certificate Templates

### 5.1 Template Entity

**Entity: CertificateTemplate**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| certificate_type | VARCHAR(20) | Certificate type |
| name | VARCHAR(200) | Template name |
| content_html | TEXT | HTML template |
| header_image | VARCHAR(500) | Header logo/image |
| footer_image | VARCHAR(500) | Footer/seal |
| page_size | VARCHAR(10) | A4, Letter, etc. |
| orientation | ENUM | portrait, landscape |
| variables | JSONB | Available variables |
| is_default | BOOLEAN | Default template |
| is_active | BOOLEAN | Active status |

### 5.2 Template Variables

```yaml
common_variables:
  - {student_name}
  - {father_name}
  - {mother_name}
  - {date_of_birth}
  - {admission_number}
  - {class}
  - {section}
  - {admission_date}
  - {school_name}
  - {school_address}
  - {certificate_number}
  - {issue_date}

tc_specific:
  - {leaving_date}
  - {reason_for_leaving}
  - {conduct}
  - {subjects_studied}
  - {last_exam_appeared}
  - {qualified_for_promotion}
  - {fees_paid_till}
  - {tc_number}
```

### 5.3 Transfer Certificate Format

```
┌─────────────────────────────────────────────────────────────┐
│                    [SCHOOL LOGO]                             │
│                 ABC PUBLIC SCHOOL                            │
│    (Affiliated to CBSE, New Delhi - Affiliation No: XXXX)   │
│              123 Main Street, Mumbai - 400001               │
│                                                              │
│                  TRANSFER CERTIFICATE                        │
│                     TC No: TC/2026/0042                      │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│ 1. Name of Student         : AARAV SHARMA                   │
│ 2. Father's Name           : RAJESH SHARMA                  │
│ 3. Mother's Name           : PRIYA SHARMA                   │
│ 4. Nationality             : INDIAN                         │
│ 5. Date of Birth           : 15-05-2010                     │
│    (in figures and words)    (Fifteenth May Two Thousand Ten)│
│ 6. Date of Admission       : 01-04-2020                     │
│ 7. Class in which admitted : V (Five)                       │
│ 8. Class in which studying : X (Ten)                        │
│ 9. Date of Leaving         : 31-01-2026                     │
│ 10. Reason for Leaving     : Parent's Transfer              │
│ 11. Subjects Studied       : English, Hindi, Mathematics,   │
│                              Science, Social Science        │
│ 12. Last Examination       : Half Yearly 2025-26           │
│ 13. Qualified for Promotion: Yes, to Class XI              │
│ 14. Fees Paid Till         : January 2026                   │
│ 15. Total Working Days     : 180                            │
│ 16. Days Present           : 165                            │
│ 17. Conduct & Character    : Good                           │
│ 18. Games & Sports         : Cricket, Athletics             │
│ 19. Extra-Curricular       : Quiz Club, Science Club        │
│                                                              │
│ Certified that the above information is correct as per      │
│ school records.                                              │
│                                                              │
│ Date: 01-02-2026                                            │
│                                                              │
│                                                              │
│ Class Teacher          Principal              [SCHOOL SEAL] │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## 6. Certificate Generation

### 6.1 Generated Certificate

**Entity: Certificate**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| request_id | UUID | Request reference |
| certificate_number | VARCHAR(50) | Unique certificate number |
| template_id | UUID | Template used |
| content_data | JSONB | Variable values |
| pdf_url | VARCHAR(500) | Generated PDF |
| generated_by | UUID | Generator |
| generated_at | TIMESTAMP | Generation time |
| is_duplicate | BOOLEAN | Duplicate issue |
| original_id | UUID | If duplicate |
| qr_code | VARCHAR(500) | Verification QR |
| verification_code | VARCHAR(20) | Manual verification code |

### 6.2 Certificate Numbering

```
TC/{YEAR}/{BRANCH}/{SEQUENCE}
Example: TC/2026/MAIN/0042

BC/{YEAR}/{SEQUENCE}
Example: BC/2026/0156
```

### 6.3 QR Code Verification

QR Code contains:
- Certificate number
- Student name
- Certificate type
- Issue date
- Verification URL

Scanning QR leads to: `https://school.edu/verify/TC-2026-0042`

---

## 7. Delivery & Collection

### 7.1 Collection Process

1. SMS/Email sent when ready
2. Parent/guardian visits school
3. Identity verification
4. Sign acknowledgement register
5. Collect certificate

### 7.2 Delivery Entity

**Entity: CertificateDelivery**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| certificate_id | UUID | Certificate reference |
| delivery_mode | ENUM | collect, courier, email |
| delivered_to | VARCHAR(200) | Recipient name |
| delivered_at | TIMESTAMP | Delivery time |
| collected_by | VARCHAR(200) | If collected |
| relation | VARCHAR(50) | Relation to student |
| id_proof | VARCHAR(100) | ID proof provided |
| signature_url | VARCHAR(500) | Digital signature |
| courier_tracking | VARCHAR(100) | Courier tracking |

---

## 8. Reports

### 8.1 Certificate Register

```
Certificate Register | January 2026

| Sr | Cert No       | Student        | Type | Issue Date | Status    |
|----|---------------|----------------|------|------------|-----------|
| 1  | TC/2026/0041  | Priya Singh    | TC   | 15-Jan     | Delivered |
| 2  | TC/2026/0042  | Aarav Sharma   | TC   | 01-Feb     | Pending   |
| 3  | BC/2026/0155  | Rahul Kumar    | BC   | 18-Jan     | Delivered |
| 4  | BC/2026/0156  | Ananya Patel   | BC   | 20-Jan     | Generated |

Total TCs Issued: 5
Total BCs Issued: 12
Total Revenue: ₹3,400
```

---

## 9. API Endpoints

```
# Certificate Requests
GET    /api/v1/certificates/requests        # List requests
POST   /api/v1/certificates/requests        # Create request
GET    /api/v1/certificates/requests/{id}   # Get request
PUT    /api/v1/certificates/requests/{id}   # Update request

# Clearances
GET    /api/v1/certificates/requests/{id}/clearances # Get clearances
POST   /api/v1/certificates/requests/{id}/clear # Clear department

# Approvals
POST   /api/v1/certificates/requests/{id}/approve # Approve
POST   /api/v1/certificates/requests/{id}/reject  # Reject

# Generation
POST   /api/v1/certificates/requests/{id}/generate # Generate PDF
GET    /api/v1/certificates/{id}/download   # Download PDF

# Verification
GET    /api/v1/certificates/verify/{code}   # Verify certificate

# Delivery
POST   /api/v1/certificates/{id}/deliver    # Mark delivered
```

---

## 10. Business Rules

| Rule | Description |
|------|-------------|
| Fee Required | Certificate not generated until fee paid |
| Clearance Required | TC needs all department clearances |
| Principal Approval | TC and CC need principal signature |
| Duplicate Fee | Duplicate certificates cost 2x |
| Processing Time | Standard 3-5 days, Urgent 24-48 hours |
| Original Record | Original TC can only be issued once |
| Verification | All certificates have verification QR/code |

---

## 11. Related Documents

- [03-student-management.md](./03-student-management.md) - Student data
- [12-fees-payments.md](./12-fees-payments.md) - Certificate fees
- [index.md](./index.md) - Main PRD index

---

**Previous**: [15-student-portal.md](./15-student-portal.md)
**Next**: [17-transport-management.md](./17-transport-management.md)
