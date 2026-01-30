/**
 * Period Attendance Page (Story 7.2)
 * Teachers can mark period-wise attendance for their classes
 */

import { Component, inject, signal, computed, OnInit, HostListener } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';

import { StudentAttendanceService } from '../../student-attendance.service';
import {
  TeacherClass,
  PeriodInfo,
  TeacherPeriodsResponse,
  PeriodAttendance,
  StudentForPeriodAttendance,
  StudentAttendanceStatus,
  PeriodAttendanceRecordRequest,
  MarkPeriodAttendanceRequest,
  STUDENT_ATTENDANCE_DOT_COLORS,
} from '../../student-attendance.model';

@Component({
  selector: 'app-period-attendance',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <div class="container mx-auto px-4 py-6">
      <!-- Header -->
      <div class="mb-6">
        <h1 class="text-2xl font-semibold text-gray-900">Period-wise Attendance</h1>
        <p class="text-sm text-gray-600 mt-1">Mark attendance for each teaching period</p>
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

      <!-- Class & Period Selection -->
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

            <!-- Period Selector -->
            @if (periods() && periods()!.length > 0) {
              <div class="flex-1 min-w-[200px]">
                <label class="block text-sm font-medium text-gray-700 mb-1">Select Period</label>
                <select
                  [ngModel]="selectedPeriodId()"
                  (ngModelChange)="onPeriodChange($event)"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                >
                  <option value="">-- Select a period --</option>
                  @for (period of periods(); track period.periodSlotId) {
                    <option [value]="period.periodSlotId">
                      {{ period.periodName }} ({{ period.startTime }} - {{ period.endTime }})
                      {{ period.subjectName ? '- ' + period.subjectName : '' }}
                      {{ period.isMarked ? '✓' : '' }}
                    </option>
                  }
                </select>
              </div>
            }

            <!-- Submit Button -->
            @if (periodAttendance() && periodAttendance()!.students.length > 0) {
              <button
                (click)="submitAttendance()"
                [disabled]="submitting()"
                class="px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
              >
                @if (submitting()) {
                  <span class="animate-spin h-4 w-4 border-2 border-white border-t-transparent rounded-full"></span>
                }
                {{ periodAttendance()!.isMarked ? 'Update' : 'Submit' }} Attendance
              </button>
            }
          </div>
        </div>
      }

      <!-- Period Info Banner -->
      @if (periodsResponse() && selectedPeriodId()) {
        <div class="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-6">
          <div class="flex items-center gap-4">
            <div>
              <span class="text-blue-800 font-medium">{{ periodsResponse()!.dayName }}, {{ formatDisplayDate(selectedDate()) }}</span>
            </div>
            @if (periodAttendance()?.subjectName) {
              <div class="text-blue-700">
                <span class="font-medium">Subject:</span> {{ periodAttendance()!.subjectName }}
                @if (periodAttendance()!.subjectCode) {
                  ({{ periodAttendance()!.subjectCode }})
                }
              </div>
            }
          </div>
        </div>
      }

      <!-- Period Tabs -->
      @if (periods() && periods()!.length > 0) {
        <div class="mb-6 overflow-x-auto">
          <div class="flex gap-2 min-w-max pb-2">
            @for (period of periods(); track period.periodSlotId) {
              <button
                (click)="onPeriodChange(period.periodSlotId)"
                [class]="getPeriodTabClass(period)"
                class="px-4 py-2 rounded-lg text-sm font-medium transition-colors whitespace-nowrap"
              >
                {{ period.periodName }}
                @if (period.subjectCode) {
                  <span class="text-xs opacity-75">({{ period.subjectCode }})</span>
                }
                @if (period.isMarked) {
                  <span class="ml-1 text-green-600">✓</span>
                }
              </button>
            }
          </div>
        </div>
      }

      <!-- Loading Periods -->
      @if (loadingPeriods()) {
        <div class="flex justify-center py-8">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        </div>
      }

      <!-- No Timetable Message -->
      @if (!loadingPeriods() && selectedSectionId() && periods().length === 0) {
        <div class="bg-yellow-50 border border-yellow-200 rounded-lg p-6 text-center">
          <p class="text-yellow-700">No timetable found for this class on {{ formatDisplayDate(selectedDate()) }}.</p>
          <p class="text-yellow-600 text-sm mt-2">Please ensure a timetable is published for this section.</p>
        </div>
      }

      <!-- Attendance Grid -->
      @if (loadingAttendance()) {
        <div class="flex justify-center py-12">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        </div>
      }

      @if (periodAttendance() && !loadingAttendance()) {
        <div class="bg-white rounded-lg shadow-sm border">
          <!-- Toolbar -->
          <div class="p-4 border-b flex flex-wrap gap-4 items-center justify-between">
            <div class="flex gap-2">
              <button
                (click)="markAllPresent()"
                class="px-3 py-1.5 text-sm bg-green-100 text-green-700 rounded hover:bg-green-200"
              >
                Mark All Present
              </button>
              <button
                (click)="markAllAbsent()"
                class="px-3 py-1.5 text-sm bg-red-100 text-red-700 rounded hover:bg-red-200"
              >
                Mark All Absent
              </button>
            </div>

            <!-- Search -->
            <div class="relative">
              <input
                type="text"
                [(ngModel)]="searchQuery"
                placeholder="Search student..."
                class="pl-8 pr-4 py-1.5 text-sm border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
              />
              <svg class="absolute left-2.5 top-2 h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
              </svg>
            </div>

            <!-- Summary -->
            <div class="flex gap-4 text-sm">
              <span class="text-gray-600">Total: {{ periodAttendance()!.summary.total }}</span>
              <span class="text-green-600">Present: {{ presentCount() }}</span>
              <span class="text-red-600">Absent: {{ absentCount() }}</span>
              <span class="text-yellow-600">Late: {{ lateCount() }}</span>
              <span class="text-blue-600">Half-day: {{ halfDayCount() }}</span>
            </div>
          </div>

          <!-- Status Info -->
          @if (periodAttendance()!.isMarked && !periodAttendance()!.canEdit) {
            <div class="px-4 py-2 bg-yellow-50 border-b text-yellow-800 text-sm">
              Attendance was marked at {{ formatDateTime(periodAttendance()!.markedAt) }} by {{ periodAttendance()!.markedByName || 'Unknown' }}.
              Edit window has expired.
            </div>
          }

          <!-- Keyboard Shortcuts Info -->
          <div class="px-4 py-2 bg-gray-50 border-b text-gray-600 text-xs">
            Keyboard shortcuts: <kbd class="px-1.5 py-0.5 bg-gray-200 rounded">P</kbd> Present,
            <kbd class="px-1.5 py-0.5 bg-gray-200 rounded">A</kbd> Absent,
            <kbd class="px-1.5 py-0.5 bg-gray-200 rounded">L</kbd> Late,
            <kbd class="px-1.5 py-0.5 bg-gray-200 rounded">H</kbd> Half-day |
            <kbd class="px-1.5 py-0.5 bg-gray-200 rounded">↑</kbd>/<kbd class="px-1.5 py-0.5 bg-gray-200 rounded">↓</kbd> Navigate
          </div>

          <!-- Student Grid -->
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead class="bg-gray-50">
                <tr>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-12">#</th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-16">Photo</th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Name</th>
                  <th class="px-4 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider w-40">Status</th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-48">Remarks</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-200">
                @for (student of filteredStudents(); track student.studentId; let i = $index) {
                  <tr
                    [class.bg-blue-50]="selectedIndex() === i"
                    class="hover:bg-gray-50 transition-colors cursor-pointer"
                    (click)="selectStudent(i)"
                  >
                    <td class="px-4 py-3 text-sm text-gray-500">{{ i + 1 }}</td>
                    <td class="px-4 py-3">
                      @if (student.photoUrl) {
                        <img [src]="student.photoUrl" alt="" class="w-10 h-10 rounded-full object-cover" />
                      } @else {
                        <div class="w-10 h-10 rounded-full bg-gray-200 flex items-center justify-center text-gray-500 text-sm font-medium">
                          {{ getInitials(student) }}
                        </div>
                      }
                    </td>
                    <td class="px-4 py-3">
                      <div class="font-medium text-gray-900">{{ student.fullName }}</div>
                      <div class="text-xs text-gray-500">{{ student.admissionNumber }}</div>
                    </td>
                    <td class="px-4 py-3">
                      <div class="flex justify-center gap-1">
                        @for (status of statusOptions; track status.value) {
                          <button
                            (click)="setStudentStatus(student, status.value); $event.stopPropagation()"
                            [class]="getStatusButtonClass(student, status.value)"
                            [disabled]="!canEdit()"
                            class="w-8 h-8 rounded-md text-sm font-semibold transition-colors disabled:opacity-50"
                            [title]="status.label"
                          >
                            {{ status.short }}
                          </button>
                        }
                      </div>
                    </td>
                    <td class="px-4 py-3">
                      <input
                        type="text"
                        [(ngModel)]="student.remarks"
                        [disabled]="!canEdit()"
                        placeholder="Add note..."
                        class="w-full px-2 py-1 text-sm border border-gray-200 rounded focus:ring-blue-500 focus:border-blue-500 disabled:bg-gray-100"
                        (click)="$event.stopPropagation()"
                      />
                    </td>
                  </tr>
                }
              </tbody>
            </table>
          </div>

          @if (filteredStudents().length === 0) {
            <div class="p-8 text-center text-gray-500">
              @if (searchQuery) {
                No students found matching "{{ searchQuery }}"
              } @else {
                No students in this class
              }
            </div>
          }
        </div>
      }

      <!-- Success Message -->
      @if (successMessage()) {
        <div class="fixed bottom-4 right-4 bg-green-100 border border-green-400 text-green-700 px-6 py-3 rounded-lg shadow-lg">
          {{ successMessage() }}
        </div>
      }
    </div>
  `,
  styles: [`
    kbd {
      font-family: monospace;
    }
  `]
})
export class PeriodAttendanceComponent implements OnInit {
  private router = inject(Router);
  private attendanceService = inject(StudentAttendanceService);

  // State
  loadingClasses = signal(true);
  loadingPeriods = signal(false);
  loadingAttendance = signal(false);
  submitting = signal(false);
  error = signal<string | null>(null);
  successMessage = signal<string | null>(null);

  classes = signal<TeacherClass[]>([]);
  periodsResponse = signal<TeacherPeriodsResponse | null>(null);
  periodAttendance = signal<PeriodAttendance | null>(null);
  selectedSectionId = signal<string>('');
  selectedPeriodId = signal<string>('');
  selectedDate = signal<string>(this.formatDate(new Date()));
  selectedIndex = signal<number>(-1);

  searchQuery = '';
  todayDate = this.formatDate(new Date());

  statusOptions: { value: StudentAttendanceStatus; short: string; label: string }[] = [
    { value: 'present', short: 'P', label: 'Present' },
    { value: 'absent', short: 'A', label: 'Absent' },
    { value: 'late', short: 'L', label: 'Late' },
    { value: 'half_day', short: 'H', label: 'Half Day' },
  ];

  periods = computed(() => this.periodsResponse()?.periods || []);

  // Computed values for summary
  presentCount = computed(() => {
    const attendance = this.periodAttendance();
    if (!attendance) return 0;
    return attendance.students.filter(s => s.status === 'present').length;
  });

  absentCount = computed(() => {
    const attendance = this.periodAttendance();
    if (!attendance) return 0;
    return attendance.students.filter(s => s.status === 'absent').length;
  });

  lateCount = computed(() => {
    const attendance = this.periodAttendance();
    if (!attendance) return 0;
    return attendance.students.filter(s => s.status === 'late').length;
  });

  halfDayCount = computed(() => {
    const attendance = this.periodAttendance();
    if (!attendance) return 0;
    return attendance.students.filter(s => s.status === 'half_day').length;
  });

  filteredStudents = computed(() => {
    const attendance = this.periodAttendance();
    if (!attendance) return [];

    if (!this.searchQuery.trim()) {
      return attendance.students;
    }

    const query = this.searchQuery.toLowerCase();
    return attendance.students.filter(
      s =>
        s.fullName.toLowerCase().includes(query) ||
        s.admissionNumber.toLowerCase().includes(query)
    );
  });

  canEdit = computed(() => {
    const attendance = this.periodAttendance();
    if (!attendance) return false;
    return !attendance.isMarked || attendance.canEdit;
  });

  ngOnInit(): void {
    this.loadClasses();
  }

  @HostListener('window:keydown', ['$event'])
  handleKeyDown(event: KeyboardEvent): void {
    const attendance = this.periodAttendance();
    if (!attendance || !this.canEdit()) return;

    const students = this.filteredStudents();
    const currentIndex = this.selectedIndex();

    switch (event.key.toLowerCase()) {
      case 'arrowdown':
        event.preventDefault();
        if (currentIndex < students.length - 1) {
          this.selectedIndex.set(currentIndex + 1);
        }
        break;
      case 'arrowup':
        event.preventDefault();
        if (currentIndex > 0) {
          this.selectedIndex.set(currentIndex - 1);
        }
        break;
      case 'p':
        if (currentIndex >= 0 && currentIndex < students.length) {
          this.setStudentStatus(students[currentIndex], 'present');
        }
        break;
      case 'a':
        if (currentIndex >= 0 && currentIndex < students.length) {
          this.setStudentStatus(students[currentIndex], 'absent');
        }
        break;
      case 'l':
        if (currentIndex >= 0 && currentIndex < students.length) {
          this.setStudentStatus(students[currentIndex], 'late');
        }
        break;
      case 'h':
        if (currentIndex >= 0 && currentIndex < students.length) {
          this.setStudentStatus(students[currentIndex], 'half_day');
        }
        break;
    }
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
    this.selectedPeriodId.set('');
    this.selectedIndex.set(-1);
    this.periodAttendance.set(null);

    if (sectionId) {
      this.loadPeriods();
    } else {
      this.periodsResponse.set(null);
    }
  }

  onDateChange(date: string): void {
    this.selectedDate.set(date);
    this.selectedPeriodId.set('');
    this.periodAttendance.set(null);

    this.loadClasses();
    if (this.selectedSectionId()) {
      this.loadPeriods();
    }
  }

  onPeriodChange(periodId: string): void {
    this.selectedPeriodId.set(periodId);
    this.selectedIndex.set(-1);

    if (periodId) {
      this.loadPeriodAttendance();
    } else {
      this.periodAttendance.set(null);
    }
  }

  loadPeriods(): void {
    if (!this.selectedSectionId()) return;

    this.loadingPeriods.set(true);
    this.error.set(null);

    this.attendanceService.getTeacherPeriods(this.selectedSectionId(), this.selectedDate()).subscribe({
      next: (response) => {
        this.periodsResponse.set(response);
        this.loadingPeriods.set(false);

        // Auto-select first unmarked period if available
        const firstUnmarked = response.periods.find(p => !p.isMarked);
        if (firstUnmarked) {
          this.onPeriodChange(firstUnmarked.periodSlotId);
        } else if (response.periods.length > 0) {
          this.onPeriodChange(response.periods[0].periodSlotId);
        }
      },
      error: (err) => {
        this.loadingPeriods.set(false);
        if (err.status === 404) {
          // No timetable - show empty state instead of error
          this.periodsResponse.set({ periods: [] } as any);
        } else {
          this.error.set(err.error?.message || 'Failed to load periods');
        }
      }
    });
  }

  loadPeriodAttendance(): void {
    if (!this.selectedSectionId() || !this.selectedPeriodId()) return;

    this.loadingAttendance.set(true);
    this.error.set(null);

    this.attendanceService.getPeriodAttendance(
      this.selectedPeriodId(),
      this.selectedSectionId(),
      this.selectedDate()
    ).subscribe({
      next: (attendance) => {
        // Set default status to 'present' for students without status
        attendance.students = attendance.students.map(s => ({
          ...s,
          status: s.status || 'present'
        }));
        this.periodAttendance.set(attendance);
        this.loadingAttendance.set(false);
      },
      error: (err) => {
        this.loadingAttendance.set(false);
        this.error.set(err.error?.message || 'Failed to load period attendance');
      }
    });
  }

  selectStudent(index: number): void {
    this.selectedIndex.set(index);
  }

  setStudentStatus(student: StudentForPeriodAttendance, status: StudentAttendanceStatus): void {
    if (!this.canEdit()) return;
    student.status = status;
    // Trigger change detection
    this.periodAttendance.update(a => a ? { ...a } : null);
  }

  markAllPresent(): void {
    if (!this.canEdit()) return;
    const attendance = this.periodAttendance();
    if (attendance) {
      attendance.students.forEach(s => s.status = 'present');
      this.periodAttendance.set({ ...attendance });
    }
  }

  markAllAbsent(): void {
    if (!this.canEdit()) return;
    const attendance = this.periodAttendance();
    if (attendance) {
      attendance.students.forEach(s => s.status = 'absent');
      this.periodAttendance.set({ ...attendance });
    }
  }

  submitAttendance(): void {
    const attendance = this.periodAttendance();
    if (!attendance || this.submitting()) return;

    this.submitting.set(true);
    this.error.set(null);

    const records: PeriodAttendanceRecordRequest[] = attendance.students.map(s => ({
      studentId: s.studentId,
      status: (s.status || 'present') as StudentAttendanceStatus,
      remarks: s.remarks || undefined,
    }));

    const request: MarkPeriodAttendanceRequest = {
      sectionId: this.selectedSectionId(),
      date: this.selectedDate(),
      records,
    };

    this.attendanceService.markPeriodAttendance(this.selectedPeriodId(), request).subscribe({
      next: (result) => {
        this.submitting.set(false);
        this.successMessage.set(result.message);

        // Reload to get updated state
        this.loadPeriodAttendance();
        this.loadPeriods();

        // Clear success message after 3 seconds
        setTimeout(() => this.successMessage.set(null), 3000);
      },
      error: (err) => {
        this.submitting.set(false);
        this.error.set(err.error?.message || 'Failed to save attendance');
      }
    });
  }

  getPeriodTabClass(period: PeriodInfo): string {
    const isSelected = this.selectedPeriodId() === period.periodSlotId;
    if (isSelected) {
      return 'bg-blue-600 text-white';
    }
    if (period.isMarked) {
      return 'bg-green-100 text-green-800 hover:bg-green-200';
    }
    return 'bg-gray-100 text-gray-700 hover:bg-gray-200';
  }

  getStatusButtonClass(student: StudentForPeriodAttendance, status: StudentAttendanceStatus): string {
    const isActive = student.status === status;
    const baseClasses = 'border-2 ';

    switch (status) {
      case 'present':
        return baseClasses + (isActive ? 'bg-green-500 text-white border-green-500' : 'border-green-300 text-green-600 hover:bg-green-50');
      case 'absent':
        return baseClasses + (isActive ? 'bg-red-500 text-white border-red-500' : 'border-red-300 text-red-600 hover:bg-red-50');
      case 'late':
        return baseClasses + (isActive ? 'bg-yellow-500 text-white border-yellow-500' : 'border-yellow-300 text-yellow-600 hover:bg-yellow-50');
      case 'half_day':
        return baseClasses + (isActive ? 'bg-blue-500 text-white border-blue-500' : 'border-blue-300 text-blue-600 hover:bg-blue-50');
      default:
        return baseClasses + 'border-gray-300 text-gray-600';
    }
  }

  getInitials(student: StudentForPeriodAttendance): string {
    return (student.firstName[0] || '') + (student.lastName[0] || '');
  }

  formatDate(date: Date): string {
    return date.toISOString().split('T')[0];
  }

  formatDisplayDate(dateString: string): string {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', { weekday: 'short', month: 'short', day: 'numeric', year: 'numeric' });
  }

  formatDateTime(dateString?: string): string {
    if (!dateString) return '';
    return new Date(dateString).toLocaleString();
  }
}
