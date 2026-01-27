/**
 * Staff List Page Component
 *
 * Displays a paginated, filterable list of staff with quick actions.
 */

import {
  Component,
  ChangeDetectionStrategy,
  inject,
  OnInit,
  signal,
  computed,
  DestroyRef,
} from '@angular/core';
import { CommonModule, TitleCasePipe } from '@angular/common';
import { Router, ActivatedRoute } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

import {
  MslsBadgeComponent,
  MslsAvatarComponent,
} from '../../../../shared/components';
import { StaffService } from '../../services/staff.service';
import {
  Staff,
  StaffStatus,
  StaffType,
  getStatusBadgeVariant,
  getStatusLabel,
  getStaffTypeLabel,
  StaffListFilter,
  calculateTenure,
} from '../../models/staff.model';

/** Simple option interface for native select */
interface FilterOption {
  value: string;
  label: string;
}

@Component({
  selector: 'msls-staff-list',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MslsBadgeComponent,
    MslsAvatarComponent,
  ],
  templateUrl: './staff-list.component.html',
  styleUrl: './staff-list.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class StaffListComponent implements OnInit {
  private staffService = inject(StaffService);
  private route = inject(ActivatedRoute);
  private destroyRef = inject(DestroyRef);
  readonly router = inject(Router);

  // =========================================================================
  // Reactive State from Service
  // =========================================================================

  readonly staffList = this.staffService.staffList;
  readonly loading = this.staffService.loading;
  readonly error = this.staffService.error;
  readonly totalCount = this.staffService.totalCount;
  readonly hasMore = this.staffService.hasMore;
  readonly isEmpty = this.staffService.isEmpty;

  // =========================================================================
  // Local State
  // =========================================================================

  readonly searchTerm = signal<string>('');
  readonly selectedStatus = signal<StaffStatus | ''>('');
  readonly selectedStaffType = signal<StaffType | ''>('');

  /** Current filter state */
  readonly currentFilters = signal<StaffListFilter>({});

  /** Whether the filters panel is visible */
  readonly showFilters = signal(false);

  // =========================================================================
  // Filter Options
  // =========================================================================

  readonly statusOptions: FilterOption[] = [
    { value: '', label: 'All Status' },
    { value: 'active', label: 'Active' },
    { value: 'inactive', label: 'Inactive' },
    { value: 'on_leave', label: 'On Leave' },
    { value: 'terminated', label: 'Terminated' },
  ];

  readonly staffTypeOptions: FilterOption[] = [
    { value: '', label: 'All Types' },
    { value: 'teaching', label: 'Teaching' },
    { value: 'non_teaching', label: 'Non-Teaching' },
  ];

  // =========================================================================
  // Lifecycle
  // =========================================================================

  ngOnInit(): void {
    // Read filters from URL query params
    this.route.queryParams
      .pipe(takeUntilDestroyed(this.destroyRef))
      .subscribe((params) => {
        const filter: StaffListFilter = {};
        if (params['search']) {
          filter.search = params['search'];
          this.searchTerm.set(params['search']);
        }
        if (params['status']) {
          filter.status = params['status'] as StaffStatus;
          this.selectedStatus.set(params['status'] as StaffStatus);
        }
        if (params['staffType']) {
          filter.staffType = params['staffType'] as StaffType;
          this.selectedStaffType.set(params['staffType'] as StaffType);
        }
        if (params['departmentId']) filter.departmentId = params['departmentId'];
        if (params['designationId']) filter.designationId = params['designationId'];
        if (params['joinDateFrom']) filter.joinDateFrom = params['joinDateFrom'];
        if (params['joinDateTo']) filter.joinDateTo = params['joinDateTo'];
        if (params['sortBy']) filter.sortBy = params['sortBy'] as StaffListFilter['sortBy'];
        if (params['sortOrder']) filter.sortOrder = params['sortOrder'] as 'asc' | 'desc';

        this.currentFilters.set(filter);
        this.loadStaffWithFilter(filter);
      });
  }

  // =========================================================================
  // Data Loading
  // =========================================================================

  loadStaff(): void {
    this.loadStaffWithFilter(this.currentFilters());
  }

  private loadStaffWithFilter(filter: StaffListFilter): void {
    const mergedFilter: StaffListFilter = {
      limit: 20,
      ...filter,
    };

    if (this.searchTerm()) {
      mergedFilter.search = this.searchTerm();
    }
    if (this.selectedStatus()) {
      mergedFilter.status = this.selectedStatus() as StaffStatus;
    }
    if (this.selectedStaffType()) {
      mergedFilter.staffType = this.selectedStaffType() as StaffType;
    }

    this.staffService.loadStaff(mergedFilter).subscribe();
  }

  loadMore(): void {
    const observable = this.staffService.loadMore();
    if (observable) {
      observable.subscribe();
    }
  }

  refresh(): void {
    this.staffService.refresh().subscribe();
  }

  /** Update URL with current filters */
  private updateUrlFilters(filter: StaffListFilter): void {
    const queryParams: Record<string, string | null> = {
      search: filter.search || null,
      status: filter.status || null,
      staffType: filter.staffType || null,
      departmentId: filter.departmentId || null,
      designationId: filter.designationId || null,
      joinDateFrom: filter.joinDateFrom || null,
      joinDateTo: filter.joinDateTo || null,
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

  onSearchChange(event: Event): void {
    const input = event.target as HTMLInputElement;
    const value = input.value;
    this.searchTerm.set(value);
    const filter = { ...this.currentFilters(), search: value || undefined };
    this.currentFilters.set(filter);
    this.updateUrlFilters(filter);
    this.loadStaffWithFilter(filter);
  }

  onStatusChange(event: Event): void {
    const select = event.target as HTMLSelectElement;
    const value = select.value as StaffStatus | '';
    this.selectedStatus.set(value);
    const filter = { ...this.currentFilters(), status: (value || undefined) as StaffStatus };
    this.currentFilters.set(filter);
    this.updateUrlFilters(filter);
    this.loadStaffWithFilter(filter);
  }

  onStaffTypeChange(event: Event): void {
    const select = event.target as HTMLSelectElement;
    const value = select.value as StaffType | '';
    this.selectedStaffType.set(value);
    const filter = { ...this.currentFilters(), staffType: (value || undefined) as StaffType };
    this.currentFilters.set(filter);
    this.updateUrlFilters(filter);
    this.loadStaffWithFilter(filter);
  }

  onFiltersToggle(): void {
    this.showFilters.update((v) => !v);
  }

  clearFilters(): void {
    this.searchTerm.set('');
    this.selectedStatus.set('');
    this.selectedStaffType.set('');
    this.currentFilters.set({});
    this.updateUrlFilters({});
    this.loadStaffWithFilter({});
  }

  // =========================================================================
  // Actions
  // =========================================================================

  onRowClick(staff: Staff): void {
    this.router.navigate(['/staff', staff.id]);
  }

  onAddStaff(): void {
    this.router.navigate(['/staff', 'new']);
  }

  onEditStaff(event: Event, staff: Staff): void {
    event.stopPropagation();
    this.router.navigate(['/staff', staff.id, 'edit']);
  }

  onViewStaff(event: Event, staff: Staff): void {
    event.stopPropagation();
    this.router.navigate(['/staff', staff.id]);
  }

  // =========================================================================
  // Helpers
  // =========================================================================

  getStatusVariant(status: StaffStatus): 'success' | 'warning' | 'danger' | 'neutral' {
    return getStatusBadgeVariant(status);
  }

  getStatusLabel(status: StaffStatus): string {
    return getStatusLabel(status);
  }

  getStaffTypeLabel(type: StaffType): string {
    return getStaffTypeLabel(type);
  }

  getTenure(joinDate: string): string {
    return calculateTenure(joinDate);
  }

  trackByStaffId(index: number, staff: Staff): string {
    return staff.id;
  }

  // =========================================================================
  // Stats Helpers
  // =========================================================================

  getActiveCount(): number {
    return this.staffList().filter(s => s.status === 'active').length;
  }

  getTeachingCount(): number {
    return this.staffList().filter(s => s.staffType === 'teaching').length;
  }

  getNonTeachingCount(): number {
    return this.staffList().filter(s => s.staffType === 'non_teaching').length;
  }

  getOnLeaveCount(): number {
    return this.staffList().filter(s => s.status === 'on_leave').length;
  }
}
