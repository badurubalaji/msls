# Epic 11: Fees & Payments

**Phase:** 1 (MVP)
**Priority:** High - Critical business operation

## Epic Goal

Enable complete fee management from structure definition to payment collection and receipt generation.

## User Value

Admins can configure fees, parents can pay online, and accounts can track dues and generate reports.

## FRs Covered

FR-FE-01 to FR-FE-09

---

## Stories

### Story 11.1: Fee Category & Structure Configuration

As an **accounts administrator**,
I want **to configure fee categories and structures**,
So that **fees are standardized across classes**.

**Acceptance Criteria:**

**Given** admin is on fee settings
**When** creating a fee category
**Then** they can enter: name (Tuition, Transport, Books, etc.)
**And** they can set: frequency (monthly, quarterly, annual, one-time)
**And** they can set: tax applicable (yes/no, rate)
**And** they can set: refundable (yes/no)

**Given** fee structure is created
**When** defining amounts
**Then** they can set: base amount per category
**And** they can assign: to specific class or all classes
**And** they can set: due date pattern (e.g., 10th of each month)
**And** structure is linked to academic year

---

### Story 11.2: Student Fee Assignment

As an **accounts administrator**,
I want **to assign fee structures to students**,
So that **each student has correct fees**.

**Acceptance Criteria:**

**Given** a student is enrolled
**When** assigning fees
**Then** default class fee structure is auto-assigned
**And** individual adjustments can be made
**And** additional fees (transport, books) can be added
**And** total annual fee is calculated

**Given** fee modifications are needed
**When** updating student fees
**Then** they can add: one-time fee (admission fee, uniform)
**And** they can remove: categories not applicable
**And** changes are logged with reason
**And** parent is notified of fee changes

---

### Story 11.3: Discount & Concession Management

As an **accounts administrator**,
I want **to apply discounts and concessions**,
So that **eligible students pay reduced fees**.

**Acceptance Criteria:**

**Given** discount types are configured
**When** creating discount
**Then** they can set: name (Sibling, Staff Child, Merit, Need-based)
**And** they can set: percentage or fixed amount
**And** they can set: applicable categories
**And** they can set: requires approval

**Given** a discount is applied to student
**When** assigning
**Then** they can select: discount type
**And** they can enter: supporting documents/reason
**And** discount reflects in fee calculation
**And** approval workflow if configured

---

### Story 11.4: Invoice Generation

As an **accounts administrator**,
I want **to generate fee invoices**,
So that **parents receive proper billing**.

**Acceptance Criteria:**

**Given** fee structure is assigned
**When** generating invoices
**Then** they can generate for: month, quarter, or full year
**And** they can generate for: individual or bulk (class/all)
**And** invoice number is auto-generated (INV/2026/0001)
**And** invoice includes: student details, fee breakdown, total

**Given** invoice is generated
**When** viewing invoice
**Then** it shows: fee categories, amounts, tax, discount, net amount
**And** due date is shown
**And** payment instructions are included
**And** invoice can be printed/emailed

---

### Story 11.5: Payment Collection - Counter

As an **accounts staff**,
I want **to collect payments at school counter**,
So that **parents can pay in person**.

**Acceptance Criteria:**

**Given** parent is at counter
**When** collecting payment
**Then** staff can search: by student name or ID
**And** staff sees: pending dues for the student
**And** staff can select: invoices to pay (partial allowed)
**And** staff can enter: payment mode (cash, cheque, card, UPI)

**Given** payment is entered
**When** confirming
**Then** payment amount is validated
**And** if cheque: cheque number and date recorded
**And** if UPI: transaction ID recorded
**And** receipt is generated immediately

---

### Story 11.6: Online Payment Integration

As a **parent**,
I want **to pay fees online**,
So that **I can pay from anywhere anytime**.

**Acceptance Criteria:**

**Given** a parent is on fee payment page
**When** selecting invoices to pay
**Then** they see: pending invoices with amounts
**And** they can select: multiple invoices to pay together
**And** total payable is calculated
**And** convenience fee (if any) is shown

**Given** payment is initiated
**When** redirected to payment gateway
**Then** Razorpay/PayU payment page opens
**And** parent can choose: UPI, card, netbanking
**And** payment is processed securely
**And** on success: redirected back with confirmation

**Given** payment succeeds
**When** callback is received
**Then** payment is recorded automatically
**And** receipt is generated
**And** email/SMS confirmation sent
**And** dues are updated

---

### Story 11.7: Receipt Generation & Management

As an **accounts staff or parent**,
I want **to generate and access receipts**,
So that **payment proof is available**.

**Acceptance Criteria:**

**Given** payment is recorded
**When** generating receipt
**Then** receipt includes: receipt number, date, student details
**And** receipt shows: payment breakdown, mode, amount
**And** receipt has: school branding, authorized signature
**And** receipt is downloadable as PDF

**Given** receipts need management
**When** viewing receipt list
**Then** they can filter: by date, student, payment mode
**And** they can search: by receipt number
**And** they can reprint: any receipt
**And** duplicate receipt is marked as "DUPLICATE"

---

### Story 11.8: Late Fee Calculation

As an **accounts administrator**,
I want **late fees to be calculated automatically**,
So that **payment deadlines are enforced**.

**Acceptance Criteria:**

**Given** late fee rules are configured
**When** setting rules
**Then** they can set: grace period (days after due date)
**And** they can set: late fee per day or fixed amount
**And** they can set: maximum late fee cap
**And** they can set: exempt categories (staff children, etc.)

**Given** due date passes
**When** calculating fees
**Then** late fee is auto-added to pending amount
**And** calculation is per invoice
**And** late fee appears as separate line item
**And** notification sent about late fee

---

### Story 11.9: Dues Tracking & Defaulter Reports

As an **accounts administrator**,
I want **to track dues and identify defaulters**,
So that **follow-up can be done**.

**Acceptance Criteria:**

**Given** dues exist
**When** viewing dues dashboard
**Then** they see: total dues amount, student count
**And** they see: aging analysis (0-30, 31-60, 61-90, 90+ days)
**And** they can drill down to individual students

**Given** defaulter report is needed
**When** generating report
**Then** they can filter: by minimum due amount, days overdue
**And** report shows: student, class, total due, days overdue
**And** export to Excel available
**And** bulk SMS to defaulters available

---

### Story 11.10: Fee Refund Processing

As an **accounts administrator**,
I want **to process fee refunds**,
So that **withdrawing students get appropriate refunds**.

**Acceptance Criteria:**

**Given** a refund is needed
**When** initiating refund request
**Then** they can select: student, payments to refund
**And** they can enter: refund amount (within paid amount)
**And** they can enter: reason and supporting details
**And** request goes for approval

**Given** refund is approved
**When** processing
**Then** refund mode is selected (original mode or cheque)
**And** refund is recorded with reference
**And** fee ledger is updated
**And** confirmation sent to parent
