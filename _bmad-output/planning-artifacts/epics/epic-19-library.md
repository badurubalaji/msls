# Epic 19: Library Management

**Phase:** 2 (Extended)
**Priority:** Medium - Academic support feature

## Epic Goal

Enable complete library operations from cataloging to circulation.

## User Value

Librarians can manage books, students/staff can borrow and return, and fines are automatically calculated.

## FRs Covered

FR-LB-01 to FR-LB-06

---

## Stories

### Story 19.1: Book Cataloging

As a **librarian**,
I want **to catalog books**,
So that **library inventory is managed**.

**Acceptance Criteria:**

**Given** librarian is adding book
**When** entering details
**Then** they can enter: ISBN, title, authors
**And** they can enter: publisher, year, edition
**And** they can select: category, subject
**And** they can enter: shelf location
**And** they can upload: cover image

**Given** book is cataloged
**When** adding copies
**Then** they can add: multiple copies
**And** each copy gets: accession number, barcode
**And** copy condition can be set
**And** total copies updated

---

### Story 19.2: Book Copy Management

As a **librarian**,
I want **to manage book copies**,
So that **each copy is tracked**.

**Acceptance Criteria:**

**Given** book exists
**When** viewing copies
**Then** they see: list of copies with status
**And** status: available, issued, reserved, lost

**Given** copy condition changes
**When** updating
**Then** they can update: condition (new, good, fair, poor)
**And** they can mark: damaged, discarded
**And** history is maintained

---

### Story 19.3: Library Member Management

As a **librarian**,
I want **to manage library members**,
So that **borrowing limits are set**.

**Acceptance Criteria:**

**Given** student/staff are members
**When** viewing member
**Then** they see: member ID, type
**And** they see: borrowing limits
**And** they see: current borrows, history
**And** they see: fines pending

**Given** member limits
**When** configured
**Then** different limits per member type
**And** max books, loan period set
**And** renewal limits set

---

### Story 19.4: Book Issue

As a **librarian**,
I want **to issue books to members**,
So that **borrowing is recorded**.

**Acceptance Criteria:**

**Given** librarian is issuing book
**When** processing
**Then** they scan/search: member card
**And** they scan/search: book barcode
**And** system checks: borrowing limit
**And** system sets: due date

**Given** issue is confirmed
**When** saved
**Then** book copy status: issued
**And** transaction recorded
**And** receipt can be printed
**And** available copies decremented

---

### Story 19.5: Book Return

As a **librarian**,
I want **to process book returns**,
So that **books are available again**.

**Acceptance Criteria:**

**Given** librarian is processing return
**When** scanning book
**Then** borrower info shown
**And** due date shown
**And** fine calculated if overdue
**And** condition can be updated

**Given** return is confirmed
**When** saved
**Then** book copy status: available
**And** transaction updated
**And** fine added to member (if any)
**And** available copies incremented

---

### Story 19.6: Book Reservation

As a **member**,
I want **to reserve books**,
So that **I get them when available**.

**Acceptance Criteria:**

**Given** book is not available
**When** reserving
**Then** member added to waitlist
**And** queue position shown
**And** notification preference set

**Given** book is returned
**When** reservation exists
**Then** first in queue notified
**And** book held for 3 days
**And** if not collected, next in queue

---

### Story 19.7: Fine Management

As a **librarian**,
I want **to manage fines**,
So that **overdue penalties are enforced**.

**Acceptance Criteria:**

**Given** fine policy exists
**When** configured
**Then** fine per day is set
**And** maximum fine cap set
**And** grace period set

**Given** book is overdue
**When** returned
**Then** fine auto-calculated
**And** displayed to librarian
**And** can be collected or waived
**And** collected fines recorded

---

### Story 19.8: Circulation Reports

As a **librarian**,
I want **to generate circulation reports**,
So that **library usage is analyzed**.

**Acceptance Criteria:**

**Given** librarian is on reports
**When** generating
**Then** they can see: books issued today/month
**And** they can see: books returned
**And** they can see: overdue books
**And** they can see: most borrowed books

**Given** export is needed
**When** exporting
**Then** reports can be exported to Excel
**And** filters can be applied
**And** date ranges selected
