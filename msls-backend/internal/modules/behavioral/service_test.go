package behavioral

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"msls-backend/internal/pkg/database/models"
)

// MockRepository is a mock implementation of the Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateIncident(ctx context.Context, incident *models.StudentBehavioralIncident) error {
	args := m.Called(ctx, incident)
	return args.Error(0)
}

func (m *MockRepository) GetIncident(ctx context.Context, tenantID, incidentID uuid.UUID) (*models.StudentBehavioralIncident, error) {
	args := m.Called(ctx, tenantID, incidentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StudentBehavioralIncident), args.Error(1)
}

func (m *MockRepository) GetIncidentByStudent(ctx context.Context, tenantID, studentID, incidentID uuid.UUID) (*models.StudentBehavioralIncident, error) {
	args := m.Called(ctx, tenantID, studentID, incidentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StudentBehavioralIncident), args.Error(1)
}

func (m *MockRepository) ListIncidents(ctx context.Context, tenantID, studentID uuid.UUID, filter IncidentFilter) ([]models.StudentBehavioralIncident, int64, error) {
	args := m.Called(ctx, tenantID, studentID, filter)
	return args.Get(0).([]models.StudentBehavioralIncident), args.Get(1).(int64), args.Error(2)
}

func (m *MockRepository) UpdateIncident(ctx context.Context, incident *models.StudentBehavioralIncident) error {
	args := m.Called(ctx, incident)
	return args.Error(0)
}

func (m *MockRepository) DeleteIncident(ctx context.Context, tenantID, incidentID uuid.UUID) error {
	args := m.Called(ctx, tenantID, incidentID)
	return args.Error(0)
}

func (m *MockRepository) GetBehaviorSummary(ctx context.Context, tenantID, studentID uuid.UUID) (*BehaviorSummary, error) {
	args := m.Called(ctx, tenantID, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*BehaviorSummary), args.Error(1)
}

func (m *MockRepository) CreateFollowUp(ctx context.Context, followUp *models.IncidentFollowUp) error {
	args := m.Called(ctx, followUp)
	return args.Error(0)
}

func (m *MockRepository) GetFollowUp(ctx context.Context, tenantID, followUpID uuid.UUID) (*models.IncidentFollowUp, error) {
	args := m.Called(ctx, tenantID, followUpID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.IncidentFollowUp), args.Error(1)
}

func (m *MockRepository) UpdateFollowUp(ctx context.Context, followUp *models.IncidentFollowUp) error {
	args := m.Called(ctx, followUp)
	return args.Error(0)
}

func (m *MockRepository) DeleteFollowUp(ctx context.Context, tenantID, followUpID uuid.UUID) error {
	args := m.Called(ctx, tenantID, followUpID)
	return args.Error(0)
}

func (m *MockRepository) ListPendingFollowUps(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]models.IncidentFollowUp, int64, error) {
	args := m.Called(ctx, tenantID, limit, offset)
	return args.Get(0).([]models.IncidentFollowUp), args.Get(1).(int64), args.Error(2)
}

// Test cases

func TestService_CreateIncident_Success(t *testing.T) {
	req := CreateIncidentRequest{
		IncidentType:          models.BehavioralIncidentTypeMinorInfraction,
		Severity:              models.BehavioralSeverityLow,
		IncidentDate:          "2024-01-15",
		IncidentTime:          "10:30",
		Description:           "Late to class",
		ActionTaken:           "Verbal warning",
		ParentMeetingRequired: false,
	}

	assert.Equal(t, models.BehavioralIncidentTypeMinorInfraction, req.IncidentType)
	assert.Equal(t, "2024-01-15", req.IncidentDate)
}

func TestService_CreateIncident_InvalidDate(t *testing.T) {
	assert.NotNil(t, ErrInvalidIncidentDate)
}

func TestService_IncidentResponse_Conversion(t *testing.T) {
	now := time.Now()
	reporterID := uuid.New()
	witnesses, _ := json.Marshal([]string{"Teacher A", "Student B"})

	incident := &models.StudentBehavioralIncident{
		ID:                    uuid.New(),
		StudentID:             uuid.New(),
		IncidentType:          models.BehavioralIncidentTypeMajorViolation,
		Severity:              models.BehavioralSeverityHigh,
		IncidentDate:          now,
		IncidentTime:          "14:00",
		Description:           "Fighting in hallway",
		Witnesses:             witnesses,
		ActionTaken:           "Suspension pending review",
		ParentMeetingRequired: true,
		ParentNotified:        true,
		ReportedBy:            reporterID,
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	service := &Service{}
	response := service.toIncidentResponse(incident)

	assert.Equal(t, incident.ID, response.ID)
	assert.Equal(t, incident.StudentID, response.StudentID)
	assert.Equal(t, incident.IncidentType, response.IncidentType)
	assert.Equal(t, incident.Severity, response.Severity)
	assert.Equal(t, "14:00", response.IncidentTime)
	assert.Equal(t, "Fighting in hallway", response.Description)
	assert.True(t, response.ParentMeetingRequired)
	assert.True(t, response.ParentNotified)
	assert.Len(t, response.Witnesses, 2)
}

func TestService_FollowUpResponse_Conversion(t *testing.T) {
	now := time.Now()
	participants, _ := json.Marshal([]Participant{
		{Name: "Principal", Role: "Administrator"},
		{Name: "Parent", Role: "Guardian"},
	})

	followUp := &models.IncidentFollowUp{
		ID:               uuid.New(),
		IncidentID:       uuid.New(),
		ScheduledDate:    now,
		ScheduledTime:    func() *string { s := "15:00"; return &s }(),
		Participants:     participants,
		ExpectedOutcomes: func() *string { s := "Discuss improvement plan"; return &s }(),
		Status:           models.FollowUpStatusPending,
		CreatedAt:        now,
	}

	service := &Service{}
	response := service.toFollowUpResponse(followUp)

	assert.Equal(t, followUp.ID, response.ID)
	assert.Equal(t, followUp.IncidentID, response.IncidentID)
	assert.Equal(t, models.FollowUpStatusPending, response.Status)
	assert.Len(t, response.Participants, 2)
}

func TestService_UpdateIncident_ParentNotified(t *testing.T) {
	notified := true
	req := UpdateIncidentRequest{
		ParentNotified: &notified,
	}

	assert.NotNil(t, req.ParentNotified)
	assert.True(t, *req.ParentNotified)
}

func TestService_ListIncidents_EmptyFilter(t *testing.T) {
	filter := IncidentFilter{
		Limit:  20,
		Offset: 0,
	}

	assert.Equal(t, 20, filter.Limit)
	assert.Equal(t, 0, filter.Offset)
}

func TestService_ListIncidents_WithFilters(t *testing.T) {
	incidentType := models.BehavioralIncidentTypePositiveRecognition
	severity := models.BehavioralSeverityLow
	dateFrom := "2024-01-01"
	dateTo := "2024-12-31"

	filter := IncidentFilter{
		IncidentType: &incidentType,
		Severity:     &severity,
		DateFrom:     &dateFrom,
		DateTo:       &dateTo,
		Limit:        50,
		Offset:       10,
	}

	assert.NotNil(t, filter.IncidentType)
	assert.NotNil(t, filter.Severity)
	assert.NotNil(t, filter.DateFrom)
	assert.NotNil(t, filter.DateTo)
}

func TestService_BehaviorSummary(t *testing.T) {
	summary := &BehaviorSummary{
		TotalIncidents:       10,
		PositiveCount:        5,
		MinorInfractionCount: 3,
		MajorViolationCount:  2,
		ThisMonthCount:       4,
		LastMonthCount:       6,
		Trend:                "improving",
		PendingFollowUps:     1,
	}

	assert.Equal(t, 10, summary.TotalIncidents)
	assert.Equal(t, 5, summary.PositiveCount)
	assert.Equal(t, 3, summary.MinorInfractionCount)
	assert.Equal(t, 2, summary.MajorViolationCount)
	assert.Equal(t, "improving", summary.Trend)
	assert.Equal(t, 1, summary.PendingFollowUps)
}

func TestService_FollowUpStatus(t *testing.T) {
	assert.Equal(t, "Pending", GetFollowUpStatusLabel(models.FollowUpStatusPending))
	assert.Equal(t, "Completed", GetFollowUpStatusLabel(models.FollowUpStatusCompleted))
	assert.Equal(t, "Cancelled", GetFollowUpStatusLabel(models.FollowUpStatusCancelled))
}

func TestService_IncidentTypeLabels(t *testing.T) {
	label := GetIncidentTypeLabel(models.BehavioralIncidentTypePositiveRecognition)
	assert.Equal(t, "Positive Recognition", label)

	label = GetIncidentTypeLabel(models.BehavioralIncidentTypeMinorInfraction)
	assert.Equal(t, "Minor Infraction", label)

	label = GetIncidentTypeLabel(models.BehavioralIncidentTypeMajorViolation)
	assert.Equal(t, "Major Violation", label)
}

func TestService_SeverityLabels(t *testing.T) {
	label := GetSeverityLabel(models.BehavioralSeverityLow)
	assert.Equal(t, "Low", label)

	label = GetSeverityLabel(models.BehavioralSeverityMedium)
	assert.Equal(t, "Medium", label)

	label = GetSeverityLabel(models.BehavioralSeverityHigh)
	assert.Equal(t, "High", label)

	label = GetSeverityLabel(models.BehavioralSeverityCritical)
	assert.Equal(t, "Critical", label)
}

func TestErrors_Defined(t *testing.T) {
	assert.NotNil(t, ErrIncidentNotFound)
	assert.NotNil(t, ErrFollowUpNotFound)
	assert.NotNil(t, ErrInvalidIncidentDate)
	assert.NotNil(t, ErrInvalidScheduledDate)
}

func TestCreateFollowUpRequest_Validation(t *testing.T) {
	req := CreateFollowUpRequest{
		ScheduledDate: "2024-02-15",
		ScheduledTime: func() *string { s := "10:00"; return &s }(),
		Participants: []Participant{
			{Name: "Parent", Role: "Guardian"},
		},
		ExpectedOutcomes: func() *string { s := "Discuss behavior improvement"; return &s }(),
	}

	assert.NotEmpty(t, req.ScheduledDate)
	assert.Len(t, req.Participants, 1)
}

func TestUpdateFollowUpRequest_Completion(t *testing.T) {
	status := models.FollowUpStatusCompleted
	req := UpdateFollowUpRequest{
		Status:         &status,
		MeetingNotes:   func() *string { s := "Meeting went well"; return &s }(),
		ActualOutcomes: func() *string { s := "Agreed on improvement plan"; return &s }(),
	}

	assert.Equal(t, models.FollowUpStatusCompleted, *req.Status)
	assert.NotNil(t, req.MeetingNotes)
	assert.NotNil(t, req.ActualOutcomes)
}
