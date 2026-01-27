// Package admission provides HTTP handlers for admission management endpoints.
package admission

import (
	"time"

	"msls-backend/internal/pkg/database/models"
)

// ============================================================================
// Decision Request DTOs
// ============================================================================

// CreateDecisionRequest represents the request body for making an admission decision.
type CreateDecisionRequest struct {
	Decision         string  `json:"decision" binding:"required,oneof=approved waitlisted rejected"`
	SectionAssigned  *string `json:"sectionAssigned,omitempty"`
	WaitlistPosition *int    `json:"waitlistPosition,omitempty"`
	RejectionReason  *string `json:"rejectionReason,omitempty"`
	OfferValidUntil  *string `json:"offerValidUntil,omitempty"`
	Remarks          *string `json:"remarks,omitempty"`
}

// BulkDecisionRequest represents the request body for making bulk admission decisions.
type BulkDecisionRequest struct {
	ApplicationIDs   []string `json:"applicationIds" binding:"required,min=1"`
	Decision         string   `json:"decision" binding:"required,oneof=approved waitlisted rejected"`
	SectionAssigned  *string  `json:"sectionAssigned,omitempty"`
	RejectionReason  *string  `json:"rejectionReason,omitempty"`
	OfferValidUntil  *string  `json:"offerValidUntil,omitempty"`
	Remarks          *string  `json:"remarks,omitempty"`
}

// GenerateOfferLetterRequest represents the request body for generating an offer letter.
type GenerateOfferLetterRequest struct {
	ValidUntil *string `json:"validUntil,omitempty"`
}

// AcceptOfferRequest represents the request body for accepting an offer.
type AcceptOfferRequest struct {
	PaymentMethod    *string `json:"paymentMethod,omitempty"`
	PaymentReference *string `json:"paymentReference,omitempty"`
	Remarks          *string `json:"remarks,omitempty"`
}

// EnrollRequest represents the request body for completing enrollment.
type EnrollRequest struct {
	SectionID     *string `json:"sectionId,omitempty"`
	RollNumber    *string `json:"rollNumber,omitempty"`
	AdmissionDate *string `json:"admissionDate,omitempty"`
	Remarks       *string `json:"remarks,omitempty"`
}

// ============================================================================
// Decision Response DTOs
// ============================================================================

// DecisionResponse represents an admission decision in API responses.
type DecisionResponse struct {
	ID               string  `json:"id"`
	ApplicationID    string  `json:"applicationId"`
	Decision         string  `json:"decision"`
	DecisionDate     string  `json:"decisionDate"`
	DecidedBy        *string `json:"decidedBy,omitempty"`
	SectionAssigned  *string `json:"sectionAssigned,omitempty"`
	WaitlistPosition *int    `json:"waitlistPosition,omitempty"`
	RejectionReason  *string `json:"rejectionReason,omitempty"`
	OfferLetterURL   *string `json:"offerLetterUrl,omitempty"`
	OfferValidUntil  *string `json:"offerValidUntil,omitempty"`
	OfferAccepted    *bool   `json:"offerAccepted,omitempty"`
	OfferAcceptedAt  *string `json:"offerAcceptedAt,omitempty"`
	Remarks          *string `json:"remarks,omitempty"`
	CreatedAt        string  `json:"createdAt"`
	UpdatedAt        string  `json:"updatedAt"`
}

// OfferLetterResponse represents the offer letter generation response.
type OfferLetterResponse struct {
	URL          string `json:"url"`
	ValidUntil   string `json:"validUntil"`
	GeneratedAt  string `json:"generatedAt"`
	ApplicationID string `json:"applicationId"`
}

// EnrollmentResponse represents the enrollment completion response.
type EnrollmentResponse struct {
	ApplicationID    string  `json:"applicationId"`
	StudentID        *string `json:"studentId,omitempty"`
	EnrollmentNumber *string `json:"enrollmentNumber,omitempty"`
	AdmissionDate    string  `json:"admissionDate"`
	ClassName        string  `json:"className"`
	SectionName      *string `json:"sectionName,omitempty"`
	RollNumber       *string `json:"rollNumber,omitempty"`
	Status           string  `json:"status"`
}

// BulkDecisionResponse represents the bulk decision response.
type BulkDecisionResponse struct {
	Successful int                `json:"successful"`
	Failed     int                `json:"failed"`
	Decisions  []DecisionResponse `json:"decisions"`
	Errors     []BulkDecisionError `json:"errors,omitempty"`
}

// BulkDecisionError represents an error in bulk decision processing.
type BulkDecisionError struct {
	ApplicationID string `json:"applicationId"`
	Error         string `json:"error"`
}

// ============================================================================
// Conversion Functions
// ============================================================================

// decisionToResponse converts a AdmissionDecision model to DecisionResponse.
func decisionToResponse(d *models.AdmissionDecision) DecisionResponse {
	resp := DecisionResponse{
		ID:               d.ID.String(),
		ApplicationID:    d.ApplicationID.String(),
		Decision:         string(d.Decision),
		DecisionDate:     d.DecisionDate.Format("2006-01-02"),
		SectionAssigned:  d.SectionAssigned,
		WaitlistPosition: d.WaitlistPosition,
		RejectionReason:  d.RejectionReason,
		OfferLetterURL:   d.OfferLetterURL,
		OfferAccepted:    d.OfferAccepted,
		Remarks:          d.Remarks,
		CreatedAt:        d.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        d.UpdatedAt.Format(time.RFC3339),
	}

	if d.DecidedBy != nil {
		decidedBy := d.DecidedBy.String()
		resp.DecidedBy = &decidedBy
	}

	if d.OfferValidUntil != nil {
		validUntil := d.OfferValidUntil.Format("2006-01-02")
		resp.OfferValidUntil = &validUntil
	}

	if d.OfferAcceptedAt != nil {
		acceptedAt := d.OfferAcceptedAt.Format(time.RFC3339)
		resp.OfferAcceptedAt = &acceptedAt
	}

	return resp
}

// decisionsToResponses converts a slice of AdmissionDecision models to responses.
func decisionsToResponses(decisions []models.AdmissionDecision) []DecisionResponse {
	responses := make([]DecisionResponse, len(decisions))
	for i := range decisions {
		responses[i] = decisionToResponse(&decisions[i])
	}
	return responses
}
