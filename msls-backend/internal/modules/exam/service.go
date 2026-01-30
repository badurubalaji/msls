package exam

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"msls-backend/internal/pkg/database/models"
)

// Service handles business logic for exam entities.
type Service struct {
	repo *Repository
}

// NewService creates a new exam service.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// ========================================
// Exam Type Service Methods
// ========================================

// ListExamTypes returns all exam types with filters.
func (s *Service) ListExamTypes(ctx context.Context, filter ExamTypeFilter) ([]models.ExamType, int64, error) {
	return s.repo.ListExamTypes(ctx, filter)
}

// GetExamType returns an exam type by ID.
func (s *Service) GetExamType(ctx context.Context, tenantID, id uuid.UUID) (*models.ExamType, error) {
	return s.repo.GetExamTypeByID(ctx, tenantID, id)
}

// CreateExamType creates a new exam type.
func (s *Service) CreateExamType(ctx context.Context, tenantID uuid.UUID, req CreateExamTypeRequest, userID uuid.UUID) (*models.ExamType, error) {
	// Validate weightage
	hundred := decimal.NewFromInt(100)
	zero := decimal.NewFromInt(0)
	if req.Weightage.LessThan(zero) || req.Weightage.GreaterThan(hundred) {
		return nil, ErrInvalidWeightage
	}

	// Validate max marks
	if req.DefaultMaxMarks <= 0 {
		return nil, ErrInvalidMaxMarks
	}

	// Validate passing marks
	if req.DefaultPassingMarks != nil && *req.DefaultPassingMarks > req.DefaultMaxMarks {
		return nil, ErrInvalidPassingMarks
	}

	// Check for duplicate code
	code := strings.ToUpper(strings.TrimSpace(req.Code))
	existing, err := s.repo.GetExamTypeByCode(ctx, tenantID, code)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrExamTypeCodeExists
	}

	// Get next display order
	maxOrder, err := s.repo.GetMaxDisplayOrder(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	examType := &models.ExamType{
		TenantID:            tenantID,
		Name:                strings.TrimSpace(req.Name),
		Code:                code,
		Description:         req.Description,
		Weightage:           req.Weightage,
		EvaluationType:      req.EvaluationType,
		DefaultMaxMarks:     req.DefaultMaxMarks,
		DefaultPassingMarks: req.DefaultPassingMarks,
		DisplayOrder:        maxOrder + 1,
		IsActive:            true,
		CreatedBy:           &userID,
		UpdatedBy:           &userID,
	}

	if err := s.repo.CreateExamType(ctx, examType); err != nil {
		return nil, err
	}

	return examType, nil
}

// UpdateExamType updates an existing exam type.
func (s *Service) UpdateExamType(ctx context.Context, tenantID, id uuid.UUID, req UpdateExamTypeRequest, userID uuid.UUID) (*models.ExamType, error) {
	examType, err := s.repo.GetExamTypeByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil {
		examType.Name = strings.TrimSpace(*req.Name)
	}

	if req.Code != nil {
		code := strings.ToUpper(strings.TrimSpace(*req.Code))
		// Check if new code conflicts with another exam type
		if code != examType.Code {
			existing, err := s.repo.GetExamTypeByCode(ctx, tenantID, code)
			if err != nil {
				return nil, err
			}
			if existing != nil && existing.ID != id {
				return nil, ErrExamTypeCodeExists
			}
			examType.Code = code
		}
	}

	if req.Description != nil {
		examType.Description = req.Description
	}

	if req.Weightage != nil {
		hundred := decimal.NewFromInt(100)
		zero := decimal.NewFromInt(0)
		if req.Weightage.LessThan(zero) || req.Weightage.GreaterThan(hundred) {
			return nil, ErrInvalidWeightage
		}
		examType.Weightage = *req.Weightage
	}

	if req.EvaluationType != nil {
		examType.EvaluationType = *req.EvaluationType
	}

	if req.DefaultMaxMarks != nil {
		if *req.DefaultMaxMarks <= 0 {
			return nil, ErrInvalidMaxMarks
		}
		examType.DefaultMaxMarks = *req.DefaultMaxMarks
	}

	if req.DefaultPassingMarks != nil {
		if *req.DefaultPassingMarks > examType.DefaultMaxMarks {
			return nil, ErrInvalidPassingMarks
		}
		examType.DefaultPassingMarks = req.DefaultPassingMarks
	}

	examType.UpdatedBy = &userID

	if err := s.repo.UpdateExamType(ctx, examType); err != nil {
		return nil, err
	}

	return examType, nil
}

// DeleteExamType deletes an exam type.
func (s *Service) DeleteExamType(ctx context.Context, tenantID, id uuid.UUID) error {
	// Check if exam type exists
	_, err := s.repo.GetExamTypeByID(ctx, tenantID, id)
	if err != nil {
		return err
	}

	// Check if exam type is in use
	usage, err := s.repo.CountExamTypeUsage(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if usage > 0 {
		return ErrExamTypeInUse
	}

	return s.repo.DeleteExamType(ctx, tenantID, id)
}

// UpdateDisplayOrders updates the display order of exam types.
func (s *Service) UpdateDisplayOrders(ctx context.Context, tenantID uuid.UUID, req UpdateDisplayOrderRequest) error {
	return s.repo.UpdateDisplayOrders(ctx, tenantID, req.Items)
}

// ToggleExamTypeActive toggles the active status of an exam type.
func (s *Service) ToggleExamTypeActive(ctx context.Context, tenantID, id uuid.UUID, isActive bool) error {
	return s.repo.ToggleActive(ctx, tenantID, id, isActive)
}
