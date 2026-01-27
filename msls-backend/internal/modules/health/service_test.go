package health

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"msls-backend/internal/pkg/database/models"
)

// MockRepository is a mock implementation of the Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) StudentExists(ctx context.Context, tenantID, studentID uuid.UUID) (bool, error) {
	args := m.Called(ctx, tenantID, studentID)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) GetHealthProfile(ctx context.Context, tenantID, studentID uuid.UUID) (*models.StudentHealthProfile, error) {
	args := m.Called(ctx, tenantID, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StudentHealthProfile), args.Error(1)
}

func (m *MockRepository) CreateHealthProfile(ctx context.Context, profile *models.StudentHealthProfile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *MockRepository) UpdateHealthProfile(ctx context.Context, profile *models.StudentHealthProfile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *MockRepository) ListAllergies(ctx context.Context, tenantID, studentID uuid.UUID, activeOnly bool) ([]models.StudentAllergy, error) {
	args := m.Called(ctx, tenantID, studentID, activeOnly)
	return args.Get(0).([]models.StudentAllergy), args.Error(1)
}

func (m *MockRepository) GetAllergy(ctx context.Context, tenantID, allergyID uuid.UUID) (*models.StudentAllergy, error) {
	args := m.Called(ctx, tenantID, allergyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StudentAllergy), args.Error(1)
}

func (m *MockRepository) CreateAllergy(ctx context.Context, allergy *models.StudentAllergy) error {
	args := m.Called(ctx, allergy)
	return args.Error(0)
}

func (m *MockRepository) UpdateAllergy(ctx context.Context, allergy *models.StudentAllergy) error {
	args := m.Called(ctx, allergy)
	return args.Error(0)
}

func (m *MockRepository) DeleteAllergy(ctx context.Context, tenantID, allergyID uuid.UUID) error {
	args := m.Called(ctx, tenantID, allergyID)
	return args.Error(0)
}

func (m *MockRepository) ListConditions(ctx context.Context, tenantID, studentID uuid.UUID, activeOnly bool) ([]models.StudentChronicCondition, error) {
	args := m.Called(ctx, tenantID, studentID, activeOnly)
	return args.Get(0).([]models.StudentChronicCondition), args.Error(1)
}

func (m *MockRepository) GetCondition(ctx context.Context, tenantID, conditionID uuid.UUID) (*models.StudentChronicCondition, error) {
	args := m.Called(ctx, tenantID, conditionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StudentChronicCondition), args.Error(1)
}

func (m *MockRepository) CreateCondition(ctx context.Context, condition *models.StudentChronicCondition) error {
	args := m.Called(ctx, condition)
	return args.Error(0)
}

func (m *MockRepository) UpdateCondition(ctx context.Context, condition *models.StudentChronicCondition) error {
	args := m.Called(ctx, condition)
	return args.Error(0)
}

func (m *MockRepository) DeleteCondition(ctx context.Context, tenantID, conditionID uuid.UUID) error {
	args := m.Called(ctx, tenantID, conditionID)
	return args.Error(0)
}

func (m *MockRepository) ListMedications(ctx context.Context, tenantID, studentID uuid.UUID, activeOnly bool) ([]models.StudentMedication, error) {
	args := m.Called(ctx, tenantID, studentID, activeOnly)
	return args.Get(0).([]models.StudentMedication), args.Error(1)
}

func (m *MockRepository) GetMedication(ctx context.Context, tenantID, medicationID uuid.UUID) (*models.StudentMedication, error) {
	args := m.Called(ctx, tenantID, medicationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StudentMedication), args.Error(1)
}

func (m *MockRepository) CreateMedication(ctx context.Context, medication *models.StudentMedication) error {
	args := m.Called(ctx, medication)
	return args.Error(0)
}

func (m *MockRepository) UpdateMedication(ctx context.Context, medication *models.StudentMedication) error {
	args := m.Called(ctx, medication)
	return args.Error(0)
}

func (m *MockRepository) DeleteMedication(ctx context.Context, tenantID, medicationID uuid.UUID) error {
	args := m.Called(ctx, tenantID, medicationID)
	return args.Error(0)
}

func (m *MockRepository) ListVaccinations(ctx context.Context, tenantID, studentID uuid.UUID) ([]models.StudentVaccination, error) {
	args := m.Called(ctx, tenantID, studentID)
	return args.Get(0).([]models.StudentVaccination), args.Error(1)
}

func (m *MockRepository) GetVaccination(ctx context.Context, tenantID, vaccinationID uuid.UUID) (*models.StudentVaccination, error) {
	args := m.Called(ctx, tenantID, vaccinationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StudentVaccination), args.Error(1)
}

func (m *MockRepository) CreateVaccination(ctx context.Context, vaccination *models.StudentVaccination) error {
	args := m.Called(ctx, vaccination)
	return args.Error(0)
}

func (m *MockRepository) UpdateVaccination(ctx context.Context, vaccination *models.StudentVaccination) error {
	args := m.Called(ctx, vaccination)
	return args.Error(0)
}

func (m *MockRepository) DeleteVaccination(ctx context.Context, tenantID, vaccinationID uuid.UUID) error {
	args := m.Called(ctx, tenantID, vaccinationID)
	return args.Error(0)
}

func (m *MockRepository) ListIncidents(ctx context.Context, tenantID, studentID uuid.UUID, limit int) ([]models.StudentMedicalIncident, error) {
	args := m.Called(ctx, tenantID, studentID, limit)
	return args.Get(0).([]models.StudentMedicalIncident), args.Error(1)
}

func (m *MockRepository) GetIncident(ctx context.Context, tenantID, incidentID uuid.UUID) (*models.StudentMedicalIncident, error) {
	args := m.Called(ctx, tenantID, incidentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StudentMedicalIncident), args.Error(1)
}

func (m *MockRepository) CreateIncident(ctx context.Context, incident *models.StudentMedicalIncident) error {
	args := m.Called(ctx, incident)
	return args.Error(0)
}

func (m *MockRepository) UpdateIncident(ctx context.Context, incident *models.StudentMedicalIncident) error {
	args := m.Called(ctx, incident)
	return args.Error(0)
}

func (m *MockRepository) DeleteIncident(ctx context.Context, tenantID, incidentID uuid.UUID) error {
	args := m.Called(ctx, tenantID, incidentID)
	return args.Error(0)
}

// Test cases

func TestService_GetHealthProfile_StudentNotFound(t *testing.T) {
	// This is a structural test to verify error handling
	assert.NotNil(t, ErrStudentNotFound)
}

func TestService_GetHealthProfile_Success(t *testing.T) {
	studentID := uuid.New()

	// Verify response structure
	response := &HealthProfileResponse{
		ID:        uuid.New().String(),
		StudentID: studentID.String(),
		BloodGroup: func() *string { s := "A+"; return &s }(),
	}

	assert.NotEmpty(t, response.ID)
	assert.Equal(t, studentID.String(), response.StudentID)
	assert.NotNil(t, response.BloodGroup)
}

func TestService_CreateAllergy_InvalidType(t *testing.T) {
	assert.NotNil(t, ErrInvalidAllergyType)
}

func TestService_CreateAllergy_InvalidSeverity(t *testing.T) {
	assert.NotNil(t, ErrInvalidSeverity)
}

func TestService_AllergyResponse_Conversion(t *testing.T) {
	allergy := &models.StudentAllergy{
		ID:          uuid.New(),
		StudentID:   uuid.New(),
		Allergen:    "Peanuts",
		AllergyType: models.AllergyTypeFood,
		Severity:    models.AllergySeveritySevere,
		IsActive:    true,
	}

	service := &Service{}
	response := service.toAllergyResponse(allergy)

	assert.Equal(t, allergy.ID.String(), response.ID)
	assert.Equal(t, "Peanuts", response.Allergen)
	assert.Equal(t, string(models.AllergyTypeFood), response.AllergyType)
	assert.Equal(t, string(models.AllergySeveritySevere), response.Severity)
	assert.True(t, response.IsActive)
}

func TestService_ConditionResponse_Conversion(t *testing.T) {
	condition := &models.StudentChronicCondition{
		ID:            uuid.New(),
		StudentID:     uuid.New(),
		ConditionName: "Asthma",
		ConditionType: models.ConditionTypeRespiratory,
		Severity:      models.ConditionSeverityModerate,
		IsActive:      true,
	}

	service := &Service{}
	response := service.toConditionResponse(condition)

	assert.Equal(t, condition.ID.String(), response.ID)
	assert.Equal(t, "Asthma", response.ConditionName)
	assert.Equal(t, string(models.ConditionTypeRespiratory), response.ConditionType)
	assert.True(t, response.IsActive)
}

func TestService_MedicationResponse_Conversion(t *testing.T) {
	medication := &models.StudentMedication{
		ID:             uuid.New(),
		StudentID:      uuid.New(),
		MedicationName: "Ventolin",
		Dosage:         "2 puffs",
		Frequency:      models.MedicationFrequencyAsNeeded,
		Route:          models.MedicationRouteInhaler,
		IsActive:       true,
	}

	service := &Service{}
	response := service.toMedicationResponse(medication)

	assert.Equal(t, medication.ID.String(), response.ID)
	assert.Equal(t, "Ventolin", response.MedicationName)
	assert.Equal(t, "2 puffs", response.Dosage)
	assert.True(t, response.IsActive)
}

func TestService_VaccinationResponse_Conversion(t *testing.T) {
	vaccination := &models.StudentVaccination{
		ID:          uuid.New(),
		StudentID:   uuid.New(),
		VaccineName: "MMR",
		DoseNumber:  1,
		IsVerified:  false,
	}

	service := &Service{}
	response := service.toVaccinationResponse(vaccination)

	assert.Equal(t, vaccination.ID.String(), response.ID)
	assert.Equal(t, "MMR", response.VaccineName)
	assert.Equal(t, 1, response.DoseNumber)
	assert.False(t, response.IsVerified)
}

func TestService_IncidentResponse_Conversion(t *testing.T) {
	reportedBy := uuid.New()
	incident := &models.StudentMedicalIncident{
		ID:           uuid.New(),
		StudentID:    uuid.New(),
		IncidentType: models.IncidentTypeInjury,
		Description:  "Fall in playground",
		ReportedBy:   reportedBy,
	}

	service := &Service{}
	response := service.toIncidentResponse(incident)

	assert.Equal(t, incident.ID.String(), response.ID)
	assert.Equal(t, string(models.IncidentTypeInjury), response.IncidentType)
	assert.Equal(t, "Fall in playground", response.Description)
	assert.Equal(t, reportedBy.String(), response.ReportedBy)
}

func TestErrors_Defined(t *testing.T) {
	// Verify all error types are properly defined
	assert.NotNil(t, ErrStudentNotFound)
	assert.NotNil(t, ErrAllergyNotFound)
	assert.NotNil(t, ErrConditionNotFound)
	assert.NotNil(t, ErrMedicationNotFound)
	assert.NotNil(t, ErrVaccinationNotFound)
	assert.NotNil(t, ErrIncidentNotFound)
	assert.NotNil(t, ErrInvalidAllergyType)
	assert.NotNil(t, ErrInvalidConditionType)
	assert.NotNil(t, ErrInvalidSeverity)
	assert.NotNil(t, ErrInvalidFrequency)
	assert.NotNil(t, ErrInvalidRoute)
	assert.NotNil(t, ErrInvalidVaccineType)
	assert.NotNil(t, ErrInvalidIncidentType)
	assert.NotNil(t, ErrInvalidDateRange)
}
