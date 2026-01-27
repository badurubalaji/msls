# 17 - Transport Management

**Parent Document**: [index.md](./index.md)
**Status**: Draft
**Phase**: 2

---

## 1. Overview

The Transport Management module handles school bus routes, vehicle management, student assignments, driver/attendant management, and real-time GPS tracking.

---

## 2. Vehicle Management

### 2.1 Vehicle Entity

**Entity: Vehicle**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| branch_id | UUID | Branch reference |
| vehicle_number | VARCHAR(20) | Registration number |
| vehicle_type | ENUM | bus, van, auto, car |
| make | VARCHAR(50) | Manufacturer |
| model | VARCHAR(50) | Model |
| year | INT | Year of manufacture |
| capacity | INT | Seating capacity |
| fuel_type | ENUM | diesel, petrol, cng, electric |
| insurance_number | VARCHAR(50) | Insurance policy |
| insurance_expiry | DATE | Insurance expiry |
| fitness_expiry | DATE | Fitness certificate expiry |
| permit_expiry | DATE | Permit expiry |
| puc_expiry | DATE | PUC expiry |
| gps_device_id | VARCHAR(100) | GPS tracker ID |
| status | ENUM | active, maintenance, inactive |
| current_odometer | INT | Current reading |
| photos | JSONB | Vehicle photos |

### 2.2 Vehicle Dashboard

```
Vehicle Fleet | 12 Vehicles

| Vehicle   | Type | Capacity | Route   | Driver      | Status      |
|-----------|------|----------|---------|-------------|-------------|
| MH01-1234 | Bus  | 40       | Route 1 | Ramesh K.   | Active      |
| MH01-5678 | Bus  | 40       | Route 2 | Suresh P.   | Active      |
| MH01-9012 | Van  | 15       | Route 3 | Mahesh D.   | Active      |
| MH01-3456 | Bus  | 40       | -       | -           | Maintenance |

Alerts:
⚠️ MH01-1234: Insurance expires in 15 days
⚠️ MH01-5678: PUC expires in 7 days
```

---

## 3. Route Management

### 3.1 Route Entity

**Entity: Route**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| branch_id | UUID | Branch reference |
| name | VARCHAR(100) | Route name |
| code | VARCHAR(20) | Route code |
| description | TEXT | Route description |
| vehicle_id | UUID | Assigned vehicle |
| driver_id | UUID | Assigned driver |
| attendant_id | UUID | Assigned attendant |
| route_type | ENUM | morning, afternoon, both |
| start_time | TIME | Departure time |
| estimated_duration | INT | Duration in minutes |
| distance_km | DECIMAL | Total distance |
| base_fee | DECIMAL | Monthly fee |
| is_active | BOOLEAN | Active status |

### 3.2 Route Stops

**Entity: RouteStop**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| route_id | UUID | Route reference |
| stop_name | VARCHAR(200) | Stop name/landmark |
| stop_order | INT | Stop sequence |
| latitude | DECIMAL | GPS latitude |
| longitude | DECIMAL | GPS longitude |
| pickup_time | TIME | Morning pickup |
| drop_time | TIME | Afternoon drop |
| distance_from_school | DECIMAL | Distance in km |
| fee_slab_id | UUID | Fee based on distance |

### 3.3 Route Map View

```
Route 1: Andheri - School | MH01-1234 | 28 Students

Stop Order:
1. 07:00 - Andheri Station (5 students)
2. 07:10 - DN Nagar (4 students)
3. 07:20 - Versova (6 students)
4. 07:35 - Lokhandwala (8 students)
5. 07:45 - Oshiwara (5 students)
6. 08:00 - School (Arrival)

Total Distance: 12 km
Estimated Duration: 60 min
Current Location: [GPS Map]
```

---

## 4. Student Transport Assignment

### 4.1 Assignment Entity

**Entity: StudentTransport**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| student_id | UUID | Student reference |
| route_id | UUID | Route reference |
| stop_id | UUID | Pickup/drop stop |
| academic_year_id | UUID | Academic year |
| transport_type | ENUM | both_ways, morning_only, drop_only |
| pickup_address | TEXT | Home address |
| emergency_contact | VARCHAR(20) | Emergency phone |
| start_date | DATE | Service start |
| end_date | DATE | Service end |
| monthly_fee | DECIMAL | Monthly charge |
| status | ENUM | active, suspended, cancelled |

### 4.2 Assignment Interface

```
Assign Transport | Aarav Sharma | Class 10-A

Current: Not assigned

Available Routes:
| Route   | Stop          | Pickup | Fee/Month |
|---------|---------------|--------|-----------|
| Route 1 | Lokhandwala   | 07:35  | ₹2,000    |
| Route 1 | Versova       | 07:20  | ₹1,800    |
| Route 2 | Goregaon      | 07:30  | ₹2,200    |

Selected: Route 1 - Lokhandwala

Transport Type: ○ Both Ways ○ Morning Only ○ Drop Only

Pickup Address: [Home address]
Emergency Contact: [9876543210]

[Assign Transport]
```

---

## 5. Driver & Attendant Management

### 5.1 Driver Entity

**Entity: TransportStaff**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| staff_type | ENUM | driver, attendant |
| name | VARCHAR(200) | Full name |
| phone | VARCHAR(20) | Phone number |
| license_number | VARCHAR(50) | Driving license |
| license_expiry | DATE | License expiry |
| badge_number | VARCHAR(50) | ID badge |
| address | TEXT | Address |
| photo_url | VARCHAR(500) | Photo |
| verification_status | ENUM | pending, verified, rejected |
| police_verification_date | DATE | Verification date |
| blood_group | VARCHAR(5) | Blood group |
| emergency_contact | VARCHAR(20) | Emergency phone |
| status | ENUM | active, on_leave, terminated |

### 5.2 Driver Dashboard

```
Driver: Ramesh Kumar | MH01-1234 | Route 1

Today's Trip:
Morning: Started 07:00 | Status: Completed
Afternoon: Starts 14:30 | Status: Upcoming

Students: 28 assigned | 26 picked | 2 absent
Next Stop: School (08:00)

License Expiry: 15-Aug-2026 (Valid)
Police Verification: Verified (20-Mar-2025)
```

---

## 6. GPS Tracking

### 6.1 GPS Data Entity

**Entity: VehicleLocation**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| vehicle_id | UUID | Vehicle reference |
| latitude | DECIMAL | Current latitude |
| longitude | DECIMAL | Current longitude |
| speed_kmph | DECIMAL | Current speed |
| heading | INT | Direction (0-360) |
| timestamp | TIMESTAMP | Location time |
| ignition_on | BOOLEAN | Engine status |
| door_open | BOOLEAN | Door status |

### 6.2 Live Tracking Interface

```
┌─────────────────────────────────────────────────────────────┐
│  LIVE TRACKING | Route 1 - MH01-1234                        │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌─────────────────────────────────────────────────────┐    │
│  │                                                      │    │
│  │                    [MAP VIEW]                        │    │
│  │                                                      │    │
│  │     🚌 Current Location                             │    │
│  │     📍 Next Stop: Lokhandwala                       │    │
│  │     🏫 Destination: School                          │    │
│  │                                                      │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                              │
│  Status: En Route | Speed: 25 km/h                          │
│  ETA at School: 08:05 (5 min delay)                         │
│  Last Updated: 07:42:15                                     │
│                                                              │
│  Students on Board: 22/28                                   │
│  Next Pickup: Lokhandwala (8 students)                      │
└─────────────────────────────────────────────────────────────┘
```

### 6.3 Parent Tracking View

```
Track Aarav's Bus | Route 1

[Map showing bus location and route]

🚌 Bus Status: On Route
📍 Current: Near Versova
⏱️ ETA at your stop (Lokhandwala): 07:33
🏫 ETA at School: 08:00

Notifications:
✓ 07:00 - Bus started from depot
✓ 07:20 - Passed Versova
```

---

## 7. Attendance (Transport)

### 7.1 Transport Attendance

**Entity: TransportAttendance**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| student_id | UUID | Student reference |
| route_id | UUID | Route reference |
| date | DATE | Date |
| trip_type | ENUM | morning, afternoon |
| status | ENUM | boarded, absent, cancelled |
| boarded_at | TIMESTAMP | Boarding time |
| stop_id | UUID | Actual boarding stop |
| marked_by | UUID | Attendant |

### 7.2 Attendant App Features

- Mark student boarding
- Scan student ID card
- Mark absent (inform parent)
- Emergency SOS button
- Route deviation alert

---

## 8. Fees Integration

### 8.1 Distance-based Slabs

**Entity: TransportFeeSlab**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| name | VARCHAR(100) | Slab name |
| min_distance | DECIMAL | From km |
| max_distance | DECIMAL | To km |
| monthly_fee | DECIMAL | Fee amount |
| academic_year_id | UUID | Academic year |

```
Transport Fee Slabs | 2025-26

| Distance      | Monthly Fee |
|---------------|-------------|
| 0 - 3 km      | ₹1,200      |
| 3 - 6 km      | ₹1,500      |
| 6 - 10 km     | ₹1,800      |
| 10 - 15 km    | ₹2,200      |
| > 15 km       | ₹2,500      |
```

---

## 9. Reports

### 9.1 Daily Transport Report

```
Transport Report | 22-Jan-2026

Morning Trip:
| Route | Vehicle   | Students | Boarded | Absent | On Time |
|-------|-----------|----------|---------|--------|---------|
| R1    | MH01-1234 | 28       | 26      | 2      | Yes     |
| R2    | MH01-5678 | 35       | 33      | 2      | 5 min late |
| R3    | MH01-9012 | 15       | 15      | 0      | Yes     |

Total Students Using Transport: 78
Absentees: 4 (5.1%)
```

### 9.2 Vehicle Utilization

```
Vehicle Utilization | January 2026

| Vehicle   | Capacity | Assigned | Utilization | Revenue   |
|-----------|----------|----------|-------------|-----------|
| MH01-1234 | 40       | 28       | 70%         | ₹56,000   |
| MH01-5678 | 40       | 35       | 88%         | ₹70,000   |
| MH01-9012 | 15       | 15       | 100%        | ₹30,000   |
```

---

## 10. API Endpoints

```
# Vehicles
GET    /api/v1/transport/vehicles           # List vehicles
POST   /api/v1/transport/vehicles           # Add vehicle
GET    /api/v1/transport/vehicles/{id}      # Get vehicle
PUT    /api/v1/transport/vehicles/{id}      # Update vehicle

# Routes
GET    /api/v1/transport/routes             # List routes
POST   /api/v1/transport/routes             # Create route
GET    /api/v1/transport/routes/{id}        # Get route
GET    /api/v1/transport/routes/{id}/stops  # Get stops
GET    /api/v1/transport/routes/{id}/students # Students on route

# Student Assignment
POST   /api/v1/transport/assign             # Assign student
PUT    /api/v1/transport/assign/{id}        # Update assignment
DELETE /api/v1/transport/assign/{id}        # Remove assignment

# Tracking
GET    /api/v1/transport/vehicles/{id}/location # Current location
GET    /api/v1/transport/vehicles/{id}/history  # Location history
GET    /api/v1/transport/routes/{id}/track      # Track route

# Attendance
POST   /api/v1/transport/attendance         # Mark attendance
GET    /api/v1/transport/attendance/today   # Today's attendance
```

---

## 11. Business Rules

| Rule | Description |
|------|-------------|
| Capacity Check | Cannot assign more than vehicle capacity |
| License Validity | Alert before driver license expires |
| Insurance Check | Vehicle cannot operate with expired insurance |
| GPS Mandatory | All vehicles must have working GPS |
| Attendant Required | Every vehicle must have attendant |
| Speed Alert | Alert if speed exceeds 40 km/h |
| Route Deviation | Alert if vehicle deviates from route |

---

## 12. Related Documents

- [03-student-management.md](./03-student-management.md) - Student data
- [12-fees-payments.md](./12-fees-payments.md) - Transport fees
- [13-communication-system.md](./13-communication-system.md) - Parent alerts
- [index.md](./index.md) - Main PRD index

---

**Previous**: [16-certificate-generation.md](./16-certificate-generation.md)
**Next**: [18-library-management.md](./18-library-management.md)
