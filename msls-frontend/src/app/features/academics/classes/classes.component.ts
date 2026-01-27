import { Component, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';

/**
 * ClassesComponent - Classes management page.
 * Placeholder component for Story 1.8 layout integration.
 */
@Component({
  selector: 'msls-classes',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="classes">
      <h1 class="text-2xl font-bold text-secondary-900 mb-4">Classes</h1>
      <p class="text-secondary-600">Class management will be displayed here.</p>
    </div>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ClassesComponent {}
