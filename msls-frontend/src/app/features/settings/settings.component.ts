import { Component, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';

/**
 * SettingsComponent - Application settings page.
 * Placeholder component for Story 1.8 layout integration.
 */
@Component({
  selector: 'msls-settings',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="settings">
      <h1 class="text-2xl font-bold text-secondary-900 mb-4">Settings</h1>
      <p class="text-secondary-600">Application settings will be displayed here.</p>
    </div>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class SettingsComponent {}
