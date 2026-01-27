package health

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// =============================================================================
// Health Profile Repository Methods
// =============================================================================

func (r *Repository) GetHealthProfile(ctx context.Context, tenantID, studentID uuid.UUID) (*models.StudentHealthProfile, error) {
	var profile models.StudentHealthProfile
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND student_id = ?", tenantID, studentID).
		First(&profile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &profile, nil
}

func (r *Repository) CreateHealthProfile(ctx context.Context, profile *models.StudentHealthProfile) error {
	return r.db.WithContext(ctx).Create(profile).Error
}

func (r *Repository) UpdateHealthProfile(ctx context.Context, profile *models.StudentHealthProfile) error {
	return r.db.WithContext(ctx).Save(profile).Error
}

func (r *Repository) DeleteHealthProfile(ctx context.Context, tenantID, studentID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND student_id = ?", tenantID, studentID).
		Delete(&models.StudentHealthProfile{}).Error
}

// =============================================================================
// Allergy Repository Methods
// =============================================================================

func (r *Repository) ListAllergies(ctx context.Context, tenantID, studentID uuid.UUID, activeOnly bool) ([]models.StudentAllergy, error) {
	var allergies []models.StudentAllergy
	query := r.db.WithContext(ctx).
		Where("tenant_id = ? AND student_id = ?", tenantID, studentID)
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}
	err := query.Order("created_at DESC").Find(&allergies).Error
	return allergies, err
}

func (r *Repository) GetAllergy(ctx context.Context, tenantID, allergyID uuid.UUID) (*models.StudentAllergy, error) {
	var allergy models.StudentAllergy
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, allergyID).
		First(&allergy).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &allergy, nil
}

func (r *Repository) CreateAllergy(ctx context.Context, allergy *models.StudentAllergy) error {
	return r.db.WithContext(ctx).Create(allergy).Error
}

func (r *Repository) UpdateAllergy(ctx context.Context, allergy *models.StudentAllergy) error {
	return r.db.WithContext(ctx).Save(allergy).Error
}

func (r *Repository) DeleteAllergy(ctx context.Context, tenantID, allergyID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, allergyID).
		Delete(&models.StudentAllergy{}).Error
}

// =============================================================================
// Chronic Condition Repository Methods
// =============================================================================

func (r *Repository) ListConditions(ctx context.Context, tenantID, studentID uuid.UUID, activeOnly bool) ([]models.StudentChronicCondition, error) {
	var conditions []models.StudentChronicCondition
	query := r.db.WithContext(ctx).
		Where("tenant_id = ? AND student_id = ?", tenantID, studentID)
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}
	err := query.Order("created_at DESC").Find(&conditions).Error
	return conditions, err
}

func (r *Repository) GetCondition(ctx context.Context, tenantID, conditionID uuid.UUID) (*models.StudentChronicCondition, error) {
	var condition models.StudentChronicCondition
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, conditionID).
		First(&condition).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &condition, nil
}

func (r *Repository) CreateCondition(ctx context.Context, condition *models.StudentChronicCondition) error {
	return r.db.WithContext(ctx).Create(condition).Error
}

func (r *Repository) UpdateCondition(ctx context.Context, condition *models.StudentChronicCondition) error {
	return r.db.WithContext(ctx).Save(condition).Error
}

func (r *Repository) DeleteCondition(ctx context.Context, tenantID, conditionID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, conditionID).
		Delete(&models.StudentChronicCondition{}).Error
}

// =============================================================================
// Medication Repository Methods
// =============================================================================

func (r *Repository) ListMedications(ctx context.Context, tenantID, studentID uuid.UUID, activeOnly bool) ([]models.StudentMedication, error) {
	var medications []models.StudentMedication
	query := r.db.WithContext(ctx).
		Where("tenant_id = ? AND student_id = ?", tenantID, studentID)
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}
	err := query.Order("created_at DESC").Find(&medications).Error
	return medications, err
}

func (r *Repository) GetMedication(ctx context.Context, tenantID, medicationID uuid.UUID) (*models.StudentMedication, error) {
	var medication models.StudentMedication
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, medicationID).
		First(&medication).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &medication, nil
}

func (r *Repository) CreateMedication(ctx context.Context, medication *models.StudentMedication) error {
	return r.db.WithContext(ctx).Create(medication).Error
}

func (r *Repository) UpdateMedication(ctx context.Context, medication *models.StudentMedication) error {
	return r.db.WithContext(ctx).Save(medication).Error
}

func (r *Repository) DeleteMedication(ctx context.Context, tenantID, medicationID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, medicationID).
		Delete(&models.StudentMedication{}).Error
}

// =============================================================================
// Vaccination Repository Methods
// =============================================================================

func (r *Repository) ListVaccinations(ctx context.Context, tenantID, studentID uuid.UUID) ([]models.StudentVaccination, error) {
	var vaccinations []models.StudentVaccination
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND student_id = ?", tenantID, studentID).
		Order("administered_date DESC").Find(&vaccinations).Error
	return vaccinations, err
}

func (r *Repository) GetVaccination(ctx context.Context, tenantID, vaccinationID uuid.UUID) (*models.StudentVaccination, error) {
	var vaccination models.StudentVaccination
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, vaccinationID).
		First(&vaccination).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &vaccination, nil
}

func (r *Repository) CreateVaccination(ctx context.Context, vaccination *models.StudentVaccination) error {
	return r.db.WithContext(ctx).Create(vaccination).Error
}

func (r *Repository) UpdateVaccination(ctx context.Context, vaccination *models.StudentVaccination) error {
	return r.db.WithContext(ctx).Save(vaccination).Error
}

func (r *Repository) DeleteVaccination(ctx context.Context, tenantID, vaccinationID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, vaccinationID).
		Delete(&models.StudentVaccination{}).Error
}

// =============================================================================
// Medical Incident Repository Methods
// =============================================================================

func (r *Repository) ListIncidents(ctx context.Context, tenantID, studentID uuid.UUID, limit int) ([]models.StudentMedicalIncident, error) {
	var incidents []models.StudentMedicalIncident
	query := r.db.WithContext(ctx).
		Where("tenant_id = ? AND student_id = ?", tenantID, studentID).
		Order("incident_date DESC, incident_time DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&incidents).Error
	return incidents, err
}

func (r *Repository) GetIncident(ctx context.Context, tenantID, incidentID uuid.UUID) (*models.StudentMedicalIncident, error) {
	var incident models.StudentMedicalIncident
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, incidentID).
		First(&incident).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &incident, nil
}

func (r *Repository) CreateIncident(ctx context.Context, incident *models.StudentMedicalIncident) error {
	return r.db.WithContext(ctx).Create(incident).Error
}

func (r *Repository) UpdateIncident(ctx context.Context, incident *models.StudentMedicalIncident) error {
	return r.db.WithContext(ctx).Save(incident).Error
}

func (r *Repository) DeleteIncident(ctx context.Context, tenantID, incidentID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, incidentID).
		Delete(&models.StudentMedicalIncident{}).Error
}

// =============================================================================
// Student Lookup
// =============================================================================

func (r *Repository) StudentExists(ctx context.Context, tenantID, studentID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Student{}).
		Where("tenant_id = ? AND id = ?", tenantID, studentID).
		Count(&count).Error
	return count > 0, err
}
