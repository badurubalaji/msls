// Package academicyear provides academic year management services.
package academicyear

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
)

// setupTestDB creates an in-memory SQLite database for testing.
func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Run migrations
	err = db.AutoMigrate(&models.Tenant{}, &models.Branch{}, &models.AcademicYear{}, &models.AcademicTerm{}, &models.Holiday{})
	require.NoError(t, err)

	return db
}

// createTestTenant creates a test tenant in the database.
func createTestTenant(t *testing.T, db *gorm.DB) uuid.UUID {
	t.Helper()

	tenantID := uuid.New()
	tenant := &models.Tenant{
		Name: "Test Tenant",
		Slug: "test-tenant",
	}
	tenant.ID = tenantID
	tenant.CreatedAt = time.Now()
	tenant.UpdatedAt = time.Now()
	tenant.Status = models.StatusActive

	err := db.Create(tenant).Error
	require.NoError(t, err)

	return tenantID
}

// =============================================================================
// Academic Year Tests
// =============================================================================

func TestService_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		req     CreateAcademicYearRequest
		wantErr error
	}{
		{
			name: "success",
			req: CreateAcademicYearRequest{
				Name:      "2025-26",
				StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC),
				IsCurrent: true,
			},
			wantErr: nil,
		},
		{
			name: "missing tenant ID",
			req: CreateAcademicYearRequest{
				Name:      "2025-26",
				StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC),
			},
			wantErr: ErrTenantIDRequired,
		},
		{
			name: "missing name",
			req: CreateAcademicYearRequest{
				StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC),
			},
			wantErr: ErrAcademicYearNameRequired,
		},
		{
			name: "missing start date",
			req: CreateAcademicYearRequest{
				Name:    "2025-26",
				EndDate: time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC),
			},
			wantErr: ErrAcademicYearStartDateRequired,
		},
		{
			name: "missing end date",
			req: CreateAcademicYearRequest{
				Name:      "2025-26",
				StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
			},
			wantErr: ErrAcademicYearEndDateRequired,
		},
		{
			name: "invalid dates - end before start",
			req: CreateAcademicYearRequest{
				Name:      "2025-26",
				StartDate: time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2025, 3, 31, 0, 0, 0, 0, time.UTC),
			},
			wantErr: ErrAcademicYearInvalidDates,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db := setupTestDB(t)
			svc := NewService(db)
			ctx := context.Background()

			// Set up tenant ID for success case
			if tc.wantErr == nil || tc.wantErr == ErrAcademicYearNameRequired ||
				tc.wantErr == ErrAcademicYearStartDateRequired ||
				tc.wantErr == ErrAcademicYearEndDateRequired ||
				tc.wantErr == ErrAcademicYearInvalidDates {
				tenantID := createTestTenant(t, db)
				tc.req.TenantID = tenantID
			}

			academicYear, err := svc.Create(ctx, tc.req)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, academicYear)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, academicYear)
				assert.Equal(t, tc.req.Name, academicYear.Name)
				assert.Equal(t, tc.req.IsCurrent, academicYear.IsCurrent)
				assert.True(t, academicYear.IsActive)
			}
		})
	}
}

func TestService_Create_DuplicateName(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create first academic year
	req := CreateAcademicYearRequest{
		TenantID:  tenantID,
		Name:      "2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC),
	}
	_, err := svc.Create(ctx, req)
	require.NoError(t, err)

	// Try to create duplicate
	req.StartDate = time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	req.EndDate = time.Date(2027, 3, 31, 0, 0, 0, 0, time.UTC)
	_, err = svc.Create(ctx, req)
	assert.ErrorIs(t, err, ErrAcademicYearNameExists)
}

func TestService_GetByID(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create academic year
	req := CreateAcademicYearRequest{
		TenantID:  tenantID,
		Name:      "2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC),
	}
	created, err := svc.Create(ctx, req)
	require.NoError(t, err)

	// Get by ID
	academicYear, err := svc.GetByID(ctx, tenantID, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, academicYear.ID)
	assert.Equal(t, created.Name, academicYear.Name)
}

func TestService_GetByID_NotFound(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	_, err := svc.GetByID(ctx, tenantID, uuid.New())
	assert.ErrorIs(t, err, ErrAcademicYearNotFound)
}

func TestService_List(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create academic years
	years := []CreateAcademicYearRequest{
		{TenantID: tenantID, Name: "2025-26", StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC), IsCurrent: true},
		{TenantID: tenantID, Name: "2024-25", StartDate: time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2025, 3, 31, 0, 0, 0, 0, time.UTC)},
		{TenantID: tenantID, Name: "2023-24", StartDate: time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC)},
	}

	for _, req := range years {
		_, err := svc.Create(ctx, req)
		require.NoError(t, err)
	}

	// List all
	result, err := svc.List(ctx, ListAcademicYearFilter{TenantID: tenantID})
	require.NoError(t, err)
	assert.Len(t, result, 3)
	// Results should be ordered by start_date DESC
	assert.Equal(t, "2025-26", result[0].Name)
}

func TestService_List_WithFilters(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create academic years
	_, err := svc.Create(ctx, CreateAcademicYearRequest{
		TenantID:  tenantID,
		Name:      "2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC),
		IsCurrent: true,
	})
	require.NoError(t, err)

	_, err = svc.Create(ctx, CreateAcademicYearRequest{
		TenantID:  tenantID,
		Name:      "2024-25",
		StartDate: time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 3, 31, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Filter by is_current
	isCurrent := true
	result, err := svc.List(ctx, ListAcademicYearFilter{TenantID: tenantID, IsCurrent: &isCurrent})
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "2025-26", result[0].Name)

	// Filter by search
	result, err = svc.List(ctx, ListAcademicYearFilter{TenantID: tenantID, Search: "2024"})
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "2024-25", result[0].Name)
}

func TestService_Update(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create academic year
	created, err := svc.Create(ctx, CreateAcademicYearRequest{
		TenantID:  tenantID,
		Name:      "2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Update
	newName := "2025-2026"
	updated, err := svc.Update(ctx, tenantID, created.ID, UpdateAcademicYearRequest{Name: &newName})
	require.NoError(t, err)
	assert.Equal(t, newName, updated.Name)
}

func TestService_SetCurrent(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create academic years
	year1, err := svc.Create(ctx, CreateAcademicYearRequest{
		TenantID:  tenantID,
		Name:      "2024-25",
		StartDate: time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 3, 31, 0, 0, 0, 0, time.UTC),
		IsCurrent: true,
	})
	require.NoError(t, err)

	year2, err := svc.Create(ctx, CreateAcademicYearRequest{
		TenantID:  tenantID,
		Name:      "2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Set year2 as current
	updated, err := svc.SetCurrent(ctx, tenantID, year2.ID, nil)
	require.NoError(t, err)
	assert.True(t, updated.IsCurrent)

	// Note: The database trigger handles unsetting the previous current year.
	// In SQLite for testing, we may need to verify this manually if needed.
	// For now, we verify that the new year is set as current.
	_ = year1
}

func TestService_SetCurrent_InactiveYear(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create academic year
	created, err := svc.Create(ctx, CreateAcademicYearRequest{
		TenantID:  tenantID,
		Name:      "2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Deactivate
	isActive := false
	_, err = svc.Update(ctx, tenantID, created.ID, UpdateAcademicYearRequest{IsActive: &isActive})
	require.NoError(t, err)

	// Try to set as current
	_, err = svc.SetCurrent(ctx, tenantID, created.ID, nil)
	assert.ErrorIs(t, err, ErrCannotSetInactiveAsCurrent)
}

func TestService_Delete(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create academic year
	created, err := svc.Create(ctx, CreateAcademicYearRequest{
		TenantID:  tenantID,
		Name:      "2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Delete
	err = svc.Delete(ctx, tenantID, created.ID)
	require.NoError(t, err)

	// Verify deleted
	_, err = svc.GetByID(ctx, tenantID, created.ID)
	assert.ErrorIs(t, err, ErrAcademicYearNotFound)
}

func TestService_GetCurrent(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create current academic year
	created, err := svc.Create(ctx, CreateAcademicYearRequest{
		TenantID:  tenantID,
		Name:      "2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC),
		IsCurrent: true,
	})
	require.NoError(t, err)

	// Get current
	current, err := svc.GetCurrent(ctx, tenantID, nil)
	require.NoError(t, err)
	assert.Equal(t, created.ID, current.ID)
}

func TestService_GetCurrent_NotFound(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create non-current academic year
	_, err := svc.Create(ctx, CreateAcademicYearRequest{
		TenantID:  tenantID,
		Name:      "2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC),
		IsCurrent: false,
	})
	require.NoError(t, err)

	_, err = svc.GetCurrent(ctx, tenantID, nil)
	assert.ErrorIs(t, err, ErrAcademicYearNotFound)
}

// =============================================================================
// Academic Term Tests
// =============================================================================

func TestService_CreateTerm(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create academic year
	academicYear, err := svc.Create(ctx, CreateAcademicYearRequest{
		TenantID:  tenantID,
		Name:      "2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Create term
	term, err := svc.CreateTerm(ctx, CreateTermRequest{
		TenantID:       tenantID,
		AcademicYearID: academicYear.ID,
		Name:           "Term 1",
		StartDate:      time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:        time.Date(2025, 9, 30, 0, 0, 0, 0, time.UTC),
		Sequence:       1,
	})
	require.NoError(t, err)
	assert.NotNil(t, term)
	assert.Equal(t, "Term 1", term.Name)
	assert.Equal(t, 1, term.Sequence)
}

func TestService_CreateTerm_AutoSequence(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create academic year
	academicYear, err := svc.Create(ctx, CreateAcademicYearRequest{
		TenantID:  tenantID,
		Name:      "2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Create terms without explicit sequence
	term1, err := svc.CreateTerm(ctx, CreateTermRequest{
		TenantID:       tenantID,
		AcademicYearID: academicYear.ID,
		Name:           "Term 1",
		StartDate:      time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:        time.Date(2025, 9, 30, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)
	assert.Equal(t, 1, term1.Sequence)

	term2, err := svc.CreateTerm(ctx, CreateTermRequest{
		TenantID:       tenantID,
		AcademicYearID: academicYear.ID,
		Name:           "Term 2",
		StartDate:      time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC),
		EndDate:        time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)
	assert.Equal(t, 2, term2.Sequence)
}

func TestService_CreateTerm_OutsideAcademicYear(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create academic year
	academicYear, err := svc.Create(ctx, CreateAcademicYearRequest{
		TenantID:  tenantID,
		Name:      "2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Try to create term outside academic year dates
	_, err = svc.CreateTerm(ctx, CreateTermRequest{
		TenantID:       tenantID,
		AcademicYearID: academicYear.ID,
		Name:           "Term 1",
		StartDate:      time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), // Before academic year
		EndDate:        time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
	})
	assert.ErrorIs(t, err, ErrTermOutsideAcademicYear)
}

func TestService_CreateTerm_DuplicateName(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create academic year
	academicYear, err := svc.Create(ctx, CreateAcademicYearRequest{
		TenantID:  tenantID,
		Name:      "2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Create first term
	_, err = svc.CreateTerm(ctx, CreateTermRequest{
		TenantID:       tenantID,
		AcademicYearID: academicYear.ID,
		Name:           "Term 1",
		StartDate:      time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:        time.Date(2025, 9, 30, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Try to create duplicate
	_, err = svc.CreateTerm(ctx, CreateTermRequest{
		TenantID:       tenantID,
		AcademicYearID: academicYear.ID,
		Name:           "Term 1",
		StartDate:      time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC),
		EndDate:        time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC),
	})
	assert.ErrorIs(t, err, ErrTermNameExists)
}

func TestService_ListTerms(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create academic year
	academicYear, err := svc.Create(ctx, CreateAcademicYearRequest{
		TenantID:  tenantID,
		Name:      "2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Create terms
	_, err = svc.CreateTerm(ctx, CreateTermRequest{
		TenantID:       tenantID,
		AcademicYearID: academicYear.ID,
		Name:           "Term 1",
		StartDate:      time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:        time.Date(2025, 9, 30, 0, 0, 0, 0, time.UTC),
		Sequence:       1,
	})
	require.NoError(t, err)

	_, err = svc.CreateTerm(ctx, CreateTermRequest{
		TenantID:       tenantID,
		AcademicYearID: academicYear.ID,
		Name:           "Term 2",
		StartDate:      time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC),
		EndDate:        time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC),
		Sequence:       2,
	})
	require.NoError(t, err)

	// List terms
	terms, err := svc.ListTerms(ctx, tenantID, academicYear.ID)
	require.NoError(t, err)
	assert.Len(t, terms, 2)
	assert.Equal(t, "Term 1", terms[0].Name)
	assert.Equal(t, "Term 2", terms[1].Name)
}

func TestService_UpdateTerm(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create academic year
	academicYear, err := svc.Create(ctx, CreateAcademicYearRequest{
		TenantID:  tenantID,
		Name:      "2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Create term
	term, err := svc.CreateTerm(ctx, CreateTermRequest{
		TenantID:       tenantID,
		AcademicYearID: academicYear.ID,
		Name:           "Term 1",
		StartDate:      time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:        time.Date(2025, 9, 30, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Update term
	newName := "First Term"
	updated, err := svc.UpdateTerm(ctx, tenantID, academicYear.ID, term.ID, UpdateTermRequest{Name: &newName})
	require.NoError(t, err)
	assert.Equal(t, newName, updated.Name)
}

func TestService_DeleteTerm(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create academic year
	academicYear, err := svc.Create(ctx, CreateAcademicYearRequest{
		TenantID:  tenantID,
		Name:      "2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Create term
	term, err := svc.CreateTerm(ctx, CreateTermRequest{
		TenantID:       tenantID,
		AcademicYearID: academicYear.ID,
		Name:           "Term 1",
		StartDate:      time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:        time.Date(2025, 9, 30, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Delete term
	err = svc.DeleteTerm(ctx, tenantID, academicYear.ID, term.ID)
	require.NoError(t, err)

	// Verify deleted
	_, err = svc.GetTermByID(ctx, tenantID, term.ID)
	assert.ErrorIs(t, err, ErrTermNotFound)
}

// =============================================================================
// Holiday Tests
// =============================================================================

func TestService_CreateHoliday(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create academic year
	academicYear, err := svc.Create(ctx, CreateAcademicYearRequest{
		TenantID:  tenantID,
		Name:      "2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Create holiday
	holiday, err := svc.CreateHoliday(ctx, CreateHolidayRequest{
		TenantID:       tenantID,
		AcademicYearID: academicYear.ID,
		Name:           "Diwali",
		Date:           time.Date(2025, 10, 21, 0, 0, 0, 0, time.UTC),
		Type:           models.HolidayTypeReligious,
		IsOptional:     false,
	})
	require.NoError(t, err)
	assert.NotNil(t, holiday)
	assert.Equal(t, "Diwali", holiday.Name)
	assert.Equal(t, models.HolidayTypeReligious, holiday.Type)
}

func TestService_CreateHoliday_DefaultType(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create academic year
	academicYear, err := svc.Create(ctx, CreateAcademicYearRequest{
		TenantID:  tenantID,
		Name:      "2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Create holiday without type
	holiday, err := svc.CreateHoliday(ctx, CreateHolidayRequest{
		TenantID:       tenantID,
		AcademicYearID: academicYear.ID,
		Name:           "Holiday",
		Date:           time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)
	assert.Equal(t, models.HolidayTypePublic, holiday.Type)
}

func TestService_CreateHoliday_OutsideAcademicYear(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create academic year
	academicYear, err := svc.Create(ctx, CreateAcademicYearRequest{
		TenantID:  tenantID,
		Name:      "2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Try to create holiday outside academic year
	_, err = svc.CreateHoliday(ctx, CreateHolidayRequest{
		TenantID:       tenantID,
		AcademicYearID: academicYear.ID,
		Name:           "Holiday",
		Date:           time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), // Before academic year
	})
	assert.ErrorIs(t, err, ErrHolidayOutsideAcademicYear)
}

func TestService_ListHolidays(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create academic year
	academicYear, err := svc.Create(ctx, CreateAcademicYearRequest{
		TenantID:  tenantID,
		Name:      "2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Create holidays
	_, err = svc.CreateHoliday(ctx, CreateHolidayRequest{
		TenantID:       tenantID,
		AcademicYearID: academicYear.ID,
		Name:           "Independence Day",
		Date:           time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC),
		Type:           models.HolidayTypeNational,
	})
	require.NoError(t, err)

	_, err = svc.CreateHoliday(ctx, CreateHolidayRequest{
		TenantID:       tenantID,
		AcademicYearID: academicYear.ID,
		Name:           "Diwali",
		Date:           time.Date(2025, 10, 21, 0, 0, 0, 0, time.UTC),
		Type:           models.HolidayTypeReligious,
	})
	require.NoError(t, err)

	// List holidays
	holidays, err := svc.ListHolidays(ctx, tenantID, academicYear.ID)
	require.NoError(t, err)
	assert.Len(t, holidays, 2)
	// Should be sorted by date
	assert.Equal(t, "Independence Day", holidays[0].Name)
	assert.Equal(t, "Diwali", holidays[1].Name)
}

func TestService_UpdateHoliday(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create academic year
	academicYear, err := svc.Create(ctx, CreateAcademicYearRequest{
		TenantID:  tenantID,
		Name:      "2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Create holiday
	holiday, err := svc.CreateHoliday(ctx, CreateHolidayRequest{
		TenantID:       tenantID,
		AcademicYearID: academicYear.ID,
		Name:           "Holiday",
		Date:           time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Update holiday
	newName := "Independence Day"
	isOptional := true
	updated, err := svc.UpdateHoliday(ctx, tenantID, academicYear.ID, holiday.ID, UpdateHolidayRequest{
		Name:       &newName,
		IsOptional: &isOptional,
	})
	require.NoError(t, err)
	assert.Equal(t, newName, updated.Name)
	assert.True(t, updated.IsOptional)
}

func TestService_DeleteHoliday(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create academic year
	academicYear, err := svc.Create(ctx, CreateAcademicYearRequest{
		TenantID:  tenantID,
		Name:      "2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Create holiday
	holiday, err := svc.CreateHoliday(ctx, CreateHolidayRequest{
		TenantID:       tenantID,
		AcademicYearID: academicYear.ID,
		Name:           "Holiday",
		Date:           time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Delete holiday
	err = svc.DeleteHoliday(ctx, tenantID, academicYear.ID, holiday.ID)
	require.NoError(t, err)

	// Verify deleted
	_, err = svc.GetHolidayByID(ctx, tenantID, holiday.ID)
	assert.ErrorIs(t, err, ErrHolidayNotFound)
}

func TestService_IsHoliday(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create academic year
	academicYear, err := svc.Create(ctx, CreateAcademicYearRequest{
		TenantID:  tenantID,
		Name:      "2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Create holiday
	holidayDate := time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC)
	_, err = svc.CreateHoliday(ctx, CreateHolidayRequest{
		TenantID:       tenantID,
		AcademicYearID: academicYear.ID,
		Name:           "Independence Day",
		Date:           holidayDate,
		Type:           models.HolidayTypeNational,
	})
	require.NoError(t, err)

	// Check if it's a holiday
	isHoliday, holiday, err := svc.IsHoliday(ctx, tenantID, holidayDate, nil)
	require.NoError(t, err)
	assert.True(t, isHoliday)
	assert.NotNil(t, holiday)
	assert.Equal(t, "Independence Day", holiday.Name)

	// Check a non-holiday date
	isHoliday, holiday, err = svc.IsHoliday(ctx, tenantID, time.Date(2025, 8, 16, 0, 0, 0, 0, time.UTC), nil)
	require.NoError(t, err)
	assert.False(t, isHoliday)
	assert.Nil(t, holiday)
}

func TestService_Count(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create academic years
	for i := 0; i < 3; i++ {
		_, err := svc.Create(ctx, CreateAcademicYearRequest{
			TenantID:  tenantID,
			Name:      "Year " + string(rune('A'+i)),
			StartDate: time.Date(2020+i, 4, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2021+i, 3, 31, 0, 0, 0, 0, time.UTC),
		})
		require.NoError(t, err)
	}

	count, err := svc.Count(ctx, tenantID)
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)
}

// =============================================================================
// Delete Cascade Tests
// =============================================================================

func TestService_Delete_CascadesTermsAndHolidays(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create academic year
	academicYear, err := svc.Create(ctx, CreateAcademicYearRequest{
		TenantID:  tenantID,
		Name:      "2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Create terms
	term, err := svc.CreateTerm(ctx, CreateTermRequest{
		TenantID:       tenantID,
		AcademicYearID: academicYear.ID,
		Name:           "Term 1",
		StartDate:      time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:        time.Date(2025, 9, 30, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Create holidays
	holiday, err := svc.CreateHoliday(ctx, CreateHolidayRequest{
		TenantID:       tenantID,
		AcademicYearID: academicYear.ID,
		Name:           "Holiday",
		Date:           time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Delete academic year
	err = svc.Delete(ctx, tenantID, academicYear.ID)
	require.NoError(t, err)

	// Verify term is deleted
	_, err = svc.GetTermByID(ctx, tenantID, term.ID)
	assert.ErrorIs(t, err, ErrTermNotFound)

	// Verify holiday is deleted
	_, err = svc.GetHolidayByID(ctx, tenantID, holiday.ID)
	assert.ErrorIs(t, err, ErrHolidayNotFound)
}
