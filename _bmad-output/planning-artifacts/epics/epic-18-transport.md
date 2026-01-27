# Epic 18: Transport Management

**Phase:** 2 (Extended)
**Priority:** Medium - Operational feature

## Epic Goal

Enable complete school transport operations with GPS tracking.

## User Value

Admins manage fleet and routes, parents track buses in real-time, and attendants mark transport attendance.

## FRs Covered

FR-TR-01 to FR-TR-07

---

## Stories

### Story 18.1: Vehicle Management

As a **transport administrator**,
I want **to manage school vehicles**,
So that **fleet is properly tracked**.

**Acceptance Criteria:**

**Given** admin is on vehicle management
**When** adding vehicle
**Then** they can enter: registration number, type (bus, van)
**And** they can enter: capacity, make, model
**And** they can enter: insurance details, expiry date
**And** they can enter: fitness certificate, permit details
**And** they can add: GPS device ID

**Given** vehicle documents expire
**When** approaching expiry
**Then** alert shown on dashboard
**And** notification sent 30 days before
**And** vehicle can be marked for renewal

---

### Story 18.2: Route Configuration

As a **transport administrator**,
I want **to configure routes with stops**,
So that **pickup/drop is organized**.

**Acceptance Criteria:**

**Given** admin is creating route
**When** configuring
**Then** they can enter: route name, code
**And** they can assign: vehicle, driver, attendant
**And** they can set: start time, estimated duration

**Given** adding stops
**When** configuring route
**Then** they can add: stop name, GPS coordinates
**And** they can set: stop order
**And** they can set: pickup/drop times
**And** they can set: fee slab for stop

---

### Story 18.3: Student Route Assignment

As a **transport administrator**,
I want **to assign students to routes**,
So that **pickup is organized**.

**Acceptance Criteria:**

**Given** admin is assigning transport
**When** selecting student
**Then** they can select: route, stop
**And** they can set: transport type (both ways, morning, drop)
**And** they can enter: pickup address, emergency contact
**And** monthly fee is calculated based on stop

**Given** assignment is made
**When** saved
**Then** student appears in route roster
**And** fee is added to student dues
**And** parent is notified of assignment

---

### Story 18.4: Driver & Attendant Management

As a **transport administrator**,
I want **to manage drivers and attendants**,
So that **transport staff is tracked**.

**Acceptance Criteria:**

**Given** admin is adding transport staff
**When** entering details
**Then** they can enter: name, phone, address
**And** they can enter: license details (drivers)
**And** they can upload: photo, ID proof
**And** they can enter: police verification status

**Given** staff documents expire
**When** approaching expiry
**Then** license expiry alerts shown
**And** renewal reminders sent

---

### Story 18.5: GPS Real-Time Tracking

As a **transport administrator**,
I want **to track vehicles in real-time**,
So that **fleet is monitored**.

**Acceptance Criteria:**

**Given** GPS devices are configured
**When** viewing tracking dashboard
**Then** they see: map with vehicle locations
**And** they see: vehicle status (moving, stopped)
**And** they see: speed, direction
**And** they see: last updated time

**Given** route is in progress
**When** tracking
**Then** route path shown on map
**And** current position highlighted
**And** stops passed shown differently
**And** ETA to next stop calculated

---

### Story 18.6: Parent Bus Tracking

As a **parent**,
I want **to track my child's bus**,
So that **I know when bus will arrive**.

**Acceptance Criteria:**

**Given** parent opens tracking
**When** viewing
**Then** they see: bus location on map
**And** they see: ETA to their stop
**And** they see: stops remaining

**Given** bus approaches stop
**When** within 5 minutes
**Then** push notification sent
**And** "Arriving soon" message shown

---

### Story 18.7: Transport Attendance

As an **attendant**,
I want **to mark transport attendance**,
So that **boarding is recorded**.

**Acceptance Criteria:**

**Given** attendant opens app
**When** marking attendance
**Then** they see: students assigned to route
**And** they can mark: boarded, absent
**And** they can scan: student ID card
**And** boarding time is recorded

**Given** student doesn't board
**When** marked absent
**Then** parent is notified immediately
**And** reason can be added (optional)

---

### Story 18.8: Transport Fee Calculation

As a **system**,
I want **to calculate distance-based fees**,
So that **billing is accurate**.

**Acceptance Criteria:**

**Given** fee slabs are configured
**When** student is assigned
**Then** distance from school to stop calculated
**And** appropriate slab applied
**And** monthly fee determined
**And** fee added to student invoices

**Given** transport changes mid-month
**When** assignment updated
**Then** pro-rata calculation done
**And** adjustment made in next invoice
