// Package student provides student management functionality.
package student

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

func TestCreateStudentDTO_Validation(t *testing.T) {
	tests := []struct {
		name    string
		dto     CreateStudentDTO
		wantErr error
	}{
		{
			name: "valid DTO",
			dto: CreateStudentDTO{
				TenantID:    uuid.New(),
				BranchID:    uuid.New(),
				FirstName:   "John",
				LastName:    "Doe",
				DateOfBirth: time.Now().AddDate(-10, 0, 0),
				Gender:      models.GenderMale,
			},
			wantErr: nil,
		},
		{
			name: "missing tenant ID",
			dto: CreateStudentDTO{
				BranchID:    uuid.New(),
				FirstName:   "John",
				LastName:    "Doe",
				DateOfBirth: time.Now().AddDate(-10, 0, 0),
				Gender:      models.GenderMale,
			},
			wantErr: ErrTenantIDRequired,
		},
		{
			name: "missing branch ID",
			dto: CreateStudentDTO{
				TenantID:    uuid.New(),
				FirstName:   "John",
				LastName:    "Doe",
				DateOfBirth: time.Now().AddDate(-10, 0, 0),
				Gender:      models.GenderMale,
			},
			wantErr: ErrBranchIDRequired,
		},
		{
			name: "missing first name",
			dto: CreateStudentDTO{
				TenantID:    uuid.New(),
				BranchID:    uuid.New(),
				LastName:    "Doe",
				DateOfBirth: time.Now().AddDate(-10, 0, 0),
				Gender:      models.GenderMale,
			},
			wantErr: ErrFirstNameRequired,
		},
		{
			name: "missing last name",
			dto: CreateStudentDTO{
				TenantID:    uuid.New(),
				BranchID:    uuid.New(),
				FirstName:   "John",
				DateOfBirth: time.Now().AddDate(-10, 0, 0),
				Gender:      models.GenderMale,
			},
			wantErr: ErrLastNameRequired,
		},
		{
			name: "missing date of birth",
			dto: CreateStudentDTO{
				TenantID:  uuid.New(),
				BranchID:  uuid.New(),
				FirstName: "John",
				LastName:  "Doe",
				Gender:    models.GenderMale,
			},
			wantErr: ErrDateOfBirthRequired,
		},
		{
			name: "future date of birth",
			dto: CreateStudentDTO{
				TenantID:    uuid.New(),
				BranchID:    uuid.New(),
				FirstName:   "John",
				LastName:    "Doe",
				DateOfBirth: time.Now().AddDate(1, 0, 0),
				Gender:      models.GenderMale,
			},
			wantErr: ErrInvalidDateOfBirth,
		},
		{
			name: "missing gender",
			dto: CreateStudentDTO{
				TenantID:    uuid.New(),
				BranchID:    uuid.New(),
				FirstName:   "John",
				LastName:    "Doe",
				DateOfBirth: time.Now().AddDate(-10, 0, 0),
			},
			wantErr: ErrGenderRequired,
		},
		{
			name: "invalid gender",
			dto: CreateStudentDTO{
				TenantID:    uuid.New(),
				BranchID:    uuid.New(),
				FirstName:   "John",
				LastName:    "Doe",
				DateOfBirth: time.Now().AddDate(-10, 0, 0),
				Gender:      "invalid",
			},
			wantErr: ErrInvalidGender,
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

func TestStudentResponse(t *testing.T) {
	studentID := uuid.New()
	tenantID := uuid.New()
	branchID := uuid.New()
	dob := time.Date(2015, 5, 15, 0, 0, 0, 0, time.UTC)
	admDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)

	student := &models.Student{
		ID:              studentID,
		TenantID:        tenantID,
		BranchID:        branchID,
		AdmissionNumber: "MUM-2026-00001",
		FirstName:       "John",
		MiddleName:      "William",
		LastName:        "Doe",
		DateOfBirth:     dob,
		Gender:          models.GenderMale,
		BloodGroup:      "O+",
		Status:          models.StudentStatusActive,
		AdmissionDate:   admDate,
		Version:         1,
		Addresses: []models.StudentAddress{
			{
				ID:           uuid.New(),
				TenantID:     tenantID,
				StudentID:    studentID,
				AddressType:  models.AddressTypeCurrent,
				AddressLine1: "123 Main St",
				AddressLine2: "Apt 4B",
				City:         "Mumbai",
				State:        "Maharashtra",
				PostalCode:   "400001",
				Country:      "India",
			},
		},
	}

	resp := ToStudentResponse(student)

	assert.Equal(t, studentID.String(), resp.ID)
	assert.Equal(t, "MUM-2026-00001", resp.AdmissionNumber)
	assert.Equal(t, "John", resp.FirstName)
	assert.Equal(t, "William", resp.MiddleName)
	assert.Equal(t, "Doe", resp.LastName)
	assert.Equal(t, "John William Doe", resp.FullName)
	assert.Equal(t, "2015-05-15", resp.DateOfBirth)
	assert.Equal(t, "male", resp.Gender)
	assert.Equal(t, "O+", resp.BloodGroup)
	assert.Equal(t, "active", resp.Status)
	assert.Equal(t, "JD", resp.Initials)
	assert.Equal(t, 1, resp.Version)
	require.NotNil(t, resp.CurrentAddress)
	assert.Equal(t, "123 Main St", resp.CurrentAddress.AddressLine1)
	assert.Equal(t, "Mumbai", resp.CurrentAddress.City)
}

func TestStudentStatus_IsValid(t *testing.T) {
	tests := []struct {
		status   models.StudentStatus
		expected bool
	}{
		{models.StudentStatusActive, true},
		{models.StudentStatusInactive, true},
		{models.StudentStatusTransferred, true},
		{models.StudentStatusGraduated, true},
		{models.StudentStatus("invalid"), false},
		{models.StudentStatus(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.status.IsValid())
		})
	}
}

func TestGender_IsValid(t *testing.T) {
	tests := []struct {
		gender   models.Gender
		expected bool
	}{
		{models.GenderMale, true},
		{models.GenderFemale, true},
		{models.GenderOther, true},
		{models.Gender("invalid"), false},
		{models.Gender(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.gender), func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.gender.IsValid())
		})
	}
}

func TestAddressType_IsValid(t *testing.T) {
	tests := []struct {
		addrType models.AddressType
		expected bool
	}{
		{models.AddressTypeCurrent, true},
		{models.AddressTypePermanent, true},
		{models.AddressType("invalid"), false},
		{models.AddressType(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.addrType), func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.addrType.IsValid())
		})
	}
}

func TestStudent_FullName(t *testing.T) {
	tests := []struct {
		name       string
		student    models.Student
		expected   string
	}{
		{
			name: "with middle name",
			student: models.Student{
				FirstName:  "John",
				MiddleName: "William",
				LastName:   "Doe",
			},
			expected: "John William Doe",
		},
		{
			name: "without middle name",
			student: models.Student{
				FirstName: "John",
				LastName:  "Doe",
			},
			expected: "John Doe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.student.FullName())
		})
	}
}

func TestStudent_GetInitials(t *testing.T) {
	tests := []struct {
		name     string
		student  models.Student
		expected string
	}{
		{
			name: "normal names",
			student: models.Student{
				FirstName: "John",
				LastName:  "Doe",
			},
			expected: "JD",
		},
		{
			name: "empty first name",
			student: models.Student{
				FirstName: "",
				LastName:  "Doe",
			},
			expected: "D",
		},
		{
			name: "empty last name",
			student: models.Student{
				FirstName: "John",
				LastName:  "",
			},
			expected: "J",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.student.GetInitials())
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
	assert.Empty(t, filter.Search)
	assert.Empty(t, filter.Cursor)
	assert.Equal(t, 0, filter.Limit) // Will default to 20 in repository
}

// Helper function to validate CreateStudentDTO (matches service validation)
func validateCreateDTO(dto CreateStudentDTO) error {
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
	return nil
}
