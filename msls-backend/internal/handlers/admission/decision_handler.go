// Package admission provides HTTP handlers for admission management endpoints.
package admission

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"msls-backend/internal/middleware"
	"msls-backend/internal/pkg/database/models"
	apperrors "msls-backend/internal/pkg/errors"
	"msls-backend/internal/pkg/response"
	admissionservice "msls-backend/internal/services/admission"
)

// DecisionHandler handles admission decision HTTP requests.
type DecisionHandler struct {
	decisionService *admissionservice.DecisionService
}

// NewDecisionHandler creates a new DecisionHandler.
func NewDecisionHandler(decisionService *admissionservice.DecisionService) *DecisionHandler {
	return &DecisionHandler{decisionService: decisionService}
}

// =============================================================================
// Decision Endpoints
// =============================================================================

// MakeDecision makes an admission decision for an application.
// @Summary Make admission decision
// @Description Make an admission decision (approve/waitlist/reject) for an application
// @Tags Admission Decisions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Application ID" format(uuid)
// @Param request body CreateDecisionRequest true "Decision details"
// @Success 201 {object} response.Success{data=DecisionResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Failure 409 {object} apperrors.AppError
// @Router /api/v1/applications/{id}/decision [post]
func (h *DecisionHandler) MakeDecision(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	applicationIDParam := c.Param("id")
	applicationID, err := uuid.Parse(applicationIDParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid application ID"))
		return
	}

	var req CreateDecisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	// Parse decision type
	decision := models.DecisionType(req.Decision)
	if !decision.IsValid() {
		apperrors.Abort(c, apperrors.BadRequest("Invalid decision type"))
		return
	}

	// Parse offer valid until date
	var offerValidUntil *time.Time
	if req.OfferValidUntil != nil {
		t, err := time.Parse("2006-01-02", *req.OfferValidUntil)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid offer valid until date format (expected YYYY-MM-DD)"))
			return
		}
		offerValidUntil = &t
	}

	createReq := admissionservice.CreateDecisionRequest{
		TenantID:         tenantID,
		ApplicationID:    applicationID,
		Decision:         decision,
		DecisionDate:     time.Now(),
		DecidedBy:        &userID,
		SectionAssigned:  req.SectionAssigned,
		WaitlistPosition: req.WaitlistPosition,
		RejectionReason:  req.RejectionReason,
		OfferValidUntil:  offerValidUntil,
		Remarks:          req.Remarks,
	}

	decisionResult, err := h.decisionService.CreateDecision(c.Request.Context(), createReq)
	if err != nil {
		switch err {
		case admissionservice.ErrApplicationNotFound:
			apperrors.Abort(c, apperrors.NotFound("Application not found"))
		case admissionservice.ErrDecisionExists:
			apperrors.Abort(c, apperrors.Conflict("Decision already exists for this application"))
		case admissionservice.ErrWaitlistPositionRequired:
			apperrors.Abort(c, apperrors.BadRequest("Waitlist position is required for waitlisted decision"))
		case admissionservice.ErrRejectionReasonRequired:
			apperrors.Abort(c, apperrors.BadRequest("Rejection reason is required for rejected decision"))
		case admissionservice.ErrInvalidDecisionType:
			apperrors.Abort(c, apperrors.BadRequest("Invalid decision type"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to make decision"))
		}
		return
	}

	response.Created(c, decisionToResponse(decisionResult))
}

// MakeBulkDecision makes bulk admission decisions.
// @Summary Make bulk admission decisions
// @Description Make admission decisions for multiple applications at once
// @Tags Admission Decisions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param request body BulkDecisionRequest true "Bulk decision details"
// @Success 200 {object} response.Success{data=BulkDecisionResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Router /api/v1/applications/bulk-decision [post]
func (h *DecisionHandler) MakeBulkDecision(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	var req BulkDecisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	// Parse decision type
	decision := models.DecisionType(req.Decision)
	if !decision.IsValid() {
		apperrors.Abort(c, apperrors.BadRequest("Invalid decision type"))
		return
	}

	// Parse offer valid until date
	var offerValidUntil *time.Time
	if req.OfferValidUntil != nil {
		t, err := time.Parse("2006-01-02", *req.OfferValidUntil)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid offer valid until date format (expected YYYY-MM-DD)"))
			return
		}
		offerValidUntil = &t
	}

	var successful int
	var failed int
	var decisions []DecisionResponse
	var errors []BulkDecisionError

	for i, appIDStr := range req.ApplicationIDs {
		appID, err := uuid.Parse(appIDStr)
		if err != nil {
			failed++
			errors = append(errors, BulkDecisionError{
				ApplicationID: appIDStr,
				Error:         "Invalid application ID",
			})
			continue
		}

		// For waitlisted decisions, auto-assign position
		var waitlistPos *int
		if decision == models.DecisionWaitlisted {
			pos := i + 1
			waitlistPos = &pos
		}

		createReq := admissionservice.CreateDecisionRequest{
			TenantID:         tenantID,
			ApplicationID:    appID,
			Decision:         decision,
			DecisionDate:     time.Now(),
			DecidedBy:        &userID,
			SectionAssigned:  req.SectionAssigned,
			WaitlistPosition: waitlistPos,
			RejectionReason:  req.RejectionReason,
			OfferValidUntil:  offerValidUntil,
			Remarks:          req.Remarks,
		}

		decisionResult, err := h.decisionService.CreateDecision(c.Request.Context(), createReq)
		if err != nil {
			failed++
			errors = append(errors, BulkDecisionError{
				ApplicationID: appIDStr,
				Error:         err.Error(),
			})
			continue
		}

		successful++
		decisions = append(decisions, decisionToResponse(decisionResult))
	}

	resp := BulkDecisionResponse{
		Successful: successful,
		Failed:     failed,
		Decisions:  decisions,
		Errors:     errors,
	}

	response.OK(c, resp)
}

// GetDecision retrieves the decision for an application.
// @Summary Get admission decision
// @Description Get the admission decision for an application
// @Tags Admission Decisions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Application ID" format(uuid)
// @Success 200 {object} response.Success{data=DecisionResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/applications/{id}/decision [get]
func (h *DecisionHandler) GetDecision(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	applicationIDParam := c.Param("id")
	applicationID, err := uuid.Parse(applicationIDParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid application ID"))
		return
	}

	decisionResult, err := h.decisionService.GetDecisionByApplication(c.Request.Context(), tenantID, applicationID)
	if err != nil {
		switch err {
		case admissionservice.ErrDecisionNotFound:
			apperrors.Abort(c, apperrors.NotFound("Decision not found"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to retrieve decision"))
		}
		return
	}

	response.OK(c, decisionToResponse(decisionResult))
}

// GenerateOfferLetter generates an offer letter for an approved application.
// @Summary Generate offer letter
// @Description Generate an offer letter PDF for an approved application
// @Tags Admission Decisions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Application ID" format(uuid)
// @Param request body GenerateOfferLetterRequest false "Offer letter parameters"
// @Success 200 {object} response.Success{data=OfferLetterResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/applications/{id}/offer-letter [post]
func (h *DecisionHandler) GenerateOfferLetter(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	applicationIDParam := c.Param("id")
	applicationID, err := uuid.Parse(applicationIDParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid application ID"))
		return
	}

	var req GenerateOfferLetterRequest
	// Bind JSON is optional here
	_ = c.ShouldBindJSON(&req)

	// Parse validity date
	var validUntil *time.Time
	if req.ValidUntil != nil {
		t, err := time.Parse("2006-01-02", *req.ValidUntil)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid date format (expected YYYY-MM-DD)"))
			return
		}
		validUntil = &t
	}

	generateReq := admissionservice.GenerateOfferLetterRequest{
		TenantID:      tenantID,
		ApplicationID: applicationID,
		ValidUntil:    validUntil,
		GeneratedBy:   &userID,
	}

	decision, err := h.decisionService.GenerateOfferLetter(c.Request.Context(), generateReq)
	if err != nil {
		switch err {
		case admissionservice.ErrDecisionNotFound:
			apperrors.Abort(c, apperrors.NotFound("Decision not found"))
		case admissionservice.ErrOfferNotFound:
			apperrors.Abort(c, apperrors.BadRequest("Application is not approved"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to generate offer letter"))
		}
		return
	}

	offerURL := ""
	if decision.OfferLetterURL != nil {
		offerURL = *decision.OfferLetterURL
	}

	validUntilStr := ""
	if decision.OfferValidUntil != nil {
		validUntilStr = decision.OfferValidUntil.Format("2006-01-02")
	}

	resp := OfferLetterResponse{
		URL:           offerURL,
		ValidUntil:    validUntilStr,
		GeneratedAt:   time.Now().Format(time.RFC3339),
		ApplicationID: applicationID.String(),
	}

	response.OK(c, resp)
}

// AcceptOffer accepts an offer for an application.
// @Summary Accept offer
// @Description Accept an admission offer
// @Tags Admission Decisions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Application ID" format(uuid)
// @Param request body AcceptOfferRequest false "Acceptance details"
// @Success 200 {object} response.Success{data=DecisionResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/applications/{id}/accept-offer [post]
func (h *DecisionHandler) AcceptOffer(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	applicationIDParam := c.Param("id")
	applicationID, err := uuid.Parse(applicationIDParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid application ID"))
		return
	}

	// Parse optional request body
	var req AcceptOfferRequest
	_ = c.ShouldBindJSON(&req)

	acceptReq := admissionservice.AcceptOfferRequest{
		TenantID:      tenantID,
		ApplicationID: applicationID,
		AcceptedBy:    &userID,
	}

	decision, err := h.decisionService.AcceptOffer(c.Request.Context(), acceptReq)
	if err != nil {
		switch err {
		case admissionservice.ErrDecisionNotFound:
			apperrors.Abort(c, apperrors.NotFound("Decision not found"))
		case admissionservice.ErrOfferNotFound:
			apperrors.Abort(c, apperrors.BadRequest("No offer found for this application"))
		case admissionservice.ErrOfferAlreadyAccepted:
			apperrors.Abort(c, apperrors.BadRequest("Offer has already been accepted"))
		case admissionservice.ErrOfferExpired:
			apperrors.Abort(c, apperrors.BadRequest("Offer has expired"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to accept offer"))
		}
		return
	}

	response.OK(c, decisionToResponse(decision))
}

// Enroll completes the enrollment process for an application.
// @Summary Complete enrollment
// @Description Complete the enrollment process after offer acceptance
// @Tags Admission Decisions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Application ID" format(uuid)
// @Param request body EnrollRequest false "Enrollment details"
// @Success 200 {object} response.Success{data=EnrollmentResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/applications/{id}/enroll [post]
func (h *DecisionHandler) Enroll(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	applicationIDParam := c.Param("id")
	applicationID, err := uuid.Parse(applicationIDParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid application ID"))
		return
	}

	// Parse optional request body
	var req EnrollRequest
	_ = c.ShouldBindJSON(&req)

	enrollReq := admissionservice.EnrollRequest{
		TenantID:      tenantID,
		ApplicationID: applicationID,
		EnrolledBy:    &userID,
	}

	application, err := h.decisionService.Enroll(c.Request.Context(), enrollReq)
	if err != nil {
		switch err {
		case admissionservice.ErrApplicationNotFound:
			apperrors.Abort(c, apperrors.NotFound("Application not found"))
		case admissionservice.ErrAlreadyEnrolled:
			apperrors.Abort(c, apperrors.BadRequest("Application is already enrolled"))
		case admissionservice.ErrInvalidApplicationStatus:
			apperrors.Abort(c, apperrors.BadRequest("Application is not in approved status"))
		case admissionservice.ErrOfferNotAccepted:
			apperrors.Abort(c, apperrors.BadRequest("Offer must be accepted before enrollment"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to complete enrollment"))
		}
		return
	}

	// Use current time as admission date
	admissionDate := time.Now().Format("2006-01-02")

	// TODO: Create actual student record and get student ID
	enrollmentNumber := fmt.Sprintf("ENR-%s-%s", time.Now().Format("2006"), application.ID.String()[:8])

	resp := EnrollmentResponse{
		ApplicationID:    application.ID.String(),
		StudentID:        nil, // TODO: Set when student module is integrated
		EnrollmentNumber: &enrollmentNumber,
		AdmissionDate:    admissionDate,
		ClassName:        application.ClassApplying,
		SectionName:      req.SectionID,
		RollNumber:       req.RollNumber,
		Status:           string(application.Status),
	}

	response.OK(c, resp)
}

// PromoteFromWaitlist promotes a waitlisted application to approved.
// @Summary Promote from waitlist
// @Description Promote a waitlisted application to approved status
// @Tags Admission Decisions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Application ID" format(uuid)
// @Param request body struct{ SectionAssigned *string `json:"sectionAssigned"` } false "Promotion details"
// @Success 200 {object} response.Success{data=DecisionResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/applications/{id}/promote [post]
func (h *DecisionHandler) PromoteFromWaitlist(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	applicationIDParam := c.Param("id")
	applicationID, err := uuid.Parse(applicationIDParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid application ID"))
		return
	}

	var req struct {
		SectionAssigned *string `json:"sectionAssigned"`
	}
	_ = c.ShouldBindJSON(&req)

	decision, err := h.decisionService.PromoteFromWaitlist(c.Request.Context(), tenantID, applicationID, req.SectionAssigned, &userID)
	if err != nil {
		switch err {
		case admissionservice.ErrDecisionNotFound:
			apperrors.Abort(c, apperrors.NotFound("Decision not found"))
		case admissionservice.ErrInvalidDecisionType:
			apperrors.Abort(c, apperrors.BadRequest("Application is not waitlisted"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to promote from waitlist"))
		}
		return
	}

	response.OK(c, decisionToResponse(decision))
}

// UpdateWaitlistPosition updates the waitlist position of a waitlisted application.
// @Summary Update waitlist position
// @Description Update the waitlist position of a waitlisted application
// @Tags Admission Decisions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Application ID" format(uuid)
// @Param request body struct{ Position int `json:"position" binding:"required,min=1"` } true "New position"
// @Success 200 {object} response.Success{data=DecisionResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/applications/{id}/waitlist-position [patch]
func (h *DecisionHandler) UpdateWaitlistPosition(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	applicationIDParam := c.Param("id")
	applicationID, err := uuid.Parse(applicationIDParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid application ID"))
		return
	}

	var req struct {
		Position int `json:"position" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	decision, err := h.decisionService.UpdateWaitlistPosition(c.Request.Context(), tenantID, applicationID, req.Position, &userID)
	if err != nil {
		switch err {
		case admissionservice.ErrDecisionNotFound:
			apperrors.Abort(c, apperrors.NotFound("Decision not found"))
		case admissionservice.ErrInvalidDecisionType:
			apperrors.Abort(c, apperrors.BadRequest("Application is not waitlisted"))
		case admissionservice.ErrWaitlistPositionRequired:
			apperrors.Abort(c, apperrors.BadRequest("Position must be greater than 0"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to update waitlist position"))
		}
		return
	}

	response.OK(c, decisionToResponse(decision))
}
