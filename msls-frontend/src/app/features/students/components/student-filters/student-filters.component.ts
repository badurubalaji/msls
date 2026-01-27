/**
 * Student Filters Component
 *
 * Provides advanced filter options for student search.
 */

import { Component, input, output, signal, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { StudentListFilter, StudentStatus, Gender } from '../../models/student.model';

@Component({
  selector: 'msls-student-filters',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <div class="filters-panel">
      <div class="filters-header">
        <h3>Filters</h3>
        @if (hasActiveFilters()) {
          <button type="button" class="clear-all-btn" (click)="clearAllFilters()">
            Clear All
          </button>
        }
      </div>

      <div class="filters-grid">
        <!-- Status Filter -->
        <div class="filter-group">
          <label>Status</label>
          <select
            [value]="currentFilters().status || ''"
            (change)="onStatusChange($event)"
            class="filter-select"
          >
            <option value="">All Statuses</option>
            <option value="active">Active</option>
            <option value="inactive">Inactive</option>
            <option value="transferred">Transferred</option>
            <option value="graduated">Graduated</option>
          </select>
        </div>

        <!-- Gender Filter -->
        <div class="filter-group">
          <label>Gender</label>
          <select
            [value]="currentFilters().gender || ''"
            (change)="onGenderChange($event)"
            class="filter-select"
          >
            <option value="">All Genders</option>
            <option value="male">Male</option>
            <option value="female">Female</option>
            <option value="other">Other</option>
          </select>
        </div>

        <!-- Admission Date From -->
        <div class="filter-group">
          <label>Admission From</label>
          <input
            type="date"
            [value]="currentFilters().admissionFrom || ''"
            (change)="onAdmissionFromChange($event)"
            class="filter-input"
          />
        </div>

        <!-- Admission Date To -->
        <div class="filter-group">
          <label>Admission To</label>
          <input
            type="date"
            [value]="currentFilters().admissionTo || ''"
            (change)="onAdmissionToChange($event)"
            class="filter-input"
          />
        </div>

        <!-- Sort By -->
        <div class="filter-group">
          <label>Sort By</label>
          <select
            [value]="currentFilters().sortBy || ''"
            (change)="onSortByChange($event)"
            class="filter-select"
          >
            <option value="">Default (Name)</option>
            <option value="name">Name</option>
            <option value="admission_number">Admission Number</option>
            <option value="created_at">Date Added</option>
          </select>
        </div>

        <!-- Sort Order -->
        <div class="filter-group">
          <label>Order</label>
          <select
            [value]="currentFilters().sortOrder || 'asc'"
            (change)="onSortOrderChange($event)"
            class="filter-select"
          >
            <option value="asc">Ascending</option>
            <option value="desc">Descending</option>
          </select>
        </div>
      </div>

      <div class="filters-footer">
        <button type="button" class="apply-btn" (click)="applyFilters()">
          Apply Filters
        </button>
      </div>
    </div>
  `,
  styles: [`
    .filters-panel {
      background: var(--color-bg-primary);
      border: 1px solid var(--color-border);
      border-radius: 0.5rem;
      padding: 1rem;
      margin-top: 0.75rem;
    }

    .filters-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 1rem;

      h3 {
        font-size: 0.875rem;
        font-weight: 600;
        color: var(--color-text-primary);
        margin: 0;
      }
    }

    .clear-all-btn {
      padding: 0.25rem 0.5rem;
      font-size: 0.75rem;
      background: none;
      border: none;
      color: var(--color-primary);
      cursor: pointer;
      border-radius: 0.25rem;

      &:hover {
        text-decoration: underline;
      }
    }

    .filters-grid {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
      gap: 1rem;
    }

    .filter-group {
      display: flex;
      flex-direction: column;
      gap: 0.375rem;

      label {
        font-size: 0.75rem;
        font-weight: 500;
        color: var(--color-text-muted);
      }
    }

    .filter-select,
    .filter-input {
      padding: 0.5rem 0.75rem;
      border: 1px solid var(--color-border);
      border-radius: 0.375rem;
      font-size: 0.875rem;
      background: var(--color-bg-primary);
      color: var(--color-text-primary);
      transition: border-color 0.2s;

      &:focus {
        outline: none;
        border-color: var(--color-primary);
      }
    }

    .filters-footer {
      display: flex;
      justify-content: flex-end;
      margin-top: 1rem;
      padding-top: 1rem;
      border-top: 1px solid var(--color-border);
    }

    .apply-btn {
      padding: 0.5rem 1rem;
      background: var(--color-primary);
      color: white;
      border: none;
      border-radius: 0.375rem;
      font-size: 0.875rem;
      font-weight: 500;
      cursor: pointer;
      transition: background-color 0.2s;

      &:hover {
        background: var(--color-primary-dark);
      }
    }
  `],
})
export class StudentFiltersComponent {
  /** Current filter values */
  filters = input<StudentListFilter>({});

  /** Emits when filters are changed */
  filtersChanged = output<StudentListFilter>();

  /** Local filter state for editing */
  currentFilters = signal<StudentListFilter>({});

  /** Whether any filters are active */
  hasActiveFilters = computed(() => {
    const f = this.currentFilters();
    return !!(f.status || f.gender || f.admissionFrom || f.admissionTo || f.sortBy);
  });

  constructor() {
    // Initialize local filters from input
    const initial = this.filters();
    this.currentFilters.set({ ...initial });
  }

  onStatusChange(event: Event): void {
    const value = (event.target as HTMLSelectElement).value as StudentStatus | '';
    this.currentFilters.update((f) => ({
      ...f,
      status: value || undefined,
    }));
  }

  onGenderChange(event: Event): void {
    const value = (event.target as HTMLSelectElement).value as Gender | '';
    this.currentFilters.update((f) => ({
      ...f,
      gender: value || undefined,
    }));
  }

  onAdmissionFromChange(event: Event): void {
    const value = (event.target as HTMLInputElement).value;
    this.currentFilters.update((f) => ({
      ...f,
      admissionFrom: value || undefined,
    }));
  }

  onAdmissionToChange(event: Event): void {
    const value = (event.target as HTMLInputElement).value;
    this.currentFilters.update((f) => ({
      ...f,
      admissionTo: value || undefined,
    }));
  }

  onSortByChange(event: Event): void {
    const value = (event.target as HTMLSelectElement).value as
      | 'name'
      | 'admission_number'
      | 'created_at'
      | '';
    this.currentFilters.update((f) => ({
      ...f,
      sortBy: value || undefined,
    }));
  }

  onSortOrderChange(event: Event): void {
    const value = (event.target as HTMLSelectElement).value as 'asc' | 'desc';
    this.currentFilters.update((f) => ({
      ...f,
      sortOrder: value,
    }));
  }

  clearAllFilters(): void {
    this.currentFilters.set({
      cursor: undefined,
      limit: this.filters().limit,
    });
    this.filtersChanged.emit(this.currentFilters());
  }

  applyFilters(): void {
    this.filtersChanged.emit(this.currentFilters());
  }
}
