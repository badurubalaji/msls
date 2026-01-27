// Package staff provides staff management functionality.
package staff

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

// MockBranchService is a mock for branch.Service
type MockBranchService struct {
	mock.Mock
}

func (m *MockBranchService) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Branch, error) {
	args := m.Called(ctx, tenantID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Branch), args.Error(1)
}

func TestCreateStaffDTO_Validation(t *testing.T) {
	tests := []struct {
		name    string
		dto     CreateStaffDTO
		wantErr error
	}{
		{
			name: "valid DTO",
			dto: CreateStaffDTO{
				TenantID:    uuid.New(),
				BranchID:    uuid.New(),
				FirstName:   "John",
				LastName:    "Doe",
				DateOfBirth: time.Now().AddDate(-30, 0, 0),
				Gender:      models.GenderMale,
				WorkEmail:   "john.doe@company.com",
				WorkPhone:   "9876543210",
				StaffType:   models.StaffTypeTeaching,
				JoinDate:    time.Now(),
			},
			wantErr: nil,
		},
		{
			name: "missing tenant ID",
			dto: CreateStaffDTO{
				BranchID:    uuid.New(),
				FirstName:   "John",
				LastName:    "Doe",
				DateOfBirth: time.Now().AddDate(-30, 0, 0),
				Gender:      models.GenderMale,
				WorkEmail:   "john.doe@company.com",
				WorkPhone:   "9876543210",
				StaffType:   models.StaffTypeTeaching,
				JoinDate:    time.Now(),
			},
			wantErr: ErrTenantIDRequired,
		},
		{
			name: "missing branch ID",
			dto: CreateStaffDTO{
				TenantID:    uuid.New(),
				FirstName:   "John",
				LastName:    "Doe",
				DateOfBirth: time.Now().AddDate(-30, 0, 0),
				Gender:      models.GenderMale,
				WorkEmail:   "john.doe@company.com",
				WorkPhone:   "9876543210",
				StaffType:   models.StaffTypeTeaching,
				JoinDate:    time.Now(),
			},
			wantErr: ErrBranchIDRequired,
		},
		{
			name: "missing first name",
			dto: CreateStaffDTO{
				TenantID:    uuid.New(),
				BranchID:    uuid.New(),
				LastName:    "Doe",
				DateOfBirth: time.Now().AddDate(-30, 0, 0),
				Gender:      models.GenderMale,
				WorkEmail:   "john.doe@company.com",
				WorkPhone:   "9876543210",
				StaffType:   models.StaffTypeTeaching,
				JoinDate:    time.Now(),
			},
			wantErr: ErrFirstNameRequired,
		},
		{
			name: "missing last name",
			dto: CreateStaffDTO{
				TenantID:    uuid.New(),
				BranchID:    uuid.New(),
				FirstName:   "John",
				DateOfBirth: time.Now().AddDate(-30, 0, 0),
				Gender:      models.GenderMale,
				WorkEmail:   "john.doe@company.com",
				WorkPhone:   "9876543210",
				StaffType:   models.StaffTypeTeaching,
				JoinDate:    time.Now(),
			},
			wantErr: ErrLastNameRequired,
		},
		{
			name: "missing date of birth",
			dto: CreateStaffDTO{
				TenantID:  uuid.New(),
				BranchID:  uuid.New(),
				FirstName: "John",
				LastName:  "Doe",
				Gender:    models.GenderMale,
				WorkEmail: "john.doe@company.com",
				WorkPhone: "9876543210",
				StaffType: models.StaffTypeTeaching,
				JoinDate:  time.Now(),
			},
			wantErr: ErrDateOfBirthRequired,
		},
		{
			name: "future date of birth",
			dto: CreateStaffDTO{
				TenantID:    uuid.New(),
				BranchID:    uuid.New(),
				FirstName:   "John",
				LastName:    "Doe",
				DateOfBirth: time.Now().AddDate(1, 0, 0),
				Gender:      models.GenderMale,
				WorkEmail:   "john.doe@company.com",
				WorkPhone:   "9876543210",
				StaffType:   models.StaffTypeTeaching,
				JoinDate:    time.Now(),
			},
			wantErr: ErrInvalidDateOfBirth,
		},
		{
			name: "missing gender",
			dto: CreateStaffDTO{
				TenantID:    uuid.New(),
				BranchID:    uuid.New(),
				FirstName:   "John",
				LastName:    "Doe",
				DateOfBirth: time.Now().AddDate(-30, 0, 0),
				WorkEmail:   "john.doe@company.com",
				WorkPhone:   "9876543210",
				StaffType:   models.StaffTypeTeaching,
				JoinDate:    time.Now(),
			},
			wantErr: ErrGenderRequired,
		},
		{
			name: "invalid gender",
			dto: CreateStaffDTO{
				TenantID:    uuid.New(),
				BranchID:    uuid.New(),
				FirstName:   "John",
				LastName:    "Doe",
				DateOfBirth: time.Now().AddDate(-30, 0, 0),
				Gender:      "invalid",
				WorkEmail:   "john.doe@company.com",
				WorkPhone:   "9876543210",
				StaffType:   models.StaffTypeTeaching,
				JoinDate:    time.Now(),
			},
			wantErr: ErrInvalidGender,
		},
		{
			name: "missing work email",
			dto: CreateStaffDTO{
				TenantID:    uuid.New(),
				BranchID:    uuid.New(),
				FirstName:   "John",
				LastName:    "Doe",
				DateOfBirth: time.Now().AddDate(-30, 0, 0),
				Gender:      models.GenderMale,
				WorkPhone:   "9876543210",
				StaffType:   models.StaffTypeTeaching,
				JoinDate:    time.Now(),
			},
			wantErr: ErrWorkEmailRequired,
		},
		{
			name: "missing work phone",
			dto: CreateStaffDTO{
				TenantID:    uuid.New(),
				BranchID:    uuid.New(),
				FirstName:   "John",
				LastName:    "Doe",
				DateOfBirth: time.Now().AddDate(-30, 0, 0),
				Gender:      models.GenderMale,
				WorkEmail:   "john.doe@company.com",
				StaffType:   models.StaffTypeTeaching,
				JoinDate:    time.Now(),
			},
			wantErr: ErrWorkPhoneRequired,
		},
		{
			name: "missing staff type",
			dto: CreateStaffDTO{
				TenantID:    uuid.New(),
				BranchID:    uuid.New(),
				FirstName:   "John",
				LastName:    "Doe",
				DateOfBirth: time.Now().AddDate(-30, 0, 0),
				Gender:      models.GenderMale,
				WorkEmail:   "john.doe@company.com",
				WorkPhone:   "9876543210",
				JoinDate:    time.Now(),
			},
			wantErr: ErrInvalidStaffType,
		},
		{
			name: "missing join date",
			dto: CreateStaffDTO{
				TenantID:    uuid.New(),
				BranchID:    uuid.New(),
				FirstName:   "John",
				LastName:    "Doe",
				DateOfBirth: time.Now().AddDate(-30, 0, 0),
				Gender:      models.GenderMale,
				WorkEmail:   "john.doe@company.com",
				WorkPhone:   "9876543210",
				StaffType:   models.StaffTypeTeaching,
			},
			wantErr: ErrJoinDateRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCreateDTO(tt.dto)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestStaffResponse(t *testing.T) {
	staffID := uuid.New()
	tenantID := uuid.New()
	branchID := uuid.New()
	dob := time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC)
	joinDate := time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC)
	now := time.Now()

	staff := &models.Staff{
		ID:               staffID,
		TenantID:         tenantID,
		BranchID:         branchID,
		EmployeeID:       "EMP00001",
		EmployeeIDPrefix: "EMP",
		FirstName:        "John",
		MiddleName:       "William",
		LastName:         "Doe",
		DateOfBirth:      dob,
		Gender:           models.GenderMale,
		BloodGroup:       "O+",
		Nationality:      "Indian",
		WorkEmail:        "john.doe@company.com",
		WorkPhone:        "9876543210",
		StaffType:        models.StaffTypeTeaching,
		JoinDate:         joinDate,
		Status:           models.StaffStatusActive,
		Version:          1,
		CreatedAt:        now,
		UpdatedAt:        now,
		CurrentAddressLine1: "123 Main St",
		CurrentAddressLine2: "Apt 4B",
		CurrentCity:         "Mumbai",
		CurrentState:        "Maharashtra",
		CurrentPincode:      "400001",
		CurrentCountry:      "India",
	}

	resp := ToStaffResponse(staff)

	assert.Equal(t, staffID.String(), resp.ID)
	assert.Equal(t, "EMP00001", resp.EmployeeID)
	assert.Equal(t, "EMP", resp.EmployeeIDPrefix)
	assert.Equal(t, "John", resp.FirstName)
	assert.Equal(t, "William", resp.MiddleName)
	assert.Equal(t, "Doe", resp.LastName)
	assert.Equal(t, "John William Doe", resp.FullName)
	assert.Equal(t, "1990-05-15", resp.DateOfBirth)
	assert.Equal(t, "male", resp.Gender)
	assert.Equal(t, "O+", resp.BloodGroup)
	assert.Equal(t, "Indian", resp.Nationality)
	assert.Equal(t, "john.doe@company.com", resp.WorkEmail)
	assert.Equal(t, "9876543210", resp.WorkPhone)
	assert.Equal(t, "teaching", resp.StaffType)
	assert.Equal(t, "2020-01-15", resp.JoinDate)
	assert.Equal(t, "active", resp.Status)
	assert.Equal(t, "JD", resp.Initials)
	assert.Equal(t, 1, resp.Version)
	require.NotNil(t, resp.CurrentAddress)
	assert.Equal(t, "123 Main St", resp.CurrentAddress.AddressLine1)
	assert.Equal(t, "Mumbai", resp.CurrentAddress.City)
}

func TestStaffStatus_IsValid(t *testing.T) {
	tests := []struct {
		status   models.StaffStatus
		expected bool
	}{
		{models.StaffStatusActive, true},
		{models.StaffStatusInactive, true},
		{models.StaffStatusTerminated, true},
		{models.StaffStatusOnLeave, true},
		{models.StaffStatus("invalid"), false},
		{models.StaffStatus(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.status.IsValid())
		})
	}
}

func TestStaffType_IsValid(t *testing.T) {
	tests := []struct {
		staffType models.StaffType
		expected  bool
	}{
		{models.StaffTypeTeaching, true},
		{models.StaffTypeNonTeaching, true},
		{models.StaffType("invalid"), false},
		{models.StaffType(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.staffType), func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.staffType.IsValid())
		})
	}
}

func TestStaff_FullName(t *testing.T) {
	tests := []struct {
		name     string
		staff    models.Staff
		expected string
	}{
		{
			name: "with middle name",
			staff: models.Staff{
				FirstName:  "John",
				MiddleName: "William",
				LastName:   "Doe",
			},
			expected: "John William Doe",
		},
		{
			name: "without middle name",
			staff: models.Staff{
				FirstName: "John",
				LastName:  "Doe",
			},
			expected: "John Doe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.staff.FullName())
		})
	}
}

func TestStaff_GetInitials(t *testing.T) {
	tests := []struct {
		name     string
		staff    models.Staff
		expected string
	}{
		{
			name: "normal names",
			staff: models.Staff{
				FirstName: "John",
				LastName:  "Doe",
			},
			expected: "JD",
		},
		{
			name: "empty first name",
			staff: models.Staff{
				FirstName: "",
				LastName:  "Doe",
			},
			expected: "D",
		},
		{
			name: "empty last name",
			staff: models.Staff{
				FirstName: "John",
				LastName:  "",
			},
			expected: "J",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.staff.GetInitials())
		})
	}
}

func TestListFilter_Defaults(t *testing.T) {
	filter := ListFilter{
		TenantID: uuid.New(),
	}

	// Verify defaults
	assert.Nil(t, filter.BranchID)
	assert.Nil(t, filter.Status)
	assert.Nil(t, filter.Gender)
	assert.Nil(t, filter.StaffType)
	assert.Empty(t, filter.Search)
	assert.Empty(t, filter.Cursor)
	assert.Equal(t, 0, filter.Limit) // Will default to 20 in repository
}

func TestStatusUpdateDTO_Validation(t *testing.T) {
	tests := []struct {
		name    string
		dto     StatusUpdateDTO
		wantErr error
	}{
		{
			name: "valid DTO - active to inactive",
			dto: StatusUpdateDTO{
				Status:        models.StaffStatusInactive,
				Reason:        "Employee resigned",
				EffectiveDate: time.Now(),
			},
			wantErr: nil,
		},
		{
			name: "valid DTO - terminate",
			dto: StatusUpdateDTO{
				Status:        models.StaffStatusTerminated,
				Reason:        "Contract ended",
				EffectiveDate: time.Now(),
			},
			wantErr: nil,
		},
		{
			name: "invalid status",
			dto: StatusUpdateDTO{
				Status:        models.StaffStatus("invalid"),
				Reason:        "Some reason",
				EffectiveDate: time.Now(),
			},
			wantErr: ErrInvalidStatus,
		},
		{
			name: "missing reason",
			dto: StatusUpdateDTO{
				Status:        models.StaffStatusInactive,
				EffectiveDate: time.Now(),
			},
			wantErr: ErrStatusReasonRequired,
		},
		{
			name: "missing effective date",
			dto: StatusUpdateDTO{
				Status: models.StaffStatusInactive,
				Reason: "Employee resigned",
			},
			wantErr: ErrEffectiveDateRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateStatusUpdateDTO(tt.dto)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestToStaffResponses(t *testing.T) {
	now := time.Now()
	staffList := []models.Staff{
		{
			ID:          uuid.New(),
			TenantID:    uuid.New(),
			BranchID:    uuid.New(),
			EmployeeID:  "EMP00001",
			FirstName:   "John",
			LastName:    "Doe",
			DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			Gender:      models.GenderMale,
			WorkEmail:   "john@example.com",
			WorkPhone:   "1234567890",
			StaffType:   models.StaffTypeTeaching,
			JoinDate:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			Status:      models.StaffStatusActive,
			Version:     1,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          uuid.New(),
			TenantID:    uuid.New(),
			BranchID:    uuid.New(),
			EmployeeID:  "EMP00002",
			FirstName:   "Jane",
			LastName:    "Smith",
			DateOfBirth: time.Date(1992, 6, 15, 0, 0, 0, 0, time.UTC),
			Gender:      models.GenderFemale,
			WorkEmail:   "jane@example.com",
			WorkPhone:   "0987654321",
			StaffType:   models.StaffTypeNonTeaching,
			JoinDate:    time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC),
			Status:      models.StaffStatusActive,
			Version:     1,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	responses := ToStaffResponses(staffList)

	assert.Len(t, responses, 2)
	assert.Equal(t, "John", responses[0].FirstName)
	assert.Equal(t, "Doe", responses[0].LastName)
	assert.Equal(t, "EMP00001", responses[0].EmployeeID)
	assert.Equal(t, "Jane", responses[1].FirstName)
	assert.Equal(t, "Smith", responses[1].LastName)
	assert.Equal(t, "EMP00002", responses[1].EmployeeID)
}

func TestToStatusHistoryResponse(t *testing.T) {
	historyID := uuid.New()
	staffID := uuid.New()
	changedBy := uuid.New()
	effectiveDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
	changedAt := time.Date(2025, 1, 10, 12, 0, 0, 0, time.UTC)

	history := &models.StaffStatusHistory{
		ID:            historyID,
		StaffID:       staffID,
		OldStatus:     "active",
		NewStatus:     "on_leave",
		Reason:        "Medical leave",
		EffectiveDate: effectiveDate,
		ChangedBy:     &changedBy,
		ChangedAt:     changedAt,
	}

	resp := ToStatusHistoryResponse(history)

	assert.Equal(t, historyID.String(), resp.ID)
	assert.Equal(t, "active", resp.OldStatus)
	assert.Equal(t, "on_leave", resp.NewStatus)
	assert.Equal(t, "Medical leave", resp.Reason)
	assert.Equal(t, "2025-01-15", resp.EffectiveDate)
	assert.Equal(t, changedBy.String(), resp.ChangedBy)
}

func TestToStatusHistoryResponses(t *testing.T) {
	now := time.Now()
	histories := []models.StaffStatusHistory{
		{
			ID:            uuid.New(),
			StaffID:       uuid.New(),
			OldStatus:     "",
			NewStatus:     "active",
			Reason:        "Initial hire",
			EffectiveDate: now.AddDate(0, -6, 0),
			ChangedAt:     now.AddDate(0, -6, 0),
		},
		{
			ID:            uuid.New(),
			StaffID:       uuid.New(),
			OldStatus:     "active",
			NewStatus:     "on_leave",
			Reason:        "Vacation",
			EffectiveDate: now.AddDate(0, -1, 0),
			ChangedAt:     now.AddDate(0, -1, 0),
		},
	}

	responses := ToStatusHistoryResponses(histories)

	assert.Len(t, responses, 2)
	assert.Equal(t, "", responses[0].OldStatus)
	assert.Equal(t, "active", responses[0].NewStatus)
	assert.Equal(t, "active", responses[1].OldStatus)
	assert.Equal(t, "on_leave", responses[1].NewStatus)
}

// Helper function to validate CreateStaffDTO (matches service validation)
func validateCreateDTO(dto CreateStaffDTO) error {
	if dto.TenantID == uuid.Nil {
		return ErrTenantIDRequired
	}
	if dto.BranchID == uuid.Nil {
		return ErrBranchIDRequired
	}
	if dto.FirstName == "" {
		return ErrFirstNameRequired
	}
	if dto.LastName == "" {
		return ErrLastNameRequired
	}
	if dto.DateOfBirth.IsZero() {
		return ErrDateOfBirthRequired
	}
	if dto.DateOfBirth.After(time.Now()) {
		return ErrInvalidDateOfBirth
	}
	if dto.Gender == "" {
		return ErrGenderRequired
	}
	if !dto.Gender.IsValid() {
		return ErrInvalidGender
	}
	if dto.WorkEmail == "" {
		return ErrWorkEmailRequired
	}
	if dto.WorkPhone == "" {
		return ErrWorkPhoneRequired
	}
	if dto.StaffType == "" {
		return ErrInvalidStaffType
	}
	if !dto.StaffType.IsValid() {
		return ErrInvalidStaffType
	}
	if dto.JoinDate.IsZero() {
		return ErrJoinDateRequired
	}
	return nil
}

// Helper function to validate StatusUpdateDTO (matches service validation)
func validateStatusUpdateDTO(dto StatusUpdateDTO) error {
	if dto.Status == "" || !dto.Status.IsValid() {
		return ErrInvalidStatus
	}
	if dto.Reason == "" {
		return ErrStatusReasonRequired
	}
	if dto.EffectiveDate.IsZero() {
		return ErrEffectiveDateRequired
	}
	return nil
}
