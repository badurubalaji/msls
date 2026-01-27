// Package admission provides HTTP handlers for admission management endpoints.
package admission

import (
	"time"

	"msls-backend/internal/pkg/database/models"
)

// ============================================================================
// Request DTOs
// ============================================================================

// GenerateMeritListRequest represents the request body for generating a merit list.
type GenerateMeritListRequest struct {
	ClassName   string   `json:"className" binding:"required,max=50"`
	TestID      *string  `json:"testId,omitempty"`
	CutoffScore *float64 `json:"cutoffScore,omitempty"`
}

// UpdateCutoffRequest represents the request body for updating cutoff score.
type UpdateCutoffRequest struct {
	CutoffScore *float64 `json:"cutoffScore"`
}

// ============================================================================
// Response DTOs
// ============================================================================

// MeritListEntryResponse represents a single entry in the merit list response.
type MeritListEntryResponse struct {
	Rank           int      `json:"rank"`
	ApplicationID  string   `json:"applicationId"`
	StudentName    string   `json:"studentName"`
	Score          float64  `json:"score"`
	TestScore      *float64 `json:"testScore,omitempty"`
	InterviewScore *float64 `json:"interviewScore,omitempty"`
	PreviousMarks  *float64 `json:"previousMarks,omitempty"`
	Status         string   `json:"status"`
	ParentPhone    string   `json:"parentPhone,omitempty"`
	ParentEmail    *string  `json:"parentEmail,omitempty"`
}

// MeritListResponse represents a merit list in API responses.
type MeritListResponse struct {
	ID          string                   `json:"id"`
	SessionID   string                   `json:"sessionId"`
	ClassName   string                   `json:"className"`
	TestID      *string                  `json:"testId,omitempty"`
	GeneratedAt string                   `json:"generatedAt"`
	GeneratedBy *string                  `json:"generatedBy,omitempty"`
	CutoffScore *float64                 `json:"cutoffScore,omitempty"`
	Entries     []MeritListEntryResponse `json:"entries"`
	IsFinal     bool                     `json:"isFinal"`
	TotalCount  int                      `json:"totalCount"`
	AboveCutoff int                      `json:"aboveCutoff"`
	CreatedAt   string                   `json:"createdAt"`
}

// MeritListListResponse represents a list of merit lists.
type MeritListListResponse struct {
	MeritLists []MeritListResponse `json:"meritLists"`
	Total      int                 `json:"total"`
}

// ============================================================================
// Conversion Functions
// ============================================================================

// meritListEntryToResponse converts a MeritListEntry model to response.
func meritListEntryToResponse(entry models.MeritListEntry) MeritListEntryResponse {
	return MeritListEntryResponse{
		Rank:           entry.Rank,
		ApplicationID:  entry.ApplicationID.String(),
		StudentName:    entry.StudentName,
		Score:          entry.Score,
		TestScore:      entry.TestScore,
		InterviewScore: entry.InterviewScore,
		PreviousMarks:  entry.PreviousMarks,
		Status:         entry.Status,
		ParentPhone:    entry.ParentPhone,
		ParentEmail:    entry.ParentEmail,
	}
}

// meritListToResponse converts a MeritList model to MeritListResponse.
func meritListToResponse(ml *models.MeritList) MeritListResponse {
	entries := make([]MeritListEntryResponse, len(ml.Entries))
	for i, entry := range ml.Entries {
		entries[i] = meritListEntryToResponse(entry)
	}

	// Count entries above cutoff
	aboveCutoff := len(entries)
	if ml.CutoffScore != nil {
		aboveCutoff = 0
		for _, entry := range ml.Entries {
			if entry.Score >= *ml.CutoffScore {
				aboveCutoff++
			}
		}
	}

	resp := MeritListResponse{
		ID:          ml.ID.String(),
		SessionID:   ml.SessionID.String(),
		ClassName:   ml.ClassName,
		GeneratedAt: ml.GeneratedAt.Format(time.RFC3339),
		CutoffScore: ml.CutoffScore,
		Entries:     entries,
		IsFinal:     ml.IsFinal,
		TotalCount:  len(entries),
		AboveCutoff: aboveCutoff,
		CreatedAt:   ml.CreatedAt.Format(time.RFC3339),
	}

	if ml.TestID != nil {
		testID := ml.TestID.String()
		resp.TestID = &testID
	}

	if ml.GeneratedBy != nil {
		generatedBy := ml.GeneratedBy.String()
		resp.GeneratedBy = &generatedBy
	}

	return resp
}

// meritListsToResponses converts a slice of MeritList models to responses.
func meritListsToResponses(lists []models.MeritList) []MeritListResponse {
	responses := make([]MeritListResponse, len(lists))
	for i := range lists {
		responses[i] = meritListToResponse(&lists[i])
	}
	return responses
}
