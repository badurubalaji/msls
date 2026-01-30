import { Component, OnInit, inject, signal, computed } from '@angular/core';
import { CommonModule, DatePipe } from '@angular/common';
import { RouterLink } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { SubstitutionService } from '../substitution.service';
import { BranchService } from '../../../../admin/branches/branch.service';
import {
  Substitution,
  SubstitutionStatus,
  SubstitutionFilter,
  SUBSTITUTION_STATUS_CONFIG,
} from '../substitution.model';

@Component({
  selector: 'app-substitution-list',
  standalone: true,
  imports: [CommonModule, RouterLink, FormsModule, DatePipe],
  template: `
    <div class="space-y-6">
      <!-- Header -->
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 class="text-2xl font-semibold text-gray-900">Substitutions</h1>
          <p class="mt-1 text-sm text-gray-500">
            Manage teacher substitutions and coverage
          </p>
        </div>
        <a
          routerLink="new"
          class="inline-flex items-center px-4 py-2 border border-transparent rounded-lg shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
        >
          <svg class="w-5 h-5 mr-2 -ml-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          New Substitution
        </a>
      </div>

      <!-- Filters -->
      <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-4">
          <!-- Branch Filter -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Branch</label>
            <select
              [(ngModel)]="filters.branchId"
              (ngModelChange)="applyFilters()"
              class="w-full rounded-lg border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
            >
              <option value="">All Branches</option>
              @for (branch of branches(); track branch.id) {
                <option [value]="branch.id">{{ branch.name }}</option>
              }
            </select>
          </div>

          <!-- Status Filter -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Status</label>
            <select
              [(ngModel)]="filters.status"
              (ngModelChange)="applyFilters()"
              class="w-full rounded-lg border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
            >
              <option value="">All Statuses</option>
              <option value="pending">Pending</option>
              <option value="confirmed">Confirmed</option>
              <option value="completed">Completed</option>
              <option value="cancelled">Cancelled</option>
            </select>
          </div>

          <!-- Start Date -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">From Date</label>
            <input
              type="date"
              [(ngModel)]="filters.startDate"
              (ngModelChange)="applyFilters()"
              class="w-full rounded-lg border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
            />
          </div>

          <!-- End Date -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">To Date</label>
            <input
              type="date"
              [(ngModel)]="filters.endDate"
              (ngModelChange)="applyFilters()"
              class="w-full rounded-lg border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
            />
          </div>

          <!-- Clear Filters -->
          <div class="flex items-end">
            <button
              (click)="clearFilters()"
              class="w-full px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              Clear Filters
            </button>
          </div>
        </div>
      </div>

      <!-- Loading State -->
      @if (loading()) {
        <div class="flex justify-center py-12">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
        </div>
      } @else {
        <!-- Stats Summary -->
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
          @for (stat of statusStats(); track stat.status) {
            <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
              <div class="flex items-center">
                <div [class]="'flex-shrink-0 w-10 h-10 rounded-lg flex items-center justify-center bg-' + stat.color + '-100'">
                  <svg [class]="'w-5 h-5 text-' + stat.color + '-600'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    @switch (stat.icon) {
                      @case ('clock') {
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                      }
                      @case ('check-circle') {
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                      }
                      @case ('badge-check') {
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z" />
                      }
                      @case ('x-circle') {
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
                      }
                    }
                  </svg>
                </div>
                <div class="ml-4">
                  <p class="text-sm font-medium text-gray-500">{{ stat.label }}</p>
                  <p class="text-2xl font-semibold text-gray-900">{{ stat.count }}</p>
                </div>
              </div>
            </div>
          }
        </div>

        <!-- Substitutions Table -->
        <div class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
          @if (substitutions().length === 0) {
            <div class="text-center py-12">
              <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
              </svg>
              <h3 class="mt-2 text-sm font-medium text-gray-900">No substitutions found</h3>
              <p class="mt-1 text-sm text-gray-500">Get started by creating a new substitution.</p>
              <div class="mt-6">
                <a
                  routerLink="new"
                  class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-lg text-white bg-indigo-600 hover:bg-indigo-700"
                >
                  <svg class="w-5 h-5 mr-2 -ml-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
                  </svg>
                  New Substitution
                </a>
              </div>
            </div>
          } @else {
            <div class="overflow-x-auto">
              <table class="min-w-full divide-y divide-gray-200">
                <thead class="bg-gray-50">
                  <tr>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Date</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Original Teacher</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Substitute</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Periods</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Reason</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                    <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
                  </tr>
                </thead>
                <tbody class="bg-white divide-y divide-gray-200">
                  @for (sub of substitutions(); track sub.id) {
                    <tr class="hover:bg-gray-50">
                      <td class="px-6 py-4 whitespace-nowrap">
                        <div class="text-sm font-medium text-gray-900">
                          {{ sub.substitutionDate | date:'mediumDate' }}
                        </div>
                        <div class="text-xs text-gray-500">
                          {{ sub.substitutionDate | date:'EEEE' }}
                        </div>
                      </td>
                      <td class="px-6 py-4 whitespace-nowrap">
                        <div class="text-sm text-gray-900">{{ sub.originalStaffName }}</div>
                        <div class="text-xs text-gray-500">{{ sub.branchName }}</div>
                      </td>
                      <td class="px-6 py-4 whitespace-nowrap">
                        <div class="text-sm text-gray-900">{{ sub.substituteStaffName }}</div>
                      </td>
                      <td class="px-6 py-4 whitespace-nowrap">
                        <div class="flex flex-wrap gap-1">
                          @for (period of sub.periods?.slice(0, 3); track period.id) {
                            <span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-800">
                              {{ period.periodSlotName || 'Period ' + ($index + 1) }}
                            </span>
                          }
                          @if ((sub.periods?.length || 0) > 3) {
                            <span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-800">
                              +{{ (sub.periods?.length || 0) - 3 }} more
                            </span>
                          }
                        </div>
                      </td>
                      <td class="px-6 py-4">
                        <div class="text-sm text-gray-900 max-w-xs truncate">{{ sub.reason || '-' }}</div>
                      </td>
                      <td class="px-6 py-4 whitespace-nowrap">
                        <span [class]="getStatusClass(sub.status)">
                          {{ getStatusConfig(sub.status).label }}
                        </span>
                      </td>
                      <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                        <div class="flex items-center justify-end space-x-2">
                          <a
                            [routerLink]="[sub.id]"
                            class="text-indigo-600 hover:text-indigo-900"
                          >
                            View
                          </a>
                          @if (sub.status === 'pending') {
                            <button
                              (click)="confirmSubstitution(sub)"
                              class="text-green-600 hover:text-green-900"
                            >
                              Confirm
                            </button>
                            <button
                              (click)="cancelSubstitution(sub)"
                              class="text-red-600 hover:text-red-900"
                            >
                              Cancel
                            </button>
                          }
                        </div>
                      </td>
                    </tr>
                  }
                </tbody>
              </table>
            </div>

            <!-- Pagination -->
            @if (total() > pageSize) {
              <div class="bg-white px-4 py-3 flex items-center justify-between border-t border-gray-200 sm:px-6">
                <div class="flex-1 flex justify-between sm:hidden">
                  <button
                    (click)="previousPage()"
                    [disabled]="currentPage() === 1"
                    class="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50"
                  >
                    Previous
                  </button>
                  <button
                    (click)="nextPage()"
                    [disabled]="currentPage() >= totalPages()"
                    class="ml-3 relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50"
                  >
                    Next
                  </button>
                </div>
                <div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
                  <div>
                    <p class="text-sm text-gray-700">
                      Showing
                      <span class="font-medium">{{ (currentPage() - 1) * pageSize + 1 }}</span>
                      to
                      <span class="font-medium">{{ Math.min(currentPage() * pageSize, total()) }}</span>
                      of
                      <span class="font-medium">{{ total() }}</span>
                      results
                    </p>
                  </div>
                  <div>
                    <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px">
                      <button
                        (click)="previousPage()"
                        [disabled]="currentPage() === 1"
                        class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50"
                      >
                        <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
                        </svg>
                      </button>
                      <span class="relative inline-flex items-center px-4 py-2 border border-gray-300 bg-white text-sm font-medium text-gray-700">
                        Page {{ currentPage() }} of {{ totalPages() }}
                      </span>
                      <button
                        (click)="nextPage()"
                        [disabled]="currentPage() >= totalPages()"
                        class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50"
                      >
                        <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                        </svg>
                      </button>
                    </nav>
                  </div>
                </div>
              </div>
            }
          }
        </div>
      }
    </div>
  `,
})
export class SubstitutionListComponent implements OnInit {
  private substitutionService = inject(SubstitutionService);
  private branchService = inject(BranchService);

  // Expose Math for template
  protected Math = Math;

  // State
  substitutions = this.substitutionService.substitutions;
  total = this.substitutionService.total;
  loading = this.substitutionService.loading;
  branches = signal<{ id: string; name: string }[]>([]);

  // Pagination
  pageSize = 20;
  currentPage = signal(1);
  totalPages = computed(() => Math.ceil(this.total() / this.pageSize));

  // Filters
  filters: SubstitutionFilter = {
    branchId: '',
    status: '',
    startDate: '',
    endDate: '',
  };

  // Stats
  statusStats = computed(() => {
    const subs = this.substitutions();
    const counts: Record<SubstitutionStatus, number> = {
      pending: 0,
      confirmed: 0,
      completed: 0,
      cancelled: 0,
    };

    subs.forEach((s) => {
      counts[s.status]++;
    });

    return Object.entries(SUBSTITUTION_STATUS_CONFIG).map(([status, config]) => ({
      status,
      ...config,
      count: counts[status as SubstitutionStatus],
    }));
  });

  ngOnInit(): void {
    this.branchService.getBranches().subscribe((branches) => this.branches.set(branches));
    this.loadSubstitutions();
  }

  loadSubstitutions(): void {
    const filter: SubstitutionFilter = {
      limit: this.pageSize,
      offset: (this.currentPage() - 1) * this.pageSize,
    };

    if (this.filters.branchId) filter.branchId = this.filters.branchId;
    if (this.filters.status) filter.status = this.filters.status;
    if (this.filters.startDate) filter.startDate = this.filters.startDate;
    if (this.filters.endDate) filter.endDate = this.filters.endDate;

    this.substitutionService.loadSubstitutions(filter).subscribe();
  }

  applyFilters(): void {
    this.currentPage.set(1);
    this.loadSubstitutions();
  }

  clearFilters(): void {
    this.filters = {
      branchId: '',
      status: '',
      startDate: '',
      endDate: '',
    };
    this.applyFilters();
  }

  previousPage(): void {
    if (this.currentPage() > 1) {
      this.currentPage.update((p) => p - 1);
      this.loadSubstitutions();
    }
  }

  nextPage(): void {
    if (this.currentPage() < this.totalPages()) {
      this.currentPage.update((p) => p + 1);
      this.loadSubstitutions();
    }
  }

  getStatusConfig(status: SubstitutionStatus) {
    return SUBSTITUTION_STATUS_CONFIG[status];
  }

  getStatusClass(status: SubstitutionStatus): string {
    const colorMap: Record<string, string> = {
      amber: 'bg-amber-100 text-amber-800',
      blue: 'bg-blue-100 text-blue-800',
      green: 'bg-green-100 text-green-800',
      red: 'bg-red-100 text-red-800',
    };
    const config = this.getStatusConfig(status);
    return `inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${colorMap[config.color]}`;
  }

  confirmSubstitution(sub: Substitution): void {
    if (confirm(`Confirm substitution for ${sub.originalStaffName}?`)) {
      this.substitutionService.confirmSubstitution(sub.id).subscribe({
        next: () => this.loadSubstitutions(),
        error: (err) => alert('Failed to confirm substitution: ' + err.message),
      });
    }
  }

  cancelSubstitution(sub: Substitution): void {
    if (confirm(`Cancel substitution for ${sub.originalStaffName}?`)) {
      this.substitutionService.cancelSubstitution(sub.id).subscribe({
        next: () => this.loadSubstitutions(),
        error: (err) => alert('Failed to cancel substitution: ' + err.message),
      });
    }
  }
}
