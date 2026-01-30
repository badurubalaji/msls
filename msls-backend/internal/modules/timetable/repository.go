// Package timetable provides timetable management functionality.
package timetable

import (
	"context"
	"errors"

	"msls-backend/internal/pkg/database/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository handles database operations for timetable entities.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new timetable repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// ========================================
// Shift Repository Methods
// ========================================

// ListShifts returns all shifts for a tenant with filters.
func (r *Repository) ListShifts(ctx context.Context, filter ShiftFilter) ([]models.Shift, int64, error) {
	var shifts []models.Shift
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Shift{}).
		Where("tenant_id = ?", filter.TenantID)

	if filter.BranchID != nil {
		query = query.Where("branch_id = ?", *filter.BranchID)
	}

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Preload("Branch").
		Order("display_order ASC, name ASC").
		Find(&shifts).Error; err != nil {
		return nil, 0, err
	}

	return shifts, total, nil
}

// GetShiftByID returns a shift by ID.
func (r *Repository) GetShiftByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Shift, error) {
	var shift models.Shift
	err := r.db.WithContext(ctx).
		Preload("Branch").
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&shift).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrShiftNotFound
	}
	return &shift, err
}

// GetShiftByCode returns a shift by code.
func (r *Repository) GetShiftByCode(ctx context.Context, tenantID, branchID uuid.UUID, code string) (*models.Shift, error) {
	var shift models.Shift
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND branch_id = ? AND code = ?", tenantID, branchID, code).
		First(&shift).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &shift, err
}

// CreateShift creates a new shift.
func (r *Repository) CreateShift(ctx context.Context, shift *models.Shift) error {
	return r.db.WithContext(ctx).Create(shift).Error
}

// UpdateShift updates an existing shift.
func (r *Repository) UpdateShift(ctx context.Context, shift *models.Shift) error {
	return r.db.WithContext(ctx).Save(shift).Error
}

// DeleteShift deletes a shift.
func (r *Repository) DeleteShift(ctx context.Context, tenantID, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&models.Shift{}).Error
}

// ========================================
// Day Pattern Repository Methods
// ========================================

// ListDayPatterns returns all day patterns for a tenant with filters.
func (r *Repository) ListDayPatterns(ctx context.Context, filter DayPatternFilter) ([]models.DayPattern, int64, error) {
	var patterns []models.DayPattern
	var total int64

	query := r.db.WithContext(ctx).Model(&models.DayPattern{}).
		Where("tenant_id = ?", filter.TenantID)

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("display_order ASC, name ASC").
		Find(&patterns).Error; err != nil {
		return nil, 0, err
	}

	return patterns, total, nil
}

// GetDayPatternByID returns a day pattern by ID.
func (r *Repository) GetDayPatternByID(ctx context.Context, tenantID, id uuid.UUID) (*models.DayPattern, error) {
	var pattern models.DayPattern
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&pattern).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrDayPatternNotFound
	}
	return &pattern, err
}

// GetDayPatternByCode returns a day pattern by code.
func (r *Repository) GetDayPatternByCode(ctx context.Context, tenantID uuid.UUID, code string) (*models.DayPattern, error) {
	var pattern models.DayPattern
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND code = ?", tenantID, code).
		First(&pattern).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &pattern, err
}

// CreateDayPattern creates a new day pattern.
func (r *Repository) CreateDayPattern(ctx context.Context, pattern *models.DayPattern) error {
	return r.db.WithContext(ctx).Create(pattern).Error
}

// UpdateDayPattern updates an existing day pattern.
func (r *Repository) UpdateDayPattern(ctx context.Context, pattern *models.DayPattern) error {
	return r.db.WithContext(ctx).Save(pattern).Error
}

// DeleteDayPattern deletes a day pattern.
func (r *Repository) DeleteDayPattern(ctx context.Context, tenantID, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&models.DayPattern{}).Error
}

// IsDayPatternInUse checks if a day pattern is used.
func (r *Repository) IsDayPatternInUse(ctx context.Context, id uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.DayPatternAssignment{}).
		Where("day_pattern_id = ?", id).
		Count(&count).Error; err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}

	if err := r.db.WithContext(ctx).Model(&models.PeriodSlot{}).
		Where("day_pattern_id = ?", id).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ========================================
// Day Pattern Assignment Repository Methods
// ========================================

// ListDayPatternAssignments returns all day pattern assignments for a branch.
func (r *Repository) ListDayPatternAssignments(ctx context.Context, tenantID, branchID uuid.UUID) ([]models.DayPatternAssignment, error) {
	var assignments []models.DayPatternAssignment
	err := r.db.WithContext(ctx).
		Preload("DayPattern").
		Where("tenant_id = ? AND branch_id = ?", tenantID, branchID).
		Order("day_of_week ASC").
		Find(&assignments).Error
	return assignments, err
}

// GetDayPatternAssignment returns a day pattern assignment by ID.
func (r *Repository) GetDayPatternAssignment(ctx context.Context, tenantID, branchID uuid.UUID, dayOfWeek int) (*models.DayPatternAssignment, error) {
	var assignment models.DayPatternAssignment
	err := r.db.WithContext(ctx).
		Preload("DayPattern").
		Where("tenant_id = ? AND branch_id = ? AND day_of_week = ?", tenantID, branchID, dayOfWeek).
		First(&assignment).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &assignment, err
}

// UpsertDayPatternAssignment creates or updates a day pattern assignment.
func (r *Repository) UpsertDayPatternAssignment(ctx context.Context, assignment *models.DayPatternAssignment) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND branch_id = ? AND day_of_week = ?",
			assignment.TenantID, assignment.BranchID, assignment.DayOfWeek).
		Assign(map[string]interface{}{
			"day_pattern_id": assignment.DayPatternID,
			"is_working_day": assignment.IsWorkingDay,
		}).
		FirstOrCreate(assignment).Error
}

// ========================================
// Period Slot Repository Methods
// ========================================

// ListPeriodSlots returns all period slots for a tenant with filters.
func (r *Repository) ListPeriodSlots(ctx context.Context, filter PeriodSlotFilter) ([]models.PeriodSlot, int64, error) {
	var slots []models.PeriodSlot
	var total int64

	query := r.db.WithContext(ctx).Model(&models.PeriodSlot{}).
		Where("tenant_id = ?", filter.TenantID)

	if filter.BranchID != nil {
		query = query.Where("branch_id = ?", *filter.BranchID)
	}

	if filter.DayPatternID != nil {
		query = query.Where("day_pattern_id = ?", *filter.DayPatternID)
	}

	if filter.ShiftID != nil {
		query = query.Where("shift_id = ?", *filter.ShiftID)
	}

	if filter.SlotType != nil {
		query = query.Where("slot_type = ?", *filter.SlotType)
	}

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Preload("DayPattern").
		Preload("Shift").
		Order("display_order ASC, start_time ASC").
		Find(&slots).Error; err != nil {
		return nil, 0, err
	}

	return slots, total, nil
}

// GetPeriodSlotByID returns a period slot by ID.
func (r *Repository) GetPeriodSlotByID(ctx context.Context, tenantID, id uuid.UUID) (*models.PeriodSlot, error) {
	var slot models.PeriodSlot
	err := r.db.WithContext(ctx).
		Preload("DayPattern").
		Preload("Shift").
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&slot).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrPeriodSlotNotFound
	}
	return &slot, err
}

// CreatePeriodSlot creates a new period slot.
func (r *Repository) CreatePeriodSlot(ctx context.Context, slot *models.PeriodSlot) error {
	return r.db.WithContext(ctx).Create(slot).Error
}

// UpdatePeriodSlot updates an existing period slot.
func (r *Repository) UpdatePeriodSlot(ctx context.Context, slot *models.PeriodSlot) error {
	return r.db.WithContext(ctx).Save(slot).Error
}

// DeletePeriodSlot deletes a period slot.
func (r *Repository) DeletePeriodSlot(ctx context.Context, tenantID, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&models.PeriodSlot{}).Error
}

// ========================================
// Timetable Repository Methods
// ========================================

// ListTimetables returns all timetables for a tenant with filters.
func (r *Repository) ListTimetables(ctx context.Context, filter TimetableFilter) ([]models.Timetable, int64, error) {
	var timetables []models.Timetable
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Timetable{}).
		Where("tenant_id = ? AND deleted_at IS NULL", filter.TenantID)

	if filter.BranchID != nil {
		query = query.Where("branch_id = ?", *filter.BranchID)
	}

	if filter.SectionID != nil {
		query = query.Where("section_id = ?", *filter.SectionID)
	}

	if filter.AcademicYearID != nil {
		query = query.Where("academic_year_id = ?", *filter.AcademicYearID)
	}

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Preload("Branch").
		Preload("Section").
		Preload("Section.Class").
		Preload("AcademicYear").
		Order("created_at DESC").
		Find(&timetables).Error; err != nil {
		return nil, 0, err
	}

	return timetables, total, nil
}

// GetTimetableByID returns a timetable by ID with all entries.
func (r *Repository) GetTimetableByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Timetable, error) {
	var timetable models.Timetable
	err := r.db.WithContext(ctx).
		Preload("Branch").
		Preload("Section").
		Preload("Section.Class").
		Preload("AcademicYear").
		Preload("Entries", func(db *gorm.DB) *gorm.DB {
			return db.Order("day_of_week ASC, period_slot_id ASC")
		}).
		Preload("Entries.PeriodSlot").
		Preload("Entries.Subject").
		Preload("Entries.Staff").
		Where("tenant_id = ? AND id = ? AND deleted_at IS NULL", tenantID, id).
		First(&timetable).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrTimetableNotFound
	}
	return &timetable, err
}

// GetPublishedTimetableForSection returns the published timetable for a section.
func (r *Repository) GetPublishedTimetableForSection(ctx context.Context, tenantID, sectionID, academicYearID uuid.UUID) (*models.Timetable, error) {
	var timetable models.Timetable
	err := r.db.WithContext(ctx).
		Preload("Branch").
		Preload("Section").
		Preload("Section.Class").
		Preload("AcademicYear").
		Preload("Entries", func(db *gorm.DB) *gorm.DB {
			return db.Order("day_of_week ASC, period_slot_id ASC")
		}).
		Preload("Entries.PeriodSlot").
		Preload("Entries.Subject").
		Preload("Entries.Staff").
		Where("tenant_id = ? AND section_id = ? AND academic_year_id = ? AND status = ? AND deleted_at IS NULL",
			tenantID, sectionID, academicYearID, models.TimetableStatusPublished).
		First(&timetable).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &timetable, err
}

// CreateTimetable creates a new timetable.
func (r *Repository) CreateTimetable(ctx context.Context, timetable *models.Timetable) error {
	return r.db.WithContext(ctx).Create(timetable).Error
}

// UpdateTimetable updates an existing timetable.
func (r *Repository) UpdateTimetable(ctx context.Context, timetable *models.Timetable) error {
	return r.db.WithContext(ctx).Save(timetable).Error
}

// DeleteTimetable soft deletes a timetable.
func (r *Repository) DeleteTimetable(ctx context.Context, tenantID, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Timetable{}).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

// ArchiveOtherTimetables archives all other published timetables for the same section.
func (r *Repository) ArchiveOtherTimetables(ctx context.Context, tenantID, sectionID, academicYearID, excludeID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Timetable{}).
		Where("tenant_id = ? AND section_id = ? AND academic_year_id = ? AND id != ? AND status = ? AND deleted_at IS NULL",
			tenantID, sectionID, academicYearID, excludeID, models.TimetableStatusPublished).
		Update("status", models.TimetableStatusArchived).Error
}

// ========================================
// Timetable Entry Repository Methods
// ========================================

// GetTimetableEntries returns all entries for a timetable.
func (r *Repository) GetTimetableEntries(ctx context.Context, timetableID uuid.UUID) ([]models.TimetableEntry, error) {
	var entries []models.TimetableEntry
	err := r.db.WithContext(ctx).
		Preload("PeriodSlot").
		Preload("Subject").
		Preload("Staff").
		Where("timetable_id = ?", timetableID).
		Order("day_of_week ASC, period_slot_id ASC").
		Find(&entries).Error
	return entries, err
}

// GetTimetableEntry returns a single entry by ID.
func (r *Repository) GetTimetableEntry(ctx context.Context, tenantID, entryID uuid.UUID) (*models.TimetableEntry, error) {
	var entry models.TimetableEntry
	err := r.db.WithContext(ctx).
		Preload("PeriodSlot").
		Preload("Subject").
		Preload("Staff").
		Where("tenant_id = ? AND id = ?", tenantID, entryID).
		First(&entry).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrTimetableEntryNotFound
	}
	return &entry, err
}

// UpsertTimetableEntry creates or updates a timetable entry.
func (r *Repository) UpsertTimetableEntry(ctx context.Context, entry *models.TimetableEntry) error {
	// Try to find existing entry for this slot
	var existing models.TimetableEntry
	err := r.db.WithContext(ctx).
		Where("timetable_id = ? AND day_of_week = ? AND period_slot_id = ?",
			entry.TimetableID, entry.DayOfWeek, entry.PeriodSlotID).
		First(&existing).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new entry
		return r.db.WithContext(ctx).Create(entry).Error
	} else if err != nil {
		return err
	}

	// Update existing entry
	entry.ID = existing.ID
	return r.db.WithContext(ctx).Save(entry).Error
}

// DeleteTimetableEntry deletes a timetable entry.
func (r *Repository) DeleteTimetableEntry(ctx context.Context, tenantID, entryID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, entryID).
		Delete(&models.TimetableEntry{}).Error
}

// DeleteTimetableEntriesByTimetableID deletes all entries for a timetable.
func (r *Repository) DeleteTimetableEntriesByTimetableID(ctx context.Context, timetableID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("timetable_id = ?", timetableID).
		Delete(&models.TimetableEntry{}).Error
}

// ========================================
// Conflict Detection Repository Methods
// ========================================

// GetTeacherConflicts checks for teacher double-booking conflicts.
func (r *Repository) GetTeacherConflicts(ctx context.Context, tenantID, staffID uuid.UUID, dayOfWeek int, periodSlotID uuid.UUID, excludeTimetableID *uuid.UUID) ([]models.TimetableEntry, error) {
	var entries []models.TimetableEntry

	query := r.db.WithContext(ctx).
		Joins("JOIN timetables ON timetables.id = timetable_entries.timetable_id").
		Preload("Timetable").
		Preload("Timetable.Section").
		Preload("Timetable.Section.Class").
		Preload("PeriodSlot").
		Preload("Subject").
		Where("timetable_entries.tenant_id = ? AND timetable_entries.staff_id = ? AND timetable_entries.day_of_week = ? AND timetable_entries.period_slot_id = ?",
			tenantID, staffID, dayOfWeek, periodSlotID).
		Where("timetables.status = ? AND timetables.deleted_at IS NULL", models.TimetableStatusPublished)

	if excludeTimetableID != nil {
		query = query.Where("timetables.id != ?", *excludeTimetableID)
	}

	err := query.Find(&entries).Error
	return entries, err
}

// GetTeacherSchedule returns all timetable entries for a teacher across published timetables.
func (r *Repository) GetTeacherSchedule(ctx context.Context, tenantID, staffID, academicYearID uuid.UUID) ([]models.TimetableEntry, error) {
	var entries []models.TimetableEntry

	err := r.db.WithContext(ctx).
		Joins("JOIN timetables ON timetables.id = timetable_entries.timetable_id").
		Preload("Timetable").
		Preload("Timetable.Section").
		Preload("Timetable.Section.Class").
		Preload("PeriodSlot").
		Preload("Subject").
		Where("timetable_entries.tenant_id = ? AND timetable_entries.staff_id = ?", tenantID, staffID).
		Where("timetables.academic_year_id = ? AND timetables.status = ? AND timetables.deleted_at IS NULL",
			academicYearID, models.TimetableStatusPublished).
		Order("timetable_entries.day_of_week ASC, timetable_entries.period_slot_id ASC").
		Find(&entries).Error

	return entries, err
}

// GetStaffIDByUserID returns the staff ID for a given user ID.
func (r *Repository) GetStaffIDByUserID(ctx context.Context, tenantID, userID uuid.UUID) (uuid.UUID, error) {
	var staff models.Staff
	err := r.db.WithContext(ctx).
		Select("id").
		Where("tenant_id = ? AND user_id = ? AND deleted_at IS NULL", tenantID, userID).
		First(&staff).Error

	if err != nil {
		return uuid.Nil, err
	}

	return staff.ID, nil
}
