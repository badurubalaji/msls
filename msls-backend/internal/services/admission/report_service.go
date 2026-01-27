// Package admission provides admission management services.
package admission

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msls-backend/internal/pkg/database/models"
)

// ReportService handles admission report and analytics operations.
type ReportService struct {
	db *gorm.DB
}

// NewReportService creates a new ReportService instance.
func NewReportService(db *gorm.DB) *ReportService {
	return &ReportService{db: db}
}

// DashboardFilter contains filters for dashboard statistics.
type DashboardFilter struct {
	TenantID       uuid.UUID
	BranchID       *uuid.UUID
	SessionID      *uuid.UUID
	AcademicYearID *uuid.UUID
	StartDate      *time.Time
	EndDate        *time.Time
}

// ConversionRates represents conversion rates in the admission funnel.
type ConversionRates struct {
	EnquiryToApplication    float64 `json:"enquiryToApplication"`
	ApplicationToApproved   float64 `json:"applicationToApproved"`
	ApprovedToEnrolled      float64 `json:"approvedToEnrolled"`
}

// DashboardStats represents the overall dashboard statistics.
type DashboardStats struct {
	TotalEnquiries    int64           `json:"totalEnquiries"`
	TotalApplications int64           `json:"totalApplications"`
	Approved          int64           `json:"approved"`
	Enrolled          int64           `json:"enrolled"`
	Pending           int64           `json:"pending"`
	Rejected          int64           `json:"rejected"`
	Waitlisted        int64           `json:"waitlisted"`
	ConversionRates   ConversionRates `json:"conversionRates"`
}

// GetDashboardStats retrieves admission dashboard statistics.
func (s *ReportService) GetDashboardStats(ctx context.Context, filter DashboardFilter) (*DashboardStats, error) {
	if filter.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}

	stats := &DashboardStats{}

	// Count enquiries
	enquiryQuery := s.db.WithContext(ctx).
		Model(&models.AdmissionEnquiry{}).
		Where("tenant_id = ?", filter.TenantID)

	if filter.BranchID != nil {
		enquiryQuery = enquiryQuery.Where("branch_id = ?", *filter.BranchID)
	}
	if filter.SessionID != nil {
		enquiryQuery = enquiryQuery.Where("session_id = ?", *filter.SessionID)
	}
	if filter.StartDate != nil {
		enquiryQuery = enquiryQuery.Where("enquiry_date >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		enquiryQuery = enquiryQuery.Where("enquiry_date <= ?", *filter.EndDate)
	}

	if err := enquiryQuery.Count(&stats.TotalEnquiries).Error; err != nil {
		return nil, fmt.Errorf("failed to count enquiries: %w", err)
	}

	// Build base application query
	appQuery := s.db.WithContext(ctx).
		Model(&models.AdmissionApplication{}).
		Where("tenant_id = ?", filter.TenantID)

	if filter.BranchID != nil {
		appQuery = appQuery.Where("branch_id = ?", *filter.BranchID)
	}
	if filter.SessionID != nil {
		appQuery = appQuery.Where("session_id = ?", *filter.SessionID)
	}
	if filter.StartDate != nil {
		appQuery = appQuery.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		appQuery = appQuery.Where("created_at <= ?", *filter.EndDate)
	}

	// Count total applications
	if err := appQuery.Count(&stats.TotalApplications).Error; err != nil {
		return nil, fmt.Errorf("failed to count applications: %w", err)
	}

	// Count by status
	statusCounts := make(map[string]int64)
	var results []struct {
		Status string
		Count  int64
	}

	statusQuery := s.db.WithContext(ctx).
		Model(&models.AdmissionApplication{}).
		Select("status, count(*) as count").
		Where("tenant_id = ?", filter.TenantID)

	if filter.BranchID != nil {
		statusQuery = statusQuery.Where("branch_id = ?", *filter.BranchID)
	}
	if filter.SessionID != nil {
		statusQuery = statusQuery.Where("session_id = ?", *filter.SessionID)
	}
	if filter.StartDate != nil {
		statusQuery = statusQuery.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		statusQuery = statusQuery.Where("created_at <= ?", *filter.EndDate)
	}

	if err := statusQuery.Group("status").Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to count by status: %w", err)
	}

	for _, r := range results {
		statusCounts[r.Status] = r.Count
	}

	stats.Approved = statusCounts[string(models.ApplicationStatusApproved)]
	stats.Enrolled = statusCounts[string(models.ApplicationStatusEnrolled)]
	stats.Pending = statusCounts[string(models.ApplicationStatusSubmitted)] +
		statusCounts[string(models.ApplicationStatusUnderReview)] +
		statusCounts[string(models.ApplicationStatusDraft)]
	stats.Rejected = statusCounts[string(models.ApplicationStatusRejected)]
	stats.Waitlisted = statusCounts[string(models.ApplicationStatusWaitlisted)]

	// Calculate conversion rates
	if stats.TotalEnquiries > 0 {
		stats.ConversionRates.EnquiryToApplication = float64(stats.TotalApplications) / float64(stats.TotalEnquiries) * 100
	}
	if stats.TotalApplications > 0 {
		totalApprovedOrEnrolled := stats.Approved + stats.Enrolled
		stats.ConversionRates.ApplicationToApproved = float64(totalApprovedOrEnrolled) / float64(stats.TotalApplications) * 100
	}
	if stats.Approved+stats.Enrolled > 0 {
		stats.ConversionRates.ApprovedToEnrolled = float64(stats.Enrolled) / float64(stats.Approved+stats.Enrolled) * 100
	}

	// Round conversion rates to 1 decimal place
	stats.ConversionRates.EnquiryToApplication = roundToDecimal(stats.ConversionRates.EnquiryToApplication, 1)
	stats.ConversionRates.ApplicationToApproved = roundToDecimal(stats.ConversionRates.ApplicationToApproved, 1)
	stats.ConversionRates.ApprovedToEnrolled = roundToDecimal(stats.ConversionRates.ApprovedToEnrolled, 1)

	return stats, nil
}

// FunnelStage represents a stage in the admission funnel.
type FunnelStage struct {
	Stage      string  `json:"stage"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

// FunnelReport represents the complete admission funnel.
type FunnelReport struct {
	Stages []FunnelStage `json:"stages"`
}

// GetFunnelReport retrieves the admission funnel report.
func (s *ReportService) GetFunnelReport(ctx context.Context, filter DashboardFilter) (*FunnelReport, error) {
	if filter.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}

	report := &FunnelReport{
		Stages: make([]FunnelStage, 0, 4),
	}

	// Count enquiries
	var enquiryCount int64
	enquiryQuery := s.db.WithContext(ctx).
		Model(&models.AdmissionEnquiry{}).
		Where("tenant_id = ?", filter.TenantID)

	if filter.BranchID != nil {
		enquiryQuery = enquiryQuery.Where("branch_id = ?", *filter.BranchID)
	}
	if filter.SessionID != nil {
		enquiryQuery = enquiryQuery.Where("session_id = ?", *filter.SessionID)
	}

	if err := enquiryQuery.Count(&enquiryCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count enquiries: %w", err)
	}

	// Count applications, approved, and enrolled
	var appCount, approvedCount, enrolledCount int64

	appQuery := s.db.WithContext(ctx).
		Model(&models.AdmissionApplication{}).
		Where("tenant_id = ?", filter.TenantID)

	if filter.BranchID != nil {
		appQuery = appQuery.Where("branch_id = ?", *filter.BranchID)
	}
	if filter.SessionID != nil {
		appQuery = appQuery.Where("session_id = ?", *filter.SessionID)
	}

	if err := appQuery.Count(&appCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count applications: %w", err)
	}

	approvedQuery := s.db.WithContext(ctx).
		Model(&models.AdmissionApplication{}).
		Where("tenant_id = ?", filter.TenantID).
		Where("status IN ?", []string{
			string(models.ApplicationStatusApproved),
			string(models.ApplicationStatusEnrolled),
		})

	if filter.BranchID != nil {
		approvedQuery = approvedQuery.Where("branch_id = ?", *filter.BranchID)
	}
	if filter.SessionID != nil {
		approvedQuery = approvedQuery.Where("session_id = ?", *filter.SessionID)
	}

	if err := approvedQuery.Count(&approvedCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count approved: %w", err)
	}

	enrolledQuery := s.db.WithContext(ctx).
		Model(&models.AdmissionApplication{}).
		Where("tenant_id = ?", filter.TenantID).
		Where("status = ?", models.ApplicationStatusEnrolled)

	if filter.BranchID != nil {
		enrolledQuery = enrolledQuery.Where("branch_id = ?", *filter.BranchID)
	}
	if filter.SessionID != nil {
		enrolledQuery = enrolledQuery.Where("session_id = ?", *filter.SessionID)
	}

	if err := enrolledQuery.Count(&enrolledCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count enrolled: %w", err)
	}

	// Calculate percentages (relative to enquiries as base)
	baseCount := enquiryCount
	if baseCount == 0 {
		baseCount = 1 // Avoid division by zero
	}

	report.Stages = []FunnelStage{
		{
			Stage:      "enquiry",
			Count:      enquiryCount,
			Percentage: 100.0,
		},
		{
			Stage:      "application",
			Count:      appCount,
			Percentage: roundToDecimal(float64(appCount)/float64(baseCount)*100, 1),
		},
		{
			Stage:      "approved",
			Count:      approvedCount,
			Percentage: roundToDecimal(float64(approvedCount)/float64(baseCount)*100, 1),
		},
		{
			Stage:      "enrolled",
			Count:      enrolledCount,
			Percentage: roundToDecimal(float64(enrolledCount)/float64(baseCount)*100, 1),
		},
	}

	return report, nil
}

// ClassWiseReport represents the class-wise admission report.
type ClassWiseReport struct {
	ClassName    string `json:"className"`
	TotalSeats   int    `json:"totalSeats"`
	Applications int64  `json:"applications"`
	Approved     int64  `json:"approved"`
	Enrolled     int64  `json:"enrolled"`
	Waitlisted   int64  `json:"waitlisted"`
	Vacant       int    `json:"vacant"`
}

// ClassWiseReportResponse wraps the class-wise report.
type ClassWiseReportResponse struct {
	Classes []ClassWiseReport `json:"classes"`
}

// GetClassWiseReport retrieves the class-wise admission report.
func (s *ReportService) GetClassWiseReport(ctx context.Context, filter DashboardFilter) (*ClassWiseReportResponse, error) {
	if filter.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}

	response := &ClassWiseReportResponse{
		Classes: make([]ClassWiseReport, 0),
	}

	// Get seat configurations
	seatQuery := s.db.WithContext(ctx).
		Model(&models.AdmissionSeat{}).
		Where("tenant_id = ?", filter.TenantID)

	if filter.SessionID != nil {
		seatQuery = seatQuery.Where("session_id = ?", *filter.SessionID)
	}

	var seats []models.AdmissionSeat
	if err := seatQuery.Find(&seats).Error; err != nil {
		return nil, fmt.Errorf("failed to get seats: %w", err)
	}

	// Get application counts by class
	type classStats struct {
		ClassName  string
		Status     string
		Count      int64
	}

	var appStats []classStats
	appQuery := s.db.WithContext(ctx).
		Model(&models.AdmissionApplication{}).
		Select("class_applying as class_name, status, count(*) as count").
		Where("tenant_id = ?", filter.TenantID)

	if filter.BranchID != nil {
		appQuery = appQuery.Where("branch_id = ?", *filter.BranchID)
	}
	if filter.SessionID != nil {
		appQuery = appQuery.Where("session_id = ?", *filter.SessionID)
	}

	if err := appQuery.Group("class_applying, status").Scan(&appStats).Error; err != nil {
		return nil, fmt.Errorf("failed to get application stats: %w", err)
	}

	// Build a map of class name to stats
	classMap := make(map[string]*ClassWiseReport)

	// Initialize from seats
	for _, seat := range seats {
		classMap[seat.ClassName] = &ClassWiseReport{
			ClassName:  seat.ClassName,
			TotalSeats: seat.TotalSeats,
			Vacant:     seat.TotalSeats - seat.FilledSeats,
		}
	}

	// Add application stats
	for _, stat := range appStats {
		report, exists := classMap[stat.ClassName]
		if !exists {
			report = &ClassWiseReport{
				ClassName: stat.ClassName,
			}
			classMap[stat.ClassName] = report
		}

		report.Applications += stat.Count

		switch models.ApplicationStatus(stat.Status) {
		case models.ApplicationStatusApproved:
			report.Approved += stat.Count
		case models.ApplicationStatusEnrolled:
			report.Enrolled += stat.Count
		case models.ApplicationStatusWaitlisted:
			report.Waitlisted += stat.Count
		}
	}

	// Convert map to slice
	for _, report := range classMap {
		response.Classes = append(response.Classes, *report)
	}

	return response, nil
}

// SourceAnalysis represents enquiry source analysis.
type SourceAnalysis struct {
	Source     string  `json:"source"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
	Converted  int64   `json:"converted"`
}

// SourceAnalysisResponse wraps the source analysis report.
type SourceAnalysisResponse struct {
	Sources    []SourceAnalysis `json:"sources"`
	TotalCount int64            `json:"totalCount"`
}

// GetSourceAnalysis retrieves the enquiry source analysis.
func (s *ReportService) GetSourceAnalysis(ctx context.Context, filter DashboardFilter) (*SourceAnalysisResponse, error) {
	if filter.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}

	response := &SourceAnalysisResponse{
		Sources: make([]SourceAnalysis, 0),
	}

	// Count enquiries by source
	type sourceCount struct {
		Source string
		Count  int64
	}

	var enquiryCounts []sourceCount
	enquiryQuery := s.db.WithContext(ctx).
		Model(&models.AdmissionEnquiry{}).
		Select("source, count(*) as count").
		Where("tenant_id = ?", filter.TenantID)

	if filter.BranchID != nil {
		enquiryQuery = enquiryQuery.Where("branch_id = ?", *filter.BranchID)
	}
	if filter.SessionID != nil {
		enquiryQuery = enquiryQuery.Where("session_id = ?", *filter.SessionID)
	}
	if filter.StartDate != nil {
		enquiryQuery = enquiryQuery.Where("enquiry_date >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		enquiryQuery = enquiryQuery.Where("enquiry_date <= ?", *filter.EndDate)
	}

	if err := enquiryQuery.Group("source").Scan(&enquiryCounts).Error; err != nil {
		return nil, fmt.Errorf("failed to count enquiries by source: %w", err)
	}

	// Count converted enquiries by source
	var convertedCounts []sourceCount
	convertedQuery := s.db.WithContext(ctx).
		Model(&models.AdmissionEnquiry{}).
		Select("source, count(*) as count").
		Where("tenant_id = ?", filter.TenantID).
		Where("status = ?", models.EnquiryStatusConverted)

	if filter.BranchID != nil {
		convertedQuery = convertedQuery.Where("branch_id = ?", *filter.BranchID)
	}
	if filter.SessionID != nil {
		convertedQuery = convertedQuery.Where("session_id = ?", *filter.SessionID)
	}

	if err := convertedQuery.Group("source").Scan(&convertedCounts).Error; err != nil {
		return nil, fmt.Errorf("failed to count converted enquiries: %w", err)
	}

	// Build converted map
	convertedMap := make(map[string]int64)
	for _, c := range convertedCounts {
		convertedMap[c.Source] = c.Count
	}

	// Calculate total
	for _, c := range enquiryCounts {
		response.TotalCount += c.Count
	}

	// Build response
	for _, c := range enquiryCounts {
		percentage := float64(0)
		if response.TotalCount > 0 {
			percentage = float64(c.Count) / float64(response.TotalCount) * 100
		}

		response.Sources = append(response.Sources, SourceAnalysis{
			Source:     c.Source,
			Count:      c.Count,
			Percentage: roundToDecimal(percentage, 1),
			Converted:  convertedMap[c.Source],
		})
	}

	return response, nil
}

// DailyTrend represents daily application trend data.
type DailyTrend struct {
	Date         string `json:"date"`
	Applications int64  `json:"applications"`
	Enquiries    int64  `json:"enquiries"`
}

// DailyTrendResponse wraps the daily trend report.
type DailyTrendResponse struct {
	Trends []DailyTrend `json:"trends"`
}

// GetDailyTrend retrieves the daily application trend.
func (s *ReportService) GetDailyTrend(ctx context.Context, filter DashboardFilter, days int) (*DailyTrendResponse, error) {
	if filter.TenantID == uuid.Nil {
		return nil, ErrTenantIDRequired
	}

	if days <= 0 {
		days = 30 // Default to 30 days
	}

	response := &DailyTrendResponse{
		Trends: make([]DailyTrend, 0),
	}

	// Calculate date range
	endDate := time.Now()
	if filter.EndDate != nil {
		endDate = *filter.EndDate
	}
	startDate := endDate.AddDate(0, 0, -days)
	if filter.StartDate != nil {
		startDate = *filter.StartDate
	}

	// Get application counts by date using raw SQL for reliability
	type dateCount struct {
		DateStr string `gorm:"column:date_str"`
		Count   int64  `gorm:"column:count"`
	}

	appMap := make(map[string]int64)
	enquiryMap := make(map[string]int64)

	// Get application counts - use raw query for better compatibility
	var appCounts []dateCount
	appSQL := `SELECT to_char(DATE(created_at), 'YYYY-MM-DD') as date_str, count(*) as count
		FROM admission_applications
		WHERE tenant_id = ? AND created_at >= ? AND created_at <= ? AND deleted_at IS NULL`
	appArgs := []interface{}{filter.TenantID, startDate, endDate}

	if filter.BranchID != nil {
		appSQL += " AND branch_id = ?"
		appArgs = append(appArgs, *filter.BranchID)
	}
	if filter.SessionID != nil {
		appSQL += " AND session_id = ?"
		appArgs = append(appArgs, *filter.SessionID)
	}
	appSQL += " GROUP BY DATE(created_at) ORDER BY date_str"

	if err := s.db.WithContext(ctx).Raw(appSQL, appArgs...).Scan(&appCounts).Error; err != nil {
		// Log but don't fail - return empty data
		appCounts = []dateCount{}
	}

	for _, c := range appCounts {
		if c.DateStr != "" {
			appMap[c.DateStr] = c.Count
		}
	}

	// Get enquiry counts
	var enquiryCounts []dateCount
	enquirySQL := `SELECT to_char(enquiry_date, 'YYYY-MM-DD') as date_str, count(*) as count
		FROM admission_enquiries
		WHERE tenant_id = ? AND enquiry_date >= ? AND enquiry_date <= ?`
	enquiryArgs := []interface{}{filter.TenantID, startDate, endDate}

	if filter.BranchID != nil {
		enquirySQL += " AND branch_id = ?"
		enquiryArgs = append(enquiryArgs, *filter.BranchID)
	}
	if filter.SessionID != nil {
		enquirySQL += " AND session_id = ?"
		enquiryArgs = append(enquiryArgs, *filter.SessionID)
	}
	enquirySQL += " GROUP BY enquiry_date ORDER BY date_str"

	if err := s.db.WithContext(ctx).Raw(enquirySQL, enquiryArgs...).Scan(&enquiryCounts).Error; err != nil {
		// Log but don't fail - return empty data
		enquiryCounts = []dateCount{}
	}

	for _, c := range enquiryCounts {
		if c.DateStr != "" {
			enquiryMap[c.DateStr] = c.Count
		}
	}

	// Generate daily trend for all days in range
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		response.Trends = append(response.Trends, DailyTrend{
			Date:         dateStr,
			Applications: appMap[dateStr],
			Enquiries:    enquiryMap[dateStr],
		})
	}

	return response, nil
}

// roundToDecimal rounds a float to the specified number of decimal places.
func roundToDecimal(val float64, places int) float64 {
	pow := 1.0
	for i := 0; i < places; i++ {
		pow *= 10
	}
	return float64(int64(val*pow+0.5)) / pow
}
