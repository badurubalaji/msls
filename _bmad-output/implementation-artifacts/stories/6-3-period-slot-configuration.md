# Story 6.3: Period Slot Configuration

**Epic:** 6 - Academic Structure & Timetable
**Status:** done
**Priority:** High

## User Story

As an **academic administrator**,
I want **to define period slots and school timing**,
So that **timetable can be created with correct time slots**.

## Acceptance Criteria

### AC1: Period Slot Definition
**Given** admin is on timetable settings
**When** defining period slots
**Then** they can enter: period number, start time, end time
**And** they can mark periods as: regular, short, assembly, break, lunch
**And** they can set: duration (minutes)

### AC2: Day Pattern Configuration
**Given** different day patterns exist
**When** configuring day types
**Then** they can create patterns: regular day, Saturday (half-day)
**And** each day type has its own period structure
**And** days of week can be assigned patterns

### AC3: Shift Configuration
**Given** school operates different shifts
**When** configuring shifts
**Then** they can define: morning shift, afternoon shift
**And** each shift has its own period slots
**And** sections can be assigned to shifts

## Technical Notes

### Database Schema
- `period_slots` table: id, tenant_id, branch_id, name, period_number, start_time, end_time, duration_minutes, slot_type, day_pattern_id, shift_id, is_active
- `day_patterns` table: id, tenant_id, name, code, description, is_active
- `day_pattern_assignments` table: id, tenant_id, day_of_week (0-6), day_pattern_id
- `shifts` table: id, tenant_id, branch_id, name, code, start_time, end_time, is_active

### API Endpoints
- GET/POST /api/v1/period-slots - List/Create period slots
- GET/PUT/DELETE /api/v1/period-slots/:id - Get/Update/Delete period slot
- GET/POST /api/v1/day-patterns - List/Create day patterns
- GET/PUT/DELETE /api/v1/day-patterns/:id - Get/Update/Delete day pattern
- GET/PUT /api/v1/day-pattern-assignments - List/Update day assignments
- GET/POST /api/v1/shifts - List/Create shifts
- GET/PUT/DELETE /api/v1/shifts/:id - Get/Update/Delete shift

### Frontend Components
- PeriodSlotsComponent - Manage period slots with visual time grid
- DayPatternsComponent - Configure day patterns
- ShiftsComponent - Configure school shifts

## Tasks

- [x] Task 1: Create database migrations for period slots, day patterns, and shifts
- [x] Task 2: Implement backend models and repositories
- [x] Task 3: Implement backend services with business logic
- [x] Task 4: Implement backend HTTP handlers and routes
- [x] Task 5: Create frontend models and services
- [x] Task 6: Create Period Slots management page
- [x] Task 7: Create Day Patterns configuration page
- [x] Task 8: Create Shifts management page
- [x] Task 9: Add navigation and integrate with academics module
- [x] Task 10: Add seed data for default period slots

## Dev Agent Record

### Started: 2026-01-29
### Completed: 2026-01-29

## File List

### Backend Files

#### New Files
- `msls-backend/migrations/000048_timetable_structure.up.sql` - Database tables for shifts, day_patterns, day_pattern_assignments, period_slots
- `msls-backend/migrations/000048_timetable_structure.down.sql` - Rollback migration
- `msls-backend/migrations/000049_timetable_seed_data.up.sql` - Seed data for timetable configuration
- `msls-backend/migrations/000049_timetable_seed_data.down.sql` - Rollback seed data
- `msls-backend/internal/pkg/database/models/timetable.go` - GORM models for Shift, DayPattern, DayPatternAssignment, PeriodSlot
- `msls-backend/internal/modules/timetable/errors.go` - Custom error types
- `msls-backend/internal/modules/timetable/dto.go` - Request/Response DTOs and mapping functions
- `msls-backend/internal/modules/timetable/repository.go` - Database operations
- `msls-backend/internal/modules/timetable/service.go` - Business logic
- `msls-backend/internal/modules/timetable/handler.go` - HTTP handlers and routes

#### Modified Files
- `msls-backend/cmd/api/main.go` - Register timetable module

### Frontend Files

#### New Files
- `msls-frontend/src/app/features/academics/timetable/timetable.model.ts` - TypeScript interfaces for Shift, DayPattern, PeriodSlot
- `msls-frontend/src/app/features/academics/timetable/timetable.service.ts` - API service for timetable operations
- `msls-frontend/src/app/features/academics/timetable/shifts/shifts.component.ts` - Shifts management component
- `msls-frontend/src/app/features/academics/timetable/day-patterns/day-patterns.component.ts` - Day patterns and week schedule component
- `msls-frontend/src/app/features/academics/timetable/period-slots/period-slots.component.ts` - Period slots management with timeline view

#### Modified Files
- `msls-frontend/src/app/features/academics/timetable/timetable.component.ts` - Updated to landing page with links
- `msls-frontend/src/app/features/academics/academics.routes.ts` - Added timetable sub-routes
- `msls-frontend/src/app/layouts/nav-config.ts` - Added timetable sub-navigation

## Change Log

- 2026-01-29: Story created and implementation started
- 2026-01-29: Backend implementation completed (models, repository, service, handler)
- 2026-01-29: Frontend implementation completed (models, service, components)
- 2026-01-29: Navigation and routing updated
- 2026-01-29: Story marked as done

## Features Implemented

### Shifts Management
- CRUD operations for school shifts
- Configure start/end times
- Visual time display with AM/PM formatting
- Duration calculation
- Status toggle (active/inactive)

### Day Patterns
- CRUD operations for day patterns (Regular, Half-Day, Assembly, etc.)
- Total periods configuration
- Week schedule assignment (assign patterns to weekdays)
- Working day toggle for each day
- Branch-specific assignments

### Period Slots
- CRUD operations for period slots
- Multiple slot types: regular, short, assembly, break, lunch, activity, zero_period
- Visual timeline view with color-coded slots
- Table view with full details
- Filter by branch, day pattern, shift, and slot type
- Automatic duration calculation from start/end times
- Display order management
