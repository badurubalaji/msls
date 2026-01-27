/**
 * Attendance Management Page (HR View)
 * Allows HR to view all staff attendance and manage regularizations
 */

import { Component, inject, signal, computed, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { AttendanceService } from '../../services/attendance.service';
import { BranchService } from '../../../admin/branches/branch.service';
import { DepartmentService } from '../../../admin/departments/department.service';
import { StaffService } from '../../../staff/services/staff.service';
import {
  Attendance,
  Regularization,
  AttendanceStatus,
  RegularizationStatus,
  ATTENDANCE_STATUS_LABELS,
  REGULARIZATION_STATUS_LABELS,
  AttendanceFilter,
  RegularizationFilter,
} from '../../models/attendance.model';

interface Branch {
  id: string;
  name: string;
}

interface Department {
  id: string;
  name: string;
}

interface Staff {
  id: string;
  fullName: string;
}

@Component({
  selector: 'app-attendance-management',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './attendance-management.component.html',
  styleUrl: './attendance-management.component.scss',
})
export class AttendanceManagementComponent implements OnInit {
  private attendanceService = inject(AttendanceService);
  private branchService = inject(BranchService);
  private departmentService = inject(DepartmentService);
  private staffService = inject(StaffService);

  loading = signal(false);
  loadingRegularizations = signal(false);
  error = signal<string | null>(null);

  activeTab = signal<'attendance' | 'regularization'>('attendance');

  // Data
  attendanceList = signal<Attendance[]>([]);
  regularizationList = signal<Regularization[]>([]);
  branches = signal<Branch[]>([]);
  departments = signal<Department[]>([]);
  staffList = signal<Staff[]>([]);

  // Pagination
  hasMoreAttendance = signal(false);
  hasMoreRegularizations = signal(false);
  attendanceCursor = signal<string | null>(null);
  regularizationCursor = signal<string | null>(null);

  // Filters
  attendanceFilter = signal<Partial<AttendanceFilter>>({});
  regularizationFilter = signal<Partial<RegularizationFilter>>({});

  // Stats
  stats = computed(() => {
    const list = this.attendanceList();
    return {
      present: list.filter(a => a.status === 'present').length,
      absent: list.filter(a => a.status === 'absent').length,
      late: list.filter(a => a.lateMinutes && a.lateMinutes > 0).length,
      total: list.length,
    };
  });

  pendingRegularizations = computed(() => {
    return this.regularizationList().filter(r => r.status === 'pending').length;
  });

  // Modals
  showMarkAttendanceModal = signal(false);
  showRejectModal = signal(false);
  selectedRegularization = signal<Regularization | null>(null);
  rejectReason = '';

  markAttendanceForm = {
    staffId: '',
    date: '',
    status: '' as AttendanceStatus | '',
    notes: '',
  };

  todayDate = new Date().toISOString().split('T')[0];

  statusOptions = [
    { value: 'present', label: 'Present' },
    { value: 'absent', label: 'Absent' },
    { value: 'half_day', label: 'Half Day' },
    { value: 'on_leave', label: 'On Leave' },
  ];

  ngOnInit(): void {
    this.loadInitialData();
  }

  loadInitialData(): void {
    this.loadBranches();
    this.loadDepartments();
    this.loadStaff();
    this.loadData();
  }

  loadData(): void {
    this.loadAttendance();
    this.loadRegularizations();
  }

  loadBranches(): void {
    this.branchService.getBranches().subscribe({
      next: (branches) => {
        this.branches.set(branches || []);
      },
      error: (err) => console.error('Failed to load branches', err),
    });
  }

  loadDepartments(): void {
    this.departmentService.getDepartments().subscribe({
      next: (departments) => {
        this.departments.set(departments || []);
      },
      error: (err) => console.error('Failed to load departments', err),
    });
  }

  loadStaff(): void {
    this.staffService.loadStaff({ limit: 100 }).subscribe({
      next: (response) => {
        this.staffList.set(response.staff || []);
      },
      error: (err) => console.error('Failed to load staff', err),
    });
  }

  loadAttendance(): void {
    this.loading.set(true);

    const filter: AttendanceFilter = {
      ...this.attendanceFilter(),
      dateFrom: this.attendanceFilter().dateFrom || this.todayDate,
      dateTo: this.attendanceFilter().dateTo || this.todayDate,
      limit: 20,
    };

    this.attendanceService.getAllAttendance(filter).subscribe({
      next: (response) => {
        this.attendanceList.set(response.attendance);
        this.hasMoreAttendance.set(response.hasMore);
        this.attendanceCursor.set(response.nextCursor || null);
        this.loading.set(false);
      },
      error: (err) => {
        console.error('Failed to load attendance', err);
        this.error.set(err.message || 'Failed to load attendance');
        this.loading.set(false);
      },
    });
  }

  loadMoreAttendance(): void {
    if (!this.attendanceCursor()) return;

    this.loading.set(true);

    const filter: AttendanceFilter = {
      ...this.attendanceFilter(),
      cursor: this.attendanceCursor()!,
      limit: 20,
    };

    this.attendanceService.getAllAttendance(filter).subscribe({
      next: (response) => {
        this.attendanceList.update((list) => [...list, ...response.attendance]);
        this.hasMoreAttendance.set(response.hasMore);
        this.attendanceCursor.set(response.nextCursor || null);
        this.loading.set(false);
      },
      error: (err) => {
        console.error('Failed to load more attendance', err);
        this.loading.set(false);
      },
    });
  }

  loadRegularizations(): void {
    this.loadingRegularizations.set(true);

    const filter: RegularizationFilter = {
      ...this.regularizationFilter(),
      limit: 20,
    };

    this.attendanceService.getRegularizations(filter).subscribe({
      next: (response) => {
        this.regularizationList.set(response.regularizations);
        this.hasMoreRegularizations.set(response.hasMore);
        this.regularizationCursor.set(response.nextCursor || null);
        this.loadingRegularizations.set(false);
      },
      error: (err) => {
        console.error('Failed to load regularizations', err);
        this.loadingRegularizations.set(false);
      },
    });
  }

  loadMoreRegularizations(): void {
    if (!this.regularizationCursor()) return;

    this.loadingRegularizations.set(true);

    const filter: RegularizationFilter = {
      ...this.regularizationFilter(),
      cursor: this.regularizationCursor()!,
      limit: 20,
    };

    this.attendanceService.getRegularizations(filter).subscribe({
      next: (response) => {
        this.regularizationList.update((list) => [...list, ...response.regularizations]);
        this.hasMoreRegularizations.set(response.hasMore);
        this.regularizationCursor.set(response.nextCursor || null);
        this.loadingRegularizations.set(false);
      },
      error: (err) => {
        console.error('Failed to load more regularizations', err);
        this.loadingRegularizations.set(false);
      },
    });
  }

  setActiveTab(tab: 'attendance' | 'regularization'): void {
    this.activeTab.set(tab);
  }

  // Filter handlers
  onDateChange(event: Event): void {
    const input = event.target as HTMLInputElement;
    this.attendanceFilter.update((f) => ({ ...f, dateFrom: input.value, dateTo: input.value }));
    this.loadAttendance();
  }

  onBranchChange(event: Event): void {
    const select = event.target as HTMLSelectElement;
    this.attendanceFilter.update((f) => ({ ...f, branchId: select.value || undefined }));
    this.loadAttendance();
  }

  onDepartmentChange(event: Event): void {
    const select = event.target as HTMLSelectElement;
    this.attendanceFilter.update((f) => ({ ...f, departmentId: select.value || undefined }));
    this.loadAttendance();
  }

  onStatusChange(event: Event): void {
    const select = event.target as HTMLSelectElement;
    this.attendanceFilter.update((f) => ({
      ...f,
      status: (select.value || undefined) as AttendanceStatus | undefined,
    }));
    this.loadAttendance();
  }

  onRegularizationStatusChange(event: Event): void {
    const select = event.target as HTMLSelectElement;
    this.regularizationFilter.update((f) => ({
      ...f,
      status: (select.value || undefined) as RegularizationStatus | undefined,
    }));
    this.loadRegularizations();
  }

  onRegularizationDateFromChange(event: Event): void {
    const input = event.target as HTMLInputElement;
    this.regularizationFilter.update((f) => ({ ...f, dateFrom: input.value || undefined }));
    this.loadRegularizations();
  }

  onRegularizationDateToChange(event: Event): void {
    const input = event.target as HTMLInputElement;
    this.regularizationFilter.update((f) => ({ ...f, dateTo: input.value || undefined }));
    this.loadRegularizations();
  }

  // Modal handlers
  openMarkAttendanceModal(): void {
    this.markAttendanceForm = {
      staffId: '',
      date: this.todayDate,
      status: '',
      notes: '',
    };
    this.showMarkAttendanceModal.set(true);
  }

  closeMarkAttendanceModal(): void {
    this.showMarkAttendanceModal.set(false);
  }

  submitMarkAttendance(): void {
    if (!this.markAttendanceForm.staffId || !this.markAttendanceForm.date || !this.markAttendanceForm.status) {
      return;
    }

    this.loading.set(true);

    this.attendanceService
      .markAttendance({
        staffId: this.markAttendanceForm.staffId,
        attendanceDate: this.markAttendanceForm.date,
        status: this.markAttendanceForm.status as AttendanceStatus,
        remarks: this.markAttendanceForm.notes || undefined,
      })
      .subscribe({
        next: () => {
          this.closeMarkAttendanceModal();
          this.loadAttendance();
        },
        error: (err) => {
          console.error('Failed to mark attendance', err);
          this.error.set(err.message || 'Failed to mark attendance');
          this.loading.set(false);
        },
      });
  }

  approveRegularization(id: string): void {
    this.loadingRegularizations.set(true);

    this.attendanceService.approveRegularization(id).subscribe({
      next: () => {
        this.loadRegularizations();
      },
      error: (err) => {
        console.error('Failed to approve regularization', err);
        this.loadingRegularizations.set(false);
      },
    });
  }

  openRejectModal(reg: Regularization): void {
    this.selectedRegularization.set(reg);
    this.rejectReason = '';
    this.showRejectModal.set(true);
  }

  closeRejectModal(): void {
    this.showRejectModal.set(false);
    this.selectedRegularization.set(null);
  }

  submitRejectRegularization(): void {
    const reg = this.selectedRegularization();
    if (!reg || !this.rejectReason) return;

    this.loadingRegularizations.set(true);

    this.attendanceService.rejectRegularization(reg.id, { rejectionReason: this.rejectReason }).subscribe({
      next: () => {
        this.closeRejectModal();
        this.loadRegularizations();
      },
      error: (err) => {
        console.error('Failed to reject regularization', err);
        this.loadingRegularizations.set(false);
      },
    });
  }

  // Helpers
  getInitials(name: string): string {
    return name
      .split(' ')
      .map((n) => n[0])
      .join('')
      .toUpperCase()
      .slice(0, 2);
  }

  getStatusLabel(status: AttendanceStatus): string {
    return ATTENDANCE_STATUS_LABELS[status] || status;
  }

  getRegularizationStatusLabel(status: RegularizationStatus): string {
    return REGULARIZATION_STATUS_LABELS[status] || status;
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

  formatDisplayDate(dateStr: string): string {
    const date = new Date(dateStr);
    return date.toLocaleDateString('en-US', {
      weekday: 'long',
      month: 'long',
      day: 'numeric',
      year: 'numeric',
    });
  }
}
