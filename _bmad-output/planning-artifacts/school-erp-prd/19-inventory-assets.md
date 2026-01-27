# 19 - Inventory & Assets

**Parent Document**: [index.md](./index.md)
**Status**: Draft
**Phase**: 2

---

## 1. Overview

The Inventory & Assets module manages school assets (furniture, equipment, lab instruments), consumables, and stock management.

---

## 2. Asset Management

### 2.1 Asset Entity

**Entity: Asset**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| branch_id | UUID | Branch reference |
| asset_code | VARCHAR(50) | Unique asset ID |
| name | VARCHAR(200) | Asset name |
| category_id | UUID | Category |
| description | TEXT | Description |
| make | VARCHAR(100) | Manufacturer |
| model | VARCHAR(100) | Model |
| serial_number | VARCHAR(100) | Serial number |
| purchase_date | DATE | Purchase date |
| purchase_price | DECIMAL | Cost |
| vendor_id | UUID | Supplier |
| warranty_expiry | DATE | Warranty end |
| depreciation_method | ENUM | straight_line, declining |
| useful_life_years | INT | Expected life |
| current_value | DECIMAL | Depreciated value |
| location_id | UUID | Current location |
| assigned_to | UUID | Assigned staff |
| condition | ENUM | new, good, fair, poor, damaged |
| status | ENUM | active, maintenance, disposed, lost |
| barcode | VARCHAR(50) | Asset barcode |
| qr_code | VARCHAR(500) | Asset QR |
| photos | JSONB | Asset photos |

### 2.2 Asset Categories

```
Asset Categories:
├── Furniture
│   ├── Desks & Tables
│   ├── Chairs
│   ├── Cabinets & Storage
│   └── Boards & Displays
├── IT Equipment
│   ├── Computers
│   ├── Laptops
│   ├── Printers
│   ├── Projectors
│   └── Networking
├── Lab Equipment
│   ├── Science Lab
│   ├── Computer Lab
│   └── Language Lab
├── Sports Equipment
├── Musical Instruments
├── Vehicles
└── Building & Infrastructure
```

### 2.3 Asset Register

```
Asset Register | IT Equipment | Computer Lab

| Asset Code | Item           | Serial      | Location | Status |
|------------|----------------|-------------|----------|--------|
| IT-PC-001  | Desktop PC     | DEL-12345   | Lab-1    | Active |
| IT-PC-002  | Desktop PC     | DEL-12346   | Lab-1    | Active |
| IT-PRJ-001 | Projector      | EPSON-789   | Lab-1    | Active |
| IT-PC-003  | Desktop PC     | DEL-12347   | Lab-1    | Maint. |

Total Assets: 45 | Active: 42 | Maintenance: 3
Total Value: ₹18,50,000
```

---

## 3. Locations

### 3.1 Location Entity

**Entity: AssetLocation**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| branch_id | UUID | Branch |
| name | VARCHAR(100) | Location name |
| code | VARCHAR(20) | Location code |
| type | ENUM | room, building, floor, outdoor |
| parent_id | UUID | Parent location |
| capacity | INT | Asset capacity |
| in_charge_id | UUID | Responsible staff |

### 3.2 Location Hierarchy

```
Main Building
├── Ground Floor
│   ├── Reception
│   ├── Office
│   └── Auditorium
├── First Floor
│   ├── Class 1-A to 1-D
│   ├── Class 2-A to 2-D
│   └── Staff Room
└── Second Floor
    ├── Computer Lab
    ├── Science Lab
    └── Library
```

---

## 4. Maintenance

### 4.1 Maintenance Request

**Entity: MaintenanceRequest**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| asset_id | UUID | Asset reference |
| request_number | VARCHAR(50) | Request ID |
| issue_description | TEXT | Problem description |
| priority | ENUM | low, medium, high, critical |
| requested_by | UUID | Reporter |
| request_date | DATE | Request date |
| assigned_to | UUID | Technician |
| status | ENUM | open, in_progress, resolved, closed |
| resolution | TEXT | Resolution details |
| cost | DECIMAL | Repair cost |
| resolved_date | DATE | Resolution date |

### 4.2 Maintenance Schedule

```
Scheduled Maintenance | January 2026

| Asset          | Type       | Due Date | Status   |
|----------------|------------|----------|----------|
| AC Units (10)  | Servicing  | 15-Jan   | Due      |
| Projectors (5) | Cleaning   | 20-Jan   | Upcoming |
| Computers (45) | Antivirus  | 01-Jan   | Done     |
| Fire Ext. (20) | Inspection | 01-Feb   | Upcoming |
```

---

## 5. Stock/Consumables

### 5.1 Stock Item

**Entity: StockItem**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant reference |
| branch_id | UUID | Branch |
| name | VARCHAR(200) | Item name |
| category_id | UUID | Category |
| unit | VARCHAR(20) | Unit of measure |
| sku | VARCHAR(50) | Stock keeping unit |
| current_stock | DECIMAL | Current quantity |
| reorder_level | DECIMAL | Reorder trigger |
| reorder_quantity | DECIMAL | Reorder amount |
| unit_price | DECIMAL | Average price |
| location_id | UUID | Storage location |

### 5.2 Stock Categories

```
Consumables:
├── Stationery
│   ├── Paper & Notebooks
│   ├── Pens & Pencils
│   └── Office Supplies
├── Cleaning Materials
├── Lab Consumables
├── Sports Consumables
├── First Aid Supplies
└── Electrical Items
```

### 5.3 Stock Movement

**Entity: StockMovement**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| stock_item_id | UUID | Item reference |
| movement_type | ENUM | in, out, adjustment |
| quantity | DECIMAL | Quantity |
| reference_type | VARCHAR(50) | PO, Issue, Adjustment |
| reference_id | UUID | Reference document |
| remarks | TEXT | Notes |
| moved_by | UUID | Staff |
| movement_date | DATE | Date |

---

## 6. Purchase & Vendors

### 6.1 Purchase Order

**Entity: PurchaseOrder**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| po_number | VARCHAR(50) | PO number |
| vendor_id | UUID | Vendor |
| order_date | DATE | Order date |
| expected_date | DATE | Expected delivery |
| total_amount | DECIMAL | Order value |
| status | ENUM | draft, approved, ordered, received, cancelled |
| approved_by | UUID | Approver |
| items | JSONB | Line items |

### 6.2 Vendor Entity

**Entity: Vendor**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| name | VARCHAR(200) | Vendor name |
| contact_person | VARCHAR(200) | Contact |
| phone | VARCHAR(20) | Phone |
| email | VARCHAR(255) | Email |
| address | TEXT | Address |
| gst_number | VARCHAR(20) | GST |
| pan_number | VARCHAR(10) | PAN |
| categories | JSONB | Supply categories |
| rating | INT | 1-5 rating |
| status | ENUM | active, inactive, blacklisted |

---

## 7. Reports

```
Inventory Summary | January 2026

Assets:
- Total Assets: 850
- Total Value: ₹1.2 Cr
- Under Maintenance: 15
- Disposed This Year: 8

Stock Alerts:
⚠️ Whiteboard Markers - Below reorder level (50 left)
⚠️ A4 Paper - Below reorder level (5 reams left)
⚠️ Lab Chemicals - Expiring soon (30 days)

Recent Purchases:
| Date    | Item              | Qty | Amount   |
|---------|-------------------|-----|----------|
| 15-Jan  | Desktop Computers | 10  | ₹4,50,000|
| 10-Jan  | Office Chairs     | 25  | ₹75,000  |
| 05-Jan  | Stationery        | -   | ₹15,000  |
```

---

## 8. API Endpoints

```
# Assets
GET    /api/v1/assets                       # List assets
POST   /api/v1/assets                       # Add asset
GET    /api/v1/assets/{id}                  # Get asset
PUT    /api/v1/assets/{id}                  # Update asset
POST   /api/v1/assets/{id}/transfer         # Transfer location
POST   /api/v1/assets/{id}/dispose          # Dispose asset

# Maintenance
POST   /api/v1/maintenance/requests         # Create request
GET    /api/v1/maintenance/requests         # List requests
PUT    /api/v1/maintenance/requests/{id}    # Update request

# Stock
GET    /api/v1/stock                        # List stock
POST   /api/v1/stock/receive                # Stock receipt
POST   /api/v1/stock/issue                  # Stock issue
GET    /api/v1/stock/alerts                 # Low stock alerts

# Purchase
POST   /api/v1/purchase-orders              # Create PO
GET    /api/v1/purchase-orders              # List POs
POST   /api/v1/purchase-orders/{id}/approve # Approve PO
```

---

## 9. Related Documents

- [10-staff-management.md](./10-staff-management.md) - Staff assignments
- [12-fees-payments.md](./12-fees-payments.md) - Purchase payments
- [index.md](./index.md) - Main PRD index

---

**Previous**: [18-library-management.md](./18-library-management.md)
**Next**: [20-visitor-gate-management.md](./20-visitor-gate-management.md)
