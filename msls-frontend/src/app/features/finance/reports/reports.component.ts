import { Component, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';

/**
 * ReportsComponent - Financial reports page.
 * Placeholder component for Story 1.8 layout integration.
 */
@Component({
  selector: 'msls-reports',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="reports">
      <h1 class="text-2xl font-bold text-secondary-900 mb-4">Financial Reports</h1>
      <p class="text-secondary-600">Financial reports will be displayed here.</p>
    </div>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ReportsComponent {}
