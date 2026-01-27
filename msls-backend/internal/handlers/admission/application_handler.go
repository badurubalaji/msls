// Package admission provides HTTP handlers for admission management endpoints.
package admission

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"msls-backend/internal/middleware"
	"msls-backend/internal/pkg/database/models"
	apperrors "msls-backend/internal/pkg/errors"
	"msls-backend/internal/pkg/response"
	admissionservice "msls-backend/internal/services/admission"
)

// ApplicationHandler handles admission application HTTP requests.
type ApplicationHandler struct {
	applicationService *admissionservice.ApplicationService
}

// NewApplicationHandler creates a new ApplicationHandler.
func NewApplicationHandler(applicationService *admissionservice.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{applicationService: applicationService}
}

// =============================================================================
// Application CRUD Endpoints
// =============================================================================

// List returns all applications for the tenant.
// @Summary List applications
// @Description Get all admission applications for the current tenant
// @Tags Admission Applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param sessionId query string false "Filter by session ID"
// @Param status query string false "Filter by status"
// @Param className query string false "Filter by class name"
// @Param search query string false "Search by name or application number"
// @Success 200 {object} response.Success{data=ApplicationListResponse}
// @Failure 401 {object} apperrors.AppError
// @Failure 403 {object} apperrors.AppError
// @Router /api/v1/applications [get]
func (h *ApplicationHandler) List(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	filter := admissionservice.ListApplicationFilter{
		TenantID: tenantID,
		Search:   c.Query("search"),
	}

	// Parse optional branch ID
	if branchIDStr := c.Query("branchId"); branchIDStr != "" {
		branchID, err := uuid.Parse(branchIDStr)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid branch ID"))
			return
		}
		filter.BranchID = &branchID
	}

	// Parse optional session ID
	if sessionIDStr := c.Query("sessionId"); sessionIDStr != "" {
		sessionID, err := uuid.Parse(sessionIDStr)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid session ID"))
			return
		}
		filter.SessionID = &sessionID
	}

	// Parse optional status
	if statusStr := c.Query("status"); statusStr != "" {
		status := models.ApplicationStatus(statusStr)
		if !status.IsValid() {
			apperrors.Abort(c, apperrors.BadRequest("Invalid status value"))
			return
		}
		filter.Status = &status
	}

	// Parse optional class name
	if className := c.Query("className"); className != "" {
		filter.ClassApplying = &className
	}

	applications, err := h.applicationService.List(c.Request.Context(), filter)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve applications"))
		return
	}

	total, err := h.applicationService.Count(c.Request.Context(), filter)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to count applications"))
		return
	}

	resp := ApplicationListResponse{
		Applications: applicationsToResponses(applications),
		Total:        total,
	}

	response.OK(c, resp)
}

// GetByID returns an application by ID.
// @Summary Get application by ID
// @Description Get an admission application with full details
// @Tags Admission Applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Application ID" format(uuid)
// @Success 200 {object} response.Success{data=ApplicationResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/applications/{id} [get]
func (h *ApplicationHandler) GetByID(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid application ID"))
		return
	}

	application, err := h.applicationService.GetByID(c.Request.Context(), tenantID, id)
	if err != nil {
		switch err {
		case admissionservice.ErrApplicationNotFound:
			apperrors.Abort(c, apperrors.NotFound("Application not found"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to retrieve application"))
		}
		return
	}

	// Get parents and documents
	parents, _ := h.applicationService.GetParents(c.Request.Context(), tenantID, id)
	documents, _ := h.applicationService.GetDocuments(c.Request.Context(), tenantID, id)

	resp := applicationToResponseWithRelations(application, parents, documents)
	response.OK(c, resp)
}

// Create creates a new application.
// @Summary Create application
// @Description Create a new admission application
// @Tags Admission Applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param request body CreateApplicationRequest true "Application details"
// @Success 201 {object} response.Success{data=ApplicationResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 409 {object} apperrors.AppError
// @Router /api/v1/applications [post]
func (h *ApplicationHandler) Create(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	var req CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	// Parse session ID
	sessionID, err := uuid.Parse(req.SessionID)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid session ID"))
		return
	}

	// Parse optional enquiry ID
	var enquiryID *uuid.UUID
	if req.EnquiryID != nil {
		eid, err := uuid.Parse(*req.EnquiryID)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid enquiry ID"))
			return
		}
		enquiryID = &eid
	}

	createReq := admissionservice.CreateApplicationRequest{
		TenantID:      tenantID,
		SessionID:     sessionID,
		EnquiryID:     enquiryID,
		StudentName:   req.StudentName,
		ClassApplying: req.ClassName,
		CreatedBy:     &userID,
	}

	// Set applicant details if provided
	if req.ApplicantDetails != nil {
		createReq.Gender = req.ApplicantDetails.Gender
		if req.ApplicantDetails.DateOfBirth != "" {
			if dob, err := time.Parse("2006-01-02", req.ApplicantDetails.DateOfBirth); err == nil {
				createReq.DateOfBirth = &dob
			}
		}
		createReq.Nationality = req.ApplicantDetails.Nationality
		createReq.Religion = req.ApplicantDetails.Religion
		createReq.Category = req.ApplicantDetails.Category
		createReq.BloodGroup = req.ApplicantDetails.BloodGroup
		createReq.AddressLine1 = req.ApplicantDetails.Address
		createReq.City = req.ApplicantDetails.City
		createReq.State = req.ApplicantDetails.State
		createReq.PostalCode = req.ApplicantDetails.PinCode
		createReq.PreviousSchool = req.ApplicantDetails.PreviousSchool
	}

	// Set parent info if provided
	if req.ParentInfo != nil {
		createReq.FatherName = req.ParentInfo.FatherName
		createReq.FatherPhone = req.ParentInfo.FatherPhone
		createReq.FatherEmail = req.ParentInfo.FatherEmail
		createReq.FatherOccupation = req.ParentInfo.FatherOccupation
		createReq.MotherName = req.ParentInfo.MotherName
		createReq.MotherPhone = req.ParentInfo.MotherPhone
		createReq.MotherEmail = req.ParentInfo.MotherEmail
		createReq.MotherOccupation = req.ParentInfo.MotherOccupation
		createReq.GuardianName = req.ParentInfo.GuardianName
		createReq.GuardianPhone = req.ParentInfo.GuardianPhone
		createReq.GuardianEmail = req.ParentInfo.GuardianEmail
		createReq.GuardianRelation = req.ParentInfo.GuardianRelation
	}

	application, err := h.applicationService.Create(c.Request.Context(), createReq)
	if err != nil {
		switch err {
		case admissionservice.ErrSessionNotFound:
			apperrors.Abort(c, apperrors.NotFound("Session not found"))
		case admissionservice.ErrSessionClosed:
			apperrors.Abort(c, apperrors.BadRequest("Session is closed"))
		case admissionservice.ErrStudentNameRequired:
			apperrors.Abort(c, apperrors.BadRequest("Student name is required"))
		case admissionservice.ErrClassApplyingRequired:
			apperrors.Abort(c, apperrors.BadRequest("Class is required"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to create application"))
		}
		return
	}

	response.Created(c, applicationToResponse(application))
}

// Update updates an application.
// @Summary Update application
// @Description Update an existing admission application (only draft applications)
// @Tags Admission Applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Application ID" format(uuid)
// @Param request body UpdateApplicationRequest true "Application updates"
// @Success 200 {object} response.Success{data=ApplicationResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/applications/{id} [put]
func (h *ApplicationHandler) Update(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid application ID"))
		return
	}

	var req UpdateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	updateReq := admissionservice.UpdateApplicationRequest{
		StudentName:   req.StudentName,
		ClassApplying: req.ClassName,
		UpdatedBy:     &userID,
	}

	// Set applicant details if provided
	if req.ApplicantDetails != nil {
		if req.ApplicantDetails.Gender != "" {
			updateReq.Gender = &req.ApplicantDetails.Gender
		}
		if req.ApplicantDetails.DateOfBirth != "" {
			// Parse date string to time.Time
			if dob, err := time.Parse("2006-01-02", req.ApplicantDetails.DateOfBirth); err == nil {
				updateReq.DateOfBirth = &dob
			}
		}
		if req.ApplicantDetails.Nationality != "" {
			updateReq.Nationality = &req.ApplicantDetails.Nationality
		}
		if req.ApplicantDetails.Religion != "" {
			updateReq.Religion = &req.ApplicantDetails.Religion
		}
		if req.ApplicantDetails.Category != "" {
			updateReq.Category = &req.ApplicantDetails.Category
		}
		if req.ApplicantDetails.BloodGroup != "" {
			updateReq.BloodGroup = &req.ApplicantDetails.BloodGroup
		}
		if req.ApplicantDetails.Address != "" {
			updateReq.AddressLine1 = &req.ApplicantDetails.Address
		}
		if req.ApplicantDetails.City != "" {
			updateReq.City = &req.ApplicantDetails.City
		}
		if req.ApplicantDetails.State != "" {
			updateReq.State = &req.ApplicantDetails.State
		}
		if req.ApplicantDetails.PinCode != "" {
			updateReq.PostalCode = &req.ApplicantDetails.PinCode
		}
		if req.ApplicantDetails.PreviousSchool != "" {
			updateReq.PreviousSchool = &req.ApplicantDetails.PreviousSchool
		}
	}

	// Set parent info if provided
	if req.ParentInfo != nil {
		if req.ParentInfo.FatherName != "" {
			updateReq.FatherName = &req.ParentInfo.FatherName
		}
		if req.ParentInfo.FatherPhone != "" {
			updateReq.FatherPhone = &req.ParentInfo.FatherPhone
		}
		if req.ParentInfo.FatherEmail != "" {
			updateReq.FatherEmail = &req.ParentInfo.FatherEmail
		}
		if req.ParentInfo.FatherOccupation != "" {
			updateReq.FatherOccupation = &req.ParentInfo.FatherOccupation
		}
		if req.ParentInfo.MotherName != "" {
			updateReq.MotherName = &req.ParentInfo.MotherName
		}
		if req.ParentInfo.MotherPhone != "" {
			updateReq.MotherPhone = &req.ParentInfo.MotherPhone
		}
		if req.ParentInfo.MotherEmail != "" {
			updateReq.MotherEmail = &req.ParentInfo.MotherEmail
		}
		if req.ParentInfo.MotherOccupation != "" {
			updateReq.MotherOccupation = &req.ParentInfo.MotherOccupation
		}
		if req.ParentInfo.GuardianName != "" {
			updateReq.GuardianName = &req.ParentInfo.GuardianName
		}
		if req.ParentInfo.GuardianPhone != "" {
			updateReq.GuardianPhone = &req.ParentInfo.GuardianPhone
		}
		if req.ParentInfo.GuardianEmail != "" {
			updateReq.GuardianEmail = &req.ParentInfo.GuardianEmail
		}
		if req.ParentInfo.GuardianRelation != "" {
			updateReq.GuardianRelation = &req.ParentInfo.GuardianRelation
		}
	}

	application, err := h.applicationService.Update(c.Request.Context(), tenantID, id, updateReq)
	if err != nil {
		switch err {
		case admissionservice.ErrApplicationNotFound:
			apperrors.Abort(c, apperrors.NotFound("Application not found"))
		case admissionservice.ErrCannotUpdateApplication:
			apperrors.Abort(c, apperrors.BadRequest("Cannot update application in current status"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to update application"))
		}
		return
	}

	response.OK(c, applicationToResponse(application))
}

// Submit submits an application.
// @Summary Submit application
// @Description Submit a draft application for review
// @Tags Admission Applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Application ID" format(uuid)
// @Success 200 {object} response.Success{data=ApplicationResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/applications/{id}/submit [post]
func (h *ApplicationHandler) Submit(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid application ID"))
		return
	}

	application, err := h.applicationService.Submit(c.Request.Context(), tenantID, id, &userID)
	if err != nil {
		switch err {
		case admissionservice.ErrApplicationNotFound:
			apperrors.Abort(c, apperrors.NotFound("Application not found"))
		case admissionservice.ErrApplicationNotInDraft:
			apperrors.Abort(c, apperrors.BadRequest("Application is not in draft status"))
		case admissionservice.ErrMissingRequiredFields:
			apperrors.Abort(c, apperrors.BadRequest("Missing required fields"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to submit application"))
		}
		return
	}

	response.OK(c, applicationToResponse(application))
}

// UpdateStage updates the stage of an application.
// @Summary Update application stage
// @Description Update the processing stage of an application
// @Tags Admission Applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Application ID" format(uuid)
// @Param request body UpdateStageRequest true "Stage update"
// @Success 200 {object} response.Success{data=ApplicationResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/applications/{id}/stage [patch]
func (h *ApplicationHandler) UpdateStage(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid application ID"))
		return
	}

	var req UpdateStageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	newStage := models.ApplicationStatus(req.NewStage)
	if !newStage.IsValid() {
		apperrors.Abort(c, apperrors.BadRequest("Invalid stage value"))
		return
	}

	updateReq := admissionservice.UpdateStageRequest{
		NewStage:  newStage,
		Remarks:   req.Remarks,
		ChangedBy: &userID,
	}

	application, err := h.applicationService.UpdateStage(c.Request.Context(), tenantID, id, updateReq)
	if err != nil {
		switch e := err.(type) {
		case *admissionservice.StageTransitionError:
			// Build a detailed error message
			validOptions := "none"
			if len(e.ValidTransitions) > 0 {
				validOptions = strings.Join(e.ValidTransitions, ", ")
			}
			errorMsg := fmt.Sprintf(
				"Cannot transition from '%s' to '%s'. Valid next stages are: %s",
				e.CurrentStatus,
				e.RequestedStatus,
				validOptions,
			)
			apperrors.Abort(c, apperrors.BadRequest(errorMsg))
		default:
			if err == admissionservice.ErrApplicationNotFound {
				apperrors.Abort(c, apperrors.NotFound("Application not found"))
			} else if err == admissionservice.ErrInvalidStageTransition {
				apperrors.Abort(c, apperrors.BadRequest("Invalid stage transition"))
			} else {
				apperrors.Abort(c, apperrors.InternalError("Failed to update stage"))
			}
		}
		return
	}

	response.OK(c, applicationToResponse(application))
}

// Delete deletes an application.
// @Summary Delete application
// @Description Delete a draft application
// @Tags Admission Applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Application ID" format(uuid)
// @Success 204 "No Content"
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/applications/{id} [delete]
func (h *ApplicationHandler) Delete(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid application ID"))
		return
	}

	err = h.applicationService.Delete(c.Request.Context(), tenantID, id)
	if err != nil {
		switch err {
		case admissionservice.ErrApplicationNotFound:
			apperrors.Abort(c, apperrors.NotFound("Application not found"))
		case admissionservice.ErrCannotDeleteSubmittedApplication:
			apperrors.Abort(c, apperrors.BadRequest("Cannot delete a submitted application"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to delete application"))
		}
		return
	}

	response.NoContent(c)
}

// =============================================================================
// Parent Endpoints
// =============================================================================

// ListParents returns all parents for an application.
// @Summary List application parents
// @Description Get all parents for an application
// @Tags Admission Applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Application ID" format(uuid)
// @Success 200 {object} response.Success{data=[]ApplicationParentResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/applications/{id}/parents [get]
func (h *ApplicationHandler) ListParents(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	idParam := c.Param("id")
	applicationID, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid application ID"))
		return
	}

	parents, err := h.applicationService.GetParents(c.Request.Context(), tenantID, applicationID)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve parents"))
		return
	}

	responses := make([]ApplicationParentResponse, len(parents))
	for i, p := range parents {
		responses[i] = parentToResponse(&p)
	}

	response.OK(c, responses)
}

// AddParent adds a parent to an application.
// @Summary Add parent to application
// @Description Add a new parent to an application
// @Tags Admission Applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Application ID" format(uuid)
// @Param request body CreateParentRequest true "Parent details"
// @Success 201 {object} response.Success{data=ApplicationParentResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/applications/{id}/parents [post]
func (h *ApplicationHandler) AddParent(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	idParam := c.Param("id")
	applicationID, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid application ID"))
		return
	}

	var req CreateParentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	relation := models.ParentRelation(req.Relation)
	if !relation.IsValid() {
		apperrors.Abort(c, apperrors.BadRequest("Invalid relation value"))
		return
	}

	addReq := admissionservice.AddParentRequest{
		TenantID:      tenantID,
		ApplicationID: applicationID,
		Relation:      relation,
		Name:          req.Name,
		Phone:         req.Phone,
		Email:         req.Email,
		Occupation:    req.Occupation,
		Education:     req.Education,
		AnnualIncome:  req.AnnualIncome,
	}

	parent, err := h.applicationService.AddParent(c.Request.Context(), addReq)
	if err != nil {
		switch err {
		case admissionservice.ErrApplicationNotFound:
			apperrors.Abort(c, apperrors.NotFound("Application not found"))
		case admissionservice.ErrInvalidParentRelation:
			apperrors.Abort(c, apperrors.BadRequest("Invalid parent relation"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to add parent"))
		}
		return
	}

	response.Created(c, parentToResponse(parent))
}

// UpdateParent updates a parent record.
// @Summary Update parent
// @Description Update a parent record
// @Tags Admission Applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Application ID" format(uuid)
// @Param parentId path string true "Parent ID" format(uuid)
// @Param request body UpdateParentRequest true "Parent updates"
// @Success 200 {object} response.Success{data=ApplicationParentResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/applications/{id}/parents/{parentId} [put]
func (h *ApplicationHandler) UpdateParent(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	idParam := c.Param("id")
	applicationID, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid application ID"))
		return
	}

	parentIDParam := c.Param("parentId")
	parentID, err := uuid.Parse(parentIDParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid parent ID"))
		return
	}

	var req UpdateParentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	updateReq := admissionservice.UpdateParentRequest{
		Name:         req.Name,
		Phone:        req.Phone,
		Email:        req.Email,
		Occupation:   req.Occupation,
		Education:    req.Education,
		AnnualIncome: req.AnnualIncome,
	}

	if req.Relation != nil {
		relation := models.ParentRelation(*req.Relation)
		if !relation.IsValid() {
			apperrors.Abort(c, apperrors.BadRequest("Invalid relation value"))
			return
		}
		updateReq.Relation = &relation
	}

	parent, err := h.applicationService.UpdateParent(c.Request.Context(), tenantID, applicationID, parentID, updateReq)
	if err != nil {
		switch err {
		case admissionservice.ErrParentNotFound:
			apperrors.Abort(c, apperrors.NotFound("Parent not found"))
		case admissionservice.ErrInvalidParentRelation:
			apperrors.Abort(c, apperrors.BadRequest("Invalid parent relation"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to update parent"))
		}
		return
	}

	response.OK(c, parentToResponse(parent))
}

// DeleteParent deletes a parent record.
// @Summary Delete parent
// @Description Delete a parent from an application
// @Tags Admission Applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Application ID" format(uuid)
// @Param parentId path string true "Parent ID" format(uuid)
// @Success 204 "No Content"
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/applications/{id}/parents/{parentId} [delete]
func (h *ApplicationHandler) DeleteParent(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	idParam := c.Param("id")
	applicationID, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid application ID"))
		return
	}

	parentIDParam := c.Param("parentId")
	parentID, err := uuid.Parse(parentIDParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid parent ID"))
		return
	}

	err = h.applicationService.DeleteParent(c.Request.Context(), tenantID, applicationID, parentID)
	if err != nil {
		switch err {
		case admissionservice.ErrParentNotFound:
			apperrors.Abort(c, apperrors.NotFound("Parent not found"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to delete parent"))
		}
		return
	}

	response.NoContent(c)
}

// =============================================================================
// Document Endpoints
// =============================================================================

// ListDocuments returns all documents for an application.
// @Summary List application documents
// @Description Get all documents for an application
// @Tags Admission Applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Application ID" format(uuid)
// @Success 200 {object} response.Success{data=[]ApplicationDocumentResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/applications/{id}/documents [get]
func (h *ApplicationHandler) ListDocuments(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	idParam := c.Param("id")
	applicationID, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid application ID"))
		return
	}

	documents, err := h.applicationService.GetDocuments(c.Request.Context(), tenantID, applicationID)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve documents"))
		return
	}

	responses := make([]AppDocumentResponse, len(documents))
	for i, d := range documents {
		responses[i] = documentToResponse(&d)
	}

	response.OK(c, responses)
}

// AddDocument adds a document to an application.
// @Summary Add document to application
// @Description Add a new document to an application (supports multipart form upload)
// @Tags Admission Applications
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Application ID" format(uuid)
// @Param file formData file true "Document file"
// @Param documentType formData string true "Document type"
// @Success 201 {object} response.Success{data=ApplicationDocumentResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/applications/{id}/documents [post]
func (h *ApplicationHandler) AddDocument(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	idParam := c.Param("id")
	applicationID, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid application ID"))
		return
	}

	// Get document type from form
	documentType := c.PostForm("documentType")
	if documentType == "" {
		apperrors.Abort(c, apperrors.BadRequest("Document type is required"))
		return
	}

	docType := models.DocumentType(documentType)
	if !docType.IsValid() {
		apperrors.Abort(c, apperrors.BadRequest("Invalid document type"))
		return
	}

	// Get file from form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("File is required"))
		return
	}
	defer file.Close()

	// Create uploads directory if it doesn't exist
	uploadDir := filepath.Join("uploads", "documents", tenantID.String(), applicationID.String())
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to create upload directory"))
		return
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	uniqueFileName := fmt.Sprintf("%s_%s%s", documentType, uuid.New().String()[:8], ext)
	filePath := filepath.Join(uploadDir, uniqueFileName)

	// Save file
	dst, err := os.Create(filePath)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to save file"))
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to save file"))
		return
	}

	// Generate URL (relative path for local storage)
	fileURL := "/" + filePath

	addReq := admissionservice.AddDocumentRequest{
		TenantID:      tenantID,
		ApplicationID: applicationID,
		DocumentType:  docType,
		FileURL:       fileURL,
		FileName:      header.Filename,
		FileSize:      header.Size,
		MimeType:      header.Header.Get("Content-Type"),
	}

	document, err := h.applicationService.AddDocument(c.Request.Context(), addReq)
	if err != nil {
		// Clean up uploaded file on error
		os.Remove(filePath)

		switch err {
		case admissionservice.ErrApplicationNotFound:
			apperrors.Abort(c, apperrors.NotFound("Application not found"))
		case admissionservice.ErrInvalidDocumentType:
			apperrors.Abort(c, apperrors.BadRequest("Invalid document type"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to add document"))
		}
		return
	}

	response.Created(c, documentToResponse(document))
}

// VerifyDocument verifies or rejects a document.
// @Summary Verify document
// @Description Verify or reject a document
// @Tags Admission Applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Application ID" format(uuid)
// @Param documentId path string true "Document ID" format(uuid)
// @Param request body AppVerifyDocumentRequest true "Verification details"
// @Success 200 {object} response.Success{data=AppDocumentResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/applications/{id}/documents/{documentId}/verify [patch]
func (h *ApplicationHandler) VerifyDocument(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	idParam := c.Param("id")
	applicationID, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid application ID"))
		return
	}

	documentIDParam := c.Param("documentId")
	documentID, err := uuid.Parse(documentIDParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid document ID"))
		return
	}

	var req AppVerifyDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	// Convert boolean IsVerified to VerificationStatus enum
	var verificationStatus models.VerificationStatus
	if req.IsVerified {
		verificationStatus = models.VerificationStatusVerified
	} else if req.RejectionReason != "" {
		verificationStatus = models.VerificationStatusRejected
	} else {
		verificationStatus = models.VerificationStatusPending
	}

	verifyReq := admissionservice.DocumentVerificationInput{
		VerificationStatus:  verificationStatus,
		VerificationRemarks: req.RejectionReason,
		VerifiedBy:          &userID,
	}

	document, err := h.applicationService.VerifyDocument(c.Request.Context(), tenantID, applicationID, documentID, verifyReq)
	if err != nil {
		switch err {
		case admissionservice.ErrDocumentNotFound:
			apperrors.Abort(c, apperrors.NotFound("Document not found"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to verify document"))
		}
		return
	}

	response.OK(c, documentToResponse(document))
}

// DeleteDocument deletes a document.
// @Summary Delete document
// @Description Delete a document from an application
// @Tags Admission Applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Application ID" format(uuid)
// @Param documentId path string true "Document ID" format(uuid)
// @Success 204 "No Content"
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/applications/{id}/documents/{documentId} [delete]
func (h *ApplicationHandler) DeleteDocument(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	idParam := c.Param("id")
	applicationID, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid application ID"))
		return
	}

	documentIDParam := c.Param("documentId")
	documentID, err := uuid.Parse(documentIDParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid document ID"))
		return
	}

	err = h.applicationService.DeleteDocument(c.Request.Context(), tenantID, applicationID, documentID)
	if err != nil {
		switch err {
		case admissionservice.ErrDocumentNotFound:
			apperrors.Abort(c, apperrors.NotFound("Document not found"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to delete document"))
		}
		return
	}

	response.NoContent(c)
}

// =============================================================================
// Public Status Check Endpoint
// =============================================================================

// CheckStatus checks the status of an application (public API).
// @Summary Check application status (public)
// @Description Check the status of an application using application number and phone
// @Tags Admission Applications
// @Accept json
// @Produce json
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param request body StatusCheckRequest true "Status check details"
// @Success 200 {object} response.Success{data=StatusCheckResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/public/applications/status [post]
func (h *ApplicationHandler) CheckStatus(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	var req StatusCheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	checkReq := admissionservice.StatusCheckRequest{
		ApplicationNumber: req.ApplicationNumber,
		Phone:             req.Phone,
	}

	status, err := h.applicationService.CheckStatus(c.Request.Context(), tenantID, checkReq)
	if err != nil {
		switch err {
		case admissionservice.ErrApplicationNotFound:
			apperrors.Abort(c, apperrors.NotFound("Application not found or phone number does not match"))
		case admissionservice.ErrInvalidPhoneNumber:
			apperrors.Abort(c, apperrors.BadRequest("Invalid phone number format"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to check status"))
		}
		return
	}

	resp := StatusCheckResponse{
		ApplicationNumber: status.ApplicationNumber,
		StudentName:       status.StudentName,
		ClassName:         status.ClassApplying,
		Status:            string(status.Status),
		SubmittedAt:       status.SubmittedAt,
		ReviewNotes:       status.ReviewNotes,
	}

	response.OK(c, resp)
}
