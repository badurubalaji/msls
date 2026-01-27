// Package admission provides HTTP handlers for admission management endpoints.
package admission

import (
	"time"

	"msls-backend/internal/services/admission"
)

// ============================================================================
// Response DTOs
// ============================================================================

// DashboardStatsResponse represents the dashboard statistics response.
type DashboardStatsResponse struct {
	TotalEnquiries    int64                   `json:"totalEnquiries"`
	TotalApplications int64                   `json:"totalApplications"`
	Approved          int64                   `json:"approved"`
	Enrolled          int64                   `json:"enrolled"`
	Pending           int64                   `json:"pending"`
	Rejected          int64                   `json:"rejected"`
	Waitlisted        int64                   `json:"waitlisted"`
	ConversionRates   ConversionRatesResponse `json:"conversionRates"`
}

// ConversionRatesResponse represents conversion rates in the response.
type ConversionRatesResponse struct {
	EnquiryToApplication  float64 `json:"enquiryToApplication"`
	ApplicationToApproved float64 `json:"applicationToApproved"`
	ApprovedToEnrolled    float64 `json:"approvedToEnrolled"`
}

// FunnelStageResponse represents a funnel stage in the response.
type FunnelStageResponse struct {
	Stage      string  `json:"stage"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

// FunnelReportResponse represents the funnel report response.
type FunnelReportResponse struct {
	Stages []FunnelStageResponse `json:"stages"`
}

// ClassWiseReportItem represents a class-wise report item.
type ClassWiseReportItem struct {
	ClassName    string `json:"className"`
	TotalSeats   int    `json:"totalSeats"`
	Applications int64  `json:"applications"`
	Approved     int64  `json:"approved"`
	Enrolled     int64  `json:"enrolled"`
	Waitlisted   int64  `json:"waitlisted"`
	Vacant       int    `json:"vacant"`
}

// ClassWiseReportResponse represents the class-wise report response.
type ClassWiseReportResponse struct {
	Classes []ClassWiseReportItem `json:"classes"`
}

// SourceAnalysisItem represents a source analysis item.
type SourceAnalysisItem struct {
	Source     string  `json:"source"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
	Converted  int64   `json:"converted"`
}

// SourceAnalysisResponse represents the source analysis response.
type SourceAnalysisResponse struct {
	Sources    []SourceAnalysisItem `json:"sources"`
	TotalCount int64                `json:"totalCount"`
}

// DailyTrendItem represents a daily trend data point.
type DailyTrendItem struct {
	Date         string `json:"date"`
	Applications int64  `json:"applications"`
	Enquiries    int64  `json:"enquiries"`
}

// DailyTrendResponse represents the daily trend response.
type DailyTrendResponse struct {
	Trends []DailyTrendItem `json:"trends"`
}

// ============================================================================
// Request Query Parameters
// ============================================================================

// ReportQueryParams represents common query parameters for report endpoints.
type ReportQueryParams struct {
	SessionID      string `form:"session_id"`
	BranchID       string `form:"branch_id"`
	AcademicYearID string `form:"academic_year_id"`
	StartDate      string `form:"start_date"`
	EndDate        string `form:"end_date"`
	Days           int    `form:"days"`
}

// ============================================================================
// Conversion Functions
// ============================================================================

// dashboardStatsToResponse converts service stats to response DTO.
func dashboardStatsToResponse(stats *admission.DashboardStats) DashboardStatsResponse {
	return DashboardStatsResponse{
		TotalEnquiries:    stats.TotalEnquiries,
		TotalApplications: stats.TotalApplications,
		Approved:          stats.Approved,
		Enrolled:          stats.Enrolled,
		Pending:           stats.Pending,
		Rejected:          stats.Rejected,
		Waitlisted:        stats.Waitlisted,
		ConversionRates: ConversionRatesResponse{
			EnquiryToApplication:  stats.ConversionRates.EnquiryToApplication,
			ApplicationToApproved: stats.ConversionRates.ApplicationToApproved,
			ApprovedToEnrolled:    stats.ConversionRates.ApprovedToEnrolled,
		},
	}
}

// funnelReportToResponse converts service funnel report to response DTO.
func funnelReportToResponse(report *admission.FunnelReport) FunnelReportResponse {
	stages := make([]FunnelStageResponse, len(report.Stages))
	for i, stage := range report.Stages {
		stages[i] = FunnelStageResponse{
			Stage:      stage.Stage,
			Count:      stage.Count,
			Percentage: stage.Percentage,
		}
	}
	return FunnelReportResponse{Stages: stages}
}

// classWiseReportToResponse converts service class-wise report to response DTO.
func classWiseReportToResponse(report *admission.ClassWiseReportResponse) ClassWiseReportResponse {
	classes := make([]ClassWiseReportItem, len(report.Classes))
	for i, c := range report.Classes {
		classes[i] = ClassWiseReportItem{
			ClassName:    c.ClassName,
			TotalSeats:   c.TotalSeats,
			Applications: c.Applications,
			Approved:     c.Approved,
			Enrolled:     c.Enrolled,
			Waitlisted:   c.Waitlisted,
			Vacant:       c.Vacant,
		}
	}
	return ClassWiseReportResponse{Classes: classes}
}

// sourceAnalysisToResponse converts service source analysis to response DTO.
func sourceAnalysisToResponse(report *admission.SourceAnalysisResponse) SourceAnalysisResponse {
	sources := make([]SourceAnalysisItem, len(report.Sources))
	for i, s := range report.Sources {
		sources[i] = SourceAnalysisItem{
			Source:     s.Source,
			Count:      s.Count,
			Percentage: s.Percentage,
			Converted:  s.Converted,
		}
	}
	return SourceAnalysisResponse{
		Sources:    sources,
		TotalCount: report.TotalCount,
	}
}

// dailyTrendToResponse converts service daily trend to response DTO.
func dailyTrendToResponse(report *admission.DailyTrendResponse) DailyTrendResponse {
	trends := make([]DailyTrendItem, len(report.Trends))
	for i, t := range report.Trends {
		trends[i] = DailyTrendItem{
			Date:         t.Date,
			Applications: t.Applications,
			Enquiries:    t.Enquiries,
		}
	}
	return DailyTrendResponse{Trends: trends}
}

// parseDateWithError parses a date string in YYYY-MM-DD format.
func parseDateWithError(dateStr string) (*time.Time, error) {
	if dateStr == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
