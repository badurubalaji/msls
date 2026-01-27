import {
  Component,
  input,
  output,
  signal,
  computed,
  ContentChildren,
  QueryList,
  TemplateRef,
  Directive,
} from '@angular/core';
import { CommonModule } from '@angular/common';

/** Sort direction type */
export type SortDirection = 'asc' | 'desc' | null;

/** Sort event interface */
export interface SortEvent {
  column: string;
  direction: SortDirection;
}

/** Table column definition */
export interface TableColumn {
  key: string;
  label: string;
  sortable?: boolean;
  width?: string;
  align?: 'left' | 'center' | 'right';
}

/**
 * Directive for custom cell templates
 */
@Directive({
  selector: '[mslsTableCell]',
  standalone: true,
})
export class MslsTableCellDirective {
  columnKey = input.required<string>({ alias: 'mslsTableCell' });
  constructor(public templateRef: TemplateRef<unknown>) {}
}

/**
 * MslsTableComponent - A data table with sorting functionality using Tailwind CSS.
 */
@Component({
  selector: 'msls-table',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="relative w-full overflow-x-auto rounded-lg border border-slate-200 bg-white"
         [class.min-h-52]="loading()">
      <!-- Loading Overlay -->
      @if (loading()) {
        <div class="absolute inset-0 flex items-center justify-center bg-white/90 backdrop-blur-sm z-10">
          <div class="h-8 w-8 animate-spin rounded-full border-2 border-slate-200 border-t-slate-900"></div>
        </div>
      }

      <!-- Table -->
      <table class="w-full text-sm">
        <!-- Header - Clean minimal style -->
        <thead class="bg-slate-50 border-b border-slate-200">
          <tr>
            @for (column of columns(); track column.key) {
              <th
                class="px-4 py-3 text-left text-xs font-semibold text-slate-500 uppercase tracking-wider"
                [class.cursor-pointer]="isColumnSortable(column)"
                [class.select-none]="isColumnSortable(column)"
                [class.hover:text-slate-900]="isColumnSortable(column)"
                [style.width]="column.width"
                [style.text-align]="column.align || 'left'"
                (click)="onHeaderClick(column)"
              >
                <div class="inline-flex items-center gap-1.5">
                  <span>{{ column.label }}</span>
                  @if (isColumnSortable(column)) {
                    <span class="inline-flex shrink-0 text-xs">
                      @switch (getSortIcon(column.key)) {
                        @case ('asc') {
                          <i class="fa-solid fa-sort-up text-slate-900"></i>
                        }
                        @case ('desc') {
                          <i class="fa-solid fa-sort-down text-slate-900"></i>
                        }
                        @default {
                          <i class="fa-solid fa-sort text-slate-300"></i>
                        }
                      }
                    </span>
                  }
                </div>
              </th>
            }
          </tr>
        </thead>

        <!-- Body -->
        <tbody class="divide-y divide-slate-100">
          @for (row of sortedData(); track $index) {
            <tr
              class="transition-colors duration-150"
              [class.hover:bg-slate-50]="hoverable()"
              [class.cursor-pointer]="hoverable()"
              (click)="onRowClick(row)"
            >
              @for (column of columns(); track column.key) {
                <td
                  class="px-4 py-3.5 text-slate-700"
                  [style.text-align]="column.align || 'left'"
                >
                  @if (getCellTemplate(column.key); as template) {
                    <ng-container
                      *ngTemplateOutlet="template; context: { $implicit: row, value: getCellValue(row, column.key) }"
                    ></ng-container>
                  } @else {
                    {{ getCellValue(row, column.key) }}
                  }
                </td>
              }
            </tr>
          } @empty {
            <tr>
              <td [attr.colspan]="columns().length" class="px-4 py-12 text-center">
                <div class="flex flex-col items-center gap-3">
                  <i class="fa-regular fa-folder-open text-4xl text-slate-300"></i>
                  <span class="text-slate-500 text-sm">{{ emptyMessage() }}</span>
                </div>
              </td>
            </tr>
          }
        </tbody>
      </table>
    </div>
  `,
})
export class MslsTableComponent {
  /** Table column definitions */
  columns = input<TableColumn[]>([]);

  /** Table data rows */
  data = input<Record<string, unknown>[]>([]);

  /** Enable sorting on all sortable columns */
  sortable = input<boolean>(true);

  /** Current sort column */
  sortColumn = signal<string | null>(null);

  /** Current sort direction */
  sortDirection = signal<SortDirection>(null);

  /** Loading state */
  loading = input<boolean>(false);

  /** Empty state message */
  emptyMessage = input<string>('No data available');

  /** Striped rows */
  striped = input<boolean>(false);

  /** Hoverable rows */
  hoverable = input<boolean>(true);

  /** Compact mode */
  compact = input<boolean>(false);

  /** Emitted when sort changes */
  sortChange = output<SortEvent>();

  /** Emitted when a row is clicked */
  rowClick = output<Record<string, unknown>>();

  /** Cell templates from content */
  @ContentChildren(MslsTableCellDirective)
  cellTemplates!: QueryList<MslsTableCellDirective>;

  /** Computed sorted data */
  sortedData = computed(() => {
    const rawData = this.data();
    const column = this.sortColumn();
    const direction = this.sortDirection();

    if (!column || !direction) {
      return rawData;
    }

    return [...rawData].sort((a, b) => {
      const aValue = a[column];
      const bValue = b[column];

      if (aValue == null && bValue == null) return 0;
      if (aValue == null) return direction === 'asc' ? 1 : -1;
      if (bValue == null) return direction === 'asc' ? -1 : 1;

      if (typeof aValue === 'string' && typeof bValue === 'string') {
        const comparison = aValue.localeCompare(bValue);
        return direction === 'asc' ? comparison : -comparison;
      }

      if (typeof aValue === 'number' && typeof bValue === 'number') {
        return direction === 'asc' ? aValue - bValue : bValue - aValue;
      }

      const comparison = String(aValue).localeCompare(String(bValue));
      return direction === 'asc' ? comparison : -comparison;
    });
  });

  /** Handle column header click for sorting */
  onHeaderClick(column: TableColumn): void {
    if (!this.sortable() || column.sortable === false) {
      return;
    }

    const currentColumn = this.sortColumn();
    const currentDirection = this.sortDirection();

    if (currentColumn !== column.key) {
      this.sortColumn.set(column.key);
      this.sortDirection.set('asc');
    } else {
      if (currentDirection === 'asc') {
        this.sortDirection.set('desc');
      } else if (currentDirection === 'desc') {
        this.sortColumn.set(null);
        this.sortDirection.set(null);
      } else {
        this.sortDirection.set('asc');
      }
    }

    this.sortChange.emit({
      column: this.sortColumn() ?? '',
      direction: this.sortDirection(),
    });
  }

  /** Handle row click */
  onRowClick(row: Record<string, unknown>): void {
    this.rowClick.emit(row);
  }

  /** Get cell template for a column */
  getCellTemplate(columnKey: string): TemplateRef<unknown> | null {
    const directive = this.cellTemplates?.find(t => t.columnKey() === columnKey);
    return directive?.templateRef ?? null;
  }

  /** Get cell value */
  getCellValue(row: Record<string, unknown>, columnKey: string): unknown {
    return row[columnKey];
  }

  /** Check if column is currently sorted */
  isSorted(columnKey: string): boolean {
    return this.sortColumn() === columnKey && this.sortDirection() !== null;
  }

  /** Get sort icon for column */
  getSortIcon(columnKey: string): 'asc' | 'desc' | 'none' {
    if (this.sortColumn() !== columnKey) {
      return 'none';
    }
    return this.sortDirection() ?? 'none';
  }

  /** Check if column is sortable */
  isColumnSortable(column: TableColumn): boolean {
    return this.sortable() && column.sortable !== false;
  }
}
