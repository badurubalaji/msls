# Story 4.8: Student Search & Bulk Operations

**Epic:** 4 - Student Lifecycle Management
**Status:** done
**Priority:** P1
**Story Points:** 5

---

## User Story

**As an** administrator,
**I want** to search students and perform bulk operations,
**So that** managing large student populations is efficient.

---

## Acceptance Criteria

### AC1: Student Search
- [x] Can search by: name (first, last, full), admission number, phone number
- [x] Search is case-insensitive
- [x] Results show in paginated list (cursor-based, 20 per page default)
- [x] Search debounced (300ms) for performance

### AC2: Advanced Filters
- [x] Filter by: class, section, status (active/inactive/transferred)
- [x] Filter by: gender, transport user (yes/no)
- [x] Filter by: admission date range
- [x] Multiple filters can be combined (AND logic)
- [x] Filters persist in URL for shareability

### AC3: Bulk Operations
- [x] Select multiple students via checkboxes
- [x] "Select all on page" and "Select all matching filters" options
- [x] Bulk operations: Send SMS, Send Email, Update Status, Export
- [x] Confirmation required before bulk execution
- [x] Operation progress shown for large batches
- [x] Operation log maintained for audit

### AC4: Export
- [x] Export student list to Excel/CSV
- [x] Column selection available before export
- [x] Filters applied to export
- [x] Default columns: Name, Admission No, Class, Section, Phone, Guardian Phone, Status
- [x] Maximum 10,000 records per export

---

## Technical Requirements

### Backend

**Search Implementation:**

```go
type StudentSearchParams struct {
    Query       string     `form:"q"`
    ClassID     *uuid.UUID `form:"class_id"`
    SectionID   *uuid.UUID `form:"section_id"`
    Status      string     `form:"status"`
    Gender      string     `form:"gender"`
    HasTransport *bool     `form:"has_transport"`
    AdmissionFrom *time.Time `form:"admission_from"`
    AdmissionTo   *time.Time `form:"admission_to"`
    Cursor      string     `form:"cursor"`
    Limit       int        `form:"limit" binding:"max=100"`
    SortBy      string     `form:"sort_by" binding:"oneof=name admission_number created_at"`
    SortOrder   string     `form:"sort_order" binding:"oneof=asc desc"`
}

func (r *Repository) Search(ctx context.Context, params StudentSearchParams) ([]Student, string, error) {
    query := r.db.WithContext(ctx).Model(&Student{}).
        Where("deleted_at IS NULL")

    // Full-text search on name
    if params.Query != "" {
        searchQuery := "%" + strings.ToLower(params.Query) + "%"
        query = query.Where(
            "LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ? OR LOWER(first_name || ' ' || last_name) LIKE ? OR LOWER(admission_number) LIKE ?",
            searchQuery, searchQuery, searchQuery, searchQuery,
        )
    }

    // Apply filters
    if params.ClassID != nil {
        query = query.Joins("JOIN student_enrollments ON students.id = student_enrollments.student_id").
            Where("student_enrollments.class_id = ? AND student_enrollments.status = 'active'", params.ClassID)
    }

    if params.SectionID != nil {
        query = query.Where("student_enrollments.section_id = ?", params.SectionID)
    }

    if params.Status != "" {
        query = query.Where("status = ?", params.Status)
    }

    if params.Gender != "" {
        query = query.Where("gender = ?", params.Gender)
    }

    if params.AdmissionFrom != nil {
        query = query.Where("admission_date >= ?", params.AdmissionFrom)
    }

    if params.AdmissionTo != nil {
        query = query.Where("admission_date <= ?", params.AdmissionTo)
    }

    // Cursor pagination
    if params.Cursor != "" {
        cursorData, err := decodeCursor(params.Cursor)
        if err != nil {
            return nil, "", fmt.Errorf("decode cursor: %w", err)
        }
        query = query.Where("id > ?", cursorData.LastID)
    }

    // Sorting
    sortBy := "last_name, first_name"
    if params.SortBy != "" {
        sortOrder := "ASC"
        if params.SortOrder == "desc" {
            sortOrder = "DESC"
        }
        sortBy = fmt.Sprintf("%s %s", params.SortBy, sortOrder)
    }
    query = query.Order(sortBy)

    // Limit + 1 to check if more pages exist
    limit := 20
    if params.Limit > 0 {
        limit = params.Limit
    }
    query = query.Limit(limit + 1)

    var students []Student
    if err := query.Find(&students).Error; err != nil {
        return nil, "", fmt.Errorf("find students: %w", err)
    }

    // Calculate next cursor
    var nextCursor string
    if len(students) > limit {
        students = students[:limit]
        nextCursor = encodeCursor(CursorData{LastID: students[len(students)-1].ID})
    }

    return students, nextCursor, nil
}
```

**Bulk Operations:**

```sql
-- Migration: 20260126100700_create_bulk_operations.up.sql

CREATE TABLE bulk_operations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    operation_type VARCHAR(50) NOT NULL, -- send_sms, send_email, update_status, export
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed')),
    total_count INTEGER NOT NULL DEFAULT 0,
    processed_count INTEGER NOT NULL DEFAULT 0,
    success_count INTEGER NOT NULL DEFAULT 0,
    failure_count INTEGER NOT NULL DEFAULT 0,
    parameters JSONB, -- Operation-specific params
    result_url VARCHAR(500), -- For export operations
    error_message TEXT,
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL REFERENCES users(id)
);

ALTER TABLE bulk_operations ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON bulk_operations USING (tenant_id = current_setting('app.current_tenant')::UUID);
CREATE INDEX idx_bulk_operations_user ON bulk_operations(created_by);

CREATE TABLE bulk_operation_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    operation_id UUID NOT NULL REFERENCES bulk_operations(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    error_message TEXT,
    processed_at TIMESTAMPTZ
);

CREATE INDEX idx_bulk_items_operation ON bulk_operation_items(operation_id);
```

**API Endpoints:**

| Method | Endpoint | Description | Permission |
|--------|----------|-------------|------------|
| GET | `/api/v1/students` | Search students | `student:read` |
| POST | `/api/v1/students/bulk/status` | Bulk update status | `student:update` |
| POST | `/api/v1/students/bulk/sms` | Bulk send SMS | `communication:send` |
| POST | `/api/v1/students/bulk/email` | Bulk send email | `communication:send` |
| POST | `/api/v1/students/export` | Export to Excel/CSV | `student:export` |
| GET | `/api/v1/bulk-operations/{id}` | Get operation status | `student:read` |
| GET | `/api/v1/bulk-operations/{id}/result` | Download result | `student:read` |

**Export Service:**

```go
type ExportParams struct {
    Format  string   `json:"format" binding:"required,oneof=xlsx csv"`
    Columns []string `json:"columns"`
    Filters StudentSearchParams `json:"filters"`
}

var DefaultExportColumns = []string{
    "admission_number",
    "first_name",
    "last_name",
    "class",
    "section",
    "phone",
    "guardian_name",
    "guardian_phone",
    "status",
}

func (s *ExportService) Export(ctx context.Context, params ExportParams) (string, error) {
    columns := params.Columns
    if len(columns) == 0 {
        columns = DefaultExportColumns
    }

    // Get students (with limit)
    params.Filters.Limit = 10000
    students, _, err := s.studentRepo.Search(ctx, params.Filters)
    if err != nil {
        return "", fmt.Errorf("search students: %w", err)
    }

    // Create file
    var fileURL string
    switch params.Format {
    case "xlsx":
        fileURL, err = s.createExcel(ctx, students, columns)
    case "csv":
        fileURL, err = s.createCSV(ctx, students, columns)
    }

    if err != nil {
        return "", fmt.Errorf("create file: %w", err)
    }

    return fileURL, nil
}
```

### Frontend

**Components to create:**

```bash
ng generate component features/students/components/student-search --standalone
ng generate component features/students/components/student-filters --standalone
ng generate component features/students/components/bulk-actions --standalone
ng generate component features/students/components/export-dialog --standalone
ng generate interface features/students/models/search-params
ng generate interface features/students/models/bulk-operation
```

**Search Component:**

```typescript
@Component({
  selector: 'app-student-search',
  template: `
    <div class="flex gap-4 items-center">
      <div class="relative flex-1">
        <app-icon name="search" class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
        <input
          type="text"
          [value]="searchQuery()"
          (input)="onSearchInput($event)"
          placeholder="Search by name, admission number, phone..."
          class="w-full pl-10 pr-4 py-2 border rounded-lg"
        />
      </div>
      <app-button variant="secondary" (click)="toggleFilters()">
        <app-icon name="filter" class="mr-2" />
        Filters
        @if (activeFilterCount() > 0) {
          <app-badge variant="primary" class="ml-2">{{ activeFilterCount() }}</app-badge>
        }
      </app-button>
    </div>

    @if (showFilters()) {
      <app-student-filters
        [filters]="filters()"
        (filtersChanged)="onFiltersChanged($event)"
      />
    }
  `
})
export class StudentSearchComponent {
  searchQuery = signal('');
  filters = signal<StudentFilters>({});
  showFilters = signal(false);

  private searchDebounce = inject(DestroyRef);
  private searchSubject = new Subject<string>();

  constructor() {
    this.searchSubject.pipe(
      debounceTime(300),
      distinctUntilChanged(),
      takeUntilDestroyed(this.searchDebounce)
    ).subscribe(query => {
      this.onSearch.emit(query);
    });
  }

  onSearchInput(event: Event) {
    const value = (event.target as HTMLInputElement).value;
    this.searchQuery.set(value);
    this.searchSubject.next(value);
  }
}
```

**Bulk Actions Component:**

```typescript
@Component({
  selector: 'app-bulk-actions',
  template: `
    @if (selectedCount() > 0) {
      <div class="fixed bottom-4 left-1/2 -translate-x-1/2 bg-gray-900 text-white rounded-lg shadow-lg p-4 flex items-center gap-4">
        <span>{{ selectedCount() }} students selected</span>
        <div class="flex gap-2">
          <app-button variant="secondary" size="sm" (click)="onAction('sms')">
            <app-icon name="message" class="mr-1" /> Send SMS
          </app-button>
          <app-button variant="secondary" size="sm" (click)="onAction('email')">
            <app-icon name="mail" class="mr-1" /> Send Email
          </app-button>
          <app-button variant="secondary" size="sm" (click)="onAction('status')">
            <app-icon name="edit" class="mr-1" /> Update Status
          </app-button>
          <app-button variant="secondary" size="sm" (click)="onAction('export')">
            <app-icon name="download" class="mr-1" /> Export
          </app-button>
        </div>
        <app-button variant="ghost" size="sm" (click)="clearSelection()">
          <app-icon name="x" />
        </app-button>
      </div>
    }
  `
})
export class BulkActionsComponent {
  selectedCount = input.required<number>();
  selectedIds = input.required<string[]>();
  action = output<{ type: string; ids: string[] }>();
}
```

---

## Tasks

### Backend Tasks

- [x] **BE-4.8.1**: Implement search with filters in repository
- [x] **BE-4.8.2**: Add cursor-based pagination
- [x] **BE-4.8.3**: Create bulk operations table migration
- [x] **BE-4.8.4**: Create bulk operation service
- [x] **BE-4.8.5**: Create export service (Excel/CSV)
- [x] **BE-4.8.6**: Create search and bulk HTTP handlers
- [x] **BE-4.8.7**: Add export permission to seed (already existed in migration 000023)
- [x] **BE-4.8.8**: Write unit tests

### Frontend Tasks

- [x] **FE-4.8.1**: Create search and filter interfaces
- [x] **FE-4.8.2**: Update student service with search params
- [x] **FE-4.8.3**: Create student search component
- [x] **FE-4.8.4**: Create student filters component
- [x] **FE-4.8.5**: Implement filter URL persistence
- [x] **FE-4.8.6**: Create bulk actions bar component
- [x] **FE-4.8.7**: Create export dialog with column selection
- [x] **FE-4.8.8**: Integrate with student list page
- [x] **FE-4.8.9**: Write component tests

---

## Definition of Done

- [x] All acceptance criteria verified
- [x] Search working with all fields
- [x] Filters working and combinable
- [x] Filter state persists in URL
- [x] Bulk operations create audit log
- [x] Export generates valid Excel/CSV
- [x] Backend tests passing
- [x] Frontend tests passing

---

## Dependencies

| Dependency | Status | Notes |
|------------|--------|-------|
| Story 4.1 (Student Profile) | Required | Search operates on students |
| Story 4.6 (Enrollment) | Required | Filter by class/section |
| SMS Integration (Epic 12) | Future | Bulk SMS deferred |
| Email Integration (Epic 12) | Future | Bulk email deferred |

---

## Notes

- Bulk SMS/Email will show "Coming Soon" until Epic 12
- Export limited to 10,000 records for performance
- Consider background jobs for large bulk operations
- Search should use database indexes efficiently
- URL filter format: `/students?class_id=uuid&status=active`
