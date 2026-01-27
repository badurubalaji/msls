# 18 - Library Management

**Parent Document**: [index.md](./index.md)
**Status**: Draft
**Phase**: 2

---

## 1. Overview

The Library Management module handles book cataloging, issue/return, member management, fines, and digital library resources.

---

## 2. Book Catalog

### 2.1 Book Entity

**Entity: Book**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| branch_id | UUID | Branch/Library |
| isbn | VARCHAR(20) | ISBN number |
| title | VARCHAR(500) | Book title |
| subtitle | VARCHAR(500) | Subtitle |
| authors | JSONB | Author names |
| publisher | VARCHAR(200) | Publisher |
| publication_year | INT | Year published |
| edition | VARCHAR(50) | Edition |
| category_id | UUID | Category |
| subject_id | UUID | Subject mapping |
| language | VARCHAR(50) | Language |
| pages | INT | Page count |
| description | TEXT | Synopsis |
| cover_image | VARCHAR(500) | Cover image URL |
| shelf_location | VARCHAR(50) | Physical location |
| total_copies | INT | Total copies |
| available_copies | INT | Currently available |
| is_reference | BOOLEAN | Reference only |
| is_digital | BOOLEAN | E-book available |
| digital_url | VARCHAR(500) | E-book link |
| barcode | VARCHAR(50) | Book barcode |
| status | ENUM | active, archived, lost |

### 2.2 Book Categories

```
Library Categories:
├── Fiction
│   ├── Classic Literature
│   ├── Modern Fiction
│   ├── Science Fiction
│   └── Fantasy
├── Non-Fiction
│   ├── Biography
│   ├── History
│   └── Science
├── Academic
│   ├── Textbooks (Class 1-12)
│   ├── Reference Books
│   └── Competitive Exams
├── Magazines & Periodicals
└── Children's Section
```

### 2.3 Book Copy

**Entity: BookCopy**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| book_id | UUID | Parent book |
| accession_number | VARCHAR(50) | Unique copy ID |
| barcode | VARCHAR(50) | Copy barcode |
| condition | ENUM | new, good, fair, poor, damaged |
| acquisition_date | DATE | When acquired |
| acquisition_type | ENUM | purchase, donation, exchange |
| price | DECIMAL | Cost |
| vendor | VARCHAR(200) | Supplier |
| status | ENUM | available, issued, reserved, lost, discarded |

---

## 3. Member Management

### 3.1 Library Member

**Entity: LibraryMember**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| user_id | UUID | User reference |
| member_type | ENUM | student, staff, external |
| membership_number | VARCHAR(50) | Member ID |
| valid_from | DATE | Membership start |
| valid_until | DATE | Membership end |
| max_books | INT | Borrowing limit |
| max_days | INT | Default loan period |
| fine_exempted | BOOLEAN | Fine exemption |
| status | ENUM | active, suspended, expired |

### 3.2 Borrowing Limits

| Member Type | Max Books | Loan Period | Renewal |
|-------------|-----------|-------------|---------|
| Student (1-5) | 2 | 7 days | 1 |
| Student (6-10) | 3 | 14 days | 1 |
| Student (11-12) | 4 | 21 days | 2 |
| Teacher | 5 | 30 days | 2 |
| Staff | 3 | 14 days | 1 |

---

## 4. Issue & Return

### 4.1 Transaction Entity

**Entity: LibraryTransaction**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| member_id | UUID | Member reference |
| book_copy_id | UUID | Copy reference |
| transaction_type | ENUM | issue, return, renew, reserve |
| issue_date | DATE | Issue date |
| due_date | DATE | Return due |
| return_date | DATE | Actual return |
| renewed_count | INT | Times renewed |
| fine_amount | DECIMAL | Fine charged |
| fine_paid | BOOLEAN | Fine payment status |
| issued_by | UUID | Librarian |
| returned_to | UUID | Receiving librarian |
| remarks | TEXT | Notes |

### 4.2 Issue Interface

```
┌─────────────────────────────────────────────────────────────┐
│  ISSUE BOOK                                                  │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Member: [Scan Card / Search]                               │
│  ┌─────────────────────────────────────────────────────┐    │
│  │ Aarav Sharma | Class 10-A | Member: LIB-2024-0542  │    │
│  │ Books: 2/3 issued | No fines pending               │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                              │
│  Book: [Scan Barcode / Search]                              │
│  ┌─────────────────────────────────────────────────────┐    │
│  │ Harry Potter and the Chamber of Secrets             │    │
│  │ Accession: ACC-2024-1234 | Available: Yes          │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                              │
│  Issue Date: 22-Jan-2026                                    │
│  Due Date: 05-Feb-2026 (14 days)                           │
│                                                              │
│  [Issue Book]                                                │
└─────────────────────────────────────────────────────────────┘
```

---

## 5. Reservations

### 5.1 Reservation Entity

**Entity: BookReservation**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| member_id | UUID | Member reference |
| book_id | UUID | Book (not copy) |
| reserved_date | DATE | Reservation date |
| queue_position | INT | Position in queue |
| status | ENUM | waiting, notified, fulfilled, cancelled, expired |
| notified_at | TIMESTAMP | When notified |
| expiry_date | DATE | Hold expiry |

---

## 6. Fine Management

### 6.1 Fine Configuration

**Entity: FinePolicy**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| member_type | ENUM | Applicable member type |
| fine_per_day | DECIMAL | Daily fine |
| max_fine | DECIMAL | Maximum fine |
| grace_days | INT | Grace period |

### 6.2 Fine Calculation

```
Fine Calculation | Aarav Sharma

Book: Harry Potter (Chamber of Secrets)
Due Date: 05-Feb-2026
Returned: 12-Feb-2026
Overdue: 7 days

Fine Rate: ₹2/day
Grace Period: 2 days
Chargeable Days: 5

Fine Amount: ₹10

[Pay Fine] [Waive Fine]
```

---

## 7. Digital Library

### 7.1 E-Resources

**Entity: DigitalResource**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| title | VARCHAR(500) | Resource title |
| type | ENUM | ebook, journal, video, audio |
| category_id | UUID | Category |
| file_url | VARCHAR(500) | File/Link |
| access_type | ENUM | open, restricted, subscription |
| allowed_downloads | INT | Download limit |
| valid_until | DATE | Access expiry |

---

## 8. Reports

### 8.1 Circulation Report

```
Library Circulation | January 2026

| Metric              | Count |
|---------------------|-------|
| Books Issued        | 450   |
| Books Returned      | 420   |
| Currently Issued    | 180   |
| Overdue Books       | 25    |
| Reservations        | 15    |
| Fines Collected     | ₹1,250|

Most Borrowed:
1. Harry Potter Series (45 issues)
2. NCERT Science Class 10 (38 issues)
3. Geronimo Stilton (32 issues)
```

---

## 9. API Endpoints

```
# Catalog
GET    /api/v1/library/books                # List books
POST   /api/v1/library/books                # Add book
GET    /api/v1/library/books/{id}           # Get book
GET    /api/v1/library/books/search         # Search

# Transactions
POST   /api/v1/library/issue                # Issue book
POST   /api/v1/library/return               # Return book
POST   /api/v1/library/renew                # Renew book
POST   /api/v1/library/reserve              # Reserve book

# Members
GET    /api/v1/library/members/{id}         # Member details
GET    /api/v1/library/members/{id}/history # Borrowing history

# Fines
GET    /api/v1/library/fines                # Pending fines
POST   /api/v1/library/fines/{id}/pay       # Pay fine
```

---

## 10. Business Rules

| Rule | Description |
|------|-------------|
| Borrowing Limit | Cannot exceed max books |
| No Issue with Fine | Block issue if fine > ₹50 |
| Reference Books | Cannot be issued |
| Renewal Limit | Max 2 renewals allowed |
| Reservation Hold | 3 days to collect reserved book |
| Lost Book | Pay replacement cost + 20% |

---

## 11. Related Documents

- [03-student-management.md](./03-student-management.md) - Student members
- [10-staff-management.md](./10-staff-management.md) - Staff members
- [index.md](./index.md) - Main PRD index

---

**Previous**: [17-transport-management.md](./17-transport-management.md)
**Next**: [19-inventory-assets.md](./19-inventory-assets.md)
