/**
 * Pay Run Detail Component
 * Story 5.6: Payroll Processing
 *
 * Displays pay run details with payslip list.
 */

import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { PayrollService } from '../../payroll.service';
import {
  PayRun,
  Payslip,
  PayRunStatus,
  PayRunSummary,
  getPayRunStatusLabel,
  getPayslipStatusLabel,
  formatPayPeriod,
  formatCurrency,
} from '../../payroll.model';
import { ToastService } from '../../../../shared/services/toast.service';

@Component({
  selector: 'msls-pay-run-detail',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterLink],
  template: `
    <div class="page">
      <!-- Breadcrumb -->
      <div class="breadcrumb">
        <a routerLink="/payroll/runs">
          <i class="fa-solid fa-arrow-left"></i>
          Back to Pay Runs
        </a>
      </div>

      @if (loading()) {
        <div class="loading-container">
          <div class="spinner"></div>
          <span>Loading pay run...</span>
        </div>
      } @else if (error()) {
        <div class="error-container">
          <i class="fa-solid fa-circle-exclamation"></i>
          <span>{{ error() }}</span>
          <button class="btn btn-secondary btn-sm" (click)="loadPayRun()">
            <i class="fa-solid fa-refresh"></i>
            Retry
          </button>
        </div>
      } @else if (payRun()) {
        <!-- Page Header -->
        <div class="page-header">
          <div class="header-content">
            <div class="header-icon">
              <i class="fa-solid fa-calendar-days"></i>
            </div>
            <div class="header-text">
              <h1>{{ formatPayPeriod(payRun()!.payPeriodMonth, payRun()!.payPeriodYear) }}</h1>
              <p>
                {{ payRun()!.branchName || 'All Branches' }}
                <span class="badge" [class]="'badge-' + getStatusClass(payRun()!.status)">
                  {{ getPayRunStatusLabel(payRun()!.status) }}
                </span>
              </p>
            </div>
          </div>
          <div class="header-actions">
            @if (payRun()!.status === 'draft') {
              <button class="btn btn-primary" (click)="calculatePayroll()">
                <i class="fa-solid fa-calculator"></i>
                Calculate Payroll
              </button>
            }
            @if (payRun()!.status === 'calculated') {
              <button class="btn btn-success" (click)="approvePayRun()">
                <i class="fa-solid fa-check"></i>
                Approve
              </button>
            }
            @if (payRun()!.status === 'approved') {
              <button class="btn btn-success" (click)="finalizePayRun()">
                <i class="fa-solid fa-lock"></i>
                Finalize
              </button>
            }
            @if (payRun()!.status === 'finalized') {
              <button class="btn btn-secondary" (click)="exportBankFile()">
                <i class="fa-solid fa-file-export"></i>
                Export Bank File
              </button>
            }
          </div>
        </div>

        <!-- Summary Cards -->
        <div class="summary-cards">
          <div class="summary-card">
            <div class="card-icon staff">
              <i class="fa-solid fa-users"></i>
            </div>
            <div class="card-content">
              <span class="card-label">Total Staff</span>
              <span class="card-value">{{ payRun()!.totalStaff }}</span>
            </div>
          </div>
          <div class="summary-card">
            <div class="card-icon gross">
              <i class="fa-solid fa-indian-rupee-sign"></i>
            </div>
            <div class="card-content">
              <span class="card-label">Gross Salary</span>
              <span class="card-value">{{ formatCurrency(payRun()!.totalGross) }}</span>
            </div>
          </div>
          <div class="summary-card">
            <div class="card-icon deduction">
              <i class="fa-solid fa-arrow-down"></i>
            </div>
            <div class="card-content">
              <span class="card-label">Total Deductions</span>
              <span class="card-value">{{ formatCurrency(payRun()!.totalDeductions) }}</span>
            </div>
          </div>
          <div class="summary-card highlight">
            <div class="card-icon net">
              <i class="fa-solid fa-wallet"></i>
            </div>
            <div class="card-content">
              <span class="card-label">Net Payable</span>
              <span class="card-value">{{ formatCurrency(payRun()!.totalNet) }}</span>
            </div>
          </div>
        </div>

        <!-- Tabs -->
        <div class="tabs">
          <button
            class="tab"
            [class.active]="activeTab() === 'payslips'"
            (click)="activeTab.set('payslips')"
          >
            <i class="fa-solid fa-file-invoice-dollar"></i>
            Payslips
          </button>
          <button
            class="tab"
            [class.active]="activeTab() === 'summary'"
            (click)="activeTab.set('summary'); loadSummary()"
          >
            <i class="fa-solid fa-chart-pie"></i>
            Department Summary
          </button>
        </div>

        <!-- Payslips Tab -->
        @if (activeTab() === 'payslips') {
          <div class="content-card">
            <!-- Search -->
            <div class="search-bar">
              <div class="search-box">
                <i class="fa-solid fa-search search-icon"></i>
                <input
                  type="text"
                  placeholder="Search by name or employee ID..."
                  [ngModel]="searchTerm()"
                  (ngModelChange)="searchTerm.set($event)"
                  class="search-input"
                />
                @if (searchTerm()) {
                  <button class="clear-search" (click)="searchTerm.set('')">
                    <i class="fa-solid fa-xmark"></i>
                  </button>
                }
              </div>
            </div>

            @if (loadingPayslips()) {
              <div class="loading-container">
                <div class="spinner"></div>
                <span>Loading payslips...</span>
              </div>
            } @else {
              <div class="table-wrapper">
                <table class="data-table">
                  <thead>
                    <tr>
                      <th>Employee</th>
                      <th class="hide-lg" style="text-align: center;">Working</th>
                      <th class="hide-md" style="text-align: center;">Present</th>
                      <th class="hide-lg" style="text-align: center;">LOP</th>
                      <th class="hide-md" style="text-align: right;">Gross</th>
                      <th class="hide-lg" style="text-align: right;">Deductions</th>
                      <th style="text-align: right;">Net</th>
                      <th class="hide-sm">Status</th>
                      <th style="text-align: right;">Actions</th>
                    </tr>
                  </thead>
                  <tbody>
                    @for (payslip of filteredPayslips(); track payslip.id) {
                      <tr>
                        <td class="employee-cell">
                          <div class="employee-wrapper">
                            <div class="employee-avatar">
                              {{ getInitials(payslip.staffName || '') }}
                            </div>
                            <div class="employee-content">
                              <span class="employee-name">{{ payslip.staffName }}</span>
                              <span class="employee-id">{{ payslip.staffEmployeeId }}</span>
                            </div>
                          </div>
                        </td>
                        <td class="days-cell hide-lg">{{ payslip.workingDays }}</td>
                        <td class="days-cell hide-md">
                          <span class="days-present">{{ payslip.presentDays }}</span>
                        </td>
                        <td class="days-cell hide-lg">
                          @if (payslip.lopDays > 0) {
                            <span class="days-lop">{{ payslip.lopDays }}</span>
                          } @else {
                            <span>-</span>
                          }
                        </td>
                        <td class="amount-cell hide-md">{{ formatCurrency(payslip.grossSalary) }}</td>
                        <td class="amount-cell deduction hide-lg">{{ formatCurrency(payslip.totalDeductions) }}</td>
                        <td class="amount-cell net">{{ formatCurrency(payslip.netSalary) }}</td>
                        <td class="hide-sm">
                          <span class="badge" [class]="'badge-' + getPayslipStatusClass(payslip.status)">
                            {{ getPayslipStatusLabel(payslip.status) }}
                          </span>
                        </td>
                        <td class="actions-cell">
                          <button class="action-btn" title="View Payslip" (click)="viewPayslip(payslip)">
                            <i class="fa-regular fa-eye"></i>
                          </button>
                          <button class="action-btn" title="Download PDF" (click)="downloadPayslip(payslip)">
                            <i class="fa-solid fa-download"></i>
                          </button>
                        </td>
                      </tr>
                    } @empty {
                      <tr>
                        <td colspan="9" class="empty-cell">
                          <div class="empty-state">
                            <i class="fa-regular fa-file-invoice"></i>
                            @if (payRun()!.status === 'draft') {
                              <p>No payslips yet</p>
                              <span>Calculate payroll to generate payslips</span>
                            } @else {
                              <p>No payslips found</p>
                            }
                          </div>
                        </td>
                      </tr>
                    }
                  </tbody>
                </table>
              </div>
            }
          </div>
        }

        <!-- Summary Tab -->
        @if (activeTab() === 'summary') {
          <div class="content-card">
            @if (loadingSummary()) {
              <div class="loading-container">
                <div class="spinner"></div>
                <span>Loading summary...</span>
              </div>
            } @else if (summary()) {
              <div class="table-wrapper">
                <table class="data-table">
                  <thead>
                    <tr>
                      <th>Department</th>
                      <th class="hide-sm" style="text-align: center;">Staff</th>
                      <th class="hide-md" style="text-align: right;">Gross</th>
                      <th class="hide-lg" style="text-align: right;">Deductions</th>
                      <th style="text-align: right;">Net</th>
                    </tr>
                  </thead>
                  <tbody>
                    @for (dept of summary()!.departmentSummaries; track dept.departmentId) {
                      <tr>
                        <td>
                          <div class="dept-wrapper">
                            <div class="dept-icon">
                              <i class="fa-solid fa-building"></i>
                            </div>
                            <span class="dept-name">{{ dept.departmentName }}</span>
                          </div>
                        </td>
                        <td class="count-cell hide-sm">
                          <span class="count-badge">{{ dept.staffCount }}</span>
                        </td>
                        <td class="amount-cell hide-md">{{ formatCurrency(dept.totalGross) }}</td>
                        <td class="amount-cell deduction hide-lg">{{ formatCurrency(dept.totalDeductions) }}</td>
                        <td class="amount-cell net">{{ formatCurrency(dept.totalNet) }}</td>
                      </tr>
                    } @empty {
                      <tr>
                        <td colspan="5" class="empty-cell">
                          <div class="empty-state">
                            <i class="fa-regular fa-building"></i>
                            <p>No department data available</p>
                          </div>
                        </td>
                      </tr>
                    }
                  </tbody>
                  <tfoot>
                    <tr class="total-row">
                      <td><strong>Total</strong></td>
                      <td class="count-cell hide-sm">
                        <strong>{{ summary()!.totalStaff }}</strong>
                      </td>
                      <td class="amount-cell hide-md"><strong>{{ formatCurrency(summary()!.totalGross) }}</strong></td>
                      <td class="amount-cell deduction hide-lg"><strong>{{ formatCurrency(summary()!.totalDeductions) }}</strong></td>
                      <td class="amount-cell net"><strong>{{ formatCurrency(summary()!.totalNet) }}</strong></td>
                    </tr>
                  </tfoot>
                </table>
              </div>
            }
          </div>
        }

        <!-- Notes Section -->
        @if (payRun()!.notes) {
          <div class="notes-section">
            <h3>
              <i class="fa-regular fa-note-sticky"></i>
              Notes
            </h3>
            <p>{{ payRun()!.notes }}</p>
          </div>
        }

        <!-- Audit Info -->
        <div class="audit-info">
          <div class="audit-item">
            <span class="audit-label">Created</span>
            <span class="audit-value">{{ formatDateTime(payRun()!.createdAt) }}</span>
          </div>
          @if (payRun()!.calculatedAt) {
            <div class="audit-item">
              <span class="audit-label">Calculated</span>
              <span class="audit-value">{{ formatDateTime(payRun()!.calculatedAt!) }}</span>
            </div>
          }
          @if (payRun()!.approvedAt) {
            <div class="audit-item">
              <span class="audit-label">Approved</span>
              <span class="audit-value">
                {{ formatDateTime(payRun()!.approvedAt!) }}
                @if (payRun()!.approvedByName) {
                  by {{ payRun()!.approvedByName }}
                }
              </span>
            </div>
          }
          @if (payRun()!.finalizedAt) {
            <div class="audit-item">
              <span class="audit-label">Finalized</span>
              <span class="audit-value">
                {{ formatDateTime(payRun()!.finalizedAt!) }}
                @if (payRun()!.finalizedByName) {
                  by {{ payRun()!.finalizedByName }}
                }
              </span>
            </div>
          }
        </div>
      }

      <!-- Processing Modal -->
      @if (processing()) {
        <div class="modal-overlay">
          <div class="modal modal--sm">
            <div class="modal__header">
              <h3>
                <i class="fa-solid fa-spinner fa-spin"></i>
                Processing
              </h3>
            </div>
            <div class="modal__body">
              <div class="processing-modal">
                <div class="processing-spinner"></div>
                <p>{{ processingMessage() }}</p>
                <span class="processing-hint">Please wait...</span>
              </div>
            </div>
          </div>
        </div>
      }
    </div>
  `,
  styles: [`
    .page {
      padding: 1.5rem;
      max-width: 1400px;
      margin: 0 auto;
    }

    .breadcrumb {
      margin-bottom: 1rem;
    }

    .breadcrumb a {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      color: #64748b;
      text-decoration: none;
      font-size: 0.875rem;
    }

    .breadcrumb a:hover {
      color: #4f46e5;
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

    .page-header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: 1.5rem;
    }

    .header-content {
      display: flex;
      align-items: center;
      gap: 1rem;
    }

    .header-icon {
      width: 3rem;
      height: 3rem;
      border-radius: 0.75rem;
      background: #dbeafe;
      color: #2563eb;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 1.25rem;
    }

    .header-text h1 {
      margin: 0;
      font-size: 1.5rem;
      font-weight: 600;
      color: #1e293b;
    }

    .header-text p {
      margin: 0.25rem 0 0;
      color: #64748b;
      font-size: 0.875rem;
      display: flex;
      align-items: center;
      gap: 0.75rem;
    }

    .header-actions {
      display: flex;
      gap: 0.75rem;
    }

    .summary-cards {
      display: grid;
      grid-template-columns: repeat(4, 1fr);
      gap: 1rem;
      margin-bottom: 1.5rem;
    }

    .summary-card {
      display: flex;
      align-items: center;
      gap: 1rem;
      padding: 1.25rem;
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 0.75rem;
    }

    .summary-card.highlight {
      background: linear-gradient(135deg, #f0fdf4 0%, #dcfce7 100%);
      border-color: #86efac;
    }

    .card-icon {
      width: 3rem;
      height: 3rem;
      border-radius: 0.75rem;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 1.25rem;
    }

    .card-icon.staff {
      background: #f1f5f9;
      color: #475569;
    }

    .card-icon.gross {
      background: #dbeafe;
      color: #2563eb;
    }

    .card-icon.deduction {
      background: #fef2f2;
      color: #dc2626;
    }

    .card-icon.net {
      background: #dcfce7;
      color: #16a34a;
    }

    .card-content {
      display: flex;
      flex-direction: column;
    }

    .card-label {
      font-size: 0.75rem;
      color: #64748b;
      text-transform: uppercase;
      letter-spacing: 0.05em;
    }

    .card-value {
      font-size: 1.25rem;
      font-weight: 600;
      color: #1e293b;
    }

    .tabs {
      display: flex;
      gap: 0.5rem;
      margin-bottom: 1rem;
    }

    .tab {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.75rem 1.25rem;
      border: none;
      background: transparent;
      color: #64748b;
      font-size: 0.875rem;
      font-weight: 500;
      border-radius: 0.5rem;
      cursor: pointer;
      transition: all 0.2s;
    }

    .tab:hover {
      background: #f1f5f9;
    }

    .tab.active {
      background: #4f46e5;
      color: white;
    }

    .content-card {
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 1rem;
      overflow: hidden;
    }

    .search-bar {
      padding: 1rem;
      border-bottom: 1px solid #e2e8f0;
    }

    .search-box {
      max-width: 400px;
      position: relative;
    }

    .search-icon {
      position: absolute;
      left: 0.875rem;
      top: 50%;
      transform: translateY(-50%);
      color: #9ca3af;
    }

    .search-input {
      width: 100%;
      padding: 0.625rem 2.5rem 0.625rem 2.5rem;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      font-size: 0.875rem;
    }

    .search-input:focus {
      outline: none;
      border-color: #4f46e5;
      box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
    }

    .clear-search {
      position: absolute;
      right: 0.5rem;
      top: 50%;
      transform: translateY(-50%);
      background: none;
      border: none;
      color: #9ca3af;
      cursor: pointer;
      padding: 0.25rem;
    }

    .data-table {
      width: 100%;
      border-collapse: collapse;
    }

    .data-table th {
      text-align: left;
      padding: 0.875rem 1rem;
      font-size: 0.75rem;
      font-weight: 600;
      text-transform: uppercase;
      letter-spacing: 0.05em;
      color: #64748b;
      background: #f8fafc;
      border-bottom: 1px solid #e2e8f0;
    }

    .data-table td {
      padding: 1rem;
      border-bottom: 1px solid #f1f5f9;
      color: #374151;
    }

    .data-table tbody tr:hover {
      background: #f8fafc;
    }

    .employee-cell {
      min-width: 200px;
    }

    .employee-wrapper {
      display: flex;
      align-items: center;
      gap: 0.75rem;
    }

    .employee-avatar {
      width: 2.5rem;
      height: 2.5rem;
      border-radius: 50%;
      background: #e0e7ff;
      color: #4f46e5;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 0.75rem;
      font-weight: 600;
      flex-shrink: 0;
    }

    .employee-content {
      display: flex;
      flex-direction: column;
    }

    .employee-name {
      font-weight: 500;
      color: #1e293b;
    }

    .employee-id {
      font-size: 0.75rem;
      color: #64748b;
      font-family: monospace;
    }

    .days-cell {
      text-align: center;
    }

    .days-present {
      color: #16a34a;
      font-weight: 500;
    }

    .days-lop {
      color: #dc2626;
      font-weight: 500;
    }

    .count-cell {
      text-align: center;
    }

    .count-badge {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      min-width: 2rem;
      height: 1.5rem;
      padding: 0 0.5rem;
      background: #f1f5f9;
      border-radius: 9999px;
      font-size: 0.75rem;
      font-weight: 500;
      color: #475569;
    }

    .amount-cell {
      text-align: right;
      font-family: monospace;
      font-size: 0.875rem;
    }

    .amount-cell.deduction {
      color: #dc2626;
    }

    .amount-cell.net {
      font-weight: 600;
      color: #16a34a;
    }

    .badge {
      display: inline-flex;
      padding: 0.25rem 0.75rem;
      border-radius: 9999px;
      font-size: 0.75rem;
      font-weight: 500;
    }

    .badge-gray { background: #f1f5f9; color: #64748b; }
    .badge-blue { background: #dbeafe; color: #1e40af; }
    .badge-yellow { background: #fef3c7; color: #92400e; }
    .badge-green { background: #dcfce7; color: #166534; }
    .badge-red { background: #fef2f2; color: #991b1b; }

    .actions-cell {
      text-align: right;
    }

    .action-btn {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      width: 2rem;
      height: 2rem;
      border: none;
      background: transparent;
      color: #64748b;
      border-radius: 0.375rem;
      cursor: pointer;
      transition: all 0.2s;
    }

    .action-btn:hover {
      background: #f1f5f9;
      color: #4f46e5;
    }

    .empty-cell {
      padding: 3rem !important;
    }

    .empty-state {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 0.5rem;
      color: #64748b;
    }

    .empty-state i {
      font-size: 2rem;
      color: #cbd5e1;
    }

    .empty-state span {
      font-size: 0.875rem;
    }

    .dept-wrapper {
      display: flex;
      align-items: center;
      gap: 0.75rem;
    }

    .dept-icon {
      width: 2rem;
      height: 2rem;
      border-radius: 0.375rem;
      background: #f1f5f9;
      color: #475569;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 0.75rem;
    }

    .dept-name {
      font-weight: 500;
      color: #1e293b;
    }

    .total-row {
      background: #f8fafc;
    }

    .total-row td {
      border-top: 2px solid #e2e8f0;
      padding-top: 1rem;
      padding-bottom: 1rem;
    }

    .notes-section {
      margin-top: 1.5rem;
      padding: 1rem;
      background: #fffbeb;
      border: 1px solid #fde68a;
      border-radius: 0.75rem;
    }

    .notes-section h3 {
      margin: 0 0 0.5rem;
      font-size: 0.875rem;
      font-weight: 600;
      color: #92400e;
      display: flex;
      align-items: center;
      gap: 0.5rem;
    }

    .notes-section p {
      margin: 0;
      color: #78350f;
      font-size: 0.875rem;
    }

    .audit-info {
      display: flex;
      flex-wrap: wrap;
      gap: 2rem;
      margin-top: 1.5rem;
      padding-top: 1rem;
      border-top: 1px solid #e2e8f0;
    }

    .audit-item {
      display: flex;
      flex-direction: column;
      gap: 0.25rem;
    }

    .audit-label {
      font-size: 0.75rem;
      color: #64748b;
      text-transform: uppercase;
      letter-spacing: 0.05em;
    }

    .audit-value {
      font-size: 0.875rem;
      color: #374151;
    }

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

    .btn-success {
      background: #16a34a;
      color: white;
    }

    .btn-success:hover {
      background: #15803d;
    }

    /* Modal Styles */
    .modal-overlay {
      position: fixed;
      inset: 0;
      background: rgba(0, 0, 0, 0.5);
      backdrop-filter: blur(4px);
      display: flex;
      align-items: center;
      justify-content: center;
      z-index: 1000;
      padding: 1rem;
    }

    .modal {
      background: white;
      border-radius: 1.25rem;
      width: 100%;
      box-shadow: 0 25px 50px rgba(0, 0, 0, 0.25);
    }

    .modal--sm { max-width: 28rem; }

    .modal__header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 1.25rem 1.5rem;
      border-bottom: 1px solid #f1f5f9;

      h3 {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        font-size: 1.125rem;
        font-weight: 600;
        color: #0f172a;
        margin: 0;

        i { color: #4f46e5; }
      }
    }

    .modal__body {
      padding: 1.5rem;
    }

    /* Processing Modal */
    .processing-modal {
      text-align: center;
      padding: 2rem;
    }

    .processing-spinner {
      width: 3rem;
      height: 3rem;
      margin: 0 auto 1rem;
      border: 4px solid #e2e8f0;
      border-top-color: #4f46e5;
      border-radius: 50%;
      animation: spin 0.8s linear infinite;
    }

    .processing-modal p {
      margin: 0 0 0.5rem;
      font-weight: 500;
      color: #1e293b;
    }

    .processing-hint {
      font-size: 0.875rem;
      color: #64748b;
    }

    /* Responsive Table */
    .table-wrapper {
      overflow-x: auto;
      -webkit-overflow-scrolling: touch;
    }

    @media (max-width: 1024px) {
      .summary-cards {
        grid-template-columns: repeat(2, 1fr);
      }

      .data-table .hide-lg {
        display: none;
      }
    }

    @media (max-width: 768px) {
      .page-header {
        flex-direction: column;
        gap: 1rem;
      }

      .header-actions {
        width: 100%;
        flex-wrap: wrap;
      }

      .header-actions .btn {
        flex: 1;
        justify-content: center;
        min-width: 120px;
      }

      .summary-cards {
        grid-template-columns: repeat(2, 1fr);
      }

      .summary-card {
        padding: 1rem;
      }

      .card-icon {
        width: 2.5rem;
        height: 2.5rem;
        font-size: 1rem;
      }

      .card-value {
        font-size: 1rem;
      }

      .tabs {
        flex-wrap: wrap;
      }

      .tab {
        flex: 1;
        justify-content: center;
        min-width: 120px;
      }

      .data-table .hide-md {
        display: none;
      }

      .data-table th,
      .data-table td {
        padding: 0.75rem 0.5rem;
        font-size: 0.8125rem;
      }

      .employee-cell {
        min-width: 140px;
      }

      .employee-avatar {
        width: 2rem;
        height: 2rem;
        font-size: 0.625rem;
      }

      .amount-cell {
        font-size: 0.75rem;
      }
    }

    @media (max-width: 480px) {
      .page {
        padding: 1rem;
      }

      .summary-cards {
        grid-template-columns: 1fr;
      }

      .data-table .hide-sm {
        display: none;
      }

      .data-table th,
      .data-table td {
        padding: 0.5rem 0.375rem;
        font-size: 0.75rem;
      }

      .badge {
        padding: 0.125rem 0.5rem;
        font-size: 0.625rem;
      }

      .search-input {
        font-size: 16px; /* Prevent iOS zoom */
      }
    }
  `],
})
export class PayRunDetailComponent implements OnInit {
  private payrollService = inject(PayrollService);
  private toastService = inject(ToastService);
  private route = inject(ActivatedRoute);
  private router = inject(Router);

  // State signals
  payRun = signal<PayRun | null>(null);
  payslips = signal<Payslip[]>([]);
  summary = signal<PayRunSummary | null>(null);
  loading = signal(true);
  loadingPayslips = signal(false);
  loadingSummary = signal(false);
  processing = signal(false);
  processingMessage = signal('');
  error = signal<string | null>(null);
  searchTerm = signal('');
  activeTab = signal<'payslips' | 'summary'>('payslips');

  // Computed filtered payslips
  filteredPayslips = computed(() => {
    const term = this.searchTerm().toLowerCase();
    if (!term) return this.payslips();

    return this.payslips().filter(
      p =>
        p.staffName?.toLowerCase().includes(term) ||
        p.staffEmployeeId?.toLowerCase().includes(term)
    );
  });

  ngOnInit(): void {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.loadPayRun(id);
    }
  }

  loadPayRun(id?: string): void {
    const payRunId = id || this.payRun()?.id;
    if (!payRunId) return;

    this.loading.set(true);
    this.error.set(null);

    this.payrollService.getPayRun(payRunId).subscribe({
      next: payRun => {
        this.payRun.set(payRun);
        this.loading.set(false);
        this.loadPayslips();
      },
      error: () => {
        this.error.set('Failed to load pay run details');
        this.loading.set(false);
      },
    });
  }

  loadPayslips(): void {
    const payRun = this.payRun();
    if (!payRun) return;

    this.loadingPayslips.set(true);

    this.payrollService.getPayslips(payRun.id).subscribe({
      next: payslips => {
        this.payslips.set(payslips);
        this.loadingPayslips.set(false);
      },
      error: () => {
        this.payslips.set([]);
        this.loadingPayslips.set(false);
      },
    });
  }

  loadSummary(): void {
    const payRun = this.payRun();
    if (!payRun || this.summary()) return;

    this.loadingSummary.set(true);

    this.payrollService.getPayRunSummary(payRun.id).subscribe({
      next: summary => {
        this.summary.set(summary);
        this.loadingSummary.set(false);
      },
      error: () => {
        this.loadingSummary.set(false);
      },
    });
  }

  calculatePayroll(): void {
    const payRun = this.payRun();
    if (!payRun) return;

    this.processing.set(true);
    this.processingMessage.set('Calculating payroll for all staff...');

    this.payrollService.calculatePayroll(payRun.id).subscribe({
      next: () => {
        this.toastService.success('Payroll calculated successfully');
        this.processing.set(false);
        this.loadPayRun();
      },
      error: err => {
        const message = err?.error?.message || 'Failed to calculate payroll';
        this.toastService.error(message);
        this.processing.set(false);
      },
    });
  }

  approvePayRun(): void {
    const payRun = this.payRun();
    if (!payRun) return;

    this.processing.set(true);
    this.processingMessage.set('Approving pay run...');

    this.payrollService.approvePayRun(payRun.id).subscribe({
      next: () => {
        this.toastService.success('Pay run approved successfully');
        this.processing.set(false);
        this.loadPayRun();
      },
      error: err => {
        const message = err?.error?.message || 'Failed to approve pay run';
        this.toastService.error(message);
        this.processing.set(false);
      },
    });
  }

  finalizePayRun(): void {
    const payRun = this.payRun();
    if (!payRun) return;

    this.processing.set(true);
    this.processingMessage.set('Finalizing pay run...');

    this.payrollService.finalizePayRun(payRun.id).subscribe({
      next: () => {
        this.toastService.success('Pay run finalized successfully');
        this.processing.set(false);
        this.loadPayRun();
      },
      error: err => {
        const message = err?.error?.message || 'Failed to finalize pay run';
        this.toastService.error(message);
        this.processing.set(false);
      },
    });
  }

  exportBankFile(): void {
    const payRun = this.payRun();
    if (!payRun) return;

    this.payrollService.exportBankFile(payRun.id).subscribe({
      next: blob => {
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `bank-transfer-${payRun.payPeriodYear}-${payRun.payPeriodMonth}.csv`;
        a.click();
        window.URL.revokeObjectURL(url);
        this.toastService.success('Bank file downloaded');
      },
      error: () => {
        this.toastService.error('Failed to download bank file');
      },
    });
  }

  viewPayslip(payslip: Payslip): void {
    this.router.navigate(['/payroll/payslip', payslip.id]);
  }

  downloadPayslip(payslip: Payslip): void {
    const payRun = this.payRun();
    this.payrollService.downloadPayslipPdf(payslip.id).subscribe({
      next: blob => {
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        // Format: payslip_monthyear.pdf
        const monthNames = ['', 'jan', 'feb', 'mar', 'apr', 'may', 'jun', 'jul', 'aug', 'sep', 'oct', 'nov', 'dec'];
        const month = payRun ? monthNames[payRun.payPeriodMonth] || 'unknown' : 'unknown';
        const year = payRun?.payPeriodYear || new Date().getFullYear();
        a.download = `payslip_${month}${year}.pdf`;
        a.click();
        window.URL.revokeObjectURL(url);
      },
      error: () => {
        this.toastService.error('Failed to download payslip');
      },
    });
  }

  getStatusClass(status: PayRunStatus): string {
    switch (status) {
      case 'draft': return 'gray';
      case 'processing': return 'blue';
      case 'calculated': return 'yellow';
      case 'approved': return 'blue';
      case 'finalized': return 'green';
      case 'reversed': return 'red';
      default: return 'gray';
    }
  }

  getPayslipStatusClass(status: string): string {
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

  formatPayPeriod = formatPayPeriod;
  formatCurrency = formatCurrency;
  getPayRunStatusLabel = getPayRunStatusLabel;
  getPayslipStatusLabel = getPayslipStatusLabel;

  formatDateTime(dateStr: string): string {
    if (!dateStr) return '-';
    return new Date(dateStr).toLocaleString('en-IN', {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  }
}
