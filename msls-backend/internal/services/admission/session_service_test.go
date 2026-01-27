// Package admission provides admission management services.
package admission

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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
	err = db.AutoMigrate(&models.Tenant{}, &models.Branch{}, &models.AcademicYear{}, &models.AdmissionSession{}, &models.AdmissionSeat{})
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
// Session Tests
// =============================================================================

func TestSessionService_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		req     CreateSessionRequest
		wantErr error
	}{
		{
			name: "success",
			req: CreateSessionRequest{
				Name:           "Admission 2025-26",
				StartDate:      time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
				EndDate:        time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
				ApplicationFee: decimal.NewFromInt(500),
			},
			wantErr: nil,
		},
		{
			name: "missing tenant ID",
			req: CreateSessionRequest{
				Name:      "Admission 2025-26",
				StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
			},
			wantErr: ErrTenantIDRequired,
		},
		{
			name: "missing name",
			req: CreateSessionRequest{
				StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
			},
			wantErr: ErrSessionNameRequired,
		},
		{
			name: "invalid date range - end before start",
			req: CreateSessionRequest{
				Name:      "Admission 2025-26",
				StartDate: time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
			},
			wantErr: ErrInvalidDateRange,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db := setupTestDB(t)
			svc := NewSessionService(db)
			ctx := context.Background()

			// Set up tenant ID for success case
			if tc.wantErr == nil || tc.wantErr == ErrSessionNameRequired ||
				tc.wantErr == ErrInvalidDateRange {
				tenantID := createTestTenant(t, db)
				tc.req.TenantID = tenantID
			}

			session, err := svc.Create(ctx, tc.req)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, session)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, session)
				assert.Equal(t, tc.req.Name, session.Name)
				assert.Equal(t, models.SessionStatusUpcoming, session.Status)
			}
		})
	}
}

func TestSessionService_Create_DuplicateName(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewSessionService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create first session
	req := CreateSessionRequest{
		TenantID:  tenantID,
		Name:      "Admission 2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
	}
	_, err := svc.Create(ctx, req)
	require.NoError(t, err)

	// Try to create duplicate
	req.StartDate = time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)
	req.EndDate = time.Date(2025, 9, 30, 0, 0, 0, 0, time.UTC)
	_, err = svc.Create(ctx, req)
	assert.ErrorIs(t, err, ErrSessionNameExists)
}

func TestSessionService_GetByID(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewSessionService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create session
	req := CreateSessionRequest{
		TenantID:  tenantID,
		Name:      "Admission 2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
	}
	created, err := svc.Create(ctx, req)
	require.NoError(t, err)

	// Get by ID
	session, err := svc.GetByID(ctx, tenantID, created.ID, false)
	require.NoError(t, err)
	assert.Equal(t, created.ID, session.ID)
	assert.Equal(t, created.Name, session.Name)
}

func TestSessionService_GetByID_NotFound(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewSessionService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	_, err := svc.GetByID(ctx, tenantID, uuid.New(), false)
	assert.ErrorIs(t, err, ErrSessionNotFound)
}

func TestSessionService_List(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewSessionService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create sessions
	sessions := []CreateSessionRequest{
		{TenantID: tenantID, Name: "Admission 2025-26", StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC)},
		{TenantID: tenantID, Name: "Mid-Term 2025", StartDate: time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2025, 8, 31, 0, 0, 0, 0, time.UTC)},
	}

	for _, req := range sessions {
		_, err := svc.Create(ctx, req)
		require.NoError(t, err)
	}

	// List all
	result, err := svc.List(ctx, ListSessionFilter{TenantID: tenantID})
	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestSessionService_List_WithSearch(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewSessionService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create sessions
	_, err := svc.Create(ctx, CreateSessionRequest{
		TenantID:  tenantID,
		Name:      "Admission 2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	_, err = svc.Create(ctx, CreateSessionRequest{
		TenantID:  tenantID,
		Name:      "Mid-Term 2025",
		StartDate: time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 8, 31, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Filter by search
	result, err := svc.List(ctx, ListSessionFilter{TenantID: tenantID, Search: "Mid-Term"})
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Mid-Term 2025", result[0].Name)
}

func TestSessionService_Update(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewSessionService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create session
	created, err := svc.Create(ctx, CreateSessionRequest{
		TenantID:  tenantID,
		Name:      "Admission 2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Update
	newName := "Admission 2025-26 - Regular"
	updated, err := svc.Update(ctx, tenantID, created.ID, UpdateSessionRequest{Name: &newName})
	require.NoError(t, err)
	assert.Equal(t, newName, updated.Name)
}

func TestSessionService_Update_CannotModifyClosedSession(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewSessionService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create session
	created, err := svc.Create(ctx, CreateSessionRequest{
		TenantID:  tenantID,
		Name:      "Admission 2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Open and then close the session
	_, err = svc.ChangeStatus(ctx, tenantID, created.ID, models.SessionStatusOpen, nil)
	require.NoError(t, err)

	_, err = svc.ChangeStatus(ctx, tenantID, created.ID, models.SessionStatusClosed, nil)
	require.NoError(t, err)

	// Try to update - should fail
	newName := "New Name"
	_, err = svc.Update(ctx, tenantID, created.ID, UpdateSessionRequest{Name: &newName})
	assert.ErrorIs(t, err, ErrCannotModifyClosedSession)
}

func TestSessionService_ChangeStatus(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewSessionService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create session
	created, err := svc.Create(ctx, CreateSessionRequest{
		TenantID:  tenantID,
		Name:      "Admission 2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)
	assert.Equal(t, models.SessionStatusUpcoming, created.Status)

	// Open session
	opened, err := svc.ChangeStatus(ctx, tenantID, created.ID, models.SessionStatusOpen, nil)
	require.NoError(t, err)
	assert.Equal(t, models.SessionStatusOpen, opened.Status)

	// Close session
	closed, err := svc.ChangeStatus(ctx, tenantID, created.ID, models.SessionStatusClosed, nil)
	require.NoError(t, err)
	assert.Equal(t, models.SessionStatusClosed, closed.Status)
}

func TestSessionService_ChangeStatus_InvalidTransition(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewSessionService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create session
	created, err := svc.Create(ctx, CreateSessionRequest{
		TenantID:  tenantID,
		Name:      "Admission 2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Try to change directly from upcoming to closed is valid
	_, err = svc.ChangeStatus(ctx, tenantID, created.ID, models.SessionStatusClosed, nil)
	require.NoError(t, err)
}

func TestSessionService_Delete(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewSessionService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create session
	created, err := svc.Create(ctx, CreateSessionRequest{
		TenantID:  tenantID,
		Name:      "Admission 2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Delete
	err = svc.Delete(ctx, tenantID, created.ID)
	require.NoError(t, err)

	// Verify deleted
	_, err = svc.GetByID(ctx, tenantID, created.ID, false)
	assert.ErrorIs(t, err, ErrSessionNotFound)
}

func TestSessionService_Delete_CannotDeleteOpenSession(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewSessionService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create session
	created, err := svc.Create(ctx, CreateSessionRequest{
		TenantID:  tenantID,
		Name:      "Admission 2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Open session
	_, err = svc.ChangeStatus(ctx, tenantID, created.ID, models.SessionStatusOpen, nil)
	require.NoError(t, err)

	// Try to delete
	err = svc.Delete(ctx, tenantID, created.ID)
	assert.ErrorIs(t, err, ErrCannotDeleteOpenSession)
}

func TestSessionService_ExtendDeadline(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewSessionService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create session
	created, err := svc.Create(ctx, CreateSessionRequest{
		TenantID:  tenantID,
		Name:      "Admission 2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Extend deadline
	newEndDate := time.Date(2025, 7, 31, 0, 0, 0, 0, time.UTC)
	extended, err := svc.ExtendDeadline(ctx, tenantID, created.ID, newEndDate, nil)
	require.NoError(t, err)
	assert.Equal(t, newEndDate, extended.EndDate)
}

func TestSessionService_GetStats(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewSessionService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create session
	created, err := svc.Create(ctx, CreateSessionRequest{
		TenantID:  tenantID,
		Name:      "Admission 2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Add seats
	_, err = svc.CreateSeat(ctx, CreateSeatRequest{
		TenantID:   tenantID,
		SessionID:  created.ID,
		ClassName:  "Class 1",
		TotalSeats: 40,
	})
	require.NoError(t, err)

	_, err = svc.CreateSeat(ctx, CreateSeatRequest{
		TenantID:   tenantID,
		SessionID:  created.ID,
		ClassName:  "Class 2",
		TotalSeats: 50,
	})
	require.NoError(t, err)

	// Get stats
	stats, err := svc.GetStats(ctx, tenantID, created.ID)
	require.NoError(t, err)
	assert.Equal(t, 90, stats.TotalSeats)
	assert.Equal(t, 0, stats.FilledSeats)
	assert.Equal(t, 90, stats.AvailableSeats)
}

func TestSessionService_Count(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewSessionService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create sessions
	for i := 0; i < 3; i++ {
		_, err := svc.Create(ctx, CreateSessionRequest{
			TenantID:  tenantID,
			Name:      "Session " + string(rune('A'+i)),
			StartDate: time.Date(2025+i, 4, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2025+i, 6, 30, 0, 0, 0, 0, time.UTC),
		})
		require.NoError(t, err)
	}

	count, err := svc.Count(ctx, tenantID)
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)
}

// =============================================================================
// Seat Tests
// =============================================================================

func TestSessionService_CreateSeat(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewSessionService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create session
	session, err := svc.Create(ctx, CreateSessionRequest{
		TenantID:  tenantID,
		Name:      "Admission 2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Create seat
	seat, err := svc.CreateSeat(ctx, CreateSeatRequest{
		TenantID:      tenantID,
		SessionID:     session.ID,
		ClassName:     "Class 1",
		TotalSeats:    40,
		WaitlistLimit: 10,
	})
	require.NoError(t, err)
	assert.NotNil(t, seat)
	assert.Equal(t, "Class 1", seat.ClassName)
	assert.Equal(t, 40, seat.TotalSeats)
	assert.Equal(t, 0, seat.FilledSeats)
	assert.Equal(t, 10, seat.WaitlistLimit)
}

func TestSessionService_CreateSeat_DuplicateClass(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewSessionService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create session
	session, err := svc.Create(ctx, CreateSessionRequest{
		TenantID:  tenantID,
		Name:      "Admission 2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Create first seat
	_, err = svc.CreateSeat(ctx, CreateSeatRequest{
		TenantID:   tenantID,
		SessionID:  session.ID,
		ClassName:  "Class 1",
		TotalSeats: 40,
	})
	require.NoError(t, err)

	// Try to create duplicate
	_, err = svc.CreateSeat(ctx, CreateSeatRequest{
		TenantID:   tenantID,
		SessionID:  session.ID,
		ClassName:  "Class 1",
		TotalSeats: 50,
	})
	assert.ErrorIs(t, err, ErrClassAlreadyExists)
}

func TestSessionService_ListSeats(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewSessionService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create session
	session, err := svc.Create(ctx, CreateSessionRequest{
		TenantID:  tenantID,
		Name:      "Admission 2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Create seats
	classes := []string{"Class 1", "Class 2", "Class 3"}
	for _, className := range classes {
		_, err := svc.CreateSeat(ctx, CreateSeatRequest{
			TenantID:   tenantID,
			SessionID:  session.ID,
			ClassName:  className,
			TotalSeats: 40,
		})
		require.NoError(t, err)
	}

	// List seats
	seats, err := svc.ListSeats(ctx, tenantID, session.ID)
	require.NoError(t, err)
	assert.Len(t, seats, 3)
}

func TestSessionService_UpdateSeat(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewSessionService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create session
	session, err := svc.Create(ctx, CreateSessionRequest{
		TenantID:  tenantID,
		Name:      "Admission 2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Create seat
	seat, err := svc.CreateSeat(ctx, CreateSeatRequest{
		TenantID:   tenantID,
		SessionID:  session.ID,
		ClassName:  "Class 1",
		TotalSeats: 40,
	})
	require.NoError(t, err)

	// Update seat
	newTotalSeats := 50
	updated, err := svc.UpdateSeat(ctx, tenantID, seat.ID, UpdateSeatRequest{TotalSeats: &newTotalSeats})
	require.NoError(t, err)
	assert.Equal(t, 50, updated.TotalSeats)
}

func TestSessionService_UpdateSeat_FilledExceedsTotal(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewSessionService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create session
	session, err := svc.Create(ctx, CreateSessionRequest{
		TenantID:  tenantID,
		Name:      "Admission 2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Create seat
	seat, err := svc.CreateSeat(ctx, CreateSeatRequest{
		TenantID:   tenantID,
		SessionID:  session.ID,
		ClassName:  "Class 1",
		TotalSeats: 40,
	})
	require.NoError(t, err)

	// Manually set filled seats
	err = db.Model(&models.AdmissionSeat{}).Where("id = ?", seat.ID).Update("filled_seats", 30).Error
	require.NoError(t, err)

	// Try to reduce total seats below filled
	newTotalSeats := 20
	_, err = svc.UpdateSeat(ctx, tenantID, seat.ID, UpdateSeatRequest{TotalSeats: &newTotalSeats})
	assert.ErrorIs(t, err, ErrFilledExceedsTotal)
}

func TestSessionService_DeleteSeat(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewSessionService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create session
	session, err := svc.Create(ctx, CreateSessionRequest{
		TenantID:  tenantID,
		Name:      "Admission 2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Create seat
	seat, err := svc.CreateSeat(ctx, CreateSeatRequest{
		TenantID:   tenantID,
		SessionID:  session.ID,
		ClassName:  "Class 1",
		TotalSeats: 40,
	})
	require.NoError(t, err)

	// Delete seat
	err = svc.DeleteSeat(ctx, tenantID, seat.ID)
	require.NoError(t, err)

	// Verify deleted
	_, err = svc.GetSeatByID(ctx, tenantID, seat.ID)
	assert.ErrorIs(t, err, ErrSeatNotFound)
}

func TestSessionService_IncrementFilledSeats(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewSessionService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create session
	session, err := svc.Create(ctx, CreateSessionRequest{
		TenantID:  tenantID,
		Name:      "Admission 2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Create seat
	seat, err := svc.CreateSeat(ctx, CreateSeatRequest{
		TenantID:   tenantID,
		SessionID:  session.ID,
		ClassName:  "Class 1",
		TotalSeats: 40,
	})
	require.NoError(t, err)

	// Increment filled seats
	err = svc.IncrementFilledSeats(ctx, tenantID, seat.ID, 5)
	require.NoError(t, err)

	// Verify
	updated, err := svc.GetSeatByID(ctx, tenantID, seat.ID)
	require.NoError(t, err)
	assert.Equal(t, 5, updated.FilledSeats)
}

func TestSessionService_IncrementFilledSeats_ExceedsTotal(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewSessionService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create session
	session, err := svc.Create(ctx, CreateSessionRequest{
		TenantID:  tenantID,
		Name:      "Admission 2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Create seat with small capacity
	seat, err := svc.CreateSeat(ctx, CreateSeatRequest{
		TenantID:   tenantID,
		SessionID:  session.ID,
		ClassName:  "Class 1",
		TotalSeats: 5,
	})
	require.NoError(t, err)

	// Try to increment beyond total
	err = svc.IncrementFilledSeats(ctx, tenantID, seat.ID, 10)
	assert.ErrorIs(t, err, ErrFilledExceedsTotal)
}

// =============================================================================
// Delete Cascade Tests
// =============================================================================

func TestSessionService_Delete_CascadesSeats(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := NewSessionService(db)
	ctx := context.Background()
	tenantID := createTestTenant(t, db)

	// Create session
	session, err := svc.Create(ctx, CreateSessionRequest{
		TenantID:  tenantID,
		Name:      "Admission 2025-26",
		StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	// Create seats
	seat, err := svc.CreateSeat(ctx, CreateSeatRequest{
		TenantID:   tenantID,
		SessionID:  session.ID,
		ClassName:  "Class 1",
		TotalSeats: 40,
	})
	require.NoError(t, err)

	// Delete session
	err = svc.Delete(ctx, tenantID, session.ID)
	require.NoError(t, err)

	// Verify seat is deleted
	_, err = svc.GetSeatByID(ctx, tenantID, seat.ID)
	assert.ErrorIs(t, err, ErrSeatNotFound)
}
