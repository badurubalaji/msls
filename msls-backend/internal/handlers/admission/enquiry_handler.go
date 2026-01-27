// Package admission provides HTTP handlers for admission management endpoints.
package admission

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	apperrors "msls-backend/internal/pkg/errors"
	"msls-backend/internal/pkg/response"
	"msls-backend/internal/middleware"
	"msls-backend/internal/services/admission"
)

// EnquiryHandler handles enquiry-related HTTP requests.
type EnquiryHandler struct {
	enquiryService *admission.EnquiryService
}

// NewEnquiryHandler creates a new EnquiryHandler instance.
func NewEnquiryHandler(enquiryService *admission.EnquiryService) *EnquiryHandler {
	return &EnquiryHandler{enquiryService: enquiryService}
}

// List returns all enquiries for the tenant with optional filters.
// @Summary List enquiries
// @Description Get all admission enquiries for the current tenant with optional filters
// @Tags Enquiries
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param branchId query string false "Filter by branch ID"
// @Param sessionId query string false "Filter by session ID"
// @Param status query string false "Filter by status (new, contacted, interested, converted, closed)"
// @Param source query string false "Filter by source"
// @Param classApplying query string false "Filter by class"
// @Param search query string false "Search by name, phone, or enquiry number"
// @Param startDate query string false "Filter by start date (YYYY-MM-DD)"
// @Param endDate query string false "Filter by end date (YYYY-MM-DD)"
// @Param assignedTo query string false "Filter by assigned user ID"
// @Param page query int false "Page number (default 1)"
// @Param pageSize query int false "Page size (default 20)"
// @Success 200 {object} response.Success{data=EnquiryListResponse}
// @Failure 401 {object} apperrors.AppError
// @Failure 403 {object} apperrors.AppError
// @Router /api/v1/enquiries [get]
func (h *EnquiryHandler) List(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	var params ListEnquiriesParams
	if err := c.ShouldBindQuery(&params); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	filter := admission.ListEnquiryFilter{
		TenantID:      tenantID,
		ClassApplying: params.ClassApplying,
		Search:        params.Search,
		Page:          params.Page,
		PageSize:      params.PageSize,
	}

	// Parse optional UUID filters
	if params.BranchID != "" {
		branchID := parseUUID(params.BranchID)
		if branchID == nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid branch ID"))
			return
		}
		filter.BranchID = branchID
	}

	if params.SessionID != "" {
		sessionID := parseUUID(params.SessionID)
		if sessionID == nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid session ID"))
			return
		}
		filter.SessionID = sessionID
	}

	if params.AssignedTo != "" {
		assignedTo := parseUUID(params.AssignedTo)
		if assignedTo == nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid assigned to ID"))
			return
		}
		filter.AssignedTo = assignedTo
	}

	// Parse status filter
	if params.Status != "" {
		status := admission.EnquiryStatus(params.Status)
		if !status.IsValid() {
			apperrors.Abort(c, apperrors.BadRequest("Invalid status value"))
			return
		}
		filter.Status = &status
	}

	// Parse source filter
	if params.Source != "" {
		source := admission.EnquirySource(params.Source)
		if !source.IsValid() {
			apperrors.Abort(c, apperrors.BadRequest("Invalid source value"))
			return
		}
		filter.Source = &source
	}

	// Parse date filters
	if params.StartDate != "" {
		startDate := parseDate(params.StartDate)
		if startDate == nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid start date format (use YYYY-MM-DD)"))
			return
		}
		filter.StartDate = startDate
	}

	if params.EndDate != "" {
		endDate := parseDate(params.EndDate)
		if endDate == nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid end date format (use YYYY-MM-DD)"))
			return
		}
		filter.EndDate = endDate
	}

	enquiries, total, err := h.enquiryService.List(c.Request.Context(), filter)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve enquiries"))
		return
	}

	// Apply defaults for pagination
	pageSize := params.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	page := params.Page
	if page <= 0 {
		page = 1
	}

	resp := EnquiryListResponse{
		Enquiries: enquiriesToResponses(enquiries),
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
	}

	response.OK(c, resp)
}

// GetByID returns an enquiry by ID.
// @Summary Get enquiry by ID
// @Description Get an enquiry with full details
// @Tags Enquiries
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Enquiry ID" format(uuid)
// @Success 200 {object} response.Success{data=EnquiryResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/enquiries/{id} [get]
func (h *EnquiryHandler) GetByID(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid enquiry ID"))
		return
	}

	enquiry, err := h.enquiryService.GetByID(c.Request.Context(), tenantID, id)
	if err != nil {
		switch err {
		case admission.ErrEnquiryNotFound:
			apperrors.Abort(c, apperrors.NotFound("Enquiry not found"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to retrieve enquiry"))
		}
		return
	}

	response.OK(c, enquiryToResponse(enquiry))
}

// Create creates a new enquiry.
// @Summary Create enquiry
// @Description Create a new admission enquiry
// @Tags Enquiries
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param request body CreateEnquiryRequest true "Enquiry details"
// @Success 201 {object} response.Success{data=EnquiryResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Router /api/v1/enquiries [post]
func (h *EnquiryHandler) Create(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	var req CreateEnquiryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	createReq := admission.CreateEnquiryRequest{
		TenantID:        tenantID,
		StudentName:     req.StudentName,
		Gender:          req.Gender,
		ClassApplying:   req.ClassApplying,
		ParentName:      req.ParentName,
		ParentPhone:     req.ParentPhone,
		ParentEmail:     req.ParentEmail,
		ReferralDetails: req.ReferralDetails,
		Remarks:         req.Remarks,
		CreatedBy:       &userID,
	}

	// Parse optional fields
	if req.BranchID != nil {
		createReq.BranchID = parseUUID(*req.BranchID)
	}
	if req.SessionID != nil {
		createReq.SessionID = parseUUID(*req.SessionID)
	}
	if req.DateOfBirth != nil {
		createReq.DateOfBirth = parseDate(*req.DateOfBirth)
	}
	if req.FollowUpDate != nil {
		createReq.FollowUpDate = parseDate(*req.FollowUpDate)
	}
	if req.AssignedTo != nil {
		createReq.AssignedTo = parseUUID(*req.AssignedTo)
	}
	if req.Source != "" {
		createReq.Source = admission.EnquirySource(req.Source)
	}

	enquiry, err := h.enquiryService.Create(c.Request.Context(), createReq)
	if err != nil {
		switch err {
		case admission.ErrStudentNameRequired:
			apperrors.Abort(c, apperrors.BadRequest("Student name is required"))
		case admission.ErrClassApplyingRequired:
			apperrors.Abort(c, apperrors.BadRequest("Class applying is required"))
		case admission.ErrParentNameRequired:
			apperrors.Abort(c, apperrors.BadRequest("Parent name is required"))
		case admission.ErrParentPhoneRequired:
			apperrors.Abort(c, apperrors.BadRequest("Parent phone is required"))
		case admission.ErrInvalidEnquirySource:
			apperrors.Abort(c, apperrors.BadRequest("Invalid source value"))
		case admission.ErrInvalidGender:
			apperrors.Abort(c, apperrors.BadRequest("Invalid gender value"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to create enquiry"))
		}
		return
	}

	response.Created(c, enquiryToResponse(enquiry))
}

// Update updates an enquiry.
// @Summary Update enquiry
// @Description Update an existing enquiry
// @Tags Enquiries
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Enquiry ID" format(uuid)
// @Param request body UpdateEnquiryRequest true "Enquiry updates"
// @Success 200 {object} response.Success{data=EnquiryResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/enquiries/{id} [put]
func (h *EnquiryHandler) Update(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid enquiry ID"))
		return
	}

	var req UpdateEnquiryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	updateReq := admission.UpdateEnquiryRequest{
		StudentName:     req.StudentName,
		Gender:          req.Gender,
		ClassApplying:   req.ClassApplying,
		ParentName:      req.ParentName,
		ParentPhone:     req.ParentPhone,
		ParentEmail:     req.ParentEmail,
		ReferralDetails: req.ReferralDetails,
		Remarks:         req.Remarks,
		UpdatedBy:       &userID,
	}

	// Parse optional fields
	if req.BranchID != nil {
		updateReq.BranchID = parseUUID(*req.BranchID)
	}
	if req.SessionID != nil {
		updateReq.SessionID = parseUUID(*req.SessionID)
	}
	if req.DateOfBirth != nil {
		updateReq.DateOfBirth = parseDate(*req.DateOfBirth)
	}
	if req.FollowUpDate != nil {
		updateReq.FollowUpDate = parseDate(*req.FollowUpDate)
	}
	if req.AssignedTo != nil {
		updateReq.AssignedTo = parseUUID(*req.AssignedTo)
	}
	if req.Source != nil {
		source := admission.EnquirySource(*req.Source)
		updateReq.Source = &source
	}
	if req.Status != nil {
		status := admission.EnquiryStatus(*req.Status)
		updateReq.Status = &status
	}

	enquiry, err := h.enquiryService.Update(c.Request.Context(), tenantID, id, updateReq)
	if err != nil {
		switch err {
		case admission.ErrEnquiryNotFound:
			apperrors.Abort(c, apperrors.NotFound("Enquiry not found"))
		case admission.ErrEnquiryClosed:
			apperrors.Abort(c, apperrors.BadRequest("Enquiry is closed and cannot be modified"))
		case admission.ErrInvalidEnquiryStatus:
			apperrors.Abort(c, apperrors.BadRequest("Invalid status value"))
		case admission.ErrInvalidEnquirySource:
			apperrors.Abort(c, apperrors.BadRequest("Invalid source value"))
		case admission.ErrInvalidGender:
			apperrors.Abort(c, apperrors.BadRequest("Invalid gender value"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to update enquiry"))
		}
		return
	}

	response.OK(c, enquiryToResponse(enquiry))
}

// Delete deletes an enquiry.
// @Summary Delete enquiry
// @Description Delete an enquiry (soft delete)
// @Tags Enquiries
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Enquiry ID" format(uuid)
// @Success 204 "No Content"
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/enquiries/{id} [delete]
func (h *EnquiryHandler) Delete(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid enquiry ID"))
		return
	}

	err = h.enquiryService.Delete(c.Request.Context(), tenantID, id)
	if err != nil {
		switch err {
		case admission.ErrEnquiryNotFound:
			apperrors.Abort(c, apperrors.NotFound("Enquiry not found"))
		case admission.ErrEnquiryAlreadyConverted:
			apperrors.Abort(c, apperrors.BadRequest("Cannot delete a converted enquiry"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to delete enquiry"))
		}
		return
	}

	response.NoContent(c)
}

// AddFollowUp adds a follow-up to an enquiry.
// @Summary Add follow-up
// @Description Add a follow-up entry to an enquiry
// @Tags Enquiries
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Enquiry ID" format(uuid)
// @Param request body CreateFollowUpRequest true "Follow-up details"
// @Success 201 {object} response.Success{data=FollowUpResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/enquiries/{id}/follow-ups [post]
func (h *EnquiryHandler) AddFollowUp(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	idParam := c.Param("id")
	enquiryID, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid enquiry ID"))
		return
	}

	var req CreateFollowUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	// Parse follow-up date
	followUpDate := parseDate(req.FollowUpDate)
	if followUpDate == nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid follow-up date format (use YYYY-MM-DD)"))
		return
	}

	createReq := admission.CreateFollowUpRequest{
		TenantID:     tenantID,
		EnquiryID:    enquiryID,
		FollowUpDate: *followUpDate,
		Notes:        req.Notes,
		CreatedBy:    &userID,
	}

	// Parse optional fields
	if req.ContactMode != "" {
		createReq.ContactMode = admission.ContactMode(req.ContactMode)
	}
	if req.Outcome != nil {
		outcome := admission.FollowUpOutcome(*req.Outcome)
		createReq.Outcome = &outcome
	}
	if req.NextFollowUp != nil {
		createReq.NextFollowUp = parseDate(*req.NextFollowUp)
	}

	followUp, err := h.enquiryService.AddFollowUp(c.Request.Context(), createReq)
	if err != nil {
		switch err {
		case admission.ErrEnquiryNotFound:
			apperrors.Abort(c, apperrors.NotFound("Enquiry not found"))
		case admission.ErrEnquiryClosed:
			apperrors.Abort(c, apperrors.BadRequest("Enquiry is closed and cannot have follow-ups added"))
		case admission.ErrFollowUpDateRequired:
			apperrors.Abort(c, apperrors.BadRequest("Follow-up date is required"))
		case admission.ErrInvalidContactMode:
			apperrors.Abort(c, apperrors.BadRequest("Invalid contact mode"))
		case admission.ErrInvalidFollowUpOutcome:
			apperrors.Abort(c, apperrors.BadRequest("Invalid follow-up outcome"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to add follow-up"))
		}
		return
	}

	response.Created(c, followUpToResponse(followUp))
}

// ListFollowUps returns all follow-ups for an enquiry.
// @Summary List follow-ups
// @Description Get all follow-ups for an enquiry
// @Tags Enquiries
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Enquiry ID" format(uuid)
// @Success 200 {object} response.Success{data=FollowUpListResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/enquiries/{id}/follow-ups [get]
func (h *EnquiryHandler) ListFollowUps(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	idParam := c.Param("id")
	enquiryID, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid enquiry ID"))
		return
	}

	followUps, err := h.enquiryService.ListFollowUps(c.Request.Context(), tenantID, enquiryID)
	if err != nil {
		switch err {
		case admission.ErrEnquiryNotFound:
			apperrors.Abort(c, apperrors.NotFound("Enquiry not found"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to retrieve follow-ups"))
		}
		return
	}

	resp := FollowUpListResponse{
		FollowUps: followUpsToResponses(followUps),
		Total:     len(followUps),
	}

	response.OK(c, resp)
}

// ConvertToApplication converts an enquiry to an application.
// @Summary Convert to application
// @Description Convert an enquiry to an admission application
// @Tags Enquiries
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Enquiry ID" format(uuid)
// @Param request body ConvertEnquiryRequest true "Conversion details"
// @Success 200 {object} response.Success{data=EnquiryResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Failure 409 {object} apperrors.AppError
// @Router /api/v1/enquiries/{id}/convert [post]
func (h *EnquiryHandler) ConvertToApplication(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	idParam := c.Param("id")
	enquiryID, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid enquiry ID"))
		return
	}

	var req ConvertEnquiryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	// Parse optional session ID
	var sessionID *uuid.UUID
	if req.SessionID != "" {
		parsed, err := uuid.Parse(req.SessionID)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid session ID"))
			return
		}
		sessionID = &parsed
	}

	// Parse optional branch ID
	var branchID *uuid.UUID
	if req.BranchID != "" {
		parsed, err := uuid.Parse(req.BranchID)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid branch ID"))
			return
		}
		branchID = &parsed
	}

	convertReq := admission.ConvertEnquiryRequest{
		TenantID:    tenantID,
		EnquiryID:   enquiryID,
		SessionID:   sessionID,
		BranchID:    branchID,
		ConvertedBy: &userID,
	}

	enquiry, err := h.enquiryService.ConvertToApplication(c.Request.Context(), convertReq)
	if err != nil {
		switch err {
		case admission.ErrEnquiryNotFound:
			apperrors.Abort(c, apperrors.NotFound("Enquiry not found"))
		case admission.ErrEnquiryAlreadyConverted:
			apperrors.Abort(c, apperrors.Conflict("Enquiry has already been converted"))
		case admission.ErrEnquiryClosed:
			apperrors.Abort(c, apperrors.BadRequest("Enquiry is closed and cannot be converted"))
		case admission.ErrSessionRequired:
			apperrors.Abort(c, apperrors.BadRequest("Session ID is required - enquiry has no session assigned"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to convert enquiry"))
		}
		return
	}

	// Return the application ID as expected by frontend
	var applicationID string
	if enquiry.ConvertedApplicationID != nil {
		applicationID = enquiry.ConvertedApplicationID.String()
	}

	response.OK(c, gin.H{
		"applicationId": applicationID,
		"enquiry":       enquiryToResponse(enquiry),
	})
}
