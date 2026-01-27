import { Component, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';

/**
 * TimetableComponent - Timetable management page.
 * Placeholder component for Story 1.8 layout integration.
 */
@Component({
  selector: 'msls-timetable',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="timetable">
      <h1 class="text-2xl font-bold text-secondary-900 mb-4">Timetable</h1>
      <p class="text-secondary-600">Timetable management will be displayed here.</p>
    </div>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class TimetableComponent {}
