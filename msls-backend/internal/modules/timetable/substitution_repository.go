package timetable

import (
	"context"
	"time"

	"msls-backend/internal/pkg/database/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ========================================
// Substitution Repository Methods
// ========================================

// ListSubstitutions returns substitutions based on filter criteria.
func (r *Repository) ListSubstitutions(ctx context.Context, filter SubstitutionFilter) ([]models.Substitution, int64, error) {
	var substitutions []models.Substitution
	var total int64

	query := r.db.WithContext(ctx).
		Model(&models.Substitution{}).
		Where("tenant_id = ?", filter.TenantID)

	if filter.BranchID != nil {
		query = query.Where("branch_id = ?", *filter.BranchID)
	}

	if filter.OriginalStaffID != nil {
		query = query.Where("original_staff_id = ?", *filter.OriginalStaffID)
	}

	if filter.SubstituteStaffID != nil {
		query = query.Where("substitute_staff_id = ?", *filter.SubstituteStaffID)
	}

	if filter.StartDate != nil {
		query = query.Where("substitution_date >= ?", *filter.StartDate)
	}

	if filter.EndDate != nil {
		query = query.Where("substitution_date <= ?", *filter.EndDate)
	}

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.
		Preload("Branch").
		Preload("OriginalStaff").
		Preload("SubstituteStaff").
		Preload("Creator").
		Preload("Periods").
		Preload("Periods.PeriodSlot").
		Preload("Periods.Subject").
		Preload("Periods.Section").
		Preload("Periods.Section.Class").
		Order("substitution_date DESC, created_at DESC")

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}

	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	err := query.Find(&substitutions).Error
	return substitutions, total, err
}

// GetSubstitutionByID returns a substitution by ID.
func (r *Repository) GetSubstitutionByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Substitution, error) {
	var substitution models.Substitution

	err := r.db.WithContext(ctx).
		Preload("Branch").
		Preload("OriginalStaff").
		Preload("SubstituteStaff").
		Preload("Creator").
		Preload("Approver").
		Preload("Periods").
		Preload("Periods.PeriodSlot").
		Preload("Periods.Subject").
		Preload("Periods.Section").
		Preload("Periods.Section.Class").
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&substitution).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrSubstitutionNotFound
		}
		return nil, err
	}

	return &substitution, nil
}

// CreateSubstitution creates a new substitution.
func (r *Repository) CreateSubstitution(ctx context.Context, substitution *models.Substitution) error {
	return r.db.WithContext(ctx).Create(substitution).Error
}

// CreateSubstitutionPeriods creates substitution periods.
func (r *Repository) CreateSubstitutionPeriods(ctx context.Context, periods []models.SubstitutionPeriod) error {
	if len(periods) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&periods).Error
}

// UpdateSubstitution updates a substitution.
func (r *Repository) UpdateSubstitution(ctx context.Context, substitution *models.Substitution) error {
	return r.db.WithContext(ctx).Save(substitution).Error
}

// DeleteSubstitution soft deletes a substitution.
func (r *Repository) DeleteSubstitution(ctx context.Context, tenantID, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&models.Substitution{}).Error
}

// DeleteSubstitutionPeriods deletes all periods for a substitution.
func (r *Repository) DeleteSubstitutionPeriods(ctx context.Context, substitutionID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("substitution_id = ?", substitutionID).
		Delete(&models.SubstitutionPeriod{}).Error
}

// CheckSubstitutionConflict checks if a substitution already exists for the same teacher, date, and period.
func (r *Repository) CheckSubstitutionConflict(ctx context.Context, tenantID, originalStaffID uuid.UUID, date time.Time, periodSlotIDs []uuid.UUID, excludeSubstitutionID *uuid.UUID) (bool, error) {
	query := r.db.WithContext(ctx).
		Model(&models.SubstitutionPeriod{}).
		Joins("JOIN substitutions ON substitutions.id = substitution_periods.substitution_id").
		Where("substitutions.tenant_id = ?", tenantID).
		Where("substitutions.original_staff_id = ?", originalStaffID).
		Where("substitutions.substitution_date = ?", date).
		Where("substitutions.status NOT IN ?", []string{"cancelled"}).
		Where("substitutions.deleted_at IS NULL").
		Where("substitution_periods.period_slot_id IN ?", periodSlotIDs)

	if excludeSubstitutionID != nil {
		query = query.Where("substitutions.id != ?", *excludeSubstitutionID)
	}

	var count int64
	err := query.Count(&count).Error
	return count > 0, err
}

// GetSubstituteTeacherConflicts checks if the substitute teacher has conflicts.
func (r *Repository) GetSubstituteTeacherConflicts(ctx context.Context, tenantID, substituteStaffID uuid.UUID, date time.Time, periodSlotIDs []uuid.UUID) ([]uuid.UUID, error) {
	var conflictingPeriods []uuid.UUID

	// Check for conflicts in timetable entries
	dayOfWeek := int(date.Weekday())
	err := r.db.WithContext(ctx).
		Model(&models.TimetableEntry{}).
		Joins("JOIN timetables ON timetables.id = timetable_entries.timetable_id").
		Where("timetable_entries.tenant_id = ?", tenantID).
		Where("timetable_entries.staff_id = ?", substituteStaffID).
		Where("timetable_entries.day_of_week = ?", dayOfWeek).
		Where("timetable_entries.period_slot_id IN ?", periodSlotIDs).
		Where("timetables.status = ?", models.TimetableStatusPublished).
		Where("timetables.deleted_at IS NULL").
		Pluck("timetable_entries.period_slot_id", &conflictingPeriods).Error

	if err != nil {
		return nil, err
	}

	// Also check for existing substitutions where they are already substitute
	var existingSubPeriods []uuid.UUID
	err = r.db.WithContext(ctx).
		Model(&models.SubstitutionPeriod{}).
		Joins("JOIN substitutions ON substitutions.id = substitution_periods.substitution_id").
		Where("substitutions.tenant_id = ?", tenantID).
		Where("substitutions.substitute_staff_id = ?", substituteStaffID).
		Where("substitutions.substitution_date = ?", date).
		Where("substitutions.status NOT IN ?", []string{"cancelled"}).
		Where("substitutions.deleted_at IS NULL").
		Where("substitution_periods.period_slot_id IN ?", periodSlotIDs).
		Pluck("substitution_periods.period_slot_id", &existingSubPeriods).Error

	if err != nil {
		return nil, err
	}

	conflictingPeriods = append(conflictingPeriods, existingSubPeriods...)
	return conflictingPeriods, nil
}

// GetAvailableTeachers returns teachers who are available for substitution.
func (r *Repository) GetAvailableTeachers(ctx context.Context, tenantID, branchID uuid.UUID, date time.Time, periodSlotIDs []uuid.UUID, excludeStaffID uuid.UUID) ([]models.Staff, error) {
	var staff []models.Staff

	// Get all teaching staff in the branch
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Where("branch_id = ?", branchID).
		Where("staff_type = ?", "teaching").
		Where("status = ?", "active").
		Where("id != ?", excludeStaffID).
		Where("deleted_at IS NULL").
		Preload("Department").
		Find(&staff).Error

	return staff, err
}

// GetTeacherTimetableEntries returns timetable entries for a teacher on a specific day.
func (r *Repository) GetTeacherTimetableEntries(ctx context.Context, tenantID, staffID uuid.UUID, dayOfWeek int) ([]models.TimetableEntry, error) {
	var entries []models.TimetableEntry

	err := r.db.WithContext(ctx).
		Joins("JOIN timetables ON timetables.id = timetable_entries.timetable_id").
		Where("timetable_entries.tenant_id = ?", tenantID).
		Where("timetable_entries.staff_id = ?", staffID).
		Where("timetable_entries.day_of_week = ?", dayOfWeek).
		Where("timetables.status = ?", models.TimetableStatusPublished).
		Where("timetables.deleted_at IS NULL").
		Preload("PeriodSlot").
		Preload("Subject").
		Preload("Timetable").
		Preload("Timetable.Section").
		Preload("Timetable.Section.Class").
		Find(&entries).Error

	return entries, err
}
