# Story 6.1: Class & Section Management

## Status: review

## Story

As an **academic administrator**,
I want **to configure classes and sections**,
So that **student groupings are properly defined**.

## Acceptance Criteria

### AC1: Class Creation
- **Given** admin is on class configuration
- **When** creating a class
- **Then** they can enter: class name (e.g., "Class 10"), numeric order
- **And** they can set: applicable academic years
- **And** they can set: streams (if senior classes - Science, Commerce, Arts)

### AC2: Section Management
- **Given** a class exists
- **When** adding sections
- **Then** they can enter: section name (A, B, C)
- **And** they can set: max capacity
- **And** they can assign: class teacher
- **And** sections are linked to the class

### AC3: Hierarchical View
- **Given** class-section setup is complete
- **When** viewing the structure
- **Then** they see hierarchical view: Class → Sections
- **And** each section shows: student count, teacher assigned
- **And** capacity utilization is displayed

## Technical Design

### Database Schema

```sql
-- Classes table
CREATE TABLE classes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    branch_id UUID REFERENCES branches(id),
    name VARCHAR(100) NOT NULL,
    code VARCHAR(20),
    numeric_order INT NOT NULL,
    description TEXT,
    has_streams BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    UNIQUE(tenant_id, branch_id, name),
    UNIQUE(tenant_id, branch_id, numeric_order)
);

-- Streams for senior classes
CREATE TABLE streams (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    name VARCHAR(100) NOT NULL,
    code VARCHAR(20),
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(tenant_id, name)
);

-- Class-Stream mapping (which streams available for which class)
CREATE TABLE class_streams (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    stream_id UUID NOT NULL REFERENCES streams(id) ON DELETE CASCADE,
    is_active BOOLEAN DEFAULT TRUE,
    UNIQUE(class_id, stream_id)
);

-- Sections table
CREATE TABLE sections (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    academic_year_id UUID NOT NULL REFERENCES academic_years(id),
    stream_id UUID REFERENCES streams(id),
    name VARCHAR(20) NOT NULL,
    max_capacity INT DEFAULT 40,
    class_teacher_id UUID REFERENCES staff(id),
    room_number VARCHAR(50),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    UNIQUE(tenant_id, class_id, academic_year_id, name, stream_id)
);

-- Class-Academic Year mapping
CREATE TABLE class_academic_years (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    academic_year_id UUID NOT NULL REFERENCES academic_years(id) ON DELETE CASCADE,
    is_active BOOLEAN DEFAULT TRUE,
    UNIQUE(class_id, academic_year_id)
);
```

### Backend Structure

```
internal/modules/academic/
├── class/
│   ├── handler.go      # HTTP handlers
│   ├── service.go      # Business logic
│   ├── repository.go   # Database operations
│   ├── dto.go          # Request/Response DTOs
│   └── errors.go       # Domain errors
├── stream/
│   ├── handler.go
│   ├── service.go
│   ├── repository.go
│   └── dto.go
└── section/
    ├── handler.go
    ├── service.go
    ├── repository.go
    └── dto.go
```

### API Endpoints

```
# Classes
GET    /api/v1/classes                    - List all classes
POST   /api/v1/classes                    - Create class
GET    /api/v1/classes/:id                - Get class details
PUT    /api/v1/classes/:id                - Update class
DELETE /api/v1/classes/:id                - Delete class
GET    /api/v1/classes/:id/sections       - Get sections for a class
POST   /api/v1/classes/:id/streams        - Assign streams to class
GET    /api/v1/classes/structure          - Get hierarchical class-section structure

# Streams
GET    /api/v1/streams                    - List all streams
POST   /api/v1/streams                    - Create stream
PUT    /api/v1/streams/:id                - Update stream
DELETE /api/v1/streams/:id                - Delete stream

# Sections
GET    /api/v1/sections                   - List all sections (with filters)
POST   /api/v1/sections                   - Create section
GET    /api/v1/sections/:id               - Get section details
PUT    /api/v1/sections/:id               - Update section
DELETE /api/v1/sections/:id               - Delete section
PUT    /api/v1/sections/:id/class-teacher - Assign class teacher
GET    /api/v1/sections/:id/students      - Get students in section
```

### Frontend Structure

```
src/app/features/academic/
├── academic.routes.ts
├── academic.model.ts
├── services/
│   ├── class.service.ts
│   ├── stream.service.ts
│   └── section.service.ts
└── pages/
    ├── class-management/
    │   └── class-management.component.ts
    ├── stream-management/
    │   └── stream-management.component.ts
    └── section-management/
        └── section-management.component.ts
```

## Tasks

- [x] Task 1: Create database migrations for classes, streams, sections
- [x] Task 2: Implement backend models and repositories
- [x] Task 3: Implement backend services with business logic
- [x] Task 4: Implement backend HTTP handlers and routes
- [x] Task 5: Add class level field (nursery, primary, middle, secondary, senior_secondary) per PRD
- [x] Task 6: Create frontend models and services
- [x] Task 7: Create Class Management page with full CRUD
- [x] Task 8: Create Stream Management page with full CRUD
- [x] Task 9: Create Section Management page with class teacher assignment
- [x] Task 10: Create hierarchical structure view (Class → Sections tree)
- [x] Task 11: Backend unit testing and validation

## Dev Notes

- Classes are templates that exist across academic years
- Sections are created per academic year for each class
- Streams are optional and only for senior classes (11th, 12th typically)
- Class teacher is assigned per section per academic year
- Student count is computed from enrollments table

## Dev Agent Record

### Implementation Plan
- Task 10: Add comprehensive unit tests for the academic module
  - Test error definitions for class, section, and stream errors
  - Test DTO conversion functions (ClassToResponse, SectionToResponse, StreamToResponse)
  - Test request DTOs (CreateClassRequest, UpdateClassRequest, CreateSectionRequest, etc.)
  - Test filter structs (ClassFilter, SectionFilter, StreamFilter)
  - Test response structs (ClassListResponse, SectionListResponse, StreamListResponse)
  - Test structure response DTOs for hierarchical view

### Completion Notes
- Created 38 unit tests covering:
  - Error definitions validation (class, section, stream, general errors)
  - DTO conversion functions with various scenarios (with/without optional fields, timestamps)
  - Request DTO field validation
  - Response DTO structure validation
  - Filter struct field validation
  - All tests pass (38/38)
  - Code compiles with go build
  - go vet passes with no issues

## File List

### New Files
- `msls-backend/internal/modules/academic/service_test.go` - Unit tests for academic module
- `msls-backend/migrations/000046_class_level.up.sql` - Migration to add class level field
- `msls-backend/migrations/000046_class_level.down.sql` - Down migration for class level
- `msls-frontend/src/app/features/academics/academic.model.ts` - TypeScript interfaces
- `msls-frontend/src/app/features/academics/services/class.service.ts` - Class API service
- `msls-frontend/src/app/features/academics/services/section.service.ts` - Section API service
- `msls-frontend/src/app/features/academics/services/stream.service.ts` - Stream API service
- `msls-frontend/src/app/features/academics/classes/classes.component.ts` - Full CRUD Classes page
- `msls-frontend/src/app/features/academics/sections/sections.component.ts` - Full CRUD Sections page
- `msls-frontend/src/app/features/academics/streams/streams.component.ts` - Full CRUD Streams page
- `msls-frontend/src/app/features/academics/structure/structure.component.ts` - Hierarchical structure view

### Existing Files (modified)
- `msls-backend/internal/modules/academic/handler.go` - Added level query parameter
- `msls-backend/internal/modules/academic/service.go` - Added level handling in Create/Update
- `msls-backend/internal/modules/academic/repository.go` - Added level filter
- `msls-backend/internal/modules/academic/dto.go` - Added Level field to DTOs
- `msls-backend/internal/modules/academic/errors.go`
- `msls-backend/internal/pkg/database/models/assignment.go` - Added ClassLevel type and Level field
- `msls-backend/migrations/000041_academic_structure.up.sql`
- `msls-backend/migrations/000041_academic_structure.down.sql`
- `msls-frontend/src/app/features/academics/academics.routes.ts` - Added all routes

## Change Log

- 2026-01-29: Task 10 - Added comprehensive unit tests for academic module (38 tests covering DTOs, errors, filters, and responses)
- 2026-01-29: Task 5 - Added class level field with migration 000046_class_level (nursery, primary, middle, secondary, senior_secondary)
- 2026-01-29: Tasks 6-10 - Implemented complete frontend:
  - Created academic.model.ts with all TypeScript interfaces (Class, Section, Stream, ClassLevel, etc.)
  - Created class.service.ts, section.service.ts, stream.service.ts for API operations
  - Created ClassesComponent with full CRUD, level filtering, and stream assignment
  - Created SectionsComponent with capacity visualization and class teacher assignment
  - Created StreamsComponent with full CRUD functionality
  - Created StructureComponent with hierarchical Class → Sections view, capacity utilization display
  - Updated academics.routes.ts with all routes
