import { Component, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';

/**
 * SubjectsComponent - Subjects management page.
 * Placeholder component for Story 1.8 layout integration.
 */
@Component({
  selector: 'msls-subjects',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="subjects">
      <h1 class="text-2xl font-bold text-secondary-900 mb-4">Subjects</h1>
      <p class="text-secondary-600">Subject management will be displayed here.</p>
    </div>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class SubjectsComponent {}
