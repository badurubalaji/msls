// Package admission provides admission management services.
package admission

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
)

// ApplicationService handles application-related operations.
type ApplicationService struct {
	db *gorm.DB
}

// NewApplicationService creates a new ApplicationService instance.
func NewApplicationService(db *gorm.DB) *ApplicationService {
	return &ApplicationService{db: db}
}

// CreateApplicationRequest represents a request to create an application.
type CreateApplicationRequest struct {
	TenantID         uuid.UUID
	BranchID         *uuid.UUID
	SessionID        uuid.UUID
	EnquiryID        *uuid.UUID
	StudentName      string
	ClassApplying    string
	DateOfBirth      *time.Time
	Gender           string
	BloodGroup       string
	Nationality      string
	Religion         string
	Category         string
	AadharNumber     string
	AddressLine1     string
	AddressLine2     string
	City             string
	State            string
	PostalCode       string
	PreviousSchool   string
	PreviousClass    string
	FatherName       string
	FatherPhone      string
	FatherEmail      string
	FatherOccupation string
	MotherName       string
	MotherPhone      string
	MotherEmail      string
	MotherOccupation string
	GuardianName     string
	GuardianPhone    string
	GuardianEmail    string
	GuardianRelation string
	CreatedBy        *uuid.UUID
}

// UpdateApplicationRequest represents a request to update an application.
type UpdateApplicationRequest struct {
	StudentName      *string
	ClassApplying    *string
	DateOfBirth      *time.Time
	Gender           *string
	BloodGroup       *string
	Nationality      *string
	Religion         *string
	Category         *string
	AadharNumber     *string
	AddressLine1     *string
	AddressLine2     *string
	City             *string
	State            *string
	PostalCode       *string
	PreviousSchool   *string
	PreviousClass    *string
	PreviousPercentage *float64
	FatherName       *string
	FatherPhone      *string
	FatherEmail      *string
	FatherOccupation *string
	MotherName       *string
	MotherPhone      *string
	MotherEmail      *string
	MotherOccupation *string
	GuardianName     *string
	GuardianPhone    *string
	GuardianEmail    *string
	GuardianRelation *string
	UpdatedBy        *uuid.UUID
}

// ListApplicationFilter contains filters for listing applications.
type ListApplicationFilter struct {
	TenantID        uuid.UUID
	BranchID        *uuid.UUID
	SessionID       *uuid.UUID
	Status          *models.ApplicationStatus
	ClassApplying   *string
	Search          string
	SubmittedAfter  *time.Time
	SubmittedBefore *time.Time
}

// UpdateStageRequest represents a request to update the application stage.
type UpdateStageRequest struct {
	NewStage  models.ApplicationStatus
	Remarks   string
	ChangedBy *uuid.UUID
}

// AddParentRequest represents a request to add a parent.
type AddParentRequest struct {
	TenantID      uuid.UUID
	ApplicationID uuid.UUID
	Relation      models.ParentRelation
	Name          string
	Phone         string
	Email         string
	Occupation    string
	Education     string
	AnnualIncome  string
}

// UpdateParentRequest represents a request to update a parent.
type UpdateParentRequest struct {
	Relation     *models.ParentRelation
	Name         *string
	Phone        *string
	Email        *string
	Occupation   *string
	Education    *string
	AnnualIncome *string
}

// AddDocumentRequest represents a request to add a document.
type AddDocumentRequest struct {
	TenantID      uuid.UUID
	ApplicationID uuid.UUID
	DocumentType  models.DocumentType
	FileURL       string
	FileName      string
	FileSize      int64
	MimeType      string
}

// DocumentVerificationInput represents input for simple document verification.
type DocumentVerificationInput struct {
	VerificationStatus  models.VerificationStatus
	VerificationRemarks string
	VerifiedBy          *uuid.UUID
}

// StatusCheckRequest represents a request to check application status.
type StatusCheckRequest struct {
	ApplicationNumber string
	Phone             string
}

// StatusCheckResponse represents the response for status check.
type StatusCheckResponse struct {
	ApplicationNumber string                   `json:"applicationNumber"`
	StudentName       string                   `json:"studentName"`
	ClassApplying     string                   `json:"classApplying"`
	Status            models.ApplicationStatus `json:"status"`
	SubmittedAt       *time.Time               `json:"submittedAt,omitempty"`
	ReviewNotes       string                   `json:"reviewNotes,omitempty"`
}

// Create creates a new admission application.
func (s *ApplicationService) Create(ctx context.Context, req CreateApplicationRequest) (*models.AdmissionApplication, error) {
	if req.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if req.SessionID == uuid.Nil {
		return nil, ErrSessionIDRequired
	}
	if req.StudentName == "" {
		return nil, ErrStudentNameRequired
	}
	if req.ClassApplying == "" {
		return nil, ErrClassApplyingRequired
	}

	// Validate session exists and is open
	var session models.AdmissionSession
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", req.TenantID, req.SessionID).
		First(&session).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to check session: %w", err)
	}
	if session.Status == models.SessionStatusClosed {
		return nil, ErrSessionClosed
	}

	// Generate application number
	appNumber, err := s.generateApplicationNumber(ctx, req.TenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate application number: %w", err)
	}

	application := &models.AdmissionApplication{
		TenantID:          req.TenantID,
		BranchID:          req.BranchID,
		SessionID:         req.SessionID,
		EnquiryID:         req.EnquiryID,
		ApplicationNumber: appNumber,
		StudentName:       req.StudentName,
		ClassApplying:     req.ClassApplying,
		DateOfBirth:       req.DateOfBirth,
		Gender:            req.Gender,
		BloodGroup:        req.BloodGroup,
		Nationality:       req.Nationality,
		Religion:          req.Religion,
		Category:          req.Category,
		AadharNumber:      req.AadharNumber,
		AddressLine1:      req.AddressLine1,
		AddressLine2:      req.AddressLine2,
		City:              req.City,
		State:             req.State,
		PostalCode:        req.PostalCode,
		PreviousSchool:    req.PreviousSchool,
		PreviousClass:     req.PreviousClass,
		FatherName:        req.FatherName,
		FatherPhone:       req.FatherPhone,
		FatherEmail:       req.FatherEmail,
		FatherOccupation:  req.FatherOccupation,
		MotherName:        req.MotherName,
		MotherPhone:       req.MotherPhone,
		MotherEmail:       req.MotherEmail,
		MotherOccupation:  req.MotherOccupation,
		GuardianName:      req.GuardianName,
		GuardianPhone:     req.GuardianPhone,
		GuardianEmail:     req.GuardianEmail,
		GuardianRelation:  req.GuardianRelation,
		Status:            models.ApplicationStatusDraft,
		CreatedBy:         req.CreatedBy,
		UpdatedBy:         req.CreatedBy,
	}

	if err := s.db.WithContext(ctx).Create(application).Error; err != nil {
		return nil, fmt.Errorf("failed to create application: %w", err)
	}

	return s.GetByID(ctx, req.TenantID, application.ID)
}

// GetByID retrieves an application by ID.
func (s *ApplicationService) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*models.AdmissionApplication, error) {
	var application models.AdmissionApplication
	err := s.db.WithContext(ctx).
		Preload("Session").
		Preload("Parents").
		Preload("Documents").
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&application).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrApplicationNotFound
		}
		return nil, fmt.Errorf("failed to get application: %w", err)
	}
	return &application, nil
}

// GetByNumber retrieves an application by application number.
func (s *ApplicationService) GetByNumber(ctx context.Context, tenantID uuid.UUID, appNumber string) (*models.AdmissionApplication, error) {
	var application models.AdmissionApplication
	err := s.db.WithContext(ctx).
		Preload("Session").
		Preload("Parents").
		Preload("Documents").
		Where("tenant_id = ? AND application_number = ?", tenantID, appNumber).
		First(&application).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrApplicationNotFound
		}
		return nil, fmt.Errorf("failed to get application: %w", err)
	}
	return &application, nil
}

// List retrieves applications with optional filtering.
func (s *ApplicationService) List(ctx context.Context, filter ListApplicationFilter) ([]models.AdmissionApplication, error) {
	query := s.db.WithContext(ctx).
		Model(&models.AdmissionApplication{}).
		Preload("Session").
		Where("tenant_id = ?", filter.TenantID).
		Order("created_at DESC")

	if filter.BranchID != nil {
		query = query.Where("branch_id = ?", *filter.BranchID)
	}
	if filter.SessionID != nil {
		query = query.Where("session_id = ?", *filter.SessionID)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.ClassApplying != nil {
		query = query.Where("class_applying = ?", *filter.ClassApplying)
	}
	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("student_name ILIKE ? OR application_number ILIKE ?", search, search)
	}
	if filter.SubmittedAfter != nil {
		query = query.Where("submitted_at >= ?", *filter.SubmittedAfter)
	}
	if filter.SubmittedBefore != nil {
		query = query.Where("submitted_at <= ?", *filter.SubmittedBefore)
	}

	var applications []models.AdmissionApplication
	if err := query.Find(&applications).Error; err != nil {
		return nil, fmt.Errorf("failed to list applications: %w", err)
	}

	return applications, nil
}

// Update updates an application.
func (s *ApplicationService) Update(ctx context.Context, tenantID, id uuid.UUID, req UpdateApplicationRequest) (*models.AdmissionApplication, error) {
	application, err := s.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	// Only allow updates for applications that are not in terminal statuses
	terminalStatuses := map[models.ApplicationStatus]bool{
		models.ApplicationStatusApproved:  true,
		models.ApplicationStatusRejected:  true,
		models.ApplicationStatusEnrolled:  true,
		models.ApplicationStatusWithdrawn: true,
	}
	if terminalStatuses[application.Status] {
		return nil, ErrCannotUpdateApplication
	}

	updates := make(map[string]interface{})

	if req.StudentName != nil {
		updates["student_name"] = *req.StudentName
	}
	if req.ClassApplying != nil {
		updates["class_applying"] = *req.ClassApplying
	}
	if req.DateOfBirth != nil {
		updates["date_of_birth"] = req.DateOfBirth
	}
	if req.Gender != nil {
		updates["gender"] = *req.Gender
	}
	if req.BloodGroup != nil {
		updates["blood_group"] = *req.BloodGroup
	}
	if req.Nationality != nil {
		updates["nationality"] = *req.Nationality
	}
	if req.Religion != nil {
		updates["religion"] = *req.Religion
	}
	if req.Category != nil {
		updates["category"] = *req.Category
	}
	if req.AadharNumber != nil {
		updates["aadhar_number"] = *req.AadharNumber
	}
	if req.AddressLine1 != nil {
		updates["address_line1"] = *req.AddressLine1
	}
	if req.AddressLine2 != nil {
		updates["address_line2"] = *req.AddressLine2
	}
	if req.City != nil {
		updates["city"] = *req.City
	}
	if req.State != nil {
		updates["state"] = *req.State
	}
	if req.PostalCode != nil {
		updates["postal_code"] = *req.PostalCode
	}
	if req.PreviousSchool != nil {
		updates["previous_school"] = *req.PreviousSchool
	}
	if req.PreviousClass != nil {
		updates["previous_class"] = *req.PreviousClass
	}
	if req.PreviousPercentage != nil {
		updates["previous_percentage"] = *req.PreviousPercentage
	}
	if req.FatherName != nil {
		updates["father_name"] = *req.FatherName
	}
	if req.FatherPhone != nil {
		updates["father_phone"] = *req.FatherPhone
	}
	if req.FatherEmail != nil {
		updates["father_email"] = *req.FatherEmail
	}
	if req.FatherOccupation != nil {
		updates["father_occupation"] = *req.FatherOccupation
	}
	if req.MotherName != nil {
		updates["mother_name"] = *req.MotherName
	}
	if req.MotherPhone != nil {
		updates["mother_phone"] = *req.MotherPhone
	}
	if req.MotherEmail != nil {
		updates["mother_email"] = *req.MotherEmail
	}
	if req.MotherOccupation != nil {
		updates["mother_occupation"] = *req.MotherOccupation
	}
	if req.GuardianName != nil {
		updates["guardian_name"] = *req.GuardianName
	}
	if req.GuardianPhone != nil {
		updates["guardian_phone"] = *req.GuardianPhone
	}
	if req.GuardianEmail != nil {
		updates["guardian_email"] = *req.GuardianEmail
	}
	if req.GuardianRelation != nil {
		updates["guardian_relation"] = *req.GuardianRelation
	}
	if req.UpdatedBy != nil {
		updates["updated_by"] = req.UpdatedBy
	}

	if len(updates) > 0 {
		if err := s.db.WithContext(ctx).Model(application).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("failed to update application: %w", err)
		}
	}

	return s.GetByID(ctx, tenantID, id)
}

// Submit submits an application.
func (s *ApplicationService) Submit(ctx context.Context, tenantID, id uuid.UUID, submittedBy *uuid.UUID) (*models.AdmissionApplication, error) {
	application, err := s.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	// Allow submission from draft or re-submission from certain statuses
	allowedStatuses := map[models.ApplicationStatus]bool{
		models.ApplicationStatusDraft:            true,
		models.ApplicationStatusSubmitted:        true, // Allow re-submission
		models.ApplicationStatusDocumentsPending: true, // After fixing documents
		models.ApplicationStatusUnderReview:      true, // If sent back for corrections
	}
	if !allowedStatuses[application.Status] {
		return nil, ErrApplicationNotInDraft
	}

	// Validate required fields
	if application.StudentName == "" || application.ClassApplying == "" {
		return nil, ErrMissingRequiredFields
	}

	now := time.Now()
	updates := map[string]interface{}{
		"status":       models.ApplicationStatusSubmitted,
		"submitted_at": now,
		"updated_by":   submittedBy,
	}

	// Create stage history entry
	historyEntry := models.StageHistoryEntry{
		Stage:     models.ApplicationStatusSubmitted,
		Timestamp: now,
		ChangedBy: submittedBy,
		Remarks:   "Application submitted",
	}

	// Get current history and append
	var currentHistory models.StageHistory
	if err := s.db.WithContext(ctx).Model(application).Select("stage_history").Scan(&currentHistory).Error; err == nil {
		currentHistory = append(currentHistory, historyEntry)
		updates["stage_history"] = currentHistory
	}

	if err := s.db.WithContext(ctx).Model(application).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to submit application: %w", err)
	}

	return s.GetByID(ctx, tenantID, id)
}

// UpdateStage updates the stage of an application.
func (s *ApplicationService) UpdateStage(ctx context.Context, tenantID, id uuid.UUID, req UpdateStageRequest) (*models.AdmissionApplication, error) {
	application, err := s.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	// Validate stage transition
	if !application.Status.CanTransitionTo(req.NewStage) {
		validTransitions := application.Status.GetValidTransitions()
		validStrings := make([]string, len(validTransitions))
		for i, v := range validTransitions {
			validStrings[i] = string(v)
		}
		return nil, NewStageTransitionError(
			string(application.Status),
			string(req.NewStage),
			validStrings,
		)
	}

	// Use transaction for enrollment to ensure atomicity
	if req.NewStage == models.ApplicationStatusEnrolled {
		return s.enrollStudent(ctx, tenantID, application, req)
	}

	updates := map[string]interface{}{
		"status":     req.NewStage,
		"updated_by": req.ChangedBy,
	}

	// Add remarks if provided (for rejections or other notes)
	if req.Remarks != "" {
		updates["remarks"] = req.Remarks
	}

	if err := s.db.WithContext(ctx).Model(application).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update stage: %w", err)
	}

	return s.GetByID(ctx, tenantID, id)
}

// enrollStudent creates a student record from an application and updates the application status to enrolled.
func (s *ApplicationService) enrollStudent(ctx context.Context, tenantID uuid.UUID, application *models.AdmissionApplication, req UpdateStageRequest) (*models.AdmissionApplication, error) {
	// Start transaction
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get primary branch for the tenant (or use application's branch if set)
	var branchID uuid.UUID
	if application.BranchID != nil {
		branchID = *application.BranchID
	} else {
		// Get primary branch
		var branch models.Branch
		if err := tx.Where("tenant_id = ? AND is_primary = ?", tenantID, true).First(&branch).Error; err != nil {
			// Fallback: get any branch
			if err := tx.Where("tenant_id = ?", tenantID).First(&branch).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("no branch found for tenant: %w", err)
			}
		}
		branchID = branch.ID
	}

	// Generate admission number
	admissionNumber, err := s.generateAdmissionNumber(tx, tenantID, branchID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to generate admission number: %w", err)
	}

	// Parse student name
	firstName, lastName := parseStudentName(application.StudentName)

	// Parse date of birth
	var dob time.Time
	if application.DateOfBirth != nil {
		dob = *application.DateOfBirth
	}

	// Map gender
	gender := models.GenderMale
	if application.Gender != "" {
		switch application.Gender {
		case "female":
			gender = models.GenderFemale
		case "other":
			gender = models.GenderOther
		}
	}

	// Create student record
	student := &models.Student{
		TenantID:        tenantID,
		BranchID:        branchID,
		AdmissionNumber: admissionNumber,
		FirstName:       firstName,
		LastName:        lastName,
		DateOfBirth:     dob,
		Gender:          gender,
		BloodGroup:      application.BloodGroup,
		AadhaarNumber:   application.AadharNumber,
		Status:          models.StudentStatusActive,
		AdmissionDate:   time.Now(),
		CreatedBy:       req.ChangedBy,
	}

	if err := tx.Create(student).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create student: %w", err)
	}

	// Create student address if available
	if application.AddressLine1 != "" {
		address := &models.StudentAddress{
			TenantID:     tenantID,
			StudentID:    student.ID,
			AddressType:  models.AddressTypeCurrent,
			AddressLine1: application.AddressLine1,
			AddressLine2: application.AddressLine2,
			City:         application.City,
			State:        application.State,
			PostalCode:   application.PostalCode,
		}
		if err := tx.Create(address).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create student address: %w", err)
		}
	}

	// Update application status
	updates := map[string]interface{}{
		"status":     models.ApplicationStatusEnrolled,
		"updated_by": req.ChangedBy,
	}
	if req.Remarks != "" {
		updates["remarks"] = req.Remarks
	}

	if err := tx.Model(application).Updates(updates).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update application status: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return s.GetByID(ctx, tenantID, application.ID)
}

// generateAdmissionNumber generates a unique admission number for a student.
func (s *ApplicationService) generateAdmissionNumber(tx *gorm.DB, tenantID, branchID uuid.UUID) (string, error) {
	year := time.Now().Year()

	// Get or create sequence
	var sequence models.StudentAdmissionSequence
	err := tx.Where("tenant_id = ? AND branch_id = ? AND year = ?", tenantID, branchID, year).First(&sequence).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new sequence
			sequence = models.StudentAdmissionSequence{
				TenantID:     tenantID,
				BranchID:     branchID,
				Year:         year,
				LastSequence: 0,
			}
			if err := tx.Create(&sequence).Error; err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	// Increment sequence
	sequence.LastSequence++
	if err := tx.Save(&sequence).Error; err != nil {
		return "", err
	}

	// Format: ADM-YYYY-NNNNN
	return fmt.Sprintf("ADM-%d-%05d", year, sequence.LastSequence), nil
}

// parseStudentName splits a full name into first and last name.
func parseStudentName(fullName string) (firstName, lastName string) {
	parts := regexp.MustCompile(`\s+`).Split(fullName, -1)
	if len(parts) == 0 {
		return "", ""
	}
	firstName = parts[0]
	if len(parts) > 1 {
		lastName = parts[len(parts)-1]
	}
	return
}

// Delete deletes an application (only draft applications can be deleted).
func (s *ApplicationService) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	application, err := s.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}

	if application.Status != models.ApplicationStatusDraft {
		return ErrCannotDeleteSubmittedApplication
	}

	// Delete in a transaction: parents, documents, then application
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete parents
		if err := tx.Where("application_id = ?", id).Delete(&models.ApplicationParent{}).Error; err != nil {
			return fmt.Errorf("failed to delete parents: %w", err)
		}

		// Delete documents
		if err := tx.Where("application_id = ?", id).Delete(&models.ApplicationDocument{}).Error; err != nil {
			return fmt.Errorf("failed to delete documents: %w", err)
		}

		// Delete application
		if err := tx.Delete(application).Error; err != nil {
			return fmt.Errorf("failed to delete application: %w", err)
		}

		return nil
	})

	return err
}

// AddParent adds a parent to an application.
func (s *ApplicationService) AddParent(ctx context.Context, req AddParentRequest) (*models.ApplicationParent, error) {
	// Validate application exists
	_, err := s.GetByID(ctx, req.TenantID, req.ApplicationID)
	if err != nil {
		return nil, err
	}

	// Validate relation
	if !req.Relation.IsValid() {
		return nil, ErrInvalidParentRelation
	}

	parent := &models.ApplicationParent{
		TenantID:      req.TenantID,
		ApplicationID: req.ApplicationID,
		Relation:      req.Relation,
		Name:          req.Name,
		Phone:         req.Phone,
		Email:         req.Email,
		Occupation:    req.Occupation,
		Education:     req.Education,
		AnnualIncome:  req.AnnualIncome,
	}

	if err := s.db.WithContext(ctx).Create(parent).Error; err != nil {
		return nil, fmt.Errorf("failed to add parent: %w", err)
	}

	return parent, nil
}

// UpdateParent updates a parent record.
func (s *ApplicationService) UpdateParent(ctx context.Context, tenantID, applicationID, parentID uuid.UUID, req UpdateParentRequest) (*models.ApplicationParent, error) {
	var parent models.ApplicationParent
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND application_id = ? AND id = ?", tenantID, applicationID, parentID).
		First(&parent).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrParentNotFound
		}
		return nil, fmt.Errorf("failed to get parent: %w", err)
	}

	updates := make(map[string]interface{})

	if req.Relation != nil {
		if !req.Relation.IsValid() {
			return nil, ErrInvalidParentRelation
		}
		updates["relation"] = *req.Relation
	}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.Occupation != nil {
		updates["occupation"] = *req.Occupation
	}
	if req.Education != nil {
		updates["education"] = *req.Education
	}
	if req.AnnualIncome != nil {
		updates["annual_income"] = *req.AnnualIncome
	}

	if len(updates) > 0 {
		if err := s.db.WithContext(ctx).Model(&parent).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("failed to update parent: %w", err)
		}
	}

	// Re-fetch parent
	if err := s.db.WithContext(ctx).First(&parent, parent.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to get updated parent: %w", err)
	}

	return &parent, nil
}

// GetParents retrieves all parents for an application.
func (s *ApplicationService) GetParents(ctx context.Context, tenantID, applicationID uuid.UUID) ([]models.ApplicationParent, error) {
	var parents []models.ApplicationParent
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND application_id = ?", tenantID, applicationID).
		Order("created_at").
		Find(&parents).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get parents: %w", err)
	}
	return parents, nil
}

// DeleteParent deletes a parent record.
func (s *ApplicationService) DeleteParent(ctx context.Context, tenantID, applicationID, parentID uuid.UUID) error {
	result := s.db.WithContext(ctx).
		Where("tenant_id = ? AND application_id = ? AND id = ?", tenantID, applicationID, parentID).
		Delete(&models.ApplicationParent{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete parent: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrParentNotFound
	}
	return nil
}

// AddDocument adds a document to an application.
func (s *ApplicationService) AddDocument(ctx context.Context, req AddDocumentRequest) (*models.ApplicationDocument, error) {
	// Validate application exists
	_, err := s.GetByID(ctx, req.TenantID, req.ApplicationID)
	if err != nil {
		return nil, err
	}

	// Validate document type
	if !req.DocumentType.IsValid() {
		return nil, ErrInvalidDocumentType
	}

	document := &models.ApplicationDocument{
		TenantID:      req.TenantID,
		ApplicationID: req.ApplicationID,
		DocumentType:  req.DocumentType,
		FileURL:       req.FileURL,
		FileName:      req.FileName,
		FileSize:      req.FileSize,
		MimeType:      req.MimeType,
	}

	if err := s.db.WithContext(ctx).Create(document).Error; err != nil {
		return nil, fmt.Errorf("failed to add document: %w", err)
	}

	return document, nil
}

// GetDocuments retrieves all documents for an application.
func (s *ApplicationService) GetDocuments(ctx context.Context, tenantID, applicationID uuid.UUID) ([]models.ApplicationDocument, error) {
	var documents []models.ApplicationDocument
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND application_id = ?", tenantID, applicationID).
		Order("created_at").
		Find(&documents).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get documents: %w", err)
	}
	return documents, nil
}

// GetDocument retrieves a specific document.
func (s *ApplicationService) GetDocument(ctx context.Context, tenantID, applicationID, documentID uuid.UUID) (*models.ApplicationDocument, error) {
	var document models.ApplicationDocument
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND application_id = ? AND id = ?", tenantID, applicationID, documentID).
		First(&document).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDocumentNotFound
		}
		return nil, fmt.Errorf("failed to get document: %w", err)
	}
	return &document, nil
}

// DeleteDocument deletes a document.
func (s *ApplicationService) DeleteDocument(ctx context.Context, tenantID, applicationID, documentID uuid.UUID) error {
	result := s.db.WithContext(ctx).
		Where("tenant_id = ? AND application_id = ? AND id = ?", tenantID, applicationID, documentID).
		Delete(&models.ApplicationDocument{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete document: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrDocumentNotFound
	}
	return nil
}

// VerifyDocument verifies or rejects a document.
func (s *ApplicationService) VerifyDocument(ctx context.Context, tenantID, applicationID, documentID uuid.UUID, req DocumentVerificationInput) (*models.ApplicationDocument, error) {
	document, err := s.GetDocument(ctx, tenantID, applicationID, documentID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	updates := map[string]interface{}{
		"verification_status":  req.VerificationStatus,
		"verified_by":          req.VerifiedBy,
		"verified_at":          now,
		"verification_remarks": req.VerificationRemarks,
	}

	if err := s.db.WithContext(ctx).Model(document).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to verify document: %w", err)
	}

	return s.GetDocument(ctx, tenantID, applicationID, documentID)
}

// CheckStatus checks the status of an application (public API).
func (s *ApplicationService) CheckStatus(ctx context.Context, tenantID uuid.UUID, req StatusCheckRequest) (*StatusCheckResponse, error) {
	// Validate phone number format (basic validation)
	phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{9,14}$`)
	if !phoneRegex.MatchString(req.Phone) {
		return nil, ErrInvalidPhoneNumber
	}

	var application models.AdmissionApplication
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND application_number = ?", tenantID, req.ApplicationNumber).
		First(&application).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrApplicationNotFound
		}
		return nil, fmt.Errorf("failed to get application: %w", err)
	}

	// Verify phone number matches (check in parent info)
	phoneMatches := false
	if application.FatherPhone == req.Phone ||
		application.MotherPhone == req.Phone ||
		application.GuardianPhone == req.Phone {
		phoneMatches = true
	}

	// Also check in separate parents table
	if !phoneMatches {
		var count int64
		s.db.WithContext(ctx).
			Model(&models.ApplicationParent{}).
			Where("application_id = ? AND phone = ?", application.ID, req.Phone).
			Count(&count)
		phoneMatches = count > 0
	}

	if !phoneMatches {
		return nil, ErrApplicationNotFound // Don't reveal that application exists
	}

	return &StatusCheckResponse{
		ApplicationNumber: application.ApplicationNumber,
		StudentName:       application.StudentName,
		ClassApplying:     application.ClassApplying,
		Status:            application.Status,
		SubmittedAt:       application.SubmittedAt,
		ReviewNotes:       application.InternalNotes,
	}, nil
}

// Count returns the total number of applications for a tenant with optional filters.
func (s *ApplicationService) Count(ctx context.Context, filter ListApplicationFilter) (int64, error) {
	query := s.db.WithContext(ctx).
		Model(&models.AdmissionApplication{}).
		Where("tenant_id = ?", filter.TenantID)

	if filter.SessionID != nil {
		query = query.Where("session_id = ?", *filter.SessionID)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count applications: %w", err)
	}
	return count, nil
}

// generateApplicationNumber generates a unique application number.
// Format: APP-YYYYMMDD-XXXX (e.g., APP-20260123-0001)
func (s *ApplicationService) generateApplicationNumber(ctx context.Context, tenantID uuid.UUID) (string, error) {
	datePrefix := time.Now().Format("20060102")

	var sequence models.ApplicationNumberSequence
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND date_prefix = ?", tenantID, datePrefix).
		First(&sequence).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new sequence for today
		sequence = models.ApplicationNumberSequence{
			TenantID:     tenantID,
			DatePrefix:   datePrefix,
			LastSequence: 0,
		}
	} else if err != nil {
		return "", fmt.Errorf("failed to get sequence: %w", err)
	}

	// Increment sequence
	sequence.LastSequence++
	newSequence := sequence.LastSequence

	// Upsert the sequence
	err = s.db.WithContext(ctx).
		Where("tenant_id = ? AND date_prefix = ?", tenantID, datePrefix).
		Assign(map[string]interface{}{"last_sequence": newSequence}).
		FirstOrCreate(&sequence).Error
	if err != nil {
		return "", fmt.Errorf("failed to update sequence: %w", err)
	}

	return fmt.Sprintf("APP-%s-%04d", datePrefix, newSequence), nil
}
