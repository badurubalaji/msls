# 12 - Fees & Payments

**Parent Document**: [index.md](./index.md)
**Status**: Draft
**Phase**: 1 (MVP)

---

## 1. Overview

The Fees & Payments module handles complete fee lifecycle including structure definition, invoicing, collection, receipts, dues tracking, and financial reporting.

---

## 2. Fee Structure

### 2.1 Fee Category

**Entity: FeeCategory**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| name | VARCHAR(100) | Category name |
| code | VARCHAR(20) | Short code |
| description | TEXT | Description |
| is_recurring | BOOLEAN | Recurring fee (monthly/term) |
| is_refundable | BOOLEAN | Can be refunded |
| accounting_code | VARCHAR(50) | For accounting integration |
| display_order | INT | Display ordering |
| is_active | BOOLEAN | Active status |

**Common Fee Categories**:
```yaml
categories:
  - name: Tuition Fee
    code: TF
    recurring: true
    refundable: false
  - name: Admission Fee
    code: AF
    recurring: false
    refundable: false
  - name: Transport Fee
    code: TRANS
    recurring: true
    refundable: true
  - name: Library Fee
    code: LIB
    recurring: false
    refundable: false
  - name: Lab Fee
    code: LAB
    recurring: true
    refundable: false
  - name: Sports Fee
    code: SPORTS
    recurring: false
    refundable: false
  - name: Exam Fee
    code: EXAM
    recurring: true
    refundable: false
  - name: Caution Deposit
    code: CD
    recurring: false
    refundable: true
```

### 2.2 Fee Structure Definition

**Entity: FeeStructure**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| branch_id | UUID | Branch (NULL = all) |
| academic_year_id | UUID | Academic year |
| class_id | UUID | Class (NULL = all classes) |
| name | VARCHAR(200) | Structure name |
| description | TEXT | Description |
| total_amount | DECIMAL | Total annual fee |
| is_active | BOOLEAN | Active status |
| effective_from | DATE | Start date |
| effective_until | DATE | End date |

### 2.3 Fee Structure Items

**Entity: FeeStructureItem**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| fee_structure_id | UUID | Parent structure |
| fee_category_id | UUID | Fee category |
| amount | DECIMAL | Fee amount |
| frequency | ENUM | one_time, monthly, quarterly, term, annual |
| installments | INT | Number of installments |
| due_day | INT | Day of month due |
| late_fee_applicable | BOOLEAN | Apply late fee |

**Example Structure**:
```
Class 5 Fee Structure - 2025-26

| Category       | Amount   | Frequency  | Installments |
|----------------|----------|------------|--------------|
| Admission Fee  | ₹10,000  | One-time   | 1            |
| Tuition Fee    | ₹60,000  | Quarterly  | 4            |
| Transport Fee  | ₹24,000  | Monthly    | 12           |
| Library Fee    | ₹2,000   | Annual     | 1            |
| Lab Fee        | ₹5,000   | Annual     | 1            |
| Sports Fee     | ₹3,000   | Annual     | 1            |
| Exam Fee       | ₹4,000   | Per-term   | 2            |
| Caution Deposit| ₹5,000   | One-time   | 1            |

Total Annual: ₹1,13,000
```

---

## 3. Student Fee Assignment

### 3.1 Fee Assignment

**Entity: StudentFee**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| student_id | UUID | Student reference |
| enrollment_id | UUID | Enrollment reference |
| fee_structure_id | UUID | Assigned structure |
| applicable_amount | DECIMAL | Amount after discounts |
| discount_id | UUID | Applied discount |
| discount_amount | DECIMAL | Discount value |
| custom_adjustments | JSONB | Custom fee changes |
| assigned_at | TIMESTAMP | Assignment date |
| assigned_by | UUID | Assigned by |

### 3.2 Discount/Concession Types

**Entity: FeeDiscount**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| name | VARCHAR(100) | Discount name |
| code | VARCHAR(20) | Discount code |
| type | ENUM | percentage, fixed, category_wise |
| value | DECIMAL | Discount value |
| applicable_categories | JSONB | Categories it applies to |
| criteria | JSONB | Eligibility criteria |
| max_students | INT | Maximum beneficiaries |
| valid_from | DATE | Start date |
| valid_until | DATE | End date |
| is_active | BOOLEAN | Active status |

**Discount Types**:
```yaml
discounts:
  - name: Sibling Discount
    type: percentage
    value: 10
    criteria: { has_sibling: true }

  - name: Staff Ward
    type: percentage
    value: 50
    criteria: { parent_is_staff: true }

  - name: Merit Scholarship
    type: percentage
    value: 25
    criteria: { previous_percentage: ">90" }

  - name: EWS Concession
    type: category_wise
    value: { tuition: 100, transport: 50 }
    criteria: { category: "EWS" }
```

---

## 4. Invoicing

### 4.1 Invoice Generation

**Entity: Invoice**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| branch_id | UUID | Branch reference |
| student_id | UUID | Student reference |
| invoice_number | VARCHAR(50) | Unique invoice number |
| invoice_date | DATE | Invoice date |
| due_date | DATE | Payment due date |
| period_from | DATE | Billing period start |
| period_to | DATE | Billing period end |
| subtotal | DECIMAL | Total before adjustments |
| discount_amount | DECIMAL | Discounts applied |
| late_fee | DECIMAL | Late fee if any |
| adjustment | DECIMAL | Other adjustments |
| total_amount | DECIMAL | Final amount |
| paid_amount | DECIMAL | Amount paid |
| balance | DECIMAL | Outstanding balance |
| status | ENUM | draft, sent, partial, paid, overdue, cancelled |
| sent_at | TIMESTAMP | When sent to parent |
| paid_at | TIMESTAMP | When fully paid |

### 4.2 Invoice Items

**Entity: InvoiceItem**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| invoice_id | UUID | Parent invoice |
| fee_category_id | UUID | Fee category |
| description | VARCHAR(200) | Item description |
| quantity | INT | Quantity (usually 1) |
| unit_price | DECIMAL | Unit price |
| discount | DECIMAL | Item discount |
| tax | DECIMAL | Tax if applicable |
| total | DECIMAL | Line total |

### 4.3 Invoice Number Format

```
INV-{BRANCH_CODE}-{YEAR}{MONTH}-{SEQUENCE}
Example: INV-MAIN-202601-00001
```

---

## 5. Payment Collection

### 5.1 Payment Entity

**Entity: Payment**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| branch_id | UUID | Branch reference |
| student_id | UUID | Student reference |
| invoice_id | UUID | Invoice reference (optional) |
| receipt_number | VARCHAR(50) | Unique receipt number |
| payment_date | DATE | Payment date |
| amount | DECIMAL | Payment amount |
| payment_mode | ENUM | cash, cheque, dd, upi, card, netbanking, wallet |
| payment_reference | VARCHAR(100) | Transaction reference |
| payment_status | ENUM | pending, completed, failed, refunded |
| collected_by | UUID | Staff who collected |
| remarks | TEXT | Payment remarks |
| created_at | TIMESTAMP | Record creation |

### 5.2 Payment Modes

| Mode | Description | Fields Required |
|------|-------------|-----------------|
| cash | Cash payment | None |
| cheque | Cheque payment | cheque_number, bank_name, cheque_date |
| dd | Demand Draft | dd_number, bank_name, dd_date |
| upi | UPI payment | upi_id, transaction_id |
| card | Card payment | card_last4, transaction_id |
| netbanking | Net banking | bank_name, transaction_id |
| wallet | Digital wallet | wallet_name, transaction_id |

### 5.3 Cheque/DD Details

**Entity: PaymentCheque**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| payment_id | UUID | Payment reference |
| cheque_number | VARCHAR(20) | Cheque/DD number |
| bank_name | VARCHAR(100) | Bank name |
| branch_name | VARCHAR(100) | Branch name |
| cheque_date | DATE | Cheque date |
| status | ENUM | pending, cleared, bounced |
| clearance_date | DATE | When cleared |
| bounce_reason | TEXT | If bounced |
| bounce_charges | DECIMAL | Bounce charges |

### 5.4 Online Payment Integration

**Supported Gateways**:
- Razorpay
- PayU
- CCAvenue
- Paytm

**Entity: OnlinePayment**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| payment_id | UUID | Payment reference |
| gateway | VARCHAR(50) | Payment gateway |
| gateway_order_id | VARCHAR(100) | Gateway order ID |
| gateway_payment_id | VARCHAR(100) | Gateway payment ID |
| gateway_signature | VARCHAR(200) | Verification signature |
| status | VARCHAR(50) | Gateway status |
| response_data | JSONB | Full gateway response |
| created_at | TIMESTAMP | Initiated at |
| completed_at | TIMESTAMP | Completed at |

---

## 6. Receipt Generation

### 6.1 Receipt Format

```
┌─────────────────────────────────────────────────────────────┐
│                    ABC PUBLIC SCHOOL                         │
│              123 Main Street, Mumbai 400001                  │
│         Phone: 022-12345678 | Email: info@abc.edu           │
├─────────────────────────────────────────────────────────────┤
│                     FEE RECEIPT                              │
│                                                              │
│ Receipt No: RCP-MAIN-202601-00001    Date: 22-Jan-2026      │
├─────────────────────────────────────────────────────────────┤
│ Student Name: Aarav Sharma           Adm No: 2024001        │
│ Class: 5-A                           Roll No: 15            │
│ Father's Name: Rajesh Sharma                                │
├─────────────────────────────────────────────────────────────┤
│ Particulars                                      Amount (₹) │
│ ─────────────────────────────────────────────────────────── │
│ Tuition Fee (Jan-Mar 2026)                         15,000   │
│ Transport Fee (Jan 2026)                            2,000   │
│ Lab Fee (Annual)                                    5,000   │
│ ─────────────────────────────────────────────────────────── │
│ Subtotal                                           22,000   │
│ Sibling Discount (10%)                             -2,200   │
│ ─────────────────────────────────────────────────────────── │
│ Total Amount                                       19,800   │
│ Amount Paid                                        19,800   │
│ Balance Due                                             0   │
├─────────────────────────────────────────────────────────────┤
│ Payment Mode: UPI | Ref: 123456789012                       │
│ Received By: Ms. Priya (Accountant)                         │
├─────────────────────────────────────────────────────────────┤
│           *** Thank You for the Payment ***                 │
│        This is a computer generated receipt                 │
└─────────────────────────────────────────────────────────────┘
```

---

## 7. Late Fee Management

### 7.1 Late Fee Rules

**Entity: LateFeeRule**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| name | VARCHAR(100) | Rule name |
| grace_days | INT | Days after due date |
| calculation_type | ENUM | fixed, percentage, per_day |
| value | DECIMAL | Fee value |
| max_amount | DECIMAL | Maximum late fee |
| is_compounding | BOOLEAN | Compound monthly |
| applicable_categories | JSONB | Fee categories |
| is_active | BOOLEAN | Active status |

**Example Rules**:
```yaml
late_fee_rules:
  - name: Standard Late Fee
    grace_days: 7
    calculation_type: fixed
    value: 100
    max_amount: 500

  - name: Per Day Penalty
    grace_days: 15
    calculation_type: per_day
    value: 10
    max_amount: 300
```

---

## 8. Dues & Defaulters

### 8.1 Dues Tracking

**Entity: StudentDue**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| student_id | UUID | Student reference |
| invoice_id | UUID | Invoice reference |
| due_amount | DECIMAL | Amount due |
| due_date | DATE | Original due date |
| days_overdue | INT | Days past due |
| late_fee_applied | DECIMAL | Late fee charged |
| last_reminder_at | TIMESTAMP | Last reminder sent |
| reminder_count | INT | Number of reminders |
| status | ENUM | pending, partial, paid, waived |

### 8.2 Defaulter Report

```
Defaulter Report | As of: 22-Jan-2026 | Branch: Main Campus

| # | Adm No | Name          | Class | Total Due | Overdue Days | Last Paid |
|---|--------|---------------|-------|-----------|--------------|-----------|
| 1 | 2024015| Rahul Verma   | 5-A   | ₹25,000   | 45 days      | 15-Nov    |
| 2 | 2024023| Priya Singh   | 6-B   | ₹18,500   | 30 days      | 01-Dec    |
| 3 | 2024031| Amit Kumar    | 7-A   | ₹32,000   | 60 days      | 01-Nov    |

Total Dues: ₹75,500 | Total Defaulters: 3
```

---

## 9. Refund Management

### 9.1 Refund Request

**Entity: Refund**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| student_id | UUID | Student reference |
| request_date | DATE | Request date |
| reason | TEXT | Refund reason |
| requested_amount | DECIMAL | Amount requested |
| approved_amount | DECIMAL | Amount approved |
| status | ENUM | pending, approved, rejected, processed |
| approved_by | UUID | Approving authority |
| approved_at | TIMESTAMP | Approval time |
| processed_at | TIMESTAMP | When refund issued |
| payment_mode | VARCHAR(50) | Refund mode |
| reference_number | VARCHAR(100) | Refund reference |
| remarks | TEXT | Admin remarks |

---

## 10. Reports

### 10.1 Collection Report

```
Daily Collection Report | Date: 22-Jan-2026

| Category      | Cash     | Cheque   | UPI      | Card     | Total    |
|---------------|----------|----------|----------|----------|----------|
| Tuition Fee   | ₹15,000  | ₹30,000  | ₹45,000  | ₹10,000  | ₹100,000 |
| Transport Fee | ₹5,000   | ₹0       | ₹12,000  | ₹3,000   | ₹20,000  |
| Others        | ₹2,000   | ₹5,000   | ₹8,000   | ₹0       | ₹15,000  |
|---------------|----------|----------|----------|----------|----------|
| Total         | ₹22,000  | ₹35,000  | ₹65,000  | ₹13,000  | ₹135,000 |
```

### 10.2 Outstanding Report

```
Class-wise Outstanding | As of: 22-Jan-2026

| Class | Students | Total Due    | 0-30 Days | 31-60 Days | 60+ Days |
|-------|----------|--------------|-----------|------------|----------|
| 1     | 12       | ₹1,20,000    | ₹80,000   | ₹30,000    | ₹10,000  |
| 2     | 8        | ₹95,000      | ₹60,000   | ₹25,000    | ₹10,000  |
| 3     | 15       | ₹1,80,000    | ₹1,00,000 | ₹50,000    | ₹30,000  |
...
| Total | 85       | ₹12,50,000   | ₹8,00,000 | ₹3,00,000  | ₹1,50,000|
```

---

## 11. API Endpoints

```
# Fee Structure
GET    /api/v1/fee-structures               # List structures
POST   /api/v1/fee-structures               # Create structure
GET    /api/v1/fee-structures/{id}          # Get structure
PUT    /api/v1/fee-structures/{id}          # Update structure

# Student Fees
GET    /api/v1/students/{id}/fees           # Student fee details
POST   /api/v1/students/{id}/fees           # Assign fee structure
GET    /api/v1/students/{id}/invoices       # Student invoices
GET    /api/v1/students/{id}/payments       # Payment history

# Invoices
POST   /api/v1/invoices/generate            # Generate invoices
GET    /api/v1/invoices/{id}                # Get invoice
POST   /api/v1/invoices/{id}/send           # Send to parent

# Payments
POST   /api/v1/payments                     # Record payment
GET    /api/v1/payments/{id}/receipt        # Get receipt PDF
POST   /api/v1/payments/online/initiate     # Start online payment
POST   /api/v1/payments/online/verify       # Verify payment

# Reports
GET    /api/v1/fees/reports/collection      # Collection report
GET    /api/v1/fees/reports/outstanding     # Outstanding report
GET    /api/v1/fees/reports/defaulters      # Defaulter list
```

---

## 12. Business Rules

| Rule | Description |
|------|-------------|
| Invoice Sequence | Invoice numbers must be sequential per branch per month |
| Receipt Mandatory | Every payment must generate a receipt |
| Edit Restriction | Cannot edit payments older than 24 hours |
| Refund Approval | Refunds above ₹5000 need principal approval |
| Cheque Clearance | Mark cheques as cleared only after bank confirmation |
| Discount Limit | Maximum discount cannot exceed 50% without special approval |
| Partial Payment | Allow partial payments, update invoice balance |

---

## 13. Related Documents

- [03-student-management.md](./03-student-management.md) - Student profiles
- [14-parent-portal.md](./14-parent-portal.md) - Parent payment access
- [21-analytics-dashboards.md](./21-analytics-dashboards.md) - Financial dashboards
- [index.md](./index.md) - Main PRD index

---

**Previous**: [11-leave-management.md](./11-leave-management.md)
**Next**: [13-communication-system.md](./13-communication-system.md)
