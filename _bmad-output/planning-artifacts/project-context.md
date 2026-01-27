# MSLS Project Context

> **Purpose**: Concise implementation rules for AI agents. Read this BEFORE writing any code.
> **Last Updated**: 2026-01-23
> **Architecture**: See `architecture.md` for full details

---

## Critical Technology Rules

### Go Backend Rules

```go
// ALWAYS use these patterns:

// 1. Error handling - ALWAYS wrap with context
if err != nil {
    return fmt.Errorf("operation_name failed: %w", err)
}

// 2. Context propagation - ALWAYS pass context
func (s *Service) DoSomething(ctx context.Context, req Request) error

// 3. Tenant isolation - ALWAYS include tenant_id
func (r *Repository) FindByID(ctx context.Context, tenantID, id uuid.UUID) (*Entity, error)

// 4. UUID v7 for all IDs
id := uuid.Must(uuid.NewV7())

// 5. Transactions for multi-table operations
tx := r.db.WithContext(ctx).Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()
```

### Go File Structure (Per Module)

```
internal/modules/{module}/
├── handler.go      # HTTP handlers, request validation
├── service.go      # Business logic, orchestration
├── repository.go   # Database operations
├── dto.go          # Request/Response structs
├── entity.go       # Database models
└── {module}_test.go
```

### Go Best Practices (MANDATORY)

Backend developers MUST follow the Google Go Style Guide:

| Resource | Link |
|----------|------|
| **Best Practices** | https://google.github.io/styleguide/go/best-practices |
| **Style Guide** | https://google.github.io/styleguide/go/guide |
| **Decisions** | https://google.github.io/styleguide/go/decisions |

**Key Requirements:**
1. **Naming**: Use MixedCaps, not underscores; keep names short but descriptive
2. **Error Handling**: Always handle errors; wrap with context using `fmt.Errorf("...: %w", err)`
3. **Documentation**: All exported functions/types must have doc comments
4. **Testing**: Table-driven tests preferred; use `t.Helper()` in test helpers
5. **Concurrency**: Prefer channels over shared memory; use `context.Context` for cancellation

### Angular Frontend Rules

```typescript
// 1. Standalone components ONLY
@Component({
  standalone: true,
  imports: [CommonModule, RouterModule],
  // ...
})

// 2. Signals for state (NOT NgRx)
export class MyComponent {
  data = signal<Data[]>([]);
  loading = signal(false);

  readonly filteredData = computed(() =>
    this.data().filter(d => d.active)
  );
}

// 3. Inject function (NOT constructor injection)
export class MyService {
  private http = inject(HttpClient);
  private auth = inject(AuthService);
}

// 4. Typed HTTP responses
return this.http.get<ApiResponse<Student[]>>('/api/v1/students');

// 5. Error handling in services
return this.http.get<T>(url).pipe(
  catchError(this.handleError.bind(this))
);
```

### Angular File Structure (Per Feature)

```
src/app/features/{feature}/
├── {feature}.routes.ts       # Lazy-loaded routes
├── components/
│   ├── {name}/{name}.component.ts
│   └── {name}/{name}.component.html
├── services/
│   └── {feature}.service.ts
├── models/
│   └── {feature}.model.ts
└── guards/
    └── {feature}.guard.ts
```

### Angular CLI Commands (MANDATORY)

Frontend developers MUST use Angular CLI (`ng`) commands for all code generation. **Never create files manually.**

```bash
# Project initialization
ng new msls-frontend --style=scss --routing=true --ssr=false --standalone=true --strict=true

# Components (always use)
ng generate component features/students/components/student-list --standalone
ng generate component shared/components/button --standalone --export

# Services
ng generate service features/students/services/student
ng generate service core/services/auth

# Guards
ng generate guard core/guards/auth --functional

# Interceptors
ng generate interceptor core/interceptors/auth --functional

# Pipes
ng generate pipe shared/pipes/date-format

# Directives
ng generate directive shared/directives/tooltip

# Interfaces/Models
ng generate interface features/students/models/student
ng generate interface core/models/api-response

# Environments
ng generate environments
```

**Rationale:**
- Ensures consistent file structure and naming
- Automatically updates module imports (where applicable)
- Generates spec files for testing
- Follows Angular conventions exactly

### Angular Best Practices (MANDATORY)

Frontend developers MUST follow these official Angular best practices:

| Category | Documentation |
|----------|---------------|
| **Style Guide** | https://angular.dev/style-guide |
| **Security** | https://angular.dev/best-practices/security |
| **Accessibility (a11y)** | https://angular.dev/best-practices/a11y |
| **Error Handling** | https://angular.dev/best-practices/error-handling |
| **Runtime Performance** | https://angular.dev/best-practices/runtime-performance |
| **Zone Pollution** | https://angular.dev/best-practices/zone-pollution |
| **Slow Computations** | https://angular.dev/best-practices/slow-computations |
| **Skipping Subtrees** | https://angular.dev/best-practices/skipping-subtrees |
| **Chrome DevTools Profiling** | https://angular.dev/best-practices/profiling-with-chrome-devtools |
| **Zoneless Angular** | https://angular.dev/guide/zoneless |
| **Tailwind CSS Integration** | https://angular.dev/guide/tailwind |
| **Angular Updates** | https://angular.dev/update |

**Key Requirements:**

1. **Performance**: Use `OnPush` change detection, avoid zone pollution, optimize computed signals
2. **Security**: Sanitize user input, use Angular's built-in XSS protection, avoid `bypassSecurityTrust*`
3. **Accessibility**: All components MUST meet WCAG 2.1 AA, use semantic HTML, proper ARIA labels
4. **Error Handling**: Implement global error handler, use `catchError` in HTTP calls
5. **Zoneless**: Project targets zoneless Angular - avoid Zone.js dependencies

---

## Naming Conventions (MANDATORY)

| Context | Convention | Example |
|---------|------------|---------|
| Go files | snake_case | `student_repository.go` |
| Go structs | PascalCase | `StudentService` |
| Go methods | PascalCase | `FindByTenantID` |
| Go variables | camelCase | `studentList` |
| TS files | kebab-case | `student-list.component.ts` |
| TS classes | PascalCase | `StudentListComponent` |
| TS methods | camelCase | `getStudents()` |
| DB tables | snake_case | `student_attendance` |
| DB columns | snake_case | `created_at` |
| API routes | kebab-case | `/api/v1/student-attendance` |
| JSON fields | camelCase | `{ "firstName": "John" }` |

---

## Database Rules

### Every Table MUST Have

```sql
CREATE TABLE {table_name} (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    -- ... other columns ...
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id)
);

-- RLS Policy (MANDATORY)
ALTER TABLE {table_name} ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON {table_name}
    USING (tenant_id = current_setting('app.tenant_id')::UUID);

-- Indexes
CREATE INDEX idx_{table_name}_tenant ON {table_name}(tenant_id);
```

### Migration Naming

```
{timestamp}_{description}.up.sql
{timestamp}_{description}.down.sql

Example:
20260123100000_create_students_table.up.sql
20260123100000_create_students_table.down.sql
```

---

## API Design Rules

### Endpoint Pattern

```
GET    /api/v1/{resource}           # List (paginated)
GET    /api/v1/{resource}/{id}      # Get single
POST   /api/v1/{resource}           # Create
PUT    /api/v1/{resource}/{id}      # Full update
PATCH  /api/v1/{resource}/{id}      # Partial update
DELETE /api/v1/{resource}/{id}      # Delete
```

### Response Format

```json
// Success (single)
{
  "data": { ... },
  "meta": { "requestId": "uuid" }
}

// Success (list)
{
  "data": [ ... ],
  "meta": {
    "pagination": {
      "cursor": "next_cursor_value",
      "hasMore": true,
      "limit": 20
    }
  }
}

// Error (RFC 7807)
{
  "type": "https://api.msls.com/errors/validation",
  "title": "Validation Error",
  "status": 400,
  "detail": "First name is required",
  "instance": "/api/v1/students",
  "errors": [
    { "field": "firstName", "message": "Required" }
  ]
}
```

---

## Testing Requirements

### Backend Tests

```go
// Unit test for service
func TestStudentService_Create(t *testing.T) {
    // Arrange
    mockRepo := mocks.NewMockStudentRepository(t)
    svc := NewStudentService(mockRepo)

    // Act
    result, err := svc.Create(ctx, dto)

    // Assert
    require.NoError(t, err)
    assert.Equal(t, expected, result)
}

// Integration test with testcontainers
func TestStudentRepository_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    // Use testcontainers for PostgreSQL
}
```

### Frontend Tests

```typescript
// Component test with Testing Library
describe('StudentListComponent', () => {
  it('should display students', async () => {
    await render(StudentListComponent, {
      providers: [
        { provide: StudentService, useValue: mockService }
      ]
    });

    expect(screen.getByText('John Doe')).toBeInTheDocument();
  });
});
```

---

## Security Rules (NEVER VIOLATE)

1. **NEVER** log sensitive data (passwords, tokens, PII)
2. **ALWAYS** validate input at handler level
3. **ALWAYS** use parameterized queries (GORM handles this)
4. **NEVER** expose internal errors to clients
5. **ALWAYS** check tenant_id matches current user's tenant
6. **ALWAYS** use HTTPS in production
7. **NEVER** store plain text passwords (use Argon2id)
8. **ALWAYS** sanitize user input before rendering (Angular does this)

---

## Agent Workflow

### Before Starting Any Task

1. Read `architecture.md` for technology decisions
2. Read the story file for requirements
3. Check existing code for patterns to follow
4. Plan your changes before coding

### While Implementing

1. Follow naming conventions exactly
2. Match existing code patterns
3. Write tests alongside code
4. Run linter after changes

### After Completing

1. Run all tests
2. Update story status
3. Commit with descriptive message

---

## Quick Reference

| Need | Solution |
|------|----------|
| State management | Angular Signals |
| HTTP client | Angular HttpClient |
| Forms | Reactive Forms |
| Routing | Lazy-loaded routes |
| Styling | Tailwind utility classes |
| Icons | Heroicons |
| Date handling | date-fns |
| Validation (FE) | class-validator |
| Validation (BE) | go-playground/validator |
| ORM | GORM (simple) + sqlc (complex) |
| Migrations | golang-migrate |
| Testing (BE) | testify + testcontainers |
| Testing (FE) | Vitest + Testing Library |
| E2E | Playwright |
