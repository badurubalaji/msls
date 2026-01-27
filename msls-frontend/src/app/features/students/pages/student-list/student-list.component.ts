/**
 * Student List Page Component
 *
 * Displays a paginated, filterable list of students with quick actions.
 * Supports bulk operations and export functionality.
 */

import {
  Component,
  ChangeDetectionStrategy,
  inject,
  OnInit,
  OnDestroy,
  signal,
  computed,
  DestroyRef,
} from '@angular/core';
import { CommonModule, DatePipe, TitleCasePipe } from '@angular/common';
import { Router, ActivatedRoute } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

import {
  MslsBadgeComponent,
  MslsAvatarComponent,
} from '../../../../shared/components';
import { StudentService } from '../../services/student.service';
import {
  Student,
  StudentStatus,
  getStatusBadgeVariant,
  getStatusLabel,
  StudentListFilter,
  ExportRequest,
  BulkStatusUpdateRequest,
} from '../../models/student.model';
import {
  StudentSearchComponent,
  StudentFiltersComponent,
  BulkActionsComponent,
  ExportDialogComponent,
  BulkActionType,
} from '../../components';

/** Simple option interface for native select */
interface StatusOption {
  value: string;
  label: string;
}

@Component({
  selector: 'msls-student-list',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MslsBadgeComponent,
    MslsAvatarComponent,
    StudentSearchComponent,
    StudentFiltersComponent,
    BulkActionsComponent,
    ExportDialogComponent,
  ],
  templateUrl: './student-list.component.html',
  styleUrl: './student-list.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class StudentListComponent implements OnInit, OnDestroy {
  private studentService = inject(StudentService);
  private route = inject(ActivatedRoute);
  private destroyRef = inject(DestroyRef);
  readonly router = inject(Router);

  // =========================================================================
  // Reactive State from Service
  // =========================================================================

  readonly students = this.studentService.students;
  readonly loading = this.studentService.loading;
  readonly error = this.studentService.error;
  readonly totalCount = this.studentService.totalCount;
  readonly hasMore = this.studentService.hasMore;
  readonly isEmpty = this.studentService.isEmpty;

  // =========================================================================
  // Local State
  // =========================================================================

  readonly searchTerm = signal<string>('');
  readonly selectedStatus = signal<StudentStatus | ''>('');
  readonly selectedBranch = signal<string>('');

  /** Current filter state for advanced filters */
  readonly currentFilters = signal<StudentListFilter>({});

  /** Whether the filters panel is visible */
  readonly showFilters = signal(false);

  /** Selected student IDs for bulk operations */
  readonly selectedIds = signal<Set<string>>(new Set());

  /** Whether the export dialog is open */
  readonly showExportDialog = signal(false);

  /** Whether the status update dialog is open */
  readonly showStatusDialog = signal(false);

  /** Status to apply in bulk update */
  readonly bulkStatus = signal<StudentStatus>('active');

  /** Whether all visible students are selected */
  readonly allSelected = computed(() => {
    const ids = this.selectedIds();
    const students = this.students();
    return students.length > 0 && students.every((s) => ids.has(s.id));
  });

  /** Array of selected IDs for components */
  readonly selectedIdsArray = computed(() => Array.from(this.selectedIds()));

  // =========================================================================
  // Filter Options
  // =========================================================================

  readonly statusOptions: StatusOption[] = [
    { value: '', label: 'All Status' },
    { value: 'active', label: 'Active' },
    { value: 'inactive', label: 'Inactive' },
    { value: 'transferred', label: 'Transferred' },
    { value: 'graduated', label: 'Graduated' },
  ];

  // =========================================================================
  // Lifecycle
  // =========================================================================

  ngOnInit(): void {
    // Read filters from URL query params
    this.route.queryParams
      .pipe(takeUntilDestroyed(this.destroyRef))
      .subscribe((params) => {
        const filter: StudentListFilter = {};
        if (params['search']) {
          filter.search = params['search'];
          this.searchTerm.set(params['search']);
        }
        if (params['status']) {
          filter.status = params['status'] as StudentStatus;
          this.selectedStatus.set(params['status'] as StudentStatus);
        }
        if (params['gender']) filter.gender = params['gender'];
        if (params['classId']) filter.classId = params['classId'];
        if (params['sectionId']) filter.sectionId = params['sectionId'];
        if (params['admissionFrom']) filter.admissionFrom = params['admissionFrom'];
        if (params['admissionTo']) filter.admissionTo = params['admissionTo'];
        if (params['sortBy']) filter.sortBy = params['sortBy'];
        if (params['sortOrder']) filter.sortOrder = params['sortOrder'] as 'asc' | 'desc';

        this.currentFilters.set(filter);
        this.loadStudentsWithFilter(filter);
      });
  }

  ngOnDestroy(): void {
    // Clear selection on destroy
    this.selectedIds.set(new Set());
  }

  // =========================================================================
  // Data Loading
  // =========================================================================

  loadStudents(): void {
    this.loadStudentsWithFilter(this.currentFilters());
  }

  private loadStudentsWithFilter(filter: StudentListFilter): void {
    const mergedFilter: StudentListFilter = {
      limit: 20,
      ...filter,
    };

    if (this.searchTerm()) {
      mergedFilter.search = this.searchTerm();
    }
    if (this.selectedStatus()) {
      mergedFilter.status = this.selectedStatus() as StudentStatus;
    }
    if (this.selectedBranch()) {
      mergedFilter.branchId = this.selectedBranch();
    }

    // Clear selection when loading new data
    this.selectedIds.set(new Set());

    this.studentService.loadStudents(mergedFilter).subscribe();
  }

  loadMore(): void {
    const observable = this.studentService.loadMore();
    if (observable) {
      observable.subscribe();
    }
  }

  refresh(): void {
    this.selectedIds.set(new Set());
    this.studentService.refresh().subscribe();
  }

  /** Update URL with current filters */
  private updateUrlFilters(filter: StudentListFilter): void {
    const queryParams: Record<string, string | null> = {
      search: filter.search || null,
      status: filter.status || null,
      gender: filter.gender || null,
      classId: filter.classId || null,
      sectionId: filter.sectionId || null,
      admissionFrom: filter.admissionFrom || null,
      admissionTo: filter.admissionTo || null,
      sortBy: filter.sortBy || null,
      sortOrder: filter.sortOrder || null,
    };

    this.router.navigate([], {
      relativeTo: this.route,
      queryParams,
      queryParamsHandling: 'merge',
      replaceUrl: true,
    });
  }

  // =========================================================================
  // Filtering
  // =========================================================================

  onSearchChange(value: string): void {
    this.searchTerm.set(value);
    const filter = { ...this.currentFilters(), search: value || undefined };
    this.currentFilters.set(filter);
    this.updateUrlFilters(filter);
    this.loadStudentsWithFilter(filter);
  }

  onStatusChange(value: string): void {
    this.selectedStatus.set(value as StudentStatus | '');
    const filter = { ...this.currentFilters(), status: (value || undefined) as StudentStatus };
    this.currentFilters.set(filter);
    this.updateUrlFilters(filter);
    this.loadStudentsWithFilter(filter);
  }

  onBranchChange(value: string): void {
    this.selectedBranch.set(value);
    const filter = { ...this.currentFilters(), branchId: value || undefined };
    this.currentFilters.set(filter);
    this.loadStudentsWithFilter(filter);
  }

  onFiltersToggle(): void {
    this.showFilters.update((v) => !v);
  }

  onAdvancedFiltersChange(filter: StudentListFilter): void {
    const merged = { ...this.currentFilters(), ...filter };
    this.currentFilters.set(merged);
    this.updateUrlFilters(merged);
    this.loadStudentsWithFilter(merged);
    this.showFilters.set(false);
  }

  clearFilters(): void {
    this.searchTerm.set('');
    this.selectedStatus.set('');
    this.selectedBranch.set('');
    this.currentFilters.set({});
    this.updateUrlFilters({});
    this.loadStudentsWithFilter({});
  }

  // =========================================================================
  // Actions
  // =========================================================================

  onRowClick(row: Record<string, unknown>): void {
    const studentId = row['id'] as string;
    this.router.navigate(['/students', studentId]);
  }

  onAddStudent(): void {
    this.router.navigate(['/students', 'new']);
  }

  onEditStudent(event: Event, student: Student): void {
    event.stopPropagation();
    this.router.navigate(['/students', student.id, 'edit']);
  }

  onViewStudent(event: Event, student: Student): void {
    event.stopPropagation();
    this.router.navigate(['/students', student.id]);
  }

  // =========================================================================
  // Selection & Bulk Actions
  // =========================================================================

  toggleSelectAll(): void {
    const ids = this.selectedIds();
    const students = this.students();

    if (this.allSelected()) {
      // Deselect all
      this.selectedIds.set(new Set());
    } else {
      // Select all visible students
      this.selectedIds.set(new Set(students.map((s) => s.id)));
    }
  }

  toggleSelectStudent(event: Event, studentId: string): void {
    event.stopPropagation();
    this.selectedIds.update((ids) => {
      const newIds = new Set(ids);
      if (newIds.has(studentId)) {
        newIds.delete(studentId);
      } else {
        newIds.add(studentId);
      }
      return newIds;
    });
  }

  isSelected(studentId: string): boolean {
    return this.selectedIds().has(studentId);
  }

  clearSelection(): void {
    this.selectedIds.set(new Set());
  }

  onBulkAction(action: { type: BulkActionType; ids: string[] }): void {
    switch (action.type) {
      case 'export':
        this.showExportDialog.set(true);
        break;
      case 'status':
        this.showStatusDialog.set(true);
        break;
      case 'sms':
      case 'email':
        // Coming soon - these are disabled in the UI
        break;
    }
  }

  onExport(request: ExportRequest): void {
    this.studentService.exportStudents(request).subscribe({
      next: (operation) => {
        this.showExportDialog.set(false);
        this.clearSelection();
        // TODO: Show notification about export started
        // Poll for completion or show download link
        if (operation.resultUrl) {
          window.open(operation.resultUrl, '_blank');
        }
      },
      error: () => {
        // Error is handled by service
      },
    });
  }

  onExportCancel(): void {
    this.showExportDialog.set(false);
  }

  onBulkStatusUpdate(): void {
    const request: BulkStatusUpdateRequest = {
      studentIds: this.selectedIdsArray(),
      newStatus: this.bulkStatus(),
    };

    this.studentService.bulkUpdateStatus(request).subscribe({
      next: () => {
        this.showStatusDialog.set(false);
        this.clearSelection();
        // Service automatically refreshes the list
      },
      error: () => {
        // Error is handled by service
      },
    });
  }

  onStatusDialogCancel(): void {
    this.showStatusDialog.set(false);
  }

  // =========================================================================
  // Helpers
  // =========================================================================

  formatClassSection(student: Student): string {
    if (student.className && student.sectionName) {
      return `${student.className} - ${student.sectionName}`;
    }
    if (student.className) {
      return student.className;
    }
    return '—';
  }

  getStatusVariant(status: StudentStatus): 'success' | 'warning' | 'info' | 'neutral' {
    return getStatusBadgeVariant(status);
  }

  getStatusLabel(status: StudentStatus): string {
    return getStatusLabel(status);
  }

  trackByStudentId(index: number, row: Record<string, unknown>): string {
    return row['id'] as string;
  }

  // =========================================================================
  // Stats Helpers
  // =========================================================================

  getActiveCount(): number {
    return this.students().filter(s => s.status === 'active').length;
  }

  getNewThisMonth(): number {
    const now = new Date();
    const currentMonth = now.getMonth();
    const currentYear = now.getFullYear();
    return this.students().filter(s => {
      if (!s.admissionDate) return false;
      const admDate = new Date(s.admissionDate);
      return admDate.getMonth() === currentMonth && admDate.getFullYear() === currentYear;
    }).length;
  }

  getGraduatedCount(): number {
    return this.students().filter(s => s.status === 'graduated').length;
  }
}
