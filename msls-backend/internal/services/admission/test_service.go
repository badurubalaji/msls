// Package admission provides admission management services.
package admission

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// TestService handles entrance test operations.
type TestService struct {
	db *gorm.DB
}

// NewTestService creates a new TestService instance.
func NewTestService(db *gorm.DB) *TestService {
	return &TestService{db: db}
}

// CreateTestRequest represents a request to create an entrance test.
type CreateTestRequest struct {
	TenantID          uuid.UUID
	SessionID         uuid.UUID
	TestName          string
	TestDate          time.Time
	StartTime         string
	DurationMinutes   int
	Venue             string
	ClassNames        []string
	MaxCandidates     int
	Subjects          []SubjectMarks
	Instructions      string
	PassingPercentage decimal.Decimal
	CreatedBy         *uuid.UUID
}

// UpdateTestRequest represents a request to update an entrance test.
type UpdateTestRequest struct {
	TestName          *string
	TestDate          *time.Time
	StartTime         *string
	DurationMinutes   *int
	Venue             *string
	ClassNames        []string
	MaxCandidates     *int
	Subjects          []SubjectMarks
	Instructions      *string
	PassingPercentage *decimal.Decimal
	Status            *EntranceTestStatus
	UpdatedBy         *uuid.UUID
}

// RegisterCandidateRequest represents a request to register a candidate for a test.
type RegisterCandidateRequest struct {
	TenantID      uuid.UUID
	TestID        uuid.UUID
	ApplicationID uuid.UUID
}

// SubmitResultsRequest represents a request to submit test results.
type SubmitResultsRequest struct {
	TenantID       uuid.UUID
	TestID         uuid.UUID
	RegistrationID uuid.UUID
	Marks          map[string]decimal.Decimal
	Remarks        string
}

// BulkSubmitResultsRequest represents a request to submit results for multiple candidates.
type BulkSubmitResultsRequest struct {
	TenantID uuid.UUID
	TestID   uuid.UUID
	Results  []SingleResultRequest
}

// SingleResultRequest represents results for a single candidate.
type SingleResultRequest struct {
	RegistrationID uuid.UUID
	Marks          map[string]decimal.Decimal
	Remarks        string
}

// TestListFilter contains filters for listing entrance tests.
type TestListFilter struct {
	TenantID   uuid.UUID
	SessionID  *uuid.UUID
	Status     *EntranceTestStatus
	ClassName  string
	FromDate   *time.Time
	ToDate     *time.Time
}

// RegistrationListFilter contains filters for listing test registrations.
type RegistrationListFilter struct {
	TenantID      uuid.UUID
	TestID        *uuid.UUID
	ApplicationID *uuid.UUID
	Status        *TestRegistrationStatus
	HasResult     *bool
}

// HallTicketData represents data for generating a hall ticket.
type HallTicketData struct {
	Registration    TestRegistration
	Test            EntranceTest
	Application     AdmissionApplication
	GeneratedAt     time.Time
}

// CreateTest creates a new entrance test.
func (s *TestService) CreateTest(ctx context.Context, req CreateTestRequest) (*EntranceTest, error) {
	if req.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if req.SessionID == uuid.Nil {
		return nil, ErrSessionIDRequired
	}
	if req.TestName == "" {
		return nil, ErrTestNameRequired
	}
	if req.TestDate.IsZero() {
		return nil, ErrTestDateRequired
	}
	if req.StartTime == "" {
		return nil, ErrTestStartTimeRequired
	}

	// Check for duplicate test name in session
	var existing EntranceTest
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND session_id = ? AND test_name = ?", req.TenantID, req.SessionID, req.TestName).
		First(&existing).Error
	if err == nil {
		return nil, ErrTestAlreadyExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing test: %w", err)
	}

	// Set defaults
	if req.DurationMinutes <= 0 {
		req.DurationMinutes = 60
	}
	if req.MaxCandidates <= 0 {
		req.MaxCandidates = 100
	}
	if req.PassingPercentage.IsZero() {
		req.PassingPercentage = decimal.NewFromFloat(33.0)
	}

	test := &EntranceTest{
		TenantID:          req.TenantID,
		SessionID:         req.SessionID,
		TestName:          req.TestName,
		TestDate:          req.TestDate,
		StartTime:         req.StartTime,
		DurationMinutes:   req.DurationMinutes,
		Venue:             req.Venue,
		ClassNames:        req.ClassNames,
		MaxCandidates:     req.MaxCandidates,
		Status:            TestStatusScheduled,
		Subjects:          req.Subjects,
		Instructions:      req.Instructions,
		PassingPercentage: req.PassingPercentage,
		CreatedBy:         req.CreatedBy,
		UpdatedBy:         req.CreatedBy,
	}

	if err := s.db.WithContext(ctx).Create(test).Error; err != nil {
		return nil, fmt.Errorf("failed to create test: %w", err)
	}

	return s.GetTestByID(ctx, req.TenantID, test.ID)
}

// GetTestByID retrieves an entrance test by ID.
func (s *TestService) GetTestByID(ctx context.Context, tenantID, testID uuid.UUID) (*EntranceTest, error) {
	var test EntranceTest
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, testID).
		First(&test).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTestNotFound
		}
		return nil, fmt.Errorf("failed to get test: %w", err)
	}
	return &test, nil
}

// ListTests retrieves entrance tests with optional filtering.
func (s *TestService) ListTests(ctx context.Context, filter TestListFilter) ([]EntranceTest, error) {
	query := s.db.WithContext(ctx).
		Model(&EntranceTest{}).
		Where("tenant_id = ?", filter.TenantID)

	if filter.SessionID != nil {
		query = query.Where("session_id = ?", *filter.SessionID)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.ClassName != "" {
		query = query.Where("class_names @> ?", fmt.Sprintf(`["%s"]`, filter.ClassName))
	}
	if filter.FromDate != nil {
		query = query.Where("test_date >= ?", *filter.FromDate)
	}
	if filter.ToDate != nil {
		query = query.Where("test_date <= ?", *filter.ToDate)
	}

	query = query.Order("test_date ASC, start_time ASC")

	var tests []EntranceTest
	if err := query.Find(&tests).Error; err != nil {
		return nil, fmt.Errorf("failed to list tests: %w", err)
	}

	return tests, nil
}

// UpdateTest updates an entrance test.
func (s *TestService) UpdateTest(ctx context.Context, tenantID, testID uuid.UUID, req UpdateTestRequest) (*EntranceTest, error) {
	test, err := s.GetTestByID(ctx, tenantID, testID)
	if err != nil {
		return nil, err
	}

	// Cannot modify completed or cancelled tests
	if test.Status == TestStatusCompleted || test.Status == TestStatusCancelled {
		return nil, ErrCannotModifyCompletedTest
	}

	updates := make(map[string]interface{})

	if req.TestName != nil {
		updates["test_name"] = *req.TestName
	}
	if req.TestDate != nil {
		updates["test_date"] = *req.TestDate
	}
	if req.StartTime != nil {
		updates["start_time"] = *req.StartTime
	}
	if req.DurationMinutes != nil {
		updates["duration_minutes"] = *req.DurationMinutes
	}
	if req.Venue != nil {
		updates["venue"] = *req.Venue
	}
	if req.ClassNames != nil {
		updates["class_names"] = StringArray(req.ClassNames)
	}
	if req.MaxCandidates != nil {
		updates["max_candidates"] = *req.MaxCandidates
	}
	if req.Subjects != nil {
		updates["subjects"] = TestSubjects(req.Subjects)
	}
	if req.Instructions != nil {
		updates["instructions"] = *req.Instructions
	}
	if req.PassingPercentage != nil {
		updates["passing_percentage"] = *req.PassingPercentage
	}
	if req.Status != nil {
		if !req.Status.IsValid() {
			return nil, ErrInvalidTestStatus
		}
		updates["status"] = *req.Status
	}
	if req.UpdatedBy != nil {
		updates["updated_by"] = req.UpdatedBy
	}

	if len(updates) > 0 {
		if err := s.db.WithContext(ctx).Model(test).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("failed to update test: %w", err)
		}
	}

	return s.GetTestByID(ctx, tenantID, testID)
}

// DeleteTest deletes an entrance test.
func (s *TestService) DeleteTest(ctx context.Context, tenantID, testID uuid.UUID) error {
	test, err := s.GetTestByID(ctx, tenantID, testID)
	if err != nil {
		return err
	}

	// Cannot delete completed tests
	if test.Status == TestStatusCompleted {
		return ErrCannotModifyCompletedTest
	}

	// Check for registrations
	var count int64
	if err := s.db.WithContext(ctx).
		Model(&TestRegistration{}).
		Where("test_id = ?", testID).
		Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count registrations: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("cannot delete test with %d registrations", count)
	}

	if err := s.db.WithContext(ctx).Delete(test).Error; err != nil {
		return fmt.Errorf("failed to delete test: %w", err)
	}

	return nil
}

// RegisterCandidate registers a candidate for an entrance test.
func (s *TestService) RegisterCandidate(ctx context.Context, req RegisterCandidateRequest) (*TestRegistration, error) {
	if req.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if req.TestID == uuid.Nil {
		return nil, ErrTestNotFound
	}
	if req.ApplicationID == uuid.Nil {
		return nil, ErrApplicationNotFound
	}

	// Get the test
	test, err := s.GetTestByID(ctx, req.TenantID, req.TestID)
	if err != nil {
		return nil, err
	}

	// Validate test status
	if test.Status != TestStatusScheduled {
		return nil, ErrTestNotScheduled
	}

	// Check if test date has passed
	if test.TestDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return nil, ErrTestDatePassed
	}

	// Check max candidates
	var regCount int64
	if err := s.db.WithContext(ctx).
		Model(&TestRegistration{}).
		Where("test_id = ? AND status != ?", req.TestID, TestRegStatusCancelled).
		Count(&regCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count registrations: %w", err)
	}
	if int(regCount) >= test.MaxCandidates {
		return nil, ErrTestFull
	}

	// Check if already registered
	var existing TestRegistration
	err = s.db.WithContext(ctx).
		Where("test_id = ? AND application_id = ?", req.TestID, req.ApplicationID).
		First(&existing).Error
	if err == nil {
		return nil, ErrAlreadyRegistered
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing registration: %w", err)
	}

	// Verify application exists
	var app AdmissionApplication
	if err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ? AND deleted_at IS NULL", req.TenantID, req.ApplicationID).
		First(&app).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrApplicationNotFound
		}
		return nil, fmt.Errorf("failed to get application: %w", err)
	}

	// Generate roll number
	rollNumber, err := s.generateRollNumber(ctx, req.TenantID, req.TestID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate roll number: %w", err)
	}

	reg := &TestRegistration{
		TenantID:      req.TenantID,
		TestID:        req.TestID,
		ApplicationID: req.ApplicationID,
		RollNumber:    rollNumber,
		Status:        TestRegStatusRegistered,
		Marks:         MarksMap{},
	}

	// Create registration and update application status
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(reg).Error; err != nil {
			return fmt.Errorf("failed to create registration: %w", err)
		}

		// Update application status to test_scheduled
		if app.Status == ApplicationStatusUnderReview || app.Status == ApplicationStatusSubmitted {
			if err := tx.Model(&AdmissionApplication{}).
				Where("id = ?", req.ApplicationID).
				Update("status", ApplicationStatusTestScheduled).Error; err != nil {
				return fmt.Errorf("failed to update application status: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.GetRegistrationByID(ctx, req.TenantID, reg.ID)
}

// generateRollNumber generates the next roll number for a test.
func (s *TestService) generateRollNumber(ctx context.Context, tenantID, testID uuid.UUID) (string, error) {
	var rollNumber string
	err := s.db.WithContext(ctx).
		Raw("SELECT get_next_roll_number(?, ?)", tenantID, testID).
		Scan(&rollNumber).Error
	if err != nil {
		return "", err
	}
	return rollNumber, nil
}

// GetRegistrationByID retrieves a test registration by ID.
func (s *TestService) GetRegistrationByID(ctx context.Context, tenantID, registrationID uuid.UUID) (*TestRegistration, error) {
	var reg TestRegistration
	err := s.db.WithContext(ctx).
		Preload("Test").
		Preload("Application").
		Where("tenant_id = ? AND id = ?", tenantID, registrationID).
		First(&reg).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRegistrationNotFound
		}
		return nil, fmt.Errorf("failed to get registration: %w", err)
	}
	return &reg, nil
}

// ListRegistrations retrieves test registrations with optional filtering.
func (s *TestService) ListRegistrations(ctx context.Context, filter RegistrationListFilter) ([]TestRegistration, error) {
	query := s.db.WithContext(ctx).
		Model(&TestRegistration{}).
		Preload("Test").
		Preload("Application").
		Where("tenant_id = ?", filter.TenantID)

	if filter.TestID != nil {
		query = query.Where("test_id = ?", *filter.TestID)
	}
	if filter.ApplicationID != nil {
		query = query.Where("application_id = ?", *filter.ApplicationID)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.HasResult != nil {
		if *filter.HasResult {
			query = query.Where("result IS NOT NULL")
		} else {
			query = query.Where("result IS NULL")
		}
	}

	query = query.Order("roll_number ASC")

	var registrations []TestRegistration
	if err := query.Find(&registrations).Error; err != nil {
		return nil, fmt.Errorf("failed to list registrations: %w", err)
	}

	return registrations, nil
}

// SubmitResults submits results for a test registration.
func (s *TestService) SubmitResults(ctx context.Context, req SubmitResultsRequest) (*TestRegistration, error) {
	if req.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if req.RegistrationID == uuid.Nil {
		return nil, ErrRegistrationNotFound
	}

	// Get the registration
	reg, err := s.GetRegistrationByID(ctx, req.TenantID, req.RegistrationID)
	if err != nil {
		return nil, err
	}

	// Check if results already submitted
	if reg.Result != nil {
		return nil, ErrResultsAlreadySubmitted
	}

	// Verify registration belongs to the test
	if req.TestID != uuid.Nil && reg.TestID != req.TestID {
		return nil, ErrRegistrationNotFound
	}

	// Get the test to calculate max marks and percentage
	test, err := s.GetTestByID(ctx, req.TenantID, reg.TestID)
	if err != nil {
		return nil, err
	}

	// Calculate total marks and max marks
	totalMarks := decimal.Zero
	maxMarks := decimal.Zero
	marksMap := MarksMap{}

	for _, subject := range test.Subjects {
		maxMarks = maxMarks.Add(subject.MaxMarks)
		if marks, ok := req.Marks[subject.Subject]; ok {
			totalMarks = totalMarks.Add(marks)
			marksMap[subject.Subject] = marks
		}
	}

	// Calculate percentage
	var percentage decimal.Decimal
	if !maxMarks.IsZero() {
		percentage = totalMarks.Div(maxMarks).Mul(decimal.NewFromInt(100))
	}

	// Determine result
	result := TestResultFail
	if percentage.GreaterThanOrEqual(decimal.NewFromFloat(75)) {
		result = TestResultDistinction
	} else if percentage.GreaterThanOrEqual(decimal.NewFromFloat(60)) {
		result = TestResultMerit
	} else if percentage.GreaterThanOrEqual(test.PassingPercentage) {
		result = TestResultPass
	}

	// Update registration
	updates := map[string]interface{}{
		"marks":       marksMap,
		"total_marks": totalMarks,
		"max_marks":   maxMarks,
		"percentage":  percentage,
		"result":      result,
		"remarks":     req.Remarks,
		"status":      TestRegStatusAppeared,
	}

	if err := s.db.WithContext(ctx).Model(&TestRegistration{}).
		Where("id = ?", req.RegistrationID).
		Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update registration: %w", err)
	}

	return s.GetRegistrationByID(ctx, req.TenantID, req.RegistrationID)
}

// BulkSubmitResults submits results for multiple candidates.
func (s *TestService) BulkSubmitResults(ctx context.Context, req BulkSubmitResultsRequest) ([]TestRegistration, error) {
	if req.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if req.TestID == uuid.Nil {
		return nil, ErrTestNotFound
	}

	var results []TestRegistration
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, singleReq := range req.Results {
			reg, err := s.SubmitResults(ctx, SubmitResultsRequest{
				TenantID:       req.TenantID,
				TestID:         req.TestID,
				RegistrationID: singleReq.RegistrationID,
				Marks:          singleReq.Marks,
				Remarks:        singleReq.Remarks,
			})
			if err != nil {
				return fmt.Errorf("failed to submit results for registration %s: %w", singleReq.RegistrationID, err)
			}
			results = append(results, *reg)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return results, nil
}

// GenerateHallTicket generates hall ticket data for a registration.
func (s *TestService) GenerateHallTicket(ctx context.Context, tenantID, registrationID uuid.UUID) (*HallTicketData, error) {
	reg, err := s.GetRegistrationByID(ctx, tenantID, registrationID)
	if err != nil {
		return nil, err
	}

	// Mark hall ticket as generated
	now := time.Now()
	if reg.HallTicketGeneratedAt == nil {
		updates := map[string]interface{}{
			"hall_ticket_generated_at": now,
			"status":                   TestRegStatusHallTicketGenerated,
		}
		if err := s.db.WithContext(ctx).Model(reg).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("failed to update registration: %w", err)
		}
		reg.HallTicketGeneratedAt = &now
		reg.Status = TestRegStatusHallTicketGenerated
	}

	hallTicket := &HallTicketData{
		Registration: *reg,
		Test:         *reg.Test,
		Application:  *reg.Application,
		GeneratedAt:  now,
	}

	return hallTicket, nil
}

// GetHallTickets retrieves hall ticket data for all registrations of a test.
func (s *TestService) GetHallTickets(ctx context.Context, tenantID, testID uuid.UUID) ([]HallTicketData, error) {
	registrations, err := s.ListRegistrations(ctx, RegistrationListFilter{
		TenantID: tenantID,
		TestID:   &testID,
	})
	if err != nil {
		return nil, err
	}

	var hallTickets []HallTicketData
	for _, reg := range registrations {
		if reg.Status == TestRegStatusCancelled {
			continue
		}
		hallTickets = append(hallTickets, HallTicketData{
			Registration: reg,
			Test:         *reg.Test,
			Application:  *reg.Application,
			GeneratedAt:  time.Now(),
		})
	}

	return hallTickets, nil
}

// CountRegistrations counts registrations for a test.
func (s *TestService) CountRegistrations(ctx context.Context, tenantID, testID uuid.UUID) (int64, error) {
	var count int64
	err := s.db.WithContext(ctx).
		Model(&TestRegistration{}).
		Where("tenant_id = ? AND test_id = ? AND status != ?", tenantID, testID, TestRegStatusCancelled).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count registrations: %w", err)
	}
	return count, nil
}

// CancelRegistration cancels a test registration.
func (s *TestService) CancelRegistration(ctx context.Context, tenantID, registrationID uuid.UUID) error {
	reg, err := s.GetRegistrationByID(ctx, tenantID, registrationID)
	if err != nil {
		return err
	}

	if reg.Result != nil {
		return fmt.Errorf("cannot cancel registration with results")
	}

	if err := s.db.WithContext(ctx).Model(reg).Update("status", TestRegStatusCancelled).Error; err != nil {
		return fmt.Errorf("failed to cancel registration: %w", err)
	}

	return nil
}
