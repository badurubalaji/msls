// Package admission provides admission management services.
package admission

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
)

// DecisionService handles admission decision operations.
type DecisionService struct {
	db *gorm.DB
}

// NewDecisionService creates a new DecisionService instance.
func NewDecisionService(db *gorm.DB) *DecisionService {
	return &DecisionService{db: db}
}

// CreateDecisionRequest represents a request to create an admission decision.
type CreateDecisionRequest struct {
	TenantID         uuid.UUID
	ApplicationID    uuid.UUID
	Decision         models.DecisionType
	DecisionDate     time.Time
	DecidedBy        *uuid.UUID
	SectionAssigned  *string
	WaitlistPosition *int
	RejectionReason  *string
	OfferValidUntil  *time.Time
	Remarks          *string
}

// CreateDecision creates an admission decision for an application.
func (s *DecisionService) CreateDecision(ctx context.Context, req CreateDecisionRequest) (*models.AdmissionDecision, error) {
	if req.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if req.ApplicationID == uuid.Nil {
		return nil, ErrApplicationIDRequired
	}
	if req.Decision == "" {
		return nil, ErrDecisionRequired
	}
	if !req.Decision.IsValid() {
		return nil, ErrInvalidDecisionType
	}
	if req.DecisionDate.IsZero() {
		return nil, ErrDecisionDateRequired
	}

	// Validate decision-specific requirements
	if req.Decision == models.DecisionWaitlisted && (req.WaitlistPosition == nil || *req.WaitlistPosition <= 0) {
		return nil, ErrWaitlistPositionRequired
	}
	if req.Decision == models.DecisionRejected && (req.RejectionReason == nil || *req.RejectionReason == "") {
		return nil, ErrRejectionReasonRequired
	}

	// Verify application exists
	var application models.AdmissionApplication
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", req.TenantID, req.ApplicationID).
		First(&application).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrApplicationNotFound
		}
		return nil, fmt.Errorf("failed to get application: %w", err)
	}

	// Check if a primary decision already exists for this application
	var existingDecision models.AdmissionDecision
	err = s.db.WithContext(ctx).
		Where("tenant_id = ? AND application_id = ?", req.TenantID, req.ApplicationID).
		First(&existingDecision).Error
	if err == nil {
		return nil, ErrDecisionExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing decision: %w", err)
	}

	// Create decision in a transaction
	var decision *models.AdmissionDecision
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		decision = &models.AdmissionDecision{
			TenantID:         req.TenantID,
			ApplicationID:    req.ApplicationID,
			Decision:         req.Decision,
			DecisionDate:     req.DecisionDate,
			DecidedBy:        req.DecidedBy,
			SectionAssigned:  req.SectionAssigned,
			WaitlistPosition: req.WaitlistPosition,
			RejectionReason:  req.RejectionReason,
			OfferValidUntil:  req.OfferValidUntil,
			Remarks:          req.Remarks,
			CreatedBy:        req.DecidedBy,
			UpdatedBy:        req.DecidedBy,
		}

		if err := tx.Create(decision).Error; err != nil {
			return fmt.Errorf("failed to create decision: %w", err)
		}

		// Update application status based on decision
		var newStatus models.ApplicationStatus
		switch req.Decision {
		case models.DecisionApproved:
			newStatus = models.ApplicationStatusApproved
		case models.DecisionWaitlisted:
			newStatus = models.ApplicationStatusWaitlisted
		case models.DecisionRejected:
			newStatus = models.ApplicationStatusRejected
		}

		updates := map[string]interface{}{
			"status":     newStatus,
			"updated_by": req.DecidedBy,
			"updated_at": time.Now(),
		}

		if req.Decision == models.DecisionApproved {
			now := time.Now()
			updates["approved_at"] = &now
			updates["approved_by"] = req.DecidedBy
		}

		if req.Decision == models.DecisionWaitlisted {
			updates["waitlist_position"] = req.WaitlistPosition
		}

		if err := tx.Model(&application).Updates(updates).Error; err != nil {
			return fmt.Errorf("failed to update application status: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return decision, nil
}

// GetDecisionByID retrieves a decision by ID.
func (s *DecisionService) GetDecisionByID(ctx context.Context, tenantID, id uuid.UUID) (*models.AdmissionDecision, error) {
	if tenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}

	var decision models.AdmissionDecision
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&decision).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDecisionNotFound
		}
		return nil, fmt.Errorf("failed to get decision: %w", err)
	}

	return &decision, nil
}

// GetDecisionByApplication retrieves the decision for an application.
func (s *DecisionService) GetDecisionByApplication(ctx context.Context, tenantID, applicationID uuid.UUID) (*models.AdmissionDecision, error) {
	if tenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if applicationID == uuid.Nil {
		return nil, ErrApplicationIDRequired
	}

	var decision models.AdmissionDecision
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND application_id = ?", tenantID, applicationID).
		First(&decision).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDecisionNotFound
		}
		return nil, fmt.Errorf("failed to get decision: %w", err)
	}

	return &decision, nil
}

// GenerateOfferLetterRequest represents a request to generate an offer letter.
type GenerateOfferLetterRequest struct {
	TenantID      uuid.UUID
	ApplicationID uuid.UUID
	ValidUntil    *time.Time
	GeneratedBy   *uuid.UUID
}

// GenerateOfferLetter generates an offer letter for an approved application.
func (s *DecisionService) GenerateOfferLetter(ctx context.Context, req GenerateOfferLetterRequest) (*models.AdmissionDecision, error) {
	if req.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if req.ApplicationID == uuid.Nil {
		return nil, ErrApplicationIDRequired
	}

	// Get the decision
	decision, err := s.GetDecisionByApplication(ctx, req.TenantID, req.ApplicationID)
	if err != nil {
		return nil, err
	}

	// Verify decision is approved
	if decision.Decision != models.DecisionApproved {
		return nil, ErrOfferNotFound
	}

	// Set default validity if not provided (30 days)
	validUntil := req.ValidUntil
	if validUntil == nil {
		defaultValidity := time.Now().AddDate(0, 0, 30)
		validUntil = &defaultValidity
	}

	// TODO: Implement actual PDF generation
	// For now, we just set a placeholder URL
	offerLetterURL := fmt.Sprintf("/api/v1/applications/%s/offer-letter.pdf", req.ApplicationID)

	updates := map[string]interface{}{
		"offer_letter_url":  offerLetterURL,
		"offer_valid_until": validUntil,
		"updated_by":        req.GeneratedBy,
		"updated_at":        time.Now(),
	}

	if err := s.db.WithContext(ctx).Model(decision).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update decision with offer letter: %w", err)
	}

	decision.OfferLetterURL = &offerLetterURL
	decision.OfferValidUntil = validUntil

	return decision, nil
}

// AcceptOfferRequest represents a request to accept an offer.
type AcceptOfferRequest struct {
	TenantID      uuid.UUID
	ApplicationID uuid.UUID
	AcceptedBy    *uuid.UUID
}

// AcceptOffer accepts an offer for an application.
func (s *DecisionService) AcceptOffer(ctx context.Context, req AcceptOfferRequest) (*models.AdmissionDecision, error) {
	if req.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if req.ApplicationID == uuid.Nil {
		return nil, ErrApplicationIDRequired
	}

	// Get the decision
	decision, err := s.GetDecisionByApplication(ctx, req.TenantID, req.ApplicationID)
	if err != nil {
		return nil, err
	}

	// Verify decision is approved
	if decision.Decision != models.DecisionApproved {
		return nil, ErrOfferNotFound
	}

	// Check if offer is already accepted
	if decision.OfferAccepted != nil && *decision.OfferAccepted {
		return nil, ErrOfferAlreadyAccepted
	}

	// Check if offer has expired
	if decision.OfferValidUntil != nil && decision.OfferValidUntil.Before(time.Now()) {
		return nil, ErrOfferExpired
	}

	accepted := true
	now := time.Now()

	updates := map[string]interface{}{
		"offer_accepted":    accepted,
		"offer_accepted_at": now,
		"updated_by":        req.AcceptedBy,
		"updated_at":        now,
	}

	if err := s.db.WithContext(ctx).Model(decision).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to accept offer: %w", err)
	}

	decision.OfferAccepted = &accepted
	decision.OfferAcceptedAt = &now

	return decision, nil
}

// EnrollRequest represents a request to complete enrollment.
type EnrollRequest struct {
	TenantID      uuid.UUID
	ApplicationID uuid.UUID
	EnrolledBy    *uuid.UUID
}

// Enroll completes the enrollment process for an application.
func (s *DecisionService) Enroll(ctx context.Context, req EnrollRequest) (*models.AdmissionApplication, error) {
	if req.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if req.ApplicationID == uuid.Nil {
		return nil, ErrApplicationIDRequired
	}

	// Get the application
	var application models.AdmissionApplication
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", req.TenantID, req.ApplicationID).
		First(&application).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrApplicationNotFound
		}
		return nil, fmt.Errorf("failed to get application: %w", err)
	}

	// Check if already enrolled
	if application.Status == models.ApplicationStatusEnrolled {
		return nil, ErrAlreadyEnrolled
	}

	// Verify application is approved
	if application.Status != models.ApplicationStatusApproved {
		return nil, ErrInvalidApplicationStatus
	}

	// Get the decision and verify offer is accepted
	decision, err := s.GetDecisionByApplication(ctx, req.TenantID, req.ApplicationID)
	if err != nil {
		return nil, err
	}

	if decision.OfferAccepted == nil || !*decision.OfferAccepted {
		return nil, ErrOfferNotAccepted
	}

	// Complete enrollment in a transaction
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		// Update application status to enrolled
		updates := map[string]interface{}{
			"status":      models.ApplicationStatusEnrolled,
			"enrolled_at": &now,
			"updated_by":  req.EnrolledBy,
			"updated_at":  now,
		}

		if err := tx.Model(&application).Updates(updates).Error; err != nil {
			return fmt.Errorf("failed to update application status: %w", err)
		}

		// Update filled seats count
		if err := tx.Exec(`
			UPDATE admission_seats
			SET filled_seats = filled_seats + 1, updated_at = NOW()
			WHERE tenant_id = ? AND session_id = ? AND class_name = ?
		`, req.TenantID, application.SessionID, application.ClassApplying).Error; err != nil {
			return fmt.Errorf("failed to update seat count: %w", err)
		}

		// TODO: Create student record from application data
		// This would be implemented when the student module is available
		// studentService.CreateFromApplication(ctx, application)

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Refresh application data
	err = s.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", req.TenantID, req.ApplicationID).
		First(&application).Error
	if err != nil {
		return nil, fmt.Errorf("failed to refresh application: %w", err)
	}

	return &application, nil
}

// DecisionFilter contains filters for listing decisions.
type DecisionFilter struct {
	TenantID      uuid.UUID
	ApplicationID *uuid.UUID
	Decision      *models.DecisionType
	DecidedBy     *uuid.UUID
	FromDate      *time.Time
	ToDate        *time.Time
}

// ListDecisions retrieves decisions with optional filtering.
func (s *DecisionService) ListDecisions(ctx context.Context, filter DecisionFilter) ([]models.AdmissionDecision, error) {
	if filter.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}

	query := s.db.WithContext(ctx).
		Model(&models.AdmissionDecision{}).
		Where("tenant_id = ?", filter.TenantID).
		Order("decision_date DESC, created_at DESC")

	if filter.ApplicationID != nil {
		query = query.Where("application_id = ?", *filter.ApplicationID)
	}

	if filter.Decision != nil {
		query = query.Where("decision = ?", *filter.Decision)
	}

	if filter.DecidedBy != nil {
		query = query.Where("decided_by = ?", *filter.DecidedBy)
	}

	if filter.FromDate != nil {
		query = query.Where("decision_date >= ?", *filter.FromDate)
	}

	if filter.ToDate != nil {
		query = query.Where("decision_date <= ?", *filter.ToDate)
	}

	var decisions []models.AdmissionDecision
	if err := query.Find(&decisions).Error; err != nil {
		return nil, fmt.Errorf("failed to list decisions: %w", err)
	}

	return decisions, nil
}

// UpdateWaitlistPosition updates the waitlist position of a waitlisted application.
func (s *DecisionService) UpdateWaitlistPosition(ctx context.Context, tenantID, applicationID uuid.UUID, newPosition int, updatedBy *uuid.UUID) (*models.AdmissionDecision, error) {
	if tenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if applicationID == uuid.Nil {
		return nil, ErrApplicationIDRequired
	}
	if newPosition <= 0 {
		return nil, ErrWaitlistPositionRequired
	}

	decision, err := s.GetDecisionByApplication(ctx, tenantID, applicationID)
	if err != nil {
		return nil, err
	}

	if decision.Decision != models.DecisionWaitlisted {
		return nil, ErrInvalidDecisionType
	}

	updates := map[string]interface{}{
		"waitlist_position": newPosition,
		"updated_by":        updatedBy,
		"updated_at":        time.Now(),
	}

	if err := s.db.WithContext(ctx).Model(decision).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update waitlist position: %w", err)
	}

	// Also update the application's waitlist position
	if err := s.db.WithContext(ctx).
		Model(&models.AdmissionApplication{}).
		Where("tenant_id = ? AND id = ?", tenantID, applicationID).
		Update("waitlist_position", newPosition).Error; err != nil {
		return nil, fmt.Errorf("failed to update application waitlist position: %w", err)
	}

	decision.WaitlistPosition = &newPosition

	return decision, nil
}

// PromoteFromWaitlist promotes a waitlisted application to approved.
func (s *DecisionService) PromoteFromWaitlist(ctx context.Context, tenantID, applicationID uuid.UUID, sectionAssigned *string, promotedBy *uuid.UUID) (*models.AdmissionDecision, error) {
	if tenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if applicationID == uuid.Nil {
		return nil, ErrApplicationIDRequired
	}

	decision, err := s.GetDecisionByApplication(ctx, tenantID, applicationID)
	if err != nil {
		return nil, err
	}

	if decision.Decision != models.DecisionWaitlisted {
		return nil, ErrInvalidDecisionType
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		// Update decision
		decisionUpdates := map[string]interface{}{
			"decision":          models.DecisionApproved,
			"section_assigned":  sectionAssigned,
			"waitlist_position": nil,
			"updated_by":        promotedBy,
			"updated_at":        now,
		}

		if err := tx.Model(decision).Updates(decisionUpdates).Error; err != nil {
			return fmt.Errorf("failed to update decision: %w", err)
		}

		// Update application status
		appUpdates := map[string]interface{}{
			"status":            models.ApplicationStatusApproved,
			"waitlist_position": nil,
			"approved_at":       &now,
			"approved_by":       promotedBy,
			"updated_by":        promotedBy,
			"updated_at":        now,
		}

		if err := tx.Model(&models.AdmissionApplication{}).
			Where("tenant_id = ? AND id = ?", tenantID, applicationID).
			Updates(appUpdates).Error; err != nil {
			return fmt.Errorf("failed to update application: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	decision.Decision = models.DecisionApproved
	decision.SectionAssigned = sectionAssigned
	decision.WaitlistPosition = nil

	return decision, nil
}
