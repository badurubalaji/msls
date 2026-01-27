# Epic 20: Inventory & Assets

**Phase:** 2 (Extended)
**Priority:** Medium - Operational feature

## Epic Goal

Enable tracking of school assets and consumable inventory.

## User Value

Admins can track all assets, request maintenance, manage stock levels, and process purchase orders.

## FRs Covered

FR-IN-01 to FR-IN-06

---

## Stories

### Story 20.1: Asset Registration

As an **inventory manager**,
I want **to register school assets**,
So that **all assets are tracked**.

**Acceptance Criteria:**

**Given** manager is adding asset
**When** entering details
**Then** they can enter: name, category, make/model
**And** they can enter: serial number, purchase date, price
**And** they can enter: warranty details, vendor
**And** they can set: location, assigned staff
**And** asset code auto-generated

**Given** asset is registered
**When** saved
**Then** QR/barcode generated
**And** asset appears in register
**And** depreciation starts (if configured)

---

### Story 20.2: Asset Location Management

As an **inventory manager**,
I want **to manage asset locations**,
So that **assets can be found**.

**Acceptance Criteria:**

**Given** locations are configured
**When** setting up
**Then** they can create: buildings, floors, rooms
**And** hierarchical structure supported
**And** in-charge can be assigned

**Given** asset is transferred
**When** moving
**Then** new location selected
**And** transfer recorded with date
**And** history maintained

---

### Story 20.3: Asset Depreciation

As an **inventory manager**,
I want **depreciation to be calculated**,
So that **asset value is tracked**.

**Acceptance Criteria:**

**Given** depreciation method configured
**When** calculated
**Then** straight-line or declining balance applied
**And** current value calculated
**And** monthly/yearly depreciation shown
**And** book value updated

**Given** asset reports
**When** viewing
**Then** they see: total asset value
**And** they see: depreciation this year
**And** they see: net current value

---

### Story 20.4: Maintenance Request

As a **staff member**,
I want **to request asset maintenance**,
So that **issues are reported**.

**Acceptance Criteria:**

**Given** staff finds asset issue
**When** raising request
**Then** they can select: asset
**And** they can describe: issue
**And** they can set: priority
**And** they can upload: photos

**Given** request is submitted
**When** saved
**Then** request number generated
**And** notification to inventory team
**And** status: open

---

### Story 20.5: Maintenance Processing

As an **inventory manager**,
I want **to process maintenance requests**,
So that **assets are repaired**.

**Acceptance Criteria:**

**Given** request exists
**When** processing
**Then** they can assign: technician
**And** they can update: status (in progress)
**And** they can record: resolution, cost
**And** they can close: request

**Given** maintenance is complete
**When** closing
**Then** asset condition updated
**And** cost recorded
**And** history maintained

---

### Story 20.6: Stock/Consumable Management

As an **inventory manager**,
I want **to manage consumable stock**,
So that **supplies don't run out**.

**Acceptance Criteria:**

**Given** stock item is created
**When** setting up
**Then** they can enter: name, category, SKU
**And** they can set: unit of measure
**And** they can set: reorder level, reorder quantity
**And** they can set: storage location

**Given** stock is received
**When** recording
**Then** quantity added to current stock
**And** receipt reference recorded
**And** stock movement logged

---

### Story 20.7: Purchase Order Management

As an **inventory manager**,
I want **to create purchase orders**,
So that **procurement is tracked**.

**Acceptance Criteria:**

**Given** purchase is needed
**When** creating PO
**Then** they can select: vendor
**And** they can add: line items (quantity, price)
**And** they can set: expected date
**And** total calculated
**And** PO number generated

**Given** PO approval workflow
**When** submitted
**Then** goes for approval
**And** approver can approve/reject
**And** on approval: sent to vendor

---

### Story 20.8: Vendor Management

As an **inventory manager**,
I want **to manage vendors**,
So that **suppliers are tracked**.

**Acceptance Criteria:**

**Given** adding vendor
**When** entering details
**Then** they can enter: name, contact, address
**And** they can enter: GST, PAN numbers
**And** they can assign: categories supplied
**And** they can rate: vendor performance

**Given** viewing vendors
**When** browsing
**Then** they see: vendor list with status
**And** they can filter: by category
**And** they can see: purchase history
