/**
 * Daily Summary Page (Story 7.2)
 * View aggregated attendance for all periods in a day
 */

import { Component, inject, signal, computed, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { StudentAttendanceService } from '../../student-attendance.service';
import {
  TeacherClass,
  DailySummary,
  DailySummaryStudent,
  PeriodInfo,
  StudentAttendanceStatus,
  STUDENT_ATTENDANCE_STATUS_SHORT,
  STUDENT_ATTENDANCE_STATUS_COLORS,
} from '../../student-attendance.model';

@Component({
  selector: 'app-daily-summary',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <div class="container mx-auto px-4 py-6">
      <!-- Header -->
      <div class="mb-6">
        <h1 class="text-2xl font-semibold text-gray-900">Daily Attendance Summary</h1>
        <p class="text-sm text-gray-600 mt-1">View period-wise attendance summary for a day</p>
      </div>

      <!-- Loading State -->
      @if (loadingClasses()) {
        <div class="flex justify-center py-12">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        </div>
      }

      <!-- Error State -->
      @if (error()) {
        <div class="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
          <p class="text-red-700">{{ error() }}</p>
          <button (click)="loadClasses()" class="mt-2 text-red-600 hover:text-red-800 underline">
            Try Again
          </button>
        </div>
      }

      <!-- Filters -->
      @if (!loadingClasses() && !error()) {
        <div class="bg-white rounded-lg shadow-sm border p-4 mb-6">
          <div class="flex flex-wrap gap-4 items-end">
            <!-- Class Selector -->
            <div class="flex-1 min-w-[200px]">
              <label class="block text-sm font-medium text-gray-700 mb-1">Select Class</label>
              <select
                [ngModel]="selectedSectionId()"
                (ngModelChange)="onSectionChange($event)"
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
              >
                <option value="">-- Select a class --</option>
                @for (cls of classes(); track cls.sectionId) {
                  <option [value]="cls.sectionId">
                    {{ cls.className }} - {{ cls.sectionName }} ({{ cls.studentCount }} students)
                  </option>
                }
              </select>
            </div>

            <!-- Date Picker -->
            <div class="w-48">
              <label class="block text-sm font-medium text-gray-700 mb-1">Date</label>
              <input
                type="date"
                [ngModel]="selectedDate()"
                (ngModelChange)="onDateChange($event)"
                [max]="todayDate"
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
              />
            </div>
          </div>
        </div>
      }

      <!-- Loading Summary -->
      @if (loadingSummary()) {
        <div class="flex justify-center py-12">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        </div>
      }

      <!-- Summary View -->
      @if (summary() && !loadingSummary()) {
        <!-- Day Info Banner -->
        <div class="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-6">
          <div class="flex items-center justify-between">
            <div>
              <span class="text-blue-800 font-semibold text-lg">{{ summary()!.className }} - {{ summary()!.sectionName }}</span>
              <span class="text-blue-700 ml-4">{{ summary()!.dayName }}, {{ formatDisplayDate(summary()!.date) }}</span>
            </div>
            <div class="text-right">
              <div class="text-blue-800 font-medium">Average Attendance</div>
              <div class="text-2xl font-bold" [class]="getAttendanceColor(summary()!.summary.averageAttendance)">
                {{ summary()!.summary.averageAttendance.toFixed(1) }}%
              </div>
            </div>
          </div>
        </div>

        <!-- Summary Stats -->
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
          <div class="bg-white rounded-lg shadow-sm border p-4">
            <div class="text-gray-600 text-sm">Total Students</div>
            <div class="text-2xl font-semibold text-gray-900">{{ summary()!.summary.totalStudents }}</div>
          </div>
          <div class="bg-white rounded-lg shadow-sm border p-4">
            <div class="text-gray-600 text-sm">Total Periods</div>
            <div class="text-2xl font-semibold text-gray-900">{{ summary()!.summary.totalPeriods }}</div>
          </div>
          <div class="bg-white rounded-lg shadow-sm border p-4">
            <div class="text-gray-600 text-sm">Full Present</div>
            <div class="text-2xl font-semibold text-green-600">{{ summary()!.summary.fullPresentCount }}</div>
          </div>
          <div class="bg-white rounded-lg shadow-sm border p-4">
            <div class="text-gray-600 text-sm">Low Attendance (&lt;50%)</div>
            <div class="text-2xl font-semibold text-red-600">{{ summary()!.summary.absentCount }}</div>
          </div>
        </div>

        <!-- Summary Grid -->
        <div class="bg-white rounded-lg shadow-sm border">
          <!-- Legend -->
          <div class="px-4 py-2 bg-gray-50 border-b flex gap-6 text-sm">
            <span class="flex items-center gap-1">
              <span class="w-4 h-4 rounded bg-green-500"></span> Present
            </span>
            <span class="flex items-center gap-1">
              <span class="w-4 h-4 rounded bg-red-500"></span> Absent
            </span>
            <span class="flex items-center gap-1">
              <span class="w-4 h-4 rounded bg-yellow-500"></span> Late
            </span>
            <span class="flex items-center gap-1">
              <span class="w-4 h-4 rounded bg-blue-500"></span> Half-day
            </span>
            <span class="flex items-center gap-1">
              <span class="w-4 h-4 rounded bg-gray-200"></span> Not marked
            </span>
          </div>

          <!-- Grid -->
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead class="bg-gray-50">
                <tr>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-12">#</th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider min-w-[200px]">Student</th>
                  @for (period of summary()!.periods; track period.periodSlotId) {
                    <th class="px-2 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider w-16" [title]="period.subjectName || period.periodName">
                      {{ period.periodNumber || period.periodName.substring(0, 2) }}
                      @if (period.subjectCode) {
                        <div class="text-xs font-normal normal-case text-gray-400">{{ period.subjectCode }}</div>
                      }
                    </th>
                  }
                  <th class="px-4 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider w-24">%</th>
                  <th class="px-4 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider w-24">Status</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-200">
                @for (student of sortedStudents(); track student.studentId; let i = $index) {
                  <tr
                    [class.bg-red-50]="student.attendancePercentage < 50"
                    [class.bg-yellow-50]="student.attendancePercentage >= 50 && student.attendancePercentage < 75"
                    class="hover:bg-gray-50 transition-colors"
                  >
                    <td class="px-4 py-3 text-sm text-gray-500">{{ i + 1 }}</td>
                    <td class="px-4 py-3">
                      <div class="flex items-center gap-3">
                        @if (student.photoUrl) {
                          <img [src]="student.photoUrl" alt="" class="w-8 h-8 rounded-full object-cover" />
                        } @else {
                          <div class="w-8 h-8 rounded-full bg-gray-200 flex items-center justify-center text-gray-500 text-xs font-medium">
                            {{ getInitials(student.fullName) }}
                          </div>
                        }
                        <div>
                          <div class="font-medium text-gray-900">{{ student.fullName }}</div>
                          <div class="text-xs text-gray-500">{{ student.admissionNumber }}</div>
                        </div>
                      </div>
                    </td>
                    @for (period of summary()!.periods; track period.periodSlotId) {
                      <td class="px-2 py-3 text-center">
                        <span
                          [class]="getStatusCellClass(student.periodStatuses[period.periodSlotId])"
                          class="inline-flex items-center justify-center w-8 h-8 rounded text-xs font-semibold"
                        >
                          {{ getStatusShort(student.periodStatuses[period.periodSlotId]) }}
                        </span>
                      </td>
                    }
                    <td class="px-4 py-3 text-center">
                      <span
                        [class]="getAttendanceColor(student.attendancePercentage)"
                        class="font-semibold"
                      >
                        {{ student.attendancePercentage.toFixed(0) }}%
                      </span>
                    </td>
                    <td class="px-4 py-3 text-center">
                      <span
                        [class]="getOverallStatusClass(student.overallStatus)"
                        class="px-2 py-1 rounded text-xs font-medium"
                      >
                        {{ getOverallStatusLabel(student.overallStatus) }}
                      </span>
                    </td>
                  </tr>
                }
              </tbody>
            </table>
          </div>

          @if (summary()!.students.length === 0) {
            <div class="p-8 text-center text-gray-500">
              No attendance data found for this date
            </div>
          }
        </div>
      }

      <!-- No Data Message -->
      @if (!loadingSummary() && selectedSectionId() && !summary()) {
        <div class="bg-yellow-50 border border-yellow-200 rounded-lg p-6 text-center">
          <p class="text-yellow-700">No period attendance data found for this class on {{ formatDisplayDate(selectedDate()) }}.</p>
          <p class="text-yellow-600 text-sm mt-2">Period-wise attendance may not be enabled or no periods have been marked yet.</p>
        </div>
      }
    </div>
  `,
})
export class DailySummaryComponent implements OnInit {
  private attendanceService = inject(StudentAttendanceService);

  // State
  loadingClasses = signal(true);
  loadingSummary = signal(false);
  error = signal<string | null>(null);

  classes = signal<TeacherClass[]>([]);
  summary = signal<DailySummary | null>(null);
  selectedSectionId = signal<string>('');
  selectedDate = signal<string>(this.formatDate(new Date()));

  todayDate = this.formatDate(new Date());

  sortedStudents = computed(() => {
    const data = this.summary();
    if (!data) return [];
    // Sort by attendance percentage ascending (lowest first to highlight low attendance)
    return [...data.students].sort((a, b) => a.attendancePercentage - b.attendancePercentage);
  });

  ngOnInit(): void {
    this.loadClasses();
  }

  loadClasses(): void {
    this.loadingClasses.set(true);
    this.error.set(null);

    this.attendanceService.getTeacherClasses(this.selectedDate()).subscribe({
      next: (classes) => {
        this.classes.set(classes);
        this.loadingClasses.set(false);

        // Auto-select first class if available
        if (classes.length > 0 && !this.selectedSectionId()) {
          this.onSectionChange(classes[0].sectionId);
        }
      },
      error: (err) => {
        this.loadingClasses.set(false);
        this.error.set(err.error?.message || 'Failed to load classes');
      }
    });
  }

  onSectionChange(sectionId: string): void {
    this.selectedSectionId.set(sectionId);
    if (sectionId) {
      this.loadSummary();
    } else {
      this.summary.set(null);
    }
  }

  onDateChange(date: string): void {
    this.selectedDate.set(date);
    this.loadClasses();
    if (this.selectedSectionId()) {
      this.loadSummary();
    }
  }

  loadSummary(): void {
    if (!this.selectedSectionId()) return;

    this.loadingSummary.set(true);
    this.error.set(null);

    this.attendanceService.getDailySummary(this.selectedSectionId(), this.selectedDate()).subscribe({
      next: (summary) => {
        this.summary.set(summary);
        this.loadingSummary.set(false);
      },
      error: (err) => {
        this.loadingSummary.set(false);
        if (err.status === 404) {
          this.summary.set(null);
        } else {
          this.error.set(err.error?.message || 'Failed to load summary');
        }
      }
    });
  }

  getStatusShort(status?: StudentAttendanceStatus | string): string {
    if (!status) return '-';
    const mapping: Record<string, string> = {
      present: 'P',
      absent: 'A',
      late: 'L',
      half_day: 'H',
    };
    return mapping[status] || '-';
  }

  getStatusCellClass(status?: StudentAttendanceStatus | string): string {
    if (!status) return 'bg-gray-100 text-gray-400';
    switch (status) {
      case 'present':
        return 'bg-green-500 text-white';
      case 'absent':
        return 'bg-red-500 text-white';
      case 'late':
        return 'bg-yellow-500 text-white';
      case 'half_day':
        return 'bg-blue-500 text-white';
      default:
        return 'bg-gray-100 text-gray-400';
    }
  }

  getAttendanceColor(percentage: number): string {
    if (percentage >= 75) return 'text-green-600';
    if (percentage >= 50) return 'text-yellow-600';
    return 'text-red-600';
  }

  getOverallStatusClass(status: StudentAttendanceStatus | string): string {
    switch (status) {
      case 'present':
        return 'bg-green-100 text-green-800';
      case 'absent':
        return 'bg-red-100 text-red-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  }

  getOverallStatusLabel(status: StudentAttendanceStatus | string): string {
    return status === 'present' ? 'Present' : 'Absent';
  }

  getInitials(name: string): string {
    const parts = name.split(' ');
    return (parts[0]?.[0] || '') + (parts[1]?.[0] || '');
  }

  formatDate(date: Date): string {
    return date.toISOString().split('T')[0];
  }

  formatDisplayDate(dateString: string): string {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' });
  }
}
