// Package admission provides admission management services including enquiries.
package admission

import (
	"context"
	"errors"
	"fmt"
	"time"

	"msls-backend/internal/pkg/database/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EnquirySource represents the source of an admission enquiry.
type EnquirySource string

// EnquirySource constants.
const (
	EnquirySourceWalkIn        EnquirySource = "walk_in"
	EnquirySourcePhone         EnquirySource = "phone"
	EnquirySourceWebsite       EnquirySource = "website"
	EnquirySourceReferral      EnquirySource = "referral"
	EnquirySourceAdvertisement EnquirySource = "advertisement"
	EnquirySourceSocialMedia   EnquirySource = "social_media"
	EnquirySourceOther         EnquirySource = "other"
)

// IsValid checks if the source is a valid value.
func (s EnquirySource) IsValid() bool {
	switch s {
	case EnquirySourceWalkIn, EnquirySourcePhone, EnquirySourceWebsite,
		EnquirySourceReferral, EnquirySourceAdvertisement,
		EnquirySourceSocialMedia, EnquirySourceOther:
		return true
	}
	return false
}

// EnquiryStatus represents the status of an admission enquiry.
type EnquiryStatus string

// EnquiryStatus constants.
const (
	EnquiryStatusNew       EnquiryStatus = "new"
	EnquiryStatusContacted EnquiryStatus = "contacted"
	EnquiryStatusIntersted EnquiryStatus = "interested"
	EnquiryStatusConverted EnquiryStatus = "converted"
	EnquiryStatusClosed    EnquiryStatus = "closed"
)

// IsValid checks if the status is a valid value.
func (s EnquiryStatus) IsValid() bool {
	switch s {
	case EnquiryStatusNew, EnquiryStatusContacted, EnquiryStatusIntersted,
		EnquiryStatusConverted, EnquiryStatusClosed:
		return true
	}
	return false
}

// ContactMode represents the mode of contact for follow-ups.
type ContactMode string

// ContactMode constants.
const (
	ContactModePhone    ContactMode = "phone"
	ContactModeEmail    ContactMode = "email"
	ContactModeWhatsApp ContactMode = "whatsapp"
	ContactModeInPerson ContactMode = "in_person"
	ContactModeOther    ContactMode = "other"
)

// IsValid checks if the contact mode is a valid value.
func (m ContactMode) IsValid() bool {
	switch m {
	case ContactModePhone, ContactModeEmail, ContactModeWhatsApp,
		ContactModeInPerson, ContactModeOther:
		return true
	}
	return false
}

// FollowUpOutcome represents the outcome of a follow-up.
type FollowUpOutcome string

// FollowUpOutcome constants.
const (
	FollowUpOutcomeInterested      FollowUpOutcome = "interested"
	FollowUpOutcomeNotInterested   FollowUpOutcome = "not_interested"
	FollowUpOutcomeFollowUpReq     FollowUpOutcome = "follow_up_required"
	FollowUpOutcomeConverted       FollowUpOutcome = "converted"
	FollowUpOutcomeNoResponse      FollowUpOutcome = "no_response"
)

// IsValid checks if the outcome is a valid value.
func (o FollowUpOutcome) IsValid() bool {
	switch o {
	case FollowUpOutcomeInterested, FollowUpOutcomeNotInterested,
		FollowUpOutcomeFollowUpReq, FollowUpOutcomeConverted,
		FollowUpOutcomeNoResponse:
		return true
	}
	return false
}

// Gender represents gender values.
type Gender string

// Gender constants.
const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

// IsValid checks if the gender is a valid value.
func (g Gender) IsValid() bool {
	switch g {
	case GenderMale, GenderFemale, GenderOther, "":
		return true
	}
	return false
}

// AdmissionEnquiry represents an admission enquiry in the database.
type AdmissionEnquiry struct {
	ID                     uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()"`
	TenantID               uuid.UUID      `gorm:"type:uuid;not null"`
	BranchID               *uuid.UUID     `gorm:"type:uuid"`
	SessionID              *uuid.UUID     `gorm:"type:uuid"`
	EnquiryNumber          string         `gorm:"type:varchar(50);not null"`
	StudentName            string         `gorm:"type:varchar(200);not null"`
	DateOfBirth            *time.Time     `gorm:"type:date"`
	Gender                 string         `gorm:"type:varchar(20)"`
	ClassApplying          string         `gorm:"type:varchar(50);not null"`
	ParentName             string         `gorm:"type:varchar(200);not null"`
	ParentPhone            string         `gorm:"type:varchar(20);not null"`
	ParentEmail            string         `gorm:"type:varchar(255)"`
	Source                 EnquirySource  `gorm:"type:enquiry_source;not null;default:'walk_in'"`
	ReferralDetails        string         `gorm:"type:text"`
	Remarks                string         `gorm:"type:text"`
	Status                 EnquiryStatus  `gorm:"type:enquiry_status;not null;default:'new'"`
	FollowUpDate           *time.Time     `gorm:"type:date"`
	AssignedTo             *uuid.UUID     `gorm:"type:uuid"`
	ConvertedApplicationID *uuid.UUID     `gorm:"type:uuid"`
	CreatedAt              time.Time      `gorm:"not null;default:now()"`
	UpdatedAt              time.Time      `gorm:"not null;default:now()"`
	CreatedBy              *uuid.UUID     `gorm:"type:uuid"`
	UpdatedBy              *uuid.UUID     `gorm:"type:uuid"`
	DeletedAt              gorm.DeletedAt `gorm:"index"`
}

// TableName returns the table name for the model.
func (AdmissionEnquiry) TableName() string {
	return "admission_enquiries"
}

// EnquiryFollowUp represents a follow-up entry for an enquiry.
type EnquiryFollowUp struct {
	ID           uuid.UUID        `gorm:"type:uuid;primaryKey;default:uuid_generate_v7()"`
	TenantID     uuid.UUID        `gorm:"type:uuid;not null"`
	EnquiryID    uuid.UUID        `gorm:"type:uuid;not null"`
	FollowUpDate time.Time        `gorm:"type:date;not null"`
	ContactMode  ContactMode      `gorm:"type:contact_mode;not null;default:'phone'"`
	Notes        string           `gorm:"type:text"`
	Outcome      *FollowUpOutcome `gorm:"type:follow_up_outcome"`
	NextFollowUp *time.Time       `gorm:"type:date"`
	CreatedAt    time.Time        `gorm:"not null;default:now()"`
	CreatedBy    *uuid.UUID       `gorm:"type:uuid"`
}

// TableName returns the table name for the model.
func (EnquiryFollowUp) TableName() string {
	return "enquiry_follow_ups"
}

// EnquiryService handles enquiry-related operations.
type EnquiryService struct {
	db *gorm.DB
}

// NewEnquiryService creates a new EnquiryService instance.
func NewEnquiryService(db *gorm.DB) *EnquiryService {
	return &EnquiryService{db: db}
}

// CreateEnquiryRequest represents a request to create an enquiry.
type CreateEnquiryRequest struct {
	TenantID        uuid.UUID
	BranchID        *uuid.UUID
	SessionID       *uuid.UUID
	StudentName     string
	DateOfBirth     *time.Time
	Gender          string
	ClassApplying   string
	ParentName      string
	ParentPhone     string
	ParentEmail     string
	Source          EnquirySource
	ReferralDetails string
	Remarks         string
	FollowUpDate    *time.Time
	AssignedTo      *uuid.UUID
	CreatedBy       *uuid.UUID
}

// UpdateEnquiryRequest represents a request to update an enquiry.
type UpdateEnquiryRequest struct {
	BranchID        *uuid.UUID
	SessionID       *uuid.UUID
	StudentName     *string
	DateOfBirth     *time.Time
	Gender          *string
	ClassApplying   *string
	ParentName      *string
	ParentPhone     *string
	ParentEmail     *string
	Source          *EnquirySource
	ReferralDetails *string
	Remarks         *string
	Status          *EnquiryStatus
	FollowUpDate    *time.Time
	AssignedTo      *uuid.UUID
	UpdatedBy       *uuid.UUID
}

// ListEnquiryFilter contains filters for listing enquiries.
type ListEnquiryFilter struct {
	TenantID      uuid.UUID
	BranchID      *uuid.UUID
	SessionID     *uuid.UUID
	Status        *EnquiryStatus
	Source        *EnquirySource
	ClassApplying string
	Search        string
	StartDate     *time.Time
	EndDate       *time.Time
	AssignedTo    *uuid.UUID
	Page          int
	PageSize      int
}

// CreateFollowUpRequest represents a request to create a follow-up.
type CreateFollowUpRequest struct {
	TenantID     uuid.UUID
	EnquiryID    uuid.UUID
	FollowUpDate time.Time
	ContactMode  ContactMode
	Notes        string
	Outcome      *FollowUpOutcome
	NextFollowUp *time.Time
	CreatedBy    *uuid.UUID
}

// ConvertEnquiryRequest represents a request to convert an enquiry.
type ConvertEnquiryRequest struct {
	TenantID    uuid.UUID
	EnquiryID   uuid.UUID
	SessionID   *uuid.UUID // Optional, uses enquiry's session if not provided
	BranchID    *uuid.UUID // Optional, uses enquiry's branch if not provided
	ConvertedBy *uuid.UUID
}

// Create creates a new admission enquiry.
func (s *EnquiryService) Create(ctx context.Context, req CreateEnquiryRequest) (*AdmissionEnquiry, error) {
	if req.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if req.StudentName == "" {
		return nil, ErrStudentNameRequired
	}
	if req.ClassApplying == "" {
		return nil, ErrClassApplyingRequired
	}
	if req.ParentName == "" {
		return nil, ErrParentNameRequired
	}
	if req.ParentPhone == "" {
		return nil, ErrParentPhoneRequired
	}

	// Validate source if provided
	if req.Source != "" && !req.Source.IsValid() {
		return nil, ErrInvalidEnquirySource
	}

	// Validate gender if provided
	if req.Gender != "" && !Gender(req.Gender).IsValid() {
		return nil, ErrInvalidGender
	}

	// Set default source if not provided
	source := req.Source
	if source == "" {
		source = EnquirySourceWalkIn
	}

	// Generate enquiry number
	var enquiryNumber string
	err := s.db.WithContext(ctx).Raw(
		"SELECT get_next_enquiry_number(?)",
		req.TenantID,
	).Scan(&enquiryNumber).Error
	if err != nil {
		return nil, fmt.Errorf("failed to generate enquiry number: %w", err)
	}

	enquiry := &AdmissionEnquiry{
		TenantID:        req.TenantID,
		BranchID:        req.BranchID,
		SessionID:       req.SessionID,
		EnquiryNumber:   enquiryNumber,
		StudentName:     req.StudentName,
		DateOfBirth:     req.DateOfBirth,
		Gender:          req.Gender,
		ClassApplying:   req.ClassApplying,
		ParentName:      req.ParentName,
		ParentPhone:     req.ParentPhone,
		ParentEmail:     req.ParentEmail,
		Source:          source,
		ReferralDetails: req.ReferralDetails,
		Remarks:         req.Remarks,
		Status:          EnquiryStatusNew,
		FollowUpDate:    req.FollowUpDate,
		AssignedTo:      req.AssignedTo,
		CreatedBy:       req.CreatedBy,
		UpdatedBy:       req.CreatedBy,
	}

	if err := s.db.WithContext(ctx).Create(enquiry).Error; err != nil {
		return nil, fmt.Errorf("failed to create enquiry: %w", err)
	}

	return s.GetByID(ctx, req.TenantID, enquiry.ID)
}

// GetByID retrieves an enquiry by ID.
func (s *EnquiryService) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*AdmissionEnquiry, error) {
	var enquiry AdmissionEnquiry
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&enquiry).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEnquiryNotFound
		}
		return nil, fmt.Errorf("failed to get enquiry: %w", err)
	}
	return &enquiry, nil
}

// GetByEnquiryNumber retrieves an enquiry by its enquiry number.
func (s *EnquiryService) GetByEnquiryNumber(ctx context.Context, tenantID uuid.UUID, enquiryNumber string) (*AdmissionEnquiry, error) {
	var enquiry AdmissionEnquiry
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND enquiry_number = ?", tenantID, enquiryNumber).
		First(&enquiry).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEnquiryNotFound
		}
		return nil, fmt.Errorf("failed to get enquiry: %w", err)
	}
	return &enquiry, nil
}

// List retrieves enquiries with optional filtering.
func (s *EnquiryService) List(ctx context.Context, filter ListEnquiryFilter) ([]AdmissionEnquiry, int64, error) {
	query := s.db.WithContext(ctx).
		Model(&AdmissionEnquiry{}).
		Where("tenant_id = ?", filter.TenantID)

	if filter.BranchID != nil {
		query = query.Where("branch_id = ?", *filter.BranchID)
	}

	if filter.SessionID != nil {
		query = query.Where("session_id = ?", *filter.SessionID)
	}

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	if filter.Source != nil {
		query = query.Where("source = ?", *filter.Source)
	}

	if filter.ClassApplying != "" {
		query = query.Where("class_applying = ?", filter.ClassApplying)
	}

	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where(
			"student_name ILIKE ? OR parent_name ILIKE ? OR parent_phone ILIKE ? OR enquiry_number ILIKE ?",
			search, search, search, search,
		)
	}

	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}

	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}

	if filter.AssignedTo != nil {
		query = query.Where("assigned_to = ?", *filter.AssignedTo)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count enquiries: %w", err)
	}

	// Apply pagination
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.PageSize

	var enquiries []AdmissionEnquiry
	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&enquiries).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list enquiries: %w", err)
	}

	return enquiries, total, nil
}

// Update updates an enquiry.
func (s *EnquiryService) Update(ctx context.Context, tenantID, id uuid.UUID, req UpdateEnquiryRequest) (*AdmissionEnquiry, error) {
	enquiry, err := s.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	// Cannot update closed or converted enquiries
	if enquiry.Status == EnquiryStatusClosed || enquiry.Status == EnquiryStatusConverted {
		return nil, ErrEnquiryClosed
	}

	updates := make(map[string]interface{})

	if req.BranchID != nil {
		updates["branch_id"] = req.BranchID
	}
	if req.SessionID != nil {
		updates["session_id"] = req.SessionID
	}
	if req.StudentName != nil {
		updates["student_name"] = *req.StudentName
	}
	if req.DateOfBirth != nil {
		updates["date_of_birth"] = req.DateOfBirth
	}
	if req.Gender != nil {
		if *req.Gender != "" && !Gender(*req.Gender).IsValid() {
			return nil, ErrInvalidGender
		}
		updates["gender"] = *req.Gender
	}
	if req.ClassApplying != nil {
		updates["class_applying"] = *req.ClassApplying
	}
	if req.ParentName != nil {
		updates["parent_name"] = *req.ParentName
	}
	if req.ParentPhone != nil {
		updates["parent_phone"] = *req.ParentPhone
	}
	if req.ParentEmail != nil {
		updates["parent_email"] = *req.ParentEmail
	}
	if req.Source != nil {
		if !req.Source.IsValid() {
			return nil, ErrInvalidEnquirySource
		}
		updates["source"] = *req.Source
	}
	if req.ReferralDetails != nil {
		updates["referral_details"] = *req.ReferralDetails
	}
	if req.Remarks != nil {
		updates["remarks"] = *req.Remarks
	}
	if req.Status != nil {
		if !req.Status.IsValid() {
			return nil, ErrInvalidEnquiryStatus
		}
		updates["status"] = *req.Status
	}
	if req.FollowUpDate != nil {
		updates["follow_up_date"] = req.FollowUpDate
	}
	if req.AssignedTo != nil {
		updates["assigned_to"] = req.AssignedTo
	}
	if req.UpdatedBy != nil {
		updates["updated_by"] = req.UpdatedBy
	}

	if len(updates) > 0 {
		if err := s.db.WithContext(ctx).Model(enquiry).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("failed to update enquiry: %w", err)
		}
	}

	return s.GetByID(ctx, tenantID, id)
}

// Delete soft-deletes an enquiry.
func (s *EnquiryService) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	enquiry, err := s.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}

	// Cannot delete converted enquiries
	if enquiry.Status == EnquiryStatusConverted {
		return ErrEnquiryAlreadyConverted
	}

	if err := s.db.WithContext(ctx).Delete(enquiry).Error; err != nil {
		return fmt.Errorf("failed to delete enquiry: %w", err)
	}

	return nil
}

// AddFollowUp adds a follow-up to an enquiry.
func (s *EnquiryService) AddFollowUp(ctx context.Context, req CreateFollowUpRequest) (*EnquiryFollowUp, error) {
	if req.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if req.EnquiryID == uuid.Nil {
		return nil, ErrEnquiryNotFound
	}
	if req.FollowUpDate.IsZero() {
		return nil, ErrFollowUpDateRequired
	}

	// Validate contact mode
	if req.ContactMode != "" && !req.ContactMode.IsValid() {
		return nil, ErrInvalidContactMode
	}

	// Validate outcome if provided
	if req.Outcome != nil && !req.Outcome.IsValid() {
		return nil, ErrInvalidFollowUpOutcome
	}

	// Verify enquiry exists
	enquiry, err := s.GetByID(ctx, req.TenantID, req.EnquiryID)
	if err != nil {
		return nil, err
	}

	// Cannot add follow-up to closed or converted enquiries
	if enquiry.Status == EnquiryStatusClosed || enquiry.Status == EnquiryStatusConverted {
		return nil, ErrEnquiryClosed
	}

	// Set default contact mode if not provided
	contactMode := req.ContactMode
	if contactMode == "" {
		contactMode = ContactModePhone
	}

	followUp := &EnquiryFollowUp{
		TenantID:     req.TenantID,
		EnquiryID:    req.EnquiryID,
		FollowUpDate: req.FollowUpDate,
		ContactMode:  contactMode,
		Notes:        req.Notes,
		Outcome:      req.Outcome,
		NextFollowUp: req.NextFollowUp,
		CreatedBy:    req.CreatedBy,
	}

	// Use transaction to create follow-up and update enquiry
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(followUp).Error; err != nil {
			return fmt.Errorf("failed to create follow-up: %w", err)
		}

		// Update enquiry status and next follow-up date
		updates := make(map[string]interface{})

		// Update status based on outcome
		if req.Outcome != nil {
			switch *req.Outcome {
			case FollowUpOutcomeInterested:
				updates["status"] = EnquiryStatusIntersted
			case FollowUpOutcomeNotInterested:
				updates["status"] = EnquiryStatusClosed
			case FollowUpOutcomeConverted:
				updates["status"] = EnquiryStatusConverted
			case FollowUpOutcomeFollowUpReq, FollowUpOutcomeNoResponse:
				updates["status"] = EnquiryStatusContacted
			}
		} else if enquiry.Status == EnquiryStatusNew {
			updates["status"] = EnquiryStatusContacted
		}

		// Update next follow-up date
		if req.NextFollowUp != nil {
			updates["follow_up_date"] = req.NextFollowUp
		}

		if len(updates) > 0 {
			updates["updated_by"] = req.CreatedBy
			if err := tx.Model(&AdmissionEnquiry{}).
				Where("id = ?", req.EnquiryID).
				Updates(updates).Error; err != nil {
				return fmt.Errorf("failed to update enquiry: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return followUp, nil
}

// ListFollowUps retrieves follow-ups for an enquiry.
func (s *EnquiryService) ListFollowUps(ctx context.Context, tenantID, enquiryID uuid.UUID) ([]EnquiryFollowUp, error) {
	// Verify enquiry exists
	if _, err := s.GetByID(ctx, tenantID, enquiryID); err != nil {
		return nil, err
	}

	var followUps []EnquiryFollowUp
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND enquiry_id = ?", tenantID, enquiryID).
		Order("follow_up_date DESC, created_at DESC").
		Find(&followUps).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list follow-ups: %w", err)
	}

	return followUps, nil
}

// ConvertToApplication creates an application from an enquiry and marks the enquiry as converted.
func (s *EnquiryService) ConvertToApplication(ctx context.Context, req ConvertEnquiryRequest) (*AdmissionEnquiry, error) {
	if req.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}
	if req.EnquiryID == uuid.Nil {
		return nil, ErrEnquiryNotFound
	}

	enquiry, err := s.GetByID(ctx, req.TenantID, req.EnquiryID)
	if err != nil {
		return nil, err
	}

	// Cannot convert already converted enquiries
	if enquiry.Status == EnquiryStatusConverted {
		return nil, ErrEnquiryAlreadyConverted
	}

	// Cannot convert closed enquiries
	if enquiry.Status == EnquiryStatusClosed {
		return nil, ErrEnquiryClosed
	}

	// Determine session ID - use request value or fall back to enquiry's session
	var sessionID uuid.UUID
	if req.SessionID != nil {
		sessionID = *req.SessionID
	} else if enquiry.SessionID != nil {
		sessionID = *enquiry.SessionID
	} else {
		return nil, ErrSessionRequired
	}

	// Determine branch ID - use request value or fall back to enquiry's branch
	var branchID *uuid.UUID
	if req.BranchID != nil {
		branchID = req.BranchID
	} else {
		branchID = enquiry.BranchID
	}

	// Generate application number
	appNumber, err := s.generateApplicationNumber(ctx, req.TenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate application number: %w", err)
	}

	// Create application from enquiry data
	application := &models.AdmissionApplication{
		TenantID:          req.TenantID,
		SessionID:         sessionID,
		BranchID:          branchID,
		EnquiryID:         &enquiry.ID,
		ApplicationNumber: appNumber,
		StudentName:       enquiry.StudentName,
		DateOfBirth:       enquiry.DateOfBirth,
		Gender:            enquiry.Gender,
		ClassApplying:     enquiry.ClassApplying,
		FatherName:        enquiry.ParentName,
		FatherPhone:       enquiry.ParentPhone,
		FatherEmail:       enquiry.ParentEmail,
		Status:            models.ApplicationStatusDraft,
		CreatedBy:         req.ConvertedBy,
		UpdatedBy:         req.ConvertedBy,
	}

	// Use transaction to ensure atomicity
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create the application
		if err := tx.Create(application).Error; err != nil {
			return fmt.Errorf("failed to create application: %w", err)
		}

		// Update enquiry with converted status and application ID
		updates := map[string]interface{}{
			"status":                   EnquiryStatusConverted,
			"converted_application_id": application.ID,
			"updated_by":               req.ConvertedBy,
		}
		if err := tx.Model(enquiry).Updates(updates).Error; err != nil {
			return fmt.Errorf("failed to update enquiry: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.GetByID(ctx, req.TenantID, req.EnquiryID)
}

// generateApplicationNumber generates a unique application number for the tenant.
func (s *EnquiryService) generateApplicationNumber(ctx context.Context, tenantID uuid.UUID) (string, error) {
	var count int64
	if err := s.db.WithContext(ctx).
		Model(&models.AdmissionApplication{}).
		Where("tenant_id = ?", tenantID).
		Count(&count).Error; err != nil {
		return "", err
	}

	// Generate application number in format: APP-YYYYMMDD-XXXXX
	now := time.Now()
	return fmt.Sprintf("APP-%s-%05d", now.Format("20060102"), count+1), nil
}

// Count returns the total number of enquiries for a tenant with optional status filter.
func (s *EnquiryService) Count(ctx context.Context, tenantID uuid.UUID, status *EnquiryStatus) (int64, error) {
	query := s.db.WithContext(ctx).
		Model(&AdmissionEnquiry{}).
		Where("tenant_id = ?", tenantID)

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count enquiries: %w", err)
	}
	return count, nil
}

// GetPendingFollowUps retrieves enquiries with pending follow-ups for a date range.
func (s *EnquiryService) GetPendingFollowUps(ctx context.Context, tenantID uuid.UUID, startDate, endDate time.Time) ([]AdmissionEnquiry, error) {
	var enquiries []AdmissionEnquiry
	err := s.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Where("status NOT IN (?, ?)", EnquiryStatusConverted, EnquiryStatusClosed).
		Where("follow_up_date BETWEEN ? AND ?", startDate, endDate).
		Order("follow_up_date ASC").
		Find(&enquiries).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get pending follow-ups: %w", err)
	}
	return enquiries, nil
}
