/**
 * Export Dialog Component
 *
 * Modal dialog for configuring and executing student exports.
 */

import { Component, input, output, signal, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { DEFAULT_EXPORT_COLUMNS, ExportColumn, ExportRequest } from '../../models/student.model';

@Component({
  selector: 'msls-export-dialog',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    @if (isOpen()) {
      <div class="modal-overlay" (click)="onCancel()">
        <div class="modal" (click)="$event.stopPropagation()">
          <!-- Header -->
          <div class="modal__header">
            <div class="header-content">
              <div class="header-icon">
                <i class="fa-solid fa-file-export"></i>
              </div>
              <div class="header-text">
                <h3>Export Students</h3>
                <p>Download student data in your preferred format</p>
              </div>
            </div>
            <button type="button" class="modal__close" (click)="onCancel()">
              <i class="fa-solid fa-xmark"></i>
            </button>
          </div>

          <div class="modal__body">
            <!-- Export Info -->
            <div class="export-info">
              <i class="fa-solid fa-users"></i>
              <span>Exporting <strong>{{ studentCount() }}</strong> {{ studentCount() === 1 ? 'student' : 'students' }}</span>
            </div>

            <!-- Format Selection -->
            <div class="form-group">
              <label class="form-label">Export Format</label>
              <div class="format-options">
                <label class="format-option" [class.selected]="selectedFormat() === 'xlsx'">
                  <input
                    type="radio"
                    name="format"
                    value="xlsx"
                    [checked]="selectedFormat() === 'xlsx'"
                    (change)="selectedFormat.set('xlsx')"
                  />
                  <div class="format-content">
                    <div class="format-icon format-icon--excel">
                      <i class="fa-solid fa-file-excel"></i>
                    </div>
                    <div class="format-text">
                      <span class="format-name">Excel</span>
                      <span class="format-ext">.xlsx</span>
                    </div>
                    <i class="fa-solid fa-check format-check"></i>
                  </div>
                </label>
                <label class="format-option" [class.selected]="selectedFormat() === 'csv'">
                  <input
                    type="radio"
                    name="format"
                    value="csv"
                    [checked]="selectedFormat() === 'csv'"
                    (change)="selectedFormat.set('csv')"
                  />
                  <div class="format-content">
                    <div class="format-icon format-icon--csv">
                      <i class="fa-solid fa-file-csv"></i>
                    </div>
                    <div class="format-text">
                      <span class="format-name">CSV</span>
                      <span class="format-ext">.csv</span>
                    </div>
                    <i class="fa-solid fa-check format-check"></i>
                  </div>
                </label>
              </div>
            </div>

            <!-- Column Selection -->
            <div class="form-group">
              <div class="form-label-row">
                <label class="form-label">Columns to Export</label>
                <button type="button" class="select-all-btn" (click)="toggleAllColumns()">
                  {{ allColumnsSelected() ? 'Deselect All' : 'Select All' }}
                </button>
              </div>
              <div class="columns-grid">
                @for (column of columns(); track column.key) {
                  <label class="column-option" [class.selected]="column.selected">
                    <input
                      type="checkbox"
                      [checked]="column.selected"
                      (change)="toggleColumn(column.key)"
                    />
                    <span class="column-label">{{ column.label }}</span>
                    <i class="fa-solid fa-check column-check"></i>
                  </label>
                }
              </div>
            </div>
          </div>

          <div class="modal__footer">
            <button type="button" class="btn btn--secondary" (click)="onCancel()">
              <i class="fa-solid fa-times"></i>
              Cancel
            </button>
            <button
              type="button"
              class="btn btn--primary"
              (click)="onExport()"
              [disabled]="!canExport()"
            >
              <i class="fa-solid fa-download"></i>
              Export {{ studentCount() }} {{ studentCount() === 1 ? 'Student' : 'Students' }}
            </button>
          </div>
        </div>
      </div>
    }
  `,
  styles: [`
    .modal-overlay {
      position: fixed;
      inset: 0;
      background: rgba(15, 23, 42, 0.6);
      backdrop-filter: blur(4px);
      display: flex;
      align-items: center;
      justify-content: center;
      z-index: 1000;
      padding: 1rem;
      animation: fadeIn 0.2s ease-out;
    }

    @keyframes fadeIn {
      from { opacity: 0; }
      to { opacity: 1; }
    }

    .modal {
      background: #ffffff;
      border-radius: 1rem;
      width: 100%;
      max-width: 520px;
      max-height: 90vh;
      display: flex;
      flex-direction: column;
      box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
      animation: slideUp 0.3s ease-out;
      overflow: hidden;
    }

    @keyframes slideUp {
      from {
        opacity: 0;
        transform: translateY(20px) scale(0.98);
      }
      to {
        opacity: 1;
        transform: translateY(0) scale(1);
      }
    }

    .modal__header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      padding: 1.25rem 1.5rem;
      background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
      border-bottom: 1px solid #e2e8f0;
    }

    .header-content {
      display: flex;
      align-items: center;
      gap: 1rem;
    }

    .header-icon {
      width: 2.75rem;
      height: 2.75rem;
      border-radius: 0.75rem;
      display: flex;
      align-items: center;
      justify-content: center;
      background: linear-gradient(135deg, #10b981 0%, #059669 100%);
      color: white;
      font-size: 1.125rem;
      box-shadow: 0 4px 12px rgba(16, 185, 129, 0.3);
    }

    .header-text h3 {
      margin: 0;
      font-size: 1.125rem;
      font-weight: 600;
      color: #0f172a;
    }

    .header-text p {
      margin: 0.25rem 0 0 0;
      font-size: 0.8125rem;
      color: #64748b;
    }

    .modal__close {
      width: 2rem;
      height: 2rem;
      display: flex;
      align-items: center;
      justify-content: center;
      background: #ffffff;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      color: #64748b;
      cursor: pointer;
      transition: all 0.15s;

      &:hover {
        background: #f1f5f9;
        color: #0f172a;
        border-color: #cbd5e1;
      }
    }

    .modal__body {
      padding: 1.5rem;
      overflow-y: auto;
      background: #ffffff;
    }

    .export-info {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      margin-bottom: 1.5rem;
      padding: 0.875rem 1rem;
      background: linear-gradient(135deg, #eff6ff 0%, #dbeafe 100%);
      border: 1px solid #bfdbfe;
      border-radius: 0.625rem;
      font-size: 0.875rem;
      color: #1e40af;

      i {
        font-size: 1rem;
      }
    }

    .form-group {
      margin-bottom: 1.5rem;
    }

    .form-label {
      display: block;
      margin-bottom: 0.625rem;
      font-size: 0.8125rem;
      font-weight: 600;
      color: #334155;
      text-transform: uppercase;
      letter-spacing: 0.025em;
    }

    .form-label-row {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 0.625rem;
    }

    .select-all-btn {
      padding: 0.375rem 0.75rem;
      font-size: 0.75rem;
      font-weight: 500;
      background: #f1f5f9;
      border: 1px solid #e2e8f0;
      border-radius: 0.375rem;
      color: #6366f1;
      cursor: pointer;
      transition: all 0.15s;

      &:hover {
        background: #eef2ff;
        border-color: #c7d2fe;
      }
    }

    .format-options {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: 0.75rem;
    }

    .format-option {
      cursor: pointer;

      input {
        display: none;
      }

      &.selected .format-content {
        border-color: #6366f1;
        background: linear-gradient(135deg, #eef2ff 0%, #e0e7ff 100%);
        box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
      }

      &.selected .format-check {
        opacity: 1;
        transform: scale(1);
      }

      &.selected .format-icon {
        transform: scale(1.1);
      }
    }

    .format-content {
      position: relative;
      display: flex;
      align-items: center;
      gap: 0.875rem;
      padding: 1rem;
      border: 2px solid #e2e8f0;
      border-radius: 0.625rem;
      background: #ffffff;
      transition: all 0.2s;

      &:hover {
        border-color: #cbd5e1;
        background: #f8fafc;
      }
    }

    .format-icon {
      width: 2.5rem;
      height: 2.5rem;
      display: flex;
      align-items: center;
      justify-content: center;
      border-radius: 0.5rem;
      font-size: 1.25rem;
      transition: transform 0.2s;
    }

    .format-icon--excel {
      background: linear-gradient(135deg, #dcfce7 0%, #bbf7d0 100%);
      color: #16a34a;
    }

    .format-icon--csv {
      background: linear-gradient(135deg, #fef3c7 0%, #fde68a 100%);
      color: #d97706;
    }

    .format-text {
      display: flex;
      flex-direction: column;
      gap: 0.125rem;
    }

    .format-name {
      font-size: 0.9375rem;
      font-weight: 600;
      color: #0f172a;
    }

    .format-ext {
      font-size: 0.75rem;
      color: #64748b;
    }

    .format-check {
      position: absolute;
      top: 0.625rem;
      right: 0.625rem;
      width: 1.25rem;
      height: 1.25rem;
      display: flex;
      align-items: center;
      justify-content: center;
      background: #6366f1;
      color: white;
      border-radius: 50%;
      font-size: 0.625rem;
      opacity: 0;
      transform: scale(0.5);
      transition: all 0.2s;
    }

    .columns-grid {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: 0.5rem;
      padding: 0.75rem;
      background: #f8fafc;
      border: 1px solid #e2e8f0;
      border-radius: 0.625rem;
      max-height: 200px;
      overflow-y: auto;
    }

    .column-option {
      position: relative;
      display: flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.5rem 0.625rem;
      border-radius: 0.375rem;
      cursor: pointer;
      transition: all 0.15s;
      background: #ffffff;
      border: 1px solid transparent;

      &:hover {
        background: #f1f5f9;
      }

      &.selected {
        background: #eef2ff;
        border-color: #c7d2fe;
      }

      &.selected .column-check {
        opacity: 1;
        transform: scale(1);
      }

      input[type="checkbox"] {
        display: none;
      }
    }

    .column-label {
      font-size: 0.8125rem;
      color: #334155;
      flex: 1;
    }

    .column-check {
      width: 1rem;
      height: 1rem;
      display: flex;
      align-items: center;
      justify-content: center;
      background: #6366f1;
      color: white;
      border-radius: 0.25rem;
      font-size: 0.5rem;
      opacity: 0;
      transform: scale(0.5);
      transition: all 0.15s;
    }

    .modal__footer {
      display: flex;
      justify-content: flex-end;
      gap: 0.75rem;
      padding: 1rem 1.5rem;
      background: #f8fafc;
      border-top: 1px solid #e2e8f0;
    }

    .btn {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.625rem 1.25rem;
      border-radius: 0.5rem;
      font-size: 0.875rem;
      font-weight: 500;
      cursor: pointer;
      transition: all 0.15s;
    }

    .btn--secondary {
      background: #ffffff;
      color: #64748b;
      border: 1px solid #e2e8f0;

      &:hover {
        background: #f8fafc;
        border-color: #cbd5e1;
        color: #475569;
      }
    }

    .btn--primary {
      background: linear-gradient(135deg, #10b981 0%, #059669 100%);
      color: white;
      border: none;
      box-shadow: 0 4px 12px rgba(16, 185, 129, 0.3);

      &:hover:not(:disabled) {
        transform: translateY(-1px);
        box-shadow: 0 6px 16px rgba(16, 185, 129, 0.4);
      }

      &:active:not(:disabled) {
        transform: translateY(0);
      }

      &:disabled {
        opacity: 0.5;
        cursor: not-allowed;
        box-shadow: none;
      }
    }

    @media (max-width: 480px) {
      .modal {
        max-height: 95vh;
        border-radius: 0.75rem;
      }

      .format-options {
        grid-template-columns: 1fr;
      }

      .columns-grid {
        grid-template-columns: 1fr;
      }

      .modal__footer {
        flex-direction: column-reverse;
      }

      .btn {
        width: 100%;
        justify-content: center;
      }
    }
  `],
})
export class ExportDialogComponent {
  /** Whether the dialog is open */
  isOpen = input<boolean>(false);

  /** IDs of students to export */
  studentIds = input<string[]>([]);

  /** Number of students to export */
  studentCount = computed(() => this.studentIds().length);

  /** Emits the export request when confirmed */
  exported = output<ExportRequest>();

  /** Emits when dialog is cancelled */
  cancelled = output<void>();

  /** Selected export format */
  selectedFormat = signal<'xlsx' | 'csv'>('xlsx');

  /** Column configuration */
  columns = signal<ExportColumn[]>([...DEFAULT_EXPORT_COLUMNS]);

  /** Whether all columns are selected */
  allColumnsSelected = computed(() => this.columns().every((c) => c.selected));

  /** Whether export can proceed */
  canExport = computed(() => this.columns().some((c) => c.selected));

  toggleColumn(key: string): void {
    this.columns.update((cols) =>
      cols.map((c) => (c.key === key ? { ...c, selected: !c.selected } : c))
    );
  }

  toggleAllColumns(): void {
    const selectAll = !this.allColumnsSelected();
    this.columns.update((cols) => cols.map((c) => ({ ...c, selected: selectAll })));
  }

  onExport(): void {
    const selectedColumns = this.columns()
      .filter((c) => c.selected)
      .map((c) => c.key);

    this.exported.emit({
      studentIds: this.studentIds(),
      format: this.selectedFormat(),
      columns: selectedColumns,
    });
  }

  onCancel(): void {
    this.cancelled.emit();
  }
}
