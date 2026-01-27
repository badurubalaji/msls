// Package promotion provides student promotion and retention processing functionality.
package promotion

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository handles database operations for promotions.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new promotion repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// DB returns the underlying database connection for transactions.
func (r *Repository) DB() *gorm.DB {
	return r.db
}

// ========================================================================
// Promotion Rules
// ========================================================================

// CreateRule creates a new promotion rule.
func (r *Repository) CreateRule(ctx context.Context, rule *PromotionRule) error {
	if err := r.db.WithContext(ctx).Create(rule).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return ErrRuleNotFound // Rule already exists for this class
		}
		return fmt.Errorf("create promotion rule: %w", err)
	}
	return nil
}

// GetRuleByID retrieves a promotion rule by ID.
func (r *Repository) GetRuleByID(ctx context.Context, tenantID, id uuid.UUID) (*PromotionRule, error) {
	var rule PromotionRule
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&rule).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRuleNotFound
		}
		return nil, fmt.Errorf("get rule by id: %w", err)
	}
	return &rule, nil
}

// GetRuleByClass retrieves a promotion rule by class ID.
func (r *Repository) GetRuleByClass(ctx context.Context, tenantID, classID uuid.UUID) (*PromotionRule, error) {
	var rule PromotionRule
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND class_id = ? AND is_active = true", tenantID, classID).
		First(&rule).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRuleNotFound
		}
		return nil, fmt.Errorf("get rule by class: %w", err)
	}
	return &rule, nil
}

// ListRules retrieves all promotion rules for a tenant.
func (r *Repository) ListRules(ctx context.Context, tenantID uuid.UUID) ([]PromotionRule, error) {
	var rules []PromotionRule
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("created_at DESC").
		Find(&rules).Error
	if err != nil {
		return nil, fmt.Errorf("list rules: %w", err)
	}
	return rules, nil
}

// UpdateRule updates a promotion rule.
func (r *Repository) UpdateRule(ctx context.Context, rule *PromotionRule) error {
	result := r.db.WithContext(ctx).Save(rule)
	if result.Error != nil {
		return fmt.Errorf("update rule: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrRuleNotFound
	}
	return nil
}

// DeleteRule deletes a promotion rule.
func (r *Repository) DeleteRule(ctx context.Context, tenantID, id uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&PromotionRule{})
	if result.Error != nil {
		return fmt.Errorf("delete rule: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrRuleNotFound
	}
	return nil
}

// UpsertRule creates or updates a promotion rule for a class.
func (r *Repository) UpsertRule(ctx context.Context, rule *PromotionRule) error {
	// Try to find existing rule
	var existing PromotionRule
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND class_id = ?", rule.TenantID, rule.ClassID).
		First(&existing).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new rule
			return r.CreateRule(ctx, rule)
		}
		return fmt.Errorf("check existing rule: %w", err)
	}

	// Update existing rule
	existing.MinAttendancePct = rule.MinAttendancePct
	existing.MinOverallMarksPct = rule.MinOverallMarksPct
	existing.MinSubjectsPassed = rule.MinSubjectsPassed
	existing.AutoPromoteOnCriteria = rule.AutoPromoteOnCriteria
	existing.IsActive = rule.IsActive
	existing.UpdatedBy = rule.UpdatedBy

	return r.UpdateRule(ctx, &existing)
}

// ========================================================================
// Promotion Batches
// ========================================================================

// CreateBatch creates a new promotion batch.
func (r *Repository) CreateBatch(ctx context.Context, batch *PromotionBatch) error {
	if err := r.db.WithContext(ctx).Create(batch).Error; err != nil {
		return fmt.Errorf("create batch: %w", err)
	}
	return nil
}

// CreateBatchWithTx creates a new promotion batch within a transaction.
func (r *Repository) CreateBatchWithTx(ctx context.Context, tx *gorm.DB, batch *PromotionBatch) error {
	if err := tx.WithContext(ctx).Create(batch).Error; err != nil {
		return fmt.Errorf("create batch: %w", err)
	}
	return nil
}

// GetBatchByID retrieves a promotion batch by ID.
func (r *Repository) GetBatchByID(ctx context.Context, tenantID, id uuid.UUID) (*PromotionBatch, error) {
	var batch PromotionBatch
	err := r.db.WithContext(ctx).
		Preload("FromAcademicYear").
		Preload("ToAcademicYear").
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&batch).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBatchNotFound
		}
		return nil, fmt.Errorf("get batch by id: %w", err)
	}
	return &batch, nil
}

// GetBatchByIDWithTx retrieves a promotion batch by ID within a transaction.
func (r *Repository) GetBatchByIDWithTx(ctx context.Context, tx *gorm.DB, tenantID, id uuid.UUID) (*PromotionBatch, error) {
	var batch PromotionBatch
	err := tx.WithContext(ctx).
		Preload("FromAcademicYear").
		Preload("ToAcademicYear").
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&batch).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBatchNotFound
		}
		return nil, fmt.Errorf("get batch by id: %w", err)
	}
	return &batch, nil
}

// ListBatches retrieves all promotion batches for a tenant.
func (r *Repository) ListBatches(ctx context.Context, tenantID uuid.UUID, status *BatchStatus) ([]PromotionBatch, error) {
	query := r.db.WithContext(ctx).
		Preload("FromAcademicYear").
		Preload("ToAcademicYear").
		Where("tenant_id = ?", tenantID)

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	var batches []PromotionBatch
	err := query.Order("created_at DESC").Find(&batches).Error
	if err != nil {
		return nil, fmt.Errorf("list batches: %w", err)
	}
	return batches, nil
}

// ListBatchesByAcademicYear retrieves batches for a specific academic year.
func (r *Repository) ListBatchesByAcademicYear(ctx context.Context, tenantID, academicYearID uuid.UUID) ([]PromotionBatch, error) {
	var batches []PromotionBatch
	err := r.db.WithContext(ctx).
		Preload("FromAcademicYear").
		Preload("ToAcademicYear").
		Where("tenant_id = ? AND from_academic_year_id = ?", tenantID, academicYearID).
		Order("created_at DESC").
		Find(&batches).Error
	if err != nil {
		return nil, fmt.Errorf("list batches by academic year: %w", err)
	}
	return batches, nil
}

// UpdateBatch updates a promotion batch.
func (r *Repository) UpdateBatch(ctx context.Context, batch *PromotionBatch) error {
	result := r.db.WithContext(ctx).Save(batch)
	if result.Error != nil {
		return fmt.Errorf("update batch: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrBatchNotFound
	}
	return nil
}

// UpdateBatchWithTx updates a promotion batch within a transaction.
func (r *Repository) UpdateBatchWithTx(ctx context.Context, tx *gorm.DB, batch *PromotionBatch) error {
	result := tx.WithContext(ctx).Save(batch)
	if result.Error != nil {
		return fmt.Errorf("update batch: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrBatchNotFound
	}
	return nil
}

// DeleteBatch deletes a promotion batch and its records.
func (r *Repository) DeleteBatch(ctx context.Context, tenantID, id uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&PromotionBatch{})
	if result.Error != nil {
		return fmt.Errorf("delete batch: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrBatchNotFound
	}
	return nil
}

// ========================================================================
// Promotion Records
// ========================================================================

// CreateRecord creates a new promotion record.
func (r *Repository) CreateRecord(ctx context.Context, record *PromotionRecord) error {
	if err := r.db.WithContext(ctx).Create(record).Error; err != nil {
		return fmt.Errorf("create record: %w", err)
	}
	return nil
}

// CreateRecordWithTx creates a new promotion record within a transaction.
func (r *Repository) CreateRecordWithTx(ctx context.Context, tx *gorm.DB, record *PromotionRecord) error {
	if err := tx.WithContext(ctx).Create(record).Error; err != nil {
		return fmt.Errorf("create record: %w", err)
	}
	return nil
}

// CreateRecordsBatch creates multiple promotion records in a batch.
func (r *Repository) CreateRecordsBatch(ctx context.Context, records []PromotionRecord) error {
	if len(records) == 0 {
		return nil
	}
	if err := r.db.WithContext(ctx).Create(&records).Error; err != nil {
		return fmt.Errorf("create records batch: %w", err)
	}
	return nil
}

// CreateRecordsBatchWithTx creates multiple promotion records in a batch within a transaction.
func (r *Repository) CreateRecordsBatchWithTx(ctx context.Context, tx *gorm.DB, records []PromotionRecord) error {
	if len(records) == 0 {
		return nil
	}
	if err := tx.WithContext(ctx).Create(&records).Error; err != nil {
		return fmt.Errorf("create records batch: %w", err)
	}
	return nil
}

// GetRecordByID retrieves a promotion record by ID.
func (r *Repository) GetRecordByID(ctx context.Context, tenantID, id uuid.UUID) (*PromotionRecord, error) {
	var record PromotionRecord
	err := r.db.WithContext(ctx).
		Preload("Student").
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, fmt.Errorf("get record by id: %w", err)
	}
	return &record, nil
}

// GetRecordByIDWithTx retrieves a promotion record by ID within a transaction.
func (r *Repository) GetRecordByIDWithTx(ctx context.Context, tx *gorm.DB, tenantID, id uuid.UUID) (*PromotionRecord, error) {
	var record PromotionRecord
	err := tx.WithContext(ctx).
		Preload("Student").
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, fmt.Errorf("get record by id: %w", err)
	}
	return &record, nil
}

// ListRecordsByBatch retrieves all promotion records for a batch.
func (r *Repository) ListRecordsByBatch(ctx context.Context, tenantID, batchID uuid.UUID) ([]PromotionRecord, error) {
	var records []PromotionRecord
	err := r.db.WithContext(ctx).
		Preload("Student").
		Where("tenant_id = ? AND batch_id = ?", tenantID, batchID).
		Order("created_at ASC").
		Find(&records).Error
	if err != nil {
		return nil, fmt.Errorf("list records by batch: %w", err)
	}
	return records, nil
}

// ListRecordsByBatchWithTx retrieves all promotion records for a batch within a transaction.
func (r *Repository) ListRecordsByBatchWithTx(ctx context.Context, tx *gorm.DB, tenantID, batchID uuid.UUID) ([]PromotionRecord, error) {
	var records []PromotionRecord
	err := tx.WithContext(ctx).
		Preload("Student").
		Where("tenant_id = ? AND batch_id = ?", tenantID, batchID).
		Order("created_at ASC").
		Find(&records).Error
	if err != nil {
		return nil, fmt.Errorf("list records by batch: %w", err)
	}
	return records, nil
}

// ListRecordsByDecision retrieves records filtered by decision.
func (r *Repository) ListRecordsByDecision(ctx context.Context, tenantID, batchID uuid.UUID, decision PromotionDecision) ([]PromotionRecord, error) {
	var records []PromotionRecord
	err := r.db.WithContext(ctx).
		Preload("Student").
		Where("tenant_id = ? AND batch_id = ? AND decision = ?", tenantID, batchID, decision).
		Order("created_at ASC").
		Find(&records).Error
	if err != nil {
		return nil, fmt.Errorf("list records by decision: %w", err)
	}
	return records, nil
}

// UpdateRecord updates a promotion record.
func (r *Repository) UpdateRecord(ctx context.Context, record *PromotionRecord) error {
	result := r.db.WithContext(ctx).Save(record)
	if result.Error != nil {
		return fmt.Errorf("update record: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

// UpdateRecordWithTx updates a promotion record within a transaction.
func (r *Repository) UpdateRecordWithTx(ctx context.Context, tx *gorm.DB, record *PromotionRecord) error {
	result := tx.WithContext(ctx).Save(record)
	if result.Error != nil {
		return fmt.Errorf("update record: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

// CountRecordsByDecision counts records by decision.
func (r *Repository) CountRecordsByDecision(ctx context.Context, tenantID, batchID uuid.UUID) (map[PromotionDecision]int, error) {
	type result struct {
		Decision PromotionDecision
		Count    int
	}

	var results []result
	err := r.db.WithContext(ctx).
		Model(&PromotionRecord{}).
		Select("decision, COUNT(*) as count").
		Where("tenant_id = ? AND batch_id = ?", tenantID, batchID).
		Group("decision").
		Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("count records by decision: %w", err)
	}

	counts := make(map[PromotionDecision]int)
	for _, r := range results {
		counts[r.Decision] = r.Count
	}
	return counts, nil
}

// GetRecordsSummary gets a summary of records for a batch.
func (r *Repository) GetRecordsSummary(ctx context.Context, tenantID, batchID uuid.UUID) (*RecordsSummary, error) {
	counts, err := r.CountRecordsByDecision(ctx, tenantID, batchID)
	if err != nil {
		return nil, err
	}

	// Count auto vs manual decisions
	var autoCount int64
	err = r.db.WithContext(ctx).
		Model(&PromotionRecord{}).
		Where("tenant_id = ? AND batch_id = ? AND auto_decided = true", tenantID, batchID).
		Count(&autoCount).Error
	if err != nil {
		return nil, fmt.Errorf("count auto decisions: %w", err)
	}

	total := 0
	for _, c := range counts {
		total += c
	}

	return &RecordsSummary{
		TotalStudents: total,
		PendingCount:  counts[DecisionPending],
		PromoteCount:  counts[DecisionPromote],
		RetainCount:   counts[DecisionRetain],
		TransferCount: counts[DecisionTransfer],
		AutoDecided:   int(autoCount),
		ManualDecided: total - int(autoCount),
	}, nil
}
