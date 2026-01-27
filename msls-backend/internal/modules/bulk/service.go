// Package bulk provides bulk operation functionality.
package bulk

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
)

// Service handles bulk operation business logic.
type Service struct {
	repo          *Repository
	exportService *ExportService
	db            *gorm.DB
}

// NewService creates a new bulk operation service.
func NewService(db *gorm.DB, exportService *ExportService) *Service {
	return &Service{
		repo:          NewRepository(db),
		exportService: exportService,
		db:            db,
	}
}

// CreateBulkStatusUpdate creates a bulk status update operation.
func (s *Service) CreateBulkStatusUpdate(ctx context.Context, dto CreateBulkOperationDTO, newStatus models.StudentStatus) (*models.BulkOperation, error) {
	if len(dto.StudentIDs) == 0 {
		return nil, ErrNoStudentsProvided
	}
	if len(dto.StudentIDs) > MaxBulkStudents {
		return nil, ErrTooManyStudents
	}

	dto.OperationType = models.BulkOperationTypeUpdateStatus
	dto.Parameters = models.BulkOperationParams{
		"newStatus": string(newStatus),
	}

	return s.createOperation(ctx, dto)
}

// CreateExport creates an export operation.
func (s *Service) CreateExport(ctx context.Context, dto CreateBulkOperationDTO, params ExportParams) (*models.BulkOperation, error) {
	if len(dto.StudentIDs) == 0 {
		return nil, ErrNoStudentsProvided
	}
	if len(dto.StudentIDs) > MaxExportRecords {
		return nil, ErrTooManyStudents
	}

	// Validate format
	if params.Format != "xlsx" && params.Format != "csv" {
		return nil, ErrInvalidExportFormat
	}

	dto.OperationType = models.BulkOperationTypeExport
	dto.Parameters = models.BulkOperationParams{
		"format":  params.Format,
		"columns": params.Columns,
	}

	return s.createOperation(ctx, dto)
}

// createOperation creates a bulk operation with items.
func (s *Service) createOperation(ctx context.Context, dto CreateBulkOperationDTO) (*models.BulkOperation, error) {
	var op *models.BulkOperation

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create operation
		op = &models.BulkOperation{
			TenantID:      dto.TenantID,
			OperationType: dto.OperationType,
			Status:        models.BulkOperationStatusPending,
			TotalCount:    len(dto.StudentIDs),
			Parameters:    dto.Parameters,
			CreatedBy:     dto.CreatedBy,
		}

		if err := s.repo.CreateWithTx(ctx, tx, op); err != nil {
			return err
		}

		// Create items
		items := make([]models.BulkOperationItem, len(dto.StudentIDs))
		for i, studentID := range dto.StudentIDs {
			items[i] = models.BulkOperationItem{
				TenantID:    dto.TenantID,
				OperationID: op.ID,
				StudentID:   studentID,
				Status:      models.BulkItemStatusPending,
			}
		}

		if err := s.repo.CreateItems(ctx, tx, items); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return op, nil
}

// GetByID retrieves a bulk operation by ID.
func (s *Service) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*models.BulkOperation, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

// ListByUser retrieves bulk operations for a user.
func (s *Service) ListByUser(ctx context.Context, tenantID, userID uuid.UUID, limit int) ([]models.BulkOperation, error) {
	return s.repo.ListByUser(ctx, tenantID, userID, limit)
}

// ProcessStatusUpdate processes a bulk status update operation synchronously.
func (s *Service) ProcessStatusUpdate(ctx context.Context, tenantID, opID uuid.UUID, newStatus models.StudentStatus) error {
	// Mark as started
	if err := s.repo.MarkStarted(ctx, opID); err != nil {
		return fmt.Errorf("mark started: %w", err)
	}

	// Get pending items in batches
	const batchSize = 100
	for {
		items, err := s.repo.GetPendingItems(ctx, opID, batchSize)
		if err != nil {
			return fmt.Errorf("get pending items: %w", err)
		}
		if len(items) == 0 {
			break
		}

		for _, item := range items {
			err := s.processStatusUpdateItem(ctx, tenantID, &item, newStatus)
			if err != nil {
				// Mark item as failed but continue processing
				item.Status = models.BulkItemStatusFailed
				item.ErrorMessage = err.Error()
			} else {
				item.Status = models.BulkItemStatusSuccess
			}

			now := time.Now()
			item.ProcessedAt = &now

			// Update item and counts in transaction
			txErr := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
				if err := s.repo.UpdateItem(ctx, tx, &item); err != nil {
					return err
				}
				return s.repo.IncrementCounts(ctx, tx, opID, item.Status == models.BulkItemStatusSuccess)
			})

			if txErr != nil {
				return fmt.Errorf("update item: %w", txErr)
			}
		}
	}

	// Mark as completed
	if err := s.repo.MarkCompleted(ctx, opID, ""); err != nil {
		return fmt.Errorf("mark completed: %w", err)
	}

	return nil
}

// processStatusUpdateItem processes a single status update item.
func (s *Service) processStatusUpdateItem(ctx context.Context, tenantID uuid.UUID, item *models.BulkOperationItem, newStatus models.StudentStatus) error {
	// Update student status
	result := s.db.WithContext(ctx).
		Model(&models.Student{}).
		Where("tenant_id = ? AND id = ?", tenantID, item.StudentID).
		Update("status", newStatus)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("student not found")
	}
	return nil
}

// ProcessExport processes an export operation and returns the file URL.
func (s *Service) ProcessExport(ctx context.Context, tenantID, opID uuid.UUID, params ExportParams) (string, error) {
	// Mark as started
	if err := s.repo.MarkStarted(ctx, opID); err != nil {
		return "", fmt.Errorf("mark started: %w", err)
	}

	// Get operation to get student IDs
	op, err := s.repo.GetByID(ctx, tenantID, opID)
	if err != nil {
		s.repo.MarkFailed(ctx, opID, err.Error())
		return "", err
	}

	// Collect student IDs
	studentIDs := make([]uuid.UUID, len(op.Items))
	for i, item := range op.Items {
		studentIDs[i] = item.StudentID
	}

	// Use export service
	columns := params.Columns
	if len(columns) == 0 {
		columns = DefaultExportColumns
	}

	fileURL, err := s.exportService.ExportStudents(ctx, tenantID, studentIDs, params.Format, columns)
	if err != nil {
		s.repo.MarkFailed(ctx, opID, err.Error())
		return "", err
	}

	// Mark all items as success
	s.db.WithContext(ctx).
		Model(&models.BulkOperationItem{}).
		Where("operation_id = ?", opID).
		Updates(map[string]interface{}{
			"status":       models.BulkItemStatusSuccess,
			"processed_at": time.Now(),
		})

	// Update counts
	s.db.WithContext(ctx).
		Model(&models.BulkOperation{}).
		Where("id = ?", opID).
		Updates(map[string]interface{}{
			"processed_count": op.TotalCount,
			"success_count":   op.TotalCount,
		})

	// Mark as completed with result URL
	if err := s.repo.MarkCompleted(ctx, opID, fileURL); err != nil {
		return "", fmt.Errorf("mark completed: %w", err)
	}

	return fileURL, nil
}
