package health

import (
	"context"
	"time"

	"github.com/google/uuid"

	"msls-backend/internal/pkg/database/models"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// =============================================================================
// Health Profile Service Methods
// =============================================================================

func (s *Service) GetHealthProfile(ctx context.Context, tenantID, studentID uuid.UUID) (*HealthProfileResponse, error) {
	// Verify student exists
	exists, err := s.repo.StudentExists(ctx, tenantID, studentID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrStudentNotFound
	}

	profile, err := s.repo.GetHealthProfile(ctx, tenantID, studentID)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, nil // No profile yet, not an error
	}

	return s.toHealthProfileResponse(profile), nil
}

func (s *Service) CreateOrUpdateHealthProfile(ctx context.Context, tenantID, studentID, userID uuid.UUID, req *CreateHealthProfileRequest) (*HealthProfileResponse, error) {
	// Verify student exists
	exists, err := s.repo.StudentExists(ctx, tenantID, studentID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrStudentNotFound
	}

	// Check if profile exists
	existing, err := s.repo.GetHealthProfile(ctx, tenantID, studentID)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		// Update existing
		s.updateProfileFromRequest(existing, req)
		existing.UpdatedBy = &userID
		existing.UpdatedAt = time.Now()
		if err := s.repo.UpdateHealthProfile(ctx, existing); err != nil {
			return nil, err
		}
		return s.toHealthProfileResponse(existing), nil
	}

	// Create new
	profile := &models.StudentHealthProfile{
		ID:        uuid.New(),
		TenantID:  tenantID,
		StudentID: studentID,
		CreatedBy: &userID,
		UpdatedBy: &userID,
	}
	s.updateProfileFromRequest(profile, req)

	if err := s.repo.CreateHealthProfile(ctx, profile); err != nil {
		return nil, err
	}

	return s.toHealthProfileResponse(profile), nil
}

func (s *Service) updateProfileFromRequest(profile *models.StudentHealthProfile, req *CreateHealthProfileRequest) {
	profile.BloodGroup = req.BloodGroup
	profile.HeightCm = req.HeightCm
	profile.WeightKg = req.WeightKg
	profile.VisionLeft = req.VisionLeft
	profile.VisionRight = req.VisionRight
	profile.HearingStatus = req.HearingStatus
	profile.MedicalNotes = req.MedicalNotes
	profile.InsuranceProvider = req.InsuranceProvider
	profile.InsurancePolicyNumber = req.InsurancePolicyNumber
	profile.PreferredHospital = req.PreferredHospital
	profile.FamilyDoctorName = req.FamilyDoctorName
	profile.FamilyDoctorPhone = req.FamilyDoctorPhone

	if req.InsuranceExpiry != nil {
		if t, err := time.Parse("2006-01-02", *req.InsuranceExpiry); err == nil {
			profile.InsuranceExpiry = &t
		}
	}
	if req.LastCheckupDate != nil {
		if t, err := time.Parse("2006-01-02", *req.LastCheckupDate); err == nil {
			profile.LastCheckupDate = &t
		}
	}
}

func (s *Service) toHealthProfileResponse(p *models.StudentHealthProfile) *HealthProfileResponse {
	return &HealthProfileResponse{
		ID:                    p.ID.String(),
		StudentID:             p.StudentID.String(),
		BloodGroup:            p.BloodGroup,
		HeightCm:              p.HeightCm,
		WeightKg:              p.WeightKg,
		VisionLeft:            p.VisionLeft,
		VisionRight:           p.VisionRight,
		HearingStatus:         p.HearingStatus,
		MedicalNotes:          p.MedicalNotes,
		InsuranceProvider:     p.InsuranceProvider,
		InsurancePolicyNumber: p.InsurancePolicyNumber,
		InsuranceExpiry:       formatDate(p.InsuranceExpiry),
		PreferredHospital:     p.PreferredHospital,
		FamilyDoctorName:      p.FamilyDoctorName,
		FamilyDoctorPhone:     p.FamilyDoctorPhone,
		LastCheckupDate:       formatDate(p.LastCheckupDate),
		CreatedAt:             p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:             p.UpdatedAt.Format(time.RFC3339),
	}
}

// =============================================================================
// Allergy Service Methods
// =============================================================================

func (s *Service) ListAllergies(ctx context.Context, tenantID, studentID uuid.UUID, activeOnly bool) (*AllergyListResponse, error) {
	allergies, err := s.repo.ListAllergies(ctx, tenantID, studentID, activeOnly)
	if err != nil {
		return nil, err
	}

	responses := make([]AllergyResponse, len(allergies))
	for i, a := range allergies {
		responses[i] = s.toAllergyResponse(&a)
	}

	return &AllergyListResponse{
		Allergies: responses,
		Total:     len(responses),
	}, nil
}

func (s *Service) GetAllergy(ctx context.Context, tenantID, allergyID uuid.UUID) (*AllergyResponse, error) {
	allergy, err := s.repo.GetAllergy(ctx, tenantID, allergyID)
	if err != nil {
		return nil, err
	}
	if allergy == nil {
		return nil, ErrAllergyNotFound
	}
	resp := s.toAllergyResponse(allergy)
	return &resp, nil
}

func (s *Service) CreateAllergy(ctx context.Context, tenantID, studentID, userID uuid.UUID, req *CreateAllergyRequest) (*AllergyResponse, error) {
	// Validate
	if !models.AllergyType(req.AllergyType).IsValid() {
		return nil, ErrInvalidAllergyType
	}
	if !models.AllergySeverity(req.Severity).IsValid() {
		return nil, ErrInvalidSeverity
	}

	allergy := &models.StudentAllergy{
		ID:                    uuid.New(),
		TenantID:              tenantID,
		StudentID:             studentID,
		Allergen:              req.Allergen,
		AllergyType:           models.AllergyType(req.AllergyType),
		Severity:              models.AllergySeverity(req.Severity),
		ReactionDescription:   req.ReactionDescription,
		TreatmentInstructions: req.TreatmentInstructions,
		EmergencyMedication:   req.EmergencyMedication,
		IsActive:              true,
		Notes:                 req.Notes,
		CreatedBy:             &userID,
		UpdatedBy:             &userID,
	}

	if req.DiagnosedDate != nil {
		if t, err := time.Parse("2006-01-02", *req.DiagnosedDate); err == nil {
			allergy.DiagnosedDate = &t
		}
	}

	if err := s.repo.CreateAllergy(ctx, allergy); err != nil {
		return nil, err
	}

	resp := s.toAllergyResponse(allergy)
	return &resp, nil
}

func (s *Service) UpdateAllergy(ctx context.Context, tenantID, allergyID, userID uuid.UUID, req *UpdateAllergyRequest) (*AllergyResponse, error) {
	allergy, err := s.repo.GetAllergy(ctx, tenantID, allergyID)
	if err != nil {
		return nil, err
	}
	if allergy == nil {
		return nil, ErrAllergyNotFound
	}

	if req.Allergen != nil {
		allergy.Allergen = *req.Allergen
	}
	if req.AllergyType != nil {
		if !models.AllergyType(*req.AllergyType).IsValid() {
			return nil, ErrInvalidAllergyType
		}
		allergy.AllergyType = models.AllergyType(*req.AllergyType)
	}
	if req.Severity != nil {
		if !models.AllergySeverity(*req.Severity).IsValid() {
			return nil, ErrInvalidSeverity
		}
		allergy.Severity = models.AllergySeverity(*req.Severity)
	}
	if req.ReactionDescription != nil {
		allergy.ReactionDescription = req.ReactionDescription
	}
	if req.TreatmentInstructions != nil {
		allergy.TreatmentInstructions = req.TreatmentInstructions
	}
	if req.EmergencyMedication != nil {
		allergy.EmergencyMedication = req.EmergencyMedication
	}
	if req.DiagnosedDate != nil {
		if t, err := time.Parse("2006-01-02", *req.DiagnosedDate); err == nil {
			allergy.DiagnosedDate = &t
		}
	}
	if req.IsActive != nil {
		allergy.IsActive = *req.IsActive
	}
	if req.Notes != nil {
		allergy.Notes = req.Notes
	}

	allergy.UpdatedBy = &userID
	allergy.UpdatedAt = time.Now()

	if err := s.repo.UpdateAllergy(ctx, allergy); err != nil {
		return nil, err
	}

	resp := s.toAllergyResponse(allergy)
	return &resp, nil
}

func (s *Service) DeleteAllergy(ctx context.Context, tenantID, allergyID uuid.UUID) error {
	allergy, err := s.repo.GetAllergy(ctx, tenantID, allergyID)
	if err != nil {
		return err
	}
	if allergy == nil {
		return ErrAllergyNotFound
	}
	return s.repo.DeleteAllergy(ctx, tenantID, allergyID)
}

func (s *Service) toAllergyResponse(a *models.StudentAllergy) AllergyResponse {
	return AllergyResponse{
		ID:                    a.ID.String(),
		StudentID:             a.StudentID.String(),
		Allergen:              a.Allergen,
		AllergyType:           string(a.AllergyType),
		Severity:              string(a.Severity),
		ReactionDescription:   a.ReactionDescription,
		TreatmentInstructions: a.TreatmentInstructions,
		EmergencyMedication:   a.EmergencyMedication,
		DiagnosedDate:         formatDate(a.DiagnosedDate),
		IsActive:              a.IsActive,
		Notes:                 a.Notes,
		CreatedAt:             a.CreatedAt.Format(time.RFC3339),
		UpdatedAt:             a.UpdatedAt.Format(time.RFC3339),
	}
}

// =============================================================================
// Chronic Condition Service Methods
// =============================================================================

func (s *Service) ListConditions(ctx context.Context, tenantID, studentID uuid.UUID, activeOnly bool) (*ConditionListResponse, error) {
	conditions, err := s.repo.ListConditions(ctx, tenantID, studentID, activeOnly)
	if err != nil {
		return nil, err
	}

	responses := make([]ConditionResponse, len(conditions))
	for i, c := range conditions {
		responses[i] = s.toConditionResponse(&c)
	}

	return &ConditionListResponse{
		Conditions: responses,
		Total:      len(responses),
	}, nil
}

func (s *Service) GetCondition(ctx context.Context, tenantID, conditionID uuid.UUID) (*ConditionResponse, error) {
	condition, err := s.repo.GetCondition(ctx, tenantID, conditionID)
	if err != nil {
		return nil, err
	}
	if condition == nil {
		return nil, ErrConditionNotFound
	}
	resp := s.toConditionResponse(condition)
	return &resp, nil
}

func (s *Service) CreateCondition(ctx context.Context, tenantID, studentID, userID uuid.UUID, req *CreateConditionRequest) (*ConditionResponse, error) {
	if !models.ConditionType(req.ConditionType).IsValid() {
		return nil, ErrInvalidConditionType
	}
	if !models.ConditionSeverity(req.Severity).IsValid() {
		return nil, ErrInvalidSeverity
	}

	condition := &models.StudentChronicCondition{
		ID:             uuid.New(),
		TenantID:       tenantID,
		StudentID:      studentID,
		ConditionName:  req.ConditionName,
		ConditionType:  models.ConditionType(req.ConditionType),
		Severity:       models.ConditionSeverity(req.Severity),
		ManagementPlan: req.ManagementPlan,
		Restrictions:   req.Restrictions,
		Triggers:       req.Triggers,
		IsActive:       true,
		Notes:          req.Notes,
		CreatedBy:      &userID,
		UpdatedBy:      &userID,
	}

	if req.DiagnosedDate != nil {
		if t, err := time.Parse("2006-01-02", *req.DiagnosedDate); err == nil {
			condition.DiagnosedDate = &t
		}
	}

	if err := s.repo.CreateCondition(ctx, condition); err != nil {
		return nil, err
	}

	resp := s.toConditionResponse(condition)
	return &resp, nil
}

func (s *Service) UpdateCondition(ctx context.Context, tenantID, conditionID, userID uuid.UUID, req *UpdateConditionRequest) (*ConditionResponse, error) {
	condition, err := s.repo.GetCondition(ctx, tenantID, conditionID)
	if err != nil {
		return nil, err
	}
	if condition == nil {
		return nil, ErrConditionNotFound
	}

	if req.ConditionName != nil {
		condition.ConditionName = *req.ConditionName
	}
	if req.ConditionType != nil {
		if !models.ConditionType(*req.ConditionType).IsValid() {
			return nil, ErrInvalidConditionType
		}
		condition.ConditionType = models.ConditionType(*req.ConditionType)
	}
	if req.Severity != nil {
		if !models.ConditionSeverity(*req.Severity).IsValid() {
			return nil, ErrInvalidSeverity
		}
		condition.Severity = models.ConditionSeverity(*req.Severity)
	}
	if req.ManagementPlan != nil {
		condition.ManagementPlan = req.ManagementPlan
	}
	if req.Restrictions != nil {
		condition.Restrictions = req.Restrictions
	}
	if req.Triggers != nil {
		condition.Triggers = req.Triggers
	}
	if req.DiagnosedDate != nil {
		if t, err := time.Parse("2006-01-02", *req.DiagnosedDate); err == nil {
			condition.DiagnosedDate = &t
		}
	}
	if req.IsActive != nil {
		condition.IsActive = *req.IsActive
	}
	if req.Notes != nil {
		condition.Notes = req.Notes
	}

	condition.UpdatedBy = &userID
	condition.UpdatedAt = time.Now()

	if err := s.repo.UpdateCondition(ctx, condition); err != nil {
		return nil, err
	}

	resp := s.toConditionResponse(condition)
	return &resp, nil
}

func (s *Service) DeleteCondition(ctx context.Context, tenantID, conditionID uuid.UUID) error {
	condition, err := s.repo.GetCondition(ctx, tenantID, conditionID)
	if err != nil {
		return err
	}
	if condition == nil {
		return ErrConditionNotFound
	}
	return s.repo.DeleteCondition(ctx, tenantID, conditionID)
}

func (s *Service) toConditionResponse(c *models.StudentChronicCondition) ConditionResponse {
	return ConditionResponse{
		ID:             c.ID.String(),
		StudentID:      c.StudentID.String(),
		ConditionName:  c.ConditionName,
		ConditionType:  string(c.ConditionType),
		Severity:       string(c.Severity),
		ManagementPlan: c.ManagementPlan,
		Restrictions:   c.Restrictions,
		Triggers:       c.Triggers,
		DiagnosedDate:  formatDate(c.DiagnosedDate),
		IsActive:       c.IsActive,
		Notes:          c.Notes,
		CreatedAt:      c.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      c.UpdatedAt.Format(time.RFC3339),
	}
}

// =============================================================================
// Medication Service Methods
// =============================================================================

func (s *Service) ListMedications(ctx context.Context, tenantID, studentID uuid.UUID, activeOnly bool) (*MedicationListResponse, error) {
	medications, err := s.repo.ListMedications(ctx, tenantID, studentID, activeOnly)
	if err != nil {
		return nil, err
	}

	responses := make([]MedicationResponse, len(medications))
	for i, m := range medications {
		responses[i] = s.toMedicationResponse(&m)
	}

	return &MedicationListResponse{
		Medications: responses,
		Total:       len(responses),
	}, nil
}

func (s *Service) GetMedication(ctx context.Context, tenantID, medicationID uuid.UUID) (*MedicationResponse, error) {
	medication, err := s.repo.GetMedication(ctx, tenantID, medicationID)
	if err != nil {
		return nil, err
	}
	if medication == nil {
		return nil, ErrMedicationNotFound
	}
	resp := s.toMedicationResponse(medication)
	return &resp, nil
}

func (s *Service) CreateMedication(ctx context.Context, tenantID, studentID, userID uuid.UUID, req *CreateMedicationRequest) (*MedicationResponse, error) {
	if !models.MedicationFrequency(req.Frequency).IsValid() {
		return nil, ErrInvalidFrequency
	}
	if !models.MedicationRoute(req.Route).IsValid() {
		return nil, ErrInvalidRoute
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, err
	}

	medication := &models.StudentMedication{
		ID:                       uuid.New(),
		TenantID:                 tenantID,
		StudentID:                studentID,
		MedicationName:           req.MedicationName,
		Dosage:                   req.Dosage,
		Frequency:                models.MedicationFrequency(req.Frequency),
		Route:                    models.MedicationRoute(req.Route),
		Purpose:                  req.Purpose,
		SpecialInstructions:      req.SpecialInstructions,
		StartDate:                startDate,
		AdministeredAtSchool:     req.AdministeredAtSchool,
		SchoolAdministrationTime: req.SchoolAdministrationTime,
		PrescribingDoctor:        req.PrescribingDoctor,
		IsActive:                 true,
		Notes:                    req.Notes,
		CreatedBy:                &userID,
		UpdatedBy:                &userID,
	}

	if req.EndDate != nil {
		if t, err := time.Parse("2006-01-02", *req.EndDate); err == nil {
			if t.Before(startDate) {
				return nil, ErrInvalidDateRange
			}
			medication.EndDate = &t
		}
	}
	if req.PrescriptionDate != nil {
		if t, err := time.Parse("2006-01-02", *req.PrescriptionDate); err == nil {
			medication.PrescriptionDate = &t
		}
	}

	if err := s.repo.CreateMedication(ctx, medication); err != nil {
		return nil, err
	}

	resp := s.toMedicationResponse(medication)
	return &resp, nil
}

func (s *Service) UpdateMedication(ctx context.Context, tenantID, medicationID, userID uuid.UUID, req *UpdateMedicationRequest) (*MedicationResponse, error) {
	medication, err := s.repo.GetMedication(ctx, tenantID, medicationID)
	if err != nil {
		return nil, err
	}
	if medication == nil {
		return nil, ErrMedicationNotFound
	}

	if req.MedicationName != nil {
		medication.MedicationName = *req.MedicationName
	}
	if req.Dosage != nil {
		medication.Dosage = *req.Dosage
	}
	if req.Frequency != nil {
		if !models.MedicationFrequency(*req.Frequency).IsValid() {
			return nil, ErrInvalidFrequency
		}
		medication.Frequency = models.MedicationFrequency(*req.Frequency)
	}
	if req.Route != nil {
		if !models.MedicationRoute(*req.Route).IsValid() {
			return nil, ErrInvalidRoute
		}
		medication.Route = models.MedicationRoute(*req.Route)
	}
	if req.Purpose != nil {
		medication.Purpose = req.Purpose
	}
	if req.SpecialInstructions != nil {
		medication.SpecialInstructions = req.SpecialInstructions
	}
	if req.StartDate != nil {
		if t, err := time.Parse("2006-01-02", *req.StartDate); err == nil {
			medication.StartDate = t
		}
	}
	if req.EndDate != nil {
		if t, err := time.Parse("2006-01-02", *req.EndDate); err == nil {
			medication.EndDate = &t
		}
	}
	if req.AdministeredAtSchool != nil {
		medication.AdministeredAtSchool = *req.AdministeredAtSchool
	}
	if req.SchoolAdministrationTime != nil {
		medication.SchoolAdministrationTime = req.SchoolAdministrationTime
	}
	if req.PrescribingDoctor != nil {
		medication.PrescribingDoctor = req.PrescribingDoctor
	}
	if req.PrescriptionDate != nil {
		if t, err := time.Parse("2006-01-02", *req.PrescriptionDate); err == nil {
			medication.PrescriptionDate = &t
		}
	}
	if req.IsActive != nil {
		medication.IsActive = *req.IsActive
	}
	if req.Notes != nil {
		medication.Notes = req.Notes
	}

	medication.UpdatedBy = &userID
	medication.UpdatedAt = time.Now()

	if err := s.repo.UpdateMedication(ctx, medication); err != nil {
		return nil, err
	}

	resp := s.toMedicationResponse(medication)
	return &resp, nil
}

func (s *Service) DeleteMedication(ctx context.Context, tenantID, medicationID uuid.UUID) error {
	medication, err := s.repo.GetMedication(ctx, tenantID, medicationID)
	if err != nil {
		return err
	}
	if medication == nil {
		return ErrMedicationNotFound
	}
	return s.repo.DeleteMedication(ctx, tenantID, medicationID)
}

func (s *Service) toMedicationResponse(m *models.StudentMedication) MedicationResponse {
	return MedicationResponse{
		ID:                       m.ID.String(),
		StudentID:                m.StudentID.String(),
		MedicationName:           m.MedicationName,
		Dosage:                   m.Dosage,
		Frequency:                string(m.Frequency),
		Route:                    string(m.Route),
		Purpose:                  m.Purpose,
		SpecialInstructions:      m.SpecialInstructions,
		StartDate:                m.StartDate.Format("2006-01-02"),
		EndDate:                  formatDate(m.EndDate),
		AdministeredAtSchool:     m.AdministeredAtSchool,
		SchoolAdministrationTime: m.SchoolAdministrationTime,
		PrescribingDoctor:        m.PrescribingDoctor,
		PrescriptionDate:         formatDate(m.PrescriptionDate),
		IsActive:                 m.IsActive,
		Notes:                    m.Notes,
		CreatedAt:                m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:                m.UpdatedAt.Format(time.RFC3339),
	}
}

// =============================================================================
// Vaccination Service Methods
// =============================================================================

func (s *Service) ListVaccinations(ctx context.Context, tenantID, studentID uuid.UUID) (*VaccinationListResponse, error) {
	vaccinations, err := s.repo.ListVaccinations(ctx, tenantID, studentID)
	if err != nil {
		return nil, err
	}

	responses := make([]VaccinationResponse, len(vaccinations))
	for i, v := range vaccinations {
		responses[i] = s.toVaccinationResponse(&v)
	}

	return &VaccinationListResponse{
		Vaccinations: responses,
		Total:        len(responses),
	}, nil
}

func (s *Service) GetVaccination(ctx context.Context, tenantID, vaccinationID uuid.UUID) (*VaccinationResponse, error) {
	vaccination, err := s.repo.GetVaccination(ctx, tenantID, vaccinationID)
	if err != nil {
		return nil, err
	}
	if vaccination == nil {
		return nil, ErrVaccinationNotFound
	}
	resp := s.toVaccinationResponse(vaccination)
	return &resp, nil
}

func (s *Service) CreateVaccination(ctx context.Context, tenantID, studentID, userID uuid.UUID, req *CreateVaccinationRequest) (*VaccinationResponse, error) {
	if req.VaccineType != nil && !models.VaccineType(*req.VaccineType).IsValid() {
		return nil, ErrInvalidVaccineType
	}
	if req.DoseNumber < 1 {
		req.DoseNumber = 1
	}

	administeredDate, err := time.Parse("2006-01-02", req.AdministeredDate)
	if err != nil {
		return nil, err
	}

	vaccination := &models.StudentVaccination{
		ID:                  uuid.New(),
		TenantID:            tenantID,
		StudentID:           studentID,
		VaccineName:         req.VaccineName,
		DoseNumber:          req.DoseNumber,
		AdministeredDate:    administeredDate,
		AdministeredBy:      req.AdministeredBy,
		AdministrationSite:  req.AdministrationSite,
		BatchNumber:         req.BatchNumber,
		HadReaction:         req.HadReaction,
		ReactionDescription: req.ReactionDescription,
		CertificateUrl:      req.CertificateUrl,
		IsVerified:          false,
		Notes:               req.Notes,
		CreatedBy:           &userID,
		UpdatedBy:           &userID,
	}

	if req.VaccineType != nil {
		vt := models.VaccineType(*req.VaccineType)
		vaccination.VaccineType = &vt
	}
	if req.NextDueDate != nil {
		if t, err := time.Parse("2006-01-02", *req.NextDueDate); err == nil {
			vaccination.NextDueDate = &t
		}
	}

	if err := s.repo.CreateVaccination(ctx, vaccination); err != nil {
		return nil, err
	}

	resp := s.toVaccinationResponse(vaccination)
	return &resp, nil
}

func (s *Service) UpdateVaccination(ctx context.Context, tenantID, vaccinationID, userID uuid.UUID, req *UpdateVaccinationRequest) (*VaccinationResponse, error) {
	vaccination, err := s.repo.GetVaccination(ctx, tenantID, vaccinationID)
	if err != nil {
		return nil, err
	}
	if vaccination == nil {
		return nil, ErrVaccinationNotFound
	}

	if req.VaccineName != nil {
		vaccination.VaccineName = *req.VaccineName
	}
	if req.VaccineType != nil {
		if !models.VaccineType(*req.VaccineType).IsValid() {
			return nil, ErrInvalidVaccineType
		}
		vt := models.VaccineType(*req.VaccineType)
		vaccination.VaccineType = &vt
	}
	if req.DoseNumber != nil {
		if *req.DoseNumber < 1 {
			return nil, ErrInvalidDoseNumber
		}
		vaccination.DoseNumber = *req.DoseNumber
	}
	if req.AdministeredDate != nil {
		if t, err := time.Parse("2006-01-02", *req.AdministeredDate); err == nil {
			vaccination.AdministeredDate = t
		}
	}
	if req.AdministeredBy != nil {
		vaccination.AdministeredBy = req.AdministeredBy
	}
	if req.AdministrationSite != nil {
		vaccination.AdministrationSite = req.AdministrationSite
	}
	if req.BatchNumber != nil {
		vaccination.BatchNumber = req.BatchNumber
	}
	if req.NextDueDate != nil {
		if t, err := time.Parse("2006-01-02", *req.NextDueDate); err == nil {
			vaccination.NextDueDate = &t
		}
	}
	if req.HadReaction != nil {
		vaccination.HadReaction = *req.HadReaction
	}
	if req.ReactionDescription != nil {
		vaccination.ReactionDescription = req.ReactionDescription
	}
	if req.CertificateUrl != nil {
		vaccination.CertificateUrl = req.CertificateUrl
	}
	if req.Notes != nil {
		vaccination.Notes = req.Notes
	}

	vaccination.UpdatedBy = &userID
	vaccination.UpdatedAt = time.Now()

	if err := s.repo.UpdateVaccination(ctx, vaccination); err != nil {
		return nil, err
	}

	resp := s.toVaccinationResponse(vaccination)
	return &resp, nil
}

func (s *Service) DeleteVaccination(ctx context.Context, tenantID, vaccinationID uuid.UUID) error {
	vaccination, err := s.repo.GetVaccination(ctx, tenantID, vaccinationID)
	if err != nil {
		return err
	}
	if vaccination == nil {
		return ErrVaccinationNotFound
	}
	return s.repo.DeleteVaccination(ctx, tenantID, vaccinationID)
}

func (s *Service) toVaccinationResponse(v *models.StudentVaccination) VaccinationResponse {
	resp := VaccinationResponse{
		ID:                  v.ID.String(),
		StudentID:           v.StudentID.String(),
		VaccineName:         v.VaccineName,
		DoseNumber:          v.DoseNumber,
		AdministeredDate:    v.AdministeredDate.Format("2006-01-02"),
		AdministeredBy:      v.AdministeredBy,
		AdministrationSite:  v.AdministrationSite,
		BatchNumber:         v.BatchNumber,
		NextDueDate:         formatDate(v.NextDueDate),
		HadReaction:         v.HadReaction,
		ReactionDescription: v.ReactionDescription,
		CertificateUrl:      v.CertificateUrl,
		IsVerified:          v.IsVerified,
		VerifiedAt:          formatDateTime(v.VerifiedAt),
		Notes:               v.Notes,
		CreatedAt:           v.CreatedAt.Format(time.RFC3339),
		UpdatedAt:           v.UpdatedAt.Format(time.RFC3339),
	}
	if v.VaccineType != nil {
		s := string(*v.VaccineType)
		resp.VaccineType = &s
	}
	return resp
}

// =============================================================================
// Medical Incident Service Methods
// =============================================================================

func (s *Service) ListIncidents(ctx context.Context, tenantID, studentID uuid.UUID, limit int) (*IncidentListResponse, error) {
	incidents, err := s.repo.ListIncidents(ctx, tenantID, studentID, limit)
	if err != nil {
		return nil, err
	}

	responses := make([]IncidentResponse, len(incidents))
	for i, inc := range incidents {
		responses[i] = s.toIncidentResponse(&inc)
	}

	return &IncidentListResponse{
		Incidents: responses,
		Total:     len(responses),
	}, nil
}

func (s *Service) GetIncident(ctx context.Context, tenantID, incidentID uuid.UUID) (*IncidentResponse, error) {
	incident, err := s.repo.GetIncident(ctx, tenantID, incidentID)
	if err != nil {
		return nil, err
	}
	if incident == nil {
		return nil, ErrIncidentNotFound
	}
	resp := s.toIncidentResponse(incident)
	return &resp, nil
}

func (s *Service) CreateIncident(ctx context.Context, tenantID, studentID, userID uuid.UUID, req *CreateIncidentRequest) (*IncidentResponse, error) {
	if !models.IncidentType(req.IncidentType).IsValid() {
		return nil, ErrInvalidIncidentType
	}

	incidentDate, err := time.Parse("2006-01-02", req.IncidentDate)
	if err != nil {
		return nil, err
	}

	incident := &models.StudentMedicalIncident{
		ID:                    uuid.New(),
		TenantID:              tenantID,
		StudentID:             studentID,
		IncidentDate:          incidentDate,
		IncidentTime:          req.IncidentTime,
		Location:              req.Location,
		IncidentType:          models.IncidentType(req.IncidentType),
		Description:           req.Description,
		Symptoms:              req.Symptoms,
		FirstAidGiven:         req.FirstAidGiven,
		FirstAidDescription:   req.FirstAidDescription,
		ActionTaken:           req.ActionTaken,
		ParentNotified:        req.ParentNotified,
		HospitalVisitRequired: req.HospitalVisitRequired,
		HospitalName:          req.HospitalName,
		StudentSentHome:       req.StudentSentHome,
		ReturnToClassTime:     req.ReturnToClassTime,
		FollowUpRequired:      req.FollowUpRequired,
		FollowUpNotes:         req.FollowUpNotes,
		Outcome:               req.Outcome,
		ReportedBy:            userID,
		Notes:                 req.Notes,
	}

	if req.ParentNotifiedAt != nil {
		if t, err := time.Parse(time.RFC3339, *req.ParentNotifiedAt); err == nil {
			incident.ParentNotifiedAt = &t
			incident.ParentNotifiedBy = &userID
		}
	}
	if req.HospitalVisitDate != nil {
		if t, err := time.Parse("2006-01-02", *req.HospitalVisitDate); err == nil {
			incident.HospitalVisitDate = &t
		}
	}
	if req.FollowUpDate != nil {
		if t, err := time.Parse("2006-01-02", *req.FollowUpDate); err == nil {
			incident.FollowUpDate = &t
		}
	}

	if err := s.repo.CreateIncident(ctx, incident); err != nil {
		return nil, err
	}

	resp := s.toIncidentResponse(incident)
	return &resp, nil
}

func (s *Service) UpdateIncident(ctx context.Context, tenantID, incidentID, userID uuid.UUID, req *UpdateIncidentRequest) (*IncidentResponse, error) {
	incident, err := s.repo.GetIncident(ctx, tenantID, incidentID)
	if err != nil {
		return nil, err
	}
	if incident == nil {
		return nil, ErrIncidentNotFound
	}

	if req.IncidentDate != nil {
		if t, err := time.Parse("2006-01-02", *req.IncidentDate); err == nil {
			incident.IncidentDate = t
		}
	}
	if req.IncidentTime != nil {
		incident.IncidentTime = *req.IncidentTime
	}
	if req.Location != nil {
		incident.Location = req.Location
	}
	if req.IncidentType != nil {
		if !models.IncidentType(*req.IncidentType).IsValid() {
			return nil, ErrInvalidIncidentType
		}
		incident.IncidentType = models.IncidentType(*req.IncidentType)
	}
	if req.Description != nil {
		incident.Description = *req.Description
	}
	if req.Symptoms != nil {
		incident.Symptoms = req.Symptoms
	}
	if req.FirstAidGiven != nil {
		incident.FirstAidGiven = *req.FirstAidGiven
	}
	if req.FirstAidDescription != nil {
		incident.FirstAidDescription = req.FirstAidDescription
	}
	if req.ActionTaken != nil {
		incident.ActionTaken = *req.ActionTaken
	}
	if req.ParentNotified != nil {
		incident.ParentNotified = *req.ParentNotified
	}
	if req.ParentNotifiedAt != nil {
		if t, err := time.Parse(time.RFC3339, *req.ParentNotifiedAt); err == nil {
			incident.ParentNotifiedAt = &t
		}
	}
	if req.HospitalVisitRequired != nil {
		incident.HospitalVisitRequired = *req.HospitalVisitRequired
	}
	if req.HospitalName != nil {
		incident.HospitalName = req.HospitalName
	}
	if req.HospitalVisitDate != nil {
		if t, err := time.Parse("2006-01-02", *req.HospitalVisitDate); err == nil {
			incident.HospitalVisitDate = &t
		}
	}
	if req.StudentSentHome != nil {
		incident.StudentSentHome = *req.StudentSentHome
	}
	if req.ReturnToClassTime != nil {
		incident.ReturnToClassTime = req.ReturnToClassTime
	}
	if req.FollowUpRequired != nil {
		incident.FollowUpRequired = *req.FollowUpRequired
	}
	if req.FollowUpDate != nil {
		if t, err := time.Parse("2006-01-02", *req.FollowUpDate); err == nil {
			incident.FollowUpDate = &t
		}
	}
	if req.FollowUpNotes != nil {
		incident.FollowUpNotes = req.FollowUpNotes
	}
	if req.Outcome != nil {
		incident.Outcome = req.Outcome
	}
	if req.Notes != nil {
		incident.Notes = req.Notes
	}

	incident.UpdatedBy = &userID
	incident.UpdatedAt = time.Now()

	if err := s.repo.UpdateIncident(ctx, incident); err != nil {
		return nil, err
	}

	resp := s.toIncidentResponse(incident)
	return &resp, nil
}

func (s *Service) DeleteIncident(ctx context.Context, tenantID, incidentID uuid.UUID) error {
	incident, err := s.repo.GetIncident(ctx, tenantID, incidentID)
	if err != nil {
		return err
	}
	if incident == nil {
		return ErrIncidentNotFound
	}
	return s.repo.DeleteIncident(ctx, tenantID, incidentID)
}

func (s *Service) toIncidentResponse(i *models.StudentMedicalIncident) IncidentResponse {
	return IncidentResponse{
		ID:                    i.ID.String(),
		StudentID:             i.StudentID.String(),
		IncidentDate:          i.IncidentDate.Format("2006-01-02"),
		IncidentTime:          i.IncidentTime,
		Location:              i.Location,
		IncidentType:          string(i.IncidentType),
		Description:           i.Description,
		Symptoms:              i.Symptoms,
		FirstAidGiven:         i.FirstAidGiven,
		FirstAidDescription:   i.FirstAidDescription,
		ActionTaken:           i.ActionTaken,
		ParentNotified:        i.ParentNotified,
		ParentNotifiedAt:      formatDateTime(i.ParentNotifiedAt),
		HospitalVisitRequired: i.HospitalVisitRequired,
		HospitalName:          i.HospitalName,
		HospitalVisitDate:     formatDate(i.HospitalVisitDate),
		StudentSentHome:       i.StudentSentHome,
		ReturnToClassTime:     i.ReturnToClassTime,
		FollowUpRequired:      i.FollowUpRequired,
		FollowUpDate:          formatDate(i.FollowUpDate),
		FollowUpNotes:         i.FollowUpNotes,
		Outcome:               i.Outcome,
		ReportedBy:            i.ReportedBy.String(),
		Notes:                 i.Notes,
		CreatedAt:             i.CreatedAt.Format(time.RFC3339),
		UpdatedAt:             i.UpdatedAt.Format(time.RFC3339),
	}
}

// =============================================================================
// Health Summary
// =============================================================================

func (s *Service) GetHealthSummary(ctx context.Context, tenantID, studentID uuid.UUID) (*HealthSummaryResponse, error) {
	profile, _ := s.GetHealthProfile(ctx, tenantID, studentID)

	allergies, _ := s.ListAllergies(ctx, tenantID, studentID, true)
	conditions, _ := s.ListConditions(ctx, tenantID, studentID, true)
	medications, _ := s.ListMedications(ctx, tenantID, studentID, true)
	vaccinations, _ := s.ListVaccinations(ctx, tenantID, studentID)
	incidents, _ := s.ListIncidents(ctx, tenantID, studentID, 5)

	summary := &HealthSummaryResponse{
		Profile:         profile,
		Allergies:       []AllergyResponse{},
		Conditions:      []ConditionResponse{},
		Medications:     []MedicationResponse{},
		Vaccinations:    []VaccinationResponse{},
		RecentIncidents: []IncidentResponse{},
	}

	if allergies != nil {
		summary.Allergies = allergies.Allergies
	}
	if conditions != nil {
		summary.Conditions = conditions.Conditions
	}
	if medications != nil {
		summary.Medications = medications.Medications
	}
	if vaccinations != nil {
		summary.Vaccinations = vaccinations.Vaccinations
	}
	if incidents != nil {
		summary.RecentIncidents = incidents.Incidents
	}

	return summary, nil
}
