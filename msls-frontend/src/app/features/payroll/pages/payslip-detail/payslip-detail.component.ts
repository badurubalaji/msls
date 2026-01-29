/**
 * Payslip Detail Component
 * Story 5.6: Payroll Processing
 *
 * Displays individual payslip with component breakdown.
 */

import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';
import { PayrollService } from '../../payroll.service';
import {
  Payslip,
  PayslipComponent,
  getPayslipStatusLabel,
  formatCurrency,
} from '../../payroll.model';
import { ToastService } from '../../../../shared/services/toast.service';

@Component({
  selector: 'msls-payslip-detail',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="page">
      @if (loading()) {
        <div class="loading-container">
          <div class="spinner"></div>
          <span>Loading payslip...</span>
        </div>
      } @else if (error()) {
        <div class="error-container">
          <i class="fa-solid fa-circle-exclamation"></i>
          <span>{{ error() }}</span>
          <button class="btn btn-secondary btn-sm" (click)="loadPayslip()">
            <i class="fa-solid fa-refresh"></i>
            Retry
          </button>
        </div>
      } @else if (payslip()) {
        <!-- Payslip Document -->
        <div class="payslip-document">
          <!-- Header -->
          <div class="payslip-header">
            <div class="company-info">
              <h1>Payslip</h1>
              <p class="company-name">My School System</p>
            </div>
            <div class="payslip-meta">
              <div class="meta-item">
                <span class="meta-label">Employee ID</span>
                <span class="meta-value">{{ payslip()!.staffEmployeeId }}</span>
              </div>
              <div class="meta-item">
                <span class="meta-label">Status</span>
                <span class="badge" [class]="'badge-' + getStatusClass(payslip()!.status)">
                  {{ getPayslipStatusLabel(payslip()!.status) }}
                </span>
              </div>
            </div>
          </div>

          <!-- Employee Details -->
          <div class="employee-section">
            <div class="employee-info">
              <div class="employee-avatar">
                {{ getInitials(payslip()!.staffName || '') }}
              </div>
              <div class="employee-details">
                <h2>{{ payslip()!.staffName }}</h2>
                <p>{{ payslip()!.staffEmployeeId }}</p>
              </div>
            </div>
          </div>

          <!-- Attendance Summary -->
          <div class="attendance-section">
            <h3>Attendance Summary</h3>
            <div class="attendance-grid">
              <div class="attendance-item">
                <span class="att-value">{{ payslip()!.workingDays }}</span>
                <span class="att-label">Working Days</span>
              </div>
              <div class="attendance-item present">
                <span class="att-value">{{ payslip()!.presentDays }}</span>
                <span class="att-label">Present</span>
              </div>
              <div class="attendance-item leave">
                <span class="att-value">{{ payslip()!.leaveDays }}</span>
                <span class="att-label">Leave</span>
              </div>
              <div class="attendance-item absent">
                <span class="att-value">{{ payslip()!.absentDays }}</span>
                <span class="att-label">Absent</span>
              </div>
              @if (payslip()!.lopDays > 0) {
                <div class="attendance-item lop">
                  <span class="att-value">{{ payslip()!.lopDays }}</span>
                  <span class="att-label">LOP Days</span>
                </div>
              }
            </div>
          </div>

          <!-- Salary Breakdown -->
          <div class="breakdown-section">
            <div class="breakdown-columns">
              <!-- Earnings -->
              <div class="breakdown-column">
                <h3>
                  <i class="fa-solid fa-arrow-up"></i>
                  Earnings
                </h3>
                <div class="component-list">
                  @for (comp of earningComponents(); track comp.id) {
                    <div class="component-item">
                      <span class="comp-name">
                        {{ comp.componentName }}
                        @if (comp.isProrated) {
                          <span class="prorated">(Prorated)</span>
                        }
                      </span>
                      <span class="comp-amount">{{ formatCurrency(comp.amount) }}</span>
                    </div>
                  } @empty {
                    <div class="no-components">No earnings</div>
                  }
                </div>
                <div class="component-total">
                  <span>Total Earnings</span>
                  <span>{{ formatCurrency(payslip()!.totalEarnings) }}</span>
                </div>
              </div>

              <!-- Deductions -->
              <div class="breakdown-column deductions">
                <h3>
                  <i class="fa-solid fa-arrow-down"></i>
                  Deductions
                </h3>
                <div class="component-list">
                  @for (comp of deductionComponents(); track comp.id) {
                    <div class="component-item">
                      <span class="comp-name">
                        {{ comp.componentName }}
                        @if (comp.isProrated) {
                          <span class="prorated">(Prorated)</span>
                        }
                      </span>
                      <span class="comp-amount">{{ formatCurrency(comp.amount) }}</span>
                    </div>
                  } @empty {
                    <div class="no-components">No deductions</div>
                  }
                  @if (parseFloat(payslip()!.lopDeduction) > 0) {
                    <div class="component-item lop-item">
                      <span class="comp-name">
                        LOP Deduction
                        <span class="prorated">({{ payslip()!.lopDays }} days)</span>
                      </span>
                      <span class="comp-amount">{{ formatCurrency(payslip()!.lopDeduction) }}</span>
                    </div>
                  }
                </div>
                <div class="component-total">
                  <span>Total Deductions</span>
                  <span>{{ formatCurrency(payslip()!.totalDeductions) }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Net Pay -->
          <div class="net-pay-section">
            <div class="net-pay-row">
              <div class="net-pay-item">
                <span class="net-label">Gross Salary</span>
                <span class="net-value gross">{{ formatCurrency(payslip()!.grossSalary) }}</span>
              </div>
              <div class="net-pay-item">
                <span class="net-label">Total Deductions</span>
                <span class="net-value deduction">- {{ formatCurrency(payslip()!.totalDeductions) }}</span>
              </div>
              <div class="net-pay-item highlight">
                <span class="net-label">Net Pay</span>
                <span class="net-value">{{ formatCurrency(payslip()!.netSalary) }}</span>
              </div>
            </div>
          </div>

          <!-- Payment Info -->
          @if (payslip()!.paymentDate || payslip()!.paymentReference) {
            <div class="payment-section">
              <h3>Payment Information</h3>
              <div class="payment-grid">
                @if (payslip()!.paymentDate) {
                  <div class="payment-item">
                    <span class="payment-label">Payment Date</span>
                    <span class="payment-value">{{ formatDate(payslip()!.paymentDate!) }}</span>
                  </div>
                }
                @if (payslip()!.paymentReference) {
                  <div class="payment-item">
                    <span class="payment-label">Reference</span>
                    <span class="payment-value">{{ payslip()!.paymentReference }}</span>
                  </div>
                }
              </div>
            </div>
          }

          <!-- Footer -->
          <div class="payslip-footer">
            <p>This is a system-generated payslip and does not require a signature.</p>
            <p class="generated-date">Generated on {{ formatDate(payslip()!.createdAt) }}</p>
          </div>
        </div>

        <!-- Actions -->
        <div class="page-actions">
          <button class="btn btn-secondary" (click)="goBack()">
            <i class="fa-solid fa-arrow-left"></i>
            Back
          </button>
          <button class="btn btn-primary" (click)="downloadPdf()">
            <i class="fa-solid fa-download"></i>
            Download PDF
          </button>
          <button class="btn btn-secondary" (click)="print()">
            <i class="fa-solid fa-print"></i>
            Print
          </button>
        </div>
      }
    </div>
  `,
  styles: [`
    .page {
      padding: 1.5rem;
      max-width: 900px;
      margin: 0 auto;
    }

    .loading-container,
    .error-container {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 1rem;
      padding: 3rem;
      color: #64748b;
    }

    .error-container {
      color: #dc2626;
      flex-direction: column;
    }

    .spinner {
      width: 24px;
      height: 24px;
      border: 3px solid #e2e8f0;
      border-top-color: #4f46e5;
      border-radius: 50%;
      animation: spin 0.8s linear infinite;
    }

    @keyframes spin {
      to { transform: rotate(360deg); }
    }

    .payslip-document {
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 1rem;
      overflow: hidden;
      margin-bottom: 1.5rem;
    }

    .payslip-header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      padding: 1.5rem;
      background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
      color: white;
    }

    .company-info h1 {
      margin: 0;
      font-size: 1.5rem;
      font-weight: 600;
    }

    .company-name {
      margin: 0.25rem 0 0;
      opacity: 0.9;
    }

    .payslip-meta {
      display: flex;
      gap: 1.5rem;
    }

    .meta-item {
      display: flex;
      flex-direction: column;
      align-items: flex-end;
    }

    .meta-label {
      font-size: 0.75rem;
      opacity: 0.8;
    }

    .meta-value {
      font-weight: 500;
    }

    .employee-section {
      padding: 1.5rem;
      border-bottom: 1px solid #e2e8f0;
    }

    .employee-info {
      display: flex;
      align-items: center;
      gap: 1rem;
    }

    .employee-avatar {
      width: 4rem;
      height: 4rem;
      border-radius: 50%;
      background: #e0e7ff;
      color: #4f46e5;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 1.25rem;
      font-weight: 600;
    }

    .employee-details h2 {
      margin: 0;
      font-size: 1.25rem;
      color: #1e293b;
    }

    .employee-details p {
      margin: 0.25rem 0 0;
      color: #64748b;
      font-family: monospace;
    }

    .attendance-section {
      padding: 1.5rem;
      border-bottom: 1px solid #e2e8f0;
    }

    .attendance-section h3 {
      margin: 0 0 1rem;
      font-size: 0.875rem;
      font-weight: 600;
      color: #374151;
      text-transform: uppercase;
      letter-spacing: 0.05em;
    }

    .attendance-grid {
      display: flex;
      gap: 1.5rem;
      flex-wrap: wrap;
    }

    .attendance-item {
      display: flex;
      flex-direction: column;
      align-items: center;
      padding: 0.75rem 1.25rem;
      background: #f8fafc;
      border-radius: 0.5rem;
      min-width: 80px;
    }

    .attendance-item.present { background: #dcfce7; }
    .attendance-item.leave { background: #fef3c7; }
    .attendance-item.absent { background: #fee2e2; }
    .attendance-item.lop { background: #fef2f2; }

    .att-value {
      font-size: 1.5rem;
      font-weight: 600;
      color: #1e293b;
    }

    .attendance-item.present .att-value { color: #16a34a; }
    .attendance-item.leave .att-value { color: #ca8a04; }
    .attendance-item.absent .att-value { color: #dc2626; }
    .attendance-item.lop .att-value { color: #991b1b; }

    .att-label {
      font-size: 0.75rem;
      color: #64748b;
      margin-top: 0.25rem;
    }

    .breakdown-section {
      padding: 1.5rem;
      border-bottom: 1px solid #e2e8f0;
    }

    .breakdown-columns {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 2rem;
    }

    .breakdown-column h3 {
      margin: 0 0 1rem;
      font-size: 0.875rem;
      font-weight: 600;
      color: #16a34a;
      display: flex;
      align-items: center;
      gap: 0.5rem;
    }

    .breakdown-column.deductions h3 {
      color: #dc2626;
    }

    .component-list {
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
    }

    .component-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 0.5rem 0;
      border-bottom: 1px dashed #e2e8f0;
    }

    .component-item:last-child {
      border-bottom: none;
    }

    .comp-name {
      color: #374151;
      font-size: 0.875rem;
    }

    .prorated {
      color: #64748b;
      font-size: 0.75rem;
      font-style: italic;
    }

    .comp-amount {
      font-family: monospace;
      font-weight: 500;
      color: #1e293b;
    }

    .lop-item {
      background: #fef2f2;
      margin: 0 -0.5rem;
      padding: 0.5rem;
      border-radius: 0.25rem;
    }

    .lop-item .comp-amount {
      color: #dc2626;
    }

    .component-total {
      display: flex;
      justify-content: space-between;
      margin-top: 1rem;
      padding-top: 1rem;
      border-top: 2px solid #e2e8f0;
      font-weight: 600;
    }

    .no-components {
      color: #64748b;
      font-size: 0.875rem;
      font-style: italic;
      padding: 1rem 0;
    }

    .net-pay-section {
      padding: 1.5rem;
      background: #f8fafc;
    }

    .net-pay-row {
      display: flex;
      justify-content: space-between;
      align-items: center;
      gap: 2rem;
    }

    .net-pay-item {
      display: flex;
      flex-direction: column;
      align-items: center;
    }

    .net-pay-item.highlight {
      padding: 1rem 2rem;
      background: linear-gradient(135deg, #16a34a 0%, #15803d 100%);
      border-radius: 0.75rem;
      color: white;
    }

    .net-label {
      font-size: 0.75rem;
      color: #64748b;
      text-transform: uppercase;
      letter-spacing: 0.05em;
    }

    .net-pay-item.highlight .net-label {
      color: rgba(255, 255, 255, 0.8);
    }

    .net-value {
      font-size: 1.25rem;
      font-weight: 600;
      font-family: monospace;
      color: #1e293b;
    }

    .net-value.gross {
      color: #16a34a;
    }

    .net-value.deduction {
      color: #dc2626;
    }

    .net-pay-item.highlight .net-value {
      color: white;
      font-size: 1.5rem;
    }

    .payment-section {
      padding: 1.5rem;
      border-bottom: 1px solid #e2e8f0;
    }

    .payment-section h3 {
      margin: 0 0 1rem;
      font-size: 0.875rem;
      font-weight: 600;
      color: #374151;
    }

    .payment-grid {
      display: flex;
      gap: 2rem;
    }

    .payment-item {
      display: flex;
      flex-direction: column;
    }

    .payment-label {
      font-size: 0.75rem;
      color: #64748b;
    }

    .payment-value {
      font-weight: 500;
      color: #1e293b;
    }

    .payslip-footer {
      padding: 1rem 1.5rem;
      background: #f8fafc;
      text-align: center;
    }

    .payslip-footer p {
      margin: 0;
      font-size: 0.75rem;
      color: #64748b;
    }

    .generated-date {
      margin-top: 0.25rem !important;
    }

    .page-actions {
      display: flex;
      justify-content: center;
      gap: 1rem;
    }

    .badge {
      display: inline-flex;
      padding: 0.25rem 0.75rem;
      border-radius: 9999px;
      font-size: 0.75rem;
      font-weight: 500;
    }

    .badge-gray { background: rgba(255,255,255,0.2); color: white; }
    .badge-yellow { background: #fef3c7; color: #92400e; }
    .badge-blue { background: #dbeafe; color: #1e40af; }
    .badge-green { background: #dcfce7; color: #166534; }

    /* Buttons */
    .btn {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.625rem 1.25rem;
      border-radius: 0.5rem;
      font-size: 0.875rem;
      font-weight: 500;
      cursor: pointer;
      transition: all 0.2s;
      border: none;
    }

    .btn-sm {
      padding: 0.375rem 0.75rem;
      font-size: 0.75rem;
    }

    .btn-primary {
      background: #4f46e5;
      color: white;
    }

    .btn-primary:hover {
      background: #4338ca;
    }

    .btn-secondary {
      background: #f1f5f9;
      color: #475569;
    }

    .btn-secondary:hover {
      background: #e2e8f0;
    }

    @media print {
      .page-actions {
        display: none;
      }

      .payslip-document {
        border: none;
        box-shadow: none;
      }
    }

    @media (max-width: 768px) {
      .payslip-header {
        flex-direction: column;
        gap: 1rem;
      }

      .payslip-meta {
        width: 100%;
        justify-content: space-between;
      }

      .meta-item {
        align-items: flex-start;
      }

      .breakdown-columns {
        grid-template-columns: 1fr;
      }

      .net-pay-row {
        flex-direction: column;
        gap: 1rem;
      }

      .net-pay-item {
        width: 100%;
        flex-direction: row;
        justify-content: space-between;
      }

      .net-pay-item.highlight {
        flex-direction: column;
        text-align: center;
      }

      .page-actions {
        flex-direction: column;
      }

      .page-actions .btn {
        width: 100%;
        justify-content: center;
      }
    }
  `],
})
export class PayslipDetailComponent implements OnInit {
  private payrollService = inject(PayrollService);
  private toastService = inject(ToastService);
  private route = inject(ActivatedRoute);
  private router = inject(Router);

  // State signals
  payslip = signal<Payslip | null>(null);
  loading = signal(true);
  error = signal<string | null>(null);

  // Computed components
  earningComponents = computed(() => {
    return (this.payslip()?.components || []).filter(c => c.componentType === 'earning');
  });

  deductionComponents = computed(() => {
    return (this.payslip()?.components || []).filter(c => c.componentType === 'deduction');
  });

  ngOnInit(): void {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.loadPayslip(id);
    }
  }

  loadPayslip(id?: string): void {
    const payslipId = id || this.payslip()?.id;
    if (!payslipId) return;

    this.loading.set(true);
    this.error.set(null);

    this.payrollService.getPayslip(payslipId).subscribe({
      next: payslip => {
        this.payslip.set(payslip);
        this.loading.set(false);
      },
      error: () => {
        this.error.set('Failed to load payslip');
        this.loading.set(false);
      },
    });
  }

  goBack(): void {
    const payslip = this.payslip();
    if (payslip?.payRunId) {
      this.router.navigate(['/payroll', payslip.payRunId]);
    } else {
      this.router.navigate(['/payroll']);
    }
  }

  downloadPdf(): void {
    const payslip = this.payslip();
    if (!payslip) return;

    // Fetch pay run to get the period info for filename
    this.payrollService.getPayRun(payslip.payRunId).subscribe({
      next: payRun => {
        this.payrollService.downloadPayslipPdf(payslip.id).subscribe({
          next: blob => {
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            // Format: payslip_monthyear.pdf
            const monthNames = ['', 'jan', 'feb', 'mar', 'apr', 'may', 'jun', 'jul', 'aug', 'sep', 'oct', 'nov', 'dec'];
            const month = monthNames[payRun.payPeriodMonth] || 'unknown';
            a.download = `payslip_${month}${payRun.payPeriodYear}.pdf`;
            a.click();
            window.URL.revokeObjectURL(url);
          },
          error: () => {
            this.toastService.error('Failed to download payslip');
          },
        });
      },
      error: () => {
        this.toastService.error('Failed to download payslip');
      },
    });
  }

  print(): void {
    window.print();
  }

  getStatusClass(status: string): string {
    switch (status) {
      case 'calculated': return 'gray';
      case 'adjusted': return 'yellow';
      case 'approved': return 'blue';
      case 'paid': return 'green';
      default: return 'gray';
    }
  }

  getInitials(name: string): string {
    return name
      .split(' ')
      .map(n => n[0])
      .join('')
      .substring(0, 2)
      .toUpperCase();
  }

  parseFloat = parseFloat;
  formatCurrency = formatCurrency;
  getPayslipStatusLabel = getPayslipStatusLabel;

  formatDate(dateStr: string): string {
    if (!dateStr) return '-';
    return new Date(dateStr).toLocaleDateString('en-IN', {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
    });
  }
}
