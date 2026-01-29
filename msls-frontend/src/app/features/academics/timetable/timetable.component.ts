import { Component, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink } from '@angular/router';

/**
 * TimetableComponent - Timetable configuration landing page.
 * Links to Shifts, Day Patterns, and Period Slots configuration.
 */
@Component({
  selector: 'msls-timetable',
  standalone: true,
  imports: [CommonModule, RouterLink],
  template: `
    <div class="page">
      <!-- Page Header -->
      <div class="page-header">
        <div class="header-content">
          <div class="header-icon">
            <i class="fa-solid fa-calendar-days"></i>
          </div>
          <div class="header-text">
            <h1>Timetable Configuration</h1>
            <p>Configure school timings, shifts, and period structures</p>
          </div>
        </div>
      </div>

      <!-- Configuration Cards -->
      <div class="config-grid">
        <!-- Shifts Card -->
        <a routerLink="shifts" class="config-card">
          <div class="card-icon shifts">
            <i class="fa-solid fa-clock-rotate-left"></i>
          </div>
          <div class="card-content">
            <h3>Shifts</h3>
            <p>Configure morning, afternoon, and custom shifts for your school</p>
          </div>
          <div class="card-arrow">
            <i class="fa-solid fa-chevron-right"></i>
          </div>
        </a>

        <!-- Day Patterns Card -->
        <a routerLink="day-patterns" class="config-card">
          <div class="card-icon patterns">
            <i class="fa-solid fa-layer-group"></i>
          </div>
          <div class="card-content">
            <h3>Day Patterns</h3>
            <p>Define different day schedules like regular, half-day, or assembly days</p>
          </div>
          <div class="card-arrow">
            <i class="fa-solid fa-chevron-right"></i>
          </div>
        </a>

        <!-- Period Slots Card -->
        <a routerLink="period-slots" class="config-card">
          <div class="card-icon slots">
            <i class="fa-solid fa-table-cells"></i>
          </div>
          <div class="card-content">
            <h3>Period Slots</h3>
            <p>Configure period timings, breaks, lunch, and activity slots</p>
          </div>
          <div class="card-arrow">
            <i class="fa-solid fa-chevron-right"></i>
          </div>
        </a>
      </div>

      <!-- Quick Info -->
      <div class="info-section">
        <h2>How It Works</h2>
        <div class="info-steps">
          <div class="info-step">
            <div class="step-number">1</div>
            <div class="step-content">
              <h4>Define Shifts</h4>
              <p>Create shifts for your school (e.g., Morning 8AM-2PM, Afternoon 12PM-6PM)</p>
            </div>
          </div>
          <div class="info-step">
            <div class="step-number">2</div>
            <div class="step-content">
              <h4>Create Day Patterns</h4>
              <p>Define different day structures (Regular Day - 8 periods, Saturday - 5 periods)</p>
            </div>
          </div>
          <div class="info-step">
            <div class="step-number">3</div>
            <div class="step-content">
              <h4>Configure Period Slots</h4>
              <p>Set up individual periods with times, breaks, and activities</p>
            </div>
          </div>
          <div class="info-step">
            <div class="step-number">4</div>
            <div class="step-content">
              <h4>Assign to Days</h4>
              <p>Map day patterns to weekdays for each branch</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  `,
  styles: [`
    .page { padding: 1.5rem; max-width: 1200px; margin: 0 auto; }
    .page-header { margin-bottom: 2rem; }
    .header-content { display: flex; align-items: center; gap: 1rem; }
    .header-icon { width: 3.5rem; height: 3.5rem; border-radius: 1rem; background: linear-gradient(135deg, #4f46e5, #7c3aed); color: white; display: flex; align-items: center; justify-content: center; font-size: 1.5rem; }
    .header-text h1 { margin: 0; font-size: 1.75rem; font-weight: 600; color: #1e293b; }
    .header-text p { margin: 0.25rem 0 0; color: #64748b; font-size: 0.9375rem; }

    .config-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 1.25rem; margin-bottom: 2.5rem; }
    .config-card { display: flex; align-items: center; gap: 1rem; padding: 1.25rem; background: white; border: 1px solid #e2e8f0; border-radius: 1rem; text-decoration: none; transition: all 0.2s; }
    .config-card:hover { border-color: #c7d2fe; box-shadow: 0 4px 12px rgba(79, 70, 229, 0.08); transform: translateY(-2px); }
    .card-icon { width: 3rem; height: 3rem; border-radius: 0.75rem; display: flex; align-items: center; justify-content: center; font-size: 1.25rem; flex-shrink: 0; }
    .card-icon.shifts { background: #dbeafe; color: #2563eb; }
    .card-icon.patterns { background: #e0e7ff; color: #4f46e5; }
    .card-icon.slots { background: #dcfce7; color: #16a34a; }
    .card-content { flex: 1; }
    .card-content h3 { margin: 0 0 0.25rem; font-size: 1rem; font-weight: 600; color: #1e293b; }
    .card-content p { margin: 0; font-size: 0.8125rem; color: #64748b; line-height: 1.4; }
    .card-arrow { color: #94a3b8; transition: transform 0.2s; }
    .config-card:hover .card-arrow { transform: translateX(4px); color: #4f46e5; }

    .info-section { background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 1rem; padding: 1.5rem; }
    .info-section h2 { margin: 0 0 1.25rem; font-size: 1.125rem; font-weight: 600; color: #1e293b; }
    .info-steps { display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap: 1.25rem; }
    .info-step { display: flex; gap: 0.75rem; }
    .step-number { width: 2rem; height: 2rem; border-radius: 50%; background: #4f46e5; color: white; display: flex; align-items: center; justify-content: center; font-weight: 600; font-size: 0.875rem; flex-shrink: 0; }
    .step-content h4 { margin: 0 0 0.25rem; font-size: 0.875rem; font-weight: 600; color: #1e293b; }
    .step-content p { margin: 0; font-size: 0.75rem; color: #64748b; line-height: 1.4; }
  `],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class TimetableComponent {}
