package behavioral

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"msls-backend/internal/pkg/database/models"
)

// Service handles business logic for behavioral incidents
type Service struct {
	repo *Repository
}

// NewService creates a new behavioral service
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// =============================================================================
// Incident Operations
// =============================================================================

// CreateIncident creates a new behavioral incident
func (s *Service) CreateIncident(ctx context.Context, tenantID, studentID, reportedBy uuid.UUID, req CreateIncidentRequest) (*IncidentResponse, error) {
	// Parse date
	incidentDate, err := time.Parse("2006-01-02", req.IncidentDate)
	if err != nil {
		return nil, ErrInvalidIncidentDate
	}

	// Set default severity if not provided
	severity := req.Severity
	if severity == "" {
		severity = models.BehavioralSeverityMedium
	}

	// Convert witnesses to JSON
	var witnessesJSON json.RawMessage
	if len(req.Witnesses) > 0 {
		witnessesJSON, _ = json.Marshal(req.Witnesses)
	}

	incident := &models.StudentBehavioralIncident{
		ID:                    uuid.Must(uuid.NewV7()),
		TenantID:              tenantID,
		StudentID:             studentID,
		IncidentType:          req.IncidentType,
		Severity:              severity,
		IncidentDate:          incidentDate,
		IncidentTime:          req.IncidentTime,
		Location:              req.Location,
		Description:           req.Description,
		Witnesses:             witnessesJSON,
		StudentResponse:       req.StudentResponse,
		ActionTaken:           req.ActionTaken,
		ParentMeetingRequired: req.ParentMeetingRequired,
		ReportedBy:            reportedBy,
	}

	if err := s.repo.CreateIncident(ctx, incident); err != nil {
		return nil, err
	}

	// Fetch with relations
	incident, err = s.repo.GetIncident(ctx, tenantID, incident.ID)
	if err != nil {
		return nil, err
	}

	return s.toIncidentResponse(incident), nil
}

// GetIncident retrieves an incident by ID
func (s *Service) GetIncident(ctx context.Context, tenantID, studentID, incidentID uuid.UUID) (*IncidentResponse, error) {
	incident, err := s.repo.GetIncidentByStudent(ctx, tenantID, studentID, incidentID)
	if err != nil {
		return nil, ErrIncidentNotFound
	}
	return s.toIncidentResponse(incident), nil
}

// ListIncidents retrieves incidents for a student
func (s *Service) ListIncidents(ctx context.Context, tenantID, studentID uuid.UUID, filter IncidentFilter) (*IncidentListResponse, error) {
	incidents, total, err := s.repo.ListIncidents(ctx, tenantID, studentID, filter)
	if err != nil {
		return nil, err
	}

	response := &IncidentListResponse{
		Incidents: make([]IncidentResponse, len(incidents)),
		Total:     int(total),
	}

	for i, incident := range incidents {
		response.Incidents[i] = *s.toIncidentResponse(&incident)
	}

	return response, nil
}

// UpdateIncident updates an incident
func (s *Service) UpdateIncident(ctx context.Context, tenantID, studentID, incidentID uuid.UUID, req UpdateIncidentRequest) (*IncidentResponse, error) {
	incident, err := s.repo.GetIncidentByStudent(ctx, tenantID, studentID, incidentID)
	if err != nil {
		return nil, ErrIncidentNotFound
	}

	// Update fields
	if req.IncidentType != nil {
		incident.IncidentType = *req.IncidentType
	}
	if req.Severity != nil {
		incident.Severity = *req.Severity
	}
	if req.IncidentDate != nil {
		date, err := time.Parse("2006-01-02", *req.IncidentDate)
		if err != nil {
			return nil, ErrInvalidIncidentDate
		}
		incident.IncidentDate = date
	}
	if req.IncidentTime != nil {
		incident.IncidentTime = *req.IncidentTime
	}
	if req.Location != nil {
		incident.Location = req.Location
	}
	if req.Description != nil {
		incident.Description = *req.Description
	}
	if req.Witnesses != nil {
		witnessesJSON, _ := json.Marshal(req.Witnesses)
		incident.Witnesses = witnessesJSON
	}
	if req.StudentResponse != nil {
		incident.StudentResponse = req.StudentResponse
	}
	if req.ActionTaken != nil {
		incident.ActionTaken = *req.ActionTaken
	}
	if req.ParentMeetingRequired != nil {
		incident.ParentMeetingRequired = *req.ParentMeetingRequired
	}
	if req.ParentNotified != nil {
		incident.ParentNotified = *req.ParentNotified
		if *req.ParentNotified {
			now := time.Now()
			incident.ParentNotifiedAt = &now
		}
	}

	incident.UpdatedAt = time.Now()

	if err := s.repo.UpdateIncident(ctx, incident); err != nil {
		return nil, err
	}

	return s.toIncidentResponse(incident), nil
}

// DeleteIncident deletes an incident
func (s *Service) DeleteIncident(ctx context.Context, tenantID, studentID, incidentID uuid.UUID) error {
	// Verify incident belongs to student
	_, err := s.repo.GetIncidentByStudent(ctx, tenantID, studentID, incidentID)
	if err != nil {
		return ErrIncidentNotFound
	}

	return s.repo.DeleteIncident(ctx, tenantID, incidentID)
}

// GetBehaviorSummary retrieves behavior statistics for a student
func (s *Service) GetBehaviorSummary(ctx context.Context, tenantID, studentID uuid.UUID) (*BehaviorSummary, error) {
	return s.repo.GetBehaviorSummary(ctx, tenantID, studentID)
}

// =============================================================================
// Follow-Up Operations
// =============================================================================

// CreateFollowUp creates a new follow-up
func (s *Service) CreateFollowUp(ctx context.Context, tenantID, incidentID uuid.UUID, createdBy uuid.UUID, req CreateFollowUpRequest) (*FollowUpResponse, error) {
	// Verify incident exists
	incident, err := s.repo.GetIncident(ctx, tenantID, incidentID)
	if err != nil {
		return nil, ErrIncidentNotFound
	}

	// Parse date
	scheduledDate, err := time.Parse("2006-01-02", req.ScheduledDate)
	if err != nil {
		return nil, ErrInvalidScheduledDate
	}

	// Convert participants to JSON
	var participantsJSON json.RawMessage
	if len(req.Participants) > 0 {
		participantsJSON, _ = json.Marshal(req.Participants)
	}

	followUp := &models.IncidentFollowUp{
		ID:               uuid.Must(uuid.NewV7()),
		TenantID:         tenantID,
		IncidentID:       incident.ID,
		ScheduledDate:    scheduledDate,
		ScheduledTime:    req.ScheduledTime,
		Participants:     participantsJSON,
		ExpectedOutcomes: req.ExpectedOutcomes,
		Status:           models.FollowUpStatusPending,
		CreatedBy:        &createdBy,
	}

	if err := s.repo.CreateFollowUp(ctx, followUp); err != nil {
		return nil, err
	}

	return s.toFollowUpResponse(followUp), nil
}

// UpdateFollowUp updates a follow-up
func (s *Service) UpdateFollowUp(ctx context.Context, tenantID, incidentID, followUpID uuid.UUID, userID uuid.UUID, req UpdateFollowUpRequest) (*FollowUpResponse, error) {
	followUp, err := s.repo.GetFollowUp(ctx, tenantID, followUpID)
	if err != nil {
		return nil, ErrFollowUpNotFound
	}

	// Verify follow-up belongs to incident
	if followUp.IncidentID != incidentID {
		return nil, ErrFollowUpNotFound
	}

	// Update fields
	if req.ScheduledDate != nil {
		date, err := time.Parse("2006-01-02", *req.ScheduledDate)
		if err != nil {
			return nil, ErrInvalidScheduledDate
		}
		followUp.ScheduledDate = date
	}
	if req.ScheduledTime != nil {
		followUp.ScheduledTime = req.ScheduledTime
	}
	if req.Participants != nil {
		participantsJSON, _ := json.Marshal(req.Participants)
		followUp.Participants = participantsJSON
	}
	if req.ExpectedOutcomes != nil {
		followUp.ExpectedOutcomes = req.ExpectedOutcomes
	}
	if req.MeetingNotes != nil {
		followUp.MeetingNotes = req.MeetingNotes
	}
	if req.ActualOutcomes != nil {
		followUp.ActualOutcomes = req.ActualOutcomes
	}
	if req.Status != nil {
		followUp.Status = *req.Status
		if *req.Status == models.FollowUpStatusCompleted {
			now := time.Now()
			followUp.CompletedAt = &now
			followUp.CompletedBy = &userID
		}
	}

	if err := s.repo.UpdateFollowUp(ctx, followUp); err != nil {
		return nil, err
	}

	return s.toFollowUpResponse(followUp), nil
}

// DeleteFollowUp deletes a follow-up
func (s *Service) DeleteFollowUp(ctx context.Context, tenantID, incidentID, followUpID uuid.UUID) error {
	followUp, err := s.repo.GetFollowUp(ctx, tenantID, followUpID)
	if err != nil {
		return ErrFollowUpNotFound
	}

	// Verify follow-up belongs to incident
	if followUp.IncidentID != incidentID {
		return ErrFollowUpNotFound
	}

	return s.repo.DeleteFollowUp(ctx, tenantID, followUpID)
}

// ListPendingFollowUps retrieves all pending follow-ups
func (s *Service) ListPendingFollowUps(ctx context.Context, tenantID uuid.UUID, limit, offset int) (*PendingFollowUpsResponse, error) {
	followUps, total, err := s.repo.ListPendingFollowUps(ctx, tenantID, limit, offset)
	if err != nil {
		return nil, err
	}

	response := &PendingFollowUpsResponse{
		FollowUps: make([]PendingFollowUpItem, len(followUps)),
		Total:     int(total),
	}

	for i, fu := range followUps {
		item := PendingFollowUpItem{
			FollowUpResponse: *s.toFollowUpResponse(&fu),
		}
		if fu.Incident != nil {
			item.StudentID = fu.Incident.StudentID
			item.IncidentType = string(fu.Incident.IncidentType)
			item.IncidentDate = fu.Incident.IncidentDate.Format("2006-01-02")
			if fu.Incident.Student != nil {
				item.StudentName = fu.Incident.Student.FullName()
			}
		}
		response.FollowUps[i] = item
	}

	return response, nil
}

// =============================================================================
// Response Converters
// =============================================================================

func (s *Service) toIncidentResponse(incident *models.StudentBehavioralIncident) *IncidentResponse {
	response := &IncidentResponse{
		ID:                    incident.ID,
		StudentID:             incident.StudentID,
		IncidentType:          incident.IncidentType,
		IncidentTypeLabel:     GetIncidentTypeLabel(incident.IncidentType),
		Severity:              incident.Severity,
		SeverityLabel:         GetSeverityLabel(incident.Severity),
		IncidentDate:          incident.IncidentDate.Format("2006-01-02"),
		IncidentTime:          incident.IncidentTime,
		Location:              incident.Location,
		Description:           incident.Description,
		StudentResponse:       incident.StudentResponse,
		ActionTaken:           incident.ActionTaken,
		ParentMeetingRequired: incident.ParentMeetingRequired,
		ParentNotified:        incident.ParentNotified,
		ParentNotifiedAt:      incident.ParentNotifiedAt,
		ReportedBy:            incident.ReportedBy,
		CreatedAt:             incident.CreatedAt,
		UpdatedAt:             incident.UpdatedAt,
	}

	// Parse witnesses
	if len(incident.Witnesses) > 0 {
		var witnesses []string
		json.Unmarshal(incident.Witnesses, &witnesses)
		response.Witnesses = witnesses
	}

	// Add reporter name
	if incident.Reporter != nil {
		response.ReporterName = incident.Reporter.FullName()
	}

	// Add follow-ups
	if len(incident.FollowUps) > 0 {
		response.FollowUps = make([]FollowUpResponse, len(incident.FollowUps))
		for i, fu := range incident.FollowUps {
			response.FollowUps[i] = *s.toFollowUpResponse(&fu)
		}
	}

	return response
}

func (s *Service) toFollowUpResponse(followUp *models.IncidentFollowUp) *FollowUpResponse {
	response := &FollowUpResponse{
		ID:               followUp.ID,
		IncidentID:       followUp.IncidentID,
		ScheduledDate:    followUp.ScheduledDate.Format("2006-01-02"),
		ScheduledTime:    followUp.ScheduledTime,
		ExpectedOutcomes: followUp.ExpectedOutcomes,
		MeetingNotes:     followUp.MeetingNotes,
		ActualOutcomes:   followUp.ActualOutcomes,
		Status:           followUp.Status,
		StatusLabel:      GetFollowUpStatusLabel(followUp.Status),
		CompletedAt:      followUp.CompletedAt,
		CreatedAt:        followUp.CreatedAt,
	}

	// Parse participants
	if len(followUp.Participants) > 0 {
		var participants []Participant
		json.Unmarshal(followUp.Participants, &participants)
		response.Participants = participants
	}

	return response
}
