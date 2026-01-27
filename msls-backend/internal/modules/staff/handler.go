// Package staff provides staff management functionality.
package staff

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	apperrors "msls-backend/internal/pkg/errors"
	"msls-backend/internal/pkg/logger"
	"msls-backend/internal/pkg/response"
	"msls-backend/internal/middleware"
	"msls-backend/internal/pkg/database/models"
)

// Handler handles staff-related HTTP requests.
type Handler struct {
	service *Service
}

// NewHandler creates a new staff handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateStaffRequest represents the request body for creating a staff member.
type CreateStaffRequest struct {
	BranchID      string `json:"branchId" binding:"required,uuid"`
	FirstName     string `json:"firstName" binding:"required,max=100"`
	MiddleName    string `json:"middleName" binding:"max=100"`
	LastName      string `json:"lastName" binding:"required,max=100"`
	DateOfBirth   string `json:"dateOfBirth" binding:"required"` // Format: YYYY-MM-DD
	Gender        string `json:"gender" binding:"required,oneof=male female other"`
	BloodGroup    string `json:"bloodGroup" binding:"max=10"`
	Nationality   string `json:"nationality" binding:"max=50"`
	Religion      string `json:"religion" binding:"max=50"`
	MaritalStatus string `json:"maritalStatus" binding:"max=20"`
	
	PersonalEmail            string `json:"personalEmail" binding:"omitempty,email,max=255"`
	WorkEmail                string `json:"workEmail" binding:"required,email,max=255"`
	PersonalPhone            string `json:"personalPhone" binding:"max=20"`
	WorkPhone                string `json:"workPhone" binding:"required,max=20"`
	EmergencyContactName     string `json:"emergencyContactName" binding:"max=200"`
	EmergencyContactPhone    string `json:"emergencyContactPhone" binding:"max=20"`
	EmergencyContactRelation string `json:"emergencyContactRelation" binding:"max=50"`
	
	CurrentAddress   *AddressRequest `json:"currentAddress"`
	PermanentAddress *AddressRequest `json:"permanentAddress"`
	SameAsCurrent    bool            `json:"sameAsCurrent"`
	
	StaffType          string  `json:"staffType" binding:"required,oneof=teaching non_teaching"`
	DepartmentID       *string `json:"departmentId" binding:"omitempty,uuid"`
	DesignationID      *string `json:"designationId" binding:"omitempty,uuid"`
	ReportingManagerID *string `json:"reportingManagerId" binding:"omitempty,uuid"`
	JoinDate           string  `json:"joinDate" binding:"required"` // Format: YYYY-MM-DD
	ConfirmationDate   *string `json:"confirmationDate"`
	ProbationEndDate   *string `json:"probationEndDate"`
	
	Bio string `json:"bio" binding:"max=1000"`
}

// UpdateStaffRequest represents the request body for updating a staff member.
type UpdateStaffRequest struct {
	FirstName     *string `json:"firstName" binding:"omitempty,max=100"`
	MiddleName    *string `json:"middleName" binding:"omitempty,max=100"`
	LastName      *string `json:"lastName" binding:"omitempty,max=100"`
	DateOfBirth   *string `json:"dateOfBirth"`
	Gender        *string `json:"gender" binding:"omitempty,oneof=male female other"`
	BloodGroup    *string `json:"bloodGroup" binding:"omitempty,max=10"`
	Nationality   *string `json:"nationality" binding:"omitempty,max=50"`
	Religion      *string `json:"religion" binding:"omitempty,max=50"`
	MaritalStatus *string `json:"maritalStatus" binding:"omitempty,max=20"`
	
	PersonalEmail            *string `json:"personalEmail" binding:"omitempty,email,max=255"`
	WorkEmail                *string `json:"workEmail" binding:"omitempty,email,max=255"`
	PersonalPhone            *string `json:"personalPhone" binding:"omitempty,max=20"`
	WorkPhone                *string `json:"workPhone" binding:"omitempty,max=20"`
	EmergencyContactName     *string `json:"emergencyContactName" binding:"omitempty,max=200"`
	EmergencyContactPhone    *string `json:"emergencyContactPhone" binding:"omitempty,max=20"`
	EmergencyContactRelation *string `json:"emergencyContactRelation" binding:"omitempty,max=50"`
	
	CurrentAddress   *AddressRequest `json:"currentAddress"`
	PermanentAddress *AddressRequest `json:"permanentAddress"`
	SameAsCurrent    *bool           `json:"sameAsCurrent"`
	
	StaffType          *string `json:"staffType" binding:"omitempty,oneof=teaching non_teaching"`
	DepartmentID       *string `json:"departmentId" binding:"omitempty,uuid"`
	DesignationID      *string `json:"designationId" binding:"omitempty,uuid"`
	ReportingManagerID *string `json:"reportingManagerId" binding:"omitempty,uuid"`
	ConfirmationDate   *string `json:"confirmationDate"`
	ProbationEndDate   *string `json:"probationEndDate"`
	
	Bio     *string `json:"bio" binding:"omitempty,max=1000"`
	Version int     `json:"version"`
}

// StatusUpdateRequest represents the request body for updating staff status.
type StatusUpdateRequest struct {
	Status        string `json:"status" binding:"required,oneof=active inactive terminated on_leave"`
	Reason        string `json:"reason" binding:"required,max=500"`
	EffectiveDate string `json:"effectiveDate" binding:"required"` // Format: YYYY-MM-DD
}

// AddressRequest represents an address in requests.
type AddressRequest struct {
	AddressLine1 string `json:"addressLine1" binding:"max=255"`
	AddressLine2 string `json:"addressLine2" binding:"max=255"`
	City         string `json:"city" binding:"max=100"`
	State        string `json:"state" binding:"max=100"`
	Pincode      string `json:"pincode" binding:"max=10"`
	Country      string `json:"country" binding:"max=100"`
}

// List returns all staff for the tenant.
// @Summary List staff
// @Description Get all staff for the current tenant with pagination and filters
// @Tags Staff
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param search query string false "Search by name, employee ID, or phone"
// @Param branch_id query string false "Filter by branch ID"
// @Param department_id query string false "Filter by department ID"
// @Param designation_id query string false "Filter by designation ID"
// @Param staff_type query string false "Filter by staff type (teaching, non_teaching)"
// @Param status query string false "Filter by status (active, inactive, terminated, on_leave)"
// @Param gender query string false "Filter by gender (male, female, other)"
// @Param cursor query string false "Cursor for pagination"
// @Param limit query int false "Number of results per page (max 100)"
// @Param sort_by query string false "Sort by field (name, employee_id, join_date)"
// @Param sort_order query string false "Sort order (asc, desc)"
// @Success 200 {object} response.Success{data=StaffListResponse}
// @Router /api/v1/staff [get]
func (h *Handler) List(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	filter := ListFilter{
		TenantID:  tenantID,
		Search:    c.Query("search"),
		Cursor:    c.Query("cursor"),
		SortBy:    c.Query("sort_by"),
		SortOrder: c.Query("sort_order"),
	}

	// Parse branch_id filter
	if branchIDStr := c.Query("branch_id"); branchIDStr != "" {
		branchID, err := uuid.Parse(branchIDStr)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid branch ID"))
			return
		}
		filter.BranchID = &branchID
	}

	// Parse department_id filter
	if deptIDStr := c.Query("department_id"); deptIDStr != "" {
		deptID, err := uuid.Parse(deptIDStr)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid department ID"))
			return
		}
		filter.DepartmentID = &deptID
	}

	// Parse designation_id filter
	if desigIDStr := c.Query("designation_id"); desigIDStr != "" {
		desigID, err := uuid.Parse(desigIDStr)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid designation ID"))
			return
		}
		filter.DesignationID = &desigID
	}

	// Parse staff_type filter
	if staffTypeStr := c.Query("staff_type"); staffTypeStr != "" {
		staffType := models.StaffType(staffTypeStr)
		if !staffType.IsValid() {
			apperrors.Abort(c, apperrors.BadRequest("Invalid staff type value"))
			return
		}
		filter.StaffType = &staffType
	}

	// Parse status filter
	if statusStr := c.Query("status"); statusStr != "" {
		status := models.StaffStatus(statusStr)
		if !status.IsValid() {
			apperrors.Abort(c, apperrors.BadRequest("Invalid status value"))
			return
		}
		filter.Status = &status
	}

	// Parse gender filter
	if genderStr := c.Query("gender"); genderStr != "" {
		gender := models.Gender(genderStr)
		if !gender.IsValid() {
			apperrors.Abort(c, apperrors.BadRequest("Invalid gender value"))
			return
		}
		filter.Gender = &gender
	}

	// Parse limit
	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err == nil && limit > 0 {
			filter.Limit = limit
		}
	}

	staffList, nextCursor, total, err := h.service.List(c.Request.Context(), filter)
	if err != nil {
		logger.Error("Failed to list staff",
			zap.String("tenant_id", tenantID.String()),
			zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve staff"))
		return
	}

	resp := StaffListResponse{
		Staff:      ToStaffResponses(staffList),
		NextCursor: nextCursor,
		HasMore:    nextCursor != "",
		Total:      total,
	}

	response.OK(c, resp)
}

// Get retrieves a staff member by ID.
// @Summary Get staff by ID
// @Description Get a staff member by their ID
// @Tags Staff
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Staff ID"
// @Success 200 {object} response.Success{data=StaffResponse}
// @Router /api/v1/staff/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid staff ID"))
		return
	}

	staff, err := h.service.GetByID(c.Request.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, ErrStaffNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Staff member not found"))
			return
		}
		logger.Error("Failed to get staff",
			zap.String("tenant_id", tenantID.String()),
			zap.String("staff_id", id.String()),
			zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve staff member"))
		return
	}

	response.OK(c, ToStaffResponse(staff))
}

// Create creates a new staff member.
// @Summary Create staff
// @Description Create a new staff member
// @Tags Staff
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param request body CreateStaffRequest true "Staff details"
// @Success 201 {object} response.Success{data=StaffResponse}
// @Router /api/v1/staff [post]
func (h *Handler) Create(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	var req CreateStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	branchID, err := uuid.Parse(req.BranchID)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid branch ID"))
		return
	}

	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid date of birth format, use YYYY-MM-DD"))
		return
	}

	joinDate, err := time.Parse("2006-01-02", req.JoinDate)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid join date format, use YYYY-MM-DD"))
		return
	}

	dto := CreateStaffDTO{
		TenantID:      tenantID,
		BranchID:      branchID,
		FirstName:     req.FirstName,
		MiddleName:    req.MiddleName,
		LastName:      req.LastName,
		DateOfBirth:   dob,
		Gender:        models.Gender(req.Gender),
		BloodGroup:    req.BloodGroup,
		Nationality:   req.Nationality,
		Religion:      req.Religion,
		MaritalStatus: req.MaritalStatus,
		PersonalEmail: req.PersonalEmail,
		WorkEmail:     req.WorkEmail,
		PersonalPhone: req.PersonalPhone,
		WorkPhone:     req.WorkPhone,
		EmergencyContactName:     req.EmergencyContactName,
		EmergencyContactPhone:    req.EmergencyContactPhone,
		EmergencyContactRelation: req.EmergencyContactRelation,
		StaffType:     models.StaffType(req.StaffType),
		JoinDate:      joinDate,
		Bio:           req.Bio,
		SameAsCurrent: req.SameAsCurrent,
		CreatedBy:     &userID,
	}

	// Parse optional UUIDs
	if req.DepartmentID != nil {
		deptID, err := uuid.Parse(*req.DepartmentID)
		if err == nil {
			dto.DepartmentID = &deptID
		}
	}
	if req.DesignationID != nil {
		desigID, err := uuid.Parse(*req.DesignationID)
		if err == nil {
			dto.DesignationID = &desigID
		}
	}
	if req.ReportingManagerID != nil {
		mgrID, err := uuid.Parse(*req.ReportingManagerID)
		if err == nil {
			dto.ReportingManagerID = &mgrID
		}
	}

	// Parse optional dates
	if req.ConfirmationDate != nil {
		confDate, err := time.Parse("2006-01-02", *req.ConfirmationDate)
		if err == nil {
			dto.ConfirmationDate = &confDate
		}
	}
	if req.ProbationEndDate != nil {
		probDate, err := time.Parse("2006-01-02", *req.ProbationEndDate)
		if err == nil {
			dto.ProbationEndDate = &probDate
		}
	}

	// Convert addresses
	if req.CurrentAddress != nil {
		dto.CurrentAddress = &AddressDTO{
			AddressLine1: req.CurrentAddress.AddressLine1,
			AddressLine2: req.CurrentAddress.AddressLine2,
			City:         req.CurrentAddress.City,
			State:        req.CurrentAddress.State,
			Pincode:      req.CurrentAddress.Pincode,
			Country:      req.CurrentAddress.Country,
		}
	}
	if req.PermanentAddress != nil {
		dto.PermanentAddress = &AddressDTO{
			AddressLine1: req.PermanentAddress.AddressLine1,
			AddressLine2: req.PermanentAddress.AddressLine2,
			City:         req.PermanentAddress.City,
			State:        req.PermanentAddress.State,
			Pincode:      req.PermanentAddress.Pincode,
			Country:      req.PermanentAddress.Country,
		}
	}

	staff, err := h.service.Create(c.Request.Context(), dto)
	if err != nil {
		switch {
		case errors.Is(err, ErrBranchNotFound):
			apperrors.Abort(c, apperrors.BadRequest("Branch not found"))
		case errors.Is(err, ErrDuplicateEmployeeID):
			apperrors.Abort(c, apperrors.Conflict("Employee ID already exists"))
		case errors.Is(err, ErrInvalidDateOfBirth):
			apperrors.Abort(c, apperrors.BadRequest("Date of birth cannot be in the future"))
		case errors.Is(err, ErrReportingManagerNotFound):
			apperrors.Abort(c, apperrors.BadRequest("Reporting manager not found"))
		default:
			logger.Error("Failed to create staff",
				zap.String("tenant_id", tenantID.String()),
				zap.Error(err))
			apperrors.Abort(c, apperrors.InternalError("Failed to create staff member"))
		}
		return
	}

	response.Created(c, ToStaffResponse(staff))
}

// Update updates a staff member.
// @Summary Update staff
// @Description Update an existing staff member
// @Tags Staff
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Staff ID"
// @Param request body UpdateStaffRequest true "Staff details"
// @Success 200 {object} response.Success{data=StaffResponse}
// @Router /api/v1/staff/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid staff ID"))
		return
	}

	var req UpdateStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	dto := UpdateStaffDTO{
		FirstName:     req.FirstName,
		MiddleName:    req.MiddleName,
		LastName:      req.LastName,
		BloodGroup:    req.BloodGroup,
		Nationality:   req.Nationality,
		Religion:      req.Religion,
		MaritalStatus: req.MaritalStatus,
		PersonalEmail: req.PersonalEmail,
		WorkEmail:     req.WorkEmail,
		PersonalPhone: req.PersonalPhone,
		WorkPhone:     req.WorkPhone,
		EmergencyContactName:     req.EmergencyContactName,
		EmergencyContactPhone:    req.EmergencyContactPhone,
		EmergencyContactRelation: req.EmergencyContactRelation,
		Bio:           req.Bio,
		SameAsCurrent: req.SameAsCurrent,
		Version:       req.Version,
		UpdatedBy:     &userID,
	}

	// Parse optional fields
	if req.DateOfBirth != nil {
		dob, err := time.Parse("2006-01-02", *req.DateOfBirth)
		if err == nil {
			dto.DateOfBirth = &dob
		}
	}
	if req.Gender != nil {
		gender := models.Gender(*req.Gender)
		dto.Gender = &gender
	}
	if req.StaffType != nil {
		staffType := models.StaffType(*req.StaffType)
		dto.StaffType = &staffType
	}
	if req.DepartmentID != nil {
		deptID, err := uuid.Parse(*req.DepartmentID)
		if err == nil {
			dto.DepartmentID = &deptID
		}
	}
	if req.DesignationID != nil {
		desigID, err := uuid.Parse(*req.DesignationID)
		if err == nil {
			dto.DesignationID = &desigID
		}
	}
	if req.ReportingManagerID != nil {
		mgrID, err := uuid.Parse(*req.ReportingManagerID)
		if err == nil {
			dto.ReportingManagerID = &mgrID
		}
	}
	if req.ConfirmationDate != nil {
		confDate, err := time.Parse("2006-01-02", *req.ConfirmationDate)
		if err == nil {
			dto.ConfirmationDate = &confDate
		}
	}
	if req.ProbationEndDate != nil {
		probDate, err := time.Parse("2006-01-02", *req.ProbationEndDate)
		if err == nil {
			dto.ProbationEndDate = &probDate
		}
	}

	// Convert addresses
	if req.CurrentAddress != nil {
		dto.CurrentAddress = &AddressDTO{
			AddressLine1: req.CurrentAddress.AddressLine1,
			AddressLine2: req.CurrentAddress.AddressLine2,
			City:         req.CurrentAddress.City,
			State:        req.CurrentAddress.State,
			Pincode:      req.CurrentAddress.Pincode,
			Country:      req.CurrentAddress.Country,
		}
	}
	if req.PermanentAddress != nil {
		dto.PermanentAddress = &AddressDTO{
			AddressLine1: req.PermanentAddress.AddressLine1,
			AddressLine2: req.PermanentAddress.AddressLine2,
			City:         req.PermanentAddress.City,
			State:        req.PermanentAddress.State,
			Pincode:      req.PermanentAddress.Pincode,
			Country:      req.PermanentAddress.Country,
		}
	}

	staff, err := h.service.Update(c.Request.Context(), tenantID, id, dto)
	if err != nil {
		switch {
		case errors.Is(err, ErrStaffNotFound):
			apperrors.Abort(c, apperrors.NotFound("Staff member not found"))
		case errors.Is(err, ErrOptimisticLockConflict):
			apperrors.Abort(c, apperrors.Conflict("Staff was modified by another user, please refresh and try again"))
		case errors.Is(err, ErrReportingManagerNotFound):
			apperrors.Abort(c, apperrors.BadRequest("Reporting manager not found"))
		default:
			logger.Error("Failed to update staff",
				zap.String("tenant_id", tenantID.String()),
				zap.String("staff_id", id.String()),
				zap.Error(err))
			apperrors.Abort(c, apperrors.InternalError("Failed to update staff member"))
		}
		return
	}

	response.OK(c, ToStaffResponse(staff))
}

// UpdateStatus updates the status of a staff member.
// @Summary Update staff status
// @Description Update the status of a staff member with reason
// @Tags Staff
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Staff ID"
// @Param request body StatusUpdateRequest true "Status update details"
// @Success 200 {object} response.Success{data=StaffResponse}
// @Router /api/v1/staff/{id}/status [patch]
func (h *Handler) UpdateStatus(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid staff ID"))
		return
	}

	var req StatusUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	effectiveDate, err := time.Parse("2006-01-02", req.EffectiveDate)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid effective date format, use YYYY-MM-DD"))
		return
	}

	dto := StatusUpdateDTO{
		Status:        models.StaffStatus(req.Status),
		Reason:        req.Reason,
		EffectiveDate: effectiveDate,
		UpdatedBy:     &userID,
	}

	staff, err := h.service.UpdateStatus(c.Request.Context(), tenantID, id, dto)
	if err != nil {
		switch {
		case errors.Is(err, ErrStaffNotFound):
			apperrors.Abort(c, apperrors.NotFound("Staff member not found"))
		case errors.Is(err, ErrInvalidStatus):
			apperrors.Abort(c, apperrors.BadRequest("Invalid status value"))
		default:
			logger.Error("Failed to update staff status",
				zap.String("tenant_id", tenantID.String()),
				zap.String("staff_id", id.String()),
				zap.Error(err))
			apperrors.Abort(c, apperrors.InternalError("Failed to update staff status"))
		}
		return
	}

	response.OK(c, ToStaffResponse(staff))
}

// GetStatusHistory returns the status history for a staff member.
// @Summary Get staff status history
// @Description Get the status change history for a staff member
// @Tags Staff
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Staff ID"
// @Success 200 {object} response.Success{data=[]StatusHistoryResponse}
// @Router /api/v1/staff/{id}/status-history [get]
func (h *Handler) GetStatusHistory(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid staff ID"))
		return
	}

	histories, err := h.service.GetStatusHistory(c.Request.Context(), tenantID, id)
	if err != nil {
		logger.Error("Failed to get status history",
			zap.String("tenant_id", tenantID.String()),
			zap.String("staff_id", id.String()),
			zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve status history"))
		return
	}

	response.OK(c, ToStatusHistoryResponses(histories))
}

// Delete soft deletes a staff member.
// @Summary Delete staff
// @Description Soft delete a staff member
// @Tags Staff
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Staff ID"
// @Success 204 "No Content"
// @Router /api/v1/staff/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid staff ID"))
		return
	}

	err = h.service.Delete(c.Request.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, ErrStaffNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Staff member not found"))
			return
		}
		logger.Error("Failed to delete staff",
			zap.String("tenant_id", tenantID.String()),
			zap.String("staff_id", id.String()),
			zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to delete staff member"))
		return
	}

	response.NoContent(c)
}

// UpdatePhoto uploads and updates the photo for a staff member.
// @Summary Upload staff photo
// @Description Upload a photo for a staff member
// @Tags Staff
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Staff ID"
// @Param photo formData file true "Photo file"
// @Success 200 {object} response.Success{data=StaffResponse}
// @Router /api/v1/staff/{id}/photo [post]
func (h *Handler) UpdatePhoto(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid staff ID"))
		return
	}

	// Verify staff exists
	_, err = h.service.GetByID(c.Request.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, ErrStaffNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Staff member not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to verify staff member"))
		return
	}

	file, header, err := c.Request.FormFile("photo")
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Photo file is required"))
		return
	}
	defer file.Close()

	// Validate file size (max 2MB)
	const maxFileSize = 2 * 1024 * 1024
	if header.Size > maxFileSize {
		apperrors.Abort(c, apperrors.BadRequest("Photo size must be less than 2MB"))
		return
	}

	// Validate content type
	contentType := header.Header.Get("Content-Type")
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
	}
	if !allowedTypes[contentType] {
		apperrors.Abort(c, apperrors.BadRequest("Invalid photo format. Only JPEG and PNG are allowed."))
		return
	}

	// Save file
	uploadDir := filepath.Join("uploads", "staff", tenantID.String())
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to create upload directory"))
		return
	}

	filename := fmt.Sprintf("%s%s", uuid.New().String(), filepath.Ext(header.Filename))
	filePath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(filePath)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to upload photo"))
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to upload photo"))
		return
	}

	// Build photo URL
	photoURL := fmt.Sprintf("/uploads/staff/%s/%s", tenantID.String(), filename)

	// Update staff record
	staff, err := h.service.UpdatePhoto(c.Request.Context(), tenantID, id, photoURL, &userID)
	if err != nil {
		os.Remove(filePath) // Clean up on error
		if errors.Is(err, ErrStaffNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Staff member not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to update photo"))
		return
	}

	response.OK(c, ToStaffResponse(staff))
}

// PreviewEmployeeID returns the next available employee ID.
// @Summary Preview next employee ID
// @Description Get a preview of the next employee ID that will be generated
// @Tags Staff
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Success 200 {object} response.Success{data=map[string]string}
// @Router /api/v1/staff/employee-id/preview [get]
func (h *Handler) PreviewEmployeeID(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	employeeID, err := h.service.GenerateEmployeeID(c.Request.Context(), tenantID)
	if err != nil {
		logger.Error("Failed to generate employee ID preview",
			zap.String("tenant_id", tenantID.String()),
			zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to generate employee ID preview"))
		return
	}

	response.OK(c, map[string]string{"employeeId": employeeID})
}
