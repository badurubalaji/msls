/**
 * Bulk Actions Component
 *
 * Floating action bar for bulk operations on selected students.
 */

import { Component, input, output, signal } from '@angular/core';
import { CommonModule } from '@angular/common';

export type BulkActionType = 'sms' | 'email' | 'status' | 'export';

@Component({
  selector: 'msls-bulk-actions',
  standalone: true,
  imports: [CommonModule],
  template: `
    @if (selectedCount() > 0) {
      <div class="bulk-actions-bar">
        <div class="selection-info">
          <span class="count">{{ selectedCount() }}</span>
          <span class="label">students selected</span>
        </div>

        <div class="actions">
          <button
            type="button"
            class="action-btn action-btn--disabled"
            (click)="onAction('sms')"
            title="Coming Soon"
            disabled
          >
            <i class="fa-solid fa-message"></i>
            <span>Send SMS</span>
          </button>

          <button
            type="button"
            class="action-btn action-btn--disabled"
            (click)="onAction('email')"
            title="Coming Soon"
            disabled
          >
            <i class="fa-solid fa-envelope"></i>
            <span>Send Email</span>
          </button>

          <button
            type="button"
            class="action-btn"
            (click)="onAction('status')"
          >
            <i class="fa-solid fa-edit"></i>
            <span>Update Status</span>
          </button>

          <button
            type="button"
            class="action-btn action-btn--primary"
            (click)="onAction('export')"
          >
            <i class="fa-solid fa-download"></i>
            <span>Export</span>
          </button>
        </div>

        <button
          type="button"
          class="close-btn"
          (click)="onClear()"
          title="Clear selection"
        >
          <i class="fa-solid fa-times"></i>
        </button>
      </div>
    }
  `,
  styles: [`
    .bulk-actions-bar {
      position: fixed;
      bottom: 1.5rem;
      left: 50%;
      transform: translateX(-50%);
      display: flex;
      align-items: center;
      gap: 1.5rem;
      padding: 0.75rem 1rem 0.75rem 1.25rem;
      background: var(--color-bg-dark, #1f2937);
      color: white;
      border-radius: 0.75rem;
      box-shadow: 0 10px 25px rgba(0, 0, 0, 0.3);
      z-index: 50;
      animation: slideUp 0.3s ease-out;
    }

    @keyframes slideUp {
      from {
        opacity: 0;
        transform: translateX(-50%) translateY(20px);
      }
      to {
        opacity: 1;
        transform: translateX(-50%) translateY(0);
      }
    }

    .selection-info {
      display: flex;
      align-items: baseline;
      gap: 0.375rem;

      .count {
        font-size: 1.25rem;
        font-weight: 600;
      }

      .label {
        font-size: 0.875rem;
        color: rgba(255, 255, 255, 0.7);
      }
    }

    .actions {
      display: flex;
      gap: 0.5rem;
    }

    .action-btn {
      display: flex;
      align-items: center;
      gap: 0.375rem;
      padding: 0.5rem 0.75rem;
      background: rgba(255, 255, 255, 0.1);
      color: white;
      border: 1px solid rgba(255, 255, 255, 0.2);
      border-radius: 0.375rem;
      font-size: 0.875rem;
      cursor: pointer;
      transition: background-color 0.2s, border-color 0.2s;

      &:hover:not(:disabled) {
        background: rgba(255, 255, 255, 0.2);
        border-color: rgba(255, 255, 255, 0.3);
      }

      &--primary {
        background: var(--color-primary);
        border-color: var(--color-primary);

        &:hover {
          background: var(--color-primary-dark);
        }
      }

      &--disabled {
        opacity: 0.5;
        cursor: not-allowed;
      }

      i {
        font-size: 0.875rem;
      }
    }

    .close-btn {
      padding: 0.5rem;
      background: none;
      border: none;
      color: rgba(255, 255, 255, 0.6);
      cursor: pointer;
      border-radius: 0.25rem;
      transition: color 0.2s;

      &:hover {
        color: white;
      }

      i {
        font-size: 1rem;
      }
    }
  `],
})
export class BulkActionsComponent {
  /** Number of selected students */
  selectedCount = input.required<number>();

  /** IDs of selected students */
  selectedIds = input.required<string[]>();

  /** Emits when an action is clicked */
  action = output<{ type: BulkActionType; ids: string[] }>();

  /** Emits when selection is cleared */
  cleared = output<void>();

  onAction(type: BulkActionType): void {
    this.action.emit({ type, ids: this.selectedIds() });
  }

  onClear(): void {
    this.cleared.emit();
  }
}
