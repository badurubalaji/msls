// Package guardian provides guardian and emergency contact management functionality.
package guardian

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	apperrors "msls-backend/internal/pkg/errors"
	"msls-backend/internal/pkg/logger"
	"msls-backend/internal/pkg/response"
	"msls-backend/internal/middleware"
	"msls-backend/internal/pkg/database/models"
)

// Handler handles guardian-related HTTP requests.
type Handler struct {
	service *Service
}

// NewHandler creates a new guardian handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// =========================================================================
// Guardian Request Types
// =========================================================================

// CreateGuardianRequest represents the request body for creating a guardian.
type CreateGuardianRequest struct {
	Relation        string `json:"relation" binding:"required,oneof=father mother grandfather grandmother uncle aunt sibling guardian other"`
	FirstName       string `json:"firstName" binding:"required,max=100"`
	LastName        string `json:"lastName" binding:"required,max=100"`
	Phone           string `json:"phone" binding:"required,max=15"`
	Email           string `json:"email" binding:"omitempty,email,max=255"`
	Occupation      string `json:"occupation" binding:"max=100"`
	AnnualIncome    string `json:"annualIncome"`
	Education       string `json:"education" binding:"max=100"`
	IsPrimary       bool   `json:"isPrimary"`
	HasPortalAccess bool   `json:"hasPortalAccess"`
	AddressLine1    string `json:"addressLine1" binding:"max=255"`
	AddressLine2    string `json:"addressLine2" binding:"max=255"`
	City            string `json:"city" binding:"max=100"`
	State           string `json:"state" binding:"max=100"`
	PostalCode      string `json:"postalCode" binding:"max=10"`
	Country         string `json:"country" binding:"max=100"`
}

// UpdateGuardianRequest represents the request body for updating a guardian.
type UpdateGuardianRequest struct {
	Relation        *string `json:"relation" binding:"omitempty,oneof=father mother grandfather grandmother uncle aunt sibling guardian other"`
	FirstName       *string `json:"firstName" binding:"omitempty,max=100"`
	LastName        *string `json:"lastName" binding:"omitempty,max=100"`
	Phone           *string `json:"phone" binding:"omitempty,max=15"`
	Email           *string `json:"email" binding:"omitempty,email,max=255"`
	Occupation      *string `json:"occupation" binding:"omitempty,max=100"`
	AnnualIncome    *string `json:"annualIncome"`
	Education       *string `json:"education" binding:"omitempty,max=100"`
	IsPrimary       *bool   `json:"isPrimary"`
	HasPortalAccess *bool   `json:"hasPortalAccess"`
	AddressLine1    *string `json:"addressLine1" binding:"omitempty,max=255"`
	AddressLine2    *string `json:"addressLine2" binding:"omitempty,max=255"`
	City            *string `json:"city" binding:"omitempty,max=100"`
	State           *string `json:"state" binding:"omitempty,max=100"`
	PostalCode      *string `json:"postalCode" binding:"omitempty,max=10"`
	Country         *string `json:"country" binding:"omitempty,max=100"`
}

// =========================================================================
// Emergency Contact Request Types
// =========================================================================

// CreateEmergencyContactRequest represents the request body for creating an emergency contact.
type CreateEmergencyContactRequest struct {
	Name           string `json:"name" binding:"required,max=200"`
	Relation       string `json:"relation" binding:"required,max=50"`
	Phone          string `json:"phone" binding:"required,max=15"`
	AlternatePhone string `json:"alternatePhone" binding:"max=15"`
	Priority       int    `json:"priority" binding:"omitempty,min=1,max=5"`
	Notes          string `json:"notes" binding:"max=500"`
}

// UpdateEmergencyContactRequest represents the request body for updating an emergency contact.
type UpdateEmergencyContactRequest struct {
	Name           *string `json:"name" binding:"omitempty,max=200"`
	Relation       *string `json:"relation" binding:"omitempty,max=50"`
	Phone          *string `json:"phone" binding:"omitempty,max=15"`
	AlternatePhone *string `json:"alternatePhone" binding:"omitempty,max=15"`
	Priority       *int    `json:"priority" binding:"omitempty,min=1,max=5"`
	Notes          *string `json:"notes" binding:"omitempty,max=500"`
}

// =========================================================================
// Guardian Handlers
// =========================================================================

// ListGuardians returns all guardians for a student.
// @Summary List guardians
// @Description Get all guardians for a student
// @Tags Guardians
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param studentId path string true "Student ID" format(uuid)
// @Success 200 {object} response.Success{data=GuardianListResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/students/{studentId}/guardians [get]
func (h *Handler) ListGuardians(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	studentIDParam := c.Param("id")
	studentID, err := uuid.Parse(studentIDParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid student ID"))
		return
	}

	guardians, err := h.service.GetGuardiansByStudentID(c.Request.Context(), tenantID, studentID)
	if err != nil {
		if errors.Is(err, ErrStudentNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Student not found"))
			return
		}
		logger.Error("Failed to list guardians",
			zap.String("student_id", studentID.String()),
			zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve guardians"))
		return
	}

	resp := GuardianListResponse{
		Guardians: ToGuardianResponses(guardians),
		Total:     len(guardians),
	}

	response.OK(c, resp)
}

// GetGuardian returns a guardian by ID.
// @Summary Get guardian by ID
// @Description Get a guardian with full details
// @Tags Guardians
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param studentId path string true "Student ID" format(uuid)
// @Param id path string true "Guardian ID" format(uuid)
// @Success 200 {object} response.Success{data=GuardianResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/students/{studentId}/guardians/{id} [get]
func (h *Handler) GetGuardian(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	idParam := c.Param("guardianId")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid guardian ID"))
		return
	}

	guardian, err := h.service.GetGuardianByID(c.Request.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, ErrGuardianNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Guardian not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve guardian"))
		return
	}

	response.OK(c, ToGuardianResponse(guardian))
}

// CreateGuardian creates a new guardian for a student.
// @Summary Create guardian
// @Description Create a new guardian for a student
// @Tags Guardians
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param studentId path string true "Student ID" format(uuid)
// @Param request body CreateGuardianRequest true "Guardian details"
// @Success 201 {object} response.Success{data=GuardianResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/students/{studentId}/guardians [post]
func (h *Handler) CreateGuardian(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	studentIDParam := c.Param("id")
	studentID, err := uuid.Parse(studentIDParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid student ID"))
		return
	}

	var req CreateGuardianRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	dto := CreateGuardianDTO{
		TenantID:        tenantID,
		StudentID:       studentID,
		Relation:        models.GuardianRelation(req.Relation),
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Phone:           req.Phone,
		Email:           req.Email,
		Occupation:      req.Occupation,
		Education:       req.Education,
		IsPrimary:       req.IsPrimary,
		HasPortalAccess: req.HasPortalAccess,
		AddressLine1:    req.AddressLine1,
		AddressLine2:    req.AddressLine2,
		City:            req.City,
		State:           req.State,
		PostalCode:      req.PostalCode,
		Country:         req.Country,
		CreatedBy:       &userID,
	}

	// Parse annual income if provided
	if req.AnnualIncome != "" {
		income, err := decimal.NewFromString(req.AnnualIncome)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid annual income format"))
			return
		}
		dto.AnnualIncome = income
	}

	guardian, err := h.service.CreateGuardian(c.Request.Context(), dto)
	if err != nil {
		switch {
		case errors.Is(err, ErrStudentNotFound):
			apperrors.Abort(c, apperrors.NotFound("Student not found"))
		case errors.Is(err, ErrInvalidRelation):
			apperrors.Abort(c, apperrors.BadRequest("Invalid guardian relation"))
		case errors.Is(err, ErrFirstNameRequired):
			apperrors.Abort(c, apperrors.BadRequest("First name is required"))
		case errors.Is(err, ErrLastNameRequired):
			apperrors.Abort(c, apperrors.BadRequest("Last name is required"))
		case errors.Is(err, ErrPhoneRequired):
			apperrors.Abort(c, apperrors.BadRequest("Phone number is required"))
		default:
			logger.Error("Failed to create guardian", zap.Error(err))
			apperrors.Abort(c, apperrors.InternalError("Failed to create guardian"))
		}
		return
	}

	response.Created(c, ToGuardianResponse(guardian))
}

// UpdateGuardian updates an existing guardian.
// @Summary Update guardian
// @Description Update an existing guardian
// @Tags Guardians
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param studentId path string true "Student ID" format(uuid)
// @Param id path string true "Guardian ID" format(uuid)
// @Param request body UpdateGuardianRequest true "Guardian updates"
// @Success 200 {object} response.Success{data=GuardianResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/students/{studentId}/guardians/{id} [put]
func (h *Handler) UpdateGuardian(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	idParam := c.Param("guardianId")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid guardian ID"))
		return
	}

	var req UpdateGuardianRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	dto := UpdateGuardianDTO{
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Phone:           req.Phone,
		Email:           req.Email,
		Occupation:      req.Occupation,
		Education:       req.Education,
		IsPrimary:       req.IsPrimary,
		HasPortalAccess: req.HasPortalAccess,
		AddressLine1:    req.AddressLine1,
		AddressLine2:    req.AddressLine2,
		City:            req.City,
		State:           req.State,
		PostalCode:      req.PostalCode,
		Country:         req.Country,
		UpdatedBy:       &userID,
	}

	// Parse relation if provided
	if req.Relation != nil {
		relation := models.GuardianRelation(*req.Relation)
		dto.Relation = &relation
	}

	// Parse annual income if provided
	if req.AnnualIncome != nil {
		income, err := decimal.NewFromString(*req.AnnualIncome)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid annual income format"))
			return
		}
		dto.AnnualIncome = &income
	}

	guardian, err := h.service.UpdateGuardian(c.Request.Context(), tenantID, id, dto)
	if err != nil {
		switch {
		case errors.Is(err, ErrGuardianNotFound):
			apperrors.Abort(c, apperrors.NotFound("Guardian not found"))
		case errors.Is(err, ErrInvalidRelation):
			apperrors.Abort(c, apperrors.BadRequest("Invalid guardian relation"))
		default:
			logger.Error("Failed to update guardian", zap.Error(err))
			apperrors.Abort(c, apperrors.InternalError("Failed to update guardian"))
		}
		return
	}

	response.OK(c, ToGuardianResponse(guardian))
}

// DeleteGuardian deletes a guardian.
// @Summary Delete guardian
// @Description Delete a guardian
// @Tags Guardians
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param studentId path string true "Student ID" format(uuid)
// @Param id path string true "Guardian ID" format(uuid)
// @Success 204 "No Content"
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/students/{studentId}/guardians/{id} [delete]
func (h *Handler) DeleteGuardian(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	idParam := c.Param("guardianId")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid guardian ID"))
		return
	}

	if err := h.service.DeleteGuardian(c.Request.Context(), tenantID, id); err != nil {
		if errors.Is(err, ErrGuardianNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Guardian not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to delete guardian"))
		return
	}

	response.NoContent(c)
}

// SetPrimaryGuardian sets a guardian as the primary guardian.
// @Summary Set primary guardian
// @Description Set a guardian as the primary guardian for a student
// @Tags Guardians
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param studentId path string true "Student ID" format(uuid)
// @Param id path string true "Guardian ID" format(uuid)
// @Success 200 {object} response.Success{data=GuardianResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/students/{studentId}/guardians/{id}/set-primary [post]
func (h *Handler) SetPrimaryGuardian(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	idParam := c.Param("guardianId")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid guardian ID"))
		return
	}

	guardian, err := h.service.SetPrimaryGuardian(c.Request.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, ErrGuardianNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Guardian not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to set primary guardian"))
		return
	}

	response.OK(c, ToGuardianResponse(guardian))
}

// =========================================================================
// Emergency Contact Handlers
// =========================================================================

// ListEmergencyContacts returns all emergency contacts for a student.
// @Summary List emergency contacts
// @Description Get all emergency contacts for a student
// @Tags Emergency Contacts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param studentId path string true "Student ID" format(uuid)
// @Success 200 {object} response.Success{data=EmergencyContactListResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/students/{studentId}/emergency-contacts [get]
func (h *Handler) ListEmergencyContacts(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	studentIDParam := c.Param("id")
	studentID, err := uuid.Parse(studentIDParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid student ID"))
		return
	}

	contacts, err := h.service.GetEmergencyContactsByStudentID(c.Request.Context(), tenantID, studentID)
	if err != nil {
		if errors.Is(err, ErrStudentNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Student not found"))
			return
		}
		logger.Error("Failed to list emergency contacts",
			zap.String("student_id", studentID.String()),
			zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve emergency contacts"))
		return
	}

	resp := EmergencyContactListResponse{
		Contacts: ToEmergencyContactResponses(contacts),
		Total:    len(contacts),
	}

	response.OK(c, resp)
}

// GetEmergencyContact returns an emergency contact by ID.
// @Summary Get emergency contact by ID
// @Description Get an emergency contact with full details
// @Tags Emergency Contacts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param studentId path string true "Student ID" format(uuid)
// @Param id path string true "Emergency Contact ID" format(uuid)
// @Success 200 {object} response.Success{data=EmergencyContactResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/students/{studentId}/emergency-contacts/{id} [get]
func (h *Handler) GetEmergencyContact(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	idParam := c.Param("contactId")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid emergency contact ID"))
		return
	}

	contact, err := h.service.GetEmergencyContactByID(c.Request.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, ErrEmergencyContactNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Emergency contact not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve emergency contact"))
		return
	}

	response.OK(c, ToEmergencyContactResponse(contact))
}

// CreateEmergencyContact creates a new emergency contact for a student.
// @Summary Create emergency contact
// @Description Create a new emergency contact for a student
// @Tags Emergency Contacts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param studentId path string true "Student ID" format(uuid)
// @Param request body CreateEmergencyContactRequest true "Emergency contact details"
// @Success 201 {object} response.Success{data=EmergencyContactResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Failure 409 {object} apperrors.AppError
// @Router /api/v1/students/{studentId}/emergency-contacts [post]
func (h *Handler) CreateEmergencyContact(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	studentIDParam := c.Param("id")
	studentID, err := uuid.Parse(studentIDParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid student ID"))
		return
	}

	var req CreateEmergencyContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	dto := CreateEmergencyContactDTO{
		TenantID:       tenantID,
		StudentID:      studentID,
		Name:           req.Name,
		Relation:       req.Relation,
		Phone:          req.Phone,
		AlternatePhone: req.AlternatePhone,
		Priority:       req.Priority,
		Notes:          req.Notes,
		CreatedBy:      &userID,
	}

	contact, err := h.service.CreateEmergencyContact(c.Request.Context(), dto)
	if err != nil {
		switch {
		case errors.Is(err, ErrStudentNotFound):
			apperrors.Abort(c, apperrors.NotFound("Student not found"))
		case errors.Is(err, ErrPriorityConflict):
			apperrors.Abort(c, apperrors.Conflict("Priority already exists for this student"))
		case errors.Is(err, ErrInvalidPriority):
			apperrors.Abort(c, apperrors.BadRequest("Priority must be between 1 and 5"))
		case errors.Is(err, ErrNameRequired):
			apperrors.Abort(c, apperrors.BadRequest("Name is required"))
		case errors.Is(err, ErrRelationRequired):
			apperrors.Abort(c, apperrors.BadRequest("Relation is required"))
		case errors.Is(err, ErrPhoneRequired):
			apperrors.Abort(c, apperrors.BadRequest("Phone number is required"))
		default:
			logger.Error("Failed to create emergency contact", zap.Error(err))
			apperrors.Abort(c, apperrors.InternalError("Failed to create emergency contact"))
		}
		return
	}

	response.Created(c, ToEmergencyContactResponse(contact))
}

// UpdateEmergencyContact updates an existing emergency contact.
// @Summary Update emergency contact
// @Description Update an existing emergency contact
// @Tags Emergency Contacts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param studentId path string true "Student ID" format(uuid)
// @Param id path string true "Emergency Contact ID" format(uuid)
// @Param request body UpdateEmergencyContactRequest true "Emergency contact updates"
// @Success 200 {object} response.Success{data=EmergencyContactResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Failure 409 {object} apperrors.AppError
// @Router /api/v1/students/{studentId}/emergency-contacts/{id} [put]
func (h *Handler) UpdateEmergencyContact(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	idParam := c.Param("contactId")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid emergency contact ID"))
		return
	}

	var req UpdateEmergencyContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	dto := UpdateEmergencyContactDTO{
		Name:           req.Name,
		Relation:       req.Relation,
		Phone:          req.Phone,
		AlternatePhone: req.AlternatePhone,
		Priority:       req.Priority,
		Notes:          req.Notes,
		UpdatedBy:      &userID,
	}

	contact, err := h.service.UpdateEmergencyContact(c.Request.Context(), tenantID, id, dto)
	if err != nil {
		switch {
		case errors.Is(err, ErrEmergencyContactNotFound):
			apperrors.Abort(c, apperrors.NotFound("Emergency contact not found"))
		case errors.Is(err, ErrPriorityConflict):
			apperrors.Abort(c, apperrors.Conflict("Priority already exists for this student"))
		case errors.Is(err, ErrInvalidPriority):
			apperrors.Abort(c, apperrors.BadRequest("Priority must be between 1 and 5"))
		default:
			logger.Error("Failed to update emergency contact", zap.Error(err))
			apperrors.Abort(c, apperrors.InternalError("Failed to update emergency contact"))
		}
		return
	}

	response.OK(c, ToEmergencyContactResponse(contact))
}

// DeleteEmergencyContact deletes an emergency contact.
// @Summary Delete emergency contact
// @Description Delete an emergency contact
// @Tags Emergency Contacts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param studentId path string true "Student ID" format(uuid)
// @Param id path string true "Emergency Contact ID" format(uuid)
// @Success 204 "No Content"
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/students/{studentId}/emergency-contacts/{id} [delete]
func (h *Handler) DeleteEmergencyContact(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	idParam := c.Param("contactId")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid emergency contact ID"))
		return
	}

	if err := h.service.DeleteEmergencyContact(c.Request.Context(), tenantID, id); err != nil {
		if errors.Is(err, ErrEmergencyContactNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Emergency contact not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to delete emergency contact"))
		return
	}

	response.NoContent(c)
}
