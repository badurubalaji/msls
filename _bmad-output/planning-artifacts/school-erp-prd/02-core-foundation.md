# 02 - Core Foundation

**Parent Document**: [index.md](./index.md)
**Status**: Draft
**Phase**: 1 (MVP)

---

## 1. Overview

This document covers the foundational modules that every other module depends on:
- Institution Setup (Tenants, Branches, Academic Years)
- User Management
- Role-Based Access Control (RBAC)
- Configuration Engine

---

## 2. Institution Setup

### 2.1 Tenant Management (SaaS Mode)

**Entity: Tenant**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| name | VARCHAR(255) | School/Organization name |
| code | VARCHAR(50) | Unique tenant code |
| domain | VARCHAR(255) | Custom domain (optional) |
| logo_url | VARCHAR(500) | School logo |
| subscription_plan_id | UUID | Current plan |
| subscription_status | ENUM | active, suspended, trial, expired |
| trial_ends_at | TIMESTAMP | Trial expiry |
| settings | JSONB | Tenant-level settings |
| created_at | TIMESTAMP | Creation date |
| updated_at | TIMESTAMP | Last update |

**Subscription Plans**:
- Free (limited students, basic features)
- Standard (all Phase 1 features)
- Premium (all features + priority support)
- Enterprise (custom, unlimited)

### 2.2 Branch Management

**Entity: Branch**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Parent tenant |
| name | VARCHAR(255) | Branch name |
| code | VARCHAR(50) | Unique within tenant |
| address | TEXT | Physical address |
| city | VARCHAR(100) | City |
| state | VARCHAR(100) | State/Province |
| country | VARCHAR(100) | Country |
| pincode | VARCHAR(20) | Postal code |
| phone | VARCHAR(20) | Contact phone |
| email | VARCHAR(255) | Contact email |
| is_main | BOOLEAN | Is this the main branch? |
| is_active | BOOLEAN | Active status |
| settings | JSONB | Branch-specific settings |

**Branch Hierarchy**:
```
Tenant (School Group)
├── Branch 1 (Main Campus)
├── Branch 2 (Junior Wing)
└── Branch 3 (Senior Wing)
```

### 2.3 Academic Year Management

**Entity: AcademicYear**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Parent tenant |
| name | VARCHAR(50) | Display name (e.g., "2025-26") |
| start_date | DATE | Year start |
| end_date | DATE | Year end |
| is_current | BOOLEAN | Currently active year |
| is_locked | BOOLEAN | Prevent modifications |
| settings | JSONB | Year-specific configurations |

**Term/Semester Management**:

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| academic_year_id | UUID | Parent year |
| name | VARCHAR(50) | Term name (e.g., "Term 1", "Semester 1") |
| start_date | DATE | Term start |
| end_date | DATE | Term end |
| is_current | BOOLEAN | Currently active term |

---

## 3. User Management

### 3.1 User Entity

**Entity: User**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Parent tenant |
| email | VARCHAR(255) | Unique email (login) |
| phone | VARCHAR(20) | Phone number |
| password_hash | VARCHAR(255) | Bcrypt hash |
| first_name | VARCHAR(100) | First name |
| last_name | VARCHAR(100) | Last name |
| display_name | VARCHAR(200) | Full display name |
| avatar_url | VARCHAR(500) | Profile picture |
| user_type | ENUM | staff, student, parent, admin |
| status | ENUM | active, inactive, suspended, pending |
| email_verified | BOOLEAN | Email verification status |
| phone_verified | BOOLEAN | Phone verification status |
| last_login_at | TIMESTAMP | Last login time |
| password_changed_at | TIMESTAMP | Last password change |
| settings | JSONB | User preferences |
| created_at | TIMESTAMP | Creation date |

### 3.2 User-Branch Association

Users can belong to multiple branches:

**Entity: UserBranch**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| user_id | UUID | User reference |
| branch_id | UUID | Branch reference |
| is_primary | BOOLEAN | Primary branch |
| assigned_at | TIMESTAMP | Assignment date |

### 3.3 User Types

| Type | Description | Portal Access |
|------|-------------|---------------|
| **admin** | System administrators | Admin dashboard |
| **staff** | Teachers, non-teaching staff | Staff portal |
| **student** | Enrolled students | Student portal |
| **parent** | Student guardians | Parent portal |

### 3.4 Authentication Flows

**Login Flow**:
```
1. User submits email/phone + password
2. Validate credentials
3. Check user status (active, not suspended)
4. Generate JWT + Refresh Token
5. Store refresh token in Redis
6. Return tokens to client
```

**Password Reset Flow**:
```
1. User requests reset (email/phone)
2. Generate OTP (6 digits, 10 min expiry)
3. Send via email/SMS
4. User submits OTP + new password
5. Validate OTP, update password
6. Invalidate all existing sessions
```

**Session Management**:
- Single session per device (optional)
- Force logout from all devices
- Session timeout: 30 minutes inactivity

---

## 4. Role-Based Access Control (RBAC)

### 4.1 Role Definition

**Entity: Role**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | NULL for system roles |
| name | VARCHAR(100) | Role name |
| code | VARCHAR(50) | Unique code |
| description | TEXT | Role description |
| is_system | BOOLEAN | System-defined (non-editable) |
| parent_role_id | UUID | Inherits from parent |
| created_at | TIMESTAMP | Creation date |

### 4.2 System Roles (Pre-defined)

| Role Code | Name | Level | Description |
|-----------|------|-------|-------------|
| super_admin | Super Admin | 0 | Platform owner (SaaS only) |
| school_admin | School Admin | 1 | Full school access |
| branch_admin | Branch Admin | 2 | Full branch access |
| principal | Principal | 3 | Academic head |
| vice_principal | Vice Principal | 4 | Deputy academic head |
| hod | Head of Department | 5 | Department head |
| teacher | Teacher | 6 | Regular teacher |
| class_teacher | Class Teacher | 6 | Class in-charge |
| accountant | Accountant | 6 | Finance access |
| librarian | Librarian | 6 | Library access |
| transport_incharge | Transport In-charge | 6 | Transport access |
| receptionist | Receptionist | 7 | Front desk |
| data_entry | Data Entry Operator | 7 | Limited data entry |

### 4.3 Permission Structure

**Entity: Permission**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| module | VARCHAR(50) | Module name |
| resource | VARCHAR(50) | Resource name |
| action | VARCHAR(50) | Action name |
| code | VARCHAR(150) | Full permission code |
| description | TEXT | Permission description |

**Permission Format**: `module:resource:action`

**Examples**:
```
student:profile:read
student:profile:write
student:profile:delete
fees:payment:collect
fees:payment:refund
exam:result:view
exam:result:publish
attendance:student:mark
attendance:student:edit
```

### 4.4 Role-Permission Mapping

**Entity: RolePermission**

| Field | Type | Description |
|-------|------|-------------|
| role_id | UUID | Role reference |
| permission_id | UUID | Permission reference |
| granted_at | TIMESTAMP | When granted |
| granted_by | UUID | Who granted |

### 4.5 User-Role Assignment

**Entity: UserRole**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| user_id | UUID | User reference |
| role_id | UUID | Role reference |
| branch_id | UUID | Branch scope (NULL = all branches) |
| valid_from | DATE | Start date |
| valid_until | DATE | End date (NULL = indefinite) |
| assigned_by | UUID | Who assigned |
| assigned_at | TIMESTAMP | When assigned |

### 4.6 Permission Check Flow

```go
func HasPermission(userID, branchID, permission string) bool {
    // 1. Get user's roles for the branch
    roles := GetUserRoles(userID, branchID)

    // 2. For each role, check if permission exists
    for _, role := range roles {
        if RoleHasPermission(role.ID, permission) {
            return true
        }
        // Check inherited permissions from parent roles
        if role.ParentRoleID != nil {
            if HasInheritedPermission(role.ParentRoleID, permission) {
                return true
            }
        }
    }
    return false
}
```

---

## 5. Configuration Engine

### 5.1 Feature Flags

**Entity: FeatureFlag**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant scope |
| branch_id | UUID | Branch scope (NULL = all) |
| feature_code | VARCHAR(100) | Feature identifier |
| is_enabled | BOOLEAN | Enabled status |
| config | JSONB | Feature-specific config |
| updated_by | UUID | Last updated by |
| updated_at | TIMESTAMP | Last update |

**Core Feature Flags**:
```yaml
features:
  - code: module.admissions
    default: true
  - code: module.library
    default: false
  - code: module.transport
    default: false
  - code: module.hostel
    default: false
  - code: module.online_quiz
    default: false
  - code: feature.biometric_attendance
    default: false
  - code: feature.sms_notifications
    default: true
  - code: feature.parent_app
    default: true
```

### 5.2 Custom Fields

**Entity: CustomFieldDefinition**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | Tenant scope |
| entity_type | VARCHAR(50) | student, staff, etc. |
| field_name | VARCHAR(100) | Internal field name |
| display_name | VARCHAR(200) | UI label |
| field_type | ENUM | text, number, date, dropdown, file, boolean |
| is_required | BOOLEAN | Mandatory field |
| is_searchable | BOOLEAN | Include in search |
| options | JSONB | For dropdown type |
| validation | JSONB | Validation rules |
| display_order | INT | UI ordering |
| is_active | BOOLEAN | Active status |

**Entity: CustomFieldValue**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| field_definition_id | UUID | Field definition |
| entity_type | VARCHAR(50) | Entity type |
| entity_id | UUID | Entity ID |
| value | JSONB | Field value |
| updated_at | TIMESTAMP | Last update |

### 5.3 System Settings

**Entity: Setting**

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| tenant_id | UUID | NULL for global |
| branch_id | UUID | NULL for tenant-wide |
| category | VARCHAR(50) | Setting category |
| key | VARCHAR(100) | Setting key |
| value | JSONB | Setting value |
| value_type | VARCHAR(20) | string, number, boolean, json |
| is_sensitive | BOOLEAN | Hide in UI |
| updated_by | UUID | Last updated by |
| updated_at | TIMESTAMP | Last update |

**Setting Categories**:
```yaml
categories:
  - general: School name, logo, contact info
  - academic: Grading scales, pass marks
  - attendance: Working days, grace period
  - fees: Late fee rules, payment modes
  - communication: SMS provider, email config
  - security: Password policy, session timeout
```

---

## 6. API Endpoints

### 6.1 Institution APIs

```
# Branches
GET    /api/v1/branches                    # List branches
POST   /api/v1/branches                    # Create branch
GET    /api/v1/branches/{id}               # Get branch
PUT    /api/v1/branches/{id}               # Update branch
DELETE /api/v1/branches/{id}               # Delete branch

# Academic Years
GET    /api/v1/academic-years              # List years
POST   /api/v1/academic-years              # Create year
GET    /api/v1/academic-years/{id}         # Get year
PUT    /api/v1/academic-years/{id}         # Update year
POST   /api/v1/academic-years/{id}/activate # Set as current
```

### 6.2 User APIs

```
# Users
GET    /api/v1/users                       # List users
POST   /api/v1/users                       # Create user
GET    /api/v1/users/{id}                  # Get user
PUT    /api/v1/users/{id}                  # Update user
DELETE /api/v1/users/{id}                  # Deactivate user
POST   /api/v1/users/{id}/reset-password   # Reset password

# Authentication
POST   /api/v1/auth/login                  # Login
POST   /api/v1/auth/logout                 # Logout
POST   /api/v1/auth/refresh                # Refresh token
POST   /api/v1/auth/forgot-password        # Request reset
POST   /api/v1/auth/verify-otp             # Verify OTP
```

### 6.3 RBAC APIs

```
# Roles
GET    /api/v1/roles                       # List roles
POST   /api/v1/roles                       # Create custom role
GET    /api/v1/roles/{id}                  # Get role
PUT    /api/v1/roles/{id}                  # Update role
GET    /api/v1/roles/{id}/permissions      # Get role permissions
PUT    /api/v1/roles/{id}/permissions      # Update permissions

# User Roles
GET    /api/v1/users/{id}/roles            # Get user roles
POST   /api/v1/users/{id}/roles            # Assign role
DELETE /api/v1/users/{id}/roles/{roleId}   # Remove role
```

---

## 7. Related Documents

- [01-technical-architecture.md](./01-technical-architecture.md) - Technical stack
- [03-student-management.md](./03-student-management.md) - Student module
- [index.md](./index.md) - Main PRD index

---

**Previous**: [01-technical-architecture.md](./01-technical-architecture.md)
**Next**: [03-student-management.md](./03-student-management.md)
