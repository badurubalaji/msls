import { Component, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink } from '@angular/router';

/**
 * PayrollComponent - Payroll module landing page.
 * Links to Pay Runs and Payslips management.
 */
@Component({
  selector: 'msls-payroll',
  standalone: true,
  imports: [CommonModule, RouterLink],
  template: `
    <div class="page">
      <!-- Page Header -->
      <div class="page-header">
        <div class="header-content">
          <div class="header-icon">
            <i class="fa-solid fa-money-check-dollar"></i>
          </div>
          <div class="header-text">
            <h1>Payroll Management</h1>
            <p>Process staff salaries, manage pay runs, and generate payslips</p>
          </div>
        </div>
      </div>

      <!-- Configuration Cards -->
      <div class="config-grid">
        <!-- Pay Runs Card -->
        <a routerLink="runs" class="config-card featured">
          <div class="card-icon runs">
            <i class="fa-solid fa-play-circle"></i>
          </div>
          <div class="card-content">
            <h3>Pay Runs</h3>
            <p>Create and process monthly payroll runs for all staff</p>
          </div>
          <div class="card-arrow">
            <i class="fa-solid fa-chevron-right"></i>
          </div>
        </a>

        <!-- Salary Configuration Card -->
        <a routerLink="/admin/salary-components" class="config-card">
          <div class="card-icon salary">
            <i class="fa-solid fa-indian-rupee-sign"></i>
          </div>
          <div class="card-content">
            <h3>Salary Components</h3>
            <p>Manage salary component types (allowances, deductions) used across staff</p>
          </div>
          <div class="card-arrow">
            <i class="fa-solid fa-chevron-right"></i>
          </div>
        </a>

        <!-- Reports Card -->
        <a routerLink="runs" class="config-card">
          <div class="card-icon reports">
            <i class="fa-solid fa-chart-line"></i>
          </div>
          <div class="card-content">
            <h3>Payroll Reports</h3>
            <p>Generate salary reports, bank files, and tax summaries</p>
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
              <h4>Configure Salary Structure</h4>
              <p>Set up basic pay, allowances (HRA, DA, TA), and deductions (PF, ESI, TDS) for each staff</p>
            </div>
          </div>
          <div class="info-step">
            <div class="step-number">2</div>
            <div class="step-content">
              <h4>Create Pay Run</h4>
              <p>Start a new pay run for a specific month and year, optionally filter by branch</p>
            </div>
          </div>
          <div class="info-step">
            <div class="step-number">3</div>
            <div class="step-content">
              <h4>Calculate Payroll</h4>
              <p>System calculates gross pay, applies deductions, and generates net salary for all staff</p>
            </div>
          </div>
          <div class="info-step">
            <div class="step-number">4</div>
            <div class="step-content">
              <h4>Review & Adjust</h4>
              <p>Review calculated salaries, make adjustments for leave, overtime, or bonuses as needed</p>
            </div>
          </div>
          <div class="info-step">
            <div class="step-number">5</div>
            <div class="step-content">
              <h4>Approve Pay Run</h4>
              <p>Manager reviews and approves the pay run after verification</p>
            </div>
          </div>
          <div class="info-step">
            <div class="step-number">6</div>
            <div class="step-content">
              <h4>Finalize & Export</h4>
              <p>Finalize the pay run, generate payslips, and export bank transfer file</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Status Flow -->
      <div class="status-flow-section">
        <h2>Pay Run Status Flow</h2>
        <div class="status-flow">
          <div class="status-item draft">
            <div class="status-icon"><i class="fa-solid fa-file-pen"></i></div>
            <span class="status-label">Draft</span>
            <span class="status-desc">Created, not calculated</span>
          </div>
          <div class="status-arrow"><i class="fa-solid fa-arrow-right"></i></div>
          <div class="status-item calculated">
            <div class="status-icon"><i class="fa-solid fa-calculator"></i></div>
            <span class="status-label">Calculated</span>
            <span class="status-desc">Salaries computed</span>
          </div>
          <div class="status-arrow"><i class="fa-solid fa-arrow-right"></i></div>
          <div class="status-item approved">
            <div class="status-icon"><i class="fa-solid fa-check-circle"></i></div>
            <span class="status-label">Approved</span>
            <span class="status-desc">Manager approved</span>
          </div>
          <div class="status-arrow"><i class="fa-solid fa-arrow-right"></i></div>
          <div class="status-item finalized">
            <div class="status-icon"><i class="fa-solid fa-lock"></i></div>
            <span class="status-label">Finalized</span>
            <span class="status-desc">Locked, payslips ready</span>
          </div>
        </div>
      </div>
    </div>
  `,
  styles: [`
    .page { padding: 1.5rem; max-width: 1200px; margin: 0 auto; }
    .page-header { margin-bottom: 2rem; }
    .header-content { display: flex; align-items: center; gap: 1rem; }
    .header-icon { width: 3.5rem; height: 3.5rem; border-radius: 1rem; background: linear-gradient(135deg, #16a34a, #15803d); color: white; display: flex; align-items: center; justify-content: center; font-size: 1.5rem; }
    .header-text h1 { margin: 0; font-size: 1.75rem; font-weight: 600; color: #1e293b; }
    .header-text p { margin: 0.25rem 0 0; color: #64748b; font-size: 0.9375rem; }

    .config-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 1.25rem; margin-bottom: 2.5rem; }
    .config-card { display: flex; align-items: center; gap: 1rem; padding: 1.25rem; background: white; border: 1px solid #e2e8f0; border-radius: 1rem; text-decoration: none; transition: all 0.2s; }
    .config-card:hover { border-color: #bbf7d0; box-shadow: 0 4px 12px rgba(22, 163, 74, 0.08); transform: translateY(-2px); }
    .card-icon { width: 3rem; height: 3rem; border-radius: 0.75rem; display: flex; align-items: center; justify-content: center; font-size: 1.25rem; flex-shrink: 0; }
    .card-icon.runs { background: linear-gradient(135deg, #16a34a, #15803d); color: white; }
    .card-icon.salary { background: #dbeafe; color: #2563eb; }
    .card-icon.reports { background: #fef3c7; color: #d97706; }
    .config-card.featured { border: 2px solid #bbf7d0; background: linear-gradient(135deg, #f0fdf4, #dcfce7); }
    .config-card.featured:hover { border-color: #86efac; box-shadow: 0 4px 16px rgba(22, 163, 74, 0.15); }
    .card-content { flex: 1; }
    .card-content h3 { margin: 0 0 0.25rem; font-size: 1rem; font-weight: 600; color: #1e293b; }
    .card-content p { margin: 0; font-size: 0.8125rem; color: #64748b; line-height: 1.4; }
    .card-arrow { color: #94a3b8; transition: transform 0.2s; }
    .config-card:hover .card-arrow { transform: translateX(4px); color: #16a34a; }

    .info-section { background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 1rem; padding: 1.5rem; margin-bottom: 1.5rem; }
    .info-section h2 { margin: 0 0 1.25rem; font-size: 1.125rem; font-weight: 600; color: #1e293b; }
    .info-steps { display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap: 1.25rem; }
    .info-step { display: flex; gap: 0.75rem; }
    .step-number { width: 2rem; height: 2rem; border-radius: 50%; background: #16a34a; color: white; display: flex; align-items: center; justify-content: center; font-weight: 600; font-size: 0.875rem; flex-shrink: 0; }
    .step-content h4 { margin: 0 0 0.25rem; font-size: 0.875rem; font-weight: 600; color: #1e293b; }
    .step-content p { margin: 0; font-size: 0.75rem; color: #64748b; line-height: 1.4; }

    .status-flow-section { background: white; border: 1px solid #e2e8f0; border-radius: 1rem; padding: 1.5rem; }
    .status-flow-section h2 { margin: 0 0 1.25rem; font-size: 1.125rem; font-weight: 600; color: #1e293b; }
    .status-flow { display: flex; align-items: center; justify-content: center; gap: 0.5rem; flex-wrap: wrap; }
    .status-item { display: flex; flex-direction: column; align-items: center; gap: 0.375rem; padding: 1rem; border-radius: 0.75rem; min-width: 120px; }
    .status-item.draft { background: #f1f5f9; }
    .status-item.calculated { background: #fef3c7; }
    .status-item.approved { background: #dbeafe; }
    .status-item.finalized { background: #dcfce7; }
    .status-icon { width: 2.5rem; height: 2.5rem; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 1rem; }
    .status-item.draft .status-icon { background: #e2e8f0; color: #64748b; }
    .status-item.calculated .status-icon { background: #fde68a; color: #92400e; }
    .status-item.approved .status-icon { background: #93c5fd; color: #1e40af; }
    .status-item.finalized .status-icon { background: #86efac; color: #166534; }
    .status-label { font-weight: 600; font-size: 0.875rem; color: #1e293b; }
    .status-desc { font-size: 0.6875rem; color: #64748b; text-align: center; }
    .status-arrow { color: #cbd5e1; font-size: 1rem; }

    @media (max-width: 768px) {
      .status-flow { flex-direction: column; gap: 0.25rem; }
      .status-arrow { transform: rotate(90deg); }
      .status-item { width: 100%; flex-direction: row; gap: 1rem; padding: 0.75rem 1rem; }
      .status-item .status-desc { text-align: left; }
    }
  `],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class PayrollComponent {}
