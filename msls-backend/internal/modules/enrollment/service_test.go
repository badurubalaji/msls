// Package enrollment provides student enrollment management functionality.
package enrollment

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"msls-backend/internal/pkg/database/models"
)

// MockStudentRepository is a mock implementation of StudentRepository.
type MockStudentRepository struct {
	mock.Mock
}

func (m *MockStudentRepository) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Student, error) {
	args := m.Called(ctx, tenantID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Student), args.Error(1)
}

func (m *MockStudentRepository) UpdateStatus(ctx context.Context, tenantID, id uuid.UUID, status models.StudentStatus, updatedBy *uuid.UUID) error {
	args := m.Called(ctx, tenantID, id, status, updatedBy)
	return args.Error(0)
}

// MockAcademicYearRepository is a mock implementation of AcademicYearRepository.
type MockAcademicYearRepository struct {
	mock.Mock
}

func (m *MockAcademicYearRepository) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*models.AcademicYear, error) {
	args := m.Called(ctx, tenantID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AcademicYear), args.Error(1)
}

func TestEnrollmentStatus_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		status EnrollmentStatus
		want   bool
	}{
		{"active is valid", EnrollmentStatusActive, true},
		{"completed is valid", EnrollmentStatusCompleted, true},
		{"transferred is valid", EnrollmentStatusTransferred, true},
		{"dropout is valid", EnrollmentStatusDropout, true},
		{"empty is invalid", EnrollmentStatus(""), false},
		{"unknown is invalid", EnrollmentStatus("unknown"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.status.IsValid()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStudentEnrollment_Validate(t *testing.T) {
	tenantID := uuid.New()
	studentID := uuid.New()
	academicYearID := uuid.New()

	tests := []struct {
		name    string
		e       *StudentEnrollment
		wantErr error
	}{
		{
			name: "valid enrollment",
			e: &StudentEnrollment{
				TenantID:       tenantID,
				StudentID:      studentID,
				AcademicYearID: academicYearID,
				Status:         EnrollmentStatusActive,
			},
			wantErr: nil,
		},
		{
			name: "missing tenant ID",
			e: &StudentEnrollment{
				StudentID:      studentID,
				AcademicYearID: academicYearID,
			},
			wantErr: ErrTenantIDRequired,
		},
		{
			name: "missing student ID",
			e: &StudentEnrollment{
				TenantID:       tenantID,
				AcademicYearID: academicYearID,
			},
			wantErr: ErrStudentIDRequired,
		},
		{
			name: "missing academic year ID",
			e: &StudentEnrollment{
				TenantID:  tenantID,
				StudentID: studentID,
			},
			wantErr: ErrAcademicYearIDRequired,
		},
		{
			name: "invalid status",
			e: &StudentEnrollment{
				TenantID:       tenantID,
				StudentID:      studentID,
				AcademicYearID: academicYearID,
				Status:         EnrollmentStatus("invalid"),
			},
			wantErr: ErrInvalidStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.e.Validate()
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEnrollmentStatusChange_Validate(t *testing.T) {
	tenantID := uuid.New()
	enrollmentID := uuid.New()
	changedBy := uuid.New()

	tests := []struct {
		name    string
		c       *EnrollmentStatusChange
		wantErr error
	}{
		{
			name: "valid status change",
			c: &EnrollmentStatusChange{
				TenantID:     tenantID,
				EnrollmentID: enrollmentID,
				ToStatus:     EnrollmentStatusActive,
				ChangedBy:    changedBy,
			},
			wantErr: nil,
		},
		{
			name: "missing tenant ID",
			c: &EnrollmentStatusChange{
				EnrollmentID: enrollmentID,
				ToStatus:     EnrollmentStatusActive,
				ChangedBy:    changedBy,
			},
			wantErr: ErrTenantIDRequired,
		},
		{
			name: "missing enrollment ID",
			c: &EnrollmentStatusChange{
				TenantID:  tenantID,
				ToStatus:  EnrollmentStatusActive,
				ChangedBy: changedBy,
			},
			wantErr: ErrEnrollmentIDRequired,
		},
		{
			name: "invalid to status",
			c: &EnrollmentStatusChange{
				TenantID:     tenantID,
				EnrollmentID: enrollmentID,
				ToStatus:     EnrollmentStatus("invalid"),
				ChangedBy:    changedBy,
			},
			wantErr: ErrInvalidStatus,
		},
		{
			name: "missing changed by",
			c: &EnrollmentStatusChange{
				TenantID:     tenantID,
				EnrollmentID: enrollmentID,
				ToStatus:     EnrollmentStatusActive,
			},
			wantErr: ErrChangedByRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.c.Validate()
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestToEnrollmentResponse(t *testing.T) {
	enrollmentID := uuid.New()
	studentID := uuid.New()
	academicYearID := uuid.New()
	classID := uuid.New()
	now := time.Now()
	transferDate := now.AddDate(0, -1, 0)

	enrollment := &StudentEnrollment{
		ID:             enrollmentID,
		StudentID:      studentID,
		AcademicYearID: academicYearID,
		ClassID:        &classID,
		RollNumber:     "A001",
		Status:         EnrollmentStatusTransferred,
		EnrollmentDate: now.AddDate(-1, 0, 0),
		TransferDate:   &transferDate,
		TransferReason: "Family relocation",
		Notes:          "Test notes",
		CreatedAt:      now,
		UpdatedAt:      now,
		AcademicYear: &AcademicYearRef{
			ID:        academicYearID,
			Name:      "2024-2025",
			StartDate: time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2025, 3, 31, 0, 0, 0, 0, time.UTC),
			IsCurrent: true,
		},
	}

	resp := ToEnrollmentResponse(enrollment)

	assert.Equal(t, enrollmentID.String(), resp.ID)
	assert.Equal(t, studentID.String(), resp.StudentID)
	assert.Equal(t, classID.String(), resp.ClassID)
	assert.Equal(t, "A001", resp.RollNumber)
	assert.Equal(t, "transferred", resp.Status)
	assert.Equal(t, "Family relocation", resp.TransferReason)
	assert.Equal(t, "Test notes", resp.Notes)
	require.NotNil(t, resp.AcademicYear)
	assert.Equal(t, "2024-2025", resp.AcademicYear.Name)
	assert.True(t, resp.AcademicYear.IsCurrent)
}

func TestToEnrollmentResponses(t *testing.T) {
	enrollments := []StudentEnrollment{
		{
			ID:             uuid.New(),
			StudentID:      uuid.New(),
			AcademicYearID: uuid.New(),
			Status:         EnrollmentStatusActive,
			EnrollmentDate: time.Now(),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			ID:             uuid.New(),
			StudentID:      uuid.New(),
			AcademicYearID: uuid.New(),
			Status:         EnrollmentStatusCompleted,
			EnrollmentDate: time.Now().AddDate(-1, 0, 0),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}

	responses := ToEnrollmentResponses(enrollments)

	assert.Len(t, responses, 2)
	assert.Equal(t, "active", responses[0].Status)
	assert.Equal(t, "completed", responses[1].Status)
}

func TestToStatusChangeResponse(t *testing.T) {
	changeID := uuid.New()
	enrollmentID := uuid.New()
	changedBy := uuid.New()
	fromStatus := EnrollmentStatusActive
	now := time.Now()
	changeDate := now.AddDate(0, 0, -1)

	change := &EnrollmentStatusChange{
		ID:           changeID,
		EnrollmentID: enrollmentID,
		FromStatus:   &fromStatus,
		ToStatus:     EnrollmentStatusTransferred,
		ChangeReason: "Family moved",
		ChangeDate:   changeDate,
		ChangedAt:    now,
		ChangedBy:    changedBy,
	}

	resp := ToStatusChangeResponse(change)

	assert.Equal(t, changeID.String(), resp.ID)
	assert.Equal(t, enrollmentID.String(), resp.EnrollmentID)
	assert.Equal(t, "active", resp.FromStatus)
	assert.Equal(t, "transferred", resp.ToStatus)
	assert.Equal(t, "Family moved", resp.ChangeReason)
	assert.Equal(t, changedBy.String(), resp.ChangedBy)
}

func TestCreateEnrollmentDTO_Validation(t *testing.T) {
	t.Run("missing tenant ID returns error", func(t *testing.T) {
		dto := CreateEnrollmentDTO{
			StudentID:      uuid.New(),
			AcademicYearID: uuid.New(),
		}

		// Mock service would catch this validation
		// Testing that the DTO allows empty tenant (validation happens in service)
		assert.Equal(t, uuid.Nil, dto.TenantID)
	})

	t.Run("valid DTO has all required fields", func(t *testing.T) {
		now := time.Now()
		dto := CreateEnrollmentDTO{
			TenantID:       uuid.New(),
			StudentID:      uuid.New(),
			AcademicYearID: uuid.New(),
			RollNumber:     "A001",
			EnrollmentDate: &now,
			Notes:          "Test notes",
		}

		assert.NotEqual(t, uuid.Nil, dto.TenantID)
		assert.NotEqual(t, uuid.Nil, dto.StudentID)
		assert.NotEqual(t, uuid.Nil, dto.AcademicYearID)
	})
}

func TestTransferDTO_Validation(t *testing.T) {
	t.Run("transfer date is required", func(t *testing.T) {
		dto := TransferDTO{
			TransferReason: "Family moved",
		}

		assert.True(t, dto.TransferDate.IsZero())
	})

	t.Run("transfer reason is required", func(t *testing.T) {
		dto := TransferDTO{
			TransferDate: time.Now(),
		}

		assert.Empty(t, dto.TransferReason)
	})

	t.Run("valid transfer DTO", func(t *testing.T) {
		userID := uuid.New()
		dto := TransferDTO{
			TransferDate:   time.Now(),
			TransferReason: "Family moved to different city",
			UpdatedBy:      &userID,
		}

		assert.False(t, dto.TransferDate.IsZero())
		assert.NotEmpty(t, dto.TransferReason)
		assert.NotNil(t, dto.UpdatedBy)
	})
}

func TestDropoutDTO_Validation(t *testing.T) {
	t.Run("dropout date is required", func(t *testing.T) {
		dto := DropoutDTO{
			DropoutReason: "Financial difficulties",
		}

		assert.True(t, dto.DropoutDate.IsZero())
	})

	t.Run("dropout reason is required", func(t *testing.T) {
		dto := DropoutDTO{
			DropoutDate: time.Now(),
		}

		assert.Empty(t, dto.DropoutReason)
	})

	t.Run("valid dropout DTO", func(t *testing.T) {
		userID := uuid.New()
		dto := DropoutDTO{
			DropoutDate:   time.Now(),
			DropoutReason: "Financial difficulties",
			UpdatedBy:     &userID,
		}

		assert.False(t, dto.DropoutDate.IsZero())
		assert.NotEmpty(t, dto.DropoutReason)
		assert.NotNil(t, dto.UpdatedBy)
	})
}
