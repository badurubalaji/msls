/**
 * Document List Component - Expiring Documents View
 * Story 5.8: Staff Document Management
 */
import { Component, OnInit, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { StaffDocumentService } from '../../staff-document.service';
import {
  ExpiringDocument,
  getVerificationStatusLabel,
  getVerificationStatusColor,
  formatFileSize,
} from '../../staff-document.model';
import { MslsIconComponent } from '../../../../shared/components/icon/icon.component';

@Component({
  selector: 'app-document-list',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterLink, MslsIconComponent],
  template: `
    <div class="container mx-auto p-6">
      <!-- Header -->
      <div class="flex items-center justify-between mb-6">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">Expiring Documents</h1>
          <p class="text-gray-600">Documents expiring within {{ selectedDays() }} days</p>
        </div>
        <div class="flex items-center gap-4">
          <select
            [(ngModel)]="selectedDays"
            (ngModelChange)="loadExpiringDocuments()"
            class="px-3 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500"
          >
            <option [ngValue]="7">Next 7 days</option>
            <option [ngValue]="30">Next 30 days</option>
            <option [ngValue]="60">Next 60 days</option>
            <option [ngValue]="90">Next 90 days</option>
          </select>
        </div>
      </div>

      <!-- Summary Cards -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
        <div class="bg-white rounded-lg shadow p-4">
          <div class="text-2xl font-bold text-gray-900">{{ documents().length }}</div>
          <div class="text-sm text-gray-500">Total Expiring</div>
        </div>
        <div class="bg-white rounded-lg shadow p-4">
          <div class="text-2xl font-bold text-red-600">{{ expiredCount() }}</div>
          <div class="text-sm text-gray-500">Already Expired</div>
        </div>
        <div class="bg-white rounded-lg shadow p-4">
          <div class="text-2xl font-bold text-amber-600">{{ urgentCount() }}</div>
          <div class="text-sm text-gray-500">Urgent (7 days)</div>
        </div>
        <div class="bg-white rounded-lg shadow p-4">
          <div class="text-2xl font-bold text-blue-600">{{ warningCount() }}</div>
          <div class="text-sm text-gray-500">Warning (30 days)</div>
        </div>
      </div>

      <!-- Loading State -->
      @if (loading()) {
        <div class="flex items-center justify-center h-64">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
        </div>
      }

      <!-- Documents List -->
      @if (!loading() && documents().length > 0) {
        <div class="bg-white shadow rounded-lg overflow-hidden">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Staff</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Document</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Expiry Date</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Days Left</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              @for (item of documents(); track item.document.id) {
                <tr class="hover:bg-gray-50" [class.bg-red-50]="item.days_to_expiry < 0">
                  <td class="px-6 py-4 whitespace-nowrap">
                    <div class="font-medium text-gray-900">{{ item.staff_name }}</div>
                    <div class="text-sm text-gray-500">{{ item.employee_id }}</div>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <div class="font-medium text-gray-900">{{ item.document.document_type?.name || 'Document' }}</div>
                    <div class="text-sm text-gray-500">{{ item.document.file_name }}</div>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {{ item.document.expiry_date | date:'mediumDate' }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <span class="px-2 py-1 text-xs font-medium rounded-full"
                          [class]="getDaysClass(item.days_to_expiry)">
                      {{ formatDaysLeft(item.days_to_expiry) }}
                    </span>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <span class="px-2 py-1 text-xs font-medium rounded-full"
                          [class]="getStatusClass(item.document.verification_status)">
                      {{ getVerificationStatusLabel(item.document.verification_status) }}
                    </span>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                    <a [routerLink]="['/staff', item.document.staff_id]"
                       class="text-indigo-600 hover:text-indigo-900">
                      View Staff
                    </a>
                  </td>
                </tr>
              }
            </tbody>
          </table>
        </div>
      }

      <!-- Empty State -->
      @if (!loading() && documents().length === 0) {
        <div class="bg-white shadow rounded-lg p-12 text-center">
          <msls-icon name="check-circle" class="w-12 h-12 mx-auto text-green-500 mb-4" />
          <h3 class="text-lg font-medium text-gray-900 mb-2">No expiring documents</h3>
          <p class="text-gray-500">No documents are expiring within the selected timeframe.</p>
        </div>
      }
    </div>
  `,
})
export class DocumentListComponent implements OnInit {
  private readonly service = inject(StaffDocumentService);
  private readonly route = inject(ActivatedRoute);

  // State
  documents = signal<ExpiringDocument[]>([]);
  loading = signal(false);
  selectedDays = signal(30);

  // Computed counts
  expiredCount = signal(0);
  urgentCount = signal(0);
  warningCount = signal(0);

  ngOnInit(): void {
    this.loadExpiringDocuments();
  }

  loadExpiringDocuments(): void {
    this.loading.set(true);
    this.service.getExpiringDocuments(this.selectedDays()).subscribe({
      next: (response) => {
        const docs = response.documents || [];
        this.documents.set(docs);

        // Calculate counts
        this.expiredCount.set(docs.filter(d => d.days_to_expiry < 0).length);
        this.urgentCount.set(docs.filter(d => d.days_to_expiry >= 0 && d.days_to_expiry <= 7).length);
        this.warningCount.set(docs.filter(d => d.days_to_expiry > 7 && d.days_to_expiry <= 30).length);

        this.loading.set(false);
      },
      error: () => {
        this.loading.set(false);
      },
    });
  }

  getVerificationStatusLabel = getVerificationStatusLabel;
  formatFileSize = formatFileSize;

  getDaysClass(days: number): string {
    if (days < 0) return 'bg-red-100 text-red-800';
    if (days <= 7) return 'bg-amber-100 text-amber-800';
    if (days <= 30) return 'bg-yellow-100 text-yellow-800';
    return 'bg-blue-100 text-blue-800';
  }

  getStatusClass(status: string): string {
    const classes: Record<string, string> = {
      pending: 'bg-yellow-100 text-yellow-800',
      verified: 'bg-green-100 text-green-800',
      rejected: 'bg-red-100 text-red-800',
    };
    return classes[status] || 'bg-gray-100 text-gray-800';
  }

  formatDaysLeft(days: number): string {
    if (days < 0) return `${Math.abs(days)} days overdue`;
    if (days === 0) return 'Expires today';
    if (days === 1) return '1 day left';
    return `${days} days left`;
  }
}
