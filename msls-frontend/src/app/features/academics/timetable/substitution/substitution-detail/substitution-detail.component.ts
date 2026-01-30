import { Component, OnInit, inject, signal } from '@angular/core';
import { CommonModule, DatePipe } from '@angular/common';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { SubstitutionService } from '../substitution.service';
import {
  Substitution,
  SubstitutionStatus,
  SUBSTITUTION_STATUS_CONFIG,
} from '../substitution.model';

@Component({
  selector: 'app-substitution-detail',
  standalone: true,
  imports: [CommonModule, RouterLink, DatePipe],
  template: `
    <div class="max-w-4xl mx-auto space-y-6">
      @if (loading()) {
        <div class="flex justify-center py-12">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
        </div>
      } @else if (!substitution()) {
        <div class="text-center py-12">
          <p class="text-gray-500">Substitution not found</p>
          <a
            routerLink="../"
            class="mt-4 inline-flex items-center text-indigo-600 hover:text-indigo-500"
          >
            Back to list
          </a>
        </div>
      } @else {
        <!-- Header -->
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div class="flex items-center space-x-4">
            <a
              routerLink="../"
              class="text-gray-400 hover:text-gray-500"
            >
              <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
              </svg>
            </a>
            <div>
              <h1 class="text-2xl font-semibold text-gray-900">
                Substitution Details
              </h1>
              <p class="mt-1 text-sm text-gray-500">
                {{ substitution()!.substitutionDate | date:'fullDate' }}
              </p>
            </div>
          </div>

          <div class="flex items-center space-x-3">
            <span [class]="getStatusClass(substitution()!.status)">
              {{ getStatusConfig(substitution()!.status).label }}
            </span>

            @if (substitution()!.status === 'pending') {
              <button
                (click)="confirmSubstitution()"
                class="inline-flex items-center px-4 py-2 border border-transparent rounded-lg shadow-sm text-sm font-medium text-white bg-green-600 hover:bg-green-700"
              >
                Confirm
              </button>
              <button
                (click)="cancelSubstitution()"
                class="inline-flex items-center px-4 py-2 border border-transparent rounded-lg shadow-sm text-sm font-medium text-white bg-red-600 hover:bg-red-700"
              >
                Cancel
              </button>
            }
          </div>
        </div>

        <!-- Main Content -->
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <!-- Left Column - Teachers & Info -->
          <div class="lg:col-span-2 space-y-6">
            <!-- Teachers Card -->
            <div class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
              <div class="px-6 py-4 border-b border-gray-200 bg-gray-50">
                <h2 class="text-lg font-medium text-gray-900">Teachers</h2>
              </div>
              <div class="p-6">
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
                  <!-- Original Teacher -->
                  <div>
                    <p class="text-sm font-medium text-gray-500 mb-2">Absent Teacher</p>
                    <div class="flex items-center">
                      <div class="flex-shrink-0 w-10 h-10 rounded-full bg-red-100 flex items-center justify-center">
                        <svg class="w-5 h-5 text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                        </svg>
                      </div>
                      <div class="ml-3">
                        <p class="text-sm font-medium text-gray-900">{{ substitution()!.originalStaffName }}</p>
                        <p class="text-xs text-gray-500">{{ substitution()!.branchName }}</p>
                      </div>
                    </div>
                  </div>

                  <!-- Substitute Teacher -->
                  <div>
                    <p class="text-sm font-medium text-gray-500 mb-2">Substitute Teacher</p>
                    <div class="flex items-center">
                      <div class="flex-shrink-0 w-10 h-10 rounded-full bg-green-100 flex items-center justify-center">
                        <svg class="w-5 h-5 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                        </svg>
                      </div>
                      <div class="ml-3">
                        <p class="text-sm font-medium text-gray-900">{{ substitution()!.substituteStaffName }}</p>
                        <p class="text-xs text-gray-500">Covering</p>
                      </div>
                    </div>
                  </div>
                </div>

                @if (substitution()!.reason) {
                  <div class="mt-6 pt-6 border-t border-gray-200">
                    <p class="text-sm font-medium text-gray-500 mb-1">Reason for Absence</p>
                    <p class="text-sm text-gray-900">{{ substitution()!.reason }}</p>
                  </div>
                }

                @if (substitution()!.notes) {
                  <div class="mt-4">
                    <p class="text-sm font-medium text-gray-500 mb-1">Notes</p>
                    <p class="text-sm text-gray-900">{{ substitution()!.notes }}</p>
                  </div>
                }
              </div>
            </div>

            <!-- Periods Card -->
            <div class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
              <div class="px-6 py-4 border-b border-gray-200 bg-gray-50">
                <h2 class="text-lg font-medium text-gray-900">Covered Periods</h2>
              </div>
              <div class="divide-y divide-gray-200">
                @for (period of substitution()!.periods; track period.id) {
                  <div class="p-4">
                    <div class="flex items-center justify-between">
                      <div class="flex items-center">
                        <div class="flex-shrink-0 w-12 h-12 rounded-lg bg-indigo-50 flex items-center justify-center">
                          <span class="text-sm font-medium text-indigo-600">{{ period.periodSlotName }}</span>
                        </div>
                        <div class="ml-4">
                          <p class="text-sm font-medium text-gray-900">
                            {{ period.className }} {{ period.sectionName }}
                          </p>
                          <p class="text-sm text-gray-500">{{ period.subjectName }}</p>
                        </div>
                      </div>
                      <div class="text-right">
                        <p class="text-sm text-gray-900">{{ period.startTime }} - {{ period.endTime }}</p>
                        @if (period.roomNumber) {
                          <p class="text-xs text-gray-500">Room: {{ period.roomNumber }}</p>
                        }
                      </div>
                    </div>
                    @if (period.notes) {
                      <p class="mt-2 text-sm text-gray-500 pl-16">{{ period.notes }}</p>
                    }
                  </div>
                }

                @if (!substitution()!.periods?.length) {
                  <div class="p-6 text-center text-sm text-gray-500">
                    No periods assigned
                  </div>
                }
              </div>
            </div>
          </div>

          <!-- Right Column - Meta Info -->
          <div class="space-y-6">
            <!-- Status Timeline -->
            <div class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
              <div class="px-6 py-4 border-b border-gray-200 bg-gray-50">
                <h2 class="text-lg font-medium text-gray-900">Activity</h2>
              </div>
              <div class="p-6">
                <div class="flow-root">
                  <ul class="-mb-8">
                    <!-- Created -->
                    <li>
                      <div class="relative pb-8">
                        <span class="absolute left-4 top-4 -ml-px h-full w-0.5 bg-gray-200"></span>
                        <div class="relative flex space-x-3">
                          <div>
                            <span class="h-8 w-8 rounded-full bg-blue-100 flex items-center justify-center">
                              <svg class="h-4 w-4 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
                              </svg>
                            </span>
                          </div>
                          <div class="min-w-0 flex-1 pt-1.5">
                            <p class="text-sm text-gray-900">Created</p>
                            <p class="text-xs text-gray-500">
                              {{ substitution()!.createdAt | date:'medium' }}
                              @if (substitution()!.createdByName) {
                                by {{ substitution()!.createdByName }}
                              }
                            </p>
                          </div>
                        </div>
                      </div>
                    </li>

                    <!-- Approved (if applicable) -->
                    @if (substitution()!.approvedAt) {
                      <li>
                        <div class="relative pb-8">
                          @if (substitution()!.status !== 'confirmed') {
                            <span class="absolute left-4 top-4 -ml-px h-full w-0.5 bg-gray-200"></span>
                          }
                          <div class="relative flex space-x-3">
                            <div>
                              <span class="h-8 w-8 rounded-full bg-green-100 flex items-center justify-center">
                                <svg class="h-4 w-4 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                                </svg>
                              </span>
                            </div>
                            <div class="min-w-0 flex-1 pt-1.5">
                              <p class="text-sm text-gray-900">Confirmed</p>
                              <p class="text-xs text-gray-500">
                                {{ substitution()!.approvedAt | date:'medium' }}
                                @if (substitution()!.approvedByName) {
                                  by {{ substitution()!.approvedByName }}
                                }
                              </p>
                            </div>
                          </div>
                        </div>
                      </li>
                    }

                    <!-- Current Status -->
                    @if (substitution()!.status === 'cancelled') {
                      <li>
                        <div class="relative">
                          <div class="relative flex space-x-3">
                            <div>
                              <span class="h-8 w-8 rounded-full bg-red-100 flex items-center justify-center">
                                <svg class="h-4 w-4 text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                                </svg>
                              </span>
                            </div>
                            <div class="min-w-0 flex-1 pt-1.5">
                              <p class="text-sm text-gray-900">Cancelled</p>
                              <p class="text-xs text-gray-500">
                                {{ substitution()!.updatedAt | date:'medium' }}
                              </p>
                            </div>
                          </div>
                        </div>
                      </li>
                    }

                    @if (substitution()!.status === 'completed') {
                      <li>
                        <div class="relative">
                          <div class="relative flex space-x-3">
                            <div>
                              <span class="h-8 w-8 rounded-full bg-green-100 flex items-center justify-center">
                                <svg class="h-4 w-4 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                                </svg>
                              </span>
                            </div>
                            <div class="min-w-0 flex-1 pt-1.5">
                              <p class="text-sm text-gray-900">Completed</p>
                              <p class="text-xs text-gray-500">
                                {{ substitution()!.updatedAt | date:'medium' }}
                              </p>
                            </div>
                          </div>
                        </div>
                      </li>
                    }
                  </ul>
                </div>
              </div>
            </div>

            <!-- Quick Info -->
            <div class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
              <div class="px-6 py-4 border-b border-gray-200 bg-gray-50">
                <h2 class="text-lg font-medium text-gray-900">Details</h2>
              </div>
              <div class="p-6 space-y-4">
                <div>
                  <p class="text-sm font-medium text-gray-500">Branch</p>
                  <p class="text-sm text-gray-900">{{ substitution()!.branchName }}</p>
                </div>
                <div>
                  <p class="text-sm font-medium text-gray-500">Date</p>
                  <p class="text-sm text-gray-900">{{ substitution()!.substitutionDate | date:'fullDate' }}</p>
                </div>
                <div>
                  <p class="text-sm font-medium text-gray-500">Total Periods</p>
                  <p class="text-sm text-gray-900">{{ substitution()!.periods?.length || 0 }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      }
    </div>
  `,
})
export class SubstitutionDetailComponent implements OnInit {
  private route = inject(ActivatedRoute);
  private router = inject(Router);
  private substitutionService = inject(SubstitutionService);

  substitution = this.substitutionService.currentSubstitution;
  loading = this.substitutionService.loading;

  ngOnInit(): void {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.substitutionService.getSubstitution(id).subscribe();
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
    return `inline-flex items-center px-3 py-1 rounded-full text-sm font-medium ${colorMap[config.color]}`;
  }

  confirmSubstitution(): void {
    const sub = this.substitution();
    if (sub && confirm(`Confirm substitution for ${sub.originalStaffName}?`)) {
      this.substitutionService.confirmSubstitution(sub.id).subscribe({
        error: (err) => alert('Failed to confirm: ' + err.message),
      });
    }
  }

  cancelSubstitution(): void {
    const sub = this.substitution();
    if (sub && confirm(`Cancel substitution for ${sub.originalStaffName}?`)) {
      this.substitutionService.cancelSubstitution(sub.id).subscribe({
        next: () => this.router.navigate(['..'], { relativeTo: this.route }),
        error: (err) => alert('Failed to cancel: ' + err.message),
      });
    }
  }
}
