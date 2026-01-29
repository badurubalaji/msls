/**
 * Compliance Report Component
 * Story 5.8: Staff Document Management
 */
import { Component, OnInit, inject, signal, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink } from '@angular/router';
import { StaffDocumentService } from '../../staff-document.service';
import {
  ComplianceReportResponse,
  ComplianceStats,
  DocumentTypeCompliance,
  StaffComplianceDetail,
  getCategoryLabel,
} from '../../staff-document.model';
import { MslsIconComponent } from '../../../../shared/components/icon/icon.component';

@Component({
  selector: 'app-compliance-report',
  standalone: true,
  imports: [CommonModule, RouterLink, MslsIconComponent],
  template: `
    <div class="container mx-auto p-6">
      <!-- Header -->
      <div class="flex items-center justify-between mb-6">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">Document Compliance Report</h1>
          <p class="text-gray-600">Overview of staff document compliance status</p>
        </div>
        <div class="flex items-center gap-2">
          <button
            (click)="toggleStaffDetails()"
            class="px-4 py-2 text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50"
          >
            {{ showStaffDetails() ? 'Hide' : 'Show' }} Staff Details
          </button>
          <button
            (click)="loadReport()"
            [disabled]="loading()"
            class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 disabled:opacity-50"
          >
            <msls-icon name="arrow-path" class="w-5 h-5 inline-block mr-1" />
            Refresh
          </button>
        </div>
      </div>

      <!-- Loading State -->
      @if (loading()) {
        <div class="flex items-center justify-center h-64">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
        </div>
      }

      @if (!loading() && stats()) {
        <!-- Summary Stats -->
        <div class="grid grid-cols-2 md:grid-cols-5 gap-4 mb-6">
          <div class="bg-white rounded-lg shadow p-4">
            <div class="text-2xl font-bold text-gray-900">{{ stats()?.total_staff }}</div>
            <div class="text-sm text-gray-500">Total Staff</div>
          </div>
          <div class="bg-white rounded-lg shadow p-4">
            <div class="text-2xl font-bold text-blue-600">{{ stats()?.documents_submitted }}</div>
            <div class="text-sm text-gray-500">Documents Submitted</div>
          </div>
          <div class="bg-white rounded-lg shadow p-4">
            <div class="text-2xl font-bold text-green-600">{{ stats()?.verified }}</div>
            <div class="text-sm text-gray-500">Verified</div>
          </div>
          <div class="bg-white rounded-lg shadow p-4">
            <div class="text-2xl font-bold text-amber-600">{{ stats()?.pending_verification }}</div>
            <div class="text-sm text-gray-500">Pending</div>
          </div>
          <div class="bg-white rounded-lg shadow p-4">
            <div class="text-2xl font-bold text-red-600">{{ stats()?.expired }}</div>
            <div class="text-sm text-gray-500">Expired</div>
          </div>
        </div>

        <!-- Compliance Percentage -->
        <div class="bg-white rounded-lg shadow p-6 mb-6">
          <div class="flex items-center justify-between mb-4">
            <h2 class="text-lg font-semibold text-gray-900">Overall Compliance</h2>
            <span class="text-2xl font-bold" [class]="getComplianceClass(stats()?.compliance_percentage || 0)">
              {{ (stats()?.compliance_percentage || 0) | number:'1.1-1' }}%
            </span>
          </div>
          <div class="w-full bg-gray-200 rounded-full h-4">
            <div
              class="h-4 rounded-full transition-all duration-500"
              [class]="getComplianceBarClass(stats()?.compliance_percentage || 0)"
              [style.width.%]="stats()?.compliance_percentage || 0"
            ></div>
          </div>
          <div class="flex justify-between text-sm text-gray-500 mt-2">
            <span>Expiring in 30 days: {{ stats()?.expiring_in_30_days }}</span>
            <span>Expiring in 60 days: {{ stats()?.expiring_in_60_days }}</span>
            <span>Expiring in 90 days: {{ stats()?.expiring_in_90_days }}</span>
          </div>
        </div>

        <!-- Compliance by Document Type -->
        <div class="bg-white rounded-lg shadow mb-6">
          <div class="px-6 py-4 border-b">
            <h2 class="text-lg font-semibold text-gray-900">Compliance by Document Type</h2>
          </div>
          <div class="overflow-x-auto">
            <table class="min-w-full divide-y divide-gray-200">
              <thead class="bg-gray-50">
                <tr>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Document Type</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Category</th>
                  <th class="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">Required</th>
                  <th class="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">Submitted</th>
                  <th class="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">Verified</th>
                  <th class="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">Pending</th>
                  <th class="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">Compliance</th>
                </tr>
              </thead>
              <tbody class="bg-white divide-y divide-gray-200">
                @for (item of byDocumentType(); track item.document_type.id) {
                  <tr class="hover:bg-gray-50">
                    <td class="px-6 py-4 whitespace-nowrap">
                      <div class="font-medium text-gray-900">{{ item.document_type.name }}</div>
                      @if (item.document_type.is_mandatory) {
                        <span class="text-xs text-red-600">Mandatory</span>
                      }
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {{ getCategoryLabel(item.document_type.category) }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-center text-sm text-gray-900">{{ item.required }}</td>
                    <td class="px-6 py-4 whitespace-nowrap text-center text-sm text-blue-600">{{ item.submitted }}</td>
                    <td class="px-6 py-4 whitespace-nowrap text-center text-sm text-green-600">{{ item.verified }}</td>
                    <td class="px-6 py-4 whitespace-nowrap text-center text-sm text-amber-600">{{ item.pending }}</td>
                    <td class="px-6 py-4 whitespace-nowrap text-center">
                      <div class="flex items-center justify-center gap-2">
                        <div class="w-24 bg-gray-200 rounded-full h-2">
                          <div
                            class="h-2 rounded-full"
                            [class]="getComplianceBarClass(item.compliance_percent)"
                            [style.width.%]="item.compliance_percent"
                          ></div>
                        </div>
                        <span class="text-sm font-medium" [class]="getComplianceClass(item.compliance_percent)">
                          {{ item.compliance_percent | number:'1.0-0' }}%
                        </span>
                      </div>
                    </td>
                  </tr>
                }
              </tbody>
            </table>
          </div>
        </div>

        <!-- Staff Details -->
        @if (showStaffDetails() && staffDetails().length > 0) {
          <div class="bg-white rounded-lg shadow">
            <div class="px-6 py-4 border-b">
              <h2 class="text-lg font-semibold text-gray-900">Staff Compliance Details</h2>
            </div>
            <div class="overflow-x-auto">
              <table class="min-w-full divide-y divide-gray-200">
                <thead class="bg-gray-50">
                  <tr>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Staff</th>
                    <th class="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">Required</th>
                    <th class="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">Submitted</th>
                    <th class="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">Verified</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Missing</th>
                    <th class="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">Compliance</th>
                    <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
                  </tr>
                </thead>
                <tbody class="bg-white divide-y divide-gray-200">
                  @for (staff of staffDetails(); track staff.staff_id) {
                    <tr class="hover:bg-gray-50" [class.bg-red-50]="staff.compliance_percent < 50">
                      <td class="px-6 py-4 whitespace-nowrap">
                        <div class="font-medium text-gray-900">{{ staff.staff_name }}</div>
                        <div class="text-sm text-gray-500">{{ staff.employee_id }}</div>
                      </td>
                      <td class="px-6 py-4 whitespace-nowrap text-center text-sm text-gray-900">{{ staff.total_required }}</td>
                      <td class="px-6 py-4 whitespace-nowrap text-center text-sm text-blue-600">{{ staff.submitted }}</td>
                      <td class="px-6 py-4 whitespace-nowrap text-center text-sm text-green-600">{{ staff.verified }}</td>
                      <td class="px-6 py-4">
                        @if (staff.missing_documents && staff.missing_documents.length > 0) {
                          <div class="flex flex-wrap gap-1">
                            @for (doc of staff.missing_documents.slice(0, 3); track doc) {
                              <span class="px-2 py-0.5 text-xs bg-red-100 text-red-800 rounded">{{ doc }}</span>
                            }
                            @if (staff.missing_documents.length > 3) {
                              <span class="px-2 py-0.5 text-xs bg-gray-100 text-gray-600 rounded">
                                +{{ staff.missing_documents.length - 3 }} more
                              </span>
                            }
                          </div>
                        } @else {
                          <span class="text-green-600 text-sm">All submitted</span>
                        }
                      </td>
                      <td class="px-6 py-4 whitespace-nowrap text-center">
                        <span class="text-sm font-medium" [class]="getComplianceClass(staff.compliance_percent)">
                          {{ staff.compliance_percent | number:'1.0-0' }}%
                        </span>
                      </td>
                      <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                        <a [routerLink]="['/staff', staff.staff_id]" class="text-indigo-600 hover:text-indigo-900">
                          View
                        </a>
                      </td>
                    </tr>
                  }
                </tbody>
              </table>
            </div>
          </div>
        }
      }
    </div>
  `,
})
export class ComplianceReportComponent implements OnInit {
  private readonly service = inject(StaffDocumentService);

  // State
  loading = signal(false);
  showStaffDetails = signal(false);
  stats = signal<ComplianceStats | null>(null);
  byDocumentType = signal<DocumentTypeCompliance[]>([]);
  staffDetails = signal<StaffComplianceDetail[]>([]);

  ngOnInit(): void {
    this.loadReport();
  }

  loadReport(): void {
    this.loading.set(true);
    this.service.getComplianceReport(this.showStaffDetails()).subscribe({
      next: (response) => {
        this.stats.set(response.stats);
        this.byDocumentType.set(response.by_document_type || []);
        this.staffDetails.set(response.staff_details || []);
        this.loading.set(false);
      },
      error: () => {
        this.loading.set(false);
      },
    });
  }

  toggleStaffDetails(): void {
    this.showStaffDetails.update(v => !v);
    if (this.showStaffDetails() && this.staffDetails().length === 0) {
      this.loadReport();
    }
  }

  getCategoryLabel = getCategoryLabel;

  getComplianceClass(percent: number): string {
    if (percent >= 80) return 'text-green-600';
    if (percent >= 60) return 'text-amber-600';
    return 'text-red-600';
  }

  getComplianceBarClass(percent: number): string {
    if (percent >= 80) return 'bg-green-500';
    if (percent >= 60) return 'bg-amber-500';
    return 'bg-red-500';
  }
}
