package behavioral

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
)

// Repository handles database operations for behavioral incidents
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new behavioral repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// =============================================================================
// Incident Operations
// =============================================================================

// CreateIncident creates a new behavioral incident
func (r *Repository) CreateIncident(ctx context.Context, incident *models.StudentBehavioralIncident) error {
	return r.db.WithContext(ctx).Create(incident).Error
}

// GetIncident retrieves an incident by ID
func (r *Repository) GetIncident(ctx context.Context, tenantID, incidentID uuid.UUID) (*models.StudentBehavioralIncident, error) {
	var incident models.StudentBehavioralIncident
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, incidentID).
		Preload("Reporter").
		Preload("FollowUps").
		First(&incident).Error
	if err != nil {
		return nil, err
	}
	return &incident, nil
}

// GetIncidentByStudent retrieves an incident for a specific student
func (r *Repository) GetIncidentByStudent(ctx context.Context, tenantID, studentID, incidentID uuid.UUID) (*models.StudentBehavioralIncident, error) {
	var incident models.StudentBehavioralIncident
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND student_id = ? AND id = ?", tenantID, studentID, incidentID).
		Preload("Reporter").
		Preload("FollowUps").
		First(&incident).Error
	if err != nil {
		return nil, err
	}
	return &incident, nil
}

// ListIncidents retrieves incidents for a student with optional filtering
func (r *Repository) ListIncidents(ctx context.Context, tenantID, studentID uuid.UUID, filter IncidentFilter) ([]models.StudentBehavioralIncident, int64, error) {
	query := r.db.WithContext(ctx).
		Where("tenant_id = ? AND student_id = ?", tenantID, studentID)

	if filter.IncidentType != nil {
		query = query.Where("incident_type = ?", *filter.IncidentType)
	}
	if filter.Severity != nil {
		query = query.Where("severity = ?", *filter.Severity)
	}
	if filter.DateFrom != nil {
		query = query.Where("incident_date >= ?", *filter.DateFrom)
	}
	if filter.DateTo != nil {
		query = query.Where("incident_date <= ?", *filter.DateTo)
	}

	var total int64
	if err := query.Model(&models.StudentBehavioralIncident{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var incidents []models.StudentBehavioralIncident
	query = query.Order("incident_date DESC, incident_time DESC").
		Preload("Reporter").
		Preload("FollowUps")

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	if err := query.Find(&incidents).Error; err != nil {
		return nil, 0, err
	}

	return incidents, total, nil
}

// UpdateIncident updates an incident
func (r *Repository) UpdateIncident(ctx context.Context, incident *models.StudentBehavioralIncident) error {
	return r.db.WithContext(ctx).Save(incident).Error
}

// DeleteIncident deletes an incident
func (r *Repository) DeleteIncident(ctx context.Context, tenantID, incidentID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, incidentID).
		Delete(&models.StudentBehavioralIncident{}).Error
}

// =============================================================================
// Follow-Up Operations
// =============================================================================

// CreateFollowUp creates a new follow-up
func (r *Repository) CreateFollowUp(ctx context.Context, followUp *models.IncidentFollowUp) error {
	return r.db.WithContext(ctx).Create(followUp).Error
}

// GetFollowUp retrieves a follow-up by ID
func (r *Repository) GetFollowUp(ctx context.Context, tenantID, followUpID uuid.UUID) (*models.IncidentFollowUp, error) {
	var followUp models.IncidentFollowUp
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, followUpID).
		First(&followUp).Error
	if err != nil {
		return nil, err
	}
	return &followUp, nil
}

// UpdateFollowUp updates a follow-up
func (r *Repository) UpdateFollowUp(ctx context.Context, followUp *models.IncidentFollowUp) error {
	return r.db.WithContext(ctx).Save(followUp).Error
}

// DeleteFollowUp deletes a follow-up
func (r *Repository) DeleteFollowUp(ctx context.Context, tenantID, followUpID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, followUpID).
		Delete(&models.IncidentFollowUp{}).Error
}

// ListPendingFollowUps retrieves all pending follow-ups for a tenant
func (r *Repository) ListPendingFollowUps(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]models.IncidentFollowUp, int64, error) {
	query := r.db.WithContext(ctx).
		Where("tenant_id = ? AND status = ?", tenantID, models.FollowUpStatusPending)

	var total int64
	if err := query.Model(&models.IncidentFollowUp{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var followUps []models.IncidentFollowUp
	query = query.Order("scheduled_date ASC").
		Preload("Incident").
		Preload("Incident.Student")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&followUps).Error; err != nil {
		return nil, 0, err
	}

	return followUps, total, nil
}

// =============================================================================
// Summary Operations
// =============================================================================

// GetBehaviorSummary retrieves behavior statistics for a student
func (r *Repository) GetBehaviorSummary(ctx context.Context, tenantID, studentID uuid.UUID) (*BehaviorSummary, error) {
	summary := &BehaviorSummary{}

	// Total counts by type
	type TypeCount struct {
		IncidentType models.BehavioralIncidentType
		Count        int64
	}
	var typeCounts []TypeCount
	if err := r.db.WithContext(ctx).
		Model(&models.StudentBehavioralIncident{}).
		Select("incident_type, COUNT(*) as count").
		Where("tenant_id = ? AND student_id = ?", tenantID, studentID).
		Group("incident_type").
		Scan(&typeCounts).Error; err != nil {
		return nil, err
	}

	for _, tc := range typeCounts {
		switch tc.IncidentType {
		case models.BehavioralIncidentTypePositiveRecognition:
			summary.PositiveCount = int(tc.Count)
		case models.BehavioralIncidentTypeMinorInfraction:
			summary.MinorInfractionCount = int(tc.Count)
		case models.BehavioralIncidentTypeMajorViolation:
			summary.MajorViolationCount = int(tc.Count)
		}
		summary.TotalIncidents += int(tc.Count)
	}

	// This month count
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	var thisMonthCount int64
	if err := r.db.WithContext(ctx).
		Model(&models.StudentBehavioralIncident{}).
		Where("tenant_id = ? AND student_id = ? AND incident_date >= ?", tenantID, studentID, startOfMonth).
		Count(&thisMonthCount).Error; err != nil {
		return nil, err
	}
	summary.ThisMonthCount = int(thisMonthCount)

	// Last month count
	startOfLastMonth := startOfMonth.AddDate(0, -1, 0)
	var lastMonthCount int64
	if err := r.db.WithContext(ctx).
		Model(&models.StudentBehavioralIncident{}).
		Where("tenant_id = ? AND student_id = ? AND incident_date >= ? AND incident_date < ?", tenantID, studentID, startOfLastMonth, startOfMonth).
		Count(&lastMonthCount).Error; err != nil {
		return nil, err
	}
	summary.LastMonthCount = int(lastMonthCount)

	// Pending follow-ups count
	var pendingCount int64
	if err := r.db.WithContext(ctx).
		Model(&models.IncidentFollowUp{}).
		Joins("JOIN student_behavioral_incidents ON student_behavioral_incidents.id = incident_follow_ups.incident_id").
		Where("incident_follow_ups.tenant_id = ? AND student_behavioral_incidents.student_id = ? AND incident_follow_ups.status = ?",
			tenantID, studentID, models.FollowUpStatusPending).
		Count(&pendingCount).Error; err != nil {
		return nil, err
	}
	summary.PendingFollowUps = int(pendingCount)

	// Calculate trend
	summary.Trend = calculateTrend(summary)

	return summary, nil
}

// calculateTrend determines the behavioral trend based on incident ratios
func calculateTrend(summary *BehaviorSummary) string {
	if summary.TotalIncidents == 0 {
		return "stable"
	}

	positiveRatio := float64(summary.PositiveCount) / float64(summary.TotalIncidents)
	if positiveRatio > 0.7 {
		return "improving"
	} else if positiveRatio < 0.3 {
		return "declining"
	}
	return "stable"
}
