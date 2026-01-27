// Package admission provides HTTP handlers for admission management.
package admission

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"msls-backend/internal/middleware"
	apperrors "msls-backend/internal/pkg/errors"
	"msls-backend/internal/pkg/response"
	"msls-backend/internal/services/admission"
)

// TestHandler handles entrance test HTTP requests.
type TestHandler struct {
	testService *admission.TestService
}

// NewTestHandler creates a new TestHandler.
func NewTestHandler(testService *admission.TestService) *TestHandler {
	return &TestHandler{testService: testService}
}

// ListTests godoc
// @Summary List entrance tests
// @Description Retrieves a list of entrance tests with optional filtering
// @Tags Entrance Tests
// @Accept json
// @Produce json
// @Param sessionId query string false "Filter by session ID"
// @Param status query string false "Filter by status (scheduled, in_progress, completed, cancelled)"
// @Param className query string false "Filter by class name"
// @Success 200 {object} response.Success{data=TestListResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 500 {object} apperrors.AppError
// @Router /api/v1/entrance-tests [get]
func (h *TestHandler) ListTests(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	filter := admission.TestListFilter{TenantID: tenantID}

	// Parse optional filters
	if sessionID := c.Query("sessionId"); sessionID != "" {
		id, err := uuid.Parse(sessionID)
		if err == nil {
			filter.SessionID = &id
		}
	}
	if status := c.Query("status"); status != "" {
		s := admission.EntranceTestStatus(status)
		if s.IsValid() {
			filter.Status = &s
		}
	}
	if className := c.Query("className"); className != "" {
		filter.ClassName = className
	}

	tests, err := h.testService.ListTests(c.Request.Context(), filter)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve entrance tests"))
		return
	}

	// Build response with registration counts
	testResponses := make([]TestResponse, len(tests))
	for i, test := range tests {
		count, _ := h.testService.CountRegistrations(c.Request.Context(), tenantID, test.ID)
		testResponses[i] = NewTestResponse(&test, count)
	}

	response.OK(c, TestListResponse{
		Tests: testResponses,
		Total: len(testResponses),
	})
}

// GetTest godoc
// @Summary Get entrance test by ID
// @Description Retrieves a single entrance test by its ID
// @Tags Entrance Tests
// @Accept json
// @Produce json
// @Param id path string true "Test ID"
// @Success 200 {object} response.Success{data=TestResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Failure 500 {object} apperrors.AppError
// @Router /api/v1/entrance-tests/{id} [get]
func (h *TestHandler) GetTest(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	testID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid test ID"))
		return
	}

	test, err := h.testService.GetTestByID(c.Request.Context(), tenantID, testID)
	if err != nil {
		if err == admission.ErrTestNotFound {
			apperrors.Abort(c, apperrors.NotFound("Entrance test not found"))
		} else {
			apperrors.Abort(c, apperrors.InternalError("Failed to retrieve entrance test"))
		}
		return
	}

	count, _ := h.testService.CountRegistrations(c.Request.Context(), tenantID, testID)
	response.OK(c, NewTestResponse(test, count))
}

// CreateTest godoc
// @Summary Create entrance test
// @Description Creates a new entrance test
// @Tags Entrance Tests
// @Accept json
// @Produce json
// @Param request body CreateTestRequest true "Test details"
// @Success 201 {object} response.Success{data=TestResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 500 {object} apperrors.AppError
// @Router /api/v1/entrance-tests [post]
func (h *TestHandler) CreateTest(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	var req CreateTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	userID, ok := middleware.GetCurrentUserID(c)
	var userIDPtr *uuid.UUID
	if ok {
		userIDPtr = &userID
	}
	svcReq, err := req.ToServiceRequest(tenantID, userIDPtr)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	test, err := h.testService.CreateTest(c.Request.Context(), *svcReq)
	if err != nil {
		if err == admission.ErrTestAlreadyExists {
			apperrors.Abort(c, apperrors.Conflict("Test with this name already exists"))
		} else {
			apperrors.Abort(c, apperrors.InternalError("Failed to create entrance test"))
		}
		return
	}

	response.Created(c, NewTestResponse(test, 0))
}

// UpdateTest godoc
// @Summary Update entrance test
// @Description Updates an existing entrance test
// @Tags Entrance Tests
// @Accept json
// @Produce json
// @Param id path string true "Test ID"
// @Param request body UpdateTestRequest true "Test updates"
// @Success 200 {object} response.Success{data=TestResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Failure 500 {object} apperrors.AppError
// @Router /api/v1/entrance-tests/{id} [put]
func (h *TestHandler) UpdateTest(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	testID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid test ID"))
		return
	}

	var req UpdateTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	userID, ok := middleware.GetCurrentUserID(c)
	var userIDPtr *uuid.UUID
	if ok {
		userIDPtr = &userID
	}
	svcReq, err := req.ToServiceRequest(userIDPtr)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	test, err := h.testService.UpdateTest(c.Request.Context(), tenantID, testID, *svcReq)
	if err != nil {
		if err == admission.ErrTestNotFound {
			apperrors.Abort(c, apperrors.NotFound("Entrance test not found"))
		} else if err == admission.ErrCannotModifyCompletedTest {
			apperrors.Abort(c, apperrors.Conflict("Cannot modify completed or cancelled test"))
		} else {
			apperrors.Abort(c, apperrors.InternalError("Failed to update entrance test"))
		}
		return
	}

	count, _ := h.testService.CountRegistrations(c.Request.Context(), tenantID, testID)
	response.OK(c, NewTestResponse(test, count))
}

// DeleteTest godoc
// @Summary Delete entrance test
// @Description Deletes an entrance test
// @Tags Entrance Tests
// @Accept json
// @Produce json
// @Param id path string true "Test ID"
// @Success 204 "No Content"
// @Failure 400 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Failure 409 {object} apperrors.AppError
// @Failure 500 {object} apperrors.AppError
// @Router /api/v1/entrance-tests/{id} [delete]
func (h *TestHandler) DeleteTest(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	testID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid test ID"))
		return
	}

	err = h.testService.DeleteTest(c.Request.Context(), tenantID, testID)
	if err != nil {
		if err == admission.ErrTestNotFound {
			apperrors.Abort(c, apperrors.NotFound("Entrance test not found"))
		} else if err == admission.ErrCannotModifyCompletedTest {
			apperrors.Abort(c, apperrors.Conflict("Cannot delete completed test"))
		} else {
			apperrors.Abort(c, apperrors.InternalError("Failed to delete entrance test"))
		}
		return
	}

	response.NoContent(c)
}

// ListRegistrations godoc
// @Summary List test registrations
// @Description Retrieves registrations for an entrance test
// @Tags Entrance Tests
// @Accept json
// @Produce json
// @Param id path string true "Test ID"
// @Param status query string false "Filter by status"
// @Success 200 {object} response.Success{data=[]TestRegistrationResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 500 {object} apperrors.AppError
// @Router /api/v1/entrance-tests/{id}/registrations [get]
func (h *TestHandler) ListRegistrations(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	testID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid test ID"))
		return
	}

	filter := admission.RegistrationListFilter{
		TenantID: tenantID,
		TestID:   &testID,
	}

	if status := c.Query("status"); status != "" {
		s := admission.TestRegistrationStatus(status)
		if s.IsValid() {
			filter.Status = &s
		}
	}

	registrations, err := h.testService.ListRegistrations(c.Request.Context(), filter)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve registrations"))
		return
	}

	regResponses := make([]TestRegistrationResponse, len(registrations))
	for i, reg := range registrations {
		regResponses[i] = NewTestRegistrationResponse(&reg)
	}

	response.OK(c, regResponses)
}

// RegisterCandidate godoc
// @Summary Register candidate for test
// @Description Registers a candidate (application) for an entrance test
// @Tags Entrance Tests
// @Accept json
// @Produce json
// @Param id path string true "Test ID"
// @Param request body RegisterCandidateRequest true "Registration details"
// @Success 201 {object} response.Success{data=TestRegistrationResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 409 {object} apperrors.AppError
// @Failure 500 {object} apperrors.AppError
// @Router /api/v1/entrance-tests/{id}/register [post]
func (h *TestHandler) RegisterCandidate(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	testID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid test ID"))
		return
	}

	var req RegisterCandidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	applicationID, err := uuid.Parse(req.ApplicationID)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid application ID"))
		return
	}

	registration, err := h.testService.RegisterCandidate(c.Request.Context(), admission.RegisterCandidateRequest{
		TenantID:      tenantID,
		TestID:        testID,
		ApplicationID: applicationID,
	})
	if err != nil {
		switch err {
		case admission.ErrTestNotFound, admission.ErrApplicationNotFound:
			apperrors.Abort(c, apperrors.NotFound(err.Error()))
		case admission.ErrAlreadyRegistered:
			apperrors.Abort(c, apperrors.Conflict("Candidate is already registered for this test"))
		case admission.ErrTestFull:
			apperrors.Abort(c, apperrors.Conflict("Test has reached maximum capacity"))
		case admission.ErrTestNotScheduled:
			apperrors.Abort(c, apperrors.Conflict("Test is not in scheduled status"))
		case admission.ErrTestDatePassed:
			apperrors.Abort(c, apperrors.Conflict("Test date has passed"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to register candidate"))
		}
		return
	}

	response.Created(c, NewTestRegistrationResponse(registration))
}

// CancelRegistration godoc
// @Summary Cancel test registration
// @Description Cancels a candidate's registration for an entrance test
// @Tags Entrance Tests
// @Accept json
// @Produce json
// @Param id path string true "Test ID"
// @Param registrationId path string true "Registration ID"
// @Success 204 "No Content"
// @Failure 400 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Failure 500 {object} apperrors.AppError
// @Router /api/v1/entrance-tests/{id}/registrations/{registrationId} [delete]
func (h *TestHandler) CancelRegistration(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	registrationID, err := uuid.Parse(c.Param("registrationId"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid registration ID"))
		return
	}

	err = h.testService.CancelRegistration(c.Request.Context(), tenantID, registrationID)
	if err != nil {
		if err == admission.ErrRegistrationNotFound {
			apperrors.Abort(c, apperrors.NotFound("Registration not found"))
		} else {
			apperrors.Abort(c, apperrors.InternalError("Failed to cancel registration"))
		}
		return
	}

	response.NoContent(c)
}

// SubmitResult godoc
// @Summary Submit test result
// @Description Submits test results for a single registration
// @Tags Entrance Tests
// @Accept json
// @Produce json
// @Param id path string true "Test ID"
// @Param request body SubmitResultRequest true "Result details"
// @Success 200 {object} response.Success{data=TestRegistrationResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Failure 409 {object} apperrors.AppError
// @Failure 500 {object} apperrors.AppError
// @Router /api/v1/entrance-tests/{id}/results [post]
func (h *TestHandler) SubmitResult(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	testID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid test ID"))
		return
	}

	var req SubmitResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	registrationID, err := uuid.Parse(req.RegistrationID)
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid registration ID"))
		return
	}

	// Convert marks to decimal
	marksDecimal := make(map[string]decimal.Decimal)
	for subject, marks := range req.Marks {
		marksDecimal[subject] = decimal.NewFromFloat(marks)
	}

	registration, err := h.testService.SubmitResults(c.Request.Context(), admission.SubmitResultsRequest{
		TenantID:       tenantID,
		TestID:         testID,
		RegistrationID: registrationID,
		Marks:          marksDecimal,
		Remarks:        req.Remarks,
	})
	if err != nil {
		switch err {
		case admission.ErrRegistrationNotFound, admission.ErrTestNotFound:
			apperrors.Abort(c, apperrors.NotFound(err.Error()))
		case admission.ErrResultsAlreadySubmitted:
			apperrors.Abort(c, apperrors.Conflict("Results have already been submitted"))
		default:
			apperrors.Abort(c, apperrors.InternalError("Failed to submit results"))
		}
		return
	}

	response.OK(c, NewTestRegistrationResponse(registration))
}

// BulkSubmitResults godoc
// @Summary Submit bulk test results
// @Description Submits test results for multiple registrations
// @Tags Entrance Tests
// @Accept json
// @Produce json
// @Param id path string true "Test ID"
// @Param request body BulkSubmitResultsRequest true "Bulk results"
// @Success 200 {object} response.Success{data=[]TestRegistrationResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 500 {object} apperrors.AppError
// @Router /api/v1/entrance-tests/{id}/results/bulk [post]
func (h *TestHandler) BulkSubmitResults(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	testID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid test ID"))
		return
	}

	var req BulkSubmitResultsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Abort(c, apperrors.BadRequest(err.Error()))
		return
	}

	// Convert to service request
	svcResults := make([]admission.SingleResultRequest, len(req.Results))
	for i, r := range req.Results {
		regID, err := uuid.Parse(r.RegistrationID)
		if err != nil {
			apperrors.Abort(c, apperrors.BadRequest("Invalid registration ID"))
			return
		}
		marksDecimal := make(map[string]decimal.Decimal)
		for subject, marks := range r.Marks {
			marksDecimal[subject] = decimal.NewFromFloat(marks)
		}
		svcResults[i] = admission.SingleResultRequest{
			RegistrationID: regID,
			Marks:          marksDecimal,
			Remarks:        r.Remarks,
		}
	}

	registrations, err := h.testService.BulkSubmitResults(c.Request.Context(), admission.BulkSubmitResultsRequest{
		TenantID: tenantID,
		TestID:   testID,
		Results:  svcResults,
	})
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to submit results"))
		return
	}

	regResponses := make([]TestRegistrationResponse, len(registrations))
	for i, reg := range registrations {
		regResponses[i] = NewTestRegistrationResponse(&reg)
	}

	response.OK(c, regResponses)
}

// GetHallTickets godoc
// @Summary Get hall tickets for test
// @Description Retrieves hall ticket data for all registrations of a test
// @Tags Entrance Tests
// @Accept json
// @Produce json
// @Param id path string true "Test ID"
// @Success 200 {object} response.Success{data=[]HallTicketResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 500 {object} apperrors.AppError
// @Router /api/v1/entrance-tests/{id}/hall-tickets [get]
func (h *TestHandler) GetHallTickets(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	testID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid test ID"))
		return
	}

	hallTickets, err := h.testService.GetHallTickets(c.Request.Context(), tenantID, testID)
	if err != nil {
		apperrors.Abort(c, apperrors.InternalError("Failed to retrieve hall tickets"))
		return
	}

	ticketResponses := make([]HallTicketResponse, len(hallTickets))
	for i, ticket := range hallTickets {
		ticketResponses[i] = NewHallTicketResponse(&ticket)
	}

	response.OK(c, ticketResponses)
}

// GetHallTicket godoc
// @Summary Get hall ticket for registration
// @Description Retrieves hall ticket data for a specific registration
// @Tags Entrance Tests
// @Accept json
// @Produce json
// @Param id path string true "Test ID"
// @Param registrationId path string true "Registration ID"
// @Success 200 {object} response.Success{data=HallTicketResponse}
// @Failure 400 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Failure 500 {object} apperrors.AppError
// @Router /api/v1/entrance-tests/{id}/hall-tickets/{registrationId} [get]
func (h *TestHandler) GetHallTicket(c *gin.Context) {
	tenantID, ok := middleware.GetCurrentTenantID(c)
	if !ok {
		apperrors.Abort(c, apperrors.BadRequest("Tenant ID is required"))
		return
	}

	registrationID, err := uuid.Parse(c.Param("registrationId"))
	if err != nil {
		apperrors.Abort(c, apperrors.BadRequest("Invalid registration ID"))
		return
	}

	hallTicket, err := h.testService.GenerateHallTicket(c.Request.Context(), tenantID, registrationID)
	if err != nil {
		if err == admission.ErrRegistrationNotFound {
			apperrors.Abort(c, apperrors.NotFound("Registration not found"))
		} else {
			apperrors.Abort(c, apperrors.InternalError("Failed to generate hall ticket"))
		}
		return
	}

	response.OK(c, NewHallTicketResponse(hallTicket))
}
