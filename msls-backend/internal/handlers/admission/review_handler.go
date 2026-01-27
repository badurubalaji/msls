// Package admission provides HTTP handlers for admission management.
package admission

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"msls-backend/internal/middleware"
	apperrors "msls-backend/internal/pkg/errors"
	"msls-backend/internal/pkg/response"
	"msls-backend/internal/services/admission"
)

// ReviewHandler handles application review HTTP requests.
type ReviewHandler struct {
	reviewService *admission.ReviewService
}

// NewReviewHandler creates a new ReviewHandler.
func NewReviewHandler(reviewService *admission.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewService: reviewService}
}

// GetApplication godoc
// @Summary Get application by ID
// @Description Retrieves an application with documents and reviews
// @Tags Applications
// @Accept json
// @Produce json
// @Param id path string true "Application ID"
// @Success 200 {object} response.Success{data=AdmissionApplicationResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Failure 500 {object} apperrors.AppError
// @Router /api/v1/applications/{id} [get]
func (h *ReviewHandler) GetApplication(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	applicationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid application ID"))
		return
	}

	app, err := h.reviewService.GetApplicationByID(c.Request.Context(), tenantID, applicationID)
	if err != nil {
		if err == admission.ErrApplicationNotFound {
			apperrors.Abort(c, apperrors.NotFound("Application not found"))
		} else {
			apperrors.Abort(c, apperrors.InternalError("Failed to retrieve application"))
		}
		return
	}

	response.OK(c, NewAdmissionApplicationResponse(app))
}

// GetReviews godoc
// @Summary Get application reviews
// @Description Retrieves all reviews for an application
// @Tags Applications
// @Accept json
// @Produce json
// @Param id path string true "Application ID"
// @Success 200 {object} response.Success{data=ReviewListResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 500 {object} apperrors.AppError
// @Router /api/v1/applications/{id}/reviews [get]
func (h *ReviewHandler) GetReviews(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	applicationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid application ID"))
		return
	}

	reviews, err := h.reviewService.GetReviewsByApplication(c.Request.Context(), tenantID, applicationID)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve reviews"))
		return
	}

	reviewResponses := make([]ApplicationReviewResponse, len(reviews))
	for i, review := range reviews {
		reviewResponses[i] = NewApplicationReviewResponse(&review)
	}

	response.OK(c, ReviewListResponse{
		Reviews: reviewResponses,
		Total:   len(reviewResponses),
	})
}

// CreateReview godoc
// @Summary Add application review
// @Description Adds a review to an application
// @Tags Applications
// @Accept json
// @Produce json
// @Param id path string true "Application ID"
// @Param request body CreateReviewRequest true "Review details"
// @Success 201 {object} response.Success{data=ApplicationReviewResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Failure 500 {object} apperrors.AppError
// @Router /api/v1/applications/{id}/review [post]
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	applicationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid application ID"))
		return
	}

	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	svcReq := req.ToServiceRequest(tenantID, applicationID, userID)

	review, err := h.reviewService.CreateReview(c.Request.Context(), *svcReq)
	if err != nil {
		if err == admission.ErrApplicationNotFound {
			apperrors.Abort(c, apperrors.NotFound("Application not found"))
		} else if err == admission.ErrInvalidReviewType || err == admission.ErrInvalidReviewStatus {
			apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		} else {
			apperrors.Abort(c, apperrors.InternalError("Failed to create review"))
		}
		return
	}

	response.Created(c, NewApplicationReviewResponse(review))
}

// GetDocuments godoc
// @Summary Get application documents
// @Description Retrieves all documents for an application
// @Tags Applications
// @Accept json
// @Produce json
// @Param id path string true "Application ID"
// @Success 200 {object} response.Success{data=DocumentListResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 500 {object} apperrors.AppError
// @Router /api/v1/applications/{id}/documents [get]
func (h *ReviewHandler) GetDocuments(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	applicationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid application ID"))
		return
	}

	docs, err := h.reviewService.GetDocumentsByApplication(c.Request.Context(), tenantID, applicationID)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve documents"))
		return
	}

	docResponses := make([]ApplicationDocumentResponse, len(docs))
	for i, doc := range docs {
		docResponses[i] = NewApplicationDocumentResponse(&doc)
	}

	response.OK(c, DocumentListResponse{
		Documents: docResponses,
		Total:     len(docResponses),
	})
}

// VerifyDocument godoc
// @Summary Verify application document
// @Description Verifies or rejects an application document
// @Tags Applications
// @Accept json
// @Produce json
// @Param id path string true "Application ID"
// @Param docId path string true "Document ID"
// @Param request body VerifyDocumentRequest true "Verification details"
// @Success 200 {object} response.Success{data=ApplicationDocumentResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Failure 500 {object} apperrors.AppError
// @Router /api/v1/applications/{id}/documents/{docId}/verify [patch]
func (h *ReviewHandler) VerifyDocument(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	applicationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid application ID"))
		return
	}

	documentID, err := uuid.Parse(c.Param("docId"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid document ID"))
		return
	}

	var req VerifyDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	svcReq := req.ToServiceRequest(tenantID, applicationID, documentID, userID)

	doc, err := h.reviewService.VerifyDocument(c.Request.Context(), *svcReq)
	if err != nil {
		if err == admission.ErrDocumentNotFound || err == admission.ErrApplicationNotFound {
			apperrors.Abort(c, apperrors.NotFound("Document not found"))
		} else if err == admission.ErrInvalidVerificationStatus {
			apperrors.Abort(c, apperrors.BadRequest("Invalid verification status"))
		} else {
			apperrors.Abort(c, apperrors.InternalError("Failed to verify document"))
		}
		return
	}

	response.OK(c, NewApplicationDocumentResponse(doc))
}

// UpdateStatus godoc
// @Summary Update application status
// @Description Updates the status of an application
// @Tags Applications
// @Accept json
// @Produce json
// @Param id path string true "Application ID"
// @Param request body UpdateApplicationStatusRequest true "Status update"
// @Success 200 {object} response.Success{data=AdmissionApplicationResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Failure 409 {object} apperrors.AppError
// @Failure 500 {object} apperrors.AppError
// @Router /api/v1/applications/{id}/status [patch]
func (h *ReviewHandler) UpdateStatus(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	applicationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid application ID"))
		return
	}

	var req UpdateApplicationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	svcReq := req.ToServiceRequest(tenantID, applicationID, userID)

	app, err := h.reviewService.UpdateApplicationStatus(c.Request.Context(), *svcReq)
	if err != nil {
		if err == admission.ErrApplicationNotFound {
			apperrors.Abort(c, apperrors.NotFound("Application not found"))
		} else if err == admission.ErrInvalidApplicationStatus || err == admission.ErrInvalidStatusTransition {
			apperrors.Abort(c, apperrors.Conflict(err.Error()))
		} else {
			apperrors.Abort(c, apperrors.InternalError("Failed to update application status"))
		}
		return
	}

	// Fetch full application with documents and reviews
	fullApp, err := h.reviewService.GetApplicationByID(c.Request.Context(), tenantID, applicationID)
	if err != nil {
		// Still return the basic app if we can't get the full one
		response.OK(c, NewAdmissionApplicationResponse(app))
		return
	}

	response.OK(c, NewAdmissionApplicationResponse(fullApp))
}
