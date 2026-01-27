// Package admission provides admission management services.
package admission

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ReviewService handles application review operations.
type ReviewService struct {
	db *gorm.DB
}

// NewReviewService creates a new ReviewService instance.
func NewReviewService(db *gorm.DB) *ReviewService {
	return &ReviewService{db: db}
}

// CreateReviewRequest represents a request to create an application review.
type CreateReviewRequest struct {
	TenantID      uuid.UUID
	ApplicationID uuid.UUID
	ReviewerID    uuid.UUID
	ReviewType    ReviewType
	Status        ReviewStatus
	Comments      string
}

// VerifyDocumentRequest represents a request to verify a document.
type VerifyDocumentRequest struct {
	TenantID      uuid.UUID
	ApplicationID uuid.UUID
	DocumentID    uuid.UUID
	VerifiedBy    uuid.UUID
	Status        VerificationStatus
	Remarks       string
}

// UpdateApplicationStatusRequest represents a request to update application status.
type UpdateApplicationStatusRequest struct {
	TenantID      uuid.UUID
	ApplicationID uuid.UUID
	Status        ApplicationStatus
	UpdatedBy     uuid.UUID
	Remarks       string
}

// ReviewListFilter contains filters for listing reviews.
type ReviewListFilter struct {
	TenantID      uuid.UUID
	ApplicationID *uuid.UUID
	ReviewerID    *uuid.UUID
	ReviewType    *ReviewType
	Status        *ReviewStatus
}

// CreateReview creates a new application review.
func (s *ReviewService) CreateReview(ctx context.Context, req CreateReviewRequest) (*ApplicationReview, error) {
	if req.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if req.ApplicationID == uuid.Nil {
		return nil, ErrApplicationIDRequired
	}
	if !req.ReviewType.IsValid() {
		return nil, ErrInvalidReviewType
	}
	if !req.Status.IsValid() {
		return nil, ErrInvalidReviewStatus
	}

	// Verify application exists
	var app AdmissionApplication
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ? AND deleted_at IS NULL", req.TenantID, req.ApplicationID).
		First(&app).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrApplicationNotFound
		}
		return nil, fmt.Errorf("failed to get application: %w", err)
	}

	review := &ApplicationReview{
		TenantID:      req.TenantID,
		ApplicationID: req.ApplicationID,
		ReviewerID:    req.ReviewerID,
		ReviewType:    req.ReviewType,
		Status:        req.Status,
		Comments:      req.Comments,
	}

	// Create review and optionally update application status
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(review).Error; err != nil {
			return fmt.Errorf("failed to create review: %w", err)
		}

		// Update application status based on review
		newStatus := s.determineNewStatus(app.Status, req.ReviewType, req.Status)
		if newStatus != app.Status {
			updates := map[string]interface{}{
				"status":     newStatus,
				"updated_by": req.ReviewerID,
			}
			if err := tx.Model(&AdmissionApplication{}).
				Where("id = ?", req.ApplicationID).
				Updates(updates).Error; err != nil {
				return fmt.Errorf("failed to update application status: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return review, nil
}

// determineNewStatus determines the new application status based on review outcome.
func (s *ReviewService) determineNewStatus(currentStatus ApplicationStatus, reviewType ReviewType, reviewStatus ReviewStatus) ApplicationStatus {
	switch reviewStatus {
	case ReviewStatusRejected:
		return ApplicationStatusRejected
	case ReviewStatusApproved:
		switch reviewType {
		case ReviewTypeInitialScreening:
			return ApplicationStatusUnderReview
		case ReviewTypeDocumentVerification:
			if currentStatus == ApplicationStatusDocumentsPending {
				return ApplicationStatusUnderReview
			}
		case ReviewTypeFinalDecision:
			return ApplicationStatusApproved
		}
	case ReviewStatusPendingInfo:
		return ApplicationStatusDocumentsPending
	}
	return currentStatus
}

// GetReviewsByApplication retrieves all reviews for an application.
func (s *ReviewService) GetReviewsByApplication(ctx context.Context, tenantID, applicationID uuid.UUID) ([]ApplicationReview, error) {
	var reviews []ApplicationReview
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND application_id = ?", tenantID, applicationID).
		Order("created_at DESC").
		Find(&reviews).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get reviews: %w", err)
	}
	return reviews, nil
}

// ListReviews retrieves reviews with optional filtering.
func (s *ReviewService) ListReviews(ctx context.Context, filter ReviewListFilter) ([]ApplicationReview, error) {
	query := s.db.WithContext(ctx).
		Model(&ApplicationReview{}).
		Where("tenant_id = ?", filter.TenantID)

	if filter.ApplicationID != nil {
		query = query.Where("application_id = ?", *filter.ApplicationID)
	}
	if filter.ReviewerID != nil {
		query = query.Where("reviewer_id = ?", *filter.ReviewerID)
	}
	if filter.ReviewType != nil {
		query = query.Where("review_type = ?", *filter.ReviewType)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	query = query.Order("created_at DESC")

	var reviews []ApplicationReview
	if err := query.Find(&reviews).Error; err != nil {
		return nil, fmt.Errorf("failed to list reviews: %w", err)
	}

	return reviews, nil
}

// VerifyDocument verifies or updates the verification status of a document.
func (s *ReviewService) VerifyDocument(ctx context.Context, req VerifyDocumentRequest) (*ApplicationDocument, error) {
	if req.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if req.ApplicationID == uuid.Nil {
		return nil, ErrApplicationIDRequired
	}
	if req.DocumentID == uuid.Nil {
		return nil, ErrDocumentNotFound
	}
	if !req.Status.IsValid() {
		return nil, ErrInvalidVerificationStatus
	}

	// Get the document
	var doc ApplicationDocument
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ? AND application_id = ?", req.TenantID, req.DocumentID, req.ApplicationID).
		First(&doc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDocumentNotFound
		}
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	// Update verification status
	now := time.Now()
	updates := map[string]interface{}{
		"verification_status":  req.Status,
		"verified_by":          req.VerifiedBy,
		"verified_at":          now,
		"verification_remarks": req.Remarks,
	}

	if err := s.db.WithContext(ctx).Model(&doc).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update document verification: %w", err)
	}

	// Refresh the document
	if err := s.db.WithContext(ctx).First(&doc, doc.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to refresh document: %w", err)
	}

	return &doc, nil
}

// GetDocumentsByApplication retrieves all documents for an application.
func (s *ReviewService) GetDocumentsByApplication(ctx context.Context, tenantID, applicationID uuid.UUID) ([]ApplicationDocument, error) {
	var docs []ApplicationDocument
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND application_id = ?", tenantID, applicationID).
		Order("created_at DESC").
		Find(&docs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get documents: %w", err)
	}
	return docs, nil
}

// UpdateApplicationStatus updates the status of an application.
func (s *ReviewService) UpdateApplicationStatus(ctx context.Context, req UpdateApplicationStatusRequest) (*AdmissionApplication, error) {
	if req.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if req.ApplicationID == uuid.Nil {
		return nil, ErrApplicationIDRequired
	}
	if !req.Status.IsValid() {
		return nil, ErrInvalidApplicationStatus
	}

	// Get the application
	var app AdmissionApplication
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ? AND deleted_at IS NULL", req.TenantID, req.ApplicationID).
		First(&app).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrApplicationNotFound
		}
		return nil, fmt.Errorf("failed to get application: %w", err)
	}

	// Validate status transition
	if !s.isValidStatusTransition(app.Status, req.Status) {
		return nil, ErrInvalidStatusTransition
	}

	// Update status
	updates := map[string]interface{}{
		"status":     req.Status,
		"updated_by": req.UpdatedBy,
	}
	if req.Remarks != "" {
		updates["remarks"] = req.Remarks
	}

	if err := s.db.WithContext(ctx).Model(&app).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update application status: %w", err)
	}

	// Refresh
	if err := s.db.WithContext(ctx).First(&app, app.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to refresh application: %w", err)
	}

	return &app, nil
}

// isValidStatusTransition checks if a status transition is valid.
func (s *ReviewService) isValidStatusTransition(from, to ApplicationStatus) bool {
	// Define valid transitions
	validTransitions := map[ApplicationStatus][]ApplicationStatus{
		ApplicationStatusDraft: {
			ApplicationStatusSubmitted,
		},
		ApplicationStatusSubmitted: {
			ApplicationStatusUnderReview,
			ApplicationStatusDocumentsPending,
			ApplicationStatusRejected,
		},
		ApplicationStatusUnderReview: {
			ApplicationStatusDocumentsPending,
			ApplicationStatusTestScheduled,
			ApplicationStatusShortlisted,
			ApplicationStatusApproved,
			ApplicationStatusRejected,
		},
		ApplicationStatusDocumentsPending: {
			ApplicationStatusUnderReview,
			ApplicationStatusRejected,
		},
		ApplicationStatusTestScheduled: {
			ApplicationStatusTestCompleted,
			ApplicationStatusRejected,
		},
		ApplicationStatusTestCompleted: {
			ApplicationStatusShortlisted,
			ApplicationStatusApproved,
			ApplicationStatusRejected,
			ApplicationStatusWaitlisted,
		},
		ApplicationStatusShortlisted: {
			ApplicationStatusApproved,
			ApplicationStatusRejected,
			ApplicationStatusWaitlisted,
		},
		ApplicationStatusApproved: {
			ApplicationStatusEnrolled,
			ApplicationStatusRejected,
		},
		ApplicationStatusWaitlisted: {
			ApplicationStatusApproved,
			ApplicationStatusRejected,
		},
	}

	allowed, exists := validTransitions[from]
	if !exists {
		return false
	}

	for _, status := range allowed {
		if status == to {
			return true
		}
	}
	return false
}

// GetApplicationByID retrieves an application by ID.
func (s *ReviewService) GetApplicationByID(ctx context.Context, tenantID, applicationID uuid.UUID) (*AdmissionApplication, error) {
	var app AdmissionApplication
	err := s.db.WithContext(ctx).
		Preload("Documents").
		Preload("Reviews").
		Where("tenant_id = ? AND id = ? AND deleted_at IS NULL", tenantID, applicationID).
		First(&app).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrApplicationNotFound
		}
		return nil, fmt.Errorf("failed to get application: %w", err)
	}
	return &app, nil
}

// CountPendingDocuments counts documents pending verification for an application.
func (s *ReviewService) CountPendingDocuments(ctx context.Context, tenantID, applicationID uuid.UUID) (int64, error) {
	var count int64
	err := s.db.WithContext(ctx).
		Model(&ApplicationDocument{}).
		Where("tenant_id = ? AND application_id = ? AND verification_status = ?",
			tenantID, applicationID, VerificationStatusPending).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count pending documents: %w", err)
	}
	return count, nil
}
