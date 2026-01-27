// Package promotion provides student promotion and retention processing functionality.
package promotion

import (
	"time"

	"github.com/google/uuid"
)

// ========================================================================
// Promotion Rules DTOs
// ========================================================================

// CreateRuleRequest represents a request to create or update a promotion rule.
type CreateRuleRequest struct {
	ClassID               uuid.UUID `json:"classId" binding:"required"`
	MinAttendancePct      *float64  `json:"minAttendancePct"`
	MinOverallMarksPct    *float64  `json:"minOverallMarksPct"`
	MinSubjectsPassed     *int      `json:"minSubjectsPassed"`
	AutoPromoteOnCriteria *bool     `json:"autoPromoteOnCriteria"`
}

// RuleResponse represents a promotion rule in API responses.
type RuleResponse struct {
	ID                    string   `json:"id"`
	ClassID               string   `json:"classId"`
	MinAttendancePct      *float64 `json:"minAttendancePct"`
	MinOverallMarksPct    *float64 `json:"minOverallMarksPct"`
	MinSubjectsPassed     int      `json:"minSubjectsPassed"`
	AutoPromoteOnCriteria bool     `json:"autoPromoteOnCriteria"`
	IsActive              bool     `json:"isActive"`
	CreatedAt             string   `json:"createdAt"`
	UpdatedAt             string   `json:"updatedAt"`
}

// RuleListResponse represents a list of promotion rules.
type RuleListResponse struct {
	Rules []RuleResponse `json:"rules"`
	Total int            `json:"total"`
}

// ========================================================================
// Promotion Batch DTOs
// ========================================================================

// CreateBatchRequest represents a request to create a promotion batch.
type CreateBatchRequest struct {
	FromAcademicYearID uuid.UUID  `json:"fromAcademicYearId" binding:"required"`
	ToAcademicYearID   uuid.UUID  `json:"toAcademicYearId" binding:"required"`
	FromClassID        uuid.UUID  `json:"fromClassId" binding:"required"`
	FromSectionID      *uuid.UUID `json:"fromSectionId"`
	ToClassID          *uuid.UUID `json:"toClassId"` // Target class for promotions
	Notes              string     `json:"notes"`
}

// BatchResponse represents a promotion batch in API responses.
type BatchResponse struct {
	ID                 string              `json:"id"`
	FromAcademicYear   *AcademicYearRefDTO `json:"fromAcademicYear"`
	ToAcademicYear     *AcademicYearRefDTO `json:"toAcademicYear"`
	FromClassID        string              `json:"fromClassId"`
	FromSectionID      string              `json:"fromSectionId,omitempty"`
	ToClassID          string              `json:"toClassId,omitempty"`
	Status             string              `json:"status"`
	TotalStudents      int                 `json:"totalStudents"`
	PromotedCount      int                 `json:"promotedCount"`
	RetainedCount      int                 `json:"retainedCount"`
	TransferredCount   int                 `json:"transferredCount"`
	PendingCount       int                 `json:"pendingCount"`
	ProcessedAt        string              `json:"processedAt,omitempty"`
	ProcessedBy        string              `json:"processedBy,omitempty"`
	CancelledAt        string              `json:"cancelledAt,omitempty"`
	CancelledBy        string              `json:"cancelledBy,omitempty"`
	CancellationReason string              `json:"cancellationReason,omitempty"`
	Notes              string              `json:"notes,omitempty"`
	CreatedAt          string              `json:"createdAt"`
	CreatedBy          string              `json:"createdBy,omitempty"`
}

// BatchListResponse represents a list of promotion batches.
type BatchListResponse struct {
	Batches []BatchResponse `json:"batches"`
	Total   int             `json:"total"`
}

// AcademicYearRefDTO represents an academic year reference in responses.
type AcademicYearRefDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	IsCurrent bool   `json:"isCurrent"`
}

// ========================================================================
// Promotion Record DTOs
// ========================================================================

// RecordResponse represents a promotion record in API responses.
type RecordResponse struct {
	ID                  string         `json:"id"`
	BatchID             string         `json:"batchId"`
	Student             *StudentRefDTO `json:"student"`
	FromEnrollmentID    string         `json:"fromEnrollmentId"`
	ToEnrollmentID      string         `json:"toEnrollmentId,omitempty"`
	Decision            string         `json:"decision"`
	ToClassID           string         `json:"toClassId,omitempty"`
	ToSectionID         string         `json:"toSectionId,omitempty"`
	RollNumber          string         `json:"rollNumber,omitempty"`
	AutoDecided         bool           `json:"autoDecided"`
	DecisionReason      string         `json:"decisionReason,omitempty"`
	AttendancePct       *float64       `json:"attendancePct,omitempty"`
	OverallMarksPct     *float64       `json:"overallMarksPct,omitempty"`
	SubjectsPassed      *int           `json:"subjectsPassed,omitempty"`
	OverrideBy          string         `json:"overrideBy,omitempty"`
	OverrideAt          string         `json:"overrideAt,omitempty"`
	OverrideReason      string         `json:"overrideReason,omitempty"`
	RetentionReason     string         `json:"retentionReason,omitempty"`
	TransferDestination string         `json:"transferDestination,omitempty"`
	CreatedAt           string         `json:"createdAt"`
	UpdatedAt           string         `json:"updatedAt"`
}

// RecordListResponse represents a list of promotion records.
type RecordListResponse struct {
	Records []RecordResponse `json:"records"`
	Total   int              `json:"total"`
	Summary *RecordsSummary  `json:"summary,omitempty"`
}

// RecordsSummary provides a summary of promotion records.
type RecordsSummary struct {
	TotalStudents  int `json:"totalStudents"`
	PendingCount   int `json:"pendingCount"`
	PromoteCount   int `json:"promoteCount"`
	RetainCount    int `json:"retainCount"`
	TransferCount  int `json:"transferCount"`
	AutoDecided    int `json:"autoDecided"`
	ManualDecided  int `json:"manualDecided"`
}

// StudentRefDTO represents a student reference in responses.
type StudentRefDTO struct {
	ID              string `json:"id"`
	AdmissionNumber string `json:"admissionNumber"`
	FullName        string `json:"fullName"`
	PhotoURL        string `json:"photoUrl,omitempty"`
}

// UpdateRecordRequest represents a request to update a promotion record.
type UpdateRecordRequest struct {
	Decision            *PromotionDecision `json:"decision"`
	ToClassID           *uuid.UUID         `json:"toClassId"`
	ToSectionID         *uuid.UUID         `json:"toSectionId"`
	RollNumber          *string            `json:"rollNumber"`
	OverrideReason      *string            `json:"overrideReason"`
	RetentionReason     *string            `json:"retentionReason"`
	TransferDestination *string            `json:"transferDestination"`
}

// BulkUpdateRecordsRequest represents a request to update multiple records.
type BulkUpdateRecordsRequest struct {
	RecordIDs   []uuid.UUID       `json:"recordIds" binding:"required"`
	Decision    PromotionDecision `json:"decision" binding:"required"`
	ToClassID   *uuid.UUID        `json:"toClassId"`
	ToSectionID *uuid.UUID        `json:"toSectionId"`
	Reason      string            `json:"reason"`
}

// ========================================================================
// Processing DTOs
// ========================================================================

// AutoDecideRequest represents a request to apply auto-decision rules.
type AutoDecideRequest struct {
	ApplyRules bool `json:"applyRules"` // If true, applies promotion rules
}

// ProcessBatchRequest represents a request to process a promotion batch.
type ProcessBatchRequest struct {
	GenerateRollNumbers bool `json:"generateRollNumbers"` // If true, generates sequential roll numbers
}

// CancelBatchRequest represents a request to cancel a promotion batch.
type CancelBatchRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// ========================================================================
// Report DTOs
// ========================================================================

// PromotionReportRow represents a row in the promotion report.
type PromotionReportRow struct {
	StudentAdmissionNo string `json:"studentAdmissionNo"`
	StudentName        string `json:"studentName"`
	FromClass          string `json:"fromClass"`
	FromSection        string `json:"fromSection"`
	Decision           string `json:"decision"`
	ToClass            string `json:"toClass"`
	ToSection          string `json:"toSection"`
	RollNumber         string `json:"rollNumber"`
	AttendancePct      string `json:"attendancePct"`
	MarksPct           string `json:"marksPct"`
	Reason             string `json:"reason"`
}

// ========================================================================
// Conversion Functions
// ========================================================================

// ToRuleResponse converts a PromotionRule to RuleResponse.
func ToRuleResponse(r *PromotionRule) RuleResponse {
	return RuleResponse{
		ID:                    r.ID.String(),
		ClassID:               r.ClassID.String(),
		MinAttendancePct:      r.MinAttendancePct,
		MinOverallMarksPct:    r.MinOverallMarksPct,
		MinSubjectsPassed:     r.MinSubjectsPassed,
		AutoPromoteOnCriteria: r.AutoPromoteOnCriteria,
		IsActive:              r.IsActive,
		CreatedAt:             r.CreatedAt.Format(time.RFC3339),
		UpdatedAt:             r.UpdatedAt.Format(time.RFC3339),
	}
}

// ToRuleResponses converts a slice of PromotionRule to RuleResponses.
func ToRuleResponses(rules []PromotionRule) []RuleResponse {
	responses := make([]RuleResponse, len(rules))
	for i, r := range rules {
		responses[i] = ToRuleResponse(&r)
	}
	return responses
}

// ToBatchResponse converts a PromotionBatch to BatchResponse.
func ToBatchResponse(b *PromotionBatch) BatchResponse {
	resp := BatchResponse{
		ID:               b.ID.String(),
		FromClassID:      b.FromClassID.String(),
		Status:           string(b.Status),
		TotalStudents:    b.TotalStudents,
		PromotedCount:    b.PromotedCount,
		RetainedCount:    b.RetainedCount,
		TransferredCount: b.TransferredCount,
		PendingCount:     b.TotalStudents - b.PromotedCount - b.RetainedCount - b.TransferredCount,
		Notes:            b.Notes,
		CreatedAt:        b.CreatedAt.Format(time.RFC3339),
	}

	if b.FromAcademicYear != nil {
		resp.FromAcademicYear = &AcademicYearRefDTO{
			ID:        b.FromAcademicYear.ID.String(),
			Name:      b.FromAcademicYear.Name,
			StartDate: b.FromAcademicYear.StartDate.Format("2006-01-02"),
			EndDate:   b.FromAcademicYear.EndDate.Format("2006-01-02"),
			IsCurrent: b.FromAcademicYear.IsCurrent,
		}
	}

	if b.ToAcademicYear != nil {
		resp.ToAcademicYear = &AcademicYearRefDTO{
			ID:        b.ToAcademicYear.ID.String(),
			Name:      b.ToAcademicYear.Name,
			StartDate: b.ToAcademicYear.StartDate.Format("2006-01-02"),
			EndDate:   b.ToAcademicYear.EndDate.Format("2006-01-02"),
			IsCurrent: b.ToAcademicYear.IsCurrent,
		}
	}

	if b.FromSectionID != nil {
		resp.FromSectionID = b.FromSectionID.String()
	}
	if b.ToClassID != nil {
		resp.ToClassID = b.ToClassID.String()
	}
	if b.ProcessedAt != nil {
		resp.ProcessedAt = b.ProcessedAt.Format(time.RFC3339)
	}
	if b.ProcessedBy != nil {
		resp.ProcessedBy = b.ProcessedBy.String()
	}
	if b.CancelledAt != nil {
		resp.CancelledAt = b.CancelledAt.Format(time.RFC3339)
	}
	if b.CancelledBy != nil {
		resp.CancelledBy = b.CancelledBy.String()
	}
	if b.CancellationReason != "" {
		resp.CancellationReason = b.CancellationReason
	}
	if b.CreatedBy != nil {
		resp.CreatedBy = b.CreatedBy.String()
	}

	return resp
}

// ToBatchResponses converts a slice of PromotionBatch to BatchResponses.
func ToBatchResponses(batches []PromotionBatch) []BatchResponse {
	responses := make([]BatchResponse, len(batches))
	for i, b := range batches {
		responses[i] = ToBatchResponse(&b)
	}
	return responses
}

// ToRecordResponse converts a PromotionRecord to RecordResponse.
func ToRecordResponse(r *PromotionRecord) RecordResponse {
	resp := RecordResponse{
		ID:               r.ID.String(),
		BatchID:          r.BatchID.String(),
		FromEnrollmentID: r.FromEnrollmentID.String(),
		Decision:         string(r.Decision),
		AutoDecided:      r.AutoDecided,
		DecisionReason:   r.DecisionReason,
		AttendancePct:    r.AttendancePct,
		OverallMarksPct:  r.OverallMarksPct,
		SubjectsPassed:   r.SubjectsPassed,
		RollNumber:       r.RollNumber,
		RetentionReason:  r.RetentionReason,
		TransferDestination: r.TransferDestination,
		CreatedAt:        r.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        r.UpdatedAt.Format(time.RFC3339),
	}

	if r.Student != nil {
		resp.Student = &StudentRefDTO{
			ID:              r.Student.ID.String(),
			AdmissionNumber: r.Student.AdmissionNumber,
			FullName:        r.Student.FullName(),
			PhotoURL:        r.Student.PhotoURL,
		}
	}

	if r.ToEnrollmentID != nil {
		resp.ToEnrollmentID = r.ToEnrollmentID.String()
	}
	if r.ToClassID != nil {
		resp.ToClassID = r.ToClassID.String()
	}
	if r.ToSectionID != nil {
		resp.ToSectionID = r.ToSectionID.String()
	}
	if r.OverrideBy != nil {
		resp.OverrideBy = r.OverrideBy.String()
	}
	if r.OverrideAt != nil {
		resp.OverrideAt = r.OverrideAt.Format(time.RFC3339)
	}
	if r.OverrideReason != "" {
		resp.OverrideReason = r.OverrideReason
	}

	return resp
}

// ToRecordResponses converts a slice of PromotionRecord to RecordResponses.
func ToRecordResponses(records []PromotionRecord) []RecordResponse {
	responses := make([]RecordResponse, len(records))
	for i, r := range records {
		responses[i] = ToRecordResponse(&r)
	}
	return responses
}
