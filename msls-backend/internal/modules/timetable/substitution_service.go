package timetable

import (
	"context"
	"time"

	"msls-backend/internal/pkg/database/models"

	"github.com/google/uuid"
)

// ========================================
// Substitution Service Methods
// ========================================

// ListSubstitutions returns substitutions based on filter criteria.
func (s *Service) ListSubstitutions(ctx context.Context, filter SubstitutionFilter) ([]models.Substitution, int64, error) {
	return s.repo.ListSubstitutions(ctx, filter)
}

// GetSubstitutionByID returns a substitution by ID.
func (s *Service) GetSubstitutionByID(ctx context.Context, tenantID, id uuid.UUID) (*models.Substitution, error) {
	return s.repo.GetSubstitutionByID(ctx, tenantID, id)
}

// CreateSubstitution creates a new substitution.
func (s *Service) CreateSubstitution(ctx context.Context, tenantID uuid.UUID, req CreateSubstitutionRequest, userID uuid.UUID) (*models.Substitution, error) {
	// Parse date
	date, err := time.Parse("2006-01-02", req.SubstitutionDate)
	if err != nil {
		return nil, err
	}

	// Extract period slot IDs
	periodSlotIDs := make([]uuid.UUID, len(req.Periods))
	for i, p := range req.Periods {
		periodSlotIDs[i] = p.PeriodSlotID
	}

	// Check for existing substitution conflict
	hasConflict, err := s.repo.CheckSubstitutionConflict(ctx, tenantID, req.OriginalStaffID, date, periodSlotIDs, nil)
	if err != nil {
		return nil, err
	}
	if hasConflict {
		return nil, ErrSubstitutionConflict
	}

	// Check if substitute teacher has conflicts
	conflictingPeriods, err := s.repo.GetSubstituteTeacherConflicts(ctx, tenantID, req.SubstituteStaffID, date, periodSlotIDs)
	if err != nil {
		return nil, err
	}
	if len(conflictingPeriods) > 0 {
		return nil, ErrSubstituteConflict
	}

	// Create substitution
	substitution := &models.Substitution{
		TenantID:          tenantID,
		BranchID:          req.BranchID,
		OriginalStaffID:   req.OriginalStaffID,
		SubstituteStaffID: req.SubstituteStaffID,
		SubstitutionDate:  date,
		Reason:            req.Reason,
		Status:            models.SubstitutionStatusPending,
		Notes:             req.Notes,
		CreatedBy:         &userID,
	}

	if err := s.repo.CreateSubstitution(ctx, substitution); err != nil {
		return nil, err
	}

	// Create substitution periods
	periods := make([]models.SubstitutionPeriod, len(req.Periods))
	for i, p := range req.Periods {
		periods[i] = models.SubstitutionPeriod{
			SubstitutionID:   substitution.ID,
			PeriodSlotID:     p.PeriodSlotID,
			TimetableEntryID: p.TimetableEntryID,
			SubjectID:        p.SubjectID,
			SectionID:        p.SectionID,
			RoomNumber:       p.RoomNumber,
			Notes:            p.Notes,
		}
	}

	if err := s.repo.CreateSubstitutionPeriods(ctx, periods); err != nil {
		return nil, err
	}

	return s.repo.GetSubstitutionByID(ctx, tenantID, substitution.ID)
}

// UpdateSubstitution updates a substitution.
func (s *Service) UpdateSubstitution(ctx context.Context, tenantID, id uuid.UUID, req UpdateSubstitutionRequest) (*models.Substitution, error) {
	substitution, err := s.repo.GetSubstitutionByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	// Only pending substitutions can be fully modified
	if req.SubstituteStaffID != nil && substitution.Status != models.SubstitutionStatusPending {
		return nil, ErrSubstitutionNotPending
	}

	// Update fields
	if req.SubstituteStaffID != nil {
		// Check for conflicts with new substitute
		periodSlotIDs := make([]uuid.UUID, len(substitution.Periods))
		for i, p := range substitution.Periods {
			periodSlotIDs[i] = p.PeriodSlotID
		}

		conflictingPeriods, err := s.repo.GetSubstituteTeacherConflicts(ctx, tenantID, *req.SubstituteStaffID, substitution.SubstitutionDate, periodSlotIDs)
		if err != nil {
			return nil, err
		}
		if len(conflictingPeriods) > 0 {
			return nil, ErrSubstituteConflict
		}

		substitution.SubstituteStaffID = *req.SubstituteStaffID
	}

	if req.Reason != nil {
		substitution.Reason = *req.Reason
	}

	if req.Notes != nil {
		substitution.Notes = *req.Notes
	}

	if req.Status != nil {
		substitution.Status = models.SubstitutionStatus(*req.Status)
	}

	if err := s.repo.UpdateSubstitution(ctx, substitution); err != nil {
		return nil, err
	}

	return s.repo.GetSubstitutionByID(ctx, tenantID, id)
}

// ConfirmSubstitution confirms a pending substitution.
func (s *Service) ConfirmSubstitution(ctx context.Context, tenantID, id uuid.UUID, userID uuid.UUID) (*models.Substitution, error) {
	substitution, err := s.repo.GetSubstitutionByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if substitution.Status != models.SubstitutionStatusPending {
		return nil, ErrSubstitutionNotPending
	}

	now := time.Now()
	substitution.Status = models.SubstitutionStatusConfirmed
	substitution.ApprovedBy = &userID
	substitution.ApprovedAt = &now

	if err := s.repo.UpdateSubstitution(ctx, substitution); err != nil {
		return nil, err
	}

	return s.repo.GetSubstitutionByID(ctx, tenantID, id)
}

// CancelSubstitution cancels a substitution.
func (s *Service) CancelSubstitution(ctx context.Context, tenantID, id uuid.UUID) (*models.Substitution, error) {
	substitution, err := s.repo.GetSubstitutionByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if substitution.Status != models.SubstitutionStatusPending && substitution.Status != models.SubstitutionStatusConfirmed {
		return nil, ErrSubstitutionNotCancellable
	}

	substitution.Status = models.SubstitutionStatusCancelled

	if err := s.repo.UpdateSubstitution(ctx, substitution); err != nil {
		return nil, err
	}

	return s.repo.GetSubstitutionByID(ctx, tenantID, id)
}

// DeleteSubstitution deletes a substitution.
func (s *Service) DeleteSubstitution(ctx context.Context, tenantID, id uuid.UUID) error {
	substitution, err := s.repo.GetSubstitutionByID(ctx, tenantID, id)
	if err != nil {
		return err
	}

	// Only pending substitutions can be deleted
	if substitution.Status != models.SubstitutionStatusPending {
		return ErrSubstitutionNotPending
	}

	// Delete periods first
	if err := s.repo.DeleteSubstitutionPeriods(ctx, id); err != nil {
		return err
	}

	return s.repo.DeleteSubstitution(ctx, tenantID, id)
}

// GetAvailableTeachers returns teachers available for substitution.
func (s *Service) GetAvailableTeachers(ctx context.Context, tenantID, branchID uuid.UUID, dateStr string, periodSlotIDs []uuid.UUID, excludeStaffID uuid.UUID) ([]AvailableTeacherResponse, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}

	dayOfWeek := int(date.Weekday())

	// Get all teaching staff
	staff, err := s.repo.GetAvailableTeachers(ctx, tenantID, branchID, date, periodSlotIDs, excludeStaffID)
	if err != nil {
		return nil, err
	}

	result := make([]AvailableTeacherResponse, 0, len(staff))

	for _, st := range staff {
		// Get their timetable for this day
		entries, err := s.repo.GetTeacherTimetableEntries(ctx, tenantID, st.ID, dayOfWeek)
		if err != nil {
			continue
		}

		// Check for conflicts with requested periods
		conflictingPeriods, err := s.repo.GetSubstituteTeacherConflicts(ctx, tenantID, st.ID, date, periodSlotIDs)
		if err != nil {
			continue
		}

		hasConflict := len(conflictingPeriods) > 0

		// Calculate free periods (assuming 8 total periods)
		teacherResp := AvailableTeacherResponse{
			StaffID:      st.ID,
			StaffName:    st.FirstName,
			TotalPeriods: len(entries),
			FreePeriods:  8 - len(entries), // Simplified calculation
			HasConflict:  hasConflict,
		}

		if st.LastName != "" {
			teacherResp.StaffName += " " + st.LastName
		}

		if st.Department != nil {
			teacherResp.DepartmentID = st.DepartmentID
			teacherResp.Department = st.Department.Name
		}

		result = append(result, teacherResp)
	}

	return result, nil
}

// GetTeacherAbsencePeriods returns the timetable entries for an absent teacher on a specific date.
func (s *Service) GetTeacherAbsencePeriods(ctx context.Context, tenantID, staffID uuid.UUID, dateStr string) ([]TimetableEntryResponse, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}

	dayOfWeek := int(date.Weekday())

	entries, err := s.repo.GetTeacherTimetableEntries(ctx, tenantID, staffID, dayOfWeek)
	if err != nil {
		return nil, err
	}

	result := make([]TimetableEntryResponse, len(entries))
	for i, e := range entries {
		result[i] = TimetableEntryToResponse(&e)
	}

	return result, nil
}
