/**
 * Student Search Component
 *
 * Provides search input with debouncing and filter toggle.
 */

import {
  Component,
  OnInit,
  OnDestroy,
  inject,
  input,
  output,
  signal,
  computed,
  DestroyRef,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Subject } from 'rxjs';
import { debounceTime, distinctUntilChanged } from 'rxjs/operators';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

import { MslsBadgeComponent } from '../../../../shared/components/badge/badge.component';
import { StudentListFilter } from '../../models/student.model';

@Component({
  selector: 'msls-student-search',
  standalone: true,
  imports: [CommonModule, FormsModule, MslsBadgeComponent],
  template: `
    <div class="search-container">
      <div class="search-input-wrapper">
        <i class="fa-solid fa-search search-icon"></i>
        <input
          type="text"
          class="search-input"
          [value]="searchQuery()"
          (input)="onSearchInput($event)"
          placeholder="Search by name, admission number, phone..."
        />
        @if (searchQuery()) {
          <button type="button" class="clear-btn" (click)="clearSearch()">
            <i class="fa-solid fa-times"></i>
          </button>
        }
      </div>

      <button type="button" class="filter-btn" (click)="toggleFilters()">
        <i class="fa-solid fa-filter"></i>
        <span>Filters</span>
        @if (activeFilterCount() > 0) {
          <msls-badge variant="primary" size="sm">{{ activeFilterCount() }}</msls-badge>
        }
      </button>
    </div>
  `,
  styles: [`
    .search-container {
      display: flex;
      gap: 0.75rem;
      align-items: center;
    }

    .search-input-wrapper {
      position: relative;
      flex: 1;
      max-width: 400px;
    }

    .search-icon {
      position: absolute;
      left: 0.75rem;
      top: 50%;
      transform: translateY(-50%);
      color: var(--color-text-muted);
      font-size: 0.875rem;
    }

    .search-input {
      width: 100%;
      padding: 0.5rem 2.25rem 0.5rem 2.25rem;
      border: 1px solid var(--color-border);
      border-radius: 0.5rem;
      font-size: 0.875rem;
      background: var(--color-bg-primary);
      color: var(--color-text-primary);
      transition: border-color 0.2s, box-shadow 0.2s;

      &:focus {
        outline: none;
        border-color: var(--color-primary);
        box-shadow: 0 0 0 3px var(--color-primary-light);
      }

      &::placeholder {
        color: var(--color-text-muted);
      }
    }

    .clear-btn {
      position: absolute;
      right: 0.5rem;
      top: 50%;
      transform: translateY(-50%);
      padding: 0.25rem;
      background: none;
      border: none;
      color: var(--color-text-muted);
      cursor: pointer;
      border-radius: 0.25rem;
      transition: color 0.2s;

      &:hover {
        color: var(--color-text-primary);
      }
    }

    .filter-btn {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.5rem 1rem;
      background: var(--color-bg-secondary);
      border: 1px solid var(--color-border);
      border-radius: 0.5rem;
      font-size: 0.875rem;
      color: var(--color-text-primary);
      cursor: pointer;
      transition: background-color 0.2s, border-color 0.2s;

      &:hover {
        background: var(--color-bg-hover);
        border-color: var(--color-border-dark);
      }

      i {
        font-size: 0.875rem;
      }
    }
  `],
})
export class StudentSearchComponent implements OnInit, OnDestroy {
  private destroyRef = inject(DestroyRef);

  /** Current filter values */
  filters = input<StudentListFilter>({});

  /** Emits when search query changes (debounced) */
  searchChanged = output<string>();

  /** Emits when filters button is clicked */
  filtersToggled = output<void>();

  /** Current search query */
  searchQuery = signal('');

  /** Whether filters panel is shown */
  showFilters = signal(false);

  /** Count of active filters */
  activeFilterCount = computed(() => {
    const f = this.filters();
    let count = 0;
    if (f.classId) count++;
    if (f.sectionId) count++;
    if (f.status) count++;
    if (f.gender) count++;
    if (f.admissionFrom) count++;
    if (f.admissionTo) count++;
    return count;
  });

  private searchSubject = new Subject<string>();

  ngOnInit(): void {
    // Initialize search from filters
    const f = this.filters();
    if (f.search) {
      this.searchQuery.set(f.search);
    }

    // Set up debounced search
    this.searchSubject
      .pipe(debounceTime(300), distinctUntilChanged(), takeUntilDestroyed(this.destroyRef))
      .subscribe((query: string) => {
        this.searchChanged.emit(query);
      });
  }

  ngOnDestroy(): void {
    this.searchSubject.complete();
  }

  onSearchInput(event: Event): void {
    const value = (event.target as HTMLInputElement).value;
    this.searchQuery.set(value);
    this.searchSubject.next(value);
  }

  clearSearch(): void {
    this.searchQuery.set('');
    this.searchSubject.next('');
  }

  toggleFilters(): void {
    this.showFilters.update((v) => !v);
    this.filtersToggled.emit();
  }
}
