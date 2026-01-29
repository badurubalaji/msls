/**
 * Pay Runs List Component
 * Story 5.6: Payroll Processing
 *
 * Displays list of pay runs with filtering and actions.
 */

import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { CommonModule, Location } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { PayrollService } from '../../payroll.service';
import {
  PayRun,
  PayRunStatus,
  CreatePayRunRequest,
  getPayRunStatusLabel,
  formatPayPeriod,
  formatCurrency,
} from '../../payroll.model';
import { ToastService } from '../../../../shared/services/toast.service';
import { BranchService } from '../../../admin/branches/branch.service';
import { Branch } from '../../../admin/branches/branch.model';

@Component({
  selector: 'msls-pay-runs',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterLink],
  template: `
    <div class="page">
      <!-- Back Navigation -->
      <div class="breadcrumb">
        <a routerLink="/dashboard">
          <i class="fa-solid fa-arrow-left"></i>
          Back to Dashboard
        </a>
      </div>

      <!-- Page Header -->
      <div class="page-header">
        <div class="header-content">
          <div class="header-icon">
            <i class="fa-solid fa-money-check-dollar"></i>
          </div>
          <div class="header-text">
            <h1>Payroll Processing</h1>
            <p>Manage monthly payroll runs and process staff salaries</p>
          </div>
        </div>
        <button class="btn btn-primary" (click)="openCreateModal()">
          <i class="fa-solid fa-plus"></i>
          New Pay Run
        </button>
      </div>

      <!-- Filters -->
      <div class="filters-bar">
        <div class="filter-group">
          <label class="filter-label">Year</label>
          <select
            class="filter-select"
            [ngModel]="yearFilter()"
            (ngModelChange)="yearFilter.set($event); loadPayRuns()"
          >
            @for (year of availableYears; track year) {
              <option [value]="year">{{ year }}</option>
            }
          </select>
        </div>
        <div class="filter-group">
          <label class="filter-label">Month</label>
          <select
            class="filter-select"
            [ngModel]="monthFilter()"
            (ngModelChange)="monthFilter.set($event); loadPayRuns()"
          >
            <option value="">All Months</option>
            @for (month of months; track month.value) {
              <option [value]="month.value">{{ month.name }}</option>
            }
          </select>
        </div>
        <div class="filter-group">
          <label class="filter-label">Status</label>
          <select
            class="filter-select"
            [ngModel]="statusFilter()"
            (ngModelChange)="statusFilter.set($event); loadPayRuns()"
          >
            <option value="">All Status</option>
            <option value="draft">Draft</option>
            <option value="calculated">Calculated</option>
            <option value="approved">Approved</option>
            <option value="finalized">Finalized</option>
          </select>
        </div>
        @if (branches().length > 1) {
          <div class="filter-group">
            <label class="filter-label">Branch</label>
            <select
              class="filter-select"
              [ngModel]="branchFilter()"
              (ngModelChange)="branchFilter.set($event); loadPayRuns()"
            >
              <option value="">All Branches</option>
              @for (branch of branches(); track branch.id) {
                <option [value]="branch.id">{{ branch.name }}</option>
              }
            </select>
          </div>
        }
      </div>

      <!-- Content -->
      <div class="content-card">
        @if (loading()) {
          <div class="loading-container">
            <div class="spinner"></div>
            <span>Loading pay runs...</span>
          </div>
        } @else if (error()) {
          <div class="error-container">
            <i class="fa-solid fa-circle-exclamation"></i>
            <span>{{ error() }}</span>
            <button class="btn btn-secondary btn-sm" (click)="loadPayRuns()">
              <i class="fa-solid fa-refresh"></i>
              Retry
            </button>
          </div>
        } @else {
          <div class="table-wrapper">
            <table class="data-table">
              <thead>
                <tr>
                  <th>Pay Period</th>
                  <th class="hide-md">Branch</th>
                  <th class="hide-sm" style="text-align: center;">Staff</th>
                  <th class="hide-lg" style="text-align: right;">Gross</th>
                  <th class="hide-lg" style="text-align: right;">Deductions</th>
                  <th style="text-align: right;">Net</th>
                  <th>Status</th>
                  <th style="text-align: right;">Actions</th>
                </tr>
              </thead>
              <tbody>
                @for (run of payRuns(); track run.id) {
                  <tr>
                    <td class="period-cell">
                      <div class="period-wrapper">
                        <div class="period-icon">
                          <i class="fa-solid fa-calendar-days"></i>
                        </div>
                        <div class="period-content">
                          <span class="period-name">{{ formatPayPeriod(run.payPeriodMonth, run.payPeriodYear) }}</span>
                          <span class="period-date">Created {{ formatDate(run.createdAt) }}</span>
                        </div>
                      </div>
                    </td>
                    <td class="hide-md">{{ run.branchName || 'All Branches' }}</td>
                    <td class="count-cell hide-sm">
                      <span class="count-badge">{{ run.totalStaff }}</span>
                    </td>
                    <td class="amount-cell hide-lg">{{ formatCurrency(run.totalGross) }}</td>
                    <td class="amount-cell deduction hide-lg">{{ formatCurrency(run.totalDeductions) }}</td>
                    <td class="amount-cell net">{{ formatCurrency(run.totalNet) }}</td>
                    <td>
                      <span class="badge" [class]="'badge-' + getStatusClass(run.status)">
                        {{ getPayRunStatusLabel(run.status) }}
                      </span>
                    </td>
                    <td class="actions-cell">
                    <button class="action-btn" title="View Details" (click)="viewPayRun(run)">
                      <i class="fa-regular fa-eye"></i>
                    </button>
                    @if (run.status === 'draft') {
                      <button class="action-btn" title="Calculate Payroll" (click)="calculatePayroll(run)">
                        <i class="fa-solid fa-calculator"></i>
                      </button>
                    }
                    @if (run.status === 'calculated') {
                      <button class="action-btn action-btn--success" title="Approve" (click)="approvePayRun(run)">
                        <i class="fa-solid fa-check"></i>
                      </button>
                    }
                    @if (run.status === 'approved') {
                      <button class="action-btn action-btn--success" title="Finalize" (click)="finalizePayRun(run)">
                        <i class="fa-solid fa-lock"></i>
                      </button>
                    }
                    @if (run.status === 'finalized') {
                      <button class="action-btn" title="Export Bank File" (click)="exportBankFile(run)">
                        <i class="fa-solid fa-file-export"></i>
                      </button>
                    }
                    @if (run.status === 'draft') {
                      <button
                        class="action-btn action-btn--danger"
                        title="Delete"
                        (click)="confirmDelete(run)"
                      >
                        <i class="fa-regular fa-trash-can"></i>
                      </button>
                    }
                  </td>
                </tr>
                } @empty {
                  <tr>
                    <td colspan="8" class="empty-cell">
                      <div class="empty-state">
                        <i class="fa-regular fa-calendar-xmark"></i>
                        <p>No pay runs found</p>
                        <button class="btn btn-primary btn-sm" (click)="openCreateModal()">
                          <i class="fa-solid fa-plus"></i>
                          Create Pay Run
                        </button>
                      </div>
                    </td>
                  </tr>
                }
              </tbody>
            </table>
          </div>
        }
      </div>

      <!-- Create Pay Run Modal -->
      @if (showCreateModal()) {
        <div class="modal-overlay" (click)="closeCreateModal()">
          <div class="modal modal--md" (click)="$event.stopPropagation()">
            <div class="modal__header">
              <h3>
                <i class="fa-solid fa-plus-circle"></i>
                Create Pay Run
              </h3>
              <button type="button" class="modal__close" (click)="closeCreateModal()">
                <i class="fa-solid fa-xmark"></i>
              </button>
            </div>
            <div class="modal__body">
              <form (ngSubmit)="createPayRun()" class="create-form">
                <div class="form-row">
                  <div class="form-group">
                    <label class="label required">Month</label>
                    <select
                      class="select"
                      [(ngModel)]="newPayRun.payPeriodMonth"
                      name="month"
                      required
                    >
                      @for (month of months; track month.value) {
                        <option [value]="month.value">{{ month.name }}</option>
                      }
                    </select>
                  </div>
                  <div class="form-group">
                    <label class="label required">Year</label>
                    <select
                      class="select"
                      [(ngModel)]="newPayRun.payPeriodYear"
                      name="year"
                      required
                    >
                      @for (year of availableYears; track year) {
                        <option [value]="year">{{ year }}</option>
                      }
                    </select>
                  </div>
                </div>

                @if (branches().length > 1) {
                  <div class="form-group">
                    <label class="label">Branch (Optional)</label>
                    <select
                      class="select"
                      [(ngModel)]="newPayRun.branchId"
                      name="branchId"
                    >
                      <option value="">All Branches</option>
                      @for (branch of branches(); track branch.id) {
                        <option [value]="branch.id">{{ branch.name }}</option>
                      }
                    </select>
                    <span class="hint">Leave empty to include all branches</span>
                  </div>
                }

                <div class="form-group">
                  <label class="label">Notes</label>
                  <textarea
                    class="textarea"
                    [(ngModel)]="newPayRun.notes"
                    name="notes"
                    rows="3"
                    placeholder="Optional notes about this pay run..."
                  ></textarea>
                </div>
              </form>
            </div>
            <div class="modal__footer">
              <button type="button" class="btn btn-secondary" (click)="closeCreateModal()">
                Cancel
              </button>
              <button type="button" class="btn btn-primary" [disabled]="creating()" (click)="createPayRun()">
                @if (creating()) {
                  <div class="btn-spinner"></div>
                  Creating...
                } @else {
                  <i class="fa-solid fa-plus"></i>
                  Create Pay Run
                }
              </button>
            </div>
          </div>
        </div>
      }

      <!-- Delete Confirmation Modal -->
      @if (showDeleteModal()) {
        <div class="modal-overlay" (click)="closeDeleteModal()">
          <div class="modal modal--sm" (click)="$event.stopPropagation()">
            <div class="modal__header">
              <h3>
                <i class="fa-solid fa-trash"></i>
                Delete Pay Run
              </h3>
              <button type="button" class="modal__close" (click)="closeDeleteModal()">
                <i class="fa-solid fa-xmark"></i>
              </button>
            </div>
            <div class="modal__body">
              <div class="delete-confirmation">
                <div class="delete-icon">
                  <i class="fa-solid fa-triangle-exclamation"></i>
                </div>
                <p>
                  Are you sure you want to delete the pay run for
                  <strong>{{ payRunToDelete() ? formatPayPeriod(payRunToDelete()!.payPeriodMonth, payRunToDelete()!.payPeriodYear) : '' }}</strong>?
                </p>
                <p class="delete-warning">This action cannot be undone.</p>
              </div>
            </div>
            <div class="modal__footer">
              <button type="button" class="btn btn-secondary" (click)="closeDeleteModal()">Cancel</button>
              <button type="button" class="btn btn-danger" [disabled]="deleting()" (click)="deletePayRun()">
                @if (deleting()) {
                  <div class="btn-spinner"></div>
                  Deleting...
                } @else {
                  <i class="fa-solid fa-trash"></i>
                  Delete
                }
              </button>
            </div>
          </div>
        </div>
      }

      <!-- Processing Modal -->
      @if (processing()) {
        <div class="modal-overlay">
          <div class="modal modal--sm" (click)="$event.stopPropagation()">
            <div class="modal__header">
              <h3>
                <i class="fa-solid fa-spinner fa-spin"></i>
                Processing Payroll
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
      transition: color 0.2s;
    }

    .breadcrumb a:hover {
      color: #4f46e5;
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
      background: #dcfce7;
      color: #16a34a;
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
    }

    .filters-bar {
      display: flex;
      gap: 1rem;
      margin-bottom: 1rem;
      flex-wrap: wrap;
    }

    .filter-group {
      display: flex;
      flex-direction: column;
      gap: 0.25rem;
    }

    .filter-label {
      font-size: 0.75rem;
      font-weight: 500;
      color: #64748b;
    }

    .filter-select {
      padding: 0.5rem 2rem 0.5rem 0.75rem;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      font-size: 0.875rem;
      background: white;
      cursor: pointer;
      min-width: 140px;
    }

    .filter-select:focus {
      outline: none;
      border-color: #4f46e5;
      box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
    }

    .content-card {
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 1rem;
      overflow: hidden;
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

    .period-cell {
      min-width: 200px;
    }

    .period-wrapper {
      display: flex;
      align-items: center;
      gap: 0.75rem;
    }

    .period-icon {
      width: 2.5rem;
      height: 2.5rem;
      border-radius: 0.5rem;
      background: #dbeafe;
      color: #2563eb;
      display: flex;
      align-items: center;
      justify-content: center;
      flex-shrink: 0;
    }

    .period-content {
      display: flex;
      flex-direction: column;
    }

    .period-name {
      font-weight: 500;
      color: #1e293b;
    }

    .period-date {
      font-size: 0.75rem;
      color: #64748b;
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

    .badge-gray {
      background: #f1f5f9;
      color: #64748b;
    }

    .badge-blue {
      background: #dbeafe;
      color: #1e40af;
    }

    .badge-yellow {
      background: #fef3c7;
      color: #92400e;
    }

    .badge-green {
      background: #dcfce7;
      color: #166534;
    }

    .badge-red {
      background: #fef2f2;
      color: #991b1b;
    }

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

    .action-btn:hover:not(:disabled) {
      background: #f1f5f9;
      color: #4f46e5;
    }

    .action-btn--success:hover:not(:disabled) {
      background: #dcfce7;
      color: #16a34a;
    }

    .action-btn--danger:hover:not(:disabled) {
      background: #fef2f2;
      color: #dc2626;
    }

    .empty-cell {
      padding: 3rem !important;
    }

    .empty-state {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 0.75rem;
      color: #64748b;
    }

    .empty-state i {
      font-size: 2rem;
      color: #cbd5e1;
    }

    /* Form Styles */
    .create-form {
      display: flex;
      flex-direction: column;
      gap: 1rem;
    }

    .form-row {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 1rem;
    }

    .form-group {
      display: flex;
      flex-direction: column;
      gap: 0.375rem;
    }

    .label {
      font-size: 0.875rem;
      font-weight: 500;
      color: #374151;
    }

    .label.required::after {
      content: ' *';
      color: #dc2626;
    }

    .select,
    .textarea {
      padding: 0.625rem 0.875rem;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      font-size: 0.875rem;
    }

    .select:focus,
    .textarea:focus {
      outline: none;
      border-color: #4f46e5;
      box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
    }

    .hint {
      font-size: 0.75rem;
      color: #64748b;
    }

    .form-actions {
      display: flex;
      justify-content: flex-end;
      gap: 0.75rem;
      padding-top: 1rem;
      border-top: 1px solid #e2e8f0;
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

    .btn-primary:hover:not(:disabled) {
      background: #4338ca;
    }

    .btn-secondary {
      background: #f1f5f9;
      color: #475569;
    }

    .btn-secondary:hover {
      background: #e2e8f0;
    }

    .btn-danger {
      background: #dc2626;
      color: white;
    }

    .btn-danger:hover:not(:disabled) {
      background: #b91c1c;
    }

    .btn:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }

    .btn-spinner {
      width: 16px;
      height: 16px;
      border: 2px solid transparent;
      border-top-color: currentColor;
      border-radius: 50%;
      animation: spin 0.8s linear infinite;
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
      max-height: 90vh;
      overflow: hidden;
      display: flex;
      flex-direction: column;
      box-shadow: 0 25px 50px rgba(0, 0, 0, 0.25);
    }

    .modal--sm { max-width: 28rem; }
    .modal--md { max-width: 36rem; }
    .modal--lg { max-width: 48rem; }

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

        i {
          color: #4f46e5;
        }
      }
    }

    .modal__close {
      width: 2rem;
      height: 2rem;
      display: flex;
      align-items: center;
      justify-content: center;
      background: #f1f5f9;
      border: none;
      border-radius: 0.5rem;
      cursor: pointer;
      color: #64748b;
      transition: all 0.2s ease;

      &:hover {
        background: #e2e8f0;
        color: #334155;
      }
    }

    .modal__body {
      padding: 1.5rem;
      overflow-y: auto;
      flex: 1;
    }

    .modal__footer {
      display: flex;
      justify-content: flex-end;
      gap: 0.75rem;
      padding: 1.25rem 1.5rem;
      background: #f8fafc;
      border-top: 1px solid #f1f5f9;
      border-radius: 0 0 1.25rem 1.25rem;
    }

    /* Delete Modal */
    .delete-confirmation {
      text-align: center;
      padding: 1rem;
    }

    .delete-icon {
      width: 4rem;
      height: 4rem;
      margin: 0 auto 1rem;
      border-radius: 50%;
      background: #fef2f2;
      color: #dc2626;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 1.5rem;
    }

    .delete-confirmation p {
      margin: 0 0 0.5rem;
      color: #374151;
    }

    .delete-warning {
      font-size: 0.875rem;
      color: #64748b;
    }

    .delete-actions {
      display: flex;
      gap: 0.75rem;
      justify-content: center;
      margin-top: 1.5rem;
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
      .data-table .hide-lg {
        display: none;
      }
    }

    @media (max-width: 768px) {
      .page-header {
        flex-direction: column;
        gap: 1rem;
      }

      .page-header .btn {
        width: 100%;
        justify-content: center;
      }

      .filters-bar {
        flex-direction: column;
      }

      .filter-group {
        width: 100%;
      }

      .filter-select {
        width: 100%;
      }

      .form-row {
        grid-template-columns: 1fr;
      }

      .data-table .hide-md {
        display: none;
      }

      .data-table th,
      .data-table td {
        padding: 0.75rem 0.5rem;
        font-size: 0.8125rem;
      }

      .period-cell {
        min-width: 150px;
      }

      .period-icon {
        display: none;
      }

      .amount-cell {
        font-size: 0.75rem;
      }

      .action-btn {
        width: 1.75rem;
        height: 1.75rem;
      }
    }

    @media (max-width: 480px) {
      .page {
        padding: 1rem;
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

      .actions-cell {
        white-space: nowrap;
      }
    }
  `],
})
export class PayRunsComponent implements OnInit {
  private payrollService = inject(PayrollService);
  private branchService = inject(BranchService);
  private toastService = inject(ToastService);
  private router = inject(Router);

  // State signals
  payRuns = signal<PayRun[]>([]);
  branches = signal<Branch[]>([]);
  loading = signal(true);
  creating = signal(false);
  deleting = signal(false);
  processing = signal(false);
  processingMessage = signal('');
  error = signal<string | null>(null);

  // Filters
  yearFilter = signal(new Date().getFullYear());
  monthFilter = signal<number | ''>('');
  statusFilter = signal<string>('');
  branchFilter = signal<string>('');

  // Modal state
  showCreateModal = signal(false);
  showDeleteModal = signal(false);
  payRunToDelete = signal<PayRun | null>(null);

  // New pay run form
  newPayRun: CreatePayRunRequest = {
    payPeriodMonth: new Date().getMonth() + 1,
    payPeriodYear: new Date().getFullYear(),
    branchId: undefined,
    notes: undefined,
  };

  // Static data
  months = [
    { value: 1, name: 'January' },
    { value: 2, name: 'February' },
    { value: 3, name: 'March' },
    { value: 4, name: 'April' },
    { value: 5, name: 'May' },
    { value: 6, name: 'June' },
    { value: 7, name: 'July' },
    { value: 8, name: 'August' },
    { value: 9, name: 'September' },
    { value: 10, name: 'October' },
    { value: 11, name: 'November' },
    { value: 12, name: 'December' },
  ];

  availableYears = Array.from({ length: 5 }, (_, i) => new Date().getFullYear() - i);

  ngOnInit(): void {
    this.loadPayRuns();
    this.loadBranches();
  }

  loadPayRuns(): void {
    this.loading.set(true);
    this.error.set(null);

    this.payrollService
      .getPayRuns({
        year: this.yearFilter(),
        month: this.monthFilter() || undefined,
        status: this.statusFilter() || undefined,
        branchId: this.branchFilter() || undefined,
      })
      .subscribe({
        next: payRuns => {
          this.payRuns.set(payRuns);
          this.loading.set(false);
        },
        error: () => {
          this.error.set('Failed to load pay runs. Please try again.');
          this.loading.set(false);
        },
      });
  }

  loadBranches(): void {
    this.branchService.getBranches().subscribe({
      next: branches => this.branches.set(branches),
      error: () => this.branches.set([]),
    });
  }

  openCreateModal(): void {
    this.newPayRun = {
      payPeriodMonth: new Date().getMonth() + 1,
      payPeriodYear: new Date().getFullYear(),
      branchId: undefined,
      notes: undefined,
    };
    this.showCreateModal.set(true);
  }

  closeCreateModal(): void {
    this.showCreateModal.set(false);
  }

  createPayRun(): void {
    this.creating.set(true);

    const request: CreatePayRunRequest = {
      payPeriodMonth: Number(this.newPayRun.payPeriodMonth),
      payPeriodYear: Number(this.newPayRun.payPeriodYear),
      branchId: this.newPayRun.branchId || undefined,
      notes: this.newPayRun.notes?.trim() || undefined,
    };

    this.payrollService.createPayRun(request).subscribe({
      next: payRun => {
        this.toastService.success('Pay run created successfully');
        this.closeCreateModal();
        this.loadPayRuns();
        this.creating.set(false);
        // Navigate to pay run detail
        this.router.navigate(['/payroll', payRun.id]);
      },
      error: err => {
        const message = err?.error?.message || 'Failed to create pay run';
        this.toastService.error(message);
        this.creating.set(false);
      },
    });
  }

  viewPayRun(run: PayRun): void {
    this.router.navigate(['/payroll', run.id]);
  }

  calculatePayroll(run: PayRun): void {
    this.processing.set(true);
    this.processingMessage.set('Calculating payroll for all staff...');

    this.payrollService.calculatePayroll(run.id).subscribe({
      next: () => {
        this.toastService.success('Payroll calculated successfully');
        this.loadPayRuns();
        this.processing.set(false);
      },
      error: err => {
        const message = err?.error?.message || 'Failed to calculate payroll';
        this.toastService.error(message);
        this.processing.set(false);
      },
    });
  }

  approvePayRun(run: PayRun): void {
    this.processing.set(true);
    this.processingMessage.set('Approving pay run...');

    this.payrollService.approvePayRun(run.id).subscribe({
      next: () => {
        this.toastService.success('Pay run approved successfully');
        this.loadPayRuns();
        this.processing.set(false);
      },
      error: err => {
        const message = err?.error?.message || 'Failed to approve pay run';
        this.toastService.error(message);
        this.processing.set(false);
      },
    });
  }

  finalizePayRun(run: PayRun): void {
    this.processing.set(true);
    this.processingMessage.set('Finalizing pay run...');

    this.payrollService.finalizePayRun(run.id).subscribe({
      next: () => {
        this.toastService.success('Pay run finalized successfully');
        this.loadPayRuns();
        this.processing.set(false);
      },
      error: err => {
        const message = err?.error?.message || 'Failed to finalize pay run';
        this.toastService.error(message);
        this.processing.set(false);
      },
    });
  }

  exportBankFile(run: PayRun): void {
    this.payrollService.exportBankFile(run.id).subscribe({
      next: blob => {
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `bank-transfer-${run.payPeriodYear}-${run.payPeriodMonth}.csv`;
        a.click();
        window.URL.revokeObjectURL(url);
        this.toastService.success('Bank file downloaded');
      },
      error: () => {
        this.toastService.error('Failed to download bank file');
      },
    });
  }

  confirmDelete(run: PayRun): void {
    this.payRunToDelete.set(run);
    this.showDeleteModal.set(true);
  }

  closeDeleteModal(): void {
    this.showDeleteModal.set(false);
    this.payRunToDelete.set(null);
  }

  deletePayRun(): void {
    const run = this.payRunToDelete();
    if (!run) return;

    this.deleting.set(true);

    this.payrollService.deletePayRun(run.id).subscribe({
      next: () => {
        this.toastService.success('Pay run deleted successfully');
        this.closeDeleteModal();
        this.loadPayRuns();
        this.deleting.set(false);
      },
      error: () => {
        this.toastService.error('Failed to delete pay run');
        this.deleting.set(false);
      },
    });
  }

  getStatusClass(status: PayRunStatus): string {
    switch (status) {
      case 'draft':
        return 'gray';
      case 'processing':
        return 'blue';
      case 'calculated':
        return 'yellow';
      case 'approved':
        return 'blue';
      case 'finalized':
        return 'green';
      case 'reversed':
        return 'red';
      default:
        return 'gray';
    }
  }

  formatPayPeriod = formatPayPeriod;
  formatCurrency = formatCurrency;
  getPayRunStatusLabel = getPayRunStatusLabel;

  formatDate(dateStr: string): string {
    if (!dateStr) return '-';
    return new Date(dateStr).toLocaleDateString('en-IN', {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
    });
  }
}
