package health

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"msls-backend/internal/middleware"
	"msls-backend/internal/pkg/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// =============================================================================
// Helper Methods
// =============================================================================

func (h *Handler) getTenantID(c *gin.Context) (uuid.UUID, error) {
	tenantID, exists := c.Get(middleware.TenantIDKey)
	if !exists {
		return uuid.Nil, errors.New("tenant ID not found in context")
	}
	// Handle both string and uuid.UUID types
	switch v := tenantID.(type) {
	case string:
		return uuid.Parse(v)
	case uuid.UUID:
		return v, nil
	default:
		return uuid.Nil, errors.New("invalid tenant ID type in context")
	}
}

func (h *Handler) getUserID(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		return uuid.Nil, errors.New("user ID not found in context")
	}
	// Handle both string and uuid.UUID types
	switch v := userID.(type) {
	case string:
		return uuid.Parse(v)
	case uuid.UUID:
		return v, nil
	default:
		return uuid.Nil, errors.New("invalid user ID type in context")
	}
}

func (h *Handler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrHealthProfileNotFound):
		response.NotFound(c, err.Error())
	case errors.Is(err, ErrAllergyNotFound):
		response.NotFound(c, err.Error())
	case errors.Is(err, ErrConditionNotFound):
		response.NotFound(c, err.Error())
	case errors.Is(err, ErrMedicationNotFound):
		response.NotFound(c, err.Error())
	case errors.Is(err, ErrVaccinationNotFound):
		response.NotFound(c, err.Error())
	case errors.Is(err, ErrIncidentNotFound):
		response.NotFound(c, err.Error())
	case errors.Is(err, ErrStudentNotFound):
		response.NotFound(c, err.Error())
	case errors.Is(err, ErrInvalidAllergyType),
		errors.Is(err, ErrInvalidSeverity),
		errors.Is(err, ErrInvalidConditionType),
		errors.Is(err, ErrInvalidFrequency),
		errors.Is(err, ErrInvalidRoute),
		errors.Is(err, ErrInvalidDateRange),
		errors.Is(err, ErrInvalidVaccineType),
		errors.Is(err, ErrInvalidDoseNumber),
		errors.Is(err, ErrInvalidIncidentType):
		response.BadRequest(c, err.Error())
	case errors.Is(err, ErrUnauthorized):
		response.Forbidden(c, err.Error())
	default:
		response.InternalServerError(c, "An unexpected error occurred")
	}
}

// =============================================================================
// Health Profile Handlers
// =============================================================================

// GetHealthProfile godoc
// @Summary Get student health profile
// @Tags Health
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {object} response.Response{data=HealthProfileResponse}
// @Router /students/{id}/health/profile [get]
func (h *Handler) GetHealthProfile(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid tenant")
		return
	}

	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid student ID")
		return
	}

	profile, err := h.service.GetHealthProfile(c.Request.Context(), tenantID, studentID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.OK(c, profile)
}

// CreateOrUpdateHealthProfile godoc
// @Summary Create or update student health profile
// @Tags Health
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param body body CreateHealthProfileRequest true "Health profile data"
// @Success 200 {object} response.Response{data=HealthProfileResponse}
// @Router /students/{id}/health/profile [put]
func (h *Handler) CreateOrUpdateHealthProfile(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid tenant")
		return
	}

	userID, err := h.getUserID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid user")
		return
	}

	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid student ID")
		return
	}

	var req CreateHealthProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	profile, err := h.service.CreateOrUpdateHealthProfile(c.Request.Context(), tenantID, studentID, userID, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.OK(c, profile)
}

// =============================================================================
// Allergy Handlers
// =============================================================================

// ListAllergies godoc
// @Summary List student allergies
// @Tags Health
// @Produce json
// @Param id path string true "Student ID"
// @Param active query bool false "Only active allergies"
// @Success 200 {object} response.Response{data=AllergyListResponse}
// @Router /students/{id}/health/allergies [get]
func (h *Handler) ListAllergies(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid tenant")
		return
	}

	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid student ID")
		return
	}

	activeOnly := c.Query("active") == "true"

	allergies, err := h.service.ListAllergies(c.Request.Context(), tenantID, studentID, activeOnly)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.OK(c, allergies)
}

// GetAllergy godoc
// @Summary Get allergy by ID
// @Tags Health
// @Produce json
// @Param id path string true "Student ID"
// @Param allergyId path string true "Allergy ID"
// @Success 200 {object} response.Response{data=AllergyResponse}
// @Router /students/{id}/health/allergies/{allergyId} [get]
func (h *Handler) GetAllergy(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid tenant")
		return
	}

	allergyID, err := uuid.Parse(c.Param("allergyId"))
	if err != nil {
		response.BadRequest(c, "Invalid allergy ID")
		return
	}

	allergy, err := h.service.GetAllergy(c.Request.Context(), tenantID, allergyID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.OK(c, allergy)
}

// CreateAllergy godoc
// @Summary Create allergy
// @Tags Health
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param body body CreateAllergyRequest true "Allergy data"
// @Success 201 {object} response.Response{data=AllergyResponse}
// @Router /students/{id}/health/allergies [post]
func (h *Handler) CreateAllergy(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid tenant")
		return
	}

	userID, err := h.getUserID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid user")
		return
	}

	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid student ID")
		return
	}

	var req CreateAllergyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	allergy, err := h.service.CreateAllergy(c.Request.Context(), tenantID, studentID, userID, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.Created(c, allergy)
}

// UpdateAllergy godoc
// @Summary Update allergy
// @Tags Health
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param allergyId path string true "Allergy ID"
// @Param body body UpdateAllergyRequest true "Allergy data"
// @Success 200 {object} response.Response{data=AllergyResponse}
// @Router /students/{id}/health/allergies/{allergyId} [put]
func (h *Handler) UpdateAllergy(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid tenant")
		return
	}

	userID, err := h.getUserID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid user")
		return
	}

	allergyID, err := uuid.Parse(c.Param("allergyId"))
	if err != nil {
		response.BadRequest(c, "Invalid allergy ID")
		return
	}

	var req UpdateAllergyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	allergy, err := h.service.UpdateAllergy(c.Request.Context(), tenantID, allergyID, userID, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.OK(c, allergy)
}

// DeleteAllergy godoc
// @Summary Delete allergy
// @Tags Health
// @Param id path string true "Student ID"
// @Param allergyId path string true "Allergy ID"
// @Success 204
// @Router /students/{id}/health/allergies/{allergyId} [delete]
func (h *Handler) DeleteAllergy(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid tenant")
		return
	}

	allergyID, err := uuid.Parse(c.Param("allergyId"))
	if err != nil {
		response.BadRequest(c, "Invalid allergy ID")
		return
	}

	if err := h.service.DeleteAllergy(c.Request.Context(), tenantID, allergyID); err != nil {
		h.handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// =============================================================================
// Chronic Condition Handlers
// =============================================================================

// ListConditions godoc
// @Summary List student chronic conditions
// @Tags Health
// @Produce json
// @Param id path string true "Student ID"
// @Param active query bool false "Only active conditions"
// @Success 200 {object} response.Response{data=ConditionListResponse}
// @Router /students/{id}/health/conditions [get]
func (h *Handler) ListConditions(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid tenant")
		return
	}

	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid student ID")
		return
	}

	activeOnly := c.Query("active") == "true"

	conditions, err := h.service.ListConditions(c.Request.Context(), tenantID, studentID, activeOnly)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.OK(c, conditions)
}

// GetCondition godoc
// @Summary Get condition by ID
// @Tags Health
// @Produce json
// @Param id path string true "Student ID"
// @Param conditionId path string true "Condition ID"
// @Success 200 {object} response.Response{data=ConditionResponse}
// @Router /students/{id}/health/conditions/{conditionId} [get]
func (h *Handler) GetCondition(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid tenant")
		return
	}

	conditionID, err := uuid.Parse(c.Param("conditionId"))
	if err != nil {
		response.BadRequest(c, "Invalid condition ID")
		return
	}

	condition, err := h.service.GetCondition(c.Request.Context(), tenantID, conditionID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.OK(c, condition)
}

// CreateCondition godoc
// @Summary Create chronic condition
// @Tags Health
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param body body CreateConditionRequest true "Condition data"
// @Success 201 {object} response.Response{data=ConditionResponse}
// @Router /students/{id}/health/conditions [post]
func (h *Handler) CreateCondition(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid tenant")
		return
	}

	userID, err := h.getUserID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid user")
		return
	}

	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid student ID")
		return
	}

	var req CreateConditionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	condition, err := h.service.CreateCondition(c.Request.Context(), tenantID, studentID, userID, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.Created(c, condition)
}

// UpdateCondition godoc
// @Summary Update chronic condition
// @Tags Health
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param conditionId path string true "Condition ID"
// @Param body body UpdateConditionRequest true "Condition data"
// @Success 200 {object} response.Response{data=ConditionResponse}
// @Router /students/{id}/health/conditions/{conditionId} [put]
func (h *Handler) UpdateCondition(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid tenant")
		return
	}

	userID, err := h.getUserID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid user")
		return
	}

	conditionID, err := uuid.Parse(c.Param("conditionId"))
	if err != nil {
		response.BadRequest(c, "Invalid condition ID")
		return
	}

	var req UpdateConditionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	condition, err := h.service.UpdateCondition(c.Request.Context(), tenantID, conditionID, userID, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.OK(c, condition)
}

// DeleteCondition godoc
// @Summary Delete chronic condition
// @Tags Health
// @Param id path string true "Student ID"
// @Param conditionId path string true "Condition ID"
// @Success 204
// @Router /students/{id}/health/conditions/{conditionId} [delete]
func (h *Handler) DeleteCondition(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid tenant")
		return
	}

	conditionID, err := uuid.Parse(c.Param("conditionId"))
	if err != nil {
		response.BadRequest(c, "Invalid condition ID")
		return
	}

	if err := h.service.DeleteCondition(c.Request.Context(), tenantID, conditionID); err != nil {
		h.handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// =============================================================================
// Medication Handlers
// =============================================================================

// ListMedications godoc
// @Summary List student medications
// @Tags Health
// @Produce json
// @Param id path string true "Student ID"
// @Param active query bool false "Only active medications"
// @Success 200 {object} response.Response{data=MedicationListResponse}
// @Router /students/{id}/health/medications [get]
func (h *Handler) ListMedications(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid tenant")
		return
	}

	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid student ID")
		return
	}

	activeOnly := c.Query("active") == "true"

	medications, err := h.service.ListMedications(c.Request.Context(), tenantID, studentID, activeOnly)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.OK(c, medications)
}

// GetMedication godoc
// @Summary Get medication by ID
// @Tags Health
// @Produce json
// @Param id path string true "Student ID"
// @Param medicationId path string true "Medication ID"
// @Success 200 {object} response.Response{data=MedicationResponse}
// @Router /students/{id}/health/medications/{medicationId} [get]
func (h *Handler) GetMedication(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid tenant")
		return
	}

	medicationID, err := uuid.Parse(c.Param("medicationId"))
	if err != nil {
		response.BadRequest(c, "Invalid medication ID")
		return
	}

	medication, err := h.service.GetMedication(c.Request.Context(), tenantID, medicationID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.OK(c, medication)
}

// CreateMedication godoc
// @Summary Create medication
// @Tags Health
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param body body CreateMedicationRequest true "Medication data"
// @Success 201 {object} response.Response{data=MedicationResponse}
// @Router /students/{id}/health/medications [post]
func (h *Handler) CreateMedication(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid tenant")
		return
	}

	userID, err := h.getUserID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid user")
		return
	}

	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid student ID")
		return
	}

	var req CreateMedicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	medication, err := h.service.CreateMedication(c.Request.Context(), tenantID, studentID, userID, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.Created(c, medication)
}

// UpdateMedication godoc
// @Summary Update medication
// @Tags Health
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param medicationId path string true "Medication ID"
// @Param body body UpdateMedicationRequest true "Medication data"
// @Success 200 {object} response.Response{data=MedicationResponse}
// @Router /students/{id}/health/medications/{medicationId} [put]
func (h *Handler) UpdateMedication(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid tenant")
		return
	}

	userID, err := h.getUserID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid user")
		return
	}

	medicationID, err := uuid.Parse(c.Param("medicationId"))
	if err != nil {
		response.BadRequest(c, "Invalid medication ID")
		return
	}

	var req UpdateMedicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	medication, err := h.service.UpdateMedication(c.Request.Context(), tenantID, medicationID, userID, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.OK(c, medication)
}

// DeleteMedication godoc
// @Summary Delete medication
// @Tags Health
// @Param id path string true "Student ID"
// @Param medicationId path string true "Medication ID"
// @Success 204
// @Router /students/{id}/health/medications/{medicationId} [delete]
func (h *Handler) DeleteMedication(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid tenant")
		return
	}

	medicationID, err := uuid.Parse(c.Param("medicationId"))
	if err != nil {
		response.BadRequest(c, "Invalid medication ID")
		return
	}

	if err := h.service.DeleteMedication(c.Request.Context(), tenantID, medicationID); err != nil {
		h.handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// =============================================================================
// Vaccination Handlers
// =============================================================================

// ListVaccinations godoc
// @Summary List student vaccinations
// @Tags Health
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {object} response.Response{data=VaccinationListResponse}
// @Router /students/{id}/health/vaccinations [get]
func (h *Handler) ListVaccinations(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid tenant")
		return
	}

	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid student ID")
		return
	}

	vaccinations, err := h.service.ListVaccinations(c.Request.Context(), tenantID, studentID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.OK(c, vaccinations)
}

// GetVaccination godoc
// @Summary Get vaccination by ID
// @Tags Health
// @Produce json
// @Param id path string true "Student ID"
// @Param vaccinationId path string true "Vaccination ID"
// @Success 200 {object} response.Response{data=VaccinationResponse}
// @Router /students/{id}/health/vaccinations/{vaccinationId} [get]
func (h *Handler) GetVaccination(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid tenant")
		return
	}

	vaccinationID, err := uuid.Parse(c.Param("vaccinationId"))
	if err != nil {
		response.BadRequest(c, "Invalid vaccination ID")
		return
	}

	vaccination, err := h.service.GetVaccination(c.Request.Context(), tenantID, vaccinationID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.OK(c, vaccination)
}

// CreateVaccination godoc
// @Summary Create vaccination
// @Tags Health
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param body body CreateVaccinationRequest true "Vaccination data"
// @Success 201 {object} response.Response{data=VaccinationResponse}
// @Router /students/{id}/health/vaccinations [post]
func (h *Handler) CreateVaccination(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid tenant")
		return
	}

	userID, err := h.getUserID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid user")
		return
	}

	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid student ID")
		return
	}

	var req CreateVaccinationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	vaccination, err := h.service.CreateVaccination(c.Request.Context(), tenantID, studentID, userID, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.Created(c, vaccination)
}

// UpdateVaccination godoc
// @Summary Update vaccination
// @Tags Health
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param vaccinationId path string true "Vaccination ID"
// @Param body body UpdateVaccinationRequest true "Vaccination data"
// @Success 200 {object} response.Response{data=VaccinationResponse}
// @Router /students/{id}/health/vaccinations/{vaccinationId} [put]
func (h *Handler) UpdateVaccination(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid tenant")
		return
	}

	userID, err := h.getUserID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid user")
		return
	}

	vaccinationID, err := uuid.Parse(c.Param("vaccinationId"))
	if err != nil {
		response.BadRequest(c, "Invalid vaccination ID")
		return
	}

	var req UpdateVaccinationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	vaccination, err := h.service.UpdateVaccination(c.Request.Context(), tenantID, vaccinationID, userID, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.OK(c, vaccination)
}

// DeleteVaccination godoc
// @Summary Delete vaccination
// @Tags Health
// @Param id path string true "Student ID"
// @Param vaccinationId path string true "Vaccination ID"
// @Success 204
// @Router /students/{id}/health/vaccinations/{vaccinationId} [delete]
func (h *Handler) DeleteVaccination(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid tenant")
		return
	}

	vaccinationID, err := uuid.Parse(c.Param("vaccinationId"))
	if err != nil {
		response.BadRequest(c, "Invalid vaccination ID")
		return
	}

	if err := h.service.DeleteVaccination(c.Request.Context(), tenantID, vaccinationID); err != nil {
		h.handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// =============================================================================
// Medical Incident Handlers
// =============================================================================

// ListIncidents godoc
// @Summary List student medical incidents
// @Tags Health
// @Produce json
// @Param id path string true "Student ID"
// @Param limit query int false "Limit results"
// @Success 200 {object} response.Response{data=IncidentListResponse}
// @Router /students/{id}/health/incidents [get]
func (h *Handler) ListIncidents(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid tenant")
		return
	}

	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid student ID")
		return
	}

	limit := 0 // No limit by default

	incidents, err := h.service.ListIncidents(c.Request.Context(), tenantID, studentID, limit)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.OK(c, incidents)
}

// GetIncident godoc
// @Summary Get medical incident by ID
// @Tags Health
// @Produce json
// @Param id path string true "Student ID"
// @Param incidentId path string true "Incident ID"
// @Success 200 {object} response.Response{data=IncidentResponse}
// @Router /students/{id}/health/incidents/{incidentId} [get]
func (h *Handler) GetIncident(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid tenant")
		return
	}

	incidentID, err := uuid.Parse(c.Param("incidentId"))
	if err != nil {
		response.BadRequest(c, "Invalid incident ID")
		return
	}

	incident, err := h.service.GetIncident(c.Request.Context(), tenantID, incidentID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.OK(c, incident)
}

// CreateIncident godoc
// @Summary Create medical incident
// @Tags Health
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param body body CreateIncidentRequest true "Incident data"
// @Success 201 {object} response.Response{data=IncidentResponse}
// @Router /students/{id}/health/incidents [post]
func (h *Handler) CreateIncident(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid tenant")
		return
	}

	userID, err := h.getUserID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid user")
		return
	}

	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid student ID")
		return
	}

	var req CreateIncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	incident, err := h.service.CreateIncident(c.Request.Context(), tenantID, studentID, userID, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.Created(c, incident)
}

// UpdateIncident godoc
// @Summary Update medical incident
// @Tags Health
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param incidentId path string true "Incident ID"
// @Param body body UpdateIncidentRequest true "Incident data"
// @Success 200 {object} response.Response{data=IncidentResponse}
// @Router /students/{id}/health/incidents/{incidentId} [put]
func (h *Handler) UpdateIncident(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid tenant")
		return
	}

	userID, err := h.getUserID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid user")
		return
	}

	incidentID, err := uuid.Parse(c.Param("incidentId"))
	if err != nil {
		response.BadRequest(c, "Invalid incident ID")
		return
	}

	var req UpdateIncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	incident, err := h.service.UpdateIncident(c.Request.Context(), tenantID, incidentID, userID, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.OK(c, incident)
}

// DeleteIncident godoc
// @Summary Delete medical incident
// @Tags Health
// @Param id path string true "Student ID"
// @Param incidentId path string true "Incident ID"
// @Success 204
// @Router /students/{id}/health/incidents/{incidentId} [delete]
func (h *Handler) DeleteIncident(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid tenant")
		return
	}

	incidentID, err := uuid.Parse(c.Param("incidentId"))
	if err != nil {
		response.BadRequest(c, "Invalid incident ID")
		return
	}

	if err := h.service.DeleteIncident(c.Request.Context(), tenantID, incidentID); err != nil {
		h.handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// =============================================================================
// Health Summary Handler
// =============================================================================

// GetHealthSummary godoc
// @Summary Get complete health summary
// @Tags Health
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {object} response.Response{data=HealthSummaryResponse}
// @Router /students/{id}/health [get]
func (h *Handler) GetHealthSummary(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		response.Unauthorized(c, "Invalid tenant")
		return
	}

	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid student ID")
		return
	}

	summary, err := h.service.GetHealthSummary(c.Request.Context(), tenantID, studentID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.OK(c, summary)
}
