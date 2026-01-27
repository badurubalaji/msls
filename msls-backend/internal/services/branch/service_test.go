// Package branch provides branch management services.
package branch

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
	err = db.AutoMigrate(&models.Tenant{}, &models.Branch{})
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

func TestService_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		req     CreateRequest
		wantErr error
	}{
		{
			name: "success",
			req: CreateRequest{
				Code:      "MAIN",
				Name:      "Main Branch",
				City:      "Mumbai",
				State:     "Maharashtra",
				Country:   "India",
				Timezone:  "Asia/Kolkata",
				IsPrimary: true,
			},
			wantErr: nil,
		},
		{
			name: "missing tenant ID",
			req: CreateRequest{
				Code: "MAIN",
				Name: "Main Branch",
			},
			wantErr: ErrTenantIDRequired,
		},
		{
			name: "missing name",
			req: CreateRequest{
				Code: "MAIN",
			},
			wantErr: ErrBranchNameRequired,
		},
		{
			name: "missing code",
			req: CreateRequest{
				Name: "Main Branch",
			},
			wantErr: ErrBranchCodeRequired,
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
			if tc.wantErr == nil || tc.wantErr == ErrBranchNameRequired || tc.wantErr == ErrBranchCodeRequired {
				tenantID := createTestTenant(t, db)
				tc.req.TenantID = tenantID
			}

			branch, err := svc.Create(ctx, tc.req)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, branch)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, branch)
				assert.Equal(t, tc.req.Code, branch.Code)
				assert.Equal(t, tc.req.Name, branch.Name)
				assert.Equal(t, tc.req.IsPrimary, branch.IsPrimary)
			}
		})
	}
}

func TestService_Create_DuplicateCode(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create first branch
	req := CreateRequest{
		TenantID: tenantID,
		Code:     "MAIN",
		Name:     "Main Branch",
	}
	_, err := svc.Create(ctx, req)
	require.NoError(t, err)

	// Try to create duplicate
	req.Name = "Another Branch"
	_, err = svc.Create(ctx, req)
	assert.ErrorIs(t, err, ErrBranchCodeExists)
}

func TestService_Create_InvalidTimezone(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	req := CreateRequest{
		TenantID: tenantID,
		Code:     "MAIN",
		Name:     "Main Branch",
		Timezone: "Invalid/Timezone",
	}
	_, err := svc.Create(ctx, req)
	assert.ErrorContains(t, err, "invalid timezone")
}

func TestService_GetByID(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create branch
	req := CreateRequest{
		TenantID: tenantID,
		Code:     "MAIN",
		Name:     "Main Branch",
	}
	created, err := svc.Create(ctx, req)
	require.NoError(t, err)

	// Get by ID
	branch, err := svc.GetByID(ctx, tenantID, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, branch.ID)
	assert.Equal(t, created.Name, branch.Name)
}

func TestService_GetByID_NotFound(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	_, err := svc.GetByID(ctx, tenantID, uuid.New())
	assert.ErrorIs(t, err, ErrBranchNotFound)
}

func TestService_GetByCode(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create branch
	req := CreateRequest{
		TenantID: tenantID,
		Code:     "MAIN",
		Name:     "Main Branch",
	}
	created, err := svc.Create(ctx, req)
	require.NoError(t, err)

	// Get by code
	branch, err := svc.GetByCode(ctx, tenantID, "MAIN")
	require.NoError(t, err)
	assert.Equal(t, created.ID, branch.ID)
}

func TestService_List(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create branches
	branches := []CreateRequest{
		{TenantID: tenantID, Code: "MAIN", Name: "Main Branch", IsPrimary: true},
		{TenantID: tenantID, Code: "SOUTH", Name: "South Branch"},
		{TenantID: tenantID, Code: "NORTH", Name: "North Branch"},
	}

	for _, req := range branches {
		_, err := svc.Create(ctx, req)
		require.NoError(t, err)
	}

	// List all
	result, err := svc.List(ctx, ListFilter{TenantID: tenantID})
	require.NoError(t, err)
	assert.Len(t, result, 3)
	// Primary should be first
	assert.True(t, result[0].IsPrimary)
}

func TestService_List_WithFilter(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create branches
	_, err := svc.Create(ctx, CreateRequest{TenantID: tenantID, Code: "MAIN", Name: "Main Branch", IsPrimary: true})
	require.NoError(t, err)
	_, err = svc.Create(ctx, CreateRequest{TenantID: tenantID, Code: "SOUTH", Name: "South Branch"})
	require.NoError(t, err)

	// Filter by primary
	isPrimary := true
	result, err := svc.List(ctx, ListFilter{TenantID: tenantID, IsPrimary: &isPrimary})
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "MAIN", result[0].Code)
}

func TestService_Update(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create branch
	created, err := svc.Create(ctx, CreateRequest{TenantID: tenantID, Code: "MAIN", Name: "Main Branch"})
	require.NoError(t, err)

	// Update
	newName := "Updated Branch"
	updated, err := svc.Update(ctx, tenantID, created.ID, UpdateRequest{Name: &newName})
	require.NoError(t, err)
	assert.Equal(t, newName, updated.Name)
}

func TestService_SetPrimary(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create branches
	main, err := svc.Create(ctx, CreateRequest{TenantID: tenantID, Code: "MAIN", Name: "Main Branch", IsPrimary: true})
	require.NoError(t, err)
	south, err := svc.Create(ctx, CreateRequest{TenantID: tenantID, Code: "SOUTH", Name: "South Branch"})
	require.NoError(t, err)

	// Set south as primary
	updated, err := svc.SetPrimary(ctx, tenantID, south.ID, nil)
	require.NoError(t, err)
	assert.True(t, updated.IsPrimary)

	// Verify main is no longer primary
	main, err = svc.GetByID(ctx, tenantID, main.ID)
	require.NoError(t, err)
	assert.False(t, main.IsPrimary)
}

func TestService_SetStatus(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create branch
	created, err := svc.Create(ctx, CreateRequest{TenantID: tenantID, Code: "MAIN", Name: "Main Branch"})
	require.NoError(t, err)

	// Deactivate
	updated, err := svc.SetStatus(ctx, tenantID, created.ID, models.StatusInactive, nil)
	require.NoError(t, err)
	assert.Equal(t, models.StatusInactive, updated.Status)
}

func TestService_SetStatus_CannotDeactivatePrimary(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create primary branch
	created, err := svc.Create(ctx, CreateRequest{TenantID: tenantID, Code: "MAIN", Name: "Main Branch", IsPrimary: true})
	require.NoError(t, err)

	// Try to deactivate
	_, err = svc.SetStatus(ctx, tenantID, created.ID, models.StatusInactive, nil)
	assert.ErrorIs(t, err, ErrCannotDeactivatePrimaryBranch)
}

func TestService_Delete(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create branch (non-primary)
	created, err := svc.Create(ctx, CreateRequest{TenantID: tenantID, Code: "SOUTH", Name: "South Branch"})
	require.NoError(t, err)

	// Delete
	err = svc.Delete(ctx, tenantID, created.ID)
	require.NoError(t, err)

	// Verify deleted
	_, err = svc.GetByID(ctx, tenantID, created.ID)
	assert.ErrorIs(t, err, ErrBranchNotFound)
}

func TestService_Delete_CannotDeletePrimary(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create primary branch
	created, err := svc.Create(ctx, CreateRequest{TenantID: tenantID, Code: "MAIN", Name: "Main Branch", IsPrimary: true})
	require.NoError(t, err)

	// Try to delete
	err = svc.Delete(ctx, tenantID, created.ID)
	assert.ErrorIs(t, err, ErrCannotDeletePrimaryBranch)
}

func TestService_GetPrimaryBranch(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create primary branch
	created, err := svc.Create(ctx, CreateRequest{TenantID: tenantID, Code: "MAIN", Name: "Main Branch", IsPrimary: true})
	require.NoError(t, err)

	// Get primary
	primary, err := svc.GetPrimaryBranch(ctx, tenantID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, primary.ID)
}

func TestService_Count(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create branches
	for i := 0; i < 3; i++ {
		_, err := svc.Create(ctx, CreateRequest{
			TenantID: tenantID,
			Code:     string(rune('A' + i)),
			Name:     "Branch " + string(rune('A'+i)),
		})
		require.NoError(t, err)
	}

	count, err := svc.Count(ctx, tenantID)
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)
}
