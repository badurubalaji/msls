/**
 * MSLS Academic Years Management Component
 *
 * Main component for managing academic years with terms and holidays.
 * Features expandable rows to show terms and holidays inline.
 */

import { Component, OnInit, inject, signal, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { MslsModalComponent } from '../../../shared/components';
import { ToastService } from '../../../shared/services';

import { AcademicYearService } from './academic-year.service';
import { AcademicYear, AcademicTerm, Holiday, HolidayType, HOLIDAY_TYPES } from './academic-year.model';
import { AcademicYearFormComponent } from './academic-year-form.component';
import { TermFormComponent } from './term-form.component';
import { HolidayFormComponent } from './holiday-form.component';

@Component({
  selector: 'msls-academic-years',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MslsModalComponent,
    AcademicYearFormComponent,
    TermFormComponent,
    HolidayFormComponent,
  ],
  template: `
    <div class="academic-years-page">
      <div class="academic-years-card">
        <!-- Header -->
        <div class="page-header">
          <div class="page-header__left">
            <h1 class="page-header__title">Academic Year Management</h1>
            <p class="page-header__subtitle">
              Configure academic years, terms, and holidays for your institution
            </p>
          </div>
          <div class="page-header__right">
            <div class="search-input">
              <i class="fa-solid fa-magnifying-glass search-icon"></i>
              <input
                type="text"
                placeholder="Search academic years..."
                [ngModel]="searchTerm()"
                (ngModelChange)="onSearchChange($event)"
                class="search-field"
              />
            </div>
            <button class="btn btn-primary" (click)="openCreateModal()">
              <i class="fa-solid fa-plus"></i>
              Add Academic Year
            </button>
          </div>
        </div>

        <!-- Loading State -->
        @if (loading()) {
          <div class="loading-container">
            <div class="spinner"></div>
            <p>Loading academic years...</p>
          </div>
        } @else if (error()) {
          <div class="error-container">
            <i class="fa-solid fa-circle-exclamation"></i>
            <p>{{ error() }}</p>
            <button class="btn btn-secondary" (click)="loadAcademicYears()">
              Retry
            </button>
          </div>
        } @else {
          <!-- Table -->
          <div class="table-container">
            <table class="data-table">
              <thead>
                <tr>
                  <th style="width: 40px;"></th>
                  <th>Name</th>
                  <th>Start Date</th>
                  <th>End Date</th>
                  <th style="width: 100px;">Status</th>
                  <th style="width: 100px;">Current</th>
                  <th style="width: 180px; text-align: right;">Actions</th>
                </tr>
              </thead>
              <tbody>
                @for (year of filteredAcademicYears(); track year.id) {
                  <!-- Main Row -->
                  <tr
                    class="main-row"
                    [class.expanded]="isExpanded(year.id)"
                    (click)="toggleExpand(year.id)"
                  >
                    <td class="expand-cell">
                      <button class="expand-btn" [class.rotated]="isExpanded(year.id)">
                        <i class="fa-solid fa-chevron-right"></i>
                      </button>
                    </td>
                    <td class="name-cell">{{ year.name }}</td>
                    <td>{{ formatDate(year.startDate) }}</td>
                    <td>{{ formatDate(year.endDate) }}</td>
                    <td>
                      <span
                        class="badge"
                        [class.badge-green]="year.isActive"
                        [class.badge-gray]="!year.isActive"
                      >
                        {{ year.isActive ? 'Active' : 'Inactive' }}
                      </span>
                    </td>
                    <td>
                      @if (year.isCurrent) {
                        <span class="badge badge-blue">
                          <i class="fa-solid fa-star" style="font-size: 0.625rem; margin-right: 0.25rem;"></i>
                          Current
                        </span>
                      } @else {
                        <span class="badge badge-gray">-</span>
                      }
                    </td>
                    <td class="actions-cell">
                      @if (!year.isCurrent) {
                        <button
                          class="action-btn action-btn--primary"
                          (click)="setAsCurrent(year); $event.stopPropagation()"
                          title="Set as current year"
                        >
                          <i class="fa-solid fa-star"></i>
                        </button>
                      }
                      <button
                        class="action-btn"
                        (click)="editAcademicYear(year); $event.stopPropagation()"
                        title="Edit academic year"
                      >
                        <i class="fa-regular fa-pen-to-square"></i>
                      </button>
                      <button
                        class="action-btn action-btn--danger"
                        (click)="confirmDelete(year); $event.stopPropagation()"
                        title="Delete academic year"
                      >
                        <i class="fa-regular fa-trash-can"></i>
                      </button>
                    </td>
                  </tr>

                  <!-- Expanded Content -->
                  @if (isExpanded(year.id)) {
                    <tr class="expanded-row">
                      <td colspan="7">
                        <div class="expanded-content">
                          <!-- Terms Section -->
                          <div class="section">
                            <div class="section-header">
                              <h3 class="section-title">
                                <i class="fa-solid fa-calendar-days"></i>
                                Terms / Semesters
                              </h3>
                              <button
                                class="btn btn-sm btn-secondary"
                                (click)="openTermModal(year); $event.stopPropagation()"
                              >
                                <i class="fa-solid fa-plus"></i>
                                Add Term
                              </button>
                            </div>
                            @if (year.terms && year.terms.length > 0) {
                              <div class="mini-table-container">
                                <table class="mini-table">
                                  <thead>
                                    <tr>
                                      <th>Seq</th>
                                      <th>Name</th>
                                      <th>Start Date</th>
                                      <th>End Date</th>
                                      <th style="width: 100px; text-align: right;">Actions</th>
                                    </tr>
                                  </thead>
                                  <tbody>
                                    @for (term of sortedTerms(year.terms); track term.id) {
                                      <tr>
                                        <td class="seq-cell">{{ term.sequence }}</td>
                                        <td>{{ term.name }}</td>
                                        <td>{{ formatDate(term.startDate) }}</td>
                                        <td>{{ formatDate(term.endDate) }}</td>
                                        <td class="actions-cell">
                                          <button
                                            class="action-btn action-btn--sm"
                                            (click)="editTerm(year, term); $event.stopPropagation()"
                                            title="Edit term"
                                          >
                                            <i class="fa-regular fa-pen-to-square"></i>
                                          </button>
                                          <button
                                            class="action-btn action-btn--sm action-btn--danger"
                                            (click)="confirmDeleteTerm(year, term); $event.stopPropagation()"
                                            title="Delete term"
                                          >
                                            <i class="fa-regular fa-trash-can"></i>
                                          </button>
                                        </td>
                                      </tr>
                                    }
                                  </tbody>
                                </table>
                              </div>
                            } @else {
                              <div class="empty-section">
                                <i class="fa-regular fa-calendar"></i>
                                <p>No terms defined yet</p>
                              </div>
                            }
                          </div>

                          <!-- Holidays Section -->
                          <div class="section">
                            <div class="section-header">
                              <h3 class="section-title">
                                <i class="fa-solid fa-umbrella-beach"></i>
                                Holidays
                              </h3>
                              <button
                                class="btn btn-sm btn-secondary"
                                (click)="openHolidayModal(year); $event.stopPropagation()"
                              >
                                <i class="fa-solid fa-plus"></i>
                                Add Holiday
                              </button>
                            </div>
                            @if (year.holidays && year.holidays.length > 0) {
                              <div class="mini-table-container">
                                <table class="mini-table">
                                  <thead>
                                    <tr>
                                      <th>Name</th>
                                      <th>Date</th>
                                      <th>Type</th>
                                      <th style="width: 80px;">Optional</th>
                                      <th style="width: 100px; text-align: right;">Actions</th>
                                    </tr>
                                  </thead>
                                  <tbody>
                                    @for (holiday of sortedHolidays(year.holidays); track holiday.id) {
                                      <tr>
                                        <td>{{ holiday.name }}</td>
                                        <td>{{ formatDate(holiday.date) }}</td>
                                        <td>
                                          <span class="badge badge-light">
                                            {{ getHolidayTypeLabel(holiday.type) }}
                                          </span>
                                        </td>
                                        <td>
                                          @if (holiday.isOptional) {
                                            <span class="badge badge-yellow">Optional</span>
                                          } @else {
                                            <span class="badge badge-gray">-</span>
                                          }
                                        </td>
                                        <td class="actions-cell">
                                          <button
                                            class="action-btn action-btn--sm"
                                            (click)="editHoliday(year, holiday); $event.stopPropagation()"
                                            title="Edit holiday"
                                          >
                                            <i class="fa-regular fa-pen-to-square"></i>
                                          </button>
                                          <button
                                            class="action-btn action-btn--sm action-btn--danger"
                                            (click)="confirmDeleteHoliday(year, holiday); $event.stopPropagation()"
                                            title="Delete holiday"
                                          >
                                            <i class="fa-regular fa-trash-can"></i>
                                          </button>
                                        </td>
                                      </tr>
                                    }
                                  </tbody>
                                </table>
                              </div>
                            } @else {
                              <div class="empty-section">
                                <i class="fa-regular fa-calendar-xmark"></i>
                                <p>No holidays defined yet</p>
                              </div>
                            }
                          </div>
                        </div>
                      </td>
                    </tr>
                  }
                } @empty {
                  <tr>
                    <td colspan="7" class="empty-cell">
                      <div class="empty-state">
                        <i class="fa-regular fa-folder-open"></i>
                        <p>No academic years found</p>
                        <button class="btn btn-primary btn-sm" (click)="openCreateModal()">
                          <i class="fa-solid fa-plus"></i>
                          Add Academic Year
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

      <!-- Create/Edit Academic Year Modal -->
      <msls-modal
        [isOpen]="showYearModal()"
        [title]="editingYear() ? 'Edit Academic Year' : 'Create Academic Year'"
        size="md"
        (closed)="closeYearModal()"
      >
        <msls-academic-year-form
          [academicYear]="editingYear()"
          [loading]="saving()"
          (save)="saveAcademicYear($event)"
          (cancel)="closeYearModal()"
        />
      </msls-modal>

      <!-- Term Modal -->
      <msls-modal
        [isOpen]="showTermModal()"
        [title]="editingTerm() ? 'Edit Term' : 'Add Term'"
        size="md"
        (closed)="closeTermModal()"
      >
        <msls-term-form
          [term]="editingTerm()"
          [academicYear]="selectedYear()"
          [loading]="savingTerm()"
          (save)="saveTerm($event)"
          (cancel)="closeTermModal()"
        />
      </msls-modal>

      <!-- Holiday Modal -->
      <msls-modal
        [isOpen]="showHolidayModal()"
        [title]="editingHoliday() ? 'Edit Holiday' : 'Add Holiday'"
        size="md"
        (closed)="closeHolidayModal()"
      >
        <msls-holiday-form
          [holiday]="editingHoliday()"
          [academicYear]="selectedYear()"
          [loading]="savingHoliday()"
          (save)="saveHoliday($event)"
          (cancel)="closeHolidayModal()"
        />
      </msls-modal>

      <!-- Delete Academic Year Modal -->
      <msls-modal
        [isOpen]="showDeleteModal()"
        title="Delete Academic Year"
        size="sm"
        (closed)="closeDeleteModal()"
      >
        <div class="delete-confirmation">
          <div class="delete-icon">
            <i class="fa-solid fa-triangle-exclamation"></i>
          </div>
          <p>
            Are you sure you want to delete the academic year
            <strong>"{{ yearToDelete()?.name }}"</strong>?
          </p>
          <p class="delete-warning">
            This will also delete all associated terms and holidays. This action cannot be undone.
          </p>
          <div class="delete-actions">
            <button class="btn btn-secondary" (click)="closeDeleteModal()">
              Cancel
            </button>
            <button
              class="btn btn-danger"
              [disabled]="deleting()"
              (click)="deleteAcademicYear()"
            >
              @if (deleting()) {
                <div class="btn-spinner"></div>
                Deleting...
              } @else {
                Delete
              }
            </button>
          </div>
        </div>
      </msls-modal>

      <!-- Delete Term Modal -->
      <msls-modal
        [isOpen]="showDeleteTermModal()"
        title="Delete Term"
        size="sm"
        (closed)="closeDeleteTermModal()"
      >
        <div class="delete-confirmation">
          <div class="delete-icon">
            <i class="fa-solid fa-triangle-exclamation"></i>
          </div>
          <p>
            Are you sure you want to delete the term
            <strong>"{{ termToDelete()?.name }}"</strong>?
          </p>
          <p class="delete-warning">This action cannot be undone.</p>
          <div class="delete-actions">
            <button class="btn btn-secondary" (click)="closeDeleteTermModal()">
              Cancel
            </button>
            <button
              class="btn btn-danger"
              [disabled]="deletingTerm()"
              (click)="deleteTerm()"
            >
              @if (deletingTerm()) {
                <div class="btn-spinner"></div>
                Deleting...
              } @else {
                Delete
              }
            </button>
          </div>
        </div>
      </msls-modal>

      <!-- Delete Holiday Modal -->
      <msls-modal
        [isOpen]="showDeleteHolidayModal()"
        title="Delete Holiday"
        size="sm"
        (closed)="closeDeleteHolidayModal()"
      >
        <div class="delete-confirmation">
          <div class="delete-icon">
            <i class="fa-solid fa-triangle-exclamation"></i>
          </div>
          <p>
            Are you sure you want to delete the holiday
            <strong>"{{ holidayToDelete()?.name }}"</strong>?
          </p>
          <p class="delete-warning">This action cannot be undone.</p>
          <div class="delete-actions">
            <button class="btn btn-secondary" (click)="closeDeleteHolidayModal()">
              Cancel
            </button>
            <button
              class="btn btn-danger"
              [disabled]="deletingHoliday()"
              (click)="deleteHoliday()"
            >
              @if (deletingHoliday()) {
                <div class="btn-spinner"></div>
                Deleting...
              } @else {
                Delete
              }
            </button>
          </div>
        </div>
      </msls-modal>
    </div>
  `,
  styles: [`
    .academic-years-page {
      padding: 1.5rem;
      max-width: 1400px;
      margin: 0 auto;
    }

    .academic-years-card {
      background: #ffffff;
      border: 1px solid #e2e8f0;
      border-radius: 1rem;
      padding: 1.5rem;
    }

    /* Header */
    .page-header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: 1.5rem;
      padding-bottom: 1.5rem;
      border-bottom: 1px solid #e2e8f0;
      flex-wrap: wrap;
      gap: 1rem;
    }

    .page-header__title {
      font-size: 1.5rem;
      font-weight: 700;
      color: #0f172a;
      margin: 0 0 0.375rem 0;
    }

    .page-header__subtitle {
      font-size: 0.875rem;
      color: #64748b;
      margin: 0;
    }

    .page-header__right {
      display: flex;
      gap: 0.75rem;
      align-items: center;
    }

    /* Search Input */
    .search-input {
      position: relative;
    }

    .search-icon {
      position: absolute;
      left: 0.875rem;
      top: 50%;
      transform: translateY(-50%);
      color: #94a3b8;
      font-size: 0.875rem;
    }

    .search-field {
      padding: 0.625rem 0.875rem 0.625rem 2.5rem;
      font-size: 0.875rem;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      background: #ffffff;
      color: #0f172a;
      width: 220px;
      transition: all 0.15s;
    }

    .search-field::placeholder {
      color: #94a3b8;
    }

    .search-field:focus {
      outline: none;
      border-color: #4f46e5;
      box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
    }

    /* Buttons */
    .btn {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      gap: 0.5rem;
      padding: 0.625rem 1rem;
      font-size: 0.875rem;
      font-weight: 500;
      border-radius: 0.5rem;
      border: none;
      cursor: pointer;
      transition: all 0.15s;
    }

    .btn-sm {
      padding: 0.375rem 0.75rem;
      font-size: 0.8125rem;
    }

    .btn-primary {
      background: #4f46e5;
      color: #ffffff;
    }

    .btn-primary:hover {
      background: #4338ca;
    }

    .btn-secondary {
      background: #ffffff;
      color: #334155;
      border: 1px solid #e2e8f0;
    }

    .btn-secondary:hover {
      background: #f8fafc;
      border-color: #cbd5e1;
    }

    .btn-danger {
      background: #dc2626;
      color: #ffffff;
    }

    .btn-danger:hover:not(:disabled) {
      background: #b91c1c;
    }

    .btn-danger:disabled {
      opacity: 0.6;
      cursor: not-allowed;
    }

    .btn-spinner {
      width: 1rem;
      height: 1rem;
      border: 2px solid rgba(255, 255, 255, 0.3);
      border-top-color: #ffffff;
      border-radius: 50%;
      animation: spin 0.6s linear infinite;
    }

    /* Loading */
    .loading-container {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 4rem;
      gap: 1rem;
    }

    .spinner {
      width: 2rem;
      height: 2rem;
      border: 3px solid #e2e8f0;
      border-top-color: #4f46e5;
      border-radius: 50%;
      animation: spin 0.8s linear infinite;
    }

    @keyframes spin {
      to { transform: rotate(360deg); }
    }

    .loading-container p {
      color: #64748b;
      font-size: 0.875rem;
      margin: 0;
    }

    /* Error */
    .error-container {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 4rem;
      gap: 1rem;
    }

    .error-container i {
      font-size: 2.5rem;
      color: #dc2626;
    }

    .error-container p {
      color: #64748b;
      font-size: 0.875rem;
      margin: 0;
    }

    /* Table */
    .table-container {
      border: 1px solid #e2e8f0;
      border-radius: 0.75rem;
      overflow: hidden;
    }

    .data-table {
      width: 100%;
      border-collapse: collapse;
    }

    .data-table thead {
      background: #f8fafc;
      border-bottom: 1px solid #e2e8f0;
    }

    .data-table th {
      padding: 0.75rem 1rem;
      text-align: left;
      font-size: 0.75rem;
      font-weight: 600;
      color: #64748b;
      text-transform: uppercase;
      letter-spacing: 0.05em;
    }

    .data-table td {
      padding: 1rem;
      font-size: 0.875rem;
      border-bottom: 1px solid #f1f5f9;
      color: #334155;
    }

    .data-table tbody tr:last-child td {
      border-bottom: none;
    }

    .main-row {
      cursor: pointer;
      transition: background 0.15s;
    }

    .main-row:hover {
      background: #f8fafc;
    }

    .main-row.expanded {
      background: #f1f5f9;
    }

    .expand-cell {
      padding: 1rem 0.5rem 1rem 1rem !important;
    }

    .expand-btn {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 1.5rem;
      height: 1.5rem;
      background: transparent;
      border: none;
      color: #64748b;
      cursor: pointer;
      transition: transform 0.2s;
    }

    .expand-btn.rotated {
      transform: rotate(90deg);
    }

    .expand-btn i {
      font-size: 0.75rem;
    }

    .name-cell {
      font-weight: 500;
      color: #0f172a;
    }

    /* Badges */
    .badge {
      display: inline-flex;
      align-items: center;
      padding: 0.25rem 0.625rem;
      font-size: 0.75rem;
      font-weight: 500;
      border-radius: 9999px;
    }

    .badge-green {
      background: #dcfce7;
      color: #166534;
    }

    .badge-blue {
      background: #dbeafe;
      color: #1e40af;
    }

    .badge-gray {
      background: #f1f5f9;
      color: #475569;
    }

    .badge-yellow {
      background: #fef3c7;
      color: #92400e;
    }

    .badge-light {
      background: #f1f5f9;
      color: #334155;
    }

    /* Actions */
    .actions-cell {
      display: flex;
      justify-content: flex-end;
      gap: 0.5rem;
    }

    .action-btn {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 2rem;
      height: 2rem;
      background: transparent;
      border: 1px solid #e2e8f0;
      border-radius: 0.375rem;
      color: #64748b;
      cursor: pointer;
      transition: all 0.15s;
    }

    .action-btn:hover {
      background: #f8fafc;
      border-color: #cbd5e1;
      color: #0f172a;
    }

    .action-btn--sm {
      width: 1.75rem;
      height: 1.75rem;
    }

    .action-btn--sm i {
      font-size: 0.75rem;
    }

    .action-btn--primary {
      border-color: #c7d2fe;
      color: #4f46e5;
    }

    .action-btn--primary:hover {
      background: #eef2ff;
      border-color: #a5b4fc;
      color: #4338ca;
    }

    .action-btn--danger:hover {
      background: #fef2f2;
      border-color: #fecaca;
      color: #dc2626;
    }

    /* Expanded Row */
    .expanded-row {
      background: #fafbfc;
    }

    .expanded-row td {
      padding: 0 !important;
      border-bottom: 1px solid #e2e8f0;
    }

    .expanded-content {
      padding: 1.5rem;
      display: flex;
      flex-direction: column;
      gap: 1.5rem;
    }

    /* Section */
    .section {
      background: #ffffff;
      border: 1px solid #e2e8f0;
      border-radius: 0.5rem;
      overflow: hidden;
    }

    .section-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 0.75rem 1rem;
      background: #f8fafc;
      border-bottom: 1px solid #e2e8f0;
    }

    .section-title {
      font-size: 0.875rem;
      font-weight: 600;
      color: #334155;
      margin: 0;
      display: flex;
      align-items: center;
      gap: 0.5rem;
    }

    .section-title i {
      color: #64748b;
    }

    /* Mini Table */
    .mini-table-container {
      overflow-x: auto;
    }

    .mini-table {
      width: 100%;
      border-collapse: collapse;
    }

    .mini-table th {
      padding: 0.625rem 1rem;
      text-align: left;
      font-size: 0.6875rem;
      font-weight: 600;
      color: #64748b;
      text-transform: uppercase;
      letter-spacing: 0.05em;
      background: #fafbfc;
      border-bottom: 1px solid #e2e8f0;
    }

    .mini-table td {
      padding: 0.625rem 1rem;
      font-size: 0.8125rem;
      color: #334155;
      border-bottom: 1px solid #f1f5f9;
    }

    .mini-table tbody tr:last-child td {
      border-bottom: none;
    }

    .seq-cell {
      font-weight: 600;
      color: #4f46e5;
    }

    /* Empty Section */
    .empty-section {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 2rem;
      color: #94a3b8;
    }

    .empty-section i {
      font-size: 1.5rem;
      margin-bottom: 0.5rem;
    }

    .empty-section p {
      margin: 0;
      font-size: 0.8125rem;
    }

    /* Empty State */
    .empty-cell {
      padding: 3rem !important;
    }

    .empty-state {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 0.75rem;
      color: #94a3b8;
    }

    .empty-state i {
      font-size: 2.5rem;
    }

    .empty-state p {
      margin: 0;
      font-size: 0.875rem;
    }

    /* Delete Confirmation */
    .delete-confirmation {
      text-align: center;
      padding: 1rem;
    }

    .delete-icon {
      width: 3.5rem;
      height: 3.5rem;
      margin: 0 auto 1rem;
      background: #fef2f2;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
    }

    .delete-icon i {
      font-size: 1.5rem;
      color: #dc2626;
    }

    .delete-confirmation p {
      color: #475569;
      margin: 0 0 0.5rem 0;
    }

    .delete-confirmation strong {
      color: #dc2626;
    }

    .delete-warning {
      color: #dc2626 !important;
      font-size: 0.8125rem;
      font-weight: 500;
      padding: 0.75rem;
      background: #fef2f2;
      border-radius: 0.5rem;
      margin-top: 1rem !important;
    }

    .delete-actions {
      display: flex;
      justify-content: center;
      gap: 0.75rem;
      margin-top: 1.5rem;
    }

    /* Responsive */
    @media (max-width: 768px) {
      .page-header {
        flex-direction: column;
        align-items: stretch;
      }

      .page-header__right {
        flex-direction: column;
      }

      .search-field {
        width: 100%;
      }

      .btn-primary {
        width: 100%;
        justify-content: center;
      }

      .table-container {
        overflow-x: auto;
      }

      .data-table {
        min-width: 700px;
      }
    }
  `],
})
export class AcademicYearsComponent implements OnInit {
  private academicYearService = inject(AcademicYearService);
  private toastService = inject(ToastService);

  // State signals
  academicYears = signal<AcademicYear[]>([]);
  loading = signal(true);
  saving = signal(false);
  savingTerm = signal(false);
  savingHoliday = signal(false);
  deleting = signal(false);
  deletingTerm = signal(false);
  deletingHoliday = signal(false);
  error = signal<string | null>(null);
  searchTerm = signal('');
  expandedRows = signal<Set<string>>(new Set());

  // Modal state - Academic Year
  showYearModal = signal(false);
  showDeleteModal = signal(false);
  editingYear = signal<AcademicYear | null>(null);
  yearToDelete = signal<AcademicYear | null>(null);

  // Modal state - Term
  showTermModal = signal(false);
  showDeleteTermModal = signal(false);
  editingTerm = signal<AcademicTerm | null>(null);
  termToDelete = signal<AcademicTerm | null>(null);
  selectedYear = signal<AcademicYear | null>(null);

  // Modal state - Holiday
  showHolidayModal = signal(false);
  showDeleteHolidayModal = signal(false);
  editingHoliday = signal<Holiday | null>(null);
  holidayToDelete = signal<Holiday | null>(null);

  // Computed filtered academic years
  filteredAcademicYears = computed(() => {
    const term = this.searchTerm().toLowerCase();
    if (!term) return this.academicYears();

    return this.academicYears().filter(
      year => year.name.toLowerCase().includes(term)
    );
  });

  ngOnInit(): void {
    this.loadAcademicYears();
  }

  // ============================================================================
  // Data Loading
  // ============================================================================

  loadAcademicYears(): void {
    this.loading.set(true);
    this.error.set(null);

    this.academicYearService.getAcademicYears().subscribe({
      next: years => {
        this.academicYears.set(years);
        this.loading.set(false);
      },
      error: err => {
        this.error.set('Failed to load academic years. Please try again.');
        this.loading.set(false);
        console.error('Failed to load academic years:', err);
      },
    });
  }

  // ============================================================================
  // Utility Methods
  // ============================================================================

  onSearchChange(term: string): void {
    this.searchTerm.set(term);
  }

  isExpanded(id: string): boolean {
    return this.expandedRows().has(id);
  }

  toggleExpand(id: string): void {
    const expanded = new Set(this.expandedRows());
    if (expanded.has(id)) {
      expanded.delete(id);
    } else {
      expanded.add(id);
    }
    this.expandedRows.set(expanded);
  }

  formatDate(dateStr: string): string {
    if (!dateStr) return '-';
    const date = new Date(dateStr);
    return date.toLocaleDateString('en-IN', {
      day: '2-digit',
      month: 'short',
      year: 'numeric',
    });
  }

  sortedTerms(terms: AcademicTerm[]): AcademicTerm[] {
    return [...terms].sort((a, b) => a.sequence - b.sequence);
  }

  sortedHolidays(holidays: Holiday[]): Holiday[] {
    return [...holidays].sort((a, b) => new Date(a.date).getTime() - new Date(b.date).getTime());
  }

  getHolidayTypeLabel(type: string): string {
    const found = HOLIDAY_TYPES.find(t => t.value === type);
    return found?.label || type;
  }

  // ============================================================================
  // Academic Year CRUD
  // ============================================================================

  openCreateModal(): void {
    this.editingYear.set(null);
    this.showYearModal.set(true);
  }

  editAcademicYear(year: AcademicYear): void {
    this.editingYear.set(year);
    this.showYearModal.set(true);
  }

  closeYearModal(): void {
    this.showYearModal.set(false);
    this.editingYear.set(null);
  }

  saveAcademicYear(data: { name: string; startDate: string; endDate: string; isCurrent?: boolean }): void {
    this.saving.set(true);

    const editing = this.editingYear();
    const operation = editing
      ? this.academicYearService.updateAcademicYear(editing.id, data)
      : this.academicYearService.createAcademicYear(data);

    operation.subscribe({
      next: () => {
        this.toastService.success(
          editing ? 'Academic year updated successfully' : 'Academic year created successfully'
        );
        this.closeYearModal();
        this.loadAcademicYears();
        this.saving.set(false);
      },
      error: err => {
        this.toastService.error(
          editing ? 'Failed to update academic year' : 'Failed to create academic year'
        );
        this.saving.set(false);
        console.error('Failed to save academic year:', err);
      },
    });
  }

  setAsCurrent(year: AcademicYear): void {
    this.academicYearService.setAsCurrent(year.id).subscribe({
      next: () => {
        this.toastService.success(`"${year.name}" is now the current academic year`);
        this.loadAcademicYears();
      },
      error: err => {
        this.toastService.error('Failed to set current academic year');
        console.error('Failed to set current year:', err);
      },
    });
  }

  confirmDelete(year: AcademicYear): void {
    this.yearToDelete.set(year);
    this.showDeleteModal.set(true);
  }

  closeDeleteModal(): void {
    this.showDeleteModal.set(false);
    this.yearToDelete.set(null);
  }

  deleteAcademicYear(): void {
    const year = this.yearToDelete();
    if (!year) return;

    this.deleting.set(true);

    this.academicYearService.deleteAcademicYear(year.id).subscribe({
      next: () => {
        this.toastService.success('Academic year deleted successfully');
        this.closeDeleteModal();
        this.loadAcademicYears();
        this.deleting.set(false);
      },
      error: err => {
        this.toastService.error('Failed to delete academic year');
        this.deleting.set(false);
        console.error('Failed to delete academic year:', err);
      },
    });
  }

  // ============================================================================
  // Term CRUD
  // ============================================================================

  openTermModal(year: AcademicYear): void {
    this.selectedYear.set(year);
    this.editingTerm.set(null);
    this.showTermModal.set(true);
  }

  editTerm(year: AcademicYear, term: AcademicTerm): void {
    this.selectedYear.set(year);
    this.editingTerm.set(term);
    this.showTermModal.set(true);
  }

  closeTermModal(): void {
    this.showTermModal.set(false);
    this.editingTerm.set(null);
    this.selectedYear.set(null);
  }

  saveTerm(data: { name: string; startDate: string; endDate: string; sequence?: number }): void {
    const year = this.selectedYear();
    if (!year) return;

    this.savingTerm.set(true);

    const editing = this.editingTerm();
    const operation = editing
      ? this.academicYearService.updateTerm(year.id, editing.id, data)
      : this.academicYearService.createTerm(year.id, data);

    operation.subscribe({
      next: () => {
        this.toastService.success(
          editing ? 'Term updated successfully' : 'Term added successfully'
        );
        this.closeTermModal();
        this.loadAcademicYears();
        this.savingTerm.set(false);
      },
      error: err => {
        this.toastService.error(
          editing ? 'Failed to update term' : 'Failed to add term'
        );
        this.savingTerm.set(false);
        console.error('Failed to save term:', err);
      },
    });
  }

  confirmDeleteTerm(year: AcademicYear, term: AcademicTerm): void {
    this.selectedYear.set(year);
    this.termToDelete.set(term);
    this.showDeleteTermModal.set(true);
  }

  closeDeleteTermModal(): void {
    this.showDeleteTermModal.set(false);
    this.termToDelete.set(null);
  }

  deleteTerm(): void {
    const year = this.selectedYear();
    const term = this.termToDelete();
    if (!year || !term) return;

    this.deletingTerm.set(true);

    this.academicYearService.deleteTerm(year.id, term.id).subscribe({
      next: () => {
        this.toastService.success('Term deleted successfully');
        this.closeDeleteTermModal();
        this.loadAcademicYears();
        this.deletingTerm.set(false);
      },
      error: err => {
        this.toastService.error('Failed to delete term');
        this.deletingTerm.set(false);
        console.error('Failed to delete term:', err);
      },
    });
  }

  // ============================================================================
  // Holiday CRUD
  // ============================================================================

  openHolidayModal(year: AcademicYear): void {
    this.selectedYear.set(year);
    this.editingHoliday.set(null);
    this.showHolidayModal.set(true);
  }

  editHoliday(year: AcademicYear, holiday: Holiday): void {
    this.selectedYear.set(year);
    this.editingHoliday.set(holiday);
    this.showHolidayModal.set(true);
  }

  closeHolidayModal(): void {
    this.showHolidayModal.set(false);
    this.editingHoliday.set(null);
    this.selectedYear.set(null);
  }

  saveHoliday(data: { name: string; date: string; type?: HolidayType; isOptional?: boolean }): void {
    const year = this.selectedYear();
    if (!year) return;

    this.savingHoliday.set(true);

    const editing = this.editingHoliday();
    const operation = editing
      ? this.academicYearService.updateHoliday(year.id, editing.id, data)
      : this.academicYearService.createHoliday(year.id, data);

    operation.subscribe({
      next: () => {
        this.toastService.success(
          editing ? 'Holiday updated successfully' : 'Holiday added successfully'
        );
        this.closeHolidayModal();
        this.loadAcademicYears();
        this.savingHoliday.set(false);
      },
      error: err => {
        this.toastService.error(
          editing ? 'Failed to update holiday' : 'Failed to add holiday'
        );
        this.savingHoliday.set(false);
        console.error('Failed to save holiday:', err);
      },
    });
  }

  confirmDeleteHoliday(year: AcademicYear, holiday: Holiday): void {
    this.selectedYear.set(year);
    this.holidayToDelete.set(holiday);
    this.showDeleteHolidayModal.set(true);
  }

  closeDeleteHolidayModal(): void {
    this.showDeleteHolidayModal.set(false);
    this.holidayToDelete.set(null);
  }

  deleteHoliday(): void {
    const year = this.selectedYear();
    const holiday = this.holidayToDelete();
    if (!year || !holiday) return;

    this.deletingHoliday.set(true);

    this.academicYearService.deleteHoliday(year.id, holiday.id).subscribe({
      next: () => {
        this.toastService.success('Holiday deleted successfully');
        this.closeDeleteHolidayModal();
        this.loadAcademicYears();
        this.deletingHoliday.set(false);
      },
      error: err => {
        this.toastService.error('Failed to delete holiday');
        this.deletingHoliday.set(false);
        console.error('Failed to delete holiday:', err);
      },
    });
  }
}
