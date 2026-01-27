// Package bulk provides bulk operation functionality.
package bulk

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
)

// Repository handles database operations for bulk operations.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new bulk operations repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create creates a new bulk operation with items.
func (r *Repository) Create(ctx context.Context, op *models.BulkOperation) error {
	if err := r.db.WithContext(ctx).Create(op).Error; err != nil {
		return fmt.Errorf("create bulk operation: %w", err)
	}
	return nil
}

// CreateWithTx creates a new bulk operation within a transaction.
func (r *Repository) CreateWithTx(ctx context.Context, tx *gorm.DB, op *models.BulkOperation) error {
	if err := tx.WithContext(ctx).Create(op).Error; err != nil {
		return fmt.Errorf("create bulk operation: %w", err)
	}
	return nil
}

// CreateItems creates bulk operation items.
func (r *Repository) CreateItems(ctx context.Context, tx *gorm.DB, items []models.BulkOperationItem) error {
	if len(items) == 0 {
		return nil
	}
	if err := tx.WithContext(ctx).Create(&items).Error; err != nil {
		return fmt.Errorf("create bulk operation items: %w", err)
	}
	return nil
}

// GetByID retrieves a bulk operation by ID.
func (r *Repository) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*models.BulkOperation, error) {
	var op models.BulkOperation
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.Student").
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&op).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOperationNotFound
		}
		return nil, fmt.Errorf("get bulk operation: %w", err)
	}
	return &op, nil
}

// GetByIDSimple retrieves a bulk operation without loading items.
func (r *Repository) GetByIDSimple(ctx context.Context, tenantID, id uuid.UUID) (*models.BulkOperation, error) {
	var op models.BulkOperation
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&op).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOperationNotFound
		}
		return nil, fmt.Errorf("get bulk operation: %w", err)
	}
	return &op, nil
}

// Update updates a bulk operation.
func (r *Repository) Update(ctx context.Context, op *models.BulkOperation) error {
	if err := r.db.WithContext(ctx).Save(op).Error; err != nil {
		return fmt.Errorf("update bulk operation: %w", err)
	}
	return nil
}

// UpdateWithTx updates a bulk operation within a transaction.
func (r *Repository) UpdateWithTx(ctx context.Context, tx *gorm.DB, op *models.BulkOperation) error {
	if err := tx.WithContext(ctx).Save(op).Error; err != nil {
		return fmt.Errorf("update bulk operation: %w", err)
	}
	return nil
}

// UpdateItem updates a bulk operation item.
func (r *Repository) UpdateItem(ctx context.Context, tx *gorm.DB, item *models.BulkOperationItem) error {
	if err := tx.WithContext(ctx).Save(item).Error; err != nil {
		return fmt.Errorf("update bulk operation item: %w", err)
	}
	return nil
}

// IncrementCounts increments the processed/success/failure counts.
func (r *Repository) IncrementCounts(ctx context.Context, tx *gorm.DB, opID uuid.UUID, success bool) error {
	updates := map[string]interface{}{
		"processed_count": gorm.Expr("processed_count + 1"),
	}
	if success {
		updates["success_count"] = gorm.Expr("success_count + 1")
	} else {
		updates["failure_count"] = gorm.Expr("failure_count + 1")
	}

	if err := tx.WithContext(ctx).
		Model(&models.BulkOperation{}).
		Where("id = ?", opID).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("increment counts: %w", err)
	}
	return nil
}

// MarkStarted marks the operation as started.
func (r *Repository) MarkStarted(ctx context.Context, opID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.BulkOperation{}).
		Where("id = ?", opID).
		Updates(map[string]interface{}{
			"status":     models.BulkOperationStatusProcessing,
			"started_at": now,
		}).Error
}

// MarkCompleted marks the operation as completed.
func (r *Repository) MarkCompleted(ctx context.Context, opID uuid.UUID, resultURL string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.BulkOperation{}).
		Where("id = ?", opID).
		Updates(map[string]interface{}{
			"status":       models.BulkOperationStatusCompleted,
			"completed_at": now,
			"result_url":   resultURL,
		}).Error
}

// MarkFailed marks the operation as failed.
func (r *Repository) MarkFailed(ctx context.Context, opID uuid.UUID, errMsg string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.BulkOperation{}).
		Where("id = ?", opID).
		Updates(map[string]interface{}{
			"status":        models.BulkOperationStatusFailed,
			"completed_at":  now,
			"error_message": errMsg,
		}).Error
}

// ListByUser retrieves bulk operations for a user.
func (r *Repository) ListByUser(ctx context.Context, tenantID, userID uuid.UUID, limit int) ([]models.BulkOperation, error) {
	if limit <= 0 {
		limit = 20
	}

	var ops []models.BulkOperation
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND created_by = ?", tenantID, userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&ops).Error
	if err != nil {
		return nil, fmt.Errorf("list bulk operations: %w", err)
	}
	return ops, nil
}

// GetPendingItems retrieves pending items for an operation.
func (r *Repository) GetPendingItems(ctx context.Context, opID uuid.UUID, limit int) ([]models.BulkOperationItem, error) {
	var items []models.BulkOperationItem
	err := r.db.WithContext(ctx).
		Where("operation_id = ? AND status = ?", opID, models.BulkItemStatusPending).
		Preload("Student").
		Limit(limit).
		Find(&items).Error
	if err != nil {
		return nil, fmt.Errorf("get pending items: %w", err)
	}
	return items, nil
}

// DB returns the underlying database connection.
func (r *Repository) DB() *gorm.DB {
	return r.db
}
