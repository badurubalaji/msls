// Package admission provides admission management services.
package admission

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
)

// MeritService handles merit list operations.
type MeritService struct {
	db *gorm.DB
}

// NewMeritService creates a new MeritService instance.
func NewMeritService(db *gorm.DB) *MeritService {
	return &MeritService{db: db}
}

// GenerateMeritListRequest represents a request to generate a merit list.
type GenerateMeritListRequest struct {
	TenantID    uuid.UUID
	SessionID   uuid.UUID
	ClassName   string
	TestID      *uuid.UUID
	CutoffScore *float64
	GeneratedBy *uuid.UUID
}

// GenerateMeritList generates a merit list for a session and class.
func (s *MeritService) GenerateMeritList(ctx context.Context, req GenerateMeritListRequest) (*models.MeritList, error) {
	if req.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if req.SessionID == uuid.Nil {
		return nil, ErrSessionIDRequired
	}
	if req.ClassName == "" {
		return nil, ErrClassNameRequired
	}

	// Verify session exists
	var session models.AdmissionSession
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", req.TenantID, req.SessionID).
		First(&session).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// Check if merit list already exists for this combination
	var existingList models.MeritList
	query := s.db.WithContext(ctx).
		Where("tenant_id = ? AND session_id = ? AND class_name = ?", req.TenantID, req.SessionID, req.ClassName)
	if req.TestID != nil {
		query = query.Where("test_id = ?", *req.TestID)
	} else {
		query = query.Where("test_id IS NULL")
	}

	err = query.First(&existingList).Error
	if err == nil {
		// List exists, check if it's finalized
		if existingList.IsFinal {
			return nil, ErrMeritListFinalized
		}
		// Delete the existing non-final list to regenerate
		if err := s.db.WithContext(ctx).Delete(&existingList).Error; err != nil {
			return nil, fmt.Errorf("failed to delete existing merit list: %w", err)
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing merit list: %w", err)
	}

	// Get applications for this session and class
	var applications []models.AdmissionApplication
	err = s.db.WithContext(ctx).
		Where("tenant_id = ? AND session_id = ? AND class_applying = ?", req.TenantID, req.SessionID, req.ClassName).
		Where("status IN ?", []string{
			string(models.ApplicationStatusSubmitted),
			string(models.ApplicationStatusUnderReview),
			string(models.ApplicationStatusApproved),
			string(models.ApplicationStatusWaitlisted),
		}).
		Find(&applications).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get applications: %w", err)
	}

	if len(applications) == 0 {
		return nil, ErrNoApplicantsForMeritList
	}

	// Build merit list entries
	entries := make(models.MeritListEntries, 0, len(applications))
	for _, app := range applications {
		entry := models.MeritListEntry{
			ApplicationID: app.ID,
			StudentName:   app.StudentName,
			Score:         0, // Will be calculated from test scores if available
			Status:        string(app.Status),
			ParentPhone:   app.FatherPhone,
		}

		// Use mother's phone if father's is not available
		if entry.ParentPhone == "" {
			entry.ParentPhone = app.MotherPhone
		}

		// Set parent email
		if app.FatherEmail != "" {
			entry.ParentEmail = &app.FatherEmail
		} else if app.MotherEmail != "" {
			entry.ParentEmail = &app.MotherEmail
		}

		// TODO: Get test scores from entrance_tests table when available
		// For now, we use a placeholder score calculation
		// In a real implementation, you would join with test_results table
		entry.Score = calculateMockScore(&app)

		entries = append(entries, entry)
	}

	// Sort entries by score (descending)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Score > entries[j].Score
	})

	// Assign ranks
	for i := range entries {
		entries[i].Rank = i + 1
	}

	// Apply cutoff filter if provided
	if req.CutoffScore != nil {
		filteredEntries := make(models.MeritListEntries, 0)
		for _, entry := range entries {
			if entry.Score >= *req.CutoffScore {
				filteredEntries = append(filteredEntries, entry)
			}
		}
		entries = filteredEntries
	}

	// Create merit list
	meritList := &models.MeritList{
		TenantID:    req.TenantID,
		SessionID:   req.SessionID,
		ClassName:   req.ClassName,
		TestID:      req.TestID,
		GeneratedAt: time.Now(),
		GeneratedBy: req.GeneratedBy,
		CutoffScore: req.CutoffScore,
		Entries:     entries,
		IsFinal:     false,
	}

	if err := s.db.WithContext(ctx).Create(meritList).Error; err != nil {
		return nil, fmt.Errorf("failed to create merit list: %w", err)
	}

	return meritList, nil
}

// calculateMockScore calculates a mock score for an application.
// This is a placeholder until test scores are integrated.
func calculateMockScore(app *models.AdmissionApplication) float64 {
	// In a real implementation, this would:
	// 1. Get test scores from entrance_tests table
	// 2. Calculate weighted average based on test weights
	// 3. Include interview scores if applicable
	// 4. Include previous academic marks if applicable

	// For now, return a deterministic score based on application data
	// This ensures consistent ordering
	score := float64(50) // Base score

	// Add some variation based on application details
	if app.Category != "" {
		score += 5
	}
	if app.PreviousSchool != "" {
		score += 10
	}

	return score
}

// GetMeritListRequest represents a request to get a merit list.
type GetMeritListRequest struct {
	TenantID  uuid.UUID
	SessionID uuid.UUID
	ClassName string
	TestID    *uuid.UUID
}

// GetMeritList retrieves the merit list for a session and class.
func (s *MeritService) GetMeritList(ctx context.Context, req GetMeritListRequest) (*models.MeritList, error) {
	if req.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if req.SessionID == uuid.Nil {
		return nil, ErrSessionIDRequired
	}

	var meritList models.MeritList
	query := s.db.WithContext(ctx).
		Where("tenant_id = ? AND session_id = ?", req.TenantID, req.SessionID)

	if req.ClassName != "" {
		query = query.Where("class_name = ?", req.ClassName)
	}

	if req.TestID != nil {
		query = query.Where("test_id = ?", *req.TestID)
	}

	err := query.Order("generated_at DESC").First(&meritList).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMeritListNotFound
		}
		return nil, fmt.Errorf("failed to get merit list: %w", err)
	}

	return &meritList, nil
}

// GetMeritListByID retrieves a merit list by ID.
func (s *MeritService) GetMeritListByID(ctx context.Context, tenantID, id uuid.UUID) (*models.MeritList, error) {
	if tenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}

	var meritList models.MeritList
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&meritList).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMeritListNotFound
		}
		return nil, fmt.Errorf("failed to get merit list: %w", err)
	}

	return &meritList, nil
}

// MeritListFilter contains filters for listing merit lists.
type MeritListFilter struct {
	TenantID  uuid.UUID
	SessionID uuid.UUID
	ClassName string
	IsFinal   *bool
}

// ListMeritLists retrieves merit lists with optional filtering.
func (s *MeritService) ListMeritLists(ctx context.Context, filter MeritListFilter) ([]models.MeritList, error) {
	if filter.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}

	query := s.db.WithContext(ctx).
		Model(&models.MeritList{}).
		Where("tenant_id = ?", filter.TenantID).
		Order("generated_at DESC")

	if filter.SessionID != uuid.Nil {
		query = query.Where("session_id = ?", filter.SessionID)
	}

	if filter.ClassName != "" {
		query = query.Where("class_name = ?", filter.ClassName)
	}

	if filter.IsFinal != nil {
		query = query.Where("is_final = ?", *filter.IsFinal)
	}

	var lists []models.MeritList
	if err := query.Find(&lists).Error; err != nil {
		return nil, fmt.Errorf("failed to list merit lists: %w", err)
	}

	return lists, nil
}

// FinalizeMeritList marks a merit list as final.
func (s *MeritService) FinalizeMeritList(ctx context.Context, tenantID, id uuid.UUID, updatedBy *uuid.UUID) (*models.MeritList, error) {
	meritList, err := s.GetMeritListByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if meritList.IsFinal {
		return nil, ErrMeritListFinalized
	}

	if err := s.db.WithContext(ctx).Model(meritList).Update("is_final", true).Error; err != nil {
		return nil, fmt.Errorf("failed to finalize merit list: %w", err)
	}

	meritList.IsFinal = true
	return meritList, nil
}

// DeleteMeritList deletes a merit list.
func (s *MeritService) DeleteMeritList(ctx context.Context, tenantID, id uuid.UUID) error {
	meritList, err := s.GetMeritListByID(ctx, tenantID, id)
	if err != nil {
		return err
	}

	if meritList.IsFinal {
		return ErrMeritListFinalized
	}

	if err := s.db.WithContext(ctx).Delete(meritList).Error; err != nil {
		return fmt.Errorf("failed to delete merit list: %w", err)
	}

	return nil
}

// UpdateMeritListCutoff updates the cutoff score of a merit list.
func (s *MeritService) UpdateMeritListCutoff(ctx context.Context, tenantID, id uuid.UUID, cutoffScore *float64) (*models.MeritList, error) {
	meritList, err := s.GetMeritListByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if meritList.IsFinal {
		return nil, ErrMeritListFinalized
	}

	if err := s.db.WithContext(ctx).Model(meritList).Update("cutoff_score", cutoffScore).Error; err != nil {
		return nil, fmt.Errorf("failed to update cutoff score: %w", err)
	}

	meritList.CutoffScore = cutoffScore
	return meritList, nil
}
