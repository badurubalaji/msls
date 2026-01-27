// Package student provides student management functionality.
package student

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
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

// Handler handles student-related HTTP requests.
type Handler struct {
	service *Service
}

// NewHandler creates a new student handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateStudentRequest represents the request body for creating a student.
type CreateStudentRequest struct {
	BranchID            string          `json:"branchId" binding:"required,uuid"`
	FirstName           string          `json:"firstName" binding:"required,max=100"`
	MiddleName          string          `json:"middleName" binding:"max=100"`
	LastName            string          `json:"lastName" binding:"required,max=100"`
	DateOfBirth         string          `json:"dateOfBirth" binding:"required"` // Format: YYYY-MM-DD
	Gender              string          `json:"gender" binding:"required,oneof=male female other"`
	BloodGroup          string          `json:"bloodGroup" binding:"max=5"`
	AadhaarNumber       string          `json:"aadhaarNumber" binding:"max=12"`
	AdmissionDate       string          `json:"admissionDate"` // Format: YYYY-MM-DD, optional
	CurrentAddress      *AddressRequest `json:"currentAddress"`
	PermanentAddress    *AddressRequest `json:"permanentAddress"`
	SameAsCurrentAddress bool            `json:"sameAsCurrentAddress"`
}

// UpdateStudentRequest represents the request body for updating a student.
type UpdateStudentRequest struct {
	FirstName           *string         `json:"firstName" binding:"omitempty,max=100"`
	MiddleName          *string         `json:"middleName" binding:"omitempty,max=100"`
	LastName            *string         `json:"lastName" binding:"omitempty,max=100"`
	DateOfBirth         *string         `json:"dateOfBirth"` // Format: YYYY-MM-DD
	Gender              *string         `json:"gender" binding:"omitempty,oneof=male female other"`
	BloodGroup          *string         `json:"bloodGroup" binding:"omitempty,max=5"`
	AadhaarNumber       *string         `json:"aadhaarNumber" binding:"omitempty,max=12"`
	Status              *string         `json:"status" binding:"omitempty,oneof=active inactive transferred graduated"`
	CurrentAddress      *AddressRequest `json:"currentAddress"`
	PermanentAddress    *AddressRequest `json:"permanentAddress"`
	SameAsCurrentAddress *bool           `json:"sameAsCurrentAddress"`
	Version             int             `json:"version"` // For optimistic locking
}

// AddressRequest represents an address in requests.
type AddressRequest struct {
	AddressLine1 string `json:"addressLine1" binding:"required,max=255"`
	AddressLine2 string `json:"addressLine2" binding:"max=255"`
	City         string `json:"city" binding:"required,max=100"`
	State        string `json:"state" binding:"required,max=100"`
	PostalCode   string `json:"postalCode" binding:"required,max=10"`
	Country      string `json:"country" binding:"max=100"`
}

// List returns all students for the tenant.
// @Summary List students
// @Description Get all students for the current tenant with pagination and filters
// @Tags Students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param search query string false "Search by name or admission number"
// @Param branch_id query string false "Filter by branch ID"
// @Param class_id query string false "Filter by class ID (via active enrollment)"
// @Param section_id query string false "Filter by section ID (via active enrollment)"
// @Param status query string false "Filter by status (active, inactive, transferred, graduated)"
// @Param gender query string false "Filter by gender (male, female, other)"
// @Param admission_from query string false "Filter by admission date from (YYYY-MM-DD)"
// @Param admission_to query string false "Filter by admission date to (YYYY-MM-DD)"
// @Param cursor query string false "Cursor for pagination"
// @Param limit query int false "Number of results per page (max 100)"
// @Param sort_by query string false "Sort by field (name, admission_number, created_at)"
// @Param sort_order query string false "Sort order (asc, desc)"
// @Success 200 {object} response.Success{data=StudentListResponse}
// @Failure 401 {object} apperrors.AppError
// @Failure 403 {object} apperrors.AppError
// @Router /api/v1/students [get]
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

	// Parse class_id filter
	if classIDStr := c.Query("class_id"); classIDStr != "" {
		classID, err := uuid.Parse(classIDStr)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid class ID"))
			return
		}
		filter.ClassID = &classID
	}

	// Parse section_id filter
	if sectionIDStr := c.Query("section_id"); sectionIDStr != "" {
		sectionID, err := uuid.Parse(sectionIDStr)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid section ID"))
			return
		}
		filter.SectionID = &sectionID
	}

	// Parse status filter
	if statusStr := c.Query("status"); statusStr != "" {
		status := models.StudentStatus(statusStr)
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

	// Parse admission date range
	if admFromStr := c.Query("admission_from"); admFromStr != "" {
		admFrom, err := time.Parse("2006-01-02", admFromStr)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid admission_from date format, use YYYY-MM-DD"))
			return
		}
		filter.AdmissionFrom = &admFrom
	}

	if admToStr := c.Query("admission_to"); admToStr != "" {
		admTo, err := time.Parse("2006-01-02", admToStr)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid admission_to date format, use YYYY-MM-DD"))
			return
		}
		filter.AdmissionTo = &admTo
	}

	// Parse limit
	if limitStr := c.Query("limit"); limitStr != "" {
		var limit int
		if _, err := parseIntFromString(limitStr, &limit); err == nil && limit > 0 {
			filter.Limit = limit
		}
	}

	students, nextCursor, total, err := h.service.List(c.Request.Context(), filter)
	if err != nil {
		logger.Error("Failed to list students",
			zap.String("tenant_id", tenantID.String()),
			zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve students"))
		return
	}

	resp := StudentListResponse{
		Students:   ToStudentResponses(students),
		NextCursor: nextCursor,
		HasMore:    nextCursor != "",
		Total:      total,
	}

	response.OK(c, resp)
}

// GetByID returns a student by ID.
// @Summary Get student by ID
// @Description Get a student with full details
// @Tags Students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Student ID" format(uuid)
// @Success 200 {object} response.Success{data=StudentResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/students/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid student ID"))
		return
	}

	student, err := h.service.GetByID(c.Request.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, ErrStudentNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Student not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve student"))
		return
	}

	response.OK(c, ToStudentResponse(student))
}

// Create creates a new student.
// @Summary Create student
// @Description Create a new student profile
// @Tags Students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param request body CreateStudentRequest true "Student details"
// @Success 201 {object} response.Success{data=StudentResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 409 {object} apperrors.AppError
// @Router /api/v1/students [post]
func (h *Handler) Create(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	var req CreateStudentRequest
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

	dto := CreateStudentDTO{
		TenantID:            tenantID,
		BranchID:            branchID,
		FirstName:           req.FirstName,
		MiddleName:          req.MiddleName,
		LastName:            req.LastName,
		DateOfBirth:         dob,
		Gender:              models.Gender(req.Gender),
		BloodGroup:          req.BloodGroup,
		AadhaarNumber:       req.AadhaarNumber,
		SameAsCurrentAddress: req.SameAsCurrentAddress,
		CreatedBy:           &userID,
	}

	// Parse admission date if provided
	if req.AdmissionDate != "" {
		admDate, err := time.Parse("2006-01-02", req.AdmissionDate)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid admission date format, use YYYY-MM-DD"))
			return
		}
		dto.AdmissionDate = &admDate
	}

	// Convert addresses
	if req.CurrentAddress != nil {
		dto.CurrentAddress = &AddressDTO{
			AddressLine1: req.CurrentAddress.AddressLine1,
			AddressLine2: req.CurrentAddress.AddressLine2,
			City:         req.CurrentAddress.City,
			State:        req.CurrentAddress.State,
			PostalCode:   req.CurrentAddress.PostalCode,
			Country:      req.CurrentAddress.Country,
		}
	}

	if req.PermanentAddress != nil {
		dto.PermanentAddress = &AddressDTO{
			AddressLine1: req.PermanentAddress.AddressLine1,
			AddressLine2: req.PermanentAddress.AddressLine2,
			City:         req.PermanentAddress.City,
			State:        req.PermanentAddress.State,
			PostalCode:   req.PermanentAddress.PostalCode,
			Country:      req.PermanentAddress.Country,
		}
	}

	student, err := h.service.Create(c.Request.Context(), dto)
	if err != nil {
		switch {
		case errors.Is(err, ErrBranchNotFound):
			apperrors.Abort(c, apperrors.BadRequest("Branch not found"))
		case errors.Is(err, ErrDuplicateAdmissionNumber):
			apperrors.Abort(c, apperrors.Conflict("Admission number already exists"))
		case errors.Is(err, ErrFirstNameRequired):
			apperrors.Abort(c, apperrors.BadRequest("First name is required"))
		case errors.Is(err, ErrLastNameRequired):
			apperrors.Abort(c, apperrors.BadRequest("Last name is required"))
		case errors.Is(err, ErrDateOfBirthRequired):
			apperrors.Abort(c, apperrors.BadRequest("Date of birth is required"))
		case errors.Is(err, ErrInvalidDateOfBirth):
			apperrors.Abort(c, apperrors.BadRequest("Date of birth cannot be in the future"))
		case errors.Is(err, ErrInvalidGender):
			apperrors.Abort(c, apperrors.BadRequest("Invalid gender value"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to create student"))
		}
		return
	}

	response.Created(c, ToStudentResponse(student))
}

// Update updates a student.
// @Summary Update student
// @Description Update an existing student profile
// @Tags Students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Student ID" format(uuid)
// @Param request body UpdateStudentRequest true "Student updates"
// @Success 200 {object} response.Success{data=StudentResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Failure 409 {object} apperrors.AppError
// @Router /api/v1/students/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid student ID"))
		return
	}

	var req UpdateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	dto := UpdateStudentDTO{
		FirstName:           req.FirstName,
		MiddleName:          req.MiddleName,
		LastName:            req.LastName,
		BloodGroup:          req.BloodGroup,
		AadhaarNumber:       req.AadhaarNumber,
		SameAsCurrentAddress: req.SameAsCurrentAddress,
		Version:             req.Version,
		UpdatedBy:           &userID,
	}

	// Parse date of birth if provided
	if req.DateOfBirth != nil {
		dob, err := time.Parse("2006-01-02", *req.DateOfBirth)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid date of birth format, use YYYY-MM-DD"))
			return
		}
		dto.DateOfBirth = &dob
	}

	// Parse gender if provided
	if req.Gender != nil {
		gender := models.Gender(*req.Gender)
		dto.Gender = &gender
	}

	// Parse status if provided
	if req.Status != nil {
		status := models.StudentStatus(*req.Status)
		dto.Status = &status
	}

	// Convert addresses
	if req.CurrentAddress != nil {
		dto.CurrentAddress = &AddressDTO{
			AddressLine1: req.CurrentAddress.AddressLine1,
			AddressLine2: req.CurrentAddress.AddressLine2,
			City:         req.CurrentAddress.City,
			State:        req.CurrentAddress.State,
			PostalCode:   req.CurrentAddress.PostalCode,
			Country:      req.CurrentAddress.Country,
		}
	}

	if req.PermanentAddress != nil {
		dto.PermanentAddress = &AddressDTO{
			AddressLine1: req.PermanentAddress.AddressLine1,
			AddressLine2: req.PermanentAddress.AddressLine2,
			City:         req.PermanentAddress.City,
			State:        req.PermanentAddress.State,
			PostalCode:   req.PermanentAddress.PostalCode,
			Country:      req.PermanentAddress.Country,
		}
	}

	student, err := h.service.Update(c.Request.Context(), tenantID, id, dto)
	if err != nil {
		switch {
		case errors.Is(err, ErrStudentNotFound):
			apperrors.Abort(c, apperrors.NotFound("Student not found"))
		case errors.Is(err, ErrOptimisticLockConflict):
			apperrors.Abort(c, apperrors.Conflict("Student was modified by another user, please refresh and try again"))
		case errors.Is(err, ErrInvalidDateOfBirth):
			apperrors.Abort(c, apperrors.BadRequest("Date of birth cannot be in the future"))
		case errors.Is(err, ErrInvalidGender):
			apperrors.Abort(c, apperrors.BadRequest("Invalid gender value"))
		case errors.Is(err, ErrInvalidStatus):
			apperrors.Abort(c, apperrors.BadRequest("Invalid status value"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to update student"))
		}
		return
	}

	response.OK(c, ToStudentResponse(student))
}

// Delete soft deletes a student.
// @Summary Delete student
// @Description Soft delete a student
// @Tags Students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Student ID" format(uuid)
// @Success 204 "No Content"
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/students/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid student ID"))
		return
	}

	err = h.service.Delete(c.Request.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, ErrStudentNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Student not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to delete student"))
		return
	}

	response.NoContent(c)
}

// UpdatePhoto uploads and updates the photo for a student.
// @Summary Upload student photo
// @Description Upload or update the student's photo (JPEG or PNG, max 2MB)
// @Tags Students
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Student ID" format(uuid)
// @Param photo formData file true "Photo image (JPEG or PNG, max 2MB)"
// @Success 200 {object} response.Success{data=StudentResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /api/v1/students/{id}/photo [post]
func (h *Handler) UpdatePhoto(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	userID, _ := middleware.GetCurrentUserID(c)

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid student ID"))
		return
	}

	// Get the uploaded file
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

	// Create upload directory
	uploadDir := filepath.Join("uploads", "students", tenantID.String())
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		logger.Error("Failed to create upload directory", zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to upload photo"))
		return
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	if ext == "" {
		if contentType == "image/jpeg" {
			ext = ".jpg"
		} else {
			ext = ".png"
		}
	}
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filePath := filepath.Join(uploadDir, filename)

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		logger.Error("Failed to create file", zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to upload photo"))
		return
	}
	defer dst.Close()

	// Copy file contents
	if _, err := io.Copy(dst, file); err != nil {
		logger.Error("Failed to save file", zap.Error(err))
		apperrors.Abort(c, apperrors.InternalError("Failed to upload photo"))
		return
	}

	// Build photo URL
	photoURL := fmt.Sprintf("/uploads/students/%s/%s", tenantID.String(), filename)

	// Update student record
	student, err := h.service.UpdatePhoto(c.Request.Context(), tenantID, id, photoURL, &userID)
	if err != nil {
		// Clean up uploaded file on error
		os.Remove(filePath)
		if errors.Is(err, ErrStudentNotFound) {
			apperrors.Abort(c, apperrors.NotFound("Student not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to update photo"))
		return
	}

	response.OK(c, ToStudentResponse(student))
}

// GetNextAdmissionNumber returns the next admission number for preview.
// @Summary Get next admission number
// @Description Preview the next admission number for a branch
// @Tags Students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param branch_id query string true "Branch ID"
// @Success 200 {object} response.Success{data=AdmissionNumberResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Router /api/v1/students/next-admission-number [get]
func (h *Handler) GetNextAdmissionNumber(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	branchIDStr := c.Query("branch_id")
	if branchIDStr == "" {
		apperrors.Abort(c, apperrors.BadRequest("Branch ID is required"))
		return
	}

	branchID, err := uuid.Parse(branchIDStr)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid branch ID"))
		return
	}

	admissionNumber, err := h.service.GenerateAdmissionNumber(c.Request.Context(), tenantID, branchID)
	if err != nil {
		if errors.Is(err, ErrBranchNotFound) {
			apperrors.Abort(c, apperrors.BadRequest("Branch not found"))
			return
		}
		apperrors.Abort(c, apperrors.InternalError("Failed to generate admission number"))
		return
	}

	response.OK(c, AdmissionNumberResponse{AdmissionNumber: admissionNumber})
}

// AdmissionNumberResponse represents the response for next admission number.
type AdmissionNumberResponse struct {
	AdmissionNumber string `json:"admissionNumber"`
}

// Helper function to parse int from string
func parseIntFromString(s string, v *int) (bool, error) {
	var temp int
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return false, errors.New("invalid number")
		}
		temp = temp*10 + int(ch-'0')
	}
	*v = temp
	return true, nil
}
