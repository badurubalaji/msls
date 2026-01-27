# Story 4.2: Guardian & Emergency Contact Management

**Epic:** 4 - Student Lifecycle Management
**Status:** done
**Priority:** P0
**Story Points:** 5

---

## User Story

**As an** administrator,
**I want** to manage student guardians and emergency contacts,
**So that** parents can be contacted and emergency situations handled.

---

## Acceptance Criteria

### AC1: Guardian Management
- [x] Admin can add guardians to a student profile (father, mother, other)
- [x] Guardian fields: name, relation, phone (required), email, occupation, address
- [x] One guardian must be marked as primary contact
- [x] Admin can enable/disable portal access for each guardian
- [x] Multiple guardians can be added (minimum 1 required)
- [x] Guardian records are soft-deletable

### AC2: Emergency Contacts
- [x] Admin can add emergency contacts beyond guardians
- [x] Each contact has: name, relation, phone (required)
- [x] Priority order can be set (1st, 2nd, 3rd)
- [x] Maximum 5 emergency contacts allowed
- [x] At least one emergency contact required (can be same as guardian)

### AC3: Parent Portal Access
- [ ] Guardian with portal access flag can login to parent portal (deferred to Epic 8)
- [ ] When portal access enabled, system creates user account linked to guardian (deferred to Epic 8)
- [ ] Guardian can update their own contact details via portal (deferred to Epic 8)
- [ ] Guardian sees only their linked children's data (deferred to Epic 8)

### AC4: Guardian Tab in Student Profile
- [x] Student profile shows "Guardians" tab
- [x] Tab displays list of guardians with primary indicator
- [x] Tab displays emergency contacts with priority
- [x] Quick actions: Add Guardian, Add Emergency Contact, Edit, Remove

---

## Technical Requirements

### Backend

**Database Tables:**

```sql
-- Migration: 20260126100100_create_guardians_table.up.sql

CREATE TABLE guardians (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id), -- Link to user account if portal access enabled
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    relation VARCHAR(20) NOT NULL CHECK (relation IN ('father', 'mother', 'guardian', 'grandparent', 'other')),
    phone VARCHAR(15) NOT NULL,
    email VARCHAR(255),
    occupation VARCHAR(100),
    address_line1 VARCHAR(255),
    address_line2 VARCHAR(255),
    city VARCHAR(100),
    state VARCHAR(100),
    postal_code VARCHAR(10),
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    has_portal_access BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id)
);

ALTER TABLE guardians ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON guardians USING (tenant_id = current_setting('app.current_tenant')::UUID);
CREATE INDEX idx_guardians_student ON guardians(student_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_guardians_user ON guardians(user_id) WHERE user_id IS NOT NULL;

-- Ensure only one primary guardian per student
CREATE UNIQUE INDEX uniq_guardians_primary ON guardians(student_id) WHERE is_primary = TRUE AND deleted_at IS NULL;

CREATE TABLE emergency_contacts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    relation VARCHAR(50) NOT NULL,
    phone VARCHAR(15) NOT NULL,
    priority SMALLINT NOT NULL DEFAULT 1 CHECK (priority BETWEEN 1 AND 5),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uniq_emergency_priority UNIQUE (student_id, priority)
);

ALTER TABLE emergency_contacts ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON emergency_contacts USING (tenant_id = current_setting('app.current_tenant')::UUID);
CREATE INDEX idx_emergency_contacts_student ON emergency_contacts(student_id);
```

**API Endpoints:**

| Method | Endpoint | Description | Permission |
|--------|----------|-------------|------------|
| GET | `/api/v1/students/{id}/guardians` | List guardians | `student:read` |
| POST | `/api/v1/students/{id}/guardians` | Add guardian | `student:update` |
| PUT | `/api/v1/students/{id}/guardians/{gid}` | Update guardian | `student:update` |
| DELETE | `/api/v1/students/{id}/guardians/{gid}` | Remove guardian | `student:update` |
| PUT | `/api/v1/students/{id}/guardians/{gid}/primary` | Set as primary | `student:update` |
| POST | `/api/v1/students/{id}/guardians/{gid}/portal-access` | Enable portal access | `student:update` |
| GET | `/api/v1/students/{id}/emergency-contacts` | List emergency contacts | `student:read` |
| POST | `/api/v1/students/{id}/emergency-contacts` | Add contact | `student:update` |
| PUT | `/api/v1/students/{id}/emergency-contacts/{cid}` | Update contact | `student:update` |
| DELETE | `/api/v1/students/{id}/emergency-contacts/{cid}` | Remove contact | `student:update` |

**Business Logic:**

```go
// Ensure at least one guardian when saving student
func (s *Service) ValidateGuardians(ctx context.Context, studentID uuid.UUID) error {
    count, err := s.guardianRepo.CountByStudent(ctx, studentID)
    if err != nil {
        return fmt.Errorf("count guardians: %w", err)
    }
    if count == 0 {
        return ErrAtLeastOneGuardianRequired
    }
    return nil
}

// Create user account when portal access enabled
func (s *Service) EnablePortalAccess(ctx context.Context, guardianID uuid.UUID) error {
    guardian, err := s.guardianRepo.GetByID(ctx, guardianID)
    if err != nil {
        return fmt.Errorf("get guardian: %w", err)
    }

    if guardian.Email == "" {
        return ErrEmailRequiredForPortalAccess
    }

    // Create user with 'parent' role
    user, err := s.userService.CreateParentUser(ctx, CreateParentUserDTO{
        Email:     guardian.Email,
        FirstName: guardian.FirstName,
        LastName:  guardian.LastName,
        Phone:     guardian.Phone,
    })
    if err != nil {
        return fmt.Errorf("create user: %w", err)
    }

    // Link user to guardian
    return s.guardianRepo.UpdateUserID(ctx, guardianID, user.ID)
}
```

### Frontend

**Components to create:**

```bash
ng generate component features/students/components/guardian-list --standalone
ng generate component features/students/components/guardian-form --standalone
ng generate component features/students/components/emergency-contact-form --standalone
ng generate interface features/students/models/guardian
ng generate interface features/students/models/emergency-contact
```

**Guardian Tab Component:**

```typescript
@Component({
  selector: 'app-guardian-list',
  standalone: true,
  imports: [CommonModule, ButtonComponent, BadgeComponent],
  template: `
    <div class="space-y-4">
      <div class="flex justify-between items-center">
        <h3 class="text-lg font-medium">Guardians</h3>
        <app-button (click)="showAddGuardian()" variant="primary" size="sm">
          Add Guardian
        </app-button>
      </div>

      @for (guardian of guardians(); track guardian.id) {
        <div class="border rounded-lg p-4">
          <div class="flex justify-between">
            <div>
              <span class="font-medium">{{ guardian.firstName }} {{ guardian.lastName }}</span>
              @if (guardian.isPrimary) {
                <app-badge variant="success" class="ml-2">Primary</app-badge>
              }
              <span class="text-gray-500 text-sm ml-2">({{ guardian.relation }})</span>
            </div>
            <div class="space-x-2">
              <app-button variant="ghost" size="sm" (click)="editGuardian(guardian)">Edit</app-button>
              <app-button variant="danger" size="sm" (click)="removeGuardian(guardian)">Remove</app-button>
            </div>
          </div>
          <div class="mt-2 text-sm text-gray-600">
            <p>Phone: {{ guardian.phone }}</p>
            @if (guardian.email) {
              <p>Email: {{ guardian.email }}</p>
            }
            <p>Portal Access: {{ guardian.hasPortalAccess ? 'Enabled' : 'Disabled' }}</p>
          </div>
        </div>
      }
    </div>
  `
})
export class GuardianListComponent {
  studentId = input.required<string>();
  guardians = signal<Guardian[]>([]);
}
```

---

## Tasks

### Backend Tasks

- [x] **BE-4.2.1**: Create guardian entity and migration
- [x] **BE-4.2.2**: Create emergency_contacts entity and migration
- [x] **BE-4.2.3**: Create guardian repository
- [x] **BE-4.2.4**: Create guardian service with portal access logic
- [x] **BE-4.2.5**: Create guardian HTTP handlers
- [x] **BE-4.2.6**: Create emergency contact handlers
- [x] **BE-4.2.7**: Add validation for at least one guardian
- [ ] **BE-4.2.8**: Write unit tests

### Frontend Tasks

- [x] **FE-4.2.1**: Create guardian and emergency contact interfaces
- [x] **FE-4.2.2**: Extend student service for guardian endpoints
- [x] **FE-4.2.3**: Create guardian list component
- [x] **FE-4.2.4**: Create guardian form component (add/edit)
- [x] **FE-4.2.5**: Create emergency contact form component
- [x] **FE-4.2.6**: Add Guardians tab to student detail page
- [ ] **FE-4.2.7**: Write component tests

---

## Definition of Done

- [x] All acceptance criteria verified (AC3 portal deferred to Epic 8)
- [x] At least one guardian validation working
- [x] Primary guardian constraint working
- [ ] Portal access creates user account correctly (deferred to Epic 8)
- [ ] Backend tests passing
- [ ] Frontend tests passing
- [x] No lint/type errors

---

## Dependencies

| Dependency | Status | Notes |
|------------|--------|-------|
| Story 4.1 (Student Profile) | Required | Guardian links to student |
| User module (Epic 2) | ✅ Done | For creating parent user accounts |
| Parent Portal (Epic 13) | Future | Portal access feature complete when Epic 13 done |

---

## Notes

- Portal access creates a user with 'parent' role
- Parent user can only see students linked via guardian records
- Email is required for portal access (used as login)
- Consider SMS notification when portal access granted
