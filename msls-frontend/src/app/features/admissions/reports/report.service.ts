/**
 * MSLS Admission Reports Service
 *
 * HTTP service for admission reports and analytics API calls.
 * Provides dashboard stats, funnel data, and class-wise reports.
 */

import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { map, catchError } from 'rxjs/operators';

import { ApiService } from '../../../core/services/api.service';
import {
  DashboardStats,
  ClassWiseReport,
  ClassWiseReportResponse,
  FunnelStage,
  SourceAnalysis,
  DailyTrendPoint,
  ReportFilterParams,
  ExportRequest,
  AdmissionSessionOption,
  DEFAULT_FUNNEL_STAGES,
} from './report.model';

/**
 * API response interfaces matching backend DTOs
 */
interface ApiFunnelStage {
  stage: string;
  count: number;
  percentage: number;
}

interface ApiFunnelResponse {
  stages: ApiFunnelStage[];
}

interface ApiClassWiseItem {
  className: string;
  totalSeats: number;
  applications: number;
  approved: number;
  enrolled: number;
  waitlisted: number;
  vacant: number;
}

interface ApiClassWiseResponse {
  classes: ApiClassWiseItem[];
}

interface ApiSourceItem {
  source: string;
  count: number;
  percentage: number;
  converted: number;
}

interface ApiSourceResponse {
  sources: ApiSourceItem[];
  totalCount: number;
}

interface ApiDailyTrendItem {
  date: string;
  applications: number;
  enquiries: number;
}

interface ApiDailyTrendResponse {
  trends: ApiDailyTrendItem[];
}

interface ApiSessionItem {
  id: string;
  name: string;
  status: string;
}

interface ApiSessionListResponse {
  sessions: ApiSessionItem[];
  total: number;
}

/**
 * AdmissionReportService - Handles all admission reporting API operations.
 */
@Injectable({ providedIn: 'root' })
export class AdmissionReportService {
  private readonly apiService = inject(ApiService);
  private readonly http = inject(HttpClient);
  private readonly basePath = '/admissions';

  /**
   * Get dashboard statistics
   */
  getDashboardStats(filters?: ReportFilterParams): Observable<DashboardStats> {
    const params = this.buildFilterParams(filters);
    return this.apiService.get<DashboardStats>(`${this.basePath}/dashboard`, { params });
  }

  /**
   * Get conversion funnel data
   */
  getFunnelData(filters?: ReportFilterParams): Observable<FunnelStage[]> {
    const params = this.buildFilterParams(filters);
    return this.apiService.get<ApiFunnelResponse>(`${this.basePath}/reports/funnel`, { params }).pipe(
      map(response => this.transformFunnelData(response))
    );
  }

  /**
   * Get class-wise admission report
   */
  getClassWiseReport(filters?: ReportFilterParams): Observable<ClassWiseReportResponse> {
    const params = this.buildFilterParams(filters);
    return this.apiService.get<ApiClassWiseResponse>(`${this.basePath}/reports/class-wise`, { params }).pipe(
      map(response => this.transformClassWiseData(response))
    );
  }

  /**
   * Get source analysis data
   */
  getSourceAnalysis(filters?: ReportFilterParams): Observable<SourceAnalysis[]> {
    const params = this.buildFilterParams(filters);
    return this.apiService.get<ApiSourceResponse>(`${this.basePath}/reports/source-analysis`, { params }).pipe(
      map(response => this.transformSourceAnalysisData(response))
    );
  }

  /**
   * Get daily trend data
   */
  getDailyTrend(filters?: ReportFilterParams, days: number = 30): Observable<DailyTrendPoint[]> {
    const params = {
      ...this.buildFilterParams(filters),
      days: days,
    };
    return this.apiService.get<ApiDailyTrendResponse>(`${this.basePath}/reports/daily-trend`, { params }).pipe(
      map(response => this.transformDailyTrendData(response))
    );
  }

  /**
   * Get available sessions for filter dropdown
   */
  getSessions(): Observable<AdmissionSessionOption[]> {
    return this.apiService.get<ApiSessionListResponse>('/admission-sessions').pipe(
      map(response => (response.sessions || []).map(s => ({
        id: s.id,
        name: s.name,
        status: this.mapSessionStatus(s.status),
      }))),
      catchError(() => {
        // If session endpoint fails, return empty array
        return of([]);
      })
    );
  }

  /**
   * Export report to file
   */
  exportReport(request: ExportRequest): Observable<Blob> {
    const params: Record<string, string> = {
      report_type: request.reportType,
      format: request.format,
    };

    if (request.filters) {
      if (request.filters.sessionId) {
        params['session_id'] = request.filters.sessionId;
      }
      if (request.filters.branchId) {
        params['branch_id'] = request.filters.branchId;
      }
      if (request.filters.fromDate) {
        params['start_date'] = request.filters.fromDate;
      }
      if (request.filters.toDate) {
        params['end_date'] = request.filters.toDate;
      }
    }

    return this.http.get(`/api${this.basePath}/export`, {
      params,
      responseType: 'blob',
    });
  }

  /**
   * Build filter params for API request
   */
  private buildFilterParams(
    filters?: ReportFilterParams
  ): Record<string, string | number | boolean> | undefined {
    if (!filters) return undefined;

    const params: Record<string, string> = {};

    if (filters.sessionId) {
      params['session_id'] = filters.sessionId;
    }
    if (filters.branchId) {
      params['branch_id'] = filters.branchId;
    }
    if (filters.fromDate) {
      params['start_date'] = filters.fromDate;
    }
    if (filters.toDate) {
      params['end_date'] = filters.toDate;
    }
    if (filters.className) {
      params['class_name'] = filters.className;
    }

    return Object.keys(params).length > 0 ? params : undefined;
  }

  /**
   * Transform API funnel response to frontend model
   */
  private transformFunnelData(response: ApiFunnelResponse): FunnelStage[] {
    const stageConfigs: Record<string, { name: string; variant: FunnelStage['variant']; icon: string }> = {
      enquiry: { name: 'Enquiries', variant: 'info', icon: 'fa-solid fa-question-circle' },
      application: { name: 'Applications', variant: 'primary', icon: 'fa-solid fa-file-alt' },
      approved: { name: 'Approved', variant: 'warning', icon: 'fa-solid fa-check-circle' },
      enrolled: { name: 'Enrolled', variant: 'success', icon: 'fa-solid fa-user-graduate' },
    };

    return response.stages.map(stage => {
      const config = stageConfigs[stage.stage] || {
        name: stage.stage,
        variant: 'info' as FunnelStage['variant'],
        icon: 'fa-solid fa-circle',
      };
      return {
        name: config.name,
        count: stage.count,
        percentage: stage.percentage,
        variant: config.variant,
        icon: config.icon,
      };
    });
  }

  /**
   * Transform API class-wise response to frontend model
   */
  private transformClassWiseData(response: ApiClassWiseResponse): ClassWiseReportResponse {
    return {
      classes: response.classes.map(c => ({
        className: c.className,
        totalSeats: c.totalSeats,
        applications: c.applications,
        approved: c.approved,
        enrolled: c.enrolled,
        waitlisted: c.waitlisted,
        vacant: c.vacant,
        fillPercentage: c.totalSeats > 0 ? (c.enrolled / c.totalSeats) * 100 : 0,
      })),
      generatedAt: new Date().toISOString(),
    };
  }

  /**
   * Transform API source analysis response to frontend model
   */
  private transformSourceAnalysisData(response: ApiSourceResponse): SourceAnalysis[] {
    const sourceLabels: Record<string, string> = {
      walk_in: 'Walk-in',
      website: 'Website',
      referral: 'Referral',
      phone: 'Phone Call',
      advertisement: 'Advertisement',
      social_media: 'Social Media',
      newspaper: 'Newspaper',
      other: 'Other',
    };

    return response.sources.map(s => ({
      source: s.source,
      label: sourceLabels[s.source] || this.formatSourceLabel(s.source),
      count: s.count,
      percentage: s.percentage,
    }));
  }

  /**
   * Transform API daily trend response to frontend model
   */
  private transformDailyTrendData(response: ApiDailyTrendResponse): DailyTrendPoint[] {
    return response.trends.map(t => ({
      date: t.date,
      enquiries: t.enquiries,
      applications: t.applications,
      enrollments: 0, // Backend doesn't track daily enrollments yet
    }));
  }

  /**
   * Format source label from snake_case to Title Case
   */
  private formatSourceLabel(source: string): string {
    return source
      .split('_')
      .map(word => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
      .join(' ');
  }

  /**
   * Map backend session status to frontend status type
   */
  private mapSessionStatus(status: string): 'upcoming' | 'open' | 'closed' {
    switch (status.toLowerCase()) {
      case 'upcoming':
      case 'draft':
        return 'upcoming';
      case 'open':
      case 'active':
        return 'open';
      case 'closed':
      case 'completed':
      case 'cancelled':
        return 'closed';
      default:
        return 'closed';
    }
  }
}
