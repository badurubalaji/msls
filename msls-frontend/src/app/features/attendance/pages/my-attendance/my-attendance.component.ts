/**
 * My Attendance Page
 * Staff's personal attendance dashboard
 */

import { Component, inject, signal, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { AttendanceService } from '../../services/attendance.service';
import { StaffService } from '../../../staff/services/staff.service';
import { AuthService } from '../../../../core/services/auth.service';
import {
  Attendance,
  TodayAttendance,
  AttendanceSummary,
  AttendanceStatus,
  ATTENDANCE_STATUS_LABELS,
} from '../../models/attendance.model';

@Component({
  selector: 'app-my-attendance',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './my-attendance.component.html',
  styleUrl: './my-attendance.component.scss',
})
export class MyAttendanceComponent implements OnInit, OnDestroy {
  private attendanceService = inject(AttendanceService);
  private staffService = inject(StaffService);
  private authService = inject(AuthService);

  private clockInterval?: ReturnType<typeof setInterval>;

  loading = signal(true);
  error = signal<string | null>(null);

  staffId = signal<string | null>(null);
  todayAttendance = signal<TodayAttendance | null>(null);
  summary = signal<AttendanceSummary | null>(null);
  attendanceList = signal<Attendance[]>([]);
  hasMore = signal(false);
  nextCursor = signal<string | null>(null);

  currentDate = '';
  currentTime = '';

  selectedYear = new Date().getFullYear();
  selectedMonth = new Date().getMonth() + 1;

  dateFrom = signal('');
  dateTo = signal('');

  showRegularizationModal = signal(false);
  regularizationDate = '';
  regularizationReason = '';

  ngOnInit(): void {
    this.updateDateTime();
    this.clockInterval = setInterval(() => this.updateDateTime(), 1000);
    this.loadStaffProfile();
  }

  ngOnDestroy(): void {
    if (this.clockInterval) {
      clearInterval(this.clockInterval);
    }
  }

  private updateDateTime(): void {
    const now = new Date();
    this.currentDate = now.toLocaleDateString('en-US', {
      weekday: 'long',
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
    this.currentTime = now.toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: true,
    });
  }

  loadStaffProfile(): void {
    this.loading.set(true);

    const user = this.authService.currentUser();
    if (user) {
      this.staffService.loadStaff({ limit: 1 }).subscribe({
        next: (response) => {
          if (response.staff.length > 0) {
            this.staffId.set(response.staff[0].id);
            this.loadData();
          } else {
            this.loading.set(false);
            this.error.set('No staff profile found');
          }
        },
        error: (err) => {
          this.loading.set(false);
          this.error.set(err.message || 'Failed to load staff profile');
        },
      });
    } else {
      this.loading.set(false);
      this.error.set('User not authenticated');
    }
  }

  loadData(): void {
    this.loadTodayAttendance();
    this.loadSummary();
    this.loadAttendanceHistory();
  }

  loadTodayAttendance(): void {
    if (!this.staffId()) return;

    this.attendanceService.getTodayAttendance(this.staffId()!).subscribe({
      next: (attendance) => {
        this.todayAttendance.set(attendance);
        this.loading.set(false);
      },
      error: (err) => {
        console.error('Failed to load today attendance', err);
        this.loading.set(false);
      },
    });
  }

  loadSummary(): void {
    if (!this.staffId()) return;

    this.attendanceService
      .getMonthlySummary(this.staffId()!, this.selectedYear, this.selectedMonth)
      .subscribe({
        next: (summary) => {
          this.summary.set(summary);
        },
        error: (err) => {
          console.error('Failed to load summary', err);
        },
      });
  }

  loadAttendanceHistory(): void {
    if (!this.staffId()) return;

    this.loading.set(true);

    const filter: { dateFrom?: string; dateTo?: string; limit: number; cursor?: string } = {
      limit: 20,
    };

    if (this.dateFrom()) filter.dateFrom = this.dateFrom();
    if (this.dateTo()) filter.dateTo = this.dateTo();

    this.attendanceService.getMyAttendance(this.staffId()!, filter).subscribe({
      next: (response) => {
        this.attendanceList.set(response.attendance);
        this.hasMore.set(response.hasMore);
        this.nextCursor.set(response.nextCursor || null);
        this.loading.set(false);
      },
      error: (err) => {
        console.error('Failed to load attendance history', err);
        this.loading.set(false);
      },
    });
  }

  loadMore(): void {
    if (!this.staffId() || !this.nextCursor()) return;

    this.loading.set(true);

    this.attendanceService
      .getMyAttendance(this.staffId()!, { cursor: this.nextCursor()!, limit: 20 })
      .subscribe({
        next: (response) => {
          this.attendanceList.update((list) => [...list, ...response.attendance]);
          this.hasMore.set(response.hasMore);
          this.nextCursor.set(response.nextCursor || null);
          this.loading.set(false);
        },
        error: (err) => {
          console.error('Failed to load more attendance', err);
          this.loading.set(false);
        },
      });
  }

  checkIn(): void {
    if (!this.staffId()) return;

    this.loading.set(true);
    this.attendanceService.checkIn({ staffId: this.staffId()! }).subscribe({
      next: () => {
        this.loadTodayAttendance();
        this.loadSummary();
        this.loadAttendanceHistory();
      },
      error: (err) => {
        console.error('Check-in failed', err);
        this.error.set(err.message || 'Check-in failed');
        this.loading.set(false);
      },
    });
  }

  checkOut(): void {
    if (!this.staffId()) return;

    this.loading.set(true);
    this.attendanceService.checkOut({ staffId: this.staffId()! }).subscribe({
      next: () => {
        this.loadTodayAttendance();
        this.loadSummary();
        this.loadAttendanceHistory();
      },
      error: (err) => {
        console.error('Check-out failed', err);
        this.error.set(err.message || 'Check-out failed');
        this.loading.set(false);
      },
    });
  }

  previousMonth(): void {
    if (this.selectedMonth === 1) {
      this.selectedMonth = 12;
      this.selectedYear--;
    } else {
      this.selectedMonth--;
    }
    this.loadSummary();
  }

  nextMonth(): void {
    const now = new Date();
    if (this.selectedYear === now.getFullYear() && this.selectedMonth === now.getMonth() + 1) {
      return;
    }

    if (this.selectedMonth === 12) {
      this.selectedMonth = 1;
      this.selectedYear++;
    } else {
      this.selectedMonth++;
    }
    this.loadSummary();
  }

  getMonthLabel(): string {
    const date = new Date(this.selectedYear, this.selectedMonth - 1);
    return date.toLocaleDateString('en-US', { month: 'long', year: 'numeric' });
  }

  getWorkingHours(): string {
    const today = this.todayAttendance();
    if (!today?.attendance?.checkInTime) return '0h 0m';

    const checkIn = new Date(today.attendance.checkInTime);
    const checkOut = today.attendance.checkOutTime ? new Date(today.attendance.checkOutTime) : new Date();

    const diffMs = checkOut.getTime() - checkIn.getTime();
    const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
    const diffMinutes = Math.floor((diffMs % (1000 * 60 * 60)) / (1000 * 60));

    return `${diffHours}h ${diffMinutes}m`;
  }

  calculateWorkingHours(checkIn?: string, checkOut?: string): string {
    if (!checkIn) return '0 hrs';
    if (!checkOut) return '-';

    const checkInDate = new Date(checkIn);
    const checkOutDate = new Date(checkOut);

    const diffMs = checkOutDate.getTime() - checkInDate.getTime();
    const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
    const diffMinutes = Math.floor((diffMs % (1000 * 60 * 60)) / (1000 * 60));

    if (diffHours === 0) return `${diffMinutes} min`;
    return `${diffHours}h ${diffMinutes}m`;
  }

  formatTime(isoString?: string): string {
    if (!isoString) return '';
    const date = new Date(isoString);
    return date.toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit',
      hour12: true,
    });
  }

  onDateFromChange(event: Event): void {
    const input = event.target as HTMLInputElement;
    this.dateFrom.set(input.value);
    this.loadAttendanceHistory();
  }

  onDateToChange(event: Event): void {
    const input = event.target as HTMLInputElement;
    this.dateTo.set(input.value);
    this.loadAttendanceHistory();
  }

  clearFilters(): void {
    this.dateFrom.set('');
    this.dateTo.set('');
    this.loadAttendanceHistory();
  }

  openRegularizationModal(): void {
    this.regularizationDate = '';
    this.regularizationReason = '';
    this.showRegularizationModal.set(true);
  }

  closeRegularizationModal(): void {
    this.showRegularizationModal.set(false);
  }

  submitRegularization(): void {
    if (!this.staffId() || !this.regularizationDate || !this.regularizationReason) return;

    this.loading.set(true);

    this.attendanceService
      .submitRegularization({
        staffId: this.staffId()!,
        requestDate: this.regularizationDate,
        requestedStatus: 'present' as AttendanceStatus,
        reason: this.regularizationReason,
      })
      .subscribe({
        next: () => {
          this.closeRegularizationModal();
          this.loading.set(false);
        },
        error: (err) => {
          console.error('Failed to submit regularization', err);
          this.error.set(err.message || 'Failed to submit regularization');
          this.loading.set(false);
        },
      });
  }

  formatDate(dateStr: string): string {
    const date = new Date(dateStr);
    return date.toLocaleDateString('en-US', {
      weekday: 'short',
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    });
  }

  getStatusLabel(status: AttendanceStatus): string {
    return ATTENDANCE_STATUS_LABELS[status] || status;
  }
}
